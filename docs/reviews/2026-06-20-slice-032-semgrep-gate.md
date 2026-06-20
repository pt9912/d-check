# Review — slice-032 Implementierung (semgrep als hermetisches Gate)

## Kopf-Metadaten

- **Review-Art:** Code-Review (Working-Tree-Diff gegen Plan/ADR/Anforderungen/
  Hard Rules — **kein Verifier**: DoD-Abhaken und Gate-Lauf-Bestätigung sind
  nicht Gegenstand; Gates werden NICHT als grün angenommen).
- **Datum:** 2026-06-20
- **Gegenstand:** Working-Tree-Diff (unstaged) der slice-032-Umsetzung plus die
  untrackten Dateien `tools/semgrep.sh` (NEU),
  `docs/plan/adr/0010-semgrep-hermetisches-gate.md` (NEU),
  `docs/plan/planning/done/slice-032-semgrep-gate.md` (NEU). semgrep wird als
  neues Gate `make semgrep` in `make gates` (und damit `make ci`) eingehängt:
  ein gepinntes semgrep-Image (`semgrep/semgrep:1.167.0`) scant über ein
  gepinntes, lokal **außerhalb des Repos** gecachtes Regelset
  (`semgrep/semgrep-rules`@Commit `d41fb34…`, Umfang `go/lang/security`) netzlos
  (`docker run --network none --disable-version-check`, `--error`). Geändert:
  `Makefile` (`.PHONY` + `semgrep`-Target + `gates`-Prerequisite/`echo`),
  `AGENTS.md` §4, `harness/README.md` Sensors-Tabelle, `docs/plan/adr/README.md`
  Index-Zeile.
- **Reviewer:** Claude (Agent), Skill `.harness/skills/reviewer.md` v1.0.0.
- **Modell:** `claude-opus-4-8[1m]`.
- **Eingangs-Kontext:** Slice-Plan
  `docs/plan/planning/done/slice-032-semgrep-gate.md`; ADR
  `docs/plan/adr/0010-semgrep-hermetisches-gate.md` (ergänzt ADR-0006);
  Anforderungen `DC-QA-02` (Determinismus — byte-identische Ausgabe bei
  identischer Eingabe) und `DC-QA-03` (Seiteneffektfreiheit/Netzlosigkeit — kein
  Netz außer Modul `external`); Hard Rules `AGENTS.md` §3 (insb. §3.1
  Docker/make-only, §3.2 Suppression-Verbot, §3.3 Lifecycle); Meta-Gate
  `tools/gate-consistency.sh`. Kein Vertrags-`DC-FA-*` (Prozess-/Qualitäts-Gate).
  **Die DoD-Abhakung des Slices lag bewusst nicht vor** — die Findings sind
  eigenständig gebildet. Tests/Gates wurden NICHT ausgeführt (make/Docker-only;
  Reviewer ist kein Verifier; `make gates` wird NICHT als grün angenommen).

## Findings

### HIGH

#### HIGH-1 — Kein Nachweis geladener Regeln: leeres/falsches Config-Verzeichnis ergibt stilles Grün

