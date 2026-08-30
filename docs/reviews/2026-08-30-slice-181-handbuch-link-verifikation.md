# Verifikation slice-181 — der Handbuch-Zeiger in `--help` und `--print-mk`

**Gegenstand:** `84f7259~1..50d461e` (5 Commits). HEAD ist `ff21044` (eingehender CR, **kein** Prüfgegenstand).
**Rolle:** unabhängiger Verifier (gegen DoD und Spec, nicht Code-Review).
**Gefahrene Sensors (alle lesend):** `make baseline-verify` · `make arch-check` · `make semgrep` · `make nightly-state` · direkte `docker run`-Läufe von `d-check:latest` für `doc-check`, `planning-check`, `gate-consistency`, `workflow-pins`, `adr-check RANGE=…`, `trace-check RANGE=…`, `completeness-check`, Closure-Profil · Docker-Build-Stages `test`/`lint`/`coverage` **gegen eine Scratch-Kopie** des Baums (eigene Tags, inzwischen entfernt).
**Nicht gefahren:** `make gates`, `make record-gates`, `make fullbuild` (Auftrag).

## 1. DoD-Tabelle

| # | Behauptet (§4) | Gemessen | Verdikt |
|---|---|---|---|
| 1 | `--help` nennt das Handbuch mit voller URL; `--print-mk` trägt sie im Kopfkommentar — **gegen das echte Binary** | `docker run --rm --network none d-check:latest --help` → **Exit 0**, stderr endet mit `Benutzerhandbuch (aufgabenorientiert, deutsch):` + `https://github.com/pt9912/d-check/blob/main/docs/user/benutzerhandbuch.md`. `--print-mk` → **Exit 0**, stderr leer, Zeilen 11–12 der Ausgabe tragen dieselbe URL, **vor** `DCHECK_IMAGE ?=` (Zeile 13). Je **genau ein** Vorkommen. Zwei Läufe byte-identisch (`sha256 e0539a8e…`). **Kein Mount nötig** — Repo-Zugriffs-Freiheit mitgemessen | **erfüllt** |
| 2 | URL zeigt auf `blob/main`, ohne Versionsangabe; ein Release muss sie nicht anfassen | URL enthält `/blob/main/`, keine Version. Release-Prep (`docs/user/releasing.md:24-100`) nennt vier Bump-Flächen; keine trifft die URL. **Gegenprobe gefahren** (§3, P1/P2): eine *versionierte* Fassung derselben URL im Handbuch bleibt über **alle elf** Default-Module **still** (613/0, Exit 0), während ein veralteter `ghcr`-Pin an derselben Datei `version-stale` meldet. Die tragende Aussage ist damit **gemessen**, nicht nur zitiert | **erfüllt, stärker belegt als zugesagt** |
| 3 | Beide Akzeptanzkriterien fordern den Zeiger; Lastenheft-Bump samt Historie | `spec/lastenheft.md:98` (*Hilfe*) und `:516` (*Happy Path*) fordern ihn; Version 0.76.0 → **0.77.0**, Historie-Zeile `:3343`. `doc-check` grün (Chronologie-Monotonie der Historie inbegriffen) | **erfüllt** — Einschränkungen: **A-2**, **A-3** (die `.a`-Verfeinerungen sind mitbetroffen und nicht nachgezogen), **A-7** (eine Aussage der Historie trägt weiter, als ihre Quelle reicht) |
| 4 | Je ein Test hält die Zusage; jeder wird von **genau** seiner Mutation rot; Vorzustand mitgeprüft | Fünf Mutationen + Vorzustands-Probe gefahren (§3). M1 → 2 rot, M2 → 2 rot, M3 → 3 rot — **exakt die Zahlen der Commit-Botschaft**. **V0** (neue Testdatei gegen den Code von `ca43027~1`): **alle 3 rot** ⇒ Vorzustand belegt. **M3b** (beide Literale versioniert): **genau 1 rot** — der Versions-Wächter ist spezifisch. **M5** (Zeile unter `DCHECK_IMAGE` verschoben): **genau 1 rot** — der Positions-Anker greift | **erfüllt** — Präzisierung **A-10** |
| 5 | Handbuch §4.16 zeigt den Kopf mit der neuen Zeile; `--print-mk` bleibt `not-replayable` | `docs/user/benutzerhandbuch.md:845-852`: der Fence zeigt jetzt `# … Einbinden, Digest-Pin` / `#` / die zwei Handbuch-Zeilen / `DCHECK_IMAGE ?=`. Die drei nicht elidierten Zeilen sind **wortgleich** mit der echten Ausgabe (Zeilen 10–13). Marker `:844` steht unverändert unmittelbar vor dem Fence; `TestHandbook_OutputBlocksClassified` grün | **erfüllt** — Einschränkungen **A-5** (die Begründung im Commit trifft nicht zu) und **A-6** (die *umgebende, mitbearbeitete* Prosa ist falsch) |
| 6 | `make gates` grün; unabhängiger Review; Verifikation — in eigenen Kontexten | Auftragsgemäß kein `gates`-Lauf. **Zehn der elf Glieder einzeln gefahren, alle grün:** `baseline-verify` (51/ok) · `workflow-pins` (613/0) · `doc-check` (613/0) · `lint` (**0 issues**) · `test` (alle Pakete `ok`) · `arch-check` (0) · `coverage-gate` (**94,70 % ≥ 93 %**) · `semgrep` (55 Regeln, 0) · `gate-consistency` (613/0) · `planning-check` (613/0). Nicht gefahren: `record-gates` | **offen (planmäßig)** — dies ist die Verifikation |

