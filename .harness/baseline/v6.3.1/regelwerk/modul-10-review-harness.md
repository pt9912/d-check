## Modul 10 — Review Harness

<!-- Quelle: [04-qualitaet/modul-10-review-harness.md](https://github.com/pt9912/ai-harness-course/blob/v6.3.1/kurs/de/04-qualitaet/modul-10-review-harness.md) -->

### Drei Review-Arten — wogegen wird geprüft

Die drei Review-Arten unterscheiden sich nicht im *Wie* (alle liefern
kategorisierte Findings), sondern im *Wogegen* und im *Wann*:

* **Plan-Review** prüft den Plan eines Slices gegen Spec und
  Accepted-ADRs — *bevor* implementiert wird. Es gibt noch keinen
  Diff; Eingabe ist der Plan selbst (Modul 9, Schritt 2).
* **Design-Review** prüft den Lösungs-Schnitt gegen die Architektur:
  Layer-Grenzen, Schnittstellen, ADR-Verträglichkeit einer neuen
  Komponente — bevor die Details festgezurrt sind.
* **Code-Review** prüft den fertigen Diff gegen Plan und Konventionen
  (AGENTS.md, Hard Rules) — die Findings-Kategorien dieses Moduls.

Merkregel: je früher die Review-Art, desto billiger das Finding —
ein Plan-Review-HIGH kostet eine Plan-Korrektur, dasselbe Finding im
Code-Review kostet den ganzen Implementierungs-Lauf.

### Finding-Kategorien

| Kategorie | Bedeutung |
|---|---|
| HIGH | blockiert Merge: Sicherheits-, Korrektheits- oder ADR-Verstoß |
| MEDIUM | sollte vor Merge geklärt werden |
| LOW | nice-to-fix, blockiert nicht |
| INFO | Hinweis, keine Aktion erwartet |

### Harness-Einordnung (Modul 10)

Review = *inferential feedback* (siehe
[`grundlagen/klassifikation.md`](grundlagen-klassifikation.md)).
Teurer als ein Linter, billiger als Verifikation. Adressiert primär die
Maintainability-Kategorie.

Die *Kategorisierung* eines Findings bleibt inferential — kein Gate
entscheidet, ob ein Verstoß wirklich HIGH ist. Die *Deckung* dagegen ist
mechanisierbar: trägt ein `done/`-Slice mit Review-Zusage tatsächlich
einen Report? Das ist *computational feedback* und prüfbar, ohne die
Kategorisierung selbst zu bewerten — Werkzeug-Beispiel: das d-check-Modul
`reviews` (Review-Report-Deckung für `done/`-Slices).

### Kernidee (Modul 10)

Ein Review ohne Kategorisierung ist eine Mängelliste. Ein Review mit
Kategorisierung ist eine Entscheidungsvorlage.

### Ziel-Form: Reviewer-Skill

Ein Reviewer-Agent ohne Skill-Datei driftet zwischen Sessions — und zwischen
**Rolleninhabern** ([Modul 8](modul-08-agentenrollen.md#rollen-regeln-modul-8)):
Füllen mehrere Menschen die Rolle, ist dieselbe Abweichung kein
Nicht-Determinismus, sondern ein **Dissens** — auch dann wird der Skill
geschärft, nicht die mildere Lesart gewählt (gleiche
Eingabe → andere Findings/Kategorien). Die Skill-Datei liegt in
`.harness/skills/reviewer.md` und ist das repo-spezifische „worauf
achtest du"; Vorlage
[`templates/.harness/skills/reviewer.template.md`](../templates/.harness/skills/reviewer.template.md)
(für die engere Closure-Note-Prüfung der Schwester-Skill
`closure-note-reviewer.md`, Modul 11). Operative Pflichtteile:

- **Kontext-Eingang (Pflicht):** Diff · `spec/lastenheft.md` · ADRs, deren
  ID im PR/Commit vorkommt · `AGENTS.md` §Hard Rules · vorherige Findings
  am gleichen Modul. Ohne den Block sieht der Reviewer Code, aber nicht
  die Verträge, gegen die er prüft.
- **Klassifikation repo-konkret**, nicht generisch: HIGH/MEDIUM/LOW je
  eine konkrete Liste, INFO kurz (Ergänzungs-Kanal, nicht Hauptkanal).
  Die HIGH-Liste muss **mindestens zwei repo-spezifische Regeln** nennen,
  die ein generischer Skill nicht abdeckt — sonst greift bei einem realen
  Diff keines der Repo-HIGHs.
- **„Was dieser Skill NICHT macht":** kein Lösungsvorschlag, kein
  Refactoring über den Diff hinaus, keine Verifikation (Verifier, Modul
  11), keine Validation (Validator) — sonst wird der Reviewer zum zweiten
  Implementer. Auffälliges außerhalb → INFO-Finding mit Rollen-Verweis.
- **Output-Schema strukturiert** (`kategorie · quelle · pfad · befund ·
  verifizierbar · klasse`; `klasse` = stabile Kurz-Bezeichnung des
  Fehlermusters, speist den Steering-Loop-Zähler — siehe Pflege unten) plus
  je betrachtetem Bereich eine **Negativbefund-Zeile**
  („geprüft, ohne Befund"; eigene Sektion unten).
- **Pflege (Steering-Loop):** Das „dreimal" zählt der Skill nicht selbst —
  jeder Lauf steht für sich; gezählt wird über Finding-Klasse →
  Slice-Closure §7 → Eintrag ins Beobachtungs-Register. Bei dreimaligem gleichem Finding
  Klassifikation schärfen / Folge-ADR bzw. `AGENTS.md`-Update / Gate
  (Modul 13). Skill-Datei selbst wird **nicht** überschrieben, sondern
  versioniert (siehe ADR-Hard-Rule, Modul 4). Ein **neuer** HIGH-Eintrag trägt
  ab Einführung einen Auflösungs-Trigger oder *permanent* — dieselbe Disziplin
  wie bei Hard Rules
  ([Modul 13](modul-13-quality-gates.md#hard-rule-doku-disziplin)); der
  Altbestand bleibt ohne. Ein HIGH-Eintrag, der aus dem
  Steering Loop kam, trägt den Herkunfts-Anker `(seit welle-<NN>)` — ohne Welle
  `(seit slice-<NNN>)`
  ([`grundlagen-traceability.md` §Herkunfts-Anker](grundlagen-traceability.md#herkunfts-anker)).

Vergleichbares Skill-Pattern für *Verifier* und *Validator* in Modul 11
bzw. [Modul 8 §"Konfliktfall"](modul-08-agentenrollen.md).

### Reviewer berichtet auch, was er nicht gefunden hat

Ein Report, der nur Findings listet, ist nicht auditierbar: „keine
Findings in `internal/auth/`" und „`internal/auth/` nicht angesehen"
sehen identisch aus — eine leere Liste. Deshalb verlangt das
Output-Schema pro betrachtetem Bereich eine **Negativbefund-Zeile**
(„geprüft, ohne Befund"). Sie macht die Abdeckung des Laufs sichtbar,
ist die Grundlage für Vertrauen in ein grünes Review — und sie ist
der Teil des Reports, den ein Reviewer-Agent am ehesten weglässt,
weil ihn niemand einfordert.

Das Dokument-Gerüst für den **ganzen Report** — Kopf-Metadaten
(Review-Art, Gegenstand, Skill-Version, Modell, Eingangs-Kontext),
Findings nach Output-Schema, Negativbefunde, Kategorie-Summary,
Verdikt — liefert
[`review-report.template.md`](../templates/docs/reviews/review-report.template.md);
**Mit der Closure der Welle, die seinen Slice einsammelt, wandert der Report
vollständig ins Archiv — ohne Stub** ([Modul 6](modul-06-roadmap.md), Schritt
4). Er hat keine Identität jenseits seines Slice; wer ihn sucht, sucht ihn
unter dem Slice, den er geprüft hat. **In einem Repo ohne Wellen archiviert ihn
die Slice-Closure selbst** — dieselbe, die seinen Slice schließt, nach den
Paarungen, nach `done/slice-<NNN>-archiv.zip` ([Modul 6](modul-06-roadmap.md),
*Wann Arbeit eine Welle braucht*). Er wartet nicht auf ein Ereignis, das es
dort nie gibt. **Ein Rang-Dokument, das einen einzelnen
Report als Beleg verlinkt, hat damit ein Problem, das älter ist als das
Archiv** — es macht ein Zeitdokument zur Quelle. Die Aussage gehört an den
zitierenden Ort, die Report-Kennung bleibt im Text.

**Der Report ist ein Lauf-Beleg, kein Wissensspeicher**: Konsument ist der
Implementer im selben Zyklus, danach der Audit. Über Läufe hinweg wird er
nicht wieder gelesen — das steuerungsrelevante Signal ist die
**Finding-Klasse** (Summary-Zeile), die über die Slice-Closure §7 ins
**Beobachtungs-Register** wandert
([Modul 6](modul-06-roadmap.md#das-beobachtungs-register-modul-6)). Bedingung: die
Klassen-Bezeichnung ist über Läufe hinweg **stabil** — der Report kennt das
Register nicht, die Zuordnung zur `BEO-<NNN>` passiert erst bei der
Slice-Closure und braucht den wiedererkennbaren Namen. Ein Archiv-Scan ist
nicht nötig — die Häufung steht im Register.

Abgelegt wird ein Report pro Lauf unter `docs/reviews/`, Folgeläufe
als neue Datei statt Überschreibung.

### Regeln gegen typische Fehlannahmen (Modul 10)

- **Gegen "Reviewer ist ein zweiter Implementer":** Reviewer kategorisiert. Vorschläge "wie ich es geschrieben hätte" sind nett, aber kein Reviewer-Ergebnis.
- **Gegen "Findings ohne Prioritätssortierung":** Implementer arbeitet sequentiell ab und bleibt am LOW hängen. HIGH zuerst, immer.
- **Gegen "Reviewer-Agent läuft ohne Skill-Datei":** Verhalten driftet zwischen Sessions und zwischen Rolleninhabern. Jeder Reviewer-Agent braucht eine Skill-Datei in `.harness/` mit "worauf achtest du in diesem Repo".
- **Gegen "Bei zwei verschiedenen Kategorisierungen nehmen wir die mildere":** Genau das belohnt Inkonsistenz. Stattdessen: Skill schärfen, bis die Klassifikation reproduzierbar ist.

