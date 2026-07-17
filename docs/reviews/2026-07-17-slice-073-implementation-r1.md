# Review — slice-073 Implementierung R1 (nachgeholt, nach Release v0.44.1)

**Datum:** 2026-07-17
**Review-Art:** unabhängiger Implementierungs-Review (kontext-getrennt;
adversarial, nicht bestätigend). **Nachgeholt** — der Gegenstand ist bereits
als v0.44.1 veröffentlicht.
**Gegenstand:**
[`slice-073`](../plan/planning/in-progress/slice-073-link-transparente-range-fortsetzung.md) —
Range `ca0d631~1..91f1a52` (`ca0d631` doc-first, `2954e4d` Fix + Tests,
`a5bcdec` Release-Prep, `91f1a52` Digest-Backfill)
**Reviewer:** Claude (kontext-unabhängiger Lauf)
**Skill:** `.harness/skills/reviewer.md` v1.2.0
**Modell:** Opus 4.8 (1M context)

**Scope-Hinweis:** Während des Reviews hat eine Parallel-Session HEAD auf
`c4c08f3` (slice-071) bewegt. Der Review-Scope ist auf `ca0d631~1..91f1a52`
gepinnt; `internal/hexagon/core/app/trace.go` ist seit `91f1a52` unverändert
(`git diff 91f1a52..HEAD -- …/trace.go` leer), die Befunde gelten also
unverändert für HEAD.

## Eingangs-Kontext

