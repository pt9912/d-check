package configyaml_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

// Invarianten der committeten Durchsetzungs-Artefakte unter `.claude/`
// (MR-048).
//
// ZUSAGE: sind die Dateien da, sind sie wohlgeformt und zeigen aufeinander.
//
// ABGRENZUNG, und sie ist der wichtigere Satz: dass die Durchsetzung LÄUFT,
// sagt hier nichts. Ob das Werkzeug den Hook ruft und seine Antwort befolgt,
// steht in keiner Datei; AGENTS.md §3.1 erklärt einen Lauf ohne dieses Werkzeug
// ausdrücklich für ungebunden. Grün schließt einen Ausfallgrund aus, nicht den
// Ausfall.
//
// ABGRENZUNG zur Anwesenheit: fehlt `.claude/`, überspringen diese Tests. Eine
// Werkzeug-Einstellung darf fehlen; kaputt sein darf sie nicht.

func repoRoot() string { return "../../../.." }

func settingsFile() string { return repoRoot() + "/.claude/settings.json" }

func guardFile() string { return repoRoot() + "/.claude/hooks/pretooluse-command-guard.sh" }

// errAbwesend trennt „gibt es nicht" von „ist kaputt". Nur das Zweite ist ein
// Befund.
var errAbwesend = errors.New("artefakt nicht vorhanden")

// settings bildet ab, was geprüft wird — nicht das ganze Schema. Unbekannte
// Felder sind kein Fehler: die Datei gehört dem Werkzeug und darf wachsen.
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
	if os.IsNotExist(err) {
		return s, errAbwesend
	}
	if err != nil {
		return s, fmt.Errorf("nicht lesbar: %w", err)
	}
	if err := json.Unmarshal(raw, &s); err != nil {
		return s, fmt.Errorf("kein gültiges JSON: %w", err)
	}
	return s, nil
}

// blockedFromGuard liest die BLOCKED-Zuweisung des Wächters.
// GRENZE: erfasst ist die eine Zuweisung in doppelten Anführungszeichen. Eine
// zweite Zuweisung oder eine Fortschreibung (`BLOCKED="$BLOCKED …"`) ist
// mehrdeutig und endet deshalb als Befund statt als stille Teilmenge.
func blockedFromGuard(path string) ([]string, error) {
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, errAbwesend
	}
	if err != nil {
		return nil, fmt.Errorf("Wächter nicht lesbar: %w", err)
	}
	const key = "BLOCKED=\""
	text := string(raw)
	if n := strings.Count(text, key); n != 1 {
		return nil, fmt.Errorf("BLOCKED-Zuweisung %d-mal gefunden, genau eine erwartet", n)
	}
	rest := text[strings.Index(text, key)+len(key):]
	j := strings.Index(rest, "\"")
	if j < 0 {
		return nil, errors.New("BLOCKED-Zuweisung nicht terminiert")
	}
	return strings.Fields(rest[:j]), nil
}

// hookRef ist ein Hook-Eintrag, reduziert auf das, was prüfbar ist.
type hookRef struct {
	Event   string
	Matcher string
	Command string
	Paths   []string
}

// hookRefs sammelt die Hook-Einträge und die darin genannten `.claude/hooks/`-
// Pfade.
//
// Der Marker trägt KEINEN führenden Schrägstrich: sonst fällt die literale
// Relativform (`bash .claude/hooks/x.sh`) heraus, und ein Hook, der ins Leere
// zeigt, bliebe still — die Schreibweise entschiede über die Prüfung, nicht der
// Zustand. Gesammelt werden ALLE Vorkommen je Kommando, nicht nur das erste.
//
// GRENZE: das ist eine Teilketten-Suche, keine Shell-Analyse. Ein Pfad, der
// erst zur Laufzeit entsteht, fällt heraus; und ein Pfad ausserhalb des Repos
// (`~/.claude/hooks/…`) wird gegen die Repo-Wurzel aufgelöst und damit falsch
// zugeordnet.
func hookRefs(s settings) []hookRef {
	const marker = ".claude/hooks/"
	var out []hookRef
	for event, entries := range s.Hooks {
		for _, e := range entries {
			for _, h := range e.Hooks {
				ref := hookRef{Event: event, Matcher: e.Matcher, Command: h.Command}
				for i := 0; ; {
					k := strings.Index(h.Command[i:], marker)
					if k < 0 {
						break
					}
					start := i + k
					p := h.Command[start:]
					if f := strings.Fields(p); len(f) > 0 {
						p = strings.TrimRight(f[0], "\"';")
					}
					ref.Paths = append(ref.Paths, p)
					i = start + len(marker)
				}
				out = append(out, ref)
			}
		}
	}
	return out
}

