# Slice slice-157: Die Konfiguration der Durchsetzungs-Schichten wird nicht geprüft

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`MR-047`](../../../../harness/conventions.md#mr-047) (die zweite
Schicht und ihre Datei); [`MR-042`](../../../../harness/conventions.md#mr-042)
(der Wächter und seine Verdrahtung);
[`AGENTS.md`](../../../../AGENTS.md) §3.1;
[`BEO-018`](../observations.md) (ein Wächter fällt aus, und niemand merkt es).

**Berührte Spec-Stellen:** — (Durchsetzungsschicht; das Produkt bleibt
unberührt).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

Beide Durchsetzungs-Schichten für [`AGENTS.md`](../../../../AGENTS.md) §3.1
hängen an **einer** Datei: `.claude/settings.json` trägt die Hook-Verdrahtung
**und** die Permission-Sperrliste. Ist sie kaputt oder verschoben, fällt beides
zugleich aus — und zwar **still**, denn kein Gate liest sie.

**Das ist keine Vermutung, sondern der zweite Anlauf derselben Klasse.** Der
Wächter war schon einmal einen Arbeitstag lang wirkungslos, ohne dass ein
Sensor etwas sagte ([`BEO-018`](../observations.md)); damals lag es an einer
Zeile im Wächter, hier läge es an der Datei, die ihn überhaupt aufruft.

**Was heute an ihrer Stelle steht, ist ein von Hand geschriebener
`awk`-Klammerzähler** — im Arbeitsverlauf zweimal gebraucht, um zu prüfen, ob
eine Bearbeitung die Datei zerschossen hat. Er beantwortet **Balance**, nicht
**Gültigkeit**, und er ist genau der „Parser durch die Hintertür", vor dem der
Extraktor-Kommentar warnt. Der Slice ist die Konsequenz aus dem offenen Rest
von [slice-154](../done/slice-154-python-stage.md).

## 2. Vorgehen

1. **Erst den Gegenstand festlegen, dann das Werkzeug.** Welche Dateien tragen
   Durchsetzungs-Zusagen — `.claude/settings.json`, `.claude/hooks/*`, die
   Verdrahtung selbst? Und welche Aussage soll das Gate halten: *gültiges JSON* ·
   *Hook ist verdrahtet und die Datei existiert* · *die Sperrliste deckt die
   Klasse aus §3.1*? Die drei sind verschieden teuer und verschieden wertvoll.
2. **Die billigste tragende Aussage zuerst.** Ein Gate, das nur „gültiges JSON"
   sagt, fängt den Zerschuss-Fall — und mehr braucht es vielleicht nicht.
3. **Kein Fremd-Parser.** `jq` wäre dieselbe Abhängigkeit unter anderem Namen
   ([`MR-046`](../../../../harness/conventions.md#mr-046)). Naheliegend ist ein
   Go-Test im vorhandenen Testbaum: die Toolchain ist da, die Prüfung ist
   getippt statt gezählt, und sie läuft in `make test` ohne neues Target.
4. **Die Gegenprobe gehört dazu:** eine kaputte Kopie muss rot werden, und der
   unveränderte Bestand grün ([`BEO-017`](../observations.md)).
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine vierte Toolchain.** Weder `jq` noch ein Skript-Interpreter.
- **Keine Bewertung der Regel-Inhalte.** Ob die Sperrliste die richtigen Namen
  führt, ist ein Urteil und gehört in den Konventionsspeicher, nicht in ein
  Gate.
- **Kein Gate auf den Wächter selbst.** Der bleibt werkzeug-lokal
  ([`MR-042`](../../../../harness/conventions.md#mr-042)); `make guard-probe`
  ist seine Probe und bewusst kein Gate.

## 4. Definition of Done

- [x] Der Gegenstand ist benannt: welche Dateien, welche Aussage, und was das
      Gate ausdrücklich **nicht** zusagt.
- [x] Die Prüfung läuft in `make gates` und ist netzlos.
- [x] Gegenprobe belegt: kaputte Kopie rot **mit gelesener Meldung**,
      unveränderter Bestand grün. *(Vier Fälle: zerschossenes JSON, driftende
      Sperrliste, Hook-Pfad ins Leere in **beiden** Schreibweisen, fremder
      `matcher`. Bestand danach jedes Mal byte-identisch.)*
- [x] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein Gate auf eine Werkzeug-Datei ist eine Grenzüberschreitung.**
  `.claude/` ist Werkzeug-Einstellung, kein Repo-Vertrag — und §3.1 sagt
  ausdrücklich, ein Lauf ohne dieses Werkzeug sei ungebunden. Ob ein
  **Repo**-Gate über eine Werkzeug-Datei wachen darf, gehört entschieden, nicht
  vorausgesetzt. — **Ausgang: eingetreten, und der erste Anlauf hatte die Frage
  tatsächlich vorausgesetzt** — der Review hat das gefangen. Entschieden als
  [`MR-048`](../../../../harness/conventions.md#mr-048): eine
  Werkzeug-Einstellung **darf fehlen**, dann überspringt die Prüfung; **kaputt
  sein darf sie nicht**. Damit ist auch der dritte Ausschluss aus §3 gewahrt,
  den der erste Anlauf gar nicht behandelt hatte: ohne das Werkzeug behauptet
  das Gate nichts, und [`MR-042`](../../../../harness/conventions.md#mr-042)s
  Warnung vor einer Zusage, „die ein Lauf ohne dieses Werkzeug nicht hält", ist
  eingelöst statt umgangen.
- **Gültiges JSON ist nicht dasselbe wie wirksame Verdrahtung.** Ein Gate, das
  nur die Syntax prüft, könnte als Zusage gelesen werden, die Durchsetzung
  stehe — das wäre schlimmer als kein Gate. Die Grenze gehört in dieselbe
  Zeile wie die Zusage. — **Ausgang: eingetreten, in einer Form, die der Punkt
  nicht vorhergesehen hatte.** Die Grenze stand von Anfang an neben der Zusage,
  war aber **unvollständig** (die blinde Notations-Klasse fehlte) und an einer
  Stelle **falsch** (der Versions-Suffix-Fall war der falschen Schicht
  zugeschrieben). Eine Grenz-Zeile, die die verkehrte Klasse benennt, ist
  gefährlicher als keine — sie erklärt eine offene Lücke für geschlossen.
  Beides korrigiert; die Zusage ist dabei zugleich enger und tragfähiger
  geworden.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls der erste Abnahme-Punkt gegen
das Gate entschieden wird — dann bleibt die Lücke benannt, und das ist ein
Ergebnis.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Durchsetzungsschicht (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27):
  [`BEO-018`](../observations.md) — wer den Rand eines Wächters anfasst, probt
  jede Form, die dieser Rand annehmen kann;
  [`BEO-017`](../observations.md) — ein rotes Gate muss vom geprüften Grund
  kommen.

Slice-ID: slice-157. Betroffene IDs: — (kein `DC-`-Bezug;
Durchsetzungsschicht). Module: — . Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — neue Prüfung an einem vorhandenen Artefakt
der Durchsetzungsschicht.

## 9. Closure-Notiz (nach `done/`)

**Die Datei, an der beide Durchsetzungs-Schichten hängen, hat einen Sensor.**
Vier Zusagen in `make test` und damit in `gates`, netzlos, ohne neue Toolchain:
`.claude/settings.json` ist gültiges JSON, jeder verdrahtete Hook-Pfad
existiert, der Befehls-Wächter hängt am `Bash`-Werkzeug, und die
Permission-Sperrliste deckt jeden Namen der Wächter-Sperrliste als **ganze**
Befehlsklasse.

**Der Review hat den teuersten Befund gemessen, nicht vermutet.** Zwei
Mutationen, die sich **nur in der Schreibweise** unterschieden — nicht im
Zustand —, ergaben grün und rot: die Pfad-Suche verlangte einen führenden
Schrägstrich, und die literale Relativform fiel damit heraus. Der Test hielt
genau den Fall nicht, für den er gebaut war. Schlimmer noch war die
Grenz-Zeile daneben: sie erklärte **Variablen**-Pfade für blind, während die
literale Form die blinde war. Wer sie las, schloss das Gegenteil des Wahren.

**Dieselbe Klasse traf eine zweite Achse, die der Slice gar nicht auf dem Zettel
hatte.** Der `matcher` wurde dekodiert und nie gelesen — ein Wächter am
matcher `Write` feuert auf keinen Bash-Aufruf mehr, ohne dass eine Datei fehlt
oder kaputt ist. Ein dekodiertes Feld liest sich wie ein geprüftes. Jetzt
geprüft, mit Gegenprobe.

**Die Kopplung der zwei Sperrlisten war schwächer, als ihre Zusage klang.** Sie
mass Namens-**Nennung**: `Bash(pip install *)` zählte als Deckung für `pip`,
obwohl `pip download …` durchliefe, und ein formloser Eintrag ohne
`Bash(`-Hülle ebenso — eine Regel, die das Werkzeug gar nicht auswertet, galt
als Schutz. Jetzt zählt nur die ganze Befehlsklasse.

**Die Grenzfrage aus §5 war vorausgesetzt statt entschieden.** Sie ist es
jetzt, als [`MR-048`](../../../../harness/conventions.md#mr-048), und die
Antwort trennt zwei Dinge, die leicht zusammenfallen: eine Werkzeug-Einstellung
darf **fehlen** — §3.1 erklärt einen Lauf ohne sie für ungebunden —, aber sie
darf nicht **kaputt** sein. Fehlt sie, überspringt die Prüfung; ist sie da, muss
sie tragen. Ohne diese Trennung stünde die Wahl zwischen zwei schlechten Enden:
ein Gate, das jeden Klon ohne dieses Werkzeug rot färbt, oder gar keines.

**Was das Gate ausdrücklich nicht sagt**, und es steht im Test, in §4 und in der
Sensors-Zeile: dass die Durchsetzung **läuft**. Ob das Werkzeug den Hook ruft
und seine Antwort befolgt, steht in keiner Datei. Grün schließt einen
Ausfallgrund aus, nicht den Ausfall. Dazu drei benannte Grenzen: die
Pfad-Prüfung ist eine Teilketten-Suche statt einer Shell-Analyse, ein Pfad
ausserhalb des Repos würde falsch zugeordnet, und
`.claude/settings.local.json` ist nicht Gegenstand.

**Eine eigene Regelverletzung im selben Zug.** Die Kommentare des ersten
Anlaufs trugen zwei Herkunfts-Felder statt einem, eine Slice-Nummer, einen
`BEO`-Verweis und Mess-Prosa — vier Formen, die §3.7 ausdrücklich ausschließt,
und keine davon grandfathered. Geräumt.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 532 Dateien, 0 Befunde);
vier Gegenproben mit gelesener Meldung, Bestand danach jedes Mal
byte-identisch. Ein unabhängiger Review ist gelaufen; seine zwei HIGH, drei
MEDIUM und drei LOW sind in `ba56d44` eingearbeitet.
