## Modul 10 — Review Harness

*Quelle: [04-qualitaet/modul-10-review-harness.md](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/04-qualitaet/modul-10-review-harness.md)*

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
[`grundlagen/klassifikation.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/klassifikation.md)).
Teurer als ein Linter, billiger als Verifikation. Adressiert primär die
Maintainability-Kategorie.

### Kernidee (Modul 10)

Ein Review ohne Kategorisierung ist eine Mängelliste. Ein Review mit
Kategorisierung ist eine Entscheidungsvorlage.

### Worked Example: eine Reviewer-Skill-Datei schreiben

Ein Reviewer-Agent ohne Skill-Datei driftet zwischen Sessions. Dieselbe
Eingabe → unterschiedliche Findings, unterschiedliche Kategorien.
Skill-Dateien leben in `.harness/` und sind das Repo-spezifische
"worauf achtest du" eines Agenten.

**Schritt 1 — Pfad und Kopf:**

```
.harness/skills/reviewer.md
```

```markdown
# Reviewer-Skill — DocSearch

* Status: Accepted
* Bezug: ADR-0007, AGENTS.md §"Review-Regeln"
* Gilt für: `agent-review`-Make-Target
```

**Schritt 2 — Eingangs-Kontext explizit machen.** Was der Reviewer
*immer* mitbringt, bevor er den Diff liest:

```markdown
## Kontext-Eingang (Pflicht)

- Diff des PR
- `spec/lastenheft.md` (für referenzierte LH-IDs)
- ADRs, deren ID im PR oder Commit-Message vorkommt
- AGENTS.md §"Hard Rules"
- vorherige Findings am gleichen Modul (letzte 5 PRs)
```

Ohne diesen Block sieht der Reviewer den Code, aber nicht *die Verträge,
gegen die er prüft*.

**Schritt 3 — Kategorien-Regeln *für dieses Repo*.** Nicht generisch,
sondern konkret:

```markdown
## Klassifikation

**HIGH** — eines der folgenden:
- ADR-Verstoß (Layer, Tool, Hard Rule)
- Sicherheits-Anti-Pattern (Injection, fehlende Auth-Prüfung)
- Korrektheitsfehler im *kritischen* Pfad (Index-Schreiben, Auth)
- Suppression eines Gates (#noqa, //nolint, [SuppressMessage]) ohne ADR

**MEDIUM** — eines der folgenden:
- unklare Fehlerbehandlung am Rand des Spec-Bereichs
- fehlende Negativtests bei neuem öffentlichen Vertrag
- Wiederholung eines Musters, das schon zweimal LOW war

**LOW** — stilistisch unschön ohne semantische Auswirkung,
einmalige Tippfehler, unbenutzte Imports.

**INFO** — Hinweis ohne erwartete Aktion (z. B. "diese Stelle hat
ein passendes ArchUnit-Pendant, das du nicht kennst").
```

Beachte: drei Kategorien-Anker (HIGH/MEDIUM/LOW) haben *jeweils* eine
konkrete Liste. INFO ist bewusst kurz — INFO ist Ergänzungs-Kanal, nicht
Hauptkanal.

**Schritt 4 — Anti-Pattern und "Was bist du nicht".** Verhindert, dass
der Reviewer zum zweiten Implementer wird:

```markdown
## Was dieser Skill NICHT macht

- Keine Lösungsvorschläge ("schreib das so") — Reviewer kategorisiert,
  Implementer entscheidet.
- Kein Refactoring-Vorschlag, der über den Diff hinausgeht.
- Keine Verifikation gegen DoD — das ist Verifier-Aufgabe (Modul 11).
- Keine Validation gegen reale Bedürfnisse — das ist Validator-Aufgabe.

Wenn etwas auffällt, das in diese Kategorien gehört, ein INFO-Finding
mit Verweis auf die zuständige Rolle.
```

**Schritt 5 — Output-Schema fixieren.** Findings sind strukturiert, nicht
Fließtext:

```markdown
## Output-Schema

Jedes Finding:

- `kategorie`: HIGH | MEDIUM | LOW | INFO
- `quelle`: ADR-ID, LH-ID, Hard-Rule-Name oder "Maintainability"
- `pfad`: Datei:Zeile
- `befund`: 1–2 Sätze, beobachtbar, ohne Lösungsvorschlag
- `verifizierbar`: ja/nein — gibt es einen Gate-Lauf, der es bestätigen würde?

Zusätzlich am Ende: eine Zeile "geprüft, ohne Befund" pro betrachtetem
Verzeichnis (Negativbefund-Zeile — siehe Modul 10 §"Reviewer berichtet
auch, was er nicht gefunden hat").
```

**Schritt 6 — Steering-Loop-Eintrag.** Skills sind nicht statisch:

```markdown
## Pflege

Bei dreimaligem Auftreten desselben Findings:
- ist die Kategorie noch richtig? → Klassifikation schärfen
- gibt es einen ADR/AGENTS.md-Eintrag, der das verhindert hätte?
  → Folge-ADR oder AGENTS.md-Update
- gibt es eine Fitness Function, die das prüfen würde?
  → Modul 13, Gate hinzufügen

Skill-Datei selbst wird **nicht** überschrieben, sondern versioniert
(siehe ADR-Hard-Rule, Modul 4).
```

Sechs Schritte, eine reproduzierbare Reviewer-Rolle. Vergleichbares
Skill-Pattern für *Verifier* und *Validator* in Modul 11 bzw. in
[Modul 8 §"Konfliktfall"](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/03-agenten/modul-08-agentenrollen.md).

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
[`review-report.template.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/templates/docs/reviews/review-report.template.md);
abgelegt wird ein Report pro Lauf unter `docs/reviews/`, Folgeläufe
als neue Datei statt Überschreibung.

### Regeln gegen typische Fehlannahmen (Modul 10)

- Reviewer kategorisiert. Vorschläge "wie ich es geschrieben hätte" sind nett, aber kein Reviewer-Ergebnis.
- Implementer arbeitet sequentiell ab und bleibt am LOW hängen. HIGH zuerst, immer.
- Verhalten driftet zwischen Sessions. Jeder Reviewer-Agent braucht eine Skill-Datei in `.harness/` mit "worauf achtest du in diesem Repo".
- Genau das belohnt Inkonsistenz. Stattdessen: Skill schärfen, bis die Klassifikation reproduzierbar ist.

