## Modul 16 — Produktiver Betrieb

<!-- Quelle: [05-betrieb/modul-16-produktiver-betrieb.md](https://github.com/pt9912/ai-harness-course/blob/v5.12.0/kurs/de/05-betrieb/modul-16-produktiver-betrieb.md) -->

### Kernidee (Modul 16)

Produktiv heißt: Du musst eine Frage in der Nacht beantworten können,
ohne den Autor zu kennen. Runbooks und Replay sind dafür da.

### Regeln gegen typische Fehlannahmen (Modul 16)

- **"Rollback ist die Standardantwort."** — Drei Fälle, in denen Rollback schadet: nicht-rückwärtskompatible DB-Migration, bereits erzeugte Buggy-Daten, ungetesteter Rollback-Pfad. Runbook entscheidet *vor* dem Incident, wann Fix-Forward gilt.
- **"Runbook beschreibt den Happy Path."** — Nein. Runbook beschreibt *Entscheidungen unter Unsicherheit*, mit Triggern. Wenn das Runbook nur sagt "Service neu starten", ist es kein Runbook.
- **"Produktionsfreigabe ist eine formale Checkbox."** — Eine Checkliste ohne *Belege* pro Item (Replay-Lauf-Link, ADR-ID, Trace-Hash) ist Bürokratie. Mit Belegen ist sie das einzige nicht-fragmentierte Audit-Artefakt.
- **"Deployt heißt produktiv."** — Nein. Deployment ist *eine* Anwendung des Container-Ankers (Modul 14), nicht das Ziel. Produktionsreife heißt *belegte Betriebsfähigkeit*: kann ein anderer Mensch nachts handeln (Runbook), ist der Lauf reproduzierbar (Replay-Beleg), entfällt die Freigabe bei einem Incident automatisch (Incident-Klausel)? Ein Service kann längst deployt und trotzdem nicht produktionsreif sein — genau diese Lücke schließt die Freigabe-Checkliste.
- **"Prompt-Injection ist eine Modell-Frage."** — Nein. Erkennung von Injection ist eine *Telemetrie-Frage*: Eingabe-Logging + Tool-Call-Audit + Output-Drift-Marker. Wer das nicht hat, erkennt Injection nur durch Glück.
- **"Postmortem ist Schuldzuweisung — also macht man's leise."** — Genau das Gegenteil. Ein produktiver Postmortem ist *blameless* (vgl. Etsy/Google SRE-Tradition): er sucht den Pfad, auf dem ein vernünftiger Mensch unter Druck dieselbe Entscheidung getroffen hätte, und fragt, *welcher Sensor oder Guide gefehlt hat*. Closure-Einträge in `done/` ([Modul 5](modul-05-planning-harness.md)) und Reflexions-Einträge sind beide *strukturell* blameless: sie fragen "welche Harness-Lücke war Ursache", nicht "wer war es". Wer Postmortems als Schuldzuweisung erlebt hat, wird Drift-Symptome zukünftig verschweigen — und genau dadurch wachsen sie. Blameless ist keine moralische Wahl; es ist eine Sensor-Schutz-Maßnahme.

### Produktionsfreigabe-Checkliste (Modul 16)

Produktionsreife heißt *belegte Betriebsfähigkeit*, nicht „deployt"; die
Freigabe-Checkliste ist das Audit-Artefakt dafür. Regeln:

- **Kein Häkchen ohne Beleg-Slot.** Jedes Item trägt einen konkreten
  Beleg (CI-Run-Link + Image-Hash · Replay-Manifest-Link · ADR-ID ·
  Trace-Hash · Frontmatter-Grep). Fehlt der Beleg, ist das Item *nicht*
  abgehakt — auch wenn es inhaltlich erfüllt wäre. Die Beleg-Pflicht ist
  der einzige Schutz gegen Bürokratie.
- **Ein Pflicht-Item pro Phase**, je mit Beleg: Spec/Architektur
  (Slices tragen `lastenheft_refs`; Accepted-ADRs referenziert oder
  superseded) · Planung (Carveouts permanent markiert oder mit
  Folge-Slice + Trigger) · Agenten (`AGENTS.md` beschreibt nur
  existierende Konventionen) · Qualität (`make gates` grün auf frischem
  Klon *und* im CI mit identischem Image-Hash; Replay-Manifest ≥ 3 Fälle
  grün) · Betrieb (Runbook mit Entscheidungs-Triggern statt „Service neu
  starten"; Trace-Fixture pro Welle archiviert).
- **Anti-Items explizit auflisten** (bewusst *nicht* Teil dieser Freigabe
  — z. B. manuelle Smoke-Tests → Validator, 100 %-Coverage → ADR), sonst
  wandern sie schleichend in die Pflicht.
- **Incident-Klausel verlinken:** wer in den ersten 15 Minuten welche der
  drei Optionen wählt (Rollback · Fix-Forward · Datenkorrektur), steht
  *vor* dem Incident fest (siehe §Rollback-vs-Fix-Forward-Regeln unten).
- **Freigabe-Eintrag** in `done/welle-NN-closure.md`: Status · Datum ·
  Checklisten-Pfad (alle Items mit Beleg) · Restrisiken (Zeiger auf
  Anti-Items + Folge-Slices) · Steering-Loop-Eintrag.

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