## 2. Nachgeprüfte Zahlen

| Zahl | Fundstelle | Behauptet | Gemessen |
|---|---|---|---|
| Scan-Menge | Commits `ca43027`, `50d461e` | **612 Dateien, 0 Befunde** | heute **613/0** — die Differenz ist genau `docs/plan/cr/2026-08-30-cr-a-check-leermenge.md` aus `ff21044` (außerhalb des Slice). 613 − 1 = **612**, rekonstruiert konsistent |
| „beide Ausgaben trugen **null** URLs" | slice §1, Lastenheft-Historie, Commit | 0 | **0** — `git show ca43027~1:…/cli.go` (`writeUsage`-Rumpf) und `…/print_mk.go`: **kein** `http`/`://`. Auch die volle `--help`-Ausgabe (inkl. `flag.PrintDefaults`) trug keine. Die Menge ist in §7 ausdrücklich auf diese beiden Ausgaben eingegrenzt — `--print-config` trägt zwei Beispiel-URLs (`config_template.go:239,241`); die Eingrenzung ist korrekt und trägt |
| „verwiesen nur auf andere Ausgaben" | slice §1 | `--print-config`/`--suggest-config` bzw. Release-Notes | bestätigt: Hilfe-Trailer vorher genau diese zwei; `mkTemplate`-Kopf: *„den Digest aus den **Release-Notes** pinnen"* |
| „`versions` hält ausschließlich `ghcr`-präfixierte Pins" | slice §2.2, `cli.go:124`, Test-Kommentar, Lastenheft `:98`/`:3343`, Commit | — | **für dieses Repo wahr**: `.d-check.yml:398` `pin-pattern: 'ghcr\\.io/[^\\s:]+:(v…)'`. **Nicht** Eigenschaft des Moduls: das Produkt hat **keinen** verdrahteten Default (`configyaml.go:130`), `DC-FA-VER-001` sagt „des **konfigurierten** Musters (z. B. …)". Siehe **A-7** |
| „zwei benannte Stellen im Handbuch" (nackte Tags) | slice §2.2, Lastenheft-Historie, Commit | **2** | **exakt 2** — `docs/user/benutzerhandbuch.md:82` (`docker pull pt9912/d-check:v0.68.0`) und `:87` (`` `:v0.68.0` ``); beide stehen namentlich in `docs/user/releasing.md:43-46`. **Erschöpfend nachgezählt:** jedes weitere `v0.68.0` im Handbuch ist entweder `ghcr`-präfixiert (Gate deckt), der Header-Link `:3` (siehe unten) oder eine §11-Historienzeile. **Beide Hälften mutativ belegt** (§3, P3/P4) |
| „`blob/main` … dieselbe Form, die die Docker-Hub-Overview schon nutzt" | slice §2.2 | — | **byte-identisch**: `packaging/dockerhub/overview.md:62` trägt exakt dieselbe URL |
| Nachtlauf-Stand (`MR-053`) | slice `:156-160` | beide `gruen`, `2026-08-30T06:08:17Z` / `2026-08-30T09:16:25Z` | `make nightly-state`: **byte-genau dieselben Zeitstempel**, beide `gruen`, Exit 0 |
| Register-Stand | slice §7 | höchste Kennung `BEO-024`; `BEO-020` = 4; `BEO-023` = 4 | alle drei bestätigt (`observations.md`) |
| `##`-annotierte Targets in `--print-mk` | Handbuch `:839`, Spezifikation `:692` | **elf** | **zwölf** (`doc-structure` fehlt in beiden Listen). Alle 11 `docker run`-Rezepte tragen `--network none` **und** `:/repo:ro`; `doc-help` läuft ohne Docker. Siehe **A-6** |

