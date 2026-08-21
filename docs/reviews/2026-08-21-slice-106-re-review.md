# Re-Review — slice-106: Baseline-Bump v5.0.0 → v5.6.0, Etappe A (Heilung der Erst-Review-Befunde)

- **Review-Art:** bestätigende Re-Review (frischer Kontext) — geprüft wird
  ausschließlich, ob die Heilung die sieben Befunde des Erst-Reviews
  ([`2026-08-21-slice-106-baseline-v560-review.md`](2026-08-21-slice-106-baseline-v560-review.md))
  schließt und ob sie neue Defekte einführt.
- **Gegenstand:** slice-106; die neu geschnittene, unpushed Commit-Kette
  `8ba8490` (Skript-Fix) → `d20accc` (Vendoring + Baum-Entfernung + Retargets +
  Tombstones) → `d4e1744` (Erst-Review-Report) → `77b9c15` (Auflagen
  F-1/F-3/F-4/F-5/F-6/F-7). `origin/main` = `0b9431f`; HEAD = `77b9c15`,
  Arbeitsbaum clean.
- **Skill:** `reviewer.md` @ 1.4.0 (2026-08-15).
- **Modell-ID:** claude-fable-5.
- **Datum:** 2026-08-21.
- **Eingangs-Kontext:** Erst-Review-Report (Findings F-1–F-7 samt Belegen),
  MR-021/023/025/026, DC-FA-REF-001 (`ignore-refs`-Ventil), `.d-check.yml`,
  `tools/harness/fetch-baseline-cache.sh`, `harness/README.md`,
  `harness/conventions.md`.
- **Selbst gefahrene Läufe:** `git show --stat` auf allen vier Commits;
  doc-check per `docker run --rm --network none … d-check:latest` gegen
  `git worktree`-Stände **aller vier** Kettenglieder: `8ba8490` → 395/0,
  `d20accc` → 396/0, `d4e1744` → 397/0, `77b9c15` (== Arbeitsbaum) → 397/0 —
  jeweils Exit 0; Worktrees danach entfernt. Grep-Zensus der
  `baseline/v5.0.0`-Fundstellen (Pfad-, URL- und Link-Achse getrennt).

---

## Befund-Schließung (je Erst-Review-Finding)

### F-2 · MEDIUM · Commit-Schnitt — GESCHLOSSEN

