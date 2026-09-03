# Ein pfad-selektiver Commit nimmt still mit, was früher gestagt wurde

**Sub-Area:** `*`

`git rm`/`git mv` stagen sofort; ein späterer Commit, der nur weitere Pfade
addet, trägt den gesamten Index — die Botschaft beschreibt einen Teil des
Diffs, der Rest reist undeklariert mit, und der Zwischenstand kann gate-rot
sein. Kein Gate sieht das: `commits` prüft nur die Message-Form.
