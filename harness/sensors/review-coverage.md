# `make review-coverage` — hält, dass ein geschlossener Slice mit Review-Zusage auch einen Report hat

## Vertrag

Via Modul `reviews` (Image, dogfood). Ein `done/`-Slice mit **Review-Zusage** —
ein DoD-Haken, dessen Zeile „unabhängiger Review" nennt, jede Bullet-Form,
Haken-Zustand egal — braucht mindestens einen Report unter
`reviews.reviews-dir` mit derselben `slice-<NNN>`-Kennung im Dateinamen
(`review-missing`). Substring-Match, 1:N zulässig.

## Grenze — was das Grün nicht abdeckt

1. **Beide Verzeichnisse werden nicht rekursiv gescannt** — ein archivierter
   Slice-Stub trägt keine DoD mehr und fällt natürlich aus der
   Kandidatenmenge. Gewollt.
2. **Fail-closed bei leerer Kandidatenmenge oder unlesbarem `reviews-dir`**,
   **nicht** bei null gefundenen Zusagen unter vorhandenen Kandidaten — ein
   junger Bestand ohne jede Zusage ist legitim.
3. **Bestands-Ausnahme mit fester Dateiliste** (fünf Funde beim
   Scharfschalten, davon zwei mit geschlossenem Haken — für den Haken-Wächter
   unsichtbar).
4. **Geprüft ist die Existenz eines Reports, nicht sein Inhalt** — die
   Kategorisierung eines Findings bleibt inferential.

## Bindung

kein Gate in `gates`/`ci` — **bewusst**: eine neue Modul-Klasse startet als
eigenständiger Fokus-Lauf, dieselbe Vorsicht wie bei `trace-check`. Netzlos,
hermetisch.
[ADR-0081](../../docs/plan/adr/0081-reviews-modul.md) ·
[`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)
