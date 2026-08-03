# Slice slice-092: `## Aktuelle Welle` auf die Template-Struktur-Felder (aktive-Welle-Form)

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** welle-68-planning-roadmap-harness (erster Slice).

**Bezug:** Nutzer-Entscheid „Template strikt" für `## Aktuelle Welle` — die
Baseline-`roadmap.template.md` führt dort die Struktur-Felder (Welle-ID · Start ·
Geplantes Ende · Closure-Trigger). Erkenntnis: eine **aktive** Welle trägt diese Felder
**ohne** Ruhe-Marker, und `planning-check` ist grün (`hasActive == hasSlices`) — **kein**
`planning`-Modul-Umbau nötig. Dieser Slice hebt die Aktive-Welle-Form auf die Template-
Felder und verfeinert die deklarierte Adaption entsprechend.

**Autor:** pt9912. **Datum:** 2026-08-03.

---

## 1. Ziel

`## Aktuelle Welle` erreicht die Template-Form: **während eine Welle läuft** trägt der
Abschnitt die Struktur-Felder (Welle-ID · Start · Geplantes Ende · Closure-Trigger) statt
der bisherigen Frei-Prosa; **im wellenlosen Zustand** bleibt der vom `planning`-Modul
erzwungene Ruhe-Marker. Die Adaption im Konventionsspeicher wird auf diese **zwei
Zustände** verfeinert.

## 2. Vorgehen

1. **Roadmap-Aktive-Welle-Form.** `## Aktuelle Welle` in die Struktur-Feld-Form heben
   (Welle-ID: welle-68… · Start · Geplantes Ende · Closure-Trigger + Slice-Liste) —
   **ohne** die Zeichenfolge des Ruhe-Markers im Abschnitt (sonst `planning-drift`, da
   ein Slice in `in-progress/` liegt). Erfolgt mit der welle-68-Eröffnung.
2. **Adaption verfeinern.** Den Konventionsspeicher-Eintrag zur `## Aktuelle Welle`-Form
   auf die **zwei Zustände** präzisieren: aktive Welle → Template-Struktur-Felder;
   wellenlos → Ruhe-Marker (gate-erzwungen). Festhalten, dass **kein** `planning`-Modul-
   Umbau nötig ist (die aktive Welle löst die Feld-Form gate-konform).
3. **Gate.** `make gates` (inkl. `planning-check` — Felder ohne Ruhe-Marker + Slice
   in-progress ⇒ grün) + `make adr-check` grün; unabhängiger Frischkontext-Review.

## 3. Abnahme-Punkte

_Keine offenen — der Nutzer-Entscheid („Template strikt", aktive-Welle-Form ohne
Modul-Umbau) ist die Vorgabe._

## 4. Definition of Done

- [ ] `## Aktuelle Welle` trägt die Template-Struktur-Felder (Welle-ID/Start/Geplantes
  Ende/Closure-Trigger + Slices), **ohne** Ruhe-Marker-Zeichenfolge, während welle-68 läuft.
- [ ] Der Konventionsspeicher-Eintrag zur `## Aktuelle Welle`-Form ist auf die zwei
  Zustände verfeinert (aktive = Felder, wellenlos = Ruhe-Marker; kein Modul-Umbau).
- [ ] `make gates` + `make adr-check` grün; unabhängiger Frischkontext-Review.

## 5. Risiken / offene Punkte

- **`planning`-Gate-Falle:** die Zeichenfolge „Keine aktive Welle" darf im
  `## Aktuelle Welle`-Block **nicht** vorkommen, solange ein Slice in `in-progress/`
  liegt (das Modul sucht sie als literalen Teilstring) — Adaptions-Erklärung ohne das
  wörtliche Marker-Zitat formulieren.

## 6. Trigger

welle-68-Eröffnung; Nutzer-Entscheid „Template strikt" für `## Aktuelle Welle`.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** berührt *Harness/Prozess*-Doku (Roadmap-Form + Konventionsspeicher)
  — greenfield, GF.
- **Offene Beobachtungen sichten:** `observations.md` = `— keine —`; nichts zu
  berücksichtigen.

## 8. Sub-Area-Modus-Begründung

GF (Repo-Default): Doc/Prozess führt. Berührt die Roadmap-Form + den Konventionsspeicher;
greenfield-Angleich an die Template-Form.

## 9. Closure-Notiz (nach `done/`)

_Ausstehend._
