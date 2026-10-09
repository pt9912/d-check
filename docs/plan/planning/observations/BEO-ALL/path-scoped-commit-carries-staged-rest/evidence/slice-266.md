**Vorgang:** slice-266
**Fund:** Ein `git add` mit einer schon entfernten Datei brach ab; die
folgende Kette schrieb die neue Commit-Botschaft nicht, und der Commit
erfasste die gestagte Löschung von `manual_test.go` unter der Botschaft der
vorigen Plan-Änderung (`d04d5c8a`, gepusht). Schritt 21 — `git diff --cached
--stat` vor dem Commit lesen — war nicht ausgeführt. Der folgende Commit nennt
es (Review R2 F-5).