// denyDeckt sammelt die Namen, für die eine Sperrregel die GANZE Befehlsklasse
// abdeckt — also genau die Form `Bash(<name> *)`. Eine verengte Regel wie
// `Bash(pip install *)` sperrt einen Teilbaum und zählt deshalb nicht: sonst
// meldete der Test Deckung, wo `pip download …` durchliefe.
func denyDeckt(deny []string) map[string]bool {
	have := map[string]bool{}
	for _, d := range deny {
		inner, ok := strings.CutPrefix(d, "Bash(")
		if !ok {
			continue
		}
		inner, ok = strings.CutSuffix(inner, ")")
		if !ok {
			continue
		}
		name, rest, ok := strings.Cut(inner, " ")
		if !ok || rest != "*" {
			continue
		}
		have[name] = true
	}
	return have
}

// assertDenyDecktGuard koppelt die zwei hand-gepflegten Sperrlisten: jeder
// Name, den der Wächter in Befehlsposition sperrt, hat eine Permission-Regel
// über seine ganze Befehlsklasse.
//
// ABGRENZUNG: geprüft wird Gleichheit zweier Artefakte, nicht die Richtigkeit
// der Namen — die bleibt Urteil und lebt im Konventionsspeicher. Die
// Gegenrichtung ist bewusst offen: die Sperrliste führt zusätzlich `git` und
// `docker`, die der Wächter gar nicht kennt (MR-047).
func assertDenyDecktGuard(blocked, deny []string) error {
	have := denyDeckt(deny)
	for _, b := range blocked {
		if !have[b] {
			return fmt.Errorf("der Wächter sperrt %q, die Permission-Sperrliste deckt es nicht "+
				"als ganze Befehlsklasse (`Bash(%s *)`) — die zweite Schicht trägt die "+
				"Klasse aus AGENTS.md §3.1 dann nur teilweise", b, b)
		}
	}
	return nil
}

// ladeOderUeberspringe hält die Anwesenheits-Abgrenzung an einer Stelle.
func ladeOderUeberspringe(t *testing.T) settings {
	t.Helper()
	s, err := loadSettings(settingsFile())
	if errors.Is(err, errAbwesend) {
		t.Skip(".claude/settings.json nicht vorhanden — Werkzeug-Einstellung, kein Repo-Vertrag")
	}
	if err != nil {
		t.Fatalf(".claude/settings.json: %v", err)
	}
	return s
}

// TestDurchsetzung_SettingsIstGueltigesJSON: eine zerschossene Datei nimmt
// beide Schichten zugleich aus dem Spiel, und ohne diese Prüfung sähe das
// niemand.
func TestDurchsetzung_SettingsIstGueltigesJSON(t *testing.T) {
	ladeOderUeberspringe(t)
}

// TestDurchsetzung_HookVerdrahtungTraegt fängt, was ein Syntax-Test nicht
// sieht: gültiges JSON, das ins Leere zeigt.
func TestDurchsetzung_HookVerdrahtungTraegt(t *testing.T) {
	refs := hookRefs(ladeOderUeberspringe(t))
	if len(refs) == 0 {
		t.Fatal("kein Hook verdrahtet — die Durchsetzungs-Schicht hätte dann keinen Einstieg")
	}
	for _, r := range refs {
		if len(r.Paths) == 0 {
			t.Errorf("Hook %q (matcher %q) nennt kein Skript unter .claude/hooks/: %q",
				r.Event, r.Matcher, r.Command)
			continue
		}
		for _, p := range r.Paths {
			if _, err := os.Stat(repoRoot() + "/" + p); err != nil {
				t.Errorf("verdrahteter Hook %q existiert nicht: %v", p, err)
			}
		}
	}
}

