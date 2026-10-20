package app

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// Eine Datei, zwei Anforderungsformen (DC-FA-REQ-001): Überschriften für FA,
// eine Tabelle für QA.
const sourcesLastenheft = `# Lastenheft

#### SFC-FA-01 — Erste
Das System MUSS laufen.

#### SFC-FA-02 — Zweite
Das System KANN mehr.

| ID | Aspekt | Anforderung | Priorität |
|---|---|---|---|
| SFC-QA-01 | Perf | Antwort in einer Sekunde. | MUSS |
| SFC-QA-02 | Ops | Logs lesbar. | KANN |
`

var (
	sfcFA     = regexp.MustCompile(`SFC-FA-\d{2}`)
	sfcQA     = regexp.MustCompile(`SFC-QA-\d{2}`)
	sfcFAorQA = regexp.MustCompile(`SFC-(?:FA|QA)-\d{2}`)
)

func sfcTable(idColumn string) *model.TraceTableConfig {
	return &model.TraceTableConfig{
		IDColumn: idColumn, TextColumns: []string{"Anforderung"},
		ModalityColumn: "Priorität", DuplicateIDs: model.TraceDuplicateError,
	}
}

func sfcSources() []model.TraceSource {
	return []model.TraceSource{
		{Source: "spec/lastenheft.md", ReqPattern: sfcFA, Format: model.TraceFormatHeadings},
		{Source: "spec/lastenheft.md", ReqPattern: sfcQA, Format: model.TraceFormatTable, Table: sfcTable("ID")},
	}
}

func rowIDs(m TraceMatrix) []string {
	ids := make([]string, len(m.Requirements))
	for i, r := range m.Requirements {
		ids[i] = r.ID
	}
	return ids
}

// Die Vereinigung führt erst die Überschriften-, dann die Tabellen-
// Anforderungen, in Listenreihenfolge; die Summe zählt über beide Einträge.
func TestTraceSourcesUnionInListOrder(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{"spec/lastenheft.md": sourcesLastenheft})
	m, err := BuildTraceMatrix(fs, model.TraceConfig{Sources: sfcSources()})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	want := []string{"SFC-FA-01", "SFC-FA-02", "SFC-QA-01", "SFC-QA-02"}
	if got := rowIDs(m); !reflect.DeepEqual(got, want) {
		t.Fatalf("Reihenfolge = %v, want %v", got, want)
	}
	if m.Total != 4 {
		t.Fatalf("Total = %d, want 4", m.Total)
	}
	// Die Listenreihenfolge bestimmt die Ausgabe, nicht die Kennungs-Sortierung.
	rev := sfcSources()
	rev[0], rev[1] = rev[1], rev[0]
	m, err = BuildTraceMatrix(fs, model.TraceConfig{Sources: rev})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	want = []string{"SFC-QA-01", "SFC-QA-02", "SFC-FA-01", "SFC-FA-02"}
	if got := rowIDs(m); !reflect.DeepEqual(got, want) {
		t.Fatalf("umgekehrte Reihenfolge = %v, want %v", got, want)
	}
}

// Dieselbe Kennung als Überschrift und als Tabellenzeile ist Exit 2 mit
// Kennung und beiden Einträgen — in beiden Reihenfolgen der Einträge, denn die
// Prüfung hat zwei Seiten.
func TestTraceSourcesDuplicateAcrossEntries(t *testing.T) {
	doc := sourcesLastenheft + "\n#### SFC-QA-01 — Auch als Überschrift\nText.\n"
	fs := coretest.NewMemFS(map[string]string{"spec/lastenheft.md": doc})
	both := func() []model.TraceSource {
		s := sfcSources()
		s[0].ReqPattern, s[1].ReqPattern = sfcFAorQA, sfcFAorQA
		return s
	}
	for name, src := range map[string][]model.TraceSource{
		"überschrift zuerst": both(),
		"tabelle zuerst":     {both()[1], both()[0]},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := BuildTraceMatrix(fs, model.TraceConfig{Sources: src})
			if err == nil {
				t.Fatal("erwartete Fehler für die doppelte Kennung über Einträge")
			}
			for _, want := range []string{"SFC-QA-01", "Eintrag 0", "Eintrag 1", "überlappende id-pattern"} {
				if !strings.Contains(err.Error(), want) {
					t.Fatalf("Meldung %q trägt %q nicht", err, want)
				}
			}
		})
	}
}

