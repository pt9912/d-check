# Verifikation slice-284 — Messung aus git (Modul `metrics`)

**Rolle:** Verifier (Modul 11) · **Datum:** 2026-10-30 · **Stand:** `18f5412d`
**Gegenstand:** DoD (§2) und Plan (§1 Abgrenzung, §3 Plan-Tabelle samt Plan-Änderung `8dcc6c17`) von
`docs/plan/planning/in-progress/slice-284-metrics-messung-aus-git.md`.
**Methode:** Sensoren selbst gefahren, Erwartung am Host mit `git` nachgebaut, Image `d-check:latest`
frisch gebaut (`make build`), Wegwerf-Repos unter dem Scratchpad (`.../scratchpad/verify284/`), nicht im Repo.

## Ergebnis je DoD-Punkt

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| Port und Adapter, Gleichheit mit `numstat` | **bestätigt** | s. (1) |
| Prüfung Messung ≤ Obergrenze | **bestätigt** | s. (2) |
| Config-Schlüssel `metrics.excluded` | **bestätigt mit Vorbehalt** | s. (7) — Unterstützung am Image: `excluded: [docs/vertraege/, src/move.txt]` wirkt (Messung 27 bzw. 47 ohne Ausschluss von `move.txt`); Fehlerfälle (Glob, `..`, `/`, Steuerzeichen) liegen im Test `configyaml/metrics_test.go` (`make test` grün), nicht selbst am Image gefahren |
| Fail-closed und Betrieb | **bestätigt** | s. (3) |
| `make gates` grün | **bestätigt** | s. (5) |
| Review durchgeführt, Report liegt vor | **bestätigt** | `docs/reviews/2026-10-30-slice-284-r1.md` (F-1…F-8, 4 MEDIUM); Einarbeitung im Fix-Commit gesehen (F-1 Binär/Submodul, F-2 Zusage, F-4 Vorfahr-Prüfung, F-6 `ValidLiteralPath`, F-8 Historie). Die Wirksamkeit der Einarbeitung ist in (1)–(3) mitgeprüft |
| Doku-Update (Handbuch/`operations.md` nicht) | **bestätigt** | Diff `8dcc6c17..HEAD` berührt weder Handbuch, `operations.md`, README noch CHANGELOG |
| Grund-Code `metrics-over-limit` (Spec §4, `AllReasons()`, Klartext) | **bestätigt** | s. (4) |
| Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | **offen** | §7 steht auf „Wird bei der Closure gefüllt", §6-Risiken tragen noch „wird bei der Closure gesetzt"; Slice liegt in `in-progress/`. Kein Mangel, sondern noch nicht fällig |
| ADR-0112 `Proposed` → `Accepted` (§5 Closure-Trigger) | **offen** | Bei der Closure; Folgepflicht der ADR (Gleichheit am Fixture) ist nach (1) erfüllt |

## (1) Port und Adapter — Gleichheit mit `git diff --numstat --no-renames`

- Fixture am Host nachgebaut (Inhalt von `linesFixture`, plus `src/move.txt` mit 40 Zeilen, in `HEAD` ein Block von
  10 Zeilen verschoben). Kommando aus dem Auftrag: `git diff --numstat --no-renames <base>..<head> -- ':(literal)src/' ':(literal)docs/' ':(exclude,literal)docs/vertraege/'`.
  Zeilen: `a.txt 2+1`, `added 4+0`, `bin.dat -`, `crlf 3+2`, `del 0+3`, `empty 1+0`, `new 5+0`, `old 0+4`, `ü ber 2+0`, leere
  Dateien und `mode.sh` 0. **Summe ohne `move.txt` = 27** (= Soll des Tests `TestLineChanges_GleichNumstat`),
  mit `move.txt` (10+10) = 47.
- Image gegen dasselbe Repo (`--enable metrics`, `:ro`, Block `path: src/` + `docs/`, `excluded` mit `docs/vertraege/` und `src/move.txt`):
  `Messung 27` (Befund bei limit 26, grün auf der Grenze bei limit 27). Ohne Ausschluss von `move.txt`:
  `Messung 47`, für `move.txt` allein `Messung 20`. **go-git == git** in allen diesen Formen.
