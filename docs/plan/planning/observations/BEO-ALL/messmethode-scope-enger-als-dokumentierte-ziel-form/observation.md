# Der Scan-Bereich einer neuen Bedingung ist enger als die Ziel-Form, die sie belegen soll

**Sub-Area:** `*`

Eine neue Prüfung wird gegen den **Anlassfall** entworfen und dort korrekt
getestet — aber der Anlassfall erfüllt die dokumentierte Ziel-Form (Kanon
oder CR) nicht wörtlich, sondern über einen unbemerkten Umweg (hier: eine
zweite, nicht vorgeschriebene Kopie der geprüften Marke am Ort, den der
Code tatsächlich scannt). Die Prüfung wird grün, das Zitat der Ziel-Form im
Plan/ADR bleibt korrekt stehen, und die Diskrepanz zwischen „was zitiert
wird" und „was der Code tatsächlich liest" bleibt unsichtbar, weil kein
Test je einen Fall konstruiert, der die beiden trennt (Coverage-Metriken
sehen keine Lücke, weil derselbe Codepfad für „Marke fehlt" und „Marke
steht nur anderswo" durchlaufen wird). Gefunden nur durch unabhängigen
Review mit einer eigens konstruierten Fixture, die exakt die zitierte
Ziel-Form nachbildet — nicht durch Ausführen der Test-Suite.
