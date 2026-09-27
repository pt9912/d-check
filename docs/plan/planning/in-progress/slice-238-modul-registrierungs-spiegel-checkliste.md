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

**Verantwortlich:** pt9912.

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

- [x] **Entscheidung (A) oder (B)** getroffen und begründet — bei (A) reicht
      eine ADR-lose Fitness-Function-Ergänzung (kein neues Modul, kein neues
      `DC-*`), bei (B) eine `AGENTS.md`-Ergänzung mit Herkunfts-Anker.
- [x] **Vollständige Fundort-Liste** einmal zusammengetragen (Diese Datei,
      `harness/README.md`, `AGENTS.md`, `spec/lastenheft.md` (§3, §6),
      `spec/spezifikation.md` §2/§4, `docs/user/*.md`, `README*.md`,
      `internal/adapter/driving/cli/config_template.go`,
      `internal/hexagon/core/app/diagnose.go`, `internal/hexagon/core/model/config.go`,
      `.d-check.yml`, `Makefile` `FOCUS_DISABLE`) — geprüft durch tatsächliches
      Lesen, nicht durch Wiederholung dieser Aufzählung.
- [x] Umsetzung nach der getroffenen Entscheidung.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — der
      Ausgang von `BEO-ALL/modulliste-spiegel-ungegated` (heute `geplant:
      slice-238`) wird auf `verkörpert` gesetzt.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen sind getragen — wellenlos hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/core/app/registry_mirror_test.go` | neu | sechs Deckungstests gegen `model.ValidModules()` (Lastenheft-Beschreibung, Glossar, Operations-Optionen-Tabelle, README.md/README.de.md-Bullets, Handbuch-Tabelle) plus Guard-Test |
| `internal/adapter/driven/configyaml/gate_consistency_test.go` | update | dritte Prüfrichtung in `assertNetlessModules` (unbekanntes, weder gelistetes noch verbotenes Modul); neuer Deckungstest `FOCUS_DISABLE` gegen `.d-check.yml` `modules:` |

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
  aufwändig werden. **Ausgang: entfallen.** Zwei Klassen bewusster
  Teilmengen wurden identifiziert und explizit AUS der mechanischen Prüfung
  genommen, statt sie mit einer Ausnahme-Liste nachzubilden: die
  Bereichskürzel-Liste (Abkürzungen, keine wörtlichen Modulnamen) und das
  „fixe Standard-Modulset" der `ai-harness`-Gerüste (bewusste Teilmenge,
  kein Vollständigkeits-Anspruch) — beide von zwei unabhängigen Reviews
  bestätigt. Der zusätzliche Aufwand blieb dadurch klein: sechs
  Deckungstests plus zwei Guard-Tests, keine Ausnahme-Liste nötig.

## 7. Closure-Notiz

**Geliefert:** Option (A) — sechs mechanische Deckungstests
(`internal/hexagon/core/app/registry_mirror_test.go`) halten
`model.ValidModules()` gegen sechs wörtlich-vollständige Doku-Fundorte
(Lastenheft-Beschreibung
[`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl),
Lastenheft-Glossar,
`docs/user/operations.md` Optionen-Tabelle, `README.md`, `README.de.md`,
`docs/user/benutzerhandbuch.md` §6 Regelmodule-Tabelle), plus eine dritte
Prüfrichtung im bestehenden Netzlos-Guard
(`internal/adapter/driven/configyaml/gate_consistency_test.go`,
`assertNetlessModules`) und ein neuer Deckungstest, der `FOCUS_DISABLE`
(Makefile) gegen `.d-check.yml` `modules:` hält — ein vierter, andersartiger
Fundort-Typ (Spiegel der aktiven Konfiguration statt der vollen Registry).
Jeder Mechanismus trägt einen Guard-Test mit synthetischen Eingaben
(fehlendes/verwaistes/unbekanntes Modul), kein neues `DC-*`, keine ADR
(ADR-lose Fitness-Function-Ergänzung, wie im Plan vorgesehen).

**Bewusst nicht mechanisiert, mit Begründung:** die Bereichskürzel-Liste
(`spec/lastenheft.md` §3) ist eine Abkürzungs-Liste, keine wörtliche
Modulnamen-Liste; die „fixe Standard-Modulset"-Stellen (`--suggest-config
ai-harness`-Gerüst, drei Fundorte) führen eine bewusste **Teilmenge** ohne
Vollständigkeits-Anspruch. Beide Ausschlüsse wurden von zwei unabhängigen
Reviews eigenständig nachgeprüft und bestätigt.

**Zwei Review-Runden.** R1 fand F-1 (HIGH, merge-blockierend): drei neue
Kommentare trugen verbotene Review-Befund-Marker (`AGENTS.md` §3.7,
teils wortgleich aus dem Bestand übernommen), sowie F-2 (MEDIUM): zwei
prominente, wörtlich-vollständige Modul-Spiegel (`README.md`/`README.de.md`,
Handbuch-Tabelle) blieben zunächst ungedeckt — genau die Art Lücke, die
diesen Slice ausgelöst hat, hätte sich sonst am wahrscheinlichsten
wiederholt. Beide behoben; R2 verifizierte die Korrektur mit einer
vollständigen Neu-Prüfung, eigenen adversariellen Regex-Proben (24/24
Treffer ohne Fehltreffer) und zwei selbst gebauten Mutationsproben und fand
einen weiteren, nicht blockierenden LOW-Fund (ein veraltendes
Ordinal-Element „3. Evidenz" in zwei Kommentaren) — sofort behoben.

**Register-Ausgang:**
[`BEO-ALL/modulliste-spiegel-ungegated`](../observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md)
wechselt von `geplant: slice-238` auf `verkörpert` — liegt in
`internal/hexagon/core/app/registry_mirror_test.go` (seit slice-238).

**Risiko-Ausgang:** das einzige Risiko aus §6 *entfallen* — siehe dort.

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

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung, 2026-09-27):
`make nightly-state` meldet `upstream-drift.yml` **ROT** (Lauf
2026-09-27T06:03:41Z, unverändert seit slice-232/233/235), `image-scan.yml`
**grün** (Lauf 2026-09-27T09:18:12Z). Derselbe veraltete Stand: vier der
fünf gemeldeten Fremd-Release-Stände sind bereits gehoben, die Kurs-Baseline
(`v6.9.0` → `v6.10.0`) bleibt bewusst zurückgestellt bis `v6.11.0`. Keiner
der fünf Punkte berührt dieses Harness-Werkzeug.

**Modus-Begründungsblock:** GF.

### Sub-Area: `*` (Harness-Werkzeug)

- **Modus:** GF
- **Konventionen-Dichte:** Mittel — die Register-Mechanik ist Kanon, ihr
  Ausgang „geplant" führt hierher.
- **Phase-Reife:** Phase 3.
- **Evidenz-/Diskrepanz-Risiko:** Niedrig.
- **Reconciliation-Aufwand:** Keiner.
