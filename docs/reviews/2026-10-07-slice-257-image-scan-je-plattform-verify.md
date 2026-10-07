# Verifikation slice-257 — CVE-Scan je Plattform des Image-Index (DoD)

- **Rolle:** Verifier (Modul 11), Frage „Bauen wir es richtig?". Geprüft gegen §2 DoD, §1 Ziel und Abgrenzung und §3 Plan des Slice-Plans (samt beider vermerkter Plan-Änderungen), gegen [ADR-0066](../plan/adr/0066-cve-scan-gegen-das-publizierte-image.md) und [ADR-0102](../plan/adr/0102-multi-arch-index-und-spiegel-per-index-digest.md) (beide `Accepted`). Review-Entscheidungen waren nicht der Maßstab; die R1-Befunde wurden auf sachliche Schließung geprüft.
- **Gegenstand:** `fd49c9c7..90377e04`, fünf Commits: `e6595333` (Plan-Änderung Plattform-Nachweis), `63400720` (Skript + Sensor-Datei), `5c2624fc` (R1-Report), `c97c2a18` (Plan-Änderung nach R1), `90377e04` (R1-Einarbeitung F-1 bis F-5).
- **Sensor-Evidence:** Alles selbst gemessen, nichts aus dem Implementer-Bericht übernommen: `bash tools/image-scan.sh --selftest`, ein echter Scan mit Default-Refs, sechs Fehlerpfade, zwei Mutationen (bewusstes Brechen), `make gates`, `make adr-check`/`make trace-check RANGE=fd49c9c7..HEAD`, `make review-coverage`, `gh run list/view` für den Nachtlauf. Kein Push, kein Tag; Netz nur für Trivy und Registry-Lesezugriffe.
- **Modell-ID:** claude-opus-5-5 (Claude Agent SDK)
- **Datum:** 2026-10-07

---

## DoD-Prüfung je Punkt (§2)

### 1. `tools/image-scan.sh` scannt je Referenz `linux/amd64` und `linux/arm64` (Plattform im Bericht genannt); die Plattformliste steht an einer Stelle: **ERFÜLLT** (Form nach Plan-Änderung, siehe V-3)

- Echter Scan, Default-Refs (`ghcr.io/pt9912/d-check:latest pt9912/d-check:latest`), Exit 0. Vier Scans mit Plattform im Label: `OK — keine behebbaren CRITICAL/HIGH in ghcr.io/pt9912/d-check:latest (linux/amd64).`, dasselbe für `(linux/arm64)` und für beide Plattformen von `pt9912/d-check:latest`. Schlusszeile: `image-scan: keine behebbaren CRITICAL/HIGH in: ghcr.io/pt9912/d-check:latest (linux/amd64) ghcr.io/pt9912/d-check:latest (linux/arm64) pt9912/d-check:latest (linux/amd64) pt9912/d-check:latest (linux/arm64)`.
- Die Plattformen kommen aus dem Index selbst (`index_plattformen`, `docker buildx imagetools inspect`). Gegengelesen: der Index von `:latest` trägt genau `linux/amd64 sha256:bf4b0074…` und `linux/arm64 sha256:d8fcde19…`, keine Variante und keinen Attestations-Eintrag. Eine Kopie der Release-Liste gibt es im Skript nicht mehr. Die Einzelstelle ist damit der Index, das ist strenger als der DoD-Wortlaut.
- `--selftest`: 11 Proben `ok`, `== Fehlschlaege: 0`, rc 0 (sieben zur Zählung, vier zur Architektur, wie die Sensor-Datei sagt).

**Bewusstes Brechen.** Beide Mechanismen sind korrektheitskritisch: eine Plattform, die ungescannt bleibt, sieht aus wie eine saubere. Beide werden ohne den Fix aus dem richtigen Grund rot bzw. still grün:

| Probe | Lauf | Ergebnis |
|---|---|---|
| M1: Plattform-Nachweis entfernt (`if [ "${got_arch}" != "${want_arch}" ]` → `if false`, Kopie im Scratchpad) | `IMAGE_SCAN_REFS=ghcr.io/pt9912/d-check:v0.83.0 IMAGE_SCAN_PLATFORMS=linux/arm64` | **Exit 0**, `OK — keine behebbaren CRITICAL/HIGH in ghcr.io/pt9912/d-check:v0.83.0 (linux/arm64).`: ein amd64-Einzel-Manifest wird still als „arm64 gescannt" gemeldet. Mit dem Nachweis (Probe B unten) ist es Exit 2. Der Nachweis trägt die Zusage |
| M2: Vorzustand `fd49c9c7:tools/image-scan.sh` | gegen `ghcr.io/pt9912/d-check:latest` (Index) | Exit 0, ein einziger Scan ohne Plattform, `OK — … in ghcr.io/pt9912/d-check:latest.` Die arm64-Variante bleibt ungescannt, und der Lauf sagt es nicht. Das ist die Lücke aus §1 |