- [`DC-FA-COV-001`](../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
  (Lastenheft, **unverändert** durch den Slice),
  [`DC-FA-COV-001.a`](../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)
  Schritt 3 (geschärft),
  [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
  (mit-betroffen über den geteilten Parser),
  [`DC-FA-ID-001`](../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
- [ADR-0039](../plan/adr/0039-link-transparente-range-fortsetzung.md) (Proposed),
  [`AGENTS.md`](../../AGENTS.md) §3, [`docs/user/releasing.md`](../user/releasing.md)
  §Release-Prep
- Prüfgegenstand: `internal/hexagon/core/app/trace.go`,
  `internal/hexagon/core/app/trace_coverage_test.go`, `spec/spezifikation.md`,
  ADR-0039, Release-Prep-Artefakte

**Verifikation gegen das echte Image.** `make build` (post-fix, `d-check:latest`)
und ein Vergleichs-Image aus `c400f18` (`make build IMAGE=d-check-pre`, entspricht
v0.44.0). Läufe netzlos/read-only:
`docker run --rm --network none -v <fixture>:/repo:ro -w /repo <img> --trace`.

---

## Findings

### F-1 — HIGH · Neue Falsch-Deckung: URL-Interna werden als Range/Enum expandiert und verstecken Waisen

**kategorie:** HIGH
**quelle:** [`DC-FA-COV-001`](../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in),
[`DC-FA-COV-001.a`](../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)
Schritt 3, [ADR-0039](../plan/adr/0039-link-transparente-range-fortsetzung.md)
Entscheidung 2 („Kein weiteres Peeling"), Hard Rule „Stilles-Grün-Pfad in einem Gate"
**pfad:** `internal/hexagon/core/app/trace.go:374` (`linkSuffix`),
`internal/hexagon/core/app/trace.go:384–389` (`skipLinkSuffix`),
`internal/hexagon/core/app/trace.go:516`

**befund:** `linkSuffix` (`^`?\]\([^)]*\)`) beendet das Linkziel an der **ersten**
`)`. Enthält das Ziel eine Klammer, konsumiert das Suffix nur den Ziel-**Anfang**;
der **Rest der URL** bleibt in `rest` stehen und wird anschließend von
`rangeSuffix`/`enumSuffix` gelesen. Ziel-Fragmente wie `/002/003.md` oder
`..003.md` werden dadurch als Fortsetzung interpretiert und expandiert, obwohl in
der Zelle **keinerlei** Range-/Enum-Notation steht. Die Expansion ist damit nicht
„eng gefasst", sondern liest an einer Stelle, die vorher unerreichbar war.

Reproduktion (Fixture: `spec/req.md` mit `GG-QA-001..003` als Headings;
`docs/plan/traceability.md` mit **einer** Zelle; `trace.coverage` `ranges: true`):

| Zelleninhalt | v0.44.0 | v0.44.1 | korrekt |
|---|---|---|---|
| `[`GG-QA-001`](../specs/Rev(2)/002/003.md)` | 2 Waisen | **0 Waisen** | 2 Waisen |
| `[`GG-QA-001`](../a(1)..003.md)` | 2 Waisen | **0 Waisen** | 2 Waisen |

Gate-Wirkung, gleiche Fixture:

```text
d-check-pre (v0.44.0)  --trace --require-complete -> exit 1   (korrekt rot)
d-check     (v0.44.1)  --trace --require-complete -> exit 0   (still grün)
```

`GG-QA-002`/`GG-QA-003` erscheinen in der RTM mit Coverage-Label `Trace`, obwohl
sie von keiner Kante gedeckt sind. Das ist die vom Slice selbst als gefährlicher
eingestufte Richtung (§4 „Kein Konsument verliert Deckung, kein Befund entsteht
neu") — hier entsteht **Deckung**, die es nicht gibt, und ein korrekt roter
`--require-complete`-Lauf wird grün. Der Kommentar an `trace.go:369–371`
(„Das ist die sichere Richtung: lieber keine Expansion als eine geratene")
beschreibt nur den `..`-Zweig; im `/`-Zweig ist die Aussage widerlegt.

**verifizierbar:** ja — die obige Fixture gegen `d-check:latest`
(`--trace` bzw. `--trace --require-complete`); kein bestehender Gate-Lauf fängt es,
weil kein Test die Zeichen-Klasse `)` im Ziel mit numerischem Ziel-Rest kombiniert.

---

### F-2 — MEDIUM · Vertrag und Code sind über die Ziel-Grammatik nicht deckungsgleich; der Range-Parser kennt eine andere Link-Definition als `ids`/`links`

**kategorie:** MEDIUM
**quelle:** [`DC-FA-COV-001.a`](../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)
Schritt 3 (§Link-Transparenz), [ADR-0039](../plan/adr/0039-link-transparente-range-fortsetzung.md)
Entscheidung 1, Konsistenz-Lücke zwischen Modulen derselben Eingabe-Klasse
**pfad:** `internal/hexagon/core/app/trace.go:374` vs.
`internal/hexagon/core/rules/markdown.go:368–393` (`parseLinkAt`)

**befund:** Spezifikation und ADR-0039 beschreiben das übersprungene Suffix
unqualifiziert als „`]`, geklammertes Ziel" bzw. „`](…)`". Der Code implementiert
eine **engere** Sprache: `[^)]*` — ein Ziel **ohne** Klammern. Der `ids`-/`links`-Pfad
derselben Eingabe-Klasse liest Linkziele dagegen **klammer-balanciert**
(`markdown.go:382`: `matchBracket(s, textEnd+1, '(', ')')`). Beide Module haben damit
unterschiedliche Vorstellungen davon, was ein Link ist. Folge für einen Konsumenten
mit `link-policy: always`: `` [`GG-QA-001`](https://x.org/A_(b))..003 `` gilt `ids`
als erfüllte Linkpflicht, während der Range-Parser **nicht** expandiert — gemessen
2 Waisen gegen `d-check:latest`. Die strukturelle Kollision zwischen Linkpflicht und
Range-Notation, die ADR-0039 §Konsequenzen als „aufgelöst" führt, besteht für diese
Ziel-Klasse fort; die Verengung steht in keinem der beiden Vertrags-Straten.

**verifizierbar:** ja — Fixture-Zelle `` [`GG-QA-001`](https://x.org/A_(b))..003 ``
gegen `d-check:latest --trace` (2 Waisen statt 0).

---

### F-3 — MEDIUM · Drei nutzersichtbare Kompatibilitäts-Zusagen sind faktisch falsch

**kategorie:** MEDIUM
**quelle:** [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus), Maintainability
**pfad:** `CHANGELOG.md` (§0.44.1 „Wirkung auf bestehende Läufe"),
`docs/user/benutzerhandbuch.md:1257 ff.` (§5 „Range-Notation unter Linkpflicht"),
`internal/hexagon/core/app/trace.go:369–371`

**befund:** Drei Aussagen desselben Release sind durch F-1 widerlegt:

1. CHANGELOG: „Wer keine verlinkten Ranges nutzt, ist nicht betroffen
   (byte-identisch)." — die F-1-Fixture nutzt **keine** Range-Notation; die RTM
   ändert sich (2 Waisen ⇒ 0) und der Exit-Code kippt 1 ⇒ 0.
2. CHANGELOG: „Kein Konsument verliert Deckung, kein Befund entsteht neu." — trifft
   die Richtung nicht, die tatsächlich eintritt: es entsteht **Deckung**, und ein
   Befund **verschwindet**.
3. Handbuch §5: „Enthält das Linkziel selbst eine Klammer, greift die Regel
   ebenfalls nicht." — die Regel greift, sie greift nur falsch (F-1). Derselbe
   Irrtum steht als Kommentar an `trace.go:369–371`.

Diese Zusagen sind zugleich die Begründungsbasis des SemVer-Patch- und
„kein CR"-Arguments (Slice §2, ADR-0039 Entscheidung 5).

**verifizierbar:** ja — dieselbe F-1-Fixture.

---

### F-4 — MEDIUM · Die Fitness-Funktion von ADR-0039 ist von keinem Test mechanisiert; die DoD-Zusage „Tests für `trace.cross-consistency`" hat keine Entsprechung

**kategorie:** MEDIUM
**quelle:** [ADR-0039](../plan/adr/0039-link-transparente-range-fortsetzung.md)
§Fitness-Funktion, fehlende Negativtests bei neuem öffentlichen Vertrag
**pfad:** `internal/hexagon/core/app/trace_coverage_test.go:97–145`;
`internal/hexagon/core/app/trace_cross_test.go` (im Slice-Scope **nicht** angefasst)

**befund:** Der Slice ändert genau **eine** Testdatei
(`git diff --name-only ca0d631~1..91f1a52 | grep test` ⇒ nur
`trace_coverage_test.go`). Alle neuen Tests rufen `expandRange` direkt auf. Damit ist
keine der beiden Fitness-Aussagen von ADR-0039 verriegelt:

- „Der ausgelieferte `trace.coverage`-Falschbefund (2 Waisen statt 0) ist weg" — kein
  Test läuft über `coverageRefs`/den Waisen-Pfad; die Zahl „2 Waisen" existiert in
  keinem Test.
- „Ein Fix, zwei Konsumenten" (Entscheidung 4) — `trace_cross_test.go` ist
  unverändert; kein Test belegt die Wirkung auf `trace.cross-consistency`.

Die Fail-closed-Tests (`TestExpandRangeLinkTransparentFailClosed`) prüfen nur, dass
`err != nil`, nicht welche Fehlerklasse — ein Breiten-Mismatch, der zu einem
AAA>BBB-Fehler mutiert, bleibt grün. `equalStrings` (`trace_coverage_test.go:75–85`)
vergleicht längenbasiert, `nil` und `[]string{}` sind ununterscheidbar; die
Negativfälle pinnen „len 0", nicht „keine Expansion". Der Fall
`"Klammer in der URL bricht das Suffix"` (`:120`) pinnt ausgerechnet den **gutartigen**
Zweig derselben Zeichen-Klasse, aus der F-1 stammt, und erzeugt so den Eindruck, die
Klasse sei abgedeckt.

**verifizierbar:** ja — `make test` bleibt grün, wenn `enumSuffix` in `expandRange`
nach `skipLinkSuffix` ausgewertet wird (F-1-Pfad); kein Testfall kippt.

---

### F-5 — LOW · ADR-0039-Indexzeile ist unvollständig

**kategorie:** LOW
**quelle:** [`docs/plan/adr/README.md`](../plan/adr/README.md) §Konventionen,
[`AGENTS.md`](../../AGENTS.md) §5 („Neue ADRs müssen den ADR-Index aktualisieren"),
Source Precedence Rang 4
**pfad:** `docs/plan/adr/README.md:54`

**befund:** Die Indexzeile für ADR-0039 trägt 4 Pipe-Zeichen, jede andere Zeile
(z. B. ADR-0038) trägt 6 — die Spalten **Datum** und **Bezug** fehlen. Die im Index
geführte Vertrags-Verknüpfung (`DC-FA-COV-001`, `DC-FA-XREF-001`, `DC-FA-ID-001`) ist
für ADR-0039 als einzige nicht nachschlagbar; kein Gate prüft die Spaltenzahl.

**verifizierbar:** ja — Sichtprüfung/`awk`-Pipe-Zählung über `docs/plan/adr/README.md`;
`make gates` fängt es nicht.

---

## Negativbefunde (geprüft, ohne Befund)

- **Fix-Platzierung / Fail-closed hinter Link:** `skipLinkSuffix` sitzt nach der
  Familien-/Breiten-Ableitung (`trace.go:514–516`), vor beiden Suffix-Zweigen.
  `AAA>BBB` und Breiten-Mismatch bleiben hinter einem Link Fehler (Exit 2) —
  Code-Pfad und Test decken sich; ohne Befund.
- **`ranges: false`:** `rangeAwareIDs` kehrt vor `expandRange` zurück
  (`trace.go:491–493`); `skipLinkSuffix` ist unerreichbar. Byte-Identität für
  Konsumenten ohne `ranges` gilt; ohne Befund.
- **„Genau ein Suffix" (Mehrfach-Peeling):** `skipLinkSuffix` ist nicht geschleift;
  `` `](a.md)](b.md)..003 `` expandiert nicht. Code und ADR-0039 Entscheidung 2
  decken sich; ohne Befund.
