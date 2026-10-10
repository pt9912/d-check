# Verifikation — slice-273: `golang.org/x/net` auf v0.60.0 — vier HTTP/2-CVEs

**Rolle:** Verifier (Baseline-Regelwerk `modul-11-verification.md`): DoD und Plan gegen den tatsächlichen Stand, nicht Diff gegen Plan
**Gegenstand:** slice-273 auf `e098ed10`; DoD-Punkt 1
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-10
**Eingangs-Kontext:** Slice-Plan §1, §2, §4, §6, §8; Commit `e098ed10` (Diff und Botschaft als
Behauptung, nicht als Beleg); [ADR-0066](../plan/adr/0066-cve-scan-gegen-das-publizierte-image.md).

Alle Läufe in einem frischen Klon bzw. Wegwerf-Klonen unter dem Session-Scratchpad, nur über
`make`/Docker. Am Arbeitsbaum ist nichts geändert außer dieser Datei.

## Messungen (Kommando und Ergebnis)

| # | Kommando | Ergebnis |
|---|---|---|
| V1 | `git clone` (HEAD `e098ed10`), `make gates` | Exit 0, `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; lint `0 issues.`; alle Testpakete `ok`; `coverage-gate: OK — Coverage 94.80% erfüllt Schwelle 93%`; semgrep `Ran 55 rules on 69 files: 0 findings.`; doc-check `0 Befund(e)` |
| V2 | `make build` im Klon | Exit 0, Image `d-check:latest` (`sha256:39b57546…`) |
| V3 | `docker save d-check:latest` → `aquasec/trivy:0.74.0@sha256:62b1e65e…` `image --input <tar> --list-all-pkgs --format json` (Vuln-DB `UpdatedAt 2026-10-10 01:03 UTC`) | Go-Binary führt `golang.org/x/net@v0.60.0`; Befunde nur `DLA-4792-1` (tzdata) und `GO-2026-5932` (`x/crypto` openpgp, ohne Fix); **keine** der vier CVEs, auch `CVE-2026-97032` nicht |
| V4 | Gegenprobe: `docker pull ghcr.io/pt9912/d-check:v0.86.0`, derselbe Trivy, dieselbe DB | `golang.org/x/net@v0.59.0`; `CVE-2026-78659`, `-78660`, `-78663`, `-78669` (und `CVE-2026-97032`) je `InstalledVersion v0.59.0`, `FixedVersion 0.60.0`. Der Scan sieht die vier — ihr Fehlen in V3 ist also ein Befund über das Image, nicht über den Scanner |
| V5 | `docker build --no-cache --target deps` (frisches `go mod download` gegen `go.sum`), danach eine Wegwerf-Stufe darauf mit `go mod verify` | Exit 0; `all modules verified` |
| V6 | Bewusstes Brechen: Wegwerf-Klon, `h1:`-Hash von `golang.org/x/net v0.60.0` in `go.sum` verfälscht, `docker build --no-cache --target deps` | Exit 1: `verifying golang.org/x/net@v0.60.0: checksum mismatch` / `SECURITY ERROR` im `go mod download` der deps-Stufe. Der Build fällt bei einem `go.sum`-Fehler aus dem richtigen Grund |
| V7 | `make tidy` im Klon, danach `git status --short` | Exit 0, Arbeitsbaum unverändert — `go.mod`/`go.sum` sind der Fixpunkt von `tidy` |
| V8 | `bash tools/image-scan.sh` (der Weg hinter `make image-scan`, eigener Cache, Vuln-DB `UpdatedAt 2026-10-09 07:05 UTC`) gegen die publizierten `:latest` | Skript-Exit 0; alle vier CVEs `UNKNOWN`, `CVE-2026-97032` `MEDIUM` — siehe Befund I1 |

## DoD-Punkt 1 — Hebung, `go.sum`, Gates, Trivy-Vollbericht

**Erfüllt.**

- `go.mod`: genau eine Zeile geändert, `golang.org/x/net v0.59.0` → `v0.60.0` (`// indirect`);
  `go.sum`: die beiden `x/net v0.60.0`-Zeilen ersetzen die v0.59.0-Zeilen. `go.sum` passt (V5,
  V7); der Build mit `-mod=readonly` läuft (V1, V2).
