# Verifikation R2 — slice-039 Fixes (PR-/Push-CI + Traceability-Gate)

## Kopf-Metadaten

- **Review-Art:** Unabhängige Verifikations-/Review-Runde **R2** (adversariale
  Bestätigung/Widerlegung der R1-Fixes + Regressions-Suche). Frischer Kontext,
  alles selbst reproduziert. **Kein Verifier im DoD-Sinn:** `make gates`/
  `make doc-check` werden NICHT als grün angenommen; Gate-Aussagen des Commits
  werden nicht übernommen. Adversariale Tests read-only in `/tmp`-Wegwerf-Repos
  (keine Änderung/kein Commit im d-check-Repo).
- **Datum:** 2026-06-21
- **Gegenstand:** Fix-Commit `ba50813`
  („slice-039 Review R1 — 2 HIGH + 1 MEDIUM behoben"). Geändert:
  `tools/trace-check.sh` (`clean_message`, fail-closed-Range, `MR-`-Muster),
  `.github/workflows/ci.yml` (Zero-SHA → `origin/<default-branch>`),
  `docs/plan/adr/0013-pr-ci-und-traceability-gate.md` (Muster-Block).
  Bezug: R1-Report
  [`docs/reviews/2026-06-21-slice-039-pr-ci-traceability.md`](2026-06-21-slice-039-pr-ci-traceability.md)
  (2 HIGH + 1 MEDIUM + 1 LOW + 2 INFO).
- **Reviewer:** Claude (Agent), Skill
  [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) v1.0.0.
- **Modell:** `claude-opus-4-8[1m]`.
- **Eingangs-Kontext:** R1-Report (HIGH-1/HIGH-2/MEDIUM-1/LOW-1/INFO-1/INFO-2);
  Fix-Diff `ba50813`; aktueller Stand `tools/trace-check.sh`,
  `.github/workflows/ci.yml`, ADR-0013, `.d-check.yml` (`ids`). Hinweis: ein
  Folge-Commit `8cbc3d8` berührt nur `harness/`-Doku und `AGENTS.md`
  (MR-`*`-Wortlaut), nicht das geprüfte Gate-Skript.

## Reproduktions-Umgebung

Alle Belege aus `bash`-Läufen gegen `tools/trace-check.sh` (Stand nach
`ba50813`) in `/tmp`-Wegwerf-Repos. `bash -n tools/trace-check.sh` → exit 0
(Syntax sauber).

## Befund-Status der R1-Findings

### HIGH-1 — Silent-Green bei neuem Branch — **BESTÄTIGT BEHOBEN**

- **Kategorie:** HIGH (R1) → behoben in R2
- **Quelle:** R1 HIGH-1; [ADR-0013](../plan/adr/0013-pr-ci-und-traceability-gate.md)
  Punkt 2 (CI-Range-Check als „robuster, klon-unabhängiger Backstop")
- **Pfad:** `tools/trace-check.sh:75` (Resolvability-Gate, Zeilen 75–78);
  `.github/workflows/ci.yml:58` (Zero-SHA-Zweig, Zeilen 58–62)
- **Beleg (reproduziert):** Wegwerf-Repo mit drei Branch-Commits (zwei ohne ID,
  HEAD `slice-099`). `trace-check.sh --range <Zero-SHA>..<HEAD>` endet jetzt mit
  **Exit 2** und „Range-Basis … nicht auflösbar" — statt dem alten still-grünen
  „1 Commit(s)"/exit 0. Auch eine **nicht-null, aber unauflösbare** Basis
  (`deadbeef…..HEAD`) und eine **fehlgeschlagene CI-Fetch-Situation**
  (`origin/main` nicht vorhanden, mimt `git fetch … || true`-Fehler in
  `.github/workflows/ci.yml:61`) enden mit **Exit 2** (fail-closed). Ein
  normaler auflösbarer Range mit ID-losem **Zwischen**-Commit endet mit Exit 1
  (Zwischen-Commit benannt). Der „nur-HEAD"-Fallback ist entfernt
  (`tools/trace-check.sh:79` listet jetzt unbedingt `git rev-list --no-merges
  "$range"` nach der Resolvability-Gate). Restpfad force-push: `before` ist eine
  reale, lokal noch auflösbare SHA → Range wird gebildet (über-approximiert ggf.
  mehr Commits, nie weniger) → kein Silent-Green.
- **Verifizierbar:** ja — `trace-check.sh --range 000…0..<HEAD>` → exit 2;
  `--range origin/main..<HEAD>` ohne `origin` → exit 2.

### HIGH-2 — Hook/CI-Divergenz (#-Kommentar/scissors) — **BESTÄTIGT BEHOBEN (für den Editor-Pfad)**; siehe NEU HIGH-A für den `-m`-Pfad

- **Kategorie:** HIGH (R1) → für den Editor-Pfad behoben in R2
- **Quelle:** R1 HIGH-2; [ADR-0013](../plan/adr/0013-pr-ci-und-traceability-gate.md)
  Punkt 3 („eine Wahrheit")
- **Pfad:** `tools/trace-check.sh:32` (`clean_message`), `:63`
- **Beleg (reproduziert):** Eine ID **nur** in einer `#`-Kommentarzeile
  (Message-Datei) → `trace-check.sh --message` endet jetzt mit **Exit 1**
  (vorher grün). Eine ID **nur** unterhalb der scissors-Zeile
  `# ------------------------ >8 ------------------------` (verbose `-v`) →
  **Exit 1**. `clean_message` (`sed -e '/^#.*>8/,$d' -e '/^#/d'`) hält dabei
  legitime Body-Zeilen mit `>8` (z. B. „throughput >8 ops") korrekt und strippt
  nur Kommentar-/Scissors-Inhalt. Eine real-committete ID im echten Betreff
  passiert weiterhin grün (Editor-Cleanup `strip`). Für den **Editor-Commit**
  (git-Default `--cleanup=default` = `strip`, da edited) bewertet der Hook jetzt
  dasselbe wie git committet — Divergenz in dieser Richtung geschlossen.
- **Verifizierbar:** ja — `--message` mit ID nur in `# …`-Zeile → exit 1; ID nur
  unter `>8`-Zeile → exit 1.

### MEDIUM-1 — Muster nicht deckungsgleich (`MR-\d{3}` fehlt) — **BESTÄTIGT BEHOBEN**

- **Kategorie:** MEDIUM (R1) → behoben in R2
- **Quelle:** R1 MEDIUM-1; [ADR-0013](../plan/adr/0013-pr-ci-und-traceability-gate.md)
  Punkt 4 (Muster = `.d-check.yml`-`ids` plus `slice-*`)
- **Pfad:** `tools/trace-check.sh:19` gegen `.d-check.yml:17`,`:21`,`:25`
- **Beleg (reproduziert):** `ID_RE` ist jetzt
  `(ADR-[0-9]{4}|MR-[0-9]{3}|DC-(FA-[A-Z]+|QA)-[0-9]+|slice-[0-9]+)` — die drei
  `.d-check.yml`-`ids`-Muster (`ADR-\d{4}`, `MR-\d{3}`, `DC-(FA-[A-Z]+|QA)-\d+`)
  plus `slice-`. Ein MR-only-Commit („docs: nachtrag MR-008 zur baseline") →
  `--message` **Exit 0** (vorher abgelehnt). Die ADR-0013-Prosa
  (`docs/plan/adr/0013-pr-ci-und-traceability-gate.md:99`, Zeilen 99–108) listet
  nun `MR-*` explizit und ist deckungsgleich mit `.d-check.yml`. Boundary: `MR-12`
  (zweistellig) matcht korrekt **nicht**.
- **Verifizierbar:** ja — `--message` mit „MR-008" → exit 0; Regex-Vergleich
  `tools/trace-check.sh:19` ↔ `.d-check.yml:21`.

## NEUES Finding (durch den Fix eingeführt)

### HIGH-A — `clean_message` über-strippt für `git commit -m`/`-F` (cleanup=*whitespace*); ein `#`-führender Betreff mit der einzigen ID wird vom Hook abgelehnt, obwohl git ihn samt ID committet — Hook/CI-Divergenz in der **Gegenrichtung**

- **Kategorie:** HIGH
- **Quelle:** Reviewer-Skill (Repo-Anker HIGH: „Stilles-Grün-Pfad in einem
  Gate-Skript — Harness-Lüge"; hier die symmetrische Lüge: stilles **Rot** für
  einen real konformen Commit); [ADR-0013](../plan/adr/0013-pr-ci-und-traceability-gate.md)
  Punkt 3 („eine Wahrheit, keine Logik-Dopplung … sonst divergiert der Hook"
  vom maßgeblichen committeten Inhalt)
- **Pfad:** `tools/trace-check.sh:32` (`clean_message`), `:63`
- **Befund:** `clean_message` (`sed … -e '/^#/d'`) entfernt **alle** mit `#`
  beginnenden Zeilen. Das modelliert git's `--cleanup=strip` (Editor-Pfad),
  aber **nicht** den `-m`/`-F`-Pfad: git's `default`-Cleanup ist laut
  `git-commit(1)` „strip if the message is to be edited, **otherwise
  whitespace**", und `whitespace` lässt `#`-Kommentar-Zeilen **stehen**. Trägt
  die einzige ID eine `#`-Zeile, die git via `-m`/`-F` **behält**, lehnt der
  Hook den Commit ab, während die committete Message die ID enthält und der
  CI-Range-Check (`git log -1 --format=%B`, `tools/trace-check.sh:84`) sie
  akzeptiert. Konkret: `git commit -m "#123 feat slice-099"` (Issue-/PR-Ref-
  Betreff, gängige Konvention) wird vom installierten Hook mit „nennt keine
  DC-/ADR-/slice-ID" und leerem Betreff (`> `) **blockiert**, obwohl git die
  Message verbatim speichert und CI sie als gültig (`slice-099`) bestätigt. Das
  ist eine Hook/CI-Divergenz derselben Klasse, die HIGH-2 schließen sollte — nur
  in der Gegenrichtung (falsches Rot statt falsches Grün). Der Happy-Path (ID im
  Nicht-`#`-Betreff, `#`-Zeile nur im Body) bleibt grün; betroffen ist nur der
  Fall „einzige ID auf einer `#`-Zeile bei `-m`/`-F`".
- **Verifizierbar:** ja — Wegwerf-Repo mit installiertem `commit-msg`-Hook
  (`core.hooksPath`): `git commit -m "#123 feat slice-099"` → Hook exit 1,
  Commit abgebrochen; dasselbe ohne Hook → committet, `git log -1 --format=%B`
  enthält `slice-099`, `grep -qE "$ID_RE"` → exit 0 (CI würde akzeptieren).
  Ebenso `git commit -m "subject ohne id" -m "# ADR-0013 ref"` → Hook exit 1,
  während die gespeicherte Message (cleanup=whitespace) `ADR-0013` trägt.

## Negativbefunde (geprüft, ohne Befund)

- **HIGH-1 Restpfade vollständig fail-closed:** Zero-SHA-Basis, nicht-null-
  unauflösbare Basis und fehlgeschlagener `origin/<default>`-Fetch enden alle
  mit Exit 2; der alte „nur-HEAD"-Zweig ist im Diff entfernt
  (`tools/trace-check.sh:75`, Zeilen 75–79). Kein verbliebener Pfad gefunden, auf
  dem ID-lose Zwischen-Commits still grün durchrutschen.
- **`set -euo pipefail`-Verhalten der neuen Zweige:** Der neue `if … then echo
  …; exit 2; fi`-Block (`tools/trace-check.sh:75`, Zeilen 75–78) und `clean_message`
  (`tools/trace-check.sh:32`, im `$( … )` vor `check_msg`) lösen keinen
  vorzeitigen `-e`-Abbruch mit falschem Status aus. Tests: leere Message-Datei →
  exit 1; nicht-existente `--message`-Datei → exit 2; fehlendes `--range`-Arg →
  exit 2; unbekannter Modus → exit 2; `--self-test` → exit 0. Kein euo-Crash mit
  falschem Grün.
- **Self-Test weiterhin schlüssig:** `self_test`
  (`tools/trace-check.sh:46`, Zeilen 46–54) prüft ID-erkannt /
  fehlende-ID-MUSS-failen /
  Merge-ausgenommen und läuft in jedem nicht-trivialen Modus vorab; mit dem neu
  ergänzten `MR-`-Muster und `clean_message` bleibt der Negativ-Selbsttest
  intakt (`--self-test` grün, exit 0). Eine fälschlich durchgewunkene fehlende
  ID bräche ihn mit exit 2.
- **Muster-Boundary unverändert (LOW-1 weiter offen, aber konsistent):** `MR-12`
  (2-stellig) matcht nicht; `MR-1234`, `ADR-00000`, `myslice-1`, `fooDC-QA-1`,
  `xADR-1234` matchen weiter (unverankertes Substring-Verhalten). Das ist die in
  R1-LOW-1 dokumentierte, mit `.d-check.yml` (`ids`) konsistente Eigenschaft —
  vom Fix weder verbessert noch verschlechtert; kein neues Finding.
- **`clean_message` strippt keine legitimen Body-Zeilen mit `>8`:** Die
  scissors-Regel `'/^#.*>8/'` greift nur auf `#`-führende Zeilen; eine
  Nicht-Kommentar-Body-Zeile „throughput >8 ops" bleibt erhalten. Indentierte
  (`  #…`) Kommentarzeilen werden von `clean_message` **nicht** gestrippt —
  konsistent mit git, das nur `#` in Spalte 0 als Kommentar behandelt.
- **MEDIUM-1-Doku-Sync:** `docs/plan/adr/0013-pr-ci-und-traceability-gate.md:94`
  listet `MR-\d{3}` im Muster-Block; die Prosa
  (`docs/plan/adr/0013-pr-ci-und-traceability-gate.md:99`, Zeilen 99–108) nennt
  `MR-*` explizit als Konventions-ID. Deckungsgleichheits-Behauptung jetzt
  zutreffend.
- **CI-Workflow-Logik (`.github/workflows/ci.yml`):** `pull_request` →
  `PR_BASE..HEAD`; Zero-SHA-`before` → `origin/<default>..HEAD` nach
  `git fetch`; sonst → `before..HEAD` (`.github/workflows/ci.yml:56-64`). Der
  `git fetch … || true` (`:61`) verschluckt zwar Fetch-Fehler, aber ein dann
  unauflösbares `origin/<default>` führt im Skript zum fail-closed Exit 2 (oben
  belegt) — kein Silent-Green-Restpfad über diesen swallow.
- **Kein neuer Doppellauf / keine Permission-Regression:** `ci.yml`-Trigger
  (`pull_request`, `push` Branches, `tags-ignore`), `permissions: {}` top-level
  + `contents: read` im Job und der SHA-Pin von `actions/checkout` sind im
  Fix-Diff `ba50813` **unverändert** — die R1-Negativbefunde dazu bleiben gültig.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 1 | 0 | 0 | 0 |

(R1-Findings HIGH-1, HIGH-2, MEDIUM-1: bestätigt behoben — nicht erneut gezählt.
LOW-1/INFO-1/INFO-2 aus R1 bleiben offen/won't-fix wie dort eingestuft, hier
nicht neu gezählt. **NEU:** HIGH-A.)

## Verdikt

**Blockiert (nicht closure-/accept-fähig).** Die drei R1-Befunde sind
**bestätigt behoben**: HIGH-1 ist fail-closed (Zero-SHA/unauflösbare Basis →
Exit 2 statt nur-HEAD-Grün; CI liefert `origin/<default>` für neue Branches),
HIGH-2 ist für den Editor-Pfad geschlossen (`clean_message` entfernt `#`-/
scissors-Inhalt, ID nur im Kommentar → Exit 1), MEDIUM-1 ist geschlossen
(`MR-\d{3}` ergänzt, MR-only → Exit 0, ADR-Deckungsgleichheit nun zutreffend).

Der Fix für HIGH-2 führt jedoch **HIGH-A** ein: `clean_message` modelliert
git's `strip`-Cleanup (Editor) und über-strippt damit für den `-m`/`-F`-Pfad
(git-Default `whitespace`, `#`-Zeilen bleiben). Ein `#`-führender Betreff mit der
einzigen ID — z. B. das gängige `git commit -m "#123 feat slice-099"` — wird vom
lokalen Hook **abgelehnt**, obwohl git die ID committet und der CI-Range-Check
sie akzeptiert: dieselbe Hook/CI-Divergenz gegen die zugesicherte „eine
Wahrheit", nur als falsches Rot statt falsches Grün. Reproduziert mit
installiertem Hook (Commit blockiert) vs. `git log -1 --format=%B` (ID präsent,
CI akzeptiert).

LOW-1 (unverankerte Substring-Treffer, konsistent mit `.d-check.yml`) und
INFO-1/INFO-2 aus R1 bleiben wie dort eingestuft (nicht blockierend). Die
Gate-Bestätigung (`make gates`/`make doc-check`) obliegt der getrennten
Verifikation; hier nicht als grün angenommen. **ADR-0013 → Accepted erst nach
Schließung von HIGH-A** (Cleanup-Modus-Annahme von `clean_message`
vereinheitlichen mit dem maßgeblichen committeten Inhalt).
