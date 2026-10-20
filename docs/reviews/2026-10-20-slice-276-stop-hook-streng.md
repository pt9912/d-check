# Review slice-276 — Stop-Hook streng mit HEAD-Stempel

- **Review-Art:** Code (gegen Plan, `MR-079`, `MR-076`, Hard Rules)
- **Gegenstand:** `slice-276`, Diff `2032fe1d..HEAD` (Code `dd92b601`, Doku `d1b49459`)
- **Skill:** `reviewer.md` @ 1.20.0
- **Modell-ID:** claude-sonnet-5-5
- **Datum:** 2026-10-20
- **Eingangs-Kontext:** Slice-Plan §1/§3/§6, `MR-079`, `MR-076`, `SPEC-094`, `SPEC-103`, `AGENTS.md` §3.1/§3.2/§3.6/§3.7/§3.8; Vergleich mit dem alten Hook (`git show 2032fe1d:.claude/hooks/stop-require-gates.sh`). Wegwerf-Proben in einem Scratch-Repo (gelöscht).

## Findings

### 1 — HIGH
- **quelle:** `AGENTS.md` §3.6; Skill Frage 1 (Pfad, auf dem ein Gate still grün wird); `MR-079` §Grenze
- **pfad:** `.claude/hooks/stop-require-gates.sh` · „if [ -e .harness/stop-gate-streng ]; then streng=1; fi"
- **befund:** Der strenge Modus hängt allein an der Existenz einer Datei im Arbeitsbaum. Wird sie unversioniert entfernt (`rm`, `git stash`, Branch-Wechsel, ein toter Symlink) und liegt der gestempelte HEAD an, gibt der Hook ein Turn-Ende mit ungedecktem Inhalt frei. Gemessen im Scratch-Repo: Nachweis geschrieben, Änderung an `f`, `rm .harness/stop-gate-streng` ergibt `approve`, Exit 0; mit toter Symlink ebenso. Mit Datei ergibt dieselbe Lage `block`. Der alte Hook hatte diesen Schalter nicht. Kein Gate meldet den Zustand: `selbstpruefung` ist nicht in `make gates` und prüft den Klon von `HEAD`, nicht den Arbeitsbaum. `harness/README.md` nennt das Löschen „keine Aufräumarbeit", verschweigt aber, dass es ohne Signal wirkt. Das ist die Lockerung, die der Auftraggeber ausgeschlossen hat, nur über einen Seiteneingang.
- **verifizierbar:** ja — Probe: Hook-Lauf mit gelöschter Datei und Änderung nach dem Nachweis
- **klasse:** `schalter-datei-ohne-halter`

### 2 — INFO
- **quelle:** `MR-079` §Begründung / SPEC-103
- **pfad:** `tools/harness/selbstpruefung.sh` · „mit der Datei gibt der Stop-Hook einen Commit ohne Nachweis frei"
- **befund:** Im strengen Zweig blockt der Commit-Fall schon über den Hash-Vergleich — dieser weicht bereits wegen der zuvor angelegten Datei `selbstpruefung-stop.txt` und der geleerten `stop-gate-streng` ab, nicht wegen des Commits. Nur im lockeren Zweig misst die Probe „neuer HEAD und abweichender Hash blockt". Die Zeile „streng wie locker" überdehnt deshalb für den strengen Zweig.
- **verifizierbar:** nein
- **klasse:** `probe-misst-nicht-ihren-namen`

### 3 — INFO
- **quelle:** Maintainability
- **pfad:** `tools/harness/selbstpruefung.sh` · „[ -f \"$klon/.harness/stop-gate-streng\" ]" gegen Hook `[ -e`
- **befund:** Selbstprüfung (`-f`) und Hook (`-e`) lesen „Datei vorhanden" verschieden; ein Verzeichnis oder Symlink dieses Namens ist für den Hook streng (gemessen: Verzeichnis ergibt `block`), für die Selbstprüfung ein Fehler. Keine Lockerung, nur Inkonsistenz.
- **verifizierbar:** ja
- **klasse:** `praedikat-drift-hook-probe`

### 4 — INFO
- **quelle:** Skill Frage 19 (Bestand, nicht Teil des Diffs)
- **pfad:** `.claude/hooks/stop-require-gates.sh` · „printf '%s' \"$input\" | grep -q"
- **befund:** `grep -q` unter `pipefail` kann bei sehr großer Eingabe mit SIGPIPE (141) enden; dann gilt `stop_hook_active` als nicht gesetzt und der Hook blockt erneut. Fällt in die sichere Richtung (blockt statt frei) und ist unverändert aus dem alten Hook übernommen.
- **verifizierbar:** nein
- **klasse:** `pipe-swallows-gate-exit-code`

## Negativbefunde (geprüft, ohne Befund)

- **Strenge gegenüber dem alten Hook, vorhandene Datei:** gemessen gedeckt, Änderung ohne Commit, detached HEAD mit Commit, Löschung committet, `stop_hook_active`, Verzeichnis als Datei, `STOP_GATE_STRENG=1` — jeweils wie der alte Hook; `" 1"` und `0` heben die Datei nicht auf und setzen den Modus ohne Datei nicht.
- **Trap (Frage 19):** `exit 0` in `approve`/`block` läuft im Hauptprozess (keine Subshell), der Trap sieht rc 0. Fehler außerhalb der benannten Lagen enden mit Exit 2 (außerhalb eines Repos: alt Exit 1, neu Exit 2 — die Behauptung in `MR-079` stimmt; geschlossenes stdout: Exit 2). JSON und Exit 2 entstehen nicht gemeinsam, da jeder Block erst druckt und dann mit 0 endet.
- **Kopplung `head_wert`:** Hook und `record-gates.sh` byte-gleich (`diff` über die Funktion).
- **Reihenfolge Hash/Stempel:** `record-gates` bricht vor beiden Schreibakten ab, wenn HEAD nicht auflösbar ist; ein halb geschriebener Zustand (Hash neu, Stempel alt) führt im lockeren Modus nur zu einem Hash-Vergleich, im strengen wird der Stempel nicht gelesen.
- **Selbstprüfung (5):** Der Nachweis kommt aus `record-gates.sh` direkt, die Hook-Entscheidung aus dem echten Hook, Exit ≠ 0 bricht ab; die Zählung „vier Zweige" stimmt mit den vier STOP-Zeilen überein.
- **SPEC-094, SPEC-103, README, `MR-079`:** stimmen mit dem Code überein, bis auf die in 1 benannte Auslassung zur Wirkung der Löschung.
- **Hard Rules:** §3.1 (nur bash/git/POSIX), §3.2 (keine `nolint`), §3.7 (Kommentare tragen Zusage/Kopplung/Grenze, Herkunft nur `MR-`/`seit slice-`), §3.8 (Hook liest Datei, schreibt sie nie; Grenze benannt) ohne Befund. §3.6: der Default des Hooks ist nicht gelockert — siehe aber Finding 1.
- **Abgrenzung §1:** `slice-mv.sh`, `extract-command.awk`, `pretooluse-command-guard.sh` unverändert (leerer Diff); Baseline-Pin unverändert.

## Kategorie-Summary

HIGH 1 · MEDIUM 0 · LOW 0 · INFO 3

## Verdikt

Blockiert durch Finding 1: Der Hook verhält sich bei vorhandener Datei wie der alte, die Datei als alleiniger Schalter lässt sich aber ohne Signal entfernen. Plan-Risiko §6 (Datei als einziger Halter) ist damit eingetreten, nicht entfallen.
