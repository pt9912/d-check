# Review-Report: slice-102 — Wellen-Lifecycle-Invariante — 2026-08-16

**Review-Art:** Code — geprüft wird der ausgelieferte Diff gegen die Verträge,
die er selbst anlegt (Anforderung, Algorithmus W1–W5, §2-Schema, §4-Codes), gegen
die referenzierten ADRs und gegen die Messung in §3a des Slice-Plans. Nicht
geprüft wird die DoD-Abhakung (getrennter Kontext, Verifikation).

**Gegenstand:** Commit-Range `9c7506d..a040eff` — **fünf** Commits
(Wellen-Eröffnung · Messung · ADR · Vertrag · Implementierung; der Auftrag nannte
vier, der ADR-Commit liegt separat); Arbeitsbaum-Stand `a040eff` (= HEAD).

**Skill:** `.harness/skills/reviewer.md` @ 1.4.0 ·
**Modell:** claude-fable-5 · **Datum:** 2026-08-16

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  §Wellen-Invariante (Lastenheft-Fassung 0.59.0, samt Akzeptanzkriterien)
- §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritte W1–W5, das §2-Schema (`planning.waves.*`), die vier §4-Zeilen
  (`wave-drift` · `wave-preview-exists` · `wave-results-missing` ·
  `wave-unregistered`) und
  §[`DC-FA-TGT-001.a`](../../spec/spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets)
  Schritt 3 (die Tabellenzeilen-Lexik, deren zweiter Konsument hier entsteht)
