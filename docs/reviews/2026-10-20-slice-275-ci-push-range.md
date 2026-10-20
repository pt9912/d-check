# Review slice-275 — CI-Push-Range fällt bei nicht auflösbarer Basis zurück

- **Review-Art:** Code (gegen Plan, ADR, Hard Rules)
- **Gegenstand:** `slice-275`, Diff `fbbbecae..HEAD` (Plan-Änderung `f7f36a84`, Code `903ed405`)
- **Skill:** `reviewer.md` @ 1.20.0
- **Modell-ID:** claude-sonnet-5-5
- **Datum:** 2026-10-20
- **Eingangs-Kontext:** `slice-275` (§1 Abgrenzung, §3 Plan samt Plan-Änderung, §6 Risiken), `ADR-0109` (Proposed bis Closure, beabsichtigt), `ADR-0013`, `SPEC-102`, `AGENTS.md` §3.1 / §3.2 / §3.6 / §3.7 / §3.9, BEO `stilles-gruen-ueber-leerer-range`, BEO `pipe-swallows-gate-exit-code`.
- **Zusätzlich gefahren (Wegwerf, gelöscht, nicht im Diff):** `bash tools/harness/ci-range-probe.sh` (6 ok); `ci-range.sh` direkt in einem Wegwerf-Repo mit `set -euo pipefail` + `RANGE="$(bash …)"` (Exit 2 bricht ab: bestätigt) sowie mit nicht gesetztem/leerem `PUSH_BEFORE` und `PR_BASE`; drei Mutationen an einer Kopie des Skripts gegen die Probe (Rot-Beleg).

## Findings

Keine HIGH, keine MEDIUM, keine LOW.

### F-1 — INFO — Zusage „mehr Commits, nie weniger“ reicht weiter als die Fallunterscheidung

- **quelle:** Maintainability (kein Konventions-Anker, daher INFO); `AGENTS.md` §3.7 (Klasse Zusage)
- **pfad:** `tools/harness/ci-range.sh` · „mehr Commits, nie weniger“
- **befund:** `origin/<default>..HEAD` schließt Commits aus, die der Default-Branch schon trägt; nach einem Force-Push, der den Branch auf einen neueren Default-Stand rebased, waren diese gegenüber `before..HEAD` neu im Push und werden nun nicht mehr geprüft. Sie sind bereits Teil des Default-Branch, der Verlust ist praktisch klein, die Zusage aber wörtlich nicht allgemein wahr.
- **verifizierbar:** nein (Szenario-Lesung; ein Rebase-Probenfall würde es zeigen)
- **klasse:** zusage-ueber-fallmenge-hinaus

### F-2 — INFO — Fetch-Fehler wird verschluckt; der Wächter prüft Existenz, nicht Frische

- **quelle:** Prüffrage 19 (neuer Leseweg `2>/dev/null || true`); Fehlerform bewertet, kein Finding
- **pfad:** `tools/harness/ci-range.sh` · „fetch --no-tags origin“ in `fetch_default`
- **befund:** Schlägt der Fetch fehl und liegt ein älterer `origin/<default>` vor, läuft der Rückfall gegen den veralteten Stand. Das vergrößert die Range (mehr Commits, nicht weniger), ist also kein stilles Grün; fehlt die Ref ganz, fängt `rev-parse --verify` mit Exit 2. Die Probe fährt den Fetch-Weg nicht (die Ref liegt dort schon vom Push), die Fehlerform ist ungeprüft, aber ungefährlich.
- **verifizierbar:** nein
- **klasse:** fetch-fehler-still-verschluckt-folge-harmlos

### F-3 — INFO — Leere/nicht gesetzte Eingaben haben keine Probe; der Schutz liegt dahinter

- **quelle:** Prüffrage 1 / 19, BEO `stilles-gruen-ueber-leerer-range` bewertet
- **pfad:** `tools/harness/ci-range.sh` · „echo "$PR_BASE..$HEAD_SHA"“
- **befund:** Nicht gesetztes `PR_BASE`/`PUSH_BEFORE` bricht laut ab (Exit 1, `unbound variable`); ein leeres `PR_BASE` liefert `..<sha>` (leere Range, Exit 0 im Skript), leeres `PUSH_BEFORE` fällt im Zweig „nicht auflösbar“ zurück. Die leere Range fängt erst `history-range-guard` in `make trace-check`/`adr-check` (Exit 2) — nicht das Skript. Auf GitHub ist `base.sha` unter `pull_request` immer gesetzt; der Pfad ist vorbestehend und nicht verschlechtert.
- **verifizierbar:** ja — Probe mit `PR_BASE=""`.
- **klasse:** weiche-schutz-liegt-nachgelagert

