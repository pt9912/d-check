# ADR-0069: Zellenlänge als `structure`-Bedingung — Spalte über den Namen, beidseitig begrenzt

**Status:** Accepted (teil-superseded: ADR-0070)

**Datum:** 2026-08-28

**Autor:** pt9912

**Bezug:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(die erweiterte Anforderung),
[ADR-0057](0057-structure-tabellen-monotonie.md) (die siebte Bedingung, deren
Tabellen-Mechanik geteilt wird),
[ADR-0059](0059-closure-waechter-weicht-structure-regel.md) (die Präzedenz: ein
Zeichenketten-Wächter weicht einer `structure`-Regel)

**Schärft:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v5.12.0/regelwerk/modul-04-adrs.md).

---

## Kontext

Eine Tabellenspalte, die einen **Titel** trägt, nimmt still die ganze
Entscheidung auf. Der ADR-Index dieses Repos ist der Beleg — gemessen über
alle 68 Zeilen:

| Größe | Median | Maximum |
|---|---|---|
| Titel-Spalte im Index | 442 | 2206 |
| H1-Titel der ADR-Dateien | 77 | 158 |

44 der 68 Zeilen liegen über 200 Zeichen. **Alle 68 H1-Titel sind genau ein
Satz** — die Datei sagt es also bereits kurz, nur die Sicht auf sie nicht.

**Das ist mehr als Kosmetik.** Der Kopf derselben Datei deklariert: *„Derivativ:
Quelle der Wahrheit sind die **ADR-Dateien**; dieser Index ist eine
Bequemlichkeits-Sicht."* Eine Zelle, die die Entscheidung wiederholt statt auf
sie zu zeigen, macht die Sicht zu einer **zweiten Quelle**, die abdriften kann,
ohne dass ein Gate es merkt. Kein bestehendes Modul misst Zellen-Umfang.

Die naheliegende Antwort war ein `forbid-pattern` — die achte Bedingung existiert
bereits, und ein Muster ist billig. Sie wurde gebaut und **gemessen**, bevor sie
verworfen wurde; die Gegenüberstellung steht unten.

## Entscheidung

Wir führen die **neunte `structure`-Bedingung** ein: `cell-max-column`
(Kopfzeilen-**Name** der Spalte) und `cell-max-chars` (Schwelle in **Zeichen**).

1. **Die Spalte wird über ihren Kopfzeilen-Namen adressiert, nicht über eine
   Position.** Eine eingefügte Spalte verschöbe eine Positions-Angabe **still**
   auf die falsche Spalte; ein umbenannter Kopf meldet **laut**
   (`section-column-missing`). Das ist bewusst **anders** als bei
   `table-column` der siebten Bedingung — dort ist die Schlüsselspalte
   typisiert, ein Griff daneben fällt als `section-cell-untyped` auf; eine
   Längen-Messung hat diese Selbstkontrolle nicht.
2. **Der Befund steht auf der Zeile der Zelle**, nicht am Abschnittskopf. Bei
   44 Treffern ist das der Unterschied zwischen einer Meldung und einer
   Arbeitsliste.
3. **Gezählt werden Zeichen, nicht Bytes.** Die Schwelle beschreibt einen Text;
   ein Umlaut ist ein Zeichen.
4. **Die Bindung erfolgt je Tabelle, und Nicht-Treffer sind kein Befund.** Eine
   Nebentabelle im selben Abschnitt schaltet die Messung nicht ab. Ein
   **doppelter** Name in einer Kopfzeile, eine **zu kurze** Datenzeile und der
   **Leerlauf** (keine Tabelle des Abschnitts bindet die Spalte) sind Befunde:
   die Bedingung zu setzen **ist** die Behauptung, dass es diese Spalte gibt.
5. **Die Zell-Zerlegung des Produkts wird auf eine Antwort zusammengeführt.**
   Sie ist ab jetzt escape- und backtick-bewusst: `\|` ist ein Zeichen der
   Zelle und **kein** Zelltrenner. Ohne diesen Zug erbte die neue Bedingung
   genau das Loch, gegen das sie gebaut ist (siehe unten).
6. **Die Bedingung begrenzt beidseitig** (`cell-min-chars`). Eine Obergrenze
   allein lässt die **leere** Zelle passieren — null Zeichen liegen unter jeder
   Schwelle —, und genau das ist der stille Grün-Pfad, den dieses Repo sonst
   jagt. Eigener Grund-Code, weil die Reparatur eine andere ist: ausfüllen
   statt kürzen. Mindestens eine der beiden Grenzen ist Pflicht.
