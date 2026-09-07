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

**Verantwortlich:** pt9912 · **Autor:** pt9912. **Datum:** 2026-09-07.

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

- [x] **(1)** `harness/sensors/baseline-verify.md` §Grenze nennt die
      Echtheits-Lücke, mit der **gemessenen** Ausgabe beider Läufe und dem
      Zeiger auf den Träger, der sie hält (`--check-latest`, mit seiner
      fail-open-Bindung).
- [x] **(2)** Die **23** übrigen Sensor-Dateien sind **einmal** gegen ihr
      Skript/Target gehalten: je Datei eine Antwort — Grenze vollständig ·
      Grenze ergänzt (mit der ergänzten Lücke) · nicht entscheidbar (mit
      Begründung). Eine Datei ohne Antwort ist ein offener Punkt.
- [x] **(3)** Der Register-Eintrag zur Wiederholung ist geschrieben: dreimal in
      Folge enthielt eine ausgeschriebene Grenzen-Liste ihre eigene größte
      Lücke nicht.
- [x] `make gates` grün.
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

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

### Die Inventur (DoD 2)

**Zuerst eine eigene Fehlmessung, weil sie den Ableiter belegt.** Der erste
Überblick zählte `## Grenze`-Abschnitte über **Listenmarker** (`^- ` und
`^N.`) und meldete für `mention-coverage` **null Einträge**. Die Datei führt
ihre Grenzen als **Fettabsätze**. Gezählt war ein Proxy, ausgesagt wurde über
den Gegenstand — genau die Klasse, die
[`AGENTS.md`](../../../../AGENTS.md) §5 seit slice-210 führt, im
**Inventur-Schritt dieses Slice selbst**.

**Und die erste Fassung dieser Inventur war ebenfalls falsch — in beide
Richtungen.** Sie meldete *„20 von 24 vollständig, vier ergänzt"*. Der
unabhängige Review hat drei Stichproben gezogen und in **allen dreien** eine
ungenannte Lücke gefunden; dazu waren **zwei** der vier Ergänzungen
inhaltlich falsch, und die als *nicht entscheidbar* abgelegte Datei war am
Code entscheidbar. **Die Inventur hat also ihre eigene These nicht getragen** —
und der Grund ist derselbe, den der Registereintrag beschreibt: Wer die
Grenzen liest, kennt den Gegenstand zu gut.

**Der Stand nach der Korrektur: 17 vollständig, 7 ergänzt.**

| Datei | Antwort | Was fehlte |
|---|---|---|
| `baseline-verify` | **ergänzt** | Echtheit vs. innere Konsistenz — der Anlass (DoD 1) |
| `adr-check` | **ergänzt** | `## Geschichte` ist bis zum Dateiende ausgenommen; **79 von 84** ADRs führen sie als letzte Sektion |
| `doc-check` | **ergänzt** | **45** `ignore-refs`-Einträge und **248** `d-check:ignore`-Marker verkleinern den Prüfbereich |
| `lint` | **ergänzt** | **fünf** Ausschluss-Regeln in `.golangci.yml` |
| `arch-check` | **ergänzt** | `exclude` nimmt `**/*_test.go` und `tools/archive-wave/**` ganz heraus |
| `semgrep` | **ergänzt** | das gepinnte Regelset **altert** |
| `review-coverage` | **ergänzt** | der Abgleich sieht die **erste** Kennung im Dateinamen, und nur sie |
| die übrigen **17** | **vollständig** | — |

**Zwei Ergänzungen mussten zurückgenommen und ersetzt werden.** Sie
beschrieben Mechanismen, die es nicht gibt:

- `review-coverage`: Die erste Fassung nannte eine **Teilzeichenketten**-Suche
  und eine mögliche Präfix-Kollision. Der Code zieht die Kennung per Muster aus
  dem Namen und vergleicht auf **Gleichheit** — die Kollision ist
  ausgeschlossen. Die **echte** Grenze ist eine andere: Ein Report mit **zwei**
  Kennungen im Namen deckt nur die erste.
- `arch-check`: Die erste Fassung sagte, eine Kante ohne Regel sei
  *„unsichtbar"*. `edges` ist eine **Erlaubnis**liste — eine undeklarierte
  Kante ist ein **Befund**. Unsichtbar ist etwas anderes: eine Datei, die
  keinen `layers`-Glob trifft.

