# slice-274: `trace.requirements` als Liste von Quellen — eine Datei, zwei Anforderungsformen in der RTM

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst; ein
repo-weiter Beleg über die DoD hinaus (Wellen-Trigger) fehlt.

**Bezug:** [`DC-FA-REQ-001`](../../../../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen)
(die Liste der Quellen),
[`DC-FA-CLI-009`](../../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
(die RTM, die ihre Vereinigung führt),
[`DC-FA-CLI-011`](../../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
(`--require-complete` wertet die Vereinigung aus),
[`DC-FA-MOD-001`](../../../../../spec/lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in)
(`modality` je Eintrag),
[ADR-0108](../../../adr/0108-trace-requirements-als-liste-von-quellen.md)
(die Entscheidung; Status `Proposed`, bis dieser Slice schließt),
[ADR-0037](../../../adr/0037-trace-tabellenquellen-nullmengen-guard.md) (der
Quellenvertrag, der je Eintrag gilt),
[`DC-QA-02`](../../../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).
Anlass: eingehender CR eines Adopters, abgelegt unter `docs/plan/cr/` (kein
Spec-Verweis, die Spec nennt ihn nur als „CR eines Adopters").

**Berührte Spec-Stellen:** [`DC-FA-REQ-001`](../../../../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen)
(Lastenheft 0.106.0); [§`DC-FA-REQ-001.a`](../../../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 8 und die Schlüsseltabelle `trace.requirements`.

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-20.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** `trace.requirements` nimmt eine nichtleere Liste von Quellen an; die
RTM führt ihre Vereinigung, und Lastenheft 0.106.0 samt
[ADR-0108](../../../adr/0108-trace-requirements-als-liste-von-quellen.md) sind
umgesetzt — so dass eine Datei Anforderungen als Überschriften **und** als
Tabellenzeilen führen kann, ohne umgeformt zu werden.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Änderung von Lastenheft, Spezifikation und ADR-Text.** Sie sind mit der
  Annahme des CR geschrieben und tragen den Vertrag; dieser Slice setzt ihn
  um. Findet die Umsetzung eine Lücke im Vertrag (etwa die Referenzseite mit
  mehreren Mustern, die erst beim Lesen des Codes auffiel), ist das eine
  Änderung am Vertrag und gehört vor den Code in einen eigenen Commit — nicht
  in den Code-Commit.
- **`CHANGELOG.md`, READMEs und Handbuch-Kopf.** Sie werden in der Release-Prep
  gepflegt, nicht im Feature-Commit
  ([`AGENTS.md`](../../../../../AGENTS.md) §5 Regel 17). Der **Inhalt** des Handbuchs
  zur Listen-Form ist dagegen Teil der DoD, weil der CR ihn verlangt.
- **Eine `first`/`last`-Politik über Einträge und die Vererbung von
  `modality`.** Das Lastenheft schließt beide aus (Out-of-Scope von
  [`DC-FA-REQ-001`](../../../../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen)); kein Adopter braucht sie.
- **`--print-config`/`--suggest-config` erzeugen keine Listen-Form.** Bestand
  bleibt bewusst stehen: Die Generatoren liefern die Objekt-Form, die für einen
  einzelnen Quellentyp richtig ist; eine Liste ist Handarbeit für den
  gemischten Bestand. Eine Anpassung wäre ein anderer Vorgang an einem anderen
  Gegenstand.
- **Die Umstellung der Konfiguration des Adopters** und die Antwort an ihn.
  Das ist ein anderer Vorgang in einem anderen Repo; erst das Release macht ihn
  möglich.
- **Strikter Modus für Listen-Einträge ohne explizite `source`.** Parität mit
  der Objekt-Form ([ADR-0108](../../../adr/0108-trace-requirements-als-liste-von-quellen.md), Grenze): ein Eintrag mit `format: headings` und
  ohne `source` bleibt bei null Treffern leer, statt Exit 2 zu liefern.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [x] **Liste laden:** Der Config-Adapter akzeptiert `trace.requirements` als
      Objekt oder nichtleere Liste, verwirft jede andere Form und die leere
      Liste mit Exit 2 und validiert jeden Eintrag mit seinem Index
      (AK-3, AK-7); das Modell trägt die Einträge.
- [x] **Vereinigung führen:** Der Kern extrahiert je Eintrag, bildet die
      Vereinigung in Listenreihenfolge, meldet die doppelte Kennung über
      Einträge und die Nullmenge je Eintrag (Exit 2, Index in der Meldung),
      wendet `modality` je Eintrag an und lässt die Referenzseite die Muster
      aller Einträge erkennen (AK-1, AK-2, AK-4, AK-6); die Objekt-Form
      liefert byte-identische Ausgabe (AK-5).
- [x] **Handbuch und ADR:** Das Benutzerhandbuch (§5 `trace`, Konfigurationsbeispiel
      zur `.d-check.yml`) beschreibt die Listen-Form, die doppelte Kennung und
      den Nullmengen-Guard je Eintrag; [ADR-0108](../../../adr/0108-trace-requirements-als-liste-von-quellen.md)
      steht auf `Accepted`, Kopf **und** Index.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review; Verifikation.
- [x] Beobachtungs-Register (`../../observations/`) fortgeschrieben — oder
      „keine Beobachtung angefallen" in §7 notiert.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/core/model/config.go` | update | `TraceConfig` trägt die Liste der Quellen-Einträge; der Nullwert bleibt die Objekt-Form (AK-5) |
| `internal/adapter/driven/configyaml/configyaml.go` | update | `rawTraceRequirements` nimmt Objekt **oder** Sequenz an (eigener `UnmarshalYAML`), strikt; leere Liste und fremde Form Exit 2; `applyTraceRequirements` je Eintrag mit Index (AK-3, AK-7) |
| `internal/adapter/driven/configyaml/` Testdatei | update | Happy/Boundary/Negative nach AK-3, AK-7; Weiche Objekt↔Liste gegen ihre Fehlerformen |
| `internal/hexagon/core/app/trace.go`, `trace_table.go` | update | `resolveTrace`/`traceRequirements` je Eintrag, Vereinigung, Kreuz-Duplikat, Nullmenge je Eintrag, `modality` je Eintrag; `reqPat` der Referenzseite wird zur Vereinigung der Muster |
| `internal/hexagon/core/app/trace_cross.go` | update | `cross-consistency` liest das Muster über dieselbe Vereinigung statt eines einzelnen `reqPat` |
| `internal/hexagon/core/app/` Testdateien | update | AK-1, AK-2, AK-4, AK-6; beide Seiten der Kreuz-Duplikat-Prüfung (Überschrift zuerst **und** Tabelle zuerst) |
| `internal/adapter/driving/cli/cli_acceptance_test.go` | update | `--trace` und `--trace --require-complete` über die Binärgrenze (AK-1, AK-6) |
| `docs/user/benutzerhandbuch.md` | update | Listen-Form, S-3, S-4 im §5 `trace` und im Konfigurationsbeispiel |
| `docs/plan/adr/0108-…md`, `docs/plan/adr/README.md` | update | `Accepted`, Geschichte-Zeile, Index-Status — erst, wenn der Mechanismus steht |

**Spiegel vor dem Editieren auflisten**
([`MR-025`](../../../../../harness/conventions.md#mr-025), und Schritt 17 des
Workflow-Skeletts): wer die Semantik von `trace.requirements` ändert, greppt
vorher `trace.requirements`/`requirements.format`/`rawTraceRequirements` über
`harness/`, `docs/user/` (auch `operations.md`), `.claude/`, die Skript-Kommentare
und `tools/` und listet jeden Treffer im Plan, bevor der erste Code-Commit
fällt. Ausgangsliste der Treffer beim Anlegen: `docs/user/benutzerhandbuch.md`,
`docs/user/operations.md`, `internal/adapter/driving/cli/cli_acceptance_test.go`,
`internal/adapter/driven/configyaml/abdeckung_test.go`,
`tools/archive-wave/stub.go` (zu prüfen, ob Spiegel oder Zufall).

**Plan-Änderung vor dem Code (2026-10-20, Schritt 15 des Workflow-Skeletts):**
Beim Lesen des Codes weicht der Zuschnitt von der Tabelle ab, ohne den Umfang
zu ändern — der Diff bleibt in den zwei Schichten Adapter und Kern.

- **Neu:** `internal/hexagon/core/app/trace_sources.go` trägt die Einträge, die
  Vereinigung, die Kreuz-Duplikat-Prüfung und die Modalität je Eintrag;
  `trace.go` bleibt dadurch schlank (die Funktionen `loadTraceRequirements`,
  `traceModalities`, `traceRequirements` wandern dorthin).
- **Entfällt:** `trace_cross.go` und `trace_table.go` bleiben unverändert —
  `cross-consistency` liest das Muster der Referenzseite über `rt.reqPat`
  (jetzt die Vereinigung), und der Tabellen-Extraktor kennt keinen Eintrag.
- **Neue Testdateien:** `internal/hexagon/core/app/trace_sources_test.go`,
  `internal/adapter/driven/configyaml/trace_requirements_list_test.go`,
  `internal/adapter/driving/cli/cli_trace_sources_test.go`.
- **Generiert:** `docs/user/abdeckung-tests.md` (`make abdeckung`); `make test`
  hält die Datei gegen ihre Ableitung, die neuen Tests ändern sie.
- **Spiegel (Schritt 17, `grep` über `harness/`, `docs/user/`, `.claude/`,
  `tools/`, `README.md` nach `trace.requirements`/`requirements.format`/
  `requirements.id-pattern`/`requirements.source`):** `docs/user/benutzerhandbuch.md`
  (§4.12 „Andere Repo-Konvention", §5 *Konfigurierbare RTM-Quellen* samt
  Nullmengen-Absatz, §7 Fehlerbild „ergab 0 Anforderungen") und
  `docs/user/operations.md` (Zeile `--trace`) beschreiben die Objekt-Form und
  werden nachgezogen. Die Zeilen der Handbuch-Versionstabelle sind Release-Prep.
  `tools/archive-wave/stub.go` nennt `trace.requirements.id-pattern` nur als
  Herkunft eines Kennungs-Musters und beschreibt keinen Mechanismus, der sich
  ändert — Zufallstreffer, bleibt.

**Plan-Änderung nach Review-Runde 1 (2026-10-20, Schritt 20):** Der Report
`docs/reviews/2026-10-20-slice-274-trace-requirements-liste.md` meldet zwei
MEDIUM und drei LOW; alle werden behoben, die Abgrenzung in §1 bleibt.

- **F-1 und F-3 (gleiche Aussage: die Neu-Serialisierung des Teilbaums):**
  `decodeStrictNode` serialisiert nicht mehr neu. Er dekodiert den Knoten mit
  `Node.Decode` (löst Aliase und Merge-Keys auf) und prüft die Schlüssel
  getrennt am Knotenbaum gegen die `yaml`-Tags des Zieltyps — Zeilen kommen
  dann aus dem Dokument, nicht aus einer Kopie. Aussagen-Suche:
  `grep -rn "Neu-Serialisierung\|decodeStrictNode\|yamlLineRef"` über
  `internal/`, `docs/user/`, `spec/`, `harness/` — Treffer nur in
  `configyaml.go`.
- **F-2:** Tests für Range-Validierung und `cross-consistency` in Listen-Form.
- **F-4:** Das Kreuz-Duplikat wird nach dem Extrahieren **aller** Einträge
  geprüft (Spezifikation Schritt 8).
- **F-5:** Die Test-Kommentare nennen Gegenstand statt `AK-<N>`
  (`grep -rn "AK-[0-9]" internal/`: 11 Treffer in drei Testdateien, ein
  Zufallstreffer `LH-AK-002` in `rules/markdown_test.go` bleibt).
- **F-6:** Die Tests belegen den Fall in AK-6 über Exit-Code und Waisenzahl und
  pinnen die [ADR-0108](../../../adr/0108-trace-requirements-als-liste-von-quellen.md)-Grenze
  (Eintrag mit `format: headings` ohne `source`). Der Fehlerzweig der
  Mustervereinigung ist gepinnt (Runde 2, F-1): ein Einzelmuster mit 999
  verschachtelten Gruppen kompiliert allein, die Vereinigung scheitert an der
  Schachtelungstiefe.
- **Zusätzlich:** Weil der Fix einen neuen Leseweg schreibt, läuft nach ihm eine
  zweite `reviewer`-Runde über genau die Fix-Commits (Schritt 14).

**Rot-Beleg** ([Modul 11 §Bewusstes Brechen](../../../../../.harness/baseline/v6.18.1/regelwerk/modul-11-verification.md#bewusstes-brechen-für-dod-testbehauptungen-modul-11)):
Die Kreuz-Duplikat-Prüfung und der Nullmengen-Guard je Eintrag werden
testweise entfernt; die benannten Tests müssen mit der behaupteten Meldung
rot laufen, nicht nur irgendwie.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen Slice.
Lastenheft 0.106.0, Spezifikation und [ADR-0108](../../../adr/0108-trace-requirements-als-liste-von-quellen.md) liegen vor und sind auf dem
Hauptzweig (Vorgang davor).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Umsetzung der
  Referenzseite (`trace_cross.go`, Coverage-Ranges) braucht mehr als eine
  mechanische Umstellung auf die Mustervereinigung — dann trennt sich diese
  Hälfte als eigener Slice ab. Wächst der Diff über das, was eine
  Review-Sitzung prüft, gilt die Form von
  [`MR-066`](../../../../../harness/conventions.md#mr-066): Grund **und** benannte
  Ersatz-Form der Prüfung im Plan, vor dem Code.
- `in-progress` → `open` (blockiert): Die Spezifikation erweist sich beim
  Umsetzen als widersprüchlich (etwa zwischen Schritt 8 und der bestehenden
  Fehlerpräzedenz) — dann zuerst die Spec berichtigen, nicht am Code vorbei
  entscheiden.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft; [ADR-0108](../../../adr/0108-trace-requirements-als-liste-von-quellen.md)
steht auf `Accepted`. Geht mit dem nächsten Release hinaus; die Release-Prep
trägt `CHANGELOG.md` und den Handbuch-Kopf.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Die Referenzseite ändert das Verhalten der Objekt-Form.** Der Umbau von
  `reqPat` auf eine Mustervereinigung kann die bisherige Ganz-Token-Erkennung
  verschieben (Alternation mit Linksvorrang statt je Muster geprüft) und damit
  die Byte-Identität der Objekt-Form brechen. — **Ausgang:** entfallen — bei
  einem Eintrag liefert `referencePattern` das Muster selbst, unverändert; die
  Verifikation verglich `--trace`, `--require-complete`, `--json` und
  `--doctor` für das Repo selbst, eine Heading- und eine Tabellen-Fixture gegen
  ein Image von `68cae247` (stdout, stderr, Exit getrennt, `cmp`): identisch.
- **Der Slice wächst über eine Review-Sitzung** (Adapter, Kern, Cross, Handbuch,
  ADR). — **Ausgang:** entfallen — `trace_cross.go` und `trace_table.go`
  blieben unverändert, der Diff blieb in zwei Schichten; Runde 1 las den
  gesamten Diff in einer Sitzung, Runde 2 genau die Fix-Commits (Schritt 14).
  Die Rückführung nach §4 wurde nicht gezogen.
- **Wechselwirkung `modality` je Eintrag mit `cross-consistency` und
  `--require-complete`.** — **Ausgang:** entfallen — `cross-consistency` liest
  keine Modalität; `--require-complete` gatet je Eintrag
  (`TestTraceSourcesModalityPerEntry`, Rot-Beleg D der Verifikation: ohne die
  Eintrags-Modalität wird er rot).

## 7. Closure-Notiz

- **Was hat funktioniert:** `trace.requirements` nimmt eine Liste von Quellen
  an; die RTM führt ihre Vereinigung in Listenreihenfolge, eine doppelte Kennung
  über Einträge ist Exit 2, der Nullmengen-Guard gilt je Eintrag mit Index, und
  die Objekt-Form liefert byte-identische Ausgabe (gemessen gegen `68cae247`).
  Die Strikt-Prüfung am Knotenbaum hält unbekannte Schlüssel auch hinter Alias
  und Merge-Key fest. Sieben Rot-Belege der Verifikation (Kreuz-Duplikat,
  Nullmenge je Eintrag, leere Liste, `modality` je Eintrag, Mustervereinigung,
  Strikt-Prüfung, Index) laufen aus dem richtigen Grund rot.
- **Was ging anders als geplant:** (1) Beim Lesen des Codes fiel auf, dass
  Referenzscan, Coverage-Ranges und `cross-consistency` ein einziges Muster
  lesen; die Vereinigung der Muster kam als Vertragszusatz in Spezifikation und
  ADR (Entscheidung 5) vor den Code. (2) Die erste Strikt-Prüfung
  serialisierte jeden Eintrag neu und zerbrach damit Aliase und Merge-Keys
  sowie die Zeilennummern der Objekt-Form (Review Runde 1, F-1, F-3); sie wurde
  durch eine Prüfung am Knotenbaum ersetzt, daraufhin lief Runde 2 über die
  Fix-Commits. (3) Die Verifikation fand in [ADR-0108](../../../adr/0108-trace-requirements-als-liste-von-quellen.md) eine Konsequenz, die
  `--print-config` eine Listen-Form zusprach (V-1); sie wurde vor `Accepted`
  berichtigt, solange die ADR editierbar war.
- **Steering-Loop-Eintrag:** keiner mit neuer Schwelle.
- **Beobachtungs-Register (`../../observations/`):** `evidence/slice-274.md` in
  [`BEO-ALL/haertung-kippt-fehlerpolitik-ungeprueft`](../../observations/BEO-ALL/haertung-kippt-fehlerpolitik-ungeprueft/state.md)
  (jetzt 5×): ein neuer Leseweg (Neu-Serialisierung) wurde nur gegen den Fall
  des CR gefahren; der fremde Leser über den Code-Commit fing ihn (Schritt 14).
- **Folge-Slices:** keiner.
- **Risiken aus §6:** alle drei entfallen (siehe §6). Trigger-Audit: kein
  Carveout, kein bootstrap-aware Gate; [ADR-0108](../../../adr/0108-trace-requirements-als-liste-von-quellen.md) ist mit diesem Slice
  `Accepted` und trägt ihren Re-Evaluierungs-Trigger; keine Hard Rule mit
  eingetretenem Trigger. Nachtlauf-Stand
  ([`MR-053`](../../../../../harness/conventions.md#mr-053)): wie in §8.
- **Drei Paarungen:** (a) Anker — kein Eintrag mit `liegt in`; (b)
  Folge-Slices — keiner genannt; (c) Register — die zitierte Beobachtung
  existiert und trägt fünf Belege.

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind der Config-Adapter, der
Trace-Kern, das Benutzerhandbuch und die ADR-Ablage — sämtlich unter dem
Default `*` (`ALL`); keine eigene Sub-Area erfüllt das Inklusionskriterium
(`tools/harness/` wird nicht berührt).

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen;
Treffer unter `BEO-ALL`:

- [`semantic-change-body-only-edges-stale`](../../observations/BEO-ALL/semantic-change-body-only-edges-stale/observation.md)
  (22×, verkörpert als [`MR-025`](../../../../../harness/conventions.md#mr-025)): trifft diesen Slice unmittelbar — eine
  Semantik-Änderung an `trace.requirements` hat Spiegel (Handbuch, Tests,
  Kommentare); sie werden in §3 **vor** dem Editieren aufgelistet.
- [`guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`](../../observations/BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen/observation.md)
  (4×, verkörpert, Schritt 17 des Workflow-Skeletts): derselbe Grep über
  `harness/`, `docs/user/`, `.claude/` und Skript-Kommentare.
- [`haertung-kippt-fehlerpolitik-ungeprueft`](../../observations/BEO-ALL/haertung-kippt-fehlerpolitik-ungeprueft/observation.md)
  (4×, verkörpert, Prüffrage 19): die Weiche Objekt↔Liste ist ein neuer
  Leseweg und wird gegen ihre Fehlerformen gefahren, nicht nur gegen den Fall
  des CR.
- [`negativtest-deckt-nur-eine-regelseite`](../../observations/BEO-ALL/negativtest-deckt-nur-eine-regelseite/observation.md)
  (1×, offen): die Kreuz-Duplikat-Prüfung hat zwei symmetrische Seiten
  (Überschrift/Tabelle); der Test variiert beide Reihenfolgen. Ein zweites
  Auftreten würde den Zähler auf 2× heben.
- [`adr-accepted-bevor-der-mechanismus-steht`](../../observations/BEO-ALL/adr-accepted-bevor-der-mechanismus-steht/observation.md)
  (1×, offen): [ADR-0108](../../../adr/0108-trace-requirements-als-liste-von-quellen.md) bleibt `Proposed`, bis die Umsetzung steht, und wird
  erst im letzten Commit `Accepted`.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../../harness/conventions.md#mr-053)):
`image-scan` grün (2026-10-20); `upstream-drift` rot — **planmäßig, gelesen**:
neuer Release der Baseline (Pin `v6.17.0`) sowie veraltete Pins für Go
(1.27.2 → 1.27.3), golangci-lint, semgrep und a-check. Fremd-Release und
Pin-Veralterung, kein Bezug zu diesem Slice, jeweils eigener Vorgang.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
