# Review slice-276 — Stop-Hook streng, Runde 2 (Fix-Commits)

- **Review-Art:** Code (gegen Plan, `MR-079`, `MR-076`, Hard Rules)
- **Gegenstand:** `slice-276`, Diff `d1b49459..HEAD` (Plan-Änderung `c581639f`, Report `4173bc00`, Fix `58a61cdf`)
- **Skill:** `reviewer.md` @ 1.20.0
- **Modell-ID:** claude-sonnet-5-5
- **Datum:** 2026-10-20
- **Eingangs-Kontext:** Runde 1 (`docs/reviews/2026-10-20-slice-276-stop-hook-streng.md`, HIGH `schalter-datei-ohne-halter`, INFO), Slice-Plan §Plan-Änderung nach Review-Runde 1, `MR-079`, `MR-076`, `SPEC-094`, `SPEC-103`, `AGENTS.md` §3.1/§3.6/§3.7/§3.8. Wegwerf-Proben in Scratch-Repos (gelöscht); `make selbstpruefung` einmal gefahren.

## Findings

### 1 — INFO
- **kategorie:** INFO
- **quelle:** `MR-079`, `AGENTS.md` §3.6
- **pfad:** `.claude/hooks/stop-require-gates.sh` · „if [ -e .harness/stop-gate-locker ]; then locker=1; fi"
- **befund:** `-e` ist wahr für ein Verzeichnis dieses Namens: im Scratch-Repo gab der Hook mit `mkdir .harness/stop-gate-locker` und ungedecktem Inhalt (gleicher HEAD-Stempel) frei, also lockerer Modus. Ein toter Symlink dagegen ist `-e` falsch und bleibt streng (gemessen). Ein Verzeichnis kann in diesem Repo nicht versehentlich entstehen (kein Werkzeug, Hook oder Make-Target legt etwas unter `.harness/` außer `state/` an, und `state/` ist ignoriert); die Lockerung bleibt ein absichtlicher Akt. Nicht dokumentiert ist nur, dass „Datei" auch Verzeichnis meint.
- **verifizierbar:** ja — `mkdir .harness/stop-gate-locker` im Wegwerf-Repo, Hook fahren.
- **klasse:** schalter-praedikat-weiter-als-benannt

### 2 — INFO
- **kategorie:** INFO
- **quelle:** `AGENTS.md` §3.6 (permanent ungeprüft), `MR-079` §Grenze
- **pfad:** `harness/conventions/MR-079-stop-hook-streng-mit-head-stempel.md` · „Wer `.harness/stop-gate-locker` anlegt, lockert das Gate, ohne dass ein Gate es meldet"
- **befund:** Die Grenze ist benannt und ehrlich. Ergänzend: Eine **unversionierte** Datei im Arbeitsbaum sieht auch `make selbstpruefung` nicht, weil sie den Klon von `HEAD` prüft; nur eine committete Datei bricht den Lauf. Die Richtung ist gegenüber Runde 1 umgekehrt und damit tragbar (Lockern braucht eine Tat, kein Verlust lockert).
- **verifizierbar:** ja
- **klasse:** grenze-nur-fuer-committete-datei

### 3 — LOW
- **kategorie:** LOW
- **quelle:** `AGENTS.md` §6 Schritt 4 (Plan-Änderung gehört in den Plan), §5 Regel 13
- **pfad:** `docs/plan/planning/in-progress/slice-276-stop-hook-haertung-aus-ai-harness-init-v0100.md` · „(`.harness/stop-gate-streng`): das Verhalten des Hooks bleibt das heutige"
- **befund:** Die Plan-Änderung ist als Abschnitt angehängt; §1 Ziel, §2 DoD-Punkt 1 („`.harness/stop-gate-streng` ist versioniert"), die Abgrenzungs-Prosa („löscht `.harness/stop-gate-streng` per eigener ADR"), die §3-Tabellenzeile und das §6-Risiko nennen die Datei weiter als Halter des strengen Modus. Der Plan widerspricht damit dem gelieferten Stand; wer DoD oder Ziel isoliert liest, bekommt die entfallene Datei.
- **verifizierbar:** ja — `grep -n "stop-gate-streng"` im Plan.
- **klasse:** plan-anhang-statt-korrektur

