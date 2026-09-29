# Spezifikation — <Projektname>

> **Template-Hinweis.** Diese Datei ist eine Vorlage. Sie ist
> **technisch verbindlich, aber ohne Lastenheft-Änderung fortschreibbar**
> (siehe Spec-Stratifizierung in
> [Baseline-Regelwerk §Spec-Stratifizierung](../../regelwerk/grundlagen-source-precedence.md#spec-stratifizierung)).
> Kopiere sie nach `spec/spezifikation.md`, ersetze `<Platzhalter>` und
> lösche diesen Block.

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
§ID-Schema als Klammer). Was **keine** einzelne Lastenheft-ID verfeinert,
trägt eine `SPEC-<NNN>` — siehe §2 bis §6.

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

Regeln dieser Sektion: Jede Struktur trägt eine `SPEC-<NNN>` — eine Adresse,
keine Anforderung (Baseline-Regelwerk `grundlagen-source-precedence.md`
§ID-Schema als Klammer). Gezählt wird fortlaufend je Datei, nicht je Sektion
(Baseline-Regelwerk `grundlagen-source-precedence.md` §Vergabe).

<!--
Konkrete JSON-Schemas, OpenAPI-Snippets, Protokoll-Definitionen.
Fortschreibbar ohne Lastenheft-Änderung, solange die Lastenheft-
Anforderung gewahrt bleibt.
-->

### SPEC-001 — <Datenstruktur 1>

```json
{
  "field": "type"
}
```

## 3. Defaults und Konstanten

Regeln dieser Sektion: Die ADR, die einen Wert festlegt, deklariert das
aufwärts in ihrem `Schärft:`-Feld — kein ADR-Rückzeiger hier
(Baseline-Regelwerk `grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)).
Die `SPEC-<NNN>` ist das, was ihr `Schärft:`-Feld **benennt**; ohne sie kann
eine ADR nur den ganzen Abschnitt nennen. Der Link zeigt weiter auf die
Sektion — Kennungen in Tabellenzellen haben keinen eigenen Anker, und kein
Sensor bemerkt, wenn eine umbenannt wird.

<!-- Werte, die in Code fest sind. -->

| ID | Name | Wert | Begründung |
|---|---|---|---|
| `SPEC-002` | `MAX_BATCH_SIZE` | 100 | <…> |

## 4. Fehler-Codes und Logging-Felder

Regeln dieser Sektion: Der Fehler-Code ist ein Laufzeit-Symbol, die
`SPEC-<NNN>` benennt die *Festlegung* darüber — beide stehen nebeneinander,
sonst hätte eine ADR, die die Fehlerbehandlung schärft, kein Ziel
(Baseline-Regelwerk `grundlagen-source-precedence.md` §ID-Schema als Klammer).

<!-- Verbindliche Codes und Felder. -->

| ID | Code | Bedingung | Aktion |
|---|---|---|---|
| `SPEC-003` | E001 | <…> | <…> |

## 5. Metriken und Tracing-Felder

Regeln dieser Sektion: verbindliche OTel-Felder pro Span
(Baseline-Regelwerk `modul-15-observability.md`).

| ID | Span | Pflicht-Attribute | Quelle |
|---|---|---|---|
| `SPEC-004` | `<service>.<operation>` | `<feldname>`, `<feldname>` | <…> |

## 6. Externe Verträge

<!-- Schnittstellen zu Drittsystemen, mit Versionsannahme. -->

| ID | System | Version | Vertrag-Datei |
|---|---|---|---|
| `SPEC-005` | <…> | <…> | <Pfad> |

## 7. Historie

Regeln dieser Sektion: **kein ADR- und kein Slice-Verweis.** Die Decken-Regel
gilt für alle drei Spec-Straten, auch hier — welche ADR eine Festlegung
schärft, deklariert die ADR aufwärts in ihrem `Schärft:`-Feld
(Baseline-Regelwerk `modul-03-spec.md` §Ziel-Form: Spezifikation).

| Datum | Änderung |
|---|---|
| YYYY-MM-DD | Initial |
