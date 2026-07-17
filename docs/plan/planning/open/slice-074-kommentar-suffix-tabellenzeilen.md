# Slice slice-074: Direktiven-Zelle in Tabellenzeilen (Reader vs. eigene Direktive)

**Status:** open — **aus `in-progress/` zurückgestellt (2026-07-17), Implementierung
zurückgenommen** (`05e1889`, `806051f`). Der Reader steht wieder byte-identisch auf
v0.45.1. Der Defekt ist **unverändert offen und ausgeliefert**.

**Welle:** keine. Wartet auf erneute Einplanung, **nachdem die Klasse verstanden
ist** — nicht auf den sechsten Anlauf.

**Bezug:** **Defekt-Fix**, **kein Change Request**: das Lastenheft definiert keine
Zellenzahl — „was eine Zelle ist" ist Spezifikations-Sache
([`DC-FA-REQ-001.a`](../../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen);
wirkt über den geteilten Reader zugleich auf
[`DC-FA-XREF-001.a`](../../../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)).
Begründende Entscheidung
[ADR-0040](../../adr/0040-kommentar-suffix-in-tabellenzeilen.md) — bleibt
`Proposed`: die Frage ist gestellt, nicht beantwortet. **SemVer-Patch.**

**Autor:** pt9912. **Datum:** 2026-07-17.

---

## 1. Ziel

d-checks eigene Direktiven-Konvention `<!-- d-check:ignore (Grund) -->` **muss** in
einer Tabellenzeile hinter der letzten Pipe stehen — in einer Zelle wäre sie
Zellinhalt. Der header-gebundene Reader zählt sie dort als zusätzliche Zelle und
bricht fail-closed mit Exit 2 ab. **Die eigene Konvention macht den eigenen Reader
blind.** GFM ignoriert überzählige Zellen; die Zeile rendert überall normal.

Ausgeliefert seit v0.43.0 für `trace.requirements.format: table`. Der Slice muss
diese Zeile lesbar machen, **ohne** den Zellenzahl-Guard aufzuweichen
([ADR-0037](../../adr/0037-trace-tabellenquellen-nullmengen-guard.md)).

## 2. Entscheidungen / Regel

**Es gibt keine tragende Regel — das ist der Stand.** Fünf Fassungen haben dieselbe
Klasse fünfmal knapp verfehlt; drei davon erzeugten ein **stilles Grün**:

| Fassung | Ansatz | Wie sie fiel |
|---|---|---|
| `44a5201` | `dropCommentSuffix` im Splitter | **R1-F-1 HIGH:** der Splitter kennt Header ≠ Datenzeile nicht ⇒ Body-Regel am Header; dessen Zellenzahl fiel unter die Trennzeile, die Tabelle wurde **wortlos** übersprungen (Exit 1 ⇒ Exit 0) |
| `1210842` | tolerieren statt strippen, in `cellCountOK` | **R2-F-1:** die Toleranz fraß den Direktiven-**Header** der Folgetabelle (Exit 2 ⇒ Exit 0) |
| `e8b66ec` | Grenze am nächsten Header (`isNewTableHeader`) | **R3-F-1 HIGH:** das Spiegelbild — die tolerierte **Daten**zeile verschluckt die **ganze** Folgetabelle (Exit 1 ⇒ Exit 0, beide Konsumenten) |
| `670ebaf` | Lookahead auf die Nachsicht verengt (Selbstbefund) | **R3-F-2 MEDIUM:** beide neuen Grenzen **ungepinnt** — Rückdrehen bleibt grün (M-3), Panic-Guard entfernen bleibt grün (M-4) |
| Vorschlag: `isNewTableHeader` **unbedingt** | — | vor dem Code an `fx-t` widerlegt: eine Datenzeile aus lauter `---` würde zum Header einer rollenlosen Tabelle ⇒ stiller Verlust |

**Der strukturelle Kern**
([R3](../../../reviews/2026-07-17-slice-074-implementation-r3.md) F-1): In v0.45.1
setzt **jede** `badLine` den Header-Scan neu auf. **Die Toleranz entfernt genau
diesen Wiederaufsetz-Punkt.** Jede bisherige Fassung fragte „ist *diese* Zeile ein
Header?" statt „verschlucke ich durch das Tolerieren einen *nachfolgenden*?". Wer
die Invariante direkt adressiert, schließt die Klasse; wer weiter Nachbarfälle
verengt, bekommt R4.

## 3. Definition of Done

Zurückgesetzt — die Implementierung ist zurückgenommen. Der Doc-Teil trägt:

- [x] **ADR + Index:** [ADR-0040](../../adr/0040-kommentar-suffix-in-tabellenzeilen.md),
  `Proposed`, im Index. Prämisse zweimal korrigiert, Historie vollständig.
- [x] **Reviews als Beleg:** [R1](../../../reviews/2026-07-17-slice-074-implementation-r1.md),
  [R2](../../../reviews/2026-07-17-slice-074-implementation-r2.md),
  [R3](../../../reviews/2026-07-17-slice-074-implementation-r3.md) — sie sind das
  Wertvollste am Slice.
