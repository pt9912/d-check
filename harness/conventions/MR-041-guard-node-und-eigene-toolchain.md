# MR-041 — Auch `node` ist ein Host-Interpreter; der Wächter selbst bleibt die benannte Ausnahme (schärft MR-040)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Fortschreibung der Härtung aus
  [`MR-040`](../conventions.md#mr-040) in der Form, die
  [`modul-13-quality-gates.md` §Guard-Härtung](../../.harness/baseline/v5.12.0/regelwerk/modul-13-quality-gates.md)
  vorgibt: jede Härtung ein neuer Eintrag, nie eine Änderung am akzeptierten.
- **Datum:** 2026-08-26
- **Geltungsbereich:** [`.claude/hooks/pretooluse-command-guard.sh`](../../.claude/hooks/pretooluse-command-guard.sh)
  und [`AGENTS.md`](../../AGENTS.md) §3.1.
- **Adaption:** Der Wächter blockiert zusätzlich `node`, `deno` und `bun` in
  Befehlsposition. Der **Hook** ruft `node` weiterhin: Hook-Subprozesse laufen
  nicht durch den Hook.
- **Begründung:** [`MR-040`](../conventions.md#mr-040) nahm `node` aus, mit der
  Begründung, ein Verbot wäre *„eine Regel, die ihr eigenes Werkzeug verböte"*.
  **Diese Begründung trägt nicht.** Der Wächter ruft `node` aus dem Hook heraus,
  nicht über das Bash-Werkzeug des Agenten; ein Verbot für Agenten-Befehle hätte
  die eigene Prüfung nie berührt. Bis zur Korrektur stand damit ein Schlupfloch
  offen, das genauso mächtig war wie das geschlossene.

  **Was bleibt, ist die tiefere Inkonsistenz — und sie ist hier benannt, nicht
  gelöst:** Der Wächter soll erzwingen, dass nur die in
  [`AGENTS.md`](../../AGENTS.md) §3.1 genannte Host-Klasse benutzt wird, und ist
  selbst in einer Toolchain geschrieben, die nicht dazugehört. Eine Regel, deren
  Durchsetzung außerhalb ihrer eigenen Klasse steht, ist nur so lange
  glaubwürdig, wie das dasteht.
- **Auflösungs-Trigger:** der Wächter läuft mit `bash` und den
  POSIX-Werkzeugen; die Host-Abhängigkeit `node` entfällt dann ersatzlos.
