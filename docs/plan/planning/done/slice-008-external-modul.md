# Slice slice-008: Modul `external` + HTTP-Adapter

**Status:** done.

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

- [x] Akzeptanzkriterien von `DC-FA-EXT-001` als Tests: Status < 400
  ok; ≥ 400 → `external-status`; Timeout → `external-timeout`;
  > 5 Redirects → `external-redirects` (Spezifikation
  §`DC-FA-EXT-001.a`: HEAD mit GET-Fallback bei 405/501, Dedupe pro
  URL, begrenzte Parallelität, Timeout konfigurierbar 1–300 s).
- [x] Opt-in-Garantie getestet: ohne aktiviertes Modul keinerlei
  Netzwerkzugriff; `external` ist nie Teil der Defaults.
- [x] HTTP-Port in `internal/hexagon/port/driven` definiert; Adapter
  in `internal/adapter/driven/httpcheck` (arch-check-Regel R2 greift
  nun positiv); Kern-Tests gegen Port-Fake, Adapter-Tests gegen
  `httptest.Server`.
- [x] `DC-QA-03`-Messmethode automatisiert: Gate-Lauf der
  Default-Module in netzwerkloser Umgebung (`--network none`) gegen
  ein Fixture — als Make-Target in `gates` aggregiert (umgesetzt als
  `--network none` im Dogfooding-Gate `doc-check`: alle Module außer
  `external` gegen das eigene Repo).
- [x] `external` lauffähig (alle fünf Vertragsmodule); der
  Interim-Mechanismus `isImplemented`/`SkippedModules` ist damit
  toter Code und wurde entfernt (statt erweitert — Steering-Loop-
  Eintrag aus slice-006/007); `make gates` grün;
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
  Byte-Identitäts-Garantie ausgenommen (Spezifikation §`DC-QA-02.a`) —
  Tests dürfen nur gegen lokale `httptest`-Server laufen, nie gegen
  echte URLs.
- Parallelität (Default 4) darf die Befund-Sortierung nicht
  beeinflussen (Sammeln → Sortieren bleibt Pflicht).

## 7. Closure-Notiz (nach `done/`)

**Umsetzung:** Commit `f4b603d` (HTTP-Port, `httpcheck`-Adapter,
Kernmodul `external`, QA-03-Netzlos-Gate, Interim-Rückbau).

- **Was hat funktioniert:** Die Port-Abstraktion trug sofort — der
  Kern testet alle Verdikte gegen einen Fake (inkl. Opt-in-Garantie
  über einen `panicChecker`), der Adapter isoliert gegen
  `httptest.Server`; kein Test berührt echte URLs. Das QA-03-Gate
  kostete eine Zeile: `--network none` im bestehenden Dogfooding-Lauf
  — der netzlose Selbstlauf aller vier Default-nahen Module ist exakt
  die Messmethode aus dem Lastenheft.
- **Anders als geplant:** (a) Statt `external` in `isImplemented`
  aufzunehmen, wurde der Mechanismus komplett entfernt (nach diesem
  Slice toter Code — wie in den Closure-Notizen 006/007 vorgemerkt);
  `Run` trägt jetzt den Checker als zweiten Port-Parameter.
  (b) Spez-Fortschreibung aus der Implementierung: Transportfehler
  (DNS/Verbindung) waren in §4 nicht gemappt → `external-status`
  (Status 0); Dedupe-Semantik expliziert (eine Prüfung pro URL,
  Befund an jedem Vorkommen).
- **Steering-Loop-Lerneintrag:** Der „sonst"-Zweig wurde diesmal vor
  der Implementierung gegen Fälle geprüft (Lehre aus slice-006/007) —
  die Transportfehler-Lücke fiel dadurch beim Schreiben des
  Verdikt-Mappings auf, nicht erst im Review. Die Lehre trägt:
  Negativ-Raum-Prüfung der Spezifikation gehört in Schritt 4
  (Plan), nicht in die Review-Runde.
- **Folge-Slices:** keine neuen; slice-009 (coverage-gate,
  gate-consistency, `DC-QA-01`-Benchmark) schließt die Welle.

**Review R1 (nach Closure, Agent-Review mit getrenntem Kontext):**
9 Findings (2 MEDIUM, 5 LOW, 2 INFO); 7 nachgeschärft in einem
Folge-Commit: Fragment-Stripping vor Prüfung/Dedupe (häufigstes
Link-Muster — Anker auf derselben Seite — erzeugte Mehrfach-Requests
gegen dieselbe Ressource); case-insensitiver Schema-Vergleich
(`HTTP://`-Links fielen still zwischen `links` und `external`);
explizit gesetzte 0 in den external-Parametern ist jetzt
Konfigurationsfehler (`*int` im Config-Schema statt stillem Default);
GET-Fallback-Drain auf 64 KB begrenzt; HTTP-Adapter wird nur bei
aktivem Modul verdrahtet (strukturelle Opt-in-Absicherung statt
reiner Kern-Disziplin); Timeout-pro-Request-Semantik spezifiziert;
QA-03-Config-Kopplung als gate-consistency-Auftrag in slice-009
eingetragen. Zwei bewusste Won't-Fix-Design-Notizen: (a)
`HTTPResult`-Exklusivität bleibt Kommentar-Invariante (der Adapter
garantiert sie konstruktiv, ein Verdikt-Enum wäre Port-Churn ohne
akuten Nutzen); (b) Goroutine-pro-URL statt Worker-Pool (bei
realistischen Doku-Repos irrelevant — erst bei zigtausenden externen
Links messbar). Steering-Loop-Beleg: zwei der drei substanziellen
Findings (Fragment, Case) sind Konsistenz-Lücken *zwischen* Modulen —
Lehre: bei Modulen, die dieselbe Eingabe-Klasse (Linkziele)
verarbeiten, die Klassifikations-Helfer teilen statt parallel
implementieren.

Alle berührten Sub-Areas GF (spec-first; siehe Kurs Modul 5 §Worked
Mini-Example).
