# Slice slice-091: Regelwerk-Migration Etappe D — Slice-`Status:`-Feld entfernen (Lifecycle = Verzeichnis)

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis**, in dem die Datei liegt
(`open/`/`next/`/`in-progress/`/`done/`) — **kein** `Status:`-Feld; Wechsel nur per
`git mv` (Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State
Machine). Dieser Slice trägt selbst schon die Ziel-Form (kein Status-Feld).

**Welle:** welle-67-baseline-v500-migration (Etappe D, **letzter** „Mini-Welle"-Slice,
nach [slice-090](../done/slice-090-review-infrastruktur.md)).

**Bezug:** Etappe D (Form-Konformität) aus
[slice-085](../done/slice-085-etappe-b-modul-delta.md) §3.2 — **D-5**: das
Slice-`**Status:**`-Feld **dupliziert** die Lifecycle-Verzeichnis-Position; die
Baseline (`modul-05` §Lifecycle: „Lifecycle-Zustand = Verzeichnis, kein Status-Feld";
`slice.template.md` führt einen `**Lifecycle:**`-Hinweis statt eines Status-Felds).
Mini-Welle abgeschlossen: slice-088 Planning-Layer · slice-089 Doc-Form · slice-090
Review-Infrastruktur → **slice-091** (dieser). **Danach schließt welle-67.**
**Kein Release.**

**Autor:** pt9912. **Datum:** 2026-08-02.

---

## 1. Ziel

Das `**Status:**`-Feld aus den Slice-Dateien entfernen — der Lifecycle-Zustand **ist**
die Verzeichnis-Position (`done/` usw.), das Feld ist reine Duplikation. **NB:** Die
**andere** `Status:`-Achse bleibt unberührt: das MR-`Status: Accepted` (Akzeptanz,
template-vorgeschrieben, C-6) und der ADR-`**Status:**` (immutable). Betrifft nur die
`docs/plan/planning/**/slice-*.md`-Dateien.

## 2. Vorgehen

1. **Umfang messen + Abnahme-Punkt (§3).** **90** Slice-Dateien tragen `**Status:**`
   (alle in `done/`); die Formate variieren (ein-/mehrzeilig, mit Datum/Welle-Prosa).
   Kein Gate/Skript liest das Slice-Status-Feld (nur ADR-Status in Go-Tests) → gate-sicher.
2. **Ziel-Form etablieren (template-forward, Abnahme-Punkt 1).** slice-091 **modelliert**
   die Ziel-Form: **kein** `**Status:**`-Feld, stattdessen der `**Lifecycle:**`-Hinweis.
   Die 90 `done/`-Slices behalten ihr Feld (ruhende Audit-Records, kein Massen-Touch);
   **künftige** Slices tragen kein Status-Feld mehr.
3. **Konvention nachziehen.** Die Erwähnung der „Status-Zeile" im Move-Commit-Eintrag
   des Konventionsspeichers (die Lifecycle-Move-Adaption) auf „Body = DoD-Haken +
   Closure-Notiz" umschreiben (das Feld existiert nicht mehr). Prüfen, ob `AGENTS.md` §5
   (Slice-Lifecycle) einen „kein Status-Feld"-Zusatz braucht.
4. **Gate.** `make gates` + `make adr-check` grün (keine ADR/Spec/Code berührt);
   unabhängiger Frischkontext-Review.

## 3. Abnahme-Punkte

1. **Retrofit-Umfang (D-5).** Das Feld aus **allen 90** `done/`-Slices entfernen
   (uniform, großer mechanischer Touch auf ruhende Audit-Records) **vs.** nur
   template-forward (neue Slices ohne Feld; die 90 behalten ihres → gemischter Bestand).
   → **Entschieden 2026-08-02: template-forward-only.** Neue Slices ohne Status-Feld
   (slice-091 modelliert es); die 90 `done/`-Slices behalten ihres — ruhende
   Audit-Records ohne Churn. Die Konvention wird go-forward verankert (AGENTS §5,
   Lifecycle-Move-Adaption). *(Ein späterer sauberer Retrofit bleibt möglich, ist aber
   bewusst nicht Teil dieses Slice.)*

## 4. Definition of Done

- [x] slice-091 + die go-forward-Konvention ohne `**Status:**`-Feld (`**Lifecycle:**`-
  Hinweis); die 90 `done/`-Slices unangetastet (template-forward, Abnahme-Punkt 1).
