# Review — slice-154 / Commits `c7730ea`, `8ea69a9` (Nebenbefund `aab0cf4`)

**Review-Art:** Code/Design · **Gegenstand:** `slice-154-python-stage.md`,
`c7730ea` (Claim), `8ea69a9` (Ergebnis), `aab0cf4` (zweite
Durchsetzungs-Schicht) · **Skill:** `reviewer.md` @ 1.10.0 · **Modell:**
claude-opus-5[1m] · **Datum:** 2026-08-27

**Selbst gefahrene Sensoren:** `make doc-check` Exit 0 (529/0),
`make planning-check` Exit 0 (529/0), `make guard-probe` Exit 0 (40 Proben, 0
Fehlschläge). `make gates`/`make fullbuild` **nicht** gefahren.

## Findings

### F-1 · HIGH · Härtung der Durchsetzungsschicht ohne `MR`-Eintrag

**Pfad:** `.claude/settings.json:2-42`; `MR-040-…md:16`; `AGENTS.md:105-117`.

`aab0cf4` ist eine Härtung (zweite Schicht, 31 `deny`- und 4 `ask`-Regeln, dazu
`uv`/`uvx` neu in der Wächter-Sperrliste) und hat **keinen** `MR`-Eintrag. Der
Kanon ist wörtlich: *„Jede Härtung landet als neuer `MR-<NNN>`, der den
vorherigen schärft"* und *„Die Grenz-Zeile wird mitgezogen. Jeder Wächter-`MR`
trägt, was der Wächter nicht kann"* (`modul-13-quality-gates.md:204-213`).
Gemessen: die Permission-Schicht steht in **keinem** lebenden Dokument;
`AGENTS.md` §3.1 sagt weiter „ein Tool-Call-Wächter" — Singular. Die Grenzen der
neuen Schicht leben ausschließlich in der Commit-Botschaft.

Zusätzlich liest `MR-040:16` — *„Die Berechtigungsliste führt keine
Python-Einträge mehr."* — seit `aab0cf4` **falsch**: die Berechtigungsliste unter
`.claude/` führt jetzt `Bash(python *)`.

### F-2 · HIGH · Die Messung beantwortet eine andere Frage als die gestellte

**Pfad:** `MR-046-…md:16-22`; Botschaft `8ea69a9`.

Der Slice schreibt vor: *„Aus der **Sitzungshistorie** die Fälle sammeln … und
**je Fall** entscheiden."* Gemessen wurde der **committete Bestand**. Null Fälle
aufgezählt, null Zuordnungen. Die Populationen sind nicht dieselbe:
`MR-040`s eigene Begründung benennt den realen Bedarfsfall — *„ein ganzer
Arbeitstag mit Host-Python"* —, und der kann im Bestand **nicht** stehen. Der
Befund „leer" ist für die gestellte Frage weitgehend tautologisch; er belegt,
dass **Gates und Build** keinen Interpreter brauchen, nicht dass der **Bedarf**
null ist.

*Nicht bestritten:* dass das Ergebnis plausibel ist. Bestritten ist die Form des
Belegs.

### F-3 · MEDIUM · „zwölf ausführbare Skripte" passt zu keiner messbaren Menge

Gemessen (`git ls-files -s`): Modus `100755` → **10**; Endung `.sh` → **12**
(davon vier `100644`, via `bash tools/…` gerufen); bash-Shebang gesamt → **14**.
Die zwei ausführbaren Nicht-`.sh` — `.githooks/commit-msg`, `.githooks/pre-commit`
— fallen aus der Zählung **und** aus der Such-Enumeration heraus: dieselbe
ausgelassene Achse zweimal. *(Substanz hält: beide Hooks sind reines `bash` und
rufen nur `make`.)*

### F-4 · MEDIUM · Zwei der drei neuen §3.1-Aussagen sind am eigenen Bestand widerlegt

- *„Messungen macht das Produkt"* — sechs Messungen macht **nicht** das Produkt:
  `workflow-pins.sh`, `fetch-baseline-cache.sh --verify`, `pin-freshness.sh
  --compare`, `working-tree-hash.sh`, `guard-probe.sh`, `coverage-gate.sh`.
- *„was ein Gate-Skript braucht, tragen `bash` und die POSIX-Werkzeuge"* —
  falsch für zwei `make`-Targets: `fetch-baseline-cache.sh` braucht `curl`
  **und** `unzip`, `pin-freshness.sh` braucht `curl`. Weder `curl` noch `unzip`
  steht in der Host-Klasse von §3.1 Absatz 1.

