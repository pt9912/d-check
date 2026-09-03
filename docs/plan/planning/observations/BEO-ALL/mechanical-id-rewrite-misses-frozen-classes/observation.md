# Eine repo-weite mechanische ID-Ersetzung kennt nur die benannten Frozen-Verzeichnisse, nicht jede Frozen-Klasse

**Sub-Area:** `*`

Beim Umhängen eines Registerpfads über den ganzen Baum (`sed` auf einen
Kennungs-Token) wurden die drei im Kanon benannten Frozen-Verzeichnisse
(`done/`, `docs/reviews/`, `harness/conventions/done/`) korrekt ausgenommen —
aber zwei weitere Klassen mit derselben „zitiert den Stand ihrer Zeit"-
Eigenschaft nicht: Accepted-ADR-Kerne (§3.5-Immutabilität) und gesendete CRs
(`docs/plan/cr/`). Der Fehler wurde durch sofortigen `git diff`-Check nach
dem Lauf gefangen, nicht durch Vorab-Analyse.
