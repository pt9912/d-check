# Ein Pfad-Nachzug auf eine gleichnamige andere Datei geht durch

**Sub-Area:** `*`

`vcs.ignore-link-targets` normiert ein auflösendes Link-Ziel auf Dateiname und
Anker. Zeigt ein Nachzug auf eine **andere** Datei mit **demselben** Namen
(`conventions.md`, `README.md`, `AGENTS.md`, `observation.md` …) oder auf eine,
die nur noch in BASE existiert, bleibt der Core gleich — die Aussage einer
`Accepted`-ADR kann sich so verschieben, ohne dass das Gate es sieht. Ableiter:
bei einem Nachzug in einer `Accepted`-ADR das neue Ziel im Review gegen das
alte lesen, nicht nur den Namen.
