package rules

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

const (
	psDockerPat = `^ARG GO_VERSION=v?(\S+)`
	psMakePat   = `^GO_VERSION\s*:?=\s*v?(\S+)`
)

// psSrc und psPin bauen die Konfiguration; psPins fasst mehrere Pins zusammen.
func psSrc(file, pattern string) model.PinsyncSource {
	return model.PinsyncSource{File: file, Pattern: pattern}
}

func psPin(name string, srcs ...model.PinsyncSource) model.PinsyncPin {
	return model.PinsyncPin{Name: name, Sources: srcs}
}

func psCfg(pins ...model.PinsyncPin) model.PinsyncConfig {
	return model.PinsyncConfig{Pins: pins}
}

// psGo ist der Regelfall: dieselbe Version in Dockerfile und Makefile.
func psGo() model.PinsyncPin {
	return psPin("go", psSrc("Dockerfile", psDockerPat), psSrc("Makefile", psMakePat))
}

func psRun(t *testing.T, files map[string]string, cfg model.PinsyncConfig) (PinsyncResult, error) {
	t.Helper()
	return CheckPinsync(coretest.NewMemFS(files), cfg)
}

func psMust(t *testing.T, files map[string]string, cfg model.PinsyncConfig) PinsyncResult {
	t.Helper()
	r, err := psRun(t, files, cfg)
	if err != nil {
		t.Fatalf("unerwarteter Eingabefehler: %v", err)
	}
	return r
}

func psReasons(r PinsyncResult) []string {
	var out []string
	for _, f := range r.Findings {
		out = append(out, fmt.Sprintf("%s@%s:%d", f.Reason, f.File, f.Line))
	}
	return out
}

func psWant(t *testing.T, r PinsyncResult, want ...string) {
	t.Helper()
	got := psReasons(r)
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("Befunde: got %v, want %v", got, want)
	}
}

func psGoFiles(docker, mk string) map[string]string {
	return map[string]string{"Dockerfile": docker, "Makefile": mk}
}

// Happy Path (DC-FA-PSYNC-001): gleicher Wert, kein Befund; die Bezugsmenge
// nennt Pins und gelesene Quellen (DC-FA-PSYNC-001).
func TestPinsyncGleicherWertKeinBefund(t *testing.T) {
	r := psMust(t, psGoFiles("FROM x\nARG GO_VERSION=1.27.3\n", "GO_VERSION := 1.27.3\n"), psCfg(psGo()))
	psWant(t, r)
	if r.Note() != "pinsync: 1 Pin(s) geprüft, 2 Quelle(n) gelesen" {
		t.Fatalf("Note: %q", r.Note())
	}
}

// Bezugsmenge: zwei Pins, vier Quellen (DC-FA-PSYNC-001).
func TestPinsyncNoteNenntDieBezugsmenge(t *testing.T) {
	files := map[string]string{"a": "x=1\n", "b": "x=1\n", "c": "y=2\n", "d": "y=2\n"}
	pa := psPin("p1", psSrc("a", `x=(\d)`), psSrc("b", `x=(\d)`))
	pb := psPin("p2", psSrc("c", `y=(\d)`), psSrc("d", `y=(\d)`))
	r := psMust(t, files, psCfg(pa, pb))
	if r.Note() != "pinsync: 2 Pin(s) geprüft, 4 Quelle(n) gelesen" {
		t.Fatalf("Note: %q", r.Note())
	}
}

// Ungleiche Werte: Befund an der abweichenden Quelle, die Meldung nennt beide
// Werte und beide Dateien (DC-FA-PSYNC-001).
func TestPinsyncMismatchNenntBeideWerteUndDateien(t *testing.T) {
	r := psMust(t, psGoFiles("ARG GO_VERSION=1.27.3\n", "\nGO_VERSION := 1.27.4\n"), psCfg(psGo()))
	psWant(t, r, "pinsync-mismatch@Makefile:2")
	f := r.Findings[0]
	if f.Rule != "pinsync" || f.Target != "go" {
		t.Fatalf("Regel/Ziel: %+v", f)
	}
	for _, s := range []string{"Dockerfile=1.27.3", "Makefile=1.27.4", `"1.27.3"`, `"1.27.4"`} {
		if !strings.Contains(f.Message, s) {
			t.Fatalf("Meldung ohne %q: %s", s, f.Message)
		}
	}
}

