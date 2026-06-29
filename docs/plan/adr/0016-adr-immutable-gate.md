# ADR-0016 — ADR-Immutable-Gate (Accepted-ADRs maschinell unveränderlich)

**Status:** Accepted
**Datum:** 2026-06-21
**Autor:** pt9912
**Bezug:** [`AGENTS.md` §3.5](../../../AGENTS.md#35-adrs-sind-nach-accepted-immutable)
(die durchzusetzende Hard Rule), Mechanik-Präzedenz
[ADR-0013](0013-pr-ci-und-traceability-gate.md) (CI-Range-Check + lokaler
Hook, eine Skript-Quelle), sowie
[`harness/README.md` §Sensors](../../../harness/README.md#sensors-feedback-gates).
Quellen-Grundlage mit nicht-DC-Bezug analog
[ADR-0007](0007-repository-lizenz-mit.md).
**Schärft:** keine Spec-Stelle — Prozess-/Durchsetzungs-ADR; verbindlich für
die maschinelle Durchsetzung von
[`AGENTS.md` §3.5](../../../AGENTS.md#35-adrs-sind-nach-accepted-immutable).

## Kontext

[`AGENTS.md` §3.5](../../../AGENTS.md#35-adrs-sind-nach-accepted-immutable)
verlangt: „Eine ADR mit Status `Accepted` wird nicht inhaltlich
überschrieben. Korrekturen entstehen als neue ADR mit `Supersedes`." Die
Regel ist **dokumentiert, aber nicht erzwungen** (Audit 2026-06-21) — dieselbe
„aspirativ vs. bindend"-Klasse wie die Traceability-Regel vor
[ADR-0013](0013-pr-ci-und-traceability-gate.md) und die Planning-Konsistenz
vor slice-040. Ein stiller Körper-Edit an einer akzeptierten Entscheidung
unterläuft die Provenienz, auf der alle abwärts gerichteten Verweise beruhen.

Beobachtung der gelebten Praxis: superseded ADR-Dateien (z. B.
[ADR-0002](0002-distribution-ghcr-image.md),
[ADR-0004](0004-architektur-pattern-hexagonal.md)) tragen im **Datei-Body**
unverändert `**Status:** Accepted` mit Original-Geschichte; die
Supersede-Annotation lebt im **Index** ([`README.md`](README.md)) und in der
ablösenden ADR. Accepted-ADR-Dateien sind also faktisch bereits voll
immutabel — der Wächter zementiert den Ist-Zustand, statt neue Disziplin zu
fordern.

## Entscheidung

1. **Zwei Quadranten, eine Skript-Quelle** (wie
   [ADR-0013](0013-pr-ci-und-traceability-gate.md)): `tools/adr-immutable-check.sh`
   ist die einzige Logik; `make adr-check` ist ein **dünner Wrapper**, der
   lokale `pre-commit`-Hook ruft dasselbe Skript direkt.
   - **CI-Range** (*feedback*, klon-unabhängiger Backstop): prüft den
     PR-/Push-Range in [`ci.yml`](../../../.github/workflows/ci.yml) über
     dieselbe Range wie `make trace-check`.
   - **Lokaler `pre-commit`-Hook** (*feedforward*): prüft den staged-Diff
     gegen `HEAD`; eingehängt über das bestehende `core.hooksPath` (`make
     hooks`), opt-in pro Klon.

2. **Immutable-Policy (core-compare statt roher Zeilen-Diff).** Geprüft wird
   je **modifizierte** `docs/plan/adr/NNNN-*.md`, deren **BASE**-Version
   `**Status:** Accepted` trägt. „Immutable core" = Dateiinhalt **ohne** den
   `## Geschichte`-Abschnitt (bis zur nächsten `## `-H2 bzw. EOF) **und ohne**
   die `**Status:**`-Zeile.
   - `core(BASE) ≠ core(HEAD)` → FAIL (Körper geändert).
   - HEAD-`**Status:**`-Zeile außerhalb `Accepted…` / `Superseded by ADR-NNNN`
     → FAIL (kein stiller Status-Rückfall).
   - **Gelöschte** oder **umbenannte** Accepted-ADR → FAIL (der ADR-Pfad
     `NNNN-…` ist stabil).
   - Erlaubt bleiben damit: `## Geschichte`-Anhänge und der
     Supersede-/Annotations-Status-Übergang.

3. **Grandfathering automatisch.** Geprüft wird nur der Range-/staged-Diff,
   nicht die Historie. `Proposed`-ADRs sind frei (BASE nicht `Accepted`) —
   inklusive der `Proposed → Accepted`-Reifung einer frischen ADR.

4. **`make adr-check` ist NICHT Teil von `make gates`/`make ci`** — wie
   `make trace-check` ein **Diff-/Commit-Zeit-Bindepunkt**, kein
   Arbeitsbaum-Check. Die CI ruft es als getrennten, billigen Schritt vor dem
   Docker-Build. [`AGENTS.md` §4](../../../AGENTS.md#4-quality-gates) +
   [`harness/README.md` §Sensors](../../../harness/README.md#sensors-feedback-gates)
   dokumentieren das Target (beide Richtungen ⇒ `make gate-consistency`).

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **CI-Range + Hook, core-compare (gewählt)** | erzwingt §3.5 in zwei Quadranten; robust gegen Reformatierung der Geschichte; eine Skript-Quelle | core-Extraktion an die ADR-Struktur (`## Geschichte`, `**Status:**`) gebunden |
| Roher `git diff`-Hunk-Klassifikator | kein Re-Parsen | fragil: jede Verschiebung/Umbruch in Geschichte erzeugt Hunks außerhalb; Whitespace-empfindlich |
| Hash-Baseline pro Accepted-ADR | simpel zu prüfen | Baseline ist mit-editierbar (gleiche Vertrauensgrenze wie die Datei); pflegeintensiv |
| git-history-Vergleich gegen den Accept-Commit | exakt | teuer/komplex; braucht den Accept-Commit pro ADR; kein lokaler staged-Modus |
| Status quo (nur §3.5-Text) | kein Aufwand | Regel bleibt aspirativ |

## Konsequenzen

- Der **Index** ([`README.md`](README.md)) wird **nicht** geprüft
  (Supersede-Annotation lebt dort legitim) — nur `NNNN-*.md`.
- **Branch Protection** (Pflicht-Status-Check) liegt außerhalb des Repos;
  ohne sie ist die CI nur *advisory* — ehrlich benannte Restlücke (wie
  [ADR-0013](0013-pr-ci-und-traceability-gate.md)).
- Die `core`-Extraktion ist an die ADR-Vorlagen-Struktur gekoppelt; ändert
  sich die Geschichte-/Status-Konvention, muss das Skript mitgezogen werden
  (Re-Evaluierungs-Trigger).
- Determinismus/Read-only: das Skript liest nur git-Objekte und schreibt
  nichts ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Geist;
  das Skript ist Harness-Werkzeug, nicht Teil des Produkts).

## Fitness Function

- `make adr-check` läuft **rot** bei einem Körper-Edit oder Status-Rückfall an
  einer `Accepted`-ADR (bzw. Löschung/Umbenennung); **grün** sonst.
  Negativ-Selbsttest bei jedem Lauf (analog
  [`tools/`](../../../tools)-Gate-Skripten).
- CI prüft denselben Range wie `make trace-check`; `make gate-consistency`
  erfasst das neue Target in beide Richtungen.

## Re-Evaluierungs-Trigger

- ADR-Vorlage ändert die Struktur (`## Geschichte`-Überschrift,
  `**Status:**`-Zeilenformat) → `core`-Extraktion anpassen.
- Legitime Klasse von Accepted-Datei-Änderungen taucht auf (über Geschichte +
  Status hinaus) → Policy per Folge-ADR erweitern, nicht still im Skript.
- Wechsel der Forge → Workflow-Portierung, Skript bleibt.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-21 | Proposed — Audit-Befund: §3.5 ohne Wächter; Topologie + Policy via Nutzer-Entscheid (CI-Range + Hook, eigene ADR) |
| 2026-06-21 | Proposed → Accepted (Umsetzung slice-041; Review R1 MERGE-FÄHIG, 1 MEDIUM/2 LOW/1 NIT behoben — `core()`-Status-Strip nur im Header, Selbsttest 6/7, Range-`..`-Guard) |
| 2026-06-29 | Skript-Mechanik (`adr-immutable-check.sh`) teil-superseded durch [ADR-0024](0024-vcs-immutable-gate.md) (slice-053): das Gate `adr-check` läuft auf das d-check-Modul `vcs` um (reine-Go, im Image verteilt). **Policy und Gate unverändert**; das Skript bleibt pfad-stabiler Fallback im Baum. |
