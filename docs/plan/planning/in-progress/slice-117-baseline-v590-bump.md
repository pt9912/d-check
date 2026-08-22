# Slice slice-117: Baseline-Pin-Hebung auf `v5.9.0` (zwei Stufen)

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-81-zustandsfelder](../welle-81-zustandsfelder.md) (zugeordnet
bei der Eröffnung).

**Bezug:** [`MR-028`](../../../../harness/conventions.md#mr-028) (der abzulösende
Pin-Eintrag), [`MR-021`](../../../../harness/conventions.md#mr-021) (in-repo-Verweise
auf die vendorte Baseline sind pin-gebunden),
[`MR-023`](../../../../harness/conventions.md#mr-023) (self-contained
Bundle-Layout, bleibt), [`MR-013`](../../../../harness/conventions.md#mr-013)
(Lifecycle-Move bündelt gekoppelte Verweise); Beobachtung BEO-008
(Drei-Klassen-Zensus einer Hebung).

**Berührte Spec-Stellen:** — (Harness-/Konventions-Ebene; keine Spec-Zeile).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

Der vendorte Baseline-Bestand steht auf `v5.9.0`: das self-contained Bundle
beider Bäume ist materialisiert und verifiziert, der Baum der Vorgänger-Version
entfernt, der Pin in §Baseline und §Adoptierte Konventions-Quellen umgestellt,
und **alle** pin-gebundenen Verweise sind retargetet — nach dem
Drei-Klassen-Zensus, nicht nach Gefühl. Der Bump ist zwei Stufen weit; das
inhaltliche Delta ist in der Welle gemessen und wird hier nur vendored, nicht
schon angewandt.

## 2. Vorgehen

1. **Vendoren:** `tools/harness/fetch-baseline-cache.sh v5.9.0` (Netz, Anlass
   Bump), danach `--verify` offline. Der Vorgänger-Baum wird entfernt — ein
   Pin, eine netzlose Lese-Form, wie bei jeder Hebung zuvor.
2. **Nachfolge-Adaption** aus der vendorten Vorlage anlegen (die Nummer ist
   die nächste freie des dichten Zählraums, vergeben beim Anlegen): was sich
   ändert (zwei Stufen, gemessenes Delta), Begründung, Auflösungs-Trigger; der
   Vorgänger wandert per `git mv` nach `conventions/done/` samt seiner
   Link-Tiefen-Fixes, die Index-Zeile von „Aktive" nach „Aufgelöste
   Adaptionen" mit ihren Voll-Slug-Ankern.
3. **Drei-Klassen-Zensus** (BEO-008): (1) `baseline/<tag>`-**Pfad**-Verweise
   (grep-bar, gate-gedeckt), (2) Release-/Tree-**URLs** mit dem Tag
   (`releases/tag/`, `releases/download/`, `tree/` — kein Gate deckt sie),
   (3) **Prosa-/Ellipsen-Pins** in MR-Körpern und Fließtext. Jede Klasse
   einzeln durchgehen und das Ergebnis je Klasse notieren.
4. **Eingefrorene Verweise:** was auf den alten Baum zeigt und nicht wandern
   darf (immutable ADRs, `done/`-Slices, Review-Reports), über das geteilte
   Referenz-Ventil quell-skopiert ausnehmen — lebende Verweise werden
   retargetet.
5. Unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Anwendung der neuen Regel** — Kopfzeilen, Register und Drift-Log
  bleiben unverändert (slice-118 bis slice-120).
- **Kein Konformitäts-Anspruch:** dieser Slice hebt den Pin, er behauptet
  keine Konformität.
- **Kein Produkt-Code, kein Release.**

## 4. Definition of Done

- [ ] `.harness/baseline/v5.9.0/` vollständig vendored, `--verify` grün,
      Vorgänger-Baum entfernt; `--check-latest` meldet keinen neueren Release.
- [ ] Nachfolge-Adaption angelegt, Vorgänger in `conventions/done/`, Index-Zeile
      umgezogen (beide Anker mit), §Baseline und §Adoptierte
      Konventions-Quellen auf den neuen Pin.
- [ ] Drei-Klassen-Zensus durchgeführt und **je Klasse** dokumentiert.
- [ ] `make gates` grün; unabhängiger Review; Closure-Notiz; Register
      gesichtet.

## 5. Abnahme-Punkte / Risiken

- **Die Klassen 2 und 3 sind gate-blind** und in zwei aufeinanderfolgenden
  Hebungen je als Review-Auflage nachgezogen worden. Der Zensus ist die
  Antwort darauf — er wird abgearbeitet, nicht erinnert. — **Ausgang:** *(bei
  Closure)*
- **Zwei Stufen auf einmal:** das Delta ist gemessen, aber die Zwischenstufe
  wird nie vendored — falls eine Regel in `v5.8.0` entstand und in `v5.9.0`
  wieder verschwand, sieht der Zensus sie nicht. Gegenmittel: der Delta-Diff
  läuft über **beide** Stufen einzeln. — **Ausgang:** *(bei Closure)*
- **Eingefrorene Verweise auf den alten Baum** brechen beim Entfernen. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): Wellen-Eröffnung (eingetreten) — erster
Slice der Welle.

**Rückführungen:** `in-progress` → `next`, falls das Bundle-Layout sich
geändert hat (dann ist es kein Pin-Nachtrag mehr, sondern eine
Layout-Adaption).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-/Konventions-Doku (`harness/`, `.harness/`, GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-22): BEO-008 ist
  **einschlägig** und steuert §2 Schritt 3; BEO-006/009/010 sind
  Arbeitsregeln; BEO-002 wirkt als Spiegel-Pflicht.

Slice-ID: slice-117. Betroffene IDs:
[`MR-021`](../../../../harness/conventions.md#mr-021),
[`MR-023`](../../../../harness/conventions.md#mr-023),
[`MR-028`](../../../../harness/conventions.md#mr-028). Module:
Konventionsspeicher, vendorte Baseline, Autoritäts-Doku. Gates:
`make doc-check` (eng), `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Adoptions-Pflege nach etablierter
Prozedur; kein Legacy-Import.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
