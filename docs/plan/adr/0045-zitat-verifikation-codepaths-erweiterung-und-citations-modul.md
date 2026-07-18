# ADR-0045 — Zitat-Verifikation: `codepaths`-Zeilen-Check erweitern, `verbatim` als eigenes Modul `citations`

**Status:** Proposed
**Datum:** 2026-07-18
**Autor:** pt9912
**Schärft:** [`DC-FA-CITE-001.a`](../../../spec/spezifikation.md#dc-fa-cite-001a--verbatim-zitat-verifikation-citations)
**Bezug:** [`DC-FA-CITE-001`](../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in), [`DC-FA-CODE-001`](../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in), [ADR-0040](0040-kommentar-suffix-in-tabellenzeilen.md), [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus), [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).

## Kontext

Ein Adopter (`ai-harness-init`, CR 2026-07-17) führt seine Baseline **committet
vendored**; damit entsteht ein Korpus von `datei:zeile`-**Zitaten** auf einen
in-tree, versionierten Fremdbaum. Diese Zitate sind heute korrekt, zeigen aber beim
nächsten Tag-Bump **still** ins Leere — der Adopter maß neun Stunden zwischen zwei
Tags und verschobene Zeilen in 35 von 42 Dateien. d-check verifiziert **Struktur**
(Links lösen auf, Anker existieren), aber **keine Behauptung**: eine `datei:zeile`
kann falsch sein, und jeder Gate-Lauf bleibt grün.

Zwei Befunde erden den Zuschnitt:

- **Die Zitate stehen in Inline-Code** (gemessen: 33 von 33 `datei:zeile`-Zitate im
  Adopter-Repo backticked, null in nackter Prosa). Das Modul `codepaths` erkennt die
  `datei:zeile`-Konvention **bereits** und **verwirft** die Zeilennummer bewusst
  ([`DC-FA-CODE-001.a`](../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code)
  Schritt 3). Die Stufen 1/2 sind damit **das Weglassen eines `trim`** plus
  Zeilenzahl-Prüfung — keine neue Detektion.
- Der Adopter hat seinen **eigenen** Sensor-Slice deshalb selbst blockiert (Prämisse
  widerlegt: d-check erkennt die Zeile schon).

Der CR beschreibt drei Stufen: **1 `exists`** (Ziel existiert mit ≥ `bis` Zeilen),
**2 `in-range`** (`von ≤ bis`), **3 `verbatim`** (ein ausgezeichneter Zitatblock wird
zeichengenau gegen die Quell-Spanne geprüft).

## Entscheidung

1. **Stufe 1/2 erweitern `codepaths`, kein neues Modul.** Der Zeilen-Check läuft auf
   **denselben** Inline-Code-Pfaden, die `codepaths` schon detektiert; er ist eine
   Einzelmodul-Erweiterung → Änderung von
   [`DC-FA-CODE-001`](../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
   (opt-in `codepaths.check-lines`, Default aus **byte-identisch** — das Suffix wird
   sonst wie bisher verworfen). Das folgt dem Kürzel-Kriterium: *querschnittlich →
   neues Kürzel; Einzelmodul-Erweiterung → bestehende Anforderung ändern.* Neue
   Grund-Codes `citation-out-of-range` (Zeile hinter Datei-Ende) und
   `citation-inverted-range` (`von > bis`).

2. **Stufe 3 (`verbatim`) als eigenes Modul `citations` — die zentrale Abgrenzung.**
   Erweiterung *oder* eigenes Modul? Entscheidend ist die **Art der Prüfung**:

   - Stufe 1/2 prüfen **Existenz/Bereich** (zeigt die Zeile ins Leere?) — das ist
     `codepaths`' Aufgabe (Pfade existieren), nur zeilen-genau.
   - Stufe 3 prüft **Inhalt** (stimmt der zitierte Text mit der Quelle überein?) —
     eine **andere** Aufgabe. Sie ist zudem **direktiven-getrieben**
     (`d-check:cite` markiert den zu prüfenden Zitatblock), nicht scan-getrieben wie
     `codepaths`' Inline-Code-Suche.

   Dieselbe Trennung wie `anchors` ↔ `links`: `anchors` greift die Ziel-Auflösung von
   `links` auf, prüft aber eine **andere** Eigenschaft und ist ein eigenes Modul.
   `citations` greift `codepaths`' `datei:zeile`-Erkennung auf (**kein** zweiter
   Detektor — genau die Duplikation, die ein re-detektierendes Modul brächte, wird so
   vermieden) und bekommt als **neues Modul** ein eigenes Bereichskürzel (`CITE`) und
   diese ADR (d-check-Konvention: jedes Modul eine ADR).

3. **`d-check:cite`-Direktive.** Eine **zweite** Direktive neben `d-check:ignore`:
   ein HTML-Kommentar `<!-- d-check:cite <pfad>:<von>-<bis> -->` unmittelbar vor einem
   `>`-Zitatblock. Die Platzierungs-/Zeilen-Regeln der Direktiven-Klasse sind mit
   [ADR-0040](0040-kommentar-suffix-in-tabellenzeilen.md) gefestigt (der
   Tabellenzeilen-Fall der ersten Direktive war fünfmal gescheitert und ist jetzt
   gelöst) — die Voraussetzung dafür, überhaupt eine zweite Direktive einzuführen.

4. **Zeichengenauer Vergleich, keine Normalisierung.** Der Wert des Features ist,
   „wortgleich" von einer **Zusage** zu einem **gemessenen Property** zu machen; jede
   Whitespace-/Formatierungs-Normalisierung würde genau das aufweichen.

5. **Fail-closed.** Fehlende Zieldatei, Spanne über das Datei-Ende oder ungültiger
   Bereich ⇒ Exit 2 — kein stiller Nicht-Vergleich; eine kaputte `d-check:cite`-
   Direktive ist ein Autoren-Fehler, kein Schweigen.

6. **Nur ausgezeichnete Blöcke.** `citations` prüft ausschließlich `d-check:cite`-
   markierte Zitate — **kein** Prosa-Scanning. Freie Zahlen mit externer
   Grundwahrheit („42 Dateien im ZIP") und Prosa-Quantoren bleiben Review-Territorium.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **1/2 in `codepaths`, 3 eigenes Modul `citations`** (gewählt) | Existenz/Bereich bei `codepaths`, Inhalt separat (wie `anchors`↔`links`); reused Detektion; ein Kürzel-Kriterium konsistent | zwei Config-Flächen; zweite Direktive |
| Alle drei in `codepaths` | ein Modul | `verbatim` ist Inhalts- statt Existenz-Prüfung + direktiven-getrieben — Scope-Bruch für „Pfade existieren" |
| Eigenes Modul für alle drei, mit eigener Detektion | klare Grenze | **dupliziert** `codepaths`' `datei:zeile`-Detektion, `roots`, Ventile |
| Prosa-Scanning statt Direktive | fängt Zitate ohne Markup | 33/33 sind Inline-Code; Prosa-Scan ist deutlich fehleranfälliger |
| `claims` (Provenienz-Pflicht für Zahlen) | fängt freie Zahlen | höheres False-Positive-Risiko; sollte auf Betriebserfahrung aufsetzen, nicht darauf wetten |

**Fitness-Funktion:**

- Ein `datei:zeile`-Zitat hinter dem Datei-Ende ⇒ `citation-out-of-range`; ein
  invertierter Bereich ⇒ `citation-inverted-range`; ein per `d-check:cite` markierter,
  vom Quelltext abweichender Zitatblock ⇒ `citation-mismatch`.
- Ein korrektes Zitat bleibt **grün** nach einem Tag-Bump, der die zitierte Datei
  nicht anfasst (kein Fehlalarm durch Nachbar-Drift).
- Ohne `check-lines` bzw. ohne aktives `citations`-Modul jeder Befundsatz
  byte-identisch ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)),
  hermetisch ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- Realdatenbeleg gegen das Adopter-Repo (Dogfood-Reichweite bei uns nur ~7) ist
  **nicht optional**.

## Konsequenzen

- **Positiv:** ein bislang ungefangener Drift-Typ (Zitat-Fäule) wird beim Verfassen
  rot; „wortgleich" ist gemessen; `citations` reused `codepaths`' Detektion; opt-in
  und byte-identisch für Nicht-Nutzer.
- **Negativ / Kosten:** die Config-Fläche wächst (`codepaths.check-lines` + neues
  Modul); die zweite Direktive `d-check:cite` vergrößert die Direktiven-Platzierungs-
  Klasse (durch ADR-0040 beherrschbar); die Dogfood-Reichweite ist gering — der
  Realdatenbeleg trägt die Absicherung.
- **Verworfen:** alles in `codepaths`, ein re-detektierendes Modul, Prosa-Scanning,
  `claims` (jeweils oben begründet).

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-07-18 | Proposed. Change Request Adopter `ai-harness-init` (2026-07-17). §4-Vorfragen entschieden: Adopter-Rückfrage empirisch (33/33 Zitate in Inline-Code ⇒ `codepaths`-Erweiterung); Zuschnitt Form (c) vom Auftraggeber. Umsetzender Slice slice-079. |
