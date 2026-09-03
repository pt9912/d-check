# Slice slice-195: Beobachtungs-Register — Datenmigration

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-88](../welle-88-baseline-v600-migration.md) — dritter von
drei Slices: migriert den Bestand auf die in slice-194 gebaute Fähigkeit.

**Bezug:** [`MR-013`](../../../../harness/conventions.md#mr-013)
(Lifecycle-Move-Bündelung — die Umstellung ist ein Move mit
Inhaltsänderung); die Kürzel-Deklaration aus slice-194.

**Berührte Spec-Stellen:** — (Datenmigration, keine neue Anforderung).

**Verantwortlich:** —.

**Autor:** pt9912. **Datum:** 2026-09-03.

---

## 1. Ziel

**`docs/plan/planning/observations.md` (~27 Einträge, eine Tabelle) wird
`docs/plan/planning/observations/` (ein Verzeichnis je Beobachtung).** <!-- d-check:ignore (Ziel-Form, entsteht erst mit diesem Slice) --> Jede
Zeile wird zu `observations/BEO-<KUERZEL>/<slug>/` mit `observation.md`
(Bezeichnung, Sub-Area — unveränderlich), `state.md` (Stand: `offen` |
`verkörpert …` | `geplant …` | `gestrichen …`) und `evidence/<vorgangs-id>.md`
je bisherigem Beleg (ein Beleg-Zähler-Eintrag wird eine Datei). Jedes
lebende `BEO-<NNN>`-Zitat im Repo wird auf den neuen Pfad umgehängt;
eingefrorene Bestände (`done/`, `docs/reviews/`, `harness/conventions/done/`)
bleiben unangetastet — sie zitieren den Stand ihrer Zeit.

**Voraussetzung:** slice-194 in `done/` (die Fähigkeit, die den neuen
Zustand überhaupt prüfen kann, muss vor der Migration stehen — sonst liefe
`make gates` nach dieser Migration gegen ein Format, das kein Sensor kennt).

## 2. Vorgehen

1. **Kürzel/Slug je Bestandseintrag festlegen** — 27 Einträge, Kürzel aus
   slice-194s Deklaration, Slug lowercase-kebab-case aus der Beobachtung
   (keine Übersetzung der Bezeichnung, nur ein kurzer, eindeutiger
   Bezeichner).
2. **Verzeichnisse anlegen** — je Eintrag `observation.md` (Bezeichnung +
   Sub-Area, wortgleich aus der alten Zeile übernommen, nicht
   umformuliert), `state.md` (Stand-Spalte übernommen), `evidence/`
   (ein File je bisherigem Beleg-Eintrag der Belege-Spalte).
3. **Alle lebenden Zitate nachziehen** — `grep -rn "BEO-[0-9]\+"` über alle
   Nicht-Frozen-Pfade, jede Fundstelle auf den neuen Pfad umgehängt.
   [`MR-054`](../../../../harness/conventions.md#mr-054),
   [`ADR-0082`](../../adr/0082-uebergangswaechter-reviews-observations.md) und
   jeder künftige Slice-Plan-Vorprüfungs-Block sind die Hauptkandidaten.
4. **`observations.md` entfernen**, `observations/README.md` neu (aus dem
   Template `observation.template.md`s Nachbarschaft — die Register-eigene
   Kopfnotiz).
5. **Templates propagieren** — die eigene Haus-Slice-Kopfform (DoD-Zeile
   „Beobachtungs-Register fortgeschrieben" + Closure-Vorlagen-Zeile, wie
   zuletzt in slice-193/194 verwendet) auf die neue Zitierform;
   `welle.template.md`/`welle-results.template.md` falls sie das Register
   erwähnen.

## 3. Ausdrücklich NICHT in diesem Slice

- **Neuvergabe abgeschlossener, eingefrorener Zitate** in `done/`,
  `docs/reviews/`, `harness/conventions/done/` — sie bleiben, wie sie sind.
- **Die in [`BEO-027`](../observations.md) benannte Sensor-Lücke** zu
  beheben — falls die neue `state.md`-Form sie nicht schon strukturell
  entschärft, bleibt sie als eigener, benannter Punkt bestehen.
- **Eine Anpassung von `tools/archive-wave`** an das neue Format — das
  Werkzeug liest das Beobachtungs-Register heute nicht.

## 4. Definition of Done

- [ ] Alle ~27 Bestandseinträge migriert (Zahl bei Ausführung exakt
      erhoben, nicht aus diesem Plan übernommen).
- [ ] `docs/plan/planning/observations.md` entfernt,
      `docs/plan/planning/observations/README.md` neu. <!-- d-check:ignore (Ziel-Form, entsteht erst mit diesem Slice) -->
- [ ] 0 lebende `BEO-<NNN>`-Zitate außerhalb der Frozen-Bestände
      (`grep -rn "BEO-[0-9]\+"` abzüglich `done/`, `docs/reviews/`,
      `harness/conventions/done/`).
- [ ] Templates (mind. `slice.template.md`) auf die neue Zitierform
      nachgezogen.
- [ ] `make gates` grün (zehn Gates) — insbesondere `planning-check` gegen
      das neue Verzeichnis.
- [ ] `make fullbuild` grün — `verify-closure-notes` gegen die neue Form.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/`.
- [ ] Unabhängige Verifikation durchgeführt.
- [ ] Closure-Notiz (§9) geschrieben, jedes Risiko aus §5 mit Ausgang.

## 5. Abnahme-Punkte / Risiken

- **Größte Einzel-Risiko: stille Bedeutungsverschiebung beim Umzug.** 27
  Einträge von Hand umzuhängen ist fehleranfällig — ein falsch
  zugeordneter Beleg oder ein verlorener Zähler-Stand wäre ein echter
  Informationsverlust am Steering-Loop selbst. Gegenmaßnahme: Diff der
  Beleg-**Anzahl** je migriertem Eintrag gegen die alte Zähler-Spalte vor
  dem Löschen von `observations.md`.
- **Zeitgleiche Arbeit an einem anderen Slice könnte einen neuen
  `BEO-<NNN>`-Eintrag anlegen**, während dieser Slice läuft — WIP-Limit 1
  schützt dagegen, ist aber nur so gut wie seine Einhaltung.
- **Der Umfang (27 Einträge) könnte die Ein-Sitzungs-Review-Grenze
  sprengen** — falls so, Rückführung nach `next` und Aufteilung nach
  Sub-Area oder Kennungsblock, nicht stilles Weiterarbeiten über die
  Grenze hinaus.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei, slice-194 in `done/`
(Kürzel-Deklaration und Produktcode-Fähigkeit stehen).

**Rückführungen:** `in-progress` → `next`, falls §5s dritter Punkt eintritt
(Umfang sprengt eine Review-Sitzung) — dann Aufteilung nach Sub-Area.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `docs/plan/planning/` (Planning-Harness selbst).
  Fällt unter den Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration).

  <!-- d-check:cite .harness/baseline/v6.0.0/regelwerk/modul-05-planning-harness.md:223-224 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Stand 2026-09-03, höchste Kennung
  `BEO-027`): [`BEO-027`](../observations.md) (**1**) — direkt einschlägig,
  diese Migration ist die Gelegenheit, die dort benannte Lücke (Zeile ohne
  Ausgang oberhalb der Schwelle) beim Umzug sichtbar zu machen, auch wenn
  ihre Behebung außerhalb des Scope bleibt (§3); [`BEO-025`](../observations.md)
  (**1**, ein Bump zerfällt in mehrere Commits, Zuordnung vor dem Editieren
  klären) — gilt für den Migrationsumfang ebenso wie für einen Baseline-Bump.

  <!-- d-check:cite .harness/baseline/v6.0.0/regelwerk/modul-05-planning-harness.md:229-229 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)) — bei
  Beanspruchung neu zu lesen. **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-195. Betroffene IDs:
[`MR-013`](../../../../harness/conventions.md#mr-013). Module: `planning`.
Gates: `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Planning-Harness fällt unter den
Default: Doc führt, Code folgt. Reine Datenmigration innerhalb einer
bereits etablierten Konvention (das Register selbst), kein neuer
Produktcode (der ist slice-194). **Evidenz-/Diskrepanz-Risiko** ist die
Achse mit Substanz — siehe §5, stille Bedeutungsverschiebung beim
Von-Hand-Umzug.

## 9. Closure-Notiz (nach `done/`)

<!-- wird erst bei Closure gefüllt -->
