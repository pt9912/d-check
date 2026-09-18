**Vorgang:** slice-230
**Fund:** Beim ersten `make gates`-Lauf auf einem macOS-Host schlug
`baseline-verify` an zwei unabhängigen Stellen in `fetch-baseline-cache.sh`
fehl — Datei-Anzahl-Vergleich (`wc -l`-Padding) und Symlink-Auflösung
(`readlink -e`), beide reine GNU-Coreutils-Annahmen. Beide behoben
(numerischer Vergleich, `readlink -f` + `[ -e ]`); `--selftest` (9 Proben)
und `make baseline-verify` liefen danach mit Standard-BSD-Werkzeugen grün.
