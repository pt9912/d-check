# `make slice-mv` scheitert am eigenen pre-commit-Hook

**Sub-Area:** `tools/harness/`

Das Werkzeug committet den reinen Move, bevor es die Verweise nachzieht; der
pre-commit-Hook fährt `make doc-check` über den Arbeitsbaum, in dem die Verweise
in diesem Moment tot sind. Der Move-Commit fällt durch, und es bleibt ein
gestagter Rename. Die Übergänge laufen deshalb von Hand: Move und Verweise in
einem Commit, die bewegte Datei unverändert.
