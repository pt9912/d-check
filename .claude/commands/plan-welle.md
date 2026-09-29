# Welle planen (Harness)

Argument: $ARGUMENTS

Dieser Command führt die **Planner**-Rolle für *eine* Welle (Modul 6 — Roadmap Engineering). Eine
Welle ist ein **Bündel von Slices**, das gemeinsam geplant und geschlossen wird. Der Welle-Status ist
die **Verzeichnis-Position**, **kein `Status:`-Feld**: eine **offene** Welle liegt **flach** in
`docs/plan/planning/<welle-id>.md`, bei Closure wandert sie per `git mv` nach `done/` (das schließt
`/close-welle`). Der Roadmap-Abschnitt *Offene Wellen* ist derivativ — ein Zeiger je flacher
Welle-Datei; woran gerade gearbeitet wird, sagt das `Welle:`-Feld der Slices in `in-progress/`.

Kanonische Quellen (vendored Regelwerk, `.harness/baseline/<tag>/regelwerk/`): Modul 6 (Roadmap),
Modul 5 (Planning-Lifecycle), Modul 7 (Carveouts). Bei Konflikt gilt der Kurs.

## Repo-lokale Adaptionen, die du beachten MUSST (d-check)

Lies den MR-Block in `harness/conventions.md`; die planungs-relevanten:

- **Neue Artefakte per `cp` aus den vendored Templates** (`.harness/baseline/<tag>/templates/…`),
  dann **in-place** ausfüllen — **keine handgeschriebenen Kopien und kein Modellieren auf ein
  bestehendes Artefakt.** Das gilt für den Welle-Plan (`welle.template.md`) **und** jeden neuen
  Slice (`slice.template.md`).
- **Strenges Doc-Gate (d-check).** Jede `DC-`/`ADR-`/`MR-`/`slice-`-Kennung in einer gescannten
  `.md` muss ein klickbarer Anker-Link sein (link-policy: always) — ein bares Kennungs-Token bricht
  `doc-check`. Der Welle-Plan wird gescannt. **`welle-NN` ist keine Traceability-Klasse** — der
  commit-msg-Hook verlangt zusätzlich eine slice-/MR-/ADR-/DC-ID in der Botschaft.
- **Commit via Message-Datei** (`git commit -F <datei>`); `make gates` endet mit `record-gates`,
  jede Inhaltsänderung nach einem Gate-Lauf macht den Stempel ungültig.
