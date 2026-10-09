package cli_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	dcheck "github.com/pt9912/d-check"
)

// Die eingebetteten Dokumente sind byte-gleich zu ihren Quelldateien
// (DC-FA-CLI-013).
func TestManual_EingebetteteDokumenteGleichDerQuelle(t *testing.T) {
	for pfad, eingebettet := range map[string]string{
		dcheck.HandbuchPfad:      dcheck.Handbuch,
		dcheck.SpezifikationPfad: dcheck.Spezifikation,
	} {
		quelle, err := os.ReadFile(filepath.Join("..", "..", "..", "..", pfad))
		if err != nil {
			t.Fatalf("%s: %v", pfad, err)
		}
		if string(quelle) != eingebettet {
			t.Errorf("%s: eingebettete Fassung weicht von der Quelldatei ab", pfad)
		}
	}
}

// Nur die Composition Root liest das Paket im Modul-Root (ADR-0107): kein
// Paket unter internal/hexagon oder internal/adapter/driven importiert es.
// a-check löst den Importpfad des Modul-Roots nicht auf, deshalb hält dieser
// Test die Kante.
func TestManual_RootPaketNurAusDerCompositionRoot(t *testing.T) {
	root := filepath.Join("..", "..", "..", "..")
	for _, dir := range []string{"internal/hexagon", "internal/adapter/driven"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(p, ".go") {
				return err
			}
			f, err := parser.ParseFile(token.NewFileSet(), p, nil, parser.ImportsOnly)
			if err != nil {
				return err
			}
			for _, imp := range f.Imports {
				if pfad, _ := strconv.Unquote(imp.Path.Value); pfad == "github.com/pt9912/d-check" {
					t.Errorf("%s importiert das Paket im Modul-Root", p)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}
