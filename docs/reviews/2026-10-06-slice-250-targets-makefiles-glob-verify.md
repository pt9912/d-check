# Verifikation — slice-250: `targets.makefiles` nimmt Glob-Muster an

- **Rolle:** Verifier (Modul 8/11). Die Frage lautet „Bauen wir es richtig?“. Geprüft
  wird gegen Plan-DoD (§2), Abgrenzung (§1), CR-Akzeptanzkriterien und
  `DC-FA-TGT-001`, nicht gegen die Maintainability (das ist R1).
- **Gegenstand:** `slice-250` @ `c67a1c02` (Kette `384a4770` feat → `276b5c79` R1-Report →
  `c67a1c02` R1-Einarbeitung).
- **Eingang:** Slice-Plan `docs/plan/planning/in-progress/slice-250-targets-makefiles-glob.md`;
  eingehender CR `docs/plan/cr/2026-10-06-cr-eingehend-ai-harness-init-targets-makefiles-glob.md`;
  `spec/lastenheft.md` §`DC-FA-TGT-001` (0.95.1); `spec/spezifikation.md`
  §`DC-FA-TGT-001.a` Schritt 1a und die Schema-Zeile; `ADR-0099` (Proposed);
  R1-Report `docs/reviews/2026-10-06-slice-250-targets-makefiles-glob-r1.md`.
- **Modell-ID:** claude-opus-5-5 · **Datum:** 2026-10-06
- **Arbeitsbaum:** Jede Mutation wurde per `git checkout --` zurückgenommen. `git status` ist am Ende bis
  auf diesen Report leer.

## Selbst gefahrene Sensoren

| Sensor | Ergebnis |
|---|---|
| `make gates` (HEAD, sauberer Baum) | Exit 0, `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; Coverage 94.60 % (Schwelle 93 %), semgrep 0 Befunde |
| `make build` | `d-check:latest` aus HEAD |
| Vorher-Image | `git archive 384a4770^` → `docker build --target runtime -t d-check:pre250` |
| `make test` × 9 Mutationen | siehe §Bewusstes Brechen |
| Black-Box, 18 Fixtures × 2 Images | siehe §Black-Box |

## Bewusstes Brechen (Modul 11)

Jede Mutation lief einzeln, danach `make test` und die gelesene Fehlerursache.

| # | Mutation | Rot in | gelesene Ursache | richtiger Grund? |
|---|---|---|---|---|
| M-pre | `targets.go` ← `384a4770^` | alle 4 Glob-Tests, 9/9 Grenz-Fälle | `kann das Makefile "harness/mk/*.mk" nicht lesen … not found` — Glob als Dateiname gelesen | ja |
| M-feat | `targets.go` ← `384a4770` | 3 Grenz-Fälle | Symlink im Präfix: `err=<nil>`, Befund aus `harness/mk/a.mk`; Symlink-Treffer: `err=<nil>`; `./Makefile`: `nicht lesen … not found: ./Makefile` | 2× ja, 1× nur teilweise (V1) |
| A | Präfix-Kette → nur `Kind(dir)` | Grenzen/Symlink im Präfix | `fail-closed mit "keine Datei" erwartet, bekam err=<nil>` | ja |
| B | `case KindSymlink` deaktiviert | Grenzen/Symlink-Treffer ist laut | `fail-closed mit "Symlink" erwartet, bekam err=<nil>` | ja |
| C | Dubletten-Schlüssel `path.Clean(p)` → `p` | Grenzen/Dublette über bereinigten Pfad | `nicht lesen … not found: ./Makefile` (MemFS) | nur teilweise (V1) |
| D | `validateSegmentGlobs(…, nil)` in `configyaml.go` | `TestDecode_TargetsMakefilesGlob` | `"harness/mk/[x.mk": Konfigurationsfehler erwartet` | ja |
| E | `isSkipDir` im Walk aus | Grenzen/SKIP_DIR unter dem Präfix | Befunde enthalten `harness/build/x.mk:skipped` | ja |
| F | Leer-Prüfung `len(hits)==0` aus | `…GlobLeer`, Grenzen/Symlink-Präfix, Präfix-Datei | `bekam <nil>` | ja |
| G | Dubletten-Erkennung ganz aus | `…Glob` (zwei `delta`-Befunde), Fragezeichen/Klasse, zwei Globs | doppelte Befunde je Regelzeile | ja |

Anmerkung zu `TestCheckTargetsMakefileGlobLeer`: Unter dem Vorher-Code war ein Glob schon
Exit 2, nur aus dem falschen Grund („nicht lesen“). Gegen diesen Stand unterscheidet der Test
deshalb allein am Wortlaut. Den eigentlichen Rot-Beleg für die Leer-Prüfung liefert Mutation F.

## Black-Box (echter Dateisystem-Adapter)

Aufruf: `docker run --rm --network none -v <fx>:/repo:ro d-check:<img> --enable targets`. Fixtures mit
`chmod -R a+rX`. Die Autoritäts-Doku listet `help`, `alpha` und `beta`.

| Fixture | `makefiles` / Lage | HEAD | `pre250` |
|---|---|---|---|
| f1 | `Makefile, "harness/mk/*.mk"` | Exit 1, `harness/mk/a.mk:4 secret gate-undocumented`; `sub/c.mk` und `r.txt` nicht erfasst | Exit 2 „nicht lesen“ |
| f2 | dazu wörtlich `harness/mk/a.mk` | Exit 1, **ein** Befund | Exit 2 |
| f2b / f2c | `"./Makefile","Make*"` bzw. umgekehrt | je **ein** `secret`, gemeldet in der Form der ersten Nennung | Exit 2 |
| f3 | Glob, Verzeichnis ohne `.mk` | Exit 2 `findet zum Makefile-Glob "harness/mk/*.mk" keine Datei` | Exit 2 (anderer Grund) |
| f4 | `harness -> real` (Präfix, erste Komponente) | Exit 2 „keine Datei“, `real/mk/a.mk` nicht gelesen | — |
| f4b | `harness/x -> ../real` (mittlere Komponente) | Exit 2 „keine Datei“ — R1-M1 behoben | — |
| f5 | `harness/mk/b.mk -> ../../other/b.mk` | Exit 2 `folgt dem Symlink "harness/mk/b.mk" … nicht … wörtlich … eintragen` — R1-M2 behoben | — |
| f5b | Symlink-Verzeichnis unter Präfix, `harness/**/*.mk` | nicht betreten, `other/b.mk` nicht gelesen | — |
| f5c | dasselbe, `harness/mk/*` trifft den Verzeichnis-Symlink | Exit 2 (Symlink-Hinweis) — deckt sich mit Schritt 1a | — |
| f6 | nur wörtliche Einträge | stdout **und** stderr byte-identisch mit `pre250` (auch `--json`) | |
| f6b | wörtlich fehlende Datei | byte-identisch, Exit 2 | |
| f6c | wörtlich `a.mk` + `./harness/mk/a.mk` | 1× `secret` (vorher 2×) — die im Lastenheft 0.95.1 benannte Ausnahme von der Byte-Identität | |
| f7a | `"harness/mk/[x.mk"` | Exit 2 beim Laden, `kein gültiges Glob` | Exit 2 erst im Lauf |
| f7b / f7c | `"../*.mk"` / `"/abs/*.mk"` | Exit 2 beim Laden, byte-identisch | |
| f8 | `"./harness/mk/*.mk"` | Exit 1, Treffer `./harness/mk/a.mk:4` | Exit 2 |
| f9 | `"**/*.mk"` mit `build/x.mk` | `build/` nicht betreten, `r.mk` erfasst | Exit 2 |
| Repo selbst | `--enable targets` und Default-Lauf | byte-identisch zwischen beiden Images | |