- `make gates` grün, selbst gefahren (V1).
- Trivy-Vollbericht gegen das neu gebaute Image ohne die vier CVEs (V3), mit Gegenprobe gegen
  `v0.86.0` (V4). Die Gegenprobe ist zugleich das Bewusste Brechen dieses DoD-Punkts: Der
  Vorzustand (`x/net` v0.59.0, publiziert) liefert mit demselben Scanner und derselben DB die
  vier Befunde, gebunden an genau dieses Modul und die Fix-Version 0.60.0.

**Plan-Konformität:** Diff beschränkt auf `go.mod`/`go.sum` (§3). Die Abgrenzung (§1) ist
eingehalten: weder go-git noch go-winio gehoben, `image-scan` unverändert. Die
Rückführungs-Bedingung (§4, neuere Toolchain nötig) trat nicht ein — der Build läuft mit der
gepinnten `GO_VERSION=1.27.2`. Das Risiko aus §6 (Nachzug von `x/sys`, `x/text`, `x/crypto`)
ist nicht eingetreten: `tidy` zog nichts nach (V7); sein Ausgang steht bei der Closure noch aus.

DoD-Punkte 2 und 3 (Review, Verifikation, Closure) sind nicht Gegenstand dieser Verifikation;
dieser Bericht deckt die Verifikations-Hälfte von Punkt 2.

## Befunde

- **I1 (INFO) — der Schweregrad hat sich innerhalb eines Tages bewegt.** Mit der Vuln-DB vom
  2026-10-09 (Cache von `image-scan`, V8) stehen **alle vier** CVEs auf `UNKNOWN`; mit der DB vom
  2026-10-10 (V4) stehen `-78659`/`-78660` auf `UNKNOWN`, `-78663` auf `MEDIUM`, `-78669` auf
  `HIGH`. Folge: Sobald der `image-scan`-Cache aktualisiert, meldet der Scan gegen `v0.86.0`
  `CVE-2026-78669` als behebbares HIGH und läuft rot, bis `v0.86.1` publiziert ist. Für die
  Beobachtung `BEO-ALL/scanner-vendor-severity-lag` ist das ein Beleg für die *Länge* des
  Rückstands (rund ein Tag), nicht nur für sein Vorkommen.
- **L1 (LOW) — „drei davon als `UNKNOWN`" deckt sich mit keiner der beiden Messungen.** Slice-Plan
  §1/§8 und die Commit-Botschaft nennen drei `UNKNOWN`; gemessen sind vier (DB 2026-10-09) bzw.
  zwei (DB 2026-10-10). In der Tabellenausgabe von Trivy sind die Schweregrad-Zellen der Zeilen
  2–4 zusammengefasst und leer — ein Lesefehler liegt nahe. Die Commit-Botschaft ist
  eingefroren; die Zahl gehört in Closure-Notiz und Evidence-Datei richtig gestellt, samt DB-Stand
  (Dokumentations-Regeln 14/15).
- **INFO — `go mod verify` steht nicht im Dockerfile.** Die deps-Stufe fährt nur
  `go mod download`. Das reicht für die Zusage „der Build fällt bei `go.sum`-Fehlern" (V6): der
  Download prüft jeden Modul-Hash gegen `go.sum`. `go mod verify` prüft darüber hinaus nur, ob
  der lokale Modul-Cache seit dem Download verändert wurde; im frischen Build-Container ist das
  gegenstandslos. Kein Handlungsbedarf.

## Verdikt

DoD-Punkt 1 **bestätigt**, aus eigenen Läufen, nicht aus dem Bericht des Implementers. Kein
blockierender Befund; L1 ist bei der Closure zu korrigieren.
