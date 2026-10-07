# Review R3 — slice-257: CVE-Scan je Plattform des Image-Index

- **Review-Art:** Code. Geprüft gegen den Slice-Plan (`slice-257`, Plan-Änderung
  nach R2), gegen
  [ADR-0066](../plan/adr/0066-cve-scan-gegen-das-publizierte-image.md) und
  [ADR-0102](../plan/adr/0102-multi-arch-index-und-spiegel-per-index-digest.md)
  sowie die Hard Rules [`AGENTS.md`](../../AGENTS.md) §3.1, §3.5, §3.7 und §5
  Regel 13. Die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-257`, Runde 3, nur der Fix der R2-Befunde · Range
  `e0d96368..680428ab` (`3d77043e` R2-Report, `2d671dd5` Plan-Änderung,
  `680428ab` Fix); Dateien `tools/image-scan.sh`,
  `harness/sensors/image-scan.md`, ADR-0066 (ein Geschichte-Anhang),
  Slice-Plan.
- **Skill:** `reviewer.md` @ 1.17.0.
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-07
- **Eingangs-Kontext:** Slice-Plan `slice-257`; ADR-0066, ADR-0102; Hard Rules;
  vorherige Findings am selben Gegenstand: R1 (F-1 bis F-5), Verify (V-1 bis
  V-3), R2 (R2-1 HIGH, R2-2 MEDIUM, R2-3 INFO).
- **Proben (Netz für Registry und Trivy; eine Wegwerf-Registry `registry:2` auf
  `127.0.0.1:5103`, danach gestoppt, `--rm`):**
  - `bash tools/image-scan.sh --selftest`: 17 Proben `ok` (7 Zählung,
    4 Architektur, 6 Plattformliste), `== Fehlschlaege: 0`, rc 0.
  - Neues Template gegen die echten Refs: `ghcr.io/…:latest` und
    `pt9912/d-check:latest` → `linux/amd64`, `linux/arm64`, rc 0;
    `ghcr.io/…:v0.83.0` → `ERROR: template: … can't evaluate field Manifests`,
    rc 1.
  - Echter Scan (Default-Refs): vier `OK`-Zeilen (2 Refs × 2 Plattformen),
    Schlusszeile nennt alle vier Labels, rc 0.
  - Präparierte Indexe in der Wegwerf-Registry, `index_plattformen` und
    `plattformen_aus_antwort` unverändert aus dem Skript extrahiert, unter
    `set -uo pipefail`:

    | Tag | Inhalt | `idx_plats` | `idx_why` |
    |---|---|---|---|
    | `latest` | Kopie des echten Index | `linux/amd64 linux/arm64` | — |
    | `partial` | amd64, Eintrag ohne `platform`, arm64 (R2-Fall) | leer | `Index-Eintrag ohne Plattform` |
    | `lastnoplat` | amd64, arm64, Eintrag ohne `platform` | leer | `Index-Eintrag ohne Plattform` |
    | `empty` | `manifests: []` | leer | `Index ohne Plattform-Eintrag (nur Attestations)` |
    | `attest` | amd64, arm64, `unknown/unknown` | `linux/amd64 linux/arm64` | — |
    | `emptyos` | amd64, arm64 mit `os: ""` | `linux/amd64 /arm64` | — |
    | `nested` | amd64, verschachtelter Index mit `platform` arm64 | `linux/amd64 linux/arm64` | — |
    | `nestednoplat` | amd64, verschachtelter Index ohne `platform` | leer | `Index-Eintrag ohne Plattform` |

  - Voller Skriptlauf gegen `dc:partial` und `dc:empty`: je Kein-Index-Zeile,
    Ursachen-Zeile, `mindestens ein Scan ist GESCHEITERT`, rc 2. Der R2-Fall
    (rc 0 mit fehlendem arm64-Scan) ist damit nicht mehr herstellbar.
  - Stderr-Probe: `DOCKER_CONFIG` auf ein Verzeichnis mit fehlerhafter
    `config.json`. `imagetools` gibt rc 0 und schreibt zwei
    `WARNING: Error parsing config file …`-Zeilen auf stderr.
    `index_plattformen` liefert sie, wortweise, als Plattformen. Der volle Lauf
    gegen `ghcr.io/…:latest` scannt amd64 und arm64 korrekt, versucht dazu 17
    Schein-Plattformen (`WARNING:`, `Error`, `parsing`, …), meldet sie je
    zweimal `GESCHEITERT`, rc 2. Siehe R3-1.

## Findings

