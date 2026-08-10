# Review-Report: slice-097 — 2026-08-10

**Review-Art:** Code — geprüft wird die Umsetzung gegen Slice-Plan, Anforderung,
Spezifikation und die begleitende ADR (Modul 10 §Drei Review-Arten). Kein
Verifikations-Lauf: die DoD-Abhakung bleibt ausdrücklich außen vor.

**Gegenstand:** [slice-097](../plan/planning/in-progress/slice-097-closure-glob-entkopplung.md),
Diff-Range `39a3c6a..HEAD` (drei Commits: Wellen-Eröffnung `39a3c6a`, CR + ADR
`70d4c53`, Implementierung `21de5a3`).

**Skill:** `.harness/skills/reviewer.md` @ 1.3.0 ·
**Modell:** claude-opus-5[1m] · **Datum:** 2026-08-10

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne diese Liste
ist der Lauf nicht reproduzierbar):

- [slice-097](../plan/planning/in-progress/slice-097-closure-glob-entkopplung.md)
  §3 Abnahme-Punkte, §4 DoD, §5 Risiken
- [ADR-0051](../plan/adr/0051-eigener-kandidaten-filter-closure.md)
  (Entscheidungen 1–6, Fitness Function, Re-Evaluierungs-Trigger) und
  [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md)
  (Entscheidungen 1, 2, 8; §Geschichte)
