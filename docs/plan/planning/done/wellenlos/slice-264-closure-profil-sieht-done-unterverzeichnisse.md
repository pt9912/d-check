# slice-264: Das Closure-Profil prüft die Slices in den Unterverzeichnissen von `done/`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** geteilt aus slice-263, der die Produkt-Schlüssel liefert; Befund
F-1 (HIGH) aus dem Review von slice-260;
[ADR-0048](../../../adr/0048-closure-note-struktur-im-planning-modul.md).

**Berührte Spec-Stellen:** [`SPEC-095`](../../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)
(der Satz zum ausgelösten Lauf).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** `make verify-closure-notes` prüft jeden geschlossenen Slice im
Volltext, gleich ob er direkt unter `done/`, unter `done/wellenlos/` oder unter
einem Wellen-Verzeichnis liegt, und lässt die Stubs aus — mit den Schlüsseln
aus slice-263 in `.d-check.closure.yml`. Die vier Altverstöße, die beim
Schnitt von slice-263 gemessen wurden (Risiko-Ausgang außerhalb des
Wortschatzes in slice-240 bis slice-243), sind behoben, und die Aussagen, dass
der Lauf die Unterverzeichnisse nicht prüft, sind zurückgenommen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Produkt-Schlüssel** — slice-263.
- **Die Archivierung der 24 Volltexte unter `done/wellenlos/`** — ein eigener
  Vorgang am Bestand, kein Teil der Prüfung.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [x] `.d-check.closure.yml` trifft die Volltexte in den Unterverzeichnissen
      und nimmt die Stubs an ihrem Inhalt aus; bewusst gebrochen: ein
      kaputter Volltext unter `done/wellenlos/` macht den Lauf rot, ein Stub
      nicht.
- [x] Die Altverstöße sind behoben; `make verify-closure-notes` und
      `make gates` grün.
- [x] Zurückgenommen: der Satz in [`SPEC-095`](../../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge), Grenze 4 und Vertrag von
      `harness/sensors/hooks.md`, Vertrag „direkt", Grenze 8 und Bindung von
      `harness/sensors/verify-closure-notes.md`, der GRENZE-Kommentar in
      `.githooks/pre-commit`.
- [x] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [x] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.
- [x] [ADR-0105](../../../adr/0105-reviews-liest-done-unterverzeichnisse.md) `Accepted`, [ADR-0081](../../../adr/0081-reviews-modul.md) als
      teil-superseded markiert, Index nachgezogen (Verifikation V-2).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.d-check.closure.yml` | update | Rekursion, Globs, Stub-Ausnahme |
| `done/wellenlos/slice-240` bis `slice-243` | update | Altverstöße |
| `spec/spezifikation.md` ([`SPEC-095`](../../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)), `harness/sensors/{hooks,verify-closure-notes}.md`, `.githooks/pre-commit` | update | Rücknahme der Grenz-Aussagen |
| `.d-check.yml` (`reviews`), `harness/sensors/review-coverage.md` | update | **Plan-Änderung nach Review R1 F-1:** `make review-coverage` liest dieselbe Kandidatenmenge wie der Übergangs-Wächter — sonst antworten zwei Läufe verschieden auf dieselbe Frage; gemessen mit Rekursion: 0 Befunde |
| `observations/BEO-ALL/` (neuer Eintrag) | create | **Plan-Änderung nach Review R1 F-2:** der offene Punkt aus slice-242 R2 bekommt den Ausgang *weiter offen* und braucht dafür einen Register-Eintrag |
| `docs/plan/adr/0105-…` (neu), [ADR-0081](../../../adr/0081-reviews-modul.md) `## Geschichte`, ADR-Index | create/update | **Plan-Änderung nach Review R2 F-9:** die Rekursion in beiden Profilen widerspricht [ADR-0081](../../../adr/0081-reviews-modul.md) Entscheidung 4 — Folge-ADR nach `AGENTS.md` §3.6, `Proposed` bis zur Closure |

## 4. Trigger

