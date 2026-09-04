# Ein archivierter Slice-Zip kann Inhalt tragen, der erst nach seinem letzten Commit-Stand entstand

**Sub-Area:** `*`

`tools/archive-wave`s Einzel-Slice-Modus (`ApplySlice()`) baut das Archiv-Zip
aus dem aktuellen Datei-Inhalt auf der Platte, nicht aus dem Git-Objekt des
letzten Commits. Wird eine Reihe unabhängiger Slices **sequenziell** in
derselben Sitzung archiviert (mehrere `-slice=<id> -apply`-Läufe ohne
zwischenzeitlichen Commit), kann ein früherer Lauf per `RewriteRepo()` einen
Verweis in einem noch nicht archivierten, später folgenden Slice umschreiben
(Pfad-Nachzug auf den neu archivierten Ort). Wird dieser spätere Slice dann
selbst archiviert, zippt das Werkzeug den bereits veränderten Stand — das
Archiv enthält damit einen Text, der zu keinem Zeitpunkt so committet war.

Der Effekt ist inhaltlich harmlos (ein Pfad wird korrekt auf den neuen
Speicherort nachgezogen), aber er macht den Zip-Inhalt **von der
Verarbeitungsreihenfolge des Batches abhängig** — ein Determinismus-Anliegen
([`DC-QA-02`](../../../../../../spec/lastenheft.md#dc-qa-02--determinismus)),
das der bestehenden Testsuite fremd ist. Gefunden vom unabhängigen Review von
slice-200 (2026-09-04): die Zips von slice-184 und slice-188 enthalten je
einen auf den neuen Speicherort von slice-183 nachgezogenen Link, wo der
Ursprungs-Commit noch auf dessen alten, flachen Pfad zeigte.
