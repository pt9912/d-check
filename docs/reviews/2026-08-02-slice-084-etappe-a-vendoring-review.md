# Review — slice-084 Etappe A (Vendoring v1.4.0 -> v5.0.0)

**Datum:** 2026-08-02 · **Reviewer:** unabhängiger Frischkontext (adversarial) ·
**Gegenstand:** Baseline-Regelwerk-Migration, Etappe A (Vendoring + Pin/Pointer),
Commits `ae582c3` (Skript + Materialisierung) und `e885a4a` (atomare Umschaltung).

**Eingangs-Kontext:** Commit-Range `ae582c3..e885a4a`; Slice-Plan
`docs/plan/planning/in-progress/slice-084-etappe-a-vendoring.md`; SOLL-Schnitt
`docs/plan/planning/done/slice-083-regelwerk-v500-migration-analyse.md` §2.7;
Hard Rules `AGENTS.md` §3; Reviewer-Skill `.harness/skills/reviewer.md` v1.2.0.

**Eigene Gate-Läufe (read-only) zur Verifikation:**

- `make gates` -> Exit 0. Alle inneren Gates grün: doc-check 271/0, lint, test,
  arch-check, coverage-gate 94.20% (Schwelle 93%), semgrep 0/0, gate-consistency
  (targets) 271/0, planning-check 271/0.
- `docker run --rm --network none -v $PWD:/repo:ro d-check:latest` (HEAD) -> 271
  Dateien, 0 Befunde.
- Dieselbe Dogfood-Ausführung gegen einen Detached-Worktree auf `ae582c3` -> 271
  Dateien, 0 Befunde (v1.4.0 UND v5.0.0 beide präsent, kein Baseline-`ignore-refs`).
- `bash tools/harness/fetch-baseline-cache.sh --verify` -> Exit 0 (50 Dateien,
  vollständig; Manifest-Deckung ok).

---

## Findings

### LOW-1 — Under-Copy-Barriere ist blind für Nicht-Markdown-Dateien im regelwerk-Baum

