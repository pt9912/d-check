**Vorgang:** slice-223
**Fund:** Der ursprüngliche Plan zitierte *„Beide Commits gehören in denselben
Push … der Zwischenstand ist zulässig, solange er nicht die Spitze eines
Push wird"* als Beleg dafür, [`MR-013`](../../../../../../../harness/conventions.md#mr-013)s drei Fälle aufzulösen — die Stelle
steht aber im Abschnitt über die **Anker-Paarung** bei der Welle-Closure,
nicht über den Slice-Lifecycle. Der Plan selbst hat das als offenes Risiko
benannt, bevor implementiert wurde (§6, vorab). Bei der Implementierung
zeigte sich: Für **zwei** der drei Fälle (Slice-Lifecycle-Move,
Beanspruchung) trägt trotzdem eine andere, tragfähige Begründung (reine
Git-Semantik — die bewegte Datei bleibt unverändert). Für den **dritten**
(MR-/Wellen-Lifecycle-Move) trägt sie nicht, und [`MR-013`](../../../../../../../harness/conventions.md#mr-013)s
`Ersetzt-Baseline-Regel`-Feld zeigte obendrein auf die falsche Kanon-Stelle
(`modul-05` statt `grundlagen-traceability.md` §Ruheort-Regel) — beide
Fehlzitate sind jetzt korrigiert.
