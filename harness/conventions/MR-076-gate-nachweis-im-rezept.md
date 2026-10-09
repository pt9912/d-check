# MR-076 — Der Gate-Nachweis entsteht im Rezept von `gates`, und der Stop-Hook blockt, wo er den Zustand nicht liest (schärft MR-004)

- **Datum:** 2026-10-08
- **Geltungsbereich:** das Target `gates` im [`Makefile`](../../Makefile),
  [`record-gates.sh`](../../tools/harness/record-gates.sh), der Stop-Hook
  [`stop-require-gates.sh`](../../.claude/hooks/stop-require-gates.sh)
- **Ersetzt-Baseline-Regel:** — *(keine; schärft
  [`MR-004`](../conventions.md#mr-004), dessen Mechanik den Nachweis als
  letzten Prerequisite von `gates` beschreibt)*
- **Adaption:** `make gates` schreibt den Nachweis im eigenen **Rezept**, nach
  allen Gliedern; das Target `record-gates` bleibt als Werkzeug. Ein Rezept
  läuft erst, wenn alle Prerequisites grün sind — auch unter `make -k`. Als
  Prerequisite lief der Nachweis unter `make -k` trotz rotem Glied, denn `-k`
  arbeitet nach einem Fehler die übrigen Prerequisites weiter ab, und der
  Stop-Hook gab danach frei. Unter `make -i`, das jeden Fehler ignoriert,
  schreibt das Rezept keinen Nachweis und meldet das (`make` endet dort
  trotzdem mit 0; der Stop-Hook gibt nur frei, wenn ein früherer grüner Lauf
  denselben Inhalt belegt); Erkennung und Schreiben
  stehen in einer Rezeptzeile, weil `-i` auch den Abbruch einer eigenen
  Prüfzeile ignoriert. Die Festlegung des Handoff-Gates steht in
  [`SPEC-094`](../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge).
  **Stop-Hook:** Lässt sich der Inhalts-Hash nicht berechnen, der Nachweis
  nicht lesen oder `git status` nicht ausführen — außerhalb eines
  Repositorys, bei einem kaputten Index —, blockt der Hook mit Grund. Vorher
  endete er im ersten Fall ohne Antwort und gab im letzten frei; beides ließ
  den Stop durch.
- **Grenze:** Wer `record-gates` von Hand ruft, schreibt einen Nachweis ohne
  Gate-Lauf — das Target urteilt nicht, es schreibt. Der Wächter ist ein
  Stolperdraht gegen „fertig ohne Gate-Lauf", keine Sperre gegen Absicht —
  ebenso `make -i gates MAKEFLAGS=`, das die Erkennung von `-i` leert.
- **Begründung:** Gemessen mit `make -k gates THRESHOLD=99`: `coverage-gate`
  rot, Exit 2 — mit dem Prerequisite entstand
  `.harness/state/gates-passed.diffsha`, mit dem Rezept nicht. Der Stop-Hook
  in einem Wegwerf-Klon: unlesbare Datei, unlesbarer Nachweis, kaputter Index
  und ein Lauf außerhalb des Repositorys blocken, ein sauberer Baum ohne
  Nachweis und ein passender Nachweis geben frei.
- **Auflösungs-Trigger:** permanent — solange der Nachweis über `make`
  entsteht.
