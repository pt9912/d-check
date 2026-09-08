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
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

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
  in seiner teuersten Form. — **Ausgang:** entfallen — die bequeme Antwort ist
  nicht gegeben worden, und der unabhängige Review hat es geprüft: Die
  Ja-Hälfte ist **keine Dehnung** (ein allgemeiner Satz auf einen Fall, den
  sein Wortlaut umfasst), die Nein-Hälfte **keine künstliche Verengung** (§6
  kennt keinen Vertrags-Umkehr-Schritt). **Der Beleg kam aus dem Wortlaut**,
  wie das Risiko es verlangt hat. Das Risiko kann nicht mehr eintreten: Die
  Frage ist entschieden und die Entscheidung geprüft.
- **Zwei der drei Belege sind nachgetragen.** Sie stammen aus den
  Review-Reports von slice-209 und slice-210, nicht aus einer Rekonstruktion —
  aber sie wurden geschrieben, **nachdem** die Klasse benannt war. Wer eine
  Klasse benennt und dann rückwärts Belege sammelt, findet sie. Die
  Kennzeichnung steht in den Dateien; ob die Schwelle damit **gültig** erreicht
  ist, gehört geprüft und nicht vorausgesetzt. — **Ausgang:** entfallen — und
  zwar durch eine **vierte** Instanz, die im Lauf entstand statt nachgetragen
  zu werden: Der Slice gab der neuen Regel eine Grenzen-Aussage mit **einer**
  Grenze, und der Review fand zwei weitere, die der Slice anderswo bereits
  wusste. **Die Schwelle ruht damit nicht mehr auf zwei rückwärts gesammelten
  Belegen.** Der unabhängige Review hat zusätzlich geprüft, ob die beiden
  nachgetragenen hergeben, was ihnen zugeschrieben wird — sie tun es, gedeckt
  durch die Reports von slice-209 und slice-210.
- **Eine vierte Instanz wäre der bessere Beleg als eine dritte nachgetragene.**
  Der Eintrag ist mit diesem Zuschnitt an der Schwelle, nicht darüber. Trägt
  die Regel, entscheidet die nächste Grenzen-Liste, die jemand schreibt.
  — **Ausgang:** entfallen — die vierte Instanz ist da, und sie ist die
  aussagekräftigste der vier: Sie entstand an dem Artefakt, das die Klasse
  **benennt**, geschrieben von jemandem, der sie zu diesem Zeitpunkt kannte.
  **Damit ist die Frage beantwortet, ob die Regel trägt** — sie hat ihren
  eigenen Autor nicht davor bewahrt, und genau das ist ihre erste Grenze.

## 7. Closure-Notiz

**Geliefert.** Eine **geteilt** beantwortete Vorfrage (DoD 1), die Regel für
den Teil, den [`AGENTS.md`](../../../../AGENTS.md) §6 **nicht** trägt — in §5
mit drei Grenzen und als Anker 18 im Reviewer-Skill (1.16.0) —, und der
Registereintrag auf *verkörpert* (DoD 2/3). Ein unabhängiger Review,
blockierend, sechs MEDIUM und ein LOW. `make gates` grün (zehn Gates, 730
Dateien).

**Was funktioniert hat: die Vorfrage hat die Regel halbiert.** Ohne sie wäre
der volle zweistufige Ableiter als Regel gelandet. Gemessen am **Wortlaut**
trägt §6 den zweiten Teil bereits — *„Kein Self-Review — anderer Kontext findet
andere Findings"* —, und eine zweite Regel daneben wäre eine zweite Quelle für
dieselbe Aussage. **Geschrieben ist nur, was nirgends sonst steht.** Der Review
hat beide Hälften geprüft: keine Dehnung, keine künstliche Verengung.

**Was Friktion war, und es ist die Pointe dieses Slice: die neue Regel tat
dreimal selbst, wogegen sie geschrieben ist.**

- **Ihre Grenze war aus einer Beschreibung abgeleitet, nicht aus dem Beleg.**
  Sie sagte, den fremden Leser richte §6 ein *„und er hat in allen sieben
  Fällen gefunden"*. Die eigene Evidence-Datei sagt wörtlich das Gegenteil: Die
  Echtheits-Lücke fand der **Auftraggeber**, *„nicht ein Review und kein
  Gate"*. Aus *„nicht der Autor"* war *„§6s Leser"* geworden — und **genau
  dieser Satz begründete, Teil 2 nicht zu schreiben**.
- **Ihre Zahl hielt nicht.** *„Sieben von sieben stand die Lücke im
  Vertrags-Teil"* — bei `review-coverage` stand sie **nur im Code**, und der
  Vertrags-Text sagte das Gegenteil des Verhaltens. Sechs von sieben; die
  siebte ist der Fall, für den die **zweite** Regelhälfte existiert.
- **Ihre Grenzen-Liste war unvollständig** — zwei Grenzen, die der Slice an
  anderer Stelle bereits wusste, standen nicht in ihr.

**Steering-Loop-Lerneintrag: die vierte Instanz ist da, und sie ist die
aussagekräftigste.**
[`grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md)
stand bei 3×, zwei davon nachgetragen — §6 dieses Plans nannte das als Risiko
und eine vierte, im Lauf entstandene Instanz als den besseren Beleg. **Sie ist
eingetreten, an dem Artefakt, das die Klasse benennt, geschrieben von jemandem,
der sie kannte.** Damit ist auch die Frage beantwortet, ob die Regel trägt: Sie
hat ihren eigenen Autor nicht davor bewahrt — und genau das ist ihre erste
Grenze, die jetzt in ihr steht.

**Ein zweiter Lerneintrag, den der Review als Kandidat markiert hat.** Die
`state.md` des Registereintrags wurde **vollständig überschrieben** statt
geändert — und holte dabei die Chronik-Hälfte zurück, die der Review des
Vorgängers zwei Commits zuvor hatte entfernen lassen. Dazu der Zähler, der im
Kopf `10×` und in §8 `9×` stand. **Beides dieselbe Wurzel** — eine Stelle
geändert, ihre Geschwister nicht gelesen — und die **dritte** Meldung in drei
aufeinanderfolgenden Slices; eingetragen bei
[`semantic-change-body-only-edges-stale`](../observations/BEO-ALL/semantic-change-body-only-edges-stale/observation.md).
**Wer eine Datei ersetzt statt sie zu ändern, verliert jede fremde Korrektur
darin, ohne es zu merken.**

**Was offen bleibt.** Die Regel ist an **Sensor-Beschreibungen** belegt, nicht
an Grenzen-Listen überhaupt — das steht in ihr und im Register. Und ``semgrep``s
Regel-Cache ohne `SHA256SUMS`-Gegenstück ist weiterhin benannt, nicht
aufgelöst (Beobachtung des slice-212-Reviews).

**Die drei Paarungen, gemessen.** **(a) Anker** — `AGENTS.md` §5 und Anker 18
des Reviewer-Skills tragen `seit slice-213`, beide lösen auf, der
Registereintrag nennt beide. **(b) Folge-Slice** — keiner genannt; nichts wurde
auf später verwiesen. **(c) Register** — alle zitierten Pfade lösen auf, die
drei neuen Belege liegen als `evidence/slice-213.md` in ihren Verzeichnissen.
Der Wachposten
[`kanal-kennung-als-inhalt-gelesen`](../observations/BEO-ALL/kanal-kennung-als-inhalt-gelesen/observation.md)
trägt weiterhin kein `evidence/` — unverändert die benannte Spannung aus
slice-208.
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
