package rules

// Modul metrics (DC-FA-MET-001): Kennzahl eines Planungsdokuments. Dieser Teil
// ist hermetisch -- er liest nur Plan-Texte und rechnet ganzzahlig, ohne git.
//
// GRENZE, ausgesprochen (AGENTS.md §3.8): Die Plan-Dateien werden als Text
// gelesen -- dort gilt die Zusage. Das Plan-Verzeichnis ist ausdruecklich
// benannt und wird nicht auf scan.roots beschraenkt; scan.ignore und die
// Skip-Liste gelten fuer dieses Verzeichnis nicht.
//
// FAIL-CLOSED: ein Eingabefehler eines Plans ist ein error (Exit 2) und
// beendet den Lauf beim ersten Fehler in der Reihenfolge der Plan-Namen; kein
// Befund aus keinem Plan erscheint.

import (
	"fmt"
	"math"
	"path"
	"sort"
	"strings"

	"github.com/pt9912/d-check/internal/hexagon/core/model"
	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// Grund-Codes des Moduls metrics (spec/spezifikation.md §4).
const (
	ReasonMetricsLimitMismatch    = "metrics-limit-mismatch"
	ReasonMetricsReasonEmpty      = "metrics-reason-empty"
	ReasonMetricsCutReasonMissing = "metrics-cut-reason-missing"
	ReasonMetricsOverLimit        = "metrics-over-limit"
)

// MetricsResult traegt die Befunde, die Hinweise (keine Befunde) und die
// Bezugsmenge des Laufs; Hinweise und Bezugsmenge gehoeren in die
// Zusammenfassung (DC-FA-CLI-004 sagt den Wortlaut von message nicht zu).
type MetricsResult struct {
	Findings  []model.Finding
	Hints     []string
	Plans     int
	WithBlock int
}

// Notes liefert die Bezugsmenge, gefolgt von den Hinweisen in
// Plan-Reihenfolge (DC-QA-02).
func (r MetricsResult) Notes() []string {
	out := []string{fmt.Sprintf("metrics: %d Plan/Pläne geprüft, %d mit Block", r.Plans, r.WithBlock)}
	return append(out, r.Hints...)
}

// metricsException ist ein Ausnahme-Eintrag `<N> | <Grund>`.
type metricsException struct {
	line   int
	limit  int
	reason string
}

// metricsBlock ist der gelesene Block eines Plans; Zeilen sind 1-basiert im
// Plan-Text.
type metricsBlock struct {
	startLine    int
	paths        []string
	base         string
	baseLine     int
	estimate     int
	estimateLine int
	hasEstimate  bool
	limit        int
	limitLine    int
	hasLimit     bool
	cutReason    string
	hasCutReason bool
	exceptions   []metricsException
}

// CheckMetrics ist das Regelmodul metrics (hermetischer Teil). Der
// error-Rueckgabewert traegt die Eingabefehler; der Aufrufer mappt ihn auf
// Exit 2.
func CheckMetrics(fsys driven.Filesystem, vcs driven.VCS, cfg model.MetricsConfig) (MetricsResult, error) {
	var res MetricsResult
	if cfg.PlansDir == "" || cfg.Block == "" {
		return res, fmt.Errorf("das Modul metrics braucht metrics.plans-dir und metrics.block (DC-FA-MET-001, fail-closed)")
	}
	dir := strings.TrimSuffix(cfg.PlansDir, "/")
	names, err := listMetricsPlans(fsys, dir)
	if err != nil {
		return res, err
	}
	if len(names) == 0 {
		res.Hints = append(res.Hints, fmt.Sprintf("metrics: kein Plan in %s", dir))
		return res, nil
	}
	fields := cfg.EffectiveFields()
	for _, name := range names {
		plan := dir + "/" + name
		if err := checkMetricsPlan(fsys, vcs, plan, cfg, fields, &res); err != nil {
			return MetricsResult{}, fmt.Errorf("metrics: %s: %v (DC-FA-MET-001, fail-closed)", plan, err)
		}
	}
	return res, nil
}

// listMetricsPlans liefert die Namen der Dateien slice-*.md unmittelbar in
// dir, sortiert (nicht rekursiv).
func listMetricsPlans(fsys driven.Filesystem, dir string) ([]string, error) {
	kind, err := fsys.Kind(dir)
	if err != nil {
		return nil, fmt.Errorf("metrics.plans-dir %q: %v (DC-FA-MET-001, fail-closed)", dir, err)
	}
	if kind != driven.KindDir {
		return nil, fmt.Errorf("metrics.plans-dir %q ist kein Verzeichnis (DC-FA-MET-001, fail-closed)", dir)
	}
	entries, err := fsys.List(dir)
	if err != nil {
		return nil, fmt.Errorf("metrics.plans-dir %q: %v (DC-FA-MET-001, fail-closed)", dir, err)
	}
	var names []string
	for _, e := range entries {
		ok, _ := path.Match("slice-*.md", e.Name)
		if ok && e.Kind == driven.KindSymlink {
			return nil, fmt.Errorf("metrics: %s/%s ist ein Symlink -- ein Plan ist eine Datei, kein Verweis (DC-FA-MET-001, fail-closed)", dir, e.Name)
		}
		if ok && e.Kind == driven.KindFile {
			names = append(names, e.Name)
		}
	}
	sort.Strings(names)
	return names, nil
}

// checkMetricsPlan liest und prueft einen Plan und haengt Befunde und
// Hinweise an res.
func checkMetricsPlan(fsys driven.Filesystem, vcs driven.VCS, plan string, cfg model.MetricsConfig, f model.MetricsFields, res *MetricsResult) error {
	content, err := fsys.ReadFile(plan)
	if err != nil {
		return err
	}
	res.Plans++
	blk, err := parseMetricsBlock(string(content), cfg.Block, f)
	if err != nil {
		return err
	}
	if blk == nil {
		res.Hints = append(res.Hints, fmt.Sprintf("metrics: %s: kein Block (nicht gemessen)", plan))
		return nil
	}
	res.WithBlock++
	fs, hints, err := checkMetricsBlock(plan, blk, vcs, cfg, f)
	if err != nil {
		return err
	}
	res.Findings = append(res.Findings, fs...)
	res.Hints = append(res.Hints, hints...)
	return nil
}

// fenceSpan ist ein umzaeunter Bereich: Info-String, 0-basierte Zeile des
// Oeffners, Inhalt, Index der Schluss-Zeile (oder len(lines)) und ob er
// geschlossen wurde.
type fenceSpan struct {
	info   string
	open   int
	body   []string
	end    int
	closed bool
}

// nextFence findet den naechsten Zaun ab Index from. Zaeune erkennt es mit dem
// Oeffner-Test des Produkts (nach Entfernen von Space und Tab vorn) und
// schliesst laengengleich wie der RTM-Tabellen-Leser.
func nextFence(lines []string, from int) (fenceSpan, bool) {
	for i := from; i < len(lines); i++ {
		trimmed := TrimFenceIndent(strings.TrimSuffix(lines[i], "\r"))
		if !FenceToggle(trimmed) {
			continue
		}
		char, run := FenceRun(trimmed)
		sp := fenceSpan{info: strings.TrimSpace(trimmed[run:]), open: i, end: len(lines)}
		for j := i + 1; j < len(lines); j++ {
			bl := strings.TrimSuffix(lines[j], "\r")
			if FenceCloses(TrimFenceIndent(bl), char, run) {
				sp.closed, sp.end = true, j
				break
			}
			sp.body = append(sp.body, bl)
		}
		return sp, true
	}
	return fenceSpan{}, false
}

// parseMetricsBlock liest den einen Block eines Plans. nil ohne Fehler heisst
// "kein Block". Zaeune anderer Info-Strings verbergen ihren Inhalt.
func parseMetricsBlock(content, block string, f model.MetricsFields) (*metricsBlock, error) {
	lines := strings.Split(content, "\n")
	var found *metricsBlock
	from := 0
	for {
		sp, ok := nextFence(lines, from)
		if !ok {
			return found, nil
		}
		from = sp.end + 1
		if sp.info != block {
			continue
		}
		if found != nil {
			return nil, fmt.Errorf("mehr als ein Block %q (zeile %d)", block, sp.open+1)
		}
		if !sp.closed {
			return nil, fmt.Errorf("block %q ab zeile %d ist nicht geschlossen", block, sp.open+1)
		}
		blk, err := readMetricsFields(sp.body, sp.open+2, f)
		if err != nil {
			return nil, err
		}
		blk.startLine = sp.open + 1
		found = blk
	}
}

// readMetricsFields liest die Zeilen `<schluessel>: <wert>` des Blocks;
// first ist die Plan-Zeile der ersten Block-Zeile.
func readMetricsFields(body []string, first int, f model.MetricsFields) (*metricsBlock, error) {
	b := &metricsBlock{}
	seen := map[string]bool{}
	for k, raw := range body {
		line := first + k
		s := strings.TrimSpace(raw)
		if s == "" {
			continue
		}
		key, val, ok := strings.Cut(s, ":")
		if !ok {
			return nil, fmt.Errorf("zeile %d: kein ':' in %q", line, s)
		}
		key, val = strings.TrimSpace(key), strings.TrimSpace(val)
		if seen[key] && isMetricsSingleKey(key, f) {
			return nil, fmt.Errorf("zeile %d: schlüssel %q steht mehrfach", line, key)
		}
		seen[key] = true
		if err := b.set(key, val, line, f); err != nil {
			return nil, err
		}
	}
	return b, nil
}

// isMetricsSingleKey meldet, ob der Schluessel hoechstens einmal stehen darf:
// bekannt und nicht Pfad oder Ausnahme-Eintrag.
func isMetricsSingleKey(key string, f model.MetricsFields) bool {
	switch key {
	case f.Base, f.Estimate, f.Limit, f.CutReason:
		return true
	}
	return false
}

// set ordnet einen Schluessel einem Feld zu; unbekannte Schluessel werden
// ignoriert (DC-FA-MET-001).
func (b *metricsBlock) set(key, val string, line int, f model.MetricsFields) error {
	switch key {
	case f.Path:
		return b.setPath(val, line)
	case f.Base:
		return b.setBase(val, line)
	case f.Estimate:
		n, err := metricsNumberAt(key, val, line)
		b.estimate, b.estimateLine, b.hasEstimate = n, line, err == nil
		return err
	case f.Limit:
		n, err := metricsNumberAt(key, val, line)
		b.limit, b.limitLine, b.hasLimit = n, line, err == nil
		return err
	case f.CutReason:
		b.cutReason, b.hasCutReason = val, true
	case f.Exception:
		return b.setException(val, line)
	}
	return nil
}

func (b *metricsBlock) setPath(val string, line int) error {
	if err := model.ValidLiteralPath(val); err != nil {
		return fmt.Errorf("zeile %d: %v", line, err)
	}
	b.paths = append(b.paths, val)
	return nil
}

func (b *metricsBlock) setBase(val string, line int) error {
	if !validMetricsBase(val) {
		return fmt.Errorf("zeile %d: basis %q ist kein Hash aus 7 bis 40 Zeichen 0-9a-f", line, val)
	}
	b.base, b.baseLine = val, line
	return nil
}

// setException liest `<N> | <Grund>`; ein leerer Grund ist erlaubt und ein
// Befund der Pruefung 4.
func (b *metricsBlock) setException(val string, line int) error {
	num, reason, ok := strings.Cut(val, "|")
	if !ok {
		return fmt.Errorf("zeile %d: ausnahme-eintrag %q ohne '|' (erwartet: <N> | <Grund>)", line, val)
	}
	n, err := parseMetricsNumber(strings.TrimSpace(num))
	if err != nil {
		return fmt.Errorf("zeile %d: ausnahme-eintrag: %v", line, err)
	}
	b.exceptions = append(b.exceptions, metricsException{line: line, limit: n, reason: strings.TrimSpace(reason)})
	return nil
}

func metricsNumberAt(key, val string, line int) (int, error) {
	n, err := parseMetricsNumber(val)
	if err != nil {
		return 0, fmt.Errorf("zeile %d: %s: %v", line, key, err)
	}
	return n, nil
}

// parseMetricsNumber liest ASCII-Ziffern ohne fuehrende Null und ohne
// Tausendertrennung; ein Ueberlauf von int ist ein Fehler.
func parseMetricsNumber(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("zahl fehlt")
	}
	if len(s) > 1 && s[0] == '0' {
		return 0, fmt.Errorf("zahl %q mit führender null", s)
	}
	n := 0
	for i := 0; i < len(s); i++ {
		d := int(s[i]) - '0'
		if d < 0 || d > 9 {
			return 0, fmt.Errorf("zahl %q enthält ein zeichen außer 0-9", s)
		}
		if n > (math.MaxInt-d)/10 {
			return 0, fmt.Errorf("zahl %q läuft über", s)
		}
		n = n*10 + d
	}
	return n, nil
}

