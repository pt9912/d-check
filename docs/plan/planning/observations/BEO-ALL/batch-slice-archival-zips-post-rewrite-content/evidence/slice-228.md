Zweites Auftreten, reproduziert beim eigenen Batch-Lauf von slice-228 (28
Slices `slice-200`–`slice-227` sequenziell ohne Zwischen-Commit): Die
Archivierung von `slice-207` zog per `RewriteRepo()` drei Verweise in
slice-208s damals noch flacher Datei nach, bevor slice-208 selbst — als
nächster Lauf derselben Sitzung — archiviert wurde. Dasselbe Muster trat mehrfach auf (`slice-207` rewrote zusätzlich
`slice-222`/`slice-224`, `slice-208` rewrote `slice-209`/`slice-210`/
`slice-222`/`slice-224`) — hier als **ein** Vorgang gezählt, wie beim
ersten Auftreten. Der Effekt blieb wie dort inhaltlich harmlos (korrekter
Pfad-Nachzug); keine Werkzeug-Änderung in diesem Slice (§1-Ausschluss von
slice-228).