// Drei Quellen mit mehreren Abweichungen: jede weicht vom Bezug (erste Quelle)
// ab und bekommt einen Befund; die dritte stimmt überein und bleibt still (DC-FA-PSYNC-001).
func TestPinsyncDreiQuellenMehrereAbweichungen(t *testing.T) {
	files := map[string]string{"a": "v=1\n", "b": "v=2\n", "c": "v=3\n", "d": "v=1\n"}
	pin := psPin("p", psSrc("a", `v=(\d)`), psSrc("b", `v=(\d)`), psSrc("c", `v=(\d)`), psSrc("d", `v=(\d)`))
	r := psMust(t, files, psCfg(pin))
	psWant(t, r, "pinsync-mismatch@b:1", "pinsync-mismatch@c:1")
	if !strings.Contains(r.Findings[1].Message, "a=1, b=2, c=3, d=1") {
		t.Fatalf("Meldung nennt nicht alle Werte: %s", r.Findings[1].Message)
	}
}

// Der Bezug ist die erste Quelle mit genau einem Treffer: eine Quelle davor
// ohne Treffer trägt nichts bei, ihr Befund steht, ihr Wert ist `?` (DC-FA-PSYNC-001).
func TestPinsyncBezugIstErsteQuelleMitGenauEinemTreffer(t *testing.T) {
	files := map[string]string{"a": "nichts\n", "b": "v=1\n", "c": "v=2\n"}
	pin := psPin("p", psSrc("a", `v=(\d)`), psSrc("b", `v=(\d)`), psSrc("c", `v=(\d)`))
	r := psMust(t, files, psCfg(pin))
	psWant(t, r, "pinsync-no-match@a:1", "pinsync-mismatch@c:1")
	if !strings.Contains(r.Findings[1].Message, "a=?") || !strings.Contains(r.Findings[1].Message, "b=1") {
		t.Fatalf("Meldung: %s", r.Findings[1].Message)
	}
}

// Weniger als zwei Quellen mit genau einem Treffer: der Vergleich entfällt, der
// Befund der anderen Quelle steht (ohne und mit mehreren Treffern) (DC-FA-PSYNC-001).
func TestPinsyncWenigerAlsZweiGenauEinTreffer(t *testing.T) {
	pin := psPin("p", psSrc("a", `v=(\d)`), psSrc("b", `v=(\d)`))
	r := psMust(t, map[string]string{"a": "v=1\n", "b": "x\n"}, psCfg(pin))
	psWant(t, r, "pinsync-no-match@b:1")
	r = psMust(t, map[string]string{"a": "v=1\n", "b": "v=2\nv=3\n"}, psCfg(pin))
	psWant(t, r, "pinsync-ambiguous@b:1")
	r = psMust(t, map[string]string{"a": "x\n", "b": "v=2\nv=3\n"}, psCfg(pin))
	psWant(t, r, "pinsync-no-match@a:1", "pinsync-ambiguous@b:1")
}

// Kein Treffer: Zeile 1, auch bei einer Gruppe, die nur leer fängt, bei leerer
// Datei und bei einer nicht beteiligten Gruppe (DC-FA-PSYNC-001).
func TestPinsyncKeinTrefferLeereGruppeLeereDatei(t *testing.T) {
	ok := "v=1\n"
	cases := map[string]struct{ content, pat string }{
		"kein Vorkommen":        {"x\ny\n", `v=(\d)`},
		"leer fangende Gruppe":  {"v=\n", `v=(\d*)`},
		"leere Datei":           {"", `v=(\d)`},
		"nicht beteiligt":       {"v=\n", `v=(\d)?`},
		"erste Passung ist leer": {"abc 12\n", `(\d*)`},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			pin := psPin("p", psSrc("a", `v=(\d)`), psSrc("b", c.pat))
			r := psMust(t, map[string]string{"a": ok, "b": c.content}, psCfg(pin))
			psWant(t, r, "pinsync-no-match@b:1")
		})
	}
}

