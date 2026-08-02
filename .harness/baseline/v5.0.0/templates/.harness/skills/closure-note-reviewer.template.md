# Closure-Note-Reviewer-Skill — <Repo-Name>

> **Template-Hinweis.** Vorlage für den *inferentiellen* Closure-Note-Reviewer
> (Modul 11 §Schritt 5, ADR-0011-Folgepflicht, Modul 15 Doku-Konsistenz-Agent).
> Kopiere nach `.harness/skills/closure-note-reviewer.md`, ersetze `<Platzhalter>`
> und lösche diesen Block. Er ist die *semantische* Schicht **über** dem
> *computational* Gate `make verify-closure-notes`
> (`tools/check_closure_notes.py`): das Gate prüft Struktur (Heading + Satzzahl),
> dieser Skill prüft *Inhalt vs. Floskel*. Sechs-Schritt-Muster wie der
> allgemeine `reviewer.md` (Modul 10), aber eng auf Closure-Notes fokussiert.

* Status: Accepted
* Bezug: ADR-0011 (Closure-Note-Pflicht), `tools/check_closure_notes.py`,
  Modul 11 §Schritt 5, Modul 15 (Doku-Konsistenz-Agent) · <!-- d-check:ignore (Kurs-/ADR-Referenzen; Pfade gelten im Ziel-Repo) -->
* Gilt für: den *inferentiellen* Nachlauf zu `make verify-closure-notes` —
  greift dort, wo Struktur allein die Floskel nicht fängt

## Kontext-Eingang (Pflicht)

Was der Reviewer *immer* mitbringt, bevor er urteilt:

- alle `closure_note`-Abschnitte der Slices in `docs/plan/planning/done/`
- das Slice-Template `docs/plan/planning/slice.template.md` §"Closure-Notiz"
  (welche drei Inhalte Pflicht sind)
- ADR-0011 — *warum* Closure-Notes existieren (Auditierbarkeit, Lernsignal)
- das Ergebnis von `make verify-closure-notes` für denselben Stand — was das
  Struktur-Gate bereits abgedeckt hat, wird **nicht** doppelt gemeldet

Ohne diesen Block prüft der Reviewer Text, aber nicht *gegen die drei
Pflicht-Inhalte, die eine Closure-Note tragen muss*.

## Prüf-Auftrag (wörtlich aus Modul 11 §Schritt 5)

> "Lies die `closure_note` jedes Slice in `done/`. Markiere alle, die *keinen*
> der folgenden Inhalte tragen: (a) ein konkretes Lernsignal (z. B. »Test rot,
> weil X«), (b) ein konkretes Folge-Slice, (c) eine konkrete
> Architektur-Beobachtung. Floskeln ohne Inhalt sind ein HIGH-Finding."

Inferentiell, weil „Inhalt vs. Floskel" semantisch ist; das computational Gate
deckt nur die Struktur (Heading, Satzzahl außerhalb Code-Blöcken, Floskel-Liste).

## Klassifikation

**HIGH** — Floskel ohne Substanz: die `closure_note` ist syntaktisch vorhanden
(überlebt `check_closure_notes.py`), trägt aber *keinen* der drei Pflicht-Inhalte.
Beispiele: „war ganz okay, läuft jetzt", „Fertig.", „wie geplant umgesetzt".

**MEDIUM** — genau *einer* der drei Pflicht-Inhalte fehlt oder ist unkonkret:
- Lernsignal ohne das „weil X" (Behauptung statt Ursache)
- Folge-Slice benannt, aber ohne zugehörigen `open/`-Eintrag
- Architektur-Beobachtung als Etikett statt als beobachtbare Aussage

**LOW** — alle drei Inhalte vorhanden, aber schwer nachvollziehbar formuliert
(Substanz da, Klarheit fehlt).

**INFO** — Hinweis ohne erwartete Aktion (z. B. „verweist auf ein Folge-Slice,
das noch nicht in `open/` liegt — Tracker-Nachtrag durch die Planning-Rolle").

## Was dieser Skill NICHT macht

- Keine Struktur-Prüfung (Heading vorhanden? ≥ 2 Sätze?) — das ist
  `tools/check_closure_notes.py`; nicht doppeln.
- Keine Bewertung, ob der Slice *fachlich* korrekt abgeschlossen wurde —
  Verifier/Validator (Modul 11).
- Keine Umschreibung der `closure_note` — der Autor formuliert nach, der
  Reviewer kategorisiert nur.
- Keine Prüfung von Slices außerhalb `done/` — `open/`/in-progress tragen noch
  keine Pflicht-Closure.

## Output-Schema

Jedes Finding:

- `kategorie`: HIGH | MEDIUM | LOW | INFO
- `quelle`: `ADR-0011` | `Closure-Inhaltspflicht (a/b/c)`
- `pfad`: `docs/plan/planning/done/<slice>.md`:<Zeile>
- `befund`: *welcher* der drei Pflicht-Inhalte fehlt, 1–2 Sätze, beobachtbar,
  ohne Formulierungs-Vorschlag
- `verifizierbar`: nein — Floskel-Erkennung ist inferentiell;
  `check_closure_notes.py` bestätigt nur die Struktur, nicht den Inhalt

Zusätzlich am Ende: eine Zeile „geprüft, ohne Befund: `done/<Charge>`" pro
betrachteter Slice-Charge (Negativbefund — macht die Abdeckung sichtbar).

## Pflege (Steering-Loop)

Bei dreimaligem HIGH derselben Floskel-Art:

- Muster in ADR-0011 bzw. `AGENTS.md` §"Closure" als benanntes Anti-Pattern
  aufnehmen
- prüfen, ob `check_closure_notes.py` die Floskel *strukturell* fangen könnte
  (Floskel-Liste erweitern) → dann ins computational Gate heben; ein
  computational Feedforward-Marker ist billiger als inferentielles Nachlesen
- Slice-Template §"Closure-Notiz" schärfen, falls das Feld die drei
  Pflicht-Inhalte nicht klar genug abfragt

Diese Skill-Datei wird **nicht** überschrieben, sondern versioniert
(ADR-Hard-Rule, Modul 4).
