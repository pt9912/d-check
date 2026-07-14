# ADR-0037 — RTM-Anforderungen aus Tabellen und fail-closed Nullmengen

**Status:** Proposed  
**Datum:** 2026-07-14  
**Autor:** pt9912  
**Schärft:**
[`DC-FA-REQ-001.a`](../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)  
**Bezug:**
[`DC-FA-REQ-001`](../../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen),
[`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix),
[`DC-FA-MOD-001`](../../../spec/lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in),
[ADR-0034](0034-trace-konfigurierbare-quellen.md),
[ADR-0036](0036-trace-modality-klassifikation.md),
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).

## Kontext

[ADR-0034](0034-trace-konfigurierbare-quellen.md) machte Quelldatei und
Kennungsregex der RTM konfigurierbar, ließ aber die Dokumentgrammatik fest:
Anforderungen entstehen ausschließlich aus ATX-Überschriften. Im Konsumenten
`m-trace` stehen 371 eindeutige Anforderungen in Markdown-Tabellen mit
`Kennung | Prioritaet | Anforderung` beziehungsweise
`Kennung | Prioritaet | Akzeptanzkriterium`. Trotz expliziter Quelle und passender
Regex erkennt d-check v0.42.0 null Anforderungen; selbst
`--trace --require-complete` endet mit Exit 0. Das ist bei einem als Gate
gebundenen Vollständigkeitslauf ein still grüner Pfad.

Ein Tabellenparser darf zugleich nicht zur zweiten RTM-Implementierung werden:
Referenzscan, Coverage, Waisenstatus, Modalitäts-Gating und Reporter sollen
dieselben Datenstrukturen und denselben Ablauf behalten. Der unkonfigurierte
Heading-Pfad ist ein bestehender Kompatibilitätsvertrag und muss byte-identisch
bleiben.

## Entscheidung

### 1. Format-Strategie vor dem gemeinsamen RTM-Modell

`trace.requirements.format` wählt genau einen Extraktor:

- leer oder `headings`: bestehender ATX-Heading-Extraktor;
- `table`: neuer Markdown-Pipe-Tabellen-Extraktor.

Beide liefern dieselbe interne Menge aus Kennung, Titel und
Modalitäts-Eingabetext. Erst danach laufen die bestehenden ADR-/Slice-/Coverage-
Scans, Waisenberechnung und Reporter. Es gibt weder ein neues Regelmodul noch
ein zweites Ausgabe-/Gate-Modell.

### 2. Spalten werden nach Namen, nicht nach Position gebunden

Der `table`-Block benennt `id-column`, genau eine von `text-column` oder
`text-columns` sowie optional `modality-column`. Die Liste deklariert
alternative Text-Header; jede relevante Tabelle muss genau einen davon tragen,
und jeder deklarierte Name muss in mindestens einer Tabelle vorkommen. So wird
ein Tippfehler nicht zu einer still unvollständigen Teilmenge.
Header werden nach Trimmen und Auflösen von `\|` exakt und
case-sensitiv verglichen; Positionen und Synonyme werden nicht geraten.
Mehrere relevante Tabellen in derselben Datei werden in Quellreihenfolge
zusammengeführt. Ein doppelter konfigurierter Header macht die Spaltenbindung
mehrdeutig und ist ein Fehler.

`duplicate-ids` ist `error` (Default), `first` oder `last`. Die Overrides sind
bewusste Brownfield-Politiken für historische Neuformulierungen; ohne explizite
Wahl bleibt eine Mehrfachdefinition fail-closed.

Die kleine Tabellenlexik erkennt Header plus unmittelbare Trennzeile
(`:---`, `---:`, `:---:` oder `---`, mindestens drei Bindestriche) außerhalb
von Fenced-Code. Führende/abschließende Pipes sind optional; `\|` und Pipes in
einem einzeiligen, korrekt geschlossenen Backtick-Code-Span teilen keine Zelle.
Mehrzeilige Zellen und Block-Markdown bleiben außerhalb des Vertrags. Diese
begrenzte Grammatik wird im App-Kern implementiert; eine vollständige
CommonMark-AST-Abhängigkeit wäre für drei Zellen unverhältnismäßig und würde den
Dependency-/Image-Scope vergrößern.

### 3. Modalität ist ein Extraktor-Ergebnis

Im Heading-Format bleibt der Body-Abschnitt die Modalitätsquelle. Im
Tabellenformat ist es bei gesetzter `modality-column` ausschließlich deren
Zelle, sonst die Textzelle. Der bestehende Matcher aus
[ADR-0036](0036-trace-modality-klassifikation.md) klassifiziert diesen Text
unverändert; Keywords, `unknown` und `require-levels` werden nicht dupliziert.

### 4. Explizite Absicht aktiviert fail-closed

Ein nichtleerer `requirements.source`-Wert oder `format: table` aktiviert den
strikten Quellenvertrag:

- Quelle fehlt/unlesbar oder null Definitionen erkannt → Exit 2;
- unbekanntes Format, inkonsistenter `table`-Block oder fehlende/doppelte
  konfigurierte Header → Exit 2;
- doppelte erkannte ID unter `duplicate-ids: error` → Exit 2; `first`/`last`
  lösen sie deterministisch auf.

`source: ""` bleibt gemäß bestehender Empty⇒Default-Regel abwesend. Ohne
strikten Vertrag behält der Heading-Pfad seine bisherige Semantik: fehlende
Default-Quelle/null Treffer → leere RTM Exit 0, doppelte Headings werden wie
bisher auf den ersten Treffer dedupliziert. Damit ist das Default-Verhalten
byte-identisch; fail-closed greift dort, wo Konfiguration die Nutzerabsicht
ausdrückt.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| Tabelle vorab in Headings generieren | kein Produktcode | zusätzlicher Generator + Drift-Gate je Konsument; keine native Modalitätsspalte |
| Vollständige CommonMark-AST-Bibliothek aufnehmen | breite Markdown-Abdeckung | neue Dependency und größerer Parser-/Supply-Chain-Scope für einen engen RTM-Fall |
| `--require-nonempty` als weiteres CLI-Flag | explizites Gate | `doc-complete` bleibt ohne Zusatzflag still grün; Konfiguration der Quelle drückt die Absicht bereits aus |
| Gewählt: kleiner Pipe-Tabellen-Extraktor + impliziter Strict-Modus | native Brownfield-Unterstützung, kein zweites RTM-Modell, fail-closed an der Absichtsgrenze | bewusst keine vollständige CommonMark-Tabellensemantik |

**Fitness-Funktion:**

- `m-trace`-Quelle mit den Text-Headern `Anforderung` und
  `Akzeptanzkriterium`, `text-columns` und `duplicate-ids: last` liefert aus
  372 Zeilen exakt 371 eindeutige Anforderungen; `Prioritaet=Muss`
  klassifiziert bei aktivem `modality` als
  `must`.
- Falscher Header oder unpassende Regex bei nichtleer expliziter Quelle liefert
  Exit 2 statt `total: 0`/Exit 0.
- Ohne `trace`-Block bleibt die bestehende Heading-Fixture byte-identisch;
  `source: ""` verhält sich wie abwesend.

## Konsequenzen

- Das Modell erhält Format, Tabellen-Spalten und die explizite Strict-Grenze;
  keine neue Reporterstruktur.
- Der Config-Adapter validiert Format-/Block-Konsistenz vor dem I/O; Header,
  Duplikate und Nullmenge werden laufzeitnah gegen die Quelle validiert.
- Der App-Kern erhält einen kleinen, isoliert testbaren Tabellenlexer. Die
  Modalitätsklassifikation wird auf einen vom Extraktor gelieferten Text
  angewandt.
- `--print-config` und Benutzerhandbuch zeigen die Tabellenkonfiguration.
- Read-only/Netzlos bleiben unverändert; nur die bereits konfigurierte
  Markdown-Quelldatei wird gelesen.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-07-14 | Entwurf zu slice-070 nach Konsumentenbefund `m-trace` (371 Tabellenanforderungen → 0/Exit 0). Plan-Review R1 schärfte zwei Kompatibilitätsgrenzen: `source: ""` bleibt abwesend; Duplicate-ID-Fehler nur bei Tabellenmodus oder nichtleer expliziter Quelle. Realdaten-Review R2 fand zwei Text-Header und 372 Zeilen/371 IDs; `text-columns` und explizite `duplicate-ids`-Politik ergänzt. Status Proposed. |
