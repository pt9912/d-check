# Verifikation slice-258 — Black-Box-Vorher/Nachher-Probe als make-Target (DoD)

- **Rolle:** Verifier (Modul 11), Frage „Bauen wir es richtig?". Geprüft gegen §2 DoD, §1 Ziel und Abgrenzung und §3 Plan des Slice-Plans (samt beider vermerkter Plan-Änderungen). Keine ADR berührt (Plan: „Harness-Werkzeug, keine Spec-Aussage"). Review-Entscheidungen waren nicht der Maßstab; die R1-Befunde wurden auf sachliche Schließung geprüft.
- **Gegenstand:** `231ce771..3c8180e0`, fünf Commits: `af3fc934` (Plan-Änderung Fixtures aus dem Scan), `87b2ada9` (Skript, Fixtures, Target, Sensor-Datei, Index-Zeile, `scan.ignore`), `ea212876` (R1-Report), `47f3e282` (Plan-Änderung nach R1), `3c8180e0` (R1-Einarbeitung F-1 bis F-6).
- **Sensor-Evidence:** Alles selbst gemessen, nichts aus dem Implementer-Bericht übernommen: `make blackbox-probe REF=HEAD`, drei `docker`-Shims im `PATH` (Scratchpad), fünf Eingabe-Fehlerpfade, eine leere Fixture-Menge in einem Wegwerf-Worktree, ein Bruch-Test mit geänderter Meldung im Code (zurückgesetzt, neu gebaut), `VERSION=9.9.9 … PROBE_FORMS=--print-mk`, `make gates`, `make trace-check`/`make adr-check RANGE=231ce771..HEAD`, `make review-coverage`. Kein Push. Docker-Client/Server `29.8.2`.
- **Modell-ID:** claude-opus-5-5 (Claude Agent SDK)
- **Datum:** 2026-10-08

---

## DoD-Prüfung je Punkt (§2)

### 1. `tools/blackbox-probe.sh` und `make blackbox-probe REF=<ref>`: Vorher-Image, Läufe, getrennter Vergleich; Exit 0/1/2; fail-closed bei leerer Fixture- oder Formenmenge: **TEILWEISE ERFÜLLT** (V-1)

- `make blackbox-probe REF=HEAD` auf HEAD `3c8180e0`, sauberer Baum: 20 Zeilen `gleich` (5 Fixtures × 4 Formen; `ids`/`links`/`targets` exit 1, `sauber`/`repo` exit 0), Schlusszeile `blackbox-probe: byte-identisch über 20 Vergleiche (Vorher 3c8180e0, stdout/stderr/Exit getrennt).`, make-rc 0, ~51 s.
- Getrennter Vergleich: im Bruch-Test unten meldet die Zeile genau `stdout` als abweichenden Strom, nicht `stdout stderr`. stdout und stderr werden in getrennte Dateien geschrieben, der Exit in eine dritte (Code gelesen, Ausgabe bestätigt).
- Exit 1 bei Abweichung: Skript-rc 1, Schlusszeile `2 von 10 Vergleichen weichen ab`. Über `make` ist der rc 2 (`make: *** … Fehler 1`), wie die Sensor-Datei sagt.
- Exit 2 bei gescheitertem Lauf. Selbst gefahren:

| Pfad | Lauf | Ergebnis |
|---|---|---|
| REF fehlt | `bash tools/blackbox-probe.sh` / `make blackbox-probe` | rc 2 / make-rc 2; `GESCHEITERT — REF fehlt` |
| REF unbekannt | `REF=gibtsnicht` | rc 2; `ist kein Commit dieses Repos` |
| REF syntaktisch gültig, kein Objekt | `REF=000…0` (40 Stellen) | rc 2; dieselbe Meldung |
| Formen nur Leerraum | `REF=HEAD PROBE_FORMS='   '` | rc 2; `PROBE_FORMS ist leer — nichts zu vergleichen ist kein gleicher Stand` |
| Fixtures leer | Worktree mit leerem `tools/blackbox-probe/fixtures/` | rc 2; `keine Fixtures unter tools/blackbox-probe/fixtures/` |
| Shim: `docker run` → `docker: Error response from daemon: …`, exit 125 | `REF=HEAD` | rc 2; `Exit 125 — kein Lauf des Werkzeugs: docker: Error response from daemon: simulierter Mount-Fehler.` |
| Shim: `docker run` → `Cannot connect to the Docker daemon at unix:///var/run/docker.sock. …`, exit 1 | `REF=HEAD` | rc 2; `Docker-Fehler — Cannot connect to the Docker daemon …` |
| **Shim: `docker run` gegen echten toten Daemon** (`DOCKER_HOST=unix:///nonexistent.sock exec docker "$@"`) | `REF=HEAD` | **rc 0**, 20× `gleich … (exit 1)`, `byte-identisch über 20 Vergleiche` (V-1) |

