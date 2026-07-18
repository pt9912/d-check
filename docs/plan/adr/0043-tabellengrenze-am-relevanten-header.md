# ADR-0043 — Tabellengrenze am relevanten Header: der stille Tabellen-Übersprung

**Status:** Accepted
**Datum:** 2026-07-18
**Autor:** pt9912
**Schärft:** [`DC-FA-REQ-001.a`](../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen) (Schritt 5, Datenzeilen/Tabellengrenze)
**Bezug:** [`DC-FA-REQ-001`](../../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen), [`DC-FA-XREF-001`](../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in), [ADR-0037](0037-trace-tabellenquellen-nullmengen-guard.md), [ADR-0040](0040-kommentar-suffix-in-tabellenzeilen.md), [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus).

## Kontext

Der header-gebundene Tabellen-Reader (`markdownTables`/`consumeTableRows`, geteilt
von `trace.requirements.format: table` und `trace.cross-consistency`) konsumiert
Datenzeilen **greedy** bis zur ersten Leer-/Nicht-Tabellenzeile oder einem
Zellenzahl-Bruch. Zwei Tabellen ohne Leerzeile dazwischen sind für ihn **eine**
Tabelle: die zweite Header- und Trennzeile werden als Datenzeilen der ersten
gelesen.

Das ist ein **stiller Waisen-Verlust**, wenn die **erste** Tabelle irrelevant ist
(ihr Header bindet keine konfigurierte Rolle) und die **zweite** relevant: die
irrelevante Tabelle verschluckt die relevante samt Header, Trennzeile und
Datenzeilen; `extractTable` meldet für die irrelevante nichts, und die
Anforderungen der zweiten existieren nicht — **kein Gate sagt es**.

Belegt gegen das **ausgelieferte** v0.47.0-Image, Fixture `fx-s`,
**ohne jeden Marker** — eine relevante Tabelle, Leerzeile, dann eine irrelevante
gleicher Breite, unmittelbar gefolgt von einer relevanten:

```text
v0.47.0  --trace  ⇒  1 Anforderung(en) (nur F-1).  F-2/F-3 verschluckt.
```

Der Defekt ist **älter** als die Direktiven-Toleranz aus
[ADR-0040](0040-kommentar-suffix-in-tabellenzeilen.md) und **unabhängig** von ihr.
Er ist zugleich deren **Wurzel:** die dortige Toleranz scheiterte fünfmal, weil das
Lesbarmachen einer Direktiven-Datenzeile den `badLine`-Abbruch entfernt, der die
Folgetabelle heute nur **zufällig** rettet (der Kommentar erzeugt eine überzählige
Zelle, jede `badLine` setzt den Header-Scan neu auf). Ohne robuste Tabellengrenze
kehrt dieser stille Pfad bei jeder Toleranz-Fassung zurück (Review R3-F-1).

**Zwei gemessene Wahrheiten, die den Lösungsraum einschränken:**

- **GFM gibt uns recht.** goldmark v1.8.4 sieht `fx-s` **exakt** wie der heutige
  Reader: ohne Leerzeile beginnt keine neue Tabelle. Die Frage ist **Policy**, nicht
  Grammatik — ein GFM-Parser schließt die Klasse nicht.
