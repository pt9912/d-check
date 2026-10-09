# Review R2 — slice-260: Festlegungen der lokalen Wächter, Hooks und Prüfer in der Spezifikation

- **Review-Art:** Code (Doku- und Harness-Diff), Folgerunde. Geprüft wird, ob die R1-Befunde F-1 bis
  F-11 in der Sache geschlossen sind und ob der Fix neue Fehler einführt. Maßstab sind der Slice-Plan
  (`slice-260`, §3 mit der Plan-Änderung nach R1, §8), `MR-004`, `MR-076`, `MR-025`, die Baseline
  `v6.17.0` · `regelwerk/grundlagen-durchsetzungsschicht.md` §Vier Design-Eigenschaften und
  `v6.17.0` · `regelwerk/modul-03-spec.md` §Ziel-Form: Spezifikation sowie die Hard Rules
  `AGENTS.md` §3.4/§3.7/§3.8 und §5 Regel 13/15/16. Die DoD-Abhakung gehört nicht zu diesem Review.
- **Gegenstand:** `slice-260` · Plan-Commit `01d7c500` und Fix-Commit `596630db`
  (`01d7c500^..596630db`).
- **Skill:** `reviewer.md` @ 1.18.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** R1-Report (`2026-10-08-…-r1.md`, Findings F-1 bis F-11); Slice-Plan
  `slice-260`; Folge-Slice `slice-263` (open); Spezifikation §7 (`SPEC-093` bis `SPEC-096`) und §8.
  Im Code gelesen: `Makefile` (`gates`, `ci`, `fullbuild`, `.NOTPARALLEL`, `blackbox-probe`),
  `a-check.mk`, `.claude/hooks/stop-require-gates.sh`, `tools/harness/{working-tree-hash,record-gates}.sh`,
  `.githooks/pre-commit`, `.github/workflows/{ci,release}.yml`, `.claude/hooks/pretooluse-command-guard.sh`,
  `tools/harness/extract-command.awk`, `tools/harness/selbstpruefung.sh`, `harness/sensors/{hooks,verify-closure-notes,blackbox-probe}.md`,
  `harness/conventions/MR-076-gate-nachweis-im-rezept.md`.
