# Antwort auf den CR aus `ai-harness-course` — teamfähige BEO-Ablage

**Absender der Antwort:** d-check
**Datum:** 2026-08-31
**Bezug:** [eingehender CR](2026-08-31-cr-ai-harness-course-observations-relational.md)
samt der Antwort des Absenders auf unsere vier Rückfragen (liegt in seinem Repo)
**Ergebnis:** **angenommen — und aufgeschoben, auf euren eigenen Vorschlag hin.**
Wir bauen §1–§5 nach euren zwei Quell-Wellen, nicht davor.
**Einordnung:** kein Lastenheft-Bump, keine ADR — es wird nichts gebaut.

---

## Vorab

Eure Antwort zieht mehr zurück, als sie verteidigt, und sie stützt zwei ihrer
Punkte auf eine Messung an unserem Repo. Beides ist selten und beides ist der
Grund, warum diese Antwort kurz sein kann: Über §6, `threshold.count` und
`accept:` müssen wir nicht mehr verhandeln, und über §1 auch nicht.

Die Trennlinie, die ihr für §5 an unserem `planning`-Modul zieht — *was das
Dokument nennt, darf Konfiguration sein; was die Regel entscheidet, nicht* —
übernehmen wir als eure Formulierung unserer eigenen Politik. Sie ist schärfer
als das, was in
[`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
dazu steht.

## Euer Hinweis auf den „nächsten festen Boden" — nachgemessen

Ihr schlagt als Zwischenschritt die maschinelle Hälfte des **flachen** Registers
vor: Form · Anzahl · Lage. Wir haben sie gegen unseren Bestand gehalten, und das
Ergebnis ändert den Vorschlag:

| Prüfung | Stand bei uns |
|---|---|
| **Lage** — der Beleg liegt in `done/` | **bereits gedeckt.** Alle 21 Beleg-Zellen sind Links nach `done/`; das Modul `links` bricht, wenn das Ziel nicht dort liegt. Es braucht dafür keinen neuen Sensor. |
| **Form** — Slice-Kennung statt Freitext | erfüllt, aber ungewächtert. Kein Freitext-Beleg im Bestand. |
| **Anzahl** — so viele Belege wie der Zähler | **die einzige echte Lücke — und die, die euer Kanon heute blockiert.** |

Zwei der drei Prüfungen sind also entweder erledigt oder folgenlos. Die dritte
ist die, die euer Antrag gegenstandslos machen würde — und genau sie können wir
nicht bauen.

## Warum wir die Anzahl-Prüfung nicht bauen: eure eigenen zwei Zeilen

Eure Messung stimmt; wir haben sie unabhängig nachgeprüft:

- [`BEO-004`](../planning/observations.md) — Zähler **3**, genau **ein**
  Slice-Beleg.
- [`BEO-023`](../planning/observations.md) — Zähler **7**, **sechs** eindeutige
  Slice-Kennungen (`slice-178` steht zweimal).

Beide Abweichungen sind in der Zeile selbst begründet, keine Drift. Und beide
sind exakt die zwei Fälle, für die ihr in derselben Antwort einräumt, eine Regel
zu schulden: ein Vorkommen **außerhalb** einer Slice-Closure, und ein zweites
Vorkommen derselben Klasse **im selben** Slice.

Ein Gate gegen eine Regel zu bauen, deren Quelle sagt, dass sie unvollständig
ist, ließe uns am ersten Tag zwei Wege: das Register passend machen — also
Information löschen — oder einen Carveout eröffnen. Beides wäre schlechter als
die heutige Lage, in der die Abweichung **sichtbar und begründet** dasteht.

**Deshalb liefern wir statt eines Sensors den Beleg:** unsere zwei Zeilen sind
der empirische Fall für die zwei Regeln, die ihr schuldet. Sie stammen nicht aus
einem Gedankenexperiment, sondern aus einem Register, das seit welle-1 läuft.

## Zu §2 — die Verschärfung, die ihr selbst benannt habt

Ihr kennzeichnet §2 als Verschärfung gegenüber eurem Kanon: Modul 6 verlangt die
**Existenz** der Slice-Datei ausdrücklich nicht, §2 verlangt sie plus die
Rückreferenz. Wir teilen die Einordnung — das ist eine Regeländerung, keine
Mechanisierung, und sie gehört an der Quelle entschieden.

Ein Datenpunkt dazu, der euch die Entscheidung erleichtern könnte: **wir wären
heute konform.** Alle 21 Beleg-Zellen unseres Registers sind Links auf
existierende Dateien in `done/`. Die Verschärfung würde bei uns nichts brechen —
sie ist trotzdem eine Regeländerung, und wir bauen sie nicht vor eurer
Entscheidung.

## Was wir zusagen

- Wir bauen **§1–§5 nach euren zwei Quell-Wellen** (Sub-Area-Kürzel;
  geschlossene Ausgangs-Menge samt den beiden Beleg-Fällen oben) — in der
  Reihenfolge, die ihr vorschlagt.
- Die **Autoritäts-Quelle für Sub-Areas** kommt in Konfiguration **und**
  Anforderung, in der Bauform von `targets.authority`. Das ist unabhängig von
  eurem Kanon-Stand richtig und bleibt gesetzt.
  **Nachtrag nach eurem Zwischenstand (erste Quell-Welle steht):** die Zusage
  gilt weiter, bekommt aber die Bedingung, die ihr benannt habt. Die
  Kürzel-Spalte ist **bedingt** — wo Kennungen kein Bereichssegment tragen (wie
  in unserem Repo heute), entfällt sie. Eine Autoritäts-Prüfung, die eine
  **fehlende** Spalte als Fehler wertete, wäre bei jedem Adopter rot, der die
  Verzeichnisform nicht fährt. Die Prüfung wird deshalb so gebaut, dass die
  Abwesenheit der Spalte *nicht anwendbar* heißt und nicht *falsch* — und diese
  Grenze steht in der Anforderung, nicht nur im Code.
- Die **CI-Voraussetzung** (Merge-Stand) kommt als benannte Voraussetzung in die
  Anforderung, so wie das Modul `vcs` seinen `--range` nennt — nicht in die
  Prosa eines CR.
- §6 nehmen wir als eigenen Ausgabemodus entgegen, wenn ihr ihn einreicht. Eure
  Begründung trägt: der angezeigte Zähler muss aus **derselben** Ableitung
  kommen wie der gatende, sonst ist die zweite Zahl wieder da, die der Antrag
  abschaffen will.

## Nachtrag — euer `vcs`-Befund ist bestätigt und behoben

Er kam nach dieser Antwort und hängt an keinem ihrer Fäden; wir tragen ihn
trotzdem hier nach, weil er den Abschnitt *Bereits gelöst* eures Antrags
zurückzieht. Festgehalten als
[eigener Befund](2026-08-31-befund-ai-harness-course-vcs-rename.md).

**Reproduziert, alle vier Fälle** — eure Tabelle stimmt Zelle für Zelle,
einschließlich der beiden Kontrollen, die den Defekt eingrenzen: nur der
**erkannte** Rename fiel durch, Löschung und Rename-mit-Umformulierung meldeten
immer. Eure gelesene Ursache trifft ebenfalls; wir haben sie am Code bestätigt.

**Behoben** im Range-Pfad, mit dem Schnitt, den ihr als ersten genannt habt:
Diff **ohne** Rename-Erkennung, damit die Delete-Hälfte entsteht. Keine
Ähnlichkeits-Erkennung, kein neuer Grund-Code, kein Lastenheft-Bump — die
Anforderung sagte den Befund bereits zu, ohne einen Modus einzuschränken; sie
war nicht falsch, sie war nicht eingelöst.

**Zwei Dinge nehmen wir über euren Befund hinaus mit.** Erstens: der Kommentar
über der Übersetzungs-Funktion behauptete genau das Verhalten, das fehlte — er
hat den Fehler nicht nur nicht verhindert, sondern **gedeckt**; wer ihn las,
prüfte nicht nach. Zweitens: die Zusage war modus-abhängig, und das stand
nirgends — die Spezifikation nennt jetzt den Mechanismus, durch den eine
Umbenennung überhaupt sichtbar wird, statt ihn vorauszusetzen.

Euer vierter Punkt — *wo der Dateiname eine Aussage trägt, wird aus der Lücke
eine Fälschung* — ist der schärfste, und er trifft euren eigenen Entwurf: §2
bindet den Beleg an den Dateinamen. Er ist damit ein Argument **für** eure
Verzeichnisform und zugleich die Bedingung, unter der sie trägt.

## Was wir nicht tun

- **Keine Migration unseres eigenen Registers**, solange die drei Regeln fehlen.
  Sie wäre keine Probe eures Vertrags, sondern ein Informationsverlust.
- **Kein Modul auf Vorrat.** Eine Regel, die an der Quelle nicht steht, prüfen
  wir nicht — das ist dieselbe Politik, mit der wir eure Anträge sonst annehmen.
- **Keine Anzahl-Prüfung am flachen Register**, aus dem Grund oben. Sobald die
  zwei Regeln stehen, ist sie ohnehin gegenstandslos: ein abgeleiteter Zähler
  kann von seiner Belegliste nicht abweichen.
