# ADR-0088: Zeilen-/Byte-Obergrenzen einer ganzen Datei werden ein eigenes Modul `file`, keine Erweiterung von `structure`

**Status:** Accepted

**Datum:** 2026-09-27

**Autor:** claude-sonnet-5

**Bezug:** Auftraggeber-Auftrag 2026-09-27 („max-rows und max-bytes einer
Datei", Ort auf Nachfrage entschieden: „ein neues Modul nehmen: file");
Anlass ist die offene Beobachtung
[`BEO-ALL/briefing-datei-ueberschreitet-lade-budget`](../planning/observations/BEO-ALL/briefing-datei-ueberschreitet-lade-budget/observation.md)
(`AGENTS.md` lädt in jeden Lauf, kein Sensor hält ihre Größe); slice-236 <!-- d-check:status-provenance -->.

**Schärft:**
[`DC-FA-FILE-001`](../../../spec/lastenheft.md#dc-fa-file-001--zeilen--und-byte-obergrenzen-einer-ganzen-datei-modul-file-opt-in)
(neue Anforderung — diese ADR trägt ihre Begründung), und die um `file`
erweiterte Modul-Liste in
[`DC-FA-CLI-002`](../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Ein Dokument, das in jeden Lauf geladen wird (`AGENTS.md`), wächst mit jeder
verkörperten Regel und schrumpft nie von selbst. Der Auftraggeber hat die
Größe als Mangel benannt; d-check kennt dafür keinen Sensor. Der naheliegende
Umweg — eine `structure[]`-Regel mit `forbid-pattern` gegen einen
Wiederholungs-Zähler — ist gemessen und trägt nicht (Plan-Anlass, slice-236 <!-- d-check:status-provenance --> §1): `structure` zählt den **bereinigten** Abschnittstext (Fenced-Code
entfernt, Inline-Code geleert), nicht die rohe Datei; es braucht einen
Abschnitts-Selektor (eine ATX-Überschrift); ein einzelner
Wiederholungs-Zähler ist bei RE2 auf 1000 begrenzt, größere Grenzen nur als
lange Verkettung; es zählt Zeichen, nicht Bytes; und es prüft nur Markdown.

Zwei Entscheide sind vorab zu treffen.

## Entscheidung

### (a) Eigenes Modul `file`, keine Erweiterung von `structure`

Auftraggeber-Entscheidung vom 2026-09-27. Drei Gründe, jeder für sich
hinreichend:

- **`structure` ist ausdrücklich abschnitts- und Markdown-bezogen**
  ([`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in):
  „prüft … ob Dokumente einer Klasse die Abschnitte
  tragen, die ihre Klasse verlangt"). Eine Zeilen-/Byte-Grenze der **ganzen**
  Datei braucht keinen Abschnitt und gilt jeder Dateiart — eine Erweiterung
  hätte ein Modul mit zwei Gegenständen erzeugt (Abschnitt **und** Datei),
  wie ADR-0084 für `targets`/`mentions` dieselbe Trennung nach *Gegenstand*
  zieht.
- **Der Umweg über `forbid-pattern` ist keine Erweiterung, sondern eine
  Zweckentfremdung** — er zählt genau das, was `structure` bewusst
  ausblendet (Fenced-Code, Inline-Code), und trifft an einer realen Grenze
  auf die RE2-Wiederholungsobergrenze (1000).
- **Nicht-Markdown-Dateien sind der Regelfall künftiger Nutzung** — ein Repo,
  das ein Skript oder eine Konfigurationsdatei begrenzen will, hat keinen
  Abschnitt, den `structure` je sehen könnte.

### (b) Zwei unabhängige Schwellen in einer Regel, nicht zwei Regel-Typen

Eine Regel trägt `max-lines` **und/oder** `max-bytes` über **demselben**
`files`-Glob — dieselbe Form wie `structure`s mehrere Bedingungen auf einem
Abschnitt. Mindestens eine der beiden ist Pflicht (halbe Aktivierung wäre
eine Regel, die nichts prüft); keine Sektions-, Muster- oder
Chronologie-Bedingung, weil es dafür keinen Gegenstand gibt (die ganze Datei
hat keine Abschnitte).

### Zählregel

Zeilen = Zeilenumbrüche, plus eine unvollständige Schlusszeile — dieselbe
Zählung, die `codepaths`/`citations` bereits teilen (`countLines`), und die
mit `wc -l` übereinstimmt, solange die Datei mit einem Zeilenumbruch endet.
Bytes = die rohe Dateigröße (`len(content)`), keine Zeichen-, Wort- oder
Token-Zählung.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Nichts tun** — der Umweg über `structure.forbid-pattern` bleibt die Antwort | kein neues Modul, keine neue Konfigurations-Fläche | zählt den bereinigten Text, nicht die rohe Datei; RE2-Wiederholungsobergrenze bei 1000; nur Markdown mit Abschnitt — an `AGENTS.md` selbst gemessen zweimal falsch (Grenzwert-Verschiebung durch Fenced-Blöcke, Verkettungs-Unlesbarkeit) |
| **`structure[]` um `max-lines`/`max-bytes` erweitern, `section`/`section-pattern` optional machen** | ein Modul weniger, Wiederverwendung der Glob-/Exempt-Mechanik | vermischt zwei Gegenstände unter einem Namen (Abschnitt und ganze Datei); eine `Accepted`-Anforderung ([`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)) müsste ihre Sektions-Pflicht aufgeben, was ihre Beschreibung („Abschnitte … tragen") widerlegt |
| **Eigenes Modul `file`** (gewählt) | eigener, einfacher Gegenstand (ganze Datei, jede Art); keine RE2-Wiederholungsgrenze, weil Zeilen/Bytes gezählt statt gematcht werden | 25. Regelmodul — mehr Modul-Fläche, mehr Spiegel (Verfügbar-Zeile, Grund-Code-Registry, `AllReasons()`) |

## Konsequenzen

- Ein **25. Regelmodul** (`validModules()` führte 24 Namen ohne `file`) — mehr
  Spiegel, die eine Modul-Aufnahme mitzieht
  ([`BEO-ALL/modulliste-spiegel-ungegated`](../planning/observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md)).
  In diesem Slice berührt: `validModules()`, die Grund-Code-Registry
  (`AllReasons()`/`reasonTexts()`, an `spec/spezifikation.md` §4 verriegelt),
  die „Verfügbar"-Zeile von `--print-config`. **Nicht** berührt, mit
  Begründung: das `.d-check.yml`/`FOCUS_DISABLE`-Spiegelpaar (`file` wird in
  diesem Slice nirgends aktiviert), ein kommentierter Beispiel-Block in
  `config_template.go` (Präzedenzfall `mentions`: kein Beispiel, weil es
  keine repo-unabhängige Schwelle gibt, die als Vorlage taugte), die
  `ai-harness`-Vorlage (`file` erfüllt K2 nicht — die Schwelle ist keine
  ableitbare oder konventionsfeste Größe).
- Neue Konfigurations-Fläche: `file[].files` (Pflicht), `max-lines`,
  `max-bytes` (mindestens eines Pflicht), `exempt-paths`, `hint`.
- **Die Implementierung ist Teil dieses Slice**, anders als ADR-0084: kein
  gesonderter Folge-Slice.
- **Keine Schwelle für `AGENTS.md` oder eine andere Datei** — diese ADR legt
  die Fähigkeit fest, nicht ihre Nutzung. Ein Wert ist eine
  Auftraggeber-Entscheidung mit Vertragswirkung (Gate-Einführung,
  `AGENTS.md` §3.6) und gehört in einen Folge-Slice.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `go test` | `TestFileModul*` (Grenzwert N/N+1 je Schwelle, Nullmengen-Härte, byte-identisch ohne Block) | `make test` |
| `go test` | `TestAllReasonsDeckungGegenSpezifikationGrundCodes` (die drei neuen Grund-Codes stehen in `spec/spezifikation.md` §4) | `make test` |
| `go test` | `TestPrintConfigVerfuegbarDecktRegistry` (`file` in der „Verfügbar"-Zeile) | `make test` |

## Re-Evaluierungs-Trigger

Zwei Bedingungen, jede für sich hinreichend:

1. **Ein Folge-Slice setzt eine Schwelle für `AGENTS.md` oder eine andere
   Datei fest.** Dann ist zu prüfen, ob die hier getroffenen Design-Entscheide
   (Zählregel, Pflicht mindestens einer Schwelle) tragen, oder ob die reale
   Anwendung eine dritte Bedingung nahelegt (z. B. eine Wortzahl).
2. **`structure` bekommt eine Fähigkeit, die den rohen Dateitext statt des
   bereinigten liest.** Dann verliert Grund 2 von (a) seine Grundlage, und
   die Trennung ist neu zu bewerten.

Ohne eines von beiden: permanent.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-09-27 | Proposed → Accepted (`slice-236`) |
