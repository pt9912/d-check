# Slice slice-236: Modul `file` — Obergrenze für Zeilen und Bytes einer ganzen Datei

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** — **wellenlos**. Die Closure-Bedingung geht über die eigene DoD
nicht hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine
Welle braucht).

**Bezug:** neue Anforderung — ihre Kennung wird beim Schreiben des Lastenhefts
vergeben, nicht im Plan. Angrenzend:
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(die abschnittsbezogene Nachbarin),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Determinismus). Auftrag des Auftraggebers am 2026-09-27: „max-rows und
max-bytes einer Datei".

**Berührte Spec-Stellen:** — (die neue Fähigkeit trägt beim Schreiben der
Spezifikation ihre eigene Kennung).

**Verantwortlich:** claude-sonnet-5.

**Autor:** claude-sonnet-5. **Datum:** 2026-09-27.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** d-check meldet eine Datei, die mehr Zeilen oder mehr Bytes hat als
konfiguriert. Gezählt wird die **rohe** Datei, für **jede** Dateiart, ohne
Abschnitts-Selektor.

**Befund, gemessen (2026-09-27).** Einen Schlüssel dafür gibt es nicht. Der
Umweg über `structure[].forbid-pattern` (RE2 gegen den Abschnittstext) trägt
nicht: (a) gezählt wird der **bereinigte** Text — an `docs/user/releasing.md`
(258 Zeilen roh, drei Fenced-Blöcke) feuert eine Obergrenze bei 230 und nicht bei
240; (b) die Regel braucht einen Abschnitt, also eine H1 als Anker; (c) ein
Wiederholungszähler ist auf 1000 begrenzt, `(?:.{1000}){40}` ⇒ Exit 2, größere
Grenzen nur als lange Verkettung; (d) gezählt werden Zeichen, nicht Bytes; (e)
`structure` prüft nur Markdown. Anlass ist die offene Beobachtung
[`briefing-datei-ueberschreitet-lade-budget`](../observations/BEO-ALL/briefing-datei-ueberschreitet-lade-budget/observation.md):
sie hat heute keinen Sensor.

**Entscheidungen vor dem Code — beim Auftraggeber, in der ADR festgehalten:**

- **Ort — entschieden (Auftraggeber, 2026-09-27): ein eigenes Modul `file`** mit
  eigener Anforderung, opt-in. Grund: `structure` ist ausdrücklich abschnitts-
  und Markdown-bezogen und verlangt einen Abschnitts-Selektor; die Grenze einer
  ganzen Datei gilt jeder Dateiart und braucht keinen. Das Kriterium der
  Historie: querschnittlich ⇒ eigenes Kürzel, Einzelmodul-Frage ⇒ bestehende
  Anforderung ändern. Die Erweiterung von `structure` bleibt in der ADR als
  verworfene Alternative mit ihrem Grund stehen.
