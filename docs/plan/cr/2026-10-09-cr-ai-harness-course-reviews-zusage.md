# Change Request an `ai-harness-course` — Die `.d-check.yml`-Vorlage prüft ihre eigene Review-Zeile nicht

**Absender:** d-check (Adopter und Werkzeug) · **Datum:** 2026-10-09
**Richtung:** ausgehend ([`MR-035`](../../../harness/conventions.md#mr-035))
**Ziel:** `lab/templates/.d-check.yml` (Block `reviews`)
**Baseline-Stand:** `v6.17.0`
**Stand:** **beantwortet** — angenommen mit einer Ergänzung (Slug-Kennungen
brauchen `match: name`), siehe [Antwort](2026-10-09-antwort-ai-harness-course-reviews-zusage.md);
Umsetzung nach d-check
`v0.85.0`. Ein
[Befund](2026-10-09-befund-ai-harness-course-reviews-gleichgewicht.md) zum
Ruhezustand archivierender Repos folgte; seine Punkte 1 und 2 sind geplant.

---

## Bitte — `reviews` in der Vorlage so einschalten, dass es die Vorlagen-DoD erkennt und einen Leerlauf meldet

Die Slice-Vorlage formuliert den Review-Punkt der DoD so
(`templates/docs/plan/planning/slice.template.md`, Zeile 92):

> `- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor`

Die `.d-check.yml`-Vorlage schaltet das Modul `reviews` mit `done-dir` und
`reviews-dir` ein, ohne weiteren Schlüssel. Bis d-check `v0.84.0` erkennt das
Modul als Review-Zusage nur die Phrase „unabhängiger Review". Die Zeile der
Vorlage trägt sie nicht.

**Folge, gemessen:** Ein Repo, das beiden Vorlagen folgt, hat ein
`reviews`-Gate, das grün über einer leeren Menge läuft — kein Slice wird
geprüft, ein fehlender Report fällt nicht auf. Gemeldet von einem Adopter
(alle zwölf getesteten DoD-Formen unerkannt) und nachgemessen an d-check
selbst: mit einem Report-Verzeichnis ohne Reports meldet das Gate null Befunde.

**Was d-check `v0.85.0` dafür liefert:**

- Der Default erkennt zusätzlich die Vorlagen-Form „Review durchgeführt".
  Damit trägt die heutige Vorlage ohne Änderung.
- `reviews.require-promises: true` meldet, wenn unter vorhandenen Kandidaten
  keine einzige Zusage erkannt wird.
- `reviews.promise-pattern`, `reviews.match`, `reviews.recursive` und
  `reviews.skip-pattern` für Repos mit eigener Form, benannten Kennungen oder
  Volltexten in Unterverzeichnissen von `done/`.

**Vorschlag:** Die Vorlage setzt `require-promises: true` im `reviews`-Block:

```yaml
# reviews:
#   done-dir: docs/plan/planning/done
#   reviews-dir: docs/reviews
#   require-promises: true   # null erkannte Zusagen unter vorhandenen Slices ⇒ Befund
```

Begründung: Vorlage und Werkzeug-Default sind zwei Orte für eine Aussage. Ändert
sich eines von beiden — die Formulierung der DoD-Zeile oder der Default —,
läuft das Gate ohne den Schalter wieder still leer; mit ihm meldet es den
Leerlauf beim ersten Lauf.

**Optional, falls der Kurs Volltexte unter `done/wellenlos/` oder
`done/<welle-id>/` vorsieht:** dazu `recursive: true` und
`skip-pattern: '(?m)^> \*\*ARCHIVIERT'`, damit die archivierten Stubs nicht als
Slices ohne DoD zählen.
