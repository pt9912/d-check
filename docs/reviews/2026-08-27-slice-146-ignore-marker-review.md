# Review — slice-146 / Commit `7b44d8e`

**Review-Art:** Plan + Closure (kein Produkt-Delta im Diff) · **Gegenstand:**
`slice-146-ignore-marker-wirkung.md`, `slice-159-ignore-marker-syntax.md`,
`observations.md` (BEO-013), Commit `7b44d8e` · **Skill:** `reviewer.md` @
1.10.0 · **Modell:** claude-opus-5[1m] · **Datum:** 2026-08-27

**Eigene Messungen:** `awk`-Nachbau von `proseLines` (`TrimFenceIndent` +
`FenceToggle` byte-genau nachgebildet) über die Scan-Menge aus `.d-check.yml`;
Probe mit ins Leere gelenkter Marker-Konstante, `make doc-check`, danach
Wiederherstellung aus Sicherungskopie (`git diff --stat` leer).

## Findings

### F-1 · MEDIUM · „233 außerhalb von Fences" ist die Menge von zwei der vier Konsumenten

**Pfad:** slice-159 `:26-29`; abgeschwächt slice-146 `:104-108`.

- `versions` liest **nicht** `proseLines`, sondern alle Zeilen
  (`versions.go:52`); die Spec sagt es zu: Vorkommen werden ausgenommen „in
  Fence wie außerhalb" (`spezifikation.md:1553-1554`). Gemessen: **236**
  Zeilen gesamt, 233 davon außerhalb von Fences — für `versions` ist 233 **zu
  klein**.
- `diagrams` liest **ausschließlich** Zeilen innerhalb von Diagramm-Fences plus
  die Öffnungszeile (`diagrams.go:33-36`). Die 233 Prosa-Zeilen unterdrücken
  für `diagrams` **gar nichts**; die Mengen sind disjunkt.

Die drei aus der Prosa-Menge fallenden Zeilen: `adr/0040-…:39` und
`reviews/2026-07-17-slice-074-implementation-r3.md:82` (in Fence, wirkungslos)
sowie **`docs/user/benutzerhandbuch.md:1125`** — eine `mermaid`-Öffnungszeile
mit **gesetztem, dokumentiertem** Marker, die weder in den 233 noch in den 48
noch in den 185 auftaucht. Zweitens ist die Wirkung in `diagrams` nicht
zeilenweise: auf der Öffnungszeile nimmt der Marker den **ganzen Block** aus.

### F-2 · MEDIUM · Die Aufteilung 87/146 trennt nicht Marker von zitierter Marker-Form

**Pfad:** slice-146 `:110-115`; slice-159 `:28-29`, `:50-52`; `observations.md:16`.

87 + 146 = 233 stimmt. Von den 87 „Kommentar-Form"-Zeilen stehen aber **25** als
zitierte Form **innerhalb von Inline-Code** (u. a. `benutzerhandbuch.md:1117,
1650`, `spezifikation.md:407`, `adr/0040-…:26,116`, `CHANGELOG.md:625,1712`).
**Echte gesetzte Marker sind 62, Erwähnungen 171.** Folge: eine Verengung auf
die rohe Teilkette `<!-- d-check:ignore` würde genau diese 25 **weiter**
unterdrücken — der diagnostizierte Defekt bliebe zu 29 % bestehen.

### F-3 · MEDIUM · Die vorgeschlagene Kommentar-Form widerspricht einer bindenden Spec-Festlegung

