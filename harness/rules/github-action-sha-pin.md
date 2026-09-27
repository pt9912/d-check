# GitHub-Action-Referenzen sind SHA-gepinnt — die lokale Ausnahme und ihre Grenzen

Ausgelagerter Körper zu `AGENTS.md` §3.9. Der operative Kern (voller
Commit-SHA + Tag-Kommentar, Begründung, Durchsetzung) steht dort; diese
Datei trägt die lokale Workflow-Ausnahme und ihre Prüf-Grenzen —
Auslagerungs-Muster aus
[ADR-0096](../../docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md).

## Die lokale Ausnahme ist keine Lockerung

([ADR-0068](../../docs/plan/adr/0068-lokale-workflow-referenzen-ohne-pin.md)):
eine **lokale** Workflow-Referenz (`uses: ./.github/workflows/x.yml`) kann
keinen SHA tragen und **braucht keinen**. Sie löst auf denselben Commit auf
wie der aufrufende Workflow und ist damit **stärker** gebunden als ein
SHA-Pin — sie kann per Konstruktion nicht driften. **An die Stelle des
Pin-Checks tritt die Frage, die hier trägt: existiert das Ziel?** Ein
vertippter Verweis fiele sonst erst zur Laufzeit auf; `make workflow-pins`
meldet ihn als `uses-local-missing`. Sie gilt **nur** dem `./`-Präfix; eine
Referenz in ein fremdes Repository
(`owner/repo/.github/workflows/x.yml@ref`) fällt unter die Regel wie jede
Action. Die Zahl der lokalen Referenzen steht in der Erfolgsmeldung, statt
stillschweigend übergangen zu werden.

## Die Existenz ist nicht die einzige Frage

([ADR-0071](../../docs/plan/adr/0071-lokale-workflow-referenz-rechte-pruefung.md)).
Ein aufgerufener Workflow bekommt nur die Rechte, die der aufrufende **Job**
selbst führt; verlangt er mehr, lehnt GitHub den **ganzen Lauf vor dem
ersten Job** ab (`startup_failure`, kein Log) — die Existenz-Prüfung allein
sieht das nicht. Geprüft wird deshalb auch die **Rechte-Anforderung des
Ziels**: ein Job ohne eigenes `permissions:`, dessen Ziel Rechte verlangt
(`uses-local-perms-undeclared`), und ein Aufrufer, der einen geforderten
Scope zu niedrig führt (`uses-local-perms-narrow`). Was der Wächter nicht
sicher liest, meldet er (`uses-local-perms-unreadable`), statt es zu
übergehen. **Er deckt damit eine Fehlerklasse, nicht die Lauffähigkeit** —
und die Zerlegung ist eine Näherung über die YAML-Block-Form, keine
Parser-Zusage.

## Grenze

`make workflow-pins` trägt die **Form** — voller SHA plus Tag-Kommentar —,
nicht die **Gültigkeit**: ob der SHA existiert und den Commit bezeichnet,
den der Tag-Kommentar behauptet, prüft er nicht. *(Auflösungs-Trigger:
permanent — die Gültigkeitsfrage ist Netz und gehört zur
Freshness-Familie.)*