// Ein Fehler im zweiten Eintrag nennt dessen Index — fehlender
// Header wie die Nullmenge.
func TestTraceSourcesErrorsNameTheEntry(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{"spec/lastenheft.md": sourcesLastenheft})
	missingHeader := sfcSources()
	missingHeader[1].Table = sfcTable("Kennung")
	zero := sfcSources()
	zero[1] = model.TraceSource{Source: "spec/lastenheft.md", ReqPattern: regexp.MustCompile(`SFC-RB-\d{2}`)}
	missingSource := sfcSources()
	missingSource[1].Source = "spec/fehlt.md"
	for name, c := range map[string]struct {
		src  []model.TraceSource
		want string
	}{
		"header fehlt":   {missingHeader, "trace.requirements[1]"},
		"nullmenge":      {zero, "trace.requirements[1]: Quelle"},
		"quelle fehlt":   {missingSource, "trace.requirements[1]: Quelle"},
		"nullmenge-text": {zero, "ergab 0 Anforderungen"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := BuildTraceMatrix(fs, model.TraceConfig{Sources: c.src})
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("Fehler = %v, want enthält %q", err, c.want)
			}
			if strings.Contains(err.Error(), "trace.requirements[0]") {
				t.Fatalf("Meldung nennt den falschen Eintrag: %v", err)
			}
		})
	}
}

// Die Objekt-Form bleibt ohne Index in den Meldungen (DC-QA-02).
func TestTraceObjectFormErrorsCarryNoIndex(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{"spec/lastenheft.md": sourcesLastenheft})
	_, err := BuildTraceMatrix(fs, model.TraceConfig{
		Source: "spec/lastenheft.md", ReqPattern: regexp.MustCompile(`SFC-RB-\d{2}`),
	})
	if err == nil || strings.Contains(err.Error(), "[") {
		t.Fatalf("Objekt-Form trägt Index oder keinen Fehler: %v", err)
	}
}

// Eine Liste mit einem Eintrag liefert dieselbe RTM wie die Objekt-Form.
func TestTraceSourcesSingleEntryEqualsObjectForm(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{"spec/lastenheft.md": sourcesLastenheft})
	obj, err := BuildTraceMatrix(fs, model.TraceConfig{
		Source: "spec/lastenheft.md", ReqPattern: sfcFA, Modality: &model.TraceModality{},
	})
	if err != nil {
		t.Fatalf("Objekt-Form: %v", err)
	}
	list, err := BuildTraceMatrix(fs, model.TraceConfig{Sources: []model.TraceSource{
		{Source: "spec/lastenheft.md", ReqPattern: sfcFA, Modality: &model.TraceModality{}},
	}})
	if err != nil {
		t.Fatalf("Liste: %v", err)
	}
	if !reflect.DeepEqual(obj, list) {
		t.Fatalf("Liste mit einem Eintrag weicht ab:\n obj=%+v\nlist=%+v", obj, list)
	}
}

// Die Referenzseite erkennt die Muster aller Einträge: ein Slice, der eine
// QA-Kennung nennt, deckt sie, obwohl das Muster des ersten Eintrags sie nicht
// trifft.
func TestTraceSourcesReferenceSideSeesAllPatterns(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{
		"spec/lastenheft.md":                     sourcesLastenheft,
		"docs/plan/planning/slice-001-x.md":      "Liefert SFC-QA-02 und SFC-FA-01.\n",
		"docs/plan/planning/slice-002-ohne-id.md": "Keine Kennung.\n",
	})
	m, err := BuildTraceMatrix(fs, model.TraceConfig{Sources: sfcSources()})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	covered := map[string]bool{}
	for _, r := range m.Requirements {
		covered[r.ID] = !r.Orphan
	}
	want := map[string]bool{"SFC-FA-01": true, "SFC-FA-02": false, "SFC-QA-01": false, "SFC-QA-02": true}
	if !reflect.DeepEqual(covered, want) {
		t.Fatalf("Deckung = %v, want %v", covered, want)
	}
}

