# Verifikation — slice-247: vcs-Modul meldet stilles Grün über leerer, auflösbaren Range

- **Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?"; Prüfgrundlage DoD + Spec + Plan
- **Gegenstand:** `slice-247` · Range `5a45b0b2..HEAD` (`5f19f081` Claim/Move · `f4522207`
  Implementierung · `cdc8ede8` R1-Report + R1-Einarbeitung · `a2219e1f` R1-Rest)
- **Datum:** 2026-09-29
- **Arbeitsbaum bei Verifikation:** sauber (`git status --porcelain` leer, vor und nach allen Läufen)
- **Sensor-Belege — alle selbst gefahren**, nicht übernommen: `make test` (gebrochen und
  intakt), `make gates` (Exit 0), `make history-range-guard` (zwei Proben), Blob-Vergleiche
  via `git rev-parse`, Diff-Prüfungen gegen `internal/hexagon/core/rules/commits*.go`,
  `tools/harness/history-range-guard.sh`, `Makefile`

---

## Ergebnis in einem Satz

DoD-Punkte 1–5 sind **bestätigt** — der Bewusst-Brech-Nachweis wurde **selbst wiederholt**
(ohne Fix rot an genau der fail-closed-Assertion, mit Fix grün). **Nicht erfüllt ist DoD 6**
(Closure-Notiz §7 leer, beide §6-Risiken ohne Ausgang, DoD-Häkchen ungesetzt) — die
Aufgaben-Premisse „§7 gefüllt, DoD-Häkchen gesetzt" ist gegen den Baum **falsch**; der
`git mv` nach `done/` ist solange blockiert.

---

## DoD-Punkt-für-Punkt

### 1. Neue `DC-FA-VCS-002` im Lastenheft + Spezifikation §DC-FA-VCS-002.a — BESTÄTIGT

- **Form (`MR-032`):** Versions-Bump `0.93.2` → `0.93.3` im Kopf (`spec/lastenheft.md:3`)
  und Historie-Zeile 0.93.3 (2026-09-29) in der Versions-Tabelle
  (`spec/lastenheft.md:3974`) — beides im selben Commit (`f4522207`) wie die Anforderung,
  vor der Closure. Lastenheft-Status (`Draft`) liegt unter `Accepted`, die Pflicht greift.
- **Keine Kollision:** `DC-FA-VCS-002` steht genau zweimal — Definition (Zeile 2195) und
  Historie-Verweis (Zeile 3974). Kein Alttreffer.
- **Aussage:** der Leerfall (aufgelöst, null Commits, inklusive Basis = Spitze,
  `rev-list --count base..head` = 0) ist als **laut zu meldender Zustand** (Exit ≠ 0,
  benannte Meldung) festgelegt; shallow-Klon `HEAD..HEAD` als Prüffall benannt;
  Out-of-Scope-Satz trägt `commits`, Vorlauf-Wächter und Range-Scoping (Zeilen 2205–2212).
- **Grund-Code-Form:** §DC-FA-VCS-002.a (`spec/spezifikation.md:1958`) — Leerfall-Prüfung
  **vor** dem Tree-Vergleich, `--staged`-Ausnahme, und (nach R1-M-3) die Shallow-Grenze
  des Vorfahren-Walks. Aussage und Code decken sich (siehe Punkt 2).

### 2. Fix umgesetzt + Bewusstes Brechen — BESTÄTIGT (Rot-Beweis selbst wiederholt)

