# MR-046 — Es gibt keine Skript-Stage, weil kein Fall sie verlangt (schärft MR-040)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon verlangt, dass eine Regel auf
  Vorhandenes zeigt, nicht auf Vorgesehenes
  ([`modul-13-quality-gates.md` §Guard-Härtung](../../.harness/baseline/v6.0.0/regelwerk/modul-13-quality-gates.md)
  — eine Doku, die mehr verspricht als der Bestand trägt, ist eine
  Harness-Lüge); *ob* ein Repo eine Skript-Toolchain führt, ist seine Sache.
- **Datum:** 2026-08-27
- **Geltungsbereich:** [`AGENTS.md`](../../AGENTS.md) §3.1 und der
  [`Dockerfile`](../../Dockerfile).
- **Adaption:** Dieses Repo führt **keine** Dockerfile-Stage für Skripte, und
  [`AGENTS.md`](../../AGENTS.md) §3.1 verweist nicht mehr auf eine. Wo ein
  Skript nötig ist, trägt es `bash` mit den POSIX-Werkzeugen; darüber hinaus
  gibt es das Produkt und die Datei-Werkzeuge des Agenten.

  **Der Bedarf ist an den Fällen gemessen, nicht am Bestand.** Das ist der
  Unterschied, auf den es hier ankommt: der Bestand kann den Bedarf gar nicht
  zeigen, weil ein Interpreter-Aufruf seit
  [`MR-040`](../conventions.md#mr-040) nicht mehr eingecheckt werden dürfte.
  Gezählt sind deshalb die Fälle, in denen tatsächlich ein Skript lief oder
  laufen sollte — neun, jeder mit Zuordnung:

  | Fall | Zuordnung |
  |---|---|
  | Datei-Änderungen per Interpreter (der Anlass von [`MR-040`](../conventions.md#mr-040)) | **Datei-Werkzeug** des Agenten |
  | Ad-hoc-Zählungen über Dokumente | **Produkt**, sonst `grep -c` |
  | Referenz-/Anker-Prüfungen | **Produkt** (`make doc-check`) |
  | JSON-Extraktion in der Hook-Eingabe | **POSIX** (`awk`-Scanner) |
  | Proben-Harnisch des Wächters | **POSIX** (`bash`) |
  | Der Wächter selbst | **POSIX** (`bash` + `awk`) |
  | Zeilen-/Muster-Zählungen im Arbeitsverlauf | **POSIX** (`grep -c`, `awk`) |
  | Komplexe Suchmuster, die der Wächter quote-blind blockt | **POSIX** (`grep -f` mit Musterdatei) |
  | Klammer-Balance einer `.json` prüfen | **offener Rest**, siehe unten |

  **Ein Rest ist geblieben, und er trägt die Stage nicht.** Für die letzte Zeile
  gibt es kein Werkzeug in der Klasse: geprüft wurde mit einem von Hand
  geschriebenen `awk`-Zähler über Klammern und Zeichenketten — also genau der
  „Parser durch die Hintertür", vor dem der Extraktor-Kommentar warnt. Er
  beantwortet **Balance**, nicht **Gültigkeit**. Das rechtfertigt keine
  Skript-Toolchain: eine Prüfung der Konfigurations-Dateien gehört an ein
  Gate, nicht an eine Shell — und dass sie fehlt, ist ein eigener Befund, kein
  Argument für eine Stage.

  **Der Bestand stützt das, ohne es zu beweisen:** über den gesamten getrackten
  Baum kommt kein Aufruf von `python`/`perl`/`ruby`/`node`/`deno`/`bun`/`uv`
  vor; die Fundstellen sind die Sperrliste des Wächters, seine Proben, die
  Berechtigungsliste und die Dokumentation. Vierzehn Dateien tragen einen
  bash-Shebang (zwölf `*.sh`, zwei unter `.githooks/`), zehn davon sind
  ausführbar gesetzt; die vier übrigen ruft das `Makefile` als `bash <pfad>`.

  **Benannte Grenze der Bestands-Messung:** die `uses:`-Schritte der
  CI-Workflows sind fremde JavaScript-Actions und laufen beim Runner unter
  `node`. Das ist keine Host-Abhängigkeit im Sinn von
  [`AGENTS.md`](../../AGENTS.md) §3.1 — die Regel gilt der Werkbank dieses
  Repos, nicht der Laufzeit eines fremden Runners — und aus dem Klon heraus
  ohnehin nicht prüfbar.

  **Tritt später ein Fall auf, ist er ein Entscheid.** Eine vierte Toolchain
  entsteht nicht nebenbei: sie bräuchte einen digest-gepinnten Pin, ein
  `make`-Target, die Deklaration in
  [`AGENTS.md`](../../AGENTS.md) §4 und
  [`harness/README.md`](../README.md), und sie zöge die drei Pin-Spiegel-Klassen
  aus [`BEO-008`](../../docs/plan/planning/observations.md) nach sich.

  **Und die Form wäre offen, nicht vorgegeben.** Die zwei jüngsten
  Fremd-Toolchains dieses Repos sind **keine** Dockerfile-Stage, sondern ein
  digest-gepinntes **externes Image**, gefahren aus einem Skript bzw. einem
  Include: `semgrep` ([ADR-0010](../../docs/plan/adr/0010-semgrep-hermetisches-gate.md))
  und `a-check` ([ADR-0029](../../docs/plan/adr/0029-arch-check-via-a-check.md)).
  Die Stage-Formulierung in [`MR-040`](../conventions.md#mr-040) beschreibt
  damit **eine** mögliche Form, nicht die gelebte Regel.
- **Begründung:** [`MR-040`](../conventions.md#mr-040) benennt richtig, **wohin**
  ein Skript gehörte, das die erlaubte Klasse sprengt. `AGENTS.md` §3.1 hat
  daraus einen Weg gemacht, den es gibt — und damit auf etwas verwiesen, das
  nicht existiert. Wer ihn gehen wollte, hätte nichts gefunden; wer die Regel
  las, hielt die Frage für beantwortet. Beides ist schlimmer als die offene
  Aussage, dass der Fall bisher nicht auftrat.
- **Auflösungs-Trigger:** die **Slice-Planung** ist der Anlass — wer einen Slice
  schneidet, der eine Fremd-Toolchain braucht, prüft diesen Eintrag mit und
  löst ihn durch einen Nachfolge-Eintrag ab, der die gewählte Form nennt.
  Bis dahin gilt er.
