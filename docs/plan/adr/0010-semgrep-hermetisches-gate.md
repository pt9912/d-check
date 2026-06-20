# ADR-0010 — semgrep als hermetisches Gate mit lokal gecachtem, gepinntem Regelset

**Status:** Accepted
**Datum:** 2026-06-19
**Autor:** pt9912
**Bezug:** [ADR-0006](0006-lint-profil-solid.md) (golangci-lint als erste
statische Analyse — semgrep ergänzt sie sprachübergreifend),
[ADR-0001](0001-implementierungssprache.md),
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
**Schärft:** keine Spec-Stelle — Prozess-/Qualitäts-ADR (Kurs-Modul 13);
verbindlich für `tools/semgrep.sh` (Cache-Hol-Logik + Pin) und das Target
`make semgrep` als Bestandteil von `make gates`.

## Kontext

Ein ad-hoc-Lauf `semgrep scan --config auto` (Nutzer, 2026-06-19) lieferte
über das Registry-Regelset einen brauchbaren Befund-Satz (1 Treffer, ein
False Positive). `--config auto` taugt aber **nicht als Gate**: es lädt die
Regeln live über Netz (bricht die netzlose Hermetik der `gates` —
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Messmethode)
und die Registry-Regeln ändern sich über Zeit (gleicher Code heute grün,
morgen rot — bricht
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).
golangci-lint ([ADR-0006](0006-lint-profil-solid.md)) deckt Go ab, aber
nicht sprachübergreifende Muster (bash, Dockerfile, Secrets, YAML).

## Entscheidung

semgrep wird ein **reproduzierbares Gate mit netzlosem Scan**:

1. **Gepinntes, lokal gecachtes Regelset (nicht ins Repo vendort):**
   `semgrep/semgrep-rules` am festen Commit
   `d41fb34cf74466e2878af5f268ebf54466a04541` wird einmalig in einen Cache
   **außerhalb des Repos** geholt
   (`~/.cache/d-check/semgrep-rules/<commit>/`, Override
   `SEMGREP_RULES_CACHE`; wie `go mod`/Image-Pulls) und per `--config`
   lokal genutzt — statt `--config auto`. Kuratierter Umfang
   **`go/lang/security`** (55 Regeln, hoch-Signal für die Go-Codebasis
   dieses Repos). Breitere Packs (`bash/`, `dockerfile/`, `generic/secrets`)
   bewusst ausgelassen: ein Probelauf `go`+`bash`+`dockerfile` lieferte auf
   d-check **13 Treffer, sämtlich False Positives** (interne `FROM deps`
   Build-Stages ohne Digest-Pin; gewolltes `path` statt `filepath` für
   Slash-Pfade; numerisch-gewollte unquoted `$(…)` in `tools/*.sh`) — reines
   Rauschen, das „grün=Boden" untergräbt. `go/lang/security`: **0 Befunde**,
   daher kein `--exclude-rule` nötig.
2. **Lizenz/Provenienz:** die Regeln (Semgrep Rules License v1.0) werden
   **nicht ins MIT-Repo aufgenommen oder weitergegeben** ([ADR-0007](0007-repository-lizenz-mit.md)),
   nur lokal gecacht und genutzt (was die Lizenz erlaubt) — keine
   Fremdlizenz-Mischung im Repo; der Commit-Pin steht in `tools/semgrep.sh`.
3. **Netzlos + gepinntes Image:** der **Scan** läuft `docker run
   --network none --disable-version-check` (ohne den Flag liefe semgreps
   Versions-Ping unter `--network none` in einen ~2-min-Timeout) mit
   `semgrep/semgrep:1.167.0` und dem lokalen Regel-Cache
   ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
   Gepinnter Commit + gepinntes Image ⇒ identische Eingabe, identische
   Befunde
   ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).
   Das einmalige Holen des Caches am Pin ist Setup (Netz, wie Image-Pull),
   nicht Teil der Analyse.
