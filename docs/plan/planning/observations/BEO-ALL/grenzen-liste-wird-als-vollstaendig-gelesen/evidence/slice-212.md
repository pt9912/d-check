**Vorgang:** slice-212
**Fund:** Der `## Grenze`-Abschnitt von `harness/sensors/baseline-verify.md`
führte **vier** Einträge — alle vier über die Alias-Auflösung — und **nicht**
den, der alle überwiegt: Der Lauf beweist **innere Konsistenz, nicht
Echtheit**; ein mitgezogenes `SHA256SUMS` passiert grün. **Gefunden hat es der
Auftraggeber** beim Lesen eines Zwischenbescheids, nicht ein Review und kein
Gate. Die Inventur über alle 24 Sensor-Dateien fand **sechs weitere** —
`adr-check`, `doc-check`, `lint`, `arch-check`, `semgrep`, `review-coverage` —
und in **allen sieben** stand die fehlende Grenze bereits im **Vertrags**-Teil
derselben Datei oder in ihrer Konfiguration. Das ist der Mechanismus, den der
Eintrag beschreibt, an sieben unabhängigen Dateien gemessen.
**Und die zweite Stufe des Ableiters stammt aus demselben Lauf:** Zwei der
Ergänzungen waren aus dem Vertrags-Teil abgeleitet und **falsch** — bei
`review-coverage` widerspricht der Code-Kommentar dem Verhalten. Wer den
Vertrag umdreht, erbt seinen Fehler.
