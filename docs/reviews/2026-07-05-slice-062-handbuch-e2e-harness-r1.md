# Review-Report — slice-062 Handbuch-E2E-Harness (R1)

## Kopf-Metadaten

- **Gegenstand:** slice-062 (Dimension B) — E2E-Verankerung der
  Handbuch-Kommando-/Ausgabe-Beispiele gegen das echte Binary.
- **Diff:** NEU `internal/adapter/driving/cli/handbook_examples_test.go`
  (Paket `cli_test`, nutzt `run`/`write`/`traceRepo` aus
  `cli_acceptance_test.go`); GEÄNDERT `docs/user/benutzerhandbuch.md`
  (eine Zeile: `d-check-test:not-replayable`-Marker vor dem
  `--print-mk`-`text`-Block, §4.13).
- **Lauf:** R1 · **Datum:** 2026-07-05 · **Reviewer:** unabhängig (nicht
  Mitautor).
- **Pflichtlektüre:** `.harness/skills/reviewer.md`,
  `docs/plan/planning/done/slice-062-handbuch-e2e-beispiele.md`,
  `docs/plan/planning/done/slice-061-doc-config-beispiel-verifikation.md`,
  `AGENTS.md` §3, `harness/conventions.md`.
- **Bezug:** `DC-FA-CLI-003` (Exit-Codes), `DC-FA-DIST-001` (Image);
  kein CR, kein ADR (E2E-Test-Erweiterung im bestehenden Schnitt).
- **Rolle:** Reviewer, **kein Verifier** — kein `make`-/`go test`-Lauf;
  rein statische Analyse + gezieltes Lesen.

---

## Findings

### F-1 (LOW) — `yaml`-CLI-Ausgabeblöcke sind von der Fail-closed-Klassifikation ausgenommen; ein künftiger yaml-Ausgabeblock rutscht still durch

- **Kategorie:** LOW (latente Wartungsfalle, die erst bei künftigem
  Edit zündet)
- **Quelle:** slice-062 §3 (DoD „Fail-closed-Guards … jeder neue
  Ausgabeblock"), §4 („Ausgabe-Matching-Stabilität"); Maintainability
- **Pfad:** `internal/adapter/driving/cli/handbook_examples_test.go:191`
  (`if info != "text" && info != "json" { continue }`) i. V. m. dem
  Datei-Kommentar Z. 34–36 und der Marker-Zusage Z. 57–59.
