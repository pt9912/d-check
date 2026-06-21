# Review — slice-036 Implementierung (`--trace` Requirements Traceability Matrix)

## Kopf-Metadaten

- **Review-Art:** Code-Review (fertiger Commit gegen Plan/Spec/Anforderungen/
  Hard Rules — **kein Verifier**: DoD-Abhaken ist nicht Gegenstand). Gates
  wurden hier abweichend vom Skill-Default zur **Beleg-Führung** real
  ausgeführt (read-only, Docker), die Ergebnisse sind in den Findings/
  Negativbefunden zitiert — die Closure-Gate-Bestätigung bleibt dennoch
  Sache der getrennten Verifikation.
- **Datum:** 2026-06-21
- **Gegenstand:** Commit `9ef67fa` — neuer read-only-Modus `--trace`: leitet
  aus den kanonischen Quellen (`spec/lastenheft.md` → Anforderungen;
  `docs/plan/adr/` + `docs/plan/planning/` → Referenzen) eine Requirements
  Traceability Matrix ab und gibt sie auf stdout aus (Default Markdown-Tabelle,
  `--trace --json`/`--yaml` maschinenlesbar über den format-neutralen
  report-Adapter aus slice-031). Neu: `app/trace.go`
  (`BuildTraceMatrix`/`TraceRow`/`TraceMatrix`, `traceRequirements`,
  `isFullReqID`, `traceTitle`, `traceRefs`), `report.go`
  (`Trace`/`TraceJSON`/`TraceYAML`, `joinOrDash`, `traceCell`), `cli.go`
  (`--trace`-Flag, `comboError`-Auslagerung, `runTrace`, Dispatch), fünf
  Akzeptanztests `TestCLI036_Trace_*`, Lastenheft-CR `DC-FA-CLI-009` (0.21.0),
  Spezifikation §DC-FA-CLI-009.a, operations.md, CHANGELOG, Slice-Plan.
- **Reviewer:** Claude (Agent), Skill `.harness/skills/reviewer.md` v1.0.0.
- **Modell:** `claude-opus-4-8[1m]`.
- **Eingangs-Kontext:** Slice-Plan
  `docs/plan/planning/in-progress/slice-036-rtm-trace.md`; Anforderung
  `DC-FA-CLI-009` (Haupt — RTM, drei AK), `DC-FA-CLI-003` (Exit-Codes),
  `DC-FA-CLI-004` (Ausgabeformate, format-neutraler Reporter), `DC-QA-02`
  (Determinismus), `DC-QA-03` (read-only/netzlos); Spezifikation
  §DC-FA-CLI-009.a; Hard Rules `AGENTS.md` §3 (§3.4 Spec-Straten nie abwärts,
  §3.6 Gate-Lockerung nur per ADR), `MR-006` (Referenzrichtung). **Kein
  neuer ADR** — Anspruch: additiv, read-only, keine neuen Import-Kanten. Die
  DoD-*Abhakung* lag bewusst nicht als Prüfgegenstand vor.

## Findings

### HIGH

Keine.

### MEDIUM

Keine.

### LOW

#### LOW-1 — `traceTitle` strippt die führende Kennung nicht, wenn der Heading-ID in Backticks steht

