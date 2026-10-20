# Review slice-274 — `trace.requirements` als Liste von Quellen

- **Review-Art:** Code (gegen Plan, ADR, Hard Rules)
- **Gegenstand:** `slice-274`, Diff `68cae247..HEAD` (Plan-Änderung, Code, Handbuch)
- **Skill:** `reviewer.md` @ 1.20.0
- **Modell-ID:** claude-sonnet-5-5
- **Datum:** 2026-10-20
- **Eingangs-Kontext:** `slice-274` (§1 Abgrenzung, §3 Plan samt Plan-Änderung, §6 Risiken), `ADR-0108`, `DC-FA-REQ-001` (Lastenheft, Absatz „Mehrere Anforderungsquellen“ und Akzeptanzkriterien), `DC-FA-REQ-001.a` Schritt 8 samt Fehlerpräzedenz, `AGENTS.md` §3.2 / §3.7 / §3.8.
- **Zusätzlich gefahren:** ein Wegwerf-Test (Probe, danach gelöscht, nicht im Diff) gegen `configyaml.Decode` mit Alias-/Merge-Key- und Kommentar-Eingaben.

## Findings

### F-1 — MEDIUM — Strikt-Dekodierung per Neu-Serialisierung bricht YAML-Aliase und Merge-Keys

- **quelle:** Prüffrage 19 (neuer Leseweg, nicht gegen seine Fehlerformen gefahren); `DC-FA-REQ-001` (Objekt-Form „bleibt gültig“)
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go` · „yaml.Marshal(n)“ in `decodeStrictNode`
- **befund:** Jeder Eintrag (und die Objekt-Form) wird isoliert neu serialisiert und dekodiert; ein Alias oder `<<: *b` auf einen Anker außerhalb des Teilbaums wird dadurch unauflösbar. Gefahren: `- &b {source: …}` gefolgt von `- <<: *b` und von `- *b` liefern beide Exit 2 mit `yaml: unknown anchor 'b' referenced`; derselbe Weg trifft eine Objekt-Form, die auf einen Anker außerhalb ihres Blocks zeigt (vor dem Slice gültig). Das ist das natürliche DRY-Idiom für zwei Einträge mit gleicher `source` — genau der Fall des CR.
- **verifizierbar:** ja — Test in `configyaml`, der Merge-Key/Alias zwischen Listen-Einträgen dekodiert. Kein vorhandener Test deckt Aliase ab (Fehlerformen-Test prüft nur Skalar, Null, unbekannte Schlüssel).
- **klasse:** weiche-neuer-leseweg-fehlerformen

### F-2 — MEDIUM — Referenzseite mit mehreren Mustern nur an einem von drei Verbrauchern getestet

- **quelle:** Prüffrage 13; `DC-FA-REQ-001.a` Schritt 8 „Muster der Referenzseite“ (Referenzscan, Range-Validierung von `trace.coverage[].ranges`, `cross-consistency`); `slice-274` §6 (dritter Punkt: Wechselwirkung mit `cross-consistency`)
- **pfad:** `internal/hexagon/core/app/trace_sources_test.go` · „TestTraceSourcesReferenceSideSeesAllPatterns“
- **befund:** Getestet sind nur der Slice-Scan und die Funktion `referencePattern` für sich. Range-Expansion (`expandRange`/`expandNumericRange` prüfen `FindString(cand) == cand` gegen die Vereinigung) und `cross-consistency` mit Listen-Form haben keinen Test; ein Regress dort (z. B. ohne `Longest()`) bliebe grün. Der CLI-Test deckt `--require-complete`, nicht diese beiden Wege.
- **verifizierbar:** ja — Test mit Liste aus zwei Einträgen und `trace.coverage` mit `ranges: true` bzw. `cross-consistency`.
- **klasse:** neuer-vertrag-ohne-test-je-verbraucher

### F-3 — LOW — Zeilennummer der Objekt-Form regrediert bei Kommentaren/Leerzeilen im Block

- **quelle:** `DC-QA-02` (Objekt-Form unverändert); Prüffrage 19 (LOW wegen reiner Diagnose)
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go` · „Grenze: Kommentare und Leerzeilen im Teilbaum zählt die Neu-Serialisierung“
- **befund:** Gefahren: Objekt-Form mit zwei Kommentarzeilen und einer Leerzeile vor `source:` und Tippfehler `sorce` in Zeile 7 meldet `line 10`; vor dem Slice stand dort die richtige Zeile. Der Code-Kommentar benennt die Grenze, der Test (`TestDecode_RequirementsListErrorLine`) prüft nur den kommentarfreien Fall, und die Abweichung geht auch in die falsche Richtung (zu groß), nicht nur „abweichen“.
- **verifizierbar:** ja
- **klasse:** diagnose-zeile-nach-neuserialisierung

### F-4 — LOW — Fehlerpräzedenz weicht von Schritt 8 ab

- **quelle:** `DC-FA-REQ-001.a` Fehlerpräzedenz („… → Nullmenge → Duplikat über Einträge → Referenz-/Coverage-Scans“; Schritt 8: „Nach dem Extrahieren **aller** Einträge wird geprüft …“)
- **pfad:** `internal/hexagon/core/app/trace_sources.go` · „if prev, dup := out.owner[id]; dup“
- **befund:** Die Kreuz-Duplikat-Prüfung läuft innerhalb der Schleife nach jedem Eintrag; ein Duplikat zwischen Eintrag 0 und 1 verdeckt damit eine fehlende Quelle oder Nullmenge in Eintrag 2, die laut Spezifikation zuerst gemeldet würde. Beide Exit 2, nur die Meldung unterscheidet sich; kein Test pinnt die Reihenfolge.
- **verifizierbar:** ja
- **klasse:** fehlerpraezedenz-schleife-statt-nach-allen

