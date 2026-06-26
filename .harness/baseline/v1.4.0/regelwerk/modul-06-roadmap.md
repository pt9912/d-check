## Modul 6 — Roadmap Engineering

*Quelle: [02-planung/modul-06-roadmap.md](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/02-planung/modul-06-roadmap.md)*

### Kernidee (Modul 6)

Eine Roadmap ist eine Reihenfolge von Wellen, keine Reihenfolge von
Terminen. Termine sind eine Folge der Wellen, nicht ihr Treiber.

### Roadmap-Regeln (Modul 6)

- Ein Welle-Eintrag braucht minimal drei Bestandteile: Slice-IDs (Inhalt) · Trigger als beobachtbare Bedingung (kein Datum) · Closure-Kriterien (z. B. Replay grün, alle Slices in `done/`). Datum darf *erwähnt* werden (Prognose), darf aber nie Trigger sein — sonst kappt die Welle halbfertige Slices am Kalendertag und das Auditierbarkeits-Versprechen bricht.
- Ein Trigger ist beobachtbar dann, wenn ein *anderer* Mensch ohne Rückfrage sagen kann, ob er eingetreten ist. "Sobald wir Zeit haben" scheitert daran; "SL-024 in `done/`" besteht. Beispiele für beobachtbare Trigger: "SL-024 liegt in `done/`" · "Replay-Lauf gegen Golden Set grün" · "Carveout `CO-007` aufgelöst".
- Welle 30 % über Schätzung — Diagnose vor Aktion: liegt es an Slice-Größe (→ neu schneiden), an Reihenfolge (→ neu planen), oder an unerwarteter Komplexität (→ Carveout)? 30 % früh können ein Steering-Loop-Signal sein (Slice-Sizing-Regel schärfen), 30 % spät (vor Welle-Closure) eher Carveout.

### Welle ≠ Meilenstein ≠ Release (Modul 6)

- **Welle** = Bündel paralleler/serialisierter Slices mit Closure-Kriterien. Eine Welle endet *durch* Closure-Kriterien.
- **Meilenstein** = extern beobachtbarer Zustand (Release, Audit-Punkt). Ein Meilenstein endet durch *Datum oder externe Bestätigung* — und genau deshalb leitet sich der Meilenstein aus Wellen ab, nicht umgekehrt.
- **Release** — Trigger: ein Artefakt verlässt das Repo in eine Umgebung (Tag + Staging). Ein Release kann mehrere Wellen umfassen, der Meilenstein liegt *neben* der Welle (externe Bestätigung), die Welle endet *durch* Closure.

### Worked Example: einen Datumswunsch in eine Trigger-Welle übersetzen

**Schritt 1 — Wunsch in Inhalt zerlegen.** Frage zurück: *Was muss
*beim Audit* gezeigt werden?* Die Antwort ist immer eine Liste von
beobachtbaren Zuständen — und genau diese werden zu Closure-Triggern.
Stakeholder antwortet konkret: "ANN-Suche funktioniert auf 100k
Einträgen unter 1 s p95; Multi-Sprach-Adapter ist konsolidiert; OTel-
Pipeline zeigt End-to-End-Traces."

Drei Zustände, drei Trigger-Anker — und keiner davon enthält ein
Datum.

**Schritt 2 — Inhalt in Slices binden.** Jeder Closure-Trigger muss auf
einen oder mehrere Slices mit eigenem DoD verweisen. Sonst ist der
Trigger ein Wunsch, kein Beleg.

| Trigger-Anker (Stakeholder) | Slice(s) (Implementer-Ebene) |
|---|---|
| ANN-Suche < 1 s p95 bei 100k | `slice-014` (ANN-Bibliothek-Integration) + `slice-019` (Latenz-Replay gegen 100k-Korpus) |
| Multi-Sprach-Adapter konsolidiert | `slice-015` (Adapter-Cleanup) |
| OTel-Pipeline E2E | `slice-017` (OTel-Collector) + `slice-018` (Trace-Schema-Pflicht) |

Mehrfachbezüge sind erlaubt — *fehlende* Bezüge nicht. Wer einen
Trigger ohne Slice formuliert, hat einen Wunsch ohne Plan.

**Schritt 3 — Abhängigkeiten gegen vorhandene Wellen messen.** Eine
Welle, die ohne fertige Vorgängerin nicht starten kann, ist eine
Phantom-Welle. Lab-Beispiel: Welle 3 (`welle-3-skalierung`) hängt an
Welle 2 (`welle-2-qualitaet`) — Property-Tests müssen *vor* der
Skalierungs-Welle stehen, weil sonst die Skalierungs-Gates auf einer
nicht-property-getesteten Basis laufen.

Im Abhängigkeitsgraphen wird das eine gerichtete Kante; in der
Roadmap-Tabelle ein expliziter Eintrag in der `Trigger`-Spalte.