- **kategorie:** LOW
- **quelle:** Maintainability (slice-083 §2.7 A1 „Under-Copy-Barriere")
- **pfad:** `tools/harness/fetch-baseline-cache.sh:90` / `:94` / `:209`
- **befund:** Der regelwerk-Kopierschritt filtert die Quelle auf `-name '*.md'`,
  und `src_n` zählt für den regelwerk-Zweig ebenfalls nur `*.md`; `dst_n` zählt
  alle Dateien im vendorten regelwerk-Baum, der aber nach dem `*.md`-Filter
  entstanden ist. Eine künftige Baseline, die im regelwerk-Baum eine
  Nicht-Markdown-Datei mitliefert (z. B. ein von einem Regelwerk-Doc eingebettetes
  Diagramm-Asset), wird still fallengelassen: sie fließt weder in `src_n` noch in
  die Kopie noch in `dst_n`, also bleibt `src_n == dst_n` und die Barriere grün.
  `--verify` bleibt ebenfalls grün (Manifest über den kopierten Bestand), und
  `--check-latest` Teil B extrahiert die Upstream-Bytes identisch `*.md`-gefiltert,
  ist also gegen dieselbe Auslassung blind. Failure-Szenario: Bump auf ein
  v-next-Regelwerk mit eingebettetem Nicht-md-Asset -> vendorte Lese-Form
  unvollständig, ohne dass irgendein Skript-Modus das meldet.
- **verifizierbar:** nein durch den heutigen Gate-Lauf — das v5.0.0-regelwerk ist
  rein Markdown (26/26 `.md`), die Auslassung ist heute leer. Bestätigbar nur über
  einen konstruierten Bump-Fixture mit Nicht-md-Datei im regelwerk-Baum.

### LOW-2 — `conventions.md` behauptet an zwei Live-Stellen „aktueller <tag> = v1.4.0"

- **kategorie:** LOW
- **quelle:** Maintainability / Doku-Drift
- **pfad:** `harness/conventions.md:517` (MR-017-Eintrag) und `harness/conventions.md:593`
  (MR-019-Eintrag)
- **befund:** Beide Einträge tragen weiterhin den Satz „aktueller `<tag>` =
  `v1.4.0`", während dieselbe Datei in §Baseline jetzt `v5.0.0` als Stand führt.
  Ein Agent, der den MR-019-Eintrag (die vendoring-tragende Adaption) statt §Baseline
  liest, entnimmt daraus den falschen aktuell vendorten Tag. Kein Gate greift
  (Prosa/Inline-Code, kein ghcr-Pin, also außerhalb des `versions`-Musters; kein
  Link). Der MR-021-Eintrag (`:685`) trägt zusätzlich ein illustratives
  `…/v1.4.0/…`-Beispiel — harmlos, aber gleicher Drift-Charakter.
- **verifizierbar:** nein (kein maschineller Sensor deckt Prosa-Tag-Aktualität ab).
- **Einordnung (nicht-blockierend):** Dieser Drift ist die ausdrücklich nach
  Etappe C geschnittene Arbeit: die Analyse §2.7 lässt die MR-017/018/019/022-Prosa
  explizit draußen, und der neue MR-023-Eintrag (`Auflösungs-Trigger`) trackt die
  vollständige Angleichung an das v5.0.0-Layout als eigene Etappe. Erwartungskonform,
  daher kein MEDIUM. Anmerkung: die Deferral-Notiz nennt MR-017/018/019/022, nicht
  MR-021 — dessen `…/v1.4.0/…`-Beispiel ist im Tracking nicht explizit erfasst.

### INFO-1 — `check-latest` Teil B nutzt eine präfigierte Zuweisung `src_n=0 extract_both_trees`

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** `tools/harness/fetch-baseline-cache.sh:155`
- **befund:** Der Content-Drift-Teil ruft die Extraktionsfunktion als
  `src_n=0 extract_both_trees …`. Bei einer Shell-Funktion ist die präfigierte
  Zuweisung invocation-lokal; die interne globale Neuzuweisung `src_n="$(…)"`
  propagiert daher nicht zuverlässig nach außen. Für die Korrektheit ist das
  irrelevant, weil Teil B `src_n` nicht konsultiert (der Vergleich `up` wird aus dem
  extrahierten Dateisystem gebildet, nicht aus `src_n`); die Under-Copy-Barriere mit
  `src_n` sitzt ausschließlich im re-vendor-Pfad, wo die Funktion ohne Präfix
  gerufen wird und die globale Zuweisung greift. Kein Failure, reine Lesbarkeitsnotiz.
- **verifizierbar:** n/a (kein Verhaltens-Effekt).

---

## Negativbefunde (geprüft, ohne blockierenden Befund)

- **Übersehene v1.4.0-Live-Verweise (Achse 1):** geprüft repo-weit. Die einzigen
  lebenden Markdown-Links auf den entfernten `.harness/baseline/v1.4.0/…`-Pfad sind
  die drei getombsteten (`docs/plan/adr/0022-…`, `done/slice-080-…`,
  `done/slice-081-…`). Alle übrigen v1.4.0-Vorkommen sind (a) historische Pin-Bump-
  Prosa in `harness/conventions.md` (MR-011/012/016, korrekt als Vergangenheit),
  (b) Inline-Code in `slice-083`/`slice-084`/Skript-Kommentar (keine Links),
  (c) ventilierte `docs/reviews/**`, (d) `go.sum` (Fremd-Dependency). doc-check
  271/0 belegt: keine brechende Live-Referenz.
- **Retarget-Korrektheit (Achse 2):** geprüft. Alle sechs retargeten Ziele +
  Anker existieren und lösen auf — manuell verifiziert und per Gate: reviewer-Skill
  -> `grundlagen-referenz-richtung.md` Heading „Referenz-Richtung (SDP): wer darf
  wen referenzieren" (Slug `referenz-richtung-sdp-wer-darf-wen-referenzieren`);
  planning/README -> `grundlagen-klassifikation.md` Heading „Steering Loop" (Slug
  `steering-loop`); roadmap -> `modul-06-roadmap.md` (bare); harness/README ->
  `regelwerk/`-Verzeichnis + `modul-13-quality-gates.md` + `modul-14-docker-harness.md`;
  conventions §Baseline -> `v5.0.0/regelwerk/`-Verzeichnis. Der links-Modul löst
  belegbar in `.harness/baseline/**` hinein auf (nur deshalb brauchte die Umschaltung
  den Tombstone), Anker zusätzlich per Hand gegen die Headings geprüft.
- **Tombstone-Ventil (Achse 3):** geprüft. Schema korrekt (`in` = String je Datei,
  `refs` = Liste). Nicht zu breit: jeder Eintrag ist auf genau eine eingefrorene
  Datei gescoped, ein echter künftiger v1.4.0-Fehler in einer anderen Datei bliebe
  ungedeckt; innerhalb der drei Dateien deckt nur das `.harness/baseline/v1.4.0/**`-
  Muster, andere brechende Refs dieser Dateien blieben sichtbar. Nicht zu eng:
  load-bearing — die drei Refs zeigen auf gelöschte Dateien, und nur weil das Ventil
  greift, ist doc-check grün (Gegenprobe: der Worktree auf `ae582c3` ist ohne dieses
  Ventil grün, weil dort v1.4.0 noch existiert). Fällt v1.4.0 real weg (jetzt der
  Fall), greift das Ventil nicht leer, sondern unterdrückt genau die drei sonst
  fälligen `target-missing`-Befunde.
