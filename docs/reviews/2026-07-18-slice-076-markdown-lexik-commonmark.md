# Review — slice-076: Markdown-Lexik an CommonMark/GFM angleichen

**Gegenstand:** slice-076 (Trennzelle `-+` + Fence-Infozeilen-Regel).
**Datum:** 2026-07-18. **Reviewer:** unabhängig, kontext-getrennt (keine
Implementierungs-Session-Kenntnis). **Skill:** `.harness/skills/reviewer.md`
v1.2.0.

**Commit-Range:** `origin/main..HEAD` = `2dbf4e1` (feat), `98169c0`/`b4f6e26`
(plan/roadmap). **Betroffene Anforderungen:**
[`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 3,
[`DC-FA-LINK-001.a`](../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
Schritt 1, [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus).
**ADR:** [ADR-0042](../plan/adr/0042-markdown-lexik-folgt-commonmark.md) (Proposed).

**Methodik:** Diff-Lektüre; Abgleich beider Regeln gegen Spec-Wortlaut und
CommonMark/GFM; Aufruf-Ketten-Analyse (`markdownTables`/`tableHeaderAt`,
`proseLines`-Konsumenten); **empirische Mutationsprüfung** — jede der drei Regeln
in einer Wegwerf-Kopie einzeln zurückgedreht und `go test ./...` über das
`--target test`-Image gefahren (Repo unangetastet); **Dogfood-Gates** `make
test`/`doc-check`/`trace`/`doc-complete` real gelaufen.

---

## Findings

### R-F-1 · Info-Zeilen-Regel doppelt implementiert
- **Kategorie:** LOW
- **Quelle:** Maintainability (latente Divergenz-Wiederkehr; vgl. ADR-0042
  Konsequenz (a) „zwei Automaten, zwei Verhalten")
- **Pfad:** `internal/hexagon/core/rules/markdown.go:29-41` (`fenceToggle`) und
  `internal/hexagon/core/app/trace_table.go:294` (Inline-Zweig in
  `markdownTableLines`)
- **Befund:** Dieselbe CommonMark-Infozeilen-Semantik („Backtick-Fence mit
  Backtick im Rest ist kein Öffner") steht in zwei getrennten Code-Stellen, weil
  die beiden Automaten unterschiedliche Trag-Strukturen haben (`fenceToggle`
  liefert ein bool für den naiven Toggle; `markdownTableLines` führt
  `fenceChar`/`fenceLen` für den längenabgeglichenen Schluss). Es gibt keinen
  gemeinsamen Prädikat-Aufruf und keinen Test, der beide Implementierungen
  gegeneinander bindet — sie sind heute deckungsgleich, driften aber bei einem
  einseitigen Edit auseinander. Genau die trace≠links-Divergenz, die dieser Slice
  schließt, kann so still zurückkehren.
- **Verifizierbar:** ja — kein Test assertiert die Gleichheit der beiden Öffner-
  Tests; nur je ein Automat-eigener Test (`TestCLI076_…VerdecktTabelleNicht` vs.
  `…VerdecktLinkNicht`) existiert. Ein künftiger Edit an nur einer Seite bliebe
  grün.
- **Übergabe (nicht Teil des Befunds):** ein gemeinsames, exportiertes Öffner-
  Prädikat erwägen, das beide Automaten aufrufen; oder einen Test, der für ein
  Zeilen-Sample beide Wege identisch entscheiden lässt.

### R-F-2 · Regel auch in `diagramFenceLines` — dritter Automat, ohne Beleg/Test
- **Kategorie:** INFO
- **Quelle:** ADR-0042 Entscheidung 3 („Die Grenze ist das Gemessene") /
  Maintainability
- **Pfad:** `internal/hexagon/core/rules/markdown.go:90` (`diagramFenceLines` ruft
  `fenceToggle`)
- **Befund:** Der geteilte `fenceToggle` wird nicht nur von `proseLines`, sondern
  auch von `diagramFenceLines` (Modul `diagrams`, `DC-FA-DIAG-001.a`) benutzt —
  ein dritter Automat, den weder DoD noch ADR nennen und für den kein Realfall in
  den 522 Messdateien gemessen wurde. **Zwei ehrliche Sichten:** (a) Konsistenz —
  ein gemeinsames Öffner-Prädikat hält die Fence-Definition über *alle* naiven
  Automaten gleich; eine abweichende `diagrams`-Definition wäre selbst eine
  latente Inkonsistenz derselben Klasse. (b) Scope — ADR-0042 mahnt „keine
  weiteren Angleichungen auf Verdacht"; die Verhaltensänderung für `diagrams`
  (ein `` ```mermaid`x` ``-artiger Öffner toggelt nicht mehr) ist unbelegt und
  ungetestet. Sicht (a) überwiegt (die Angleichung ist *dieselbe* Regel, nicht
  eine *weitere* CommonMark-Regel), daher INFO, nicht MEDIUM.
