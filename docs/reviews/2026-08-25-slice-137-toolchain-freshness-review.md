# Review-Report: slice-137 — Zwei Toolchain-Achsen im Nachtlauf — 2026-08-25

**Review-Art:** Code-Review (Diff gegen Slice-Plan, Hard Rules und die
zitierten `DC-*`-Verträge, Modul 10 §Drei Review-Arten) · **Gegenstand:**
Commit `26015b4` (`HEAD~1..HEAD`) — der Feature-Commit von slice-137
(`feat(harness): zwei Toolchain-Achsen im Nachtlauf — die andere Hälfte der
Blindheit (slice-137)`).

**Skill:** `.harness/skills/reviewer.md` @ 1.10.0 (Stand `9a7654a`) ·
**Modell-ID:** `claude-sonnet-5` · **Datum:** 2026-08-25

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan
  `docs/plan/planning/in-progress/slice-137-toolchain-freshness.md`
  (§1–§9, insb. §3 NICHT, §4 DoD, §5 Risiken)
- `tools/harness/pin-freshness.sh` (vollständig — Kopf-Kommentar, `compare()`,
  Dispatch, `--github`-/`--godev`-Zweige)
- `tools/harness/fetch-baseline-cache.sh` (Schwester-Achse, insb.
  `check_latest()` als Vergleichs-Präzedenz für Parse-Robustheit)
- `Makefile` (Targets `freshness-go`/`freshness-golangci`, `GO_VERSION`/
  `GOLANGCI_LINT_VERSION`, `.PHONY`)
- `.github/workflows/upstream-drift.yml` (drei Schritte, `if: always()`,
  `permissions`)
- `AGENTS.md` §3.1 (Docker/make-only, Host-Werkzeug-Klasse), §4
  (Gate-Tabelle)
- `harness/README.md` §Sensors
- `.d-check.yml` → Block `targets`

**Vom Reviewer selbst gefahren** (nur Lesekommandos im Repo — `git show`/
`git log`/`grep`/`sed -n` —, sowie **isolierte Reproduktionen** außerhalb
des Repos): `tools/harness/pin-freshness.sh` wurde unverändert in ein
temporäres Scratchpad-Verzeichnis kopiert
(`/tmp/claude-1000/…/scratchpad/pfr/`, danach entfernt) und dort gegen (a)
den netzlosen `--compare`-Einstieg mit den drei DoD-Fixtures sowie (b) ein
Fake-`curl`-Binary (`FAKE_MODE`-gesteuert, kein echtes Netz) für beide
Quellen-Zweige (`--github`, `--godev`) gefahren, um die Fail-open-/
Fail-safe-Logik empirisch zu prüfen. Zusätzlich wurde generisches
GNU-`make`-Verhalten (Exit-Normalisierung eines fehlschlagenden Recipes)
gegen ein Wegwerf-`TestMakefile` verifiziert. **Nicht** im geprüften Repo
ausgeführt: `bash tools/harness/pin-freshness.sh --compare …` direkt, kein
`make`-Target, keine echte Netzanfrage. `git status` vor/nach: clean, keine
Repo-Datei geändert.

**Verdikt: blockierend** — ein HIGH, kein MEDIUM, zwei LOW, ein INFO.

---

## Findings

### F-1 — `--godev` validiert die Form der Upstream-Antwort nicht; jede nicht-leere Nicht-Versions-Antwort (HTTP 200, Fehlerseite/Interstitial/Müll) wird als „upstream" behandelt und erzeugt ein Falsch-Positiv (`VERALTET`, Exit 3) statt des zugesagten `SKIP`

