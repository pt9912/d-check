# Slice slice-145: Ein Sensor auf Pfad-Token in der Architektur-Sicht

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §3.4 (zweiter Auflösungs-Trigger); [`MR-033`](../../../../harness/conventions.md#mr-033); [slice-136](../done/slice-136-agents-34-klaerung.md).

**Berührte Spec-Stellen:** — (abhängig von der Wegwahl; ein Produkt-Delta bräuchte Bump und ADR).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-25.

---

## 1. Ziel

§3.4 verbietet der Architektur-Sicht Modul-Pfade — eine Verschärfung gegenüber
der Baseline, geführt als [`MR-033`](../../../../harness/conventions.md#mr-033).
Die Regel zerfällt in zwei Hälften: *Rollen statt Technologie* ist Urteil,
**Modul-Pfad ja/nein** ist ein detektierbarer Zustand. Gemessen: die Sicht trägt
heute **null** solcher Token.

Kein heutiges Modul trägt die Prüfung — `codepaths` **findet** solche Token,
verbietet sie aber nicht.

## 2. Vorgehen

0. **Die Schwelle steht, bevor gemessen wird** — festgehalten am 2026-08-26,
   vor dem ersten Lauf und vor jedem Blick in die Sicht
   ([`BEO-011`](../observations.md); dieselbe Form wie in
   [slice-144](../done/slice-144-commit-hash-muster.md), wo sie getragen hat):

   - **Null Falsch-Positive** in `spec/architecture.md`. Die Sicht ist die
     einzige bewachte Datei; ein Befund dort, der kein Modul-Pfad ist, kippt den
     Kandidaten.
   - **Positivkontrolle:** ein eingesetzter Code-Modul-Pfad wird gemeldet, die
     Ursache gelesen, der Rückbau grün. Ohne sie ist „null Befunde" von „die
     Regel tut nichts" nicht zu unterscheiden.
   - **Die Falsch-Negativ-Klasse ist zu benennen**, nicht zu minimieren.
   - **Keine dritte Mechanik** (§3): die Wahl fällt zwischen der aus
     [slice-143](../done/slice-143-structure-abschnitts-skopus.md) und der aus
     [slice-144](../done/slice-144-commit-hash-muster.md), und sie wird
     **begründet**, nicht bequem getroffen.

1. **Ort und Form klären**: ein Muster-Verbot (dieselbe Fähigkeit, die
   [slice-144](../done/slice-144-commit-hash-muster.md) und
   [slice-143](../done/slice-143-structure-abschnitts-skopus.md) berühren) oder ein
   Ventil an `codepaths`. Die drei Slices teilen eine Frage — das gehört
   gesehen, bevor drei Mechaniken entstehen.
2. Am Bestand messen; die Sicht startet grün, andere Dateien womöglich nicht.
3. Konstruierter Verstoß mit gelesener Ursache; Rückbau grün.
4. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine dritte Muster-Mechanik.** Fällt die Entscheidung in slice-143 oder
  slice-144 anders, folgt dieser Slice ihr — er erfindet nichts Eigenes.
- **Keine Ausweitung auf die anderen Spec-Straten.** Die Spezifikation **darf**
  Pfade führen.

## 4. Definition of Done

- [ ] Die Wegwahl ist mit den beiden verwandten Slices **abgeglichen**, nicht
      unabhängig getroffen.
- [ ] Konstruierter Verstoß rot mit gelesener Ursache; Rückbau grün.
- [ ] §3.4 trägt den eingelösten Trigger für **diese** Hälfte; die andere bleibt
      *permanent*.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Drei Slices, eine Frage.** Wer sie einzeln löst, baut drei Antworten auf
  dieselbe Lexik-Frage — genau der Defekt, den dieses Repo an fremden Skripten
  ablehnt. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Wegwahl ein Produkt-Delta verlangt, das in einen der verwandten Slices gehört.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Module (GF), Spec-Straten (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-25): [`BEO-012`](../observations.md), weil der Slice auf eine Verschärfung aufsetzt, deren Grundlage eine Lesart ist.

Slice-ID: slice-145. Betroffene IDs: — (abhängig von der Wegwahl). Module: `codepaths` bzw. Muster-Verbot.
Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Sensor auf eigenem Bestand.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