// Mehrdeutig: zwei Zeilen treffen; der Befund steht an der ersten, die Meldung
// nennt alle Zeilen. Mehrere Vorkommen in einer Zeile zählen als eine Zeile (DC-FA-PSYNC-001).
func TestPinsyncMehrdeutigUndMehrfachVorkommenInEinerZeile(t *testing.T) {
	pin := psPin("p", psSrc("a", `v=(\d)`), psSrc("b", `v=(\d)`))
	r := psMust(t, map[string]string{"a": "v=1\n", "b": "x\nv=1\ny\nv=1\n"}, psCfg(pin))
	psWant(t, r, "pinsync-ambiguous@b:2")
	if !strings.Contains(r.Findings[0].Message, "2, 4") {
		t.Fatalf("Meldung nennt nicht alle Zeilen: %s", r.Findings[0].Message)
	}
	r = psMust(t, map[string]string{"a": "v=1\n", "b": "v=1 v=1 v=1\n"}, psCfg(pin))
	psWant(t, r)
	// Die erste Passung der Zeile zählt, nicht die letzte.
	r = psMust(t, map[string]string{"a": "v=1\n", "b": "v=2 v=1\n"}, psCfg(pin))
	psWant(t, r, "pinsync-mismatch@b:1")
}

// Randformen des Lesens: CRLF, letzte Zeile ohne Zeilenende, leere Datei (DC-FA-PSYNC-001).
func TestPinsyncZeilenRandformen(t *testing.T) {
	pin := psPin("p", psSrc("a", `v=(.+)`), psSrc("b", `v=(.+)`))
	cases := map[string]string{
		"CRLF":                  "x\r\nv=1\r\ny\r\n",
		"ohne Schlusszeilenende": "x\nv=1",
		"nur eine Zeile":        "v=1\n",
	}
	for name, content := range cases {
		t.Run(name, func(t *testing.T) {
			r := psMust(t, map[string]string{"a": "v=1\n", "b": content}, psCfg(pin))
			psWant(t, r) // ein \r gehört nicht zum Wert; `.+` fängt es ohne das Abschneiden mit
		})
	}
	// Mit `$` verankert: das \r am Zeilenende stört nicht.
	pin = psPin("p", psSrc("a", `^v=(\d)$`), psSrc("b", `^v=(\d)$`))
	psWant(t, psMust(t, map[string]string{"a": "v=1\r\n", "b": "v=1\n"}, psCfg(pin)))
	// Nur ein einzelnes \r wird abgeschnitten.
	pin = psPin("p", psSrc("a", `v=(.+)`), psSrc("b", `v=(.+)`))
	r := psMust(t, map[string]string{"a": "v=1\r\r\n", "b": "v=1\n"}, psCfg(pin))
	psWant(t, r, "pinsync-mismatch@b:1")
}

// v-Präfix: außerhalb der Gruppe (v?) ist der Wert gleich, innerhalb der
// Gruppe nicht -- das Modul normalisiert nichts (DC-FA-PSYNC-001).
func TestPinsyncVPraefixAusserhalbUndInnerhalbDerGruppe(t *testing.T) {
	files := map[string]string{"a": "ARG V=v1.27\n", "b": "V=1.27\n"}
	aussen := psPin("p", psSrc("a", `V=v?(\S+)`), psSrc("b", `V=v?(\S+)`))
	psWant(t, psMust(t, files, psCfg(aussen)))
	innen := psPin("p", psSrc("a", `V=(v?\S+)`), psSrc("b", `V=(v?\S+)`))
	psWant(t, psMust(t, files, psCfg(innen)), "pinsync-mismatch@b:1")
}

// Kommentarzeile desselben Wortlauts und FROM-Zeile ähnlichen Werts: das Modul
// liest Text und kennt kein Format -- ein Treffer in einem Kommentar ist ein
// Treffer, und ein zweiter Treffer macht den Ausdruck mehrdeutig (DC-FA-PSYNC-001).
func TestPinsyncKommentarUndFromZeileSindTreffer(t *testing.T) {
	docker := "# ARG GO_VERSION=1.27.3\nARG GO_VERSION=1.27.3\nFROM golang:1.27.3\n"
	unverankert := psPin("go", psSrc("Dockerfile", `GO_VERSION=v?(\S+)`), psSrc("Makefile", psMakePat))
	r := psMust(t, psGoFiles(docker, "GO_VERSION := 1.27.3\n"), psCfg(unverankert))
	psWant(t, r, "pinsync-ambiguous@Dockerfile:1")
	// Verankert trifft die Kommentarzeile nicht: nur die ARG-Zeile gilt.
	r = psMust(t, psGoFiles(docker, "GO_VERSION := 1.27.3\n"), psCfg(psGo()))
	psWant(t, r)
	// Nur der Kommentar trifft: er gilt als Wert der Quelle.
	pin := psPin("p", psSrc("a", `GO_VERSION=(\S+)`), psSrc("b", `GO_VERSION=(\S+)`))
	r = psMust(t, map[string]string{"a": "# GO_VERSION=1.27.9\n", "b": "GO_VERSION=1.27.3\n"}, psCfg(pin))
	psWant(t, r, "pinsync-mismatch@b:1")
	// FROM-Zeile mit ähnlichem Wert, Ausdruck auf `golang:` angelegt.
	pin = psPin("p", psSrc("a", `^FROM golang:(\S+)`), psSrc("b", `^FROM golang:(\S+)`))
	r = psMust(t, map[string]string{"a": "FROM golang:1.27.3 AS b\n", "b": "FROM golang:1.27.4-alpine\n"}, psCfg(pin))
	psWant(t, r, "pinsync-mismatch@b:1")
}

