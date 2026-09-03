# Slice slice-173: Review-Report-Deckung wird ein Modul

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-86](welle-86-closure-uebergang-durchsetzen.md) — der
Closure-Trigger der Welle verlangt vier Slices in `done/`; dieser ist der
zweite von ihnen.

**Bezug:** [`DC-FA-RVW-001`](../../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)
(die neue Anforderung), [ADR-0081](../../adr/0081-reviews-modul.md) (das neue
Modul), slice-172 <!-- d-check:status-provenance --> §3 (die ausdrückliche
Abgrenzung, die diesen Slice benennt).

**Berührte Spec-Stellen:**
[`DC-FA-RVW-001`](../../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)
und seine `.a`-Verfeinerung in der Spezifikation.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-03.

---

## 1. Ziel

Ein neues, eigenständiges Modul `reviews` prüft, ob jeder `done/`-Slice mit
einer Review-Zusage tatsächlich einen Report unter `docs/reviews/` hat — die
Deckungsprüfung, die slice-172 <!-- d-check:status-provenance --> §3
ausdrücklich als eigenen Slice auswies, weil sie eine Deckung zweier Mengen
ist und keine `structure`-Regel sein kann.

## 2. Vorgehen

1. Modul-Wiring nach dem Muster von `workflows`
   ([ADR-0072](../../adr/0072-workflows-modul.md)): `ReviewsConfig` in
   `internal/hexagon/core/model/config.go`, `CheckReviews` in
   `internal/hexagon/core/rules/reviews.go`, Registrierung in `run.go`,
   `rawReviews`/`applyReviews` in `configyaml.go`.
2. Erkennung als rohe Zeilen-Form (dieselbe Lexik-Entscheidung wie
   `max-open-tasks`/[ADR-0074](../../adr/0074-offene-tasks-auf-rohen-zeilen.md)):
   eine DoD-Zeile mit Checkbox (jede der drei Bullet-Formen, Haken-Zustand
   egal), die die Phrase „unabhängiger Review" trägt.
3. Deckung als Substring-Match der `slice-<NNN>`-Kennung gegen
   `reviews.reviews-dir` — dieselbe Form wie `tools/archive-wave`s
   `CollectReviews`, unabhängig nachgebaut (kein Fremd-Werkzeug im Kern).
4. Vor dem Scharfschalten rot messen: `make review-coverage` gegen den
   echten Bestand, nicht behauptet.
5. Bestands-Ausnahme mit fester Dateiliste für die gemessenen Funde.
6. `.d-check.yml`-Config-Block, neues Makefile-Ziel `review-coverage`
   (bewusst **nicht** in `gates`), Doku (ADR, Lastenheft, Spezifikation,
   beide READMEs, `AGENTS.md` §4, `harness/README.md` Sensors-Tabelle).

