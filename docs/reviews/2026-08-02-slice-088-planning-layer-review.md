# Review — slice-088 Planning-Layer-Form (welle-67, Etappe D-1)

- **Review-Art:** unabhängiges, kontext-getrenntes Frischkontext-Review
  (adversarial). Gegenstand: die d-check-eigene **Planning-Schicht** an die
  v5.0.0-Baseline-Form gehoben — Wellendokument (D-2), Beobachtungs-Register
  (D-3), Roadmap auf die sechs Baseline-Abschnitte (D-1), Slice-Vorprüfungen
  (D-4), Carveout-/Trigger-Audit im Welle-Closure-Trigger (D-9), plus die
  deklarierte Ruhe-Marker-Adaption `MR-024`.
- **Skill:** `.harness/skills/reviewer.md` v1.2.0.
- **Modell:** claude-opus-4-8 (Opus 4.8).
- **Datum:** 2026-08-02.
- **Range:** `1c535e5..HEAD` — sieben Commits (`e81fc21`→`deadd21`).
- **Gates (selbst gefahren):**
  - Netzloser Selbst-Scan `docker run --rm --network none -v "$PWD":/repo:ro d-check:latest`
    grün: 314 Datei(en), 0 Befund(e) — deckt `links`/`anchors`/`ids`/`matrix`/
    `planning` ab, d. h. alle Roadmap-Zeiger lösen auf, alle Anker existieren,
    der Planning-Invariant hält.
  - `make gates` grün: doc-check + lint + test + arch-check + coverage-gate
    (94.20 % ≥ 93 %) + semgrep (0/0) + gate-consistency + planning-check.
  - `make adr-check RANGE=1c535e5..HEAD` grün: 314 Datei(en), 0 Befund(e)
    (Modul vcs); `git diff --name-only 1c535e5..HEAD | grep -i adr` ist leer —
    **keine** Accepted-ADR berührt.
- **Änderungs-Umfang (13 Dateien, +508/-29):** 7× `done/welle-6N-results.md`
  (Backfill), `roadmap.md`, `welle-67-baseline-v500-migration.md`,
  `observations.md`, `slice-088-…md` (Plan), `harness/conventions.md` (+1 Zeile),
  `harness/conventions/MR-024-…md`. Keine Produkt-, Spec- oder ADR-Datei.
- **Prüf-Achsen:** Template-Konformität (Welle · Register) · Roadmap-Struktur ·
  Ehrlichkeit des Backfills · `MR-024`-Korrektheit · Konsistenz mit `modul-06`/
  `modul-05` · Slice-Doc↔Umsetzung. Verifiziert mit grep, `git show`/`git diff`
  über den Range und den drei Gate-Läufen — nicht gegen Zusammenfassung.

## Findings

### F-1 (INFO) — Backfill zementiert Einzel-Slice-Läufe (welle-61…66) als „Wellen" im Closure-Log, ohne die Reklassifikation zu benennen

