# Slice slice-279: Der Trace erkennt Slices mit Namens-Kennung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-92

**Bezug:** [`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix) (Scope),
[`DC-FA-CLI-011`](../../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code) (der Exit-Code, an dem das Feld es bemerkt hat),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus) (Byte-Identität für den Ziffern-Fall). Keine *aktive* ADR trägt die
Erkennungsform; eine begleitende ADR hält die Entscheidung fest (Nummer bei der Anlage).

**Berührte Spec-Stellen:** [`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
(Beschreibung, Akzeptanzkriterium *Default byte-identisch*) ·
[`DC-FA-CLI-009.a`](../../../../spec/spezifikation.md#dc-fa-cli-009a--requirements-traceability-matrix) Schritt 2 und die Schema-Zeile
`trace.slices.file-pattern` (`spezifikation.md §2`).
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-29.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `d-check --trace` ordnet ohne `trace.slices.file-pattern` auch einen
Slice mit Namens-Kennung (`slice-<name>.md`) einer Anforderung zu — die Kennung
ist der Dateiname ohne `.md` —, während Ziffern-Slices (`slice-NNN-titel.md`)
ihre Kennung `slice-NNN` und die RTM eines Repos ohne Namens-Slices ihre
Ausgabe behalten.

**Gemessener Anlass** (GitHub-Issue 11, am Image `v0.87.1` nachgestellt): ein
Slice `slice-qa-sicherheit.md` mit `TDS-QA-03` im Kopf gilt als `WAISE`;
dieselbe Datei unter `slice-001-sicherheit.md` als `ok`. Mit
`trace.slices.file-pattern: '^(slice-[a-z0-9-]+)\.md$'` ist sie `ok` und
`--trace --require-complete` endet mit Exit 0 — der Schlüssel, den das Issue
als Wunsch nennt, existiert also; offen ist der **Default**.

**Entscheid des Auftraggebers (2026-10-29):** den Default erweitern, statt nur
den Schlüssel auffindbar zu machen.

**Die Erkennung — Entwurf, im Review zu prüfen, nicht gesetzt:**

- Passt der Basisname auf `^slice-\d+-.*\.md$`, gilt der **Ziffern-Fall** wie
  bisher: Kennung `slice-\d+`.
- Jede andere Datei `slice-*.md` ist ein **Namens-Fall**: Kennung = Basisname
  ohne `.md`.
- Dateien unter einem Pfadsegment `observations` unterhalb von `slices.dir`
  sind im Default **keine** Slices. Der Grund ist gemessen, nicht vermutet:
  im Beobachtungs-Register heißt jeder Beleg `evidence/slice-<Kennung>.md`. In
  diesem Repo sind es 270 von 549 Dateien `slice-*.md` unter
  `docs/plan/planning/` (Zählkommando in §3); in einem Repo mit Namens-Kennung
  trüge jeder davon den Namen eines Slice und deckte, was er nennt. Ohne den
  Ausschluss würde der Slice Waisen **verdecken**, wo das Issue Slices
  **übersah** — die schlimmere der beiden Fehlerrichtungen.
- Ein gesetztes `trace.slices.file-pattern` gilt wie bisher unverändert und
  ohne Ausschluss: wer ein Muster setzt, entscheidet die Menge selbst.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein neuer Konfig-Schlüssel** — der Wunsch 2 des Issues ist durch
  `trace.slices.file-pattern` erfüllt (am Image belegt); ein zweiter Schlüssel
  für denselben Zustand wäre eine zweite Quelle.
- **Eine Warnung für `slice-*.md`-Dateien, die das Muster nicht trifft** — eine
  neue Ausgabe-Klasse im `--trace`-Lauf bräche die Byte-Identität des
  Ziffern-Falls, und der Slice beseitigt die Ursache, statt sie zu melden.
- **`reviews`, `--suggest-config` und `tools/archive-wave`** tragen dieselbe
  ziffern-only Form (`slice-([0-9]+)` bzw. `slice-(\d+)` bzw. `slice-\d{3}`;
  Zählkommando in §3). Das sind andere Module und andere Anforderungen, und
  das Feld hat nur den Trace gemeldet — Klasse 3, ein anderer Vorgang. Die
  Review-Deckung kennt mit `match: name` bereits einen Namens-Weg.
- **Handbuch, README und CHANGELOG** — Release-Prep, nicht Feature-Commit
  (`AGENTS.md` §5 Zeile 17); der Handbuch-Inhalt zum Default wird dort
  nachgezogen.
- **Die Prüfung, ob eine `slice-*.md` ein Slice-Plan ist** — am Dateinamen
  nicht entscheidbar; wer Notiz-Dateien mit `slice-`-Präfix im Planning-Baum
  führt, setzt `file-pattern`. Benannte Grenze, kein Mangel dieses Slice.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix) erfüllt, Test referenziert:
      Happy (Namens-Slice deckt die Anforderung, `--require-complete` Exit 0) ·
      Boundary (Ziffern-Slice unverändert `slice-NNN`) · Negative (Beleg unter
      `observations/` deckt nichts; gesetztes `file-pattern` übersteuert).
      Rot-Beleg: der Namens-Test läuft ohne die Änderung aus dem richtigen
      Grund rot (`WAISE` statt `ok`).
- [ ] Fassung nachgezogen: Lastenheft ([`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix),
      Kriterium *Default byte-identisch* auf den Ziffern-Fall gefasst, Version,
      Historie-Zeile), Spezifikation ([`DC-FA-CLI-009.a`](../../../../spec/spezifikation.md#dc-fa-cli-009a--requirements-traceability-matrix) Schritt 2, Schema-Zeile),
      begleitende ADR samt ADR-Index.
- [ ] Spiegel der Default-Form nachgezogen: das Konfigurations-Gerüst
      (`--print-config`) beschreibt den neuen Default.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für die Schnittstelle `--trace` falls öffentlicher Vertrag berührt (Spec und Gerüst oben; Handbuch ist Release-Prep).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/core/app/trace.go` | update | Default-Erkennung: Ziffern-Fall wie bisher, sonst Namens-Fall, `observations`-Segment im Default ausgenommen — [`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix) |
| neue Testdatei `trace_slice_names_test` im Paket `internal/hexagon/core/app` | neu | Happy/Boundary/Negative nach [`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix); Negativliste unten |
| `internal/adapter/driving/cli/cli_acceptance_test.go` | update | Akzeptanz *Default byte-identisch* am Ziffern-Fall, neues Kriterium am Namens-Fall |
| `internal/adapter/driving/cli/config_template.go` | update | Gerüst-Kommentar zu `trace.slices.file-pattern` nennt den bisherigen Default (Spiegel) |
| `spec/lastenheft.md` | update | [`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix): Beschreibung, Kriterium, Version, Historie-Zeile ([`MR-032`](../../../../harness/conventions.md#mr-032)) |
| `spec/spezifikation.md` | update | [`DC-FA-CLI-009.a`](../../../../spec/spezifikation.md#dc-fa-cli-009a--requirements-traceability-matrix) Schritt 2, Schema-Zeile `trace.slices.file-pattern`, Historie |
| `docs/plan/adr/` · `docs/plan/adr/README.md` | neu · update | begleitende ADR (Alternativen: Schlüssel auffindbar machen · Default erweitern · Warnung); Index nachgezogen |

**Zählung der Spiegel (Schritt 17/18 des Workflow-Skeletts), vor dem Code
nachzuzählen:**

- Alle `slice-*.md` unter `docs/plan/planning/`: `find docs/plan/planning -name 'slice-*.md' | wc -l` → 549.
- Davon auf den Ziffern-Fall: `find docs/plan/planning -name 'slice-*.md' -printf '%f\n' | grep -cE '^slice-[0-9]+-.*\.md$'` → 279.
- Der Rest, alle unter `observations/<BEO>/evidence/`: `find docs/plan/planning -name 'slice-*.md' ! -regex '.*/slice-[0-9]+-[^/]*\.md' -printf '%h\n' | sed -E 's#/BEO-[A-Z]+/[^/]+/evidence#/<BEO>/evidence#' | sort | uniq -c` → eine Zeile, 270.
- Fundstellen der Default-Form: `grep -rnF '(slice-\d+)-' --include=*.go --include=*.md --include=*.yml spec docs/user harness .claude internal .d-check.yml | grep -v _test` → `trace.go`, `config_template.go`, `spec/spezifikation.md` (Schema-Zeile) und `.d-check.yml` (dieses Repo setzt das Muster **selbst**, bleibt unberührt).
- Weitere ziffern-only Formen **ohne** diesen Gegenstand: `internal/hexagon/core/rules/reviews.go` (`slice-([0-9]+)`), `tools/archive-wave/collect.go` (`slice-(\d+)`), `internal/hexagon/core/app/suggest.go` (`slice-\d{3}`).

**Negativliste der Erkennung (Schritt 19, vor dem Code-Commit zu testen):**
Beleg `observations/BEO-X/y/evidence/slice-qa-sicherheit.md` · Datei
`slice-273.md` ohne Titel außerhalb von `observations/` (heute nicht erkannt,
künftig Namens-Fall mit Kennung `slice-273` — im Eigenbestand null Treffer,
Kommando oben) · Basisname mit Punkt im Stamm (`slice-a.b.md`) · `slice-.md` ·
Großschreibung (`Slice-X.md`, kein Treffer wie heute) · Archiv-Stub und Plan mit
gleichem Basisnamen (eine Kennung, kein Doppel) · Präfix-Paar (`slice-a-foo.md`
und `slice-a-foo-bar.md`: zwei Kennungen, kein Präfix-Treffer) · Datei im
Verzeichnis `observations` einer konfigurierten fremden `slices.dir`. Dazu der
Lauf über den eigenen Bestand: Ausgabe von `--trace` ohne `trace`-Block gegen
die Planning-Baum-Kopie vorher/nachher byte-gleich.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): welle-92 ist eröffnet, und die Lifecycle-Verzeichnisse
tragen keinen anderen Slice in `in-progress/`. Der Übergang landet auf dem Hauptzweig,
vor der Arbeit.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Negativliste
  deckt mehr als eine Klasse falscher Treffer außerhalb von `observations/`
  auf (etwa Notiz-Dateien, die sich am Namen nicht trennen lassen) — dann
  gehört die Frage *was ist ein Slice* in einen eigenen Slice, nicht in eine
  wachsende Regex.
- `in-progress` → `open` (blockiert — Carveout?): der Auftraggeber widerruft den
  Entscheid zum Default (§1) oder die Fassung des Lastenhefts wird nicht
  freigegeben.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, Review-Report liegt vor, und am gebauten
Image liefert das Fixture aus Issue 11 `ok` (Exit 0 bei `--require-complete`),
während die RTM einer Planning-Baum-Kopie ohne `trace`-Block gegen `v0.87.1`
byte-gleich ist. Dazu die Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Ein Beleg- oder Stub-Dokument mit Namens-Form deckt eine Anforderung, die
  kein Slice-Plan trägt, und verdeckt damit eine Waise still. Die
  `observations`-Ausnahme fängt den gemessenen Fall; ob es einen zweiten Ort
  gibt (Review-Reports liegen unter `docs/reviews/`, außerhalb von
  `slices.dir`), ist an der Menge der Dateien unter `docs/plan/planning/`
  nachzuzählen, bevor der Ausschluss als vollständig gilt — **Ausgang:** wird
  bei der Closure gesetzt (eingetreten / entfallen / weiter offen)
- Ein Repo mit Namens-Slices und ohne `trace`-Block bekommt eine andere RTM als
  zuvor — beabsichtigt (Waisen werden zu `ok`), aber nicht byte-identisch zum
  Vorgänger. Das Kriterium *Default byte-identisch* im Lastenheft muss darauf
  gefasst werden, sonst widerspricht die Spec dem Verhalten — **Ausgang:** wird
  bei der Closure gesetzt (eingetreten / entfallen / weiter offen)
- Eine Datei `slice-*.md`, die kein Slice-Plan ist (Notiz, Vorlage), wird zum
  Slice und deckt zu großzügig. Am Namen nicht entscheidbar; der Weg ist ein
  gesetztes `file-pattern` — **Ausgang:** wird bei der Closure gesetzt
  (eingetreten / entfallen / weiter offen)

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist der Default `*` (`ALL`) —
Trace-Kern, Spec, Gerüst; keine Sub-Area mit eigener Modus-Deklaration
(`tools/harness/` bleibt unberührt). Eine Sub-Area, die ausdifferenziert werden
müsste, gibt es hier nicht.

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen, Sub-Area
`*`. Zwei Treffer:
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/observation.md)
(5×, verkörpert in Schritt 19 des Workflow-Skeletts — die Negativliste in §3
ist seine Anwendung; der Slice schreibt eine Erkennungsform aus einem
Issue, das genau die Positivfälle nennt) und
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
(3×, verkörpert in `AGENTS.md` §3.8 — der `observations`-Ausschluss ist eine
Eingabe, die der Trace liest und nicht als Slice zählt; die Frage *gilt dort
dieselbe Zusage?* beantwortet §1). Erreicht keiner der beiden mit diesem Slice
eine neue Schwelle, bevor ein Beleg geschrieben ist; der Zähler folgt aus den
Dateien der Closure.

### Sub-Area: `*` (Default)

Alle berührten Sub-Areas GF.
