# Slice slice-198: `tools/archive-wave` — Modus für eigenständige Review-Archivierung

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-90](../welle-90-eigenstaendige-review-archivierung.md) —
erster von zwei Slices: baut den Werkzeug-Modus, den slice-199 auf den
Bestand anwendet.

**Bezug:** [`BEO-ALL/review-collection-misses-non-slice-filenames`](../observations/BEO-ALL/review-collection-misses-non-slice-filenames/observation.md).

**Berührte Spec-Stellen:** — (Tooling unter `tools/`, kein Produktcode-Modul
von d-check selbst).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-04.

---

## 1. Ziel

`tools/archive-wave` bekommt einen dritten Betriebsmodus neben `-welle` und
`-slice`: `-review=<dateiname>` archiviert einen **eigenständigen**
Review-Report — einen ohne `slice-<NNN>` im Dateinamen, der keinem Slice
zugeordnet ist — nach `docs/reviews/archiv/<basisname>-archiv.zip`,
gekürzter Stub im selben Verzeichnis. Alle drei Modi bleiben mutually
exclusive.

**Design-Unterschied zu den beiden bestehenden Modi:** ein Slice- oder
Wellen-Review bekommt beim Archivieren **keinen** Stub (seine Identität kommt
vom Slice/von der Welle). Ein eigenständiger Review **ist** selbst der
abgeschlossene Vorgang und bekommt deshalb, wie ein Slice, einen eigenen
Stub — sonst verschwände er spurlos, ohne Zeiger auf sein Archiv.

## 2. Vorgehen

1. **CLI-Flag `-review` ergänzen** (`main.go`), mutually exclusive zu
   `-welle` **und** `-slice` (`validateModeFlags` auf drei Flags erweitert:
   genau eines von dreien).
2. **Zugehörigkeits-Prüfung**: die genannte Datei muss unter
   `docs/reviews/` existieren und darf **kein** `slice-<NNN>`-Muster im
   Dateinamen tragen (`sliceIDInNameRE`) — sonst gehört sie in den
   `-slice`-Modus (dessen Sammel-Logik sie ohnehin schon findet, sobald ihr
   Slice archiviert wird).
3. **Archiv + Stub** (`archive.go`, `stub.go`): `docs/reviews/archiv/` <!-- d-check:ignore (Ziel-Form, entsteht erst mit diesem Slice) -->
   (neue Konstante analog zu `WellenlosArchiveDir`) nimmt Zip und Stub auf.
   Stub-Titel via `ExtractTitle` (funktioniert bereits generisch für einen
   einfachen `# ...`-Header, keine Slice/Welle-Form nötig — geprüft an allen
   elf realen Überschriften). `Hervorgegangen:`-Feld wie bei Slices
   (`ExtractSurvivingIDs`/`FormatHervorgegangen`), da ein Review die einzige
   Zitierstelle einer Anforderung sein kann.
4. **Tests**: Fixture mit einem eigenständigen Review + einem Review mit
   `slice-<NNN>` im Namen (Gegenprobe: letzterer wird abgelehnt), Vorher/
   Nachher-Vergleich (Stub-Form, Archiv-Inhalt), Negativtest (zwei oder null
   von drei Modus-Flags gesetzt ⇒ Exit 2), ein Test gegen mindestens zwei
   der elf realen Überschriftenformen (Fixture-Kopie, kein Fremd-Import).

## 3. Ausdrücklich NICHT in diesem Slice

- **Die Anwendung auf den echten Bestand** — das ist slice-199.
- **Eine Änderung der beiden bestehenden Modi** (`-welle`, `-slice`).
- **Eine automatische Umleitung** eines fälschlich über `-review`
  angesprochenen Slice-Reviews in den `-slice`-Modus — der Aufrufer sagt,
  welchen Modus er will; das Werkzeug rät nicht (dieselbe Linie wie bei
  slice-196).

## 4. Definition of Done

- [x] `-review=<dateiname>`-Modus implementiert, mutually exclusive zu
      `-welle` **und** `-slice` (`validateModeFlags` auf drei Flags erweitert:
      genau eines von dreien, Exit 2 sonst).
