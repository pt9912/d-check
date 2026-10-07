# Review R2 — slice-257: CVE-Scan je Plattform des Image-Index

- **Review-Art:** Code. Geprüft gegen den Slice-Plan (`slice-257`, §1 Ziel und
  Abgrenzung, §3 samt Plan-Änderung nach R1, §6), gegen
  [ADR-0066](../plan/adr/0066-cve-scan-gegen-das-publizierte-image.md) und
  [ADR-0102](../plan/adr/0102-multi-arch-index-und-spiegel-per-index-digest.md)
  sowie die Hard Rules [`AGENTS.md`](../../AGENTS.md) §3.1, §3.5, §3.7 und §5
  Regel 13. Die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-257`, Runde 2, nur die Einarbeitungen seit R1 · Range
  `5c2624fc..e0d96368` (`c97c2a18` Plan-Änderung, `90377e04` R1-Einarbeitung,
  `2d6fdfd1` Verify-Report, `e0d96368` Verify-Befund V-2); Dateien
  `tools/image-scan.sh`, `harness/sensors/image-scan.md`, ADR-0066, ADR-0102
  (je ein Geschichte-Anhang), Slice-Plan. Mitgelesen, nicht im Diff:
  `.github/workflows/image-scan.yml` (Parser der Ausgabe),
  `tools/image-verify-published.sh` (Plattform-Prüfung im Release-Pfad).
- **Skill:** `reviewer.md` @ 1.17.0.
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-07
- **Eingangs-Kontext:** Slice-Plan `slice-257`; ADR-0066, ADR-0102; Hard Rules;
  vorherige Findings am selben Gegenstand: R1 `slice-257` (F-1 bis F-5) und
  Verify `slice-257` (V-1 verlangt genau diese Runde, V-2 Ursachen-Ausgabe,
  V-3 DoD-Wortlaut).
- **Proben (nur lesend; Netz für Registry und Trivy; eine Wegwerf-Registry
  `registry:2` auf `127.0.0.1:5055`, danach gestoppt):**
  - `bash tools/image-scan.sh --selftest`: 11 Proben `ok`, `== Fehlschlaege: 0`.
  - `imagetools inspect` mit dem Template aus `index_plattformen`:
    `ghcr.io/…:latest` und `pt9912/d-check:latest` → `linux/amd64`, `linux/arm64`,
    rc 0. `ghcr.io/…:v0.83.0` → `ERROR: template: … can't evaluate field
    Manifests`, rc 1, stdout leer. Nicht existierender Ref → `ERROR: …: not
    found`, rc 1, stdout leer.
  - Ursachen-Zeile (`{{.Manifest.MediaType}}`): v0.83.0 →
    `application/vnd.docker.distribution.manifest.v2+json`, fehlender Ref →
    `ERROR: …: not found`. Die Trennung, die die Sensor-Datei zusagt, ist für
    diese zwei Fälle gemessen.
  - Gemischter Lauf `IMAGE_SCAN_REFS="…:v0.83.0 …:latest"`: Kein-Index-Zeile
    und Ursachen-Zeile für v0.83.0, zwei `OK` für `latest`, dann
    `mindestens ein Scan ist GESCHEITERT`, rc 2.
  - **Teilweise Index-Antwort:** Index von `:latest` in die Wegwerf-Registry
    kopiert und als `dc:partial` mit drei Einträgen neu geschrieben: amd64,
    dann ein Eintrag **ohne** `platform` (OCI-zulässig, das Feld ist optional),
    dann arm64. `imagetools inspect` mit dem Skript-Template gibt
    `linux/amd64` auf stdout aus, dann `ERROR: … nil pointer evaluating
    *v1.Platform.OS`, rc 1. Die aus dem Skript extrahierte Funktion
    `index_plattformen` liefert `[linux/amd64 ]` mit rc 0. Der volle Skriptlauf
    (Docker-Shim im Scratchpad, der `docker run` nur `--network host` und
    `TRIVY_INSECURE` mitgibt, sonst unverändert durchreicht) endet mit
    `OK — … (linux/amd64)` und `image-scan: keine behebbaren CRITICAL/HIGH in:
    localhost:5055/dc:partial (linux/amd64)`, **rc 0**. Die arm64-Variante des
    Index wurde nicht gescannt, und der Lauf sagt das nicht.

## Findings

| # | Kategorie | Befund | Quelle | Pfad | verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| R2-1 | HIGH | `index_plattformen` verwirft den Exit-Code von `imagetools inspect` (`2>/dev/null`, `\|\| true` über die ganze Pipeline). Eine **teilweise** Ausgabe gilt deshalb als vollständige Plattformliste. buildx schreibt das Template zeilenweise und bricht erst am fehlerhaften Eintrag ab. **Failure-Szenario, gemessen:** Ein Index mit einem Eintrag ohne `platform` vor arm64 ergibt `[linux/amd64 ]`. Der Lauf scannt nur amd64, meldet „keine behebbaren CRITICAL/HIGH“ und endet mit Exit 0. Die Ausgabe enthält weder `GESCHEITERT` noch den Befund-Marker, also setzt der Workflow `status=clean`. Das ist genau die stille Plattform-Lücke, die der Slice schließen soll, jetzt über die Index-Ableitung. **Reichweite:** Heute lehnt `tools/image-verify-published.sh` (1) ein solches Index-Format im Release-Pfad ab, bevor getaggt wird. Der Scan stützt sich damit stillschweigend auf eben die Kopplung an den Release-Pfad, die F-1 auflösen sollte. | Skill Prüffrage 1 (Stilles Grün); [`AGENTS.md`](../../AGENTS.md) §5 Regel 13 | `tools/image-scan.sh` · „2>/dev/null \| grep -v '^unknown/' \| sed '/^$/d' \| tr '\n' ' ' \|\| true“ | ja, die Probe „Teilweise Index-Antwort“ oben (Wegwerf-Registry, rc 0 bei fehlendem arm64-Scan) | teilantwort-als-vollstaendig |
| R2-2 | MEDIUM | Drei Aussagen versprechen mehr, als der Code hält: der Funktionskommentar („leer, wenn der Ref kein lesbarer Index ist“), die Sensor-Datei Grenze 4 („Ohne lesbaren Index gibt es keinen Scan“) und der Geschichte-Anhang von ADR-0066 („ein Ref ohne lesbaren Index gilt als gescheitert“). Ein **teilweise** lesbarer Index ist der Fall aus R2-1 und ergibt einen Teil-Scan mit grünem Ausgang. Keine der drei Grenzen-Aussagen nennt ihn. **Failure-Szenario:** Wer nach Sensor-Datei oder ADR beurteilt, ob ein grüner Nachtlauf alle Plattformen des Index abdeckt, liest „ja oder Exit 2“ und rechnet mit keinem dritten Ausgang. | [`AGENTS.md`](../../AGENTS.md) §5 Regel 13; Skill §MEDIUM *Grenzen-Liste ohne ihre größte Lücke* | `harness/sensors/image-scan.md` · „**Ohne lesbaren Index gibt es keinen Scan.**“; `tools/image-scan.sh` · „leer, wenn der Ref kein lesbarer Index ist“ | nein, kein Gate. Probe wie R2-1 | grenze-gegen-beschreibung-statt-gegenstand |
| R2-3 | INFO | Die Ursachen-Zeile ist ein **zweiter**, unabhängiger `imagetools`-Aufruf. Ein vorübergehender Fehler beim ersten Aufruf kann beim zweiten verschwunden sein, oder der Ref enthält nur Attestations-Einträge. Dann steht neben „kein lesbarer Multi-Plattform-Index“ die Ursache `application/vnd.oci.image.index.v1+json`, die Diagnose widerspricht also der Zeile davor. Ausgang und Exit-Code bleiben richtig (Exit 2), betroffen ist nur die Diagnose. | Sensor-Datei Grenze 4 („trennt ein Einzel-Manifest von einem fehlenden Ref oder einem Netz-Fehler“) | `tools/image-scan.sh` · „why=\"$(docker buildx imagetools inspect“ | nein | diagnose-aus-zweitem-aufruf |

## R1-Befunde: sachlich geschlossen?

- **F-1 (MEDIUM), im Kern geschlossen, mit neuem Pfad.** Die feste Kopie ist
  weg, und die Plattformen kommen je Ref aus dem Index. Das R1-Szenario (eine
  dritte Plattform im Index wird nicht mitgezogen) ist für einen **vollständig**
  lesbaren Index beseitigt. Durch die Ableitung ist aber der Teil-Antwort-Pfad
  R2-1 entstanden.
- **F-2 (LOW), geschlossen.** Der Geschichte-Anhang von ADR-0066 nennt drei
  Läufe je Plattform, alle mit `--exit-code 0`. Das stimmt mit dem Wrapper
  `trivy()` überein. Zum Satz über den „lesbaren Index“ siehe R2-2.
- **F-3 (LOW), geschlossen.** Der Geschichte-Anhang von ADR-0102 hält fest,
  dass die Konsequenz *Offen* eingelöst ist.
- **F-4 (INFO), geschlossen.** In §6 steht jetzt „verdreifacht seine
  Trivy-Läufe je Image“. Nachgerechnet: vorher 2 Läufe je Image, jetzt 3 × 2
  Plattformen = 6, also das Dreifache. Die Aussage stimmt.
- **F-5 (INFO), geschlossen.** `want_arch` nimmt das zweite Segment, eine
  Variante führt deshalb nicht mehr zu einer falschen Abweichung. Die
  Reihenfolge-Annahme steht weiterhin als Zusage am Kopf von `arch_aus_json`.

## Negativbefunde

- **Stille Grün-Pfade (Prüffrage 1), außer R2-1 ohne Befund:**
  - Einzel-Manifest: Das Template fällt an `.Manifest.Manifests`, stdout bleibt
    leer, der Ref gilt als gescheitert, Exit 2 (gemessen).
  - Ein fehlender Ref, ein Netz- oder Auth-Fehler und ein fehlendes
    buildx-Plugin ergeben stdout leer und damit denselben Zweig.
  - Ein Index nur aus Attestations-Einträgen wird durch `grep -v '^unknown/'`
    leer und damit gescheitert.
  - `|| true` sitzt bei `index_plattformen` und `why` auf Diagnose- bzw.
    Ableitungs-Pfaden. Das Urteil fällt danach am Leer-Test bzw. an
    `errored=1`, außer im Teil-Fall R2-1.
  - Ein Leerraum-Wert in `IMAGE_SCAN_PLATFORMS` fällt auf die Index-Ableitung
    zurück. Der weggefallene Leer-Check öffnet damit keinen Null-Iterations-Pfad.
  - `errored` wird nie zurückgesetzt.
- **Übersteuerung `IMAGE_SCAN_PLATFORMS`, ohne Befund:** Sie scannt bewusst eine
  Teilmenge. Die Schlusszeile nennt jetzt nur tatsächlich gescannte Labels
  (`scanned`) statt der Eingabeliste, sie behauptet also nichts über nie
  besuchte Paare. Der Workflow setzt die Variable nicht.
- **Variante, ohne Befund:** `linux/arm/v7` → `want_arch=arm`, und das ist der
  Wert, den Trivy als Architektur meldet. Der Nachweis prüft die Architektur,
  nicht die Variante. So steht es in Grenze 3 („vergleicht die gemeldete
  Architektur“), die Grenze ist also benannt.
- **Workflow-Parser (`image-scan.yml`), ohne Befund außer R2-1:**
  - Die neuen Zeilen „kein lesbarer Multi-Plattform-Index“ und „imagetools
    meldet: …“ enthalten keinen der beiden Marker.
  - Jeder Kein-Index-Fall setzt `errored=1` und erreicht deshalb die
    Schlusszeile `GESCHEITERT`, was `status=error` ergibt (gemessen im
    gemischten Lauf).
  - Die geänderte Grün-Schlusszeile wertet niemand aus: `grep` über Skripte und
    Workflows findet keinen Leser.
  - Ein Registry-Fehlertext in der Ursachen-Zeile kann `GESCHEITERT` nicht
    vortäuschen, und selbst wenn, schlüge er in die sichere Richtung aus.
- **Kommentar-Klassen §3.7, ohne Befund:**
  - Kopfblock über `IMAGE_SCAN_PLATFORMS`: Zusage, Abgrenzung (keine Kopie der
    Release-Liste, mit Grund), Grenze (gemessen) und Kopplung zum Nachweis.
  - `index_plattformen`: Zusage.
  - Ursachen-Kommentar: Zusage.
  - `want_arch`-Zeile: Zusage.
  - Es gibt keine Review-Historie, keinen Befund-Marker und keine Slice-Nummer.
    „(bis v0.83.0)“ grenzt die betroffenen Images ab und erzählt keine Herkunft.
  - Ob die Aussagen auch stimmen, prüft R2-2.
- **ADR-Immutabilität §3.5, ohne Befund:** Beide ADRs erhalten nur eine neue
  Zeile in `## Geschichte`, der Kern ist unverändert.
