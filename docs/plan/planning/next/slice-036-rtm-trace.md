# Slice slice-036: RTM als `d-check --trace`

**Status:** next (geplant — noch nicht in Arbeit).

**Welle:** welle-25-rtm-trace (Trigger: Nutzer-Entscheid 2026-06-20 — RTM als
d-check-Modus statt separatem Skript; ein Prototyp-Skript bewies die
Machbarkeit: 25 `DC-*` → ADR/Slice/Test, alle belegt).

**Bezug:**
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
Matrix** als **Markdown-Tabelle auf stdout** ausgibt — **kein Dokument
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
- **Kombinierbarkeit:** `--trace --json`/`--yaml` für maschinenlesbare
  Trace-Ausgabe (analog `--doctor --json`)?

## 4. Definition of Done (vorläufig)

- [ ] Lastenheft-CR + Spezifikation für `--trace` (neuer CLI-Vertrag).
- [ ] Modus im Paket `app` (wie diagnose/suggest), nutzt den `ids`/`matrix`-
  Referenzgraphen aus `rules`; Markdown-Tabelle auf stdout, deterministisch
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)),
  read-only ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- [ ] Lücken-Erkennung (Waisen: Anforderung ohne Slice/ADR).
- [ ] Akzeptanztests; `make gates` grün; unabhängiges Review R1; ADR nach
  Bedarf.

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