## 3. Mutations- und Kontrollproben

Alle Proben gegen eine Scratch-Kopie des Baums (`internal/` byte-identisch verifiziert), Stage `test` mit `--no-cache-filter`. Basislauf **grün**.

| Probe | Was geändert | Rot | Welche |
|---|---|---|---|
| **M1** | die zwei `Fprintln`-Zeilen des Handbuch-Trailers aus `writeUsage` entfernt | **2** | `TestCLI001_HilfeNenntHandbuch`, `TestHandbuchURL_TraegtKeineVersion` |
| **M2** | die zwei Template-Zeilen aus `mkTemplate` entfernt | **2** | `TestCLI010_PrintMKNenntHandbuch`, `TestHandbuchURL_TraegtKeineVersion` |
| **M3** | `blob/main` → `blob/v0.68.0`, **nur** im Produktions-Const | **3** | alle drei |
| **M3b** *(schärfer als zugesagt)* | dieselbe Versionierung **in beiden** Literalen (Produktion **und** Test) | **1** | **nur** `TestHandbuchURL_TraegtKeineVersion` |
| **M5** *(zusätzlich)* | Handbuch-Zeilen **unter** `DCHECK_IMAGE ?=` verschoben | **1** | **nur** `TestCLI010_PrintMKNenntHandbuch` |
| **M4 — Kontrolle** | in **beiden** Literalen auf ein **fremdes Repo/Pfad**, weiter `blob/main` | **0** | — (siehe **A-9**) |
| **V0 — Vorzustand** | neue Testdatei gegen `cli.go`/`print_mk.go` aus `ca43027~1` | **3** | alle drei |

M3b und M5 sind der eigentliche Beleg: sie zeigen, dass der Versions-Wächter und der Positions-Wächter **je einzeln** greifen und nicht nur als Kollateral eines geänderten Literals. M4 zeigt die Grenze der Trias.

**Gate-Gegenproben am echten Baum** (Bind-Mount-Überlagerung **einer** Datei, alle elf Default-Module):

| Probe | Ergebnis |
|---|---|
| **P0** Basis | 613 Dateien, **0 Befunde**, Exit 0 |
| **P1** Handbuch-URL in der §4.16-Illustration auf `blob/v0.60.0` | **0 Befunde, Exit 0** — kein Gate deckt eine versionierte URL |
| **P2 — Positivkontrolle** `ghcr`-Pin `:67` auf `v0.60.0` | **`version-stale` :67, Exit 1** — das Vorher-Grün in P1 ist kein taubes Grün |
| **P3** die zwei nackten Tags `:82`/`:87` auf `v0.60.0` | **0 Befunde, Exit 0** — die zwei „benannten Stellen" driften tatsächlich still |
| **P4** Header-Stempel `:3` auf `[v0.60.0](…#v0.60.0)` | **`anchor-missing` :3, Exit 1** — der Header ist gedeckt und gehört zu Recht **nicht** zu den zwei |

## 4. Spec-Konformität

**`DC-FA-CLI-001` *Hilfe* — als Verhalten erfüllt.** Gemessen gegen `d-check:latest`: Exit 0 · Ausgabe auf **stderr** · Synopsis `d-check [optionen] [pfad]` · Pfad-Argument-Zeile · Verweis auf `--print-config` · **URL des Benutzerhandbuchs** · `/blob/main/`, keine Version.