**Start** (`next` → `in-progress`): slice-263 in `done/`; `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `open` (blockiert): die Prüfung der Volltexte zeigt mehr
  Altverstöße, als ein Slice beheben sollte — dann die Behebung in einen
  eigenen Slice, mit befristeter Ausnahme.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- **Weitere Altverstöße** — die Messung beim Schnitt galt nur dem Modul
  `structure`; die `planning`-Hälfte (Platzhalter, Floskeln, dünne Notiz) kann
  weitere zeigen. — **Ausgang:** entfallen — die `planning`-Hälfte meldet über
  die 24 Volltexte nichts, gemessen beim Beanspruchen (§8) und in der
  Verifikation (M6).

## 7. Closure-Notiz

- **Was hat funktioniert:** `make verify-closure-notes` prüft die Volltexte in
  `done/wellenlos/` und den Wellen-Verzeichnissen und nimmt die 254 Stubs am
  Marker `> **ARCHIVIERT** — Volltext:` aus. Gemessen vor dem Schnitt:
  genau vier Altverstöße (Risiko-Ausgang `*(offen)*` in slice-240 bis
  slice-243), keine in der `planning`-Hälfte. Die Verifikation brach jede
  Hälfte im Klon: das alte Profil bleibt auf einem kaputten Volltext unter
  `wellenlos/` grün, das neue wird rot — für `planning`, `structure` und
  `reviews`; ein Stub mit Marker bleibt still.
- **Was ging anders als geplant:** Der Slice wuchs zweimal nach Review, beide
  Male als Plan-Änderung vor dem Code. R1 fand, dass `make review-coverage`
  im Hauptprofil nicht rekursiv blieb und auf dieselbe Frage anders antwortete
  als der Übergangs-Wächter; R2, dass die Rekursion der akzeptierten
  [ADR-0081](../../../adr/0081-reviews-modul.md) Entscheidung 4 widersprach —
  daraus [ADR-0105](../../../adr/0105-reviews-liest-done-unterverzeichnisse.md),
  `Proposed` bis zu dieser Closure. Die nachgetragenen Ausgänge der vier
  Altverstöße waren im ersten Anlauf aus dem Schweigen von §7 abgeleitet und
  trugen bei eingetretenen Risiken das Wort *entfallen* (R1 F-2, F-3); jetzt
  belegt jeder Ausgang seine Quelle, und der offene Punkt aus slice-242 steht
  im Register. Meine Zählung der Volltexte verlor slice-263, weil dessen Text
  das Marker-Wort im Fließtext trägt — gezählt wird gegen die volle
  Marker-Form (24).
- **Steering-Loop-Eintrag:** keiner mit neuer Schwelle. Die Klasse
  `BEO-ALL/semantic-change-body-only-edges-stale` ist als
  [`MR-025`](../../../../../harness/conventions.md#mr-025) verkörpert und trat
  erneut auf: die Rekursion wurde im Closure-Profil gesetzt, ihre Ränder im
  Hauptprofil, in der Sensor-Doku und in der ADR blieben stehen.
- **Beobachtungs-Register (`../../observations/`):** `evidence/slice-264.md` in
  [`BEO-ALL/semantic-change-body-only-edges-stale`](../../observations/BEO-ALL/semantic-change-body-only-edges-stale/state.md);
  neu
  [`BEO-ALL/aufnahme-kriterium-ohne-bestands-beleg`](../../observations/BEO-ALL/aufnahme-kriterium-ohne-bestands-beleg/state.md)
  (1×, Beleg slice-242) und
  [`BEO-ALL/eingetreten-ohne-folge-kennung`](../../observations/BEO-ALL/eingetreten-ohne-folge-kennung/state.md)
  (1×, Verifikation V-1).
- **Folge-Slices:** keiner. Produkt-Verhalten unverändert. Die Zeile im
  Benutzerhandbuch, `reviews` lese nicht rekursiv, beschreibt den
  Produkt-Default und wird in der nächsten Release-Prep präzisiert (R1 F-8).
- **Risiken aus §6:** entfallen (siehe §6). Trigger-Audit: kein Carveout,
  kein bootstrap-aware Gate; [ADR-0105](../../../adr/0105-reviews-liest-done-unterverzeichnisse.md)
  neu und `Accepted`, ihre Trigger nicht eingetreten; keine Hard Rule mit
  eingetretenem Trigger. Nachtlauf-Stand
  ([`MR-053`](../../../../../harness/conventions.md#mr-053)): wie in §8.
- **Drei Paarungen:** (a) Anker — kein Eintrag mit `liegt in`; (b)
  Folge-Slices — keiner genannt; (c) Register — die drei zitierten
  Beobachtungen existieren und tragen Belege.

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind die Closure-Konfiguration,
vier `done/`-Slices und die Harness-Doku samt Spezifikation §7 — alle unter dem
Default `*` (`ALL`). Eine eigene Sub-Area erfüllt keine davon: Keine trägt eine
eigene Konvention, einen eigenen Modus oder eine eigene Inventur-Linie.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Zwei berühren den Slice.
[`BEO-ALL/messmethode-scope-enger-als-dokumentierte-ziel-form`](../../observations/BEO-ALL/messmethode-scope-enger-als-dokumentierte-ziel-form/observation.md)
(1×) beschreibt genau die Lage, die dieser Slice schließt — ein Scan-Bereich,
der enger ist als das, wofür sein Grün gelesen wird; der Slice ist ihre
Behebung, kein zweites Auftreten.
[`BEO-ALL/leeres-verzeichnis-lokal-verdeckt-ci-befund`](../../observations/BEO-ALL/leeres-verzeichnis-lokal-verdeckt-ci-befund/observation.md)
(1×) betrifft das bewusste Brechen: Es wird auch im frischen Klon gemessen.
Keine erreicht mit diesem Slice 3×.

**Messung beim Beanspruchen:** Mit Rekursion, `**/`-Globs und der
Stub-Ausnahme `(?m)^> \*\*ARCHIVIERT` in einer Wegwerf-Kopie des Profils
meldet der Lauf über 974 Dateien genau **vier** Befunde — den Risiko-Ausgang
`*(offen)*` in §6 von slice-240 bis slice-243. Die `planning`-Hälfte meldet
nichts, die `reviews`-Hälfte ebenfalls nichts. Unter `done/wellenlos/` liegen
24 Volltexte (gezählt gegen die volle Marker-Form); die Wellen-Verzeichnisse tragen nur Stubs.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../../harness/conventions.md#mr-053)):
`upstream-drift` grün (2026-10-09T07:17Z). `image-scan` rot (2026-10-09T10:37Z)
— der Lauf liegt vor dem Release v0.85.0, das die gemeldeten CVEs behebt;
der Slice berührt das Image nicht.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