### 2. Sensor-Datei nennt Plattformen und Grenze; `make gates` grün; Nachtlauf manuell ausgelöst und grün: **ERFÜLLT**

- `harness/sensors/image-scan.md`: „Je Registry **und je Plattform** des Index", Plattformen aus dem Index, `IMAGE_SCAN_PLATFORMS` als Übersteuerung, drei Läufe je Plattform. Die Grenzen 3 (Plattform steht im Bild, nicht im Flag) und 4 (ohne lesbaren Index kein Scan) stehen da, ebenso die Selbsttest-Zählung. Gegen den Gegenstand geprüft: jede Aussage ist durch die Proben unten gedeckt.
- `make gates` selbst gefahren auf HEAD `90377e04`, sauberer Baum: Exit 0, Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `d-check: 978 Datei(en) geprüft, 0 Befund(e)`, `Ran 55 rules on 65 files: 0 findings.`
- Nachtlauf, gelesen mit `gh run list --workflow image-scan.yml`. Zwei `workflow_dispatch`-Läufe, beide `success`:
  - `37659622795` auf `63400720` (vor R1, noch mit fester Liste).
  - **`37660782225` auf `90377e04` = HEAD = `origin/main`**, 17:39:38–17:40:23Z (~45 s). Im Log: vier `OK — …`-Zeilen je Registry × Plattform, die Schlusszeile mit allen vier Labels und `image-scan: sauber.` Der Workflow-Parser (`GESCHEITERT` → error, `behebbare CRITICAL/HIGH-Befunde` → findings) ist unverändert, und die neuen Fehlermeldungen münden alle in die Schlusszeile `mindestens ein Scan ist GESCHEITERT`, die er trifft.

### 3. Review durchgeführt, Report unter `docs/reviews/`: **ERFÜLLT** (siehe V-1)

