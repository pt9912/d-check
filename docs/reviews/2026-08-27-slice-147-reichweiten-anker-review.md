# Review — slice-147 / Commit `7869b0d`

**Review-Art:** Plan + Doku (gegen `AGENTS.md` §1/§5, Reviewer-Skill, Baseline
`grundlagen-source-precedence.md` / `modul-09` / `modul-06`) ·
**Gegenstand:** `docs/plan/planning/in-progress/slice-147-reviewer-anker-reichweite.md`,
Commit `7869b0d` (AGENTS.md §5, `.harness/skills/reviewer.md` 1.10.0→1.11.0,
`observations.md` BEO-012) · **Skill:** `reviewer.md` @ 1.11.0 — dieselbe
Fassung, die dieser Commit einführt (dieselbe Lage wie bei slice-131; hier
**nicht** tragbar, weil sich eine Kategorie und eine Prüffrage bewegt haben) ·
**Modell:** claude-opus-5[1m] · **Datum:** 2026-08-27

**Selbst gefahren:** `make doc-check` (`d-check: 535 Datei(en) geprüft, 0
Befund(e)`), `git show`/`git log`, Lesekommandos. Arbeitsbaum unverändert
(`git status --porcelain` leer).

Der neue Anker wurde zuerst auf diesen Diff selbst angewandt. Er greift.

## Findings

### F-1 · HIGH · Die tragende Aussage des Slice ist falsch: die Regel stand normativ im Kanon

**Pfad:** Slice `:93-96` und `:130-133`; Commit-Botschaft `7869b0d`.

Der Kanon trägt die Regel wörtlich, seit dem gepinnten Stand —
`grundlagen-source-precedence.md:145-149`:

> **Wie weit trägt ein zitierter Satz?** … Sie ist an **jede** zitierte Aussage
> zu stellen, auch an einen Satz der Baseline: *Gilt er auch außerhalb des
> Falls, für den er geschrieben wurde?*

`docs/plan/planning/done/slice-149-baseline-v5120-delta-audit.md:137` führt ihn
als Delta-Punkt **4 — die Reichweitenfrage als Frage** und weist **slice-147**
ausdrücklich als Träger seiner Feedforward-Hälfte aus. Der Slice hatte die
Fundstelle in einem geschlossenen Artefakt vor sich und schreibt trotzdem
„nirgends". Gemessene Menge: drei Orte (`AGENTS.md`, `harness/README.md`,
Konventionsspeicher — dort steht sie wirklich nicht); Schluss: universal.
**Klasse:** `universale-nicht-existenz-behauptung-aus-drei-messpunkten`.

### F-2 · HIGH · Dieselbe falsche Aussage steht jetzt im lebenden Register

**Pfad:** `observations.md:17`.

Das Register liest jeder Lauf. Die Zeile über die Klasse „eine Quelle wird über
ihren Geltungsbereich hinaus zitiert" behauptet dauerhaft etwas Falsches über
den Bestand — die fünfte Instanz der Klasse entsteht im Eintrag der Klasse.

### F-3 · MEDIUM · Zwei Zeilen über BEO-012 steht die entgegengesetzte Entscheidung, ungemeldet

**Pfad:** `observations.md:12` (BEO-017) gegen `:17`; `AGENTS.md:425-435`.

BEO-017 steht in derselben Lage (Kanon trägt die Regel seit `v5.12.0`) und
zieht den umgekehrten Schluss: *„Darum hier keine eigene Regelzeile: sie zu
duplizieren verstieße gegen `AGENTS.md` §1."* slice-147 legt für seine Klasse
eine eigene Regelzeile an, ohne die Gegenposition zu nennen. Ein tragfähiges
Gegenargument existiert (`modul-09-implementierung.md:170-175`,
Zwei-Quadranten-Pflicht, so von slice-136 angewandt) — es steht nur nirgends im
Diff. Das ist die Gegenrichtung, die der neue Anker in derselben Änderung
einfordert.

### F-4 · MEDIUM · Die neue Hard Rule trägt keinen Zeiger auf ihre Quelle

