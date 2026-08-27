# Closure-Note-Reviewer-Skill — d-check

**Version:** 1.0.0 · **Datum:** 2026-08-09 ·
**Baseline:** `modul-11-verification.md` §Fitness Function ohne Standard-Tool
(inferentielle Schicht über dem computational Gate) · Ziel-Form
`closure-note-reviewer.template.md` · Schwester-Skill zu
[`reviewer.md`](reviewer.md) (Modul 10).

Dies ist die **semantische** Schicht über dem **strukturellen** Gate
`make verify-closure-notes`
([`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in),
[ADR-0048](../../docs/plan/adr/0048-closure-note-struktur-im-planning-modul.md)).
Das Gate prüft **Struktur**, dieser Skill prüft **Inhalt vs. Floskel**. Der   <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-11-verification.md:64-64 -->
Grund für die Zweiteilung ist konkret: „war ganz okay, läuft jetzt" ist
syntaktisch ein vollständiger Satz und überlebt jede Zählung.

## Eingangs-Kontext (Pflicht — sonst nicht reproduzierbar)

- die Closure-Notiz-Abschnitte der Slices in `docs/plan/planning/done/`
- der `**Lifecycle:**`-/Closure-Abschnitt der Slice-Vorlage — *welche* Inhalte
  eine Notiz tragen soll
- die Ausgabe von `make verify-closure-notes` **für denselben Stand**: was das
  Struktur-Gate bereits gemeldet hat, wird **nicht** doppelt gemeldet
- das Beobachtungs-Register [`observations.md`](../../docs/plan/planning/observations.md)
  — eine dort schon geführte Beobachtung wird zitiert, nicht neu benannt

Ohne den Gate-Lauf im Eingang prüfst du Text, aber du weißt nicht, welcher
Befund schon maschinell abgedeckt ist — und produzierst Doppelmeldungen.

## Prüf-Auftrag

Lies die Closure-Notiz jedes Slice in `done/`. Markiere jede, die **keinen** der
folgenden Inhalte trägt:

1. ein konkretes **Lernsignal** — mit Ursache, nicht nur Behauptung
   („Test rot, **weil** die Grenze am Fence-Ende fehlte")
2. ein konkretes **Folge-Slice** oder eine benannte Folge-Entscheidung
3. eine konkrete **Architektur-/Design-Beobachtung**

Eine Notiz, die nur referiert, *was* getan wurde, erfüllt keinen der drei
Punkte: das steht schon im Diff.

## Kategorien (repo-spezifische Anker)

- **HIGH** — Floskel ohne Substanz: die Notiz überlebt das Struktur-Gate
  (Abschnitt da, Sätze gezählt, keine deklarierte Phrase), trägt aber *keinen*
  der drei Pflicht-Inhalte. Sie ist damit als Audit-Record wertlos.
- **MEDIUM** — genau *einer* der drei Inhalte fehlt oder ist unkonkret:
  Lernsignal ohne das „weil X"; Folge-Slice benannt, aber ohne Eintrag in
  `open/`; Architektur-Beobachtung als Etikett statt als beobachtbare Aussage.
- **LOW** — alle drei Inhalte da, aber schwer nachvollziehbar formuliert
  (Substanz vorhanden, Klarheit fehlt).
- **INFO** — Hinweis ohne erwartete Aktion, z. B. ein Folge-Slice, das noch
  nicht angelegt ist und in die Planung gehört.

**Eskalation:** dieselbe Floskel-Art zum dritten Mal ist kein Einzelbefund mehr,
sondern ein Steering-Loop-Signal — siehe §Pflege.

## Anti-Pattern — was du nicht bist

- **Kein Struktur-Prüfer.** Abschnitt vorhanden? Genug Sätze? Deklarierte
  Floskel? Das ist `make verify-closure-notes`. Doppelmeldung ist ein Fehler,
  kein Gründlichkeitsbeweis.
- **Kein Verifier.** Ob der Slice *fachlich* korrekt abgeschlossen wurde, prüft
  die Verifikation gegen DoD und Spec — nicht du.
- **Kein Ghostwriter.** Du kategorisierst; der Autor formuliert nach. Kein
  Formulierungs-Vorschlag im Befund.
- **Kein Prüfer offener Slices.** `open/`, `next/` und `in-progress/` tragen
  noch keine Closure-Pflicht; ihr `_Ausstehend._` ist korrekt, nicht dünn.
- **Kein Finding ohne Failure-Szenario.** „Wer diese Notiz in einem Jahr liest,
  erfährt X nicht" ist ein Szenario; „wirkt dünn" ist keins.

## Output-Schema (pro Finding)

`kategorie` (HIGH/MEDIUM/LOW/INFO) · `quelle` (`Closure-Inhalt (a)`, `(b)` oder
`(c)` — welcher der drei Pflicht-Inhalte fehlt) · `pfad`
(`docs/plan/planning/done/<slice>.md:<Zeile>`) · `befund` (1–2 Sätze,
beobachtbar, ohne Formulierungs-Vorschlag) · `verifizierbar` (**nein** —
Floskel-Erkennung ist inferentiell; das Struktur-Gate bestätigt nur die Form) ·
`klasse` (stabile Kurz-Bezeichnung der Floskel-Art, über Reviews hinweg
wiederauffindbar).

## Negativbefunde (Pflicht)

<!-- d-check:cite .harness/baseline/v5.12.0/templates/.harness/skills/closure-note-reviewer.template.md:83-83 -->
Eine Zeile „geprüft, ohne Befund: `done/<Charge>`" pro betrachteter
Slice-Charge — sonst ist „keine Findings" nicht von „nicht geprüft"
unterscheidbar.

## Pflege (Steering-Loop)

Bei **dreimaligem** HIGH derselben Floskel-Art:

1. Prüfe, ob die Phrase **strukturell** fangbar ist. Wenn ja, gehört sie in
   `planning.closure.boilerplate` des Closure-Profils
   ([`.d-check.closure.yml`](../../.d-check.closure.yml)) — ein computational
   Marker ist billiger als inferentielles Nachlesen. **Vorher messen:** die
   Phrase darf im Bestand keine Treffer haben, sonst erzeugt sie
   Falschbefunde in Notizen, die sehr wohl tragen.
2. Wenn nein (die Art ist semantisch, nicht lexikalisch): Eintrag im
   Beobachtungs-Register statt stiller Wiederholung.

## Ablage

Ein Report pro Lauf unter `docs/reviews/<YYYY-MM-DD>-<gegenstand>.md`, mit den
Kopf-Metadaten aus [`reviewer.md`](reviewer.md) §Ablage — **Review-Art** hier
stets `Closure-Note`. Nie überschreiben; Folgeläufe bekommen eine neue Datei.
