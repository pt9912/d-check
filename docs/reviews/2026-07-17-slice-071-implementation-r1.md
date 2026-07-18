# Review — slice-071 Implementierung R1

**Datum:** 2026-07-17
**Review-Art:** unabhängiger Implementierungs-Review (kontext-getrennt; frisches
Kontextfenster, nicht der Implementierer-Kontext)
**Gegenstand:**
[`slice-071`](../plan/planning/done/slice-071-trace-cross-consistency-gate.md) —
Trace-Kreuzverweis-Konsistenz-Gate; Feat-Commit `1bfa7f8` (Range `6c4ccf5..HEAD`)
**Reviewer:** Claude (kontext-unabhängiger Lauf)
**Skill:** `.harness/skills/reviewer.md` v1.2.0
**Modell:** Opus 4.8 (1M context)

## Eingangs-Kontext

- Slice-Vertrag §2/§4/§5
  [`slice-071`](../plan/planning/done/slice-071-trace-cross-consistency-gate.md)
  (DoD §3 bewusst **nicht** erhalten — anderes Prüf-Artefakt)
- Lastenheft
  [`DC-FA-XREF-001`](../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in),
  [`DC-FA-CLI-009`](../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix),
  [`DC-FA-CLI-011`](../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code),
  [`DC-FA-REQ-001`](../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen),
  [`DC-FA-COV-001`](../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in),
  [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus),
  [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
- Spezifikation
  [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
  + Schema-Zeilen `trace.cross-consistency.*`
- [ADR-0038](../plan/adr/0038-trace-cross-consistency.md) (Proposed),
  [ADR-0005](../plan/adr/0005-modul-layout-hexagon-ordner.md),
  [ADR-0012](../plan/adr/0012-kern-paketschnitt-model-rules-app.md),
  [ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md)
- Hard Rules [`AGENTS.md`](../../AGENTS.md) §3, [`harness/conventions.md`](../../harness/conventions.md)
- Prüfgegenstand: `internal/hexagon/core/app/trace_cross.go`,
  `internal/hexagon/core/app/trace_table.go`,
  `internal/hexagon/core/app/trace.go`,
  `internal/hexagon/core/rules/matrix.go`,
  `internal/hexagon/core/model/config.go`,
  `internal/adapter/driven/configyaml/configyaml.go`,
  `internal/adapter/driven/report/report.go`,
  `internal/adapter/driving/cli/cli.go`,
  `internal/adapter/driving/cli/config_template.go` + zugehörige Tests

**Ausgeführte Sensoren:** `make lint` (0 issues), `make test` (grün),
`make arch-check` (0 Befunde), `make build` + manuelle Repro-Läufe gegen das
gebaute Image (`d-check:latest`, `--network none`, `:ro`-Mount).

## Findings

### F-1 — Namensraum-Vorbedingung ist unmechanisiert: der Abgleich läuft leer und meldet „0 Differenz(en)." (stilles Grün)

- **kategorie:** HIGH
- **quelle:** [`DC-FA-XREF-001`](../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
  (§Namensraum-Vorbedingung), [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
  Schritt 3, [ADR-0038](../plan/adr/0038-trace-cross-consistency.md)
  §Fitness-Funktion, [ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md)
  (Nullmengen-Lehre)
- **pfad:** `internal/hexagon/core/app/trace_cross.go:129-132` (`forwardView`),
  `:181-185` (`backwardView`), `:288-317` (`diffViews`)
- **befund:** Der fail-closed-Guard beider Sichten endet auf **Tabellen**-Ebene
  (`found` = „≥1 Tabelle trägt die konfigurierten Header"). Auf **Kanten**-Ebene
  gibt es keinen: Ein `design-pattern`, das zwar kompiliert, aber im falschen
  Namensraum liegt, liefert in `forwardRows` (`FindAllString`, Zeile 145) und
  `backwardRows` (`FindString` + `continue`, Zeile 192-195) für jede Zeile null
  Treffer. Beide Sichten sind dann kantenleer, `diffViews` liefert die leere
  Menge, und der Reporter druckt `0 Differenz(en).` — dieselbe Zeile, die er
  laut Kommentar (`report.go:280-283`) bewusst als **Beleg** eines gelaufenen
  Abgleichs setzt („Schweigen wäre nicht von ‚Abgleich lief nicht'
  unterscheidbar"). Die Spezifikation nennt den gemeinsamen Namensraum
  ausdrücklich eine **Vorbedingung** („sonst wäre der Mengen-Diff inhärent
  leer/voll und bedeutungslos"); im Code prüft sie nichts.
- **verifizierbar:** ja — reproduziert gegen `d-check:latest`. Repo mit
  reellem Drift (`F = {GG-AR-COMP-CORE, GG-AR-COMP-DOMAIN}` vorwärts,
  `B = {GG-AR-COMP-SCHED}` rückwärts) und `design-pattern:
  'GG-ARCH-COMP-[A-Z]+'` (Tippfehler `GG-ARCH-` statt `GG-AR-`, gültige RE2):

  ```text
  ## Kreuzverweis-Konsistenz

  0 Differenz(en).
  EXIT=0                       # --trace --require-complete
  ```

  Gegenprobe mit `design-pattern: 'GG-AR-COMP-[A-Z]+'` auf **identischem**
  Doku-Stand: 3 Differenzen, Exit 1. Der Unterschied zwischen „Gate grün" und
  „Gate blind" ist damit ein Zeichen in der Config, ohne jedes Signal.
  Ein Go-Test in `trace_cross_test.go` mit einem nicht-treffenden
  `DesignPattern` bestätigt den Befund ebenso.

### F-2 — `artifact-id-column` als Header-Name: Tabelle mit `edge-column`, aber ohne ID-Header wird still übersprungen statt Exit 2

- **kategorie:** MEDIUM
- **quelle:** [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
  Schritt 3 („zählt **jede** Tabelle mit einem `edge-column`-Header") +
  §Fehlerpräzedenz („Header-Bindung … sowie `artifact-id-column`, wenn ≠
  `first`, je genau einmal");
  [`DC-FA-XREF-001`](../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
  („fehlende Spalte … ⇒ Exit 2")
- **pfad:** `internal/hexagon/core/app/trace_cross.go:248-262`
  (`bindBackwardColumns`) i. V. m. `:221-243` (`bindCrossColumns`)
- **befund:** Bei `artifact-id-column: <Header-Name>` nimmt
  `bindBackwardColumns` den Namen in dieselbe Relevanz-Prüfung wie
  `edge-column` auf. `bindCrossColumns` liefert für ein fehlendes `names`-Element
  `(nil, false, nil)` — also „nicht relevant, kein Fehler". Eine Tabelle, die
  `Bezug` trägt, aber den konfigurierten ID-Header nicht (genau der von
  [ADR-0038](../plan/adr/0038-trace-cross-consistency.md) §Kontext beschriebene
  Realfall heterogener Header `Kennung`/`Port-ID`/`Tabu-ID`/`Komponente`), wird
  dadurch samt aller ihrer Rück-Kanten lautlos verworfen. Nach Spezifikation ist
  sie relevant und die fehlende Spalte Exit 2. Der `found`-Guard fängt nur den
  Fall, dass **gar keine** Tabelle bindet. Der Pfad hat keinen Kern-Test — die
  einzige Abdeckung von `artifact-id-column: Kennung` ist
  `configyaml/trace_cross_test.go:72-95` (Decode-Ebene).
- **verifizierbar:** ja — reproduziert gegen `d-check:latest`.
  `spec/architecture.md` mit Tabelle „4 Komponenten" (`| Kennung | Bezug |`) und
  Tabelle „5 Ports" (`| Port-ID | Bezug |`), Config `artifact-id-column: Kennung`
  + `mode: superset`. Die echte Rück-Kante `GG-AR-P-005 → GG-ARCH-006` hat keinen
  RTM-Eintrag und ist ein `B\F`-Befund; Ausgabe: `0 Differenz(en).`, Exit 0.

### F-3 — Fehlerpräzedenz: Vorwärts-Range-Expansion schlägt vor der Rückwärts-Header-Bindung zu

- **kategorie:** LOW
- **quelle:** [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
  §Fehlerpräzedenz („… → Header-Bindung … → Range-Expansion … → Diff")
- **pfad:** `internal/hexagon/core/app/trace_cross.go:83-90`
- **befund:** `crossConsistency` fährt `forwardView` **vollständig** (inkl.
  Range-Expansion in `forwardRows`) durch, bevor `backwardView` überhaupt Header
  bindet. Die Präzedenz der Spezifikation ist stufen-, nicht sicht-orientiert: ein
  Header-Bindungs-Defekt der Rück-Sicht rangiert vor einem Range-Defekt der
  Vorwärts-Sicht. Beobachtbar ist allein die Fehlermeldung (beide Fälle Exit 2) —
  dieselbe Klasse, die in slice-070 als R-Nit F-2 (Header→Duplicate-ID) gepatcht
  wurde.
- **verifizierbar:** ja — Config mit `GG-SIM-009..003` in einer Vorwärts-Zelle
  **und** doppeltem `Bezug`-Header in der Rück-Datei; gemeldet wird `AAA>BBB`,
  nach Vertrag der Mehrfach-Header.

### F-4 — `trace_cross_test.go` ist nicht gofmt-formatiert

- **kategorie:** LOW
- **quelle:** Maintainability (Präzedenz: slice-070-Review F-1, gepatcht in
  v0.43.1)
- **pfad:** `internal/hexagon/core/app/trace_cross_test.go:55-56`
- **befund:** Die Map-Literal-Einträge in `crossFS` tragen zwei Leerzeichen nach
  dem Doppelpunkt, obwohl beide Schlüssel gleich lang sind; `gofmt` richtet auf
  ein Leerzeichen aus. Kein `make`-Gate fängt das (das
  golangci-Profil führt keinen Formatter), der nächste Editor-Speichervorgang
  erzeugt daher Diff-Rauschen in einer fremden Zeile.
- **verifizierbar:** ja — `gofmt -l internal/hexagon/core/app/trace_cross_test.go`
  in der Toolchain listet die Datei; der Diff betrifft ausschließlich diese zwei
  Zeilen (die übrigen `gofmt -l`-Treffer des Repos sind Import-Gruppierungs-
  Artefakte des bloßen `gofmt` und Bestand).

### F-5 — `--json`/`--yaml` tragen keinen „Abgleich lief"-Beleg

- **kategorie:** INFO
- **quelle:** [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) i. V. m.
  [`DC-FA-CLI-009`](../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
- **pfad:** `internal/hexagon/core/app/trace.go:52`, `:66-69`
- **befund:** `CrossConsistency` trägt `omitempty`, `CrossActive` ist `json:"-"`.
  Ein aktiver Block mit null Differenzen ist in der Maschinen-Ausgabe damit nicht
  von „kein Block konfiguriert" zu unterscheiden — die Markdown-Ausgabe leistet
  über `0 Differenz(en).` genau das Gegenteil und begründet es im Kommentar. Der
  Trade-off (Byte-Identität schlägt Beleg) ist plausibel und vom Vertrag gedeckt
  (die Spezifikation verlangt kein Feld), aber nirgends notiert.
- **verifizierbar:** nein (kein Vertragsbruch — dokumentationswürdige Annahme).

## Negativbefunde (geprüft, ohne Befund)

- **Refactor `traceTableRequirements` → `markdownTables`
  (`trace_table.go:29-33`, `:66-83`, `:123-160`):** verhaltenserhaltend.
  Schleifen-Arithmetik identisch (alt: `extractTableAt` liefert `next-1`, außen
  `i = next-1` + `i++` ⇒ `next`; neu: `i = next - 1` + `i++` ⇒ `next`; Fortschritt
  garantiert, da `consumeTableRows` ab `i+2` startet). Der entfallene
  `TrimSpace(text) == ""`-Abbruch ist redundant, nicht verloren:
  `splitPipeTableLine` liefert für leere/whitespace-Zeilen `(nil, false)` und
  bricht am selben Index ab (`:313-317`). Fences: `markdownTableLines` bleibt
  unverändert und positionsgetreu. `badLine` bei **nicht**-relevanten Tabellen
  bleibt folgenlos (`extractTable` kehrt vor der `badLine`-Prüfung zurück,
  `:71-78`) — exakt das alte `relevant`-Flag. Reihenfolge `foundTable` /
  `usedTextHeaders` / `badLine` / `dupErr` unverändert; Teil-Ergebnisse vor einem
  Fehler werden in beiden Fassungen verworfen. Fehler-Reihenfolge über mehrere
  Tabellen bleibt Dokument-Reihenfolge.
- **Refactor `SelectSections` → `rules.SectionMask` (`matrix.go:337-366`):**
  äquivalent. `mask[i] = (len(include) == 0 || inRanges(inc, no)) && !inRanges(exc, no)`
  ist die wörtliche Umkehrung der alten `continue`-Kette; der
  Kurzschluss `len(include) == 0 && len(exclude) == 0 ⇒ return content` bleibt
  vor dem Masken-Bau und erhält die Identität des Rückgabe-Slices.
  Masken-Grenzen: `maskAllows` (`trace_table.go:163-165`) kann nicht
  out-of-range laufen — `SectionMask` und `markdownTableLines` splitten dieselbe
  `content` an `"\n"`, Längen sind gleich, `no` ist 1-basiert.
- **Determinismus ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)):**
  ohne Befund. `diffViews` sortiert über `(Requirement, Artifact, Direction)`;
  der Schlüssel ist eindeutig (je `crossMissing`-Aufruf ist `artifact` ein
  Map-Schlüssel, die beiden Aufrufe tragen verschiedene Richtungen), daher ist
  die Instabilität von `sort.Slice` folgenlos. `crossView.add` behält die
  **erste** Fundstelle, und `markdownTables` liefert Dokument-Reihenfolge — die
  gemeldete `Datei:Zeile` ist reproduzierbar. `requireCrossFields`
  (`configyaml.go:548-560`) sortiert die Schlüssel vor der Map-Iteration; die
  Meldung des ersten leeren Pflichtfelds ist deterministisch.
- **Byte-Identität ohne Block
  ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)):** ohne Befund.
  `crossConsistency` kehrt bei `cc == nil` mit `(nil, nil)` zurück
  (`trace_cross.go:72-74`); `CrossConsistency` steht mit `omitempty` **hinter**
  `Orphans` in der Struct-Feld-Reihenfolge, `CrossActive` ist nicht serialisiert,
  `traceCross` (`report.go:278-280`) kehrt ohne `CrossActive` sofort zurück.
  Markdown-Pfad ist per `TestCLI071_Cross_DefaultByteIdentisch` belegt; der
  `--json`/`--yaml`-Pfad ist nicht getestet, aber durch `omitempty` auf einem
  nil-Slice garantiert (auch der aktive 0-Differenzen-Fall liefert ein
  Slice der Länge 0 und entfällt).
- **Modi `equal`/`superset` (`trace_cross.go:301-304`):** vertragstreu — `equal`
  meldet beide Richtungen, `superset` unterdrückt allein `F\B`. Per CLI-Test und
  eigener Repro bestätigt.
- **`exclude-req`-Ventil (`trace_cross.go:298-300`):** vertragstreu — die
  Kennungen fallen aus der **Schlüsselmenge** der Vereinigung, vor dem Diff, aus
  beiden Sichten (Spezifikation Schritt 4).
- **Geteilte `design-pattern`-Vorbedingung als Mechanik
  (`trace_cross.go:87`, `:158`, `:192`):** ohne Befund — die Rückwärts-Artefakt-ID
  wird tatsächlich über `cc.Forward.DesignPattern` extrahiert, nicht über ein
  zweites Muster; ein zweiter Namensraum ist damit nicht ausdrückbar. (Dass das
  eine **gemeinsame** Treffermenge nicht garantiert, ist F-1.)
- **`artifact-id-column: first` (`trace_cross.go:250-261`):** vertragstreu —
  Sentinel ⇒ `idIdx = 0` (erste Spalte), Artefakt-ID = **erster**
  `design-pattern`-Treffer der Zelle (`FindString`), Zeilen ohne Treffer werden
  übersprungen. Default-Auflösung in `applyCrossBackward` korrekt.
- **Sortierung `(R, Artefakt, Richtung)` und Befund-Inhalt:** vertragstreu —
  `CrossFinding` trägt Anforderungs-ID, Artefakt, Richtungslabel und
  `Datei:Zeile` der nennenden Sicht; per CLI-Test und Repro geprüft.
- **Fail-closed-Guards, die greifen:** fehlende `forward.file`/`backward.file`
  (beide **vor** jeder Sicht-Auswertung gelesen, `trace_cross.go:75-82` — die
  Präzedenz „Quellen lesen vor Header-Bindung" stimmt), keine Tabelle mit den
  konfigurierten Headern, mehrfacher Rollen-Header, Sektionsname ohne
  Heading-Treffer (`checkCrossSections`, konsistent zum
  [`DC-FA-COV-001`](../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)-Pendant
  `checkSectionNames`), Zellenzahl-Bruch in relevanter Tabelle, fehlender
  Pflichtblock, unbekannter `mode`, nicht kompilierendes Regex, leeres
  Pflichtfeld, Pfad außerhalb der Repo-Wurzel, unbekannter YAML-Schlüssel
  (`KnownFields(true)`). Jeder mit Negativtest belegt.
- **Gate-Bindung ([`DC-FA-CLI-011`](../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)):**
  ohne Befund — `requireCompleteExit` (`cli.go:203-234`) gatet ausschließlich
  unter `--require-complete`, ohne block-lokalen Schalter; Waisen- und
  Kreuzverweis-Ursache melden **getrennt** auf stderr und verdecken einander
  nicht (`exit = 1` statt frühem `return`). `--trace` allein bleibt Exit 0.
- **Erfundener Vertrag:** ohne Befund — keine Validierung gefunden, die eine
  nach Schema legitime Config ablehnt. Der `found`-Guard („keine Tabelle mit den
  konfigurierten Headern ⇒ Exit 2") steht nicht wörtlich in der Spezifikation,
  ist aber durch die Fail-closed-Klausel des Slice-Vertrags §2 und die
  [ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md)-Lehre
  gedeckt. `requireCrossFields` prüft `TrimSpace` — reine Leer-Prüfung, das
  Muster selbst wird ungetrimmt kompiliert (keine stille Mutation).
- **Hexagon-Import-Regeln ([ADR-0005](../plan/adr/0005-modul-layout-hexagon-ordner.md),
  [ADR-0012](../plan/adr/0012-kern-paketschnitt-model-rules-app.md)):** ohne
  Befund — `trace_cross.go` (Schicht `app`) importiert `model`, `rules`,
  `port/driven`; `SectionMask` liegt korrekt in `rules`, nicht in `app`. Keine
  Rückwärts-Kante. `make arch-check` (a-check-Image, digest-gepinnt): 0 Befunde.
- **[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit):**
  ohne Befund — reiner Lese-Pfad über den Filesystem-Port (`pathExists`,
  `ReadFile`), kein Netz, kein Schreibzugriff. Repro-Läufe liefen mit
  `--network none` und `:ro`-Mount durch.
- **[`AGENTS.md`](../../AGENTS.md) §3.2 (Suppression-Verbot):** ohne Befund —
  keine `//nolint`-Direktive im Diff, keine neue `.golangci.yml`-Ausnahme.
- **[`AGENTS.md`](../../AGENTS.md) §3.6 (Gate-Lockerung nur per ADR):** ohne
  Befund — keine Schwelle gesenkt; die Coverage-Bindung bleibt unverändert.
- **`--print-config`-Vorlage (`config_template.go:180-197`):** ohne Befund — der
  Block ist auskommentiert (keine Verhaltensänderung), die Schlüssel decken sich
  mit dem Schema, `design-pattern` ist als geteilt annotiert.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 1 | F-1 |
| MEDIUM | 1 | F-2 |
| LOW | 2 | F-3, F-4 |
| INFO | 1 | F-5 |

## Verdikt

**REJECT.**

F-1 blockiert: Der Slice existiert, um einen Drift zu fangen, den „kein Gate
bemerkt" (Slice §5). In der vorliegenden Fassung entscheidet ein einziges
Config-Zeichen darüber, ob das Gate prüft oder blind ist — und der blinde Fall
druckt `0 Differenz(en).` und Exit 0, also genau die Zeile, die der Reporter als
Beleg eines gelaufenen Abgleichs setzt. Das ist ein Stilles-Grün-Pfad im
Gate-Pfad (Harness-Lüge) nach dem HIGH-Anker des Reviewer-Skills; die
Namensraum-Kongruenz ist von Lastenheft, Spezifikation (Schritt 3) und
[ADR-0038](../plan/adr/0038-trace-cross-consistency.md) als **Vorbedingung**
benannt, aber nirgends mechanisiert. Bemerkenswert: der Slice hat den Riss selbst
als Risiko §4 „Namensraum-Kongruenz" notiert — er ist von der Doku in den Code
nicht mitgewandert.

F-2 blockiert nach der Regel „HIGH und MEDIUM blockieren typischerweise"
ebenfalls: eine relevante Rück-Tabelle wird entgegen der ausdrücklichen
Fehlerpräzedenz still verworfen. Die Wirkung ist konfigurations-bedingt
(`artifact-id-column` ≠ `first`) und trifft die dokumentierte Default-Nutzung
(`first`) nicht — die Einordnung als MEDIUM statt HIGH ist bewusst; unter
`mode: superset` ist die Wirkung allerdings ebenfalls ein stilles Grün.

F-3/F-4 sind Nits ohne Blockade-Anspruch; F-5 ist eine Notiz.

Nicht beanstandet und ausdrücklich bestätigt: der riskanteste Teil des Diffs —
der Umbau auf den gemeinsamen `markdownTables`-Walker und `rules.SectionMask` —
ist nach Zeilen-Analyse aller benannten Randfälle (Fences, Tabellen-Ende,
Zellenzahl-Bruch bei nicht-relevanten Tabellen, Masken-Grenzen, leere Zeilen,
`i = next - 1`) verhaltenserhaltend, und die Byte-Identität der RTM ohne Block
hält in allen drei Ausgabeformaten.
