Gefunden vom unabhängigen Review von slice-229
([Report](../../../../../../reviews/2026-09-18-slice-229-matrix-instanz-identitaet-review-r1.md),
F-1): [ADR-0087](../../../../../adr/0087-matrix-instanz-identitaets-ausnahme.md)
nannte `make doc-check` als Fitness-Function-Beleg für „das erweiterte
YAML-Config-Beispiel in `spec/spezifikation.md` bleibt syntaktisch gültig" —
`doc-check` dekodiert kein gefenctes YAML, und
`TestDocExamples_ConfigBeispieleValidieren` deckt nur vier Nutzer-Doku-Dateien
ab, `spec/spezifikation.md` nicht darunter. Da die ADR bereits `Accepted`
war, wurde die falsche Zeile nicht überschrieben, sondern per
`## Geschichte`-Nachtrag als falsch deklariert.