// Zwei verschiedene Ausdrücke derselben Datei, die dieselbe Zeile treffen, sind
// erlaubt (verschiedenes pattern) und liefern je einen eigenen Wert (DC-FA-PSYNC-001).
func TestPinsyncZweiAusdruckeDerselbenDateiAufDerselbenZeile(t *testing.T) {
	files := map[string]string{"a": "X=1 Y=2\n", "b": "Z=1\n"}
	pin := psPin("p", psSrc("a", `X=(\d)`), psSrc("a", `Y=(\d)`), psSrc("b", `Z=(\d)`))
	r := psMust(t, files, psCfg(pin))
	psWant(t, r, "pinsync-mismatch@a:1")
	if !strings.Contains(r.Findings[0].Message, "a=1, a=2, b=1") {
		t.Fatalf("Meldung: %s", r.Findings[0].Message)
	}
	if r.Sources != 3 {
		t.Fatalf("Quellen gelesen: %d", r.Sources)
	}
}

// Reihenfolge: Pins nach Namen, Quellen in Konfigurations-Reihenfolge,
// unabhängig von der Reihenfolge in der Konfiguration (DC-QA-02).
func TestPinsyncReihenfolgeDeterminismus(t *testing.T) {
	files := map[string]string{"a": "v=1\n", "b": "v=2\n", "c": "v=3\n", "d": "v=4\n"}
	pz := psPin("z", psSrc("a", `v=(\d)`), psSrc("b", `v=(\d)`))
	pa := psPin("a", psSrc("c", `v=(\d)`), psSrc("d", `v=(\d)`))
	r1 := psMust(t, files, psCfg(pz, pa))
	r2 := psMust(t, files, psCfg(pa, pz))
	psWant(t, r1, "pinsync-mismatch@d:1", "pinsync-mismatch@b:1")
	if fmt.Sprint(r1.Findings) != fmt.Sprint(r2.Findings) {
		t.Fatalf("nicht deterministisch:\n%v\n%v", r1.Findings, r2.Findings)
	}
}

