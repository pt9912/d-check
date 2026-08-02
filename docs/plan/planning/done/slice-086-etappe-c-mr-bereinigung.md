# Slice slice-086: Regelwerk-Migration Etappe C — MR-Bereinigung + Datei-Migration

**Status:** In Arbeit (welle-67).

**Welle:** welle-67-baseline-v500-migration (dritte Umsetzungs-Etappe, nach
[slice-085](../done/slice-085-etappe-b-modul-delta.md)).

**Bezug:** Umsetzung von **Etappe C** des in
[slice-083](../done/slice-083-regelwerk-v500-migration-analyse.md) §2.7
abgenommenen Schnitts, **verbindlich gemacht** durch die
[slice-085](../done/slice-085-etappe-b-modul-delta.md)-Finding-Liste (§3.1
C-1…C-8 + die §3.4-Korrekturen). Berührt die **Konventions-Identität** des Repos.
**Kein Change Request**, **kein ADR**, **kein Release** — Harness-/Konventions-
Änderung.

**Autor:** pt9912. **Datum:** 2026-08-02.

---

## 1. Ziel

`harness/conventions.md` von der Inline-Sammeldatei auf die v5.0.0-**Default-Form**
heben: **Index + eine Datei je Adaption** (`harness/conventions/MR-<NNN>-<titel>.md`),
aufgelöste Einträge → `conventions/done/` (Zustand = Verzeichnis-Position), die
neuen Pflichtfelder je Eintrag, **Forks reklassifiziert** (nach dem in slice-085
§3.4 korrigierten Kriterium), entwertete Adaptionen gestrichen/verschmolzen, die
**173** `conventions.md#mr-…`-Links via **Voll-Slug-Anker-Kompatibilitäts-Block**
erhalten und den d-check-`matrix`-Scope gegen die **8×8-Matrix** als bewusste
Grenze deklariert (C-4). **Nur die Konventions-Identität** — die Form-Struktur
(die 11 D-Findings) ist Etappe D; die **Historie-Provenance-Enforcement (C-3)** ist
als Code-Feature in einen **Folge-Produkt-Slice** herausgeschnitten (§3).

## 2. Vorgehen (slice-083 §2.7 Etappe C, gefüttert durch die C-Findings)

1. **Verbindlich machen** — die §2.4-MR-Einschätzungen gegen die slice-085-Findings
   festziehen; besonders die **korrigierte** Fork-Frage (§3.4: die Guard-Härtungs-
   Adaptionen sind keine Forks; Fork = *pauschale* Nichtanwendbarkeit).
2. **In Dateien splitten** — die Inline-Adaptionen aus `conventions.md` je in eine
   `harness/conventions/MR-<NNN>-<titel>.md` überführen; im Index bleibt die
   Adoptions-Erklärung + eine Zeile je aktiver Adaption; aufgelöste per `git mv`
   nach `conventions/done/`. **Anker-Erhalt (C-5/Review-N-1):** `conventions.md`
   behält je **von eingefrorener Doku** (immutable ADRs, `done/`-Slices,
   `docs/reviews/`) referenziertem MR einen **Voll-Slug-`<a id>`-Anker** — auch für
   die aufgelösten — in einem eigenen **Anker-Kompatibilitäts-Block**, unabhängig
   vom aktiv/`done`-Schnitt. So lösen die 173 Links **ohne Retarget, ohne ADR-Edit**
   auf; der Block wird als eigener `MR` deklariert (migrationsspezifisch).
3. **Neue Pflichtfelder** (C-6) je Eintrag: `Ersetzt-Baseline-Regel` (genau **eine**
   v5.0.0-Regel, Anker-Link in den vendored Baum), `Status: Accepted`, für
   Ablösungen `Löst auf` + `Ausgelöst durch Baseline-Stand` bzw. `(schärft …)` im
   Titel.
4. **Forks reklassifizieren** (C-Schritt 4 / §3.4) — **Abnahme-Punkt (§3.1)**.
5. **Entfall/Verschmelzung** (C-1/C-2/C-8) — die entwerteten Adaptionen (eigene
   Spec-Straten-Schicht → Default; Spec-Decke → Default; Template-Freiheit/-Cache)
   streichen bzw. zu Provenienz abrüsten; Vendoring/Currency-Audit in den Default
   überführen.
6. **Nummern-Kollision** (§2.4) — **Abnahme-Punkt (§3.1)**.
7. **Prosa entpinnen + Referenzrichtung** — die interne `v1.4.0`-Prosa (slice-084
   LOW-2a: `<tag>=v1.4.0` in den Alt-Einträgen) auf `v5.0.0` ziehen; die
   Referenzrichtung gegen die **8×8-Matrix** stellen und den d-check-`matrix`-Scope
   (Spec-Decke + markierte ADR→Slice-Kante) als **bewusste Grenze** in der
   Referenzrichtungs-Adaption deklarieren (C-4; die Planungs-Kanten ADR→{Carveout,
   Welle,Roadmap}/Slice→Roadmap bleiben unbewacht, dokumentiert). **Die
   Historie-Provenance-Enforcement (C-3) ist herausgeschnitten** (§3 Abnahme-Punkt
   3 → Folge-Produkt-Slice).
