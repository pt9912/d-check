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
(9×).

**Berührte Spec-Stellen:** — *(keine; der Slice entscheidet über eine Regel und
ändert kein Verhalten)*

**Verantwortlich:** pt9912 · **Autor:** pt9912. **Datum:** 2026-09-07.

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
  [slice-212](../done/slice-212-grenzen-liste-nennt-ihre-groesste-luecke.md)
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

- [x] **(1)** Die **Vorfrage ist beantwortet**, und zwar aus dem **Wortlaut**
      von [`AGENTS.md`](../../../../AGENTS.md) §6, nicht aus seiner Absicht:
      Deckt er den Ableiter — beide Teile, oder nur den zweiten? Ein **Ja**
      schließt den Slice ohne neue Regel und setzt den Register-Ausgang auf
      *verkörpert* mit §6 als Zielort.
- [x] **(2)** Fällt die Antwort auf **Nein**: Die Regel steht am entschiedenen
      Ort, mit Herkunfts-Anker `seit slice-213` und mit ihrer Grenze — **und
      diese Grenzen-Liste wird nach ihrem eigenen Ableiter geprüft**, bevor der
      Review sie sieht.
- [x] **(3)** Der Registereintrag trägt den Ausgang mit auflösbarem Zielort.
- [x] `make gates` grün.
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

### Die Vorfrage, beantwortet (DoD 1)

**Die Antwort ist geteilt, und das war keine der beiden vorgesehenen.** Die DoD
sah **Ja** (§6 trägt es, kein neuer Eintrag) oder **Nein** (Regel schreiben)
vor. Gemessen am Wortlaut trägt §6 den **zweiten** Teil des Ableiters
vollständig und den **ersten** gar nicht. Geschrieben wird deshalb **nur der
erste** — eine deutlich kleinere Regel, als der Registereintrag nahelegt.

**Teil 2 — *„die Liste braucht einen fremden Leser"* — trägt §6, wörtlich.**
Der Absatz nach Schritt 8 lautet:

> **Kein Self-Review** — anderer Kontext findet andere Findings, derselbe
> Kontext dieselben blinden Flecken

Das ist die Aussage, und sie gilt **jedem** Artefakt, das der Slice
hervorbringt — eine Grenzen-Liste ist keine Ausnahme, die eigens zu benennen
wäre. **Eine zweite Regel daneben wäre eine zweite Quelle für dieselbe
Aussage**, und genau davor warnt die Source Precedence. Der Bestand stützt das:
In allen sieben Fundstellen aus slice-212 fand die Lücke jemand anderes — der
Mechanismus, den §6 einrichtet, **hat funktioniert**. Was fehlte, war nie der
fremde Leser.

**Teil 1 — *„den Vertrags-Teil umdrehen, und wo der Gegenstand Code ist, gegen
den Code prüfen"* — steht nirgends.** §6 kennt acht Schritte; der nächste
Verwandte ist Schritt 5 (*„Engsten nützlichen Sensor laufen lassen"*), und der
meint einen **Sensor**, nicht das Lesen eines Vertragstexts. §5 regelt, was
eine Aussage behaupten darf; §3.7, was ein Kommentar trägt. **Keine der drei
sagt, was man tut, bevor man eine Grenzen-Liste aus der Hand gibt.**

**Die Gegenprobe, und sie ist hier die wichtigere Hälfte.** Die bequeme
Antwort wäre **Ja** gewesen — sie spart einen Eintrag. Sie hätte §6 von *„das
Arbeitsergebnis geht an einen fremden Leser"* auf *„also braucht der Autor
vorher nichts zu tun"* gedehnt. Das ist
[`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
(18×) in seiner teuersten Form: eine dreimal belegte Beobachtung mit einem
Satz löschen, der von etwas anderem handelt. **Die Gegenrichtung ist ebenso
geprüft:** §6 enger zu lesen, um auch Teil 2 als Regel zu rechtfertigen, wäre
derselbe Fehler mit umgekehrtem Vorzeichen — deshalb fällt Teil 2 weg.

**Folge:** DoD (2) greift, aber **nur für Teil 1**. Der Registereintrag behält
seinen zweistufigen Ableiter — er beschreibt die Klasse —, die **Regel** trägt
davon nur, was nirgends sonst steht.

## 4. Trigger

**Start** (`open` → `in-progress`): [slice-212](../done/slice-212-grenzen-liste-nennt-ihre-groesste-luecke.md)
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

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende. Dieses Repo führt **drei** Prüfungen — die
zwei kanonischen und, als Adaption, den Nachtlauf-Stand
([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:268-269 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area: `*` (Repo-Default). Der Slice ändert höchstens
[`AGENTS.md`](../../../../AGENTS.md) und eine Register-Datei.
`tools/harness/` ist **nicht** berührt — §1 schließt einen Sensor aus.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:274-274 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, **38** Verzeichnisse über beide
Kürzel — der neue Eintrag ist der 38.). **Vier** sind einschlägig:

- [`grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md)
  (3×, Ausgang *geplant*) — **der Gegenstand.** Der Eintrag nennt diesen Slice
  namentlich. **Zwei seiner drei Belege sind nachgetragen**, und §6 führt das
  als Risiko: Wer eine Klasse benennt und dann rückwärts sammelt, findet sie.
- [`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
  (18×) — **der einschlägigste für DoD (1)**, und zwar in **beide** Richtungen:
  [`AGENTS.md`](../../../../AGENTS.md) §6 enger zu lesen, um eine eigene Regel
  zu rechtfertigen, ist derselbe Fehler wie ihn weiter zu lesen, um sie zu
  sparen. Der Beleg muss aus dem **Wortlaut** kommen.
- [`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md)
  (9×) — slice-212 hat gezeigt, was er wert ist: Dort änderte er den **Umfang**
  und machte aus einer Lücke sieben. Hier ist die Inventur klein und schon
  getan — die 24 Sensor-Dateien liegen vor, der Bestand ist gemessen.
- [`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
  (13×, Stand *gemischt*) — dreimal in slice-211 und einmal in slice-212
  eingetreten. Dieser Slice zählt wenig; wo er es tut, gilt der Test.

**Keiner der vier erreicht mit diesem Slice die Schwelle erstmalig.**

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-07 gelesen: **beide Nachtläufe grün** —
`upstream-drift.yml` (2026-09-07T05:33:45Z) und `image-scan.yml`
(2026-09-07T08:21:32Z). Nichts zu tun.

**Modus-Begründungsblock.** Alle berührten Sub-Areas GF — ein Block genügt.

### Sub-Area: `*`

- **Modus:** GF (Repo-Default).
- **Konventions-Dichte:** hoch, und diesmal ist sie der **Gegenstand**:
  [`AGENTS.md`](../../../../AGENTS.md) §5 und §6 regeln beide etwas, das an
  die Klasse grenzt. Die Frage ist gerade, wo die Grenze verläuft.
- **Phase-Reife:** Phase 5 für die Slice-Mechanik; Phase 4 für den Gegenstand —
  die Klasse ist benannt, dreimal belegt und in ihrem Ableiter zweistufig
  ausformuliert, aber noch nirgends als Regel geschrieben.
- **Evidenz-/Diskrepanz-Risiko:** **niedrig für den Bestand, hoch für DoD (1).**
  Am Bestand ist nichts zu inventarisieren. Das Risiko sitzt im Urteil über
  §6 — und `citation-stretched-beyond-scope` steht bei 18×.
- **Reconciliation-Aufwand:** keiner (GF). Graduation entfällt.
