# Review R3 — slice-263: Nachzüge nach R2

- **Review-Art:** Code. Geprüft werden die Nachzüge nach R2 gegen die R2-Befunde, den Slice-Plan
  (`slice-263` §3 samt Plan-Änderungen nach R2 und nach der Verifikation) und die Hard Rules
  `AGENTS.md` §3.4/§3.7 sowie §5 Regel 3/13/15. Die DoD-Abhakung gehört nicht zu diesem Review.
- **Gegenstand:** `slice-263` · Commit `e87dc084` (Wortlaut-Fixes Lastenheft 0.99.1, Spezifikation,
  Tests) und `77c5a2c3` (Test der Exit-2-Meldung); Kontext `835020c5` (Plan-Änderung nach der
  Verifikation).
- **Skill:** `reviewer.md` @ 1.19.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** R2-Report dieses Slice (F-1 bis F-5); `spec/lastenheft.md` 0.99.1 —
  `DC-FA-PLAN-001` (Kandidatenmenge, fail-closed, Kriterien), `DC-FA-STRUCT-001` (fail-closed-Liste,
  Kriterien), Historie; `spec/spezifikation.md` §`DC-FA-PLAN-001.a` C2, Grund-Code-Zeile `SPEC-039`,
  Historie. Im Code gelesen: `internal/hexagon/core/rules/planning.go` (`CheckPlanningClosure`,
  `closureCandidates`), `internal/adapter/driven/configyaml/configyaml.go`
  (`structureBedingungsFehler`, `planning.closure.skip-pattern`-Prüfung), die Tests
  `planning_closure_rekursiv_test.go`, `structure_skip_test.go`, `configyaml_test.go`.
- **Probe** (Wegwerf-Kopie per `git archive HEAD` im Scratchpad, `make test IMAGE=dcheck-r3-probe`;
  Arbeitsbaum unverändert, Kopie und Image danach entfernt): Mutation `{"skip-muster", r.SkipPattern}`
  in `structureBedingungsFehler` ⇒ **rot** — `TestDecode_SkipPatternMeldungNenntSchluessel` mit
  `Fehler mit "skip-pattern" erwartet, got … structure[0]: skip-muster "^([" ist kein gültiges Regex`.
  Alle übrigen Pakete grün.

## Stand der R2-Befunde

| R2 | Stand | Beleg |
|---|---|---|
| F-1 | geschlossen | `DC-FA-STRUCT-001` fail-closed-Liste endet auf „`/skip-pattern`"; Kriterium „fail-closed (Inhalts-Ausnahme)"; Code `structureBedingungsFehler` führt `skip-pattern`, Probe rot |
| F-2 | geschlossen | Lastenheft „Gezählt wird die Gesamtmenge; ein einzelnes Unterverzeichnis ohne Kandidaten ist kein Befund", C2 „über die **gesamte** Kandidatenmenge", `SPEC-039` trennt die drei Prädikate; am Code: `closureCandidates` sammelt rekursiv in eine Liste, `len(names) == 0` steht nach `closureSkip` einmal für die Gesamtmenge — stimmt |
| F-3 | geschlossen an den benannten Stellen; ein Rest — siehe F-1 dieses Reports | „VORZUSTAND", „wie bisher", „der Vorzustand" entfernt, Testname `…OhneSchalterIstKandidat` |
| F-4 | geschlossen | (a) „—;" weg; (b) „Unterverzeichnissen von `planning.closure.dir`" explizit; (c) „in jeder geprüften Datei"; (d) Apposition hängt an `planning.closure.dir`; C2 „wie zuvor" ersetzt |
| F-5 | geschlossen | zwei neue Kriterien, je durch Test gedeckt (`TestClosureRecursive_UnlesbaresUnterverzeichnisFailClosed`; `TestDecode_SkipPatternMeldungNenntSchluessel` für beide Module); Historie beider Straten je eine neue Zeile |

## Findings

### F-1 — LOW — „neuen" im Test-Kommentar

