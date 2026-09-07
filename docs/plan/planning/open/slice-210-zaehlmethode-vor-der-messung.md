# Slice slice-210: Die Form des Gegenstands steht vor seiner Zählung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
(3× erreicht, Ausgang *geplant*),
[`BEO-ALL/eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
(der Geschwister-Eintrag: falsche **Menge**, richtige Kategorie).

**Berührte Spec-Stellen:** — *(keine; der Slice verkörpert eine
Steering-Loop-Regel, er ändert keine Anforderung)*

**Verantwortlich:** — · **Autor:** pt9912. **Datum:** 2026-09-07.

**Form-Hinweis.** Erster Plan dieses Repos in der **Baseline-Form**
(`v6.5.0`, acht Abschnitte). Die Haus-Form ist mit
[slice-208](../done/slice-208-v650-regel-adoption.md) aufgelöst —
template-forward, also ab diesem Plan.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die Regel schreiben, die drei Vorgänge verlangt haben: **Vor einer
Messung steht die Form ihres Gegenstands** — was macht eine Kennung zu einer
Kennung, einen Verweis zu einem Folge-Verweis, eine Regel zu einer Regel? —
**und die Trefferliste wird stichprobenweise gegen diese Form gehalten, nicht
nur ihre Zahl gelesen.**

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Geschwister-Eintrag `eigene-menge-gemessen-fremde-behauptet`.** Er
  steht bei 10× und ist bereits *verkörpert*; seine Frage ist die **Menge**
  (*wer ändert, was ich zähle?*), diese hier ist die **Kategorie** (*ist das
  Gezählte überhaupt der Gegenstand?*). Zwei Fragen, zwei Regeln — eine
  gemeinsame Formulierung machte beide unschärfer.
- **Ein Sensor.** Ob jemand vor einer Messung ihre Form ausgeschrieben hat, ist
  eine Aussage über einen **Akt**, nicht über einen ruhenden Zustand — dieselbe
  Lage wie bei [`AGENTS.md`](../../../../AGENTS.md) §3.6. Es gäbe nichts zu
  bauen, und ein Gate zu behaupten wäre schlimmer als keines.
- **Ein Retrofit der drei Vorgänge.** slice-205, slice-207 und slice-208 liegen
  in `done/` und sind eingefrorene Lauf-Belege; ihre Fehlmessungen stehen in den
  Evidence-Dateien und bleiben dort.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** Der **Zielort ist entschieden und begründet**: trägt die Regel
      [`AGENTS.md`](../../../../AGENTS.md) §5 (wie die beiden
      Geschwister-Klassen), der Reviewer-Skill, oder beide? Ein Ort, der aus
      der Gewohnheit gewählt ist, ist keiner.
- [ ] **(2)** Die Regel steht am entschiedenen Ort, mit Herkunfts-Anker
      `seit slice-210` und mit der **Grenze**, was sie nicht leistet.
- [ ] **(3)** Der Registereintrag trägt den Ausgang *verkörpert* mit
      auflösbarem Zielort.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| der in DoD (1) entschiedene Zielort | update | die Regel selbst |
| [`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/state.md) | update | Ausgang von *geplant* auf *verkörpert*, mit Zielort und Anker |

**Die drei Vorgänge, und was sie gemeinsam haben.** In slice-205 wurden
`###`-Überschriften als Anforderungs-Kennungen gezählt — deutsche Komposita
zählten mit, und daraus entstand ein fünftes ID-Schema, das es nicht gibt. In
slice-207 wurden geänderte **Zeilen** gezählt und über geänderte **Regeln**
ausgesagt. In slice-208 wurden **Slice-Kennungen** gezählt und über
**Folge-Slice-Verweise** ausgesagt. Dreimal ähnelte das Muster dem Gegenstand
genug, um eine plausible Zahl zu liefern, und dreimal wurde aus den
Fehltreffern eine Aussage abgeleitet, die es nie gab.

**Dreimal von vier fand es der unabhängige Review, nicht der Zählende** — das
ist die Eigenschaft, die den Zielort mitbestimmt: Eine Regel, die nur im
Briefing steht, adressiert genau den, der sie beim Zählen nicht liest.

## 4. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei. Kein inhaltlicher
Vorgänger — die Schwelle ist mit slice-208 erreicht, der Eintrag trägt den
Ausgang *geplant* mit dieser Kennung.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Zeigt DoD (1), dass die Regel an **zwei**
  Orten stehen muss und der zweite eine eigene Form verlangt (etwa ein
  Skill-Anker mit eigener Kategorie), wird der zweite Ort ein eigener Slice.
- `in-progress` → `open` (blockiert): Ergibt DoD (1), dass die Regel von
  [`AGENTS.md`](../../../../AGENTS.md) §5 bereits mitgetragen wird — der
  Absatz über die *gemessene Menge* liegt nah —, ruht der Slice bis zum
  Entscheid, ob eine Schärfung dort die Klasse deckt oder verwässert.

## 5. Closure-Trigger

Zwei beobachtbare Kriterien und ein Lerneintrag: (a) die Regel steht am
entschiedenen Ort und `make gates` ist grün; (b) der Registereintrag trägt
*verkörpert* mit einem Zielort, der auflöst, und dem Anker `seit slice-210`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Die Regel könnte in `AGENTS.md` §5 schon stecken.** Dort steht, der Schluss
  reiche *„nicht weiter als die gemessene Menge"*. Das ist nah, aber es ist die
  Aussage über den **Schluss**, nicht über die **Messmethode** — der Unterschied
  ist genau der zwischen `commit-message-overclaims-work` und diesem Eintrag.
  Ob er trägt, ist DoD (1). — **Ausgang:** \<offen\>
- **Eine Regel aus drei Vorgängen ist keine Inventur**
  ([`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md),
  7×). Alle drei Vorgänge sind Doku-/Planning-Messungen; ob die Regel für eine
  Code-Messung trägt, ist unbelegt und gehört als Grenze in die Formulierung.
  — **Ausgang:** \<offen\>
- **Kein Gate, und das ist eine Aussage über die Wirkung.** Die Regel verschiebt
  einen Fehler von *unsichtbar* nach *vermeidbar*, nicht nach *unmöglich*; drei
  von vier Instanzen fand ein Review, und daran ändert eine geschriebene Regel
  zunächst nichts. Das gehört ausgeschrieben, sonst liest sie sich stärker, als
  sie ist. — **Ausgang:** \<offen\>

## 7. Closure-Notiz

\<wird vor dem `git mv` nach `done/` gefüllt\>

## 8. Sub-Area-Prüfungen und Modus-Begründung

\<die drei Vorprüfungen und der Modus-Block entstehen spätestens bei der
Beanspruchung — ein Plan in `open/` trägt sie noch nicht\>
