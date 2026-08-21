# Lastenheft — <Projektname>

> **Template-Hinweis.** Diese Datei ist eine Vorlage. Kopiere sie nach
> `spec/lastenheft.md` deines Repos und ersetze alle `<Platzhalter>`.
> Lösche diesen Hinweis-Block und alle Kommentar-Zeilen (`<!-- -->`)
> nach dem Ausfüllen.

**Version:** 0.1.0 (`Major.Minor.Patch` — erhöhen bei Änderungen, siehe
Baseline-Regelwerk `modul-03-spec.md`).

**Status:** Draft | In Review | Accepted (entwurfsstatus, der die
Verbindlichkeit der IDs steuert).

**Autor:** <Name>, **Datum:** YYYY-MM-DD.

---

## 1. Zweck und Geltungsbereich

<!--
Ein bis zwei Absätze: Was leistet das System, für wen, gegen welche
Annahme. Konkret, aber nicht implementierungs-nah. "Wir bauen einen
Service, der …" ist OK.

Nicht hier: wie das System gebaut wird (das gehört in spezifikation.md
oder die ADRs).
-->

## 2. Stakeholder

<!--
Wer hat ein Interesse am Ergebnis? Pro Stakeholder: Rolle, Erwartung
in einem Satz.
-->

| Stakeholder | Rolle | Erwartung |
|---|---|---|
| <Beispiel: Vertrieb> | <Auftraggeber> | <Verkürzte Time-to-Market> |
| | | |

## 3. Funktionale Anforderungen

Regeln dieser Sektion: ID-Schema `<PREFIX>-FA-<NN>` — mit Bereichssegment
`<PREFIX>-FA-<BEREICH>-<NNN>`, wenn dein Repo den Zählraum je Sub-Area führt
(Baseline-Regelwerk `grundlagen-source-precedence.md` §Vergabe). Das Präfix ist im ganzen
Repo dasselbe und taucht in Make-Target-Kommentaren, ADRs und Commits wieder auf
(Baseline-Regelwerk `grundlagen-source-precedence.md` §ID-Schema als Klammer). Jede
Anforderung trägt drei Pfade — Happy · Boundary · Negative — plus Out-of-Scope
(Baseline-Regelwerk `modul-03-spec.md` §Ziel-Form: Akzeptanzkriterium).

<!-- Format: ID — Titel — Beschreibung — Akzeptanzkriterien
     (Given/When/Then, Boundary, Negative). -->

### LH-FA-01 — <Titel der Anforderung>

**Beschreibung:** <Was muss das System leisten?>

**Akzeptanzkriterien:**

- **Happy Path:** Given <Vorbedingung>, when <Aktion>, then <Erwartung>.
- **Boundary:** Given <Randfall>, when <Aktion>, then <definiertes Verhalten>.
- **Negative:** Given <ungültige Eingabe>, when <Aktion>, then <expliziter Fehlerpfad>.

**Out-of-Scope:** <Was explizit nicht gefordert ist.>

---

### LH-FA-02 — <…>

<!-- Weitere Anforderungen analog. -->

---

## 4. Nichtfunktionale Anforderungen

<!--
Format: ID — Kategorie — messbare Anforderung — Messmethode.

ID-Schema: <PREFIX>-QA-<NN>.

Kategorien (typische): Performance, Skalierbarkeit, Verfügbarkeit,
Sicherheit, Wartbarkeit, Betriebskosten.
-->

### LH-QA-01 — <Performance>

- **Anforderung:** <z.B. p95-Latenz < 200 ms bei 100 RPS.>
- **Messmethode:** <z.B. Lasttest unter Standardlast, definiert in spec/spezifikation.md.>

### LH-QA-02 — <Sicherheit>

- **Anforderung:** <…>
- **Messmethode:** <…>

---

## 5. Globale Out-of-Scope-Punkte

<!--
Explizite Nicht-Anforderungen, die für das Gesamtsystem gelten.
Ohne diesen Abschnitt baut der Agent gerne Plausibles.
-->

- <Beispiel: Multi-Mandanten-Fähigkeit ist nicht Teil der ersten Version.>
- <Beispiel: Keine Echtzeit-Streaming-API.>

## 6. Glossar

<!-- Begriffe, die im Lastenheft präzise verwendet werden. -->

| Begriff | Bedeutung im Lastenheft |
|---|---|
| <Begriff 1> | <Definition> |

## 7. Historie

Regeln dieser Sektion: Jede Änderung an *angenommenen* Anforderungen ist eine
Vertragsänderung. Sie entsteht **nur** aus einem Change Request, nie aus einem
ADR oder Slice. Fußabdruck pro angenommenem CR: Version-Bump oben plus eine
Zeile hier, mit dem CR unter „Verweis" (Baseline-Regelwerk
`grundlagen-source-precedence.md` §Spec-Stratifizierung).

Sind Auftraggeber und Entwickler **dieselbe Person**, fehlt nur die
Ticket-Form, nicht der Vorgang. Dann trägt der **Commit** die Trennung: Die
Änderung an dieser Datei liegt in einem eigenen Commit, **vor** dem Slice, der
sie umsetzt — und „Verweis" nennt diesen Vorgang statt eines Tickets.

**Auch hier gilt die Decken-Regel:** keine ADR, kein Slice, kein Carveout,
keine Welle, kein Verweis auf `spezifikation.md` oder `architecture.md` — in
keiner Spalte. Kein Spec-Stratum nimmt seine Historie davon aus. Die Spalte
„Verweis" trägt den **externen** CR; der steht außerhalb des Repos und damit
außerhalb der Referenz-Richtung. Wer im Repo bemerkt hat, dass eine Änderung
nötig wird, hält das auf seiner Seite fest — etwa in der Closure-Notiz des
Slice (Baseline-Regelwerk `modul-03-spec.md`
§Ziel-Form: Akzeptanzkriterium).

| Version | Datum | Änderung | Verweis |
|---|---|---|---|
| 0.1.0 | YYYY-MM-DD | Initiale Fassung | — |
