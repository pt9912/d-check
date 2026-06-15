# Slice slice-022: Inline-HTML-Anker im Modul `anchors`

**Status:** open (Doc-Schicht erstellt; Implementierung nach Review).

**Welle:** welle-12-html-anker (Trigger: Change Request des Auftraggebers
— Falsch-Befunde `anchor-missing` auf manuell gesetzte HTML-Anker in der
Doku; GitHub rendert sie als Sprungziele, d-check kannte sie bisher nicht).

**Bezug:** [`DC-FA-ANCH-001`](../../../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors)
(Schärfung 0.12.0 — Inline-HTML-Anker zählen zur gültigen Anker-Menge),
Spezifikation §[`DC-FA-ANCH-001.b`](../../../../spec/spezifikation.md#dc-fa-anch-001b--inline-html-anker);
mittelbar [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
(geteiltes Anker-Verfahren), Determinismus-Vertrag
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus).

**Autor:** pt9912. **Datum:** 2026-06-15.

---

## 1. Ziel

Das Modul `anchors` erkennt heute nur ATX-Heading-Slugs als gültige
Anker. Manuell gesetzte HTML-Sprungziele — `<a name="…">` (klassisch)
und `id="…"` an einem beliebigen Element (modern) — rendert GitHub als
Anker; ein Link darauf meldet d-check jedoch als `anchor-missing`
(Falsch-Befund). Dieser Slice erweitert die **gültige Anker-Menge** um
Inline-HTML-Anker (GitHub-Parität), wörtlicher Vergleich. Da `codepaths`
dasselbe Anker-Verfahren nutzt, gilt die Erweiterung dort konsistent mit.

## 2. Definition of Done

- [x] **Lastenheft-Schärfung**
  [`DC-FA-ANCH-001`](../../../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors)
  (0.12.0): gültige Anker-Menge um Inline-HTML-Anker erweitert,
  Akzeptanzkriterien-Trio (HTML), Out-of-Scope neu gefasst, §5-Abgrenzung.
- [x] **Spezifikation**
  §[`DC-FA-ANCH-001.b`](../../../../spec/spezifikation.md#dc-fa-anch-001b--inline-html-anker):
  Extraktion (`id` beliebig, `name` an `<a>`), wörtlicher Vergleich,
  konservativ/zeilenbasiert; die `codepaths`-Anker-Prüfung auf die
  gemeinsame Anker-Menge fortgeschrieben.
- [ ] **Implementierung** in `internal/hexagon/core/anchors.go`:
  Inline-HTML-Anker-Extraktion (fence-bewusst, wie die Heading-Extraktion);
  Anker-Menge als Union aus Heading-Slugs (geslugged) und HTML-Ankern
  (wörtlich); bestehender Datei-Cache (`slugsFor`) unverändert nutzbar.
- [ ] **Tests** in `internal/hexagon/core/anchors_test.go`: die drei neuen
  HTML-Akzeptanzkriterien (Happy/Boundary/Negative) plus Fence-Ausschluss
  und Case-Sensitivität; bestehende Heading-Tests bleiben grün.
- [ ] `make gates` grün; bei nutzersichtbarer Wirkung
  [`CHANGELOG.md`](../../../../CHANGELOG.md) gepflegt.

## 3. Betroffene Artefakte

- **Slice-ID:** slice-022.
- **IDs:**
  [`DC-FA-ANCH-001`](../../../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors)
  (Schärfung),
  §[`DC-FA-ANCH-001.b`](../../../../spec/spezifikation.md#dc-fa-anch-001b--inline-html-anker)
  (neu),
  [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
  (mittelbar),
  [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
  (bleibt erfüllt).
- **ADR:** keiner. Kein neuer Port/Adapter — der Hexagon-Schnitt
  ([ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md)) bleibt
  unverändert; keine neue Dependency (konservative Regex-Erkennung wie
  bei `codepaths`/`spans`/`hostpaths`, kein HTML-Parser).
- **Module:** `anchors` (+ `codepaths` mittelbar).
- **Gates:** `make gates` (doc-check, test, coverage-gate, lint,
  arch-check, gate-consistency).

## 4. Vorgehen und Risiken

- **Erkennung:** zeilenbasierter, fence-bewusster Scan analog
  `ExtractHeadings`; Attribut-Werte `id="…"`/`id='…'` (beliebiges Element)
  und `name="…"`/`name='…'` (nur an `<a>`). Über Zeilengrenzen verteilte
  Tags bleiben unerkannt (konservativ).
- **Anker-Menge:** Union aus geslugten Heading-Slugs und wörtlichen
  HTML-Ankern; Fragment-Treffer = Mitgliedschaft in der Union.
- **Risiko `id=` in Prosa:** ein irrtümlich erkannter Anker vergrößert
  nur die Menge und führt zum Schweigen des Moduls (nie ein Falsch-Befund)
  — bewusst akzeptiert (Spezifikation
  §[`DC-FA-ANCH-001.b`](../../../../spec/spezifikation.md#dc-fa-anch-001b--inline-html-anker)).
- **Determinismus:** Mengenbildung ist reihenfolgeunabhängig
  (`map[string]bool`); identische Eingabe ⇒ identische Ausgabe
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- **Dogfooding:** nach der Implementierung läuft `make doc-check` über das
  eigene Repo; etwaige bisher unsichtbare HTML-Anker werden dann mitgeführt.

## 5. Lifecycle

`open/` → `next/` → `in-progress/` → `done/` per reinem `git mv`
(`AGENTS.md` §3.3). Closure-Notiz (§7) folgt in `done/` mit dem
Commit-Hash der Umsetzung.
