# Eingehender Change Request — mehrere Anforderungsquellen bzw. -formen in `trace.requirements`

**Absender:** Adopter `sf-connector` · **Eingegangen:** 2026-10-20
**Richtung:** eingehend — dieses Repo ist der **Empfänger**, nicht der Bittsteller.
**Ziel-Dokument:** [`spec/lastenheft.md`](../../../spec/lastenheft.md)
**Berührt:** [`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
(RTM, `--trace`), [`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
(`--require-complete`)
**Stand:** **angenommen** — Lastenheft-Änderung in Version 0.106.0
([`DC-FA-REQ-001`](../../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen)).
Spezifikation, ADR und Umsetzung stehen aus. Entscheidung je Punkt unten.

**Ablage-Hinweis.** Ein **eingehender** CR liegt neben
[`MR-035`](../../../harness/conventions.md#mr-035) (ausgehend) und
[`MR-036`](../../../harness/conventions.md#mr-036) (Antworten darauf). Der
Wortlaut steht unten sinngemäß unverändert; Kennungen des Absender-Repos tragen
ein `d-check:ignore`, die doppelt eingefügte Tabellenzeile der Vorlage ist
einmal geführt.

---

## Wortlaut

> **Change Request an d-check: mehrere Anforderungsquellen bzw. -formen in `trace.requirements`**
>
> Entwurf — eingebracht vom Nutzer von sf-connector. Stand der Messung: d-check v0.86.1 (neuestes Release, 2026-10-10), Image
> `ghcr.io/pt9912/d-check@sha256:3e0b9779a71e2fa942960e8b961428fba036455c535513d34a75bfb799ffce0e`,
> `--network none`, gemessen am 2026-10-20 an einer Kopie des Lastenhefts von sf-connector;
> Benutzerhandbuch gelesen auf `origin/main` (`be0737ee`). Nichts im Repo von d-check geändert.
> Anfragender Nutzer: sf-connector. Dort trägt `make doc-trace` die Matrix
> (`docker run … --trace`); sie ist Werkzeug, kein Gate.
>
> **Anlass**
>
> Das Lastenheft von sf-connector führt seine Anforderungen in zwei Formen in einer Datei:
>
> - 21 funktionale Anforderungen als ATX-Überschrift `#### SFC-FA-<NN> — <Titel>` (Body darunter), <!-- d-check:ignore (Kennungs-Schema des Absender-Repos) -->
> - 6 Qualitätsanforderungen als Zeilen einer Pipe-Tabelle mit den Spalten
>   `| ID | Aspekt | Anforderung | Priorität |`.
>
> Das Lastenheft ist vertraglich abnahmebindend und wird nicht umgeformt, auch nicht in der Form.
> `trace.requirements` nimmt genau eine Quelle mit genau einem `format` an; damit fehlt in
> `--trace` je nach Konfiguration eine der beiden Gruppen.
>
> **Ist**
>
> | Messung | Konfiguration `trace.requirements` | Ergebnis |
> | --- | --- | --- |
> | I-1 | `source: spec/lastenheft.md`, `id-pattern: 'SFC-(?:FA\|QA\|RB)-\d{2}'`, `format: headings` | Exit 0, 21 Anforderung(en): alle FA, kein QA |
> | I-2 | dieselbe Quelle, `format: table`, `table: {id-column: ID, text-column: Anforderung, modality-column: Priorität}` | Exit 0, 6 Anforderung(en): alle QA, kein FA |
> | I-3 | `requirements` als Liste aus zwei Einträgen (I-1 mit `SFC-FA-\d{2}`, I-2 mit `SFC-QA-\d{2}`) | Exit 2: `d-check: error: .d-check.yml: yaml: unmarshal errors: line 3: cannot unmarshal !!seq into configyaml.rawTraceRequirements` |
>
> Das Benutzerhandbuch (§5, `requirements.format`) beschreibt genau das: „`headings` (Default) oder
> `table`"; für gemischte Bestände empfiehlt es, die Tabelle in Überschriften zu überführen oder ein
> separates Heading-Dokument zu erzeugen, dessen Drift ein eigener Sensor absichert. Beides kann
> sf-connector nicht bzw. nur mit einem zusätzlichen erzeugten Artefakt samt Sensor.
>
> **Soll**
>
> **S-1 — Liste von Anforderungsquellen.** `trace.requirements` darf statt eines Objekts eine
> nichtleere Liste solcher Objekte sein. Jeder Eintrag hat die heutigen Felder (`source`,
> `id-pattern`, `format`, bei `table` die `table.*`-Felder, `modality`) mit derselben Bedeutung und
> denselben Fehlerbildern wie heute. Dieselbe `source` darf in mehreren Einträgen stehen.
>
> ```yaml
> trace:
>   requirements:
>     - source: spec/lastenheft.md
>       id-pattern: 'SFC-FA-\d{2}'
>       format: headings
>     - source: spec/lastenheft.md
>       id-pattern: 'SFC-QA-\d{2}'
>       format: table
>       table:
>         id-column: ID
>         text-column: Anforderung
>         modality-column: Priorität
> ```
>
> **S-2 — Vereinigung und Reihenfolge.** Die RTM führt die Vereinigung aller Einträge. Reihenfolge:
> die der Einträge in der Liste, innerhalb eines Eintrags die heutige. Die Summenzeile zählt über
> alle Einträge (`27 Anforderung(en), …`).
>
> **S-3 — Doppelte Kennung über Einträge hinweg.** Erkennen zwei Einträge dieselbe Kennung, ist das
> ein Konfigurationsfehler (Exit 2, fail-closed) mit Kennung und beiden Einträgen in der Meldung —
> dieselbe Härte wie `table.duplicate-ids: error` innerhalb einer Tabelle. (Eine
> `first`/`last`-Politik über Einträge braucht sf-connector nicht.)
>
> **S-4 — Nullmengen-Guard je Eintrag.** Der heutige Guard (nichtleer explizite `source` oder
> `format: table` ⇒ fehlende Datei oder 0 Treffer ⇒ Exit 2) gilt je Eintrag, mit Index des
> Eintrags in der Meldung. Ein Eintrag, der nichts findet, verschwindet so nicht still in der
> Vereinigung.
>
> **S-5 — Rückwärtskompatibilität.** Die heutige Objekt-Form bleibt gültig und ergibt
> byte-identische Ausgabe; ebenso ohne `trace`-Block. `--require-complete` (`make doc-complete`)
> wertet die Vereinigung aus.
>
> **Akzeptanzkriterien**
>
> Jeweils an einer Kopie mit dem Lastenheft von sf-connector (oder einem Fixture derselben Form:
> Überschriften `#### X-FA-NN — …` und eine Tabelle `| ID | Aspekt | Anforderung | Priorität |` in
> derselben Datei):
>
> | AK | Befehl / Konfiguration | erwartet |
> | --- | --- | --- |
> | AK-1 | Konfiguration aus S-1, `--trace` | Exit 0; 21 Zeilen `SFC-FA-*` gefolgt von 6 Zeilen `SFC-QA-*`; Summenzeile 27 Anforderung(en) |
> | AK-2 | Konfiguration aus S-1, beide Einträge mit `id-pattern: 'SFC-(?:FA\|QA)-\d{2}'` und eine Kennung, die als Überschrift und als Tabellenzeile steht | Exit 2, Meldung nennt die Kennung und beide Einträge (S-3) |
> | AK-3 | Konfiguration aus S-1, zweiter Eintrag mit `id-column: Kennung` (gibt es nicht) | Exit 2 wie heute bei fehlendem Header, Meldung nennt den Eintrag |
> | AK-4 | Konfiguration aus S-1, zweiter Eintrag mit `id-pattern: 'SFC-RB-\d{2}'` (0 Treffer) | Exit 2, Nullmengen-Meldung nennt den Eintrag (S-4) |
> | AK-5 | heutige Objekt-Form (I-1 bzw. I-2) | Ausgabe byte-identisch zu v0.86.1 |
> | AK-6 | Konfiguration aus S-1, `--trace --require-complete` | Exit 1, solange eine der 27 ohne Slice- und ohne Coverage-Referenz ist; Exit 0, wenn keine |
> | AK-7 | `requirements: []` | Exit 2 (leere Liste ist Konfigurationsfehler) |
>
> Dazu: das Benutzerhandbuch (§5 `trace`, Konfigurationsbeispiel im Abschnitt zur `.d-check.yml`)
> beschreibt die Listen-Form, S-3 und S-4.

---

## Entscheidung

Angenommen, S-1 bis S-5 und AK-1 bis AK-7 vollständig.

- **S-1, S-2:** übernommen. Die Objekt-Form bleibt gültig.
- **S-3:** übernommen, Exit 2 auch bei überlappenden `id-pattern`; eine
  Politik `first`/`last` über Einträge hinweg ist ausdrücklich ausgeschlossen.
- **S-4:** übernommen, je Eintrag mit Index in der Meldung.
- **S-5:** übernommen.
- **`modality`:** je Eintrag, ohne Vererbung. Das Lastenheft nennt es in
  [`DC-FA-REQ-001`](../../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen); der Wortlaut von [`DC-FA-MOD-001`](../../../spec/lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in) bleibt unverändert.
- **Ausschluss entfällt:** der bisherige Out-of-Scope-Satz „genau ein `format`
  / mehrere Quelldateien" widersprach dem CR und ist ersetzt.
- **Offen:** Spezifikation (Schlüsseltabelle `trace.requirements.*`), ADR
  (Form Objekt-oder-Liste, Schärfung von [ADR-0037](../adr/0037-trace-tabellenquellen-nullmengen-guard.md)), Slice, Handbuch (in der
  Release-Prep).
