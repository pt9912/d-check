# ADR-0081: Review-Report-Deckung wird das Modul `reviews`

**Status:** Accepted

**Datum:** 2026-09-03

**Autor:** pt9912

**Bezug:** [`DC-FA-RVW-001`](../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)
(die neue Anforderung), slice-172 <!-- d-check:status-provenance --> §3 (die
ausdrückliche Abgrenzung, die diesen Slice benennt — der einsammelnde
Vorgang selbst benannte drei Slices mit offenem Review-Haken als
Anlassfall), [ADR-0072](0072-workflows-modul.md) (das strukturell nächste
Modul — Filesystem-Port, keine VCS-Bindung).

**Schärft:** [`DC-FA-RVW-001`](../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)

**Regeln:** Baseline-Regelwerk
[`modul-10-review-harness.md` §Reviewer berichtet auch, was er nicht gefunden hat](../../../.harness/baseline/v5.18.0/regelwerk/modul-10-review-harness.md#reviewer-berichtet-auch-was-er-nicht-gefunden-hat).

---

## Kontext

Der Vorgänger-Slice slice-172 <!-- d-check:status-provenance --> hat
gemessen: **87 von 95** `done/`-Slices mit Review-Zusage tragen tatsächlich
einen Report in `docs/reviews/` — die Konvention wird gelebt, aber ohne
Wächter. Der Anlassfall — drei Slices gingen mit offenem Review-Haken nach
`done/`, ohne dass ein Gate es meldete — führte dort zu einem
**offenen-Haken**-Wächter (`structure`, `max-open-tasks: 0`), der
ausdrücklich abgrenzte:

> **Keine Review-Report-Deckung.** Die Prüfung „jeder `done/`-Slice mit
> Review-Zusage hat einen Report in `docs/reviews/`" ist gemessen tragfähig
> (87/95), aber sie ist eine **Deckung zwischen zwei Mengen** — das kann keine
> `structure`-Regel. Es wäre eine neue Fähigkeit und damit eine neue
> Anforderung: eigener Slice.

Ein geschlossener Haken ist außerdem eine **schwächere** Zusage als ein
Report: Beim Scharfschalten dieses Moduls fand sich ein **geschlossener**
("`[x]`") Review-Haken ohne jeden Report unter `docs/reviews/`. Der
Haken-Wächter sieht das nicht, weil der Haken gesetzt ist; das ist genau die
oben zitierte Lücke.

## Entscheidung

Ein neues, eigenständiges Modul `reviews`
([`DC-FA-RVW-001`](../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)):

1. **Zwei Verzeichnisse, ein Modul.** `reviews.done-dir` trägt die
   Slice-Pläne mit ihrer Zusage, `reviews.reviews-dir` die tatsächlichen
   Reports — beide sind der **Aktivierungs-Schalter** (`done-dir` leer ⇒
   inert) und **nicht** verdrahtet: die Ablage ist Repo-Konvention, kein
   CI-System-Standard wie bei `workflows`.
2. **Die Zusage ist eine Zeilen-Form, keine Struktur.** Ein DoD-Haken, dessen
   Zeile die Phrase „unabhängiger Review" trägt — jede der drei
   CommonMark-Bullet-Formen, Haken-Zustand egal (dieselbe Lexik-Entscheidung
   wie `max-open-tasks`/[ADR-0074](0074-offene-tasks-auf-rohen-zeilen.md)).
   **Bloßes „Review" ist zu breit** — gemessen an einem Bestandsfall, dessen
   „Adaptions-Review" ein anderes, in der Slice-Datei selbst dokumentiertes
   Konzept ohne externen Report ist. Die engere Phrase trägt exakt die oben
   gemessene Konvention.
3. **Die Deckung ist ein Substring-Match auf die `slice-<NNN>`-Kennung** —
   dieselbe Form wie `tools/archive-wave/collect.go`s `CollectReviews`
   (1:N zulässig, z. B. `-r1`/`-r2`-Suffixe), hier unabhängig nachgebaut: das
   Modul liest kein Fremd-Werkzeug (Hexagon-Schnitt).
4. **Beide Verzeichnisse werden NICHT rekursiv gescannt.** Ein bereits
   archivierter Slice-Stub (`done/<welle-id>/`) trägt keine DoD mehr und
   fällt damit natürlich aus der Kandidatenmenge — kein Sonderfall nötig
   (AGENTS.md §3.8-Grenze, ausgesprochen im Modul-Kommentar).
5. **Fail-closed bei leerer Kandidatenmenge**, wie jedes andere Modul dieser
   Familie — **nicht** bei null gefundenen Zusagen unter vorhandenen
   Kandidaten: Ein kleiner oder junger Bestand ohne jede Review-Zusage ist ein
   legitimer Zustand, anders als bei `workflows`s `refs == 0` (eine
   Workflow-Datei ganz ohne `uses:`-Zeile ist dort ein Anomalie-Signal).
6. **Bestands-Ausnahme, feste Dateiliste** (analog
   [`MR-049`](../../../harness/conventions.md#mr-049)/[`MR-056`](../../../harness/conventions.md#mr-056),
   aber ohne Ziffernbereich — die fünf Treffer streuen nicht zusammenhängend),
   genannt in `.d-check.yml`s `reviews.exempt-paths`. Kein Nachtrag — ein
   nachträglich erzeugter Report behauptete eine Prüfung, die es nicht gab.
7. **Noch nicht Teil von `gates`.** Eine neue Modul-Klasse startet als
   eigenständiger Fokus-Lauf (`make review-coverage`), dieselbe Vorsicht wie
   bei `commits`/`trace-check`. Aufnahme in `gates` ist eine spätere, eigene
   Entscheidung — sie setzt voraus, dass die Bestands-Ausnahme nicht wächst.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, weiter manuell prüfen | kein Eingriff | genau der Zustand, den der Vorgänger-Slice als Lücke misst — 87/95 ist gelebt, nicht gewächtert |
| B — Erweiterung von `structure` (z. B. ein zweites `max-open-tasks`-Muster) | kein neues Modul | vom Vorgänger-Slice explizit ausgeschlossen: eine Deckung zweier Mengen ist keine `structure`-Regel (die prüft EIN Dokument) |
| C — Erweiterung von `tracked` (Link-Ziel gegen git-Index) | ähnliche Form (Deckung gegen einen externen Bestand) | `tracked` prüft Markdown-**Linkziele**, nicht Dateinamen-Substrings über zwei Verzeichnisse; die Umwidmung wäre unklarer als ein eigenes, kleines Modul |
| **D — eigenständiges Modul `reviews`, Filesystem-Port, Substring-Match (gewählt)** | eine Zusage, ein Ort; folgt `workflows`s bewährter Form (Dir-Aktivierungsschalter, hermetisch, `Reason...`-Consts lokal); über `d-check.mk` an Schwester-Repos verteilbar | neue Anforderung samt Spec, Codes, Tests, Doku; die Phrasen-Erkennung ist an einer Konvention geeicht, die ein anderes Repo anders schreiben könnte |

## Konsequenzen

- **Positiv, gemessen:** fünf reale Funde am Bestand, davon zwei über die drei
  bekannten Anlassfälle hinaus — beide mit **geschlossenem** Haken, also von
  `structure`s Wächter unsichtbar. Das ist der Beleg, dass die
  Deckungsprüfung mehr sieht als der Haken-Zustand.
- **Positiv:** die Prüfung ist über `d-check.mk` an Schwester-Repos
  verteilbar, ohne ein Skript zu kopieren.
- **Negativ, benannt:** die Phrasen-Erkennung („unabhängiger Review") ist
  Konvention dieses Repos, nicht Kanon — ein Adopter mit anderer Formulierung
  bekäme falsche Negative. `AGENTS.md`/Handbuch nennen die Phrase deshalb
  ausdrücklich.
- **Negativ, benannt:** einer der beiden neuen Funde ist **kein** Altfall wie
  die drei bekannten — er entstand, nachdem ein früherer Slice die
  korrigierte Praxis markiert hatte. Er bleibt trotzdem in der
  Bestands-Ausnahme stehen ([`BEO-011`](../planning/observations.md)-Familie:
  ein Slice, der eine Lücke gewächtert, kann eine **andere** Instanz
  derselben Lücke übersehen — hier: ein Haken kann `[x]` sein, ohne dass ein
  Report existiert).

## Fitness Function (falls maschinell prüfbar)

`make review-coverage` (noch **nicht** in `gates`) fährt das Modul über den
eigenen Bestand — Dogfooding. Dazu die getippten Tests in
`internal/hexagon/core/rules/reviews_test.go`:

- `TestReviewsInert` — ohne `done-dir` wird keine Datei geöffnet.
- `TestReviewsHappyPath` — Zusage plus passender Report ⇒ befundfrei.
- `TestReviewsMissing` — Zusage ohne Report ⇒ `review-missing`, auf der
  Zusage-Zeile.
- `TestReviewsBulletForms` — alle drei CommonMark-Bullet-Formen tragen die
  Zusage.
- `TestReviewsNoPromiseNoFinding` — ein Kandidat ohne Zusage ist kein Befund
  UND löst für sich allein kein Fail-Closed aus.
- `TestReviewsEmptyScopeFailsClosed` — keine Kandidaten ⇒ fail-closed.
- `TestReviewsExemptPaths` — ein ausgenommener Kandidat verschwindet aus der
  Menge, der Leerlauf-Befund bleibt.
- `TestReviewsIgnoresArchivedSubdirs` — ein archivierter Stub in
  `done/<welle-id>/` ist kein Kandidat (nicht-rekursiver Scan).

**Was keine Fitness Function prüft:** ob „unabhängiger Review" die einzige
Phrase ist, die je im Bestand eine echte Zusage trug — das ist an einem
Zeitpunkt gemessen (2026-09-03), nicht dauerhaft garantiert.

## Re-Evaluierungs-Trigger

**Ein sechster Bestands-Fund nach dem Scharfschalten.** Die Ausnahme ist eine
feste Liste zu einem Datum; wächst sie, ist entweder die Phrasen-Erkennung zu
breit, oder die Konvention „unabhängiger Review" wird nicht mehr durchgehend
befolgt — beides ein Entscheid, kein stiller Nachtrag.

**Zweiter Trigger: Aufnahme in `gates`.** Sobald ein voller Wellen-Zyklus mit
`review-coverage` grün gelaufen ist, ist die Aufnahme in `gates` fällig zu
prüfen — dieselbe Reifung wie bei `workflows` (`workflow-pins` kam erst mit
diesem Modul in `gates`, nicht mit seinem Vorgänger-Skript).

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-09-03 | Accepted. |
| 2026-09-03 | Unabhängiger Review (slice-173) fand zwei HIGH, vor der Closure behoben: die Erkennung verlangte Checkbox und Phrase auf derselben Zeile, obwohl der überwiegende Bestand DoD-Punkte über mehrere Zeilen schreibt (Fix: Item-Span bis zur nächsten Checkbox/Leerzeile/Dateiende; Neumessung liefert weiterhin exakt fünf Funde); ein unlesbares `reviews-dir` mit vorhandenen Zusagen erzeugte zusätzlich zu den Pro-Kandidat-Befunden eine textlich widersprüchliche „leere Prüfmenge"-Meldung (Fix: diese Meldung feuert nur noch, wenn sie die einzige wäre). Vier neue Tests in `internal/hexagon/core/rules/reviews_test.go`. **Policy und Entscheidung unverändert** — beide Funde sind Implementierungs-Korrekturen an der bereits getroffenen Entscheidung, keine neue Abwägung. |
