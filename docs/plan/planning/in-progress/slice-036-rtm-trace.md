# Slice slice-036: RTM als `d-check --trace`

**Status:** in-progress (seit 2026-06-21; Code/Tests/Spec-CR/Spezifikation/
Doku fertig, `make gates` grün; Review R1 + Closure ausstehend).

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
- [ ] Unabhängiges Review R1; Closure.

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