- **Schlüsselnamen:** der Auftrag nennt `max-rows`; d-check spricht überall von
  **Zeilen** (`line`-Feld, „Zeilen") — Vorschlag `max-lines`, bestätigt der
  Auftraggeber `max-rows`, gilt das.
- **Zählregel:** Zeilen = Zeilenumbrüche, plus eine unvollständige Schlusszeile;
  Bytes = Dateigröße. Beides an Grenzwert **N und N+1** zu belegen.
- **Grund-Codes:** je Grenze einer, weil jede eine andere Reparatur verlangt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Schwellen für `AGENTS.md` und andere Dateien festlegen** — eine Schwelle ist
  eine Gate-Einführung mit Zahl und braucht ihre eigene ADR (`AGENTS.md` §3.6);
  sie hängt außerdem an der offenen Entscheidung, wie groß die Datei sein darf.
  Übernimmt ein Folge-Slice, der noch keine Kennung hat und erst nach
  Closure dieses Slice geschnitten wird.
- **Zeichen-, Wort- oder Token-Zählung** — nicht beauftragt; Bytes und Zeilen
  sind die Zusage.
- **Zählen des bereinigten Textes** (ohne Fences, ohne Inline-Code) — das ist
  die Eigenschaft von `structure`, die diese Fähigkeit gerade **nicht** haben soll.
- **Reparatur (`--repair`)** — eine zu große Datei ist keine mechanische
  Reparatur.
- **Handbuch, README, CHANGELOG, Release** — Release-Prep, kein Feature-Commit
  (`AGENTS.md` §5).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — höchstens drei Liefer-Punkte.

- [x] **Vertrag:** Lastenheft (neue Anforderung für das Modul `file`, mit Happy · Boundary · Negative
      und Out-of-Scope, Versions-Bump mit Historie-Zeile nach
      [MR-032](../../../../harness/conventions/MR-032-historie-vor-accepted.md)),
      Spezifikation (Zählregel, Schema-Schlüssel, Grund-Codes, Grenze nach
      `AGENTS.md` §3.8: die Datei-Menge kommt aus dem Glob, nicht aus
      `scan.roots`) und eine neue ADR samt Index-Eintrag. **Form nach dem Kanon:**
      die Historie-Zeile nennt weder ADR noch Slice, die Spezifikation in keinem
      Abschnitt; die ADR trägt `Schärft:` aufwärts, mindestens drei verglichene
      Alternativen (unter anderem: `structure` erweitern · eigenes Modul ·
      Muster-Umweg), eine Fitness Function und einen `Re-Evaluierungs-Trigger`.
- [x] **Umsetzung:** die Fähigkeit samt Konfigurations-Parser. Tests: Grenzwert
      **N und N+1** je Schlüssel (Rot-Beleg gegen den Stand davor), Datei ohne
      Schlusszeilenumbruch, leere Datei, Umlaute (Bytes ≠ Zeichen), Datei mit
      Fenced-Blöcken (roh gezählt — die Abweichung aus §1 tritt **nicht** auf),
      Nicht-Markdown-Datei, Glob ohne Treffer (Befund statt Stille), ohne
      Konfigurationsblock Befundsatz byte-identisch
      ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- [x] **Spiegel:** jede Aufzählung der Module und Schlüssel im Repo ist gelesen
      und, wo sie das Modul führen muss, im Feature-Diff mitgezogen (Liste in
      §3 **vor** dem Editieren); was Release-Prep ist, steht dort benannt.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8, kein
      Self-Review.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — oder „keine
      Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      wellenlos hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | neue Anforderung; Bump + Historie (eigener Commit vor dem Code) |
| `spec/spezifikation.md` | update | Zählregel, Schema, Grund-Codes, Grenze |
| `docs/plan/adr/` (neu) + `README.md` | neu / update | Ort, Namen, Zählregel; Alternativen |
| Regel-Paket, Konfigurations-Parser, Modul-Registrierung unter `internal/` | neu / update | Umsetzung; genaue Dateien nach dem Lesen des Bestands |
| Akzeptanz- und Modul-Tests | neu | Grenzwert-Belege aus §2 |

**Reihenfolge der Vertrags-Änderung.** Das Lastenheft steht auf `Draft`: der
Kanon (`grundlagen-source-precedence.md`, *Wann die CR-Pflicht beginnt*) lässt
es vor `Accepted` frei änderbar; [MR-032](../../../../harness/conventions/MR-032-historie-vor-accepted.md)
verlangt Bump und Historie-Zeile trotzdem. Die Änderung steht in einem eigenen
Commit vor dem Code.

**Vor dem Editieren — Spiegel listen**
([MR-025](../../../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)):
die Modul-Aufzählungen (Lastenheft-Modulliste, Konfigurations-Gerüst
`--print-config`, `--suggest-config`, `--print-mk`, Handbuch, Betriebsdoku,
README beider Sprachen) und die Schlüsseltabelle der Spezifikation. Was davon
Release-Prep ist, wird **benannt**, nicht übergangen.

## 4. Trigger

**Start** (`open` → `in-progress`): Auftraggeber-Priorisierung 2026-09-27 —
dieser Slice startet **vor** `slice-232`/`-233`/`-234`, die bislang vor ihm
in der Warteschlange standen. Alle vier bumpen das Lastenheft; da sie
nacheinander (WIP-Limit 1) implementiert werden, kollidiert kein
Versions-Bump — es ändert sich nur die Reihenfolge. Bei der Beanspruchung
entsteht der dritte Vorprüfungs-Block (Nachtlauf-Stand).

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): wenn die Spiegel-Liste zeigt, dass die
  Modul-Registrierung mehr als eine Review-Sitzung an Nachzug verlangt.
