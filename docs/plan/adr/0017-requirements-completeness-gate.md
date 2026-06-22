# ADR-0017 — Requirements-Completeness-Gate (Waisen als Closure-Invariante)

**Status:** Accepted
**Datum:** 2026-06-22
**Autor:** pt9912
**Bezug:** [`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
(RTM-Modus `--trace`, dessen Waisen-Markierung das Gate liest),
Mechanik-Präzedenz [ADR-0013](0013-pr-ci-und-traceability-gate.md) und
[ADR-0016](0016-adr-immutable-gate.md) (eigener Bindepunkt außerhalb von
`make gates`/`ci`, eine Skript-Wahrheit + Negativ-Selbsttest), sowie
[`harness/README.md` §Sensors](../../../harness/README.md#sensors-feedback-gates).
**Schärft:** keine Spec-Stelle — Prozess-/Durchsetzungs-ADR; macht die
advisory Waisen-Markierung von
[`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
zu einer erzwungenen Closure-Invariante, ohne deren Exit-0-Vertrag zu ändern.

## Kontext

Die RTM (`--trace`, slice-036) markiert je Anforderung Waisen (kein
referenzierender Slice), ist aber **advisory**: Exit 0 auch bei Waisen
(spec-fixierte Akzeptanzkriterien). Damit ist „jede Anforderung hat einen
umsetzenden Slice" — die Abdeckungs-Hälfte von „Arbeit abgeschlossen" —
dokumentiert und sichtbar, aber **nicht erzwungen**. Dieselbe „aspirativ vs.
bindend"-Klasse wie Traceability vor [ADR-0013](0013-pr-ci-und-traceability-gate.md)
und Planning-Konsistenz vor slice-040.

Nutzer-Beobachtung 2026-06-22: „damit kann man prüfen, ob die Arbeit
abgeschlossen ist." Die offene Frage ist der **Bindepunkt**: ein per-Commit-Gate
widerspräche dem Greenfield-Prinzip „Doc führt, Code folgt" — eine frisch ins
Lastenheft geschriebene Anforderung ist legitim **transient Waise**, bis ihr
Slice landet.

## Entscheidung

Wir führen ein **Closure-Meta-Gate `make completeness-check`** ein:

1. **Bindepunkt Closure, nicht per Commit.** Das Target ist an `make fullbuild`
   (volle Closure vor Welle-Merge/Release) gebunden und **NICHT** Teil von
   `make gates`/`make ci`. So bleibt der per-Commit-Inner-Loop GF-konform
   (transiente Waisen erlaubt), während der Abschluss-Checkpoint „0 Waisen"
   erzwingt — Schwester-Logik zu [ADR-0013](0013-pr-ci-und-traceability-gate.md)/[ADR-0016](0016-adr-immutable-gate.md),
   die bewusst Bindepunkte außerhalb von `gates` wählen.
2. **Mechanik: dünner Wrapper über `--trace --json`, eine Skript-Wahrheit.**
   Das Wächter-Skript `completeness-check.sh` (unter `tools/`) ruft das
   Runtime-Image `d-check --trace --json`
   (`--network none`, ro-Mount) und liest das Feld `orphans` (int) des
   RTM-JSON: `orphans > 0` ⇒ FAIL mit der Liste der Waisen-IDs
   (`requirements[].orphan == true`). `make completeness-check` ist der dünne
   Wrapper. **Parsing rein mit bash/grep** (kein `jq`/`python` — keine
   Host-Binary-Abhängigkeit, konsistent mit den bestehenden `tools/`-Gate-
   Skripten); ein fehlendes oder nicht-numerisches `orphans` ⇒ FAIL, nie
   stilles „0".
3. **`--trace` bleibt advisory (Exit 0).** Die Spec-Akzeptanz von
   [`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
   (Exit 0 bei Waisen) wird **nicht** angetastet — die Durchsetzung lebt im
   Gate-Wrapper, nicht im Produkt-Exit-Code. Kein Change Request am Lastenheft.
4. **Negativ-Selbsttest bei jedem Lauf (fail-closed, beide Richtungen).** Auf
   der Parser-Kernlogik (synthetisches JSON, ohne Container): `orphans: 1` muss
   feuern, `orphans: 0` nicht — **und** die Stilles-Grün-Vektoren müssen FAIL
   ergeben: leeres stdout (Image-Exit ≠ 0 / leerer Lauf), JSON ohne
   `orphans`-Feld, nicht-numerischer Wert, kaputtes JSON. Analog den
   [`tools/`](../../../tools)-Gate-Skripten; ohne diese Richtungen wäre der
   Wächter silent-green-anfällig (dieselbe Falle wie slice-040/041).
5. **Doku-Kopplung.** [`AGENTS.md` §4](../../../AGENTS.md#4-quality-gates) +
   [`harness/README.md` §Sensors](../../../harness/README.md#sensors-feedback-gates)
   dokumentieren das Target (beide Richtungen ⇒ `make gate-consistency`); in der
   Gate-Taxonomie als **Meta-/Governance-Gate** mit einem **dritten Bindepunkt
   „Closure"** (neben „in `gates`" und „Commit-/Diff-Bindepunkt") geführt. Die
   Bindepunkt-Klassen-Zuordnung ist eine Form-Setzung, die `gate-consistency`
   nicht prüft (es misst nur Target-Existenz) — daher in der Taxonomie explizit
   benannt.

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **Wrapper über `--trace --json`, Closure-gebunden (gewählt)** | erzwingt Abdeckung am richtigen Bindepunkt; eine Skript-Wahrheit; nutzt das released Feature; GF-Inner-Loop bleibt frei | greift nur am Closure-Checkpoint, nicht per Commit (bewusst) |
| `--trace` selbst auf Exit ≠ 0 bei Waisen | kein Wrapper nötig | Spec-Change an der advisory-Akzeptanz von [`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix); bräche jeden GF-Zwischenstand mit transienter Waise |
| In `make gates`/`ci` aufnehmen | per Commit erzwungen | GF-feindlich: eine neue Anforderung vor ihrem Slice ist nicht mehr commit-bar |
| Status quo (nur advisory `--trace`) | kein Aufwand | Abdeckung bleibt unerzwungen (aspirativ) |

## Konsequenzen

- **Semantik-Grenze (ehrlich benannt):** „0 Waisen" = jede Anforderung von ≥1
  Slice **beansprucht** — nicht, dass der Slice *done* ist (Slice-Status
  out-of-scope in
  [`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix))
  noch dass die Implementierung die Anforderung *erfüllt* (das prüfen
  Tests/Verification). Das Gate ist die Abdeckungs-, nicht die
  Fertigstellungs-Garantie.
- Eine bewusst vorgezogene Zukunfts-Anforderung (Lastenheft vor Slice) failt
  `make fullbuild`/Release — gewollt (kein Release mit unimplementierter
  Vertrags-Anforderung), als Konsequenz benannt.
- **JSON-Schema-Kopplung:** ändert sich das RTM-JSON
  ([`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)),
  muss das Skript mit (Re-Eval-Trigger; fail-closed bei Parse-Fehler).
- Determinismus/Read-only: das Skript ruft nur den read-only/netzlosen
  `--trace`-Lauf und schreibt nichts
  ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Geist;
  Harness-Werkzeug, nicht Produkt).