// validMetricsBase prueft 7 bis 40 Zeichen 0-9a-f.
func validMetricsBase(s string) bool {
	if len(s) < 7 || len(s) > 40 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !strings.ContainsRune("0123456789abcdef", rune(s[i])) {
			return false
		}
	}
	return true
}

// checkMetricsBlock wendet die Pruefungen 1, 4 und 5 an; die Zustaende folgen
// aus den vorhandenen Feldern.
func checkMetricsBlock(plan string, b *metricsBlock, vcs driven.VCS, cfg model.MetricsConfig, f model.MetricsFields) ([]model.Finding, []string, error) {
	hint, done, err := metricsState(plan, b, f)
	if err != nil || done {
		return nil, hint, err
	}
	out, err := checkMetricsLimit(plan, b, cfg, f)
	if err != nil {
		return nil, nil, err
	}
	out = append(out, checkMetricsExceptions(plan, b, f)...)
	cut, err := checkMetricsCut(plan, b, cfg, f)
	if err != nil {
		return nil, nil, err
	}
	out = append(out, cut...)
	mf, hints, err := checkMetricsMeasure(plan, b, vcs, cfg, f)
	if err != nil {
		return nil, nil, err
	}
	return append(out, mf...), hints, nil
}

// metricsState bestimmt den Zustand: done heisst "Marke ausstehend" (nur ein
// Hinweis), ein Fehler eine ungueltige Feld-Kombination.
func metricsState(plan string, b *metricsBlock, f model.MetricsFields) ([]string, bool, error) {
	if len(b.paths) == 0 {
		return nil, false, fmt.Errorf("block ab zeile %d ohne pfad (%s)", b.startLine, f.Path)
	}
	hasBase := b.base != ""
	if !hasBase && !b.hasEstimate && !b.hasLimit {
		return []string{fmt.Sprintf("metrics: %s: Marke ausstehend", plan)}, true, nil
	}
	if !hasBase || !b.hasEstimate || !b.hasLimit {
		return nil, false, fmt.Errorf("block ab zeile %d: %s, %s und %s stehen nur zusammen oder gar nicht", b.startLine, f.Base, f.Estimate, f.Limit)
	}
	return nil, false, nil
}

