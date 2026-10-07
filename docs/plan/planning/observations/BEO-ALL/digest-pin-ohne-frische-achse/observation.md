# Ein digest-gepinntes Fremd-Image ohne Frische-Achse

**Sub-Area:** `*`

Ein Fremd-Image wird nach [ADR-0011](../../../../adr/0011-digest-pins-build-gate-images.md) per Digest gepinnt, aber an einer Stelle,
die weder Dependabot liest noch ein Nachtlauf-Target abfragt — etwa als
`with:`-Eingabe einer GitHub Action. Der Pin ist damit reproduzierbar und
zugleich unbewacht: ein neuerer Bau desselben Tags oder eine neuere Version
fällt niemandem auf. Die Form gleicht den gewächterten Pins, die Deckung
nicht.