8. **Gate** — `make gates` + `make adr-check` grün; die neuen `conventions/`-Dateien
   erfüllen die `ids`-Linkpflicht; **keine** `Accepted`-ADR inhaltlich berührt (die
   173 Links via Anker-Block, die eingefrorene ADR-Alt-Pfad-Fundstelle via den
   Tombstone aus Etappe A); unabhängiger Frischkontext-Review.

## 3. Abnahme-Punkte (Maintainer-Entscheid vor/während der Umsetzung)

1. **Fork-Klassifikation** (Vorgehen 4). Nach dem in slice-085 §3.4 am Quelltext
   geschärften Kriterium (*pauschale* Nichtanwendbarkeit der Baseline) ist **keine**
   d-check-Adaption ein echter Fork; die Guard-Härtungs-Adaptionen benennen
   `Ersetzt-Baseline-Regel` auf `grundlagen-durchsetzungsschicht`. **Entscheid:**
   (a) alle als **Adaptionen** mit `Ersetzt-Baseline-Regel` führen (empfohlen,
   quell-belegt), oder (b) die Rest-Kandidaten als Forks behalten.
   → **Entschieden 2026-08-02: (a)** — alle als Adaptionen, **keine Forks**; jeder
   Eintrag benennt eine `Ersetzt-Baseline-Regel`.
2. **Nummern-Kollision** (Vorgehen 6, §2.4). Eine `MR`-Nummer ist in der Baseline
   anders belegt als in d-checks aufgelöster Vorgänger-Adaption. **Entscheid:**
   eigene Nummern behalten (Provenienz) oder an die Baseline angleichen.
   → **Entschieden 2026-08-02: eigene Nummern behalten (Provenienz)** —
   [`MR-003`](../../../../harness/conventions.md#mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh)
   bleibt der historische Bootstrap-Sensor, das Vendoring bleibt
   [`MR-023`](../../../../harness/conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
   (die Nummern tragen die 173 Links; die Baseline-Nummer ist ein
   Vorlagen-Beispiel, keine Regel).
3. **Historie-Provenance-Spec-Bereinigung** (Vorgehen 7 / C-3). Die ADR-/Slice-
   Verweise aus den Spec-§7-Historien entfernen **berührt die kanonische Spec**.
   **Entscheid:** jetzt bereinigen (Baseline-Default sticht) — und ob der interne
   Slice-Bezug nach v5.0.0-Verfahren in die Slice-Closure-Notiz wandert.
   → **Revidiert 2026-08-02 (nach Messung): in einen Folge-Produkt-Slice
   herausgeschnitten.** Die Messung (`exclude-sections` testweise entfernt) zeigte:
   `matrix.exclude-sections` ist **global**, nicht pro Klasse — global entfernen
   bräche **65 Befunde in `Accepted`-immutablen ADRs** (0022–0046, deren
   `Geschichte` Slices als marker-lose Provenance nennt). C-3 braucht ein
   **d-check-Code-Feature** (per-Klasse/-Pfad exclude-sections) **+** die gekoppelte
   Spec-§7-Bereinigung → eigener Produkt-Slice (CR + ADR + Slice am `matrix`-Modul),
   **nicht** diese Doc-Migrations-Etappe C. Die anderen zwei Entscheide bleiben.

## 4. Definition of Done

- [ ] `conventions.md` = Index (Adoptions-Erklärung + Zeile je aktiver Adaption) +
  Anker-Kompatibilitäts-Block; je Adaption eine Datei
  `harness/conventions/MR-<NNN>-<titel>.md`, aufgelöste unter `conventions/done/`.
- [ ] Neue Pflichtfelder je Eintrag; Forks reklassifiziert (Abnahme-Punkt 1);
  entwertete gestrichen/verschmolzen; Nummern-Kollision aufgelöst (Abnahme-Punkt 2).
- [ ] Referenzrichtung gegen die 8×8-Matrix gestellt; d-check-`matrix`-Scope als
  bewusste Grenze deklariert (C-4). (C-3-Enforcement herausgeschnitten →
  Folge-Produkt-Slice, Abnahme-Punkt 3.)
- [ ] Die 173 `conventions.md#mr-…`-Links lösen auf (Anker-Block); **keine**
  `Accepted`-ADR inhaltlich berührt.
- [ ] `make gates` + `make adr-check` grün; unabhängiger Frischkontext-Review.

## 5. Risiken / offene Punkte

- **Größter, identitäts-berührender Schnitt.** 24 Einträge, 173 Links, 12 immutable
  ADRs — der Anker-Kompatibilitäts-Block (Vorgehen 2) ist der load-bearing Kern;
  ein Fehler dort bricht `anchors` repo-weit.
- **Abnahme-Punkte sind Voraussetzung** für Vorgehen 4/6/7 — ohne Entscheid keine
  Umsetzung dieser Schritte.
- **Kanonische Spec** (Vorgehen 7): der einzige Schritt, der `spec/`-Dateien ändert
  — vorsichtig, review-pflichtig.

## 6. Trigger

Abschluss von slice-085 (Etappe B): die Finding-Liste macht C verbindlich.

## 7. Sub-Area-Modus-Begründung

GF (Repo-Default): Doc/Prozess führt. Berührt *Harness/Konventionen* (der
Konventionsspeicher) und punktuell die *Spec*-Sub-Area (Vorgehen 7) — greenfield,
ohne Brownfield-Spec.

## 8. Closure-Notiz (nach `done/`)

_Ausstehend._
