# Review-Report — slice-060 (`--doctor`-Klartext-Vollständigkeit) — R1

**Datum:** 2026-07-04
**Reviewer-Rolle:** unabhängig/adversarial (Code + Doku; nicht der Autor).
**Gegenstand:** slice-060 (welle-49), zwei Commits —
`1e5cf41` (doc-first: [slice-060](../plan/planning/done/slice-060-doctor-klartexte.md)
+ Roadmap-Flip) und `f1017f8` (feat:
`internal/hexagon/core/app/diagnose.go` + `diagnose_test.go` + DoD-Häkchen).
**Anforderungs-Anker:** [`DC-FA-CLI-007`](../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus),
[Spezifikation §`DC-FA-CLI-007.a`](../../spec/spezifikation.md#dc-fa-cli-007a--diagnose-modus)
Schritt 3, [Spezifikation §4](../../spec/spezifikation.md#4-grund--und-fehler-codes)
(Grund-Code-Tabelle), Modul-Spezifikationen `DC-FA-DIAG/VER/PIN/IMM/VCS/COMMITS/PLAN-001.a`.
**Baseline:** `.harness/skills/reviewer.md` v1.2.0 (Kategorien/Schema/Negativbefund-Pflicht),
`AGENTS.md` §3.
**NICHT geprüft:** DoD-Abhakung (Verifikations-Rolle), Gate-Lauf-Bestätigung
(nicht Reviewer-Rolle). Keine `make`-Läufe; ausschließlich lesende Verifikation.

**Verifikations-Schritte (lesend):**

- **Mengen-Abgleich §4 ↔ `AllReasons()`:** die §4-Tabelle
  (`spec/spezifikation.md:1401`–`1423`) trägt exakt **23** Codes; `AllReasons()`
  (`diagnose.go:70`–`81`) und `reasonTexts()` (`diagnose.go:87`–`113`) tragen
  exakt dieselben 23 (die sieben neuen Konstanten in
  `internal/hexagon/core/model/finding.go:18`–`24` string-identisch zu §4).
  Zeilen-Zählprobe: `grep -c '^| \x60[a-z0-9-]*\x60 |'` über die §4-Spanne = 23.
- **Historien-Probe der Defekt-Behauptung:** `git show v0.25.0:…/diagnose.go`
  ⇒ `AllReasons()` hatte 14 Einträge, während v0.25.0 laut CHANGELOG
  `diagram-id-undefined` einführte — „hinkt seit v0.25 hinterher" stimmt.
- **Parser gegen die echte Spec:** Überschrift `## 4. Grund- und Fehler-Codes`
  existiert genau einmal (Zeile 1395); der Zeilen-Regex
  `^\| \x60([a-z0-9-]+)\x60 \|` matcht innerhalb §4 (bis `## 5.`) genau die
  23 Code-Zeilen und nichts sonst (Header-/Trennzeilen und die
  §3-Konstanten-Tabelle matchen nicht: Großbuchstaben bzw. kein Backtick-Lead).
- **Pfad-/Kontext-Probe:** Paket `internal/hexagon/core/app` → 4×`..` = Repo-Wurzel;
  `go test` läuft mit cwd = Paket-Verzeichnis. Docker-`test`-Stage
  (`Dockerfile`, `FROM deps AS test; COPY . .`) kopiert die volle Wurzel nach
  `/src`; **kein** `.dockerignore` vorhanden ⇒ `spec/` liegt im Build-Kontext.
  Gleiches gilt für die `coverage`-Stage.
- **Konsum-Pfad:** `app.ReasonText` wird in
  `internal/adapter/driven/report/report.go:75` (Prosa-Diagnose) und `:207`
  (JSON/YAML-`reasonText`) gerufen — die sieben Codes fielen dort bis v0.37.0
  auf den fail-safe-Rohcode; die Slice-Behauptung stimmt.
- **Hard Rules:** kein `//nolint`, kein `d-check:ignore` im Diff; `os`-Import
  nur im Testfile; `.a-check.yml` schließt `**/*_test.go` aus (dokumentierte
  Skript-Parität) — der R4-Bann (`os` nur im fs-Adapter) betrifft Produktcode
  und bleibt unverletzt.

---

## Findings

### LOW-1 — Verriegelung greift nicht bei Format-Abweichung einer einzelnen künftigen §4-Zeile (Rest-Lücke des historischen Fehlermodus)

- **kategorie:** LOW
- **quelle:** Maintainability (latente Wartungsfalle im Verriegelungs-Test)
- **pfad:** `internal/hexagon/core/app/diagnose_test.go:71` (+ `:86`–`:96`)
- **befund:** Der Zeilen-Regex verriegelt nur Zeilen der exakten Gestalt
  `| `Code` |` (genau ein Leerzeichen nach dem Pipe, Zeichenvorrat
  `[a-z0-9-]`). Die fail-closed-Guards (fehlende/mehrdeutige Überschrift,
  leere Tabelle) decken nur den Totalausfall: wird künftig **eine einzelne**
  neue §4-Zeile abweichend formatiert (z. B. ohne Leerzeichen nach dem Pipe,
  mit Fett-Auszeichnung oder einem Zeichen außerhalb des Vorrats), fällt
  genau dieser Code still aus `specCodes` heraus, während die Tabelle
  nicht-leer bleibt. Vergisst der Implementierer zugleich den
  `AllReasons()`-Eintrag — exakt der historische Fehlermodus, den der Slice
  schließt —, bleiben beide Tests grün und `--doctor` zeigt wieder den
  Rohcode. Die Risiko-Notiz des Slice-Plans („ändert sich das
  Tabellen-Format, wird der Test rot") gilt nur für den Ganz-Tabellen-Fall,
  nicht für die Einzel-Zeile.
- **failure-szenario:** Neues Modul trägt seinen Code doc-first als
  `|`neu-code`| modul | …` (ohne Leerzeichen) in §4 ein, `AllReasons()` wird
  vergessen ⇒ `make test` grün, Klartext-Lücke Nr. 8 entsteht unbemerkt.
- **verifizierbar:** ja — Mutations-Probe: eine bestehende §4-Zeile
  entsprechend umformatieren **und** ihr Paar aus
  `AllReasons()`/`reasonTexts()` entfernen ⇒ `make test` bleibt grün
  (Verriegelung greift nicht); Rückformatierung ⇒ rot.

### LOW-2 — `planning-drift`-Klartext unterschlägt die Mehrdeutig-Alternative der §4-Bedingung

- **kategorie:** LOW
- **quelle:** DC-FA-PLAN-001 / Spezifikation §4 (Bedingungsspalte, Zeile 1422)
- **pfad:** `internal/hexagon/core/app/diagnose.go:110`
- **befund:** §4 nennt drei fail-closed-Auslöser: Überschrift **fehlt/ist
  mehrdeutig** bzw. Roadmap-Datei fehlt
  (ebenso §`DC-FA-PLAN-001.a` Schritt 3: „oder kommt sie mehrfach vor").
  Der Klartext sagt nur „(oder Roadmap/Überschrift **fehlt** — fail-closed)".
  In der Prosa-Diagnose ist der Klartext die einzige Grund-Erklärung
  (`report.go:75` rendert `f.Message` nicht) — beim Mehrdeutig-Fall trifft
  keine der genannten Alternativen zu.
- **failure-szenario:** Roadmap trägt `## Aktuelle Welle` zweimal ⇒
  `planning-drift`; `--doctor` behauptet eine fehlende Überschrift, die
  zweifach existiert — der Nutzer sucht am falschen Ort.
- **verifizierbar:** ja — `--doctor --enable planning` gegen eine
  Roadmap-Probe mit dupliziertem Heading zeigt den Klartext-Wortlaut.

### INFO-1 — Undokumentierte Annahme: §4 trägt genau eine Backtick-Tabelle

- **kategorie:** INFO
- **quelle:** Maintainability (dokumentationswürdige Annahme)
- **pfad:** `internal/hexagon/core/app/diagnose_test.go:86`–`:93`
- **befund:** Der Parser erntet **jede** Backtick-Lead-Tabellenzeile zwischen
  der §4-Überschrift und der nächsten `## `-H2. Die Überschrift heißt
  „Grund- **und Fehler**-Codes"; die Fehler-Codes sind heute nur Prosa
  (Zeile 1425 f.). Würde ein künftiger Slice die Exit-2-Fehler-Codes in §4
  tabellieren, landeten sie in `specCodes` und der Test forderte sie in
  `AllReasons()` ein — falsch-rot (fail-closed-Richtung, laut, nicht still),
  aber die Annahme „eine Tabelle je Sektion" steht nirgends.
- **failure-szenario:** Fehler-Code-Tabelle in §4 ⇒ Test rot mit
  irreführender Forderung, Nicht-Befund-Codes ins Klartext-Mapping
  aufzunehmen.
- **verifizierbar:** ja — Probe-Edit einer zweiten Backtick-Tabelle in §4
  macht den Test rot.

---

## Negativbefunde (geprüft, ohne Befund)

- **Spec-Treue der Deckung:** `AllReasons()` deckt exakt die 23 §4-Codes,
  beidseitig, string-identisch über die `Reason*`-Konstanten (keine
  String-Literale — DoD-/Plan-Vorgabe eingehalten).
- **Fachliche Korrektheit der sieben Klartexte:** gegen die
  §4-Bedingungsspalte und §`DC-FA-DIAG-001.a`/`VER`/`PIN`/`IMM`/`VCS`/
  `COMMITS`/`PLAN-001.a` geprüft — inhaltlich korrekt (`core-drift-vcs`
  deckt alle drei Verletzungsarten: Range-Änderung, D/R,
  Status-Übergang); einzige Abweichung ist die Mehrdeutig-Kondensation
  (LOW-2). Stil-Parität zu den 16 Bestands-Klartexten gegeben (deutsch,
  einzeilig, Bedingungs-benennend, ohne Punkt).
- **Verriegelungs-Test, fail-closed-Zweige:** unlesbare Spec ⇒ Fatal;
  Überschrift fehlt ⇒ Fatal; Überschrift mehrfach ⇒ Fatal (Zeilennummern
  im Fehlertext); leere Code-Menge ⇒ Fatal — kein stilles Grün mit leerer
  Menge. Abbruch korrekt bei der nächsten `## `-H2 (`###` stoppt nicht,
  in §4 folgt aber direkt `## 5.` — verifiziert).
- **Pfad/Build-Kontext:** `../../../../spec/spezifikation.md` stimmt in der
  Tiefe (4 Ebenen ab `internal/hexagon/core/app`); Docker-`test`- und
  `coverage`-Stage tragen `spec/` im Kontext (kein `.dockerignore`);
  falscher Pfad fiele fail-closed auf ReadFile-Fatal.
- **Determinismus:** reiner Datei-Read + Mengenvergleich; Map-Iteration
  beeinflusst nur die Reihenfolge der Fehlermeldungen, nicht das
  Bestehen/Scheitern.
- **Stilles-Grün-Suche im Bestand:** der bestehende Deckungs-Test
  (`reasonTexts` ↔ `AllReasons`, beidseitig via Längenvergleich) bleibt
  unverändert; die Kombination beider Tests schließt alle
  Einzel-Vergessens-Pfade (Code nur in §4 / nur in `AllReasons` / Klartext
  fehlt / Klartext verwaist). Einziger Restpfad ist die kombinierte
  Format-plus-Vergessens-Mutation (LOW-1).
- **Hexagon/Hard Rules:** `os`/`path/filepath` nur im Testfile;
  `.a-check.yml` nimmt `**/*_test.go` mit dokumentierter Begründung aus
  (Skript-Parität, Bestandsmuster: git-Adapter-Test importiert `os`) —
  tragfähig, kein R4-Bruch im Produktcode; keine Inline-Suppressions;
  kein Netz-Zugriff (DC-QA-03 unberührt).
- **Kein Verhaltens-Delta außerhalb der Diagnose:** Diff fasst nur
  `AllReasons`/`reasonTexts` + Kommentar an; `FixCandidateFor`,
  Default-/JSON-/YAML-Befundausgabe und Exit-Codes unverändert.
- **Doc-Konsistenz Slice-Plan ↔ Implementierung:** Klartext-Tabelle in §2
  des Plans ist wortidentisch mit `reasonTexts()`; Mutations-Beleg-Plan
  (a)–(c) — inkl. der im feat-Commit nachgeschärften Paar-Entfernung —
  passt zum Testaufbau; „erster Test, der die echte Spec liest" als Risiko
  benannt und begründet.
- **Roadmap-Eintrag:** die sieben Codes vollständig und korrekt benannt;
  „seit v0.25" historisch verifiziert (v0.25.0: 14 Einträge bei
  gleichzeitigem `diagram-id-undefined`-Release); Links/Anker auf
  Lastenheft/Spezifikation korrekt aufgelöst; welle-49-Flip konsistent zur
  planning-Invariante (aktive Welle ∧ Slice in `in-progress/`).
- **SemVer-Einordnung:** Patch v0.37.1 plausibel — Defekt-Fix am
  bestehenden `DC-FA-CLI-007`-Vertrag, keine neue Config-/CLI-Surface,
  kein neues Modul, kein CR/ADR nötig (Codes waren in §4 bereits
  dokumentiert).
- **Commit-Grenzen:** feat-Commit fasst README/Handbuch/CHANGELOG nicht an
  (Release-Prep-Territorium) — konform zur Repo-Arbeitsweise.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 1 |

## Verdikt

**ACCEPT.** Keine HIGH-/MEDIUM-Findings; die Verriegelung schließt den
historischen Fehlermodus real (beidseitig, fail-closed), die sieben
Klartexte sind fachlich korrekt und stil-paritätisch. LOW-1 (Einzelzeilen-
Format-Restlücke) und LOW-2 (`planning-drift`-Wortlaut) sind nice-to-fix —
LOW-2 wäre als Ein-Wort-Ergänzung noch im Slice billig, blockiert aber
nicht; INFO-1 ist eine Won't-Fix-taugliche Designnotiz.
