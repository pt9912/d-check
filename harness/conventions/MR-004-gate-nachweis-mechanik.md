# MR-004 — Gate-Nachweis-Mechanik und `.claude`-Hooks nach b-cad-Vorbild

- **Status:** Accepted
- **Datum:** 2026-06-10
- **Ersetzt-Baseline-Regel:** [`grundlagen-durchsetzungsschicht.md` §Drei Bindepunkte](../../.harness/baseline/v5.11.0/regelwerk/grundlagen-durchsetzungsschicht.md#drei-bindepunkte)
- **Geltungsbereich:** [`tools/harness/`](../../tools/harness/), [`.claude/`](../../.claude/), `make record-gates`
- **Adaption:** Übernahme der Working-Tree-Hash-Mechanik (`record-gates` als
  letzter `gates`-Prerequisite, Stop-Hook vergleicht den Hash) und der
  `.claude`-Hooks (PreToolUse-Guard, Stop-Gate) aus dem Repo `b-cad`. Das
  **realisiert** die drei Bindepunkte der Durchsetzungsschicht (Tool-Call-Gate
  via PreToolUse, Handoff-Gate via Stop) für dieses Repo mit konkreten Skripten.
- **Begründung:** Bewährte Mechanik gegen „Erfolgsmeldung ohne Gate-Lauf"; keine
  Logik-Dopplung zwischen Makefile und Hook.
- **Auflösungs-Trigger:** permanent.
