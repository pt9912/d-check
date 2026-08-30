# Verifikation slice-178 — offene Task-Items auf dem rohen Abschnitt

**Gegenstand:** `3f8049e^..HEAD` (3 Commits: Prep `3f8049e`, Claim-Move `b2530f7`, feat `23b0f56`).
**Rolle:** unabhängiger Verifier — geprüft wird gegen **DoD und Spec** (Behaviour / Architecture Fitness), nicht gegen Plan/ADR-Maintainability.
**Gefahrene Sensors:** `make gates` · `make fullbuild` · `make completeness-check` · `make verify-closure-notes` · `make adr-check RANGE=…` · `make trace-check RANGE=…` · über 30 Produktläufe des lokal gebauten Images gegen **eigene** Probe-Dokumente und Probe-Konfigurationen · 7 Mutations- und 1 Kontrollprobe gegen eine **Scratch-Kopie** des Baums (Build-Stage `test`) · ein selbst aus `b2530f7` gebautes **Vorgänger-Image** für die Byte-Identität.
**Nicht verändert:** kein Repo-Artefakt außer dieser Datei. Alle Kopien, Mutationen und Probe-Bäume liegen außerhalb des Repos.

---

## 1. DoD-Tabelle (§4 des Slice-Plans)

| # | Behauptet | Gemessen | Ergebnis |
|---|---|---|---|
| 1 | `max-open-tasks` im Schema, im Lastenheft (Bump + Historie) und in der Spezifikation geführt, mit eigenem Grund-Code; **explizit** < 0 ⇒ Exit 2, mit Test | Schema: `rawStructure.MaxOpenTasks` (`configyaml.go`, Zeiger-Typ) → `model.StructureRule.MaxOpenTasks` (`config.go`). Lastenheft `0.78.0` → **`0.79.0`**, Bedingungs-Tabellenzeile bei `:2524`, Prosa-Block, drei Akzeptanzkriterien, Config-Rand-Satz, Historie-Zeile bei `:3400` — **unter `## 7. Historie`** geprüft, nicht nur an der Zeilennummer. Spezifikation: §2-Schema `:2916` (unter `### SPEC-005`), §4 `SPEC-078` `:3064` (unter `## 4. Grund- und Fehler-Codes`), Historie `:3093` (unter `## 7. Historie`), Ablauf-Block und Schritt 7. Grund-Code `section-tasks-open` in `model`, `AllReasons()`, `reasonTexts()`. Produkt mit `-1`: `d-check: error: … structure[0]: max-open-tasks -1 muss >= 0 sein`, **Exit 2**; Test-Fall `max-open-tasks negativ` in `TestDecode_StructureFehler` | **erfüllt** — mit **V-1** und **V-2** (Spezifikation widerspricht sich an zwei Stellen weiter) |
| 2 | Die Blindstelle ist gemessen geschlossen: derselbe Backtick-Fall, an dem der Vorgänger auf 0 Befunde fiel, meldet jetzt — **Ausgabe vorher und nachher in der Commit-Botschaft** | Sache **belegt**, Form **nicht**. Eigene Probe (Backtick vor und hinter dem Item, ein Absatz): mit `max-open-tasks: 0` → `docs/backtick.md:6 … section-tasks-open`, Exit 1; dieselbe Datei mit `forbid-pattern: '- \[ \]'` → **kein Befund** für diese Datei, während drei Nachbardateien melden. Der Vorzustand steht zusätzlich **im Test selbst** (`TestMaxOpenTasks_BacktickSchaltetNichtAb` fährt beide Formen in einer Funktion). **Die Commit-Botschaft trägt weder eine Vorher- noch eine Nachher-Ausgabe** — sie nennt den Backtick nur als Begründung | **in der Sache erfüllt, in der geforderten Form nicht** (**V-5**) |
| 3 | Alle Bullet-Formen gemessen: `-`, `*`, `+` und die geordnete Liste, je offen und je gehakt; eingerückt und mit Tab-Trenner. Erwartung und Ergebnis je Fall | Eine Probe-Datei mit elf Zeilen, ein Lauf: **sechs** Befunde auf den Zeilen 5–10 (`-`, `*`, `+`, `1.`, vier Leerzeichen eingerückt, Tab als Trenner) und **kein** Befund für `- [x]`, `- [X]`, den blanken Listenpunkt, `> - [ ]` und `- [\t]`. Tabellengetriebener Test mit `meldet`-Erwartung je Fall deckt dieselben elf | **erfüllt** |
| 4 | Fence-Treue: ein Task-Item **innerhalb** eines Fenced-Blocks zählt **nicht**. Gemessen, nicht behauptet | Probe mit einem Item im Fence (Zeile 6) und einem dahinter (Zeile 9): genau **ein** Befund, auf `:9`. Mutation *Fence-Gate entfernt* macht genau `TestMaxOpenTasks_FenceBleibtAussen` rot (§3) | **erfüllt** |
| 5 | Ein Befund **je offenem Item**, auf **seiner** Zeile; zwei offene Items ⇒ zwei Befunde | Probe mit offen/gehakt/offen: Befunde auf `:5` und `:7`, keiner auf `:6`, keiner auf der Abschnitts-Überschrift `:3`. Zum Vergleich meldet `forbid-pattern` über dieselbe Datei **einen** Befund auf `:3` | **erfüllt** |
| 6 | Umkehr-Probe: je Zusage eine Mutation, die genau einen Test rot macht | **Sieben** der acht in der Commit-Botschaft genannten Mutationen unabhängig nachgefahren, Kontrolllauf der unmutierten Kopie grün. Alle sieben liefern **exakt** die behauptete Zahl roter Tests (§3). Die achte (*bereinigt statt roh gelesen ⇒ 2 rot*) ist nicht mechanisch reproduzierbar; ihr Verhalten ist stattdessen am Produkt belegt (Zeile 2 dieser Tabelle) | **erfüllt** |
| 7 | Eine ADR begründet die drei Entscheide aus §2 und ist im ADR-Index eingetragen | ADR-0074 `Accepted`, Index-Zeile `docs/plan/adr/README.md:84`. Die drei Entscheide dieses Baus stehen als **`## Geschichte`-Zeilen** (kein zweiter Grund-Code-Entscheid als neue ADR — §3.5 verbietet den Kern-Edit). `make adr-check RANGE=3f8049e^..HEAD` → **Exit 0** | **erfüllt** — mit **V-6** (die Fitness Function des Kerns ist nicht mitkorrigiert) |
| 8 | Das Benutzerhandbuch führt die Bedingung dort, wo es die übrigen `structure`-Schlüssel führt | `docs/user/benutzerhandbuch.md`: drei Kommentarzeilen im Konfigurations-Fence beim Schlüssel und drei Absätze mit der Abgrenzung zu `max-tasks`, der Lexik-Bequemlichkeit und den drei Grenzen | **erfüllt** — mit **V-3** (die Modul-Tabelle derselben Datei blieb stehen) |
| 9 | `make gates` grün (Exit explizit); unabhängiger Review; Verifikation gegen DoD/Spec — beide in eigenen Kontexten | `make gates` → **Exit 0** (§5). `make fullbuild` → **Exit 0**. Diese Datei ist die Verifikation; der Review lief in einem eigenen Kontext | **erfüllt, soweit hier prüfbar** |

