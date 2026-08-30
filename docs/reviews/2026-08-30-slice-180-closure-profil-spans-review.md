# Review — slice-180 „`spans` am Closure-Bindepunkt"

**Review-Art:** Code/Design (gegen Slice-Plan, ADR-0076, Hard Rules)
**Gegenstand:** `085459e~1..6835944` (4 Commits), HEAD `6835944`, ungepusht
**Skill:** `.harness/skills/reviewer.md` @ v1.13.0 · **Modell:** claude-opus-5[1m] · **Datum:** 2026-08-30
**Eingangs-Kontext:** `AGENTS.md` §3.1/3.3/3.4/3.5/3.6/3.7/3.8/§4/§5, `harness/conventions.md`, `.harness/baseline/v5.12.0/regelwerk/modul-05/06/08`, ADR-0042/0048/0059/0076, `DC-FA-SPAN-001`, `DC-FA-TGT-001`, Beobachtungs-Register (`BEO-009`, `BEO-012`, `BEO-013`, `BEO-020`, `BEO-023`), Vor-Befund `docs/reviews/2026-08-22-slice-115-arc-vergabe-review.md` F-4.

## Eigener Lauf

| Lauf | Ausgabe |
|---|---|
| `make doc-check` | `608 Datei(en) geprüft, 0 Befund(e)` · Exit 0 |
| `make verify-closure-notes` | `546 Datei(en) geprüft, 0 Befund(e)` · Exit 0 |
| `make gate-consistency` | `608 Datei(en) geprüft, 0 Befund(e)` · Exit 0 |
| `make adr-check RANGE=085459e~1..6835944` | `608 Datei(en) geprüft, 0 Befund(e)` · Exit 0 |
| `make test` | alle Pakete `ok` · Exit 0 |
| eigene Differential-Messung (awk, beide Fence-Automaten) | `678 Dateien, 120284 Zeilen, 0 divergente Zeilen` |
| Proben-Läufe gegen `d-check:latest` in isolierten Klonen | s. u. |

`make gates`/`make record-gates` bewusst **nicht** gefahren (Nachweis-Datei).

## Urteil

**BLOCKIERT** — 3 HIGH · 2 MEDIUM · 2 LOW · 2 INFO.

Der Bau selbst ist korrekt und die diskriminierende Probe hält. Blockierend sind **drei falsche Aussagen über Messungen und Sensoren**, alle drei nachgemessen, zwei davon in einer bereits `Accepted`-ADR.

---

## Findings

### H-1 · Die „Richtigstellung" an slice-178 §1 misst den falschen Abschnitt

- **quelle:** `BEO-020`, `AGENTS.md` §5 (Botschaft behauptet nicht mehr, als die Arbeit trägt)
- **pfad:** `docs/plan/planning/open/slice-178-offene-tasks-roh.md:52-61`
- **zugesagt:** *„Die Expositions-Behauptung dieses Absatzes war falsch gezählt und ist korrigiert … Gezählt wurde **abschnittsweise**, gepaart wird aber **absatzweise** … Nachgemessen: 21 bzw. 4 Backticks."*
- **gemessen:**

  | Datei | §3 *Definition of Done* | §4 *Risiken / offene Punkte* |
  |---|---|---|
  | `slice-061` (§3 = Z. 80–99, §4 = Z. 100–113) | **25** Backticks | **21** |
  | `slice-076` (§3 = Z. 64–106, §4 = Z. 107–128) | **45** Backticks | **4** |

  Die neuen Zahlen **21 bzw. 4** sind exakt die Zahlen von **§4**, nicht von §3 (*Definition of Done*), über den die korrigierte Aussage spricht. Die alten Zahlen **25 bzw. 45** reproduzieren dagegen exakt.
  Zweitens trägt der behauptete Zähl-Fehler hier gar nicht: eine Absatz-Zählung nach Produkt-Semantik (`proseParagraphs`) ergibt in **beiden** Dateien einen **einzigen** Absatz über den ganzen DoD-Abschnitt — 25 Backticks ab `slice-061:82`, 45 ab `slice-076:66`. Abschnitt **ist** hier Absatz; „eine Abschnitts-Summe sagt über den Mechanismus nichts" ist in diesem Fall gegenstandslos.
