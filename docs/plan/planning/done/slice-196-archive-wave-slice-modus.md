# Slice slice-196: `tools/archive-wave` — wellenloser Einzel-Slice-Archivierungs-Modus

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-89](../welle-89-wellenlose-review-archivierung.md) — erster
von zwei Slices: baut den Werkzeug-Modus, den slice-197 auf den Bestand
anwendet.

**Bezug:** Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht, Zeile „Zeitdokumente archivieren … Ohne Wellen tut es die
Slice-Closure selbst … Schlüssel ist der Slice: `done/slice-<NNN>-archiv.zip`,
flach neben dem Stub".

**Berührte Spec-Stellen:** — (Tooling unter `tools/`, kein Produktcode-Modul
von d-check selbst).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-04.

---

## 1. Ziel

`tools/archive-wave` bekommt einen zweiten Betriebsmodus neben dem
bestehenden Wellen-Modus (`-welle=<id>`): `-slice=<slice-id>` archiviert
**einen einzelnen** `done/`-Slice ohne Wellen-Zugehörigkeit — seinen
Volltext und seine Review-Reports nach `done/slice-<NNN>-archiv.zip`
(flach, keine Unterverzeichnis-Ebene wie beim Wellen-Modus), gekürzter Stub
an seiner Stelle, Verweis-Nachzug wie im bestehenden Modus. Die beiden
Modi sind exklusiv (`-welle` und `-slice` gleichzeitig ⇒ Exit 2).

## 2. Vorgehen

1. **CLI-Flag `-slice` ergänzen** (`main.go`), mutually exclusive zu
   `-welle`, gleiche `-apply`/`-root`-Semantik.
2. **Sammel-Logik für den Ein-Slice-Fall** (`collect.go`): der Slice selbst
   (muss in `done/` liegen und `ohne Welle` sein — ein Slice mit
   `**Welle:**`-Feld gehört in den Wellen-Modus, nicht hierher) plus seine
   Review-Reports (dieselbe Erkennung wie im Wellen-Modus, `slice-<NNN>` im
   Dateinamen).
3. **Archiv-Pfad ohne Wellen-Ebene** (`archive.go`): `done/slice-<NNN>-archiv.zip`
   statt `done/<welle-id>/archiv.zip` — flach neben dem Stub, wie der Kanon
   es fordert.
4. **Stub-Form für den wellenlosen Fall** (`stub.go`): das vorhandene
   `archiv-stub-slice.template.md` sieht `**Welle:** ohne Welle` bereits vor,
   aber sein `**Archiviert mit:** <welle-id>`-Feld setzt eine Welle voraus.
   Für den Selbst-Archivierungs-Fall trägt das Feld stattdessen
   `**Archiviert:** <JJJJ-MM-TT> (eigene Closure)` — kein Widerspruch zum
   Template (das Feld beschreibt die *Einsammlung*, hier eben durch die
   eigene Closure statt durch eine Welle), aber eine neue Wert-Form, die im
   Slice-Plan dokumentiert wird, weil das Template sie nicht wörtlich
   vorwegnimmt.
5. **Tests**: Fixture mit einem wellenlosen `done/`-Slice + zwei
   Review-Reports, Vorher/Nachher-Vergleich (Stub-Form, Archiv-Inhalt,
   Verweis-Nachzug), Negativtest (`-welle` und `-slice` gleichzeitig ⇒
   Exit 2), Negativtest (Slice mit `**Welle:**`-Feld über `-slice`
   angesprochen ⇒ Fehlermeldung statt stillem Fehlverhalten).

## 3. Ausdrücklich NICHT in diesem Slice

- **Die Anwendung auf den echten Bestand** — das ist slice-197.
- **Eine Änderung des Wellen-Modus** (`-welle`) — bleibt unverändert, beide
  Modi teilen sich nur die Hilfsfunktionen, die schon modusneutral sind.
- **Eine automatische Erkennung „wellenlos, also -slice"** ohne explizites
  Flag — der Aufrufer sagt, welchen Modus er will; das Werkzeug rät nicht.

## 4. Definition of Done

- [x] `-slice=<slice-id>`-Modus implementiert, mutually exclusive zu
      `-welle` (Exit 2 bei beiden gesetzt oder keinem gesetzt) —
      `validateModeFlags` in `main.go`, separat testbar von `main()`.
- [x] Archiv liegt flach unter `done/slice-<NNN>-archiv.zip`, Stub trägt die
      neue `**Archiviert:** … (eigene Closure)`-Form — gegen einen echten
      Bestands-Slice (`slice-137`) im Dry-Run smoke-getestet.
- [x] Ein Slice mit gesetztem, echtem `**Welle:**`-Feld wird über `-slice`
      abgelehnt (Fehlermeldung, kein stilles Falsch-Archivieren).
- [x] Unit-Tests inkl. Fixture-Vorher/Nachher-Vergleich und beiden
      Negativtests, plus ein dritter, im Vorgehen nicht benannter Test: eine
      Meldung (nicht Behebung) toter externer Verweise auf die zu löschenden
      Review-Reports, da für sie kein Move-Ziel existiert (`FindReferencesToPaths`).
- [x] `make archive-wave-test` grün.
- [x] `make gates` grün (zehn Gates).
- [x] `make fullbuild` grün.
- [x] Unabhängiger Review durchgeführt
      ([Report](../../../reviews/2026-09-04-slice-196-archive-wave-slice-modus-code-r1.md),
      1 HIGH + 2 MEDIUM — alle drei behoben, siehe §9).
