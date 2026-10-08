# Verifikation 2 — slice-258: Black-Box-Vorher/Nachher-Probe als make-Target (DoD)

- **Rolle:** Verifier (Modul 11), Frage „Bauen wir es richtig?". Abschließende Verifikation in frischem Kontext gegen §2 DoD, §1 Ziel und Abgrenzung und §3 Plan des Slice-Plans samt aller Plan-Änderungen (nach R1, Verifikation, R2, R3, R4). Geprüft wurde außerdem die Sensor-Datei gegen den Code ([`AGENTS.md`](../../AGENTS.md) §5 Regel 13). Keine ADR berührt. Die R4-Befunde wurden auf sachliche Schließung geprüft, Review-Entscheidungen waren nicht der Maßstab.
- **Gegenstand:** `231ce771..d70b0f8b`, 17 Commits. Endstand: `tools/blackbox-probe.sh`, `tools/blackbox-probe/` (README, vier Fixtures), `harness/sensors/blackbox-probe.md`, Target `blackbox-probe`, Index-Zeile in `harness/README.md`, `scan.ignore` in `.d-check.yml`, drittes Ventil in `harness/sensors/doc-check.md`.
- **Sensor-Evidence:** Alles selbst gemessen. Gefahren wurden `make gates` und `make blackbox-probe REF=HEAD`, dazu Probe-Läufe gegen `REF=v0.83.0` und `REF=v0.40.0`. Für den Bruch-Test wurde eine Meldung im Code geändert, danach zurückgesetzt und das Image mit `make build` neu gebaut. Acht `docker`-Shims im `PATH` (Scratchpad) und ein Klon unter `umask 077` (auch mit gelöschter getrackter Datei bzw. verschobenen Fixtures) kamen hinzu, außerdem `make trace-check`/`make adr-check RANGE=231ce771..HEAD`. Kein Push. Docker-Client/Server `29.8.2`.
- **Modell-ID:** claude-opus-5-5 (Claude Agent SDK)
- **Datum:** 2026-10-08

---

## DoD-Prüfung je Punkt (§2)

### 1. Skript und `make blackbox-probe REF=<ref>`: Vorher-Image, Läufe, getrennter Vergleich, Exit 0/1/2, fail-closed bei leerer Fixture- oder Formenmenge: **ERFÜLLT** (mit V2-1 als LOW)

- `make blackbox-probe REF=HEAD` auf `d70b0f8b` mit sauberem Baum: 20× `gleich` (`ids`/`links`/`targets` exit 1, `sauber`/`repo` exit 0). Schlusszeile `byte-identisch über 20 Vergleiche (Vorher d70b0f8b, stdout/stderr/Exit getrennt, Kanarienlauf vorher und nachher).`, rc 0, 47 s.
- Älterer REF:
  - `REF=v0.83.0`: 20× `gleich`, rc 0.
  - `REF=v0.40.0 PROBE_FORMS=-`: beide Kanarienläufe grün, dann `4 von 5 Vergleichen weichen ab`, rc 1. Die Abweichungen sind die fehlende Meldungs-Spalte der alten Ausgabe sowie `repo/-: stderr exit(2→0)`: Das alte Binary versteht die heutige `.d-check.yml` nicht, wie Grenze 3 sagt.
- Getrennte Ströme: Im Bruch-Test (Punkt 2) und im Shim „Zusammenfassungs-Zeile nur im Nachher" (unten) nennt die Zeile genau den betroffenen Strom.
- Exit-2-Pfade, selbst gefahren:

| Pfad | Ergebnis |
|---|---|
| Shim: leerer Mount für **alle** Läufe | rc 2. `Kanarienlauf vorher auf d-check:probe-vorher (sauber): Exit 2 statt 0 — d-check: error: Scan-Wurzel ist leer` |
| Shim: Mount ohne Markdown | rc 2. `Kanarienlauf vorher … (links): Exit 0 statt 1 — d-check: 0 Datei(en) geprüft` |
| Shim: toter Socket (`DOCKER_HOST=unix:///nonexistent.sock`, echte Client-Meldung) | rc 2. `… (sauber): Exit 1 statt 0 — failed to connect to the docker API …`. Das stille Grün aus V-1 ist damit geschlossen. |
| Shim: `docker run` → Exit 125 | rc 2. `Exit 125, kein Lauf des Werkzeugs — docker: Error response from daemon: simulated` |
| Shim: echter Lauf, Exit danach auf 137 gesetzt | rc 2. `Exit 137, kein Lauf des Werkzeugs` |
| `PROBE_FORMS=" "` direkt am Skript | rc 2. `PROBE_FORMS ist leer — nichts zu vergleichen ist kein gleicher Stand` |
| Fixture-Verzeichnis leer (Klon) | rc 2. `keine Fixtures unter tools/blackbox-probe/fixtures/` |

