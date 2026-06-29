# Review — slice-054 (`codepaths.ignore-refs`) · R2 (Code- / Config-Korrektheit)

## Kopf-Metadaten

- **Reviewer:** R2 (unabhängig, Schwerpunkt Code- und Config-Korrektheit) — **nicht**
  Doc-first-Prozess (separater Reviewer R1).
- **Datum:** 2026-06-29
- **Reviewer-Skill:** `.harness/skills/reviewer.md` v1.2.0.
- **Gegenstand:** Working-Tree-Änderungen slice-054 (vor Commit), `git diff HEAD` + untracked.
- **Quellen:** `spec/spezifikation.md` §DC-FA-CODE-001.a (Schritt 5 + §2-Schema),
  `spec/lastenheft.md` DC-FA-CODE-001, ADR-0025, Slice-Plan slice-054.
- **Rollen-Abgrenzung:** kein Verifier — Gate-Status (`make ci` Cov 93,40 %, completeness,
  semgrep 0, arch-check grün) als gegeben übernommen.

## Findings

### F-1 — MEDIUM — `repo-escape`/`anchor-missing`-Unterdrückung durch `ignore-refs` untestiert

- **Quelle:** DC-FA-CODE-001 / ADR-0025 (öffentlicher Vertrag) · Skill-Anker „fehlende
  Negativtests bei neuem öffentlichen Vertrag"
- **Pfad:** `internal/hexagon/core/rules/codepaths_test.go` (`TestCodepathsIgnoreRefs`)
- **Befund:** Der Vertrag (Schritt 5: „nicht existenz- **oder anker**-geprüft … Escape →
  `repo-escape`" *nach* der Auslassung; ADR-0025: „`codepath-missing` (**und
  `anchor-missing`/`repo-escape`**) entfallen") nennt drei unterdrückte Grund-Codes. Der
  Test deckte nur die `codepath-missing`-Unterdrückung ab. **Failure-Szenario:** Ein
  künftiger Refactor verschiebt `if ignored(rel, ignoreRefs) { return nil }` hinter den
  `escaped`- oder Anker-Block. Ein ignorierter, eskapierender Pfad meldete dann wieder
  `repo-escape`, ein ignorierter `.md#frag`-Pfad wieder `anchor-missing` — beides bricht den
  Vertrag, doch die Suite bliebe grün. Die `repo-escape`-Unterdrückung ist der
  sicherheits-nächste der drei Codes.
- **Verifizierbar:** ja — ein Testfall mit (a) ignoriertem `../…`-eskapierenden Ziel und (b)
  ignoriertem `…md#fehlt`-Ziel würde unter einer Umsortierung fehlschlagen; `make test`.

### F-2 — INFO — `Target` trägt den rohen Span, `ignore-refs` matcht den aufgelösten `rel`

- **Quelle:** DC-FA-CODE-001 (Match-Semantik) · Maintainability/Auffindbarkeit
- **Pfad:** `internal/hexagon/core/rules/codepaths.go` (Match auf `rel` vs. `Target: value`)
- **Befund:** `ignored(rel, …)` matcht gegen den aufgelösten **Wurzel-relativen** Pfad
  (korrekt, spec-konform), während der Befund `Target` den **rohen** Span-Wert trägt. Für
  `./`/`../`-Formen weichen beide ab. **Failure-Szenario:** Ein Adopter sieht `Target:
  ../weg.sh`, trägt `../weg.sh` ein, und der Befund bleibt, weil intern gegen `weg.sh`
  gematcht wird. Reibung, kein Korrektheitsfehler; in der Spezifikation gedeckt, im
  Benutzerhandbuch nicht explizit als Resolved-vs-Roh herausgestellt.
- **Verifizierbar:** nein (Beobachtung am Code/Doc).

## Negativbefunde (geprüft, ohne Befund)

- **Skip-Platzierung:** `ignored(rel, ignoreRefs)` steht **nach** der Auflösung in BEIDEN
  Zweigen (`ResolveConfigPath` für rootRel, `ResolveTarget` datei-relativ) und **vor**
  Escape/Existenz/Anker. Die `repo-escape`-Unterdrückung ist spec/ADR-konform gewollt. Korrekt.
