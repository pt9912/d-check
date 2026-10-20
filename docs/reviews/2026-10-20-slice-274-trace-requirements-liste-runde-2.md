# Review slice-274 Runde 2 — Fix-Commits zu `trace.requirements` als Liste

- **Review-Art:** Code (gegen Plan, ADR, Hard Rules)
- **Gegenstand:** `slice-274`, Diff `88a6dacd..5914ae90` (Plan-Notiz `c4134726`, Fix `5914ae90`); Befunde der Runde 1 (`docs/reviews/2026-10-20-slice-274-trace-requirements-liste.md`)
- **Skill:** `reviewer.md` @ 1.20.0
- **Modell-ID:** claude-sonnet-5-5
- **Datum:** 2026-10-20
- **Eingangs-Kontext:** `slice-274` (Abschnitt „Plan-Änderung nach Review-Runde 1“), `ADR-0108`, `DC-FA-REQ-001` / `DC-FA-REQ-001.a` Schritt 8 samt Fehlerpräzedenz, `ADR-0104`, `AGENTS.md` §3.2 / §3.7 / §3.8, vorherige Findings F-1 bis F-6.
- **Zusätzlich gefahren (alles danach gelöscht bzw. zurückgesetzt, `git status` sauber):** Wegwerf-Tests gegen `configyaml.Decode` (15 Formen, 3 zyklische Anker), ein Wegwerf-Test gegen `referencePattern`, eine Mutation von `referencePattern` (Vereinigung abgeschaltet) mit `make test`.

## Schließung der Befunde aus Runde 1 (je nach Aussage)

- **F-1 (Aliase/Merge-Keys): geschlossen.** `Node.Decode` löst auf; Probe: Alias auf Skalar (`source: *s`), Merge-Key mit Anker, Merge-Key mit Sequenz `<<: [*a, *b]`, ganzer Eintrag als Alias, Alias auf Sequenz in `modality.levels` werden angenommen (die dort gemeldete „nondeterministisch“-Meldung ist die fachliche Prüfung des doppelten Keywords, nicht der Leseweg).
- **F-3 (Zeile): geschlossen.** Zeilen kommen aus dem Knoten des Dokuments; Test mit Kommentaren/Leerzeile in Objekt- und Listen-Form. Die Aussagen-Suche des Plans (`Neu-Serialisierung`, `yamlLineRef`) trifft nichts mehr; Grenze steht nirgends mehr im Handbuch.
- **F-4 (Präzedenz): geschlossen.** `loadTraceRequirements` extrahiert alle Einträge, bevor das Kreuz-Duplikat geprüft wird; Test `TestTraceSourcesLaterEntryErrorPrecedesDuplicate` erzeugt ein Duplikat zwischen 0 und 1 und einen Fehler in 2 — mit der Schleifen-Form des Vorstands rot.
- **F-5 (`AK-<N>`): geschlossen** (`grep` über `internal/`: nur `LH-AK-002`, Zufallstreffer).
- **F-2 (Referenzseite): geschlossen, Mutation belegt.** `referencePattern` auf Einzel-Muster gestellt: rot werden `TestTraceSourcesCoverageRangesUseAllPatterns`, `TestTraceSourcesCrossConsistencyUsesAllPatterns`, `TestTraceSourcesReferenceSideSeesAllPatterns`, `TestReferencePatternLongestMatch` und der CLI-Test zu `--require-complete`.
- **F-6 (vakuum/Grenze): teilweise geschlossen.** Prädikat ist jetzt `"1 Requirements-Waise(n)"`, die ADR-0108-Grenze ist gepinnt; der Fehlerzweig der Mustervereinigung bleibt ohne Test, siehe F-3 unten (die Begründung dafür ist falsch).

## Findings

### F-1 — LOW — Plan-Notiz behauptet „kein Auslöser gefunden“, der Zweig ist erreichbar

- **quelle:** `AGENTS.md` §5 Regel 13 / 15 (Grenze gegen den Gegenstand prüfen; Behauptung nicht größer als die Arbeit)
- **pfad:** `docs/plan/planning/in-progress/slice-274-trace-requirements-als-liste-von-quellen.md` · „ein Auslöser wurde nicht gefunden“
- **befund:** Ein einzelnes `id-pattern` mit 999 verschachtelten Gruppen kompiliert allein, die Vereinigung `(?:…)|(?:…)` in `referencePattern` scheitert mit „expression nests too deeply“ (gefahren: n=999 scheitert, n=998 nicht). Der Zweig ist damit mit einem Test pinnbar; zudem enthält die Meldung dann das gesamte Muster.
- **verifizierbar:** ja — Wegwerf-Test in `app`.
- **klasse:** unerreichbar-behauptet-erreichbar

### F-2 — LOW — Fixture-Kennung `SFC-QA-01` wird als Anforderung in die Abdeckungsdatei abgeleitet

