# Review-Report: slice-023 — `ids`-Ventile für nackte Vorkommen — 2026-06-16

**Review-Art:** Code-Review der Implementierung über den Working-Tree-Diff
(`/code-review high`, 8 Finder-Angles → 4 gebündelte Finder → Synthese).
Gegenstand ist die Umsetzung von slice-023, nicht erneut der Vertrag.

**Gegenstand:** uncommitteter Diff (8 Dateien + neuer Slice):
`internal/hexagon/core/rules/ids.go`, `ids_test.go`, `config.go`,
`spec/lastenheft.md`, `spec/spezifikation.md`, `docs/user/operations.md`,
`CHANGELOG.md`, `docs/plan/planning/in-progress/`.

**Skill:** `.harness/skills/reviewer.md` · **Datum:** 2026-06-16

**Eingangs-Kontext:** `DC-FA-ID-001` (CR 0.13.0), `DC-QA-02`/`DC-QA-03`,
`AGENTS.md` §3, `MR-006`, slice-023.

---

## Findings

| # | Kategorie | Quelle | Pfad | Befund | Verifizierbar |
|---|---|---|---|---|---|
| 1 | 🟡 MEDIUM | `DC-FA-ID-001` / fehlender Negativtest bei neuem Vertrag | `internal/hexagon/core/rules/ids_test.go` | Die Politik-Unabhängigkeit des `d-check:ignore`-Ventils für nackte Vorkommen (Vertrag: gilt auch unter Default `prose`) war untestet; `exempt-paths` hatte einen prose-Default-Test, der Marker nur einen `always`-Test. Eine Regression, die den Marker politik-gated macht, wäre durch die Suite gekommen. | ja — `make test` mit Marker-Zeile unter Default-Politik |
| 2 | 🟢 LOW | Maintainability (Altitude) | `internal/hexagon/core/rules/ids.go:73` | Die Zeilen-Ausnahme „Zeile trägt `d-check:ignore`" wurde an zwei Stellen über dasselbe Prädikat geprüft (Prosa-Pfad via Map in `checkIDs`, Inline-Code-Pfad via `strings.Contains` in `alwaysLineFindings`). Eine spätere Änderung der Marker-Semantik müsste beide Stellen treffen — Divergenz-Risiko. | ja — Code-Lesung; `make test` deckt beide Pfade |
| 3 | 🟢 LOW | `DC-FA-ID-001` / Doku-Drift | `spec/lastenheft.md:315` | Der `always`-Aufzählungssatz der linkpflichtfreien Fälle (Fence/Heading/`target`) las als erschöpfend, nannte die beiden Ventile aber nicht — ein Leser konnte folgern, ein Backtick-Vorkommen in einer `exempt-paths`-Datei sei unter `always` weiter meldepflichtig. | ja — Doc-Lesung gegen §`DC-FA-ID-001.a` |
| 4 | 🟢 LOW | Maintainability (Perf) | `internal/hexagon/core/rules/ids.go:128` | `ignored(file, p.ExemptPaths)` wurde je ID-Vorkommen im inneren Match-Loop ausgewertet, obwohl es nur von (Datei, Muster) abhängt — anders als das bereits gehoistete `inTarget`. | ja — Code-Lesung |
| 5 | 🟢 LOW | `DC-FA-ID-001` / Doku-Drift | `docs/user/operations.md:65` | Der (korrekt politik-unabhängig formulierte) Ventil-Block stand allein unter der Überschrift „`link-policy: always`"; ein Nutzer unter Default `prose` (genau das Repro-Szenario) hätte ihn dort nicht gesucht. | nein — Struktur-/Auffindbarkeitsnotiz |
| 6 | 🔵 INFO | `DC-QA-01` / Maintainability (Perf) | `internal/hexagon/core/rules/ids.go:147` | `proseLines(content)` lief im `ids`-Modul zweimal (Marker-Scan + `checkIDsAlways`); der `PreprocessMarkdown`-Aufruf (`run.go`) ist ein dritter, geteilter Lauf der gesamten Pipeline. | ja — `make bench`, kein Gate |

