# Slice slice-073: Link-transparente Range-Fortsetzung (ausgelieferter Coverage-Defekt)

**Status:** open (Backlog; noch keiner Welle zugeordnet).

**Welle:** keine; Kandidat für welle-60 — der Fix sitzt im **geteilten** Parser und
ist damit Voraussetzung des Realdatenbelegs von
[`slice-071`](../in-progress/slice-071-trace-cross-consistency-gate.md), nicht eine Parallelbaustelle.

**Bezug:** **Defekt-Fix**, **kein Change Request**: das Lastenheft verspricht die
Range-Expansion unqualifiziert
([`DC-FA-COV-001`](../../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
„`<FAM>-AAA..BBB` deckt alle …"); die Verengung „**unmittelbar**" stand allein in
der Spezifikation und ist dort geschärft
([`DC-FA-COV-001.a`](../../../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)
Schritt 3, wirkt über den geteilten Parser zugleich auf
[`DC-FA-XREF-001.a`](../../../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)).
Begründende Entscheidung
[ADR-0039](../../adr/0039-link-transparente-range-fortsetzung.md) (Proposed).
Mit-betroffen:
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
(die Linkpflicht ist die Ursache der Kollision). **SemVer-Patch** — die Zusage wird
hergestellt, nicht geändert.

**Autor:** pt9912. **Datum:** 2026-07-17.

---

## 1. Ziel

Eine Range-/Enum-Fortsetzung hinter einer **verlinkten** Kennung expandiert nicht:
`` [`GG-UI-001`](…)..009 `` liefert nur `GG-UI-001`. Ursache ist die
Spec-Verengung „die Fortsetzung muss der Fundstelle **unmittelbar** folgen" — hinter
einem Link folgt sie eben nicht unmittelbar.

Das kollidiert **strukturell** mit d-checks eigenem Modul `ids`: wo Linkpflicht
gilt, trägt jede Kennung ein Link-Suffix. Wer Linkpflicht **und** Range-Notation
nutzt, verliert die Expansion still. Der Slice macht den Parser link-transparent
(genau ein Suffix) und stellt die Lastenheft-Zusage her.

## 2. Entscheidungen / Regel

- **Genau ein** Link-Suffix `](…)` darf die Fortsetzung unterbrechen; dahinter gilt
  wieder „unmittelbar". Kein weiteres Peeling (Whitespace, Emphasis, zweites
  Suffix, Text dazwischen) — jede weitere Toleranz **rät** die Autor-Absicht
  ([ADR-0039](../../adr/0039-link-transparente-range-fortsetzung.md)).
- **Ein Fix, zwei Konsumenten:** der Parser ist geteilt, der Fix wirkt auf
  `trace.coverage` **und** `trace.cross-consistency`. Derselbe Defekt saß in beiden
  — das ist der Preis der geteilten Mechanik und zugleich ihr Nutzen.
- **Kein Lastenheft-Bump:** die Zusage steht bereits unqualifiziert im Lastenheft;
  geschärft wird Rang 2. Der Slice ist ein **Patch**, kein Minor.

## 3. Definition of Done

- [x] **Spezifikation:** [`DC-FA-COV-001.a`](../../../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)
  Schritt 3 um die Link-Transparenz geschärft, inkl. der **negativen** Abgrenzung
  (was bewusst nicht übersprungen wird) + Historie.
- [x] **ADR + Index:** [ADR-0039](../../adr/0039-link-transparente-range-fortsetzung.md)
  (Regel, Alternativen, Verhaltensänderung für Bestandskonsumenten), Status
  Proposed, im Index.
- [ ] **Implementierung:** der geteilte Range-Parser überspringt höchstens **ein**
  Link-Suffix; unverlinkte Ranges, Enum-Notation und die Fail-closed-Fälle
  (`AAA>BBB`, Breiten-Mismatch) bleiben unverändert.
- [ ] **Tests (positiv):** identische Quelle, einmal `GG-UI-001..003`, einmal
  `` [`GG-UI-001`](…)..003 `` ⇒ **gleiches** Ergebnis — je einmal für
  `trace.coverage` (die ausgelieferte Regression: 2 Waisen ⇒ 0) und für
  `trace.cross-consistency` (2 Differenzen ⇒ 0). Enum-Form `` [`ID`](…)/004/005 ``
  ebenso.
- [ ] **Tests (negativ, gegen das Raten):** zwei Link-Suffixe hintereinander, ein
  Zeichen zwischen `)` und `..`, Whitespace davor ⇒ **keine** Expansion. Ohne diese
  Tests wäre „genau eins" eine Behauptung.
- [ ] **Mutations-Härte:** die Suffix-Überspringung einzeln herausmutiert kippt
  genau einen Test; die „genau eins"-Grenze ebenso.
- [ ] **Nutzerdoku:** Handbuch (§5 `trace.coverage`/`cross-consistency`:
  Range-Notation unter Linkpflicht) + CHANGELOG (als **Fixed**, mit dem Hinweis,
  dass ausgelieferte Läufe grüner werden können).
- [ ] **Release:** v0.44.1, Release-Prep + Tag + GHCR + Digest-Backfill.
- [ ] **Qualität:** unabhängiger, kontext-getrennter Review; `make gates`/`make ci`
  grün.

## 4. Risiken / offene Punkte

- **Verhaltensänderung für Bestandskonsumenten:** wer heute verlinkte Ranges in
  `trace.coverage` führt, sieht **weniger** Waisen — ein fälschlich roter Lauf wird
  grün. Kein Konsument verliert Deckung, kein Befund entsteht neu; trotzdem ändert
  sich Ausgabe, die vorher „stabil falsch" war. Gehört sichtbar in den CHANGELOG.
- **„Genau eins" ist eine gesetzte Zahl.** Sie ist durch Realdaten gedeckt (40/40
  Fortsetzungen des Konsumenten haben exakt diese Form), aber ein Repo mit anderer
  Schreibweise stößt erneut an. Bewusst: die Alternative wäre Raten.
