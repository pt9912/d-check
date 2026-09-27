# ADR-0089: `structure` bekommt eine zwölfte Bedingung `max-lines` (Zeilenbudget eines Abschnitts), Untergrenze 1

**Status:** Accepted

**Datum:** 2026-09-27

**Autor:** claude-sonnet-5

**Bezug:** Eingehender Change Request des Adopters `ai-harness-course`
(`docs/plan/cr/2026-09-27-cr-eingehend-ai-harness-course-structure-max-lines.md`,
2026-09-27, Priorität niedrig, additiv); slice-237 <!-- d-check:status-provenance -->.

**Schärft:**
[`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(Erweiterung — kein neues Kürzel, nach dem etablierten Schnitt-Kriterium:
Einzelmodul-Frage, dieselbe Zusagen-Familie wie `min-sentences`/`max-tasks`).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der CR bittet um eine Obergrenze für die Länge eines Abschnitts, damit ein
wachsendes Hard-Rule-Dokument (Anlass beim Adopter: eine `AGENTS.md`-Sektion)
auffällt, bevor sie unhandlich wird. `structure` kennt bereits zwei zählende
Bedingungen auf demselben bereinigten Text (`min-sentences` als Untergrenze,
`max-tasks` als Obergrenze über Task-Items) — eine dritte, `max-lines`, ist
strukturell dieselbe Familie: eine Schwelle über eine Zahl, die aus dem
bereinigten Abschnittstext abgeleitet wird.

Zwei Design-Fragen sind vorab zu klären, weil sie den Umfang der Bedingung
festlegen, bevor der Code sie festschreibt.

## Entscheidung

### (a) Erweiterung von `structure`, kein eigenes Modul

Anders als bei ADR-0088 (`file`) ist hier **kein** neuer Gegenstand im Spiel:
`max-lines` zählt denselben bereinigten Abschnittstext, den `min-sentences`
und `max-tasks` bereits zählen — dieselbe Grundmenge, derselbe Selektor
(`section`/`section-pattern`), dieselbe Kardinalitäts-Behandlung. Es gibt
keinen zweiten Gegenstand zu trennen (der Unterschied zu `file[].max-lines`
ist gerade *derselbe Name, anderer Gegenstand* — geregelt in der
Spezifikation, nicht in der Modul-Grenze).

### (b) Grundmenge: bereinigter Text, Fenced-Code zählt nicht mit

`max-lines` misst auf demselben Text wie `min-sentences`
(`SectionProse`, `internal/hexagon/core/rules/sections.go`), der Fenced-Code
vollständig entfernt statt zu maskieren. Ein Abschnitt mit einem langen
Beispiel im Codeblock wächst dadurch nicht gegen `max-lines`. Für den Anlass
(Prosa-Wachstum) ist das die richtige Eigenschaft: Codeblöcke sind Beleg, kein
Text, den ein Leser überfliegt. Die Kehrseite ist eine benannte Grenze (siehe
Konsequenzen).

### (c) Untergrenze `1`, nicht `0`

`max-lines: 0` wäre am Config-Rand syntaktisch zulässig, aber praktisch ein
Abschnitts-Verbot: Der bereinigte Body eines Abschnitts mit einer Überschrift
ist so gut wie nie leer (die Roh-Überschrift selbst zählt nicht mit, siehe
Zählregel), und die einzige Konfiguration, die je befundfrei liefe, wäre ein
leerer Abschnitt — dieselbe Zusage, die `non-empty` schon abdeckt, nur über
den Umweg einer Zeilengrenze. Eine Bedingung, die einen legitimen Nutzwert
nur in einem Fall hat, den eine andere Bedingung bereits trägt, ist keine
sinnvolle Obergrenze, sondern ein verkleidetes Verbot. Analog zu
`min-sentences`/`max-tasks`/`max-open-tasks` (alle < 1 bzw. < 0 ⇒ Exit 2)
zieht `max-lines` seine Untergrenze dort, wo die Bedingung noch etwas anderes
sagt als eine Nachbarbedingung.

### Zählregel

Gezählt werden die Zeilenumbrüche (`\n`) im bereinigten Abschnittstext, keine
`wc -l`-Korrektur um eine unvollständige Schlusszeile: `SectionProse`
rekonstruiert den Text zeilenweise mit einem Zeilenumbruch **hinter jeder**
erhaltenen Zeile, auch der letzten — anders als eine rohe Datei
(`file[].max-lines`, `codepaths`/`citations`), die ohne abschließenden
Zeilenumbruch enden kann.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Nichts tun** — der Adopter dehnt Abschnitte manuell im Review | kein neuer Code, keine neue Konfigurations-Fläche | genau der Anlass des CR: ein wachsender Abschnitt fällt niemandem auf, bevor er unhandlich ist — dieselbe Lücke, die `max-tasks`/`max-open-tasks` für Task-Listen bereits schließen |
| **`forbid-pattern`-Wiederholungs-Muster** (ein RE2 wie `(.*\n){N,}` gegen den bereinigten Text) | kein neuer Schlüssel, kein neuer Grund-Code | zählt Zeilen nur indirekt über eine Wiederholungs-Klammer, die bei RE2 auf 1000 begrenzt ist; die Meldung nennt kein Ist/Soll (`forbid-pattern` sagt nur „Muster trifft", keine Zahl); dieselbe Zweckentfremdung, die ADR-0088 für `file` bereits gegen `structure.forbid-pattern` verworfen hat |
| **Zeichen- statt Zeilen-Zählung** (`max-chars` auf dem bereinigten Text, analog zu `cell-max-chars`) | eine Schwelle, keine Rundungsfrage bei langen/kurzen Zeilen | der Anlass ist ein **visuelles** Zeilenbudget (Leser überfliegt eine Sektion), keine Byte-/Zeichen-Zusage wie bei einer Tabellenzelle; eine Zeichenzahl korreliert nur lose mit der Zeilenzahl, die ein Editor anzeigt, und der CR fragt ausdrücklich nach Zeilen |
| **`max-lines` als zwölfte Bedingung, Grundmenge = bereinigter Text, Untergrenze 1** (gewählt) | dieselbe Familie wie `min-sentences`/`max-tasks` (ein Selektor, mehrere Zahlen-Bedingungen); direkte Antwort auf den CR; Fenced-Code-Ausnahme schließt an eine etablierte Eigenschaft an | Fenced-Code-lastige Abschnitte bleiben unter der Schwelle, obwohl die Datei wächst (siehe Konsequenzen); Namens-Kollision mit `file[].max-lines` |

## Konsequenzen

- Neuer Grund-Code `section-lines-exceeded` ([`SPEC-087`](../../../spec/spezifikation.md#4-grund--und-fehler-codes)),
  neue Schema-Zeile `structure[].max-lines`. Ohne den Schlüssel byte-identisches
  Verhalten.
- **Benannte Grenze:** Ein Abschnitt mit viel Beispielcode und wenig Prosa
  kann `max-lines` nie auslösen, obwohl die Datei insgesamt wächst — dafür
  ist `file[].max-lines` (ADR-0088) der richtige Träger, nicht diese
  Bedingung.
- **Namens-Kollision, keine Vertrags-Kollision:** `structure[].max-lines` und
  `file[].max-lines` leben unter verschiedenen Top-Level-Blöcken und
  kollidieren syntaktisch nicht, tragen aber unterschiedliche Grundmengen
  (bereinigter Abschnitt vs. rohe Datei). Die Abgrenzung steht in der
  Spezifikation an beiden Schema-Zeilen.
- **Keine Schwelle für `AGENTS.md` § Harte Regeln oder eine andere konkrete
  Datei** — diese ADR legt die Fähigkeit fest, nicht ihre Nutzung. Eine
  Schwelle ist eine Auftraggeber-Entscheidung mit Vertragswirkung
  (`AGENTS.md` §3.6) und gehört in einen Folge-Slice.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `go test` | `TestStructureMaxLines_*` (Grenzwert N/N+1, Fenced-Code zählt nicht mit, fehlender Abschnitt bleibt `section-missing`, Mehrdeutigkeit bleibt `section-ambiguous`, `hint` gewinnt, byte-identisch ohne den Schlüssel) | `make test` |
| `go test` | `TestAllReasonsDeckungGegenSpezifikationGrundCodes` (`section-lines-exceeded` steht in `spec/spezifikation.md` §4) | `make test` |

## Re-Evaluierungs-Trigger

Zwei Bedingungen, jede für sich hinreichend:

1. **Ein Folge-Slice setzt eine Schwelle für eine konkrete Datei fest**
   (z. B. `AGENTS.md`). Dann ist zu prüfen, ob die Fenced-Code-Ausnahme für
   diesen Anwendungsfall trägt, oder ob ein zweiter Modus (rohe Zeilen im
   Abschnitt) nötig wird.
2. **Die Namens-Kollision mit `file[].max-lines` erzeugt wiederholte
   Verwechslung** (Support-Anfrage, Fehlkonfiguration im Feld). Dann ist eine
   Umbenennung eines der beiden Schlüssel neu zu bewerten — trotz des in
   slice-237 <!-- d-check:status-provenance --> benannten Ausschlusses, diesen Slice damit zu verquicken.

Ohne eines von beiden: permanent.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-09-27 | Proposed → Accepted (`slice-237`) |
