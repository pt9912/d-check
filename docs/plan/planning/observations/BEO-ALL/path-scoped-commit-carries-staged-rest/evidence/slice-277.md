**Vorgang:** slice-277
**Fund:** `git rm -r .harness/baseline/v6.17.0` stagte 55 Löschungen sofort; der
folgende pfad-selektive `git add` mit `git commit` für die Link-Tiefen von
[MR-074](../../../../../../../harness/conventions.md#mr-074) nahm sie mit (56 Dateien, Botschaft: „Link-Tiefen"), dazu war
`--no-verify` gesetzt, das den Hook-Lauf ohnehin übersprang. `git diff --cached --stat`
vor dem Commit (Schritt 21) hätte es gezeigt; der Commit war lokal und wurde mit
`git reset --soft` und `git restore --staged` aufgetrennt, bevor er den Remote erreichte.
