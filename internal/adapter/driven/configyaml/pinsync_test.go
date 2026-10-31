package configyaml_test

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/adapter/driven/configyaml"
)

const psYAML = "pinsync:\n" +
	"  pins:\n" +
	"    - name: go-version\n" +
	"      sources:\n" +
	"        - {file: Dockerfile, pattern: '^ARG GO_VERSION=v?(\\S+)'}\n" +
	"        - {file: Makefile, pattern: '^GO_VERSION\\s*:?=\\s*v?(\\S+)'}\n"

// Der Block wird gelesen (DC-FA-PSYNC-001, DC-FA-CONF-001).
func TestDecode_PinsyncVollblock(t *testing.T) {
	cfg, err := configyaml.Decode([]byte(psYAML))
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	pins := cfg.Pinsync.Pins
	if len(pins) != 1 || pins[0].Name != "go-version" || len(pins[0].Sources) != 2 {
		t.Fatalf("Pins nicht gelesen: %+v", pins)
	}
	if pins[0].Sources[1].File != "Makefile" || pins[0].Sources[1].Pattern != `^GO_VERSION\s*:?=\s*v?(\S+)` {
		t.Fatalf("Quelle nicht gelesen: %+v", pins[0].Sources[1])
	}
}

// Ohne Block bleibt die Konfiguration leer; ein leerer Block ist erst im
// aktiven Modul ein Fehler (DC-FA-PSYNC-001).
func TestDecode_PinsyncOhneBlockUndLeerePinListe(t *testing.T) {
	for _, in := range []string{"modules: [links]\n", "pinsync:\n  pins: []\n"} {
		cfg, err := configyaml.Decode([]byte(in))
		if err != nil || len(cfg.Pinsync.Pins) != 0 {
			t.Fatalf("%q: err=%v pins=%+v", in, err, cfg.Pinsync.Pins)
		}
	}
}

// Jede Randbedingung der Konfiguration endet am Rand mit einem Fehler (Exit 2) (DC-FA-PSYNC-001).
func TestDecode_PinsyncEingabefehler(t *testing.T) {
	src := func(f, p string) string { return "{file: " + f + ", pattern: '" + p + "'}" }
	block := func(name string, srcs ...string) string {
		return "pinsync:\n  pins:\n    - name: " + name + "\n      sources: [" + strings.Join(srcs, ", ") + "]\n"
	}
	ok := src("a", `v=(\d)`)
	cases := map[string]struct{ in, want string }{
		"eine Quelle":          {block("p", ok), "mindestens zwei"},
		"Pfad mit ..":          {block("p", src("../a", `v=(\d)`), ok), "segment"},
		"abschließender Strich": {block("p", src("a/", `v=(\d)`), ok), "endet auf '/'"},
		"keine Fanggruppe":     {block("p", src("a", `v=\d`), ok), "0 Fanggruppen"},
		"zwei Fanggruppen":     {block("p", src("a", `(v)=(\d)`), ok), "2 Fanggruppen"},
		"Ausdruck ungültig":    {block("p", src("a", `v=(\d`), ok), "RE2"},
		"dieselbe Stelle":      {block("p", ok, ok), "gegen sich selbst"},
		"unbekannter Schlüssel": {"pinsync:\n  pins: []\n  extra: 1\n", "extra"},
		"unbekannter Quellenschlüssel": {"pinsync:\n  pins:\n    - name: p\n      sources: [{file: a, pattern: 'v=(\\d)', flag: x}, " + ok + "]\n", "flag"},
		"doppelter Name": {"pinsync:\n  pins:\n    - {name: p, sources: [" + ok + ", " + src("b", `v=(\d)`) + "]}\n    - {name: p, sources: [" + ok + ", " + src("b", `v=(\d)`) + "]}\n", "mehrfach"},
		"Pin ohne Namen": {"pinsync:\n  pins:\n    - sources: [" + ok + ", " + src("b", `v=(\d)`) + "]\n", "keinen Namen"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := configyaml.Decode([]byte(c.in))
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("erwartet Fehler mit %q, got %v", c.want, err)
			}
		})
	}
}
