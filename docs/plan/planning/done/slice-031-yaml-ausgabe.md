# Slice slice-031: YAML-Ausgabeformat (`--yaml`)

**Status:** done (abgeschlossen 2026-06-19).

**Welle:** welle-20-yaml-ausgabe (Trigger: Change Request 0.19.0 akzeptiert;
baut auf der JSON-Ausgabe und slice-029 (`--doctor --json`) auf).

**Bezug:**
[`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)
(Hauptanforderung — YAML als Ausgabeformat),
[`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)
(`--doctor --yaml`-Variante),
[`DC-FA-CLI-003`](../../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes)
(Exit-Codes),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(byte-identisch),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(read-only). ADR: yaml.v3 im report-Adapter (Folge-ADR zu [ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md)).

**Autor:** pt9912. **Datum:** 2026-06-19.

---

## 1. Ziel

`d-check --yaml` gibt die Befunde als YAML auf stdout aus — **dieselbe
Struktur wie `--json`** (`findings`/`summary`/`exitCode`), nur die
Serialisierung unterscheidet sich. Volle Parität: `--doctor --yaml` rendert
die Diagnose maschinenlesbar (mit `reasonText`/`fixCandidate`, analog
`--doctor --json`). `--json`+`--yaml` und `--repair`+`--yaml` sind
Nutzungsfehler (Exit 2). Deterministisch, read-only.

## 2. Definition of Done

- [x] **CLI:** `--yaml` (boolesches Flag); Kombi-Check: `--json`+`--yaml`
  → Exit 2, `--repair`+`--yaml` → Exit 2; `--doctor`+`--yaml` erlaubt.
  render-Switch um `--doctor --yaml` und `--yaml` ergänzt.
- [x] **Renderer:** YAML im report-Adapter über `gopkg.in/yaml.v3`
  ([ADR-0009](../../adr/0009-yaml-im-report-adapter.md)). Format-neutrale
  Output-Structs (json+yaml-Tags) für `report.YAML` und
  `report.DoctorYAML`; eingebettetes `core.Finding` per `yaml:",inline"`
  flach (wie die JSON-Promotion); `fixCandidate: null` explizit.
- [x] **arch-check:** `tools/arch-check.sh` R3 erlaubt yaml.v3 zusätzlich
  in `internal/adapter/driven/report`; Kommentar/Regel-Text nachgezogen.
