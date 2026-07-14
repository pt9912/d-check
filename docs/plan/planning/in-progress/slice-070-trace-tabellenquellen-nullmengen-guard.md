# Slice slice-070: Trace-Tabellenquellen und Nullmengen-Guard

**Status:** in-progress (welle-59-trace-tabellenquellen).

**Welle:** aktiv; Vorgänger
[`slice-069`](../done/slice-069-trace-handbuch-parsergrenzen.md) ist
abgeschlossen.

**Bezug:** neuer Lastenheft-Change-Request
[`DC-FA-REQ-001`](../../../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen)
(Lastenheft 0.43.0), Mit-Änderung
[`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
und Interaktion mit
[`DC-FA-MOD-001`](../../../../spec/lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in).
Produkt-/Config-Delta; Spezifikation und begründende ADR folgen doc-first vor
der Implementierung. Release v0.43.0 geplant.

**Autor:** pt9912. **Datum:** 2026-07-14.

---

## 1. Ziel

`--trace` erkennt in d-check v0.42.0 Anforderungen ausschließlich aus
ATX-Überschriften. Eine explizite Quelle und eine passende Kennungs-Regex
reichen deshalb für tabellenbasierte Lastenhefte nicht aus: Im Konsumenten
`m-trace` ergeben 371 Zeilen der Tabelle
`Kennung | Prioritaet | Anforderung` eine leere RTM und selbst
`--require-complete` endet mit Exit 0.

Der Slice liefert zwei zusammengehörige Korrekturen: eine native, über
Header-Namen konfigurierte Tabellenquelle und einen fail-closed
Nullmengen-Guard, sobald der Nutzer eine Anforderungsquelle oder den
Tabellenmodus explizit auswählt.

## 2. Vertrag aus dem Change Request

- **Zwei explizite Formate:** `headings` bleibt Default;
  `trace.requirements.format: table` liest Markdown-Pipe-Tabellen.
- **Spalten nach Header:** `table.id-column` und genau eine von
  `table.text-column` oder `table.text-columns` sind im Tabellenmodus Pflicht;
  die Liste bildet explizite alternative Text-Header ab.
  `table.modality-column` ist optional. Keine positionsbasierte oder
  heuristische Spaltenwahl.
- **Ganzzellen-ID:** Nur eine ID-Zelle, die vollständig auf
  `requirements.id-pattern` passt, definiert eine Anforderung.
- **Modalitätsquelle:** Bei gesetzter `modality-column` wird ausschließlich
  diese Zelle klassifiziert; sonst die Textzelle. Keywords und
  `require-levels` bleiben die Semantik von
  [`DC-FA-MOD-001`](../../../../spec/lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in).
- **Gemeinsames RTM-Modell:** Nach der Extraktion bleiben ADR-/Slice-/Coverage-
  Scan, Waisenstatus und Reporter unverändert.
- **Duplikatpolitik:** `table.duplicate-ids` ist `error` (Default), `first` oder
  `last`. Nur die beiden expliziten Override-Werte erlauben historische
  Mehrfachdefinitionen deterministisch; der sichere Default bleibt
  fail-closed.
- **Fail-closed:** Nichtleer explizite `requirements.source` oder Tabellenmodus
  plus null erkannte Anforderungen ⇒ Exit 2. Unbekanntes Format,
  fehlende/leere/nicht gefundene oder doppelte konfigurierte Header sowie
  doppelte IDs im Tabellenmodus oder bei nichtleer expliziter Quelle ⇒ Exit 2.
- **Empty bleibt Default:** `requirements.source: ""` gilt wie Abwesenheit und
  aktiviert weder Nullmengen- noch Duplicate-ID-Guard.
- **Kompatibilität:** Ohne `trace`-Block bleibt die Heading-Ausgabe inklusive
  der bisherigen unkonfigurierten Nullmenge byte-identisch.

## 3. Definition of Done

- [x] **Lastenheft-CR:**
  [`DC-FA-REQ-001`](../../../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen)
  mit Happy/Boundary/Negative,
  Out-of-Scope, Bereich `REQ`, Version 0.43.0 und Historie angelegt;
  [`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
  auf die Definitionsformate geschärft.
- [x] **Spezifikation:** Algorithmusabschnitt zur neuen Anforderung,
  Tabellenlexik/-Escaping, alternative Text-Header, Duplikatpolitik,
  Header-Normalisierung, Fehlerpräzedenz und Config-Schema für `format` /
  `table.*` festlegen; Spezifikationshistorie ergänzen.
- [x] **ADR + Index:** Parser-/Config-Entscheidung, Default-Kompatibilität und
  Nullmengen-Präzedenz begründen; Status zunächst Proposed.
- [x] **Modell/Config:** Format und Spaltennamen abbilden, vollständig
  validieren und in `--print-config` sichtbar machen.
- [x] **Parser:** Markdown-Pipe-Tabellen deterministisch/read-only extrahieren;
  dieselbe `RequirementTrace`-Form wie der Heading-Pfad liefern.
- [x] **Nullmengen-Guard:** explizite Quelle bzw. Tabellenmodus von Default-
  Fallback unterscheiden; Exit 2 vor Reporter/Gate-Auswertung.
- [x] **Tests:** 371-Zeilen-Happy-Path inklusive zweier Text-Header und
  historischem Duplikat, Modalitätsspalte, Text-Fallback,
  null Treffer, fehlende Header, doppelte IDs, Escapes/Ränder sowie
  byte-identischer Heading-Default.
- [x] **Realdatenbeleg `m-trace`:** mit dessen Konfiguration exakt 371
  Anforderungen statt 0; ein absichtlich falscher Spaltenname und eine
  unpassende Regex brechen jeweils fail-closed.
- [x] **Nutzerdoku:** Handbuch-Brownfield-Hinweis durch native Config ergänzt,
  Warnung auf Versionsgrenze aktualisiert sowie Changelog/Operations gepflegt.
- [x] **Release:** Versionsregister/Release-Prep, Tag und GHCR-Release v0.43.0
  samt Digest-Backfill abgeschlossen (Release-Run 29340521688).
- [x] **Qualität:** Self-Reviews liegen vor; `make gates` und `make ci` sind
  grün. Für Closure bleibt ein personell unabhängiger Review erforderlich.

## 4. Risiken / offene Designpunkte

- **Markdown-Tabellen sind kein CSV:** escaped Pipes, Code-Spans und
  führende/abschließende Pipes brauchen eine festgeschriebene kleine Grammatik;
  Block-Markdown und mehrzeilige Zellen bleiben Out-of-Scope.
- **Header-Kollisionen:** doppelte gleichnamige Header müssen fail-closed sein,
  sonst wäre die Spaltenwahl mehrdeutig.
- **Doppelte IDs:** Anders als ein stilles Last-write-wins muss die Quelle als
  fehlerhaft abgewiesen werden, sofern Tabellenmodus oder nichtleer explizite
  Quelle den fail-closed Vertrag aktiviert. Nur die explizite Tabellenpolitik
  `first`/`last` löst historische Mehrfachdefinitionen; der unkonfigurierte
  Heading-Pfad behält seine bestehende Deduplizierung.
- **Fehlerpräzedenz:** Config-Schemafehler vor Dateilesen, Tabellenstruktur vor
  Nullmengen-Guard; eine stabile Reihenfolge ist für deterministische Diagnose
  zu spezifizieren.
- **Lifecycle:** `slice-069` ist nach ACCEPT-Folgereview geschlossen; dieser
  Slice kann nun regulär über `next/` nach `in-progress/` aufgenommen werden.

## 5. Trigger

Auftraggeber-Befund 2026-07-14 in `m-trace`: d-check v0.42.0 liest aus dessen
explizit konfiguriertem `spec/lastenheft.md` trotz passender
`(F|NF|MVP|AK|RAK)-[0-9]+`-Regex keine der 371 tabellarischen Anforderungen.
Die leere RTM endet mit Exit 0 und wäre als Vollständigkeits-Gate irreführend
grün.

## 6. Sub-Area-Modus-Begründung

GF für die additive `trace.requirements`-Config und den zweiten
Requirements-Parserpfad: Der neue Vertrag führt, Implementierung folgt. Der
bestehende Heading-Pfad ist Kompatibilitätsbaseline und wird durch
byte-identische Tests gegen Regression geschützt.

## 7. Realdatenbeleg

Read-only-Lauf des lokalen Runtime-Images gegen eine Kopie des
`m-trace`-Arbeitsstands unter `/tmp/m-trace-slice070`:

- gültige Config mit `text-columns: [Anforderung, Akzeptanzkriterium]` und
  `duplicate-ids: last`: `371 Anforderung(en), 98 Waise(n).`, Exit 0;
- Tippfehler nur im zweiten Text-Header (`Akzeptanzkriteriuum`): Exit 2 mit
  `konfigurierter Text-Header … kommt in keiner Tabelle … vor`;
- unpassende ID-Regex `ZZZ-[0-9]+`: Exit 2 mit `ergab 0 Anforderungen`.

Das Original-Repository wurde nicht verändert; der Container mountete nur die
Kopie read-only und lief mit `--network none`.

Das publizierte v0.43.0-Digest-Image wurde zusätzlich mit einer read-only
Tabellen-Fixture smoke-verifiziert: zwei Anforderungen wurden aus den Spalten
`Kennung | Prioritaet | Anforderung` erkannt und als `must` beziehungsweise
`may` klassifiziert. Release-Run 29340521688 war grün; Digest-Pin:
`ghcr.io/pt9912/d-check@sha256:2963f882c40a0b34d1fc03ba0e91feaf18423e55a35084bead1efa9d5500bd53`.

## 8. Closure-Notiz (nach `done/`)

Implementierung, Realdatenbeleg, Gates, Release v0.43.0 und Digest-Backfill sind
abgeschlossen. Offen bleiben das personell unabhängige Review und danach der
zweistufige Lifecycle-Move nach `done/`; bis dahin bleibt der Slice ehrlich in
`in-progress/`.
