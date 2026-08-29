# ADR-0070: Tabellen-Bedingungen unter einer Klammer, Zellengrenzen als Spaltenliste

**Status:** Accepted

**Datum:** 2026-08-29

**Autor:** pt9912

**Bezug:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(die umgeformte Anforderung),
[ADR-0057](0057-structure-tabellen-monotonie.md) (die Chronologie-Bedingung,
deren Schlüssel mitwandern),
[ADR-0069](0069-zellenlaenge-als-strukturbedingung.md) (die Zellenlängen-Bedingung,
deren Sachentscheide **unberührt** bleiben)

**Schärft:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v5.12.0/regelwerk/modul-04-adrs.md).

---

## Kontext

Die Zellenlängen-Bedingung aus
[ADR-0069](0069-zellenlaenge-als-strukturbedingung.md) war einen Tag alt, als
ihr erster echter Konsument zeigte, was ihre Form kostet. Vier Spalten des
ADR-Index zuzusagen brauchte **vier Regeln**:

```yaml
- files: docs/plan/adr/README.md
  section: "# ADR-Index — d-check"
  cell-max-column: Titel
  cell-max-chars: 200
  cell-min-chars: 10
- files: docs/plan/adr/README.md      # derselbe Selektor,
  section: "# ADR-Index — d-check"    # wortgleich wiederholt
  cell-max-column: Status
  …
```

**Drei Befunde, und der dritte ist der Anlass:**

1. **Der Selektor steht viermal identisch da.** Ändert sich die
   Abschnitts-Überschrift, müssen vier Stellen mit. Das ist heute genau
   **einmal** real — aber es ist die Form, die **jede** künftige
   Mehrspalten-Zusage erbt.
2. **`cell-max-column` ist falsch benannt, und der eigene Bestand belegt es.**
   Der Schlüssel benennt die Spalte *und* schaltet die Bedingung scharf; „max"
   im Namen ist eine dritte, unwahre Aussage. Die Regel über die Spalte `Bezug`
   trug `cell-max-column: Bezug` mit **ausschließlich** `cell-min-chars: 1` —
   dort log der Name.
3. **Fünf flache Schlüssel sprechen über dieselbe Sache, ohne es zu zeigen.**
   `table-order`, `table-column`, `cell-max-column`, `cell-max-chars`,
   `cell-min-chars` — dass sie alle die **Tabelle** eines Abschnitts meinen,
   trugen nur ihre Namenspräfixe.

**Das Zeitfenster ist der eigentliche Grund, das jetzt zu entscheiden.** Die
Bedingung stand unter `[Unreleased]`; außerhalb der eigenen `.d-check.yml` kam
kein Schlüssel vor, und keines der neun Schwester-Repos verwendete `table-order`
oder `table-column` (gemessen). Heute ist es eine Korrektur **vor** der
Auslieferung; nach dem Release wäre es eine Migration mit Konsumenten.

## Entscheidung

Alle tabellenbezogenen Bedingungen leben unter **einer Klammer** `table`, und
die Zellengrenzen werden zu einer **Liste je Spalte**:

```yaml
- files: docs/plan/adr/README.md
  section: "# ADR-Index — d-check"
  table:
    order: desc                  # war: table-order
    order-column: 2              # war: table-column
    column:                      # war: cell-max-column & Co., aber als LISTE
      - name: Titel
        cell-max-chars: 200
        cell-min-chars: 10
      - name: Status
        cell-max-chars: 40
```

1. **Die Klammer benennt, was zusammengehört.** Beide Bedingungen sprechen über
   dieselbe Tabelle — das stand vorher nirgends, außer in den Namenspräfixen.
2. **`column` ist eine Liste, und das ist der Kern.** Mehrere Spalten desselben
   Abschnitts stehen unter **einem** Selektor. Jeder Eintrag ist eine eigene
   Zusage und wird unabhängig ausgewertet, in Listen-Reihenfolge.
3. **Die Spalte wandert von der Regel-Identität ins Befund-Ziel.**
   [ADR-0069](0069-zellenlaenge-als-strukturbedingung.md) Entscheidung 7 legte
   sie in die **Identität**, weil zwei Spalten zwei Regeln brauchten. Das ist
   jetzt gegenstandslos: Spalten einer Liste können keine zwei Regeln
   kollidieren lassen. Getrennt werden müssen weiterhin die **Befunde** — zwei
   zu lange Zellen derselben Zeile trügen sonst dieselbe Adresse (Datei, Zeile,
   Regel, Ziel, Grund). Das leistet `ColumnTarget` je Befund, und **in
   derselben Form**: `… :: Spalte <name>`. Kein `target` ändert sich.