- [x] Unabhängige Verifikation durchgeführt
      ([Report](../../../reviews/2026-09-04-slice-196-archive-wave-slice-modus-verifikation.md),
      bestanden, 1 LOW ohne Handlungsbedarf).
- [x] Closure-Notiz (§9) geschrieben, jedes Risiko aus §5 mit Ausgang.

## 5. Abnahme-Punkte / Risiken

- **Das Feld, das im Wellen-Modus `**Archiviert mit:**` heißt, bekommt für
  den wellenlosen Fall einen eigenen Namen (`**Archiviert:**`) und eine
  Bedeutung, die das vendorte Template nicht wörtlich vorwegnimmt —
  entfallen.** Umgesetzt genau wie in §2 Punkt 4 begründet:
  `SliceStubStandalone` schreibt `**Archiviert:** <Datum> (eigene
  Closure)`, das vendorte Template selbst bleibt unverändert (committet
  vendored) — die Adaption lebt allein im Werkzeug-Code, nicht im Template.
- **Ein Slice, dessen Review-Report-Namen nicht dem `slice-<NNN>`-Muster
  folgen, wird von der Sammel-Logik nicht gefunden — entfallen als neues
  Risiko, benannt als Bestandsgrenze.** Dieselbe Grenze wie im bestehenden
  Wellen-Modus (unverändert), hier nicht neu geschaffen; ob sie beim
  Bestand aus slice-197 real zutrifft, zeigt sich erst dort.
- **Zeitgleiche Arbeit an einem anderen Slice könnte einen neuen,
  wellenlosen `done/`-Slice anlegen — entfallen.** WIP-Limit 1 hielt; kein
  zweiter Slice lief parallel.

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

- **Offene Beobachtungen sichten** (Stand 2026-09-04): keine Treffer zu
  `archive-wave` oder wellenloser Archivierung im Register.

  <!-- d-check:cite .harness/baseline/v6.0.0/regelwerk/modul-05-planning-harness.md:229-229 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)) — bei
  Beanspruchung neu zu lesen. **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-196. Betroffene IDs: —. Module: —.
Gates: `make gates`, `make fullbuild`, `make archive-wave-test`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — `tools/archive-wave/` fällt unter den
Default: Doc führt, Code folgt. Erweiterung eines bereits etablierten,
getesteten Werkzeugs um einen zweiten Modus, keine neue Konvention.

## 9. Closure-Notiz (nach `done/`)

**Geliefert:** `tools/archive-wave` bekam einen zweiten Modus (`-slice=<id>`),
mutually exclusive zu `-welle` (`validateModeFlags`, separat testbar von
`main()`). Neue Funktionen: `FindSlice` (collect.go), `ApplySlice`
(archive.go), `SliceStubStandalone` (stub.go), `runSlice` (main.go),
`FindReferencesToPaths` (rewrite.go, meldet tote Verweise auf gelöschte
Review-Reports, da für sie kein Move-Ziel existiert). Beide Makefiles
(`tools/archive-wave/Makefile`, Root-`Makefile`) nehmen `SLICE=` als
Alternative zu `WELLE=` an. `AGENTS.md`s `make archive-wave`-Zeile
nachgezogen.

**Was anders lief als geplant:** nichts Wesentliches — der Umfang deckte
sich mit §2. Ein Fixture-Bug in der eigenen Testdatei (falsche
Pfad-Tiefe im externen Verweis, `../../` statt `../../../`) führte zu
einem zunächst falsch-negativen Test (`TestRunSlice_DanglingReviewReference`)
— sofort selbst gefunden und korrigiert, kein Produktcode-Fehler.

**Nach Review/Verifikation behoben (1 HIGH + 2 MEDIUM, unabhängig
gefunden):**

- **Fünf neue Kommentare trugen Slice-Nummer-Provenienz** (`slice-196` als
  Herkunfts-Prosa in `main.go`, `stub.go`, `slice_mode_test.go`) — §3.7
  verbietet das ausdrücklich. Behoben: Slice-Bezug entfernt, nur die
  funktionale Beschreibung bleibt stehen. Unterschied zu den unveränderten
  Alt-Kommentaren desselben Werkzeugs (`gemessen an welle-70`, `slice-075`):
  jene zitieren empirischen Bestandsbefund, die neuen erzählten die eigene
  Autorschaft — genau die vom Kanon ausgeschlossene Klasse.
- **`FindSlice`s beide Fehlerpfade (0 Treffer, mehrdeutig) waren ungetestet**
  — anders als das strukturell analoge `FindWellePlan`. Zwei Tests
  nachgetragen (`TestFindSlice_Keine`, `TestFindSlice_Mehrdeutig`).
- **Die exclude-Menge von `FindReferencesToPaths` war nie geprüft** — ein
  Regressions-Test nachgetragen
  (`TestFindReferencesToPaths_ExcludesSelfAndReviews`), der Vorgänger
  (`BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt`) verlangt die
  Umkehr-Probe: mit auskommentiertem Exclude-Check schlägt der neue Test
  fehl, mit dem echten Code besteht er.

**Lerneintrag:** keiner über das Behobene hinaus — reine Werkzeug-Erweiterung
nach etabliertem Muster (Wellen-Modus als Vorlage), keine neue Fehlerklasse.

**Risiko-Ausgänge:** siehe §5 — alle drei entfallen.
