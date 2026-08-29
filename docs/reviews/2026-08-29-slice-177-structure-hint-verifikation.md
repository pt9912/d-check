# Verifikation slice-177 — DoD, Spec-Konformität und Mutationsproben

**Gegenstand:** [slice-177](../plan/planning/done/slice-177-structure-hint.md), Stand `c894c7f` … `86d2985`.
**Datum:** 2026-08-29. **Verifikation:** unabhängiger Subagent, eigener Kontext.
**Aufbau:** zwei Images — `d-check:latest` (HEAD) und ein aus dem Stand vor dem Feat-Commit gebautes Vergleichs-Image. „Vorher/nachher" ist damit selbst gefahren, nicht abgeschrieben. Arbeitsbaum am Ende unverändert (`sha256sum -c` auf allen temporär mutierten Dateien).

---

## DoD-Tabelle

| # | DoD-Punkt | Urteil | Beleg |
|---|---|---|---|
| 1 | vierte Spalte nur wenn gefüllt · `Hinweis:` in `--doctor` · `--json`/`--yaml` unverändert · je mit Test | **teilweise** | Verhalten erfüllt: Feldzählung liefert 4 (mit) / 3 (ohne), JSON+YAML-Vergleich pre↔post identisch, `--doctor`-Vergleich zeigt genau **eine** neue Zeile. **Aber „nur wenn gefüllt" hält kein Test:** Mutation `if f.Message != ""` → `if true` ⇒ `make test` **Exit 0**, Suite grün |
| 2 | `hint` in Schema, Lastenheft (Bump/Historie), Spezifikation; leer ⇒ Exit 2, mit Test | **erfüllt** | `hint: ''` und `hint: '   '` ⇒ je Exit 2 mit sprechender Meldung; fehlender Schlüssel ⇒ Exit 1 mit normalem Befund |
| 3 | Vorrang gemessen, als Test | **erfüllt** | Ohne `hint`: `… section-forbidden  verbotenes Muster trifft: TODO`; mit: der Hinweis. Auch für alle drei Zellen-Codes. Mutation `MessageFor` → `if false` macht den Kern-Test **rot** — er beißt |
| 4 | Byte-Identität gemessen, Grenze benannt | **erfüllt** | Grüner Lauf: stdout beidseits leer und hash-gleich, stderr identisch. Befund ohne `message` (`fence-unclosed`): Zeilen-Hash pre==post. Das vorher/nachher-Paar der Commit-Botschaft **wörtlich** reproduziert |
| 5 | ADR begründet die drei Entscheide, im Index | **erfüllt** | ADR-0073 samt Alternativen-Tabelle; Index-Zeile vorhanden |
| 6 | Handbuch: Feld bei den `structure`-Schlüsseln **und** vierte Spalte beim Ausgabeformat | **erfüllt**, mit Rändern | Beide Beispiele gegen ein Probe-Repo nachgefahren und **zeichengleich** reproduziert; Ränder s. A5–A8 |
| 7 | `make gates` grün (Exit explizit); unabhängiger Review; Verifikation | **teilweise** | Exit 0, `593 Datei(en), 0 Befund(e)`. **Kein Review-Artefakt** zu slice-177 in `docs/reviews/` |

## Spec-Konformität

**DC-FA-CLI-004** — Beschreibung und das neue Kriterium „Vierte Spalte" treffen
zu, gemessen mit `links`/`target-missing` genau wie im Kriterium genannt.
`message,omitempty` deckt die neue JSON-Formulierung. Das bestehende Kriterium
„genau zwei Befund-Zeilen" bleibt unberührt.

**DC-FA-CLI-007 / .a** — Schritt-4-Einschub und Neunummerierung konsistent, alle
Rückverweise nachgezogen, kein Rest.

**SPEC-001 / SPEC-002** — der neue Zwei-Zeilen-Block trifft zu; SPEC-002
unverändert, gemessen identisch.

**DC-FA-STRUCT-001** — Prosa ergänzt, aber **kein einziges neues
Akzeptanzkriterium** (s. A9).

## Abweichungen Zusage ↔ Zustand

1. **Ein dritter „hat gar nicht gemessen"-Fall ist nicht ausgenommen.**
   `checkStructureFile` meldet die unlesbare **Einzeldatei** über
   `structureFinding`, also über `MessageFor`. Reproduziert. Widerspricht
   wörtlich ADR-0073, der Lastenheft-Prosa, der Historie-Zeile, der
   §2-Schemazeile, dem Feat-Commit und dem Doc-Kommentar.
2. **`TestHint_OhneErlaeuterungDreiSpalten` kann nicht rot werden** — `t.Logf`
   statt Assertion, und die Prämisse ist falsch: das Szenario liefert real vier
   Spalten. Belegt durch die `if true`-Mutation.
3. **`TestStructureHintGiltNichtFuerLeerlaufendeRegel` deckt eine von zwei
   genannten Grenzen** — sein Kommentar sagt „die **beiden** Befunde".
4. **Fehlzitat ADR-0069 statt ADR-0070** im Slice und in der Plan-Revision.
5. **Zustellung unvollständig** (`BEO-022`, vom Slice selbst zitiert):
   `--print-config` führt `hint` nicht.
6. **Handbuch-JSON/YAML-Beispiele widersprechen jetzt der Spec** — sie zeigen
   `target-missing` **ohne** `message`; die echte Ausgabe trägt es, und
   DC-FA-CLI-004 sagt es seit 0.75.0 zu.
7. **Zwei weitere Text-Beispiele zeigen drei Spalten** (`structure`- und
   `workflows`-Illustration); beide tragen den `not-replayable`-Marker, sind
   also deklariert — aber ausgerechnet `structure` bekommt `hint`.
8. **Die neuen Handbuch-Beispiele sind gemessen, aber nicht verankert** — die
   Form-Anker des Replay-Harness blieben unverändert.
9. **Geänderte Anforderung ohne Akzeptanzkriterium:** DC-FA-STRUCT-001 bekam
   neuen Schlüssel, Exit-2-Rand, Vorrang-Regel und zwei Ausnahmen, aber kein
   neues Kriterium. Zum Vergleich: die neunte Bedingung brachte zehn mit.
10. **Stabilitäts-Vorbehalt wandert nicht mit** — SPEC-001 führt `message` als
    „nicht stabilitätsgarantiert"; DC-FA-CLI-004 nimmt das Feld ins
    Zeilen-**Format** auf, ohne den Vorbehalt zu wiederholen.
11. **Wortlaut-Spannung** in DC-FA-CLI-007.a Schritt 4: „verfasst — modul-eigen
    oder aus `structure[].hint`", während ADR-0073 gerade unterscheidet.
12. **Nicht als Risiko geführter Informationsverlust:** mit `hint` fallen
    unterscheidbare Diagnosen zu einem Text zusammen — gemessen an zwei
    verschiedenen `section-column-missing`-Ursachen.
13. **Präzision der Zählung:** „22 von 31" zählt Dateien mit dem Literal;
    `structure_tableorder.go` setzt seine Meldung über `structureFinding` und
    ist nicht mitgezählt.
14. **`CHANGELOG.md` nicht gepflegt** — die vierte Spalte ist nutzersichtbar.
15. **DoD 7 offen:** kein unabhängiges Review-Artefakt.

## Nebenbefund (keine slice-177-Regression)

Der Formatierungs-Prüfer meldet 37 Dateien — **identisch vor und nach dem
Feat-Commit**; das Lint-Profil führt keinen Formatter. Die Ausrichtung in
`configyaml_test.go` und die doppelte Leerzeile vor `derefString` fallen in
diesen vorbestehenden Bestand.
