# Review — slice-155 / Commits `617fbfa`, `49dc022`

**Review-Art:** Code + Doku · **Gegenstand:** `slice-155-guard-ohne-node.md`,
`617fbfa` (Wächter in `bash` + `awk`, Extraktor, Proben-Target, `MR-042`,
`MR-041` → `done/`), `49dc022` (Werkzeug-Einstieg, `MR-043`) · **Skill:**
`reviewer.md` @ 1.10.0 · **Modell:** claude-opus-5[1m] · **Datum:** 2026-08-27

## 0. Kontext-Beleg

Vor dem ersten eigenen Tool-Aufruf lagen im Kontext vor:

1. **`CLAUDE.md`** — vollständig. Erste Zeile: `# Claude Code Einstieg — d-check`
2. **`AGENTS.md`** — vollständig. Erste Zeile:
   `# AGENTS.md — Briefing für AI-Coding-Agenten`
3. Der Memory-Index des Nutzers.

**Nicht** vorgelegen: `harness/README.md`, `harness/conventions.md`, der
Reviewer-Skill, das Regelwerk, der Slice.

**Nebenbefund, zugleich Beleg für `MR-043`:** dass **AGENTS.md** im Kontext lag,
ist die direkte Wirkung des `@AGENTS.md`-Imports aus `49dc022`. Die Mechanik ist
damit an dieser Sitzung **gemessen**, nicht erschlossen.

**Gefahren:** `make guard-probe` (`== Fehlschläge: 0`, **28** Fall-Zeilen),
`make doc-check` (525/0), `make gate-consistency` (525/0), `make planning-check`
(525/0), `make trace-check RANGE=de978a7..HEAD` (525/0), `make help`, eigene
Proben (~45 Fälle), mechanischer Wortlisten-`diff` alt↔neu. Die alte Fassung
konnte **nicht** ausgeführt werden (sie braucht `node`); der Alt/Neu-Vergleich
ist Zeile-für-Zeile gelesen plus mechanischer Listen-`diff`.

## Findings

### F-1 · MEDIUM · „Eine einzige Host-Abhängigkeit" trägt nicht — und die zweite ist fail-**open**

`MR-042-…md:19-21`; Guard `:6`, `:33`, `:112`. Der Wächter ruft neben `awk` auch
`dirname` und `cat`. Gemessen mit `env -i PATH=/nonexistent /bin/bash <guard>`:
`exit 127`, beide „command not found", **keine Block-Ausgabe**. Ein
PreToolUse-Hook, der mit 127 und ohne JSON endet, blockiert nicht — während
`MR-042` schreibt *„Fail-closed ist die Voreinstellung, nicht die Ausnahme"*.
*Gegenprobe:* die beiden benannten Fälle blocken korrekt.

### F-2 · MEDIUM · Die Grenz-Zeile beschreibt die Lücke zu eng — und Modul 13 verlangt sie im Wächter-`MR`

