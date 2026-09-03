# Eine Härtung am Rand eines Wächters kippt seine Fehlerpolitik, unbemerkt von den Proben

**Sub-Area:** `*`

Eine Härtung entfernte eine Host-Abhängigkeit, damit ein fehlendes Werkzeug
den Wächter nicht mehr ohne Urteil enden lässt; der neue Leseweg
funktionierte aber nur über eine Pipe — unter `set -e` beendete er den
Wächter still bei anderer Eingabeform. Jeder Befehl lief unbemerkt durch,
weil die Proben nur die eine noch funktionierende Form fuhren. Die Klasse
ist nicht „Regression", sondern: ein Prüfling, den man am Rand anfasst,
wird an einer Stelle geprüft, die der Rand nicht berührt.
