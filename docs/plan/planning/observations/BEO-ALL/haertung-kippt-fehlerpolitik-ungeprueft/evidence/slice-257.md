**Vorgang:** slice-257
**Fund:** Der R1-Fix las die Plattformen eines Image-Index aus der Registry,
statt sie aus einer Kopie zu nehmen — eine Härtung. Der neue Leseweg
verwarf stderr und Exit von `imagetools` (`2>/dev/null … || true`); gefahren
wurde er nur gegen saubere Indexe. Ein Index mit einem Eintrag ohne Plattform
lieferte eine Teil-Liste, und der Scan meldete Exit 0, ohne arm64 je gescannt
zu haben (R2-1 HIGH, gemessen). Die Behebung fügte selbst den nächsten Rand
hinzu: stderr in der Plattformliste (R3-1).
