**Vorgang:** slice-235
**Fund:** Direkt nach dem Anlegen von `.claude/agents/reviewer.md` scheiterte
ein Review-Aufruf mit `subagent_type: "reviewer"` mit „Agent type 'reviewer'
not found. Available agents: claude, claude-code-guide, Explore,
general-purpose, Plan, statusline-setup". Der Review lief stattdessen unter
einem generischen Typ. Ein späterer Aufruf in **derselben** Sitzung (R2, ohne
weitere Repo-Änderung an den Agent-Dateien dazwischen außer dem eigentlichen
Fix) löste `subagent_type: "reviewer"` bereits erfolgreich auf — ein
System-Hinweis kündigte die neuen Typen zwischenzeitlich von sich aus an.
