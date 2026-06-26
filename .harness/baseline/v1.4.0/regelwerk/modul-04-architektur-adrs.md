## Modul 4 — Architektur und ADRs

*Quelle: [01-spec-und-architektur/modul-04-architektur-adrs.md](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/01-spec-und-architektur/modul-04-architektur-adrs.md)*

### Mini-Glossar für dieses Modul (Modul 4)

| Begriff | Ein-Satz-Definition | Bild im Kopf |
|---|---|---|
| **MADR** | Markdown-basiertes ADR-Format mit Kopf-Feldern (Status, Datum, Bezug, Supersedes) und Body-Blöcken (Kontext, Optionen mit Trade-offs, Entscheidung, Konsequenzen). | ein Formular, das die Entscheidung zwingt, ihre Belege mitzubringen. |
| **Nygard-Format** | Das ursprüngliche, schlankere ADR-Format nach Michael Nygard: Kontext, Entscheidung, Konsequenzen. | der Urahn von MADR — gleiche Idee, weniger Felder. |
| **superseded** | ADR-Status: Entscheidung ist durch eine *neue* ADR abgelöst — der Bedarf bleibt, die Antwort wechselt. | Schild "ersetzt durch Nr. N" am alten Protokoll. |
| **deprecated** | ADR-Status: Entscheidung entfällt *ersatzlos* — der zugrunde liegende Bedarf existiert nicht mehr. | Akte geschlossen, kein Nachfolger nötig. |
| **Fitness-Function-Werkzeuge** | ArchUnit (Java), dep-cruiser (JS/TS), import-linter (Python) — prüfen Architektur-Aussagen maschinell, z. B. Layer-Importregeln. | der Prüfstand, auf den die ADR-Aussage geschnallt wird. |

### Harness-Einordnung (Modul 4)

ADR = *inferential feedforward* (für den Implementation-Agent) und
gleichzeitig Quelle für *computational feedback* (ArchUnit/Fitness
Functions, wenn die Entscheidung maschinell prüfbar ist). Eine ADR ohne
Fitness Function ist eine Absichtserklärung.

### Kernidee (Modul 4)

Ein ADR ist die einzige Stelle, an der "weil" gegen "ist halt so" gewinnt.
Wenn dein Reviewer-Agent den Grund nicht findet, kann er die Entscheidung
nicht verteidigen.

### Hard Rule (Beispiel aus c-hsm-doc, ADR 0001)

Begriff *Hard Rule* siehe Glossar in
[`grundlagen/konventionen.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/konventionen.md).

*"Eine ADR mit Status `Accepted` wird nicht inhaltlich überschrieben.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
explizitem Verweis auf die abgelöste oder geschärfte Vorgängerin."*

Wirkung: ADRs sind Geschichtsdokumente, kein Wiki. Reviewer-Agent kann
auf ältere Entscheidungen vertrauen, ohne Versionsstände zu vergleichen.

### Regeln gegen typische Fehlannahmen (Modul 4)

- Nein. ADRs begründen die *Lösung*. Anforderungen begründet die Spec. Wer ADRs zur Spec macht, kann später keine Architektur ohne Lastenheft-Änderung wechseln.
- Hard Rule: Accepted-ADRs werden nicht überschrieben. Folge-ADR mit `supersedes ADR-N`. Sonst kann der Reviewer-Agent nicht auf ältere Entscheidungen vertrauen.
- Eine ADR ohne Fitness Function ist eine Absichtserklärung. Wer architecture fitness im Kopf hat, schreibt parallel den ArchUnit-Test.
- MADR ist ein Format unter mehreren (auch Nygard, Tyree/Akerman). Wichtig ist, dass dein Repo *eines* konsequent benutzt.
- Diagramme sind *eine* Output-Form, nicht die Sache selbst. Architektur in diesem Kurs heißt: *Entscheidungen mit Begründung (ADR), prüfbar gemacht (Fitness Function), versioniert (Accepted-Hard-Rule)*. Ein Diagramm ohne ADRs hinter sich ist Wandtapete; eine ADR ohne Fitness Function ist Absichtserklärung. `spec/architecture.md` ist explizit *diagrammatisch und enthält keine eigenen Anforderungen* (siehe Spec-Stratifizierung in [`grundlagen/konventionen.md#spec-stratifizierung`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/konventionen.md#spec-stratifizierung)) — genau weil sonst Bilder anfangen würden, die ADR-Schicht zu ersetzen.
- Eine ADR ohne maschinelle Durchsetzung ist eine *Absichtserklärung*, die der Implementation-Agent freundlich liest und dann ignoriert, wenn ein anderer Pfad "einfacher" wirkt. Eine ADR *mit* Fitness Function ist ein Constraint — die Layering-Regel, die ArchUnit dem Agenten als roten Build entgegenhält. Worked Example in [Modul 13 §Worked Example "ADR → import-linter"](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/04-qualitaet/modul-13-quality-gates.md#worked-example-vom-adr-satz-zur-fitness-function) zeigt, was die Übersetzung kostet (kleine Tabelle: ADR-Satz, Werkzeug, Make-Target, Failure-Beispiel). Wer das nicht macht, dokumentiert *Hoffnung*.

### Worked Example: vom Diskussionsfaden zum prüfbaren ADR

**Schritt 1 — Triggerschwelle erreichen.** Drei Vorfälle = Symptom →
Lücke im Harness. Architect-Agent legt ADR-Entwurf an: `0007-service-adapter-layer.md`.

**Schritt 2 — MADR-Kopf:**
```markdown
# ADR-0007 — Service-Schicht spricht externe APIs nur über Adapter

* Status: Accepted
* Datum: 2026-06-15
* Bezug: LH-QA-COUPLING-002
* Supersedes: —
```

**Schritt 3 — Kontext (Spec-Verweis statt -Wiederholung):**
> Wiederholter Wunsch, im `service/`-Layer direkt `http.Client` zu
> instanziieren. LH-QA-COUPLING-002 verlangt, dass externe Abhängigkeiten
> austauschbar bleiben (für Replay und für Provider-Wechsel).

**Schritt 4 — Optionen mit Trade-offs:**
> 1. **Direkt-Calls in Service-Schicht** — minimal Boilerplate; bricht LH-QA-COUPLING-002 (kein Replay ohne API-Mocks).
> 2. **Adapter-Schicht mit Interface** — etwas Boilerplate; erfüllt LH-QA-COUPLING-002; Replay-fähig.
> 3. **Service-Mesh / Sidecar** — verschiebt das Problem in Infrastruktur; überdimensioniert für aktuelle Repo-Größe.

**Schritt 5 — Entscheidung:**
> Option 2. Service-Layer importiert ausschließlich aus `adapter/`-Paket.
> HTTP-Client lebt unter `adapter/http/`.

**Schritt 6 — Konsequenz mit Fitness Function:**
> ArchUnit-Test `arch_no_direct_http_in_service`:
> Keine Klasse in `service.*` darf `java.net.http.*` oder `okhttp.*`
> importieren.
> Gate: `make arch-check` (vergleichbar mit dep-cruiser/dep-rule für
> Python/Go).

**Schritt 7 — Lerneintrag in `done/`** (nach Schließen der Welle, in
der ADR-7 implementiert wird):
> Steering-Loop-Beleg: drei Vorfälle in zwei Wochen → ADR-7 →
> ArchUnit-Test → kein weiterer Vorfall in 6 Wochen.

Sieben Schritte, eine geprüfte Entscheidung. Vergleich:
[`/lab/example/docs/plan/adr/`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/docs/plan/adr/).
