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

- [ ] Bestand exakt erhoben (Zahl im DoD eingetragen).
- [ ] Alle erhobenen Reviews archiviert (`docs/reviews/archiv/<basisname>-archiv.zip`
      + Stub je Review).
- [ ] `docs/reviews/` enthält (außerhalb `docs/reviews/archiv/`) <!-- d-check:ignore (Ziel-Form, entsteht erst mit slice-198) --> keinen
      eigenständigen Review mehr.
- [ ] `make gates` grün (zehn Gates) auf dem Endstand.
- [ ] `make fullbuild` grün (`--require-complete` insbesondere:
      `Hervorgegangen:`-Felder verhindern neue Trace-Waisen).
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/`.
- [ ] Unabhängige Verifikation durchgeführt.
- [ ] Closure-Notiz (§9) geschrieben, jedes Risiko aus §5 mit Ausgang.

## 5. Abnahme-Punkte / Risiken

- **Dieselbe Klasse „stille Bedeutungsverschiebung" wie bei slice-197.**
  Gegenmaßnahme: Voll-Dry-Run vor der Anwendung (§2 Punkt 2), `make
  fullbuild` auf dem Endstand statt 11 Einzel-Prüfungen.
- **Ein Review könnte die einzige Zitierstelle einer `DC-*`/`ADR-*`-Anforderung
  sein** — `Hervorgegangen:` fängt das strukturell auf, `--require-complete`
  bestätigt es empirisch.
- **Der Umfang (11 Reviews) ist klein genug, um NICHT dieselbe
  Ein-Sitzungs-Grenze wie bei slice-195/197 zu sprengen** — falls doch,
  Rückführung nach `next`.

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

<!-- wird erst bei Closure gefüllt -->
