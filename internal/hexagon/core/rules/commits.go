package rules

import (
	"regexp"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// commitScissorsRE erkennt die git-`strip`-scissors-Zeile (`# ----- >8 -----`,
// verbose-Diff): ab ihr entfällt der Rest der Message (DC-FA-COMMITS-001.a
// Schritt 2). Package-Var wie die übrigen Modul-Regexe (immutable/pins/matrix).
var commitScissorsRE = regexp.MustCompile(`^#.*>8`)

// CheckCommits ist das Regelmodul commits (DC-FA-COMMITS-001): es prüft, dass
// jede Nicht-Merge-Commit-Message der Range base..head eine Traceability-Kennung
// (cfg.IDPatterns) auf einer Inhalts-Zeile trägt; sonst `commit-untraceable`.
// Opt-in (leere IDPatterns ⇒ inert); es liest die Commit-Messages über den
// VCS-Port (nicht-hermetisch, aber lokal/lesend/deterministisch) — ein Port-Fehler
// (fehlendes `.git`/Range) wird als error zurückgegeben, den der Aufrufer auf
// Exit 2 abbildet (fail-closed). Diagnose-only: kein `--repair`-Hunk.
func CheckCommits(vcs driven.VCS, cfg model.CommitsConfig, base, head string) ([]model.Finding, error) {
	if len(cfg.IDPatterns) == 0 || vcs == nil {
		return nil, nil // inert (Modul ohne ID-Muster oder ohne Port)
	}
	metas, err := vcs.CommitMessages(base, head)
	if err != nil {
		return nil, err
	}
	var findings []model.Finding
	for _, m := range metas {
		if f, bad := CheckCommitMessage(m.ShortSHA, []byte(m.Message), cfg); bad {
			findings = append(findings, f)
		}
	}
	return findings, nil
}

// CheckCommitMessage prüft eine einzelne **rohe** Commit-Message; label ist der
// Ort (Commit-Kurz-SHA im Range-Modus, "pending" im --commit-msg-Kurzschluss-Modus).
// bad=true mit dem Befund `commit-untraceable`, wenn die bereinigte Message keine
// Kennung (cfg.IDPatterns) auf einer Inhalts-Zeile trägt und ihr Betreff nicht per
// ExemptPattern ausgenommen ist (DC-FA-COMMITS-001.a Schritte 2–4). Exportiert für
// den --commit-msg-Modus des CLI.
func CheckCommitMessage(label string, raw []byte, cfg model.CommitsConfig) (model.Finding, bool) {
	cleaned := cleanCommitMessage(string(raw))
	subject := commitSubject(cleaned)
	if cfg.ExemptPattern != nil && cfg.ExemptPattern.MatchString(subject) {
		return model.Finding{}, false
	}
	for _, ln := range strings.Split(cleaned, "\n") {
		for _, re := range cfg.IDPatterns {
			if re.MatchString(ln) {
				return model.Finding{}, false
			}
		}
	}
	return model.Finding{
		File: label, Line: 1, Rule: "commits", Target: label,
		Reason:  model.ReasonCommitUntraceable,
		Message: subject,
	}, true
}

// cleanCommitMessage bereinigt eine Message wie git-`strip` (DC-FA-COMMITS-001.a
// Schritt 2): alles ab der ersten scissors-Zeile (`^#.*>8`) entfällt, danach alle
// `#`-Kommentarzeilen. Uniform in beiden Quellen (Range + --commit-msg), damit sie
// dieselbe (kommentar-bereinigte) Message bewerten — keine Divergenz je
// git-Cleanup-Modus.
func cleanCommitMessage(raw string) string {
	var out []string
	for _, ln := range strings.Split(raw, "\n") {
		if commitScissorsRE.MatchString(ln) {
			break
		}
		if strings.HasPrefix(ln, "#") {
			continue
		}
		out = append(out, ln)
	}
	return strings.Join(out, "\n")
}

// commitSubject liefert den Betreff (erste Zeile) der bereinigten Message
// (DC-FA-COMMITS-001.a Schritt 3; die Ausnahme prüft nur den Betreff).
func commitSubject(cleaned string) string {
	if i := strings.IndexByte(cleaned, '\n'); i >= 0 {
		return cleaned[:i]
	}
	return cleaned
}
