# Slice slice-074: Kommentar-Suffix in Tabellenzeilen (Reader vs. eigene Direktive)

**Status:** in-progress (welle-60-trace-cross-consistency).

**Welle:** aktiv (welle-60), **vorrangig** — der Defekt blockiert den
Realdatenbeleg von [`slice-071`](slice-071-trace-cross-consistency-gate.md)
vollständig (Exit 2 beim Lesen der Rück-Sicht).

**Bezug:** **Defekt-Fix**, **kein Change Request**: das Lastenheft definiert keine
Zellenzahl — „was eine Zelle ist" ist Spezifikations-Sache und dort geschärft
([`DC-FA-REQ-001.a`](../../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 3; wirkt über den geteilten Reader zugleich auf
[`DC-FA-XREF-001.a`](../../../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)).
Begründende Entscheidung
[ADR-0040](../../adr/0040-kommentar-suffix-in-tabellenzeilen.md) (Proposed).
**SemVer-Patch.**

**Autor:** pt9912. **Datum:** 2026-07-17.

---

## 1. Ziel

d-checks eigene Direktiven-Konvention `<!-- d-check:ignore (Grund) -->` **muss** in
einer Tabellenzeile hinter der letzten Pipe stehen — in einer Zelle wäre sie
Zellinhalt. Der header-gebundene Reader zählt sie dort als zusätzliche Zelle und
bricht fail-closed mit Exit 2 ab. **Die eigene Konvention macht den eigenen Reader
blind.** GFM ignoriert überzählige Zellen; die Zeile rendert überall normal.

Der Slice macht den abschließenden Kommentar zu einem Zeilen-**Suffix** statt zu
einer Zelle — und lässt den Zellenzahl-Guard sonst unangetastet.

## 2. Entscheidungen / Regel

- **Genau einer, nur am Ende, nur ganzzellig.** Ein Kommentar innerhalb einer Zelle
  oder in der Zeilenmitte bleibt Zellinhalt
  ([ADR-0040](../../adr/0040-kommentar-suffix-in-tabellenzeilen.md)).
- **Der Guard bleibt scharf:** nach dem Entfernen des Suffixes ist eine abweichende
  Zellenzahl unverändert Exit 2. Die Ausnahme ist die eigene Direktive, keine
  allgemeine Aufweichung — die GFM-weite Lösung wurde bewusst verworfen, weil sie
  genau die Klasse still machte, gegen die [ADR-0037](../../adr/0037-trace-tabellenquellen-nullmengen-guard.md) den Guard gebaut hat.
- **Ein Fix, zwei Konsumenten** (geteilter Reader): `trace.requirements.format:
  table` (ausgeliefert seit v0.43.0) und `trace.cross-consistency`.

## 3. Definition of Done

- [x] **Spezifikation:** Kommentar-Suffix in
  [`DC-FA-REQ-001.a`](../../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
  Schritt 3, inkl. der negativen Abgrenzung (Kommentar in der Zelle) + Historie.
- [x] **ADR + Index:** [ADR-0040](../../adr/0040-kommentar-suffix-in-tabellenzeilen.md),
  Status Proposed, im Index.
- [ ] **Implementierung:** im geteilten Zeilen-Splitter, vor dem Zellen-Zählen —
  Header-, Trenn- und Datenzeilen gleichermaßen.
- [ ] **Tests (positiv):** Zeile mit `<!-- d-check:ignore (…) -->` wird gelesen wie
  ohne Marker — je einmal für `trace.requirements.format: table` und
  `trace.cross-consistency`, auf **Konsumenten-Ebene** (nicht nur am Splitter).
- [ ] **Tests (negativ):** echt verrutschte Zeile ohne Kommentar bleibt Exit 2;
  Kommentar **in** einer Zelle bleibt Zellinhalt; zwei Trailing-Kommentare
  expandieren nicht still.
- [ ] **Mutations-Härte:** Suffix-Abstreifung entfernt kippt genau einen Test; die
  „nur am Ende"-Grenze ebenso.
- [ ] **Realdatenbeleg:** der Lauf gegen grid-gyms echte `spec/architecture.md`
  bricht nicht mehr an Zeile 913 ab.
- [ ] **Nutzerdoku:** Handbuch (§5 Tabellen-Grammatik) + CHANGELOG (Fixed).
- [ ] **Release:** v0.45.2, Release-Prep + Tag + GHCR + Digest-Backfill.
- [ ] **Qualität:** unabhängiger, kontext-getrennter Review **vor** dem Release
  (die slice-073-Lehre: ein Fix, der Befunde entfernen kann, verdient ihn zuerst);
  `make gates`/`make ci` grün.

## 4. Risiken / offene Punkte

- **Der Splitter trägt ein weiteres Stück Markdown-Wissen** (nach Code-Span und
  Escape nun der Kommentar). Jede Regel dort ist eine Kopplung an die
  Markdown-Grammatik.
- **„Genau einer, nur am Ende" ist gesetzt.** Durch die Konvention gedeckt, aber ein
  Konsument mit zwei Trailing-Kommentaren stößt an. Bewusst — die Alternative wäre
  Raten.
- **Dogfood-Lücke:** d-check nutzt in den eigenen Doku-Tabellen keinen
  Trailing-Marker und konnte den Defekt daher an sich selbst nicht bemerken —
  dieselbe Klasse wie die Range-Notation (slice-073 §4). Offener Punkt: ob die
  Reader-Grammatik einen Konsumenten-nahen Fixture-Anker braucht.
- **Der Guard ist jetzt zweimal aufgeweicht worden** (Link-Suffix, Kommentar-Suffix),
  beide Male begründet und eng. Ein drittes Mal wäre ein Signal, dass die
  „lexikalischer Splitter"-Annahme selbst falsch ist und der Reader einen echten
  Markdown-Parser braucht — das wäre eine eigene, große Entscheidung.

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
durch die Akzeptanztests aus [slice-070](../done/slice-070-trace-tabellenquellen-nullmengen-guard.md) geschützt.

## 7. Closure-Notiz (nach `done/`)

_Ausstehend — wird bei Abschluss mit Commit-Hash, Review-Verdikt und Lerneintrag
gefüllt._
