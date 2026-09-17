# Review-Report: slice-225 — 2026-09-17 (R1)

**Review-Art:** Code-/Plan-Review — geprüft wird die Implementierung der
neuen `structure`-Bedingung `open-tasks-require-marker`
(`internal/hexagon/core/rules/structure.go`), die begleitende
`ADR-0085`, die Spec-Erweiterungen (`spec/lastenheft.md` 0.87.0,
`spec/spezifikation.md` SPEC-083), die `CO-002`-Teilauflösung und die
Antwort auf den eingehenden CR von `ai-harness-init`, gegen den
Slice-Plan `slice-225`, `AGENTS.md` §3.4/§5, `harness/conventions.md`
und Baseline `v6.9.0` `modul-05-planning-harness.md`.

**Gegenstand:** `docs/plan/planning/in-progress/slice-225-gegenstand-entfallen-uebernommen.md`,
Commit-Kette `d0af7e12..834fcbb1` (`d0af7e12`, `b2874d18`, `ec057821`,
`da7205bc`, `834fcbb1`).

**Skill:** `.harness/skills/reviewer.md` @ v1.16.0 / 2026-09-07 ·
**Modell:** `claude-sonnet-5` · **Datum:** 2026-09-17

> **Zitier-Form.** Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb Kennung statt Adresse: `slice-NNN` statt Lifecycle-Pfad,
> `make <target>` als Token statt Link auf die Sensor-Datei, eine
> Baseline-Stelle als Tag plus Pfad in Inline-Code
> (`v6.9.0` · `regelwerk/<datei>.md` §<Abschnitt>).

**Eingangs-Kontext:**

- Slice-Plan `slice-225` (Fassung `in-progress/`, Stand des Prüfzeitpunkts,
  inkl. Plan-Revision `d0af7e12`)
- `ADR-0085` (neu, Accepted), `ADR-0074`, `ADR-0075` (Vorgänger-Muster
  des Moduls `structure`)
- `CO-002` (technisch aufgelöst, noch nicht nach `carveouts/done/`
  bewegt — laut Auftrag bewusst der Closure dieses Slice vorbehalten)
- eingehender CR `docs/plan/cr/2026-09-17-cr-eingehend-ai-harness-init-stilllegungs-form.md`
  (Adopter `ai-harness-init`), samt Antwort
- `DC-FA-STRUCT-001` (`spec/lastenheft.md` 0.87.0), `SPEC-083`
  (`spec/spezifikation.md`)
- `AGENTS.md` §3.4 (Referenzrichtung), §5 (Kommentar-/Zitat-/
  Commit-Disziplin), §6 (Workflow)
- Baseline `v6.9.0` · `modul-05-planning-harness.md` §Ein Slice, dessen
  Gegenstand ein anderer übernimmt; §Ziel-Form: Slice
- Beobachtungs-Register: `BEO-ALL/citation-stretched-beyond-scope`
  (vom Plan selbst als einschlägig geführt), `BEO-ALL/commit-message-overclaims-work`,
  `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`

---

## Verifikationen dieses Laufs

- `make test`: grün (alle Pakete, inkl. `internal/hexagon/core/rules` und
  `internal/adapter/driven/configyaml`).
- `make doc-check`: **791 Datei(en) geprüft, 0 Befund(e).**
- `make gates`: grün — `baseline-verify + workflow-pins + doc-check + lint +
  test + arch-check + coverage-gate + semgrep + gate-consistency +
  planning-check`; Coverage 94,60 % (Schwelle 93 %,
  `structureOpenTasks`/`hasMarker` je 100,0 %); Semgrep 0 Findings.
- `make verify-closure-notes`: **687 Datei(en) geprüft, 0 Befund(e)** —
  bestätigt `CO-002`s Verifikations-Häkchen 1+2 (Gate aktiviert, grün ohne
  `exempt-paths`) für den aktuellen Baumzustand.
