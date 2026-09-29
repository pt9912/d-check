## Modul 4 — ADRs

<!-- Quelle: [01-spec-und-architektur/modul-04-adrs.md](https://github.com/pt9912/ai-harness-course/blob/v6.13.0/kurs/de/01-spec-und-architektur/modul-04-adrs.md) -->

### Mini-Glossar für dieses Modul (Modul 4)

| Begriff | Ein-Satz-Definition | Bild im Kopf |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------------------------------------------------------------------- |
| **MADR** | Markdown-basiertes ADR-Format mit Kopf-Feldern (Status, Datum, Bezug, Supersedes) und Body-Blöcken (Kontext, Optionen mit Trade-offs, Entscheidung, Konsequenzen). | ein Formular, das die Entscheidung zwingt, ihre Belege mitzubringen. |
| **Nygard-Format** | Das ursprüngliche, schlankere ADR-Format nach Michael Nygard: Kontext, Entscheidung, Konsequenzen. | der Urahn von MADR — gleiche Idee, weniger Felder. |
| **superseded** | ADR-Status: Entscheidung ist durch eine *neue* ADR abgelöst — der Bedarf bleibt, die Antwort wechselt. | Schild "ersetzt durch Nr. N" am alten Protokoll. |
| **deprecated** | ADR-Status: Entscheidung entfällt *ersatzlos* — der zugrunde liegende Bedarf existiert nicht mehr. | Akte geschlossen, kein Nachfolger nötig. |
| **Fitness-Function-Werkzeuge** | [a-check](https://github.com/pt9912/a-check), ArchUnit (Java), dep-cruiser (JS/TS), import-linter (Python) — prüfen Architektur-Aussagen maschinell, z. B. Layer-Importregeln. | der Prüfstand, auf den die ADR-Aussage geschnallt wird. |

### Harness-Einordnung (Modul 4)

ADR = *inferential feedforward* (für den Implementer-Agent) und
gleichzeitig Quelle für *computational feedback* (ArchUnit/Fitness
Functions, wenn die Entscheidung maschinell prüfbar ist). Eine ADR ohne
Fitness Function ist eine Absichtserklärung.

### Kernidee (Modul 4)

Ein ADR ist die einzige Stelle, an der "weil" gegen "ist halt so" gewinnt.
Wenn dein Reviewer-Agent den Grund nicht findet, kann er die Entscheidung
nicht verteidigen.

- **Jede ADR trägt einen Re-Evaluierungs-Trigger** — beobachtbare Bedingung,
  unter der die Entscheidung erneut geprüft wird, oder ausdrücklich
  *permanent*. Eine ADR ohne Trigger gilt unbefristet weiter, auch wenn ihre
  Voraussetzung weg ist. Geprüft im **Trigger-Audit** der Welle-Closure
  ([Modul 6](modul-06-roadmap.md)); bei Eintreten: bestätigen oder Folge-ADR
  mit `supersedes` (Accepted-ADRs werden nie überschrieben).
- **Eine ADR entsteht nicht nur aus Architektur-Fragen.** Auch ein
  Rollen-Konflikt endet als ADR, wenn sein Verdikt bestritten wird — das
  Terminal des Konflikt-Pfads
  ([Modul 8](modul-08-agentenrollen.md#konflikt-pfad-als-rollen-sequenz-modul-8)):
  Die Entscheidung wird immutabel, Widerspruch braucht danach eine Folge-ADR
  mit neuer Evidenz.
- **Eine Gate-Erweiterung ist nicht automatisch ein ADR-Anlass.** Nimmt
  `make gates` ([Modul 13](modul-13-quality-gates.md)) einen bereits
  existierenden, unabhängig lauffähigen Wächter neu in den
  PR-blockierenden Satz auf, ist das keine neue Entscheidung — ein Verweis
  auf die ADR reicht, die den Wächter ursprünglich trägt. Existiert keine,
  trägt schon die *Einführung* des Wächters selbst eine, nicht seine
  spätere Aufnahme in `make gates`. Eine neue Fehlerklasse, ein neuer Scope
  oder ein Widerspruch zu einer bestehenden ADR braucht weiterhin eine
  eigene ADR. Im Zweifel: ADR.

### Hard Rule für Accepted-ADRs

Begriff *Hard Rule* siehe Glossar in [`grundlagen-begriffe.md`](grundlagen-begriffe.md#kernbegriffe).

**Eine ADR mit Status `Accepted` wird nicht inhaltlich überschrieben.**
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
explizitem Verweis auf die abgelöste oder geschärfte Vorgängerin.

Wirkung: ADRs sind Geschichtsdokumente, kein Wiki. Reviewer-Agent kann
auf ältere Entscheidungen vertrauen, ohne Versionsstände zu vergleichen.

### Nachzug ist keine Überschreibung

Die Hard Rule verbietet eine *inhaltliche* Überschreibung — nicht jede
spätere Berührung einer `Accepted`-ADR. Zwei Fälle bleiben davon
ausdrücklich getrennt, beide ohne Folge-ADR:

- **Referenz-/Pfad-Nachzug.** Eine ADR verweist über eine Kennung, nicht
  über eine Adresse (dieselbe Doktrin wie bei einfrierenden Artefakten,
  [`grundlagen-harness-dateien.md`](grundlagen-harness-dateien.md)) — wo
  das nicht möglich ist, etwa bei einem Verweis auf ein vergängliches
  Artefakt wie einen Review-Report, trägt die ADR einen deklarierten
  Vermerk: Der Verweis *verfällt*, wenn sein Ziel verschwindet; die
  Entscheidung bleibt davon unberührt gültig.
- **Template-Feld-Nachzug.** Gewinnt die Ziel-Form ein neues Pflichtfeld
  (etwa den Re-Evaluierungs-Trigger oben), tragen bestehende `Accepted`-ADRs es
  nicht rückwirkend. Ein Grandfathering-Vermerk (Feld fehlt, Datum des
  Feld-Zugangs) reicht.

Beide Fälle ändern die Entscheidung nicht — nur ihre Adresse oder ihre
Form. Was mehr ist als das, bleibt bei der Hard Rule: Folge-ADR mit
`supersedes`. Eine Formulierung wie „ist ja nur ein Pfad-Nachzug" für eine
tatsächlich geänderte Begründung ist selbst ein Bruch der Hard Rule, keine
Ausnahme von ihr.

### Regeln gegen typische Fehlannahmen (Modul 4)

- **"Eine ADR begründet eine Anforderung."** — Nein. ADRs begründen die *Lösung*. Anforderungen begründet die Spec. Wer ADRs zur Spec macht, kann später keine Architektur ohne Lastenheft-Änderung wechseln.
- **"Wenn ich die Entscheidung ändere, schreibe ich die ADR um."** — Hard Rule: Accepted-ADRs werden nicht überschrieben. Folge-ADR mit `supersedes ADR-N`. Sonst kann der Reviewer-Agent nicht auf ältere Entscheidungen vertrauen.
- **"Eine ADR ohne Fitness Function ist eine ADR."** — Eine ADR ohne Fitness Function ist eine Absichtserklärung. Wer architecture fitness im Kopf hat, schreibt parallel den ArchUnit-Test.
- **"MADR ist Pflicht."** — MADR ist ein Format unter mehreren (auch Nygard, Tyree/Akerman). Wichtig ist, dass dein Repo *eines* konsequent benutzt.
- **"Architektur ist Bilder zeichnen."** — Diagramme sind *eine* Output-Form, nicht die Sache selbst. Architektur heißt in diesem Regelwerk: *Entscheidungen mit Begründung (ADR), prüfbar gemacht (Fitness Function), versioniert (Accepted-Hard-Rule)*. Ein Diagramm ohne ADRs hinter sich ist Wandtapete; eine ADR ohne Fitness Function ist Absichtserklärung. `spec/architecture.md` ist explizit *diagrammatisch und enthält keine eigenen Anforderungen* (siehe Spec-Stratifizierung in [`grundlagen-source-precedence.md`](grundlagen-source-precedence.md#spec-stratifizierung)) — genau weil sonst Bilder anfangen würden, die ADR-Schicht zu ersetzen. Die `ARC-*`-Kennungen der Sicht ändern daran nichts: Sie adressieren Komponenten, damit man auf sie zeigen kann, und behaupten selbst nichts ([`grundlagen-source-precedence.md` §ID-Schema](grundlagen-source-precedence.md#id-schema-als-klammer)).
- **"ADRs sind Dokumentation, nicht Constraints."** — Eine ADR ohne maschinelle Durchsetzung ist eine *Absichtserklärung*, die der Implementer-Agent freundlich liest und dann ignoriert, wenn ein anderer Pfad "einfacher" wirkt. Eine ADR *mit* Fitness Function ist ein Constraint — die Layering-Regel, die ArchUnit dem Agenten als roten Build entgegenhält. Die Übersetzung (ADR-Satz → Werkzeug → Make-Target → Failure-Beispiel) steht kompakt in [Modul 13 §Fitness Function aus einem ADR-Satz](modul-13-quality-gates.md#adr-zur-fitness-function). Wer das nicht macht, dokumentiert *Hoffnung*.

### Ziel-Form: ADR (MADR)

Die Form liefert die Vorlage [`templates/docs/plan/adr/NNNN-titel.template.md`](../templates/docs/plan/adr/NNNN-titel.template.md):
Kopf (Status · Datum · Bezug · Supersedes) plus Body (Kontext · Verglichene
Alternativen · Entscheidung · Konsequenz mit Fitness Function). Operative
Regeln zur Form:

- Der Kontext *referenziert* die Anforderung, wiederholt sie nicht — und zwar
  **aufwärts** auf die stabilere Quelle, nie abwärts auf einen Slice
  ([`grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)](grundlagen-referenz-richtung.md#referenz-richtung-sdp-wer-darf-wen-referenzieren)).
- Mindestens drei Verglichene Alternativen, jede mit Trade-off.
- Jede Entscheidung mit Architektur-Wirkung bekommt eine Fitness Function —
  sonst ist sie Absichtserklärung.
- `Accepted` wird nie überschrieben — Korrektur = Folge-ADR mit `Supersedes`.
