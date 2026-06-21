# ADR-0013 — PR-/Push-CI und Traceability-Gate (DC-*/ADR-* in Commits)

**Status:** Accepted
**Datum:** 2026-06-21
**Autor:** pt9912
**Bezug:** [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)
(make-ci-Vertrag bei jeder Integration),
[`DC-FA-DIST-001`](../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
(`make ci` ⊃ `make image-test`),
[ADR-0011](0011-digest-pins-build-gate-images.md) (digest-gepinnte Engine
⇒ lokal = PR = Release), [ADR-0002](0002-distribution-ghcr-image.md)
(GHCR-Release-Pipeline), sowie
[`harness/README.md` §Traceability rules](../../../harness/README.md#traceability-rules)
(die durchzusetzende Regel; Quellen-Grundlage analog
[ADR-0007](0007-repository-lizenz-mit.md) mit nicht-DC-Bezug).
**Schärft:** keine Spec-Stelle — Prozess-/Durchsetzungs-ADR; verbindlich
für die CI-Topologie des Repos und die maschinelle Durchsetzung von
[`harness/README.md` §Traceability rules](../../../harness/README.md#traceability-rules).

## Kontext

Zwei Regeln des Harness sind heute **aspirativ statt bindend** (Audit
2026-06-21):

1. [`harness/README.md` §Traceability rules](../../../harness/README.md#traceability-rules)
   verlangt „PRs/Commits **müssen** mindestens eine `DC-*`- oder
   `ADR-*`-ID nennen". Es gibt aber **keinen** Wächter dafür — kein
   `commit-msg`-Hook, keinen CI-Check. Die Disziplin wird gelebt (zuletzt
   6/15 Commits nennen eine `DC-*`/`ADR-*`-ID, weitere 6/15 nur eine
   `slice-NNN`-Planungs-ID), aber nichts erzwingt sie. Das ist die
   Lücke „aspirativ vs. bindend" der Durchsetzungsschicht (Kurs-Modul:
   Durchsetzungsschicht; Traceability-Constraint als *computational
   feedforward* per Commit-Hook).
2. Im Repo existiert nur ein **Tag-Release-Workflow**
   ([`release.yml`](../../../.github/workflows/release.yml): Tag-Push
   `v*` → `make ci` → GHCR-Push mit Digest). Normale Änderungen
   (`pull_request`, `push` auf Branches) haben **keinen maschinellen
   Boden** im Repo. `make gates`/`make ci` laufen heute nur lokal; der
   `Stop`-Hook
   ([`stop-require-gates.sh`](../../../.claude/hooks/stop-require-gates.sh))
   ist Claude-spezifisch und gibt frische, cleane Klone ohne lokalen
   Gate-State **frei** — dokumentierte Restlücke, für die der Kurs
   explizit **CI als Netz** vorsieht.

`make ci` (gates + image-test) existiert bereits und ist die Engine der
Release-Pipeline — die CI-Topologie muss also nichts Neues *bauen*,
sondern denselben Vertrag auf einen früheren Lebenszyklus-Punkt
(*pre-integration*) ziehen.

## Entscheidung

1. **PR-/Push-CI-Workflow** neben dem bestehenden Release-Workflow.
   Trigger: `pull_request` (gegen den Default-Branch) und `push` auf
   Branches (**ohne** Tags — Tag-Läufe bleiben
   [`release.yml`](../../../.github/workflows/release.yml) vorbehalten,
   keine Doppelläufe). Der Job ruft **`make ci`** — identische,
   digest-gepinnte Engine ([ADR-0011](0011-digest-pins-build-gate-images.md))
   wie Release. Damit gilt lokal = PR = Release; die Stop-Hook-Restlücke
   (frischer Klon) ist über CI geschlossen.

2. **Traceability-Gate in zwei Quadranten** (durchgesetzt, nicht nur
   behauptet):
   - **Lokaler `commit-msg`-Hook** (*computational feedforward*): prüft
     die Commit-Message beim Commit gegen die Kennungs-Muster. Ein
     versioniertes Hook-Verzeichnis wird per `make hooks`
     (`git config core.hooksPath`) eingehängt — opt-in pro Klon, daher
     ist (1) der klon-unabhängige Backstop.
   - **CI-Check** (*computational feedback*): prüft den Commit-Range des
     PR/Push auf mindestens eine ID — robust, klon-unabhängig.

3. **Eine Wahrheit, keine Logik-Dopplung.** Hook und CI rufen dasselbe
   Prüf-Skript unter `tools/` (analog zur gemeinsamen
   Working-Tree-Hash-Quelle für `make gates` und Stop-Hook). Das
   `make`-Target ist ein **dünner Wrapper** über dieses Skript — keine
   eigene Logik, sonst divergiert der Hook (der das Skript direkt ruft)
   vom Target. `make gate-consistency` erfasst das Target erst, wenn es
   **sowohl** in [`harness/README.md`](../../../harness/README.md)
   §Sensors **als auch** in [`AGENTS.md`](../../../AGENTS.md) §4 steht
   (beide Richtungen):

   ```
   tools/trace-check.sh        # gemeinsame Prüf-Quelle (Message | Range)
   .githooks/commit-msg        # ruft tools/trace-check.sh <message-file>
   .github/workflows/ci.yml    # ruft `make ci` + `make trace-check` (Range)
   Makefile: trace-check, hooks # Target-Wrapper + Installer
   ```

   „Eine Wahrheit" heißt auch **uniforme Message-Bereinigung**: Hook
   (`--message`) und CI (`--range`) wenden dieselbe `#`-Kommentar-/scissors-
   Bereinigung an, bevor sie auf eine ID prüfen. Sonst divergieren sie an
   Messages mit ID nur auf einer `#`-Zeile, je nach git-Cleanup-Modus
   (`-m`=whitespace behält `#`, Editor=strip entfernt `#`) — eine ID gehört
   auf eine Inhalts-Zeile, nicht in einen Kommentar (Review R1 HIGH-2 /
   R2 HIGH-A, beide dieselbe Divergenz-Klasse).

4. **Kennungs-Muster = eine Quelle.** Das Gate nutzt deckungsgleiche
   Muster mit `.d-check.yml` (`ids`), damit „gültige ID" nur **einmal**
   definiert ist:

   ```
   ADR-\d{4}              (ADR-IDs)
   MR-\d{3}               (Konventions-Adaptionen)
   DC-(FA-[A-Z]+|QA)-\d+  (Anforderungs-IDs)
   slice-\d+              (Planungs-/Slice-IDs)
   ```

   Die drei `.d-check.yml`-`ids`-Muster (ADR/MR/DC) **plus** `slice-*`:
   `MR-*` lässt Konventions-Commits (`harness/conventions.md`) zu,
   `slice-*` die Planning-Lifecycle-Commits (6 der 15 jüngsten Commits
   tragen nur eine Slice-ID) — ein reines `DC-*`/`ADR-*`-Gate würde beide
   blocken. **Sync-Trigger (Folge-Slice):** der Wortlaut in
   [`harness/README.md`](../../../harness/README.md) §Traceability rules
   und die Commit-Disziplin in [`AGENTS.md`](../../../AGENTS.md) werden um
   `MR-*`/`slice-*` als zulässige IDs nachgezogen, damit Regel und Gate
   deckungsgleich sind. Ausnahmen (kein ID-Zwang): Merge- und
   `Revert`-Commits.

5. **`make trace-check` ist NICHT Teil von `make gates` oder `make ci`.**
   Gates/`make ci` prüfen den **Arbeitsbaum-Inhalt** (Handoff-Bindepunkt)
   und bleiben unverändert (gates bzw. gates + image-test); Traceability
   prüft **Commit-Messages** (Commit-Zeit-Bindepunkt). Der **CI-Workflow**
   ruft daher `make ci` *und* `make trace-check` als getrennte Schritte
   (so im Layout in Punkt 3); der lokale Hook prüft nur die Message.

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **PR-CI + Hook, eine Skript-Quelle (gewählt)** | Hard Rule in zwei Quadranten (feedforward + feedback); klon-unabhängiger Backstop; lokal lauffähig | zwei Bindepunkte zu pflegen; CI-Laufzeit pro PR |
| Nur CI-Check | minimal; klon-unabhängig | kein lokales Frühwarn-Feedback vor dem Push (Schleife erst im PR) |
| Nur lokaler `commit-msg`-Hook | leichtgewichtig | nicht aus dem Klon erzwingbar (opt-in); kein Boden für Fremd-Beiträge |
| Traceability in `make gates` falten | ein Target weniger | falscher Bindepunkt — Gates prüfen Tree, nicht Commit-Messages |
| Status quo (nur Tag-Release) | kein Aufwand | Regel bleibt aspirativ; normale Änderungen ohne maschinellen Boden |

## Konsequenzen

- Diese ADR **baut nichts** (Status Proposed). Workflow, Hook, Skript und
  `make`-Targets entstehen in einem Folge-Slice, dort mit DC-Bindung,
  DoD und `make gates`-Nachweis.
- `make trace-check`/`make hooks` werden nach Umsetzung in
  [`harness/README.md`](../../../harness/README.md) §Sensors **und** in
  [`AGENTS.md`](../../../AGENTS.md) §4 dokumentiert — `make gate-consistency`
  prüft beide Richtungen (Target↔Doku), sonst rotes `make gates`.
- **Branch Protection** (Pflicht-Status-Checks auf dem Default-Branch)
  liegt außerhalb des Repos und ist nicht aus dem Klon auditierbar.
  Ohne sie ist die PR-CI nur *advisory* — ehrlich benannte Restlücke
  (Empfehlung: extern konfigurieren, im Folge-Slice als Betriebshinweis
  vermerken).
- **CI-Kosten:** `make ci` baut Docker-Images (Minuten) pro PR.
  Mitigation: Layer-Cache über digest-gepinnte Stages
  ([ADR-0011](0011-digest-pins-build-gate-images.md)); ein optionaler
  schnellerer `gates`-only-Pfad bleibt dem Folge-Slice überlassen.
- Determinismus-Vertrag
  ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus))
  wird erstmals bei **jeder** Integration maschinell geprüft, nicht erst
  beim Release.

## Fitness Function

- `make trace-check` läuft **rot** bei einem Commit/Range ohne
  `DC-*`/`ADR-*`-ID; **grün** sonst (Selbsttest analog
  [`tools/`](../../../tools)-Gate-Skripten).
- PR-/Push-CI führt `make ci` aus; ein roter Gate blockiert den Merge
  (sobald Branch Protection gesetzt ist).
- `make gate-consistency` erfasst die neuen Targets in beide Richtungen.

## Re-Evaluierungs-Trigger

- Wechsel der Forge (GitHub → andere) → Workflow-Portierung, Skript bleibt.
- Commit-Range in CI unzuverlässig (shallow clone) → `fetch-depth`
  anpassen oder Range aus dem Event-Payload lesen.
- Häufige legitime ID-lose Commit-Klassen tauchen auf → Ausnahmeliste in
  Punkt 4 erweitern (per Folge-ADR, nicht still im Skript).

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-21 | Proposed — Audit-Befund: Traceability-Regel ohne Wächter, kein PR-/Push-CI |
| 2026-06-21 | Review (docs/reviews/2026-06-21-adr-0013-pr-ci-traceability.md) eingearbeitet: Index-Zeile ergänzt; `slice-*` ins Muster (Planning-Commits); `make ci`/`make trace-check` entkoppelt; gate-consistency-Bewachung in `harness/README.md` + `AGENTS.md`; ID-Quote korrigiert |
| 2026-06-21 | Proposed → Accepted (Umsetzung slice-039; Review R1: 2 HIGH/1 MEDIUM behoben, R2: HIGH-A behoben + `MR-*` ergänzt — beide Runden adversarial verifiziert) |
