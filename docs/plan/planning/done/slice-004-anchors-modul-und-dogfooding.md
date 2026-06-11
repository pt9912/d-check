# Slice slice-004: Modul `anchors` und Dogfooding-Umstellung

**Status:** done.

**Welle:** welle-02-mvp.

**Bezug:** [`DC-FA-ANCH-001`](../../../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors),
[`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl),
[`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools);
[`MR-003`](../../../../harness/conventions.md#mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh) (vendorter Bootstrap-Sensor).

**Autor:** pt9912. **Datum:** 2026-06-10.

---

## 1. Ziel

Das Modul `anchors` ist implementiert und `make doc-check` läuft über
`d-check` selbst (Dogfooding) — der vendorte Bootstrap-Sensor
`tools/verify-doc-refs.sh` wird gelöscht.

## 2. Definition of Done

- [x] Akzeptanzkriterien von [`DC-FA-ANCH-001`](../../../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors) als Tests umgesetzt und grün (inkl. Duplikat-Suffix-Boundary).
- [x] `make doc-check` ruft `d-check` (Default-Module `links` + `anchors`) auf dem eigenen Repo auf und ist grün.
- [x] `tools/verify-doc-refs.sh` gelöscht; [`harness/conventions.md`](../../../../harness/conventions.md): Auflösung von [`MR-003`](../../../../harness/conventions.md#mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh) als neuer MR-Eintrag nachgetragen, Modus-Tabelle bereinigt (BF-Sub-Area entfällt).
- [x] Vergleichslauf dokumentiert: `d-check` findet auf diesem Repo mindestens die Befunde des Alt-Skripts, keine neuen False-Positives (erster Datenpunkt für [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)).
- [x] `make gates` grün.
- [x] [`CHANGELOG.md`](../../../../CHANGELOG.md) aktualisiert.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Modul `anchors` (Pfad gemäß [ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md)) | neu | [`DC-FA-ANCH-001`](../../../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors) |
| Test-Suite (Anker-Fixtures, Slug-Kantenfälle) | neu | testbare Akzeptanzkriterien |
| [`Makefile`](../../../../Makefile) (`doc-check` auf `d-check` umstellen) | update | Dogfooding, Ablösung Bootstrap-Sensor |
| `tools/verify-doc-refs.sh` | löschen | [`MR-003`](../../../../harness/conventions.md#mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh)-Auflösungs-Trigger |
| [`harness/conventions.md`](../../../../harness/conventions.md) | update | MR-Eintrag zur Auflösung, Modus-Tabelle |
| [`harness/README.md`](../../../../harness/README.md) | update | Sensors-Tabelle: `doc-check`-Bindung aktualisieren |

## 4. Trigger

slice-003 done (CLI-Kern und Modul `links` vorhanden).

## 5. Closure-Trigger

DoD vollständig + Commit(s) auf `main` + Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Das eigene Repo nutzt Anker-Links auf gedashte Heading-Slugs
  (z. B. `#dc-fa-link-001--…`) — der Slug-Algorithmus muss vor der
  Umstellung gegen GitHub-Verhalten verifiziert sein, sonst bricht
  `doc-check` mit False-Positives.
- Falls `d-check` zum Umstellungszeitpunkt nur als lokales Build
  vorliegt (noch kein GHCR-Image, welle-04), läuft `doc-check`
  übergangsweise gegen das lokal gebaute Image — im Makefile
  dokumentieren.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung:** Commit `3fd2067` (Modul `anchors`, Dogfooding-Schnitt,
`MR-007`, Spez-Fortschreibung); Vorarbeiten in `2a2bd3f`/`52df5c3`
(Black-Box-Review-Fixes am CLI).

- **Was hat funktioniert:** Der normierte Slug-Algorithmus aus der
  Spezifikation ließ sich direkt als Tabelle von Slug-Tests gegen die
  realen Heading-Formen dieses Repos absichern; der erste
  Dogfooding-Lauf validierte auf Anhieb alle ~100 Anker-Links des
  eigenen Repos (23 Dateien, 0 Befunde) — und bestätigte damit
  maschinell das Doppel-Bindestrich-Urteil aus der Spec-Review-Runde.
  Der Vergleichslauf gegen das Alt-Skript (0/0, bei strikt größerer
  Abdeckung) ist der erste Datenpunkt für `DC-QA-04`.
- **Anders als geplant:** `scan.roots` brauchte den Wert `"."`
  (gesamte Wurzel), damit das Dogfooding die Abdeckung des
  Alt-Skripts (harness/, Top-Level) erreicht — Spez fortgeschrieben.
  Slug-Schritt 1 wurde auf GitHub-Verhalten präzisiert (literale
  Unterstriche bleiben). Der `MR-003`-Geltungsbereich-Link wurde nach
  der Skript-Löschung auf einen Code-Span umgestellt (Form-Fix,
  dokumentiert in `MR-007`).
- **Steering-Loop-Lerneintrag:** Der dogfoodete Sensor schlug schon
  während des Slice produktiv zu (Lifecycle-Move von slice-004 brach
  den `MR-003`-Link — doc-check rot, sofort gefixt): Der Sensor
  erzwingt jetzt die Link-Pflege, die vorher Review-Disziplin war.
  Lehre: Lifecycle-Moves immer mit `make doc-check` vor dem Commit.
- **Folge-Slices:** keine neuen; welle-02-mvp ist mit diesem Slice
  abgeschlossen (M2 erreicht), welle-03-regelmodule wird aktiv —
  deren Slices (ids, matrix, external, coverage-gate,
  gate-consistency) werden beim Wellen-Start geschnitten.

## 8. Sub-Area-Modus-Begründung

Berührt eine Sub-Area in BF; daher voller Begründungsblock.

### Sub-Area: `tools/verify-doc-refs.sh` (vendorter Bootstrap-Sensor)

- **Modus:** BF
- **Konventionen-Dichte:** hoch — Vendoring in [`harness/conventions.md`](../../../../harness/conventions.md)
  [`MR-003`](../../../../harness/conventions.md#mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh) deklariert, Auflösungs-Trigger ist dieser Slice.
- **Phase-Reife:** Bootstrap-Phase — der Sensor existierte vor jeder
  Implementierung und hat keine eigene Spec in diesem Repo.
- **Evidenz-/Diskrepanz-Risiko:** niedrig — ~100 Zeilen Shell; das
  Soll-Verhalten ist durch
  [`DC-FA-LINK-001`](../../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
  abgedeckt; der
  Vergleichslauf (DoD) macht Diskrepanzen sichtbar.
- **Reconciliation-Aufwand:** klein — das Skript wird ersatzlos
  gelöscht, sobald `make doc-check` auf `d-check` umgestellt ist
  (Graduation = Löschung, kein Doc-Nachzug nötig).

### Übrige berührte Sub-Areas

GF (neuer Modul-Code spec-first; siehe Kurs Modul 5 §Worked
Mini-Example).