// Eingabefehler der Konfiguration (Exit 2, kein Befund). Jede Zeile ist eine
// Randbedingung der Anforderung (DC-FA-PSYNC-001).
func TestPinsyncEingabefehlerDerKonfiguration(t *testing.T) {
	ok := psSrc("a", `v=(\d)`)
	ok2 := psSrc("b", `v=(\d)`)
	cases := map[string]struct {
		cfg  model.PinsyncConfig
		want string
	}{
		"kein Pin":                 {psCfg(), "pinsync.pins ist leer"},
		"Pin ohne Namen":           {psCfg(psPin("", ok, ok2)), "keinen Namen"},
		"doppelter Name":           {psCfg(psPin("p", ok, ok2), psPin("p", ok, ok2)), "mehrfach"},
		"keine Quelle":             {psCfg(psPin("p")), "0 Quelle(n)"},
		"eine Quelle":              {psCfg(psPin("p", ok)), "1 Quelle(n)"},
		"file fehlt":               {psCfg(psPin("p", psSrc("", `v=(\d)`), ok2)), "Pflicht"},
		"pattern fehlt":            {psCfg(psPin("p", psSrc("a", ""), ok2)), "Pflicht"},
		"absoluter Pfad":           {psCfg(psPin("p", psSrc("/a", `v=(\d)`), ok2)), "beginnt mit '/'"},
		"Pfad mit ..":              {psCfg(psPin("p", psSrc("../a", `v=(\d)`), ok2)), "segment"},
		"Pfad mit .":               {psCfg(psPin("p", psSrc("./a", `v=(\d)`), ok2)), "segment"},
		"leeres Segment":           {psCfg(psPin("p", psSrc("x//a", `v=(\d)`), ok2)), "segment"},
		"Glob-Zeichen":             {psCfg(psPin("p", psSrc("a*", `v=(\d)`), ok2)), "glob-zeichen"},
		"Steuerzeichen":            {psCfg(psPin("p", psSrc("a\tb", `v=(\d)`), ok2)), "steuer"},
		"abschließender Schrägstrich": {psCfg(psPin("p", psSrc("Dockerfile/", `v=(\d)`), ok2)), "endet auf '/'"},
		"Ausdruck nicht übersetzbar": {psCfg(psPin("p", psSrc("a", `v=(\d`), ok2)), "kein gültiges RE2"},
		"Rückwärts-Verweis (RE2)":  {psCfg(psPin("p", psSrc("a", `(a)\1`), ok2)), "kein gültiges RE2"},
		"Lookahead (RE2)":          {psCfg(psPin("p", psSrc("a", `v(?=1)(\d)`), ok2)), "kein gültiges RE2"},
		"keine Fanggruppe":         {psCfg(psPin("p", psSrc("a", `v=\d`), ok2)), "0 Fanggruppen"},
		"nur nicht fangende Gruppe": {psCfg(psPin("p", psSrc("a", `v=(?:\d)`), ok2)), "0 Fanggruppen"},
		"zwei Fanggruppen":         {psCfg(psPin("p", psSrc("a", `(v)=(\d)`), ok2)), "2 Fanggruppen"},
		"benannte Gruppen zählen":  {psCfg(psPin("p", psSrc("a", `(?P<x>v)=(\d)`), ok2)), "2 Fanggruppen"},
		"dieselbe Stelle zweimal":  {psCfg(psPin("p", ok, psSrc("b", `w=(\d)`), ok)), "gegen sich selbst"},
	}
	files := map[string]string{"a": "v=1\n", "b": "v=1\n"}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			r, err := psRun(t, files, c.cfg)
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("erwartet Fehler mit %q, got %v", c.want, err)
			}
			if len(r.Findings) != 0 {
				t.Fatalf("ein Eingabefehler liefert keinen Befund: %+v", r.Findings)
			}
		})
	}
}

// Eine benannte Gruppe allein und eine nicht fangende neben einer fangenden
// sind genau eine Fanggruppe (DC-FA-PSYNC-001).
func TestPinsyncGenauEineFanggruppeAkzeptiert(t *testing.T) {
	files := map[string]string{"a": "v=1\n", "b": "v=1\n"}
	for _, pat := range []string{`v=(?P<wert>\d)`, `(?:v)=(\d)`} {
		pin := psPin("p", psSrc("a", pat), psSrc("b", pat+` *`))
		psWant(t, psMust(t, files, psCfg(pin)))
	}
}

// Dieselbe Datei mit verschiedenem Ausdruck ist keine Stelle gegen sich selbst (DC-FA-PSYNC-001).
func TestPinsyncGleicheDateiAnderesPatternIstErlaubt(t *testing.T) {
	pin := psPin("p", psSrc("a", `x=(\d)`), psSrc("a", `y=(\d)`))
	psWant(t, psMust(t, map[string]string{"a": "x=1\ny=1\n"}, psCfg(pin)))
}

