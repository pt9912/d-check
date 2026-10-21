# `make review-coverage` — hält, dass ein geschlossener Slice mit Review-Zusage auch einen Report hat

## Vertrag

Via Modul `reviews` (Image, dogfood). Ein `done/`-Slice mit **Review-Zusage**
braucht mindestens einen Report unter `reviews.reviews-dir` mit derselben
`slice-<NNN>`-Kennung im Dateinamen (`review-missing`), Substring-Match, 1:N
zulässig. Was als Zusage gilt, legt
[`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)
fest; die Konfiguration dieses Repos setzt kein eigenes Muster und nutzt den
Default.

## Grenze — was das Grün nicht abdeckt

1. **Ein archivierter Stub ist kein Kandidat** — das Modul liest `done/` samt
   Unterverzeichnissen (`reviews.recursive`) und nimmt den Stub über seinen
   Marker `> **ARCHIVIERT** — Volltext:` aus (`reviews.skip-pattern`,
   [`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)).
   Ein Volltext, der diesen Marker am Zeilenanfang zitiert, fällt ebenso still
   heraus. KOPPLUNG: dieselbe Kandidatenmenge liest
   [`verify-closure-notes`](verify-closure-notes.md) am Übergang — beide
   Profile tragen dieselben drei Schlüssel (`recursive`, `skip-pattern`, `require-promises`).
2. **Ob das Gate überhaupt Zusagen sieht, zeigt kein grüner Lauf.** Die
   DoD-Punkte dieses Repos tragen die Form der Slice-Vorlage („Review
   durchgeführt"), die der Default erkennt. Gemessen (2026-10-21, an einer
   Wegwerf-Kopie ohne Reports): 42 erkannte Zusagen, also 42 `review-missing`.
   Ändert sich der Wortlaut der Vorlage, so dass das Muster keine einzige Zeile
   mehr trifft, bleibt der Lauf ohne `require-promises` grün; mit dem Schalter,
   den dieses Repo seit slice-278 in beiden Profilen setzt, meldet er „keine
   Review-Zusage unter N Kandidat(en)". **Grenze:** der Schalter fängt den
   **Totalausfall**, nicht den einzelnen Slice — ein Slice mit anderer Formulierung
   neben anderen, die das Muster treffen, fällt weiter still heraus. Zeigen lässt
   sich die Zahl mit einer Probe-Konfiguration, deren `reviews-dir` auf ein
   Verzeichnis ohne Reports zeigt: die Zahl der Befunde ist die Zahl der
   erkannten Zusagen.
3. **Fail-closed bei leerer Kandidatenmenge oder unlesbarem `reviews-dir`.**
   Null gefundene Zusagen unter vorhandenen Kandidaten sind hier ein Befund
   (`require-promises`); ohne den Schalter wären sie für einen jungen Bestand
   legitim. Auch eine Menge, die erst `skip-pattern` leert, bleibt ein Befund
   („leere Prüfmenge … fail-closed", gemessen): Dieses Repo setzt
   `reviews.skip-allows-empty` nicht — heute liegen 47 lebende Slices als Volltext
   in `done/`, und der Schlüssel nähme dem Modul die Fail-closed-Antwort, wenn ein
   falsches `skip-pattern` alles ausnimmt. **Grenze:** `require-promises` zählt die
   Zusagen **nach** dem `skip-pattern`. Heute tragen 43 der 47 lebenden Slices eine
   Zusage; archiviert eine Wellen-Closure alle, die eine tragen, und lässt nur
   Slices ohne Zusage liegen, melden `make review-coverage` und
   `make verify-closure-notes` „keine Review-Zusage unter N Kandidat(en)" — auch ohne
   dass etwas fehlt, und `skip-allows-empty` greift dort nicht.
4. **Geprüft ist die Existenz eines Reports, nicht sein Inhalt** — die
   Kategorisierung eines Findings bleibt inferential.
5. **Der Abgleich vergleicht die **erste** Kennung im Dateinamen, und nur
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
[ADR-0105](../../docs/plan/adr/0105-reviews-liest-done-unterverzeichnisse.md) ·
[`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)
