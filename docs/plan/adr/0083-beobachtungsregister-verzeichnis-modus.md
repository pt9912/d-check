# ADR-0083: Das Beobachtungs-Register bekommt einen zweiten, additiven Verzeichnis-Modus

**Status:** Accepted

**Datum:** 2026-09-03

**Autor:** pt9912

**Bezug:** [ADR-0082](0082-uebergangswaechter-reviews-observations.md) (die
vierte `planning`-Fähigkeit, die dieser ADR erweitert), welle-88 <!-- d-check:status-provenance -->, slice-193 <!-- d-check:status-provenance -->
(Baseline-Pin auf `v6.0.0`, adoptiert dessen Regelwerk-Delta ohne die
Register-Neugestaltung selbst umzusetzen).

**Schärft:**
[`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(neue, fünfte Fähigkeit — additiv, keine Änderung der vierten).

**Regeln:** Baseline-Regelwerk
[`modul-06-roadmap.md` §Das Beobachtungs-Register](../../../.harness/baseline/v6.0.0/regelwerk/modul-06-roadmap.md#das-beobachtungs-register-modul-6).

---

## Kontext

Baseline `v6.0.0` (Kurs-Wellen 115–116) gestaltet das Beobachtungs-Register
komplett neu: aus einer Tabellen-Datei (`observations.md`, Kennung
`BEO-<NNN>`, ein gepflegter Zähler) wird eine Verzeichnis-Ablage
(`observations/BEO-<KUERZEL>/<slug>/{observation.md,state.md,evidence/}`,
Kennung = Pfad, Zähler abgeleitet aus der Zahl der `evidence/`-Dateien).

d-checks eigenes Produkt trägt diese Fähigkeit bereits — die **vierte**
`planning`-Fähigkeit
([`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
ADR-0082, ausgeliefert seit `v0.72.0`): `internal/hexagon/core/rules/planning_observations.go`
prüft, ob eine zitierte Kennung eine Zeile in einer **Tabellen-Datei**
(`observations.register`) hat. Das ist genau die **alte** Kanon-Form.

**Zwei Zwänge treffen sich hier:**

1. **Der eigene Bestand besteht bis slice-195 <!-- d-check:status-provenance --> im alten Format fort.**
   `docs/plan/planning/observations.md` bleibt eine Tabelle, bis die
   Datenmigration (slice-195 <!-- d-check:status-provenance -->, blockiert von diesem Slice) sie umzieht. Ein
   Bruch der Tabellen-Form JETZT machte `make gates` sofort rot, bevor die
   Migration überhaupt beginnt.
2. **Andere Adopter-Repos verwenden die Tabellen-Form produktiv.** Die
   vierte Fähigkeit ist seit `v0.72.0` ausgeliefert; eine ersatzlose
   Ablösung wäre ein Breaking Change ohne Übergangsfrist für jeden Adopter,
   der `observations.register` bereits konfiguriert hat.

## Entscheidung

Das Modul `planning` bekommt einen **zweiten, additiven Konfigurations-Zweig**
für die vierte Fähigkeit: `observations.dir` (Verzeichnis-Modus) steht
gleichberechtigt neben `observations.register` (Tabellen-Modus, unverändert).

1. **Mutually exclusive, fail-closed bei Kollision.** Sind beide gesetzt,
   bricht der Lauf mit Exit 2 — zwei Quellen für dieselbe Fähigkeit wären
   unentscheidbar, welche gilt.
2. **Verzeichnis-Modus, „deklariert" heißt jetzt „Verzeichnis existiert".**
   Statt Tabellenzeilen zu parsen, prüft die Fähigkeit, ob
   `<observations.dir>/<zitierter-pfad>/observation.md` existiert. Die
   übrige Scan-Logik (Zitat-Erkennung, Prosa-vs-Code-Unterscheidung,
   `citedWithoutRow`) ist **unverändert** — sie braucht nur eine
   `declared`-Menge, unabhängig davon, wie diese entsteht.
3. **Die Kennungs-Gestalt ändert sich mit dem Modus.** Der bestehende
   `Pattern`-Schlüssel bleibt der Mechanismus; sein Default für den
   Verzeichnis-Modus wird `BEO-[A-Z][A-Z0-9]*/[a-z][a-z0-9-]*` (Kürzel/Slug)
   statt `BEO-\d{3}` (Tabellen-Modus-Default bleibt unverändert für
   Rückwärtskompatibilität).