- **§8 der Slice-Pläne trägt die beiden `d-check:cite`-Spannen** ([`MR-054`](../../harness/conventions.md#mr-054)),
  der Plan-Kopf die drei Vorprüfungen samt Nachtlauf ([`MR-053`](../../harness/conventions.md#mr-053)).

## Kontext lesen

1. `CLAUDE.md` (falls vorhanden), `harness/README.md`, `AGENTS.md`, `harness/conventions.md` lesen.
2. Den Regelwerk-Index (`.harness/baseline/<tag>/regelwerk/README.md`) und **Modul 6** on-demand lesen
   (Source Precedence). Nicht den ganzen Baum laden.
3. Die Roadmap (`docs/plan/planning/in-progress/roadmap.md`) lesen: steht die zu planende Welle schon
   als Zeiger in *Offene Wellen* (dann existiert ihre Datei) oder als Zeile in *Nächste Wellen*
   (dann wird sie jetzt eröffnet)? Liegen ihre Slices bereits in `open/`?

## Eröffnung: die drei Pflichtteile (Modul 6 — vor dem Schreiben benennen)

**Vor den Pflichtteilen: das Beobachtungs-Register sichten** (`docs/plan/planning/observations/README.md`,
Modul 6, Eröffnungs-Schritt 2). Betrifft eine offene Beobachtung eine Sub-Area dieser Welle, gehört
sie in die Slice-Planung (Risiko im betroffenen Slice) — und erreicht der Eintrag **mit dieser Welle**
3× (die Zahl der Dateien unter seinem `evidence/`; ein Zähler-Feld gibt es nicht), ist er keine Notiz
mehr, sondern eine Lücke und braucht einen eigenen Slice. **Bei der ersten Welle entfällt der Schritt
nicht:** das Register existiert ab Repo-Beginn; trägt die Ablage nur ihre `README.md`, ist genau das
die Antwort und wird notiert.

4. Eine Welle braucht **minimal drei Bestandteile**, sonst ist sie keine Welle:
   - **Slice-IDs** (der Inhalt) — welche Slices bündelt sie?
   - **Trigger** (Welle startet) — eine **beobachtbare Bedingung**, kein Datum. Beobachtbar heißt:
     *ein anderer Mensch kann ohne Rückfrage sagen, ob er eingetreten ist* (z. B. „Welle X done",
     „Replay grün"). „Sobald wir Zeit haben" scheitert daran.
   - **Closure-Kriterien** (Welle schließt) — Aktion, kein Termin (z. B. alle Slices in `done/`,
     `make fullbuild` grün, das Trigger-Audit vollzogen). Ein Datum darf als *Schätzung* erscheinen,
     triggert nie.
5. Berichten: Welle-ID · Zielmeilenstein · Slice-IDs · Trigger · Closure-Kriterien.

## Slices bereitstellen

6. Existiert ein Slice der Welle noch nicht, ihn **per `cp` aus `slice.template.md`** anlegen
   (`docs/plan/planning/open/slice-<NN>-<titel>.md`), dann füllen. Nie hand-authoren. **§8 des
   Plans trägt dieselbe Register-Sichtung noch einmal je Slice** (Modul 5, *Zwei Schritte vor der
   Modus-Begründung*) — und ist damit für alles **unter** 3× der einzige Leser; keine Treffer sind
   dort ebenfalls eine Antwort und werden notiert.

## Welle-Plan per cp anlegen und füllen (der Kern-Schritt)

7. **`cp` aus `.harness/baseline/<tag>/templates/docs/plan/planning/welle.template.md` nach
   `docs/plan/planning/<welle-id>.md`** — flach in `planning/`, kein Lifecycle-Ordner. Provenienz mit
   `diff -q <template> <ziel>` belegen (byte-identisch, dann füllen).
8. **In-place füllen** (Edits, **kein** Voll-Überschreiben): den `> **Template-Hinweis.**`-Block
   strippen, alle Platzhalter ersetzen, die `<!-- -->`-Guidance-Kommentare entfernen. Die
   **`Lifecycle:`-Note der Vorlage behalten** (Zustand = Verzeichnis-Position, **kein
   `Status:`-Feld**), Zielmeilenstein und Verantwortlich/Datum setzen. Die Abschnitte mit den drei
   Pflichtteilen aus Schritt 4 füllen, dazu **§6 Out-of-Scope** — was nicht ausdrücklich
   ausgeschlossen ist, dehnt die Welle, bis der Closure-Trigger unerreichbar wird. Kennungen als
   Anker-Links.

## Roadmap verdrahten und gaten

9. Roadmap fortschreiben: mit dem Anlegen der flachen Welle-Datei verlässt die Welle-Zeile *Nächste
   Wellen*, und unter *Offene Wellen* steht der Zeiger auf die Datei. Ist die Welle nur geplant,
   bleibt ihre Zeile in *Nächste Wellen* (Trigger · wichtigste Slices · Aufwand S/M/L, kein Termin).
   Die Welle-Verweise der zugehörigen Slices auf die neue Plan-Datei ziehen.
10. `make gates` laufen lassen (grün). Beim **Anlegen** ist der Welle-Plan **Inhalt** → ein einzelner
    Commit, hier **kein `git mv`** (die neue Welle entsteht flach = aktiv/geplant). Der `git mv` nach
    `done/` kommt erst bei der **Closure** (`/close-welle`). Commit via `-F`.

**Merke (Modul 6):** Eine Welle endet durch **Closure-Kriterien**, nicht durch ein Datum
(Welle ≠ Sprint). Ein Trigger ist eine beobachtbare Bedingung, kein Kalendertag. Die fertige Welle
schließt `/close-welle`.

Gates nicht überspringen. Keine Erfolgsmeldung ohne Command-Ausgabe.
