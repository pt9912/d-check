# Review-Report: slice-112 — Wellen-Invariante als Kennungs-Mengen (Lastenheft 0.62.1) + Roadmap-Entscheid Mehr-Wellen-Betrieb

**Datum:** 2026-08-22 · **Review-Art:** Design/Doku (Wortlaut-Präzision gegen den Code als Wahrheit des Ist-Verhaltens; geprüft gegen Lastenheft 0.62.1 §Wellen-Invariante, `spec/spezifikation.md` §`DC-FA-PLAN-001.a` W2/W3/§2/§4, ADR-0055 Entscheidung 6, Hard Rules `AGENTS.md` §3.4/§3.5/§3.7, `MR-025`) und Code (Pinning-Test)
**Gegenstand:** Commit-Range `86a9afc..9a92b7c` — `4a6506f` (Roadmap-Entscheid, eigener Auftraggeber-Entscheid, nur Kennungsfreiheit/Marker-Wortlaut geprüft) → `4877d8b`/`605f55c` (slice-112 eröffnet/beansprucht) → `ae41ca8` (Lastenheft 0.62.1) → `ce750d4` (Spezifikation + ADR-0055) → `33cb425` (Pinning-Test) → `9a92b7c` (Handbuch 1.54 + CHANGELOG), **vor** der Closure, kein Release
**Skill:** `reviewer.md` @ 1.5.0 (Stand `f845e8b`) · **Modell-ID:** `claude-fable-5`
**Eingangs-Kontext:** Slice-Plan `docs/plan/planning/in-progress/slice-112-wellen-invariante-kennungs-menge.md` (§2 Schritt 3 = MR-025-Spiegel-Liste, §3 NICHT-Liste, §5 Risiken); `DC-FA-PLAN-001` §Wellen-Invariante (Zeile 1/2 + Akzeptanzkriterien + §7 oberste Zeile); `spec/spezifikation.md` W2/W3 + §4 `wave-drift`; [ADR-0055](../plan/adr/0055-wellen-invariante-artefakt-und-grund-codes.md) Entscheidung 6 + Geschichte; `internal/hexagon/core/rules/planning_waves.go` (`waveSets`, `waveDrift`, `waveBijection`) als Ist-Wahrheit; Vorgänger-Review `docs/reviews/2026-08-21-slice-111-waves-mode-review.md` F-5 (Anlass). Nicht erhalten: DoD-Abhakung (Verifikations-Rolle).

## Findings

### F-1 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** `MR-025` (Spiegel ist die Stelle, nicht die Datei) / `DC-FA-PLAN-001`
- **pfad:** `spec/lastenheft.md:2186` (Absatz „Zwei Kardinalitäts-Modelle, ein Prädikat"); `docs/user/benutzerhandbuch.md:1238-1239` (§6-Absatz „Dritte Fähigkeit, opt-in über `waves.dir`")
- **befund:** In den beiden bearbeiteten Dateien steht je **eine weitere** Stelle derselben Semantik, die weiterhin die Datei-Lesart behauptet: das Lastenheft sagt im Prosa-Absatz unterhalb der geänderten Tabelle für `one` „verglichen wird gegen **genau ein** flaches Dokument", das Handbuch im Absatz direkt **über** dem geänderten sagt „der Aktiv-Status-Abschnitt nennt eine Welle genau dann, wenn **genau ein** flaches Wellendokument liegt" — und der geänderte Absatz verweist mit „das Singleton-Prädikat oben" genau auf diesen Satz. Failure-Szenario: ein Leser des Lastenheft-Absatzes (oder des Handbuch-Absatzes) erwartet im Doppel-Dokument-Fall unter `one` einen `wave-drift`, den weder Zeile 1/2, noch das neue Akzeptanzkriterium, noch der Code liefern — die Präzisierung hebt den Widerspruch innerhalb derselben Anforderung auf, statt ihn zu beseitigen.
- **verifizierbar:** nein (kein Gate misst Aussagen gegen Aussagen — genau die `MR-025`-Klasse); Beleg sind die drei Zitate
- **klasse:** spiegel-stelle-nicht-datei

