# Ein neu angelegter `.claude/agents/`-Typ löst nicht sofort in derselben Sitzung auf

**Sub-Area:** `*`

Wird eine neue `.claude/agents/<typ>.md`-Datei innerhalb einer laufenden
Sitzung angelegt, meldet ein sofortiger Aufruf mit `subagent_type: "<typ>"`
„Agent type '<typ>' not found" — die Discovery neuer Agent-Typen geschieht
offenbar nicht live, sondern verzögert (in derselben Sitzung, nach einer
gewissen Zeit, ohne erkennbares Neustart-Ereignis). Für den Implementer
bedeutet das: der erste Handoff an einen gerade erst angelegten Rollen-Typ
braucht einen generischen Agenten als Rückfall, und ein späterer Aufruf mit
demselben Typ kann bereits wieder funktionieren, ohne dass sich am Repo etwas
geändert hätte.
