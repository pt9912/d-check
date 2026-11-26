# MR-081 — Baseline-Pin-Hebung auf `v7.0.0` (achtzehnter Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Datum:** 2026-11-26
- **Geltungsbereich:** §Baseline, pin-gebundene Verweise,
  `.harness/baseline/v7.0.0/`
- **Ersetzt-Baseline-Regel:** — *(Pin-Fortschreibung im Bundle-Layout des
  Vorgängers; keine inhaltliche Abweichung)*
- **Adaption:** Der vendorte Bestand steht auf
  [`v7.0.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v7.0.0)
  (Kurs-Welle 168, 2026-11-23). Zwischen dem alten Pin und diesem liegt kein
  weiterer Release-Tag (die Release-Liste, gelesen am 2026-11-26, führt nach
  `v6.17.0` nur `v6.18.0`, `v6.18.1` und `v7.0.0`).
  **Der alte Baum `v6.18.1` entfällt im Bump** — die Abweichung vom Kanon trägt
  [`MR-021`](../conventions.md#mr-021) (1) schon; er liegt in der git-Historie
  beim Vor-Bump-Commit `738d2bca` ([`MR-052`](../conventions.md#mr-052)).

  **Der Delta, gemessen und nicht geschätzt** (Vergleich beider Bäume, beide
  noch vorhanden; der alte Baum aus `git archive 738d2bca`; ohne die
  `<!-- Quelle:`-Zeilen, Versionsstempel auf einen Platzhalter normalisiert):
  54 Dateien im alten Bundle, 55 im neuen (`find` über `regelwerk/` und
  `templates/`; vendort je eine mehr für `SHA256SUMS`: 55 und 56), **47**
  gleichnamig in beiden, **7 umbenannt**, **1 neu**. Von den 47 sind **38**
  byte-verschieden (`cmp` je Datei; 862 Diff-Zeilen, `diff` mit `<`/`>`), davon
  **31** inhaltlich verschieden (19 Regelwerk, 12 Vorlagen; **804** Diff-Zeilen)
  und **7** nur im Versionsstempel und Quell-Verweis. Größte Bewegungen:
  `regelwerk/modul-14-docker-harness.md` mit 258 Diff-Zeilen,
  `templates/.d-check.yml` mit 74, `regelwerk/modul-04-adrs.md` mit 57,
  `regelwerk/modul-05-planning-harness.md` mit 49, `templates/README.md` mit 35.

  **Inhaltliche Bewegungen — vollständig nach Datei, keine davon übernommen**
  (Abgrenzung, slice-295; die Namensregel führen
  [`slice-297`](../../docs/plan/planning/done/slice-297-spezifikation-heisst-specification.md)
  und [`slice-298`](../../docs/plan/planning/open/slice-298-neue-namen-wave-archive-anker.md),
  die Klassifikation mit Empfehlung zu allen übrigen
  [`slice-296`](../../docs/plan/planning/open/slice-296-baseline-v700-neuerungen-sichtung.md);
  die Zahl hinter dem Pfad sind die Diff-Zeilen):

  1. **Dateinamen sind englisch** — `grundlagen-harness-dateien.md` (23;
     neuer Absatz in §Verzeichniskonvention: `wave-<Kennung>.md`,
     `specification.md`, das Präfix der Welle folgt der Regel; dazu
     `seit wave-<Kennung>` im Kommentar-Abschnitt, die Platzhalter
     `<wave-id>` und der Vorlagen-Name; der Rest der 23 gehört zu 3). Nachgezogen in:
     `grundlagen-source-precedence.md` (14; `specification.md`,
     `slice-<wave-name>-…`), `grundlagen-referenz-richtung.md` (12),
     `grundlagen-traceability.md` (24; `seit wave-`, `done/wave-<Kennung>-results.md`,
     `wave-results.template.md`), `modul-02-harness-bootstrap.md` (14; Dateinamen,
     Vorlagen-Namen), `modul-03-spec.md` (21; `specification.md`, Anteile zu 2),
     `modul-04-adrs.md` (57; `NNNN-title.template.md`, Anteile zu 2),
     `modul-05-planning-harness.md` (49; `done/slice-<Kennung>-archive.zip`,
     Anteile zu 2 und 4), `modul-06-roadmap.md` (28; `done/wave-<Kennung>-results.md`,
     `done/<wave-id>/archive.zip`, `archive-stub-slice.template.md`,
     `archive-stub-wave.template.md`, `wave.template.md`, `wave-results.template.md`),
     `modul-07-carveouts.md` (4), `modul-08-agentenrollen.md` (14; Anteil zu 4),
     `modul-09-implementierung.md` (4; Anteil zu 3), `modul-10-review-harness.md`
     (28; `seit wave-`, `done/slice-<Kennung>-archive.zip`, Anteile zu 2),
     `modul-13-quality-gates.md` (14; `seit wave-`, Anteile zu 2 und 3),
     `modul-16-produktiver-betrieb.md` (4; `done/wave-<Kennung>-closure.md`, Anteil zu 3);
     die sieben umbenannten Vorlagen (unten) und
     `templates/.d-check.yml` (74; Globs `wave-*.md`, Token `wave-…`),
     `templates/README.md` (35), `templates/AGENTS.template.md` (24),
     `templates/docs/plan/planning/slice.template.md` (17),
     `templates/spec/lastenheft.template.md` (13),
     `templates/docs/plan/planning/roadmap.template.md` (10),
     `templates/docs/plan/planning/README.template.md` (6),
     `templates/harness/conventions.template.md` (4),
     `templates/harness/README.template.md` (4),
     `templates/docs/plan/planning/observation.template.md` (4).
     **Dieses Repo ist betroffen und hat entschieden** (Auftraggeber,
     2026-11-25, „ab jetzt für Neues, der Bestand bleibt"): die Spezifikation
     [`slice-297`](../../docs/plan/planning/done/slice-297-spezifikation-heisst-specification.md),
     Wellen, Archive und Anker
     [`slice-298`](../../docs/plan/planning/open/slice-298-neue-namen-wave-archive-anker.md).
     **In diesem Slice nur die Namen der sieben Vorlagen-Dateien in den Verweisen** — ein
     Verweis auf eine Datei, die es nicht mehr gibt, ist ein toter Verweis.
  2. **Bedrohung, Grenze, Ausgang** (Vertrauensgrenzen in der Sicht, drei
     Ausgänge je Bedrohung, kein eigenes Register) — neuer Abschnitt in
     `modul-04-adrs.md`; dazu `modul-03-spec.md` (Fehlannahme, Ziel-Form der
     Architektur-Sicht), `modul-05-planning-harness.md` (eine Bedrohung ist kein
     Risiko des Slice), `modul-10-review-harness.md` (Prüffrage im Plan- und
     Design-Review, *ein akzeptierter Ausgang ist kein Verstoß*, Kontext-Eingang
     um `spec/architecture.md` und aktive ADRs), `modul-13-quality-gates.md`
     (ein Bedrohungsmodell speist Anforderung, Sicht und ADR, für einen Wächter
     höchstens Kandidaten). Nachgezogen in `templates/spec/architecture.template.md`
     (21; neuer §6 *Vertrauensgrenzen und Datenflüsse*),
     `templates/spec/lastenheft.template.md` (die Sicherheits-Anforderung als
     Eigenschaft statt Bedrohung), `templates/docs/plan/adr/NNNN-title.template.md`
     (7; eine ADR, die eine Bedrohung akzeptiert),
     `templates/docs/plan/planning/slice.template.md` (Kommentar zu den Risiken),
     `templates/.harness/skills/reviewer.template.md` (15; Kontext-Eingang,
     Ausnahme beim Sicherheits-Anti-Pattern, MEDIUM-Fall) und
     `templates/README.md`. **Nicht übernommen**; Klassifikation und Empfehlung
     in [`slice-296`](../../docs/plan/planning/open/slice-296-baseline-v700-neuerungen-sichtung.md).
  3. **Umgebungs-Zusage, Tool-Pin und Tool-/Versions-Check** — `modul-14-docker-harness.md`
     (258; neu: *Die Umgebungs-Zusage — Anforderung vor Träger*, *Der zweite Träger:
     Tool-Pin und Tool-/Versions-Check* samt Prüfformen, Inventar, Ausführungsort und
     Probe im leeren Kontext; die Überschrift *Zwei Formen des Reproduzierbarkeits-Ankers*
     heißt jetzt *Zwei Formen des Belegs beim Docker-Träger*), `grundlagen-begriffe.md`
     (1; neue Zeile *Umgebungs-Zusage*), `grundlagen-harness-dateien.md`
     (Reproduzierbarkeits-Bindung), `modul-09-implementierung.md` (die
     Begründung der Docker-only-Regel trägt die Umgebungs-Zusage statt
     „Supply-Chain-Defense"), `modul-12-replay-evaluierung.md` (6; Rang 1 hat eine
     maschinelle Vorstufe), `modul-13-quality-gates.md` (neue Zeile in der
     Gate-Typ-Tabelle und die Abgrenzung des Checks zum Wächter),
     `modul-16-produktiver-betrieb.md`. Nachgezogen in
     `templates/AGENTS.template.md` (§3.1), `templates/.d-check.yml` (die
     Blöcke `planning` und `pinsync` als Beispiel, `doc-manual` in der Modulliste),
     `templates/README.md` und dem **neuen**
     `templates/harness/sensors/environment-check.template.md` (74 Zeilen).
     **Nicht übernommen**; Klassifikation und Empfehlung in
     [`slice-296`](../../docs/plan/planning/open/slice-296-baseline-v700-neuerungen-sichtung.md).
  4. **Paralleler Lauf: ein Zweig, ein Arbeitsbaum** (`git worktree`) —
     `modul-05-planning-harness.md` (drei neue Absätze unter §Lifecycle als State
     Machine, der WIP-Satz), `modul-08-agentenrollen.md` (der Rolleninhaber hängt
     an Zweig **und** Arbeitsbaum); nachgezogen in `templates/AGENTS.template.md`
     (Absatz *Parallele Läufe*), `templates/docs/plan/planning/slice.template.md`
     (`Verantwortlich:` nennt Person und Zweig) und `templates/.d-check.yml`
     (`.claude/worktrees/**` im Scan-Ausschluss) sowie `templates/README.md` (die
     `.d-check.yml`-Zeile). **Nicht übernommen**; Klassifikation und Empfehlung in
     [`slice-296`](../../docs/plan/planning/open/slice-296-baseline-v700-neuerungen-sichtung.md).
  5. **Welle-Vorlage: Slice-Tabelle mit Links** —
     `templates/docs/plan/planning/wave.template.md` (28; die Kennung in der
     Tabelle ist ein Link auf die Slice-Datei samt `<lifecycle>`-Verzeichnis).
  6. **Bundle-Stand** — `regelwerk/README.md` (2; Kurs-Welle 161 → 168, Datum
     2026-10-10 → 2026-11-23).

  **Die sieben umbenannten Vorlagen** (Dateiname; Inhalt: nur Namen und
  Platzhalter nachgezogen, Diff-Zeilen in Klammern):
  `templates/docs/plan/planning/welle.template.md` → `wave.template.md` (28; dazu
  Bewegung 5), `welle-results.template.md` → `wave-results.template.md` (16),
  `archiv-stub-slice.template.md` → `archive-stub-slice.template.md` (10),
  `archiv-stub-welle.template.md` → `archive-stub-wave.template.md` (10),
  `templates/docs/plan/adr/NNNN-titel.template.md` → `NNNN-title.template.md` (7; dazu
  Bewegung 2), `templates/harness/conventions/MR-NNN-titel.template.md` →
  `MR-NNN-title.template.md` (2), `templates/spec/spezifikation.template.md` →
  `specification.template.md` (2). **Nachgezogen in den Verweisen dieses Repos**
  (gemessen: **9** Zeilen in **4** Dateien — `.claude/commands/close-welle.md` 4,
  `.claude/commands/plan-welle.md` 2, `harness/conventions.md` 1,
  `tools/archive-wave/stub.go` 2 Kommentarzeilen, Kopplung an die Vorlagen-Namen):
  Befund `git grep -n "welle\.template\|welle-results\.template\|archiv-stub-\(slice\|welle\)\.template\|titel\.template" -- . ':!.harness/baseline' ':!docs/plan/planning' ':!docs/reviews' ':!harness/conventions/done' ':!harness/conventions/MR-081-baseline-v700.md'`
  war vor dem Nachzug 9 Zeilen, danach 0 (dieser Eintrag nennt die Namen selbst und ist
  deshalb ausgenommen; ohne die Ausnahme zeigt das Kommando 7 Zeilen, alle hier). Die
  Treffer in den Plänen dieser Welle, in `docs/reviews/` und in `conventions/done/` nennen die
  Namen als Vergangenheits-Aussage und bleiben.

  **Die sieben Dateien, die nur ihren Stempel ändern:** `regelwerk/modul-00-einfuehrung.md`,
  `modul-01-entwicklungszyklus.md`, `modul-11-verification.md`, `modul-15-observability.md`,
  `grundlagen-bootstrap.md`, `grundlagen-klassifikation.md` und
  `grundlagen-durchsetzungsschicht.md` — Versionsstempel und Quell-Verweis
  (0 Diff-Zeilen nach der Normalisierung).

  **Die entfallene Überschrift trifft kein lebendes Dokument:** `git grep -n
  'zwei-formen-des-reproduzierbarkeits' -- . ':!.harness/baseline'` ist leer
  (Stand vor dem Bump), und nach dem Entfernen des alten Baums meldet
  `make doc-check` keinen Anker-Bruch.

  **Spiegel-Klassen, gemessen am echten Vorzustand** (HEAD `738d2bca`, ohne
  `.harness/baseline/`, **nicht** ohne `.harness/` — die Skills liegen dort;
  `git grep -l` und `git grep -o 'v6\.18\.1'`, gezählt wird die Datei und das
  Vorkommen, nicht der Verweis auf den Baum): **77** Dateien außerhalb des
  vendorten Baums nannten `v6.18.1` — **190** Vorkommen (175 Zeilen). **Frozen
  (17 Dateien, 51 Vorkommen):** die `done/`-Slices und -Wellen (8 / 29), die
  Review-Reports (8 / 21) und die CHANGELOG-Historie (1 / 1). **Davon bleiben 14
  Dateien, 40 Vorkommen byte-stabil** (7 / 28, 6 / 11, 1 / 1); die **drei**
  übrigen tragen je eine geänderte Zeile, nicht im Ventil (unten). Gemessen mit
  `for f in $(git grep -l 'v6\.18\.1' 738d2bca -- docs/plan/planning/done
  docs/reviews CHANGELOG.md | sed 's/^738d2bca://'); do git diff --quiet 738d2bca
  HEAD -- "$f" && git grep -o 'v6\.18\.1' 738d2bca -- "$f" | wc -l; done`
  (eine Zahl je byte-stabiler Datei: 14 Zeilen, Summe 40; mit nur
  `docs/plan/planning/done` bzw. nur `docs/reviews` als Pfad 7 und 6 Dateien).
  `conventions/done/` und eingehende CRs nannten den Tag nicht.
  **Frozen hinter dem Ventil (5 / 5):**
  die `Accepted`-ADRs 0112 bis 0116, je mit einem Link in der Regeln-Zeile — **nicht**
  retargetet ([`AGENTS.md`](../../AGENTS.md#35-adrs-sind-nach-accepted-immutable) §3.5,
  siehe unten). **Pläne dieser Welle (5 / 25):** die Welle und die vier Slices nennen den
  alten Pin als Vorzustand und als Gegenstand; nachgezogen sind allein ihre
  `d-check:cite`-Spannen (unten). **Der Vorgänger-Eintrag (1 / 9):** [`MR-080`](../conventions.md#mr-080)
  nennt den Pin als Tatsache seiner Zeit und bleibt unverändert; er wandert nach
  `conventions/done/`. **Lebend retargetet (49 Dateien, 100 Vorkommen):**
  `AGENTS.md`, `harness/README.md`, `harness/conventions.md`, die drei
  `harness/rules/`-Dateien, 35 aktive `MR-*`, `.claude/agents/reviewer.md`, die beiden
  Skills unter `.harness/skills/`, `roadmap.md`, `planning/README.md`,
  `observations/README.md`, `spec/architecture.md` und `spec/spezifikation.md`
  (Rolle-Zeile; eingebettet, siehe [`slice-295`](../../docs/plan/planning/done/slice-295-baseline-v700-bump.md)
  §1). **Zeilenweise gelesen, nicht pauschal ersetzt:** `harness/conventions.md`
  (die Zeilen zum Pin retargetet; die Index-Zeile des Vorgängers wandert nach
  *Aufgelöste Adaptionen*), `observations/README.md` Z. 18, `roadmap.md` Z. 13,
  [`MR-021`](../conventions.md#mr-021) Z. 20 („aktuell …/v7.0.0/…"),
  [`MR-049`](../conventions.md#mr-049) Z. 14, [`MR-053`](../conventions.md#mr-053)
  Z. 9 und `AGENTS.md` Z. 345 („Baseline-Form `v7.0.0` §1"): jede Nennung ist eine
  Pin-Aussage über die **Form** des Abschnitts, nicht über den Stand der Zeit, und
  steht gegen den neuen Baum (die Überschriften §1 bis §8 der
  `slice.template.md` stehen unverändert, `grep -n '^## '` beider Bäume; nur die
  Zeilennummern verschieben sich; `modul-06-roadmap.md` §Roadmap-Struktur trägt
  die Aussage zu *Offene Wellen* im alten Wortlaut). **Anderer Gegenstand:**
  `tools/harness/selbstpruefung.sh` nennt die Version des Schwester-Werkzeugs,
  aus dem es adoptiert ist. **Die 7 Baseline-Symlinks unter `.claude/rules/`**
  lösen jetzt gegen `v7.0.0` (`for l in .claude/rules/*; do readlink "$l"; done |
  grep -c 'v6\.18\.1'` ergibt 0; `git grep` liest keinen Symlink-Inhalt).
  **Was nach dem Bump den alten Pin noch nennt** (`git grep -l 'v6\.18\.1' -- .
  ':!docs/plan/planning/done' ':!docs/reviews' ':!CHANGELOG.md'
  ':!harness/conventions/done' ':!.harness/baseline'`, Stand des Bump-Commits):
  **10** Dateien — dieser Eintrag, die `.d-check.yml` (die fünf Ventil-Einträge),
  die fünf `Accepted`-ADRs hinter dem Ventil und drei Pläne dieser Welle, die den
  alten Pin als Vorzustand und als Gegenstand nennen (`slice-295`, `slice-296`,
  die Welle selbst; ein Vergangenheits-Beleg, der mit ihrem Umzug nach `done/`
  aus der Menge fällt).

  **`ignore-refs` wächst, um fünf Einträge** — gemessen: nach dem Entfernen des
  alten Baums meldet `make doc-check` sieben eingefrorene Dateien mit einem
  Markdown-Link auf `.harness/baseline/v6.18.1/`; fünf davon, die `Accepted`-ADRs
  0112, 0113, 0114, 0115 und 0116 (je in der Regeln-Zeile), bleiben hinter dem
  Ventil, die zwei übrigen sind retargetet (unten). Die Form ist die von
  [`MR-069`](../conventions.md#mr-069): datei-skopiert, je ein Eintrag. Der Grund:
  die ADRs zitieren den Stand ihrer Zeit, und der zitierte Abschnitt (`modul-04`
  §Ziel-Form: ADR) ist in `v7.0.0` **nicht wortgleich** — er nennt die Vorlage
  `NNNN-title.template.md` statt `NNNN-titel.template.md` (je Abschnitt von der
  Überschrift bis zur nächsten Überschrift gelesen, alt gegen neu: eine
  abweichende Zeile, der Dateiname der Vorlage); ein Nachziehen verfälschte das Zitat
  ([`MR-069`](../conventions.md#mr-069),
  [`AGENTS.md`](../../AGENTS.md#35-adrs-sind-nach-accepted-immutable) §3.5), und der
  Nachzug verschöbe nur die Wand auf die nächste Hebung
  ([Nachtrag V1](../../docs/reviews/2026-11-26-welle-98-architect-verdikt.md)).
  Bestand: **35** statt 30 Einträge auf entfernte Bäume, über **13** statt 12 Tags
  (`grep -c 'refs: \[\?"\?\.harness/baseline/v' .d-check.yml`; Tags:
  `grep -o 'baseline/v[0-9.]*/\*\*' .d-check.yml | sort -u | wc -l`; Vorzustand
  mit `git show 738d2bca:.d-check.yml | …`). Der
  Auflösungs-Trigger von [`MR-069`](../conventions.md#mr-069) („kein neuer Eintrag über drei
  aufeinanderfolgende Pin-Hebungen") beginnt neu: die beiden vorigen Hebungen
  ([`MR-074`](../conventions.md#mr-074), [`MR-080`](../conventions.md#mr-080)) waren ohne
  Wachstum, diese ist es nicht (die fünf ADRs). [`slice-274`](../../docs/plan/planning/done/wellenlos/slice-274-trace-requirements-als-liste-von-quellen.md)
  und das [Verdikt zu welle-93](../../docs/reviews/2026-10-30-welle-93-architect-verdikt.md)
  sind wie bei [`MR-080`](../conventions.md#mr-080) auf den neuen Baum retargetet: die
  zitierte Stelle ist gelesen wortgleich (`modul-11` §Bewusstes Brechen, `modul-05`
  §Ziel-Form: Slice; `diff` je Abschnitt leer, 18 bzw. 58 Zeilen), und der Tausch hält
  die Gate-Senkung um zwei Dateien kleiner.

  **Drei Ausnahmen von der Frozen-Regel, nicht im Ventil** — je **1 Datei, 1 Zeile**,
  nur Link-Ziel: (1) der Verifikations-Report
  `docs/reviews/2026-10-20-slice-277-baseline-v6181-verify.md` verlinkt die
  Eintrags-Datei von [`MR-080`](../conventions.md#mr-080) am alten Pfad (Z. 4); mit
  dem Move nach `conventions/done/` wäre der Link tot, er ist auf den neuen Pfad
  nachgezogen, das Ziel lebt weiter. (2) `slice-274` und (3) das Verdikt zu welle-93
  (oben) tragen das Link-Ziel auf den neuen Baum. Die Alternative, je ein weiterer
  Ventil-Eintrag, wäre eine Gate-Senkung für einen Link, der auflösen kann (Plan-Änderung in
  [`slice-295`](../../docs/plan/planning/done/slice-295-baseline-v700-bump.md) §3).

  **`d-check:cite`-Neu-Ankern** — gemessen **15** lebende Direktiven auf den
  Baum (0 tag-only; ohne `conventions/done/`, `done/`-Slices und Reviews, die der
  `citations`-Scope ausnimmt, und ohne die acht Direktiven in den vier Plänen der
  Welle), **7** davon neu geankert, je reine Zeilenverschiebung bei **gelesenem,
  wortgleichem** Text (alter und neuer Bereich nebeneinander): Reviewer-Skill
  `templates/.harness/skills/reviewer.template.md` 87–88 → 98–99 und 99–100 → 110–111,
  `modul-10-review-harness.md` 99 → 115; [`MR-031`](../conventions.md#mr-031)
  `grundlagen-harness-dateien.md` 196 → 208; [`MR-035`](../conventions.md#mr-035)
  `grundlagen-begriffe.md` 50 → 51; [`MR-049`](../conventions.md#mr-049)
  `modul-05-planning-harness.md` 180–181 → 223–224 und 170 → 213. Die **8** übrigen
  stehen auf unverschobenen Zeilen (derselbe Bereich beider Bäume wortgleich):
  Closure-Note-Skill 3 (`modul-11` 83; die Vorlage 77–78 und 84),
  [`MR-031`](../conventions.md#mr-031) 2 (`modul-09` 20 und 203),
  [`MR-005`](../conventions.md#mr-005) 1, [`MR-039`](../conventions.md#mr-039) 1,
  [`MR-043`](../conventions.md#mr-043) 1. **Dazu die 8 in den Plänen:** `modul-05`
  373–374 → 416–417 (4×) und 379 → 422 (4×), wortgleich.

  **Zitat-Delta** ([`MR-039`](../conventions.md#mr-039)): wörtliche Zitate in
  `„…"`-Form (ab 25 Zeichen) aus den lebenden Markdown- und YAML-Dateien
  (`git ls-files -- '*.md' '*.yml' ':!docs/plan/planning/done' ':!docs/reviews'
  ':!CHANGELOG.md' ':!harness/conventions/done' ':!docs/plan/adr' ':!docs/plan/cr'
  ':!.harness/baseline'`, Symlinks ausgenommen). Je Datei über Zeilenumbrüche
  normalisiert (`tr '\n\t' '  ' | tr -s ' '`) und mit
  `grep -o '„[^„"“]*["“]'` gezogen; jedes Zitat wird in beiden Bäumen gesucht
  (alle `*.md` je Baum, gleich normalisiert, `grep -qF`). Gemessen: **36**
  Fundstellen stehen in mindestens einem der beiden Bäume, **35** im Wortlaut
  beider Bäume gleich, **1** nur im alten — die Überschrift *Zwei Formen des
  Reproduzierbarkeits-Ankers* im Plan dieses Slice, der ihren Wegfall nennt. Weitere
  **381** Zitate stehen in keinem der beiden Bäume (andere Quellen) und gehören nicht
  zur Menge. Die Zitate ohne Anführungszeichen-Form (Paraphrase) sind damit nicht gemessen.
- **Begründung:** Der Nachtlauf `upstream-drift` meldete den neuen Release
  (`make baseline-freshness`); das Release `v7.0.0` lag vor. Die Hebung ist eine
  reine Fortschreibung; die inhaltlichen Bewegungen (oben) warten auf
  [`slice-296`](../../docs/plan/planning/open/slice-296-baseline-v700-neuerungen-sichtung.md)
  und seine Folge-Slices, die Namensregel auf `slice-297` und `slice-298`.
- **Auflösungs-Trigger:** die nächste Pin-Hebung.
- **Löst auf:** [`MR-080`](../conventions.md#mr-080)
- **Ausgelöst durch Baseline-Stand:** `v7.0.0`
