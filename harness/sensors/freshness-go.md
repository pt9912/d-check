# `make freshness-go` u. a. — meldet, ob upstream ein neuerer Release existiert als der gepinnte

## Vertrag

**Fünf Versions-Achsen über _einen_** parametrierten Sensor
([`pin-freshness.sh`](../../tools/harness/pin-freshness.sh)): die
Toolchain-Version gegen eine **Sonderquelle** — golang/go publiziert keine
Release-Objekte, der `releases/latest`-Pfad liefe ins Leere —, dazu
`GOLANGCI_LINT_VERSION`, `SEMGREP_VERSION`, `A_CHECK_VERSION` und
`TRIVY_VERSION` gegen den `releases/latest`-Redirect. Die drei Action-Pins
derselben Mechanik beschreibt [`checkout-pin-freshness`](checkout-pin-freshness.md).

Quellen, Präfix-Behandlung, Vergleich (Gleich/Ungleich, keine
Versions-Ordnung), fail-open und Exit legt
[`SPEC-098`](../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)
fest; diese Datei führt sie nicht ein zweites Mal.

## Grenze — was das Grün nicht abdeckt

1. **Fail-open mit Zeitgrenzen** — jede fremde Störung endet als `SKIP` mit
   Exit 0. Ein Sensor, der bei fremder Störung rot wird, wird abgeschaltet.
2. **Eine Versions-Achse sagt nichts über den Digest** — derselbe Tag kann neu
   gebaut sein; das ist die Frage der Digest-Achsen
   ([`runtime-base-digest`](runtime-base-digest.md)), und dass sie eine andere
   ist, ist **gemessen**, nicht vermutet.
3. **Gemeldet wird, nicht gehoben** — die Hebung bleibt ein bewusster Akt, und
   bei einer Toolchain-Hebung zieht das `golangci`-Pendant mit.
4. **`freshness-trivy` läuft nicht im Nachtlauf** — `upstream-drift.yml` führt
   die anderen vier Achsen; der Scanner-Pin wird nur gemeldet, wenn jemand die Achse
   selbst ruft.

**Netzlos prüfbar** ist der Vergleich über `--compare <name> <gepinnt>
<upstream>` — die Präfix-Behandlung der Quell-Zweige nicht, sie gehört zu
`--github`/`--godev`; ohne diesen
Einstieg wäre der Vergleich nur mit Netz zu prüfen und damit gar nicht.

## Bindung

kein Gate, bewusst **nicht** in `gates`; Netz, gerufen vom Nachtlauf
(außer `freshness-trivy`, Grenze 4).
[`SPEC-098`](../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)
