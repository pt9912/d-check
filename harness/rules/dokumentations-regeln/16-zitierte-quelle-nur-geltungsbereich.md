# Eine zitierte Quelle trägt nur, was in ihrem Geltungsbereich steht

Ausgelagerte Regel zu `AGENTS.md` §5 — Auslagerungs-Muster aus
[ADR-0096](../../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md).

Vor jedem Verweis das **Feld** lesen, nicht den Titel: bei `MR-<NNN>` den
`Geltungsbereich` **und** `Ersetzt-Baseline-Regel`, bei einer Kanon-Stelle
den Absatz, bei einer ADR die Unterscheidung **Akt gegen stehendes
Verbot**. Und die **direkteste** Quelle wählen — für eine Regel dieser
Datei ist das ihre Vorlage, nicht eine andere. Ein Zitat sieht aus wie ein
Beleg, auch wenn es keiner ist; das macht die Klasse beim Schreiben
unsichtbar und im Review auffindbar. Urteil, kein `grep`; der
Reviewer-Skill trägt den Anker dazu. Kanon:
[`grundlagen-source-precedence.md` §Wie weit trägt ein zitierter Satz](../../../.harness/baseline/v6.17.0/regelwerk/grundlagen-source-precedence.md)
— dort als Frage an **jede** zitierte Aussage, hier als operative Form für
den Implementer. *(Hard Rule aus dem Steering Loop,
[`BEO-ALL/citation-stretched-beyond-scope`](../../../docs/plan/planning/observations/BEO-ALL/citation-stretched-beyond-scope/observation.md),
seit slice-147; Auflösungs-Trigger: permanent.)*