- [x] Archiv liegt unter `docs/reviews/archiv/<basisname>-archiv.zip`, Stub
      im selben Verzeichnis mit Titel, Archiv-Zeiger, Datum,
      `Hervorgegangen:`-Feld.
- [x] Ein Review mit `slice-<NNN>` im Dateinamen wird über `-review`
      abgelehnt (Fehlermeldung, kein stilles Falsch-Archivieren).
- [x] **Präzisiert:** nicht `ExtractTitle` (Slice-/Wellen-Präfixschema),
      sondern eine neue Funktion `ExtractFullHeading` — `ExtractTitle` hätte
      das führende Wort einer Review-Überschrift ("Review-Report:", "Review
      —") verschluckt, weil es dessen Gruppe 1 nie zurückgibt. Gegen drei
      der elf realen Überschriftenformen unit-getestet; der Dry-Run gegen
      alle elf testet nur `FindReview`/den Ablehnungs-Pfad, **nicht**
      `ExtractFullHeading` selbst (kehrt vorher zurück) — die unabhängige
      Verifikation hat die tatsächliche Titel-Extraktion separat gegen alle
      elf realen Überschriften bestätigt (siehe Verifikations-Report).
- [x] Unit-Tests inkl. Fixture-Vorher/Nachher-Vergleich und beiden
      Negativtests.
- [x] `make archive-wave-test` grün.
- [x] `make gates` grün (zehn Gates).
- [x] `make fullbuild` grün.
- [x] Unabhängiger Review durchgeführt
      ([Report](../../../reviews/2026-09-04-slice-198-archive-wave-review-modus-code-r1.md),
      1 HIGH + 2 LOW + 1 INFO — alle behoben, siehe §9).
- [x] Unabhängige Verifikation durchgeführt
      ([Report](../../../reviews/2026-09-04-slice-198-archive-wave-review-modus-verifikation.md),
      bestanden, ein Formulierungs-Widerspruch §5/§9 gemeldet — behoben,
      siehe §9).
- [ ] Closure-Notiz (§9) geschrieben, jedes Risiko aus §5 mit Ausgang.

## 5. Abnahme-Punkte / Risiken

- **`ExtractTitle`s Regex ist auf Slice-/Wellen-Header zugeschnitten und
  wurde nie gegen uneinheitliche Review-Überschriften geprüft —
  entfallen.** `ExtractTitle` wurde für Reviews nie eingesetzt: die
  Prüfung (manuelles Durchspielen der Regex gegen alle elf realen
  Überschriften, siehe §9) zeigte das Verschlucken des führenden Worts
  bereits am Entwurf, vor jeder Implementierung — eine eigene Funktion
  (`ExtractFullHeading`) trat von Anfang an an ihre Stelle. Kein
  Produktionscode hat je die falsche Funktion verwendet.
- **Ein Review, dessen einzige Zitierstelle einer Anforderung er ist, könnte
  die Anforderung zur Trace-Waise machen — entfallen für diesen Slice**
  (kein Bestand berührt), bleibt Prüfpunkt für slice-199 (`--require-complete`
  nach der Anwendung).
- **Zeitgleiche Arbeit an einem anderen Slice könnte einen neuen,
  eigenständigen Review anlegen — entfallen.** WIP-Limit 1 hielt.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** keine vorab erkennbare Bedingung — reiner Tooling-Slice
mit festem Umfang.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `tools/archive-wave/` — kein `tools/harness/`-Pfad,
  fällt unter den Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration).

  <!-- d-check:cite .harness/baseline/v6.0.0/regelwerk/modul-05-planning-harness.md:223-224 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Stand 2026-09-04): direkt einschlägig
  ist
  [`BEO-ALL/review-collection-misses-non-slice-filenames`](../observations/BEO-ALL/review-collection-misses-non-slice-filenames/observation.md)
  (Zähler 1) — dieser Slice behebt sie, statt auf die 3×-Schwelle zu warten.

  <!-- d-check:cite .harness/baseline/v6.0.0/regelwerk/modul-05-planning-harness.md:229-229 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)) — bei
  Beanspruchung neu zu lesen. **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-198. Betroffene IDs: —. Module: —.