// TestDurchsetzung_WaechterHaengtAmBashWerkzeug: ein anderer `matcher` nimmt
// den Wächter aus jedem Bash-Aufruf heraus, ohne dass eine Datei fehlt oder
// kaputt ist — derselbe stille Ausfall wie ein Pfad ins Leere.
func TestDurchsetzung_WaechterHaengtAmBashWerkzeug(t *testing.T) {
	refs := hookRefs(ladeOderUeberspringe(t))
	for _, r := range refs {
		if r.Event != "PreToolUse" {
			continue
		}
		for _, p := range r.Paths {
			if !strings.Contains(p, "command-guard") {
				continue
			}
			if r.Matcher != "Bash" {
				t.Errorf("der Befehls-Wächter hängt am matcher %q statt an \"Bash\" — "+
					"er feuert dann auf keinen Bash-Aufruf mehr", r.Matcher)
			}
			return
		}
	}
	t.Error("kein PreToolUse-Hook auf den Befehls-Wächter verdrahtet")
}

// TestDurchsetzung_SperrlisteDecktWaechter koppelt die zwei Sperrlisten. Ohne
// die Kopplung driften sie, und nichts wird rot.
func TestDurchsetzung_SperrlisteDecktWaechter(t *testing.T) {
	s := ladeOderUeberspringe(t)
	blocked, err := blockedFromGuard(guardFile())
	if errors.Is(err, errAbwesend) {
		t.Skip("Wächter nicht vorhanden — Werkzeug-Einstellung, kein Repo-Vertrag")
	}
	if err != nil {
		t.Fatalf("Wächter: %v", err)
	}
	if len(blocked) == 0 {
		t.Fatal("BLOCKED-Liste des Wächters ist leer")
	}
	if err := assertDenyDecktGuard(blocked, s.Permissions.Deny); err != nil {
		t.Fatalf("Durchsetzungs-Schichten driften: %v", err)
	}
}

// TestDurchsetzung_GuardsMeldenDieMutation ist die Gegenprobe zu den
// Live-Tests: dieselben Guards, gegen konstruierte Eingaben. Ohne sie bliebe
// offen, ob ein grüner Lauf oben etwas prüft oder nur nichts findet.
func TestDurchsetzung_GuardsMeldenDieMutation(t *testing.T) {
	if err := assertDenyDecktGuard([]string{"pip"}, []string{"Bash(npm *)"}); err == nil {
		t.Error("fehlender Name in der Sperrliste wurde nicht gemeldet")
	}
	if err := assertDenyDecktGuard([]string{"pip"}, []string{"Bash(pip install *)"}); err == nil {
		t.Error("verengte Regel zählte als Deckung der ganzen Befehlsklasse")
	}
	if err := assertDenyDecktGuard([]string{"pip"}, []string{"pip"}); err == nil {
		t.Error("Eintrag ohne Bash-Hülle zählte als Deckung")
	}
	if err := assertDenyDecktGuard([]string{"pip"}, []string{"Bash(pip *)"}); err != nil {
		t.Errorf("gedeckter Name wurde gemeldet: %v", err)
	}
	if _, err := loadSettings(repoRoot() + "/.claude/gibt-es-nicht.json"); !errors.Is(err, errAbwesend) {
		t.Errorf("fehlende Datei wurde nicht als abwesend gemeldet: %v", err)
	}
	if _, err := blockedFromGuard(repoRoot() + "/AGENTS.md"); err == nil {
		t.Error("Datei ohne BLOCKED-Zuweisung wurde nicht gemeldet")
	}
	if len(hookRefs(settings{})) != 0 {
		t.Error("leere Verdrahtung ergab Einträge")
	}
	mitPfaden := mustHookRefs(t, `{"hooks":{"PreToolUse":[{"matcher":"Bash","hooks":[`+
		`{"type":"command","command":"bash .claude/hooks/a.sh && bash x/.claude/hooks/b.sh"}]}]}}`)
	if len(mitPfaden) != 1 || len(mitPfaden[0].Paths) != 2 {
		t.Errorf("beide Schreibweisen müssen erfasst sein, gefunden: %+v", mitPfaden)
	}
}

func mustHookRefs(t *testing.T, raw string) []hookRef {
	t.Helper()
	var s settings
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatalf("Testeingabe dekodiert nicht: %v", err)
	}
	return hookRefs(s)
}
