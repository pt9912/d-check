**Vorgang:** slice-211
**Fund:** Die Beanspruchungs-Commits von slice-210 **und** slice-211
reklamieren einen *„Lauf unter [`MR-070`](../../../../../../../harness/conventions.md#mr-070)"* — dessen **Geltungsbereich den
Pfad-Nachzug eines Lifecycle-Moves ausdrücklich ausnimmt**. Beide Commits taten
genau das und nichts anderes. **Die Ausnahme stammt aus slice-209, drei
Commits zuvor, und wurde dort von mir selbst geschrieben**, nachdem der Review
die Kollision mit [`MR-013`](../../../../../../../harness/conventions.md#mr-013)
gemeldet hatte. Der Eintrag gilt damit auch für die **eigene** frische Regel:
Ein Geltungsbereich, den man selbst formuliert hat, wird beim Anwenden nicht
automatisch gelesen. Commit-Botschaften sind nicht änderbar — der Fall bleibt
als Beleg stehen.
