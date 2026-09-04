# Slice slice-200: Verbliebene wellenlose Slices archivieren

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Reine Bestandspflege mit einem bereits gebauten
und getesteten Werkzeug (`tools/archive-wave -slice=<id>`, slice-196), kein
neuer Produktcode; der Closure-Grund geht über die eigene DoD nicht hinaus
(Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).

**Bezug:** slice-197 (welle-89) — dessen Bestandserhebung war ausdrücklich auf
wellenlose Slices **mit** noch unarchivierten Review-Reports beschränkt (§2
Punkt 1 dort: „… dessen zugehörige Review-Reports noch in `docs/reviews/`
liegen"). Sieben wellenlose Slices **ohne** Review-Report blieben dadurch
flach in `done/` liegen, obwohl sie für den `-slice=<id>`-Modus ebenso
archivierbar sind (per Dry-Run bestätigt).

**Berührte Spec-Stellen:** — (Bestandspflege, keine neue Anforderung).

**Verantwortlich:** claude-sonnet-5.

**Autor:** pt9912. **Datum:** 2026-09-04.

---

## 1. Ziel

Jeder flach in `docs/plan/planning/done/` liegende wellenlose Slice (kein
`**Welle:**`-Feld bzw. `— wellenlos`) wird per
`tools/archive-wave -slice=<id> -apply` nach
`docs/plan/planning/done/wellenlos/` archiviert — unabhängig davon, ob er
einen Review-Report trägt.

## 2. Vorgehen

1. **Bestand exakt erheben** — jeder `done/*.md`-Slice mit `**Welle:**` =
   `wellenlos`/`— wellenlos`, der noch nicht unter
   `docs/plan/planning/done/wellenlos/` liegt (Zahl bei Ausführung erhoben;
   7 zum Zeitpunkt der Planung: slice-141, slice-168, slice-169, slice-170,
   slice-183, slice-184, slice-188).
2. **Dry-Run über den ganzen Bestand vor der Anwendung** — dieselbe Lehre
   wie aus welle-89/slice-199: ein Voll-Dry-Run deckt Fehlerpfade ab, bevor
   irgendetwas geschrieben wird.
3. **Je Slice ein Lauf** — `tools/archive-wave -slice=<id> -apply`.
4. **`make gates` und `make fullbuild`** auf dem Endstand (nicht nach jedem
   Einzel-Lauf — dieselbe Begründung wie bei slice-199: ein Voll-Dry-Run vor
   der Anwendung plus eine Prüfung auf dem Endstand ist gleichwertig
   sicherer als sieben Einzel-Prüfungen bei einer uniformen,
   werkzeug-verifizierten Operation).
5. **Commit-Granularität**: ein gebündelter Commit für alle 7 (dieselbe
   Begründung wie bei slice-197/slice-199 — mechanisch, werkzeug-verifiziert,
   keine inhaltliche Unterscheidung zwischen den Läufen).

## 3. Ausdrücklich NICHT in diesem Slice

- **Eine nachträgliche Korrektur der Closure-Notiz von slice-197** — dessen
  eigene Scope-Abgrenzung (§3 dort) bleibt stehen; dieser Slice schließt die
  fachliche Lücke, ohne den archivierten Vorgänger anzufassen.
- **Eine Änderung des `-slice=<id>`-Modus selbst.**

## 4. Definition of Done

- [x] Bestand exakt erhoben: **7** wellenlose Slices ohne Review-Report.
- [x] Alle erhobenen Slices archiviert (`done/wellenlos/<basisname>.md`
      Stub + `done/wellenlos/slice-<NNN>-archiv.zip` je Slice).
- [x] `docs/plan/planning/done/` enthält (außerhalb `done/wellenlos/`) keinen
      flachen wellenlosen Slice mehr.
- [x] `make gates` grün (zehn Gates) auf dem Endstand.
- [x] `make fullbuild` grün (`--require-complete`: 0 neue Trace-Waisen).
      Bestätigt erst nach dem `git mv` dieses Slices nach `done/` — davor
      traf der nicht-rekursive `done/slice-*.md`-Glob des Closure-Profils auf
      0 Dateien, weil alle anderen flachen wellenlosen Slices bereits
      archiviert waren.
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/`.
- [x] Unabhängige Verifikation durchgeführt.
- [x] Closure-Notiz (§9) geschrieben, jedes Risiko aus §5 mit Ausgang.

## 5. Abnahme-Punkte / Risiken

- **Dieselbe Klasse „stille Bedeutungsverschiebung" wie bei slice-197/199.**
  **Entfallen** — Voll-Dry-Run vor der Anwendung lief sauber; unabhängiger
  Review bestätigt 5/7 Zips byte-identisch, 2/7 (slice-184, slice-188) mit
  einem korrekt nachgezogenen Link statt eines veralteten (siehe §9 — eine
  neue, benannte Nebenwirkung, kein Bedeutungsverlust).
- **Ein Slice könnte die einzige Zitierstelle einer `DC-*`/`ADR-*`-Anforderung
  sein.** **Entfallen** — `make fullbuild --require-complete` bestätigt 0
  neue Trace-Waisen (siehe §9).
- **Der Umfang (7 Slices) ist klein genug, um NICHT dieselbe
  Ein-Sitzungs-Grenze wie bei slice-195/197 zu sprengen.** **Entfallen** —
  alle 7 Archivierungen liefen in einem gebündelten Commit, keine
  Rückführung nach `next` nötig.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls §5s dritter Punkt eintritt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `docs/plan/planning/` (Planning-Harness selbst) fällt
  unter den Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration).

  <!-- d-check:cite .harness/baseline/v6.0.0/regelwerk/modul-05-planning-harness.md:223-224 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Stand 2026-09-04): kein Eintrag im
  Beobachtungs-Register betrifft die Bestandserhebungs-Grenze von
  slice-197 direkt (dort als reine Scope-Abgrenzung in Prosa geführt, ohne
  einen der drei kanonischen Risiko-Ausgänge). Keine Treffer sind ebenfalls
  eine Antwort.

  <!-- d-check:cite .harness/baseline/v6.0.0/regelwerk/modul-05-planning-harness.md:229-229 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)) — bei Beanspruchung
  neu zu lesen. **Dieser Block trägt bewusst keine `cite`-Direktive** — sein
  Ziel ist eine Repo-Adaption
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-200. Betroffene IDs: —. Module: —.
Gates: `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Planning-Harness fällt unter den
Default: Doc führt, Code folgt. Reine Bestandspflege mit einem bereits
gebauten und getesteten Werkzeug (slice-196), kein neuer Produktcode.

## 9. Closure-Notiz (nach `done/`)

Alle 7 zum Zeitpunkt der Planung erhobenen wellenlosen Slices ohne
Review-Report (slice-141, slice-168, slice-169, slice-170, slice-183,
slice-184, slice-188) sind archiviert: Voll-Dry-Run vor der Anwendung lief
sauber, dann 7 `-slice=<id> -apply`-Läufe, ein gebündelter Commit
(`ce5fb50`), `make gates` (zehn Module) grün auf dem Endstand.

**Nach Review/Verifikation behoben bzw. dokumentiert** (unabhängig
übereinstimmend gefunden):

- **F-1 (MEDIUM):** Die Bündelung von 7 Einzel-Slice-Archiven in einem
  Commit folgte der gelebten Praxis von slice-197 (45×), ohne dass
  [`MR-062`](../../../../harness/conventions.md#mr-062--wellenloser-einzel-slice-archiv-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)
  diese Bündelung — anders als sein eigener „Grenze"-Absatz ausdrücklich
  festhält — bereits deckte. Behoben durch
  [`MR-064`](../../../../harness/conventions.md#mr-064--bündelung-mehrerer-einzel-slice-archiv-moves-in-einem-commit-ist-zulässig-nachtrag-zu-mr-062),
  mit derselben Bündelungs-Klausel wie
  [`MR-063`](../../../../harness/conventions.md#mr-063--eigenständiger-review-archiv-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)
  sie für den Review-Modus trägt.
- **F-2 (MEDIUM):** Die Zips von slice-184 und slice-188 enthalten je einen
  bereits auf slice-183s neuen Speicherort nachgezogenen Link, weil
  `ApplySlice()` den aktuellen Platten-Inhalt zippt statt eines Git-Objekts
  — slice-183 wurde vor ihnen archiviert, sein `RewriteRepo()`-Durchlauf
  hatte ihre Verweise bereits umgeschrieben. Inhaltlich harmlos (der Link
  bleibt korrekt auflösbar), aber vom Batch-Verarbeitungsstand abhängig statt
  deterministisch gegen einen festen Commit — als
  [`BEO-ALL/batch-slice-archival-zips-post-rewrite-content`](../observations/BEO-ALL/batch-slice-archival-zips-post-rewrite-content/observation.md)
  registriert, nicht in diesem Slice behoben (Werkzeug-Änderung außerhalb
  des Bestandspflege-Scopes, siehe §3).
- **F-3:** Vom Review vermutete verirrte Duplikat-Datei erwies sich als
  Artefakt eines parallel laufenden Verifikations-Tests (dieser hatte sie
  selbst wieder entfernt) — kein eigener Befund.
- **F-4 (INFO):** Fünf jetzt gegenstandslose `exempt-paths`-Einträge im
  `reviews:`-Block (`.d-check.yml` und `.d-check.closure.yml`) für exakt die
  fünf hier archivierten Slices mit bekannter Review-Lücke entfernt.

**Beobachtungs-Register:** ein neuer Eintrag
([`BEO-ALL/batch-slice-archival-zips-post-rewrite-content`](../observations/BEO-ALL/batch-slice-archival-zips-post-rewrite-content/observation.md))
für F-2, Zähler 1, weiter offen.

**Lerneintrag:** Eine gelebte Commit-Bündelungs-Praxis braucht ihre
Deklaration **beim ersten Mal**, nicht erst beim zweiten oder dritten
Vorkommen — slice-197 lebte sie ungedeckt, slice-200 wiederholte es, bevor
[`MR-064`](../../../../harness/conventions.md#mr-064--bündelung-mehrerer-einzel-slice-archiv-moves-in-einem-commit-ist-zulässig-nachtrag-zu-mr-062)
sie nachträglich schloss. Geschärfte Regel: dieselbe Regel selbst.
