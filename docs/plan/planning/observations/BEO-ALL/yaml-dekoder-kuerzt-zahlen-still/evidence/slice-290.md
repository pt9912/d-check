**Vorgang:** slice-290
**Fund:** `linelength` liest `max-chars` nur als YAML-Ganzzahl und die Schalter nur als `true`/`false`, weil yaml.v3 sonst kürzt; `file[].max-lines`, `cell-max-chars` und Wahrheitswerte anderer Module nehmen `80.5` bzw. `yes` weiter still an ([Code-Review R1](../../../../../../../docs/reviews/2026-11-23-slice-290-code-r1.md) F-5). Außerhalb der Welle nicht behoben.