- **MR-023 (Achse 4):** geprüft. Jede Behauptung des Eintrags ist umgesetzt
  (beide Bäume vendored; Skript aufs Bundle gehoben; entfallene Quellzeiger
  umgeschrieben; Tombstone gesetzt). Der Heading-Slug
  `mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout` löst aus
  §Baseline, §Adoptierte-Quellen und harness/README auf; die `ids`-Linkpflicht für
  alle MR-023- und DC-FA-REF-001-Erwähnungen ist erfüllt (doc-check grün). Der
  „bekannt offen für Etappe C"-Vermerk deckt sich mit dem SOLL-Schnitt der
  abgenommenen Analyse §2.7 und kaschiert nichts, das in Etappe A gehört.
- **Skript-Korrektheit (Achse 5):** geprüft. Tolerantes Entpacken am
  `modul-00`-Marker mit templates-Geschwister-Pflicht; rekursive Kopie erhält die
  Struktur. SHA256SUMS-Rekonstruktion in `--check-latest` Teil B ist byte-vergleichbar:
  beide Seiten tragen die Präfixe `regelwerk/`+`templates/` (jeweils `cd` in den
  Wurzel-Baum + `find regelwerk templates`), beide werden final `LC_ALL=C`-
  liniensortiert (das Manifest ist pfadsortiert geschrieben, der zusätzliche
  Schluss-Sort auf beiden Seiten kanonisiert auf Zeilen-Sortierung) — kein
  Ordnungs-Fehlalarm. Leerer `up` -> `authenticity=skip` (kein Falsch-Grün). Kein
  stilles Grün im gate-relevanten Pfad: `--verify` meldet fehlendes Manifest / 0
  Dateien / Manifest-Deckungs-Bruch je mit Exit 1; `--check-latest` ist per Design
  informativ/kein Gate (SKIP je Teil bei Netz-/Werkzeug-Ausfall, Exit = schlimmster
  Fall 4>3>0). `set -e` global, `check_latest` bewusst unter `set +e`; Kopierfehler
  im `while`-Loop propagieren über die Subshell (kein verschluckter `cp`-Fehler).
  (Rest-Notizen: INFO-1 oben; LOW-1 oben.)
- **Entfallene Quellzeiger (Achse 6):** geprüft an den vendorten v5.0.0-Dateien.
  `agents-regelwerk.md` ist im v5.0.0-Bundle nicht vorhanden (Regelwerk-Baum: 8
  `grundlagen-*` + 17 Module + `README.md`, kein Agenten-Destillat). Die frühere
  Sammeldatei `konventionen.md` ist real in `grundlagen-*` aufgeteilt (u. a.
  `grundlagen-referenz-richtung.md`, das die §Referenz-Richtung (SDP) trägt). Die
  Umschreibung statt Retarget ist damit belegt.
- **Per-Commit-Grün + Atomizität (Achse 7):** geprüft. `ae582c3` allein grün
  (271/0, beide Bäume präsent, Pointer noch auf v1.4.0, kein Baseline-`ignore-refs`
  nötig); `e885a4a`/HEAD grün (271/0). Da ein Commit atomar angewandt wird, existiert
  kein `target-missing`-Zwischenstand zwischen Alt-Stand-Entfernung und Pointer-Umzug.
- **DoD-Abgleich (informativ, nicht Reviewer-Rolle):** die sechs
  Live-Fundstellen-Zählung (§2.7, nach Dokument) und die „7 Live-Pfad-Links"-Zählung
  der Commit-Message (nach einzelnem Link, harness/README trägt drei) sind
  konsistent, keine übersehene Stelle.

---

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 1 |

---

## Verdikt

**Abnahmereif.** Kein HIGH, kein MEDIUM. Die Etappe erfüllt ihren SOLL-Schnitt:
beide Bäume netzlos vendored, Alt-Stand atomar entfernt, Pin und alle Live-Pointer
auf v5.0.0 (Ziele + Anker verifiziert), entfallene Quellzeiger umgeschrieben,
eingefrorene Historie korrekt und load-bearing getombstoned, Vendoring-Adaption
unter MR-023 deklariert. `make gates` (Exit 0, 271/0), der Dogfood-Selbstscan und
`--verify` sind unabhängig reproduziert; die Per-Commit-Grün-Eigenschaft ist gegen
`ae582c3` empirisch bestätigt. Die beiden LOW-Befunde blockieren nicht: LOW-2 ist
die ausdrücklich nach Etappe C geschnittene und in MR-023 getrackte Prosa-Drift,
LOW-1 eine latente Barriere-Lücke, die am heutigen rein-markdown-Regelwerk leer ist.
Empfehlung: LOW-1 und den MR-021-Teil von LOW-2 in den Etappe-C-Track aufnehmen.
