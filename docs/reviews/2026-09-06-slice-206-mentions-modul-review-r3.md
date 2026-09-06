# Review-Report R3 — slice-206 (Modul `mentions`, dritte Fassung der Grenz-Semantik)

- **Review-Art:** Code — geprüft wird die **dritte Fassung der Kern-Semantik** gegen den R2-Report, gegen Anforderung, Entscheid und Hard Rules; Maßstab sind [`DC-FA-MENT-001`](../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) (Lastenheft jetzt 0.86.4, fünfzehn Akzeptanzkriterien), [ADR-0084](../plan/adr/0084-mentions-eigenes-modul.md), [`spec/spezifikation.md`](../../spec/spezifikation.md) §`DC-FA-MENT-001.a` Schritt 5, [`AGENTS.md`](../../AGENTS.md) §3.1/§3.7/§3.8/§4/§5 und der Baseline-Kanon (`modul-13-quality-gates.md` §Hard Rule, `modul-11-verification.md` §Fitness Function ohne Standard-Tool)
- **Gegenstand:** `0f13f0a` (`fix(mentions)`, „die Haertung kippte die Fehlerpolitik — Grenze neu geschnitten", 8 Dateien, +407/−69)
- **Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) @ 1.13.0 · **Modell-ID:** claude-opus-5[1m] · **Datum:** 2026-09-06
- **Eingangs-Kontext:** [R1-Report](2026-09-06-slice-206-mentions-modul-review.md) (3 HIGH / 8 MEDIUM / 6 LOW), [R2-Report](2026-09-06-slice-206-mentions-modul-review-r2.md) (2 HIGH / 6 MEDIUM / 4 LOW / 4 INFO, blockierend), Slice-Plan [`slice-206`](../plan/planning/in-progress/slice-206-mentions-modul.md), `internal/hexagon/core/rules/mentions.go` + `mentions_test.go`, `internal/hexagon/core/rules/{workflows.go,reviews.go,paths.go,scan.go,planning.go}`, `.d-check.yml`, `.golangci.yml`, `Makefile`, `harness/sensors/mention-coverage.md`, `spec/lastenheft.md`, `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`
- **Vorherige Findings am gleichen Gegenstand:** R1 und R2 zu slice-206, davor R1/R2 zu slice-205. Dominante Klassen: `wortlaut-behauptet-pruefung-die-fehlt`, `selbstauskunft-zahl-reproduziert-nicht`, `semantic-change-body-only-edges-stale` — **die dritte erscheint in diesem Lauf wieder, und zwar am Kern.**

**Gefahrene Läufe.**
`make gates` → Exit 0; `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` (zehn Gates); `d-check: 680 Datei(en) geprüft, 0 Befund(e)`; `coverage-gate: OK — Coverage 94.60% erfüllt Schwelle 93%`.
`make mention-coverage` → Exit 0; `d-check: mentions: 84 von 84 Artefakt(en) erwähnt, über 1 Dokument(e)`.
Alle vier Zahlen der Botschaft — `zehn Gates`, `680 Dateien`, `94,60 %`, `84 von 84` — stehen wörtlich so im Lauf.

**Sechs eigene Bruch-Tests am Produktivcode**, je mit Wiederherstellung; sie beantworten, ob die neuen Tests fallen **können**.

| # | Eingriff in `mentions.go` | Erwartung | Ergebnis |
|---|---|---|---|
| M1 | `mentionsLeftFree`: `/` wieder blockieren | rot | **rot** — `TestMentionsPathTraegtRelativeVerlinkung`, `TestMentionsPfadPraefixDecktMitglied`, `TestMentionsBasenameKeineTeilzeichenkette` |
| M2 | `mentionsLeftFree` byte-basiert (`before[len-1]`) | rot | **rot** — `TestMentionsGrenzeIstRunenBasiert` |
| M4 | `mentionsRightFree` wieder weit (`-`, `.`, `_` blockieren) | rot | **rot** — `TestMentionsGrenzeMeldetHaeufigeFormenNicht`, mit **je Form benannter** Zeile |
| M5 | `mentionsRightFree` **byte-basiert** (`rune(after[0])`) | rot | **grün, Exit 0** — kein Test fällt. Siehe R3-M-1 |
| M6 / M6b | `ignored` **oder** `dirIgnored` einzeln aus `mentionsWalk` entfernt | rot | **grün, Exit 0** — je einzeln fällt nichts. Siehe R3-L-2 |
| M6c | **beide** entfernt | rot | **rot** — `TestMentionsHonoriertScanIgnore`. R2-N-3 ist damit geschlossen |

**Über sechzig Mess-Läufe gegen eigene Fixtures und einen Fremd-Bestand** (Container gegen einen eigenen Mount, `--network none`, read-only), dazu zwei aus dem Quellstand gebaute Vergleichs-Images (`d-check:old` = Grenze der Fassung 0.86.3, `d-check:sym` = symmetrisch-laxe Gegenprobe). Die vollständige Form-Tabelle steht unten.

**Arbeitsbaum.** Am Ende ist die einzige Änderung dieser Report (`git status --short` zeigt nur ihn); die Vergleichs-Images sind entfernt.

---

## R2-Finding → Status