- **quelle:** `AGENTS.md` §3.7 („Ein Kommentar beschreibt, was da ist"); R2 F-3
- **pfad:** `internal/hexagon/core/rules/planning_closure_rekursiv_test.go` · „Ohne die neuen Schlüssel
  bleibt die Meldung eines unlesbaren closure.dir"
- **befund:** Dieselbe Klasse wie R2 F-3 steht an einer Stelle weiter: „neu" hat für einen späteren
  Leser keinen Bezugspunkt (gemeint sind `recursive` und `skip-pattern`). Die Botschaft von `e87dc084`
  sagt „Chronik-Woerter aus den neuen Tests entfernt" — für diese Stelle stimmt das nicht ganz. Der
  Kommentar trägt eine Zusage (byte-identische Meldung), daher LOW.
- **verifizierbar:** nein.
- **klasse:** `kommentar-mit-chronik-wort`

### F-2 — INFO — „mit dem Pfad des Unterverzeichnisses" lässt offen, in welchem Feld

- **quelle:** Maintainability
- **pfad:** `spec/lastenheft.md` (`DC-FA-PLAN-001`, Kriterium „fail-closed (Unterverzeichnis und
  Muster)") · „then `closure-note-missing` mit dem Pfad des Unterverzeichnisses"
- **befund:** Der Code setzt `file` = `planning.closure.dir` und nennt das Unterverzeichnis in der
  Meldung (`closureFinding(dir, 1, dir, …)`, `"Closure-Verzeichnis "+err.Error()+…`); die Spezifikation
  C2 sagt das genau so („die Meldung nennt seinen Pfad"), der Test prüft `Message`. Das Kriterium auf
  Rang 1 lässt offen, ob `file` oder Meldung gemeint ist — eine Schärfung durch Rang 2, kein
  Widerspruch.
- **verifizierbar:** nein.
- **klasse:** `vertragssatz-schwer-lesbar`

### F-3 — INFO — „dann" in C2 hat seinen Bezug verloren

- **quelle:** Maintainability
- **pfad:** `spec/spezifikation.md` §`DC-FA-PLAN-001.a` C2 · „ein Unterverzeichnis ohne Kandidaten ist
  kein Befund; die Meldung nennt"
- **befund:** Nach dem eingeschobenen Satz über das Unterverzeichnis schließt „die Meldung nennt dann
  das Muster" grammatisch an „kein Befund" an; gemeint ist die Nullmengen-Meldung bei gesetztem
  `skip-pattern`. Keine falsche Aussage, der Code (`closureSkipNote`) tut das Gemeinte.
- **verifizierbar:** nein.
- **klasse:** `vertragssatz-schwer-lesbar`

## Negativbefunde

- **Nullmengen-Aussage gegen Code:** „gesamte Kandidatenmenge" in Lastenheft, C2 und `SPEC-039` gegen
  `CheckPlanningClosure` — eine Liste über alle Ebenen, ein Nullmengen-Test nach dem Abzug, ein leeres
  Unterverzeichnis trägt nichts bei und meldet nichts. Stimmt. Ohne Befund.
- **`SPEC-039` neu gefasst:** drei Prädikate getrennt durch „oder"; „ein Kandidat (Schritt C2)" verweist
  auf die Definition statt sie zu kopieren; Ausschluss von `-thin`/`-boilerplate` erhalten. Ohne Befund.
- **Exit-2-Aufzählung:** `DC-FA-STRUCT-001` nennt `skip-pattern`, `DC-FA-PLAN-001` unverändert
  vollständig; Code und Test decken beide (Probe rot). Ohne Befund.
- **Neue Akzeptanzkriterien gegen Tests:** „fail-closed (Unterverzeichnis und Muster)" — unlesbares
  Unterverzeichnis per `listErrSubFS`, Pfad in der Meldung; Exit 2 mit Schlüssel per `configyaml`-Test.
  „fail-closed (Inhalts-Ausnahme)" — derselbe Test, `structure`-Fall. Ohne Befund außer F-2.
- **Lastenheft-Satzbau:** Kandidaten-Absatz in drei Sätzen, Bezüge eindeutig; Absatz „Geprüft wird
  ausschließlich" mit Apposition am richtigen Bezugswort. Ohne Befund.
- **Historie (`MR-032`):** Lastenheft 0.99.0 → 0.99.1 bei Status Draft, neue Zeile oben, die
  0.99.0-Zeile byte-gleich; Spezifikation eine neue Zeile 2026-10-09 über der bestehenden, diese
  unverändert. Beide Zeilen beschreiben, was sich ändert, ohne Chronik der Entstehung über „Nachzug
  nach Review" hinaus — dieselbe Form wie die Bestandszeile 0.97.2. Ohne Befund.
- **Referenz-Richtung (§3.4):** in den hinzugefügten Spec-Zeilen kein `slice-`, `welle-`, `ADR-`-Token,
  kein Commit-Hash. Ohne Befund.
- **Kommentare (§3.7) in `77c5a2c3`:** Testkommentar trägt eine Zusage, kein Befund-Marker, keine
  Slice-Nummer; die Kennung „V-2" steht nur in der Botschaft. Ohne Befund.
- **Test-Determinismus:** Map-Iteration in `TestDecode_SkipPatternMeldungNenntSchluessel` ist
  ordnungsunabhängig (je Eintrag eigene Prüfung mit `t.Errorf`). Ohne Befund.
- **Code:** keine Produkt-Code-Änderung in beiden Commits; kein Import, kein `//nolint`, keine
  Schwelle. Ohne Befund.
- **Abgrenzung:** der Test-Commit steht nach der Plan-Änderung `835020c5`. Ohne Befund.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 0 | — |
| LOW | 1 | F-1 |
| INFO | 2 | F-2, F-3 |

Wiederkehrende Finding-Klasse für die Closure-Notiz: `kommentar-mit-chronik-wort` (R1 F-9, R2 F-3,
hier F-1 — jede Runde schließt die benannte Stelle und lässt eine Nachbarstelle stehen).

## Verdikt

**Freigegeben.** R2 F-1 bis F-5 sind geschlossen; die Nachzüge bringen keinen neuen HIGH- oder
MEDIUM-Befund. F-1 (LOW) ist eine Einwort-Korrektur, F-2/F-3 sind Lesehilfen ohne falsche Aussage —
keiner blockiert.
