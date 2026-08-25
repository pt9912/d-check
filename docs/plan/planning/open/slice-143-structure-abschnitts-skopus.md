# Slice slice-143: Der Platzhalter-Erkenner des Produkts sieht nur einen Abschnitt

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [slice-139](../done/slice-139-closure-ausgang-waechter.md); `tools/harness/closure-outcomes.sh`; [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md); [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in), [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in).

**Berührte Spec-Stellen:** `spec/lastenheft.md` — eine Anforderung wächst (Abschnitts-Skopus bzw. Muster-Verbot); Bump und Historie nach [`MR-032`](../../../../harness/conventions.md#mr-032).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-25.

---

## 1. Ziel

Das Produkt erkennt Vorlagen-Platzhalter fence- und inline-code-bewusst — aber
**nur im Abschnitt der Closure-Notiz**. Für den Rest des Slice steht seit
[slice-139](../done/slice-139-closure-ausgang-waechter.md) ein Zeichenketten-
Wächter als Bash-Skript daneben, mit zwei benannten Schwächen: seine
Platzhalter-Liste ist eine **Liste**, und er behandelt Fenced Code nicht.

Die saubere Form ist der **Abschnitts-Skopus im Produkt**: dieselbe Erkennung,
weiter gefasst — oder ein `forbid-match` in `structure` als Gegenstück zum
vorhandenen `headings-match`. Welche der beiden, entscheidet der Slice.

## 2. Vorgehen

1. **Die zwei Wege gegeneinander stellen**, mit ihren Folgen: den Skopus des
   `planning`-Platzhalter-Erkenners weiten, oder `structure` um ein
   Muster-Verbot ergänzen. Beides ist ein Produkt-Delta mit ADR.
2. **Am Bestand messen**, bevor scharfgeschaltet wird — beide Wege treffen mehr
   als den einen Abschnitt.
3. Das Bash-Skript **ersatzlos entfernen**, sobald der Skopus es abdeckt: sein
   Auflösungs-Trigger steht in seinem Kopf.
4. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Doppelführung.** Am Ende trägt **eine** Mechanik die Frage; die
  andere fällt weg.
- **Keine Ausweitung auf andere Dokumentklassen** in diesem Zug.

## 4. Definition of Done

- [ ] Die Wegwahl ist **begründet**, mit der benannten Folge des verworfenen Wegs.
- [ ] Der Bestand ist gemessen; jede Fundstelle ist geräumt oder ausgewiesen.
- [ ] Das Bash-Skript ist entfernt, sein Auflösungs-Trigger eingelöst.
- [ ] Lastenheft-Bump + Historie; ADR geschrieben; `make fullbuild` grün.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein weiter gefasster Skopus meldet mehr.** Was heute in Slice-Abschnitten
  legitim in Winkelklammern steht, wird dann zum Befund — die Messung entscheidet,
  ob der Weg tragbar ist. — **Ausgang:** *(bei Closure)*
- **Zwei Mechaniken gleichzeitig sind schlechter als eine.** Bleibt das Skript
  neben dem Produkt stehen, ist die Doppelung dauerhaft statt übergangsweise. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Messung zeigt, dass der weitere Skopus den Bestand breit trifft.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Module (GF), Spec-Straten (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-25): [`BEO-011`](../observations.md) für die Bestandsmessung; [`BEO-015`](../observations.md), weil dieser Slice die dortige mechanische Form einlöst.

Slice-ID: slice-143. Betroffene IDs: [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in), [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in). Module: `planning`, `structure`.
Gates: `make doc-check`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Erweiterung einer bestehenden Modul-Fähigkeit.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