- **Die naheliegende Regel ist mehrdeutig.** `fx-s` (`[Textzeile][Trennzeile][Daten…]`,
  soll trennen) und die Gegenprobe `fx-t` (eine Datenzeile aus lauter `-` gefolgt
  von einer weiteren, soll **nicht** trennen) sind **syntaktisch identisch**. Jede
  rein strukturelle Grenze („eine Tabelle endet, wo ein Header+Trennzeile beginnt")
  zerreißt `fx-t` — genau daran wurde sie in den fünf früheren Fassungen vor dem
  Code widerlegt.

## Entscheidung

1. **Die Tabellengrenze liegt am *relevanten* Header.** Während der Datenzeilen-
   Konsumtion beendet eine Zeile die laufende Tabelle genau dann, wenn sie mit ihrer
   Folgezeile eine gültige Header+Trennzeile bildet (bestehendes `tableHeaderAt`)
   **und** ihr Header eine konfigurierte **Rolle bindet** (relevant — `bindTableColumns`
   für `format: table`, `bindCrossColumns` für `cross-consistency`). Der Re-Scan
   erkennt die neue Tabelle. Ein Header, der **keine** Rolle bindet, beendet die
   laufende Tabelle **nicht** — er wird wie bisher als Datenzeile konsumiert.

2. **Das „relevant" bricht die Mehrdeutigkeit.** In `fx-t` bindet die `-`-Zeile
   (Kandidaten-Header) **keine** Rolle → keine Trennung → `fx-t` bleibt eine Tabelle
   (byte-identisch zu v0.47.0). In `fx-s` bindet der verschluckte Header `| ID |
   Anforderung | Notiz |` die Rollen → Trennung → die relevante Tabelle wird erkannt.
   Dieselbe Struktur, unterschiedliche **Bindung** — das ist der Diskriminator, den
   die fünf rein strukturellen Fassungen nicht hatten.

3. **Die Invariante, positiv formuliert:** *Ein relevanter Header beendet die
   laufende Tabelle.* Damit wird **jede relevante Tabelle erkannt** — ihre
   Anforderungen können nicht mehr still in einer vorangehenden Tabelle verschwinden.
   Das adressiert den von [Review R3](../../reviews/2026-07-17-slice-074-implementation-r3.md)
   benannten Kern **direkt** (der Wiederaufsetz-Punkt des Header-Scans), statt
   Nachbarfälle zu verengen.

4. **Der Guard bleibt scharf.** Eine echt verrutschte Zeile (falsche Zellenzahl,
   kein neuer relevanter Header) bleibt Exit 2 ([ADR-0037](0037-trace-tabellenquellen-nullmengen-guard.md)).
   Die Grenze **ergänzt** die `badLine`-Regel, sie ersetzt sie nicht. Das
   Grenz-Prädikat feuert zudem **fail-closed auf den Bind-Fehler:** ein Header,
   dessen Rollen-Spalten mehrdeutig sind (z. B. eine doppelte Rollen-Spalte),
   bindet nicht sauber, beendet die laufende Tabelle aber dennoch — sonst würde
   eine Tabelle, die standalone Exit 2 wäre, hinter einer irrelevanten still
   verschluckt.

5. **Config-abhängig und deterministisch.** „Relevant" ist durch die konfigurierten
   Header-Namen definiert — dieselbe Wahrheit, die die Extraktion ohnehin leitet.
   Gleiche Config + gleiche Eingabe ⇒ gleiche Grenzen ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).

6. **Kein Lastenheft-Change-Request.** Das Lastenheft definiert keine Tabellengrenze;
   „wo eine Tabelle endet" ist Spezifikations-Sache (Rang 2, fortschreibbar).
   **SemVer-Minor:** d-check **findet mehr** — bisher verschluckte relevante Tabellen
   liefern Anforderungen (neue Waisen möglich). Ein grüner Konsumentenlauf kann
   danach rot sein, laut und in der sicheren Richtung.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Grenze am relevanten Header** (gewählt) | bricht die `fx-s`≡`fx-t`-Mehrdeutigkeit über die **Bindung**; jede relevante Tabelle wird erkannt; adressiert den Wiederaufsetz-Punkt direkt; über 10 Fixtures + volle Suite gemessen | die Grenze ist config-abhängig; eine irrelevante Tabelle absorbiert weiter benachbarten irrelevanten Inhalt (benigne — keine Anforderung geht verloren) |
| Grenze am **beliebigen** Header (`isNewTableHeader` unbedingt) | rein strukturell, config-frei | zerreißt `fx-t`: eine `-`-Datenzeile wird zum Header einer rollenlosen Tabelle, die laufende endet davor, echte Anforderungen verschwinden **still** (vor dem Code widerlegt) |
| Trennzeile im Body ⇒ Exit 2 | fängt den `fx-s`-Fall laut | trifft `fx-t` (heute grün ⇒ dann Exit 2); braucht ein Signal auch für irrelevante Tabellen, das `badLine` nicht liefert |
| Echter GFM-Parser | standard-treu | **schließt die Klasse nicht** (gemessen: goldmark sieht `fx-s` wie wir — Policy, nicht Grammatik) |
| Nichts ändern | kein Risiko | ein stiller Waisen-Verlust bleibt ausgeliefert; der Toleranz-Aufsatz aus [ADR-0040](0040-kommentar-suffix-in-tabellenzeilen.md) bleibt blockiert |

**Fitness-Funktion** (gegen das ausgelieferte v0.47.0-Image gemessen, Prototyp 2026-07-18):