### F-5 — LOW — Test-Kommentare tragen `AK-<N>` als Herkunftsmarke

- **quelle:** `AGENTS.md` §3.7 (Herkunft nur als auflösbares Feld `DC-*`/`ADR-*`/`MR-*`/`seit welle-<NN>`; keine Mess-Labels)
- **pfad:** `internal/hexagon/core/app/trace_sources_test.go` · „// AK-1: die Vereinigung führt erst“ (ebenso `cli_trace_sources_test.go`, `trace_requirements_list_test.go`)
- **befund:** `AK-1`…`AK-7` lösen nur im Slice-Plan/CR auf, nicht im Lastenheft (dort sind die Kriterien benannt, nicht nummeriert). Die Kommentare tragen daneben echte Zusagen; die Marke ist Überschuss, kein Totalausfall.
- **verifizierbar:** nein (kein Gate)
- **klasse:** kommentar-label-ohne-aufloesung

### F-6 — INFO — Vakuumes Prädikat und ungepinnte Grenze

- **quelle:** Maintainability (kein Konventions-Anker)
- **pfad:** `internal/adapter/driving/cli/cli_trace_sources_test.go` · „strings.Contains(stderr, "1")“
- **befund:** Die Stderr-Prüfung auf `"1"` ist praktisch immer wahr; getragen wird der Test vom Exit-Code. Zudem pinnt kein Test die in ADR-0108 benannte Grenze (Listen-Eintrag mit `format: headings` ohne `source` ⇒ null Treffer, kein Fehler), und der Fehlerzweig „Muster lassen sich nicht vereinigen“ in `referencePattern` ist ohne Test.
- **verifizierbar:** nein
- **klasse:** test-prädikat-vakuum

## Negativbefunde (geprüft, ohne Befund)

- **Objekt-Form byte-identisch (RTM-Ausgabe):** `resolveEntries` bildet die Objekt-Form als Ein-Eintrag-Liste mit Index -1 ab; `referencePattern` liefert bei einem Eintrag das Muster unverändert (kein `Longest()`, kein Umbau); `ModalityActive`, Gating und `isFullReqID` unverändert; Index-Präfix entfällt (Test `TestTraceObjectFormErrorsCarryNoIndex`, `…SingleEntryEqualsObjectForm`). Die Gegenprobe gegen `v0.86.1` (Plan §6) habe ich nicht gefahren.
- **Kreuz-Duplikat-Prüfung, beide Seiten:** Test variiert Überschrift-zuerst und Tabelle-zuerst; Meldung nennt Kennung, beide Indizes und den Muster-Hinweis (ADR-0108 §3).
- **Nullmenge/Quelle fehlt je Eintrag, Index in Meldung:** `indexed()` setzt den Index nur in der Listen-Form; Tests für Header-, Nullmengen- und Quellenfehler im Eintrag 1, mit Gegenprobe „nicht Eintrag 0“.
- **`modality` je Eintrag:** `traceModalities` bildet je Eintrag Matcher und Map; ohne `modality` in einem Eintrag gaten dessen Waisen alle; Test `TestTraceSourcesModalityPerEntry` trägt beide Hälften (Orphans/Gating 4/3 und 4/4).
- **Muster-Vereinigung mit `Longest()`:** kurzes Muster vor langem verschluckt die Kennung nicht (`TestReferencePatternLongestMatch`).
- **Fehlerformen der Weiche (außer F-1/F-3):** leere Liste, Null-Eintrag, Skalar, Skalar als Eintrag, unbekannte Schlüssel in Eintrag/Objekt/`table`-Block, Index bei Regex-/Format-/Pfadfehlern — getestet; `requirements:` (null) bleibt nil (yaml.v3 ruft den Unmarshaler für Null nicht).
- **Hard Rules:** §3.2 keine `//nolint`; §3.1 nur `make`; ADR-0005 (Importe nur innerhalb Kern/Adapter, `arch-check` laut Implementer grün); §3.8 keine neue ungescannte Eingabe — dieselbe `source` wird je Eintrag neu gelesen, was Schritt 8 so festlegt; §3.7 Produktions-Kommentare tragen Zusage/Grenze/Kopplung mit `DC-FA-REQ-001(.a)`-Anker (nur F-5 in Tests); die bestehenden `slice-067/068`-Marken in der angefassten Zeile in `trace.go` sind Bestand.
- **Abgrenzung (§1):** kein `trace_cross.go`/`trace_table.go`-Edit (Plan-Änderung vor dem Code dokumentiert), keine Generator-Änderung, kein CHANGELOG/Kopf, ADR-Text unverändert; Handbuch- und `operations.md`-Änderungen gedeckt (Spiegel aus der Plan-Änderung).
- **Handbuch:** Beschreibung deckt sich mit Schritt 8 (Vereinigung, Index, Duplikat, leere Liste, `modality` je Eintrag); die Grenze der Zeilennummer (F-3) steht dort nicht.
- **Netz/Gate-Pfade:** kein Netzzugriff, keine neue Gate-Logik.

## Kategorie-Summary

HIGH 0 · MEDIUM 2 (F-1, F-2) · LOW 3 (F-3, F-4, F-5) · INFO 1 (F-6)

## Verdikt

Blockiert durch die beiden MEDIUM (F-1 Aliase/Merge-Keys, F-2 Referenzseite-Tests) bis zur Klärung; die LOW/INFO sind nicht blockierend. Plan-/ADR-Konformität der Mechanik ist sonst gegeben; DoD-Abhaken und Gate-Läufe gehören der Verifikation.