- **quelle:** `ADR-0104` (Test-Nachweise entlasten in der RTM); `docs/user/abdeckung-tests.md` („Anforderungs-Kennung im Doc-Kommentar“)
- **pfad:** `docs/user/abdeckung-tests.md` · „`SFC-QA-01` | `TestTraceSourcesCrossConsistencyUsesAllPatterns`“
- **befund:** Der Doc-Kommentar des neuen Tests nennt die Fixture-Kennung `SFC-QA-01`; `make abdeckung` führt sie als Zeile mit Link auf das Lastenheft, in dem sie nicht existiert. Der Test steht dadurch unter einer nicht vorhandenen Anforderung und nicht unter `DC-FA-REQ-001` bzw. `DC-FA-XREF-001`.
- **verifizierbar:** ja — `grep -rn "SFC-QA-01" spec/` ohne Treffer.
- **klasse:** fixture-kennung-als-anforderung

### F-3 — INFO — Grenzen von `unknownKeys` nicht benannt

- **quelle:** `AGENTS.md` §3.7 (Klasse Grenze) — kein Anker für LOW, daher INFO
- **pfad:** `internal/adapter/driven/configyaml/yamlstrict.go` · „func yamlFieldTypes“
- **befund:** Der Helfer kennt `,inline` und Zieltypen mit eigenem `UnmarshalYAML` nicht (ein inline-Feld würde als unbekannt gemeldet — fail-laut, nicht still). Heute trägt kein Zieltyp beides; die Grenze steht weder am Helfer noch in einem Test. Zyklische Anker laufen gefahren nicht in `unknownKeys` (Decode meldet zuvor „contains itself“; kein rekursiver Zieltyp).
- **verifizierbar:** nein
- **klasse:** helfer-grenze-unbenannt

### INFO — Kommentar-Präfixe aus mechanischem Austausch

Mehrere Test-Kommentare tragen nach dem `AK-<N>`-Austausch ein Wort-Präfix, das den Satz doppelt (z. B. „Vereinigung: die Vereinigung führt …“, „Leere Liste: eine leere Liste …“). Kein Konventions-Anker; nicht blockierend.

## Negativbefunde (geprüft, ohne Befund)

- **Fehlerformen des neuen Lesewegs (Prüffrage 19):** Alias auf Skalar / Sequenz, ganzer Eintrag als Alias, Merge-Key mit Mapping und mit Sequenz von Aliasen, unbekannter Schlüssel hinter Merge-Key (Zeile der Definition), unbekannter Schlüssel an Tiefe 1 (Eintrag), 2 (`table`, `modality`), in einem flow-Mapping im Merge, Objekt-Form, Map-Schlüssel in `modality.levels` (bleibt frei, korrekt), `table:` null, nicht-String-Schlüssel (`1: x` ⇒ gemeldet), doppelter Schlüssel (Decode meldet, Zeile korrekt), Typfehler plus unbekannter Schlüssel (beide gemeldet), zyklische Anker (Merge, Sequenz, Map ⇒ Fehler, kein Stack-Überlauf). `yaml.Node`-Felder und Zeiger-Typen sind im Code berücksichtigt; ein `yaml.Node`-Feld trägt kein Zieltyp.
- **Objekt-Form gegenüber Vorzustand:** Meldungstext `field X not found in type configyaml.<Typ>` und Zeile des Dokuments entsprechen der Form von `KnownFields` (Typname je Struktur, Zeile des Schlüssels); der Zeilenversatz entfällt.
- **Hard Rules:** §3.2 keine `//nolint`; §3.1 nur `make` (Mutation und Probes über `make test`); §3.7 Kommentare in `yamlstrict.go` und `trace_sources.go` tragen Zusage/Abgrenzung/Kopplung mit einem `DC-FA-REQ-001(.a)`-Anker, keine Review-Historie; §3.8 keine neue ungescannte Eingabe (der Knotenbaum desselben Dokuments).
- **Abgrenzung §1:** keine Ausweitung; Fix liegt in `yamlstrict.go`, `trace_sources.go` und Tests.
- **Tests wirklich rot ohne Fix:** F-2-Tests per Mutation belegt (siehe oben); Alias-/Zeilen-Tests und der Präzedenz-Test waren gegen den Vorstand (Neu-Serialisierung, Schleifen-Prüfung) aus den Befunden der Runde 1 beschrieben rot — nicht erneut gegen den alten Commit gefahren.

## Kategorie-Summary

HIGH 0 · MEDIUM 0 · LOW 2 (F-1, F-2) · INFO 1 (F-3) plus ein INFO ohne Nummer

## Verdikt

Nicht blockierend. F-1 bis F-5 der Runde 1 sind der Aussage nach geschlossen, F-6 teilweise (Zweig-Begründung falsch, F-1 dieser Runde). Die beiden LOW sind vor Merge billig behebbar. Plan-/ADR-Konformität gegeben; DoD und Gate-Läufe gehören der Verifikation.
