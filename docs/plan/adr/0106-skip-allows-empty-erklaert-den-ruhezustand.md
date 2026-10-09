# ADR-0106: Der Ruhezustand nach `skip-pattern` wird erklärt, nicht still angenommen

**Status:** Accepted

**Datum:** 2026-10-09

**Autor:** pt9912

**Bezug:** [der Befund von `ai-harness-course`](../cr/2026-10-09-befund-ai-harness-course-reviews-gleichgewicht.md)
(Punkt 1, Antrag und Messung des Absenders);
[ADR-0048](0048-closure-note-struktur-im-planning-modul.md) (Entscheidung 8,
der Nullmengen-Wächter der Closure-Prüfung, dessen Reichweite hier
zurückgeschnitten wird); [ADR-0081](0081-reviews-modul.md) (Entscheidung 5,
derselbe Wächter in `reviews`); [ADR-0078](0078-erklaerte-leermenge-mit-zahl.md)
(das Muster: eine Leere, die ein generisches Muster erzeugt, wird deklariert —
hier als Schalter, nicht als Zahl, siehe Verglichene Alternativen);
`AGENTS.md` §3.6; slice-271 <!-- d-check:status-provenance -->.

**Schärft:** [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
[`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
[`DC-FA-RVW-001`](../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`planning.closure`, `structure` und `reviews` melden eine leere
Kandidatenmenge fail-closed, auch wenn erst `skip-pattern` sie geleert hat.
[ADR-0048](0048-closure-note-struktur-im-planning-modul.md) begründet den
Wächter mit einem Umzug: Wandert der Bestand in ein Unterverzeichnis, liefe
das Gate sonst leer und grün.

Ein Repo, das jeden abgeschlossenen Slice archiviert, erreicht diese Leere
aber regelmäßig, ohne dass etwas fehlt: An der Stelle jedes Slice steht ein
Stub, `skip-pattern` nimmt ihn heraus, die Menge ist leer. Der Kurs
`ai-harness-course` hat das gemessen — es gibt heute keine Konfiguration, die
dieses Repo grün lässt.

Die erste Umsetzung machte die Leere nach `skip-pattern` per Default still.
Der Review zeigte, dass das den Umzugs-Fall abschaltet: Ein flacher Stub
neben Volltexten, die in ein nicht gelesenes Verzeichnis gewandert sind,
meldete nichts mehr, wo `v0.85.0` drei Befunde meldete. Außerdem verlangt
[ADR-0078](0078-erklaerte-leermenge-mit-zahl.md) für dieselbe Art Leere — ein
generisches Muster leert die Menge — eine Deklaration, keine stille Erlaubnis.

## Entscheidung

1. Der Default bleibt fail-closed: Eine Menge, die erst `skip-pattern` leert,
   ist ein Befund, wie bisher.
2. Ein neuer Schlüssel `skip-allows-empty: true` — in `planning.closure`, je
   `structure`-Regel und in `reviews` — erklärt diese Leere zum Ruhezustand:
   Hat `skip-pattern` mindestens eine Datei genommen und bleibt keine übrig,
   gibt es keinen Befund.
3. Ohne `skip-pattern` ist der Schlüssel ein Nutzungsfehler (Exit 2), dieselbe
   halbe Aktivierung wie `exempt-expect-count` ohne `exempt-section-pattern`.
4. Fail-closed bleiben auch mit dem Schlüssel: eine Menge, aus der
   `skip-pattern` nichts genommen hat, eine Menge, die `exempt-paths` leert,
   ein unlesbares Verzeichnis.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon (Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR
(MADR)).

| Option | Pro | Contra |
|---|---|---|
| Nichts tun | Der Wächter bleibt ungeteilt | Ein archivierendes Repo hat keine Konfiguration, die grün bleibt; es müsste die Module abwählen |
| Still per Default | Keine neue Konfiguration nötig | Schaltet den Umzugs-Wächter ab — gemessen im Review; widerspricht ADR-0078 |
| Eine Zahl wie in ADR-0078 (`skip-expect-count: <n>`) | Drift ist laut: übersprungen werden muss genau die erklärte Zahl, sonst Befund | Die Zahl der Stubs wächst mit jeder Archivierung; jede Closure müsste die Konfiguration nachziehen, und ein vergessener Nachzug machte das Gate rot, ohne dass etwas fehlt — der Befund des Kurses, nur verschoben |
| Still nur mit `recursive` | Ein Umzug ins Unterverzeichnis bleibt sichtbar | `structure` kennt kein `recursive`; die Erlaubnis bleibt still statt erklärt |
| **Gewählt:** Opt-in-Schlüssel neben `skip-pattern` | Der Default bleibt streng; wer archiviert, erklärt es an einer Stelle | Mit dem Schlüssel wird ein zu breites Muster oder ein Umzug still |

## Konsequenzen

- Ein Repo, das jeden Slice archiviert, setzt `skip-pattern` **und**
  `skip-allows-empty: true` und bleibt im Ruhezustand grün.
- Ohne den Schlüssel ändert sich nichts; der Befundsatz ist byte-identisch.
- **Grenze:** Mit dem Schlüssel wird der Lauf still, wenn das Muster auch
  Volltexte trifft oder die Volltexte außerhalb der Kandidatenmenge liegen.
  Die Deklaration nimmt dieses Risiko bewusst in Kauf. Der Schalter ist damit
  **schwächer** als die Zahl aus ADR-0078: Er macht Drift nicht laut.

## Fitness Function (falls maschinell prüfbar)

Je Modul ein Test: lauter Stubs mit `skip-pattern` und ohne Schlüssel ⇒
Befund; mit Schlüssel ⇒ kein Befund; keine passende Datei ⇒ Befund auch mit
Schlüssel. Der Config-Rand weist den Schlüssel ohne `skip-pattern` mit Exit 2
ab.

## Re-Evaluierungs-Trigger

- Ein Repo mit gesetztem Schlüssel verliert einen Volltext still aus der
  Prüfung.
- Das Archiv-Werkzeug ändert die Form des Stub-Markers.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-10-09 | Proposed (Review R1 zu slice-271: die stille Fassung lockerte den Wächter ohne ADR; Auftraggeber-Entscheid für den Opt-in-Schlüssel) |
| 2026-10-09 | Proposed → Accepted (nach Review R1–R2 und Verifikation; die Fitness Function ist im Klon gebrochen, ohne Schlüssel byte-identisch zu v0.85.0) |