**`DC-FA-CLI-010` *Happy Path* — als Verhalten erfüllt.** `DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v0.0.0-dev` (lokaler Build-Default, spec-konform) · `DCHECK_DIGEST` · `TRACE_FLAGS` · **zwölf** `##`-annotierte Targets inkl. `doc-structure` und `doc-help` · alle 11 Docker-Rezepte mit `--network none` und `:/repo:ro` · Exit 0 · **ohne Mount lauffähig** · deterministisch · **Kopfkommentar nennt die URL**.

**Die Verfeinerungen gibt es — und sie sind nicht nachgezogen.** `spec/spezifikation.md` führt beide Pendants:
- `DC-FA-CLI-001.a` (`:49-56`) zählt die Usage-Ausgabe **„(in dieser Reihenfolge)"** als fünf Elemente auf; ausgeliefert werden jetzt **sechs**. → **A-2**
- `DC-FA-CLI-010.a` Punkt 1 (`:683-684`) zählt den Kommentar-Kopf als „Einbindung via `include`, Hinweis zum Digest-Pin"; er trägt jetzt ein **drittes** Element. → **A-3**

Die Spezifikation ist im ganzen Range **nicht angefasst** (`git diff --name-status`). Der Slice-Kopf *Berührte Spec-Stellen* nennt nur `DC-FA-CLI-001` und `DC-FA-CLI-010` — die `<DC-ID>.<Buchstabe>`-Form ist nach `AGENTS.md` §5 ausdrücklich zulässig und wäre hier fällig gewesen.

**Abwärts-Sperre (§3.4) eingehalten:** der neue Lastenheft-Text nennt keine ADR, Welle, Slice-ID und keinen Commit-Hash; `matrix` grün. Dass er über *Release-Prep dieses Repos* argumentiert, ist Prozess-Prosa in einem Produkt-Vertrag — keine Regelverletzung, aber die Wurzel von **A-7**.

**Traceability:** `trace-check` über den ganzen Range grün; `adr-check` grün (keine ADR berührt — konsistent mit §3 „Keine ADR"). `completeness-check` Exit 0.

## 5. Abweichungen Zusage ↔ Zustand

**A-1 — `internal/adapter/driving/cli/cli.go:116-127` (MEDIUM).**
Der Doc-Kommentar *„`writeUsage` gibt die Hilfe aus (DC-FA-CLI-001.a): …"* hängt jetzt an **`const handbuchURL`**, nicht an `writeUsage` — zwischen Kommentar und `const` steht keine Leerzeile, zwischen `const` und `func` ebenfalls nicht. Gemessen mit `go doc -all -u`: der Const trägt die writeUsage-Prosa, und `func writeUsage(flags *flag.FlagSet)` erscheint **ohne jede** Dokumentation. `make lint` meldet **0 issues** — kein Gate hält das. Nach `AGENTS.md` §3.7 beschreibt ein Kommentar, *was da ist*; hier beschreibt er, was **darunter** steht.

**A-2 — `spec/spezifikation.md:49-56` (MEDIUM).**
`DC-FA-CLI-001.a` zählt die Usage-Ausgabe „(in dieser Reihenfolge)" abschließend als fünf Elemente auf. Das Binary liefert ein sechstes. Die Verfeinerung ist unvollständig und weder aktualisiert noch im Slice-Kopf benannt.

**A-3 — `spec/spezifikation.md:683-684` (MEDIUM).**
`DC-FA-CLI-010.a` Punkt 1 beschreibt den Kommentar-Kopf abschließend mit zwei Bestandteilen; es sind jetzt drei. Dasselbe Bild wie A-2 — und das Fragment ist das Artefakt, das in **fremde** Repos reist, also die Stelle, an der eine unvollständige Spec am teuersten ist.

**A-4 — `spec/lastenheft.md:3343`, Slice `:55-59`, Commit `ca43027` (MEDIUM).**
*„**Keine Sprach-Kennzeichnung:** … sind Hilfe und Fragment **selbst deutsch**."* Gemessen tragen **beide** Ausgaben wörtlich `Benutzerhandbuch (aufgabenorientiert, **deutsch**)`. Der Zeiger führt also eine Sprach-Kennzeichnung — nur in anderer Form als das `(German)` der README. Die Entscheidung ist vertretbar; die Aussage über sie ist falsch, und sie steht jetzt in der **Lastenheft-Historie**.

**A-5 — Commit `50d461e`, Botschaft (MEDIUM).**
*„Die Illustration in Paragraph 4.16 zeigte den Kopf des erzeugten d-check.mk ohne die neue Zeile — seit dem Feat-Commit falsch."* Gemessen an `50d461e~1`: die Illustration zeigte **keinen Kopf**; ihr Fence begann bei `DCHECK_IMAGE ?=`. Eine vollständig elidierte Passage kann nicht falsch werden — der Commit **ergänzt** den Kopf, er repariert nichts. §2.5 des Plans formuliert die Bedingung korrekt („eine Illustration, **die den Kopf zeigt**"); die Commit-Botschaft behauptet die Bedingung als Tatsache. Klasse `BEO-020`: der Schluss reicht weiter als die Messung.