- [x] Die „Status-Zeile"-Erwähnung im Lifecycle-Move-Konventionseintrag nachgezogen;
  `AGENTS.md` §5 geprüft.
- [x] `make gates` + `make adr-check` grün; unabhängiger Frischkontext-Review.

## 5. Risiken / offene Punkte

- **Mechanischer Massen-Touch** (90 Dateien): der `**Status:**`-Absatz variiert
  (mehrzeilig) → Skript + Diff-Sicht, kein Blind-`sed`; Datum/Welle-Info darf nicht
  verloren gehen, ohne dass sie anderswo (Closure-Notiz/git) steht.
- **Zwei `Status:`-Achsen nicht verwechseln:** MR-`Status: Accepted` (bleibt) und
  ADR-`**Status:**` (immutable) sind **nicht** betroffen.

## 6. Trigger

Abschluss von [slice-090](../done/slice-090-review-infrastruktur.md) (Review-Infra);
letzter Etappe-D-Slice.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** berührt *Harness/Prozess*-Doku (Slice-Header-Form) — greenfield, GF.
- **Offene Beobachtungen sichten:** `observations.md` = `— keine —`; nichts zu
  berücksichtigen.

## 8. Sub-Area-Modus-Begründung

GF (Repo-Default): Doc/Prozess führt. Berührt die Slice-Header-Konvention
(Harness/Prozess); greenfield-Form-Angleich an die adoptierte Baseline.

## 9. Closure-Notiz (nach `done/`)

Umgesetzt: das Slice-`**Status:**`-Feld ist **go-forward abgeschafft** — der
Lifecycle-Zustand **ist** die Verzeichnis-Position.

- **slice-091 modelliert die Ziel-Form:** kein `**Status:**`-Feld, stattdessen der
  `**Lifecycle:**`-Hinweis (wie die Baseline-`slice.template.md`). Dieser Slice ist
  zugleich der erste, dessen Closure **ohne** Status-Flip abläuft (Body-Commit = nur
  DoD-Haken + diese Notiz).
- **Nutzer-Entscheid: template-forward-only** — **kein** Retrofit der 90 `done/`-Slices
  (ruhende Audit-Records ohne Churn); nur neue Slices ohne Feld. *(Ein späterer sauberer
  Retrofit bleibt möglich, bewusst nicht Teil dieses Slice.)*
- **Konvention verankert:** `AGENTS.md` §5 (neue Slices ohne `**Status:**`, dafür der
  `**Lifecycle:**`-Hinweis; Alt-Slices behalten ihres), `AGENTS.md` §3.3 (Move-Commit-Body
  = DoD + Closure-Notiz; historische Slices auch die Status-Zeile) und die
  Lifecycle-Move-Adaption im Konventionsspeicher nachgezogen.
- Die **zwei anderen `Status:`-Achsen** bleiben unberührt: das MR-`Status: Accepted`
  (Akzeptanz, template-vorgeschrieben) und der ADR-`**Status:**` (immutable). Kein
  Gate/Skript liest das Slice-Status-Feld.

**Review:** unabhängiger Frischkontext-Review
(`docs/reviews/2026-08-02-slice-091-status-feld-review.md`) → **1 MEDIUM** (F-1: die
`AGENTS.md` §3.3-„Status-Zeile"-Spiegel-Stelle war beim ersten Wurf nicht mitgezogen) →
geheilt; **bestätigende Re-Review abnahmereif**. `make gates` + `make adr-check` grün,
**keine** ADR/Spec/Code berührt.

**Lehre:** eine Konventions-Erwähnung („Status-Zeile") steht oft an **mehreren
gekoppelten Stellen** (Lifecycle-Move-Adaption + `AGENTS.md` §5 **und** §3.3) — beim
Nachziehen **alle Spiegel** finden, nicht nur die Kanonisch-Quelle.

**Anschluss:** **Etappe D ist vollständig** — alle welle-67-Slices (084–091) liegen in
`done/`. Es folgt die **Welle-Closure-Prozedur** (Trigger-Audit → `welle-67-results.md`
→ Wellendokument `git mv` → `done/` → Wave-Self-Close-Commit → Roadmap fortschreiben);
danach ist die Baseline-Migration `v1.4.0` → `v5.0.0` **komplett**.
