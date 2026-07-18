# Slice slice-076: Markdown-Lexik an CommonMark/GFM angleichen (Trennzeile + Fence-Infozeile)

**Status:** done — **welle-60, abgeschlossen 2026-07-18** (Review ACCEPT-WITH-NITS,
Release v0.47.0). Closure-Notiz in §7. welle-60 bleibt **pausiert, nicht
abgeschlossen** (slice-071 weiter blockiert). Aus `next/` in Arbeit genommen,
nachdem die slice-075-Closure den WIP-Slot freigab (Modul 5, WIP-Limit = 1).
**Vorgeschichte:** am 2026-07-17 war der Slice bei **belegtem** WIP-Limit
versehentlich nach `in-progress/` gezogen und wieder nach `next/` zurückgeführt
worden; die Implementierung hatte nie begonnen. Doc-first
([ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md), Spezifikation) lag
fertig vor; offen war nur die Implementierung.

**Welle:** welle-60-trace-cross-consistency. Blockiert niemanden und wird von
niemandem blockiert — der Defekt trifft denselben geteilten Reader und dieselbe
Vorverarbeitung, ist aber unabhängig; gefunden bei der Rücknahme von slice-074,
nicht von ihr verursacht.

**Bezug:** **Defekt-Fix**, **kein Change Request**: das Lastenheft sagt weder, was
eine Trennzeile ist, noch was einen Fence öffnet — beides ist
Spezifikations-Sache (Rang 2). Geschärft in
[`DC-FA-REQ-001.a`](../../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 3 (Trennzelle) und
[`DC-FA-LINK-001.a`](../../../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
Schritt 1 (Fence-Infozeile). Begründende Entscheidung
[ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) (Proposed).
**SemVer-Minor** — d-check findet danach **mehr**.

**Autor:** pt9912. **Datum:** 2026-07-17.

---

## 1. Ziel

d-checks handgeschriebene Markdown-Lexik weicht an zwei Stellen von
CommonMark/GFM ab — und zwar **still**: eine Regel, die eine Struktur nicht
erkennt, meldet nichts. Sie meldet **weniger**. Beide Abweichungen sind
ausgeliefert.

**A · Trennzelle.** `^:?-{3,}:?$` verlangt drei Bindestriche; GFM verlangt
**einen**. Jede reale Tabelle mit `| -- |` oder `| - |` ist für d-check keine
Tabelle — ihre Anforderungen, Links und IDs existieren nicht.

**B · Fence-Infozeile.** Eine `` ``` ``-Zeile, deren Rest einen Backtick enthält,
ist in CommonMark **kein** Fence-Öffner. d-check hält sie für einen und markiert
**den ganzen Rest der Datei als Nicht-Prosa** — für **alle** Module
([`DC-FA-LINK-001.a`](../../../../spec/spezifikation.md#dc-fa-link-001a--markdown-vorverarbeitung-und-link-extraktion)
Schritt 1: „Zeilen im Fence-Zustand werden von allen Modulen ignoriert").

## 2. Entscheidungen / Regel

Vollständig in [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md).
Kurz:

- **Trennzelle `^:?-+:?$`** — rein erweiternd, jede heute erkannte Trennzeile
  bleibt erkannt.
- **Backtick-Fence-Öffner mit Backtick in der Infozeile ist Fließtext.** Für
  `~~~` gilt die Regel nicht (CommonMark-Asymmetrie, keine d-check-Eigenheit).
- **Die Grenze ist das Gemessene.** Keine weiteren CommonMark-Angleichungen auf
  Verdacht — das ist die Lehre aus
  [slice-074](../done/slice-074-kommentar-suffix-tabellenzeilen.md) §2: fünf Fassungen,
  die eine Regel knapp neben dem belegten Problem platzierten.
- **Kein Parser** — gemessen, nicht behauptet ([ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) Entscheidung 4).

## 3. Definition of Done

- [x] **Spezifikation:** beide Regeln geschärft + Historie.
- [x] **ADR + Index:** [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md),
  Proposed, im Index.
- [x] **Tests aus den Realdaten (zuerst, rot):** die `carveouts.md`-Form
  (`| ------- | ------------- | -- |`) und die
  `2026-06-19-slice-030-…md:179`-Form (```` ```yaml-Fence (`datei.md`) — … ````).
  **Konsumenten-Ebene**, nicht nur am Lexer: die Fence-Form über ein
  **anderes** Modul als `trace` (der kaputte Link dahinter muss gemeldet
  werden) — der Defekt ist modulübergreifend, ein Tabellentest allein belegt ihn
  nicht.
- [x] **Implementierung:** `tableDelimiterCell`
  (`internal/hexagon/core/app/trace_table.go`) und der Fence-Automat. **Achtung,
  zwei Automaten:** `proseLines` (`internal/hexagon/core/rules/markdown.go`,
  geteilt von allen Modulen) und `markdownTableLines`
  (`internal/hexagon/core/app/trace_table.go`). Die Infozeilen-Regel gehört in
  **beide**, sonst sieht `trace` das Dokument anders als `links`.
- [x] **Mutations-Härte:** Trennzeilen-Lockerung zurückgedreht ⇒ mindestens ein
  Test kippt; Infozeilen-Regel entfernt ⇒ mindestens ein Test kippt. **Gemessen,
  nicht zugesagt** — die Suite ist heute gegen **beide** blind (die Lockerung
  ließ `make test` grün). Das ist die R3-F-2-Lehre.
- [ ] **Differential-Gegenprobe — nicht re-run (dokumentierte Ausnahme):** der
  goldmark-Spike war eine Wegwerf-Messung über einen **Out-of-Repo-Korpus** (522
  Dateien = d-check-Doku + grid-gym-Kopie); goldmark ist bewusst **keine**
  Repo-Abhängigkeit ([ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md)
  Entscheidung 4), also ist das exakte „8→2" nicht in-repo reproduzierbar. Die
  Vollständigkeit der zwei Regeln ist stattdessen in-repo belegt:
  Konsumenten-Tests auf den realen Formen, isolierte Mutations-Härte, der
  unabhängige Review (rein erweiternd, CommonMark-treu) und der saubere
  Dogfood-Lauf über 244 eigene Dateien. Bewusst offen ausgewiesen, nicht still
  übersprungen — siehe §7.
- [x] **Nutzerdoku:** Handbuch — die Trennzeilen- und Fence-Regel stehen dort als
  Nutzer-Vertrag; CHANGELOG mit dem Release-Prep. **Die Minor-Ansage gehört in
  die Release-Notiz:** ein grüner Konsumentenlauf kann danach rot werden.
- [x] **Release:** **v0.47.0** (Minor, **kein** Patch — v0.46.0 war bei
  Slice-Erstellung noch frei, ist aber inzwischen von slice-075 vergeben) + Tag +
  GHCR (Digest `sha256:ad42432d…423eed`) + Digest-Backfill (Handbuch §4).
  Release-Run 29632470921 grün.
- [x] **Qualität:** unabhängiger, kontext-getrennter Review **vor** dem Release
  (ACCEPT-WITH-NITS, alle Nits eingearbeitet bzw. begründet ausgewiesen);
  `make gates`/`make ci` grün.

## 4. Risiken / offene Punkte

- **Der Lauf wird lauter, nicht stiller** — das ist der Zweck, aber es ist eine
  Verhaltensänderung für Bestandskonsumenten: bisher unsichtbare Tabellen
  liefern Anforderungen (⇒ neue Waisen), bisher unsichtbare Prosa liefert Links
  und IDs (⇒ neue Befunde). Deshalb Minor, und deshalb muss die Release-Notiz es
  sagen.
- **Zwei Fence-Automaten, zwei Verhalten.** `proseLines` ist ein **naiver
  Toggle** (jede Fence-artige Zeile kippt den Zustand, ohne Zeichen- oder
  Längenabgleich); `markdownTableLines` prüft beides. Die Spec beschreibt den
  naiven ⇒ **der Tabellen-Reader weicht heute schon von seiner eigenen Spec ab.**
  Dieser Slice fasst das **nicht** an: unbelegt, kein Realfall in den 522
  Dateien. Aber es ist eine Divergenz, die jemand einmal auflösen muss.
- **Die Lexik bleibt handgeschrieben.** Der Spike fand zwei Familien, weil zum
  ersten Mal gegen einen echten Parser gemessen wurde. Ein **Differential-Sensor**
  gegen goldmark (nur im Test, ohne Produktiv-Abhängigkeit) wäre der Sensor, der
  die dritte Familie findet, bevor ein Konsument sie findet. Offener Kandidat,
  nicht Teil dieses Slices.
- **Dogfood-Lücke, wieder:** die Fence-Fundstelle liegt in d-checks **eigener**
  Doku und blendet diese Datei seit Monaten — bemerkt hat es kein Gate, sondern
  erst ein Parser-Differential. Dieselbe Klasse wie slice-073 §4 und slice-074 §4.

## 5. Trigger

Nutzer-Frage „Brauchen wir einen besseren Markdown-Parser?" (2026-07-17), gestellt
während der Rücknahme von [slice-074](../done/slice-074-kommentar-suffix-tabellenzeilen.md).
Der daraufhin gefahrene Differential-Spike (goldmark v1.8.4 gegen 522 reale
Dateien, 490 Tabellen) lieferte **8 Abweichungen, alle in dieselbe Richtung**:
d-check sieht Tabellen nicht, die jeder Renderer zeigt. Die Antwort auf die Frage
ist zweigeteilt und gemessen — für die Policy nein, für die Grammatik ja; dieser
Slice nimmt die Grammatik-Hälfte **ohne** den Parser.

## 6. Sub-Area-Modus-Begründung

GF (Repo-Default): Die Spezifikation führt, der Code folgt. Beide Regeln sind
zuerst in der Spezifikation geschärft; der Code zieht nach. Die Kompatibilität
der heute erkannten Formen ist durch die bestehenden Akzeptanztests geschützt —
die Trennzeilen-Änderung ist rein erweiternd.

## 7. Closure-Notiz (nach `done/`)

**Abschluss 2026-07-18, Release v0.47.0** (nicht v0.46.0 — die Nummer war bei
Slice-Erstellung frei, wurde aber inzwischen von slice-075 belegt; in DoD, ADR und
Doku korrigiert).

**Commit-Kette:** Aktivierungs-Move `b4f6e26` → Slice-Body `98169c0` → feat
`2dbf4e1` → Review-Report `a9b4f07` → Review-Fixes `6a084db` → Release-Prep
`ab22a8c` → Closure-Move `f8805ab` → Closure-Body (dieser) → ADR-Accepted →
Post-Release (Tag/GHCR/Digest-Backfill).

**Geliefert:** zwei gemessene, still ausgelieferte Blindstellen zu — (A) die
Tabellen-Trennzelle folgt GFM (`^:?-+:?$` statt `^:?-{3,}:?$`); (B) eine
Backtick-Fence-Zeile mit Backtick in der Infozeile ist Fließtext, kein Öffner.
Regel B lebt nach dem Review an **einer** Stelle: das exportierte Prädikat
`rules.FenceToggle`, das alle drei Fence-Automaten (`proseLines`,
`diagramFenceLines`, `markdownTableLines`) als Öffner-Test rufen. Tests
rot-zuerst auf Konsumenten-Ebene (`links` **und** `trace`), Mutations-Härte pro
Regel isoliert gemessen.

**Review — ACCEPT-WITH-NITS**
([Report](../../../reviews/2026-07-18-slice-076-markdown-lexik-commonmark.md)),
0 HIGH/MEDIUM; der kontext-getrennte Reviewer fuhr die Mutations-Härte selbst
empirisch nach. Drei Nits, alle vor der Closure entschieden:

- **R-F-1 (LOW, eingearbeitet):** die Infozeilen-Regel stand doppelt (`proseLines`
  + Inline-Kopie in `markdownTableLines`). Aufgelöst zu **einem** exportierten
  `rules.FenceToggle`; die trace≠links-Divergenz, die dieser Slice schließt, kann
  nicht mehr durch einen einseitigen Edit still zurückkehren. Der `app`→`rules`-
  Import ist [ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md)-konform.
- **R-F-2 (INFO, eingearbeitet):** `FenceToggle` wirkt auch auf
  `diagramFenceLines`. Bewusst so — dieselbe Regel auf denselben naiven Toggle,
  keine *weitere* CommonMark-Regel auf Verdacht
  ([ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) Entscheidung 3
  gewahrt).
  Mutations-Pin ergänzt; die defined-in-Quelle liegt bewusst in einer **separaten**
  Datei, sonst maskiert der `proseLines`-Definitions-Scan die Mutation.
- **R-F-3 (LOW, begründet ausgewiesen):** der feat-Commit form-fixt die dattierte
  Review-Doku `2026-06-19-…:179` (die literale Backtick-Fence-Schreibweise wird zu
  einem `yaml`-Inline-Span). **Erzwungen, nicht optional:** der Fix un-blindet die
  Zeile, worauf der `spans`-Gate (kein Opt-out) den vorbestehenden, bis dahin
  verdeckten `span-unclosed` meldet — und er kann **nicht** in ein eigenes Commit,
  weil er ohne die Code-Änderung gar nicht sichtbar (also nicht behebbar) ist.
  Aussage unverändert; der in ADR/Spec zitierte Beleg bleibt **im ADR** wörtlich.

**Differential-Gegenprobe nicht re-run** (offen ausgewiesen, DoD §3): der
goldmark-„8→2"-Spike lief über einen Wegwerf-Out-of-Repo-Korpus; goldmark ist
bewusst keine Abhängigkeit. Vollständigkeit stattdessen in-repo belegt.

**Dogfood-Lücke geschlossen (slice §4):** die Fence-Fundstelle lag in d-checks
eigener Doku und blendete die Datei seit Monaten — kein Gate hatte es bemerkt, nur
ein Parser-Differential. Nach dem Fix sieht der eigene doc-check die Datei; die
eine erwartete Folge (`span-unclosed`) ist behoben.

**Reusable Lehren:**

- **Un-blindende Fixes ziehen Dogfood-Nachzug nach sich:** wer eine Lese-Regel
  erweitert, macht bisher verdeckte Doku-Zeilen sichtbar, und der eigene Gate
  meldet dort vorbestehende Defekte. Solche Form-Fixes gehören **in denselben
  Commit** wie die Code-Änderung (sonst ist der feat-Commit gate-rot), nicht in
  ein separates Hygiene-Commit — anders als die generische Prozess-Erwartung.
- **Ein Mutations-Pin muss den echten Maskierungs-Pfad umgehen:** der erste
  diagrams-Pin schlug fehl, weil ein **zweiter**, nicht mutierter Automat
  (`proseLines` im Definitions-Scan) die Mutation maskierte. Erst eine Fixture,
  die diesen Pfad ausschließt, fing die Rückdrehung. „Gemessen, nicht zugesagt"
  heißt auch: prüfen, dass der Pin wirklich kippt — sonst gibt er Falsch-Vertrauen.
- **Doppelte Regel-Kopien in getrennten Paketen driften still:** ein geteiltes,
  exportiertes Prädikat ist die einzige strukturelle Garantie gegen die Wiederkehr
  genau der Divergenz, die man gerade schließt (R-F-1).
