# Review R2 — slice-258: Black-Box-Vorher/Nachher-Probe als make-Target

- **Review-Art:** Code. Runde 2, nur über die Einarbeitungen seit R1. Geprüft
  gegen den Slice-Plan (`slice-258`, §1 und beide Plan-Änderungen), die R1-Befunde
  F-1 bis F-6, die Verifier-Befunde V-1 bis V-3 und die Hard Rules
  [`AGENTS.md`](../../AGENTS.md) §3.1, §3.7 sowie §5 Regel 13. Die DoD-Abhakung
  prüft dieses Review nicht.
- **Gegenstand:** `slice-258` · Range `87b2ada9..f24fd1ef`, eingeschränkt auf
  `tools/blackbox-probe.sh`, `Makefile`, `harness/sensors/blackbox-probe.md`,
  `harness/sensors/doc-check.md` und den Slice-Plan. Die beiden Fix-Commits sind
  `3c8180e0` (R1) und `f24fd1ef` (Verifikation). Mitgelesen, nicht im Diff:
  `tools/image-publish.sh` (Gegenprobe), `Makefile` `build`/`image-test`, die
  CLI-Hilfe des Images.
- **Skill:** `reviewer.md` @ 1.18.0.
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-08
- **Eingangs-Kontext:** R1-Report und Verifikations-Report zu `slice-258`, beide
  vom 2026-10-08. Vorherige Findings am selben Werkzeug: R1 F-1 und V-1, beide
  zur Klasse *stilles Grün bei gleichem Scheitern*.
- **Proben** (Docker 29.8.2; Shims im Scratchpad, kein Code geändert):
  - `make blackbox-probe REF=HEAD`: 20× `gleich`, `byte-identisch über 20
    Vergleiche`, rc 0, 49 s.
  - **Shim „toter Daemon nur für `run`"** (`DOCKER_HOST=unix:///nonexistent.sock exec docker "$@"`):
    rc 2, `Exit 1 ohne Befund-Ausgabe — kein Lauf des Werkzeugs: failed to connect
    to the docker API …`. Der Fall aus V-1 endet damit laut.
  - **Shim „Mount erreicht den Container nicht"**: Er ersetzt die Quelle jedes
    `-v …:/repo:ro` durch ein leeres Verzeichnis. Damit ist simuliert, dass ein
    Daemon den Client-Pfad nicht sieht, etwa bei Remote-`DOCKER_HOST` oder DinD.
    Ergebnis: 20× `gleich … (exit 2)`, `byte-identisch über 20 Vergleiche`,
    **rc 0**. Das Werkzeug meldet auf beiden Seiten `d-check: error: Scan-Wurzel
    ist leer: /repo …` mit Exit 2.
  - **Nicht lesbare Quelle** (Fixture-Kopie mit `chmod -R go-rwx`) direkt gegen
    das Image: rc 2, `d-check: error: lstat /repo/docs: permission denied`.
  - Exit, stdout-Größe und `^d-check:` auf stderr wurden je Fixture über elf Formen
    gemessen: `-`, `--json`, `--yaml`, `--doctor`, `--repair`, `--repair-broad`,
    `--trace`, `--require-complete`, `--bogus`, `--range=x..y` und `--print-mk`.
    `--repair` und `--repair-broad` enden auf `links` und `targets` mit **Exit 1 und
    leerem stdout** (`d-check: 0 Reparatur-Hunk(s), 0 review-pflichtig`). Alle
    anderen Exit-1-Fälle tragen stdout, alle Exit-2-Fälle eine `d-check:`-Zeile.
  - `REF=HEAD PROBE_FORMS=--repair bash tools/blackbox-probe.sh`: rc 2,
    `… links (--repair): Exit 1 ohne Befund-Ausgabe — kein Lauf des Werkzeugs:
    d-check: 0 Reparatur-Hunk(s) …`.
  - `make -n blackbox-probe REF=HEAD VERSION=9.9.9`: Der rekursive Build läuft mit
    `VERSION=0.0.0-dev -t d-check:latest`. Auch eine Version auf der Kommandozeile
    wird überstimmt.

