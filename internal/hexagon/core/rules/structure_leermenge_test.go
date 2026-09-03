package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// acDrei traegt drei gleichartige Abschnitte — die Gestalt des Bestands, um
// den der Antrag gestellt wurde: alle sind grandfathert, keiner ist neu.
const acDrei = "# Lastenheft\n\n" +
	"### AC-001 — alt\n\nHappy: ja.\n\n" +
	"### AC-002 — alt\n\nHappy: ja.\n\n" +
	"### AC-003 — alt\n\nHappy: ja.\n"

func leerRule() model.StructureRule {
	return model.StructureRule{
		Files: "docs/*.md", SectionPattern: `^### AC-`, Sections: "each",
		RequireAll:           []string{"Boundary"},
		ExemptSectionPattern: `^### AC-00[123]\b`,
	}
}

// DIE DEKLARIERTE LEERMENGE IST STUMM — und die Umkehr steht in derselben
// Funktion: ohne die Zahl meldet dieselbe Eingabe section-missing (BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt).
func TestExemptExpectCount_DeklarierteLeermengeIstStumm(t *testing.T) {
	r := leerRule()
	f := laufe(t, acDrei, r)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionMissing {
		t.Fatalf("VORZUSTAND: ohne die Zahl meldet die Nullmengen-Haerte, got %+v", f)
	}
	r.ExemptExpectCount = ptr(3)
	if f := laufe(t, acDrei, r); f != nil {
		t.Fatalf("drei von drei ausgenommen, drei deklariert ⇒ befundfrei, got %+v", f)
	}
}

// DIE DRIFT IST BEIDSEITIG. Eine veraltete Aufzaehlung nimmt zu wenig, eine
// erweiterte zu viel — beides ist dieselbe Luecke, nur in andere Richtung, und
// eine einseitige Pruefung waere ein halber Waechter.
func TestExemptExpectCount_DriftIstBeidseitig(t *testing.T) {
	for name, tc := range map[string]struct {
		muster string
		zahl   int
		will   string
	}{
		"zu wenig ausgenommen": {`^### AC-00[12]\b`, 3, "nimmt 2 von 3 Abschnitten aus, deklariert sind 3"},
		"zu viel ausgenommen":  {`^### AC-00[123]\b`, 2, "nimmt 3 von 3 Abschnitten aus, deklariert sind 2"},
	} {
		t.Run(name, func(t *testing.T) {
			r := leerRule()
			r.ExemptSectionPattern = tc.muster
			r.ExemptExpectCount = ptr(tc.zahl)
			f := laufe(t, acDrei, r)
			if len(f) != 1 || f[0].Reason != model.ReasonSectionExemptMismatch {
				t.Fatalf("erwartet section-exempt-mismatch, got %+v", f)
			}
			if !strings.Contains(f[0].Message, tc.will) {
				t.Fatalf("Meldung nennt die Zahlen nicht: %q", f[0].Message)
			}
		})
	}
}

// Die Zahl wird auch dann geprueft, wenn NOCH ABSCHNITTE UEBRIG sind — sonst
// bliebe genau der Fall stumm, den der Antrag als einzigen benennt: die
// Aufzaehlung veraltet, waehrend der Bestand waechst.
func TestExemptExpectCount_PrueftAuchBeiRestmenge(t *testing.T) {
	body := acDrei + "\n### AC-042 — neu\n\nHappy: ja.\n"
	r := leerRule()
	r.ExemptExpectCount = ptr(3)
	f := laufe(t, body, r)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionMarkerMissing {
		t.Fatalf("die Zahl stimmt ⇒ nur der neue Abschnitt meldet, got %+v", f)
	}
	r.ExemptExpectCount = ptr(4)
	f = laufe(t, body, r)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionExemptMismatch {
		t.Fatalf("falsche Zahl bei Restmenge muss melden, got %+v", f)
	}
}

// DIE TRENNUNG HAELT: trifft schon der Selektor nichts, ist das ein
// Konfigurationsdefekt und bleibt section-missing — auch mit gesetzter Zahl.
// Genau diese Trennung ist der Gegenstand des Antrags.
func TestExemptExpectCount_SelektorOhneTrefferBleibtDefekt(t *testing.T) {
	r := leerRule()
	r.SectionPattern = `^### GIBTESNICHT-`
	r.ExemptExpectCount = ptr(0)
	f := laufe(t, acDrei, r)
	if len(f) != 1 || f[0].Reason != model.ReasonSectionMissing {
		t.Fatalf("Selektor ohne Treffer bleibt section-missing, got %+v", f)
	}
	if !strings.Contains(f[0].Message, "kein Abschnitt passt auf den Selektor") {
		t.Fatalf("die Meldung muss den Selektor nennen, got %q", f[0].Message)
	}
}

// Eine deklarierte NULL bedeutet etwas: „das Muster soll heute noch nichts
// treffen" — ein Bestand, der erst waechst. Sie ist deshalb erlaubt und von
// „nicht deklariert" unterscheidbar.
func TestExemptExpectCount_NullIstEineAussage(t *testing.T) {
	r := leerRule()
	r.ExemptSectionPattern = `^### AC-9\d\d\b`
	r.ExemptExpectCount = ptr(0)
	f := laufe(t, acDrei, r)
	if len(f) != 3 {
		t.Fatalf("null ausgenommen, null deklariert ⇒ alle drei werden geprueft, got %+v", f)
	}
	for _, x := range f {
		if x.Reason != model.ReasonSectionMarkerMissing {
			t.Fatalf("erwartet nur Marken-Befunde, got %+v", x)
		}
	}
	// Die Gegenrichtung ist die schaerfere: eine deklarierte NULL, die das
	// Muster widerlegt, MUSS melden. Ohne sie liesse sich die Zaehlung auf
	// `> 0` einschraenken, ohne dass ein Test rot wird -- die Null waere dann
	// von „nicht deklariert" ununterscheidbar, und genau das ist ihr Zweck.
	r.ExemptSectionPattern = `^### AC-001\b`
	if f := laufe(t, acDrei, r); len(f) != 1 || f[0].Reason != model.ReasonSectionExemptMismatch {
		t.Fatalf("null deklariert, einer ausgenommen ⇒ Mismatch, got %+v", f)
	}
}

// Der neue Befund traegt den verfassten Hinweis — anders als der
// Nullmengen-Befund, der zur Klasse „die Regel hat nicht gemessen" gehoert.
// Hier HAT sie gemessen und eine Deklaration widerlegt.
func TestExemptExpectCount_HintGiltFuerDenMismatch(t *testing.T) {
	r := leerRule()
	r.ExemptExpectCount = ptr(2)
	r.Hint = "Aufzaehlung oder Zahl nachziehen"
	f := laufe(t, acDrei, r)
	if len(f) != 1 || f[0].Message != "Aufzaehlung oder Zahl nachziehen" {
		t.Fatalf("der Hinweis muss den Mismatch erklaeren duerfen, got %+v", f)
	}
}

// Zwei Regeln, die sich NUR in der erwarteten Zahl unterscheiden, sind ein
// Widerspruch und kein Paar — die Identitaet trennt sie bewusst nicht.
func TestExemptExpectCount_TraegtNichtDieRegelIdentitaet(t *testing.T) {
	a := leerRule()
	b := leerRule()
	a.ExemptExpectCount = ptr(3)
	b.ExemptExpectCount = ptr(2)
	if a.Identity() != b.Identity() {
		t.Fatalf("zwei erwartete Zahlen ueber derselben Regel sind ein Widerspruch, kein Paar: %q vs %q",
			a.Identity(), b.Identity())
	}
}
