# Verifikation — slice-274: `trace.requirements` als Liste von Quellen

- **Rolle:** Verifier (Modul 11) — DoD, Plan und Spec gegen den Stand, nicht der Diff gegen Plan/ADR
- **Gegenstand:** `docs/plan/planning/done/wellenlos/slice-274-trace-requirements-als-liste-von-quellen.md` §1, §2, §3 (samt Plan-Änderungsnotizen), §4, §6
- **Range:** `68cae247..HEAD` (10 Commits, Spitze `17c7ce7c`)
- **Eingang:**
  [`DC-FA-REQ-001`](../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen)
  (Absatz „Mehrere Anforderungsquellen“, Akzeptanzkriterien Liste);
  [§`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen) Schritt 8;
  [ADR-0108](../plan/adr/0108-trace-requirements-als-liste-von-quellen.md) (Proposed);
  Reviews Runde 1 und 2 nur als Kontext
- **Datum:** 2026-10-20

## Messungen (Kommando und Ergebnis)

Alle Läufe selbst gefahren, nur über make/Docker und Lesewerkzeuge. Der Vorzustand
(`68cae247`) lief als eigenes Image aus einem Wegwerf-Worktree im Scratchpad, die
Brech-Proben in einem zweiten Worktree auf `17c7ce7c`; nach jeder Probe
`git checkout -- .`. Worktree, Probe-Images und Fixture-Verzeichnis sind entfernt,
`git status` des Repos ist leer.

| # | Kommando | Ergebnis |
|---|---|---|
| M1 | `make gates` auf `17c7ce7c` | Exit 0; `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; lint `0 issues.`; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; semgrep `0 findings`; doc-check `1167 Datei(en) geprüft, 0 Befund(e)` |
| M2 | Objekt-Form, Vorher-Image (`68cae247`) gegen Nachher-Image (`HEAD`), stdout, stderr und Exit getrennt verglichen: das Repo selbst (`--trace`, `--trace --require-complete`, `--trace --json`, `--json`) · Heading-Fixture (`--trace`, `--require-complete` Exit 1, `--json`, `--doctor`) · Tabellen-Fixture mit `modality` (dieselben) | alle `IDENT` (`cmp` auf stdout, stderr, Exit); das Repo selbst 12548 B Matrix, Exit 0, `55 Anforderung(en), 0 Waise(n)` |
| M3 | Fehlermeldungen kaputter Objekt-Form-Konfigurationen, beide Images, je mit `--trace`, `--trace --require-complete`, `--trace --json`, `--doctor`, `--json`: unbekannter Schlüssel (`sourse`), falsches `format` (`csv`), falsche Form (`requirements: "text"`), Muster ohne Treffer, fehlende Quelle, `format: table` ohne `text-column` | alle `IDENT`, Meldungen und Exit 2 wortgleich, z. B. `.d-check.yml: yaml: unmarshal errors: line 3: field sourse not found` |
| M4 | Listen-Probe nach dem Lastenheft-Beispiel: eine Datei, 21 Überschriften (`LH-FA-nn`) + 6 Tabellenzeilen (`LH-QA-nnn`), Liste aus zwei Einträgen, `--trace` | Exit 0; 27 Zeilen, erst die 21 Überschriften, dann die 6 Tabellenzeilen; `27 Anforderung(en), 24 Waise(n).` |
| M5 | wie M4, `--trace --require-complete`, zwei Slices belegen zwei Kennungen | Exit 1, `24 Requirements-Waise(n) ohne referenzierenden Slice (--require-complete)`; mit einem Slice, der alle 27 Kennungen nennt: Exit 0, `0 Waise(n)` |
| M6 | wie M4, Fehlerproben | doppelte Kennung (beide Muster `LH-(?:FA\|QA)-\d+`, `LH-QA-003` auch als Überschrift): Exit 2, `Kennung "LH-QA-003" kommt in Eintrag 0 und in Eintrag 1 vor (doppelt definiert oder überlappende id-pattern)` · falscher Header in Eintrag 1: Exit 2, `trace.requirements[1]: keine Tabelle mit den konfigurierten Headern "Nope" …` · Eintrag ohne Treffer (Fixture-Fehler, versehentlich): Exit 2, `trace.requirements[1]: Quelle … im Format table ergab 0 Anforderungen` · `requirements: []`: Exit 2, `line 2: trace.requirements darf als Liste nicht leer sein` · unbekannter Schlüssel im Eintrag: Exit 2, `line 4: field bogus not found` (Zeile stimmt) · `- null`: Exit 2, `trace.requirements[0] ist leer` · `- "x"`: Exit 2, `cannot unmarshal !!str …` |
| M7 | Bewusstes Brechen A–G, je `make test` im Wegwerf-Worktree (Tabelle unten) | jede Probe rot, jeweils mit dem benannten Test und der behaupteten Meldung |

### Bewusstes Brechen (M7)

| Probe | Fix testweise zurückgenommen | Rot durch |
|---|---|---|
| A | Kreuz-Duplikat-Prüfung abgeschaltet (`dup && false`) | `trace_sources_test.go:103: erwartete Fehler für die doppelte Kennung über Einträge`; `cli_trace_sources_test.go:80` (Exit 0 statt 2, `28 Anforderung(en)`, `SFC-QA-03` doppelt in der Matrix) |
| B | Nullmengen-Guard je Eintrag abgeschaltet | `trace_sources_test.go:136: Fehler = <nil>, want enthält "ergab 0 Anforderungen"`; `:152`; `cli_trace_sources_test.go:107` (Exit 0, 21 Waisen statt Exit 2) |
| C | Leere-Liste-Prüfung abgeschaltet | `trace_requirements_list_test.go:85: Decode-Fehler = <nil>, want enthält "darf als Liste nicht leer sein"`; `cli_trace_sources_test.go:144` (`0 Anforderung(en)`, Exit 0) |
| D | `modality` global statt je Eintrag (`mods[0]`) | `trace_sources_test.go:239: Stufen = map[… SFC-QA-01:unknown …], want … SFC-QA-01: …` |
| E | Referenzseite liest nur das Muster des ersten Eintrags | `trace_sources_test.go:195` (Deckung), `:211: FindString = "X-1", want X-123`, `:272` (Coverage-Range), `cli_trace_sources_test.go:129` |
| F | Strikt-Prüfung am Knotenbaum entfernt | `trace_requirements_list_test.go:85: … want enthält "field sorce not found"`, `:97: want line 7` |
| G | Index aus den Meldungen entfernt | `cli_trace_sources_test.go:107` und `trace_sources_test.go:136: … want enthält "trace.requirements[1]: Quelle"` |

## DoD-Punkte

**DoD 1 — Liste laden (AK-3, AK-7).** **Erfüllt.**
`rawTraceRequirementsForm.UnmarshalYAML` nimmt Objekt oder nichtleere Sequenz an;
`applyTraceRequirementsForm` validiert je Eintrag und setzt den Index in die
Meldung; das Modell trägt `TraceConfig.Sources []TraceSource`. Fremde Form, leere
Liste, `null`-Eintrag und unbekannter Schlüssel enden mit Exit 2 (M6), die Zeile
stammt aus dem Dokument. Rot-Beleg: Probe C (leere Liste), F (unbekannter
Schlüssel), G (Index).

**DoD 2 — Vereinigung führen (AK-1, AK-2, AK-4, AK-6; AK-5).** **Erfüllt.**
`loadTraceRequirements` extrahiert alle Einträge, prüft erst danach die Kennung
über Einträge (Fehlerpräzedenz der Spezifikation), führt die Vereinigung in
Listenreihenfolge, `modality` je Eintrag, und die Referenzseite nutzt die
Mustervereinigung. Belege: M4 (21 + 6 = 27, Überschriften zuerst), M5 (Exit 1 bei
Waise, Exit 0 bei voller Abdeckung), M6 (Kreuz-Duplikat mit Kennung und beiden
Einträgen, Index bei Header-/Nullmengen-Fehler). Rot-Beleg: Proben A, B, D, E, G,
jeweils aus dem richtigen Grund. Objekt-Form byte-identisch: M2, M3.

**DoD 3 — Handbuch und ADR.** **Handbuch erfüllt, ADR-Zeile offen.**
Das Handbuch beschreibt die Listen-Form (§5 *Mehrere Quellen*, Beispiel, Vereinigung,
doppelte Kennung, leere Liste, `modality` je Eintrag, Referenzseite), den
Nullmengen-Guard je Eintrag (§5 und §7) und das Konfigurationsbeispiel. Jede
Aussage deckt sich mit den Proben M4–M6. `docs/user/operations.md` (Zeile
`--trace`) ist nachgezogen. **Offen, wie vor der Closure zu erwarten:** ADR-0108
steht noch auf `Proposed` (Kopf und Index). Das ist keine Verletzung. Siehe V-1
zum Inhalt der ADR vor dem Statuswechsel.

**DoD 4 — `make gates` grün.** **Erfüllt** (M1).

**DoD 5–9 — Review, Register, Closure-Notiz, Risiko-Ausgänge, Paarungen.** **Offen**,
wie es vor der Closure sein muss: §7 ist noch Platzhalter, die drei Risiken aus §6
stehen auf *offen bis Closure*. Mit diesem Bericht ist die Verifikation geliefert;
Review-Reports Runde 1 und 2 liegen vor.

## Akzeptanzkriterien gegen Tests

| Akzeptanzkriterium des Lastenhefts | Test | Rot ohne Fix |
|---|---|---|
| Happy Path (Liste) | `TestTraceSourcesUnionInListOrder`, `TestCLI_TraceListe_VereinigungInListenreihenfolge` | E (Referenzmuster); Reihenfolge durch die Umkehr-Probe im Test gepinnt |
| Negative (doppelte Kennung) | `TestTraceSourcesDuplicateAcrossEntries`, `TestCLI_TraceListe_DoppelteKennungExit2`, `TestTraceSourcesLaterEntryErrorPrecedesDuplicate` | A |
| Negative (Eintrag) | `TestTraceSourcesErrorsNameTheEntry`, `TestCLI_TraceListe_FehlerNenntDenEintrag` | B, G |
| Negative (leere Liste) | `TestDecode_RequirementsListErrors`, `TestCLI_TraceListe_LeereListeExit2` | C |
| Boundary (Vollständigkeit) | `TestCLI_TraceListe_RequireCompleteWertetVereinigungAus` | E |
| Objekt-Form byte-identisch | `TestDecode_RequirementsObjectFormUnchanged`, `TestTraceObjectFormErrorsCarryNoIndex`, `TestTraceSourcesSingleEntryEqualsObjectForm` | Kein Mutationsnachweis möglich, weil die Objekt-Form heute noch den Vorzustand ist; ersatzweise M2/M3 gegen das Vorher-Image |
| `modality` je Eintrag (Absatz Mehrere Quellen) | `TestTraceSourcesModalityPerEntry` | D |

## Plan gegen Code

*Was §1 ausschließt, wurde nicht geliefert:*
- **Lastenheft, Spezifikation, ADR-Text:** im Range unverändert (`git diff --stat` auf `spec/`, `docs/plan/adr/`).
- **`CHANGELOG.md`, READMEs, Handbuch-Kopf, Versionstabelle:** unberührt.
- **`--print-config`/`--suggest-config`:** unverändert, erzeugen weiter keine Liste.
- **`first`/`last` über Einträge, Vererbung von `modality`:** nicht vorhanden; Probe D belegt die Trennung.
- **Strikter Modus für Listen-Einträge ohne `source`:** nicht eingeführt; `TestTraceSourcesHeadingsEntryWithoutSourceStaysLenient` pinnt die Parität.

*Jede Datei im Diff ist durch §3 oder eine Plan-Änderungsnotiz gedeckt:*
`model/config.go`, `configyaml.go`, `trace.go` (Tabelle); `trace_sources.go` und die
drei neuen Testdateien (Plan-Änderung vor dem Code); `yamlstrict.go` (Runde 1, F-1/F-3);
`abdeckung-tests.md` (generiert, notiert); `benutzerhandbuch.md`, `operations.md`
(Tabelle und Spiegel-Liste). `trace_cross.go` und `trace_table.go` sind laut Notiz
unverändert und sind es im Diff auch. Die Tabellenzeile `cli_acceptance_test.go` ist
durch `cli_trace_sources_test.go` ersetzt, ebenfalls notiert.

## Befunde

- **V-1 — MEDIUM — ADR-0108 widerspricht dem Slice-Plan und dem Stand.** Der Abschnitt
  *Konsequenzen* der ADR sagt: „`--print-config` und das Benutzerhandbuch zeigen die
  Listen-Form; das Handbuch wird in der Release-Prep nachgezogen.“ Der Slice schließt
  `--print-config` ausdrücklich aus (§1), und `--print-config` zeigt heute keine
  Liste (Probe: nur das Objekt-Beispiel). Das Handbuch wurde in diesem Slice
  geliefert, nicht in der Release-Prep (der Kopf und die Versionstabelle bleiben
  dort). Solange die ADR `Proposed` ist, darf sie berichtigt werden; nach dem
  Statuswechsel nur per Folge-ADR (`AGENTS.md` §3.5). **Vor** dem Wechsel auf
  `Accepted` den Satz auf den Stand bringen (Generatoren liefern weiter die
  Objekt-Form; Handbuch ist geliefert).
- **V-2 — LOW — Fehlermeldung eines falschen Listen-Eintrags nennt keinen Index.**
  `requirements: ["x"]` meldet `line 3: cannot unmarshal !!str "x" into
  configyaml.rawTraceRequirements`, ohne `trace.requirements[0]`; dieselbe interne
  Typbezeichnung steht schon in der Objekt-Form (M3, byte-identisch zum Vorzustand),
  dort ist sie Bestand. Das Lastenheft-Kriterium „Negative (Eintrag)“ verlangt den
  Index für Header-Fehler und Nullmenge, nicht für eine falsche YAML-Form des
  Eintrags; daher keine DoD-Verletzung.
- **V-3 — INFO (Closure-Eingang) — Risiko 1 aus §6 ist bestätigt als entfallen.** Die
  Objekt-Form liefert über Repo, Heading- und Tabellen-Fixture sowie alle kaputten
  Konfigurationen byte-identische Ausgabe (M2, M3). Die Gegenprobe der Heading-Fixture
  gegen den Vorzustand ist damit erbracht (§6 nennt `v0.86.1`; gemessen wurde gegen
  `68cae247`, den unmittelbaren Vorzustand — strenger, nicht schwächer).
- **V-4 — INFO (Closure-Eingang) — Risiko 3 aus §6.** Die Wechselwirkung von
  `modality` je Eintrag mit `--require-complete` ist durch Probe D und
  `TestTraceSourcesModalityPerEntry` gepinnt: kein bestehender Pfad nimmt Modalität
  global an. Risiko 2 (Slice wächst über eine Review-Sitzung) ist mit zwei
  Review-Runden eingetreten und aufgelöst worden; ein Ausgang steht noch in §7 aus.

Keine DoD-Verletzung.

## Kategorie-Summary

- **DoD-Verletzung:** 0
- **Plan-/ADR-Inkonsistenz:** 1 (V-1, MEDIUM; vor `Accepted` zu beheben)
- **Mängel unter der DoD-Schwelle:** 1 (V-2, LOW)
- **Closure-Eingänge:** 2 (V-3, V-4, INFO)
- **Offen bis Closure (keine Verletzung):** ADR-0108 `Accepted` (Kopf und Index), §7, Risiko-Ausgänge, Paarungen, Register

## Verdikt

**DoD 1, 2, 4 sind bestätigt, DoD 3 zur Hälfte (Handbuch); die ADR-Zeile und die
Closure-Punkte sind offen, wie vorgesehen.** Alle sieben Brech-Proben liefen mit der
behaupteten Meldung rot, die Objekt-Form ist gegen den Vorzustand byte-identisch,
die Listen-Probe nach dem Lastenheft-Beispiel liefert 21 + 6 = 27 mit den
zugesagten Exit-Codes, `make gates` ist grün. Der Slice kann in die Closure, wenn
V-1 vor dem Wechsel von ADR-0108 auf `Accepted` berichtigt ist.
