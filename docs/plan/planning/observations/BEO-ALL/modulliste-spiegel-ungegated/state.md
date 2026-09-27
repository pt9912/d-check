**Stand:** verkörpert — sechs Deckungstests gegen `model.ValidModules()`
(`internal/hexagon/core/app/registry_mirror_test.go`) plus eine dritte
Prüfrichtung im Netzlos-Guard und ein `FOCUS_DISABLE`-Deckungstest
(`internal/adapter/driven/configyaml/gate_consistency_test.go`) — liegt in
`internal/hexagon/core/app/registry_mirror_test.go` (seit slice-238).
Bewusst nicht mechanisiert: die Bereichskürzel-Liste (Abkürzungen) und die
„fixe Standard-Modulset"-Stellen (bewusste Teilmenge) — beide von zwei
unabhängigen Reviews bestätigt.