| R2 | Kategorie | Status | Beleg |
|---|---|---|---|
| **N-1** Grenz-Prüfung meldet erwähnte Artefakte als unerwähnt | HIGH | **behoben — mit neuem Preis, der nicht vollständig nachgezogen ist** | Alle vier Formen zählen (Form-Tabelle Z. 2–5, 26). Gegenprobe `d-check:old` unter `path` gegen `harness/README.md`: `0 von 84` → jetzt `21 von 84`. **Offen:** R3-H-1, R3-M-4 |
| **N-2** Review-Befund-Marker in `.d-check.yml` | HIGH | **behoben** | Der `mentions:`-Block trägt keinen Marker und keine Chronik der verworfenen Fassung mehr (`.d-check.yml:882–898` gelesen). Rest-Bestand außerhalb des Blocks: R3-I-1 |
| **N-3** `scan.ignore`-Test kann nicht fallen | MEDIUM | **behoben** | M6c: nach Entfernen **beider** Auswertungen fällt `TestMentionsHonoriertScanIgnore`. Residuum je einzelner Mechanismus: R3-L-2 |
| **N-4** Grenze byte- statt runen-basiert | MEDIUM | **behoben im Verhalten** | `aÄx.md` ⇒ MELDET, `x.mdÄ` ⇒ MELDET (Form-Tabelle Z. 20, 27); Lastenheft, Spezifikation und Code stimmen jetzt zeichengenau überein. **Ungewächtert auf der rechten Seite:** R3-M-1 |
| **N-5** Historie 0.86.3 behauptet Beleg als „nachgeholt" | MEDIUM | **behoben, und die Korrektur stimmt** | Neue Fassung: „erstmals mit dem Modul selbst gefahren; als Messung führte ihn die Zeile 0.86.1 bereits". Nachgeprüft: 0.86.1 stammt aus `08fb419` (slice-205), das per `git merge-base --is-ancestor` **vor** `ac12993` (Geburt des Moduls) liegt — beide Halbsätze tragen |
| **N-6** Erkennungsform-Messung nicht nachgezogen, Stichprobe unbenannt | MEDIUM | **halb behoben** | Erste Hälfte eingelöst (Nachmessung behauptet). Zweite Hälfte **offen**: `spec/lastenheft.md:3480` nennt weiterhin nicht, welche zwölf Artefakte aus welchen drei Repos. Siehe R3-L-3 |
| **N-7** Drei Vertrags-Aussagen ohne Akzeptanzkriterium | MEDIUM | **behoben** | Nachgezählt über beide Revisionen: 10 → **15** Kriterien. Eines davon misst aber die falsche Sache: R3-M-3 |
| **N-8** Handbuch §5 vs. `matchAnyGlob`-Kommentar | MEDIUM | **behoben in der Richtung, neu falsch im Beispiel** | Die Ausnahme nennt genau die richtigen Felder (Inventur unten). Das gegebene Beispiel wirkt aber nicht: R3-M-2 |
| **N-9** Sensors-Tabellen-Invariante ohne Träger | LOW | **offen** | `git log -- docs/plan/planning/observations/` endet bei slice-205; `open/` enthält nur slice-207/208 (Baseline-Sprung). Siehe R3-L-4 |
| **N-10** fail-closed-Fläche = ganzes Repo | LOW | **offen** | vom Diff nicht berührt |
| **N-11** `scan.ignore` leert die Ist-Menge, Meldung zeigt auf `documents` | LOW | **offen** | vom Diff nicht berührt |
| **N-12** Test-Zahlen der Botschaft reproduzieren nicht | LOW | **offen (Botschaft ist historisch)** | betrifft `e5cf62d`; die Botschaft von `0f13f0a` nennt keine Test-Zahlen und ist insoweit sauber |
| **I-1** `seen`-Map strukturell tot | INFO | **offen** | `mentions.go:112–118` unverändert |
| **I-2** `Config.Mentions` ohne Feld-Kommentar | INFO | **offen** | unverändert |
| **I-3** Handbuch-Referenzblock ohne `mentions:` | INFO | **offen** | `grep -n "^mentions:" docs/user/benutzerhandbuch.md` ⇒ kein Treffer |
| **I-4** Paralleler Schreibzugriff (slice-207/208) | INFO | erledigt | als `af29684` committet, gehört nicht zum Gegenstand |

---

## Nennungsform → zählt / zählt nicht / erwartet

Gemessen mit dem ausgelieferten Modul gegen eigene Fixtures. Soll-Menge ein Mitglied; `basename`-Zeilen mit Needle `x.md`, `path`-Zeilen mit Needle `a/x.md`. **„erwartet"** ist das Urteil, das die Anforderung für diese Form zusagt oder das ein Leser der Anforderung ableiten würde.

