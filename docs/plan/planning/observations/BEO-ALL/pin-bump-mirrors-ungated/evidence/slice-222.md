**Vorgang:** slice-222
**Fund:** **Beide gate-blinden Richtungen traten in derselben Hebung ein** —
vergessene Hebung *und* Über-Hebung, obwohl der Eintrag genau davor warnt und
der Slice-Plan ihn in DoD (2) zitiert.

**Über-Hebung:** Ein `sed` über die lebenden Dateien setzte den
`Geltungsbereich` von [`MR-067`](../../../../../../../harness/conventions.md#mr-067) auf `v6.6.0` — der Eintrag **ist** aber die
Hebung *auf v6.5.0*. Dasselbe in seiner Index-Zeile. Eine
Vergangenheits-Aussage wurde damit still falsch, und kein Gate meldet das:
Der Pfad existiert ja.

**Vergessene Hebung, durch eine übersehene Form:**
`.harness/skills/reviewer.md` schreibt `../baseline/v6.5.0/` <!-- d-check:ignore (der Baum ist mit slice-222 entfernt) --> **relativ**. Das
Suchmuster war `.harness/baseline/v6.5.0/` und traf es nicht. Gefunden hat es
erst der `links`-Befund nach dem Entfernen des alten Baums — also der Zufall,
dass die Ziele verschwanden, nicht die Suche.

**Der Ableiter ist eine Änderung der Suchform:** Nach der **Version** suchen,
nicht nach dem **Pfad**. Erst `grep -rhoE '[^ ("]*baseline/v6\.5\.0'` mit
Gruppierung nach Präfix machte die fünf Schreibweisen sichtbar — vier
Pfad-Formen unterschiedlicher Tiefe plus die relative.

**Die dritte Klasse hielt:** Die fünf Release-/Tree-URLs wurden gefunden, weil
sie nach dem Pfad-Retarget als einzige Inkonsistenz übrig blieben — eine
Datei sagte in Zeile 47 `v6.6.0` und in Zeile 41 `v6.5.0`. Auch das war
Augenschein, kein Sensor.
