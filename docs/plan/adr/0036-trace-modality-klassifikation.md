# ADR-0036 — Modalitäts-Klassifikation der Anforderungen (`trace.requirements.modality`)

**Status:** Accepted
**Datum:** 2026-07-11
**Autor:** pt9912
**Bezug:** [`DC-FA-MOD-001`](../../../spec/lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in)
(die neue Anforderung; Spezifikation
[§`DC-FA-MOD-001.a`](../../../spec/spezifikation.md#dc-fa-mod-001a--modalitäts-klassifikation-tracerequirementsmodality));
[`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
(Modality-Spalte) und
[`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
(gatende Stufen), beide mit-modifiziert;
[ADR-0035](0035-trace-coverage-quellen.md) (die Coverage-Klasse, deren
Rest-Waisen die Modalitäts-Frage auslösten);
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).

## Kontext

Nach [ADR-0035](0035-trace-coverage-quellen.md) sank die Waisen-Zahl des
Konsumenten grid-gym auf 10. Deren Analyse zeigte: die RTM behandelt jede
Anforderung gleich, obwohl das Lastenheft **RFC-2119-Modalität** trägt (gemessen:
112 MUSS, 50 SOLLTE, 9 KANN, 3 DARF NICHT je Anforderung). Die 10 „Waisen" waren
**5× GG-FUTURE (KANN)**, **4× GG-NONGOAL (deklaratives Nicht-Ziel, kein
Modal-Verb)** und **1× GG-MVP-004 (DARF NICHT)** — nur die letzte ist eine echte
Pflicht ohne Nachweis. Eine reine **KANN**-Anforderung ohne Slice/Coverage als
„Waise" wie eine unabgedeckte **MUSS**-Pflicht zu zählen ist irreführend; ein
Vollständigkeits-Gate, das darauf bricht, ist nicht benutzbar.

Die Modalität lebt im Anforderungs-**Body** (`… MUSS …`), nicht im Heading, das
die RTM heute allein liest.

## Entscheidung

### 1. Opt-in `trace.requirements.modality` mit konfigurierbaren Keywords + Defaults

Ein opt-in Unter-Block klassifiziert jede Anforderung anhand von **Modal-Verb-
Schlüsselwörtern**:

- **`levels`** — Map Stufen-Name → Schlüsselwort-Liste. **Built-in-Default** (wenn
  `modality` gesetzt, `levels` aber leer): eine DE+EN-RFC-2119-Menge —
  `must: [MUSS, MUESSEN, MÜSSEN, "DARF NICHT", "DÜRFEN NICHT", MUST, SHALL, "MUST NOT", "SHALL NOT"]`,
  `should: [SOLLTE, SOLLTEN, "SOLLTE NICHT", "SOLLTEN NICHT", SHOULD, "SHOULD NOT"]`,
  `may: [KANN, KÖNNEN, "MUSS NICHT", "MÜSSEN NICHT", MAY, OPTIONAL]`. Voll
  überschreib-/erweiterbar (sprach-/konventionsspezifisch).
- **`require-levels`** — welche Stufen `--require-complete` gaten; Default `[must]`.

Warum **konfigurierbar mit Defaults** statt fix: die Keywords sind sprach- und
konventionsabhängig (DE vs. EN, projektspezifische Verben); die Defaults machen
den Normalfall bequem (`modality: {}` genügt), die Config den Sonderfall möglich —
dasselbe Muster wie `trace.requirements.id-pattern` (Default + Override).

### 2. Klassifikation: erster Treffer, längste Phrase, Body-Abschnitt

Je Anforderung wird der **Body-Abschnitt** (Überschrift bis zur nächsten
gleich-/höherrangigen — dieselbe `rules`-Span-Mechanik wie `trace.coverage`) auf
das **erste** Schlüsselwort gescannt. An einer Position gewinnt die **längste**
Phrase (`DARF NICHT` vor `DARF`; `MUSS NICHT` [may] vor `MUSS` [must] — sonst
würde die Negation als Pflicht verschluckt). Schlüsselwörter case-insensitiv und
**wortgrenzen-genau** (`MUSS` matcht nicht `musste`). **Eine** Stufe je
Anforderung (erster Treffer) — keine NLP, keine Mehrfach-Modalität.

### 3. `unknown` sichtbar, Gating explizit (fail-open bewusst begrenzt)

Eine Anforderung **ohne** Treffer erhält die Stufe **`unknown`** — in der
Modality-Spalte **sichtbar**, nicht still. Ob `unknown` gatet, entscheidet
`require-levels`: Default `[must]` ⇒ `unknown` **advisory** (grid-gyms Nicht-Ziele
ohne Modal-Verb gaten nicht). Wer strikt will, setzt `require-levels: [must, unknown]`.
Das ist die bewusste **fail-open-Grenze**: ein echtes MUSS mit **unaufgeführtem**
Verb fiele auf `unknown` und entkäme dem Gate — Gegenmittel sind (a) umfassende
Defaults, (b) die sichtbare Spalte (ein Mensch sieht `unknown` und ergänzt das
Keyword), (c) der `unknown`-in-`require-levels`-Strikt-Modus. Ein `require-levels`-
Eintrag, der weder deklarierte Stufe noch `unknown` ist, ist ein
**Config-Fehler** (Exit 2, fail-closed) — verhindert einen Tippfehler, der still
nichts gaten würde. Ebenso fail-closed (aus dem Doc-first-Review): der
**reservierte** Stufen-Name `unknown` in `levels` (kollidiert mit dem Sentinel)
und **dasselbe Keyword in zwei Stufen** (die gewinnende Stufe hinge sonst an der
Map-Iteration ⇒ Nondeterminismus,
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).

### 3a. Body-Normalisierung vor dem Matching (Doc-first-Review MEDIUM)

Der Body-Abschnitt wird **vor** dem Keyword-Matching normalisiert: Markdown-
Emphasis (`*`/`` ` ``) entfernt und Whitespace-/Umbruch-Folgen zu einem
Leerzeichen zusammengezogen. Sonst matchte eine hart umbrochene oder emphasierte
Negation (`MUSS\nNICHT`, `**MUSS** NICHT`) die Phrase `MUSS NICHT` **nicht** und
fiele still auf `MUSS` (must) zurück — genau die Negation, die der
Längster-Treffer-Mechanismus schützt. Wortgrenze via RE2-`\b` (**ASCII**; die
Default-Keywords haben ASCII-Ränder und sind sicher, ein konfiguriertes
Umlaut-Rand-Keyword nicht — im Vertrag caveatet).

### 4. Konditionale Spalte, byte-identisch im Default

Modality-Spalte + `modality`-json/yaml-Feld erscheinen **nur** bei aktivem
`modality`; **aktiv = Schlüssel-Präsenz** (`modality: {}` ⇒ Defaults + Spalte),
**nicht** `len(levels)>0` (sonst wäre `{}` fälschlich inaktiv — Doc-first-Review
LOW). Ohne ist die RTM **byte-identisch** und `--require-complete` gatet **alle**
Waisen wie bisher
([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)) — dieselbe
Konditional-Mechanik wie die Coverage-Spalte (ADR-0035).

### Verglichene Alternativen

| Option | Pro | Contra |
| --- | --- | --- |
| **(A) Modalität ignorieren (Status quo)** | null Code | KANN/MUSS gleichbehandelt; `--require-complete` für modalitäts-tragende Repos unbenutzbar; grid-gyms 10 „Waisen" bleiben Rauschen |
| **(B) Feste Keyword-Menge (nicht konfigurierbar)** | einfacher | sprach-/konventions-gebunden; ein Repo mit anderen Verben (oder EN-only) mis-klassifiziert still |
| **(C) Familien-/Sektions-Ausschluss (statt Modalität)** | löst grid-gyms NONGOAL/FUTURE | unprinzipiell (keine MUSS/SOLLTE-Unterscheidung **innerhalb** der echten Anforderungen); enumeriert Familien (RE2 ohne Negative-Lookahead) |
| **(D, gewählt) opt-in `modality`, konfigurierbare Keywords + Defaults, `require-levels`, `unknown` sichtbar** | RFC-2119-prinzipiell; DE+EN-Default bequem, override möglich; byte-identisch im Default; fail-open sichtbar/konfigurierbar | Body-Heuristik (erster Treffer); Negations-Kanten via Längster-Treffer |

**Fitness-Funktion:**

- grid-gym mit `modality: {}` (Defaults) + `require-levels: [must]`: die 5 KANN +
  4 Nicht-Ziele (`unknown`) werden **advisory**, nur GG-MVP-004 (`DARF NICHT`)
  bleibt gatende Waise (Beleg im umsetzenden Slice, an Realdaten gemessen).
- `… MUSS NICHT …` ⇒ `may` (nicht `must`); `… DARF NICHT …` ⇒ `must` (Negations-
  Beweis, Längster-Treffer).
- Ungültiges `require-levels` / leeres Keyword ⇒ Exit 2.
- **Ohne** `modality`: RTM byte-identisch, `--require-complete` gatet alle Waisen.

## Konsequenzen

- **Neue Anforderung + zwei Modifikationen (Lastenheft-CR, Versions-Bump +
  Historie).** [`DC-FA-MOD-001`](../../../spec/lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in)
  neu (Bereich `MOD`);
  [`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
  (Modality-Spalte) und
  [`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
  (gatende Stufen) mit-modifiziert. Nutzersichtbare Config + geänderter Output/
  Gate-Verhalten bei aktivem `modality` → **Release**.
- **Modell** (`config.go`): `TraceModality` (Levels map + RequireLevels) +
  `TraceConfig.Modality`; `TraceRow.Modality string` (`omitempty`);
  `TraceMatrix.ModalityActive`.
- **Config-Decode** (`configyaml.go`): `modality`-Block; Default-Keyword-Menge bei
  leerem `levels`; `require-levels` gegen Stufen-Namen + `unknown` validiert
  (fail-closed).
- **Klassifikator** (`trace.go`): Body-Abschnitt je Anforderung (Span-
  Wiederverwendung), erster/längster Keyword-Treffer → Stufe; `unknown`-Fallback.
- **Waise-Gating** (`cli.go`/`trace.go`): `--require-complete` filtert Waisen auf
  `require-levels`. **Reporter** (`report.go`): konditionale Modality-Spalte.
- **`--print-config`**: kommentierter `modality`-Block. **Determinismus/Read-only**
  unberührt (nur Body-Lesen im gemounteten Baum). **Reversibel** im Verhalten
  (Default byte-identisch), aber Vertrags-Änderung → Lastenheft-CR.

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-07-11 | Entwurf (slice-068, welle-57; Nutzer-Frage „unterscheidet Traceability MUSS/SOLLTE?" + Hinweis, dass das Lastenheft MUSS/SOLLTE/DARF NICHT trägt — meine erste Analyse [„Modalität hülfe nicht"] war durch einen Body-Scan-Bug falsch; die 10 Coverage-Rest-Waisen sind 5 KANN + 4 Nicht-Ziele + 1 DARF NICHT). Opt-in `trace.requirements.modality`: konfigurierbare Modal-Verb-Keywords (DE+EN-RFC-2119-Defaults, `levels`/`require-levels`), erster/längster Treffer im Body-Abschnitt, `unknown`-Fallback sichtbar; Modality-Spalte konditional; `--require-complete` gatet nur `require-levels`-Stufen (Default `[must]`). Byte-identisch ohne `modality`; fail-closed (Config/`require-levels`). Design spiegelt `trace`-Config (Default+Override). Lastenheft-CR (v0.42.0), Release geplant. Status Proposed. |
| 2026-07-11 | **Accepted** (slice-068 umgesetzt + zwei Reviews eingearbeitet). Präzisierungen aus R1: Body-**Markup-Normalisierung** vor dem Match (`**MUSS** NICHT`/`MUSS\nNICHT` ⇒ `may`, nicht fälschlich `must`); dasselbe Keyword in zwei Stufen ist **fail-closed** (Exit 2, Nondeterminismus); Aktivierung = **Schlüssel-Präsenz**. End-to-End gegen grid-gym: 10 Coverage-Rest-Waisen ⇒ unter `modality: {}` + `require-levels: [must]` gaten nur **2** (GG-MVP-004 `DARF NICHT` + GG-NONGOAL-005 `muessen`-Klausel), acht advisory. Release v0.42.0. |
