**Vorgang:** slice-212
**Fund:** Der `## Grenze`-Abschnitt von `harness/sensors/baseline-verify.md`
führte **vier** Einträge — alle vier über die Alias-Auflösung — und **nicht**
den, der alle überwiegt: Der Lauf beweist **innere Konsistenz, nicht
Echtheit**; ein mitgezogenes `SHA256SUMS` passiert grün. **Gefunden hat es der
Auftraggeber** beim Lesen eines Zwischenbescheids, nicht ein Review und kein
Gate. Die Inventur über alle 24 Sensor-Dateien fand **drei weitere**
(`semgrep`, `review-coverage`, `arch-check`) — und in **allen vier** stand die
fehlende Grenze bereits im **Vertrags**-Teil derselben Datei. Das ist der
Mechanismus, den der Eintrag beschreibt, an vier unabhängigen Dateien
gemessen.
