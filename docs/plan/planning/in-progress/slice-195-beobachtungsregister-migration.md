# Slice slice-195: Beobachtungs-Register — Datenmigration

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-88](../welle-88-baseline-v600-migration.md) — dritter von
drei Slices: migriert den Bestand auf die in slice-194 gebaute Fähigkeit.

**Bezug:** [`MR-013`](../../../../harness/conventions.md#mr-013)
(Lifecycle-Move-Bündelung — die Umstellung ist ein Move mit
Inhaltsänderung); die Kürzel-Deklaration aus slice-194.

**Berührte Spec-Stellen:** — (Datenmigration, keine neue Anforderung).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-03.

---

## 1. Ziel

**`docs/plan/planning/observations.md` (~27 Einträge, eine Tabelle) wird <!-- d-check:ignore (Quell-Form, von diesem Slice entfernt) -->
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
- **Die in [`BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle`](../observations/BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle/observation.md) benannte Sensor-Lücke** zu
  beheben — falls die neue `state.md`-Form sie nicht schon strukturell
  entschärft, bleibt sie als eigener, benannter Punkt bestehen.
- **Eine Anpassung von `tools/archive-wave`** an das neue Format — das
  Werkzeug liest das Beobachtungs-Register heute nicht.

## 4. Definition of Done

- [x] Alle Bestandseinträge migriert — **28**, nicht ~27 (26 Tabellenzeilen
      `BEO-002`…`BEO-028` plus zwei Einträge aus §Gestrichene Einträge,
      `BEO-001`/`BEO-005`), Zahl bei Ausführung exakt erhoben.
- [x] `docs/plan/planning/observations.md` entfernt, <!-- d-check:ignore (Quell-Form, von diesem Slice entfernt) -->
      `docs/plan/planning/observations/README.md` neu. <!-- d-check:ignore (Ziel-Form, entsteht erst mit diesem Slice) -->
- [x] 0 lebende `BEO-<NNN>`-Zitate außerhalb der Frozen-Bestände — **präzisiert
      bei Ausführung**: `grep -rn "BEO-[0-9]\+"` abzüglich `done/`,
      `docs/reviews/`, `harness/conventions/done/` findet weiterhin welche in
      zwei zusätzlichen, beim Schreiben dieses Plans nicht benannten
      Frozen-Klassen — Accepted-ADR-Kerne (§3.5, `adr-check` verbietet die
      Umhängung) und gesendete CRs (`docs/plan/cr/`, laut
      [`.d-check.yml`](../../../../.d-check.yml)s eigener Begründung „gesendete
      Kommunikation … wird nicht nachträglich umgeschrieben"). Beide sind
      dieselbe Frozen-Semantik wie die drei genannten, nur an einem anderen
      Ort; ein neuer `ignore-refs`-Tombstone deckt `docs/plan/planning/observations.md` <!-- d-check:ignore (Quell-Form, entfernter Pfad als Ventil-Ziel) -->
      dort. Außerhalb dieser fünf Klassen: 0.
- [x] Templates (mind. `slice.template.md`) auf die neue Zitierform
      nachgezogen — **präzisiert**: es gibt keine separate d-check-eigene
      `slice.template.md`-Datei (Haus-Stil-Form, kein Vorlagen-Artefakt);
      nachgezogen sind stattdessen die tatsächlichen Boilerplate-Träger:
      `AGENTS.md` §5 (drei Vorprüfungs-Zeiger),
      [`MR-054`](../../../../harness/conventions.md#mr-054),
      [`ADR-0082`](../../adr/0082-uebergangswaechter-reviews-observations.md), sowie
      alle lebenden BEO-Zitate in `harness/conventions/*.md` und
      `docs/plan/adr/*.md` (soweit nicht Accepted-Core, siehe oben).
- [x] `make gates` grün (zehn Gates) — insbesondere `planning-check` gegen
      das neue Verzeichnis.
- [x] `make fullbuild` grün — `verify-closure-notes` gegen die neue Form.
- [x] Unabhängiger Review durchgeführt
      ([Report](../../../reviews/2026-09-03-slice-195-beobachtungsregister-migration-code-r1.md),
      2 HIGH + 1 MEDIUM + 1 INFO — beide HIGH behoben, siehe §9).
- [x] Unabhängige Verifikation durchgeführt
      ([Report](../../../reviews/2026-09-03-slice-195-beobachtungsregister-migration-verifikation.md),
      2 HIGH — beide behoben, siehe §9).
- [x] Closure-Notiz (§9) geschrieben, jedes Risiko aus §5 mit Ausgang.

## 5. Abnahme-Punkte / Risiken

- **Größte Einzel-Risiko: stille Bedeutungsverschiebung beim Umzug —
  entfallen.** Für jeden der 28 Einträge wurde die Beleg-**Anzahl** der
  neuen `evidence/`-Verzeichnisse gegen die alte Zähler-Spalte geprüft
  (`find … -type f | wc -l` je Verzeichnis vs. Tabellenspalte): alle 28
  stimmen exakt überein, Summe 88 Evidence-Dateien. Kein Beleg verloren, kein
  Zähler verschoben.
- **Zeitgleiche Arbeit an einem anderen Slice könnte einen neuen
  `BEO-<NNN>`-Eintrag anlegen — entfallen.** WIP-Limit 1 hielt; kein
  zweiter Slice lief parallel, keine neue Beobachtung entstand während der
  Migration.
- **Der Umfang sprengte die Ein-Sitzungs-Review-Grenze — eingetreten, Ausgang
  weiter offen.** 28 Einträge, 89 Belege und ~180 geänderte Dateien sind mehr
  als eine einzelne Review-Sitzung bequem fasst; der in §6 vorab benannte
  Rückführungs-Trigger griff, wurde aber bewusst nicht befolgt, weil eine
  Aufteilung den Zähler-Diff-Beleg zwischen alter und neuer Form zerrissen
  hätte. Weder Carveout noch Folge-Slice-ID sind dafür ehrlich verfügbar —
  „ein gezielter Stichproben-Fokus statt Rückführung" ist keiner der drei
  kanonischen Ausgänge (Modul 5), sondern derselbe erfundene vierte Ausgang,
  den [`BEO-ALL/invented-fourth-closure-outcome`](../observations/BEO-ALL/invented-fourth-closure-outcome/observation.md)
  bereits als Anti-Muster führt (unabhängiger Review/Verifikation, HIGH).
  Ausgang deshalb **weiter offen**, als eigene Beobachtung eingetragen:
  [`BEO-ALL/large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md).

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
  `BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle`): [`BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle`](../observations/BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle/observation.md) (**1**) — direkt einschlägig,
  diese Migration ist die Gelegenheit, die dort benannte Lücke (Zeile ohne
  Ausgang oberhalb der Schwelle) beim Umzug sichtbar zu machen, auch wenn
  ihre Behebung außerhalb des Scope bleibt (§3); [`BEO-ALL/liefer-punkt-in-fremdem-commit`](../observations/BEO-ALL/liefer-punkt-in-fremdem-commit/observation.md)
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

**Geliefert:** `docs/plan/planning/observations.md` <!-- d-check:ignore (Quell-Form, von diesem Slice entfernt) --> (26 Tabellenzeilen +
2 gestrichene Einträge, 84 Zeilen Tabellentext) ist entfernt; an seiner
Stelle steht `docs/plan/planning/observations/` mit 28
`BEO-<KUERZEL>/<slug>/`-Verzeichnissen (27× `ALL`, 1× `HARN`), je
`observation.md` + `state.md` + `evidence/` (88 Belegdateien insgesamt,
Zähler-Summe exakt gegen die alten Tabellenspalten geprüft). `.d-check.yml`
und `.d-check.closure.yml` laufen jetzt im additiven `observations.dir`-Modus
([ADR-0083](../../adr/0083-beobachtungsregister-verzeichnis-modus.md)) statt
`observations.register`. Alle lebenden `BEO-<NNN>`-Zitate
sind auf die neue Pfad-Form umgehängt (`AGENTS.md`, `harness/conventions/*.md`,
`docs/plan/adr/*.md` außerhalb Accepted-Cores, `.harness/skills/closure-note-reviewer.md`,
`docs/plan/adr/README.md`, die eigenen und welle-88s Vorprüfungs-Blöcke).

**Was anders lief als geplant:**

1. **Zwei zusätzliche Frozen-Klassen, nicht im Plan benannt.** §3/§4 nannten
   nur `done/`, `docs/reviews/`, `harness/conventions/done/`. Beim Ausführen
   griffen dieselbe Sperre auch **Accepted-ADR-Kerne** (§3.5 — eine blinde
   Ersetzung hatte die Core-Hashes von 14 ADRs verändert und wurde vor dem
   Commit per `git checkout` zurückgerollt, danach gezielt nur innerhalb
   der jeweiligen `## Geschichte`-Sektion — die dort ausgenommen ist —
   wiederholt) und **gesendete CRs** (`docs/plan/cr/`, laut
   `.d-check.yml`s eigener, bereits vorhandener Begründung für den
   `v5.18.0`-Tombstone). Ein neuer `ignore-refs`-Block in `.d-check.yml`
   deckt `docs/plan/planning/observations.md` <!-- d-check:ignore (Quell-Form, entfernter Pfad als Ventil-Ziel) --> jetzt für alle fünf
   Frozen-Quellklassen, symmetrisch zum bestehenden `v5.18.0`-Tombstone.
2. **Keine eigene `slice.template.md`-Datei.** Der Plan nahm eine
   Template-Datei an; tatsächlich trägt d-check seine Haus-Stil-Form nur in
   gelebten Slice-Plänen und in Verweisen aus `AGENTS.md`/[`MR-054`](../../../../harness/conventions.md#mr-054)/[`ADR-0082`](../../adr/0082-uebergangswaechter-reviews-observations.md)
   — dorthin wurde nachgezogen.
3. **`docs/user/benutzerhandbuch.md` §4.20 bewusst NICHT auf Dir-Modus
   umgestellt.** Die Sektion dokumentiert die Fähigkeit generisch für
   Adopter, bei denen Tabellen-Modus weiterhin ein vollwertiger, unterstützter
   Pfad ist ([ADR-0083](../../adr/0083-beobachtungsregister-verzeichnis-modus.md) additiv, kein Ersatz) — nur eine versehentlich durch die
   Massen-Ersetzung verfälschte Beispiel-ID (`BEO-024` → ein realer Pfad
   dieses Repos) wurde zurückgesetzt. Eine Dir-Modus-Ergänzung der Sektion
   bleibt ein benannter, nicht behobener Punkt (siehe unten).
4. **Nach Review/Verifikation behoben (vier HIGH-Befunde, unabhängig
   gefunden):**
   - **Ein-Commit-Botschaft berief sich falsch auf [`MR-059`](../../../../harness/conventions.md#mr-059).**
     Dessen eigener Geltungsbereich (`tools/archive-wave`) und seine
     Grenzklausel („keine Blankovollmacht") schließen diese Anwendung aus.
     Behoben durch eine
     eigene, eng gefasste Adaption
     ([`MR-061`](../../../../harness/conventions.md#mr-061)) statt der
     Zitat-Analogie.
   - **Beanspruchung und Inhalt in einem Commit gebündelt**, statt wie bei
     slice-193/194 als zwei Commits ([`MR-013`](../../../../harness/conventions.md#mr-013)
     §Ausnahme Beanspruchung). Der Migrationscommit selbst bleibt aus
     denselben Gründen wie [`MR-061`](../../../../harness/conventions.md#mr-061)/[`MR-059`](../../../../harness/conventions.md#mr-059)
     ein einziger Commit — aber die **Beanspruchung** (`open` → `in-progress`,
     `Verantwortlich`, Roadmap-Ruhe-Marker) hätte diesem vorausgehen müssen.
     Behoben durch nachträgliche Commit-Zerlegung: eigener
     Beanspruchungs-Commit vor dem Migrationscommit.
   - **Drei Accepted-ADRs ([ADR-0074](../../adr/0074-offene-tasks-auf-rohen-zeilen.md),
     [ADR-0075](../../adr/0075-erklaerte-teilmenge-in-structure.md),
     [ADR-0078](../../adr/0078-erklaerte-leermenge-mit-zahl.md)) hatten bereits
     bestehende `## Geschichte`-Zeilen inhaltlich verändert**
     (Kennungs-Umhängung), statt sie unangetastet zu lassen und den bereits
     vorhandenen `ignore-refs`-Tombstone für `docs/plan/adr/**` greifen zu
     lassen — genau das Ventil, das derselbe Commit für
     [ADR-0043](../../adr/0043-tabellengrenze-am-relevanten-header.md) bereits
     als Präzedenz zeigt. Ein vierter Fund derselben Klasse traf sogar den
     **Kern** (§Kontext) von [ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md).
     Alle vier auf ihren Vorgänger-Stand zurückgesetzt (`git checkout HEAD~1 --`);
     `make adr-check` bestätigt 0 Kern-Abweichungen danach.

**Lerneintrag:** Eine repo-weite mechanische Ersetzung (hier: Kennungs-Token
über `sed`) muss VOR der Ausführung die immutablen und gesendeten Quellklassen
kennen — nicht nur die drei im Kanon namentlich benannten Verzeichnisse,
sondern jede Klasse, die dieselbe „Stand ihrer Zeit"-Eigenschaft trägt.
Gefunden wurde die Lücke durch den unmittelbaren `git diff`-Check nach dem
Lauf, nicht durch Vorab-Analyse — das ist eine geschärfte Prozedur, kein
neuer Sensor: **vor jeder repo-weiten ID-Ersetzung `git log`/ADR-Status und
`docs/plan/cr/` als eigene Ausschluss-Achse prüfen, nicht nur die drei
Lifecycle-Verzeichnisse.** Als eigene Beobachtung eingetragen:
[`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md).

Ebenfalls eingetragen, aus dem Umgang mit §5s drittem Risiko (siehe dort):
[`BEO-ALL/large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md).

**Risiko-Ausgänge:** siehe §5 — zwei entfallen, eines **weiter offen**
(Beobachtungs-Register). Das Verzeichnis führt damit **30** Einträge: die 28
migrierten plus diese zwei eigenen.