- **Stilles-Ventil-Risiko:** Die `.d-check.yml`-Glob ist der exakte Pfad
  `tools/adr-immutable-check.sh` — silenct genau die klassifizierten Referenzen plus die
  referenz-weit gleichlautenden in `done/`/Slices/ADR-0025. Bare `adr-immutable-check.sh` und
  `bash tools/…` (Space) werden ohnehin nicht klassifiziert. Kein zu breiter Glob; das
  Über-Breit-Risiko ist in ADR-0025/Slice §4 dokumentiert und akzeptiert.
- **Match-Semantik:** Match gegen aufgelöstes `rel`, nicht gegen den rohen Span;
  `matchGlob`→`path.Match` identisch zu `exempt-paths`/`scan.ignore` (geteilter Helfer in
  `paths.go`, unverändert).
- **Default-Identität / DC-QA-02:** `ignored(rel, nil)` → false; bei leerem/nil `IgnoreRefs`
  kein Verhaltenswechsel. Bestätigt durch `cfgPlain` + alle Bestands-codepaths-Tests.
- **Test-Abdeckung Kern:** Happy, Negative, Glob, **referenz-weit** (derselbe Pfad in zwei
  Dateien) und Default-leer abgedeckt; Lücke nur bei Anker-/Escape-Interaktion → F-1.
- **configyaml:** `rawCodepaths.IgnoreRefs` (`yaml:"ignore-refs"`) geparst und durchgereicht.
  Keine separate Glob-Validierung — konsistent mit `exempt-paths`.
- **Config-Surface:** `--print-config` und `suggest.renderHarness` schreiben die Zeile
  **auskommentiert** (ignore-refs startet leer — korrekt, Tombstone-Register ist
  repo-historie-spezifisch). Round-Trip-Decodes der Ausgaben in `make ci` grün.
- **arch-check / Hexagon:** Keine neuen Imports/Dependencies; `ignored`/`matchGlob` im selben
  Paket `rules`. Kein Import-Regel-Verstoß (ADR-0005), kein Stilles-Grün-Pfad. Die
  Gate-Berührungen (`Makefile`, `.githooks/pre-commit`, `tools/completeness-check.sh`) sind
  reine Kommentar-Änderungen; `VCS_DISABLE` unverändert.
- **Realer Fall / Wirksamkeit:** `tools/adr-immutable-check.sh` ist `git rm`'t; die zwei
  klassifizierten Frozen-Referenzen werden durch den exakten ignore-refs-Eintrag still — die
  Falle ist am realen Fall ohne Edit an immutabler Doku aufgelöst.

## Kategorie-Summary

| Kategorie | Anzahl |
| --- | --- |
| HIGH | 0 |
| MEDIUM | 1 (F-1) |
| LOW | 0 |
| INFO | 1 (F-2) |

## Verdikt

Code/Config der Kern-Mechanik ist korrekt und spec-treu — Skip-Platzierung in beiden
Auflösungs-Zweigen, Resolved-`rel`-Match-Semantik, Default-Byte-Identität, configyaml-Parität
und Hexagon-Sauberkeit sind verifiziert; der reale Tombstone-Fall greift wie spezifiziert. Ein
MEDIUM (F-1) blockiert typischerweise: nur die `codepath-missing`-Vertragshälfte war durch
einen Test verriegelt. F-2 ist INFO und blockiert nicht.

## Einarbeitung (Implementation, 2026-06-29)

- **F-1 — behoben:** Neuer Test `TestCodepathsIgnoreRefsUnterdruecktEscapeUndAnker` verriegelt
  **alle drei** Grund-Codes-Unterdrückungen: ein ignoriertes eskapierendes `../oben.md` (kein
  `repo-escape`) und ein ignoriertes `docs/real.md#fehlt` (kein `anchor-missing`) — ohne
  Eintrag feuern beide (2 Befunde), mit Eintrag keiner. Eine Umsortierung des
  `ignored()`-Aufrufs hinter Escape oder Anker bricht den Test. `make test` grün.
- **F-2 — behoben:** Das Benutzerhandbuch stellt jetzt explizit heraus, dass das Glob den
  **aufgelösten, Wurzel-relativen** Ziel-Pfad matcht (nicht den rohen `Target`-Span), mit
  Beispiel für die `./`/`../`-Auflösung.
