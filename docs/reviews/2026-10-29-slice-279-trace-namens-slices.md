# Review slice-279: Der Trace erkennt Slices mit Namens-Kennung

**Review-Art:** Code (Diff gegen Plan, ADR und Hard Rules; nicht gegen die DoD)
**Gegenstand:** slice-279, Range `cebdfb39..HEAD` (2724a57e feat, 16932f8c docs/spec)
**Skill:** `reviewer.md` @ 1.20.0
**Modell-ID:** claude-sonnet-5-5
**Datum:** 2026-10-29
**Eingangs-Kontext:** slice-279 (§1 Ziel und Abgrenzung, §3 Plan inkl. Negativliste, §6 Risiken), ADR-0110 (Proposed), welle-92, `DC-FA-CLI-009`, `DC-QA-02`, `AGENTS.md` §3.4/§3.7/§3.8.

## Findings

### F-1 — MEDIUM — Zusage „Ziffern-Repo byte-identisch" gilt nicht für Ziffern-Dateien unter `observations`

- **quelle:** Reviewer-Skill Anker „Botschaft verallgemeinert über die Messung hinaus" (Prüffrage 8), `DC-QA-02`, slice-279 §1 Ziel
- **pfad:** `spec/lastenheft.md` · „Default byte-identisch: … ein Repo, dessen Slice-Dateien der Ziffern-Form `slice-NNN-titel.md` folgen … byte-identisch zur Fassung vor dieser Erweiterung"; Gegenstück `internal/hexagon/core/app/trace.go` · „if underSegment(dir, file, observationsSegment) {"
- **befund:** Der `observations`-Ausschluss steht in der Default-Funktion vor dem Ziffern-Fall und gilt für beide. Eine Datei `observations/…/evidence/slice-007-foo.md` war bisher ein erkannter Slice (Ziffern-Fall) und ist es nun nicht mehr. Der Test `TestTraceSliceNameNegatives` pinnt das ausdrücklich (`slice-007-foo.md` unter `evidence/`). Das Kriterium, ADR-0110 §Konsequenzen und slice-279 §1 sagen dagegen, nur ein Repo mit Namens-Slices bekomme eine andere RTM.
- **Failure-Szenario:** Ein Adopter führt Register-Belege mit Titel (`evidence/slice-007-foo.md`), und eine Anforderung wird nur dort genannt. Bisher war sie `ok`; nach dem Upgrade ist sie `WAISE`, und `--trace --require-complete` kippt von Exit 0 auf 1, ohne `trace`-Block und ohne Namens-Slice im Repo. Der Eigenbestand zeigt das nicht, weil seine Belege titellos (`slice-143.md`) heißen und er `file-pattern` selbst setzt.
- **verifizierbar:** ja. Ein Test mit Ziffern-Beleg unter `observations/`, der den Alt-Stand grün und den Neu-Stand rot zeigt.
- **klasse:** `default-aenderung-ueber-benannten-fall-hinaus`

### F-2 — LOW — `underSegment` setzt eine bereinigte `slices.dir` voraus; die Negativprobe deckt nur diese Form

- **quelle:** `AGENTS.md` §3.8 (Eingaben, die das Modul nicht normalisiert); Plan §3 Negativliste „Datei im Verzeichnis `observations` einer konfigurierten fremden `slices.dir`"
- **pfad:** `internal/hexagon/core/app/trace.go` · „rel := strings.TrimPrefix(file, strings.TrimSuffix(dir, "/")+"/")"
- **befund:** `validateTracePath` lässt `./observations` oder `a//b` zu und lehnt nur `..` und führendes `/` ab. Bei einer nicht bereinigten `dir` trimmt `TrimPrefix` nichts. Dann steht das `observations`-Segment der `dir` selbst in `rel`, und das Kriterium „das Segment in `dir` selbst zählt nicht" gilt nicht mehr. `TestTraceSliceDirNamedObservations` fährt nur `SliceDir: "observations"`.
- **Failure-Szenario:** `trace.slices.dir: ./observations` ohne gesetztes `file-pattern`. Alle Slices fallen aus, und alle Anforderungen erscheinen als `WAISE` (sichtbar, nicht still grün). Das ist unwahrscheinlich und im Ergebnis laut.
- **verifizierbar:** ja, per Unit-Test mit `SliceDir: "./observations"`.
- **klasse:** `pfad-vergleich-ohne-normalisierung`

