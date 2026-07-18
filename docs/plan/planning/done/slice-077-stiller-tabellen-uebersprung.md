# Slice slice-077: Stiller Tabellen-Übersprung — eine irrelevante Tabelle verschluckt die nächste

**Status:** done — **welle-60, abgeschlossen 2026-07-18** (Review R1
ACCEPT-WITH-NITS, alle drei Nits eingearbeitet, Release v0.48.0). Closure-Notiz in
§7. welle-60 bleibt **pausiert** (slice-074/071 offen). Erfasst 2026-07-17 als
Messbefund mit bewusst offener Regel; die Regel ist gemessen entschieden
(relevant-Header-Grenze).

**Bezug:** betrifft
[`DC-FA-REQ-001.a`](../../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 5 und — über den geteilten Reader —
[`DC-FA-XREF-001.a`](../../../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency).
Berührt den Guard aus
[ADR-0037](../../adr/0037-trace-tabellenquellen-nullmengen-guard.md). Begründende
Entscheidung [ADR-0043](../../adr/0043-tabellengrenze-am-relevanten-header.md)
(Proposed) — die Regel wurde erst **gemessen belegt**, dann benannt (der Bedingung
aus [slice-074](../open/slice-074-kommentar-suffix-tabellenzeilen.md) §2 genügend).
**Defekt-Fix, kein CR** (das Lastenheft definiert keine Tabellengrenze).
**SemVer-Minor.**

**Autor:** pt9912. **Datum:** 2026-07-17.

---

## 1. Ziel

Eine **irrelevante** Tabelle (ihr Header bindet keine konfigurierte Rolle)
verschluckt die unmittelbar folgende **relevante** Tabelle samt Header,
Trennzeile und allen Datenzeilen — und weil die verschluckende Tabelle
irrelevant ist, meldet `extractTable` nichts. **Die Anforderungen der zweiten
Tabelle existieren nicht, und kein Gate sagt es.**

Gemessen gegen das **ausgelieferte** `ghcr.io/pt9912/d-check:v0.45.1` <!-- d-check:ignore (historischer Beleg gegen das damals ausgelieferte Image; darf nicht auf die aktuelle Version mitwandern) -->, Fixture
`fx-s`, **ohne jeden Marker**:

```markdown
| ID | Anforderung | Notiz |          ← relevante Tabelle 1, gedeckt
|---|---|---|
| F-1 | Alpha | ok |
                                       ← Leerzeile
| Werkzeug | Zweck | Stand |          ← IRRELEVANTE Tabelle, gleiche Breite
|---|---|---|
| foo | bar | baz |
| ID | Anforderung | Notiz |          ← relevante Tabelle 3: als Datenzeile gefressen
|---|---|---|
| F-2 | Beta | ok |
| F-3 | Gamma | ok |
```

```text
v0.45.1 --trace --require-complete   ⇒  1 Anforderung(en), 0 Waise(n).   EXIT=0
                                        (F-2 und F-3 sind echte Waisen)
```

Ein Lauf, der zwei echte Waisen verschweigt und **grün** endet. Vorbedingung ist
allein: gleiche Breite, keine Leerzeile. Kein Marker, keine Direktive, keine
Sonderform.

## 2. Entscheidungen / Regel

**Die Tabellengrenze liegt am *relevanten* Header** — vollständig in
[ADR-0043](../../adr/0043-tabellengrenze-am-relevanten-header.md). Kurz: eine Zeile,
die mit ihrer Folgezeile einen gültigen Header + Trennzeile bildet **und** deren
Header eine konfigurierte Rolle bindet (relevant), beendet die laufende Tabelle —
auch bei passender Zellenzahl. Ein rollenloser Header beendet **nicht**.

Warum genau **relevant** — das ist der Diskriminator, den die fünf slice-074-Fassungen
nicht hatten:

- **`fx-s` und `fx-t` sind syntaktisch identisch** (`[Textzeile][Trennzeile][Daten…]`);
  jede rein **strukturelle** Grenze zerreißt `fx-t` (eine `-`-Datenzeile würde zum
  Header). Das „relevant" trennt sie: in `fx-t` bindet die `-`-Zeile keine Rolle ⇒
  keine Trennung; in `fx-s` bindet der verschluckte Header die Rollen ⇒ Trennung.
- **GFM gibt uns recht** (goldmark v1.8.4 sieht `fx-s` wie wir — ohne Leerzeile keine
  neue Tabelle): die Frage ist **Policy**, nicht Grammatik.
- **Invariante:** *Ein relevanter Header beendet die laufende Tabelle* ⇒ jede
  relevante Tabelle wird erkannt, keine Anforderung verschwindet still in einer
  vorangehenden. Das adressiert den Kern aus
  [Review R3](../../../reviews/2026-07-17-slice-074-implementation-r3.md) (den
  Wiederaufsetz-Punkt des Header-Scans) **direkt**, statt Nachbarfälle zu verengen.

**Gemessen** (Prototyp gegen das ausgelieferte v0.47.0-Image, 2026-07-18) — die
**erste** Regel, die `fx-s` **und** `fx-t` **und** `fx-m` gleichzeitig besteht:

| Fixture | v0.47.0 | mit Regel |
|---|---|---|
| `fx-s` (der Defekt) | `total 1` (still) | `total 3` ✓ |
| `fx-t` (Gegenprobe, Attempt-5-Killer) | `total 2` | `total 2` ✓ (keine Falsch-Trennung) |
| `fx-m` (slice-074-Aufsatz) | Exit 2 | `total 3` ✓ (mit der Toleranz) |
| `fx-m3` (echte Verrutschung) | Exit 2 | Exit 2 ✓ (Guard scharf) |
| volle Suite (slice-070) | grün | grün ✓ |

**Bedingung an die Umsetzung** (die R3-F-2-Lehre): jede Grenze **per Mutation
gepinnt**, gemessen; das Relevanz-Prädikat an **beide** Konsumenten durchgereicht
(`format: table` via `bindTableColumns`, `cross-consistency` via `bindCrossColumns`).
Die früheren Kandidaten (`isNewTableHeader` unbedingt, Trennzeile-im-Body,
Zeilen-Relevanz) sind in [ADR-0043](../../adr/0043-tabellengrenze-am-relevanten-header.md)
§Verglichene Alternativen als widerlegt protokolliert.

## 3. Definition of Done

- [x] **Befund erfasst** samt Reproduktion gegen das ausgelieferte Image.
- [x] **Tragende Regel** benannt und begründet ⇒
  [ADR-0043](../../adr/0043-tabellengrenze-am-relevanten-header.md) (Proposed).
- [x] **Spezifikation:** [`DC-FA-REQ-001.a`](../../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
  Schritt 5 geschärft (Grenze am relevanten Header) + Historie.
- [x] **Regel-Kandidat gemessen** (Prototyp, 10 Fixtures + volle Suite grün gegen
  v0.47.0) — der Schritt, den die fünf slice-074-Fassungen übersprangen.
- [x] **Fixtures zuerst, rot:** `fx-s` (der Befund, rot→grün), `fx-t` (all-dashes-
  Gegenprobe, no-regress), `fx-adj` (breitensensitiv, echter Pin), `fxdup`
  (fail-closed), plus ein cross-Pin — auf **Konsumenten-Ebene** (`format: table`
  **und** `cross-consistency`).
- [x] **Implementierung:** relevant-Header-Grenze in `consumeTableRows`; das
  Relevanz-Prädikat an **beide** Konsumenten durchgereicht (fail-closed auf den
  Bind-Fehler, Review R-F-3).
- [x] **Mutations-Härte:** **gemessen** — Grenze aus ⇒ `fx-s`/`fx-adj`/`fxdup`/cross
  kippen, `fx-t` bleibt grün; fail-open ⇒ nur `fxdup` kippt (die R3-F-2-Lehre,
  diesmal empirisch nachgefahren, auch vom Reviewer).
- [x] **Realdatenbeleg** (grid-gym read-only, kein Regress) + **unabhängiger,
  kontext-getrennter Review vor** dem Release (R1 ACCEPT-WITH-NITS);
  `make gates`/`make ci` grün.
- [x] **Release** (SemVer-Minor, **v0.48.0**) + CHANGELOG mit Rot-werden-Ansage;
  Tag + GHCR (Digest `sha256:644b0c71…25d9`) + Digest-Backfill (Handbuch §4).
  Release-Run 29635549944 grün.

## 4. Risiken / offene Punkte

- **Der Defekt bleibt bis dahin ausgeliefert.** Er ist von den drei offenen
  Tabellen-Defekten der einzige, der heute **still Waisen verschweigt** —
  slice-074 ist laut (Exit 2), [slice-076](../done/slice-076-markdown-lexik-commonmark.md)
  ist Blindheit ohne Falschaussage. Dieser hier sagt „0 Waise(n)", während zwei
  existieren.
- **Praxis-Häufigkeit unbekannt.** In den 522 gemessenen Realdateien kam die
  Vorbedingung **nicht** vor (grid-gyms `spec/` ist sauber). Der Befund ist
  konstruiert-reproduzierbar, nicht real beobachtet — das senkt die Dringlichkeit,
  nicht die Gültigkeit.
- **Kopplung an slice-074:** dessen Toleranz **verbreitert** dieses Loch. Wer
  slice-074 wieder aufnimmt, muss diesen Slice mitdenken, sonst wird R3-F-1 zum
  vierten Mal reproduziert.
- **Die Verlockung ist eine schnelle Regel.** Sie ist der Grund, warum dieser
  Slice ohne ADR erfasst wird.

## 5. Trigger

Isolierungs-Messung während der Rücknahme von
[slice-074](../open/slice-074-kommentar-suffix-tabellenzeilen.md) (2026-07-17): auf die
Frage, ob R3-F-1 die Klasse **einführt** oder eine bestehende **verbreitert**,
wurde `fx-s` ohne Marker gegen v0.45.1 gefahren — und das ausgelieferte Image
verschwieg zwei echte Waisen. Der Befund war in drei unabhängigen Reviews und
fünf Implementierungsanläufen nicht aufgetaucht, weil alle die **Marker**-Frage
untersuchten.

## 6. Sub-Area-Modus-Begründung

GF (Repo-Default): Die Spezifikation führt, der Code folgt. Hier besonders: die
Regel gehört **zuerst** in
[`DC-FA-REQ-001.a`](../../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 5 und in eine ADR — der Code ist die letzte Station, nicht die erste.

## 7. Closure-Notiz (nach `done/`)

**Abschluss 2026-07-18, Release v0.48.0.** slice-077 ist die **Wurzel** der
slice-074-Klasse: erst die robuste Tabellengrenze macht die Direktiven-Toleranz
(slice-074) sicher und slice-071s Realdatenbeleg fahrbar.

**Commit-Kette:** Aktivierung `af25012`/`a05b2f9` → Doc-first `69fcddc` → feat
`3ea0090` → Review-Report `43f0005` → Review-Fixes `99bbad8` → Release-Prep
`0751ba3` → Closure-Move `30b75d4` → Closure-Body (dieser) → ADR-Accepted →
Post-Release (Tag/GHCR/Digest-Backfill).

**Geliefert:** die Tabellengrenze am **relevanten** Header (`consumeTableRows`,
Prädikat durchgereicht an `format: table` via `bindTableColumns` und
cross-consistency forward/backward via `bindCrossColumns`). Ein Header, der eine
Rolle bindet, beendet die laufende Tabelle — jede relevante Tabelle wird erkannt.

**Der Durchbruch (gemessen):** das „relevant" bricht die syntaktische
`fx-s`≡`fx-t`-Mehrdeutigkeit — die **erste** Regel, die `fx-s` **und** `fx-t`
gleichzeitig besteht. Fünf rein strukturelle Fassungen scheiterten genau hier. Der
Weg war **messen zuerst** (Prototyp gegen das ausgelieferte v0.47.0-Image, 10
Fixtures + volle Suite) **vor** dem ersten Produktivcode — der Schritt, den die fünf
Fassungen übersprangen.

**Review R1 ACCEPT-WITH-NITS**
([Report](../../../reviews/2026-07-18-slice-077-implementation-r1.md)), 0
HIGH/MEDIUM. Der kontext-getrennte Reviewer belegte den Kern-Fix empirisch gegen die
Images (`fx-s` Elter `total 1` → HEAD `total 3`) und fuhr die Mutations-Härte selbst
nach. Drei LOW, alle eingearbeitet:

- **R-F-1:** `fx-adj` war grenzenblind (behauptete einen Pin, der nicht hielt — die
  R3-F-2-Klasse in genau diesem Slice). Neu mit **unterschiedlicher Breite** ⇒
  echter, gemessener Pin.
- **R-F-2:** die Vorwärts-Cross-Durchreichung war ungetestet. Die Cross-Fixture trägt
  nun in **beiden** Sichten eine irrelevante Vor-Tabelle ⇒ beide Cross-Pfade gepinnt.
- **R-F-3:** ein mehrdeutig-relevanter Header (doppelte Rollen-Spalte) wurde hinter
  einer irrelevanten Tabelle **still** verschluckt (fail-open). Das Prädikat feuert
  jetzt **fail-closed auf den Bind-Fehler** ⇒ Exit 2 statt still (`fxdup`-Pin).

**Reusable Lehren:**

- **Bei einer mehrfach gescheiterten Klasse zuerst messen, dann proponieren.** Fünf
  Fassungen schlossen je aus einem geprüften Zweig auf die Klasse; der gemessene
  Prototyp über 10 Fixtures fand die erste Regel, die alle besteht. „Gemessen, nicht
  behauptet" gilt auch für die **Regel-Wahl**, nicht nur die Mutations-Härte.
- **Der Diskriminator lag nicht in der Struktur, sondern in der Semantik.** `fx-s`
  und `fx-t` sind syntaktisch identisch; erst die **Rollen-Bindung** (config-getragen)
  trennt sie. Wer eine strukturelle Mehrdeutigkeit rein strukturell lösen will,
  bekommt die nächste Fassung.
- **Ein Mutations-Pin, der grenzenblind ist, ist keiner** (R-F-1) — der Test muss bei
  genau dieser Breite/Config kippen, sonst pinnt er eine grenzenunabhängige
  Eigenschaft. Empirisch gegenfahren, nicht dem Kommentar glauben.
- **Fail-open an einer Bind-Fehler-Kante ist ein stiller Pfad** (R-F-3): ein Prädikat,
  das einen Fehler als „nicht relevant" wertet, verschluckt genau die Tabellen, die
  laut sein sollten. Fail-closed auf den Fehler.
