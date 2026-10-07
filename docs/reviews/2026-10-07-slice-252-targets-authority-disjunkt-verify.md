# Verifikation — slice-252: `targets` prüft opt-in die Disjunktheit der Autoritäts-Dateien

- **Rolle:** Verifier („Bauen wir es richtig?" — gegen Plan, DoD und Spec; Baseline-Regelwerk
  `modul-11-verification.md`). Nicht der Reviewer: Maintainability und Diff-gegen-Plan trägt
  R1 (`2026-10-07-slice-252-targets-authority-disjunkt-r1.md`).
- **Gegenstand:** `slice-252` im Stand `a73144dc` (Kette `59d21c07` feat → `c41f8725` R1-Report →
  `a73144dc` R1-Einarbeitung).
- **Eingang:** Slice-Plan §2 (DoD), Commit-Botschaften als Sensor-Behauptung des Implementers,
  `DC-FA-TGT-001` (Lastenheft 0.97.1), Spezifikation §DC-FA-TGT-001.a (Schritte 1, 3, 5, 5a, 6),
  Schema-Zeile unter `SPEC-005`, `SPEC-088`, `ADR-0101` (Proposed), Geschichte-Anhang `ADR-0100`,
  Baseline `ai-harness-course` `v6.16.0` `kurs/de/grundlagen/harness-dateien.md`
  §„Ein Index, mehrere Eigentümer" (lokaler Klon, `git show`).
- **Modell-ID:** claude-opus-5-5 · **Datum:** 2026-10-07

## Selbst gefahrene Sensoren

| Sensor | Ergebnis |
|---|---|
| `make gates` (Stand `a73144dc`, sauberer Baum) | Exit 0; letzte Zeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; doc-check `946 Datei(en) geprüft, 0 Befund(e)` |
| `make test` unter zehn Mutationen (unten) | jede Mutation Exit 2 mit dem erwarteten Test und der erwarteten Ursache, bis auf die `--print-config`-Zeile (Exit 0) |
| Black-Box: drei Runtime-Images aus `git archive` (`59d21c07^` = Vorher, `59d21c07`, `a73144dc`), je `docker run --rm --network none -v <fx>:/repo:ro <img> …`, Fixtures `chmod -R a+rX` | siehe §Black-Box |
| `make verify-closure-notes` | nicht anwendbar — der Slice liegt in `in-progress/` |

Zusatz-Images wieder entfernt, Fixtures nur im Scratchpad; der Arbeitsbaum ist bis auf diesen
Report unverändert.

## Bewusstes Brechen (Modul 11)

Je Schutzprüfung in `internal/hexagon/core/rules/targets.go` (bzw. Config-Rand/Lexik) gezielt
entfernt, `make test`, Ursache gelesen, Datei zurückgesetzt:

| Mutation | rot | Ursache (gelesen) |
|---|---|---|
| Schalter-Abfrage `!cfg.AuthorityDisjoint` entfernt | `TestCheckTargetsAuthorityDisjunktGrenzen` Fall 0, `TestCheckTargetsAuthorityListe` | `gate-declared-twice` ohne Schalter — richtig |
| Führende-Zelle-Filter `if !r.lead {continue}` entfernt | `…DisjunktFuehrendeZelle:518` | `gates` aus der Vertrag-Spalte gemeldet — richtig |
| Makefiles-Unabhängigkeit zurückgenommen (`len(cfg.Makefiles)==0` vorn) | `…DisjunktOhneMakefiles:543` | „bekam []" — richtig |
| `path.Clean`-Dublette (`key := a`) | `…DisjunktGrenzen` Fall 2 | `./harness/README.md` gegen `harness/README.md` gemeldet — richtig |
| Heimat-Vergleich `first == a` entfernt | `…DisjunktGrenzen` Fall 2 | Doppelung innerhalb der Heimat gemeldet — richtig (nur über den kombinierten Fall, siehe V5) |
| `exempt` wirkt (Ausschluss `clean`) | `TestCheckTargetsAuthorityDisjunkt:476` | `clean`-Befund fehlt — richtig |
| Kern aus (`if true { return nil }`) | alle drei Disjunkt-Tests | „Befunde = []" — bestätigt die Feat-Botschaft |
| `targets.go` aus `59d21c07` (vor R1) | `…FuehrendeZelle`, `…OhneMakefiles` | bestätigt die Behauptung der R1-Einarbeitung wörtlich |
| `targets.go` aus `59d21c07^` | Compile-Fehler (`undefined: ReasonGateDeclaredTwice`) | erwartbar; die „richtige Ursache" trägt die Kern-aus-Mutation |
| `applyTargets` reicht den Schalter nicht durch | `TestDecode_TargetsAuthorityDisjoint:1345` | richtig |
| Grund-Code aus `AllReasons` | `TestAllReasonsDeckungGegenSpezifikationGrundCodes`, `TestReasonTextDeckungGegenAllReasons` | §4-Lockstep — richtig |
| `--doctor`-Klartext entfernt | `TestReasonTextDeckungGegenAllReasons` | richtig |
| `--print-config`-Zeile entfernt | **grün** | kein Test hält die Zeile (V5) |

