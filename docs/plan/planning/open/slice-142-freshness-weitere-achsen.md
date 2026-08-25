# Slice slice-142: Zwei ungewachte Pin-Klassen — Basis-Images und Action-Pins

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [slice-137](../done/slice-137-toolchain-freshness.md); `tools/harness/pin-freshness.sh`; [`AGENTS.md`](../../../../AGENTS.md) §3.9, §4.

**Berührte Spec-Stellen:** — (Harness-Sensoren; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-25.

---

## 1. Ziel

Drei gepinnte Fremd-Bestände sind gewacht (Baseline, Go, golangci-lint). Zwei
Klassen sind es nicht: die **Docker-Basis-Images** (digest-gepinnt im
`Dockerfile`) und die **Action-Pins** der Workflows, deren *Form* seit
`make workflow-pins` geprüft ist, deren **Gültigkeit** aber nicht.

Der Sensor ist parametriert; je Achse wäre es eine Zeile. **Ob sie eine Achse
brauchen, ist trotzdem eine Entscheidung** — ein Sensor, den niemand liest, ist
Aufwand ohne Wirkung.

## 2. Vorgehen

1. **Je Klasse messen, bevor gebaut wird:** wie oft hat sich der Pin
   tatsächlich bewegt, und was wäre die Folge einer verpassten Hebung?
2. **Quellen-Form klären.** Ein Digest hat keinen Release-Tag; die Frage
   *„gibt es einen neueren Digest für denselben Tag"* ist eine andere als
   *„gibt es einen neueren Tag"*.
3. Nur bauen, was die Messung trägt — und die Entscheidung gegen eine Achse
   **ausweisen**, nicht weglassen.
4. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Auto-Bump.** Wie bei den drei bestehenden Achsen: melden, nicht heben.
- **Keine Aufnahme in `gates`.** Netzhaltiges gehört in den Nachtlauf.

## 4. Definition of Done

- [ ] Je Klasse eine **Messung** und daraus eine begründete Entscheidung.
- [ ] Was gebaut wird, hängt im bestehenden Nachtlauf und ist fail-open.
- [ ] Was **nicht** gebaut wird, ist mit Grund ausgewiesen.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Eine Achse ist billig zu bauen und teuer zu ignorieren.** Wächst die Zahl
  der Nachtlauf-Meldungen, sinkt die Aufmerksamkeit für jede einzelne. —
  **Ausgang:** *(bei Closure)*
- **Digest-Achsen haben keine Tag-Semantik.** Die Gleich/Ungleich-Form der
  bestehenden Achsen trägt hier womöglich nicht. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Messung zeigt, dass eine Klasse eine andere Vergleichsform braucht.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Werkzeuge (GF), CI (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-25): [`BEO-011`](../observations.md) für die Bewegungs-Messung; [`BEO-007`](../observations.md) für jeden Beleg-Lauf.

Slice-ID: slice-142. Betroffene IDs: — (Harness-Sensoren; keine Anforderung). Module: Harness-Werkzeuge, CI.
Gates: `make gate-consistency`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Erweiterung an bereits gebauter Mechanik.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
