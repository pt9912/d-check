# Ein dünn angelegtes Werkzeug-Artefakt wächst beim nächsten Anfassen wieder

**Sub-Area:** `*`

Eine bewusst schmal gehaltene Datei — etwa ein `.claude/agents/<rolle>.md`,
das nur auf seine Skill-/Regelwerk-Quelle verweist statt Inhalt zu
duplizieren — bleibt schmal nur, solange niemand sie anfasst. In
Schwester-Repos (`ai-harness-init`, `pg-change-feed`) sind vergleichbare
Rollen-Dateien im Lauf der Zeit auf 17–28 KB gewachsen: Herleitungen,
konkrete Zahlen und einzelne Kandidatenlauf-Notizen sammeln sich an, wenn ein
Implementer beim nächsten Fix „nur kurz" etwas ergänzt, statt die Ergänzung
an ihren gerankten Ort zu schreiben. Anders als bei `AGENTS.md`
([`BEO-ALL/briefing-datei-ueberschreitet-lade-budget`](../briefing-datei-ueberschreitet-lade-budget/observation.md))
trägt diese Datei niemanden Steering-Loop-Regeln — der Treiber ist nicht
Verkörperung, sondern gewöhnliche Pflege-Trägheit.