// Die Vereinigung nimmt den längsten Treffer: ein kurzes Muster vor einem
// langen verschluckt die Kennung nicht (Ganz-Token auf mindestens ein Muster).
func TestReferencePatternLongestMatch(t *testing.T) {
	entries := []reqEntry{
		{reqPat: regexp.MustCompile(`X-\d`)},
		{reqPat: regexp.MustCompile(`X-\d{3}`)},
	}
	re, err := referencePattern(entries)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if got := re.FindString("siehe X-123."); got != "X-123" {
		t.Fatalf("FindString = %q, want X-123", got)
	}
	single, err := referencePattern(entries[:1])
	if err != nil || single != entries[0].reqPat {
		t.Fatalf("ein Eintrag muss sein Muster unverändert liefern: %v %v", single, err)
	}
}

// modality gilt je Eintrag: der Eintrag ohne modality klassifiziert nicht, seine
// Waisen gaten alle; der Eintrag mit modality gatet nur die Stufen aus seinen
// require-levels.
func TestTraceSourcesModalityPerEntry(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{"spec/lastenheft.md": sourcesLastenheft})
	src := sfcSources()
	src[0].Modality = &model.TraceModality{} // FA: Body-Klassifikation, gatet [must]
	m, err := BuildTraceMatrix(fs, model.TraceConfig{Sources: src})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if !m.ModalityActive {
		t.Fatal("ModalityActive fehlt, obwohl ein Eintrag modality trägt")
	}
	levels := map[string]string{}
	for _, r := range m.Requirements {
		levels[r.ID] = r.Modality
	}
	want := map[string]string{"SFC-FA-01": "must", "SFC-FA-02": "may", "SFC-QA-01": "", "SFC-QA-02": ""}
	if !reflect.DeepEqual(levels, want) {
		t.Fatalf("Stufen = %v, want %v", levels, want)
	}
	// 4 Waisen; gatend: FA-01 (must) und beide QA (Eintrag ohne modality).
	if m.Orphans != 4 || m.GatingOrphans != 3 {
		t.Fatalf("Orphans/Gating = %d/%d, want 4/3", m.Orphans, m.GatingOrphans)
	}
	// Ohne modality in irgendeinem Eintrag gaten alle (byte-identisch).
	m, err = BuildTraceMatrix(fs, model.TraceConfig{Sources: sfcSources()})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if m.ModalityActive || m.GatingOrphans != 4 {
		t.Fatalf("ohne modality: active=%v gating=%d, want false/4", m.ModalityActive, m.GatingOrphans)
	}
}

// Die Range-Validierung der Coverage prüft die expandierten Kennungen gegen die
// Vereinigung der Muster: eine Range über die QA-Reihe zählt, obwohl das Muster
// des ersten Eintrags sie nicht trifft.
func TestTraceSourcesCoverageRangesUseAllPatterns(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{
		"spec/lastenheft.md":     sourcesLastenheft,
		"docs/traceability.md": "Abdeckung: SFC-QA-01..02 und SFC-FA-01..02.\n",
	})
	m, err := BuildTraceMatrix(fs, model.TraceConfig{
		Sources:  sfcSources(),
		Coverage: []model.TraceCoverage{{Files: []string{"docs/traceability.md"}, Label: "Trace", Ranges: true}},
	})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	for _, r := range m.Requirements {
		if r.Orphan {
			t.Fatalf("%s bleibt Waise trotz Range in der Coverage-Quelle: %+v", r.ID, m.Requirements)
		}
	}
}

