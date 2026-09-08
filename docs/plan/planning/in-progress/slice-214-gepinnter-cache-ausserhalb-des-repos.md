# Slice slice-214: Was einen gepinnten Cache außerhalb des Repos prüft — und wann

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md) (hermetisches
semgrep-Gate), [ADR-0011](../../adr/0011-digest-pins-build-gate-images.md)
(Digest-Pins), [`MR-011`](../../../../harness/conventions.md#mr-011)-Kette
(`baseline-verify` — die **andere** Hälfte derselben Frage),
[`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md)
(4×, verkörpert — der Anlass kam aus seinem Review).

**Berührte Spec-Stellen:** — *(keine; der Slice beschreibt eine bestehende
Eigenschaft und ändert kein Verhalten)*

**Verantwortlich:** pt9912 · **Autor:** pt9912. **Datum:** 2026-09-08.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Für **jedes gepinnte Fremd-Artefakt** beantworten, **was seine
Unversehrtheit prüft und zu welchem Zeitpunkt** — und die Antwort dort
hinschreiben, wo sie fehlt. Der Anlass ist der `semgrep`-Regel-Cache: Er liegt
als **einziger** außerhalb des Repos, und die Frage nach ihm ist im Review von
slice-212 gestellt und nicht beantwortet worden.

**Die Achse ist neu und die Antwort vermutlich nicht überall dieselbe.**
`baseline-verify` prüft **jeden Lauf** und kann die **Echtheit** nicht beweisen
(slice-212). Beim `semgrep`-Cache ist es genau umgekehrt: Der Bezug über einen
**git-Commit-Pin** bindet die Echtheit stärker als jedes Manifest — ein
Commit-SHA ist ein Hash über den Baum —, aber **nach** dem Holen prüft ihn
nichts mehr; die Bedingung ist `[ ! -d "$RULES_DIR/$RULES_SUBSET" ]`, also
Existenz eines Verzeichnisses. **Zwei Artefakte, zwei entgegengesetzte
Lücken** — das ist der Grund für diesen Slice und nicht nur für eine Zeile.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Re-Verifikation einbauen.** Ob der `semgrep`-Cache bei jedem Lauf
  gegen seinen Commit-Pin geprüft werden soll, ist ein **Entscheid mit ADR**
  (Laufzeit gegen Sicherheit, und [ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md)
  nennt den einmaligen Bezug ausdrücklich als Eigenschaft). Dieser Slice
  beschreibt den Ist-Zustand.
- **Der `ignore-refs`-/`exempt-paths`-CR.** Er liegt unentschieden in
  [`docs/plan/cr/`](../../cr/); dieses Repo ist dort der **Empfänger** und
  schuldet den Entscheid, nicht der Absender eine Antwort —
  ein anderer Vorgang.
- **Eine erneute Inventur der `## Grenze`-Abschnitte.** Sie liegt in
  [slice-212](../done/slice-212-grenzen-liste-nennt-ihre-groesste-luecke.md);
  hier wird **eine** Achse über **wenige** Artefakte gemessen, nicht der
  ganze Bestand noch einmal.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [x] **(1)** Die **Inventur** liegt vor: je gepinntem Fremd-Artefakt **Ort**
      (im Repo / außerhalb), **Bindung** (Digest · Commit-SHA · Manifest),
      **Prüfzeitpunkt** (jeder Lauf · einmalig · nie) — aus der Konfiguration
      bzw. dem Skript gelesen, nicht aus der Prosa darüber
      ([`AGENTS.md`](../../../../AGENTS.md) §5, `seit slice-213`).
- [x] **(2)** Wo die Antwort in der Sensor-Beschreibung fehlt, steht sie dort —
      mit dem Zeitpunkt, nicht nur mit der Bindung.
- [x] `make gates` grün.
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

**Zuerst die Form des Gegenstands** ([`AGENTS.md`](../../../../AGENTS.md) §5,
`seit slice-210`), denn DoD (1) ist eine Messung.

**Was zählt als *gepinntes Fremd-Artefakt*?** Etwas, das (a) **nicht** in
diesem Repo entsteht, (b) an einer **festen Kennung** bezogen wird
(Digest, Commit-SHA, Tag+Manifest), und (c) in einen **Lauf** eingeht — nicht
jede Abhängigkeit überhaupt. Go-Module fallen damit heraus: Sie sind über
`go.sum` gebunden und werden von der Toolchain bei **jedem** Bau geprüft; sie
sind kein offener Punkt, sondern die Referenz-Antwort.

**Was zählt als *Prüfzeitpunkt*?** Der Moment, in dem etwas die Bindung
**tatsächlich nachrechnet** — nicht der, in dem sie in einer Datei steht. Ein
Digest im `Dockerfile` ist eine Deklaration; nachgerechnet wird sie von Docker
beim **Pull**, und beim Lauf aus dem lokalen Store.

**Der vermutete Bestand ist klein** — vier Klassen, und die Zahl gehört
geprüft statt vorausgesetzt: vier digest-gepinnte `FROM`-Zeilen im
[`Dockerfile`](../../../../Dockerfile), das digest-gepinnte a-check-Image in
[`a-check.mk`](../../../../a-check.mk), das digest-gepinnte semgrep-Image plus
sein **git-gepinntes Regelset**, und der vendorte Baseline-Baum. **Nur eines
davon liegt außerhalb des Repos.**

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`harness/sensors/semgrep.md`](../../../../harness/sensors/semgrep.md) | update | der Anlass: Bindung stark, Prüfzeitpunkt einmalig |
| weitere Sensor-Dateien | update, wo die Antwort fehlt | Ergebnis von DoD (1) |

### Die Inventur (DoD 1)

**Die erste Fassung war unvollständig und in ihrer Kernaussage falsch.** Sie
zählte **sieben** Artefakte, fand an drei **Fundorten** statt an der in §3
definierten **Form**, und schloss daraus, das `semgrep`-Regelset sei *„das
einzige gepinnte Fremd-Artefakt außerhalb des Repos"*. Der unabhängige Review
hat fünf fehlende Mitglieder gefunden und den Schluss widerlegt. Beides ist
unten korrigiert; die Fehlfassung steht in der Closure-Notiz als Beleg.

**Gelesen aus Konfiguration und Skript**, nicht aus der Prosa darüber
([`AGENTS.md`](../../../../AGENTS.md) §5, `seit slice-213`). **Zwölf**
Artefakte erfüllen die Form aus §3; die letzte Zeile ist die Referenz-Antwort
und zählt nicht mit.

| Artefakt | Ort | Bindung | Prüfzeitpunkt |
|---|---|---|---|
| `golang`-Basis (`deps`) | Docker-Store | Digest im [`Dockerfile`](../../../../Dockerfile) | **beim Bezug** |
| `golangci-lint`-Image | Docker-Store | Digest im `Dockerfile` | **beim Bezug** |
| `distroless`-Runtime | Docker-Store | Digest im `Dockerfile` | **beim Bezug** |
| `golang`-Basis des Werkzeugs | Docker-Store | **anderer** Digest in [`tools/archive-wave/Dockerfile`](../../../../tools/archive-wave/Dockerfile) | **beim Bezug** |
| `a-check`-Image | Docker-Store | Digest in [`a-check.mk`](../../../../a-check.mk) | **beim Bezug** |
| `semgrep`-Image | Docker-Store | Digest in [`tools/semgrep.sh`](../../../../tools/semgrep.sh) | **beim Bezug** |
| `trivy`-Image | Docker-Store | Digest in [`tools/image-scan.sh`](../../../../tools/image-scan.sh) | **beim Bezug** |
| drei GitHub-Actions | Runner des Anbieters | **Commit-SHA** im `uses:` ([`AGENTS.md`](../../../../AGENTS.md) §3.9) | **jeder CI-Lauf** — der Runner hat keinen Cache |
| **`semgrep`-Regelset** | **Nutzer-Cache**, außerhalb von Repo **und** Docker-Store | **git-Commit-SHA** | **beim Bezug** — lokal einmalig, in CI jedes Mal |
| `# syntax=docker/dockerfile:1.7` | Docker-Store | **keine** — beweglicher Tag | **beim Bezug**, ohne Bindung |
| vendorte Baseline | **im Repo** | `SHA256SUMS` (kommt mit dem Baum) | **jeder `make gates`** — Echtheit nur im Nachtlauf |
| *(Referenz: Modul-Abhängigkeiten)* | *Modul-Cache* | *Prüfsummen-Datei* | *jeder Bau, über den read-only-Schalter* |

**Die Achse der ersten Fassung ist zusammengebrochen, und die richtige ist
schärfer.** Sie lautete *„stark gebunden, selten geprüft"* gegen *„schwach
gebunden, oft geprüft"* und setzte voraus, dass die Digest-Images bei **jedem
Bau** nachgerechnet werden. **Gemessen im eigenen `gates`-Lauf: kein Pull, 22
Schichten aus dem Cache.** Ein Digest wird **beim Bezug** eingelöst, nicht bei
jeder Verwendung — genau wie der git-Commit-Pin des Regelsets.

**Damit steht die eigentliche Aussage:** Von zwölf gepinnten Fremd-Artefakten
wird **genau eines bei jedem Lauf erneut geprüft** — die vendorte Baseline, und
ausgerechnet ihre Bindung kann die Echtheit nicht beweisen (slice-212). **Alle
übrigen werden einmal beim Bezug geprüft und danach geglaubt.** Das ist keine
Eigenheit des `semgrep`-Caches, sondern der Normalfall; die Ausnahme ist das
Repo-interne Artefakt.

**Was das Regelset trotzdem unterscheidet — schmaler als behauptet.** Es liegt
als einziges weder im Repo **noch** im Docker-Store, sondern im Nutzer-Cache:
ein Verzeichnis, das kein Werkzeug verwaltet und dessen Inhalt nach dem Holen
von nichts mehr adressiert wird. Der Docker-Store ist inhaltsadressiert, das
Cache-Verzeichnis ist es nicht. **Ob der Docker-Store spätere Veränderung
bemerkt, ist hier nicht gemessen** und wird nicht behauptet.

**Was gemessen ist und was nicht — die Grenze der Inventur.** Gemessen sind
**Ort und Bindung** (aus den Dateien gelesen), der Prüfzeitpunkt **dort, wo ein
eigenes Skript ihn setzt** (Regelset, vendorte Baseline), und die
**Cache-Wirkung** im eigenen Lauf. **Nicht gemessen** und als Fremd-Eigenschaft
gekennzeichnet: dass Docker beim Pull nachrechnet, dass git einen
Commit-Bezug gegen den SHA prüft, und dass die Modul-Prüfsummen bei jedem Bau
greifen. **Auch die git-Aussage trägt diesen Vorbehalt** — die erste Fassung
gab ihn nur den Docker-Zeilen und stellte die eigene Lieblings-These
ungeprüft daneben.

## 4. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.
[slice-213](../done/slice-213-grenzen-liste-braucht-fremden-leser.md) liegt in
`done/` — die Regel, nach der DoD (1) misst, ist damit geschrieben.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Ergibt die Inventur, dass mehrere
  Artefakte eine **Verhaltens**-Änderung bräuchten statt einer Beschreibung,
  ist das ein ADR-Vorgang und dieser Slice zu klein geschnitten.
- `in-progress` → `open` (blockiert): Zeigt sich, dass der Prüfzeitpunkt eines
  Artefakts von hier aus **nicht feststellbar** ist (fremde Werkzeug-Innerei),
  ruht der Slice bis zum Entscheid, ob eine Messung oder eine benannte Grenze
  die Antwort ist.

## 5. Closure-Trigger

Zwei beobachtbare Kriterien und ein Lerneintrag: (a) zu **jedem** Artefakt der
Inventur stehen Ort, Bindung und Prüfzeitpunkt im Slice, jeweils aus
Konfiguration oder Skript gelesen; (b) `make gates` ist grün und die
Sensor-Beschreibungen tragen, was ihnen fehlte.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Die Antwort könnte lauten: alles in Ordnung, nur unbeschrieben.** Dann ist
  der Slice eine Doku-Ergänzung und kein Fund — und das wäre ein **gutes**
  Ergebnis, das als solches dastehen muss. Wer eine Lücke sucht, findet eine;
  die Inventur ist gegen die Konfiguration zu führen, nicht gegen die
  Erwartung. — **Ausgang:** eingetreten, und zwar in der **anderen** Richtung
  als befürchtet. Das Risiko warnte, die Antwort könne *„alles in Ordnung, nur
  unbeschrieben"* lauten und der Slice dann kein Fund sein. Gefunden wurde
  stattdessen etwas Größeres — **und nicht vom Slice, sondern vom Review**: Die
  Inventur war um fünf Mitglieder zu klein, und ihre These fiel. **Der Satz
  *„wer eine Lücke sucht, findet eine"* traf zu, nur umgekehrt:** Ich habe die
  Lücke gefunden, die ich erwartete (der Cache ist ungeprüft), und dabei
  übersehen, dass sie der Normalfall ist.
- **Der Prüfzeitpunkt fremder Werkzeuge ist von hier aus schwer zu belegen.**
  Dass Docker einen Digest beim Pull nachrechnet, ist bekannt, aber nicht in
  diesem Repo gemessen; dass `go.sum` bei jedem Bau greift, ebenso. **Wo der
  Beleg fehlt, gehört das gesagt statt behauptet** — sonst ist die Inventur
  eine Aufzählung von Vermutungen mit Tabellen-Rahmen. — **Ausgang:**
  eingetreten, und der Review hat es an der schärfsten Stelle gemessen. Die
  Tabelle schrieb *„jeder Image-Bau"*, wo **kein Pull** stattfindet — eine
  Prüfung behauptet, die es nur beim ersten Bezug gibt. **Der
  Quarantäne-Marker war da und stand am falschen Ort:** Er galt den
  Docker-Zeilen, während die eigene Lieblings-These (der git-Commit-Pin)
  ungekennzeichnet daneben stand. Jetzt trägt er beide. Eingetragen bei
  [`wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md).
- **Vierter Slice in Folge an Grenzen-Beschreibungen.** slice-212, slice-213
  und dieser berühren dieselbe Familie. Die Gefahr ist nicht Wiederholung,
  sondern **Selbstbezug**: ein Harness, der nur noch sich selbst beschreibt.
  Der Unterschied hier ist der Gegenstand — eine **Supply-Chain**-Frage nach
  [`AGENTS.md`](../../../../AGENTS.md) §3.1, die zufällig in einer
  Grenzen-Zeile landet. — **Ausgang:** eingetreten — **nicht** als Selbstbezug,
  sondern als Sachfund: Die Inventur hat fünf gepinnte Artefakte sichtbar
  gemacht, die in keiner Sensor-Beschreibung stehen, darunter die drei
  SHA-gepinnten Actions und ein **zweiter** `golang`-Digest im Werkzeug-Baum.
  **Das ist Supply-Chain-Bestand, keine Nabelschau.** Was der Slice
  **nicht** widerlegt hat, ist die Gefahr selbst: Vier von fünf Slices dieser
  Folge ändern nur Beschreibungen. Ob die Reihe damit endet, ist eine
  Planungs-Entscheidung und steht in der Closure-Notiz.

## 7. Closure-Notiz

**Geliefert.** Eine Inventur über **zwölf** gepinnte Fremd-Artefakte mit Ort,
Bindung und Prüfzeitpunkt (DoD 1) und die fehlende Antwort in
`harness/sensors/semgrep.md` (DoD 2). Ein unabhängiger Review, blockierend,
**zwei HIGH** und sechs MEDIUM — der erste HIGH-Befund dieser Slice-Folge.
`make gates` grün (zehn Gates, 736 Dateien).

**Was der Slice wert war, steht am Ende und nicht am Anfang.** Seine
Ausgangs-These — *„der `semgrep`-Cache ist das einzige ungeprüfte Artefakt
außerhalb des Repos"* — war **falsch**. Was an ihre Stelle trat, ist schärfer
und unbequemer: **Von zwölf gepinnten Fremd-Artefakten wird genau eines bei
jedem Lauf erneut geprüft — die vendorte Baseline, und ausgerechnet deren
Bindung beweist die Echtheit nicht** (slice-212). Alle übrigen werden **einmal
beim Bezug** geprüft und danach geglaubt. Der Cache ist der Normalfall; die
Ausnahme ist das Repo-interne Artefakt.

**Was Friktion war — und es ist diesmal der ganze Slice.** Beide Liefer-Punkte
trugen nicht:

- **Die Menge war um fünf Mitglieder zu klein.** §3 schrieb die **Form** eines
  gepinnten Fremd-Artefakts vorher aus — und gesucht wurde dann an **drei
  Fundorten** statt an ihr. Es fehlten drei SHA-gepinnte GitHub-Actions (mit
  eigener Hard Rule in §3.9 und drei Freshness-Achsen in §4), das `trivy`-Image
  und ein **zweiter** `golang`-Digest im Werkzeug-Baum.
- **Die Aussage je Mitglied war dreifach.** *„jeder Image-Bau"*, *„jeder
  Lauf"*, *„bei jedem Bezug"* — für denselben Mechanismus, im selben Artefakt.
  Gemessen: **kein Pull, 22 Schichten aus dem Cache.**
- **Und der falsche Schluss stand bereits im Gate-Vertrag**, als der Review ihn
  widerlegte.

**Steering-Loop-Lerneintrag: eine Form vorher auszuschreiben nützt nichts, wenn
man danach an Fundorten sucht.** Der Slice hat getan, was
[`AGENTS.md`](../../../../AGENTS.md) §5 seit slice-210 verlangt — die Form des
Gegenstands stand in §3, vor der Messung. **Gesucht wurde trotzdem dort, wo
schon einmal etwas gefunden worden war.** Die drei Fundorte waren die drei aus
dem Kopf; die Form hätte `.github/workflows/`, `tools/image-scan.sh` und
`tools/archive-wave/Dockerfile` mitgenommen. Eingetragen als 15. Beleg bei
[`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md).

**Der zweite Lerneintrag ist der peinlichste dieser Serie.** §8 schrieb
*„Register durchgegangen"* und folgerte, für `HARN` stehe **kein** Eintrag
darin. Es steht einer — offen —, und er betrifft `--check-latest`, also
**genau den Träger, den derselbe Slice als einzige Echtheits-Prüfung führt**.
Die Gesamtzahl 38 stimmte und deckte die nicht durchgeführte Teilprüfung zu.
**Eine korrekte Zahl neben einer angenommenen Aussage, in einem Satz** —
eingetragen bei
[`commit-message-overclaims-work`](../observations/BEO-ALL/commit-message-overclaims-work/observation.md)
(10×), weil die genannte Probe für dieses Kürzel nicht stattgefunden hat.

**Und die fünfte Instanz der frischen Regel widerlegt eine bequeme Lesart von
ihr.** *„Gegen den Gegenstand prüfen, nicht gegen die Beschreibung"* war
**angewandt**: Die neue Grenze stand gegen das Skript, und `[ ! -d … ]` steht
dort wirklich. Falsch war ihre Aussage über den **Bestand**. **Gegen den
Gegenstand zu prüfen genügt nicht, wenn der Gegenstand größer ist als die
Menge, die man gezählt hat** — das gehört zur Regel und steht im Register
([`grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md),
5×).

**Eine Planungs-Entscheidung, die hier hingehört.** Vier der letzten fünf
Slices ändern nur Beschreibungen. §6 hat das als Risiko benannt, und der
Ausgang sagt: Dieser Slice hat Sachbestand sichtbar gemacht — fünf gepinnte
Artefakte, die in keiner Sensor-Beschreibung stehen. **Die Gefahr des
Selbstbezugs ist damit nicht widerlegt, nur diesmal nicht eingetreten.** Der
nächste Vorgang sollte ein anderer sein; der unentschiedene eingehende CR
liegt bereit.

**Die drei Paarungen, gemessen.** **(a) Anker** — vakant: Der Slice verkörpert
keine Regel. **(b) Folge-Slice** — keiner genannt. **(c) Register** — alle
zitierten Pfade lösen auf; die vier neuen Belege liegen als
`evidence/slice-214.md` in ihren Verzeichnissen. Der Wachposten
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

**Zwei** Sub-Areas, und das ist neu gegenüber den drei Vorgängern:

- **`*`** (Repo-Default) — die Sensor-Beschreibungen unter
  [`harness/sensors/`](../../../../harness/sensors/), die geändert werden.
- **`tools/harness/`** (Kürzel `HARN`,
  [`MR-004`](../../../../harness/conventions.md#mr-004)) — **gelesen, nicht
  geändert.** Der Gegenstand von DoD (1) sind die Skripte und
  Konfigurationen, die den Bezug herstellen ([`tools/semgrep.sh`](../../../../tools/semgrep.sh),
  [`tools/harness/fetch-baseline-cache.sh`](../../../../tools/harness/fetch-baseline-cache.sh),
  [`Dockerfile`](../../../../Dockerfile), [`a-check.mk`](../../../../a-check.mk)).
  **Die Sub-Area wird geführt, weil sie den Gegenstand trägt, nicht weil sie
  ein Diff bekommt** — genau das verlangt das Inklusionskriterium, und §1
  schließt Verhaltens-Änderungen dort aus.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:274-274 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, **38** Verzeichnisse über beide
Kürzel). **Fünf** Einträge sind einschlägig — und der erste ist eine
Korrektur: Die erste Fassung behauptete, für `HARN` stehe **kein** Eintrag im
Register. Das war nicht nachgesehen, sondern angenommen.

- [`BEO-HARN/check-latest-blind-before-pin`](../observations/BEO-HARN/check-latest-blind-before-pin/observation.md)
  (Stand *offen*) — **der einschlägigste Eintrag überhaupt, und er wäre
  beinahe übersehen worden.** Er betrifft
  [`fetch-baseline-cache.sh`](../../../../tools/harness/fetch-baseline-cache.sh)
  `--check-latest` — also genau den Träger, den die Inventur als einzige
  Echtheits-Prüfung des vendorten Baums führt. Wer über Prüfzeitpunkte urteilt
  und diesen Eintrag nicht liest, urteilt an der eigenen Beobachtung vorbei.
- [`grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md)
  (4×, seit slice-213 verkörpert) — **der Anlass**, und seine Regel gilt diesem
  Slice unmittelbar: DoD (1) liest Konfiguration und Skript, **nicht** die
  Prosa darüber. Genau daran ist die Vermutung *„der `semgrep`-Cache ist
  ungeprüft"* schon einmal zerbrochen — er ist per Commit-SHA gebunden.
- [`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
  (14×, Stand *gemischt*) — **der gefährlichste hier.** Die Inventur hat eine
  **Menge** (welche Artefakte zählen?) und eine **Aussage je Mitglied**
  (Prüfzeitpunkt). slice-212 hat gezeigt, dass die Gefahr in der zweiten sitzt.
  **In diesem Slice traf es beide** — die Menge war um fünf Mitglieder zu
  klein, die Aussage je Mitglied in drei Varianten formuliert.
- [`wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md)
  (8×, Ausgang *geplant*) — §6 führt es als Risiko: Dass Docker einen Digest
  beim Bezug nachrechnet, ist **bekannt**, nicht in diesem Repo **gemessen**.
  Eine Inventur, die Vermutungen in Tabellenzellen schreibt, behauptet
  Prüfungen.
- [`zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
  (4×, verkörpert) — deshalb schreibt §3 vorher aus, was als *gepinntes
  Fremd-Artefakt* und was als *Prüfzeitpunkt* zählt. **Die erste Fassung hat
  die Form geschrieben und dann an drei Fundorten gesucht statt an ihr.**

**Keiner der fünf erreicht mit diesem Slice die Schwelle erstmalig.**

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-08 gelesen: **beide Nachtläufe grün** —
`upstream-drift.yml` (jüngster Lauf 2026-09-07T05:33:45Z) und `image-scan.yml`
(2026-09-07T08:21:32Z). **Beide sind für diesen Slice einschlägig, nicht
Routine:** Sie sind die einzigen Träger, die gepinnte Fremd-Artefakte gegen
upstream halten — und ihr Grün ist von gestern. **Genau das ist die
Prüfzeitpunkt-Frage, die dieser Slice stellt**, angewandt auf seine eigene
Vorprüfung.

**Modus-Begründungsblock.** Beide berührten Sub-Areas GF — ein Block je
Sub-Area, der zweite kurz, weil dort nichts geändert wird.

### Sub-Area: `*`

- **Modus:** GF (Repo-Default).
- **Konventions-Dichte:** hoch — die Form der Sensor-Dateien gibt die vendorte
  `gate.template.md` vor, und seit slice-213 sagt
  [`AGENTS.md`](../../../../AGENTS.md) §5, wie eine Grenze zu prüfen ist.
- **Phase-Reife:** Phase 5.
- **Evidenz-/Diskrepanz-Risiko:** **mittel.** Der Bestand liegt offen; das
  Risiko sitzt in der Antwort je Artefakt, und §6 führt es.
- **Reconciliation-Aufwand:** keiner (GF).

### Sub-Area: `tools/harness/`

- **Modus:** GF ([`MR-004`](../../../../harness/conventions.md#mr-004)).
- **Konventions-Dichte:** hoch — die Skripte sind über
  [`MR-004`](../../../../harness/conventions.md#mr-004),
  [`MR-005`](../../../../harness/conventions.md#mr-005) und
  [`MR-042`](../../../../harness/conventions.md#mr-042) konventionsgetragen.
- **Phase-Reife:** Phase 5.
- **Evidenz-/Diskrepanz-Risiko:** **niedrig** — die Sub-Area wird **gelesen**,
  nicht geändert; eine Diskrepanz zwischen Skript und Doku ist genau der Fund,
  den der Slice sucht, und kein Risiko seines Vorgehens.
- **Reconciliation-Aufwand:** keiner (GF).