## Beantwortung der Fragen des Auftrags

1. **Stilles Grün?** Nein. `RANGE="$(bash …)"` unter `set -e` bricht bei Exit 2 des Skripts den Schritt ab (gefahren: Zuweisung mit Befehlssubstitution gibt den Exit-Status weiter, `NOT-ABORTED` wurde nicht erreicht). Keine Pipe im Weg, `pipefail`-Falle entfällt. Eine leere Rückfall-Range (Branch ⊂ Default) fängt `history-range-guard` laut (Makefile: `trace-check`/`adr-check` rufen ihn zuerst). Hinweise gehen auf stderr, stdout trägt nur die Range.
2. **Fallunterscheidung:** `^0+$` für neuen Branch korrekt; `cat-file -e "…^{commit}"` löst Tag-Objekte auf Commits und weist Nicht-Commits ab (→ Rückfall bzw. Exit 2 auf dem Default-Branch); `pull_request` kehrt vor jeder `REF_NAME`-Nutzung zurück (Merge-Ref irrelevant); Default-Branch-Vergleich exakt, Branch-Namen mit `/` unproblematisch; Fetch-Fehler siehe F-2.
3. **Probe:** trägt die in der ADR genannten fünf Fälle plus Rückfall ohne origin (ADR-Zeile 74 deckt sich). Rot-Belege selbst gefahren an einer Kopie: (A) Default-Ausnahme entfernt → `force-push-default` FAIL; (B) `origin`-Wächter entfernt → `force-push-ohne-origin` FAIL; (C) Rückfall entfernt → `force-push-feature` und `force-push-ohne-origin` FAIL. Alle aus dem richtigen Grund (Erwartungswert Range/Exit verletzt).
4. **Hard Rules:** §3.1 bash+git, kein Host-Interpreter jenseits der POSIX-Klasse, keine Docker-Umgehung; §3.2 keine Suppression; §3.6 keine Schwellen-Senkung — die Range wird nie kleiner als ein Fehler, die Aufnahme als Ausnahme trägt ADR-0109; §3.7 Skript-Köpfe tragen Zusage/Grenze/Eingabe mit ADR-Verweis, keine Historie (Kommentar-Grenze der Probe „ob die Workflow-Datei das Skript ruft“ ist ehrlich benannt); §3.9 `uses:`-Pins in `ci.yml` unverändert (Diff trifft nur den `run`-Block und `env`). `harness/README.md`-Zeile ist als „kein Gate“ markiert, Target existiert im Makefile samt `.PHONY`.
5. **Abgrenzung:** `history-range-guard.sh` unverändert; Makefile und `REF_NAME` stehen in der Plan-Änderungsnotiz vor dem Code samt Spiegel-Grep. Keine Mitnahme außerhalb des Plans gefunden.

## Negativbefunde

- `.github/workflows/ci.yml` (Workflow-Weiche, Fehlerformen, Pins) — geprüft, ohne Befund.
- `tools/harness/ci-range.sh` (Fallunterscheidung, Exit-Codes, Kommentare) — geprüft, nur INFO oben.
- `tools/harness/ci-range-probe.sh` (Abdeckung, Rot-Beleg, Wegwerf-Repo ohne Netz) — geprüft, ohne Befund.
- `Makefile` / `harness/README.md` / ADR + Index (Target-Deklaration, Bindung „kein Gate“) — geprüft, ohne Befund.
- Hexagon-Import, Netzzugriff außerhalb `external`, Gate-Suppression — nicht berührt.

## Kategorie-Summary

HIGH 0 · MEDIUM 0 · LOW 0 · INFO 3

## Verdikt

Kein blockierendes Finding. Merge-fähig aus Review-Sicht; die drei INFO sind Kenntnisnahme, DoD-Konformität prüft der Verifier.
