# Slice slice-036: RTM als `d-check --trace`

**Status:** done (Closure 2026-06-21; Review R1+R2, kein ADR nötig).

**Welle:** welle-27-rtm-trace (Trigger: Nutzer-Entscheid 2026-06-20 — RTM als
d-check-Modus statt separatem Skript; ein Prototyp-Skript bewies die
Machbarkeit: 25 `DC-*` → ADR/Slice/Test, alle belegt). *(welle-Nummer von
ursprünglich „welle-25" korrigiert — welle-25 ist welle-25-pr-ci-traceability,
slice-039.)*

**Bezug:**
[`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
(die umgesetzte Anforderung, CR 0.21.0),
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
(Modul `ids` liefert Kennungen + ihre Referenzen),
[`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
(Modul `matrix` liefert Dokumentklassen + Referenzgraph),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(deterministisch),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(read-only).

**Autor:** pt9912. **Datum:** 2026-06-20.

---

## 1. Ziel

Ein read-only-Modus `d-check --trace`, der eine **Requirements Traceability
Matrix** auf stdout ausgibt (**Default Markdown-Tabelle**, optional
maschinenlesbar `--trace --json` / `--trace --yaml`) — **kein Dokument
erzeugt**, immer frisch aus den kanonischen Quellen abgeleitet. Reiht sich in
die Advisory-Modi (`--print-config`/`--suggest-config`/`--doctor`) ein.

## 2. Was die Matrix zeigt

Je Anforderung (`DC-*` im Lastenheft): Titel · referenzierende **ADRs** ·
**Slices** · **Status/Lücke** (Anforderung ohne Slice = Waise). Abgeleitet
aus dem **Referenzgraphen, den `ids` + `matrix` bereits bauen** („welche
Kennung wird von welcher Dokumentklasse referenziert") — der Modus *rendert*
ihn, kaum neue Logik.

## 3. Zu entscheiden (im Slice)

- **Lastenheft-CR:** ein neuer CLI-Vertrag (`DC-FA-CLI`-Familie) für
  `--trace`, wie die übrigen CLI-Modi; Akzeptanzkriterien Happy/Boundary/
  Negative.
- **Test-Spalte (Domänen-Grenze):** der Prototyp grept `*_test.go`
  (Go-Code) nach Kennungen — das verlässt die Markdown-Domäne. Entweder
  weglassen (reine Doku-Traceability Anforderung→ADR→Slice) oder über
  `scan.roots` konfigurierbar scannen. Empfehlung: Kern-RTM doku-only;
  Code-Abdeckung optional.
- **Ausgabeformate (entschieden):** Default Markdown-Tabelle; **optional**
  maschinenlesbar `--trace --json` / `--trace --yaml` (analog
  `--doctor --json`/`--yaml`) — nutzt den vorhandenen report-Adapter
  (slice-031, format-neutrale Output-Structs), kein neuer Serialisierer.

## 4. Definition of Done (vorläufig)

- [x] Lastenheft-CR (0.21.0) + Spezifikation (`spec/spezifikation.md`) für
  `--trace` (neuer CLI-Vertrag; Bezug oben).
- [x] Modus im Paket `app` (`trace.go`, wie diagnose/suggest); **eigene
  Ableitung** (ids/matrix liefern nur Findings, keinen Graphen) mit
  `rules`-Helfern; Markdown-Tabelle auf stdout, deterministisch
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)),
  read-only ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- [x] Lücken-Erkennung (Waise = Anforderung ohne referenzierenden Slice).
- [x] **Maschinenlesbare Ausgabe** `--trace --json` / `--trace --yaml`
  (strukturgleich, über den format-neutralen report-Adapter, slice-031).
- [x] Akzeptanztests (`TestCLI036_Trace_*`); `make gates` grün; **Doku-only**
  entschieden (Code-Test-Spalte verworfen); kein ADR nötig (additiv,
  read-only, keine neuen Import-Kanten).
- [x] Unabhängiges Review R1 (0 HIGH/0 MEDIUM/1 LOW) + R2 (LOW-1/LOW-2
  behoben); Closure.

## 5. Risiken / offene Punkte

- Setzt voraus, dass der `matrix`/`ids`-Graph die Klassen (ADR/Slice/…)
  sauber trennt und die Definitions-/Referenz-Richtung je Klasse exponiert —
  prüfen, was die Module schon liefern.
- **Dogfooding:** d-check tracet seine *eigenen* Anforderungen → der
  `--trace`-Lauf über das eigene Repo ist zugleich Akzeptanztest.
- Abgrenzung: rein additiver read-only-Modus; kein neuer Vertrag an den
  Prüf-Modulen.

## 6. Trigger

Nutzer-Entscheid 2026-06-20: RTM ins d-check (gleiche Doku-Domäne, nutzt
`ids`/`matrix`); arch-check bewusst NICHT (braucht `go list`/Go-Toolchain,
bräche das I/O-freie distroless-Design — eigenes stack-weites Produkt).
Prototyp-Skript als Machbarkeitsbeleg.

## 7. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (CLI-/Doku-Arbeit; Greenfield-Default).

## 8. Closure-Notiz (nach `done/`)

**Umsetzung.** Neuer read-only-Modus `--trace`
([`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)):
leitet die RTM aus den kanonischen Quellen ab (Anforderungen aus
`spec/lastenheft.md`, Referenzen aus `docs/plan/adr/` + `docs/plan/planning/`)
und rendert sie — Markdown-Default, `--trace --json`/`--yaml` über den
format-neutralen report-Adapter (slice-031). **Doku-only** (Auftraggeber-
Entscheidung); Code-Test-Spalte verworfen (bräuchte Go-Toolchain). Eigene
Ableitung im `app`-Paket (ids/matrix liefern nur Findings, keinen Graphen),
präfix-agnostisch, deterministisch sortiert; `runTrace`/`comboError` aus
`cli.Run`/`parseOptions` ausgelagert (Komplexität). **Kein ADR** (additiv,
read-only, keine neuen Import-Kanten; arch-check R1–R6 grün bestätigt).

**Belege.** `make gates` grün (doc-check, lint, test, arch-check, coverage
94,20 %, semgrep 28/0, gate-consistency); read-only
([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)),
deterministisch
([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus),
mehrfach byte-identisch). Dogfooding: `d-check --trace` über das eigene Repo
→ 26 Anforderungen, 1 Waise (die neue Anforderung selbst — jetzt im Bezug
oben nachgezogen, künftig belegt).

**Review R1** (`docs/reviews/2026-06-21-slice-036-rtm-trace.md`):
0 HIGH/0 MEDIUM/1 LOW/3 INFO. **R2**
(`docs/reviews/2026-06-21-slice-036-r2-verifikation.md`): R1 bestätigt;
LOW-1 (Backtick-Heading-Titel) behoben, dessen Fix LOW-2 (Titel-initialer
Code-Span) einführte → mit der wrapper-nur-Strip-Lösung **beide**
geschlossen. INFO (dangling refs verworfen, heading-level-agnostisch,
Waisen → Exit 0, Em-Dash ohne Leerzeichen) won't-fix: spec-konform/bewusst.

**Lerneintrag.** Der RTM-Modus verifiziert sich beim Dogfooding selbst — er
fand seine **eigene** noch unbelegte Anforderung als Waise. Und: ein
Trim-Fix gegen ein Heading-Artefakt kann eines derselben Klasse erzeugen
(LOW-1 → LOW-2); die wrapper-präzise Lösung schlägt die
Zeichenklassen-Erweiterung.
