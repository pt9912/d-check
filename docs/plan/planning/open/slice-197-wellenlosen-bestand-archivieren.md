# Slice slice-197: Wellenlosen Review-Bestand archivieren

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-89](../welle-89-wellenlose-review-archivierung.md) —
zweiter von zwei Slices: wendet den in slice-196 gebauten Modus auf den
echten Bestand an.

**Bezug:** slice-196 (Werkzeug-Modus, Voraussetzung).

**Berührte Spec-Stellen:** — (Planungs-Bestandspflege, keine neue
Anforderung).

**Verantwortlich:** —.

**Autor:** pt9912. **Datum:** 2026-09-04.

---

## 1. Ziel

Jeder wellenlos geschlossene `done/`-Slice, dessen Review-Report(s) noch
unarchiviert in `docs/reviews/` liegen, wird per
`tools/archive-wave -slice=<id> -apply` archiviert: Volltext + Reports nach
`done/slice-<NNN>-archiv.zip`, gekürzter Stub an seiner Stelle.

**Voraussetzung:** slice-196 in `done/` (der Werkzeug-Modus muss stehen).

## 2. Vorgehen

1. **Bestand exakt erheben** — jeder `done/*.md`-Slice ohne `**Welle:**`-Feld
   (oder mit `ohne Welle`), dessen zugehörige Review-Reports noch in
   `docs/reviews/` liegen (Substring-Match `slice-<NNN>` im Dateinamen, wie
   das Werkzeug selbst sammelt). Zahl bei Ausführung erhoben, nicht aus
   diesem Plan übernommen (~43 geschätzt).
2. **Je Slice ein Lauf** — `tools/archive-wave -slice=<id> -apply`,
   `git status` prüfen, `make gates` nach **jedem** Lauf (dieselbe
   Welle-für-Welle-Disziplin wie bei slice-191: jeder Fehler wird am Ort
   seines Auftretens sichtbar, nicht erst am Ende).
3. **Commit-Granularität**: ein Commit je archiviertem Slice (mehrere
   `D`/`A`-Paare, keine Renames — dieselbe Begründung wie bei
   [`MR-059`](../../../../harness/conventions.md#mr-059--wellen-archiv-stub-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013),
   angewendet auf den neuen Ein-Slice-Modus) **oder** ein gebündelter Commit
   für mehrere Slices, falls die Einzel-Commit-Granularität bei ~43
   Läufen unverhältnismäßig wird — Entscheidung bei Ausführung, mit
   Begründung in der Commit-Botschaft.
4. **Abschluss-Prüfung**: `docs/reviews/` enthält keinen Report mehr zu
   einem archivierbaren wellenlosen Slice.

## 3. Ausdrücklich NICHT in diesem Slice

- **Slices, deren Reports bereits über eine Welle archiviert sind** —
  unberührt, sie liegen schon in einem `archiv.zip`.
- **Eine nachträgliche `**Welle:**`-Feld-Vergabe** für tatsächlich
  wellenlose Slices — sie bleiben `ohne Welle`.
- **Reviews ohne erkennbaren Slice-Bezug** (z. B. reine CR-/Baseline-Reviews
  ohne `slice-<NNN>` im Dateinamen) — fallen außerhalb der Sammel-Logik des
  Werkzeugs und bleiben unangetastet; benannt, nicht behoben.

## 4. Definition of Done

- [ ] Bestand exakt erhoben (Zahl im DoD eingetragen).
- [ ] Alle erhobenen Slices archiviert (`done/slice-<NNN>-archiv.zip` +
      Stub je Slice).
- [ ] `docs/reviews/` enthält keinen Report mehr zu einem archivierbaren
      wellenlosen Slice (Reste benannt, falls vorhanden — siehe §3).
- [ ] `make gates` grün (zehn Gates) auf dem Endstand.
- [ ] `make fullbuild` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/`.
- [ ] Unabhängige Verifikation durchgeführt.
- [ ] Closure-Notiz (§9) geschrieben, jedes Risiko aus §5 mit Ausgang.

## 5. Abnahme-Punkte / Risiken

- **Größtes Risiko: ~43 Einzel-Läufe sind fehleranfällig für dieselbe
  Klasse „stille Bedeutungsverschiebung" wie bei einer Datenmigration.**
  Gegenmaßnahme: `make gates` nach jedem Lauf, kein Sammel-Commit ohne
  zwischenzeitliche Prüfung.
- **Der Umfang (~43 Slices) könnte die Ein-Sitzungs-Review-Grenze sprengen**
  — falls so, Rückführung nach `next` und Aufteilung in Batches (z. B. nach
  Slice-Nummernblock), nicht stilles Weiterarbeiten über die Grenze hinaus.
- **Ein Review-Report ohne erkennbaren `slice-<NNN>`-Bezug bleibt
  unarchiviert** (siehe §3) — benannt als bekannte Grenze, kein Blocker für
  die Closure dieses Slice.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei, slice-196 in `done/`.

**Rückführungen:** `in-progress` → `next`, falls §5s zweiter Punkt eintritt
(Umfang sprengt eine Review-Sitzung) — dann Aufteilung in Batches.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `docs/plan/planning/` (Planning-Harness selbst).
  Fällt unter den Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration).

  <!-- d-check:cite .harness/baseline/v6.0.0/regelwerk/modul-05-planning-harness.md:223-224 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Stand 2026-09-04): keine Treffer zu
  wellenloser Archivierung; die frisch aus welle-88 registrierten
  Beobachtungen (`BEO-ALL/large-migration-exceeds-session-review-limit`,
  `BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`) gelten direkt für
  diesen Slice — großer Umfang über viele Einzel-Läufe, dieselbe
  Ein-Sitzungs-Vorsicht.

  <!-- d-check:cite .harness/baseline/v6.0.0/regelwerk/modul-05-planning-harness.md:229-229 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)) — bei
  Beanspruchung neu zu lesen. **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-197. Betroffene IDs: —. Module: —.
Gates: `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Planning-Harness fällt unter den
Default: Doc führt, Code folgt. Reine Bestandspflege mit einem bereits
gebauten und getesteten Werkzeug (slice-196), kein neuer Produktcode.
**Evidenz-/Diskrepanz-Risiko** ist die Achse mit Substanz — siehe §5.

## 9. Closure-Notiz (nach `done/`)

<!-- wird erst bei Closure gefüllt -->