- `make … PROBE_FORMS=" "` fährt die Default-Formen (20 Vergleiche), weil make den Leerraum abstreift. Das ist „leer heißt Default", wie die Sensor-Datei sagt. Die Menge ist nie leer, also gibt es kein stilles Grün. `PROBE_FORMS=-` erreicht das Skript über `make` (5 Vergleiche).

### 2. Grundmenge an Fixtures plus Repo; Bruch-Test: geänderte Meldung wird Abweichung, unveränderter Stand byte-identisch: **ERFÜLLT**

- Es gibt vier Fixtures (`ids`, `links`, `sauber`, `targets`) und das Repo als Kopie seiner nicht ignorierten Dateien. Die Kopie und der Arbeitsbaum ergeben im selben Image dasselbe (`1000 Datei(en) geprüft, 0 Befund(e)`, Exit 0), die Kopie verliert also nichts, was der Default-Lauf liest.
- **Bewusstes Brechen:** In `internal/hexagon/core/rules/links.go:56` wurde `"Linkziel existiert nicht"` zu `"Linkziel fehlt (Mutation)"` geändert, dann lief `make blackbox-probe REF=HEAD`. Ergebnis: `4 von 20 Vergleichen weichen ab`. Betroffen sind genau `links/-`, `--json`, `--yaml` und `--doctor`, jeweils mit `stdout` und dem passenden Diff (`message:`-Feld bzw. `Hinweis:`-Zeile). Die übrigen 16 blieben `gleich`. Das Skript endete mit rc 1, make mit rc 2 (`Fehler 1`), wie die Sensor-Datei sagt. Die Probe wurde also aus dem richtigen Grund rot, an der richtigen Stelle.
- Danach `git checkout` und `make build`. `d-check:latest` hat wieder die ID `a2e38087`, identisch mit dem Build aus `make gates`. Der Folgelauf (`PROBE_FORMS=-`) war byte-identisch, rc 0.
- Rauschen (§6): Zwei getrennt gebaute Images desselben Stands waren in sechs Läufen byte-identisch, darunter zweimal die vollen 20 Vergleiche.

### 3. Index-Zeile (`kein Gate`), Sensor-Datei (Vertrag, Grenzen), `make gates` grün: **ERFÜLLT**

- `make gates` selbst gefahren, rc 0: `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`. Coverage 94,70 % ≥ 93 %, doc-check `994`/`1000 Datei(en) geprüft, 0 Befund(e)`.
- Die Index-Zeile in `harness/README.md` liegt in der Werkzeug-Tabelle und trägt `kein Gate`. `blackbox-probe` steht in keinem der Targets `gates`, `ci`, `fullbuild` und in keinem Workflow, die Abgrenzung aus §1 ist also eingehalten.
- `make trace-check` und `make adr-check` über `231ce771..HEAD`: 17 Commits, 0 Befunde.
- Den Abgleich der Sensor-Datei mit dem Code zeigt der Abschnitt unten.

### 4. Review und Verifikation: **ERFÜLLT mit diesem Report**

R1 bis R4 und die erste Verifikation liegen unter `docs/reviews/`. Dieser Report ist die abschließende Verifikation.

### 5. Closure-Notiz, Register, Risiko-Ausgänge, Paarungen: **OFFEN (erwartet)**

§7 ist leer, die Ausgänge in §6 stehen auf *(offen)*. Das gehört zur Closure und blockiert diese Verifikation nicht. Messwerte für die Ausgänge: Laufzeit 47 s (alle Formen, HEAD) bis etwa 80 s (älterer REF, Wert aus R4). Kein Rauschen über zwei Images, was [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) stützt.

---

## Sensor-Datei gegen Code (Regel 13)