- `make fullbuild` braucht das Runtime-Image (ohnehin via `ci`/`image-test`
  gebaut) — keine neue Toolchain.

## Fitness Function

- `make completeness-check` läuft **rot** bei ≥1 Requirements-Waise, **grün**
  bei 0; Negativ-Selbsttest bei jedem Lauf (analog
  [`tools/`](../../../tools)-Gate-Skripten).
- `make fullbuild` schlägt fehl, wenn der Closure-Stand Waisen trägt;
  `make gate-consistency` erfasst das neue Target in beide Richtungen.

## Re-Evaluierungs-Trigger

- RTM-JSON-Schema ändert sich
  ([`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix))
  → Parser anpassen.
- Bedarf, den Slice-*Status* (done vs. in-progress) in die Abdeckung
  einzubeziehen → Folge-ADR (Semantik-Erweiterung), nicht still im Skript.
- Greenfield-Phase vorbei / Bedarf nach per-Commit-Erzwingung → Bindepunkt per
  Folge-ADR verschieben.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-22 | Proposed — Nutzer-Wunsch „prüfen, ob die Arbeit abgeschlossen ist"; Bindepunkt-Policy (Closure/`fullbuild`, nicht per Commit) + Wrapper über `--trace --json` (advisory bleibt) via Nutzer-Entscheid; Umsetzung slice-042 |
| 2026-06-22 | Proposed (Nachbesserung) — unabhängiges Review R1 zu slice-042 (NACHBESSERN, 1 MEDIUM): Negativ-Selbsttest auf Stilles-Grün-Vektoren erweitert (F-1), bash/grep-Parsing festgeschrieben (F-2), Closure-Bindepunkt als dritte Taxonomie-Klasse benannt (F-3) |
| 2026-06-22 | Proposed → Accepted (Umsetzung slice-042: `make completeness-check` + `tools/completeness-check.sh`, an `fullbuild` gehängt; `make gates` grün, adversariale Waisen-Fixture → Exit 1; R1-Findings F-1/F-2/F-3 eingearbeitet) |
