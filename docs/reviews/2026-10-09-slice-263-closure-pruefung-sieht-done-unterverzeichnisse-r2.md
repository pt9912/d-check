# Review R2 — slice-263: Closure-Prüfung über Unterverzeichnisse und Inhalts-Ausnahme (Fix-Runde)

- **Review-Art:** Code. Geprüft wird der Fix-Commit gegen den Slice-Plan (`slice-263`, §3 samt
  Plan-Änderung nach R1, §6), die R1-Befunde, `MR-025` und die Hard Rules `AGENTS.md` §3.4/§3.7/§3.8
  sowie §5 Regel 3/13/15. Die DoD-Abhakung gehört nicht zu diesem Review.
- **Gegenstand:** `slice-263` · Fix-Commit `5c3febc1` (Kontext: Plan-Änderung `f4606192`, R1-Report
  `31d6abe2`, Produkt-Commit `487ee4de`).
- **Skill:** `reviewer.md` @ 1.19.0 (`980d6314`)
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** R1-Report dieses Slice (F-1 bis F-10); Slice-Plan `slice-263` §3 und §6;
  `spec/lastenheft.md` 0.99.0 — `DC-FA-PLAN-001` (Closure-Note-Struktur, fail-closed, Kriterien) und
  `DC-FA-STRUCT-001` (Kandidaten-Menge, fail-closed-Liste, Kriterien), Historie; `spec/spezifikation.md`
  §`DC-FA-PLAN-001.a` C1/C2, §`DC-FA-STRUCT-001.a` Schritte 1/2, §2-Schema, §4-Zeile `SPEC-039`,
  Historie. Im Code gelesen: `internal/hexagon/core/rules/planning.go` (`CheckPlanningClosure`,
  `closureCandidates`, `closureSkip`), `structure.go` (`structureSkipped`),
  `internal/adapter/driven/configyaml/configyaml.go` (`skip-pattern`-Validierung beider Module),
  `config_template.go`, die Tests `planning_closure_rekursiv_test.go`, `structure_skip_test.go`,
  `configyaml_test.go`.
- **Proben** (Wegwerf-Kopie per `git archive HEAD` im Scratchpad, `make test IMAGE=dcheck-r2-probe`;
  Arbeitsbaum unverändert, Kopie und Probe-Image danach entfernt):
  1. Mutation `cur := path.Join(dir, sub)` (Stand vor dem Fix): **rot** —
     `TestClosureDir_UngereinigtMeldungUnveraendert` mit `Message:Closure-Verzeichnis docs/fehlt fehlt …`
     gegen erwartetes `docs/fehlt/`.
  2. Mutation `if false && isSkipDir(e.Name)`: **rot** — `TestClosureRecursive_SkipDirsBleibenUnbetreten`.
  3. Probe-Test: `recursive: true`, `skip-pattern` = Stub-Marker, `done/slice-001-a.md` (volle Notiz),
     `done/2024/slice-002-b.md` (nur Stub), `done/leer/notiz.txt` — **kein Befund**: ein
     Unterverzeichnis ohne Kandidaten meldet nichts.

## Stand der R1-Befunde

| R1 | Stand | Beleg |
|---|---|---|
| F-1 | geschlossen | oberstes `dir` geht roh an `List` und in die Meldung; Test hält es, Probe 1 rot |
| F-2 | **teilweise** — siehe F-1 dieses Reports | `DC-FA-PLAN-001` vollständig nachgezogen; in `DC-FA-STRUCT-001` fehlt `skip-pattern` in der Exit-2-Liste |
| F-3 | geschlossen | beide Vorlagen-Zeilen tragen `(?m)`; der Marker steht in den Stubs am Zeilenanfang (`grep -c '^> \*\*ARCHIVIERT'` trifft die Stubs unter `done/`) |
| F-4 | geschlossen | Grenze in C2 und in `structure` Schritt 2, Risiko in §6 des Plans |
| F-5 | geschlossen | `TestClosureRecursive_SkipDirsBleibenUnbetreten`, Probe 2 rot |
| F-6 | übergeben | geht laut Botschaft in die Closure-Notiz; zu prüfen bei Closure |
| F-7 | geschlossen | Grenze in C2 und im Kommentar von `closureCandidates`, am Code (`KindDir`-Abfrage) stimmig |
| F-8 | geschlossen in der Sache, neuer Fehler im Wortlaut — siehe F-2 dieses Reports | `SPEC-039` und §2-Zeile `structure[].files` nachgezogen |
| F-9 | geschlossen an der benannten Stelle; Klasse steht in den Tests weiter — siehe F-3 | `planning.go` ohne „bisher" |

