package configyaml_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

// Repo-Artefakt-Invarianten der Durchsetzungsschicht (MR-047, slice-157).
// Gegenstand ist NICHT d-checks eigene Konfiguration, sondern `.claude/` —
// dieselbe Familie wie gate_consistency_test.go: eine committete Datei, die
// eine Zusage der Dokumentation trägt, gegen ihren Inhalt gehalten. Sie liegt
// hier statt in einem eigenen Paket, weil eine Invariante kein Paket
// rechtfertigt und dieses die Repo-Datei-Invarianten bereits führt.
//
// WAS DIESE TESTS ZUSAGEN: dass die committeten Dateien wohlgeformt sind und
// aufeinander zeigen.
//
// WAS SIE NICHT ZUSAGEN — und das ist der wichtigere Satz: dass die
// Durchsetzung LÄUFT. Ob das Werkzeug den Hook aufruft und seine Antwort
// befolgt, steht in keiner Datei; AGENTS.md §3.1 sagt ausdrücklich, ein Lauf
// ohne dieses Werkzeug sei ungebunden. Ein grüner Test hier ist kein Beleg für
// einen scharfen Wächter — er ist der Ausschluss EINES Ausfallgrunds.

func repoRoot() string { return "../../../.." }

// settingsFile ist die Datei, an der seit MR-047 BEIDE Schichten hängen: die
// Hook-Verdrahtung und die Permission-Sperrliste. Ein Syntaxfehler darin nimmt
// beide zugleich aus dem Spiel, und zwar still.
func settingsFile() string { return repoRoot() + "/.claude/settings.json" }

func guardFile() string { return repoRoot() + "/.claude/hooks/pretooluse-command-guard.sh" }

// settings bildet nur ab, was geprüft wird — nicht das ganze Schema. Ein
// unbekanntes Feld ist kein Fehler: die Datei gehört dem Werkzeug, nicht diesem
// Test, und sie darf wachsen.
type settings struct {
	Permissions struct {
		Deny []string `json:"deny"`
		Ask  []string `json:"ask"`
	} `json:"permissions"`
	Hooks map[string][]struct {
		Matcher string `json:"matcher"`
		Hooks   []struct {
			Type    string `json:"type"`
			Command string `json:"command"`
		} `json:"hooks"`
	} `json:"hooks"`
}

func loadSettings(path string) (settings, error) {
	var s settings
	raw, err := os.ReadFile(path)
	if err != nil {
		return s, fmt.Errorf("nicht lesbar (fail-closed): %w", err)
	}
	if err := json.Unmarshal(raw, &s); err != nil {
		return s, fmt.Errorf("kein gültiges JSON (fail-closed): %w", err)
	}
	return s, nil
}

// blockedFromGuard liest die BLOCKED-Zuweisung des Wächters. Bewusst eng: die
// Zuweisung ist eine Zeichenkette über eine oder mehrere Zeilen, und mehr Form
// als das braucht die Prüfung nicht. GRENZE: der Versions-Suffix-Fall
// (`python`, `python3.12`) steht im Wächter als Muster daneben und ist hier
// nicht erfasst — er wird über die Sperrliste separat abgedeckt.
func blockedFromGuard(path string) ([]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Wächter nicht lesbar (fail-closed): %w", err)
	}
	const key = "BLOCKED=\""
	i := strings.Index(string(raw), key)
	if i < 0 {
		return nil, fmt.Errorf("keine BLOCKED-Zuweisung im Wächter gefunden (fail-closed)")
	}
	rest := string(raw)[i+len(key):]
	j := strings.Index(rest, "\"")
	if j < 0 {
		return nil, fmt.Errorf("BLOCKED-Zuweisung nicht terminiert (fail-closed)")
	}
	return strings.Fields(rest[:j]), nil
}

// hookPaths zieht die auf `.claude/hooks/` zeigenden Pfade aus den
// Hook-Kommandos. GRENZE: das ist eine Teilketten-Suche, keine Shell-Analyse —
// ein Hook, der seinen Pfad zusammensetzt oder aus einer Variablen holt, fällt
// still heraus.
func hookPaths(s settings) []string {
	var out []string
	const marker = "/.claude/hooks/"
	for _, entries := range s.Hooks {
		for _, e := range entries {
			for _, h := range e.Hooks {
				k := strings.Index(h.Command, marker)
				if k < 0 {
					continue
				}
				p := h.Command[k+1:]
				p = strings.TrimRight(strings.Fields(p)[0], "\"'")
				out = append(out, p)
			}
		}
	}
	return out
}

