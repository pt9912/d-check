# Slice slice-039: PR-/Push-CI und Traceability-Gate

**Status:** in-progress (seit 2026-06-21; Artefakte gebaut, Hooks
installiert, `make gates` grün; unabhängiges Review R1 +
[ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md)→`Accepted`
ausstehend).

**Welle:** welle-25-pr-ci-traceability (Trigger: Audit 2026-06-21 —
`harness/README.md` §Traceability rules ohne Wächter + kein
PR-/Push-CI-Workflow;
[ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md) Proposed).

**Bezug:**
[ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md) (Entscheidung —
CI-Topologie + Traceability-Gate),
[ADR-0011](../../adr/0011-digest-pins-build-gate-images.md) (digest-gepinnte
Engine ⇒ lokal = PR = Release),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Determinismus bei jeder Integration prüfen),
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
(`make ci` ⊃ `make image-test`).

**Autor:** pt9912. **Datum:** 2026-06-21.

---

## 1. Ziel

Die zwei heute **aspirativen** Harness-Regeln in **bindende** Kontrollen
überführen ([ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md)):
(1) ein PR-/Push-CI-Workflow, der `make ci` auf jede Integration zieht
(*pre-integration* statt nur Release); (2) ein Traceability-Gate, das eine
`DC-*`/`ADR-*`/`slice-*`-ID in Commit-Messages erzwingt — als lokaler
`commit-msg`-Hook (*feedforward*) **und** CI-Range-Check (*feedback*) über
**eine** gemeinsame Skript-Wahrheit.

## 2. Ausgangslage (Ist)

- `.github/workflows/` enthält nur `release.yml` (Tag-Push `v*` →
  `make ci` → GHCR). Normale `pull_request`/`push` haben **keinen**
  maschinellen Boden im Repo.
- `harness/README.md` §Traceability rules verlangt `DC-*`/`ADR-*` in
  Commits — es gibt **keinen** `commit-msg`-Hook und **keinen** CI-Check;
  6 der 15 jüngsten Commits tragen nur eine `slice-NNN`-Planungs-ID.
- `make ci` (gates + image-test) existiert und ist die Release-Engine
  (digest-gepinnt,
  [ADR-0011](../../adr/0011-digest-pins-build-gate-images.md)) —
  wiederverwendbar, nichts Neues zu bauen.
- Der Stop-Hook (`.claude/hooks/stop-require-gates.sh`) gibt cleane Klone
  ohne Gate-State frei (dokumentierte Restlücke) — CI ist das Netz.

## 3. Definition of Done (vorläufig)

Artefakte (geplant):

```
.github/workflows/ci.yml     # pull_request + push (ohne Tags) -> make ci + make trace-check
.githooks/commit-msg         # ruft die gemeinsame Skript-Quelle für eine Message
tools/trace-check.sh         # eine Wahrheit: Message- und Range-Modus
Makefile: trace-check, hooks # duenner Target-Wrapper + core.hooksPath-Installer
```

- [x] `ci.yml`: Trigger `pull_request` (gegen Default-Branch) + `push`
  auf Branches **ohne** Tags (keine Doppelläufe mit `release.yml`); Job
  ruft `make ci` **und** `make trace-check` (Range) als getrennte Schritte.
- [x] Gemeinsame Skript-Quelle mit Message- und Range-Modus;
  `make trace-check` ist ein **dünner Wrapper** ohne eigene Logik (sonst
  divergiert der Hook, der das Skript direkt ruft).
- [x] `make hooks` setzt `core.hooksPath` auf das versionierte
  Hook-Verzeichnis; der `commit-msg`-Hook prüft die Message gegen die
  Kennungs-Muster aus
  [ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md) (deckungsgleich
  mit `.d-check.yml`, inkl. `slice-`); Merge- und `Revert`-Commits
  ausgenommen.
- [x] Selbsttest des Skripts (grün bei vorhandener ID, rot ohne) — analog
  zum Negativ-Selbsttest in `tools/gate-consistency.sh`.
- [x] **Doku-Sync (Sync-Trigger aus
  [ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md)):**
  `harness/README.md` §Traceability rules und die Commit-Disziplin in
  `AGENTS.md` um `slice-*` als zulässige Planungs-ID ergänzt;
  `make trace-check`/`make hooks` in `harness/README.md` §Sensors **und**
  `AGENTS.md` §4 dokumentiert (`make gate-consistency` prüft beide
  Richtungen).
- [x] `make gates` grün (inkl. `gate-consistency` über die neuen Targets);
  `make trace-check` grün auf HEAD.
- [ ] Unabhängiges Review R1; Closure-Notiz;
  [ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md) auf `Accepted`
  (Acceptance-Trigger nach Umsetzung).

## 4. Risiken und offene Punkte

- **Branch Protection** (Pflicht-Status-Checks) liegt extern in den
  GitHub-Settings — nicht aus dem Klon auditierbar. Ohne sie ist die
  PR-CI nur *advisory*. Restlücke, **kein** Carveout (kein gebrochenes
  Gate) — als Betriebshinweis in Closure/`docs/user/` vermerken.
- **make-only vs. Host-Toolchain:** der git-Hook läuft auf dem Host (wie
  `.claude/hooks/*`). Bewusst ein **dünner Stolperdraht** (bash + git),
  keine zweite Toolchain; die robuste Kontrolle bleibt der CI-Range-Check.
- **Commit-Range in CI:** shallow clone → `fetch-depth` setzen oder den
  Range aus dem Event-Payload lesen (push: `before..after`; PR:
  `base..head`).
- **Fork-PRs:** `make ci` baut Docker; Secrets/Permissions in fork-PRs
  prüfen — gates/trace-check sind read-only und brauchen keine
  Push-Secrets.
- **CI-Kosten:** `make ci` baut Images (Minuten) pro PR; Layer-Cache via
  digest-gepinnte Stages
  ([ADR-0011](../../adr/0011-digest-pins-build-gate-images.md)); ein
  optionaler gates-only-Fast-Path bleibt offen.

## 5. Trigger

Audit 2026-06-21 (Traceability-Regel ohne Wächter; nur Tag-Release-CI) →
[ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md) (Proposed,
Review R1 eingearbeitet). Die Umsetzung dieses Slices schaltet
[ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md) auf `Accepted`.

## 6. Sub-Area-Modus-Begründung

Berührte Sub-Areas — alle **GF** (Greenfield-Default; Tooling/Build/Doku):

- **CI-/Build-Infrastruktur** (`.github/workflows/`): GF — neuer Workflow
  neben `release.yml`; die Entscheidung
  ([ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md)) führt, Code
  folgt.
- **Durchsetzungsschicht** (Hook + Skript + `make`-Targets): GF — neue
  Kontrolle aus der Entscheidung; kein Bestandscode-Inventar nötig.
- **Doku-Sync** (`harness/README.md`, `AGENTS.md`): GF — der
  Regel-Wortlaut wird der Entscheidung nachgezogen (Sync-Trigger).

Kein BF/Hybrid (kein Diskrepanz-Cluster, kein „Code vor Doku"); kein
Reconciliation-Aufwand.
