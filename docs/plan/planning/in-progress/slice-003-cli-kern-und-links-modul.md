# Slice slice-003: CLI-Kern und Modul `links`

**Status:** done.

**Welle:** welle-02-mvp.

**Bezug:** [`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel),
[`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl),
[`DC-FA-CLI-003`](../../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes),
[`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate),
[`DC-FA-SCAN-001`](../../../../spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln),
[`DC-FA-LINK-001`](../../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links),
[`DC-FA-LINK-002`](../../../../spec/lastenheft.md#dc-fa-link-002--symlink-ablehnung),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit); [ADR-0001](../../adr/0001-implementierungssprache.md)–[ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md).

**Autor:** pt9912. **Datum:** 2026-06-10.

---

## 1. Ziel

Ein lauffähiges `d-check`-CLI mit Datei-Scan, Modul `links`
(inkl. Symlink-Ablehnung), Exit-Codes und Text-/JSON-Ausgabe — das
erste implementierte Inkrement.

## 2. Definition of Done

- [x] Akzeptanzkriterien (Happy/Boundary/Negative) der bezogenen `DC-FA-*` als automatisierte Tests umgesetzt und grün.
- [x] Determinismus-Test ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)): wiederholter Lauf, identische Ausgabe-Hashes.
- [x] `make lint`, `make typecheck` (sofern die Toolchain es vorsieht — bei Go deckt `go build`/`govet` das ab, kein eigenes Target) und `make test` existieren (Docker-basiert gemäß [ADR-0001](../../adr/0001-implementierungssprache.md)/[ADR-0002](../../adr/0002-distribution-ghcr-image.md)), tragen ID-Kommentare und sind in `make gates` aggregiert.
- [x] `make arch-check` existiert als Fitness Function zu [ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md) (Pfad-Revision von ADR-0004): Kern ohne I/O-Imports, Netz nur im HTTP-Adapter — strukturelle Durchsetzung von [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).
- [x] Sensors-Tabelle in [`harness/README.md`](../../../../harness/README.md) und Gates-Tabelle in [`AGENTS.md`](../../../../AGENTS.md) §4 aktualisiert — keine behaupteten Targets ohne Existenz.
- [x] `make gates` grün.
- [x] [`CHANGELOG.md`](../../../../CHANGELOG.md) aktualisiert.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Quellbaum gemäß [ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md) (Modul-Pfade, u-boot-Konvention) und [`spec/architecture.md`](../../../../spec/architecture.md) (Rollen/Constraints) | neu | erstes Implementierungs-Inkrement |
| [`docs/plan/adr/0005-modul-layout-hexagon-ordner.md`](../../adr/0005-modul-layout-hexagon-ordner.md) | neu | Pfad-Revision von ADR-0004 (immutable) auf u-boot-Ordnerkonvention — während der Umsetzung entschieden |
| Test-Suite (Unit + Fixture-Repos) | neu | Akzeptanzkriterien sind testbar formuliert |
| [`Makefile`](../../../../Makefile) (`lint`, `typecheck`, `test`, `arch-check`, Aggregation in `gates`) | update | neue Gates entstehen mit dem Code |
| [`harness/README.md`](../../../../harness/README.md), [`AGENTS.md`](../../../../AGENTS.md) | update | Sensors-/Gates-Tabellen nachziehen |
| `Dockerfile` | neu | Toolchain- und Gate-Stages gemäß [ADR-0002](../../adr/0002-distribution-ghcr-image.md) |

## 4. Trigger

welle-01-fundament done (ADRs und Spec-Straten liegen vor).

## 5. Closure-Trigger

DoD vollständig + Commit(s) auf `main` + Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Markdown-Parsing-Kantenfälle (verschachtelte Klammern, Referenz-Links,
  HTML) — Scope bewusst auf das Verhalten der Quell-Tools begrenzen,
  Kantenfälle als Fixtures dokumentieren.
- Der Plan auf Datei-Ebene wird nach slice-002 konkretisiert
  (Architektur legt Modulpfade fest).

## 7. Closure-Notiz (nach `done/`)

**Umsetzung:** Commit `9354d85` (32 Dateien: Quellbaum, Tests,
Dockerfile, Makefile, arch-check, Doku-Nachzug); ADR-0005 entstand
während der Umsetzung.

- **Was hat funktioniert:** Die Spezifikations-Algorithmen
  (Slug/Link-Extraktion/Modul-Auflösung) ließen sich 1:1 in Code und
  Tests übersetzen; das In-Memory-FS aus ADR-0004-Konsequenz machte
  die Kern-Akzeptanztests trivial. `make run` liefert nebenbei den
  ersten Dogfooding-Datenpunkt: d-check prüft das eigene Repo aus dem
  distroless-Image (19 Dateien, 0 Befunde, Exit 0, read-only).
- **Anders als geplant:** (a) Während der Umsetzung kam die
  User-Entscheidung, das Modul-Layout auf die u-boot-Ordnerkonvention
  zu heben (`internal/hexagon/…`, `adapter/{driven,driving}`) — da
  ADR-0004 immutable ist, als [ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md)
  (Teil-Supersede der Pfad-Tabelle). (b) Kein eigenes
  `make typecheck` — Go deckt das über Build/govet ab. (c) Default-
  Modul `anchors` wird bis slice-004 mit stderr-Hinweis übersprungen
  (dokumentierter Interim in `core.ImplementedModules`).
- **Steering-Loop-Lerneintrag:** Das eigene `arch-check`-Gate schlug
  beim ersten Lauf auf `net/url` an — die Fitness-Function-Wildcard
  (`net/*`) war strenger als das ADR-Wording (I/O-Verbot). Lehre für
  Gate-Autoren: Fitness Functions exakt am ADR-Wortlaut kalibrieren;
  ADR-0005 nennt die `net/url`-Ausnahme jetzt explizit.
- **Folge-Slices:** keine neuen; slice-004 (anchors + Dogfooding) ist
  durch diesen Slice entsperrt.

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (neuer Code entsteht spec-first gegen
Lastenheft + Architektur; siehe Kurs Modul 5 §Worked Mini-Example).