### F-2 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** `MR-025` (Spiegel-Zeile „Config-Schema — `spec/spezifikation.md` §2")
- **pfad:** `spec/spezifikation.md:2529` (§2-Schema-Zeile `planning.waves.mode`)
- **befund:** Die §2-Schema-Zeile beschreibt `one` als „`hasActive` ⇔ genau ein flaches Dokument" und `many` als „Kennungs-Mengen-Gleichheit `planning.heading`-Block ⇔ flache Dokumente" — beides die Datei-Lesart, die W3 und §4 derselben Datei seit `ce750d4` nicht mehr tragen. Die Spiegel-Liste des Slice (§2 Schritt 3) nennt für die Spezifikation nur W3, §4, §7 und Kopf; die Config-Schema-Zeile ist eine eigene Zeile der `MR-025`-Tabelle und fehlt in der Liste. Failure-Szenario: wer den Schlüssel aus dem Schema heraus konfiguriert, liest eine andere Zusage als die, die der Algorithmus-Schritt derselben Spezifikation gibt.
- **verifizierbar:** nein (Prosa-Spiegel); Beleg ist das Zitat gegen W3 (`spec/spezifikation.md:1787-1791`)
- **klasse:** mr-025-spiegel-zeile-fehlt

### F-3 · LOW