- **Fix (`f4522207`, `internal/adapter/driven/git/git.go` +24 Zeilen):** `AllPaths` löst
  vor dem Tree-Vergleich beide Refs via `ResolveRevision` (Fehler ⇒ „nicht auflösbar",
  wie vor dem Fix Exit 2), bricht bei `*baseHash == *headHash` ab („Basis und Spitze
  benennen denselben Commit") und bei `baseAnc[*headHash]` (Spitze ∈ ancestors(Basis),
  „0 Commits, es wurde nichts geprüft"). Deckungsgleich mit §DC-FA-VCS-002.a.
- **Eigener Rot-Lauf:** Pre-Fix-Blob (`20059a01`, identisch an `5f19f081` und
  `f4522207^`) via `git checkout 5f19f081 -- internal/adapter/driven/git/git.go`
  eingesetzt, `make test` gefahren:
  `--- FAIL: TestAllPathsLeereRange` … `git_test.go:671: leere Range (Basis = Spitze)
  hätte fail-closed liefern müssen` — **genau die fail-closed-Assertion, der richtige
  Grund**; einziger Fehler der Suite. Danach Fix restauriert (Blob `1628200c` = HEAD,
  Status leer).
- **Grün-Hälfte:** `make test`/`make gates` auf HEAD: `ok … internal/adapter/driven/git`
  — `TestAllPathsLeereRange` (Basis = Spitze, invertierte Range, Normalfall) grün.
- **Anmerkung:** der Implementer-Beweis (Commit-Botschaft) nennt „git stash, Suite rot,
  Fix zurück"; ich habe über den Pre-Fix-Blob reproduziert — äquivalent, gleicher
  assertion-genauer Rot-Lauf. Die R1-Notiz (Report) und die Commit-Botschaft behaupten
  damit Belegbares.
- **R1-L-1-Abweichung sauber:** der Test fährt ein volles In-memory-Fixture (zwei
  Commits), nicht slice-245s shallow-Clone-Aufbau — im Plan §3 als `ABWEICHUNG`
  deklariert; der shallow-Fall ist über die deklarierte Spec-Grenze + den Wächter
  abgedeckt (nicht automatisiert). Deklariert, nicht still.

### 3. Regression (commits-Modul, history-range-guard) — BESTÄTIGT

- **commits-Modul:** `git diff 5a45b0b2..HEAD -- internal/hexagon/core/rules/commits.go
  internal/hexagon/core/rules/commits_test.go` **leer**. Vertragstest
  `TestCheckCommitsEmptyRange` (`commits_test.go:137`: leere Range ⇒ 0 Befunde, kein
  Fehler) unverändert; sein Paket `internal/hexagon/core/rules` lief in **beiden** Läufen
  (mit und ohne Adapter-Fix) grün — der Still-grün-Vertrag ist unangetastet, und die
  Plan-Abgrenzung (§1, nach R1-M-2 korrigiert: still grün auf **vollem** Fixture, laut nur
  an der Shallow-Grenze) stimmt jetzt mit dem Code überein.
- **history-range-guard:** `tools/harness/history-range-guard.sh` und `Makefile` im
  Range-Diff **leer**; das Skript ist reines bash + git („Kein Docker, kein Netz"),
  hat keine Produkt-Abhängigkeit. **Live-Probe:** `make history-range-guard
  RANGE=HEAD..HEAD` → laut: „ist aufloesbar, aber LEER (0 Commits)" mit
  fetch-depth-Hinweis, Exit ≠ 0; `RANGE=HEAD~1..HEAD` → „aufgeloest, 1 Commit(s) — OK",
  Exit 0. Verhalten wie in slice-245 (Probe D) — unverändert.

### 4. `make gates` grün — BESTÄTIGT

- Eigener Lauf nach Restaurierung: **Exit 0**, Abschlusszeile `[gates] baseline-verify +
  workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep +
  gate-consistency + planning-check green`, `record-gates` als letzter Schritt.
- Coverage-Gate: `OK — Coverage 94.60% erfüllt Schwelle 93%`.
- `git status --porcelain` nach dem Lauf: **leer**.

### 5. Review-Report liegt vor — BESTÄTIGT

`docs/reviews/2026-09-29-slice-247-vcs-leere-range-r1.md` (HIGH 0 · MEDIUM 3 · LOW 2 ·
INFO 1). Einarbeitung geprüft: R1-M-1 (commits-Spiegel, `spec/spezifikation.md`,
`cdc8ede8`), R1-M-2 (Plan-§1-Korrektur, `cdc8ede8`), R1-M-3 (Shallow-Grenze in
Lastenheft-Grenze + §DC-FA-VCS-002.a, `cdc8ede8`), R1-L-1 (Abweichung im Plan §3,
`cdc8ede8`/`a2219e1f`), R1-L-2 (`Verantwortlich:` gesetzt, `cdc8ede8`). I-1 ist als
Release-Prep-Restposten geführt (Dokumentations-Regel 17) — korrekt **nicht** im
Feature-Commit gepflegt.

### 6. Closure-Notiz mit Lerneintrag; jedes Risiko aus §6 mit Ausgang — NICHT ERFÜLLT

- §7 Closure-Notiz: alle sieben Felder `—` (Platzhalter), kein Lerneintrag.
- §6: beide Risiken tragen `**Ausgang:** *(offen)*` — weder *eingetreten* noch
  *entfallen* noch *weiter offen* (Regelwerk `modul-05` §Offene Risiken).
- Die DoD-Häkchen (§2) sind sämtlich ungesetzt (`- [ ]`).
- **Die Aufgaben-Premisse „§7 gefüllt, DoD-Häkchen gesetzt" ist damit gegen den
  Baum falsch** — die Behauptung läuft der Datei voraus. Das ist die Verifier-Klasse
  „Behauptung ohne Bestätigung", hier am Plan-Artefakt selbst.

**Folgerung:** Der Slice liegt korrekt in `in-progress/` — DoD-Häkchen, §7-Ausfüllung
und §6-Ausgänge sind **Bedingung** für den `git mv` nach `done/` (Baseline-Regelwerk
`modul-05` §Lifecycle), nicht seine Folge. Vor der Closure fehlen außerdem die
Nachtlauf-/Trigger-Angaben der §4-Zusage nur als §7-Inhalt, nicht als Blocker.

---

## Plan-vs-Code-Diff

| Plan §3 | Ist | Urteil |
|---|---|---|
| `spec/lastenheft.md` update | `DC-FA-VCS-002` + 0.93.3 + Historie-Zeile | konform |
| `spec/spezifikation.md` update | §DC-FA-VCS-002.a + commits-Spiegel-Korrektur (R1-M-1) | konform |
| `internal/adapter/driven/git` update | `AllPaths`-Leerfall-Zweig (+24 Zeilen) | konform |
| Testdatei neu/update | `TestAllPathsLeereRange` + `TestAllPathsUnlesbarerUnterbaum` auf zwei Commits umgestellt; ABWEICHUNG (volles Fixture) deklariert | konform |

- **Abgrenzung (§1) eingehalten:** Range-Diff umfasst genau sechs Dateien — die vier
  geplanten, den Slice-Plan selbst und den R1-Report (DoD 5, pro Slice konstant). Das
  `commits`-Modul, der history-range-guard, die Range-Semantik und fremde CI-Workflows
  sind unangetastet; keine stillsche Ausweitung, die Plan-Änderungen (R1) sind
  Review-Bearbeitung, keine Scope-Vergrößerung.
- **Commit-Traceability:** `f4522207` nennt `DC-FA-VCS-002` und `slice-247`; die
  übrigen Commits nennen `slice-247` — `trace-check`-Form erfüllt.

## Verbleibende Risiken / offene Punkte

1. **Closure-Schuld (DoD 6)** — §7, §6-Ausgänge, Häkchen: vor dem `git mv` nachzutragen;
   solange blockiert.
2. **I-1 Restposten:** Handbuch-Fehlerbild „Range-Leerfall …" bei der Release-Prep
   ergänzen (Dokumentations-Regel 17) — geführt, nicht vergessen werden.
3. Die shallow-seitige Inversion (Range leer, aber Fehlerpfad „Range-Basis-Vorfahren
   nicht lesbar" statt benannter Leerfall-Meldung) ist nur über die deklarierte Grenze
   abgedeckt, nicht automatisiert — R1-L-1/§3-Abweichung; die Gleichheits-Komponente
   fängt `HEAD..HEAD` vor jedem Walk.
4. Nebenbefund zur Übergabe: die im Auftrag genannten Hashes `e5787517`/`cf6f6908`
   existieren in der Historie nicht — die tatsächlichen Commits sind `5f19f081`,
   `f4522207`, `cdc8ede8`, `a2219e1f`. Ohne Sachwert, aber die Übergabe zitierte die
   Range nicht byte-genau.
