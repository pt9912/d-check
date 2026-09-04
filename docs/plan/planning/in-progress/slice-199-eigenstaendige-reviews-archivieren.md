# Slice slice-199: Eigenständige Reviews archivieren

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-90](../welle-90-eigenstaendige-review-archivierung.md) —
zweiter von zwei Slices: wendet den in slice-198 gebauten Modus auf den
echten Bestand an.

**Bezug:** slice-198 (Werkzeug-Modus, Voraussetzung).

**Berührte Spec-Stellen:** — (Bestandspflege, keine neue Anforderung).

**Verantwortlich:** claude-sonnet-5.

**Autor:** pt9912. **Datum:** 2026-09-04.

---

## 1. Ziel

Jeder eigenständige Review-Report (kein `slice-<NNN>` im Dateinamen), der
nach welle-89 noch in `docs/reviews/` liegt, wird per
`tools/archive-wave -review=<dateiname> -apply` archiviert.

**Voraussetzung:** slice-198 in `done/` (der Werkzeug-Modus muss stehen).

## 2. Vorgehen

1. **Bestand exakt erheben** — jede Datei in `docs/reviews/` ohne
   `slice-<NNN>` im Namen (Zahl bei Ausführung erhoben; 11 zum Zeitpunkt der
   Planung).
2. **Dry-Run über den ganzen Bestand vor der Anwendung** — dieselbe Lehre
   wie aus welle-89: ein Voll-Dry-Run deckt Fehlerpfade (insbesondere
   Titel-Extraktion) ab, bevor irgendetwas geschrieben wird.
3. **Je Review ein Lauf** — `tools/archive-wave -review=<dateiname> -apply`.
4. **`make gates` und `make fullbuild`** auf dem Endstand (nicht zwingend
   nach jedem Einzel-Lauf, siehe welle-89s Lerneintrag: ein Voll-Dry-Run vor
   der Anwendung plus eine Prüfung auf dem Endstand ist gleichwertig
   sicherer als 11 Einzel-Prüfungen bei einer uniformen, werkzeug-
   verifizierten Operation).
5. **Commit-Granularität**: ein gebündelter Commit für alle 11 (dieselbe
   Begründung wie bei slice-197 — mechanisch, werkzeug-verifiziert, keine
   inhaltliche Unterscheidung zwischen den Läufen).

## 3. Ausdrücklich NICHT in diesem Slice

- **Reviews mit `slice-<NNN>` im Dateinamen** — unberührt, gehören in den
  `-slice`-Modus.
- **Eine rückwirkende Umbenennung** der archivierten Reviews.

## 4. Definition of Done

- [x] Bestand exakt erhoben: **11** eigenständige Reviews.
- [x] Alle erhobenen Reviews archiviert (`docs/reviews/archiv/<basisname>-archiv.zip`
      + Stub je Review).
- [x] `docs/reviews/` enthält (außerhalb `docs/reviews/archiv/`) <!-- d-check:ignore (Ziel-Form, entsteht erst mit slice-198) --> keinen
      eigenständigen Review mehr.
- [x] `make gates` grün (zehn Gates) auf dem Endstand.
- [x] `make fullbuild` grün (`--require-complete` insbesondere:
      `Hervorgegangen:`-Felder verhindern neue Trace-Waisen).
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/`.
- [x] Unabhängige Verifikation durchgeführt.
- [x] Closure-Notiz (§9) geschrieben, jedes Risiko aus §5 mit Ausgang.

## 5. Abnahme-Punkte / Risiken

- **Dieselbe Klasse „stille Bedeutungsverschiebung" wie bei slice-197.**
  **Entfallen** — Voll-Dry-Run vor der Anwendung lief sauber, Stub-Stichproben
  (unabhängiger Review + Verifikation) bestätigen unveränderten Titel und
  byte-identischen Volltext im Zip für alle geprüften Fälle.
- **Ein Review könnte die einzige Zitierstelle einer `DC-*`/`ADR-*`-Anforderung
  sein.** **Entfallen** — `make fullbuild --require-complete` meldet 0
  Trace-Waisen bei unveränderter Anforderungszahl (51). Präzisierung aus dem
  unabhängigen Review (F-4, LOW): `Hervorgegangen:` fängt nicht *jede* im
  Original zitierte Kennung strukturell auf, sondern nur die außerhalb von
  Fenced-/Inline-Code-Spannen bzw. als Link-Label stehenden — der
  tatsächliche Schutz in diesem Slice ist der empirische
  `--require-complete`-Lauf, nicht die Feld-Konstruktion allein.
- **Der Umfang (11 Reviews) ist klein genug, um NICHT dieselbe
  Ein-Sitzungs-Grenze wie bei slice-195/197 zu sprengen.** **Entfallen** —
  alle 11 Archivierungen liefen in einem gebündelten Commit, keine
  Rückführung nach `next` nötig.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei, slice-198 in `done/`.

**Rückführungen:** `in-progress` → `next`, falls §5s dritter Punkt eintritt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `docs/plan/planning/` (Planning-Harness selbst) bzw.
  `docs/reviews/` — beide fallen unter den Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration).

  <!-- d-check:cite .harness/baseline/v6.0.0/regelwerk/modul-05-planning-harness.md:223-224 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Stand 2026-09-04):
  [`BEO-ALL/review-collection-misses-non-slice-filenames`](../observations/BEO-ALL/review-collection-misses-non-slice-filenames/observation.md)
  (Zähler 1) — dieser Slice löst sie ein;
  [`BEO-ALL/large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md)
  (Zähler 2) gilt der Größenordnung nach hier voraussichtlich nicht (11 statt
  45+ Läufe), wird aber am Ende gegengeprüft.

  <!-- d-check:cite .harness/baseline/v6.0.0/regelwerk/modul-05-planning-harness.md:229-229 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)) — bei
  Beanspruchung neu zu lesen. **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-199. Betroffene IDs: —. Module: —.
