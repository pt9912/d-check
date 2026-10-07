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
als Einstiegspunkt, Bedingung *Disjunkt*; **umgesetzt am 2026-10-07**, Träger
`slice-252` (Antwort unten).

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

## Antwort

**Umgesetzt** — `targets` prüft die Disjunktheit opt-in über
`targets.authority-disjoint: true`; Begründung in
[ADR-0101](../adr/0101-targets-authority-disjunkt.md).

- **Grund-Code:** `gate-declared-twice`.
- **Was als „steht in einem Teil" zählt:** die Tabellenzeile, die das Target
  in ihrer **ersten Zelle** führt — wie die Baseline es für die Target-Zelle
  beschreibt. Eine Erwähnung in einer anderen Zelle (etwa „eingehängt in
  `make gates`" in der Vertrag-Spalte eines Werkzeug-Teils) zählt nicht.
- **Fundstelle:** die Datei, die ein Target in Konfigurations-Reihenfolge
  zuerst führt, ist seine Heimat; gemeldet wird jede führende Zeile in einer
  späteren Datei, die Meldung nennt die Heimat.
- **`exempt-targets`** nimmt von der Disjunktheit nicht aus.
- **Unabhängig von `targets.makefiles`:** die Prüfung braucht keine
  Regelmenge.
- **Grenzen:** eine Doppelung innerhalb **einer** Datei ist kein Fall; zwei
  Einträge sind dieselbe Datei, wenn ihr bereinigter Pfad gleich ist — ein
  symbolischer Link auf eine andere Autoritäts-Datei zählt als zweite Datei.
- **Ohne Schalter** bleibt die Doppelnennung still, byte-identisch zum
  bisherigen Verhalten.

Für die Probe s30d heißt das: mit `targets.authority-disjoint: true` wird das
doppelt geführte Target laut (Exit 1, `gate-declared-twice`), sofern beide
Teile es in der ersten Zelle führen.