- **kategorie:** LOW
- **quelle:** `MR-025` (Spiegel-Zeile „Autoritäts-Doku") / Maintainability (Doku-Drift)
- **pfad:** `AGENTS.md:113` („`wave-drift` misst unter `mode: many` Zeiger ⟺ Dateien in beide Richtungen"); `README.md:98-99` („identifier bijection active block ⟺ files"); `README.de.md:97-99` („Kennungs-Bijektion Aktiv-Block ⟺ Dateien"); `docs/user/benutzerhandbuch.md:1776` (§6-Modul-Tabelle: „Kennungs-Bijektion Aktiv-Block ⟺ flache Dokumente")
- **befund:** Vier Kurzform-Spiegel setzen die rechte Seite der Bijektion mit „Dateien"/„Dokumenten" gleich; nach 0.62.1 ist sie die Kennungs-Menge. Keine der Stellen zählt („genau ein"), die Drift ist daher eine Lesart-, keine Kardinalitäts-Aussage — im Doppel-Dokument-Fall bleibt das Gate-Ergebnis mit „je Datei ein Zeiger" sogar grün (Mehrfachnennung zählt einmal). Failure-Szenario ist die falsche Erwartung, nicht ein falscher Befund: ein Adopter mit zwei gleich-kennigen Dateien und **einem** Zeiger liest „⟺ Dateien" und hält den ausbleibenden Befund für ein Loch.
- **verifizierbar:** nein (Prosa); bestätigt per `grep -n "wave-drift" AGENTS.md README.md README.de.md` und Handbuch-Zeile 1776
- **klasse:** kurzform-spiegel-datei-lesart

### F-4 · INFO

- **kategorie:** INFO
- **quelle:** `DC-FA-PLAN-001` / Slice §3 (bewusste Abgrenzung)
- **pfad:** `internal/hexagon/core/rules/planning_waves.go:172-177` (Meldungen des `one`-Pfads); `internal/hexagon/core/app/diagnose.go:134` (`--doctor`-Klartext `wave-drift`)
- **befund:** Die Klartexte sind — wie in Slice §3 angekündigt — unverändert: „liegen N flache Wellendokumente — genau eines ist erwartet" bzw. „liegt N flaches Wellendokument". N ist `len(flach)`, also die **Kennungs**-Zahl (so sagt es jetzt auch W3: „bei mehr als einer Kennung deren Zahl"), das Nomen ist aber „Wellendokument". Im Regelfall (eine Kennung je Datei) gibt es keinen Widerspruch; im Doppel-Dokument-Fall (drei Dateien, zwei Kennungen, Aktiv-Status genannt) sagt die Meldung „liegen 2 flache Wellendokumente — genau eines ist erwartet", und das Entfernen **einer** Datei des Paars ändert nichts. Die Abgrenzung des Slice trägt (kein Klartext behauptet das Datei-Prädikat als Regel; eine Textänderung wäre Verhaltensänderung und Release-Anlass) — die Label-Drift gehört zu der in Lastenheft 0.62.1 ausdrücklich offen gelassenen Doppel-Dokument-Frage und sollte mit ihr entschieden werden.
- **verifizierbar:** ja — ein Fixture mit drei Dateien/zwei Kennungen unter `one` zeigt den Wortlaut
- **klasse:** klartext-label-vs-zaehlgroesse

## Negativ-Proben (geprüft, ohne Befund)

Alle Kommandos selbst gelaufen; die Gate-Ausgabe in eine Datei umgeleitet, Exit explizit geprüft.

1. **Gate:** `make test > make-test.log 2>&1; echo EXIT=$?` → **EXIT=0**; das Log zeigt `RUN CGO_ENABLED=0 go test ./...` und `ok github.com/pt9912/d-check/internal/hexagon/core/rules` — der neue Test ist kompiliert und gelaufen (`maps`-Import gegen `go 1.25.0` aus `go.mod`). Kein `make gates`, kein Runtime-Image-Build (`make test` baut nur die Test-Stufe `d-check:test`), nichts committet.
2. **Leitfrage 1 — Code-Wahrheit:** `waveSets` (`planning_waves.go:127ff`) füllt `flach` über `waveID(prefix, name)` — zwei Basisnamen derselben Kennung schreiben denselben Schlüssel; `waveDrift` (`:161-164`) prüft `hasActive == (len(flach) == 1)`, `waveBijection` (`:67-78`) vergleicht `ids` gegen `flach` per Kennung. Beide Modi zählen also Kennungen; der neue Wortlaut von Lastenheft Zeile 1/2, Akzeptanzkriterium, Spezifikation W3 + §4, ADR-0055 Entscheidung 6, Handbuch-§6-Absatz, §11-Zeile und CHANGELOG sagt an allen **geänderten** Stellen dasselbe (Kennungs-Menge, Doppel-Dokument = ein Element, auch unter `one`).
3. **Leitfrage 2 — Residual-grep:** `grep -rn -i "genau ein flaches\|Menge der flachen Wellendokumente\|Datei-Menge\|Datei-Zahl\|gegen die flachen Wellendokumente\|gegen die flachen Dokumente" --include=*.md --include=*.yml --include=*.go .` (ohne `docs/reviews/`, `.harness/baseline/`, `.git/`) → Treffer ausschließlich (a) Historie-Zeilen, die den 0.62.0-Stand korrekt datieren (Lastenheft §7 0.62.0, Spezifikation §7 2026-08-21, Handbuch §11 1.53, CHANGELOG [0.62.0]), (b) geschlossene Slices/Wellen unter `done/` und Reviews (zeitliche Schicht, nicht zu redigieren), (c) das Akzeptanzkriterium „Wellen-Happy-Path" (dort ist „genau ein flaches Wellendokument" die Fixture-Angabe, nicht das Prädikat) und (d) die in F-1 bis F-3 gemeldeten Stellen. Die eigenen Config-Kommentare (`.d-check.yml:261-272`, `:284-288`) sagen „Kennungs-Bijektion (Zeiger <=> flache Wellendokumente)" — die Slice-Angabe „bleibt korrekt" trägt.
4. **Leitfrage 2 — Klartexte nicht angefasst:** `git diff --quiet 86a9afc..9a92b7c -- internal/hexagon/core/rules/planning_waves.go internal/hexagon/core/app/diagnose.go` → **Exit 0**; `git diff --name-only` der Range zeigt als einzige `.go`-Datei `planning_waves_test.go`. Kein Produktions-Code, kein Grund-Code, kein `--doctor`-Text geändert (Slice §3 eingehalten). Ob ein Klartext dem neuen Wortlaut **aktiv** widerspricht: nein im Regelfall, Label-Drift im Doppel-Fall (F-4).
5. **Leitfrage 3 — Test-Aussagekraft, statisch am Code:** Fixture `welle-9-neu.md` + `welle-9-nachtrag.md` ⇒ `flach = {welle-9}`. Subtest `one`: Block „welle-9-neu in Arbeit." ohne Marker ⇒ `hasActive = true`, `len(flach) == 1` ⇒ 0 Befunde. Subtest `many`: `planningBlockIDs` liest `welle-9` aus „- `[welle-9-neu](../welle-9-neu.md)`", Marker daneben geht nicht ein ⇒ beide Differenzen leer ⇒ 0 Befunde. Die gemeldete Mutation `flach[id]`→`flach[name]` ist plausibel rot in beiden Subtests: unter `one` wird `len(flach) = 2 ≠ 1` (Befund), unter `many` fehlt `welle-9` in der Namens-Map und beide Namen fehlen in `ids` (drei Befunde). Zusatz-Beobachtung: `TestWavesHappyPath` (eine Datei) würde dieselbe Mutation unter `one` **nicht** sehen — der neue Test fügt für den `one`-Pfad echte Trennschärfe hinzu, nicht nur eine Wiederholung. `coretest.MemFS.List` (`memfs.go:73ff`) liefert für fehlende Verzeichnisse keinen Fehler ⇒ kein verdeckter `dirErrs`-Befund im Fixture (derselbe Aufbau wie `TestWavesDrift`).
6. **Leitfrage 4 — Hard Rules:** §3.4: `git diff 86a9afc..9a92b7c -- spec/ | grep "^+" | grep -q "slice-\|ADR-\|welle-[0-9]\|[0-9a-f]\{7\}"` → **Exit 1** (kein Treffer; die Lastenheft-Historie sagt „Review-Hinweis zur 0.62.0-Erweiterung", kein Slice/ADR/Hash). §3.5: ADR-0055 trägt `**Status:** Proposed` (`:3`) — der Core-Edit in Entscheidung 6 ist zulässig, die Fortschreibung steht zusätzlich als `## Geschichte`-Zeile (`:175-179`), Status unverändert. §3.7: der Kommentar des neuen Tests (`planning_waves_test.go:482-484`) ist eine **Zusage** mit **einem** Herkunfts-Feld (`DC-FA-PLAN-001` + Kriteriums-Name) — kein Befund-Marker, keine Slice-Nummer, keine Review-Historie; Stil wie die Nachbar-Kommentare (`:322-323`).
7. **Leitfrage 5 — Redaktion:** Lastenheft `**Version:** 0.62.1` (`:3`), Historie-Zeile 0.62.1 ist die oberste Datenzeile über 0.62.0 (`:2816-2817`, newest first); Spezifikation Kopf „Letzte Änderung: 2026-08-22" und §7-Zeile 2026-08-22 oben (`:2698`); Handbuch Kopf 1.54/Stand 2026-08-22/Software-Version v0.62.0 (Anker `version.md#v0.62.0` existiert, `version.md:35`), §11-Zeile 1.54 **unter** 1.53 als letzte Zeile (chronologisch); CHANGELOG `[Unreleased]` → `### Changed` → „slice-112 — **…**" in derselben Form wie der `[0.62.0]`-Eintrag (Keep-a-Changelog-konform, kein Tag).
8. **Leitfrage 6 — Roadmap-Commit `4a6506f`:** `git show --stat` → nur `docs/plan/planning/in-progress/roadmap.md` (9+/1−). §Offene-Wellen-Prosa (Zeilen 12-36) gegen `welle-[0-9]` → **Exit 1** (kennungsfrei); gegen den Marker-Wortlaut „Nichts in Arbeit" → **Exit 1** (nicht literal; die beiden Treffer der Datei liegen in §Historie-Zeilen 115/117 außerhalb des Blocks). Der Move-Commit `605f55c` nimmt den Ruhe-Marker (slice-112 beansprucht) — konsistent zur Marker-Hälfte; `ls docs/plan/planning/welle-*.md` → 0 flache Wellendokumente, die Block-Prosa nennt keine Kennung ⇒ Bijektion leer in beide Richtungen.
9. **Commit-Schnitt:** `git show --stat` je Commit — `4877d8b` nur die Slice-Datei (open/), `605f55c` reiner `git mv` + zwei Roadmap-Zeilen (MR-013-Muster), `ae41ca8` nur `spec/lastenheft.md` (4+/2−), `ce750d4` nur Spezifikation + ADR-0055, `33cb425` nur der Test (26+), `9a92b7c` nur Handbuch + CHANGELOG. Lastenheft liegt als erster Spec-Commit vor Spezifikation/ADR/Test (Doc führt).
10. **Spiegel-Liste des Slice gegen den Diff:** alle in §2 Schritt 3 gelisteten Stellen sind im Diff geändert (Lastenheft Zeile 1/2 + AK + §7; Spezifikation W3 + §4 + §7 + Kopf; ADR-0055 E6 + Geschichte; Handbuch §6-Absatz + Kopf + §11; CHANGELOG). Die Liste war vollständig **abgearbeitet**, aber nicht vollständig **abgeleitet** — F-1 bis F-3 sind Stellen, die der `MR-025`-Ableiter per `grep` nach dem vorigen Wortlaut („genau ein flaches") gefunden hätte.
11. **ADR-0055 übrige Entscheidungen:** Entscheidung 4 („Aktiv-Status gegen Plan-Dokument") nennt die Rolle, keine Kardinalität; Entscheidungen 1–3, 5, Alternativen und Fitness-Function-Teile tragen keine Datei-Mengen-Aussage — kein weiterer ADR-Spiegel offen.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 2 | F-1 (Lastenheft-Prosa + Handbuch-Absatz „oben" sagen noch Datei-Lesart), F-2 (Spezifikation §2-Schema-Zeile fehlt in der Spiegel-Liste) |
| LOW | 1 | F-3 (Kurzform-Spiegel AGENTS/README/README.de/Handbuch-Modul-Tabelle „⟺ Dateien") |
| INFO | 1 | F-4 (Klartexte: Kennungs-Zahl mit Dokument-Nomen, bewusst ausgenommen) |

## Verdikt

**APPROVE mit Auflagen.**

Der Kern stimmt: der präzisierte Wortlaut deckt sich an jeder geänderten Stelle mit dem Code (Kennungs-Map in beiden Modi), der Pinning-Test trägt die Aussage und ist für den `one`-Pfad trennschärfer als der bestehende Happy-Path, die Hard Rules §3.4/§3.5/§3.7 sind eingehalten, die Redaktion (Versionen, Historien-Positionen, CHANGELOG-Form) ist korrekt, und der Roadmap-Entscheid ist kennungsfrei und marker-frei. Kein HIGH.

Auflagen vor der Closure: **F-1** und **F-2** beheben — drei Sätze (Lastenheft-Absatz Zeile 2186, Handbuch-Absatz Zeile 1238-1239, Spezifikation §2-Zeile 2529), die dieselbe Semantik noch in der Datei-Lesart führen; ohne sie widerspricht sich das Lastenheft innerhalb der Anforderung, und die Spezifikation zwischen Schema und Algorithmus — das ist genau die Lücke, die der Slice schließen soll. Die Lehre gehört in die Closure-Notiz: die `MR-025`-Liste wurde aus dem Slice-Ziel abgeleitet, nicht per `grep` nach dem **alten** Wortlaut über den Baum; das `grep` hätte F-1 bis F-3 geliefert. **F-3** ist nicht blockierend (Kurzformen, keine Kardinalitätsaussage) und kann mit derselben Redaktion oder der nächsten Release-Prep mitgehen. **F-4** bleibt bewusst offen und ist zusammen mit der Doppel-Dokument-Frage zu entscheiden (eigener CR, wie in Lastenheft 0.62.1 benannt).
