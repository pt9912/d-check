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

Do not skip gates.
Do not claim completion without command output.
