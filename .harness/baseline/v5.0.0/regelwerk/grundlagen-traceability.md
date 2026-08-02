## Traceability-Constraint
<!-- Quelle: [grundlagen/traceability.md](https://github.com/pt9912/ai-harness-course/blob/v5.0.0/kurs/de/grundlagen/traceability.md) -->

## Traceability-Constraint

Keine relevante Änderung ohne Bezug zu mindestens einem der folgenden Punkte:

* Requirement-ID
* Architektur-ID oder Architekturprinzip
* ADR-ID
* Test, Gate oder Demo-Artefakt
* Dokumentations-Update, falls ein öffentlicher Vertrag betroffen ist

Das ist eine *computational feedforward*-Kontrolle (siehe
[`klassifikation.md`](grundlagen-klassifikation.md#klassifikation-und-steering-loop)): ein Commit-Hook prüft, dass
die Nachricht mindestens eine ID enthält. Billig, deterministisch, und
sie zwingt den Implementer-Agent in die Source-Precedence-Kette zurück.

<a id="herkunfts-anker"></a>

### Herkunfts-Anker für Steering-Loop-Regeln

Der Traceability-Constraint bindet **Änderungen** an eine ID. Der
Herkunfts-Anker ist dieselbe Regel auf dem **Artefakt**: Eine Regel, die
aus dem Steering Loop entstand, nennt die Welle, in der sie entstand — oder,
wenn sie ohne Welle verkörpert wurde, den Slice: `seit welle-<NN>` bzw.
`seit slice-<NNN>`.

- **Geltungsbereich — eng.** Nur Regeln, die die 3×-Schwelle erreicht
  haben. Was aus Lastenheft, Spezifikation oder ADR folgt, trägt bereits
  eine ID und braucht keinen zweiten Anker.
- **Form** — ein Feld, kein Konstrukt:
  `noqa-gate:  ## LH-QA-SUP-002 · seit welle-3` (Make-Target) ·
  `coverage-floor: ## LH-QA-SUP-004 · seit slice-047` (wellenlos) ·
  `### 3.3 <Hard Rule>   (seit welle-3)` (AGENTS.md) ·
  `- <HIGH-Regel>  (seit welle-3)` (Reviewer-Skill). Der Adaptions-Block
  trägt das Muster bereits über sein Feld *Begründung*.
- **Die Welle ist der Regelfall, der Slice die Ausnahme.**
  `done/welle-<NN>-results.md` §Steering-Loop-Einträge nennt beim
  Schwellen-Übertritt *Regel · stabile Bezeichnung · Slice-Belege* — ein Anker
  löst damit in einem Hop auf und bleibt grob genug, um nicht zu verrotten.
  Wurde die Regel **ohne Welle** verkörpert, gibt es diese Datei nicht; dann
  ist der Slice die einzige auflösbare Herkunft (`seit slice-<NNN>`, löst über
  `done/slice-<NNN>-<kurzer-titel>.md` §7 auf — die Nummer ist eindeutig, der
  Titelrest gehört zum Dateinamen; maschinell also `done/slice-<NNN>-*.md`).
- **Ab Einführung, kein Nachrüsten.** Altbestand bleibt ohne Anker;
  `seit unbekannt` wäre eine Harness-Lüge, der leere Zustand ist die
  ehrliche Information.

**Sensor 1 — Anker-Paarung** (*computational feedback*). Die Prüfung läuft
**von der Closure-Notiz nach außen**, nicht von der Regel nach innen: von
der Regel aus ist nicht entscheidbar, ob sie einen Anker braucht.
**Ausgelöst wird durch ein Feld, nicht durch die Semantik des Eintrags und
nicht durch Prosa:** durch das Pflichtfeld **`liegt in <Zielort>`** — in
`## Steering-Loop-Einträge` jeder `welle-<NN>-results.md` und, für wellenlos
verkörperte Regeln, in §7 jeder `done/slice-<NNN>-<kurzer-titel>.md`; die
kanonischen Formen liefern `welle-results.template.md` bzw.
`slice.template.md` §7 (siehe Ziel-Form unten).

- **Das Feld gilt nur in diesen beiden Sektionen.** Überall sonst sind es
  gewöhnliche Wörter und lösen nichts aus — der Trigger-Sprachgebrauch
  „`SL-024` liegt in `done/`" (Modul 6) ebenso wenig wie eine bloße Erwähnung
  eines Pfades im Fließtext. Der Sektions-Scope grenzt den Auslöser ein,
  ersetzt ihn aber nicht: *innerhalb* der Sektion entscheidet das Feld.
- **Die Ruheort-Regel — für jede Datei, die per `git mv` wandert.** Ein
  Slice-Plan und ein Welle-Plan werden an einem Ort geschrieben und an einem
  anderen gelesen: Bei der Closure wandern sie nach `done/`. Jeder relative
  Pfad darin ist deshalb so zu schreiben, wie er **vom Ruheort** auflöst, nicht
  vom Schreibort — die Ergebnis-Notiz liegt in `done/` als Geschwister (ohne
  Präfix), das Beobachtungs-Register eine Ebene höher (Eltern-Verzeichnis, also mit `..`-Präfix).
  Ein im Schreibmoment richtiges `done/…` bricht für jeden Leser danach, und
  zwar still: Der Pfad bleibt syntaktisch intakt und zeigt ins Leere.
- **In den Backticks steht ein Zielort, nicht immer eine Datei** — drei
  kanonische Füllungen: `AGENTS.md §<N>` · `Makefile:<target>` ·
  `.harness/skills/<name>.md`.
- **Geprüft wird:** (1) der Pfad existiert, **ab Repo-Wurzel** — nicht relativ
  zur Closure-Notiz: Der Zielort zeigt aus dem Planungs-Baum hinaus und wandert
  nicht mit, wenn die Notiz nach `done/` wandert. Dafür wird ein Suffix ab
  ` §` oder ab `:` abgetrennt und der Rest als Pfad geprüft. (Die Pfade auf
  Nachbar-Artefakte — der Zeiger aufs Beobachtungs-Register — bleiben
  datei-relativ und folgen der Ruheort-Regel.) (2) Das Ziel trägt `seit welle-<NN>` bzw.
  `seit slice-<NNN>` — beim Make-Target auf dessen Target-Zeile, beim
  Abschnitt in dessen Überschrift, bei einer Datei ohne Suffix irgendwo in ihr.
- **Fehlt das Feld**, ist der Eintrag *gezählt, nicht verkörpert* und kein
  Gegenstand der Paarung. Ausnahme ohne Gegenausnahme: Eine **benannte
  Spec-Lücke** trägt kein `liegt in` und ist trotzdem verkörpert — in einer
  versionierten Spec statt an einem Zielort. Ihr Gegenstück ist die
  `LH-*`-ID; an der Register-Paarung (Modul 6) nimmt sie teil wie jeder
  andere Eintrag.

Rot bei: Regel nie geschrieben · still gelöscht · Anker vergessen —
dieselbe Klasse wie ein halluziniertes Gate
([Modul 13](modul-13-quality-gates.md)).

> **Grenze:** Der Sensor erzwingt den Anker nur für **deklarierte**
> Steering-Loop-Regeln. Wer die Closure-Notiz nicht schreibt, wird nicht
> erwischt. Das ist die Grenze der Deklaration, nicht ein Fehler des
> Sensors — und sie gehört benannt.

**Sensor 2 — Retirement-Check** (*inferential feedback*,
ereignis-getriggert, kein periodischer Sweep): Eine Regel mit
Herkunfts-Anker wird **nicht entfernt oder gelockert**, ohne dass die
Herkunft konsultiert und das Ergebnis dokumentiert wurde — *„Regel seit
`welle-3` — ist die Beobachtung seither wieder aufgetreten?"*. Dieselbe
Bauart wie „Gates dürfen nicht ohne ADR gelockert werden", aber **kumulativ,
nicht ersetzend**: ist das verankerte Artefakt selbst ein Gate, gilt die
ADR-Pflicht unverändert weiter — der Retirement-Check beantwortet eine andere
Frage („ist der Grund entfallen?", nicht „darf ich?"). Er ist der
**Konsument** des Ankers; ohne ihn wäre der Anker eine zweite
write-only-Ablage.

Ziel-Form des Eintrags mit dem Pflichtfeld `liegt in <Zielort>` — zwei Orte, zwei
Vorlagen: für die Welle-Closure
[`../templates/docs/plan/planning/welle-results.template.md`](../templates/docs/plan/planning/welle-results.template.md),
für wellenlos verkörperte Regeln
[`../templates/docs/plan/planning/slice.template.md`](../templates/docs/plan/planning/slice.template.md)
§7.

<a id="jedes-artefakt-hat-einen-konsumenten"></a>