## Black-Box

**Ohne Schalter, Vorher (`59d21c07^`) gegen Nachher (`a73144dc`), stdout + stderr + Exit:**
vier Fixtures (Werkzeug-Teil mit „eingehängt in `make gates`"; fünf Autoritäts-Einträge mit
`./`-Formen, `exempt-targets`, Phantom und Undokumentiert; zwei Dateien ohne `makefiles`;
fehlende Datei ohne `makefiles`) je mit `--enable targets`, `--json`, `--doctor`,
`--doctor --json`, Default und Default `--json`, dazu der Selbstlauf übers Repo in sechs
Formen: **alle byte-identisch** — **außer** einer Eingabe: eine Tabellen-Datei mit einer Zeile,
die nur aus `|` besteht. Vorher Exit 0, Nachher Go-Panic, Exit 2 (V1).

**Mit Schalter (Nachher):**

| Fall | Ergebnis | Zusage |
|---|---|---|
| Werkzeug-Teil erwähnt `make gates` in der Vertrag-Spalte | Exit 0 (Stand `59d21c07`: Befund `mk/tool.md:5 gates`) | erfüllt (H1 behoben) |
| dieselbe Datei führt `make ci` zusätzlich in der ersten Zelle | `mk/tool.md:6 ci gate-declared-twice … (zuerst in README.md)`, auch `--json` | erfüllt |
| `[./z.md, a.md, z.md, b.md, ./sub/c.md]` | Heimat `./z.md` für `x`/`clean`, `a.md` für `y` (z führt `y` nicht); `a.md:3` zwei `make x` in der ersten Zelle ⇒ ein Befund; Erwähnung `a.md:4` still; `clean` trotz `exempt-targets` gemeldet; `z.md`/`./z.md` einmal gelesen | erfüllt |
| zwei Dateien, kein `makefiles` | Befund (Stand `59d21c07`: still) | erfüllt (M1 behoben) |
| `[a.md, ./a.md, a.md]` | Exit 0 | erfüllt |
| `[a.md]` | Exit 0 (wirkungslos) | erfüllt |
| Symlink-Alias `alias.md -> a.md`, echter Adapter | `alias.md:3 x gate-declared-twice (zuerst in a.md)` | wie benannt (L1) |
| Verzeichnis-Symlink `lnk -> real`, `[real/a.md, lnk/a.md]` | Befund an `lnk/a.md:3` | wie benannt |
| `[a.md, nope.md]`, Schalter, kein `makefiles` | Exit 2 mit Dateinamen | fail-closed erfüllt |
| `authority-disjoint: "true"` | Exit 2 (strikt Bool) | erfüllt |
| Fence, `\|` in der ersten Zelle, Code-Span mit `|`, eingerückte Zeile | Fence und eingerückt still; `\|`- und Span-Zelle führen beide Tokens | erfüllt (Zell-Zerlegung wie Schritt 5a sagt) |

`--doctor` zeigt den Klartext „Target in mehr als einer Autoritäts-Doku deklariert (Teile des
Gate-Index nicht disjunkt)"; `--print-config` führt die Zeile `authority-disjoint: false`.

**Trifft die Umsetzung „disjunkt" der Baseline?** Ja. `v6.16.0` sagt „Kein Target steht in zwei
Teilen. Sonst führen zwei Zeilen dasselbe Target" und beschreibt den Werkzeug-Teil als einen,
der Gates in `make gates` einhängt. Die erste Zelle als führende Zelle trifft beides: der
Einhänge-Vermerk ist still, die zweite führende Zeile laut; die verlinkte Target-Zelle
(eine als Link gesetzte Target-Zelle wie in `harness/README.md`) zählt ebenfalls als führend. Datei- statt Abschnitts-Granularität ist
Out-of-Scope und benannt.

## Verdikt je DoD-Punkt

