# slice-259: Festlegungen der Harness-Werkzeuge in der Spezifikation — Abschnitt und die Gates von `make gates`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [`MR-074`](../../../../harness/conventions.md#mr-074) (Bewegung 2:
Festlegungen der Harness-Werkzeuge gehören in die Spezifikation), Baseline
`v6.17.0` · `regelwerk/grundlagen-referenz-richtung.md` §Spec-Straten,
`regelwerk/modul-03-spec.md` §Ziel-Form: Spezifikation, Vorlagen
`templates/spec/spezifikation.template.md` (§7 neu, Historie §8) und
`templates/harness/sensors/gate.template.md`;
[`MR-0098`](../../../../harness/conventions.md#mr-0098) (Kopplung an die
Historie-Ausnahme); Auftraggeber-Freigabe 2026-10-07.

**Berührte Spec-Stellen:** `spec/spezifikation.md` §7 (neu) und §8 (Historie,
bisher §7).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-08.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die Spezifikation bekommt den Abschnitt **§7 Festlegungen der
Harness-Werkzeuge** — was ein Gate prüft und wie es an seinen Randformen
entscheidet —, die Historie rückt nach **§8**. Erste Einträge sind die vier
Gates aus `make gates`, deren Festlegung bisher nur in ihrer Sensor-Datei
steht: `coverage-gate` (Schwelle, Messbasis), `lint` (Profil, Ausnahmen),
`semgrep` (Regelset, Umfang, Befund-Politik) und `baseline-verify` (die drei
Fragen). Ihre Sensor-Dateien verlinken danach die `SPEC`-Kennung, statt
Schwelle und Randform selbst zu führen; sie sagen, wie ein Lauf zu **lesen**
ist.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Werkzeuge außerhalb von `make gates`** (`image-scan`, `image-test`-Teil
  ohne DC-Verfeinerung, `blackbox-probe`, Wächter und Hooks, Frische-Achsen,
  `nightly-state`, `record-gates`, `archive-wave`) — Folge-Slice slice-260;
  zusammen wären es über ein Dutzend Einträge, mehr als ein Review fasst.
- **Gates mit eigener Anforderung** (`doc-check`, `gate-consistency`,
  `planning-check`, `workflow-pins`, `trace-check`, …) — ihre Festlegung steht
  schon als Verfeinerung der Anforderung in §1; die Vorlage verlangt für sie
  keinen §7-Eintrag.
- **`Schärft:`-Feld der Gate-ADRs** ([ADR-0006](../../adr/0006-lint-profil-solid.md), [ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md), …) — sie sind
  `Accepted` und unveränderlich; der Bestand bleibt bewusst stehen, neue
  Gate-ADRs zeigen auf §7.
- **Neue Anforderungen im Lastenheft** — eine Festlegung eines Werkzeugs legt
  fest, wie geprüft wird, nicht, was das Produkt verspricht.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] `spec/spezifikation.md`: neuer §7 nach der Vorlage, Historie als §8;
      die Kopplungen nachgezogen — `matrix.exclude-sections` und die
      `structure`-Regel der Historie in der `.d-check.yml`, dazu ein
      Konventions-Nachtrag zu [`MR-0098`](../../../../harness/conventions.md#mr-0098) (die Ausnahme gilt seit der
      Umnummerierung für zwei verschiedene Überschriften).
- [ ] §7-Einträge (`SPEC-<NNN>`) für `coverage-gate`, `lint`, `semgrep`,
      `baseline-verify`, je mit dem, was als Treffer gilt, und wie die
      Randformen entschieden sind — am Code und an der Konfiguration geprüft,
      nicht aus der Sensor-Datei abgeschrieben.
- [ ] Die Sensor-Dateien der vier (bzw. die Index-Zeile von `coverage-gate`,
      das keine Sensor-Datei hat) verlinken die Kennung; Schwelle und
      Randform stehen nur noch in der Spezifikation. `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft;
      [MR-074](../../../../harness/conventions.md#mr-074) Bewegung 2 als
      teilweise eingelöst vermerkt (Rest: slice-260).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` | update | §7 neu, Historie §8, vier Einträge, Historie-Zeile |