- **Befund:** `TestHandbook_OutputBlocksClassified` iteriert nur über
  `text`-/`json`-Blöcke; `yaml` wird per `continue` übersprungen. Die
  Marker-Zusage im Kommentar („Opt-out, nicht Opt-in: ein neuer echter
  Ausgabeblock ist automatisch prüfpflichtig", Z. 57–59) hält damit für
  die yaml-Ausgabe-Klasse nicht. Die eine heutige CLI-Ausgabe-YAML
  (§4.11 `--yaml`, @550) ist über E7 im Replay verankert, aber ein
  **neuer** yaml-Ausgabeblock ist es nicht.
- **Failure-Szenario:** Ein Maintainer dokumentiert das im Handbuch
  bereits zugesagte `--doctor --yaml`-Beispiel (§4.9 „Dieselben …
  Varianten gibt es als YAML") als neuen `yaml`-Block. slice-061s
  Harness (`docexamples_test.go`) lehnt ihn per `configyaml.Decode`
  (strikt) ab ⇒ rot ⇒ Maintainer setzt den `d-check-test:not-config`-
  Marker; slice-061 überspringt ihn nun. slice-062 sieht ihn nie
  (Sweep überspringt yaml; Replay fährt nur die 7 fest verdrahteten
  Beispiele). Ergebnis: der dokumentierte yaml-Ausgabeblock ist an
  **kein** Binary gekoppelt — driftet später `reasonText`→`reason_text`
  o. Ä. in der yaml-Ausgabe (oder ist die Doku-Schlüsselmenge von
  Anfang an falsch), bleibt beides **still grün**. Das ist exakt die
  Silent-Drift-Klasse, die der Slice für text/json schließt — für die
  yaml-CLI-Ausgabe bleibt sie offen. (Grenzt an MEDIUM, weil es die
  Kern-Messmethode betrifft; als LOW eingestuft, weil der heutige
  Bestand vollständig gedeckt ist und die Lücke erst bei einem künftigen
  yaml-Ausgabeblock zündet.)
- **Verifizierbar:** nein (latent; kein heutiger Gate-Lauf färbt rot —
  genau das ist der Punkt). Reproduzierbar durch Hinzufügen eines
  zweiten `yaml`-CLI-Ausgabeblocks mit `not-config`-Marker: `make test`
  bleibt grün.

### F-2 (INFO) — Strukturelle Kopplung ist nur doc ⊆ echt; ein zusätzliches Ausgabe-Feld des Binaries bleibt grün

- **Kategorie:** INFO (dokumentationswürdige, im Kommentar teils
  benannte Grenze der Messmethode)
- **Quelle:** Maintainability; Datei-Kommentar Z. 14–19 („Drift in
  beide Richtungen") und Z. 71–72/276–289 (Schlüssel-Menge ⊆)
- **Pfad:** `internal/adapter/driving/cli/handbook_examples_test.go:284`
  (`for k := range docKeys { if !actKeys[k] … }`)
- **Befund:** Die json/yaml-Kopplung prüft nur `docKeys ⊆ actKeys`.
  Umbenennungen/Entfernungen dokumentierter Schlüssel fängt sie in
  beide Richtungen (Kommentar Z. 277 korrekt), aber ein vom Binary
  **neu hinzugefügtes** Ausgabe-Feld erscheint nicht in `docKeys` und
  löst nichts aus.
- **Failure-Szenario:** `model.Finding` bekommt ein neues
  serialisiertes Feld (z. B. `column`); das Handbuch @532/@550 zeigt es
  nicht. Der Test bleibt grün, obwohl die dokumentierte Ausgabe-Form
  unvollständig gegenüber dem Binary ist (eine Drift-Richtung „Binary
  gewinnt Feld", die die Slice-Prämisse „Binary driftet ⇒ rot" nicht
  abdeckt). Der Form-Anker-Teil (`formTokens`) ist echt bidirektional;
  nur die strukturelle Schlüssel-Kopplung ist einseitig — die Header-
  Formulierung „Drift in beide Richtungen" (Z. 18) gilt streng nur für
  die formTokens, nicht für die Schlüssel-Menge.
- **Verifizierbar:** nein (latent; low-impact — zusätzliches, nicht
  falsches Feld).

---

## Negativbefunde (geprüft, ohne Befund)

- **Kopplungs-Ehrlichkeit / bidirektionale Form-Kopplung (text/json):**
  geprüft — `TestHandbook_AnchoredExamplesReplay` prüft jeden
  `formToken` in **beidem** (Doku-Block Z. 266 und echter Ausgabe
  Z. 270); fehlt er im Doku-Block ⇒ rot, fehlt er in der echten Ausgabe
  ⇒ rot. Alle `formTokens` sind nicht-leer (keine vakuum-wahre
  `Contains("")`-Kopplung). Kein vakuum-grüner Pfad (Beispiel-Liste
  hart, 7 Einträge; `mustBlocks` fatalt bei leerer Block-Menge).
- **Fail-closed-Vollständigkeit für text/json:** geprüft — jeder
  `text`-/`json`-Block muss genau **ein** Beispiel matchen (`matches
  != 1` ⇒ Fehler) oder den `not-replayable`-Marker tragen; ein neuer
  json-Ausgabeblock mit unbekanntem Grund (z. B. `--trace --json`)
  matcht 0 Beispiele ⇒ rot. Disc-Kollision (zwei Blöcke, ein Disc)
  ⇒ `TestHandbook_AnchoredExamplesReplay` `docBlocks != 1` ⇒ rot. Beide
  Tests zusammen sind fail-closed sowohl für „Block ohne Beispiel"
  (Sweep) als auch „Beispiel ohne Block" (Replay). `outputBlocks == 0`
  und `len(blocks) == 0` fatal. (Einzige Klasse ohne diese Eigenschaft:
  yaml — siehe F-1.)
- **Extraktor-Korrektheit (`extractFencedBlocks`):** geprüft — echte
  Fence-Zustandsverfolgung (öffnet auf Fence mit beliebigem Info-String,
  schließt nur auf nacktes Fence); `json` in `markdown` ist Body, kein
  Block (Handbuch @870 korrekt behandelt); unbalanciert ⇒ Fehler
  (fail-closed); `TrimSpace` toleriert eingerückte Fences (Handbuch
  @226/@235/@508 als bash-Bodies erfasst — nützlich für die
  Flag-Kopplung) und CRLF; `firstToken` normalisiert Info-Suffixe
  (`yaml title=x` → `yaml`). `TestExtractFencedBlocks` deckt die
  Kanten inkl. Marker-Vorzeilen-Regel (adjazent, keine Leerzeile).
- **Diskriminator-Eindeutigkeit:** geprüft — die 7 `outputDisc`
  (`geprüft, 0 Befund(e)`, `target-missing`, `Fix-Kandidat:`,
  `# Requirements Traceability Matrix`, `reasonText`, `exitCode:`) sind
  je Fence-Typ eindeutig gegen die realen Blöcke; @429 (doctor-json,
  `id-unlinked`) enthält kein `target-missing`, @532 kein `reasonText`;
  `exitCode:` erscheint nur in der CLI-yaml @550, nicht in den
  Config-yaml. Eine harmlose Doku-Erweiterung (zweiter yaml-
  Ausgabeblock mit `exitCode:`) bricht E7 **laut** (Replay-Uniqueness),
  nicht still — mit dem Slice-Designentscheid „laut > still" konsistent.
- **stdout/stderr-Merge (`combined = stdout+"\n"+stderr`):** geprüft,
  akzeptabel — die verankerten Doku-Ausgabeblöcke mischen selbst die
  Ströme (@116: Befund-Zeile[stdout] + Zusammenfassung[stderr] in einem
  Block), es gibt also keine pro-Strom-Doku-Behauptung zu koppeln; kein
  `formToken` eines Beispiels liegt im „falschen" Strom (E3/E4-Tokens
  nur auf stdout, kein Leck aus stderr). Die dokumentierte
  Strom-Trennung (§3) ist redundant über `TestCLI001_Happy`
  (Zusammenfassung/stderr), `TestCLI003_EinBefund` (Befund/stdout),
  `TestCLI007_*` (stdout vs. stderr) verankert.
- **Abgrenzung zu Bestehendem / Beispiel-Auswahl-Ehrlichkeit:**
  geprüft — Mehrwert real (koppelt Handbuch-Blöcke an das Binary;
  `cli_acceptance_test.go` liest das Handbuch nie). Der Auswahl-Kommentar
  (7 verankert, Rest begründet ausgenommen) deckt sich mit der
  Fence-Sondierung: 5 `text` + 2 `json` = E1–E6 + der markierte
  `--print-mk`-Block; E7 ist die einzige CLI-yaml; §4.5–4.8 tragen keine
  `text`/`json`-Ausgabeblöcke (nur bash + Config-yaml); die übrigen
  10 `yaml` sind Config (von slice-061 gedeckt).
- **Hard Rules (`AGENTS.md` §3):** geprüft — kein Netzzugriff (kein
  Beispiel aktiviert `external`), read-only (nur `t.TempDir()`),
  Determinismus (feste Fixtures, keine Map-Iteration im Prüfpfad); keine
  Inline-Suppression (§3.2); Paket `cli_test` extern, Import nur
  std + `gopkg.in/yaml.v3` (ADR-0005-konform); keine
  Symbol-Kollision mit den fünf weiteren `cli_test`-Dateien
  (`firstToken`/`keysOf`/`collectKeys`/`anyContains`/`mustBlocks`/
  `handbookExamples`/`fencedBlock`/`extractFencedBlocks` je einmalig).
- **Handbuch-Diff:** geprüft — Marker `d-check-test:not-replayable` steht
  in der Zeile **unmittelbar** vor dem Fence (@627/@628, keine
  Leerzeile), mit Grund (abgekürzte Illustration/Elision); der Marker-
  String stimmt mit `notReplayableMarker` überein; der Block ist
  tatsächlich abgekürzt (`# … doc-trace, doc-complete` /
  `# … doc-help`), also zu Recht nicht replaybar.

---

## Kategorie-Summary

| Kategorie | Anzahl |
| --------- | ------ |
| HIGH      | 0      |
| MEDIUM    | 0      |
| LOW       | 1 (F-1) |
| INFO      | 1 (F-2) |

---

## Verdikt

**ACCEPT.**

Der Harness koppelt die dokumentierten `text`-/`json`-Ausgabeblöcke
bidirektional und fail-closed an das echte Binary; die Extraktor-,
Diskriminator- und Guard-Logik ist korrekt und die Beispiel-Auswahl im
Kommentar ehrlich und vollständig. Kein HIGH/MEDIUM (kein Silent-Green im
gedeckten Bereich, keine Hard-Rule-Verletzung), daher nicht blockierend.

Die eine LOW-Beobachtung (F-1) ist bemerkenswert, weil sie die Kern-
Absicht des Slice betrifft: für die **yaml**-CLI-Ausgabe-Klasse ist die
Fail-closed-Zusage nicht eingelöst — ein künftiger yaml-Ausgabeblock
(z. B. das im Handbuch bereits zugesagte `--doctor --yaml`) landet
zwischen slice-061 (per `not-config` ausgenommen) und slice-062 (yaml
nicht im Sweep) ungekoppelt. Empfehlung zur Übergabe (nicht Merge-
blockierend): entweder den Sweep auf yaml-Ausgabeblöcke ausweiten
(mit eigenem Ausnahme-Marker gegen die Config-yaml) oder die Grenze
im Handbuch-Autoren-Hinweis / Datei-Kommentar explizit als bewusste
Nicht-Deckung benennen, damit die Marker-Zusage „automatisch
prüfpflichtig" nicht über die yaml-Klasse überreicht. F-2 (INFO) ist
eine bewusst einseitige Schlüssel-Kopplung — festhalten, nicht handeln.
