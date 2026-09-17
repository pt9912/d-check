# Ein Adapter verlässt sich auf das interne, nicht-öffentliche Verhalten einer Abhängigkeit — ohne Compile-Zeit-Warnung bei Drift

**Sub-Area:** `internal/adapter/driven/git`

Ein Adapter kompensiert eine Lücke in einer Fremdbibliothek, indem er ihr
**internes, undokumentiertes** Verhalten nachbildet (hier: wie go-gits
`DotGit`-Schicht Packs über `ReadDir` + einen rekonstruierten Dateinamen
entdeckt und öffnet) — nicht ihre öffentliche, SemVer-gebundene API. Ein
Bump der Abhängigkeit kann diesen Mechanismus ändern, ohne die öffentliche
API zu brechen: Der Adapter kompiliert weiter, aber der Kompensations-Code
läuft ins Leere, ohne dass irgendein Signal — Compile-Fehler, Lint-Befund —
darauf hinweist. Nur ein laufender Test gegen das tatsächliche Verhalten der
gepinnten Version deckt das auf, und nur dann, wenn der Test bei jedem
Dependency-Bump erneut läuft.

**Von verwandten Klassen abgegrenzt:** anders als
[`BEO-ALL/spec-randbedingung-ohne-test`](../spec-randbedingung-ohne-test/observation.md)
(eine zugesagte Randbedingung ohne eigenen Test) hat dieser Fall einen Test —
die Lücke liegt nicht im Testen, sondern darin, dass **kein Trigger** einen
erneuten Lauf gerade beim Anlass (Dependency-Bump) erzwingt.
