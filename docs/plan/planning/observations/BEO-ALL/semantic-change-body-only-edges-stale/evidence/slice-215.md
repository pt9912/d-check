**Vorgang:** slice-215
**Fund:** Die Exit-3-Semantik wurde geändert — sie deckt jetzt **zwei**
Zustände statt einem — und war in **zwei von fünf** Spiegeln nachgezogen.
Stehengeblieben: [`AGENTS.md`](../../../../../../../AGENTS.md) §4 und die
`Makefile`-Zeile nannten weiter nur den *„neueren Release-Tag"*, und der Kopf
von `.github/workflows/upstream-drift.yml` sagte sogar **Falsches** — dort
stand, bei den Versions-Achsen heiße *„anders"* stets *„neuer"*, was seit der
Änderung für `baseline-freshness` nicht mehr stimmt.
**Ein Spiegel, der eine Semantik nicht nur wiederholt, sondern über sie
generalisiert, wird durch eine Änderung nicht nur veraltet, sondern falsch** —
und er stand in einem Artefakt, das der Slice ohnehin gelesen hatte, um seinen
Konsumenten zu prüfen.
