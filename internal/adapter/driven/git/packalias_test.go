package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5/osfs"
)

// Weißbox-Test gegen packAliases() selbst (R1-F-2 aus dem Review dieses
// Slice): zwei Dateien mit identischem Hash-Suffix, aber unterschiedlichem
// Fremdpräfix, dürfen nicht von Gos randomisierter Map-Iteration abhängen —
// die lexikografisch kleinste gewinnt, über mehrere Läufe stabil.
func TestPackAliasesKollisionIstDeterministisch(t *testing.T) {
	dir := t.TempDir()
	packDirAbs := filepath.Join(dir, "objects", "pack")
	if err := os.MkdirAll(packDirAbs, 0o755); err != nil {
		t.Fatal(err)
	}
	hash := "0123456789abcdef0123456789abcdef01234567"
	for _, stem := range []string{"beta-" + hash, "alpha-" + hash} {
		if err := os.WriteFile(filepath.Join(packDirAbs, stem+packExt), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(packDirAbs, stem+idxExt), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	fs, ok := newPackAliasFS(osfs.New(dir)).(*packAliasFS)
	if !ok {
		t.Fatal("newPackAliasFS lieferte keinen *packAliasFS")
	}
	wantPack := "alpha-" + hash + packExt
	wantIdx := "alpha-" + hash + idxExt
	for i := 0; i < 5; i++ {
		aliases := fs.packAliases()
		if got := aliases[packPrefix+hash+packExt]; got != wantPack {
			t.Fatalf("Lauf %d: kanonischer Pack-Alias = %q, erwartet %q (lexikografisch kleinster Name)", i, got, wantPack)
		}
		if got := aliases[packPrefix+hash+idxExt]; got != wantIdx {
			t.Fatalf("Lauf %d: kanonischer Idx-Alias = %q, erwartet %q", i, got, wantIdx)
		}
	}
}