Die beiden Fälle aus dem Auftrag (125 und „Cannot connect", exit 1) enden mit Exit 2. Der dritte Shim lässt den installierten Docker-Client seine echte Meldung schreiben. Er meldet einen nicht erreichbaren Daemon mit `failed to connect to the docker API at unix:///nonexistent.sock; check if the path is correct and if the daemon is running: …` und rc 1 (direkt gemessen). Diese Form trifft das Muster `^docker: |Cannot connect to the Docker daemon|Error response from daemon` nicht. Das stille Grün aus R1 F-1 besteht damit auf diesem Host in der Form fort, die er tatsächlich erzeugt.

- `PROBE_FORMS=''` (gesetzt, aber leer) und `make … PROBE_FORMS='  '` fallen auf die Default-Formen zurück (20 Vergleiche, rc 0). Make entfernt führenden Leerraum aus dem Wert, und `${PROBE_FORMS:-…}` greift. Das ist kein stilles Grün, weil die Menge nicht leer ist. Der fail-closed-Zweig ist über `make` aber nicht erreichbar (V-3).

### 2. Grundmenge an Fixtures plus das Repo; Bruch-Test: geänderte Meldung wird als Abweichung gemeldet, der unveränderte Stand als byte-identisch: **ERFÜLLT**

- Vier Fixtures (`ids`, `links`, `sauber`, `targets`) plus `repo`. Die Exit-Codes zeigen, dass jedes Fixture auslöst, was es soll; R1 hat die Befunde einzeln gemessen.
- **Bewusstes Brechen.** `internal/hexagon/core/rules/links.go:56` von `"Linkziel existiert nicht"` auf `"Linkziel fehlt (Bruch-Test)"` geändert, dann `make blackbox-probe REF=HEAD PROBE_FORMS='- --json'` gefahren. Ergebnis: `ABWEICHUNG links/-: stdout` und `ABWEICHUNG links/--json: stdout` mit dem richtigen Diff-Auszug (`-…Linkziel existiert nicht` / `+…Linkziel fehlt (Bruch-Test)` bzw. das `"message"`-Feld). Die übrigen 8 Vergleiche bleiben `gleich`, Schlusszeile `2 von 10 Vergleichen weichen ab`, Skript-rc 1. Rot aus dem richtigen Grund, an der richtigen Stelle. Nebenbei belegt: `PROBE_FORMS` erreicht das Skript über `make`.
- Zurückgesetzt (`git checkout`), `make build VERSION=0.0.0-dev`, erneut gefahren: `byte-identisch über 10 Vergleiche`, rc 0. Der Arbeitsbaum ist sauber.
- Rauschen (§6): zwei getrennt gebaute Images desselben Stands sind über 20 Vergleiche byte-identisch, viermal gemessen (`REF=HEAD` vollständig, zweimal mit Teilmengen, einmal `--print-mk`).

### 3. `harness/README.md` (Werkzeug-Zeile, `kein Gate`), `harness/sensors/blackbox-probe.md` (Vertrag, Grenzen); `make gates` grün: **ERFÜLLT**, Grenze 4 mit V-1

- Index-Zeile in der Werkzeug-Tabelle, Bindung `kein Gate · meldet, urteilt nicht über den Repo-Zustand`. `gate-consistency` ist in `gates` grün.
- Sensor-Datei: Vertrag, fünf Grenzen, Exit-Tabelle, Bindung. Gegen den Gegenstand geprüft: Formen, Fixture-Menge, Exit-2-Gründe und der Hinweis „make normalisiert auf 2" stimmen. Grenze 3 (zwei Binaries über heutige Eingaben) stimmt mit `fixtures+=("repo=$PWD")`. Grenze 4 nennt die erkannten Meldungsformen richtig und benennt „eine andere Form des Daemon-Ausfalls mit Exit 1" als Lücke. Sie liest sich aber als Randfall, während es auf dem installierten Docker der Regelfall ist (V-1).
- `make gates` selbst gefahren, HEAD `3c8180e0`, sauberer Baum: Exit 0, `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `d-check: 990 Datei(en) geprüft, 0 Befund(e)`; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; `fetch-baseline-cache: verify ok (54 Dateien, vollständig)`.

### 4. Review durchgeführt, Report unter `docs/reviews/`; Verifikation: **ERFÜLLT** mit diesem Bericht (siehe V-2)

R1-Report `2026-10-08-slice-258-blackbox-vorher-nachher-probe-r1.md` (MEDIUM 2, LOW 1, INFO 3; „nicht freigegeben, F-1 und F-2 blockieren"); `make review-coverage` → `0 Befund(e)`. Einarbeitung in `3c8180e0`, Stelle für Stelle geprüft:

- **F-1 (MEDIUM), teilweise geschlossen.** Exit außerhalb 0/1/2 bricht mit rc 2 ab (125-Shim). Die Meldungen `docker: …`, `Error response from daemon` und `Cannot connect to the Docker daemon` brechen ab („Cannot connect"-Shim). Ein nicht erreichbarer Daemon in der Meldungsform von Docker 29.8.2 bleibt still grün (V-1). Die Commit-Botschaft sagt „ein Daemon-Ausfall mit Exit 1 [endet] in 'GESCHEITERT'". Das ist mit einem Shim gemessen, der die alte Meldungsform schreibt, nicht am Gegenstand (`AGENTS.md` §5 Regel 13/15).
- **F-2 (MEDIUM), geschlossen bis auf V-1.** Grenze 4 steht in der Sensor-Datei, und die Exit-Tabelle nennt „ein Lauf ohne Exit des Werkzeugs oder mit Docker-Fehler".
- **F-3 (LOW), geschlossen.** `harness/sensors/doc-check.md` Punkt 4 heißt jetzt „drei Ventile" und nennt `scan.ignore` mit Baseline, Cache und den Fixtures.
- **F-4 (INFO), geschlossen.** Plan §1 heißt jetzt „Dockerfile und Build-Defaults dieses Stands". Das Target baut Nachher über `$(MAKE) build VERSION=0.0.0-dev`. Gemessen: `VERSION=9.9.9 make blackbox-probe REF=HEAD PROBE_FORMS=--print-mk` → Build-Zeile mit `--build-arg VERSION=0.0.0-dev`, 5× `gleich`, rc 0.
- **F-5 (INFO), geschlossen.** Grenze 3 der Sensor-Datei.
- **F-6 (INFO), geschlossen.** `clean` nennt `$(IMAGE):probe-vorher`; als Text geprüft, `make clean` nicht ausgeführt.

**Fix-Commit in der Größe seiner Befunde?** Weitgehend ja: +8 Zeilen Skript (F-1), Doku (F-2, F-3, F-5), Makefile (F-4, F-6), 4 Dateien, +32/−6. Die Plan-Änderung `47f3e282` liegt vor dem Fix-Commit. Zwei Stellen gehen über die Befunde hinaus, beide klein und in der Botschaft deklariert:
- `clean` räumt zusätzlich `$(IMAGE):arm64` ab. R1 F-6 nennt das als „dieselbe Lage", ohne es zu fordern.
- Der F-4-Fix ersetzt die Prerequisite `build` durch einen rekursiven `make build VERSION=0.0.0-dev`. Damit überschreibt `make blackbox-probe` `$(IMAGE):latest` immer mit einem `0.0.0-dev`-Image, auch wenn der Aufrufer `VERSION` exportiert hat (V-2).

### 5. Closure-Notiz, Register, Risiko-Ausgänge, Paarungen: **OFFEN, wie erwartet**, kein Befund

§6 steht auf `(offen)`, §7 trägt `—`; `make verify-closure-notes` greift erst in `done/`. Belege für die Closure: **Laufzeit** ~50 s für 20 Vergleiche bei gecachtem Vorher-Build. **Rauschen** ist über zwei Images desselben Stands nicht aufgetreten (siehe DoD 2).

---

## Abgrenzung (§1): gehalten

- Kein Gate: `blackbox-probe` steht nicht in `gates`, `ci`, `fullbuild` oder einem Workflow (Makefile gelesen; Index-Zeile in der Werkzeug-Tabelle).
- Keine Workflow-Pflicht: `AGENTS.md` und die Slice-Vorlage sind im Range unberührt.
- Jede Datei im Diff (18) steht in §3 oder einer der beiden Plan-Änderungen (`scan.ignore`, `tools/blackbox-probe/README.md`, `harness/sensors/doc-check.md`, R1-Report, Plan).
- Rückführung aus §4 nicht eingetreten: das Vorher-Image baut aus `git archive HEAD` ohne eigene Versions-Logik.
- `make trace-check RANGE=231ce771..HEAD` → `5 Commit(s) — OK`, `0 Befund(e)`; `make adr-check` (gleiche Range) → `0 Befund(e)`.

---

## Verifier-Befunde

### V-1: MEDIUM

- **Kategorie:** MEDIUM
- **Quelle:** DoD Punkt 1 („2 bei gescheitertem Lauf"); R1 F-1; `AGENTS.md` §5 Regel 13 und 15
- **Pfad:** `tools/blackbox-probe.sh` · „`grep -qE '^docker: |Cannot connect to the Docker daemon|Error response from daemon'`"
- **Befund:** Der installierte Docker-Client (29.8.2) meldet einen nicht erreichbaren Daemon mit `failed to connect to the docker API at …` und rc 1. Das Muster trifft diese Form nicht. Mit einem Shim, der `docker run` unverändert gegen einen toten Socket weiterreicht, endet die Probe mit `byte-identisch über 20 Vergleiche`, rc 0. Jeder Vergleich trägt dabei `exit 1`, auch `sauber` und `repo`, die sonst 0 liefern. Damit besteht das stille Grün, das R1 F-1 schließen sollte, in der realen Fehlerform dieses Hosts fort. Der Fix ist gegen einen Shim mit der älteren Meldungsform gemessen worden, nicht gegen den Gegenstand. Grenze 4 der Sensor-Datei benennt „eine andere Form" ehrlich, stellt sie aber als Randfall dar. MEDIUM, nicht HIGH: Das Target ist deklariert kein Gate, und ein Daemon, der *zwischen* `image inspect`/`build` und den Läufen ausfällt, ist selten.
- **Verifizierbar:** ja. Shim `docker`: `if [ "$1" = run ]; then DOCKER_HOST=unix:///nonexistent.sock exec <echtes docker> "$@"; fi; exec <echtes docker> "$@"`, dann `REF=HEAD bash tools/blackbox-probe.sh`.
- **Erwartete Aktion:** Das Muster um `failed to connect to the docker API` erweitern. Robuster wäre ein Kriterium, das nicht an Meldungstexten hängt: etwa `docker info` bzw. `docker version --format '{{.Server.Version}}'` nach den Läufen (fail bei rc ≠ 0), oder der Nachweis, dass jeder Lauf eine Ausgabe des Werkzeugs trägt (bei Exit 0/1 ist das die Schlusszeile `d-check: N Datei(en) geprüft` auf stderr; für `--json`/`--yaml` zu prüfen). Die Sensor-Datei (Grenze 4) danach am Gegenstand nachziehen. Die Messung gegen den echten Client gehört in die Closure.
- **Klasse:** stilles Grün bei gleichem Scheitern (verwandt mit `BEO-ALL/stilles-gruen-ueber-leerer-range`); Grenze gegen Beschreibung statt Gegenstand geprüft

### V-2: LOW

- **Kategorie:** LOW
- **Quelle:** Prüfauftrag „Fix-Commit in der Größe seiner Befunde"; Register [`BEO-ALL/fix-commit-ausserhalb-review-range`](../plan/planning/observations/BEO-ALL/fix-commit-ausserhalb-review-range/state.md)
- **Pfad:** `Makefile` · „`@$(MAKE) --no-print-directory build VERSION=0.0.0-dev`"
- **Befund:** Der F-4-Fix ändert die Mechanik des Targets: Statt der Prerequisite `build` (mit `VERSION` des Aufrufers) baut es jetzt rekursiv mit festem `VERSION=0.0.0-dev` nach `$(IMAGE):latest`. Das ist das Image, das `image-test` prüft und `image-publish` als `TESTED_AMD64` erwartet. Wer `make blackbox-probe` zwischen `make ci VERSION=x` und `make image-publish` fährt, ersetzt das geprüfte Image durch ein `0.0.0-dev`-Image. Der Publish-Pfad vergleicht Digests und fiele dann laut. Das ist Bedienfehler-Territorium und nicht still, wurde aber von keinem Review gesehen. Dazu kommt `clean` mit `:arm64` als kleine Mitnahme. Beides ist in der Botschaft deklariert, und der Fix bleibt klein (+32/−6). Es wäre das dritte Auftreten der Klasse „Fix trägt Mechanik außerhalb der Review-Range" (Register 2×: `slice-256`, `slice-257`).
- **Verifizierbar:** ja (Makefile lesen; `make blackbox-probe` gibt die Build-Zeile mit `VERSION=0.0.0-dev -t d-check:latest` aus).
- **Erwartete Aktion:** In der Closure begründen, warum dieser Bericht eine R2 ersetzt (Umfang: eine Recipe-Zeile), oder R2 über `ea212876..3c8180e0`. Optional: Nachher unter eigenem Tag bauen (`$(IMAGE):probe-nachher`) statt `:latest` zu überschreiben. Register-Beleg bei der Closure. Mit diesem Beleg erreicht die Klasse 3×, und der Lese-Schritt ist fällig.
- **Klasse:** Fix-Commit trägt neue Mechanik außerhalb der Review-Range

### V-3: INFO

- **Kategorie:** INFO
- **Quelle:** DoD Punkt 1 („fail-closed bei leerer … Formenmenge"); `AGENTS.md` §5 Regel 13
- **Pfad:** `tools/blackbox-probe.sh` · „`PROBE_FORMS="${PROBE_FORMS:-- --json --yaml --doctor}"`"
- **Befund:** Eine leere `PROBE_FORMS` bedeutet „Default", nicht „leer". Nur reiner Leerraum führt in den fail-closed-Zweig, und über `make` ist auch das nicht möglich: make entfernt führenden Leerraum, `make blackbox-probe REF=HEAD PROBE_FORMS='  '` fährt 20 Vergleiche, rc 0. Ein stilles Grün ohne Vergleich entsteht dadurch nicht. Die Sensor-Datei sagt aber nur „`PROBE_FORMS` übersteuert".
- **Verifizierbar:** ja (oben).
- **Erwartete Aktion:** optional einen Halbsatz in die Sensor-Datei: „leer = Default".
- **Klasse:** —

Negativ geprüft ohne Befund:
- Byte-Identität `REF=HEAD` über 20 Vergleiche, rc 0.
- Bruch-Test: geänderte Meldung → genau `links/-` und `links/--json` als `stdout`-Abweichung mit richtigem Diff, rc 1; zurückgesetzt → identisch.
- 125-Shim und „Cannot connect"-Shim → rc 2.
- REF fehlt, unbekannt, Null-SHA; Formen nur Leerraum; leere Fixture-Menge → jeweils rc 2, laut.
- F-4 mit exportiertem `VERSION=9.9.9` und `--print-mk` → identisch.
- F-2/F-3/F-5/F-6 als Text geschlossen.
- Abgrenzung (kein Gate, keine Pflicht).
- `make gates`, `trace-check`, `adr-check`, `review-coverage`.

---

## Verdict

- **DoD 1 teilweise erfüllt:** Die Mechanik und die Fehlerpfade aus dem Auftrag stimmen. Der reale Daemon-Ausfall des installierten Docker endet still grün (V-1, MEDIUM).
- **DoD 2 erfüllt**, mit Bruch-Beleg (rot aus dem richtigen Grund) und Rückkehr zu identisch.
- **DoD 3 erfüllt**: Index-Zeile, Sensor-Datei, `make gates` selbst grün. Grenze 4 ist nach V-1 nachzuziehen.
- **DoD 4 erfüllt** mit R1 und diesem Bericht. R1 F-2 bis F-6 sind sachlich geschlossen, F-1 nur teilweise (V-1).
- **DoD 5 erwartungsgemäß offen.**

Kein HIGH. V-1 vor der Closure beheben oder als Grenze am Gegenstand umformulieren und dabei die reale Meldungsform nennen. V-2 annehmen oder begründen (Register 3×). V-3 ohne Pflicht-Aktion.