- [ ] **Tragende Regel:** benannt, die den Wiederaufsetz-Punkt **nicht** entfernt.
- [ ] **Spezifikation:** Schritte 3/5 — erst wenn die Regel steht (die Fassung von
  `1210842`/`e8b66ec` ist mit zurückgenommen).
- [ ] **Mutations-Härte:** **jede** neue Grenze kippt einen Test. Die Zusage aus
  `670ebaf` war für beide Grenzen falsch (R3-F-2) — der Sensor war blind an genau
  der Stelle, an der die Klasse zuletzt zuschlug.
- [ ] **Realdatenbeleg:** grid-gym `architecture.md:913`.
- [ ] **Release** + **unabhängiger, kontext-getrennter Review vor** dem Release.

## 4. Risiken / offene Punkte

- **Der Defekt bleibt ausgeliefert.** Die Rücknahme ist die ehrliche Zwischenlage,
  nicht die Lösung: eine Zeile, die jeder Renderer normal darstellt, ist weiter
  Exit 2. Sie blockiert den von
  [ADR-0038](../../adr/0038-trace-cross-consistency.md) Entscheidung 7 geforderten
  Realdatenbeleg von
  [`slice-071`](slice-071-trace-cross-consistency-gate.md).
- **Der stille Pfad ist älter als dieser Slice und braucht keinen Marker.** Belegt
  am **ausgelieferten** v0.45.1 (Fixture `fx-s`, 2026-07-17): eine irrelevante
  Tabelle gleicher Breite, ohne Leerzeile gefolgt von einer relevanten, verschluckt
  deren Anforderungen lautlos — `1 Anforderung(en), 0 Waise(n)`, Exit 0, während
  zwei echte Waisen existieren. R3-F-1 ist die **Verbreiterung** eines bestehenden
  Lochs, nicht dessen Ursprung; v0.45.1s korrektes Verhalten an `fx-m` war Zufall,
  weil die kaputte Zeile den Scan versehentlich rettete. **Eigener Defekt, eigener
  Slice.**
- **Marker auf der Vorzeile ist gemessen tot** (Fixture `fx-p`, beide Images): eine
  eigenständige Kommentarzeile trägt keine Pipe, der Reader beendet die Tabelle,
  die restlichen Anforderungen verschwinden **lautlos**. goldmark liest sie genauso
  (1 statt 3 Datenzeilen) — kein d-check-Artefakt, sondern GFM. Zudem prüfen alle
  Ignore-Konsumenten (`ids`, `codepaths`, `versions`) den Marker **auf der Zeile
  selbst**; eine Vorzeilen-Semantik wäre ein neuer Vertrag in jedem Modul.
- **Ein echter GFM-Parser schließt diese Klasse nicht** (Spike 2026-07-17, goldmark
  v1.8.4 gegen 522 reale Dateien): auf `fx-s` und `fx-p` stimmt goldmark **exakt**
  mit dem heutigen Reader überein. Die Grenze ist **Policy**, nicht Grammatik — GFM
  gibt uns dort recht. Für `fx-913` verwirft goldmark die überzählige Zelle
  stillschweigend und nähme uns den Guard aus
  [ADR-0037](../../adr/0037-trace-tabellenquellen-nullmengen-guard.md). Der Spike fand allerdings
  **zwei echte Grammatik-Defekte** an anderer Stelle ⇒ slice-076.
- **Dogfood-Lücke:** d-check nutzt in den eigenen Doku-Tabellen keinen
  Trailing-Marker und konnte den Defekt an sich selbst nicht bemerken — dieselbe
  Klasse wie die Range-Notation (slice-073 §4).

## 5. Trigger

Realdatenbeleg gegen grid-gym (2026-07-17) — der von
[ADR-0038](../../adr/0038-trace-cross-consistency.md) Entscheidung 7 geforderte
Beleg, der den Generator freigibt. v0.45.1 brach an `spec/architecture.md:913` ab:
`Tabellenzeile 913 hat 4 statt 3 Zellen`. Die Zeile ist eine legale 3-Spalten-Zeile
mit d-checks eigenem Ignore-Marker dahinter. Dritter Defekt, den erst die Realdaten
zeigten — die ersten beiden (`forward.req-pattern`, Link-Ranges) lagen in den
Mustern, dieser in der Grammatik.

## 6. Sub-Area-Modus-Begründung

GF (Repo-Default): Die Spezifikation führt, der Code folgt. Der Reader ist
bestehender, spezifizierter Code; die Kompatibilität der markerlosen Formen ist
durch die Akzeptanztests aus
[slice-070](../done/slice-070-trace-tabellenquellen-nullmengen-guard.md) geschützt.

## 7. Closure-Notiz (nach `done/`)

_Ausstehend — der Slice ist nicht abgeschlossen, sondern zurückgestellt. Die
Rücknahme ist in `05e1889` (Code/Spec) und `806051f` (Handbuch) dokumentiert;
Anlass, Verlauf und die fünf Fehlschläge stehen in §2 und in den drei Reviews._
