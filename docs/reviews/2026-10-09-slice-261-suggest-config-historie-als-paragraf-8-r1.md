# Review slice-261 — `--suggest-config` kennt die Spezifikations-Historie als §8 (R1)

- **Review-Art:** Code (gegen Plan, ADRs, Hard Rules)
- **Gegenstand:** slice-261, Commit `fa995ce8`
- **Skill:** `.harness/skills/reviewer.md` @ 1.20.0 (`94f798ef`)
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** Slice-Plan slice-261 (§1 Ziel und Abgrenzung, §3 Plan);
  [`DC-FA-CLI-006`](../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten),
  [`DC-FA-CLI-006.a`](../../spec/spezifikation.md#dc-fa-cli-006a--konfigurations-vorschlag);
  [`MR-025`](../../harness/conventions.md#mr-025) (Spiegel), [`MR-075`](../../harness/conventions.md#mr-075)
  (Historie-Ausnahme folgt der Umnummerierung); `AGENTS.md` §3.4, §3.7, §5 Regel 17;
  Vorbefund F-6 aus dem Review von slice-259 (Anlass dieses Slice);
  Baseline `v6.17.0` · `templates/spec/spezifikation.template.md` (Historie als `## 8. Historie`).

## Findings

### F-1 — LOW: Der Vorschlag nennt sich an Baseline v1.3.0 gebunden und trägt jetzt eine Überschrift der v6.17.0-Vorlage

- **kategorie:** LOW
- **quelle:** [`DC-FA-CLI-006`](../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten) Out-of-Scope · „sie sind an **eine** adoptierte Baseline-Version gebunden (im Kommentar-Header genannt)"
- **pfad:** `internal/hexagon/core/app/suggest.go` · „harnessBaseline   = \"v1.3.0\""; `spec/spezifikation.md` §`DC-FA-CLI-006.a` · „ai-harness-course-Konvention (Baseline\n**v1.3.0**)"
- **befund:** Der Slice begründet `"8. Historie"` mit der Vorlage seit Baseline `v6.17.0` (Plan §1, Testkommentar „der Ort nach der Baseline-Vorlage"), der Kommentar-Header der Ausgabe und die Spezifikation binden die Vorlage aber weiter an `v1.3.0`, die diese Überschrift nicht kennt. Die Drift ist Bestand (`"7. Historie"` stand schon vor slice-035 im Vorschlag); der Diff verbreitert sie um einen Eintrag, den die genannte Baseline nicht trägt.
- **verifizierbar:** nein — kein Gate hält den Header gegen die Herkunft der Einträge.
- **klasse:** baseline-version-header-drift

### F-2 — INFO: `docs/user/abdeckung-tests.md` steht nicht in Plan §3

- **kategorie:** INFO
- **quelle:** Slice-Plan slice-261 §3
- **pfad:** `docs/user/abdeckung-tests.md` · „`TestCLI006_AiHarness_SchlaegtBeideHistorienVor`"
- **befund:** Der Diff ändert eine dritte Datei, die die Plantabelle nicht führt. Sie ist abgeleitet (`make abdeckung`) und `make test` verlangt sie zum neuen Test; eine Abgrenzungsverletzung ist das nicht, eine Plan-Unschärfe schon.
- **verifizierbar:** ja — `make test` hält die Datei gegen ihre Ableitung.
- **klasse:** plan-tabelle-ohne-abgeleitete-datei

## Negativbefunde

- **Test-Grund (bewusstes Brechen):** geprüft, ohne Befund. In einem Wegwerf-Worktree `suggest.go` auf `fa995ce8~1` zurückgesetzt, `make test`: genau `TestCLI006_AiHarness_SchlaegtBeideHistorienVor` rot mit `exclude-sections ohne "8. Historie": [Historie 7. Historie Geschichte]` — der richtige Grund, kein Dekodier- oder Exit-Fehler. Der Test prüft über den eigenen Parser (`configyaml.Decode`) den dekodierten Wert, nicht einen String-Treffer — eine kommentierte oder falsch gequotete Zeile fiele auf.
- **Semantik der Ausnahme:** geprüft, ohne Befund. `exclude-sections` vergleicht den vollen Heading-Klartext (Lastenheft `DC-FA-CLI-009`-Abschnitt „exakt wie `matrix.exclude-sections`"); `Historie` deckt `8. Historie` also nicht — der Eintrag ist nötig. Die vendorten Vorlagen `v6.17.0` führen `## 7. Historie` (Lastenheft), `## 8. Historie` (Spezifikation), `## Geschichte` (ADR); die vorgeschlagene Liste deckt alle drei.
- **Spiegel nach `MR-025`:** geprüft. Lastenheft `DC-FA-CLI-006` zählt die `exclude-sections` nicht auf — kein Spiegel, keine Versionspflicht. Spezifikation `DC-FA-CLI-006.a`: kanonisches Beispiel nachgezogen (Zeile `exclude-sections`), Historie-Zeile ergänzt; die Vorlage bleibt 1:1 zur Ausgabe. `--print-config`-Vorlage (`config_template.go`) zeigt `exclude-sections: [Historie]` als generisches Gerüst, nicht den ai-harness-Vorschlag — kein Spiegel. Handbuch: keine Fundstelle des Vorschlagswerts (das `matrix`-Beispiel in §5 ist generisch) — für die Release-Prep nichts nachzuziehen außer der §11-Versionszeile. CHANGELOG: kein Eintrag im Feature-Commit, wie Regel 17 es verlangt; der Eintrag gehört in die Release-Prep nach `0.85.0`.
- **Abgrenzung:** geprüft, ohne Befund. Nur die eine Vorschlagszeile geändert, keine anderen `--suggest-config`-Vorschläge, `.d-check.yml` des Repos unberührt (trägt `"8. Historie"` seit slice-259).
- **Kommentar-Regel §3.7:** geprüft, ohne Befund. Der neue Testkommentar ist eine Zusage mit `DC-*`-Feld, ohne Slice-Nummer, Befund-Marker oder Herkunfts-Prosa; der Kommentar über `renderHarnessMatrix` stimmt weiter („exclude-sections sind pfad-unabhängig und immer aktiv").
- **Spec-Straten abwärts (§3.4):** geprüft, ohne Befund. Die neue Historie-Zeile der Spezifikation nennt weder Slice noch ADR noch Commit.
- **Commit-Botschaft (Anker 8):** geprüft, ohne Befund. „blackbox-probe … byte-identisch" behauptet nur Unverändertheit der Pfade ohne die Option, nicht mehr; „Test vor dem Fix rot" ist oben nachgemessen.
- **Hexagon/Netz/Suppression/Gate-Pfade:** nicht berührt.

## Kategorie-Summary

HIGH 0 · MEDIUM 0 · LOW 1 · INFO 1

## Verdikt

Freigabe. F-1 ist Bestand, den der Diff um einen Eintrag verbreitert; ob der
Header gehoben oder die Bindung umformuliert wird, ist ein eigener Vorgang
(außerhalb der Abgrenzung von slice-261, die nur diese Überschrift betrifft).
F-2 ist Information.
