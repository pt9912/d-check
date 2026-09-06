# MR-047 — Eine zweite Durchsetzungs-Schicht neben dem Wächter (schärft MR-040)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon verlangt für die
  Durchsetzungsschicht, dass jede Härtung als neuer Eintrag landet
  ([`modul-13-quality-gates.md` §Guard-Härtung](../../.harness/baseline/v6.3.1/regelwerk/modul-13-quality-gates.md));
  *wie viele* Schichten ein Repo führt, ist seine Sache.
- **Datum:** 2026-08-27
- **Geltungsbereich:** [`.claude/settings.json`](../../.claude/settings.json)
  (`permissions.deny` / `permissions.ask`).
- **Adaption:** Neben dem Tool-Call-Wächter lehnt eine **Permission-Sperrliste**
  ab: die Host-Toolchain aus [`AGENTS.md`](../../AGENTS.md) §3.1 samt der
  uv-Runner, dazu Force-Push als `deny` und die destruktiven git-/docker-Formen
  als `ask`. Beide Schichten sind voneinander unabhängig — sie teilen keinen
  Code und keinen Zustand.

  **Warum überhaupt zwei:** der Wächter war der einzige Punkt, an dem die ganze
  Durchsetzung hing, und als er ausfiel, sagte niemand etwas
  ([`BEO-ALL/haertung-kippt-fehlerpolitik-ungeprueft`](../../docs/plan/planning/observations/BEO-ALL/haertung-kippt-fehlerpolitik-ungeprueft/observation.md)). Eine Schicht, deren
  Ausfall unsichtbar ist, ist keine Zusage.

  **Verhältnis zu [`MR-040`](../conventions.md#mr-040):** dessen Satz *„Die
  Berechtigungsliste führt keine Python-Einträge mehr"* gilt der **Erlaubnis**-
  Hälfte und gilt dort unverändert weiter — `permissions.allow` führt keinen
  Interpreter. Die neuen Einträge stehen in der **Sperr**-Hälfte, die es damals
  nicht gab.
- **Grenzen — was diese Schicht NICHT kann.** Eine Permission-Regel matcht den
  **ganzen Befehl** ab dem Anfang. Damit fallen durch:

  | Klasse | Beispiel |
  |---|---|
  | Präfix-Formen | `sudo pip install x` |
  | absolute Pfade | `/usr/bin/pip install x` |
  | zusammengesetzte Kommandos | `echo a; pip install x` |
  | Versions-Suffixe | `python3.12 x.py` |
  | Flag hinter dem Positional | `git push origin main --force` |
  | Langform statt Kurzflag | `git clean --force` |
  | Regel ohne Folge-Token | `git reset --hard` (ohne Argument) |

  Für die Interpreter-Hälfte deckt der Wächter die ersten vier — **die
  git-/docker-Hälfte ist einschichtig**, denn seine Sperrliste kennt weder
  `git` noch `docker`. Dort gilt die Aufzählung, nicht die Klasse: `git branch
  -D`, `git checkout -- .`, `git reflog expire`, `git filter-branch`,
  `docker rmi`, `docker volume prune` und ihresgleichen haben **keinen**
  Eintrag.

  **Und beide Schichten sind über `make` konstruktionsbedingt durchlässig:** die
  Regel sieht `make clean`, der Wächter sieht `make` in Befehlsposition,
  ausgeführt wird `docker image rm`. Das ist gewollt — `make` ist der
  sanktionierte Einstieg —, heißt aber, dass die `ask`-Einträge nur **direkte**
  Aufrufe binden.

  **Keine wiederholbare Probe.** `make guard-probe` fährt den Wächter; die
  Regeln dieser Schicht wertet das Werkzeug aus, nicht ein Skript, und ein
  Proben-Harnisch dafür existiert nicht. Belegt ist **eine** Regel, einmalig:
  ein Force-Push auf ein nicht existierendes Remote wurde abgelehnt, und der
  Wächter kennt `git` nicht — die Ablehnung kann also nur von hier kommen.
  Damit ist die Unabhängigkeit der Schicht gezeigt, ihre Vollständigkeit nicht.
- **Begründung:** Zwei Schichten kosten eine Datei und decken verschiedene
  Gestalten. Die eine liest den Befehl syntaktisch und findet ihn auch hinter
  `sudo`, in einer Sub-Shell oder nach einem Trenner; die andere hängt an keinem
  Hook und überlebt dessen Ausfall. Keine ersetzt die andere, und keine ist eine
  Sandbox.
- **Auflösungs-Trigger:** das Werkzeug führt keine Permission-Regeln mehr, oder
  eine der beiden Schichten wird abgeschafft. Beides ist ein Entscheid und ein
  neuer Eintrag.
