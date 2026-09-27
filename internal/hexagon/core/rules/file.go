package rules

import (
	"sort"
	"strconv"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Grund-Codes des Moduls file (DC-FA-FILE-001, spec/spezifikation.md §4).
const (
	// ReasonFileNoMatch: die Regel trifft keine Datei (auch nach Abzug von
	// exempt-paths) — dieselbe Nullmengen-Härte wie section-missing bei
	// structure, nur ohne Abschnitts-Konzept.
	ReasonFileNoMatch = "file-no-match"
	// ReasonFileLinesExceeded: die Datei hat mehr Zeilen, als max-lines
	// erlaubt.
	ReasonFileLinesExceeded = "file-lines-exceeded"
	// ReasonFileBytesExceeded: die Datei ist größer, als max-bytes erlaubt.
	ReasonFileBytesExceeded = "file-bytes-exceeded"
)

// CheckFile prüft Zeilen- und Byte-Obergrenzen GANZER Dateien
// (DC-FA-FILE-001, §DC-FA-FILE-001.a). Post-Pass wie structure: die Regeln
// benennen ihre Dateien selbst über eigene Globs, unabhängig von
// scan.roots/scan.ignore — deshalb kennt das Modul kein <modul>.scope.
// Anders als structure ist die Kandidaten-Menge NICHT auf Markdown
// beschränkt: eine Zeilen-/Byte-Zahl gibt es für jede Dateiart.
func CheckFile(fsys driven.Filesystem, rules []model.FileRule) []model.Finding {
	if len(rules) == 0 {
		return nil
	}
	all, err := fileTree(fsys)
	if err != nil {
		// Wie structure: der Dateibaum ist für KEINE Regel messbar. Je Regel
		// ein Befund, damit die Identität der unmessbaren Behauptung erhalten
		// bleibt (kein zusammengefasster Sammel-Befund).
		out := make([]model.Finding, 0, len(rules))
		for _, r := range rules {
			out = append(out, fileRawFinding(r, r.Files, 1, ReasonFileNoMatch,
				"Dateibaum nicht lesbar ("+err.Error()+") — Regel nicht messbar (fail-closed)"))
		}
		return out
	}
	var out []model.Finding
	for _, r := range rules {
		out = append(out, checkFileRule(fsys, r, all)...)
	}
	return out
}

// fileTree listet ALLE Dateien des Baums (SKIP_DIRS gelten, die Scan-
// Wurzeln nicht) — die Grundmenge, gegen die die Regel-Globs laufen.
func fileTree(fsys driven.Filesystem) ([]string, error) {
	var out []string
	if err := walkAllFiles(fsys, "", &out); err != nil {
		return nil, err
	}
	sort.Strings(out)
	return out, nil
}

// walkAllFiles wandert den ganzen Baum und sammelt JEDE Datei (anders als
// walkMarkdown, das auf .md filtert) — SKIP_DIRS gelten unverändert.
func walkAllFiles(fsys driven.Filesystem, dir string, out *[]string) error {
	entries, err := fsys.List(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		rel := e.Name
		if dir != "" {
			rel = dir + "/" + e.Name
		}
		switch e.Kind {
		case driven.KindDir:
			if isSkipDir(e.Name) {
				continue
			}
			if err := walkAllFiles(fsys, rel, out); err != nil {
				return err
			}
		case driven.KindFile:
			*out = append(*out, rel)
		}
	}
	return nil
}

// checkFileRule wertet eine Regel aus: Kandidaten (Glob minus exempt-paths),
// dann je Datei die Schwellen (DC-FA-FILE-001.a).
func checkFileRule(fsys driven.Filesystem, r model.FileRule, all []string) []model.Finding {
	var cands []string
	for _, f := range all {
		if !matchGlob(r.Files, f) {
			continue
		}
		if fileExempt(r, f) {
			continue
		}
		cands = append(cands, f)
	}
	// Nullmengen-Härte: eine Regel zu setzen IST die Behauptung, dass sie
	// Dateien trifft — auch dann, wenn erst exempt-paths die Menge geleert
	// hat (wie bei structure).
	if len(cands) == 0 {
		return []model.Finding{fileRawFinding(r, r.Files, 1, ReasonFileNoMatch,
			"Regel trifft keine Datei (auch nach Abzug von exempt-paths) — das Gate liefe leer")}
	}
	var out []model.Finding
	for _, f := range cands {
		out = append(out, checkFileOne(r, fsys, f)...)
	}
	return out
}

func fileExempt(r model.FileRule, file string) bool {
	for _, g := range r.ExemptPaths {
		if matchGlob(g, file) {
			return true
		}
	}
	return false
}

// checkFileOne prüft eine Kandidaten-Datei gegen die Schwellen der Regel.
// Eine unlesbare Einzeldatei ist wie bei structure FAIL-CLOSED: sie still zu
// überspringen wäre genau der Grün-Pfad, den die Nullmengen-Härte an anderer
// Stelle ausschließt — ein Repo, das gerade deshalb rot werden soll, weil eine
// Datei nicht mehr lesbar ist, bliebe sonst grün. Derselbe Grund-Code wie bei
// leerer Kandidaten-Menge: die Regel hat hier nicht gemessen, kein Hint.
func checkFileOne(r model.FileRule, fsys driven.Filesystem, file string) []model.Finding {
	content, err := fsys.ReadFile(file)
	if err != nil {
		return []model.Finding{fileRawFinding(r, file, 1, ReasonFileNoMatch,
			"Datei ist unlesbar (fail-closed)")}
	}
	var out []model.Finding
	if r.MaxLines != nil {
		if n := countLines(content); n > *r.MaxLines {
			out = append(out, fileFinding(r, file, 1, ReasonFileLinesExceeded,
				"Datei hat "+strconv.Itoa(n)+" Zeilen, erlaubt sind "+strconv.Itoa(*r.MaxLines)))
		}
	}
	if r.MaxBytes != nil {
		if n := len(content); n > *r.MaxBytes {
			out = append(out, fileFinding(r, file, 1, ReasonFileBytesExceeded,
				"Datei hat "+strconv.Itoa(n)+" Bytes, erlaubt sind "+strconv.Itoa(*r.MaxBytes)))
		}
	}
	return out
}

func fileRawFinding(r model.FileRule, file string, line int, reason, msg string) model.Finding {
	return model.Finding{
		File: file, Line: line, Rule: "file",
		Target: r.Identity(), Reason: reason, Message: msg,
	}
}

func fileFinding(r model.FileRule, file string, line int, reason, msg string) model.Finding {
	return model.Finding{
		File: file, Line: line, Rule: "file",
		Target: r.Identity(), Reason: reason, Message: r.MessageFor(msg),
	}
}
