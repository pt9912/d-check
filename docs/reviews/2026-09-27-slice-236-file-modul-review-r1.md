# Review-Report: slice-236 — 2026-09-27

**Review-Art:** Plan/Code/Doku — geprüft gegen den Slice-Plan
(`slice-236-datei-obergrenze-zeilen-bytes.md`: §1 Ziel und Abgrenzung, §2 DoD,
§6 Risiken), `docs/plan/adr/0088-file-modul-groessengrenzen.md`, `AGENTS.md` §3
(Hard Rules, insbes. §3.4/§3.8) und §5, sowie die Reviewer-Skill-Prüffragen
(Modul 10).

**Gegenstand:** slice-236 — Vertrags-Commit `bdc47b19` (`spec/lastenheft.md`,
`spec/spezifikation.md`, `docs/plan/adr/0088-file-modul-groessengrenzen.md`,
`docs/plan/adr/README.md`) und Umsetzungs-Commit `92ec638a`
(`internal/hexagon/core/rules/file.go`, `file_test.go`, `model/config.go`,
`configyaml.go`, `run.go`, `diagnose.go`, `config_template.go`, `README.md`,
`README.de.md`). Ausgangsstand: `73d07b9a`. Verglichen wurde
`git diff 73d07b9a bdc47b19` bzw. `git diff bdc47b19 92ec638a`, Arbeitsbaum
gleich HEAD (`92ec638a`).