// checkMetricsLimit ist Pruefung 1: die Marke entspricht der Formel.
func checkMetricsLimit(plan string, b *metricsBlock, cfg model.MetricsConfig, f model.MetricsFields) ([]model.Finding, error) {
	want, err := metricsLimit(b.estimate, cfg)
	if err != nil {
		return nil, err
	}
	if b.limit == want {
		return nil, nil
	}
	return []model.Finding{{
		File: plan, Line: b.limitLine, Rule: "metrics", Target: f.Limit,
		Reason:  ReasonMetricsLimitMismatch,
		Message: fmt.Sprintf("%s %d entspricht nicht der Formel (erwartet %d)", f.Limit, b.limit, want),
	}}, nil
}

// checkMetricsExceptions ist Pruefung 4: ein Ausnahme-Eintrag mit leerem
// Grund ist ein Befund, unabhaengig von der Messung.
func checkMetricsExceptions(plan string, b *metricsBlock, f model.MetricsFields) []model.Finding {
	var out []model.Finding
	for _, e := range b.exceptions {
		if e.reason != "" {
			continue
		}
		out = append(out, model.Finding{
			File: plan, Line: e.line, Rule: "metrics", Target: f.Exception,
			Reason:  ReasonMetricsReasonEmpty,
			Message: "Ausnahme-Eintrag mit leerem Grund",
		})
	}
	return out
}

