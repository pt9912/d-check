# ADR-0052 — Platzhalter-Erkennung: Inline-Code ist keine Prosa, Nachfilter sind Code

**Status:** Proposed
**Datum:** 2026-08-10
**Autor:** pt9912
**Bezug:** [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(vierte Struktur-Bedingung), [§`DC-FA-PLAN-001.a`](../../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
Schritt C4b. **Verfeinert** [ADR-0048](0048-closure-note-struktur-im-planning-modul.md)
um eine vierte Bedingung derselben Fähigkeit.
**Change Request** des Konsumenten `ai-harness-course` (CR 2).

## Kontext

Ein unausgefüllter Vorlagen-Rumpf ist syntaktisch vollständig und passiert
**alle drei** bestehenden Struktur-Bedingungen. Gemessen gegen v0.52.0 mit der
Notiz `Ergebnis: <ergebnis>. Belege: <belege>. Offen: <offen>. Ende: <ende>.`:
**0 Befunde** — der Abschnitt existiert, vier Satzende-Zeichen erreichen die
Schwelle, keine deklarierte Floskel kommt vor.

Der Konsument liefert eine an realen Falsch-Positiven gehärtete Regex — mit
**Lookbehind und Lookahead**. Go nutzt RE2; das Muster wird fail-closed
abgelehnt. Beide Lookarounds sind feste Zeichen-Prüfungen und lassen sich
ersetzen, indem man das Vorzeichen **konsumiert** statt hineinzuschauen. Das ist
eine Portierung mit eigener Testlast, keine Übernahme.

**Winkelklammern sind in technischer Prosa häufig.** Ein Falsch-Positiv ist hier
teurer als ein übersehener Platzhalter: es macht das Gate unglaubwürdig, und ein
unglaubwürdiges Gate wird abgeschaltet.

**Bestandsmessung am eigenen Repo (2026-08-10)**, über die 96 Closure-Notizen in
`done/`, mit der portierten Form:

| Fassung | Treffer | echte Platzhalter |
|---|---|---|
| naiv (`<` … `>`) | 24 | 0 |
| portiert (Vorzeichen + erstes Zeichen) | 12 | 0 |
| portiert **ohne** Inline-Code | **0** | 0 |

**Alle zwölf** Treffer der portierten Fassung liegen in Inline-Code —
``<PREFIX>`, ``<a id>`, ``<datei>` —, **keiner** außerhalb.

## Entscheidung

1. **Inline-Code ist keine Prosa.** Vor der Suche werden die Backtick-Spans des
   Abschnitts geleert. Das ist keine Heuristik, sondern die Rolle von
   Inline-Code: dort wird Syntax **gezeigt**. Ein Platzhalter, den jemand
   ausfüllen soll, steht im Fließtext. Die Messung entscheidet die Frage
   eindeutig — ohne die Einschränkung wären zwölf von zwölf Treffern
   Falsch-Positive, und das Gate wäre am ersten Tag unglaubwürdig.

2. **Die Nachfilter sind Code, nicht Muster.** Autolink/Adresse (`://`, `@`) und
   HTML-Tag-Namen werden **nach** dem Treffer geprüft. In die Regex gepresst
   wären sie unlesbar und unprüfbar; als Code sind sie einzeln testbar, und die
   Tag-Liste ist eine Liste statt einer Alternation.

3. **Die Substanz-Zählung bleibt unberührt.** Sie zählt Satzende-Zeichen
   weiterhin einschließlich Inline-Code. Die engere Sicht gilt **nur** der neuen
   Bedingung.
   Das ist eine bewusste Asymmetrie, keine Nachlässigkeit: die Zählung ist eine
   **ausgelieferte** Schwelle, und sie zu verengen lässt bisher grüne Notizen rot
   werden — in beide Richtungen wirksam und darum eine eigene Entscheidung mit
   eigener Bestandsmessung. Diese Bedingung ist dagegen rein additiv (Default
   `false`).

4. **Erster Treffer je Kandidat**, wie bei der Floskel. Mehrere Platzhalter
   derselben Notiz sind dieselbe Reparatur; ein Befund je Platzhalter machte den
   Befundsatz laut, ohne mehr zu sagen.

5. **Opt-in, Default `false`.** Anders als die drei bestehenden Bedingungen ist
   diese abhängig von der Schreibkultur des Repos. Der Default schaltet nichts
   frei, was ein Adopter nicht angefordert hat — dieselbe Zurückhaltung wie bei
   der per Default **leeren** Floskel-Liste.

6. **SemVer: Minor.** Neuer Grund-Code; ohne den Schalter ist der Befundsatz
   byte-identisch.

## Verglichene Alternativen

| Alternative | Warum verworfen |
|---|---|
| Konsumenten-Regex übernehmen | Nutzt Lookarounds; RE2 lehnt sie fail-closed ab. Keine Wahl, sondern eine Sprachgrenze |
| Naive Form (`<` … `>`) | 24 Treffer auf dem eigenen Bestand, **null** davon echt |
| Inline-Code mitzählen | Zwölf Falsch-Positive auf dem eigenen Bestand — das Gate wäre am ersten Tag unglaubwürdig |
| Nachfilter in die Regex | Unlesbar und unprüfbar; die HTML-Tag-Liste würde zur Alternation, die niemand pflegt |
| Whitespace im Inneren zulassen (`<mein feld>`) | Die Messwert-Prosa `<1 s und der Recall >0,9` wird dann zum Treffer — gemessen. Ein Platzhalter ist ein **Feldname**; mehrwortige Formen sind ein Re-Evaluierungs-Trigger, keine Erweiterung des Musters |
| Eingerückte Code-Blöcke hier behandeln | d-check modelliert sie **nirgends**; eine sechste Sicht auf denselben Text zu bauen wäre genau die Klasse, gegen die diese ADR argumentiert. Als Grenze benannt |
| Ein Befund je Platzhalter | Dieselbe Reparatur, mehrfach gemeldet; die Floskel-Bedingung hat dieselbe Frage bereits entschieden |
| Default `true` | Winkelklammern sind schreibkultur-abhängig; ein Adopter bekäme beim Update Befunde, die er nicht bestellt hat |
| Zugleich die Substanz-Zählung verengen | Bewegt eine **ausgelieferte** Schwelle in beide Richtungen — eigene Entscheidung, eigene Messung, eigener Slice |

## Konsequenzen

- **Der eigene Bestand ist messbar frei von Falsch-Positiven** (0 von 96
  Notizen). Der Schalter kann in der eigenen Konfiguration gesetzt werden, ohne
  eine Sanierung auszulösen.
- **Zwei Sichten auf denselben Abschnitt** leben nebeneinander: die Zählung sieht
  Inline-Code, die Platzhalter-Suche nicht. Das ist dokumentiert und begründet,
  aber es ist eine Naht — wer sie später zusammenführt, muss beide Richtungen
  messen.
- **Zwei Grenzen sind benannt, nicht geschlossen:** der eingerückte Code-Block
  (vier Leerzeichen — in d-check nirgends modelliert) und die ungerade
  Backtick-Parität im Absatz. Beide sind Eigenschaften der geteilten Lexik und
  gehören dorthin, nicht in diese Bedingung.
- **Die Erkennung ist eine Portierung, kein Import.** Sie trägt eigene Tests je
  Falsch-Positiv-Klasse; die Vorlage des Konsumenten bleibt Vorlage.

## Fitness Function

- **Der Template-Rumpf meldet** — genau ein Befund, an der Zeile des ersten
  Treffers.
- **Technische Prosa bleibt grün:** Vergleichszeichen in **beiden**
  Schreibweisen (mit und ohne Leerzeichen), Tabellenzellen, Generics, Autolink,
  Adresse, HTML-Tag mit und ohne Attribut, Linkziel in Winkelklammern und in
  Inline-Code gezeigte Meta-Syntax, **einzeln** geprüft, nicht als
  Sammel-Fixture.
- **Jeder Eintrag der HTML-Tag-Liste wirkt** — über die Liste iteriert, nicht
  stichprobenartig.
- **Ohne den Schalter byte-identischer Befundsatz.**
- **Der eigene Bestand bleibt bei null** — jede Abweichung ist ein echter Fund.

## Re-Evaluierungs-Trigger

- Wenn ein Konsument eine Platzhalter-Form meldet, die diese Erkennung verfehlt
  (etwa `{{feld}}` oder `TODO:`), ist zu entscheiden, ob die Form
  konfigurierbar wird — **nicht**, ob die Regex wächst.
- Wenn die Substanz-Zählung auf dieselbe engere Sicht umgestellt wird, ist
  Entscheidung 3 erledigt und die Naht verschwindet.
- Wenn Falsch-Positive trotz der drei Einschränkungen auftreten, ist die
  Erkennung zu **verengen**, nicht abzuschalten.

## Geschichte

- 2026-08-10: Proposed (doc-first, `slice-098`).
- 2026-08-10: nach unabhängigem Review verengt — das Innere muss **frei von
  Whitespace** sein (die Zusage zu den Vergleichszeichen hielt nur für die
  Schreibweise mit Leerzeichen, gemessen), ein Winkelklammer-Linkziel ist ein
  dritter Nachfilter, und die beiden offenen Grenzen sind benannt.
