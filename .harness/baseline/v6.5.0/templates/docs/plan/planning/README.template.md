# Planning — <Projektname>

> **Template-Hinweis.** Vorlage für `docs/plan/planning/README.md`. Kopiere
> nach `docs/plan/planning/README.md`, ersetze `<Platzhalter>` und lösche
> diesen Block. **Derivativ:** dokumentiert die Konvention; Quelle der
> Wahrheit sind die Dateien in den Verzeichnissen selbst.

Slice-Lifecycle: `open/` → `next/` → `in-progress/` → `done/`.

Reine `git mv`-Commits beim Wechsel zwischen Verzeichnissen — siehe Hard
Rule „git mv + Inhaltsänderung = zwei Commits" in
[`../../../AGENTS.md`](../../../AGENTS.md).

## Lifecycle-Bedeutungen

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand ist das Verzeichnis, kein Status-Feld.

| Verzeichnis | Bedeutung |
|---|---|
| `open/` | Geplant, noch nicht priorisiert. Keine Garantie auf Umsetzung. |
| `next/` | Als Nächstes priorisiert. Verantwortlicher zugeordnet (`Verantwortlich:`-Feld im Slice-Kopf). |
| `in-progress/` | Beansprucht: Der `git mv` hierher liegt auf dem **Hauptzweig, vor der Arbeit** — Branch/PR entsteht danach. |
| `done/` | DoD erfüllt, gemerged, Closure-Notiz vorhanden. |

## Slices vs. Wellen — zwei Ablagen, dieselbe Regel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

- **Slices** tragen ihren Zustand über das **Verzeichnis**
  (`open/` → `next/` → `in-progress/` → `done/`).
- Eine **Welle** (Bündel von Slices) ebenso: Der Zustand ist die
  Verzeichnis-Position, kein `Status:`-Feld. Der Welle-Plan (`<welle-id>.md`)
  liegt **flach** in `planning/`, solange die Welle läuft, und wandert bei
  Closure per `git mv` nach `done/` — neben seine
  `welle-<NN>-results.md`. Den aktiven Durchlauf `open/` → `next/` →
  `in-progress/` durchläuft er nicht; `done/` ist sein einziges
  Lifecycle-Verzeichnis. **Geplante** Wellen haben noch keine Datei — sie
  stehen in der Roadmap, die auch Sequenzierungs-Autorität bleibt
  ([`in-progress/roadmap.md`](in-progress/roadmap.md): Meilensteine, nächste
  Wellen, Zeiger auf die offenen).
- Der aktive Durchlauf `open/` → `next/` → `in-progress/` nimmt ausschließlich
  **Slices** auf; `done/` archiviert **zusätzlich** abgeschlossene
  **Nicht-Slice-Records** — Welle-Plan und Welle-Closure
  `done/welle-<NN>-results.md`. Aufgelöste Carveouts wandern **nicht** hierher,
  sondern in ihr eigenes `docs/plan/carveouts/done/` (Baseline-Regelwerk
  `modul-07-carveouts.md`).

Neben den Lifecycle-Verzeichnissen liegt **flach** in `planning/` das
**Beobachtungs-Register** (`observations/`): Es trägt den
Steering-Loop-Zähler, wird bei jeder Slice-Closure fortgeschrieben und
überlebt jede Welle (Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register).

Ein Repo, das aus einem **Brownfield-Bootstrap** kommt, trägt hier flach ein
zweites Register: `reconciliation.md` mit den offenen Funden des Rückbaus. Es
wird beim Auflösen fortgeschrieben und ist leer, sobald alle Sub-Areas
graduiert sind (Baseline-Regelwerk `modul-02-harness-bootstrap.md`
§Das Reconciliation-Register). Greenfield-Repos haben die Datei nicht.

## Aktueller Stand

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — kein Snapshot; der Stand ergibt sich aus den
Verzeichnissen.

Nicht als Snapshot hier eintragen — der Stand ergibt sich aus den
`open/`/`next/`/`in-progress/`/`done/`-Verzeichnissen (optional ein
`plan-status`-Target wie im Kurs-`lab/example`), sonst driftet die Tabelle.

## Roadmap

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

Siehe [`in-progress/roadmap.md`](in-progress/roadmap.md).
