# Reviewer-Skill — <Repo-Name>

> **Template-Hinweis.** Vorlage für den allgemeinen Code-/Plan-/Design-Reviewer
> (Modul 10, Worked Example „eine Reviewer-Skill-Datei schreiben"). Kopiere nach
> `.harness/skills/reviewer.md`, ersetze `<Platzhalter>` und lösche diesen Block.
> Ein Reviewer ohne Skill-Datei driftet zwischen Sessions (gleiche Eingabe →
> andere Findings/Kategorien); diese Datei ist das repo-spezifische „worauf
> achtest du". Für die *engere* Closure-Note-Prüfung gibt es den
> Schwester-Skill `closure-note-reviewer.md` (Modul 11). Report-Gerüst pro Lauf:
> `docs/reviews/review-report.template.md`.

* Status: Accepted
* Bezug: <ADR-NNNN, AGENTS.md §"Review-Regeln"> · <!-- d-check:ignore (Kurs-/ADR-Referenzen; Anker gelten im Ziel-Repo) -->
* Gilt für: `<agent-review-Make-Target>`

## Kontext-Eingang (Pflicht)

Was der Reviewer *immer* mitbringt, bevor er den Diff liest:

- Diff des PR
- `spec/lastenheft.md` (für referenzierte `LH-*`-IDs)
- ADRs, deren ID im PR oder in der Commit-Message vorkommt
- `AGENTS.md` §"Hard Rules"
- vorherige Findings am gleichen Modul (letzte ~5 PRs)

Ohne diesen Block sieht der Reviewer den Code, aber nicht *die Verträge, gegen
die er prüft*.

## Klassifikation

Jeder Anker HIGH/MEDIUM/LOW hat eine *konkrete* Liste — nicht generisch. INFO ist
bewusst kurz (Ergänzungs-Kanal, nicht Hauptkanal).

**HIGH** — eines der folgenden:
- ADR-Verstoß (Layer, Tool, Hard Rule)
- Sicherheits-Anti-Pattern (Injection, fehlende Auth-Prüfung)
- Korrektheitsfehler im *kritischen* Pfad (<z. B. Index-Schreiben, Auth>)
- Suppression eines Gates (`#noqa`, `//nolint`, `[SuppressMessage]`) ohne ADR
- **Norm nur im Template-Kommentar** — eine Regel steht im `<!-- -->`-Block
  eines `.template.md` und nirgends sonst. Sie ist beim Adopter weg, sobald er
  die Kommentare entfernt. Kein Gate fängt das (siehe Baseline-Regelwerk
  `grundlagen-harness-dateien.md` §Template-Schichtung)
- <**repo-spezifisch #1** — eine Regel, die ein generischer Skill nicht abdeckt,
  z. B. „git mv + Inhalt = zwei Commits" oder „Accepted-ADRs immutable">
- <**repo-spezifisch #2** — eine zweite solche Regel>

**MEDIUM** — eines der folgenden:
- unklare Fehlerbehandlung am Rand des Spec-Bereichs
- fehlende Negativtests bei neuem öffentlichem Vertrag
- Wiederholung eines Musters, das schon zweimal LOW war

**LOW** — stilistisch unschön ohne semantische Auswirkung, einmalige Tippfehler,
unbenutzte Imports.

**INFO** — Hinweis ohne erwartete Aktion (z. B. „diese Stelle hat ein passendes
ArchUnit-Pendant, das du nicht kennst").

> **Pflicht beim Ausfüllen (Modul 10 §Übungen):** Die HIGH-Liste muss mindestens
> *zwei* repo-spezifische Regeln nennen, die ein generischer Skill nicht abdeckt.
> Ist der Skill ohne sie, ist er noch nicht scharf genug — dann kommt bei einem
> Lauf auf einem realen Diff keines deiner Repo-HIGHs zur Anwendung.

## Was dieser Skill NICHT macht

- Keine Lösungsvorschläge („schreib das so") — Reviewer kategorisiert,
  Implementer entscheidet.
- Kein Refactoring-Vorschlag, der über den Diff hinausgeht.
- Keine Verifikation gegen DoD — das ist Verifier-Aufgabe (Modul 11).
- Keine Validation gegen reale Bedürfnisse — das ist Validator-Aufgabe.

Wenn etwas auffällt, das in diese Kategorien gehört: ein INFO-Finding mit Verweis
auf die zuständige Rolle.

## Output-Schema

Jedes Finding:

- `kategorie`: HIGH | MEDIUM | LOW | INFO
- `quelle`: ADR-ID, `LH-*`-ID, Hard-Rule-Name oder „Maintainability"
- `pfad`: Datei:Zeile
- `befund`: 1–2 Sätze, beobachtbar, ohne Lösungsvorschlag
- `verifizierbar`: ja/nein — gibt es einen Gate-Lauf, der es bestätigen würde?

Zusätzlich am Ende: eine Zeile „geprüft, ohne Befund" pro betrachtetem
Verzeichnis (Negativbefund-Zeile — sonst ist „keine Findings" nicht von „nicht
geprüft" unterscheidbar). Report-Gerüst für den ganzen Lauf:
`docs/reviews/review-report.template.md`, ein Report pro Lauf, Folgeläufe als
neue Datei statt Überschreibung.

## Pflege (Steering-Loop)

Bei dreimaligem Auftreten desselben Findings:

- ist die Kategorie noch richtig? → Klassifikation schärfen
- gibt es einen ADR/`AGENTS.md`-Eintrag, der das verhindert hätte?
  → Folge-ADR oder `AGENTS.md`-Update
- gibt es eine Fitness Function, die das prüfen würde? → Modul 13, Gate hinzufügen

Diese Skill-Datei wird **nicht** überschrieben, sondern versioniert
(ADR-Hard-Rule, Modul 4).