| `.d-check.yml` | update | `exclude-sections`, `structure`-Abschnitt der Historie |
| Konventions-Nachtrag zu [MR-0098](../../../../harness/conventions.md#mr-0098) | neu | die Ausnahme trägt jetzt `7. Historie` (Lastenheft) und `8. Historie` (Spezifikation) |
| `harness/sensors/{lint,semgrep,baseline-verify}.md`, `harness/README.md` | update | Verweis auf die Kennung statt eigener Festlegung |

*(Plan-Änderung nach R1, vor dem Code: `tools/coverage-gate.sh` wird an die Festlegung angepasst statt umgekehrt — eine leere, nicht numerische oder negative Schwelle bestand bisher still grün, und der Exit-2-Zweig für einen unlesbaren Wert war tot (R1 F-1); die Kommentare in Skript, `Dockerfile` und `Makefile` zeigen auf [`SPEC-089`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge) (F-4). Die Schlagwort-Liste von `--suggest-config` (Produkt-Code) übernimmt ein eigener Folge-Slice (F-6). Mitgenommen: `LC_ALL=C` im Skript — unter einer Locale mit Dezimalkomma endete selbst der bestandene Fall mit 1 (gemessen auf dem Host; im Container C-Locale).)*

*(Plan-Änderung nach R2, vor dem Code: der Skript-Kommentar zur Schwellen-Prüfung trägt die Zusage statt des früheren Verhaltens, die Muster-Zeile im Kopf entfällt (R2 F-1); [`SPEC-089`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge) und [`SPEC-091`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge) nennen ihre Exit-Codes als die des Skripts — über `make` endet jedes Scheitern mit 2, die Unterscheidung trägt die Meldung (F-2); [`SPEC-091`](../../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge) und Grenze 2 von `harness/sensors/semgrep.md` nennen die Voreinstellung von semgrep als Ausschluss-Klasse samt der gemessenen Verzeichnis-Ausnahmen, nicht nur `*_test.go` (F-3); die Hilfe-Zeile von `make coverage-gate` nennt die Schwelle nicht mehr als Zahl (F-5).)*

**Spiegel vor dem Editieren** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
die Überschrift `7. Historie` der Spezifikation steht in `.d-check.yml`
(zweimal), in [`MR-0098`](../../../../harness/conventions.md#mr-0098) (Titel, Geltungsbereich, Adaption), in der Index-Zeile
von [`MR-0098`](../../../../harness/conventions.md#mr-0098) in `harness/conventions.md` und in den Kommentaren der
`.d-check.yml`; ein Anker-Link auf `#7-historie` der Spezifikation existiert
nicht (gemessen). Die Schwelle `93 %` steht in `harness/README.md` und
`tools/coverage-gate.sh`/`Dockerfile` (`COVERAGE_THRESHOLD`) — die
Spezifikation wird die Zusage, der Code bleibt der Träger.

## 4. Trigger

**Start** (`next` → `in-progress`): Implementer übernimmt; `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): die Umnummerierung reißt mehr Kopplungen
  auf als die gemessenen (etwa ein Werkzeug, das die Überschrift wörtlich
  erwartet) — dann Umnummerierung und Einträge trennen.
- `in-progress` → `open` (blockiert): keiner bekannt.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- **Zwei Orte für eine Schwelle** — die Coverage-Schwelle steht als Zusage in
  der Spezifikation und als Wert im Build (`COVERAGE_THRESHOLD`); kein Gate
  hält beide gleich. — **Ausgang:** *(offen)*
- **Die Matrix sieht die Umnummerierung nicht** — ohne Nachzug fiele die
  Historie der Spezifikation still aus der Ausnahme und meldete ihre frozen
  Verweise. — **Ausgang:** *(offen)*

## 7. Closure-Notiz

*(gefüllt vor dem `git mv` nach `done/`)*

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** eine berührte Sub-Area: die
Spezifikation und die Harness-Doku unter dem Default `*` (`ALL`); deklariert.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:**
[`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/state.md)
— die Festlegungen werden am Code geprüft, nicht aus den Sensor-Dateien
übernommen;
[`BEO-ALL/semantic-change-body-only-edges-stale`](../observations/BEO-ALL/semantic-change-body-only-edges-stale/state.md)
— die Spiegel der Umnummerierung stehen in §3.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-08 aus dem **jüngsten** Lauf — `upstream-drift` rot
(semgrep 1.180.0, a-check v0.23.0 upstream; beide Pins vor diesem Slice
gehoben, Lauf neu ausgelöst), `image-scan` grün.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
