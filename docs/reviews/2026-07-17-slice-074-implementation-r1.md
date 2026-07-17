# Review slice-074 (R1) — Kommentar-Suffix in Tabellenzeilen

**Datum:** 2026-07-17 · **Rolle:** unabhängiger Reviewer (kontext-getrennt,
nicht der Autor) · **Lauf:** R1, **vor** Release v0.45.2.

**Gegenstand:** `44a5201` (Fix + Tests), `a594144` (doc-first: ADR-0040 + Spec),
`4f76133`/`ae847f8` (Lifecycle). Diff `0379970..44a5201`.

**Kontext-Hinweis:** Während des Laufs bewegte sich `HEAD` von `44a5201` auf
`0e26e27` (`docs(trace): … slice-075`, Parallel-Session). `git diff --stat
44a5201..HEAD -- internal/` ist leer — der Prüfgegenstand ist unberührt.

**Quellen:** [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 3, [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency),
[ADR-0040](../plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md),
[ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md),
[`AGENTS.md`](../../AGENTS.md) §3.

**Verifikations-Basis:** `make build` (HEAD) vs. `make build` in einem Worktree
auf `44a5201~1` (Tag `d-check:prefix`), Läufe `docker run --rm --network none -v
<fixture>:/repo:ro -w /repo <image> --trace [--require-complete]`. Mutationen
via `make test` in einem Scratch-Worktree (Produktivbaum unverändert).

---

## Findings

### F-1 · HIGH · Kommentar-Suffix auf der **Headerzeile** macht eine GFM-renderbare Tabelle unsichtbar — Anforderungen, Waisen und Kreuzverweis-Differenzen verschwinden still, das Gate wird grün

**quelle:** [ADR-0040](../plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md)
§Konsequenzen („**Kein Lauf wird stiller**: der Guard feuert weiter, nur nicht
mehr auf die eigene Direktive"); [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 3; Reviewer-Anker HIGH („Stilles-Grün-Pfad in einem Gate — Harness-Lüge").

**pfad:** `internal/hexagon/core/app/trace_table.go:341`
(`dropCommentSuffix` **innerhalb** `splitPipeTableLine`) im Zusammenspiel mit
`internal/hexagon/core/app/trace_table.go:100-102`
(`len(delimiter) != len(header)` ⇒ keine Tabelle).

**befund:** `dropCommentSuffix` sitzt in `splitPipeTableLine` und greift daher
auch auf der **Headerzeile**. `tableHeaderAt` vergleicht danach die Zellenzahl
von Header und Trennzeile. Trägt die Headerzeile ein Kommentar-Suffix, sinkt die
Header-Zahl auf N, während die Trennzeile bei N+1 bleibt — die Tabelle wird
**gar nicht mehr als Tabelle erkannt** und kommentarlos übersprungen. N+1 ist
dabei **die einzige GFM-renderbare Form**: GFM zählt den Kommentar als Zelle und
verlangt Gleichstand von Header- und Trennzeile. Solange **eine andere** Tabelle
relevant bleibt, greift weder der `foundTable`-Guard aus
[ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md) noch die
Vakuitäts-Stufe aus ADR-0038 — der Lauf endet Exit 0.

Belegt gegen das echte Image, Fixture `fx-c` (`format: table`; Tabelle 1 sauber,
Tabelle 2 mit Header-Suffix und Trennzeile `|---|---|---|`; `F-1` gedeckt,
`F-2`/`F-3` echte Waisen):

```text
=== PRE-FIX (44a5201~1), --trace --require-complete ===
d-check: 2 Requirements-Waise(n) ohne referenzierenden Slice (--require-complete)
| F-1 | Alpha | ADR-0001 | slice-001 | ok    |
| F-2 | Beta  | —        | —         | WAISE |
| F-3 | Gamma | —        | —         | WAISE |
3 Anforderung(en), 2 Waise(n).                      EXIT=1

=== POST-FIX (HEAD), --trace --require-complete ===
| F-1 | Alpha | ADR-0001 | slice-001 | ok |
1 Anforderung(en), 0 Waise(n).                      EXIT=0
```

Der Guard war **vor** dem Fix laut (`Tabellenzeile 13 hat 2 statt 3 Zellen`,
Exit 2, Fixture `fx-a`) und ist **nach** dem Fix stumm. Die Richtung ist damit
laut → still, nicht umgekehrt.

Derselbe geteilte Reader trägt den Defekt in die zweite Sicht. Fixture `fx-x`
(`trace.cross-consistency`, zwei Rück-Tabellen, die zweite mit Header-Suffix):

```text
=== PRE-FIX  ===  | F-2 | C-2 | Rück-Kante, ohne RTM-Eintrag | spec/architecture.md:11 |
                  1 Differenz(en).
=== POST-FIX ===  0 Differenz(en).
```

Die Rück-Kante `C-2 → F-2` verschwindet lautlos; `bindBackwardTables` meldet
nichts, weil die **erste** Rück-Tabelle relevant bleibt.

Die Inversion ist scharf: die **nicht** GFM-renderbare Form (Trennzeile mit N
Zellen, Fixture `fx-d`) liest der Fix korrekt; die GFM-renderbare Form macht er
unsichtbar. Der Slice fordert die Abstreifung für „Header-, Trenn- und
Datenzeilen gleichermaßen" (§3 DoD) — genau diese Gleichbehandlung erzeugt den
Befund; die Spec-Regel steht in §3 „Tabellen-Lexik" und adressiert die
Header-/Trennzeilen-Kopplung nicht.

**verifizierbar:** ja — `docker run --rm --network none -v fx-c:/repo:ro -w /repo
d-check:latest --trace --require-complete` liefert Exit 0/„0 Waise(n)", das
Image auf `44a5201~1` Exit 1/„2 Waise(n)". Als Regressionstest würde eine
Konsumenten-Ebene-Fixture mit Header-Suffix und Trennzeile N+1 den Befund
dauerhaft binden.

---

### F-2 · MEDIUM · Bestandsregression in der lauten Richtung: eine Datenzeile, deren legitime letzte Zelle ganz aus einem Kommentar besteht, war bis v0.45.1 lesbar und ist jetzt Exit 2

**quelle:** [ADR-0040](../plan/adr/0040-kommentar-suffix-in-tabellenzeilen.md)
§Konsequenzen („Verhaltensänderung für Bestandskonsumenten: eine Quelle, die
**heute** mit `Tabellenzeile N hat X statt Y Zellen` abbricht, läuft danach
durch") — die Zusage ist einseitig formuliert; [`DC-QA-04`](../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-nahe
Erkennungs-Differenz zum ausgelieferten Stand.

**pfad:** `internal/hexagon/core/app/trace_table.go:359-364`.

**befund:** `dropCommentSuffix` unterscheidet nicht zwischen einem Suffix
**hinter** der letzten Pipe und einer vollwertigen letzten **Zelle**, deren
Inhalt zufällig ganz ein HTML-Kommentar ist — beide erreichen den Splitter als
identische Zeichenkette. Eine 3-spaltige Tabelle mit der Zeile
`| F-2 | Beta | <!-- offen --> |` wurde von v0.45.1 korrekt gelesen und bricht
jetzt ab (Fixture `fx-e`):

```text
=== PRE-FIX  ===  2 Anforderung(en), 2 Waise(n).                     EXIT=0
=== POST-FIX ===  d-check: error: trace.requirements: Tabellenzeile 6 hat 2 statt 3 Zellen
```

Die Verhaltensänderung ist damit **beidseitig** (kaputt→läuft **und**
läuft→kaputt), nicht einseitig wie in der ADR zugesagt. Das berührt zugleich die
SemVer-Einstufung: ein Patch, der eine bislang gültig gelesene Quelle in Exit 2
kippt, ist keine reine Fehlerbehebung. Laut ist der Pfad immerhin — daher
MEDIUM, nicht HIGH.

**verifizierbar:** ja — `docker run --rm --network none -v fx-e:/repo:ro -w /repo
<image> --trace` gegen beide Images.

---

### F-3 · MEDIUM · Der Header-/Trennzeilen-Zweig ist in **beide** Richtungen ungepinnt, obwohl die DoD ihn ausdrücklich als Ziel benennt

**quelle:** Reviewer-Anker MEDIUM („fehlende Negativtests bei neuem öffentlichen
Vertrag"); slice-074 §3 DoD („Implementierung: … Header-, Trenn- und Datenzeilen
gleichermaßen").

**pfad:** `internal/hexagon/core/app/trace_coverage_test.go:220-247`,
`internal/hexagon/core/app/trace_cross_test.go:587-615`.

**befund:** Die drei vom Autor behaupteten Grenzen sind tatsächlich gepinnt —
nachgemessen per Mutation gegen `make test`:

| Mutation | kippt |
|---|---|
| M1 Abstreifung entfernt (`return cells, true`) | `TestSplitPipeTableLineKommentarSuffix/Suffix_*` **und** `TestCrossConsistencyKommentarSuffix` |
| M2 „nur am Ende" aufgehoben (jede Kommentarzelle fällt) | `…/Kommentar_in_der_Zeilenmitte` |
| M3 „genau einer" aufgehoben (`^<!--.*-->$`) | `…/zwei_Trailing-Kommentare` |

Ungepinnt ist die **vierte** Grenze, die den Defekt aus F-1 trägt: kein Test
konstruiert eine **Headerzeile** mit Kommentar-Suffix. Beide Kommentar-Tests
setzen den Marker auf eine **Datenzeile** (`crossFS(header, row)`), der
Splitter-Test kennt den Zeilen-Typ konstruktionsbedingt nicht.

Verschärfend pinnt der Splitter-Test die **Einbau-Stelle** statt des Vertrags:
Mutation M4 (Abstreifung nur in `consumeTableRows`, also nicht auf Header/Trenn-
zeile) lässt den Konsumententest `TestCrossConsistencyKommentarSuffix` **grün**
und kippt ausschließlich `…/Suffix_ohne_Schluss-Pipe` und
`…/Suffix_mit_Schluss-Pipe`. Die Testmenge hält damit die defekte Platzierung
fest und weist eine Variante ab, die den motivierenden Realfall weiterhin löst.
Die Konsumenten-Ebene ist gegenüber dem Vorlauf zwar nachgezogen
(`TestCrossConsistencyKommentarSuffix`, `TestCrossConsistencyZellenzahlGuardBleibtScharf`),
deckt aber nur `trace.cross-consistency` — für `trace.requirements.format: table`
(ausgeliefert seit v0.43.0, der ältere der beiden Konsumenten) gibt es **keinen**
Konsumententest, obwohl die DoD „je einmal für … und …" verlangt.

Der Slice benennt in §4 selbst die **Dogfood-Lücke** (d-checks eigene Tabellen
tragen keinen Trailing-Marker) — die Selbstprüfung konnte den Zweig also
bekanntermaßen nicht abfangen, und `make gates` bleibt grün.

**verifizierbar:** ja — die Mutationsläufe oben sind mit `make test` in einem
Worktree reproduzierbar.

---

### F-4 · LOW · Doppelte Backlog-Zeile in der Roadmap

**quelle:** Maintainability (Doku-Drift).

**pfad:** `docs/plan/planning/in-progress/roadmap.md:60` und `:62`.

**befund:** Der Diff fügt die Zeile „**Im Backlog (`next/`), auf Aufnahme in eine
Welle wartend:** derzeit keiner." ein, obwohl sie zwei Zeilen darüber bereits
steht; beide Vorkommen sind wortgleich.

**verifizierbar:** ja — `grep -n "Im Backlog" docs/plan/planning/in-progress/roadmap.md`
liefert zwei Treffer.

---

### F-5 · INFO · Ein Kommentar mit Pipe oder mit ungerader Backtick-Parität ist kein Suffix — die Direktiven-Konvention erlaubt beides im Grund-Text

**quelle:** [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 3 („ganzzellig"); dokumentationswürdige, undokumentierte Annahme.

**pfad:** `internal/hexagon/core/app/trace_table.go:320-341`.

**befund:** Splitter-Regeln laufen **vor** `dropCommentSuffix`. Eine unescapte
Pipe im Grund-Text zerteilt das Suffix in mehrere Zellen, bevor die Regel greift
— `| F-1 | Alpha |` + `| F-2 | Beta | <!-- d-check:ignore (a | b) -->` ergibt
`Tabellenzeile 6 hat 4 statt 2 Zellen` (Exit 2, Fixture `fx-f`). Das ist mit dem
Wortlaut „ganzzellig" vereinbar und laut, aber die Konvention
`<!-- d-check:ignore (Grund) -->` stellt an den Grund keine Zeichen-Einschränkung;
weder Spec noch Handbuch nennen sie. Backticks im Grund sind unkritisch, solange
ihre Parität gerade ist (der Regex-Body akzeptiert sie über `[^-]`); bei ungerader
Parität greift die vorbestehende Code-Span-Klasse, unverändert durch slice-074.

**verifizierbar:** ja — Fixture `fx-f` gegen das HEAD-Image.

---

## Negativbefunde (geprüft, ohne Befund)

- **`htmlCommentCell`-Regex, Engegrad in der gefährlichen Richtung:** geprüft, ohne
  Befund. `^<!--(?:[^-]|-[^-]|--[^>])*-->$` lässt im Body kein `-->` zu (jede
  Alternative scheitert an `-->`; `$` erzwingt den letzten `-->` als Terminator).
  Eine matchende Zeichenkette ist damit **exakt** ein Kommentar — es gibt keine
  echte Zelle, die fälschlich matcht. Durchgespielt: `<!---->` (Body leer, matcht),
  `<!-- -->`, `<!-- a -- b -->`, `<!-- a -> b -->` (matchen, korrekt),
  `<!--->` und `<!-- x --->` (matchen **nicht**, bleiben Zelle — konservativ),
  `-->` ohne `<!--` (matcht nicht), leere Zelle `""` (matcht nicht, überlebt).
- **Zwei/mehrfache und verschachtelte Kommentare:** geprüft, ohne Befund — kein
  Suffix, laufen in den Zellenzahl-Guard (M3 pinnt es).
- **Zeile, die nur aus einem Kommentar besteht:** geprüft, ohne Befund.
  `| <!-- x -->` liefert `([], true)`; `len(cells)=0 != len(header)` ⇒ `badLine`,
  und `isTableDelimiter(nil)` gibt über `len(cells) > 0` `false` zurück — kein
  Phantom-Table aus zwei leeren Zeilen. Kein Panic-Pfad: `t.rows` wird nur bei
  `len(cells) == len(t.header)` gefüllt, `cols.id/text/modality` bleiben also
  indexsicher.
- **Kommentar über Zeilengrenzen:** geprüft, ohne Befund — `| a | b | <!--` endet
  auf einer Nicht-Kommentar-Zelle, läuft in den Guard. Block-Kommentare um ganze
  Zeilen sind vorbestehend außerhalb des Reader-Modells (nur Fences werden
  maskiert), von slice-074 unverändert.
- **`\|`-Escape und `hasUnescapedTrailingPipe`:** geprüft, ohne Befund —
  `dropCommentSuffix` läuft nach dem Leading-/Trailing-Strip; `| a | b | <!-- x --> |`
  und ohne Schluss-Pipe liefern identisch `[a b]` (beide Formen getestet).
- **Code-Span-Erkennung (`consumeBackticks`, `hasClosingBacktickRun`):** geprüft,
  ohne Befund für slice-074 — die Backtick-Paritäts-Klasse ist vorbestehend, laut
  (Zellen verschmelzen ⇒ Guard) und wird durch die neue Regel weder verschoben
  noch stiller.
- **Realdatenbeleg grid-gym:** geprüft, ohne Befund — `spec/architecture.md:913`
  ist eine **Datenzeile** bei sauberem 3/3-Header; der Fix trifft den
  motivierenden Fall. (Nur gelesen, nichts geändert.)
- **ADR-0005-Import-Regeln / Modul-Layout:** geprüft, ohne Befund — keine neuen
  Imports, die Änderung bleibt in `core/app`.
- **[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  (Netz/Seiteneffekte):** geprüft, ohne Befund — nur `regexp`/`strings`, alle
  Verifikationsläufe mit `--network none`.
- **Gate-Suppression / Schwellen-Senkung (AGENTS.md §3.6):** geprüft, ohne Befund —
  keine `//nolint`, keine Gate-/Schwellen-Änderung.
- **ADR-Immutabilität (AGENTS.md §3.5):** geprüft, ohne Befund — ADR-0040 ist neu
  und `Proposed`; bestehende `Accepted`-ADRs unangetastet.
- **„Defekt-Fix, kein CR" (Vertragstreue):** geprüft, **REFUTED** als Finding und
  haltbar. [`DC-FA-REQ-001`](../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen)
  nennt Zellen (ID-Zelle, Zelleninhalt, Rollen-Header), sagt aber **keine
  Zellenzahl** und keinen Zellenzahl-Guard zu; die Zusage steht allein in
  [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
  Schritt 5 (Rang 2, fortschreibbar). Der Weg über Spec + ADR ist korrekt.
- **ADR-0040 ↔ Spec-Widerspruch:** geprüft, ohne Widerspruch im Wortlaut — beide
  sagen „genau einer, nur am Ende, nur ganzzellig". Beide schweigen jedoch zur
  Header-/Trennzeilen-Kopplung (Ursache von F-1) und zur Rückrichtung der
  Verhaltensänderung (F-2).
- **Referenz-Richtung (SDP), Marker-Ehrlichkeit:** geprüft, ohne Befund — ADR-0040
  nennt slice-074 nur in `## Geschichte` als Provenance („Umsetzender Slice"),
  nicht als Entscheidungsgrundlage.
- **CHANGELOG / Benutzerhandbuch fehlen im Diff:** geprüft, **kein Finding** — die
  Commit-Grenzen-Konvention des Repos legt Handbuch/CHANGELOG in den
  Release-Prep-Commit; die DoD führt „Nutzerdoku" und „Release" als eigene,
  offene Punkte. Der Review läuft vertragsgemäß davor.
- **DoD-Abhakung / Gate-Lauf-Bestätigung:** nicht geprüft — Rolle der
  Verifikation, nicht des Reviews (Reviewer-Skill §Anti-Pattern).

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 1 | F-1 |
| MEDIUM | 2 | F-2, F-3 |
| LOW | 1 | F-4 |
| INFO | 1 | F-5 |

---

## Verdikt

**BLOCK — v0.45.2 darf nicht getaggt werden.**

F-1 ist ein Stilles-Grün-Pfad in einem Gate und trifft den HIGH-Anker des
Reviewer-Skills unmittelbar: ein Lauf, der vor dem Fix laut Exit 2 meldete,
endet danach Exit 0 und verschweigt zwei echte Waisen bzw. eine echte
Kreuzverweis-Differenz. Er trifft beide Konsumenten, die der Slice als „ein Fix,
zwei Konsumenten" adressiert — nur in der Gegenrichtung. Damit ist die zentrale
Zusage von ADR-0040 („Kein Lauf wird stiller") **falsifiziert**, und die
Kontext-Eskalation greift zusätzlich (Gate-Pfad).

Es ist zudem die **dritte Wiederholung derselben Klasse in derselben
Code-Region** in dieser Woche: eine lexikalische Verengung im geteilten Reader,
deren gutartiger Zweig geprüft und deren Gegenzweig vom Ergebnis her erschlossen
wurde — hier der Datenzeilen-Zweig geprüft, der Header-/Trennzeilen-Zweig nicht,
obwohl die DoD ihn ausdrücklich mit einschließt. Nach dem Reviewer-Skill ist das
ein **Steering-Loop-Signal**: die Konsequenz gehört über den Einzelfix hinaus in
Guide/Sensor (der Slice benennt in §4 mit der Dogfood-Lücke bereits den passenden
Ansatzpunkt, und sein eigener §4-Satz „ein drittes Mal wäre ein Signal, dass die
lexikalischer-Splitter-Annahme selbst falsch ist" ist mit diesem Befund
eingetreten).

F-2 blockiert zusätzlich als Bestandsregression und stellt die SemVer-Patch-
Einstufung in Frage; F-3 blockiert als fehlender Negativtest am neuen
öffentlichen Vertrag. F-4/F-5 sind nicht release-blockierend.

Der motivierende Realfall (grid-gym `architecture.md:913`) wird vom Fix korrekt
gelöst — der Befund betrifft die **Verallgemeinerung** auf Header- und
Trennzeilen, nicht die Absicht des Slices.
</content>