### F-5 · MEDIUM · Die Grenz-Zeile der neuen Schicht ist halbseitig

Die vier genannten Durchfall-Klassen betreffen **alle** die Interpreter-Hälfte,
die der Wächter zweitdeckt. Die `git`-/`docker`-Hälfte ist **einschichtig**
(`BLOCKED` kennt weder `git` noch `docker`) und hat eigene, ungenannte Formen:
`git reset --hard` **ohne Argument**; `git push origin main --force` (Flag
hinter dem Positional); `git push origin +main`; `git clean --force` und
`git clean -x -f`. Ganz ohne Eintrag: `git branch -D`, `git checkout -- .`,
`git restore`, `git reflog expire`, `git gc --prune=now`, `git filter-branch`,
`git update-ref -d`; `docker rmi`, `docker image|volume|container|network prune`.

### F-6 · MEDIUM · Die neue Schicht hat keine wiederholbare Probe

`make guard-probe` fährt 40 Proben — **alle 40** treffen den Wächter; die 35
Permission-Regeln sind unbeprobt. Belegt ist **eine** Regel, einmalig.
`BEO-018`s eigene, „sofort wirksam" deklarierte Prozedur verlangt eine Probe
**je Form**; angewandt auf 1 von 35.

### F-7 · LOW · `uv`/`uvx` im Kopf-Kommentar und im Ablehnungs-Text nicht nachgezogen

### F-8 · LOW · `MR-046`s Auflösungs-Trigger nennt kein beobachtbares Ereignis

Er reproduziert die Urteilsfrage ohne Zähler, Schwelle oder Gelegenheit; und
*„dieser Eintrag wird von **ihr** [der Stage] abgelöst"* — ablösen tut ein
**Nachfolge-Eintrag**.

### F-9 · INFO · Die genannten Bedingungen beschreiben eine Stage; die Präzedenzfälle sind keine

`tools/semgrep.sh` (ADR-0010) und `a-check.mk` (ADR-0029) fahren digest-gepinnte
**externe Images**. `MR-040`s *„wie jede andere Toolchain dieses Repos"* trifft
für zwei von drei Fremd-Toolchains nicht zu.

### F-10 · INFO (vermutet) · `uses:`-Schritte laufen beim Runner unter `node`

Keine Host-Abhängigkeit im Sinn von §3.1, aber die N+1-te Form zur Aussage „kein
Interpreter in den Workflows"; aus dem Klon nicht verifizierbar.

### F-11 · INFO · Beide Schichten sind über `make` konstruktionsbedingt umgehbar

Die Regel sieht `make clean`, der Wächter `make` in Befehlsposition; ausgeführt
wird `docker image rm`. Gewollt — steht aber in keiner Grenz-Aussage.

### F-12 · LOW · §Offene Wellen trägt 19 Leerzeilen

Gate-grün, aber ein Roadmap-Leser sieht nicht, was in Arbeit ist.

## Negativbefunde

- **Interpreter-Sweep über den gesamten getrackten Baum** (19 Muster, via
  `grep -f`): **kein einziger Aufruf**. Ein Treffer ist ein Falsch-Positiv
  (`planning.go:230` — `"ruby"` als HTML-Elementname). Die inhaltliche Aussage
  von `MR-046` trägt also **weiter** als seine Enumeration.
- **`.githooks/`**, **Makefile-Includes**, **Docker-Basis-Images**,
  **`go:generate`** (0 Vorkommen), **Tests/Tool-Configs**, **Untracked-Dateien**:
  ohne Befund.
- **`c7730ea`** reiner Move (R100) plus Roadmap-Flip; nichts nachzuziehen.
- **Form von `MR-046`** entspricht der etablierten; `doc-check` grün.
- **Datei-Zahlen der Botschaften** konsistent (528 → 529 = die neue Datei).

## Kategorie-Summary

HIGH 2 · MEDIUM 4 · LOW 3 · INFO 3

## Urteil

**Schließbar nach Nacharbeit.** Der Kern ist richtig — eine Regel, die auf eine
nicht existierende Stage zeigte, wurde zurückgeschnitten; unabhängig bestätigt,
dass im gesamten Baum kein Interpreter gerufen wird. Blockierend: **F-1** (eine
ganze zweite Durchsetzungs-Schicht ohne `MR`-Eintrag, gegen einen wörtlichen
Kanon-Satz) und **F-2** (der Bedarfs-Beleg misst eine andere Population als die,
die DoD-Haken 1 verlangt).
