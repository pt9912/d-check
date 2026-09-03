# `fetch-baseline-cache.sh --check-latest` sieht einen neuen Release-Tag nicht, solange der eigene Pin noch dahinter liegt

**Sub-Area:** `tools/harness/`

Gemessen beim `v6.0.0`-Bump: `make baseline-freshness` meldete weiterhin nur
`v5.19.0`/`v5.20.0` als neuere Tags, obwohl `v6.0.0` laut direkter
GitHub-Release-API zu dem Zeitpunkt bereits publiziert war. Sobald der Pin
selbst auf `v6.0.0` gehoben war, meldete `check-latest` korrekt „ist der
neueste Tag" — die Lücke betrifft offenbar nur die Sicht *vor* dem eigenen
Pin, nicht danach. Ursache nicht untersucht (kein Blocker für den Bump
selbst, da die Existenz des Tags unabhängig per API bestätigt wurde).
