# Slice slice-120: `Stand`-Zellen und Drift-Log — Zustand und Beleg statt Chronik

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-81-zustandsfelder](welle-81-zustandsfelder.md) (zugeordnet
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
4a. **Fünfter Treffer, nachgemessen:** die `**Stand:**`-Zeile des
   Konventionsspeichers (§Baseline) ist ebenfalls ein Zustandsfeld — sie nennt
   den Zustand und den Beleg, hängt aber die Kette **aller** bisherigen
   Pin-Hebungen an. Der Kanon eröffnet allgemein („ein Feld, das einen Zustand
   trägt"), nicht auf die drei genannten Tabellen beschränkt; die Welle hatte
   vier Treffer gezählt, es sind fünf. Zustand und Beleg bleiben, die Kette
   wird auf ihren Ort gezogen — die Ableitungs-Kette steht in den
   MR-Einträgen selbst (`Löst auf:` / `Ausgelöst durch Baseline-Stand:`).
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

- [x] Jede `Stand`-Zelle nennt Zustand und Beleg; keine erzählt eine Chronik.
      Zählerstände und Belege unverändert (gemessen).
- [x] Drift-Log trägt nur Umplanungen; Zeilenzahl vorher/nachher dokumentiert,
      jede entfernte Zeile ist im Closure-Log gedeckt (gemessen, nicht
      angenommen).
- [x] Meilenstein-Status-Form übernommen.
- [x] `make gates` grün (insbesondere die Chronologie-Regeln auf beiden
      Tabellen); unabhängiger Review; Closure-Notiz; Register gesichtet.

## 5. Abnahme-Punkte / Risiken

- **Beim Kürzen geht Substanz verloren, nicht nur Chronik.** Eine `Stand`-Zelle
  trägt heute auch das *Gegenmittel* und die *3×-Form* — beides ist Zustand
  bzw. vorgeschlagene Handlung und muss bleiben. Der Review prüft je Zeile
  gegen die Vorfassung. — **Ausgang:** **eingetreten.** Der Review hat
  Substanzverlust gefunden, den ich nicht gesehen hatte: zwei Zellen
  deklarierten ihre Abweichung von der Beleg-Anzahl-Form („die Klasse ist
  dichter als der Zähler") — Zustand, nicht Chronik. Zurückgeholt. Alle
  übrigen Zellen hat er Zeile für Zeile gegen die Vorfassung geprüft; Zähler,
  Belege und Sub-Area sind byte-identisch.
- **Eine entfernte Drift-Log-Zeile könnte die einzige Spur sein.** Vor dem
  Entfernen wird je Zeile geprüft, ob das Closure-Log sie deckt. —
  **Ausgang:** entfallen — die Deckung ist mechanisch nachgezählt (jede
  genannte Welle im Closure-Log, jeder Slice in `done/`), und der Review hat
  die vier Umplanungs-Verdachtsfälle unter den entfernten Zeilen einzeln
  gegengeprüft: keine war die einzige Spur.
- **Die Chronologie-Regeln liegen auf beiden Tabellen** — Spalten-Lage und
  Richtung dürfen sich nicht ändern. — **Ausgang:** entfallen — Spalten-Lage
  und Richtung unberührt, und die Gegenprobe auf der **zurückgeschnittenen**
  Tabelle beißt unverändert (zwei vertauschte Zeilen ⇒ zwei Befunde).

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

**Geliefert:** die beiden lebenden Register sagen, **was ist**. Acht
`Stand`-Zellen tragen Zustand, Gegenmittel und — wo noch nicht verkörpert —
die mechanische Form; die Chronik ist fort (577–3 011 → 326–636 Zeichen), und
sie ist nirgendwohin kopiert worden: sie steht in `git`, in den
Closure-Notizen und in den Review-Reports. Das Drift-Log ist von **69 auf 10**
Zeilen zurückgeschnitten — es führte überwiegend Lifecycle-Protokoll und war
damit das zweite Closure-Log, vor dem die Regel warnt. Dazu die
Meilenstein-Status-Form und der fünfte Treffer, den erst der slice-118-Review
gefunden hat.

**Review** ([Report](../../../reviews/2026-08-22-slice-120-register-und-drift-log-review.md)):
merge-blockierend — 0 HIGH, 5 MEDIUM, 1 LOW, 2 INFO; die vier blockierenden
lagen **alle** auf Flächen, die dieser Slice selbst angefasst hat. Eingearbeitet.

**Was ging anders als geplant — dreimal habe ich beim Aufräumen selbst
gefehlt:**
1. **Ein Zeiger ins Leere.** Ich ersetzte die Pin-Kette durch den Satz „jeder
   Eintrag nennt seinen Vorgänger im Feld `Löst auf:`" — dieses Feld gibt es im
   ganzen Repo **einmal**. Die Kette hätte nach einem Hop geendet. Sie steht in
   Wahrheit im Index, in der Spalte *aufgelöst durch*. Wer eine Nacherzählung
   durch einen Verweis ersetzt, muss prüfen, ob das Verwiesene existiert.
2. **Substanzverlust beim Kürzen.** Zwei Zellen deklarierten, warum ihr Zähler
   und ihre Beleg-Zahl auseinanderfallen. Das ist Zustand, nicht Chronik — und
   fiel meinem Rotstift zum Opfer. Beim Kürzen ist die Frage nicht „ist das
   alt?", sondern „ist das Zustand?".
3. **Die neue Prosa widersprach ihrer eigenen Tabelle.** Ich schrieb „erreichte
   Meilensteine bleiben hier stehen" über eine leere Tabelle, aus der drei
   erreichte vor dieser Regel entfernt worden waren. Eine Regel, die man
   aufstellt, gilt auch für den Abschnitt, in dem sie steht.
Und der Trim der drei Wellen-Zeilen war halb: er erfasste Spalte 2 und ließ
Spalte 3 mit über tausend Zeichen Protokoll stehen — die Botschaft
überzeichnete den Diff.

- **Steering-Loop-Eintrag:** kein neuer Träger — die Regel liegt seit
  slice-118 im Briefing und im Reviewer-Skill; dieser Slice wendet sie an, und
  der HIGH-Anker hat in allen vier blockierenden Befunden gegriffen.
- **Beobachtungs-Register (`../observations.md`):** keine neue Beobachtung;
  Zähler und Belege sind unverändert (das war ausdrücklich Nicht-Ziel).
- **Folge-Slices:** keiner — die Welle ist mit diesem Slice vollständig. Als
  **benannte Folgepunkte** aus dem Review, ausdrücklich ohne Slice: 90 von 118
  `done/`-Slices tragen ein historisches `Status`-Kopffeld (elf davon
  widersprechen ihrem Verzeichnis) — das Briefing erklärt das historische Feld
  für Alt-Slices, die Widersprüche bleiben offen; dazu ein doppelter Anker in
  der Versions-Datei und ein ADR-Statusfeld mit Chronik, das mit der
  Immutabilitäts-Regel kollidiert.
- **Risiken aus §6:** alle drei mit Ausgang (§5) — eines eingetreten, zwei
  entfallen.
- **Drei Paarungen:** Wellen-Slice — die Paarungen prüft die Welle-Closure.
