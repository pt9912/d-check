**Vorgang:** slice-215
**Fund:** **Eine fail-open-Zusage aus genau einer Ausfall-Form.** Der Slice
änderte einen Exit-Code und rechtfertigte das mit der Messung *„ein Netz-/
API-Ausfall landet im `skip`-Zweig"*. Gemessen war **eine** Form: ein
unerreichbarer API-Host. Der Satz stand danach in **drei** Artefakten —
Skript-Kopf, Sensor-Grenze, Commit-Botschaft.

**Die N+1-te Form widerlegt ihn.** Ein **abgebrochener** Transfer liefert
Teildaten *und* einen Fehlerstatus (gemessen: 35 010 von 322 579 Bytes,
`curl`-Exit 28). Der Code wertete nur den Pipeline-Ausgang, hielt die
abgeschnittene Liste für vollständig und meldete **Exit 3 auf einen
Netzausfall** — also genau das falsche Rot, dessen Ausbleiben die Zusage
behauptete.

**Der Ableiter des Eintrags ist wörtlich das, was gefehlt hat:** *„Prüfe den
Schluss gegen die Proben-Menge, nicht gegen die Proben: suche die N+1-te
Form."* Der Slice-Plan hatte das Risiko in §6 sogar **vorab benannt** — und
der Lauf erklärte es für ausgeräumt, nachdem er eine Form gemessen hatte.
