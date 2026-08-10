# Welle welle-70-fence-lexik: Unbalancierter Fence als stiller Grün-Pfad

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** und liegt **flach**
unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach `done/`
(neben ihre `welle-70-results.md`). Der Zustand ist die Verzeichnis-Position —
kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug (Defekt-Behebung an ausgeliefertem
Verhalten).

**Verantwortlich:** pt9912. **Datum:** 2026-08-09.

---

## 1. Welle-Ziel

Einen **ausgelieferten** stillen Grün-Pfad schließen: ein unbalancierter oder
verschachtelter Fenced-Code-Block verschluckt alles hinter sich, und die
Bedingungen der Closure-Note-Struktur laufen darüber grün.

**Zur Ehrlichkeit über diese Welle:** sie bündelt **einen** Slice, und der Grund
ist *nicht* ein Mehr gegenüber seiner Definition of Done. Der Grund ist eine
repo-eigene Kopplung: [`MR-024`](../../../../harness/conventions.md#mr-024--aktuelle-welle-ruhe-marker-im-wellenlosen-zustand-aktive-welle-template-konform)
kennt genau zwei Zustände — laufende Welle mit Struktur-Feldern oder Ruhe-Marker
—, und `make planning-check` erzwingt „kein `slice-*` in `in-progress/` ⟺
Marker". Ein Slice in Arbeit **ohne** benannte Welle ist damit nicht
darstellbar, obwohl das Baseline-Regelwerk Arbeit ohne Wellen-Betrieb
ausdrücklich vorsieht. Die drei Slices, die „ohne Welle" in ihrem Kopf trugen,
sind entsprechend korrigiert. Ob die Kopplung so bleiben soll, gehört bei der
Closure ins Beobachtungs-Register, nicht in diese Welle.

## 2. Trigger (Welle startet)

- [slice-096](slice-096-structure-modul-analyse.md) ist geschlossen; der
  Defekt ist dort als eigener Slice ausgewiesen und **bindende** Vorbedingung
  für die Umsetzungs-Welle.
- Kein WIP-Konflikt: `in-progress/` trug beim Start keinen Slice.

## 3. Closure-Trigger (Welle schließt)

- [slice-101](slice-101-fence-unbalanciert.md) liegt in `done/`.
- Der belegte Reproduktionsfall meldet; der Test ist mutations-echt.
- `make gates` + `make verify-closure-notes` grün; Trigger-Audit durchlaufen.
- Closure-Notiz `done/welle-70-results.md` geschrieben.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| [slice-101](slice-101-fence-unbalanciert.md) | Unbalancierter Fence verschluckt still den Rest des Abschnitts | [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in), [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md) |

## 5. Abhängigkeiten

- **Blockiert:** die structure-Umsetzungs-Welle — deren Slice darf nicht vor
  diesem laufen, sonst erbt das neue Modul den Defekt über die geteilte Mechanik.
- **Wird blockiert von:** nichts.

## 6. Out-of-Scope für diese Welle

- **Das Modul `structure`** — eigener Strang, eigene Welle.
- **Eine allgemeine Markdown-Parser-Ablösung.** [ADR-0042](../../adr/0042-markdown-lexik-folgt-commonmark.md)
  hat gemessen, dass ein echter Parser die Policy-Klasse nicht löst; diese Welle
  ändert die Fence-Lexik gezielt, nicht das Verfahren.
- **Sanierung des Bestands**, falls die Messung unbalancierte Fences findet —
  das wäre ein eigener Slice.

## 7. Closure-Notiz

Geschlossen am 2026-08-10. Die Ergebnis-Notiz der Welle steht — der Baseline-Form
folgend — in einer **eigenen** Datei neben dieser: [`welle-70-results.md`](welle-70-results.md). Sie
trägt, was geliefert wurde, was funktioniert hat, was anders lief, den
Lese-Schritt am Beobachtungs-Register und die Verifikation.

Zwei neue Register-Einträge sind dabei entstanden: **BEO-003** und
**BEO-004**, letzterer bei Zähler 3 und damit an der Verkörperungs-Schwelle.

Diese Plan-Datei hält nur noch fest, **dass** die Welle geschlossen ist; ihr
Zustand ist die Verzeichnis-Position.