7. **Die Regel-Identität trägt die ausdrücklich benannte Spalte.** Ohne das
   wäre nur **eine** Spalte je Abschnitt konfigurierbar (zweite Regel ⇒
   Konfigurations-Duplikat, Exit 2) — und schlimmer: zwei zu lange Zellen
   **derselben Zeile** trügen dieselbe Befund-Adresse (Datei, Zeile, Regel,
   Ziel, Grund) und fielen unter die Deduplikation zusammen. Eine Regel, die
   **keine** Spalte nennt, behält ihre Identität unverändert; damit ändert sich
   für keine Bestandsregel ein `target`. Das löst zugleich die in
   [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
   §Out-of-Scope benannte Grenze *„zwei Chronologie-Zusagen über denselben
   Abschnitt"*, die dort ausdrücklich einen Change Request einlud.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v5.12.0/regelwerk/modul-04-adrs.md)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, die Spalte per Review kurz halten | kein Code | die Wachsamkeit hat 68-mal nicht gehalten; genau das ist die Bauform, die dieses Repo als `BEO-013` führt |
| B — `forbid-pattern` mit einer Längen-Klasse | eine Zeile Konfiguration, kein Code | der Befund nennt den **Abschnitt** statt der Zeile (gemessen: Zeile 1 gegen Zeile 5); eine escapte Pipe zerteilt den Lauf und die Zelle entkommt (gemessen an einer konstruierten Zelle); das Muster hängt an der Form der **Nachbar**-Zelle |
| C — die Bedingung bauen, Spalte über die **Position** (wie `table-column`) | konsistent mit der siebten Bedingung | eine eingefügte Spalte verschiebt die Messung **still** auf die falsche — bei einer Längen-Messung fällt das nie auf |
| D — die Bedingung bauen, aber auf dem **naiven** Zerleger | kleinerer Eingriff, keine Fremd-Fläche berührt | erbt Bs Escape-Loch: die Bedingung hätte genau die Schwäche, mit der sie begründet wird |
| E — nur eine **Obergrenze**, ohne `cell-min-chars` | ein Schlüssel weniger | die **leere** Zelle passiert grün (0 liegt unter jeder Schwelle) — ein Gate, das „die Spalte ist kurz" mit „die Spalte ist gefüllt" verwechselt |
| **F — Namens-Bindung, beidseitige Grenze, zusammengeführte Zell-Antwort, Spalte in der Identität (gewählt)** | Befund auf der Zeile; kein Escape-Loch; ein umbenannter Kopf meldet laut; die leere Zelle fällt auf; mehrere Spalten je Abschnitt sind konfigurierbar **und** unterscheidbar; die benannte Grenze aus [ADR-0057](0057-structure-tabellen-monotonie.md) **entfällt**, statt vererbt zu werden | berührt zwei Bestands-Bedingungen (`table-order`, `planning.waves`); drei Grund-Codes und sieben Config-Ränder mehr |

## Konsequenzen

- **Positiv:** die Sicht kann wieder derivativ sein, und ein Rückfall meldet
  sich auf der Zeile, auf der er entsteht.
- **Positiv:** das Produkt hat **eine** Zell-Antwort statt zwei. Die in
  [ADR-0057](0057-structure-tabellen-monotonie.md) §Konsequenzen benannte
  Grenze — *ein Pipe in einem Backtick-Span verschiebt die Spaltenadresse* —
  entfällt damit für `table-order` und `planning.waves` mit.
- **Negativ, benannt:** gemessen wird die Zelle, **wie sie dasteht** —
  Markdown-Syntax eingeschlossen. Eine Zelle aus einem einzigen langen Link ist
  lang, auch wenn ihr sichtbarer Text kurz ist; wer eine Link-Spalte begrenzen
  will, misst etwas anderes, als er meint.
- **Positiv, gemessen:** die Spalten `Datum` und `Bezug` decken eine Hälfte ab,
  die vorher niemand sah — **vier** Zeilen des Index trugen nur 3 von 5 Zellen,
  ohne Datum und ohne Bezug. Kein Gate hat das je gemeldet: die vorderen
  Spalten standen ja da.
- **Negativ:** eine Schwelle ist eine Zahl und kein Urteil. 200 kommt aus dem
  heutigen Bestand (längster H1: 158); ein künftiger längerer Titel fällt auf.
  Das ist gewollt — *„ein Satz"* ist die Regel, 200 ihr grober Wächter.
- **Negativ:** die Bindung an den Kopfzeilen-Namen macht die Regel abhängig von
  einer Zeichenkette in der geprüften Datei. Sie fällt dabei **laut** aus
  (`section-column-missing`), nicht still — das ist der Unterschied zu C, aber
  eine Abhängigkeit bleibt es.
- **Folgepflicht:** wird die Bedingung auf eine Spalte mit Links angewandt, ist
  vorher zu entscheiden, ob die Syntax mitgezählt werden soll — sonst misst das
  Gate etwas anderes, als sein Konfigurator meint.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `d-check` (Modul `structure`, Dogfooding) | im ADR-Index trägt `Titel` 10–200 Zeichen, `Status` 8–40, `Datum` genau 10, `Bezug` mindestens 1 (`section-cell-oversized`/`section-cell-undersized` auf ihrer Zeile); jede Spalte ist adressierbar, und eine Zeile, die nicht bis zu ihr reicht, meldet `section-column-missing` | `make doc-check` (in `gates`) |

## Re-Evaluierungs-Trigger

Meldet die Bedingung wiederholt Zeilen, deren Titel **berechtigt** länger als
200 Zeichen sind, ist nicht die Zeile falsch, sondern die Schwelle — dann ist
sie neu aus dem dann gültigen Bestand zu ziehen.

Zweiter Trigger: wird eine Spalte gemessen, deren Inhalt überwiegend
Markdown-Syntax ist (Links, Bilder), greift die oben benannte Grenze. Dann ist
zu entscheiden, ob die Bedingung eine zweite Zähl-Semantik braucht — oder ob
die Anwendung falsch war.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-08-28 | Accepted — Anlass war der Hinweis des Auftraggebers auf eine missbrauchte Titel-Spalte in einem Schwester-Repo; die Messung ergab, dass dieses Repo beim Median schlechter dasteht (442 gegen 305) | [ADR-0059](0059-closure-waechter-weicht-structure-regel.md) |