**Beide Fehlfassungen sind aus dem Vertrags-Teil abgeleitet worden statt aus
dem Code** — und das ist die Ironie dieses Slice: Sein eigener Ableiter sagt,
man solle den Vertrags-Teil umdrehen. Er setzt voraus, dass der Vertrag
**stimmt**. Bei `review-coverage` stimmte er nicht (*„Substring-Match"* steht
so im Code-Kommentar, das Verhalten ist Gleichheit). **Der Ableiter braucht
eine zweite Stufe: gegen den Code, nicht gegen die Beschreibung** — das gehört
in den Registereintrag und in die Vorfrage von slice-213.

**`adr-check` war entscheidbar, und die erste Fassung hat sich gedrückt.** Sie
schrieb *„nicht entscheidbar ohne Bruch-Test"*. Nötig waren zwei Messungen:
`exclude-sections: [Geschichte]` in der Konfiguration, und die Zählung, wie
viele ADRs `## Geschichte` als **letzte** Sektion führen (79 von 84). Kein
Bruch-Test, keine manipulierte ADR. **Die dritte Antwort der DoD ist für den
Fall da, dass etwas wirklich unentscheidbar ist — nicht dafür, dass die
Messung teuer aussieht.**

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
  ein erwarteter Ausgang. — **Ausgang:** eingetreten, aber **nicht** in der
  Richtung, die das Risiko beschrieb. Die Zahl blieb tragbar — **sieben** von
  24 Dateien brauchten eine Ergänzung, nicht zehn —, und die Rückführung war
  nicht nötig. Gesprengt hat den Zuschnitt etwas anderes: die **Gründlichkeit
  je Datei**. Die erste Fassung meldete vier Ergänzungen und *„20 von 24
  vollständig"*; der Review fand drei weitere Lücken in drei Stichproben und
  zwei **falsche** unter meinen vier. Der Aufwand lag nicht in der Zahl der
  Funde, sondern darin, jeden gegen den **Code** statt gegen die Beschreibung
  zu prüfen. **Das Risiko hat die richtige Achse verfehlt** — und dass es
  überhaupt eine Achse benannte, ist der Grund, warum die Rückführung vorab
  dastand.
- **„Vollständig" bleibt an der Kante ein Urteil.** Die Form in §3 macht den
  Kern entscheidbar (kennt das Skript die Lücke?), aber eine Lücke, die
  **niemand** bisher benannt hat, findet auch diese Inventur nicht. Sie
  verschiebt den Fehler von *unbenannt* nach *einmal geprüft*, nicht nach
  *ausgeschlossen*. — **Ausgang:** eingetreten, und der Review hat gezeigt, wie
  weit. Die Form in §3 machte den Kern entscheidbar (*kennt das Skript die
  Lücke?*) — und trotzdem gingen **drei** Lücken durch, die das Skript bzw.
  seine Konfiguration sehr wohl kannte (`.golangci.yml`, `ignore-refs`,
  `exclude`). Die Form war richtig, ihre **Anwendung** unvollständig: Ich habe
  die `## Grenze`-Abschnitte gegen meine Lektüre gehalten, nicht gegen die
  Konfigurationsdateien. **Die Verschiebung ist damit von *unbenannt* nach
  *einmal geprüft, davon drei Stichproben nachgeprüft*** — schmaler, als die
  erste Fassung klang.
- **Der Anlass kam von außen, und das ist selbst ein Befund.** Die Lücke fand
  der Auftraggeber beim Lesen eines Zwischenbescheids, nicht ein Review und
  kein Gate. Ob die zwei Vorgänger-Instanzen und diese dieselbe Klasse sind
  oder zwei, entscheidet DoD (3) — sie zu verschmelzen wäre bequem und
  vielleicht falsch. — **Ausgang:** eingetreten — und die Entscheidung fiel
  gegen das Verschmelzen. Der neue Eintrag
  [`grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md)
  grenzt sich in seinem Kopf gegen **beide** Nachbarn ab, und der unabhängige
  Review hat die Abgrenzung geprüft und als **echt** bestätigt. **Dass der
  Anlass von außen kam**, steht als Beleg in der Evidence-Datei: Gefunden hat
  die Lücke der Auftraggeber, nicht ein Review und kein Gate — bei einer
  Klasse, deren zweiter Ableiter-Teil *„die Liste braucht einen fremden Leser"*
  lautet. Der fremdeste Leser war diesmal kein Reviewer.

## 7. Closure-Notiz

**Geliefert.** Die Echtheits-Grenze in `harness/sensors/baseline-verify.md`
mit beiden gemessenen Ausgaben (DoD 1), eine Inventur über **alle 24**
Sensor-Dateien mit **sieben** Ergänzungen (DoD 2), und ein neuer
Registereintrag mit drei Belegen und dem Ausgang *geplant* → slice-213
(DoD 3). Ein unabhängiger Review, blockierend, vier MEDIUM und drei LOW.
`make gates` grün (zehn Gates, 726 Dateien).

**Was funktioniert hat: der Zuschnitt trug die Inventur statt des Anlasses.**
[`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md)
(9×) stand in slice-209 und slice-210 als Grenze im Eintrag und blieb beide
Male folgenlos. Hier hat er zum **ersten Mal den Umfang** eines Slice geändert.
Und er hat sich gelohnt: Aus **einer** gefundenen Lücke wurden **sieben** —
sechs davon hätte kein Anlass zutage gefördert.