4. **Das Kürzel ist repo-deklariert, nicht baseline-vorgegeben.** Es kommt
   aus `harness/conventions.md`s Sub-Area-Modus-Deklaration — eine neue
   Spalte `BEO-Kürzel` neben Modus/Begründung/Graduation. Ohne Kürzel keine
   Migration einer Beobachtung dieser Sub-Area (Vorbedingung für slice-195 <!-- d-check:status-provenance -->).
5. **Die dritte Prüfung der Beleg-Form entfällt strukturell.** Ein
   gespeicherter Zähler existiert im Verzeichnis-Modus nicht — er wird aus
   der Zahl der `evidence/`-Dateien abgeleitet und kann seiner eigenen
   Ableitung nicht widersprechen. Das deckt sich mit dem Kanon-Text: „eine
   dritte Prüfung ist entfallen".
6. **Kein Automatik-Umstieg.** Ein Repo bleibt im Tabellen-Modus, bis es
   `observations.dir` explizit setzt. d-check selbst tut das erst in
   slice-195 <!-- d-check:status-provenance -->, wenn der Bestand migriert ist.

## Verglichene Alternativen

| Alternative | Verworfen, weil |
|---|---|
| **Vollständige Ablösung** der Tabellen-Form durch die Verzeichnis-Form | Bricht sowohl den eigenen Zwischenzustand (Tabelle besteht bis slice-195 <!-- d-check:status-provenance -->) als auch jeden externen Adopter, der `observations.register` bereits produktiv nutzt — ohne Übergangsfrist. |
| **Neues, eigenständiges Modul** statt Erweiterung der vierten `planning`-Fähigkeit | Verdoppelt die Konzept-Fläche (zwei Module für dieselbe Register-Deckungs-Frage) ohne Gewinn — [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)s Rahmen („eine zitierte Kennung hat eine Zeile im Register") ist modusunabhängig gültig. |
| **Migrations-Skript statt Doppel-Modus**, das die Tabelle beim ersten Lauf automatisch in Verzeichnisse konvertiert | Verschiebt die Kürzel-Entscheidung (Sub-Area → Kürzel, nicht automatisierbar) in einen Lauf, der genau dafür keinen Menschen fragt — der Kanon selbst sagt „nachgeschlagen, nicht erfunden". Bleibt slice-195 <!-- d-check:status-provenance -->, mit Kürzel-Deklaration davor. |

## Konsequenzen

- Neue Konfigurationsfläche (`observations.dir`) in `model.PlanningConfig`;
  additiv, `observations.register` unverändert im Verhalten.
- `harness/conventions.md`s Sub-Area-Tabelle bekommt eine neue Spalte
  (Kürzel) — Vorbedingung für slice-195 <!-- d-check:status-provenance -->, kein Bruch für Repos ohne
  Verzeichnis-Modus.
- Lastenheft ([`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)) wächst um eine fünfte Fähigkeit,
  Akzeptanzkriterien-Trio, Versions-Bump.
- Zwei Modi im selben Modul zu pflegen ist echter Mehraufwand — begrenzt
  durch den Re-Evaluierungs-Trigger unten.

## Fitness Function (falls maschinell prüfbar)

Unit-Tests in `internal/hexagon/core/rules/planning_observations_test.go`:
Verzeichnis-Modus meldet eine zitierte Kennung ohne `observation.md`
(Umkehr-Probe), meldet **nicht**, wenn die Datei existiert, und der
Kollisions-Fall (`Dir` und `Register` beide gesetzt) bricht mit dem
erwarteten Grund-Code.

## Re-Evaluierungs-Trigger

Sobald slice-195 <!-- d-check:status-provenance --> den eigenen Bestand migriert **und** ein voller
Wellen-Zyklus mit dem Verzeichnis-Modus im eigenen Betrieb vergangen ist,
ist zu prüfen, ob der Tabellen-Modus für d-checks eigene Nutzung noch
gebraucht wird oder als Altlast markiert werden kann — dieselbe Reifung
wie bei ADR-0081s eigenem Trigger. Die Entscheidung über andere
Adopter-Repos liegt nicht bei d-check.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-09-03 | Proposed → Accepted (slice-194) |
| 2026-09-03 | Korrektur zu §Fitness Function: der Kollisions-Fall (`Dir` und `Register` beide gesetzt) bricht beim Konfigurations-Laden mit einem generischen `error`, nicht mit einem strukturierten Grund-Code — Grund-Codes (`model.Reason*`) existieren nur für Scan-Befunde, nicht für Config-Fehler. Der ursprüngliche Wortlaut suggerierte fälschlich einen Grund-Code auch für diesen Fall (unabhängiger Review, F-3) |
