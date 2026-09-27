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

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** claude-sonnet-5. **Datum:** 2026-09-27.

---

## 1. Ziel und Abgrenzung

**Ziel:** Für `AGENTS.md` (oder die Sub-Area, die der Auftraggeber
bestimmt) eine tatsächliche Zeilen-/Größenschwelle aktivieren — über
`file[].max-lines` und/oder `structure[].max-lines` in `.d-check.yml` —,
damit das seit slice-231 benannte Wachstum der Briefing-Datei einen
Sensor **trägt**, nicht nur einen Sensor **hat**.

**Ausdrücklich NICHT in diesem Slice, bis zur Klärung:**

- **Die konkrete Schwelle selbst ist hier nicht vorentschieden** — eine
  Gate-Schwelle mit Vertragswirkung ist Auftraggeber-Entscheidung
  (`AGENTS.md` §3.6). Dieser Slice-Kopf hält nur die Kennung, an die das
  Beobachtungs-Register (`BEO-ALL/briefing-datei-ueberschreitet-lade-budget`)
  seinen `geplant`-Ausgang hängt (Baseline-Regelwerk `modul-06-roadmap.md`
  §Das Beobachtungs-Register, 3×-Schwelle erreicht mit slice-237).
- **Kein neuer Sensor-Code** — beide Werkzeuge (`file`, `structure`)
  existieren bereits
  ([ADR-0088](../../adr/0088-file-modul-groessengrenzen.md),
  [ADR-0089](../../adr/0089-structure-max-lines-zwoelfte-bedingung.md)); dieser
  Slice ist reine Konfiguration + ADR für die Schwellen-Entscheidung selbst.

## 2. Definition of Done

*(Wird bei Beanspruchung präzisiert — abhängig von der
Auftraggeber-Entscheidung zu Schwelle und Werkzeug-Wahl.)*

- [ ] Auftraggeber-Entscheidung eingeholt: Schwelle (Zeilen), Werkzeug
      (`file[].max-lines` ganze Datei vs. `structure[].max-lines` je
      Abschnitt), Ziel-Datei(en).
- [ ] `.d-check.yml` entsprechend erweitert, ADR mit der Begründung.
- [ ] `make gates` grün.

## 3. Plan (vor Code)

*(Noch nicht ausgearbeitet — abhängig von §2.)*

## 4. Trigger

**Start** (`open` → `next`): Auftraggeber priorisiert und entscheidet die
Schwelle.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

- Ohne Auftraggeber-Entscheidung bleibt dieser Slice in `open/` liegen —
  das ist der Grund, warum er existiert (Kennung für den Register-Ausgang),
  nicht ein Versäumnis.

## 7. Closure-Notiz

*(Bei der Closure zu füllen.)*