- [ADR-0055](../plan/adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
  (Proposed — Entscheidungen 1–5, Konsequenzen),
  [ADR-0028](../plan/adr/0028-planning-lifecycle-modul.md),
  [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md),
  [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  (Entscheidungen 1 und 4)
- [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  (Spiegel einer Semantik-Änderung), [`AGENTS.md`](../../AGENTS.md) §3 Hard
  Rules, [`CLAUDE.md`](../../CLAUDE.md)
- Der Slice-Plan
  [slice-102](../plan/planning/done/slice-102-wellen-lifecycle-invariante.md)
  (§3 Abnahme-Punkte mit Ausgängen, §3a Messung) und das Wellendokument
  [welle-75](../plan/planning/done/welle-75/welle-75-wellen-register.md)

**Läufe dieses Reviews.** Alle Fixtures und alle Mutations-Kopien in einem
Scratch-Verzeichnis außerhalb des Repos; keine Probe-Config wurde im Repo
abgelegt. Gefahren: `make test` als Baseline (grün) plus **sechs**
Mutationsläufe (je über Dateikopie mutieren → `make test` → Kopie
zurückschreiben, Ergebnis am **Exit-Code** abgelesen, nie per `git checkout`);
rund zwanzig Fixture-Läufe gegen das HEAD-Image (netzlos, read-only); der
`planning-check`-äquivalente Lauf über den eigenen Baum (378 Dateien, 0
Befunde); die §3a-Messung repliziert über **lesende Kopien** der Planungs-Bäume
des Schwester-Repos und des Kurs-Beispiels (an beiden Originalen nichts
geändert). Nach dem letzten Rückbau ist die Repo-Datei byte-identisch zur
Sicherungskopie (`diff` leer); `git status --short` ist am Ende leer bis auf
diesen Report.

---

## Findings

### F-1 — Ein unlesbares `waves.dir` schaltet die Fähigkeit im Ruhe-Zustand still ab — genau in der zweimal belegten Defekt-Richtung; der Kommentar behauptet das Gegenteil

- `kategorie`: HIGH
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  (fail-closed-Disziplin der Modul-Familie; Schritt C2 desselben Moduls:
  „Fehlt das gesetzte Verzeichnis oder ist es unlesbar ⇒ … fail-closed, kein
  stilles Grün") · Reviewer-Anker „Stilles-Grün-Pfad in einem Gate"
- `pfad`: `internal/hexagon/core/rules/planning_waves.go:41` (Kommentar) und
  `:218` (`dirNames`); Vertragslücke in
  §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  W1/W2
- `befund`: `dirNames` liefert bei einem Listing-Fehler eine leere Menge ohne
  jeden Befund. Scratch-Fixture: Roadmap mit Ruhe-Marker, ein flaches
  Wellendokument im **realen** Verzeichnis, `waves.dir` mit Tippfehler ⇒ **0
  Befunde, Exit 0** — die Aussage-2-Verletzung, die in den Closures von
  welle-68 und welle-69 zweimal real eingetreten ist und die diesen Slice
  motiviert hat, ist mit einem Pfad-Tippfehler unsichtbar, dauerhaft und
  kommentarlos. Der Docstring von `waveSets` behauptet „der daraus folgende
  Befund ist fail-closed, nicht still" — im Ruhe-Zustand folgt kein Befund. Im
  Aktiv-Zustand kommt zwar ein `wave-drift`, aber mit falscher Ursachen-Nennung
  („liegt kein flaches Wellendokument", tatsächlich ist das Verzeichnis
  unlesbar). Die Closure-Fähigkeit desselben Moduls meldet den identischen Fall
  ausdrücklich („sonst schaltete ein Tippfehler im Pfad die ganze Prüfung ab");
  W1/W2 sagen für `waves.dir`/`done-dir` nichts dergleichen zu — Code- und
  Vertragslücke liegen im selben Diff.
- `verifizierbar`: ja — das Scratch-Fixture (Ruhe-Marker + flache Datei +
  Tippfehler-`dir`) gegen das Image: Exit 0; mit korrektem `dir`: `wave-drift`,
  Exit 1.
- `klasse`: „Tippfehler im Pfad schaltet ein Gate still ab"

### F-2 — Fehlt eine Register-Überschrift, entfallen alle drei Register-Aussagen kommentarlos

- `kategorie`: MEDIUM
- `quelle`: §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  W4/W5 gegen den Heading-Guard aus Schritt 3 (fail-closed war dort eine der
  drei Richtungen des abgelösten Skripts)
- `pfad`: `internal/hexagon/core/rules/planning_waves.go:148` und `:158`
  (`no == 0` ⇒ leere Menge bzw. `closedNo == 0` ⇒ W5c entfällt)
- `befund`: Trifft `next-heading` oder `closed-heading` keine Zeile der
  Roadmap (Tippfehler, umbenannter Abschnitt, abweichende Konsumenten-Konvention),
  werden W5a–c wortlos übersprungen. Scratch-Fixture: Abschluss-Register unter
  „## Fertige Wellen", eine Register-Zeile ohne Notiz **und** eine Notiz ohne
  Register-Zeile ⇒ 0 Befunde, Exit 0. Die Slice-Invariante behandelt die
  fehlende kanonische Überschrift als `planning-drift` (fail-closed); die neue
  Fähigkeit übernimmt diese Disziplin für ihre beiden Abschnitts-Anker nicht,
  und der Vertrag regelt den Fall nicht.
- `verifizierbar`: ja — Scratch-Fixture mit umbenanntem Register: Exit 0 trotz
  beidseitiger Verletzung.
- `klasse`: „Abschnitts-Anker ohne fail-closed-Guard"

### F-3 — Die Kennungs-Ableitung ist unterspezifiziert und fällt still auf ein hartes `welle-` zurück; der Präfix stammt allein aus dem Plan-Glob

- `kategorie`: MEDIUM
- `quelle`: §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  W2 („Kennung ist jeweils das Zahlen-Präfix `welle-<n>`") und W4
- `pfad`: `internal/hexagon/core/rules/planning_waves.go:194` (`wavePrefix`,
  Rückfall `welle-`), `:46` und `:97` (Präfix nur aus `EffectiveGlob()`)
- `befund`: Drei belegte Folgen. (1) Ein Glob mit **führendem** Platzhalter
  (`*.md` — am Config-Rand gültig) liefert den Rückfall-Präfix `welle-`:
  Scratch-Fixture mit `wave-*`-Dateien meldet ein falsches `wave-drift`
  („kein flaches Wellendokument", obwohl eines liegt) und schaltet die
  Register-Aussagen still ab. (2) Ein `results-glob` mit anderem Präfix als der
  Plan-Glob (etwa `ergebnis-*.md`) liefert leere Ergebnis-Kennungen:
  Scratch-Fixture meldet `wave-results-missing`, obwohl die Notiz nach dem
  konfigurierten Muster im Ruheort liegt, und `wave-unregistered` kann für
  diese Notizen nie feuern. (3) Der Vertrag legt `welle-<n>` **wörtlich** fest
  und sagt zur Ableitung des Präfixes aus dem Glob nichts — für den regulären
  `wave-*.md`-Konsumenten (funktioniert, gemessen) beschreibt der Vertragstext
  die falsche Kennung, für die Fälle 1/2 beschreibt er das Verhalten gar nicht.
- `verifizierbar`: ja — die beiden Scratch-Fixtures (Glob `*.md`;
  `results-glob` mit fremdem Präfix) gegen das Image.
- `klasse`: „Kennungs-Lexik unterspezifiziert, Ableitung mit stillem Rückfall"

### F-4 — `wave-unregistered` trägt die Kennung als `target`, der Vertrag verlangt die Notiz

- `kategorie`: MEDIUM
- `quelle`: §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  W5c („Befund an der `closed-heading`-Zeile mit der **Notiz** als `target`")
- `pfad`: `internal/hexagon/core/rules/planning_waves.go:127` (Aufruf mit `id`);
  Ursache `:56`–`:64` (`waveSets` wirft die Basisnamen weg)
- `befund`: Der gelieferte Befund lautet `file` = Roadmap, `line` =
  `closed-heading`-Zeile (vertragskonform), `target` = `welle-9` — die
  **Kennung**, nicht die Ergebnisnotiz. Die Notiz ist zur Meldezeit nicht mehr
  greifbar, weil `waveSets` nur Kennungs-Mengen liefert; die Meldung nennt
  ersatzweise das Verzeichnis. Ein Konsument, der (wie beim Ventil-Paritäts-
  Muster der übrigen Module) auf das `target` filtert oder es auflöst, bekommt
  einen Wert, der keine Datei ist. Befund-Formen der drei übrigen Codes sind
  vertragskonform (geprüft, siehe Negativbefunde).
- `verifizierbar`: ja — Scratch-Fixture „Notiz ohne Register-Zeile": die
  `target`-Spalte der Ausgabe zeigt die Kennung.
- `klasse`: „Befund-Form weicht vom Algorithmus-Vertrag ab"

### F-5 — Eine eingeschleuste eigene Aktiv-Status-Bestimmung überlebt `make test`: die geteilte Antwort ist am neuen Konsumenten nicht festgenagelt

- `kategorie`: MEDIUM
- `quelle`: [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 4 („je Konsument eine Assertion … sonst ist die Invariante nur
  ein Kommentar");
  [ADR-0055](../plan/adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
  Entscheidung 5
- `pfad`: `internal/hexagon/core/rules/planning_waves.go:30` (Aufruf von
  `planningActiveStatus`); `internal/hexagon/core/rules/planning_waves_test.go`
  (kein Test bindet den Status an die geteilte Mechanik)
- `befund`: Rückbau-Probe über Dateikopie: der `planningActiveStatus`-Aufruf in
  `CheckPlanningWaves` durch eine naive lokale Bestimmung ersetzt (roher
  `## `-Präfix-Vergleich, keine Fence-Behandlung, kein Mehrfach-Guard) —
  `make test` bleibt **grün** (Exit 0). Der strukturelle Teil der Lieferung
  stimmt (es gibt genau eine Funktion, beide Fähigkeiten rufen sie auf), aber
  die von ADR-0054 verlangte Verankerung fehlt: kein Wellen-Test legt den
  Status in eine Lage, in der die geteilte und eine naive Antwort
  auseinanderfallen (Marker hinter einer Fence-Zeile im Aktiv-Block,
  eingerückte oder tab-getrennte H2, doppelte Überschrift). Die zentrale
  Zusage des Commits („der Aktiv-Status wird nicht zweimal bestimmt") ist
  damit heute wahr und morgen unbewacht — exakt die Klasse, die welle-74
  sechsmal geheilt hat.
- `verifizierbar`: ja — der beschriebene Mutant gegen `make test`: Exit 0.
- `klasse`: „geteilte Antwort ohne Assertion am neuen Konsumenten"

### F-6 — Die flach-Hälfte von W5a ist untestet: ein Mutant, der nur den Ruheort prüft, überlebt

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  (W5a: Datei „flach **oder** im Ruheort") · Reviewer-Anker „fehlende
  Negativtests bei neuem öffentlichen Vertrag"
- `pfad`: `internal/hexagon/core/rules/planning_waves.go:101`
  (`flach[r.id] || ruhe[r.id]`);
  `internal/hexagon/core/rules/planning_waves_test.go` (Vorschau-Test nur mit
  Ruheort-Datei)
- `befund`: Rückbau-Probe: die Bedingung auf `ruhe[r.id]` verkürzt —
  `make test` bleibt grün (Exit 0). Der einzige Vorschau-Positiv-Test legt die
  Datei in den Ruheort; der nächstliegende reale Fall — die Vorschau-Zeile
  nennt die **laufende** Welle, deren Plan-Dokument flach liegt (die bei der
  Eröffnung nicht geleerte Vorschau-Zeile) — ist ungedeckt und würde unter dem
  Mutanten nicht mehr gemeldet.
- `verifizierbar`: ja — der beschriebene Mutant gegen `make test`: Exit 0.
- `klasse`: „Vertrags-Oder ohne Test für eine Hälfte"

### F-7 — Die W1-Zusage „Exit 2 bei leerem `next-heading`/`closed-heading`" ist unerreichbar: explizit leer fällt still auf den Default zurück

- `kategorie`: MEDIUM
- `quelle`: §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  W1; Vergleichsmaßstab Schritt C2 desselben Moduls („ein Rückfall auf den
  Default wäre ein stilles Übergehen einer gesetzten Aussage")
- `pfad`: `internal/adapter/driven/configyaml/configyaml.go` (`rawWaves`,
  Klartext-Strings; `applyWaves`, Prüfung nur `h != "" && TrimSpace == ""`)
- `befund`: `next-heading: ""` und `closed-heading: ""` laufen mit Exit 0
  durch und verwenden still die Defaults (Scratch-Fixture belegt es); nur eine
  **Whitespace**-Überschrift bricht ab. Der Vertrag verspricht Exit 2 „bei
  leerem `next-heading`/`closed-heading`". `rawWaves` führt die Felder als
  Klartext-Strings und kann „explizit leer" von „abwesend" nicht
  unterscheiden — die Nachbar-Fähigkeit löst dieselbe Frage mit einem
  Zeiger-Feld (`rawClosure.Glob *string`) und macht die gesetzte Null-Aussage
  zum Config-Fehler. Dieselbe Unschärfe trifft `glob: ""`/`results-glob: ""`
  (still Default), wo der Vertrag immerhin nichts Gegenteiliges zusagt.
- `verifizierbar`: ja — Scratch-Fixture mit `next-heading: ""`: Exit 0 statt 2.
- `klasse`: „Config-Rand-Zusage ohne erreichbare Unterscheidung
  gesetzt/abwesend"

### F-8 — Der zwölfte Schwester-Repo-Befund ist ein Marker-Artefakt der Probe-Konfiguration, keine Verletzung von Aussage 1

- `kategorie`: MEDIUM
- `quelle`: Slice-Plan §3a;
  [ADR-0055](../plan/adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
  Konsequenzen; `CHANGELOG.md` (Unreleased)
- `pfad`: [slice-102 §3a](../plan/planning/done/slice-102-wellen-lifecycle-invariante.md)
  (Zeilen 123 und 148–152) ·
  `docs/plan/adr/0055-wellen-invariante-artefakt-und-grund-codes.md:86` ·
  `CHANGELOG.md:31`
- `befund`: Die Replikation über eine lesende Kopie des Schwester-Baums
  bestätigt die **11** fehlenden Ergebnisnotizen (13 Register-Zeilen, 2
  Notizen) und reproduziert die **12** — aber nur unter d-checks
  Default-Marker. Der Aktiv-Block des Schwester-Repos beginnt mit „**Keine.**"
  und erklärt die letzte Welle ausdrücklich für geschlossen; mit einem
  konsument-gerechten Marker verschwindet das `wave-drift` und es bleiben 11.
  Die Aussage „a-check verletzt Aussage 1 heute (aktive Welle in der Roadmap,
  kein Wellendokument)" — im Slice, in den ADR-Konsequenzen („eine Roadmap,
  die eine aktive Welle ohne Wellendokument nennt") und als „12" im
  CHANGELOG — beschreibt damit kein Bestands-Faktum, sondern ein Artefakt der
  Probe-Konfiguration: die Roadmap nennt keine aktive Welle, sie ruht unter
  einer abweichenden Marker-Formulierung. Die tragenden Zahlen des Slice (19,
  11, 0, 15/15, 1/1) halten der Nachmessung stand.
- `verifizierbar`: ja — zwei Läufe über die Kopie des Schwester-Baums
  (Default-Marker: 12 Wellen-Befunde; Marker „**Keine.**": 11).
- `klasse`: „Messwert trägt Konfigurations-Artefakt als Bestands-Befund"

### F-9 — Der MR-025-Spiegel „Emittierte Vorlage" fehlt: `--print-config` kennt `planning.closure`, aber kein `planning.waves`

- `kategorie`: MEDIUM
- `quelle`: [`MR-025`](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  (Spiegel-Liste: „Emittierte Vorlage — `--print-config` / `--suggest-config`");
  [`DC-FA-CLI-005`](../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben)
- `pfad`: `internal/adapter/driving/cli/config_template.go:119`–`:132`
  (planning-Block endet bei den `closure`-Schlüsseln)
- `befund`: Die im Binary ausgelieferte Konfigurations-Vorlage führt den
  `planning`-Block mit allen `closure`-Schlüsseln, aber keinem der sechs
  `waves`-Schlüssel (`grep waves` über die Vorlagen-Datei: null Treffer). Die
  Vorlage ist Produkt-Code, nicht Release-Prep-Doku — der bewusste
  Handbuch-/README-Aufschub deckt sie nicht. Ein Konsument, der die neue
  Fähigkeit nach der Vorlage konfiguriert, findet sie dort nicht.
  `--suggest-config` emittiert für `planning` nur den `roadmap`-Schlüssel und
  behandelt `closure` ebenfalls nicht — dort ist der Verzicht auf `waves`
  konsistent (geprüft, kein Befund).
- `verifizierbar`: ja — `--print-config`-Ausgabe des Images enthält keinen
  `waves`-Schlüssel.
- `klasse`: „Semantik-Spiegel bei der Lieferung ausgelassen"

### F-10 — Die Tabellenzeilen-Erkennung ist eine zweite Inline-Implementierung, kein geteiltes Prädikat

- `kategorie`: MEDIUM
- `quelle`: [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidungen 1 und 4;
  §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  W4 („dieselbe Lexik wie
  [`DC-FA-TGT-001.a`](../../spec/spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets)")
- `pfad`: `internal/hexagon/core/rules/planning_waves.go:178` gegen
  `internal/hexagon/core/rules/targets.go:156`
- `befund`: Beide Konsumenten beantworten „ist das eine Tabellenzeile?" mit
  demselben inline geschriebenen Ausdruck (Prosa-Zeile und erstes Zeichen
  Pipe) — die Antwort ist heute **identisch** (verglichen, einschließlich
  „Einrückung zählt nicht"), und die Assertion am neuen Konsumenten existiert
  (der Fence-Test wird beim Entfernen des Prosa-Filters rot, Mutationslauf
  belegt). Was fehlt, ist das geteilte Prädikat selbst: ändert `targets` seine
  Zeilen-Erkennung, folgt `planning.waves` nicht, und kein Test bemerkt es —
  die Zusage „dieselbe Lexik" hängt an zwei Kopien statt an einer Funktion.
  Nach ADR-0054 ist die Reparatur „Übernahme der vorhandenen Antwort, nicht
  eine zweite Implementierung"; die Vorgänger-Welle hat dieselbe Halbform
  („Assertion ohne geteiltes Prädikat" bzw. umgekehrt) mehrfach als Defekt
  behandelt. Das Wellendokument kündigt für den dritten Konsumenten (BEO-005)
  ohnehin einen Kopplungs-Test an — der braucht die Funktion.
- `verifizierbar`: ja — strukturell am Code (kein gemeinsames Symbol); die
  Drift-Richtung wäre erst nach einer künftigen `targets`-Änderung messbar.
- `klasse`: „geteilte Antwort halb übernommen — zwei Kopien statt einer
  Funktion"

### F-11 — Fehlende Roadmap meldet doppelt, unbestimmbarer Aktiv-Status bewusst einfach: die Randfall-Disziplin ist uneinheitlich

- `kategorie`: LOW
- `quelle`: §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  (W-Schritte regeln den Fall nicht) · Maintainability
- `pfad`: `internal/hexagon/core/rules/planning_waves.go:24`–`:35`
- `befund`: Bei fehlender/unlesbarer Roadmap melden **beide** Fähigkeiten
  (`planning-drift` und `wave-drift`, Scratch-Fixture belegt es) für dieselbe
  eine Ursache mit derselben einen Reparatur; beim nicht bestimmbaren
  Aktiv-Status schweigt die Wellen-Fähigkeit dagegen ausdrücklich („kein
  Doppelbefund"). Beide Formen sind vertretbar, aber sie widersprechen sich,
  und der Vertrag entscheidet keinen der beiden Fälle.
- `verifizierbar`: ja — Scratch-Fixture mit fehlender Roadmap: zwei Befunde.
- `klasse`: „Randfall-Disziplin uneinheitlich"

### F-12 — Der W5-Körpersatz „Grund-Codes (§4) folgen mit der Implementierung" ist nach der Lieferung stehengeblieben

- `kategorie`: LOW
- `quelle`: §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  W5 gegen die §4-Tabelle
- `pfad`: `spec/spezifikation.md:1760`
- `befund`: Der Implementierungs-Commit hat die vier §4-Zeilen ergänzt, den
  Zukunfts-Satz im Algorithmus-Körper aber stehen lassen; bei den
  Vorgänger-Fähigkeiten steht dieser Satz nur noch in den Historien-Zeilen
  (dort korrekt, als Zeitpunkt-Aussage). Ein Leser des Körpers schließt heute,
  die Codes seien noch nicht ausgeliefert.
- `verifizierbar`: ja — Textbefund; §4 führt alle vier Codes.
- `klasse`: „Rand auf der alten Fassung stehengeblieben"

### F-13 — `CheckPlanningWaves` liest und parst die Roadmap ein zweites Mal

- `kategorie`: LOW
- `quelle`: Maintainability; Geist von
  [ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
  Entscheidung 1 (eine Frage, eine Antwort — hier: eine Eingabe, eine Lesung)
- `pfad`: `internal/hexagon/core/rules/planning_waves.go:24` und `:95`–`:96`
- `befund`: Beide Fähigkeiten lesen dieselbe Roadmap-Datei je einmal über den
  Port und berechnen `splitLines`/`proseLineSet`/`planningActiveStatus`
  getrennt. Unter der zugesagten Umgebung
  ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus):
  identischer, read-only Arbeitsbaum) ist das reine Ineffizienz und keine
  beobachtbare Divergenz; die Bauform hält die „eine Größe, eine
  Quelle"-Zusage aber nur deshalb, weil beide Lesungen zufällig denselben
  Inhalt sehen — die naheliegende Härtung wäre, den Status (oder den Inhalt)
  zu übergeben statt neu zu gewinnen. Latente Wartungsfalle, kein aktueller
  Fehler.
- `verifizierbar`: ja — strukturell am Code (zwei `ReadFile`-Aufrufe für
  dieselbe Datei im selben Lauf).
- `klasse`: „doppelte Lesung derselben Eingabe"

### F-14 — Kopf- und Trennzeile sind nur durch den Kennungs-Filter ausgeschlossen, nicht strukturell

- `kategorie`: INFO
- `quelle`: §[`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  W4 („ohne Kopf- und Trennzeile")
- `pfad`: `internal/hexagon/core/rules/planning_waves.go:177`–`:188`
- `befund`: `waveRowsFrom` überspringt Kopf- und Trennzeile nicht als solche;
  sie fallen nur heraus, weil ihre erste Spalte üblicherweise keine Kennung
  trägt. Ein Register, dessen **Spaltentitel** eine Kennung enthält (etwa eine
  Kopfzeile „| welle-75-Artefakte | … |"), würde als Datenzeile gewertet und
  könnte ein falsches `wave-results-missing` erzeugen. Konstruiert, im
  Bestand nicht vorhanden — als dokumentationswürdige Abweichung vom
  Vertragstext notiert.
- `verifizierbar`: ja — konstruierbares Fixture; im Bestand ohne Treffer.
- `klasse`: „Vertragswort ohne strukturelle Entsprechung"

## Negativbefunde

- geprüft, ohne Befund: **W2-Mengenbildung und Kennungs-Grenzfälle.**
  Präfix-Kollision ein-/zweistellig (`welle-7` gegen `welle-75`) in beiden
  Richtungen sauber getrennt (greedy Ziffernlauf plus exakter Mengen-Vergleich;
  je zwei Scratch-Fixtures); doppelte Register-Zeilen derselben Kennung melden
  je Vorkommen an ihrer Zeile; dieselbe Kennung in Vorschau **und**
  Abschluss-Register erzeugt genau die zutreffenden Befunde; „zwei Dateien mit
  demselben Präfix zählen als eine Welle" hält; `waves.dir == done-dir` trennt
  die Rollen korrekt über die Globs.
- geprüft, ohne Befund: **Rollen-Trennung und W3.** Der `results-glob`-Abzug
  ist implementiert und getestet (Mutationslauf: Abzug entfernt ⇒ Test rot,
  Exit 2); „genau eins" statt „mindestens eins" ist getestet (Mutationslauf:
  `>= 1` ⇒ Drift-Subtest rot); die W3-Meldung nennt die Richtung und bei
  mehreren Dokumenten deren Zahl, Befund an der `planning.heading`-Zeile mit
  `dir` als `target` — vertragskonform.
- geprüft, ohne Befund: **Tabellenzeilen-Lexik am neuen Konsumenten
  (Verhalten).** Erste-Spalte-Adresse statt Zeilen-Scan getestet
  (Mutationslauf: ganze Zeile ⇒ Namens-Vorschau-Test rot); Fence-Filter
  getestet (Mutationslauf: Prosa-Filter entfernt ⇒ Fence-Test rot); die
  Antwort ist mit der `targets`-Lexik verhaltensgleich verglichen
  (Spalte-0-Pipe, Einrückung zählt nicht). Offen bleibt allein die
  Prädikat-Teilung (F-10).
- geprüft, ohne Befund: **Grund-Code-Flächen.** `AllReasons()`, `reasonTexts()`
  und die vier §4-Zeilen sind deckungsgleich (vier Codes, Reihenfolge im
  Lockstep); der `--doctor`-Klartext zu `wave-drift` rendert am Fixture korrekt;
  die Vier-Codes-Begründung (Deduplikation über Datei, Zeile, Regel, Ziel,
  Grund) trägt — die einzige konstruierbare Selbe-Zeile-Kollision entsteht bei
  identisch konfigurierten Register-Überschriften und wird durch die
  getrennten Codes aufgelöst.
- geprüft, ohne Befund: **Config-Rand, positive Seite.** `dir` mit
  `..`-Segment und ein ungültiges Glob brechen mit Exit 2 und benanntem
  Schlüssel ab (zwei Scratch-Fixtures); `done-dir`-Default `<dir>/done` und
  alle Effective-Defaults entsprechen dem §2-Schema; YAML-abwesende Schlüssel
  liefern die Konventions-Defaults.
- geprüft, ohne Befund: **Opt-in und Byte-Identität.** Ohne `planning.waves.dir`
  wird kein Wellendokument geöffnet (Test vorhanden; Kern-Guard in W1); die
  Fähigkeit liegt außerhalb der Default-`modules`, `make doc-check` aktiviert
  sie nicht; der eigene Lauf ohne `--enable planning` ist unverändert.
- geprüft, ohne Befund: **Messung §3a, tragende Zahlen.** Eigener Baum: der
  `planning-check`-äquivalente Lauf meldet 0 Befunde (378 Dateien); 15
  Register-Zeilen gegen 15 Ergebnisnotizen; genau ein flaches Wellendokument
  bei aktiver Welle; die eigene Vorschau-Zeile nennt `welle-75` in der
  Trigger-Spalte (der Fall, der die Spalten-Adresse erzwingt). Kurs-Beispiel:
  0 Befunde über die lesende Kopie. Schwester-Repo: 11 fehlende
  Ergebnisnotizen bestätigt; nur der zwölfte Befund ist strittig (F-8).
- geprüft, ohne Befund: **Hexagon-Schnitt und Hermetik.**
  `planning_waves.go` importiert ausschließlich Kern-Model und den
  Filesystem-Port (keine Adapter, kein git, kein Netz) — konform zu
  [ADR-0005](../plan/adr/0005-modul-layout-hexagon-ordner.md) und zur
  Hermetik-Zusage aus
  [ADR-0028](../plan/adr/0028-planning-lifecycle-modul.md); Determinismus über
  sortierte Listings und sortierte Fehlend-Liste
  ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)),
  Wiederholungsläufe identisch.
- geprüft, ohne Befund: **Scharfschaltung im eigenen Baum.** Der
  `waves`-Block in der eigenen Konfiguration setzt nur `dir` und dokumentiert
  Messung und Anlass; `make planning-check` bleibt der Bindepunkt in `gates`;
  der Kommentar widerspricht keiner Vertragsfläche.
- geprüft, ohne Befund: **CHANGELOG und Lastenheft-Historie.** Der
  Unreleased-Eintrag und die 0.59.0-Zeile beschreiben Lieferung, Opt-in-Form
  und die zwei gemessenen Festlegungen konsistent zu Anforderung und
  Algorithmus — bis auf die in F-8 benannte Zwölf. Handbuch, READMEs und
  Operations-Doku sind erklärtermaßen Release-Prep und wurden nur auf aktiven
  Widerspruch geprüft: keiner gefunden.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 9 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Tippfehler im Pfad schaltet ein Gate still ab ·
Abschnitts-Anker ohne fail-closed-Guard · Kennungs-Lexik unterspezifiziert,
Ableitung mit stillem Rückfall · Befund-Form weicht vom Algorithmus-Vertrag ab ·
geteilte Antwort ohne Assertion am neuen Konsumenten · Vertrags-Oder ohne Test
für eine Hälfte · Config-Rand-Zusage ohne erreichbare Unterscheidung
gesetzt/abwesend · Messwert trägt Konfigurations-Artefakt als Bestands-Befund ·
Semantik-Spiegel bei der Lieferung ausgelassen · geteilte Antwort halb
übernommen · Randfall-Disziplin uneinheitlich · Rand auf der alten Fassung
stehengeblieben · doppelte Lesung derselben Eingabe · Vertragswort ohne
strukturelle Entsprechung.

**Wiederholungs-Signal.** Zwei Klassen sind Wiedergänger aus den Reviews der
Vorgänger-Welle: „geteilte Antwort ohne Assertion am Konsumenten" (F-5, F-10 —
dort mehrfach Gegenstand von
[ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)) und
„Rand auf der alten Fassung stehengeblieben" (F-12). Dass beide in dem Slice
wiederkehren, der ADR-0054 ausdrücklich auf sich selbst anwendet, ist ein
Steuerungs-Signal in Richtung eines mechanischen Kopplungs-Tests — den das
Wellendokument für den dritten Lexik-Konsumenten bereits ankündigt.

## Verdikt

**Merge-blockierend:** ja — ein HIGH und neun MEDIUM.

**Was trägt:** Der Kern der Lieferung ist richtig gebaut und gut belegt. Die
vier Aussagen sind korrekt geschnitten (vier Codes, vier Reparaturen, die
Deduplikations-Begründung hält), die Rollen-Trennung Plan-Dokument/Ergebnisnotiz
funktioniert einschließlich der Präfix-Kollisionen ein- und zweistelliger
Nummern, die Spalten-Adresse verhindert nachweislich den Zeilen-Scan-Fehler,
und vier von sechs eigenen Mutations-Gegenproben werden von genau dem
zugehörigen Test gefangen. Die tragenden Messzahlen (19 · 11 · 0 · 15/15 · 1/1)
halten einer unabhängigen Replikation stand; die Scharfschaltung im eigenen
Baum ist durch den grünen Lauf gedeckt.

**Was blockiert:** F-1 ist die Harness-Lüge im engen Sinn — ein
Pfad-Tippfehler macht im Ruhe-Zustand genau die zweimal real eingetretene
Verletzung unsichtbar, für die der Slice gebaut wurde, und der Code-Kommentar
behauptet fail-closed. F-2 und F-7 sind dieselbe Disziplin-Lücke am
Abschnitts- bzw. Config-Rand; alle drei stehen im Kontrast zur
Closure-Fähigkeit desselben Moduls, die jeden dieser Fälle ausdrücklich meldet.
F-5 und F-6 zeigen mit überlebenden Mutanten, dass zwei zentrale Zusagen
(geteilter Aktiv-Status; „flach oder im Ruheort") heute nur Kommentare sind.
F-4 und F-3 sind Vertragstreue-Lücken in Befund-Form und Kennungs-Lexik, F-8
korrigiert eine Zahl, die auf drei Flächen als Bestands-Befund steht, F-9 den
fehlenden Vorlage-Spiegel, F-10 die halbe Übernahme der geteilten Lexik.

**Release-Empfehlung: noch nicht releasen.** Die Minor-Einordnung (opt-in, ohne
die neuen Schlüssel byte-identisch) ist korrekt und bleibt es. Vor dem Tag
gehören F-1, F-2, F-5, F-6 und F-7 in diesen Slice (Code plus die zugehörigen
W-Schritte — [ADR-0055](../plan/adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
ist noch Proposed, die Verträge sind ohne Supersede korrigierbar); F-3 und F-4
sind je eine Entscheidung (Vertrag präzisieren oder Code angleichen) auf
demselben Diff; F-8 ist eine Textkorrektur auf drei Flächen, F-9 ein
Vorlagen-Nachtrag, F-10 die Extraktion einer Funktion mit bestehender
Testdeckung. F-11 bis F-14 sind billig und ohne Release-Risiko.

**Übergabe:** Die Findings gehen an die Implementation; die Finding-Klassen
gehören in die Slice-Closure und das Beobachtungs-Register (F-5/F-10 als
Fortschreibung der ADR-0054-Klasse, F-8 als Hinweis für die
Einführungs-Prozedur im Schwester-Repo: dessen Marker-Konvention weicht ab und
gehört vor dem ersten Lauf konfiguriert). Dieser Report ist ein Lauf-Beleg
(dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) und ersetzt keine
Verifikation.
