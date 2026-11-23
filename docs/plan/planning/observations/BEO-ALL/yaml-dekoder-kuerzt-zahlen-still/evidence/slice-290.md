**Vorgang:** slice-290
**Fund:** `linelength` liest `max-chars` nur als YAML-Ganzzahl und die Schalter nur als `true`/`false`, weil yaml.v3 sonst kürzt; `file[].max-lines`, `cell-max-chars` und Wahrheitswerte anderer Module nehmen `80.5` bzw. `yes` weiter still an (Code-Review R1 (im Archiv von welle-95) F-5). Außerhalb der Welle nicht behoben.
