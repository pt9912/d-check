# ADR-0108: `trace.requirements` ist ein Objekt oder eine Liste von Quellen

**Status:** Accepted

**Datum:** 2026-10-20

**Autor:** pt9912

**Bezug:** [`DC-FA-REQ-001`](../../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen)
(die Anforderung, die die Liste trägt);
[`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
(die RTM, die die Vereinigung führt);
[`DC-FA-MOD-001`](../../../spec/lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in)
(`modality`, die je Eintrag gilt);
[ADR-0037](0037-trace-tabellenquellen-nullmengen-guard.md) (der
Quellenvertrag, der hier je Eintrag gilt);
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus).

**Schärft:**
[`DC-FA-REQ-001.a`](../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0037](0037-trace-tabellenquellen-nullmengen-guard.md) lässt
`trace.requirements` genau eine Quelle mit genau einem `format` annehmen. Ein
Adopter führt sein vertraglich bindendes Lastenheft in **einer** Datei in zwei
Formen: 21 funktionale Anforderungen als Überschriften, 6 Qualitätsanforderungen
als Zeilen einer Tabelle. Je Konfiguration fehlt in der RTM eine der beiden
Gruppen; die Liste aus zwei Einträgen scheitert heute beim Laden der
Konfiguration (Exit 2). Das Lastenheft darf nicht umgeformt werden.

Die Entscheidung betrifft die **Form** des Schlüssels, weil sie jede künftige
Konfiguration bindet, und drei Härte-Fragen, die ADR-0037 für **einen** Eintrag
beantwortet hat und die je Eintrag neu zu beantworten sind.

## Entscheidung

### 1. Der Schlüssel nimmt ein Objekt oder eine Liste von Objekten an

`trace.requirements` bleibt ein Schlüssel. Er ist entweder das heutige Objekt
oder eine **nichtleere** Liste solcher Objekte. Das Objekt verhält sich wie
eine Liste mit einem Eintrag, jedoch ohne Index in den Meldungen und mit
byte-identischer Ausgabe. Jede andere YAML-Form und die leere Liste sind
Konfigurationsfehler (Exit 2). Jeder Eintrag trägt dieselben Schlüssel mit
derselben Bedeutung wie das Objekt; es entsteht kein zweites Schema.

### 2. Jeder Eintrag durchläuft den bestehenden Quellenvertrag für sich

Die Schritte von ADR-0037 — Format-Auflösung, Spaltenbindung, `strictSource`,
Duplikate innerhalb der Tabelle, Nullmenge — gelten **je Eintrag**; Meldungen
nennen den Index. `table.duplicate-ids` wirkt weiter nur innerhalb eines
Eintrags. Ein Eintrag, der nichts findet, verschwindet so nicht in der
Vereinigung.

### 3. Doppelte Kennung über Einträge ist ein Fehler, keine Politik

Erkennen zwei Einträge dieselbe Kennung, endet der Lauf mit Exit 2; die
Meldung nennt Kennung und beide Einträge. Es gibt keine `first`/`last`-Politik
über Einträge. Die Vereinigung ist die Zusage „jede Kennung zählt einmal";
eine stille Auflösung machte aus einer Überlappung der `id-pattern` einen
unsichtbaren Verlust, genau die Klasse, die ADR-0037 geschlossen hat. Die
Meldung weist auf eine mögliche Überlappung der Muster hin; ein
Subsumtions-Test der Regexe findet nicht statt, denn erst die gefundene
Kennung ist ein belegter Konflikt.

### 4. Reihenfolge und Modalität

Die Vereinigung folgt der Listenreihenfolge, innerhalb eines Eintrags der
bisherigen Sortierung — keine globale Neusortierung. `modality` gilt je
Eintrag und wird nicht vererbt; `require-levels` eines Eintrags gatet nur
dessen Anforderungen. Hat kein Eintrag `modality`, gaten wie bisher alle
Waisen. So kann die Modalität einer Gruppe im Body (Überschrift), der anderen
in einer Spalte liegen.

### 5. Die Referenzseite erkennt die Vereinigung der Muster

Referenzscan, Range-Validierung der Coverage und `cross-consistency` lesen
heute ein einziges `id-pattern`. Bei einer Liste gilt eine Kennung dort als
erkannt, wenn sie als Ganz-Token auf das Muster **mindestens eines** Eintrags
passt. Ein eigenes „Referenzmuster" entsteht nicht: es wäre eine zweite
Wahrheit neben den Mustern der Einträge, und die Referenzseite darf nichts
erkennen, was die Anforderungsseite nicht kennt.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| Neuer Schlüssel `trace.sources:` neben dem Objekt | Objekt-Schema bleibt unberührt | zwei Schlüssel für dasselbe, Vorrangregel nötig, `--print-config`/Handbuch doppelt |
| Tabelle vorab in Headings generieren | kein Produktcode | Generator + Drift-Sensor je Konsument; verlangt ein erzeugtes Artefakt neben dem Lastenheft |
| Mehrere `source`-Werte im Objekt (`source: [a, b]`) | kleine Änderung | trägt weder ein eigenes `format` noch ein eigenes `id-pattern` je Quelle — der Fall des Adopters (eine Datei, zwei Formen) bleibt ungelöst |
| Politik `first`/`last` über Einträge | deckt historische Überlappung | verdeckt Musterfehler; kein Adopter braucht sie |
| Gewählt: Objekt oder Liste, je Eintrag der bestehende Vertrag | ein Schlüssel, ein Schema, Objekt-Form byte-identisch | zwei erlaubte Formen für dasselbe; Meldungen unterscheiden Objekt und Liste |

**Fitness-Funktion:**

- Eine Fixture mit 21 Überschriften-Anforderungen und einer Tabelle mit 6
  Zeilen in derselben Datei, Liste aus zwei Einträgen, liefert in `--trace`
  21 + 6 Zeilen und die Summe 27.
- Eine Kennung als Überschrift **und** als Tabellenzeile, beide Einträge mit
  überlappendem `id-pattern`, endet mit Exit 2 und nennt Kennung und beide
  Einträge — die Prüfung wird rot, sobald die Kreuz-Duplikat-Prüfung entfällt.
- `requirements: []` und ein Eintrag ohne Treffer enden mit Exit 2.
- Die Objekt-Form liefert gegenüber der Fassung vor dieser ADR byte-identische
  Ausgabe.

## Konsequenzen

- Der Config-Adapter akzeptiert beide YAML-Formen und validiert Liste und
  Einträge vor dem I/O; ADR-0003 (striktes YAML) bleibt unberührt, die
  Form-Weiche wird getestet.
- Der App-Kern führt den Extraktor je Eintrag aus und bildet danach die
  Vereinigung samt Kreuz-Duplikat-Prüfung; Referenzscan, Coverage und Reporter
  sehen weiter genau eine Menge.
- Das Benutzerhandbuch zeigt die Listen-Form; `--print-config` und
  `--suggest-config` erzeugen weiter die Objekt-Form — eine Liste ist
  Handarbeit für den gemischten Bestand.
- **Grenze:** Ein Listen-Eintrag mit `format: headings` und ohne explizite
  `source` ist wie in der Objekt-Form nicht strikt; null Treffer ergeben dort
  keinen Fehler. Das ist Parität mit ADR-0037, kein Schutz — wer den Eintrag
  nicht still leer haben will, benennt `source` oder wählt `format: table`.

## Re-Evaluierungs-Trigger

- Ein Adopter braucht eine Auflösung doppelter Kennungen über Einträge hinweg.
- Ein Eintrag soll Felder von einem anderen erben (etwa `modality`).
- Die Fehlerbilder von Objekt- und Listen-Form driften auseinander.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-10-20 | Entwurf nach eingehendem CR eines Adopters (eine Datei, zwei Anforderungsformen). Lastenheft 0.106.0 und Spezifikation tragen die Anforderung; Status Proposed, bis die Umsetzung steht. |
| 2026-10-20 | **Accepted** — umgesetzt (Config-Adapter, Kern, Handbuch), zwei Review-Runden und Verifikation durchlaufen. Vor der Annahme berichtigt: die Konsequenz zu `--print-config` (die Generatoren bleiben bei der Objekt-Form). Die Referenzseite (Entscheidung 5) kam beim Umsetzen hinzu, nachdem Referenzscan, Range-Validierung und `cross-consistency` als Verbraucher eines einzigen Musters auffielen. |
