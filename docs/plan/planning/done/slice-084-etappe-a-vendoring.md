# Slice slice-084: Regelwerk-Migration Etappe A — Vendoring auf v5.0.0

**Status:** In Arbeit (welle-67).

**Welle:** welle-67-baseline-v500-migration (Trigger: Abnahme von slice-083 am
2026-08-01; erste Umsetzungs-Etappe).

**Bezug:** Umsetzung von **Etappe A** des in
[slice-083](../done/slice-083-regelwerk-v500-migration-analyse.md) §2.7
abgenommenen Migrations-Schnitts (Baseline `v1.4.0` → `v5.0.0`). Betrifft das
Materialisierungs-Skript, den vendorten Baseline-Baum und die Pin-/Pointer-
Fundstellen. **Kein Change Request**, **kein ADR**, **kein Release** — reine
Harness-/Konventions-Änderung.

**Autor:** pt9912. **Datum:** 2026-08-01.

---

## 1. Ziel

Der netzlose Boden der Migration: die adoptierte Baseline von `v1.4.0` auf
`v5.0.0` **vendoren**, sodass ab hier jeder Schritt (Etappen B–D) netzlos gegen
die vendorte v5.0.0-Quelle arbeitet. Danach ist der Alt-Stand entfernt, Pin und
Pointer stehen auf `v5.0.0`, die entfallenen Quellzeiger und die Reviewer-Skill-
Fundstelle sind neu geschrieben, die historischen Regelwerk-Pfad-Verweise
abgefangen. **Nur Vendoring + Pointer** — der `MR-*`-Datei-Umbau (Etappe C) und
die Form-Konformität (Etappe D) bleiben ausdrücklich draußen.

## 2. Entscheidungen / Regel

Die acht Schritte stehen detailliert in
[slice-083](../done/slice-083-regelwerk-v500-migration-analyse.md) §2.7 „Etappe A".
Kern-Entscheide dieses Slice:

- **Skript-Referenz u-boot.** `fetch-baseline-cache.sh` wird auf das
  self-contained v5.0.0-Bundle-Layout gehoben (tolerantes Entpacken **beider**
  Bäume, Manifest über den **tatsächlichen** Bestand, Under-Copy-Barriere,
  Currency-Teil auf die Release-Liste). Vorbild ist u-boots Fassung — Anregung,
  kein Kanon; bei Konflikt gilt das Bundle.
- **Alt-Stand entfernen — atomar.** `.harness/baseline/v1.4.0/` fällt
  (Pin-Adaptions-Pflicht); die pin-gebundenen Verweise wandern **im selben grünen
  Zustand** mit, sonst `target-missing` im Push-Tip.
- **Entfallene Quellzeiger umschreiben, nicht retargeten** (slice-083 §2.2 Bruch 6
  / Etappe A4): der `agents-regelwerk.md`-Eintrag (v2.0.0 retired) und der
  Kurs-Konventionsdatei-Zeiger (v5.0.0 entfallen) werden auf das vendorte
  Modul-Bundle umgeschrieben bzw. entfernt — **nicht** auf tote Ziele umgehängt.
