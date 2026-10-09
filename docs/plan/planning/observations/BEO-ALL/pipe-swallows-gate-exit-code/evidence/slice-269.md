**Vorgang:** slice-269
**Fund:** `make doc-check … | tail -3 && git commit …` — `doc-check` meldete
einen Befund, `tail` gab 0 zurück, die Kette lief in den Commit. Der
`pre-commit`-Hook fuhr `doc-check` selbst und lehnte ab; die Verkörperung hat
getragen, die Arbeitsregel nicht.
