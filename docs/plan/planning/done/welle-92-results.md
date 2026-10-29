# Welle welle-92 — Namens-Slices im Trace und Home-Variablen in hostpaths — Closure-Notiz

**Welle:** welle-92
**Abschluss:** 2026-10-29
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — Ergebnis, nicht Tätigkeit.

- **slice-279** (Der Trace erkennt Slices mit Namens-Kennung):
  [`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
  erfüllt, Lastenheft 0.107.0, ADR
  [ADR-0110](../../adr/0110-trace-default-erkennt-namens-slices.md) Accepted. Ohne
  `trace.slices.file-pattern` gilt der Ziffern-Fall `slice-NNN-…md` wie bisher, sonst
  der Namens-Fall `slice-<name>.md`; Namens-Belege unter `observations/` sind keine
  Slices (GitHub-Issue 11).
- **slice-280** (`hostpaths` erkennt `$HOME` und `${HOME}`):
  [`DC-FA-HOST-001`](../../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in)
  erfüllt, Lastenheft 0.108.0, ADR
  [ADR-0111](../../adr/0111-hostpaths-variablen-home-opt-in.md) Accepted. Der neue
  opt-in-Schlüssel `hostpaths.variables` meldet `$NAME/…`, `${NAME}/…` und
  `"$NAME"/…` als Home-Verweis; ohne den Schlüssel ist der Befundsatz
  byte-identisch (GitHub-Issue 12).

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- **Beide Erkennungs-Änderungen als „ohne Schalter unverändert“ gemessen, nicht
  behauptet:** die RTM über eine Kopie des Planning-Baums ohne `trace`-Block (549
  Dateien `slice-*.md`, 270 Belege unter `observations/`) und der `hostpaths`-Lauf
  über Eigenbestand und Fixture sind gegen `v0.87.1` in Text, `--json` und `--yaml`
  byte-gleich. Das Mehr der Welle gegenüber den Slice-DoDs ist genau dieser gemeinsame
  Lauf.
- **Review in frischem Kontext, danach Verifier mit eigenen Mutationen:** beide Slices
  hatten Findings, die der Autor nicht sah (MEDIUM in beiden), und der Verifier hat
  jede behauptete Rot-Stelle einzeln gebrochen und gesehen.
- **Die Reihenfolge der Slices** (erst slice-279, dann slice-280) hielt die zwei
  Lastenheft-Fassungen konfliktfrei.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile möglichst mit der Konsequenz, die
daraus schon gezogen wurde.

- **Beide Negativlisten prüften, was die neue Regel nicht treffen darf, nicht was sie
  am bereits Erkannten ändert.** slice-279: der `observations`-Ausschluss stand vor dem
  Ziffern-Fall und hätte einen titelhaltigen Beleg aus der RTM genommen (Review F-1).
  slice-280: das Ziel eines Funds hinter `${HOME}/…` wechselt, ein Bestands-Glob trifft
  nicht mehr (Review F-2), überlappende Treffer meldeten doppelt (F-3). Konsequenz:
  beide behoben, Lastenheft und ADR benennen es; Beleg im Register (siehe unten).
- **Eine Plan-Annahme stimmte nur für eine von drei Schreibweisen:** die Verdrängung
  eines Unix-Treffers wurde mit dem `E` von `HOME` begründet, das gilt nur für
  `$HOME/…`. Ein Mutationstest, der grün blieb, deckte es auf (Plan-Änderung vor §4 in
  slice-280).
- **Das Repo prüft seine eigene Doku mit `hostpaths` und meldete meine Beispiele** mit
  dem Default-Präfix; sie stehen mit Platzhalter. Dogfooding als Gegenprobe.
- **Closure-Trigger-Punkt „Beide Issues tragen eine Antwort“ ist nicht erfüllt und
  verschoben:** eine Antwort mit Verweis auf das Ergebnis wäre vor Push, Tag und
  Release falsch (der Stand liegt nur lokal). Die Verschiebung steht im Drift-Log der
  Roadmap; die Antworten gehören zur Release-Prep.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — hier stehen **nur** Beobachtungen, die im Register
3× erreicht haben.

Kein Eintrag über der Schwelle, der noch auf einen Ausgang wartet. Gelesen wurde das
Register vollständig: Einträge mit Stand `offen` und mindestens drei Belegen gibt es
keine. Gezählt, nicht neu verkörpert:

- `BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet` steht bei 7×, ist
  verkörpert (Schritt 19 des Workflow-Skeletts, Prüffrage 21 im Reviewer-Skill) und
  bekommt zwei Belege derselben **neuen Richtung** (slice-279, slice-280): die
  Wirkung der neuen Regel auf den bisher erkannten Bestand. Tritt sie ein drittes Mal
  auf, gehört die Frage *was ändert sich am Bestand der Funde?* in Schritt 19.
- `BEO-ALL/begruendung-traegt-entscheidung-nicht` steht bei 4×, ist verkörpert
  (Prüffrage 20) und bekommt einen Beleg (slice-280).
- `BEO-ALL/dateiname-als-inhalts-beleg-akzeptiert` steht bei 2×, offen, mit einem
  Beleg aus slice-279 (der Dateiname entscheidet nicht, ob eine `slice-*.md` ein
  Slice ist).

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Zähler steht in [`../observations/`](../observations/).
Was in dieser Welle **3×** erreicht hat, steht oben unter
*Steering-Loop-Einträge*: nichts Neues.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur, das
Original ist die Slice-Datei.

Keine. Offen bleibt keine Arbeit der Welle; die Erkennung in `reviews`,
`--suggest-config` und `tools/archive-wave` bleibt ziffern-only (Abgrenzung beider
Slices, andere Module).

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren Anker.

- `make fullbuild` Exit 0 am Stand der letzten Slice-Closure — Image-Hash
  `sha256:fb85652169b8cef7f3220b52fc99c88eadc7ddc665110f550b87952fb917d055`
  (`make gates`, `make image-test`, `make bench`, `make completeness-check` und
  `make verify-closure-notes` darin grün).
- `make blackbox-probe REF=v0.87.1`: byte-identisch über 20 Vergleiche (stdout, stderr
  und Exit getrennt, Kanarienlauf vorher und nachher).
- Die RTM einer Planning-Baum-Kopie ohne `trace`-Block ist in Text, `--json` und
  `--yaml` byte-gleich zu `v0.87.1`; das Fixture aus Issue 11 kippt von `WAISE` und
  Exit 1 auf `ok` und Exit 0 unter `--trace --require-complete`.
- `hostpaths` ohne den Schlüssel: Fixture und Eigenbestand byte-gleich zu `v0.87.1`;
  das Fixture aus Issue 12 mit `variables: [HOME]` meldet fünf Funde (drei Formen,
  Inline-Code, Tilde).
- Replay-Lauf gegen Golden Set: entfällt, d-check führt keinen (deterministischer Kern,
  [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- Coverage gesamt: 94,70 % bei Schwelle 93 % (`make coverage-gate`), offene Carveouts: keine.

**Trigger-Audit (Schritt 2), je Klasse geprüft:**

- **Carveout:** keine aktiven (`docs/plan/carveouts/` trägt nur `done/` und die README).
- **Bootstrap-aware Gate:** `coverage-gate` ist kalibriert (Schwelle 93 %, Senkung nur per
  ADR); kein Hochschalt-Trigger gefeuert, keine neue Stufe.
- **ADR:** 91 ADRs tragen einen Re-Evaluierungs-Trigger; keine Bedingung feuerte in der
  Welle. Gelesen wurden die berührten:
  [ADR-0098](../../adr/0098-hostpaths-home-relativ-und-ziel-ventil.md) (Ventil wird nicht
  zur Pflicht-Konfiguration; permanent), die Trace-ADRs 0034 und 0035 tragen als
  Bestand vor Einführung der Pflicht keinen Trigger (grandfathered); die neuen
  [ADR-0110](../../adr/0110-trace-default-erkennt-namens-slices.md) und
  [ADR-0111](../../adr/0111-hostpaths-variablen-home-opt-in.md) tragen ihre Trigger.
- **Hard Rule:** in `AGENTS.md` §3 tragen alle einen Auflösungs-Trigger *permanent* oder
  benennen die ungedeckte Kategorie; nichts feuerte, keine Zeile wird entfernt.

**Die drei Paarungen:** (a) Anker — kein Steering-Loop-Eintrag mit `liegt in`-Feld, keine
Paarung auszulösen; (b) Folge-Slice — keiner genannt; (c) Register — die zitierten
Beobachtungen existieren als Verzeichnisse und tragen alle Belege.
