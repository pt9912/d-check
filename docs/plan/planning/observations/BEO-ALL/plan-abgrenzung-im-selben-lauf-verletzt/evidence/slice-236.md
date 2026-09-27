**Vorgang:** slice-236
**Fund:** §1 des Slice-Plans schließt „Handbuch, README, CHANGELOG, Release —
Release-Prep-Vorgang, kein Feature-Commit" ausdrücklich aus. Der
Umsetzungs-Commit (`92ec638a`) hat trotzdem README.md/README.de.md geändert
(neue `file`-Modul-Bullets). Gefunden vom unabhängigen Review (F-1, HIGH,
merge-blockierend), nicht vom Implementer-Lauf selbst — obwohl Schritt 8 des
Workflows genau dafür da ist. Behoben durch Zurücknahme der README-Hunks in
einem Fix-Commit vor dem Review-Abschluss.
