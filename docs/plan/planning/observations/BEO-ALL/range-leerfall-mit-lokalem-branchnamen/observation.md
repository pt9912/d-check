# `--range <sha>..<lokaler Branchname>` meldet den Range-Leerfall trotz Commit

**Sub-Area:** `*`

In einem Wegwerf-Klon meldeten altes und neues Image für
`--range <sha>..<lokaler-Branchname>` den Range-Leerfall (0 Commits, Exit 2),
obwohl ein Commit dazwischen lag. Fail-closed, aber womöglich eine Lücke in
der Ref-Auflösung des Adapters. Ableiter: dieselbe Range einmal mit
Branchname und einmal mit `HEAD` fahren und vergleichen.