| # | Kategorie | Befund | Quelle | Pfad | verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| R3-1 | LOW | Mit `2>&1` landet stderr von `imagetools` in derselben Variable, aus der die Plattformliste entsteht. Bei rc 0 wird jede stderr-Zeile wortweise zur „Plattform“. **Gemessen:** Eine CLI-Warnung über eine fehlerhafte Docker-Config macht den Lauf rot (rc 2, `status=error`). Dazu kommen 17 Schein-Labels wie `(WARNING:)` und `(parsing)`, und die eigentliche Ursache steht in keiner Diagnose-Zeile. Das ist laut, nicht still, und fällt fail-closed aus: Die echten Plattformen werden trotzdem gescannt. Der Befund widerspricht aber der Zusage im Code, die Liste sei die „Plattformliste aus der Antwort von `imagetools inspect`“, und dem neuen Anhang an ADR-0066 („die Plattformen kommen aus dem Index selbst“). Entstanden ist er mit der R2-3-Einarbeitung (Ursache aus demselben Aufruf). | [`AGENTS.md`](../../AGENTS.md) §5 Regel 13 | `tools/image-scan.sh` · „2>&1)\" \|\| rc=$?“ | nein, kein Gate. Probe `DOCKER_CONFIG` mit fehlerhafter `config.json` | stderr-in-nutzdaten |
| R3-2 | INFO | Ein Index mit **null** Einträgen (`manifests: []`, rc 0, leere Antwort) wird richtig als gescheitert gewertet. Die Diagnose lautet dann aber „Index ohne Plattform-Eintrag (nur Attestations)“, obwohl es gar keine Einträge gibt. Ausgang und Exit-Code stimmen, nur der Klammertext ist ungenau. | Sensor-Datei Grenze 4 („trennt die Fälle voneinander“) | `tools/image-scan.sh` · „Index ohne Plattform-Eintrag (nur Attestations)“ | nein | diagnose-text-ungenau |

## R2-Befunde: sachlich geschlossen?

- **R2-1 (HIGH), geschlossen.** Drei Mechanismen tragen das Ergebnis, jeder
  einzeln gemessen:
  - **Exit-Code:** `local antwort rc=0` steht in einer eigenen Zeile, die
    Zuweisung `antwort="$(…)" \|\| rc=$?` folgt getrennt. `local` überdeckt den
    Exit-Code also nicht. In der Substitution läuft keine Pipe, `pipefail`
    spielt deshalb keine Rolle. Bei `set -u` sind alle Variablen gesetzt.
  - **Eintrag ohne Plattform:** `{{if .Platform}}…{{else}}?{{end}}` verhindert
    den nil-Pointer-Abbruch aus R2. Statt eines Teil-Outputs mit rc 1 kommt
    eine vollständige Ausgabe mit `?` und rc 0. Das `case *'?'*` leert dann die
    ganze Liste. Gemessen an drei Stellen: mittlerer Eintrag, letzter Eintrag,
    verschachtelter Index ohne `platform`.
  - **Abbruch mitten in der Liste:** rc ≠ 0 leert die Liste
    (Selbsttest-Probe „Abbruch nach erster Zeile“). Der echte Abbruch ist durch
    das `{{if}}` nicht mehr herstellbar, der Zweig bleibt aber als Netz.
  - **Attestations:** `unknown/unknown` wird ausgenommen, ohne die Liste zu
    leeren (gemessen, Tag `attest`).
  - Ein legitimes `?` in einer gültigen Antwort ist nicht möglich: Weder
    OS- noch Architektur- oder Varianten-Namen enthalten ein `?`. Sollte doch
    eines auftauchen, etwa über stderr, schlägt es in die sichere Richtung
    aus (leere Liste, Exit 2).
- **R2-2 (MEDIUM), geschlossen.**
  - Funktionskommentar: „leer = kein vollstaendig lesbarer Index“.
  - Sensor-Datei Grenze 4: nennt das Einzel-Manifest, den Eintrag ohne
    Plattform und den Abbruch, mit „eine Teil-Liste zählt nicht“.
  - ADR-0066-Anhang: „‚Lesbar‘ heißt vollständig lesbar“.
  - Alle drei stimmen mit dem gemessenen Verhalten überein, einschließlich des
    Satzes „gemessen gegen einen präparierten Index“ (hier nachgemessen).
  - Der ältere Kopfkommentar über `IMAGE_SCAN_PLATFORMS` („Ein Ref ohne
    lesbaren Index gilt als gescheitert“) ist unverändert. Er bleibt wahr, weil
    er keine Umkehrung behauptet, und ist deshalb kein Befund.
- **R2-3 (INFO), geschlossen.** `idx_why` stammt aus demselben Aufruf, der
  zweite `imagetools`-Aufruf ist weg. Der Preis ist R3-1.

## Negativbefunde

