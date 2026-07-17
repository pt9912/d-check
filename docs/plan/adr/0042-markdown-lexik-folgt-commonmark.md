# ADR-0042 — Markdown-Lexik folgt CommonMark/GFM: Trennzeile und Fence-Infozeile

**Status:** Proposed
**Datum:** 2026-07-17
**Autor:** pt9912
**Schärft:** [`DC-FA-REQ-001.a`](../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen) (Trennzeile), [`DC-FA-LINK-001.a`](../../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion) (Fence)
**Bezug:** [`DC-FA-REQ-001`](../../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen), [`DC-FA-LINK-001`](../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links), [`DC-FA-XREF-001`](../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in), [ADR-0037](0037-trace-tabellenquellen-nullmengen-guard.md), [ADR-0040](0040-kommentar-suffix-in-tabellenzeilen.md), [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus).

## Kontext

d-check schreibt seine Markdown-Lexik von Hand. Das ist tragfähig, solange die
handgeschriebene Regel dem entspricht, was jeder Renderer tut — sonst ist d-check
an einer Stelle blind, die für den Autor **normal aussieht**. Genau das ist der
Fall, und zwar **still**: eine Regel, die eine Struktur nicht erkennt, meldet
nichts. Sie meldet **weniger**.

Belegt durch einen Differential-Spike (2026-07-17): [goldmark](https://github.com/yuin/goldmark)
v1.8.4 gegen den heutigen Reader über **522 reale Dateien** (d-checks eigene Doku
+ eine Kopie von grid-gyms `spec/`+`docs/`), **490 erkannte Tabellen** ⇒ **8
Abweichungen — alle in dieselbe Richtung: goldmark sieht Tabellen, die d-check
nicht sieht.** Zwei Ursachen-Familien:

**A · Die Trennzelle verlangt drei Bindestriche, GFM verlangt einen.**
`^:?-{3,}:?$` ([`DC-FA-REQ-001.a`](../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 3). Jede reale Tabelle mit `| -- |` oder `| - |` in der Trennzeile ist
für d-check **keine Tabelle**. Sechs der acht Abweichungen; darunter zwei
Tabellen mit dem Header `| ID | Lastenheft-Anforderung | Substanz-Pfad |
Test-Pfad | Status |` — Anforderungs-Gestalt. Das ist **kein Code-Bug**: die
Spezifikation schreibt `-{3,}` vor. Die **Spec** weicht von GFM ab.

**B · Der Fence-Automat kennt die Infozeilen-Regel nicht.**
[`DC-FA-LINK-001.a`](../../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
Schritt 1: „Zeilen, deren erste Nicht-Leerzeichen-Folge mit ` ``` ` oder `~~~`
beginnt, schalten den Fence-Zustand um". CommonMark verlangt zusätzlich: **die
Infozeile einer Backtick-Fence darf keine Backticks enthalten** — sonst ist die
Zeile kein Fence-Öffner, sondern Fließtext. In d-checks eigener Review-Doku
(`docs/reviews/`, 2026-06-19) steht Zeile 179 genau so — ein Satz **über** einen
Fence, der mit ```` ```yaml-Fence (`datei.md`) — … ```` beginnt. d-check hält ihn
für einen Öffner und markiert **den ganzen Rest der Datei als Nicht-Prosa**.

Die Tragweite von B ist **nicht** auf Tabellen begrenzt: Schritt 1 sagt „Zeilen
im Fence-Zustand werden von **allen Modulen** ignoriert". Gemessen an einer
Fixture gegen das Image: ein kaputter Link **nach** einer solchen Zeile wird
still verschluckt — `1 Befund(e)`/Exit 1 **ohne** die Zeile, `0 Befund(e)`/Exit 0
**mit** ihr. Ein Stilles-Grün-Pfad quer durch alle Module.

Umfang, ehrlich gemessen: **eine** Fundstelle in beiden Repos (d-checks eigene).
Der Defekt ist selten, aber er ist real, still und trifft jedes Modul; und er
skaliert mit genau der Prosa, die Harness-Repos schreiben — Text **über** Fences.

## Entscheidung

1. **Die Trennzelle folgt GFM:** `^:?-+:?$` — ein Bindestrich genügt, optionale
   Doppelpunkte bleiben. Rein erweiternd: jede heute erkannte Trennzeile bleibt
   erkannt.

2. **Ein Backtick-Fence-Öffner, dessen Infozeile einen Backtick enthält, ist kein
   Fence.** Er ist Fließtext und wird von den Modulen normal gelesen. Für
   `~~~`-Fences gilt die Regel nicht (CommonMark erlaubt dort Backticks in der
   Infozeile) — die Asymmetrie ist CommonMark, keine d-check-Eigenheit.

3. **Die Grenze ist das Gemessene.** Es werden **keine** weiteren
   CommonMark-Angleichungen auf Verdacht vorgenommen. Bekannt und **bewusst
   offen** gelassen (siehe §Konsequenzen): der Zeichen-/Längen-Abgleich beim
   Fence-Schluss, die 4-Leerzeichen-Einrückungsgrenze, und die **Divergenz der
   beiden Fence-Automaten** im eigenen Code. Jede davon braucht ihren eigenen
   Beleg, bevor sie ihre eigene Regel bekommt. Das ist die Lehre aus der zurückgenommenen Direktiven-Toleranz ([ADR-0040](0040-kommentar-suffix-in-tabellenzeilen.md)):
   fünf Fassungen, die eine Regel knapp neben dem belegten Problem platzierten.

4. **Kein Parser.** Die Alternative ist gemessen und in
   [ADR-0040](0040-kommentar-suffix-in-tabellenzeilen.md) `## Geschichte`
   protokolliert: goldmark stimmt auf den Policy-Fixtures exakt mit dem heutigen
   Reader überein und verwürfe für die Direktiven-Zelle die überzählige Zelle
   **still** — er nähme uns den Guard aus
   [ADR-0037](0037-trace-tabellenquellen-nullmengen-guard.md). Ein Parser
   beantwortet die **Grammatik**, nicht die **Policy**; hier reichen zwei
   Grammatik-Regeln, für die er nicht gebraucht wird.

5. **Kein Lastenheft-Change-Request.** Beide Regeln sind Spezifikations-Sache
   (Rang 2, fortschreibbar): das Lastenheft sagt weder, was eine Trennzeile ist,
   noch was einen Fence öffnet.

6. **SemVer-Minor.** Anders als bei einem Defekt-Fix, der Falschbefunde entfernt,
   **findet** d-check danach **mehr**: bisher unsichtbare Tabellen liefern
   Anforderungen, bisher unsichtbare Prosa liefert Links. Ein Konsumentenlauf,
   der heute grün ist, kann danach rot sein — laut, nicht still, und in der
   sicheren Richtung. Das ist kein Patch, und die Release-Notiz muss es sagen.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Zwei gezielte Regeln, gegen Realdaten belegt** (gewählt) | schließt 7 der 8 gemessenen Abweichungen; die Trennzeile ist ein Zeichen; kein neues Vokabular, keine Abhängigkeit; jede Regel hat eine reale Fundstelle | die Lexik bleibt handgeschrieben — die nächste Familie ist wahrscheinlich, nur noch nicht belegt |
| goldmark-Hybrid (Struktur vom Parser, Zellenzahl auf der Rohzeile) | schlösse die Grammatik-Klasse **als Klasse**; upstream-getestet; null transitive Abhängigkeiten | Portierung von `markdownTables`/`tableHeaderAt`/`markdownTableLines`; Masken-Semantik (`sections`) neu zu klären; löst **keine** der beiden hier belegten Familien besser als zwei Zeilen es tun |
| Reiner goldmark-Port (GFM-Semantik voll) | standard-treu | verwürfe überzählige Zellen **still** ⇒ Guard aus ADR-0037 weg — genau die Klasse, gegen die er gebaut wurde |
| Nichts ändern | kein Risiko, mehr zu finden | zwei belegte Stilles-Grün-Pfade bleiben; einer davon blendet **alle** Module ab der Fundstelle |
| Alle Divergenzen zu CommonMark auf einmal angleichen | „einmal richtig" | keine Belege für die übrigen; fünf Fassungen der Direktiven-Toleranz haben gezeigt, was eine Regel ohne Beleg anrichtet ([ADR-0040](0040-kommentar-suffix-in-tabellenzeilen.md)) |

**Fitness-Funktion:**

- Eine Tabelle, deren Trennzeile `| -- |` oder `| - |` trägt, wird erkannt; ihre
  Anforderungen erscheinen in der RTM. Belegt an der realen `carveouts.md`-Form.
- Jede heute erkannte Trennzeile (`---` und länger, mit/ohne Doppelpunkte) bleibt
  erkannt — die Änderung ist rein erweiternd.
- Eine Zeile, die mit ` ``` ` beginnt und einen Backtick in der Infozeile trägt,
  öffnet **keinen** Fence: ein kaputter Link dahinter wird gemeldet (Exit 1).
  Belegt an der realen Fließtext-Form ```` ```yaml-Fence (`datei.md`) — … ```` in
  d-checks eigener Review-Doku (2026-06-19, Zeile 179).
- Ein **echter** Fence (` ```yaml `, ` ``` `, `~~~`) verdeckt seinen Inhalt
  unverändert — kein Modul liest plötzlich Fence-Inhalt.
- Beide Regeln sind per **Mutation** gepinnt: wird die Trennzeilen-Lockerung
  zurückgedreht oder die Infozeilen-Regel entfernt, kippt je mindestens ein Test.
  (Die Suite war gegen **beide** blind — die Lockerung ließ sie grün.)
- Der Differential gegen goldmark über die 522 Realdateien fällt von 8 auf 2
  Abweichungen; die verbleibenden zwei sind **Policy**, nicht Grammatik
  (überzählige Zelle: der [ADR-0037](0037-trace-tabellenquellen-nullmengen-guard.md)-Guard, by design).

## Konsequenzen

- **Positiv:** zwei belegte Stilles-Grün-Pfade sind zu. d-check sieht neun
  Tabellen mehr allein in den 522 Messdateien; die Fence-Regel wirkt über die
  geteilte Vorverarbeitung auf **alle** Module.
- **Negativ / Kosten:** die Lexik bleibt handgeschrieben, und dieses ADR sagt
  ausdrücklich, dass das eine offene Wette ist. Der Spike hat zwei Familien
  gefunden, weil er zum ersten Mal gegen einen echten Parser gemessen hat — ein
  Differential-Sensor gegen goldmark (nur im Test, ohne Produktiv-Abhängigkeit)
  wäre der Sensor, der die dritte Familie findet, bevor ein Konsument sie findet.
  **Offener Punkt, kein Teil dieser Entscheidung.**
- **Bewusst offen gelassen:** (a) `proseLines`
  (`internal/hexagon/core/rules/markdown.go`) ist ein **naiver Toggle** — jede
  ```/`~~~`-Zeile kippt den Zustand, ohne Zeichen- oder Längenabgleich;
  `markdownTableLines` (`internal/hexagon/core/app/trace_table.go`) prüft beides.
  **Zwei Automaten, zwei Verhalten**, und die Spec beschreibt den naiven ⇒ der
  Tabellen-Reader weicht heute schon von seiner eigenen Spec ab. (b) Die
  4-Leerzeichen-Einrückungsgrenze (`    ```" ` ist Code, kein Fence). Beide sind
  **unbelegt** — kein Realfall in den 522 Dateien — und bekommen erst eine Regel,
  wenn einer existiert.
- **Verhaltensänderung für Bestandskonsumenten:** d-check **findet mehr**. Eine
  bisher unsichtbare Tabelle kann neue Waisen liefern, bisher unsichtbare Prosa
  neue Link-/ID-Befunde. Ein grüner Lauf kann rot werden — das ist der Zweck.
  Release-Notiz und CHANGELOG müssen es benennen; SemVer-Minor.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-07-17 | Proposed. Anlass: Differential-Spike goldmark v1.8.4 gegen 522 reale Dateien (490 Tabellen ⇒ 8 Abweichungen, alle „d-check ist blind"), gefahren auf die Frage des Auftraggebers „brauchen wir einen besseren Markdown-Parser?" während der Rücknahme von slice-074. Die Antwort ist zweigeteilt und gemessen: für die **Policy** (slice-074, slice-077) nein — goldmark stimmt dort exakt mit dem heutigen Reader überein; für die **Grammatik** ja — er hat zwei ausgelieferte, stille Familien gefunden, die drei Reviews und fünf Implementierungsanläufe nicht gesehen haben. Dieses ADR nimmt die zwei belegten Regeln und **nicht** den Parser. Umsetzender Slice slice-076 |
