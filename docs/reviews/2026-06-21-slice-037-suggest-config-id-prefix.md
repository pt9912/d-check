<!-- d-check:ignore (Review-Report: enthält illustrative IDs/Pfade als Belege; docs/reviews/** ist ids/exempt) -->

# Review — slice-037 Implementierung (`--suggest-config` Kennungs-Präfix-Option)

## Kopf-Metadaten

- **Review-Art:** Code-Review (Diff `e87fa77` gegen Plan/Spec/Anforderungen/
  Hard Rules — adversarial; kein Verifier. DoD-Abhaken und Gate-Lauf-
  Bestätigung sind nicht Gegenstand; Gates werden NICHT pauschal als grün
  angenommen, einzelne Belege sind separat ausgewiesen).
- **Datum:** 2026-06-21
- **Gegenstand:** Commit `e87fa77` „feat(suggest): slice-037 --suggest-config
  Kennungs-Praefix-Option (DC-FA-CLI-006, ADR-0015)". Neues Flag `--id-prefix`,
  Präfix-Ableitung aus dem Lastenheft (`deriveReqPrefix`), Platzhalter
  `<PREFIX>` + TODO statt fixem `DC-`. Geänderte Dateien:
  `internal/hexagon/core/app/suggest.go`,
  `internal/adapter/driving/cli/cli.go`,
  `internal/adapter/driving/cli/cli_acceptance_test.go`, `spec/lastenheft.md`
  (v0.20.0), `docs/plan/adr/0015-suggest-config-id-prefix.md`,
  `docs/plan/adr/README.md`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`,
  `docs/plan/planning/in-progress/slice-037-suggest-config-id-prefix.md`.
- **Reviewer:** Claude (Agent), Skill `.harness/skills/reviewer.md` v1.0.0.
- **Modell:** `claude-opus-4-8[1m]`.
- **Eingangs-Kontext:** Slice-Plan
  `docs/plan/planning/in-progress/slice-037-suggest-config-id-prefix.md`;
  Anforderung `DC-FA-CLI-006` (CR 0.20.0, AK), `DC-FA-ID-001` (emittiertes
  `ids`-Muster), `DC-QA-02` (Determinismus), `DC-QA-03` (read-only);
  Spezifikation §`DC-FA-CLI-006.a`; ADR-0015 (Proposed); Konvention
  `harness/conventions.md` (MR-006, Anforderungs-Anlege-Prozess), `AGENTS.md`
  §3.5/§3.6; Referenz `.d-check.yml`. Die DoD-Abhakung des Slices lag bewusst
  nicht als Prüf-Artefakt vor — die Findings sind eigenständig gebildet.
- **Verifikations-Setup:** adversariale Läufe gegen das lokale Image
  `d-check:latest` über read-only-gemountete `/tmp`-Test-Repos; volle
  Test-Suite per `make test` (Regress-Check); `make doc-check` (Selbst-
  konformität der Markdown-Edits). Das d-check-Repo wurde nicht verändert.

## Findings

### HIGH

Keine.

### MEDIUM

#### MEDIUM-1 — Messmethoden-Spec `DC-FA-CLI-006.a` nicht nachgezogen; ADR-0015 „keine Spec-Stelle" unzutreffend

- **Kategorie:** MEDIUM
- **Quelle:** `DC-FA-CLI-006` (Spec-Treue einer Messmethode); MR-001/
  Anforderungs-Anlege-Prozess (ADR schärft die **Spezifikation**)
- **Pfad:** `spec/spezifikation.md:163` (Vorlage hardcodet
  `DC-(FA-[A-Z]+|QA)-\d+`); `spec/spezifikation.md:97` (Extraktions-
  Schritte ohne Präfix-Ableitung); `docs/plan/adr/0015-suggest-config-id-prefix.md:14`
  (`Schärft: keine Spec-Stelle`)
- **Befund:** Die Spezifikations-Schicht `DC-FA-CLI-006.a` ist die
  Messmethoden-Spec für `--suggest-config`; ihre „Kanonische Vorlage"
  zeigt das Anforderungs-Muster weiterhin als fixes `DC-(FA-[A-Z]+|QA)-\d+`
  und ihre Schritte 1–6 kennen weder das Flag `--id-prefix`, noch die
  Präfix-Ableitung, noch den `<PREFIX>`-Platzhalter. Der Commit ließ
  `spec/spezifikation.md` unverändert (`git show --stat e87fa77`: Datei
  nicht enthalten; `grep -niE 'id-prefix|<PREFIX>' spec/spezifikation.md` →
  keine Treffer). Damit beschreibt die Messmethoden-Spec ein Verhalten, das
  der Code nicht mehr zeigt: real emittiert der Default `<PREFIX>-(FA-…`
  plus TODO bzw. das abgeleitete Präfix. Zugleich behauptet ADR-0015
  `Schärft: keine Spec-Stelle`, obwohl §`DC-FA-CLI-006.a` exakt diese
  Mechanik dokumentiert — die Spec-Stelle existiert und wurde stehen
  gelassen, statt geschärft zu werden (Konvention: ADRs schärfen die
  Spezifikation, `harness/conventions.md`).
- **Verifizierbar:** ja — ein Spec-Diff von §`DC-FA-CLI-006.a` gegen die
  reale Render-Ausgabe (`--suggest-config ai-harness-init` → `<PREFIX>` +
  TODO; `git show --stat e87fa77` zeigt keinen `spezifikation.md`-Edit)
  belegt die Drift; ein Leser von §`DC-FA-CLI-006.a` erhält die alte
  `DC-`-Vorlage als verbindliche Messmethode.

#### MEDIUM-2 — Override- und Konflikt-Bypass-Pfad von `ai-harness` + `--id-prefix` ungetestet

- **Kategorie:** MEDIUM
- **Quelle:** `DC-FA-CLI-006` (fehlende Negativtests/Contract-Test bei neuem
  öffentlichen Vertrag); ADR-0015 §Konsequenzen
- **Pfad:** `internal/adapter/driving/cli/cli_acceptance_test.go:1095`
  (alle fünf `TestCLI037_*`); `internal/hexagon/core/app/suggest.go:91`
  (`reqPrefix := idPrefix; if reqPrefix == "" && harness { … }`)
- **Befund:** `--id-prefix` wird in den Tests nur mit `ai-harness-init`
  kombiniert (`:1098` Happy `AC`, `:1163` ungültig). Kein Test deckt
  `ai-harness` (repo-bewusst) **zusammen mit** `--id-prefix` ab. Damit sind
  zwei vertraglich relevante Pfade ungeprüft: (a) explizites `--id-prefix`
  **überschreibt** die Lastenheft-Ableitung (Code: `idPrefix` gewinnt vor
  `deriveReqPrefix`); (b) ein **mehrdeutiges** Lastenheft, das ohne Flag
  Exit 2 erzeugt, wird mit `--id-prefix` **stillschweigend übergangen**
  (Exit 0). Beide Pfade verhalten sich in adversarialen Läufen wie erwartet
  (Override: `--suggest-config ai-harness --id-prefix ZZ` gegen DC-Lastenheft
  ⇒ `regex: 'ZZ-(FA-…'`, Exit 0; Konflikt-Bypass: konfligierendes Lastenheft
  + `--id-prefix ZZ` ⇒ `ZZ-…`, Exit 0), aber eine Regression in der
  Ableitungs-/Konflikt-Verzweigung würde von der Test-Suite nicht gefangen.
- **Verifizierbar:** ja — ein `make test` mit einem zusätzlichen Fall
  (`ai-harness` + `--id-prefix` über ein DC-/konfligierendes Lastenheft)
  würde den Pfad binden; sein Fehlen ist im Test-Bestand sichtbar
  (`grep 'ai-harness".*id-prefix\|id-prefix".*ai-harness'` findet nur die
  `-init`-Fälle).

### LOW

#### LOW-1 — Aktiver `<PREFIX>`-Platzhalter ergibt bei direkter Übernahme stilles Grün für die Anforderungs-Linkpflicht

- **Kategorie:** LOW
- **Quelle:** `DC-FA-ID-001` (das emittierte `ids`-Muster); Maintainability
  (latente Footgun)
- **Pfad:** `internal/hexagon/core/app/suggest.go:323`
  (`harnessIDPatterns`: bei leerem `reqPrefix` → `reqPrefix = "<PREFIX>"`,
  Muster bleibt **aktiv**, nur ein `# TODO`-Kommentar davor);
  `suggest.go:400` (aktiver Zweig schreibt `regex: '<PREFIX>-…'`)
- **Befund:** Ohne Präfix/Ableitung wird das Anforderungs-Muster als
  **aktives** `regex: '<PREFIX>-(FA-[A-Z]+|QA)-\d+'` gerendert (nicht
  auskommentiert; davor ein TODO). Der Regex compiliert (`<`/`>` sind
  literale Zeichen) und matcht keine reale Kennung. Übernimmt ein Nutzer die
  `ai-harness-init`-Vorlage, baut die Struktur auf, ersetzt aber `<PREFIX>`
  nicht, läuft d-check grün, ohne die Anforderungs-Linkpflicht tatsächlich zu
  prüfen. Reproduziert: dasselbe Repo mit einer un­verlinkten `DC-FA-X-001`
  in einer Prosa-Datei unter dem ids-`scope` liefert mit korrektem
  `DC-`-Muster 1 Befund (Exit 1,
  `id-unlinked`), mit dem `<PREFIX>`-Muster 0 Befunde (Exit 0). Die Stufe
  bleibt LOW, nicht HIGH: betroffen ist die **advisory** Vorschlags-Ausgabe
  (kein d-check-eigenes Gate), der TODO-Kommentar steht unmittelbar davor,
  und es ist die ausdrücklich gewählte Anti-Footgun-Variante (ADR-0015);
  das Versagen erfordert, dass der Nutzer das markierte TODO ignoriert.
- **Verifizierbar:** ja — `--suggest-config ai-harness-init > .d-check.yml`
  ohne `<PREFIX>`-Ersatz, dann `d-check` gegen ein Repo mit unverlinkter
  Anforderungs-Kennung: 0 Befunde / Exit 0 (gegen 1 Befund / Exit 1 mit
  ersetztem Präfix) belegt das stille Grün.

### INFO

#### INFO-1 — `--id-prefix` ohne `ai-harness`-Modus wirkungslos (still ignoriert), bei Ungültigkeit dennoch Exit 2

- **Kategorie:** INFO
- **Quelle:** Maintainability (undokumentierte Verhaltensvariante)
- **Pfad:** `internal/adapter/driving/cli/cli.go:299` (Validierung im
  `suggestConfig`-Zweig, vor dem Modus-Switch); `suggest.go:87`
  (`reqPrefix` nur im `initMode || harness`-Zweig genutzt)
- **Befund:** `--id-prefix` wirkt nur in den `ai-harness`-Modi. Mit nur
  echten Quellen (`--suggest-config spec/lastenheft.md --id-prefix ZZ`) wird
  der Wert **still ignoriert** (Ausgabe behält die abgeleiteten Muster, hier
  `DC-FA-X-\d+`), ohne Hinweis; ein *ungültiger* Wert (`--id-prefix zz`)
  führt im selben Fall dennoch zu Exit 2, weil die Validierung modus-
  unabhängig im `suggestConfig`-Block läuft. Ohne `--suggest-config`
  (normaler Scan) wird `--id-prefix` komplett ignoriert (keine Validierung,
  kein Hinweis). Eindeutiges, stabiles Verhalten — kein Korrektheitsdefekt,
  aber eine in Spec/AK/Handbuch nicht benannte Variante (ein Nutzer, der
  `--id-prefix` an eine echte Quelle hängt, erwartet evtl. eine Wirkung).
- **Verifizierbar:** ja — `--suggest-config spec/lastenheft.md --id-prefix ZZ`
  (Wert ohne Effekt, Exit 0) vs. `--id-prefix zz` (Exit 2) zeigt die
  Asymmetrie; ein Spec-/Handbuch-Diff belegt das Fehlen der Regel.

## Negativbefunde (geprüft, ohne Befund)

- **Ableitungs-Korrektheit (`deriveReqPrefix`/`reqShape`):** `reqShape`
  (`suggest.go:23`) fängt das Projekt-Präfix korrekt:
  `DC-FA-CLI-001`→`DC`, `DC-QA-02`→`DC` (adversarial verifiziert, beide ⇒ ein
  Muster `DC-(FA-…`); ein Suffix-Buchstabe (`DC-FA-X-001a`) und ein
  ziffernhaltiges Präfix (`D2-FA-X-001`→`D2`) werden konsistent behandelt.
  Nicht-Anforderungs-IDs werden ignoriert: ein Lastenheft mit nur
  `ADR-0001`, `MR-000`, `slice-001` ergibt **kein** abgeleitetes Präfix →
  Platzhalter, Exit 0 (verifiziert). Der gemischte Fall `DC-FA-X-001` +
  `DC-QA-02` = **ein** Präfix (kein Konflikt, verifiziert); `AC-FA-…` +
  `DC-QA-…` = Konflikt → Exit 2 (verifiziert). Die Token-Säuberung
  (`strings.Trim(fields[0], "\`.,:;")`, `:195`) und `StripHeadingLinks`
  (`:191`) sind dieselben wie im etablierten `extractDefinedIDs`-Pfad.
- **Konflikt-Determinismus (DC-QA-02):** die Fehlermeldung sortiert die
  Präfixe (`sort.Strings(ps)`, `suggest.go:214`); deklariert in der
  Reihenfolge `ZZ, AA, MM` ⇒ Meldung `(AA, MM, ZZ)` (verifiziert). Die
  Konflikt-Erkennung greift wirklich (Exit 2, kein Gerüst — verifiziert über
  `TestCLI037_IDPrefix_KonfliktFehler` und adversarial).
- **Determinismus des neuen Pfads (DC-QA-02):** `--suggest-config ai-harness`
  über ein abgeleitetes Lastenheft 10× byte-identisch (eine sha256-Zeile);
  `--suggest-config ai-harness-init` (Platzhalter-Pfad) 10× byte-identisch.
  Die einzige Map (`prefixes` in `deriveReqPrefix`) wird nur zum Sammeln/
  Zählen genutzt; die Ausgabe läuft über `sort.Strings` bzw. das Single-
  Element-Extrakt — keine Map-Iteration in der gerenderten Reihenfolge.
- **Platzhalter-Dekodierbarkeit (DC-FA-CONF-001):** das aktive
  `<PREFIX>-(FA-[A-Z]+|QA)-\d+`-Muster dekodiert über den eigenen Parser
  (`TestCLI037_IDPrefix_Platzhalter` ruft `configyaml.Decode`; adversarial:
  das vollständige Platzhalter-Gerüst lädt als `.d-check.yml` ohne Decode-
  Fehler — der Lauf scheitert erst an fehlendem ids-`target docs/plan/adr/`,
  also am erwarteten `ai-harness-init`-Zielbild-Verhalten, nicht am Muster).
- **Override (explizit gewinnt):** `--id-prefix` hat Vorrang vor der
  Ableitung — `ai-harness --id-prefix ZZ` gegen ein `DC`-Lastenheft ⇒
  `ZZ-(FA-…` (verifiziert; Code `suggest.go:91`, `idPrefix` wird nur bei
  leerem Wert durch `deriveReqPrefix` ersetzt). (Untestet — siehe MEDIUM-2.)
- **CLI-Validierung/`ValidIDPrefix`:** `idPrefixShape` `^[A-Z][A-Z0-9]*$`
  (`suggest.go:27`) lehnt Kleinbuchstaben (`ac`→Exit 2, verifiziert) und
  Leerwert (leer = Platzhalter/Ableitung, kein Fehler) konsistent ab; die
  Gestalt deckt sich mit `reqShape`s Präfix-Capture (`[A-Z][A-Z0-9]*`), ist
  also weder strenger noch laxer als das, was die Ableitung erzeugen kann.
- **`reorderArgs`/`valueFlags`:** `--id-prefix`/`-id-prefix` sind als
  wertnehmende Flags eingetragen (`cli.go:69`); ein Flag **nach** dem Pfad
  (`… ai-harness-init AC` als Tail-Wert) wird korrekt vorgezogen, ein
  wertloses Flag am Ende ist ein Nutzungsfehler (vorhandener
  `flag needs an argument`-Guard, `cli.go:77`). Die Option ist in
  `options` (`cli.go:56`), `parseOptions` (`:133`, `:158`) sauber verdrahtet.
- **read-only (DC-QA-03):** der Ableitungs-Pfad liest `spec/lastenheft.md`
  über `fsys.ReadFile`/`pathExists` (`suggest.go:182-185`) — kein Schreiben,
  kein git, kein Netz; nach `--suggest-config ai-harness` über ein
  gemountetes Repo entsteht kein `.d-check.yml` (verifiziert,
  `ls .d-check.yml` → nicht vorhanden). `SuggestConfig` schreibt nur in den
  `strings.Builder`.
- **ADR-0005-Imports:** `suggest.go` ergänzt keine neuen Importe (weiterhin
  `fmt`, `regexp`, `sort`, `strings`, der `rules`-/`model`-Kern und der Port
  `…/port/driven`) — kein `os`/`io/fs`/`net`, kein `internal/adapter/*`. Der
  Kern bleibt I/O-frei; `make test` (Suite grün) schließt den separaten
  `arch-check` zwar nicht ein, der Import-Bestand ist aber unverändert.
- **Geltungsbereich (nur Anforderungs-Muster parametrisiert):** in
  `harnessIDPatterns` trägt ausschließlich das `spec/lastenheft.md`-Muster
  `reqPrefix` (`suggest.go:332`); `ADR-\d{4}`, `MR-\d{3}`, `slice-\d{3}`,
  `CO-\d{3}` bleiben konventions-fest (verifiziert: Output zeigt unverändert
  `ADR-\d{4}` etc.). Das deckt sich mit ADR-0015 §Entscheidung 4 und der
  Lastenheft-Prosa.
- **Spec/Lastenheft-CR-Konsistenz:** die vier neuen AK (`:253-256`) decken
  Happy (`--id-prefix AC`), Ableitung, Boundary (Platzhalter), Konflikt und
  spiegeln je einen `TestCLI037_*`. Versions-Bump 0.19.0→0.20.0 mit
  Historie-Zeile (`:800`) und Breaking-Ausweis; Anforderungs-Anlege-Prozess
  (`harness/conventions.md`: ID-Schema, AK, Versions-Bump, „CR statt ADR
  fürs Lastenheft") gewahrt.
- **MR-006 (Referenzrichtung):** der neue Prosa-Block in `DC-FA-CLI-006`
  (`:226-236`) enthält keinen Markdown-Link abwärts auf `docs/plan/adr/*`
  oder Slices; die `docs/plan/adr/`-Vorkommen im Block sind inline-Code-
  `target`-Pfade/Beispiel-IDs (Happy-AK mit `d-check:ignore`), keine Links.
  Das `slice-037`-Token in der Historie-Zeile ist Klartext-Label im Muster
  aller Vorzeilen.
- **ADR-0015-Gestalt/Index:** Status `Proposed`; `Bezug` verlinkt nur
  Lastenheft-Anker (auf-/seitwärts), keine superseded/deprecated ADR; der
  ADR-Index (`docs/plan/adr/README.md:29`) ist um die Zeile ergänzt (Status,
  Datum, Bezüge konsistent). Verglichene Alternativen + Re-Evaluierungs-
  Trigger + Geschichte vorhanden. (Einschränkung „keine Spec-Stelle" siehe
  MEDIUM-1.)
- **Slice-Status/Ablage:** Slice steht unter `…/in-progress/` mit Status
  `in-progress` und offenem DoD-Punkt „Review R1 / Closure / ADR→Accepted" —
  korrekt vor dem Review; kein verfrühtes `done`/`Accepted`.
- **doc-check-Selbstkonformität:** `make doc-check` grün (93 Datei(en), 0
  Befund(e)) auf dem aktuellen Stand — die Markdown-Edits (CHANGELOG, ADR,
  Handbuch, Lastenheft, Slice) brechen keine Links/Anker/Linkpflicht.
- **Regress (alte ai-harness-Tests):** `make test` ⇒ alle Pakete `ok`
  (`…/cli`, `…/core/app`, `…/core/rules`, `…/configyaml` u. a.), Exit 0; die
  bestehenden `TestCLI006_AiHarness*`-Fälle laufen unverändert mit. Adversarial
  bestätigt: `ai-harness`/`ai-harness-init` ohne neue Optionen liefern weiter
  das erwartete (repo-bewusste/Voll-Kanon) Gerüst, Exit 0.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 2 | 1 | 1 |

## Verdikt

**Bedingt freigabefähig — zwei MEDIUM vor Closure zu klären.** Der Kern der
Implementierung ist tragfähig: die Präfix-Ableitung (`deriveReqPrefix`/
`reqShape`) ist korrekt und deterministisch (Single-/Mehrfach-Präfix,
Nicht-Anforderungs-IDs ignoriert, sortierte Konflikt-Meldung), der
Platzhalter-Pfad dekodiert über den eigenen Parser, das explizite Flag
gewinnt vor der Ableitung, die Validierung (`ValidIDPrefix`) und die
`reorderArgs`-Verdrahtung greifen, read-only und Layer-Trennung bleiben
gewahrt, Regress ist ausgeschlossen (`make test` grün), und die Markdown-
Edits sind doc-check-konform.

Blockierend für die Closure sind die beiden MEDIUM: **MEDIUM-1** — die
Messmethoden-Spezifikation §`DC-FA-CLI-006.a` wurde nicht nachgezogen (zeigt
weiter fixes `DC-`, kennt weder `--id-prefix` noch Platzhalter), und ADR-0015s
`Schärft: keine Spec-Stelle` ist angesichts der existierenden, einschlägigen
Spec-Schicht unzutreffend; **MEDIUM-2** — der Vertrag „`ai-harness` +
`--id-prefix` überschreibt die Ableitung / übergeht den Konflikt" ist
funktionsfähig, aber ungetestet. LOW-1 (aktiver `<PREFIX>` = stilles Grün bei
ignoriertem TODO) und INFO-1 (`--id-prefix` außerhalb der ai-harness-Modi
wirkungslos) sind dokumentations-/härtungswürdig, blockieren aber nicht. Die
Gate-Bestätigung im Übrigen obliegt der getrennten Verifikation (hier nur die
explizit benannten Läufe als Beleg, nicht pauschal als grün angenommen).
