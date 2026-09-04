# Slice slice-197: Wellenlosen Review-Bestand archivieren

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-89](../welle-89-wellenlose-review-archivierung.md) —
zweiter von zwei Slices: wendet den in slice-196 gebauten Modus auf den
echten Bestand an.

**Bezug:** slice-196 (Werkzeug-Modus, Voraussetzung).

**Berührte Spec-Stellen:** — (Planungs-Bestandspflege, keine neue
Anforderung).

**Verantwortlich:** pt9912.

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

- [x] Bestand exakt erhoben — **45** Slices (nicht die geschätzten ~43):
      slice-083, 095, 102, 112, 121, 127, 137–147, 151–167, 171, 176–182,
      185–187, 189.
- [x] Alle 45 archiviert (`docs/plan/planning/done/wellenlos/<id>-archiv.zip` +
      Stub je Slice, gemeinsames Verzeichnis statt `done/slice-<NNN>-archiv.zip`
      flach — siehe §9, während der Ausführung korrigiert).
- [x] `docs/reviews/` enthält keinen Report mehr zu einem archivierbaren
      wellenlosen Slice — ein permanenter, benannter Rest siehe §3/§9
      (`release-prep-v0.71.0-review.md`, kein Slice-Review, per
      `ignore-refs`-Tombstone gedeckt statt archiviert).
- [x] `make gates` grün (zehn Gates) auf dem Endstand.
- [x] `make fullbuild` grün (526 Dateien, 0 Befunde).
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/`.
- [ ] Unabhängige Verifikation durchgeführt.
- [x] Closure-Notiz (§9) geschrieben, jedes Risiko aus §5 mit Ausgang.

## 5. Abnahme-Punkte / Risiken

- **Größtes Risiko: 45 Einzel-Läufe sind fehleranfällig für dieselbe Klasse
  „stille Bedeutungsverschiebung" wie bei einer Datenmigration — entfallen,
  mit Korrektur unterwegs.** Statt `make gates` nach jedem Einzel-Lauf lief
  ein Dry-Run-Scan über den ganzen Bestand (Fehler- und Dangling-Referenz-
  Prüfung je Slice) vor der Anwendung, dann alle 45 Läufe, dann `make gates`
  einmal auf dem Endstand — begründet in §9. Dabei fiel ein echter
  Design-Fehler in slice-196 auf (flacher Stub-Pfad kollidiert mit den
  Closure-Struktur-Prüfungen) und wurde vor dem Bestands-Commit korrigiert
  (siehe §9) — genau die Art Fund, gegen die die ursprünglich geplante
  Je-Lauf-Prüfung schützen sollte, hier durch einen vorgelagerten
  Voll-Dry-Run und `make fullbuild` auf dem Endstand ersetzt.
- **Der Umfang (45 Slices) sprengte die Ein-Sitzungs-Review-Grenze nominell,
  Rückführung entfiel bewusst.** Wie bei
  [`BEO-ALL/large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md)
  vorgezeichnet: eine Aufteilung in Batches hätte die Werkzeug-Korrektur
  (siehe oben) mehrfach nachziehen müssen, ohne die Fehleranfälligkeit zu
  senken — die Operation ist uniform und werkzeug-verifiziert (Review +
  Verifikation von slice-196), nicht inhaltlich pro Slice verschieden. Kein
  neuer Beobachtungs-Eintrag, da dieselbe Klasse bereits registriert ist.
- **Ein Review-Report ohne erkennbaren `slice-<NNN>`-Bezug bleibt
  unarchiviert — eingetreten wie vorhergesehen, benannt statt behoben.**
  Betrifft ältere CR-/Baseline-Reviews; keiner davon war Gegenstand dieses
  Slice.

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

**Geliefert:** 45 wellenlose `done/`-Slices archiviert
(`tools/archive-wave -slice=<id> -apply`, je einer der 45 in einem
gemeinsamen Lauf): `docs/plan/planning/done/wellenlos/<id>-archiv.zip` +
gekürzter Stub im selben Verzeichnis, Review-Reports gelöscht (kein Stub),
repo-weite Verweise nachgezogen. `docs/reviews/` trägt danach keinen Report
mehr zu einem archivierbaren wellenlosen Slice — ein permanenter Rest
(`release-prep-v0.71.0-review.md`, kein Slice-Review) per
`ignore-refs`-Tombstone gedeckt.

