# `make verify-closure-notes` — hält die Struktur des `done/`-Bestands am Closure-Bindepunkt

## Vertrag

Kein Slice liegt in `done/`, dessen Closure-Notiz fehlt, zu dünn ist, aus einer
deklarierten Floskel besteht oder einen unausgefüllten Vorlagen-Platzhalter
trägt; keiner trägt einen offenen DoD-Haken; kein notiertes Risiko steht ohne
einen der drei Kanon-Ausgänge da; keine zitierte Beobachtung ist unregistriert;
keine Review-Zusage steht ohne Report.

Drei Module tragen das über **ein eigenes Profil**
([`.d-check.closure.yml`](../../.d-check.closure.yml), per `--config`):

- **`planning`** — der erste `closure.heading-pattern`-Abschnitt je Slice:
  vorhanden (`closure-note-missing`), genug Satzende-Zeichen außerhalb der
  Fenced Blocks (`closure-note-thin`), keine deklarierte Floskel
  (`closure-note-boilerplate`), kein Vorlagen-Platzhalter
  (`closure-note-placeholder`, opt-in); dazu `observations` — eine zitierte,
  nicht registrierte Beobachtung (`observation-unregistered`).
- **`structure`** — die Abschnitts-Invarianten (`section-*`). Darunter die
  **urteilsfreie Hälfte der Drei-Ausgänge-Regel** (Baseline-Regelwerk
  `modul-05` §Offene Risiken werden bei Closure aufgelöst) als
  `forbid-pattern` über **jeden** H1-Abschnitt (`sections: each` — ein Selektor
  auf die Titelzeile ließe eine zweite H1 die Spanne still kappen), und seit
  slice-172 der DoD-Haken (`max-open-tasks: 0` ⇒ `section-tasks-open`, je Haken
  auf **seiner** Zeile, mit verfasstem Reparatur-Hinweis;
  [`MR-056`](../conventions/MR-056-dod-haken-waechter.md)).
- **`spans`** — nicht die Closure-Aussage, sondern die **Lesbarkeit des
  Textes, über den die beiden anderen urteilen**: `fence-unclosed`,
  `span-unclosed`, `span-nested-link`.
- **`reviews`** — eine Review-Zusage ohne Report (`review-missing`).

Warum ein eigenes Profil: Der closure-Block darf nicht in der konventionellen
Config wohnen, sonst liefe er über `make planning-check` im inneren Loop mit
und meldete beim Arbeiten an einem laufenden Slice.

## Grenze — was das Grün nicht abdeckt

Zugesagt ist **Struktur, nicht Bedeutung**. Ein grüner Lauf sagt „Form
erfüllt", nicht „die Notizen sind gut". Die semantische Schicht ist der Skill
[`closure-note-reviewer.md`](../../.harness/skills/closure-note-reviewer.md).

1. **Ob ein eingetragener Ausgang inhaltlich trägt** — permanent. Geprüft ist,
   *dass* einer dasteht und *welcher der drei* es ist; ob das Risiko wirklich
   nicht mehr eintreten kann, bleibt Urteil.
2. **Ob ein gesetzter Haken einen Review belegt** — permanent. Der Haken ist
   eine **Selbstauskunft**; die Regel verschiebt die Lücke von *unsichtbar*
   nach *behauptet*.
3. **Ein Haken innerhalb eines wohlgeformten Fenced-Blocks** ist unsichtbar,
   und dort meldet auch `fence-unclosed` nichts — dieselbe Fence-Treue, die
   eine Illustration schützt, ist der Weg, einen offenen Haken zu verstecken.
   Permanent, gewollt.
4. **Ein Platzhalter zwischen zwei Backticks desselben Absatzes** — Inline-Code
   wird **absatzweise** gepaart; der ausgewiesene Preis der Ablösung
   ([ADR-0059](../../docs/plan/adr/0059-closure-waechter-weicht-structure-regel.md),
   gemessen), wo das abgelöste Skript ihn meldete.
5. **Die Platzhalter-Alternation ist eine Liste** — ändert die Vorlage ihre
   Form, schweigt sie. Heilbar: beim Vorlagen-Bump mitprüfen.
6. **Ein Risiko ganz ohne Ausgangs-Marker** sieht sie nicht; der Altbestand vor
   slice-140 ist ausgenommen, ebenso der DoD-Haken-Bestand bis `slice-170`.

**Von diesen Grenzen gelten dem DoD-Haken-Wächter die Fence-Hälften, nicht die
Inline-Code-Hälften:** er liest die **rohen** Zeilen, also verschluckt die
absatzweise Paarung bei ihm nichts. Seine eigene Grenze ist der **vergessene**
Schluss-Fence — er macht ihn blind, und genau dort meldet `spans`.

**Wie groß der Ausschnitt ist, sagt das Kommando, nicht diese Datei:**
`make verify-closure-notes` nennt die geprüfte Dateizahl in seiner
Schluss-Zeile. Eine eingefrorene Zahl stünde hier falsch, sobald jemand
committet. Die Zahl ist eine Aussage über **diesen Ausschnitt**, nicht über das
Repo: Der Bindepunkt scannt eine **Teilmenge** dessen, was `make doc-check`
sieht — `make doc-check` nennt seine eigene Zahl daneben.

**Was `spans` hier hinzufügt, ist Unabhängigkeit, nicht Deckung.** Gemessen war
der Zuwachs an gefundenen Defekten **null** — `make doc-check` meldet denselben
Befund schon beim Commit. Gekauft ist, dass die Zusage dieses Profils nicht
mehr daran hängt, dass ein **anderes** Profil `spans` führt
([ADR-0077](../../docs/plan/adr/0077-spans-am-bindepunkt-die-begruendung-traegt-anders.md),
supersedes [ADR-0076](../../docs/plan/adr/0076-spans-am-closure-bindepunkt.md)).
Die vierte Grenze oben deckt es **nicht** — ein wohlgeformter Span ist kein
Defekt, sondern Code.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | keine Befunde — Form erfüllt |
| 1 | mindestens ein Befund; die Grund-Codes stehen je Zeile |
| 2 | fail-closed: Prüfverzeichnis fehlt oder Profil unlesbar |

Diagnose-only: Der Lauf schreibt nichts ins Repo.

## Sperren

- `fehlendes done/-Verzeichnis` — fail-closed statt stiller Leermenge: eine
  Zusage über null Dateien wäre wahr, ohne etwas gesehen zu haben.

## Bindung

**Closure-Bindepunkt** — in `make fullbuild`, bewusst **nicht** in
`gates`/`ci`. Zusätzlich am **Übergang** selbst: der `pre-commit`-Hook löst ihn
aus, sobald ein Rename/Add nach `docs/plan/planning/done/slice-*.md` gestagt
ist (nicht rekursiv), und dieselbe Bindung läuft in der PR-/Push-CI über die
Commit-Range.

[ADR-0048](../../docs/plan/adr/0048-closure-note-struktur-im-planning-modul.md) ·
[ADR-0059](../../docs/plan/adr/0059-closure-waechter-weicht-structure-regel.md) ·
[ADR-0077](../../docs/plan/adr/0077-spans-am-bindepunkt-die-begruendung-traegt-anders.md) ·
[ADR-0082](../../docs/plan/adr/0082-uebergangswaechter-reviews-observations.md) ·
[`MR-049`](../conventions/MR-049-ausgangs-wortschatz.md) ·
[`MR-056`](../conventions/MR-056-dod-haken-waechter.md) ·
[`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) ·
[`DC-FA-CLI-012`](../../spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben) ·
[`DC-FA-SPAN-001`](../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in) ·
[`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in)
