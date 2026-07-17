# Slice slice-076: Markdown-Lexik an CommonMark/GFM angleichen (Trennzeile + Fence-Infozeile)

**Status:** next (Backlog, auf Aufnahme in eine Welle wartend).

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
  [slice-074](../open/slice-074-kommentar-suffix-tabellenzeilen.md) §2: fünf Fassungen,
  die eine Regel knapp neben dem belegten Problem platzierten.
- **Kein Parser** — gemessen, nicht behauptet ([ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) Entscheidung 4).

## 3. Definition of Done

- [x] **Spezifikation:** beide Regeln geschärft + Historie.
- [x] **ADR + Index:** [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md),
  Proposed, im Index.
- [ ] **Tests aus den Realdaten (zuerst, rot):** die `carveouts.md`-Form
  (`| ------- | ------------- | -- |`) und die
  `2026-06-19-slice-030-…md:179`-Form (```` ```yaml-Fence (`datei.md`) — … ````).
  **Konsumenten-Ebene**, nicht nur am Lexer: die Fence-Form über ein
  **anderes** Modul als `trace` (der kaputte Link dahinter muss gemeldet
  werden) — der Defekt ist modulübergreifend, ein Tabellentest allein belegt ihn
  nicht.
- [ ] **Implementierung:** `tableDelimiterCell`
  (`internal/hexagon/core/app/trace_table.go`) und der Fence-Automat. **Achtung,
  zwei Automaten:** `proseLines` (`internal/hexagon/core/rules/markdown.go`,
  geteilt von allen Modulen) und `markdownTableLines`
  (`internal/hexagon/core/app/trace_table.go`). Die Infozeilen-Regel gehört in
  **beide**, sonst sieht `trace` das Dokument anders als `links`.
- [ ] **Mutations-Härte:** Trennzeilen-Lockerung zurückgedreht ⇒ mindestens ein
  Test kippt; Infozeilen-Regel entfernt ⇒ mindestens ein Test kippt. **Gemessen,
  nicht zugesagt** — die Suite ist heute gegen **beide** blind (die Lockerung
  ließ `make test` grün). Das ist die R3-F-2-Lehre.
- [ ] **Differential-Gegenprobe:** der goldmark-Spike über die 522 Realdateien
  fällt von 8 auf 2 Abweichungen; die zwei Reste sind Policy
  ([slice-077](../open/slice-077-stiller-tabellen-uebersprung.md)), nicht Grammatik.
- [ ] **Nutzerdoku:** Handbuch — die Trennzeilen- und Fence-Regel stehen dort als
  Nutzer-Vertrag; CHANGELOG mit dem Release-Prep. **Die Minor-Ansage gehört in
  die Release-Notiz:** ein grüner Konsumentenlauf kann danach rot werden.
- [ ] **Release:** v0.46.0 (Minor, **kein** Patch) + Tag + GHCR + Digest-Backfill.
- [ ] **Qualität:** unabhängiger, kontext-getrennter Review **vor** dem Release;
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
während der Rücknahme von [slice-074](../open/slice-074-kommentar-suffix-tabellenzeilen.md).
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

_Ausstehend — wird bei Abschluss mit Commit-Hash, Review-Verdikt und Lerneintrag
gefüllt._