**A-6 — `docs/user/benutzerhandbuch.md:839-842` und `spec/spezifikation.md:692` (MEDIUM, Bestand — in der *mitbearbeiteten* Zeile).**
Beide sagen **elf** `##`-annotierte Targets und zählen elf Namen ohne `doc-structure`; das Binary emittiert **zwölf**. Der Satz im Handbuch ist **genau der**, den `50d461e` verlängert hat. Bemerkenswert: `spec/lastenheft.md:3377` (Historie 0.57.1) protokolliert, dass exakt diese Drift schon einmal *„nach unabhängigem Review, vor dem Release"* in den Akzeptanzkriterien saniert wurde — *„Dieselbe Stelle war schon in 0.37.1 als Selbstwiderspruch saniert worden."* Es ist damit die **dritte** stehende Instanz derselben Klasse, in zwei Dokumenten, die die Sanierung von 0.57.1 nicht erreicht hat.

**A-7 — `spec/lastenheft.md:98` und `:3343`, `cli.go:124`, `cli_handbuch_link_test.go:55-58`, Slice `:49-52`, Commit `ca43027` (MEDIUM).**
*„`DC-FA-VER-001` hält ausschließlich `ghcr`-präfixierte Pins."* Gemessen: `DC-FA-VER-001` hält, was `versions.pin-pattern` sagt; das Produkt hat **keinen** verdrahteten Default, und das `ghcr`-Muster steht in `.d-check.yml:398` — der **Konfiguration dieses Repos**. Der Satz steht jetzt als **Akzeptanzkriterium** im Lastenheft und lässt damit den Produkt-Vertrag eine Eigenschaft der Repo-Konfiguration behaupten. `AGENTS.md` §5 / `BEO-012`: *eine zitierte Quelle trägt nur, was in ihrem Geltungsbereich steht* — der Geltungsbereich ist hier `.d-check.yml`, nicht `DC-FA-VER-001`. **Die Entscheidung bleibt richtig** (P1 belegt sie unabhängig vom Zitat); nur die Begründungskette ist an einer Stelle zu weit gefasst.

**A-8 — `CHANGELOG.md` (MEDIUM).**
Im ganzen Range nicht angefasst; es gibt keinen `## [Unreleased]`-Abschnitt. Die Änderung ist **nutzersichtbar** (Hilfe-Ausgabe; der Kopf eines Fragments, das Adopter in ihre Repos legen). `AGENTS.md` §5 verlangt Pflege bei nutzersichtbaren Änderungen, und die gelebte Form ist der Feat-Commit (`e4b5c99` legte `## [Unreleased]` genau so an). §3 des Slice nimmt den `CHANGELOG` **nicht** aus.

**A-9 — `internal/adapter/driving/cli/cli_handbuch_link_test.go` (LOW).**
Die drei Tests halten die **Gestalt** der URL, nicht ihr **Ziel**. Kontrolle **M4**: beide Literale konsistent auf `github.com/fremd/anderes-repo/blob/main/irgendwo/benutzerhandbuch.md` → **alle Tests grün**. Da `external` strikt opt-in ist, hält im Repo **nichts** die Aussage „diese URL zeigt auf dieses Handbuch". §5 nennt das Risiko eine Stufe später (Umbenennung des fremden Systems); die Stufe davor — ein von Anfang an falsches Ziel — ist unbenannt.