Die Haken in §4 des Slice-Plans stehen noch offen — planmäßig: der Closure-Body ist ein eigener Commit nach dem Lifecycle-Move.

## 2. Akzeptanzkriterien `DC-FA-STRUCT-001` (0.79.0) gegen das laufende Produkt

Alle Läufe: lokal gebautes Image, `--network none`, read-only-Mount, **eigene** Probe-Bäume und Probe-Konfigurationen außerhalb des Repos.

| Akzeptanzkriterium | Probe | Ergebnis |
|---|---|---|
| **Roh gezählt:** mehrzeilige Inline-Spanne um ein offenes Item ⇒ `section-tasks-open` auf der **Item-Zeile**, während `forbid-pattern` über denselben Abschnitt nichts meldet | Ein Absatz, Backtick vor und hinter dem Item | `…:6 … section-tasks-open`, Exit 1 · `forbid-pattern`: kein Befund für diese Datei, Exit 1 nur wegen der Nachbardateien | **hält** |
| **Abwesender Schlüssel ⇒ Bedingung aus; die explizit gesetzte Null ist davon unterscheidbar** | Dieselbe Datei mit vier offenen Items, einmal ohne Schlüssel (`non-empty: true` statt dessen), einmal mit `max-open-tasks: 0` | ohne Schlüssel `0 Befund(e)`, **Exit 0** · mit `: 0` vier Befunde, Exit 1 | **hält** |
| **Lexik:** alle vier Marker, eingerückt und mit Tabulator ⇒ jedes offene meldet, kein gehaktes | Elf-Zeilen-Probe | 6 Befunde auf `:5`–`:10`; `[x]`, `[X]` stumm | **hält** |
| **Granularität:** zwei offene Items ⇒ **zwei** Befunde auf je eigener Zeile | offen/gehakt/offen | `:5` und `:7` | **hält** |
| **Schwelle:** `max-open-tasks: 2` bei vier offenen Items ⇒ genau **zwei** Befunde, auf dem dritten und vierten Item in Dokument-Reihenfolge | vier Items auf `:5`–`:8` | Befunde auf `:7` und `:8`, Meldung *„offenes Task-Item über der Grenze von 2"* · dieselbe Datei mit `: 4` ⇒ **0 Befunde, Exit 0** | **hält** |
| **Fence:** Item im Fenced-Block ⇒ kein Befund; dasselbe Item außerhalb ⇒ einer | Fence auf `:5`–`:7`, Item dahinter auf `:9` | nur `:9` | **hält** |
| **Abschnittsgrenze:** ein offenes Item im **folgenden** Abschnitt zählt nicht mit | Zusatzprobe mit drei Ebenen: Item im Abschnitt (`:5`), Item in einer **Unterebene** `###` (`:9`), Item hinter der nächsten `##`-Überschrift (`:13`) | Befunde auf `:5` **und** `:9`, **nicht** auf `:13` — die Bedingung endet an der nächsten Überschrift gleicher oder höherer Ebene und schließt Unterabschnitte ein, genau wie die Spezifikation sagt | **hält** |