**Skill:** `.harness/skills/reviewer.md` @ Version 1.16.0
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-27

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-236-datei-obergrenze-zeilen-bytes.md` (§1 Abgrenzung
  „Ausdrücklich NICHT", §2 DoD, §3 Plan/Spiegel-Liste, §6 Risiken)
- `DC-FA-FILE-001`/`DC-FA-FILE-001.a`, `DC-FA-CLI-002`, `DC-FA-STRUCT-001`/
  `.a` (Vergleichsmodul), `DC-QA-02`
- `ADR-0088` (`Schärft:` DC-FA-FILE-001, DC-FA-CLI-002)
- `AGENTS.md` §3.4 (Referenz-Richtung), §3.8 (Modul verspricht nur über
  Scan-Menge), §5 (Grenzen-Regel, Spiegel/`MR-025`, README/CHANGELOG-Regel,
  Zitat-Geltungsbereich)
- `.d-check.yml` (`matrix`-Klassen `slice`/`adr`, Provenance-Marker)
- Unabhängig nachgefahrene Gate-Läufe: `make doc-check`, `make test`,
  `make lint`, `make arch-check`, `make coverage-gate`

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Der Umsetzungs-Commit `92ec638a` ändert `README.md` und `README.de.md` (neuer Absatz zum Modul `file`), obwohl der Slice-Plan §1 unter „Ausdrücklich NICHT in diesem Slice" wörtlich ausschließt: „**Handbuch, README, CHANGELOG, Release** — Release-Prep, kein Feature-Commit (`AGENTS.md` §5)". `AGENTS.md` §5 sagt für README/Handbuch dieselbe Release-Prep-Bindung zu, die dort für `CHANGELOG.md` explizit gemacht wird. Der Plan benennt diese Grenze selbst — der Code hält sie nicht ein. | Plan §1 / `AGENTS.md` §5 | `README.md:216-227`, `README.de.md:218-231` (Commit `92ec638a`) | ja — `git show 92ec638a --stat` listet beide Dateien; Plan-Text liegt vor | Feature-Commit fasst README/Handbuch-Klasse trotz expliziter Plan-Abgrenzung an |
| F-2 | MEDIUM | `spec/lastenheft.md` §3 führt in der Schema-Konvention-Blockquote alle Bereichskürzel (`CLI`, `SCAN`, …, `STRUCT`, `WF`, `RVW`, `MENT`, `CONF`, `DIST`) einzeln auf — **`FILE` fehlt**, obwohl `DC-FA-FILE-001` dieses Kürzel führt. Jede vorherige Modul-Einführung hat sich hier selbst eingetragen (z. B. „Bereich `SRC` in §3" für `sources`, „Bereich `STRUCT` in §3" für `structure`); dieser Vertrags-Commit editiert genau diesen Block (er liegt in derselben Datei, wenige Zeilen über der neuen Anforderung) und zieht ihn nicht nach. | MR-025 (Spiegel vor dem Editieren); Modul 5 §Ziel-Form Slice (Spiegel-Pflicht) | `spec/lastenheft.md:57-81` | ja — `grep -n '\`FILE\`' spec/lastenheft.md` liefert keinen Treffer | Spiegel unvollständig — Bereichskürzel-Deklaration nicht nachgezogen |
| F-3 | MEDIUM | `spec/lastenheft.md` §6 Glossar, Zeile „Regelmodul": die Aufzählung bricht bei `structure` ab (20 Module: `links` … `structure`) — `workflows`, `reviews`, `mentions` (aus früheren Slices) **und jetzt `file`** fehlen. Der Vertrags-Commit editiert dieselbe Datei (Lastenheft), zieht diese Zeile aber nicht nach, obwohl frühere Historie-Einträge „… in `DC-FA-CLI-002` + Glossar" als Teil der Modul-Einführung nennen (z. B. bei `sources`, `citations`, `targets`). | MR-025; Modul 5 §Ziel-Form Slice | `spec/lastenheft.md:3869` | ja — Zeile enthält nur 20 der 25 gültigen Modulnamen (`model.ValidModules()` als Referenz) | Spiegel unvollständig — Glossar-Zeile veraltet (vorbestehend, durch diesen Slice nicht behoben) |
| F-4 | MEDIUM | `docs/user/operations.md:33` listet die Options-Tabelle für `--enable`/`--disable` mit den 24 alten Modulnamen, ohne `file`. Der Plan §3 nennt „Betriebsdoku" ausdrücklich als zu prüfenden Spiegel-Kandidaten; §1 nennt als Release-Prep-Ausnahme aber nur „Handbuch, README, CHANGELOG, Release" — „Betriebsdoku"/`operations.md` steht dort **nicht**. Die DoD-Zeile „Spiegel" verlangt: „was Release-Prep ist, steht dort benannt" — für `operations.md` fehlt diese Benennung, das DoD-Häkchen ist trotzdem `[x]`. (Im Plan §6 als offenes Risiko „Ein neues Modul fehlt in einer Aufzählung, die kein Gate hält" mit Ausgang „bei Closure zu vergeben" bereits vermerkt — der Punkt ist damit nicht verschwiegen, aber die konkrete Fundstelle fehlt noch.) | Plan §2 DoD „Spiegel", §3, §6; MR-025 | `docs/user/operations.md:33` | ja | Spiegel unvollständig — Betriebsdoku-Modulliste, Release-Prep-Status nicht benannt |
| F-5 | MEDIUM | §`DC-FA-FILE-001.a` Schritt 3 sagt für den unlesbaren Dateibaum ausdrücklich zu: „ein Befund je Regel, kein Sammel-Befund". `CheckFile` implementiert genau das (Schleife über `rules` im `fileTree`-Fehlerzweig). Es fehlt aber ein Test, der `fileTree`/`walkAllFiles` scheitern lässt — `TestCheckFile_UnlesbareDateiFailClosed` nutzt `readErrFS`, das nur `ReadFile` bricht, nicht `List`. `make coverage-gate` bestätigt die Lücke: `CheckFile` 66.7 %, `fileTree` 80.0 %, `walkAllFiles` 78.6 % — während das strukturell identische Verhalten bei `structure` durch `TestStructureUnlesbarerBaumFailClosed` (mit `listErrFS`, inkl. Mehrfachregel-Assertion „ein Befund je Regel") explizit zu 100 % abgedeckt ist. | DC-FA-FILE-001.a Schritt 3; Reviewer-Anker #13 (fehlender Negativtest zu neuem Vertrag) | `internal/hexagon/core/rules/file.go:34-45` (`CheckFile`); `file_test.go` (Test fehlt); Vergleich `structure_test.go:285-300` | ja — `make coverage-gate` zeigt die genannten Prozentsätze; `grep -rn listErrFS internal/hexagon/core/rules/*_test.go` trifft nur `structure_test.go` | Fehlender Negativtest zu einer explizit zugesagten Vertragseigenschaft (Nullmengen-Härte des Dateibaums) |
| F-6 | LOW | §`DC-FA-FILE-001.a` Schritt 2 wiederholt den globalen Symlink-Ausschluss (§1: „Symlinks werden beim Scan weder verfolgt noch als Dateien gewertet … in keinem Modul") nicht ausdrücklich — anders als `structure`s Schritt 2, der explizit „und keine Symlinks (§1)" sagt. Der Code ist gleichwertig sicher (`walkAllFiles`s `switch` kennt nur `KindDir`/`KindFile`, ein `KindSymlink`-Eintrag aus `List()`/Lstat wird weder betreten noch aufgenommen — kein Folgen, kein Zyklus-Risiko, identisch zu `walkMarkdown`), die Doku-Stelle ist aber weniger explizit als beim Nachbarmodul. | DC-FA-FILE-001.a Schritt 2 vs. DC-FA-STRUCT-001.a Schritt 2 | `spec/spezifikation.md:2949-2952` (Schritt 2 `file`) vs. `spec/spezifikation.md:2233-2234` (Schritt 2 `structure`) | ja (Textvergleich) | Weniger explizite Parität mit Nachbarmodul (inhaltlich durch §1 gedeckt) |
| F-7 | LOW | Die DoD verspricht „Grenzwert N und N+1 je Schlüssel". Für `max-lines` ist das ein einziger, symmetrischer Test (`TestCheckFile_MaxLinesGrenzwert`: N=5 grün, N+1=6 rot). Für `max-bytes` ist die Symmetrie über mehrere Tests mit unterschiedlichen Werten verstreut (grüne Seite nur implizit über `TestCheckFile_LeereDatei`, N=0; rote Seite über `MaxBytesZaehltFencedCodeMit`/`BeideSchwellenGleichzeitig`/`BytesZaehltRohNichtRunen` mit jeweils anderem N) — kein Test zeigt am selben, nichttrivialen Byte-Wert beide Seiten der Grenze. Funktional identisch (beide Schwellen nutzen denselben `>`-Vergleich), aber die DoD-Formulierung ist für `max-bytes` nur implizit erfüllt. | Plan §2 DoD „Umsetzung" | `internal/hexagon/core/rules/file_test.go:150-175` | ja (Testlogik gelesen) | DoD-Testversprechen nur implizit erfüllt |

## Negativbefunde (die zehn Prüffragen des Auftrags)

| Prüffrage | Ergebnis |
|---|---|
| 1. Vertrag vs. Code (Zählregel, Grund-Codes, Schema-Schlüssel, Nullmengen-Härte, fail-closed einer unlesbaren Einzeldatei) | geprüft, ohne Befund über F-5 hinaus. `file.go` zählt Zeilen exakt wie in `spec/spezifikation.md` Schritt 5 zugesagt (`countLines`, geteilt mit `codepaths`/`citations`), Bytes exakt wie in Schritt 6 zugesagt (`len(content)`, keine Bereinigung — durch `TestCheckFile_MaxBytesZaehltFencedCodeMit` belegt). Die drei Grund-Codes (`file-no-match`, `file-lines-exceeded`, `file-bytes-exceeded`) sind in `AllReasons()`/`reasonTexts()` **und** `spec/spezifikation.md` §4 (`SPEC-084`–`086`) deckungsgleich benannt. Eine unlesbare Einzeldatei ist fail-closed (`checkFileOne`, Zeile 133-138) — Code und Zusage stimmen; nur der Test für den unlesbaren **Baum** fehlt (F-5). |
| 2. Referenz-Richtung / Decken-Regel | geprüft, ohne Befund. Die neue Historie-Zeile in `spec/lastenheft.md` §7 (0.89.0) trägt „Begründung in begleitender ADR" **ohne** Nummer, wortgleich mit den Nachbar-Zeilen (0.88.0, 0.87.1, 0.87.0); ebenso die neue Zeile in `spec/spezifikation.md` §7 (2026-09-27). Kein Slice-/Wellen-Verweis in beiden Straten. Die zwei `slice-236 <!-- d-check:status-provenance -->`-Erwähnungen liegen ausschließlich in `docs/plan/adr/0088-file-modul-groessengrenzen.md` (Klasse `adr`, für die der Marker per `.d-check.yml` explizit als Ausnahme deklariert ist) und zeigen jeweils auf den Entstehungs-/Beleg-Ort (Anlass bzw. „gemessen … Plan-Anlass, slice-236 §1"), nicht auf eine Entscheidungsbegründung — die eigentliche Begründung folgt im Text danach. `make doc-check` bestätigt die Auflösbarkeit (819 Dateien, 0 Befunde, inkl. aktivem `matrix`-Modul). |
| 3. ADR-Form | geprüft, ohne Befund. ADR-0088 trägt drei verglichene Alternativen inkl. „Nichts tun" mit Pro/Contra, eine Fitness-Function-Tabelle (drei `go test`-Ziele, alle real vorhanden — siehe Prüffrage 7), einen zweiteiligen Re-Evaluierungs-Trigger mit „Ohne eines von beiden: permanent" und ein `Schärft:`-Feld, das korrekt aufwärts auf `DC-FA-FILE-001`/`DC-FA-CLI-002` zeigt (kein Abwärtsverweis). Der ADR-Index-Eintrag (`docs/plan/adr/README.md:98`) folgt der etablierten Form (nur `DC-FA-FILE-001` in der Bezug-Spalte, nicht zusätzlich `DC-FA-CLI-002` — deckt sich mit dem Präzedenzfall ADR-0084/`DC-FA-MENT-001`, das dieselbe Ellipse fährt). |
| 4. §3.8-Frage (Scan-Menge) | geprüft, ohne Befund über F-6 hinaus. §`DC-FA-FILE-001.a` Schritt 2 benennt die Grenze korrekt und vollständig für die zentrale Eigenschaft: „unabhängig von `scan.roots`/`scan.ignore`", „die feste Skip-Liste gilt". `isSkipDir` ist identisch mit `structure`s Funktion (`scan.go:13`, von beiden Modulen aufgerufen — keine Divergenz). Symlink-Behandlung ist über den gemeinsamen `driven.Filesystem`-Port (`Lstat`-basiertes `List()`, `KindSymlink` als dritter, im `switch` nicht behandelter Fall) identisch zu `structure` und damit strukturell zyklensicher (kein Folgen von Verzeichnis-Symlinks, kein `KindDir`-Ergebnis für einen Symlink). MemFS in `coretest` kann Symlink-Zyklen tatsächlich nicht abbilden — das ist eine ungeprüfte, aber durch die Portkonstruktion (nur `Lstat`, nie `Stat`) strukturell geschlossene Lücke, nicht nur an dieser Stelle „ungetestet, aber vielleicht falsch". Die fehlende explizite Wiederholung des §1-Symlink-Satzes im Algorithmus-Schritt ist als F-6 (LOW) geführt. |
| 5. Reuse von `countLines` | geprüft, ohne Befund — Prämisse der Frage korrigiert (siehe unten). `string(content)` in `countLines` (`codepaths.go:200`) ist **keine** Rune-Konvertierung: eine Go-Konvertierung `string([]byte)` kopiert die Bytes unverändert (keine UTF-8-Decodierung/-Validierung — das passiert erst bei `[]rune(s)` oder `range s`). `strings.Count(s, "\n")`/`strings.HasSuffix(s, "\n")` suchen das einzelne ASCII-Byte `0x0A` byte-genau, identisch zu einer `bytes.Count`/`bytes.HasSuffix`-Variante auf dem rohen `[]byte`. Die Wiederverwendung hat daher **keine** Auswirkung auf `file`s unabhängige Byte-Zählung (`len(content)`, `file.go:147`, arbeitet auf dem unveränderten `[]byte`) und **keine** Auswirkung auf ungültiges UTF-8 (Bytes bleiben in der `string`-Kopie exakt erhalten; `TestCheckFile_BytesZaehltRohNichtRunen` belegt zusätzlich empirisch, dass Byte- ≠ Runen-Zahl korrekt behandelt wird). Der einzige Preis ist eine zusätzliche Speicher-Kopie (Go-Strings sind immutable) — eine Performance-, keine Korrektheitsfrage, und ohne Bezug zu `DC-QA-01`. |
| 6. Konfigurations-Validierung (`applyFileRule`) | geprüft, ohne Befund. Alle sechs in Schritt 1 der Spezifikation genannten Exit-2-Fälle sind abgedeckt: fehlendes `files` (Zeile „files ist Pflicht"), ungültiges Glob in `files` **und** in `exempt-paths` (je eigene `path.Match`-Probe), beide Schwellen abwesend, negative `max-lines`/`max-bytes` (je eigene Prüfung), leerer **oder** mehrzeiliger `hint` (`TrimSpace`-Leere bzw. `ContainsAny("\t\r\n")`), doppelte Identität (`applyFile`, `seen`-Map über `Identity()`). Keine der sechs Prüfungen fehlt. |
| 7. Test-Abdeckung gegen die DoD-Liste | geprüft, ohne Befund über F-5/F-7 hinaus. 14 Tests in `file_test.go`, jeder einzeln gedanklich durchgeführt: Grenzwert (F-7-Einschränkung bei Bytes), unvollständige Schlusszeile, Bytes roh trotz Fenced-Code, Nullmengen-Härte, `exempt-paths` bis auf null, unlesbare Einzeldatei fail-closed (nicht: unlesbarer Baum, F-5), `hint` gewinnt/gilt nicht, beide Schwellen gleichzeitig, ohne Regeln byte-identisch, Nicht-Markdown-Kandidat, leere Datei, Bytes ≠ Runen (Umlaute, empirisch mit Assertion auf die Test-Voraussetzung selbst), Regel-Identität im `target`. Alle Assertions prüfen tatsächlich das behauptete Verhalten (keine Tautologien, keine zu laxen Vergleiche). Die drei ADR-Fitness-Function-Tests (`TestAllReasonsDeckungGegenSpezifikationGrundCodes`, `TestPrintConfigVerfuegbarDecktRegistry`) existieren real und sind **nicht** neu in diesem Slice — beide sind generische, vorbestehende Register-Deckungstests, die durch den erweiterten `validModules()`/`AllReasons()`-Bestand automatisch auch `file` mitprüfen; unabhängig nachgefahren (`make test`, `make coverage-gate`: 94.20 % — exakt wie in der Commit-Botschaft behauptet, kein Overclaim). |
| 8. Spiegel-Vollständigkeit | siehe F-2/F-3/F-4 für die gefundenen Lücken. Sonst geprüft, ohne Befund: `internal/hexagon/core/model/config.go` (`validModules()`), `internal/adapter/driven/configyaml/configyaml.go`, `internal/adapter/driving/cli/config_template.go` („Verfügbar"-Zeile, durch `TestPrintConfigVerfuegbarDecktRegistry` gewächtert), `internal/hexagon/core/app/diagnose.go` (`AllReasons()`/`reasonTexts()`, durch `TestAllReasonsDeckungGegenSpezifikationGrundCodes` gewächtert), `.d-check.yml`/`FOCUS_DISABLE`-Paar (bewusst unberührt, ADR-Konsequenz korrekt), `Makefile` (kein `doc-file`-Target — konsistent mit dem Präzedenzfall `mentions`/`reviews`, die ebenfalls keines haben), `internal/adapter/driving/cli/cli_acceptance_test.go` (keine `--enable file`-Akzeptanztests — ebenfalls konsistent mit dem Präzedenzfall `mentions`/`structure`/`reviews`, die bei ihrer Einführung auch keine CLI-Blackbox-Tests bekamen, nur Unit-Tests im jeweiligen `rules`-Paket). `packaging/dockerhub/overview.md` zählt „21 rule modules" und ist bereits vor diesem Slice um vier Module veraltet (`workflows`/`reviews`/`mentions` fehlen ebenfalls) — dieser Rückstand ist nicht durch slice-236 verursacht und liegt außerhalb dessen Abgrenzung (Docker-Hub-Beschreibung, eigener Release-Prep-Pfad laut `spec/lastenheft.md` `DC-FA-DIST-001`-Out-of-Scope), daher hier nur benannt, nicht als Finding gezählt. |
| 9. Commit-Zerlegung | geprüft — für den Vertrags-Commit ohne Befund, für den Umsetzungs-Commit siehe F-1. `git show bdc47b19 --stat` zeigt ausschließlich `docs/plan/adr/0088-…`, `docs/plan/adr/README.md`, den Slice-Plan (DoD-Häkchen) und die beiden Spec-Straten — kein `internal/`-Pfad, kein `README*.md`, kein `Makefile`. `bdc47b19` ist damit eine reine Vertrags-Änderung. `92ec638a` ist für seinen Code-Teil in sich konsistent (kein halb fertiger Vertrag, keine offene Konfigurations-Fläche), trägt aber zusätzlich die README-Änderung aus F-1, die nicht in diesen Commit gehört (weder inhaltlich zum „Vertrag" noch zur erlaubten „Umsetzung"). |
| 10. Abgrenzung (§1 „Ausdrücklich NICHT") | F-1 ist ein Verstoß (README/Handbuch-Klasse). Die übrigen vier Ausschlüsse sind eingehalten: keine Schwelle für `AGENTS.md` oder eine andere konkrete Datei gesetzt (`.d-check.yml` unberührt, kein `max-lines`/`max-bytes`-Wert irgendwo aktiviert), keine Zeichen-/Wort-/Token-Zählung (nur `\n`-Zeilen und rohe Bytes), keine Zählung des bereinigten Textes (`TestCheckFile_MaxBytesZaehltFencedCodeMit` belegt das Gegenteil positiv), kein `--repair`-Kandidat (Modul ist rein diagnose-only, kein Repair-Pfad im Diff), keine Erweiterung von `structure` (eigenes Modul, eigene Grund-Codes, eigener Typ `FileRule`). |
| Unabhängig nachgefahrene Gates | `make doc-check`: „819 Datei(en) geprüft, 0 Befund(e)". `make test`: alle Pakete `ok`, keine Fehlschläge. `make lint`: „0 issues." `make arch-check` (a-check `v0.19.0`): „gesamt: 0 Befund(e)". `make coverage-gate`: „94.20% erfüllt Schwelle 93%" — exakt der in der Commit-Botschaft genannte Wert, kein Overclaim gegenüber der gemessenen Menge (Anker „Botschaft verallgemeinert nicht über die Messung hinaus"). `baseline-verify`, `workflow-pins`, `semgrep`, `gate-consistency`, `planning-check` wurden in diesem Review nicht erneut ausgeführt (kein Bezug zum Diff-Gegenstand dieses Reviews). |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 4 |
| LOW | 2 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Feature-Commit fasst README/Handbuch-Klasse
trotz expliziter Plan-Abgrenzung an · Spiegel unvollständig (Bereichskürzel-
Deklaration) · Spiegel unvollständig (Glossar-Zeile) · Spiegel unvollständig
(Betriebsdoku-Modulliste, Release-Prep-Status nicht benannt) · fehlender
Negativtest zu einer zugesagten Vertragseigenschaft (unlesbarer Dateibaum,
kein Sammel-Befund) · weniger explizite Symlink-Parität mit Nachbarmodul ·
DoD-Testversprechen (Grenzwert N/N+1) nur implizit erfüllt.

## Verdikt

**Merge-blockierend:** ja, wegen F-1 (HIGH) — ein klarer, vom Plan selbst
benannter Abgrenzungs-Verstoß, unabhängig von seiner inhaltlichen Harmlosigkeit
(der README-Text selbst ist korrekt). Die vier MEDIUM-Findings (F-2 bis F-5)
sind vor der Closure zu klären, nicht notwendig vor jedem einzelnen Commit;
F-3 und der `packaging/dockerhub/overview.md`-Rückstand sind vorbestehend und
nicht durch diesen Slice verursacht, gehören aber, wenn dieser Slice sie
aufdeckt, in seine Closure-Notiz oder ins Beobachtungs-Register.

**Kern der Implementierung:** Die Fähigkeit selbst (Zählregel, Grund-Codes,
Konfigurations-Validierung, Nullmengen-Härte, `structure`-Parität, hexagonale
Schichtung, Determinismus) ist vertragskonform und durch unabhängig
nachgefahrene Gate-Läufe belegt. Die Findings betreffen ausschließlich
Prozess-/Spiegel-Disziplin (F-1 bis F-4) und Testabdeckung eines
Randverhaltens (F-5, F-7) — kein Korrektheitsfehler im gelieferten Verhalten
selbst wurde gefunden. Übergabe an Verifier/Planner: DoD-Abhaken (insbes. das
bereits gesetzte `[x]` bei „Spiegel" gegen F-2/F-3/F-4 prüfen), Closure-Notiz,
Register-Eintrag und Risiko-Ausgänge bleiben deren Aufgabe.
