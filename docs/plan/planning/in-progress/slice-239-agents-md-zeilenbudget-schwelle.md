# Slice slice-239: Eine Zeilenobergrenze für `AGENTS.md` aktivieren (Auftraggeber-Entscheidung)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** — **wellenlos** (voraussichtlich; abschließend erst bei
Beanspruchung zu beurteilen).

**Bezug:**
[`DC-FA-FILE-001`](../../../../spec/lastenheft.md#dc-fa-file-001--zeilen--und-byte-obergrenzen-einer-ganzen-datei-modul-file-opt-in),
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(beide Sensoren existieren bereits — dieser Slice aktiviert nur die
Konfiguration).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-27.

---

## 1. Ziel und Abgrenzung

**Ziel:** `AGENTS.md` unter eine tatsächliche Zeilenschwelle bringen und
diese über `file[].max-lines` aktivieren — Auftraggeber-Entscheid
2026-09-27: Schwelle **400 Zeilen**, Werkzeug `file[].max-lines`
(ganze Datei — passt zur Beobachtung, die die *geladene* Rohdatei meint,
nicht ihren bereinigten Fließtext), Ziel-Dateien **`AGENTS.md` und
`harness/README.md`** (beide laut §Leseordnung in jedem Lauf gelesen).
Da `AGENTS.md` die Schwelle heute überschreitet, trägt dieser Slice die
Kürzung selbst: Regel-Begründung, Durchsetzungs-Detail und Grenzen wandern
nach `harness/rules/<slug>.md` — derselbe Index-Muster wie
`harness/conventions.md`s Adaptions-Block (kurzer Eintrag im Briefing,
Volltext eine Datei weiter) —, während Überschrift und der operative
Kern (Falsch/Richtig, Begründung in einem Satz, Pointer) in `AGENTS.md`
stehen bleiben. Kein Wortlaut geht verloren, nur der Ort wechselt.

**Ausdrücklich NICHT in diesem Slice:**

- **`harness/README.md` wird nicht gekürzt** — mit 233 Zeilen liegt es
  bereits unter der Schwelle; nur die Gate-Aktivierung betrifft es.
- **Keine inhaltliche Änderung der ausgelagerten Regeln** — Zusage,
  Durchsetzung und Grenzen bleiben wortgleich, nur der Ort wechselt. Eine
  inhaltliche Schärfung wäre ein eigener Vorgang (eigene Begründung, eigener
  Commit).
- **`harness/conventions.md` bleibt unangetastet** — sein Adaptions-Block
  ist die Vorlage für das Muster, nicht sein Gegenstand.
- **Keine weiteren `structure`/`file`-Regeln über den Bestand hinaus** — nur
  die beiden hier entschiedenen Ziel-Dateien.

## 2. Definition of Done

- [x] [ADR-0096](../../adr/0096-agents-md-regel-auslagerung-harness-rules.md)
      dokumentiert das Auslagerungs-Muster (`harness/rules/<slug>.md`) und
      die Schwellen-Entscheidung (400 Zeilen, `file[].max-lines`, Scope
      `AGENTS.md` + `harness/README.md`).
- [x] `AGENTS.md` liegt unter 400 Zeilen (355); alle heute referenzierten
      `AGENTS.md#3.x`/`§5`-Anker bleiben aufgelöst (Überschriften bleiben
      stehen, nur die Abschnitts-Körper schrumpfen — bestätigt durch das
      grüne `links`/`anchors`/`ids`-Ergebnis über den ganzen Baum).
- [x] `.d-check.yml`: zwei `file[]`-Regeln (`AGENTS.md`, `harness/README.md`,
      je `max-lines: 400`); `make gates` grün.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/0096-agents-md-regel-auslagerung-harness-rules.md` | neu | Muster + Schwellen-Entscheidung, schärft [`DC-FA-FILE-001`](../../../../spec/lastenheft.md#dc-fa-file-001--zeilen--und-byte-obergrenzen-einer-ganzen-datei-modul-file-opt-in) nicht (reine Konfigurationsentscheidung) |
| `docs/plan/adr/README.md` | update | Index-Eintrag [ADR-0096](../../adr/0096-agents-md-regel-auslagerung-harness-rules.md) |
| `harness/rules/docker-make-only.md` | neu | Begründung/Durchsetzung/Grenzen aus §3.1 |
| `harness/rules/kommentare-fuenf-klassen.md` | neu | Begründung/Bestandsgrenze aus §3.7 |
| `harness/rules/github-action-sha-pin.md` | neu | Begründung/Grenzen aus §3.9 |
| `harness/rules/dokumentations-regeln/*.md` | neu | je Regel aus §5 eine Datei, `AGENTS.md` §5 wird zur Index-Tabelle |
| `AGENTS.md` | update | §3.1/§3.7/§3.9/§5 auf operativen Kern + Pointer gekürzt |
| `.d-check.yml` | update | `file[]`-Regeln für `AGENTS.md`, `harness/README.md`; `file` neu in `modules:` |
| `Makefile` | update | `FOCUS_DISABLE` um `--disable file` ergänzt (spiegelt `.d-check.yml` `modules:`, Makefile-Kommentar-Pflicht) |
| `internal/adapter/driven/configyaml/gate_consistency_test.go` | update | `file` in `netlessDocModules()` — sonst meldet `TestQA03_NetlessModuleList_Live` das neu aktivierte Modul als unklassifiziert |
| `.harness/skills/reviewer.md` | update | `d-check:cite AGENTS.md:366-367` neu geankert auf `281-281` — die Kürzung verschiebt die zitierte Zeile |

**Ansatz:** Die Auslagerung ist ein mechanisches Muster, pro Regel wiederholt
— Überschrift und Anker bleiben unverändert (keine Rewrite-Pflicht für die
~24 Dateien, die `AGENTS.md §3.x`/`§5` heute zitieren), nur der Abschnitts-
**Körper** wird auf den operativen Kern gekürzt und um einen Pointer auf die
neue Datei ergänzt.

## 4. Trigger

**Start** (`next` → `in-progress`): direkt beansprucht, Auftraggeber-
Entscheid liegt vor (2026-09-27).

**Rückführungen:**

- `in-progress` → `next`: entfällt — die Arbeit ist ein mechanisches Muster
  ohne Sub-Schnitt, der eine Zerlegung nahelegen würde.
- `in-progress` → `open`: falls die Anker-Stabilität sich als falsch
  herausstellt (Überschriften-Text müsste sich doch ändern) und eine
  Nachzugs-Kampagne über die ~24 zitierenden Dateien nötig würde — dann zurück
  zur Neu-Planung als eigener, größerer Slice.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

- Ein extrahierter Abschnitt verliert beim Kürzen operative Substanz, die ein
  Implementer beim schnellen Lesen von `AGENTS.md` gebraucht hätte, ohne den
  Pointer zu verfolgen. — **Ausgang:** entfallen: unabhängiger Review R1 hat
  alle 13 Volltext-Dateien Zeile für Zeile gegen die AGENTS.md-Vorfassung
  geprüft (12/13 wortgleich, zwei kleine Wortlaut-Lücken F-1/F-2 gefunden und
  behoben) — keine operative Substanz ging verloren.
- Ein Sensor (`ids`/`matrix`/`codepaths`) meldet einen der ~24 zitierenden
  Bestandsdateien, weil ein Anker sich doch verschiebt. — **Ausgang:**
  entfallen: `make gates` lief nach der Kürzung grün, kein Befund auf
  Bestandsdateien.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Anker-Stabilität trug wie geplant — weil
  jede `### 3.x`/`## 5`-Überschrift unverändert stehen blieb, brauchte keine
  der ~24 zitierenden Bestandsdateien einen Nachzug; `make gates` bestätigte
  das leer (kein Befund auf einer Bestandsdatei). Das Muster aus
  `harness/conventions.md`s Adaptions-Block (Index im Hauptdokument,
  Volltext eine Datei weiter) übertrug sich mechanisch auf `AGENTS.md` §3.x
  und §5.
- **Was ging anders als geplant:** Die Aktivierung des Moduls `file` in der
  `.d-check.yml`-`modules:`-Liste zog zwei ungeplante Mirror-Pflichten nach
  sich (`Makefile` `FOCUS_DISABLE`, `netlessDocModules()` in
  `gate_consistency_test.go`) — im Plan (§3) nachgetragen, bevor der Code
  geschrieben wurde, wie es der Fund verlangte, nicht danach.
- **Steering-Loop-Eintrag:** Modul `file` scharf geschaltet: `AGENTS.md` und
  `harness/README.md` tragen jetzt eine echte Zeilenobergrenze (400) statt
  eines bloß existierenden, nie aktivierten Sensors — liegt in
  `.d-check.yml` (Modul `file`). Auslöser:
  `BEO-ALL/briefing-datei-ueberschreitet-lade-budget` (slice-231,
  slice-236, slice-237 — 3×).
- **Beobachtungs-Register (`../observations/`):**
  `BEO-ALL/briefing-datei-ueberschreitet-lade-budget/state.md` von `geplant`
  auf `verkörpert` gesetzt (Zielort + Anker siehe oben).
- **Folge-Slices:** keine.
- **Risiken aus §6:** beide entfallen — siehe §6.
- **Drei Paarungen** (wellenlos, hier geprüft):
  - **Anker:** `.d-check.yml` trägt `seit slice-239` im Kommentar über dem
    `file:`-Block — geprüft.
  - **Folge-Slice:** keiner genannt — entfällt.
  - **Register:** `BEO-ALL/briefing-datei-ueberschreitet-lade-budget`
    existiert und trägt Belege (`evidence/`) — geprüft.
- **Review:** R1 (`docs/reviews/2026-09-27-slice-239-agents-md-regel-auslagerung-review-r1.md`),
  4 LOW, nicht merge-blockierend. F-1 (verlorener Halbsatz in der §5-Tabelle)
  und F-2 (uneinheitlich behandelte `(§N)`-Querverweise beim Verschieben)
  sowie F-3 (Plan-Tabelle ohne den `reviewer.md`-Nachzug) behoben. F-4
  (Fitness-Function-Tabelle in [ADR-0096](../../adr/0096-agents-md-regel-auslagerung-harness-rules.md)
  beschreibt die Aktivierung als `--enable file` statt über
  `.d-check.yml`s `modules:`-Liste) **bewusst nicht behoben** —
  [ADR-0096](../../adr/0096-agents-md-regel-auslagerung-harness-rules.md)
  ist bereits `Accepted` und committet; eine Korrektur bräche
  `make adr-check` (getestet: `core-drift-vcs`) und
  verlangte eine eigene Folge-ADR für eine reine Wortlaut-Nuance in einer
  Tabellenzelle. Unverhältnismäßig für den Befund — dokumentiert statt
  korrigiert.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:363-364 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: das
Harness-Werkzeug selbst (`*`, Kürzel `ALL`); bereits deklariert.

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:369-369 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:**
[`BEO-ALL/briefing-datei-ueberschreitet-lade-budget`](../observations/BEO-ALL/briefing-datei-ueberschreitet-lade-budget/observation.md)
ist der Auslöser dieses Slice selbst (3×, Ausgang `geplant`, jetzt
`verkörpert`).

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung, 2026-09-27):
`make nightly-state` — derselbe Stand wie bei slice-238 (unverändert seit
slice-232/233/235/238): `upstream-drift.yml` rot (Baseline-Hebung
`v6.9.0`→`v6.10.0` bewusst zurückgestellt bis `v6.11.0`), `image-scan.yml`
grün. Keiner der gemeldeten Punkte berührt dieses Harness-Werkzeug.

**Modus-Begründungsblock:** GF.

### Sub-Area: `*` (Harness-Werkzeug)

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — das Auslagerungs-Muster kopiert eine
  bestehende Konvention (`harness/conventions.md` §Adaptions-Block: „Diese
  Datei trägt den Index, nicht die Einträge").
- **Phase-Reife:** Phase 3.
- **Evidenz-/Diskrepanz-Risiko:** Niedrig — reine Verschiebung, kein neuer
  Bestand.
- **Reconciliation-Aufwand:** Keiner.