- Eigene Formen, die der Test nicht hat — git (Standard-Algorithmus) gegen Image:
  Hälften getauscht (20+20 → 40 ✓), Block nach vorn (10+10 → 20 ✓), repetitive Datei (1+1 → 2 ✓), Hälften getauscht in 20 Zeilen
  (10+10 → 20 ✓), Block 21–40 hinter Block 41–60 (20+20 → 40 ✓) und sechs Dateien mit 80 zufälligen Zeilen aus 25 Werten
  (57/57/57/56/56/55 je Seite → 114/114/114/112/112/110, alle ✓). **Keine Abweichung gefunden.** Das belegt die
  Gleichheit für diese Formen, nicht für jede Eingabe; die Grenze (Halb-Treffer-Heuristik) benennt die Spec jetzt
  selbst (Lastenheft 0.110.2), und eine Abweichung ist damit nicht ausgeschlossen, nur hier nicht aufgetreten.
- Kurz-Hash als Basis (7 Zeichen) löst auf (Plan mit `${B:0:7}`), `nope/` erscheint als Hinweis „weder an der Basis noch an HEAD", Binärdatei als „1 Datei(en) nicht gezählt".

## (2) Prüfung gegen die Obergrenze

Am Image (Repo mit Plänen 1…6, `rounding: 1`, Formel 135/10000): Messung 27 / limit 27 grün (Grenze); limit 26 ⇒
`metrics-over-limit` mit „Messung 27 liegt ueber limit 26"; Ausnahme `27 | Grund da` ⇒ kein Befund; Ausnahme `27 |`
(leerer Grund) ⇒ **nur** `metrics-reason-empty`, die Überschreitung ist gedeckt (wie in der Spec); Ausnahme `26 | zu klein` ⇒
`metrics-over-limit`. Negativliste (§3) im Test: `make test` und `make gates` grün; stichprobenweise am Image: Pfad ohne
Vorkommen ⇒ 0 mit Hinweis, Kurz-/Voll-Hash, Basis = HEAD (Test), nicht eingecheckte Änderung zählt nicht (Test).

## (3) Fail-closed und Betrieb

- Ohne `.git`: `d-check: error: kein lesbares git-Repository unter /repo: stat /repo/.git: …` — **Exit 2**.
- Basis `deadbeef`: `basis "deadbeef" nicht auflösbar: reference not found (DC-FA-MET-001, fail-closed)` — **Exit 2**.
- `docker run --network none --user 12345:12345 -v …:ro`: Lauf gelingt, Exit 1 mit den erwarteten Befunden.
- `make image-test` selbst gefahren: Fall (6) „metrics — Messung unter :ro und fremdem Benutzer, Exit 1; ohne .git Exit 2", Gesamt `OK`, Exit 0.

## (4) Grund-Code

Zeile `SPEC-107 | metrics-over-limit | metrics | …` in `spec/spezifikation.md` §4 (Zeile 3959), `rules.ReasonMetricsOverLimit` in
`AllReasons()` und `reasonTexts()` (`internal/hexagon/core/app/diagnose.go`); `TestAllReasonsDeckungGegenSpezifikationGrundCodes`
liegt in der grünen `make test`-/`make gates`-Suite. **bestätigt.**

## (5) `make gates`

Selbst gefahren: `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`, `GATES_EXIT=0`. Nebenbefund: semgrep meldet
„3 timeout error(s) in … configyaml/configyaml.go" (Regeln `filepath-clean-misuse`, `raw-html-format`, `tainted-sql-string`),
Findings 0 — der Scan überspringt diese Regeln für die Datei; ob das neu ist (die Datei wuchs um 13 Zeilen), habe ich nicht
gegen den Vorzustand geprüft.

## Bewusstes Brechen (`make test`, Arbeitsbaum danach mit `git checkout --` zurückgenommen)

