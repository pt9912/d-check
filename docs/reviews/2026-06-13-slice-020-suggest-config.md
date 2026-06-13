# Review-Report: slice-020 `--suggest-config` (`DC-FA-CLI-006`) — 2026-06-13

**Review-Art:** Code-Review (`/code-review`, high effort) — 4 unabhängige
Finder-Winkel (Korrektheit line-by-line, YAML-Output/Edge-Cases,
Cross-File/Vertrag, Cleanup/Altitude) + Verifikation.

**Gegenstand:** Commit-Range `origin/main...HEAD` (slice-020:
`internal/hexagon/core/suggest.go` neu, `internal/adapter/driving/cli/cli.go`
geändert; Vertrag/Spez/Doku).

**Skill:** `/code-review` · **Datum:** 2026-06-13

**Eingangs-Kontext:** [`spec/lastenheft.md`](../../spec/lastenheft.md)
`DC-FA-CLI-006`, [`spec/spezifikation.md`](../../spec/spezifikation.md)
§`DC-FA-CLI-006.a`, `DC-QA-02`/`DC-QA-03`,
[`.d-check.yml`](../../.d-check.yml).

---

## Findings und Disposition

| # | Schwere | Finding | Disposition (Review R1) |
|---|---|---|---|
| 1 | 🔴 hoch | Gerüst emittiert `ids:`-Block, aber `modules:` ohne `ids` → abgeleitete Muster im erzeugten Config **inaktiv** | **gefixt** — `renderSuggestion` nimmt `ids` in die Modul-Liste auf, wenn Muster vorliegen; Test prüft es |
| 2 | 🟠 mittel | `probeOptInModules` scannt Default-Wurzeln (nil Roots), Gerüst empfiehlt `roots: ["."]` → Scope-Mismatch, Module übersehen; `Config{Modules}` totes Feld | **gefixt** — Probe mit `Roots: ["."]`, totes Feld entfernt |
| 3 | 🟠 mittel | `target: %s` unquoted → Quellpfad mit `:`/` #` bricht Decode oder verfälscht still | **gefixt** — `target: %q` (gequotet) |
| 4 | 🟠 mittel | Heading-Token mit Satzzeichen/Markup (`ADR-0001:`, `[ADR-0001](…)`) verfehlt die Kennung | **gefixt** — `stripHeadingLinks` + Trim von `` ` `` `.,:;`; Test mit Doppelpunkt-Heading |
| 5 | 🟡 niedrig | `--suggest-config ,` → stilles leeres Gerüst, Exit 0 | **gefixt** — leere Quellenliste = Nutzungsfehler (Exit 2); Test |
| 6 | 🟡 niedrig | `hasLetter` global → `[A-Za-z]?` über alle Präfixe (`ADR-0042x`-Over-Match in Inhalt) | **akzeptiert** — advisory Ausgabe, Round-Trip hält, der Mensch verengt (dokumentierte „Scaffold, kein Orakel"-Eigenschaft) |
| 7 | 🟡 niedrig | `--suggest-config .`/`/` → ganzes Repo statt benannter Quelle | **akzeptiert** — Whole-Repo-Ableitung ist ein zulässiger Nutzungsfall, kein Schaden |
| 8 | 🟡 niedrig | Symlink-Quelle (`KindSymlink`) fällt als Datei durch | **akzeptiert** — Edge; `ReadFile`/Exit 2 ist sicher, kein Leak |
| 9 | 🟡 niedrig | `--print-config` + `--suggest-config` → print-config gewinnt still | **akzeptiert** — Kurzschluss-Reihenfolge dokumentiert; harmlos |
| 10 | 🟡 cleanup | `idShape` als zweite ID-Wahrheit; `splitSources` ≠ `multiFlag` | **akzeptiert** — `suggest` ist bewusst *breiter* als die strenge `ids`-Regex; bewusste Trennung |

## Negativbefunde (verifiziert sauber)

- **Read-only/netzlos** (`DC-QA-03`): kein Schreibpfad; Probe übergibt
  `nil`-HTTPChecker, nur `external` (nicht im opt-in-Set) berührt Netz.
- **Round-Trip-Invariante:** der abgeleitete `regex` matcht jede
  Quell-Kennung — kein ID-Shape bricht das Matching.
- **Determinismus** (`DC-QA-02`): alle Maps werden sortiert/geordnet
  ausgegeben (`sort.Strings` in Extraktion und Ableitung, feste
  Modul-Reihenfolge).
- **`regex: '%s'`-Quoting** sicher: der Regex-Körper enthält nie ein
  einfaches Anführungszeichen (`QuoteMeta` über alphanumerische Präfixe).

## Ergebnis

5 echte Defekte/Robustheits-Punkte (1–5) als Review R1 gefixt, Tests
ergänzt, `make gates` grün. 5 Edge-/Cleanup-Punkte (6–10) bewusst
akzeptiert und hier dokumentiert.
