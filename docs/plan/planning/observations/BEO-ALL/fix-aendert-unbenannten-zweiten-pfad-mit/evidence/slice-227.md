`slice-227`s Reihenfolge-Fix in `CheckCommits` (`internal/hexagon/core/rules/commits.go`)
sollte laut Root-Cause-Analyse und DoD nur `--enable commits --range
<unauflösbar>` betreffen. Der unabhängige Review (R1-F-2, MEDIUM) fand
empirisch, dass `--enable commits --staged` mit leerer `commits:`-Config
dasselbe Verhalten teilte (still Exit 0 vor dem Fix, korrekt fail-closed
Exit 2 danach) — weder im Slice-Plan noch in den ursprünglichen zwei
Regressionstests benannt. Noch im selben Slice mit einem dritten
Regressionstest geschlossen (`TestCheckCommitsStagedTrotzLeererConfigAufgeloest`).
