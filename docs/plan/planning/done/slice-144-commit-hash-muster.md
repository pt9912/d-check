# Slice slice-144: Commit-Hashes in den Spec-Straten — ein Muster mit vertretbarer Falsch-Positiv-Last

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [slice-138](../done/slice-138-matrix-wellen-klasse.md); [`AGENTS.md`](../../../../AGENTS.md) §3.4; [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix).

**Berührte Spec-Stellen:** — (Konfigurations-Profil; eine Anforderung wächst nur, falls das Muster eine neue Fähigkeit braucht).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-25.

---

## 1. Ziel

§3.4 verbietet den Spec-Straten fünf Referenz-Kategorien; drei sind gedeckt.
Der **Commit-Hash** wäre als Token-Klasse ausdrückbar — die Mechanik existiert —,
aber ein Muster über Hex-Zeichenketten träfe jedes Wort, das wie ein Hash
aussieht. Was fehlt, ist **Präzision**, nicht die Fähigkeit.

Der Slice beantwortet genau eine Frage: **gibt es ein Muster, dessen
Falsch-Positiv-Last am eigenen Bestand vertretbar ist?** Wenn nein, ist das das
Ergebnis, und die Kategorie bleibt ausgewiesen.

## 2. Vorgehen

0. **Die Schwelle steht, bevor gemessen wird** — festgehalten am 2026-08-26,
   vor dem ersten Lauf, weil „vertretbar" sonst nachträglich an das Ergebnis
   angepasst wird ([`BEO-011`](../observations.md), §7):

   - **Null Falsch-Positive** im Bestand der drei Spec-Straten. Ein Gate in
     `gates`, das legitime Sätze meldet, erzwingt Umformulierung — und §3
     verbietet die Ausnahme-Liste als Ersatz für Präzision.
   - **Positivkontrolle:** ein konstruierter echter Commit-Hash wird gemeldet,
     die Ursache gelesen, der Rückbau grün.
   - **Die Falsch-Negativ-Klasse ist zu benennen**, nicht zu minimieren: welche
     echten Hashes das Muster nicht fängt, gehört ins Ergebnis.

   Ein Kandidat, der die erste Bedingung nur mit Ausnahmen erfüllt, ist
   durchgefallen.

1. **Kandidaten-Muster formulieren** (Mindestlänge, Wortgrenzen, Ausschluss von
   Inline-Code und Links) und je Kandidat am Bestand messen.
2. **Beide Fehlerrichtungen zählen**, nicht nur die Treffer: was würde gemeldet,
   das legitim ist — und was bliebe unentdeckt?
3. Nur scharfschalten, wenn die Messung trägt; sonst §3.4s Ausweisung schärfen.
4. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Ausnahme-Liste als Ersatz für Präzision.** Ein Muster, das nur mit
  einer langen Ausnahme-Liste grün wird, ist das falsche Muster.
- **Keine Ausweitung auf andere Dokumentklassen.**

## 4. Definition of Done

- [ ] Je Kandidat eine **Messung beider Fehlerrichtungen** am eigenen Bestand.
- [ ] Entweder scharfgeschaltet mit konstruiertem Verstoß und gelesener Ursache —
      oder die Kategorie bleibt ausgewiesen, mit der Messung als Begründung.
- [ ] §3.4 sagt danach, was gilt.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein Hex-Muster ist ein Heuristik-Wächter.** Genau die Sorte, die welle-84
  ausgeschlossen hat — hier ist sie zulässig **nur**, wenn die Messung sie trägt.
  — **Ausgang:** *(bei Closure)*
- **„Vertretbar" ist ein Urteil.** Die Schwelle gehört vorher benannt, nicht
  nachträglich an das Ergebnis angepasst. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls kein Muster die Messung besteht — dann ist die Ausweisung das Ergebnis.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Konfigurations-Profil (GF), Spec-Straten (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-25): [`BEO-011`](../observations.md) ist zentral: die Schwelle für „vertretbar" gehört **vor** die Messung.

Slice-ID: slice-144. Betroffene IDs: [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix). Module: `matrix`, Konfigurations-Profil.
Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Messung an bestehender Modul-Mechanik.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
