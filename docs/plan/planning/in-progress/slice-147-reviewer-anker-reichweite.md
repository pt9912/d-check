# Slice slice-147: Ein Anker gegen Zitat-Reichweite — Feedforward, nicht nur Feedback

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`BEO-012`](../observations.md) (Zähler 4); Baseline-Regelwerk
[`grundlagen-klassifikation.md` §Die 2×2-Matrix](../../../../.harness/baseline/v5.12.0/regelwerk/grundlagen-klassifikation.md#2x2-matrix)
(Feedforward ↔ Feedback); `.harness/skills/reviewer.md` §Repo-spezifische Anker
pro Kategorie.

**Berührte Spec-Stellen:** — (Harness-Regeltext; keine Anforderung des Produkts).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-25.

---

## 1. Ziel

[`BEO-012`](../observations.md) steht bei Zähler **4** — eine Quelle wird über
ihren Geltungsbereich hinaus zitiert. **Alle vier Male fand es der zweite
Leser.** Das ist die entscheidende Eigenschaft der Klasse: sie ist im Review
zuverlässig zu finden und beim Schreiben zuverlässig zu übersehen, weil ein
Zitat wie ein Beleg aussieht.

Der Reviewer-Skill trägt für sie **keinen** Anker, während er für
[`BEO-009`](../observations.md) und [`BEO-004`](../observations.md) je einen
MEDIUM-Eintrag mit Arbeitsanweisung führt. Die vierte Instanz zeigt zugleich,
dass Feedback allein nicht reicht: sie entstand in einem Dokument, dessen
Gegenstand diese Klasse **ist**.

## 2. Vorgehen

1. **Zuerst die Kategorie entscheiden**, nicht den Wortlaut. `BEO-009`(b) und
   `BEO-004` stehen als MEDIUM; ob Reichweite dieselbe Stufe trägt, ist ein
   Urteil und gehört begründet.
2. Anker im Reviewer-Skill mit **Prüffrage an den Diff**, nicht mit einer
   Fallliste — die Fallliste ist bei `BEO-004` dreimal unvollständig gewesen.
3. **Die Feedforward-Hälfte prüfen:** trägt `AGENTS.md` §5 die Regel schon, oder
   ist der Skill wieder ihr einziger Ort? Beides ist Pflicht
   (Baseline `grundlagen-klassifikation.md`, Feedforward **und** Feedback);
   eine Waise im Skill ist
   genau der Befund, den [slice-131](../done/slice-131-reviewer-skill-waisen.md)
   behandelt hat.
4. **Messen statt annehmen**, ob eine mechanische Hälfte existiert: die
   Wortlaut-Probe aller Blockzitate eines Dokuments gegen die zitierte Quelle
   ist in [slice-140](../done/slice-140-konsumenten-cr.md) einmal von Hand
   gelaufen und hat einen Befund gebracht, den beide Leser übersahen. Sie findet
   die **Tilgung ohne Auslassungszeichen** — nicht die überzogene Reichweite.
   Was sie deckt und was nicht, gehört gemessen, bevor jemand ein Target baut.
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Gate auf Reichweite.** Ob ein Satz weiter trägt als seine Quelle, ist
  ein Urteil, kein `grep` — die Registerzeile sagt das seit ihrer Anlage.
- **Kein Nachziehen der vier Belege.** Sie liegen in `done/` und sind Protokoll.
- **Keine Ableitung aus dem Zähler.** Die Registerzeile hält ausdrücklich fest,
  dass die Entscheidung über die Form eigens zu treffen ist; ein Zähler von 4
  ist ein Anlass, kein Beschluss.

## 4. Definition of Done

- [ ] Die Kategorie ist **begründet** gewählt, nicht von `BEO-009` abgeschrieben.
- [ ] Der Anker steht als **Prüffrage**, nicht als Fallliste.
- [ ] Beide Quadranten belegt: Feedforward-Ort benannt, Feedback-Ort benannt —
      oder ausdrücklich ausgewiesen, warum einer entfällt.
- [ ] Die Deckung der Wortlaut-Probe ist **gemessen** ausgewiesen: was sie
      findet und was sie nicht findet.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein Anker mehr im Skill ist kein Schutz, sondern Kontext-Last.** Der Skill
  trägt bereits mehrere HIGH- und MEDIUM-Anker; der fünfzehnte wird nicht mehr
  gelesen. Zu prüfen ist, ob dieser Anker einen bestehenden **schärft**, statt
  neben ihn zu treten. — **Ausgang:** *(bei Closure)*
- **Die Klasse ist beim Schreiben blind — auch mit Anker.** Der Anker wirkt im
  Review, also erst nachdem der Fehler geschrieben ist. Ob die
  Feedforward-Hälfte überhaupt greifen kann, ist die eigentliche Frage dieses
  Slice und darf nicht im Anker untergehen. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Kategorie-Frage eine
Regelwerk-Klärung braucht — dann hängt sie am ausstehenden Konsumenten-CR
([slice-140](../done/slice-140-konsumenten-cr.md), Punkt 4).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Regeltext (GF), Reviewer-Skill (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-25):
  [`BEO-012`](../observations.md) ist der Anlass;
  [`BEO-011`](../observations.md) für jede Aussage darüber, was der Anker
  „immer" fängt.

Slice-ID: slice-147. Betroffene IDs: — (kein `DC-`/`ADR-`-Bezug). Module:
Harness-Regeltext. Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Regeltext-Arbeit am eigenen Skill.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
