`slice-226`s `packAliasFS` (`internal/adapter/driven/git/packalias.go`)
kompensiert, dass go-gits `DotGit`-Schicht (`v5.19.2`) Packs ausschließlich
über `ReadDir("objects/pack")` und den rekonstruierten kanonischen Namen
`pack-<hash>.{pack,idx}` entdeckt — ein interner, nicht dokumentierter
Mechanismus (`storage/filesystem/dotgit.ObjectPacks`/`objectPackPath`,
kein SemVer-Vertrag). Der Slice-Plan (§6) benennt das Risiko selbst: ein
künftiger go-git-Bump könnte diesen Mechanismus intern ändern, ohne dass ein
Compile-Fehler warnt — der Dekorator liefe dann still ins Leere. Zum
Zeitpunkt der Closure kein Freshness-Trigger, der einen erneuten Lauf der
beiden Regressionstests (`TestAllPathsPackUnterFremdemPraefix`,
`TestAllPathsPackMitUnbrauchbaremPraefixBleibtFehlerhaft`) gerade bei einem
go-git-Bump erzwingt.