// checkMetricsCut ist Pruefung 5: ueberschreitet Zaehler · Schaetzung die
// Schwelle, ist das Pflichtfeld zu setzen.
func checkMetricsCut(plan string, b *metricsBlock, cfg model.MetricsConfig, f model.MetricsFields) ([]model.Finding, error) {
	if !cfg.CutEnabled {
		return nil, nil
	}
	scaled, err := metricsScaled(b.estimate, cfg)
	if err != nil {
		return nil, err
	}
	if scaled <= cfg.CutThreshold || strings.TrimSpace(b.cutReason) != "" {
		return nil, nil
	}
	return []model.Finding{{
		File: plan, Line: b.estimateLine, Rule: "metrics", Target: f.CutReason,
		Reason:  ReasonMetricsCutReasonMissing,
		Message: fmt.Sprintf("%s fehlt oder ist leer, obwohl die Schwelle überschritten ist", f.CutReason),
	}}, nil
}

// metricsScaled liefert Zaehler · Schaetzung mit Ueberlauf-Pruefung.
func metricsScaled(estimate int, cfg model.MetricsConfig) (int, error) {
	if estimate != 0 && cfg.SurchargeNumerator > math.MaxInt/estimate {
		return 0, fmt.Errorf("%d · %d läuft über", cfg.SurchargeNumerator, estimate)
	}
	return cfg.SurchargeNumerator * estimate, nil
}

