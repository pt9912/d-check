<!-- d-check:ignore (Review-Report: enthält illustrative IDs/Pfade/Regex als Belege; docs/reviews/** ist ids/exempt) -->

# Review R2 — slice-037 Fix-Verifikation (`--suggest-config` Kennungs-Präfix, MEDIUM-1/MEDIUM-2)

## Kopf-Metadaten

- **Review-Art:** Adversariale Fix-Verifikation (Runde R2). Gegenstand:
  bestätigen oder widerlegen, dass die zwei MEDIUM aus R1 behoben sind, und
  Regressionen durch die Fixes suchen. Kein Verifier im DoD-Sinn — geprüft
  wird die Code-/Spec-/Test-Substanz; Gate-Belege sind einzeln ausgewiesen,
  nicht pauschal als grün angenommen.
- **Datum:** 2026-06-21
- **Gegenstand:** Fix-Commit `be0fba7` „fix(suggest): slice-037 Review R1 —
  2 MEDIUM behoben (DC-FA-CLI-006, ADR-0015)" über der Implementierung
  `e87fa77`. HEAD = `be0fba7`. Geänderte Dateien im Fix:
  `docs/plan/adr/0015-suggest-config-id-prefix.md`, `spec/spezifikation.md`,
  `internal/adapter/driving/cli/cli_acceptance_test.go` (zwei neue
  `TestCLI037_*`), plus der R1-Report.
- **Reviewer:** Claude (Agent), Skill `.harness/skills/reviewer.md` v1.0.0.
- **Modell:** `claude-opus-4-8[1m]`.
- **R1-Befund (Ausgang):** 0 HIGH / 2 MEDIUM / 1 LOW / 1 INFO
  (`docs/reviews/2026-06-21-slice-037-suggest-config-id-prefix.md`). MEDIUM-1:
  Messmethoden-Spec §`DC-FA-CLI-006.a` nicht nachgezogen + ADR-0015
  „keine Spec-Stelle" unzutreffend. MEDIUM-2: Override-/Konflikt-Bypass-Pfad
  von `ai-harness` + `--id-prefix` ungetestet.
- **Verifikations-Setup:** Doc-Gates und Test-Suite über `make doc-check` /
  `make test` (Docker, `--network none`); adversariale Live-Läufe gegen das
  gebaute Image `d-check:latest` über read-only `/tmp`-Test-Repos
  (ENTRYPOINT setzt `/repo`, Flags angehängt); **Mutationsprobe** der neuen
  Tests (Kopie unter `/tmp`, nicht das Repo) zum Nachweis ihrer Wirksamkeit;
  Negativ-Kontrolle, dass `doc-check` den ADR-Anker tatsächlich validiert
  (gebrochene Kopie ⇒ rot). Das d-check-Repo wurde nicht verändert.

## Findings

### HIGH

Keine.

### MEDIUM

Keine.

### LOW

Keine neuen. (R1-LOW-1 bleibt offen — siehe Negativbefunde / Won't-fix.)

### INFO

Keine neuen. (R1-INFO-1 bleibt offen — siehe Negativbefunde / Won't-fix.)

## Negativbefunde (geprüft, ohne Befund)

- **MEDIUM-1 — Spec nachgezogen (BEHOBEN):** `spec/spezifikation.md`
  §`DC-FA-CLI-006.a` enthält jetzt den Absatz „**Anforderungs-Präfix**"
  (`:144-152`): nur das Anforderungs-Muster trägt `<PREFIX>`,
  `ADR-`/`MR-`/`slice`/Carveout sind konventions-fest; Quelle des Präfix =
  `--id-prefix` (explizit, gewinnt immer) bzw. im Modus `ai-harness` die
  eindeutige FA-/QA-Ableitung (mehrere ⇒ Nutzungsfehler); ohne beides der
  markierte Platzhalter `<PREFIX>` + `# TODO`, kein stiller `DC-`. Das deckt
  sich exakt mit der Implementierung (`suggest.go:91-98` Vorrang
  `idPrefix`→`deriveReqPrefix`→Platzhalter) und dem Lastenheft-CR
  (`spec/lastenheft.md`, Zeilen 226-236; AK 253-256).
- **MEDIUM-1 — Vorlagen-Muster parametrisiert (BEHOBEN):** in der
  „Kanonischen Vorlage" steht das Anforderungs-Muster jetzt als
  `regex: '<PREFIX>-(FA-[A-Z]+|QA)-\d+'` (`:175`) mit erklärendem Kommentar
  davor (`:173-174`), nicht mehr fix `DC-(FA-…`. `ADR-\d{4}`, `MR-\d{3}`,
  `slice-\d{3}` bleiben unparametrisiert (`:165/169/179`) — Geltungsbereich
  wie ADR-0015 §4.
- **MEDIUM-1 — ADR-0015 `Schärft` korrekt + Anker löst auf (BEHOBEN):**
  `docs/plan/adr/0015-suggest-config-id-prefix.md` (Zeilen 14-17) zeigt jetzt
  `Schärft: spec/spezifikation.md §DC-FA-CLI-006.a … — nicht das Lastenheft`
  (statt „keine Spec-Stelle"). Der Link-Anker
  `…/spec/spezifikation.md#dc-fa-cli-006a--konfigurations-vorschlag` löst auf:
  Slugify des Headings `### DC-FA-CLI-006.a — Konfigurations-Vorschlag` ergibt
  `dc-fa-cli-006a--konfigurations-vorschlag` (Punkt + Em-Dash entfallen, die
  zwei Leerzeichen um den Em-Dash ⇒ `--`). Negativ-Kontrolle: in einer
  `/tmp`-Kopie den Anker absichtlich gebrochen ⇒ `make`-Image meldet
  `docs/plan/adr/0015-…:15 … anchor-missing`, Exit 1 — d.h. `doc-check`
  validiert diesen Anker wirklich; auf dem realen Tree ist er grün (s. u.).
- **MEDIUM-1 — Regelrichtung (BEHOBEN/konform):** ADR schärft die
  **Spezifikation** (Spec-Stratum), nicht das Lastenheft (MR-006 /
  AGENTS §3.5). Der neue Spec-Prosa-Block (`:144-152`) enthält **keinen**
  Markdown-Link abwärts auf ADR/Slice (geprüft: kein `](`-Vorkommen). Die
  `docs/plan/adr/`/`slice-`-Tokens in §`DC-FA-CLI-006.a` liegen ausschließlich
  im fenced ```yaml-Block (Zeilen 157–194: `target:`/`regex:`/matrix-Pfade) —
  Fence-Inhalt, von ids/links/codepaths per Vorverarbeitung ausgenommen.
- **MEDIUM-2 — neue Tests decken den Vertrag (BEHOBEN):**
  `TestCLI037_IDPrefix_UeberschreibtAbleitung`
  (`internal/adapter/driving/cli/cli_acceptance_test.go:1171`) prüft, dass
  `ai-harness --id-prefix ZZ` über ein DC-Lastenheft das Muster `ZZ-FA-CLI-001`
  matcht und `DC-FA-CLI-001` **nicht** (Override der Ableitung).
  `TestCLI037_IDPrefix_FlagUebergehtKonflikt` (`:1187`) prüft, dass ein
  mehrdeutiges Lastenheft (DC + AC) mit `--id-prefix ZZ` Exit 0 liefert und
  `ZZ-FA-A-001` matcht (Konflikt übergangen). Beide nutzen den realen
  Decode-Pfad (`reqPattern`→`configyaml.Decode`).
- **MEDIUM-2 — Tests sind aussagekräftig (Mutationsprobe):** in einer
  `/tmp`-Kopie das Override-Guard (`if reqPrefix == "" && harness`,
  `suggest.go:92`) zu „Ableitung gewinnt immer" mutiert ⇒ beide neuen Tests
  werden **rot** mit passender Diagnose: `UeberschreibtAbleitung` →
  „überschreibt die DC-Ableitung nicht: \"DC-(FA-[A-Z]+|QA)-\d+\"";
  `FlagUebergehtKonflikt` → „Exit = 2 … mehrdeutiges Anforderungs-Präfix
  (AC, DC)". Sie würden eine Regression in der Ableitungs-/Konflikt-
  Verzweigung also fangen.
- **MEDIUM-2 — Live-Verhalten am Image bestätigt:** gegen `d-check:latest`
  (Build aus `be0fba7`): `--suggest-config ai-harness --id-prefix ZZ` über ein
  DC-Lastenheft ⇒ `regex: 'ZZ-(FA-[A-Z]+|QA)-\d+'`, Exit 0 (Override);
  dasselbe über ein konfligierendes (DC + AC) Lastenheft ⇒ `ZZ-…`, Exit 0
  (Konflikt-Bypass); Negativ-Kontrolle ohne Flag ⇒ Exit 2,
  „mehrdeutiges Anforderungs-Präfix im Lastenheft (AC, DC)" (Präfixe sortiert
  — DC-QA-02).
- **Regression — Test-Suite (geprüft, ohne Befund):** `make test` ⇒ alle
  Pakete `ok` (`…/cli`, `…/core/app`, `…/core/rules`, `…/configyaml`,
  `…/httpcheck`), Exit 0. Die bestehenden `TestCLI037_*` (Explizit,
  Platzhalter, AbleitungAiHarness, KonfliktFehler, Ungueltig) und die
  `TestCLI006_*`-ai-harness-Fälle laufen unverändert mit; kein Regress an der
  Alt-Familie.
- **Regression — doc-check (geprüft, ohne Befund):** `make doc-check` ⇒
  „94 Datei(en) geprüft, 0 Befund(e)", Exit 0. Die Spec-/ADR-Edits brechen
  keine Links/Anker/Linkpflicht; der neue ADR→Spec-Anker ist gültig
  (s. MEDIUM-1).
- **Regression — `<PREFIX>` im Vorlagen-Block (geprüft, ohne Befund):** die
  Muster-Änderung `DC-…`→`<PREFIX>-…` steht im fenced ```yaml-Block
  (`:157-194`); Fence-Zeilen werden von allen Modulen ignoriert, daher löst
  `<PREFIX>`/`<…>` keinen ids-/codepaths-Befund aus. Beleg: doc-check auf dem
  realen Tree bleibt 0 Befunde (sonst hätte das illustrative Muster oder ein
  Pfad-Token im Block angeschlagen).
- **R1-LOW-1 (aktiver `<PREFIX>`-Platzhalter) — Won't-fix korrekt
  eingeordnet:** unverändert advisory-Ausgabe (kein d-check-eigenes Gate),
  markierter Platzhalter + unmittelbar vorangestellter `# TODO`-Kommentar,
  ausdrücklich gewählte Anti-Footgun-Variante (ADR-0015 §3, §„Verglichene
  Alternativen"). Keine Eskalation: das Versagen erfordert weiterhin, dass der
  Nutzer das markierte TODO ignoriert; die Entscheidung ist in Spec
  (`:144-152`), Lastenheft (`:232-236`) und ADR dokumentiert.
- **R1-INFO-1 (`--id-prefix` außerhalb der ai-harness-Modi wirkungslos) —
  Won't-fix korrekt eingeordnet:** unverändertes, stabiles Verhalten; im
  Fix-Commit als „im Flag-Hilfetext verortet" deklariert. Keine Eskalation —
  kein Korrektheitsdefekt, nur eine Verhaltensvariante; Stufe bleibt INFO.
- **Slice-Status/Ablage (geprüft, ohne Befund):** Slice steht weiter unter
  `docs/plan/planning/in-progress/`; ADR-0015 Status `Proposed`. Korrekt vor
  der Closure — kein verfrühtes `done`/`Accepted`. Mit 0 HIGH / 0 MEDIUM in
  R2 ist der Übergang ADR→`Accepted` und die Closure jetzt begründet.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 0 | 0 | 0 |

## Verdikt

**Freigabefähig — beide R1-MEDIUM bestätigt behoben, keine neuen Findings, keine
Regression.** MEDIUM-1 ist auf allen drei Ebenen geschlossen: die
Messmethoden-Spec §`DC-FA-CLI-006.a` beschreibt jetzt `--id-prefix`/Ableitung/
Platzhalter und das `<PREFIX>`-parametrisierte Vorlagen-Muster konsistent mit
`suggest.go` und dem Lastenheft-CR; ADR-0015 `Schärft` zeigt korrekt auf die
Spezifikation (nicht das Lastenheft, MR-006-konform) und der Anker löst
nachweislich auf (Negativ-Kontrolle: brechen ⇒ `anchor-missing`). MEDIUM-2 ist
durch zwei aussagekräftige Tests geschlossen — die Mutationsprobe zeigt, dass
beide bei einer Regression der Override-/Konflikt-Verzweigung rot werden; das
Live-Verhalten am Image (`ZZ`-Override + Konflikt-Bypass, Exit 0; sortierte
Konflikt-Meldung ohne Flag) ist bestätigt. `make test` und `make doc-check`
sind grün (94/0), die `<PREFIX>`-Änderung im fenced Vorlagen-Block bricht
nichts. R1-LOW-1 und R1-INFO-1 bleiben als bewusster Won't-fix korrekt
eingeordnet (keine Eskalation). Damit ist slice-037 **closure-/accept-fähig**
(ADR-0015 → Accepted).
