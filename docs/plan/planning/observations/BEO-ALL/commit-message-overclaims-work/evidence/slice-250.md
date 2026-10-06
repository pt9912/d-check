**Vorgang:** slice-250
**Fund:** Richtung (b), zweimal im selben Vorgang. Die Feature-Botschaft
sagte „keine Symlinks" und „SKIP_DIRS wie file" zu, ohne dass der Code es
trug (R1). Die Fix-Botschaft nannte als Rot-Grund eines Testfalls „doppelt
gelesen", der Test fing die Mutation aber an einem Ersatz-Symptom (die
Test-Attrappe kannte `./`-Pfade nicht); das behauptete Symptom zeigte erst
die Black-Box des Verifiers (V1). Behoben durch eine Test-Attrappe, die die
Pfad-Schreibweise wie das echte Dateisystem auflöst.
