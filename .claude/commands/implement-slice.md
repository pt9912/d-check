# Implement Harness Slice

Argument: $ARGUMENTS

Follow this exact workflow:

1. Read `CLAUDE.md`.
2. Read `harness/README.md`.
3. Read `AGENTS.md`.
4. Read `harness/conventions.md`.
5. Read the slice file passed as argument.
6. Read all referenced ADRs and requirements.
7. Report:
   - slice id
   - DC ids
   - ADR ids
   - affected modules
   - gates to run
8. Implement the smallest viable diff.
9. Run the narrowest relevant gate first.
10. Run `make gates`.
11. Update docs, ADR index, and planning lifecycle if a public contract is
    touched (`AGENTS.md` §5 — CHANGELOG is Release-Prep, not the feature
    commit).
12. Report changed files, gates run, failures, and residual risks.
13. Hand off to the `reviewer` agent type (no self-review, `AGENTS.md` §6
    Schritt 8/Modul 8), then to the `verifier` agent type before closure.
14. When a fix for review or verification findings writes more than those
    findings ask for — a new mechanism, a rework, a new reading path —, run
    another `reviewer` round over exactly the fix commit(s) before the
    `verifier` sees it; otherwise the largest code share of the slice is
    read by no review. (seit slice-258; permanent — the size of a fix is a
    judgement, no gate can see it)
15. Before each code commit, hold `git diff --cached --name-only` against the
    plan table in §3 and its plan-change notes. Every file not named there
    gets a plan-change note first, committed on its own before the code
    (`AGENTS.md` §6 Schritt 4). (seit slice-259; permanent — whether a
    plan-change note covers a file is a judgement)
16. Before each code commit, read every comment line the diff touches —
    also one only rewrapped or moved — against the five classes of
    `AGENTS.md` §3.7: no story of earlier behaviour, no origin prose, no
    external reference. (seit slice-259; permanent — chronicle is a
    judgement, no gate can see it)

Do not skip gates.
Do not claim completion without command output.
