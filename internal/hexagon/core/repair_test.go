package core

import (
	"regexp"
	"testing"
)

func adrCfg() Config {
	return Config{IDPatterns: []IDPattern{
		{Regex: regexp.MustCompile(`ADR-\d{4}`), Target: "docs/plan/adr/"},
	}}
}

// Konservativ: ein nacktes Prosa-id-unlinked wird zum Definitions-Link;
// genau ein Edit, kein ReviewRequired.
func TestRepairEdits_Conservative_IDUnlinked(t *testing.T) {
	fsys := newMemFS(map[string]string{"docs/a.md": "nacktes ADR-0042 hier\n"})
	f := Finding{File: "docs/a.md", Line: 1, Rule: "ids", Target: "ADR-0042", Reason: ReasonIDUnlinked}
	edits, err := RepairEdits(fsys, []Finding{f}, adrCfg(), false)
	if err != nil {
		t.Fatal(err)
	}
	if len(edits) != 1 {
		t.Fatalf("erwartet 1 Edit, got %d: %+v", len(edits), edits)
	}
	if edits[0].NewLine != "nacktes [`ADR-0042`](plan/adr) hier" || edits[0].ReviewRequired {
		t.Fatalf("Edit unerwartet: %+v", edits[0])
	}
}

// Sicherheit: ein id-unlinked-Vorkommen INNERHALB von Inline-Code wird
// NICHT repariert (sonst zerrissener Code-Span). Der vorverarbeitete Text
// hat den Span geleert → keine nackte Fundstelle.
func TestRepairEdits_InCode_NichtRepariert(t *testing.T) {
	fsys := newMemFS(map[string]string{"docs/a.md": "siehe `ADR-0042` im Code\n"})
	f := Finding{File: "docs/a.md", Line: 1, Rule: "ids", Target: "ADR-0042", Reason: ReasonIDUnlinked}
	edits, err := RepairEdits(fsys, []Finding{f}, adrCfg(), false)
	if err != nil {
		t.Fatal(err)
	}
	if len(edits) != 0 {
		t.Fatalf("Inline-Code darf nicht repariert werden, got %+v", edits)
	}
}

// Sicherheit: eine über-matchende Kennung (ADR-0001 in ADR-00012) wird
// NICHT repariert — die Wortgrenzen-Prüfung verhindert eine falsch
// platzierte Ersetzung.
func TestRepairEdits_Overmatch_NichtRepariert(t *testing.T) {
	fsys := newMemFS(map[string]string{"docs/a.md": "siehe ADR-00012 hier\n"})
	f := Finding{File: "docs/a.md", Line: 1, Rule: "ids", Target: "ADR-0001", Reason: ReasonIDUnlinked}
	edits, err := RepairEdits(fsys, []Finding{f}, adrCfg(), false)
	if err != nil {
		t.Fatal(err)
	}
	if len(edits) != 0 {
		t.Fatalf("über-matchende Kennung darf nicht repariert werden, got %+v", edits)
	}
}

// Konservativ lässt target-missing unangetastet (kein eindeutiger Fix).
func TestRepairEdits_Conservative_SkipTargetMissing(t *testing.T) {
	fsys := newMemFS(map[string]string{"docs/a.md": "[x](alt.md)\n", "docs/sub/alt.md": "ok\n"})
	f := Finding{File: "docs/a.md", Line: 1, Rule: "links", Target: "alt.md", Reason: ReasonTargetMissing}
	edits, err := RepairEdits(fsys, []Finding{f}, Config{}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(edits) != 0 {
		t.Fatalf("konservativ darf target-missing nicht anfassen, got %+v", edits)
	}
}

// Breit: target-missing → eindeutige Basisnamen-Datei, ReviewRequired.
func TestRepairEdits_Broad_TargetMissing(t *testing.T) {
	fsys := newMemFS(map[string]string{"docs/a.md": "[x](alt.md)\n", "docs/sub/alt.md": "ok\n"})
	f := Finding{File: "docs/a.md", Line: 1, Rule: "links", Target: "alt.md", Reason: ReasonTargetMissing}
	edits, err := RepairEdits(fsys, []Finding{f}, Config{}, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(edits) != 1 {
		t.Fatalf("erwartet 1 Best-Guess-Edit, got %d: %+v", len(edits), edits)
	}
	if edits[0].NewLine != "[x](sub/alt.md)" || !edits[0].ReviewRequired {
		t.Fatalf("Best-Guess-Edit unerwartet: %+v", edits[0])
	}
}

// Breit ohne eindeutigen Basisnamen-Treffer → kein Edit (kein Raten ins
// Leere).
func TestRepairEdits_Broad_KeinEindeutigerTreffer(t *testing.T) {
	fsys := newMemFS(map[string]string{"docs/a.md": "[x](weg.md)\n"})
	f := Finding{File: "docs/a.md", Line: 1, Rule: "links", Target: "weg.md", Reason: ReasonTargetMissing}
	edits, err := RepairEdits(fsys, []Finding{f}, Config{}, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(edits) != 0 {
		t.Fatalf("ohne eindeutigen Treffer kein Edit, got %+v", edits)
	}
}