- [x] **Spezifikation:** [§2 JSON-Ausgabe](../../../../spec/spezifikation.md#json-ausgabe---json)
  um den YAML-Hinweis (gleiche Struktur, nur Serialisierung) ergänzt.
- [x] **Tests:** `--yaml` (parsbar, gleiche Struktur wie JSON, Exit 1);
  `--doctor --yaml` (`reasonText`/`fixCandidate`); `--json --yaml` → Exit 2;
  `--repair --yaml` → Exit 2; Determinismus 10× ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- [x] **Doku:** `docs/user/operations.md` (Options-Tabelle) + Handbuch
  (§4.11 Maschinenlesbare Ausgabe) um `--yaml` ergänzt.
- [x] **[ADR-0009](../../adr/0009-yaml-im-report-adapter.md)** auf `Accepted` (mit diesem Slice); `make gates` grün;
  unabhängiges Review R1; Closure-Notiz.

## 3. Plan (vor Code)

| Datei | Art | Begründung |
|---|---|---|
| `internal/hexagon/core/finding.go` | update | yaml-Tags an `Finding` (Schlüssel-Parität zu JSON) |
| `internal/adapter/driven/report/report.go` | update | format-neutrale Structs (+yaml-Tags), `YAML`/`DoctorYAML`, yaml.v3-Import |
| `internal/adapter/driving/cli/cli.go` | update | `--yaml`-Flag, Kombi-Check, render-Switch |
| `tools/arch-check.sh` | update | R3: yaml.v3 auch in `report` ([ADR-0009](../../adr/0009-yaml-im-report-adapter.md)) |
| `spec/spezifikation.md` | update | §2 YAML-Hinweis |
| `docs/user/operations.md`, `docs/user/benutzerhandbuch.md` | update | `--yaml` dokumentieren |

Lastenheft ist mit Change Request 0.19.0 gesetzt.

## 4. Trigger

Change Request 0.19.0 akzeptiert (Bezug oben); Folge-ADR (Proposed).

## 5. Closure-Trigger

DoD vollständig inkl. grüner Gates, Review R1 und der Folge-ADR auf `Accepted`.

## 6. Risiken und offene Punkte

- **Schlüssel-Parität YAML↔JSON:** yaml.v3 würde Feldnamen ohne Tag
  klein schreiben (`filesChecked` → `fileschecked`); explizite yaml-Tags
  sichern dieselben Schlüssel wie JSON.
- **Embedded-Flattening:** yaml.v3 flacht anonyme Structs **nur** mit
  `yaml:",inline"` ab — sonst würde `core.Finding` verschachtelt; Tag nötig.
- **Determinismus** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  Struct-Marshal ist feldreihenfolge-stabil, keine Map.
- **arch-check-Drift:** die R3-Lockerung muss exakt `report` zulassen,
  nichts sonst (Selbsttest des Gates bleibt scharf).

## 7. Closure-Notiz (nach `done/`)

**Umsetzung:** Vertrag
[`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)
(Lastenheft 0.19.0) + Spezifikation §2 „YAML-Ausgabe". `--yaml` rendert
die Befunde als YAML, strukturgleich zu `--json`; `--doctor --yaml` analog
`--doctor --json`. Formatneutrale Output-Structs (json+yaml-Tags),
gemeinsamer `buildDiagDoc`, je ein Encoder; eingebettetes `core.Finding`
per `yaml:",inline"` flach. yaml.v3 ist über
[ADR-0009](../../adr/0009-yaml-im-report-adapter.md) (Accepted) zusätzlich
im report-Adapter erlaubt; `tools/arch-check.sh` R3 erweitert. `render`
wurde in `render`+`renderStdout` aufgeteilt (gocyclo). `make gates` grün
(Coverage 94,10 %).

**Belege:**

- `--yaml` (parsbar, Schlüssel-Parität zu JSON inkl. camelCase),
  `--doctor --yaml` (`reasonText`/`fixCandidate`, `null` explizit),
  `--json --yaml`/`--repair --yaml` → Exit 2, Determinismus 10×
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- **Am Image verifiziert:** `--yaml` liefert YAML mit gleicher Struktur
  wie JSON; `--json --yaml` → Exit 2.
- read-only ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit));
  `make arch-check` grün mit yaml.v3 im report (R3 erweitert), Kern bleibt
  yaml-frei.

**Lerneintrag:** YAML-Schlüssel-Parität braucht explizite yaml-Tags
(yaml.v3 schreibt sonst Feldnamen klein: `filesChecked` → `fileschecked`)
und `yaml:",inline"` für die flache Promotion des eingebetteten
`core.Finding` (anders als JSON flacht yaml.v3 anonyme Structs nicht
automatisch). Die Hexagon-Importregel (yaml.v3 nur in `configyaml`;
Adapter importieren einander nicht, Regel 5) ließ sich nicht umgehen —
daher die Folge-ADR-Erweiterung statt hand-gerolltem Emitter.

**Review R1** (unabhängiger Reviewer-Subagent, eigener Kontext ohne
DoD-Wissen,
[Report](../../../reviews/2026-06-19-slice-031-yaml-ausgabe.md)): HIGH 0 /
MEDIUM 0 / LOW 2 / INFO 1 — freigegeben. Beide LOW (Doku-Drift: §4.9-
Vergleichstabelle ohne YAML; §4.11-Heading nur `--json`) im selben Stand
geschlossen; INFO-1 (Text-Render-Fehler erzeugt nach dem render-Refactor
eine stderr-Meldung) als bewusste Angleichung akzeptiert.

**Welle:** welle-20-yaml-ausgabe ist damit vollständig (slice-031).

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Code-/Spec-/Doku-Arbeit; Greenfield-Default).
