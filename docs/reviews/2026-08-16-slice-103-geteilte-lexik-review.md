# Review-Report: slice-103 — Geteilte Lexik, drei Konsumenten — 2026-08-16

**Review-Art:** Code — geprüft wird der ausgelieferte Diff gegen die Verträge, die
er selbst ändert, und gegen die Leitfrage des Gegenstands: *ist die Klasse
geschlossen?* Nicht geprüft wird die DoD-Abhakung (getrennter Kontext, Verifikation).

**Gegenstand:** Commit-Range `d2aaf90..6461bd6` (vier Commits: Wellen-Eröffnung ·
Messung · Vertrag · Implementierung); Arbeitsbaum-Stand `6461bd6`.

**Skill:** `.harness/skills/reviewer.md` @ 1.4.0 ·
**Modell:** claude-opus-5[1m] · **Datum:** 2026-08-16

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`DC-FA-CITE-001`](../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in),
  [`DC-FA-VER-001`](../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
  [`DC-FA-PIN-001`](../../spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in),
  [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  und [`DC-FA-ANCH-001`](../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors)
  (Lastenheft-Fassung 0.58.0)
- §[`DC-FA-CITE-001.a`](../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations)
  Schritt 2, §[`DC-FA-VER-001.a`](../../spec/spezifikation.md#dc-fa-ver-001a--versions-pin-konsistenz-versions)
  Schritt 1, §[`DC-FA-PIN-001.a`](../../spec/spezifikation.md#dc-fa-pin-001a--content-pin-gegen-inhaltlichen-drift-pins)
  Schritt 2, §[`DC-FA-ANCH-001.b`](../../spec/spezifikation.md#dc-fa-anch-001b--inline-html-anker)
  und §[`DC-FA-LINK-001.a`](../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
  Schritt 2
- [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  (Entscheidungen 1–4, Alternativen, Konsequenzen, Re-Evaluierungs-Trigger),
  [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md),
  [ADR-0019](../plan/adr/0019-versions-pin-fence-ausnahme.md),
  [ADR-0020](../plan/adr/0020-content-pin-fence-ausnahme.md),
  [ADR-0024](../plan/adr/0024-vcs-immutable-gate.md)
- [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  (Spiegel einer Semantik-Änderung), [`AGENTS.md`](../../AGENTS.md) §3 Hard Rules,
  [`CLAUDE.md`](../../CLAUDE.md)
- Der geprüfte Slice-Plan und sein Wellendokument (§3 Befunde, §4 Abnahme-Punkte,
  §4a Messung) sowie der
  [Vorgänger-Report R3](2026-08-09-slice-101-fence-review-r3.md) als Maßstab für
  R-2, R-3 und R-7

**Läufe dieses Reviews.** Alle Fixtures in einem Scratch-Verzeichnis außerhalb des
Repos, alle Läufe netzlos und read-only. Gefahren: `make build`, `make test`
(sieben Mutationsläufe), `make gates` sowie rund 30 Fixture-Läufe gegen drei
Images — den HEAD-Build, zwei gezielt mutierte Builds und das **veröffentlichte**
Vor-Bild `v0.52.0` als Gegenprobe für die Richtungsfrage. `make gates` grün
(372 Dateien, 0 Befunde). Für jede Mutations-Gegenprobe wurde die Repo-Datei
**vor** dem Eingriff in das Scratch-Verzeichnis kopiert und danach aus der Kopie
zurückgeschrieben (**kein** `git checkout`). **Der Arbeitsbaum ist
wiederhergestellt:** `git status --short` und `git diff --stat` sind leer, und der
Neubau liefert dieselbe Image-ID wie vor dem ersten Eingriff.

---

## Findings

### L-1 — Die Anker-Frage ist nur zur Fence-Hälfte vereinheitlicht; der neue Vertrag sagt volle Parität zu, und ein Lauf gibt weiter zwei Antworten

- `kategorie`: HIGH
- `quelle`: §[`DC-FA-ANCH-001.b`](../../spec/spezifikation.md#dc-fa-anch-001b--inline-html-anker)
  („zeilenbasiert und außerhalb von Fenced-Code-Blöcken **und Inline-Code-Spans**";
  `id` an einem HTML-**Element**; `name` nur an `<a>`, Tag-Name exakt) gegen die
  neue Zusage in §[`DC-FA-VER-001.a`](../../spec/spezifikation.md#dc-fa-ver-001a--versions-pin-konsistenz-versions)
  Schritt 1 und §[`DC-FA-PIN-001.a`](../../spec/spezifikation.md#dc-fa-pin-001a--content-pin-gegen-inhaltlichen-drift-pins)
  Schritt 2 („Beide Anker-Formen werden erkannt **wie in** `DC-FA-ANCH-001.b`");
  [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 1 („Die Reparatur ist die **Übernahme der vorhandenen Antwort**");
  Reviewer-Skill §HIGH (Korrektheitsfehler in Kern-Modulen mit falschen
  Befunden/Exit-Codes) und §MEDIUM (Konsistenz-Lücke zwischen Modulen derselben
  Eingabe-Klasse)
- `pfad`: `internal/hexagon/core/rules/versions.go:141`–`162` (`htmlAnchorSection`)
  gegen `internal/hexagon/core/rules/anchors.go:106`–`135` (`htmlTagRE`,
  `htmlAttrIDRE`, `htmlAttrNameRE`, `htmlAnchors`)
- `befund`: Der Commit übernimmt aus der geteilten Antwort nur die
  **Zeilen-Menge** (Prosa-Zeilen statt aller Zeilen) und behält die eigene
  Erkennung auf der **rohen** Zeile: ein einzelner Regex über `id` oder `name`
  mit Anführungszeichen. Damit bleiben vier Unterschiede zu `htmlAnchors` stehen,
  das auf den vorverarbeiteten Zeilen arbeitet und einen Tag-Kontext verlangt —
  ein Anker in einem **Inline-Code-Span**, ein `data-id`-Attribut (keine
  Wortgrenze vor `id`), ein `name` an einem beliebigen Element und eine
  anker-förmige Zeichenfolge **ohne jedes Tag** gelten für
  `versions`/`pins` als Anker und für `anchors` nicht. Die Zusage im Vertrag ist
  aber nicht „außerhalb von Fences", sondern „erkannt wie in `DC-FA-ANCH-001.b`" —
  sie behauptet die Parität, die der Code nicht liefert, und schließt damit die
  Suche.
- `verifizierbar`: ja — vier Fixtures (Scratch außerhalb des Repos), Profil
  `versions` + `anchors` + `links` in **einem** Lauf, `current-from` auf
  `version.md#aktuell`, Ziel-Datei ohne passenden Heading-Slug: (a) Anker in
  einem Inline-Code-Span, (b) `data-id`, (c) reine Prosa ohne Tag, (d) `name` an
  `<area>`. In allen vier Fällen liefert derselbe Lauf **zwei Befunde**:
  `anchor-missing` für genau diesen Anker **und** `version-stale`, weil
  `versions.current-from` ihn auflöst und die „aktuelle Version" aus dem
  Beispieltext zieht (Exit 1). Gegenprobe mit demselben Anker als echtes,
  uncodiertes `<a id=…>`: 0 Befunde, Exit 0. **Bestandsmaß (eigene Methode,
  Fence-Automat nachgebaut, drei Repos):** die reparierte Achse trägt **1**
  Vorkommen (Fence), die stehengebliebene **40** (34 + 6 in Inline-Code) — unter
  ihnen ein konkreter Versions-Anker `<a id="v0.36.0"></a>` in Inline-Code.
  Heutige Konsumenten sind nicht betroffen (der eigene `current-from` löst über
  einen Heading-Slug auf, die zwei `dpin`-Fragmente im Bestand sind
  Doku-Beispiele) — die Latenz ist dieselbe wie in der reparierten Hälfte.
- `klasse`: „geteilte Lexik, vom Konsumenten selbst vorbereitet"

### L-2 — Die Richtungs-Zusage „findet mehr, und weniger an keiner Stelle" ist falsch; zwei Richtungen verlieren Befunde, eine davon still

- `kategorie`: HIGH
- `quelle`: [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  §Konsequenzen („finden nach dieser Änderung **mehr**"), Lastenheft-Historie
  0.58.0 („**Die Änderung findet mehr** — ein grüner Konsumentenlauf kann
  fail-closed werden"), `CHANGELOG.md` (`[Unreleased]`, beide neuen Einträge:
  „**Die Änderung findet mehr**");
  [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) (Befundsatz-Zusage);
  Reviewer-Skill §HIGH (Stilles-Grün-Pfad), §MEDIUM
- `pfad`: `docs/plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md`
  §Konsequenzen, `spec/lastenheft.md` (Historien-Zeile 0.58.0), `CHANGELOG.md`
  (`[Unreleased]`); Verhalten in
  `internal/hexagon/core/rules/citations.go:187`–`200` (`citationBlockquote`) und
  `internal/hexagon/core/rules/pins.go:128`–`148` (`spanHash`)
- `befund`: Zwei konstruierbare Fälle verlieren beim Update einen Befund, und
  **kein** konsumentensichtbarer Rand nennt diese Richtung. (a) Ein
  `>`-Blockquote, das von einem Fenced-Block unterbrochen wird, wurde bis
  `v0.52.0` als **ein** Zitat gelesen und meldete `citation-mismatch` (Exit 1);
  HEAD bricht den Block am Fence, das gekürzte Zitat passt, und der Lauf endet
  mit 0 Befunden. (b) Ein `dpin`-Link, dessen Ziel-Anker ausschließlich in einem
  Fenced-Block steht, meldete bis `v0.52.0` `link-stale` (Exit 1); HEAD löst den
  Anker nicht mehr auf, und `pins` schweigt zum unauflösbaren Ziel — der
  Drift-Schutz fällt **kommentarlos** weg (Exit 0, kein Ersatz-Befund; das ist
  die zugesagte Schweige-Semantik des Moduls, nicht ihr Bruch). Beide Richtungen
  sind fachlich vertretbare Folgen der vereinheitlichten Antwort — falsch ist die
  Aussage „und **weniger** an keiner Stelle", die als universelle Behauptung
  formuliert ist und nicht als Bestandsaussage. Genau diese Klasse hat der
  Vorgänger-Report als R-5 gemeldet.
- `verifizierbar`: ja — dieselben zwei Fixtures gegen das veröffentlichte Bild
  `v0.52.0` und gegen den HEAD-Build: (a) alt Exit 1 `citation-mismatch`, HEAD
  Exit 0 / 0 Befunde; (b) alt Exit 1 `link-stale`, HEAD Exit 0 / 0 Befunde. Die
  „findet mehr"-Richtung ist im selben Lauf bestätigt (Fence zwischen Direktive
  und Zitat: alt Exit 0, HEAD Exit 2).
- `klasse`: „Verhaltensänderung eines ausgelieferten Moduls ohne gedeckte Richtung"

### L-3 — Die Fence-Grenze in `citations` hat vier Wirkstellen und zwei Assertionen; zwei Einzeiler stellen den alten Zustand wieder her, ohne dass ein Test rot wird

- `kategorie`: MEDIUM
- `quelle`: [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 4 („ein geteiltes Prädikat allein genügt nicht… Jeder reparierte
  Konsument bekommt einen Test, der die geteilte Antwort **an ihm** festnagelt");
  §[`DC-FA-CITE-001.a`](../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations)
  Schritt 2; Reviewer-Skill §HIGH (Stilles-Grün-Pfad) mit Abstufung, s. u.
- `pfad`: `internal/hexagon/core/rules/citations.go:210` (`citationParagraph`) und
  `internal/hexagon/core/rules/citations.go:190` (`citationBlockquote`) gegen
  `internal/hexagon/core/rules/citations_test.go:190`–`225`
- `befund`: Die Reparatur zieht die Grenze an **vier** Stellen — der
  Eintritts-Prüfung `fencedBlockWithin`, dem Absatz-Sammler und dem
  Blockquote-Sammler. Assertiert ist nur die Eintritts-Prüfung (drei Tests, alle
  mit dem Fence **vor** dem Zitat). Wird die Fence-Grenze im Absatz-Sammler
  entfernt, bleibt `make test` grün, und eine Direktive, deren Absatz von einem
  Fence durchschnitten wird, kippt von Exit 2 (fail-closed, die Zusage) auf
  Exit 0. Wird sie im Blockquote-Sammler entfernt, bleibt `make test` ebenfalls
  grün, und derselbe Fall kippt von Exit 0 auf ein falsches `citation-mismatch`
  (Exit 1). Das ist dieselbe Bauform, die der Vorgänger-Report als R-1 gemeldet
  hat, im selben Commit, der sie schließen soll. Eingestuft als MEDIUM statt HIGH,
  weil `citations` in keinem Gate dieses Repos aktiv ist (`.d-check.yml`
  `modules:` führt es nicht) — in einem Konsumenten-Gate wäre es der stille
  Grün-Pfad selbst.
- `verifizierbar`: ja — beide Rückbauten einzeln über eine Dateikopie angewendet,
  `make test` je Exit 0; aus jedem Rückbau zusätzlich ein Image gebaut und gegen
  das jeweilige Fixture gefahren (Zahlen oben). Arbeitsbaum danach
  wiederhergestellt (`git status --short` leer).
- `klasse`: „geteiltes Prädikat ohne Assertion gegen Wieder-Divergenz"

### L-4 — Die Blockquote-Terminatoren in §`DC-FA-CITE-001.a` Schritt 2 sind eine geschlossene Aufzählung, und die Implementierung bricht an einem vierten

- `kategorie`: MEDIUM
- `quelle`: §[`DC-FA-CITE-001.a`](../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations)
  Schritt 2: „zusammenhängende `>`-Zeilen … **eine Leer- oder Nicht-`>`-Zeile
  beendet ihn"; Reviewer-Skill §MEDIUM (Spec-Treue-Lücke einer Messmethode)
- `pfad`: `spec/spezifikation.md:1190`–`1196` gegen
  `internal/hexagon/core/rules/citations.go:187`–`200`
- `befund`: Der eingeschobene Fence-Satz steht im **Inline**-Zweig („Andernfalls
  der nächste inline-Zitat-Span im selben Absatz — Absatz im Sinne von
  `DC-FA-LINK-001.a` Schritt 2 …"). Der Blockquote-Zweig darüber blieb
  unverändert und zählt seine Terminatoren abschließend auf; der Fence ist keiner
  davon. `citationBlockquote` bricht aber genau dort — und der Unterschied ist
  am Ergebnis messbar, nicht bloß intern. Damit ist eine Verhaltensänderung im
  Vertrag nicht nur unbenannt, sondern von der geschlossenen Aufzählung
  ausgeschlossen.
- `verifizierbar`: ja — Fixture mit Direktive, Leerzeile, `>`-Zeile, Fenced-Block,
  zweiter `>`-Zeile: HEAD Exit 0 (Block endet am Fence, gekürztes Zitat passt),
  Rückbau des Fence-Bruchs Exit 1 `citation-mismatch`. Der Vertrag beschreibt die
  zweite Variante.
- `klasse`: „Rand auf der alten Fassung stehengeblieben"

### L-5 — Die Historien-Zeile 0.58.0 schreibt drei Anforderungen eine Änderung zu, die in keiner ihrer Beschreibungen und Akzeptanzkriterien steht

- `kategorie`: MEDIUM
- `quelle`: [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  (Spiegel-Tabelle, Zeile „Anforderung — `spec/lastenheft.md` (Beschreibung
  **und** Akzeptanzkriterien)"); [`AGENTS.md`](../../AGENTS.md) §2 (das Lastenheft
  ist vertraglich abnahmebindend)
- `pfad`: `spec/lastenheft.md:1342`–`1380`
  ([`DC-FA-CITE-001`](../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in)),
  der Abschnitt zu
  [`DC-FA-VER-001`](../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in)
  und der zu
  [`DC-FA-PIN-001`](../../spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in)
  gegen die Historien-Zeile 0.58.0 in derselben Datei
- `befund`: Die Versions-Zeile 0.58.0 nennt vier Anforderungen; geändert wurde der
  Körper genau **einer** (die neue Grenze in `DC-FA-VCS-001`). Bei
  `DC-FA-CITE-001` steht weiter nur „unmittelbar vor dem Zitat" ohne jede
  Absatz- oder Fence-Grenze; bei `DC-FA-VER-001` steht weiter der Satz „Anders als
  die übrigen Module liest `versions` die Pins **auch innerhalb von
  Fenced-Code-Blöcken**" — genau die Stelle, an der ein Leser die neue
  Anker-Regel für einen Widerspruch halten muss.
  [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 2 erkennt dieses Missverständnis ausdrücklich und beantwortet es in
  der ADR — nicht an der Anforderung, an der es entsteht. Der Spiegel, den
  `MR-025` an erster Stelle nennt, steht damit auf der alten Fassung, während die
  Delta-Zeile das Gegenteil behauptet.
- `verifizierbar`: ja — Volltext beider Anforderungen gegen die Historien-Zeile
  0.58.0 gelesen; die Zeichenketten „Fence", „Absatz" und „Anker" kommen in den
  drei Anforderungs-Abschnitten nicht in der zugeschriebenen Bedeutung vor. Der
  Diff der Range berührt in `spec/lastenheft.md` nur die Versionszeile, den
  `DC-FA-VCS-001`-Grenzabsatz und die Historien-Zeile.
- `klasse`: „Semantik im Körper geändert, der Rand referiert die andere Fassung"

### L-6 — Die benannte `vcs`-Grenze beschreibt eine verschobene Maske, nicht das stille Grün; ihr Re-Evaluierungs-Trigger ist in genau dieser Richtung unbeobachtbar

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  §Grenze (neu); [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 3 und §Re-Evaluierungs-Trigger („Ein unbalancierter Fence in einer
  immutablen Revision **wird beobachtet**");
  §[`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
  Schritt 4; Reviewer-Skill §HIGH (Stilles-Grün-Pfad in einem Gate) mit Abstufung, s. u.
- `pfad`: `spec/lastenheft.md` (Grenz-Absatz zu `DC-FA-VCS-001`),
  `internal/hexagon/core/rules/matrix.go:293`–`317` (`excludedRanges`; ohne
  folgende Überschrift bleibt `to` gleich 0) und
  `internal/hexagon/core/rules/matrix.go:320`–`327` (`inRanges`; `to == 0`
  bedeutet **bis Dateiende**)
- `befund`: Die Grenze sagt, ein unbalancierter Fence „verschiebt" die Maske und
  „kein Wächter meldet das". Die praktisch gefährliche Ausprägung ist stärker:
  liegt der offene Fence **innerhalb** des ausgenommenen Abschnitts, findet
  `excludedRanges` keine folgende Überschrift mehr, die Ausnahme läuft bis zum
  Dateiende, und eine reale Änderung am Core einer `Accepted`-ADR passiert das
  Immutabilitäts-Gate **mit Exit 0 und ohne Ausgabe**. Der Vorgänger-Report R-7
  hat nur die Falsch-Rot-Richtung belegt; die stille Richtung steht weder in der
  Grenze noch in der ADR noch im `CHANGELOG`. Sie macht zugleich den ersten
  Re-Evaluierungs-Trigger wirkungslos: „wird beobachtet" setzt eine Beobachtung
  voraus, und in dieser Richtung gibt es nichts zu beobachten. Eingestuft als
  MEDIUM statt HIGH, weil der stille Pfad **älter** ist als dieser Diff und die
  Nicht-Lieferung mit einer ADR gedeckt ist (kein Gate-Suppression ohne ADR) —
  gemeldet wird, dass die bewusste Abnahme auf einer zu schwachen Beschreibung
  des abgenommenen Risikos beruht.
- `verifizierbar`: ja — git-Fixture außerhalb des Repos, `exclude-sections:
  [Geschichte]`, `immutable-when` auf die `Accepted`-Kopfzeile, zwei Commits, die
  einzige Änderung ist eine Zeile in einem `## Kern`-Abschnitt **nach** dem
  Geschichts-Abschnitt: mit einem unbalancierten Fence in `## Geschichte` meldet
  `--enable vcs --range HEAD~1..HEAD` **0 Befunde, Exit 0**; dieselben zwei
  Commits ohne den Fence melden `core-drift-vcs`, Exit 1.
- `klasse`: „Modul-Grenze benannt, aber nicht in ihrer stillen Richtung"

### L-7 — Die Aufzählung, die den Klassen-Abschluss belegen soll, nennt zwei von vier Dateien mit eigenem Fence-Automaten

- `kategorie`: LOW
- `quelle`: [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  §Re-Evaluierungs-Trigger („Dann genügt die **Aufzählung** als Beleg nicht
  mehr"); Reviewer-Skill §LOW (latente Wartungsfalle)
- `pfad`: `internal/hexagon/core/rules/sections.go:19`–`37` und
  `internal/hexagon/core/rules/sections.go:42`–`57` (zwei eigene
  Fence-Zustands-Schleifen), `internal/hexagon/core/rules/spans.go:32`–`92`
- `befund`: Der Beleg-Text zählt die Stellen der Fence-Erkennung als „nur
  `markdown.go` … und `trace_table.go`" auf. Tatsächlich führen **vier** Dateien
  eine eigene Fence-Zustands-Schleife: neben diesen beiden auch `sections.go`
  (zweimal, für Abschnitts-Kopf und Abschnitts-Ende, konsumiert von `planning`
  und `structure`) und `spans.go` (der Wächter). Alle vier speisen sich aus
  denselben geteilten Prädikaten, die **Antwort** stimmt also überein — die
  Aufzählung ist trotzdem der erklärte Beleg für „die Klasse ist geschlossen",
  und ein Beleg, der zwei von vier Kandidaten nicht nennt, trägt die Aussage
  nicht. Das Failure-Szenario ist die nächste Änderung an der Fence-Lexik: wer
  die Konsumenten aus dieser Liste ableitet, findet `sections.go` nicht.
- `verifizierbar`: ja — Vollzählung der Aufrufstellen von `FenceToggle`,
  `FenceRun`, `FenceCloses` und `TrimFenceIndent` im Nicht-Test-Code des
  `internal`-Baums.
- `klasse`: „Aufzählung als Beleg, aus dem Gedächtnis abgeleitet"

### L-8 — Die Nutzer-Doku beschreibt weiter die Paarung, die es nicht mehr gibt

- `kategorie`: LOW
- `quelle`: [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  (Spiegel „Nutzer-Doku"); Reviewer-Skill §LOW (Doku-Drift)
- `pfad`: `docs/user/benutzerhandbuch.md:1139`–`1148` (§5-Aufgabe),
  `docs/user/benutzerhandbuch.md:1659` (§6-Modul-Tabelle),
  `README.de.md:104` und `README.md:103`
- `befund`: Alle vier Stellen sagen weiterhin, die Direktive markiere „das
  folgende Zitat — die nächste nicht-leere Zeile als `>`-Blockquote oder den
  nächsten inline-Zitat-Span im selben Absatz", ohne die Fence-Grenze. Genau diese
  Stellen liest ein Konsument, dessen bisher grüne Direktive nach dem Update mit
  Exit 2 abbricht. Die Release-Prep-Checkliste in `docs/user/releasing.md` §4 deckt
  den Fall nur halb: sie nennt die §11-Verlaufszeile und „ggf. **neue**
  Feature-Abschnitte (§5/§6)" — eine **geänderte** Zusage in einem bestehenden
  Abschnitt steht dort nicht. Bewusst LOW: das Repo zieht Handbuch und README
  planmäßig erst im Release-Prep nach, und der Stand ist unreleased
  (`CHANGELOG.md` `[Unreleased]`).
- `verifizierbar`: ja — die vier Stellen gegen den HEAD-Lauf des Fixtures
  „Fence zwischen Direktive und Zitat" (Exit 2, „d-check:cite ohne folgendes
  Zitat").
- `klasse`: „Rand auf der alten Fassung stehengeblieben"

### L-9 — Der `--doctor`-Klartext von `anchor-missing` nennt nur die Heading-Slug-Hälfte, während dieselbe Meldung im Befund beide Hälften nennt

- `kategorie`: LOW
- `quelle`: [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  (Spiegel „Klartexte"); §[`DC-FA-ANCH-001.b`](../../spec/spezifikation.md#dc-fa-anch-001b--inline-html-anker)
- `pfad`: `internal/hexagon/core/app/diagnose.go:104` gegen
  `internal/hexagon/core/rules/anchors.go:225`
- `befund`: Der Klartext lautet „Anker entspricht keinem Heading-Slug"; die
  Befund-Meldung desselben Codes lautet „Anker entspricht keinem Heading-Slug
  **und keinem HTML-Anker** der Zieldatei". Der Nutzer, der über `--doctor`
  nachschlägt, warum ein Anker fehlt, erfährt die HTML-Hälfte nicht — also genau
  die Achse, die dieser Diff bei zwei weiteren Modulen vereinheitlicht.
  **Älter als dieser Diff** und außerhalb der Range; genannt, weil `MR-025` diesen
  Spiegel für die hier geänderte Semantik ausdrücklich führt.
- `verifizierbar`: ja — `--doctor anchor-missing` gegen die Befund-Meldung eines
  Laufs mit einem HTML-Anker als einzigem Treffer.
- `klasse`: „Rand auf der alten Fassung stehengeblieben"

### L-10 — Der eingeschobene Fence-Satz verwaist das Pronomen im normativen Schritt

- `kategorie`: LOW
- `quelle`: §[`DC-FA-CITE-001.a`](../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations)
  Schritt 2; [`AGENTS.md`](../../AGENTS.md) §2 (die Spezifikation ist technisch
  verbindlich)
- `pfad`: `spec/spezifikation.md:1193`–`1199`
- `befund`: Der Schritt liest nach dem Einschub: „… es folgt dann kein Zitat, und
  der Fall unten greift — **er** darf Prosa vor sich haben und über mehrere Zeilen
  laufen, begrenzt durch ein Anführungs-Paar". Bezugswort von „er" war der
  inline-Zitat-Span; nächstes Substantiv ist jetzt „der Fall unten". Wer den
  Schritt implementiert, muss raten, welche der beiden Größen ein
  Anführungs-Paar begrenzt.
- `verifizierbar`: ja — Volltext des Schritts gegen die Fassung vor der Range
  (`git diff` auf `spec/spezifikation.md`).
- `klasse`: „Einschub bricht einen normativen Satz"

### L-11 — In derselben Registerzelle zeigt der Link nach `in-progress/`, die Prosa sagt `open/`

- `kategorie`: LOW
- `quelle`: `docs/plan/planning/observations.md` (Register-Zelle BEO-003);
  Reviewer-Skill §LOW (Doku-Drift)
- `pfad`: `docs/plan/planning/observations.md` (Zeile der Beobachtung BEO-003)
- `befund`: Der Commit zieht in der BEO-003-Zelle das Link-Ziel von `open/` nach
  `in-progress/` nach, lässt aber den Satzteil „dafür liegt … in `open/`" stehen.
  Wer das Register liest statt dem Link zu folgen, sucht am falschen Ort. Die
  Zelle ist dabei die einzige Änderung an `observations.md` in der ganzen Range —
  der Zähler von BEO-003 steht unverändert auf 2, obwohl das Wellendokument die
  Entscheidung „geschlossen oder auf 3" ausdrücklich zum Closure-Trigger macht
  (das ist Wellen-Closure, nicht Slice-DoD, und darum hier nur als Beobachtung).
- `verifizierbar`: ja — `git diff` der Range auf `docs/plan/planning/observations.md`.
- `klasse`: „Rand auf der alten Fassung stehengeblieben"

### L-12 — Die Grundmenge der Fall-1-Messung besteht aus Marker-Erwähnungen, nicht aus produktiven Direktiven

- `kategorie`: INFO
- `quelle`: die Messtabelle des geprüften Slice-Plans (Fall 1, Grundmenge
  „18 `d-check:cite`-Direktiven"); Reviewer-Skill §INFO
- `pfad`: keine Code-Stelle — Messmethode
- `befund`: Von den 18 Marker-Vorkommen im Bestand ist **keines** eine produktive
  Direktive: alle 18 stehen in Inline-Code oder Prosa (Handbuch, READMEs,
  `CHANGELOG`, Spec, ADR, Review-Reports), vier davon in voll geformter
  Beispielform. Ein Lauf mit aktivem `citations` über das eigene Repo bricht
  darum schon in Schritt 1 fail-closed ab — identisch mit dem Bild `v0.52.0`. Das
  **Ergebnis** der Messung (null betroffene Fälle) bleibt richtig; die Zahl 18
  suggeriert eine geprüfte Grundmenge, die es nicht gibt.
- `verifizierbar`: ja — Auflistung aller 18 Fundstellen; Lauf mit `citations`
  über die drei Repos gegen HEAD und gegen `v0.52.0` liefert dieselbe
  Schritt-1-Meldung an derselben Stelle.
- `klasse`: „Grundmenge zählt Erwähnungen statt Fälle"

### L-13 — Ein äquivalenter Mutant in der neuen Grenze

- `kategorie`: INFO
- `quelle`: [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 4 (Mutations-Echtheit); Reviewer-Skill §INFO (dokumentationswürdige
  Annahme)
- `pfad`: `internal/hexagon/core/rules/citations.go:190` und
  `internal/hexagon/core/rules/citations.go:210` (die Wache `k > j`)
- `befund`: Die Wache `k > j` in beiden Sammlern ist verhaltensneutral: die
  Eintritts-Prüfung hat den Schritt bis `j` bereits ausgewertet, und `j` ist immer
  mindestens 1, also gibt es weder einen zweiten Treffer noch einen negativen
  Index. Ihr Entfernen ist ein echter äquivalenter Mutant — kein Test kann ihn
  fangen, und er ist harmlos. Festgehalten, damit ein späterer Mutationslauf ihn
  nicht als Testlücke fehlliest.
- `verifizierbar`: ja — Wache in beiden Sammlern über eine Dateikopie entfernt,
  `make test` Exit 0; Arbeitsbaum wiederhergestellt.
- `klasse`: „äquivalenter Mutant"

## Negativbefunde

- geprüft, ohne Befund: **Das Absatz-Prädikat ist wirklich geteilt, nicht nur
  gleich benannt.** Der Rückbau von `fencedBlockBetween` auf eine konstante
  Antwort macht **drei** Tests rot — die zwei neuen `citations`-Assertionen
  **und** den bestehenden Absatz-Test der Vorverarbeitung. Die entsprechende
  Zusage der Definition of Done hält.
- geprüft, ohne Befund: **Die fünf zugesagten Assertionen fangen ihren Rückbau.**
  Eigener Nachvollzug über Dateikopien: `fencedBlockWithin` konstant falsch ⇒ zwei
  `citations`-Tests rot; `fencedBlockWithin` auf die naive Ein-Schritt-Form ⇒ der
  Gegenproben-Test „ohne Fence paart normal" rot (die naive Form würde jede
  Direktive mit mehr als einer Leerzeile Abstand fälschlich trennen); die
  Prosa-Wache im HTML-Anker-Zweig entfernt ⇒ der `versions`- **und** der
  `pins`-Test rot. Fünf von fünf.
- geprüft, ohne Befund: **Der Fence-Randfall-Fächer der Reparatur.** Fence
  unmittelbar nach der Direktive ohne Leerzeile, `~~~` statt Backticks,
  eingerückter Fence, zwei aufeinanderfolgende Fences, unbalancierter Fence und
  CRLF-Zeilenenden — alle sechs liefern den zugesagten fail-closed Exit 2 mit
  derselben Meldung; die Gegenprobe ohne Fence bleibt bei 0 Befunden.
- geprüft, ohne Befund: **Der Heading-/Slug-Zweig der Anker-Auflösung.** Er lief
  schon vorher über den fence-bewussten Prosa-Automaten und ist unverändert; die
  Umstellung betrifft ausschließlich den HTML-Zweig. Der Slug-Zweig hat Vorrang,
  ein Heading-Anker maskiert also einen gleichnamigen HTML-Anker weiterhin.
- geprüft, ohne Befund: **Die drei gescopten Roh-Lesungen sind unangetastet.** Der
  Pin-Scan von `versions` läuft weiter über alle Rohzeilen einschließlich Fences
  ([ADR-0019](../plan/adr/0019-versions-pin-fence-ausnahme.md)), der gehashte
  Ziel-Span von `pins` bleibt roh
  ([ADR-0020](../plan/adr/0020-content-pin-fence-ausnahme.md)), und `immutable`
  normalisiert weiter den rohen Core. Der Diff berührt in beiden Modulen nur die
  Anker-**Erkennung**, nicht den Span.
- geprüft, ohne Befund: **Fall 3 im Bestand, mit weiterer Grundmenge als die
  Messung.** Alle Blobs, die in irgendeiner erreichbaren Revision unter dem
  ADR-Verzeichnis existiert haben, dedupliziert nach Blob-Hash extrahiert
  (**248** statt der gemessenen 152) und mit aktivem `spans` geprüft: **0
  Befunde**. Die Messaussage „null unbalancierte Fences in immutablen Revisionen"
  hält auf einer echten Obermenge.
- geprüft, ohne Befund: **Fall 2 im Bestand.** Eigene Zählung mit nachgebautem
  Fence-Automaten über die Markdown-Dateien der drei Repos: genau **1**
  anker-förmiges Vorkommen innerhalb eines Fence — dieselbe Zahl wie die Messung.
- geprüft, ohne Befund: **Keine neuen Falsch-Positiven durch die
  `citations`-Änderung.** Läufe mit aktivem `citations` über die drei Repos gegen
  HEAD und gegen das veröffentlichte Bild `v0.52.0` liefern an jeder Stelle
  dasselbe Ergebnis; die einzigen Abbrüche sind die schon vorher bestehenden
  Schritt-1-Meldungen auf Prosa-Erwähnungen der Direktive.
- geprüft, ohne Befund: **Kein dritter Anker-Automat.** Im Produkt gibt es genau
  zwei Stellen, die „ist das ein Anker" beantworten — die Anker-Menge in
  `anchors.go` (geteilt von `anchors` und `codepaths`) und `htmlAnchorSection` in
  `versions.go` (geteilt von `versions` und `pins`). Eine dritte existiert nicht;
  die verbleibende Differenz zwischen den beiden ist L-1.
- geprüft, ohne Befund: **Heading- und Abschnitts-Lexik.** Die
  Überschriften-Erkennung läuft überall über `parseATXHeading` auf Prosa-Zeilen;
  auch die zweite Schleifenform in `sections.go` liefert dieselbe Antwort (sie
  ruft dieselben Prädikate und startet den Fence-Zustand an einer Überschrift,
  die per Konstruktion außerhalb eines Fence liegt). Ihre bloße Nicht-Nennung in
  der Beleg-Aufzählung ist L-7, keine Divergenz.
- geprüft, ohne Befund: **Referenz-Richtung (SDP) und Marker-Ehrlichkeit.** Die
  neue ADR nennt eine Slice-Kennung ausschließlich in ihrem
  `## Geschichte`-Abschnitt, den die Referenzmatrix-Konfiguration ausnimmt; einen
  Provenance-Marker trägt der Diff nicht, also gibt es keine Deklaration, deren
  Ehrlichkeit zu prüfen wäre. Das Matrix-Modul ist im grünen Gate-Lauf aktiv.
- geprüft, ohne Befund: **Hard Rules.** Kein Spec-Stratum nennt eine ADR — der
  neue Grenz-Absatz schreibt „Begründung in begleitender ADR" statt einer Kennung
  ([`AGENTS.md`](../../AGENTS.md) §3.4). Die neue ADR steht auf `Proposed`, ist im
  Index eingetragen und trägt den geforderten Re-Evaluierungs-Trigger-Abschnitt
  (§3.5, §5). Keine Inline-Suppression, keine Gate-Lockerung, kein Netzzugriff
  außerhalb der Netz-Module.
- geprüft, ohne Befund: **Keine neue Vertragsfläche unterschlagen, wo keine
  entsteht.** Der Diff führt keinen neuen Grund-Code, keinen neuen
  Config-Schlüssel und kein neues Modul ein; §2-Schema, §4-Grund-Code-Tabelle,
  die Klartext-Liste und die `--print-config`-Vorlage bleiben zu Recht
  unverändert. Die tatsächlich betroffenen Spiegel sind L-5 (Anforderung), L-8
  (Nutzer-Doku) und L-9 (Klartext der berührten Anker-Semantik).
- geprüft, ohne Befund: **Gate-Stand.** `make gates` Exit 0 (372 Dateien, 0
  Befunde; `doc-check`, `lint`, `test`, `arch-check`, `coverage-gate`, `semgrep`,
  `gate-consistency`, `planning-check`). Der Neubau nach allen
  Mutations-Gegenproben liefert dieselbe Image-ID wie vor dem ersten Eingriff.
- geprüft, ohne Befund: **Arbeitsbaum.** Nach jeder Mutation aus der
  Scratch-Kopie zurückgeschrieben; `git status --short` und `git diff --stat` sind
  am Ende leer (der Report selbst ist die einzige neue Datei).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 4 |
| LOW | 5 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** geteilte Lexik, vom Konsumenten selbst
vorbereitet · Verhaltensänderung eines ausgelieferten Moduls ohne gedeckte
Richtung · geteiltes Prädikat ohne Assertion gegen Wieder-Divergenz · Rand auf der
alten Fassung stehengeblieben · Semantik im Körper geändert, der Rand referiert die
andere Fassung · Modul-Grenze benannt, aber nicht in ihrer stillen Richtung ·
Aufzählung als Beleg, aus dem Gedächtnis abgeleitet · Einschub bricht einen
normativen Satz · Grundmenge zählt Erwähnungen statt Fälle · äquivalenter Mutant.

**Wiederholungs-Signal.** Drei der zehn Klassen sind Wiedergänger aus der
unmittelbar vorangegangenen Review-Runde desselben Vorhabens: „geteiltes Prädikat
ohne Assertion" (dort R-1, hier L-3), „Richtung der Verhaltensänderung nicht
gedeckt" (dort R-5, hier L-2) und „Rand auf der alten Fassung" (dort R-4/R-6, hier
L-4/L-8/L-9/L-11). Die vierte Stelle, die eine Lexik-Frage selbst beantwortet
(L-1), ist genau der Fall, den der Re-Evaluierungs-Trigger von
[ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) als Ende
der Aufzählungs-Beweisführung benennt.

## Verdikt

**Merge-blockierend:** ja — zwei HIGH und vier MEDIUM.

**Ist die Klasse geschlossen? Nein — aber sie ist kleiner geworden, und die
Messung trägt.** Was geliefert wurde, ist gut belegt: das Absatz-Prädikat ist
wirklich geteilt (der Rückbau macht auch den *bestehenden* Absatz-Test rot, und
das ist der stärkste Beleg im ganzen Slice), die fünf zugesagten Assertionen sind
unabhängig nachvollzogen und alle fünf rot, der Fence-Randfall-Fächer hält über
sieben konstruierte Varianten einschließlich CRLF und `~~~`, und die
Bestandszahlen halten einer eigenen Methode stand — Fall 2 auf denselben Wert,
Fall 3 sogar auf einer um zwei Drittel größeren Grundmenge. Der Verzicht auf Code
für Fall 3 ist als Entscheidung tragfähig: es gibt dort keine vorhandene richtige
Antwort zu übernehmen, und der Anlassfall existiert im Bestand nicht.

Offen bleibt die Klasse an drei Stellen, und alle drei liegen **innerhalb** dessen,
was der Diff zu schließen behauptet. Die Anker-Frage ist nur zur Fence-Hälfte
vereinheitlicht; die drei anderen Unterschiede zwischen den beiden Antworten sind
geblieben, im Bestand vierzigmal statt einmal vertreten, und der Vertrag sagt
seit diesem Commit die volle Parität zu — ein Lauf liefert weiter `anchor-missing`
und `version-stale` für denselben Anker (L-1). Die Fence-Grenze in `citations`
wirkt an vier Stellen und ist an zweien nicht assertiert; jede der beiden
Rückbau-Zeilen ist grün in `make test` und rot in der Wirklichkeit (L-3). Und die
Richtungs-Zusage, die auf jeder release-sichtbaren Fläche steht, ist gegen das
veröffentlichte Vor-Bild in zwei Fällen widerlegt, einer davon ein
kommentarloser Verlust eines Drift-Checks (L-2).

Die Vertragsflächen sind ungleich gepflegt: Algorithmus und ADR sind sorgfältig,
die **Anforderung** — der erste Spiegel, den `MR-025` nennt — steht für drei der
vier genannten Kennungen unverändert da, während die Delta-Zeile ihre Änderung
behauptet (L-5). Die neue `vcs`-Grenze ist die richtige Lieferung an der richtigen
Stelle, beschreibt aber die harmlosere ihrer beiden Richtungen und hängt ihre
Re-Evaluierung an eine Beobachtung, die es in der stillen Richtung nicht geben
kann (L-6).

**Release-Empfehlung: noch nicht releasen.** Die Einordnung als **Minor** ist
richtig, ihre Begründung nicht: „findet mehr, und weniger an keiner Stelle" ist
falsifiziert, und ein Konsument verliert nach dem Update unter Umständen einen
`pins`-Drift-Befund, ohne dass irgendetwas es ihm sagt. L-1, L-2 und L-3 gehören
vor den Tag in diesen Slice; L-4, L-5 und L-6 sind Vertragstext und liegen auf
demselben Diff. L-8 und L-9 sind Release-Prep-Material, L-7, L-10 und L-11 sind
klein und billig.

**Übergabe:** Die Findings gehen an den Implementer; die Finding-Klassen gehören
in die Slice-Closure und von dort in das Beobachtungs-Register — mit der
ausdrücklichen Notiz, dass L-1 eine **vierte** Stelle derselben Lexik-Klasse ist
und damit den ersten Re-Evaluierungs-Trigger von
[ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) auslöst,
noch bevor die ADR `Accepted` wird. Dieser Report ist ein Lauf-Beleg (dieser Diff,
dieser Skill, dieses Modell, dieses Verdikt) und ersetzt keine Verifikation —
DoD- und Plan-Konformität prüft der Verifier separat.