- **Plan-Änderung vor Code (`AGENTS.md` §6 Schritt 4), ohne Befund:** `c97c2a18`
  liegt vor `90377e04`, und die Mitnahme (Index-Ableitung, ADR-Anhänge) ist im
  Plan vermerkt.
- **Hard Rule §3.1, ohne Befund:** Neu sind nur `docker buildx`, `grep`, `sed`,
  `tr` und `tail`, kein Host-Interpreter.
- **Prüffragen 3, 5, 12, 14 und 16** sind nicht anwendbar: kein Go-Code, Netz
  ist der deklarierte Zweck außerhalb von `gates`, kein Kern-Modul, kein
  Provenance-Marker, kein neues `Schärft:`-Feld.
- **Prüffrage 13 (Negativtest):** Die Index-Ableitung hat keine netzlose Probe.
  Das ist kein eigener Befund, denn der fehlende Fall ist der aus R2-1, und
  dort ist er gemessen.

## Kategorie-Summary

HIGH 1 · MEDIUM 1 · LOW 0 · INFO 1.
Wiederkehrende Klasse: Die Einarbeitung eines Befunds öffnet einen neuen Pfad
derselben Fehlerklasse (stilles Grün über einer Teilmenge). Das deckt sich mit
Verify V-1 (Fix-Commit trägt neue Mechanik außerhalb der Review-Range) und
bestätigt dessen Annahme, dass die R2 nötig war.

## Verdikt

**Nicht freigegeben, solange R2-1 offen ist.** Die Index-Ableitung macht
eine teilweise `imagetools`-Antwort zu einer vollständigen Plattformliste. Der
Lauf endet dann grün, obwohl eine Plattform des Index ungescannt bleibt
(gemessen, rc 0). Heute verhindert der Release-Pfad ein solches Index-Format.
Der Scan sollte nach F-1 aber gerade nicht von dieser Kopplung abhängen. R2-2
hängt an R2-1: Die drei Grenzen-Aussagen stimmen erst, wenn der Code ihnen
folgt oder sie den Teil-Fall nennen. R2-3 ist INFO. F-2 bis F-5 sind
geschlossen, F-1 im Kern.
