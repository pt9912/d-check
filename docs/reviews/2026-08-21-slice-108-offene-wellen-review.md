# Review — slice-108: Roadmap auf §Offene Wellen (Etappe C-1, Commit `a46f09c`)

- **Review-Art:** Code-/Form-Review (unabhängiger Erst-Review, frischer Kontext) —
  geprüft gegen die vendorte Baseline v5.6.0 (`modul-06-roadmap.md`
  §Roadmap-Struktur + `roadmap.template.md`), die `planning`-Config-Fläche
  (`DC-FA-PLAN-001`, Spez §2) und das reale Gate-Verhalten des Images.
- **Gegenstand:** slice-108 (welle-78, Etappe C-1) / Commit `a46f09c` (HEAD,
  1 vor origin, unpushed).
- **Skill:** `reviewer.md` v1.4.0 (2026-08-15). **Modell-ID:** claude-fable-5.
- **Datum:** 2026-08-21.
- **Eingangs-Kontext:** Slice-Plan
  `docs/plan/planning/in-progress/slice-108-roadmap-offene-wellen.md`;
  Audit-Befund C-1 (`docs/plan/planning/done/slice-107-baseline-v560-delta-audit.md` §9);
  `MR-024`/`MR-013`; `DC-FA-PLAN-001` (nur Config, kein Produkt-Code);
  Hard Rules `AGENTS.md` §3. Kein DoD-Abhaken (Verifikations-Rolle).
- **Messläufe (selbst gefahren, Image `d-check:latest` = v0.61.0-Build):**
  Hauptprofil `--enable planning` + Fokus-Disables: **402 Dateien, 0 Befunde**,
  Exit 0. Closure-Profil `--config .d-check.closure.yml --enable planning
  --enable structure`: **371/0**, Exit 0. **Negativ-Probe** (Dateikopie nach
  Scratch, §Offene-Wellen-Zeiger durch `Nichts in Arbeit.` ersetzt bei
  beanspruchtem slice-108): Exit 1 mit **zwei** Befunden — `planning-drift`
  (roadmap.md:12) **und** `wave-drift` (Marker bei liegender welle-78-Datei).
  Restore per `cp` aus dem Scratch-Original, verifiziert: SHA-256 identisch
  (`63eca843…aef6`), `git status --porcelain` leer, Wiederholungs-Lauf 402/0
  grün. Alle drei Commit-Claims (402/0, 371/0, Negativ-Probe rot) reproduziert.

---

## Findings

### F-1 (LOW) — MR-024-Datei in `done/` ohne „Aufgelöst durch:"-Zeile; der Body liest sich weiter aktiv

- **kategorie:** LOW
- **quelle:** `MR-024` / Konventionsspeicher-Form (`### Aufgelöste`-Muster)
- **pfad:** `harness/conventions/done/MR-024-aktuelle-welle-ruhe-marker-form.md:3` und `:34–35`
- **befund:** Alle **15** anderen Dateien in `harness/conventions/done/` tragen
  eine „**Aufgelöst durch:**"-Zeile im Kopf (z. B. `MR-022`: „Baseline-Stand
  v5.0.0 …"); die bewegte MR-024-Datei nicht. Ihr Body sagt weiter
  `Status: Accepted` und „Auflösungs-Trigger: **permanent**, solange
  `make planning-check` den … Invariant erzwingt" — diese Bedingung gilt
  unverändert, d. h. wer nur die Datei liest (nicht den Index), schließt auf
  eine weiterhin bindende Adaption; die tatsächliche Auflösungs-Begründung
  (Baseline v5.6.0 trägt Marker+Wächter selbst) steht nur in der Index-Zeile.
- **verifizierbar:** nein (kein Gate; Beleg: `grep -L "Aufgelöst durch" harness/conventions/done/*.md` trifft genau diese Datei)
- **klasse:** aufloesung-nur-im-index

### F-2 (LOW) — Produkt-Befundtext verdrahtet „§Aktuelle Welle" statt der konfigurierten Überschrift

