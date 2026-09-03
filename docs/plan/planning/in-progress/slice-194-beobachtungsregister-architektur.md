# Slice slice-194: Beobachtungs-Register — neue Architektur

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-88](../welle-88-baseline-v600-migration.md) — zweiter von
drei Slices: liefert die **Fähigkeit**, migriert noch keine Daten (das ist
slice-195, blockiert von diesem Slice).

**Bezug:** neue ADR (Nummer bei Anlage vergeben, [MR-008](../../../../harness/conventions.md#mr-008)); Erweiterung von
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
um eine fünfte Fähigkeit; Kürzel-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md) §Modus-Deklaration
pro Sub-Area (nur dort sind Sub-Areas benannt, das neue Kürzel hängt sich
daran).

**Berührte Spec-Stellen:** [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) (Erweiterung).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-03.

---

## 1. Ziel

**Das Beobachtungs-Register wird eine Verzeichnis-Ablage statt einer
Tabellen-Datei — die Fähigkeit dafür entsteht hier, der Bestand zieht erst
in slice-195 um.** Baseline `v6.0.0` (`modul-06-roadmap.md` §Das
Beobachtungs-Register) legt die Ziel-Form fest:

```text
docs/plan/planning/observations/
├── README.md
└── BEO-<KUERZEL>/<slug>/
    ├── observation.md
    ├── state.md
    └── evidence/<vorgangs-id>.md
```

Kennung = Pfad, kein zentraler Zähler mehr (er wird aus der Zahl der
`evidence/`-Dateien abgeleitet). `<KUERZEL>` ist **nachgeschlagen, nicht
erfunden** — es kommt aus der Sub-Area-Modus-Deklaration in
`harness/conventions.md`; heute deklariert dort nur `*` (Default) und
`tools/harness/` — beide brauchen ein Kürzel, bevor die erste Beobachtung
migriert werden kann (Vorbedingung für slice-195).

**Produktseitig** liest `internal/hexagon/core/rules/planning_observations.go`
(`CheckPlanningObservations`, vierte `planning`-Fähigkeit, [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in))
heute eine **Tabellen-Datei** (`observations.register`, ein Regex gegen
`declaredObservationIDs`, das Tabellenzeilen parst). Diese Fähigkeit muss ein
**Verzeichnis** prüfen können: eine zitierte Kennung (jetzt ein Pfad-Muster
`BEO-<KUERZEL>/<slug>`) hat eine Zeile im Register, wird zu — existiert
`observations/<pfad>/observation.md`. Ob die alte tabellenbasierte Form als
zweiter unterstützter Modus erhalten bleibt (Rückwärtskompatibilität für
andere Adopter-Repos) oder vollständig abgelöst wird, ist Teil der ADR
dieses Slice, keine Vorentscheidung hier.

## 2. Vorgehen

1. **ADR schreiben** — Entscheidung, Alternativen (u. a.: Verzeichnis-Modus
   als zweite Option neben der Tabellen-Form vs. vollständige Ablösung),
   Konsequenzen, Re-Evaluierungs-Trigger. Vorbild: [ADR-0082](../../adr/0082-uebergangswaechter-reviews-observations.md) (vergleichbare
   Größe — neue Konfigurations-Fläche für eine bestehende Modul-Fähigkeit).
2. **Kürzel deklarieren** — `harness/conventions.md` §Modus-Deklaration pro
   Sub-Area bekommt eine Spalte oder einen Zusatz-Eintrag für das
   BEO-Kürzel je Sub-Area (`*` → z. B. `ALL`, `tools/harness/` → z. B.
   `HARN` — exakte Wahl bei Ausführung, Kollisionsfreiheit prüfen).
3. **Produktcode erweitern** (`internal/hexagon/core/rules/planning_observations.go`
   + `model.PlanningConfig`): neuer Konfigurationszweig für die
   Verzeichnis-Form, Existenz-Prüfung statt Tabellenzeilen-Parsing,
   Grund-Code(s) unverändert oder erweitert (Entscheidung: ADR).
4. **Tests** — Unit-Tests für den neuen Zweig, mindestens eine Umkehr-Probe
   (erfundener Pfad ohne `observation.md` ⇒ Befund).
5. **Lastenheft-Eintrag** — [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) um die fünfte Fähigkeit
   fortschreiben (Versions-Bump, Historie-Zeile), Akzeptanzkriterien-Trio.
6. **Handbuch/README/operations.md** — die neue Fähigkeit dokumentieren
   (gehört zur Release-Prep des nächsten Produkt-Release, nicht zwingend
   diesem Slice — siehe Abgrenzung in AGENTS.md §5 zu Feature- vs.
   Release-Prep-Commits).

## 3. Ausdrücklich NICHT in diesem Slice

- **Die eigentliche Datenmigration** von `docs/plan/planning/observations.md`
  nach `observations/` — das ist slice-195, das diesen Slice als
  Vorbedingung hat.
- **Der Nachzug aller lebenden `BEO-<NNN>`-Zitate** im Repo — ebenfalls
  slice-195.
- **CHANGELOG-Eintrag / Release** — Release-Prep folgt separat, wenn diese
  Fähigkeit tatsächlich released wird (nach slice-195, wenn der eigene
  Bestand sie auch nutzt).

## 4. Definition of Done

- [ ] ADR geschrieben und `Accepted`.
- [ ] Kürzel für beide heute deklarierten Sub-Areas (`*`, `tools/harness/`)
      in `harness/conventions.md` festgelegt, kollisionsfrei.
- [ ] `internal/hexagon/core/rules/planning_observations.go` erweitert um
      den Verzeichnis-Modus (oder ersetzt — je nach ADR-Entscheidung).
- [ ] Unit-Tests für den neuen Zweig, inkl. Umkehr-Probe.
- [ ] [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) in `spec/lastenheft.md` fortgeschrieben
      (Akzeptanzkriterien-Trio, Versions-Bump, Historie-Zeile).
- [ ] `make gates` grün (zehn Gates), `make test` deckt den neuen Zweig.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/`.
- [ ] Unabhängige Verifikation durchgeführt.
- [ ] Closure-Notiz (§9) geschrieben, jedes Risiko aus §5 mit Ausgang.

## 5. Abnahme-Punkte / Risiken

- **Die Kürzel-Wahl ist eine einmalige Entscheidung mit Bestandswirkung** —
  wird sie später geändert, wandern alle bis dahin migrierten
  Beobachtungs-Pfade mit. Risiko: zu granulare oder zu grobe Kürzel-Wahl
  jetzt, die slice-195 oder einen künftigen Slice zur Nacharbeit zwingt.
- **Rückwärtskompatibilität vs. Ablösung** ist eine echte Weggabelung mit
  Implementierungsaufwand in beide Richtungen — die ADR muss das mit
  Alternativen belegen, nicht als Formsache abhaken.
- **Der bestehende Bestand (`observations.md`) bleibt bis slice-195 die
  Quelle der Wahrheit** — dieser Slice darf `make planning-check`/
  `make gates` nicht rot machen, obwohl er einen zweiten, noch ungenutzten
  Modus einführt (additiv, nicht `observations.register` selbst ändernd).

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei, slice-193 in `done/`
(Pin steht auf `v6.0.0`, das Regelwerk-Delta ist damit verfügbar).

**Rückführungen:** `in-progress` → `next`, falls die ADR-Alternativenwahl
(Rückwärtskompatibilität vs. Ablösung) den Produktcode-Umfang über eine
Review-Sitzung hinaus wachsen lässt — dann wird der Produktcode-Teil ein
eigener Folge-Slice, die ADR bleibt hier.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `internal/hexagon/core/rules/` (Produktcode) und
  `harness/conventions.md` (Konventionsspeicher). Beide fallen unter den
  Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration).

  <!-- d-check:cite .harness/baseline/v6.0.0/regelwerk/modul-05-planning-harness.md:223-224 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Stand 2026-09-03, höchste Kennung
  `BEO-028` — bei der Beanspruchung neu gelesen, seit der Eröffnung um
  eine Instanz gewachsen): [`BEO-027`](../observations.md) (**1**, Sub-Area
  `docs/plan/planning/observations.md`) — beschreibt eine Sensor-Lücke der
  **heutigen** Register-Form (Zeile übersteht die Schwelle ohne Ausgang);
  relevant als Kontext, ob die neue `state.md`-Form sie strukturell
  entschärft, ist Teil der ADR-Abwägung, kein eigener Fix hier.
  [`BEO-028`](../observations.md) (**1**, Sub-Area `tools/harness/`) —
  betrifft `fetch-baseline-cache.sh`s `check-latest`, nicht das
  Beobachtungs-Register; kein Treffer für die Sub-Areas dieses Slice. Keine
  weiteren Treffer für die berührten Sub-Areas.

  <!-- d-check:cite .harness/baseline/v6.0.0/regelwerk/modul-05-planning-harness.md:229-229 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)) — bei
  Beanspruchung neu zu lesen. **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-194. Betroffene IDs: neue ADR (Nummer bei Anlage),
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in).
Module: `planning`. Gates: `make gates`, `make test`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter
den Default: Doc führt, Code folgt. Konventionen-Dichte hoch (Modul-Layout,
Test-Konventionen sind etabliert). **Evidenz-/Diskrepanz-Risiko** ist die
Achse mit Substanz: eine neue Konfigurationsfläche in einem bestehenden,
bereits ausgelieferten Modul ([`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in), seit `v0.72.0`) — Bestand
und Doku dürfen dabei nicht auseinanderlaufen (BEO-008-Klasse, hier auf
Produktcode statt Baseline-Pin angewendet).

## 9. Closure-Notiz (nach `done/`)

<!-- wird erst bei Closure gefüllt -->