// cross-consistency liest das Muster der Vorwärts-Sicht aus der Vereinigung:
// Kennungen beider Einträge stehen in beiden Sichten, die Differenz ist leer.
// Mit dem Muster nur des ersten Eintrags fiele die QA-Kennung vorwärts heraus und
// erschiene rückwärts als Befund.
func TestTraceSourcesCrossConsistencyUsesAllPatterns(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{
		"spec/lastenheft.md":   sourcesLastenheft,
		"docs/traceability.md": fwdDoc("| SFC-FA-01 | GG-AR-A |\n| SFC-QA-01 | GG-AR-B |\n"),
		"spec/architecture.md": bwdDoc("| GG-AR-A | SFC-FA-01 |\n| GG-AR-B | SFC-QA-01 |\n"),
	})
	cc := crossCfg()
	cc.Backward.ReqPattern = sfcFAorQA
	m, err := BuildTraceMatrix(fs, model.TraceConfig{Sources: sfcSources(), CrossConsistency: cc})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if len(m.CrossConsistency) != 0 {
		t.Fatalf("Differenz trotz konsistenter Sichten: %+v", m.CrossConsistency)
	}
}

// Ein Eintrag mit `format: headings` und ohne explizite source ist nicht strikt:
// null Treffer ergeben dort keinen Fehler (ADR-0108, Grenze) — im Gegensatz zum
// Eintrag mit expliziter source.
func TestTraceSourcesHeadingsEntryWithoutSourceStaysLenient(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{"spec/lastenheft.md": "# Lastenheft\n\nKeine Anforderung.\n"})
	m, err := BuildTraceMatrix(fs, model.TraceConfig{Sources: []model.TraceSource{{ReqPattern: sfcFA}}})
	if err != nil || m.Total != 0 {
		t.Fatalf("Default-Quelle ohne Treffer: Total=%d err=%v, want 0/nil", m.Total, err)
	}
	_, err = BuildTraceMatrix(fs, model.TraceConfig{Sources: []model.TraceSource{{Source: "spec/lastenheft.md", ReqPattern: sfcFA}}})
	if err == nil || !strings.Contains(err.Error(), "trace.requirements[0]") {
		t.Fatalf("explizite Quelle ohne Treffer: Fehler = %v, want Eintrag 0", err)
	}
}

// Ein Fehler eines späteren Eintrags geht dem Duplikat über Einträge voraus: die
// Duplikat-Prüfung beginnt erst, wenn alle Einträge extrahiert sind.
func TestTraceSourcesLaterEntryErrorPrecedesDuplicate(t *testing.T) {
	doc := sourcesLastenheft + "\n#### SFC-QA-01 — Auch als Überschrift\nText.\n"
	fs := coretest.NewMemFS(map[string]string{"spec/lastenheft.md": doc})
	src := sfcSources()
	src[0].ReqPattern, src[1].ReqPattern = sfcFAorQA, sfcFAorQA
	src = append(src, model.TraceSource{Source: "spec/fehlt.md", ReqPattern: sfcFA})
	_, err := BuildTraceMatrix(fs, model.TraceConfig{Sources: src})
	if err == nil || !strings.Contains(err.Error(), "trace.requirements[2]: Quelle") {
		t.Fatalf("Fehler = %v, want der fehlenden Quelle von Eintrag 2", err)
	}
}

// Jedes Einzelmuster kompiliert, die Vereinigung scheitert aber an der
// Schachtelungstiefe des Regex-Parsers: ein Fehler der Konfiguration, kein
// stilles Weglassen.
func TestReferencePatternUnionFailure(t *testing.T) {
	deep := strings.Repeat("(", 999) + "X-1" + strings.Repeat(")", 999)
	entries := []reqEntry{{reqPat: regexp.MustCompile(deep)}, {reqPat: sfcFA}}
	if _, err := referencePattern(entries); err == nil || !strings.Contains(err.Error(), "vereinigen") {
		t.Fatalf("erwartete Vereinigungs-Fehler, got %v", err)
	}
}
