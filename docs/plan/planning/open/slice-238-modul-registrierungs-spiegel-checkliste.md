# Slice slice-238: Ein Sensor oder eine vollständige Checkliste für die Spiegel der Modul-Registrierung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** — **wellenlos**. Die Closure-Bedingung geht über die eigene DoD
nicht hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine
Welle braucht).

**Bezug:** — kein `DC-*`, keine aktive ADR (dieser Slice entscheidet erst,
ob eine neue Anforderung oder nur ein Harness-Werkzeug entsteht). Folge-Slice
der Beobachtung
[`BEO-ALL/modulliste-spiegel-ungegated`](../observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md),
die mit `slice-236` ihre dritte Evidenz erreicht hat.

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** claude-sonnet-5. **Datum:** 2026-09-27.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Wer ein Regelmodul in `validModules()` aufnimmt, bekommt eine
**vollständige, an einer Stelle geführte** Liste der Orte, die mitgezogen
werden müssen — mechanisch geprüft, wo das trägt, sonst als Checkliste im
Workflow.

**Befund, der diesen Slice auslöst.** Die Beobachtung
[`BEO-ALL/modulliste-spiegel-ungegated`](../observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md)
hat mit `slice-236` ihre dritte Evidenz erreicht. Die drei Vorgänge zeigen
**unterschiedliche** Fundorte: `slice-115`/`slice-152` die bereits benannten
Nebenstellen (`FOCUS_DISABLE`, der [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Netzlos-Test, Gate-Doku-Prosa);
`slice-236` drei **weitere**, dort nicht genannte (`spec/lastenheft.md` §3
Bereichskürzel-Liste, §6 Glossar-Zeile „Regelmodul", `docs/user/operations.md`
Options-Tabelle). Der bestehende Vollständigkeitstest
(`TestAllReasonsDeckungGegenSpezifikationGrundCodes`,
`internal/hexagon/core/app/diagnose_test.go`) deckt nur die **Grund-Code**-Achse,
nicht die **Modulnamen**-Achse. `TestPrintConfigVerfuegbarDecktRegistry`
deckt eine einzelne Doku-Zeile, kein repo-weites Muster.

**Zwei mögliche Antworten, zu entscheiden vor dem Code:**

- **(A) Ein Vollständigkeitstest**, analog zu
  `TestPrintConfigVerfuegbarDecktRegistry` und
  `TestAllReasonsDeckungGegenSpezifikationGrundCodes`: für jede bekannte
  Fundort-Datei eine `grep`-artige Prüfung, dass jeder Name aus
  `model.ValidModules()` dort vorkommt. Trägt für Prosa-Listen (Glossar,
  Bereichskürzel, README-Bullets) nur, wenn die Namen dort **wörtlich** und
  **vollständig** stehen — keine Abkürzung, keine Teilmenge nach Sinn
  (`external`/`sources`/`vcs` sind in mehreren Listen bewusst ausgenommen und
  bräuchten eine Ausnahme-Liste).
- **(B) Eine Checkliste im `AGENTS.md`/Reviewer-Skill**, die die vollständige
  Fundort-Liste nennt und beim Aufnehmen eines Moduls durchgegangen wird —
  kein Sensor, `inferential feedforward` wie die drei Vorprüfungen eines
  Slice-Plans.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Rückwirkende Vollständigkeit für alle 25 Module** — dieser Slice baut den
  Mechanismus bzw. die Checkliste; ob der heutige Bestand an jeder Stelle
  bereits vollständig ist, ist eine separate Bestandsaufnahme (dieser Slice
  hat `file` an den drei gefundenen Stellen bereits nachgezogen, das reicht
  als Startpunkt).
- **Eine automatische Generierung** der Prosa-Listen aus `validModules()` —
  größerer Eingriff in mehrere Dateien, eigener Anlass nötig.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — höchstens drei Liefer-Punkte.

- [ ] **Entscheidung (A) oder (B)** getroffen und begründet — bei (A) reicht
      eine ADR-lose Fitness-Function-Ergänzung (kein neues Modul, kein neues
      `DC-*`), bei (B) eine `AGENTS.md`-Ergänzung mit Herkunfts-Anker.
- [ ] **Vollständige Fundort-Liste** einmal zusammengetragen (Diese Datei,
      `harness/README.md`, `AGENTS.md`, `spec/lastenheft.md` (§3, §6),
      `spec/spezifikation.md` §2/§4, `docs/user/*.md`, `README*.md`,
      `internal/adapter/driving/cli/config_template.go`,
      `internal/hexagon/core/app/diagnose.go`, `internal/hexagon/core/model/config.go`,
      `.d-check.yml`, `Makefile` `FOCUS_DISABLE`) — geprüft durch tatsächliches
      Lesen, nicht durch Wiederholung dieser Aufzählung.
- [ ] Umsetzung nach der getroffenen Entscheidung.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — der
      Ausgang von `BEO-ALL/modulliste-spiegel-ungegated` (heute `geplant:
      slice-238`) wird auf `verkörpert` gesetzt.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen sind getragen — wellenlos hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| noch offen | — | abhängig von Entscheidung (A)/(B) in §2 |

## 4. Trigger

**Start** (`open` → `in-progress`): keine Abhängigkeit; niedrige Priorität
(Harness-Werkzeug, kein Produkt-Vertrag). Bei der Beanspruchung entsteht der
dritte Vorprüfungs-Block (Nachtlauf-Stand).

**Rückführungen — vorab benennen:**

- `in-progress` → `next`: wenn (A) sich als so aufwändig erweist, dass es in
  mehrere Regeln zerfällt (eine je Fundort-Art).
- `in-progress` → `open`: wenn die Entscheidung (A)/(B) eine
  Auftraggeber-Präferenz braucht.

## 5. Closure-Trigger

DoD vollständig (§2) + Review-Report liegt vor + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Eine wörtliche Vollständigkeitsprüfung (A) kann durch bewusste Ausnahmen
  (`external`/`sources`/`vcs` fehlen in manchen Listen mit Absicht) selbst
  aufwändig werden. **Ausgang:** bei Closure zu vergeben.

## 7. Closure-Notiz

*(Bei der Closure zu füllen.)*

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
[`BEO-ALL/modulliste-spiegel-ungegated`](../observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md)
ist der Auslöser dieses Slice selbst (3×, Ausgang `geplant`).

**Modus-Begründungsblock:** GF.

### Sub-Area: `*` (Harness-Werkzeug)

- **Modus:** GF
- **Konventionen-Dichte:** Mittel — die Register-Mechanik ist Kanon, ihr
  Ausgang „geplant" führt hierher.
- **Phase-Reife:** Phase 3.
- **Evidenz-/Diskrepanz-Risiko:** Niedrig.
- **Reconciliation-Aufwand:** Keiner.
