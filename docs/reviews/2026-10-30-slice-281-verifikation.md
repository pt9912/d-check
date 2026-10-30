# Verifikation slice-281 — `--print-mk` `doc-manual`

**Rolle:** Verifier (Modul 8/11) · **Gegenstand:** Plan und DoD aus
`docs/plan/planning/in-progress/slice-281-print-mk-doc-manual.md`, Diff `4c955848..HEAD` ·
**Bezug:** [`DC-FA-CLI-010`](../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben),
[`DC-FA-CLI-013`](../../spec/lastenheft.md#dc-fa-cli-013--handbuch-und-spezifikation-aus-dem-werkzeug-lesen),
[`DC-FA-CLI-010.a`](../../spec/spezifikation.md#dc-fa-cli-010a--makefile-fragment) · Datum 2026-10-30.
Der Review-Report r1 wurde nicht übernommen; alle Belege sind selbst gefahren.

## DoD-Punkte

| DoD | Urteil | Beleg |
|---|---|---|
| Fragment | bestätigt | `docker run --rm d-check:latest --print-mk`: `grep -c '^\.PHONY'` = **14**; `grep -c -- '--disable'` = **6** (sechs von vierzehn, `doc-manual` keine); Recipe-Zeile: `@docker run --rm --network none -v "$(CURDIR):/repo:ro" $(DCHECK_REF) --manual "$(MANUAL)"`, `##`-Annotation vorhanden, kein `--enable`/`--disable`. |
| Lauf | bestätigt | Wegwerf-Fragment im Scratchpad, aus dem Repo-Verzeichnis: `make -s -f <frag> DCHECK_IMAGE=d-check:latest doc-manual MANUAL=planning` → Exit 0, 565 Zeilen stdout (Abschnitte aus `docs/user/benutzerhandbuch.md` und `spec/spezifikation.md`), stderr leer. Ohne `MANUAL` → `d-check: error: --manual braucht einen Begriff`, make-Exit 2. |
| Anforderung und Spiegel | bestätigt | Lastenheft 0.109.0 (Zeile 4236 Historie, Zeile 520 Aufzählung `doc-manual`), Spezifikation geändert (Diffstat 21 Zeilen). Siehe Spiegel-Befund unten (keiner). |
| Test | bestätigt | Mutation: Anführungszeichen um `$(MANUAL)` im Template entfernt, `make test` → `--- FAIL: TestCLI281_PrintMK_ManualTarget`, `cli_acceptance_test.go:3101: doc-manual-Recipe ohne "--manual \"$(MANUAL)\""` mit dem mutierten Recipe — Rot aus dem behaupteten Grund. Datei per `git checkout` zurückgesetzt, `git status --short` leer. Der Test prüft zudem Echo-Unterdrückung, `--network none`, Mount, Abwesenheit von `--enable`/`--disable` und 14 `.PHONY`. Hinweis: der Rot-Beleg ohne das Target insgesamt wurde nicht gefahren; die Mutation deckt die Recipe-Form, die erste Assertion (`\ndoc-manual: ## `) ist offensichtlich target-abhängig. |
| `make gates` | bestätigt | Schlusszeile: `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` (Exit 0). |

## Spiegel (4)

`grep -rnE 'dreizehn|Dreizehn|thirteen' spec internal docs/user`: verbleibende Treffer sind
(a) „dreizehn Quell-Tools/Alt-Tool" (`DC-QA-04`, nicht betroffen), (b) `finding.go:47`
„dreizehn Grund-Codes" (anderer Gegenstand), (c) Historie-Zeilen (Lastenheft 0.80.0,
Spezifikation 2026-08-31, Handbuch §11 — Chronik, bewusst), (d) `docs/user/benutzerhandbuch.md:883`
und `docs/user/operations.md:41` — **bewusste Release-Prep** laut Plan-Änderung
([`AGENTS.md`](../../AGENTS.md) §5 Regel 17). Kein doc-manual-relevanter Treffer außerhalb
dieser Ausnahme. Kopfkommentar `print_mk.go` steht auf „vierzehn" mit `doc-manual` in der Liste.

## Abgrenzung §1 (5)

Diffstat: `print_mk.go`, `cli_acceptance_test.go`, `abdeckung-tests.md`, Lastenheft,
Spezifikation, Slice-Plan, Review-Report. Keine Änderung an der `--manual`-Implementierung,
keine Vorbelegung von `MANUAL` (Recipe liest `$(MANUAL)` ohne `?=`), kein `$(error …)`-Guard,
kein ADR, Handbuch/`operations.md` unberührt. Eingehalten.

## Findings

- HIGH / MEDIUM: keine.
- LOW: Das Risiko „Begriff mit Leerzeichen" (§6) wurde in diesem Lauf nicht gemessen
  (`MANUAL="Schritt 1"`); der Plan nennt den Lauf als Bedingung für den Ausgang. Beim
  Planner zu klären, nicht Verifier-Befund gegen den Code.
- INFO: Zwischen Feature-Stand und Release-Prep nennen Handbuch §4.16 und `operations.md`
  noch dreizehn Targets — im Plan so entschieden und dort benannt.

## Offene DoD-Punkte (Planner, nicht bewertet)

Review durchgeführt/Report, Beobachtungs-Register, Closure-Notiz mit Lerneintrag, Risiko-Ausgänge
aus §6, drei Paarungen.

**Gesamt:** alle fünf geprüften DoD-Punkte bestätigt; keine Abweichung von Plan oder Abgrenzung.
