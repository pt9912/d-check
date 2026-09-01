## Modul 15 — Observability

<!-- Quelle: [05-betrieb/modul-15-observability.md](https://github.com/pt9912/ai-harness-course/blob/v5.18.0/kurs/de/05-betrieb/modul-15-observability.md) -->

### Harness-Einordnung

Observability ist Eingangs- und Ausgangskanal für *Entropy Management*
(siehe [`grundlagen/klassifikation.md`](grundlagen-klassifikation.md)):
ohne Telemetrie weißt du nicht, wo der Harness rostet.

### Kernidee (Modul 15)

Ein Agenten-Lauf ohne Trace ist ein Vorgang ohne Beleg. Du weißt, dass
es passiert ist; du weißt nicht, *was* passiert ist.

Konkret beobachtest du einen Agentenlauf als **Trace aus Spans** — einen
pro Tool-Call. Der teuerste Span trägt Korrelations-IDs (`slice.id`,
`requirement.id`, `adr.id`, `agent.role`); über sie verfolgst du ihn bis
zur Anforderung zurück und prüfst, ob die Kette **maschinell** hält.
Dafür braucht jeder Span ein Audit-Schema, jede Rolle eine Token-Bilanz,
jeder Cache einen Counter — die Regeln unten.

### Regeln gegen typische Fehlannahmen (Modul 15)

- **"Logs reichen."** — Logs sagen *was passierte*, nicht *wer wen wann rief*. Trace ist die Antwort darauf.
- **"Metriken sind nur für Performance."** — Metriken sind auch für *Kosten* (Token, Cache-Hit-Rate) und *Drift* (AGENTS.md-Konsistenz-Score).
- **"Prompt-Caching ist Modell-Sache."** — Nein. Cache-Hits zeigen sich erst in Metriken, wenn du sie misst. Wer Cache-Miss-Spikes nicht beobachtet, sieht Injection-Versuche und Drift-Symptome nicht.
- **"Trace teurer Tool-Call = unnötiger Tool-Call."** — Falsch. Manche teuren Calls sind nötig. Frage: lässt er sich durch Caching, Vorab-Filter oder Kontext-Verdichtung billiger machen?

### Span-/Audit-Attribut-Regeln

- **Drei Telemetrie-Typen und ihre Fragen:** Logs (*was passierte*) · Metriken (*wie oft, wie schnell, wie viel*) · Traces (*wer rief wen, in welcher Reihenfolge*). Drei verschiedene Fragen, drei verschiedene Werkzeuge. Operative Folge: Wer nur Logs hat, kann Cost-Attribution nicht durchführen (braucht Metriken) und Tool-Call-Ketten nicht rekonstruieren (braucht Traces). Ein Agent-System mit nur einem Typ ist forensisch nicht antwortfähig.
- **Mindestfelder eines Tool-Call-Spans:** `tool.name`, `tool.arguments` (redacted), `tool.result.status` plus Korrelations-IDs zu Slice/PR/Agent-Rolle. Begründung: Ohne `slice.id` / `requirement.id` ist Token-Attribuierung pro Slice nicht möglich; ohne `agent.role` bricht die Rollen-Trennung in der Forensik.
- **Audit-Span-Schema:** liste jeden Attribut-Namen, markiere ihn als *Pflicht* oder *Optional* und nenne pro Attribut die *Incident-Frage*, die es beantwortet (z. B. `slice.id` → "auf wessen Rechnung lief der Schreibzugriff?"; `tool.arguments.redacted` → "was wurde wohin geschrieben — ohne Secrets im Log?"). Pflicht-Minimum: Slice-ID, Agent-Rolle, Cache-Status, `requirement.id` — jede Abweichung davon begründest du. Ein Attribut ohne Incident-Frage fliegt raus: Schema-Felder ohne Abnehmer sind Telemetrie-Boilerplate, kein Audit.

- **Der Emissions-Pfad ist Repo-Entscheidung** (Exporter, Collector, Sampling,
  Aufbewahrung): Mitzunehmen ist das **Schema**, nicht das Setup.

### Token-Attributions-Regeln

Summiere Input- und Output-Token pro `agent.role` — attribuiert wird damit auf
**Kontexte**, nicht auf Personen; die Rollen sind die aus
[Modul 8](modul-08-agentenrollen.md), festgelegt durch das gestartete
Rollen-Artefakt — und gib an, welche Rolle den größten
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
| **Lebenszyklus** | ist das ein Pre-commit-Check, Pre-integration, oder Continuous (vgl. [`grundlagen/klassifikation.md`](grundlagen-klassifikation.md))? |

Mindestens *eine* Regel muss die Hard Rule aus
[Modul 13 §"Hard Rule (Doku-Disziplin)"](modul-13-quality-gates.md#hard-rule-doku-disziplin)
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

