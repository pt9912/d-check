# Review — slice-039 Implementierung (PR-/Push-CI + Traceability-Gate)

## Kopf-Metadaten

- **Review-Art:** Code-Review (Commit-Diff gegen Plan/ADR/Anforderungen/
  Hard Rules — **kein Verifier**: DoD-Abhaken und Gate-Lauf-Bestätigung sind
  nicht Gegenstand; Gates werden NICHT als grün angenommen). Adversariale
  Tests des Gates wurden read-only in `/tmp`-Wegwerf-Repos gefahren (keine
  Änderung am d-check-Repo, keine Commits).
- **Datum:** 2026-06-21
- **Gegenstand:** Commit `3e256b6` (slice-039). Neu: `tools/trace-check.sh`,
  `.githooks/commit-msg`, `.github/workflows/ci.yml`. Geändert: `Makefile`
  (Targets `trace-check`, `hooks`, `.PHONY`), `AGENTS.md` (§4-Zeilen, §5
  Commit-Disziplin um `slice-*`), `harness/README.md` (§Sensors-Zeilen,
  §Traceability rules um `slice-*`), Slice-Plan
  `docs/plan/planning/done/slice-039-pr-ci-traceability-gate.md`
  (Status next→in-progress, DoD-Häkchen).
- **Reviewer:** Claude (Agent), Skill `.harness/skills/reviewer.md` v1.0.0.
- **Modell:** `claude-opus-4-8[1m]`.
- **Eingangs-Kontext:** Slice-Plan
  `docs/plan/planning/done/slice-039-pr-ci-traceability-gate.md`; ADR
  [ADR-0013](../plan/adr/0013-pr-ci-und-traceability-gate.md) (Status
  Proposed; Bezug [ADR-0011](../plan/adr/0011-digest-pins-build-gate-images.md),
  [ADR-0002](../plan/adr/0002-distribution-ghcr-image.md)); Anforderungen
  [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) (make-ci-Vertrag
  bei jeder Integration) und
  [`DC-FA-DIST-001`](../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
  (`make ci` ⊃ `make image-test`); Hard Rules `AGENTS.md` §3 (insb. §3.1
  Docker/make-only, §3.3 Lifecycle); Meta-Gate `tools/gate-consistency.sh`;
  Muster-Quelle `.d-check.yml` (`ids`); Doppellauf-Abgrenzung
  `.github/workflows/release.yml`. Prozess-/Durchsetzungs-ADR ohne eigene
  `DC-FA-*`-Vertragsbindung.

## Findings

### HIGH

#### HIGH-1 — Neuer-Branch-Push prüft nur EINEN Commit; ID-lose Zwischen-Commits rutschen still grün durch

- **Kategorie:** HIGH
- **Quelle:** Reviewer-Skill (Repo-Anker HIGH: „Stilles-Grün-Pfad in einem
  Gate-Skript — Harness-Lüge"); [ADR-0013](../plan/adr/0013-pr-ci-und-traceability-gate.md)
  Punkt 2 (CI-Check „prüft den Commit-Range … robust, klon-unabhängig" als
  Backstop) und Fitness Function („rot bei einem Commit/Range ohne ID")
- **Pfad:** `tools/trace-check.sh:62` (Fallback-Zweig, Zeilen 62–63)
- **Befund:** Im `--range`-Modus prüft das Skript bei Zero-SHA-Basis **oder**
  nicht-auflösbarer Basis nicht den ganzen Push, sondern nur
  `git rev-list --no-merges -n1 HEAD` — also genau **einen** Commit. Bei
  einem CI-`push`-Event ist `github.event.before` für den **ersten Push eines
  neuen Branches** die Zero-SHA (`000…`,
  `.github/workflows/ci.yml:51`,`:58`); damit greift dieser Fallback im
  häufigsten Mehr-Commit-Fall. Ein neuer Branch mit mehreren Commits, von
  denen nur HEAD eine ID trägt, passiert grün — die ID-losen
  Zwischen-Commits werden nie gegen das Muster gehalten. Der CI-Range-Check
  ist laut ADR der „robuste, klon-unabhängige Backstop"; in diesem Pfad
  liefert er Sicherheit, die er nicht prüft.
- **Verifizierbar:** ja — Wegwerf-Repo mit drei neuen Branch-Commits (zwei
  ohne ID, HEAD mit `slice-099`),
  `trace-check.sh --range 000…0..<HEAD>` endet mit `exit 0` und meldet
  „1 Commit(s)"; `git rev-list --no-merges -n1 HEAD` zeigt, dass nur der
  HEAD-Commit geprüft wurde, die beiden ID-losen Commits ungeprüft blieben.

#### HIGH-2 — commit-msg-Hook greppt die ROHE Message-Datei; ID nur in `#`-Kommentar/`-v`-Diff täuscht Grün vor (Hook ≠ CI, „eine Wahrheit" verletzt)

- **Kategorie:** HIGH
- **Quelle:** Reviewer-Skill (Repo-Anker HIGH: „Stilles-Grün-Pfad in einem
  Gate-Skript — Harness-Lüge"); [ADR-0013](../plan/adr/0013-pr-ci-und-traceability-gate.md)
  Punkt 3 („Eine Wahrheit, keine Logik-Dopplung … sonst divergiert der Hook")
- **Pfad:** `.githooks/commit-msg:12` zusammen mit `tools/trace-check.sh:53`
  (`check_msg "$(cat "$2")"`)
- **Befund:** Git übergibt dem `commit-msg`-Hook die **noch nicht
  bereinigte** Message-Datei (inkl. `#`-Kommentarzeilen und, bei
  `git commit -v`, des kompletten Diffs unter der Scissors-Linie). Das
  Skript greppt diese Datei unverändert (keine `git stripspace
  --strip-comments`/Scissors-Bereinigung, `tools/trace-check.sh:53`,`:18`,`:23`).
  Steht eine ID nur in einer `#`-Kommentarzeile (z. B. in der git-Template-
  Hilfe versehentlich notiert) oder im `-v`-Diff-Anteil, zählt der Hook sie
  als gültig und winkt den Commit grün durch. Git entfernt diese Zeilen
  beim Editor-Commit (`--cleanup=strip`, Default) **nach** dem Hook — die
  committete Message hat dann keine ID. Der CI-Range-Check liest die
  bereinigte Message (`git log --format=%B`, `tools/trace-check.sh:71`,`:82`)
  und lehnt **denselben** Commit ab. Hook (feedforward) und CI (feedback)
  bewerten denselben Commit gegensätzlich, obwohl ADR-0013 ausdrücklich
  „eine Wahrheit" zusichert; der lokale Wächter erzeugt ein stilles Grün auf
  dem Pfad, der den Autor vor dem Push warnen soll.
- **Verifizierbar:** ja — Wegwerf-Repo mit installiertem Hook
  (`core.hooksPath`): `git commit --cleanup=strip -F <datei>` mit
  Betreff ohne ID und einer `# … ADR-0013`-Kommentarzeile committet grün
  (Hook exit 0); `git log -1 --format=%B` zeigt die bereinigte Message
  **ohne** ID; `trace-check.sh --range 000…0..<HEAD>` gegen genau diesen
  Commit endet mit `exit 1`.

### MEDIUM

#### MEDIUM-1 — Muster ist NICHT deckungsgleich mit `.d-check.yml` (`MR-\d{3}` fehlt); ADR-Behauptung der Deckungsgleichheit unzutreffend

- **Kategorie:** MEDIUM
- **Quelle:** [ADR-0013](../plan/adr/0013-pr-ci-und-traceability-gate.md)
  Punkt 4 („Kennungs-Muster = eine Quelle … deckungsgleiche Muster mit
  `.d-check.yml` (`ids`), damit ‚gültige ID' nur **einmal** definiert ist");
  Skript-Header (`tools/trace-check.sh:17` „deckungsgleich mit .d-check.yml")
- **Pfad:** `tools/trace-check.sh:18` gegen `.d-check.yml:15-28`
- **Befund:** `.d-check.yml` (`ids.patterns`) definiert **drei**
  ID-Klassen: `ADR-\d{4}`, `MR-\d{3}` und `DC-(FA-[A-Z]+|QA)-\d+`
  (`.d-check.yml:17`,`:22`,`:25`). Das Gate-Muster
  (`tools/trace-check.sh:18`) führt `ADR`/`DC` plus `slice-[0-9]+`, lässt
  aber `MR-\d{3}` weg. Damit ist die im Skript-Header und in ADR-0013
  Punkt 4 behauptete Deckungsgleichheit mit `.d-check.yml` faktisch nicht
  gegeben. Folge: ein Commit, der ausschließlich eine `MR-*`-Konventions-ID
  nennt (z. B. eine Änderung an `harness/conventions.md`, die nur `MR-008`
  referenziert), wird von Hook und CI mit „keine DC-/ADR-/slice-ID"
  abgewiesen, obwohl `MR-*` eine reguläre, definierte Kennung des Repos ist.
  Aktuell zeigt die Historie der jüngsten 50 Commits keinen MR-only-Commit,
  also kein laufendes Versagen — aber eine latente, von der ADR-Aussage
  verdeckte Falle. (ADR-0013 Punkt 4 listet selbst nur die drei Muster ohne
  `MR`; der Widerspruch liegt zwischen der ADR-Prosa „deckungsgleich mit
  `.d-check.yml`" und dem tatsächlichen Umfang von `.d-check.yml`.)
- **Verifizierbar:** ja — `trace-check.sh --message <datei>` mit Inhalt
  „docs: nachtrag MR-008 zur baseline" endet mit `exit 1`; Regex-Vergleich
  `tools/trace-check.sh:18` gegen `.d-check.yml:22` belegt die fehlende
  `MR-\d{3}`-Alternative ohne Gate-Lauf.

### LOW

#### LOW-1 — Unverankertes Muster matcht Teilwörter (`myslice-1`, `ADR-00000`, `fooDC-QA-1`); Gate-Umgehung über zufällige Substrings möglich

- **Kategorie:** LOW
- **Quelle:** „Maintainability" (latente Wartungsfalle); [ADR-0013](../plan/adr/0013-pr-ci-und-traceability-gate.md)
  Punkt 4 (Muster als Definition „gültige ID")
- **Pfad:** `tools/trace-check.sh:18`,`:23`
- **Befund:** Das Muster trägt keine Wortgrenzen-/Anker. `grep -qE`
  akzeptiert deshalb auch Teilwort-Treffer: `myslice-1` (Substring
  `slice-1`), `ADR-00000` (`[0-9]{4}` matcht vier der fünf Ziffern),
  `xADR-1234`, `fooDC-QA-1`, `ADR-1234x`. Eine Commit-Message, die keine
  echte Kennung trägt, aber zufällig (oder absichtlich) eine solche
  Zeichenfolge enthält, passiert das Gate. Die Schärfe der Durchsetzung ist
  damit geringer als die Prosa „eine gültige ID" suggeriert. Eingestuft als
  LOW, nicht MEDIUM/HIGH, weil dasselbe unverankerte Verhalten auch in
  `.d-check.yml` (`ids`) gilt — das Gate ist insofern konsistent mit der
  bestehenden Muster-Konvention; ein Schärfen müsste beide Quellen treffen
  (Steering-Signal, nicht isolierter Slice-Defekt).
- **Verifizierbar:** ja — `grep -qE '(ADR-[0-9]{4}|DC-(FA-[A-Z]+|QA)-[0-9]+|slice-[0-9]+)'`
  gegen `myslice-1`, `ADR-00000`, `fooDC-QA-1` liefert je Exit 0 (Treffer).

### INFO

#### INFO-1 — Leerer/0-Commit-Range endet grün; eine fehlberechnete Range, die 0 Commits liefert, wird still grün

- **Kategorie:** INFO
- **Quelle:** dokumentationswürdige Annahme; [ADR-0013](../plan/adr/0013-pr-ci-und-traceability-gate.md)
  Re-Evaluierungs-Trigger („Commit-Range in CI unzuverlässig")
- **Pfad:** `tools/trace-check.sh:67` (Schleifen-/Exit-Block, Zeilen 67–74)
- **Befund:** Liefert die Range 0 Commits (z. B. `base==head`, oder eine in
  CI fehlberechnete Range), durchläuft die `while`-Schleife nullmal,
  `fail=0` bleibt, das Skript meldet „0 Commit(s) tragen eine … ID" und
  endet mit `exit 0`. Das ist für einen echten Leer-Range korrekt (nichts
  zu prüfen), aber es gibt keinen Hinweis/Schwellwert, der einen
  **unerwartet** leeren Range (Fehlberechnung) von einem echten Leer-Push
  unterscheidet — ein stiller-Grün-Pfad bei kaputter Range-Bildung. Bewusst
  als INFO statt HIGH notiert: Leer-Range ist kein Defekt für sich, und der
  ADR nennt die Range-Zuverlässigkeit bereits als Re-Evaluierungs-Trigger.
- **Verifizierbar:** ja — `trace-check.sh --range <H>..<H>` endet mit
  „0 Commit(s)" und `exit 0`.

#### INFO-2 — Hook/Skript laufen auf dem Host (bash+git); §3.1 nicht verletzt, aber dünner Stolperdraht statt make-only-Gate

- **Kategorie:** INFO
- **Quelle:** `AGENTS.md` §3.1 (Docker/make-only); Slice-Plan §4 (bewusste
  Designnotiz „dünner Stolperdraht")
- **Pfad:** `.githooks/commit-msg:10-12`, `tools/trace-check.sh:15`,`:20`
- **Befund:** Hook und Prüf-Skript laufen direkt auf dem Host über `bash`
  und `git`, nicht in Docker. §3.1 verbietet „Host-Go und Host-Paketmanager
  (`go`, `pip`, `npm`, `cargo`, `apt`, `brew`)" und stellt fest „Der Host
  braucht nur `git`, GNU `make`, `bash` und Docker" (`AGENTS.md:45-50`).
  `bash`+`git` sind die explizit erlaubten Host-Werkzeuge — kein Verstoß;
  Präzedenz: `tools/gate-consistency.sh`, `.claude/hooks/*` laufen ebenso
  auf dem Host. Bewusst notiert als die im Slice-Plan §4 selbst benannte
  Restspannung (dünner Stolperdraht, robuste Kontrolle bleibt CI). REFUTED
  als Hard-Rule-Verstoß mit §3.1-Zitat.
- **Verifizierbar:** nein (Designnotiz; keine Gate-Aussage).

## Negativbefunde (geprüft, ohne Befund)

- **`make ci`/`make gates` unverändert (trace-check NICHT eingehängt):** Der
  Diff von `3e256b6` ändert die Target-Zeilen `gates:` / `ci:` / `fullbuild:`
  nicht; `gates: doc-check lint test arch-check coverage-gate semgrep
  gate-consistency record-gates` (`Makefile:121`) und `ci: gates image-test`
  (`Makefile:128`) tragen kein `trace-check`. ADR-0013 Punkt 5
  (Commit-Zeit-Bindepunkt getrennt vom Arbeitsbaum-Inhalt) ist eingehalten.
  Verifiziert per `grep` über `Makefile` und über den Commit-Diff.
- **Hook real installiert (Faktencheck):** `git config --get core.hooksPath`
  liefert `.githooks` (exit 0) — der versionierte `commit-msg`-Hook ist im
  Arbeits-Klon aktiv; `make hooks` (`Makefile:hooks`-Target) setzt genau
  diesen Wert. Der Hook-Vertrag ist korrekt: `$1` = Message-Datei,
  `cd "$(git rev-parse --show-toplevel)"` vor dem Aufruf, `exec bash
  tools/trace-check.sh --message "$1"` (`.githooks/commit-msg:11-12`). Bei
  fehlender ID bricht der Commit real ab (Wegwerf-Repo-Test: Betreff ohne
  ID → Hook exit 1 → Commit abgebrochen).
- **gate-consistency beidseitig grün über die neuen Targets:** `bash
  tools/gate-consistency.sh` läuft real grün (exit 0, „Doku ↔ Makefile
  konsistent"). `make trace-check` und `make hooks` stehen sowohl in
  `AGENTS.md` §4 (`AGENTS.md:109`,`:110`) als auch in `harness/README.md`
  §Sensors — beide Richtungen (Doku→Makefile, Makefile→AGENTS §4) erfüllt.
  `.PHONY` (`Makefile:41`) listet `trace-check` und `hooks`. Die
  `.d-check.yml`-Modulliste ist unverändert (`links, anchors, ids, matrix,
  codepaths, spans, hostpaths`, `.d-check.yml:7`) — Check (3) intakt.
- **Merge-/Revert-Ausnahme greift nur auf die erste Zeile:** `is_exempt`
  nutzt `head -n1 … | grep -qE '^(Merge |Revert )'` (`tools/trace-check.sh:22`).
  Test: `Merge branch 'foo'` → exit 0; `Revert "…"` → exit 0; eine Message
  mit „Merge" nur im Body (zweite Zeile) und ohne ID → exit 1; „Reverting …"
  (kein Leerzeichen nach `Revert`) → exit 1. Der Anker `^` plus das
  Leerzeichen verhindern Body-/Präfix-Umgehung.
- **Self-Test beweist Fang einer fehlenden ID:** `self_test`
  (`tools/trace-check.sh:37`, Zeilen 37–45) prüft drei Fälle: ID erkannt, fehlende ID
  **muss** fehlschlagen (`if check_msg "chore: ohne bezug" …; then … exit 2`,
  `:40-42`), Merge ausgenommen. Würde `check_msg` eine fehlende ID
  fälschlich durchwinken, bricht der Self-Test mit exit 2 ab. Er läuft in
  jedem nicht-trivialen Modus (`--message`,`--range`,`--self-test`,leer)
  vorab. Schlüssig.
- **`set -euo pipefail`-Fallen:** `is_exempt`/`has_id` werden ausschließlich
  in `&&`/`||`/`return`-Kontexten und in `if`-Bedingungen ausgewertet, wo
  `set -e` einen Nicht-Null-Status nicht zum Abbruch eskaliert; `check_msg`
  im Range-Modus ist mit `|| fail=1` abgefangen (`:71`), im
  `--message`/leeren Modus als letzte/abschließende Anweisung bewusst
  abbruch-wirksam. Tests: leere Message-Datei → exit 1 (kein euo-Crash),
  nicht-existente `--message`-Datei → exit 2, fehlendes `--range`-Argument →
  exit 2, unbekannter Modus → exit 2. Kein vorzeitiger euo-Abbruch mit
  falschem Grün beobachtet.
- **Range-Modus mit echtem Range fängt Zwischen-Commits:** Bei auflösbarer
  Basis nutzt das Skript `git rev-list --no-merges "$range"`
  (`tools/trace-check.sh:65`) und iteriert über **alle** Commits; ein
  ID-loser Zwischen-Commit setzt `fail=1` und liefert exit 1 (Test:
  3-Commit-Range mit ID-losem mittleren Commit → exit 1). Die HIGH-1-Lücke
  betrifft ausschließlich den Zero-SHA-/unauflösbar-Fallback-Pfad, nicht den
  normalen Range.
- **CI-Trigger ohne Doppellauf mit `release.yml`:** `ci.yml` triggert
  `pull_request` und `push` auf `branches: ['**']` mit `tags-ignore:
  ['**']` (`.github/workflows/ci.yml:18-24`); `release.yml` triggert
  ausschließlich `push: tags: ['v*']`
  (`.github/workflows/release.yml:21-24`). Tag-Pushes laufen damit nur in
  `release.yml`, Branch-/PR-Ereignisse nur in `ci.yml` — disjunkt, kein
  Doppellauf.
- **CI-Permissions minimal; Fork-PR-sicher:** `ci.yml` setzt top-level
  `permissions: {}` und im Job nur `contents: read`
  (`.github/workflows/ci.yml:27`,`:34-35`) — kein `packages:write`/
  `contents:write` wie in `release.yml`. Event ist `pull_request` (nicht
  `pull_request_target`), die Range nutzt
  `github.event.pull_request.base.sha` + `github.sha` (`:50`,`:52`); Fork-PRs
  laufen ohne Push-Secrets und ohne Schreibrechte. `fetch-depth: 0` (`:44`)
  stellt die volle Historie für den Range bereit. Action `actions/checkout`
  ist SHA-gepinnt (`@de0fac2…`, `:41`, identischer Pin wie `release.yml:46`).
- **Reihenfolge fail-fast:** Das Traceability-Gate (`make trace-check
  RANGE=…`) läuft als erster Step vor dem teuren `make ci`
  (`.github/workflows/ci.yml:47-61` vor `:64-65`) — ein fehlender ID-Befund
  spart den Docker-Build.
- **Lifecycle §3.3:** Der reine Move next→in-progress geschah in eigenem
  Commit (`a6e7db1`, „reiner Move, §3.3"); `3e256b6` ändert nur den Inhalt
  der bereits verschobenen Slice-Datei (Status + Häkchen) — keine
  git-mv-+-Inhalt-Kombination, §3.3 nicht verletzt. ADR-0013 ist im
  [ADR-Index](../plan/adr/README.md) eingetragen (Status Proposed,
  deckungsgleich mit der ADR), Acceptance bleibt offener DoD-Punkt.
- **Makefile RANGE-Durchreichung:** `trace-check:` reicht `RANGE` per
  `$(if $(RANGE),--range $(RANGE),)` durch (`Makefile:142`) — leer →
  Selbsttest+HEAD, gesetzt → `--range a..b`. Der CI-Aufruf `make trace-check
  RANGE="$RANGE"` (`.github/workflows/ci.yml:61`) trifft den Range-Pfad.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 2 | 1 | 1 | 2 |

## Verdikt

**Blockiert** (HIGH-1 + HIGH-2). Die Topologie des Slices ist im Kern
tragfähig: getrennter PR-/Push-CI ohne Doppellauf mit `release.yml`,
minimale Permissions, SHA-gepinnte Action, fail-fast-Reihenfolge, eine
Skript-Quelle für drei Aufrufer, beidseitig grünes gate-consistency,
schlüssiger Negativ-Selbsttest, korrekt nur-erste-Zeile-greifende
Merge-/Revert-Ausnahme, und `make ci`/`make gates` bleiben unverändert
(ADR-0013 Punkt 5). Es blockiert jedoch **HIGH-1**: der CI-Range-Check —
laut ADR der klon-unabhängige, robuste Backstop — prüft beim ersten Push
eines neuen Branches (`before`=Zero-SHA, der häufigste Mehr-Commit-Fall) nur
einen einzigen Commit; ID-lose Zwischen-Commits rutschen still grün durch.
**HIGH-2** blockiert als Hook/CI-Divergenz gegen die zugesicherte „eine
Wahrheit": der `commit-msg`-Hook greppt die rohe, noch nicht bereinigte
Message-Datei, sodass eine ID nur in einer `#`-Kommentar- oder `-v`-Diff-
Zeile den Hook grün durchwinkt, während git diese Zeilen nach dem Hook
entfernt und der CI-Range-Check denselben Commit ablehnt — ein stilles Grün
auf dem feedforward-Pfad, der vor dem Push warnen soll. **MEDIUM-1**
(Muster nicht deckungsgleich mit `.d-check.yml`, `MR-\d{3}` fehlt; ADR-/
Skript-Aussage der Deckungsgleichheit unzutreffend) ist vor Merge zu klären.
LOW-1 (unverankerte Teilwort-Treffer) und INFO-1/2 sind nicht blockierend.
Die Gate-Bestätigung obliegt der getrennten Verifikation (hier NICHT als
grün angenommen; `make gates`/`make doc-check` wurden im Repo nicht
ausgeführt — Docker/make-only, Reviewer ist kein Verifier).
