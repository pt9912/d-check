# Realdatenbeleg slice-079 — Zitat-Verifikation gegen `ai-harness-init`

**Datum:** 2026-07-18
**Slice:** slice-079 (Zitat-Verifikation)
**Anforderungen:** `DC-FA-CODE-001` (opt-in `codepaths.check-lines`), `DC-FA-CITE-001` (Modul `citations`)
**ADR:** ADR-0045
**Getestetes Artefakt:** `d-check:latest` gebaut aus HEAD `bd257d1` (feat-Commit)
**Konsument:** `ai-harness-init` — **read-only-Kopie** im Scratchpad; das Original
wurde nicht angefasst (mögliche Parallel-Session, nie amenden).

Der Realdatenbeleg ist per ADR-0045 (Fitness-Funktion) und slice-079 §2.2 **nicht
optional**: die eigene Dogfood-Reichweite ist gering (≈7 nicht-`reviews/`-Zitate),
und der Stufe-3-Vertrag wurde nach dem Design-Review R1 auf das
whitespace-normalisierte Teilstring-Modell umgestellt — dieses Modell muss am
**realen** Adopter-Zitat gemessen werden, nicht nur an Fixtures.

## Warum das Adopter-Repo der ehrliche Prüfstein ist

`ai-harness-init` führt seine Baseline **committet vendored**
(`.harness/baseline/v3.1.0/{regelwerk,templates}/`) und erzeugt so einen Korpus
von `datei:zeile`-Zitaten auf einen in-tree, versionierten Fremdbaum — genau die
Drift-Klasse, die den CR ausgelöst hat. Die Zitate sind **inline** (33/33 in
Inline-Code gemessen), re-wrapped und beginnen/enden mitten in der Zeile; ein
zeilenweise-exaktes Modell hätte das eine reale „wortgleich"-Zitat falsch geflaggt
(Design-Review R1, BLOCK). Deshalb hier gegen das Reale gefahren.

## Stufe 3 — Modul `citations` (das Kern-Modell)

Reales Zitat: `docs/plan/planning/done/slice-011-baseline-vendoring.md` zitiert
inline (`„…"`, über die Doc-Zeilen 43–45 re-wrapped) aus
`.harness/baseline/v3.1.0/regelwerk/modul-02-harness-bootstrap.md:173-175`.

| Lauf | Konfiguration | Ergebnis |
|---|---|---|
| **Baseline** | `--enable citations`, **keine** `d-check:cite`-Direktive im Repo | **0 Befunde** (61 Dateien) |
| **Korrekt** | Direktive `<!-- d-check:cite .harness/baseline/v3.1.0/regelwerk/modul-02-harness-bootstrap.md:173-175 -->` unmittelbar vor den Zitat-Absatz | **0 Befunde** — das re-wrapped inline-Zitat ist ein zusammenhängender Teilstring der normalisierten Quell-Spanne |
| **Drift** | ein Wort im Zitat geändert (`committet vendored` → `committet gebündelt`), Quelle unberührt | **1 Befund** `slice-011…:41 …modul-02-harness-bootstrap.md:173-175 citation-mismatch` |

Der **Baseline-Lauf 0/0 belegt den ehrlich ausgewiesenen Caveat** (ADR-0045
Konsequenzen): der Adopter hat **null** `d-check:cite`-Direktiven — das Modell
trägt, aber Stufe 3 prüft ausschließlich Ausgezeichnetes; die produktive Adoption
(Direktiven setzen) ist Adopter-Sache. Der Korrekt-/Drift-Lauf beweist das Modell
am realen Zitat: **korrekt bleibt grün** (kein Fehlalarm durch Re-Wrapping),
**ein gedriftetes Wort wird rot**.

## Stufe 1/2 — `codepaths.check-lines`

Konfiguration der Kopie: `codepaths.roots: [spec, docs, harness]`,
`codepaths.check-lines: true`.

| Lauf | Ergebnis |
|---|---|
| **Reale `datei:zeile`-Refs** (alle im Repo, z. B. `harness/conventions.md:18`, `docs/plan/planning/README.md:26`) | **0 Befunde** — check-lines feuert an keinem realen, in-range Ref (kein Fehlalarm) |
| **Drift** eines realen Refs in einer nicht-exempten Datei (`harness/conventions.md:18` → `:99999`, `slice-012`) | **1 Befund** `slice-012…:25 harness/conventions.md:99999-99999 citation-out-of-range` |

Nicht durch Wegschauen: der zunächst in `docs/reviews/**` injizierte Drift feuerte
**nicht** — korrekt, denn diese Domäne ist per `codepaths.exempt-paths` datei-weit
ausgenommen; erst der Drift in einer nicht-exempten `done/`-Datei erzeugte den
Befund.

## Feature ist nachweislich neu

`ghcr.io/pt9912/d-check:v0.49.0 --enable citations` bricht mit
`error: unbekanntes Modul "citations"` ab — das Modul existiert vor slice-079
nicht. `codepaths.check-lines` war in v0.49.0 kein Config-Schlüssel (Default aus,
byte-identisch).

## Fazit

Beide Stufen sind am realen Adopter-Korpus belegt: korrekt = grün, Drift = rot,
Baseline byte-identisch/leer. Das nach R1 redesignte Teilstring-Modell trägt das
reale, re-wrapped inline-Zitat, das das alte zeilenweise-exakte Modell falsch
geflaggt hätte. Der Beleg deckt die per ADR-0045 geforderte Fitness-Funktion.