**Pfad:** `AGENTS.md:425-435`. Die Form für den Kanon-Fall steht drei Absätze
höher: §3.7 schließt mit „Kanon: [Baseline §Was ein Kommentar trägt]". Ohne
Zeiger ist der §5-Absatz eine Kopie ohne Anker — die Drift, gegen die §1
geschrieben ist.

### F-5 · MEDIUM · Der Herkunfts-Anker fehlt

**Pfad:** `AGENTS.md:433-435`. `modul-09-implementierung.md:167-169`: eine Hard
Rule aus dem Steering Loop trägt `(seit welle-<NN>)` — ohne Welle
`(seit slice-<NNN>)`. Der Slice ist wellenlos; der Anker müsste
`seit slice-147` lauten. Beide Nachbarn führen ihn.

### F-6 · MEDIUM · Der Verweis auf slice-131 trägt die Waisen-Aussage nicht

**Pfad:** Slice `:96-99`, `:133-135`. slice-131s Kriterium: buchstabiert der
Anker eine **gerankte Fundstelle** aus (zulässig) oder legt er fest (Waise)?
Sein Zensus akzeptiert die vendorte Baseline mehrfach als gültige Fundstelle.
Mit der Fundstelle aus F-1 wäre der neue Anker ein „buchstabiert aus", keine
Waise.

### F-7 · MEDIUM · Der neue Anker bündelt zwei Prüffragen — das Argument gegen den Zusammenzug

**Pfad:** `.harness/skills/reviewer.md:77-80`. Der Risiko-Ausgang lehnt den
Zusammenzug ab, weil er „zwei Prüffragen unter einer Überschrift" trüge. Der
Anker tut es dann selbst: nach der Reichweiten-Prüffrage folgt „Prüfe zusätzlich
die Gegenrichtung …" — anderer Gegenstand (übersehener Konflikt + Meldepflicht
aus §1), andere Arbeitsanweisung, in §5 und §9 nicht gewogen.

### F-8 · MEDIUM · Der Kontext-Last-Abnahmepunkt ist umgangen, nicht behandelt

**Pfad:** Slice `:77-88`. Gezählt: **16** benannte HIGH-/MEDIUM-Anker nach der
Änderung, 15 davor — 7 HIGH, 2 einzeln stehende MEDIUM, 7 MEDIUM im
Sammelpunkt. Der neue ist der sechzehnte und mit 19 Zeilen der **längste
Einzelanker** der Datei. Der Risikosatz nennt „der fünfzehnte wird nicht mehr
gelesen"; die Schwelle war vorher überschritten. Der Ausgang beantwortet nur die
zweite Hälfte und **zählt nie**.

### F-9 · MEDIUM · DoD-Haken „unabhängiger Review" gesetzt, bevor der Review lief

**Pfad:** Slice `:73`. Der Haken steht im Feature-Commit; der Review lief erst
danach. Wörtlich die Instanz, mit der `BEO-009`(a) auf 5 gezählt wurde.
**Systemisch:** slice-146 hat es zwei Commits vorher genauso gemacht.
Randbeobachtung: für slice-146/151/154–156 existiert **kein** Report unter
`docs/reviews/`, obwohl der Skill §Ablage „ein Report pro Lauf" zusagt.

### F-10 · MEDIUM · Das Register behauptet eine Schließung, die nicht stattgefunden hat

**Pfad:** `observations.md:17`. Der Slice liegt in `in-progress/`, sein §6 sieht
die Rückführung vor. Der Halbsatz „vorher stand die Regel nirgends normativ, nur
hier" ist zudem reine Chronik — und falsch. Die Belege-Spalte ist **korrekt**
(Zähler 4, vier Kennungen, alle in `done/`, Zähler zu Recht nicht erhöht).

### F-11 · LOW · „Beides ist Pflicht" der falschen Kanon-Datei zugeschrieben

**Pfad:** Slice `:42-44`. Der Satz steht in `modul-09-implementierung.md:207`
und meint dort **AGENTS.md + Fitness Function**, nicht einen zweiten
inferentiellen Ort. `grundlagen-klassifikation.md` trägt die 2×2-Matrix
deskriptiv und das Wort „Pflicht" gar nicht. *(Vorbestehend, aus `791b1ec`.)*