| Aussage | Stand | Beleg |
|---|---|---|
| Vertrag: Vorher aus `git archive`, Dockerfile und Defaults des Stands | stimmt | `docker build … --build-arg VERSION=0.0.0-dev --target runtime`; die ARG-Defaults des Dockerfile (`GO_VERSION=1.27.1`, `GOLANGCI_LINT_VERSION=v2.14.0`) stimmen an HEAD mit dem Makefile überein |
| Formen und `PROBE_FORMS` übersteuert, „leer heißt Default" | stimmt | über `make` gemessen (5 bzw. 20 Vergleiche) |
| Nebenwirkung: `:latest` wird mit `0.0.0-dev` neu gebaut, auch bei `VERSION` auf der Kommandozeile | stimmt | `make -n blackbox-probe VERSION=9.9.9` zeigt nur `VERSION=0.0.0-dev` |
| Grenze 1 (nur Fixtures und Formen) | stimmt | — |
| Grenze 2 (falsche Erwartung) | stimmt | Urteil, nicht prüfbar |
| Grenze 3 (zwei Binaries, heutige Eingaben) | stimmt | `REF=v0.40.0`: `repo/-` exit(2→0) wegen `.d-check.yml` |
| Grenze 4: eine Mount-Quelle, `umask 077` ändert nichts | stimmt | Klon mit `drwx------`: `repo/-` und `repo/--json` jetzt `gleich (exit 0)` statt vorher `exit 2` |
| Grenze 4: leerer Mount → 2, ohne Markdown → `links` 0, Daemon-Ausfall → `sauber` 1 | stimmt | drei Shims, je rc 2 im Kanarienlauf |
| Grenze 4: Abbruch bei jedem Exit außer 0/1/2 (125, 137) | stimmt | zwei Shims |
| Grenze 4: Exit-Änderung an `sauber`/`links` lässt den Kanarienlauf scheitern | stimmt | Shim (`sauber` nur im Nachher Exit 1): rc 2, kein Vergleich. Siehe V2-2 |
| Grenze 4 Restgrenze (zeitlicher Ausfall zwischen den Kanarienläufen) | stimmt, ehrlich benannt | — |
| Grenze 5 / Bindung (kein Gate) | stimmt | — |
| Exit-Tabelle 0/1 | stimmt | — |
| Exit-Tabelle 2: Aufzählung der Ursachen | **unvollständig** | Es fehlt „Kopie des Arbeitsbaums gescheitert", der mit R4 neu eingeführte Pfad (dazu „git archive gescheitert"). Siehe V2-1 |
| `make` normalisiert auf 2 | stimmt | Bruch-Test: make rc 2, Skript rc 1 |

## Schließung R4

| R4 | Stand | Beleg |
|---|---|---|
| R4-1 (Repo als zweite, ungeprüfte Mount-Quelle) | **sachlich geschlossen** | `$PWD` wird nicht mehr gemountet. Das Repo ist eine Kopie unter `$WORK/fx` mit derselben `chmod`-Behandlung wie die Fixtures, die der Kanarienlauf belegt. Im `umask 077`-Klon gibt der Repo-Fall jetzt Exit 0 statt Exit 2. Rest siehe V2-3 (INFO). |
| R4-2 (Kanarienlauf am Wortlaut der Zusammenfassungs-Zeile) | **sachlich geschlossen** für die Ausgabe, **benannt** für den Exit | Shim mit `Datei(en)` → `Dateien` nur im Nachher: rc 1, `5 von 5` `ABWEICHUNG … stderr` mit Diff, kein Kanarien-Abbruch. Ändert ein Vorgang den Exit von `sauber`, scheitert die Probe weiterhin mit rc 2. Grenze 4 nennt das, die Meldung diagnostiziert es aber falsch, siehe V2-2. |

## Befunde

### V2-1 — LOW — Gelöschte, noch getrackte Datei bricht die Probe ab; der Pfad fehlt in der Exit-Tabelle

- **pfad:** `tools/blackbox-probe.sh` · `git ls-files -z -co --exclude-standard | tar --null -T - -cf - | …` · `fail "Kopie des Arbeitsbaums gescheitert"`; `harness/sensors/blackbox-probe.md` Exit-Tabelle Zeile 2.
- **befund:** `-c` listet auch getrackte Dateien, die im Arbeitsbaum gelöscht und nicht mit `git rm` aus dem Index genommen sind. `tar` kann sie nicht lesen, und mit `pipefail` scheitert die Probe. Gemessen im Klon: `rm CHANGELOG.md`, dann `GESCHEITERT — Kopie des Arbeitsbaums gescheitert` (`tar: CHANGELOG.md: Funktion stat fehlgeschlagen`), rc 2. Derselbe Abbruch kommt, wenn ein Fixture verschoben, aber nicht aus dem Index genommen wird. Vor R4 (Mount von `$PWD`) fehlte die Datei dann einfach. Der Fehler ist laut, fail-closed und kein stilles Grün. Er trifft aber einen gewöhnlichen Zwischenstand während eines Slice, der eine Datei entfernt. Die Exit-Tabelle nennt diesen Abbruchgrund nicht, ebenso wenig `git archive … gescheitert`.
- **verifizierbar:** ja, `rm <getrackte Datei>; REF=HEAD PROBE_FORMS=- bash tools/blackbox-probe.sh` gibt rc 2.
- **Vorschlag:** gelöschte Dateien aus der Liste nehmen (etwa `git ls-files -co --exclude-standard` ohne die Einträge aus `git ls-files -d`) und die Exit-Tabelle um beide Pfade ergänzen.

