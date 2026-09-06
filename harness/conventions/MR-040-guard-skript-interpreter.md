# MR-040 — Der Tool-Call-Wächter blockiert auch Host-Skript-Interpreter (schärft MR-005)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Die Baseline verlangt für die
  Durchsetzungsschicht, dass gegen eine **beobachtete** Umgehung gehärtet wird
  und jede Härtung als neuer Eintrag landet
  ([`modul-13-quality-gates.md` §Guard-Härtung](../../.harness/baseline/v6.3.1/regelwerk/modul-13-quality-gates.md));
  *welche* Interpreter ein Repo führt, ist seine Sache.
- **Datum:** 2026-08-26
- **Geltungsbereich:** [`.claude/hooks/pretooluse-command-guard.sh`](../../.claude/hooks/pretooluse-command-guard.sh),
  [`AGENTS.md`](../../AGENTS.md) §3.1 und die Berechtigungsliste unter
  `.claude/`.
- **Adaption:** Der Wächter blockiert zusätzlich `python`/`python3` samt
  Versions-Suffix (als Muster, weil eine Liste dort unvollständig wäre), `perl`
  und `ruby` in Befehlsposition — auch rekursiv in Sub-Shell-Strings. Die
  Berechtigungsliste führt keine Python-Einträge mehr.

  **Warum der Skopus und nicht die Sprache:** [`AGENTS.md`](../../AGENTS.md)
  §3.1 nennt als Host-Klasse `git`, GNU `make`, `bash`, Docker und die
  POSIX-Standardwerkzeuge, die die Gate-Skripte rufen. Ein
  General-Purpose-Interpreter gehört nicht dazu und **hebelt die Klasse aus**:
  er kann alles, was die genannten Werkzeuge können, ohne deren Grenze zu erben.
  Wo ein Skript wirklich nötig ist, gehört es als eigene Dockerfile-Stage
  gepinnt — wie jede andere Toolchain dieses Repos.

  **`node` steht bewusst nicht auf der Liste.** Der Wächter führt seine eigene
  Prüfung damit aus und deklariert es fail-closed als Host-Abhängigkeit; ein
  Verbot wäre eine Regel, die ihr eigenes Werkzeug verböte. Das ist eine
  benannte Inkonsistenz, kein Versehen.

  **Benannte Falsch-Positiv-Klasse, gemessen beim Scharfschalten:** Die
  Segmentierung trennt auch an `(`. Ein Kommando, das eine Zeichenkette wie
  `Bash(python3 …)` als *Daten* trägt — etwa ein `sed`-Muster über die
  Berechtigungsliste —, wird deshalb blockiert, obwohl in Befehlsposition kein
  Interpreter steht. Der erste Versuch, die Liste per `sed` zu bereinigen, ist
  genau daran gescheitert. Das ist der Preis der konservativen Segmentierung und
  bleibt so: ein Wächter, der Daten von Befehlen sicher unterscheiden wollte,
  wäre ein Parser.
- **Begründung:** Ein Stolperdraht, der die drei häufigsten Interpreter nicht
  kennt, schützt vor Tippfehlern und nicht vor Gewohnheit. Die Gewohnheit war
  der reale Fall: ein ganzer Arbeitstag mit Host-Python für Datei-Edits und
  Ad-hoc-Messungen, ohne dass ein Gate daran hing — der Schaden war nicht die
  Reproduzierbarkeit, sondern **nachgebaute Lexik**: mehrere dieser Messungen
  erwiesen sich als unzuverlässig, weil sie das Produkt imitierten statt es zu
  benutzen.
- **Auflösungs-Trigger:** permanent für die Klasse; die konkrete Liste wächst
  mit jedem beobachteten Fall.
