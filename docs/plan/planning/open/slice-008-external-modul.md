# Slice slice-008: Modul `external` + HTTP-Adapter

**Status:** open.

**Welle:** welle-03-regelmodule.

**Bezug:** [`DC-FA-EXT-001`](../../../../spec/lastenheft.md#dc-fa-ext-001--externe-links-modul-external-opt-in),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(Messmethode als automatisierter Test);
[ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md)
(Pfad `internal/adapter/driven/httpcheck`).

**Autor:** pt9912. **Datum:** 2026-06-10.

---

## 1. Ziel

Das opt-in-Regelmodul `external` ist implementiert (HTTP-Port im
Hexagon, `httpcheck`-Adapter), und die `DC-QA-03`-Messmethode läuft
als automatisierter Test.

## 2. Definition of Done

- [ ] Akzeptanzkriterien von `DC-FA-EXT-001` als Tests: Status < 400
  ok; ≥ 400 → `external-status`; Timeout → `external-timeout`;
  > 5 Redirects → `external-redirects` (Spezifikation
  §DC-FA-EXT-001.a: HEAD mit GET-Fallback bei 405/501, Dedupe pro
  URL, begrenzte Parallelität, Timeout konfigurierbar 1–300 s).
- [ ] Opt-in-Garantie getestet: ohne aktiviertes Modul keinerlei
  Netzwerkzugriff; `external` ist nie Teil der Defaults.
- [ ] HTTP-Port in `internal/hexagon/port/driven` definiert; Adapter
  in `internal/adapter/driven/httpcheck` (arch-check-Regel R2 greift
  nun positiv); Kern-Tests gegen Port-Fake, Adapter-Tests gegen
  `httptest.Server`.
- [ ] `DC-QA-03`-Messmethode automatisiert: Gate-Lauf der
  Default-Module in netzwerkloser Umgebung (`--network none`) gegen
  ein Fixture — als Make-Target in `gates` aggregiert.
- [ ] `external` in `isImplemented` (alle fünf Vertragsmodule
  lauffähig); `make gates` grün;
  [`CHANGELOG.md`](../../../../CHANGELOG.md); Closure-Notiz.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/port/driven/http.go` | neu | HTTP-Port (Architektur §2) |
| `internal/adapter/driven/httpcheck/` (+ Tests) | neu | HEAD/GET, Timeout, Redirect-Limit |
| `internal/hexagon/core/external.go` (+ Tests) | neu | URL-Sammlung, Dedupe, Befund-Mapping |
| [`Makefile`](../../../../Makefile)/[`Dockerfile`](../../../../Dockerfile) | update | QA-03-Netzlos-Gate |

## 4. Trigger

Sofort — unabhängig von slice-006/007 (welle-03 aktiv).

## 5. Closure-Trigger

DoD vollständig + Commit(s) auf `main` + Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Netz-Nichtdeterminismus: `external` ist von der
  Byte-Identitäts-Garantie ausgenommen (Spezifikation §DC-QA-02.a) —
  Tests dürfen nur gegen lokale `httptest`-Server laufen, nie gegen
  echte URLs.
- Parallelität (Default 4) darf die Befund-Sortierung nicht
  beeinflussen (Sammeln → Sortieren bleibt Pflicht).

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (spec-first; siehe Kurs Modul 5 §Worked
Mini-Example).