- [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  (Lastenheft-Version 0.53.0) und
  [§`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritt C2 samt §2-Schema-Zeile
- [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) (Byte-Identität)
- [`AGENTS.md`](../../AGENTS.md) (Hard Rules)

**Läufe dieses Reviews.** Alle Fixtures liegen in einem Temp-Verzeichnis
außerhalb des Repos, alle Läufe netzlos und read-only. Gebaut wurden **zwei**
Vergleichs-Images aus dem Repo selbst: der Stand **vor** der Änderung (`39a3c6a`,
über ein `git archive`-Auspack ins Temp-Verzeichnis) und HEAD. Gefahren:
`make build`, `make test` (fünfmal, davon viermal mutiert), `make gates`,
`make verify-closure-notes`, `make adr-check` sowie rund 40 Fixture-Läufe gegen
die beiden Images. `make gates` ist grün (355 Dateien, 0 Befunde),
`make verify-closure-notes` grün (326 Dateien, 0 Befunde),
`make adr-check RANGE=39a3c6a..HEAD` grün. Der Arbeitsbaum ist am Ende
unverändert (`git status --short` leer).

---

## Findings

### F-1 — Bei elf der fünfzehn neu aufgenommenen Dokumente fällt die Abschnitts-Grenze mit der Dokument-Grenze zusammen

- `kategorie`: MEDIUM
- `quelle`: [ADR-0051](../plan/adr/0051-eigener-kandidaten-filter-closure.md)
  Entscheidung 5 und §Fitness Function („Der eigene Bestand bleibt bei null,
  jetzt über 111 statt 96 Dokumente")
- `pfad`: `.d-check.closure.yml:36` und `.d-check.closure.yml:43`
- `befund`: In den elf `welle-*-results.md` ist die auf `^#{1,3} .*Closure-Notiz`
  passende Überschrift der **Dokument-Titel** (H1); der gemessene Abschnitt ist
  damit die gesamte Datei. Gemessen an einer Fixture nach genau diesem Muster:
  ein Ergebnisdokument, dessen drei Retrospektiv-Abschnitte sämtlich
  `_Ausstehend._` tragen, bleibt **befundfrei**, während derselbe Platzhalter
  unter einer H2 `closure-note-thin` meldet. In der Gegenrichtung greift die
  Floskel-Prüfung jetzt auf beliebigen Dateitext: das Zitat „war ganz okay" in
  einem Nicht-Notiz-Abschnitt desselben Dokumententyps meldet
  `closure-note-boilerplate`. Für diese elf Dokumente ist die Substanz-Schwelle
  damit nicht mehr an die Notiz gebunden, sondern an die Dateilänge.
- `verifizierbar`: ja — `make verify-closure-notes` mit einer nach dem Muster der
  `welle-*-results.md` gebauten Datei im Closure-Verzeichnis (H1-Titel mit
  „Closure-Notiz", Platzhalter in den Unterabschnitten): Exit 0.
- `klasse`: Abschnitts-Grenze = Dokument-Grenze — Substanz-Schwelle wird
  tautologisch

### F-2 — Zwei Akzeptanzkriterien für dieselbe Nullmengen-Härte, eines auf dem abgelösten Schlüssel

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  §Akzeptanzkriterien
- `pfad`: `spec/lastenheft.md:2005` (neben dem neuen Kriterium in Zeile 2003)
- `befund`: Das Kriterium „fail-closed (kein Kandidat)" verlangt weiterhin einen
  Befund für ein Verzeichnis **ohne eine einzige Datei nach
  `planning.slice-glob`**, während die Umsetzung (und das drei Zeilen höher
  ergänzte Kriterium) den Kandidatenbegriff an `planning.closure.glob` bindet.
  Nachgemessen: Closure-Verzeichnis mit ausschließlich `welle-01-x.md`,
  `closure.glob: "*.md"` gesetzt — das Kriterium verlangt `closure-note-missing`,
  der Lauf meldet 0 Befunde, Exit 0. Beide Kriterien beschreiben denselben
  Mechanismus und widersprechen sich für jede Konfiguration, in der die beiden
  Globs auseinanderfallen.
- `verifizierbar`: ja — Fixture wie beschrieben gegen das gebaute Image; Exit 0
  gegen ein Kriterium, das Exit 1 zusagt.
- `klasse`: Akzeptanzkriterium auf widerrufener Fassung stehengeblieben

### F-3 — Die Bestandszahlen 96/111 sind nicht reproduzierbar (gemessen 95/110)

- `kategorie`: LOW
- `quelle`: [ADR-0051](../plan/adr/0051-eigener-kandidaten-filter-closure.md)
  §Kontext und §Fitness Function
- `pfad`: `docs/plan/adr/0051-eigener-kandidaten-filter-closure.md:43`, `:48`,
  `:126`; `docs/plan/planning/in-progress/slice-097-closure-glob-entkopplung.md:76`,
  `:107`; `.d-check.closure.yml:41`
- `befund`: Das Closure-Verzeichnis enthält 110 Markdown-Dateien, von denen der
  geerbte Filter 95 sieht (zweifach gemessen: Verzeichnis-Listing und Lauf mit
  einem absichtlich nie passenden `heading-pattern`, der je Kandidat genau einen
  Befund erzeugt — 110 bzw. 95 Befunde). Dokumentiert sind 110 **und** 96 in
  derselben ADR-Passage sowie „111 statt 96" in Fitness Function, DoD und
  Konfigurations-Kommentar; die Differenz von 15 ist dagegen korrekt (11
  `welle-*-results.md` + 4 Wellen-Plan-Dokumente). Die Fitness Function benennt
  damit eine Zahl, gegen die eine spätere Nachmessung nicht abgleichen kann.
- `verifizierbar`: ja — Lauf mit `closure.glob: "*.md"` und einem
  `heading-pattern`, das nichts trifft: die Zahl der `closure-note-missing`-Befunde
  ist die Kandidatenzahl.
- `klasse`: Bestandszahl in der Fitness Function nicht reproduzierbar

### F-4 — Ein YAML-`null` umgeht die Exit-2-Zusage des Schlüssels

- `kategorie`: LOW
- `quelle`: [ADR-0051](../plan/adr/0051-eigener-kandidaten-filter-closure.md)
  Entscheidung 3 („Explizit leer oder ungültig ⇒ Exit 2, kein stiller Rückfall")
- `pfad`: `internal/adapter/driven/configyaml/configyaml.go:968`
- `befund`: Die Unterscheidung *gesetzt* vs. *abwesend* läuft über den Zeiger;
  ein YAML-`null` (`glob:` ohne Wert) dekodiert als **abwesend**. Gemessen:
  `glob: ''` und `glob: ""` brechen mit Exit 2 und Nennung des Schlüssels ab,
  `glob:` läuft mit dem geerbten `slice-*.md` grün durch (Exit 0). Wer den
  Schlüssel schreibt und den Wert leert, prüft still die engere Menge — genau die
  Alternative, die die ADR-Tabelle als „übergeht stillschweigend eine gesetzte
  Aussage" verwirft. Der Rand ist im Vertrag nicht benannt: er kennt nur „nicht
  gesetzt" und „explizit leer".
- `verifizierbar`: ja — dieselbe Config einmal mit `glob:` und einmal mit
  `glob: ''` gegen das Image: Exit 0 gegen Exit 2.
- `klasse`: YAML-null als dritter Zustand neben gesetzt und abwesend

### F-5 — Die Nullmengen-Meldung gibt den Glob unquotiert wieder

- `kategorie`: LOW
- `quelle`: [§`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritt C2 (Leerlauf-Hinweis)
- `pfad`: `internal/hexagon/core/rules/planning.go:104`
- `befund`: Der Meldungstext setzt den effektiven Glob unquotiert ein. Bei einem
  Glob aus reinem Whitespace — am Config-Rand zulässig, weil weder leer noch
  syntaktisch ungültig — lautet die Meldung „enthält keine Datei nach     — das
  Gate liefe leer (fail-closed; ist der Bestand umgezogen?)" und legt dem Leser
  eine Bestands-Wanderung nahe, während die Ursache der unsichtbare
  Konfigurationswert ist. Der Befund selbst greift korrekt (Exit 1, gemessen).
- `verifizierbar`: ja — `closure.glob: '   '` gegen das Image, Ausgabe mit
  `--json` gelesen.
- `klasse`: Meldung zitiert Konfigurationswert unquotiert

### F-6 — Grund-Code-Zeile und `--doctor`-Klartext nennen den Kandidaten weiterhin „Slice"

- `kategorie`: LOW
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  §Warum ein eigener Filter („kann auch Wellen- oder Etappen-Dokumente meinen")
- `pfad`: `spec/spezifikation.md:2375` und `internal/hexagon/core/app/diagnose.go:119`
- `befund`: Die §4-Zeile beschreibt `closure-note-missing` als „Slice im
  `planning.closure.dir` ohne … Abschnitt", der `--doctor`-Klartext als
  „Abgeschlossener Slice ohne Closure-Notiz-Abschnitt". Mit dem eigenen Filter
  ist die Kandidatenmenge ausdrücklich nicht mehr auf Slices beschränkt; im
  eigenen Repo sind 15 der 110 Kandidaten Wellen-Dokumente. Ein Konsument, dessen
  Pakete anders heißen, liest im Klartext von einem Artefakt, das sein Bestand
  nicht führt. Beide Texte nennen zudem den Leerlauf-Fall (null Kandidaten) nach
  wie vor nicht, den dieser Slice mit einem eigenen Akzeptanzkriterium bekräftigt.
- `verifizierbar`: ja — Lauf mit geweitetem Glob über ein Wellen-Dokument ohne
  Notiz, Ausgabe mit `--doctor`.
- `klasse`: Grund-Code-Klartext hinter geweiteter Kandidatenmenge

### F-7 — Der Handbuch-Konfigurationsblock der Closure-Fähigkeit führt den neuen Schlüssel nicht

- `kategorie`: LOW
- `quelle`: Maintainability (Release-Prep-Rand; `docs/user/releasing.md`
  §Release-Prep, Punkt 4)
- `pfad`: `docs/user/benutzerhandbuch.md:1167` (Schlüssel-Block der zweiten
  `planning`-Fähigkeit)
- `befund`: Der Block zählt `dir`, `heading-pattern`, `min-sentences` und
  `boilerplate` auf und beschreibt anschließend die Exit-2-Ränder von
  `min-sentences` und `boilerplate`; `glob` fehlt in beidem, ebenso im
  Aufgaben-Beispiel (`docs/user/benutzerhandbuch.md:873`). Der DoD-Punkt
  „Release" ist offen — ohne Nachzug erreicht der Schlüssel den Konsumenten im
  Image, aber nicht im Handbuch, und die dortige Aufzählung liest sich als
  vollständig. Kein Gate deckt das ab (der `versions`-Gate prüft nur ghcr-Pins).
- `verifizierbar`: nein (kein Gate) — Sichtprüfung gegen die §2-Schema-Tabelle
  der Spezifikation.
- `klasse`: Handbuch-Schlüsselliste hinter dem Schema

### F-8 — Sechster Rückbau: der Config-Rand lässt sich an `closure.dir` koppeln, ohne dass ein Test rot wird

- `kategorie`: LOW
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  §fail-closed (Config-Rand)
- `pfad`: `internal/adapter/driven/configyaml/configyaml.go:969`
- `befund`: Die Prüfung des Schlüssels läuft heute unabhängig vom
  Aktivierungs-Schalter — gemessen: eine Config mit `closure.glob: ""` **ohne**
  `closure.dir` bricht mit Exit 2 ab. Wird die Bedingung auf
  `c.Glob != nil && c.Dir != ""` verengt, bleibt `make test` grün (selbst
  ausgeführt, Exit 0), obwohl derselbe Rand dann still durchliefe. Alle
  Config-Rand-Fälle der Testtabelle setzen `dir: done` und pinnen die
  Unabhängigkeit deshalb nicht; für die Nachbar-Schlüssel `heading-pattern`,
  `min-sentences` und `boilerplate` gilt dieselbe Lücke bereits vorher.
- `verifizierbar`: ja — die genannte Verengung anwenden und `make test` fahren
  (bleibt grün), dann dieselbe Config gegen das Image (Exit 2 heute).
- `klasse`: Config-Rand-Verhalten bei inerter Fähigkeit ungetestet

### F-9 — Nachbar-Schlüssel behandeln denselben Rand gegenteilig

- `kategorie`: INFO
- `quelle`: [ADR-0051](../plan/adr/0051-eigener-kandidaten-filter-closure.md)
  Entscheidung 3
- `pfad`: `internal/adapter/driven/configyaml/configyaml.go:930` gegen
  `internal/adapter/driven/configyaml/configyaml.go:969`
- `befund`: `planning.slice-glob: ''` fällt still auf `slice-*.md` zurück
  (gemessen, Exit 0), `planning.closure.glob: ''` bricht mit Exit 2 ab. Die
  Begründung der ADR („den Schlüssel zu setzen ist eine Aussage") trifft auf
  beide gleichermaßen zu; die Asymmetrie steht in keinem der beiden Dokumente.
  Für ungültige Muster verhalten sich beide gleich (Exit 2, gemessen).
- `verifizierbar`: ja — je eine Config mit explizit leerem Wert gegen das Image.
- `klasse`: gleichartige Nachbar-Schlüssel, gegenteiliger Config-Rand

## Negativbefunde

- geprüft, ohne Befund: **Byte-Identität ohne den Schlüssel.** Vor-Image
  (`39a3c6a`) und HEAD-Image gegen denselben ausgepackten Baum: Standard-Lauf,
  `--json` und `--doctor` je byte-gleich auf stdout **und** stderr, gleiche
  Exit-Codes; ebenso der Closure-Profil-Lauf (`--config .d-check.closure.yml
  --enable planning`) in allen drei Ausgabeformen.
- geprüft, ohne Befund: **Vollständigkeit der Entkopplung.** Der Filter wird an
  genau einer Produktionsstelle gelesen (`internal/hexagon/core/rules/planning.go:92`),
  die Lifecycle-Invariante liest weiterhin ausschließlich `EffectiveSliceGlob`
  (`internal/hexagon/core/rules/planning.go:45`); kein weiterer Aufrufer von
  `EffectiveClosureGlob` oder `Closure.Glob` außerhalb von Modell, Adapter und
  Tests.
- geprüft, ohne Befund: **Rückrichtung.** Ein gesetztes `closure.glob: "*.md"`
  samt geweitetem `heading-pattern` liefert denselben `planning-drift`-Befund
  (Datei, Zeile, Ziel, Grund) wie ein Lauf ohne den Schlüssel — auch über die
  konventionelle `.d-check.yml` statt über `--config`. Die Kandidaten-Sortierung
  bleibt alphabetisch stabil; zwei aufeinanderfolgende Läufe über die geweitete
  Menge liefern byte-gleiches JSON.
- geprüft, ohne Befund: **Anlass-Messung der ADR.** Nachgestellt: `slice-glob`
  auf `*.md` geweitet erzeugt den falsch-roten `planning-drift` auf der
  Roadmap-Datei; dieselbe Weitung über `closure.glob` lässt die Invariante
  unberührt.
- geprüft, ohne Befund: **Config-Rand.** `glob: ''` und `glob: ""` ⇒ Exit 2 mit
  Nennung des Schlüssels und des Auswegs; `glob: '[a-'` ⇒ Exit 2 mit
  Muster-Zitat; ein syntaktisch gültiger, aber nirgends passender Glob
  (`'*.txt'`) und ein Whitespace-Glob laufen in die Nullmengen-Härte (Exit 1,
  Befund auf dem Verzeichnis) — kein stiller Grün-Pfad in beiden Fällen.
- geprüft, ohne Befund: **Verweis-Default.** `slice-glob` gesetzt und
  `closure.glob` gesetzt (beide unabhängig wirksam), `slice-glob` abweichend und
  `closure.glob` abwesend (Closure erbt und läuft fail-closed in die Nullmenge,
  statt still zu schweigen), beide abwesend (Konventions-Default), Herkunft über
  `--config` wie über `.d-check.yml` — vier Wege, gleicher effektiver Glob.
- geprüft, ohne Befund: **Inerte Fähigkeit.** `closure.glob` ohne `closure.dir`
  bleibt wirkungslos (Exit 0, keine Datei geöffnet); die Aktivierung hängt
  unverändert allein am Verzeichnis-Schlüssel.
- geprüft, ohne Befund: **Mutations-Echtheit.** Drei der fünf behaupteten
  Rückbauten selbst nachvollzogen (Dateikopie im Temp-Verzeichnis, Mutation,
  `make test`, Rückkopie — kein `git checkout`): Filter zurück auf `slice-glob` ⇒
  zwei rote Tests; Default als Literal statt Verweis ⇒ ein roter Test; leerer
  Glob fällt still zurück ⇒ ein roter Test. Ein sechster, **nicht** gefangener
  Rückbau ist als F-8 gemeldet.
- geprüft, ohne Befund: **Eindeutigkeit im eigenen Bestand.** Alle 110 Kandidaten
  tragen **genau eine** auf `^#{1,3} .*Closure-Notiz` passende Überschrift; keine
  fremde Überschrift wird getroffen, kein Dokument wird durch die Weitung
  mehrdeutig (relevant, weil `closure-note-ambiguous` spezifiziert, aber noch
  nicht umgesetzt ist — heute gewönne still die erste Überschrift). Die 99
  H2-Treffer verteilen sich auf 95 Slices und 4 Wellen-Plan-Dokumente, die 11
  H1-Treffer auf die Ergebnis-Notizen (siehe F-1).
- geprüft, ohne Befund: **Vor-Stand-Messung der ADR-Tabelle.** Der Vor-Commit-Baum
  liefert unter `*.md` mit Default-Muster exakt die dokumentierten 12 Befunde
  (11× H1-`closure-note-missing`, 1× echter `closure-note-thin` in
  `docs/plan/planning/done/welle-70-fence-lexik.md`); mit dem geweiteten Muster
  bleibt nur der echte übrig, den der Implementierungs-Commit gefüllt hat.
- geprüft, ohne Befund: **ADR-Immutabilität.** `make adr-check
  RANGE=39a3c6a..HEAD` grün — der Verfeinerungs-Zeiger in
  [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md)
  §Geschichte ist zulässig (Sektion ausgenommen, Rumpf unberührt), und die
  Formulierung nimmt das Ergebnis der Entscheidung ausdrücklich nicht zurück.
- geprüft, ohne Befund: **SemVer-Einordnung Minor.** Additiver Schlüssel,
  Byte-Identität ohne ihn belegt; ältere Fassungen brechen an der strikten
  Feld-Prüfung fail-closed ab, statt den unbekannten Schlüssel still zu ignorieren
  — kein stiller Prüfumfangs-Wechsel bei Versions-Rückschritt.
- geprüft, ohne Befund: **Vertrags-Konsistenz.** Spezifikation Schritt C2, die
  §2-Schema-Zeile (Position direkt nach `closure.dir`, Default als Verweis
  ausgewiesen), die Index-Zeile in `docs/plan/adr/README.md`, die
  `--print-config`-Vorlage (`internal/adapter/driving/cli/config_template.go:128`)
  und der Slice-Plan sagen dasselbe; die Lastenheft-Historie führt 0.53.0 mit dem
  CR-Bezug. Ausnahme: F-2 und F-7.
- geprüft, ohne Befund: **Gates.** `make gates` grün (355 Dateien, 0 Befunde,
  alle acht Gates), `make verify-closure-notes` grün (326 Dateien, 0 Befunde).
- geprüft, ohne Befund: **Formatierung.** `internal/adapter/driven/configyaml/configyaml_test.go`
  endet mit einer zusätzlichen Leerzeile; das Lint-Profil führt keinen
  Formatier-Linter, und ohne Konventions-Anker ist das kein Finding (Skill
  §Anti-Pattern).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 6 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Abschnitts-Grenze = Dokument-Grenze ·
Akzeptanzkriterium auf widerrufener Fassung stehengeblieben · Bestandszahl in der
Fitness Function nicht reproduzierbar · YAML-null als dritter Zustand ·
Meldung zitiert Konfigurationswert unquotiert · Grund-Code-Klartext hinter
geweiteter Kandidatenmenge · Handbuch-Schlüsselliste hinter dem Schema ·
Config-Rand-Verhalten bei inerter Fähigkeit ungetestet · gleichartige
Nachbar-Schlüssel, gegenteiliger Config-Rand

## Verdikt

**Merge-blockierend:** ja — die beiden MEDIUM-Befunde. **Nicht** blockierend ist
der Produkt-Code: die zugesagte Byte-Identität hält am Lauf gegen ein aus dem
Vor-Commit gebautes Image (Standard, `--json`, `--doctor`, Closure-Profil), die
Entkopplung ist vollständig, der Verweis-Default trägt in allen geprüften
Kombinationen, und die Exit-2-Ränder greifen. Blockierend sind eine
Vertrags-Selbstwidersprüchlichkeit (F-2) und eine Zusage der eigenen
Konfiguration, die der Lauf so nicht einlöst (F-1) — beide sind Text- bzw.
Konfigurations-Ränder, nicht Kern-Korrektheit.

**Release-Empfehlung:** **noch nicht freigeben.** Die Reihenfolge ist wichtig,
weil [ADR-0051](../plan/adr/0051-eigener-kandidaten-filter-closure.md) mit der
Closure auf `Accepted` und damit auf immutabel geht: F-1 (Reichweite der Zusage)
und F-3 (Bestandszahlen) berühren ihre Fitness Function und ihren Kontext und
sind **vor** dem Statuswechsel zu klären, sonst friert eine nicht nachmessbare
Zahl ein. F-2 gehört in dieselbe Runde, weil das Lastenheft mit 0.53.0
veröffentlicht wird. F-7 gehört in die Release-Prep-Runde (Handbuch), F-4, F-5,
F-6, F-8 und F-9 sind release-verträglich und können in einen Folge-Schnitt —
naheliegend zusammen mit
[slice-098](../plan/planning/open/slice-098-closure-note-placeholder.md), das
dieselbe Platzhalter-Klasse adressiert wie F-1.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die
Slice-Closure §7 und von dort in den Zähler des Beobachtungs-Registers. Dieser
Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill, dieses Modell, dieses
Verdikt) und ersetzt keine Verifikation — die DoD-/Spec-Konformität prüft der
Verifier separat.