### 2.1 Die Grenzen, wie sie deklariert sind — alle vier gemessen

| Deklarierte Grenze | Messung | Ergebnis |
|---|---|---|
| Eine **einzeilige** Inline-Spanne meldet **nicht** | `So sieht es aus: ` + Item in Backticks, eine Zeile | 0 Befunde, Exit 0 | **hält** |
| Eine **mehrzeilige** Spanne **zählt mit** | Backtick vor und hinter dem Item, drei Zeilen | 1 Befund auf der Item-Zeile | **hält** |
| Ein Item im **Blockquote** zählt für **keine** der beiden Bedingungen | Probe-Datei nur mit `> - [ ]` und `- [\t]`, dazu eine Kontrolldatei | `max-open-tasks: 0` → nur die Kontrolldatei meldet · `max-tasks: 0` → nur die Kontrolldatei meldet (`section-oversized`) | **hält** |
| Ein **Tabulator in der Box** zählt für **keine** der beiden | dieselbe Probe, mit `cat -A` als Tabulator verifiziert | wie oben, beide Bedingungen stumm | **hält** |

### 2.2 Der Zusatznutzen gegenüber der abgelösten Form — gemessen, nicht behauptet

Eine Probe-Datei, die **nur** `* [ ]` und `+ [ ]` trägt:

| Konfiguration | Ausgabe |
|---|---|
| `max-open-tasks: 0` | zwei Befunde, `:5` und `:6`, **Exit 1** |
| `forbid-pattern: '- \[ \]'` | `1 Datei(en) geprüft, 0 Befund(e)`, **Exit 0** |

Das ist die `BEO-003`-Klasse, mit der die Bedingung begründet ist, am eigenen Produkt belegt.

### 2.3 CRLF

