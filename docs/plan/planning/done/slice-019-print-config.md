# Slice slice-019: Konfigurations-Gerüst ausgeben (`--print-config`)

**Status:** in-progress.

**Welle:** welle-09-config-geruest (per Roadmap-Fortschreibung; Start
bei Priorisierung durch den Auftraggeber — Adoptions-Reibung in neuen
Repos ohne `.d-check.yml`).

**Bezug:** [`DC-FA-CLI-005`](../../../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben)
(Change Request 0.9.0),
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(das Gerüst dekodiert über den eigenen Parser),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(kein Repo-Zugriff, kein Schreiben),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(statisch → deterministisch).

**Autor:** pt9912. **Datum:** 2026-06-13.

---

## 1. Ziel

Ein neues Repo ohne `.d-check.yml` soll schnell zu einer Config kommen.
`d-check --print-config` gibt ein kommentiertes Startgerüst auf stdout
aus — der Aufrufer leitet um (`> .d-check.yml`). Das Werkzeug schreibt
**nie** selbst (read-only-Kernvertrag bleibt); das Gerüst ist statisch
(nicht aus Repo-Inhalt abgeleitet) und macht zugleich die verfügbaren
Optionen sichtbar.

## 2. Definition of Done

- [x] **Lastenheft-Change-Request** [`DC-FA-CLI-005`](../../../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben)
  (Version 0.9.0, Historie): `--print-config` → statisches Gerüst auf
  stdout, Exit 0, kein Repo-Zugriff; drei AKs + Out-of-Scope (Ableitung
  aus Repo-Inhalt bleibt späterer Modus).
- [x] **Spezifikation** [`DC-FA-CLI-005.a`](../../../../spec/spezifikation.md#dc-fa-cli-005a--konfigurations-gerüst): Kurzschluss-Modus vor jedem
  Repo-Zugriff, eingebettetes Gerüst, dekodiert via eigenem Parser.
- [x] **Implementierung**: Flag `--print-config`, eingebettetes
  `configTemplate`, Kurzschluss in `Run` vor `openRoot`; die drei AKs
  als Tests (gültiges YAML, kein Repo-Zugriff trotz Pfad-Argument,
  Determinismus).
- [x] **End-to-End-Beleg**: Image → `--print-config` > `.d-check.yml` →
  `d-check` läuft mit dem erzeugten Config (2 Dateien, 0 Befunde).
- [x] **Doku**: Option in [`docs/user/operations.md`](../../../../docs/user/operations.md)
  (Tabelle + Umleitungs-Hinweis).
- [x] `make gates` grün; [`CHANGELOG.md`](../../../../CHANGELOG.md);
  Closure-Notiz.

## 3. Plan (vor Code)

| Datei | Art | Begründung |
|---|---|---|
| [`spec/lastenheft.md`](../../../../spec/lastenheft.md) | update | CR [`DC-FA-CLI-005`](../../../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben) (0.9.0) |
| [`spec/spezifikation.md`](../../../../spec/spezifikation.md) | update | [`DC-FA-CLI-005.a`](../../../../spec/spezifikation.md#dc-fa-cli-005a--konfigurations-gerüst) |
| `internal/adapter/driving/cli/config_template.go` | neu | eingebettetes Gerüst |
| `internal/adapter/driving/cli/cli.go` | update | Flag + Kurzschluss |
| `internal/adapter/driving/cli/cli_acceptance_test.go` | update | drei AKs |
| [`docs/user/operations.md`](../../../../docs/user/operations.md) | update | Option dokumentiert |

## 4. Trigger

Priorisierung durch den Auftraggeber (Dialog 2026-06-13): neue Repos
ohne Config brauchen einen Startpunkt; Option statt Schreiben (das
Werkzeug bleibt read-only).

## 5. Closure-Trigger

DoD vollständig inkl. End-to-End-Beleg und grüner Gates.

## 6. Risiken und offene Punkte

- **Drift Gerüst ↔ Schema:** Das eingebettete Template könnte gegenüber
  dem realen `.d-check.yml`-Schema veralten. Sensor dagegen: der
  AK-Test schickt das Gerüst durch den **eigenen** Parser — ein
  inkompatibles Gerüst rötet `make test`.
- **Folge:** slice-020 (`--suggest-config`, Ableitung aus
  Autoritäts-Docs) baut auf diesem Ausgabe-Format auf.

## 7. Closure-Notiz (nach `done/`)

*(folgt mit dem Lifecycle-Übergang nach `done/`.)*

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Spec-/Code-/Doku-Arbeit; Greenfield-Default).
