## Modul 3 — Lastenheft und Spezifikation

*Quelle: [01-spec-und-architektur/modul-03-lastenheft.md](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/01-spec-und-architektur/modul-03-lastenheft.md)*

### Harness-Einordnung (Modul 3)

Spec = *inferential feedforward* (siehe
[`grundlagen/klassifikation.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/klassifikation.md)).
Sie ist die billigste Kontrolle: Was die Spec sauber ausschließt, kommt
im Review nicht mehr vor.

### Kernidee (Modul 3)

Ein Agent ist ein extrem buchstabengetreuer Praktikant. Was nicht in der
Spec steht, existiert für ihn nicht — Lopopolos Maxime: *"Was der Agent
nicht im Kontext erreicht, existiert für ihn nicht."* Was zweideutig in der Spec
steht, wird auf die für dich ungünstigste Weise interpretiert.

**Grenze der Metapher.** Die Praktikant-Metapher trägt nur die
*Buchstabentreue*. Anders als ein echter Praktikant **vergisst** der
Agent zwischen den Aufgaben — was nicht im Kontext steht, war für ihn
nie da (siehe Glossar in
[`grundlagen/konventionen.md#kernbegriffe`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/konventionen.md#kernbegriffe):
LLM ist *stateless*). Wer die Metapher zu weit treibt, erwartet
"Mitlernen" — und plant Reviews, als würden sie *einmal* erklärt
ausreichen. Sie reichen nicht. Jeder Lauf beginnt bei Null.

### Regeln gegen typische Fehlannahmen (Modul 3)

- Happy Path widerlegt nur die These "es funktioniert gar nicht". Boundary und Negative widerlegen die stillen Annahmen, *die ein Agent am liebsten als selbstverständlich behandelt*.
- Im Gegenteil: ein Satz "das System *darf nicht* …" spart später drei Reviews. Negativ ist genauso präzise wie positiv.
- Nein, Performance gehört in den nichtfunktionalen Block der Spec (oder in `spec/spezifikation.md`, wenn stratifiziert). Der ADR begründet, *wie* man die Schwelle einhält.
- Was nicht explizit ausgeschlossen ist, baut der Agent plausibel mit. Das ist die häufigste Quelle für "wir hatten das nie gefordert"-PRs.
- Falsch. Lopopolos Maxime *"Was der Agent nicht im Kontext erreicht, existiert für ihn nicht"* ist ein Plädoyer *für* Kontext-Verfügbarkeit — und sagt damit, dass Spec und Prompt *unterschiedliche* Lebenszyklen haben: Spec wird *gepflegt* (Versions-Geschichte, Bezüge, Audit), Prompt wird *für einen Lauf zusammengestellt*. Was im Prompt steht, aber nicht in der Spec, gilt nur für *diesen* Lauf — der nächste Agent sieht es nicht. Engage-Geschichte oben (Spec sagte *speichert*, Agent baute PostgreSQL) wäre mit einem Mega-Prompt nicht besser geworden — der Prompt würde im nächsten Lauf vergessen.

### Worked Example: vom vagen Satz zum prüfbaren Akzeptanzkriterium

**Ausgangstext (vage):**
> "Das System speichert die Konfiguration."

**Schritt 1 — Mehrdeutigkeiten markieren:** *speichert* (DB? Datei? Cache?), *die* (welche?), *Konfiguration* (welche Felder?).

**Schritt 2 — ID vergeben:** `LH-FA-CFG-001`.

**Schritt 3 — Happy Path konkret:**
> Given eine gültige Konfigurationsdatei `config.yaml` mit den Pflichtfeldern `name`, `version`,
> When das System startet,
> Then liest es die Datei *aus dem Arbeitsverzeichnis* und gibt `name@version` auf stdout aus.

**Schritt 4 — Boundary:**
> Given `config.yaml` ist leer,
> When das System startet,
> Then bricht es mit Exit-Code 2 ab und meldet `LH-FA-CFG-001: empty config`.

**Schritt 5 — Negative (zwei Sätze):**
> Given keine `config.yaml` existiert,
> When das System startet,
> Then bricht es mit Exit-Code 1 ab und schreibt keine Datei.
>
> Das System *darf nicht* Konfiguration in Datenbanken, externen APIs oder versteckten Verzeichnissen ablegen.

**Schritt 6 — Out-of-Scope:**
> Out-of-Scope (LH-FA-CFG-001): Schreiboperationen, Migration zwischen Versionen, Verschlüsselung.

Sechs Schritte, ein vollständig prüfbares Akzeptanzkriterium. Vergleich
mit dem Lab-Beispiel: [`/lab/example/spec/lastenheft.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/spec/lastenheft.md).

### Spec-Stratifizierung — Drei Schichten (Modul 3)

Nimm dein Mini-Feature aus der ersten Übung und
verteile seinen Inhalt auf drei Dateien — `lastenheft.md` (vertragliches
*Was*), `spezifikation.md` (präzisiertes *Wie genau*), `architektur.md`
(strukturelles *Wodurch*). Pflicht pro Schicht: *ein* Inhalt, der dort
zwingend gehört, und *ein* Inhalt, der dort fehl am Platz wäre (z. B.
gehört "Antwort als gültiges JSON" ins Lastenheft, "Service-Layer ruft
nie direkt die DB" in die Architektur). Formuliere zum Schluss die
*Konfliktregel*: Was gilt, wenn dieselbe Aussage in zwei Schichten
auftaucht (Lastenheft sticht Spezifikation sticht Architektur — die
untere Schicht darf *präzisieren*, nie *erweitern*)? Vorbild:
Spec-Stratifizierung in `c-hsm-doc`
([`grundlagen/fallstudien.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/fallstudien.md)).
Vorlagen: [`spec/`-Templates](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/templates/spec/).