Eine Probe mit CRLF-Zeilenenden meldet das offene Item auf `:5` und lässt das gehakte auf `:6` stehen — die Zeilen-Zerlegung trägt.

## 3. Mutations- und Kontrollproben (Umkehr-Probe, `BEO-023`)

Gefahren gegen eine `git archive`-Kopie des Baums außerhalb des Repos, je ein voller `make test`-Lauf über die Build-Stage `test` unter eigenem Image-Tag. **Kontrolllauf der unmutierten Kopie: Exit 0, null rote Tests.**

| Mutation | behauptet | gemessen | fangender Test |
|---|---|---|---|
| Fence-Gate entfernt | 1 rot | **1** | `TestMaxOpenTasks_FenceBleibtAussen` |
| Box nicht auf leer verengt | 3 rot | **3** | `…_AlleMarkerFormen`, `…_EinBefundJeItem`, `…_NurImAbschnitt` |
| Schwelle ignoriert | 1 rot | **1** | `TestMaxOpenTasks_SchwelleMeldetNurDenUeberhang` |
| über die ganze Datei statt im Abschnitt | 1 rot | **1** | `TestMaxOpenTasks_NurImAbschnitt` |
| ein Befund je Datei statt je Item | 2 rot | **2** | `…_EinBefundJeItem`, `…_SchwelleMeldetNurDenUeberhang` |
| negativer Wert geschluckt | 1 rot | **1** | `TestDecode_StructureFehler` |
| Verdrahtung Config → Modell entfernt | 1 rot | **1** | `TestDecode_StructureFehler` |
| bereinigt statt roh gelesen | 2 rot | **nicht gefahren** | — (Verhalten stattdessen am Produkt belegt, §2) |

Die vorletzte Zeile ist die, die dem ersten Bau fehlte: sie ist die einzige, die den Weg vom YAML-Schlüssel bis ins Modell prüft. Sie beißt.

Coverage der neuen Funktionen aus dem `coverage-gate`-Lauf: `offenerHaken` **100,0 %**, `structureOpenTasks` **100,0 %**, `structureAmAbschnitt` **100,0 %**, `checkStructureFile` **100,0 %**.

## 4. Rückwärtskompatibilität und Bestandszahlen

**Vorgänger-Image selbst gebaut** aus `b2530f7` (dem Commit **vor** dem feat) über `make build` mit eigenem Tag; derselbe Korpus, dieselbe Konfiguration, getrennte Ströme:

| Messung | Vorgänger | HEAD |
|---|---|---|
| 178 Slice-Pläne, `max-tasks: 3`, Selektor `^## `, `sections: each` | 164 Befunde | 164 Befunde |
| `sha256` des Befundsatzes (stdout) | `51fb272a622a1b73…` | `51fb272a622a1b73…` |
| `diff` stdout / stderr | leer / leer | — |

**Byte-Identität ohne den Schlüssel ist damit unabhängig belegt.** Die in der Commit-Botschaft genannte Zahl **276** reproduziert dabei **nicht** — sie ist korpus- und konfigurationsabhängig, und die Botschaft nennt weder Glob noch Selektor. Drei nachgefahrene Varianten liefern 164, 331 bzw. 161 Befunde. Die tragende Aussage (`diff` leer) hält; die Zahl ist ohne die fehlende Angabe nicht nachprüfbar (**V-7**).

**Die Bestandszahlen reproduzieren dagegen exakt.** Korpus: die 176 `done/`-Slice-Pläne, flach kopiert.

| Konfiguration | Befunde | Dateien |
|---|---|---|
| `max-open-tasks: 0`, Selektor `^## ` | **144** | **37** |
| dieselbe Menge, `forbid-pattern: '- \[ \]'` | 37 | **37** |

Die beiden 37er-Mengen sind **elementweise identisch** (`diff` der sortierten Dateilisten leer). Beide Behauptungen der Commit-Botschaft — *144 Befunde in 37 Dateien* und *dieselben 37 über die abgelöste Form* — halten. Das Ergebnis ist zugleich unabhängig vom Selektor: `^#` und `^## \d` liefern für `max-open-tasks` dieselben 144/37.