- **Der Parser trägt jetzt Markdown-Wissen**, war vorher rein lexikalisch. Ein
  zweiter Ort, an dem Link-Syntax verstanden wird (neben dem Reader) — klein, aber
  eine Kopplung, die bei Markdown-Änderungen mitwandern muss.
- **Der Defekt war drei Releases lang unentdeckt** (v0.41.0 → v0.44.0), obwohl
  `trace.coverage` einen Realdatenbeleg hatte. Er fiel erst auf, als der
  Kreuzverweis-Abgleich **zwei** Sichten verglich und die Asymmetrie sichtbar
  machte. Offener Punkt: ob die Range-Fähigkeit einen Dogfood-Anker im eigenen Repo
  braucht — d-check nutzt selbst keine Range-Notation und kann den Defekt daher
  nicht an sich selbst bemerken.

## 5. Trigger

Realdaten-Lauf des Konsumenten grid-gym gegen v0.44.0 (2026-07-17): 218
Differenzen, davon zwei belegte Defekte. Defekt 2 — `` [`GG-UI-001`](…)..009 ``
expandiert nicht — wurde isoliert reproduziert (verlinkte Range ⇒ 2 Differenzen,
unverlinkt ⇒ 0). Bei der Nachprüfung zeigte sich, was der Konsument **nicht** sehen
konnte: derselbe geteilte Parser bedient `trace.coverage`, dort ist der Defekt seit
**v0.41.0 ausgeliefert** (identische Quelle, nur verlinkt: 0 Waisen ⇒ 2 falsche
Waisen). Der Konsument zählte sein Korpus aus statt zu schätzen: 40 von 40
Fortsetzungen haben die Form `` [`ID`](url) `` + `..NNN`/`/NNN`.

## 6. Sub-Area-Modus-Begründung

GF (Repo-Default): Die Spezifikation führt, der Code folgt. Der Range-Parser ist
bestehender, spezifizierter Code — die Schärfung ist ein Vertrags-, kein
Rückbau-Zug; die Kompatibilität der unverlinkten Formen ist durch die vorhandenen
Akzeptanztests (slice-067) geschützt.

## 7. Closure-Notiz (nach `done/`)

_Ausstehend — wird bei Abschluss mit Commit-Hash, Review-Verdikt und Lerneintrag
gefüllt._