- `8ba8490` enthält **nur** `tools/harness/fetch-baseline-cache.sh`
  (1 Datei, +14/−7); keine Baum-Löschung mehr. `d20accc` trägt Vendoring
  (v5.6.0-Bäume + `SHA256SUMS`), die v5.0.0-Entfernung, alle Retargets
  (`AGENTS.md`, `harness/README.md`, `harness/conventions.md`, MR-Dateien,
  Planning-Docs, `reviewer.md`) und die `.d-check.yml`-Tombstones **atomar**
  (73 Dateien) und **deklariert** die Entfernung in der Message („UND der
  v5.0.0-Baum in DIESEM Commit entfernt … Entfernung und Verweis-Hebung
  atomar").
- Beide Messages sind ehrlich inkl. Ersetzungs-Verweis: `8ba8490` nennt
  „Ersetzt den gleichnamigen unpushed Commit 960600d" samt F-2-Begründung;
  `d20accc` nennt „Ersetzt den unpushed Commit 8a0477a (Review-Befund F-2)".
- Worktree-Gegenprobe: der Zwischenstand `8ba8490` ist doc-check-**grün**
  (395 Dateien, 0 Befunde) — der im Erst-Review rote Zwischenstand existiert
  in der neuen Historie nicht mehr; die gesamte Kette ist bisektierbar
  (alle vier Stände grün, siehe Läufe oben).

### F-1 · MEDIUM · Download-URL-Spiegel — GESCHLOSSEN

- `harness/README.md:42`: die `lab-regelwerk.zip`-URL lautet jetzt
  `releases/download/v5.6.0/lab-regelwerk.zip`; die beiden anderen
  URL-Spiegel (`AGENTS.md:32`, `harness/conventions.md:48`) stehen ebenfalls
  auf v5.6.0. `grep -rn "releases/download/v5"` über `harness/`, `AGENTS.md`,
  `tools/`, `.harness/skills/`: **kein** v5.0.0-Treffer mehr.
- Sweep nach weiteren nicht gehobenen Pin-Spiegeln: alle verbleibenden
  v5.0.0-Nennungen auf lebenden Oberflächen sind korrekte Historie bzw.
  Prozess-Prosa (MR-023-Körper als historische Hebung, `conventions.md:34/36`
  Ketten-Historie, aufgelöste MR-Zeilen 124–127, Roadmap-Chronik,
  Skript-Kopf-Kommentare, der slice-106-Plan als Auftragsbeschreibung,
  `.d-check.yml`-Tombstone-Pfade) — kein Spiegel behauptet v5.0.0 als
  aktuellen Stand.

### F-3 · LOW · MR-021-Beispiel — GESCHLOSSEN

- `harness/conventions/MR-021-vendored-verweise-pin-gebunden.md:17` sagt
  jetzt „tragen den **konkreten** Pin (aktuell `…/v5.6.0/…`)"; die
  `Ersetzt-Baseline-Regel`-Zeile (Z. 4) zeigt konsistent auf v5.6.0. Das
  „aktuell" macht die Zelle zudem robuster gegen die nächste Hebung.

### F-4 · LOW · Tombstone-Breite — GESCHLOSSEN

- `.d-check.yml:76–93`: das 19-Dateien-Glob ist durch **drei exakte**
  `in:`-Report-Pfade ersetzt; zusammen mit ADR-0047/0048 fünf Einträge, alle
  mit `refs: [".harness/baseline/v5.0.0/**"]`. Exaktheits-Zensus
  (`grep -rEln '\]\([^) ]*baseline/v5\.0\.0'` über `docs/ harness/ spec/ …`):
  genau die fünf getombstoneten Dateien tragen Markdown-**Links** auf den
  entfernten Baum — die zwei zusätzlichen grep-Treffer sind
  Falsch-Positive der Suchmaske (`roadmap.md:98` matcht über
  Link-plus-Inline-Code hinweg; `2026-08-21-…-review.md:96` ist das
  grep-Kommando-Beispiel in Inline-Code, kein Link). Keine lebende Datei
  wird mit-entschuldigt (2 immutable Accepted-ADRs + 3 eingefrorene
  Reports); kein gebrochener Link übrig (doc-check 397/0). Der
  Config-Kommentar erklärt die Link-vs-Inline-Code-Abgrenzung jetzt selbst.

### F-5 · LOW · MR-023-Indexzeile — GESCHLOSSEN

- `harness/conventions.md:101`: die Regel-Spalte lautet jetzt
  „self-contained Bundle-Layout (vendored; historische Pin-Hebung auf
  v5.0.0 — den aktuellen Pin trägt der Nachtrag darunter)" — die
  Index-Oberfläche widerspricht MR-026 nicht mehr und verweist auf die
  Auflösung, ohne den Körper lesen zu müssen.

### F-6 · INFO · Nachtrag-Zählung — GESCHLOSSEN

- `harness/conventions/MR-026-baseline-v560.md:22–23`: Klammer-Klarstellung
  „(Zählung der Titel-Serie; die v5.0.0-Hebung lief unter dem eigenen Titel
  des Layout-Eintrags außerhalb dieser Serie)" — die Mehrdeutigkeit ist für
  die nächste Hebung dokumentiert aufgelöst.

### F-7 · INFO · leerer SKIP-Satz — GESCHLOSSEN

- `tools/harness/fetch-baseline-cache.sh:169`: der leere-Extraktion-Zweig
  setzt jetzt `note="Upstream-Extraktion lieferte keine Dateien
  (Bundle-Layout?)"`; damit trägt jeder der vier SKIP-(Content)-Pfade
  (fehlendes Host-Werkzeug, fehlendes Manifest, leere Extraktion,
  nicht ladbares Asset) eine nicht-leere Ursachen-Angabe. Der
  `xargs -r`-Kommentar (Z. 162–163) dokumentiert den Fix-Grund.

---

## Neue Defekte durch die Heilung (geprüft, ohne Befund)

1. **Anker-Stabilität:** der explizite `<a id>` der geänderten
   MR-023-Indexzeile ist byte-identisch zum `origin/main`-Stand
   (`mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout`);
   7 Dateien verlinken darauf, doc-check grün. Der MR-026-Anker in der
   Indexzeile entspricht weiter dem unveränderten H1 (die F-6-Klarstellung
   sitzt im Körper, nicht im Titel).
2. **Tombstone-Verengung ohne Unter-Deckung:** der Übergang
   Glob → drei exakte Pfade (`d4e1744` → `77b9c15`) lässt doc-check grün
   (397/0 beiderseits); die per Glob mit-gematchten, jetzt ungedeckten
   Reports (u. a. `2026-08-03-slice-092…`) nennen den Pfad nur in
   Inline-Code und brauchen das Ventil nicht — vom Zensus bestätigt.
3. **Skript-Verhalten unverändert außer der Note:** der 77b9c15-Diff am
   Skript ist reine Note-Setzung + Kommentar; `authenticity`-Zustände und
   Exit-Codes unberührt.
4. **Kein Konformitäts-Überclaim in den neuen Messages:** beide
   Feat-Messages und die Auflagen-Message behaupten keine
   v5.6.0-Konformität; die Etappe-B/C-Grenze bleibt eingehalten.
5. **Historien-Hygiene:** `origin/main` = `0b9431f` unberührt; nur die vier
   unpushed Commits wurden neu geschnitten — kein published-History-Rewrite.
6. **Gates:** doc-check auf HEAD selbst gefahren (397/0, Exit 0), deckungs-
   gleich mit der `77b9c15`-Message („make gates gruen (397/0)"). Der volle
   `make gates`-Lauf und die DoD-Abhakung bleiben Verifikations-Rolle.

## Kategorie-Summary

| Kategorie | Erst-Review | davon geschlossen | neue Findings |
|---|---|---|---|
| HIGH | 0 | — | 0 |
| MEDIUM | 2 (F-1, F-2) | 2 | 0 |
| LOW | 3 (F-3, F-4, F-5) | 3 | 0 |
| INFO | 2 (F-6, F-7) | 2 | 0 |

## Verdikt

**APPROVE (bestätigend).** Alle sieben Erst-Review-Befunde sind geschlossen;
die Heilung führt keine neuen Defekte ein. Der neu geschnittene Commit-Bogen
ist ehrlich beschriftet (inkl. Ersetzungs-Verweise auf `960600d`/`8a0477a`)
und an jedem Zwischenstand doc-check-grün — die Blockade des Erst-Reviews
ist aufgehoben. Push/Closure aus Review-Sicht frei.