`--print-config` zeigt `makefiles: [Makefile, "harness/mk/*.mk"]` mit dem Hinweis auf Glob und Exit 2.

## CR-Akzeptanzkriterien und Abgrenzung

| CR | Verdikt | Beleg |
|---|---|---|
| 1 Expansion, echte Fundstelle | erfüllt | f1, f8, f9; `TestCheckTargetsMakefileGlob`, `…Doppelstern` |
| 2 wörtlich unverändert, fehlend ⇒ fail-closed | erfüllt | f6/f6b byte-identisch; `TestCheckTargetsFailClosed` (Bestand). Einzige Abweichung ist f6c, benannt |
| 3 leerer Glob nicht still | erfüllt (Exit 2) | f3, Mutation F |
| 4 Dublette einmal, kein falsches `gate-phantom` | erfüllt | f2, f2b, f2c, Mutation G |
| 5 Pfad-Regel beim Laden | erfüllt | f7a–c, Mutation D; ein Symlink aus der Wurzel hinaus ist zur Laufzeit Exit 2 (f4, f4b) |
| 6 rotes Gegenbeispiel | erfüllt | f1 (Exit 1, Fundstelle im Fragment), M-pre |
| Abgrenzung | eingehalten | `doc-tables`/`authority` ohne Glob-Validierung und ohne Expansion; `makefileRuleRe` und `exempt-targets` unverändert (Diff); keine `include`-Auflösung; Handbuch/README/CHANGELOG nicht berührt (`git diff 384a4770^ HEAD --stat -- docs/user README.md CHANGELOG.md` leer) |

Beide offenen CR-Fragen sind beantwortet: Exit 2 (ADR-0099 Entscheidung 2) und Glob-Beispiel im Gerüst.

## R1-Einarbeitung gegen den R1-Report

