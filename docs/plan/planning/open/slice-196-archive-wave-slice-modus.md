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

**Verantwortlich:** —.

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

- [ ] `-slice=<slice-id>`-Modus implementiert, mutually exclusive zu
      `-welle` (Exit 2 bei beiden gesetzt oder keinem gesetzt).
- [ ] Archiv liegt flach unter `done/slice-<NNN>-archiv.zip`, Stub trägt die
      neue `**Archiviert:** … (eigene Closure)`-Form.
- [ ] Ein Slice mit gesetztem `**Welle:**`-Feld wird über `-slice`
      abgelehnt (Fehlermeldung, kein stilles Falsch-Archivieren).
- [ ] Unit-Tests inkl. Fixture-Vorher/Nachher-Vergleich und beiden
      Negativtests.
- [ ] `make archive-wave-test` grün.
- [ ] `make gates` grün (zehn Gates).
- [ ] `make fullbuild` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/`.
- [ ] Unabhängige Verifikation durchgeführt.
- [ ] Closure-Notiz (§9) geschrieben, jedes Risiko aus §5 mit Ausgang.

## 5. Abnahme-Punkte / Risiken

- **Das Stub-Feld `**Archiviert mit:**` bekommt eine Bedeutung, die das
  vendorte Template nicht wörtlich trägt** (Wert ist keine Welle-Kennung,
  sondern ein Datum + Erläuterung). Gegenmaßnahme: im Slice-Plan explizit
  begründet (siehe §2 Punkt 4); keine Änderung am vendorten Template selbst
  (das bleibt committet vendored).
- **Ein Slice, dessen Review-Report-Namen nicht dem `slice-<NNN>`-Muster
  folgen, wird von der Sammel-Logik nicht gefunden** — dieselbe Grenze wie
  im bestehenden Wellen-Modus, hier nicht neu, aber erstmals gegen einzelne
  Slices statt gegen einen ganzen Wellen-Bestand gemessen.
- **Zeitgleiche Arbeit an einem anderen Slice könnte einen neuen,
  wellenlosen `done/`-Slice anlegen**, während slice-197 später den
  Bestand aufnimmt — WIP-Limit 1 schützt dagegen.

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

<!-- wird erst bei Closure gefüllt -->