### 4.1 Kosten der rohen Lesung im eigenen Bestand: null gemessene Falsch-Positive

Das zweite Risiko aus §5 des Slice-Plans ist, dass die rohe Lesung Falsch-Positive zurückholt — ein Task-Item, das nur Prosa **über** den Marker ist. Über die 144 Befunde einzeln gemessen: für jeden wurde die Backtick-Parität des vorangehenden Absatzes bestimmt. Nur **eine** Zeile liegt hinter einer ungeraden Zahl, und sie ist bei Sicht **kein** Falsch-Positiver, sondern ein echter offener DoD-Punkt — die ungerade Parität stammt aus den Nachbar-Items derselben Liste.

Das dreht die Aussage um: diese eine Zeile ist ein **Bestands-Beleg für die Blindstelle**, nicht gegen sie. Die bereinigt lesenden Bedingungen sehen sie nicht, weil die Spanne über ihr offen ist; `max-open-tasks` sieht sie. Der Slice-Plan sagt in §1, er stehe „auf dem konstruierten Fall, nicht auf einem Bestands-Fund" — nach dieser Messung gibt es einen, und er ist ein echter unquittierter Haken.

**Grenze dieser Messung:** die Absatz-Erkennung ist die einfache (bis zur nächsten Leerzeile) und trennt Listen-Items nicht; sie ist eine Näherung an die Paarung des Produkts, keine Nachbildung.
## 5. Gate-Läufe (echte Ausgabe)

| Gate | Exit | Ausgabe |
|---|---|---|
| `make gates` | **0** | `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` · `d-check: 620 Datei(en) geprüft, 0 Befund(e)` · `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%` · `Ran 55 rules on 55 files: 0 findings.` |
| `make fullbuild` | **0** | `[fullbuild] green — image-hash sha256:5305087a11fc…`; enthält `image-test` (1000/0), `bench`, `completeness-check`, `verify-closure-notes` |
| `make completeness-check` | **0** | `50 Anforderung(en), 0 Waise(n).` |
| `make verify-closure-notes` | **0** | `d-check: 558 Datei(en) geprüft, 0 Befund(e)` (Profil `.d-check.closure.yml`, Module `planning`, `structure`, `spans`) |
| `make adr-check RANGE=3f8049e^..HEAD` | **0** | `d-check: 620 Datei(en) geprüft, 0 Befund(e)` — die `## Geschichte`-Anhänge an der `Accepted`-ADR passieren regelkonform |
| `make trace-check RANGE=3f8049e^..HEAD` | **0** | `d-check: 620 Datei(en) geprüft, 0 Befund(e)` |

`lint`, `test`, `arch-check`, `coverage-gate` und `gate-consistency` sind Glieder von `make gates` und in dessen Exit 0 enthalten; die Deklarations-Konsistenz Doku ↔ Makefile (`gate-consistency`, Modul `targets`) meldete `620/0`.

**Deklarations-Konsistenz des neuen Grund-Codes**, je einzeln nachgesehen:

| Ort | Stand |
|---|---|
| `model.ReasonSectionTasksOpen` | vorhanden, mit Abgrenzungs-Kommentar |
| `AllReasons()` | vorhanden — und der Go-Test `TestAllReasonsDeckungGegenSpezifikationGrundCodes` hält §4 ↔ `AllReasons()` maschinell, er lief grün |
| `--doctor`-Klartext | vorhanden; am Produkt gegengeprüft: *„Der Abschnitt trägt mehr offene Task-Items, als die Regel erlaubt — den Haken setzen oder den Punkt auflösen"* |
| `--print-config`-Gerüst | vorhanden, drei Kommentarzeilen (Ausgabe des Images, Zeilen 180–182) |
| Spezifikation §2-Schema | vorhanden |
| Spezifikation §4 (`SPEC-078`) | vorhanden |
| Lastenheft-Bedingungstabelle | vorhanden |
| Benutzerhandbuch, Konfigurations-Fence | vorhanden |
| Benutzerhandbuch, **Modul-Tabelle** | **fehlt** (**V-3**) |
| `README.md` / `README.de.md` | **fehlt** (**V-4**, nach Repo-Praxis Release-Prep) |