- **kategorie:** HIGH
- **quelle:** Reviewer-Skill §HIGH-Anker „Korrektheitsfehler in
  Kern-Modulen mit falschen Befunden/Exit-Codes"; das Skript bricht seine
  **eigene** Kopf-Zusage (`tools/harness/pin-freshness.sh:24-28`: „FAIL-OPEN
  … Netz-, Werkzeug- **oder Parse-Ausfall** meldet SKIP und Exit 0"); dieselbe
  Zusage steht wortgleich in `AGENTS.md` §4 (Zeile zu `make freshness-go` ·
  `make freshness-golangci`: „fail-open (Ausfall ⇒ `SKIP`, Exit 0)") und
  `harness/README.md` §Sensors („**Fail-open** mit Zeitgrenzen (`SKIP`, Exit
  0) — ein Sensor, der bei fremder Störung rot wird, wird abgeschaltet")
- **pfad:** `tools/harness/pin-freshness.sh:82-86` (`--godev`-Zweig: `curl …
  | head -1 | tr -d '\r' | sed 's/^go//'` — keine Form-Prüfung des
  Ergebnisses vor der Weitergabe an `compare()` in Zeile 93)
- **befund:** Empirisch reproduziert (Fake-`curl`, `FAKE_MODE=godev-html`
  liefert `<!DOCTYPE html>` mit Exit 0 statt Plaintext-Version): das Skript
  meldet `pin-freshness: go VERALTET — Pin 1.27.0, upstream <!DOCTYPE html>`,
  Exit 3. Der `--godev`-Zweig nimmt jeden nicht-leeren String als gültigen
  Upstream-Stand — anders als der `--github`-Zweig, dessen `case`-Muster
  `*/releases/tag/*` eine Struktur-Prüfung ist und Müll zuverlässig auf
  `upstream=""` (SKIP) abbildet (empirisch gegengeprüft, s. Negativbefunde),
  und anders als die Schwester-Achse `fetch-baseline-cache.sh:151`, deren
  `check_latest()` Kandidaten-Tags explizit gegen `grep -E
  '^v[0-9]+\.[0-9]+\.[0-9]+$'` filtert, bevor sie als Upstream-Stand gelten.
  Ein `SKIP`-Zweig existiert im `--godev`-Pfad nur für den leeren String
  (`[ -z "$upstream" ]` in `compare()`); ein **geparster, aber unsinniger**
  Wert (HTML-Fehlerseite, CDN-/WAF-Interstitial, Wartungsseite — alles mit
  HTTP 200 plausibel) durchläuft dieselbe Vergleichslogik wie eine echte
  Version und erzeugt einen falschen roten Nachtlauf. Das ist exakt das
  Szenario, vor dem der Skript-Kopf selbst warnt („ein Sensor, der bei
  fremder Störung rot wird, wird abgeschaltet — und ein abgeschalteter
  Wächter ist schlechter als ein löchriger"): eine fremde, transiente Störung
  (nicht Netz-, nicht Werkzeug-, sondern Content-Ausfall) wird hier **nicht**
  zu `SKIP`, sondern zu einem falschen `VERALTET`.
- **verifizierbar:** ja — `curl` durch ein Fake-Binary ersetzen, das für
  `https://go.dev/VERSION?m=text` beliebigen Nicht-Versions-Text mit Exit 0
  liefert; `NAME=go PINNED=<beliebig> bash pin-freshness.sh --godev` meldet
  `VERALTET`/Exit 3 statt `SKIP`/Exit 0.
- **klasse:** fail-open-luecke-fehlende-form-validierung-godev

### F-2 — Normalisierung ist einseitig: nur der Upstream-Wert wird um das `go`-Präfix bereinigt, der Pin-Wert nie — ein künftiger Pin-Format-Wechsel (`GO_VERSION=go1.27.0`) würde zu einem dauerhaften Falsch-`VERALTET` führen, ohne dass echter Drift vorliegt

- **kategorie:** LOW
- **quelle:** Maintainability (latente Wartungsfalle, Reviewer-Skill
  §LOW-Anker „hart verdrahteter Wert, der erst bei künftigem Edit zündet")
- **pfad:** `tools/harness/pin-freshness.sh:82-86` (`sed 's/^go//'` nur auf
  `upstream` angewandt); `compare()` in Zeile 40-53 (`pinned` geht
  unverändert in den Gleich/Ungleich-Vergleich); Kopf-Kommentar Zeile
  14-18 („Auf EIN Format bringen heisst hier: fuehrendes `go` strippen")
- **befund:** Die im Kopf beschriebene Normalisierung reduziert faktisch nur
  eine Seite des Vergleichs. Für die golangci-Achse ist das unschädlich
  (`v`-Präfix steht dokumentiert auf beiden Seiten und wird bewusst nicht
  gestrippt). Für die Go-Achse gilt die Symmetrie nur, solange
  `GO_VERSION` im Makefile ohne `go`-Präfix gepinnt bleibt (heute:
  `1.27.0` — korrekt). Würde der Pin künftig im go.dev-eigenen Format
  geschrieben (`GO_VERSION=go1.27.0`, ein naheliegender Kopier-Fehler, da
  genau dieses Format die Upstream-Quelle selbst verwendet), verglichen
  „go1.27.0" gegen den geschrumpften Upstream-Wert „1.27.0" — Dauer-Mismatch,
  Dauer-`VERALTET`, ohne dass ein echter Rückstand besteht.
- **verifizierbar:** ja — `bash pin-freshness.sh --compare go go1.27.0
  1.27.0` meldet `VERALTET` statt `ok`, obwohl beide Werte dieselbe Version
  bezeichnen.
- **klasse:** normalisierung-nur-eine-seite-des-vergleichs

### F-3 — Job-Anzeigename im Nachtlauf trägt weiterhin nur die Baseline-Achse, obwohl der Job jetzt drei Achsen prüft

- **kategorie:** LOW
- **quelle:** Reviewer-Skill §LOW-Anker „Doku-Drift (Prosa-Modullisten,
  veraltete Beispiele)"
- **pfad:** `.github/workflows/upstream-drift.yml:42-43` (`jobs:
  baseline-freshness: name: Baseline-Pin gegen upstream`)
- **befund:** Der Workflow-Kopf-Kommentar (Zeile 3-4) wurde korrekt auf
  „Baseline-Tag … Go-Toolchain und golangci-lint" erweitert, ebenso die
  Gate-Tabellen in `AGENTS.md`/`harness/README.md` (beide nennen explizit
  „Go-Version bzw. … `golangci-lint`-Release", nicht pauschal „die
  Toolchain" — BEO-011-Falle vermieden). Der **Job**-Anzeigename, der in der
  GitHub-Actions-Übersicht sichtbar ist, blieb dagegen unverändert bei
  „Baseline-Pin gegen upstream" und benennt die beiden neu hinzugekommenen
  Achsen nicht.
- **verifizierbar:** ja — Textvergleich `name:`-Feld gegen den erweiterten
  Kopf-Kommentar derselben Datei.
- **klasse:** doku-prosa-name-veraltet

### F-4 — Die Zusage „golang/go publiziert keine Release-Objekte" (Begründung für die Sonderquelle `--godev`) steht ohne Beleg im Skript-Kopf, im Slice-Plan und im Commit; für die Korrektheit des Skripts ist sie folgenlos, aber unbelegt

- **kategorie:** INFO
- **quelle:** Reviewer-Skill §INFO-Anker „dokumentationswürdige, aber
  undokumentierte Annahme"
- **pfad:** `tools/harness/pin-freshness.sh:9-12` (Kopf-Kommentar); ebenso
  `docs/plan/planning/in-progress/slice-137-toolchain-freshness.md:46-47`
  und die Commit-Botschaft
- **befund:** Die Aussage ist in diesem Review nicht netzlos nachprüfbar
  (keine GitHub-API-Abfrage erlaubt) und wird an drei Stellen wortgleich
  wiederholt, ohne dass eine der drei auf einen Beleg (Issue, Doku-Zitat)
  verweist — sie liest sich wie eine übernommene Tatsachenbehauptung
  (BEO-012-Nachbarschaft: Reichweite einer Quelle ungeprüft mitgeführt).
  Funktional folgenlos: Der `--godev`-Pfad funktioniert unabhängig davon, ob
  golang/go tatsächlich keine GitHub-Releases führt — die Wahl der
  Plaintext-Quelle wäre auch bei einer falschen Prämisse weiterhin korrekt.
- **verifizierbar:** teilweise — ein einzelner `curl -fsSLo /dev/null -w
  '%{url_effective}' https://github.com/golang/go/releases/latest`-Lauf
  (Netz) würde die Prämisse direkt bestätigen oder widerlegen; hier nicht
  ausgeführt (Netzverbot dieses Reviews).
- **klasse:** unbelegte-uebernommene-quellenaussage

## Negativbefunde

- geprüft, ohne Befund: `set -euo pipefail` in Kombination mit `$(… |
  … || true)` im `--godev`-Zweig — empirisch reproduziert (Fake-`curl`
  mit Exit 6, „DNS-Fehler"): die Pipe gibt unter `pipefail` den
  curl-Exit-Code an die Kette weiter, `|| true` fängt ihn korrekt ab,
  `upstream=""` bleibt leer → `SKIP`, Exit 0. Kein „rot statt SKIP" auf
  diesem Pfad.
- geprüft, ohne Befund: `--github`-Zweig bei Repo ohne Releases (Redirect
  auf `/releases` ohne `/tag/`) — Fake-`curl`-Reproduktion: `case`-Muster
  `*/releases/tag/*` greift nicht, `upstream=""` → `SKIP`, Exit 0
- geprüft, ohne Befund: `--github`-Zweig bei 404 (curl `-f` schlägt fehl,
  kein Stdout) — Fake-`curl`-Reproduktion: `eff=""` nach `|| true` →
  `SKIP`, Exit 0; auch bei einem Fake-`curl`, der trotz Fehler-Exit den
  ursprünglich angefragten `releases/latest`-URL auf Stdout ausgibt, matcht
  das `case`-Muster nicht (kein `/releases/tag/`-Substring) → weiterhin
  `SKIP`
- geprüft, ohne Befund: `--github`-Zweig bei Tag-Suffix
  (`v2.13.1-rc1`) — Fake-`curl`-Reproduktion: wird korrekt als
  String-Ungleichheit zum Pin erkannt (`VERALTET`); das ist die im
  Skript-Kopf explizit dokumentierte Entwurfsentscheidung
  („VERGLEICH IST GLEICH/UNGLEICH, kein Semver-Sort"), kein stiller Bug —
  und `releases/latest` zeigt laut GitHub-Produktverhalten ohnehin nie auf
  einen als Prerelease markierten Tag
- geprüft, ohne Befund: netzloser Vergleicher `--compare` — alle drei
  DoD-Fixtures reproduziert (gleich → `ok`/Exit 0, ungleich → `VERALTET`/
  Exit 3, leer → `SKIP`/Exit 0), deckungsgleich mit Commit-Botschaft und
  Slice-Plan §4
- geprüft, ohne Befund: Zeitgrenzen — beide Netz-Zweige tragen
  `--connect-timeout "$CT" --max-time "$MT"` (10s/60s); ein hängender
  (nicht fehlschlagender) TCP-Pfad wird von `curl` selbst beendet, nicht
  erst vom Job-`timeout-minutes: 10` — die in slice-133/F-5 gefundene
  Lücke (fehlendes Timeout in `fetch-baseline-cache.sh`) ist hier von
  Anfang an vermieden
- geprüft, ohne Befund: Exit-Code-Normalisierung durch `make` — generisch
  gegen ein Wegwerf-`TestMakefile` reproduziert (`exit 3` im Recipe → `make`
  selbst beendet mit Exit 2); die Doku-Aussage „Exit-Codes 0/3 sind die des
  Skripts; `make` normalisiert auf seinen eigenen" ist für
  `freshness-go`/`freshness-golangci` korrekt und vermeidet die in
  slice-133/F-3 gefundene falsche Exit-Code-Zuschreibung
- geprüft, ohne Befund: Nachtlauf, drei Schritte — Schritt 2 und 3 tragen
  `if: always()`, laufen also unabhängig vom Ausgang von Schritt 1 (und
  voneinander); jede der drei Achsen gibt in der Actions-Übersicht ihr
  eigenes Urteil ab, kein Abbruch verdeckt eine spätere Achse
- geprüft, ohne Befund: `permissions: {}` (Top-Level) + `permissions:
  contents: read` (Job-Ebene) — Job-Scope gilt für **alle** Steps
  desselben Jobs einschließlich der zwei neuen; keine fehlende Berechtigung
  für `freshness-go`/`freshness-golangci`
- geprüft, ohne Befund: keine neuen `uses:`-Einträge im Diff — §3.9
  (Action-SHA-Pinning) ist von diesem Commit nicht berührt
- geprüft, ohne Befund: `AGENTS.md` §3.1 Host-Werkzeug-Klasse — `curl`,
  `sed`, `head`, `tr` sind Fortsetzung des bereits etablierten Präzedenzfalls
  `fetch-baseline-cache.sh` (`curl`, `sed`, `unzip`, `sha256sum`), kein neu
  eingeführtes Werkzeug für diese Sensor-Familie
- geprüft, ohne Befund: `.PHONY` — `freshness-go freshness-golangci` beide
  eingetragen (`Makefile`-Diff), Namen decken sich exakt mit den
  Target-Definitionen
- geprüft, ohne Befund: `.d-check.yml` → `targets:`-Block — `doc-tables:
  [AGENTS.md, harness/README.md]` deckt beide geänderten Doku-Tabellen ab;
  beide neuen Targets sind in **beiden** Tabellen mit identischer
  Zeilen-Bindung (`make freshness-go` · `make freshness-golangci`)
  eingetragen — keine `BEO-010`-Drift zwischen den beiden Tabellen
- geprüft, ohne Befund: Slice-Plan §3 „Ausdrücklich NICHT" — kein
  Auto-Bump-Code, kein neues `gates`-Mitglied (`freshness-go`/
  `freshness-golangci` stehen nicht in der `gates:`-Kette des Makefiles),
  kein drittes Achsen-Ziel (Docker-Basis-Images/Action-Pins unberührt),
  kein neues d-check-Modul angelegt
- geprüft, ohne Befund: `ADVICE`-Text in beiden Targets — reine
  Handlungs-Hinweise ohne Befund-Nummern/Slice-IDs/Mess-Labels (Kommentar-
  Klassen-Regel), erscheint nur bei `VERALTET` auf stderr

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:**
`fail-open-luecke-fehlende-form-validierung-godev` ·
`normalisierung-nur-eine-seite-des-vergleichs` ·
`doku-prosa-name-veraltet` · `unbelegte-uebernommene-quellenaussage`

## Verdikt

**Merge-blockierend:** ja — ein HIGH offen. Die Kern-Konstruktion des Slice
(ein parametrierter Sensor, netzlos prüfbarer Vergleicher, zwei Targets im
bestehenden Nachtlauf statt einem zweiten, `if: always()` gegen
Achsen-Verdeckung, Doku in beiden Tabellen nachgezogen) trägt und ist an
zahlreichen Stellen nachweislich korrekt — insbesondere hat der Slice zwei
konkrete Lehren aus dem slice-133-Review (fehlende Curl-Timeouts, falsche
Exit-Code-Zuschreibung in der Doku) bereits vermieden, bevor sie hier erneut
auftreten konnten. F-1 ist der schwere Befund: Der `--godev`-Zweig — anders
als `--github` (Struktur-Prüfung via `case`-Muster) und anders als die
Schwester-Achse `fetch-baseline-cache.sh` (Regex-gefilterte Tag-Kandidaten)
— validiert die Form der Upstream-Antwort nicht, bevor sie mit dem Pin
verglichen wird. Ein nicht-leerer, aber unsinniger Antwortkörper (jede
Content-Störung, die keine Netz- oder HTTP-Fehler-Störung ist) erzeugt einen
falschen `VERALTET`-Befund statt des an drei Stellen (Skript-Kopf,
`AGENTS.md`, `harness/README.md`) zugesagten `SKIP` — genau das Szenario,
das der Skript-Kopf selbst als Anlass für einen abgeschalteten Wächter
benennt. F-2/F-3 sind Wartungsfallen bzw. Doku-Nachzug ohne akute
Fehlwirkung. F-4 ist eine unbelegte, aber für die Korrektheit folgenlose
Quellenaussage.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser
Report selbst ist ein Lauf-Beleg — DoD-/Spec-Konformität prüft der Verifier
separat.
