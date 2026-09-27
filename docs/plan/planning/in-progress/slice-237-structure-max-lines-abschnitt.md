# Slice slice-237: `structure` — zwölfte Bedingung `max-lines` (Zeilenbudget eines Abschnitts)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** — **wellenlos**. Die Closure-Bedingung geht über die eigene DoD
nicht hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine
Welle braucht).

**Bezug:**
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(Erweiterung — kein neues Kürzel, nach dem etablierten Schnitt-Kriterium:
Einzelmodul-Frage). Eingehender Change Request des Konsumenten
`ai-harness-course` (2026-09-27, Priorität niedrig, additiv).

**Berührte Spec-Stellen:**
[`DC-FA-STRUCT-001.a`](../../../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
(Schritt der Prosa-Bedingungen).

**Verantwortlich:** claude-sonnet-5.

**Autor:** claude-sonnet-5. **Datum:** 2026-09-27.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** `structure` bekommt eine Zeilenobergrenze `max-lines` je Abschnitt
— eine Regel meldet `section-lines-exceeded`, wenn der **bereinigte**
Abschnittstext mehr Zeilen trägt, als die Regel erlaubt.

**Der CR-Zähler ist falsch und wird korrigiert.** Der CR nennt die neue
Bedingung „zehnte" und zählt neun bestehende, indem er `max-tasks`/
`max-open-tasks` als einen Slot bündelt und `open-tasks-require-marker`/
`open-tasks-require-marker-section` (elfte Bedingung, Lastenheft 0.87.0)
gar nicht führt. Die **etablierte Ordinal-Zählung des Lastenhefts selbst**
(jede Erweiterung nennt ihre Nummer in der Historie — achte: `headings-match`
0.64.0, neunte: `cell-max-chars` 0.72.0, zehnte: `max-open-tasks` 0.79.0
unbeziffert, elfte: `open-tasks-require-marker` 0.87.0) macht `max-lines` zur
**zwölften**. Dieser Slice übernimmt die eigene Zählung, nicht die des CR.

**Der Out-of-Scope-Verweis des CR trägt nur zum Teil.** Der CR zitiert die
bestehende Out-of-Scope-Zeile „eine Stichtags-Regel …, die die
Kennungs-Konvention des Adopters interpretieren müsste" als Beleg dafür,
dass eine Alters-/Frische-Mechanik ausgeschlossen ist. Der zitierte Satz
schließt tatsächlich etwas anderes aus — einen ID-Stichtag („erst ab Kennung
N"), nicht die Frage „wann zuletzt geprüft". Der Schluss des CR (keine
Alters-Mechanik gewünscht) bleibt richtig, nur der Beleg dafür ist ungenau;
dieser Slice zitiert die Zeile nicht als Beleg für die eigene Abgrenzung
weiter unten.

**Eine Grenze, die der CR nicht benennt, gehört in den Vertrag: Fenced-Code
zählt nicht mit.** `max-lines` misst auf demselben bereinigten Text, auf dem
`non-empty`/`min-sentences` bereits operieren (`SectionProse`,
`internal/hexagon/core/rules/sections.go`) — und der ist über
`PreprocessMarkdown` gebildet, das Fenced-Code-**Zeilen vollständig
entfernt**, nicht nur maskiert. Ein Abschnitt mit einem langen Beispiel in
einem Codeblock wächst dadurch nicht gegen `max-lines`. Für den Anlass
(Prosa-Wachstum einer Hard-Rule-Datei) ist das die richtige Eigenschaft, aber
sie ist nicht selbsterklärend und gehört als **benannte Grenze** in
Lastenheft und Spezifikation, nicht nur in den Code.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Stichtags- oder Alters-Mechanik** — vom CR selbst ausgeschlossen und
  vom Kanon (Out-of-Scope) als zweiter Regel-Interpreter ohnehin verboten.
- **Eine Aussage über den Inhalt der Zeilen** — reine Struktur-Invariante wie
  jede andere Bedingung des Moduls.
- **Eine Schwelle für `AGENTS.md` §Harte Regeln oder einen anderen konkreten
  Abschnitt dieses Repos** — eine Schwelle ist eine Auftraggeber-Entscheidung
  mit Vertragswirkung (`AGENTS.md` §3.6) und gehört in einen Folge-Slice, der
  noch keine Kennung hat.
- **Eine Umbenennung von `file[].max-lines`** (Modul `file`, slice-236) — die
  beiden Schlüssel leben unter verschiedenen Top-Level-Blöcken (`structure:`
  vs. `file:`) und kollidieren syntaktisch nicht; die Spezifikation nennt die
  Verwechslungsgefahr trotzdem ausdrücklich, damit ein Leser beide nicht für
  denselben Mechanismus hält (`structure[].max-lines` zählt den bereinigten
  Text eines **Abschnitts**, `file[].max-lines` die rohen Zeilen einer
  **ganzen** Datei).
- **`tasks-ignore-pattern`-artige Teilausnahme** — der CR verlangt keine, und
  eine erklärte Teilmenge bräuchte einen eigenen Anlass (wie bei `max-tasks`).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — höchstens drei Liefer-Punkte.

- [ ] **Vertrag:** `spec/lastenheft.md` (Struktur-Anforderung: Bedingungs-Tabelle
      um `max-lines` (int ≥ 1) ⇒ `section-lines-exceeded` erweitert, drei
      Akzeptanzkriterien aus dem Break-Test des CR — Grenzwert N/N+1,
      fehlender Abschnitt, `sections: one` mit mehreren Treffern —, die
      Fenced-Code-Grenze benannt, Versions-Bump mit Historie-Zeile
      (**zwölfte Bedingung**, nicht „zehnte") nach
      [MR-032](../../../../harness/conventions/MR-032-historie-vor-accepted.md)),
      `spec/spezifikation.md` (Schritt der Prosa-Bedingungen um `max-lines`,
      Schema-Zeile `structure[].max-lines`, ein neuer Grund-Code
      `section-lines-exceeded`, die Fenced-Code-Grenze, die Abgrenzung zu
      `file[].max-lines`) und eine neue ADR samt Index-Eintrag. **Form nach
      dem Kanon:** Historie-Zeile nennt weder ADR noch Slice, Spezifikation in
      keinem Abschnitt; ADR trägt `Schärft:` aufwärts, mindestens drei
      verglichene Alternativen (u. a. `forbid-pattern`-Wiederholungs-Muster,
      eine Zeichen- statt Zeilen-Zählung, „nichts tun"), eine Fitness Function,
      einen `Re-Evaluierungs-Trigger`.
- [ ] **Umsetzung:** `MaxLines *int` an `model.StructureRule`, Validierung am
      Config-Rand (`max-lines` explizit < 1 ⇒ Exit 2 — Untergrenze 1, nicht 0:
      ein Abschnitt hat immer mindestens die Überschriftenzeile im rohen Text,
      aber der bereinigte Body kann bei `max-lines: 0` nie befundfrei sein,
      was die Bedingung praktisch zum Verbot des Abschnitts machte; zu klären
      **vor** dem Code, in der ADR festgehalten), ein neuer Zweig in
      `structureConditions` (`internal/hexagon/core/rules/structure.go`),
      Grund-Code `ReasonSectionLinesExceeded` = `section-lines-exceeded`.
      Tests: Grenzwert N/N+1, Fenced-Code zählt **nicht** mit (Rot-Beleg:
      derselbe Abschnitt ohne Fenced-Block meldet, mit Fenced-Block derselben
      Zeilenzahl nicht), fehlender Abschnitt (bestehendes `section-missing`,
      keine neue Form), `sections: one` mit mehreren Treffern (bestehendes
      `section-ambiguous`, keine Messung), `hint` gewinnt gegen die
      modul-eigene Meldung, ohne den Schlüssel byte-identisch
      ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- [ ] **Spiegel:** `AllReasons()`/`reasonTexts()` (`internal/hexagon/core/app/diagnose.go`)
      um den neuen Grund-Code; Spiegel-Liste in §3 **vor** dem Editieren, was
      Release-Prep ist (Handbuch, CHANGELOG, README) steht dort benannt statt
      übergangen.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine
      Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      wellenlos hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | Bedingungs-Tabelle, AKs, Bump + Historie (eigener Commit vor dem Code) |
| `spec/spezifikation.md` | update | Prosa-Bedingungen-Schritt, Schema-Zeile, Grund-Code, zwei Grenzen (Fenced-Code, Abgrenzung zu `file`) |
| `docs/plan/adr/` (neu) + `README.md` | neu / update | Untergrenze 1 statt 0, Alternativen, `max-lines`-Namensteilung mit `file` |
| `internal/hexagon/core/model/config.go` | update | `StructureRule.MaxLines *int` |
| `internal/adapter/driven/configyaml/configyaml.go` | update | `rawStructure.MaxLines`, Validierung in `structureBedingungsFehler` |
| `internal/hexagon/core/rules/structure.go` | update | `ReasonSectionLinesExceeded`, Zweig in `structureConditions` |
| `internal/hexagon/core/app/diagnose.go` | update | `AllReasons()`/`reasonTexts()` |
| Tests (`structure_test.go` oder neue Datei) | neu | Fälle aus §2 |

**Reihenfolge der Vertrags-Änderung.** Das Lastenheft steht auf `Draft`: der
Kanon (`grundlagen-source-precedence.md`, *Wann die CR-Pflicht beginnt*) lässt
es vor `Accepted` frei änderbar; [MR-032](../../../../harness/conventions/MR-032-historie-vor-accepted.md)
verlangt Bump und Historie-Zeile trotzdem. Die Änderung darf im Slice liegen;
sie steht in einem eigenen Commit vor dem Code.

**Vor dem Editieren — Spiegel listen**
([MR-025](../../../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)):
die Grund-Code-Tabelle in `spec/spezifikation.md` §4 (Vollständigkeitstest
`TestAllReasonsDeckungGegenSpezifikationGrundCodes` hält sie ohnehin gate-tragfähig)
und die Bedingungs-Tabelle der Struktur-Anforderung (siehe Bezug); Handbuch-Referenzabschnitt
§4.x je Bedingung und Modul-Übersichtstabelle sind Release-Prep.

## 4. Trigger

**Start** (`open` → `in-progress`): Auftraggeber-Priorisierung 2026-09-27 —
dieser Slice startet nach `slice-236` und vor `slice-232`/`-233`/`-234`/`-235`,
die vor ihm in der Warteschlange standen. Bumpt das Lastenheft wie mehrere
andere in der Warteschlange; läuft nacheinander (WIP-Limit 1), kein
Versions-Bump-Konflikt. Bei der Beanspruchung entsteht der dritte
Vorprüfungs-Block (Nachtlauf-Stand).

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): wenn die Untergrenze von `max-lines`
  (1 vs. 0) oder die Fenced-Code-Grenze eine Auftraggeber-Entscheidung
  braucht, die den Umfang der ADR sprengt.
- `in-progress` → `open` (blockiert): wenn eine Schwellen-Entscheidung für
  eine konkrete Datei vorgezogen werden soll und diesen Slice damit
  verquickt.

## 5. Closure-Trigger

DoD vollständig (§2) + Review-Report liegt vor + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Die Untergrenze `max-lines >= 1` ist eine Design-Entscheidung dieses Plans,
  keine Vorgabe des CR. **Ausgang:** bei Closure zu vergeben (die ADR trägt
  die Begründung; ein anderer Wert bliebe Review-Sache).
- Die Fenced-Code-Ausnahme kann bei einem Abschnitt mit viel Beispielcode und
  wenig Prosa dazu führen, dass `max-lines` nie greift, obwohl die Datei
  insgesamt wächst — das deckt eher `file[].max-lines` (slice-236) ab.
  **Ausgang:** bei Closure zu vergeben (benannte Grenze, kein Defekt).
- Zwei gleich benannte Schlüssel (`structure[].max-lines`,
  `file[].max-lines`) mit unterschiedlicher Semantik sind ein
  Verwechslungsrisiko in Doku und Support. **Ausgang:** bei Closure zu
  vergeben (die Abgrenzung steht im Vertrag, ist aber keine Sensor-Zusage).

## 7. Closure-Notiz

*(Bei der Closure zu füllen.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:363-364 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: das Produkt
(`spec/`, `internal/`) unter dem Default `*` (Kürzel `ALL`); bereits
deklariert, keine Ausdifferenzierung nötig.

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:369-369 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(Stichworte `structure`, Zeilenbudget, Ratchet, Modul-Wachstum,
`briefing-datei`): Treffer
[`briefing-datei-ueberschreitet-lade-budget`](../observations/BEO-ALL/briefing-datei-ueberschreitet-lade-budget/observation.md)
(2× nach `slice-236`s Closure) — derselbe Anlass, hier als Abschnitts- statt
Datei-Grenze. Kein weiterer Treffer.

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung, 2026-09-27):
`make nightly-state` meldet `upstream-drift.yml` **ROT** (Lauf
2026-09-27T06:03Z, unverändert seit `slice-236`), `image-scan.yml` **grün**.
Vier planmäßige Fremd-Release-Meldungen (`golangci-lint`, `semgrep`,
`a-check` VERALTET; `golang`-Basis-Digest ABWEICHEND), keine unerwarteten;
sie berühren dieses Modul nicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

### Sub-Area: `*` (Produkt: Spezifikation und Code)

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — Lastenheft, Spezifikation und ADR sind
  Kanon; die Bedingungs-Tabelle und die Ordinal-Zählung sind etablierte Form.
- **Phase-Reife:** Phase 5.
- **Evidenz-/Diskrepanz-Risiko:** Niedrig — additiv, ohne den Schlüssel
  byte-identisch; das Risiko liegt in der Namens-Doppelung mit `file` (§6).
- **Reconciliation-Aufwand:** Keiner — kein Brownfield-Bestand.