- `fx-s` (irrelevant verschluckt relevant): v0.47.0 `total 1` → mit Regel `total 3`
  (F-1, F-2, F-3). Der Waisen-Verlust ist zu.
- `fx-t` (all-dashes-Datenzeile, Gegenprobe): `total 2` unverändert — **keine**
  Falsch-Trennung.
- `fx-t2` (Datenzeile wiederholt Header-Namen, gefolgt von all-dashes): kein Verlust.
- `fx-adj` (zwei benachbarte **relevante** Tabellen ohne Leerzeile): beide erkannt.
- `fx-m` (Aufsatz der Direktiven-Toleranz aus [ADR-0040](0040-kommentar-suffix-in-tabellenzeilen.md):
  tolerierte Direktiven-Datenzeile, dann relevante Folgetabelle): mit der Toleranz
  `total 3` — der zuvor stille Übersprung-Pfad ist zu.
- `fx-m3` (echte Verrutschung, 4. Zelle **kein** Kommentar): bleibt Exit 2.
- Jede Grenze ist per **Mutation** gepinnt: das Entfernen des relevanten-Header-Checks
  lässt `fx-s` wieder auf `total 1` kippen (mindestens ein Test kippt).
- Die volle bestehende Tabellen-Akzeptanzsuite bleibt grün.

## Konsequenzen

- **Positiv:** ein ausgelieferter stiller Waisen-Verlust ist zu; jede relevante
  Tabelle wird erkannt; die Regel gilt über den geteilten Reader für **beide**
  Konsumenten; und sie macht die Direktiven-Toleranz aus
  [ADR-0040](0040-kommentar-suffix-in-tabellenzeilen.md) **strukturell sicher** (die
  Folgetabelle kann nicht mehr verschluckt werden) — der Grund, warum diese Grenze
  **vor** jener Toleranz steht.
- **Negativ / Kosten:** die Grenze ist config-abhängig (kein rein struktureller
  Test). Eine irrelevante Tabelle absorbiert weiter benachbarten **irrelevanten**
  Inhalt ohne Leerzeile — benigne, weil per Definition keine Rolle gebunden und damit
  keine Anforderung betroffen ist; das Verhalten ist byte-identisch zu v0.47.0.
- **Verhaltensänderung für Bestandskonsumenten:** d-check **findet mehr** — eine
  bisher verschluckte relevante Tabelle liefert nun Anforderungen (neue Waisen
  möglich). SemVer-Minor; Release-Notiz und CHANGELOG müssen es benennen.
- **Bewusst offen gelassen:** die Grenze greift nur bei einem **relevanten** neuen
  Header. Zwei benachbarte **irrelevante** Tabellen bleiben konfliert (kein Realfall,
  kein Anforderungsverlust). Eine allgemeine „Leerzeile-zwischen-Tabellen-Pflicht"
  wäre eine breitere Policy ohne Beleg — nicht Teil dieser Entscheidung.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-07-18 | Proposed. Anlass: die Isolierungs-Messung `fx-s` während der Rücknahme von slice-074 (2026-07-17) zeigte den stillen Waisen-Verlust im ausgelieferten Image, ohne Marker. Nach fünf gescheiterten slice-074-Fassungen (alle rein strukturell, dreimal stilles Grün) wurde der relevant-Header-Diskriminator als Kandidat **gemessen** (Prototyp gegen v0.47.0, 10 Fixtures inkl. `fx-s`/`fx-t`/`fx-m` + volle Suite grün) — die erste Regel, die `fx-s` **und** `fx-t` gleichzeitig besteht. Umsetzender Slice slice-077; slice-074 setzt darauf auf, slice-071s Realdatenbeleg wird damit fahrbar. |
| 2026-07-18 | Accepted. slice-077 umgesetzt, reviewt (R1 ACCEPT-WITH-NITS, 0 HIGH/MEDIUM) und als **v0.48.0** freigegeben. Alle drei LOW-Nits eingearbeitet: R-F-1 (fx-adj zu einem echten, breitensensitiven Mutations-Pin gemacht), R-F-2 (Vorwärts-Cross-Durchreichung gepinnt), R-F-3 (das Grenz-Prädikat feuert **fail-closed auf den Bind-Fehler** — ein mehrdeutig-relevanter Header wird nicht mehr still verschluckt, sondern Exit 2; in §Entscheidung 4 aufgenommen). |