4. **Im `make gates`** (per-Commit, Stop-Hook), zusätzlich zum
   eigenständigen `make semgrep`. `--error` ⇒ Befund bricht das Gate.
5. **Zentrale Ausnahmen mit Begründung** (keine Inline-Suppression,
   `AGENTS.md` §3.2): triffe künftig eine Regel nachweislich falsch, wird
   sie per `--exclude-rule` mit Why in `tools/semgrep.sh` ausgeschlossen.
   Beim gewählten Umfang `go/lang/security` aktuell **nicht nötig** (0
   Befunde).

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **Gepinntes, lokal gecachtes `go/lang/security`-Regelset, netzloser Scan, im Gate (gewählt)** | reproduzierbar (Commit-Pin); netzloser Scan; hoch-Signal, 0 FP-Rauschen; **keine** Fremdlizenz-/Vendor-Last im Repo | einmaliges Holen des Caches am Pin (Setup-Netz) nötig; nur Go-Abdeckung |
| Regeln dauerhaft ins Repo vendoren | offline ohne Cache-Hol-Schritt | Fremdlizenz-Mischung (Semgrep Rules License) im MIT-Repo + großer Diff + Wartungslast |
| `--config auto` im Gate | null Pflege, volle Abdeckung | Netz + nicht-reproduzierbar — verletzt Determinismus + Netzlosigkeit; flaky |
| Nur eigenständiges `make semgrep` (kein Gate) | einfach | Security-Regression fällt erst beim manuellen Lauf auf |
| Breiterer Umfang `go`+`bash`+`dockerfile` | sprachübergreifende Abdeckung | auf d-check 13 Treffer, **alle False Positives** — Gate-Rauschen + FP-Ausschluss-Last |
| Voll-Mirror des Registry-Sets | Abdeckung wie auto | großer Diff, Wartungslast, rauscht im Gate |

## Konsequenzen

- `make gates` enthält semgrep; ein neuer Security-/Static-Analysis-Befund
  bricht den per-Commit-Lauf. Hermetik bleibt (netzlos, gepinntes Image +
  gepinnte Regeln).
- **Update-Politik:** Regel-Commit-Pin und Image-Pin (beide in
  `tools/semgrep.sh`) werden **bewusst** gehoben (eigener Commit,
  Befund-Delta sichtbar), nicht automatisch — analog zur Pin-Politik der
  übrigen Toolchain (`GO_VERSION`, `GOLANGCI_LINT_VERSION`).
- Der Regel-Cache liegt **außerhalb des Repos** (kein getrackter Inhalt,
  keine Interferenz mit `go list ./...` oder dem d-check-Selbstscan).
- Neue False Positives → zentrale Ausnahme (`--exclude-rule`) mit Why in
  `tools/semgrep.sh` (keine Inline-Suppression).
- Go-Dependency-Fläche unberührt: semgrep läuft als externes Image, kein
  Go-Import (gomodguard/[ADR-0003](0003-config-format.md) unberührt).

## Fitness Function

`make semgrep` (Bestandteil von `make gates`): `docker run --network none`
mit gepinntem Image + gepinntem, lokal gecachtem Regelset, `--error`;
rot = Befund.

## Re-Evaluierungs-Trigger

- Wiederkehrende False Positives einer gecachten Regel → Regel
  entfernen/ausschließen (mit Why) oder Umfang neu kuratieren.
- Bedarf an weiteren Sprachen/Packs (bash/Dockerfile) → Cache-Umfang
  erweitern, sobald deren FP-Rauschen vertretbar/ausgeschlossen ist.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-19 | Proposed (slice-032) |
| 2026-06-20 | Accepted — `make gates` grün inkl. semgrep offline (0 Befunde, reproduzierbar); Review R1 (`docs/reviews/2026-06-20-slice-032-semgrep-gate.md`): HIGH-1 Regel-Nachweis + MEDIUM-1 Index-Wortlaut behoben |
