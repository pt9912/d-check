# ADR-0006 — Lint-Profil: SOLID-nah nach u-boot-Vorbild

**Status:** Accepted
**Datum:** 2026-06-10
**Autor:** pt9912
**Bezug:** [ADR-0001](0001-implementierungssprache.md) (golangci-lint
als Lint-Werkzeug), [ADR-0005](0005-modul-layout-hexagon-ordner.md)
(Architektur-Fitness-Function)
**Schärft:** keine Spec-Stelle — Prozess-/Qualitäts-ADR
(Kurs-Modul 13); verbindlich für `.golangci.yml` und `make lint`.

## Kontext

Seit slice-003 lief golangci-lint nur mit den 5 Default-Lintern —
ein dokumentierter Interim („Erweiterung Richtung SOLID-Profil ist
eine eigene ADR-Entscheidung"). Das Ökosystem-Vorbild
[`u-boot`](https://github.com/pt9912/u-boot) fährt per dortigem
ADR-0003 ein SOLID-nahes Profil (5 Default- + 24 weitere Linter,
kalibrierte Schwellen, ausschließlich zentrale, Why-kommentierte
Ausnahmen). Inline-Suppressions sind in d-check ohnehin verboten
(`AGENTS.md` §3.2).

## Entscheidung

Übernahme des u-boot-Profils mit identischer Kalibrierung
(cyclop 15, dupl 150, funlen 100/100·60, gocognit 20, gocyclo 15,
interfacebloat 10, maintidx 20, nestif 5; revive mit explizitem
Default-Regelblock + `unused-receiver`; gomodguard mit
Anti-Modul-Liste; forbidigo gegen direktes `fmt.Print*`) — mit zwei
bewussten Abweichungen:

1. **Kein `depguard`:** Die Architektur-Import-Regeln laufen in
   d-check bereits als eigene Fitness Function
   (`tools/arch-check.sh`, [ADR-0005](0005-modul-layout-hexagon-ordner.md)
   §Fitness Function). Eine depguard-Kopie wäre eine zweite Quelle
   derselben Regeln und müsste bei jeder Layout-Änderung synchron
   gehalten werden.
2. **`gochecknoglobals` ohne Code-Carveouts außerhalb `cmd/`:**
   Statt Pfad-Ausnahmen werden die bisherigen Paket-Globals
   (Modul-/Wurzel-/Skip-Tabellen) zu Funktionen refactort — das
   Profil soll den Code formen, nicht umgekehrt.

Zentrale Ausnahmen (alle in `.golangci.yml` mit `Why:`):
Komplexitäts-/Kontext-Linter für `_test.go` (table-driven Tests);
`testpackage` für `internal/hexagon/core/` (White-Box-Tests der
unexportierten Analyse-Primitiven; Black-Box-Abdeckung über die
CLI-Akzeptanztests).

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **u-boot-Profil ohne depguard (gewählt)** | Ökosystem-Parität, erprobte Kalibrierung; eine Quelle für Architektur-Regeln | zwei Abweichungen müssen dokumentiert bleiben (dieses ADR) |
| Nur 5 Default-Linter (bisher) | minimaler Pflegeaufwand | lässt SOLID-Verstöße (Komplexität, Globals, Interface-Bloat) unbeanstandet |
| u-boot-Profil inkl. depguard | 1:1-Parität | dupliziert die ADR-0005-Fitness-Function; Drift-Risiko zwischen zwei Regelquellen |
| `golangci-lint` „enable-all" | maximale Abdeckung | unkalibriert, hohe False-Positive-Last, Profil nicht begründbar |

## Konsequenzen

- `make lint` wird strenger; Komplexitäts-Verstöße erzwingen
  Refactorings statt Ausnahmen.
- Neue Ausnahmen brauchen einen `Why:`-Kommentar in `.golangci.yml`;
  eine Lockerung der Kalibrierung ist eine neue ADR
  (`AGENTS.md` §3.6).
- Die Anti-Modul-Liste (gomodguard) hält die Dependency-Fläche bei
  der einen gepinnten YAML-Bibliothek
  ([ADR-0003](0003-config-format.md)).

## Fitness Function

`make lint` (Dockerfile-Stage `lint`, golangci-lint mit diesem
Profil); rot = Profil verletzt.

## Re-Evaluierungs-Trigger

- u-boot hebt sein Lint-Profil substanziell (Parität neu bewerten).
- Ein Linter des Profils produziert systematische False Positives im
  d-check-Code (dann Ausnahme mit Why oder neue ADR).

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-10 | Proposed → Accepted (slice-005) |