4. **`order-column` heißt so, weil der Name die Kopplung trägt.** Neben einer
   `column`-Liste wäre ein bloßes `column: 1` die zweite Bedeutung desselben
   Wortes. Der Präfix macht sichtbar, was der Config-Rand ohnehin erzwingt: ohne
   `order` ist der Schlüssel wirkungslos. Die **Abgrenzung** ist damit im Schema
   ablesbar — `order-column` adressiert über die **Position** (ein Griff daneben
   fällt als `section-cell-untyped` auf), `column[].name` über den **Namen**
   (eine Längen-Messung hat diese Selbstkontrolle nicht, deshalb meldet ein
   umbenannter Kopf laut).
5. **Die Chronologie-Spalte bleibt in der Regel-Identität.** Zwei
   Chronologie-Zusagen über denselben Abschnitt sind zwei Regeln und brauchen
   zwei Identitäten — anders als die Zellen-Spalten, die eine Liste sind.
6. **Die fünf Vorgänger-Schlüssel werden mit dem NEUEN Ort abgewiesen**, nicht
   still ignoriert und nicht der generischen `KnownFields`-Meldung überlassen
   (*„field table-order not found in type rawStructure"* — wahr, aber ohne den
   neuen Ort). Sie stehen dafür weiter im Config-Struct.
7. **Zwei neue Config-Ränder, die erst die Liste möglich macht:** eine leere
   `table`-Klammer (sie sagt nichts zu) und **derselbe Spaltenname zweimal** in
   einer Liste (beide Einträge trügen dasselbe Befund-Ziel und fielen unter die
   Deduplikation zusammen).

**Was diese ADR NICHT anfasst:** jeden Sachentscheid aus
[ADR-0069](0069-zellenlaenge-als-strukturbedingung.md) — Namens- statt
Positions-Bindung, Zeichen statt Bytes, Befund auf der Zeile, beidseitige
Grenze mit eigenen Grund-Codes, die zusammengeführte Zell-Zerlegung. Sie sind
unverändert gültig; diese ADR ersetzt ihre **Konfigurationsform**, nicht ihre
Begründung. Ebenso unberührt bleibt die Chronologie-Mechanik aus
[ADR-0057](0057-structure-tabellen-monotonie.md); dort wandern nur die
Schlüsselnamen.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v5.12.0/regelwerk/modul-04-adrs.md)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, vier Regeln stehen lassen | kein Eingriff in einen tags zuvor geschriebenen Vertrag | die Redundanz ist die Form, die jede künftige Mehrspalten-Zusage erbt; der Namens-Defekt bleibt, und nach dem Release kostet dieselbe Korrektur eine Migration mit Konsumenten |
| B — YAML-Anker/Alias für den wiederholten Selektor | keine Produktänderung | löst nur die Redundanz, nicht den Namens-Defekt; `gopkg.in/yaml.v3` führt **keine** Merge-Keys, ein Alias trägt also nur den **ganzen** Wert — und die vier Regeln unterscheiden sich ja gerade |
| C — nur umbenennen (`cell-max-column` → `cell-column`) | kleinster Eingriff, behebt den Namens-Defekt | lässt die vierfache Selektor-Wiederholung stehen — also genau die Beobachtung, die den Anlass gab |
| D — `cell-limits`-Liste ohne `table`-Klammer | löst Redundanz und Namens-Defekt, lässt den **released** `table-order` unangetastet | zwei Welten nebeneinander: die Zellen-Bedingung verschachtelt, die Chronologie flach — die Zusammengehörigkeit bliebe unsichtbar, und die nächste Tabellen-Bedingung müsste sich entscheiden |
| **E — volle `table`-Klammer, `column` als Liste (gewählt)** | eine Klammer für alles, was über die Tabelle spricht; ein Selektor für beliebig viele Spalten; der Namens-Defekt entfällt; die beiden Spalten-Begriffe sind im Schema unterscheidbar; zwei neue Config-Ränder fallen als Nebenprodukt an | migriert **fünf** Schlüssel statt drei, davon zwei aus einer **released** Bedingung ([ADR-0057](0057-structure-tabellen-monotonie.md)); braucht das Fangnetz für die Vorgänger-Namen |