**Schritt 4 — Welle-Eintrag mit den drei Pflicht-Bestandteilen
schreiben.** Closure-Kriterien · Slice-IDs · Abhängigkeits-Trigger.
Vorbild aus dem Lab
([`roadmap.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/docs/plan/planning/in-progress/roadmap.md)):

```markdown
## Aktuelle Welle

**Welle-ID:** welle-3-skalierung
**Geplantes Ende:** 2026-07-24 (Schätzung)

**Closure-Trigger:**
- slice-014 (ANN-Bibliothek) done in allen Sprachen.
- slice-015 (Multi-Sprach-Adapter-Cleanup) done.
- slice-019 (Latenz-Replay) grün: p95 < 1 s bei 100k Korpus.
- ADR-0004 (ANN-Bibliothek-Wahl) `Accepted`.

**Vorgänger-Trigger:** welle-2-qualitaet done.
```

Datum *erscheint* als "Geplantes Ende (Schätzung)" — es triggert
nichts, es prognostiziert. Wenn die Schätzung kippt, kippt sie als
Schätzung, nicht als Closure-Kriterium.

**Schritt 5 — Meilenstein neben die Welle setzen, nicht in sie.** Der
Audit-Termin ist *Meilenstein M3*, nicht *Welle 3*. Welle und
Meilenstein verhalten sich orthogonal:

| Welle | Meilenstein |
|---|---|
| endet durch Closure-Kriterien (intern) | endet durch externe Bestätigung (Audit, Release, Kunde) |
| Inhalt vollständig im Repo | Inhalt zeigt sich an einer Außengrenze |
| `welle-3-skalierung` | M3 — Skalierbar |

Tabelle aus dem Lab:

```markdown
| Meilenstein | Welle(n) | Trigger | Status |
|---|---|---|---|
| M3 — Skalierbar | welle-3-skalierung | p95 < 1 s auch bei 100k Einträgen | offen |
```

Der Audit-Termin (`2026-07-31`) ist Anhang im Meilenstein-Eintrag, nicht
Trigger der Welle. Das hat eine harte Konsequenz: wenn das Audit-Datum
gehalten werden *muss*, aber die Closure-Trigger nicht erreichbar sind,
ist die richtige Antwort ein *Carveout* (Modul 7), nicht ein halb
fertiges `done/`.

**Schritt 6 — Drift-Tabelle als Pflicht-Anhang.** Eine Roadmap, die
sich nie korrigiert, hat den Steering Loop nicht durchlaufen.
Pflicht-Block am Ende:

```markdown
## Historische Trigger-Verschiebungen

| Datum | Was wurde geändert? | Warum? |
|---|---|---|
| 2026-06-12 | slice-019 in welle-3 nachgenommen | Stakeholder ergänzte Audit-Anforderung; Trigger wäre sonst nicht beweisbar gewesen |
```

Diese Tabelle ist nicht Hilfsmittel; sie ist das Audit-Signal. Wer sie
leer hat, hat eine starre Roadmap. Wer sie *jeden* Eintrag voll hat,
hat eine treibende Roadmap.

**Schritt 7 — Datum als Trigger geschrieben: die drei möglichen
Ausgänge.** Ein Closure-Trigger als Datum formuliert (*"welle-3-
skalierung schließt am 2026-07-24"*) — am 24. Juli ist
slice-019 noch nicht grün. Drei mögliche Antworten:

| Antwort | Diagnose |
|---|---|
| Welle wird trotzdem geschlossen, slice-019 wandert in welle-4. | Datum hat Closure überschrieben — der Audit fällt durch, weil slice-019 nicht belegt ist. Trigger-Disziplin ist Theorie geblieben. |
| Welle bleibt offen, das Datum wird verschoben. | Trigger-Disziplin wirkt, aber die Roadmap-Drift-Tabelle muss den Eintrag bekommen — sonst ist die Verschiebung still. |
| Carveout `CO-009` für die fehlende Latenz, Welle schließt mit Carveout. | Sauber: das Versprechen wird offen reduziert, Folge-Slice ist verdrahtet, Audit weiß, was er ansieht. |

*Eine Roadmap ist nicht "wann?", sondern "in welcher Reihenfolge
wovon?"*. Vergleich:
[`roadmap.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/docs/plan/planning/in-progress/roadmap.md).

### Regeln gegen typische Fehlannahmen (Modul 6)

- **Gegen "Roadmap ist eine Datumsleiste":** Datum ist Output, nicht Input. Wer Datumsleisten plant, plant Wunschdenken.
- **Gegen "Burndown ist Fortschritt":** Burndown ist *Tempo*. Fortschritt ist, ob die Welle das verspricht, was sie sollte.
- **Gegen "Eine Roadmap ist statisch":** Eine Roadmap, die nach drei Wellen nicht angepasst wurde, hat den Steering Loop nicht durchlaufen.
- **Gegen "Welle = Sprint":** Ein Sprint endet durch *Datum* (zwei Wochen sind um). Eine Welle endet durch *Closure-Kriterien* (alle ihre Slices in `done/`, Replay-Lauf grün, Closure-Einträge geschrieben). Wer Wellen wie Sprints schneidet, kappt halbfertige Slices am Datum — und produziert genau die Auditierbarkeits-Lücke, die der Harness verhindern soll.
- **Gegen "Trigger = Datum":** Ein Trigger ist eine *beobachtbare Bedingung* ("SL-024 liegt in `done/`", "Replay-Lauf gegen Golden Set grün", "Carveout `CO-007` aufgelöst"). Ein Datum ist kein Trigger, sondern eine Prognose. Wenn das einzige Trigger-Kriterium ein Kalendertag ist, plant die Roadmap nicht — sie hofft.

