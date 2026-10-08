# `make lint` — fährt das SOLID-nahe golangci-lint-Profil dieses Repos

## Vertrag

Hält das SOLID-nahe Lint-Profil dieses Repos ([ADR-0006](../../docs/plan/adr/0006-lint-profil-solid.md)).
Welche Linter mit welchen Schwellen laufen und welche Ausnahmen gelten, legt
[`SPEC-090`](../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge) fest; Inline-Suppressions verbietet [`AGENTS.md`](../../AGENTS.md) §3.2.

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

3. **Die Ausschluss-Regeln verkleinern den Prüfbereich** — zentral in
   [`.golangci.yml`](../../.golangci.yml) deklariert, mit Begründung, wie
   `AGENTS.md` §3.2 es verlangt; welche es sind, zählt
   [`SPEC-090`](../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)
   auf. Ein grüner Lauf sagt „sauber außerhalb dieser Ausnahmen".

## Bindung

Bestandteil von `make gates`.
[ADR-0006](../../docs/plan/adr/0006-lint-profil-solid.md) ·
[`AGENTS.md`](../../AGENTS.md) §3.2
