# Review R1 — slice-273: `golang.org/x/net` auf v0.60.0 — vier HTTP/2-CVEs im publizierten Image

**Review-Art:** Code (Diff gegen Plan, ADRs, Konventionen und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-273, Commit `e098ed10` (`git show HEAD`)
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-10
**Eingangs-Kontext:** Slice-Plan slice-273 (§1 Ziel und Abgrenzung, §3 Plan, §4
Rückführung, §6 Risiko, §8 Vorprüfungen);
[ADR-0066](../plan/adr/0066-cve-scan-gegen-das-publizierte-image.md) (Entscheidung 4:
Rot nur bei `CRITICAL`/`HIGH` mit Fix),
[ADR-0067](../plan/adr/0067-dependabot-als-hebender-kanal.md) (Dependabot als hebender
Kanal, indirekte Abhängigkeiten); `AGENTS.md` §3.1 (Docker/make-only), §5 Regel 15;
Beobachtung `BEO-ALL/scanner-vendor-severity-lag`. Vorherige Findings am selben
Gegenstand: die Hebungen von slice-166 und slice-167 (x/net- und x/crypto-Befunde,
Reichweite von Dependabot bei indirekten Abhängigkeiten).

## Messungen (Kommando und Ergebnis)

Gelesen bzw. gefahren über Docker und `make`; Scan-Ausgaben im Scratchpad, nicht im
Arbeitsbaum.

| # | Messung | Ergebnis |
|---|---|---|
| A | `git show HEAD --stat` | genau `go.mod` (1 Zeile) und `go.sum` (2 Zeilen ersetzt), keine weitere Datei |
| B | `make tidy` erneut auf HEAD, danach `git status --porcelain` | leer — `go.sum` ist ein Fixpunkt von `tidy`, kein Handedit-Rest |
| C | `go.sum`-Zeilen `golang.org/x/net v0.59.0/go.mod` und `v0.60.0/go.mod` | identischer Hash `h1:2DA/G1UfVbCpQPeWTmMPGY7Cs2PkBkwu743bVX5PIVg=` — das `go.mod` von x/net hat sich zwischen den Versionen nicht geändert, es kommen keine neuen Requirements dazu |
| D | eingebettete Build-Info (`grep -a 'dep␉golang.org/x/…'`) aus `/d-check` in `d-check:latest` und `ghcr.io/pt9912/d-check:v0.86.0` | neu: x/net v0.60.0; alt: v0.59.0; x/crypto v0.57.0, x/sys v0.48.0, go-git v5.19.2 in beiden gleich |
| E | Trivy 0.74.0 (Digest-Pin aus `tools/image-scan.sh`), `--scanners vuln --format json` gegen `v0.86.0` | x/net v0.59.0: CVE-2026-78669 `HIGH`, -78663 `MEDIUM`, -97032 `MEDIUM`, -78659 `UNKNOWN`, -78660 `UNKNOWN`; jeweils Fix 0.60.0 |
| F | dasselbe gegen `d-check:latest` (HEAD) | nur noch `DLA-4792-1` (tzdata) und `GO-2026-5932` (x/crypto, ohne Fix); kein x/net-Befund |
| G | `.harness/state/gates-passed.diffsha` gegen `tools/harness/working-tree-hash.sh` | gleich (`29603e45…`), geschrieben 04:47:51, Commit 04:48:35 |

## Findings

### F-1 — INFO · eine fünfte CVE ist mitbehoben und nur in Klammern genannt

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** Commit `e098ed10`, Botschaft · „(auch CVE-2026-97032 nicht)"; Slice-Plan slice-273 · „**Bezug:** Fund eines Fremd-Scanners"
- **befund:** Messung E führt CVE-2026-97032 (HTTP/2 HPACK, `MEDIUM`, Fix 0.60.0) im publizierten Image; Plan-Bezug und Titel nennen vier CVEs, die Botschaft erwähnt die fünfte nur als Klammer, ohne zu sagen, dass sie in `v0.86.0` steckte. Keine Überbehauptung — die Botschaft sagt eher weniger, als gemessen ist.
- **verifizierbar:** ja — Trivy-Vollbericht gegen `v0.86.0` und HEAD (Messung E/F)
- **klasse:** fund-menge-unvollstaendig-benannt

### F-2 — INFO · der Schweregrad hat sich seit der Plan-Messung bewegt

- **kategorie:** INFO
- **quelle:** ADR-0066 (Entscheidung 4); `BEO-ALL/scanner-vendor-severity-lag`
- **pfad:** Slice-Plan slice-273, §8 · „drei davon mit Schweregrad `UNKNOWN`"; Commit `e098ed10` · „Trivy fuehrt drei davon als UNKNOWN, make image-scan blieb gruen"
- **befund:** Am 2026-10-10 führt die Trivy-DB CVE-2026-78669 als `HIGH` mit Fix — nach Entscheidung 4 von ADR-0066 die Rot-Bedingung —, -78663 (beim Fremd-Scanner 9.1) als `MEDIUM`, nur noch zwei als `UNKNOWN`. Die Aussage in Plan und Botschaft ist eine datierte Messung und war zu ihrem Zeitpunkt nach Plan §8 richtig; für den Beleg der Beobachtung ist die Bewegung selbst der Inhalt (Verzug schließt sich teilweise, eine Einstufungs-Differenz CRITICAL↔MEDIUM bleibt). `make image-scan` selbst wurde hier **nicht** gefahren.
- **verifizierbar:** ja — Messung E; `make image-scan` gegen `v0.86.0`
- **klasse:** scanner-schweregrad-zeitabhaengig

### F-3 — INFO · das Risiko aus §6 hat einen messbaren Ausgang

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** Slice-Plan slice-273, §6 · „Die Hebung zieht weitere Module nach"
- **befund:** Messung C zeigt identische `go.mod`-Hashes von x/net v0.59.0 und v0.60.0; die Hebung kann deshalb keine weiteren Module nachziehen, und Messung A/B bestätigt das am Diff. Auch die §4-Rückführung (neuere Toolchain) ist nicht eingetreten: die `go`-Direktive bleibt 1.27.1.
- **verifizierbar:** ja — Messung C, `git show HEAD`
- **klasse:** risiko-ausgang-belegbar

## Negativbefunde

- **Minimalität und Vollständigkeit der Hebung** — geprüft, ohne Befund: nur die x/net-Zeile in `go.mod`, in `go.sum` nur die zwei x/net-Zeilen ersetzt, alte Zeilen entfernt, `make tidy` ist Fixpunkt (A, B, C); das Binary trägt v0.60.0 (D).
- **Abgrenzung §1** — geprüft, ohne Befund: go-git (v5.19.2) und go-winio (v0.6.2) unverändert, also nichts aus dem Dependabot-PR #7 mitgenommen; `tools/image-scan.sh`, das Makefile und die Beobachtung sind im Diff nicht berührt.
- **Docker/make-only (`AGENTS.md` §3.1)** — geprüft, ohne Befund: der Weg ist `go.mod`-Edit plus `make tidy` (Container `golang:$(GO_VERSION)`, `GOTOOLCHAIN=local`); kein `go get` nötig, weil die neue Version keine neuen Requirements bringt (C).
- **Regel 15 (Botschaft vs. Messung)** — geprüft, ohne Befund: „keine weiteren Module nachgezogen" (A, D), „keiner der vier mehr (auch CVE-2026-97032 nicht)" (F), „make gates gruen" (G — der Nachweis-Hash passt zu HEAD) sind gedeckt. „HTTP/2" ist für -78663, -78669, -97032 über die Trivy-Titel bestätigt; für -78659 und -78660 sind die Titel abgeschnitten, dort trägt die Aussage der Plan-Bezug.
- **ADR-0066/ADR-0067** — geprüft, ohne Befund: kein Gate gelockert, keine Schwelle verändert (`AGENTS.md` §3.6); die Handhebung neben dem Dependabot-Kanal widerspricht ADR-0067 nicht — der Kanal ist ein Weg, kein ausschließlicher.
- **Kommentare/Zustandsfelder (`AGENTS.md` §3.7)** — geprüft, ohne Befund: der Diff trägt keine Kommentare und keine Zustandsfelder.
- **Traceability** — geprüft, ohne Befund: die Botschaft nennt slice-273 und ADR-0066.

## Kategorie-Summary

HIGH 0 · MEDIUM 0 · LOW 0 · INFO 3.

## Verdikt

Freigabe. Die Hebung ist minimal, vollständig und über den Docker-/make-Weg
reproduzierbar; die Abgrenzung ist eingehalten, die Botschaft ist durch die Messungen
gedeckt. Die drei INFO-Punkte sind Material für die Closure-Notiz (fünfte CVE,
Beleg für die Beobachtung, Ausgang des §6-Risikos), kein Nacharbeits-Anlass am Diff.
