# Review R1 — slice-260: Festlegungen der lokalen Wächter, Hooks und Prüfer in der Spezifikation

- **Review-Art:** Code (Doku- und Harness-Diff). Geprüft gegen den Slice-Plan (`slice-260`, §1
  Abgrenzung, §3 mit den drei Plan-Änderungs-Notizen, §8), gegen `MR-004`, `MR-005`, `MR-047`,
  den neuen Eintrag `MR-076`, `MR-025` (Spiegel), die Baseline `v6.17.0` ·
  `regelwerk/modul-13-quality-gates.md` §Guard-Härtung und
  `v6.17.0` · `regelwerk/grundlagen-durchsetzungsschicht.md` §Grenzen — ehrlich benannt sowie die
  Hard Rules `AGENTS.md` §3.4/§3.7/§3.8 und §5 Regel 13/15/16. Gegenstand ist die
  Maintainability; die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-260` · Range `b6a1c21e..e4505c7f` (zehn Commits; inhaltlich
  `ffedcb95` Übergangs-Wächter, `5bef299c` Gate-Nachweis im Rezept, `e4505c7f` Spezifikation
  §7/§8, Sensor-Dateien, `MR-076`).
- **Skill:** `reviewer.md` @ 1.18.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-08
- **Eingangs-Kontext:** Slice-Plan `slice-260`; Spezifikation §7 (`SPEC-093`..`SPEC-096`) und
  §8; die Träger im Code: `.claude/hooks/pretooluse-command-guard.sh`,
  `tools/harness/extract-command.awk`, `.claude/settings.json`, `Makefile` (`gates`,
  `record-gates`, `.NOTPARALLEL`, `blackbox-probe`, `hooks`, `verify-closure-notes`),
  `tools/harness/record-gates.sh`, `tools/harness/working-tree-hash.sh`,
  `.claude/hooks/stop-require-gates.sh`, `.githooks/{pre-commit,commit-msg}`,
  `.github/workflows/ci.yml`, `tools/blackbox-probe.sh`, `.d-check.closure.yml`. Vorherige
  Findings am selben Gegenstand: R1/R2 zu `slice-259` (Klasse
  `grenzen-liste-wird-als-vollstaendig-gelesen`, dreimal).
- **Proben (eigene Läufe, echte Ausgabe; alle in Wegwerf-Kopien im Scratchpad, Arbeitsbaum
  unverändert):**
  - Wegwerf-Klon von `e4505c7f`, zwei absichtlich kaputte Slices angelegt (offener DoD-Haken,
    Platzhalter `<…>` in der Closure-Notiz) — `done/slice-998-kaputt.md` und
    `done/wellenlos/slice-999-kaputt.md`; `make verify-closure-notes`: Exit 2,
    `d-check: 882 Datei(en) geprüft, 5 Befund(e)`, **alle fünf** auf `slice-998`, **null** auf
    `slice-999`.
  - Modell-Makefile (`.NOTPARALLEL`, Glied `a` rot, Rezept von `gates` schreibt `NACHWEIS`):
    `make gates` Exit 2 ohne Nachweis; `make -k gates` Exit 2 ohne Nachweis (`Das Ziel „gates“
    wurde wegen Fehlern nicht neugemacht`); `make -k -j4 gates` ebenso; **`make -i gates` Exit 0
    mit `NACHWEIS`**; `MAKEFLAGS=i make gates` ebenso.
  - Erkennungs-Pipeline des Übergangs-Wächters unter `set -euo pipefail`, Eingabe: ein
    `done/wellenlos/slice-900-…`-Treffer in Zeile 1, dann `n` weitere Stub-Zeilen unter
    `done/welle-99/`, je fünf Läufe: `n=0` 5/5 erkannt, `n=10` 5/5, `n=100` 2/5, `n=1000` 0/5,
    `n=5000` 0/5.
  - Stop-Hook in einem Wegwerf-Repo mit gültigem Nachweis: sauber → `{"decision":"approve"}`,
    Exit 0; mit einer unlesbaren, nicht ignorierten Datei (`chmod 000`) →
    `sha256sum: geheim.txt: Keine Berechtigung`, **Exit 1, keine JSON-Antwort**.
  - Tool-Call-Wächter direkt mit JSON-Eingaben: `timeout 5 pip …`, `if true; then pip …; fi`,
    `env -i pip …`, `! pip …` → je Exit 0, kein `deny`; `{"tool_input":{"command":"ls"} x}` →
    Exit 2, `deny`.
  - `tools/blackbox-probe.sh` mit `IMAGE=gibtsnicht-xyz REF=HEAD`: `PROBE_FORMS=` (leer) →
    Abbruch erst bei `Nachher-Image … fehlt`; `PROBE_FORMS=" "` → `PROBE_FORMS ist leer`. Beide
    Exit 2.
  - Gezählt: Sperrliste im Wächter 26 Wörter; `permissions.deny` 31 Einträge (26 + `python`,
    `python3` + drei `git push --force`-Formen); `fail`-Aufrufe in `tools/blackbox-probe.sh` 11;
    Closure-Moves seit 2026-09-29: 20 (17 `done/wellenlos/`, 3 `done/welle-91/`), keiner direkt
    unter `done/`.

## Findings

### F-1 — HIGH: Der Übergangs-Wächter löst jetzt für Slices unter `done/`-Unterverzeichnissen aus, aber `verify-closure-notes` liest diese Slices nicht

- **kategorie:** HIGH
- **quelle:** Prüffrage 1 (Stilles-Grün-Pfad in einem Gate); `AGENTS.md` §3.8; `MR-025`
- **pfad:** `.githooks/pre-commit` · „haengen am UEBERGANG selbst" und „er prueft dann den
  ganzen Bestand, der ohnehin gruen sein muss"; `.d-check.closure.yml` ·
  `files: "docs/plan/planning/done/slice-*.md"` und `closure:` → `dir: docs/plan/planning/done`;
  `harness/sensors/hooks.md` · „die Vorbedingungen hängen damit am **Übergang** selbst"
- **befund:** Das Profil von `verify-closure-notes` listet `planning.closure.dir` nicht
  rekursiv, und alle `structure`-Regeln haben einen Glob, dessen `*` keinen `/` überspringt.
  Damit prüft der ausgelöste Lauf den Slice nicht, dessen Übergang ihn ausgelöst hat:
  `done/wellenlos/slice-999` mit offenem DoD-Haken und Platzhalter ergibt null Befunde, dieselbe
  Datei direkt unter `done/` fünf. Seit 2026-09-29 sind das alle 20 Closures, und der neue
  Kommentar, die Sensor-Datei und der Plan sagen jetzt, die Vorbedingungen hingen am Übergang
  bzw. der Lauf prüfe den ganzen Bestand. Der Fix wurde nur auf die Erkennung gefahren (fünf
  Commits, alt/neu), nicht darauf, was der ausgelöste Lauf dann liest.
- **verifizierbar:** ja — `make verify-closure-notes` in einem Klon mit einem kaputten Slice unter
  `done/wellenlos/` (siehe Proben).
- **klasse:** `ausloeser-und-pruefmenge-auseinander`

### F-2 — MEDIUM: Die Erkennungs-Pipeline fällt unter `pipefail` bei großen Ranges still aus, und das erweiterte Muster macht den frühen Treffer zum Regelfall

- **kategorie:** MEDIUM
- **quelle:** Prüffrage 19 (Fehlerformen eines Lesewegs), Prüffrage 1
- **pfad:** `.github/workflows/ci.yml` · „`| awk -F'\t' '{print $NF}' | grep -qE`" unter
  `set -euo pipefail`; dieselbe Zeile in `.githooks/pre-commit`
- **befund:** `grep -q` beendet sich beim ersten Treffer. Schreibt `awk` danach weiter, endet es
  durch SIGPIPE, `pipefail` meldet die Pipeline als gescheitert, und das `if` überspringt
  `verify-closure-notes`. Gemessen: ein Treffer in Zeile 1 und 100 weitere Zeilen werden in 2 von
  5 Läufen erkannt, ab 1000 Zeilen nie. Die Pipeline gab es vorher schon. Mit dem neuen Muster
  treffen aber auch die Stub-Zeilen unter `done/<welle-id>/` zu, die vorher nie trafen; ein
  CI-Range mit einer Archivierung und einer Closure trifft also früh und verliert das Ergebnis.
  Gefahren wurde der Fix nur auf Ranges mit wenigen Zeilen.
- **verifizierbar:** ja — Probe-Skript aus den Proben; eine Range mit ≥ 100 A/R-Einträgen unter
  `done/` im CI-Schritt.
- **klasse:** `pipefail-grep-q-sigpipe`

### F-3 — MEDIUM: `SPEC-094` und `MR-076` sagen „nur, wenn alle Glieder grün sind"; unter `make -i` entsteht der Nachweis trotz rotem Glied

- **kategorie:** MEDIUM
- **quelle:** Prüffrage 8 (Botschaft verallgemeinert über die Messung), Prüffrage 18; `AGENTS.md`
  §5 Regel 15
- **pfad:** `spec/spezifikation.md` · `SPEC-094` „und damit nur, wenn alle Glieder grün sind —
  auch unter `make -k`"; `harness/conventions/MR-076-gate-nachweis-im-rezept.md` · „Ein Rezept
  läuft erst, wenn alle Prerequisites grün sind"; `Makefile` · „ein Rezept läuft erst, wenn alle
  Prerequisites grün sind"
- **befund:** Gemessen wurde `-k`; die Aussage ist allgemein gefasst. Unter `make -i gates`
  (oder `MAKEFLAGS=i`) gilt ein rotes Glied als erledigt, das Rezept läuft und schreibt den
  Nachweis, Exit 0 — der Stop-Hook gibt danach frei. Die `Grenze:`-Zeile von `MR-076` nennt nur den
  Handaufruf von `record-gates`. Das ist die N+1-te Form zu der gemessenen.
- **verifizierbar:** ja — Modell-Makefile aus den Proben; im Repo `make -i gates` mit einem roten
  Glied, danach `.harness/state/gates-passed.diffsha` vorhanden.
- **klasse:** `botschaft-ueberdehnt-gemessene-form`

### F-4 — MEDIUM: `SPEC-093` nennt als einzige Grenze die quote-blinde Falsch-Positiv-Seite; die Durchlass-Lücken, die der Wächter selbst aufzählt, fehlen

- **kategorie:** MEDIUM
- **quelle:** Prüffrage 18 (Grenzen-Liste ohne ihre größte Lücke);
  `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`
- **pfad:** `spec/spezifikation.md` · `SPEC-093` „Quote-blind: ein Trenner innerhalb eines
  Arguments beginnt ein neues Segment"; `.claude/hooks/pretooluse-command-guard.sh` · „GRENZE,
  Umfang — gemessen, nicht geschätzt. Ungeprüft bleiben:"
- **befund:** Die Festlegung sagt „blockiert, wenn … an Befehlsposition" und nennt danach nur
  eine Grenze, die zu viel blockt. Was durchgelassen wird, steht im Kommentar des Wächters und
  fehlt in `SPEC-093`: Segment-Köpfe mit Shell-Schlüsselwörtern (`if …; then pip`, `! pip`),
  Wrapper außerhalb der Präfixliste (`timeout 5 pip`), wortinterne Splices und verschachtelte
  Escapes. Dazu kommt eine Form, die auch der Kommentar nicht nennt: ein Flag hinter einem
  übersprungenen Präfix wird selbst zum Kopf (`env -i pip`, gemessen Exit 0). Wer die
  Festlegung liest, hält diese Aufrufe für geblockt. Wer den Wächter „an die Spezifikation
  angleicht", hat keinen Text, an dem die Lücken auffallen.
- **verifizierbar:** ja — die vier JSON-Proben aus dem Proben-Abschnitt gegen den Wächter.
- **klasse:** `grenzen-liste-wird-als-vollstaendig-gelesen`

### F-5 — MEDIUM: `SPEC-094` lässt aus, dass der Stop-Hook fail-open endet, wenn der Hash scheitert

- **kategorie:** MEDIUM
- **quelle:** Prüffrage 18; Baseline `v6.17.0` · `regelwerk/grundlagen-durchsetzungsschicht.md`
  §Vier Design-Eigenschaften · „fail-closed"
- **pfad:** `spec/spezifikation.md` · `SPEC-094` „blockt, wenn er abweicht oder kein Nachweis
  existiert, aber der Arbeitsbaum Änderungen trägt"; `.claude/hooks/stop-require-gates.sh` ·
  `current="$(bash tools/harness/working-tree-hash.sh)"`
- **befund:** Scheitert `working-tree-hash.sh` (z. B. an einer unlesbaren, nicht ignorierten
  Datei), bricht der Hook unter `set -e` mit Exit 1 ab, ohne JSON-Antwort (gemessen). Ein
  Stop-Hook mit Exit ≠ 0/2 blockt nicht, der Stop geht also durch. Diesen Fall schreibt
  `SPEC-094` nicht auf. Die Freigabe-Fälle (frischer Klon, `stop_hook_active`, außerhalb von git)
  sind aufgezählt und lesen sich deshalb als vollständige Menge. Szenario: ein Docker-Lauf
  hinterlässt eine root-eigene Datei im Baum; `make gates` endet rot (auch `record-gates`
  scheitert), und der Stop-Hook lässt den Handoff trotzdem durch.
- **verifizierbar:** ja — Wegwerf-Repo mit `chmod 000`-Datei, Hook-Exit 1 ohne Ausgabe.
- **klasse:** `grenzen-liste-wird-als-vollstaendig-gelesen`

### F-6 — MEDIUM: `SPEC-096` legt „`PROBE_FORMS` ist leer ⇒ Exit 2" fest; leer heißt im Code Default

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §5 Regel 13 (Grenze gegen den Gegenstand prüfen); Prüffrage 10
- **pfad:** `spec/spezifikation.md` · `SPEC-096` „`PROBE_FORMS` ist leer"; `tools/blackbox-probe.sh`
  · `PROBE_FORMS="${PROBE_FORMS:-- --json --yaml --doctor}"`
- **befund:** Ist `PROBE_FORMS` leer gesetzt, ersetzt `:-` den Wert durch den Default, und der Lauf
  geht weiter (gemessen). Exit 2 gibt es nur, wenn der Wert aus Leerraum besteht. Die alte
  Sensor-Datei sagte richtig „leer heißt Default", und genau dieser Satz ist mit der Umstellung
  auf die Kennung weggefallen. Jetzt steht die gegenteilige Aussage im verbindlichen Stratum.
  Szenario: Wer den Code an die Festlegung angleicht, macht aus dem dokumentierten Default einen
  Abbruch.
- **verifizierbar:** ja — `PROBE_FORMS= REF=HEAD IMAGE=<fehlt> bash tools/blackbox-probe.sh` bricht
  erst am Nachher-Image ab.
- **klasse:** `grenze-gegen-beschreibung-statt-gegenstand-geprueft`

### F-7 — MEDIUM: Spiegel nicht nachgezogen — die Sensor-Datei von `verify-closure-notes` sagt weiter „nicht rekursiv"

- **kategorie:** MEDIUM
- **quelle:** `MR-025`; `BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`
- **pfad:** `harness/sensors/verify-closure-notes.md` · „sobald ein Rename/Add nach
  `docs/plan/planning/done/slice-*.md` gestagt ist (nicht rekursiv)"
- **befund:** Die Erkennung erfasst jetzt jede Tiefe unter `done/`. Die Bindung-Sektion der
  Sensor-Datei des ausgelösten Targets beschreibt noch das alte Verhalten, und der Plan §3 führt
  diese Datei nicht als Spiegel. Szenario: Wer den Sensor nachschlägt, schließt daraus, eine
  Closure nach `done/wellenlos/` löse keinen Lauf aus. Laut Plan §8 steht die Beobachtung bei 2×
  („ein dritter Treffer wäre eine Lücke"); das hier ist der dritte Treffer.
- **verifizierbar:** ja — `grep -n "nicht rekursiv" harness/sensors/verify-closure-notes.md`.
- **klasse:** `guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`

### F-8 — LOW: Fail-closed-Aufzählung in `SPEC-093` lässt „Müll außerhalb eines Strings" aus

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §5 Regel 13; `MR-025`
- **pfad:** `spec/spezifikation.md` · `SPEC-093` „(kein Objekt, abgeschnitten, ein `\u`-Escape,
  zwei Strings ohne Trenner)"
- **befund:** Der Extraktor blockt auch auf jedes Zeichen außerhalb eines Strings, das keine
  JSON-Struktur ist (gemessen: Exit 2, `deny`). Die alte `guard-probe.md` nannte den Fall, die
  neue nicht mehr, und die Spezifikation hat ihn nicht übernommen. `guard-probe.sh` trägt eine
  Probe dafür, deshalb nur LOW.
- **verifizierbar:** ja — `make guard-probe` (Probe „Müll ausserhalb String").
- **klasse:** `grenzen-liste-wird-als-vollstaendig-gelesen`

### F-9 — LOW: Der umgeschriebene `.NOTPARALLEL`-Kommentar zeigt auf `MR-005`, das die Aussage nicht trägt

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §5 Regel 16; `AGENTS.md` §3.7 (Rang-Zeiger)
- **pfad:** `Makefile` · „Prerequisites laufen nacheinander, in der Reihenfolge ihrer Liste, auch
  unter `make -j` (MR-005)."
- **befund:** In `MR-005` stehen Inhalts-Hash und Sub-Shell-Rekursion, aber weder `.NOTPARALLEL`
  noch `-j`. Der ursprüngliche Grund für die Sequenz (Nachweis nach den Gliedern) ist mit
  `MR-076` weggefallen. Der Kommentar sagt jetzt *was*, aber der Zeiger führt zu keinem *warum*.
  Wer prüft, ob die Zeile noch nötig ist, findet keine Quelle.
- **verifizierbar:** nein (Urteil).
- **klasse:** `rang-zeiger-nach-autoritaets-umzug-nicht-nachgezogen`

### F-10 — LOW: Plan §8 nennt `HARN` weiter als nicht berührt

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §5 Regel 12; `MR-031`
- **pfad:** Plan `slice-260` §8 · „berührt ist `HARN` damit nicht"
- **befund:** Mit den Plan-Änderungen ändert der Slice `tools/harness/record-gates.sh` (Sub-Area
  `HARN`), dazu `Makefile`, `.githooks/` und `ci.yml`. Die Sub-Area-Prüfung wurde nicht
  nachgezogen, obwohl §8 selbst sagt, eine solche Änderung sei eine Plan-Änderung. Folgenlos für
  den Modus, weil beide Sub-Areas GF sind.
- **verifizierbar:** nein (Urteil).
- **klasse:** `vorpruefung-nach-planaenderung-nicht-nachgezogen`

### F-11 — INFO: Exit-Zuordnung von `make blackbox-probe` über `make`

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** `spec/spezifikation.md` · `SPEC-096` „über `make` endet jedes Scheitern mit 2";
  `Makefile` · „1 = Abweichung, 2 = gescheitert"
- **befund:** Über `make` endet auch eine Abweichung (Skript-Exit 1) mit 2. Die Spezifikation sagt
  „jedes Scheitern", der Hilfetext des Targets nennt die Skript-Codes, als gälten sie für das
  Target. Bestand und nicht falsch, aber missverständlich.
- **verifizierbar:** ja — `make blackbox-probe` mit einer Abweichung endet mit 2.
- **klasse:** `make-normalisiert-exit`

## Negativbefunde (geprüft, ohne Befund)

- **`SPEC-093` Mechanik:** 26 Wörter der `BLOCKED`-Liste nachgezählt, `python`-Muster mit
  Versions-Suffix, die Trenner (`;` `&` `&&` `|` `||` `$(` Backtick `(`, Zeilenenden inkl. `\r`),
  die acht Präfixe, Zuweisungen und `{`/`}`, Basisname eines Pfads, die fünf Shells, `-c` im
  Flag-Bündel, Tiefe > 3 blockt, die zwei Kanäle (`deny` + Exit 2) und der stumme Durchlass mit
  Exit 0. Alles stimmt mit dem Code (abgesehen von F-4 und F-8).
- **`SPEC-093` zweite Schicht:** die 26 Wörter und `python`/`python3` stehen vollständig in
  `permissions.deny`. Die drei `git push --force`-Einträge gehören zur git-Hälfte aus `MR-047`
  und nicht zu §3.1, sie fehlen hier zu Recht.
- **`SPEC-094` Hash:** getrackte und nicht ignorierte ungetrackte Dateien, `sort -zu`, `LINK`
  mit Ziel, `GONE` für fehlende getrackte Dateien. Freigabe-Fälle (sauber ohne Nachweis,
  `stop_hook_active`, außerhalb von git) und `decision: block` mit Grund stimmen mit dem Hook
  (abgesehen von F-5).
- **`SPEC-095`:** `make hooks` setzt `core.hooksPath`; `commit-msg` ruft
  `make trace-check MSGFILE`; `pre-commit` ruft `adr-check STAGED=1` und `doc-check` unter
  `set -e`; die CI fährt dieselbe Erkennung über die Range. Stimmt mit dem Code.
- **`SPEC-096`:** die elf Abbrüche des Skripts einzeln den neun Klauseln der Festlegung
  zugeordnet, alle elf gedeckt; Kanarien (`sauber` 0, `links` 1, vorher und nachher, beide
  Images), die getrennten Vergleiche und `VERSION=0.0.0-dev` stimmen (abgesehen von F-6).
- **`make -k` und `-j`:** Der Fix hält, was er behaupten will. Unter `-k` und unter `-k -j4` läuft
  das Rezept nach einem roten Glied nicht.
- **Kommentare (§3.7):** Neue und geänderte Kommentare in `.githooks/pre-commit`, `ci.yml`,
  `Makefile` (`gates`) und `record-gates.sh` tragen Zusage, Kopplung, Grenze oder Rang-Zeiger;
  die Slice-/Wellen-Nummern aus den berührten Kommentaren sind entfernt, neue sind nicht
  hinzugekommen. Inhaltlich falsche Zusagen stehen unter F-1/F-3, die Form ist in Ordnung.
- **Spiegel von „letzter Prerequisite":** `MR-004` bleibt unverändert und wird von `MR-076`
  geschärft, wie es Baseline `v6.17.0` · `regelwerk/modul-13-quality-gates.md` §Guard-Härtung
  verlangt. Die `record-gates`-Kopfzeile ist nachgezogen; `harness/README.md` („`record-gates`
  als letzter Schritt") und `.claude/commands/plan-welle.md` („endet mit `record-gates`") stimmen
  weiter.
- **`MR-076` Form und Index:** Pflichtfelder vorhanden, kein Status-Feld, `Grenze:` als eigene
  Zeile (zulässig nach §Guard-Härtung), Nummer dicht nach `MR-075`, Index-Zeile mit Anker.
- **Referenz-Richtung (§3.4):** `SPEC-093`..`SPEC-096` und die §8-Zeile nennen weder ADR noch
  Slice, Welle oder Commit-Hash.
- **Kennungs-Vergabe (Prüffrage 16):** `SPEC-093`..`SPEC-096` dicht nach `SPEC-092`; neue
  Verweise nennen die Kennung und verlinken den Abschnitt.
- **Zählungen in den Plan-Notizen und Botschaften:** 20 Closures (17/3) seit 2026-09-29,
  26 Wörter, elf Abbrüche — alle nachgezählt und richtig.
- **Zustandsfelder (Prüffrage 7):** Die Roadmap verliert nur den Ruhe-Marker (Beanspruchung),
  keine Chronik.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 1 | F-1 |
| MEDIUM | 6 | F-2, F-3, F-4, F-5, F-6, F-7 |
| LOW | 3 | F-8, F-9, F-10 |
| INFO | 1 | F-11 |

Wiederkehrende Klasse: `grenzen-liste-wird-als-vollstaendig-gelesen` (F-4, F-5, F-8) — wie in
`slice-259` R1 wurde die Festlegung am Code geprüft, aber nur auf der Seite, die der Autor als
Prüfgegenstand kannte (was blockt, was freigibt). Die Fehlerformen (was durchrutscht, was beim
Scheitern passiert) wurden nicht gefahren. Das ist der dritte Lauf in Folge mit dieser Klasse,
ein Steering-Loop-Signal.

## Verdikt

**Blockiert.** F-1 ist ein Stilles-Grün-Pfad: Der Fix stellt die Auslösung wieder her, aber der
ausgelöste Lauf liest die Datei nicht, um deren Übergang es geht. Kommentar, Sensor-Datei und
Plan versprechen jetzt eine Bindung, die es seit 2026-09-29 nicht gibt. Die Zusage ist von
„schweigt" auf „sagt grün" umgeschlagen. F-2 bis F-7 sind vor dem Merge zu klären; F-6 und F-7
sind reine Text-Nachzüge, F-3 und F-5 brauchen eine Entscheidung, ob nachgezogen oder benannt
wird.
