# Realdatenbeleg slice-078 — `ignore-refs` mit Quell-Skopus gegen ai-harness-course

**Datum:** 2026-07-18. **Slice:** slice-078 (geteiltes Referenz-Ventil, DC-FA-REF-001).
**Prüfling:** d-check v-next (`d-check:latest` aus dem feat-Stand, HEAD `d7e779e`),
Kontrast gegen das ausgelieferte v0.48.0-Image. **Gegenstand:** der von der Slice-DoD
geforderte Lauf gegen die realen Quellen des Konsumenten `ai-harness-course` — 0
Findings bei tatsächlich geprüften Verweisen, nicht durch Wegschauen.

## 1. Aufbau

Gefahren gegen eine read-only-Arbeitskopie des `ai-harness-course`-Repos (Kopie ohne
`.git`; das Original wurde nicht zurückgeschrieben). Der Konsument führt heute das
**ganze** Template-Verzeichnis `lab/templates/**` in `scan.ignore` — „by-design
symbolisch", weil die Ziel-Repo-Platzhalter im Quell-Repo nicht auflösen. Damit ist
auch die prüfbare Klasse (Kurs- und Template-interne Verweise) blind.

Zwei Config-Varianten:

- **Baseline:** `lab/templates/**` aus `scan.ignore` entfernt, **kein** neues Ventil
  (nur der bestehende modul-lokale `codepaths.ignore-refs`-Tombstone).
- **Ventil:** zusätzlich der geteilte Top-Level-Block (die Slice-§2-Config):

  ```yaml
  ignore-refs:
    - in: "lab/templates/**"
      refs: ["lab/templates/**", "tools/check_*.py"]
      keep: ["lab/templates/**/*.template.md", "lab/templates/README.md"]
  ```

## 2. Läufe

| Lauf | Image | Ergebnis |
|---|---|---|
| Baseline (un-ignoriert, kein Ventil) | v-next | Exit 1 — **42 Befunde** (37 `target-missing` + 5 `codepath-missing`) |
| Ventil (Slice-§2-Config) | v-next | Exit 0 — **0 Befunde** (158 Dateien geprüft) |

Die **Baseline reproduziert exakt die im CR dokumentierte Messung** (42 = 37 + 5) —
der Aufbau ist damit als treu belegt. Die 42 Befunde sind die Ziel-Repo-Platzhalter:
Links wie `harness/conventions.md`/`spec/lastenheft.md` lösen dateirelativ unter
`lab/templates/` auf und fehlen dort; Inline-Code wie `tools/check_*.py` ebenso. Mit
dem Ventil sind sie alle still (Exit 0).

## 3. Nicht durch Wegschauen (der DoD-Kern)

Damit die 0 nicht „alles ignoriert" bedeutet, zwei Tippfehler in **echte Ziele**
injiziert und der Ventil-Lauf wiederholt:

| Injizierter Tippfehler | Klasse | Ergebnis |
|---|---|---|
| `README.md` → `spec/TIPPFEHLER.template.md` (statt `spezifikation.template.md`) | gekeepter `.template.md`-Verweis | **`target-missing` gemeldet** |
| `project-readme.template.md` → `…/modul-02-TIPPFEHLER.md` | Kurs-Verweis | **`target-missing` gemeldet** |

Beide feuern (Exit 1, genau 2 Befunde), **kein** Platzhalter kommt hinzu. Das belegt:

- **`keep` prüft real** — die per `keep` zurückgeholten `.template.md`-Verweise (die
  sonst-blinde Klasse) werden tatsächlich geprüft, ein Tippfehler wird ERROR. Genau
  das ist der CR-Kern: wer eine Überschrift/Datei umbenennt und das Template vergisst,
  merkt es jetzt.
- **Kurs-Verweise werden geprüft** — sie liegen nicht unter `refs`, werden also nicht
  ignoriert.
- **Muster statt Heuristik** — eine „ignoriere, was nicht auflöst"-Heuristik hätte
  diese Tippfehler ebenfalls verschluckt.

Ein zunächst in einen **Link-Text** (Code-Span als Label) injizierter Tippfehler
erzeugte korrekt **keinen** Befund — Link-Text ist kein Verweis; erst das Ziel nach
`](…)` wird geprüft. (Notiert, weil es die saubere Trennung Text/Ziel bestätigt.)

## 4. Das ausgelieferte Image kann das Ventil nicht ausdrücken

Das v0.48.0-Image gegen die Ventil-Config: `error: field ignore-refs not found` —
das Top-Level-`ignore-refs` (Liste von `{in, refs, keep}`) ist ihm unbekannt, es
lehnt die Config strikt ab (Exit 2). Der modul-lokale `codepaths.ignore-refs`-Alias
allein erreicht die 37 `links`-Befunde ohnehin nicht und kennt keinen Quell-Skopus.
Die Fähigkeit ist damit nachweislich neu.

## 5. Ergebnis

Der Realdatenbeleg ist erbracht: gegen die realen `ai-harness-course`-Quellen meldet
das geteilte Ventil **0 Befunde** — die symbolischen Ziel-Repo-Platzhalter sind still,
die prüfbare Klasse (Kurs- und `.template.md`-Verweise) bleibt scharf (zwei injizierte
Tippfehler beider Klassen werden gefangen). Der bisher pauschale `scan.ignore`-Schnitt
über das ganze Verzeichnis wird durch den ziel-genauen `in`/`refs`/`keep`-Schnitt
ersetzt, ohne eine echte Prüfung zu opfern. Der entsprechende DoD-Punkt ist erfüllt.
