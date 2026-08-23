# Review-Report: slice-128 — Etappe A, Bundle `v5.11.0` vendored und Pin gehoben

**Datum:** 2026-08-23 · **Review-Art:** Doku-/Konventions-Review (geprüft gegen Slice-Plan slice-128 §1/§2/§3/§4/§5/§7, Wellendokument welle-83 §1/§3/§6, `MR-011`/`MR-013`/`MR-021`/`MR-023`/`MR-030`, Hard Rules AGENTS §3.3/§3.7/§4, vendorte Vorlage `MR-NNN-titel.template.md`, Beobachtungs-Register `BEO-002`/`BEO-005`/`BEO-008`/`BEO-009`/`BEO-011`), unabhängiger Reviewer ohne Anteil an der Arbeit
**Gegenstand:** `git diff 1f534ea..HEAD` — vier Commits: `cc8fffa` (Lifecycle-Move), `f10f471` (Bundle vendored), `5331466` (Pin gehoben, MR-030), `96cbd6a` (Alt-Baum entfernt)
**Skill:** `.harness/skills/reviewer.md` @ 1.9.0 (`5331466`) · **Modell-ID:** `claude-opus-5[1m]`

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-128-baseline-v5110-vendoring.md`
- Wellendokument `docs/plan/planning/welle-83-baseline-v5110-migration.md`
- `harness/conventions.md` §Baseline, §Adoptierte Konventions-Quellen, §Adaptions-Block (beide Index-Tabellen)
- `harness/conventions/MR-030-baseline-v5110.md`, `harness/conventions/done/MR-029-baseline-v590.md`, `MR-011`, `MR-013`, `MR-021`, `MR-023`
- vendorte Vorlage `.harness/baseline/v5.11.0/templates/harness/conventions/MR-NNN-titel.template.md`
- `docs/plan/planning/observations.md` (`BEO-008` zentral; `BEO-002`, `BEO-005`, `BEO-009`, `BEO-011`)
- `tools/harness/fetch-baseline-cache.sh` (Modi `--verify` / `--check-latest`)
- `.d-check.yml` (`scan.ignore`, `ignore-refs`, `codepaths.roots`/`exempt-paths`, `versions`, `structure`)
- `AGENTS.md` §3.3 (Move-Commit-Grenzen), §3.7, §4
- Nicht erhalten: die DoD-Abhakung (Verifikations-Rolle, getrennter Kontext)

**Vom Reviewer selbst gefahren** (Exit je Lauf in eine Datei umgeleitet und direkt gelesen, `BEO-007`):
`make doc-check` Exit 0 (457/0) · `make gates` Exit 0 (acht Gates, Coverage 94,80 % ggü. 93 %) · `bash tools/harness/fetch-baseline-cache.sh --verify` Exit 0 · derselbe Aufruf mit `v5.9.0` Exit 1 · `bash tools/harness/fetch-baseline-cache.sh --check-latest` Exit 0 (Netz, informativ) · neun eigene Gate-Gegenproben mit einer temporären Sonde unter `docs/` (danach entfernt, Arbeitsbaum wieder sauber) · Datei-für-Datei-`diff` des aus `f10f471` rekonstruierten Alt-Baums gegen den neuen.

**Verdikt: blockierend** — kein HIGH, vier MEDIUM, drei LOW, ein INFO.

---

## Findings

**F-1**

- **kategorie:** MEDIUM
- **quelle:** Slice-Plan §2 Schritt 3 / §4 DoD („Verweis-Hebung über **alle drei** Klassen belegt") · `MR-030` §Geltungsbereich (nennt `harness/README.md`) · `BEO-008` Richtung 1, Klasse 2
- **pfad:** `harness/README.md:60`
- **befund:** Die `lab-regelwerk.zip`-Download-URL in der Tabellenzeile steht weiter auf `releases/download/v5.9.0/lab-regelwerk.zip`, während das Link-Ziel **derselben Zeile** im selben Commit auf `.harness/baseline/v5.11.0/regelwerk/` gehoben wurde. Die Zelle sagt im Präsens, das adoptierte Betriebsregelwerk sei „vendored aus dem self-contained `lab-regelwerk.zip`" — der vendorte Baum stammt aber aus dem `v5.11.0`-Asset, was `--check-latest` (Content-Teil) unabhängig bestätigt. Eigene Zählung: auf der `+`-Seite des Diffs stehen genau vier gehobene Tag-URLs (`AGENTS.md` 1, `harness/conventions.md` 3), auf der `-`-Seite fünf lebende Vorkommen — die Botschaft von `5331466` („Klasse 2 (Release-/Tree-URLs): 4 gehoben") beschreibt die Menge, die gehoben wurde, nicht die Menge, die zu heben war. Es ist wortwörtlich die Form, die `BEO-008` benennt („teils in sich widersprüchlich in derselben Zeile — Link-Ziel neu, Quell-URL alt"), und die dritte Wiederholung an **derselben Datei:Zeile**: `docs/reviews/2026-08-21-slice-110-baseline-v570-review.md:12-19` führt sie als F-1, Klasse `pin-spiegel-nicht-gehoben`, und verweist ihrerseits auf slice-106 F-1/F-3.
- **verifizierbar:** nein — kein Gate deckt Versionen in externen URLs (`external` ist im `doc-check`-Profil nicht aktiv, `versions.pin-pattern` matcht nur `ghcr.io`-Pins). Belegt durch die eigene `+`/`-`-Zählung der Tag-URLs über den Diff und den Volltext der Zeile.
- **klasse:** pin-spiegel-nicht-gehoben

**F-2**

- **kategorie:** MEDIUM
- **quelle:** Slice-Plan §2 Schritt 3 (Klasse 3: Prosa-/„Stand"-Pins) · `BEO-008` Richtung 1
- **pfad:** `docs/plan/planning/in-progress/slice-127-claude-md-pointer.md:155-159`
- **befund:** Der Start-Trigger des lebenden `next/`-Slice sagt im Indikativ über die Gegenwart: „**Neuer Wartegrund:** dieser Kanon ist bei uns noch nicht **gepinnt** — wir stehen auf `v5.9.0` (Kurs-Welle 86)" und begründet die Blockade damit, sonst führte `AGENTS.md` Stand 94, „während der Konventionsspeicher `v5.9.0` pinnt". Nach `5331466` pinnt der Konventionsspeicher `v5.11.0`; beide Aussagen sind falsch. Derselbe Commit hat in **derselben Datei** die zwei Bezug-Links (Z. 11 und Z. 13) auf `v5.11.0` gehoben und die Prosa stehen gelassen — wieder die Form „Link-Ziel neu, Prosa-Pin alt". Die Botschaft von `5331466` nennt als einzige lebende Klasse-3-Stelle den Ellipsen-Pin in `MR-021` und führt diese Stelle weder unter den gehobenen noch unter den bewusst stehen gelassenen.
- **verifizierbar:** nein — kein Gate liest Zeitform oder Prosa-Pins. Belegt durch `grep` über alle `v5.9.0`-Nennungen außerhalb `docs/reviews/**`, `docs/plan/planning/done/**` und `harness/conventions/done/**`, Stelle für Stelle gegen ihre Zeitform gehalten.
- **klasse:** prosa-pin-nicht-gehoben

**F-3**

- **kategorie:** MEDIUM
- **quelle:** `BEO-009` (Botschaft behauptet mehr, als die Arbeit trägt) · `DC-FA-CODE-001` (Modul `codepaths`)
- **pfad:** Commit-Botschaft `96cbd6a`, Absatz 2–3
- **befund:** Die Botschaft begründet den Verzicht auf einen `ignore-refs`-Eintrag damit, dass „`docs/reviews/**` vom `codepaths`-Ventil ausgenommen" sei „und die verbliebene Zeile in `slice-117` ein VERZEICHNIS nennt, das die Pfad-Prüfung nicht auflöst". Beide Halbsätze sind am Modul widerlegt. Gegenprobe: `docs/gibtesnicht/` in Inline-Code in einer nicht ausgenommenen Datei liefert `codepath-missing`, Exit 2 — Verzeichnisse werden sehr wohl aufgelöst. Und `.harness/baseline/v5.9.0/regelwerk/modul-06-roadmap.md` (eine **Datei**, nicht ein Verzeichnis) in derselben Position liefert Exit 0, 0 Befunde; ebenso `.harness/skills/gibt-es-nicht.md`, das in keinem Ventil steht. Der tatsächliche Grund steht in `internal/hexagon/core/rules/codepaths.go:192` (`classifyCodepath`): als Pfad erkannt wird nur, was mit `./`, `../` oder einer der `codepaths.roots` (`docs|spec|tools|harness|internal|cmd`) beginnt — `.harness/…` fällt durch alle drei. Dieselbe Zeile in relativer Schreibweise (`../.harness/baseline/v5.9.0/…`) wird dagegen gefangen (Exit 2). Damit trägt „`make doc-check` nach dem Entfernen Exit 0 (457 Dateien, 0 Befunde)" die gezogene Schlussfolgerung nicht: der grüne Lauf ist für die gesamte Ziel-Klasse `.harness/…`-in-Inline-Code blind, unabhängig von Ventilen und unabhängig von Datei/Verzeichnis. Die Unterscheidung, die trägt, ist die andere im selben Absatz genannte (Markdown-Link gegen Inline-Code) — ein Link auf den entfernten Baum liefert `target-missing`, Exit 2.
- **verifizierbar:** ja — Gegenproben B/D/E/F/I mit einer temporären Sonde unter `docs/`, jeder Exit direkt gelesen; zusätzlich der Quelltext von `classifyCodepath`.
- **klasse:** gruenes-gate-als-beleg-fuer-blinde-achse

**F-4**

- **kategorie:** MEDIUM
- **quelle:** `BEO-009` Richtung (a) (behauptete Prüfung fand so nicht statt) · Slice-Plan §2 Schritt 3
- **pfad:** Commit-Botschaft `5331466`, Absatz „ZWEI STELLEN BEWUSST OFFEN GELASSEN" · `docs/plan/planning/in-progress/roadmap.md:13` · `.d-check.yml:332`
- **befund:** Die Begründung lautet: die beiden Stellen zitierten eine Regel-**Fassung**, sie zu heben hieße zu behaupten, die Regel sei in `v5.11.0` unverändert, „das ist **nicht geprüft** und gehört ins Delta-Audit". Geprüft ist es — im selben Commit. Beide Stellen zitieren `modul-06-roadmap.md` §Roadmap-Struktur bzw. dieselbe Offene-Wellen-Form; `modul-06-roadmap.md` steht in der Delta-Messung derselben MR in der Gruppe der 23 Dateien, die „ausschließlich den Versions-Stempel" ändern. Eigene Gegenprobe: die Sektion §Roadmap-Struktur ist zwischen altem und neuem Baum byte-identisch, `modul-06-roadmap.md` und `modul-05-planning-harness.md` sind es als Ganzes (ohne die Quell-Stempel-Zeile). Die Unterscheidung Pfad-Pin ↔ Regel-Fassung ist als solche tragfähig; die Aussage, der Befund darüber liege nicht vor, ist es nicht. Folge nach `96cbd6a`: zwei lebende Artefakte — die Roadmap in `in-progress/` und die Prüf-Config — zitieren eine Fassung `v5.9.0`, deren Baum im selben Slice entfernt wurde und die im netzlosen Lesepfad nicht mehr existiert.
- **verifizierbar:** nein — kein Gate liest Fassungs-Angaben in Prosa bzw. Config-Kommentaren. Belegt durch Sektions-`diff` und Ganzdatei-`diff` alt↔neu.
- **klasse:** begruendung-von-eigener-messung-widerlegt

**F-5**

- **kategorie:** LOW
- **quelle:** `MR-029` §Adaption (Präzedenz: „Von diesen elf trägt eine … im Rumpf ebenfalls nur einen Versions-Zeiger — **zehn Dateien tragen also echten Regel-Inhalt**")
- **pfad:** `harness/conventions/MR-030-baseline-v5110.md:26-38` · gleichlautend Commit-Botschaft `5331466`
- **befund:** „**sechs** tragen echten Rumpf-Inhalt". Gemessen tragen fünf Rumpf-Inhalt. `regelwerk/README.md` (+3/−3) besteht aus zwei `blob/v5.x`-Quell-URL-Stempeln und der `**Stand:**`-Zeile (`Kurs-Welle 86 · 2026-08-22` → `Kurs-Welle 94 · 2026-08-23`) — kein Regel-Inhalt. Die MR-eigene Umfangs-Tabelle annotiert die Zeile mit „(Stand-Zeile)", die Zahl im Fließtext und in der Commit-Botschaft nimmt die Annotation nicht auf; der unmittelbare Vorgänger hat genau diese Subtraktion ausdrücklich vorgenommen und die bereinigte Zahl genannt. Die übrigen fünf Umfangs-Angaben und alle drei Bucket-Zahlen stimmen exakt (siehe Negativbefunde 1 und 2).
- **verifizierbar:** nein — Prosa-Zahl. Belegt durch Datei-für-Datei-`diff` des aus `f10f471` rekonstruierten Alt-Baums gegen `.harness/baseline/v5.11.0/`.
- **klasse:** delta-kategorie-ueberzaehlt

**F-6**

- **kategorie:** LOW
- **quelle:** `BEO-005` (chronologische Tabelle kippt still ihre Richtung) / Maintainability
- **pfad:** `harness/conventions.md:139`
- **befund:** Die neue Zeile für `MR-029` ist in §Aufgelöste Adaptionen zwischen `MR-026` und `MR-027` eingefügt worden. Die Spalte war bis dahin strikt aufsteigend (`001 002 003 008 009 010 011 012 014 016 017 018 019 020 022 024 026 027 028`); sie lautet jetzt `… 024 026 029 027 028`. Die Tabelle ist die Auffind-Form der Pin-Serie („nur ID und Nachfolger, damit die Kette auffindbar bleibt"); ein Leser, der die Kette der Reihe nach abgeht, trifft `MR-029` vor `MR-027`/`MR-028` und findet den Nachfolger-Zeiger auf `MR-030` vor den Zeilen, die auf `MR-029` verweisen. Keine der sechs `structure.table-order`-Regeln deckt diese Tabelle (`.d-check.yml:258-277` nennt Lastenheft §7, Spezifikation §7, Roadmap-Drift-Log, Roadmap §Abgeschlossene Wellen, `version.md` §Verlauf, Handbuch §11).
- **verifizierbar:** nein — die Tabelle ist im `structure`-Profil nicht konfiguriert; `make gates` bleibt grün. Belegt durch Spalten-Extraktion an `HEAD` gegen `1f534ea`.
- **klasse:** register-tabelle-ausser-reihe

**F-7**

- **kategorie:** LOW
- **quelle:** `harness/conventions.md` §Adaptions-Block („was hier steht, liest **jeder** Agentenlauf, aufgelöste Adaptionen gehören nicht in diesen Pfad") · `MR-029` §Adaption, Bullet *Hebungs-Zensus (Checkliste für den Nachfolger)*
- **pfad:** `harness/conventions/MR-030-baseline-v5110.md` (durchgehend) · `harness/conventions/done/MR-029-baseline-v590.md:36-46`
- **befund:** `MR-029` trägt den Drei-Klassen-Zensus ausdrücklich als „Checkliste für den Nachfolger". Mit dem Lifecycle-Move liegt sie in `conventions/done/` — dem Pfad, den derselbe Konventionsspeicher als *nicht* von jedem Lauf gelesen deklariert. `MR-030` schreibt sie nicht fort: es übernimmt die zweite Richtung der Klasse („**Historische Verweise bleiben stehen**", ein eigener Absatz) und lässt die erste — die Aufzählung der drei Spiegel-Klassen — aus. Der aktive Konventionsspeicher trägt damit nach dem Move nur noch die Hälfte der Prozedur; die andere Hälfte hängt an `BEO-008` und am jeweiligen Slice-Plan. Die ausgelassene Richtung ist die, in der die Klasse in diesem Slice erneut eingetreten ist (F-1, F-2).
- **verifizierbar:** nein — kein Gate misst den Inhalt einer MR gegen ihren Vorgänger. Belegt durch Volltext-Vergleich beider Einträge.
- **klasse:** prozedur-wandert-mit-nach-done

**F-8**

- **kategorie:** INFO
- **quelle:** `AGENTS.md` §3.3, Ausnahme MR-/Wellen-Lifecycle-Move („Alles Übrige bleibt Commit 2") · `MR-013`
- **pfad:** Commit `5331466`
- **befund:** Der Commit bündelt den `git mv` von `MR-029` nach `conventions/done/` samt dessen zehn Link-Tiefen-Fixes mit 38 Verweis-Hebungen in 19 weiteren Dateien, dem neu angelegten `MR-030`, beiden Index-Tabellen und dem Retarget in `done/slice-117`. Die benannte Ausnahme erlaubt dem Move-Commit ausdrücklich nur die Link-Tiefen-Fixes der bewegten Datei selbst; alles Übrige gehört in einen zweiten Commit, und die Ausnahme trägt ihre Begründung im Rename-Score. Beobachtbare Folge gibt es hier keine: `git` weist den Move mit `similarity index 75%` als Rename aus, `git log --follow` bleibt tragfähig, und die Botschaft deklariert den Move. Die Grenze selbst ist ohne Deklaration überschritten und wird als Präzedenz für die nächste Pin-Hebung gelesen, deren bewegter Eintrag größer sein kann.
- **verifizierbar:** ja — `git show --find-renames` weist Rename und Similarity aus; die Commit-Zusammensetzung ist aus `git show --stat` ablesbar.
- **klasse:** move-commit-grenze-ueberschritten

---

## Negativbefunde

1. **Delta-Messung, selbst nachgezählt.** Alt-Baum aus `f10f471` rekonstruiert, Datei-für-Datei gegen `.harness/baseline/v5.11.0/` gehalten: beide Bäume tragen **52** Dateien mit **identischer** Dateiliste, **30** unterscheiden sich. Alle sechs Umfangs-Angaben der MR-Tabelle stimmen exakt: `grundlagen-source-precedence.md` +75/−4, `grundlagen-referenz-richtung.md` +30/−1, `grundlagen-durchsetzungsschicht.md` +8/−1, `grundlagen-begriffe.md` +7/−2, `templates/spec/lastenheft.template.md` +22/−7, `regelwerk/README.md` +3/−3. `SHA256SUMS` als Manifest: +29/−29. Ohne Befund außer der Kategorisierung in F-5.
2. **„23 nur Versions-Stempel" ist haltbar — mechanisch geprüft.** Genau 23 Dateien ändern genau eine Zeile; jede dieser Zeilen ist nach Normalisierung `v5.9.0`↔`v5.11.0` byte-identisch zur Vorgänger-Fassung. Es verbirgt sich dort kein Inhalt. (Randnotiz ohne Befundcharakter: alle 23 sind Quell-**URLs**, keine Pfade — die MR-Formulierung „Quell-URL bzw. Pfad" ist weiter gefasst als die Messung.)
3. **Bundle-Integrität und Authentizität.** `--verify` ohne Argument (liest den neuen Pin) Exit **0**, „verify ok (51 Dateien, vollständig)"; das Manifest zählt 51 Zeilen, der Baum 51 Dateien ohne `SHA256SUMS` — die „51 Dateien"-Angabe und die „52 Bundle-Dateien"-Angabe der MR meinen dasselbe und sind beide richtig. `--verify v5.9.0` Exit **1** (Baum entfernt, erwartet). `--check-latest` Exit **0**: Currency („Pin `v5.11.0` ist der neueste Release-Tag") **und** Content („gepinnter Tag upstream unverändert, Bytes == vendored `SHA256SUMS`") — damit ist „Kein Handanlegen an den Bäumen" unabhängig belegt und der vendorte Baum als das Release-Asset am Tag ausgewiesen.
4. **Klasse 1, alle 34 Pfad-Hebungen einzeln auf Zeitform geprüft — keine Über-Hebung.** Eigene Zählung über den Diff ergibt exakt 34 gehobene `baseline/v5.9.0/`-Pfade in 20 Dateien. Jede Stelle spricht über den **gegenwärtigen** Kanon: die neun `Ersetzt-Baseline-Regel`-Links der aktiven `MR-*`-Dateien (das ist die Regel, an deren Stelle die Adaption *heute* tritt), die fünf Kanon-Zeiger in `AGENTS.md`, die Ziel-Form-Zeiger in `spec/architecture.md:7` und `spec/spezifikation.md:9`, der Format-Regel-Zeiger der Roadmap, der Steering-Loop-Zeiger in `docs/plan/planning/README.md:17`, die Bezug-Felder der drei lebenden Slices und die drei relativen Zeiger in `.harness/skills/reviewer.md`. Nicht gehoben und zu Recht nicht gehoben: `AGENTS.md:169/199/201` („seit dem `v5.9.0`-Bump", Herkunfts-Anker), `.d-check.yml:106` (Historie einer früheren Hebung), `roadmap.md:80` (Graphknoten der geschlossenen welle-81), sämtliche `done/`-Slices, Review-Reports und aufgelösten `MR-*`. Ohne Befund.
5. **Klasse 2, Text-gegen-Ziel.** Auf der `+`-Seite stehen exakt vier gehobene Tag-URLs; bei beiden `conventions.md`-Stellen mit Versions-Linktext stimmen Text und Ziel überein (`v5.11.0` ↔ `releases/tag/v5.11.0`, `ai-harness-course@v5.11.0` ↔ `tree/v5.11.0`). Ein repo-weiter Scan aller Markdown-Links, deren **Text** eine andere `vX.Y.Z` nennt als ihr **Ziel**, liefert genau einen Treffer — `docs/reviews/2026-07-17-slice-075-implementation-r1.md:47`, ein Produkt-Versions-Zitat ohne Baseline-Bezug, vorbestehend und eingefroren. Der Rest von Klasse 2 steht in F-1.
6. **Ziel-Achse der Anker ist hier nicht blind.** Gegenprobe: ein erfundener Anker in einem Link **in den vendorten Baum** (`../.harness/baseline/v5.11.0/regelwerk/modul-06-roadmap.md#gibt-es-diesen-anker-nicht`) wird als `anchor-missing` gemeldet, Exit 2 — obwohl `.harness/baseline/**` in `scan.ignore` steht. Die 34 gehobenen Pfad-Links samt ihren Abschnitts-Ankern sind damit vom grünen `doc-check` tatsächlich gedeckt; die Blindheit aus F-3 betrifft ausschließlich `codepaths` und ausschließlich Inline-Code.
7. **Keine Zwischenstufe.** Repo-weit **0** Treffer für `v5.10.0` (Markdown, YAML, Skripte, `Makefile`, `.github/`). Das §5-Risiko „Verweise auf die Zwischenstufe kann es nicht geben" ist damit gemessen statt geglaubt.
8. **Kein toter Markdown-Link auf den entfernten Baum.** `grep` über alle Link-Ziele mit `baseline/v5.9.0`: **0** Treffer. Die verbliebenen Nennungen sind Inline-Code/Prosa in fünf eingefrorenen Review-Reports (`2026-08-22-slice-117/-118/-119/-120/-121`) und `docs/plan/planning/done/slice-117-baseline-v590-bump.md:65`; hinzu kommt ein synthetisches Test-Fixture (`internal/hexagon/core/rules/versions_test.go:312`), das keinen Verweis darstellt. Die Aufzählung der Botschaft ist insoweit vollständig und richtig — die Begründung dafür, warum das folgenlos bleibt, ist es nicht (F-3).
9. **MR-013-Link-Tiefen der bewegten Datei sind vollständig.** Alle zehn relativen Ziele in `harness/conventions/done/MR-029-baseline-v590.md` sind umgestellt und lösen vom neuen Ort auf: `../../conventions.md` (sechsmal, mit Ankern), `../../../AGENTS.md`, `../../README.md`, `../../../.harness/skills/reviewer.md`, `../../../spec/lastenheft.md`, sowie die drei Geschwister-Zeiger `MR-027-…`/`MR-028-…` ohne `done/`-Präfix. Kein Rest auf altem Tiefen-Stand; `make doc-check` Exit 0 bestätigt es.
10. **Der Eingriff in `done/slice-117` ist vertretbar.** `done/`-Slices sind nicht gate-immutabel (`immutable.paths: ["docs/plan/adr/[0-9]*.md"]`); die in AGENTS §3.7/§5 benannte Bestands-Ausnahme betrifft das `**Status:**`-Feld, nicht den Dateikörper. AGENTS §3.3 und `MR-013` schreiben für Lifecycle-Moves ausdrücklich vor, **Pfad-Verweise auf die bewegte Datei** mitzuziehen; die im Slice-Plan §2 Schritt 4 genannte `ignore-refs`-Alternative gilt Verweisen auf ein **entferntes** Artefakt, für die es kein neues Ziel gibt — hier existiert das Ziel und ist nur umgezogen. Link-Text und Link-Ziel wurden gemeinsam gezogen, die Aussage des Steering-Loop-Eintrags bleibt wahr und auflösbar. Ohne Befund.
11. **Scope-Treue gegen §3.** Außerhalb des vendorten Baums enthält der Diff keine einzige Zeile, die eine Regel des neuen Stands anwendet. Maschinell geprüft: jede geänderte Zeile außerhalb `.harness/baseline/**` und außerhalb der neuen `MR-030` ist entweder eine reine `v5.9.0`→`v5.11.0`-Substitution oder Lifecycle-Buchführung (Ruhe-Marker der Roadmap, drei `open/`→`in-progress/`-Zeiger, beide Index-Zeilen, die MR-029-Tiefen, das slice-117-Retarget). Insbesondere ist die Vollständigkeits-Zusage aus Kurs-Welle 94 **nicht** angewandt, obwohl ihre Verletzung bekannt ist; `MR-030` sagt das ausdrücklich („Diese MR hebt den Pin, sie behauptet keine Konformität"). Kein Delta-Audit im Diff.
12. **Gate-Belege der Botschaften exakt reproduziert.** `make doc-check` Exit **0**, „457 Datei(en) geprüft, 0 Befund(e)". `make gates` Exit **0**, acht Gates, „coverage-gate: OK — Coverage 94.80% erfüllt Schwelle 93%" — beide Zahlen wortgleich zur Botschaft von `96cbd6a`. Die Datei-Arithmetik der Kette geht auf: `MR-030` ist die einzige neu **gescannte** Datei über die vier Commits (die 52 Bundle-Dateien nimmt `scan.ignore` aus), also 456 vor und 457 ab `5331466` — genau die Zahlen, die `cc8fffa`/`f10f471` (456) und `5331466`/`96cbd6a` (457) nennen. „Acht Gates" stimmt gegen `Makefile:158`.
13. **`MR-030`-Form gegen Vorlage und Präzedenz.** Alle Pflichtfelder der vendorten `MR-NNN-titel.template.md` vorhanden (Status, Datum, Geltungsbereich, Ersetzt-Baseline-Regel, Adaption, Begründung, Auflösungs-Trigger); `Löst auf` und `Ausgelöst durch Baseline-Stand` als Paar gesetzt — die Auflage, an der der Vorgänger-Review anschlug. `Löst auf` zeigt korrekt auf die **Index-Zeile** (`../conventions.md#mr-029`), nicht auf die Eintrags-Datei, und trägt die Begründung dafür mit. Nummer ist die nächste freie des dichten Zählraums. Der Geltungsbereich nennt die tatsächlich berührten Träger (`AGENTS.md`, `harness/README.md`, aktive `MR-*`, `.harness/skills/reviewer.md`, Spec-Straten, Planning-Docs) und deckt sich mit den 34 gehobenen Pfaden.
14. **Beide Index-Tabellen sind sachlich nachgezogen.** Die aktive Zeile ist ersetzt (`MR-029` raus, `MR-030` rein, Geltungsbereich auf `.harness/baseline/v5.11.0/` gehoben, Ersetzt-Baseline-Regel-Zelle wortgleich zur Vorgänger-Zeile); die aufgelöste Zeile trägt ihren Nachfolger als Index-Anker `[MR-030](#mr-030)`. Beide `<a id>`-Anker je Zeile (Voll-Slug **und** Kurzform) sind erhalten, `#mr-029` löst weiter auf und wird von `slice-128` §Bezug genutzt; kein Verweis auf den alten Pfad `harness/conventions/MR-029-…` außerhalb eingefrorener Reports. Die einzige Beanstandung an den Tabellen ist die Zeilenposition (F-6).
15. **Inhaltliche Stichprobe an den fünf Rumpf-Dateien.** Der größte Block (`grundlagen-source-precedence.md` +75) trägt tatsächlich die Vollständigkeits-Zusage samt Prüffrage, Waisen-Regel und der Rolle des Werkzeug-Einstiegs; `grundlagen-durchsetzungsschicht.md` trägt dieselbe Rolle als ausdrücklich **nicht** vierten Bindepunkt; `grundlagen-begriffe.md` ergänzt den Begriff *Change Request*. Die Zuordnung der MR („Kurs-Welle 94 bringt die Vollständigkeits-Zusage in `grundlagen-source-precedence.md` und die Rolle der Werkzeug-Einstiegsdatei in `grundlagen-durchsetzungsschicht.md`") ist am Bytes-Bestand belegt.
16. **Nicht lokal verifizierbar, deshalb nur benannt.** „Kurs-Welle 94 nennt ihn im Kurs-CHANGELOG als Auslöser" — das self-contained Bundle trägt kein CHANGELOG, und der Kurs-Arbeitsbaum gehört nicht zum Prüf-Gegenstand; belegt ist nur der *Inhalt* der Zusage (Negativbefund 15). Ebenso sind die Prozess-Erzählungen aus `5331466` („der erste Zensus war unvollständig … zweiter Zensus fand sie", „der `doc-check` hat im ersten Anlauf vier Befunde aus dem MR-Move gefangen") aus einem Ein-Commit-Diff nicht nachprüfbar; sie sind in sich stimmig, weil `.harness/skills/reviewer.md` tatsächlich relativ (`../baseline/…`) verweist und ein Pfad-`grep` auf `.harness/baseline/v5.9.0/` diese drei Stellen zwangsläufig verfehlt.
17. **Zähl-Randnotiz ohne Befundcharakter.** Der Betreff von `5331466` nennt „38 Verweise"; das ist die Summe der Klassen 1 und 2 (34 + 4). Der Ellipsen-Pin der Klasse 3 in `MR-021:20` ist zusätzlich gehoben, steht aber nicht in der 38. Der Botschaftskörper zählt alle drei Klassen einzeln auf, die Zahl ist damit nicht irreführend.
18. **Die benannte mechanische Form ist weiterhin nicht gefahren — deklariert, kein Befund.** `.d-check.yml:241-246` konfiguriert `versions` mit einem einzigen `pin-pattern` (`ghcr.io`), nicht mit der seit slice-122 möglichen Muster-Liste. `BEO-008` hält ausdrücklich fest, das eigene Profil fahre den Baseline-Abgleich nicht und der Entscheid sei ein eigener; der Slice-Plan §7 wiederholt es. F-1 und F-2 sind genau in der Lücke eingetreten, die das Register benennt.

---

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
| --- | --- | --- |
| HIGH | 0 | — |
| MEDIUM | 4 | F-1, F-2, F-3, F-4 |
| LOW | 3 | F-5, F-6, F-7 |
| INFO | 1 | F-8 |

## Verdikt

**Blockierend.** Kein HIGH; vier MEDIUM. Das Handwerk der Etappe ist über weite Strecken sauber und stellenweise besser belegt als bei den Vorgänger-Hebungen: die Delta-Messung ist bis auf eine Kategorie-Zuordnung exakt reproduzierbar, die Bundle-Authentizität ist upstream gegengeprüft, die Klasse-1-Hebung ist vollständig und in **beide** Zeitform-Richtungen korrekt, und die Scope-Grenze gegen das Delta-Audit hält ohne eine einzige Ausnahme.

Blockierend sind zwei Sachbefunde und zwei Beleg-Befunde. `harness/README.md:60` (F-1) ist die dritte Wiederholung derselben Klasse an derselben Zeile und bringt `BEO-008` auf Zähler 4; `slice-127:155-159` (F-2) ist dieselbe Klasse in ihrer dritten Ausprägung. Beide entstehen dort, wo der Zensus seine Vollständigkeit behauptet, statt sie zu messen — das ist zugleich `BEO-011`. F-3 und F-4 betreffen nicht die Arbeit, sondern die Belege: einmal wird ein grünes Gate als Beleg für eine Achse geführt, auf der es blind ist, einmal wird eine Prüfung für nicht durchgeführt erklärt, die im selben Commit vorliegt. Beide Sätze reisen in der dauerhaften Aufzeichnung mit und würden von der Etappe B als gesichert übernommen.

Der Slice-Plan hat die eintretende Klasse in §5 als erstes Risiko benannt und ihren Zähler-Stand mitgeführt; der Ausgang dieses Risikos ist damit für die Closure-Notiz bestimmbar. Die drei LOW und das INFO sind unabhängig davon zu entscheiden.
