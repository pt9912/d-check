`slice-226`s `packHashSuffixRE` (`internal/adapter/driven/git/packalias.go`)
akzeptiert jede Datei, deren Name mit einem gültigen SHA1/SHA256-Hex-Suffix
endet, als aliasfähigen Pack — ohne zu prüfen, ob dieser Hash tatsächlich der
Content-Hash der Datei ist. Der Slice-Plan (§6) benennt das Risiko selbst.
Der unabhängige Review (R1-F-2) fand dazu einen konkreten, behobenen
Determinismus-Mangel bei zwei Dateien mit identischem Suffix
(`packAliases()` iterierte über eine randomisierte Map) und stellte fest,
dass go-gits eigener Idx-Checksummen-Vergleich einen Fehlgriff meist
fail-closed auffängt — das ist aber, wie dort ausdrücklich benannt, eine
Folge der nachgelagerten Schicht, keine Garantie der eigenen Heuristik.
