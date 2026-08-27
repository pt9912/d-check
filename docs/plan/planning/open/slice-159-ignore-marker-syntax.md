# Slice slice-159: Der `d-check:ignore`-Marker hat keine Syntax

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:**
[`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
(der Marker in seiner ersten Heimat);
[ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md) (das
geteilte Ventil); [ADR-0019](../../adr/0019-versions-pin-fence-ausnahme.md);
[`BEO-013`](../observations.md) (der Anlass);
[slice-146](../in-progress/slice-146-ignore-marker-wirkung.md) (die Messung).

**Berührte Spec-Stellen:** die Marker-Definition in
[`spec/spezifikation.md`](../../../../spec/spezifikation.md) für die vier
Konsumenten `codepaths`, `ids`, `versions`, `diagrams`.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

Der Unterdrückungs-Marker wird als **blanke Teilkette der Rohzeile** gematcht:
`strings.Contains(raw, "d-check:ignore")`. Damit unterdrückt **jede Erwähnung**
des Markers — in Prosa, in Inline-Code, in einer Tabellenzelle — genauso wie ein
gesetzter Marker. Die Dokumentation des Ventils schaltet die Prüfung ab, über
die sie schreibt.

**Gemessen** ([slice-146](../in-progress/slice-146-ignore-marker-wirkung.md)):
**233** Zeilen außerhalb von Fences tragen die Zeichenkette und unterdrücken
damit für `codepaths`, `ids`, `versions` und `diagrams`. Davon stehen nur
**87** in der HTML-Kommentar-Form; **146** sind blanke Erwähnungen. Wird die
Konstante ins Leere gelenkt, treten **58** Befunde auf **48** Zeilen hervor
(38 `id-unlinked`, 18 `codepath-missing`, 2 `repo-escape`) — die übrigen
**185** Marker-Zeilen unterdrücken nichts.

**Das ist dieselbe Klasse wie
[slice-158](../open/slice-158-citations-inline-code.md)**, nur mit umgekehrtem
Vorzeichen: dort bricht ein Modul an seiner eigenen Doku laut ab, hier schweigt
es lautlos. Die laute Variante fällt beim ersten Lauf auf; diese nicht.

## 2. Vorgehen

1. **Die Form festlegen, bevor Code entsteht.** Naheliegend ist die
   HTML-Kommentar-Form (`<!-- d-check:ignore … -->`), weil sie 87 der heutigen
   Marker bereits tragen und weil sie in gerendertem Markdown unsichtbar ist.
   Zu prüfen ist, ob das den Bestand trifft: welche der 87 stehen **nicht** am
   Zeilenanfang, welche tragen Zusatztext, und gibt es eine gesetzte
   Unterdrückung, die die neue Form **verlöre**?
2. **Den Preis zählen, nicht schätzen.** Nach der Verengung entstehen die
   Befunde wieder, die eine blanke Erwähnung heute deckt. Wie viele der 58
   hängen an einer Erwähnung statt an einem Marker? Das ist die Zahl, die die
   Umstellung kostet — und sie gehört vor der Entscheidung gemessen.
3. Trägt die Änderung: Spezifikation (die Marker-Definition je Konsument), eine
   ADR mit `Schärft:`-Feld, dann Code und Tests. Ob das Lastenheft berührt ist,
   entscheidet sich daran, ob die Marker-Form dort zugesagt ist.
4. **Die vier Konsumenten gemeinsam.** `codepaths`, `ids`, `versions` und
   `diagrams` teilen die Konstante; eine Verengung in nur einem wäre die zweite
   Antwort auf dieselbe Frage und damit ein Defekt nach
   [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md).
5. Bewusstes Brechen je Konsument, **Ursache gelesen**; Rückbau grün.
6. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Regel „dieser Marker unterdrückt nichts".** Sie ist erst sinnvoll,
  wenn ein Marker ein Marker ist — vorher meldete sie 185 Zeilen Prosa. Das ist
  [slice-146](../in-progress/slice-146-ignore-marker-wirkung.md)s Frage, und sie wartet
  auf dieses Ergebnis.
- **Keine Räumung eingefrorener Dokumente.** `done/` und `docs/reviews/`
  bleiben, wie sie sind; was dort nach der Verengung rot würde, braucht ein
  Ventil oder eine benannte Grenze.
- **Keine Änderung an den Ventil-Semantiken selbst** — welche Befunde ein
  gesetzter Marker unterdrückt, bleibt unangetastet.

## 4. Definition of Done

- [ ] Die Marker-Form ist festgelegt und gegen den Bestand geprüft — inklusive
      der Frage, ob eine gesetzte Unterdrückung verloren geht.
- [ ] Der Preis der Umstellung ist **gezählt**: wie viele der heutigen
      Unterdrückungen an einer blanken Erwähnung hängen.
- [ ] Spezifikation, ADR, Code und Tests hängen zusammen; alle vier Konsumenten
      tragen dieselbe Antwort.
- [ ] Ein konstruierter Verstoß je Konsument mit **gelesener Ursache**.
- [ ] `make gates` grün (Exit explizit), `make fullbuild` grün; unabhängiger
      Review.

## 5. Abnahme-Punkte / Risiken

- **Die Verengung deckt Befunde auf, die niemand bestellt hat.** Was heute eine
  Erwähnung deckt, wird danach rot — in lebenden Dokumenten ist das ein Gewinn,
  in eingefrorenen ein Problem ohne Adressaten. Die Menge gehört gezählt und
  entschieden, nicht erlitten. — **Ausgang:** *(bei Closure)*
- **Eine geteilte Konstante zu verengen, ändert vier Module auf einmal.** Die
  Gegenprobe muss je Konsument zeigen, dass ein **gesetzter** Marker weiter
  wirkt — sonst wird aus einem stillen Grün ein stilles Rot. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls der gezählte Preis einen
Auftraggeber-Entscheid verlangt — dann bleibt die Lücke benannt, und
[`BEO-013`](../observations.md) trägt sie weiter.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Module (GF), Spec-Straten (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27):
  [`BEO-013`](../observations.md) ist der Anlass;
  [`BEO-011`](../observations.md) — die Form gehört aus dem Bestand, nicht aus
  der einen Fundstelle, die aufgefallen ist.

Slice-ID: slice-159. Betroffene IDs:
[`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in).
Module: `codepaths`, `ids`, `versions`, `diagrams`. Gates: `make gates`,
`make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Verengung einer vorhandenen Erkennungs-Form
an vorhandenen Modulen.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
