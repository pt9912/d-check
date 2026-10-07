# Eingehender Hinweis — die Baseline führt die Disjunktheit des Gate-Index

**Absender:** Baseline `ai-harness-course` · **Eingegangen:** 2026-10-07
**Richtung:** eingehend — ein Hinweis, kein Change Request: er bittet um
nichts, er meldet, dass ein Re-Evaluierungs-Trigger dieses Repos eingetreten
ist.
**Bezug:** [ADR-0100](../adr/0100-targets-authority-liste.md)
(Re-Evaluierungs-Trigger), Beobachtung
[`BEO-ALL/disjunktheit-geteilter-gate-index-ungeprueft`](../planning/observations/BEO-ALL/disjunktheit-geteilter-gate-index-ungeprueft/state.md)
**Stand:** **eingegangen und geprüft** — die Regel steht in `v6.16.0`
(veröffentlicht 2026-10-06), `grundlagen-harness-dateien.md` §harness/README.md
als Einstiegspunkt, Bedingung *Disjunkt*; Umsetzung trägt `slice-252`.

**Ablage-Hinweis.** Abgelegt neben den eingehenden CRs, weil er dieselbe Frage
beantwortet wie sie: was von außen kam und wie dieses Repo darauf reagiert.
Im Wortlaut sind Kennungen verlinkt (Linkpflicht des eigenen Gates), sonst
unverändert.

---

## Wortlaut

> **An d-check — Hinweis zu [ADR-0100](../adr/0100-targets-authority-liste.md)**
>
> Die Baseline führt seit v6.16.0 (Welle 159) die Regel „kein Target steht in
> zwei Teilen des Gate-Index" (grundlagen-harness-dateien.md
> §harness/README.md als Einstiegspunkt). Damit ist der
> Re-Evaluierungs-Trigger von [ADR-0100](../adr/0100-targets-authority-liste.md) erreicht. Eine opt-in-Prüfung der
> Disjunktheit mit eigenem Grund-Code hätte jetzt eine Grundlage.
>
> Der Kurs führt die Disjunktheit bis dahin als benannte Grenze. Der Team-Sim
> probt die heutige Stille als s30d: Das doppelt genannte Target bleibt still,
> während Phantom-Proben in beiden Teilen im selben Aufbau laut werden.
> Liefert d-check die Prüfung, dreht s30d auf laut, und der Kurs zieht das
> nach.