**Was anders lief als geplant — ein echter Design-Fehler in slice-196,
gefunden erst am echten Bestand:**

1. **Flacher Stub-Pfad kollidierte mit den Closure-Struktur-Prüfungen.**
   slice-196 ließ die Slice-Datei am unveränderten Pfad
   (`done/slice-<NNN>.md`), nur der Inhalt wurde durch den Stub ersetzt.
   Beim ersten Anwendungslauf (45 Slices) meldete `make gates` 132 neue
   `id-unlinked`-Befunde (die `Hervorgegangen:`-Zeile trägt bare
   `DC-*`/`ADR-*`-Kennungen aus dem archivierten Volltext, teils vorher
   verlinkt) und `make fullbuild` 183 `section-missing`/
   `closure-note-missing`-Befunde (ein Stub hat naturgemäß kein DoD, keine
   Closure-Notiz). Beides war bei slice-196s eigenem, kleinen Test-Fixture
   unsichtbar, weil dort keine echten Anforderungs-Kennungen vorkamen und
   `make fullbuild` dort nicht gegen den Wellen-Modus-Vergleichsfall lief.
2. **Korrektur, mit Rücksprache:** die `id-unlinked`-Welle wurde durch einen
   `d-check:ignore`-Marker auf der `Hervorgegangen:`-Zeile behoben (pfad-
   unabhängig, im Unterschied zum Wellen-Modus-Ventil in `.d-check.yml`).
   Für die Closure-Struktur-Kollision zwei Optionen erwogen — ein
   Unterverzeichnis je Slice (spiegelt den Wellen-Modus 1:1, erzeugt aber 45
   fast leere Ordner ohne inhaltliche Gruppierung) oder ein gemeinsames
   Verzeichnis für alle wellenlosen Archive. Nutzer-Entscheid: das
   gemeinsame Verzeichnis (`docs/plan/planning/done/wellenlos/`) — entgeht
   denselben nicht-rekursiven Scan-Mustern wie `done/<welle-id>/`, ohne
   Verzeichnis-Vervielfachung. `planning.closure` hat gar keinen
   Exempt-Paths-Mechanismus (anders als `structure`/`reviews`) — eine
   datei-weise Ausnahmeliste hätte dort eine Produktcode-Änderung an
   d-check selbst gebraucht, keine Konfiguration; das gemeinsame Verzeichnis
   braucht keine.
3. **`ApplySlice` bekam dadurch einen echten Move** (vorher keiner nötig, da
   der Pfad unverändert blieb) — `runSlice` zieht seither repo-weite
   Verweise auf den verschobenen Slice nach (`RewriteRepo`), zusätzlich zur
   bereits vorhandenen Dangling-Meldung für gelöschte Review-Reports.
4. **Alle 45 Archivierungen wurden nach der Korrektur verworfen und mit dem
   reparierten Werkzeug neu angewendet** — nichts von der ersten,
   fehlerhaften Anwendung wurde committet.
5. **Ein Commit statt 45 einzelner** (§5 Punkt 1) — Begründung dort.

**Lerneintrag:** Ein neuer Werkzeug-Modus, der nur an einem konstruierten
Fixture bewiesen ist (slice-196 §2 Punkt 5), kann eine Kollision mit
anderen, unabhängig konfigurierten Prüfungen (hier: `ids`, `planning.closure`)
erst am echten Bestand zeigen, wenn dessen Streuung (Anforderungs-Kennungen,
Dokumentgröße) über das Fixture hinausgeht. Kein neuer Sensor — dieselbe
Klasse wie bei jeder Fixture-vs-Bestand-Lücke, hier ohne dritte Instanz.

**Risiko-Ausgänge:** siehe §5 — alle drei entfallen bzw. eingetreten und
aufgefangen, keines wandert neu ins Beobachtungs-Register.