- **Kategorie:** INFO (dokumentationswürdige, aber undokumentierte Annahme).
- **Quelle:** Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
  braucht („Bei einem einzelnen Slice … läuft ohne Welle"; „Wellenlose Arbeit
  erscheint nicht in der Roadmap — weder beim Start noch beim Abschluss … das
  Closure-Log der Roadmap ist für Wellen").
- **Pfad:** `docs/plan/planning/done/welle-65-results.md:1` und
  `docs/plan/planning/done/welle-66-results.md:1` (sowie `welle-61…64` analog);
  `docs/plan/planning/in-progress/roadmap.md:88-94` (Closure-Log-Zeilen).
- **Befund:** welle-61…66 tragen je genau **einen** Slice (078 / 079 / 080 /
  081 / 072 / 082); welle-65 und welle-66 sind reine Einzel-Slice-Doku **ohne
  Release**, deren Closure-Trigger die eigene Slice-DoD nur abgeschrieben hätte —
  nach dem jetzt adoptierten `modul-06`-Kriterium wellenlose Arbeit, die nicht ins
  Closure-Log gehört. Die sieben Backfill-Notizen markieren sich sauber als
  „retroaktiv/minimal", benennen diese Klassifikations-Spannung aber nicht; sie
  reproduzieren die historische „Welle NN"-Rahmung unkommentiert.
- **Failure-Szenario:** Ein späterer Planner liest das Closure-Log als Präzedenz,
  eröffnet für einen einzelnen reinen-Doku-Slice eine „Welle" plus
  `welle-NN-results.md` und widerspricht damit `modul-06` §Wann Arbeit eine Welle
  braucht — die Rahmung sagt ihm nicht, dass 65/66 gerade der Grenzfall waren, den
  die Regel als Nicht-Welle führt.
- **Nicht blockierend, weil:** die Rahmung historisch ist (die Drift-Tabelle
  eröffnete/schloss diese Wellen bereits als Wellen), der Backfill vom
  Auftraggeber ausdrücklich als **minimale** Log-Auflösung beschlossen ist
  (`welle-67-…md` §6 Out-of-Scope) und die Notizen ehrlich als retroaktiv
  gekennzeichnet sind — kein Fabrikat, nur eine unbenannte Reklassifikation.
- **Verifizierbar:** nein (Klassifikations-/Ehrlichkeitsurteil; kein Gate prüft
  „ist dies wirklich eine Welle"). grep über die Closure-Log-Slice-Zuordnung
  belegt die Einzel-Slice-Natur.

## Negativbefunde (geprüft, ohne Befund)

- **Wellendokument-Template-Konformität** (`welle-67-baseline-v500-migration.md`
  gegen `welle.template.md`): alle Pflichtteile vorhanden — §1 Welle-Ziel · §2
  Trigger · §3 Closure-Trigger · §4 Slices (Spalten `Slice | Titel | Bezug`,
  **keine** Status-Spalte) · §5 Abhängigkeiten · §6 Out-of-Scope · §7
  Closure-Notiz. **Keine** `**Status:**`-Kopfzeile (Zustand = Verzeichnis).
  §7 korrekt **ausstehend** — die `done/`-Pointer sind noch **nicht** geschrieben,
  die Ruheort-Regel wird nicht verletzt (kein `target-missing` von flach). Der
  Start-Trigger (§2, `slice-083` abgenommen) steht **nicht** in der Slice-Liste
  §4 (084–091) — kein zirkulärer Trigger.
- **Register-Form** (`observations.md` gegen `observations.template.md`): Kopf
  (`**Status:** Aktiv`), Spaltensatz `Kennung | Beobachtung | Sub-Area | Zähler |
  Belege | Stand`, Sektion „Gestrichene Einträge" mit eigener Spaltenzeile, beide
  Tabellen mit `— keine —`. Form-korrekt.
- **Roadmap-Struktur** (`roadmap.md`): sechs `##`-Abschnitte in Template-Reihenfolge
  (Aktuelle Welle · Nächste Wellen · Meilensteine · Abhängigkeitsgraph ·
  Abgeschlossene Wellen · Historische Trigger-Verschiebungen — deckt sich mit
  `roadmap.template.md`; `modul-06`-Prosa nennt „fünf", das Template und der
  Struktur-Bullet fügen `Abhängigkeitsgraph` explizit hinzu, daher sechs `##` —
  keine Divergenz). §Aktuelle Welle nennt die drei Pflicht-Bestandteile
  (Slice-IDs 088–091 · Start-Trigger · Closure-Kriterien).
- **Vorgänger-Prosa-Migration:** `git show b84b2cd` belegt, dass die entfernte
  §Aktuelle-Welle-„Vorgänger"-Prosa (welle-60…66) **vollständig** in
  §Abgeschlossene Wellen überführt ist — sieben Zeilen, nichts verloren, keine
  Dublette (die alte Prosa ist restlos entfernt).
- **Closure-Log-Zeiger:** alle sieben `[welle-NN-results.md](../done/welle-NN-results.md)`
  lösen auf (Dateien existieren, `links`/`anchors` grün); die Ziel-Notizen sind
  frei von erfundenen Belegen — keine behaupteten Digests/Run-IDs, Belege an
  „Slice-`done/`-Dateien + git-Historie" delegiert, die Versionsangaben
  (v0.49.0…v0.51.1) decken sich mit der bestehenden Drift-Tabelle.
- **planning-Guard:** `make planning-check` grün — der Roadmap-Umbau und die neuen
  flachen `welle-*.md`/`observations.md` brechen weder den `## Aktuelle Welle`-
  Heading-Invariant noch die in-progress-↔-Ruhe-Marker-Invariante (aktuell
  `slice-088` in `in-progress/` ⟺ Welle benannt, kein Ruhe-Marker — konsistent).
- **`MR-024`-Korrektheit:** benennt **genau eine** Ersetzt-Baseline-Regel
  (`modul-06` §Roadmap-Struktur (Aktuelle Welle), Anker
  `#roadmap-struktur-fünf-abschnitte-modul-6` — von `anchors` als existent
  bestätigt). Adaption ist Prosa + Verweis aufs Wellendokument (kein Feld-Duplikat)
  plus gate-erzwungener Ruhe-Marker; berührt **keine** `Accepted`-ADR (referenziert
  `ADR-0028`/`DC-FA-PLAN-001` nur als Geltungsbereich). Index-Zeile in
  `conventions.md:100` formkonform (Link + Voll-Slug-`<a id>` + Titel +
  Geltungsbereich + Ersetzt-Baseline-Regel), kein `#NAME?`, Anker deckt sich mit
  dem Datei-Titel.
- **D-4-Vorprüfungen:** slice-088 §7 „Vorgelagert" trägt die zwei `modul-05`-Schritte
  (Sub-Area prüfen · offene Beobachtungen sichten) **vor** der §8-Modus-Begründung
  — Reihenfolge korrekt; „keine offene Beobachtung" ist als Antwort notiert. Die
  Reichweite (Modell hier, AGENTS-Konvention → slice-089) ist in `deadd21` sauber
  abgegrenzt.
- **D-9-Trigger-Audit:** `welle-67-…md` §3 nennt den Carveout-/bootstrap-aware-Gate-/
  ADR-Re-Eval-Audit (`modul-06` Closure-Schritt 2) mit „0 aktive Carveouts, latent" —
  benannt, nicht still.
- **Konsistenz Meilensteine/Abhängigkeitsgraph:** „— keine offenen —" und „d-checks
  Wellen sind weitgehend unabhängig / keine Phantom-Welle" sind ehrliche Aussagen
  über d-checks kontinuierliches-Release-/kapazitätssequenziertes Betriebsmodell,
  keine Fehlbehauptung und keine Leerzeremonie.
- **Slice-Doc↔Umsetzung:** Ziel/Vorgehen/DoD/Abnahme-Punkte des Slice-Plans decken
  die tatsächliche Umsetzung; welle-67 selbst ist eine **echte** Welle (084–091 mit
  repo-weitem Mehr `make gates` + `make adr-check`), der Abnahme-Punkt 1
  (Ruhe-Marker ↔ Template) ist über `MR-024` erfüllt.
- **Immutabilität/Scope:** `adr-check` grün über den Range, keine ADR/Spec/Code-Datei
  im Diff — reine Planning-/Konventions-Doku wie deklariert.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

## Verdikt

**Abnahmereif.** Die Planning-Schicht ist form-konform an die v5.0.0-Baseline
gehoben: Wellendokument und Register tragen die Pflichtteile ohne verbotene
Status-Felder, die Roadmap führt die sechs Abschnitte in korrekter Reihenfolge mit
auflösenden Closure-Zeigern, der Backfill ist ehrlich als retroaktiv/minimal
markiert und frei von erfundenen Belegen, `MR-024` deklariert die Ruhe-Marker-Adaption
regelkonform ohne ADR-Berührung. Alle drei Gate-Läufe grün. Der einzige Befund (F-1)
ist INFO und **nicht blockierend** — die unbenannte Wellen-/Nicht-Wellen-Reklassifikation
der historischen Einzel-Slice-Läufe 61…66 ist eine dokumentationswürdige, aber
bewusst und ausdrücklich vom Auftraggeber scope-begrenzte Annahme, kein Fabrikat.
