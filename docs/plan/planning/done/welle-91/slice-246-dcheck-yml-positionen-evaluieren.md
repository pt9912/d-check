# slice-246: `.d-check.yml`-Positionen aus ai-harness-init evaluieren

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle as State Machine.

**Welle:** welle-91.

**Bezug:** [`welle-91`](../../welle-91-adoption-ai-harness-init.md),
[`MR-073`](../../../../../harness/conventions.md#mr-073),
[`MR-034`](../../../../../harness/conventions.md#mr-034--die-referenzmatrix-bewacht-auch-die-kante-adr--welle)
(Verwandt: die welle-Kante). Keine `DC-*` — falls das
exclude-sections-Scoping einen Produkt-Umfang ergibt, trägt der
**Folge-Slice** das Anforderungs-Delta.

**Berührte Spec-Stellen:** — *(Config-Evaluierung ist kein Spec-Stratum)*.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-29.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die fünf `.d-check.yml`-Positionen, mit denen die Schwester über
die Baseline-Vorlage hinausgeht, sind je für d-checks eigenen Betrieb
bewertet — adoptiert (Config-Update) oder abgelehnt (Begründung):

1. **welle-Klasse** samt Präfix `welle-` und den Regeln `adr→welle` /
   `spec-straten→welle` — der Wellen-Betrieb kehrt mit welle-91 zurück;
   die Matrix-Klasse ist die Gate-Deckung der Verweis-Richtung.
2. **Klasse `aussen`** (First-Match, letzte) + Regel `spec-straten→aussen` —
   schließt die Referenz-Richtung nach außen, nicht nur nach unten.
3. **Klasse `adaptionsblock`** + `token: 'MR-\d{3}'` + Regel
   `spec-straten→adaptionsblock` — fängt nackte MR-Kennungen in den Spec-Straten.
4. **Segment-tolerante ids** (`ADR-([A-Z]+-)?\d{4}`) — d-checks ids-Muster
   ist heute `'ADR-\d{4}'`.
5. **exclude-sections-Scoping** — die Schwester dokumentiert als Grenze,
   dass exclude-sections global über alle Klassen gilt (kein je-Klasse/Regel-
   Scoping); ihr Fall: ADR-Geschichte braucht die Slice-Ausnahme, das
   Spec-Historie gerade nicht. Bewertung: Produkt-Kandidat (Scoping je
   Klasse) oder Ablehnung.

**Ausdrücklich NICHT in diesem Slice:**

- **Kein exclude-sections-Scoping-Umbau ohne ADR** — ergibt die Bewertung
  einen Produkt-Umfang, trägt ihn ein Folge-Slice mit `DC-*`/ADR.
- **Keine Übernahme der Schwester-Kommentarteile** — die Begründungen sind
  ihr Dokument; d-checks Config führt eigene, am eigenen Bestand geprüfte.

## 2. Definition of Done

- [x] Je Position (5) eine belegte Entscheidung: adoptiert (Config-Update,
      gates grün, Befund-Effekt gegen den Bestand geprüft) oder abgelehnt
      (Begründung).
- [x] Bei adoptierten Positionen: die Befund-Wirkung gegen den eigenen
      Bestand gemessen (d-check-Lauf vor/nach, Differenz notiert). — *keine
      Position wurde adoptiert; die Messung (40 Befunde) trägt die
      Positionen 2+3 — Grundsatz bejaht, Umfang an slice-248 delegiert.*
- [x] `make gates` grün.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.d-check.yml` (nach Evaluierung) | update | adoptierte Positionen |
| dieser Plan | update | etwaige Anker-Korrekturen |

## 4. Trigger

**Start** (`open` → `in-progress`): direkt beansprucht — die
Welle-Eröffnung (welle-91, Auftraggeber 2026-09-29) ist die Beanspruchung;
der Nachtlauf-Stand wird bei der Beanspruchung gelesen
([`MR-053`](../../../../../harness/conventions.md#mr-053)).

**Rückführungen:** `in-progress` → `next` (zu groß): ergibt das
exclude-sections-Scoping einen Produkt-Umfang, wird es geschnitten und
dieser Slice schließt mit den übrigen Bewertungen. `in-progress` → `open`:
keiner bekannt.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag (drei Paarungen
wellenlos hier geprüft).

## 6. Risiken und offene Punkte

- Adoptierte Positionen färben den Bestand neu (z. B. adaptionsblock fängt
  nackte MR-Kennungen in den Spec-Straten — Bestands-Funde möglich) —
  DoD 2 misst vor/nach. — **Ausgang:** *(offen)*
- Die welle-Klasse ohne Wellen-Bestand ist ein leerer Fang — sie braucht
  erst Welle-Dateien; Timing der Adoption je nach slice-244/245-Ergebnis. —
  **Ausgang:** *(offen)*

## 7. Closure-Notiz

- **Was hat funktioniert:** die Probe — die Kandidaten-Konfiguration (nach
  dem Snapshot `/tmp/aih-v6.13.0`) lief gegen den Bestand, bevor eine
  Entscheidung fiel: **40 Befunde** — 29 `matrix-forbidden` (28
  `aussen`-Links aus den Straten: 20 im Lastenheft, 7 in der Spezifikation,
  1 in der Sicht; Ziele 14× Konventionsspeicher, 6× AGENTS.md, 2×
  Baseline-Zitat, je 1 Harness-README, Packaging, Carveout, Register, 2× CR;
  dazu 1 nacktes MR-Token) und 11 `matrix-inactive` (ADR-Index ×7, CHANGELOG
  ×3, releasing ×1). Die fünf Entscheidungen:
  (1) **welle-Klasse** — bereits Bestand in schärferer Form
  ([`MR-034`](../../../../../harness/conventions.md#mr-034--die-referenzmatrix-bewacht-auch-die-kante-adr--welle),
  Token `welle-\d{2,}`, beide Regeln); (2) **aussen** — Grundsatz bejaht,
  Bestands-Nachzug trägt [slice-248](../../open/slice-248-matrix-aussen-adaptionsblock.md);
  (3) **adaptionsblock** — dieselbe Lage, im selben Folge-Slice; (4)
  **segment-tolerante ids** — abgelehnt, keine segmentierte ADR-Form im
  Bestand (gemessen: null Treffer), die Weite wäre unbelegte Vorsorge;
  dasselbe Urteil trägt die übrigen Schwester-Abweichungen — der
  adr-Klassen-Glob mit Buchstaben-Segment und der slice-Token-Prefix, je ohne
  Instanz im Bestand; (5)
  **exclude-sections-Scoping** — abgelehnt, kein eigener Fall: die
  Spec-Historie trägt ein anderes Heading als die ADR-Geschichte, die
  globale Ausnahme `[Geschichte]` trifft sie nicht. Grenze der Ablehnung:
  trägt künftig ein Spec-Stratum ein `## Geschichte`, weitet
  `exclude-sections: [Geschichte]` still aus — dann tritt der Fall ein.
- **Was ging anders als geplant:** Position 1 war bereits Bestand — der Plan
  nahm an, die welle-Klasse „kehrt mit welle-91 zurück"; die Kante
  `adr→welle` ist seit
  [`MR-034`](../../../../../harness/conventions.md#mr-034--die-referenzmatrix-bewacht-auch-die-kante-adr--welle)
  bewacht. Die Probe-Ordnung hatte `aussen`
  vor `adaptionsblock` (First-Match) — die Befund-Labels sind darum teilweise
  vertauscht, die Zählung ist gültig; die Ordnungskorrektur trägt slice-248.
  Die Bestands-Wirkung von Position 2+3 überstieg die Adoptions-Schwelle
  dieses Slice — die Rückführungs-Denke des Plans (entworfen für Position 5)
  greift bei 2+3 sinngemäß.
- **Steering-Loop-Eintrag:** keine Verkörperung — keine der Klassen erreichte
  die Schwelle; die First-Match-Lektion (Klassen-Ordnung) ist im Folge-Slice
  verankert, nicht als Regel.
- **Beobachtungs-Register (`../../observations/`):** keine Beobachtung angefallen.
- **Folge-Slices:** [slice-248](../../open/slice-248-matrix-aussen-adaptionsblock.md)
  (matrix-Klassen `aussen` + `adaptionsblock` adoptieren — samt
  Bestands-Nachzug) — ist eine Datei in `open/`.
- **Risiken aus §6:** Risiko 1 (Adoptierte Positionen färben den Bestand
  neu): eingetreten — gemessen (40 Befunde), getragen von slice-248. Risiko 2
  (welle-Klasse ohne Wellen-Bestand ist leerer Fang): entfallen — die Klasse
  ist bereits Bestand, der Fall existiert nicht. Nachtlauf-Stand
  ([`MR-053`](../../../../../harness/conventions.md#mr-053)): beide Läufe grün am
  2026-09-29 (upstream-drift 06:29 UTC, image-scan 09:52 UTC) — nichts zu
  lesen.
- **Drei Paarungen:** (a) Anker — kein `liegt in`-Feld, nichts verkörpert,
  die Paarung trifft nicht zu; (b) Folge-Slice — slice-248 existiert in
  `open/`; (c) Register — keine neuen Einträge, keine Zitate.

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: das
Harness-Werkzeug selbst (`*`, Kürzel `ALL`); bereits deklariert.

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** keine berührte Sub-Area
mit Treffern in offenen Beobachtungen.

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung):
`make nightly-state` wird bei der Beanspruchung gelesen und der Stand in
§7 notiert.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.

### Sub-Area: `*` (Harness-Werkzeug)

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — die fünf Positionen sind durch die
  Schwester dokumentiert, die Messung ist `grep`/Diff gegen den Bestand.
- **Phase-Reife:** Phase 3.
- **Evidenz-/Diskrepanz-Risiko:** Niedrig — je Position eine belegte
  Entscheidung, vor/nach gemessen.
- **Reconciliation-Aufwand:** Keiner — ein Produkt-Umfang trägt sein
  Anforderungs-Delta als Folge-Slice.