Guard `:26-27`; `AGENTS.md:100-108`; `MR-042` (ganzer Eintrag). Gemessen
(Erwartung „block", Ergebnis **pass**): `if true; then pip install x; fi` ·
`for f in a; do go build; done` · `while true; do npm i; done` · `! pip install x` ·
`nohup pip install x` · `timeout 5 go build ./...` · `stdbuf -o0 pip install x` ·
`setsid go build` · `p"i"p install x` · `g\o build` ·
`bash -c "bash -c \"bash -c pip\""`.

Keine dieser Klassen ist „ein Interpreter, den die Liste nicht kennt" — es sind
**gelistete** Werkzeuge über Konstruktionen, die die Segmentierung nicht
modelliert. Alle **geerbt**. `modul-13-quality-gates.md:186ff` verlangt für
jeden Wächter-`MR` die Grenz-Zeile; **MR-042 trägt gar keine**.

### F-3 · MEDIUM · Der Extraktor ist bei zwei aufeinanderfolgenden Strings nicht fail-closed

`extract-command.awk:56-68`, `:82-87`. Gemessen:
`{"tool_input":{"command":"pip x""make gates"}}` → **pass**. Nach dem Wert-Ende
bleibt `wantkey[depth] == 0`; jeder weitere String wird wieder als Wert gelesen
und **überschreibt** `cmdval`. *Erreichbarkeit:* kein Weg gefunden, diese Form
über den Befehlsinhalt zu erzeugen — eine Zusagen-Lücke, kein Angriffspfad.

### F-4 · MEDIUM · `\u`-Escape im **Schlüssel** wird still falsch dekodiert

`extract-command.awk:45-51`, `:58-59`. Für **Werte** fängt Zeile 62 das ab, für
**Schlüssel** nicht: `curkey[depth]` bekommt die verkürzte Zeichenkette. Zwei
gemessene Eingaben mit escaptem Schlüssel laufen durch — gültiges JSON, kein
Urteil, kein Zweifel.

### F-5 · MEDIUM · `guard-probe.sh` verwirft den Exit des Wächters hinter der Pipe

`guard-probe.sh:12`, `:20`. Das Urteil hängt allein am Wort `decision` auf
stdout. Ein Wächter, der mit 127 und ohne Ausgabe endet, ergibt **`pass`**. Der
Slice führt diese Klasse in §7 selbst als sichtungspflichtig (`BEO-007`).

### F-6 · MEDIUM · `@AGENTS.md` ist kein Markdown-Link mehr — das Ziel fällt aus der Gate-Deckung

`CLAUDE.md:3` (nach `49dc022`); `MR-043-…md:13-22`. Vorher ein Link, den
`doc-check` auflöst; jetzt blanker Text. `doc-check` läuft grün, weil es
**nichts mehr zu prüfen gibt**. Failure-Szenario: `AGENTS.md` wird umbenannt —
kein Gate meldet etwas, jeder Folgelauf startet ohne die Hard Rules. MR-043
nennt zwei Grenzen des Imports, diese dritte nicht.

### F-7 · LOW · `§Artefakt-Set` gibt es nicht — der Abschnitt heißt `§Referenz-Implementierung`

`MR-042-…md:8`; `MR-043-…md:6`. „Das vollständige Artefakt-Set" ist der erste
Satz des Abschnitts, keine Sektion. Beide Links tragen **keinen** Anker, deshalb
ist der falsche Sektionsname für `anchors` unsichtbar.
*Gegenprobe:* `modul-13 §Guard-Härtung` existiert wirklich; `modul-02
§Freshness-Audit` ebenso.

### F-8 · LOW · Blockzitat in `MR-043` setzt eine Hervorhebung, die die Quelle nicht hat

`MR-043-…md:7-8`: *„Sie **bringt** …"* — `bringt` ist in der Quelle **nicht**
fett. Die Fettung ist die Beweislast des Eintrags, im Zitat gesetzt statt
daneben (`BEO-012`-nah).

### F-9 · LOW · Der `MR-041`-Move liegt im Feature-Commit statt in einem eigenen Move-Commit

§3.3 verlangt den `git mv` als eigenen Commit. Praktischer Schaden: keiner — der
Rename wird mit **75 %** erkannt, `git log --follow` hält.

### F-10 · LOW · `MR-042`s `Geltungsbereich` nennt das `Makefile` nicht

Das `guard-probe`-Target lebt dort, und der Eintrag argumentiert mehrfach über
*„Die Proben sind ein `make`-Target"*.

### F-11 · LOW · Verwaiste `node`-Erlaubnis in der Berechtigungsliste

`.claude/settings.local.json:82` — `"Bash(node tools/docs-check.js)"`. `MR-040`
führt die Berechtigungsliste im Geltungsbereich; das Ziel existiert nicht mehr.
Wirkungslos (der Hook greift vorher), Datei untracked — lokale Hygiene.

### F-12 · LOW · Die Sensors-Zeile zählt vier von fünf Gegenkontrollen

`harness/README.md:81` nennt vier; im Probenlauf sind es fünf.

## Negativbefunde

- **Denylists Alt↔Neu** — `BLOCKED` (24 Wörter), `PREFIXES` (8), `SHELLS` (5)
  mechanisch diffed: **byte-identisch**, keine Auslassung.
- **Segmentierung** — sequentielle `${s//…}`-Ersetzungen äquivalent zur
  JS-Alternation; Reihenfolge korrekt; Grenzfälle `a &&& b`, `a |& b` von Hand
  durchgerechnet.
- **Präfix-Überspringen**, **Basename-Kürzung**, **Versions-Suffix-Muster**,
  **`strip_quotes`** — Zeile für Zeile äquivalent.
- **Rekursion** — `-c` in Flag-Bündeln gefunden; Tiefenlimit greift.
- **Fail-closed-Kern** — fehlender Extraktor, fehlendes `awk`, leere Eingabe,
  malformes/abgeschnittenes JSON, `\u` **im Befehl**, Müll außerhalb eines
  Strings: alles einzeln gemessen, alles Block.
- **Feldwahl** — `description` vor `command` → richtiger Wert; Schlüsselname als
  *Daten* im Wert → korrekt geblockt.
- **Prozess-Substitution/Heredoc** — `diff <(go env) b` → block; mehrzeilige
  Heredoc-Botschaft mit `pip` → block.
- **Gate-Deklarationen** — `gate-consistency` grün ⇒ `make guard-probe` ist in
  `Makefile`, `AGENTS.md` §4 und §Sensors konsistent deklariert; `make help`
  führt es.
- **`planning-check`** und **`trace-check`** grün.
- **§3.7-Zustandsfelder**, **CHANGELOG** (korrekt unangetastet),
  **`Ersetzt-Baseline-Regel: keine`**: ohne Befund.

## Kategorie-Summary

HIGH 0 · MEDIUM 6 · LOW 6

## Urteil

**Schließbar nach Nacharbeit.** Die teuerste Frage ist beantwortet und fällt gut
aus: **die Übersetzung verliert keine Prüfung.** Denylists byte-identisch, Logik
äquivalent; die einzige gemessene Verhaltensdifferenz geht überwiegend in die
*schärfere* Richtung. Die Zeichen-Whitelist ist eine echte Härtung gegenüber der
Vorlage. Nachzuarbeiten sind **Zusagen gegen Bestand**, nicht die Prüflogik.