## Findings

### R2-1 — MEDIUM — Gleich gescheiterte Läufe des Werkzeugs (Exit 2) zählen weiter als „byte-identisch"

- **quelle:** Skill Prüffragen 1 und 19; `harness/sensors/blackbox-probe.md`
  §Ausgabe und Ausgänge (0 = „byte-identisch über alle Vergleiche")
- **pfad:** `tools/blackbox-probe.sh` · „`2) grep -q '^d-check:' "$4.err" || why="Exit 2 ohne Meldung des Werkzeugs" ;;`"
- **befund:** Das Lebenszeichen-Kriterium nimmt jeden Exit 2 mit `d-check:`-Zeile
  als gültigen Lauf. Exit 2 ist aber der Code, mit dem das Werkzeug meldet, dass es
  seine Eingabe nicht prüfen konnte. Erreicht der Mount den Container nicht, etwa
  bei Remote-Daemon oder DinD, enden alle Läufe beidseitig gleich mit „Scan-Wurzel
  ist leer". Gemessen: 20× `gleich (exit 2)`, „byte-identisch über 20
  Vergleiche", rc 0, ohne einen einzigen Vergleich auf echten Eingaben. Dasselbe
  trifft in einem Checkout mit `umask 077` die vier `repo/*`-Vergleiche, denn
  `repo=$PWD` wird anders als die Fixtures nicht lesbar gemacht und endet mit
  `permission denied`, Exit 2. Die Ausgabe ist dann `rc 0` mit 16 echten und vier
  hohlen Vergleichen. Es ist dieselbe Klasse wie R1 F-1, nur verschoben vom
  Docker-Exit auf den Werkzeug-Exit. MEDIUM statt HIGH, weil das Target
  deklariert kein Gate ist.
- **verifizierbar:** ja. Shim `docker`, der für `run` die Quelle von `…:/repo:ro`
  durch ein leeres, lesbares Verzeichnis ersetzt; danach `REF=HEAD bash tools/blackbox-probe.sh`.
- **klasse:** stilles-gruen-bei-gleichem-scheitern

### R2-2 — MEDIUM — Grenze 4 nennt die größte verbleibende Lücke nicht

- **quelle:** Skill Prüffrage 18; [`AGENTS.md`](../../AGENTS.md) §5 Regel 13
- **pfad:** `harness/sensors/blackbox-probe.md` · „Ein Docker-Ausfall, der mit Exit 0 endete, bliebe unerkannt"
- **befund:** Grenze 4 führt als einzige Restlücke einen Docker-Ausfall mit
  Exit 0 an, und den gibt es in keiner gemessenen Form. Die Lücke, die tatsächlich
  offen ist, fehlt: Ein Lauf, in dem das Werkzeug auf beiden Seiten gleich
  scheitert (Exit 2, R2-1), gilt als „gleich". Die Exit-Tabelle sagt weiter
  „0 = byte-identisch über alle Vergleiche". Wer die Datei liest, nimmt rc 0
  deshalb als gelaufenen Vergleich. Die Grenze ist am Text des Kriteriums geprüft
  worden, nicht an den Exit-Pfaden des Werkzeugs.
- **verifizierbar:** ja, wie R2-1.
- **klasse:** grenzen-liste-ohne-groesste-luecke

### R2-3 — MEDIUM — „Exit 1 trägt in jeder Form einen Befund auf stdout" stimmt für `--repair` nicht: falsch rot mit falscher Diagnose

- **quelle:** Skill Prüffrage 19; [`AGENTS.md`](../../AGENTS.md) §5 Regel 13;
  `harness/sensors/blackbox-probe.md` („`PROBE_FORMS` übersteuert")
- **pfad:** `tools/blackbox-probe.sh` · „`# Exit 1 traegt in jeder Form einen Befund auf stdout`"
  und „`1) [ -s "$4.out" ] || why="Exit 1 ohne Befund-Ausgabe" ;;`"
- **befund:** `--repair` und `--repair-broad` enden auf `links` und `targets` mit
  Exit 1 und leerem stdout, weil es keinen eindeutigen Hunk gibt und der Befund
  bestehen bleibt. Mit `PROBE_FORMS=--repair` bricht die Probe deshalb mit rc 2
  ab und meldet „kein Lauf des Werkzeugs". Die angehängte stderr-Zeile
  (`d-check: 0 Reparatur-Hunk(s)`) zeigt aber, dass das Werkzeug lief. Ein Slice,
  der die Reparatur-Ausgabe ändert, kann seine Zusage „ohne Schalter unverändert"
  mit der dokumentierten Übersteuerung nicht prüfen. Der Kommentar und Grenze 4
  stellen eine Allaussage über „jede Form" auf, gemessen ist sie nur über die vier
  Default-Formen. Das ist laut, also kein stilles Grün. MEDIUM, weil das Kriterium
  ein neuer Leseweg ist, der nur gegen den Fall gefahren wurde, den er beheben
  soll (Prüffrage 19, außerhalb des Gate-Pfads).
- **verifizierbar:** ja. `REF=HEAD PROBE_FORMS=--repair bash tools/blackbox-probe.sh` liefert rc 2.
- **klasse:** allaussage-ueber-gemessene-teilmenge

### R2-4 — INFO — Der Nachher-Build überstimmt auch ein `VERSION` auf der Kommandozeile

- **quelle:** Slice-Plan `slice-258` (Plan-Änderung nach der Verifikation, V-2);
  `harness/sensors/blackbox-probe.md` §Nebenwirkung
- **pfad:** `Makefile` · „`@$(MAKE) --no-print-directory build VERSION=0.0.0-dev`"
- **befund:** Die Nebenwirkung ist jetzt beschrieben, und die beschriebene Folge
  ist laut. `image-test` baut `:latest` über seine Prerequisite `build` neu.
  `image-publish` prüft per `image-verify-published.sh` gegen `TESTED_AMD64` und
  setzt bei Abweichung keinen Tag. Kein stiller Pfad. Die Sensor-Datei spricht von
  einem „zuvor mit Release-Version gebauten `:latest`". Dass auch
  `make blackbox-probe VERSION=x` die Version auf `0.0.0-dev` setzt (gemessen mit
  `make -n`), erschließt sich daraus, steht aber nicht ausdrücklich da.
- **verifizierbar:** ja, `make -n blackbox-probe REF=HEAD VERSION=9.9.9`.
- **klasse:** —

## Schließung R1 / Verifikation

- **F-1:** geschlossen für Docker-Fehler: 125 und ein toter Daemon (rc 1, leeres
  stdout) enden mit rc 2, gemessen am echten Client. **Offen in neuer Form** über
  Exit 2 (R2-1).
- **F-2:** teilweise geschlossen. Grenze 4 steht da, nennt aber die verbliebene
  Lücke nicht (R2-2).
- **F-3:** geschlossen. Punkt 4 von `harness/sensors/doc-check.md` nennt jetzt
  drei Ventile samt `scan.ignore`, mit Baseline, Cache und Fixtures.
- **F-4:** geschlossen. Plan §1 lautet jetzt „Build-Defaults dieses Stands", und
  `VERSION` ist auf beiden Seiten fest gesetzt.
- **F-5:** geschlossen über Grenze 3. Sie stimmt mit `fixtures+=("repo=$PWD")`
  überein.
- **F-6:** geschlossen. `clean` nennt `:probe-vorher`, die Mitnahme von `:arm64`
  ist harmlos, weil `2>/dev/null || true` dem Bestand entspricht.
- **V-1:** geschlossen. Das Kriterium hängt nicht mehr am Meldungstext, und der
  echte Client mit totem Socket endet mit rc 2.
- **V-2:** geschlossen über die Dokumentation, mit Rest R2-4 (INFO).
- **V-3:** geschlossen: „leer heißt Default" steht in der Sensor-Datei.

## Negativbefunde

- **Falsch grün über einen Docker-Fehler:** keiner gefunden. Der Docker-Client
  schreibt seine Fehler auf stderr, nie mit `d-check:`-Präfix. Der Start-Exit ist
  125/126/127, der Daemon-Ausfall endet mit 1 bei leerem stdout. Keine dieser
  Formen erfüllt das Kriterium. Gemessen: toter Socket, und die Exit-Codes
  stammen aus R1/V-1. Der falsch grüne Pfad geht über das Werkzeug selbst (R2-1).
- **Falsch rot in den Default-Formen:** keiner. In den Formen `-`, `--json`,
  `--yaml` und `--doctor` trägt jeder Exit 1 stdout, auf allen vier Fixtures
  gemessen. `--trace` und `--print-mk` enden mit 0. `--bogus` und
  `--require-complete` ohne `--trace` enden mit 2 und `d-check: error:`.
- **Fehlermeldung des Abbruchs:** `$(tail -n 1 "$4.err")` steht innerhalb der
  Meldung von `fail`, die auf stderr geht. Ein leeres stderr ergibt eine leere
  Nachricht, aber keinen Fehlpfad. Geprüft, ohne Befund.
- **Rekursiver Build:** `IMAGE` und `REF` reisen korrekt. `IMAGE` kommt über
  `MAKEFLAGS` in den Sub-Make, `REF` wird nur an das Skript gereicht. Der
  Sub-Make-Fehler bricht das Recipe ab, bevor das Skript läuft. Geprüft, ohne
  Befund.
- **Kommentar-Klassen §3.7:** Der neue Kommentar in `lauf()` trägt eine Zusage
  („nur ein Lauf des Werkzeugs zählt"), eine Abgrenzung („nicht am Wortlaut einer
  Docker-Meldung", begründet am Gegenstand: der Wortlaut wechselt mit der
  Docker-Version) und eine Grenze. Es gibt keine Befund-Marker, keine
  Slice-Nummer und keine Review-Historie. Die verworfene Mustervariante aus
  `3c8180e0` wird nicht erzählt. Ebenfalls ohne Befund: der Makefile-Kommentar zur
  VERSION-Kopplung. Die inhaltliche Unrichtigkeit „in jeder Form" ist R2-3, kein
  Klassen-Befund.
- **§3.1:** neu sind nur `grep`, `tail` und `case`, kein Host-Interpreter.
  Geprüft, ohne Befund.
- **Plan-Abgrenzung:** kein Gate-Bindepunkt, keine Workflow-Pflicht. Die
  Plan-Änderungen liegen vor ihren Fix-Commits (`47f3e282` vor `3c8180e0`,
  `af691492` vor `f24fd1ef`). Geprüft, ohne Befund.
- **`harness/sensors/doc-check.md`:** Die Zahlen 45/248 sind unverändert, und das
  dritte Ventil ist sachlich richtig beschrieben. Die Liste in `scan.ignore`
  stimmt mit der `.d-check.yml` überein. Geprüft, ohne Befund.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 3 | 0 | 1 |

Wiederkehrende Klassen: `stilles-gruen-bei-gleichem-scheitern` (R1 F-1 → V-1 →
R2-1; dieselbe Lücke ist zum dritten Mal in einer neuen Form aufgetreten, über
drei Prüfläufe desselben Vorgangs) und `grenzen-liste-ohne-groesste-luecke`
(R1 F-2, F-3 → R2-2). Das ist ein Steering-Loop-Signal: Jeder Fix schloss die
gemessene Fehlerform, und die nächste Form trat auf, weil das Kriterium gegen
den Fall gebaut war, nicht gegen die Exit-Pfade des Werkzeugs (Prüffrage 19).
Im Register zählt das als **ein** Vorgang (`slice-258`).

## Verdikt

**Nicht freigegeben.** R2-1 und R2-2 blockieren: Die Probe meldet bei einem auf
beiden Seiten gleich scheiternden Werkzeug weiterhin rc 0 „byte-identisch", und
die Sensor-Datei nennt das nicht. R2-3 blockiert ebenfalls (MEDIUM). Der Pfad ist
laut, aber eine dokumentierte Eingabe (`PROBE_FORMS`) führt zu einem falschen
Ausgang mit falscher Diagnose, und Kommentar und Grenze behaupten mehr, als
gemessen ist. R2-4 blockiert nicht.
