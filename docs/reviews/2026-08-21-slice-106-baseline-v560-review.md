# Review — slice-106: Baseline-Bump v5.0.0 → v5.6.0, Etappe A (Vendoring, Pin, Verweis-Hebung)

- **Review-Art:** unabhängiger Erst-Review (Code + Doku, frischer Kontext) — geprüft
  gegen Slice-Plan, Wellen-Plan, MR-Kette und die Hard Rules; **vor** der Slice-Closure.
- **Gegenstand:** slice-106; Commits `960600d` (fix: `fetch-baseline-cache.sh`
  check-latest) und `8a0477a` (Vendoring v5.6.0, Pin, MR-026, Verweis-Retargets,
  Tombstones). Beide zum Review-Zeitpunkt **nicht gepusht** (`origin/main` = `0b9431f`).
- **Skill:** `reviewer.md` @ 1.4.0 (2026-08-15).
- **Modell-ID:** claude-fable-5.
- **Datum:** 2026-08-21.
- **Eingangs-Kontext:** `docs/plan/planning/in-progress/slice-106-baseline-v560-vendoring.md`
  (inkl. DoD und Abnahme-Punkte 1–3), `docs/plan/planning/welle-78-baseline-v560-migration.md`,
  MR-011/012/016/021/023/025/026, DC-FA-REF-001 (`ignore-refs`-Ventil), `.d-check.yml`,
  `tools/harness/fetch-baseline-cache.sh`. Kein `DC-*`-Produkt-Bezug (reine Harness-Arbeit).
