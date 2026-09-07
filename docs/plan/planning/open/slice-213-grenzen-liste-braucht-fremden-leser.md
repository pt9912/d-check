# Slice slice-213: Braucht die Grenzen-Liste eine eigene Regel?

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md)
(3× erreicht, Ausgang *geplant* — dieser Slice ist der Ausgang),
[`BEO-ALL/rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md)
(10×).

**Berührte Spec-Stellen:** — *(keine; der Slice entscheidet über eine Regel und
ändert kein Verhalten)*

**Verantwortlich:** — · **Autor:** pt9912. **Datum:** 2026-09-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die **Vorfrage** beantworten — braucht die Klasse *„eine
Grenzen-Liste wird als vollständig gelesen"* eine eigene Regel, oder trägt
[`AGENTS.md`](../../../../AGENTS.md) §6 sie mit? §6 verlangt den Handoff an
einen Reviewer und verbietet den Self-Review; der zweite Teil des Ableiters
(*„die Liste braucht einen fremden Leser"*) könnte darin bereits stehen. **Der
erste Teil steht nirgends** — dass man nach dem Schreiben einer Grenzen-Liste
den **Vertrags**-Teil desselben Artefakts umdreht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine erneute Inventur über die Sensor-Dateien.** Sie liegt in
  [slice-212](../in-progress/slice-212-grenzen-liste-nennt-ihre-groesste-luecke.md)
  und ist dort abgeschlossen; dieser Slice urteilt über die **Regel**, nicht
  über den Bestand.
- **Ein Sensor darauf.** Ob eine Grenzen-Liste vollständig ist, ist ein Urteil
  über eine Aussage — dieselbe Lage wie bei
  [`AGENTS.md`](../../../../AGENTS.md) §3.6.
- **Der offene Bruch-Test zu `adr-check`.** slice-212 hat ihn als *nicht
  entscheidbar* stehen lassen; er braucht eine manipulierte `Accepted`-ADR im
  Arbeitsbaum und ist ein eigener Vorgang mit eigenem Risiko.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** Die **Vorfrage ist beantwortet**, und zwar aus dem **Wortlaut**
      von [`AGENTS.md`](../../../../AGENTS.md) §6, nicht aus seiner Absicht:
      Deckt er den Ableiter — beide Teile, oder nur den zweiten? Ein **Ja**
      schließt den Slice ohne neue Regel und setzt den Register-Ausgang auf
      *verkörpert* mit §6 als Zielort.
- [ ] **(2)** Fällt die Antwort auf **Nein**: Die Regel steht am entschiedenen
      Ort, mit Herkunfts-Anker `seit slice-213` und mit ihrer Grenze — **und
      diese Grenzen-Liste wird nach ihrem eigenen Ableiter geprüft**, bevor der
      Review sie sieht.
- [ ] **(3)** Der Registereintrag trägt den Ausgang mit auflösbarem Zielort.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| der in DoD (1) entschiedene Zielort | update, falls **Nein** | die Regel |
| [`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/state.md) | update | Ausgang mit Zielort |

**Die Vorfrage hat einen Präzedenzfall und eine Falle.** slice-209 stellte
dieselbe Form (*„macht der Kanon die Regel entbehrlich?"*) und beantwortete
sie mit **Nein**, belegt aus dem Wortlaut des Zielartefakts. slice-210 tat
dasselbe gegen `AGENTS.md` §5 und fand die Antwort in der **Selbstauskunft**
des Trägers. Die Falle ist beide Male dieselbe:
[`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
(18×) — §6 enger zu lesen, um die eigene Regel zu retten, oder weiter, um sie
zu sparen.

## 4. Trigger

**Start** (`open` → `in-progress`): [slice-212](../in-progress/slice-212-grenzen-liste-nennt-ihre-groesste-luecke.md)
liegt in `done/` — der Registereintrag und seine drei Belege müssen stehen,
bevor über die Regel geurteilt wird. WIP-Limit frei.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Ergibt DoD (1), dass der Ableiter **zwei**
  Orte braucht (Briefing für den Autor, Skill für den Prüfenden), und verlangt
  der zweite eine eigene Kategorie, wird er ein eigener Slice.
- `in-progress` → `open` (blockiert): Zeigt sich, dass §6 die Klasse
  **teilweise** trägt und die Differenz nur mit einer Änderung an §6 selbst
  auszudrücken wäre, ruht der Slice bis zum Entscheid — eine Hard Rule zu
  ändern ist kein Nachtrag.

## 5. Closure-Trigger

Zwei beobachtbare Kriterien und ein Lerneintrag: (a) die Vorfrage ist
beantwortet, die Antwort steht im Slice und ist aus dem **Wortlaut** des
Zielartefakts belegt; (b) der Registereintrag trägt einen Ausgang mit
auflösbarem Zielort und `make gates` ist grün.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Die bequeme Antwort ist *Ja*, und sie wäre die teuerste.** §6 zu sagen
  *„deckt das schon"* spart einen Eintrag und löscht eine Beobachtung, die
  dreimal aufgetreten ist. Der Beleg muss aus dem **Wortlaut** kommen, nicht
  aus der Absicht — sonst ist es
  [`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
  in seiner teuersten Form. — **Ausgang:** \<offen\>
- **Zwei der drei Belege sind nachgetragen.** Sie stammen aus den
  Review-Reports von slice-209 und slice-210, nicht aus einer Rekonstruktion —
  aber sie wurden geschrieben, **nachdem** die Klasse benannt war. Wer eine
  Klasse benennt und dann rückwärts Belege sammelt, findet sie. Die
  Kennzeichnung steht in den Dateien; ob die Schwelle damit **gültig** erreicht
  ist, gehört geprüft und nicht vorausgesetzt. — **Ausgang:** \<offen\>
- **Eine vierte Instanz wäre der bessere Beleg als eine dritte nachgetragene.**
  Der Eintrag ist mit diesem Zuschnitt an der Schwelle, nicht darüber. Trägt
  die Regel, entscheidet die nächste Grenzen-Liste, die jemand schreibt.
  — **Ausgang:** \<offen\>

## 7. Closure-Notiz

\<wird vor dem `git mv` nach `done/` gefüllt\>

## 8. Sub-Area-Prüfungen und Modus-Begründung

\<die drei Vorprüfungen und der Modus-Block entstehen spätestens bei der
Beanspruchung — ein Plan in `open/` trägt sie noch nicht\>
