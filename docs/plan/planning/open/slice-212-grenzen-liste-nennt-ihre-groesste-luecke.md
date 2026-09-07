# Slice slice-212: Die Grenzen-Liste eines Sensors nennt ihre größte Lücke

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [`MR-011`](../../../../harness/conventions.md#mr-011)-Kette
(`baseline-verify`),
[`BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md)
(8×),
[`BEO-ALL/rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md)
(9×, zweimal in Folge eingetreten — deshalb §1 die Inventur statt des Anlasses).

**Berührte Spec-Stellen:** — *(keine; der Slice korrigiert eine
Sensor-Beschreibung und ändert keine Anforderung und kein Verhalten)*

**Verantwortlich:** — · **Autor:** pt9912. **Datum:** 2026-09-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Der `## Grenze`-Abschnitt von
[`harness/sensors/baseline-verify.md`](../../../../harness/sensors/baseline-verify.md)
nennt seine **größte** Lücke — dass ein **mitgezogenes Manifest grün
passiert** —, und die übrigen **23** Sensor-Dateien werden **einmal** dagegen
gehalten, statt auf den nächsten Anlass zu warten.

**Der Anlass, gemessen:** `sha256sum -c` gegen ein `SHA256SUMS`, das mit der
Änderung nachgezogen wurde, meldet **Exit 0** (*„verify ok (54 Dateien,
vollständig)"*). Der netzlose Gate beweist **innere Konsistenz**, nicht
**Echtheit**; die Echtheit hält `--check-latest` (**Exit 4**,
`UPSTREAM-CONTENT-DRIFT`) — **Netz, fail-open, kein Gate**. Der Skript-Kopf
sagt das (*„Integrität ist nicht Aktualität"*), der `## Grenze`-Abschnitt der
Sensor-Datei nicht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`--check-latest` an ein Gate binden.** Die fail-open-Wahl ist bewusst: Ein
  Netzausfall soll den inneren Loop nicht rot machen, und `gates` ist netzlos
  ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-nah).
  Das umzustellen wäre ein **Entscheid mit ADR**, kein Doku-Nachtrag — und er
  gehört nicht in denselben Slice wie die Beschreibung des Ist-Zustands.
- **Eine Änderung an `baseline-verify` selbst.** Der Sensor tut, was er kann;
  falsch ist nur, was über ihn geschrieben steht. Wer beides in einem Slice
  anfasst, kann hinterher nicht sagen, was die Korrektur war.
- **Die `## Grenze`-Abschnitte inhaltlich neu schreiben.** Geprüft wird, ob
  eine **benannte** Lücke fehlt — nicht, ob die vorhandenen Formulierungen
  besser gingen. Ein Stil-Durchgang über 24 Dateien wäre ein anderer Vorgang.
- **Ein Sensor darauf.** Ob eine Grenzen-Liste vollständig ist, ist ein
  Urteil über eine Aussage, kein prüfbarer Zustand — dieselbe Lage wie bei
  [`AGENTS.md`](../../../../AGENTS.md) §3.6.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** `harness/sensors/baseline-verify.md` §Grenze nennt die
      Echtheits-Lücke, mit der **gemessenen** Ausgabe beider Läufe und dem
      Zeiger auf den Träger, der sie hält (`--check-latest`, mit seiner
      fail-open-Bindung).
- [ ] **(2)** Die **23** übrigen Sensor-Dateien sind **einmal** gegen ihr
      Skript/Target gehalten: je Datei eine Antwort — Grenze vollständig ·
      Grenze ergänzt (mit der ergänzten Lücke) · nicht entscheidbar (mit
      Begründung). Eine Datei ohne Antwort ist ein offener Punkt.
- [ ] **(3)** Der Register-Eintrag zur Wiederholung ist geschrieben: dreimal in
      Folge enthielt eine ausgeschriebene Grenzen-Liste ihre eigene größte
      Lücke nicht.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

**Zuerst die Form des Gegenstands** ([`AGENTS.md`](../../../../AGENTS.md) §5,
`seit slice-210`): Gezählt und beurteilt wird **je Sensor-Datei ein
`## Grenze`-Abschnitt** — alle **24** führen einen. *Vollständig* heißt
**nicht** „nennt alles Denkbare", sondern: **jede Lücke, die das Skript oder
das Target selbst kennt** — als Kommentar, als Exit-Code, als benannter
Ausgang —, steht auch im Abschnitt. Das ist entscheidbar, weil beide Seiten
lesbar sind; alles Weitergehende wäre ein Urteil und fällt unter DoD (2)s
dritte Antwort.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`harness/sensors/baseline-verify.md`](../../../../harness/sensors/baseline-verify.md) | update | der gemessene Anlass |
| die übrigen 23 unter [`harness/sensors/`](../../../../harness/sensors/) | update, wo eine Lücke fehlt | die Inventur statt des Anlasses |
| Beobachtungs-Register | update/neu | die Wiederholung |

**Warum die Inventur mit im Zuschnitt ist und nicht als Folge-Slice.** Der
Eintrag `rule-drawn-from-occasion-not-inventory` steht bei 9× und ist in
slice-209 **und** slice-210 eingetreten, beide Male mit vorab benannter
Grenze und ohne Konsequenz. Eine dritte Regel aus einem dritten Anlass wäre
die Instanz, die den Eintrag zum vierten Mal belegt. **24 Dateien sind eine
Menge, die sich in einer Sitzung lesen lässt** — die Inventur ist hier
bezahlbar, und genau das ist die Bedingung, unter der der Eintrag sie
verlangt.

## 4. Trigger

**Start** (`open` → `in-progress`): [slice-211](../done/slice-211-obermengen-nachweis-md013.md)
liegt in `done/` — WIP-Limit frei.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Ergibt die Inventur mehr als eine Handvoll
  fehlender Grenzen, ist DoD (2) kein Nachtrag mehr, sondern eigene Arbeit —
  dann wird sie ein eigener Slice und dieser schließt mit (1) und (3).
- `in-progress` → `open` (blockiert): Zeigt sich, dass eine fehlende Grenze
  nur mit einer **Verhaltens**-Änderung ehrlich zu beschreiben ist, ruht der
  Slice bis zum Entscheid — §1 schließt Verhaltens-Änderungen aus.

## 5. Closure-Trigger

Zwei beobachtbare Kriterien und ein Lerneintrag: (a) zu **jeder** der 24
Sensor-Dateien steht eine Antwort im Slice und `make gates` ist grün; (b) der
Registereintrag zur Wiederholung trägt einen Ausgang mit auflösbarem Zielort.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Die Inventur kann den Zuschnitt sprengen.** 24 Dateien sind lesbar, aber
  wenn zehn davon eine Lücke tragen, ist DoD (2) kein Nachtrag mehr. Die
  Rückführung ist deshalb vorab benannt und **nicht** die Ausnahme, sondern
  ein erwarteter Ausgang. — **Ausgang:** \<offen\>
- **„Vollständig" bleibt an der Kante ein Urteil.** Die Form in §3 macht den
  Kern entscheidbar (kennt das Skript die Lücke?), aber eine Lücke, die
  **niemand** bisher benannt hat, findet auch diese Inventur nicht. Sie
  verschiebt den Fehler von *unbenannt* nach *einmal geprüft*, nicht nach
  *ausgeschlossen*. — **Ausgang:** \<offen\>
- **Der Anlass kam von außen, und das ist selbst ein Befund.** Die Lücke fand
  der Auftraggeber beim Lesen eines Zwischenbescheids, nicht ein Review und
  kein Gate. Ob die zwei Vorgänger-Instanzen und diese dieselbe Klasse sind
  oder zwei, entscheidet DoD (3) — sie zu verschmelzen wäre bequem und
  vielleicht falsch. — **Ausgang:** \<offen\>

## 7. Closure-Notiz

\<wird vor dem `git mv` nach `done/` gefüllt\>

## 8. Sub-Area-Prüfungen und Modus-Begründung

\<die drei Vorprüfungen und der Modus-Block entstehen spätestens bei der
Beanspruchung — ein Plan in `open/` trägt sie noch nicht\>
