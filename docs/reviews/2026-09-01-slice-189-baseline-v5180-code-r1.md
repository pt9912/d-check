# Review-Report: slice-189 — 2026-09-01

**Review-Art:** Code — geprüft gegen Plan (`slice-189-baseline-v5180.md`),
`AGENTS.md` §3 und die aktiven `MR-*`-Adaptionen.

**Gegenstand:** `fe00f17..HEAD -- . ':(exclude).harness/baseline/**'` — ein
einzelner Commit `55f3639` (feat(harness): Baseline-Pin auf v5.18.0).

**Skill:** `.harness/skills/reviewer.md` @ v1.13.0
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-01

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-189-baseline-v5180.md` (vollständig gelesen)
- `harness/conventions/MR-013-lifecycle-move-buendelung.md`,
  `harness/conventions/MR-058-baseline-v5180.md`,
  `harness/conventions/done/MR-057-baseline-v5150.md`
- `docs/plan/planning/observations.md` (`BEO-008`)
- `AGENTS.md` §3.3, §3.4, §3.5, §3.7
- `make gates` (Live-Lauf, grün), `make doc-check` (Live-Lauf, grün)

---

## Findings

### F-1 — Externe Provenienz-URLs (Release-/Tree-Links) nicht auf v5.18.0 gezogen — BEO-008, vierte/fünfte Instanz

- `kategorie`: MEDIUM
- `quelle`: `BEO-008` (Beobachtungs-Register, Zähler 4 vor diesem Slice), DoD-Punkt 1 des Slice-Plans
- `pfad`: `AGENTS.md:36`, `harness/conventions.md:41` und `:46`, `harness/README.md:60`, `harness/conventions/MR-021-vendored-verweise-pin-gebunden.md:20`
- `befund`: Fünf lebende Vorkommen von `v5.15.0` außerhalb des `.harness/baseline/v5.15.0/`-Pfadmusters sind stehengeblieben: die externe Kurs-Tree-URL
  (`ai-harness-course@v5.15.0`) und die `lab-regelwerk.zip`-Release-Download-URL
  (dreifach: `AGENTS.md`, zweimal `harness/conventions.md`, zusätzlich
  `harness/README.md`) verweisen weiterhin auf das **superseded** v5.15.0-Release-Asset,
  nicht auf v5.18.0 — obwohl `tools/harness/fetch-baseline-cache.sh` den Download
  strikt tag-parametrisiert (`releases/download/${tag}/lab-regelwerk.zip`), also
  ein anderes Artefakt ist als das tatsächlich vendorte. Dazu `MR-021` selbst:
  *„Diese Links tragen den konkreten Pin (aktuell `…/v5.15.0/…`)"* — die eigene
  Beispiel-Aussage der Adaption ist jetzt falsch. Genau das ist BEO-008s
  benannte zweite Spiegel-Klasse (Release-/Tree-URLs, „gehoben wird nur die
  grep-bare"): Das Slice-Vorgehen (§2 Punkt 2) misst ausdrücklich nur
  `baseline/v5.15.0` — exakt das Muster, das BEO-008 als unvollständig
  benennt — und §7 des Slice-Plans sichtet BEO-008 selbst als offene
  Beobachtung, ohne die dort empfohlene Gegenmaßnahme (Drei-Klassen-Zensus)
  in die eigene Messmethode zu übernehmen. Das DoD-Statement „alle 44 gehoben,
  keines stehen gelassen" überdehnt damit die tatsächlich gelaufene Messung
  (Q8: Botschaft verallgemeinert über die gelaufene Messung hinaus). Historisch
  identische Klasse bereits in `slice-128`-Review F-1 gefunden und behoben
  (`harness/README.md:60`, damals `v5.9.0`) — dieselbe Zeile ist jetzt erneut
  betroffen.
- `verifizierbar`: ja — `grep -rn "v5\.15\.0" AGENTS.md harness/conventions.md harness/README.md harness/conventions/MR-021*.md` zeigt alle fünf Fundstellen; kein Gate meldet sie (bestätigt: `make gates`/`make doc-check` laufen grün trotz der Stale-Links).
- `klasse`: Pin-Hebung übersieht Release-/Tree-URL-Spiegelklasse (BEO-008)

### F-2 — Neuer Tombstone-Kommentar trägt Herkunfts-Prosa statt nur Zusage/Abgrenzung

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (fünf Kommentar-Klassen; „keine … Herkunfts-Prosa")
- `pfad`: `.d-check.yml:159-163`
- `befund`: Der neue Kommentarblock zum vierten Tombstone-Eintrag erzählt den
  Hergang des Fehlers („der urspruengliche Bulk-Zensus hatte … nicht
  ausgenommen und … versehentlich … gehoben") statt nur den gegenwärtigen
  Zustand/die Abgrenzung zu benennen („MR-038/MR-041 sind bei v5.12.0
  eingefroren; diese Zeile hält sie davon aus, target-missing zu werden").
  Das ist dieselbe Prosa-Form, die §3.7 als HIGH ausschließt, hier aber im
  Umfang milder (ein Nachbar-Kommentar im selben Block, für den dritten
  Eintrag, trägt eine ähnliche „Erstmals …"-Erzählung und ist grandfathered,
  weil er vor dieser Schärfung entstand — der neue Kommentar fällt unter den
  Anker, weil er ein **Neuzugang** ist).
- `verifizierbar`: nein — kein Gate prüft die fünf Kommentar-Klassen (Urteil, kein `grep`, wie im Skill vermerkt).
- `klasse`: Konfigurationskommentar trägt Bug-Historie statt Zustand/Abgrenzung

### F-3 — MR-057-Move + gesamter Bump-Rest in einem Commit statt zwei

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.3, `MR-013` (Ausnahme MR-/Wellen-Lifecycle-Move: „Alles Übrige bleibt Commit 2")
- `pfad`: Commit `55f3639` (gesamter Diff)
- `befund`: `MR-013`s Ausnahme erlaubt, dem Move-Commit von
  `MR-057-baseline-v5150.md` → `conventions/done/` die Link-Tiefen-Fixes der
  bewegten Datei selbst mitzugeben — verlangt aber ausdrücklich, dass „alles
  Übrige" in einem zweiten Commit bleibt. Hier ist der komplette restliche
  Bump (52 weitere Dateien: Pfad-Bumps, `MR-058` neu, Tombstone-Ergänzung,
  Slice-Body) im selben Commit wie der Move. Kein Gate prüft die
  Commit-Zerlegung (`AGENTS.md` §3.3 vermerkt das selbst). Derselbe Zuschnitt
  findet sich unverändert in den beiden Vorgänger-Bumps
  (`08373c9`/slice-183, `9ee805b`/slice-148) — kein neues Muster dieses
  Slice, sondern eine seit mehreren Bumps unadressierte Lücke zwischen
  Hard-Rule-Text und gelebter Praxis.
- `verifizierbar`: nein (kein Gate; Beobachtung am Commit-Log).
- `klasse`: Lifecycle-Move-Commit bündelt mehr als die benannte Ausnahme erlaubt

## Negativbefunde

- geprüft, ohne Befund: §2a-Argumentation (MR-013/MR-057-Kollisionsauflösung) — gegen `MR-013`s eigenen Text nachgelesen; die Zweiteilung (Push-CI-Begründung für Slice-Lifecycle-Move/Beanspruchung vs. eigenständige, lokale `doc-check`-Begründung für MR-/Wellen-Lifecycle-Move) ist in MR-013s Text tatsächlich so angelegt, nicht nachträglich hineininterpretiert — die beiden Begründungs-Absätze in MR-013 selbst sind bereits getrennt formuliert.
- geprüft, ohne Befund: Bulk-Zensus-Fehler-Korrektur — `git log --follow -p` für `MR-038-zitate-pin-gebunden.md` und `MR-041-guard-node-und-eigene-toolchain.md` bestätigt: beide wurden mit `v5.12.0` angelegt, von `08373c9` (slice-183) versehentlich auf `v5.15.0` gehoben, und von `55f3639` (slice-189) korrekt auf `v5.12.0` zurückgesetzt. Der neue `.d-check.yml`-Tombstone-Eintrag `harness/conventions/done/** → v5.12.0/**` deckt beide betroffenen Dateien vollständig ab und ist korrekt in den bestehenden `v5.12.0`-Tombstone-Block einsortiert.
- geprüft, ohne Befund: Pfad-Konsistenz der `.harness/baseline/v5.15.0`-Pfadmuster — `grep -rn "baseline/v5\.15\.0"` außerhalb `.harness/baseline/v5.15.0/` selbst trifft nur noch den Tombstone-Eintrag in `.d-check.yml` und die eingefrorene `docs/plan/planning/done/slice-183-baseline-v5150.md`; die eine verbleibende Erwähnung in `slice-189` selbst ist eine historische Ist-Stand-Aussage über den Zeitpunkt der Anlage, kein toter Verweis.
- geprüft, ohne Befund: vier `d-check:cite`-Direktiven gegen ihre Quelle nachgemessen — `MR-058` (`grundlagen-traceability.md:113-115`), `AGENTS.md` (`modul-05-planning-harness.md:136→142`), `slice-189` §7 (`modul-05-planning-harness.md:219-220`, `:225`) — alle vier Wortlaute halten exakt gegen `.harness/baseline/v5.18.0/`.
- geprüft, ohne Befund: `git show --find-renames` für `MR-057-baseline-v5150.md` → `conventions/done/` bestätigt `R74` (über der 50 %-Schwelle); Rename-Detection hält, kein Wortlaut außerhalb der Pfad-Tiefen-Korrektur verändert.
- geprüft, ohne Befund: Stichprobe von fünf lebenden `MR`-Dateien (`MR-004`, `MR-013`, `MR-021`, `MR-033`, `MR-049`) — alle `Ersetzt-Baseline-Regel`-Anker lösen in `.harness/baseline/v5.18.0/` korrekt auf denselben, unveränderten Abschnitt auf.
- geprüft, ohne Befund: `AGENTS.md` §3.4/§3.5/§3.7 sowie Source Precedence — keine ADR-Datei im Diff berührt, keine Spec-Stratum-Inhaltsänderung außer Link-Version-Bumps, keine Abwärts-Referenz neu eingeführt.
- geprüft, ohne Befund (Live-Lauf): `make gates` — zehn Gates, 649 Dateien, 0 Befunde, Exit 0. `make doc-check` — 649 Dateien, 0 Befunde, Exit 0.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Pin-Hebung übersieht Release-/Tree-URL-Spiegelklasse (BEO-008) · Konfigurationskommentar trägt Bug-Historie statt Zustand/Abgrenzung · Lifecycle-Move-Commit bündelt mehr als die benannte Ausnahme erlaubt

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, ein MEDIUM ohne Gate-Bezug (reine
Doku-Korrektheit dreier externer Provenienz-Links, keine Auswirkung auf
Integrität/Funktion des vendorten Baums selbst). Empfehlung: F-1 vor Closure
nachziehen (fünf Zeilen), da es exakt die in `BEO-008` benannte, bereits
einmal (slice-128) gefundene und dokumentierte Fehlerklasse ist und der
Zähler damit eine weitere, vermeidbare Instanz bekäme. F-2/F-3 sind
LOW/Konsistenz-Beobachtungen ohne Blockierungs-Anspruch.

**Übergabe:** Findings gehen an den Implementer. Die Finding-Klassen gehen
in die Slice-Closure §7 und von dort in `docs/plan/planning/observations.md`
(F-1 zitiert `BEO-008`, dessen Zähler bereits bei 4 steht — Beleg für die
fünfte Instanz wäre `slice-189` selbst, sobald geschlossen). Dieser Report
ersetzt keine Verifikation.