**A-10 — DoD-Punkt 4, Präzisierung (INFO).**
„jeder wird von **genau der** Mutation rot" trägt für die zwei *Nennt*-Tests. `TestHandbuchURL_TraegtKeineVersion` fällt **auch** bei M1 und M2 (fehlende Zeile ⇒ `zeile == ""` ⇒ `Fatalf`) — fail-closed und richtig so, aber er ist ein Obermengen-Wächter. Die scharfe Probe für ihn ist **M3b** (1 rot), nicht M3 (3 rot); M3 fängt zwei Tests nur, weil das Test-Literal *nicht* mitgezogen wurde.

**A-11 — Beobachtung (INFO).**
Der Slice führt sich als **wellenlos**, während `welle-86` offen ist. Kanonisch zulässig, Planner-Urteil — notiert, weil dieselbe Beobachtung schon bei slice-180 stand.

**A-12 — `361ed39` (INFO).**
`Verantwortlich:` und der Nachtlauf-Block entstanden, während der Plan noch in `open/` lag; der Move folgte 21 s später. `AGENTS.md` §5 sagt „spätestens bei der Beanspruchung" (erfüllt) **und** „ein Plan in `open/` trägt ihn noch nicht" (für einen Commit lang nicht zutreffend). Nicht blockierend.

**A-13 — Commit `ca43027` (INFO).**
*„die erste Fassung war daran gescheitert, und das war der Test, der sich selbst korrigiert hat"* ist eine Prozess-Aussage ohne Artefakt im Range (kein Zwischenstand committet) — aus dem Repo weder bestätigbar noch widerlegbar. Der *Sachgehalt* dahinter ist gemessen: `DCHECK_IMAGE` kommt schon in `mk.out:8` im Digest-Hinweis vor, der Anker `\
DCHECK_IMAGE ?=` ist also nötig, und M5 zeigt, dass er wirkt.

## 6. Lifecycle (E)

- **`3458610` ist ein reiner Move:** `git diff --raw -M` meldet `R100`, Blob-Hash `0c9ad12` vor und nach — die Slice-Datei ist byte-identisch.
- **Der gekoppelte Roadmap-Flip liegt im selben Commit** (`MR-013`): `Nichts in Arbeit.` verlässt §Offene Wellen, der `welle-86`-Zeiger bleibt. Bijektion gemessen: **eine** flache Welle-Datei ⟺ **ein** Zeiger; `in-progress/` trägt genau `slice-181`. `planning-check` **613/0, Exit 0**.
- **Nachtlauf-Block (`MR-053`) stimmt byte-genau** — beide Achsen `gruen`, `2026-08-30T06:08:17Z` und `2026-08-30T09:16:25Z`, reproduziert. Er trägt korrekt **keine** `cite`-Direktive und sagt das ausdrücklich.
- **Die zwei `d-check:cite`-Direktiven sitzen auf den vorschreibenden Zeilen** (`MR-054`): `modul-05:213-214` = *„**Sub-Area-Wahl prüfen.** … Schwelle ≥ 2"*, `modul-05:219` = *„**Offene Beobachtungen sichten.**"* — gegen die Baseline-Datei nachgeschlagen und vom `citations`-Modul im grünen `doc-check` wortgleich bestätigt.
- Commit-Kennungen: `trace-check` über `84f7259~1..50d461e` grün. Die Commit-Zerlegung ist sauber (Plan → Vorbereitung → reiner Move → Feat → Doku); `AGENTS.md` §3.3 eingehalten.

## 7. Was besser ist als zugesagt (F)