- **Stille Grün-Pfade (Prüffrage 1), ohne Befund:**
  - Leere Antwort mit rc 0 (Tag `empty`): leere Liste, Exit 2.
  - Nur Attestations: leere Liste, Exit 2 (Selbsttest).
  - Eine Plattform mit leerem OS (`/arm64`) wird nicht still übergangen. Sie
    bleibt in der Liste, und der Plattform-Nachweis prüft die Architektur
    (`want_arch=arm64`).
  - Verschachtelter Index mit `platform`: erscheint als Plattform. Der Scan
    fällt dann laut aus oder besteht den Architektur-Nachweis, still bleibt er
    nicht.
  - stderr-Rauschen (R3-1) wirkt nur in die rote Richtung.
- **Globale Variablen über Schleifen-Iterationen, ohne Befund:**
  - `idx_why` wird zu Beginn jeder Iteration neu gesetzt (`IMAGE_SCAN_PLATFORMS
    gesetzt`).
  - `idx_plats` wird nur gelesen, wenn `index_plattformen` in derselben
    Iteration lief, und die Funktion setzt es bei jedem Aufruf. Ein Wert aus
    einem früheren Ref kann deshalb nicht übrig bleiben.
  - `errored` wird nie zurückgesetzt. Ein gescheiterter erster Ref bleibt
    gescheitert, auch wenn ein späterer grün ist (in R2 gemessen, der Pfad ist
    unverändert).
- **Selbsttest-Zählung gegen die Sensor-Datei, ohne Befund:** „sieben Proben zur
  Zählung, vier zur Architektur, sechs zur Plattformliste“ = 7 `probe`, 4
  `arch_probe`, 6 `plat_probe`, gezählt im Skript und in der Ausgabe.
- **Grenze der Selbsttest-Deckung, benannt statt gemeldet:** Die
  `plat_probe`-Fälle prüfen `plattformen_aus_antwort`, nicht das
  `{{if .Platform}}`-Template. Ob buildx für einen Eintrag ohne Plattform
  wirklich `?` rendert, prüft nur die Registry-Probe. Dieselbe Klasse wie die
  benannte Trivy-Feldnamen-Grenze, hier mit Probe belegt.
- **Kommentar-Klassen §3.7, ohne Befund:**
  - `plattformen_aus_antwort`: Zusage und Grenze („ALLES oder NICHTS“, mit
    Grund).
  - `index_plattformen`: Zusage und Kopplung (`idx_plats`/`idx_why`, „aus
    DEMSELBEN Aufruf“).
  - Ursachen-Kommentar in der Schleife: Zusage.
  - Keine Review-Historie, keine Befund-Marker, keine Slice-Nummern.
- **ADR-Immutabilität §3.5, ohne Befund:** ADR-0066 erhält nur eine neue Zeile in
  `## Geschichte`, der Kern ist unberührt.
- **Plan-Änderung vor Code (`AGENTS.md` §6 Schritt 4), ohne Befund:** `2d671dd5`
  liegt vor `680428ab`. Die Mitnahme (Vollständigkeit, Ursache aus demselben
  Aufruf) ist im Plan vermerkt.
- **Hard Rule §3.1, ohne Befund:** Neu ist kein Werkzeug außer `docker buildx`,
  `grep`, `sed`, `tr` und `tail`.
- **Workflow-Parser, ohne Befund:** Die geänderten Diagnose-Zeilen („kein
  vollstaendig lesbarer …“, „imagetools: …“) enthalten keinen der beiden
  Marker. Jeder Kein-Index-Fall erreicht `GESCHEITERT`.
- **Prüffragen 3, 5, 12, 14 und 16** sind nicht anwendbar: kein Go-Code, Netz
  ist der deklarierte Zweck außerhalb von `gates`, kein Kern-Modul, kein
  Provenance-Marker, kein neues `Schärft:`-Feld.

## Kategorie-Summary

HIGH 0 · MEDIUM 0 · LOW 1 · INFO 1.
Wiederkehrende Klasse aus R2 („Einarbeitung öffnet einen neuen Pfad derselben
Fehlerklasse“): Sie tritt wieder auf, jetzt aber in der sicheren Richtung. R3-1
entstand mit der R2-3-Einarbeitung und macht den Lauf falsch **rot**, nicht
falsch grün.

## Verdikt

**Freigegeben.** R2-1 ist sachlich geschlossen: Eine Teil-Antwort und ein
Index-Eintrag ohne Plattform führen gemessen zu Exit 2. Einen neuen
Stilles-Grün-Pfad fand die Runde nicht. R2-2 und R2-3 sind geschlossen. R3-1
(LOW) und R3-2 (INFO) blockieren nicht. Ob sie angenommen oder begründet
abgelehnt werden, entscheidet der Implementer.
