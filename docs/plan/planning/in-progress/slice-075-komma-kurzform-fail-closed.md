# Slice slice-075: Komma-Kurzform fail-closed statt still verschluckt

**Status:** in-progress (welle-60-trace-cross-consistency, 2026-07-17).

**Welle:** aktiv (welle-60), **vorrangig** in Arbeit genommen: der Auftraggeber
(grid-gym) meldet den stillen Komma-Enum-Drop als **produktiv verdrahtetes**
`trace.coverage`-Problem — aktiv falsche Ergebnisse bei einem realen Konsumenten,
der einzige offene Punkt, der still falsche Zahlen liefert und dabei verdrahtet
ist.

**Bezug:** **Change Request** (neue fail-closed-Klasse = neues Akzeptanzkriterium):
[`DC-FA-COV-001`](../../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
(Lastenheft 0.46.0), über den geteilten Reader zugleich
[`DC-FA-XREF-001`](../../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in);
Algorithmus in
[`DC-FA-COV-001.a`](../../../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)
Schritt 3; begründende Entscheidung
[ADR-0041](../../adr/0041-komma-kurzform-fail-closed.md) (Proposed).
**SemVer-Minor** — Vertrags-Zuwachs, anders als die Defekt-Fixes
[slice-073](../done/slice-073-link-transparente-range-fortsetzung.md)/[slice-074](../open/slice-074-kommentar-suffix-tabellenzeilen.md).

**Autor:** pt9912. **Datum:** 2026-07-17.

---

## 1. Ziel

`GG-SCN-001, 007` deckt nur `GG-SCN-001` — `007` fällt **still** heraus. Die
Komma-Kurzform war nie zugesagt (Vertrag: `..`-Range und `/`-Aufzählung); der
Defekt ist nicht die fehlende Unterstützung, sondern dass der Autor **kein Signal**
bekommt. Der Slice macht die Gestalt fail-closed: Exit 2 mit Hinweis auf die
zugesagten Notationen — statt stillem Drop **oder** geratener Expansion.

## 2. Entscheidungen / Regel

- **Nur Komma + Ziffern triggert.** Ein Komma vor einer **vollständigen** Kennung
  (`GG-AR-COMP-CORE, GG-AR-COMP-DOMAIN`) ist keine Kurzform und bleibt unberührt —
  das ist in realen Design-Zellen der Normalfall, eine Regel die ihn bricht wäre
  unbrauchbar.
- **Kein Notations-Ausbau.** Komma-Enum zu unterstützen hieße raten
  (`GG-QA-001, 007 Sekunden`); die `id-pattern`-Validierung fängt breiten-gleiche
  Prosa-Zahlen nicht ([ADR-0041](../../adr/0041-komma-kurzform-fail-closed.md)).
- **Dieselbe Logik wie `AAA>BBB`:** die Syntax triggert, der Inhalt ist ungültig ⇒
  rot.

## 3. Definition of Done

- [x] **Lastenheft-CR:** neue fail-closed-Klasse + Negative-/Boundary-Kriterium,
  Version 0.46.0 und Historie.
- [x] **Spezifikation:** Komma-Kurzform in Schritt 3 + Historie.
- [x] **ADR + Index:** [ADR-0041](../../adr/0041-komma-kurzform-fail-closed.md),
  Proposed, im Index.
- [x] **Implementierung:** im geteilten Range-Parser (`expandRange`), Single-Check
  auf dem nach der Notation verbleibenden Rest — feuert **an allen drei
  Positionen** (nackte Kennung, hinter Range, hinter Enum), nach der
  Link-Transparenz. **Schärfung gegenüber der Erst-Formulierung:** der reale
  grid-gym-Fall ist `GG-SCN-001..005, 007, 008` (Range + Komma-Schwanz), nicht
  `GG-SCN-001, 007` — die enge Fassung hätte ihn verfehlt (§2.1 / ADR-0041
  Geschichte). Range-Expansion in `expandNumericRange` ausgelagert (nestif).
- [x] **Tests (positiv/negativ):** `TestExpandRangeCommaShortform` (table-driven,
  inkl. `..005, 007` / `/003, 007` / `..005, GG-QA-007` als Gegenprobe),
  `TestCoverageRefsCommaShortform` (einfach + **realer gemischter Fall** +
  `ranges: false` unberührt), `TestCrossConsistencyCommaShortform` (Bezug-Zelle).
- [x] **Mutations-Härte** (gemessen im Scratch-Worktree): Regel entfernt ⇒ 10
  Testzeilen kippen; `\d`-Bedingung entfernt (`^,\s*\d`→`^,`) ⇒ die
  „volle Kennung"/„Prosa"-Fälle kippen (auch der gemischte `..005, GG-QA-007`).
- [x] **Realdatenbeleg** (gegen das gebaute Image): `GG-SCN-001..005, 007, 008`
  (grid-gyms echte Form) ⇒ Exit 2 mit Notations-Hinweis; saubere Range und
  Range + volle Kennung ⇒ Exit 0.
- [x] **Nutzerdoku:** Handbuch §5 (die zwei zugesagten Notationen + die Komma-
  Fehlerklasse, einfach **und** hinter Range) + Historie 1.34; CHANGELOG 0.46.0.
- [ ] **Release:** v0.46.0 (Minor), Release-Prep **erledigt** (CHANGELOG,
  version.md, ghcr-Tag-Pins nachgezogen); Tag + GHCR-Push + Digest-Backfill
  stehen aus (Nutzer).
- [ ] **Qualität:** `make gates`/`make ci` grün (erreicht); unabhängiger,
  kontext-getrennter Review **vor** dem Release steht aus.

## 4. Risiken / offene Punkte

- **Neuer Falsch-Rot-Fall:** eine Prosa-Zahl hinter einer Kennung
  (`GG-QA-001, 2026`) läuft künftig auf Exit 2. Bewusst: laut und in Sekunden
  behebbar sticht still und unbemerkt. Gehört sichtbar in den CHANGELOG — es ist
  ein Vertrags-Zuwachs, kein Fix.
- **Der Konsument muss seine Quellen anfassen.** Das ist der Zweck, nicht der Preis
  — aber es ist Arbeit, die wir auslösen.
- **Dritte Notations-Kollision in Folge** (Link-Suffix, Kommentar-Suffix, Komma).
  Jede einzeln begründet, aber das Muster ist ein Signal: der lexikalische
  ID-Leser trifft laufend auf Markdown-/Prosa-Realität. Offener Punkt — gemeinsam
  mit dem gleichlautenden Risiko in
  [slice-074](../open/slice-074-kommentar-suffix-tabellenzeilen.md) §4: ab
  wann ist das die falsche Abstraktion?

## 5. Trigger

Konsumenten-Report grid-gym gegen v0.45.1 (2026-07-17), HIGH: „Komma-Aufzählung
wird still verschluckt und verfälscht das produktiv verdrahtete `trace.coverage`."
Nachgemessen: `GG-SCN-001, 007` ⇒ 1 Waise, `…/007` und `…..007` ⇒ 0. Der Konsument
ordnete selbst ein: die Kurzform ist out of spec, der stille Drop ist der Defekt —
„Stiller Drop ist die schlechteste der drei Optionen." Nutzer-Entscheid:
fail-closed statt Notations-Ausbau.

## 6. Sub-Area-Modus-Begründung

GF (Repo-Default): Der Vertrag führt, der Code folgt — Lastenheft-CR und
Spezifikation stehen vor der Implementierung.

## 7. Closure-Notiz (nach `done/`)

_Ausstehend — wird bei Abschluss mit Commit-Hash, Review-Verdikt und Lerneintrag
gefüllt._