1. **Die tragende Begründung ist mutativ belegt, nicht nur zitiert.** P1/P2 zeigen: eine versionierte Handbuch-URL passiert **alle elf** Default-Module still, während derselbe Lauf einen veralteten `ghcr`-Pin an derselben Datei rot macht. Die DoD verlangt nur „geprüft, dass ein Release sie nicht anfassen muss".
2. **Die zweite Behauptung ist ebenfalls mutativ belegt** (P3/P4): die zwei nackten Tags driften still, der Header-Stempel nicht — die Zahl „zwei" ist damit nicht nur aus `releasing.md` abgeschrieben, sondern gemessen und erschöpfend nachgezählt.
3. **Der Versions-Wächter hält der schärfsten Probe stand** (M3b): eine *konsistente* Versionierung beider Literale — die realistische Fehlbedienung — macht **genau ihn** rot. Das ist der Fall, den die Commit-Botschaft verspricht („fällt, sobald jemand `blob/v0.68.0` schreibt") und den M3 allein nicht belegt hätte.
4. **Der Positions-Anker ist eigenständig wirksam** (M5): ein Zeiger unterhalb `DCHECK_IMAGE` macht **genau** den Fragment-Test rot. Die Zusage „im Kopf, nicht irgendwo" trägt.
5. **Der Vorzustand ist vollständig belegt** (V0): die neue Testdatei kompiliert gegen den alten Code und ist dreifach rot — der Test hängt nicht an der Konstante des Prüflings.
6. **Die Illustration ist stärker als „abgekürzt korrekt":** die drei nicht elidierten Kopfzeilen sind **wortgleich** mit der echten Ausgabe, obwohl der `not-replayable`-Marker das gar nicht verlangt.
7. `--print-mk` läuft **ohne jeden Mount** und ist über zwei Läufe byte-identisch — beide `DC-FA-CLI-010`-Nebenzusagen mitgemessen.

## 8. Repo unverändert

`git status --porcelain` leer · `HEAD = ff21044` unverändert · `.harness/state/gates-passed.diffsha` trägt weiter `mtime 2026-08-30 11:29` (der Lauf des Autors) — `record-gates` lief nicht · `d-check:latest` unverändert (`7b8b931dd129`). Alle Repo-Läufe mit `:ro`-Mount; alle Mutationen in einer Scratch-Kopie unter `…/scratchpad/verify181/`, die samt allen dort erzeugten Docker-Tags (`d-check-v181:*`) entfernt ist.

## 9. Empfehlung

**Vor der Closure zu schließen:**
- **A-2 / A-3** — die beiden `.a`-Verfeinerungen in `spec/spezifikation.md` nachziehen (Usage-Enumeration um das sechste Element, Kopf-Enumeration um das dritte) samt Historie-Zeile, und den Slice-Kopf *Berührte Spec-Stellen* um `DC-FA-CLI-001.a`/`DC-FA-CLI-010.a` ergänzen. Ohne das ist die Spezifikation an genau der Stelle unvollständig, die dieser Slice geändert hat.
- **A-8** — `## [Unreleased]`-Eintrag im `CHANGELOG.md`.
- **A-4** — die Aussage „Keine Sprach-Kennzeichnung" in der Lastenheft-Historie richtigstellen; entweder die Formulierung („kein *fremdsprachiger* Marker wie `(German)`") oder die Ausgabe.
- **A-7** — den `DC-FA-VER-001`-Satz im Akzeptanzkriterium `:98` auf seinen Geltungsbereich zurückschneiden: *„… hält die Pins, die `versions.pin-pattern` beschreibt — in diesem Repo `ghcr`-präfixierte."*
- **A-1** — Leerzeile zwischen `const handbuchURL` und `func writeUsage`, Doc-Kommentar an die Funktion zurückgeben (die URL-Begründung als eigener Const-Kommentar).

**Empfohlen, außerhalb der DoD:**
- **A-6** — „elf" → „zwölf" plus `doc-structure` in `docs/user/benutzerhandbuch.md:839-842` und `spec/spezifikation.md:692`. Dritte Instanz derselben Klasse; ein Folge-Slice oder ein Registereintrag wäre hier besser angelegt als eine weitere Punktkorrektur (der `targets`-Gate misst `make`-Targets, nicht `--print-mk`-Targets — die Klasse hat keinen Wächter).
- **A-5** — die Commit-Botschaft ist geschrieben; die Richtigstellung gehört in die Closure-Notiz.
- **A-9 / A-10** — als Grenze in §5 bzw. in der Closure-Notiz benennen: die Trias hält Gestalt und Position, nicht das Ziel; der Versions-Wächter ist fail-closed und damit ein Obermengen-Wächter.

**Nicht blockierend, sachlich in Ordnung:** die Entscheidung `blob/main`, ihre Begründung, die Testarchitektur (Literal statt Import), die Commit-Zerlegung, der Lifecycle-Move und der Nachtlauf-Block.
