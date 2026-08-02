# Spezifikation — <Projektname>

> **Template-Hinweis.** Diese Datei ist eine Vorlage. Sie ist
> **technisch verbindlich, aber ohne Lastenheft-Änderung fortschreibbar**
> (siehe Spec-Stratifizierung in
> [Baseline-Regelwerk §Spec-Stratifizierung](../../regelwerk/grundlagen-source-precedence.md#spec-stratifizierung)).
> Kopiere sie nach `spec/spezifikation.md`, ersetze `<Platzhalter>` und
> lösche diesen Block.

**Status:** Aktiv. **Letzte Änderung:** YYYY-MM-DD.

**Bezug zum Lastenheft:** Diese Spezifikation präzisiert die in
`spec/lastenheft.md` formulierten Anforderungen (`LH-*`-IDs). Bei
Konflikt gewinnt das Lastenheft — präzisieren ja, erweitern nie.

**Rolle:** Technik-Stratum — fortschreibbar ohne Change Request; eine ADR darf
sie schärfen, das Lastenheft nicht. Regeln: Baseline-Regelwerk
`modul-03-spec.md` §Ziel-Form: Spezifikation.

---

## 1. Algorithmen und Datenflüsse

Regeln dieser Sektion: Tatsächlicher Code gehört in `src/`, nicht hierher.
ID-Schema `<PREFIX>-FA-<NN>.<Buchstabe>` für Verfeinerungen einzelner
Lastenheft-IDs (Baseline-Regelwerk `grundlagen-source-precedence.md`
§ID-Schema als Klammer).

<!--
Wie wird die funktionale Anforderung *technisch* erfüllt? Pseudocode
oder Sequenzbeschreibung erlaubt; tatsächlicher Code gehört in
src/.

ID-Schema: <PREFIX>-FA-<NN>.<Buchstabe> für Verfeinerungen einzelner
Lastenheft-IDs, z.B. LH-FA-03.a für eine konkrete Algorithmus-Variante.
-->

### LH-FA-01.a — Algorithmus für <…>

**Eingabe:** <…>. **Ausgabe:** <…>. **Schritte:**

1. <…>
2. <…>

**Komplexität:** <O(n log n)>, **Fehlermodi:** <…>.

---

## 2. Datenstrukturen und Schemas

<!--
Konkrete JSON-Schemas, OpenAPI-Snippets, Protokoll-Definitionen.
Fortschreibbar ohne Lastenheft-Änderung, solange die Lastenheft-
Anforderung gewahrt bleibt.
-->

### <Datenstruktur 1>

```json
{
  "field": "type"
}
```

## 3. Defaults und Konstanten

Regeln dieser Sektion: Die ADR, die einen Wert festlegt, deklariert das
aufwärts in ihrem `Schärft:`-Feld — kein ADR-Rückzeiger hier
(Baseline-Regelwerk `grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)).

<!-- Werte, die in Code fest sind. -->

| Name | Wert | Begründung |
|---|---|---|
| `MAX_BATCH_SIZE` | 100 | <…> |

## 4. Fehler-Codes und Logging-Felder

<!-- Verbindliche Codes und Felder. -->

| Code | Bedingung | Aktion |
|---|---|---|
| E001 | <…> | <…> |

## 5. Metriken und Tracing-Felder

Regeln dieser Sektion: verbindliche OTel-Felder pro Span
(Baseline-Regelwerk `modul-15-observability.md`).

| Span | Pflicht-Attribute | Quelle |
|---|---|---|
| `<service>.<operation>` | `<feldname>`, `<feldname>` | <…> |

## 6. Externe Verträge

<!-- Schnittstellen zu Drittsystemen, mit Versionsannahme. -->

| System | Version | Vertrag-Datei |
|---|---|---|
| <…> | <…> | <Pfad> |

## 7. Historie

Regeln dieser Sektion: **kein ADR- und kein Slice-Verweis.** Die Decken-Regel
gilt für alle drei Spec-Straten, auch hier — welche ADR eine Festlegung
schärft, deklariert die ADR aufwärts in ihrem `Schärft:`-Feld
(Baseline-Regelwerk `modul-03-spec.md` §Ziel-Form: Spezifikation).

| Version | Datum | Änderung |
|---|---|---|
| 0.1.0 | YYYY-MM-DD | Initial |
