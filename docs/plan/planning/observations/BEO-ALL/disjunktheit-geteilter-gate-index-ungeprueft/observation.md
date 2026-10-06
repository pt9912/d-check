# Mehrere Autoritäts-Dateien eines Gate-Index werden vereinigt, ihre Disjunktheit prüft niemand

**Sub-Area:** `*`

Sobald das Modul `targets` die Vollständigkeit gegen mehrere
Autoritäts-Dateien misst, kann dasselbe Target in zwei Teilen stehen — etwa
im Index des Adopters und in dem, den ein Bootstrap-Werkzeug schreibt. Das
Modul meldet das bewusst nicht, weil die Regel „kein Target steht in beiden
Teilen" nur in einem offenen CR an die Baseline steht. Nimmt die Baseline sie
an, ist die Doppelnennung eine Drift zwischen zwei Eigentümern, und kein
Sensor fängt sie.