**Pfad:** slice-159 `:49-53`, `:66-70`. `spezifikation.md:1470-1477`
(§DC-FA-DIAG-001.a Schritt 5, „festgelegt, nicht abgeleitet"):

> Der Marker ist ein **Token**, kein HTML-Kommentar. … Eine Kommentar-Lexik je
> Fence-Sprache wäre ein Grammatik-Parser und widerspräche Schritt 2.

Dasselbe als Nutzer-Zusage in `benutzerhandbuch.md:1115-1132` und als Kommentar
in `diagrams.go:21-24`; die gesetzte Fundstelle existiert
(`benutzerhandbuch.md:1125`). Die in §2.1 gestellte Frage „gibt es eine gesetzte
Unterdrückung, die die neue Form verlöre?" hat damit eine im Repo dokumentierte
Antwort — und slice-159 nennt sie nicht.

### F-4 · MEDIUM · „107 tragen gar keinen §5-Abschnitt" ist falsch

**Pfad:** slice-146 `:37-39`; `.d-check.closure.yml:163-164`.

**Alle 150** `done/`-Slices tragen genau eine `## 5.`-Überschrift; 107 unter
einem **anderen Titel** — `## 5. Trigger` 43 · `## 5. Closure-Trigger` 31 ·
`## 5. Risiken / offene Punkte` 19 · `## 5. Offene Punkte / Risiken` 4 ·
`## 5. Definition of Done` 3 · sieben Einzelfälle. Davon sind **25** echte
Risiko-Abschnitte unter älterem Titel.

Zweitens verdeckt die Sammelzahl den schwersten Bestands-Befund: **3 der 14**
(`slice-106`, `slice-110`, `slice-111`) tragen ein wohlbenanntes §5 mit Risiken
und **überhaupt keinen** Ausgangs-Marker. Das ist der Kanon-Kernsatz, dreimal
verletzt. Die 14 zerfallen in: **4** wirklich nur-Freitext · **3** ganz ohne
Ausgang · **7** Muster-Artefakte.

### F-5 · MEDIUM · Die benannte Grenze nennt eine Richtung; sie ist in beide weiter

**Pfad:** slice-146 `:28-34`; `AGENTS.md:349`; `.d-check.closure.yml:149-153`.

„Ein §5 mit zwei Risiken, eines kanonisch, eines Freitext, läuft grün durch" ist
korrekt. „Gedeckt ist der Abschnitt, in dem **kein** Ausgang kanonisch ist" ist
es **nicht**: die Bedingung ist eine Existenzaussage über den ganzen Text, nicht
über Ausgänge. Probe: ein §5 mit **einem** Risiko (Ausgang `*behoben*`) plus dem
Prosasatz „Die Vorlage verlangt den Wortlaut **Ausgang:** entfallen …" ⇒
**grün**.

### F-6 · MEDIUM · „belegt nicht ausdrückbar" — gemessen ist eine Formulierung, nicht die Frage

**Pfad:** Commit-Botschaft „Messung 1"; slice-146 `:28-34`.

Gemessen: RE2 weist `(?!` ab. Behauptet: die Korrelation Risiko ↔ Ausgang sei
nicht ausdrückbar. Das „damit" trägt nicht — `forbid-pattern` ist über **jedes
Vorkommen** quantifiziert und damit je Risiko wirksam; das Komplement einer
endlichen Wortmenge ist in RE2 ohne Lookahead darstellbar. Gegengemessen mit
einer Präfix-Komplement-Regel: über die 12 nicht ausgenommenen `done/`-Slices
**Exit 0, 0 Befunde**; Probe `m1` (ein kanonischer, ein Freitext-Ausgang) ⇒ die
`require-pattern` bleibt grün, die Komplement-Regel meldet `section-forbidden`;
Probe `m2` (drei kanonische Ausgänge in drei Schreibweisen) ⇒ grün.

### F-7 · MEDIUM · „schärft MR-006" hat im Eintrag keine Stütze

**Pfad:** `harness/conventions.md:129`; `MR-049-…md:1`. MR-006 regelt die
**Referenzrichtung**; MR-049 nennt ihn im Körper kein einziges Mal, sein Feld
*Ersetzt-Baseline-Regel* sagt „keine". Das ist die Gestalt von `BEO-012` in der
Erzeugungsrichtung.

### F-8 · LOW · Ein Ausgang in Inline-Code ist unsichtbar

`SectionProse` leert Inline-Code-Spannen; bei der Nachbarregel ist das ein
gewollter Preis, hier kehrt sich die Wirkung um — die Bereinigung verdeckt
**Erfüllung** statt Verstoß.

### F-9 · LOW · Ein §5 ohne Risiken ist rot

Probe: §5 mit „Keine offenen Risiken; dieser Slice ist rein redaktionell." ⇒
`section-pattern-missing`.

### F-10 · LOW · Der Selektor ist eine exakte Zeichenkette, die die Baseline-Vorlage nicht trägt

Die vendorte `slice.template.md` führt die Risiken unter
`## 6. Risiken und offene Punkte`; ein daraus geschriebener Slice wäre
`section-missing`.

### F-11 · INFO · Der Befund zeigt auf die Überschrift, nicht auf das Risiko

Bei einer Regel je Abschnitt geht es nicht anders; „an der richtigen Zeile"
liest sich aber, als zeigte der Befund auf den Ausgang.

## Negativbefunde (geprüft, ohne Befund)

- **„Frühes `continue` VOR der Prüfung" in allen vier Konsumenten** — trägt.
  `codepaths.go:44`, `ids.go:41` **und** `ids.go:83`, `versions.go:53`,
  `diagrams.go:35`. Genau vier Verwender.
- **Fence-Erkennung** — der `awk`-Nachbau ist weder strenger noch laxer:
  exakt **233** auf `7b44d8e^`, **235** auf `HEAD`.
- **58 / 48 / 38 / 18 / 2** — exakt reproduziert.
- **Arithmetik** — 87+146 = 233 ✔; 48+185 = 233 ✔; alle 48 Befund-Zeilen liegen
  in den 233 (`comm -23` leer).
- **Register-Form BEO-013** — Zähler 1 ↔ ein Beleg, Form gehalten, §3.7
  eingehalten, Zähler korrekt nicht erhöht.
- **Drei-Ausgänge-Regel** erfüllt; Folge-Slice mit ID genannt.

## Kategorie-Summary

MEDIUM 6 · LOW 3 · INFO 1

## Urteil

**Schließbar nach Nacharbeit.** Alle sechs Zahlen sind reproduzierbar, die
Kernaussage über die vier Konsumenten stimmt, die Entscheidung gegen den Bau ist
im Ergebnis richtig, der Carve-out schneidet den richtigen Gegenstand.
Nachzuarbeiten vor dem `git mv`: die falsche `DC-*`-Kante (F-5 dieses Reports
bezieht sich auf slice-146s Kopf-Feld `DC-FA-MTX-003`), Population und
Marker-Begriff je Konsument, der Spec-Widerspruch, der Schluss auf die gemessene
Regel-Gestalt, und der DoD-Haken.