## Geprüft, ohne Befund

- **HIGH `schalter-datei-ohne-halter` geschlossen (Aussage, nicht Fundstelle).** Gefahren im Scratch-Repo mit Nachweis und passendem HEAD-Stempel, ungedeckter Änderung: ohne Datei, nach `rm`, toter Symlink, `STOP_GATE_STRENG=1` mit Datei ⇒ `block`; `STOP_GATE_STRENG=0`/`true` mit Datei ⇒ `approve` (Variable sagt nichts, wie dokumentiert). Löschung, Branch-Wechsel und `git stash` entfernen nur eine Datei, die es im Normalzustand nicht gibt; sie können nichts lockern. Nur Anlegen lockert.
- **Versehentliches Anlegen.** `grep` über Makefile, `*.mk`, `tools/`, `.githooks/`, `.claude/`, `.github/`: einziger Schreiber unter `.harness/` ist `record-gates.sh` (`.harness/state/`, per `.gitignore` ausgenommen); `git clean` löscht, legt nie an. Kein Befund.
- **Default gleich altem Hook.** Alt (`2032fe1d`) und neu, je 14 Lagen im Scratch-Repo: sauber, Änderung ohne Nachweis, gedeckt, Änderung nach Nachweis mit und ohne Commit, Stempel fehlt, Stempel leer, Stempel eines anderen Commits mit Änderung, `stop_hook_active`, unlesbarer Nachweis, fehlendes Hash-Skript, kaputter Index mit und ohne Nachweis: Entscheidung und Exit 0 byte-gleich. Dokumentierte Abweichungen stimmen: außerhalb eines Repos alt Exit 1, neu Exit 2; ohne `git` im `PATH` alt Exit 127, neu `block` mit Grund (beide gemessen).
- **Selbstprüfung (5).** `make selbstpruefung` grün, alle vier STOP-Zeilen. STRENG und LOCKER belegen sich gegenseitig (gleicher Stand, ungedeckte Änderung: blockt ohne Datei, frei mit Datei); VARIABLE belegt die Aufhebung; BLOCK trägt, was die Zeile sagt (neuer HEAD, abweichender Hash; die Zeile behauptet für den strengen Zweig nur den Hash, INFO 1 aus Runde 1 erledigt). Die Vorbedingung „Datei darf im Klon nicht liegen" ist sinnvoll gesetzt (committete Lockerung bricht den Lauf). Die letzte Prüfung nach `rm -f` ist für sich nicht trennscharf (der Hash weicht ohnehin ab), aber ohne Behauptung darüber und harmlos.
- **Spiegel.** `grep -rn "stop-gate-streng"` außer `docs/reviews` und `planning/done`: nur der Plan (Finding 3) und MR-079, dort als historische Begründung des Mangels. `harness/README.md`, `SPEC-094`, `SPEC-103`, `conventions.md`-Zeile, Hook-, `record-gates.sh`- und `selbstpruefung.sh`-Kommentare stimmen mit dem Code überein; `.claude/commands` und `Makefile` ohne Treffer. Die Datei `.harness/stop-gate-streng` ist gelöscht.
- **Hard Rules.** §3.1: keine Host-Interpreter in Hook/Skripten. §3.6: der Fix ist keine Lockerung; er verschärft die Ausfallrichtung, Default und Entscheidungen sind unverändert (gemessen); der Lockerungs-Pfad ist als ADR-pflichtig benannt. §3.7: Kommentare tragen Zusage/Kopplung/Grenze, Herkunft nur `MR-079`/`MR-076`/`seit slice-245`, keine Review-Historie im Code. §3.8: der Hook liest Datei, Variable, Stempel, `git`; die Zusage gilt für jede, Fehlerformen blocken (Trap, Exit 2).
- **Abgrenzung §1.** `git diff d1b49459..HEAD --stat` für `slice-mv.sh`, `extract-command.awk`, `pretooluse-command-guard.sh`: unverändert.

## Kategorie-Summary

HIGH 0 · MEDIUM 0 · LOW 1 · INFO 2

## Verdikt

Das HIGH aus Runde 1 ist geschlossen; die LOW/INFO-Punkte blockieren nicht. Freigabe zur Verifikation. Finding 3 vor Closure im Plan nachziehen (Verifier liest dort die DoD).