**Was Friktion war: die Inventur trug ihre eigene These nicht.** Sie meldete
*„20 von 24 vollständig, vier ergänzt"*. Der Review zog drei Stichproben aus
der Vollständig-Menge und fand in **allen dreien** eine ungenannte Lücke; zwei
meiner vier Ergänzungen beschrieben **Mechanismen, die es nicht gibt**; und die
als *nicht entscheidbar* abgelegte Datei war mit zwei Messungen entscheidbar.

**Steering-Loop-Lerneintrag: der Ableiter braucht eine zweite Stufe, und der
Slice hat sie sich selbst beigebracht.** Der neue Eintrag sagt, man solle nach
dem Schreiben einer Grenzen-Liste den **Vertrags**-Teil desselben Artefakts
umdrehen. **Das setzt voraus, dass der Vertrag stimmt.** Bei `review-coverage`
stimmte er nicht — dort steht *„Substring-Match"* im Code-Kommentar, und das
Verhalten ist ein **Gleichheits**-Vergleich. Meine Grenze erbte den Fehler und
beschrieb eine Kollision, die es nicht geben kann. Ergänzt: **wo der Gegenstand
Code ist, wird gegen den Code geprüft, nicht gegen seine Beschreibung.**

**Und der zweite Teil des Ableiters hat sich an diesem Slice selbst bewiesen.**
Er lautet: *die Liste braucht einen fremden Leser*. Die Echtheits-Lücke fand
der **Auftraggeber** beim Lesen eines Zwischenbescheids; die drei
Stichproben-Lücken und die zwei falschen Ergänzungen fand der **Review**. **In
keinem der sieben Fälle war es der Autor** — und der Autor hatte die Klasse
zu diesem Zeitpunkt bereits benannt.

**Die dritte DoD-Antwort ist kein Ausweichgleis.** *„Nicht entscheidbar ohne
Bruch-Test"* stand für `adr-check` da, und nötig waren zwei Messungen:
`exclude-sections: [Geschichte]` in der Konfiguration und die Zählung, wie
viele ADRs die Sektion als **letzte** führen (**79 von 84** — bei ihnen liegt
der ganze Rest hinter dem Wächter). Keine manipulierte ADR, kein eigener
Vorgang. Wer *„unentscheidbar"* schreibt, weil eine Messung teuer **aussieht**,
benutzt eine ehrliche Antwort als bequeme.

**Was offen bleibt.** Die Vorfrage — braucht die Klasse eine eigene Regel, oder
trägt [`AGENTS.md`](../../../../AGENTS.md) §6 sie mit? — liegt in
[slice-213](../open/slice-213-grenzen-liste-braucht-fremden-leser.md), samt
dem unbequemen Risiko, dass **zwei der drei Belege nachgetragen** sind. Und
`semgrep`s Regel-Cache hat kein `SHA256SUMS`-Gegenstück; das ist eine
Beobachtung des Reviews, kein Befund dieses Slice, und sie ist hier **benannt,
nicht aufgelöst**.

**Die drei Paarungen, gemessen.** **(a) Anker** — vakant: Der Slice verkörpert
keine Regel; sein Eintrag steht auf *geplant*. **(b) Folge-Slice** —
[slice-213](../open/slice-213-grenzen-liste-braucht-fremden-leser.md) existiert
in `open/` und trägt eine DoD, die den Ausgang einlöst. **(c) Register** — alle
zitierten Pfade lösen auf; die neuen Belege liegen als `evidence/slice-212.md`
in ihren Verzeichnissen. Der Wachposten
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

