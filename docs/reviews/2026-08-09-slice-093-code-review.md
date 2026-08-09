# Review-Report: slice-093 — Go-Implementierung Closure-Note-Gate + `--config` — 2026-08-09

- **Review-Art:** **Code**-Review (Produkt-Code gegen Plan, Lastenheft/Spezifikation
  und Konventionen). Die Vertrags-**Prosa** prüft ein zweiter Reviewer parallel und
  ist hier ausdrücklich **nicht** Gegenstand.
- **Gegenstand:** slice-093, Diff `18489ee..HEAD` (`ae4cc09`) beschränkt auf
  `internal/`, dazu die neue Datei `.d-check.closure.yml` und das neue
  `Makefile`-Target `verify-closure-notes` (`Makefile:186`).
- **Skill:** `.harness/skills/reviewer.md` @ Version 1.3.0.
- **Modell-ID:** claude-opus-5[1m].
- **Datum:** 2026-08-09.
- **Eingangs-Kontext (Verträge, gegen die geprüft wurde):**
  - `spec/lastenheft.md` §DC-FA-PLAN-001 (inkl. §Closure-Note-Struktur, neun
    Akzeptanzkriterien) und §DC-FA-CLI-012 (fünf Akzeptanzkriterien).
  - `spec/spezifikation.md` §DC-FA-PLAN-001.a Schritte C1–C5,
    §DC-FA-CLI-012.a Schritte 1–6, §2-Schema-Zeilen `planning.closure.dir` /
    `heading-pattern` / `min-sentences` / `boilerplate`, §4-Grund-Code-Tabelle.
  - `spec/architecture.md` §2 (Schichten/Constraints, Hexagon-Schnitt) und §5
    (Fehlermodelle, fail-closed).
  - `AGENTS.md` §3 (Hard Rules) und §4 (Gate-Tabelle).
  - `docs/plan/adr/0048-closure-note-struktur-im-planning-modul.md`,
    `docs/plan/adr/0005-modul-layout-hexagon-ordner.md`,
    `docs/plan/adr/0012-kern-paketschnitt-model-rules-app.md`.
  - Slice-Plan `docs/plan/planning/in-progress/slice-093-closure-note-gate.md`.

---

## Verifikation (selbst ausgeführt)

- `make lint` → `0 issues.`
- `make test` → alle Pakete `ok`.
- `make arch-check` (a-check, digest-gepinnt) → `gesamt: 0 Befund(e)`.
- `make gate-consistency` → `328 Datei(en) geprüft, 0 Befund(e)`.
- `make verify-closure-notes` → `299 Datei(en) geprüft, 0 Befund(e)`.
- `make gates` → grün, `coverage-gate: OK — Coverage 94.30% erfüllt Schwelle 93%`.
- **Synthetische Läufe** über `docker run --network none -v <lab>:/repo:ro d-check:latest`
  gegen ein Wegwerf-Repo (Roadmap + `done/`-Bestand + zwei Profile) — die belegten
  Beobachtungen unten stammen aus diesen Läufen.
- **Mutations-Gegenproben:** drei gezielte Rückbauten im Arbeitsbaum, jeweils
  `make test`, danach `git checkout --` (Arbeitsbaum wieder sauber). Ergebnis je
  Mutation unter F-5, F-6 und F-7.
- **Bestands-Scan** des `done/`-Korpus (92 Slices) auf die unter F-1 beschriebene
  Eingabe-Gestalt: 0 aktuelle Treffer.

## Findings

### F-1 — Abschnitts-Grenze endet an jeder `#`-Zeile, nicht an jeder Überschrift

