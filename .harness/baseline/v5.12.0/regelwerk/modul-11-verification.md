## Modul 11 — Verification Harness

<!-- Quelle: [04-qualitaet/modul-11-verification.md](https://github.com/pt9912/ai-harness-course/blob/v5.12.0/kurs/de/04-qualitaet/modul-11-verification.md) -->

### Begriffe: Pre-completion Checklist Middleware und DoD-Verletzung

* **Pre-completion Checklist Middleware** — eine vom Implementer-Agent
  selbst durchlaufene Checkliste *vor* der "fertig"-Meldung. Sie ist
  Schritt 8 des 8-Schritt-Workflows (siehe
  [Modul 9 §Minimal Agent Workflow](modul-09-implementierung.md#minimal-agent-workflow-8-schritte)).
  In diesem Modul betrachten wir sie als *Eingabe* für die Verifikation:
  was die Checkliste *behauptet*, ist von der Verifikation maschinell
  oder semantisch zu *bestätigen*. Behauptung ohne Bestätigung ist die
  häufigste Verifier-Lücke.
* **DoD-Verletzung** — Differenz zwischen DoD-Punkten des Slice
  (Modul 5) und tatsächlichem Code-/Artefakt-Stand. Wichtig: eine
  DoD-Verletzung ist *kein* Review-Finding (Reviewer prüft gegen
  Plan/ADR, nicht gegen DoD/Spec) — sie ist eine eigene Klasse, die
  *nur* die Verifikation fängt.

### Harness-Einordnung (Modul 11)

Verifikation = primär *inferential feedback* in der Behaviour-Kategorie,
unterstützt durch *computational feedback* (Fitness Functions für die
Architecture-Fitness-Kategorie). Dies ist die anspruchsvollste Schicht
— und laut Böckeler die am wenigsten ausgereifte. Siehe
[`grundlagen/klassifikation.md`](grundlagen-klassifikation.md).

### Kernidee (Modul 11)

Verifikation ist die Stelle, an der der Harness *gegen sich selbst*
misst: "Hat das, was gebaut wurde, das umgesetzt, was geplant war?" —
nicht: "Ist es gut?"

### Regeln gegen typische Fehlannahmen (Modul 11)

- Tests prüfen ob *Code tut, was Tests testen*. Verifikation prüft, ob *Code tut, was Plan/DoD/Spec verlangt*. Lücken zwischen Tests und Spec sind genau das, was Verifikation findet.
- Nein. Reviewer hat *Plan + ADR*. Verifier hat *DoD + Spec + Plan*. Andere Eingabe, andere Findings.
- Falsch. Die wahrscheinlichere Erklärung: Reviewer hat gegen einen veralteten Plan geprüft, oder der Plan hat eine DoD-Lücke. Architect klärt — *nicht* "wir nehmen das mildere Ergebnis".

### Fitness Function ohne Standard-Tool (Modul 11)

Wenn eine ADR-Aussage kein Standard-Tool zum Prüfen hat (Beispiel:
„Closure-Note mit mindestens zwei Sätzen"), heißt Verifizieren: die
Fitness Function selbst bauen. Der Ablauf:

- **Operationalisieren** — die eigentliche Arbeit: aus der ADR-Aussage
  (*was*) eine prüfbare Form (*prüfbar was*) machen. „Mindestens zwei
  Sätze" wird z. B. zu „Frontmatter-Schlüssel `closure_note` vorhanden ·
  ≥ 2 Satzendezeichen außerhalb von Code · keine der bekannten Floskeln".
- **Sensor-Schicht nach Kosten wählen:**

| Option | Kosten | Wann sinnvoll |
|---|---|---|
| Pre-commit-Hook (Autoren-Maschine) | niedrig | nur lokale Disziplin gefragt |
| Make-Target im `make gates`/`verify`-Block | mittel | auch CI soll prüfen — Standardweg |
| Doku-Konsistenz-Agent (Modul 15) | hoch | semantische Prüfung nötig (Floskel-Erkennung) |

- **Skript + Gate verdrahten:** ID-Kommentar zeigt die ADR; eine DoD-/
  Closure-Frage hängt an `verify:` (nicht `make gates` — das ist für
  Code-Architektur-Fragen).
- **Inferentielle Schicht für Semantik:** deterministische Struktur
  deterministisch prüfen, semantische „Inhalt vs. Floskel"-Erkennung
  inferentiell — denn Floskeln wie „war ganz okay, läuft jetzt" sind
  syntaktisch zwei Sätze. Prompt-Anker in
  [`.harness/skills/closure-note-reviewer.md`](../templates/.harness/skills/closure-note-reviewer.template.md)
  (Schwester-Skill zum Reviewer, Modul 10).
- **Hard Rule in zwei Quadranten:** *inferential feedforward*
  (`AGENTS.md` sagt es) + *computational feedback* (Make-Target prüft
  es); der Implementer-Agent läuft `make verify-*` **selbst** vor der
  „fertig"-Meldung (Pre-completion Checklist, Modul 9 Schritt 8). So
  fängt der Verifier genau das, was Tests nicht prüfen und der Reviewer
  übersieht — die fehlende Closure-Note ist kein Diff-Symptom.

**Ein selbstgebautes Gate ist auf Zeit gebaut.** „Kein Standard-Tool prüft das"
ist eine Aussage über *heute*. Erscheint später eines, ist die Frage nicht, ob
das eigene Skript stört, sondern ob das Werkzeug eine **Obermenge** ist — und
das hat drei Teile, die einzeln nachgewiesen werden: dieselbe
**Kandidaten-Menge** (welche Dateien werden überhaupt geprüft), dieselben
**Bedingungen**, und dieselbe **Schwelle, wie die ADR sie setzt** — nicht die
Vorbelegung des Werkzeugs. Ein Gate, das schärfer ist als seine ADR, ist
genauso falsch wie eines, das lascher ist; beides prüft eine Entscheidung, die
niemand getroffen hat.

Der Nachweis ist nicht der Datenblatt-Vergleich, sondern **je Verstoßklasse ein
Break-Test mit beiden Sensoren nebeneinander**, plus der unveränderte Bestand,
auf dem beide schweigen müssen. Ist das Werkzeug Obermenge, wird das Skript
retired — sonst benennt man die fehlende Klasse und behält es. Beides ist ein
Ergebnis; was nicht zählt, ist die Vermutung.

Und eine dritte Antwort gibt es auch: Ein Skript kann **einen anderen
Konsumenten bekommen**, als es hatte. Dient es inzwischen der Lehre, dem
Onboarding oder einem Fixture, dann trägt es nicht mehr die Deckung, sondern
eine Rolle — und die gehört dann ausgeschrieben, sonst liest die nächste Person
es weiterhin als Gate
([`grundlagen-harness-dateien.md` §Jedes Artefakt hat einen Konsumenten](grundlagen-harness-dateien.md#jedes-artefakt-hat-einen-konsumenten)).