### V2-2 — LOW — Ein Exit-Wechsel an `sauber`/`links` wird als Umgebungsfehler gemeldet, obwohl das Skript ihn unterscheiden könnte

- **pfad:** `tools/blackbox-probe.sh` · `kanarie()`, Meldung „die Umgebung traegt keinen Vergleich (Daemon, Mount oder gesehener Inhalt)".
- **befund:** Grenze 4 benennt den Fall ehrlich. Die Meldung zur Laufzeit schreibt die Ursache aber der Umgebung zu. Im Shim (nur das Nachher-Image endet auf `sauber` mit Exit 1) stand der Vorher-Kanarienlauf auf `d-check:probe-vorher` in derselben Schleife bereits grün, der Abbruch kam erst auf `d-check:latest`. Besteht das eine Image den Kanarienlauf und scheitert das andere, liegt eher eine Verhaltensänderung vor als ein Umgebungsausfall. Die Meldung könnte das sagen, statt den Leser zum Daemon zu schicken. Das ist dieselbe Klasse wie R3-2/R4-2 (`echter-lauf-als-ausfall-abgewiesen`), jetzt auf den dokumentierten Exit-Fall eingeengt.
- **verifizierbar:** ja, mit einem Shim, der für `…/fx/sauber:/repo:ro d-check:latest` Exit 1 liefert: rc 2 mit der Umgebungs-Meldung.

### V2-3 — INFO — Der Kanarienlauf belegt den Mechanismus der Repo-Kopie, nicht ihren Inhalt

- **befund:** Ein Shim, der **nur** den Mount `…/fx/repo` leert, ergibt weiterhin `gleich repo/… (exit 2)` und rc 0. Eine reale Ursache dafür fand sich nicht mehr. Ein Kopierfehler scheitert über `pipefail` (V2-1). Rechte und Mount-Pfad teilt die Repo-Kopie mit den Kanarien-Fixtures. Ein Absturz nach Umfang (137) bricht ab. Die Formulierung im Code und im Plan, „Mount-Quelle …, die der Kanarienlauf mitbelegt", ist für den Mechanismus richtig. Den Inhalt der Repo-Kopie prüft sie nicht. Kein Handlungsbedarf, nur benannt.

## Negativbefunde

- **Abgrenzung (§1):** eingehalten. Die Probe ist kein Gate, steht nicht in `AGENTS.md` §6, und es gibt keine weiteren Fixtures über die Grundmenge hinaus.
- **Reihenfolge Plan vor Code:** Jede Plan-Änderung liegt vor ihrem Fix-Commit, nach R4 `bfd07aa6` vor `d70b0f8b`.
- **Docker/make-only (§3.1):** Das Skript nutzt nur `git`, `tar`, `docker`, `cmp`, `diff`, `sed`, `head`, `tr` und `chmod`, keinen Host-Interpreter.
- **Kommentare (§3.7):** Die Blöcke im Skript und im Makefile tragen Zusage, Abgrenzung, Grenze und Kopplung. Befund-Nummern und Review-Historie kommen nicht vor.
- **Aufräumen:** `trap` entfernt `$WORK` bei jedem `fail`. `make clean` nimmt `:probe-vorher` mit.

## Verdikt

DoD-Punkte 1 bis 3 sind mit eigenen Belegen erfüllt. Der Bruch-Test wurde aus dem richtigen Grund rot, und der unveränderte Stand war byte-identisch. R4-1 und R4-2 sind sachlich geschlossen. Für R4-2 bleibt der dokumentierte Exit-Fall, dessen Meldung die Ursache falsch zuordnet (V2-2). Offen sind zwei LOW (V2-1, V2-2) und eine INFO (V2-3). Keiner der Befunde ist ein stilles Grün, und keiner blockiert die Closure. V2-1 ist eine kleine Korrektur samt Ergänzung der Exit-Tabelle und gehört vor die Closure oder als benanntes Risiko in §6. DoD-Punkt 5 (Closure) steht aus.
