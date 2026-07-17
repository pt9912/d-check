# Slice slice-073: Link-transparente Range-Fortsetzung (ausgelieferter Coverage-Defekt)

**Status:** **done** (2026-07-17, welle-60-trace-cross-consistency). Ausgeliefert
als v0.45.1; R2 ACCEPT-WITH-NITS, R2-F-1 geschlossen.

**Welle:** welle-60 — der Fix sitzt im **geteilten** Parser und ist damit
Voraussetzung des Realdatenbelegs von
[`slice-071`](../open/slice-071-trace-cross-consistency-gate.md), nicht eine
Parallelbaustelle.

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
- [x] **Implementierung:** der geteilte Range-Parser überspringt höchstens **ein**
  Link-Suffix; unverlinkte Ranges, Enum-Notation und die Fail-closed-Fälle
  (`AAA>BBB`, Breiten-Mismatch) bleiben unverändert.
- [x] **Tests (positiv):** identische Quelle, einmal `GG-UI-001..003`, einmal
  `` [`GG-UI-001`](…)..003 `` ⇒ **gleiches** Ergebnis für `trace.coverage` (die
  ausgelieferte Regression: 2 Waisen ⇒ 0). Enum-Form `` [`ID`](…)/004/005 `` ebenso.
  Der zugesagte `trace.cross-consistency`-Test (verlinkte Range durch den
  Kreuzverweis-Abgleich) **fehlte im v0.45.1-Stand** und wurde erst nach dem
  R2-Review nachgezogen (`b8c503a`, R2-F-1) — inkl. des scharfen Klammer-URL-Falls,
  der bei der Fix-Mutation kippt (die Coverage-Tests allein deckten die Achse
  nicht).
- [x] **Tests (negativ, gegen das Raten):** zwei Link-Suffixe hintereinander, ein
  Zeichen zwischen `)` und `..`, Whitespace davor ⇒ **keine** Expansion. Ohne diese
  Tests wäre „genau eins" eine Behauptung.
- [x] **Mutations-Härte:** verifiziert — die Suffix-Überspringung entfernt kippt
  `range hinter Link mit Code-Span`, die „genau eins"-Grenze aufgehoben kippt
  `zwei Link-Suffixe`.
- [x] **Nutzerdoku:** Handbuch §5 (Range-Notation unter Linkpflicht, mit dem
  Klammer-URL-Beispiel) + Handbuch-Historie 1.33; CHANGELOG als **Fixed** (0.44.1
  Erst-Fix; 0.45.1 klammer-balancierte Nachbesserung, „Betroffen: v0.44.1 und
  v0.45.0"). Die durch R1-F-3 widerlegten Zusagen des 0.44.1-Blocks sind dort
  sichtbar als „war falsch"/„traf nicht zu" korrigiert.
- [x] **Release:** v0.44.1 (Erst-Fix), **überholt durch v0.45.1** (R1-F-1:
  Ziel klammer-balanciert) — der finale ausgelieferte Stand ist v0.45.1. Beide
  getaggt, in `version.md` registriert, auf GHCR (v0.45.1-RepoDigest im R2-Report
  belegt). Der Test-Nachzug `b8c503a` härtet nur die Suite und ändert kein
  ausgeliefertes Verhalten ⇒ kein weiterer Bump.
- [x] **Qualität:** R1 (REJECT, F-1…F-5) → **R2** (ACCEPT-WITH-NITS, alle fünf
  R1-Befunde durch eigene Messung geschlossen; einziger Rest R2-F-1 mit `b8c503a`
  behoben), beide kontext-getrennt; `make gates`/`make ci` grün.

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

## 7. Closure-Notiz

**Abgeschlossen 2026-07-17**, welle-60. Ausgelieferter Stand **v0.45.1**.

**Commit-Kette:** `2954e4d` (feat: Range-Fortsetzung link-transparent, ein
Suffix) · `6925987` (fix R1-F-1: Ziel klammer-balanciert über `LinkSuffixEnd`
statt Regex `[^)]*`; zugleich v0.45.1-Tag-Commit) · `b8c503a` (test R2-F-1:
cross-consistency-Achse verriegelt). [ADR-0039](../../adr/0039-link-transparente-range-fortsetzung.md) mit der Closure auf `Accepted`
(Folge-Commit) — die Entscheidung ist umgesetzt und ausgeliefert.

**Review-Verlauf:** R1 **REJECT** (1 HIGH, 3 MEDIUM, 1 LOW), nachgeholt nach der
Auslieferung von v0.44.1 — der HIGH war eine **neue** Falsch-Deckung: die naive
„bis zur ersten `)`"-Abgrenzung riss bei Klammern im Linkziel den URL-Rest in den
Range-Parser und versteckte Waisen (stiller-Grün-Pfad). Behoben in v0.45.1. R2
**ACCEPT-WITH-NITS** (kontext-getrennter Subagent, ohne Zugriff auf die
Session-Analyse): alle fünf R1-Befunde durch **eigene Messung** gegen das
ausgelieferte Image geschlossen, per Mutation belegt; einziger Rest R2-F-1 (LOW,
die Cross-Seite test-blind) mit `b8c503a` geschlossen.

**Lerneintrag (reusable):**
- **Zwei Konsumenten, zwei Vorläufe — ein geteilter Fix verriegelt nicht
  automatisch beide.** Der Fix wirkt über `skipLinkSuffix` auf `trace.coverage`
  (Prosa-Text) **und** `trace.cross-consistency` (Tabellenzellen via
  `rangeAwareIDs`). Die Coverage-Tests deckten die Cross-Achse **nicht** — R2-F-1.
  Ein „ein Fix, zwei Konsumenten" braucht **je einen** Test pro Konsument, sonst
  ist die zweite Seite gegen einen seitenspezifischen Umbau blind.
- **Ein grüner Test beweist nichts über die Sensor-Härte.** Mein erster
  Cross-Test nutzte ein Linkziel ohne innere Klammer und kippte bei der
  Fix-Mutation **nicht** — er pinnte den gutartigen Zweig (die R1-F-4-Falle). Erst
  der Klammer-URL-Fall, per Mutation gemessen **vor** dem Commit, verriegelt die
  Achse wirklich. Sensor-Härte gehört gemessen, nicht behauptet.
- **Ein Defekt kann drei Releases still überleben** (v0.41.0 → v0.44.0), obwohl
  `trace.coverage` einen Realdatenbeleg hatte. Er fiel erst auf, als der
  Kreuzverweis-Abgleich **zwei** Sichten verglich und die Asymmetrie sichtbar
  machte — die Dogfood-Lücke (§4) besteht fort: d-check nutzt selbst keine
  Range-Notation.
- **Ein nachgeholter Review nach dem Release ist kein Widerspruch, sondern die
  Rettung:** R1 fand nach v0.44.1 einen HIGH, der zu v0.45.1 führte. Die
  slice-073-Lehre für den Flow: der unabhängige Review gehört **vor** den Tag —
  hier kam er zu spät und kostete einen Patch-Release.