| R1 | Verdikt | Beleg |
|---|---|---|
| M1 Symlink im Präfix | behoben | `realDirChain`; f4b; Mutation A |
| M2 Symlink-Treffer still | behoben (Exit 2 mit Hinweis) | f5, f5c; Mutation B; Lastenheft und Spezifikation nachgezogen |
| M3 Grenzen ohne Test | behoben | `TestCheckTargetsMakefileGlobGrenzen` (9 Fälle); Mutationen A, B, E, F, G rot. Rest siehe V1 |
| L1 Dublette per Zeichenkette | behoben | f2b/f2c; im Kern über `path.Clean` |
| L2 SKIP_DIRS-Formulierung | Vertrag behoben | Lastenheft 0.95.1, Schritt 1a und ADR-0099 beschreiben jetzt das Verhalten im Code. Die Laufzeitmeldung bei `harness/*/build/*.mk` bleibt „keine Datei“ (V2) |
| I1 Fehlerursache aus `Kind` | behoben | eigener Fehlertext „kann das Präfix … nicht prüfen“ mit `%w` |
| I2 Release-Prep | offen, erwartet | AGENTS §5 Regel 17 |

## Verdikt je DoD-Punkt

| DoD | Verdikt |
|---|---|
| 1 Glob-Expansion im Kern + Tests nach CR 1–6 + Bewusstes Brechen | **bestätigt** — mit V1 (LOW) zum Rot-Grund des `./Makefile`-Falls |
| 2 Config-Rand, Pfad-Regel unverändert, `--print-config`-Gerüst | **bestätigt** |
| 3 Lastenheft (0.95.0 → 0.95.1, Historie, `MR-032`), Spezifikation, Schema, ADR | **bestätigt**. **Antwort-Vermerk im CR: offen** (`**Stand:** eingegangen`) — erwartet, kein Befund |
| 4 `make gates` grün | **bestätigt** (selbst gefahren, Exit 0) |
| 5 Review mit Report | **bestätigt** (R1-Report liegt vor, Einarbeitung geprüft) |
| 6 Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | **offen** — erwartet, kein Befund (§6 und §7 leer, ADR `Proposed`) |

## Neue Befunde

### V1 — LOW — Der `./Makefile`-Dublettenfall wird im Test aus einem MemFS-Grund rot, nicht wegen doppelten Lesens

- **pfad:** `internal/hexagon/core/rules/targets_test.go:323-326`; Commit-Botschaft `c67a1c02` („./Makefile doppelt gelesen“)
- **befund:** Mutation C (Schlüssel ohne `path.Clean`) und M-feat machen den Fall rot mit
  `kann das Makefile "./Makefile" nicht lesen … not found`. `coretest.MemFS.ReadFile` kennt
  `./`-Pfade nicht. Der Test erkennt die Mutation, aber am Symptom „`./Makefile` wird überhaupt
  gelesen“ und nicht an „zwei Befunde je Regelzeile“. Den behaupteten Grund (doppelte Befunde) zeigt erst
  die Black-Box über den echten Adapter (f2b/f2c gegen R1-Fixture C). Die Commit-Botschaft beansprucht also
  mehr, als der Test misst (AGENTS §5 Regel 15). Hängt die Konfigurations-Reihenfolge anders, wäre der Test
  mit Fix rot (erst `./Makefile` → MemFS-Lesefehler). Das ist eine Testfixture-Abhängigkeit.
- **verifizierbar:** ja — Mutation C.

### V2 — INFO — Wird ein SKIP_DIRS-Name unterhalb des Präfixes genannt, meldet der Lauf weiter „keine Datei“

- **pfad:** `internal/hexagon/core/rules/targets.go` (`makefileGlobHits`, Leer-Meldung)
- **befund:** Der Vertrag beschreibt das Verhalten jetzt korrekt (L2). Die Meldung zu
  `harness/*/build/*.mk` nennt aber weiterhin nur „keine Datei“, ohne Hinweis auf die
  übersprungenen Verzeichnisse. Das bleibt fail-closed und ist kein Vertragsbruch, nur ein Diagnose-Komfort.

## Negativbefunde

- Hexagon-Richtung, Netzlosigkeit, Suppressions: keine Änderung außerhalb von Plan §3. `make gates` (arch-check, semgrep, lint) ist grün.
- Determinismus (`DC-QA-02`): Bei wiederholten Läufen sind die Ausgaben gleich. Ohne Glob sind sie byte-identisch zum Vorher-Image (f6, f6b, f7b, Repo-Selbstlauf). Ausnahme ist der benannte Fall f6c.
- Rückführungs-Bedingung §4 (neuer Port, Adapter-Änderung): nicht eingetreten. `internal/adapter/driven/fs/` ist unberührt.
- Release-Prep-Flächen, ADR-Status, CR-Stand, §6/§7: offen wie erwartet.