- **Proben** (eigene Läufe mit echter Ausgabe, alle in Wegwerf-Kopien im Scratchpad; der Arbeitsbaum
  blieb unverändert):
  - **`gates`-Rezept, Modell-Makefile** mit der wörtlichen `case`-Zeile, Glied grün bzw. rot, je
    36 Aufrufformen. Gefahren mit GNU Make 4.3 (Host) und 4.4.1 (`golang:1.27.1` in Docker). Formen:
    ohne Flag, `-i`, `-k`, `-ik`, `--ignore-errors`, `-j2 -i`, `-j 4`, `-e`, `-e -i`, `--debug=b`
    (mit und ohne `-i`), `--trace`, `-I inc` (mit und ohne `-i`), `-C .`, `--no-print-directory`,
    `--warn-undefined-variables` (mit und ohne `-i`), `-Otarget`, `--output-sync=line -j2`, `-l 3`,
    `--eval=X:=1`, `-r -R`, `-s`, eine Variable mit `i` im Wert, `-- i=1`, ein Sub-Make über `$(MAKE)`
    (MAKELEVEL 1, mit und ohne `-i`, mit und ohne `--no-print-directory`), dazu `MAKEFLAGS` aus der
    Umgebung (`i`, `-i`, `--ignore-errors`, `' -i'`, `'s -- X=1'` mit `-i`), `GNUMAKEFLAGS=-i` und
    `MFLAGS=-i`. **Ergebnis auf beiden Versionen gleich:** Jede Form mit `-i` schreibt keinen Nachweis
    und gibt kein `[gates] … green` aus. Jede Form ohne `-i` und mit grünem Glied schreibt den
    Nachweis. Jede Form ohne `-i` und mit rotem Glied endet mit Exit 2 ohne Nachweis.
    `MFLAGS=-i` wirkt nicht als Flag, und die Erkennung verhält sich dazu richtig. Ein falsch positiver
    Treffer kam nicht vor. Die erste Wortgruppe war ohne Einzelbuchstaben-Flags immer `-`; lange
    Optionen erscheinen nie im ersten Wort.
  - **Umgehung über die Kommandozeile:** `make -i gates MAKEFLAGS=` und `make -i gates MAKEFLAGS=k`
    ergeben auf 4.3 und 4.4.1 mit rotem Glied den Nachweis, `[gates] … green` und Exit 0.
    `--eval='MAKEFLAGS:='` wird dagegen erkannt.
  - **Ausgabe ohne Nachweis:** Unter `make -n gates` und `make --trace -i gates` steht der Text
    `[gates] … green` als Echo der Rezeptzeile in der Ausgabe, ausgeführt wird er nicht. `make -i ci`
    im Modell gibt `gates: kein Nachweis …` aus, danach `[ci] gates + image-test green`, Exit 0.
  - **Übergangs-Erkennung (F-2)** unter `set -euo pipefail`, ein Treffer in Zeile 1 und 20000
    Stub-Zeilen: alte Form 0/10, neue 10/10, wie in der Botschaft. Fehlerform der Quelle (Ausgabe,
    dann Exit 128): die neue Form erkennt den Übergang nicht. Leere Quelle: nicht erkannt, richtig.
  - **Stop-Hook, Wegwerf-Repo mit den echten Skripten.** Ohne Nachweis und sauber: approve. Ohne
    Nachweis mit Änderung: block. Passender Nachweis: approve. Geänderter Inhalt: block. Unlesbare
    Datei (`chmod 000`): block mit dem neuen Grund. Unlesbarer Nachweis: block. Leerer Nachweis:
    block. Fehlendes Hash-Skript: block. Nachweis als Verzeichnis: block. Jeder Lauf endete mit Exit 0
    und gültigem JSON. **Kaputter `.git/index`, Änderung vorhanden, kein Nachweis:** `git status`
    endet mit Exit 128, und der Hook gibt `{"decision":"approve"}` aus.
  - **`SPEC-093`, Durchlass-Klassen gegen den Wächter (JSON per `awk` maskiert):** `then`, `do`, `!`,
    `timeout 5`, `nohup`, `env -i`, `sudo -u root`, `find -exec`, `p"i"p` und
    `bash -c "bash -c \"…\""` ergeben je Exit 0 ohne `deny`; ein Interpreter außerhalb der Liste
    (`php`) ebenfalls. Die Kontrollen `pip …` und `ls; pip …` ergeben Exit 2 mit `deny`. Ein
    `awk`-Programm mit `system("pip …")` wird geblockt, weil `(` als Trenner zählt; das deckt sich mit
    dem Satz der Festlegung („das nicht an der so bestimmten Befehlsposition steht").
  - **Extraktor außerhalb von Strings:** `,"a":fuse}` → Exit 0, `,"a":1e+-.}` → Exit 0, `} x}` →
    Exit 2, `,"a":true}` → Exit 0.

## Status der R1-Befunde

| R1 | Kategorie | Stand nach `596630db` |
|---|---|---|
| F-1 | HIGH | **teilweise geschlossen.** Die Behebung ist an `slice-263` adressiert (Auftraggeber-Entscheid, Plan-Änderung vor dem Code). `SPEC-095`, `hooks.md` Grenze 4, der `pre-commit`-Kommentar und die Bindung-Sektion von `verify-closure-notes.md` sind ehrlich geworden. Vertrag und Grenzen-Liste von `verify-closure-notes.md` sowie der Vertrag von `hooks.md` versprechen weiter das Gegenteil: siehe R2-F-1 |
| F-2 | MEDIUM | geschlossen (nachgefahren; zur Fehlerform der Quelle siehe Negativbefunde) |
| F-3 | MEDIUM | in der Sache geschlossen; Rest als INFO R2-F-6 |
| F-4 | MEDIUM | geschlossen (14 Formen gegen den Wächter gefahren) |
| F-5 | MEDIUM | für den Hash-Zweig geschlossen; der Schwester-Zweig ist offen, siehe R2-F-2 |
| F-6 | MEDIUM | geschlossen (`SPEC-096` und `blackbox-probe.md` stimmen mit `${PROBE_FORMS:-…}` und der Leerraum-Prüfung überein) |
| F-7 | MEDIUM | Text geschlossen; die Register-Folge steht offen, siehe R2-F-8 |
| F-8 | LOW | Fall genannt; die Formulierung ist zu weit, siehe R2-F-5 |
| F-9 | LOW | geschlossen (der Kommentar sagt eine Zusage, ohne falschen Zeiger) |
| F-10 | LOW | geschlossen (§8-Nachtrag nennt `HARN` als berührt, GF) |
| F-11 | INFO | geschlossen (die Hilfe-Zeile nennt Skript-Codes und die Normalisierung durch `make`) |

## Findings

### R2-F-1 — MEDIUM: Vertrag und Grenzen-Liste von `verify-closure-notes.md` und der Vertrag von `hooks.md` versprechen weiter die Prüfung der Slices, die der Lauf nicht liest

- **kategorie:** MEDIUM
- **quelle:** Prüffrage 18 (Grenzen-Liste ohne ihre größte Lücke); Prüffrage 8 / `AGENTS.md` §5
  Regel 15 (Botschaft); `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`
- **pfad:** `harness/sensors/verify-closure-notes.md` · „Kein Slice liegt in `done/`, dessen
  Closure-Notiz fehlt" und die Liste unter „## Grenze — was das Grün nicht abdeckt" (Punkte 1–7);
  `harness/sensors/hooks.md` · „die Vorbedingungen hängen damit am **Übergang** selbst"; Botschaft
  `596630db` · „Die Aussagen sagen das jetzt"
- **befund:** Die Lücke aus R1 F-1 steht in `verify-closure-notes.md` nur in der Bindung-Sektion. Der
  Vertrag sagt weiter „Kein Slice liegt in `done/` …, kein notiertes Risiko steht ohne einen der drei
  Kanon-Ausgänge da". Die sieben Punkte der Grenzen-Liste nennen die eine Lücke nicht, die alle
  Closures seit 2026-09-29 betrifft. Laut der Messung in `slice-263` §1 tragen vier davon
  (`slice-240` bis `slice-243`) heute einen Ausgang außerhalb des Wortschatzes, der Vertragssatz ist
  also schon jetzt falsch. Im Vertrag von `hooks.md` steht weiter der Satz, den R1 als Fundstelle
  zitiert hat, und Grenze 4 derselben Datei widerspricht ihm. Szenario: Wer den Sensor über Vertrag
  und Grenzen-Liste nachschlägt, wie der Abschnitt es anbietet, hält eine Closure nach
  `done/wellenlos/` für geprüft. Die Botschaft meldet die Aussagen als nachgezogen.
- **verifizierbar:** ja — `grep -n "Kein Slice liegt in" harness/sensors/verify-closure-notes.md`;
  ein kaputter Slice unter `done/wellenlos/` und `make verify-closure-notes` (R1-Probe).
- **klasse:** `grenzen-liste-wird-als-vollstaendig-gelesen`

### R2-F-2 — MEDIUM: Der Stop-Hook gibt in einem git-Repo frei, wenn `git status` scheitert und kein Nachweis existiert; `SPEC-094` nennt nur „außerhalb eines git-Repositorys"

- **kategorie:** MEDIUM
- **quelle:** Prüffrage 19 (die Härtung wurde nur auf den Fall gefahren, den sie behebt); Prüffrage 18;
  Baseline `v6.17.0` · `regelwerk/grundlagen-durchsetzungsschicht.md` §Vier Design-Eigenschaften ·
  „fail-closed"
- **pfad:** `.claude/hooks/stop-require-gates.sh` · `if [ -z "$(git status --porcelain=v1)" ]; then`;
  `spec/spezifikation.md` · `SPEC-094` „Außerhalb eines git-Repositorys gibt er frei: `git status`
  liefert dort nichts"
- **befund:** Der Fix schließt den Hash-Zweig mit `if ! current=…`. Der Zweig ohne Nachweis liest
  weiter `$(git status …)` ohne Prüfung des Exit-Codes, und jeder Fehler sieht dort aus wie ein
  sauberer Baum. Gemessen: In einem Repo mit kaputtem Index und einer Änderung endet `git status` mit
  Exit 128, der Hook gibt `approve` aus. `SPEC-094` beschreibt das Freigeben nur für den Fall
  außerhalb eines Repos. Die Sätze „blockt, wenn … kein Nachweis existiert, aber der Arbeitsbaum
  Änderungen trägt" und „Lässt sich der Hash nicht berechnen …, blockt er ebenso" lesen sich
  zusammen als vollständige fail-closed-Menge. Szenario: Der Nachweis fehlt (frischer Worktree
  mit Änderungen, `.harness/state/` entfernt) und ein Index-Schaden oder eine Rechte-Lage lässt `git status` scheitern.
  Der Handoff geht dann ohne Gate-Lauf durch. Die Vorbedingung ist selten, deshalb MEDIUM und nicht
  HIGH.
- **verifizierbar:** ja — Wegwerf-Repo, `.harness/state/` leer, eine Änderung, `printf x > .git/index`,
  `echo '{}' | bash .claude/hooks/stop-require-gates.sh` → `approve`.
- **klasse:** `haertung-kippt-fehlerpolitik-ungeprueft`

### R2-F-3 — LOW: Eine bereits committete Historie-Zeile der Spezifikation wurde rückwirkend umgeschrieben

- **kategorie:** LOW
- **quelle:** Baseline `v6.17.0` · `regelwerk/modul-03-spec.md` §Ziel-Form: Spezifikation · „Eine
  Historie-Zeile ist ein Protokoll und wird nicht rückwirkend geändert"
- **pfad:** `spec/spezifikation.md` §8 · Zeile `| 2026-10-08 | §7 um die lokalen Wächter, Hooks und
  Prüfer erweitert: …` mit „unter `make -k` und `make -i` trotz rotem Glied, der Stop-Hook endete ohne
  Antwort"
- **befund:** Die Zeile aus `e4505c7f` (2026-10-08) trägt jetzt Inhalt aus `596630db` (2026-10-09):
  `make -i`, den Stop-Hook und die Grenze aus `SPEC-095`. Die Änderung vom 2026-10-09 hat keine
  eigene Zeile, und das Datum der letzten Änderung, das laut Ziel-Form die letzte Historie-Zeile
  trägt, ist dadurch falsch. Der Nachzug nach Review in `slice-259` legte dafür eine neue Zeile an.
- **verifizierbar:** ja — `git show 596630db -- spec/spezifikation.md` zeigt `-`/`+` an der
  2026-10-08-Zeile.
- **klasse:** `historie-zeile-rueckwirkend-geaendert`

### R2-F-4 — LOW: Die vier Zwischen-Aussagen zur Lücke haben keine Adresse, die sie zurücknimmt

- **kategorie:** LOW
- **quelle:** `MR-025` (Spiegel vor dem Editieren auflisten);
  `BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`
- **pfad:** `spec/spezifikation.md` · `SPEC-095` „Der ausgelöste Lauf prüft nur die Slices direkt
  unter `done/`"; `harness/sensors/hooks.md` · Grenze 4; `harness/sensors/verify-closure-notes.md` ·
  „Der Lauf selbst liest nur `done/`"; `.githooks/pre-commit` · „GRENZE: der ausgeloeste Lauf prueft
  nur die"; `slice-263` §3 (Plan-Tabelle)
- **befund:** Der Fix legt vier Aussagen an, die nur bis zum Abschluss von `slice-263` stimmen. Die
  Plan-Tabelle von `slice-263` nennt Modul, Spec, `.d-check.closure.yml` und vier Altverstöße, aber
  keine dieser vier Stellen. Wenn `slice-263` schließt, ohne sie anzufassen, beschreiben vier Spiegel
  eine Grenze, die es nicht mehr gibt. Das ist dieselbe Klasse wie R1 F-7, nur in Gegenrichtung.
- **verifizierbar:** nein (Plan-Urteil).
- **klasse:** `zwischen-aussage-ohne-ruecknahme-adresse`

### R2-F-5 — LOW: `SPEC-093` beschreibt die Prüfung außerhalb von Strings weiter, als der Extraktor sie fährt

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §5 Regel 13 (Grenze gegen den Gegenstand prüfen)
- **pfad:** `spec/spezifikation.md` · `SPEC-093` „ein Zeichen außerhalb eines Strings, das kein
  JSON-Wert ist"; `tools/harness/extract-command.awk` · `if (c !~ /^[ \t\r\neEtrufalsn0-9.+-]$/) exit 3`
- **befund:** Geprüft wird jedes Zeichen einzeln gegen das Alphabet der Literale und Zahlen, nicht
  das Token als Wert. Die Eingaben `"a":fuse` und `"a":1e+-.` sind keine JSON-Werte und passieren mit
  Exit 0 (gemessen). Die Festlegung liest sich so, als würden sie geblockt. Bei maschinell erzeugter
  Eingabe ist der Fall kaum erreichbar, deshalb LOW.
- **verifizierbar:** ja — `printf '{"tool_input":{"command":"ls"},"a":fuse}' | bash .claude/hooks/pretooluse-command-guard.sh; echo $?` → 0.
- **klasse:** `grenze-gegen-beschreibung-statt-gegenstand-geprueft`

### R2-F-6 — INFO: Unter `make -i gates MAKEFLAGS=` entsteht der Nachweis trotz rotem Glied

- **kategorie:** INFO
- **quelle:** `MR-076` (`Grenze:`-Zeile); `SPEC-094`
- **pfad:** `Makefile` · `@case "$(firstword -$(MAKEFLAGS))" in`
- **befund:** Die Erkennung liest die Variable `MAKEFLAGS`. Eine Zuweisung auf der Kommandozeile
  überschreibt diese Variable, während der interne Ignore-Zustand von `make` bleibt (gemessen auf 4.3
  und 4.4.1: Nachweis, `green`, Exit 0). Das geht nur absichtlich und fällt unter „keine Sperre gegen
  Absicht" aus `MR-076`. `SPEC-094` sagt „unter `make -i` … kein Nachweis" aber ohne Einschränkung.
- **verifizierbar:** ja — Modell-Makefile aus den Proben.
- **klasse:** `erkennung-liest-ueberschreibbare-variable`

### R2-F-7 — INFO: Neben `[gates] … green` gibt es weitere Grün-Zeilen ohne Nachweis

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** `Makefile` · `@echo "[ci] gates + image-test green"`
- **befund:** Unter `make -i ci` steht nach `gates: kein Nachweis …` die Zeile
  `[ci] gates + image-test green`; unter `make -i fullbuild` gilt dasselbe für `[fullbuild] green`.
  Unter `-n` und `--trace` erscheint der `[gates]`-Text als Echo des Rezepts. Den Nachweis berührt
  keine dieser Zeilen, und der Stop-Hook liest nur den Nachweis. Diese Ausgabe gab es schon vor dem
  Fix.
- **verifizierbar:** ja — Modell-Makefile mit `ci: gates`.
- **klasse:** `gruen-zeile-ohne-urteil`

### R2-F-8 — INFO: Der dritte Treffer von `guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen` hat im Plan noch keine Folge

- **kategorie:** INFO
- **quelle:** Plan `slice-260` §8 · „ein dritter Treffer wäre eine Lücke"; Baseline `v6.17.0` ·
  `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register
- **pfad:** Plan `slice-260` §3 · „Pointer- und Abschnitts-Korrekturen F-7, F-9, F-10, F-11"
- **befund:** R1 F-7 war der dritte Treffer dieser Klasse. Die Plan-Änderung nach R1 behandelt ihn
  nur als Text-Korrektur. Den Ausgang (verkörpert, geplant oder gestrichen) verlangt erst die Closure
  §7. Notiert, damit der Lese-Schritt ihn nicht verliert.
- **verifizierbar:** nein (Closure-Urteil).
- **klasse:** `register-schwelle-erreicht-ohne-ausgang`

## Negativbefunde (geprüft, ohne Befund)

- **`gates`-Rezept gegen alle Flag-Formen:** Auf Make 4.3 und 4.4.1 gibt es keinen falsch negativen
  und keinen falsch positiven Treffer, gefahren mit 36 Formen × 2 Glied-Zuständen (Proben). Auch der
  Sub-Make (MAKELEVEL 1, `w` im ersten Wort), `-I dir`, `-C`, `-j N`, `--debug`, Variablen mit `i`
  und `MAKEFLAGS`/`GNUMAKEFLAGS` aus der Umgebung verhalten sich richtig. `a-check.mk` setzt
  `MAKEFLAGS` nicht, `.IGNORE` gibt es nicht. Die `-`-Präfixe der Rezepte stehen nur bei `clean`.
- **Rezept-Logik:** `record-gates.sh && echo …` in einer Zeile: Scheitert der Nachweis, endet `gates`
  rot und gibt kein `green` aus, weil der Exit-Status des `case` der des `&&` ist.
- **Konsumenten der Ausgabe/Form:** `gate-consistency` liest die Target-Zeile, und die ist
  unverändert. `ci.yml` und `release.yml` rufen `make ci` ohne Flags. `selbstpruefung.sh` wertet
  den Exit von `make gates` aus, nicht dessen Text. `.claude/commands/{plan,close}-welle.md`
  („endet mit `record-gates`") und `harness/README.md` („als letzter Schritt") stimmen weiter.
  `.claude/settings.local.json` (`make gates 2>&1 | tail -1`) sieht im Erfolgsfall weiter die
  `[gates]`-Zeile als letzte.
- **Stop-Hook-Fix:** Neun Fehlerformen gefahren (Proben), jede mit gültigem JSON und Exit 0. Das
  `if !` steht im Hook, das `set -e` von `working-tree-hash.sh` wirkt im eigenen Prozess, also
  greift der Exit ≠ 0 des Hash-Skripts. Ausnahme: der Schwester-Zweig in R2-F-2.
- **F-2-Fix:** `grep -E … > /dev/null` liest die Liste ganz, nachgefahren mit 10/10. Die Fehlerform
  „Quelle endet mit Exit ≠ 0" fällt unter `pipefail` weiter auf „kein Übergang". In der CI steht
  davor der `history-range-guard` aus `make trace-check`/`make adr-check` (unter `set -e`), lokal ist
  das ein `git diff --cached`-Fehler. Erreichbar ist das kaum, und es war vor dem Fix genauso.
  Dateinamen mit Nicht-ASCII-Zeichen, die `core.quotePath` in Anführungszeichen setzen würde, gibt
  es unter `docs/plan/planning/` nicht.
- **`SPEC-093` Durchlass-Satz:** Alle genannten Klassen wurden gegen den Wächter gefahren und gehen
  durch. Das `awk`-Beispiel stimmt unter dem einschränkenden Nebensatz.
- **`SPEC-096` und `blackbox-probe.md`:** „leer heißt diese vier" und „besteht nur aus Leerraum"
  stimmen mit `PROBE_FORMS="${PROBE_FORMS:-…}"` und der Leerraum-Prüfung des Skripts überein.
- **`MR-076`-Änderung:** Der Eintrag entsteht in diesem Slice. Die Ergänzung zu `make -i` steht vor
  der Closure im Ursprungs-Slice, und §Guard-Härtung (neuer MR statt Edit) betrifft akzeptierte
  Einträge aus früheren Vorgängen. Die Pflichtfelder sind erhalten, die `Grenze:`-Zeile trägt den
  Handaufruf.
- **Kommentare (§3.7), auch die nur umbrochenen Zeilen:** Neu und geändert sind der
  `.NOTPARALLEL`-Kommentar (Zusage), der `gates`-Kommentar (Zusage, Kopplung „in einer Zeile, denn
  …", Rang-Zeiger `SPEC-094, MR-076`), der Stop-Hook-Kommentar (Zusage mit Grenze der
  `set -e`-Semantik), der `pre-commit`-Kommentar (Rang-Zeiger `SPEC-095`, Kopplung, Abgrenzung
  gegen `-q`, Grenze) und der `ci.yml`-Kommentar (Kopplung). Keine Slice- oder Befund-Nummern, keine
  Review-Historie. Die bestehenden Slice-Nummern in nicht geänderten Kommentarzeilen (`slice-011`,
  `slice-053`) sind nicht Gegenstand.
- **Referenz-Richtung (§3.4):** Die neuen Texte in `SPEC-093` bis `SPEC-096` und in der §8-Zeile
  nennen keine ADR, keinen Slice, keine Welle und keinen Commit-Hash.
- **Plan-Änderung nach R1 (`01d7c500`):** Sie steht vor dem Code-Commit und nennt die drei
  zusätzlichen Dateien und ihre Gründe. Die Abgrenzung (Behebung von F-1 in `slice-263`) hat eine
  Adresse, die annimmt: `slice-263` §1 nennt das Ziel. Zu den Zwischen-Aussagen siehe R2-F-4.
- **Botschaft `596630db`, Proben-Aussagen:** F-2 (0/10 → 10/10), F-3 (die Formen am Modell, der
  echte Lauf mit Exit 0 ohne Nachweis) und F-5 (unlesbare Datei, unlesbarer Nachweis, passender
  Nachweis) habe ich nachgefahren, die Ergebnisse stimmen. Überdehnt ist nur „Die Aussagen sagen das
  jetzt" (R2-F-1). „make gates gruen" ist Sache des Verifiers.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 2 | R2-F-1, R2-F-2 |
| LOW | 3 | R2-F-3, R2-F-4, R2-F-5 |
| INFO | 3 | R2-F-6, R2-F-7, R2-F-8 |

**Wiederkehrende Klasse:** `grenzen-liste-wird-als-vollstaendig-gelesen` (R2-F-1) tritt in der
vierten Review-Runde in Folge auf (`slice-259` R1/R2, `slice-260` R1/R2). R2-F-2 hat dieselbe Form
wie R1 F-5 eine Verzweigung weiter: Gehärtet wurde der gemeldete Zweig, nicht die Fehlerpolitik des
ganzen Hooks. Steering-Loop-Signal: Die Arbeitsanweisung „Vertrag umdrehen" (Anker 18) wurde am
geänderten Absatz angewendet, aber nicht am Vertrag und an der Grenzen-Liste derselben Datei.

## Verdikt

**Nicht blockierend für HIGH, vor Merge zu klären: zwei MEDIUM.** Der Stilles-Grün-Pfad aus R1 F-1
ist an `slice-263` adressiert, und die Hälfte der Aussagen ist ehrlich. Vertrag und Grenzen-Liste von
`verify-closure-notes.md` und der Vertrag von `hooks.md` versprechen aber weiter die Prüfung, die es
nicht gibt (R2-F-1). R2-F-2 ist der fail-open-Zweig neben dem gerade gehärteten. Er braucht eine
Entscheidung: nachziehen oder in `SPEC-094` benennen. Die neue `make -i`-Abwehr hält auf beiden
gefahrenen Make-Versionen gegen alle geprüften Flag-Formen. Neue Fehler führt sie nicht ein, nur die
absichtliche Umgehung aus R2-F-6.