- **Selbst gefahrene Läufe:** `bash tools/harness/fetch-baseline-cache.sh --verify`
  (grün, „51 Dateien, vollständig"); `make doc-check` (grün, „396 Datei(en) geprüft,
  0 Befund(e)"); Gegenprobe per `git worktree` auf dem Zwischenstand `960600d`
  (doc-check **rot**, siehe F-2). `--check-latest` wurde **nicht** gefahren (Netz);
  die Commit-Behauptung „Currency OK + Content OK" bleibt hier unverifiziert
  (informativ, kein Gate — Verifikations-Rolle).

---

## Findings

### F-1 · MEDIUM · Zip-Download-URL in `harness/README.md` nennt weiter v5.0.0

- **kategorie:** MEDIUM
- **quelle:** MR-025 (Spiegel-Liste; der Slice-Plan §7 nennt `harness/README.md`
  ausdrücklich als Spiegel der Pin-Hebung), MR-021 (Pin-Bindung)
- **pfad:** `harness/README.md:42`
- **befund:** Die Guides-Tabellenzeile wurde nur halb gehoben: der Verzeichnis-Link
  zeigt auf `.harness/baseline/v5.6.0/regelwerk/`, aber die `lab-regelwerk.zip`-URL
  in derselben Zelle lautet weiter `releases/download/v5.0.0/lab-regelwerk.zip`.
  Die beiden anderen Spiegel derselben Aussage (`AGENTS.md` §1,
  `harness/conventions.md` §Adoptierte Konventions-Quellen) wurden gehoben; die
  Commit-Message von `8a0477a` behauptet „Alle lebenden Verweise pin-gebunden
  retargetet". Wer der URL folgt, lädt das v5.0.0-Bundle als vermeintlich
  aktuellen Stand — ein Vergleich gegen das vendored v5.6.0-`SHA256SUMS` meldet
  dann scheinbaren Drift.
- **verifizierbar:** nein durch Gates (externe URLs sind nicht gate-gedeckt —
  bekannter Blind-Spot der Klasse „Release-Prep-Doku-Currency"); Beleg per
  `grep -rn "releases/download/v5.0.0" harness/` (genau ein Treffer).
- **klasse:** pin-spiegel-ausserhalb-gate

### F-2 · MEDIUM · Commit-Schnitt: v5.0.0-Baum-Entfernung fährt undeklariert im Skript-Fix-Commit; Zwischenstand ist gate-rot

- **kategorie:** MEDIUM
- **quelle:** Maintainability (Commit-Grenzen-Disziplin), MR-021 (Bump-Prozedur:
  Entfernen (1) und Retargeten (2) bilden **einen** gate-ehrlichen Bogen)
- **pfad:** Commit `960600d` (52 Dateien, −6 283 Zeilen) vs. Commit-Message `8a0477a`
- **befund:** `960600d` ist als Skript-Fix deklariert (Message nennt ausschließlich
  die zwei check-latest-Defekte), enthält aber zusätzlich die komplette Entfernung
  von `.harness/baseline/v5.0.0/` (51 Dateien); `8a0477a` reklamiert diese
  Entfernung („der v5.0.0-Baum entfernt") für sich, enthält sie aber nicht. Der
  Zwischenstand `960600d` ist dadurch inkonsistent: Pin und sämtliche Verweise
  zeigen auf den bereits entfernten Baum — ein doc-check-Lauf gegen einen
  Worktree auf `960600d` meldet reihenweise `target-missing` (u. a.
  `harness/README.md:42`, `harness/conventions.md:47`, alle aktiven `MR-*`-Dateien).
  Die Historie ist an dieser Stelle weder bisektierbar noch trägt sie ehrliche
  Commit-Aussagen.
- **verifizierbar:** ja — `git worktree add <dir> 960600d` + doc-check-Lauf gegen
  den Worktree (im Review ausgeführt: rot); `git show 960600d --stat` zeigt die
  undeklarierten Deletions.
- **klasse:** commit-boundary-content-mismatch

### F-3 · LOW · Aktives MR-021 trägt den abgelösten Pin als „konkreten Pin"-Beispiel

- **kategorie:** LOW
- **quelle:** MR-021, MR-025
- **pfad:** `harness/conventions/MR-021-vendored-verweise-pin-gebunden.md:17`
- **befund:** Der Adaptions-Text der **aktiven** MR erklärt „Diese Links tragen den
  **konkreten** Pin (`…/v5.0.0/…`)" — nach der Hebung ist das genannte Beispiel
  nicht mehr der konkrete Pin. Die `Ersetzt-Baseline-Regel`-Zeile derselben Datei
  wurde retargetet, der Beispiel-Pin im Körper nicht (Inline-Code, von keinem
  Gate gedeckt).
- **verifizierbar:** nein durch Gates; Beleg per grep.
- **klasse:** pin-spiegel-ausserhalb-gate

### F-4 · LOW · Review-Tombstone-Glob deckt 19 Reports, belegt sind 3

- **kategorie:** LOW
- **quelle:** DC-FA-REF-001 (quell-skopiertes Ventil), Slice-Plan §2.4
  („quell-skopiert ausnehmen")
- **pfad:** `.d-check.yml:86`
- **befund:** Das Glob `docs/reviews/2026-08-0[239]-*.md` matcht 19 Report-Dateien;
  tatsächliche v5.0.0-**Links** (das Einzige, was links/anchors dort prüfen —
  ids/codepaths sind für `docs/reviews/**` ohnehin exempt) tragen nur drei:
  `2026-08-02-slice-091-status-feld-review.md`, `2026-08-09-backlog-schnitt-review.md`,
  `2026-08-09-slice-096-vertrag-review.md`. Der Datums-Zweig `[3]`
  (`2026-08-03-slice-092…`) hat ausschließlich Inline-Code-Nennungen und braucht
  das Ventil gar nicht. Der Config-Kommentar begründet mit „drei Review-Reports",
  das Ventil ist breiter als der deklarierte Bedarf; die MR-023-Präzedenz-Tombstones
  (v1.4.0) sind dagegen exakte `in:`-Pfade. Mildernd: alle mit-gematchten Dateien
  sind selbst eingefrorene Reports — **kein lebendes Dokument** wird mit-entschuldigt.
- **verifizierbar:** ja — `grep -l '](…baseline/v5.0.0' docs/reviews/` liefert die
  drei Dateien; ein auf drei exakte `in:`-Einträge verengtes Ventil ließe
  `make doc-check` weiterhin grün (Gegenprobe möglich).
- **klasse:** ventil-breiter-als-bedarf

### F-5 · LOW · Index-Zeile MR-023 behauptet in der aktiven Tabelle weiter „Baseline-Pin `v5.0.0`"

- **kategorie:** LOW
- **quelle:** MR-026 (Abgrenzung „MR-023 trägt das Layout, nicht den Pin-Wert"),
  Konventionsspeicher-Regel („was hier steht, liest jeder Agentenlauf")
- **pfad:** `harness/conventions.md:101`
- **befund:** Die Regel-Spalte der **aktiven** Adaptions-Tabelle führt für MR-023
  weiterhin „Baseline-Pin `v5.0.0` + self-contained Bundle". Nach der Hebung stehen
  damit zwei aktive Zeilen mit unterschiedlichen Pin-Aussagen nebeneinander
  (MR-023: v5.0.0, MR-026: v5.6.0); die Auflösung („Layout ja, Pin-Wert nein")
  steht nur im MR-026-Körper (§Auflösungs-Trigger), nicht an der Index-Oberfläche,
  die der Schnell-Leser konsultiert.
- **verifizierbar:** nein durch Gates (Prosa-Zelle); Beleg per Lesen der Tabelle.
- **klasse:** index-oberflaeche-vs-koerper-drift

### F-6 · INFO · Nachtrag-Zählung „dritter Nachtrag zu MR-011" ist mehrdeutig

- **kategorie:** INFO
- **quelle:** MR-011/012/016/023/026 (Ketten-Konsistenz)
- **pfad:** `harness/conventions/MR-026-baseline-v560.md:1`
- **befund:** MR-026 zählt sich als „dritter Nachtrag zu MR-011" (nach MR-012,
  MR-016). MR-023 beschreibt sich im eigenen Körper aber ebenfalls als „der von
  MR-011 vorgesehene Nachtrag" — nach dieser Lesart wäre MR-026 der vierte. Die
  Zählung ist nur als Titel-Serien-Zählung konsistent (allein MR-012/MR-016 tragen
  „Nachtrag zu MR-011" im Titel; MR-023 nicht, dafür sagt MR-026 zusätzlich
  „Nachtrag zu MR-023"). Dokumentationswürdige Mehrdeutigkeit für die nächste
  Pin-Hebung, kein Fehler mit eigenem Versagens-Szenario.
- **verifizierbar:** nein (Prosa); Beleg: Titel/Körper der vier MRs.
- **klasse:** ketten-zaehlung-ambig

### F-7 · INFO · check-latest: der neue leere-Listen-SKIP druckt eine leere Begründung

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** `tools/harness/fetch-baseline-cache.sh:164-166,188`
- **befund:** Im durch `xargs -r` neu erreichbaren Pfad (leere Datei-Liste ⇒ `up`
  leer ⇒ `authenticity` bleibt `skip`) ist `note` leer; die SKIP-Meldung lautet
  dann „check-latest SKIP (Content) — . Pin: …" ohne Ursachen-Angabe. Der
  Operator kann diesen SKIP nicht von anderen Ausfalls-Ursachen unterscheiden.
  Verhaltens-Richtung des Fixes ist korrekt (SKIP statt Phantom-Drift).
- **verifizierbar:** ja (nur synthetisch: leerer `vendored/`-Baum); im Normallauf
  unerreichbar, da eine erfolgreiche Extraktion ≥ 1 Regelwerks-Datei garantiert.
- **klasse:** diagnose-stummer-skip

---

## Negativbefunde (geprüft, ohne Befund)

1. **Verweis-Hebung (Pfad-Achse):** `grep -rn "baseline/v5.0.0"` über den
   Arbeitsbaum — außer Tombstone-Config, eingefrorenen ADRs 0047/0048,
   `done/`-Slices (nur Inline-Code), Review-Reports, Roadmap-Drift-Text und der
   Migrations-Doku (welle-78, slice-106, MR-026) nennt **kein lebendes Dokument**
   mehr einen `baseline/v5.0.0`-Pfad. Die verbleibenden v5.0.0-Nennungen in
   `harness/conventions.md:34/36` und den aufgelösten MR-Zeilen (124–127) sind
   korrekte Ketten-/Auflösungs-Historie. (Rest auf der URL-/Beispiel-Achse: F-1, F-3.)
2. **Gegenrichtung v5.6.0:** jeder im Repo referenzierte
   `baseline/v5.6.0/…`-Pfad existiert im vendorten Baum (skriptgestützte
   Existenz-Prüfung aller Fundstellen: 0 fehlend); `make doc-check` (396/0)
   bestätigt zusätzlich alle Anker in den retargeteten Zielen (u. a.
   `grundlagen-harness-dateien.md#template-schichtung…`,
   `modul-06-roadmap.md#roadmap-struktur…` existieren in v5.6.0 unverändert).
3. **Vendoring-Integrität:** `.harness/baseline/` enthält ausschließlich `v5.6.0/`
   mit `regelwerk/`, `templates/`, `SHA256SUMS` (51 Zeilen = 51 Dateien);
   `--verify` offline grün inkl. Manifest-Deckungs-Prüfung; kein v5.0.0-Rest
   (weder Dateien noch leere Verzeichnisse).
4. **MR-026-Bauform:** Felder (Status/Ersetzt/Datum/Geltungsbereich/Adaption/
   Begründung/Auflösungs-Trigger) analog MR-023; Index-Zeile vorhanden; der
   `<a id>`-Voll-Slug `mr-026--baseline-pin-hebung-auf-v560-dritter-nachtrag-zu-mr-011-nachtrag-zu-mr-023`
   entspricht der Slug-Bildung des tatsächlichen H1 (Kleinschreibung,
   Interpunktion/Punkte entfallen — „v5.6.0"→„v560" —, Leerzeichen→Bindestriche,
   Doppel-Bindestrich nach dem Gedankenstrich); die §Baseline-Ketten-Links darauf
   lösen auf (doc-check grün). Die Abgrenzung zu MR-023 ist mit dessen eigenem
   Auflösungs-Trigger („permanent, solange … self-contained … vendored") formal
   konsistent — der Trigger bindet das Layout, nicht den Pin-Wert (Rest: F-5, F-6).
5. **Tombstone-Deckung (Unter-Deckung):** die gebrochenen v5.0.0-Link-Fundstellen
   sind vollständig gedeckt — ADR-0047 (1 Link), ADR-0048 (3 Links), drei Reports;
   kein gebrochener Link bleibt (doc-check 396/0). Alle `in:`-Ziele der neuen
   Einträge sind eingefrorene Dateien — keine lebende Datei wird mit-entschuldigt
   (Über-Deckung innerhalb der eingefrorenen Klasse: F-4). Die v1.4.0-Alt-Tombstones
   bleiben unberührt und decken weiter genau ADR-0022/slice-080/slice-081.
6. **Skript-Fix (960600d):** Absolut-Pfad-Fallunterscheidung korrekt für beide
   Modi — vendor: repo-relatives `dest` wird über `$root` verankert (cwd ist per
   `git rev-parse --show-toplevel` die Repo-Wurzel, Verhalten unverändert);
   check-latest: absolutes mktemp-`dest` bleibt unangetastet (zuvor landete die
   Kopie unter `<repo>/<mktemp-Pfad>`). Die Under-Copy-Barriere
   (`src_n`/`dst_n`-Vergleich) sitzt unverändert im vendor-Pfad. `xargs -r`
   macht den leeren Fall zum SKIP statt Phantom-Drift (`up` leer ⇒
   `[ -n "$up" ]` schlägt fehl ⇒ `authenticity` bleibt `skip`; nur F-7 kosmetisch).
   `extract_both_trees` hat genau die zwei in-Skript-Aufrufer (Zeilen 160, 210);
   kein Makefile-/CI-/Skript-Aufrufer außerhalb. Die Schwester-Stelle Zeile 222
   (`xargs sha256sum` ohne `-r` bei der Manifest-Generierung) ist **kein**
   Finding: eine leere Liste ist dort strukturell unerreichbar (modul-00-Marker
   garantiert ≥ 1 Datei) und würde zusätzlich vom `verify`-Null-Guard
   („0 Dateien — leeres/kaputtes Vendoring") gefangen.
7. **Kein Konformitäts-Überclaim:** MR-026 grenzt explizit ab („diese MR hebt den
   Pin, sie behauptet keine Konformität"), der Slice-Plan §3 schließt den
   Abgleich aus, beide Commit-Messages behaupten keine v5.6.0-Konformität der
   eigenen Dokumente — Etappe-B/C-Grenze eingehalten.
8. **Skript-Kommentar-Spiegel:** die v5.0.0-Nennungen im Kopf von
   `fetch-baseline-cache.sh` (Zeilen 4, 15, 38) sind korrekte Historie bzw. ein
   weiterhin gültiges Tag-Format-Beispiel — kein Retarget nötig.
9. **Gates:** `make doc-check` selbst gefahren: 396 Dateien, 0 Befunde (Endstand
   `8a0477a` == Arbeitsbaum, clean). Der volle `make gates`-Lauf und die
   DoD-Abhakung sind Verifikations-Rolle, nicht Teil dieses Reviews.

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 2 | F-1, F-2 |
| LOW | 3 | F-3, F-4, F-5 |
| INFO | 2 | F-6, F-7 |

## Verdikt

**Blockierend (MEDIUM) — vor der Slice-Closure zu klären.** Kein HIGH; die
Substanz der Etappe A ist belegt in Ordnung (Vendoring integer, Verweis-Hebung
auf der gate-gedeckten Pfad-Achse vollständig, Tombstones ohne Unter-Deckung,
Skript-Fix korrekt). Blockierend sind die zwei MEDIUMs: der halb gehobene
Spiegel in `harness/README.md` (F-1) widerspricht der Commit-Behauptung
„alle lebenden Verweise", und der Commit-Schnitt (F-2) hinterlässt — solange
ungepusht noch korrigierbar — einen gate-roten, falsch beschrifteten
Zwischenstand in der Historie einer Migration, die sich auf Commit-Belege
beruft. Beide sind vor Push/Closure mit geringem Aufwand behebbar; LOW/INFO
sind nice-to-fix bzw. Notiz.