**Zu E, benannt:** die released-Hälfte ist der Preis, und er ist **gemessen
klein** — kein Schwester-Repo und kein anderer Konsument führte einen der
Schlüssel. Wäre er es nicht, hätte D gewonnen.

## Konsequenzen

- **Positiv:** eine weitere Spalte kostet einen Listen-Eintrag statt einer
  Regel. Der eigene ADR-Index sagt vier Spalten mit **einem** Selektor zu, wo
  vorher `files`/`section` viermal wortgleich standen.
- **Positiv:** die Zusammengehörigkeit der Tabellen-Bedingungen ist im Schema
  ablesbar statt in Namenspräfixen — samt der Abgrenzung, dass die eine über
  Position und die andere über Namen adressiert.
- **Positiv, als Nebenprodukt:** zwei Config-Ränder, die es vorher nicht geben
  konnte. Der doppelte Spaltenname ist der interessantere — er war unter der
  alten Form ein Regel-Duplikat und wäre es unter einer nachlässigen Listen-Form
  **nicht** mehr gewesen.
- **Negativ, benannt:** ein Bruch am Konfigurations-Vertrag. Er ist
  fail-closed (Exit 2 mit dem neuen Ort), aber er ist einer — und er trifft
  auch die **released** Chronologie-Schlüssel, deren eigene Form nichts falsch
  gemacht hatte.
- **Negativ, benannt:** die fünf Vorgänger-Schlüssel bleiben als
  Migrations-Fangnetz im Config-Struct stehen. Das ist Ballast mit Ablaufdatum
  (siehe Trigger) — ein Wächter, der irgendwann nichts mehr fängt, ist genau die
  Bauform, die dieses Repo als
  [`BEO-013`](../planning/observations.md) führt.
- **Grenze des Fangnetzes:** es greift nur beim **exakten** Vorgänger-Namen.
  Eine Konfiguration, die den Schlüssel bereits anders verschrieben hatte,
  fällt weiter in die generische Decoder-Meldung.

## Fitness Function (falls maschinell prüfbar)

`make gates` — das Modul `structure` prüft sich über die eigene `.d-check.yml`
selbst (Dogfooding), und der Config-Rand ist getippt abgedeckt:

- `TestConfigYAMLStructureZellenlaenge` — die sieben Exit-2-Ränder der
  Spalten-Liste, darunter die **beiden neuen** (leere Klammer, doppelter Name),
  plus der Nachweis, dass zwei Spalten **einer** Regel durchgereicht werden.
- `TestConfigYAMLStructureAltSchluessel` — jeder der fünf Vorgänger-Schlüssel
  bricht **mit dem neuen Ort** in der Meldung.
- `TestCellMaxZweiSpaltenEineZeile` — zwei Befunde derselben Zeile behalten
  verschiedene Ziele.
- `TestStructureIdentitaetOhneSpalte` / `TestStructureIdentitaetMitOrderColumn`
  — die Zellen-Spalte steht **nicht** in der Identität, die explizite
  Chronologie-Spalte schon.
- `TestDocExamples_ConfigBeispieleValidieren` — jedes YAML-Beispiel des
  Handbuchs geht durch den echten Validator; es war der Sensor, der die
  Handbuch-Beispiele beim Umbau rot meldete.

**Was keine Fitness Function prüft:** ob die Klammer *verständlicher* ist als
fünf flache Schlüssel. Das ist ein Urteil.

## Re-Evaluierungs-Trigger

**Das Migrations-Fangnetz hat ein Ablaufdatum.** Es ist ein Hilfsmittel für
Konfigurationen aus der Zeit vor dieser ADR, kein dauerhafter Wächter. Beim
nächsten MAJOR-Bump — spätestens zu v1.0.0 — ist zu entscheiden, ob es
entfällt; danach fängt es messbar niemanden mehr, und ein Wächter ohne Beute
bleibt sonst stehen, weil ihn niemand in Frage stellt
([`BEO-013`](../planning/observations.md)).

**Zweiter Trigger:** kommt eine **dritte** tabellenbezogene Bedingung hinzu
(etwa eine Aussage über Zeilen), ist zu prüfen, ob `table` als Klammer sie
trägt oder ob die Struktur dann eine eigene Achse braucht. Der hier
ausdrücklich **nicht** aufgenommene `row`-Schlüssel ist der Kandidat: er wurde
weggelassen, weil es heute keine Bedingung gibt, die ihn füllt — ein leerer
Schlüssel wäre Vorratsbau und am Config-Rand ohnehin Exit 2.