// Eingabefehler beim Lesen: fehlende Datei, Verzeichnis, Symlink (auch als
// Verzeichnis-Segment), Datei mitten im Pfad, Lesefehler (DC-FA-PSYNC-001).
func TestPinsyncEingabefehlerBeimLesen(t *testing.T) {
	good := psSrc("b", `v=(\d)`)
	newFS := func() *coretest.MemFS {
		m := coretest.NewMemFS(map[string]string{"b": "v=1\n", "d/x": "v=1\n", "file": "v=1\n", "real/y": "v=1\n"})
		m.AddSymlink("link")
		m.AddSymlink("linkdir")
		return m
	}
	cases := map[string]struct{ file, want string }{
		"fehlende Datei":             {"nope", "fehlt"},
		"fehlendes Verzeichnis":      {"nodir/x", "fehlt"},
		"Verzeichnis als Quelle":     {"d", "keine Datei"},
		"Symlink als Quelle":         {"link", "Symlink"},
		"Symlink als Segment":        {"linkdir/y", "Symlink"},
		"Datei als Verzeichnis-Segment": {"file/x", "kein Verzeichnis"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			pin := psPin("p", psSrc(c.file, `v=(\d)`), good)
			r, err := CheckPinsync(newFS(), psCfg(pin))
			if err == nil || !strings.Contains(err.Error(), c.want) || !strings.Contains(err.Error(), `pin "p"`) {
				t.Fatalf("erwartet Fehler mit %q und Pin, got %v", c.want, err)
			}
			if len(r.Findings) != 0 {
				t.Fatalf("Befund trotz Eingabefehler: %+v", r.Findings)
			}
		})
	}
}

// failReadFS meldet Kind=Datei, scheitert aber beim Lesen.
type failReadFS struct{ driven.Filesystem }

func (failReadFS) ReadFile(string) ([]byte, error) { return nil, fmt.Errorf("permission denied") }

// Ein Lesefehler außer „fehlend" ist Exit 2 mit Pin, Quelle und Ursache (DC-FA-PSYNC-001).
func TestPinsyncLesefehlerIstEingabefehler(t *testing.T) {
	fs := failReadFS{coretest.NewMemFS(map[string]string{"a": "v=1\n", "b": "v=1\n"})}
	pin := psPin("p", psSrc("a", `v=(\d)`), psSrc("b", `v=(\d)`))
	_, err := CheckPinsync(fs, psCfg(pin))
	if err == nil || !strings.Contains(err.Error(), "permission denied") || !strings.Contains(err.Error(), "a") {
		t.Fatalf("erwartet Lesefehler mit Ursache, got %v", err)
	}
}

// failKindFS scheitert schon bei Kind (etwa ENOTDIR).
type failKindFS struct{ driven.Filesystem }

func (failKindFS) Kind(string) (driven.EntryKind, error) { return driven.KindMissing, fmt.Errorf("not a directory") }

// Ein Fehler von Kind (etwa ENOTDIR) ist ein Eingabefehler mit Ursache (DC-FA-PSYNC-001).
func TestPinsyncKindFehlerIstEingabefehler(t *testing.T) {
	fs := failKindFS{coretest.NewMemFS(map[string]string{"a": "v=1\n", "b": "v=1\n"})}
	pin := psPin("p", psSrc("a", `v=(\d)`), psSrc("b", `v=(\d)`))
	if _, err := CheckPinsync(fs, psCfg(pin)); err == nil || !strings.Contains(err.Error(), "not a directory") {
		t.Fatalf("erwartet Kind-Fehler, got %v", err)
	}
}

// Reihenfolge der Fehler: die Konfiguration aller Pins vor dem ersten Lesen;
// beim Lesen Pins nach Namen, Quellen in Konfigurations-Reihenfolge (DC-FA-PSYNC-001).
func TestPinsyncFehlerReihenfolge(t *testing.T) {
	files := map[string]string{"ok": "v=1\n"}
	// Pin "a" hätte einen Lesefehler, Pin "z" einen Konfigurationsfehler: erst die Konfiguration.
	pa := psPin("a", psSrc("fehlt-a", `v=(\d)`), psSrc("ok", `v=(\d)`))
	pz := psPin("z", psSrc("ok", `v=(`), psSrc("ok", `v=(\d)`))
	if _, err := psRun(t, files, psCfg(pa, pz)); err == nil || !strings.Contains(err.Error(), `pin "z"`) {
		t.Fatalf("erwartet den Konfigurationsfehler von z, got %v", err)
	}
	// Zwei Lesefehler: der des Pins mit dem kleineren Namen kommt zuerst, im Pin der erste Pfad.
	pb := psPin("b", psSrc("fehlt-b", `v=(\d)`), psSrc("ok", `v=(\d)`))
	pa = psPin("a", psSrc("ok", `v=(\d)`), psSrc("fehlt-a2", `v=(\d)`), psSrc("fehlt-a3", `v=(\d)`))
	if _, err := psRun(t, files, psCfg(pb, pa)); err == nil || !strings.Contains(err.Error(), "fehlt-a2") {
		t.Fatalf("erwartet den Lesefehler von a/fehlt-a2, got %v", err)
	}
}

