# Verifikation — slice-275: Push-Range der CI bei Force-Push

- **Rolle:** Verifier (Modul 11) — DoD, Plan und ADR gegen den Stand
- **Gegenstand:** `docs/plan/planning/done/wellenlos/slice-275-ci-push-range-bei-force-push.md` §1, §2, §3 (samt Plan-Änderung), §4, §6
- **Range:** `fbbbecae..HEAD` (Spitze `3329d6aa`)
- **Eingang:** [ADR-0109](../plan/adr/0109-ci-push-range-faellt-bei-force-push-zurueck.md) (Proposed); Review-Report nur als Kontext
- **Datum:** 2026-10-20

## Messungen (Kommando und Ergebnis)

| # | Kommando | Ergebnis |
|---|---|---|
| M1 | `bash tools/harness/ci-range-probe.sh` | Exit 0, sechs Proben ok: pull-request, neuer-branch, normaler-push, force-push-feature, force-push-default, force-push-ohne-origin |
| M2 | `make ci-range-probe` | Exit 0, dieselben sechs ok |
| M3 | Bewusstes Brechen an Kopien im Scratchpad (Repo unberührt), siehe unten | je eine Probe rot, aus dem richtigen Grund |
| M4 | CI-Schritt lokal nachgebaut (`set -euo pipefail; RANGE="$(bash tools/harness/ci-range.sh)"`) in Wegwerf-Repo mit tatsächlich nicht mehr vorhandenem `before`-Commit (reflog expire, gc; `git cat-file -e` scheitert) | Feature-Branch: Exit 0, `RANGE=origin/main..<head>`, 1 Commit; Default-Branch: Exit 2, Schritt bricht ab, `REACHED` nicht erreicht |
| M5 | `bash tools/harness/history-range-guard.sh origin/main..HEAD` im selben Repo | `aufgeloest, 1 Commit(s) — OK`, Exit 0 |
| M6 | Rebase des Feature-Branchs auf einen weitergewanderten `main`, `git rev-list --count origin/main..HEAD` | 1 (nur der Branch-eigene Commit; Default-Commits fallen heraus) |
| M7 | `make gate-consistency`, `make gates` | beide grün; `[gates] baseline-verify + … + planning-check green`, doc-check `1172 Datei(en) geprüft, 0 Befund(e)`, semgrep 0 findings |

### Bewusstes Brechen (M3)

| Probe | Zurückgenommen | Rot durch |
|---|---|---|
| a | Rückfall (Ausgabe wieder `$PUSH_BEFORE..$HEAD_SHA`) | `FAIL force-push-feature -- … Ausgabe '1111…..<head>' (erwartet 'origin/main..<head>')`; alle anderen ok |
| b | Default-Branch-Ausnahme (`if false`) | `FAIL force-push-default -- Exit 0 (erwartet 2), Ausgabe 'origin/main..<head>'` — genau das stille Grün über leerer Range |
| c | origin-Wächter (`if false`) | `FAIL force-push-ohne-origin -- Exit 0 (erwartet 2), Ausgabe 'origin/main..<head>'`; alle anderen ok |

Jede Probe fällt einzeln und nur an der benannten Stelle; Arbeitsbaum danach unverändert.

## DoD-Punkte

- **Range auflösen — erfüllt.** `tools/harness/ci-range.sh` führt nur ein nicht als Commit auflösbares `before` auf `origin/<default>..HEAD` zurück; `ci.yml` ruft es (`RANGE="$(bash tools/harness/ci-range.sh)"` unter `set -euo pipefail`; Exit 2 bricht ab, M4). PR, neuer Branch, normaler Push behalten ihre Range.
- **Beleg — erfüllt.** Die vier geforderten Fälle plus Pull-Request und fehlendes origin sind in der Probe; der laute Fall endet mit Exit 2 (M1–M3).
- **Entscheidung — Inhalt erfüllt, Statuswechsel offen.** ADR-0109 deckt Entscheidung 1–3, Fitness, Konsequenzen und Grenze mit dem Code. Kopf und Index stehen noch auf `Proposed` — offen, vorgesehen für die Closure (keine Verletzung).
- **`make gates` grün — erfüllt** (M7).
- Review-Report liegt vor. Offen und keine Verletzung: Closure-Notiz §7, Risiko-Ausgänge §6, Register, Paarungen.

## Plan-vs-Code (Abgrenzung §1)

Dateien im Diff: `ci.yml`, `Makefile`, `harness/README.md`, `ci-range.sh`, `ci-range-probe.sh`, ADR + Index, Slice-Datei, Review-Report. Alle durch §3 bzw. die Plan-Änderung vor dem Code gedeckt (Makefile-Target und `REF_NAME` sind dort genannt). `history-range-guard.sh` im Diff unverändert, kein PR-#9-Merge, keine Ruleset-Änderung; in `ci.yml` nur der Range-Schritt (inline-Auflösung ersetzt, `REF_NAME` ergänzt; die folgenden Befehle unverändert). Abgrenzung eingehalten.

## ADR-0109 gegen Code

- Entscheidung 1–3 stimmen mit dem Skript überein. Die Aussage zur Range-Größe nach Rebase stimmt: sie ist `origin/<default>..HEAD`, Commits, die der Default-Branch schon trägt, fallen heraus (M6).
- Die Fitness-Funktion nennt zwei gebrochene Stellen; der origin-Wächter (c) ist ebenfalls probegedeckt und rot — die ADR ist damit nicht falsch, nur knapper als der Beleg.

## Befunde

- **HIGH / MEDIUM:** keine.
- **LOW:** ADR-0109 Fitness-Funktion spricht von „fünf Fällen und einer Rückfall-Lücke" und zählt dann sechs Proben; Formulierung ist missverständlich (`docs/plan/adr/0109-ci-push-range-faellt-bei-force-push-zurueck.md`). Als Proposed noch änderbar.
- **INFO:** Probe c (origin-Wächter) fehlt in der Fitness-Aussage „beides am Skript gebrochen und gesehen" (dort Rückfall und Default-Ausnahme); gebrochen und gesehen hier zusätzlich.
- **INFO:** `trace-check`/`adr-check` mit der Fallback-Range selbst nicht gefahren (Docker-Image-Build ausgeklammert); der Guard (`history-range-guard`) lief. Die Grenze „Auflösbarkeit, nicht Eignung" und „Workflow-Aufruf wird von keiner Probe gelesen" steht in ADR und Skript.

## Offen (keine Verletzung)

ADR Accepted (Kopf + Index), Closure-Notiz §7, Risiko-Ausgänge §6, Register, drei Paarungen.

## Urteil

DoD-Liefer-Punkte erfüllt, Plan und ADR konsistent mit dem Code; keine DoD-Verletzung.
