# `make review-coverage` — hält, dass ein geschlossener Slice mit Review-Zusage auch einen Report hat

## Vertrag

Via Modul `reviews` (Image, dogfood). Ein `done/`-Slice mit **Review-Zusage** —
ein DoD-Haken, dessen Zeile „unabhängiger Review" nennt, jede Bullet-Form,
Haken-Zustand egal — braucht mindestens einen Report unter
`reviews.reviews-dir` mit derselben `slice-<NNN>`-Kennung im Dateinamen
(`review-missing`). Substring-Match, 1:N zulässig.

## Grenze — was das Grün nicht abdeckt

1. **Beide Verzeichnisse werden nicht rekursiv gescannt** — ein archivierter
   Slice-Stub trägt keine DoD mehr und fällt aus der Kandidatenmenge; ebenso
   aber jeder Volltext unter `done/wellenlos/` und den Wellen-Verzeichnissen.
   Das Modul kann beides (`reviews.recursive`, `reviews.skip-pattern`,
   [`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in));
   die Konfiguration dieses Repos setzt es noch nicht.
2. **Die DoD-Punkte dieses Repos tragen die Default-Phrase nicht — das Gate
   prüft heute keinen Slice.** Gemessen am 2026-10-09: eine Probe-
   Konfiguration mit `recursive`, der Stub-Ausnahme und einem `reviews-dir`
   ohne Reports meldet mit der Default-Phrase **null** Befunde — die Phrase
   steht in den Volltexten nur im Fließtext, nie in einem DoD-Punkt. Mit
   `promise-pattern: "Review durchgeführt|[Uu]nabhängiger Review"` meldet sie
   26 (11 unter `done/`, 15 unter `done/wellenlos/`), mit dem echten
   `reviews-dir` null — alle 26 haben ihren Report. `require-promises` hätte
   den Leerlauf gemeldet; die Konfiguration setzt die Schlüssel, sobald sie
   released sind.
3. **Fail-closed bei leerer Kandidatenmenge oder unlesbarem `reviews-dir`**,
   **nicht** bei null gefundenen Zusagen unter vorhandenen Kandidaten — ein
   junger Bestand ohne jede Zusage ist legitim.
4. **Bestands-Ausnahme mit fester Dateiliste** (fünf Funde beim
   Scharfschalten, davon zwei mit geschlossenem Haken — für den Haken-Wächter
   unsichtbar).
5. **Geprüft ist die Existenz eines Reports, nicht sein Inhalt** — die
   Kategorisierung eines Findings bleibt inferential.
6. **Der Abgleich vergleicht die **erste** Kennung im Dateinamen, und nur
   sie** — das Modul zieht per Muster die `slice-<NNN>`-Kennung aus dem Namen
   und vergleicht sie auf **Gleichheit**. Ein Report, dessen Name **zwei**
   Kennungen trägt, deckt deshalb nur die **erste**; die zweite Zusage bleibt
   offen, ohne dass etwas meldet. **Eine Präfix-Kollision ist dagegen
   ausgeschlossen** — der Vergleich ist Gleichheit, kein Enthaltensein.

## Bindung

kein Gate in `gates`/`ci` — **bewusst**: eine neue Modul-Klasse startet als
eigenständiger Fokus-Lauf, dieselbe Vorsicht wie bei `trace-check`. Netzlos,
hermetisch.
[ADR-0081](../../docs/plan/adr/0081-reviews-modul.md) ·
[`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)