// Scan-Menge: die Quellen liegen unter scan.ignore, außerhalb von scan.roots
// und ohne Markdown-Endung und werden trotzdem gelesen und geprüft (DC-FA-PSYNC-001).
func TestPinsyncQuellenAusserhalbDerScanMenge(t *testing.T) {
	files := map[string]string{
		"docs/a.md":      "# x\n",
		"vendor/Dockerfile": "ARG GO_VERSION=1.27.3\n",
		"build/Makefile": "GO_VERSION := 1.27.4\n",
	}
	pin := psPin("go", psSrc("vendor/Dockerfile", psDockerPat), psSrc("build/Makefile", psMakePat))
	cfg := model.Config{Pinsync: psCfg(pin), Ignore: []string{"vendor/**", "build/**"}, Roots: []string{"docs"}}
	res, err := Run(coretest.NewMemFS(files), nil, cfg, []string{"pinsync"})
	if err != nil {
		t.Fatalf("Lauf: %v", err)
	}
	if len(res.Findings) != 1 || res.Findings[0].Reason != ReasonPinsyncMismatch || res.Findings[0].File != "build/Makefile" {
		t.Fatalf("Befunde: %+v", res.Findings)
	}
	if len(res.Notes) != 1 || res.Notes[0] != "pinsync: 1 Pin(s) geprüft, 2 Quelle(n) gelesen" {
		t.Fatalf("Notes: %v", res.Notes)
	}
}

// Ein pinsync-Block bei nicht aktivem Modul bleibt wirkungslos: nichts wird
// gelesen, keine Zeile in den Hinweisen, kein Fehler -- auch nicht bei einer
// fehlenden Quelle (DC-FA-PSYNC-001).
func TestPinsyncBlockBeiInaktivemModulIstInert(t *testing.T) {
	pin := psPin("p", psSrc("fehlt", `v=(\d)`), psSrc("fehlt2", `v=(\d)`))
	cfg := model.Config{Pinsync: psCfg(pin)}
	res, err := Run(coretest.NewMemFS(map[string]string{"README.md": "# x\n"}), nil, cfg, []string{"links"})
	if err != nil || len(res.Findings) != 0 || len(res.Notes) != 0 {
		t.Fatalf("inaktiv: err=%v findings=%+v notes=%v", err, res.Findings, res.Notes)
	}
}

// Ein aktives Modul ohne Pin ist ein Eingabefehler des Laufs (fail-closed) (DC-FA-PSYNC-001).
func TestPinsyncAktivOhnePinIstExitZwei(t *testing.T) {
	_, err := Run(coretest.NewMemFS(map[string]string{"README.md": "# x\n"}), nil, model.Config{}, []string{"pinsync"})
	if err == nil || !strings.Contains(err.Error(), "pinsync.pins ist leer") {
		t.Fatalf("erwartet Fehler, got %v", err)
	}
}

// Zwei Läufe auf demselben Baum liefern dieselben Befunde (DC-FA-PSYNC-001).
func TestPinsyncLaeufeSindGleich(t *testing.T) {
	fs := coretest.NewMemFS(psGoFiles("ARG GO_VERSION=1.0\n", "GO_VERSION := 2.0\n"))
	a, _ := CheckPinsync(fs, psCfg(psGo()))
	b, _ := CheckPinsync(fs, psCfg(psGo()))
	if fmt.Sprint(a) != fmt.Sprint(b) || len(a.Findings) != 1 {
		t.Fatalf("Läufe verschieden: %v / %v", a, b)
	}
}

// Grenze: zwei verschiedene, aber gleichwertige Ausdrücke derselben Datei
// vergleichen dieselbe Zeile mit sich selbst -- grün, eine Stelle geprüft
// (DC-FA-PSYNC-001).
func TestPinsyncGleichwertigeAusdruckeSindKeinSchutz(t *testing.T) {
	pin := psPin("p", psSrc("a", `^V=(\S+)`), psSrc("a", `^V=(\S+) *`))
	r := psMust(t, map[string]string{"a": "V=1\n"}, psCfg(pin))
	psWant(t, r)
	if r.Sources != 2 {
		t.Fatalf("Quellen: %d", r.Sources)
	}
}