## Negativbefunde (geprüft, ohne Befund)

- **Korrektheit (Finder-Angle A/B/C):** Zeilen-Nummerierung deckt sich —
  `markerLines`-Schlüssel und `Line.No` stammen beide aus `proseLines`
  (`i+1`); `nil`-Map-Lesung ist sicheres `false` (Tests mit `content=nil`);
  Prosa- und Inline-Code-Pfad sind disjunkt (kein Doppelzählen);
  `claimed`-Überlappung/Muster-Präzedenz unverändert. **Kein Korrektheitsbug.**
- **Konventionen (Hard Rules):** CR nur im Lastenheft; AK-Trio +
  Out-of-Scope + Versions-Bump 0.12.0→0.13.0 + §7-Historie vorhanden;
  keine Abwärts-Referenz aus einem Spec-Stratum (`MR-006`/§3.4) — die
  spezifikation.md-Ergänzungen tragen keinen ADR/Slice/Welle/Datum-Token;
  Slice korrekt in `in-progress/`, Gates ehrlich offen. **Kein Verstoß.**
- **Determinismus (`DC-QA-02`):** `ignored`/`matchGlob` unverändert
  wiederverwendet (reihenfolge-/plattformstabil); `markerLines`-Map ist
  Mengen-Mitgliedschaft, nicht ausgabewirksam iteriert.
- **Abwärtskompatibilität:** Configs ohne gesetzte Ventile —
  `ignored(file, nil)` ⇒ `false`, leere Marker-Menge — sind byte-identisch;
  die Erweiterung entfernt Befunde nur, fügt keine hinzu.
- **Dogfooding (`DC-QA-03`):** netzloser `doc-check`-Lauf über das eigene
  Repo grün (52 Dateien / 0 Befunde), inkl. der exempten Review-Reports.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 1 | 4 | 1 |

## Verdikt

**Kein HIGH; ein MEDIUM (Testlücke) → vor Closure zu schließen.** Korrektheit
und Konventionen sind sauber; die übrigen Befunde sind Doku-Drift und
Qualitäts-Härtung am offenen Code. `make gates` ist auf dem Review-Stand
grün (Coverage 94,80 %).

## Disposition (Review R1 — 2026-06-16)

- **F1 — gefixt:** `TestIDsIgnoreMarkerProseDefault` — `d-check:ignore`
  nimmt eine nackte ID auch unter Default `prose` aus (Gegenstück zu
  `TestIDsExemptPathsProseDefault`).
- **F2 — gefixt:** Die Zeilen-Ausnahme lebt jetzt in **einem** Satz
  (`markerLines` → `ignoreLines`-Map), den beide Pfade konsultieren;
  `alwaysLineFindings` prüft `ignoreLines[pl.no]` statt eigenem
  `strings.Contains`.
- **F3 — gefixt:** Der `always`-Aufzählungssatz nennt nun „… sowie die
  beiden unten genannten Ventile (`exempt-paths`, `d-check:ignore`)".
- **F4 — gefixt:** `exemptFile[i] := ignored(file, p.ExemptPaths)` wird
  einmal je Muster vorberechnet (wie `inTarget`) und im Loop indiziert —
  in beiden Pfaden.
- **F5 — gefixt:** operations.md nennt explizit, dass beide Ventile auch
  unter Default `prose` gelten (Review-Report-Beispiel).
- **F6 — teilweise/akzeptiert:** Die `ids`-internen `proseLines`-Läufe
  sind von zwei auf einen reduziert (Prosa einmal in `checkIDs` gewonnen,
  an `checkIDsAlways` durchgereicht); der `PreprocessMarkdown`-Lauf ist
  geteilte Pipeline-Infrastruktur, kein riskanter Eingriff zugunsten einer
  INFO-Notiz. Die `nil`-Lazy-Init in `markerLines` bleibt bewusst (spart
  die Map-Allokation für die markerfreie Mehrheit der Dateien).

`make gates` nach R1 grün (52 Dateien / 0 Befunde, Coverage 94,90 %).
