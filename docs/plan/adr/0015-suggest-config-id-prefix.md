# ADR-0015 — `--suggest-config`: Kennungs-Präfix parametrisierbar, kein stiller `DC-`-Default

**Status:** Proposed
**Datum:** 2026-06-21
**Autor:** pt9912
**Bezug:**
[`DC-FA-CLI-006`](../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)
(die geschärfte Anforderung, CR 0.20.0),
[`DC-FA-ID-001`](../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
(das emittierte `ids`-Muster),
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus) (deterministisch),
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(read-only).
**Schärft:** keine Spec-Stelle — Default-/Prozess-ADR; verbindlich für das
Default-Verhalten von `--suggest-config ai-harness[-init]` (Präfix-Quelle,
kein stiller `DC-`).

## Kontext

`--suggest-config ai-harness-init` ist „Voll-Kanon fürs **leere** Repo" —
gedacht zum Bootstrappen *fremder* Harness-Repos, backte aber d-checks
**eigenes** Anforderungs-Präfix `DC` fix ein. Jedes Nicht-d-check-Repo
(a-check = `AC`, b-cad = `BC`, …) bekam das falsche Muster und musste
nacheditieren. Realer Anlass: a-check-Bootstrap (2026-06-20) — das
Init-Template emittierte `DC-(FA-[A-Z]+|QA)` in ein `AC-`-Repo. Das ist die
**stille Fehlsetzung**, die das Regelwerk gerade verbietet („kein stiller
falscher Default"). Nur das **Anforderungs**-Muster trägt ein
projektspezifisches Präfix; `ADR-`/`MR-`/`slice`/Carveout sind die
konventions-festen ai-harness-course-Muster.

## Entscheidung

1. **Präfix parametrisierbar** statt fix `DC`: neues Flag
   **`--id-prefix <PREFIX>`** (z. B. `AC`). Die Quelle (`ai-harness[-init]`)
   bleibt ein reiner Bezeichner; das Präfix ist orthogonale Konfiguration
   (Flag, nicht Token an der Quelle).
2. **Ableitung im Modus `ai-harness`** (repo-bewusst): ohne `--id-prefix`
   wird das Präfix aus `spec/lastenheft.md` abgeleitet — das **eindeutige**
   Projekt-Präfix der FA-/QA-Kennungen. **Mehrere** verschiedene Präfixe ⇒
   Nutzungsfehler (Exit 2), der Mensch gibt `--id-prefix` explizit an
   (deterministisch + ehrlich statt stiller Heuristik).
3. **Platzhalter statt stillem Default:** ist kein Präfix angegeben **und**
   keins ableitbar (insbesondere `ai-harness-init` fürs leere Repo), wird
   ein markierter Platzhalter `<PREFIX>` mit `# TODO`-Hinweis emittiert —
   **kein** `DC-`. Das ist eine **Verhaltensänderung (Breaking)** ggü.
   0.18.1 (`ai-harness-init` ohne Präfix lieferte zuvor `DC-`).
4. **Geltungsbereich:** nur das Anforderungs-Muster ist parametrisiert;
   `ADR-`/`MR-`/`slice`/Carveout-Muster bleiben unverändert konventions-fest.

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **Platzhalter `<PREFIX>` + TODO (gewählt)** | Anti-Footgun; nie ein falscher Default in Fremd-Repos | Breaking ggü. 0.18.1 |
| `DC-` behalten, nur warnen | rückwärtskompatibel | emittiert für Fremd-Repos weiter still das falsche Muster |
| Präfix erzwingen (Fehler ohne Flag) | maximal strikt | bricht jeden heutigen `ai-harness-init`-Aufruf ohne Flag |
| Token an der Quelle (`ai-harness-init:AC`) | kein neues Flag | vermischt Modus und Präfix in einem Bezeichner |

## Konsequenzen

- **Breaking:** `ai-harness-init` ohne `--id-prefix` liefert künftig den
  Platzhalter statt `DC-`. In Lastenheft-Historie (0.20.0) und `CHANGELOG`
  als Breaking-Change ausgewiesen. Der alte Default war für Fremd-Repos —
  den einzigen Adressaten von `ai-harness-init` — ohnehin falsch.
- `ai-harness` (repo-bewusst) gegen d-check selbst bleibt im Ergebnis
  unverändert: das Präfix `DC` wird aus dem eigenen Lastenheft abgeleitet.
- Flag-Wert wird validiert (Großbuchstaben-Präfix); ungültig ⇒ Exit 2.

## Fitness Function

- `make test`: `internal/adapter/driving/cli/cli_acceptance_test.go`
  deckt Happy (`--id-prefix AC`), Boundary (Platzhalter), Negative
  (Konflikt, ungültiger Wert) und die `ai-harness`-Ableitung ab
  (`TestCLI037_*`). Implementierung in
  `internal/hexagon/core/app/suggest.go`.

## Re-Evaluierungs-Trigger

- Bedarf an mehrteiligen/abweichenden Präfix-Gestalten oder einer Heuristik
  im Konfliktfall → neue ADR.
- Automatische Baseline-Versions-Hebung berührt die Vorlage → mitziehen.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-21 | Proposed — slice-037 (Anlass: a-check-Bootstrap 2026-06-20) |
