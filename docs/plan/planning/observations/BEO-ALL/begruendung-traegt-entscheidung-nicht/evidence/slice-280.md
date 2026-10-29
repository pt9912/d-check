**Vorgang:** slice-280
**Fund:** Die Verdrängung eines Unix-Treffers innerhalb eines Variablen-Treffers war
richtig; ihre Begründung in Plan, Test-Kommentar und ADR („vor dem Präfix steht das `E`
von `HOME`“) galt nur für `$HOME/…`. Bei `${HOME}/…` und `"$HOME"/…` steht `}` bzw. `"`
davor, das Unix-Muster träfe dort wirklich, und der Test fuhr nur die Form, in der die
Verdrängung nichts bewirkt. Ein Mutationstest, der die Verdrängung entfernte und grün
blieb, deckte es auf; der Review (F-1) nannte die beiden anderen Formen.
