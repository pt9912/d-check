# ADR-0114: Das Modul `linelength` begrenzt die Zeichen je Zeile in den Markdown-Dateien der Scan-Menge

**Status:** Proposed

**Datum:** 2026-11-22

**Autor:** pt9912

**Bezug:** Auftrag des Auftraggebers vom 2026-11-22 (Zeilenlänge in Markdown-Dateien, nur Markdown),
[`DC-FA-LLEN-001`](../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in),
[`DC-FA-SCAN-001`](../../../spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln),
[`DC-FA-CONF-002`](../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope),
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit),
[ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) (die geteilte Lexik bindet ihre Konsumenten),
[ADR-0088](0088-file-modul-groessengrenzen.md) und [ADR-0089](0089-structure-max-lines-zwoelfte-bedingung.md)
(die zwei Nachbar-Achsen), [ADR-0069](0069-zellenlaenge-als-strukturbedingung.md) (dieselbe Zähleinheit),
[ADR-0063](0063-marker-form-folgt-der-kommentar-lexik-der-eingabe.md) (der Zeilen-Marker),
[ADR-0113](0113-pinsync-modul.md) (dieselbe Zeilen-Zerlegung),
slice-289 <!-- d-check:status-provenance -->.

**Schärft:**
[`DC-FA-LLEN-001.a`](../../../spec/spezifikation.md#dc-fa-llen-001a--zeilenlänge-in-markdown-dateien-linelength)
(neues Modul-Verfahren samt Grund-Code `line-chars-exceeded`) und das Schema `linelength.*` in
[`SPEC-005`](../../../spec/spezifikation.md#spec-005--d-checkyml); die Zeile des Grund-Codes in
[§4](../../../spec/spezifikation.md#4-grund--und-fehler-codes) kommt mit dem registrierten Modul.

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v6.18.1/regelwerk/modul-04-adrs.md).

---

## Kontext

Ein Repo, das für seine Markdown-Dateien eine Zeilenlänge vorgibt, hat dafür in d-check keinen Sensor.
Drei Achsen messen schon eine Länge, keine davon die Breite einer Zeile: `file` die Zeilen und Bytes einer
ganzen Datei jeder Art, `structure` die Zeilen eines Abschnitts im bereinigten Text und die Zeichen einer
Tabellenzelle. Der Auftraggeber hat die Prüfung auf Markdown-Dateien beschränkt und die Empfehlungen zu
Modulname, Eingabemenge, Zähleinheit, Schlüssel, Ausnahmen, Befund, Bezugsmenge und leerer Menge
bestätigt; das Architect-Verdikt vom selben Tag hat Entscheid 8 gegen das Lastenheft gestellt, Entscheid 9
ergänzt und die Ränder von 5 (c) und 7 geschärft. Entscheid 8 in dieser Fassung und Entscheid 9 hat der
Auftraggeber am 2026-11-22 bestätigt. Entscheid 5 (c) hat er am 2026-11-23 in der strengen Fassung neu
bestätigt: Die zuerst bestätigte Fassung nahm, am eigenen Bestand gemessen, zu fast der Hälfte Zeilen aus,
die ein Umbruch unter die Grenze gebracht hätte (Alternative M).

Drei Randbedingungen gelten unabhängig davon. Ein Lauf ohne `linelength` bleibt byte-identisch. Das Modul
schreibt nie. Und die Fragen „ist das Fenced-Code“ und „ist das eine Tabellenzeile“ hat das Produkt schon
beantwortet: `proseLineSet` und `tableRowLine` in der geteilten Lexik; eine zweite Antwort daneben wäre ein
Defekt ([ADR-0054](0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 1).

## Entscheidung

Wir führen ein **eigenes opt-in-Modul `linelength`** mit diesen Festlegungen (die Nummern dienen dieser
Liste als Bezug):

1. **Name und Kürzel:** Modul `linelength`, Bereichskürzel `LLEN` (`LINK` ist belegt).
2. **Eingabemenge:** die Markdown-Dateien der Scan-Menge (`scan.roots`/`scan.ignore`), mit modul-lokalem
   `scope` nach [`DC-FA-CONF-002`](../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope) und
   `exempt-paths` wie die übrigen Prosa-Module. Das Modul liest keine Datei außerhalb dieser Menge.
3. **Zähleinheit:** Unicode-Codepunkte der **rohen** Zeile, Einrückung eingeschlossen. Zeilen enden an `\n`;
   ein einzelnes `\r` am Zeilenende gehört nicht zur Zeile — auch am Ende einer letzten Zeile ohne `\n`,
   dieselbe Zerlegung wie bei `pinsync` ([ADR-0113](0113-pinsync-modul.md) Entscheidung 2). Ein Tab ist ein
   Zeichen. Ungültiges UTF-8 zählt je Byte als ein Zeichen.
4. **Schwelle:** `linelength.max-chars`, Pflicht bei aktivem Modul, ganzzahlig ≥ 1, sonst Exit 2. Befund bei
   `Länge > max-chars`; genau auf der Grenze ist befundfrei.
5. **Ausnahmen, je einzeln schaltbar, Default „ausgenommen“.** Überschriften zählen mit.
   - (a) `exempt-code-blocks`: Zeilen, die nicht in `proseLineSet` liegen — Fenced-Code samt Zäunen, nach dem
     Umschalt-Automaten der geteilten Lexik.
   - (b) `exempt-tables`: Zeilen, für die `tableRowLine` gilt — außerhalb von Fenced-Code, die rohe Zeile
     beginnt mit `|`. Eine Trennzeile wird nicht verlangt, eingerückte Zeilen sind keine.
   - (c) `exempt-unbreakable`: überlange Zeilen, die ein **Wort** tragen, das allein länger als `max-chars`
     ist. Ein Wort ist eine längste Folge von Codepunkten ohne Leerraum; Leerraum sind Space und Tab, kein
     anderes Zeichen. Eine solche Zeile bringt kein Umbruch unter die Grenze. Umgekehrt gilt das nicht: nicht
     jede andere überlange Zeile bringt ein Umbruch darunter — eine Überschrift lässt sich nicht umbrechen, und
     ein Präfix der Markdown-Form kann das erste Wort eines Listenpunkts oder Zitats über die Grenze schieben
     (Konsequenzen). Jede andere überlange Zeile ist ein Befund; die Ausnahme prüft das Wort, nicht die Form.
6. **Befund:** je überlanger Zeile ein Befund `line-chars-exceeded`, Ort = Datei und Zeile, die Meldung nennt
   gemessene Länge und Grenze. Reihenfolge nach Pfad, dann Zeile.
7. **Bezugsmenge:** der Lauf trägt eine Zeile mit der Zahl der geprüften Dateien, der geprüften Zeilen und der
   über (c) ausgenommenen Zeilen. Zeilen zählen wie bei `file` (Zeilenumbrüche plus eine unvollständige
   Schlusszeile). „Geprüft“ heißt: nach `scope`, `exempt-paths`, (a) und (b). Form wie bei `mentions` und
   `pinsync`: stderr im Default-Modus, `summary.notes` unter `--json`/`--yaml`.
8. **Leere Menge ist kein Fehler.** Trifft die Eingabemenge keine Markdown-Datei, meldet das Modul nichts, und
   die Bezugsmenge nennt `0`. Eine explizit leere `linelength.scope.roots` prüft nichts, wie
   [`DC-FA-CONF-002`](../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope) es jedem scannenden
   Modul zusagt. Unberührt bleibt der Umgebungsfehler einer gänzlich leeren Repository-Wurzel (Exit 2), der
   für jeden Lauf gilt, nicht für dieses Modul.
9. **Kein Zeilen-Marker.** `linelength` honoriert `d-check:ignore` nicht; die benannte Liste der
   honorierenden Module bleibt bei vier. Eine überlange Zeile ist deterministisch und behebbar; die Ventile
   sind (a)–(c), `scope` und `exempt-paths`.

### Abgrenzung der Längen-Achsen

| Schlüssel | misst | Gegenstand | Text |
|---|---|---|---|
| `file[].max-lines` | Zeilen (Höhe) | ganze Datei jeder Art, eigener Glob | roh |
| `file[].max-bytes` | Bytes | ganze Datei jeder Art, eigener Glob | roh |
| `structure[].max-lines` | Zeilenumbrüche (Höhe) | ein Abschnitt einer Markdown-Datei | bereinigt, ohne Fenced-Code |
| `structure[].table.column[].cell-max-chars` | Codepunkte | eine Tabellenzelle | getrimmt |
| `linelength.max-chars` | Codepunkte (Breite) | jede Zeile einer Markdown-Datei der Scan-Menge | roh, Einrückung zählt |

Keine zwei Zeilen teilen Gegenstand und Maß. Der Name `max-chars` meidet eine dritte Lesart von `max-lines`
und teilt mit `cell-max-chars` die Einheit, nicht den Gegenstand.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, ein fremder Markdown-Linter | keine neue Modul-Fläche | eine zweite Toolchain im Repo, und deren Fence- und Tabellen-Antwort ist nicht die des Produkts: zwei Antworten auf dieselbe Frage in einem Lauf |
| B — dritte Schwelle im Modul `file` | kein neues Modul | `file` liest jede Dateiart roh über eigene Globs und verspricht „unabhängig von der Dateiart“; Markdown-Ausnahmen brächen diese Zusage, und `file` kennt die Scan-Menge nicht |
| C — neue Bedingung in `structure` | Glob- und Ausnahme-Mechanik vorhanden | jede Regel braucht einen Abschnitts-Selektor, `structure` wählt seine Dateien selbst und kennt kein `scope`, und es misst den bereinigten Text; eine Breite der ganzen Datei über der Scan-Menge hat keinen dieser Anker |
| D — Zähleinheit Bytes | trivial, keine Kodierungsfrage | ein Umlaut zählte doppelt; die Grenze beschriebe die Kodierung statt den Text |
| E — Zähleinheit Anzeigebreite | entspricht dem, was ein Editor zeigt | braucht eine Breiten-Tabelle (East-Asian-Wide, Emoji), die sich mit der Unicode-Version ändert; das gefährdet [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus) |
| F — eigener Fence- und Tabellen-Leser nach CommonMark-strenger Lesart | schließt Zäune längengleich, erkennt GFM-Tabellen ohne führende Pipe | eine zweite Antwort neben `proseLineSet`/`tableRowLine` — genau der Defekt aus ADR-0054 Entscheidung 1; was gemessen wird, wiche von dem ab, was `links` oder `structure` als Fenced-Code und Tabelle lesen |
| G — (c) als „nur der Teil bis zum letzten Leerzeichen zählt“ | liest sich als eigene Regel | dieselbe Bedingung wie M, bis auf eine Randposition; **keine** eigene Option |
| I — keine Ausnahme (c), stattdessen der Marker `d-check:ignore` | jede Ausnahme steht sichtbar an ihrer Zeile | der Marker verlängert die Zeile, die er ausnimmt; er verlässt die benannte Liste aus vier Modulen; und das Lastenheft schaltet deterministische Befunde nicht zeilenweise stumm |
| J — leere Menge ⇒ Exit 2 (Empfehlung des Plans) | kein Lauf ohne Gegenstand passiert grün | widerspricht [`DC-FA-CONF-002`](../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope) für eine explizit leere `scope.roots`; Exit 2 ist ein Eingabefehler, und eine leere Scan-Menge ist keiner — auch die übrigen scannenden Module laufen darüber ohne Fehler |
| K — leere Menge ⇒ Befund (Exit 1), ausgenommen eine explizit leere `scope.roots` | fail-closed und mit [`DC-FA-CONF-002`](../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope) verträglich | ein Grund-Code mehr für einen Fall, den die Bezugsmenge schon zeigt; die Präzedenz (`file-no-match`) gilt Modulen, deren Regel Dateien behauptet — `linelength` behauptet keine |
| M — (c) nach dem letzten Leerraum: ausgenommen ist eine Zeile, deren Codepunkte an den Positionen `max-chars + 1` bis Zeilenende keinen Leerraum tragen (zuerst bestätigt am 2026-11-22) | eine Zeile mit einem langen Link am Ende muss nicht vor dem Link umbrechen; eine Prüfung nur ab der Grenze | gröber, als ihr Name sagt: am eigenen Bestand (Stand `f3dc94da`) bei `max-chars: 120` über die Markdown-Dateien der Scan-Menge, nach (a) und (b), nimmt sie von 2952 überlangen Zeilen 828 aus (2124 bleiben Befunde) — von den 828 tragen 151 keinen Leerraum, 316 ein Wort über 120 und 361 keines, sodass ein Umbruch sie unter die Grenze brächte; bei 141 davon ist das letzte Wort höchstens 20 Zeichen lang, gewöhnliche Prosa, die ein Wort über die Grenze ragt |
| **H — eigenes Modul, Codepunkte der rohen Zeile, Ausnahmen aus der geteilten Lexik, (c) nur bei einem Wort über der Grenze, leere Menge kein Fehler, kein Marker** | nimmt nur Zeilen aus, die kein Umbruch unter die Grenze bringt; am Bestand oben nimmt H 513 Zeilen aus und lässt 2439 Befunde — 361 Zeilen, die M ausnahm, werden Befunde, und 46, die M meldete, fallen weg (ein Wort über der Grenze vor umbrechbarem Text hinter ihr), netto 315 Befunde mehr; gleicher Aufwand wie M; die Bezugsmenge sagt, was gemessen wurde | erbt die Grenzen der geteilten Lexik; eine Zeile mit umbrechbarem Text **und** einem Wort über der Grenze bleibt ganz ausgenommen (Konsequenzen) |

## Konsequenzen

- Positiv: ein Repo bekommt die Zeilen-Grenze mit einer Konfiguration; das Modul braucht keinen Port und
  keine neue Lexik. Eine überlange Prosa-Zeile ist ein Befund, sobald ein Umbruch sie unter die Grenze bringt.
- Negativ — Grenzen der geteilten Lexik, in der Anforderung zu nennen:
  - Der Fence-Automat schaltet bei jeder ```` ``` ````- oder `~~~`-Zeile um. Ein ````` ```` `````-Block mit einer
    inneren ```` ``` ````-Zeile und ein ```` ``` ````-Block mit einer `~~~`-Zeile schließen früh; die Zeilen dahinter
    werden gemessen.
  - `TrimFenceIndent` entfernt jede Einrückung aus Space und Tab: auch ein Zaun mit vier und mehr
    Leerzeichen davor schaltet um. Ein eingerückter Code-Block ohne Zaun und ein Zaun in einem Zitat sind
    dagegen kein Fenced-Code; ihre Zeilen werden gemessen.
  - Ein nie geschlossener Zaun nimmt den Rest der Datei aus. Das meldet `spans` mit `fence-unclosed`;
    `linelength` meldet es nicht zum zweiten Mal.
  - Eine `|`-Zeile ohne Trennzeile ist ausgenommen. Eine eingerückte Tabellenzeile, eine in einem Zitat und
    eine GFM-Zeile ohne führende Pipe werden gemessen. Mit (a) aus und (b) an ist eine `|`-Zeile im Zaun keine
    Tabellenzeile.
- Negativ — Grenzen der Entscheidung selbst:
  - (c) nimmt eine Zeile **ganz** aus, sobald sie ein Wort über der Grenze trägt — auch wenn umbrechbarer Text
    davor oder dahinter steht, der auf eine eigene Zeile passte. Die Ausnahme prüft das Wort, nicht den Rest.
  - (c) liest keine Markdown-Form. Ein Wort von höchstens `max-chars` Zeichen, das das **erste** Wort hinter
    einem Listen- oder Zitat-Marker (`- `, `1. `, `> `) ist und erst mit dem Marker und der Einrückung davor
    über die Grenze ragt, macht die Zeile zu einem Befund, den kein Umbruch behebt: das erste Wort bleibt auf
    der Zeile seines Markers. Dasselbe gilt für jedes Wort einer Überschrift, die sich nicht umbrechen lässt,
    und für jedes Wort eines eingerückten Code-Blocks, dessen Zeilen alle ihre Einrückung tragen. Ein späteres
    Wort eines Listenpunkts oder Zitats lässt sich dagegen umbrechen: eine Fortsetzungszeile braucht nach
    CommonMark (*lazy continuation*) weder Marker noch Einrückung. Die Reparatur der Befunde, die kein
    Umbruch behebt, ist Kürzen des Worts oder der Form.
  - Eine Überschrift lässt sich nicht umbrechen; die Reparatur ist Kürzen, und Kürzen ändert ihren Anker. Die
    gebrochenen Verweise melden `links` und `anchors`.
  - Ein `d-check:ignore`-Kommentar, den ein anderes Modul verlangt, verlängert seine Zeile und zählt mit. Für
    (c) ist er Text wie jeder andere: trägt er ein Wort, das länger als `max-chars` ist, nimmt (c) die Zeile über
    dieses Wort aus. Das geschieht, sobald die Grenze kleiner ist als ein Wort des Kommentars —
    `d-check:ignore` hat 14 Zeichen, `d-check:status-provenance` 25, ein Pfad in `d-check:cite` oft mehr;
    zulässig ist `max-chars` ab 1. Ist die Grenze mindestens so groß wie jedes Wort des Kommentars, nimmt (c)
    die Zeile nur aus, wenn ein anderes Wort auf ihr länger als die Grenze ist.
  - Das eigene Profil des Repos schaltet das Modul nicht ein; diese Entscheidung legt die Fähigkeit fest, nicht
    ihre Nutzung.
- Folgepflicht: jede Randbedingung der Anforderung (Grenzwert, Einheit, jede Ausnahme an und aus, die Ränder
  von (c), Bezugsmenge mit und ohne abschließendes `\n`, leere Menge, Marker-Zeile) hat einen Test im
  Modulkern, der ohne die Regel aus dem richtigen Grund rot läuft; ein Kopplungs-Test fährt dieselbe Eingabe
  durch `linelength` und die bisherigen Konsumenten von `proseLineSet`/`tableRowLine` (ADR-0054
  Entscheidung 4); das Modul wird zusammen mit der Modul-Aufzählung (Lastenheft, READMEs, Handbuch,
  `operations.md`, `--print-config`) registriert.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `go test` | Grenzwert N/N+1, je Ausnahme an und aus, Ränder von (c), Bezugsmenge, leere Menge | `make test` |
| `go test` | Kopplungs-Test: `linelength` und die bisherigen Konsumenten antworten auf dieselbe Eingabe gleich, welche Zeilen Fenced-Code und welche Tabellenzeilen sind | `make test` |
| `go test` | `model.ValidModules()` stimmt mit den Modul-Spiegeln überein; `AllReasons()` mit §4 der Spezifikation | `make test` |
| a-check | der Kern liest Dateien nur über den Filesystem-Port und trägt keine YAML-Bibliothek | `make arch-check` |

## Re-Evaluierungs-Trigger

Die Entscheidung wird neu geprüft, wenn (a) die geteilte Lexik ihre Fence- oder Tabellen-Antwort ändert
(etwa längengleiches Schließen oder eine verlangte Trennzeile) — dann ändern sich (a) und (b) mit, und es ist
zu prüfen, ob das gewollt ist; (b) ein Konsument eine Zeile ausnehmen muss, die weder Code noch Tabelle noch
unumbrechbar ist (dann der Marker, Alternative I); (c) ein Konsument, ein Review oder eine Verifikation
überlange Zeilen, die (c) wegen eines einzelnen Wortes über der Grenze ganz ausnimmt, als Mangel liest — etwa
umbrechbaren Text vor einer langen URL (dann eine Prüfung des Rests neben dem Wort) —, oder umgekehrt den
Umbruch vor einem langen Link am Zeilenende als Last (dann Alternative M); (d) ein Konsument Text mit breiten
Zeichen prüft und die Codepunkte dort nicht tragen (dann Alternative E); oder (e) ein Konsument die Prüfung
für eine andere Dateiart verlangt (dann ist die Abgrenzung zu `file` neu zu ziehen, Alternative B).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-11-22 | Proposed: Entscheide 1–7 wie vom Auftraggeber bestätigt, 5 und 7 an den Rändern geschärft; Entscheid 8 gegen [`DC-FA-CONF-002`](../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope) geändert und Entscheid 9 neu — beide brauchen die Bestätigung des Auftraggebers vor `Accepted` | slice-289 <!-- d-check:status-provenance --> |
| 2026-11-22 | Proposed, fortgeschrieben: Entscheid 8 und 9 vom Auftraggeber bestätigt; Bezug und Schärft verlinken die Anforderung und ihr Verfahren; Entscheid 3 an die Zeilen-Zerlegung von `pinsync` angeglichen; die Konsequenz von 5 (c) gemessen und Alternative H geschärft — für 5 (c) ist eine Neu-Bestätigung empfohlen (Empfehlung H), die Fassung oben gilt bis dahin | slice-289 <!-- d-check:status-provenance --> |
| 2026-11-23 | Proposed, fortgeschrieben: Entscheid 5 (c) vom Auftraggeber als Alternative H neu bestätigt — ausgenommen ist nur eine Zeile mit einem Wort über der Grenze; die zuerst bestätigte Fassung ist mit ihren Messzahlen verworfene Alternative M; Konsequenzen und Trigger (c) nachgezogen | slice-289 <!-- d-check:status-provenance --> |
| 2026-11-23 | Proposed, berichtigt: die Zahlen zu H in den Alternativen — H nimmt am Bestand 513 Zeilen aus, nicht 467 (die 467 zählten nur innerhalb der Ausnahmen von M); netto 315 Befunde mehr; Stand der Messung benannt, die Zahl der überlangen Zeilen auf diesem Stand ist 2952 | slice-289 <!-- d-check:status-provenance --> |
| 2026-11-23 | Proposed, nach Review R1 fortgeschrieben: die Begründung von 5 (c) behauptet nicht mehr, jede andere überlange Zeile brächte ein Umbruch unter die Grenze (Überschriften, Präfixe der Markdown-Form); die Marker-Konsequenz gilt für jede Grenze ab 1; das Schema ist über [`SPEC-005`](../../../spec/spezifikation.md#spec-005--d-checkyml) adressiert; die Entscheidung ist unverändert | slice-289 <!-- d-check:status-provenance --> |
| 2026-11-23 | Proposed, nach Review R2 fortgeschrieben: die Konsequenz zum Präfix der Markdown-Form gilt nur für das erste Wort eines Listenpunkts oder Zitats, für Überschriften und eingerückten Code; ein späteres Wort lässt sich auf eine Fortsetzungszeile ohne Marker umbrechen; die Entscheidung ist unverändert | slice-289 <!-- d-check:status-provenance --> |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
