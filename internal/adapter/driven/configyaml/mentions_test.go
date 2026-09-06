package configyaml_test

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driven/configyaml"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// Der Voll-Block wird gelesen (DC-FA-MENT-001, DC-FA-CONF-001).
func TestDecode_MentionsVollblock(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(
		"mentions:\n  artifacts: [tools/*.sh]\n  documents: [docs/handbuch.md]\n  match: basename\n"))
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(cfg.Mentions.Artifacts) != 1 || len(cfg.Mentions.Documents) != 1 {
		t.Fatalf("Mengen nicht gelesen: %+v", cfg.Mentions)
	}
	if cfg.Mentions.EffectiveMatch() != model.MentionsMatchBasename {
		t.Fatalf("match nicht gelesen: %q", cfg.Mentions.EffectiveMatch())
	}
}

// Ohne Block bleibt die Konfiguration leer — das Modul ist dann inert im Sinne
// von "es gibt nichts zu tun"; der Abbruch bei AKTIVEM Modul faellt in der
// Regel selbst (rules.CheckMentions), nicht hier.
func TestDecode_MentionsOhneBlock(t *testing.T) {
	cfg, err := configyaml.Decode([]byte("modules: [links]\n"))
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(cfg.Mentions.Artifacts) != 0 || len(cfg.Mentions.Documents) != 0 {
		t.Fatalf("erwartet leer, got %+v", cfg.Mentions)
	}
}

// Ein halber Block ist ein Fehler und kein halber Lauf: die Deckungs-Frage
// braucht beide Seiten (DC-FA-MENT-001, Exit 2). Beide Richtungen geprueft.
func TestDecode_MentionsHalberBlockFailClosed(t *testing.T) {
	for name, yaml := range map[string]string{
		"nur artifacts": "mentions:\n  artifacts: [tools/*.sh]\n",
		"nur documents": "mentions:\n  documents: [docs/a.md]\n",
	} {
		_, err := configyaml.Decode([]byte(yaml))
		if err == nil {
			t.Fatalf("%s: erwartet Fehler, got nil", name)
		}
		if !strings.Contains(err.Error(), "mentions.") {
			t.Fatalf("%s: Meldung nennt den Schluessel nicht: %v", name, err)
		}
	}
}

// Ein ungueltiges Glob und ein leeres Glob brechen am Config-Rand ab, nicht
// erst im Lauf (DC-FA-CONF-001).
func TestDecode_MentionsGlobValidierung(t *testing.T) {
	for name, yaml := range map[string]string{
		"ungueltig":  "mentions:\n  artifacts: ['[unclosed']\n  documents: [docs/a.md]\n",
		"leer":       "mentions:\n  artifacts: ['']\n  documents: [docs/a.md]\n",
		"ist-seite":  "mentions:\n  artifacts: [tools/*.sh]\n  documents: ['[unclosed']\n",
	} {
		if _, err := configyaml.Decode([]byte(yaml)); err == nil {
			t.Fatalf("%s: erwartet Fehler, got nil", name)
		}
	}
}

// Ein unbekannter match-Wert ist ein Konfigurations-Fehler, kein stiller
// Rueckfall auf den Default — sonst pruefte der Lauf eine andere Form als die
// deklarierte (DC-FA-MENT-001).
func TestDecode_MentionsUnbekannterMatch(t *testing.T) {
	_, err := configyaml.Decode([]byte(
		"mentions:\n  artifacts: [tools/*.sh]\n  documents: [docs/a.md]\n  match: suffix\n"))
	if err == nil || !strings.Contains(err.Error(), "match") {
		t.Fatalf("erwartet match-Fehler, got %v", err)
	}
}
