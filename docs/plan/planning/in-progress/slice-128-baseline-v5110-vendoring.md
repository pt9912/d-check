# Slice slice-128: Etappe A — Bundle `v5.11.0` vendoren und den Pin heben

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-83-baseline-v5110-migration](../welle-83-baseline-v5110-migration.md)
(zugeordnet bei der Eröffnung).

**Bezug:** [`MR-011`](../../../../harness/conventions.md#mr-011) (Pin auf
Release-Tag) und die Kette bis
[`MR-029`](../../../../harness/conventions.md#mr-029);
[`MR-021`](../../../../harness/conventions.md#mr-021) (pin-gebundene Verweise);
[`MR-023`](../../../../harness/conventions.md#mr-023) (Bundle-Layout);
Baseline-Regelwerk
[`modul-02-harness-bootstrap.md` §Freshness-Audit](../../../../.harness/baseline/v5.9.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2).

**Berührte Spec-Stellen:** — (Harness-Vendoring und Konventionsspeicher; keine
Anforderung, kein Spec-Stratum, kein Produkt-Code).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

Der vendorte Baum trägt `v5.9.0` (Kurs-Welle 86). Diese Etappe hebt ihn auf
`v5.11.0` (Kurs-Welle 94) — **nur den Pin und seine Verweise**, ohne eine
einzige Regel anzuwenden. Was der neue Stand inhaltlich verlangt, beantwortet
[slice-129](../open/slice-129-baseline-v5110-delta-audit.md); wer beides mischt, kann
später nicht mehr trennen, ob eine Änderung aus dem Bump oder aus dem Audit kam.

## 2. Vorgehen

1. **Bundle materialisieren** nach `.harness/baseline/v5.11.0/`
   (`{regelwerk,templates}/` + `SHA256SUMS`) über
   `tools/harness/fetch-baseline-cache.sh`; Integrität **offline** verifizieren
   (`--verify`). Kein Handanlegen an den entpackten Bäumen.
2. **Pin heben** in [`harness/conventions.md`](../../../../harness/conventions.md)
   §Baseline und §Adoptierte Konventions-Quellen; neuer `MR-`Eintrag als
   nächster Schritt der Pin-Serie, Index-Zeile ergänzt, Vorgänger nach
   `conventions/done/` mit seiner Nachfolger-Zeile.
3. **Pin-gebundene Verweise heben** ([`MR-021`](../../../../harness/conventions.md#mr-021)) —
   und zwar nach der **Drei-Klassen-Prüfung** aus `BEO-008`, nicht nur per
   Pfad-`grep`: Pfad-Verweise, Release-/Tree-**URLs** mit dem Tag, und
   **Prosa-/Ellipsen-Pins**. Dazu die Gegenprobe, ob eine Stelle über die
   **Gegenwart** oder über die **Vergangenheit** spricht — eine gehobene
   Vergangenheits-Aussage ist ein neuer Fehler.
4. **Alt-Baum entfernen** und prüfen, ob eingefrorene Verweise darauf zeigen
   (immutable ADRs, `done/`-Slices) — falls ja, quell-skopiert über
   `ignore-refs` ausnehmen statt die eingefrorene Doku zu editieren.
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Regel des neuen Stands anwenden.** Auch nicht die
  Vollständigkeits-Zusage aus Kurs-Welle 94, obwohl wir eine Verletzung bereits
  kennen — sie gehört [slice-127](../next/slice-127-claude-md-pointer.md).
- **Kein Delta-Audit.** Das ist Etappe B.
- **Keine Zwischenstufe über `v5.10.0`.**

## 4. Definition of Done

- [ ] `.harness/baseline/v5.11.0/` liegt vollständig; `SHA256SUMS` offline
      verifiziert (Exit explizit).
- [ ] Pin gehoben, `MR-`Eintrag angelegt, Index-Zeilen (aktiv **und**
      aufgelöst) nachgezogen.
- [ ] Verweis-Hebung über **alle drei** Klassen belegt, samt
      Vergangenheits-Gegenprobe.
- [ ] Alt-Baum entfernt; kein Verweis läuft ins Leere.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **`BEO-008` steht bei Zähler 3 und ist genau hier einschlägig.** Die
  Klasse ist bei jeder der letzten drei Pin-Hebungen eingetreten. — **Ausgang:**
  *(bei Closure)*
- **Die Über-Hebung ist die zweite Richtung derselben Klasse:** ein `grep`
  kennt keine Zeitform und hebt Vergangenheits-Aussagen mit. — **Ausgang:**
  *(bei Closure)*
- **Zwei Minors auf einmal** heißt: der Alt-Baum `v5.9.0` verschwindet, ohne
  dass je ein `v5.10.0`-Baum existierte. Verweise auf die Zwischenstufe kann es
  nicht geben — aber die Annahme gehört geprüft, nicht geglaubt. — **Ausgang:**
  *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): sofort — die Welle ist eröffnet,
`in-progress/` frei.

**Rückführungen:** `in-progress` → `next`, falls die Integritätsprüfung des
Bundles fehlschlägt — dann ist es kein Vendoring-, sondern ein Upstream-Problem.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Vendoring (GF), Konventionsspeicher (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23): **`BEO-008`**
  ist die zentrale — Zähler 3, Schwelle erreicht, und ihre benannte
  mechanische Form ist seit slice-122 **baubar** (`versions.patterns`). Ob
  dieser Slice sie baut, ist ein eigener Entscheid; ihn hier zu treffen wäre
  bequem und falsch. **`BEO-002`** für die Ränder der Pin-Hebung.
  **`BEO-011`** für jede Vollständigkeits-Aussage über die gehobenen Verweise.

Slice-ID: slice-128. Betroffene IDs:
[`MR-011`](../../../../harness/conventions.md#mr-011),
[`MR-021`](../../../../harness/conventions.md#mr-021). Module:
Harness-Vendoring, Konventionsspeicher. Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Vendoring nach etablierter, viermal
gefahrener Prozedur.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