| # | DoD-Punkt | Verdikt |
|---|---|---|
| 1 | Kern `gate-declared-twice` (Fundstelle, Reihenfolge, exempt, eine Datei), Tests rot ohne Änderung | **logisch bestätigt** (Black-Box + Mutationen), **aber nicht abnahmefähig wegen V1**: die Änderung am geteilten Extraktor bricht die Zusage „ohne Schalter byte-identisch" für eine Eingabeklasse |
| 2 | Config-Rand und Lexik | **bestätigt** (Bool strikt, Default aus, `--print-config`, `AllReasons`, `--doctor`); Testlücke an der Gerüst-Zeile (V5) |
| 3 | Lastenheft/Spezifikation/ADR, Antwort im Hinweis, Black-Box-Probe mit Ergebnis im Plan | **teilweise**: Lastenheft 0.97.1 mit Historie (`MR-032`), Schritt 5a, Schema, `SPEC-088`, `ADR-0101` vorhanden; Spiegel-Lücke V3 und Doppelsatz V2; Antwort im Hinweis offen (erwartet); Probe gefahren (Botschaft `59d21c07`), **Ergebnis nicht im Plan notiert** (V4); ihr Schluss „byte-identisch" ist durch V1 für die R1-Fassung widerlegt — die Probe lief vor der Einarbeitung |
| 4 | `make gates` grün | **bestätigt**, selbst gefahren, Exit 0, Coverage 94,70 % |
| 5 | Review, Report unter `docs/reviews/` | **bestätigt** (R1, eigener Kontext) |
| 6 | Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | **offen** (erwartet, kein Befund) |

**R1-Einarbeitung gegen den Report:** H1 (führende Zelle) und M1 (ohne `makefiles`) im Code
behoben und je durch einen Test gehalten, der gegen `59d21c07` aus dem richtigen Grund rot wird;
L1 (Symlink-Grenze in Lastenheft, Schritt 5a, ADR, Kommentar), L2 (Rumpf-Satz, `--repair`-Satz,
`rawTargets`-Kommentar) und L3 (Versionsnummer aus dem CR-Nachtrag) eingearbeitet; I1/I2 zur
Kenntnis, wie vorgesehen. **Aber:** die H1-Behebung bringt V1 mit, und die H1-Präzisierung ist
nicht in alle Spiegel gewandert (V3).

## Findings

### V1 — HIGH — Eine Tabellenzeile, die nur aus `|` besteht, lässt `targets` panisch abbrechen — auch ohne Schalter

- **pfad:** `internal/hexagon/core/rules/targets.go:250`
- **befund:** `extractDocTargets` greift seit `a73144dc` unbedingt auf `tableCells(line)[0]` zu.
  Für die Zeile `|` (auch `|` mit Leerzeichen dahinter) liefert `SplitPipeTableLine` eine leere
  Zellliste (führender und abschließender Pipe sind dasselbe Zeichen) — Index außerhalb des
  Bereichs. Der Extraktor dient auch Richtung 1 und 2, also trifft es jeden `targets`-Lauf mit
  einer solchen Zeile in `doc-tables` oder `authority`, **unabhängig vom Schalter**. Black-Box:
  Vorher Exit 0 `0 Befund(e)`, Nachher `panic: runtime error: index out of range [0] with length 0`
  mit Stacktrace über `phantomFindings`, Exit 2. Das bricht die Zusage „ohne Schalter
  byte-identisch" (Lastenheft, Schritt 5a, `ADR-0101` Entscheidung 1, Plan §1/§2). Die übrigen
  Aufrufer von `tableCells` (`planning_waves.go:280`, `structure_tableorder.go:162`,
  `structure_tablecell.go:120`) prüfen die Länge, dieser nicht. Kein Test hält den Fall; die
  R1-Probe lief gegen `59d21c07`, die Probe der Feat-Botschaft ebenso.
- **verifizierbar:** ja — Fixture `README.md` mit Tabelle und einer Zeile `|`, `--enable targets`,
  ohne `authority-disjoint`.
- **klasse:** `regression-im-geteilten-extraktor`

### V2 — LOW — Doppelter Satzrest im Lastenheft

- **pfad:** `spec/lastenheft.md:3360-3361`
- **befund:** „Ohne den Schalter ist der Befundsatz byte-identisch ([`DC-QA-02`]…)." wird von
  einer zweiten Zeile „Befundsatz byte-identisch ([`DC-QA-02`]…)." gefolgt — eingeführt in
  `59d21c07`, in R1 und der Einarbeitung nicht bemerkt. Vertragstext des Rang-1-Dokuments.
- **verifizierbar:** nein — Lesen.
- **klasse:** `textfehler-im-vertrag`

### V3 — LOW — Die H1-Präzisierung („führt", erste Zelle) fehlt in fünf Spiegeln

