## Modul 3 — Die Spec: Lastenheft, Spezifikation, Architektur

<!-- Quelle: [01-spec-und-architektur/modul-03-spec.md](https://github.com/pt9912/ai-harness-course/blob/v6.3.1/kurs/de/01-spec-und-architektur/modul-03-spec.md) -->

### Harness-Einordnung (Modul 3)

Spec = *inferential feedforward* (siehe
[`grundlagen/klassifikation.md`](grundlagen-klassifikation.md)).
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
[`grundlagen-begriffe.md#kernbegriffe`](grundlagen-begriffe.md#kernbegriffe):
LLM ist *stateless*). Wer die Metapher zu weit treibt, erwartet
"Mitlernen" — und plant Reviews, als würden sie *einmal* erklärt
ausreichen. Sie reichen nicht. Jeder Lauf beginnt bei Null.

### Regeln gegen typische Fehlannahmen (Modul 3)

- Happy Path widerlegt nur die These "es funktioniert gar nicht". Boundary und Negative widerlegen die stillen Annahmen, *die ein Agent am liebsten als selbstverständlich behandelt*.
- Im Gegenteil: ein Satz "das System *darf nicht* …" spart später drei Reviews. Negativ ist genauso präzise wie positiv.
- Nein, Performance gehört in den nichtfunktionalen Block der Spec (oder in `spec/spezifikation.md`, wenn stratifiziert). Der ADR begründet, *wie* man die Schwelle einhält.
- Was nicht explizit ausgeschlossen ist, baut der Agent plausibel mit. Das ist die häufigste Quelle für "wir hatten das nie gefordert"-PRs.
- Falsch. Lopopolos Maxime *"Was der Agent nicht im Kontext erreicht, existiert für ihn nicht"* ist ein Plädoyer *für* Kontext-Verfügbarkeit — und sagt damit, dass Spec und Prompt *unterschiedliche* Lebenszyklen haben: Spec wird *gepflegt* (Versions-Geschichte, Bezüge, Audit), Prompt wird *für einen Lauf zusammengestellt*. Was im Prompt steht, aber nicht in der Spec, gilt nur für *diesen* Lauf — der nächste Agent sieht es nicht. Das Muster (Spec sagte *speichert*, Agent baute PostgreSQL) wäre mit einem Mega-Prompt nicht besser geworden — der Prompt würde im nächsten Lauf vergessen.

### Ziel-Form: Akzeptanzkriterium

Anforderungen leben im Lastenheft
([`templates/spec/lastenheft.template.md`](../templates/spec/lastenheft.template.md)).
Ein funktionales Kriterium trägt eine `<PREFIX>-FA-NNN`-ID und dann drei Pfade
im Given/When/Then-Stil — **Happy · Boundary · Negative** — plus einen
**Out-of-Scope**-Block. Vagen Satz zuerst auf Mehrdeutigkeiten prüfen (*was
genau · welche Felder · welcher Speicherort*), bevor die Pfade formuliert
werden; das Negative (`darf nicht …`) spart die spätere Review.

Das Lastenheft ist die **Decke** der Straten-Ordnung: Es referenziert nur
innerhalb der eigenen `LH-*`-Reihe — keine ADRs, Slices, Carveouts, Wellen,
und auch nicht `spezifikation.md` oder `architecture.md`
([`grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)](grundlagen-referenz-richtung.md#referenz-richtung-sdp-wer-darf-wen-referenzieren);
die drei Spec-Zeilen der Matrix tragen in jeder fremden Spalte ein ❌ **ohne**
Kontext-Ausnahme).

**Das gilt in jedem Abschnitt, auch in der Historie — und für alle drei
Straten gleich, kein Stratum nimmt seine Provenance-Sektion aus.** Eine
Historie-Zeile ist ein Protokoll und wird nicht rückwirkend geändert; ein
dort genannter ADR-Verweis zeigt nach einer Supersedure dauerhaft auf eine
Entscheidung, die nicht mehr gilt. Wer eine Anforderung mit einer ADR
begründet, hat die Entscheidung zur Anforderung gemacht — und wer den
auslösenden Slice in der Historie nennt, tut dasselbe eine Zeile später.

Der Anlass geht damit nicht verloren, er liegt nur am richtigen Ende. Beim
Lastenheft ist es der **externe CR** in der Verweis-Spalte; wer im Repo
bemerkt hat, dass er nötig wird, hält es auf seiner Seite fest (Closure-Notiz
des Slice). Bei Technik und Sicht deklariert die ADR ihre Wirkung aufwärts in
`Schärft:`. Das Lastenheft sagt, **was** zugesagt ist — nicht, wer es bemerkt
hat.

### Ziel-Form: Spezifikation

Technische Festlegungen leben in der Spezifikation
([`templates/spec/spezifikation.template.md`](../templates/spec/spezifikation.template.md)):
Algorithmen und Datenflüsse · Datenstrukturen und Schemas · Defaults und
Konstanten · Fehler-Codes und Logging-Felder · Metriken und Tracing-Felder ·
externe Verträge · Historie. Operative Regeln:

* **Fortschreibbar** — kein Change Request nötig; eine ADR darf die
  Spezifikation schärfen, das Lastenheft nicht
  ([`grundlagen-source-precedence.md` §Spec-Stratifizierung](grundlagen-source-precedence.md#spec-stratifizierung)).
* **Präzisieren, nie erweitern** — Konfliktregel *Lastenheft › Spezifikation ›
  Architektur*.
* **Obligatorisch** — alle drei Straten sind Pflicht. Wer technische
  Festlegungen ins Lastenheft faltet, verschiebt nicht Inhalt, sondern deren
  Änderungs-Prozess: Sie wären dort abnahmebindend und nur per Change Request
  änderbar, und keine ADR dürfte sie je schärfen. Die Spezifikation existiert
  deshalb auch dünn; ein Repo mit zwei Straten deklariert das als `MR-<NNN>`.
* **Kein ADR-Verweis, auch nicht in der Historie** — die Spezifikation steht
  im Stabilitäts-Rang über der ADR (*Vertrag › Technik › Sicht › ADR ›
  Slice*), Referenzen zeigen nur aufwärts
  ([`grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)](grundlagen-referenz-richtung.md#referenz-richtung-sdp-wer-darf-wen-referenzieren)).
  Die **Decken-Regel gilt für alle drei Spec-Straten**: Kein Spec-Dokument
  nennt eine ADR oder einen Slice. Welche ADR eine Festlegung schärft,
  deklariert **die ADR** in ihrem `Schärft:`-Feld — die einzige Kante, und
  sie zeigt aufwärts. Wer wissen will, welche ADR eine Spec-Stelle bewegt
  hat, sucht in den `Schärft:`-Feldern.
* **Zwei Kennungs-Arten** — ein Abschnitt, der eine einzelne Anforderung
  technisch ausführt, trägt deren Verfeinerung (`LH-FA-03.a` zu `LH-FA-03`).
  Alles Übrige — Datenschemas, Defaults, Fehler-Codes, Metrik-Felder, externe
  Verträge — trägt eine `SPEC-<NNN>`, weil es keine einzelne Anforderung
  verfeinert, aber referenzierbar sein muss: Ohne Kennung kann eine ADR im
  `Schärft:`-Feld nur auf den Abschnitt zeigen, und zwei Defaults im selben
  Abschnitt werden ununterscheidbar. Eine `SPEC-*` ist **keine**
  Anforderungs-ID; sie benennt, was gilt, und verspricht nichts
  ([`grundlagen-source-precedence.md` §ID-Schema](grundlagen-source-precedence.md#id-schema-als-klammer)).
* **Kein Kopf-Datum, kein Kopf-Status** — die Spezifikation trägt ihre
  Änderungen in der Historie; die letzte Zeile ist das Datum, ein
  Frische-Marker im Kopf (`Letzte Änderung`) wäre dasselbe Datum ein zweites
  Mal, und zwei Felder für eines driften. Das unterscheidet sie von der Sicht,
  die keine Historie hat (§Ziel-Form: Architektur-Sicht). Einen eigenen Status
  trägt sie nicht: Verbindlich ist sie, solange das Lastenheft es ist — dessen
  Status steuert die Verbindlichkeit der IDs, die sie präzisiert. Ihre
  Historie führt Datum und Änderung — keine Version: Versionen hat der
  Vertrag, dessen Version-Bump der Fußabdruck des Change Requests ist
  ([`grundlagen-source-precedence.md` §Spec-Stratifizierung](grundlagen-source-precedence.md#spec-stratifizierung));
  die Technik ist fortschreibbar ohne ihn.

### Ziel-Form: Architektur-Sicht

Die Sicht folgt der Vorlage
[`templates/spec/architecture.template.md`](../templates/spec/architecture.template.md):
Komponenten- und Sequenzsicht ohne **eigene Anforderungen** (Sicht-Stratum,
[`grundlagen-source-precedence.md` §Spec-Stratifizierung](grundlagen-source-precedence.md#spec-stratifizierung)).
Operative Regeln, die die Vorlage nicht selbst erzwingt:

* Derivativ — Konfliktregel *Lastenheft › Spezifikation › Architektur*; die
  untere Schicht präzisiert, erweitert nie.
* Sprach- und meilensteinfrei — **darf** Pfade zu **Code-Modulen**
  referenzieren (`src/service/`), aber **keine** Wellen, Slices,
  Commit-Hashes oder Closure-Daten. Die Erlaubnis ist keine Pflicht: Eine
  Sicht, die ihre Komponenten nur über Rollen und `ARC-*` führt, ist ebenso
  konform.
* Keine Historie, nur `**Letzte Änderung:**` im Kopf — ein Frische-Marker,
  kein Protokoll. Vertrag und Technik führen eine, weil ihr Änderungs-Prozess
  einen benennbaren Urheber hat (externer Change Request bzw. schärfende ADR);
  die Sicht hat keinen, jede ihrer Änderungen folgt aus einer Änderung
  darüber. Eine Verweis-Spalte trüge hier nichts Zulässiges.
* Kein ADR-Bezug: Die Sicht steht im Stabilitäts-Rang **über** der ADR
  (*Vertrag › Technik › Sicht › ADR › Slice*), normative Referenzen zeigen nur
  aufwärts ([`grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)](grundlagen-referenz-richtung.md#referenz-richtung-sdp-wer-darf-wen-referenzieren)).
  Welche ADR eine Aussage der Sicht verbindlich macht, deklariert **die ADR**
  aufwärts in ihrem `Schärft:`-Feld.
* `ARC-*` widerspricht dem nicht: Komponenten und Schnittstellen tragen
  Kennungen, damit ein Slice sagen kann, *welche* Komponente er
  berührt. Das ist eine Adresse, keine Zusage — die Sicht bleibt derivativ,
  weil eine Kennung nichts verspricht. Wer unter einer `ARC-<NNN>` eine
  Anforderung formuliert, hat nicht die Kennung falsch benutzt, sondern eine
  Anforderung am falschen Ort abgelegt
  ([`grundlagen-source-precedence.md` §ID-Schema](grundlagen-source-precedence.md#id-schema-als-klammer)).

### Spec-Stratifizierung — Drei Schichten (Modul 3)

Jede Spec zerfällt in drei obligatorische Schichten mit eigener Precedence:
`lastenheft.md` (vertragliches *Was*) › `spezifikation.md` (präzisiertes
*Wie genau*) › `architektur.md` (strukturelles *Wodurch*). Konfliktregel:
**Lastenheft sticht Spezifikation sticht Architektur** — die untere
Schicht darf *präzisieren*, nie *erweitern*. Vollform (Straten-Klassen,
Referenz-Richtung, Durchsetzung) in
[`grundlagen-source-precedence.md` §Spec-Stratifizierung](grundlagen-source-precedence.md#spec-stratifizierung).
Vorlagen: [`spec/`-Templates](../templates/spec/).

