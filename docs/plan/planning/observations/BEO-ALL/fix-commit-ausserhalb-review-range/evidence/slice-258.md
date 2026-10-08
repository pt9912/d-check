**Vorgang:** slice-258
**Fund:** Der R1-Fix brachte neben den Befunden einen rekursiven
`make build VERSION=0.0.0-dev` und die `:arm64`-Bereinigung in `make clean`
mit, außerhalb der R1-Range (Verifikation V-2); die Runden R2 bis R4 über die
Fix-Commits fanden danach jeweils neue MEDIUM-Befunde.