- **pfad:** `spec/spezifikation.md:3463` (Schema-Zeile `targets.authority-disjoint`:
  „als Tabellenzeile in mehr als einer `targets.authority`-Datei steht … an jeder späteren
  Tabellenzeile"); `spec/spezifikation.md:3602` (`SPEC-088`: „Target als Tabellenzeile in mehr
  als einer …"); `spec/lastenheft.md:3379` (Akzeptanzkriterium „Disjunktheit (Negative)": „als
  Tabellenzeile in zwei Autoritäts-Dateien steht"); `internal/hexagon/core/model/config.go:858`;
  `internal/adapter/driving/cli/config_template.go:258`
- **befund:** Gegen den Gegenstand geprüft: ein Target, das in einer zweiten Datei *als
  Tabellenzeile steht*, aber nur in der Vertrag-Spalte, meldet nichts (Black-Box, erste Zeile der
  Schalter-Tabelle). Diese fünf Stellen sagen das Gegenteil zu; nur Lastenheft-Rumpf, Schritt 5a,
  `ADR-0101` und der Funktionskommentar sind präzisiert. Gleiche Klasse wie R1-L2 (`MR-025`) und
  zum fünften Mal am Modul eine Grenze, deren Beschreibung weiter reicht als der Gegenstand.
- **verifizierbar:** ja — Black-Box (Fixture f1 mit Schalter).
- **klasse:** `grenze-gegen-beschreibung-statt-gegenstand-geprueft`

### V4 — LOW — DoD 3 verlangt das Probe-Ergebnis im Plan; es steht nur in der Commit-Botschaft

- **pfad:** `docs/plan/planning/in-progress/slice-252-targets-authority-disjunkt.md` §2, dritter Punkt
- **befund:** „Ergebnis im Plan notiert" — der Plan trägt keine Zeile dazu; die 16 Vergleiche
  und zwei Selbstläufe nennt nur `59d21c07`. Die Probe lief vor der R1-Einarbeitung; nach V1 ist
  sie für den abzuschließenden Stand ohnehin neu zu fahren, dann mit der `|`-Zeile im Fixture.
- **verifizierbar:** nein — Lesen.
- **klasse:** `dod-beleg-am-falschen-ort`

### V5 — INFO — Zwei Zusagen ohne eigenen Test

- **pfad:** `internal/adapter/driving/cli/config_template.go:258`;
  `internal/hexagon/core/rules/targets_test.go:482-498`
- **befund:** (a) Die `--print-config`-Zeile (DoD 2) hält kein Test — ihr Entfernen bleibt grün;
  vorhanden ist sie (Black-Box). (b) Das Akzeptanzkriterium „mit dem Schalter und einer
  Doppelung nur innerhalb einer Datei" ist in `…DisjunktGrenzen` Fall 1 nur mit **einer**
  Autoritäts-Datei gestellt, wo die `< 2`-Kurzschluss-Rückkehr greift; den Heimat-Vergleich fängt
  nur Fall 2 über die `path.Clean`-Dublette. Ein Fall mit zwei verschiedenen Dateien und einer
  Doppelung in der Heimat fehlt. Kein Fehlverhalten (Black-Box f3 korrekt).
- **verifizierbar:** ja — Mutationen „print-config" und „Heimat".
- **klasse:** `grenze-ohne-negativtest`

## Erwartet offen, kein Befund

Closure-Notiz (§7), Risiko-Ausgänge (§6), Beobachtungs-Register, drei Paarungen, `ADR-0101`
`Proposed`, Antwort im Hinweis-Dokument, Release-Prep-Flächen (Handbuch, README, CHANGELOG;
`AGENTS.md` §5 Regel 17).

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 1 | V1 |
| MEDIUM | 0 | — |
| LOW | 3 | V2, V3, V4 |
| INFO | 1 | V5 |

## Verdikt

**Nicht bestätigt.** Die Disjunktheits-Logik ist richtig gebaut, trifft die Baseline-Lesart von
„disjunkt", und jede ihrer Schutzprüfungen ist durch einen Test gehalten, der ohne sie aus dem
richtigen Grund rot wird. V1 blockiert: die R1-Einarbeitung hat eine Regression in den geteilten
Extraktor gebracht, die die Kernzusage „ohne Schalter byte-identisch" für jede Konfiguration
mit `targets` bricht — Längenprüfung wie bei den anderen `tableCells`-Aufrufern, Test mit einer
`|`-Zeile, Black-Box-Probe neu fahren (V4). V2/V3 vor der Closure nachziehen.
