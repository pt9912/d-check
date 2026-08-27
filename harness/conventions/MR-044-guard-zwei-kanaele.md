# MR-044 — Der Wächter blockt über zwei Kanäle, und die Antwortform ist die aktuelle (schärft MR-042)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon verlangt die Härtung als neuen
  Eintrag ([`modul-13-quality-gates.md` §Guard-Härtung](../../.harness/baseline/v5.12.0/regelwerk/modul-13-quality-gates.md));
  *wie* der Guard seine Ablehnung übermittelt, sagt er nicht — die Form-Frage
  tritt die Rangliste an diesen Speicher ab.
- **Datum:** 2026-08-27
- **Geltungsbereich:**
  [`.claude/hooks/pretooluse-command-guard.sh`](../../.claude/hooks/pretooluse-command-guard.sh)
  und [`tools/harness/guard-probe.sh`](../../tools/harness/guard-probe.sh).
- **Adaption:** Jeder Block läuft über **beide** Kanäle: die JSON-Antwort in der
  aktuellen Form (`hookSpecificOutput.permissionDecision` = `deny`, Grund in
  `permissionDecisionReason`) **und** einen Nicht-Null-Exit. Die Proben führen
  dafür ein eigenes Verdikt `halb` — eine Ablehnung, der der Exit-Riegel fehlt,
  ist kein `block`.

  **Beide Formen wirken heute, gemessen:** die veraltete Top-Level-Form
  (`decision`/`reason`) über die Abwärtskompatibilitäts-Abbildung, die aktuelle
  direkt, der Nicht-Null-Exit unabhängig von beiden. Zusammen sind sie
  konfliktfrei — der Grund aus der JSON-Antwort erscheint beim Aufrufer, der
  Exit liegt darunter.

  **Warum trotzdem beide und nicht die eine, die reicht:** ein einzelner Kanal
  macht jede Ablehnung von einer Auslegung abhängig, die dieser Guard nicht
  kontrolliert. Die veraltete Form hängt zusätzlich an einer Abbildung, die
  ausdrücklich als Übergang deklariert ist. Zwei Kanäle kosten eine Zeile.
- **Begründung:** Die Frage war nicht akademisch. Der Guard hat einen Tag lang
  **jede** Ablehnung emittiert und **keine** durchgesetzt — nicht wegen der
  Antwortform, sondern weil er vor ihr starb. Ein einziger Kanal macht den
  Unterschied zwischen „hat nicht geurteilt" und „wurde nicht gehört" von außen
  unsichtbar; der Exit-Riegel trennt die beiden Fälle und deckt den zweiten.

  **Was der Eintrag NICHT behauptet:** dass zwei Kanäle die Klasse schließen.
  Sie decken den Ausfall der Antwortform, nicht den Ausfall des Guards vor der
  Antwort. Dagegen hilft nur, dass der Guard nichts tut, was fehlschlagen kann,
  bevor er urteilt — und das ist eine Eigenschaft seines Codes, keine seiner
  Ausgabe.
- **Auflösungs-Trigger:** das Werkzeug führt nur noch **einen** Kanal, oder der
  Kanon schreibt die Form vor. Dann gilt seine.