## Findings

### F-1 — MEDIUM — Lastenheft `DC-FA-STRUCT-001`: Exit-2-Liste ohne `skip-pattern`

- **quelle:** `AGENTS.md` §5 Regel 3 (Anforderungen nur im Lastenheft) und Regel 15 (Botschaft behauptet
  nicht mehr, als die Arbeit trägt); R1 F-2; Skill-Frage 8
- **pfad:** `spec/lastenheft.md` (`DC-FA-STRUCT-001`) · „kompilierendes `section-pattern`/`forbid-pattern`/`require-pattern`/`tasks-ignore-pattern`/`exempt-section-pattern`/`open-tasks-require-marker-section`";
  Commit `5c3febc1` · „nennen recursive und skip-pattern in der Kandidatenmenge, am Config-Rand und im
  fail-closed-Teil"
- **befund:** Die fail-closed-Aufzählung der Rang-1-Anforderung führt jedes RE2-Feld der Regel einzeln
  auf, `skip-pattern` nicht; Spezifikation (`structure` Schritt 1, §2-Zeile `structure[].skip-pattern`)
  und Code (`configyaml_test.go` · „`"skip-pattern RE2"`" im `structure`-Fall) brechen mit Exit 2 ab.
  Die Spezifikation erweitert damit wieder die abschließende Liste des Lastenhefts — dieselbe Klasse,
  die R1 F-2 schließen sollte; für `DC-FA-PLAN-001` ist sie geschlossen („sowie ein nicht kompilierendes
  `skip-pattern`"), für `structure` nicht, obwohl die Botschaft „am Config-Rand" für beide Anforderungen
  behauptet. Wer gegen Rang 1 abnimmt oder nachbaut, findet für ein kaputtes Muster keine Zusage und
  darf es still als „aus" lesen — das Ventil schaltete dann nichts ab, aber auch nichts zu, ohne Meldung.
- **verifizierbar:** nein — `awk 'NR>=3040 && NR<=3060' spec/lastenheft.md | grep -c skip-pattern` ergibt 0;
  das Urteil ist keines, das ein Gate fällt.
- **klasse:** `vertrag-nur-im-niedrigeren-stratum`

### F-2 — MEDIUM — `SPEC-039` sagt einem Unterverzeichnis ohne Kandidaten einen Befund zu, den der Code nicht liefert

- **quelle:** `AGENTS.md` §5 Regel 13 (Grenze gegen den Gegenstand prüfen); Skill-Frage 10
- **pfad:** `spec/spezifikation.md` §4 · „oder das gesetzte `planning.closure.dir` (unter `recursive` auch
  ein Unterverzeichnis) fehlt, ist unlesbar oder enthält **keinen** Kandidaten unter dem effektiven
  Filter (fail-closed)"; Gegenstand `internal/hexagon/core/rules/planning.go` · „`if len(names) == 0 {`"
- **befund:** Die Klammer hängt an allen drei Prädikaten; gelesen wie geschrieben, löst ein
  Unterverzeichnis **ohne Kandidaten** `closure-note-missing` aus. Der Code prüft die Nullmenge nur über
  die Gesamtmenge unter `closure.dir`; ein leeres oder nur aus Stubs bestehendes Unterverzeichnis trägt
  nichts bei und meldet nichts (Probe 3). C2 und die §2-Zeile `planning.closure.recursive` sagen es
  richtig (nur „unlesbar" gilt je Unterverzeichnis), die Grund-Code-Tabelle widerspricht ihnen im selben
  Stratum. Failure-Szenario: Ein Konsument liest die Tabelle, verlässt sich darauf, dass ein Ordner, aus
  dem der Bestand weggewandert ist, gemeldet wird, und liest das Grün als Deckung — es ist keine.
  Dieselbe Mehrdeutigkeit schwächer im Lastenheft: „(unter `recursive` auch ein unlesbares
  Unterverzeichnis) — **und ebenso ein Verzeichnis ohne einen einzigen Kandidaten**" lässt offen, ob
  „ein Verzeichnis" das Unterverzeichnis einschließt.
- **verifizierbar:** ja — Probe 3 als Test (kein Befund) gegen den Tabellen-Wortlaut.
- **klasse:** `grenze-nicht-am-gegenstand-geprueft`

### F-3 — LOW — Chronik-Wörter in den Test-Kommentaren derselben Fähigkeit

- **quelle:** `AGENTS.md` §3.7 („Ein Kommentar beschreibt, was da ist"); R1 F-9
- **pfad:** `internal/hexagon/core/rules/planning_closure_rekursiv_test.go` · „Kandidat wie bisher
  (byte-identisch)"; dieselbe Datei · „unsichtbar — der Vorzustand, und"; Fehlertexte „VORZUSTAND:" und
  „unlesbar wie bisher"
- **befund:** R1 F-9 benannte die Stelle in `planning.go`; der Fix hat genau sie umformuliert. Dieselbe
  Klasse steht in den Kommentaren und Meldungen der Tests aus `487ee4de` weiter — „bisher" und
  „Vorzustand" haben für einen späteren Leser keinen Bezugspunkt. Die Kommentare tragen eine Zusage,
  daher LOW, nicht HIGH.
- **verifizierbar:** nein.
- **klasse:** `kommentar-mit-chronik-wort`

### F-4 — INFO — Satzbau der umgeschriebenen Lastenheft-Absätze

- **quelle:** Maintainability
- **pfad:** `spec/lastenheft.md` (`DC-FA-PLAN-001`) · „dessen Default `planning.slice-glob` ist —; mit";
  „auch in seinen Unterverzeichnissen"; „Verlangt wird in ihr"; „mit `recursive` samt
  seinen Unterverzeichnissen (per Konvention das"
- **befund:** Vier Lesestörungen, keine mit falscher Aussage: (a) Gedankenstrich-Semikolon „—;";
  (b) „seinen" folgt unmittelbar auf „Kandidaten-Filter"/„Default", gemeint ist `closure.dir`;
  (c) „in ihr" greift über den eingeschobenen Plural „die Dateien, deren Inhalt …" hinweg auf „jede
  Datei" zurück; (d) die Apposition „(per Konvention das Verzeichnis der abgeschlossenen Slices)" hängt
  jetzt an „Unterverzeichnissen" statt an `planning.closure.dir`. Ebenso `SPEC-039`: Gedankenstrich-
  Einschub, Klammer und „— oder" in einer Zelle (siehe F-2). In der Spezifikation C2 steht weiter „Ohne
  den Schlüssel bleibt alles wie zuvor" — Bestand aus `487ee4de`, kein Kommentar, aber derselbe
  fehlende Bezugspunkt wie F-3.
- **verifizierbar:** nein.
- **klasse:** `vertragssatz-schwer-lesbar`

### F-5 — INFO — Kriterien und Historie decken nicht jeden neuen Satz

- **quelle:** Maintainability
- **pfad:** `spec/lastenheft.md` · „**fail-closed (Config-Rand):** Given ein `planning.closure.dir`
  außerhalb der Repo-Wurzel"; `spec/spezifikation.md` Historie-Zeile 2026-10-09 · „Ohne die Schlüssel
  byte-identisch; kein neuer Grund-Code"
- **befund:** Zwei neue Zusagen des Lastenhefts haben kein Akzeptanzkriterium, wohl aber einen Test:
  Exit 2 bei nicht kompilierendem `skip-pattern` (beide Module, `configyaml_test.go`) und
  `closure-note-missing` bei unlesbarem Unterverzeichnis
  (`TestClosureRecursive_UnlesbaresUnterverzeichnisFailClosed`). Die drei neuen Kriterien sind durch
  Tests gedeckt (`TestClosureRecursive_SiehtUnterverzeichnis`; `TestClosureSkipPattern_*`;
  `TestStructureSkipPattern_*`). Die Spezifikations-Historie nennt weder die zwei Grenzen noch die
  geänderte `SPEC-039`-Zeile; die bestehende Zeile desselben Tages trägt sie nicht.
- **verifizierbar:** nein.
- **klasse:** `zusage-ohne-kriterium`

## Negativbefunde

- **F-1-Fix am Code:** nur das oberste Verzeichnis bleibt roh, Unterverzeichnisse werden per
  `path.Join` gebildet; Kandidaten bleiben `dir`-relativ, `closureSkip` und `checkClosureNote` lesen
  über `path.Join(dir, name)` wie vor dem Slice. Nullmengen-Meldung nennt `dir` roh wie zuvor. Ohne Befund.
- **Lastenheft gegen Code (`DC-FA-PLAN-001`):** Kandidatenmenge (Basisname-Filter, Abstieg unter
  `recursive`, Abzug nach Inhalt), Reihenfolge Abzug → Nullmenge, „unlesbare Datei bleibt Kandidatin",
  unlesbares Unterverzeichnis ⇒ `closure-note-missing`, Exit 2 für `skip-pattern`, `recursive` ohne
  ungültigen Wert — am Code nachgelesen, stimmen. Abweichung nur F-2 (Lesart „Verzeichnis ohne
  Kandidaten").
- **Lastenheft gegen Code (`DC-FA-STRUCT-001`):** Reihenfolge `exempt-paths` vor `skip-pattern`
  (`structureExempt(r, f) || structureSkipped(…)`), unlesbare Datei bleibt Kandidatin
  (`structureSkipped` gibt bei Lesefehler `false`), Nullmengen-Regel nach beiden Abzügen — stimmen.
  Abweichung nur F-1.
- **Lastenheft gegen Spezifikation:** `SKIP_DIRS` und der nicht verfolgte Symlink stehen nur in der
  Spezifikation; das Lastenheft sagt „auch in seinen Unterverzeichnissen". Das ist dieselbe Schichtung,
  die `DC-FA-STRUCT-001` seit jeher führt („über den gesamten Baum" auf Rang 1, `SKIP_DIRS` und
  „keine Symlinks" auf Rang 2) — Schärfung, kein Widerspruch. Ohne Befund.
- **Akzeptanzkriterien gegen Tests:** „Happy Path (Unterverzeichnisse)" beide Hälften in einer
  Testfunktion; „Boundary (Stub ausgenommen)" Stub/Volltext, Nullmenge, unlesbare Datei je ein Test;
  „Boundary (Inhalts-Ausnahme)" vier Tests einschließlich der Meldung ohne Schlüssel. Ohne Befund außer F-5.
- **Neue Tests aus dem richtigen Grund rot:** Probe 1 und 2. Ohne Befund.
- **Vorlage (`--print-config`):** `(?m)` in beiden Blöcken, YAML-Einfachquote lässt die Backslashes
  stehen; die Vorlage bleibt Kommentar, kein Parse-Pfad betroffen. Ohne Befund.
- **Spiegel (`MR-025`, Schritt 17):** `grep` nach `closure.dir`/`closure.glob`/`skip-pattern` außerhalb
  `done/`, Reviews, ADRs und Handbuch: `README.md`/`README.de.md` (Modul-Überblick, Release-Prep),
  `CHANGELOG.md` (Release-Prep), `.d-check.closure.yml` (Repo-Konfiguration, laut Plan `slice-264`),
  `harness/sensors/verify-closure-notes.md` (nennt nur das fehlende/unlesbare `closure.dir`, bleibt
  wahr). Keine weitere Stelle beschreibt die Kandidatenmengen. Ohne Befund.
- **Referenz-Richtung (§3.4):** in den hinzugefügten Zeilen beider Spec-Straten kein `slice-`,
  `welle-`, `ADR-`-Token und kein Commit-Hash; die Lastenheft-Historie nennt den Anlass ohne Planungs-
  Kennung. Ohne Befund.
- **Versions-Bump und Historie (`MR-032`):** 0.98.0 → 0.99.0 bei Status Draft, Zeile oben in absteigender
  Ordnung, nennt die drei neuen Kriterien wörtlich. Ohne Befund.
- **Kommentare (§3.7) im Fix:** neuer Kommentar in `closureCandidates` (Zusage), GRENZE-Satz (Grenze),
  Test-Kommentare der zwei neuen Tests (Zusage) — ohne Befund-Marker, ohne Slice-Nummer. Ohne Befund
  außer F-3 (Bestand).
- **Hexagon, Netz, Suppression:** Fix fügt keinen Import hinzu, kein `//nolint`, keine Schwelle. Ohne Befund.
- **Abgrenzung des Plans:** Lastenheft-Zeile steht in §3 nach der Plan-Änderung `f4606192` vor dem Code;
  keine Repo-Konfiguration geändert. Ohne Befund.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 2 | F-1, F-2 |
| LOW | 1 | F-3 |
| INFO | 2 | F-4, F-5 |

Wiederkehrende Finding-Klasse für die Closure-Notiz: `vertrag-nur-im-niedrigeren-stratum` (R1 F-2,
hier F-1 — die Korrektur selbst wurde nur für eine der zwei Anforderungen vollständig gezogen) und
`grenze-nicht-am-gegenstand-geprueft` (Regel-13-Klasse, `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`).

## Verdikt

**Nicht freigegeben, Nachzug klein.** Sieben der acht adressierten R1-Befunde sind in der Sache
geschlossen; F-2 aus R1 ist es für `DC-FA-STRUCT-001` nicht (F-1), und die nachgezogene
`SPEC-039`-Zeile sagt mehr zu als der Code (F-2). Beides ist Wortlaut in den Spec-Straten, kein Code;
der Code des Fixes ist ohne Befund.
