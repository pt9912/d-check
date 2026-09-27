# Review-Report: slice-237 — 2026-09-27

**Review-Art:** Plan/Code/Doku — geprüft gegen den Slice-Plan
(`slice-237-structure-max-lines-abschnitt.md`: §1 Ziel und Abgrenzung, §2 DoD,
§6 Risiken), `docs/plan/adr/0089-structure-max-lines-zwoelfte-bedingung.md`,
`AGENTS.md` §3 (Hard Rules, insbes. §3.4/§3.6/§3.7/§3.8) und §5, sowie die
Reviewer-Skill-Prüffragen (Modul 10).

**Gegenstand:** slice-237 — CR-Commit `3449cf57` (eingehender Change Request
`docs/plan/cr/2026-09-27-cr-eingehend-ai-harness-course-structure-max-lines.md`),
Vertrags-Commit `2be3fda2` (`spec/lastenheft.md`, `spec/spezifikation.md`,
`docs/plan/adr/0089-structure-max-lines-zwoelfte-bedingung.md`,
`docs/plan/adr/README.md`) und Umsetzungs-Commit `fb478768`
(`internal/hexagon/core/rules/structure.go`, `structure_maxlines_test.go`,
`model/config.go`, `model/finding.go`, `configyaml.go`, `app/diagnose.go`).
Ausgangsstand: `5c8399b4`. Verglichen wurde `git diff 5c8399b4 fb478768` sowie
die drei Einzel-Commits, Arbeitsbaum gleich HEAD (`fb478768`).