## 6. Spec-Straten-Richtung (§3.4)

Alle in `spec/lastenheft.md` und `spec/spezifikation.md` **addierten** Zeilen des feat-Commits gegen die vier detektierbaren Kategorien geprüft (ADR-Kennung, Slice-Kennung, Wellen-Kennung, Commit-Hash): **keine Fundstelle**. Die Historie-Zeile des Lastenhefts sagt *„Begründung in begleitender ADR"* ohne Kennung — die Form, die dieses Repo dafür fährt. Gegenprobe mit dem Modul `matrix` allein über den ganzen Baum: `620 Datei(en) geprüft, 0 Befund(e)`, Exit 0.

**Historie-Zeilen in der richtigen Tabelle** — nicht an der Zeilennummer, sondern an der umschließenden Überschrift geprüft, weil eine neue `SPEC`-Kennung die Nummern verschiebt:

| Zeile | umschließende Überschrift |
|---|---|
| Lastenheft `0.79.0` | `## 7. Historie` |
| Lastenheft-Bedingungszeile | `### DC-FA-STRUCT-001 — …` |
| Spezifikation `SPEC-078` | `## 4. Grund- und Fehler-Codes` |
| Spezifikation Schema-Zeile | `### SPEC-005 — .d-check.yml` |
| Spezifikation Historie 2026-08-30 | `## 7. Historie` |

Alle fünf sitzen richtig. Die in diesem Repo zweimal aufgetretene Fehlplatzierung ist hier **nicht** eingetreten.

## 7. Plan-vs-Code-Diff

