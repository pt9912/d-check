# Review R1 — slice-257: CVE-Scan je Plattform des Image-Index

- **Review-Art:** Code. Geprüft gegen den Slice-Plan (`slice-257`, §1 Ziel und
  Abgrenzung, §3 Plan samt Plan-Änderung „Plattform-Nachweis", §6 Risiken), gegen
  [ADR-0066](../plan/adr/0066-cve-scan-gegen-das-publizierte-image.md) (Entscheidungen
  3–6), [ADR-0102](../plan/adr/0102-multi-arch-index-und-spiegel-per-index-digest.md)
  (§Konsequenzen) und die Hard Rules [`AGENTS.md`](../../AGENTS.md) §3.1, §3.5, §3.7
  sowie §5 Regel 13/15. Die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-257` · Range `fd49c9c7..63400720` (`e6595333` Plan-Änderung,
  `63400720` Skript + Sensor-Datei); Dateien `tools/image-scan.sh`,
  `harness/sensors/image-scan.md`, Slice-Plan. Mitgelesen, nicht im Diff:
  `.github/workflows/image-scan.yml` (liest die Ausgabe des Skripts).
- **Skill:** `reviewer.md` @ 1.17.0.
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-07
- **Eingangs-Kontext:** Slice-Plan `slice-257`; ADR-0066, ADR-0102; Hard Rules.
  Vorherige Findings am selben Gegenstand: R1/Verify `slice-256` (F-5 —
  Spiegel-Satz „gleicher Config-Digest" in `tools/image-scan.sh`; im geprüften
  Stand nachgezogen zu „ab ADR-0102 auch derselbe Index-Digest").
- **Proben (nur lesend, Netz für Trivy):**
  - `bash tools/image-scan.sh --selftest` — 11 Proben ok, `Fehlschlaege: 0`, rc 0.
  - Echter Scan `IMAGE_SCAN_REFS=ghcr.io/pt9912/d-check:latest` — beide Plattformen
    gescannt, Nachweis bestanden, je `OK — keine behebbaren CRITICAL/HIGH in … (linux/amd64)`
    bzw. `(linux/arm64)`, Schlusszeile `… — je Plattform: linux/amd64 linux/arm64`, rc 0.
  - Einzel-Manifest-Image `ghcr.io/pt9912/d-check:v0.83.0` — amd64 `OK`, arm64:
    `Trivy meldet Architektur [amd64] statt arm64 — die Plattform ist NICHT gescannt.`,
    dann `mindestens ein Scan ist GESCHEITERT`, rc 2. Der Nachweis fängt den gemessenen
    Fall laut.
  - Fehlende Plattform `IMAGE_SCAN_PLATFORMS=linux/s390x` gegen `latest` — Trivy `FATAL`,
    `Scan von … (linux/s390x) ist GESCHEITERT`, rc 2.
  - Trivy-JSON `--platform linux/arm64` gegen `latest`: erste Fundstelle
    `"architecture": "arm64"` in Zeile 41 unter `Metadata.ImageConfig`; die
    Paket-Felder in `Results` heißen `"Arch"` (andere Schreibung, vom Muster nicht
    getroffen).

## Findings

| # | Kategorie | Befund | Quelle | Pfad | verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Plattformmenge ist eine fest verdrahtete Kopie und wird nicht aus dem Index gelesen. Der Kommentar sagt „JEDE Plattform des Index“, die Sensor-Datei „je Plattform des Index“. Dieselbe Liste steht unabhängig in `tools/image-publish.sh` (`--platform linux/amd64,linux/arm64`) und `tools/image-verify-published.sh` (Plattform-Gleichheit), und keine der drei Stellen nennt die anderen als Kopplung. **Failure-Szenario:** Ein Release nimmt eine dritte Plattform in den Index auf. Weil `image-verify-published.sh` laut fällt, wird er mitgezogen, `image-scan.sh` aber nicht. Der Nachtlauf meldet dann grün „je Plattform: linux/amd64 linux/arm64“, und die neue Variante bleibt ungescannt, genau die stille Lücke, die der Slice schließen sollte. Die Grenzen-Liste der Sensor-Datei nennt das nicht. Eskaliert von LOW (latente Wartungsfalle), weil es den Sicherheitspfad betrifft und das Szenario sich erzählen lässt. | [`AGENTS.md`](../../AGENTS.md) §3.7 (Kopplung), §5 Regel 13; Skill §MEDIUM *Grenzen-Liste ohne ihre größte Lücke* | `tools/image-scan.sh` · „JEDE Plattform des Index (ADR-0102)“; `harness/sensors/image-scan.md` · „Je Registry **und je Plattform** des Index“ | nein, kein Gate. Probe: `grep -rn 'linux/arm64' tools/` zeigt drei unabhängige Listen | zusage-ueber-kopie-statt-gegenstand |
| F-2 | LOW | Die Plan-Änderung listet als Spiegel von „BEIDE Trivy-Läufe“ nur den Skriptkopf und die Sensor-Datei. ADR-0066 führt denselben Satz als Entscheidung 3 („Beide Trivy-Läufe fahren `--exit-code 0`“) und als Entscheidung 4 („Zwei Läufe je Image“). Dazu ist weder ein `## Geschichte`-Anhang noch ein Vermerk entstanden. Wer die ADR liest, kennt den dritten Lauf, den Plattform-Nachweis, nicht. | [`MR-025`](../../harness/conventions.md#mr-025); [`AGENTS.md`](../../AGENTS.md) §1 (Widerspruch melden), §3.5 | ADR-0066 · „**Beide Trivy-Läufe fahren `--exit-code 0`.**“ | nein | spiegel-ausserhalb-der-liste |
| F-3 | LOW | ADR-0102 §Konsequenzen führt „**Offen:** `make image-scan` scannt bis zum Folge-Vorgang nur die Host-Variante des Index.“ Dieser Slice ist der Folge-Vorgang (Commit-Botschaft nennt ADR-0102). Der Diff hält das Erledigt-Sein aber nirgends fest, und das Kopf-Feld `**Bezug:**` des Plans nennt nur ADR-0066. Der Offen-Punkt bleibt als Zustand stehen, der nicht mehr stimmt. | [`AGENTS.md`](../../AGENTS.md) §1 (Widerspruch melden); ADR-0102 §Konsequenzen | ADR-0102 · „scannt bis zum Folge-Vorgang nur die“ | nein | offen-punkt-nicht-geschlossen |
| F-4 | INFO | Laut §6 verdoppelt sich die Scan-Zeit. Seit der Plan-Änderung laufen je Paar aus Registry und Plattform drei Trivy-Läufe. Bei 2 Registries × 2 Plattformen sind das 12 statt vorher 4 Läufe, also das Dreifache. Der Risiko-Text wurde nach der Plan-Änderung nicht nachgezogen. Der Ausgang ist Sache der Closure. `timeout-minutes: 30` des Workflows wurde nicht gemessen. | Slice-Plan §6 | Slice-Plan `slice-257` · „Der Nachtlauf verdoppelt seine Scan-Zeit.“ | ja, ein manuell ausgelöster Nachtlauf zeigt die Dauer | risiko-nach-plan-aenderung-veraltet |
| F-5 | INFO | `arch_aus_json` liest JSON mit `grep`/`sed`. ADR-0066 Entscheidung 6 hat das Template-Format gewählt, um gerade kein JSON zu parsen. Das ist kein Fremd-Interpreter, §3.1 bleibt eingehalten, und der Ausfallweg ist laut: ein umbenanntes Feld ergibt einen leeren Nachweis und damit Exit 2. Ungeschrieben bleibt die Annahme, dass `Metadata.ImageConfig.architecture` vor jedem anderen `"architecture"`-Schlüssel steht. Gemessen gilt das für 0.74.0, der Digest-Pin hält es. Außerdem ergibt eine Plattform mit Variante (`linux/arm64/v8`) über `${plat#*/}` den Wert `arm64/v8` und wird immer als Abweichung gemeldet. Das ist laut, nur beim Überschreiben relevant und nicht dokumentiert. | ADR-0066 Entscheidung 6 | `tools/image-scan.sh` · „die erste Fundstelle ist das Bild selbst“ | ja, `--selftest` (Muster) bzw. echter Scan (Reihenfolge) | undokumentierte-parser-annahme |

## Negativbefunde

- **Stille Grün-Pfade (Prüffrage 1), ohne Befund:**
  - Leere `IMAGE_SCAN_PLATFORMS` endet mit Exit 2 und einer Meldung.
  - Ein Trivy-Fehler im Nachweis-Lauf setzt `errored=1` und Exit 2.
  - Eine abweichende oder leere Architektur setzt `errored=1` und Exit 2. Gemessen an
    v0.83.0.
  - Eine fehlende Plattform ergibt Trivy `FATAL` und Exit 2. Gemessen mit `s390x`.
  - `continue` in der inneren Schleife springt zur nächsten Plattform derselben
    Referenz. Kein Paar wird übersprungen, ohne `errored` zu setzen, und `errored`
    wird nie zurückgesetzt.
  - Die Schlusszeile „keine behebbaren …“ ist nur bei `errored=0` und `findings=0`
    erreichbar.
- **Kopplung `--exit-code 0`, ohne Befund:** Der Wrapper `trivy()` setzt
  `--exit-code 0` für jeden Aufruf, also auch für den neuen JSON-Lauf. Die Aussage
  „ALLE Trivy-Laeufe“ in Skriptkopf und Sensor-Datei stimmt mit dem Code überein. Zum
  Spiegel in ADR-0066 siehe F-2.
- **Workflow-Parser (`image-scan.yml`), ohne Befund:**
  - Die Nachweis-Meldung selbst enthält kein `GESCHEITERT`. Die Schlusszeile
    „mindestens ein Scan ist GESCHEITERT“ steht aber bei jedem `errored=1` im Log, und
    der erste `grep` setzt `status=error`.
  - `behebbare CRITICAL/HIGH-Befunde` ist unverändert, nur das Label ist um die
    Plattform erweitert. `status=findings` greift also weiter.
  - Das JSON des Nachweises wird per `$(…)` gefangen und erreicht das Log nicht. Es
    kann also keinen Marker vortäuschen.
- **Kommentar-Klassen §3.7, ohne Befund:**
  - Der neue Block über `IMAGE_SCAN_PLATFORMS` trägt eine Zusage, eine Grenze
    (gemessen) und eine Kopplung („prueft der Plattform-Nachweis unten“).
  - `arch_aus_json` trägt Zusage und Grenze.
  - Die Schleifen-Kommentare tragen eine Kopplung.
  - Es gibt keine Review-Historie, keine Slice-Nummer und keinen Befund-Marker.
  - „(bis v0.83.0)“ grenzt die betroffenen Images ab und erzählt keine Herkunft.
- **Aussagen der Sensor-Datei gegen den Code (§5 Regel 13), ohne Befund außer F-1:**
  - „sieben Proben zur Zählung, vier zur Architektur“ entspricht der Selbsttest-Ausgabe
    (7 + 4).
  - „Exit 2“ bei Abweichung oder fehlender Plattform ist gemessen.
  - „laut, nicht still“ bei fehlendem Feld wird von der Probe „Feld fehlt“ → `[]` und
    vom Vergleich `"" != arm64` getragen.
- **Abgrenzung §1, ohne Befund:** Keine Schwelle, kein neues rotes Urteil auf Befunde,
  die Entscheidungslogik ist unverändert. Index und Release-Pfad sind nicht berührt.
- **Hard Rule §3.1, ohne Befund:** kein Host-Interpreter, nur `grep`, `sed`, `head` und
  `docker`.
- **Prüffragen 3, 5, 12, 14 und 16** sind nicht anwendbar: kein Go-Code, kein
  Kern-Modul, kein Provenance-Marker, kein neues `Schärft:`-Feld. Für Frage 5 ist Netz
  der deklarierte Zweck, außerhalb von `gates` (ADR-0066).
- **Prüffrage 13 (Negativtest), ohne Befund:** Die vier Architektur-Proben enthalten
  den Fall „Feld fehlt“. Die Netz-Pfade, also Abweichung und fehlende Plattform, sind
  nicht netzlos geprüft, hier aber gemessen (siehe Proben).

## Kategorie-Summary

HIGH 0 · MEDIUM 1 · LOW 2 · INFO 2.
Wiederkehrende Klasse: `spiegel-ausserhalb-der-liste`, zweites Auftreten am selben
Gegenstand nach R1 `slice-256` F-5.

## Verdikt

**Nicht freigegeben, solange F-1 offen ist.** F-1 ist MEDIUM mit erzählbarem
Failure-Szenario im Sicherheitspfad: Ein grüner Nachtlauf behauptet „je Plattform“,
gestützt auf eine Kopie, die niemand an die Quelle koppelt. F-2 und F-3 sind
Doku-Drift gegen höherrangige Quellen und vor der Closure zu behandeln, durch
Annahme oder Begründung. F-4 und F-5 sind INFO. Die Kernmechanik selbst, also der
Plattform-Nachweis, die Fehlerpfade und die Kompatibilität mit dem Workflow-Parser,
ist gemessen und trägt.
