# Eingehender Change Request — `structure` bekommt ein Abschnitts-Zeilenbudget

**Absender:** Adopter `ai-harness-course` · **Eingegangen:** 2026-09-27
**Richtung:** eingehend — dieses Repo ist der **Empfänger**, nicht der Bittsteller.
**Ziel-Dokument:** [`spec/lastenheft.md`](../../../spec/lastenheft.md)
**Berührt:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(Modul `structure`)
**Stand:** **angenommen mit zwei Korrekturen** — die eigene Zählung des CR
(„zehnte Bedingung") und sein Out-of-Scope-Beleg sind nicht übernommen, der
Schluss (`max-lines`, keine Alters-Mechanik) ist es; Träger `slice-237`.

**Ablage-Hinweis.** Ein **eingehender** CR ist die dritte Klasse neben
[`MR-035`](../../../harness/conventions.md#mr-035) (ausgehend) und
[`MR-036`](../../../harness/conventions.md#mr-036) (Antworten darauf). Die
Datei liegt hier aus demselben Grund wie ihre Vorgänger in diesem
Verzeichnis: Was gebeten wurde und wie entschieden wird, soll im Repo
stehen und nicht nur im Vorgang.

---

## Wortlaut

> CR an d-check: structure — Abschnitts-Zeilenbudget als zehnte Bedingung ([`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in))
>
> Von: ai-harness-course (Konsument) · Stand: 2026-09-27
> Bezug d-check: [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) (Modul structure, Bedingungen 1–9: non-empty, min-sentences, max-tasks/max-open-tasks, forbid-pattern, require-pattern, require-all, table-order/table-column, headings-match/headings-level, cell-max-chars/cell-min-chars)
> Priorität: niedrig, kein Blocker · Art: additiv, Erweiterung einer bestehenden Anforderung, kein neues Kürzel
>
> ## Anlass
>
> Guide-Dateien, die ein Agent bei jedem Lauf liest — AGENTS.md (Hard Rules), harness/README.md §Sensors —, wachsen in der Praxis unbemerkt: Eine Regel, die einen Einzelfall behebt, bleibt stehen, auch wenn ein Gate sie längst deckt. Gemessen an einem Konsumenten-Repo (pg-change-feed, .claude/commands/implement-slice.md, dieselbe Artefaktklasse „Workflow-Skelett"): 175 → 367 Zeilen über 26 Commits in 18 Tagen, mehr als verdoppelt, ohne einen einzigen Rückbau-Commit dazwischen. Erst ein bewusster, einmaliger Aufräum-Zug (mit eigener ADR) brachte die Datei auf 338 Zeilen zurück — keine Routine, ein Einzelfall.
>
> structure prüft mit max-tasks/max-open-tasks bereits eine Zähl-Eigenschaft innerhalb eines Abschnitts (Task-Items). Eine analoge Bedingung für die Zeilenzahl eines Abschnitts gibt es nicht — die neun bestehenden Bedingungen prüfen Vollständigkeit, Reihenfolge, Form und Zellenlänge, aber keine misst die schiere Größe eines wiederkehrend gelesenen Abschnitts.
>
> ## Vorschlag: zehnte Bedingung max-lines
>
> ```yaml
> structure:
>   - files: AGENTS.md
>     section: "## Harte Regeln"
>     max-lines: 180
> ```
>
> - Kardinalität: sections: one (Default) reicht — bei AGENTS.md §Harte Regeln gibt es genau einen solchen Abschnitt. each wäre für ein Repo mit mehreren so benannten Abschnitten trotzdem sinnvoll, wenn structure das ohnehin unterstützt.
> - Gemessen wird dieselbe Abschnitts-Fassung, auf der non-empty/min-sentences schon operieren (der „bereinigte Abschnitts-Text") — keine neue Bereinigungs-Logik, nur eine neue Zählung darauf.
> - Neuer Grund-Code: section-lines-exceeded (bewusst nicht section-oversized — das trägt max-tasks schon). Meldung auf der Überschriften-Zeile des Abschnitts, wie bei section-heading-mismatch bei Leerlauf.
> - Fail-closed-Ränder: max-lines ≤ 0, ein section, der im Dokument nicht existiert (→ eigener Befund, kein stilles Bestehen), sections: one mit mehreren Treffern (→ section-ambiguous, bestehendes Verhalten, keine Messung).
>
> ## Die Ratchet-Wirkung entsteht nicht im Werkzeug
>
> max-lines prüft nur „Ist ≤ Deklariert". Dass eine Anhebung sichtbar bleibt (statt eines stillen Wachstums), ist eine Eigenschaft von Git, nicht von d-check — eine Änderung an max-lines steht im Diff, eine Änderung an Prosa-Zeilen sonst nicht zwingend. Kein neues Konzept in d-check nötig, nur die neue Zählung.
>
> ## Bewusst nicht gewünscht
>
> - Keine Stichtags-/Alters-Mechanik (wann wurde der Abschnitt zuletzt geprüft). [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) schließt das laut Lastenheft ausdrücklich aus („eine Stichtags-Mechanik... hieße, die Kennungs-Konvention des Adopters zu interpretieren"); dieser CR bleibt bei einer reinen Struktur-Invariante.
> - Keine Aussage über Inhalt/Qualität der Zeilen — nur ihre Zahl.
> - Kein Zwang zum Rückbau. Das Gate meldet nur den Ist-Zustand gegen die eingecheckte Zahl.
>
> ## Prüfung (Break-Test)
>
> - Abschnitt mit 181 Zeilen bei max-lines: 180 → section-lines-exceeded, gemeldet auf der Überschriften-Zeile.
> - Abschnitt mit 180 Zeilen → kein Befund (Grenzwert inklusiv, wie bei cell-max-chars).
> - section kommt im Dokument nicht vor → eigener Befund, kein stilles Bestehen (fail-closed, wie bei section-column-missing).
> - max-lines gesenkt, Abschnitt unverändert → sofort rot, bis der Abschnitt tatsächlich kürzer ist (kein „Gnadenfrist"-Verhalten).
>
> ## Kompatibilität und Aufwand
>
> Rein additiv — ohne max-lines-Schlüssel ändert sich nichts am bestehenden Befundsatz. Aufwand: eine Zählfunktion auf bereits vorhandenem, bereinigtem Abschnitts-Text, ein neuer Grund-Code, ein Testfall, Lastenheft-/Handbuch-/CHANGELOG-Zeile — vom Umfang her vergleichbar mit der max-tasks-Erweiterung.

---

## Entscheidung

**Angenommen, mit zwei Korrekturen gegenüber dem Wortlaut:**

1. **Der Ordinal-Zähler ist falsch.** Der CR bündelt `max-tasks`/`max-open-tasks`
   zu einem Slot und führt `open-tasks-require-marker`/
   `-marker-section` gar nicht — macht `max-lines` zur „zehnten". Die eigene,
   durchgehende Ordinal-Zählung dieses Lastenhefts (achte: `headings-match`,
   neunte: `cell-max-chars`, zehnte: `max-open-tasks`, elfte:
   `open-tasks-require-marker`) macht `max-lines` zur **zwölften**. Übernommen
   ist die eigene Zählung.
2. **Der Out-of-Scope-Beleg trägt nicht.** Der zitierte Satz („eine
   Stichtags-Regel, die die Kennungs-Konvention des Adopters interpretieren
   müsste") schließt einen **ID-Stichtag** aus, nicht eine Alters-/Frische-Prüfung.
   Der Schluss des CR (keine Alters-Mechanik gewünscht) bleibt richtig, der
   Beleg dafür wird nicht übernommen.

**Eine vom CR nicht benannte Grenze wird ergänzt:** `max-lines` zählt auf dem
**bereinigten** Abschnittstext (`SectionProse`), der Fenced-Code-Zeilen
**vollständig entfernt** — ein Abschnitt mit viel Beispielcode und wenig
Prosa erreicht die Schwelle dadurch nie. Für den genannten Anlass
(Prosa-Wachstum) ist das die richtige Eigenschaft; sie steht als benannte
Grenze in Lastenheft und Spezifikation.

Begründung in begleitender [ADR-0089](../adr/0089-structure-max-lines-zwoelfte-bedingung.md).
