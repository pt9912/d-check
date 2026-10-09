**Vorgang:** slice-260
**Fund:** `make verify-closure-notes` verspricht über `done/`, liest aber nur
die Slices direkt darin — die 17 Volltexte unter `done/wellenlos/` blieben
ungeprüft (R1 F-1, HIGH; Behebung slice-263). Der Übergangs-Wächter in
`pre-commit` und CI erkannte dieselben Unterverzeichnisse nicht und löste seit
2026-09-29 bei keiner Closure aus.
