package dcheck_test

import (
	"os"
	"testing"

	dcheck "github.com/pt9912/d-check"
)

// Die eingebetteten Dokumente sind byte-gleich zu ihren Quelldateien
// (DC-FA-CLI-013).
func TestEingebetteteDokumenteGleichDerQuelle(t *testing.T) {
	for pfad, eingebettet := range map[string]string{
		dcheck.HandbuchPfad:      dcheck.Handbuch,
		dcheck.SpezifikationPfad: dcheck.Spezifikation,
	} {
		quelle, err := os.ReadFile(pfad)
		if err != nil {
			t.Fatalf("%s: %v", pfad, err)
		}
		if string(quelle) != eingebettet {
			t.Errorf("%s: eingebettete Fassung weicht von der Quelldatei ab", pfad)
		}
	}
}
