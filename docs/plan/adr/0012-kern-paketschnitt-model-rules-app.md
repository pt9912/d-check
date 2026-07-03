# ADR-0012 — Kern-Paketschnitt: `model` / `rules` / `app`

**Status:** Accepted
**Datum:** 2026-06-20
**Autor:** pt9912
**Bezug:** [ADR-0004](0004-architektur-pattern-hexagonal.md) („Hexagon
light" — diese ADR schärft die Kern-Binnenstruktur, der hexagonale
Außenschnitt bleibt),
[ADR-0005](0005-modul-layout-hexagon-ordner.md) (Modul-Layout/arch-check —
wird auf drei Kern-Pakete erweitert),
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(I/O-Freiheit des Kerns bleibt Vertrag)
**Schärft:** [`spec/architecture.md`](../../../spec/architecture.md) §Kern
(die „Kern"-Zeile wird in drei Ringe aufgelöst).

## Kontext

Der Kern `internal/hexagon/core` ist ein **einziges** Go-Paket mit
**5.212 Nicht-Test-Zeilen** in 17 Dateien — Befund-Modell, Config-Typen,
acht Regelmodule, Markdown-/Pfad-Primitive und die Orchestrierung
(Discover, `--doctor`/`--repair`/`--suggest-config`) liegen ungetrennt
nebeneinander. Bei [ADR-0004](0004-architektur-pattern-hexagonal.md)
(„light", ein Kern-Paket) war das angemessen; mit dem Wachstum kostet die
fehlende Binnenstruktur Übersicht und lässt die Schichtrichtung
(Daten ← Regeln ← Orchestrierung) **unerzwungen**. Eine Messung zeigt: das
Modell (`finding.go`/`config.go`, 228 Z.) ist ein sauberes Blatt, und die
Aufruf-Richtung ist bereits faktisch einbahnig — bis auf **einen** geteilten
Glob-Helfer (`ignored`, `scan.go`), den sowohl `ids` (Regel) als auch die
Discovery (Orchestrierung) nutzt.

## Entscheidung

Der Kern wird in **drei Pakete** unter `internal/hexagon/core/` geschnitten,
mit **strikt einbahniger** Importrichtung:

1. **`model`** — `finding.go`, `config.go`: Befund-Modell, Config-Typen,
   reine Wert-Funktionen (`SortFindings`, `EffectiveModules`, `Effective*`).
   Importiert **nichts** aus dem Kern (nur Standardbibliothek) — innerster
   Ring.
2. **`rules`** — die acht Regelmodule (`links`/`anchors`/`ids`/`matrix`/
   `codepaths`/`spans`/`hostpaths`/`external`), die Primitive `markdown.go`,
   `paths.go` (mit dem dorthin verschobenen Glob-Helfer `ignored`) **sowie
   die Prüf-Orchestrierung `run.go` und die Discovery `scan.go`** — die
   I/O-freie Prüf-Engine. Importiert `model` und die Ports (`port/driven`),
   **nicht** `app`.
3. **`app`** — die Anwendungs-Modi auf den Befunden: `diagnose.go`
   (`--doctor`), `repair.go` (`--repair`), `suggest.go`
   (`--suggest-config`). Importiert `model`, `rules` und die Ports.

**Zur Lage von `run`/`scan`** (Abweichung vom Erst-Plan, der sie `app`
zuordnete): die Modul-Tests (White-box) koppeln den Orchestrator `Run` mit
Modul-Interna; in `rules` bleiben sie ohne Interna-Export testbar. „Engine"
(Module + ihre Ausführung) vs. „Modi" ist die vom Code getragene Naht.

**Importregel (neu, maschinell erzwungen):**
`app → rules → model`, `app → model`, `rules/app → port/driven`;
`port/driven → model`/∅. **Verboten:** `model → {rules, app, port}`,
`rules → app`, `port → {rules, app}`. Der Go-Compiler verweigert Zyklen;
zusätzlich kodiert `tools/arch-check.sh` die Richtung explizit (Fitness
Function, [ADR-0005](0005-modul-layout-hexagon-ordner.md)).

Die **I/O-Freiheit** des gesamten Kerns
([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit))
bleibt: alle drei Pakete reine Berechnung + Port-Interfaces, kein
Dateisystem-/Netz-/Prozess-Zugriff.

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **3 Pakete `model`/`rules`/`app` (gewählt)** | trennt Daten/Regeln/Orchestrierung; Richtung maschinell erzwingbar; adressiert den 5k-Z.-Brocken | größter Refactor (17 + ~13 Test-Dateien umpaketieren, arch-check neu) |
| Status quo (ein Kern-Paket) | null Aufwand | 5.212 Z. ungetrennt; Schichtrichtung unerzwungen |
| 2 Pakete `model`/`service` | kleiner | isoliert nur die 228-Z.-Domain; der Brocken (Regeln+Orchestrierung) bleibt ungeteilt |
| Aufteilung je Regelmodul (ein Paket pro Modul) | feinste Granularität | viele Mini-Pakete, gemeinsame Primitive (markdown/paths) erzwingen erneut ein Sammel-Paket |

## Konsequenzen

- `internal/hexagon/core/{model,rules,app}/`; `package core` entfällt.
  Adapter/CLI importieren künftig `…/core/model` (Reporter/Config-Typen)
  bzw. `…/core/app` (CLI ruft `app.Run`).
- `ignored` (und ggf. weitere von `rules` geteilte reine Helfer) wandern aus
  `scan.go` nach `paths.go` (`rules`), damit die Richtung `rules → app`
  vermieden wird.
- White-box-Tests folgen ihrem Paket; der geteilte In-Memory-FS-Helfer
  (`memfs_test.go`) wird ein exportierter Test-Helfer in einem von allen
  drei Test-Paketen importierbaren Ort (oder `port/driven`-nahes Testpaket).
- `tools/arch-check.sh` und [`spec/architecture.md`](../../../spec/architecture.md)
  §Kern werden auf die drei Ringe nachgezogen.
- Rein strukturell: kein Verhaltens-/Befund-Delta (Determinismus
  [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus) unberührt);
  Beleg über byte-identischen Befundsatz vor/nach.

## Fitness Function

`make arch-check` erzwingt die Importrichtung der drei Kern-Pakete (R1–R5
nachgezogen + neue Richtungsregel); `make gates` grün; Befundsatz auf einem
Fixture vor/nach byte-identisch.

## Re-Evaluierungs-Trigger

- `rules` wächst weiter unhandlich → Aufteilung je Modul-Familie erwägen.
- Ein künftiger Bedarf bricht die Schichtrichtung → ADR neu bewerten.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-20 | Proposed (slice-035) |
| 2026-06-20 | Accepted — Umsetzung slice-035: `model`/`rules`/`app` geschnitten, `run`/`scan` in `rules` (Engine), arch-check R6 erzwingt die Richtung; `make gates` grün, alle Tests bestehen (kein Befund-Delta) |
| 2026-07-03 | R6-Durchsetzungs-Mechanik teil-superseded durch [ADR-0029](0029-arch-check-via-a-check.md) (slice-058): `go list`-Skript → a-check-`edges` (`model`←`rules`←`app` als Richtungs-Allowlist); die R6-Policy dieser ADR bleibt unverändert |