**Skill:** `.harness/skills/reviewer.md` @ Version 1.16.0
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-27

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-237-structure-max-lines-abschnitt.md` (§1 Abgrenzung
  „Ausdrücklich NICHT", §2 DoD, §3 Plan/Spiegel-Liste, §6 Risiken)
- `DC-FA-STRUCT-001`/`DC-FA-STRUCT-001.a` (erweiterte Anforderung), `DC-QA-02`
- `ADR-0089` (`Schärft:` DC-FA-STRUCT-001)
- `AGENTS.md` §3.4 (Referenz-Richtung), §3.6 (Gate-Lockerung nur mit ADR),
  §3.7 (Kommentar-Klassen), §3.8 (Modul verspricht nur über Scan-Menge), §5
  (Grenzen-Regel, `MR-025` Spiegel/vor-dem-Editieren, README/CHANGELOG-Regel,
  Zitat-Geltungsbereich, Commit-Botschaft-Overclaim)
- Vorherige Findings am Schwester-Feature (`docs/reviews/2026-09-27-slice-236-file-modul-review-r1.md`,
  Modul `file`, ADR-0088) — gezielt gegen F-1 (README/Handbuch im Feature-Commit
  trotz Plan-Abgrenzung) und F-2/F-3/F-4 (unvollständiger Spiegel:
  Bereichskürzel, Glossar, Betriebsdoku) geprüft
- Unabhängig nachgefahrene Gate-Läufe: `make test`, `make lint`,
  `make doc-check`, `make coverage-gate`, `make arch-check`,
  `make gate-consistency`, `make planning-check`, `make baseline-verify`,
  `make workflow-pins`, `make semgrep`

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Spezifikation verspricht ein eigenes Akzeptanzkriterium „Zeilenbudget (Config-Rand): Given `max-lines` explizit `0` oder negativ, when `d-check` startet, then Exit 2" (`spec/lastenheft.md`) sowie dieselbe Zusage in `spec/spezifikation.md` (`structure[].max-lines`-Schema-Zeile: „explizit < 1 ⇒ Exit 2"). Die Implementierung erfüllt das (`structureBedingungsFehler`, neuer Zweig `r.MaxLines != nil && *r.MaxLines < 1`), aber **kein Test** exerciert diesen Zweig: `grep -rl max-lines internal/` trifft nur `structure.go`/`structure_maxlines_test.go` (Regel-Ebene) und `configyaml.go` selbst — die Config-Decode-Ebene (`configyaml_test.go`) bleibt ohne jeden `max-lines`-Fall, obwohl `TestDecode_StructureFehler` genau diese Config-Rand-Fälle für jede strukturell gleichartige Bedingung führt (z. B. „max-open-tasks negativ", „exempt-expect-count negativ"). `make coverage-gate` bestätigt die Lücke empirisch: `structureBedingungsFehler` bleibt bei 70.6 % (der neue Zweig ist einer der ungedeckten Pfade). | DC-FA-STRUCT-001.a Zeilenbudget (Config-Rand); Reviewer-Anker #13 (fehlender Negativtest zu neuem Vertrag) | `internal/adapter/driven/configyaml/configyaml.go:351-353` (`structureBedingungsFehler`); `internal/adapter/driven/configyaml/configyaml_test.go:250-309` (`TestDecode_StructureFehler`, kein `max-lines`-Fall) | ja — `make coverage-gate` zeigt `structureBedingungsFehler 70.6%`; `grep -rn 'max-lines' internal/adapter/driven/configyaml/*_test.go` liefert keinen Treffer | Fehlender Negativtest zu einer explizit im Vertrag zugesagten Config-Rand-Eigenschaft |

## Negativbefunde (die Prüffragen des Auftrags)

| Prüffrage | Ergebnis |
|---|---|
| Zählregel (bereinigter Text, Zeilenumbrüche, kein `wc -l`-Stil, kein `+1`) | geprüft, ohne Befund. `maxLinesViolation` zählt `strings.Count(body, "\n")` auf `body`, das über `SectionProse` rekonstruiert wird — jede erhaltene Zeile trägt einen Zeilenumbruch dahinter, **auch die letzte** (`sections.go:66-77`, verifiziert durch Lesen), also keine unvollständige Schlusszeile wie bei `file[].max-lines`/`countLines`. `TestStructureMaxLines_Grenzwert` belegt N=5 grün / N+1=6 rot mit `abschnittMitZeilen`, das die genaue Zeilenzahl konstruiert. |
| Grenzwert inklusiv (`<=`) | geprüft, ohne Befund. `if got <= *maxLines { return "", false }` — der Grenzwert selbst löst keinen Befund aus, exakt wie in Lastenheft/Spezifikation zugesagt („Grenzwert ist inklusiv, wie bei `cell-max-chars`"). |
| Untergrenze `1`, nicht `0` — Konsistenz über Config-Validierung/Spec/Test | **F-1**: Config-Validierung (`configyaml.go`) und beide Spec-Straten (`lastenheft.md`, `spezifikation.md`) sind konsistent bei `< 1 ⇒ Exit 2`; der Test fehlt. |
| Grund-Code `section-lines-exceeded` — Konsistenz `finding.go`/`diagnose.go`/`spezifikation.md` §4 | geprüft, ohne Befund. `ReasonSectionLinesExceeded = "section-lines-exceeded"` steht in `AllReasons()` und `reasonTexts()` (`diagnose.go:89,158`) sowie in `spec/spezifikation.md` als `SPEC-087` (Zeile 3414, Modul `structure`). `TestAllReasonsDeckungGegenSpezifikationGrundCodes` parst die §4-Tabelle rein über die Zeilenform `\| \`SPEC-NNN\` \| \`code\` \|` — die im Diff sichtbare Tabellen-Platzierung von `SPEC-087` direkt hinter `SPEC-083` (vor `SPEC-079`) ist zahlenmäßig nicht monoton, folgt aber demselben vorbestehenden Cluster-nach-Modul-Muster wie der Rest der Tabelle (z. B. `SPEC-058`→`SPEC-067`, `SPEC-086`→`SPEC-069` weiter unten) und ist für den Parser irrelevant. `make test` bestätigt grün. |
| Refaktorierung `maxLinesViolation` (gocognit-Auslagerung) | geprüft, ohne Befund. Die Funktion ist eine reine Extraktion ohne Bedeutungsänderung: gleicher `nil`-Check, gleiche Schwellen-Semantik, gleiche Meldungsform (`"Abschnitt trägt N Zeilen (...), erlaubt sind M"`), passend zum bestehenden Muster der Nachbarbedingungen in `structureConditions`. Der Docstring benennt die Auslagerung korrekt als Kopplungs-/Abgrenzungs-Kommentar (§3.7-konform: Zusage + Abgrenzung, kein Slice-Bezug, ein ADR-Zeiger). |
| Abgrenzung §1 „Ausdrücklich NICHT" (keine Datei-Schwelle, keine Umbenennung von `file[].max-lines`, kein `tasks-ignore-pattern`-Pendant) | geprüft, ohne Befund. `.d-check.yml` unberührt (kein `max-lines`-Wert irgendwo aktiviert), `file.go`/`file_test.go` nicht im Diff des Umsetzungs-Commits, kein neuer Teilausnahme-Schlüssel (`tasks-ignore-pattern`-artig) im Schema. |
| Commit-Zerlegung (CR rein / Vertrag rein / Code+Test rein, kein README/Handbuch/CHANGELOG im Feature-Commit) | geprüft, ohne Befund — die im Auftrag benannte Wiederholungsgefahr von slice-236/F-1 tritt hier **nicht** ein: `git show fb478768 --stat` zeigt ausschließlich `internal/`-Pfade (6 Dateien), kein `README*.md`, kein `CHANGELOG.md`, kein `docs/user/`. `3449cf57` trägt ausschließlich die CR-Datei, `2be3fda2` ausschließlich ADR/Spec/Index. Konsistent mit der im Auftrag genannten Erwartung, dass diese Fehlerklasse hier nicht greift (keine README-Passage zu einem neuen Modulnamen nötig, da `structure` nicht neu ist). |
| Abgrenzung `structure[].max-lines` vs. `file[].max-lines` (Schema- und Historie-Zeile) | geprüft, ohne Befund. Beide Schema-Zeilen (`spec/spezifikation.md`) tragen wechselseitig die Abgrenzung (rohe ganze Datei vs. bereinigter Abschnitt), ebenso beide Historie-Einträge (0.89.0/0.90.0) und ADR-0089 §Konsequenzen („Namens-Kollision, keine Vertrags-Kollision"). |
| Spiegel-Vollständigkeit (Bereichskürzel §3, Glossar §6, Betriebsdoku-Modulliste) — gezielt gegen slice-236 F-2/F-3/F-4 geprüft | geprüft, ohne Befund und **zu Recht nicht anwendbar**: `structure`/`STRUCT` ist kein neues Kürzel und kein neuer Modulname — `grep -n 'STRUCT' spec/lastenheft.md` und die Glossar-Zeile (§6 „Regelmodul") führen `structure` bereits vor diesem Slice; `docs/user/operations.md` braucht keinen neuen Options-Eintrag, weil `--enable structure`/`--disable structure` bereits existieren. Die Diagnose des Auftrags („die Bereichskürzel-/Glossar-Punkte gelten hier wahrscheinlich nicht") bestätigt sich empirisch: kein Modul-Namens-Spiegel ist berührt. |
| Test-Abdeckung gegen DoD-Liste (Grenzwert N/N+1, Fenced-Code, section-missing, section-ambiguous, hint, byte-identisch) | geprüft, ohne Befund über F-1 hinaus. Alle sechs Fälle sind in `structure_maxlines_test.go` als sechs eigenständige Tests vorhanden, jede Assertion prüft tatsächlich das behauptete Verhalten (konkrete Zeilenzahl, konkreter Grund-Code, konkrete Meldungs-Teilstrings) — keine Tautologien, keine zu laxen Vergleiche. `abschnittMitZeilen` ist bewusst konstruiert, um eine Zählartefakt-Falle (Phantom-Zeile am Dateiende) zu vermeiden, mit erklärendem Kommentar. |
| ADR-0089-Form (≥ 3 Alternativen, Fitness Function mit real existierenden Tests, Re-Evaluierungs-Trigger, `Schärft:` aufwärts) | geprüft, ohne Befund. Drei verglichene Alternativen neben der gewählten Option (Nichts tun · `forbid-pattern`-Wiederholungsmuster · Zeichen- statt Zeilen-Zählung), je mit Pro/Contra. Fitness-Function-Tabelle nennt `TestStructureMaxLines_*` (sechs real existierende Funktionen, alle im Diff verifiziert) und `TestAllReasonsDeckungGegenSpezifikationGrundCodes` (vorbestehend, real, durch die Erweiterung automatisch mitgeprüft). Re-Evaluierungs-Trigger: zwei disjunkte, je hinreichende Bedingungen plus „ohne eines von beiden: permanent" — korrekte Form. `Schärft:` zeigt auf `DC-FA-STRUCT-001` (Kennung, nicht nur „§N" — Reviewer-Anker #16 erfüllt). |
| Ordinal-Zählung „zwölfte Bedingung" (Anker #17: Messung zählt Muster, das dem Gegenstand nur ähnelt) | geprüft, ohne Befund. Die Kette achte=`headings-match` (0.64.0), neunte=`cell-max-chars` (0.72.0), elfte=`open-tasks-require-marker` (0.87.0, „Erweiterung der elften Bedingung" bestätigt in 0.87.1) ist in `spec/lastenheft.md`s eigener Historie wörtlich auffindbar; `max-open-tasks` (0.79.0, „zehnte") trägt in der Historie tatsächlich **keine** explizite Ordinalzahl — exakt wie der neue 0.90.0-Eintrag selbst einräumt („zehnte: `max-open-tasks` 0.79.0 unbeziffert"). Die Zählung ist nachvollziehbar, keine Übertreibung. |
| Kommentar-Klassen (§3.7) an den vier neuen Code-Kommentaren | geprüft, ohne Befund. Alle vier (`config.go` Feld-Kommentar, `finding.go` Grund-Code-Kommentar, `structure.go` Funktions-Kopf und `structureConditions`-Update) tragen Zusage/Abgrenzung mit **einem** auflösbaren ADR-Verweis (`ADR-0089`), keine Slice-Nummer, keine Review-Historie, keine Deliberation über Verworfenes im Kommentartext selbst (die verworfenen Alternativen stehen in der ADR, nicht im Code-Kommentar). |
| Referenz-Richtung / Provenance-Marker | geprüft, ohne Befund. Beide `slice-237 <!-- d-check:status-provenance -->`-Stellen in ADR-0089 zeigen auf den Entstehungs-/Beleg-Ort (Anlass-Feld, bzw. „trotz des in slice-237 … benannten Ausschlusses"), nicht auf eine Entscheidungsbegründung. `make doc-check` bestätigt 829 Dateien / 0 Befunde inkl. aktivem `matrix`-Modul. |
| Config-Struct-Feld-Konsistenz (`model.StructureRule`, `rawStructure`, `applyStructureRule`) | geprüft, ohne Befund. `MaxLines *int` ist in allen drei Stellen durchgereicht (`config.go:520`, `configyaml.go:235`, `configyaml.go:536`); `gofmt`/`make lint` bestätigen 0 Issues, keine Ausrichtungs-Abweichung. |
| Unabhängig nachgefahrene Gates | `make test`: alle Pakete `ok`. `make lint`: „0 issues." `make doc-check`: „829 Datei(en) geprüft, 0 Befund(e)". `make coverage-gate`: „94.30% erfüllt Schwelle 93%" (siehe F-1 für die Zweig-Ebene darunter). `make arch-check` (a-check `v0.19.0`): „gesamt: 0 Befund(e)". `make gate-consistency`, `make planning-check`, `make workflow-pins`: je „0 Befund(e)"/OK. `make baseline-verify`: alle Regelwerk-/Template-Dateien OK, `fetch-baseline-cache --verify`: „54 Dateien, vollständig". `make semgrep`: „0 findings" über 55 Regeln/65 Dateien. `record-gates` wurde nicht ausgeführt (Verifier-Aufgabe). |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Fehlender Negativtest zu einer zugesagten
Vertragseigenschaft (Config-Rand-Untergrenze einer neuen Zähl-Bedingung)

## Verdikt

**Merge-blockierend:** nein (aus Reviewer-Sicht) — F-1 ist die einzige
Auffälligkeit, klärungsbedürftig vor Closure, aber kein Korrektheitsfehler im
gelieferten Verhalten: Die Config-Validierung selbst ist korrekt implementiert
und deckt sich mit Spec/ADR, nur ihr Test fehlt. Der etablierte Hausstandard
(`TestDecode_StructureFehler` führt exakt diese Art Config-Rand-Test für jede
strukturell gleichartige Bedingung, u. a. `max-open-tasks negativ`) legt nahe,
dass hier ein Fall analog „max-lines 0"/„max-lines negativ" fehlt, bevor der
Slice nach `done/` geht.

**Kern der Implementierung:** Zählregel, Grenzwert-Inklusivität, Fenced-Code-
Ausschluss, Grund-Code-Spiegel, Refaktorierung, Abgrenzung zu `file[].max-lines`
und zur eigenen §1-Abgrenzung sowie die Commit-Zerlegung sind vertragskonform
und durch unabhängig nachgefahrene Gate-Läufe belegt. Die im Auftrag gezielt
geprüften Wiederholungsrisiken aus dem Schwester-Review (slice-236: README/
Handbuch im Feature-Commit, unvollständiger Bereichskürzel-/Glossar-/
Betriebsdoku-Spiegel) treten hier **nicht** auf — erwartungsgemäß, da
`max-lines` eine Erweiterung einer bestehenden Anforderung ist, kein neues
Modul. Übergabe an Verifier/Planner: DoD-Abhaken, Closure-Notiz, Register-
Eintrag und Risiko-Ausgänge (§6 des Slice-Plans) bleiben deren Aufgabe; F-1
sollte vor der Closure entweder durch einen Test geschlossen oder bewusst als
Abweichung von der `TestDecode_StructureFehler`-Konvention begründet werden.
