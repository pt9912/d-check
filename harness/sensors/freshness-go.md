# `make freshness-go` u. a. — meldet, ob upstream ein neuerer Release existiert als der gepinnte

## Vertrag

**Vier Versions-Achsen über _einen_** parametrierten Sensor
([`pin-freshness.sh`](../../tools/harness/pin-freshness.sh)): die
Toolchain-Version gegen eine **Sonderquelle** — golang/go publiziert keine
Release-Objekte, der `releases/latest`-Pfad liefe ins Leere —, dazu
`GOLANGCI_LINT_VERSION`, `SEMGREP_VERSION` und `A_CHECK_VERSION` gegen den
`releases/latest`-Redirect.

**Der Vergleich ist Gleich/Ungleich, kein Semver-Sort:** Alle vier Reihen sind
monoton, ein „neuer, aber älter" existiert dort nicht. Ein führendes Präfix
wird **symmetrisch auf beiden Seiten** gestrippt — das macht den Vergleich
nicht großzügiger, sondern richtig, weil die Pins es uneinheitlich führen
(`v2.13.1`, aber `1.175.0`).

## Grenze — was das Grün nicht abdeckt

1. **Fail-open mit Zeitgrenzen** — jede fremde Störung endet als `SKIP` mit
   Exit 0. Ein Sensor, der bei fremder Störung rot wird, wird abgeschaltet.
2. **Eine Versions-Achse sagt nichts über den Digest** — derselbe Tag kann neu
   gebaut sein; das ist die Frage der Digest-Achsen
   ([`runtime-base-digest`](runtime-base-digest.md)), und dass sie eine andere
   ist, ist **gemessen**, nicht vermutet.
3. **Gemeldet wird, nicht gehoben** — die Hebung bleibt ein bewusster Akt, und
   bei einer Toolchain-Hebung zieht das `golangci`-Pendant mit.

**Netzlos prüfbar** über `--compare <name> <gepinnt> <upstream>`; ohne diesen
Einstieg wäre die Semantik nur mit Netz zu prüfen und damit gar nicht.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | Pin aktuell — oder `SKIP` (fail-open) |
| 3 | veraltet |

Das ist der Code des **Skripts**; `make` normalisiert auf seinen eigenen.

## Bindung

kein Gate, bewusst **nicht** in `gates`; Netz, gerufen vom Nachtlauf.
