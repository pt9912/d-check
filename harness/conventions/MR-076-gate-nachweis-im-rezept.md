# MR-076 — Der Gate-Nachweis entsteht im Rezept von `gates`, nicht als Prerequisite (schärft MR-004)

- **Datum:** 2026-10-08
- **Geltungsbereich:** das Target `gates` im [`Makefile`](../../Makefile),
  [`record-gates.sh`](../../tools/harness/record-gates.sh)
- **Ersetzt-Baseline-Regel:** — *(keine; schärft
  [`MR-004`](../conventions.md#mr-004), dessen Mechanik den Nachweis als
  letzten Prerequisite von `gates` beschreibt)*
- **Adaption:** `make gates` schreibt den Nachweis im eigenen **Rezept**, nach
  allen Gliedern; das Target `record-gates` bleibt als Werkzeug. Ein Rezept
  läuft erst, wenn alle Prerequisites grün sind — auch unter `make -k`. Als
  Prerequisite lief der Nachweis unter `make -k` trotz rotem Glied, denn `-k`
  arbeitet nach einem Fehler die übrigen Prerequisites weiter ab, und der
  Stop-Hook gab danach frei. Die Festlegung des Handoff-Gates steht in
  [`SPEC-094`](../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge).
- **Grenze:** Wer `record-gates` von Hand ruft, schreibt einen Nachweis ohne
  Gate-Lauf — das Target urteilt nicht, es schreibt. Der Wächter ist ein
  Stolperdraht gegen „fertig ohne Gate-Lauf", keine Sperre gegen Absicht.
- **Begründung:** Gemessen mit `make -k gates THRESHOLD=99`: `coverage-gate`
  rot, Exit 2 — mit dem Prerequisite entstand
  `.harness/state/gates-passed.diffsha`, mit dem Rezept nicht.
- **Auflösungs-Trigger:** permanent — solange der Nachweis über `make`
  entsteht.