**Eine** Sub-Area: `*` (Repo-Default). Der Slice ändert
[`harness/sensors/`](../../../../harness/sensors/)-Beschreibungen und eine
Register-Datei. **`tools/harness/` ist nicht berührt, und das ist hier eine
Aussage und keine Formalie:** §1 schließt jede Verhaltens-Änderung aus — die
Skripte werden **gelesen**, um die Grenzen zu prüfen, nicht angefasst. Die
Sub-Area der Skripte wäre `tools/harness/` (BEO-Kürzel `HARN`); sie bleibt
außen vor.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:274-274 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, **37** Verzeichnisse über beide
Kürzel). **Fünf** Einträge sind einschlägig, und der erste ist der Grund für
den Zuschnitt:

- [`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md)
  (9×, zuletzt slice-210) — **er bestimmt den Zuschnitt, statt nur als Risiko
  danebenzustehen.** In slice-209 und slice-210 stand er als Grenze im Eintrag
  und blieb folgenlos. Hier trägt DoD (2) die **Inventur** über alle 24
  Sensor-Dateien; das ist die erste Instanz, in der der Eintrag den Umfang
  eines Slice geändert hat statt nur seinen Text.
- [`wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md)
  (8×, Ausgang *geplant*) — die nächste Verwandte des Gegenstands: Dort
  behauptet ein Wortlaut eine Prüfung, die es nicht gibt; hier **verschweigt**
  eine Grenzen-Liste, was die Prüfung nicht kann. Ob das dieselbe Klasse ist
  oder eine zweite, entscheidet DoD (3) — §6 führt es als Risiko.
- [`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
  (12×, Stand *gemischt*) — dreimal in slice-211 eingetreten. Die Inventur ist
  eine Messung über 24 Dateien: Die **Menge** ist hier trivial (das
  Verzeichnis), die Gefahr sitzt in der **Antwort** je Datei.
- [`zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
  (4×, seit slice-210 *verkörpert*) — deshalb schreibt §3 aus, was
  *vollständig* heißt, **bevor** die 24 Dateien durchgegangen werden. Die
  Lehre aus slice-211 steht dabei: Die Form vorher auszuschreiben schützt
  nicht davor, sie unterwegs zu verschieben — der Abgleich am Ende gehört
  dazu.
- [`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
  (18×) — für DoD (1): Die Aussage über `baseline-verify` muss den
  **Geltungsbereich** des Skripts wiedergeben, nicht seinen Titel. Der
  Skript-Kopf sagt *„Integrität ist nicht Aktualität"*; genau dieser Satz fehlt
  in der Sensor-Datei.

**Keiner der fünf erreicht mit diesem Slice die Schwelle erstmalig.**

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-07 gelesen: **beide Nachtläufe grün** —
`upstream-drift.yml` (2026-09-07T05:33:45Z) und `image-scan.yml`
(2026-09-07T08:21:32Z). **Für diesen Slice ist der Nachtlauf mehr als
Routine:** `upstream-drift.yml` fährt `make baseline-freshness`, also genau den
Träger, der die Echtheits-Lücke dieses Slice hält. Sein Grün heißt: Der
gepinnte Baum entspricht dem Release-Asset — heute.

**Modus-Begründungsblock.** Alle berührten Sub-Areas GF — ein Block genügt.

### Sub-Area: `*`

- **Modus:** GF (Repo-Default).
- **Konventions-Dichte:** hoch für die **Form** der Sensor-Dateien (die
  vendorte `gate.template.md` gibt sie vor, seit slice-203 adoptiert), **null**
  für die Frage, wann eine Grenzen-Liste vollständig ist. Genau diese
  Asymmetrie ist der Slice.
- **Phase-Reife:** Phase 5 für die Slice-Mechanik. Phase 3 für den Gegenstand:
  Die Sensor-Dateien sind seit slice-203 in Gebrauch, ihre Grenzen-Abschnitte
  aber nie als **Menge** geprüft worden — dies ist die erste Inventur.
- **Evidenz-/Diskrepanz-Risiko:** **mittel bis hoch.** Nicht am Bestand — der
  liegt offen —, sondern in der Frage, wie viele Lücken die Inventur findet.
  §6 führt das als erstes Risiko mit vorab benannter Rückführung.
- **Reconciliation-Aufwand:** keiner (GF). Graduation entfällt.
