# MR-048 — Ein Repo-Gate über eine Werkzeug-Datei prüft Wohlgeformtheit, nicht Anwesenheit (schärft MR-047)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon zählt `.claude/settings.json` und
  `.claude/hooks/*` zum Artefakt-Set der Durchsetzungsschicht
  ([`grundlagen-durchsetzungsschicht.md` §Das vollständige Artefakt-Set](../../.harness/baseline/v5.18.0/regelwerk/grundlagen-durchsetzungsschicht.md))
  und benennt daneben **drei** Bindepunkte, zu denen diese Artefakte
  ausdrücklich **nicht** gehören. Ob ein Repo trotzdem über sie wacht, sagt er
  nicht.
- **Datum:** 2026-08-27
- **Geltungsbereich:** jede Repo-Prüfung, deren Gegenstand eine Datei unter
  `.claude/` ist.
- **Adaption:** Ein Gate über eine Werkzeug-Datei sagt: **ist sie da, ist sie
  wohlgeformt und zeigt auf Vorhandenes.** Fehlt sie, überspringt die Prüfung.
  Sie fordert die Datei nicht ein.

  **Die Trennlinie in einem Satz:** eine Werkzeug-Einstellung **darf fehlen** —
  [`AGENTS.md`](../../AGENTS.md) §3.1 erklärt einen Lauf ohne dieses Werkzeug
  ausdrücklich für ungebunden. **Kaputt sein darf sie nicht**, denn eine
  committete Datei, die nichts mehr tut, ist eine Zusage ohne Deckung — und sie
  fällt genau dann aus, wenn niemand hinsieht.

  **Was daraus folgt und was nicht.** Die Prüfung deckt einen **Ausfallgrund**
  ab, nicht den Ausfall: dass das Werkzeug den Hook ruft und seine Antwort
  befolgt, steht in keiner Datei und kann kein Gate zusagen. Ein grüner Lauf ist
  deshalb kein Beleg für eine scharfe Durchsetzung.

  **Verhältnis zu [`MR-042`](../conventions.md#mr-042):** dessen Satz, ein Gate
  über den Wächter wäre *„eine Zusage, die ein Lauf ohne dieses Werkzeug nicht
  hält"*, gilt seinem **Verhalten** — dafür bleibt `make guard-probe` die
  werkzeug-lokale Probe. Das Skip-bei-Abwesenheit ist die Form, in der diese
  Warnung auch für die Datei-Prüfung eingelöst ist: ohne das Werkzeug wird
  nichts behauptet.
- **Begründung:** Ohne die Trennung stünde die Wahl zwischen zwei schlechten
  Enden. Fordert das Gate die Datei ein, macht es eine Werkzeug-Einstellung zur
  Repo-Pflicht und färbt jeden Klon rot, der ohne sie arbeitet — aus einem
  Grund, der das Produkt nicht berührt. Prüft es gar nichts, bleibt der Fall
  offen, an dem beide Durchsetzungs-Schichten zugleich still ausfallen
  ([`MR-047`](../conventions.md#mr-047)).
- **Auflösungs-Trigger:** der Kanon nimmt die Werkzeug-Artefakte in die
  Bindepunkte auf. Dann sind sie Repo-Pflicht, und die Skip-Hälfte entfällt.
