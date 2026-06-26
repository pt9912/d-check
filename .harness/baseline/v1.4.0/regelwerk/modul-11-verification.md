## Modul 11 — Verification Harness

*Quelle: [04-qualitaet/modul-11-verification.md](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/04-qualitaet/modul-11-verification.md)*

### Begriffe: Pre-completion Checklist Middleware und DoD-Verletzung

* **Pre-completion Checklist Middleware** — eine vom Implementation-Agent
  selbst durchlaufene Checkliste *vor* der "fertig"-Meldung. Sie ist
  Schritt 8 des 8-Schritt-Workflows (siehe
  [Modul 9 §Minimal Agent Workflow](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/03-agenten/modul-09-implementierung.md#minimal-agent-workflow-8-schritte)).
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
[`grundlagen/klassifikation.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/klassifikation.md).

### Kernidee (Modul 11)

Verifikation ist die Stelle, an der der Harness *gegen sich selbst*
misst: "Hat das, was gebaut wurde, das umgesetzt, was geplant war?" —
nicht: "Ist es gut?"

### Regeln gegen typische Fehlannahmen (Modul 11)

- Tests prüfen ob *Code tut, was Tests testen*. Verifikation prüft, ob *Code tut, was Plan/DoD/Spec verlangt*. Lücken zwischen Tests und Spec sind genau das, was Verifikation findet.
- Nein. Reviewer hat *Plan + ADR*. Verifier hat *DoD + Spec + Plan*. Andere Eingabe, andere Findings.
- Falsch. Die wahrscheinlichere Erklärung: Reviewer hat gegen einen veralteten Plan geprüft, oder der Plan hat eine DoD-Lücke. Architect klärt — *nicht* "wir nehmen das mildere Ergebnis".

### Worked Example: eine ADR-Aussage ohne fertiges Tool als Fitness Function bauen

**Ausgangs-ADR:** ADR-0011 sagt:
> "Der Implementation-Agent darf nur Slices in `done/` verschieben,
> wenn das Slice-Frontmatter ein Feld `closure_note` mit mindestens
> zwei Sätzen enthält (Lerneintrag-Pflicht; Modul 1 §Closure)."

Es gibt kein "Closure-Note-Linter". Diese Regel zu verifizieren heißt:
sie selbst bauen.

**Schritt 1 — Aussage in eine prüfbare Operationalisierung übersetzen.**
*"Mindestens zwei Sätze"* ist nicht selbsterklärend. Operationalisierung:

> Für jede Datei in `docs/plan/planning/done/*.md` gilt:
> - Frontmatter enthält Schlüssel `closure_note` (string).
> - Der String enthält mindestens **zwei Satzendezeichen** (`.`, `!`, `?`),
>   außerhalb von Code-Blöcken und Inline-Code.
> - Keiner der Sätze ist *leer* oder einer der bekannten Floskeln
>   ("see PR", "n/a", "siehe Ticket").

Operationalisierung ist die Stelle, an der Erschaffen passiert: die ADR
sagt *was*, die Operationalisierung sagt *prüfbar was*.

**Schritt 2 — Sensor-Schicht wählen.** Optionen, mit Kosten:

| Option | Kosten | Wann sinnvoll |
|---|---|---|
| Pre-commit-Hook auf der Autoren-Maschine | niedrig | wenn nur lokale Disziplin gefragt ist |
| Make-Target im `make gates`-Block | mittel | wenn auch CI prüfen soll — Standardweg |
| Doku-Konsistenz-Agent (Modul 15) | hoch | wenn semantische Prüfung nötig ist (z. B. "Floskel-Erkennung") |

Der Worked Example wählt **Make-Target + optional Doku-Konsistenz-Agent**:
deterministische Sätze deterministisch, semantische Floskel-Erkennung
inferentiell.

**Schritt 3 — Skript schreiben (Python-Beispiel, kein neues Framework).**

```python
# tools/check_closure_notes.py
import re
import sys
import pathlib
import yaml

DONE = pathlib.Path("docs/plan/planning/done")
FLOSKELN = {"see pr", "n/a", "siehe ticket", "wird nachgereicht"}

def sentences(text: str) -> list[str]:
    no_code = re.sub(r"`[^`]+`|```.*?```", "", text, flags=re.S)
    parts = re.split(r"[.!?]+", no_code)
    return [p.strip() for p in parts if p.strip()]

def errors_for(path: pathlib.Path) -> list[str]:
    front, _, _ = path.read_text().partition("---")[2].partition("---")
    note = (yaml.safe_load(front) or {}).get("closure_note", "")
    if not note:
        return [f"{path}: closure_note fehlt"]
    sents = sentences(note)
    if len(sents) < 2:
        return [f"{path}: closure_note hat nur {len(sents)} Satz"]
    if any(s.lower() in FLOSKELN for s in sents):
        return [f"{path}: closure_note enthält Floskel"]
    return []

errs = [e for p in DONE.glob("*.md") for e in errors_for(p)]
for e in errs:
    print(e)
sys.exit(1 if errs else 0)
```

Sieben Funktionszeilen, drei Fehlertypen. Keine externe Abhängigkeit
außer `pyyaml`.

**Schritt 4 — Als Gate verdrahten:**

```makefile
verify-closure-notes:  ## ADR-0011 — Closure-Note-Pflicht
	python tools/check_closure_notes.py

verify: verify-closure-notes
```

ID-Kommentar zeigt die ADR. *Verify* hängt das Sub-Target ein —
damit greift es genau dort, wo Verifikation läuft, nicht in `make gates`
(weil es keine Code-Architektur-Frage ist, sondern eine DoD-/Closure-
Frage).

**Schritt 5 — Floskel-Erkennung mit inferentieller Schicht ergänzen.**
Floskeln wie *"war ganz okay, läuft jetzt"* sind syntaktisch zwei Sätze
und entgehen Schritt 3. Hier kommt der Doku-Konsistenz-Agent (Modul 15)
ins Spiel:

> Prompt-Anker (in `.harness/skills/closure-note-reviewer.md`):
> "Lies die `closure_note` jedes Slice in `done/`. Markiere alle, die
> *keinen* der folgenden Inhalte tragen: (a) ein konkretes Lernsignal
> (z. B. "Test rot, weil X"), (b) ein konkretes Folge-Slice, (c) eine
> konkrete Architektur-Beobachtung. Floskeln ohne Inhalt sind ein
> HIGH-Finding."

Inferentiell, weil "Inhalt vs. Floskel" semantisch ist; computational
deckt nur die Struktur.

**Schritt 6 — Bewusstes Brechen.** Ein Slice landet in `done/` mit
`closure_note: "Fertig."`. `make verify-closure-notes` läuft rot mit
`docs/plan/planning/done/SL-024.md: closure_note hat nur 1 Satz`. Der
Verifier hat *genau das* erkannt, was Tests nicht erkannt hätten und
Reviewer übersehen würde (Reviewer prüft Diff gegen Plan/ADR — der
fehlende Closure-Eintrag ist *kein* Diff-Symptom).

**Schritt 7 — Pre-completion Checklist-Bezug.** Der Implementation-Agent
soll vor "fertig"-Meldung *selbst* `make verify-closure-notes` laufen
lassen. AGENTS.md-Eintrag:

```markdown
## Closure (Modul 5, ADR-0011)
- Vor `done/`-Verschiebung: `make verify-closure-notes` muss grün sein.
- Floskeln vermeiden — der Doku-Konsistenz-Agent prüft Inhalte.
```

Damit liegt die Hard Rule in zwei Quadranten: *inferential feedforward*
(AGENTS.md sagt es) + *computational feedback* (Make-Target prüft es).

Sieben Schritte, eine Fitness Function für eine ADR-Aussage, die kein
Standard-Tool prüft. Vergleich im Lab:
[`../../lab/example/docs/plan/adr/0011-closure-note-pflicht.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/docs/plan/adr/0011-closure-note-pflicht.md),
[`../../lab/example/tools/check_closure_notes.py`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/tools/check_closure_notes.py)
und das `verify-closure-notes`-Target im
[`../../lab/example/Makefile`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/Makefile).
Das Lab wählt bewusst **Option C** (Closure-Sektion im Markdown-Body,
siehe ADR-0011 §Verglichene Alternativen) statt des oben gezeigten
Frontmatter-Schemas — beide operationalisieren dieselbe ADR-Aussage,
die Wahl ist Repo-spezifisch.