- **Whitespace / Zeichen zwischen `)` und `..`, mehrzeilige Zellen:** `rangeSuffix`
  und `enumSuffix` sind `^`-verankert; ein Newline oder beliebiges Zeichen bricht.
  Kennung am Zellenende (`[`ID`](a.md) |`) expandiert nicht; ohne Befund.
- **`..` ohne Ziffern:** `rangeSuffix` = `^\.\.(\d+)`; ein relativer Pfad `../foo`
  hinter einem Link expandiert nicht; ohne Befund.
- **Referenz-Links `[x][y]` — REFUTED als Finding:** `parseLinkAt`
  (`markdown.go:379`) verwirft alles, dessen `]` nicht `(` folgt
  (`if !ok || textEnd+1 >= len(s) || s[textEnd+1] != '('`). Referenz-Links erfüllen die
  `ids`-Linkpflicht also gar nicht erst — die Nicht-Expansion im Range-Parser ist mit
  `ids` konsistent, kein Widerspruch. (Kontrast zu F-2, wo `ids` den Link **akzeptiert**.)
- **Bildlinks `![x](y)`:** ein Anforderungs-Kennungs-Treffer im Alt-Text mit
  anschließender Range hat kein erzählbares Versagen; nicht gemeldet
  (Anti-Pattern „kein Finding ohne Failure-Szenario").
- **Hard Rules ([`AGENTS.md`](../../AGENTS.md) §3):** keine `//nolint`-Direktive
  (§3.2); keine Host-Toolchain (§3.1 — Verifikation lief ausschließlich über
  `make build` + Image); `2954e4d` fügt `trace.go` **keine** Imports hinzu, der
  ADR-0005-Schnitt ist unberührt; ADR-0039 ist neu und `Proposed`, §3.5 nicht
  berührt; keine Gate-Schwelle gesenkt (§3.6); keine `git mv`-Kollision (§3.3).
  Ohne Befund.
- **Netzlosigkeit ([`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)):**
  der Fix ist rein lexikalisch, kein neuer I/O- oder Netz-Pfad; alle
  Verifikationsläufe liefen mit `--network none` erfolgreich. Ohne Befund.
- **Release-Prep-Checkliste** ([`releasing.md`](../user/releasing.md) §Release-Prep),
  Punkt für Punkt — **vollständig, ohne Befund:**
  1. `version.md` — §Aktuell auf `v0.44.1`, neue §Verlauf-Zeile, `<a id>`-Anker von
     `v0.44.0` **auf `v0.44.1` gewandert** (die v0.44.0-Zeile hat ihn verloren). ✓
  2. Alle `ghcr`-Pins gezogen: `README.de.md:163`, `README.md:162`, Handbuch
     (16 Vorkommen inkl. `--print-mk`-Illustration `DCHECK_IMAGE ?=`). ✓
  3. `CHANGELOG.md` — `[Unreleased]` sauber unter `[0.44.1] — 2026-07-17`
     geschnitten. ✓
  4. Prosa-Currency: Handbuch-Header-Stempel (1.30 ⇒ **1.31**, Software-Version
     `v0.44.1`); **bare-Tag-Beispiel** `:v0.44.1` in §Versionen und Tags gezogen
     (der bekannte Blind Spot des `versions`-Gates); §11-Verlaufszeile **1.31
     chronologisch unter 1.30** angefügt, nicht oben; neuer Feature-Abschnitt in §5
     ergänzt; `README.de.md` (kanonisch) und `README.md` synchron.
     `operations.md` korrekt **nicht** angefasst — weder neues Modul noch neue
     CLI-Option. ✓
  5. Digest-Pin als **Folge**-Commit nach dem Tag (`91f1a52`,
     `sha256:d2b1e53d…62427`), wie vorgeschrieben. ✓
- **SemVer-Klassifikation — ohne Befund (zur Absicht):** Patch ist für die
  *beabsichtigte* Änderung korrekt. Die Argumentationskette trägt: das Lastenheft
  verspricht die Expansion unqualifiziert, die Verengung „unmittelbar" stand allein
  in Rang 2 (`DC-FA-COV-001.a`), und Rang 2 ist fortschreibbar
  ([`AGENTS.md`](../../AGENTS.md) §2) — eine Schärfung, die eine bestehende
  Lastenheft-Zusage **herstellt**, ist kein CR und braucht keinen Lastenheft-Bump.
  Die Gegenthese („das Lastenheft verspricht unqualifiziert, also war die Verengung
  der Vertragsbruch") stützt genau diese Lesart. Das Argument ist **unabhängig** von
  F-1: F-1 ist ein Implementierungsdefekt, keine Fehlklassifikation.

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 1 | F-1 |
| MEDIUM | 3 | F-2, F-3, F-4 |
| LOW | 1 | F-5 |
| INFO | 0 | — |

---

## Verdikt

**REJECT** — ein HIGH und drei MEDIUM. Empfehlung: **v0.44.2**.

Begründung: F-1 ist keine Regression im Nebenpfad, sondern eine Umkehr der
Schutzrichtung. v0.44.1 hat einen Falschbefund beseitigt, der laut war (falsche
Waisen ⇒ Gate rot ⇒ der Konsument sieht es) und dabei einen Falschbefund eingebaut,
der still ist (falsche Deckung ⇒ Gate grün ⇒ der Konsument sieht **nichts**).
Der Slice selbst benennt in §4 genau diese Richtung als die gefährlichere, und
ADR-0039 wählt „genau ein Suffix" ausdrücklich, um Falsch-Expansionen zu vermeiden —
das gewählte Mittel (`[^)]*`) verfehlt sein eigenes Ziel. Ein `--require-complete`,
das bei realem Deckungs-Loch grün liefert, ist ein stiller Grün-Pfad in einem Gate;
das ist die HIGH-Klasse dieses Repos.

F-2 und F-1 haben **eine** Wurzel: der Range-Parser definiert „Link" anders als
`rules/markdown.go` es für `ids`/`links` bereits tut. ADR-0039 §Konsequenzen führt
„der Parser trägt jetzt ein Stück Markdown-Wissen" als bekannte Kosten — die Kosten
sind höher als dort veranschlagt, weil das zweite Markdown-Wissen dem ersten
**widerspricht**. Für die Behebung ist damit ein ADR-0039-Nachfolger (Supersedes)
oder eine `## Geschichte`-Ergänzung vor der Accept-Umstellung zu erwarten; F-3
(falsche Zusagen) und F-5 fallen mit demselben Zug.

F-4 ist der Grund, warum F-1 durchkam, und der ernsteste Prozess-Befund: die
Fitness-Funktion von ADR-0039 ist prosaisch formuliert und **nirgends mechanisiert**.
Die DoD hakt „Tests (positiv) … je einmal für `trace.coverage` … und für
`trace.cross-consistency`" ab — beide Tests existieren nicht. Das ist eine
Harness-Lüge-Vorstufe: die Abhakung behauptet eine Verriegelung, die der Baum nicht
enthält. (Die DoD-Prüfung ist Aufgabe der Verifikation, nicht des Reviews — der
Befund steht hier, weil die **fehlende Verriegelung** eines öffentlichen Vertrags
Reviewer-Kategorie MEDIUM ist.)

---

## Zur Abweichung „Review vor Release" (ausdrücklich beantwortet)

**Die Abweichung war nicht vertretbar — und sie hat genau das gekostet, wogegen die
Hausregel schützt.**

Die *Abwägung* war vernünftig: ein seit v0.41.0 ausgelieferter Falschbefund, der
Konsumenten fälschlich gatet, ist dringlich, und die Richtung des Fixes (weniger
Waisen) sah risikoarm aus. Genau diese Risiko-Einschätzung ist aber die, die der
Review hätte prüfen sollen — sie stammte vom Autor, und sie war falsch. Die
Dringlichkeit war zudem geringer als angenommen: der Defekt war drei Releases alt und
laut (rote Gates), der eingebaute ist still. Ein um Stunden verzögerter Fix eines
lauten Defekts ist billiger als ein sofort ausgelieferter stiller.

**Ist etwas durchgerutscht, das ein Review vorher gefangen hätte? Ja — F-1, und zwar
auf dem direktesten denkbaren Weg.** Der Review-Auftrag stellt unter Prüf-Schwerpunkt 2
wörtlich die Frage „Kann eine Range jetzt EXPANDIEREN, wo sie es vorher zu Recht nicht
tat?" und nennt unter Schwerpunkt 1 „verschachtelte Klammern in URLs" als zu prüfende
Eingabe. Beides zusammen **ist** F-1. Ein Review vor dem Tag hätte denselben
Zwei-Zeilen-Fixture gegen dasselbe Image gefahren und wäre auf dieselbe Zeile
gestoßen. Der Autor hatte die Klammer-Klasse sogar gesehen — Kommentar, Handbuch-Satz
und Testfall existieren — hat sie aber nur im `..`-Zweig durchgespielt und aus dem
gutartigen Ergebnis auf die ganze Klasse geschlossen. Das ist der klassische
Selbstbestätigungs-Fehler, gegen den ein kontext-getrennter Reviewer das Gegenmittel
ist: der Autor testet seine Hypothese, der Reviewer testet ihre Grenze.

Verschärfend: das Release fasst **beide** Konsumenten des geteilten Parsers an
(ADR-0039 Entscheidung 4) und ist in GHCR + `:latest` publiziert; jeder Konsument mit
`trace.coverage`/`ranges: true` zieht F-1 mit dem Patch, den er wegen der Dringlichkeit
zieht. Die Reichweite der Abweichung ist damit größer als die des Defekts, den sie
beheben wollte.

**Belastbare Lehre für die Hausregel:** die Dringlichkeits-Ausnahme sollte nicht an
der Schwere des Bestandsdefekts hängen, sondern an der **Lautstärke der beiden
Fehlerrichtungen**. Ein Fix, der einen Befund *entfernen* kann, verdient den Review
mehr als einer, der Befunde hinzufügt — denn nur die erste Richtung kann sich selbst
verstecken. slice-073 ist der Beleg: der Bestandsdefekt meldete sich drei Releases
lang von selbst (und wurde am Ende auch so gefunden), F-1 hätte das nie getan.
