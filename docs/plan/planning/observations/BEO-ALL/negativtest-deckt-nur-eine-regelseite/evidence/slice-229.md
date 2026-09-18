Gefunden vom unabhängigen Review von slice-229
([Report](../../../../../../reviews/2026-09-18-slice-229-matrix-instanz-identitaet-review-r1.md),
F-2): `TestDecode_MatrixAllowIfSameIDFailClosed` brach in allen drei
ursprünglichen Unterfällen ausschließlich das `token` der `from`-Klasse
(`slice`); die `to`-Klasse (`review`) trug in jedem Unterfall ein gültiges
Token. `validateMatrixAllowIfSameID` prüft beide Seiten in einer Schleife
und kehrt beim ersten Fehler zurück — der `to`-seitige Codepfad war damit
nie isoliert getroffen, obwohl das Lastenheft die `to`-Klasse wörtlich als
eigene Akzeptanzkriterium-Bedingung nennt. Ein vierter Unterfall (nur die
`to`-Klasse gebrochen) hat die Lücke geschlossen.
