# Eine Festlegung im Slice-Plan zitiert die Stelle, die sie trägt

Ausgelagerte Regel zu `AGENTS.md` §5 — Auslagerungs-Muster aus
[ADR-0096](../../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md).

Eine **Festlegung** im Slice-Plan ist jeder Satz, nach dem sich Code, Test oder eine Vorlage an den
Auftraggeber richten wird: ein Entscheid in §3, eine Plan-Änderungs-Notiz, ein Testfall mit erwartetem
Ergebnis, eine Zusammenfassung der Richtung einer Änderung. Wer eine schreibt — Planner wie Implementer —,
beantwortet vorher: **Welche Stelle des Lastenhefts oder der ADR trägt sie?** Diese Stelle wird genannt
(verlinkte Kennung) und **zitiert**, nicht nacherzählt; eine Nacherzählung rundet ab, macht absolut oder
verschiebt die Richtung, und der Code folgt dem Plan, nicht der Quelle. Ein Zitat lässt sich als
`<!-- d-check:cite <datei>:<von>-<bis> -->` deklarieren; dann hält das Modul `citations` es gegen die
Quelle, und eine spätere Änderung des Lastenhefts macht den Plan rot statt still alt. **Trägt keine Stelle
die Festlegung, ist sie eine Entscheidung und keine Notiz:** Sie geht an den Architect (ADR) oder an den
Auftraggeber (Lastenheft), bevor Code ihr folgt.

**Grenzen:** Ob ein Satz eine Festlegung ist und ob das Zitat sie trägt, ist ein Urteil; `citations`
prüft nur, dass ein deklariertes Zitat noch in der Quelle steht, nicht, dass eine Festlegung eines hat.
Die Regel gilt dem Slice-Plan; Spec, ADR und Handbuch folgen ihren eigenen Regeln (3, 13, 16). Ob eine
neue Anforderung einem globalen Ausschluss des Lastenhefts widerspricht, ist eine andere Frage und hier
nicht gedeckt. *(Hard Rule aus dem Steering Loop,
[`BEO-ALL/plan-entscheid-steht-gegen-lastenheft`](../../../docs/plan/planning/observations/BEO-ALL/plan-entscheid-steht-gegen-lastenheft/observation.md),
seit welle-96; Auflösungs-Trigger: permanent.)*
