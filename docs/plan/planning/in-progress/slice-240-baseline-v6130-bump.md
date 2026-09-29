# slice-240: Baseline-Pin auf `v6.13.0` — mechanisch, ohne Urteil über den Delta

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — sein Closure-Trigger würde die eigene DoD abschreiben;
derselbe Stand wie die vierzehn Vorgänger-Hebungen (Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht: kein repo-weites Mehr).

**Bezug:** die
[`MR-011`](../../../../harness/conventions.md#mr-011)-Kette
(§Baseline-Pin, [`MR-072`](../../../../harness/conventions.md#mr-072) als
Vorgänger), [`MR-051`](../../../../harness/conventions.md#mr-051)
(Re-Ankern), [`MR-069`](../../../../harness/conventions.md#mr-069)
(Ventil), [`MR-055`](../../../../harness/conventions.md#mr-055)
(Symlinks). Keine `DC-*` — keine Produkt-Anforderung berührt.

**Berührte Spec-Stellen:** — *(der vendierte Baseline-Baum ist kein
Spec-Stratum; `spec/` bleibt unverändert)*.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-29.

---

## 1. Ziel und Abgrenzung

**Ziel:** Den committet vendorten Baseline-Bestand auf
[`v6.13.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v6.13.0)
heben (Re-Vendor, den aktuellen Tag — nicht die vier Zwischen-Tags; das Muster
verkörpert [`MR-072`](../../../../harness/conventions.md#mr-072)), den Pin in
`harness/conventions.md` §Baseline fortschreiben, alle lebenden pin-gebundenen
Referenzen retargeten und die Hebung als neuer MR-Eintrag führen — mit frozen
Klassen und Zitat-Delta gemessen am echten Vorzustand (nach
[`MR-070`](../../../../harness/conventions.md#mr-070),
[`MR-039`](../../../../harness/conventions.md#mr-039)).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Keine Adoption der inhaltlichen Deltas** (Hard Rule als vierte
  Trigger-Audit-Klasse, „Verkörpert ≠ automatisiert", `-RB`-Format,
  ids-Block-Muster für ADR-Kennungen) — jedes davon ist eine
  Konventions-/Produkt-Entscheidung mit eigener Begründung; der Bump ist
  mechanisch (Klasse 3, anderer Vorgang; Adoptieren gehört in eigene
  Folge-Slices, wie es der Vorgänger offen ließ).
- **Frozen-Dateien bleiben byte-stabil** — `done/`-Slices (inkl. slice-224),
  `docs/reviews/`, Register-Evidence, `carveouts/done/CO-002`, CR
  `2026-09-17`, `conventions/done/MR-*` nennen den alten Pin als
  Vergangenheits-Aussage korrekt; eine Über-Hebung fälschte ihren Stand
  (Klasse 2, Bestand bewusst).
- **Kein Release/Tag/GHCR-Lauf** — Release-Prep ist ein eigener Vorgang mit
  eigener Commit-Klasse (Klasse 3); der Bump berührt keine Distribution-Fläche
  (`Dockerfile`, Workflows unverändert).

## 2. Definition of Done

- [ ] Vendierter Baum steht auf `v6.13.0` (`SHA256SUMS` generiert und
      verifiziert), der `v6.9.0`-Baum ist entfernt; §Baseline-Pin (Stand-URL
      + Datum) nachgezogen.
- [ ] Alle deklarierten lebenden Referenz-Klassen nennen `v6.13.0`: Pfad-Verweise
      (`links`/`codepaths`-deckungsgleich), Release-/Tree-URLs, bare
      Versionsnennungen in lebenden Dokumenten; die 7 Alias-Symlinks unter
      `.claude/rules/` lösen gegen `v6.13.0` auf.
- [ ] `make gates` grün; `make baseline-freshness` meldet den Pin aktuell
      (Exit 0).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.harness/baseline/{v6.9.0 → v6.13.0}/` | neu + entfernen | Re-Vendor beider Bäume + `SHA256SUMS`; alter Baum weg |
| `harness/conventions.md` | update | §Baseline-Pin; MR-Index-Zeile für den neuen Eintrag; der MR-Vorgänger verlässt §Aktive Adaptionen |
| neue MR-Datei unter `harness/conventions/` | neu | der Hebung-Eintrag (gemessener Delta, Frozen-Liste, Zitat-Delta) |
| [`harness/conventions/done/MR-072-baseline-v690.md`](../../../../harness/conventions/done/MR-072-baseline-v690.md) | move | Auflösungs-Trigger („die nächste Pin-Hebung") tritt — reiner `git mv` nach `conventions/done/` |
| `AGENTS.md`, `harness/README.md`, `harness/rules/*.md`, aktive `harness/conventions/MR-*.md`, `.claude/agents/reviewer.md`, `roadmap.md`, `planning/README.md`, `observations/README.md`, `.d-check.closure.yml` | update | lebende Token-Swaps `v6.9.0`→`v6.13.0`, je Datei einzeln (kein pauschales sed über index-tragende Dateien — die Lehre aus [`MR-072`](../../../../harness/conventions.md#mr-072)) |
| `.claude/rules/*.md` (7 Symlinks) | update | Alias-Aliase nach [`MR-055`](../../../../harness/conventions.md#mr-055) auf `v6.13.0` umgehängt |
| `.d-check.yml` | update | nur, falls gemessen: nächste `ignore-refs`-Stufe der Bump-Kette (nach [`MR-069`](../../../../harness/conventions.md#mr-069)) |
| `slice-240`-Plan selbst | update | die beiden `d-check:cite`-Spannen neu geankert, sobald `modul-05` im neuen Baum liegt (nach [`MR-051`](../../../../harness/conventions.md#mr-051)) |

## 4. Trigger

**Start** (`next` → `in-progress`): direkt beansprucht — der Nachtlauf
(`make baseline-freshness`, Exit 3, 2026-09-29) liegt vor und die
Auftraggeber-Anfrage „auf das neueste Regelwerk umstellen" ist die
Beanspruchung selbst.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): Eine lebende Datei trägt den alten Pin
  gleichzeitig als lebenden Verweis und als Frozen-Vergangenheits-Aussage,
  sodass sich die Klassen nicht je Fundstelle trennen lassen (eine
  Massen-Ersetzung würde eine der beiden fälschen) — dann ist der Nachzug
  ein eigener, feinerer Slice.
- `in-progress` → `open` (Blocker): Der Vendor-Lauf scheitert an Netz/
  Integrität (`SHA256SUMS` passt nicht zum Asset) — dann Carveout-Pflicht
  statt stiller Altstand.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag (drei Paarungen wellenlos
hier geprüft).

## 6. Risiken und offene Punkte

- Eine lebende Referenz sitzt in einer Datei, die auch eine Frozen-Vergangenheits-Aussage
  trägt, und eine Ersetzung hebt sie mit (Über-Hebung, Lehre aus
  [`MR-070`](../../../../harness/conventions.md#mr-070)). —
  **Ausgang:** *(offen)*
- Die `d-check:cite`-Spannen verschieben sich durch den Re-Vendor
  (`citation-mismatch`, nach [`MR-051`](../../../../harness/conventions.md#mr-051)). —
  **Ausgang:** *(offen)*
- Frozen-Dateien tragen Markdown-Links auf den entfernten `v6.9.0`-Baum und
  werden `target-missing` (nach [`MR-069`](../../../../harness/conventions.md#mr-069)). —
  **Ausgang:** *(offen)*
- Eine Ersetzung trifft die Index-Tabellenzeile des MR-Vorgängers in
  `harness/conventions.md` (dieselbe Klasse wie der Fund aus
  [`MR-072`](../../../../harness/conventions.md#mr-072)). —
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

**Vorgelagert — offene Beobachtungen sichten:**
[`BEO-ALL/pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
ist die stehende Beobachtung zu genau diesem Vorgang — **7×** (slice-106,
slice-110, slice-117, slice-148, slice-189, slice-222, slice-224), dieses
Slice wird das 8. Auftreten. Keine weitere berührte Sub-Area mit Treffern;
der Ausgang bleibt ausgeschildert (siehe state.md: „kein formgültiger
Ausgang").

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung, 2026-09-29):
`make baseline-freshness` selbst ist der Anlass (Exit 3): vier neuere Tags
(`v6.10.0`–`v6.13.0`), Content am gepinnten Tag unverändert. `make
nightly-state` wird bei der Beanspruchung gelesen und der Stand in §7
notiert.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.

### Sub-Area: `*` (Harness-Werkzeug)

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — die Prozedur ist 14-fach verkörpert
  (Kette ab [`MR-011`](../../../../harness/conventions.md#mr-011) in
  `harness/conventions.md` samt Template `MR-<NNN>-titel.template.md`).
- **Phase-Reife:** Phase 3.
- **Evidenz-/Diskrepanz-Risiko:** Niedrig — reine Pin-Fortschreibung; die
  Diskrepanz-Frage (lebend vs. frozen) ist selbst der Gegenstand der
  Messung in §6/§7.
- **Reconciliation-Aufwand:** Keiner — Referenz-Nachzug ist Teil des Slice;
  Graduation-Trigger bleibt `BEO-ALL/pin-bump-mirrors-ungated`.
