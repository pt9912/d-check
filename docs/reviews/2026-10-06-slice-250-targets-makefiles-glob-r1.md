# Review R1 — slice-250: `targets.makefiles` nimmt Glob-Muster an

- **Review-Art:** Code. Geprüft wurde gegen den Slice-Plan (`slice-250`, §1 Abgrenzung,
  §3 Spiegel-Liste, §6 Risiken), gegen den eingehenden CR von `ai-harness-init` vom
  2026-10-06 (sechs Akzeptanzkriterien, Abgrenzung), gegen `ADR-0099` (Proposed) und
  `ADR-0044`/`ADR-0058` (Bezug), gegen `MR-025`/`MR-032` und gegen die Hard Rules
  `AGENTS.md` §3.1/§3.2/§3.7/§3.8 sowie §5 Regel 15/17. Gegenstand ist die
  Maintainability; die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-250` · Commit `384a4770`. Er umfasst die Kern-Expansion in
  `internal/hexagon/core/rules/targets.go`, den Config-Rand, das Gerüst, Lastenheft
  0.95.0, die Spezifikation (Schritt 1a und die Schema-Zeile), `ADR-0099` und den
  ADR-Index.
- **Skill:** `reviewer.md` @ 1.16.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-06
- **Eingangs-Kontext:** Slice-Plan (`slice-250`). Aus dem Lastenheft `DC-FA-TGT-001`,
  aus der Spezifikation `DC-FA-TGT-001.a` (Schritt 1a) und `SPEC-005` (Schema-Zeile).
  Dazu `ADR-0099`, der CR-Wortlaut, der Filesystem-Adapter
  (`internal/adapter/driven/fs/fs.go`) und das MemFS der Tests
  (`internal/hexagon/core/coretest/memfs.go`). Vorherige Findings am Modul `targets`:
  Unter `docs/reviews/` liegt außerhalb des Archivs kein Report dazu.
- **Proben:** `make build`, danach Black-Box-Läufe des Images
  (`docker run --rm --network none -v <fx>:/repo:ro d-check:latest --enable targets`)
  gegen neun Fixtures im Scratchpad. Den Bestand der Unit-Tests hat `make test`
  bestätigt (grün). Der Arbeitsbaum wurde nicht verändert, `git status` ist bis auf
  diesen Report leer.

Probe-Ergebnisse (Root-Makefile mit `help:`, Autoritäts-Doku mit `make help` und `make alpha`):

```text
Fixture  makefiles                                   Lage                                         Ergebnis
A        Makefile, "harness/mk/*.mk"                 harness -> real (Symlink), real/mk/a.mk      Treffer: gate-undocumented secret @ harness/mk/a.mk:3
A2       Makefile, "**/mk/*.mk"                      dieselbe Lage                                 kein Treffer im Symlink-Zweig, 0 Befunde
B        Makefile, "harness/mk/*.mk"                 harness/mk/b.mk -> ../../other/b.mk (secret)  0 Befunde — b.mk still ausgelassen
B2       Makefile, harness/mk/a.mk, harness/mk/b.mk  dieselbe Lage, wörtlich                       gate-undocumented secret @ harness/mk/b.mk:1
C        "./Makefile", "Make*"                       Makefile mit secret                           ZWEI gate-undocumented secret (./Makefile:3 und Makefile:3)
D        "harness/*/build/*.mk"                      harness/x/build/a.mk existiert                Exit 2 „findet … keine Datei"
D2       "build/*.mk"                                build/a.mk                                    Treffer (build/ gewandert)
E        Makefile, "Makefile/*.mk"                   Präfix ist eine Datei                         Exit 2 „keine Datei" (wie zugesagt)
F        "harness/*/*.mk", "harness/mk/*"            zwei Globs treffen a.mk                       ein Befund je Regelzeile (wie zugesagt)
```

## Findings

### M1 — MEDIUM — Ein Symlink im Präfix wird verfolgt, sobald er nicht das letzte Präfix-Segment ist

- **quelle:** `DC-FA-TGT-001` (Grenze „symbolische Links werden nicht verfolgt"), `ADR-0099` Entscheidung 4, `AGENTS.md` §5 Regel 13/15
- **pfad:** `internal/hexagon/core/rules/targets.go:206`; Vertrag in `spec/lastenheft.md:3312`, `spec/spezifikation.md:2724`, `docs/plan/adr/0099-targets-makefiles-glob.md:59`
- **befund:** `fsys.Kind(dir)` klassifiziert per `os.Lstat` nur die **letzte** Komponente des Präfixes. Ein Symlink davor wird vom Betriebssystem aufgelöst, und die Wanderung läuft durch ihn hindurch. Fixture A meldet deshalb einen Befund aus `real/mk/a.mk` über `harness -> real`. Fixture A2 (`**/mk/*.mk`) sieht dieselbe Datei nicht. Lastenheft, Spezifikation („ist es ein Symlink, gibt es keine Treffer"), ADR und Commit-Botschaft („keine Symlinks") sagen das Gegenteil. Ein solcher Symlink kann außerdem aus der Repo-Wurzel hinauszeigen, obwohl CR-Kriterium 5 das Hinausführen nur beim Laden abfängt.
- **verifizierbar:** ja — mit der Black-Box-Probe A gegen A2. `make test` deckt den Fall nicht, weil das MemFS keine Pfadauflösung über Zwischen-Segmente kennt.
- **klasse:** `grenze-gegen-beschreibung-statt-gegenstand-geprueft`

### M2 — MEDIUM — Ein Fragment, das per Symlink eingebunden ist, fällt bei einem Glob still aus, wörtlich wird es gelesen

- **quelle:** `DC-FA-TGT-001` („Jede Treffer-Datei wird wie ein wörtlicher Eintrag gelesen"), CR-Anlass („eine vergessene Datei macht deren Targets still ungeprüft"), `AGENTS.md` §3.8
- **pfad:** `internal/hexagon/core/rules/targets.go:207`. Die Ursache liegt in `walkAllFiles` (`internal/hexagon/core/rules/file.go:66`): Ein Eintrag mit `KindSymlink` ist weder Verzeichnis noch Datei.
- **befund:** `harness/mk/b.mk -> ../../other/b.mk` passt textlich auf `harness/mk/*.mk`. Trotzdem endet der Lauf mit 0 Befunden und Exit 0, weil die Regel `secret` nie gelesen wird (Fixture B). Ist derselbe Pfad wörtlich eingetragen, meldet der Lauf ihn, denn `ReadFile` folgt dem Link (Fixture B2). Die Lastenheft-Grenze „symbolische Links werden nicht verfolgt" benennt diesen Unterschied zwischen wörtlichem Eintrag und Glob nicht. Sie verschweigt auch, dass der Ausfall ohne Meldung geschieht, solange irgendein anderes Fragment trifft. Das ist genau die Lücke, die der CR mit dem Glob schließen wollte, und sie liegt in einem Gate-Pfad.
- **verifizierbar:** ja — mit der Black-Box-Probe B gegen B2.
- **klasse:** `wörtlich-vs-glob-asymmetrie-still`

### M3 — MEDIUM — Die benannten Grenzen von Schritt 1a und die Dublette aus zwei Globs haben keinen Test

- **quelle:** Reviewer-Skill Prüffrage 13; `BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet` (2×, im Plan §8 gesichtet)
- **pfad:** `internal/hexagon/core/rules/targets_test.go:221-275`
- **befund:** Mehrere Zusagen der Spezifikation tragen keinen Test: SKIP_DIRS unterhalb des Präfixes, Symlink als Präfix, Symlink darunter, Präfix als Datei, Glob im ersten Segment, Treffer über `?`/`[` sowie die Dublette aus „zwei Globs", die das Lastenheft ausdrücklich zusagt. Das MemFS kann Symlinks abbilden (`memfs.go:52`), wird dafür aber nicht genutzt. Weil diese Fälle fehlen, sind M1 und M2 im Testbestand unsichtbar. Plan §8 behauptet, „die Tests tragen die Nachbar-Fälle"; das gilt für die Dublette aus wörtlichem Eintrag und Glob, nicht für die Grenzen.
- **verifizierbar:** ja — `make test` bleibt grün, auch wenn das Symlink- oder SKIP_DIRS-Verhalten gekippt wird.
- **klasse:** `grenze-ohne-negativtest`

### L1 — LOW — Die Dubletten-Erkennung vergleicht Zeichenketten statt Pfade

- **quelle:** `DC-FA-TGT-001`, Akzeptanzkriterium „Glob (Dublette)" (`spec/lastenheft.md:3326`), CR-Kriterium 4
- **pfad:** `internal/hexagon/core/rules/targets.go:190`
- **befund:** Der Config-Rand nimmt `./Makefile` an. Ein Glob liefert dieselbe Datei als `Makefile`. `seen` vergleicht die Zeichenketten, deshalb wird die Datei doppelt gelesen, und für dieselbe Regelzeile entstehen zwei `gate-undocumented` (Fixture C). Das widerspricht der Zusage „je Regelzeile höchstens ein Befund". Eine Dublette aus zwei Globs wird korrekt erkannt (Fixture F).
- **verifizierbar:** ja — mit der Black-Box-Probe C.
- **klasse:** `dublette-ohne-pfad-normalisierung`

### L2 — LOW — Ob SKIP_DIRS greifen, hängt davon ab, wo im Muster das erste Glob-Segment steht; „wie das Modul `file`" stimmt nicht

- **quelle:** `DC-FA-TGT-001` / Spezifikation Schritt 1a; `AGENTS.md` §5 Regel 15 (Commit-Botschaft „SKIP_DIRS wie file")
- **pfad:** `internal/hexagon/core/rules/targets.go:206-207`; `spec/lastenheft.md:3309-3311`
- **befund:** `build/*.mk` trifft, weil `build` im Präfix steht und nur Verzeichnisse darunter gefiltert werden (Fixture D2). Das Modul `file` wandert dagegen immer ab der Wurzel und betritt `build/` nie. `harness/*/build/*.mk` nennt `build` ebenso ausdrücklich, endet aber in Exit 2 mit „findet … keine Datei", obwohl die Datei existiert (Fixture D). Die Spezifikation deckt das formal ab („unterhalb des Präfixes"). Die Formel „wie beim Modul `file`" in Lastenheft und Commit-Botschaft beschreibt ein anderes Verhalten, und die Fehlermeldung lenkt die Diagnose auf einen falschen Grund.
- **verifizierbar:** ja — mit den Black-Box-Proben D und D2.
- **klasse:** `skip-dirs-positionsabhängig`

### I1 — INFO — Ein Fehler aus `Kind` wird als „keine Datei" gemeldet

- **quelle:** Maintainability
- **pfad:** `internal/hexagon/core/rules/targets.go:206`
- **befund:** Liefert `Kind` einen Fehler, der kein Nicht-Existieren ist (etwa EACCES beim Lstat), wird er verworfen, und der Lauf meldet „findet zum Makefile-Glob … keine Datei". Das Verhalten bleibt fail-closed (Exit 2), nur die Ursache in der Meldung ist falsch.
- **verifizierbar:** nein (Urteil; kein Gate)
- **klasse:** `fehlerursache-verschluckt`

### I2 — INFO — Release-Prep-Flächen sind unvollständig, aber nicht falsch

- **quelle:** `AGENTS.md` §5 Regel 17
- **pfad:** `docs/user/benutzerhandbuch.md:1512`, `docs/user/benutzerhandbuch.md:2645`
- **befund:** Das Handbuch nennt für `makefiles` nur Pfade und „fail-closed bei fehlender Datei". Glob und Exit 2 bei leerem Glob fehlen dort. Falsch ist keine der Aussagen. `operations.md` trifft keine Aussage zu `targets.makefiles`. CHANGELOG und README bleiben unberührt, so wie Plan §1 es vorsieht. Nachzuziehen in der Release-Prep.
- **verifizierbar:** nein
- **klasse:** `release-prep-nachzug`

## Negativbefunde

- **CR-Akzeptanzkriterien 1–6:** 1 (Expansion, echte Fundstelle) erfüllt (Fixtures A/B2/F, Unit-Test). 2 (wörtliche Einträge unverändert, fehlende Datei Exit 2) erfüllt. 3 (leerer Glob nicht still) erfüllt mit Exit 2 (Fixtures D/E, Unit-Test). 4 (Dublette) erfüllt bis auf L1. 5 (Pfad-Regel beim Laden) erfüllt; zur Laufzeit siehe M1. 6 (rotes Gegenbeispiel) erfüllt. Die offenen Fragen des CR beantworten `ADR-0099` (Exit 2) und das Gerüst (Glob-Beispiel).
- **CR-Abgrenzung und Plan §1:** `doc-tables`/`authority` bleiben wörtlich, `makefileRuleRe` und `exempt-targets` sind unverändert, `include` wird nicht aufgelöst, Handbuch/README/CHANGELOG sind nicht angefasst. Der Diff berührt nur die Dateien aus Plan §3. Der Filesystem-Adapter ist nicht geändert, die Rückführungs-Bedingung aus §4 ist also nicht eingetreten. Geprüft, ohne Befund.
- **„Ohne Glob-Eintrag byte-identisch":** Die neue `seen`-Deduplizierung wirkt auch auf wörtliche Doppeleinträge. Der Befundsatz bleibt trotzdem identisch, weil `model.SortFindings` (`run.go:122`) gleiche Befunde ohnehin zusammenführt. Ein fehlender wörtlicher Eintrag bleibt in der Liste und bricht weiterhin mit Exit 2 ab. Geprüft, ohne Befund.
- **Randfälle ohne Befund:** Ein Glob im ersten Segment (`*.mk`, `**/…`) wandert ab der Wurzel. Bei Glob-Zeichen im Verzeichnis-Teil endet das Präfix am ersten Glob-Segment. Ist das Präfix eine Datei oder fehlt es, endet der Lauf mit Exit 2 (E). Die Dublette aus zwei Globs wird erkannt (F). Die Reihenfolge ist deterministisch: `List` ist sortiert, Treffer werden je Muster sortiert, Befunde über `SortFindings`. `DC-QA-02` ist gewahrt.
- **Adapter gegen MemFS:** `Kind`/`List` haben in beiden Lstat-Semantik für die letzte Komponente. Der einzige Unterschied ist die Auflösung von Zwischen-Segmenten (M1).
- **Hexagon-Richtung (`ADR-0005`):** Der Kern importiert nur `sort` und den Port `driven`. Geprüft, ohne Befund.
- **Netz außerhalb `external`, Inline-Suppression, Schwellen-Senkung (§3.2/§3.6):** keine. Geprüft, ohne Befund.
- **Kommentare (§3.7):** Die neuen Go-Kommentare zu `expandMakefiles`, `makefileGlobHits`, `globBaseDir`, zum Config-Rand und zu den Tests tragen Zusage oder Grenze. Sie enthalten keine Review-Historie, keine Slice-Nummern und keine Mess-Labels. Geprüft, ohne Befund.
- **§3.8 (Zusage nur über die Scan-Menge):** Die Anforderung nennt die Unabhängigkeit von `scan.*`, SKIP_DIRS und Symlinks als Grenzen. Die Lücken liegen in deren Genauigkeit (M1, M2, L2), nicht in ihrem Fehlen.
- **`MR-032`:** Lastenheft-Bump auf 0.95.0 und Historie-Zeile sind vorhanden, der Status ist Draft. Die Historie der Spezifikation ist nachgezogen. Geprüft, ohne Befund.
- **`MR-025` (Spiegel):** Lastenheft (Beschreibung, Akzeptanzkriterien, Out-of-Scope), Spezifikation (Schritt 1a, Schritt 2 mit Verweis, Schema-Zeile, Historie) und das `--print-config`-Gerüst sind nachgezogen. Die `--doctor`-Klartexte sind nicht betroffen, weil Exit 2 kein Grund-Code ist. Handbuch und `operations.md` gehören zur Release-Prep (I2).
- **`ADR-0099` Form:** Das `Schärft:`-Feld nennt die Kennungen `DC-FA-TGT-001.a` und `SPEC-005`. Der Marker `d-check:status-provenance` an `slice-250` zeigt die Entstehung und begründet keine Entscheidung. Der Re-Evaluierungs-Trigger ist vorhanden, `Proposed` ist eine bewusste Wahl. Geprüft, ohne Befund außer der Grenz-Aussage in M1.
- **Commit-Botschaft (§5 Regel 15):** Die Proben-Aussagen (rot vor der Änderung, Black-Box Exit 1/2) kann dieses Review nicht nachmessen; das ist Sache des Verifiers. Überdehnt sind „keine Symlinks" (M1) und „SKIP_DIRS wie file" (L2).

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 3 | M1, M2, M3 |
| LOW | 2 | L1, L2 |
| INFO | 2 | I1, I2 |

Wiederkehrende Finding-Klasse für die Closure: `grenze-gegen-beschreibung-statt-gegenstand-geprueft`
(M1, nahe `AGENTS.md` §5 Regel 13) und `grenze-ohne-negativtest` (M3; Nachbar von
`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`).

## Verdikt

**Nicht merge-bereit in dieser Form.** Kein Finding ist HIGH. M1 und M2 sind
Abweichungen zwischen Vertrag und Code an einer benannten Grenze eines Gates. Beide
müssen vor der Closure geklärt sein: Entweder ändert sich der Code, oder die Grenze
wird so geschrieben, wie der Code sich verhält. M3 gehört zu beiden. L1/L2 sollten
behoben werden, blockieren aber nicht. I1/I2 dienen nur der Information.