### F-3 — INFO — Fitness-Aussage der ADR nicht nachvollzogen

- **quelle:** ADR-0110 §Fitness-Funktion („je einzeln rot — jede Stelle gebrochen und gesehen")
- **pfad:** `docs/plan/adr/0110-trace-default-erkennt-namens-slices.md` · „Der Namens-Test wird ohne den Namens-Fall und ohne den `observations`-Ausschluss je einzeln rot"
- **befund:** Der Review hat die Mutationen nicht gefahren, weil er keinen Produktionscode ändern darf. Die Aussage ist ein Beleg-Anspruch für den Verifier (Rot-Beleg). Aus dem Code ist plausibel, dass beide Stellen von Tests erfasst werden.
- **verifizierbar:** ja, durch den Verifier (Bewusstes Brechen).
- **klasse:** `rot-beleg-ungeprueft`

## Beantwortung der Schwerpunkte

1. **Negativliste aus §3 vollständig umgesetzt und getestet?** Ja, jeder Punkt hat einen Test:
   - Beleg unter `observations`: `TestTraceSliceNameNegatives` und `TestCLI009_Trace_NamensSlice`.
   - `slice-273.md` ohne Titel: `TestTraceSliceDigitCaseUnchanged` (Kennung `slice-273`).
   - `slice-a.b.md`, `slice-.md`, `Slice-X.md`: `TestTraceSliceNameNegatives`.
   - Stub und Plan mit gleichem Basisnamen, sowie das Präfix-Paar `slice-a-foo` / `slice-a-foo-bar`: `TestTraceSliceNameIdentity`.
   - Fremde `slices.dir` namens `observations`: `TestTraceSliceDirNamedObservations`, mit F-2 als Lücke bei unbereinigter `dir`.
   - Gesetztes `file-pattern`: `TestTraceSliceFilePatternOverrides`.
   - Der Lauf „Planning-Baum-Kopie vorher/nachher" ist ein Closure-Kriterium und nicht Teil dieses Diffs. Er liegt nicht vor.
2. **Weitere Orte mit Slice-Namen (nachgezählt, Kommandos aus §3):**
   - `find docs/plan/planning -name 'slice-*.md' | wc -l` ergibt 549.
   - `find docs/plan/planning -name 'slice-*.md' ! -path '*/observations/*'` zusammen mit dem Ausschluss des Ziffern-Musters ergibt null Treffer. Alle 279 Nicht-Beleg-Dateien folgen `slice-NNN-….md`.
   - Alle 270 Rest-Dateien liegen unter `observations/BEO-*/…/evidence/`.
   - Außerhalb von `docs/plan/planning` gibt es keine `slice-*.md` (`docs/reviews` liegt außerhalb von `slices.dir`).
   - Ein zweiter Ort neben `observations` besteht im Eigenbestand nicht. Das Plan-Risiko 1 lässt sich damit nach Befund ausschließen.
3. **Spiegel untereinander und mit dem Code:** Lastenheft 0.107.0 (Kriterien und Historie), Spezifikation (Schritt 2, Schema-Zeile, Historie), Gerüst-Kommentar und ADR-0110 nennen dieselbe Erkennung wie der Code:
   - Ziffern-Fall zuerst.
   - Namens-Fall `[a-z0-9]+(-[a-z0-9]+)*`.
   - `observations`-Ausschluss unterhalb von `dir`.
   - Ein gesetztes Muster gilt unverändert.

   Abweichung nur bei F-1 (Byte-Identitäts-Zusage). Der ADR-Index trägt die Zeile mit Status Proposed, das Datum stimmt, die Bezüge sind aufwärts gerichtet. Die Spec-Straten nennen weder ADR noch Slice (§3.4 gehalten).
4. **Wirkung von `sliceFile *regexp.Regexp` → `ownerFunc`:**
   - ADR-Erkennung: `shapeOwner(adrFile, "ADR-")` ist semantisch identisch zum früheren Code (`path.Base`, Capture 1, Präfix).
   - Gesetztes `file-pattern`: `rt.sliceFile` ist nur bei gesetztem Muster nicht nil, denn `resolveTrace` hat den Default `sliceFileShape` nicht mehr im Feld. Der Pfad nutzt `shapeOwner(rt.sliceFile, rt.slicePrefix)` ohne Ausschluss und ohne Namens-Fall, also unverändert. Die Config-Schicht setzt `SliceFile` auch bei leerem Muster nicht (siehe Test), kein Verhaltensunterschied gefunden.
   - Einziger Aufrufer außerhalb von `trace.go` ist der neue Test.
5. **Eigenbestand:** `.d-check.yml` setzt `file-pattern: '^(slice-\d+)-.*\.md$'`. Damit läuft er über `shapeOwner` und damit unverändert, auch der `observations`-Ausschluss greift nicht. Die Planning-Baum-Kopie ohne `trace`-Block wurde nicht gemessen (kein Vorher-Image im Review-Kontext).

## Negativbefunde (geprüft, ohne Befund)

- **Hexagon-Richtung (ADR-0005):** `trace.go` importiert nichts Neues. Kein Befund.
- **Kommentarklassen (§3.7):** Neue Kommentare in `trace.go`, `config_template.go` und den Tests tragen Zusage, Abgrenzung oder Grenze mit `DC-FA-CLI-009` als Herkunft. Keine Slice-Nummern, Messlabels oder Review-Historie. Kein Befund.
- **Netz/Schreibzugriff (`DC-QA-03`):** nur lesende Pfade. Kein Befund.
- **§3.8 (Eingaben ohne Scan):** `observations` als nicht gezählte Eingabe ist in §1, Lastenheft, Spezifikation und ADR benannt. Die Grenze „Notiz-Dateien mit `slice-`-Präfix" steht als Grenze in der ADR und in Plan §1. Einzige Lücke: F-2.
- **ADR-0110:** Re-Evaluierungs-Trigger vorhanden, Alternativen-Tabelle vorhanden, Schärft-Feld mit Kennung (`DC-FA-CLI-009.a`) nach Adressierungs-Form. Provenance-Marker bei `slice-279` zeigt statt zu begründen. Kein Befund.
- **Quellen im Geltungsbereich:** `DC-FA-CLI-011` und `ADR-0108` werden als Bezug genannt und tragen die Aussage (Exit-Code beziehungsweise dieselbe Konfiguration). Kein Befund.
- **Eigene Erkennung statt Produktes (Frage 21):** Die Namens-Erkennung ist eine Dateinamen-Regel, keine Link-/Heading-Erkennung. Kein Befund.
- **Lint und Doc-Check:** `make lint` und `make doc-check` (1206 Dateien, 0 Befunde) liefen sauber durch, `make test` baute und lief. Das ist ein Beleg, kein Verifier-Urteil.
- **Doku-Drift (Handbuch):** `benutzerhandbuch.md` nennt weiter nur `slice-NNN-….md`. Laut Plan §1 und `AGENTS.md` §5 Zeile 17 ist das Release-Prep, kein Finding.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 1 | 1 | 1 |

## Verdikt

Nicht blockierend außer F-1. Das MEDIUM verlangt vor dem Merge entweder den `observations`-Ausschluss auf den Namens-Fall zu begrenzen oder Lastenheft, ADR und Plan so zu fassen, dass sie die Änderung für Ziffern-Belege unter `observations` benennen, samt Test. F-2 und F-3 sind optional beziehungsweise dem Verifier übergeben.
