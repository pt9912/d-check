# MR-043 — Der Werkzeug-Einstieg importiert `AGENTS.md`, statt auf ihn zu verweisen

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon verlangt von der
  Wurzel-Einstiegsdatei genau dies —
  [`grundlagen-durchsetzungsschicht.md` §Das vollständige Artefakt-Set](../../.harness/baseline/v6.3.1/regelwerk/grundlagen-durchsetzungsschicht.md):   <!-- d-check:cite .harness/baseline/v6.3.1/regelwerk/grundlagen-durchsetzungsschicht.md:100-102 -->
  *„Sie bringt `AGENTS.md` in den Lauf-Kontext, wo Modul 9 es für jeden Lauf
  verlangt. Sie **verweist** dorthin und legt nichts fest"* — die Fettung im
  zweiten Satz steht so in der Quelle, das tragende Verb *bringt* nicht; die
  Beweislast dieses Eintrags hängt an ihm und wird hier daneben genannt statt im
  Zitat gesetzt. Womit sie das tut, ist werkzeugabhängig und bleibt offen; die
  Form-Frage tritt die Rangliste an diesen Speicher ab.
- **Datum:** 2026-08-27
- **Geltungsbereich:** [`CLAUDE.md`](../../CLAUDE.md).
- **Adaption:** [`CLAUDE.md`](../../CLAUDE.md) führt `AGENTS.md` als
  `@`-Import — die Import-Syntax von Claude Code, die den Inhalt der Datei in
  den Kontext lädt —, nicht als Markdown-Link mit Lese-Aufforderung.

  **Die Grenze des Imports gehört mitgelesen:** er lädt genau eine Datei.
  Alles, wohin [`AGENTS.md`](../../AGENTS.md) weiterroutet — dieser Speicher,
  [`harness/README.md`](../README.md), die Spec-Straten, die Slices —, steht als
  Markdown-Link da und wird **nicht** mitgeladen. Der Import ersetzt das Lesen
  also nicht, er stellt nur sicher, dass die Hard Rules und der Routing-Pfad im
  Kontext stehen, bevor der erste Werkzeugaufruf passiert.

  **Kein zweiter Kandidat.** [`CLAUDE.md`](../../CLAUDE.md) importiert
  ausschließlich [`AGENTS.md`](../../AGENTS.md). Jede weitere importierte Datei
  wäre ein Pflichtanteil am Kontext **jedes** Laufs — genau die Kosten, mit
  denen der Kanon den Ein-Eintrag-je-Datei-Schnitt dieses Speichers begründet.
- **Begründung:** Ein Verweis, den ein Lauf **befolgen muss**, ist eine Bitte;
  der Kanon verlangt aber, dass die Einstiegsdatei die Hard Rules *bringt*. Der
  Unterschied ist nicht theoretisch: Ein Lauf, der die Regel nicht im Kontext
  hat, verletzt sie nicht aus Widerspruch, sondern weil er sie nicht kennt —
  und merkt es erst, wenn ihn jemand darauf hinweist. Die Bitte trägt genau so
  weit, wie der Lauf ihr folgt; der Import trägt unabhängig davon.
- **Auflösungs-Trigger:** das Werkzeug lädt die Einstiegsdatei nicht mehr, oder
  seine Import-Syntax ändert sich. Dann ist die Form neu zu bestimmen — die
  Anforderung des Kanons bleibt.
