# Welle welle-95 — Modul `linelength`: höchstens N Zeichen je Zeile in Markdown-Dateien — Closure-Notiz

> **Zitier-Form** *(bleibt stehen — Norm, kein Ausfüll-Hinweis).* Dieses
> Artefakt friert ein; was es zitiert, bewegt sich weiter. Deshalb: **Kennung,
> nicht Adresse** — `slice-<Kennung>` statt seines Lifecycle-Pfads, `make <target>`
> statt eines Links auf die Sensor-Datei, eine Baseline-Stelle als
> `v6.18.1` · `regelwerk/<datei>.md` §<Abschnitt> statt als Link
> (Baseline-Regelwerk `grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt).

**Welle:** welle-95
**Abschluss:** 2026-11-23
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — Ergebnis, nicht Tätigkeit.

- **slice-289** (Anforderung, Spezifikation, ADR):
  [`DC-FA-LLEN-001`](../../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in)
  im Lastenheft, Verfahren
  [`DC-FA-LLEN-001.a`](../../../../spec/specification.md#dc-fa-llen-001a--zeilenlänge-in-markdown-dateien-linelength)
  mit Schema `linelength.*` unter [`SPEC-005`](../../../../spec/specification.md#spec-005--d-checkyml),
  [ADR-0114](../../adr/0114-linelength-modul.md). Neun Entscheide des Auftraggebers stehen vor dem
  Code fest: Codepunkte der rohen Zeile, drei schaltbare Ausnahmen (Fenced-Code, Zeile mit `|`,
  ein Wort über der Grenze — Variante H), eine leere Menge ist kein Fehler, kein Zeilen-Marker.
  Lastenheft §5 nimmt die Zeilenlänge ausdrücklich vom globalen Lint-Ausschluss aus.
- **slice-290** (Modulkern): Konfigurations-Rand (nur YAML-Ganzzahl, nur `true`/`false`, ein
  eigener Schlüssel ohne Wert ist Exit 2), Messung über die geteilte Lexik ohne eigene Erkennung,
  Grund-Code `line-chars-exceeded` ([`SPEC-112`](../../../../spec/specification.md#4-grund--und-fehler-codes)), Bezugsmenge mit Zahl der Ausnahmen, Registrierung
  in allen Modul-Aufzählungen, Handbuch-Abschnitt 4.28; die Zeilen-Zerlegung von `pinsync` und
  `linelength` ist jetzt eine geteilte Funktion. [ADR-0114](../../adr/0114-linelength-modul.md) ist `Accepted`.
- **Außerhalb der Slices, in dieser Welle:** Phase (9) in `make image-test` — Grenze, Zähleinheit
  und die drei Ausnahmen, je einzeln geschaltet, unter `:ro`, fremdem Benutzer und ohne Netz. Der
  Closure-Trigger verlangt sie.
- **Harness:** das Workflow-Skelett `.claude/commands/implement-slice.md` nennt den
  Architect-Schritt und den Schreiber einer ADR ausdrücklich (Anpassung des Auftraggebers);
  Regel 17 nimmt Modul-Aufzählungen mit Spiegel-Test aus (Steering-Loop-Eintrag unten).

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- **Die Gegenprobe am eigenen Bestand trug eine Entscheidung, nicht nur einen Beleg.** Ihre
  Aufschlüsselung zeigte, dass 44 % der Ausnahmen der zuerst bestätigten Fassung von (c)
  umbrechbare Prosa waren; der Auftraggeber wählte daraufhin Variante H. Am gebauten Image traf das
  Modul die `awk`-Zahlen exakt (1237 Dateien, 513 ausgenommen, 2439 Befunde, Stand `f3dc94da`).
- **Jede Rolle fand, was eine andere nicht sah:** der Architect vor dem Code den Widerspruch von
  Entscheid 8 zu
  [`DC-FA-CONF-002`](../../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope), die
  Reviews die Ränder der Ausnahme (c) am Gegenstand (CommonMark) und eine Plan-Notiz gegen das
  Lastenheft, die Verifikation — weil sie gegen die Spec liest — den Lint-Ausschluss in Lastenheft §5.
- **Brechproben in drei Rollen unabhängig:** Implementer, Reviewer und Verifier zeigten das Rot je
  Testfall selbst (30, 20 und 24 Mutationen in slice-290).
- **Die geteilte Lexik blieb unverändert;** die Kopplungs-Tests binden das neue Modul an die
  bisherigen Konsumenten statt an eine Kopie.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile möglichst mit der Konsequenz, die
daraus schon gezogen wurde.

- **Rollenfolge in slice-289:** der Slice ging nach `in-progress/`, bevor ein Architect-Verdikt
  vorlag, und die ADR war zuerst dem Implementer zugedacht. Der Auftraggeber bemerkte es, bevor
  etwas geschrieben war. Konsequenz: das Workflow-Skelett nennt Architect-Schritt und ADR-Schreiber
  (Anpassung des Auftraggebers); das Verdikt führt die Abweichung als einmaliges akzeptiertes Negativ.
- **Vier Entscheide neu bestätigt, einer im Lastenheft geschärft:** 8 (leere Menge kein Fehler),
  9 (kein Marker, neu), 5 (c) (Variante H) und Lastenheft §5 — jedes Mal, weil eine Rolle am
  Gegenstand einen Widerspruch fand.
- **Eine vorgelegte Zahl war ein Proxy:** für Variante H wurden 467 Ausnahmen genannt, gezählt
  innerhalb der Teilmenge der alten Fassung; direkt gemessen 513. Beleg im Register
  (`zaehlmethode-misst-proxy-statt-gegenstand`).
- **Zweimal stand eine Plan-Festlegung gegen das Lastenheft** (slice-289 Entscheid 9, slice-290
  Notiz zum leeren Wert). Neuer Eintrag im Register, 2×.
- **`make slice-mv` scheiterte am eigenen pre-commit-Hook;** die Übergänge liefen von Hand
  (Register, Sub-Area `tools/harness/`).
- **Die Grenzen der Ausnahme (c) wuchsen über vier Runden,** jede fand die nächste am Gegenstand.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — hier stehen **nur** Beobachtungen, die im Register
3× erreicht haben.

- **Guide** geschärft: Regel 17 nimmt Modul-Aufzählungen, die ein Spiegel-Test gegen die
  Modul-Registrierung hält, von „READMEs/Handbuch-Kopf sind Release-Prep“ aus — sie kommen mit dem
  Feature-Commit, der das Modul registriert; die Kurzform in `AGENTS.md` §5 Zeile 17 zieht mit —
  liegt in `harness/rules/dokumentations-regeln/17-changelog-release-prep.md`.
  Auslöser: `BEO-ALL/regeltext-gegen-spiegel-test` (slice-283, slice-287, slice-290 — 3×).

Gelesen wurde das Register vollständig. Kein weiterer Eintrag dieser Welle über der Schwelle ohne
Ausgang; die drei Klassen, die Belege dieser Welle bekamen und schon verkörpert sind
(`begruendung-traegt-entscheidung-nicht`, `zaehlmethode-misst-proxy-statt-gegenstand`,
`grenzen-liste-wird-als-vollstaendig-gelesen`), bleiben es. Drei ältere Einträge über der Schwelle
mit begründetem Stand ohne Standard-Ausgang (`eigene-menge-gemessen-fremde-behauptet`,
`invented-fourth-closure-outcome`, `pin-bump-mirrors-ungated`) bekamen keinen Beleg und bleiben
unverändert. `BEO-ALL/plan-entscheid-steht-gegen-lastenheft` steht bei 2×.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Zähler steht in [`../observations/`](../observations/).
Was in dieser Welle **3×** erreicht hat, steht oben unter
*Steering-Loop-Einträge*.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur, das
Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine. Offen bleiben, ohne Arbeit der Welle: die Release-Prep des Moduls (CHANGELOG, Handbuch-Kopf,
Image-Tags; dabei in `docs/maintainer/releasing.md` bei der README-Modul-Liste ein Halbsatz, dass sie
nach Regel 17 schon mit dem Feature-Commit vorliegt), die roten Nachtläufe vom 2026-11-22 (Pins, CVEs
im publizierten Image) und die Defekt-Klasse `BEO-ALL/yaml-dekoder-kuerzt-zahlen-still` an den
Nachbarmodulen.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

Beleg des Verifiers: [`docs/reviews/2026-11-23-welle-95-verifikation.md`](../../../reviews/2026-11-23-welle-95-verifikation.md)
(Stand `f98ca836`).

- Beide Slices liegen in `done/`; `make fullbuild` Exit 0 — Image-Hash
  `sha256:7ffebd74950ad2639eae893e2373848a5719c2e2aebe997915979aaf4fbff4c5` (`make gates`,
  `make image-test`, `make bench` mit Median 1074 ms, `make completeness-check` mit 58
  Anforderungen und 0 Waisen, `make verify-closure-notes` darin grün).
- **Ohne den Block nichts Neues:** `make blackbox-probe REF=v0.91.0` byte-identisch über 20
  Vergleiche (Fixtures `ids`, `links`, `sauber`, `targets` und der Eigenlauf, je vier Formen;
  stdout, stderr und Exit getrennt). Das Fixture der Phase (9) ohne Konfiguration: Vorher und
  Nachher 4/4 gleich, Exit 0.
- **Die Entscheide tragen am gebauten Image** (`make image-test`, Phase 9, `max-chars: 40`, `:ro`,
  fremder Benutzer, ohne Netz):

  | Probe | Zeile | Länge | Default | mit Schalter aus |
  |---|---|---|---|---|
  | knapp unter | 2 | 39 | kein Befund | — |
  | genau auf | 3 | 40 | kein Befund | — |
  | knapp über | 4 | 41 | `line-chars-exceeded` | in allen Läufen |
  | Mehrbyte an der Grenze | 5 | 40 Codepunkte, 90 Bytes | kein Befund | — |
  | Mehrbyte über der Grenze | 6 | 41 Codepunkte, 91 Bytes | `line-chars-exceeded` | in allen Läufen |
  | Fenced-Code | 8 | 62 | kein Befund | `exempt-code-blocks: false` ⇒ Befund |
  | Zeile mit `\|` | 10 | 63 | kein Befund | `exempt-tables: false` ⇒ Befund |
  | lange URL (Wort mit 55) | 11 | 66 | kein Befund, zählt in U | `exempt-unbreakable: false` ⇒ Befund |

  Bezugsmenge im Default: `linelength: 1 Datei(en), 7 Zeile(n) geprüft, 1 davon mit einem Wort über
  der Grenze ausgenommen`; Exit 1 in allen vier Läufen.
- Replay-Lauf gegen Golden Set: entfällt, d-check führt keinen (deterministischer Kern,
  [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- Coverage gesamt: 94,90 % bei Schwelle 93 % (`make coverage-gate`), offene Carveouts: keine.

**Trigger-Audit (Schritt 2), je Klasse geprüft:**

- **Carveout:** keine aktiven (`docs/plan/carveouts/` trägt nur `done/` und die README).
- **Bootstrap-aware Gate:** `coverage-gate` kalibriert (93 %, gemessen 94,90 %); kein
  Hochschalt-Trigger, keine neue Stufe. Die Gate-Konfiguration ist über die Welle unverändert.
- **ADR** und **Hard Rule:** vom Architect geprüft,
  [`docs/reviews/2026-11-23-welle-95-trigger-audit.md`](../../../reviews/2026-11-23-welle-95-trigger-audit.md):
  kein Re-Evaluierungs-Trigger von [ADR-0114](../../adr/0114-linelength-modul.md), [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md), [ADR-0113](../../adr/0113-pinsync-modul.md), [ADR-0088](../../adr/0088-file-modul-groessengrenzen.md) oder [ADR-0089](../../adr/0089-structure-max-lines-zwoelfte-bedingung.md) eingetreten,
  keine Folge-ADR; kein Auflösungs-Trigger einer Hard Rule, keine Zeile entfernt.

**Die drei Paarungen:** (a) Anker — der Steering-Loop-Eintrag nennt
`harness/rules/dokumentations-regeln/17-changelog-release-prep.md`; die Datei existiert und trägt
`seit welle-95`. (b) Folge-Slice — keiner genannt. (c) Register — jede genannte Kennung existiert als
Verzeichnis unter `../observations/` und trägt ein nicht leeres `evidence/`.