- **Historische Verweise per Tombstone** (§2.2 #4): die drei eingefrorenen
  Regelwerk-Pfad-Verweise (eine immutable ADR + zwei `done/`-Slices) über ein
  gescoptes `ignore-refs` abfangen — kein Editieren eingefrorener Doku.
- **Vendoring-Adaption** unter der **nächsten freien** `MR`-Nummer deklarieren (das
  v5.0.0-Layout: self-contained Bundle, beide Bäume, entfallenes
  `lab-templates.zip`); die Nummern-Kollision mit der Baseline-Vendoring-Adaption
  bleibt Etappe C.
- **Anker-Kompatibilität ist NICHT hier.** Der Voll-Slug-`<a id>`-Block gehört zum
  `conventions.md`-Datei-Umbau (Etappe C); Etappe A lässt den inline-`MR`-Block
  unangetastet.

## 3. Definition of Done

- [x] `fetch-baseline-cache.sh` entpackt das v5.0.0-Bundle korrekt (beide Bäume,
  Manifest über tatsächlichen Bestand, Under-Copy-Barriere); `--verify` grün.
- [x] `.harness/baseline/v5.0.0/` mit `regelwerk/` + `templates/` + `SHA256SUMS`
  vendored; `.harness/baseline/v1.4.0/` entfernt.
- [x] `§Baseline`-Stand + die sechs Live-Fundstellen auf `v5.0.0`; entfallene
  Quellzeiger umgeschrieben; Reviewer-Skill retargetet.
- [x] Historische Verweise per `ignore-refs`/Tombstone abgefangen (kein
  eingefrorenes Doc editiert).
- [x] Vendoring-`MR` deklariert (Index-Eintrag + Begründung + `Ersetzt-Baseline-Regel`).
- [x] `make gates` grün; unabhängiger Frischkontext-Review (slice-083 §4: für die
  inhaltlichen Etappen Pflicht).

## 4. Risiken / offene Punkte

- **Netz für die Materialisierung.** Der Bundle-Download ist Setup (wie ein
  Image-Pull), nicht Teil der netzlosen Analyse; der `--verify`/Selbst-Scan bleibt
  netzlos.
- **Atomarität.** Alt-Stand-Entfernung + Pointer-Umzug müssen im selben grünen
  Zustand landen — die PR-/Push-CI prüft den Push-Tip.
- **Skript-Rewrite ist Fremd-Layout-Logik.** u-boots Fassung ist Vorbild, nicht
  kanonische Quelle; die Bundle-Struktur ist die Wahrheit.

## 5. Trigger

Abnahme von slice-083 (2026-08-01): erste Umsetzungs-Etappe des abgenommenen
Migrations-Schnitts A–D.

## 6. Sub-Area-Modus-Begründung

GF (Repo-Default): Doc/Prozess führt. Berührt die Sub-Areas *Harness-Tooling*
(Skript) und *Harness/Konventionen* (Pin/Pointer) — greenfield gewachsen, ohne
Brownfield-Spec.

## 7. Closure-Notiz (nach `done/`)

**Abgeschlossen 2026-08-02.** Etappe A der Baseline-Migration `v1.4.0` → `v5.0.0`
(Vendoring + Pin/Pointer). Reine Harness-/Konventions-Änderung, kein Release/ADR.

**Umgesetzt.** Das Skript `fetch-baseline-cache.sh` ist aufs self-contained
v5.0.0-Bundle gehoben (tolerantes Entpacken, Under-Copy-Barriere, Manifest über
beide Bäume; `--check-latest` Currency via Release-Liste + Content-Drift über
beide Bäume); `v5.0.0` vendored (`regelwerk/` 26 + `templates/` 24 + `SHA256SUMS`
50), `v1.4.0` entfernt. Atomare Umschaltung: §Baseline-Pin auf `v5.0.0`, sieben
Live-Pfad-Links retargetet, die zwei entfallenen Quellzeiger
(`agents-regelwerk.md` im Kurs abgelöst, Kurs-`konventionen.md` in acht
`grundlagen-*` gesplittet) umgeschrieben statt retargetet, Reviewer-Skill auf
`grundlagen-referenz-richtung.md`, die Vendoring-/Pin-Hebungs-Adaption im
Konventionsspeicher deklariert. Drei eingefrorene Verweise (eine immutable ADR +
zwei `done/`-Slices) auf den entfernten `v1.4.0`-Pfad via geteiltem
`ignore-refs`-Ventil quell-skopiert getombstoned.

**Review.** Unabhängiger Frischkontext-Review (Report unter `docs/reviews/`):
**abnahmereif**, HIGH 0 / MEDIUM 0 / LOW 2 / INFO 1; Per-Commit-Grün + Atomizität
per Detached-Worktree empirisch bestätigt, Tombstone load-bearing, kein stilles
Grün im Skript-Pfad. LOW-2b + INFO eingearbeitet.

**Nach Etappe C übergeben (im Konventionsspeicher-Eintrag getrackt):**
- LOW-1: die Under-Copy-Barriere filtert den Regelwerk-Baum auf `*.md` (heute
  inert — das v5.0.0-Regelwerk ist rein Markdown; ein künftiger
  Nicht-Markdown-Bump bliebe ungeprüft).
- LOW-2a: die MR-Prosa führt intern noch `<tag> = v1.4.0`- und
  `…/v1.4.0/…`-Beispiele (interne Drift, kein Gate-Bruch — die Prosa-Angleichung
  ist die Konventionsspeicher-Etappe).

**Nächster Schritt.** Etappe B (Modul-Delta lesen) als eigener Slice.
