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
      slice-083, 095, 102, 112, 121, 127, 137–140, 142–147, 151–167, 171,
      176–182, 185–187, 189.
- [x] Alle 45 archiviert — **präzisiert**: `docs/plan/planning/done/wellenlos/<id>-archiv.zip`
      + Stub im selben, gemeinsamen Verzeichnis statt eines flachen
      `done/slice-<NNN>-archiv.zip` — während der Ausführung korrigiert,
      siehe §9.
- [x] `docs/reviews/` enthält keinen Report mehr zu einem archivierbaren
      wellenlosen Slice — ein permanenter, benannter Rest siehe §3/§9
      (`release-prep-v0.71.0-review.md`, kein Slice-Review, per
      `ignore-refs`-Tombstone gedeckt statt archiviert).
- [x] `make gates` grün (zehn Gates) auf dem Endstand.
- [x] `make fullbuild` grün (528 Dateien, 0 Befunde).
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/`
      ([Report](../../../reviews/2026-09-04-slice-197-wellenlosen-bestand-archivieren-code-r1.md),
      1 HIGH — behoben, siehe §9).
- [x] Unabhängige Verifikation durchgeführt
      ([Report](../../../reviews/2026-09-04-slice-197-wellenlosen-bestand-archivieren-verifikation.md),
      2 Befunde zur Risiko-Ausgangs-Disziplin plus
      [`MR-059`](../../../../harness/conventions.md#mr-059--wellen-archiv-stub-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)-Zitat-Überdehnung
      — alle behoben, siehe §9).
- [x] Closure-Notiz (§9) geschrieben, jedes Risiko aus §5 mit Ausgang.

## 5. Abnahme-Punkte / Risiken

- **Größtes Risiko: 45 Einzel-Läufe sind fehleranfällig für dieselbe Klasse
  „stille Bedeutungsverschiebung" wie bei einer Datenmigration — eingetreten,
  gefunden und behoben, nicht entfallen.** Zwei echte Werkzeug-Fehler traten
  am echten Bestand auf (mehrzeiliges `**Welle:**`-Feld wurde abgeschnitten;
  die Zugehörigkeits-Prüfung griff danach fälschlich auf Kontext-Erwähnungen
  fremder Wellen) — beide von unabhängigem Review/Verifikation bzw. beim
  eigenen erneuten Anwendungslauf gefunden, siehe §9. Die geplante
  Je-Lauf-`make gates`-Prüfung hätte den ersten Fehler (Text-Kollision mit
  `ids`/`planning.closure`) gefangen, tatsächlich lief ein Voll-Dry-Run vor
  der Anwendung plus `make gates`/`make fullbuild` auf dem Endstand — beide
  Fehler wurden dennoch vor dem endgültigen Commit korrigiert, kein
  fehlerhafter Zustand ist dauerhaft eingegangen.
- **Der Umfang (45 Slices) sprengte die Ein-Sitzungs-Review-Grenze,
  Rückführung entfiel bewusst — weiter offen im Beobachtungs-Register.**
  Eingetragen als
  [`BEO-ALL/large-migration-exceeds-session-review-limit`](../observations/BEO-ALL/large-migration-exceeds-session-review-limit/observation.md)
  (zweites Auftreten nach slice-195, Zähler 1 → 2): eine Aufteilung in
  Batches hätte die zwischenzeitlich nötigen Werkzeug-Korrekturen mehrfach
  nachziehen müssen, ohne die Fehleranfälligkeit zu senken.
- **Ein Review-Report ohne erkennbaren `slice-<NNN>`-Bezug bleibt
  unarchiviert — eingetreten wie vorhergesehen, weiter offen im
  Beobachtungs-Register.** Eingetragen als
  [`BEO-ALL/review-collection-misses-non-slice-filenames`](../observations/BEO-ALL/review-collection-misses-non-slice-filenames/observation.md)
  — eine strukturelle Eigenschaft der Sammel-Logik (beide Modi), kein
  einmaliger Fund; nicht als eigener, erfundener vierter Ausgang belassen
  (unabhängige Verifikation, Bezug
  [`BEO-ALL/invented-fourth-closure-outcome`](../observations/BEO-ALL/invented-fourth-closure-outcome/observation.md)).

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
(`tools/archive-wave -slice=<id> -apply`): je
`docs/plan/planning/done/wellenlos/<id>-archiv.zip` + gekürzter Stub im
selben, gemeinsamen Verzeichnis, Review-Reports gelöscht (kein Stub),
repo-weite Verweise nachgezogen. `docs/reviews/` trägt danach keinen Report
mehr zu einem archivierbaren wellenlosen Slice — ein permanenter Rest
(`release-prep-v0.71.0-review.md`, kein Slice-Review) per
`ignore-refs`-Tombstone gedeckt.

**Was anders lief als geplant — drei echte, am Bestand gefundene
Werkzeug-Fehler in slice-196, alle vor der Closure behoben:**

1. **Flacher Stub-Pfad kollidierte mit den Closure-Struktur-Prüfungen.**
   Der erste Anwendungslauf ließ die Slice-Datei am unveränderten Pfad;
   `make gates`/`make fullbuild` meldeten daraufhin eine `id-unlinked`-Welle
   (bare Kennungen in `Hervorgegangen:`) und zahlreiche
   `section-missing`/`closure-note-missing`-Befunde (ein Stub hat kein DoD,
   keine Closure-Notiz). Nutzer-Entscheid: gemeinsames Verzeichnis
   `docs/plan/planning/done/wellenlos/` statt eines Unterverzeichnisses je
   Slice — entgeht denselben nicht-rekursiven Scan-Mustern wie
   `done/<welle-id>/` es für den Wellen-Modus tut, ohne
   Verzeichnis-Vervielfachung; `planning.closure` hat gar keinen
   Exempt-Paths-Mechanismus, eine datei-weise Ausnahmeliste hätte dort eine
   Produktcode-Änderung an d-check selbst gebraucht. `ApplySlice` liefert
   seither einen echten `Move` für den repo-weiten Verweis-Nachzug; die
   `Hervorgegangen:`-Zeile trägt zusätzlich einen `d-check:ignore`-Marker.
2. **`ReadWelleField` schnitt mehrzeilige Welle-Felder ab.** Die Funktion las
   nur die erste Zeile; die Haus-Stil-Form dieses Repos schreibt das Feld
   für wellenlose Slices meist als mehrsätzigen, umgebrochenen Absatz — 44
   von 45 Stubs im zweiten Anwendungslauf trugen dadurch einen mitten im
   Satz abgeschnittenen, teils sinnentstellten Text (unabhängiger Review, F-1
   HIGH; slice-095s Stub sagte sinngemäß das Gegenteil des Originals).
   Behoben: liest jetzt bis zur ersten Leerzeile.
3. **Die Zugehörigkeits-Prüfung griff danach auf Kontext-Erwähnungen
   fremder Wellen.** Mit dem vollen, mehrzeiligen Feldtext lehnte
   `runSlice` sechs echte wellenlose Slices (083, 151, 152, 176, 177, 178)
   fälschlich ab, weil ihr Absatz eine fremde Welle als Kontext nannte
   ("Anlass liegt in welle-86"). Beim eigenen dritten Anwendungslauf
   gefunden, nicht vom Review. Behoben: die Prüfung liest wie `CollectSlices`
   nur die erste Zeile.
4. **Jeder der drei Fehler löste eine vollständige Neuanwendung aller 45
   Slices aus** — nichts von den beiden vorherigen, fehlerhaften Ständen
   wurde committet oder blieb im Endstand stehen.
5. **Die im Plan (§2 Punkt 3) zitierte
   [`MR-059`](../../../../harness/conventions.md#mr-059--wellen-archiv-stub-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)-Analogie
   für die Commit-Granularität trug nicht** (unabhängige Verifikation):
   jene Regel grenzt ihren Geltungsbereich ausdrücklich auf den Wellen-Modus
   ein. Behoben durch eine eigene, dem Wellen-Fall strukturell gleiche, aber
   im Geltungsbereich getrennte Adaption
   ([`MR-062`](../../../../harness/conventions.md#mr-062--wellenloser-einzel-slice-archiv-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)) —
   anders als die einmalige, erschöpfte
   [`MR-061`](../../../../harness/conventions.md#mr-061--register-formatmigration-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)
   ist dies ein wiederkehrender Vorgang, wie jenes.
6. **Zwei Risiko-Ausgänge trugen zunächst keinen der drei kanonischen
   Formen** (unabhängige Verifikation): die BEO-„weiter offen"-Zuweisung für
   das Ein-Sitzungs-Risiko blieb ohne tatsächliche `evidence/slice-197.md`
   (Zähler fälschlich bei 1 statt 2) und der Review-Sammel-Grenze-Fund war
   als freitextlicher vierter Ausgang formuliert statt als echter
   Register-Eintrag. Beide nachgetragen — siehe §5.
7. **Ein verwaister `ignore-refs`-Eintrag** für den alten, flachen
   `slice-112`-Pfad in `.d-check.yml` wurde entfernt (die referierende Datei
   liegt seit der Archivierung an anderer Stelle mit anderem Inhalt, der
   Eintrag feuerte nie wieder).

**Lerneintrag:** Ein neuer Werkzeug-Modus, der nur an einem konstruierten
Fixture bewiesen ist (slice-196 §2 Punkt 5), kann Kollisionen mit anderen,
unabhängig konfigurierten Prüfungen (`ids`, `planning.closure`) und mit der
tatsächlichen Streuung realer Feld-Formen (mehrzeilige Prosa, Kontext-
Erwähnungen fremder Wellen) erst am echten Bestand zeigen. Keine dritte
Instanz dieser Klasse — kein neuer Sensor, dieselbe bereits gemessene
Fixture-vs-Bestand-Lücke.

**Risiko-Ausgänge:** siehe §5 — Risiko 1 eingetreten und behoben (kein
Register-Eintrag, da die Werkzeug-Korrektur selbst der Ausgang ist), Risiko 2
und 3 eingetreten und ins Beobachtungs-Register überführt (`weiter offen`).
