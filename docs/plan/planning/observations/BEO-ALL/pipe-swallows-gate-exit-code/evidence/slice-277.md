**Vorgang:** slice-277
**Fund:** `make doc-check 2>&1 | tail -1 && git add … && git commit … && git push` —
`doc-check` endete mit Exit 2 (die eben geschriebene Plan-Notiz nannte
[MR-074](../../../../../../../harness/conventions.md#mr-074) ohne Link), `tail` gab 0 zurück, die Kette lief in Commit und Push.
Die CI meldete den Commit `failure`, der Folge-Commit war grün. Dazu: der lokale
`pre-commit`-Hook, der den vollen `doc-check` fährt und diese Klasse fängt, war in
diesem Klon nicht aktiviert (`core.hooksPath` leer) — die Verkörperung greift nur, wo
`make hooks` gelaufen ist.
