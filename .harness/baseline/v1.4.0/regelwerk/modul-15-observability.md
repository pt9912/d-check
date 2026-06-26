## Modul 15 — Observability

*Quelle: [05-betrieb/modul-15-observability.md](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/05-betrieb/modul-15-observability.md)*

### Harness-Einordnung

Observability ist Eingangs- und Ausgangskanal für *Entropy Management*
(siehe [`grundlagen/klassifikation.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/klassifikation.md)):
ohne Telemetrie weißt du nicht, wo der Harness rostet.

### Kernidee (Modul 15)

Ein Agenten-Lauf ohne Trace ist ein Vorgang ohne Beleg. Du weißt, dass
es passiert ist; du weißt nicht, *was* passiert ist.

### Regeln gegen typische Fehlannahmen (Modul 15)

- **"Logs reichen."** — Logs sagen *was passierte*, nicht *wer wen wann rief*. Trace ist die Antwort darauf.
- **"Metriken sind nur für Performance."** — Metriken sind auch für *Kosten* (Token, Cache-Hit-Rate) und *Drift* (AGENTS.md-Konsistenz-Score).
- **"Prompt-Caching ist Modell-Sache."** — Nein. Cache-Hits zeigen sich erst in Metriken, wenn du sie misst. Wer Cache-Miss-Spikes nicht beobachtet, sieht Injection-Versuche und Drift-Symptome nicht.
- **"Trace teurer Tool-Call = unnötiger Tool-Call."** — Falsch. Manche teuren Calls sind nötig. Frage: lässt er sich durch Caching, Vorab-Filter oder Kontext-Verdichtung billiger machen?

### Worked Example: ein Span zurück bis zur Lastenheft-ID

**Ausgangs-Span:** Du öffnest den Trace zu `sl-009-agent-run`. Der
teuerste Span trägt:

```json
{
  "span_id": "impl-2",
  "name": "tool_call:writer.write_index",
  "duration_ms": 412,
  "tool.name": "writer.write_index",
  "tool.arguments.redacted": {"docs": 100, "target": "internal/index/store.bin"},
  "tool.result.status": "ok",
  "requirement.id": "LH-FA-IDX-003",
  "adr.id": "ADR-0012",
  "tokens": {"input": 2480, "output": 187},
  "cache": {"hit": false}
}
```

**Schritt 1 — Slice-ID aus dem Trace lesen.**
Der Trace-Header trägt `slice.id = slice-009` (Lab-Schreibweise mit
Bindestrich) und `requirement.refs = ["LH-QA-02", "LH-FA-IDX-003"]`.
Innerhalb des Spans selbst hängt der teure `writer.write_index`-Call
zusätzlich an `requirement.id = LH-FA-IDX-003` — der direkte Anker
zur konkreten Anforderung.

**Schritt 2 — Slice-Datei finden.**
[`docs/plan/planning/done/slice-009-tie-break-determinismus.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/docs/plan/planning/done/slice-009-tie-break-determinismus.md).
Die Lab-Datei trägt keine YAML-Frontmatter, sondern eine Klartext-
Bezug-Zeile:
```markdown
**Bezug:** LH-QA-02 (Reproduzierbarkeit, primär), LH-FA-IDX-003
(Index-Schreib-Idempotenz, sekundär — deterministischer Tie-Break ist
Voraussetzung für bit-identische Schreib-Ergebnisse), ADR-0003
(Index-Format), ADR-0012 (Index-Write-Strategie, sekundär).
```
Das ist eine *alternative Operationalisierung* zum Frontmatter-Schema
(siehe Notiz unten). Lesepfad bleibt derselbe: die Zeile führt zu den
ADRs und Anforderungen.

**Schritt 3 — ADR aufrufen.**
[`docs/plan/adr/0012-index-write-strategy.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/docs/plan/adr/0012-index-write-strategy.md).
Kopf:
```markdown
**Status:** Accepted
**Bezug:** LH-FA-IDX-003 (Index-Schreib-Idempotenz und Atomarität),
ADR-0003 (Index-Storage-Format)
```
Bestätigt: die ADR begründet die Index-Write-Strategie (Temp-File +
Atomic Rename) und referenziert dieselbe LH-ID, die der Span trägt.
Die Kette schließt sich.

**Schritt 4 — Lastenheft prüfen.**
[`spec/lastenheft.md` § `LH-FA-IDX-003`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/spec/lastenheft.md):
```markdown
### LH-FA-IDX-003 — Index-Schreib-Idempotenz und Atomarität
Anforderung: Index-Schreiboperationen sind idempotent (gleicher
Datei-Hash bei gleicher Eingabe) und atomar (kein partieller
Index-Stand beobachtbar).
Akzeptanzkriterien: Happy / Boundary (Crash-Recovery via fsync+rename)
/ Negative (E099 bei nicht beschreibbarem Verzeichnis).
```
Bestätigt: der teuerste Tool-Call (im Span) bedient eine konkrete
Lastenheft-Anforderung mit Akzeptanzkriterien.

**Schritt 5 — Make-Target-Kommentar gegenprüfen.**
ADR-0012 §Fitness Function definiert die maschinelle Prüfung:
Architekturtest pro Sprache erzwingt die `rename`-Sequenz im
Writer-Code; Property-Test (slice-013) vergleicht zwei aufeinander
folgende `writer.write_index`-Hashes. Damit ist die Kette **auch
maschinell prüfbar**: ein Commit, der den `rename`-Aufruf entfernt,
würde `make arch-check` rot machen
([`konventionen.md` §Traceability-Constraint](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/konventionen.md#traceability-constraint)).

**Schritt 6 — Bruchpunkt benennen.**
Vollständige Kette:
```
trace.slice.id            →  slice-009
span.requirement.id       →  LH-FA-IDX-003
                          →  done/slice-009-tie-break-determinismus.md (Bezug-Zeile)
                          →  ADR-0012 (Bezug: LH-FA-IDX-003)
                          →  LH-FA-IDX-003 (Akzeptanzkriterien)
                          ↩  ADR-0012 §Fitness Function prüft Architekturregel
```
Schwächste Stelle in *diesem* Beispiel-Repo: der Lesepfad zwischen
Slice-Datei und ADRs läuft über eine *Klartext*-Zeile, nicht über ein
maschinell parsbares Frontmatter. Wenn die Bezug-Zeile umformuliert
wird, fehlt der direkte Anker. Steering-Loop-Aktion: entweder
Frontmatter-Pflichtfeld als Gate ergänzen (computational feedforward),
oder einen Doku-Konsistenz-Agenten die Bezug-Zeile prüfen lassen
(inferential feedback).

> **Operationalisierungs-Variante.** Ein anderes Repo kann denselben
> Lesepfad über YAML-Frontmatter abbilden:
> ```yaml
> ---
> id: SL-009
> adr_refs: [ADR-0012]
> lastenheft_refs: [LH-FA-IDX-003]
> ---
> ```
> Vorteil: maschinell trivial parsbar. Nachteil: zusätzliche Disziplin
> im Slice-Template. Das Lab wählt die Klartext-Variante, weil die
> bestehenden Slices ohne Frontmatter angelegt waren — Migration wäre
> teurer als der Komfort des Schemas. Die *Erschaffens-Leistung* dieses
> Moduls ist das Span-Schema mit IDs; *welche* Schreibvariante die
> Slice-Seite wählt, ist Repo-spezifisch.

Sechs Schritte, eine durchgängige Traceability. Vergleich im Lab:
[`../../lab/example/otel/sl-009-agent-run.trace.json`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/otel/sl-009-agent-run.trace.json)
(Span `impl-2` ist der `writer.write_index`-Call mit `requirement.id`
und `adr.id`).

### Span-/Audit-Attribut-Regeln

- **Drei Telemetrie-Typen und ihre Fragen:** Logs (*was passierte*) · Metriken (*wie oft, wie schnell, wie viel*) · Traces (*wer rief wen, in welcher Reihenfolge*). Drei verschiedene Fragen, drei verschiedene Werkzeuge. Operative Folge: Wer nur Logs hat, kann Cost-Attribution nicht durchführen (braucht Metriken) und Tool-Call-Ketten nicht rekonstruieren (braucht Traces). Ein Agent-System mit nur einem Typ ist forensisch nicht antwortfähig.
- **Mindestfelder eines Tool-Call-Spans:** `tool.name`, `tool.arguments` (redacted), `tool.result.status` plus Korrelations-IDs zu Slice/PR/Agent-Rolle. Begründung: Ohne `slice.id` / `requirement.id` ist Token-Attribuierung pro Slice nicht möglich; ohne `agent.role` bricht die Rollen-Trennung in der Forensik.
- **Audit-Span-Schema:** liste jeden Attribut-Namen, markiere ihn als *Pflicht* oder *Optional* und nenne pro Attribut die *Incident-Frage*, die es beantwortet (z. B. `slice.id` → "auf wessen Rechnung lief der Schreibzugriff?"; `tool.arguments.redacted` → "was wurde wohin geschrieben — ohne Secrets im Log?"). Pflicht-Minimum aus dem Worked Example: Slice-ID, Agent-Rolle, Cache-Status, `requirement.id` — jede Abweichung davon begründest du. Ein Attribut ohne Incident-Frage fliegt raus: Schema-Felder ohne Abnehmer sind Telemetrie-Boilerplate, kein Audit.

### Token-Attributions-Regeln

Summiere Input- und Output-Token pro `agent.role` (Planner · Architect ·
Implementer · Reviewer · Verifier) und gib an, welche Rolle den größten
Anteil trägt — als Zahl *und* als Prozentsatz der Gesamtsumme. Wo ein
Span keinen Rollen-Tag trägt (Sammelposten), entscheide begründet, wie
du ihn aufteilst (anteilig nach Tool-Calls? dem auslösenden Slice
zugeschlagen?) — genau das ist das Buchhaltungs-Splitting eines
Sammelpostens auf Kostenstellen.

### Cache-Counter-Regeln

Die *drei* OTel-Counter, die du brauchst, um Cache-Hit-Rate *und*
Cache-Miss-Spikes zu unterscheiden — pro Counter:

| Frage | Antwort |
|---|---|
| Name | z. B. `prompt_cache_hits_total` |
| Unit | Cardinality (Counter, Gauge, Histogram?) |
| Labels | mindestens `slice.id`, `agent.role`, `model.version` |
| Aggregation | Hit-Rate als `hits / (hits + misses)` — wo wird die Division ausgeführt: in der Metrik-DB oder im Dashboard? |

Eine *einzelne* Metrik `cache.hit_ratio` reicht nicht: ohne separate
Counter für Hits *und* Misses kannst du Cache-Miss-Spikes
(Sicherheits-Indikator!) nicht von Cache-Hit-Rückgängen
(Kosten-Indikator) trennen.

**Cache-Miss in den Metriken erkennen:** Anstieg der
Token-Eingabe-Metrik *ohne* Anstieg der Cache-Hit-Rate-Metrik
(`cache.hit_ratio` fällt). Zweck: Cache-Miss-Spikes sind oft
Injection-Symptome (variable Eingaben umgehen Cache absichtlich) —
Metrik dient also gleichzeitig Kosten- *und* Sicherheitsüberwachung.

### Doku-Konsistenz-Drift-Regeln

Konsistenz-Regeln, die ein Doku-Konsistenz-Agent zwischen AGENTS.md und
realen Make-Targets / Skill-Dateien / `harness/README.md` prüft — pro
Regel:

| Feld | Inhalt |
|---|---|
| **Regel-Name** | z. B. *"AGENTS.md-Befehl existiert im Makefile"* |
| **Quelle** | welche Datei wird gelesen (z. B. `AGENTS.md` §Tool-Regeln) |
| **Vergleichs-Ziel** | welche Datei wird dagegen geprüft (z. B. `Makefile`-Target-Namen) |
| **Drift-Symptom** | wie sieht ein Drift-Treffer konkret aus (z. B. *"AGENTS.md nennt `make fullbuild`, Makefile kennt nur `make build`"*) |
| **Lebenszyklus** | ist das ein Pre-commit-Check, Pre-integration, oder Continuous (vgl. [`grundlagen/klassifikation.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/klassifikation.md))? |

Mindestens *eine* Regel muss die Hard Rule aus
[Modul 13 §"Hard Rule (Doku-Disziplin)"](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/04-qualitaet/modul-13-quality-gates.md#hard-rule-doku-disziplin)
durchsetzen ("keine Befehle behaupten, die es nicht gibt").

**Drift-Signal und Schwelle:** Konkretes Signal: Doku-Konsistenz-Agent
meldet AGENTS.md-Befehl ohne passendes Make-Target (z. B.
`make fullbuild` behauptet, Makefile kennt nur `make build`);
Konsistenz-Score als Metrik (`agents_md.consistency_ratio`) fällt unter
einen Schwellwert. Schwelle begründet: jeder behauptete-aber-fehlende
Befehl ist *sofort* gate-relevant (Hard Rule Modul 13, keine Befehle
erfinden), nicht erst ab einem Prozentsatz — Score-Verfall ist nur das
Aggregat-Signal. Gegenbeispiel-Rauschen: ein neu hinzugefügtes Target
ohne AGENTS.md-Eintrag ist *Vorwärts*-Drift (Doku hinkt nach), andere
Härte als behauptete Geister-Befehle.