// metricsLimit rechnet Rounding · ⌈Zaehler · Schaetzung / Nenner⌉ ganzzahlig.
func metricsLimit(estimate int, cfg model.MetricsConfig) (int, error) {
	scaled, err := metricsScaled(estimate, cfg)
	if err != nil {
		return 0, err
	}
	den := cfg.SurchargeDenominator
	if scaled > math.MaxInt-(den-1) {
		return 0, fmt.Errorf("aufrundung von %d läuft über", scaled)
	}
	q := (scaled + den - 1) / den
	if q != 0 && cfg.Rounding > math.MaxInt/q {
		return 0, fmt.Errorf("obergrenze %d · %d läuft über", cfg.Rounding, q)
	}
	return cfg.Rounding * q, nil
}

// checkMetricsMeasure ist Pruefung 3 samt Hinweisen: die Messung aus git
// (Zeilen-Summe zwischen Basis und HEAD ueber die Pfade des Blocks) liegt auf
// oder unter der Marke, oder ein Ausnahme-Eintrag mit N >= Messung deckt sie.
// Ein Ausnahme-Eintrag mit leerem Grund deckt trotzdem (er ist selbst ein
// Befund der Pruefung 4).
func checkMetricsMeasure(plan string, b *metricsBlock, vcs driven.VCS, cfg model.MetricsConfig, f model.MetricsFields) ([]model.Finding, []string, error) {
	if vcs == nil {
		return nil, nil, fmt.Errorf("das Modul metrics braucht den git-Port (.git) fuer die Messung (DC-FA-MET-001, fail-closed)")
	}
	stat, err := vcs.LineChanges(b.base, "HEAD", b.paths, cfg.Excluded)
	if err != nil {
		return nil, nil, err
	}
	hints := metricsMeasureHints(plan, b, stat, f)
	if stat.Lines <= b.limit || metricsCovered(b, stat.Lines) {
		return nil, hints, nil
	}
	return []model.Finding{{
		File: plan, Line: b.limitLine, Rule: "metrics", Target: f.Limit,
		Reason:  ReasonMetricsOverLimit,
		Message: fmt.Sprintf("Messung %d liegt ueber %s %d und kein Ausnahme-Eintrag deckt sie", stat.Lines, f.Limit, b.limit),
	}}, hints, nil
}

// metricsCovered meldet, ob ein Ausnahme-Eintrag mit N >= Messung die
// Ueberschreitung deckt.
func metricsCovered(b *metricsBlock, measured int) bool {
	for _, e := range b.exceptions {
		if e.limit >= measured {
			return true
		}
	}
	return false
}

// metricsMeasureHints formuliert die Hinweise der Messung: Binaerdateien,
// Pfade ohne Vorkommen an Basis und HEAD, Ausnahme-Eintraege ohne
// Ueberschreitung.
func metricsMeasureHints(plan string, b *metricsBlock, stat driven.LineStat, f model.MetricsFields) []string {
	var hints []string
	if n := len(stat.Binary); n > 0 {
		hints = append(hints, fmt.Sprintf("metrics: %s: %d Binärdatei(en) nicht gezählt", plan, n))
	}
	for _, p := range stat.Missing {
		hints = append(hints, fmt.Sprintf("metrics: %s: Pfad %s existiert weder an der Basis noch an HEAD", plan, p))
	}
	if stat.Lines <= b.limit && len(b.exceptions) > 0 {
		hints = append(hints, fmt.Sprintf("metrics: %s: %s ohne Überschreitung (Messung %d, %s %d)", plan, f.Exception, stat.Lines, f.Limit, b.limit))
	}
	return hints
}
