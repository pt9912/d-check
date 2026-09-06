# `make lint` — fährt das SOLID-nahe golangci-lint-Profil dieses Repos

## Vertrag

24 kalibrierte Linter: 5 Default- plus 23 aus
[ADR-0006](../../docs/plan/adr/0006-lint-profil-solid.md), dazu `nolintlint`
aus [`AGENTS.md`](../../AGENTS.md) §3.2. Ausnahmen leben **zentral** in
`.golangci.yml` (`exclude-rules`) mit Begründung; Inline-Suppressions sind
verboten.

## Grenze — was das Grün nicht abdeckt

1. **`nolintlint` prüft die Form der Direktive, nicht ihre Berechtigung.**
   Gemeldet wird eine Direktive ohne benannten Linter, ohne Begründung oder
   ohne Wirkung — sie wird damit sichtbar und zurechenbar. Eine
   **wohlgeformte** Direktive unterdrückt einen echten Verstoß weiterhin, und
   der Lauf bleibt grün. Verboten bleibt sie durch
   [`AGENTS.md`](../../AGENTS.md) §3.2, nicht durch dieses Gate. Permanent —
   die Berechtigungsfrage ist ein Urteil.
2. **Ein Linter macht lokale Mustererkennung** — Datenfluss über
   Funktionsgrenzen und Struktur-Regeln fängt er nicht; dafür stehen
   [`semgrep`](semgrep.md) und [`arch-check`](arch-check.md) daneben.

## Bindung

Bestandteil von `make gates`.
[ADR-0006](../../docs/plan/adr/0006-lint-profil-solid.md) ·
[`AGENTS.md`](../../AGENTS.md) §3.2