- **kategorie:** HIGH
- **quelle:** DC-FA-PLAN-001 §Closure-Note-Struktur / Spezifikation
  §DC-FA-PLAN-001.a Schritt C3 („bis zur nächsten **Überschrift** gleicher oder
  höherer Ebene")
- **pfad:** `internal/hexagon/core/rules/planning.go:183` (Abschnitts-Ende),
  `internal/hexagon/core/rules/planning.go:157` (Kandidaten-Filter),
  `internal/hexagon/core/rules/planning.go:193` (`headingLevel`)
- **befund:** Die Abschnitts-Grenze wird über `strings.HasPrefix(trimmed, "#")` plus
  reine `#`-Zählung bestimmt; die ATX-Gültigkeit (Leerzeichen/Tab nach der
  `#`-Folge, Ebene ≤ 6) wird nicht geprüft, obwohl dasselbe Paket dafür bereits
  `parseATXHeading` in `internal/hexagon/core/rules/markdown.go:286` führt, das
  `anchors` und `matrix` nutzen. Eine gewöhnliche Zeile wie `#1 war ein Thema` oder
  ein `#`-Kommentar in einem eingerückten (nicht gefencten) Code-Absatz gilt damit
  als H1 und beendet den Closure-Abschnitt vorzeitig; alles danach fällt aus der
  Messung. Belegt im Wegwerf-Repo: eine Notiz mit vier Satzende-Zeichen, danach der
  Zeile `#1 war ein Thema` und danach der deklarierten Floskel „alles gut" liefert
  **0 Befunde**, während dieselbe Notiz ohne die `#`-Zeile
  `closure-note-boilerplate` liefert. Für `closure-note-thin` wirkt derselbe
  Mechanismus in die Gegenrichtung (Falschbefund an einer inhaltlich vollen Notiz).
  Der Bestands-Scan über `docs/plan/planning/done` zeigt aktuell 0 Treffer — der
  Pfad ist vorhanden, aber (noch) nicht ausgelöst.
- **verifizierbar:** ja — `make verify-closure-notes` gegen einen `done/`-Bestand,
  der eine solche Zeile zwischen Prosa und Floskel trägt: Exit 0 statt
  `closure-note-boilerplate`; die Gegenprobe ohne die Zeile ist Exit 1.
- **klasse:** Silent-Grün durch eigenen Parser statt geteiltem Scanner

### F-2 — `--config` mit leerem Wert fällt still auf die konventionelle Datei zurück

- **kategorie:** MEDIUM
- **quelle:** DC-FA-CLI-012 („**kein** stiller Rückfall auf Defaults und **keinen**
  auf die konventionelle Datei"), Spezifikation §DC-FA-CLI-012.a Schritt 3
- **pfad:** `internal/adapter/driving/cli/cli.go:369`
- **befund:** Die Fallunterscheidung ist `case configPath != ""`; ein explizit
  übergebener, aber leerer Wert (`--config ""`, in der Praxis eine nicht gesetzte
  Make-/Shell-Variable) ist damit nicht von „Flag abwesend" unterscheidbar und der
  Lauf liest die konventionelle Datei. Belegt: im selben Wegwerf-Repo liefert
  `--config profil.yml --enable planning` Exit 1 mit `closure-note-thin`, während
  `--config "" --enable planning` Exit 0 liefert — dasselbe Repo, dasselbe Modul,
  grün statt rot, ohne jede Meldung auf stderr.
- **verifizierbar:** ja — zwei Läufe desselben Images gegen dasselbe Repo, einmal
  mit gesetztem und einmal mit leerem `--config`-Wert; die Exit-Codes
  unterscheiden sich, ohne dass eine Diagnose ausgegeben wird.
- **klasse:** Leerer Flag-Wert = Flag abwesend (stiller Rückfall)

### F-3 — Config-Fehlermeldungen nennen `.d-check.yml`, auch wenn eine andere Datei geladen wurde

- **kategorie:** MEDIUM
- **quelle:** Spezifikation §DC-FA-CLI-012.a Schritt 6 („benennt er die
  **tatsächlich geladene** Datei — nicht den konventionellen Namen (sonst zeigte
  die Meldung auf eine Datei, die der Lauf nie gelesen hat)")
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go:25` (`FileName`),
  Verwendung in 51 Fehler-Formaten, u. a. `configyaml.go:936` und
  `configyaml.go:951` in der neu hinzugefügten `applyClosure`
- **befund:** Die Herkunfts-Nachführung wurde für die `sources`-Befunde umgesetzt
  (`model.Config.ConfigFile`), nicht aber für die Validierungs-Diagnosen des
  Config-Adapters: diese präfixen weiterhin den konstanten Konventionsnamen.
  Belegt: `--config kaputt.yml` mit `min-sentences: 0` meldet
  `d-check: error: .d-check.yml: planning.closure.min-sentences 0 muss >= 1 sein`;
  ein YAML-Syntaxfehler in derselben Datei meldet
  `d-check: error: .d-check.yml: yaml: line 1: …`. In genau dem Repo, für das der
  Schalter gebaut wurde, existieren beide Dateien nebeneinander — die Meldung zeigt
  auf die falsche.
- **verifizierbar:** ja — `make verify-closure-notes` nach einem eingebauten
  Schema-Fehler in `.d-check.closure.yml`: stderr nennt `.d-check.yml`.
- **klasse:** Hart verdrahteter Konventionsname in der Fehler-Provenance

### F-4 — Wurzel-Constraint des `--config`-Pfads greift nicht durch ein Verzeichnis-Symlink

- **kategorie:** MEDIUM
- **quelle:** DC-FA-CLI-012 („muss **innerhalb** der Scan-Wurzel liegen"),
  Spezifikation §DC-FA-CLI-012.a Schritt 2 (fail-closed)
- **pfad:** `internal/adapter/driving/cli/cli.go:430` (`normalizeConfigPath`),
  `internal/adapter/driving/cli/cli.go:403` (`readConfigOverride`)
- **befund:** Die Constraint ist rein lexikalisch (`path.Clean` plus
  `../`-Präfixtest); die anschließende Klassifikation `Kind` liegt per `Lstat` nur
  auf der **letzten** Pfad-Komponente. Ein Symlink als Konfigurationsdatei wird
  dadurch korrekt abgelehnt, ein Symlink als **Zwischenverzeichnis** nicht: mit
  einem Verzeichnis-Symlink `esc` in der Wurzel liest der Lauf unter
  `--config esc/etc/passwd` nachweislich `/etc/passwd` außerhalb des Mounts (die
  Meldung zitiert dessen Inhalt in der YAML-Fehlermeldung). Der Lauf hängt damit an
  Eingaben außerhalb der Scan-Wurzel, entgegen der als fail-closed zugesagten
  Hermetik. Dieselbe Lstat-Mechanik gilt auch für `scan.roots` und
  `planning.roadmap`; neu ist, dass `--config` die Constraint als eigenen
  fail-closed Schritt zusagt.
- **verifizierbar:** ja — Lauf im read-only-Mount mit einem Verzeichnis-Symlink in
  der Wurzel und einer gültigen YAML-Datei hinter ihm: die Konfiguration wird
  geladen statt mit Exit 2 abgelehnt.
- **klasse:** Wurzel-Constraint nur auf der letzten Pfad-Komponente

### F-5 — Die neue Befund-Provenance in `sources` ist von keinem Test gehalten

- **kategorie:** MEDIUM
- **quelle:** Spezifikation §DC-FA-CLI-012.a Schritt 6; Reviewer-Skill §MEDIUM
  („fehlende Negativtests bei neuem öffentlichen Vertrag")
- **pfad:** `internal/hexagon/core/rules/sources.go:217`,
  `internal/hexagon/core/rules/sources_test.go:382` und `:406`
- **befund:** Die Signatur-Erweiterung von `CheckSources`/`allSourceRefs` ist der
  einzige Zweck der Änderung an `sources.go`, aber kein Test übergibt jemals einen
  nicht-leeren `configFile`; die beiden angepassten Aufrufe reichen `""` durch, und
  `model.Config.ConfigFile` kommt in keinem `_test.go` des Repos vor.
  Mutations-Gegenprobe: der Rückbau von `allSourceRefs` auf den konstanten
  Konventionsnamen (Wegfall des Durchreichens) lässt `make test` vollständig grün —
  der Vertragszusatz ist mechanisch unbelegt.
- **verifizierbar:** ja — Mutation wie beschrieben, danach `make test`: grün.
- **klasse:** Neuer Vertragszusatz ohne haltenden Test

### F-6 — `TestConfigPath_AusserhalbDerWurzel` belegt die Wurzel-Constraint nicht

- **kategorie:** MEDIUM
- **quelle:** DC-FA-CLI-012 §Akzeptanzkriterien, „Negative (außerhalb der Wurzel)"
- **pfad:** `internal/adapter/driving/cli/cli_config_path_test.go:76`
- **befund:** Beide Eingaben des Tests erreichen Exit 2 auch **ohne** die
  Wurzel-Constraint: `../raus.yml` existiert nicht (Existenz-Prüfung greift) und
  `../../etc/passwd` ist kein gültiges YAML (Decoder greift). Mutations-Gegenprobe:
  entfernt man den `..`-Test in `normalizeConfigPath` ersatzlos, bleibt das
  CLI-Test-Paket grün. Das Akzeptanzkriterium ist damit formal abgedeckt, aber
  nicht gehalten — genau die Lücke, hinter der F-4 sitzt.
- **verifizierbar:** ja — Mutation wie beschrieben, danach `make test`: grün.
- **klasse:** Negativtest besteht aus dem falschen Grund

### F-7 — Der Determinismus-Test der Closure-Kandidaten ist tautologisch

- **kategorie:** LOW
- **quelle:** DC-QA-02; Reviewer-Skill §LOW („latente Wartungsfalle")
- **pfad:** `internal/hexagon/core/rules/planning_closure_test.go:199`,
  geprüfte Stelle `internal/hexagon/core/rules/planning.go:96`
- **befund:** Der Test vergleicht Wiederholungsläufe gegen dasselbe `MemFS`, dessen
  `List` bereits nach Namen sortiert (`internal/hexagon/core/coretest`), ebenso wie
  der reale Adapter und der abschließende `model.SortFindings` in
  `internal/hexagon/core/rules/run.go:112`. Mutations-Gegenprobe: Entfernen von
  `sort.Strings(names)` samt Import lässt `make test` vollständig grün. Der Test
  misst die Sortierung des Doubles, nicht die der Regel; zusätzlich indiziert die
  Schleife `again` über die Länge von `first`, sodass ein kürzeres Ergebnis in einen
  Index-Panic statt in ein Test-Versagen liefe.
- **verifizierbar:** ja — Mutation wie beschrieben, danach `make test`: grün.
- **klasse:** Determinismus-Test tautologisch (Double garantiert die Ordnung)

### F-8 — Die vier neuen Konfigurations-Schlüssel fehlen im `--print-config`-Gerüst

- **kategorie:** LOW
- **quelle:** DC-FA-CLI-005 („es dokumentiert die verfügbaren Module und Optionen
  als Kommentare, damit sie sichtbar sind")
- **pfad:** `internal/adapter/driving/cli/config_template.go:121` bis `:125`
- **befund:** Der `planning`-Block des Startgerüsts führt `roadmap`, `heading`,
  `marker` und `slice-glob`, aber keinen `closure`-Unterblock; `dir`,
  `heading-pattern`, `min-sentences` und `boilerplate` sind über
  `d-check --print-config` nicht auffindbar, obwohl das Gerüst für jedes andere
  Modul auch die optionalen Schlüssel auskommentiert mitführt (etwa der
  `trace`-Baum ab `config_template.go:165`). Ein Adopter, der das Gerüst als
  Optionsübersicht nutzt, sieht die zweite `planning`-Fähigkeit nicht.
- **verifizierbar:** ja — `docker run … d-check:latest --print-config` und Suche
  nach `closure`: kein Treffer.
- **klasse:** Neue Config-Schlüssel fehlen im `--print-config`-Gerüst

### INFO-1 — Die Substanz-Schwelle ist ohne Prosa erreichbar

- **kategorie:** INFO
- **quelle:** DC-FA-PLAN-001 §Out-of-Scope („Struktur, nicht Bedeutung")
- **pfad:** `internal/hexagon/core/rules/planning.go:205` (`countSentenceEnds`)
- **befund:** Gezählt werden `.`, `!` und `?` im gesamten bereinigten Abschnitt,
  also auch die Punkte in Datei-Endungen, Link-Zielen, Versionsnummern und
  Abkürzungen. Ein Closure-Abschnitt, der nur aus vier Verweiszeilen der Form
  `- [ADR-0048](../adr/0048-….md)` besteht, überschreitet die Default-Schwelle ohne
  einen Satz Prosa. Das ist vertragskonform (die Bedeutung ist ausdrücklich nicht
  zugesagt) und im Code als bewusste Vereinfachung kommentiert; als **operative**
  Decke der Messung ist es an keiner Betriebsstelle notiert.
- **verifizierbar:** ja — synthetischer Slice mit vier reinen Verweiszeilen:
  `make verify-closure-notes` grün.
- **klasse:** Substanz-Schwelle über Zeichenzählung ohne Prosa erreichbar

## Negativbefunde

- geprüft, ohne Befund: **Byte-Identität ohne Aktivierung (C1)** —
  `CheckPlanningClosure` verlässt vor jedem `List`/`ReadFile`, wenn
  `cfg.Closure.Dir` leer ist; die Aktiv-Status-Prüfung liest weiterhin nur die
  Roadmap-Datei und das Listing ihres Verzeichnisses, nie Slice-Inhalt. Der einzige
  weitere neue Kern-Zustand (`model.Config.ConfigFile`) wird ausschließlich von
  `sources` gelesen und trägt ohne `--config` denselben Wert wie die bisherige
  Konstante — der Befundsatz ist unverändert.
- geprüft, ohne Befund: **fail-closed am Config-Rand** — `min-sentences` 0 und −1,
  leerer bzw. nur aus Leerzeichen bestehender `boilerplate`-Eintrag, nicht
  kompilierendes `heading-pattern`, absoluter oder `..`-haltiger `dir`: alle
  end-to-end mit Exit 2 belegt. Die Zeiger-Modellierung von `min-sentences` hält
  „abwesend" und „explizit 0" korrekt auseinander. Die Guard-Form entspricht dem
  Haus-Muster (fünf gleichartige Fundstellen in `configyaml.go`).
- geprüft, ohne Befund: **fail-closed zur Laufzeit** — gesetztes, aber fehlendes
  Verzeichnis (`closure-note-missing` auf dem Verzeichnis, Exit 1), unlesbare bzw.
  kaputt verlinkte Slice-Datei (`closure-note-missing`, Exit 1),
  `closure.dir` auf eine Datei statt ein Verzeichnis (Listing-Fehler ⇒ Befund).
  Ein **leeres** Verzeichnis ist wie zugesagt befundfrei.
- geprüft, ohne Befund: **Fence-Behandlung** — der geteilte `FenceToggle` wird
  genutzt (nicht nachgebaut), damit gelten die CommonMark-Infozeilen-Regel und die
  `~~~`-Asymmetrie identisch zum übrigen Scanner. Belegt: Satzende-Zeichen in einem
  Backtick- und in einem `~~~`-Fence-Block zählen nicht; eine `#`-Zeile **im** Fence
  eröffnet keinen Abschnitt; eine unbalancierte Fence schluckt den Rest und führt
  zu `closure-note-thin`, also fail-loud.
- geprüft, ohne Befund: **Zeilen-Randfälle** — Überschrift in Zeile 1, Abschnitt bis
  Dateiende, Datei ohne abschließenden Zeilenumbruch und CRLF-Zeilenenden liefern
  jeweils das erwartete Ergebnis (`splitLines` plus `TrimSpace` neutralisieren das
  `\r`).
- geprüft, ohne Befund: **Abschnitts-Ebenen-Semantik** — eine tiefere Überschrift
  gehört zum Abschnitt, eine gleich- oder höherrangige beendet ihn; C5 wird
  eingehalten (`file`, `line` = Überschriften-Zeile, `rule` = `planning`,
  `target` = das Closure-Verzeichnis), `closure-note-missing` schließt die beiden
  Mess-Codes aus, `thin` und `boilerplate` treten nebeneinander auf.
- geprüft, ohne Befund: **Determinismus des Gesamtlaufs (DC-QA-02)** — Kandidaten
  kommen aus einem sortierten Listing, und der Lauf schließt mit
  `model.SortFindings`. Die Aussage der Regel selbst ist allerdings ungetestet
  (F-7).
- geprüft, ohne Befund: **Hexagon-Import-Regeln** — `rules/planning.go` nimmt nur
  `regexp`, `sort`, `strconv`, `strings`, `path` aus der Standardbibliothek plus
  `model` und `port/driven`; `app/diagnose.go` bleibt auf `model` + `rules`;
  `model` bekommt nur Datentypen und reine Methoden. Kein Adapter-Import im Kern,
  keine Netz-/Prozess-API. `make arch-check` bestätigt es.
- geprüft, ohne Befund: **Grund-Code-Lockstep** — die drei neuen Codes stehen in
  `AllReasons`, in `reasonTexts` und in der §4-Tabelle der Spezifikation; der
  beidseitige Lockstep-Test läuft grün.
- geprüft, ohne Befund: **`--config`-Pfadformen** — `./`-Präfix, absoluter Pfad
  innerhalb der Wurzel, Backslash-Trenner (unter Linux literal, daher nicht
  gefunden ⇒ Exit 2), Verzeichnis statt Datei (Exit 2), Symlink **als** Datei
  (Exit 2 über die `Lstat`-Klassifikation), ungültiges Schema (Exit 2). Nur der
  Zwischenverzeichnis-Symlink fällt durch (F-4).
- geprüft, ohne Befund: **Argument-Reihenfolge** — `-config`/`--config` sind im
  `valueFlags`-Satz von `reorderArgs` eingetragen, das Flag verschluckt also kein
  Pfad-Argument.
- geprüft, ohne Befund: **DC-QA-03 der zweiten Konfiguration** — der neue Go-Test
  verriegelt, dass `.d-check.closure.yml` kein Netz-/Range-Modul aktiviert **und**
  dass `planning.closure.dir` gesetzt ist (ein Profil ohne Schlüssel wäre ein leer
  laufendes Gate). Das ist die passende Invariante für ein Profil, das den
  Netzlos-Modulsatz bewusst nicht führt.
- geprüft, ohne Befund: **Bindepunkt-Verdrahtung** — `verify-closure-notes` hängt
  in `fullbuild`, nicht in `gates`/`ci`, steht in `.PHONY` und in der
  Gate-Tabelle; `make gate-consistency` bestätigt die Doku-↔-Makefile-Konsistenz.
- geprüft, ohne Befund: **Coverage der neuen Pfade** — `CheckPlanningClosure` 94,4 %,
  `checkClosureNote` 94,7 %, `loadConfig` 90,5 %, `readConfigOverride` 84,6 %,
  `normalizeConfigPath` 90,0 %; unbedeckt bleiben ausschließlich die defensiven
  Zweige (nicht erreichbarer `regexp.Compile`-Fehler, `filepath.Rel`-Fehler).
- geprüft, ohne Befund: **Suppression-Verbot und Gate-Lockerung (Hard Rules §3.2,
  §3.6)** — keine `nolint`-Direktive, keine Schwellen-Änderung, keine neue
  `.golangci.yml`-Ausnahme im Diff.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 5 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Silent-Grün durch eigenen Parser statt geteiltem
Scanner · Leerer Flag-Wert = Flag abwesend (stiller Rückfall) · Hart verdrahteter
Konventionsname in der Fehler-Provenance · Wurzel-Constraint nur auf der letzten
Pfad-Komponente · Neuer Vertragszusatz ohne haltenden Test · Negativtest besteht
aus dem falschen Grund · Determinismus-Test tautologisch (Double garantiert die
Ordnung) · Neue Config-Schlüssel fehlen im `--print-config`-Gerüst ·
Substanz-Schwelle über Zeichenzählung ohne Prosa erreichbar

## Verdikt

**Merge-blockierend: ja.** Ein HIGH und fünf MEDIUM; nach der Skill-Regel blockieren
beide Stufen typischerweise, und hier gibt es keinen Grund für eine Abweichung: F-1
und F-2 sind zwei unabhängige Wege, auf denen dieses Gate grün meldet, obwohl es
hätte melden müssen — genau die Klasse, deren Vermeidung der Slice zusagt. F-3 und
F-4 sind Vertrags-Zusagen von DC-FA-CLI-012, die der Code nicht einlöst; F-5 und
F-6 sind die Testlücken, hinter denen F-3 bzw. F-4 unentdeckt bleiben konnten.

Bemerkenswert positiv und ausdrücklich nicht relativiert: die fail-closed Ränder
sind vollständig und end-to-end belegt, die Zeiger-Modellierung von `min-sentences`
schließt die Null-Schwelle sauber aus, der geteilte `FenceToggle` wurde
wiederverwendet statt nachgebaut, und die Inertheit ohne `planning.closure.dir` ist
strukturell garantiert, nicht nur behauptet. Die Substanz des Slices trägt; die
Findings betreffen Ränder, nicht den Entwurf.

**Übergabe:** Findings gehen an den Implementer. Die Finding-Klassen gehören in die
Slice-Closure §7 und von dort in den Beobachtungs-Zähler; die Klasse „Silent-Grün
durch eigenen Parser statt geteiltem Scanner" ist ein Wiederholungskandidat, weil
das Repo für Markdown-Lexik bereits geteilte Bausteine führt. Dieser Report ist ein
Lauf-Beleg (dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) und ersetzt
keine Verifikation — DoD- und Plan-Konformität prüft der Verifier separat.
