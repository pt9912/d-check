**Vorgang:** slice-262
**Fund:** Review R1 F-1 und die Verifikation: im Wegwerf-Repo mit einem
Commit und in einem Klon der Tiefe 1 bricht der Hook einen Commit mit Kennung
mit Exit 2 ab. Dazu sagen der Kommentar über `check_latest` in
`fetch-baseline-cache.sh` und der Hilfetext von `baseline-freshness` im
`Makefile` weiter pauschal „SKIP" bei einem Ausfall (R2, Verifikation B-4) —
dieselbe Werkzeug-Datei-Grenze.
