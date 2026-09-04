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

- [ ] Bestand exakt erhoben (Zahl im DoD eingetragen).
- [ ] Alle erhobenen Slices archiviert (`done/wellenlos/<basisname>.md`
      Stub + `done/slice-<NNN>-archiv.zip` je Slice).
- [ ] `docs/plan/planning/done/` enthält (außerhalb `done/wellenlos/`) keinen
      flachen wellenlosen Slice mehr.
- [ ] `make gates` grün (zehn Gates) auf dem Endstand.
- [ ] `make fullbuild` grün (`--require-complete`: 0 neue Trace-Waisen).
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/`.
- [ ] Unabhängige Verifikation durchgeführt.
- [ ] Closure-Notiz (§9) geschrieben, jedes Risiko aus §5 mit Ausgang.

## 5. Abnahme-Punkte / Risiken

- **Dieselbe Klasse „stille Bedeutungsverschiebung" wie bei slice-197/199.**
  Gegenmaßnahme: Voll-Dry-Run vor der Anwendung (§2 Punkt 2), `make
  fullbuild` auf dem Endstand statt sieben Einzel-Prüfungen.
- **Ein Slice könnte die einzige Zitierstelle einer `DC-*`/`ADR-*`-Anforderung
  sein** — `Hervorgegangen:` fängt das teilweise strukturell auf (siehe
  slice-199 §5, F-4: nicht jede im Original zitierte Kennung, nur die
  außerhalb von Fenced-/Inline-Code-Spannen), `--require-complete`
  bestätigt es empirisch.
- **Der Umfang (7 Slices) ist klein genug, um NICHT dieselbe
  Ein-Sitzungs-Grenze wie bei slice-195/197 zu sprengen** — falls doch,
  Rückführung nach `next`.

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

<!-- wird erst bei Closure gefüllt -->