### F-12 · LOW · Die `citations`-Messung nennt die kleineren Grenzen, lässt die größte weg

**Pfad:** Slice `:154-156`. Genannt sind 16-Zeichen-Schwelle und „Verweis ohne
Zitat". Die schärfere steht im Spec-Kopfsatz: *„prüft **nur** per Direktive
ausgezeichnete Zitate — kein Prosa- oder Voll-Scan."* Ein wörtliches Zitat
**ohne** Direktive ist der Normalfall im Bestand.

### F-13 · LOW · Hängender Bezug nach der Streichung

**Pfad:** `observations.md:17` — „die Entscheidung **darüber**" bezog sich auf
den gestrichenen Satzteil und zeigt jetzt ins Leere.

### F-14 · LOW · Wortlaut-Drift zwischen Register und Anker

**Pfad:** `.harness/skills/reviewer.md:72-73` gegen `observations.md:17`. Das
Register: eine ADR-Entscheidung, die einen einmaligen **Akt beschreibt**. Der
Anker: die ADR **verbietet** einen Akt. Das verengt auf verbietende ADRs.

## Negativbefunde (geprüft, ohne Befund)

- **`citations`, vier Aussagen gegen Spec und Code** — alle vier tragen:
  zusammenhängender Teilstring der whitespace-normalisierten Quell-Spanne;
  Reichweite ist nicht Gegenstand; „kürzer als 16 Zeichen ⇒ nicht geprüft";
  Verweis ohne Direktive ⇒ „prüft das Modul nichts". Auch „heute nicht
  aktivierbar" ist belegt. Einzige Auslassung: F-12.
- **Zahlen der Botschaft** — `doc-check` reproduzierbar 535/0; `Makefile:206`
  listet genau zehn Glieder in `gates`.
- **Register-Form BEO-012** — Zähler 4, vier `slice-<NNN>`-Belege, alle in
  `done/`; Zähler korrekt **nicht** erhöht.
- **Kategorie-Wahl MEDIUM** — begründet und tragfähig: kein stiller Grün-Pfad,
  kein Gate berührt.
- **Anker als Prüffrage** — erfüllt; die universelle Frage steht da und ist
  typ-unabhängig. *(Anmerkung: kein Nicht-Abschließend-Vorbehalt zur
  Drei-Typen-Aufzählung.)*
- **Zusammenzug mit `BEO-009`(b) abgelehnt** — die Begründung trägt:
  verschiedene Gegenstände, unvereinbare Arbeitsanweisungen. Der Einwand liegt
  woanders (F-7).
- **`MR-045`** eingehalten; Traceability-Kennung vorhanden; §7-Vorprüfungen
  vorhanden; kein Skill-Versions-Spiegel außerhalb der Datei.

## Kategorie-Summary

HIGH 2 · MEDIUM 8 · LOW 4

## Urteil

**Nicht schließbar** — Nacharbeit vor der Closure erforderlich.

Der Anker ist gut gebaut, die Kategorie stimmt, die `citations`-Messung trägt,
die Ablehnung des Zusammenzugs ist sauber begründet. Was nicht trägt, ist die
**Begründung, warum es diese Änderung braucht**. Die Änderung selbst kann
bleiben — `modul-09`s Zwei-Quadranten-Pflicht deckt einen `AGENTS.md`-Eintrag
auch dann, wenn der Kanon die Regel trägt —, aber sie muss auf **dieser**
Begründung stehen, den Kanon-Zeiger führen, den Herkunfts-Anker tragen und den
Widerspruch zu BEO-017 melden.

Der Slice hat die Klasse, gegen die er antritt, in seiner eigenen Begründung
begangen — an der Stelle, an der eine Quelle nicht überdehnt, sondern für nicht
existent erklärt wurde. Das gehört bei der Closure als **fünfte Instanz** an den
BEO-012-Zähler, mit slice-147 als Beleg.
