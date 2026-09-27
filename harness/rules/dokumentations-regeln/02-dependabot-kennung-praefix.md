# Dependabot-Commits tragen die Kennung im Präfix, nicht in einer Ausnahme

Ausgelagerte Regel zu `AGENTS.md` §5 — Auslagerungs-Muster aus
[ADR-0096](../../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md).

`commit-message.prefix` in
[`.github/dependabot.yml`](../../../.github/dependabot.yml) lautet
`build(deps) [ADR-0067]` bzw. `build(ci) [ADR-0067]`; damit erfüllt <!-- d-check:ignore (literale Konfigurationswerte, keine Verweise) -->
jeder Bump-Commit dieselbe Regel wie jeder andere
([ADR-0067](../../../docs/plan/adr/0067-dependabot-als-hebender-kanal.md)).
Eine Erweiterung von `commits.exempt-pattern` machte den Gate für eine
**ganze Commit-Klasse** blind; das ist der Grund gegen sie, nicht §3.6 — der
**verbietet** eine Lockerung nicht, sondern verlangt eine ADR dafür. Der
Grund ist ein sachlicher, kein verfahrensmäßiger. **Die Kennung gilt dem
Kanal, nicht dem Inhalt des einzelnen Bumps**; wer mehr Bezug hineinliest,
liest zu viel.
