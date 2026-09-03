# Review-Report: slice-173 — 2026-09-03

**Review-Art:** Code — geprüft gegen Plan, ADR und Konventionen (Modul 10
§Drei Review-Arten).

**Gegenstand:** Commit `85b1fce` (feat: Modul `reviews`)

**Skill:** `.harness/skills/reviewer.md`
**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-03

**Eingangs-Kontext:**

- [slice-173](../plan/planning/done/slice-173-review-report-deckung.md)
- [ADR-0081](../plan/adr/0081-reviews-modul.md)
- [`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)
- `AGENTS.md` (Hard Rules)

---

## Findings

### F-1 — Mehrzeilige DoD-Items werden nicht erkannt

- `kategorie`: HIGH
- `quelle`: DC-FA-RVW-001, Maintainability
- `pfad`: `internal/hexagon/core/rules/reviews.go:31` (Ausgangsfassung)
- `befund`: `reviewLineRE` verlangte Checkbox und Phrase „unabhängiger
  Review" auf derselben Zeile. Der überwiegende Bestand schreibt lange
  DoD-Punkte als Fließtext über mehrere Zeilen; die Phrase steht dabei
  häufig auf einer Folgezeile, nicht auf der Checkbox-Zeile selbst
  (gemessen: sechs Fälle im Bestand, u. a. `slice-138`).
- `verifizierbar`: ja — `make review-coverage` gegen den echten Bestand.
- `klasse`: „Zeilen-basierte Erkennung übersieht mehrzeilige
  Markdown-Konstrukte"

### F-2 — Widersprüchliche Doppel-Meldung bei unlesbarem `reviews-dir`

- `kategorie`: HIGH
- `quelle`: DC-QA-02, Maintainability
- `pfad`: `internal/hexagon/core/rules/reviews.go:82` (Ausgangsfassung)
- `befund`: Ein unlesbares `reviews.reviews-dir` mit vorhandenen
  Review-Zusagen erzeugte zusätzlich zu den Pro-Kandidat-Befunden eine
  generische „leere Prüfmenge"-Meldung — textlich widersprüchlich, weil
  die Menge in diesem Fall gerade nicht leer ist.
- `verifizierbar`: ja — konstruierter Test mit unlesbarem Verzeichnis und
  vorhandenen Kandidaten.
- `klasse`: „Fail-Closed-Meldung feuert redundant neben spezifischeren
  Befunden"

### F-3 — `.d-check.yml`-Kommentar unterschlägt die tragende Einschränkung

- `kategorie`: MEDIUM
- `quelle`: Maintainability
- `pfad`: `.d-check.yml:737` (Ausgangsfassung)
- `befund`: Der Block-Kommentar nannte nur „Review" statt „unabhängiger
  Review" — die ADR selbst nennt die engere Phrase als tragend
  (bloßes „Review" ist zu breit, gemessen an `slice-183`s
  „Adaptions-Review").
- `verifizierbar`: nein — Dokumentations-Genauigkeit, kein Gate-Lauf.
- `klasse`: „Konfigurationskommentar dehnt eine engere Vertrags-Phrase"

## Negativbefunde

- geprüft, ohne Befund: ADR-0081 — drei echte Alternativen verglichen,
  keine verkappte Einzeloption.
- geprüft, ohne Befund: Bestands-Ausnahme (fünf Dateien) — alle Pfade
  existieren, jede trägt tatsächlich eine „unabhängiger Review"-DoD-Zeile
  ohne passenden Report; live gegen eine exemptionsfreie Kopie der
  Konfiguration verifiziert (exakt fünf Funde, keiner mehr, keiner
  weniger).
- geprüft, ohne Befund: §3.4-Reinheit von `lastenheft.md`/`spezifikation.md`
  — keine bloßen `slice-`/`welle-`/Commit-Hash-Token in den neuen
  Abschnitten.
- geprüft, ohne Befund: `applyReviews`-Config-Validierung — Weißraum-Pfad,
  fehlendes `reviews-dir` bei gesetztem `done-dir`, ungültiges Glob — alle
  drei korrekt vor dem Lauf abgewiesen.
- geprüft, ohne Befund: `hasMatchingReview`/`sliceIDRE` — korrekte
  vollständige Zahlenerfassung (kein Abschneiden bei `slice-1` vs.
  `slice-100`), Substring-Match für `-r1`/`-r2`-Suffixe funktioniert gegen
  den echten `docs/reviews/`-Bestand.
- geprüft, ohne Befund: `AGENTS.md`, `harness/README.md`, beide READMEs —
  Fail-Closed-Bedingungen, nicht-rekursiver Scan und „noch nicht in
  gates"-Status konsistent über alle vier Dateien beschrieben.
- geprüft, ohne Befund: `TestReviewsBulletForms` — Groß-/Kleinschreibung
  von „Review" ist im gesamten `done/`-Bestand einheitlich (Substantiv,
  großgeschrieben); die feste Groß-/Kleinschreibungs-Form der Phrase ist
  damit keine praktische Lücke.
- geprüft, ohne Befund: archivierte Stubs (`done/<welle-id>/…`) — sowohl
  der eigene Test als auch der nicht-rekursive `List()`-Aufruf schließen
  sie korrekt aus.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Zeilen-basierte Erkennung übersieht
mehrzeilige Markdown-Konstrukte · Fail-Closed-Meldung feuert redundant
neben spezifischeren Befunden · Konfigurationskommentar dehnt eine engere
Vertrags-Phrase

## Verdikt

**Merge-blockierend:** ja (vor Behebung) — beide HIGH-Findings sind vor
der Closure behoben (Commit `da2149a`), mit vier neuen Regressionstests.
Neumessung gegen den echten Bestand nach dem Fix: weiterhin exakt fünf
Funde, keiner der sechs neu erfassten Mehrzeilen-Fälle war tatsächlich
ungedeckt.

**Übergabe:** Findings gingen an den Implementer; die Finding-Klassen
gehen in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und
ersetzt keine Verifikation — DoD-/Spec-Konformität prüfte der Verifier
separat, in eigenem Kontext.
