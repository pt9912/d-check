# Review — slice-071 Closure (unabhängig, kontext-getrennt) vor Tag v0.44.0

**Datum:** 2026-07-17
**Review-Art:** unabhängiger **Closure**-Review (kontext-getrennt; adversarial —
Auftrag war, Gründe zum Stoppen des Releases zu finden)
**Gegenstand:** der **vollständige** Slice
[`slice-071`](../plan/planning/open/slice-071-trace-cross-consistency-gate.md),
Diff `2562dc2..HEAD` (`c4439e8`), neun Commits, 29 Dateien
**Vor-Reviews (Kontext, nicht meine Meinung):**
[R1](2026-07-17-slice-071-implementation-r1.md) (REJECT),
[R2](2026-07-17-slice-071-implementation-r2.md) (REJECT),
[R3](2026-07-17-slice-071-implementation-r3.md) (REJECT),
[R4](2026-07-17-slice-071-implementation-r4.md) (REJECT)
**Reviewer:** Claude (kontext-unabhängiger Lauf; weder Autor noch Vor-Reviewer)
**Skill:** `.harness/skills/reviewer.md` v1.2.0
**Modell:** Opus 4.8 (1M context)

## Eingangs-Kontext

- Vertrag: `DC-FA-XREF-001` ([`spec/lastenheft.md`](../../spec/lastenheft.md),
  Version **0.44.2**), `DC-FA-XREF-001.a`
  ([`spec/spezifikation.md`](../../spec/spezifikation.md), Schritte 1–7 +
  Fehlerpräzedenz), [ADR-0038](../plan/adr/0038-trace-cross-consistency.md)
  (Entscheidungen 1–8, Status `Proposed`)
- Mit-Verträge: `DC-FA-CLI-009`, `DC-FA-CLI-011`, `DC-FA-REQ-001`,
  `DC-FA-COV-001`, `DC-QA-02`, `DC-QA-03`
- Prozess: [`docs/user/releasing.md`](../user/releasing.md) §Release-Prep,
  [`AGENTS.md`](../../AGENTS.md) §3 (Harte Regeln) + §5 (Dokumentations-Regeln)

**Ausgeführte Sensoren (echte Ausgabe, kein Papier-Grün):**

- `make gates` → **grün**: `doc-check` (224 Dateien, 0 Befunde) + `lint` +
  `test` + `arch-check` + `coverage-gate` (**93,90 %** ≥ 93 %) + `semgrep`
  (55 Regeln, 42 Dateien, **0 Findings**) + `gate-consistency` +
  `planning-check`.
- `make build` → `d-check:latest`, danach **sieben** Repro-Läufe gegen das
  gebaute Image (`--network none`, `-v <fixture>:/repo:ro`): Handbuch-§4.12-
  Beispiel 1:1, übergriffiges `exclude-req`, fehlgreifendes `design-pattern`,
  Bootstrap-Zustand (einseitig leere Vorwärts-Sicht), konsistenter Stand in
  Markdown/JSON/YAML, Drift in JSON.
- `git status --porcelain` → leer; `git clean -nd` → nur gitignorierte
  Verzeichnisse (`.agents/`, `.codex/`, leeres `docs/plan/planning/open/`).

## Findings

### F-1 — `--json`/`--yaml` verschweigen den Abgleich-Beleg bei 0 Differenzen

- **kategorie:** MEDIUM
- **quelle:** `DC-FA-XREF-001` / `DC-FA-CLI-009` (maschinenlesbare Matrix);
  Konsistenz-Lücke zwischen den drei Renderern derselben Eingabe-Klasse
- **pfad:** `internal/hexagon/core/app/trace.go:55` (`omitempty`) gegen
  `internal/adapter/driven/report/report.go:285-297`
