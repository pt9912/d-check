# Impl-Review R1 — slice-064 (gate-consistency.sh Voll-Tombstone)

**Datum:** 2026-07-06 · **Reviewer:** unabhängiger Subagent · **Verdikt:
NACHBESSERN → behoben** (1 MEDIUM eingearbeitet; Kern korrekt und vollständig).

**Gegenstand:** slice-064 / ADR-0032 — das letzte `tools/*.sh` (gate-consistency.sh)
wird entfernt, seine Rest-Prüfung (DC-QA-03-Modulliste der `.d-check.yml`) wird ein
getippter Go-Test; `make gate-consistency` fährt nur noch `--enable targets`. Kein
Lastenheft-CR, kein Release.

## Befund (behoben)

**MEDIUM — harness/README.md §Sensors: Bindungsspalte nicht migriert.** Die
Vertrags-Spalte war ehrlich umgeschrieben, die **Bindungs**-Spalte (Spalte 3)
aber unverändert: die `make gate-consistency`-Zeile band noch `DC-QA-03`, obwohl
der Vertrag es als „voll abgelöst → `make test`" ausweist (Row-interner
Widerspruch), und `make test` trug `DC-QA-03` nicht in der Bindung. ADR-0032/DoD-F-3
verlangen die Umschreibung der **Bindungsspalte** „nicht nur die Prosa". Kein Gate
wird rot (das Modul `targets` prüft nur `make X`↔Regel) — genau deshalb überlebte
die Harness-Lüge, die der Slice tilgt.

**Fix (eingearbeitet):** `gate-consistency`-Bindung → `DC-FA-TGT-001`/ADR-0031
(parallel zu `planning-check`, das sein Modul-Requirement bindet); `make test`-Bindung
→ `DC-QA-02`/`DC-QA-03`.

## Verifiziert sauber

- **Paritäts-Äquivalenz:** der Go-Test (`assertNetlessModules`) ist ein Superset
  des Skripts — verlangt alle **8** netzlosen Doku-Module präsent (statt 5) **und**
  `external`/`vcs` abwesend, über echten `configyaml.Decode`. Fällt
  `spans`/`hostpaths`/`versions` aus der Liste, feuert er (das Skript tat das
  nicht, R1-F-6). Keine Regression.
- **Fail-closed:** `os.ReadFile`-/`Decode`-Fehler → `t.Fatalf`; Relativpfad
  `../../../../.d-check.yml` korrekt (4-Segment-Tiefe wie `docexamples_test.go`).
- **Guard-Isolation (slice-057-R3):** synthetische Modul-Listen treffen
  ausschließlich die Assertion; „Modul fehlt"/„external gesetzt"/„vcs gesetzt"
  sterben je an einem Case; Präsenz- und Forbidden-Schleifen-Mutationen scheitern.
- **Tombstone-Vollständigkeit:** `git rm` gestaged; `codepaths.ignore-refs` +=
  `tools/gate-consistency.sh` fängt die Inline-Verweise; 0 verbliebene
  Markdown-Links aufs Skript; kein dangling Produkt-Code-Verweis (targets.go
  bereinigt, config_template.go nennt nur das make-Target).
- **Makefile-Rückbau:** `gate-consistency`-Regel bleibt + in beiden `doc-tables`
  dokumentiert → `targets`-Dogfood bleibt grün; `make gates` aggregiert weiter
  `test` + `gate-consistency` (keine Deckungslücke).

## INFO (nicht blockierend)

ADR-0032-Status bei Review noch `Proposed` (Closure-Flip erwartungsgemäß offen);
`.d-check.yml`-YAML-Kommentar nennt „ADR-0031" statt Volltombstone — kosmetisch,
nicht irreführend.