- **warum es zählt:** Eine korrekte, reproduzierbare Zahl wird durch eine nicht reproduzierbare ersetzt, und der Vorgang wird als Auflösung genau der Klasse verbucht, die er begeht (`BEO-020` — *„gemessen wird die eine Menge, ausgesagt wird über die andere"*). Die falsche Zahl steht in einem `open/`-Plan, in drei Commit-Botschaften und wird im Slice-Plan §7 als BEO-020-Instanz geführt.
- **verifizierbar:** ja — `awk`-Zählung je Abschnitt über beide Dateien; kein Gate fängt es.
- **klasse:** `korrektur-misst-falschen-abschnitt`
- **billigster Fix:** Zahlen auf den DoD-Abschnitt zurückstellen (die Commits sind ungepusht).

### H-2 · ADR-0076 §Fitness Function schreibt `gate-consistency` eine Deckung zu, die es nicht hat

- **quelle:** `AGENTS.md` §4 (*„Halluzinierte Gates sind die häufigste Form von Harness-Lüge"*), `DC-FA-TGT-001`
- **pfad:** `docs/plan/adr/0076-spans-am-closure-bindepunkt.md:138`; gleichlautend `docs/plan/planning/done/slice-180-closure-profil-spans.md:79` und die DoD-Zeile `:118-120`
- **zugesagt:** *„`make gate-consistency` hält dafür den `##`-Hilfetext gegen die Doku."* („dafür" = als Ersatz dafür, dass die `spans`-Verortung im Rezept maschinell ungeprüft ist.)
- **gemessen:** In einem vollständigen Klon den `##`-Hilfetext des Targets durch `## irgendwas ganz anderes.` ersetzt → `--enable targets` + Fokus-`--disable`: `608 Datei(en) geprüft, 0 Befund(e)`, Exit 0. Danach zusätzlich `--enable spans` **aus dem Rezept entfernt** → weiterhin `608 / 0 Befunde`, Exit 0; auch der volle `doc-check` bleibt grün. Der Quelltext bestätigt es: `internal/hexagon/core/rules/targets.go:24-28` liest Makefile-**Regelzeilen** (`makefileRuleRe`) und Doku-**Tabellenzeilen** (`docTargetRe`); der `##`-Text kommt in keinem der beiden Pfade vor.
  Das Repo hat diesen Befund bereits: `docs/reviews/2026-08-22-slice-115-arc-vergabe-review.md:43` — *„kein Gate liest die Beschreibungsspalte (`targets` vergleicht Target-**Namen**, nicht Texte)"*.
- **warum es zählt:** Die `Fitness Function`-Sektion ist genau das Register, das sagt, was eine Maschine hält. Ein falscher Eintrag dort ist für jeden Leser dasselbe wie ein halluziniertes Gate — und er steht in einer `Accepted`-ADR, die §3.5 einfriert. Faktisch hält **kein** Sensor die vier Deklarations-Flächen gegen das Rezept (gemessen).
- **verifizierbar:** ja — die zwei Mutationsläufe oben.
- **klasse:** `gate-deckung-ueberdehnt`
- **billigster Fix:** Satz vor dem Push streichen oder auf das reduzieren, was `targets` hält (Target-**Name** in Makefile ↔ Doku-Tabelle).

### H-3 · „Der Defekt fällt erst bei `make fullbuild` auf" ist falsch — der Bindepunkt schließt keine Lücke, die nicht schon geschlossen wäre

- **quelle:** `AGENTS.md` §5 (Botschaft trägt nicht weiter als die Messung), Reviewer-Skill Frage 10
- **pfad:** `docs/plan/adr/0076-spans-am-closure-bindepunkt.md:112-116` (§Verglichene Alternativen); gleichlautend `docs/plan/planning/done/slice-180-closure-profil-spans.md:144-148` (§5, Risiko 4)
- **zugesagt:** *„Den Bindepunkt in `gates` ziehen … Verworfen … Der Preis bleibt: der Defekt fällt erst bei `make fullbuild` auf."* Dazu §Kontext `:54` *„die einzige, die heute niemand sieht"*.
- **gemessen:** Vollständiger Klon, an `docs/plan/planning/done/slice-177-structure-hint.md` ein offener Fence angehängt:

  | Lauf | Ergebnis |
  |---|---|
  | Hauptprofil (= `make doc-check`, in `gates`, im `pre-commit`-Hook) | `slice-177…:313  ```text  fence-unclosed` · Exit 1 |
  | Closure-Bindepunkt (`planning`+`structure`+`spans`) | **identischer Befund** · Exit 1 |

  `spans` steht seit langem in `modules:` der `.d-check.yml` (`:29`), `gates` enthält `doc-check` (`Makefile:305`), und `.githooks/pre-commit:20` ruft `make doc-check`. Die Scan-Menge des Bindepunkts ist zudem eine **Teilmenge** der von `gates` (546 ⊂ 608; das Closure-Profil hat keinen `scan`-Block, also `defaultRoots() = {docs, spec}`, `rules/scan.go:25`). Es gibt damit **keine** Datei, für die der Bindepunkt einen `spans`-Befund sieht, den `gates` nicht schon vorher rot meldet.
- **warum es zählt:** Der Satz benennt den Preis der gewählten Alternative — und er ist der einzige Ort, an dem die ADR über den Zeitpunkt der Entdeckung spricht. Gemessen fällt der Defekt am **Commit** auf, nicht bei `make fullbuild`. In der Folge liest §Konsequenzen *„Positiv"* wie neu gewonnene Deckung, während der Zuwachs an gefundenen Defekten über den Repo-Lauf **null** ist. Der Slice-Plan selbst hatte es in §1 (`:33`) noch ehrlich: *„zwei Defekte …, die das Produkt längst findet"* — diese Einschränkung ist auf dem Weg in die ADR verlorengegangen.
- **verifizierbar:** ja — die zwei Läufe oben gegen denselben Klon.
- **klasse:** `entdeckungszeitpunkt-falsch-behauptet`
- **billigster Fix:** vor dem Push den Satz auf das Gemessene stellen und die tragende Begründung dorthin verschieben, wo sie hält (Selbstgenügsamkeit des Bindepunkts, nicht neue Deckung).

### M-1 · „Die Exposition dieses Repos ist damit heute null" folgt nicht aus „`spans` meldet nichts"

- **quelle:** `BEO-009`(b) / `BEO-020`, `AGENTS.md` §5
- **pfad:** `docs/plan/planning/open/slice-178-offene-tasks-roh.md:59-62`
- **zugesagt:** *„… und `spans` meldet in beiden Dateien **nichts**; kein Absatz ist unbalanciert. **Die Exposition dieses Repos ist damit heute null**."*
- **gemessen:** Das Verschlucken entsteht durch **Paarung**, nicht durch Unbalance. Eine ungeschlossene Folge ist literal; `span-unclosed` feuert überdies nur, wenn sie an Nicht-Whitespace klebt (`rules/spans.go:114-118`, `sticksToText`). Gegenprobe im Klon: an `slice-061` DoD-Absatz ` Ein \\`Rest (bei Closure) hier\\` steht da.` angehängt → `structure`: **0** `section-forbidden`; `spans`: **0 Befunde**. Genau die Klasse, die ADR-0076 selbst als Tabellenzeile 3 („wohlgeformter Span, der Prosa umschließt") als **nicht** deckbar ausweist. „`spans` still" kann Exposition also grundsätzlich nicht ausschließen.
  (Positiv-Kontrolle: an denselben Absatz ohne Span angehängt → `section-forbidden`, Exit 1; alle 17 Zeilenenden des Absatzes einzeln geprüft, überall sichtbar — die *aktuelle* Blindheit ist tatsächlich null, aber nicht aus dem angegebenen Grund.)
- **warum es zählt:** Der Schluss stützt die Neubewertung der Dringlichkeit von slice-178 („der Slice steht auf dem konstruierten Fall"). Er ist die zweite Richtung von `BEO-009` und trifft dieselbe Datei wie H-1.
- **verifizierbar:** ja — die Gegenprobe oben.
- **klasse:** `stiller-sensor-als-abwesenheitsbeleg`

### M-2 · Der Kopfkommentar des Profils legt eine zweite, ungewächterte Modul-Liste an — und widerspricht drei Zeilen später sich selbst

- **quelle:** ADR-0076 Entscheidung 2, `BEO-010` (Spiegel-Disziplin), `AGENTS.md` §3.7 (Abgrenzung)
- **pfad:** `.d-check.closure.yml:13-22`
- **zugesagt:** `:21` *„… welche Module laufen, steht an EINER Stelle — im Makefile-Rezept. Zwei Orte drifteten."*
- **gemessen:** `:14` derselben Kommentarblocks nennt die Module: *„heute `planning`, `structure` und `spans`"*. Das ist der zweite Ort. Dass er driftet, ist **belegt, nicht vermutet**: die abgelöste Fassung dieses Satzes stammt aus `e101e09` (slice-093) und sprach von *„per `--enable planning` dazuschaltet"*; `--enable structure` kam am 2026-08-15 mit `e93d6a9` (slice-099) ins Rezept — der Kommentar war zwei Wochen lang falsch, und kein Lauf hat es gemeldet. Gemessen bleibt es so: `--enable spans` aus dem Rezept entfernt ⇒ `doc-check` und `gate-consistency` grün (s. H-2), während vier Flächen weiter `spans` behaupten.
- **warum es zählt:** Die Entscheidung, `spans` **nicht** in `modules:` zu setzen, wird mit „ein Ort" begründet; die Umsetzung schafft im selben Commit den zweiten. Ein `modules:`-Eintrag wäre wenigstens wirksam gewesen, dieser Satz kann nur falsch werden.
- **verifizierbar:** nein maschinell — genau das ist der Befund.
- **klasse:** `prosa-spiegel-ohne-waechter`

### L-1 · `span-unclosed` in AGENTS.md als Artefakt geführt, „das die Bereinigung still macht"

- **pfad:** `AGENTS.md:386`
- **zugesagt:** *„`spans` meldet die Span- und Fence-Artefakte, die deren Bereinigung **still** machen … (`fence-unclosed`, `span-unclosed`)"*
- **gemessen:** Nach der eigenen Tabelle in `docs/plan/adr/0076-…:45` ist der Fall „ungerader Backtick im Absatz" bei `structure` **nicht** still (`section-forbidden`). `span-unclosed` deckt also keinen Stillstellungs-Fall; das tut nur `fence-unclosed`. Der Leser nimmt daraus die Inferenz mit, die M-1 begeht.
- **klasse:** `code-liste-ueberdehnt`

### L-2 · „alle 676 Markdown-Dateien dieses Repos" — die Grundgesamtheit ist nicht benannt und ist nicht die des Repos

- **pfad:** `docs/plan/adr/0076-…:62`, `docs/plan/adr/0042-…:149`, `docs/plan/planning/done/slice-180-…:52`
- **gemessen:** heute `git ls-files '*.md'` = **661**; `find . -name '*.md' -type f` = **674** (17 davon gitignoriert unter `.harness/cache/`, `.gitignore:4`); `find -L … -type f` = **678** (löst 4 Symlinks unter `.claude/rules/` auf, die bereits gezählte Regelwerk-Module dupliziert). Rekonstruktion für `9d22a44` (Elter des Plan-Commits): 655 echte + 4 Symlink-Auflösungen + 17 gitignorierte = **676** — die Zahl ist also `find -L`-Semantik über eine Menge, die 17 Dateien **außerhalb** des Repos enthält und 4 doppelt zählt.
- **warum nur LOW:** die Menge ist eine **Obermenge**; „0 Abweichungen" gilt damit erst recht für das Repo. Aber die Zahl ist heute nicht reproduzierbar, und `BEO-020` verlangt ausdrücklich *„die gezählte Menge benennen, bevor die Zahl fällt"* — genau die Prozedur, die dieser Slice zitiert. Zwei der drei Fundstellen sind bereits immutabel bzw. `Accepted`.
- **klasse:** `zahl-ohne-nenner`

### I-1 · Kommentar-Überschrift trägt eine Ortsangabe, die nicht stimmt

`.d-check.closure.yml:19` — *„WARUM `spans` HIER STEHT, ABER NICHT IN DIESER LISTE"*: `spans` steht in dieser Datei überhaupt nicht; „hier" meint den Bindepunkt, das einzige „hier" der Datei ist aber die Liste.

### I-2 · Der Go-Test, auf den Entscheidung 2 sich stützt, trägt einen veralteten Kommentar

`internal/adapter/driven/configyaml/gate_consistency_test.go` (`TestQA03_ClosureProfil_KeineZweiteNetzTuer`) — *„es ist ein fokussiertes Profil, das nur `planning` per Kommandozeile dazuschaltet"*. Seit slice-099 falsch, jetzt doppelt. **Vom Slice nicht angefasst** (Bestand), aber ADR-0076 `:80` beruft sich auf diesen Test; substantiell prüft er nur `modules:` ⊄ {external, sources, vcs} (bei `modules: []` trivial erfüllt) und `planning.closure.dir ≠ ""`. Die tatsächlich laufende Modul-Menge kommt von der Kommandozeile und wird von ihm nicht gelesen.

---

## Negativbefunde (geprüft, ohne Befund — mit der Messung)

1. **Die tragende Fence-Messung reproduziert.** Beide Automaten in `awk` nachgebaut (`FenceToggle` / `FenceRun`+`FenceCloses`, `TrimFenceIndent`, `FenceCloses`-Whitespace-Regel) und zeilenweise gegeneinander gehalten: **678 Dateien, 120 284 Zeilen, 0 divergente Zeilen, keine Lesart am Dateiende offen.**
2. **Die Messung kann finden.** Vier konstruierte Fälle je 1 Divergenz gemeldet — Infostring hinter dem Schluss-Fence, ` ``` `→`~~~`, `` ```` ``→` ``` `, `~~~`→` ``` ` —, die saubere Datei still.
3. **`spans` öffnet den strengen Automaten korrekt.** `rules/spans.go:80-93` (`stepStrict`) öffnet über `FenceToggle`, also inklusive CommonMark-Infozeilen-Regel — identisch zu `app/trace_table.go:329-341`. Die Modellierung ist damit die richtige; der Repo-Bestand hat genau **eine** Zeile, an der die Regel trägt (`docs/plan/adr/0042-…:131`), und beide Lesarten behandeln sie gleich. Die veröffentlichte 0 entspricht der korrekten Modellierung, nicht der naiven.
4. **Die diskriminierende Probe hält.** Echter `done/`-Slice (`slice-177`), unverändert `537 Dateien, 0 Befunde, Exit 0`; mit angehängtem offenem Fence `fence-unclosed`, Exit 1; `planning`+`structure` allein bleiben bei **0 Befunden, Exit 0**. Die erste Probe war also zu Recht verworfen.
5. **Die Motivation ist real und gemessen.** Platzhalter `(bei Closure)` im H1-Abschnitt eines `done/`-Slice ⇒ `section-forbidden`, Exit 1. Derselbe Platzhalter **hinter** einem offenen Fence ⇒ `0 Befunde, Exit 0` (stilles Grün). Mit `spans` ⇒ `fence-unclosed`, Exit 1. Der ADR-0059-Preis hat tatsächlich die zweite Hälfte, die die ADR beschreibt.
6. **§3.6 — keine Lockerung.** AGENTS.md-Zeile und Sensors-Zeile behalten „**Vier Bereiche sieht sie nicht**" bzw. „**Vier benannte Grenzen**" wortgleich; der Zusatz nennt ausdrücklich, dass er die **vierte** Grenze nicht deckt. Der Profil-Kopfkommentar wird präziser, nicht schwächer. Keine gesenkte Schwelle im Diff.
7. **Reichweiten-Richtung ist sauber (§3.8).** Closure-Profil ohne `scan`-Block ⇒ `defaultRoots() = {docs, spec}` ⇒ 546 Dateien ⊂ 608 des Hauptprofils. „0 über den ganzen Baum" ⇒ „0 am Bindepunkt" ist damit die zulässige Richtung (nachgezählt). Dass `spans` die **ganze Datei** und nicht nur den Closure-Abschnitt prüft, steht in §Konsequenzen.
8. **§3.5 hält.** `make adr-check RANGE=085459e~1..6835944` grün; die Änderung an `ADR-0042` ist ein reiner `## Geschichte`-Anhang (eine Zeile, Kern unberührt).
9. **§3.4 / Referenz-Richtung.** ADR-0076 trägt **kein** `slice-`Token (`grep`); `matrix` im `doc-check` grün.
10. **§3.3 / MR-013.** `6b485a6` ist ein reiner `git mv` (Slice-Datei 0 Änderungen) plus der gekoppelte Roadmap-Flip (Ruhe-Marker entfernt, Zeiger-Liste bleibt) — genau die deklarierte Ausnahme.
11. **Zitat-Treue (`BEO-012`).** Die Zitate aus ADR-0042 („bewusst offen gelassen … unbelegt — kein Realfall in den 522 Dateien", „naiver Toggle", „wenn einer existiert") stehen wortgleich in der Quelle und im Geltungsbereich des offenen Punkts (a); der Eintrag sagt korrekt, dass der Punkt **offen bleibt**. Eine Paraphrase setzt „Sensor" für das Quell-Wort „Regel" — nicht meldenswert, aber notiert.
12. **§3.7 / Kommentar-Klassen.** Der neue Makefile-Kommentar (`Makefile:334-338`) trägt Zusage/Kopplung/Grenze, **eine** auflösbare Herkunft (`ADR-0076`), keine Slice-Nummer, keine Review-Historie, kein Mess-Label. Der `##`-Hilfetext führt weiterhin drei ADR-Kennungen — Bestandsform, keine Neuerung.
13. **Adressierungs-Form (Skill-Frage 16).** ADR-0076 §Bezug nennt durchgehend Kennungen (`ADR-*`, `DC-FA-*`), kein blankes „§N"; `Schärft: —` ist begründet. ADR-Index-Zeile vorhanden und formgleich.
14. **Zustandsfelder (§3.7).** Slice-Kopf ohne `Status:`-Feld, mit `Lifecycle:`; die geänderten Doku-Zeilen nennen Zustand und Beleg (`gemessen: … vorher grün, nachher rot`, `Bestands-Rauschen: null`), keine Chronik. Kein Drift-Log-Eintrag missbraucht.
15. **Nicht-Ziele (§3 des Plans).** Die Abgrenzung „umschließender Span bekommt keine Antwort" ist ehrlich und in drei Flächen wiederholt; die Klasse ist nach Produkt-Semantik tatsächlich nicht unterscheidbar (gemessen: `structure` und `spans` beide still) — sie verdeckt keine billig schließbare Lücke.
16. **`make test` grün** (alle Pakete `ok`), einschließlich der `configyaml`-Suite, die das Closure-Profil liest.

## Repo-unverändert-Beleg

```
$ git status --porcelain
(leer)
$ git rev-parse HEAD
683594449ab72ea8537546798e4d64a1501b8b45
```

Alle Proben liefen in `/tmp/…/scratchpad/rev180/` (eigenes Unterverzeichnis, danach entfernt); das Repo war in jedem Container-Lauf `:ro` gemountet, `--network none`.