- **befund:** Der Markdown-Reporter druckt bei aktivem
  `trace.cross-consistency`-Block auch dann `## Kreuzverweis-Konsistenz` +
  `0 Differenz(en).`, wenn nichts gefunden wurde — mit ausdrücklicher
  Begründung im Code („Schweigen wäre nicht von ‚Abgleich lief nicht'
  unterscheidbar", `report.go:281-284`; wortgleich als Test-Invariante in
  `cli_acceptance_test.go:2589-2590`). Die `--json`/`--yaml`-Fassung hält
  diese Invariante **nicht**: `CrossConsistency` trägt `omitempty`, und
  `CrossActive` ist `json:"-"`. Verifiziert gegen `d-check:latest` — bei
  aktivem Block und 0 Differenzen ist das JSON **byte-identisch** zu einem
  Lauf ganz ohne Block: der Schlüssel `crossConsistency` fehlt in beiden.
  Kein Test pinnt diese Nullmengen-Form; die einzige JSON-Zusicherung
  (`crossDoc`, `cli_acceptance_test.go:2478-2488`) prüft nur den
  **nicht**-leeren Fall.
- **failure-szenario:** Eine CI-Stufe konsumiert `--trace --json` und wertet
  `crossConsistency` aus. Die `.d-check.yml` wird versehentlich nicht
  gemountet oder der Block fällt bei einem Merge heraus ⇒ der Schlüssel fehlt
  ⇒ die Pipeline liest „keine Differenzen" und bleibt grün, obwohl nie
  verglichen wurde. Genau die Silent-Green-Klasse, gegen die R1–R4 vier Runden
  lang den Vakuitäts-Guard erkämpft haben — hier im maschinenlesbaren Kanal.
- **verifizierbar:** ja —
  `docker run --rm --network none -v <fixture>:/repo:ro d-check:latest --trace --json`
  mit aktivem Block und `F = B`; Ausgabe endet bei `"orphans": 0`.

### F-2 — Handbuch-Feldliste für `--trace --json`/`--yaml` nennt `crossConsistency` nicht

