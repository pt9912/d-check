# slice-246: `.d-check.yml`-Positionen aus ai-harness-init evaluieren

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle as State Machine.

**Welle:** welle-91.

**Bezug:** [`welle-91`](../welle-91-adoption-ai-harness-init.md),
[`MR-073`](../../../../harness/conventions.md#mr-073),
[`MR-034`](../../../../harness/conventions.md#mr-034--die-referenzmatrix-bewacht-auch-die-kante-adr--welle)
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

- [ ] Je Position (5) eine belegte Entscheidung: adoptiert (Config-Update,
      gates grün, Befund-Effekt gegen den Bestand geprüft) oder abgelehnt
      (Begründung).
- [ ] Bei adoptierten Positionen: die Befund-Wirkung gegen den eigenen
      Bestand gemessen (d-check-Lauf vor/nach, Differenz notiert).
- [ ] `make gates` grün.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.d-check.yml` (nach Evaluierung) | update | adoptierte Positionen |
| dieser Plan | update | etwaige Anker-Korrekturen |

## 4. Trigger

**Start** (`open` → `in-progress`): direkt beansprucht — die
Welle-Eröffnung (welle-91, Auftraggeber 2026-09-29) ist die Beanspruchung;
der Nachtlauf-Stand wird bei der Beanspruchung gelesen
([`MR-053`](../../../../harness/conventions.md#mr-053)).

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

*(gefüllt vor dem `git mv` nach `done/`)*

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

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
