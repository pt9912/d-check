# ADR-0033 — ai-harness-Vorlage: Modulset an die gelebte Konvention angleichen (Template-Eignungs-Kriterium)

**Status:** Proposed
**Datum:** 2026-07-06
**Autor:** pt9912
**Bezug:** [`DC-FA-CLI-006`](../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)
(die geänderte Anforderung; Spezifikation
[§`DC-FA-CLI-006.a`](../../../spec/spezifikation.md#dc-fa-cli-006a--konfigurations-vorschlag));
[ADR-0015](0015-suggest-config-id-prefix.md) (die andere Achse derselben
Anforderung, Anforderungs-Präfix). Vor-Schärfungen: die reservierten Quellen kamen mit
Lastenheft 0.18.0, die Kommentar-Klassifikation der situativen Module mit
Lastenheft 0.26.0 — diese ADR revidiert die dortige Einordnung von
`spans`/`hostpaths`.

## Kontext

Die reservierten Quellen `ai-harness` / `ai-harness-init` emittieren eine an die
ai-harness-course-Konvention angelehnte `.d-check.yml`-Vorlage. Ihr Modulset ist
seit Einführung der Vorlage (Lastenheft 0.18.0) **fix**: `links`, `anchors`,
`ids`, `matrix`, `codepaths` (fünf Module). Die Schärfung auf
Lastenheft 0.26.0 nannte die übrigen als „situative opt-in-Module" in einem
Kommentar mit Verweis auf `--print-config`.

Zwei Beobachtungen machen eine Angleichung fällig:

1. **Die Vorlage bildet weniger ab, als die Konvention lebt.** d-checks eigene
   `.d-check.yml` (der Dogfood, also die Konvention „in echt") aktiviert im
   Standard-Scan **acht** Module —
   `links, anchors, ids, matrix, codepaths, spans, hostpaths, versions` — und
   konfiguriert zusätzlich `vcs`, `commits`, `planning`, `targets` als dedizierte
   Gates. Die Vorlage liefert nur fünf. `spans`/`hostpaths` stehen also in der
   0.26.0-„nicht aktiviert"-Liste, obwohl das Referenz-Repo sie im Default führt —
   ein Selbstwiderspruch zwischen gepredigter und gelebter Konvention.
2. **Der normative Vertrag deckt die emittierte Ausgabe nicht 1:1.** Die
   kanonische Vorlage in
   [§`DC-FA-CLI-006.a`](../../../spec/spezifikation.md#dc-fa-cli-006a--konfigurations-vorschlag)
   enthält weder die „Weitere opt-in-Module"-Kommentarzeile noch den
   `codepaths`-`exempt-paths`/`ignore-refs`-Block, die der Renderer tatsächlich
   ausgibt. Das „Warum nicht" der ausgeschlossenen Module lebt nur in der
   Lastenheft-Historie und im Code-Kommentar — nicht im operativen §-Body. Der
   Code läuft der Norm voraus.

Beides gehört in **eine** Entscheidung, weil beide dieselbe Frage betreffen:
**welche Module gehören in die kanonische Vorlage — und woran macht man das
fest?**

## Entscheidung

### 1. Ein explizites Eignungs-Kriterium

Ein Modul gehört in das emittierte ai-harness-Modulset genau dann, wenn **alle**
gelten:

- **K1 Konventions-kanonisch** — die ai-harness-course-Konvention nutzt es
  tatsächlich (Beleg: d-checks eigener Dogfood `.d-check.yml`).
- **K2 Ableitungs-frei oder konventions-feste Config** — es läuft mit Defaults
  oder mit Config, deren Werte **Konventions-Konstanten** sind (die das
  Harness-Layout fixiert: Pfade, Marker), **nicht** repo-spezifisch-unableitbare
  Werte.
- **K3 Baum-Scan-tauglich** — es arbeitet in einer statischen `.d-check.yml`
  über den gescannten Baum, **ohne** eine Laufzeit-Range/-Flag (`--range`,
  `--staged`, `--commit-msg`).
- **K4 Hermetisch/netzlos** — es verletzt die Read-only-/Kein-Netz-Erwartung des
  Default-Laufs nicht ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

### 2. Anwendung auf die 17 Module

- **Ins fixe Modulset (beide Modi):** `spans`, `hostpaths` erfüllen K1–K4
  (konfigfrei, hermetisch, Baum-Scan; d-check führt sie). Sie werden dem fixen
  Set hinzugefügt: `links, anchors, ids, matrix, codepaths, spans, hostpaths`.
  **Dies revidiert die 0.26.0-Einordnung** für genau diese zwei.
- **Repo-bewusster Block (wie `matrix`):** `planning` erfüllt K1–K4 — seine
  einzige nötige Config ist der `roadmap`-Pfad (`docs/plan/planning/in-progress/roadmap.md`),
  `heading`/`marker`/`slice-glob` sind Konventions-Defaults. Aktiv, wenn Roadmap
  + `docs/plan/planning/` existieren (`ai-harness`), sonst auskommentiert mit
  Hinweis; im Voll-Kanon (`ai-harness-init`) aktiv. Es ist gerade die
  Layout-Konvention, die die Vorlage in `ids`/`matrix` ohnehin schon kodiert.
- **Nach `--print-mk` verwiesen, NICHT ins Modulset:** `vcs`, `commits` scheitern
  an **K3** (Laufzeit-Range). Sie gehören nicht in ein statisches `modules:`,
  sondern sind über `--print-mk` (`doc-immutable`, `doc-commits`) verteilt. Die
  Vorlagen-Kommentarzeile verweist zusätzlich auf `--print-mk` (bisher nur
  `--print-config`).
- **Draußen, bewusst vertagt (dokumentiert):** `versions` scheitert an **K2**
  (`versions.pin-pattern` ist repo-spezifisch — d-check: `ghcr.io/…`); `targets`
  liegt an der K2-Grenze (`authority`/`doc-tables` sind semi-spezifisch). Beide
  bleiben aus dem fixen Set, werden aber im §-Body als **geprüft und begründet
  vertagt** benannt (statt still im Code entschieden). `targets` ist ein
  späterer repo-bewusster-Block-Kandidat.
- **Draußen (unverändert):** `external` (K4, Netz), `diagrams` (K2,
  unableitbare `patterns`/`defined-in` — der Muster-Fall der Norm),
  `pins`/`immutable` (pro-Marker; wirken nur, wo Marker stehen — nichts im
  Gerüst zu deklarieren), `tracked` (K3/K4, fail-closed ohne `.git` — bricht
  `ai-harness-init` fürs leere Repo).

### 3. Die Norm deckt die Ausgabe 1:1

Die kanonische Vorlage in
[§`DC-FA-CLI-006.a`](../../../spec/spezifikation.md#dc-fa-cli-006a--konfigurations-vorschlag)
wird um die tatsächlich emittierten Teile ergänzt (die „Weitere opt-in-Module"-
Kommentarzeile inkl. `--print-mk`-Verweis und der `codepaths`-Block), und der
operative Body benennt das Eignungs-Kriterium (K1–K4) samt der **geschlossenen**
Aktiv-Menge. Damit läuft der Code der Norm nicht länger voraus.

## Konsequenzen

- **Geänderte Anforderung (Lastenheft-CR, Versions-Bump + Historie).** Das
  ausgegebene Modulset ändert sich (nutzersichtbar) → geänderte Anforderung im
  Lastenheft, **kein** neues Modul, **keine** neue `DC-*`-ID, **kein** neuer
  Grund-Code. **Release** (v0.39.0, geänderte `--suggest-config`-Ausgabe).
- **Renderer (`internal/hexagon/core/app/suggest.go`).** `renderHarness` nimmt
  `spans`/`hostpaths` ins fixe `modules:`; ein `renderHarnessPlanning` (repo-
  bewusst, analog `renderHarnessMatrix`) kommt hinzu; die „Weitere
  opt-in-Module"-Kommentarzeile verweist zusätzlich auf `--print-mk` und nennt
  `versions`/`targets` als bewusst vertagt. Determinismus
  ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)) und
  Read-only ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit))
  unberührt (nur mehr Ausgabe, kein neuer Eingabe-Scope).
- **Akzeptanz-Beleg.** Das dekodierte Gerüst muss weiterhin über den eigenen
  Parser laufen (Round-Trip); der neue `planning`-Block wird im repo-bewussten
  Modus gegen ein Fixture mit/ohne Roadmap geprüft (aktiv ⟷ auskommentiert).
- **Abgrenzung.** Die Änderung betrifft **nur** die `ai-harness`-Vorlage
  (`renderHarness`), **nicht** den generischen `--suggest-config <quelle>`-Modus
  (`renderSuggestion` mit seinem `codepaths`/`spans`/`hostpaths`-Probelauf) und
  **nicht** d-checks eigene `.d-check.yml`.
- **Reversibel** im Verhalten (advisory-Ausgabe), aber Vertrags-Änderung — daher
  Lastenheft-CR statt stiller Code-Änderung.

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-07-06 | Entwurf (slice-065, welle-54; Auftraggeber-Analyse „welche Module werden in `--suggest-config ai-harness` nicht verwendet, und warum nicht"): das Eignungs-Kriterium K1–K4 macht die Template-Modul-Zuordnung explizit; `spans`/`hostpaths` ins fixe Modulset (revidiert die Lastenheft-0.26.0-Einordnung), `planning` als repo-bewusster Block, `vcs`/`commits` → `--print-mk` (Laufzeit-Range, K3), `versions`/`targets` dokumentiert vertagt (K2); zugleich deckt die kanonische Vorlage der Spezifikation die emittierte Ausgabe 1:1 (Normativitäts-Spalt geschlossen). Lastenheft-CR, Release v0.39.0 geplant. Status Proposed. |