- **kategorie:** LOW
- **quelle:** [`releasing.md`](../user/releasing.md) §Release-Prep Punkt 4
  („Prosa-Currency von Hand nachziehen — kein Gate erzwingt sie")
- **pfad:** `docs/user/benutzerhandbuch.md:656-659`
- **befund:** Die §4.12-Feldliste ist eine **geschlossene** Aufzählung
  („`requirements` mit `id`/`title`/`adrs`/`slices`, konditional `coverage`
  und `modality`, sowie `orphan`; auf Matrix-Ebene `total`/`orphans`") und
  wurde im Release-Prep-Commit nicht um das neue Top-Level-Feld
  `crossConsistency` erweitert, das das Binary nachweislich emittiert. Die
  §4.12-Prosa und die `operations.md`-`--trace`-Zeile beschreiben den Abgleich
  ausschließlich als Markdown-Abschnitt.
- **failure-szenario:** Ein Konsument liest die Feldliste als vollständig,
  baut seinen Parser darauf und erfährt vom Kreuzverweis-Kanal nichts — bzw.
  hält ein auftauchendes `crossConsistency` für undokumentiertes Rauschen.
- **verifizierbar:** ja — Feldliste gegen `--trace --json` mit aktivem Block.

### F-3 — Lastenheft-AK 0.44.2 überschreibt die reale Ventil-Diagnose

- **kategorie:** LOW
- **quelle:** `DC-FA-XREF-001`, Akzeptanzkriterium „Negative (übergriffiges
  Ventil, fail-closed)"
- **pfad:** `spec/lastenheft.md:789-790` gegen
  `internal/hexagon/core/app/trace_cross.go:308-330`
- **befund:** Das AK sagt: „die Meldung benennt das **Ventil**, *nicht die
  korrekten Muster*." Die tatsächliche Meldung ist für beide Ursachen
  **byte-identisch** und nennt die (korrekten) Muster **zuerst**:
  `… beide Sichten ergaben 0 Kanten … prüfe, ob trace.cross-consistency.forward.design-pattern
  den Artefakt-Namensraum beider Sichten trifft — und ob exclude-req nicht alle
  Anforderungen verschluckt`. Die tragende Hälfte des AK (Exit 2 + Nennung des
  Ventils) hält und ist getestet
  (`TestCrossConsistencyVakuumDurchUebergriffigesVentil`); die
  Abgrenzungs-Hälfte („nicht die korrekten Muster") ist unerfüllt und
  untestbar. Der Widerspruch entstand aus der Reihenfolge: der Hint stammt aus
  `815a2a7`, das AK aus `ea0f343`.
- **failure-szenario:** Ein Betreiber mit korrekten Mustern und übergriffigem
  `exclude-req` wird von der Meldung zuerst in die Muster-Config geschickt —
  exakt die Fehldiagnose, die `vacuityHint` laut eigenem Kommentar vermeiden
  wollte. Nachgelagert: wer das AK als Test formuliert („Meldung enthält
  `design-pattern` **nicht**"), bekommt Rot bei korrektem Code.
- **verifizierbar:** ja — reproduziert mit `exclude-req: '.'` bei sonst
  korrekter Config (Meldung identisch zum Fall `design-pattern: 'ZZ-XX-…'`).

### F-4 — Slice-§2-Vertrag und §Bezug nicht auf Lastenheft 0.44.2 nachgezogen

- **kategorie:** LOW
- **quelle:** Maintainability / `AGENTS.md` §5 (Doku-Regeln)
- **pfad:**
  `docs/plan/planning/open/slice-071-trace-cross-consistency-gate.md:7-9,
  55-59, 63-66`
- **befund:** §Bezug und die DoD-Zeile führen weiterhin „Lastenheft 0.44.0"
  und „ADR-0038 (Proposed)"; §2 „Vertrag aus dem Change Request" listet die
  Fail-closed-Klassen ohne die **Vakuitäts-Stufe** und §3 die Tests ohne den
  Vakuitäts-Fall — obwohl beide der eigentliche Ertrag der vier Review-Runden
  sind und Lastenheft/Spezifikation/ADR/CHANGELOG/Handbuch sie führen. Das
  Slice-Dokument ist damit das einzige Vertrags-Artefakt auf dem Stand vor R1.
- **failure-szenario:** Die Closure-Notiz und die `done/`-Ablage frieren einen
  §2-Vertrag ein, der die zentrale Fail-closed-Klasse nicht nennt; ein späterer
  Leser rekonstruiert aus dem Slice einen Gate-Umfang, den das Produkt
  übererfüllt — und hält den Vakuitäts-Guard für unbelegt.
- **verifizierbar:** nein (kein Gate deckt Slice-Prosa) — Sicht-Vergleich.

### F-5 — ADR-0038-Historie: die beiden `2026-07-17`-Zeilen stehen verkehrt herum

- **kategorie:** LOW
- **quelle:** Maintainability (Historien-Konvention der ADR-Familie)
- **pfad:** `docs/plan/adr/0038-trace-cross-consistency.md:177-179`
- **befund:** Die Tabelle läuft aufsteigend (`2026-07-16` zuerst), aber der
  **spätere** 0.44.2-Nachtrag (Zeile 178) steht **über** dem früheren
  0.44.1-Nachtrag (Zeile 179). Wer die Historie von oben liest, erfährt zuerst
  von der Korrektur und danach von dem, was korrigiert wurde.
- **failure-szenario:** Nach `Accepted` ist die ADR immutable (`AGENTS.md`
  §3.5); erlaubt bleiben nur `## Geschichte`-**Anhänge**. Eine Umsortierung
  der bestehenden Zeilen ist dann ein `core-drift-vcs` im `adr-check` — die
  Reihenfolge ist ab dem Statuswechsel eingefroren.
- **verifizierbar:** ja — `make adr-check` nach einem späteren Sortier-Versuch.

### F-6 — Kopfkommentar der Handbuch-Beispiel-Verankerung ist stale

- **kategorie:** LOW
- **quelle:** Maintainability (Doku-Drift, veraltete Enumeration)
- **pfad:** `internal/adapter/driving/cli/handbook_examples_test.go:22-26`
- **befund:** Der Kommentar deklariert „VERANKERT (**8**)" samt Aufzählung;
  `c4439e8` hat ein neuntes Beispiel (§4.12 Kreuzverweis) ergänzt, ohne Zahl
  oder Aufzählung nachzuziehen. Der Kommentar ist das Inventar, gegen das die
  Fail-closed-Klassifikation begründet wird.
- **failure-szenario:** Der nächste Autor liest die Liste als vollständig und
  nimmt an, für §4.12 gebe es keinen E2E-Anker — und ergänzt einen zweiten.
- **verifizierbar:** nein (Kommentar) — Abzählen der `handbookExamples()`.

### F-7 — Roadmap-Closure-Trigger fordert, was der Slice ehrlich als offen führt

- **kategorie:** INFO
- **quelle:** `DC-FA-PLAN-001`-Domäne / Lifecycle-Konsistenz
- **pfad:** `docs/plan/planning/in-progress/roadmap.md:29-38` gegen
  `…/slice-071-…:72-73, 85-93`
- **befund:** Der §Closure-Trigger verlangt „Der reale grid-gym-Drift wird
  geflaggt" und „ADR-0038 ist `Accepted`". Der Realdatenbeleg ist in der DoD
  **korrekt** als offen markiert und hängt laut §4 an Konsumenten-Vorarbeit
  (§27.1 auf konkrete IDs restrukturieren) — er ist vor dieser Vorarbeit nicht
  erreichbar. Die DoD-Zeile „ADR + Index" ist ihrerseits mit „Status Proposed"
  abgehakt. Welle-60 kann damit nach ihrem eigenen Trigger nicht schließen.
  **Kein Tag-Blocker** (der Tag liegt vor der Closure), aber vor dem
  `git mv` nach `done/` ist eine Entscheidung fällig: Trigger auf den
  reproduzierten Teilstand nachziehen oder die Welle offen halten. Der
  `Accepted`-Schwenk ist zudem eine Wiederholungsklasse — bei ADR-0037 blieb er
  bei der Closure liegen und wurde erst mit `13aa58e` nachgezogen.
- **verifizierbar:** ja — `make planning-check` beim Lifecycle-Move.

## Negativbefunde (geprüft, ohne Befund)

- **`releasing.md` §Release-Prep Punkt 1 (`version.md`):** §Aktuell auf
  `v0.44.0`/2026-07-17, neue Verlaufszeile vorhanden, der `<a id="v0.44.0">`-
  Anker ist **gewandert** — die `v0.43.1`-Zeile trägt keinen Anker mehr.
  Genau ein Anker im Register. Korrekt.
- **Punkt 2 (alle gepinnten `ghcr`-Verweise):** `grep -rn '0\.43\.1'` über den
  Arbeitsbaum findet **keinen** vergessenen Pin. Die 11 Resttreffer sind
  sämtlich legitim: CHANGELOG-Abschnitt, `version.md`-Verlauf,
  ADR-0037-Historie, Roadmap-Vorgängersatz, Review-Zitat sowie vier
  Handbuch-**Prosa**-Stellen („Ab v0.43.1 …", Mindestversions-Hinweis) und die
  §11-Zeilen 1.28/1.29. 27 `ghcr.io/pt9912/d-check:v…`-Pins in README (beide
  Sprachen) und Handbuch stehen auf `v0.44.0`.
- **Punkt 4 bare-Tag-Beispiel:** `benutzerhandbuch.md:74` (`:v0.44.0` ohne
  `ghcr`-Präfix, vom `versions`-Gate **nicht** erfasst) ist nachgezogen.
- **Punkt 4 Handbuch-Kopfstempel:** `1.30` / `[v0.44.0](../../version.md#v0.44.0)`
  / Stand 2026-07-17 — der Anker existiert (Anker-Wanderung konsistent).
- **Punkt 4 §11-Zeile:** `1.30 | v0.44.0 | 2026-07-17` steht **chronologisch
  unter** der letzten (1.29), nicht oben. Inhalt deckt sich mit CHANGELOG und
  Handbuch-§5.
- **Punkt 4 README beide Sprachen:** `README.de.md` (kanonisch) und
  `README.md` synchron. Status-Zeile/Modul-Liste/Intro-Framing **korrekt
  unangetastet** — `trace.cross-consistency` ist eine `trace`-Unterfähigkeit,
  **kein** 18. Regelmodul (Modul-Tabelle Handbuch §6 unverändert, `--enable`-
  Enumeration unverändert).
- **Punkt 4 `operations.md`:** Modul-Enumeration korrekt unberührt (kein neues
  Modul); die `--trace`- und `--require-complete`-Zeilen sind beide um die neue
  Semantik erweitert (inkl. „beide Ursachen melden getrennt auf stderr");
  keine neue CLI-Option ⇒ Optionen-Tabelle vollständig.
- **Punkt 3 CHANGELOG-Schnitt:** `[Unreleased]` leer, `## [0.44.0] — 2026-07-17`
  darunter. Die Aussagen sind **wahr** und am Image geprüft: „advisory (Exit 0)",
  „gatet allein über `--require-complete`", „`--print-config` führt den neuen
  Block" (`config_template.go:180-197` — verifiziert), die Fail-closed-Liste
  inklusive Vakuum, und „Eine einseitig leere Vorwärts-Sicht … ist **kein**
  Fehler" (Bootstrap-Lauf reproduziert: 1 Differenz, Exit 0).
- **Punkt 5 (`make ci` lokal):** `make gates` grün abgefahren (siehe oben);
  `make ci` = gates + image-test.
- **Release-Prep-Commit-Grenze:** `c4439e8` ist **ein** Commit und **kein**
  Slice-Commit; er fasst ausschließlich `version.md`, CHANGELOG, beide READMEs,
  Handbuch, `operations.md` und den Handbuch-Beispiel-Anker an. Der
  Feat-Commit `1bfa7f8` fasst README/Handbuch **nicht** an. Saubere Trennung.
- **Doku-Wahrheit Handbuch §4.12 (Kernprüfung):** Das dokumentierte
  YAML-Beispiel wurde 1:1 in ein Fixture gegossen und gegen `d-check:latest`
  gefahren. Der dokumentierte Ausgabeblock stimmt **zeichengenau** — beide
  Befundzeilen, beide Richtungslabels, beide `Datei:Zeile` (`:7`/`:7`), die
  Zeile `2 Differenz(en).`, die stderr-Zeile und Exit **1**. Auch die Prosa
  hält: die `GG-SPEC-042`-Mittelschicht-Kante fehlt in der Ausgabe (Ventil), die
  RTM darüber bleibt unverändert.
- **Doku-Wahrheit Handbuch §5:** Die dort behauptete Fail-closed-Liste ist
  vollständig gegen `applyCrossConsistency`/`crossConsistency` abgeglichen;
  **keine** verschwiegene Fail-closed-Klasse gefunden. Die Abgrenzung „**Nicht**
  fail-closed ist eine einseitig leere Vorwärts-Sicht" ist am Image belegt
  (Exit 0, 1 Differenz laut gemeldet). Der `ranges`-Default (true) und der
  `artifact-id-column`-Default (`first`) stimmen mit `crossRanges`/
  `applyCrossBackward` überein.
- **Vertrags-Konsistenz Lastenheft 0.44.2 ↔ Spezifikation ↔ ADR-0038 ↔ Code:**
  Die Schritt-Nummerierung ist nach dem Einschub der Vakuitäts-Stufe **überall**
  korrekt nachgezogen (Schritt 5 = Vakuität, 6 = Diff, 7 = Gatung; der
  Rückverweis in Schritt 3 zeigt jetzt auf Schritt 6, nicht mehr auf 5). Die
  Fehlerpräzedenz der Spezifikation entspricht Zeile für Zeile der Phasenfolge
  in `crossConsistency:73-114` — inklusive der Zusicherung „jede Stufe läuft über
  beide Sichten": `readCrossFile` ×2 → `spanCrossSource` ×2 → `bind*Tables` ×2 →
  `*Edges` ×2 → `excludeReqs` ×2 → `crossVacuity` → `diffViews`. `mode: superset`,
  `TraceCrossArtifactFirst`, Sortierschlüssel `(R, Artefakt, Richtung)` und die
  „erster Treffer gewinnt"-Determinismus-Regel decken sich mit dem Vertrag.
- **`SelectSections`-Refactor (Regressionsrisiko für `matrix`/`trace.coverage`):**
  `SectionMask` ist zur Altfassung logisch äquivalent
  (`(len(include)==0 || inRanges(inc,no)) && !inRanges(exc,no)` ≡ die alten
  beiden `continue`-Zweige); `SelectSections` konsumiert sie. Kein
  Verhaltens-Delta.
- **`markdownTables`-Refactor (Regressionsrisiko für `DC-FA-REQ-001`):** Die
  Schleifen-Fortschaltung (`i = next`), die Relevanz-/`foundTable`-/
  `usedTextHeaders`-Reihenfolge, die Bad-Row-Semantik (Fehler nur bei
  **relevanter** Tabelle, sonst Tabellenende) und die Fehlerpräzedenz
  (Header → Duplicate-ID, der v0.43.1-Fix) sind erhalten. `mask == nil` auf dem
  Requirements-Pfad ⇒ identisches Verhalten.
- **SemVer v0.44.0 (Minor) — kein verborgener Breaking Change:** Der Diff ist
  rein additiv hinter einem opt-in Block. `--require-complete` gatet die neue
  Ursache **nur** bei präsentem `trace.cross-consistency`; `requireCompleteExit`
  hat die Waisen-Ursache unverändert gelassen (kein `return 1` mehr, aber
  `exit = 1` mit identischem Ergebnis) und meldet beide Ursachen getrennt statt
  eine zu verschlucken — für Bestandskonsumenten ohne Block ist das ein No-op.
  Belegt durch `TestCLI071_Cross_DefaultByteIdentisch` (Markdown + Exit) und
  `TestBuildTraceMatrixWithoutCrossBlock`; `omitempty` hält die JSON/YAML-
  Byte-Identität ohne Block. **Minor ist richtig.**
- **DoD-Ehrlichkeit:** Kein Grün, das nicht trägt. Alle acht `[x]`-Punkte
  einzeln gegengeprüft — insbesondere „Tests" (21 Vakuitäts-/Diff-/Config-/
  CLI-Tests, jede genannte Klasse hat einen Test) und „`--print-config`
  sichtbar" (Template-Block verifiziert). Der **nicht** abgehakte
  Realdatenbeleg ist **korrekt** offen markiert und **ehrlich** begründet: der
  Teilstand benennt präzise, was reproduziert ist (die Drift-*Gestalt* als
  CLI-Akzeptanztest, `crossRepo` bildet `{COMP-CORE, COMP-DOMAIN}` vs.
  `{P-005, P-009, COMP-SCHED}` + Mittelschicht-Kante 1:1 ab) und was nicht (der
  Lauf gegen das echte Repo), samt Grund (Konsumenten-Vorarbeit). Die Punkte
  „Nutzerdoku"/„Release"/„Qualität" stehen offen, obwohl die Nutzerdoku faktisch
  steht — Untertreibung, keine Lüge.
- **`AGENTS.md` §3.1/§3.2/§3.6:** kein Host-Go (alles über `make`), keine
  `//nolint`-Direktive im Diff, keine Gate-Lockerung (Coverage-Schwelle 93 %
  unverändert, 93,90 % erreicht).
- **`AGENTS.md` §3.3 (Lifecycle-Moves):** `29daac6` und `6c4ccf5` sind
  `similarity index 100%`-Renames; der Roadmap-Flip ist per `MR-013`-Ausnahme
  gebündelt, der Slice-Body folgt im Feat-Commit. Korrekt.
- **`AGENTS.md` §3.4/§3.5:** ADR-0038 ist `Proposed` ⇒ inhaltliche Nachträge
  erlaubt (`ea0f343` argumentiert das explizit); `make gates` (inkl. `matrix`)
  meldet keine Abwärts-Referenz aus einem Spec-Stratum.
- **`AGENTS.md` §5:** Alle neun Commit-Betreffe tragen eine `DC-*`/`slice-*`-ID;
  ADR-Index geführt; CHANGELOG gepflegt.
- **Reste/Hygiene:** keine toten Pfade (`coverageIDs` sauber zu `rangeAwareIDs`
  umbenannt und beide Aufrufer gezogen; `tableHeaderAt`/`consumeTableRows` in
  `markdownTables` weiterverwendet), keine Debug-Reste, keine Fixture-Dateien
  im Repo (alle Fixtures via `t.TempDir()`), keine `.gitignore`-Verstöße,
  Arbeitsbaum sauber, `gofmt`/`lint` still.
- **`DC-QA-03` (Seiteneffektfreiheit/Netz):** Der Kreuzverweis-Pfad liest
  ausschließlich über `driven.Filesystem`; kein Schreibzugriff, kein
  Netzaufruf; alle Repro-Läufe mit `--network none` und `:ro`.
- **`DC-QA-02` (Determinismus):** `crossView.add` (erste Fundstelle gewinnt),
  `requireCrossFields` (sortierte Schlüssel statt Map-Iteration) und
  `diffViews`' dreistufiger Sortierschlüssel schließen die drei
  Map-Iterations-Nichtdeterminismen, die dieser Diff einführen könnte.
- **[ADR-0005](../plan/adr/0005-modul-layout-hexagon-ordner.md)-Importregeln:**
  `trace_cross.go` importiert nur `model`, `rules`, `port/driven` — kein
  Adapter-Import im Kern; `make arch-check` (a-check) bestätigt es.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 1 | F-1 |
| LOW | 5 | F-2, F-3, F-4, F-5, F-6 |
| INFO | 1 | F-7 |

## Verdikt

**ACCEPT-WITH-NITS** — der Tag **v0.44.0 darf gesetzt werden**.

Ich habe diesen Slice mit dem Auftrag gelesen, ihn zu stoppen, und den Kern
nicht zum Kippen gebracht. Die `releasing.md`-§Release-Prep-Checkliste ist
Punkt für Punkt abgearbeitet — inklusive der drei Stellen, die kein Gate
erfasst und die laut eigener Erfahrung genau deshalb driften: die
Anker-Wanderung in `version.md`, das bare-Tag-Beispiel in Handbuch §2 und die
§11-Zeile an chronologisch richtiger Position. Ich habe gezielt nach
vergessenen Pins gesucht; es gibt keinen. Die Handbuch-§4.12-Ausgabe habe ich
nicht geglaubt, sondern nachgefahren: sie stimmt zeichengenau gegen das
gebaute Image, samt Exit-Code und stderr-Zeile. Die Vakuitäts-Semantik, um die
vier Runden gerungen wurde, hält in allen drei Grenzfällen (beide Sichten leer,
`superset` mit leerer Rück-Sicht, Bootstrap-Zustand). SemVer ist richtig: der
Diff ist additiv hinter einem opt-in Block, und die neue `--require-complete`-
Ursache kann einen Bestandskonsumenten ohne Block nicht erreichen.

**Warum F-1 (MEDIUM) den Tag trotzdem nicht blockiert** (begründete Abweichung
von der Regel „MEDIUM blockiert typischerweise"): Der **dokumentierte** Gate-
Weg ist der Exit-Code — `--require-complete`, `doc-complete`,
`make completeness-check` —, und der ist formatunabhängig korrekt
(`requireCompleteExit` läuft hinter dem Reporter, unabhängig von
`--json`/`--yaml`). Das Feld `crossConsistency` ist heute **nirgends**
zugesagt (weder Handbuch, `operations.md`, Spezifikation noch CHANGELOG nennen
es — siehe F-2), es kann also niemanden brechen, dem etwas versprochen wurde.
Und die Behebung ist keine Hotfix-Zeile, sondern ein Vertrags-Entscheid: „leeres
Array statt weggelassen" berührt `DC-QA-02` (Byte-Identität ohne Block) und
gehört als Mit-Änderung zu `DC-FA-CLI-009` in eine eigene CR — nicht in einen
Release-Prep-Commit eine Stunde vor dem Tag. **Empfehlung:** F-1 und F-2
gemeinsam als Folge-CR schneiden (Feld dokumentieren **und** die
„Beleg-statt-Schweigen"-Invariante auf die maschinenlesbaren Renderer ziehen,
mit Test auf die Nullmengen-Form) — und das Feld nicht dokumentieren, bevor die
Invariante steht.

**Vor dem `git mv` nach `done/`** (nicht vor dem Tag) abzuräumen: F-4
(Slice-§2/§Bezug auf 0.44.2 nachziehen — das Slice-Dokument ist das letzte
Artefakt auf Vor-R1-Stand), F-5 (**jetzt** sortieren; nach `Accepted` friert
`adr-check` die Zeilenfolge ein) und die F-7-Entscheidung (Closure-Trigger vs.
ehrlich offener Realdatenbeleg; der `Accepted`-Schwenk von ADR-0038 ist die
Klasse, die bei ADR-0037 liegenblieb). F-3 und F-6 sind Einzeiler und passen in
denselben Durchgang.