| # | Nennungsform im Korpus | Ist | erwartet | Bewertung |
|---|---|---|---|---|
| 1 | `siehe x.md hier` | ZÄHLT | zählt | ✔ |
| 2 | `siehe x.md.` (Satz-Schlusspunkt) | ZÄHLT | zählt | ✔ N-1 eingelöst |
| 3 | `die x.md-Datei` (dt. Kompositum) | ZÄHLT | zählt | ✔ N-1 eingelöst |
| 4 | `_x.md_` (Kursivierung) | ZÄHLT | zählt | ✔ N-1 eingelöst |
| 5 | `[A](../a/x.md)` unter `path` | ZÄHLT | zählt | ✔ N-1 eingelöst |
| 6 | `(x.md)` · `**x.md**` · `~~x.md~~` · `` `x.md` `` | ZÄHLT | zählt | ✔ |
| 7 | `„x.md"` · `»x.md«` | ZÄHLT | zählt | ✔ |
| 8 | `x.md:12` · `x.md,` · `x.md;` · `x.md?` · `x.md's` | ZÄHLT | zählt | ✔ |
| 9 | `<!-- x.md -->` (HTML-Kommentar) | ZÄHLT | zählt | ✔ |
| 10 | `[^1]: x.md` (Fußnote) · `\| x.md \|` (Tabellenzelle) | ZÄHLT | zählt | ✔ |
| 11 | `\x.md` (Backslash davor) · `* x.md` · `\|-- x.md` (Baum) | ZÄHLT | zählt | ✔ |
| 12 | `x.md` am **Korpus-Anfang** / am **Korpus-Ende** ohne `\n` | ZÄHLT | zählt | ✔ |
| 13 | `[A](x.md#abschnitt)` (Anker) · `https://ex.com/a/x.md` | ZÄHLT | zählt | ✔ |
| 14 | `[A](./a/x.md)` · `/a/x.md` unter `path` | ZÄHLT | zählt | ✔ |
| 15 | `ax.md` (Buchstabe links) | MELDET | meldet | ✔ |
| 16 | `1x.md` (Ziffer links) | MELDET | meldet | ✔ |
| 17 | `a-x.md` (Bindestrich links) | MELDET | meldet | ✔ R1-H-2 geschlossen |
| 18 | `a.x.md` (Punkt links) | MELDET | meldet | ✔ trägt `.harness/…` ≠ `harness/…` |
| 19 | `x.mdx` · `x.md2` (Buchstabe/Ziffer rechts) | MELDET | meldet | ✔ |
| 20 | `aÄx.md` (Nicht-ASCII-Buchstabe links) | MELDET | meldet | ✔ N-4 geschlossen |
| 21 | `a_x.md` (Unterstrich links) | ZÄHLT | *meldet* | ⚠ **benannte** Grenze („Unterstrich-Kompositum") |
| 22 | `x.md.bak` (Punkt rechts) | ZÄHLT | *meldet* | ⚠ **benannte** Grenze („angehängte Endung") |
| 23 | `x/a/x.md` unter `path` (fremdes Präfix) | ZÄHLT | *meldet* | ⚠ **benannte** Grenze („fremdes Pfad-Präfix") |
| 24 | `x.md_alt` (Unterstrich **rechts**) | ZÄHLT | *meldet* | ✖ **unbenannt** — R3-M-4 |
| 25 | `x.md-alt` (Bindestrich **rechts**) | ZÄHLT | *meldet* | ✖ **unbenannt** — R3-M-4 |
| 26 | `x.md~` · `x.md+2` · `x.md/y` (rechts) | ZÄHLT | *meldet* | ✖ **unbenannt** — R3-M-4 |
| 27 | `x.mdÄ` (Nicht-ASCII rechts) | MELDET | meldet | ✔ (aber ungewächtert — R3-M-1) |
| 28 | `a+x.md` · `a~x.md` · `a@x.md` · `a%x.md` (links) | ZÄHLT | *meldet* | ✖ **unbenannt** — R3-M-4 |
| 29 | `a,x.md` · `a=x.md` · `a#x.md` · `a$x.md` · `a\x.md` (links) | ZÄHLT | *meldet* | ✖ **unbenannt** — R3-M-4 |
| 30 | `a x.md` (Leerzeichen links; Leerzeichen ist ein zulässiges Datei-Zeichen) | ZÄHLT | *meldet* | ✖ prinzipiell unentscheidbar, aber unbenannt |
| 31 | `ADR-x.md` (Präfix-Kompositum) | MELDET | *zählt* | ✖ **Falsch-Negativ, unbenannt** — R3-I-3 |
| 32 | `...x.md` (ASCII-Ellipse) | MELDET | *zählt* | ✖ Falsch-Negativ; `…x.md` (Unicode) ZÄHLT — zwei Schreibungen, zwei Urteile |
| 33 | `[B](mein%20bericht.md)` (URL-kodiert) | MELDET | *zählt* | ✖ Falsch-Negativ; von „wörtlich, ohne Markdown-Lexik" gedeckt, aber unbenannt |
| 34 | Zeilenumbruch **innerhalb** der Nennung (`../a/\nx.md`) | MELDET | *zählt* | ✖ `strings.Contains`-Grenze, unbenannt |
| 35 | `-x.md` (Liste ohne Leerzeichen) · `v1.x.md` | MELDET | Grenzfall | — |
| 36 | Nennung **über zwei Dokumente** hinweg (Korpus-Fuge) | MELDET | meldet | ✔ `\n`-Trenner in `mentionsCorpus` verhindert den Fugen-Treffer |

---

## Findings

### HIGH

**R3-H-1 · Die dokumentierte Vorbedingung für `basename` deckt die neue Kollisions-Menge nicht mehr — eine als sicher zertifizierte Konfiguration wird still grün**
`quelle`: [`DC-FA-MENT-001`](../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) · Reviewer-Prüffrage 1 (stilles Grün) und 2 (falsche Menge) · `klasse`: `semantic-change-body-only-edges-stale`
`pfad`: `.d-check.yml:896-898`, `docs/user/benutzerhandbuch.md:2342` (Modul-Tabelle `mentions`), `spec/lastenheft.md` §Historie 0.86.3 Punkt (4) — gegen `internal/hexagon/core/rules/mentions.go:249-257`
`verifizierbar`: **ja** — zwei Fixture-Läufe, unten wiedergegeben.

Anforderung, Handbuch und Konfiguration sagen einem Adopter übereinstimmend, wie er die Wahl `match: basename` absichert: *„Vor der Wahl sind zwei Dinge zu prüfen, nicht eines: die Basisnamen müssen eindeutig sein UND keiner darf **Endstück** eines anderen sein"*. Unter der Grenze der Fassung 0.86.3 war das die **vollständige** Kollisionsbedingung: rechts blockierte `.`, `-`, `_` und `/`, eine Fundstelle konnte also nur am Namensende enden. Dieser Commit öffnet die rechte Seite für `.`, `-`, `_`, `~`, `+` und `/` — und lässt alle drei Fundstellen der Vorbedingung unverändert stehen. „Endstück" ist damit die falsche Prüfung, und der Adopter, der sie korrekt ausführt, bekommt ein grünes Gate über nicht erwähnte Mitglieder.

Gemessen, `match: basename`, Soll-Menge `a/*`, Ist-Menge ein Dokument, das **ausschließlich** die längeren Namen nennt:

```
Soll-Menge:   a/a.md, a/a.md.bak, a/b.sh, a/b.sh.in
Ist-Dokument: "nur [X](a.md.bak) und [Y](b.sh.in) sind hier genannt"
⇒ d-check: mentions: 4 von 4 Artefakt(en) erwähnt, über 1 Dokument(e)
⇒ d-check: 1 Datei(en) geprüft, 0 Befund(e)                        Exit 0
```

`a.md` und `b.sh` kommen im Korpus nirgends vor. Beide bestehen die dokumentierte Vorbedingung: die vier Basisnamen sind eindeutig, und `a.md` ist **kein** Endstück von `a.md.bak` (`b.sh` keines von `b.sh.in`). Derselbe Befund unter dem **Default** `path`, wo die Anforderung überhaupt keine Vorbedingung nennt:

```
Soll-Menge:   a/a.md, a/a.md.bak     Ist-Dokument: "nur a/a.md.bak ist genannt"
⇒ d-check: mentions: 2 von 2 Artefakt(en) erwähnt, 0 Befund(e)     Exit 0
```

Die Paare `x.md`/`x.md.bak`, `x.yml`/`x.yml.example`, `x.sh`/`x.sh.in` sind gewöhnliche Bestände. Versagensbild: Ein Adopter führt die zugesagte Zwei-Punkte-Prüfung durch, aktiviert `mentions`, bekommt Exit 0 mit `N von N` und schließt daraus, dass jedes Artefakt genannt ist — während die Nennung des längeren Geschwisters das kürzere deckt. Das ist exakt die Klasse, die R1-H-2 als HIGH führte, um 180° gedreht: dort deckte ein längerer Name **links**, hier deckt er **rechts**, und diesmal steht die entwertete Vorbedingung noch dreifach im Vertrag. **Am eigenen Bestand greift es nicht** — die 84 ADR-Basisnamen haben unter alter wie neuer Regel **null** Kollisionen (nachgezählt) —, deshalb bleibt `make mention-coverage` grün und der Fehler unsichtbar.

### MEDIUM

**R3-M-1 · Die Runen-Basiertheit der **rechten** Grenze ist von keinem Test gehalten — eine byte-basierte Regression läuft grün durch**
`quelle`: [`DC-FA-MENT-001`](../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) · Reviewer-Anker „fehlende Negativtests bei neuem öffentlichen Vertrag" · `klasse`: `wortlaut-behauptet-pruefung-die-fehlt`
`pfad`: `internal/hexagon/core/rules/mentions_test.go` (`TestMentionsGrenzeIstRunenBasiert`) gegen `internal/hexagon/core/rules/mentions.go:251-257`
`verifizierbar`: **ja** — Bruch-Test M5.

`spec/spezifikation.md` Schritt 5 sagt zu: *„Die Grenze ist **asymmetrisch** und wird auf **Runen** geprüft, nicht auf Bytes"* — die Aussage gilt der Grenze als ganzer, und `mentionsRightFree` setzt sie mit `utf8.DecodeRuneInString` auch um. Der einzige Test dazu prüft nur die linke Hälfte: `TestMentionsGrenzeIstRunenBasiert` stellt `Ätest.md` **vor** die Fundstelle. Gemessen: Ersetzt man in `mentionsRightFree` die Runen-Dekodierung durch `rune(after[0])`, bleibt `make test` **grün, Exit 0** — kein Test fällt. Genau diese Regression stellte das stille Grün wieder her, dessen Fortbestehen die Botschaft als Grund für den Wechsel nennt (*„das stille Gruen bestand also bei jedem Nicht-ASCII-Nachbarn fort"*), nur eben rechts statt links; erreichbar wird sie für jede Schrift, deren UTF-8-Führungsbyte kein Buchstabe ist (`0xD7`, Hebräisch). Das ist dieselbe Klasse, die R2-N-3 für den `scan.ignore`-Test führte — Verhalten in Ordnung, Zusage ungewächtert —, hier am Kern-Vertrag statt an einem Ventil.

**R3-M-2 · Die neue Glob-Ausnahme im Handbuch gibt ein Beispiel, das nichts ändert, und nennt einen Fall, den es nicht gibt**
`quelle`: [`DC-FA-CONF-001`](../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei) · `klasse`: `semantic-change-body-only-edges-stale`
`pfad`: `docs/user/benutzerhandbuch.md:1218-1222` gegen `internal/hexagon/core/rules/workflows.go:235-256` und `internal/hexagon/core/rules/reviews.go:129-146`
`verifizierbar`: **ja** — Fixture-Lauf unten.

Der N-8-Fix schreibt: *„Wer dort einen **Teilbaum** ausnehmen will, schreibt die Segmente aus (`a/*/*.yml` statt `a/**/*.yml`)."* Beides trifft nicht.

**Erstens** sind die beiden genannten Formen unter blankem `path.Match` **äquivalent**: `**` ist dort zwei aufeinanderfolgende `*` innerhalb *eines* Segments und matcht dieselbe Menge wie ein einzelnes `*`. Gemessen mit `workflows.dir: ci/a/b/workflows` (Zielpfad `ci/a/b/workflows/x.yml`, fünf Segmente):

```
exempt-paths: ["ci/**/*.yml"]          ⇒ NICHT ausgenommen — Befund uses-pin-untagged
exempt-paths: ["ci/*/*/*/*.yml"]       ⇒ ausgenommen (leere Prüfmenge, fail-closed)
exempt-paths: ["ci/**/**/**/*.yml"]    ⇒ ausgenommen (leere Prüfmenge, fail-closed)
```

Wirksam ist die **Zahl der Segmente**, nicht der Ersatz von `**` durch `*`; die dritte Zeile zeigt, dass `**` mit passender Segmentzahl genauso trägt. Ein Adopter, der die Anweisung befolgt und `a/**/*.yml` durch `a/*/*.yml` ersetzt, ändert nichts.

**Zweitens** gibt es in beiden Modulen keinen Teilbaum: `workflowCandidates` und die Kandidaten-Sammlung von `reviews` rufen `fsys.List(dir)` und steigen **nicht ab** (das Handbuch sagt für `reviews` an anderer Stelle selbst „nicht rekursiv"). Die Pfade, gegen die `exempt-paths` dort je geprüft wird, sind ausschließlich `<dir>/<name>`. Versagensbild: Der Adopter schreibt eine Ausnahme nach Anleitung, sie greift nicht, und das Modul meldet weiter — laut, aber die Anleitung führt ihn nicht zur Ursache.

**R3-M-3 · Das neue Akzeptanzkriterium für die relative Verlinkung führt den Preis vor, nicht die Eigenschaft — und friert ihn als Vertrag ein**
`quelle`: [`DC-FA-MENT-001`](../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) · Reviewer-Prüffrage 10 (Messmethode klafft gegen die Spec-Stelle) · `klasse`: `zaehlmethode-misst-proxy-statt-gegenstand`
`pfad`: `spec/lastenheft.md`, Akzeptanzkriterium *„Eigenständige Nennung (was zählen muss)"*
`verifizierbar`: **ja** — Wortlaut gegen die Grenz-Liste derselben Anforderung.

Das Kriterium lautet: *„Given ein Mitglied `x.md`, when ein Dokument es … oder — unter `path` — als `../a/x.md` nennt, then **kein** Befund"*. Unter `path` ist der Needle der Pfad ab der Scan-Wurzel; ein Mitglied namens `x.md` liegt also **in der Wurzel**, und `../a/x.md` bezeichnet dann eine andere Datei. Was das Kriterium prüft, ist damit nicht die Eigenschaft *„eine relative Verlinkung auf das Mitglied zählt"*, sondern genau die Grenze, die dieselbe Anforderung drei Absätze höher als Preis ausweist: *„ein fremdes Pfad-Präfix deckt das Mitglied darunter"*. Der Unit-Test daneben (`TestMentionsPathTraegtRelativeVerlinkung`) macht es richtig — dort ist das Mitglied `docs/plan/a.md` und der Link `../docs/plan/a.md`. Folge: Eine spätere Implementierung, die den Preis verkleinerte (etwa indem sie links nur `./`, `../` oder den eigenen Verzeichnis-Anfang zuließe), müsste dieses Akzeptanzkriterium **brechen**, um die Anforderung besser zu erfüllen. Das Kriterium macht die benannte Grenze zur zugesagten Eigenschaft und verriegelt sie.

**R3-M-4 · „Drei benannte Grenzen" ist eine geschlossene Zahl über einer offenen Menge**
`quelle`: [`AGENTS.md`](../../AGENTS.md) §5 (nicht weiter schließen als die gemessene Menge) · `klasse`: `selbstauskunft-zahl-reproduziert-nicht`
`pfad`: `spec/lastenheft.md` §Was „kommt vor" heißt (*„**Und drei benannte Grenzen sind der Preis dafür**"*), `spec/spezifikation.md` Schritt 5, `harness/sensors/mention-coverage.md:37-40`, `docs/user/benutzerhandbuch.md:2342`, `internal/hexagon/core/rules/mentions.go:219-223`
`verifizierbar`: **ja** — Form-Tabelle Z. 24–26, 28–30.

Die vier Dokumente und der Modul-Kommentar nennen den Preis wortgleich als **drei** Fälle: Unterstrich-Kompositum, angehängte Endung, fremdes Pfad-Präfix. Die Regel lässt links jede Rune durch, die weder Buchstabe noch Ziffer noch `-` noch `.` ist, und rechts jede, die weder Buchstabe noch Ziffer ist. Alle diese Zeichen sind in Dateinamen zulässig, und mehrere sind gebräuchlich. Gemessen decken zusätzlich: `x.md_alt`, `x.md-alt`, `x.md~`, `x.md+2` (rechte Seite — die drei benannten Fälle nennen **nur** die linke `_`-Variante und die rechte `.`-Variante) sowie `a+x.md`, `a~x.md`, `a@x.md`, `a%x.md`, `a,x.md`, `a=x.md`, `a#x.md`, `a$x.md`, `a\x.md` und `a x.md` (linke Seite). Die begleitende Prosa (*„eine rein textuelle Regel kann Namensteil und Satzzeichen nicht trennen"*) ist richtig und trüge die offene Formulierung — die Zahl **drei** daneben behauptet eine Vollzähligkeit, die die Messung nicht stützt. Dieselbe Zählung ist zugleich die Wurzel von R3-H-1: Wer den Preis für dreifach hält, hält die alte „Endstück"-Vorbedingung für ausreichend.

### LOW

**R3-L-1 · Eine Zeile ist beim Editieren im Lastenheft verdoppelt worden**
`pfad`: `spec/lastenheft.md:3456-3457` · `quelle`: [`DC-FA-MENT-001`](../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) · `klasse`: `semantic-change-body-only-edges-stale`
Der Satz *„Geprüft wird das Vorkommen einer Zeichenkette im Dokument-Text — und zwar als"* steht zweimal untereinander; der Diff behält die Kontextzeile und fügt sie zusätzlich ein. Betroffen ist der Absatz, der die geänderte Semantik trägt, in der Quelle mit Rang 1. Kein Gate sieht es (`make doc-check` grün). Dieselbe Editier-Klasse hat die Historie-Zeile 0.86.2 unter Punkt (3) schon einmal für diese Anforderung festgehalten — zweites Auftreten. `verifizierbar`: ja (`grep -n "Geprüft wird das Vorkommen einer Zeichenkette" spec/lastenheft.md` ⇒ zwei Treffer).

**R3-L-2 · `TestMentionsHonoriertScanIgnore` hält nur die Konjunktion beider Ignorier-Mechanismen**
`pfad`: `internal/hexagon/core/rules/mentions_test.go` gegen `internal/hexagon/core/rules/mentions.go:150-162` · `klasse`: `wortlaut-behauptet-pruefung-die-fehlt`
R2-N-3 ist geschlossen (M6c fällt), aber die Fixture wählt mit `fremdbaum/**` ein Muster, das **beide** Mechanismen treffen: `dirIgnored` prunt das Verzeichnis, `ignored` schlösse die Datei ohnehin aus. Gemessen: Entfernt man `ignored` allein (M6) oder `dirIgnored` allein (M6b), bleibt `make test` grün. Die Botschaft sagt *„gegengeprueft durch Entfernen der Auswertung"* — für keine der beiden einzeln trifft das zu. Load-bearing ist die Verzeichnis-Hälfte bei Mustern **ohne** `/**` (`scan.ignore: ["vendor"]` prunt das Verzeichnis, matcht aber keine Datei darunter); genau diese Klasse bleibt ungewächtert. `verifizierbar`: ja.

**R3-L-3 · Die Erkennungsform-Stichprobe ist weiterhin nicht benannt, die Nachmessung damit nicht nachvollziehbar**
`pfad`: `spec/lastenheft.md:3480` · `quelle`: [`AGENTS.md`](../../AGENTS.md) §5 · `klasse`: `selbstauskunft-zahl-reproduziert-nicht`
Die Botschaft schließt: *„die Erkennungsform-Messung ist unter der neuen Semantik nachgefahren … zwoelf Artefakte aus drei Repos, EINE Abweichung — die Aussage haelt."* Welche zwölf Artefakte aus welchen drei Repos, steht weiterhin nirgends — das war die zweite Hälfte von R2-N-6 und ist unverändert. Die *andere* Fremd-Messung derselben Anforderung ist benannt genug und reproduziert unter der neuen Semantik exakt: `scripts/*.sh` + `tools/*.py` gegen `docs/user/*.md` (u-boot) ⇒ **5 Mitglieder, 8 Dokumente, 4 Funde, `path` und `basename` identisch**. Für die Zwölfer-Stichprobe ist weder Bestätigung noch Widerlegung möglich. `verifizierbar`: ja (die Nicht-Reproduzierbarkeit ist der Befund).

**R3-L-4 · Die als „nicht gehalten" deklarierte Sensors-Tabellen-Invariante hat weiterhin keinen Träger**
`pfad`: `harness/sensors/mention-coverage.md:27-28` · `quelle`: Baseline `modul-05-planning-harness.md` §Offene Risiken werden bei Closure aufgelöst · `klasse`: `registerzeile-ohne-ausgang-nach-schwelle`
R2-N-9 unverändert: `git log -- docs/plan/planning/observations/` endet bei slice-205, `next/` ist leer, `open/` trägt nur slice-207 und slice-208 (Baseline-Sprung). Die Sensor-Datei sagt weiterhin *„sie bräuchte einen zweiten Lauf"*, ohne dass ein Vorgang das einlöst. Der Slice steht noch in `in-progress/`; einlösbar bleibt es bis zur Closure. `verifizierbar`: ja.

### INFO

- **R3-I-1 · Die Klasse aus N-2 lebt außerhalb des korrigierten Blocks weiter.** `.d-check.yml` trägt an drei Stellen Chronik statt Zustand: `:205` (*„done/slice-200 und done/slice-201 je zwei d-check:cite-Direktiven"*), `:647` (*„sie vor slice-203 verloren hatten"*), `:850` (*„Die fuenf Eintraege sind mit slice-200 entfernt"*). Alle drei liegen außerhalb des `mentions:`-Blocks und sind vom Commit nicht berührt — unter der §3.7-Bestandsgrenze („geräumt beim nächsten Anfassen der Zeile") also **kein** Verstoß. Notiert, weil die Botschaft die Einsicht als *„§3.7 gilt … für die Config"* formuliert und der Fix nur den Block trägt; die nächste Lesart könnte die Datei für geräumt halten.
- **R3-I-2 · Kein einziger Falsch-Negativ-Fall ist irgendwo benannt.** Die drei „benannten Grenzen" liegen alle auf der Falsch-Positiv-Seite. Messbar sind mindestens drei auf der anderen: das Präfix-Kompositum (`ADR-x.md` deckt `x.md` nicht — Spiegelbild des zugesagten `x.md-Datei`), der URL-kodierte Link-Ziel (`mein%20bericht.md`) und der Zeilenumbruch innerhalb einer Nennung. Die zweite und dritte sind von *„wörtlich im rohen Text — ohne Markdown-Lexik"* implizit gedeckt, die erste von nichts. Zusätzlich urteilt die Regel über zwei Schreibungen derselben Interpunktion verschieden: `…x.md` (U+2026) zählt, `...x.md` nicht.
- **R3-I-3 · R2-I-1 bis I-3 unverändert**: die `seen`-Map in `mentionsFilter` bleibt strukturell tot, `Config.Mentions` ohne Feld-Kommentar, der Konfigurations-Referenzblock des Handbuchs ohne `mentions:`-Eintrag (`grep -n "^mentions:"` ⇒ kein Treffer).
- **R3-I-4 · Formatierungs-Drift, ausdrücklich nicht als Finding gemeldet.** Das Map-Literal in `TestMentionsHonoriertScanIgnore` (`mentions_test.go:308-312`) ist nicht am längsten Schlüssel ausgerichtet, wie es das Standard-Format täte. `.golangci.yml` führt **keinen** Formatierer (24 Linter, kein `formatters:`-Block), und im Repo gibt es keinen Konventions-Anker dafür — nach dem Anti-Pattern „Kein Stil-Polizist" ist das kein Finding. Notiert, damit der nächste Reviewer die Ableitung nicht wiederholt.

---

## Negativbefunde

- **Die Asymmetrie selbst trägt, und das ist gemessen statt argumentiert.** Beide Gegenrichtungen wurden als eigenes Image gebaut und gefahren. *Symmetrisch streng* (rechts blockiert wieder `-`, `.`, `_`, Bruch-Test M4): die Formen Satz-Schlusspunkt, Kompositum und Kursivierung fallen durch — die Fassung 0.86.3, deren Fehlerpolitik R2-N-1 beanstandete. *Symmetrisch lax* (`d-check:sym`, links blockiert nur noch Buchstabe/Ziffer): `a-x.md` und `a.x.md` decken wieder `x.md` — die Rückkehr von R1-H-2 —, **ohne** dass eine der vier reparierten Formen gewinnt (sie zählen in beiden Fassungen). Die gewählte Asymmetrie ist damit die einzige der drei Positionen, die beide Randbedingungen erfüllt, und sie ist sachlich begründet: Pfad- und Namens-Qualifikatoren stehen **links**, die deutsche Kompositum-Fuge steht **rechts**.
- **Der Kern-Fix löst N-1 wirklich ein.** Alle vier gemeldeten Formen zählen jetzt (Form-Tabelle Z. 2–5). Die Zahl „21 von 84" reproduziert exakt: mit `match: path`, Soll-Menge `docs/plan/adr/[0-9]*.md`, Ist-Menge `harness/README.md` liefert das aus dem Quellstand gebaute Vorgänger-Image `0 von 84` und der ausgelieferte Stand `21 von 84`. Über alle Markdown-Dateien gemessen 66 → 69, über `harness/**/*.md` 0 → 36 — die Richtung des Befundes ist an drei unabhängigen Mengen bestätigt.
- **Alle drei neuen Tests können fallen, und der Map-Test meldet die richtige Form.** M1 kippt `TestMentionsPathTraegtRelativeVerlinkung` und `TestMentionsPfadPraefixDecktMitglied`, M2 kippt `TestMentionsGrenzeIstRunenBasiert`, M4 kippt `TestMentionsGrenzeMeldetHaeufigeFormenNicht`. Die Sorge um die zufällige Map-Iteration ist **widerlegt**: Der Test verwendet `t.Errorf` statt `t.Fatalf` und führt den Formnamen im Meldungstext, meldet also **alle** fehlgeschlagenen Formen, jede benannt — bei M4 wörtlich `deutsches Kompositum:`, `Kursivierung:`, `Satz-Schlusspunkt:`. Zufällig ist allein die Reihenfolge der Zeilen.
- **Der Korpus-Fugen-Treffer ist ausgeschlossen.** `mentionsCorpus` schreibt nach jedem Dokument ein `\n`; ein Needle kann eine Dokumentgrenze nicht überspannen. Gegengeprüft mit einer Fixture, deren erstes Dokument auf `x.` und deren zweites mit `md` beginnt — Befund, wie es sein muss.
- **Die eigene Soll-Menge ist unter beiden Regeln kollisionsfrei.** Alle 84 ADR-Basisnamen paarweise geprüft: **0** Kollisionen unter der neuen Grenze, **0** unter der alten Endstück-Regel, 84 eindeutig. `make mention-coverage` ist deshalb von R3-H-1 nicht betroffen — und das Gate kann fallen: Entfernt man die ADR-0084-Zeile aus dem Index, meldet es `83 von 84` und einen `artifact-unmentioned` auf dem richtigen Pfad.
- **Die Runen-Grenze wirkt zur Laufzeit auf beiden Seiten.** `aÄx.md` ⇒ MELDET, `x.mdÄ` ⇒ MELDET. Beanstandet ist nur, dass die rechte Hälfte von keinem Test gehalten wird (R3-M-1).
- **Lastenheft, Spezifikation, Sensor-Datei, Handbuch-Modultabelle und Modul-Kommentar beschreiben die Regel zeichengenau gleich** — links Buchstabe/Ziffer/`-`/`.`, rechts Buchstabe/Ziffer, runen-basiert. Der `basename`-Sonderfall für `/` ist überall entfallen; die Regel ist für beide Erkennungsformen identisch, und keine der beiden ist dadurch strenger geworden. Unter `path` ist die linke Seite gelockert (das ist der Zweck und der benannte Preis), unter `basename` sind links `_` und rechts `_`, `-`, `.`, `~`, `+`, `/` hinzugekommen (das ist R3-M-4/H-1).
- **Die N-8-Inventur ist vollständig.** Ausgezählt über alle `exempt-paths`-Verwendungen: `codepaths`, `diagrams`, `versions`, `ids`, `matrix` und `structure` gehen über `ignored`/`matchGlob` (volle `**`-Semantik), ausschließlich `workflows` und `reviews` über `matchAnyGlob`/blankes `path.Match`. Die Ausnahme des Handbuchs nennt genau diese beiden — kein weiteres Feld dieser Familie ist der vollen Semantik fälschlich zugeschlagen. Beanstandet ist allein das Beispiel (R3-M-2). Die separat mit blankem `path.Match` arbeitenden Glob-Felder von `planning`/`planning_waves` matchen gegen einen **Basisnamen**, wo `**` gegenstandslos ist; der Handbuch-Satz wird dort nicht falsch.
- **Die N-5-Korrektur an der Historie-Zeile 0.86.3 ist chronologisch belegt.** 0.86.1 stammt aus `08fb419` (slice-205); `git merge-base --is-ancestor 08fb419 ac12993` bestätigt, dass die Zeile **vor** der Geburt des Moduls entstand. Beide Halbsätze der neuen Fassung — „als Messung führte ihn die Zeile 0.86.1 bereits" und „erstmals mit dem Modul selbst gefahren" — tragen.
- **N-7 reproduziert.** Die Akzeptanzkriterien-Liste von `DC-FA-MENT-001` zählt über beide Revisionen 10 → 15; die Botschaft sagt „fünf nachgetragen" und „die Liste stand unverändert bei zehn". Beides stimmt.
- **Die Gate-Zahlen der Botschaft sind echt.** `make gates` Exit 0 mit zehn Gates, `680 Datei(en) geprüft, 0 Befund(e)`, `Coverage 94.60%`; `make mention-coverage` Exit 0 mit `84 von 84`. Der Baum enthält 84 ADR-Dateien.
- **Der `mentions:`-Block der `.d-check.yml` ist frei von der beanstandeten Klasse.** Kein Review-Befund-Marker, keine Chronik der verworfenen Fassung; was bleibt, ist die Grenze („Ein Block ist EIN Paar", „Die Ist-Menge wird als VEREINIGUNG gelesen") und die Mengen-Begründung. Die Kommentare in `mentions.go` und `mentions_test.go` tragen durchweg eine der fünf Klassen (Zusage, Abgrenzung, Grenze); keine Slice-Nummer, kein Mess-Label, keine Review-Historie.
- **Der Hexagon-Schnitt ist unberührt**: `mentions.go` importiert weiterhin nur `model` und `port/driven` plus Standardbibliothek (neu `unicode`, `unicode/utf8`); `make arch-check` grün. Kein Netzzugriff, keine Inline-Suppression, keine gesenkte Schwelle im Diff.

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | **1** | R3-H-1 |
| MEDIUM | **4** | R3-M-1, R3-M-2, R3-M-3, R3-M-4 |
| LOW | **4** | R3-L-1, R3-L-2, R3-L-3, R3-L-4 |
| INFO | **4** | R3-I-1, R3-I-2, R3-I-3, R3-I-4 |

---

## Verdikt

**Blockierend** (ein HIGH, vier MEDIUM).

**Die dritte Fassung der Grenze trägt — als Regel.** Anders als ihre beiden Vorgänger ist sie an beiden Rändern gemessen: Die vier von R2-N-1 gemeldeten Formen zählen, die von R1-H-2 gemeldete Kollision bleibt geschlossen, und beide Gegenrichtungen sind als gebaute Images gefahren und fallen jeweils an einer der beiden Bedingungen. Die Asymmetrie ist nicht willkürlich, sondern folgt der Stelle, an der Namens-Qualifikatoren stehen (links) und an der die deutsche Kompositum-Fuge steht (rechts). Die Runen-Prüfung schließt die von R2-N-4 gemeldete Restklasse. **Die Fehlerpolitik ist diesmal nicht gekippt, sondern austariert.**

**Nicht getragen hat der Vertrag um die Regel herum.** Die Grenze wurde zum dritten Mal geändert und die Vorbedingung, die dem Adopter die Wahl `basename` absichert, zum dritten Mal nicht — sie stammt aus der Fassung 0.86.3, war dort vollständig und ist es jetzt nicht mehr. Das Ergebnis ist ein stilles Grün auf einer Konfiguration, die Anforderung, Handbuch und Beispiel-Konfiguration als sicher zertifizieren (R3-H-1). Dieselbe Wurzel trägt R3-M-4: Wer den Preis auf drei Fälle beziffert, kann die alte Vorbedingung für ausreichend halten. Die drei übrigen MEDIUM sind Nachzieh-Lücken derselben Art — eine Zusage, die kein Test hält (R3-M-1), ein Kriterium, das den Preis statt der Eigenschaft misst und ihn damit verriegelt (R3-M-3), und eine frisch geschriebene Handbuch-Passage, deren Beispiel unter der beschriebenen Semantik nichts bewirkt (R3-M-2).

Die Klasse `semantic-change-body-only-edges-stale` erscheint damit in R1, R2 und R3 — **drittes Auftreten in Folge an derselben Anforderung**. Nach dem Reviewer-Skill (§Kontext-Eskalation) ist das ein Steering-Loop-Signal: Es genügt nicht mehr, sie einzeln zu melden. [`MR-025`](../../harness/conventions.md#mr-025) verlangt bereits, bei einer Semantik-Änderung die Spiegel **vor** dem Editieren aufzulisten; für die Erkennungs-Form dieses Moduls hat das dreimal nicht gegriffen, weil ein Spiegel kein Verweis ist, sondern eine **abgeleitete Aussage** (die Vorbedingung „Endstück"). Das gehört ins Beobachtungs-Register und nicht in einen vierten Report.

**Nicht Gegenstand dieses Reports** ist die DoD-Abhakung von slice-206 — das prüft die Verifikation in getrenntem Kontext.
