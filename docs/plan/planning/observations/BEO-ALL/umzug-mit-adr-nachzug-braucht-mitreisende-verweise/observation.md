# Ein Umzug mit Link-Nachzug in einer ADR braucht die Verweise im Move-Commit

**Sub-Area:** `*`

`AGENTS.md` §3.3 verlangt für einen Umzug im Regelfall den reinen `git mv` und
danach einen eigenen Commit für die Inhaltsänderung; das Mitreisen anderer
Dateien im Move-Commit erlaubt der Wortlaut nur beim Übergang nach `done/`.
Zieht der Folge-Commit aber Link-Ziele in einer `Accepted`-ADR nach, prüft der
gestagte ADR-Check ihn gegen einen BASE-Stand, in dem der alte Pfad nicht mehr
existiert: das alte Ziel löst nicht auf, der Nachzug ist Drift. Ableiter: bei
einem Umzug, den eine `Accepted`-ADR verlinkt, die eingehenden Verweise im
Move-Commit mitnehmen und die bewegte Datei dort unverändert lassen.