R1-Report `2026-10-07-slice-257-image-scan-je-plattform-r1.md` (MEDIUM 1, LOW 2, INFO 2; „nicht freigegeben, solange F-1 offen"); `make review-coverage` → `0 Befund(e)`. Einarbeitung in `90377e04`, Stelle für Stelle geprüft:

- **F-1 (MEDIUM), sachlich geschlossen und belegt.** `IMAGE_SCAN_PLATFORMS` ist per Default leer, und die Plattformen kommen je Ref aus dem Index (`unknown/*` ausgenommen). Gibt es keinen lesbaren Index, gilt der Ref als gescheitert (Probe A). Das Failure-Szenario aus R1 (eine dritte Plattform im Index wird nicht mitgescannt) kann nicht mehr eintreten: der Scan liest dieselbe Quelle, die `image-verify-published.sh` prüft.
- **F-2 (LOW), geschlossen.** ADR-0066 hat einen neuen `## Geschichte`-Anhang: Entscheidungen 3/4 gelten je Plattform, es sind drei Läufe, alle mit `--exit-code 0`, die Plattformen kommen aus dem Index. Der Kern ist unverändert; `make adr-check RANGE=fd49c9c7..HEAD` → `0 Befund(e)`.
- **F-3 (LOW), geschlossen.** ADR-0102 hat einen Geschichte-Anhang „Konsequenz *Offen* … ist eingelöst". Die Zeile steht unter `## Geschichte` (letzter Abschnitt, Z. 170), `adr-check` grün. Das Kopf-Feld `**Bezug:**` des Plans nennt ADR-0102 weiterhin nicht, nur die Plan-Änderung verlinkt sie. Das ist eine Kleinigkeit und kein eigener Befund.
- **F-4 (INFO), geschlossen.** Der Text in §6 heißt jetzt „verdreifacht seine Trivy-Läufe". Die Dauer ist gemessen: der Nachtlauf auf HEAD brauchte ~45 s bei `timeout-minutes: 30`. Das Risiko ist damit für die Closure belegbar.
- **F-5 (INFO), sachlich geschlossen, soweit gefordert.** `want_arch` nimmt jetzt das zweite Segment (`linux/arm64/v8` → `arm64`). Gemessen: `IMAGE_SCAN_PLATFORMS=linux/arm64/v8` gegen `:latest` (Index ohne Variante) → Trivy findet die Plattform nicht, `Scan von … (linux/arm64/v8) ist GESCHEITERT`, Exit 2, also laut. Die Reihenfolge-Annahme („die erste Fundstelle ist das Bild selbst") steht im Kommentar von `arch_aus_json` als Zusage. Dass sie am Trivy-Pin hängt, steht dort nicht ausdrücklich, aber der Ausfallweg ist laut. Bleibt ohne Befund.

### 4. Closure-Notiz, Register, Risiko-Ausgänge, Paarungen: **OFFEN, wie erwartet**, kein Befund

§6 steht auf `(offen)`, §7 trägt `—`; `make verify-closure-notes` greift erst in `done/`. Für die Closure belegt dieser Bericht zum einzigen Risiko: drei Läufe je Plattform, gemessene Nachtlauf-Dauer ~45 s gegen 30 min Timeout.

---

## Fehlerpfade, selbst gefahren

| Pfad | Lauf | Ergebnis |
|---|---|---|
| A: Einzel-Manifest ohne Übersteuerung | `IMAGE_SCAN_REFS=ghcr.io/pt9912/d-check:v0.83.0` | Exit 2; `…:v0.83.0: kein lesbarer Multi-Plattform-Index — die Plattformen sind UNBEKANNT, nichts gescannt.` + `mindestens ein Scan ist GESCHEITERT` |
| B: Einzel-Manifest, arm64 erzwungen | dasselbe + `IMAGE_SCAN_PLATFORMS=linux/arm64` | Exit 2; `Trivy meldet Architektur [amd64] statt arm64 — die Plattform ist NICHT gescannt.` |
| C: Plattform fehlt im Index | `:latest` + `IMAGE_SCAN_PLATFORMS=linux/s390x` | Exit 2; Trivy `FATAL … unable to find the specified image`, `Scan von … (linux/s390x) ist GESCHEITERT` |
| D: Ref existiert nicht | `IMAGE_SCAN_REFS=ghcr.io/pt9912/d-check:gibt-es-nicht-999` | Exit 2; `kein lesbarer Multi-Plattform-Index` (siehe V-2) |
| E: Ref existiert nicht, Plattform erzwungen | dasselbe + `IMAGE_SCAN_PLATFORMS=linux/amd64` | Exit 2; `Scan von … (linux/amd64) ist GESCHEITERT` |
| F: Refs nur Leerraum | `IMAGE_SCAN_REFS='  '` | Exit 2; `IMAGE_SCAN_REFS ist leer — nichts zu pruefen ist KEIN gruener Befundstand.` |
| V: Plattform mit Variante | `:latest` + `IMAGE_SCAN_PLATFORMS=linux/arm64/v8` | Exit 2, laut (F-5 oben) |

Kein Pfad endet still grün. `errored` wird nie zurückgesetzt, und die Schlusszeile „keine behebbaren …" ist nur bei `errored=0` und `findings=0` erreichbar.

---

## Abgrenzung (§1): gehalten

- Keine Schwelle und kein neues rotes Urteil auf Befunde: der Entscheidungslauf (`CRITICAL,HIGH --ignore-unfixed`) und `zaehle` sind unverändert, nur `--platform` und das Label sind neu.
- Index und Release-Pfad (`image-publish.sh`, `image-verify-published.sh`, `release.yml`) sind im Range unberührt. `.github/workflows/image-scan.yml` ist ebenfalls unberührt.
- Jede Datei im Diff steht in §3 oder in einer der beiden Plan-Änderungen (Sensor-Datei, Skript, zwei ADR-Geschichte-Anhänge, Plan, R1-Report). Die Plan-Änderung nach R1 (`c97c2a18`) liegt **vor** dem Fix-Commit (`90377e04`), das entspricht `AGENTS.md` §6 Schritt 4.
- Die Rückführung aus §4 ist nicht eingetreten: Trivy wählt die Plattform per Flag (gemessen, Index-Fall). Für den Einzel-Manifest-Fall, in dem das Flag nicht greift, trägt der Nachweis.
- `make trace-check RANGE=fd49c9c7..HEAD` → `5 Commit(s) — OK`, `0 Befund(e)`.

---

## Verifier-Befunde

### V-1: LOW

- **Kategorie:** LOW
- **Quelle:** `AGENTS.md` §6 Schritt 8 (kein Self-Review); DoD Punkt 3
- **Pfad:** `tools/image-scan.sh` · „index_plattformen() {"
- **Befund:** Die Index-Ableitung (`index_plattformen`, die Kein-Index-Verzweigung, der Wegfall des Default `linux/amd64 linux/arm64`) entstand erst in `90377e04` und liegt außerhalb der R1-Range `fd49c9c7..63400720`. Funktional ist sie hier bestätigt (grüner Lauf lokal und im Nachtlauf, Pfade A/D, Mutationen). Auf Kommentar-Klassen und Maintainability hat sie kein Review gelesen. Das ist dieselbe Klasse wie `slice-256` V-1 (Fix-Commit trägt neue Mechanik): zweites Auftreten in Folge.
- **Verifizierbar:** ja, eine R2 über `5c2624fc..90377e04`.
- **Erwartete Aktion:** R2 über den Fix-Commit, oder in der Closure-Notiz begründen, warum dieser Bericht sie ersetzt (Umfang: eine Funktion von vier Zeilen plus eine Verzweigung). Bei der Closure ins Beobachtungs-Register (Kandidat, zweites Auftreten).
- **Klasse:** Fix-Commit trägt neue Mechanik außerhalb der Review-Range

### V-2: INFO

- **Kategorie:** INFO
- **Quelle:** `AGENTS.md` §5 Regel 13 (Grenze gegen den Gegenstand)
- **Pfad:** `tools/image-scan.sh` · „2>/dev/null | grep -v '^unknown/'"
- **Befund:** `index_plattformen` verwirft stderr von `imagetools inspect`. Ein nicht existierender Ref, ein Netz- oder Auth-Fehler und ein Docker-Hub-Rate-Limit melden deshalb alle `kein lesbarer Multi-Plattform-Index` (gemessen: Pfad D). Das ist fail-closed (Exit 2) und für den Nachtlauf `status=error`, also richtig. Die Diagnose zeigt aber auf das Image-Format statt auf die Ursache. Außerdem ist `docker buildx` (Plugin) eine neue Vorbedingung auf dem Host bzw. Runner, die die Sensor-Datei nicht nennt. Fehlt das Plugin, fällt der Lauf laut mit derselben Meldung.
- **Verifizierbar:** ja (Pfad D).
- **Erwartete Aktion:** optional: die Meldung um „oder nicht erreichbar" ergänzen bzw. stderr bei leerem Ergebnis mit ausgeben; `buildx` in der Sensor-Datei unter Bindung nennen.
- **Klasse:** —

### V-3: INFO

- **Kategorie:** INFO
- **Quelle:** Slice-Plan §2 gegen §3 (Plan-Änderung nach R1)
- **Pfad:** `docs/plan/planning/in-progress/slice-257-image-scan-je-plattform.md` · „scannt je Referenz `linux/amd64` und"
- **Befund:** Der DoD-Wortlaut beschreibt noch die feste Liste („`linux/amd64` und `linux/arm64` … die Plattformliste steht an einer Stelle"). Die Plan-Änderung nach R1 ersetzt sie durch die Index-Ableitung. Inhaltlich ist der Punkt übererfüllt (die eine Stelle ist jetzt der Index selbst), wörtlich passt er nicht mehr zum Code.
- **Verifizierbar:** nein.
- **Erwartete Aktion:** beim Abhaken in der Closure den Haken mit dem Verweis auf die Plan-Änderung versehen oder den Wortlaut vor dem Haken angleichen.
- **Klasse:** DoD-Wortlaut hinter der Plan-Änderung

Negativ geprüft ohne Befund:
- Selbsttest 11/11; Sensor-Datei-Zählung (7 + 4) stimmt.
- Plattform-Nachweis: Mutation M1 zeigt, dass er die Zusage trägt.
- Vorzustand M2: die Lücke aus §1 bestand und ist geschlossen.
- Sieben Fehlerpfade (Tabelle), alle Exit 2 und laut.
- Workflow-Parser gegen die neuen Meldungen (über die Schlusszeile `GESCHEITERT`).
- Nachtlauf auf HEAD grün, mit vier Plattform-Labels im Log.
- ADR-0066/0102 nur Geschichte-Anhänge (`adr-check`).
- Traceability (`trace-check`), `review-coverage`, `make gates`.
- Abgrenzung (keine Schwelle, Release-Pfad und Workflow unberührt).

---

## Verdict

- **DoD 1 erfüllt**, mit Bruch-Beleg für Nachweis und Lücke; V-3 betrifft den Wortlaut.
- **DoD 2 erfüllt**: Sensor-Datei, `make gates` selbst grün, Nachtlauf `37660782225` auf HEAD `90377e04` grün.
- **DoD 3 erfüllt**; R1 F-1 bis F-5 sachlich geschlossen, F-1 mit Mess-Beleg. V-1 empfiehlt eine R2 über den Fix-Commit oder eine Begründung in der Closure.
- **DoD 4 erwartungsgemäß offen.**

Kein HIGH, kein MEDIUM. V-1 bitte annehmen oder begründen; V-2 und V-3 ohne Pflicht-Aktion.