- **Verifizierbar:** ja — eine Mutation, die *nur* die `diagramFenceLines`-
  Aufrufstelle auf den alten `HasPrefix("```")||HasPrefix("~~~")`-Test
  zurückdreht, wird von keinem Test gefangen (Mutations-Lücke am dritten
  Automaten).
- **Übergabe (nicht Teil des Befunds):** entweder einen `diagrams`-Test mit
  Backtick-Infozeile ergänzen, oder in ADR-0042/Slice notieren, dass die Regel
  bewusst auf alle `fenceToggle`-Konsumenten wirkt.

### R-F-3 · Feat-Commit editiert dattierte historische Review-Doku (Zeile 179)
- **Kategorie:** LOW
- **Quelle:** Doku-Currency / Maintainability
- **Pfad:** `docs/reviews/2026-06-19-slice-030-suggest-config-ai-harness.md:179`
- **Befund:** Der Feat-Commit ändert in einem **dattierten Alt-Review** (2026-06-19)
  die Prosa ```` ```yaml-Fence (`spezifikation.md:145-180`) ```` → `` `yaml`-Fence
  (…) `` — genau die „Zeile 179", die ADR-0042 (§Kontext B, §Fitness-Funktion),
  die Spezifikation und der Slice als **realen Beleg** für Defekt B zitieren. Der
  Beleg-Zeiger verweist danach auf eine korrigierte Zeile; die Fundstelle in
  ihrer belegten Gestalt existiert nur noch in den synthetischen Test-Fixtures.
  Der Edit ist rein formal (Aussage unverändert, kein Befund/Verdikt berührt) und
  dient der Dogfood-Hygiene, aber (i) er stalet den zitierten Beleg und (ii) das
  Ändern eines historischen Review-Artefakts im Feat-Commit ist eine milde
  Prozess-Unschärfe (die Regel-Analyse selbst braucht ihn nicht — der neue Code
  liest die alte Zeilenform ohnehin korrekt als Prosa).
- **Verifizierbar:** ja — `git show 2dbf4e1 -- docs/reviews/2026-06-19-…md` zeigt
  die Inhaltsänderung; `doc-check` bleibt grün (243 Dateien, 0 Befunde), d. h.
  keine funktionale Regression, nur Currency.
- **Übergabe (nicht Teil des Befunds):** entscheiden, ob historische Review-Docs
  form-editierbar sind; falls nein, den Fix in ein separates Hygiene-Commit
  ziehen und die Beleg-Zitate auf die synthetische Fixture-Form umhängen.

---

## Negativbefunde (geprüft, ohne Befund)

- **Trennzelle `^:?-+:?$` vs. Spec/GFM:** exakt deckungsgleich mit
  `spec/spezifikation.md:371`; `-{3,}` ⊂ `-+`, also rein erweiternd. Kein neuer
  False-Positive: `tableHeaderAt`/`splitPipeTableLine` verlangen, dass die
  Trennzeile **selbst** eine Pipe-Zeile ist und **jede** Zelle `-+` matcht — eine
  thematische Trennlinie `---` (ohne Pipe) oder eine Datenzeile mit einer
  `-`-Zelle wird nicht zur Trennzeile. `make test` grün ⇒ kein bestehender
  Negativtest gegen die alte Strenge brach. **Geprüft, ohne Befund.**
- **Infozeilen-Regel vs. CommonMark:** ``strings.IndexByte(trimmed[run:], '`')``
  prüft den Rest hinter der **vollen** Backtick-Folge (kein Off-by-one); leere
  Infozeile (` ``` `), reine Backtick-Zeilen (`` `````` ``) und ≥4-Backtick-Öffner
  öffnen korrekt; `~~~`-Asymmetrie korrekt (Tilde-Zweig kurzschließt in beiden
  Automaten). **Geprüft, ohne Befund.**
- **Regel in beiden Automaten (DoD):** `proseLines` (via `fenceToggle`, geteilt
  von links/ids/codepaths/matrix/hostpaths/spans/anchors) UND `markdownTableLines`
  (trace). Öffner-Seite konsistent; die Rest-Divergenz naiver-Toggle-vs-strikter-
  Schluss ist **vorbestehend** und wird durch die Regel sogar leicht **reduziert**
  (Fall `` ```x` `` im Fence stimmt jetzt in beiden überein), nicht vergrößert.
  **Geprüft, ohne Befund.**
- **`else if`-Umbau in `markdownTableLines`:** verhaltenserhaltend für den
  Schluss-Pfad (im Ursprung schloss der erste `if fenceLen==0`-Zweig immer mit
  `continue`, der zweite `if` war bei `fenceLen==0` unerreichbar; bei blockierter
  Info-Zeile fällt der Code korrekt auf `out[i].prose = fenceLen==0` = true
  durch). Bestehende Fence-Schluss-Tests im `app`-Paket bleiben grün.
  **Geprüft, ohne Befund.**
- **Mutations-Härte (empirisch, nicht zugesagt):** drei Einzel-Rückdrehungen im
  Wegwerf-Baum gefahren:
  - `-+`→`-{3,}` ⇒ **nur** `TestCLI076_TrennzeileEinBindestrich` kippt (Exit 2).
  - Tabellen-Automat-Regel entfernt ⇒ **nur**
    `TestCLI076_FenceInfozeileVerdecktTabelleNicht` kippt (Exit 2); Prosa-Seite
    grün.
  - `fenceToggle`-Regel entfernt ⇒ `TestCLI076_FenceInfozeileVerdecktLinkNicht`
    (Exit 0 statt 1 — der stilles-Grün-Pfad kehrt sichtbar zurück) **und**
    `TestPreprocessMarkdown_FenceInfozeileMitBacktickIstFliesstext` kippen;
    Tabellen-Seite grün.
  Jede Regel ist isoliert gepinnt; die R3-F-2-Lehre ist erfüllt. **Geprüft,
  bestätigt.**
- **Konsumenten-Ebene ≠ trace (DoD):** `TestCLI076_FenceInfozeileVerdecktLinkNicht`
  fährt das `links`-Modul (`--disable anchors`) und assertiert
  `docs/a.md:4\tfehlt.md\ttarget-missing`; Zeilennummer korrekt (Fixture-Zeile 4).
  **Geprüft, ohne Befund.**
- **Determinismus (DC-QA-02):** keine neuen Maps/Iterationen; regex- und
  Byte-Scans sind deterministisch; `trace`-Ausgabe bleibt über `sort.Strings`
  sortiert. `make trace`/`doc-complete` byte-stabil grün. **Geprüft, ohne Befund.**
- **Verhaltensänderung an d-checks eigenen Gates:** `make test`, `doc-check` (243
  Dateien, 0 Befunde), `trace` und `doc-complete` (je 43 Anforderungen, 0 Waisen)
  real gelaufen und grün — die neu sichtbare Prosa/Tabellen in der eigenen Doku
  brechen keinen Dogfood-Gate. **Geprüft, ohne Befund.**
- **Fehlerpräzedenz / Edge Cases:** leere Datei (kein Panic; `trimmed[run:]`
  slice-to-len sicher), Fence am Dateiende (unclosed → bis EOF, vorbestehend),
  4-Leerzeichen-/Tab-Einrückung und naiv-vs-strikt-Schluss sind von ADR-0042
  §Konsequenzen **bewusst offen** gelassen (dokumentierte Won't-Fix, kein
  Realfall). **Geprüft, ohne Befund.**
- **ADR-0005-Importregeln / Netz / Gate-Suppression:** unberührt — Änderungen rein
  in `core/rules` und `core/app`, keine neuen Ports/Adapter, kein Netz, keine
  Gate-Ausnahme. **Geprüft, ohne Befund.**

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 0 | — |
| LOW | 2 | R-F-1, R-F-3 |
| INFO | 1 | R-F-2 |

---

## Verdikt: **ACCEPT-WITH-NITS**

Beide Regeln sind spec- und CommonMark/GFM-treu, rein erweiternd (Trennzelle) bzw.
schließen einen belegten, modulübergreifenden stilles-Grün-Pfad (Infozeile). Die
DoD-Kern-Forderung — Regel in **beiden** Automaten, **isolierte** Mutations-Härte,
Beleg auf **Konsumenten-Ebene ≠ trace** — ist erfüllt und wurde hier **empirisch**
gegengeprüft (drei Einzel-Mutationen kippen je genau den zuständigen Test, die
Rückkehr des Exit-0-Blindpfads inklusive). Kein HIGH/MEDIUM-Befund; die Dogfood-
Gates bleiben grün.

Die drei Nits blockieren den Merge nicht, verdienen aber eine Entscheidung vor der
Closure: die **Doppel-Implementierung** der Infozeilen-Regel (R-F-1) und ihre
stille **Ausweitung auf den `diagrams`-Automaten ohne Test** (R-F-2) sind beide
latente Wiederkehr-Vektoren genau der Divergenz-Klasse, die dieser Slice
adressiert; der **Edit an der historischen Review-Doku** (R-F-3) ist formal
legitim, stalet aber den in ADR/Spec zitierten „Zeile 179"-Beleg und gehört
prozessual eher in ein eigenes Hygiene-Commit.
