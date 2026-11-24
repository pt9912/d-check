# Welle welle-97 — `linelength` nimmt eingerückte Tabellenzeilen aus — Closure-Notiz

> **Zitier-Form** *(bleibt stehen — Norm, kein Ausfüll-Hinweis).* Dieses
> Artefakt friert ein; was es zitiert, bewegt sich weiter. Deshalb: **Kennung,
> nicht Adresse** — `slice-<Kennung>` statt seines Lifecycle-Pfads, `make <target>`
> statt eines Links auf die Sensor-Datei, eine Baseline-Stelle als
> `v6.18.1` · `regelwerk/<datei>.md` §<Abschnitt> statt als Link
> (Baseline-Regelwerk `grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt).

**Welle:** welle-97
**Abschluss:** 2026-11-24
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — Ergebnis, nicht Tätigkeit.

- **slice-293** (Anforderung, Spezifikation, Folge-ADR): Lastenheft
  [`DC-FA-LLEN-001`](../../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in)
  von 0.114.3 auf 0.115.1 — `linelength.exempt-tables` nimmt auch eine `|`-Zeile aus, vor der höchstens
  drei Leerzeichen (U+0020) stehen; sieben Kriterien der Familie (b), die Grenzen mit ihrer Bedingung als
  Sätze im Punkt „Tabellenzeile heißt …“; Spezifikation (Schritt 4, „Nicht Gegenstand“, „Grenzen“,
  Schema-Zeile). [ADR-0116](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md) mit vier
  Entscheiden löst nur den Satz „eingerückte Zeilen sind keine“ aus Entscheidung 5 (b) von
  [ADR-0114](../../adr/0114-linelength-modul.md) ab.
- **slice-294** (Kern): ein neuer benannter Helfer `unbreakableTableRow` in der geteilten Lexik, der
  `tableRowLine` aufruft und um null bis drei Leerzeichen erweitert, außerhalb von Fenced-Code;
  `tableRowLine` bleibt zeichengleich (28 hinzugefügte Zeilen, keine entfernte); `linelength` ruft für
  die Ausnahme (b) nur den Helfer. Kopplungs-Test mit zwei Assertionen, Tests je Kriterium und Grenze,
  Handbuch §4.28 und `--print-config`-Vorlage nachgezogen.
- **Profil und ADR-Übergang:** [ADR-0116](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md) ist
  `Accepted`, [ADR-0114](../../adr/0114-linelength-modul.md) trägt eine Zeile in `Geschichte` und im
  Index den Vermerk „teil-superseded“, ihr Status bleibt `Accepted`.
- **Harness:** Regel 13 in `AGENTS.md` §5 nimmt die Bedingung einer wiedergegebenen Grenze mit
  (Steering-Loop-Eintrag unten).
- **Anlass:** Issue #20 — Tabellen in Listenpunkten, deren Zeilen sich nicht umbrechen lassen.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- **Der kleinste Zuschnitt hielt.** Aus der Lehre von welle-96 wurde zuerst die kleinste Antwort auf die
  offenen Fragen des Issues ausgeschrieben und nur sie geplant: kein neuer Schlüssel, `> |` bleibt
  gemessen, höchstens drei Leerzeichen. Es blieben sieben Kriterien und vier Entscheide; jeder Slice
  brauchte **eine** Review-Runde, die erste ohne HIGH. Dem Auftraggeber lag nur eine Frage vor, die er
  entscheiden musste.
- **Messen vor dem Schreiben.** Die Fixture am Image legte den Ist-Zustand fest, Reviewer und Verifier
  fuhren sie unabhängig nach.
- **Die Mutation in drei Rollen** (10 · 12 · 19) deckte den äquivalenten Mutanten auf (Doppelaufruf
  von Helfer und `tableRowLine`) und zeigte, dass erst die zweite Kopplungs-Assertion ihn zusammen mit
  „Helfer ruft `tableRowLine` nicht“ rot macht.
- **Der Architect entschied Frage 4 am Code:** `tableRowLine` nicht weiten, sonst rutschten `targets`,
  `planning.waves` und `structure` gegen ihre Zusage „Spalte 0“.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile möglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- **Der Auftraggeber bestätigte die enge Fassung gegen den Wortlaut des Issues** („2 bis 5
  Leerzeichen“ gegen „höchstens drei“), weil die geteilte Lexik keinen Listenkontext kennt. Eine Tabelle
  unter `10. ` oder in einem verschachtelten Listenpunkt bleibt ein Befund; ob der Bestand des Melders
  solche Zeilen trägt, ist nicht gemessen. Konsequenz: Beobachtung
  [`BEO-ALL/enge-einrueckungs-grenze-laesst-listen-tabellen-stehen`](../observations/BEO-ALL/enge-einrueckungs-grenze-laesst-listen-tabellen-stehen/observation.md).
- **Das Verdikt des Architects fand zwei Plan-Fehler vor dem Code:** der Plan sagte „zusätzlich zu
  `tableRowLine`“ (damit bliebe die Mutation zu Kriterium 6 grün), und die Probe mit `REF=v0.92.0`
  konnte nie byte-identisch enden, weil das Vorher-Image das Profil seit welle-96 ablehnt; gemessen wird
  gegen den Self-Close von welle-96.
- **Ein Satz im Lastenheft ging über den Gegenstand hinaus** („die übrigen Module lesen eine
  eingerückte Zeile weiter nicht als Tabellenzeile“): der RTM-Leser und der Beobachtungs-Leser in
  `planning` trimmen die Zeile. Gekürzt auf die drei Module, Lastenheft 0.115.1 mit eigener
  Historie-Zeile.
- **Die Blackbox-Probe fährt `linelength` nicht.** Sie fährt im Eigenlauf `structure`, weder
  `planning.waves` noch `targets`; der Vergleich des Moduls ist ein eigener (unten).
- **Zwei äquivalente Mutanten blieben grün,** die Bedingung `indent >= 1` im Helfer ist tot, weil
  `tableRowLine` für Einrückung null schon antwortet (Verifikation V-1).
- **Ein Kommentar begründete mehr, als er trägt** (R1 F-1): Tab, U+00A0 und `>` mit „die Lexik hat
  keinen Listenkontext“, was nur für vier und mehr Leerzeichen gilt. Das war der dritte Beleg der
  Klasse `grenz-aussage-ohne-ihre-bedingung` (Steering-Loop-Eintrag).

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<KUERZEL>/<slug>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks).

- **Hard Rule (Dokumentations-Regel 13)** ergänzt: wer eine Grenze an einer zweiten Stelle wiedergibt,
  nimmt ihre Bedingung mit — in ein Kriterium, eine Historie-Zeile, die Spezifikation, eine ADR, einen
  Kommentar oder einen Testfall; für eine Begründung gilt dasselbe — liegt in `harness/rules/dokumentations-regeln/13-grenzen-gegen-gegenstand-pruefen.md`.
  Auslöser:
  [`BEO-ALL/grenz-aussage-ohne-ihre-bedingung`](../observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/observation.md)
  (slice-291, slice-293, slice-294 — 3×). Als Erweiterung der bestehenden Regel statt einer neuen Regel,
  weil dieselbe Frage nicht zwei Orte haben soll. Auflösungs-Trigger: permanent.

Keine weitere Beobachtung erreicht in dieser Welle die Schwelle 3× ohne Ausgang.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Zähler steht in [`../observations/`](../observations/).
Was in dieser Welle **3×** erreicht hat, steht oben unter
*Steering-Loop-Einträge*.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

- Keine neuen Slices. Die Listen-Erkennung für Tabellen unter längeren Listenmarken wäre ein eigener
  Vorgang, falls der Bestand des Melders aus Issue #20 nach der Hebung noch Befunde dieser Form trägt.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

Beleg des Verifiers: [`docs/reviews/2026-11-24-welle-97-verifikation.md`](../../../reviews/2026-11-24-welle-97-verifikation.md);
Trigger-Audit des Architects: [`docs/reviews/2026-11-24-welle-97-trigger-audit.md`](../../../reviews/2026-11-24-welle-97-trigger-audit.md).

- Beide Slices liegen in `done/`; `make fullbuild` Exit 0 auf `5c8ded99` — Image-Hash
  `sha256:dee681826eea27a2163f4e4b3416e65ceef5140a5c740f7739fec9067fb9aac9` (Coverage 94,90 % bei
  Schwelle 93 %, `make bench` mit Median 1000 ms, `make completeness-check` mit 58 Anforderungen und 0
  Waisen, `make verify-closure-notes` darin grün). Danach kamen nur noch Dokumente hinzu (`192cf345`:
  `AGENTS.md`, ein Regel-Volltext, ein Register-Zustand, der Trigger-Audit-Report); `make gates`,
  `make image-test`, `make coverage-gate`, `make completeness-check` und `make verify-closure-notes`
  liefen dort erneut mit Exit 0. `make fullbuild` selbst wurde auf diesem letzten Stand nicht
  wiederholt, und der gedruckte Image-Hash ist nach späteren Läufen per `docker image inspect` nicht
  mehr abrufbar, weil `d-check:latest` neu getaggt wurde.
- **Ohne die Änderung nichts Neues.** `make blackbox-probe REF=69c53e0f` (der Self-Close von
  welle-96): byte-identisch über 20 Vergleiche; sie fährt `linelength` nicht. Deshalb der eigene
  Vergleich des Moduls gegen ein Image aus `69c53e0f`: 14 Fixtures × 7 Schalter-Stellungen ×
  `max-chars` 20/40/80/120 × Text, `--json`, `--yaml` = 1176 Vergleiche, je stdout, stderr und Exit
  getrennt. 996 gleich, 180 abweichend. Aufgeteilt: 672 gleich über acht Fixtures ohne eingerückte `|`-Zeile, Kopien von
  `docs/user/`, `docs/maintainer/`, `harness/rules/` und `docs/plan/adr/`, und 252 gleich über Zaun-Rand,
  null und vier Leerzeichen; 252 über ein bis drei Leerzeichen, davon 180 abweichend — genau die
  Zellen mit `exempt-tables` an; die 72 mit `exempt-tables: false` sind gleich. Die Abweichung ist
  immer dieselbe: Befunde der eingerückten Zeilen fallen weg, die Bezugsmenge sinkt (im Beispiel von 7
  auf 4), der Exit wechselt von 1 auf 0.
- **Der Anlass trägt am gebauten Image** (`--network none`, `:ro`, `max-chars: 30`, überlange Datenzeile
  in Zeile 5):

  | Probe | Schalter | Image | Ergebnis |
  |---|---|---|---|
  | Tabelle mit 3 Leerzeichen | Default | Nachher | Exit 0, 0 Befunde, 4 Zeilen geprüft |
  | Tabelle mit 3 Leerzeichen | `exempt-tables: false` | Nachher | Exit 1, 1 Befund (Zeile 5), 7 Zeilen |
  | Tabelle mit 3 Leerzeichen | Default | Vorher (`69c53e0f`) | Exit 1, 1 Befund |
  | 4 Leerzeichen | Default | Nachher | Exit 1, 1 Befund |
  | Tab davor | Default | Nachher | Exit 1, 1 Befund |
  | `> \|` und `>   \|` | Default | Nachher | Exit 1, je 1 Befund |
  | 0 bzw. 1 Leerzeichen | Default | Nachher | Exit 0, 0 Befunde |
  | 3 Leerzeichen, im Zaun | `exempt-code-blocks: false` | Nachher | Exit 0, 0 Befunde |
  | 3 Leerzeichen, ein Wort über der Grenze | `exempt-tables: false` | Nachher | Exit 0, `1 davon mit einem Wort über der Grenze ausgenommen` (bestehende Regel) |
  | dasselbe | `exempt-tables: false`, `exempt-unbreakable: false` | Nachher | Exit 1, 1 Befund |
  | 3 Leerzeichen | Default, `--json` | Nachher / Vorher | Nachher `findings: []`, Exit 0; Vorher ein Befund `line-chars-exceeded`, Exit 1 |

  Der Verifier fuhr 24 Proben; die Tabelle nennt die tragenden.
- **Trigger-Audit:** ADR: ein Trigger eingetreten, beantwortet (Trigger (b) von
  [ADR-0114](../../adr/0114-linelength-modul.md), durch [ADR-0116](../../adr/0116-linelength-eingerueckte-tabellenzeilen.md)),
  sonst 0 offen; Hard Rule 0 offen. Carveouts: keine offen. Bootstrap-aware Gates: keines berührt
  (`.d-check.yml`, `Makefile`, `.golangci.yml`, `tools/coverage-gate.sh` ohne Diff über die Welle).