- **Empirische Gegenprobe der Kernbehauptung** (nicht Teil der Test-Suite,
  gegen eine isolierte Fixture-Kopie außerhalb des Repos, mit derselben
  `.d-check.closure.yml`-Bedingung nachgebaut):

  ```yaml
  structure:
    - files: "docs/*.md"
      section-pattern: '^## [0-9]+\. Definition of Done'
      sections: each
      max-open-tasks: 0
      open-tasks-require-marker: "Gegenstand"
  ```

  ```markdown
  ## 2. Definition of Done

  - [ ] eins
  - [ ] zwei

  ## 3. Plan (vor Code)

  Text.

  ## 7. Closure-Notiz

  **Gegenstand:** entfallen: Beispiel.
  ```

  Ergebnis: `section-open-tasks-marker-missing` auf der Überschriftszeile
  von „## 2. Definition of Done" — **trotz** vorhandener, wohlgeformter
  `**Gegenstand:**`-Zeile in § 7. Siehe Finding R1-F-1.

## Findings

### R1-F-1 (HIGH) — Die Marken-Kopplung erkennt Baseline `v6.9.0`s eigene Ziel-Form nicht

- **Kategorie:** HIGH
- **Quelle:** Baseline `v6.9.0` · `regelwerk/modul-05-planning-harness.md`
  §Ein Slice, dessen Gegenstand ein anderer übernimmt (*„§7 trägt dann eine
  Zeile `Gegenstand:`"*); `DC-FA-STRUCT-001`
- **Pfad:** `internal/hexagon/core/rules/structure.go:372-395`
  (`structureOpenTasks`/`hasMarker`); `.d-check.closure.yml:226-238`
  (DoD-Regel, `section-pattern` bindet an „Definition of Done")
- **Befund:** `open-tasks-require-marker` prüft die Marke ausschließlich
  im Prosa-Text **derselben Section**, die `max-open-tasks` bereits
  scannt (`structureAmAbschnitt` ruft beide Bedingungen mit demselben
  `body` für dieselbe Überschrift auf). In diesem Repo ist das die
  DoD-Section (§2, `section-pattern: '^#{2,3} [0-9]+\. Definition of
  Done'`). Die Baseline — und der eingehende CR, der sie zitiert — verortet
  die Zeile `Gegenstand:` jedoch explizit in **§7 Closure-Notiz**, einer
  *anderen* Section. Ein Slice, der exakt der Baseline-Ziel-Form folgt
  (Marke nur in §7, DoD-Haken offen), erhält deshalb fälschlich
  `section-open-tasks-marker-missing` statt der beabsichtigten Erlaubnis —
  empirisch bestätigt (siehe Probe oben). Der einzige Grund, warum
  `slice-221` (der Anlassfall dieses Slice) trotzdem grün läuft, ist eine
  **zweite, nicht dokumentierte** `**Gegenstand:**`-Zeile, die direkt unter
  der DoD-Checkliste eingefügt wurde (`docs/plan/planning/done/slice-221-agents-md-tabellenzellen.md`,
  zwischen `## 2. Definition of Done` und `## 3. Plan (vor Code)`) — **neben**
  der kanonischen Zeile in `## 7. Closure-Notiz`. Weder `ADR-0085` noch
  `slice-225`s Plan noch die CR-Antwort nennen diese Ko-Lokations-Pflicht als
  Voraussetzung der Erkennung; `ADR-0085` §Konsequenzen (negativ) listet zwei
  andere Einschränkungen (Reichweite der Form-Erkennung, zwei Grund-Codes für
  denselben Rohzustand), aber nicht diese — die nach der Probe gewichtigste.
  Sogar `slice-225`s eigenes Risiko 1 (§6) geht implizit von einer Prüfung
  „direkt unter `## 7. Closure-Notiz`" aus (*„Die Form muss eng genug sein
  (z. B. nur als Zeilenanfang direkt unter `## 7. Closure-Notiz`)"*) — ein
  Beleg, dass der Plan selbst zum Schreibzeitpunkt von einem anderen
  Scan-Bereich ausging, als der Code tatsächlich prüft. Kein Test der neuen
  Suite (`TestOpenTasksRequireMarker_*`) platziert die Marke außerhalb der
  DoD-Section; die Lücke ist deshalb strukturell unsichtbar für
  Coverage-Metriken (100 % Zeilenabdeckung trotz fehlendem Szenario, weil
  „Marke fehlt im Abschnitt" und „Marke steht nur in einem anderen
  Abschnitt" denselben Codepfad durchlaufen).
- **Verifizierbar:** ja — reproduzierbar über die oben dokumentierte
  Fixture gegen `d-check:latest` (`--enable structure` mit der gezeigten
  `.d-check.yml`), oder als neuer Unit-Test gegen `structureOpenTasks`
  mit einer Marke außerhalb der geprüften Section.
- **Klasse:** `messmethode-scope-enger-als-dokumentierte-ziel-form`

### R1-F-2 (MEDIUM) — Commit-Botschaft zählt die neuen Tests falsch

- **Kategorie:** MEDIUM
- **Quelle:** `AGENTS.md` §5 (*„eine genannte Probe muss gelaufen sein …
  wer N Formen geprüft hat, berichtet N"*); `BEO-ALL/commit-message-overclaims-work`
- **Pfad:** Commit `b2874d18` (Botschaft, zweiter Absatz)
- **Befund:** Die Commit-Botschaft des Hauptcommits behauptet „Fünf neue
  Tests der Kopplung (Erlaubnis, Pflicht, Normalfall unberührt, abwesender
  Schlüssel) plus zwei Config-Ränder" — die Klammer nennt aber nur **vier**
  Zustände, und `internal/hexagon/core/rules/structure_offene_tasks_test.go`
  trägt genau **vier** neue `TestOpenTasksRequireMarker_*`-Funktionen, keine
  fünf. Die spätere CR-Antwort (`da7205bc`) und `spec/lastenheft.md`s
  Historie-Zeile (0.87.0, „Fünf neue **Akzeptanzkriterien**") zählen
  jeweils korrekt (vier Tests bzw. vier Kopplungs- plus ein
  Config-Akzeptanzkriterium = fünf) — die Diskrepanz sitzt ausschließlich
  im Feature-Commit selbst.
- **Verifizierbar:** ja — `grep -c '^func TestOpenTasksRequireMarker_' internal/hexagon/core/rules/structure_offene_tasks_test.go` liefert 4.
- **Klasse:** `commit-message-overclaims-work` (Zähl-Variante)

### R1-F-3 (MEDIUM) — Größenrisiko in §6 deckt nur eine der beiden Baseline-Achsen ab

- **Kategorie:** MEDIUM
- **Quelle:** Baseline `v6.9.0` · `regelwerk/modul-05-planning-harness.md`
  §Ziel-Form: Slice (*„Zu groß, wenn eines zutrifft: mehr als drei
  Liefer-Punkte … · mehrere Schichten betroffen … · nicht in einer
  Review-Sitzung prüfbar"*); `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`
- **Pfad:** `docs/plan/planning/in-progress/slice-225-gegenstand-entfallen-uebernommen.md:208-213`
  (drittes Risiko in §6)
- **Befund:** Das Risiko „Die CR-Aufnahme könnte den Zuschnitt sprengen"
  argumentiert ausschließlich über die Liefer-Punkte-Zahl und die
  Review-Sitzungs-Prüfbarkeit. Die dritte, unabhängige Baseline-Achse
  („mehrere Schichten betroffen") bleibt unerwähnt, obwohl DoD (1)
  tatsächlich vier Go-Pakete über drei Hexagon-Schichten bündelt (Core:
  `hexagon/core/rules`, `hexagon/core/model`, `hexagon/core/app`; Driven
  Adapter: `adapter/driven/configyaml`; Driving Adapter:
  `adapter/driving/cli`). Für dieses Modul ist das nach eigener Aussage des
  Plans „etabliertes Muster" (vergleichbar mit `ADR-0074`/`ADR-0075`), was
  die Bündelung in der Sache rechtfertigt — der Risiko-Eintrag benennt
  diese Rechtfertigung aber nicht gegen die einschlägige Baseline-Achse,
  sondern lässt sie ganz aus.
- **Verifizierbar:** nein (Urteilsfrage, kein Gate) — nachvollziehbar über
  `git show b2874d18 --stat` gegen die Paketgrenzen.
- **Klasse:** `grenzen-liste-wird-als-vollstaendig-gelesen`

### R1-F-4 (LOW) — Vorprüfung nennt die Register-Verzeichniszahl leicht falsch

- **Kategorie:** LOW
- **Quelle:** Baseline `v6.9.0` · `regelwerk/modul-05-planning-harness.md`
  §Zwei Schritte vor der Modus-Begründung (Sichtungs-Schritt)
- **Pfad:** `docs/plan/planning/in-progress/slice-225-gegenstand-entfallen-uebernommen.md:257`
  (§8, „gemergter Stand, **39** Verzeichnisse")
- **Befund:** `find docs/plan/planning/observations -mindepth 2 -maxdepth 2
  -type d | wc -l` liefert zum Prüfzeitpunkt **40** (39 unter `BEO-ALL/`,
  1 unter `BEO-HARN/`), nicht 39 — alle 40 Verzeichnisse datieren laut
  `git log` vor dem 2026-09-08, also lange vor der Sichtung. Ändert nichts
  am Ergebnis (derselbe einzige einschlägige Eintrag bleibt korrekt
  benannt), ist aber ein Zähler, der stimmen sollte, wo er zitiert wird.
- **Verifizierbar:** ja — s.o.
- **Klasse:** `zaehlmethode-misst-proxy-statt-gegenstand` (Randfall: reiner
  Zählfehler, kein Proxy)

## Negativbefunde

- **Hexagon-Import-Richtung** (`ADR-0005`): geprüft — Go-Code bleibt in
  seiner Schicht (`configyaml.go` parst, `structure.go` wertet aus, kein
  Cross-Layer-Import). Kein Befund.
- **Inline-Suppression / Gate-Lockerung ohne ADR** (`AGENTS.md` §3.2/§3.6):
  geprüft — kein `//nolint`, keine Schwellen-Senkung. Kein Befund.
- **Netzzugriff außerhalb `external`** (`DC-QA-03`): geprüft — reine
  Datei-/Konfigurationslogik. Kein Befund.
- **Kommentar-Klassen** (`AGENTS.md` §3.7, fünf Klassen): geprüft in
  `structure.go`, `config.go`, `finding.go`, `configyaml.go` — alle neuen
  Kommentare tragen Zusage/Kopplung/Grenze-Anker (`ZUSAGE`, `KOPPLUNG`,
  `GRENZE`), keine Review-Historie oder Herkunfts-Prosa außer dem einen
  zulässigen `(ADR-0085)`-Feld. Kein Befund.
- **Referenzrichtung/Provenance-Marker-Ehrlichkeit** (`AGENTS.md` §3.4,
  `DC-FA-MTX-003`): geprüft — alle `slice-225`/`slice-221`-Nennungen in
  `ADR-0085`, `spec/lastenheft.md` (§7 Historie) und
  `spec/spezifikation.md` (§7 Historie) tragen
  `<!-- d-check:status-provenance -->` und zeigen ausschließlich Herkunft
  (Träger/Fundort), nie eine Entscheidungsgrundlage. `make doc-check` (0
  Befunde) bestätigt zusätzlich, dass keine unmarkierte Abwärts-Referenz
  besteht. Kein Befund.
- **Zitier-Geltungsbereich der `d-check:cite`-Direktiven** (`MR-054`,
  `AGENTS.md` §5): geprüft — beide Zeilenspannen
  (`modul-05-planning-harness.md:363-364` und `:369-369`) lösen exakt auf
  das zitierte Wort auf. Kein Befund.
- **Grund-Code-Tabellen-Reihenfolge** (`spec/spezifikation.md` §4): geprüft
  — `SPEC-083` steht zwischen `SPEC-078` und `SPEC-079`, aber die Tabelle
  ist ohnehin modul-gruppiert statt streng numerisch (`SPEC-069`, `SPEC-059`
  u. a. stehen ebenfalls außerhalb ihrer Nummernfolge); die Platzierung
  neben den anderen `structure`-Codes entspricht dem Bestand. Kein Befund.
- **DoD-Bruch-Test-Deckung der vier behaupteten Zustände** (DoD (3)):
  geprüft gegen `structure_offene_tasks_test.go` — alle vier Zustände
  (Erlaubnis, Pflicht, Normalfall, abwesender Schlüssel) sind tatsächlich
  als Tests vorhanden und grün. Kein Befund (unabhängig von R1-F-1, das
  einen von der Suite nicht abgedeckten *fünften* Zustand betrifft, den
  keine der beiden CR-Parteien als eigenen Test benannt hatte).
- **CO-002-Reihenfolgeentscheidung** (Move erst bei `slice-225`s eigener
  Closure): geprüft gegen `modul-05-planning-harness.md` §Lifecycle als
  State Machine — zulässig, reiner `git mv` bleibt für die Closure
  reserviert; nicht Gegenstand dieses Reviews ist, ob DoD (2)s Wortlaut
  „der Carveout liegt in `docs/plan/carveouts/done/`" mit dem aktuellen
  (noch nicht bewegten) Baumzustand übereinstimmt — das ist DoD-/
  Verifikations-Territorium, nicht Plan-/ADR-Review (Anti-Pattern „Kein
  Verifier").
- **`gofmt`-Formatierung** (`internal/hexagon/core/app/diagnose.go`):
  geprüft — `gofmt -l` markiert die Datei (die neue `reasonTexts`-Zeile
  bricht die Spalten-Ausrichtung des Map-Literals). Kein Konventions-Anker
  bindet `gofmt`/`gofumpt` an dieses Repo (`.golangci.yml` führt keinen
  Formatter-Abschnitt), daher laut Anti-Pattern „Formatierung ohne
  Konventions-Anker ist kein Finding" **kein Finding**, hier nur der
  Vollständigkeit halber vermerkt.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 0 |

## Verdikt

**Blockiert.** R1-F-1 (HIGH) ist vor Closure zu beheben: Der Kern-Zweck
dieses Slice — Baseline `v6.9.0`s vierten Lifecycle-Zweig „Gegenstand
entfallen/übernommen" gate-tragfähig zu machen — wird für die
Baseline-kanonische Form (Marke ausschließlich in §7 Closure-Notiz) derzeit
**nicht** eingelöst; nur eine unbelegte, an keiner Stelle vorgeschriebene
Doppel-Platzierung der Marke (zusätzlich direkt unter der DoD-Checkliste)
macht `slice-221` zufällig grün. Vor der nächsten Closure-Runde ist zu
entscheiden und zu dokumentieren, welche der beiden Reparaturen greift:
(a) die Erkennung auf einen zweiten, benannten Abschnitt (z. B. „Closure-
Notiz") ausdehnen, oder (b) die Ko-Lokations-Pflicht (Marke **im selben
Abschnitt** wie die offenen Task-Items) explizit als Repo-Konvention
festschreiben — in `ADR-0085`, `spec/lastenheft.md`/`spezifikation.md`,
der `--print-config`-Doku und, falls zutreffend, in der CR-Antwort an
`ai-harness-init` (die Baseline-Wortlaut zitiert, den der Code so nicht
einlöst). R1-F-2 bis R1-F-4 sind kleinere Korrekturen, die die Closure
nicht für sich allein aufhalten.