| Mutation in `internal/adapter/driven/git/git.go` | Ergebnis |
|---|---|
| `strings.HasPrefix(name, p+"/")` → `strings.HasPrefix(name, p)` (Segmentgrenze weg) | rot: `--- FAIL: TestLineChanges_LiteralePfade` — `Pfad "src/a" darf nichts treffen: {Lines:7 …}` (traf `src/a.txt`). Richtiger Grund. |
| `!underAny(name, paths) \|\| underAny(name, excluded)` → `!underAny(name, paths)` (Ausschluss weg) | rot: `--- FAIL: TestLineChanges_GleichNumstat` — `Summe = 29, erwartet 27`. Richtiger Grund (die 2 Zeilen aus `docs/vertraege/v.txt`). |

Beide Fixes tragen damit je einen benannten Test, der ohne sie aus dem behaupteten Grund rot läuft. Der Arbeitsbaum war nach dem Zurücknehmen sauber.

## (6) Abgrenzung und Unverändertheit ohne Konfiguration

- Lauf ohne `metrics`: `make blackbox-probe REF=v0.89.0` — „byte-identisch über 20 Vergleiche … stdout/stderr/Exit getrennt" (EXIT 0). Ebenso ein Lauf des Images ohne `--enable metrics` gegen das Fixture-Repo (mit `metrics:`-Block in der Config): 0 Befunde, Exit 0. **bestätigt.**
- `resolveVCS` wächst nur um `|| moduleActive(modules, "metrics")`; `runNoteModules` bekommt nur den `vcs`-Parameter; der Port wächst um genau
  eine Methode (+ `LineStat`). Verträge von `vcs`/`commits`/`tracked` nicht berührt. Keine Rename-Erkennung nach git-Art, keine
  Binärgewichtung, keine Release-Prep-Dateien im Diff. **bestätigt.**
- **Nicht ganz im Rahmen von §1 („Block-Lexik/Config-Leser … ändert nicht"):** `validMetricsPath` aus slice-283 ist nach `model.ValidLiteralPath` gewandert und
  dabei **verschärft** (leeres, `.`-Segment und `src//x` sind jetzt Fehler; ein Schrägstrich am Ende wird abgeschnitten) — Reaktion auf R1 F-6. Der Plan benennt die
  Verlagerung, nicht die Verhaltensänderung für die Pfade des Blocks (Exit 2 statt Messung 0). Die Spec-Historie nennt sie (SPEC-Historie 2026-10-30). Kein Schaden, aber eine
  Plan-Ausweitung ohne Plan-Änderungs-Zeile.

## Nicht bestätigt / Befunde

1. **Defekt im Lastenheft (neu im Fix-Commit):** `spec/lastenheft.md` Zeilen 4287 und 4288 tragen beide
   „Inhalt der Befunde und Hinweise." — doppelte Zeile im Absatz zur Fachlogik von `DC-FA-MET-001`. Kein Gate hat sie gefangen (`make gates` grün).
   Zu beheben, bevor die Spec (0.110.2) als Stand gilt.
2. **Plan-/Diff-Drift in der Tabelle §3:** geplant war `docs/user/abdeckung-tests.md`; geändert wurde `docs/user/abdeckung-e2e.md` (eine Zeile für
   Fall (6); sachlich richtig, denn der Beleg ist `tools/image-test.sh`). Ebenfalls im Diff, aber ohne Tabellenzeile: `spec/lastenheft.md` (0.110.1 → 0.110.2) und
   `docs/plan/adr/0112-metrics-modul.md` (ein Satz bei den Konsequenzen). Die Plan-Tabelle sollte bei der Closure nachgezogen oder in §7 benannt werden.
3. Die DoD-Zeile „Config-Schlüssel `metrics.excluded`" steht in der DoD, der Plan-Kopf zählt jedoch ≤ 3 Liefer-Punkte; DoD führt vier Liefer-Punkte (Port/Adapter, Prüfung, `excluded`, Fail-closed). Größen-Abweichung nach `MR-066`: §7 sollte Grund und Ersatz-Form der Prüfung (ein Reviewer-Durchlauf, R1) festhalten.

## Arbeitsbaum

Außer diesem Report unverändert (Befehl `git status --short` zeigt nur diese Datei als neu).
