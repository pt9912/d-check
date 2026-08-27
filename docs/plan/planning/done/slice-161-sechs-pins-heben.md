# Slice slice-161: Sechs gemeldete Pin-Rückstände heben — der Nachtlauf ist dauerrot

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [slice-142](../done/slice-142-freshness-weitere-achsen.md) (die Achsen, die melden); [ADR-0011](../../adr/0011-digest-pins-build-gate-images.md) (Digest-Pin-Politik); [ADR-0029](../../adr/0029-arch-check-via-a-check.md) (a-check-Fragment); [`AGENTS.md`](../../../../AGENTS.md) §3.9, §4.

**Berührte Spec-Stellen:** — (Pin-Stand; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

Der Nachtlauf wacht seit [slice-142](../done/slice-142-freshness-weitere-achsen.md)
über **zwölf** Achsen. **Sechs davon melden beim ersten Lauf Rückstand:**

| Achse | Pin | upstream |
|---|---|---|
| `checkout-pin-freshness` | `v6.0.2` | `v7.0.1` |
| `login-pin-freshness` | `v4.2.0` | `v4.6.0` |
| `freshness-semgrep` | `1.167.0` | `1.175.0` |
| `freshness-a-check` | `v0.8.0` | `v0.17.0` |
| `runtime-base-digest` | `sha256:d093aa3e…` | anderer Bau desselben Tags |
| `go-base-digest` | `sha256:65b6f280…` | anderer Bau desselben Tags |

Ein Nachtlauf, der **jede Nacht dasselbe Rot** zeigt, ist derselbe verwaiste
Sensor, gegen den er gebaut wurde: er verliert das Signal für die nächste Achse,
die *neu* rot wird. Die Achsen melden nur — das Heben ist ein bewusster Akt, und
dieser Slice ist er.

## 2. Vorgehen

1. **Je Pin eine eigene Entscheidung**, nicht ein Sammel-Bump. Ein
   Major-Sprung (`checkout` v6→v7, `a-check` v0.8→v0.17) ist eine andere
   Entscheidung als ein Digest-Nachzug am selben Tag.
2. **Form je Klasse aus [`BEO-008`](../observations.md):** eine Pin-Hebung hat
   Spiegel. Vor dem Bump je Pin die Spiegel zählen, nicht nur die grep-baren.
3. **Der `a-check`-Sprung ist der teuerste** — neun Minor-Releases, und das Gate
   fährt gegen unser Repo. Erst gegen den neuen Digest laufen lassen, dann
   pinnen; ein rotes `arch-check` aus fremder Regel-Verschärfung ist eine
   eigene Entscheidung, kein Nebeneffekt.
4. **Nach jedem Bump die zugehörige Achse fahren** und ihre echte Ausgabe lesen.
5. `make gates`, `make ci`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Auto-Bump-Mechanismus.** Die Achsen melden; wer hebt, entscheidet.
- **Keine Aufnahme netzhaltiger Achsen in `gates`.** Der Nachtlauf bleibt der
  Bindepunkt.
- **Keine neuen Achsen.** Die Menge ist mit slice-142 abgeschlossen.

## 4. Definition of Done

- [ ] Je gemeldeter Pin **entschieden**: gehoben, oder mit Grund stehen gelassen.
- [ ] Je Hebung die Spiegel aus [`BEO-008`](../observations.md) mitgezogen.
- [ ] Der Nachtlauf ist nach dem Slice **nicht mehr dauerrot** — oder das
      verbleibende Rot hat einen benannten Adressaten und ein Datum.
- [ ] `make gates` und `make ci` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein Major-Sprung kann ein Gate rot machen.** `a-check` v0.8→v0.17 sind neun
  Releases fremder Regel-Entwicklung; `actions/checkout` v6→v7 kann
  Runner-Verhalten ändern. — **Ausgang:**
- **Ein Digest-Nachzug am selben Tag ist nicht folgenlos.** Der Tag ist
  derselbe, der Inhalt nicht; ein neu gebautes Basis-Image kann eine andere
  libc-Ecke tragen. — **Ausgang:**
- **Das Heben schließt die Lücke nicht, es verschiebt sie.** Morgen ist der
  nächste Release da. Was dieser Slice nicht leistet, ist eine **Kadenz** — wer
  das Rot wann liest. — **Ausgang:**

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls ein Bump ein Gate rot macht und
die Ursache eine eigene Entscheidung braucht.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Werkzeuge (GF), CI (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27): [`BEO-008`](../observations.md) für die Pin-Spiegel; [`BEO-020`](../observations.md) für jede Aussage darüber, wie oft sich ein Pin bewegt; [`BEO-007`](../observations.md) für jeden Beleg-Lauf.

Slice-ID: slice-161. Betroffene IDs: [ADR-0011](../../adr/0011-digest-pins-build-gate-images.md), [ADR-0029](../../adr/0029-arch-check-via-a-check.md). Module: Harness-Werkzeuge, CI.
Gates: `make arch-check`, `make semgrep`, `make gates`, `make ci`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Pflege an bereits gebauter Mechanik.

## 9. Closure-Notiz (nach `done/`)

— (offen)
