# Verifikation slice-182 — die erklärte Leermenge bekommt eine Zahl

**Gegenstand:** `a062fe8^..HEAD` (4 Commits: Claim-Move `a062fe8`, feat `9c383ba`, Handbuch `307311c`, CR-Antwort `47b0ccb`).
**Rolle:** unabhängiger Verifier — geprüft wird gegen **DoD und Spec**, nicht gegen Plan/ADR-Maintainability.
**Gefahrene Sensors:** `make gates` (alle elf Glieder) · `make fullbuild` · `make completeness-check` · `make verify-closure-notes` · `make trace-check RANGE=…` · `make adr-check RANGE=…` · Produktläufe des lokal gebauten Images gegen **eigene Probe-Konfigurationen** · sechs Mutations- und eine Kontrollprobe gegen eine **Scratch-Kopie** des Baums (Build-Stage `test`) · ein aus `a062fe8` gebautes **Vorgänger-Image** für die Byte-Identität.
**Nicht verändert:** kein Repo-Artefakt außer dieser Datei. Alle Mutationen liefen gegen eine Kopie außerhalb des Repos.

---

## 1. DoD-Tabelle

| # | Behauptet (§4 des Slice-Plans) | Gemessen | Verdikt |
|---|---|---|---|
| 1 | `exempt-expect-count` im Schema, im `--print-config`-Gerüst, im Lastenheft (Bump + Historie mit CR-Bezug) und in der Spezifikation (Schritt 3, §2-Schema, §4-Grund-Code) geführt | Schema: `configyaml.go:243-246` (Zeiger-Typ) und `config.go:495-504`. `--print-config`: Zeile 199 der Ausgabe des Images trägt `exempt-expect-count: 19` samt Rand-Hinweis. Lastenheft: Version `0.77.0` → **`0.78.0`**, Historie-Zeile `:3364` nennt den Anlass als eingehenden CR (ohne Pfad — §3.4-konform), Grund-Code-Tabelle `:2532`, Prosa `:2569-2584`, zwei neue Akzeptanzkriterien `:2815-2816`, Config-Rand `:2779`. Spezifikation: Schritt 3 `:2165-2172`, §2-Schema `:2902`, §4 `SPEC-077` `:3033`, Historie `:3062` | **erfüllt** — mit **V-1**: derselbe Spezifikations-Edit hat die **Kopfzeile der §6-Tabelle zerstört** |
| 2 | Default byte-identisch, **gegen das Vorgänger-Image**, nicht gegen einen grünen Lauf | Vorgänger-Image aus `a062fe8` gebaut. **Vier Messungen, alle byte-identisch:** (a) Profil der ADR (`max-tasks: 3`, Selektor auf den DoD-Abschnitt, Glob über die Slice-Pläne) → beidseits **169 Befunde** (86 `section-missing` + 83 `section-oversized`), `diff` **leer** — die Zahl der ADR reproduziert **exakt**; (b) dieselbe Prüfmenge in `--json` (53 000 B), `--yaml` (44 343 B) und `--doctor` (46 125 B): stdout **und** stderr je `cmp`-gleich; (c) Repo-Konfiguration über den ganzen Baum: 618/**0** beidseits; (d) Repo-Konfiguration **ohne** `scan.ignore` über den getrackten Baum: 667 Dateien, **105 Befunde**, sortierter Befundsatz `diff`-leer | **erfüllt, stärker belegt als zugesagt** |
| 3 | Die deklarierte Leermenge ist stumm: N Abschnitte, N ausgenommen, `exempt-expect-count: N` ⇒ kein Befund, Exit 0. Mit Test | Produktlauf gegen eine eigene Probe (drei Abschnitte, alle vom Ventil getroffen, Zahl `3`): `1 Datei(en) geprüft, 0 Befund(e)`, **Exit 0**. Ohne den Schlüssel dieselbe Eingabe: `section-missing`, Exit 1. Test `TestExemptExpectCount_DeklarierteLeermengeIstStumm` prüft **beide Hälften in einer Funktion** und belegt damit den Vorzustand am selben Fixture | **erfüllt** |
| 4 | Die Drift ist beidseitig laut: weniger **und** mehr ausgenommen als deklariert ⇒ `section-exempt-mismatch`. Je ein Test | Produkt, Zahl `2` → `section-exempt-mismatch`, Meldung *„exempt-section-pattern nimmt 3 von 3 Abschnitten aus, deklariert sind 2"*, Exit 1. Zahl `4` → derselbe Code, *„… deklariert sind 4"*. Ein Ventil, das **weniger** nimmt, bei Zahl `3` → *„nimmt 2 von 3 … deklariert sind 3"*. Test `TestExemptExpectCount_DriftIstBeidseitig` mit zwei Unterfällen | **erfüllt** |
| 5 | Die Trennung hält: trifft `section-pattern` nichts, bleibt es `section-missing` — auch mit gesetztem Schlüssel. Mit Test | Produkt, Selektor ohne Treffer **plus** gesetzter Zahl `0`: `section-missing`, Meldung *„kein Abschnitt passt auf den Selektor"*, Exit 1 — der Konfigurationsdefekt bleibt einer. Test `TestExemptExpectCount_SelektorOhneTrefferBleibtDefekt` prüft zusätzlich den **Meldungstext**, nicht nur den Code | **erfüllt** |
| 6 | Zwei Config-Ränder: ohne `exempt-section-pattern` ⇒ Exit 2; Wert < 0 ⇒ Exit 2. Je ein Test | Produkt: ohne Muster → `structure[0]: exempt-expect-count ist ohne exempt-section-pattern wirkungslos (halbe Aktivierung)`, **Exit 2**. Wert `-1` → `structure[0]: exempt-expect-count -1 muss >= 0 sein`, **Exit 2**. Dritter, nicht zugesagter Rand mitgemessen: ein nicht-ganzzahliger Wert scheitert am strikten Decoding, ebenfalls Exit 2. Beide zugesagten Ränder als Fälle in `TestDecode_StructureFehler` | **erfüllt** |
| 7 | Die zwei offenen Fragen sind entschieden und begründet (Identität, raw-vs-hint) — im Code-Kommentar **und** in der ADR | Identität: Kommentar in `config.go:589-593` (*„zwei erwartete Zahlen … sind ein Widerspruch, kein Paar"*) und ADR §Entscheidung 8. raw-vs-hint: Kommentar in `structure.go:138-146` und ADR §Entscheidung 7. **Beide am Produkt gegengeprüft:** die Regel-Identität ist bei den Zahlen `2`, `3` und `4` **dieselbe** Zeichenkette, und zwei Regeln, die sich nur in der Zahl unterscheiden, weist der Config-Adapter als `Regel-Identität … kommt doppelt vor` mit **Exit 2** ab — genau die in der ADR als richtig bezeichnete Antwort. Der Mismatch trägt einen gesetzten `hint` (Meldung wird ersetzt), der Nullmengen-Befund **nicht** (behält seine modul-eigene Meldung) | **erfüllt** |
| 8 | Umkehr-Proben je Zusage, jede von den Tests gefangen, die dagegen stehen; je Regressions-Test der Beleg, dass der Vorzustand am Fixture scheitert | **Sechs Mutationen gegen eine Scratch-Kopie**, Kontrolllauf grün (§3). Alle **fünf** in der ADR §Fitness Function genannten Zahlen reproduzieren **exakt**: 5 · 3 · 6 · 1 · 1. Eine **sechste**, in der ADR nicht genannte Mutation (die Zahl geht doch in die Identität ein) fängt genau den siebten Test. Der Vorzustands-Beleg steht in Test 1 selbst | **erfüllt** |
| 9 | Eine ADR begründet Verortung, den neuen Grund-Code, die beidseitige Drift und die Abweichung von der beantragten Form; im ADR-Index eingetragen | ADR-0078 `Accepted`, acht nummerierte Entscheidungen, vier verglichene Alternativen, Fitness Function, drei Re-Evaluierungs-Trigger. Index-Zeile `docs/plan/adr/README.md:88`. `make adr-check RANGE=a062fe8^..HEAD` **Exit 0** | **erfüllt** — mit **V-2**: eine tragende Messung der ADR ist zur Hälfte falsch |
| 10 | Das Benutzerhandbuch führt den Schlüssel samt der Falle, die ihn nötig macht | `docs/user/benutzerhandbuch.md`: Schlüssel im Konfigurations-Fence (drei Kommentarzeilen) und ein eigener Absatz mit beiden Richtungen, der Null-Semantik, der Trennung **und** dem Preis (*„die Zahl altert wie jeder andere Autoren-Text, und wer sie blind mitzieht, hat einen Wächter, der nur noch sich selbst bestätigt"*) | **erfüllt** |
| 11 | Der Absender bekommt eine Antwort — angenommen in der Sache, abgelehnt in der Form, mit der `--doctor`-Messung | `docs/plan/cr/2026-08-30-antwort-a-check-leermenge.md`: drei benannte Abweichungen, Lage-Tabelle, Mess-Tabelle, Re-Evaluierungs-Angebot. **Bemerkenswert:** die Mess-Tabelle des ausgehenden Dokuments ist die **präziseste** Fassung im ganzen Bestand — sie sagt *„kein `doc-doctor` im Gate-Pfad"* und ist damit als einzige richtig (siehe **V-2**) | **erfüllt** |
| 12 | `make gates` und `make fullbuild` grün (Exit explizit); unabhängiger Review; Verifikation — beide in eigenen Kontexten | `make gates` → **Exit 0**, Schlusszeile nennt alle elf Glieder; Coverage **94,70 % ≥ 93 %**; semgrep 55 Regeln, 0 Befunde. `make fullbuild` → **Exit 0**, `[fullbuild] green — image-hash sha256:f59c1333ad9a…`. `make completeness-check` → **50 Anforderung(en), 0 Waise(n)**, Exit 0. `make verify-closure-notes` → 556/0, Exit 0. **Review-Report liegt noch nicht in `docs/reviews/`** (siehe **V-6**) | **hälftig offen (planmäßig)** — dies ist die Verifikation |

## 2. Akzeptanzkriterien von `DC-FA-STRUCT-001` (0.78.0) gegen das laufende Produkt

Alle Läufe: lokal gebautes Image, `--network none`, read-only-Mount, eigene Probe-Konfiguration und ein eigenes Probe-Dokument mit drei gleichartigen Abschnitten.

| Akzeptanzkriterium (Lastenheft) | Zusage | Gemessen | Verdikt |
|---|---|---|---|
| *Erklärte Leermenge* — drei Abschnitte, alle getroffen, Zahl `3` | kein Befund, Exit 0; **ohne** den Schlüssel `section-missing` | `0 Befund(e)`, **Exit 0**. Ohne Schlüssel: `section-missing`, Exit 1 | **erfüllt** |
| dieselbe Datei mit Zahl `2` **oder** `4` | je ein `section-exempt-mismatch`, „dessen Meldung **beide Zahlen** nennt" | Zahl `2` → *„nimmt **3** von 3 … deklariert sind **2**"*; Zahl `4` → *„nimmt **3** von 3 … deklariert sind **4**"*. Die Meldung nennt sogar **drei** Zahlen (ausgenommen, Grundmenge, deklariert) | **erfüllt, mehr als zugesagt** |
| vierter, **nicht** getroffener Abschnitt, weiterhin Zahl `3` | „nur dieser wird geprüft" | genau **ein** Befund, `section-marker-missing` auf Zeile 15 — dem vierten Abschnitt | **erfüllt** |
| derselbe Bestand bei Zahl `4` | Mismatch, „obwohl noch eine Restmenge da ist" | `section-exempt-mismatch`, *„nimmt 3 von **4** Abschnitten aus, deklariert sind 4"*, Zeile 1 | **erfüllt** |
| *Trennung* — `section-pattern` trifft nichts, Zahl gesetzt | `section-missing`, „ein Konfigurationsdefekt bleibt einer" | `section-missing`, *„kein Abschnitt passt auf den Selektor"* — **nicht** die Nullmengen-Meldung, also die richtige der beiden | **erfüllt** |
| *Null* — Muster nimmt nichts, Zahl `0` | „alle Abschnitte werden geprüft, kein Mismatch" | **drei** `section-marker-missing`, kein Mismatch | **erfüllt** |
| Schema `:2902` — Befund auf Zeile 1 | `line` = 1 | alle Mismatch-Befunde auf Zeile 1 | **erfüllt** |
| Schema `:2902` — „geht **nicht** in die Regel-Identität ein" | zwei Zahlen ⇒ Widerspruch, kein Paar | Identitäts-Zeichenkette bei den Zahlen `2`/`3`/`4` **identisch**; zwei Regeln, die sich nur in der Zahl unterscheiden ⇒ `Regel-Identität … kommt doppelt vor`, **Exit 2** | **erfüllt** |
| Spezifikation Schritt 3 — „nicht von `hint` ausgenommen" | ein verfasster Hinweis darf den Mismatch erklären | mit `hint` trägt der Mismatch den `hint`; der **Nullmengen**-Befund behält daneben seine modul-eigene Meldung — die Abgrenzung hält in beide Richtungen | **erfüllt** |
| Grund-Code-Tabelle `:2532` — Reparatur „Aufzählung oder Zahl nachziehen" | `--doctor`-Klartext | `--doctor`: *„Die Ausnahme nimmt eine andere Zahl von Abschnitten, als exempt-expect-count deklariert — Aufzählung oder Zahl nachziehen"* | **erfüllt** |
| Config-Ränder `:2779` | ohne Muster ⇒ Exit 2; Wert < 0 ⇒ Exit 2 | beide **Exit 2**, beide mit Klartext und Regel-Index | **erfüllt** |

**Kein Akzeptanzkriterium ist unerfüllt.** Alle elf Zusagen halten gegen das Binary, nicht nur gegen die Tests.

## 3. Mutations- und Kontrollproben

Gegen eine Scratch-Kopie des HEAD-Baums außerhalb des Repos, Build-Stage `test` in Docker. Gezählt sind eindeutige `--- FAIL:`-Einträge inklusive Unterfälle.

| Probe | Was geändert | Rot | Welche | ADR behauptet |
|---|---|---|---|---|
| **M0 — Kontrolle** | nichts | **0** | — | — |
| **M1** | Zählung deaktiviert (der `case` bleibt, der Mismatch feuert nie) | **5** | Drift (mit 2 Unterfällen), Hint, Restmenge | 5 — **trifft** |
| **M2** | Drift nur einseitig (Gleichheits- durch Kleiner-Vergleich ersetzt) | **3** | Drift (mit 1 Unterfall), Hint | 3 — **trifft** |
| **M3** | `case`-Guard entfernt, die Nullmengen-Härte feuert wieder | **6** | zusätzlich Leermenge-stumm und alle Unterfälle | 6 — **trifft** |
| **M4** | Mismatch als **raw** Finding | **1** | nur `TestExemptExpectCount_HintGiltFuerDenMismatch` | 1 — **trifft** |
| **M5** | beide Config-Ränder aus dem Adapter entfernt | **1** | `TestDecode_StructureFehler` | 1 — **trifft** |
| **M6** *(zusätzlich, in der ADR nicht genannt)* | die Zahl geht **doch** in die Regel-Identität ein | **1** | nur `TestExemptExpectCount_TraegtNichtDieRegelIdentitaet` | — |

**Die Fitness Function der ADR hält, was sie behauptet** — fünf Zahlen, fünf Treffer. M6 zeigt zusätzlich, dass auch die siebte Zusage (Identität) einen eigenen, spezifischen Wächter hat: keiner der übrigen sechs Tests fällt mit ihr. Die Aussage der ADR, der entfernte Frühausstieg für die geleerte Menge sei **toter Code** gewesen, ist mit M3 konsistent: der Zweig wirkt, weil er die Härte **überspringt**, nicht weil er unterdrückt.

## 4. Gate-Läufe (echte Ausgabe)

| Lauf | Ergebnis |
|---|---|
| `make gates` | **Exit 0** — `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` |
| darin `coverage-gate` | `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%` |
| darin `semgrep` | `Ran 55 rules on 55 files: 0 findings.` |
| darin `gate-consistency` / `planning-check` | je `618 Datei(en) geprüft, 0 Befund(e)` |
| `make fullbuild` | **Exit 0** — `[fullbuild] green — image-hash sha256:f59c1333ad9a…` |
| `make completeness-check` | **Exit 0** — `50 Anforderung(en), 0 Waise(n)`; die Matrix führt `DC-FA-STRUCT-001` mit `ADR-0078` und `slice-182` |
| `make verify-closure-notes` | **Exit 0** — `556 Datei(en) geprüft, 0 Befund(e)` |
| `make trace-check RANGE=a062fe8^..HEAD` | **Exit 0** — 618/0 |
| `make adr-check RANGE=a062fe8^..HEAD` | **Exit 0** — 618/0 |

**Deklarations-Konsistenz der Oberflächen.** Der neue Grund-Code steht in `AllReasons()` (`diagnose.go:94`), im `--doctor`-Klartext (`:159`), im `--print-config`-Gerüst (`config_template.go:207-210`), in der Lastenheft-Grund-Code-Tabelle (`:2532`) und als `SPEC-077` im Spezifikations-§4 (`:3033`). Die Verriegelung `TestAllReasonsDeckungGegenSpezifikationGrundCodes` hält die Deckung zwischen §4 und `AllReasons()` maschinell und ist grün. **Es fehlt keine Oberfläche** — mit der einen Einschränkung aus **V-3** (`CHANGELOG.md` und die beiden READMEs sind Release-Prep-Flächen, in diesem Repo planmäßig später).

## 5. Spec-Konformität und Abwärts-Sperre (§3.4)

**Kein verbotenes Token in den Spec-Deltas.** Über den gesamten Range gemessen: die hinzugefügten Zeilen in `spec/lastenheft.md` und `spec/spezifikation.md` enthalten **keine** vierstellige ADR-Kennung, **keine** Slice- oder Wellen-Kennung und **keinen** Commit-Hash. Beide Historie-Zeilen sagen *„Begründung in begleitender ADR"* ohne Kennung — die im Repo eingelebte Form. Das Modul `matrix` ist grün.

**Historie-Positionen.** Die Lastenheft-Zeile `0.78.0` steht korrekt an der Spitze der Versions-Tabelle. Die Spezifikations-Zeile steht korrekt an der Spitze der §7-Historie — **und ein zweites Mal an einer Stelle, an der sie nicht hingehört** (**V-1**).

**Lifecycle.** `a062fe8` ist ein **reiner Move**: `git diff --raw -M` meldet `R100`, Blob-Hash `8034621` vor und nach. Der gekoppelte Roadmap-Flip liegt im selben Commit (`MR-013`): der Ruhe-Marker verlässt den Abschnitt *Offene Wellen*. `planning-check` grün.

## 6. Plan-vs-Code-Diff

**Geliefert wie geplant, ohne Rest.** Die elf Punkte aus §2 des Plans sind alle eingelöst; die fünf Ausschlüsse aus §3 sind eingehalten:

- **Kein Schweregrad im Befund-Modell** — `finding.go` bekommt eine Konstante, kein Feld.
- **Kein `exempt-may-empty`** — kommt im Code nirgends vor.
- **Nichts für `exempt-paths`** — unverändert.
- **Kein Anwenden auf den eigenen Bestand** — die `.d-check.yml` führt weder den Schlüssel noch eine Regel mit `exempt-section-pattern`.
- **„Keine `--doctor`-Zeile"** — eingehalten im gemeinten Sinn: hinzugekommen ist der **Grund-Code-Klartext**, den die Verriegelung zwischen `AllReasons`, `reasonTexts` und dem Spezifikations-§4 für **jeden** neuen Code erzwingt. Das ist nicht die vom Antrag gewünschte **Zustands**-Zeile. Der Unterschied ist real, steht aber nirgends ausgesprochen — siehe **V-7**.

**Nicht im Plan, aber geliefert:** die Duplikats-Abweisung bei zwei verschiedenen Zahlen ist kein neuer Code, sondern eine Folge der Identitäts-Entscheidung; sie ist in der ADR benannt und am Produkt belegt (§1 Punkt 7).

**Risiken (§5 des Plans).** Alle vier tragen noch den Platzhalter *(bei Closure)* — planmäßig, der Slice liegt in `in-progress/`. Zum Stand der Verifikation:

| Risiko | Stand |
|---|---|
| „Die geprüfte Menge zu verkleinern bleibt eine Lockerung … wer die Zahl mitzieht, ohne die Aufzählung zu prüfen" | **weiter offen** — kein Sensor kann das halten, die ADR sagt das ausdrücklich. Kandidat für den dritten Ausgang (Beobachtungs-Register), nicht für einen Folge-Slice |
| „Eine deklarierte Zahl ist Autoren-Text und altert" | **weiter offen**, deckungsgleich mit dem ersten |
| „Neue Bauform ohne Präzedenz" | **weiter offen** — bestätigt: kein anderer Schlüssel dieses Moduls führt eine erwartete Anzahl. Die ADR führt dafür bereits einen Re-Evaluierungs-Trigger |
| „Die Abweichung ist ein Urteil über einen fremden Bestand" | **eingetreten, in verschärfter Form** — die Messung, die das Urteil trägt, ist zur Hälfte falsch (**V-2**). Das Urteil selbst bleibt tragfähig, seine Begründung nicht in der Form, in der sie in den kanonischen Quellen steht |

## 7. Befunde

**V-1 — `spec/spezifikation.md:3052` (HIGH).**
Der Commit `9c383ba` hat die **Kopfzeile der §6-Tabelle *Externe Verträge*** durch die Historie-Zeile ersetzt. Vorher stand dort eine vierspaltige Kopfzeile mit den Namen *Kennung*, *System*, *Version/Stand* und *Vertrag*; heute steht dort der vollständige Text der `2026-08-30`-Historie-Zeile — dieselbe, die (richtig) auch in §7 steht. Die Trennzeile darunter ist unverändert vierspaltig, die Kopfzeile trägt jetzt **zwei** Zellen. Nach GFM muss die Trennzeile in der Zellenzahl zur Kopfzeile passen, sonst wird die Tabelle **nicht als Tabelle erkannt**; unabhängig vom Renderer sind die Spaltennamen weg und `SPEC-064`, `SPEC-065`, `SPEC-066` stehen ohne Kopf. Die Baseline-Vorlage für dieses Stratum bestätigt die vierspaltige Form.
Das ist **genau die Fehlerklasse, die das Einfügen von `SPEC-077` erzeugt**: die neue §4-Zeile verschiebt die Zeilennummern, und der nächste Edit trifft die falsche Zeile. **Kein Gate fängt es** — `make doc-check`, `make gates` und `make fullbuild` sind grün, weil `links`, `anchors` und `ids` von einer kaputten Tabellen-Kopfzeile nichts wissen und keine `structure`-Regel dieses Repos den §6 der Spezifikation adressiert. Es gibt keine eingehende Referenz auf §6 oder auf die drei `SPEC-06x`-Kennungen, der Schaden ist also auf die Datei begrenzt — aber er steht in einer **kanonischen Quelle Rang 2**.
*Stand bei Abgabe dieses Reports:* der Befund gilt dem **committeten** Range; im Arbeitsbaum liegt bereits eine **uncommittete** Wiederherstellung der vierspaltigen Kopfzeile aus dem parallel laufenden Review-Kontext. Der Befund ist damit unabhängig zweimal gefunden worden und ist im Begriff, behoben zu werden — er bleibt hier stehen, weil er die Lücke benennt, die kein Sensor hält.

**V-2 — ADR-0078 §Entscheidung 2, `spec/lastenheft.md:3364`, Slice-Plan (HIGH).**
Die Messung, die die **Abweichung vom Antrag trägt**, ist zur Hälfte falsch. Behauptet wird: *„`--doctor` läuft in keinem Gate dieses Repos … und `--print-mk` verteilt kein Doctor-Target an Konsumenten."*
Gemessen gegen das Binary: `--print-mk` **verteilt sehr wohl** ein Doctor-Target — das erzeugte Fragment trägt ein eigenes, als *„erklärende Diagnose mit Fix-Kandidaten"* annotiertes `doc-doctor`-Target, dessen Rezept das Image mit `--doctor` fährt. Die erste Hälfte stimmt: `Makefile` 0 Fundstellen, `.github/workflows/` 0, `.githooks/pre-commit` 0.
Die **Schlussfolgerung** trägt weiterhin: `doc-doctor` ist ein Einzeltarget, das in keinem Sammel-Gate des Fragments steht — wer es fährt, sieht ohnehin nach. Falsch ist der **Satz**, der die Schlussfolgerung trägt, und er steht an drei Stellen: in einer **immutablen** `Accepted`-ADR (§3.5 — korrigierbar nur über einen `## Geschichte`-Anhang oder eine Folge-ADR), in der **Lastenheft-Historie** (Rang 1) und im Slice-Plan.
**Die richtige Formulierung existierte bereits:** das ausgehende CR-Antwortdokument sagt in seiner Mess-Tabelle *„kein `doc-doctor` im Gate-Pfad"* — präzise und wahr. Die Aussage ist auf dem Weg in die kanonischen Quellen **von der Messung weg** verallgemeinert worden. Klasse `BEO-009` (der Schluss reicht weiter als die gemessene Menge) und `BEO-012` (Reichweite einer zitierten Aussage).

**V-3 — `CHANGELOG.md`, `README.md`, `README.de.md` (LOW).**
Keine der drei Dateien nennt `exempt-expect-count` oder `section-exempt-mismatch`. Die Änderung ist **nutzersichtbar** (neuer opt-in-Schlüssel, neuer Grund-Code, zwei neue Exit-2-Ränder), und `AGENTS.md` §5 verlangt `CHANGELOG`-Pflege dafür; §3 des Slice nimmt sie nicht aus.
**Der nächstliegende Präzedenzfall stützt den aktuellen Stand:** der unmittelbare Vorgänger slice-179 hat denselben Schlüssel-Typ eingeführt, und sein `CHANGELOG`-Eintrag entstand erst im Release-Prep-Commit; die beiden READMEs zählen die Grundmengen-Schlüssel ebenfalls erst seither auf. Es gibt aber auch den anderen Präzedenzfall — ein Feat-Commit, der den `CHANGELOG` direkt fortschreibt —, und die Verifikation von slice-181 hat dieselbe Lücke als MEDIUM geführt. **Die Frage ist im Repo unentschieden**; als Befund gegen *diesen* Slice wiegt sie leicht, als **fällige Release-Prep-Fläche** ist sie zu notieren: beide READMEs führen heute eine abschließende Aufzählung der Grundmengen-Schlüssel, die der neue Schlüssel unvollständig macht.

**V-4 — ADR-0078 §Entscheidung 3 (LOW).**
*„`DC-FA-STRUCT-001` schreibt für genau das einen eigenen Code vor."* Das Lastenheft stellt den Grundsatz — jede Bedingung mit eigenem Grund-Code, weil jede eine andere Reparatur verlangt — für die Tabelle der **Bedingungen im Abschnitt** auf und wendet ihn ein zweites Mal auf die Zellengrenzen an. `exempt-expect-count` ist keine Bedingung, sondern eine Erklärung der Grundmenge, und für **diese** Familie sagt dieselbe Anforderung ausdrücklich das Gegenteil: `tasks-ignore-pattern` und `exempt-section-pattern` tragen *„keinen eigenen Grund-Code"*.
Die **Entscheidung ist richtig** und die Analogie trägt (andere Reparatur ⇒ anderer Code); nur das Wort *„schreibt vor"* reicht einen Schritt weiter als der Geltungsbereich der zitierten Stelle. `BEO-012` in milder Ausprägung. Das mitgeführte Dedup-Argument ist sachlich korrekt, für dieses Paar aber theoretisch: `section-missing` und `section-exempt-mismatch` sind einander ausschließende Zweige und könnten nie an derselben Befund-Adresse kollidieren.

**V-5 — Slice-Plan, Zeitpunkt der Beanspruchung (INFO).**
Das Feld `Verantwortlich:` und der Nachtlauf-Block entstanden im Commit `4d5a564`, **24 Sekunden vor** dem Move-Commit `a062fe8`, also während der Plan noch in `open/` lag. `AGENTS.md` §5 sagt *„spätestens bei der Beanspruchung"* (erfüllt) **und** *„ein Plan in `open/` trägt ihn noch nicht"* (für einen Commit lang nicht zutreffend). Nicht blockierend — aber **die Verifikation von slice-181 hat exakt dieselbe Beobachtung notiert**; es ist damit die zweite Instanz derselben Form.

**V-6 — Prozess und Lauf-Umgebung (INFO).**
Während der Messungen dieser Verifikation lag **kein Review-Report** für slice-182 in `docs/reviews/`; bei Abgabe ist er da, ungelesen und ohne Einfluss auf die Befunde oben. DoD-Punkt 12 ist damit in beiden Hälften bedient, beide in eigenen Kontexten.
**Zur Lauf-Umgebung, weil es die Messungen berührt:** während dieser Verifikation wurde `internal/hexagon/core/rules/structure.go` für rund eine Minute von einem **parallel laufenden Kontext** mutiert (dieselbe Mutation wie M2) und danach zurückgesetzt. Der Baum war vor **und** nach jedem hier berichteten Gate-Lauf sauber, und das Runtime-Image trug am Anfang wie am Ende dieselbe ID `sha256:f59c1333ad9a…` — da ein Docker-Build inhaltsadressiert ist, belegt das die Identität des Baums über alle Läufe. **Keine Messung dieses Reports ist davon berührt.** Notiert, weil zwei Kontexte gleichzeitig am selben Arbeitsbaum mutiert haben und das WIP-Limit die Rollen trennt, nicht den Baum.

**V-7 — Slice-Plan §3, zweiter Ausschluss (INFO).**
*„Keine `--doctor`-Zeile."* Geliefert **ist** eine `--doctor`-Zeile: der Klartext des neuen Grund-Codes in `reasonTexts()`. Das ist kein Bruch des Ausschlusses — der Ausschluss gilt der vom Antrag gewünschten **Zustands**-Zeile, während der Grund-Code-Klartext von der Verriegelung zwischen `AllReasons`, `reasonTexts` und dem Spezifikations-§4 für jeden neuen Code **erzwungen** ist. Die Unterscheidung steht nirgends; ein späterer Leser des Plans findet im Diff eine `--doctor`-Ergänzung, die der Plan ausschließt.

## 8. Verdikt

**Konform gegen DoD und Spec — mit einem Defekt in einer kanonischen Quelle, der vor der Closure zu beheben ist.**

Das **Verhalten** ist vollständig und richtig: alle elf Akzeptanzkriterien von `DC-FA-STRUCT-001` in Version 0.78.0 halten gegen das laufende Binary, nicht nur gegen die Tests; beide Config-Ränder sind fail-closed mit Exit 2; die Null ist erlaubt und von „nicht gesetzt" unterscheidbar; die Trennung zwischen Konfigurationsdefekt und Bestandszustand hält auch mit gesetzter Zahl; die beidseitige Drift greift auch bei verbleibender Restmenge. Die Byte-Identität ohne den Schlüssel ist gegen ein aus `a062fe8` gebautes Vorgänger-Image in **vier** Prüfmengen und **vier** Ausgabeformen gemessen, und die von der ADR genannte Zahl **169** reproduziert exakt. Die Fitness Function der ADR hält alle fünf ihrer Mutations-Zahlen; eine sechste, selbst gewählte Mutation zeigt, dass auch die siebte Zusage einen eigenen Wächter hat. Alle acht gefahrenen Gates sind grün mit Exit 0.

Zwei Dinge stehen dem sauberen Abschluss entgegen:

1. **V-1** ist ein echter Schaden an `spec/spezifikation.md` — die §6-Tabelle hat keinen Kopf mehr, weil eine Historie-Zeile ihn überschrieben hat. Kein Gate fängt das. Er ist in einem Commit reparierbar und sollte **vor** der Closure repariert werden, nicht danach: `verify-closure-notes` sieht ihn ebenso wenig wie `doc-check`.
2. **V-2** ist kein Schaden am Produkt, aber ein Riss in der Begründungskette: der Satz, der die Ablehnung der beantragten Form trägt, ist gegen das eigene Binary widerlegbar — mit **einem** Aufruf. Er steht in einer immutablen ADR und in der Lastenheft-Historie; die präzise Fassung existiert im ausgehenden CR-Dokument und hätte nur übernommen werden müssen. Die getroffene Entscheidung bleibt richtig; ihre Begründung braucht die Korrektur, die §3.5 vorsieht.

**Der Gegenstand des Slice — die Trennung zweier Zustände, die sich einen Grund-Code teilten — ist erreicht und am Produkt belegt.** Die beiden Befunde betreffen die Dokumentations-Schicht, nicht das Verhalten.