- `in-progress` → `open` (blockiert): solange die Ort- oder Namens-Entscheidung
  des Auftraggebers aussteht.

## 5. Closure-Trigger

DoD vollständig (§2) + Review-Report liegt vor + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Ein neues Modul fehlt in einer Aufzählung, die kein Gate hält — der Bestand
  zeigt, dass die Modul-Listen von Hand nachgezogen werden.
  **Ausgang:** eingetreten — der unabhängige Review fand drei bisher nicht
  benannte Fundorte (Bereichskürzel-Liste, Glossar, `operations.md`;
  F-2/F-3/F-4), behoben im Fix-Commit. Die verwandte Beobachtung
  [`BEO-ALL/modulliste-spiegel-ungegated`](../observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md)
  erreicht damit 3× — Ausgang `geplant`,
  [slice-238](slice-238-modul-registrierungs-spiegel-checkliste.md).
- Die Datei-Menge kommt aus dem Glob und nicht aus dem Scan-Bereich; eine
  Datei außerhalb der Scan-Wurzeln wird geprüft, eine gelöschte nicht.
  **Ausgang:** entfallen — das ist die bewusst gewählte, in der Spezifikation
  benannte Eigenschaft (wie bei `structure`/`mentions`), kein unentdeckter
  Mangel.
- Ohne die Schwellen-Entscheidung bleibt die Fähigkeit ungenutzt, und
  die Beobachtung zur Größe der Briefing-Datei hat weiter keinen Sensor.
  **Ausgang:** weiter offen —
  [`BEO-ALL/briefing-datei-ueberschreitet-lade-budget`](../observations/BEO-ALL/briefing-datei-ueberschreitet-lade-budget/observation.md)
  (jetzt 2×): der Sensor existiert, seine Nutzung für `AGENTS.md` ist ein
  eigener Folge-Slice ohne Kennung.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Formvorlage bestehender Module (`structure`,
  `reviews`, `mentions`) hat die Umsetzung geführt — Zähler-Wiederverwendung
  (`countLines`, geteilt mit `codepaths`/`citations`) statt einer eigenen
  Implementierung, `Identity()`/`MessageFor()` nach demselben Muster wie
  `StructureRule`. Der unabhängige Review hat die Kernimplementierung
  (Zählung, Config-Validierung, Hexagon-Schichtung, Determinismus,
  Referenzrichtung, ADR-Form) ohne Beanstandung bestätigt.
- **Was ging anders als geplant:** (1) Die geplante Ort-Entscheidung
  (eigenes Modul statt `structure`-Erweiterung) stand vor dem Code fest
  (Auftraggeber, 2026-09-27) und hat sich im Nachhinein bestätigt — kein
  Rückbau nötig. (2) Zwei eigene Fehler traten während der Umsetzung auf,
  vor dem Review korrigiert: eine `cut -c1-200`-Terminaltrunkierung hat beim
  ersten Versuch, „file" in die `--print-config`-Verfügbar-Zeile
  einzufügen, das Zeilenende (`external, sources`) verschluckt; die
  Spezifikation sagte zunächst „unlesbare Einzeldatei liefert keinen Befund"
  zu, inkonsistent mit `structure`s eigenem fail-closed-Vorbild für denselben
  Fall — beim Vergleich mit der Nullmengen-Härte aufgefallen und auf
  fail-closed vereinheitlicht. (3) Trotz eigener Sorgfalt fand der
  unabhängige Review sieben weitere Befunde (§DoD, unten) — am schwersten
  wog F-1: der Umsetzungs-Commit hat README-Bullets ergänzt, obwohl §1
  dieses Plans „Handbuch, README, CHANGELOG, Release — Release-Prep, kein
  Feature-Commit" ausdrücklich ausschließt. Ich habe die eigene Abgrenzung
  im Lauf verletzt, nicht der Kanon war unklar.
