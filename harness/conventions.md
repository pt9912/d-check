# Harness-Konventionen

## Purpose

Diese Datei deklariert die *repo-lokalen* Strukturregeln dieses Repos
gegenüber der adoptierten Harnesskonvention (Baseline):

- **Adaptionen** ggü. der Baseline (mit Begründung und Auflösungs-Trigger).
- **Zusatzklassen-Deklarationen** für die Sensors-Bindung.
- **Modus-Deklarationen** pro Sub-Area (Greenfield / Brownfield / Hybrid).

Bei Konflikt zwischen dieser Datei und einer kanonischen Quelle gilt
die kanonische Quelle (Source Precedence, siehe
[`README.md`](README.md)).

## Baseline

- **Konvention:** AI-Harness-Kurs
- **Stand:** Template-Set 2026-06
- **Datum der Adoption:** 2026-06-10

## Adoptierte Konventions-Quellen

- **Extern (Lehrmaterial):** `/Development/KI/ai-harness-course`
  (Templates: `lab/templates/`, Konventionen:
  `kurs/de/grundlagen/konventionen.md`) — Pfade außerhalb dieses
  Repos, daher bewusst nicht verlinkt.
- **In-Repo (verkörperte Form):** `AGENTS.md`, `harness/README.md`,
  Verzeichniskonvention `spec/` + `docs/plan/` + `harness/`.

## Adaptions-Block

### MR-000 — Baseline-Aussage

- **Datum:** 2026-06-10
- **Geltungsbereich:** gesamtes Repo
- **Adaption:** keine inhaltlichen Adaptionen ggü. Baseline-Default
  für Verzeichniskonvention, Lifecycle-Regeln, Carveout-Disziplin.
- **Begründung:** Initial-Setzung. Spätere Adaptionen werden als
  `MR-<NNN>` nachgetragen.
- **Auflösungs-Trigger:** permanent.

### MR-001 — Source Precedence mit eigener Spezifikations-Schicht

- **Datum:** 2026-06-10
- **Geltungsbereich:** `harness/README.md` §Source precedence
- **Adaption:** Die Source-Precedence-Tabelle führt
  `spec/spezifikation.md` als eigenen **Rang 2** zwischen Lastenheft
  (Rang 1) und Architektur (Rang 3). Der Kurs-Default setzt zwei
  Spec-Ränge; dieses Repo nutzt drei. Beide Dateien der Ränge 2–3
  entstehen mit slice-002; bis dahin sind sie in den Tabellen als
  „geplant" markiert und nicht verlinkt.
- **Begründung:** Spec-Stratifizierung mit drei Spec-Dateien; die
  ADR-Schärfungs-Regel („ADR darf Spezifikation schärfen, nicht
  Lastenheft") soll strukturell sichtbar sein.
- **Auflösungs-Trigger:** permanent.

### MR-002 — ID-Schema mit Bereichskürzeln ab initialer Fassung

- **Datum:** 2026-06-10
- **Geltungsbereich:** `spec/lastenheft.md`, alle Traceability-Verweise
- **Adaption:** Funktionale Anforderungen verwenden von Beginn an
  Bereichskürzel: `DC-FA-<BEREICH>-<NNN>` (z. B. `DC-FA-LINK-001`)
  statt des zweistelligen Kurs-Defaults `<PREFIX>-FA-<NN>`.
  Nichtfunktionale Anforderungen bleiben beim Kurs-Default
  (`DC-QA-<NN>`).
- **Begründung:** Das Lastenheft konsolidiert zwölf Quell-Tools und hat
  dadurch von Anfang an viele Funktionsbereiche; das Kurs-Beispiel
  (DocSearch) zeigt, dass eine spätere Schema-Migration teurer ist als
  ein Bereichsschema ab Welle 1.
- **Auflösungs-Trigger:** permanent.

### MR-003 — Vendorter Bootstrap-Sensor `tools/verify-doc-refs.sh`

- **Datum:** 2026-06-10
- **Geltungsbereich:** [`tools/verify-doc-refs.sh`](../tools/verify-doc-refs.sh), `make doc-check`
- **Adaption:** Bis `d-check` sich selbst prüfen kann, läuft
  `make doc-check` über ein aus `d-migrate` vendortes Shell-Skript
  (Markdown-Linkziel-Prüfung). Das ist Fremd-Code ohne eigene Spec in
  diesem Repo (Sub-Area in BF, siehe Modus-Tabelle).
- **Begründung:** Ein Doku-Repo ohne Doku-Sensor wäre ein blinder
  Bootstrap; das vendorte Skript ist dependency-frei (bash/awk) und
  deckt den Kern von
  [`DC-FA-LINK-001`](../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
  ab.
- **Auflösungs-Trigger:** slice-004 — `make doc-check` läuft über
  `d-check` selbst (Dogfooding), das Skript wird gelöscht
  ([Slice-Plan](../docs/plan/planning/open/slice-004-anchors-modul-und-dogfooding.md)).

### MR-004 — Gate-Nachweis-Mechanik und `.claude`-Hooks nach b-cad-Vorbild

- **Datum:** 2026-06-10
- **Geltungsbereich:** `tools/harness/`, `.claude/`, `make record-gates`
- **Adaption:** Übernahme der Working-Tree-Hash-Mechanik
  (`record-gates` als letzter `gates`-Prerequisite, Stop-Hook
  vergleicht den Hash) und der `.claude`-Hooks (PreToolUse-Guard,
  Stop-Gate) aus dem Repo `b-cad`.
- **Begründung:** Bewährte Mechanik gegen „Erfolgsmeldung ohne
  Gate-Lauf"; keine Logik-Dopplung zwischen Makefile und Hook.
- **Auflösungs-Trigger:** permanent.

## Zusatzklassen-Deklaration für Sensors-Bindung

— keine — (die kanonischen Klassen ADR, Carveout, Schwelle,
Reproduzierbarkeit sowie Slice-IDs genügen).

## Modus-Deklaration pro Sub-Area

| Sub-Area (Pfad / Modul) | Modus | Begründung | Graduation-Bedingung / Folge-Slice |
|---|---|---|---|
| `*` (Default für gesamtes Repo) | Greenfield | Projekt startet spec-first; Doc führt, Code folgt | n/a (GF) |
| `tools/verify-doc-refs.sh` | Brownfield | vendorter Fremd-Code ohne eigene Spec (siehe [`MR-003`](#mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh)) | slice-004: Ersatz durch `d-check` selbst, Skript wird gelöscht |
| `tools/harness/` | Greenfield | adoptierte Harness-Mechanik, konventionsgetragen über [`MR-004`](#mr-004--gate-nachweis-mechanik-und-claude-hooks-nach-b-cad-vorbild) | n/a (GF) |
