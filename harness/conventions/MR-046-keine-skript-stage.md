# MR-046 — Es gibt keine Skript-Stage, weil kein Fall sie verlangt (schärft MR-040)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon verlangt, dass eine Regel auf
  Vorhandenes zeigt, nicht auf Vorgesehenes
  ([`modul-13-quality-gates.md` §Guard-Härtung](../../.harness/baseline/v5.12.0/regelwerk/modul-13-quality-gates.md)
  — eine Doku, die mehr verspricht als der Bestand trägt, ist eine
  Harness-Lüge); *ob* ein Repo eine Skript-Toolchain führt, ist seine Sache.
- **Datum:** 2026-08-27
- **Geltungsbereich:** [`AGENTS.md`](../../AGENTS.md) §3.1 und der
  [`Dockerfile`](../../Dockerfile).
- **Adaption:** Dieses Repo führt **keine** Dockerfile-Stage für Skripte, und
  [`AGENTS.md`](../../AGENTS.md) §3.1 verweist nicht mehr auf eine. Wo ein
  Skript nötig ist, trägt es `bash` mit den POSIX-Werkzeugen; darüber hinaus
  gibt es das Produkt und die Datei-Werkzeuge des Agenten.

  **Gemessen statt angenommen:** Der Bestand kennt zwölf ausführbare Skripte,
  **alle** mit `#!/usr/bin/env bash`, dazu einen POSIX-`awk`-Extraktor. Ein
  Aufruf von `python`/`perl`/`ruby`/`node`/`deno`/`bun`/`uv` kommt in Skripten,
  `Makefile`, `a-check.mk`, den CI-Workflows und dem `Dockerfile` **nicht** vor
  — die einzigen Fundstellen sind die Sperrliste des Wächters selbst und seine
  Proben.

  **Tritt später ein Fall auf, ist er ein Entscheid.** Eine vierte Toolchain
  entsteht nicht nebenbei: sie bräuchte eine digest-gepinnte `FROM`-Zeile, ein
  `make`-Target, die Deklaration in
  [`AGENTS.md`](../../AGENTS.md) §4 und
  [`harness/README.md`](../README.md), und sie zöge die drei Pin-Spiegel-Klassen
  aus [`BEO-008`](../../docs/plan/planning/observations.md) nach sich.
- **Begründung:** [`MR-040`](../conventions.md#mr-040) benennt richtig, **wohin**
  ein Skript gehörte, das die erlaubte Klasse sprengt. `AGENTS.md` §3.1 hat
  daraus einen Weg gemacht, den es gibt — und damit auf etwas verwiesen, das
  nicht existiert. Wer ihn gehen wollte, hätte nichts gefunden; wer die Regel
  las, hielt die Frage für beantwortet. Beides ist schlimmer als die offene
  Aussage, dass der Fall bisher nicht auftrat.
- **Auflösungs-Trigger:** ein Fall tritt auf, der `bash`, die POSIX-Werkzeuge,
  das Produkt und die Datei-Werkzeuge übersteigt. Dann entsteht die Stage, und
  dieser Eintrag wird von ihr abgelöst.
