package rules

import (
	"reflect"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// hostTargets liefert die gemeldeten Ziele von hostpaths über eine Datei.
func hostTargets(t *testing.T, content string, cfg model.HostpathsConfig) []string {
	t.Helper()
	m := coretest.NewMemFS(map[string]string{"docs/v.md": content})
	res, err := Run(m, nil, model.Config{Roots: []string{"docs"}, Hostpaths: cfg}, []string{"hostpaths"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		got = append(got, f.Target)
	}
	return got
}

func homeVars() model.HostpathsConfig { return model.HostpathsConfig{Variables: []string{"HOME"}} }

// Reihenfolge der Ziele: die des Moduls (Sortierung je Datei), nicht die der Zeile.
// Happy: alle drei Formen aus dem Feld melden mit gesetztem variables, in
// Prosa und in Inline-Code; das schließende Anführungszeichen gehört nicht
// zum gemeldeten Pfad.
func TestHostpathsVariablesHappy(t *testing.T) {
	got := hostTargets(t,
		"Log in $HOME/Dev/x/y.log ablegen\n"+
			"und ${HOME}/Dev/z sowie \"$HOME\"/Dev/w\n"+
			"im Code: `$HOME/code/z`\n", homeVars())
	want := []string{"$HOME/Dev/x/y.log", "$HOME/Dev/w", "${HOME}/Dev/z", "$HOME/code/z"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Targets = %q\nwant   %q", got, want)
	}
}

// Boundary: ohne den Schlüssel meldet nichts, auch nicht dieselben Zeilen;
// Werkzeug-Konventionen (erstes Segment mit Punkt) und die nackte Variable
// bleiben auch mit dem Schlüssel still.
func TestHostpathsVariablesBoundary(t *testing.T) {
	content := "Log in $HOME/Dev/x und ${HOME}/Dev/z und \"$HOME\"/Dev/w\n"
	if got := hostTargets(t, content, model.HostpathsConfig{}); len(got) != 0 {
		t.Fatalf("ohne variables kein Fund erwartet, got %q", got)
	}
	if got := hostTargets(t, "Konfig in $HOME/.config/x und nacktes $HOME hier\n", homeVars()); len(got) != 0 {
		t.Fatalf("Punkt-Segment und nackte Variable bleiben still, got %q", got)
	}
}

// Negativliste: kein Fund für längere Namen, Default-Syntax, Wortgrenze,
// Doppel-`$`, Maskierung, URL, nicht gelistete Variable, Fence, Zeilenumbruch
// und doppelten Schrägstrich.
func TestHostpathsVariablesNegatives(t *testing.T) {
	content := "$HOMEDIR/x/y\n" +
		"${HOME:-/tmp}/x/y\n" +
		"${HOME}x/y/z\n" +
		"x$HOME/y/z\n" +
		"$$HOME/x/y\n" +
		"\\$HOME/x/y\n" +
		"https://h.example/$HOME/y/z\n" +
		"_$HOME/x/y_\n" +
		"$OTHER/x/y\n" +
		"$HOME//x/y\n" +
		"$HOME\n/x/y\n" +
		"```\ncd $HOME/src/fence\n```\n"
	if got := hostTargets(t, content, homeVars()); len(got) != 0 {
		t.Fatalf("kein Fund erwartet, got %q", got)
	}
}

// Die Hervorhebung am Rand verhält sich wie bei der Tilde: `**$HOME/x/y**`
// meldet mit den Stern-Zeichen (benannte Grenze der Erkennung, keine neue).
func TestHostpathsVariablesEmphasisLikeTilde(t *testing.T) {
	got := hostTargets(t, "**$HOME/x/y** und **~/x/y**\n", homeVars())
	want := []string{"$HOME/x/y**", "~/x/y**"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Targets = %q\nwant   %q", got, want)
	}
}

// Mehrere Variablen, und eine nicht gelistete bleibt still.
func TestHostpathsVariablesListe(t *testing.T) {
	cfg := model.HostpathsConfig{Variables: []string{"HOME", "WORK"}}
	got := hostTargets(t, "$WORK/a/b und ${WORK}/c/d und $HOME/e/f und $OTHER/g/h\n", cfg)
	want := []string{"$HOME/e/f", "$WORK/a/b", "${WORK}/c/d"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Targets = %q\nwant   %q", got, want)
	}
}


// Variablen-Funde nehmen dem Unix-Muster nichts: ein Unix-Pfad in derselben
// Zeile meldet weiter.
func TestHostpathsVariablesOhneUnixVerlust(t *testing.T) {
	got := hostTargets(t, "$HOME/Dev/x und /home/alice/y\n", homeVars())
	want := []string{"$HOME/Dev/x", "/home/alice/y"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Targets = %q\nwant   %q", got, want)
	}
}

// exempt-targets mit `$`-Glob nimmt genau diesen Fund aus, auch in der Form
// `"$HOME"/…` (das Anführungszeichen ist nicht Teil des Ziels).
func TestHostpathsVariablesExemptTargets(t *testing.T) {
	cfg := model.HostpathsConfig{
		Variables:     []string{"HOME"},
		ExemptTargets: []string{"$HOME/Dev/**", "${HOME}/Lib/**"},
	}
	got := hostTargets(t, "\"$HOME\"/Dev/a/b und $HOME/Dev/c/d und ${HOME}/Lib/e und $HOME/src/f\n", cfg)
	want := []string{"$HOME/src/f"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Targets = %q\nwant   %q", got, want)
	}
}

// Ein Unix-Treffer, der innerhalb eines Variablen-Treffers beginnt, entfällt:
// `$HOME/a=/mnt/x` meldet genau einmal und in voller Form. (Bei einem
// Präfix-Namen als erstem Segment, `$HOME/Development/x`, träfe das Unix-Muster
// ohnehin nicht — vor `/Development/` steht das `E` von `HOME`.)
func TestHostpathsVariablesVerdraengtUnix(t *testing.T) {
	got := hostTargets(t, "Pfad $HOME/a=/mnt/x/y und $HOME/Development/d-check/README.md\n", homeVars())
	want := []string{"$HOME/Development/d-check/README.md", "$HOME/a=/mnt/x/y"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Targets = %q\nwant   %q", got, want)
	}
}

// Ist das erste Segment des Variablen-Pfads ein Präfix-Name, trifft das
// Unix-Muster bei `}` oder `"` vor dem Präfix (beide stehen nicht in seiner
// Sperrklasse): die Verdrängung sorgt für genau einen Fund in voller Form. Bei
// `$HOME/…` verhindert schon das `E` des Namens den Doppeltreffer.
func TestHostpathsVariablesPraefixAlsErstesSegment(t *testing.T) {
	got := hostTargets(t,
		"${HOME}/Development/a/b und \"$HOME\"/Development/c/d und $HOME/Development/e/f\n", homeVars())
	want := []string{"$HOME/Development/c/d", "$HOME/Development/e/f", "${HOME}/Development/a/b"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Targets = %q\nwant   %q", got, want)
	}
}

// Überlappende Treffer in einer Zeile (ein Komma ist Pfadzeichen und Wortgrenze
// zugleich) melden genau einmal, der früher beginnende gewinnt — in allen
// Richtungen: Tilde in Variable, Variable in Tilde, Variable in Unix.
func TestHostpathsVariablesUeberlappungGenauEinmal(t *testing.T) {
	got := hostTargets(t,
		"$HOME/a/b,~/c/d\n"+
			"~/e/f,$HOME/g/h\n"+
			"/home/i/j,$HOME/k/l\n", homeVars())
	want := []string{"$HOME/a/b,~/c/d", "~/e/f,$HOME/g/h", "/home/i/j,$HOME/k/l"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Targets = %q\nwant   %q", got, want)
	}
}

// Mit gesetztem variables wechselt das Ziel eines bisherigen Unix-Funds hinter
// einer Variablen: er meldet in voller Form, ein Bestands-Glob auf den
// Unix-Teil trifft dann nicht mehr; das Glob muss mit der Variablen beginnen.
func TestHostpathsVariablesZielWechselt(t *testing.T) {
	line := "${HOME}/Development/x/y\n"
	if got := hostTargets(t, line, model.HostpathsConfig{}); !reflect.DeepEqual(got, []string{"/Development/x/y"}) {
		t.Fatalf("ohne variables bleibt der Unix-Fund: %q", got)
	}
	cfg := model.HostpathsConfig{Variables: []string{"HOME"}, ExemptTargets: []string{"/Development/**"}}
	if got := hostTargets(t, line, cfg); !reflect.DeepEqual(got, []string{"${HOME}/Development/x/y"}) {
		t.Fatalf("mit variables meldet die volle Form, das Unix-Glob trifft nicht: %q", got)
	}
	cfg.ExemptTargets = []string{"${HOME}/Development/**"}
	if got := hostTargets(t, line, cfg); len(got) != 0 {
		t.Fatalf("das Variablen-Glob nimmt den Fund aus: %q", got)
	}
}
