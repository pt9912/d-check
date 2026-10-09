# Eingehender Change Request — Modul `reviews` erkennt keine Review-Zusage im DoD

**Absender:** Adopter `sf-connector` · **Eingegangen:** 2026-10-09
**Richtung:** eingehend — dieses Repo ist der **Empfänger**, nicht der Bittsteller.
**Ziel-Dokument:** [`spec/lastenheft.md`](../../../spec/lastenheft.md)
**Berührt:** [`DC-FA-RVW-001`](../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)
(Modul `reviews`)
**Stand:** **angenommen, in Planung** — Punkte 1 bis 3 trägt slice-265,
Punkt 4 slice-266; beide in `open/`, ein gemeinsames Release mit slice-263.

**Ablage-Hinweis.** Ein **eingehender** CR ist die dritte Klasse neben
[`MR-035`](../../../harness/conventions.md#mr-035) (ausgehend) und
[`MR-036`](../../../harness/conventions.md#mr-036) (Antworten darauf). Was
gebeten wurde und wie entschieden wird, soll im Repo stehen und nicht nur im
Vorgang. Im Wortlaut sind Kennungen verlinkt (Linkpflicht des eigenen Gates),
sonst unverändert.

---

## Wortlaut

> **CR: Modul reviews erkennt keine Review-Zusage im DoD (d-check v0.84.0)**
>
> Kontext. In sf-connector soll das Modul reviews prüfen, dass jeder Slice in
> docs/plan/planning/done/, dessen DoD einen Review-Haken trägt, unter
> docs/reviews/ einen Report mit seiner Kennung im Dateinamen hat. Grundlage
> ist ADR-0048 im Repo. <!-- d-check:ignore (Kennung des Absender-Repos) -->
>
> Beobachtet (Pin ghcr.io/pt9912/d-check:v0.84.0, 2026-10-09, Container ohne
> Netz):
>
> - Konfiguration: done-dir: docs/plan/planning/done, reviews-dir:
>   docs/reviews, Modul in modules: und zusätzlich mit -enable reviews aktiv.
> - Auf dem Bestand meldet das Modul 0 Befunde, auch wenn reviews-dir leer ist.
> - Fehlt reviews-dir, meldet der Leerlauf-Befund „35 Kandidat(en), 0
>   Review-Zusage(n)". Das Modul findet die Slice-Dateien, erkennt aber in
>   keiner eine Zusage.
> - Getestet wurden die Zeile aus dem Bestand und elf weitere Formen in einer
>   Ein-Datei-Probe, jeweils mit und ohne Überschrift ## 2. Definition of Done.
>   Keine wurde erkannt:
>   - `- [x] Review durchgeführt, Report unter docs/reviews/ liegt vor` (Bestand)
>   - `- [ ] Review`
>   - `- [x] Review-Report liegt vor`
>   - `* [ ] Review`
>   - `- [ ] **Review**`
>   - `- [ ] Code-Review`
>   - weitere Varianten, darunter Dateinamen mit und ohne Titel nach der Nummer
>
> Folge. Das Modul ist auf diesem Repo grün über einer leeren Menge. Ein
> Slice ohne Report fällt nicht auf.
>
> Beispiel aus dem Bestand:
>
> - Slice-Datei:
>   docs/plan/planning/done/slice-039-abholzustand-und-lesepfad-publication.md.
>   Die DoD in §2 trägt die Zeile: `- [x] Review durchgeführt, Report unter
>   docs/reviews/ liegt vor (.harness/skills/reviewer.md) — Rollenwechsel nach
>   Schritt 8 …`
> - Report:
>   docs/reviews/2026-10-08-slice-039-abholzustand-und-lesepfad-publication-review.md.
>   Die Kennung steht im Dateinamen, davor das Datum, dahinter ein Suffix.
> - Kennungen haben zwei Formen: slice-&lt;NNN&gt;-&lt;titel&gt; (nummeriert)
>   und slice-&lt;welle&gt;-&lt;titel&gt; (Name).
>
> Gewünscht:
>
> 1. Erkennung dokumentieren oder konfigurierbar machen: Welche Form einer
>    DoD-Zeile gilt als Review-Zusage? Ideal wäre ein Muster in der
>    Konfiguration, zum Beispiel `promise-pattern: '^\s*- \[[ x]\] Review durchgeführt'`.
> 2. Zuordnung Report zu Slice: Ein Report gilt als vorhanden, wenn sein
>    Dateiname die volle Slice-Kennung enthält, auch mit Datums-Präfix und
>    Suffix wie -review, -korrekturen-review oder -verifikation. Am besten ist
>    auch das konfigurierbar.
> 3. Leerlauf nicht grün: Findet das Modul Kandidaten, aber keine einzige
>    Zusage, sollte das sichtbar werden, etwa als Befund oder über einen
>    Schalter require-promises: true. Grün über leerer Menge täuscht eine
>    Prüfung vor.
> 4. Handbuch im Image: Bitte die Doku der Module ins Image legen oder über
>    d-check -help &lt;modul&gt; ausgeben. Das Repo arbeitet netzlos und kann
>    die Erkennungsregel heute nur durch Probieren ermitteln.
>
> Abnahme in sf-connector: Ein roter Fall liefert review-missing für einen
> Slice mit Review-Haken ohne passenden Report, ein grüner Fall liefert 0
> Befunde. Danach heben wir den Pin in d-check.mk. Bis dahin liefert
> slice-025 nur Teil L, die Obergrenze für die Länge des Reviewer-Skills. <!-- d-check:ignore (Kennung des Absender-Repos) -->

---

## Erste Einordnung (am Code gemessen, 2026-10-09)

- **Die Erkennung ist eine feste Phrase.** Als Zusage gilt ein DoD-Punkt,
  dessen Text „unabhängiger Review" trägt; so steht es in Lastenheft und
  Spezifikation. Keine der zwölf Formen im CR enthält sie.
- **Die mitgelieferte Beschreibung ist falsch.** Die Vorlage von
  `--print-config` und das Benutzerhandbuch sagen „eine Zeile, die *Review*
  nennt" — die Regel, nach der der Absender gesucht hat. Das ist ein Fehler
  dieses Repos, nicht des Absenders.
- **Benannte Kennungen fallen still heraus.** Die Kennung wird nur als
  `slice-<Ziffern>` gelesen; eine Datei `slice-<welle>-<titel>` wird ohne
  Meldung übersprungen.
- **Punkt 2 trägt für nummerierte Kennungen schon:** Ein Report gilt als
  vorhanden, wenn sein Dateiname dieselbe `slice-<NNN>`-Kennung enthält, mit
  beliebigem Präfix und Suffix.
- **Dieselbe Blindheit hier:** Kein Volltext-Slice dieses Repos trägt die
  Phrase in einem DoD-Punkt — sie steht nur im Fließtext —, und
  `reviews.done-dir` liest keine Unterverzeichnisse; das eigene Gate prüft
  damit keinen einzigen Slice (nachgemessen in slice-265; eine erste Zählung
  hatte Fließtext-Treffer mitgezählt).

Die Entscheidung je Punkt folgt mit der Umsetzung.
