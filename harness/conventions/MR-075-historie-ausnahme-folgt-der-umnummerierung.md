# MR-075 — Die Historie-Ausnahme folgt der Umnummerierung der Spezifikation (Nachtrag zu MR-0098)

- **Datum:** 2026-10-08
- **Geltungsbereich:** `matrix.exclude-sections` und die `structure`-Regel der
  Spezifikations-Historie in der [`.d-check.yml`](../../.d-check.yml)
- **Ersetzt-Baseline-Regel:** — *(keine; Nachtrag zu
  [`MR-0098`](../conventions.md#mr-0098), dessen Ausnahme an eine
  Überschrift gebunden ist)*
- **Adaption:** Seit `v6.17.0` trägt die Spezifikation einen Abschnitt
  §7 *Festlegungen der Harness-Werkzeuge*, und ihre Historie ist §8. Die
  Ausnahme aus `MR-0098` gilt der Historie **beider** Spec-Straten und ist an
  die Überschrift gebunden; sie nennt deshalb jetzt zwei: `7. Historie` im
  Lastenheft und `8. Historie` in der Spezifikation. Die `structure`-Regel, die
  die Historie der Spezifikation absteigend nach Datum hält, zeigt auf
  `## 8. Historie`. Der Gegenstand der Ausnahme ist unverändert — die
  Historie als Zeitdokument —, nur ihr Name in einer der beiden Dateien.
- **Begründung:** Ohne den Nachzug fiele die Historie der Spezifikation still
  aus der Ausnahme, und die Matrix meldete ihre eingefrorenen Verweise; die
  `structure`-Regel liefe ins Leere. `MR-074` hatte diese Kopplung bei der
  Hebung benannt.
- **Auflösungs-Trigger:** wie `MR-0098` — ändert die Baseline die
  Historie-Semantik, werden beide aufgelöst.