- **Kategorie:** HIGH
- **Quelle:** Reviewer-Skill (Repo-Anker HIGH: „Stilles-Grün-Pfad in einem
  Gate-Skript — Harness-Lüge"); `DC-QA-02`/ADR-0010 (Fitness Function „rot =
  Befund")
- **Pfad:** `tools/semgrep.sh:33` (Cache-Guard) und `tools/semgrep.sh` Zeilen 47–53
  (Scan ohne Regel-Nachweis)
- **Befund:** Das Gate verlässt sich darauf, dass
  `--config /rules/go/lang/security` ein nicht-leeres Regelset lädt, prüft das
  aber nirgends. Der Guard (`if [ ! -d "$RULES_DIR/$RULES_SUBSET" ]`, Zeile 33)
  testet nur die **Existenz** des Subset-Verzeichnisses, nicht, dass darin
  ladbare `.yaml`-Regeln liegen; es gibt keinen `find … -name '*.yaml' | wc`-,
  `--strict`- oder Regel-Zähl-Check. Existiert `…/go/lang/security` als
  Verzeichnis ohne ladbare Regel (z. B. nach einer künftigen Pin-Hebung, bei der
  der Upstream-Pfad umbenannt/verschoben wurde, oder bei einem Cache, dessen
  Subset-Ordner leer ist), lädt semgrep 0 Regeln und meldet „0 findings" → Exit
  0 → Gate grün, obwohl faktisch nichts gescannt wurde. `--error` (Zeile 51)
  ändert nur den Exit-Code **bei Befunden**, nicht bei 0 geladenen Regeln. Der
  per-Commit-Lauf täuscht damit Sicherheits-Abdeckung vor, die nicht stattfand.
- **Verifizierbar:** ja — ein Lauf mit präpariertem Cache, in dem
  `$RULES_DIR/go/lang/security` existiert, aber keine `.yaml`-Regel enthält
  (bzw. `SEMGREP_RULES_CACHE` auf ein solches Verzeichnis gesetzt), zeigt: das
  Gate endet mit Exit 0 trotz null geladener Regeln. (Das `.tmp`+`mv`-Muster
  schützt nur gegen einen **halb geschriebenen** Cache eines abgebrochenen
  Fetch, nicht gegen einen vollständig ausgecheckten, aber regel-leeren
  Subset-Pfad.)

### MEDIUM

#### MEDIUM-1 — ADR-Index beschreibt das Regelset als „vendortes" — Gegenteil der ADR-Entscheidung

- **Kategorie:** MEDIUM
- **Quelle:** ADR-0010 (Entscheidung Punkt 1 „nicht ins Repo vendort");
  `AGENTS.md` §2 Source Precedence (ADR-Index Rang 4 — kanonische Quelle)
- **Pfad:** `docs/plan/adr/README.md:24`
- **Befund:** Die neue Index-Zeile betitelt ADR-0010 als „semgrep als
  hermetisches Gate (**vendortes**, gepinntes Regelset)". Die ADR selbst
  entscheidet das Gegenteil: „Gepinntes, lokal gecachtes Regelset (**nicht ins
  Repo vendort**)" (`docs/plan/adr/0010-…:32`), führt „Regeln dauerhaft ins Repo
  vendoren" ausdrücklich als **verworfene** Alternative
  (`docs/plan/adr/0010-…:74`), und `tools/semgrep.sh:7` sowie
  `harness/README.md:57` formulieren ebenfalls „NICHT ins Repo vendort" /
  „außerhalb des Repos gecacht". Der ADR-Index ist eine Source-Precedence-Quelle
  (`AGENTS.md` §2); seine Kurzbeschreibung behauptet den Gegenstand der
  Entscheidung invers. Ein Leser, der den Index als Einstieg nimmt, erhält die
  falsche Provenienz-/Lizenz-Aussage (die ganze ADR-Begründung — keine
  Fremdlizenz-Mischung im MIT-Repo — hängt am Gegenteil von „vendort").
- **Verifizierbar:** ja — Textvergleich `README.md:24` („vendortes") gegen
  `docs/plan/adr/0010-…:32` („nicht ins Repo vendort") und `:74` (vendoren =
  verworfen) belegt den Widerspruch; kein Gate-Lauf nötig.

### LOW

Keine.

### INFO

#### INFO-1 — Image-Pin ist ein veränderlicher Tag, kein Digest (Determinismus-Restlücke, konventionskonform)

- **Kategorie:** INFO
- **Quelle:** `DC-QA-02` (Determinismus); ADR-0010 Konsequenz „Update-Politik …
  analog zur Pin-Politik der übrigen Toolchain"
- **Pfad:** `tools/semgrep.sh:20` / `:50`
- **Befund:** Der Scanner ist über den **Tag** `semgrep/semgrep:1.167.0`
  gepinnt, nicht über `@sha256:`-Digest. Ein neu unter denselben Tag gepushtes
  Image (anderes Scanner-Binary/gebündelte Engine) könnte bei identischem
  Repo-Stand und identischem Regel-Commit andere Befunde liefern — eine
  Rest-Determinismuslücke, die der Regel-Commit-Pin nicht abdeckt. Bewusst
  notiert als dokumentationswürdige Annahme: dasselbe Tag-statt-Digest-Muster
  gilt für `GO_VERSION`/`GOLANGCI_LINT_VERSION` und die Basis-Image-Pins des
  Repos; die ADR ordnet semgrep ausdrücklich in diese Pin-Politik ein. Kein
  Korrektheitsdefekt gegen die bestehende Konvention — daher INFO, nicht LOW.
- **Verifizierbar:** ja — `docker inspect` der Image-ID vor/nach einem
  hypothetischen Tag-Repush zeigt verschiedene Digests bei gleichem Tag; im
  Normalbetrieb nicht beobachtbar (lokaler Image-Cache).

## Negativbefunde (geprüft, ohne Befund)

- **`git fetch`-Fehler bricht das Gate (kein falsches Grün):** Der Setup-Block
  läuft unter `set -euo pipefail` (`tools/semgrep.sh:18`); `git … fetch …`
  (`:39`) und `checkout` (`:40`) tragen **kein** `|| true`. Schlägt der Fetch
  fehl (Netz weg, Pin ungültig), bricht der Skript mit Non-Zero ab, **bevor**
  `mv "${RULES_DIR}.tmp" "$RULES_DIR"` (`:41`) läuft — der Commit-gekeyte
  Cache-Pfad entsteht nie, der nächste Lauf versucht erneut zu holen. Ein
  abgebrochener Fetch hinterlässt nur `${RULES_DIR}.tmp` (vorher per `rm -rf`
  bereinigt, `:35`), nie einen halben „fertigen" Cache. (Die verbleibende
  Lücke — vollständiger Fetch eines regel-leeren Subset — ist HIGH-1, nicht
  dieser Pfad.)
- **Netzlosigkeit des Scans (`DC-QA-03`):** Der Scan-Container trägt
  `--network none` (`tools/semgrep.sh:47`); `--metrics off` und
  `--disable-version-check` (`:51`) unterbinden die beiden einzigen ausgehenden
  Calls (Telemetrie, Versions-Ping). `--disable-version-check` maskiert nichts
  Sicherheitsrelevantes — es ist ein Phone-Home-Ping auf die neueste
  semgrep-**Version**, kein Regel-Update; ohne den Flag liefe er unter
  `--network none` nur in einen Timeout (ADR-0010 Punkt 3). Das einmalige
  Cache-Holen (Host-`git`, Netz) steht oberhalb des `exec docker run`
  (`:33-42`) und ist damit sauber als Setup vom netzlosen Scan getrennt.
- **Hard Rule §3.1 (Host-`git`-Fetch zulässig):** `tools/semgrep.sh` Zeilen 37–40 nutzen
  Host-`git` für den Regel-Fetch. Das ist **kein** Verstoß: §3.1 verbietet
  „Host-Go und Host-Paketmanager (`go`, `pip`, `npm`, `cargo`, `apt`, `brew`)"
  und stellt explizit fest „Der Host braucht nur `git`, GNU `make`, `bash` und
  Docker" (`AGENTS.md:50`). `git` ist hier Kern-VCS-Fetch (wie ein Clone, analog
  Image-Pull/`go mod`-Setup), kein Sprach-Paketmanager. REFUTED mit §3.1-Zitat.
- **gate-consistency bleibt beidseitig konsistent:** `make semgrep` existiert im
  Makefile als `.PHONY`-Eintrag (`Makefile:41`) und als Target (`Makefile:82`),
  steht in `AGENTS.md` §4 als Tabellenzeile mit `` `make semgrep` ``
  (`AGENTS.md:112`) und in `harness/README.md` Sensors mit `` `make semgrep` ``
  (`harness/README.md:57`). Beide Richtungen des Meta-Gates (Doku→Makefile,
  Makefile→AGENTS §4) sind damit erfüllt. Die `DC-QA-03`-Modulliste der
  `.d-check.yml` ist unverändert (`modules: [links, anchors, ids, matrix,
  codepaths, spans, hostpaths]`, `.d-check.yml:7`) — enthält alle fünf vom
  gate-consistency geforderten Module und kein `external`; Check (3) bleibt
  intakt.
- **Determinismus-Kette (`DC-QA-02`), soweit ohne HIGH-1/INFO-1:** Bei geladenem
  Regelset garantieren Regel-Commit-Pin (`:21`) + cache-pfad-Keying am Commit
  (`:27`) + gepinntes Image identische Eingabe ⇒ semgrep ist auf festem
  Regel-/Engine-Stand deterministisch. Die verbleibenden Lücken sind separat
  notiert (regel-leerer Pfad: HIGH-1; Tag-statt-Digest: INFO-1).
- **`gates`-Einhängung korrekt platziert:** `semgrep` steht im
  `gates`-Prerequisite vor `gate-consistency record-gates` (`Makefile:120`), die
  `echo`-Bestätigungszeile nennt semgrep (`Makefile:121`); `record-gates` bleibt
  letzter Prerequisite (Nachweis nur bei grünen Gates). Der Kommentar zur
  `-j`-Ordnung (`Makefile:118-119`) bleibt gültig.
- **Doku-Kohärenz Umfang/FP (außer MEDIUM-1):** ADR-Text, `tools/semgrep.sh` Zeilen 15–17,
  `harness/README.md:57` und `AGENTS.md:112` nennen durchgängig Umfang
  `go/lang/security`, 0 Befunde, kein `--exclude-rule` nötig, breitere Packs
  bewusst ausgelassen (13 FP). Der `--config`-Pfad im Skript
  (`/rules/${RULES_SUBSET}` mit `RULES_SUBSET=go/lang/security`, `:22`/`:52`)
  deckt sich mit dem ADR-Umfang. Kein Widerspruch außer der vendort-Index-Zeile
  (MEDIUM-1).
- **ADR-Status/Index-Verlinkung:** ADR-0010 trägt `Status: Proposed`
  (`docs/plan/adr/0010-…:3`) konsistent mit der Index-Zeile (Spalte „Proposed",
  `README.md:24`) und der Geschichte-Zeile (`:113`); der Slice-Plan fordert
  `Accepted` erst zur Closure (DoD-Punkt). IDs in ADR und Slice erscheinen als
  verlinkte Backtick-Spans (`DC-QA-02`/`DC-QA-03` mit Anker, ADR-Querverweise);
  keine bare `DC-`/`ADR-`/`slice-`-ID außerhalb Link/Fence in den neuen
  Prosa-Zeilen erkennbar. (Doc-check/Anker-/Linkpflicht ist Verifier-Sache, hier
  nicht ausgeführt — nur Augenscheinprüfung.)
- **Lifecycle/§3.3:** Drei neue Dateien (`tools/semgrep.sh`, ADR-0010,
  Slice-Plan) sind reine Neuanlagen (keine git-mv-+-Inhalts-Kombination); §3.3
  nicht berührt. Der Slice lag zum Review-Zeitpunkt unter `docs/plan/planning/done/`
  (Status „open (geplant)").
- **Setup-Trennung außerhalb des Repos:** Cache-Wurzel ist
  `${XDG_CACHE_HOME:-$HOME/.cache}/d-check/semgrep-rules` (`:26`), nicht im
  Repo — die `.go`-Fixtures der Regeln stören weder `go list ./...` noch den
  d-check-Selbstscan/Dogfooding (ADR-0010 Konsequenz). Keine Interferenz mit dem
  `.d-check.yml`-Selbstscan beobachtbar.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 1 | 1 | 0 | 1 |

## Verdikt

**Blockiert** (HIGH-1 + MEDIUM-1). Das Hermetik-Design des Slices ist im Kern
tragfähig: netzloser Scan (`--network none` + `--disable-version-check`/`--metrics
off`), sauber als Setup getrenntes Cache-Holen via Host-`git` (§3.1-konform),
Commit-/Image-Pin für Reproduzierbarkeit, atomares `.tmp`+`mv` gegen halbe
Caches, `set -euo pipefail` lässt einen Fetch-Fehler hart abbrechen, und das
Meta-Gate bleibt beidseitig konsistent. Es blockiert jedoch **HIGH-1**: das Gate
prüft nicht, dass `--config` tatsächlich Regeln geladen hat — ein
regel-leeres/umbenanntes Subset-Verzeichnis (insb. nach künftiger Pin-Hebung)
führt zu „0 findings" → Exit 0 und damit zu stillem Grün, das Sicherheits-
Abdeckung vortäuscht (Harness-Lüge, Repo-HIGH-Anker). **MEDIUM-1** blockiert als
Doku-Inversion in einer Source-Precedence-Quelle: der ADR-Index betitelt das
Regelset als „vendortes", obwohl die ADR „nicht ins Repo vendort" entscheidet und
Vendoring ausdrücklich verwirft. INFO-1 (Tag statt Digest) ist eine
konventionskonforme Restlücke ohne Blockwirkung. Die Gate-Bestätigung obliegt der
getrennten Verifikation (hier NICHT als grün angenommen; Tests/Gates wurden nicht
ausgeführt).
