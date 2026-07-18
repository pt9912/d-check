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

   Dieselbe Modul-Grenze wie `anchors` ↔ `links`: `anchors` hat eine **eigene**
   Erkennung (Links mit Fragment), teilt aber die **Pfad-Auflösung** von `links` und
   prüft eine **andere** Eigenschaft. `citations` hat ebenso eine **eigene** Erkennung
   (es parst die `d-check:cite`-Direktive — ein eigener Detektor) und teilt die
   **Pfad-/Zeilen-Auflösung**; es re-implementiert **nicht** die
   Inline-Code-Pfad-Detektion von `codepaths`. Als **neues Modul** bekommt es ein
   eigenes Bereichskürzel (`CITE`) und diese ADR (d-check-Konvention: jedes Modul eine
   ADR).

3. **`d-check:cite`-Direktive markiert inline- oder Block-Zitate.** Eine **zweite**
   Direktive neben `d-check:ignore`: `<!-- d-check:cite <pfad>:<von>-<bis> -->`
   unmittelbar vor dem zu prüfenden Zitat — einem `>`-Blockquote **oder** dem folgenden
   inline-Zitat-Span (`„…"` / `"…"`). Die realen Adopter-Zitate sind **inline**, nicht
   `>`-Blöcke — ein reiner Block-Bezug träfe sie nicht. Die Platzierungs-/Zeilen-Regeln
   der Direktiven-Klasse sind mit
   [ADR-0040](0040-kommentar-suffix-in-tabellenzeilen.md) gefestigt — die Voraussetzung
   dafür, überhaupt eine zweite Direktive einzuführen.

4. **Vergleich: whitespace-normalisierter Teilstring, nicht zeilenweise-exakt.**
   „wortgleich" soll von einer **Zusage** zu einem **gemessenen Property** werden —
   aber die realen Zitate sind **inline** und beim Verfassen **umgebrochen** (andere
   Spalten als die Quelle), beginnen und enden **mitten in einer Zeile**. Ein
   zeilenweise zeichengenauer Vergleich gegen die vollen Quell-Zeilen träfe genau das
   **korrekte** Zitat nicht (belegt am einen realen „wortgleich"-Zitat des Adopters).
   Daher: Quell-Spanne und Zitattext werden **whitespace-normalisiert** (Läufe aus
   Leerzeichen/Zeilenumbruch → ein Leerzeichen, getrimmt), und das normalisierte Zitat
   muss ein **zusammenhängender Teilstring** der normalisierten Quelle sein. So bleibt
   der Vergleich **zeichengenau auf dem Nicht-Whitespace-Inhalt**, tolerant nur
   gegenüber Umbruch/Whitespace: jede echte Wort-Abweichung bricht den Teilstring
   (`citation-mismatch`), Re-Wrapping nicht. Gemessen am realen Zitat: korrekt = grün,
   ein gedriftetes Wort = rot.

5. **Fail-closed.** Fehlende Zieldatei, Spanne über das Datei-Ende oder ungültiger
   Bereich ⇒ Exit 2 — kein stiller Nicht-Vergleich; eine kaputte `d-check:cite`-
   Direktive ist ein Autoren-Fehler, kein Schweigen.

6. **Nur ausgezeichnete Blöcke.** `citations` prüft ausschließlich `d-check:cite`-
   markierte Zitate — **kein** Prosa-Scanning. Freie Zahlen mit externer
   Grundwahrheit („42 Dateien im ZIP") und Prosa-Quantoren bleiben Review-Territorium.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **1/2 in `codepaths`, 3 eigenes Modul `citations`** (gewählt) | Existenz/Bereich bei `codepaths`, Inhalt separat (wie `anchors`↔`links`); teilt Pfad-/Zeilen-Auflösung; ein Kürzel-Kriterium konsistent | zwei Config-Flächen; zweite Direktive; Stufe-3-Substrat (Direktiven) muss der Adopter erst schaffen |
| Alle drei in `codepaths` | ein Modul | `verbatim` ist Inhalts- statt Existenz-Prüfung + direktiven-getrieben — Scope-Bruch für „Pfade existieren" |
| Eigenes Modul für alle drei, mit eigener Detektion | klare Grenze | **dupliziert** `codepaths`' `datei:zeile`-Detektion, `roots`, Ventile |
| Prosa-Scanning statt Direktive | fängt Zitate ohne Markup | 33/33 sind Inline-Code; Prosa-Scan ist deutlich fehleranfälliger |
| `claims` (Provenienz-Pflicht für Zahlen) | fängt freie Zahlen | höheres False-Positive-Risiko; sollte auf Betriebserfahrung aufsetzen, nicht darauf wetten |

**Fitness-Funktion:**

- Ein `datei:zeile`-Zitat hinter dem Datei-Ende ⇒ `citation-out-of-range`; ein
  invertierter Bereich ⇒ `citation-inverted-range`; ein per `d-check:cite` markiertes
  Zitat (inline oder Block), dessen normalisierter Text **kein** Teilstring der
  normalisierten Quell-Spanne ist ⇒ `citation-mismatch`.
- Ein korrektes Zitat bleibt **grün** nach einem Tag-Bump, der die zitierte Datei
  nicht anfasst (kein Fehlalarm durch Nachbar-Drift).
- Ohne `check-lines` bzw. ohne aktives `citations`-Modul jeder Befundsatz
  byte-identisch ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)),
  hermetisch ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- Realdatenbeleg gegen das Adopter-Repo (Dogfood-Reichweite bei uns nur ~7) ist
  **nicht optional**.

## Konsequenzen

- **Positiv:** ein bislang ungefangener Drift-Typ (Zitat-Fäule) wird beim Verfassen
  rot; „wortgleich" ist gemessen; `citations` teilt die Pfad-/Zeilen-Auflösung von
  `codepaths`/`links`; opt-in und byte-identisch für Nicht-Nutzer.
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
| 2026-07-18 | **Doc-first-Design-Review R1** (BLOCK scoped auf Stufe 3) eingearbeitet, Status weiterhin `Proposed`: das Verbatim-Modell von „zeilenweise zeichengenau gegen `>`-Blöcke" auf **whitespace-normalisierten Teilstring** (inline/re-wrapped/Teilzeilen-tolerant) umgestellt und am **realen** Adopter-Zitat validiert (korrekt = grün, ein gedriftetes Wort = rot); die Direktive markiert nun inline **oder** Block-Zitate; das Reuse-Argument (F-5) auf „eigener Detektor, geteilte Pfad-Auflösung" korrigiert. Stufe 1/2 waren geerdet. Offen ausgewiesen: das Stufe-3-Substrat (Direktiven) muss der Adopter erst schaffen. |
