# Slice slice-136: §3.4 gegen den Kanon halten — Doppelung, Verschärfung oder keins von beidem?

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-84-durchsetzung](../welle-84-durchsetzung.md).

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §3.4; Baseline-Regelwerk
[`modul-03-spec.md` §Ziel-Form: Architektur-Sicht](../../../../.harness/baseline/v5.11.0/regelwerk/modul-03-spec.md)
und
[`grundlagen-referenz-richtung.md`](../../../../.harness/baseline/v5.11.0/regelwerk/grundlagen-referenz-richtung.md);
die vendorte `architecture.template.md`; der Zensus in
[slice-132](../done/slice-132-hard-rule-zensus.md).

**Berührte Spec-Stellen:** — (Harness-Regeltext; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

Der Zensus hat §3.4 als *teilgedeckt* eingestuft und die ungedeckte Hälfte —
die **Sprachfreiheit der Architektur-Sicht** — mit *Auflösungs-Trigger:
permanent* versehen, also für grundsätzlich unmechanisierbar erklärt. Auf
Auftraggeber-Nachfrage stellte sich heraus: **das ist zu grob.** Zwei Fragen
hängen daran, und keine war beantwortet.

**Erstens: Ist §3.4 eine Doppelung?** Seine zweite Aussage — *kein Spec-Stratum
referenziert ADRs, Wellen, Slices, Commit-Hashes oder Closure-Daten* — steht
im Kanon **wörtlich**, in `modul-03` und in `grundlagen-referenz-richtung.md`.

**Zweitens: Ist §3.4 eine Verschärfung?** Wir schreiben *keine
Sprach-/Modul-Pfade*; `modul-03` schreibt, die Sicht *„referenziert
Modul-Pfade"*. Meint das dasselbe? Wenn ja, verschärfen wir — und
[`MR-031`](../../../../harness/conventions.md#mr-031) verlangt dafür einen
Eintrag im Konventionsspeicher, den es nicht gibt.

## 2. Vorgehen

1. **Die Doppelungs-Frage entscheiden** — an `modul-09`, nicht am Gefühl: eine
   Hard Rule soll in **zwei** Quadranten liegen, und der Feedforward-Quadrant
   *ist* der `AGENTS.md`-Eintrag. Prüfen, ob daraus folgt, dass eine gerankte
   Wiederholung des Kanons hier **erwünscht** statt verboten ist.
2. **Die Verschärfungs-Frage messen, nicht auslegen.** Der Bullet-Text ist
   mehrdeutig; die vendorte `architecture.template.md` ist es nicht — sie
   **praktiziert** eine Lesart. Zählen, welche Pfade sie führt und welche
   Kennungen ihre Komponenten-Tabelle trägt.
3. **Je nach Ergebnis:** `MR-`Eintrag (bei Verschärfung), Kürzung (bei
   unnötiger Doppelung) — oder eine berichtigte Ausweisung, wenn beides
   entfällt.
4. **Die Ausweisung berichtigen.** *Permanent* gilt allenfalls für die
   Technologie-Hälfte; die Pfad-Hälfte ist ein **detektierbarer Zustand**
   (heute gemessen: `spec/architecture.md` null, `spec/spezifikation.md` fünf).
   Was davon baubar ist und womit, gehört benannt statt pauschal verneint.
5. **CR-Kandidat festhalten**, falls die Mehrdeutigkeit im Kanon bleibt: sie
   trifft jedes Adopter-Repo, nicht nur uns.
6. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Sensor.** Selbst wenn die Pfad-Hälfte baubar ist — sie zu bauen ist
  ein eigener Slice mit eigener Ortswahl.
- **Kein Konsumenten-CR.** Ein CR-**Kandidat** wird benannt; ob und wann er
  geschrieben wird, ist eine Auftraggeber-Entscheidung.
- **Keine Änderung an den Spec-Straten selbst.** Der Slice prüft eine Regel,
  nicht den Bestand, den sie regelt.

## 4. Definition of Done

- [ ] Beide Fragen sind **mit Beleg** beantwortet — Doppelung ja/nein,
      Verschärfung ja/nein —, und der Beleg ist ein Zitat oder eine Zählung,
      keine Auslegung.
- [ ] Die Folge ist gezogen: `MR-`Eintrag, Kürzung oder berichtigte Ausweisung.
- [ ] §3.4 trägt keine Aussage mehr, die eine Unmöglichkeit behauptet, ohne sie
      geprüft zu haben.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Die Versuchung ist, die bequeme Lesart zu wählen.** *„Modul-Pfade meint
  Dokument-Pfade"* macht unsere Regel konform, *„Modul-Pfade meint
  Code-Pfade"* macht sie zur unerklärten Verschärfung. Wer die erste wählt,
  weil sie weniger Arbeit macht, hat `BEO-012` begangen. — **Ausgang:** *(bei
  Closure)*
- **Eine Vorlage ist kein Regeltext.** Was `architecture.template.md`
  praktiziert, ist starkes Indiz und keine Definition. Die Grenze dieses
  Belegs gehört benannt. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-135](../done/slice-135-uses-pin-sensor.md)
in `done/`, WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Messung eine echte
Verschärfung zeigt und die Frage, ob wir sie behalten wollen, eine
Auftraggeber-Entscheidung ist.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Regeltext (GF), Konventionsspeicher (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23):
  [`BEO-012`](../observations.md) ist **zentral** — dieser Slice besteht ganz
  aus der Frage, wie weit ein Kanon-Satz reicht, und genau daran ist die Welle
  schon fünfmal gescheitert. [`BEO-011`](../observations.md) für jede Aussage
  über den Bestand der Spec-Straten.

Slice-ID: slice-136. Betroffene IDs: — (Harness-Regeltext; keine Anforderung).
Module: Harness-Regeltext, Konventionsspeicher. Gates: `make doc-check`,
`make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Prüfung eigener Regeln gegen die Baseline.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