Gates: `make gates`, `make fullbuild`, `make archive-wave-test`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — `tools/archive-wave/` fällt unter den
Default: Doc führt, Code folgt. Erweiterung eines bereits etablierten,
getesteten Werkzeugs um einen dritten Modus, keine neue Konvention.

## 9. Closure-Notiz (nach `done/`)

**Geliefert:** `tools/archive-wave` bekam einen dritten Modus
(`-review=<dateiname>`), mutually exclusive zu `-welle`/`-slice`. Neue
Funktionen: `FindReview` (collect.go), `ApplyReview`, `ReviewArchiveDir`
(archive.go), `ReviewStub`, `ExtractFullHeading` (stub.go), `runReview`
(main.go). Beide Makefiles nehmen `REVIEW=` als dritte Alternative zu
`WELLE=`/`SLICE=` an. `AGENTS.md`s `make archive-wave`-Zeile nachgezogen —
dabei auch eine bereits veraltete Passage korrigiert (beschrieb noch den
flachen Vor-welle-89-Pfad des `-slice`-Modus).

**Was anders lief als geplant:** Das im Plan (§5) benannte Risiko wurde
bereits **beim Entwurf** erkannt, bevor `ExtractTitle` je für Reviews
verwendet wurde: die Funktion (für Slice-/Wellen-Header gebaut) hätte bei
einer Review-Überschrift wie "# Review-Report: Change Request …" das Wort
"Review-Report" verschluckt, weil sie nur ihre erste Capture-Gruppe
zurückgibt. Statt die bestehende Regex zu erweitern (und beide Aufrufer
subtil zu koppeln), eigene Funktion `ExtractFullHeading` (liest die ganze
Zeile nach `# `) für Reviews — von Anfang an, kein Produktionscode hat je
`ExtractTitle` für einen Review aufgerufen. Gegen drei der elf realen
Überschriftenformen unit-getestet; der Dry-Run gegen alle elf prüft nur den
Sammel-/Ablehnungs-Pfad, nicht die Extraktion selbst — die unabhängige
Verifikation hat sie separat gegen alle elf realen Überschriften bestätigt.

**Nach Review/Verifikation behoben (1 HIGH + 2 LOW + 1 Formulierungs-Befund,
unabhängig gefunden):**

- **Ein neuer Testkommentar trug erneut Slice-Nummer-Provenienz**
  (`review_mode_test.go`, „belegt slice-198 §2 Punkt 4/§4") — dieselbe
  Fehlerklasse wie bei slice-196 (§3.7), in derselben Sitzung ein zweites
  Mal. Behoben: Slice-Bezug entfernt; die Zahlenangabe „gegen zwei" auf die
  tatsächliche „gegen drei" korrigiert.
- **Ein Kommentar über `validateModeFlags` war seit der Drei-Flag-Erweiterung
  veraltet** (beschrieb noch zwei Flags) — nachgezogen.
- **Die DoD/§9-Formulierung überzeichnete den Dry-Run-Beleg**: der Dry-Run
  gegen alle elf testet nur `FindReview`/den Ablehnungs-Pfad, nicht
  `ExtractFullHeading` selbst (kehrt vorher zurück) — die unabhängige
  Verifikation deckte die tatsächliche Extraktion separat ab. Beide Stellen
  präzisiert.
- **§5 und §9 widersprachen sich wörtlich** über denselben Risiko-Ausgang
  (§5: „eingetreten, aufgefangen statt entfallen"; §9: „entfallen") —
  unabhängige Verifikation. Auf „entfallen" vereinheitlicht: kein
  Produktionscode hat je `ExtractTitle` für einen Review aufgerufen, das
  Risiko wurde vor jeder Implementierung erkannt.

**Lerneintrag:** keiner über das Vermerkte hinaus — die Lehre aus
slice-196/197 (Fixture-vs-Bestand-Lücken vor der Bestandsanwendung
schließen, nicht danach) wurde hier direkt angewendet, kein neuer
Steering-Loop-Eintrag nötig.

**Risiko-Ausgänge:** siehe §5 — alle drei entfallen (das erste durch
rechtzeitiges Erkennen vor der Implementierung, nicht durch einen
Bestandsfehler).