- **kategorie:** LOW
- **quelle:** `DC-FA-PLAN-001` / Maintainability
- **pfad:** `internal/hexagon/core/rules/planning.go:42` und `:45`
- **befund:** Beide `planning-drift`-Meldungen nennen die Sektion literal
  („die Roadmap §Aktuelle Welle trägt den Ruhe-Marker …") statt
  `cfg.EffectiveHeading()`. In der Negativ-Probe (`--yaml`) lautet die Meldung
  im eigenen Repo: „… die Roadmap **§Aktuelle Welle** trägt den Ruhe-Marker
  „Nichts in Arbeit"" — eine Sektion, die es in d-checks Roadmap seit `a46f09c`
  nicht mehr gibt. Befund-Ort, Grund-Code und Exit sind korrekt; nur die
  Diagnose-Prosa führt den Reparierenden in die falsche Sektion. d-check ist
  mit C-1 der erste Nicht-Default-`heading`-Konsument; der Commit-Claim „das
  Produkt deckt beides per Config" gilt für die Prüfung, nicht für den
  Meldungs-Text. Nicht von `a46f09c` eingeführt (Plan-Grenze „kein
  Produkt-Code" eingehalten) — aber von ihm erstmals scharf gestellt.
- **verifizierbar:** ja (Negativ-Probe mit `--yaml`, oben reproduziert)
- **klasse:** befundtext-default-statt-config

### F-3 (LOW) — Kleingeschriebene semantische Spiegel: drei lebende Stellen beschreiben die Roadmap weiter als Träger der „aktuellen Welle"

- **kategorie:** LOW
- **quelle:** Maintainability (Spiegel-Vollständigkeit; MR-025-Klasse)
- **pfad:** `AGENTS.md:57` · `harness/README.md:25` · `docs/plan/planning/README.md:18`
- **befund:** Die Source-Precedence-Zeile Rang 5 („`roadmap.md` — aktuelle
  Welle") in beiden Dateien und der Planning-README-Satz („Die aktuelle Welle
  und die Wellen-Reihenfolge stehen in `in-progress/roadmap.md`") beschreiben
  die Rolle, die die adoptierte Form der Roadmap gerade **entzieht**: derivativ
  sagt das `Welle:`-Feld der Slices, woran gearbeitet wird; die Roadmap listet
  offene Wellen. Der Spiegel-Sweep des Commits (grep „Aktuelle Welle"/„Keine
  aktive Welle") ist case-sensitiv und konnte diese Stellen nicht treffen.
- **verifizierbar:** nein (Doku-Prosa; Beleg: `grep -in "aktuelle Welle"` über die drei Dateien)
- **klasse:** token-grep-verfehlt-semantik-spiegel

### F-4 (LOW) — `.d-check.yml`-Blockkommentar behauptet weiter „heading/marker … per Default", direkt über den neuen Overrides

- **kategorie:** LOW
- **quelle:** Maintainability
- **pfad:** `.d-check.yml:247` (gegen `:257–258`)
- **befund:** Der Einleitungs-Kommentar des `planning`-Blocks sagt
  „heading/marker/slice-glob per Default (`## Aktuelle Welle` / „Keine aktive
  Welle" / `slice-*.md`)" — zehn Zeilen darüber, dass genau zwei der drei
  Schlüssel jetzt explizit überschrieben sind. Zwei Antworten auf dieselbe
  Frage in einem Bildschirm, in der Datei, die der Commit editiert hat.
- **verifizierbar:** nein (Kommentar; kein Gate liest ihn)
- **klasse:** kommentar-widerspricht-config

### F-5 (LOW) — Die wave-drift-Halbseite des Gates bleibt in beiden neuen Texten unbenannt

- **kategorie:** LOW
- **quelle:** `DC-FA-PLAN-001` (W3, ADR-0055) / `MR-013` / `AGENTS.md` §3.3
- **pfad:** `AGENTS.md:98–101` · `.d-check.yml:252–256`
- **befund:** Der baseline-legitime **Ein-Wellen**-Zustand „Welle offen +
  `in-progress/` leer" (modul-06: der Abschnitt trägt dann den Ruhe-Marker) ist
  unter d-checks W3-Kopplung maschinell rot — die eigene Negativ-Probe belegt
  es live (`wave-drift` feuerte **zusätzlich**, weil die flache
  welle-78-Datei lag). Der neue §3.3-Text („zurück auf den Ruhe-Marker …,
  sofern kein Slice mehr beansprucht ist") nennt nur die planning-drift-Hälfte:
  wer ihm beim **letzten** Slice einer Welle wörtlich folgt (Marker setzen,
  Wellendokument liegen lassen), produziert einen gate-roten Commit. Der
  Config-Kommentar benennt als Grenze nur den „Mehr-Wellen-Betrieb"; die
  nähere Grenze ist dieser Ein-Wellen-Pausen-Zustand, den die bisherige Praxis
  (Wellen-Closure samt `git mv` im selben Commit, kombinierte Move+Eröffnung)
  stillschweigend vermeidet. Kein Verhaltens-Regress — `a46f09c` ändert nur
  Config-Strings, W3 rechnet identisch; die Lücke ist die Beschreibung.
- **verifizierbar:** ja (Negativ-Probe oben: zweiter Befund `wave-drift`)
- **klasse:** gate-beschreibung-nennt-eine-von-zwei-kopplungen

### F-6 (LOW) — Roadmap §Abhängigkeitsgraph ist stehen geblieben (Vor-Commit-Drift, live)

- **kategorie:** LOW
- **quelle:** Maintainability (Roadmap-Currency)
- **pfad:** `docs/plan/planning/in-progress/roadmap.md:53–61`
- **befund:** Der Mermaid-Graph führt „slice-107 - Etappe B (in Arbeit)" und
  kennt die C-Etappen nicht, während slice-107 in `done/` liegt und slice-108
  beansprucht ist. Eingeführt vom Eröffnungs-Commit `5a41cd3` (die
  Vorgänger-Übergaben `0b9431f`/`37e1fd4` haben den Graphen je Etappe
  nachgeführt; `5a41cd3` und `a46f09c` nicht) — kein `a46f09c`-Regress, aber
  live in der Datei, deren Form der Commit gerade auf Baseline hebt.
- **verifizierbar:** nein (kein Gate prüft Graph-Currency)
- **klasse:** doku-drift-graph

### F-7 (INFO) — Plan-Schritt 1 nennt einen Drift-Log-Eintrag; `a46f09c` trägt keinen

- **kategorie:** INFO
- **quelle:** slice-108-Plan §2 Schritt 1
- **pfad:** `docs/plan/planning/in-progress/slice-108-roadmap-offene-wellen.md:32–34`
- **befund:** Die Umstellung selbst hat keine eigene Zeile in §Historische
  Trigger-Verschiebungen; sie steht nur **prospektiv** am Ende der
  `5a41cd3`-Zeile („C-1 stellt die Roadmap auf die v5.6.0-Form … um").
  Plan-/DoD-Konformität ist Verifikations-Rolle — hier nur protokolliert,
  nicht gewertet.
- **verifizierbar:** nein
- **klasse:** plan-artefakt-ohne-traeger

### F-8 (INFO) — MR-024-Geltungsbereichs-Link: historischer Text zielt auf den neuen Anker

- **kategorie:** INFO
- **quelle:** Maintainability (eingefrorene Datei, Link-Integrität)
- **pfad:** `harness/conventions/done/MR-024-aktuelle-welle-ruhe-marker-form.md:6`
- **befund:** Der Link-Text „`roadmap.md` §Aktuelle Welle" zielt jetzt auf
  `#offene-wellen` — Text/Ziel-Spreizung in einer eingefrorenen Datei. Es ist
  der minimale Fix, der `doc-check` grün hält (der alte Anker existiert nicht
  mehr), im Commit-Text offengelegt; die Alternative (Tombstone) verlöre den
  Zeiger. Bewusste Won't-Fix-würdige Designnotiz.
- **verifizierbar:** ja (`make doc-check` bliebe bei `#aktuelle-welle` rot)
- **klasse:** frozen-link-retarget

## Negativbefunde (geprüft, ohne Befund)

- **Form gegen Baseline:** Die neue Sektion entspricht `modul-06-roadmap.md`
  §Roadmap-Struktur (v5.6.0) und der `roadmap.template.md`-Form exakt:
  H2 `## Offene Wellen` (Zeile 12), derivative Regeln-Paraphrase
  (Zustand = flache Welle-Dateien, `Welle:`-Feld, Ziel/Trigger/Closure in der
  Welle-Datei, Geplantes Ende = Schätzung), Zeiger-Bullet in Template-Syntax
  (`- [welle-78-…](../welle-78-….md)`), keine Struktur-Felder mehr.
- **Marker-Paraphrase ehrlich:** Der Sektions-Regeln-Text trägt den Wortlaut
  „Nichts in Arbeit" **nicht** (grep über den Block: kein Treffer; die
  einzigen Roadmap-Vorkommen liegen im Drift-Log außerhalb des
  `heading`-Blocks). Die Paraphrase ist zudem **notwendig**, nicht nur Kür:
  das Baseline-Template selbst führt den Marker literal im
  BEDIENHINWEIS-HTML-Kommentar — eine wörtliche Template-Übernahme hätte den
  Substring-Match ausgelöst und die Sektion bei beanspruchtem Slice als ruhend
  gelesen. Der grüne 402/0-Lauf bei liegendem slice-108 belegt die aktive Lesart.
- **Beide Profile, eine Antwort:** `.d-check.yml` und `.d-check.closure.yml`
  tragen identische Werte (`## Offene Wellen` / `Nichts in Arbeit`). Die
  Schlüssel existieren (`configyaml.go` `rawPlanning` `yaml:"heading"`/
  `yaml:"marker"`, durchgereicht in `applyPlanning`; Spez §2-Tabelle) und
  **wirken** in beiden Profilen: das Closure-Grün 371/0 wäre bei ignoriertem
  `heading` unmöglich (Default-Überschrift fehlt in der Roadmap ⇒ fail-closed).
- **Wächter-Beleg:** eigenständig reproduziert (Messläufe im Kopf) inkl.
  byte-verifiziertem Restore und Grün-Wiederholung.
- **Spiegel-Sweep (Groß-Schreibung):** Alle verbliebenen „Aktuelle Welle"/
  „Keine aktive Welle"-Fundstellen klassifiziert legitim — Produkt-Defaults
  (Spec §2/`DC-FA-PLAN-001`, Benutzerhandbuch, README.de/README-Modullisten,
  `print-config`-Template, Go-Code/-Tests: der Produkt-Default bleibt per
  Plan-Schritt 5 unverändert), Historie (`done/`-Slices, `docs/reviews/`,
  CHANGELOG, Drift-Log-Zeilen, Handbuch-§11), vendorte Baseline, eingefrorene
  MR-Dateien in `conventions/done/`. `ADR-0028` (immutable) beschreibt die
  Konventions-**Defaults** des Moduls — mit Nicht-Default-Config kein
  Live-Widerspruch. `.harness/skills/` ohne Fundstelle;
  `.harness/cache/v1.4.0/` ist untracked Host-Rest, kein Repo-Bestand.
  Die fünf im Commit genannten Spiegel (AGENTS §3.3 + §4, harness/README-
  Sensors-Zeile, Makefile-Kommentar, MR-013) sind alle korrekt nachgezogen.
- **MR-024-Auflösung (Substanz):** tragfähig — modul-06 v5.6.0 trägt
  Offene-Wellen-Form **und** Ruhe-Marker **mit Wächter** („Ein Doku-Sensor hält
  den Marker gegen das Verzeichnis") selbst; die Adaption ist damit
  Baseline-Default, die `Ersetzt-Baseline-Regel`-Grundlage entfällt.
  Index-Zeile korrekt aus `### Aktive` entfernt und in `### Aufgelöste`
  eingefügt (numerisch/chronologisch konsistent hinter MR-022), Voll-Slug-
  `<a id>` erhalten — die historischen `conventions.md#mr-024-…`-Links
  (inkl. immutable ADRs und des welle-78-Wellendokuments) lösen weiter auf
  (`doc-check` 402/0). Link-Tiefen der bewegten Datei stimmen
  (`../../../` zu Repo-Wurzel-Zielen, `../../conventions.md`); Rest siehe F-1/F-8.
- **Move-Hygiene:** Rename-Detection hielt (similarity 76 % > 50 %,
  `git log --follow` intakt); die Bündelung Move + erzwungene Link-Fixes folgt
  derselben Begründung wie die MR-013-Lifecycle-Ausnahme.
- **wave-drift-Kopplung (Verhalten):** Kein Zustand, den die alte Form
  maschinell erlaubte und die neue rot macht oder umgekehrt — `a46f09c` ändert
  ausschließlich Config-Strings, `CheckPlanningWaves` rechnet identisch; die
  W3-Semantik-Divergenz zur Baseline-Prosa ist F-5 (Beschreibungs-, keine
  Verhaltens-Lücke).
- **AGENTS §3.3-Konditional:** „sofern kein Slice mehr beansprucht ist" ist
  gegen die kombinierte Move+Eröffnungs-Praxis **korrekt** (der alte
  unbedingte Flip war es nicht — bei nicht-leerendem `in-progress/` wäre er
  falsch); die Neuformulierung ist eine Verbesserung. Rest siehe F-5.
- **Kein Produkt-Code berührt:** der Diff fasst ausschließlich Config, Doku,
  Makefile-Kommentar und Konventionsspeicher an — plan-konform.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 6 (F-1…F-6) |
| INFO | 2 (F-7, F-8) |

## Verdikt

**APPROVE.** Der Kern der Etappe C-1 hält: die Sektion ist template-treu, die
Marker-Paraphrase ist ehrlich (und mechanisch notwendig), beide Profile geben
eine Antwort, der Wächter ist in beide Richtungen belegt (inkl. selbst
wiederholter Negativ-Probe mit verifiziertem Restore), und die
MR-024-Auflösung ist in der Substanz korrekt. Die sechs LOWs sind Nachzüge
ohne Gate-Wirkung; am ehesten vor dem Push wert sind F-1 (die `done/`-Datei
widerspricht als einzige dem 15/15-Auflösungs-Muster und liest sich allein
gelesen als aktiv) und F-3 (drei lebende Rollen-Beschreibungen der Roadmap
tragen noch die alte Semantik).
