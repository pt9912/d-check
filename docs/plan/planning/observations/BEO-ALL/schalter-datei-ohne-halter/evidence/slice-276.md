**Vorgang:** slice-276
**Fund:** Die erste Fassung des Stop-Hooks hielt den strengen Modus über die
versionierte Datei `.harness/stop-gate-streng` (so liefert es `ai-harness-init`
v0.10.0 mit); ohne sie galt der lockere Zweig. Der Review (HIGH) maß im
Scratch-Repo: Datei gelöscht ⇒ ein Turn-Ende mit ungedeckter Änderung gibt frei.
Das Plan-Risiko „Datei als einziger Halter" stand im Plan und war als offen bis
Closure notiert, nicht vor dem Code beantwortet. Behoben im Design: streng ohne
Schalter, locker nur mit einer angelegten Datei.
