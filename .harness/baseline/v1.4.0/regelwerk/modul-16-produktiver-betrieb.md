## Modul 16 — Produktiver Betrieb

*Quelle: [05-betrieb/modul-16-produktiver-betrieb.md](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/05-betrieb/modul-16-produktiver-betrieb.md)*

### Kernidee (Modul 16)

Produktiv heißt: Du musst eine Frage in der Nacht beantworten können,
ohne den Autor zu kennen. Runbooks und Replay sind dafür da.

### Regeln gegen typische Fehlannahmen (Modul 16)

- **"Rollback ist die Standardantwort."** — Drei Fälle, in denen Rollback schadet: nicht-rückwärtskompatible DB-Migration, bereits erzeugte Buggy-Daten, ungetesteter Rollback-Pfad. Runbook entscheidet *vor* dem Incident, wann Fix-Forward gilt.
- **"Runbook beschreibt den Happy Path."** — Nein. Runbook beschreibt *Entscheidungen unter Unsicherheit*, mit Triggern. Wenn das Runbook nur sagt "Service neu starten", ist es kein Runbook.
- **"Produktionsfreigabe ist eine formale Checkbox."** — Eine Checkliste ohne *Belege* pro Item (Replay-Lauf-Link, ADR-ID, Trace-Hash) ist Bürokratie. Mit Belegen ist sie das einzige nicht-fragmentierte Audit-Artefakt.
- **"Deployt heißt produktiv."** — Nein. Deployment ist *eine* Anwendung des Container-Ankers (Modul 14), nicht das Ziel. Produktionsreife heißt *belegte Betriebsfähigkeit*: kann ein anderer Mensch nachts handeln (Runbook), ist der Lauf reproduzierbar (Replay-Beleg), entfällt die Freigabe bei einem Incident automatisch (Incident-Klausel)? Ein Service kann längst deployt und trotzdem nicht produktionsreif sein — genau diese Lücke schließt die Freigabe-Checkliste.
- **"Prompt-Injection ist eine Modell-Frage."** — Nein. Erkennung von Injection ist eine *Telemetrie-Frage*: Eingabe-Logging + Tool-Call-Audit + Output-Drift-Marker. Wer das nicht hat, erkennt Injection nur durch Glück.
- **"Postmortem ist Schuldzuweisung — also macht man's leise."** — Genau das Gegenteil. Ein produktiver Postmortem ist *blameless* (vgl. Etsy/Google SRE-Tradition): er sucht den Pfad, auf dem ein vernünftiger Mensch unter Druck dieselbe Entscheidung getroffen hätte, und fragt, *welcher Sensor oder Guide gefehlt hat*. Closure-Einträge in `done/` ([Modul 5](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/02-planung/modul-05-planning-harness.md)) und Reflexions-Einträge ([`grundlagen/reflexion-vorlage.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/reflexion-vorlage.md)) sind beide *strukturell* blameless: sie fragen "welche Harness-Lücke war Ursache", nicht "wer war es". Wer Postmortems als Schuldzuweisung erlebt hat, wird Drift-Symptome zukünftig verschweigen — und genau dadurch wachsen sie. Blameless ist keine moralische Wahl; es ist eine Sensor-Schutz-Maßnahme.

### Worked Example: eine Produktionsfreigabe-Checkliste schreiben

**Ausgangs-Situation:** Du sollst die Freigabe-Checkliste für Welle 1
deines Projekts schreiben. Das Repo hat: ein abgeschlossenes Slice in
`done/`, eine ADR, einen Carveout, ein Replay-Set, ein Trace-Fixture.

**Schritt 1 — Item-Form festlegen: keine Häkchen ohne Beleg.**
Schlechtes Format:
```markdown
- [ ] Tests grün.
- [ ] Replay gelaufen.
```
Gutes Format — jedes Item trägt einen **Beleg-Slot**:
```markdown
- [ ] Tests grün. **Beleg:** Link zum CI-Run + Image-Hash.
- [ ] Replay gelaufen. **Beleg:** Link zum Replay-Manifest (Modul 12).
```
Die Beleg-Pflicht ist der einzige Schutz gegen Bürokratie.

**Schritt 2 — Pflicht-Items aus Phasen ableiten.**
Eines pro Phase des Kurses:

```markdown
## Freigabe-Checkliste — Welle 1

### Spec / Architektur (Phase 01)
- [ ] Alle abgeschlossenen Slices haben `lastenheft_refs`.
      **Beleg:** Frontmatter-Grep über `done/`.
- [ ] Alle Accepted-ADRs sind referenziert oder superseded.
      **Beleg:** `make adr-graph`.

### Planung (Phase 02)
- [ ] Carveouts sind alle entweder permanent gekennzeichnet
      oder haben Folge-Slice + Trigger.
      **Beleg:** `make carveout-audit`.

### Agenten (Phase 03)
- [ ] AGENTS.md beschreibt nur existierende Konventionen.
      **Beleg:** Doku-Konsistenz-Agent-Lauf (Modul 15).

### Qualität (Phase 04)
- [ ] `make gates` grün auf frischem Klon + im CI mit
      identischem Image-Hash.
      **Beleg:** zwei Run-Links (Klon-Run + CI-Run).
- [ ] Replay-Manifest (Modul 12) mit ≥3 Fällen, alle grün.
      **Beleg:** Link zum manifest.yaml + Run-Output.

### Betrieb (Phase 05)
- [ ] Runbook für *mindestens* den wahrscheinlichsten
      Incident-Typ existiert mit Entscheidungs-Triggern
      (nicht "Service neu starten").
      **Beleg:** Pfad zur Runbook-Datei.
- [ ] Trace-Fixture pro Welle archiviert.
      **Beleg:** OTel-Endpoint oder Pfad zur JSONL-Datei.
```

**Schritt 3 — Anti-Items hinzufügen (was *nicht* gefragt wird).**
Eine Liste der bewusst weggelassenen Häkchen — sonst wandern sie
schleichend in die Pflicht:

```markdown
### Bewusst NICHT in dieser Freigabe
- Manuelle Smoke-Tests in Produktion (delegiert an Validator).
- Aktualität von Stakeholder-Slides (delegiert an Produktmanagement).
- 100 %-Coverage (siehe ADR-0019 zu Critical Coverage).
```

**Schritt 4 — Incident-Klausel verlinken.**
```markdown
### Incident-Bereitschaft
- [ ] Bereitschafts-Dokument zeigt, wer in den ersten 15 Min
      welche der drei Optionen wählt: Rollback · Fix-Forward
      · Datenkorrektur.
      **Beleg:** Link zum Bereitschafts-Dokument.
```
Die Drei-Optionen-Tabelle gehört *vor* den Incident geschrieben — nicht
im Stress entschieden.

**Schritt 5 — Item-für-Item belegen.**
Jetzt durchgehen und *jeden* Beleg-Slot tatsächlich füllen. Wenn ein
Beleg fehlt, ist das Item *nicht* abgehakt — auch wenn das Item
inhaltlich erfüllt wäre. Eine Checkliste ohne Belege ist die
Bürokratie-Form, gegen die der Kurs sich wendet.

**Schritt 6 — Freigabe-Eintrag in `done/welle-1-closure.md`.**
```markdown
# Welle 1 — Closure-Eintrag
Status: released
Datum: 2026-06-30
Checkliste: docs/release/welle-1-checkliste.md (alle Items mit Beleg)
Restrisiken: zwei (siehe §"Bewusst NICHT", plus Folge-Slice SL-027 für
             Coverage-Erhöhung auf Critical-Pfad).
Steering-Loop-Eintrag (Modul 15 Doku-Konsistenz-Agent meldete vor
Freigabe einen Drift in AGENTS.md — wurde behoben, vor Freigabe geprüft).
```

Sechs Schritte, eine Freigabe mit Belegen pro Item. Vergleich:
[`../../lab/example/runbooks/`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/runbooks/) und
[`../../lab/example/Makefile`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/Makefile) Target
`make release`.

### Rollback-vs-Fix-Forward-Regeln

- **Drei Antwortoptionen bei produktivem Incident:** Rollback · Fix-Forward · Datenkorrektur. Drei *verschiedene* Antwortklassen, mit jeweils anderen Voraussetzungen (Rückwärtskompatibilität, Test-Coverage des Fix, Vorhandensein des Originaldatensatzes). Welche der drei greift, ist *vor* dem Incident im Runbook festzulegen — mit Triggern wie "DB-Migration rückwärtskompatibel?" und "Buggy-Daten bereits ausgeliefert?". Wer im Incident wählt, wählt typischerweise unter Stress die teuerste Option.
- **Drei Anti-Rollback-Szenarien:** nicht-rückwärtskompatible DB-Migration, bereits erzeugte Buggy-Daten, ungetesteter Rollback-Pfad. Folge: Rollback gehört *vor* den Incident im Runbook entschieden — als bedingte Regel mit Trigger, nicht als Universal-Reflex. Wer im Incident entscheidet, entscheidet schlecht.
- **Runbook-Form:** die Fälle als *bedingte Regeln* in einer Runbook-Tabelle ("**wenn** Migration nicht rückwärtskompatibel → **dann** kein Rollback").

### Injection-Symptome und Telemetrie-Zuordnung

Telemetrie für nachträgliche Injection-Erkennung — drei Spuren:
Eingabe-Roh-Logging (mit Redaction), Tool-Call-Audit-Log,
Output-vs-Eingabe-Konsistenz-Marker. Ergänzende Indikatoren:
Cache-Miss-Spike, Tool-Allowlist-Reject-Counter — ohne mindestens
*eines* der drei Pflicht-Felder bleibt Erkennung Glücksache.

---

