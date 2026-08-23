# MR-005 — Härtung ggü. b-cad: inhaltsbasierter Gate-Nachweis, Sub-Shell-Prüfung

- **Status:** Accepted
- **Datum:** 2026-06-10
- **Ersetzt-Baseline-Regel:** [`grundlagen-durchsetzungsschicht.md` §Vier Design-Eigenschaften](../../.harness/baseline/v5.11.0/regelwerk/grundlagen-durchsetzungsschicht.md#vier-design-eigenschaften)
- **Geltungsbereich:** [`tools/harness/working-tree-hash.sh`](../../tools/harness/working-tree-hash.sh), [`.claude/hooks/`](../../.claude/hooks/)
- **Adaption:** Zwei Abweichungen von der per
  [`MR-004`](../conventions.md#mr-004--gate-nachweis-mechanik-und-claude-hooks-nach-b-cad-vorbild)
  übernommenen b-cad-Mechanik:
  (a) Der Working-Tree-Hash ist **inhaltsbasiert** (sha256 über alle getrackten +
  untracked Dateiinhalte) statt diff-basiert — genau die Design-Eigenschaft
  „Nachweis über Inhalt, nicht Diff" der Durchsetzungsschicht. Damit gilt der
  Gate-Nachweis über Commits hinweg (gleicher Inhalt = gleicher Hash), und ein
  Commit *ohne* Gate-Lauf macht den Stop-Hook nicht mehr grün. Restlücke bleibt:
  frischer Klon bzw. gelöschter `.harness`-State mit cleanem Tree wird freigegeben —
  dort ist CI das Netz.
  (b) Der PreToolUse-Guard prüft Sub-Shell-Strings (`bash -c "…"`, `sh -c '…'`)
  rekursiv (Tiefe ≤ 3, darüber fail-closed).
- **Begründung:** Review-R2-Beobachtungen (User): Commit-Bypass des Stop-Hooks und
  Guard-Umgehung via `bash -c`.
- **Auflösungs-Trigger:** permanent. Rückport beider Härtungen nach b-cad steht aus.