// assertDenyDecktGuard hält die zwei hand-gepflegten Listen aneinander: jeder
// Name, den der Wächter in Befehlsposition sperrt, hat eine Permission-Regel.
// Nicht geprüft wird die Gegenrichtung — die Sperrliste führt zusätzlich `git`
// und `docker`, die der Wächter gar nicht kennt (MR-047).
//
// Das ist KEINE Bewertung der Namen (die bleibt Urteil und lebt im
// Konventionsspeicher), sondern eine Kopplung zweier Artefakte, die beide
// dieselbe Regel umsetzen — dieselbe Bauform wie gate-consistency.
func assertDenyDecktGuard(blocked, deny []string) error {
	have := map[string]bool{}
	for _, d := range deny {
		name := strings.TrimPrefix(d, "Bash(")
		if k := strings.IndexAny(name, " )"); k >= 0 {
			name = name[:k]
		}
		have[name] = true
	}
	for _, b := range blocked {
		if !have[b] {
			return fmt.Errorf("der Wächter sperrt %q, die Permission-Sperrliste nicht — "+
				"die zweite Schicht deckt die Klasse aus AGENTS.md §3.1 dann nur teilweise", b)
		}
	}
	return nil
}

// TestDurchsetzung_SettingsIstGueltigesJSON ist die billigste tragende Aussage:
// eine zerschossene Datei nimmt beide Schichten zugleich aus dem Spiel, und
// ohne diesen Test sähe das niemand.
func TestDurchsetzung_SettingsIstGueltigesJSON(t *testing.T) {
	if _, err := loadSettings(settingsFile()); err != nil {
		t.Fatalf(".claude/settings.json: %v", err)
	}
}

// TestDurchsetzung_HookPfadeExistieren fängt den Fall, den ein Syntax-Test
// nicht sieht: gültiges JSON, das auf eine Datei zeigt, die es nicht gibt.
func TestDurchsetzung_HookPfadeExistieren(t *testing.T) {
	s, err := loadSettings(settingsFile())
	if err != nil {
		t.Fatalf(".claude/settings.json: %v", err)
	}
	paths := hookPaths(s)
	if len(paths) == 0 {
		t.Fatal("kein Hook mit einem .claude/hooks/-Pfad verdrahtet (fail-closed) — " +
			"entweder ist die Verdrahtung weg oder die Teilketten-Suche trägt nicht mehr")
	}
	for _, p := range paths {
		if _, err := os.Stat(repoRoot() + "/" + p); err != nil {
			t.Errorf("verdrahteter Hook %q existiert nicht: %v", p, err)
		}
	}
}

// TestDurchsetzung_SperrlisteDecktWaechter koppelt die zwei hand-gepflegten
// Listen. Ohne diese Kopplung driften sie: bei ihrer Einführung fehlten elf
// Namen in der Sperrliste, ohne dass etwas rot wurde.
func TestDurchsetzung_SperrlisteDecktWaechter(t *testing.T) {
	s, err := loadSettings(settingsFile())
	if err != nil {
		t.Fatalf(".claude/settings.json: %v", err)
	}
	blocked, err := blockedFromGuard(guardFile())
	if err != nil {
		t.Fatalf("Wächter: %v", err)
	}
	if len(blocked) == 0 {
		t.Fatal("BLOCKED-Liste des Wächters ist leer (fail-closed)")
	}
	if err := assertDenyDecktGuard(blocked, s.Permissions.Deny); err != nil {
		t.Fatalf("Durchsetzungs-Schichten driften: %v", err)
	}
}

// TestDurchsetzung_GuardsMeldenDieMutation ist die Gegenprobe zu den drei
// Live-Tests: dieselben Guards, gegen konstruierte Eingaben. Ohne sie wäre
// unbekannt, ob ein grüner Lauf oben etwas prüft oder nur nichts findet
// (BEO-017).
func TestDurchsetzung_GuardsMeldenDieMutation(t *testing.T) {
	if err := assertDenyDecktGuard([]string{"pip"}, []string{"Bash(npm *)"}); err == nil {
		t.Error("fehlender Name in der Sperrliste wurde nicht gemeldet")
	}
	if err := assertDenyDecktGuard([]string{"pip"}, []string{"Bash(pip *)"}); err != nil {
		t.Errorf("gedeckter Name wurde gemeldet: %v", err)
	}
	if _, err := loadSettings(repoRoot() + "/.claude/gibt-es-nicht.json"); err == nil {
		t.Error("fehlende Datei wurde nicht gemeldet")
	}
	if _, err := blockedFromGuard(repoRoot() + "/AGENTS.md"); err == nil {
		t.Error("Datei ohne BLOCKED-Zuweisung wurde nicht gemeldet")
	}
	if len(hookPaths(settings{})) != 0 {
		t.Error("leere Verdrahtung ergab Pfade")
	}
}
