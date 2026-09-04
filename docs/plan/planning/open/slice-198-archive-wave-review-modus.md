# Slice slice-198: `tools/archive-wave` — Modus für eigenständige Review-Archivierung

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-90](../welle-90-eigenstaendige-review-archivierung.md) —
erster von zwei Slices: baut den Werkzeug-Modus, den slice-199 auf den
Bestand anwendet.

**Bezug:** [`BEO-ALL/review-collection-misses-non-slice-filenames`](../observations/BEO-ALL/review-collection-misses-non-slice-filenames/observation.md).

**Berührte Spec-Stellen:** — (Tooling unter `tools/`, kein Produktcode-Modul
von d-check selbst).

**Verantwortlich:** —.

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

- [ ] `-review=<dateiname>`-Modus implementiert, mutually exclusive zu
      `-welle` **und** `-slice` (Exit 2 bei mehr oder weniger als einem
      gesetzten Flag).
- [ ] Archiv liegt unter `docs/reviews/archiv/<basisname>-archiv.zip`, Stub
      im selben Verzeichnis mit Titel, Archiv-Zeiger, Datum,
      `Hervorgegangen:`-Feld.
- [ ] Ein Review mit `slice-<NNN>` im Dateinamen wird über `-review`
      abgelehnt (Fehlermeldung, kein stilles Falsch-Archivieren).
- [ ] `ExtractTitle` funktioniert nachweislich für mindestens zwei der elf
      realen, uneinheitlichen Überschriftenformen (Fixture-Kopie).
- [ ] Unit-Tests inkl. Fixture-Vorher/Nachher-Vergleich und beiden
      Negativtests.
- [ ] `make archive-wave-test` grün.
- [ ] `make gates` grün (zehn Gates).
- [ ] `make fullbuild` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/`.
- [ ] Unabhängige Verifikation durchgeführt.
- [ ] Closure-Notiz (§9) geschrieben, jedes Risiko aus §5 mit Ausgang.

## 5. Abnahme-Punkte / Risiken

- **`ExtractTitle`s Regex ist auf Slice-/Wellen-Header zugeschnitten
  (`^#\s+(?:Slice\s+|Welle\s+)?[\w-]+:?\s*(.*)$`) und wurde nie gegen
  uneinheitliche Review-Überschriften geprüft.** Gegenmaßnahme: expliziter
  Test gegen reale Formen (siehe §2 Punkt 4) vor der Anwendung auf den
  Bestand — dieselbe Lehre wie bei slice-196/197 (Fixture-vs-Bestand-Lücke
  erst am echten Bestand sichtbar).
- **Ein Review, dessen einzige Zitierstelle einer Anforderung er ist, könnte
  die Anforderung zur Trace-Waise machen**, wenn `Hervorgegangen:` fehlt
  oder falsch extrahiert. Gegenmaßnahme: `make fullbuild`
  (`--require-complete`) nach der Anwendung in slice-199, nicht nur `make
  gates`.
- **Zeitgleiche Arbeit an einem anderen Slice könnte einen neuen,
  eigenständigen Review anlegen**, während slice-199 später den Bestand
  aufnimmt — WIP-Limit 1 schützt dagegen.

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

<!-- wird erst bei Closure gefüllt -->
