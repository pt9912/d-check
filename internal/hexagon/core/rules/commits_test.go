package rules

import (
	"errors"
	"regexp"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// commitsConfig spiegelt die .d-check.yml-commits-Klasse: die drei ids-Muster
// (ADR/MR/DC) plus slice-NNN, Betreff-Ausnahme Merge/Revert (DC-FA-COMMITS-001).
func commitsConfig() model.CommitsConfig {
	return model.CommitsConfig{
		IDPatterns: []*regexp.Regexp{
			regexp.MustCompile(`ADR-\d{4}`),
			regexp.MustCompile(`MR-\d{3}`),
			regexp.MustCompile(`DC-(FA-[A-Z]+|QA)-\d+`),
			regexp.MustCompile(`slice-\d+`),
		},
		ExemptPattern: regexp.MustCompile(`^(Merge |Revert )`),
	}
}

// TestCheckCommitMessage deckt die Selbsttest-Klassen des abgelösten
// tools/trace-check.sh als Modul-Tests ab (ADR-0027): ID erkannt, fehlende ID
// gefangen, Merge/Revert ausgenommen, uniforme #-/scissors-Bereinigung.
func TestCheckCommitMessage(t *testing.T) {
	cfg := commitsConfig()
	cases := []struct {
		name string
		msg  string
		bad  bool
	}{
		{"adr-id", "fix(x): siehe ADR-0001", false},
		{"slice-id", "docs(plan): slice-056 Body", false},
		{"dc-id", "feat: DC-FA-COMMITS-001 Modul", false},
		{"mr-id", "chore: MR-013 move", false},
		{"qa-id", "test: DC-QA-02 Determinismus", false},
		{"no-id", "chore: ohne bezug", true},
		{"merge-exempt", "Merge branch 'x'", false},
		{"revert-exempt", "Revert \"feat: x\"", false},
		{"id-only-in-comment", "chore: ohne bezug\n# ADR-0001 im Kommentar", true},
		{"id-after-scissors", "chore: ohne bezug\n# ------------------------ >8 ------------------------\n# ADR-0001 im Diff", true},
		{"id-in-body-line", "chore: betreff\n\nBody nennt slice-056.", false},
		{"empty-message", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, bad := CheckCommitMessage("abc1234", []byte(tc.msg), cfg)
			if bad != tc.bad {
				t.Fatalf("bad=%v, erwartet %v (msg %q)", bad, tc.bad, tc.msg)
			}
			if !bad {
				return
			}
			if f.Reason != model.ReasonCommitUntraceable || f.Rule != "commits" {
				t.Errorf("Befund falsch: %+v", f)
			}
			if f.File != "abc1234" || f.Target != "abc1234" || f.Line != 1 {
				t.Errorf("File/Target/Line = %q/%q/%d, erwartet abc1234/abc1234/1", f.File, f.Target, f.Line)
			}
		})
	}
}

// TestCheckCommitMessageNoExempt: ohne exempt-pattern ist auch ein Merge-Commit
// kennungspflichtig (die Ausnahme ist Config, nicht fest verdrahtet).
func TestCheckCommitMessageNoExempt(t *testing.T) {
	cfg := commitsConfig()
	cfg.ExemptPattern = nil
	if _, bad := CheckCommitMessage("x", []byte("Merge branch 'x'"), cfg); !bad {
		t.Fatal("ohne exempt-pattern ist auch ein Merge-Commit ID-pflichtig")
	}
}

// TestCheckCommitsRange prüft den Range-Modus über den Fake-Port: nur der
// kennungslose Nicht-Merge-Commit erzeugt einen Befund.
func TestCheckCommitsRange(t *testing.T) {
	cfg := commitsConfig()
	vcs := &fakeVCS{commits: []driven.CommitMeta{
		{ShortSHA: "aaaaaaa", Message: "feat: ADR-0001 ok"},
		{ShortSHA: "bbbbbbb", Message: "chore: kein bezug"},
	}}
	findings, err := CheckCommits(vcs, cfg, "BASE", "HEAD")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(findings) != 1 || findings[0].Target != "bbbbbbb" || findings[0].Reason != model.ReasonCommitUntraceable {
		t.Fatalf("erwartet 1 Befund an bbbbbbb, bekam %+v", findings)
	}
}

// TestCheckCommitsInert: leere id-patterns ⇒ inert, ohne den Port zu berühren
// (der Fake würde bei Aufruf einen Fehler liefern).
func TestCheckCommitsInert(t *testing.T) {
	findings, err := CheckCommits(&fakeVCS{err: errors.New("darf nicht gerufen werden")}, model.CommitsConfig{}, "BASE", "HEAD")
	if err != nil || findings != nil {
		t.Fatalf("inert erwartet, bekam findings=%v err=%v", findings, err)
	}
}

// TestCheckCommitsNilPort: ohne Port ⇒ inert (byte-identisch zum Lauf ohne das Modul).
func TestCheckCommitsNilPort(t *testing.T) {
	findings, err := CheckCommits(nil, commitsConfig(), "BASE", "HEAD")
	if err != nil || findings != nil {
		t.Fatalf("nil-Port inert erwartet, bekam %v/%v", findings, err)
	}
}

// TestCheckCommitsFailClosed: ein Port-Fehler (fehlendes .git/Range) wird
// durchgereicht — der Aufrufer bildet ihn auf Exit 2 ab (fail-closed).
func TestCheckCommitsFailClosed(t *testing.T) {
	if _, err := CheckCommits(&fakeVCS{err: errors.New(".git fehlt")}, commitsConfig(), "BASE", "HEAD"); err == nil {
		t.Fatal("fail-closed: Port-Fehler muss durchgereicht werden")
	}
}

// TestCheckCommitsEmptyRange: eine gültige, leere Range ⇒ keine Befunde (Exit 0).
func TestCheckCommitsEmptyRange(t *testing.T) {
	findings, err := CheckCommits(&fakeVCS{commits: nil}, commitsConfig(), "BASE", "HEAD")
	if err != nil || len(findings) != 0 {
		t.Fatalf("leere Range ⇒ keine Befunde, bekam %v/%v", findings, err)
	}
}
