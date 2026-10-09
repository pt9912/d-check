**Vorgang:** slice-268
**Fund:** Der getrennte Schnitt (reiner Move committet, dann der Nachzug)
scheiterte im pre-commit-Hook an [ADR-0014](../../../../../adr/0014-latest-tag-fuer-stabile-releases.md) (`core-drift-vcs`); geliefert wurde
der Move-Commit mit allen eingehenden Verweisen, die bewegte Datei darin
unverändert (R100). Der Schnitt stand nur in der Commit-Botschaft, nicht als
Plan-Änderung (R1 LOW-1, Verifikation V-1).
