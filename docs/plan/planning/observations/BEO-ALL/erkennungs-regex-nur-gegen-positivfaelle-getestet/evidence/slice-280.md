**Vorgang:** slice-280
**Fund:** Die Negativliste prüfte, was die neue Variablen-Erkennung nicht treffen darf,
nicht was sie an **bereits Erkanntem** ändert: ein Unix-Fund hinter `${HOME}/…` meldet
mit gesetztem `variables` unter dem Variablen-Pfad, ein Bestands-Glob auf den Unix-Teil
trifft dann nicht mehr (Review F-2); überlappende Treffer mit Komma als Trenner
meldeten doppelt (F-3). Dieselbe Richtung wie in slice-279: die Wirkung der neuen
Regel auf den bisher erkannten Bestand stand nicht in der Liste.