- **Kategorie:** LOW
- **Quelle:** `DC-FA-CLI-009` (AK „eine Zeile … mit … Titel"); Maintainability
- **Pfad:** `internal/hexagon/core/app/trace.go:125` (`traceTitle`)
- **Befund:** `traceRequirements` gewinnt die ID via
  `strings.Trim(fields[0], "`.,:;")` (Backticks entfernt), übergibt
  `traceTitle` aber den *unveränderten* Heading-Klartext `plain` und die
  *getrimmte* ID. `traceTitle` macht `strings.TrimPrefix(plain, id)`; beginnt
  `plain` mit einem Backtick (`` `DC-FA-X-001` — Titel``), schlägt das
  Prefix-Matching fehl und der gesamte Heading-Text (inkl. ID und Backticks)
  landet als Titel. Adversarialer Beleg (`/tmp/t037/adv2`, Heading
  ``### `DC-FA-BAR-004` — Backtick-Kennung im Heading``): `--trace --json`
  liefert `"id":"DC-FA-BAR-004"`, aber
  `"title":"`DC-FA-BAR-004` — Backtick-Kennung im Heading"`. Keine Auswirkung
  auf d-checks eigene Daten (alle Lastenheft-Anforderungs-Headings tragen die
  ID bare am Zeilenanfang, geprüft via `grep -nE '^### `DC-' spec/lastenheft.md`
  → keine Treffer); fremde, präfix-agnostisch getracete Repos (z. B. a-check),
  die ihre Kennung im Heading per Backtick-Span auszeichnen, erhalten einen
  unsauberen Titel.
- **Verifizierbar:** ja — `docker run --rm -v "<repo mit Backtick-Heading>:/repo:ro"
  d-check:latest --trace --json` zeigt die ID dupliziert im `title`-Feld.

### INFO

#### INFO-1 — Referenzen auf NICHT im Lastenheft definierte Kennungen werden still verworfen (keine Reverse-Lücke)

- **Kategorie:** INFO
- **Quelle:** `DC-FA-CLI-009` (Out-of-Scope deckt diesen Fall nicht ab —
  dokumentationswürdige Annahme)
- **Pfad:** `internal/hexagon/core/app/trace.go:66` (`BuildTraceMatrix`
  iteriert nur über `order` aus dem Lastenheft)
- **Befund:** Die RTM iteriert ausschließlich über die im Lastenheft
  definierten Kennungen; `traceRefs` sammelt zwar *alle*
  Anforderungs-förmigen Vorkommen in ADR-/Slice-Dateien, doch ein Treffer
  ohne korrespondierende Lastenheft-Definition (dangling reference auf eine
  nicht existente Anforderung) erscheint nirgends. Adversarialer Beleg
  (`/tmp/t037/adv1`): ADR-0099 und slice-100 referenzieren beide
  `DC-FA-GHOST-777`, das nicht im Lastenheft steht — die JSON-RTM nennt
  ausschließlich `DC-FA-FOO-001` und `DC-QA-01`, `GHOST-777` fehlt
  vollständig. Eine RTM könnte diese Richtung (verwaiste Referenz) als
  Lücke ausweisen; die Anforderung verlangt es nicht (AK nur
  Anforderung→Refs, Out-of-Scope schweigt dazu), daher INFO statt Defekt.
  Bewusst undokumentierte Annahme.
- **Verifizierbar:** ja — ein Repo mit einer ADR/Slice-Referenz auf eine
  nicht im Lastenheft definierte Kennung; `--trace --json` enthält die
  Kennung nicht.

#### INFO-2 — Anforderungs-Erkennung ist heading-level-agnostisch

- **Kategorie:** INFO
- **Quelle:** Spezifikation §DC-FA-CLI-009.a Pkt. 1 („führende
  Heading-Kennungen", ohne Level-Einschränkung)
- **Pfad:** `internal/hexagon/core/app/trace.go:95` (`traceRequirements`
  iteriert über alle `rules.ExtractHeadings`)
- **Befund:** Jedes Heading beliebiger Tiefe (`#`…`######`), dessen erstes
  Feld eine vollständige Anforderungs-Kennung ist, gilt als
  Anforderungs-Definition. Adversarialer Beleg (`/tmp/t037/adv3`): das H6
  `###### DC-FA-DEEP-001 — H6 Anforderung` wird als Anforderung gezählt.
  d-checks Lastenheft führt Anforderungen durchgängig als `###` (H3), daher
  ohne praktische Wirkung; ein Fremd-Lastenheft mit Kennungs-förmigen
  Sub-Headings würde diese mitzählen. Spec-konform (keine Level-Bindung
  formuliert), daher INFO. Korrekt abgegrenzt sind dagegen mid-heading-IDs
  (`## Anforderung zu DC-FA-MID-001` — nicht gezählt) und `.a`-Varianten
  (`### DC-FA-CLI-009.a` — `isFullReqID` schließt sie über das exakte
  Voll-Match aus; beide im selben Lauf belegt).

#### INFO-3 — Waisen lassen den Exit-Code unberührt (Exit 0, advisory)

- **Kategorie:** INFO
- **Quelle:** `DC-FA-CLI-003`/`DC-FA-CLI-009` (AK Happy/Boundary = Exit 0;
  bewusste Advisory-Designnotiz)
- **Pfad:** `internal/adapter/driving/cli/cli.go:141` (`runTrace` gibt bei
  Erfolg stets `0` zurück)
- **Befund:** `--trace` liefert auch bei vorhandenen Waisen Exit 0
  (`/tmp/t037/adv1`, 1 Waise: `--trace` und `--trace --json` je `exit=0`).
  Das ist konsistent mit der Anforderung (AK nennen nur 0/2; `--trace` reiht
  sich explizit in die Advisory-Modi ein) — ein CI-Konsument, der Waisen als
  Gate-Fehlschlag erwartet, bekommt kein Signal über den Exit-Code, sondern
  müsste das `orphans`-Feld auswerten. Bewusste, in der Spec implizit
  getragene Designentscheidung; dokumentationswürdig.

## Negativbefunde (geprüft, ohne Befund)

- **Determinismus (DC-QA-02):** Kein Map-Leak in der Ausgabe. `order`
  (`trace.go:113`) und jede Referenzliste (`trace.go:167`) sind
  `sort.Strings`-sortiert; `Trace`/`TraceJSON`/`TraceYAML` iterieren nur über
  die geordnete `m.Requirements`-Slice (`report.go:248`,
  `encodeJSON`/`encodeYAML` serialisieren Structs in fester Feldreihenfolge).
  **Verifiziert:** je 5 Markdown- und 3 JSON-Läufe über das eigene Repo
  byte-identisch (gleicher `sha256sum`).
- **read-only / netzlos (DC-QA-03):** `BuildTraceMatrix` ruft ausschließlich
  `pathExists`/`fsys.ReadFile`/`rules.DiscoverFiles` (alle Lese-Pfade); kein
  Schreib-/git-/Netz-Zugriff, keine HTTP-Verdrahtung (der `external`-Pfad in
  `Run` wird vor `runTrace` gar nicht erreicht). `TestCLI036_Trace_Markdown`
  (`cli_acceptance_test.go:325`) prüft, dass keine `.d-check.yml` entsteht.
  **Verifiziert:** `make doc-check` ruft `--network none` + `:ro`-Mount,
  95 Dateien, 0 Befunde — der Trace-Pfad ist über dasselbe Image ohne
  Schreibrecht lauffähig.
- **Waisen-Logik korrekt (Waise = kein Slice; ADR-only zählt als Waise):**
  `trace.go:72` setzt `Orphan` allein bei leerer Slice-Liste, unabhängig von
  ADRs. **Verifiziert** (`/tmp/t037/adv2`): `DC-FA-BAR-002`/`-004` (nur ADR,
  kein Slice) → `orphan:true`; `DC-FA-BAR-003`/`-005` (nur Slice) →
  `orphan:false`. Deckt sich mit Lastenheft-AK und §DC-FA-CLI-009.a Pkt. 3.
- **Dublette einer Anforderung:** der `if _, seen := titles[tok]; !seen`-Guard
  (`trace.go:105`) behält die erste Definition. **Verifiziert**
  (`/tmp/t037/adv1`, `DC-FA-FOO-001` zweimal als Heading): genau eine Zeile,
  Titel der ersten Definition.
- **`isFullReqID`/`reqIDFull` erfasst nur Anforderungen:** das exakte
  Voll-Match (`reqIDFull.FindString(tok) == tok`, `trace.go:120`) schließt
  `.a`-Varianten aus (`DC-FA-CLI-009.a` → Match endet bei `…009`, ungleich
  Token). ADR-/MR-/slice-Headings werden nicht als Anforderung gezählt
  (Präfix-Form `…-FA-…`/`…-QA-…` greift bei `ADR-0099`/`MR-006`/`slice-100`
  nicht). **Verifiziert** (`/tmp/t037/adv1` mit ADR-/slice-Headings;
  `/tmp/t037/adv3`).
- **`traceTitle`-Trenner robust:** Em-Dash (—), Hyphen (-), Doppelpunkt (:)
  werden über `TrimLeft(rest, "—-:· ")` (`trace.go:127`) abgeräumt.
  **Verifiziert** (`/tmp/t037/adv2`): `Em-Dash-Titel`/`Hyphen-Titel`/
  `Doppelpunkt-Titel` sauber extrahiert (Backtick-Sonderfall: LOW-1).
- **`traceRefs` Owner-Ableitung + Skip:** `adrFileShape`/`sliceFileShape`
  (`trace.go:21`) leiten die Owner-Kennung aus dem Dateinamen ab
  (`0099-x.md`→`ADR-0099`, `slice-100-x.md`→`slice-100`); Dateien ohne
  Treffer (README.md, roadmap.md, fremde) werden übersprungen
  (`trace.go:147`). **Verifiziert:** Dogfooding-Lauf zieht Slices aus
  `done/` *und* `in-progress/` (slice-036), die `README.md`-Dateien in
  `adr/`/`planning/` erscheinen in keiner Owner-Zelle. Referenz-Treffer
  unabhängig von der Form (bare/Backtick/Markdown-Link/`.a`-Suffix), Beleg
  `/tmp/t037/adv2` (ADR-0050 referenziert via Span, Link und `BAR-004.a`).
- **Rendering — Pipe im Titel:** `traceCell` (`report.go:278`) ersetzt `|`
  durch `\|`. **Verifiziert** (`/tmp/t037/adv2`, Titel „Titel mit | Pipe
  darin"): Markdown-Zelle `Titel mit \| Pipe darin`, Tabelle bleibt intakt.
- **Rendering — leere Listen `[]` statt `null`:** `BuildTraceMatrix`
  normalisiert nil-ADR-/Slice-Listen auf `[]string{}` (`trace.go:69`/`:73`)
  und initialisiert `Requirements` via `make(…,0,len)` (`trace.go:66`).
  **Verifiziert:** YAML-Lauf (`/tmp/t037/adv1`) zeigt `adrs: []`/`slices: []`;
  Boundary ohne Lastenheft (`/tmp/t037/adv4`) liefert JSON
  `"requirements": []` (nicht `null`), `total:0`.
- **JSON/YAML strukturgleich:** beide serialisieren dasselbe
  `app.TraceMatrix`-Struct (`report.go:263`/`:267`) über `encodeJSON`/
  `encodeYAML` (slice-031); identische Feldnamen via json+yaml-Tags
  (`trace.go:27`–`:39`, camelCase-stabil). `TestCLI036_Trace_JSON`/`_YAML`
  prüfen dieselben Felder/Werte.
- **CLI-Kombi-Vertrag (Regression):** `comboError` (`cli.go:125`) deckt alle
  Alt-Kombinationen ab. **Verifiziert** (`/tmp/t037/adv2`): `--json --yaml`,
  `--repair --json`, `--doctor --repair` je Exit 2 mit unveränderter Meldung;
  neu `--trace --doctor`/`--trace --repair-broad` Exit 2
  („nicht mit --doctor/--repair"); `--trace --json --yaml` Exit 2 (json+yaml
  zuerst geprüft). `--trace --json`/`--yaml` erlaubt.
- **`--trace` nach dem Pfad-Argument (reorderArgs):** `--trace` ist ein
  bool-Flag, steht nicht in `valueFlags` (`cli.go:66`) und wird daher beim
  Reorder als Flag-Token ohne Wert-Verbrauch behandelt. **Verifiziert:** der
  Image-Aufruf injiziert via ENTRYPOINT `["/d-check","/repo"]` das Pfad-Arg
  *vor* `--trace` (`docker run … d-check:latest --trace` ⇒ args
  `[/repo, --trace]`) — der Dogfooding-Lauf funktioniert, also greift
  reorderArgs für die nachgestellte Option.
- **Dispatch-Position (`Run`):** `runTrace` (`cli.go:346`) steht nach
  `openRoot` (Wurzel validiert: existiert/Verzeichnis/nicht leer) und vor
  `loadConfig` — `--trace` braucht keine `.d-check.yml` und scannt nicht über
  die Konfiguration. `--print-config` (vor `openRoot`) bleibt der einzige
  Pre-Mount-Kurzschluss. Konsistent mit „read-only, kein Config-Bedarf".
- **Spec/Doku-Treue + MR-006:** Lastenheft §DC-FA-CLI-009
  (`spec/lastenheft.md:330`) und Spezifikation §DC-FA-CLI-009.a
  (`spec/spezifikation.md:288`) sind deckungsgleich mit dem Code (Quellen,
  Ableitungsschritte, Waisen-Definition, Format-Trias, Exit-Semantik). **Kein
  Abwärts-Verweis** auf ADR/Slice im bindenden Spec-Körper (geprüft:
  `grep -nE 'ADR-[0-9]|slice-[0-9]'` über beide Abschnitts-Körper ohne
  Treffer); die einzige `slice-036`-Nennung steht in der Lastenheft-Historie
  (§7), die der `matrix`-Default via `exclude-sections` ausnimmt
  (`.d-check.yml`). Die neuen Links (CHANGELOG, operations.md, Spec→Spec,
  Spec→Lastenheft) lösen auf — `make doc-check` (links/anchors/ids/matrix
  über 95 Dateien) meldet 0 Befunde.
- **Versions-Bump/Historie regelkonform:** Lastenheft `0.20.0`→`0.21.0` mit
  neuer Historie-Zeile (`spec/lastenheft.md:829`), CR-Charakter, korrektes
  Schema (`DC-FA-CLI-009`, CLI-Familie), drei AK (Happy/Boundary/Negative) +
  Out-of-Scope — deckt sich mit dem Anforderungs-Anlege-Prozess
  (`harness/conventions.md` §Anforderungs-Anlege-Prozess).
- **Kein ADR vertretbar (§3.6):** `make arch-check` grün — „Import-Regeln
  R1–R5 (ADR-0005) + R6 (ADR-0012) eingehalten". `trace.go` importiert nur
  `core/rules` und `port/driven` (R6-konform: app darf rules; rules/model
  unverändert), `report.go` keine neue externe Lib. Keine Gate-Lockerung,
  keine neue Import-Kante — additiv, der „kein ADR nötig"-Anspruch ist
  belegt.
- **Faktencheck Dogfooding:** `--trace` über das eigene Repo liefert 26
  Anforderungen, **0 Waisen** (der Commit-Text nannte 1 Waise =
  `DC-FA-CLI-009` selbst, „jetzt im slice-036-Bezug nachgezogen" — der
  Slice-Plan referenziert die Anforderung, daher zur Recht 0). Stichprobe:
  `DC-FA-CLI-004` → ADR-0009 + slice-002/003/025/026/029/031;
  `DC-FA-CLI-008` → ADR-0008 + slice-026/027 — beide gegen Lastenheft/ADR
  plausibel.
- **Keine Regression (`make test`):** `go test ./...` vollständig grün (alle
  Pakete `ok`); `make coverage-gate` grün (Schwelle 93 %). Der report-Adapter
  hat keine eigenen Unit-Tests; die neuen `Trace`/`joinOrDash`/`traceCell`
  sind nur über die CLI-Akzeptanztests abgedeckt — konsistent mit dem
  Bestand (das Paket trug auch vor slice-036 keine `_test.go`-Datei), daher
  kein Finding.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 0 | 1 | 3 |

## Verdikt

**Freigegeben — closure-fähig.** Keine HIGH/MEDIUM offen, keine blockierenden
Punkte. Der Kern des Slices ist tragfähig: korrekte, präfix-agnostische
Anforderungs-Ableitung mit exaktem Voll-Match (ADR-/MR-/slice-Headings und
`.a`-Varianten sauber ausgeschlossen), Dubletten-Schutz, robuste
Trenner-Behandlung, korrekte Waisen-Semantik (Waise = kein Slice; ADR-only =
Waise), deterministische Sortierung über alle Ausgaben (byte-identisch
verifiziert), strikt read-only/netzloser Lese-Pfad, strukturgleiche
JSON/YAML-Serialisierung mit `[]`-statt-`null` und Pipe-sicherem
Markdown-Rendering, vollständiger und regressionsfreier Kombi-Vertrag
(`comboError`-Auslagerung deckt alle Alt- und Neu-Kombinationen, Exit 2), und
ein additiver, ADR-loser Schnitt ohne neue Import-Kante (arch-check grün,
§3.6 gewahrt). Spec-Straten verweisen nicht abwärts (MR-006/§3.4 gewahrt),
doc-check und make test sind grün. LOW-1 (Titel-Strip versagt bei
Backtick-Heading-Kennung) ist ein latenter Robustheits-Defekt ohne Wirkung
auf d-checks eigene Daten, fällt aber im präfix-agnostischen Fremd-Repo-Einsatz
an; INFO-1..3 (keine Reverse-Lücke für dangling references;
heading-level-agnostische Erkennung; Waisen ohne Exit-Code-Signal) sind
dokumentationswürdige Annahmen bzw. bewusste Advisory-Designentscheidungen.
Keiner blockiert die Freigabe. Die Closure-Gate-Bestätigung obliegt der
getrennten Verifikation.
