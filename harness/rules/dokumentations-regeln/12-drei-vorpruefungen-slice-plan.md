# Jeder Slice-Plan trägt vor der Sub-Area-Modus-Begründung drei Vorprüfungen

Ausgelagerte Regel zu `AGENTS.md` §5 — Auslagerungs-Muster aus
[ADR-0096](../../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md).

Jeder Slice-Plan trägt **vor** der Sub-Area-Modus-Begründung die **drei**
Vorprüfungen: Sub-Area prüfen · offene Beobachtungen im Register
(`docs/plan/planning/observations/`) sichten (beide Baseline-Regelwerk
Modul 5/6, unabhängig vom Sub-Area-Modus) · den **Nachtlauf-Stand** lesen
(`make nightly-state`). Die dritte ist eine Adaption
([`MR-053`](../../conventions.md#mr-053)): der Kanon kennt keinen
Nachtlauf. Sie hängt an diesem Moment, weil dort ohnehin gelesen wird —
**benannte Grenze:** in einer Pause liest niemand. Der dritte Block
entsteht **spätestens bei der Beanspruchung**; ein Plan in `open/` trägt
ihn noch nicht.

**Die beiden kanonischen Blöcke belegen ihre Regel** mit einer
`d-check:cite`-Direktive auf die **vorschreibende** Regelwerk-Zeile, samt
wörtlichem Zitat darunter ([`MR-054`](../../conventions.md#mr-054));
`citations` prüft es wortgleich im inneren Loop, ein falsch angekerter
Beleg wird rot. Der dritte Block trägt bewusst keine — sein Ziel ist
repo-eigen und meldete bei jeder Änderung. **Kein Sensor hält das:** ein
Plan ganz ohne Direktiven ist grün.
