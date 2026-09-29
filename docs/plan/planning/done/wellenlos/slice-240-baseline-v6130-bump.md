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

- [x] Vendierter Baum steht auf `v6.13.0` (`SHA256SUMS` generiert und
      verifiziert), der `v6.9.0`-Baum ist entfernt; §Baseline-Pin (Stand-URL
      + Datum) nachgezogen.
- [x] Alle deklarierten lebenden Referenz-Klassen nennen `v6.13.0`: Pfad-Verweise
      (`links`/`codepaths`-deckungsgleich), Release-/Tree-URLs, bare
      Versionsnennungen in lebenden Dokumenten; die 7 Alias-Symlinks unter
      `.claude/rules/` lösen gegen `v6.13.0` auf. *(gemessen 8 statt 7 —
      Abweichung in
      [`MR-073`](../../../../harness/conventions.md#mr-073--baseline-pin-hebung-auf-v6130-fünfzehnter-nachtrag-zu-mr-011-nachtrag-zu-mr-023)
      dokumentiert)*
- [x] `make gates` grün; `make baseline-freshness` meldet den Pin aktuell
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

- **Was hat funktioniert:** Messen vor Schreiben — Delta (28 von 55 Dateien
  mit Inhalt) und Spiegel-Klassen (75 Dateien / 161 Vorkommen) am echten
  Vorzustand, Frozen-Klassen per Eigenschaft statt per Verzeichnis
  ([`MR-070`](../../../../harness/conventions.md#mr-070)); die Commit-Zerlegung
  §3.3 (reiner Move, ±0 Zeilen) und der
  [`MR-069`](../../../../harness/conventions.md#mr-069--das-ignore-refs-ventil-ist-eine-deklarierte-gate-senkung-und-es-wächst-mit-jedem-bump)-Ventil-Nachzug
  (4 Einträge, messbegründet) hielten, `make gates` nach Implementierung und
  nach Review-Korrektur grün.
- **Was ging anders als geplant:** Drei Instanzen derselben Klasse
  record-claim-vs-diff am MR-Eintrag: der Vorgänger-Swap traf
  [`MR-072`](../../../../harness/conventions.md#mr-072--baseline-pin-hebung-auf-v690-vierzehnter-nachtrag-zu-mr-011-nachtrag-zu-mr-023)s
  Index-Zeile und -Datei (Plan-Risiko 4, revertet); die erste
  [`MR-073`](../../../../harness/conventions.md#mr-073--baseline-pin-hebung-auf-v6130-fünfzehnter-nachtrag-zu-mr-011-nachtrag-zu-mr-023)-Fassung
  behauptete „kein cite-Neu-Ankern" (Review-R1-F-1); und die Korrektur
  15/6/9 war selbst nicht diff-genau (Verifier V-1: 14/7/7, die
  [`MR-043`](../../../../harness/conventions.md#mr-043--der-werkzeug-einstieg-importiert-agentsmd-statt-auf-ihn-zu-verweisen)-Spanne
  fehlte). Jede Instanz klein, alle drei erst durch Zählen gegen den Diff
  gefangen.
- **Steering-Loop-Eintrag:** [`BEO-ALL/pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
  — das Signal (MR-Eintrag behauptet Abwesenheit, der Diff zählt Vorkommen)
  trat in diesem Slice zweimal auf (R1-F-1, V-1) nach dem Vorgänger-Fund
  slice-224; der Ausgang bleibt „kein formgültiger", der Eintrag bleibt
  stehen.
- **Beobachtungs-Register (`../observations/`):** kein neuer Eintrag — der
  stehende BEO-ALL/pin-bump-mirrors-ungated trägt den Vorgang.
- **Folge-Slices:** vier inhaltliche Deltas der Hebung sind ohne Urteil
  geblieben und warten je auf einen eigenen Konventions-Slice: Trigger-Audit
  der Welle (`modul-06`), „Gate-Erweiterung ist kein ADR-Anlass" (`modul-04`),
  AGENTS.template-Schnitt-Prinzip (Kurzzeile/Volltext-Trennung), ids/matrix-
  Linkpflicht für ADR-Kennungen.
- **Risiken aus §6:** R1 (Ersetzung trifft MR-Vorgänger-Zeile/-Datei) —
  eingetreten, revertet; R2 (cite-Verankerung verschiebt sich) — 7 statt 0
  neu geankert, Wortlaut identisch; R3 (frozen target-missing) — 4 Markdown-
  Links, Ventil gewachsen; R4 (Zitat-Delta) — nicht eingetreten, Zitatzeile
  byte-identisch.
- **Drei Paarungen:** Lerneintrag „record-claim-vs-diff am MR-Eintrag, zweite
  Instanz am selben Vorgang" — Folge-Slice: keiner formuliert, das Signal
  trägt der stehende BEO-Eintrag — Register: BEO-ALL/
  pin-bump-mirrors-ungated, unverändert offen.

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
