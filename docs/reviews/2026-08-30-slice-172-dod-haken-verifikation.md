# Verifikation slice-172 — der Closure-Übergang wird gewächtert

**Gegenstand:** `fd4cbb4^..HEAD` (3 Commits: Prep `fd4cbb4`, Claim-Move `e8401c7`, feat `79bf375`).
**Rolle:** unabhängiger Verifier — geprüft wird gegen **DoD und Spec** (Behaviour / Architecture Fitness), nicht gegen Plan/ADR-Maintainability.
**Gefahrene Sensors:** `make gates` · `make doc-check` · `make verify-closure-notes` · `make fullbuild` · `make trace-check RANGE=…` · `make adr-check RANGE=…` · 19 Produktläufe des lokal gebauten Images gegen eine **Kopie** des Baums und gegen **eigene** Probe-Bäume und Probe-Konfigurationen, alle außerhalb des Repos · ein aus `e93d6a9` entpackter **Alt-Baum** für die historische Gegenrechnung.
**Nicht verändert:** kein Repo-Artefakt außer dieser Datei. `git status --short` ist vor und nach den Proben leer; alle Mutationen liegen in Kopien außerhalb des Repos.

---

## 1. DoD-Tabelle (§4 des Slice-Plans)

| # | Behauptet | Gemessen | Ergebnis |
|---|---|---|---|
| 1 | Eine `structure`-Regel im Closure-Profil meldet einen offenen DoD-Haken in `done/slice-*` mit eigener, reparatur-benennender Meldung | Regel liegt in `.d-check.closure.yml` als fünfte `structure`-Bedingung: `section-pattern` auf `^#{2,3} [0-9]+\. Definition of Done`, `sections: each`, `max-open-tasks: 0`, `hint`. Positiv-Probe (§4 dieses Berichts): Befund trägt **Datei, Zeile des Hakens, Grund-Code `section-tasks-open`** und den verfassten Hinweis-Text; Exit 1 | **erfüllt** |
| 2 | `spans` läuft im selben Profil, **belegt statt angenommen**; Beleg ist eine Probe mit offenem Haken hinter vergessenem Fence — `structure` allein schweigt, das Profil meldet | Makefile-Rezept `verify-closure-notes` fährt `--enable planning --enable structure --enable spans`. Isolierte Probe (§5): nur `structure` ⇒ `1 Datei(en) geprüft, 0 Befund(e)`, **Exit 0**; `structure` + `spans` ⇒ `fence-unclosed` auf der Öffnungszeile, **Exit 1**; Kontrolllauf mit geschlossenem Fence ⇒ `section-tasks-open` auf der Haken-Zeile | **erfüllt** — mit **V-6** (im Repo-Profil ist der Fall schon ohne `spans` rot; `spans` liefert die Diagnose, nicht die Sichtbarkeit — die Commit-Botschaft sagt das, `AGENTS.md` und `MR-056` sagen es nicht) |
| 3 | Retro gemessen: ohne Ausnahme 144 Befunde in 37 Dateien, mit Ausnahme null; **beide Ausgaben stehen in der Commit-Botschaft** | Selbst nachgefahren im echten Closure-Profil (§3): ohne Ausnahme `560 Datei(en) geprüft, 144 Befund(e)`, alle 144 mit Grund-Code `section-tasks-open`, **37** distinkte Dateien; mit Ausnahme `560 Datei(en) geprüft, 0 Befund(e)`, Exit 0. **Zahlen bestätigt.** Die Commit-Botschaft trägt die **Zahlen**, nicht die Ausgabe-Zeilen | **in der Sache erfüllt** — **V-8** (Form: Zahlen statt Ausgaben) und **V-2** (die Teil-Zahl „34 Stück" in §2 des Plans ist gemessen 33) |
| 4 | Die Bestands-Ausnahme ist als Adaption geführt und nennt den Grund: gehobenes Regelwerk, nicht nachgezogene Dokumente | `harness/conventions/MR-056-dod-haken-waechter.md` existiert, Status `Accepted`, Geltungsbereich benannt, Grund wörtlich genannt, vier Grenzen ausgeschrieben, Auflösungs-Trigger gesetzt. Index-Zeile in `harness/conventions.md` vorhanden. Die `d-check:cite`-Direktive ist **kein stiller Beleg**: eine veränderte Zitat-Zeile erzeugt `citation-mismatch` (§7) | **erfüllt** — mit **V-4** (die Index-Zeile trägt eine überzählige leere Spalte) und **V-3** (die Zahl „elf Closures" ist gemessen 8) |
| 5 | Die drei Folge-Kandidaten aus §3 sind als solche benannt | §3 des Plans nennt fünf Ausschlüsse, davon drei als Folge-Schnitt (Review-Report-Deckung, `BEO`-Deckung, Reihenfolge über git). `welle-86` §4 führt sie mit Kennung: slice-173, slice-174, slice-175 | **erfüllt** |
| 6 | Die neue Bedingung ist **zugestellt**, bevor sie blockiert: `AGENTS.md` §5 und die Sensors-Tabelle in `harness/README.md` nennen sie | `AGENTS.md` §5 trägt den Absatz *„Seit slice-172 hält das ein Sensor"* samt Bindepunkt, Ausnahme-Grenze und zwei Grenzen; `AGENTS.md` §4 trägt sie zusätzlich in der `verify-closure-notes`-Zeile; `harness/README.md` §Sensors trägt sie mit Grund-Code, Ausnahme und den gemessenen Zahlen | **erfüllt** — mit **V-5** (die abgeschlossene Zahl „zwei Grenzen" unterbietet die vier aus `MR-056`) und **V-7** (Platzierung in `harness/README.md`) |
| 7 | `make gates` grün (Exit explizit); **unabhängiger Review**; **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten | `make gates` → `[gates] … green`, **Exit 0**; zusätzlich `make fullbuild` → **Exit 0** (§2). Diese Datei ist die Verifikation, in eigenem Kontext gefahren. Der Review-Report `docs/reviews/2026-08-30-slice-172-dod-haken-review.md` ist während dieses Laufs im Arbeitsbaum erschienen — er lief in einem eigenen, parallelen Kontext; sein **Inhalt** ist hier bewusst nicht gelesen (Rollen-Trennung) | **erfüllt, soweit hier prüfbar** (**V-10**) |

Die Haken in §4 des Slice-Plans stehen noch offen, und die Risiko-Ausgänge in §5 tragen noch den Platzhalter — beides planmäßig: der Closure-Body ist ein eigener Commit nach dem Lifecycle-Move.

## 2. Gate-Läufe mit echter Ausgabe

Alle Läufe gegen `79bf375`, Arbeitsbaum sauber.

| Target | Ausgabe (Schlusszeilen) | Exit |
|---|---|---|
| `make verify-closure-notes` | `d-check: 560 Datei(en) geprüft, 0 Befund(e)` | **0** |
| `make doc-check` | `d-check: 623 Datei(en) geprüft, 0 Befund(e)` | **0** |
| `make gates` | `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` | **0** |
| `make fullbuild` | `50 Anforderung(en), 0 Waise(n).` · `d-check: 560 Datei(en) geprüft, 0 Befund(e)` · `[fullbuild] green — image-hash sha256:fb5c3b907001c300e9dd7cb11134ca2b10b9354a2c315ed7e436bdcc3b880f68` | **0** |
| `make trace-check RANGE=fd4cbb4^..HEAD` | `d-check: 623 Datei(en) geprüft, 0 Befund(e)` | **0** |
| `make adr-check RANGE=fd4cbb4^..HEAD` | `d-check: 623 Datei(en) geprüft, 0 Befund(e)` | **0** |

`make gate-consistency` läuft als Glied von `gates` (Modul `targets`, Ausgabe `623 Datei(en) geprüft, 0 Befund(e)`) — die Deklarations-Konsistenz Doku ↔ Makefile hält also über die neuen `AGENTS.md`-Zeilen hinweg. Der Bau ist unverändert (`sha256:fb5c3b90…` in allen Läufen).

## 3. Die Retro-Messung nachgefahren

Die Ausnahme wurde in einer **Kopie** des Baums aus der DoD-Regel entfernt, alles andere unverändert; gefahren wurde das echte Rezept (`--config .d-check.closure.yml --enable planning --enable structure --enable spans`).

| Lauf | Ausgabe | Exit |
|---|---|---|
| ohne Bestands-Ausnahme | `d-check: 560 Datei(en) geprüft, 144 Befund(e)` | **1** |
| mit Bestands-Ausnahme (unverändertes Repo) | `d-check: 560 Datei(en) geprüft, 0 Befund(e)` | **0** |

**144 Befunde, 37 distinkte Dateien, und alle 144 tragen denselben Grund-Code** `section-tasks-open` — kein einziges `section-missing`. Damit ist die Vorarbeits-Aussage des Plans, das Überschriften-Muster decke den Bestand vollständig, unabhängig belegt: alle **177** `done/`-Slices tragen eine passende DoD-Überschrift, in **sieben** distinkten Formen (`## 2.`, `## 3.`, `## 4.`, `## 5.` sowie die Zusätze *(vorläufig)* und *(R1 eingearbeitet)*).

Die Nummern der 37 Dateien, gemessen:

```text
025 026 027 028 040 041 042 043 044 045 046 048 049 050 051 052 053 054 055
057 064 065 066 067 068 076 082 094 097 098 101 102 104 160 168 169 170
```

Das sind **33** Dateien in der Spanne 025–104, dann 160, 168, 169, 170. Der Slice-Plan schreibt in §2 „025–104 (34 Stück)" — die Gesamtzahl 37 stimmt, die Teil-Zahl nicht (**V-2**). Die Commit-Botschaft, der Kopfkommentar und `MR-056` nennen die Teil-Zahl nicht und sind davon unberührt.

## 4. Positiv- und Negativ-Probe

**Außerhalb der Ausnahme** — ein Haken in `slice-181` geöffnet (Zeile 96), sonst nichts verändert:

```text
docs/plan/planning/done/slice-181-handbuch-link-in-cli.md:96	…	section-tasks-open	offener DoD-Haken in einem geschlossenen Slice: Haken setzen, wenn der Punkt erledigt ist — sonst gehört der Slice zurück nach in-progress/
d-check: 560 Datei(en) geprüft, 1 Befund(e)
```

Exit 1. Zurückgesetzt: `560 Datei(en) geprüft, 0 Befund(e)`, Exit 0. Datei, **Zeile des Hakens** und Hinweis stehen im Befund — die Zusage der DoD-Zeile 1 hält wörtlich.

**Innerhalb der Ausnahme** — zwei Proben:

- Die 37 Bestands-Dateien melden im unveränderten Repo **nichts** (Lauf oben, 0 Befunde), obwohl sie 144 offene Haken tragen.
- Ein **neu** geöffneter Haken in `slice-150` — einer Datei innerhalb der Ausnahme, die bisher keinen offenen Haken trug — läuft ebenfalls still durch: `560 Datei(en) geprüft, 0 Befund(e)`, Exit 0. Die Ausnahme ist also keine Liste der 37, sondern eine **Nummernspanne**; sie deckt auch künftige Öffnungen in Alt-Slices.

## 5. Die `spans`-Abhängigkeit — isoliert gemessen

Probe-Baum außerhalb des Repos: ein `done/`-Slice, in dessen DoD-Abschnitt ein Fence geöffnet und nie geschlossen wird, dahinter ein offener Haken. Konfiguration: **nur** die DoD-Regel.

| Lauf | Ausgabe | Exit |
|---|---|---|
| nur `structure` | `1 Datei(en) geprüft, 0 Befund(e)` | **0** |
| `structure` + `spans` | `…:5	```text	fence-unclosed` · `1 Datei(en) geprüft, 1 Befund(e)` | **1** |
| Kontrolle: derselbe Haken, Fence geschlossen, nur `structure` | `…:9	…	section-tasks-open	…` · `1 Datei(en) geprüft, 1 Befund(e)` | **1** |

Die Blindheit ist damit dem Fence zuzuschreiben, nicht der Probe. **Die Zusage der DoD-Zeile 2 hält.**

**Was `spans` dabei nicht tut:** es meldet den **Fence**, nicht den Haken. Der offene Haken bleibt in beiden Läufen ungemeldet; gewonnen ist die Sichtbarkeit des *Grundes*, nicht die der *Bedingung*.

**Und im echten Repo-Profil ist der Fall auch ohne `spans` rot.** Derselbe unbeendete Fence, in `slice-181` eingesetzt:

- ohne `spans`: **4 Befunde** — `closure-note-missing`, zweimal `section-missing`, `section-empty`. Exit 1.
- mit `spans`: dieselben vier **plus** `fence-unclosed` auf der Öffnungszeile. Exit 1.

Die Commit-Botschaft sagt das ausdrücklich („nur verdeckt … das ist Zufall, nicht Konstruktion"). `MR-056` und `AGENTS.md` sagen es nicht — dort liest sich die Grenze so, als sei der Bindepunkt ohne `spans` blind. Er ist es nicht; blind ist die **Regel** (**V-6**).

## 6. Die vier Grenzen aus `MR-056`, gegen das Produkt

| Grenze | Messung | Ergebnis |
|---|---|---|
| **Ein Haken ist eine Selbstauskunft** | Nicht maschinell prüfbar, und das ist keine Lücke der Umsetzung, sondern der Gegenstand: die Bedingung liest ein Zeichen. Gemessen ist nur die Folge — der Befund aus §4 verschwindet, sobald `[ ]` zu `[x]` wird, unabhängig davon, ob ein Report existiert. Die Deckung *Zusage ⇒ Report* ist ausdrücklich slice-173 | **zutreffend deklariert** |
| **Ein vergessener Schluss-Fence schaltet sie ab** | §5: isoliert 0 Befunde / Exit 0; gefangen wird der Fall von `fence-unclosed` | **zutreffend, mit V-6** |
| **Sie prüft den Zustand, nicht den Übergang** | Belegt an drei Stellen: (a) `.githooks/pre-commit` fährt `make adr-check STAGED=1` und `make doc-check` — **nicht** das Closure-Profil; (b) derselbe geöffnete Haken aus §4 im **`doc-check`-Profil** gemessen: `623 Datei(en) geprüft, 0 Befund(e)`, Exit 0 — der Hook sieht ihn also nicht; (c) `ci.yml` ruft `make trace-check`, `make adr-check` und `make ci` (gates + image-test), keines davon `verify-closure-notes`. Der `mv`-Commit passiert Hook **und** CI; gesehen wird der Haken erst beim nächsten `make fullbuild` bzw. `make verify-closure-notes` | **zutreffend deklariert, belegt** |
| **Die Ausnahme altert** | Gemessen: die drei Muster nehmen **169** der 177 `done/`-Slices heraus; nur **37** davon tragen heute einen offenen Haken. **132** Dateien sind also ausgenommen, ohne dass es dafür einen Bestandsbefund gäbe — und die `slice-150`-Probe aus §4 zeigt, dass eine **neue** Öffnung dort still durchläuft. Unter der Regel stehen heute **8** Dateien: slice-171, 176, 177, 178, 179, 180, 181, 182 | **zutreffend deklariert, und die Reichweite ist größer als der Anlass** |

## 7. Die Ausnahme-Globs, elementweise

Je ein Lauf mit **genau einem** der drei Muster, sonst identisch:

| Muster | verbleibende Befunde | verbleibende Dateien | herausgenommen |
|---|---|---|---|
| `slice-0??-*` | 11 | 7 (101, 102, 104, 160, 168, 169, 170) | **30** |
| `slice-1[0-6]?-*` | 134 | 31 | **6** (101, 102, 104, 160, 168, 169) |
| `slice-170-*` | 143 | 36 | **1** (170) |
| alle drei | **0** | 0 | 37 |

Die drei Beiträge sind **disjunkt** (30 + 6 + 1 = 37) und decken die Trefferliste **exakt** — kein Muster ist redundant, keines fehlt. Die Aussage aus `MR-056`, es brauche drei statt zwei Muster, ist damit belegt.

**Nimmt eines mehr heraus als gemeint?** Ja, und zwar bewusst: als Nummernspanne nehmen die drei zusammen `slice-001` bis `slice-170` heraus, also 169 Dateien statt der 37 mit Befund (§6, vierte Zeile). Über die Spanne hinaus greift keines — die Muster sind auf drei Ziffern und den Bindestrich fixiert, und die `files`-Menge ist ohnehin auf `done/slice-*.md` beschränkt. Die aus `MR-049` bekannte Falle (`slice-1[0-3]*` nähme auch `slice-1000` heraus) ist hier vermieden.

**Gegenprobe zur `cite`-Direktive:** eine geänderte Zitat-Zeile in `MR-056` erzeugt

```text
harness/conventions/MR-056-dod-haken-waechter.md:5	…modul-05-planning-harness.md:33-34	citation-mismatch	Zitattext ist kein zusammenhängender Teilstring der Quell-Spanne (Zitat-Fäule)
```

Der Beleg ist also nicht still.

## 8. Zusätzliche Produkt-Proben zur Bedingung selbst

Alle gegen den Probe-Baum, `structure` + `spans`, nur die DoD-Regel.

| Probe | Ergebnis |
|---|---|
| **Unbalancierter Backtick im selben Absatz**, Haken dahinter — die Blindstelle, an der der zweite Bau scheiterte | Haken **gemeldet** auf seiner Zeile, dazu `span-unclosed`. Die Begründung für `max-open-tasks` statt `forbid-pattern` hält |
| **Alle vier Listen-Marker** (`-`, `*`, `+`, `1.`), je offen | **4 Befunde**, je auf eigener Zeile. Die zweite Blindstelle des zweiten Baus ist geschlossen |
| **Absatzweite Inline-Code-Spanne** um den echten Haken (Backtick auf Zeile 5, zu auf Zeile 7, Haken auf 6) | Haken **gemeldet**. Die vierte Grenze aus ADR-0059 gilt für diese Bedingung **nicht** |
| Haken **in Inline-Code auf derselben Zeile** | kein Befund — richtig, das ist Dokumentation über die Form |
| Haken **innerhalb eines geschlossenen Fence** | kein Befund — richtig, gleiche Begründung |
| **Balanciertes Streu-Fence-Paar**, das die DoD-Überschrift verschluckt | **`section-missing`**, Exit 1 — die Regel wird laut, statt still zu werden. Aber die Meldung trägt den DoD-Hinweis-Text (**V-9**) |

Die dritte Zeile ist die wichtigste: sie widerlegt die Lesart, die die Platzierung in `harness/README.md` nahelegt (**V-7**).

## 9. Der ersetzte Kopfkommentar — überschreibt er eine frühere Messung?

**Nein, und das ist nachgerechnet.** Der abgelöste Kommentar stammt aus `e93d6a9` (slice-099, 2026-08-15) und führte die Regel als bewusst verworfen mit *„32 Befunde — und die sind BERECHTIGT"*. Ich habe den Baum dieses Commits ausgepackt und die **heutige** Bedingung dagegen gefahren:

| Baum | done/-Slices | Dateien mit offenem DoD-Haken |
|---|---|---|
| `e93d6a9` (slice-099) | 99 | **32** |
| `79bf375` (heute) | 177 | **37** |

Die Zahl 32 ist also reproduzierbar und stimmt. Die Differenz sind genau fünf **neu** hinzugekommene Dateien (102, 160, 168, 169, 170); **keine** der damals gezählten 32 ist seither still repariert worden. Der Bestand ist gewachsen, die Messung nicht umgedeutet.

**Was der neue Kommentar nicht sagt:** dass sich die Einheit geändert hat. „32 Befunde" waren **Dateien** (die alte Form meldete je Abschnitt), „144 Befunde in 37 Dateien" sind **Haken**. Wer die beiden Zahlen nebeneinander liest, vergleicht zwei Größen. Der Kommentar zitiert die alte Aussage und erklärt sie für überholt — er benennt aber nicht, dass 32 und 37 dasselbe messen und 144 etwas anderes.

**Die beiden namentlich genannten Fälle sind heute korrekt behandelt:**

- `slice-094` Zeile 141: `- [ ] **Release als Minor** — **Wellen-Trigger, nicht Slice-Trigger.**`
- `slice-104` Zeile 120: `- [ ] **Release** — **Wellen-Trigger**, gemeinsam mit …`

Beide liegen in der Bestands-Ausnahme (`slice-0??-*` bzw. `slice-1[0-6]?-*`) und werden nie gemeldet. Die Form-Regel in `AGENTS.md` §5, die solche Punkte künftig in den Wellen-Closure-Trigger schiebt, trägt eine ausdrückliche Bestands-Grenze für genau diese Dateien. Alt-Beleg und neue Regel widersprechen sich also nicht.

## 10. Plan-vs-Code-Diff

**Geliefert = geplant.** Der feat-Commit fasst fünf Dateien an: `.d-check.closure.yml`, `AGENTS.md`, `harness/README.md`, `harness/conventions.md`, `harness/conventions/MR-056-dod-haken-waechter.md`. Das entspricht §2 Punkt 1–5 des Plans.

**§3 (Out-of-Scope) ist eingehalten.** Keine Review-Report-Deckung, keine `BEO`-Deckung, keine Reihenfolge-Prüfung über git, kein Nachrüsten der 37 Bestands-Slices, keine Aussage über Review-Qualität. Der Claim-Move-Commit fasst zwar zwei `done/`-Slices an — die Änderungen sind ausschließlich Pfad-Verweise `open/` → `in-progress/`, kein DoD-Haken ist berührt (nachgesehen im Diff).

**§5 (Risiken) — alle vier stehen noch auf dem Platzhalter** *(bei Closure)*. Das ist planmäßig; der Closure-Body ist ein eigener Commit. Nach heutigem Stand:

| Risiko | Stand nach dieser Messung |
|---|---|
| Haken ist Selbstauskunft | **weiter offen** — nicht behebbar, deklariert; Deckung ist slice-173 |
| Bestands-Ausnahme altert | **weiter offen** — und die Messung in §6 beziffert sie: 169 ausgenommene Dateien, 132 davon ohne Befund |
| Sensor macht den inneren Loop rot | **entfallen** — gemessen: die Regel liegt im Closure-Profil, `make gates` (623 Dateien) meldet sie nicht; `make doc-check` mit einem offenen Haken bleibt bei 0 Befunden |
| Regel aus dem Anlass gezogen | **weiter offen als Notiz** — die Regel selbst ist mit 37 Bestandsfällen belegt, die Dringlichkeit stammt aus drei Slices |

**§7 (Vorgelagert) — ein Block ist veraltet.** Die Beobachtungs-Sichtung deklariert *„Register-Stand 2026-08-29, höchste Kennung `BEO-023`"*. Das Register führt `BEO-024` seit `b35662e` vom **2026-08-29 18:37** — also bereits vor dem Prep-Commit `fd4cbb4` vom **2026-08-30 17:43**, mit dem die dritte Beanspruchung §2 und §7 neu geschrieben hat. Der Nachtlauf-Block wurde bei dieser Beanspruchung aktualisiert, der Sichtungs-Block nicht (**V-1**).

## 11. Deklarations-Konsistenz zwischen den vier Zustell-Orten

| Aussage | `.d-check.closure.yml` | `AGENTS.md` §4 | `AGENTS.md` §5 | `harness/README.md` | `MR-056` |
|---|---|---|---|---|---|
| 144 Befunde / 37 Dateien | ja | — | — | ja | — |
| Ausnahme bis `slice-170` | ja | ja | ja | ja | ja |
| Bindepunkt ist `verify-closure-notes`, nicht `gates` | ja | ja | ja | ja | ja |
| Grund-Code `section-tasks-open` | ja | ja | ja | ja | — |
| Zustand statt Übergang | — | ja | **nein** | **nein** | ja |
| Ausnahme altert | — | **nein** | **nein** | **nein** | ja |
| Zahl der Grenzen | — | „zwei" | „zwei" | „zwei" (implizit) | **vier** |

Keine Zahl widerspricht einer anderen. Die Zusagen unterscheiden sich in der **Reichweite**: `AGENTS.md` §5 und `harness/README.md` schließen die Grenzen-Aufzählung mit „zwei", während `MR-056` vier führt (**V-5**).

## 12. Befunde

| ID | Schwere | Befund |
|---|---|---|
| **V-1** | **MEDIUM** | §7 des Slice-Plans deklariert *„Register-Stand 2026-08-29, höchste Kennung `BEO-023`"*; das Register führte `BEO-024` bereits seit dem Vorabend (`b35662e`, 2026-08-29 18:37), also vor dem Prep-Commit der dritten Beanspruchung. Der Sichtungs-Block ist aus der zweiten Beanspruchung übernommen statt neu gelesen — genau die Klasse, gegen die `MR-054` die `cite`-Pflicht setzt, nur an dem Block, der bewusst keine trägt. Erschwerend: `BEO-024` handelt von **Zustell-Kanälen**, also von der Frage, die DoD-Punkt 6 dieses Slice beantwortet |
| **V-2** | **LOW** | §2 des Slice-Plans schreibt „025–104 (34 Stück)"; gemessen sind es **33**. Die Gesamtzahl 37 stimmt. Kein anderes Dokument trägt die Teil-Zahl |
| **V-3** | **LOW** | `MR-056` schreibt „Gemessen halten die Slices 171–182 ihre Haken gesetzt — **elf** Closures ohne einen einzigen offenen". In `done/` liegen in dieser Spanne **acht** Dateien (171, 176–182); 172 ist in Arbeit, 173–175 existieren nicht. Die Aussage *„keine offenen Haken"* stimmt, die Zahl nicht |
| **V-4** | **LOW** | Die `MR-056`-Zeile in `harness/conventions.md` (Zeile 136) endet mit einer überzähligen leeren Spalte: sechs Trennzeichen gegen fünf im Tabellenkopf und in allen Nachbarzeilen. Kein Gate fängt das (`make doc-check` grün); die Zeile rendert mit einer leeren fünften Spalte |
| **V-5** | **LOW** | `AGENTS.md` §4 und §5 sowie `harness/README.md` schließen die Grenzen-Aufzählung mit „**Zwei** Grenzen"; `MR-056` führt **vier**. Wer nur §5 liest, erfährt weder, dass die Regel den Übergang nicht bindet, noch dass die Ausnahme altert. Eine geschlossene Zahl an der Zustell-Fläche unterbietet die kanonische Adaption |
| **V-6** | **LOW** | Die Fence-Grenze ist an der Zustell-Fläche schärfer formuliert als gemessen: `AGENTS.md` und `MR-056` sagen, ein vergessener Schluss-Fence „schaltet sie ab" — gemessen ist der **Bindepunkt** in diesem Repo trotzdem rot (vier Befunde aus Nachbarregeln), blind ist nur die **Bedingung**. Die Commit-Botschaft trägt die Unterscheidung, die Dokumente nicht |
| **V-7** | **LOW** | In `harness/README.md` steht der neue DoD-Satz **vor** dem Absatz „Vier benannte Grenzen … Inline-Code wird absatzweise gepaart". Gemessen gilt diese vierte Grenze für `max-open-tasks` **nicht** (§8, dritte Zeile): eine absatzweite Spanne verschluckt den Haken nicht. Die Platzierung lädt zu einer Lesart ein, die das Produkt widerlegt |
| **V-8** | **INFO** | DoD-Punkt 3 verlangt, dass „beide **Ausgaben**" in der Commit-Botschaft stehen. Die Botschaft nennt die **Zahlen** („144 Befunde in 37 Dateien … mit ihr NULL"), keine Ausgabe-Zeile. Sachlich vollständig, formal knapp an der eigenen Formulierung vorbei |
| **V-9** | **INFO** | Der `hint` hängt an der **Regel**, nicht am Grund-Code (dokumentierte Produkt-Eigenschaft, Benutzerhandbuch §4). Emittiert dieselbe Regel `section-missing` — gemessen möglich, wenn ein Streu-Fence-Paar die DoD-Überschrift verschluckt —, bekommt der Leser den Reparatur-Hinweis für den **falschen** Defekt: „Haken setzen … sonst gehört der Slice zurück". Im heutigen Bestand tritt der Fall nicht auf (0 `section-missing` über 177 Dateien) |
| **V-10** | **INFO** | Beim ersten Sensor-Lauf dieser Verifikation lag kein Review-Report zu slice-172 vor (`make doc-check`: 623 Dateien); beim abschließenden Lauf lag er da (625 Dateien). Beide Hälften von DoD-Punkt 7 existieren damit, in getrennten Kontexten. Der Inhalt des Reviews ist hier nicht gelesen — ein Verifier, der die Review-Findings übernimmt, prüft nicht mehr unabhängig |
| **V-11** | **INFO** | Nicht von diesem Slice eingeführt, aber von ihm berührt: die Datei-Zahlen in `harness/README.md` („546 Dateien des Bindepunkts … seiner 608", „über 674 Dateien gemessen") sind veraltet — gemessen sind es heute **560** und **623** |

## 13. Verdikt

**Konform.**

Der gelieferte Wächter tut, was der Slice über ihn behauptet, und die tragenden Behauptungen halten der Nachmessung stand — jede einzeln nachgefahren, nicht nachgelesen:

- **144 Befunde in 37 Dateien ohne Ausnahme, null mit ihr** — im echten Closure-Profil reproduziert, nicht in einer Probe-Konfiguration.
- **Die Positiv-Probe trägt** — Datei, Zeile des Hakens, Grund-Code und verfasster Hinweis; zurückgesetzt wieder grün.
- **Die `spans`-Abhängigkeit ist gemessen, nicht angenommen** — isoliert 0 Befunde/Exit 0 ohne, `fence-unclosed` mit.
- **Die drei Ausnahme-Muster decken die Trefferliste exakt und disjunkt** (30 + 6 + 1).
- **Die abgelöste Messung ist nicht überschrieben, sondern nachgerechnet**: 32 Dateien im Baum von slice-099, 37 heute, Differenz sind fünf Neuzugänge — keine stille Reparatur.
- **Alle Gates grün mit echter Ausgabe**, `make fullbuild` eingeschlossen.

Die elf Befunde sind sämtlich **Dokumentations- und Zahlen-Präzision**, keiner betrifft das Verhalten der Bedingung. Zwei verdienen vor der Closure einen Ausgang: **V-1** (der Sichtungs-Block ist eine Deklaration ohne Sensor, und er ist hier nachweislich stehengeblieben — das ist die `BEO-022`/`BEO-024`-Familie in eigener Sache) und **V-5** (eine geschlossene Grenzen-Zahl an der Zustell-Fläche, die die kanonische Adaption unterbietet — der Slice liefert Zustellung als DoD-Punkt und sollte sie nicht knapper zustellen als sie ist).

**Was dieser Bericht nicht sagt:** ob ein gesetzter Haken einen Review belegt. Das prüft die Bedingung nicht, das prüft niemand, und der Slice sagt das selbst. Die Lücke ist von *unsichtbar* nach *behauptet* gewandert — mehr ist hier nicht geliefert und mehr war nicht versprochen.