**Plan (Datei-Ebene):**

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/core/model/config.go` | update | `ReviewsConfig`, `validModules()`, `Config.Reviews` |
| `internal/hexagon/core/rules/reviews.go` | neu | `CheckReviews`, `ReasonReviewMissing`, Kandidaten-/Deckungs-Logik |
| `internal/hexagon/core/rules/reviews_test.go` | neu | acht Testfälle |
| `internal/hexagon/core/rules/run.go` | update | `active["reviews"]`-Zweig |
| `internal/hexagon/core/rules/workflows.go` | refactor | `matchAnyWorkflowGlob` → `matchAnyGlob` (geteilt mit `reviews`) |
| `internal/adapter/driven/configyaml/configyaml.go` | update | `rawReviews`, `applyReviews` |
| `internal/adapter/driving/cli/config_template.go` | update | Modul-Liste, `--print-config`-Beispielblock |
| `internal/hexagon/core/app/diagnose.go` | update | `AllReasons()`, `reasonTexts()` |
| `.d-check.yml` | update | `reviews:`-Block inkl. Bestands-Ausnahme |
| `Makefile` | update | Ziel `review-coverage` |
| `docs/plan/adr/0081-reviews-modul.md` | neu | Entscheidung, Alternativen, Fitness Function |
| `docs/plan/adr/README.md` | update | Index-Zeile |
| `spec/lastenheft.md` | update | die neue Anforderung, Historie 0.84.0, Schema-Kürzel `RVW` |
| `spec/spezifikation.md` | update | ihre `.a`-Verfeinerung, ein neuer Grund-Code-Eintrag, Config-Tabelle |
| `README.md`, `README.de.md` | update | Modul-Bullet (DE zuerst) |
| `AGENTS.md` | update | §4-Tabellenzeile |
| `harness/README.md` | update | Sensors-Tabellenzeile, Governance-Gates-Zeile |

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein rekursiver Scan** beider Verzeichnisse — ein archivierter Slice-Stub
  trägt keine DoD mehr und fällt natürlich aus der Kandidatenmenge.
- **Keine Aufnahme in `gates`.** Eine neue Modul-Klasse startet als
  eigenständiger Fokus-Lauf; die Reifung ist ein späterer, eigener Schritt
  (dieselbe Vorsicht wie bei `commits`/`trace-check`).
- **Kein Nachrüsten der fünf Bestandsfälle.** Ein nachträglich erzeugter
  Report behauptete eine Prüfung, die es nicht gab.
- **Keine Aussage über Review-Qualität.** Der Sensor prüft, dass ein Report
  existiert, nicht was darin steht.

## 4. Definition of Done

- [x] `internal/hexagon/core/rules/reviews.go` implementiert `CheckReviews`
      mit acht referenzierten Tests, alle grün.
- [x] `make review-coverage` läuft gegen den echten Bestand: fünf reale
      Funde vor der Bestands-Ausnahme (gemessen, nicht behauptet), null
      danach.
- [x] [`DC-FA-RVW-001`](../../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)
      in `spec/lastenheft.md` (Historie 0.84.0), `.a`-Verfeinerung in
      `spec/spezifikation.md`, [`SPEC-081`](../../../../spec/spezifikation.md#4-grund--und-fehler-codes)
      in §4.
- [x] [ADR-0081](../../adr/0081-reviews-modul.md) trägt Kontext, Entscheidung,
      mindestens drei verglichene Alternativen, Konsequenzen,
      Fitness Function, Re-Evaluierungs-Trigger.
- [x] Doku-Update: beide READMEs (Modul-Bullet, DE zuerst), `AGENTS.md` §4,
      `harness/README.md` Sensors-Tabelle und Governance-Gates-Zeile.
- [x] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register
      fortgeschrieben (oder Fehlanzeige begründet); jedes Risiko aus §5 mit
      Ausgang; die drei Paarungen laufen bei der Welle-Closure.

## 5. Abnahme-Punkte / Risiken

- **Die Phrasen-Erkennung ist Konvention, nicht Kanon** — ein Adopter mit
  anderer Formulierung als „unabhängiger Review" bekäme falsche Negative.
  — **Ausgang:** *entfallen* — benannt in [ADR-0081](../../adr/0081-reviews-modul.md)
  „Konsequenzen" und in beiden READMEs; kein Repo-Verhalten, das behoben
  werden könnte, ohne die Konvention selbst zu ändern.
- **Ein Bestandsfund ist kein reiner Altfall wie die drei bekannten** — er
  entstand, nachdem ein früherer Slice die korrigierte DoD-Haken-Praxis
  markiert hatte, und blieb trotzdem unsichtbar, weil sein Haken geschlossen
  ist. — **Ausgang:** *weiter offen* → siehe §9; kein neuer `BEO`-Eintrag
  (Einzelfall, keine dritte Instanz einer bereits geführten Beobachtung).
- **Aufnahme in `gates` ist vertagt.** — **Ausgang:** *entfallen* für diesen
  Slice — [ADR-0081](../../adr/0081-reviews-modul.md) benennt den
  Re-Evaluierungs-Trigger (voller Wellen-Zyklus grün) explizit; das ist eine
  spätere, eigene Entscheidung, kein offenes Risiko dieses Slice.
- **Die drei Paarungen (Anker/Folge-Slice/Register) stehen noch aus** —
  **Ausgang:** *entfallen* für diesen Slice-Plan — sie laufen bei der
  Welle-Closure von welle-86 (Repo mit Wellen-Betrieb), nicht je Slice.

## 6. Trigger

**Start** (`open` → `in-progress`): direkt beansprucht — WIP-Limit frei,
`in-progress/` trug nur die Roadmap.

**Rückführungen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls die
  Doku-Nachzüge den Umfang über eine Review-Sitzung heben — trat nicht ein.
- `in-progress` → `open` (blockiert): keine Bedingung erkannt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `internal/hexagon/core/rules/` (das neue Modul),
  `spec/` (die neue Anforderung samt `.a`-Verfeinerung) und
  `docs/plan/adr/` (die neue ADR). Alle drei fallen unter den Default `*` =
  **Greenfield** ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration pro Sub-Area) — keine feinere Ausdifferenzierung nötig.
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:219-220 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-09-03, höchste
  Kennung `BEO-027`): keine Beobachtung trifft speziell die berührten Pfade
  über das hinaus, was jeder Slice ohnehin berührt — keine Treffer. Die
  Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:225-225 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): `upstream-drift.yml`
  meldet **ROT** (Lauf 2026-09-03T05:23:56Z), zwei planmäßige
  `VERALTET`-Meldungen (`go` 1.27.0→1.27.1, `semgrep` 1.175.0→1.176.0),
  `image-scan.yml` grün. Gelesen statt weggeklickt — keiner der beiden Funde
  berührt diesen Slice; eine Go-Patch-Hebung ist ein eigener, unabhängiger
  Akt. **Dieser Block trägt bewusst keine `cite`-Direktive** — sein Ziel ist
  eine Repo-Adaption, kein Baseline-Abschnitt
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-173. Betroffene IDs:
[`DC-FA-RVW-001`](../../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in),
[ADR-0081](../../adr/0081-reviews-modul.md). Module: `reviews`. Gates:
`make test`, `make review-coverage`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — alle drei berührten Sub-Areas fallen unter
den Default: Doc führt, Code folgt. Das Modul ist eine neue Fähigkeit ohne
Fremdsystem und ohne Reconciliation; die Konventions-Dichte ist hoch (das
`workflows`-Vorbild ist vollständig in `AGENTS.md`/`harness/README.md`
verankert), die Phase-Reife der berührten Abschnitte ist etabliert.

## 9. Closure-Notiz (nach `done/`)

- **Was hat funktioniert:** Das `workflows`-Muster (Dir-Aktivierungsschalter,
  hermetisch, `Reason...`-Consts lokal, Filesystem-Port) übertrug sich fast
  unverändert auf ein zwei-Verzeichnisse-Modul. Die „vor dem Scharfschalten
  rot messen"-Disziplin (slice-172s eigene Lehre) fand sofort einen echten
  Bug in der eigenen Regex (`Review` zu breit, gemessen an einem
  „Adaptions-Review"-Fund) **vor** dem Commit, nicht danach.
- **Was ging anders als geplant:** Zwei Dinge. (1) Die erste ADR-Nummer
  (0080) kollidierte mit einer zwischenzeitlich vergebenen — `ls | sort |
  tail` liest die **höchste vorhandene** Nummer, nicht „die nächste freie";
  bei mehrdeutigem Ergebnis (`uniq -c` > 1) ist ein Zwischenstand
  wahrscheinlicher als ein Zähler-Fehler. (2) Die ADR selbst verletzte beim
  ersten Schreiben `matrix-forbidden` siebzehnfach (ADR→Slice- und
  ADR→Welle-Referenzen als Markdown-Links bzw. Bare-Tokens) — behoben durch
  Bare-Token-Zitate mit `<!-- d-check:status-provenance -->`-Marker
  (ADR→Slice) bzw. vollständigem Verzicht auf die Kennung (ADR→Welle, das
  ist ein flaches Verbot ohne Marker-Ausweg).
- **Der unabhängige Review fand zwei HIGH und ein MEDIUM, alle eingearbeitet
  vor der Closure:** (1) `reviewLineRE` verlangte Checkbox und Phrase auf
  **derselben** Zeile — der überwiegende Bestand schreibt lange DoD-Punkte
  aber als Fließtext über mehrere Zeilen, und „unabhängiger Review" steht
  dabei häufig auf einer Folgezeile (gemessen: sechs Fälle im Bestand, u. a.
  `slice-138`). Behoben durch eine Item-Span-Erkennung (Checkbox-Zeile plus
  lose Folgezeilen bis zur nächsten Checkbox/Leerzeile/Dateiende); die
  Neumessung gegen den echten Bestand liefert weiterhin genau die
  **gleichen** fünf Funde — keiner der sechs neu erfassten Fälle war
  tatsächlich ungedeckt. (2) Ein unlesbares `reviews-dir` mit vorhandenen
  Zusagen erzeugte zusätzlich zu den Pro-Kandidat-Befunden eine
  textlich widersprüchliche „leere Prüfmenge"-Meldung — behoben, die
  generische Meldung feuert jetzt nur noch, wenn sie die einzige wäre. (3)
  Der `.d-check.yml`-Kommentar nannte nur „Review" statt „unabhängiger
  Review" und unterschlug damit die tragende Einschränkung — korrigiert.
  Vier neue Tests decken beide HIGH-Fälle inklusive Gegenproben.
- **Steering-Loop-Eintrag:** keiner verkörpert — alle vier Lehren dieser
  Session (zwei eigene, zwei vom Review gefunden) sind Einzelfälle, keine
  dritte Instanz einer bereits geführten Beobachtung.
- **Beobachtungs-Register (`../observations.md`):** keine Beobachtung
  angefallen, die eine neue Kennung oder einen Zähler-Schritt rechtfertigt.
- **Folge-Slices:** keine.
- **Risiken aus §5:** siehe §5, je Zeile ein Ausgang.
- **Drei Paarungen:** nicht hier geprüft — Repo mit Wellen-Betrieb, prüft die
  Closure von welle-86.