- **Retirement-Check:** keine bestehende Regel wurde gelockert oder entfernt;
  entfällt.
- **Steering-Loop-Eintrag:** gezählt, nicht verkörpert für den Fall F-1
  (Plan-Abgrenzung im eigenen Lauf verletzt, neuer Eintrag
  [`BEO-ALL/plan-abgrenzung-im-selben-lauf-verletzt`](../observations/BEO-ALL/plan-abgrenzung-im-selben-lauf-verletzt/observation.md))
  — Beobachtung steht bei 1×, unter der Schwelle.
- **Beobachtungs-Register (`../observations/`):**
  `evidence/slice-236.md` in
  [`BEO-ALL/modulliste-spiegel-ungegated`](../observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md)
  ergänzt — Zähler steht damit bei 3×, Ausgang `geplant`
  ([slice-238](slice-238-modul-registrierungs-spiegel-checkliste.md));
  `evidence/slice-236.md` in
  [`BEO-ALL/briefing-datei-ueberschreitet-lade-budget`](../observations/BEO-ALL/briefing-datei-ueberschreitet-lade-budget/observation.md)
  ergänzt — Zähler steht bei 2×, weiter offen;
  [`BEO-ALL/plan-abgrenzung-im-selben-lauf-verletzt`](../observations/BEO-ALL/plan-abgrenzung-im-selben-lauf-verletzt/observation.md)
  neu angelegt, Beleg `evidence/slice-236.md` (1×).
- **Folge-Slices:** [slice-238](slice-238-modul-registrierungs-spiegel-checkliste.md)
  (Sensor oder Checkliste für die Modul-Registrierungs-Spiegel) — ist eine
  Datei in `open/`.
- **Risiken aus §6:** eines eingetreten (Folge-Slice slice-238), eines
  entfallen (benannte Grenze, kein Mangel), eines weiter offen (Register).
- **Drei Paarungen:** Anker — kein `liegt in`-Feld verwendet, nichts
  verkörpert, kein Gegenstand. Folge-Slice — `slice-238` existiert als Datei
  in `docs/plan/planning/open/`, geprüft ok. Register — beide genannten
  Beobachtungs-Verzeichnisse existieren und tragen je einen neuen Beleg für
  diesen Vorgang, geprüft ok.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
**Treffer:**
[`briefing-datei-ueberschreitet-lade-budget`](../observations/BEO-ALL/briefing-datei-ueberschreitet-lade-budget/observation.md)
(1×) — der Anlass; die Fähigkeit ist der fehlende Sensor, ihre Schwelle ist es
nicht (§1).

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung, 2026-09-27):
`make nightly-state` meldet `upstream-drift.yml` **ROT** (Lauf
2026-09-27T06:03Z), `image-scan.yml` **grün**. Lokal nachgefahren: `golangci-lint`
VERALTET (Pin 2.13.2, upstream 2.14.0), `semgrep` VERALTET (Pin 1.177.0,
upstream 1.178.0), `a-check` VERALTET (Pin 0.19.0, upstream 0.20.0),
`golang:1.27.1` ABWEICHEND (Digest unter demselben Tag neu gebaut). Vier
planmäßige Fremd-Release-Meldungen, keine unerwarteten; sie berühren dieses
Modul nicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

### Sub-Area: `*` (Produkt: Spezifikation und Code)

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — Lastenheft, Spezifikation und ADR sind Kanon.
- **Phase-Reife:** Phase 5.
- **Evidenz-/Diskrepanz-Risiko:** Niedrig bis mittel — additiv und opt-in; das
  Risiko liegt in den von Hand gepflegten Modul-Listen (§6).
- **Reconciliation-Aufwand:** Keiner — kein Brownfield-Bestand.