| §2-Punkt des Plans | Geliefert |
|---|---|
| 1 — neue Bedingung `max-open-tasks`, eigener Grund-Code | ja, unverändert |
| 2 — der Kommentar über `structureConditions` wird „auf drei gezogen" | geliefert als **vier**. Abweichung **in die richtige Richtung**: §2 des Plans korrigiert die Zahl selbst („Es sind vier Bedingungen auf rohem Text, nicht drei") |
| 3 — kein generischer `raw`-Schalter | ja |
| 4 — Befund auf der Item-Zeile, einer je Item | ja |
| 5 — vor dem Scharfschalten rot messen, je Fall Erwartung und Ergebnis | ja, im Test tabellengetrieben; hier unabhängig am Produkt nachgefahren |
| 6 — ADR, Lastenheft-Bump, Spezifikations-Verfeinerung, Handbuch | geliefert, **aber**: die ADR ist kein neues Dokument, sondern fünf `## Geschichte`-Zeilen an ADR-0074. Das ist §3.5-konform und in der Sache richtig; der DoD-Wortlaut („eine ADR begründet … und ist im ADR-Index eingetragen") ist damit erfüllt, weil ADR-0074 indiziert ist |
| 7 — `make gates`, Review, Verifikation, Closure | Gates grün, Review separat gelaufen, dies ist die Verifikation, Closure steht aus |

**Über den Plan hinaus geliefert** (im Plan nicht genannt, in der Commit-Botschaft benannt): die Zerlegung von `checkStructureFile` in `structureAmAbschnitt`, weil der neue Zweig die Komplexitätsschwelle riss. Gemessen: keine Suppression, keine gesenkte Schwelle, `lint` grün, beide Funktionen 100 % gedeckt. Das ist eine strukturelle Änderung an bestehendem Code, die kein DoD-Punkt vorsah — sie ist tragbar, weil ihre Wirkung durch die Byte-Identitäts-Messung in §4 abgedeckt ist.

**Im Plan, aber nicht geliefert:** nichts.

### §5 — Risiken

Alle vier Risiken stehen unverändert mit `**Ausgang:** *(bei Closure)*` — planmäßig offen, keiner ist im Lauf eingetreten:

| Risiko | Stand nach dieser Verifikation |
|---|---|
| zwei Bedingungen über dieselbe Frage, verschiedene Antworten | **weiter offen** — und die Gefahr ist gemessen real: über denselben Bestand liefert die eine 144, die andere 37 Befunde |
| die rohe Lesung holt Falsch-Positive zurück | **weiter offen**, im eigenen Bestand nicht eingetreten: unter den 144 Befunden über `done/` ist **kein** Falsch-Positiver (§4.1); die Inline-Grenze ist kleiner als der Plan zunächst annahm (nur mehrzeilig) |
| die Fähigkeit entsteht für **einen** Konsumenten | **weiter offen** — kein Konsument in diesem Commit, das Closure-Profil führt den Schlüssel nicht (planmäßig, §3 des Plans) |
| der Grund-Code-Raum wächst um einen weiteren `section-*`-Code | **weiter offen** — der Raum steht jetzt bei **16** `section-*`-Grund-Codes |

## 8. Befunde

**V-1 (HIGH, Spec-Stratum widerspricht sich) — die Spezifikation sagt an der tragenden Stelle weiterhin „zwei".**
`spec/spezifikation.md`, Schritt 5 der `.a`-Verfeinerung: *„Alle Bedingungen arbeiten auf **diesem** Text — mit **zwei** benannten Ausnahmen: die Chronologie-Bedingung … die Überschriften-Bedingung …"*. Es sind **vier** — die Zellenlänge und `max-open-tasks` fehlen in der Aufzählung. Der Lastenheft-Kopf ist in genau diesem Commit von zwei auf vier korrigiert, und die Commit-Botschaft sagt ausdrücklich *„Lastenheft und Spezifikation führten ZWEI Bedingungen … Es sind VIER"* — die Spezifikations-Hälfte dieser Zusage ist nicht eingelöst. **Dies ist wörtlich derselbe Befund, den §2 des Slice-Plans für den ersten Bau protokolliert** (*„Die `.a`-Verfeinerung Schritt 5 sagte weiter ‚zwei benannte Ausnahmen'"*). Er ist wiedergekehrt.

**V-2 (MEDIUM, Enumerations-Drift in der Spec) — die Bedingungs-Tabelle in Schritt 6 hat keine Zeile für `max-open-tasks`.**
Die Tabelle in Schritt 6 führt **jede** andere Bedingung, auch die, die zusätzlich einen Prosa-Block haben (`headings-match`, `table.column`). `max-open-tasks` steht nur als Prosa-Block dahinter. Wer die Tabelle als Aufzählung liest — und sie ist als solche eingeführt —, findet die Bedingung nicht. Auch das steht als Befund des ersten Baus in §2 des Plans (*„Schritt 6 hatte keine Zeile für die neue Bedingung"*); geliefert ist der Prosa-Block, nicht die Zeile.

**V-3 (MEDIUM, dieselbe Datei, die der Commit anfasst) — die Modul-Tabelle des Benutzerhandbuchs ist stehengeblieben.**
`docs/user/benutzerhandbuch.md`: die `structure`-Zeile sagt *„bis zu neun Bedingungen"* (es sind zehn) und ihre Grund-Code-Spalte listet 14 Codes ohne `section-tasks-open`. **Sie fehlt dort schon einen Schlüssel länger:** auch `section-exempt-mismatch` aus `0.78.0` steht nicht darin, obwohl derselbe Tabellentext ihn in der Prosa nennt — die Spalte hat bereits einen Release-Prep überlebt. Der Commit editiert diese Datei an zwei anderen Stellen; die dritte blieb liegen.

**V-4 (LOW, Release-Prep-Fläche) — `README.md` und `README.de.md` sagen ebenfalls „neun Bedingungen"** und zählen die Grund-Codes ohne den neuen auf. Nach der Commit-Historie dieses Repos ist das planmäßig Release-Prep-Arbeit (*„Release-Prep v0.66.0 — … die neunte Bedingung erreicht die READMEs"*), nicht Sache des feat-Commits. Notiert, damit es beim nächsten Release-Prep nicht zwei Rückstände statt einem sind.

**V-5 (LOW, DoD-Wortlaut) — die Commit-Botschaft trägt die geforderte Vorher-/Nachher-Ausgabe nicht.**
DoD-Punkt 2 verlangt sie ausdrücklich. Die Botschaft nennt den Backtick-Ausfall als Begründung und die Immunität dagegen als Ergebnis, zeigt aber keine Ausgabe. Die Sache ist belegt — im Test und, unabhängig, in dieser Verifikation —, die zugesagte Form nicht.

**V-6 (LOW, ADR-Kern altert) — die Fitness Function von ADR-0074 beschreibt weiter den zurückgenommenen Bau.**
Der Kern sagt *„Sieben Kern-Tests"* und zählt auf, was sie decken; es sind **neun** Testfunktionen, und die Aufzählung nennt weder die **Abschnittsgrenze** noch die **Inline-Code-Grenze** — also gerade die zwei, die der erste Bau nicht hatte. Der `## Geschichte`-Anhang korrigiert die **Mutationszahl** (drei ⇒ acht), nicht die Testzahl und nicht die Deckungsliste. Da §3.5 den Kern-Edit verbietet, wäre die Reparatur eine weitere Geschichte-Zeile.

**V-7 (INFO) — die Zahl „276" ist nicht nachprüfbar.**
Die Commit-Botschaft nennt sie ohne Glob und ohne Abschnitts-Selektor; drei nachgefahrene Korpus-/Selektor-Varianten liefern 164, 331 und 161. Die tragende Aussage (`diff` leer) ist unabhängig belegt, die Zahl ist es nicht. Das ist die `BEO-020`-Nachbarschaft: eine Zahl, die niemand nachrechnen kann, belegt nichts.

**Nicht gefunden:** kein Verhaltens-Defekt. Jede Zusage der Spezifikation über `max-open-tasks` — roh statt bereinigt, geteilte Lexik, Fence-Treue, ein Befund je Item auf seiner Zeile, Schwelle als Überhang, Abschnittsgrenze, kein Leerlauf-Fall, Config-Rand — ist am laufenden Produkt eingelöst. Alle vier deklarierten Grenzen sind genau so, wie sie dastehen: nicht enger und nicht weiter.

## 9. Verdikt

**Konform im Verhalten, nicht abschlussreif in der Spezifikation.**

Das Produkt tut, was `DC-FA-STRUCT-001` in `0.79.0` und die `.a`-Verfeinerung sagen — über 30 Produktläufe, sieben Mutationsproben, ein selbst gebautes Vorgänger-Image und zwei exakt reproduzierte Bestandszahlen belegen es. Kein Akzeptanzkriterium fällt, keine deklarierte Grenze ist über- oder untertrieben, `make gates` und `make fullbuild` schließen mit Exit 0, die Abwärts-Sperre §3.4 ist gewahrt und alle fünf neuen Tabellenzeilen sitzen in der richtigen Tabelle.

**Was dem Abschluss entgegensteht, ist V-1**: die Spezifikation — Rang 2 der Source Precedence — widerspricht sich in derselben Verfeinerung, die dieser Commit erweitert, und zwar an der Stelle, die erklärt, **welcher Text gelesen wird**. Das ist die Kern-Aussage dieser Bedingung. Dass die Commit-Botschaft die Korrektur für beide Spec-Straten behauptet und nur eines liefert, macht den Befund schwerer, nicht leichter: er ist als erledigt deklariert. Zusammen mit **V-2** heißt das, dass ein Leser der Spezifikation die neue Bedingung weder in der Ausnahmen-Aufzählung noch in der Bedingungs-Tabelle findet.

**Empfehlung:** V-1 und V-2 vor dem Lifecycle-Move nach `done/` beheben (beides Spezifikations-Prosa, kein Code, kein Bump — die Bedingung ist schon 0.79.0). V-3 im selben Zug, weil der Commit die Datei ohnehin anfasst und der Rückstand dort bereits zwei Schlüssel tief ist. V-5 und V-6 sind Form-Befunde ohne Wirkung auf das Verhalten; V-6 braucht eine Geschichte-Zeile, keinen Edit. V-4 gehört in den nächsten Release-Prep.
