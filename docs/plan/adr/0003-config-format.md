# ADR-0003 — Config-Parsing: striktes YAML via gopkg.in/yaml.v3

**Status:** Accepted
**Datum:** 2026-06-10
**Autor:** pt9912
**Bezug:** [`DC-FA-CONF-001`](../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei),
[ADR-0001](0001-implementierungssprache.md)
**Schärft:** `spec/spezifikation.md` (entsteht mit slice-002:
`.d-check.yml`-Schema und Fehlertexte) — nicht das Lastenheft.

## Kontext

Das Lastenheft fixiert die Konfigurationsdatei als `.d-check.yml`
(YAML) und fordert vollständige Validierung mit Zeilenangabe im
Fehlerfall sowie Exit-Code 2 ohne stillschweigende Defaults
([`DC-FA-CONF-001`](../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)).
Zu entscheiden sind Parser-Bibliothek und Validierungsstrategie —
Gos Standard-Library enthält kein YAML (Konsequenz aus
[ADR-0001](0001-implementierungssprache.md)).

## Entscheidung

1. **Parser:** `gopkg.in/yaml.v3`, per `go.sum`/Version gepinnt — die
   im Go-Ökosystem etablierteste YAML-Bibliothek.
2. **Striktes Decoding:** `KnownFields(true)` — unbekannte Schlüssel
   sind Konfigurationsfehler, keine stillen No-Ops (Tippfehler wie
   `modul:` statt `modules:` dürfen nicht leise Defaults aktivieren).
3. **Zweistufige Validierung, fail-closed:**
   - Stufe 1 (Syntax): YAML-Parse-Fehler → Exit 2 mit Zeilen-/
     Spaltenangabe (yaml.v3-Node-API liefert Positionsdaten — der
     Grund gegen JSON-Umweg-Bibliotheken).
   - Stufe 2 (Semantik): Kennungs-Regexe kompilieren, Modulnamen gegen
     Whitelist, deklarierte Scan-Wurzeln auf Existenz, Matrix-Klassen
     auf Konsistenz → jeder Fehler Exit 2.
4. **Keine Config-Vererbung,** eine Datei in der Repo-Wurzel — wie im
   Lastenheft abgegrenzt.

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **gopkg.in/yaml.v3 (gewählt)** | etabliert, striktes Decoding, Zeilen-/Spalteninfo via Node-API | externes Modul (gepinnt) |
| goccy/go-yaml | schneller, gute Fehlertexte | geringere Verbreitung im eigenen Stack, Mehrwert hier irrelevant (Config ist winzig) |
| sigs.k8s.io/yaml | JSON-Schema-Tooling nutzbar | YAML→JSON-Umweg verliert Zeilennummern — verletzt das Negative-Kriterium von [`DC-FA-CONF-001`](../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei) |
| JSON/TOML statt YAML | Stdlib (JSON) bzw. simpleres Format | verletzt den Lastenheft-Vertrag `.d-check.yml`; Formatwechsel wäre Change Request, kein ADR |

## Konsequenzen

- Einzige geplante Fremd-Dependency des Kerns ist damit deklariert;
  weitere Dependencies brauchen erneut eine ADR-Abwägung.
- Das maschinenlesbare Schema von `.d-check.yml` wird in
  `spec/spezifikation.md` (slice-002) normiert; dieser ADR legt nur
  Parser + Strategie fest.
- Der Config-Adapter kapselt yaml.v3 vollständig
  ([ADR-0004](0004-architektur-pattern-hexagonal.md)) — der Kern sieht
  nur validierte Strukturen.

## Fitness Function

- Tests gegen die Akzeptanzkriterien von [`DC-FA-CONF-001`](../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei) (Happy/
  Boundary/Negative), inkl. „unbekannter Schlüssel → Exit 2 mit
  Zeilenangabe".
- `make arch-check` (ab slice-003): yaml.v3 darf nur im Config-Adapter
  importiert werden.

## Re-Evaluierungs-Trigger

- yaml.v3 wird unmaintained oder bekommt eine nicht behebbare
  Sicherheitslücke.
- Das Lastenheft ändert per Change Request das Konfigurationsformat.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-10 | Proposed → Accepted (slice-001) |
