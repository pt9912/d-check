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
2. **Feld entfernen.** Je Slice-Datei den `**Status:**`-**Absatz** (die Zeile + etwaige
   Folgezeilen bis zur nächsten Leerzeile) tilgen — robust per Skript, danach Diff-Sicht.
   Neue Slices tragen stattdessen den `**Lifecycle:**`-Hinweis (Ziel-Form; dieser Slice
   modelliert ihn).
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
   → **Default/Empfehlung: retrofit-all** — die Baseline-Regel „kein Status-Feld" ist
   eindeutig, das Feld ist überall reine Lifecycle-Duplikation; die Datum/Welle-Info
   liegt in §Closure-Notiz + git.

## 4. Definition of Done

- [ ] Kein `docs/plan/planning/**/slice-*.md` trägt mehr ein `**Status:**`-Feld
  (gemäß Abnahme-Punkt 1).
- [ ] Die „Status-Zeile"-Erwähnung im Lifecycle-Move-Konventionseintrag nachgezogen;
  `AGENTS.md` §5 geprüft.
- [ ] `make gates` + `make adr-check` grün; unabhängiger Frischkontext-Review.

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

_Ausstehend._
