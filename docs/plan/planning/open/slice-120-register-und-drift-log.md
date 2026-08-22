# Slice slice-120: `Stand`-Zellen und Drift-Log — Zustand und Beleg statt Chronik

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-81-zustandsfelder](../welle-81-zustandsfelder.md) (zugeordnet
bei der Eröffnung).

**Bezug:** Baseline-Regelwerk `grundlagen-harness-dateien.md` §Was ein
Kommentar trägt („Dieselbe Regel für Zustandsfelder") und `modul-06-roadmap.md`
§Roadmap-Struktur / §Das Beobachtungs-Register; die Ziel-Formen der vendorten
Roadmap- und Register-Vorlage; die in slice-118 verkörperte Briefing-Regel.

**Berührte Spec-Stellen:** — (Planungs-Register; keine Spec-Zeile).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

Die beiden lebenden Register sagen, **was ist**, und belegen es mit einem
auflösbaren Anker — nicht, wie es dazu kam. Gemessen: acht `Stand`-Zellen des
Beobachtungs-Registers zwischen 169 und 3 011 Zeichen, überwiegend Chronik
(„erstes Eintreten … zweites Eintreten … gefunden vom Review …"); und ein
Drift-Log mit 71 Zeilen, dessen oberste **Schließungen** protokollieren
(„welle-79 geschlossen", „slice-N abgeschlossen") — also ein zweites
Closure-Log neben dem echten. Beides wird auf die Ziel-Form gezogen. Die
Meilenstein-Tabelle bekommt die vorgesehene Status-Form (`offen` bzw.
`erreicht <Datum>` plus Beleg-Anker).

**Was verschwindet, geht nicht verloren:** die Chronik steht in `git`, in den
Closure-Notizen der Slices und Wellen und in den Review-Reports. Nichts davon
wird nachgetragen, umgeschrieben oder anderswohin kopiert — genau das wäre die
Kopie, vor der die Regel warnt.

## 2. Vorgehen

1. **Je Register-Zeile entscheiden**, was Zustand ist (offen/verkörpert/
   gestrichen, Zähler, Gegenmittel, benannte 3×-Form) und was Chronik — die
   Belege stehen bereits als Slice-Kennungen in der Beleg-Spalte und bleiben
   der Anker.
2. **`Stand`-Zellen neu fassen:** knapp, im Indikativ über das, was ist, mit
   dem Beleg als Anker. Die 3×-Form bleibt, weil sie eine *vorgeschlagene
   Handlung* ist und beim Vorhaben gehört — sie steht beim Eintrag, nicht in
   einer Chronik.
3. **Drift-Log zurückschneiden:** jede Zeile prüfen — ist sie eine
   **Umplanung** (Trigger verschoben, präzisiert, ersetzt; Slice oder Welle
   umgehängt)? Dann bleibt sie. Ist sie eine **Schließung** oder ein
   erreichter Meilenstein? Dann entfällt sie hier; sie steht bereits im
   Closure-Log. Vorher zählen, nachher zählen, beides notieren.
4. **Meilenstein-Tabelle** auf die Status-Form ziehen (sie führt derzeit keinen
   offenen Eintrag — die Form gilt trotzdem für den nächsten).
5. **Gegenprobe:** die Abschnitts-Regeln der Prüf-Config, die auf diesen
   Tabellen liegen (Chronologie-Richtung der Drift- und Closure-Tabelle),
   laufen weiter grün — die Spalten-Lage darf sich nicht verschieben.
6. Unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Beobachtung wird gestrichen, keine Zählerstände ändern sich** — das
  ist eine Form-, keine Inhalts-Änderung. Was gestrichen gehört, entscheidet
  die Register-Lesung einer Closure, nicht dieser Slice.
- **Keine Chronik wird umgezogen.** Sie steht in `git` und in den
  Closure-Notizen; ein neuer Ablageort wäre eine Kopie.
- **Kein Produkt-Code.**

## 4. Definition of Done

- [ ] Jede `Stand`-Zelle nennt Zustand und Beleg; keine erzählt eine Chronik.
      Zählerstände und Belege unverändert (gemessen).
- [ ] Drift-Log trägt nur Umplanungen; Zeilenzahl vorher/nachher dokumentiert,
      jede entfernte Zeile ist im Closure-Log gedeckt (gemessen, nicht
      angenommen).
- [ ] Meilenstein-Status-Form übernommen.
- [ ] `make gates` grün (insbesondere die Chronologie-Regeln auf beiden
      Tabellen); unabhängiger Review; Closure-Notiz; Register gesichtet.

## 5. Abnahme-Punkte / Risiken

- **Beim Kürzen geht Substanz verloren, nicht nur Chronik.** Eine `Stand`-Zelle
  trägt heute auch das *Gegenmittel* und die *3×-Form* — beides ist Zustand
  bzw. vorgeschlagene Handlung und muss bleiben. Der Review prüft je Zeile
  gegen die Vorfassung. — **Ausgang:** *(bei Closure)*
- **Eine entfernte Drift-Log-Zeile könnte die einzige Spur sein.** Vor dem
  Entfernen wird je Zeile geprüft, ob das Closure-Log sie deckt. —
  **Ausgang:** *(bei Closure)*
- **Die Chronologie-Regeln liegen auf beiden Tabellen** — Spalten-Lage und
  Richtung dürfen sich nicht ändern. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): slice-118 in `done/`; unabhängig von
slice-119 (andere Flächen).

**Rückführungen:** `in-progress` → `next`, falls sich beim Zeilen-Zensus
zeigt, dass das Closure-Log Lücken hat (dann erst die Lücken schließen).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Planungs-Register (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-22): dieser Slice
  **bearbeitet das Register selbst** — die Einträge bleiben inhaltlich
  unberührt (§3); BEO-002 wirkt als Spiegel-Pflicht, BEO-006/009/010 als
  Arbeitsregeln.

Slice-ID: slice-120. Betroffene IDs: — (Form-Regel der Baseline). Module:
Beobachtungs-Register, Roadmap. Gates: `make doc-check` (eng), `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Form-Angleichung an die adoptierte
Baseline an eigenen Artefakten.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