Gates: `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Planning-Harness fällt unter den
Default: Doc führt, Code folgt. Reine Bestandspflege mit einem bereits
gebauten und getesteten Werkzeug (slice-198), kein neuer Produktcode.

## 9. Closure-Notiz (nach `done/`)

Alle 11 zum Zeitpunkt der Planung erhobenen eigenständigen Reviews (kein
`slice-<NNN>` im Dateinamen) sind archiviert: Voll-Dry-Run vor der Anwendung
lief sauber, dann 11 `-review=<datei> -apply`-Läufe, ein gebündelter Commit
(`6c4146f`), `make gates` (zehn Module) und `make fullbuild`
(`--require-complete`, 0 Trace-Waisen bei unveränderten 51 Anforderungen)
beide grün auf dem Endstand.

**Nach Review/Verifikation behoben** (unabhängig übereinstimmend gefunden):

- **F-1 (HIGH):** Der `-review=<datei>`-Modus hatte keine der fünf in
  `AGENTS.md` §3.3 enumerierten Ein-Commit-Ausnahmen — derselbe Lücken-Typ,
  der bei slice-197 bereits
  [`MR-062`](../../../../harness/conventions.md#mr-062--wellenloser-einzel-slice-archiv-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)
  nötig machte. Behoben durch
  [`MR-063`](../../../../harness/conventions.md#mr-063--eigenständiger-review-archiv-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)
  plus sechste Ausnahme-Bullet in `AGENTS.md` §3.3 (Fix-Commit `d0639fd`).
- **F-2 (MEDIUM):** Vier durch die Archivierung verwaiste `ignore-refs`-
  Einträge in `.d-check.yml` (deren Quelldateien jetzt Stubs ohne die
  ignorierten Verweise sind) entfernt.
- **F-3 (MEDIUM):** Die immutable
  [ADR-0013](../../adr/0013-pr-ci-und-traceability-gate.md) zitierte einen der
  archivierten Reviews als Klartext-Pfad (kein Modul scannt das). Ein
  Geschichte-Anhang vermerkt die Verlegung, ohne die bestehende Zeile zu
  ändern (§3.5-konform).
- **F-4 (LOW):** Plan-Formulierung zu `Hervorgegangen:` präzisiert (siehe §5).

**Beobachtungs-Register:**
[`BEO-ALL/review-collection-misses-non-slice-filenames`](../observations/BEO-ALL/review-collection-misses-non-slice-filenames/observation.md)
auf **gestrichen** gesetzt — der neue `-review`-Modus deckt genau die Lücke,
die die Beobachtung benannte, unabhängig vom 3×-Zähler.
[`BEO-ALL/large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md)
bleibt bei Zähler 2 — die angekündigte Gegenprüfung bestätigt: der Umfang
(11 Läufe, ein Commit) sprengte die Ein-Sitzungs-Grenze nicht (§5 dritter
Punkt).

**Lerneintrag:** Eine neue `tools/archive-wave`-Betriebsart braucht ihre
Ein-Commit-Rechtfertigung als **eigene** `MR-0XX` von Anfang an, nicht erst
nach einem Review-Fund — dieselbe Lehre, die slice-197s Korrektur-Commit
(`9980aa9` →
[`MR-062`](../../../../harness/conventions.md#mr-062--wellenloser-einzel-slice-archiv-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013))
schon einmal hinterließ, hier aber noch nicht auf den nächsten neuen Modus
übertragen wurde. Geschärfte Regel:
[`MR-063`](../../../../harness/conventions.md#mr-063--eigenständiger-review-archiv-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)
selbst sowie die Bündelungs-Klärung in `AGENTS.md` §3.3.
