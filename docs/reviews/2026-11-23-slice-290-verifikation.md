# Verifikation slice-290 — Modulkern `linelength`

**Rolle:** Verifier (Modul 8/11) · **Gegenstand:** Plan und DoD aus
`docs/plan/planning/done/slice-290-linelength-modulkern.md`, geprüft auf dem Stand
`45d651df` mit sauberem Arbeitsbaum · **Bezug:**
[`DC-FA-LLEN-001`](../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in),
[`DC-FA-LLEN-001.a`](../../spec/spezifikation.md#dc-fa-llen-001a--zeilenlänge-in-markdown-dateien-linelength),
[`SPEC-005`](../../spec/spezifikation.md#spec-005--d-checkyml),
`SPEC-112` ([§4](../../spec/spezifikation.md#4-grund--und-fehler-codes)),
[ADR-0114](../plan/adr/0114-linelength-modul.md) (`Proposed`),
[ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) · Datum 2026-11-23.
Code-Diff: `git diff c3c08c7b..HEAD` ohne `docs/reviews/`. Die Code-Reviews
[R1](2026-11-23-slice-290-code-r1.md) und [R2](2026-11-23-slice-290-code-r2.md) habe ich gelesen
und nicht neu bewertet. Jeder Beleg unten stammt aus einem Lauf, den ich selbst gefahren habe.
Übernommen ist keiner: Der Implementer schreibt in seinen Commit-Botschaften, jeder neue Test sei
mit einer rot färbenden Mutation gelaufen. Diese Ausgabe liegt nirgends vor, deshalb habe ich die
Mutationen selbst gefahren (§Bewusstes Brechen).

## DoD-Punkte

| DoD | Urteil | Beleg |
|---|---|---|
| 1 — Konfigurations-Rand | bestätigt | `configyaml.go` `applyLinelength` mit `linelengthInt`/`linelengthBool`/`linelengthValue`, dazu `validateSegmentGlobs` und `applyScopes` (Zeile `{"linelength", …}`); `checkLinelengthConfig` vor dem Lauf. Tests: `TestDecode_LinelengthEingabefehler` (14 Fälle), `…SchluesselOhneWert`, `…Alias`, `…ScopeUndExemptPathsOhneWert`, `TestLinelengthOhneGrenze`, `TestCLILinelengthEnable`. Am Image: ohne Grenze, `max-chars: 0`, `"40"`, Schalter ohne Wert und unbekannter Schlüssel enden je mit Exit 2 und Klartext auf stderr, ohne stdout (Tabelle unten). `--print-config` nennt `linelength` in „Verfügbar“ und führt den Block (Z. 299 ff. der Vorlage); `TestPrintConfigLinelengthBlockIstGueltig` dekodiert ihn. `--enable linelength` läuft am Image. Rot-Belege: M11–M14. |
| 2 — Regel | bestätigt, mit V-1 und V-2 | Jedes Akzeptanzkriterium hat einen Test (Zuordnung unten). Für jedes der 17 Kriterien ist ein Rot ohne die Regel gezeigt, durch 24 Mutationen am Code, alle rot. Die Bedingung beim Reihenfolge-Test steht in V-2. Den Kopplungs-Test nach ADR-0054 Entscheidung 4 gibt es zweimal: `TestTabellenzeilenFrageHatEineAntwort` mit `targets`, `planning.waves`, `structure` und `linelength`, `TestFencedCodeFrageHatEineAntwort` mit `targets`, `structure`, `planning.waves`, `planning`, `vcs`, `vcs.status` und `linelength`. Beide werden rot, sobald `linelength` eine eigene Antwort gibt (M06, M22). |
| 3 — Spiegel | bestätigt | `AllReasons()` und Klartext in `diagnose.go`; die Zeile `SPEC-112` in §4; `DC-FA-CLI-002` und Glossar im Lastenheft 0.113.2 mit Historie-Zeile; READMEs, Handbuch §4.28/§5/§6, `operations.md`, `abdeckung-tests.md`. Dass die Spiegel-Tests das neue Modul **fangen**, zeigt je ein Entfernen: S1 README ⇒ `TestModulRegistrySpiegel_ReadmeEn`, S2 §4-Zeile ⇒ `TestAllReasonsDeckungGegenSpezifikationGrundCodes`, S3 Klartext ⇒ `TestReasonTextDeckungGegenAllReasons`, S4 `DC-FA-CLI-002` ⇒ `TestModulRegistrySpiegel_LastenheftCLI002`. Am Image erscheint der Klartext unter `--doctor`. |
| 4 — `make gates` grün | bestätigt | Selbst gefahren auf `45d651df`, Exit 0: `verify ok (54 Dateien, vollständig)`; workflows und doc-check je `0 Befund(e)` (1254 bzw. 1260 Dateien); `0 issues.`; alle Test-Pakete `ok`; a-check `gesamt: 0 Befund(e)`; `coverage-gate: OK — Coverage 94.90% erfüllt Schwelle 93%` (`linelength.go`: alle fünf Funktionen 100.0 %); Semgrep `Ran 55 rules on 75 files: 0 findings.`; gate-consistency und planning-check je `1254 Datei(en) geprüft, 0 Befund(e)`. Schlusszeile: `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`. |

`make completeness-check`, selbst gefahren: `58 Anforderung(en), 0 Waise(n).`, Exit 0. Die Zeile
`DC-FA-LLEN-001 | … | ADR-0114 | slice-289, slice-290 | Tests | ok`. `make verify-closure-notes`
habe ich nicht gefahren, weil der Slice noch nicht in `done/` liegt.

## Akzeptanzkriterium → Test → Rot-Beleg

Testnamen ohne Paket stehen in `internal/hexagon/core/rules/linelength_test.go`. `TestDecode_*`
steht in `configyaml/linelength_test.go`, `TestCLI*`/`TestPrintConfig*` in
`cli_acceptance_test.go`, die beiden `…FrageHatEineAntwort` in `lexikon_kopplung_test.go`. Die
Mutationen stehen in der Liste darunter.

| Kriterium | Test | Rot ohne die Regel |
|---|---|---|
| Happy Path | `TestLinelengthGrenzeGenauAufDerSchwelle` (Zeile 1, 80 Zeichen ⇒ kein Befund); Exit 0 am Image (Probe P10) | M01: `docs/a.md: Befund-Zeilen [1 2], erwartet [2]`, Meldung `Zeile hat 80 Zeichen, erlaubt sind 80` |
| Boundary (Grenzwert) | `TestLinelengthGrenzeGenauAufDerSchwelle` (Ziel `max-chars=80`, Meldung 81/80); Exit 1: `TestCLILinelengthEnable` | M01 (wie oben) |
| Boundary (Zähleinheit) | `TestLinelengthZaehltCodepunkte` (Umlaut/Emoji, Tab, ungültiges UTF-8), `…ZeilenendeOhneCR`, `…NurEinCRFaelltWeg`, `TestTextLinesZaehltWieFile` | M08 (Bytes statt Codepunkte) ⇒ `ZaehltCodepunkte`; M09 (alle `\r` weg) ⇒ `NurEinCRFaelltWeg`, `TextLinesZaehltWieFile`, `TestPinsyncZeilenRandformen`; M10 (kein `\r` weg) ⇒ zusätzlich `ZeilenendeOhneCR` |
| Ausnahme (a), an und aus | `TestLinelengthAusnahmeFencedCode` | M05 (Schalter ignoriert): `Befund-Zeilen [], erwartet [2 3 6 9]` |
| Ausnahme (a), Grenzen | `TestLinelengthFenceGrenzenDerLexik`, `…OffenerZaunNimmtDenRestAus`, `TestFencedCodeFrageHatEineAntwort` | M22 (eigene Fence-Erkennung, nur die Zaun-Zeile): `FenceGrenzenDerLexik`, `OffenerZaun…`, `AusnahmeFencedCode`, Kopplung in 4 Fällen |
| Ausnahme (b), an und aus | `TestLinelengthAusnahmeTabellenzeile` | M07 (Schalter ignoriert): `Befund-Zeilen [], erwartet [1]` |
| Ausnahme (b), Grenzen | `TestLinelengthTabellenzeileGrenzen`, `TestTabellenzeilenFrageHatEineAntwort` | M06 (eigene Erkennung mit `TrimLeft`): `TabellenzeileGrenzen`, Kopplung `eingerueckt` und `im_Fence`; M05 trifft den Fall „(a) aus, (b) an“ (`docs/b.md: … erwartet [2]`) |
| Ausnahme (c), an und aus | `TestLinelengthEinWortUeberDerGrenze`, `…AusnahmeWortUeberDerGrenze` | M23 (Schalter ignoriert): `Befund-Zeilen [], erwartet [1]` |
| Ausnahme (c), Ränder | `TestLinelengthAusnahmeWortUeberDerGrenze`, `…WortTrennerSindSpaceUndTab`, `…EinrueckungUndUeberschrift`, `…FormPraefix` | M02 (`>=` statt `>`): Wort genau 80 ausgenommen, `erwartet [2 4]`; M03 (Vorfassung „kein Leerraum ab der Grenze“): 9 Tests, darunter `FormPraefix` und `WortTrennerSindSpaceUndTab`; M04 (Tab trennt nicht) ⇒ `WortTrennerSindSpaceUndTab` |
| Überschrift | `TestLinelengthEinrueckungUndUeberschrift` (Zeilen 2 und 3) | M02 ⇒ `Befund-Zeilen [2], erwartet [1 2 3]` |
| Kein Zeilen-Marker | `TestLinelengthKeinZeilenMarker` | M21 (Marker honoriert): `Befund-Zeilen [], erwartet [1]`; M02 trifft den 14er-Rand |
| Eingabemenge | `TestLinelengthEingabemenge` (`exempt-paths`, `scope`, `.txt`, `.MD`, `.markdown`) | M15 (`exempt-paths` wirkungslos): Befund aus `docs/frei/b.md` |
| Leere Menge | `TestLinelengthLeereMenge` | M19 (Zeile nur bei D > 0): `Bezugsmengen-Zeile fehlt: []` |
| Bezugsmenge | `TestLinelengthBezugsmenge` (mit und ohne `\n`), `…LeereDatei`, Kanal: `TestCLILinelengthEnable` (stderr; `--json`/`--yaml` in `notes`) | M16 (U nicht gezählt) ⇒ `Bezugsmenge`, 3 weitere; M17 (Schlusszeile hinter `\n` gezählt) ⇒ `Bezugsmenge`, `TestCLILinelengthEnable`, `TextLinesZaehltWieFile` |
| Negative (Eingabefehler) | `TestDecode_LinelengthEingabefehler`, `…SchluesselOhneWert`, `…Alias`, `TestLinelengthOhneGrenze`, `TestCLILinelengthEnable` | M11 (Grenzen-Prüfung aus): `OhneGrenze` und `CLILinelengthEnable` rot, als nil-Panic, siehe V-3; M12 (`< 0`): `Grenze_0`, `err = <nil>`; M13 (Tag `!!bool` nicht geprüft): `Schalter_yes`; M14 (`null` nicht erkannt): 5 Fälle `SchluesselOhneWert`, `Alias/null` |
| Reihenfolge und Determinismus | `TestLinelengthReihenfolgeUndDeterminismus`, `TestCLILinelengthByteIdentisch` | M24 (Sortierung aus **und** je Datei rückwärts): `Reihenfolge [docs/a.md:3 docs/a.md:1 …], erwartet [docs/a.md:1 docs/a.md:3 …]`. M20 allein bleibt für diesen Test grün, siehe V-2 |
| Default byte-identisch | `TestLinelengthInaktivIstByteIdentisch`, `TestCLILinelengthByteIdentisch`; gegen das Vorher-Image: `make blackbox-probe` | M18 (Zeile auch bei inaktivem Modul): `CLILinelengthByteIdentisch`, `OhneGrenze`, zwei `pinsync`-Tests |

**Die Mutationen.** Jede Mutation habe ich per `sed` gesetzt und danach `make test` gefahren.
Danach `git checkout` der Datei. Nach der Reihe meldete `git diff --exit-code` keinen Unterschied.
Jede Zeile nennt Exit 2 von `make test` und die rot gewordenen Tests:

- M01 `length <= limit` ⇒ `<` — `GrenzeGenauAufDerSchwelle`, `ZaehltCodepunkte`, `ZeilenendeOhneCR`, `NurEinCRFaelltWeg`
- M02 `longestWord(line) > limit` ⇒ `>=` — `AusnahmeWortUeberDerGrenze`, `EinrueckungUndUeberschrift`, `KeinZeilenMarker`
- M03 Vorfassung von (c): ausgenommen, wenn ab Position `max-chars + 1` kein Space/Tab steht — 9 Tests
- M04 Tab kein Worttrenner — `WortTrennerSindSpaceUndTab`, `ZaehltCodepunkte`
- M05 `exempt-code-blocks` ignoriert — `AusnahmeFencedCode`, `TabellenzeileGrenzen`, Tabellen-Kopplung `im_Fence`
- M06 eigene Tabellen-Erkennung — `TabellenzeileGrenzen`, Tabellen-Kopplung `eingerueckt`, `im_Fence`
- M07 `exempt-tables` ignoriert — `AusnahmeTabellenzeile`, Fence-Kopplung in 3 Fällen
- M08 `len(line)` statt Codepunkte — `ZaehltCodepunkte`
- M09 `TrimRight` statt `TrimSuffix` (`\r`) in `textLines` — `NurEinCRFaelltWeg`, `TextLinesZaehltWieFile`, `TestPinsyncZeilenRandformen`
- M10 kein `\r` entfernt — zusätzlich `ZeilenendeOhneCR`, `TestPinsyncZeilenRandformen/CRLF`
- M11 `if cfg.MaxChars == nil` ⇒ `if false` — `OhneGrenze`, `TestCLILinelengthEnable` (Panic)
- M12 `*maxChars < 1` ⇒ `< 0` — `TestDecode_LinelengthEingabefehler/Grenze_0`
- M13 Tag `!!bool` nicht geprüft — `…/Schalter_yes`
- M14 `!!null` nicht erkannt — `TestDecode_LinelengthSchluesselOhneWert` (5), `TestDecode_LinelengthAlias/null`
- M15 `exempt-paths` wirkungslos — `Eingabemenge`
- M16 `count.Unbreakable++` entfernt — `AusnahmeWortUeberDerGrenze`, `Bezugsmenge`, `EinWortUeberDerGrenze`, `KeinZeilenMarker`
- M17 `textLines` ohne `TrimSuffix(…, "\n")` — `Bezugsmenge`, `AusnahmeWortUeberDerGrenze`, `Eingabemenge`, `TextLinesZaehltWieFile`, `TestCLILinelengthEnable`
- M18 Bezugsmengen-Zeile auch bei inaktivem Modul — `TestCLILinelengthByteIdentisch`, `OhneGrenze`, `TestPinsyncBlockBeiInaktivemModulIstInert`, `TestPinsyncQuellenAusserhalbDerScanMenge`
- M19 Bezugsmengen-Zeile nur bei D > 0 — `LeereMenge`
- M20 `model.SortFindings` aus — 8 Tests anderer Module, **nicht** `ReihenfolgeUndDeterminismus` (V-2)
- M21 `d-check:ignore` honoriert — `KeinZeilenMarker`
- M22 eigene Fence-Erkennung — `AusnahmeFencedCode`, `FenceGrenzenDerLexik`, `OffenerZaunNimmtDenRestAus`, Fence-Kopplung in 4 Fällen
- M23 `exempt-unbreakable` ignoriert — `EinWortUeberDerGrenze`
- M24 M20 plus Befunde je Datei rückwärts — `ReihenfolgeUndDeterminismus` und 16 weitere
- S1–S4 Spiegel (DoD 3): je ein Spiegel-Test rot, siehe oben

## Plan gegen Code

`git diff c3c08c7b..HEAD --name-only` ohne `docs/reviews/` ergibt 21 Dateien. Jede ist im Plan
geführt: in der Tabelle von §3 oder in einer Plan-Änderungs-Notiz.

| Datei | Plan-Stelle |
|---|---|
| `model/config.go`, `configyaml.go`, `rules/linelength.go`, `rules/run.go`, `app/diagnose.go`, `cli/config_template.go`, `spec/spezifikation.md`, `spec/lastenheft.md`, `README.md`, `README.de.md`, `benutzerhandbuch.md`, `operations.md`, `abdeckung-tests.md`, `rules/linelength_test.go`, `configyaml/linelength_test.go` | Tabelle §3 |
| `cli/cli_acceptance_test.go`, `rules/lexikon_kopplung_test.go` | Notiz „Spiegel-Suche vor dem ersten Edit“ |
| `rules/codepaths.go`, `rules/pinsync.go` (`textLines`) | Notiz „Fix nach Code-Review R1“, F-2 |
| `docs/plan/adr/0114-linelength-modul.md` | Die Notiz nennt die ADR „unberührt (Architect)“ für den Implementer. Die zwei Änderungen sind die `docs(adr)`-Commits `18214ff1` und `faf38a83`, Verdikte des Architect zu Entscheid 4 nach R1/R2, und die Notizen verweisen auf sie. Das ist keine Abweichung des Implementers. Der Status bleibt `Proposed`. |
| der Slice-Plan selbst | Planung |

Die Abgrenzung aus §1 ist eingehalten: `.d-check.yml`, `CHANGELOG.md`, Handbuch-Kopf und §11 sowie
`packaging/dockerhub/overview.md` sind unberührt. Es gibt keine eigene Fence- oder
Tabellen-Erkennung, die Regel ruft `proseLineSet`/`tableRowLine`. `--fix-patch` und
`--suggest-config` sind nicht erweitert. Die Gegenprobe aus der Notiz „zwei Festlegungen“ habe ich
nachgefahren. `git archive f3dc94da`, der damalige `scan`-Block, `modules: []` und
`linelength: {max-chars: 120}` am gebauten Image mit `--network none` ergeben
`linelength: 1237 Datei(en), 79250 Zeile(n) geprüft, 513 davon mit einem Wort über der Grenze
ausgenommen` und `2439 Befund(e)`. Das sind genau die Zahlen aus dem Plan von slice-289.

## Verhalten am gebauten Image

Gebaut mit `make build` aus `45d651df` (`d-check:latest`, erstellt 2026-11-23T09:29:32). Das
Fixture hat eine Datei `probe.md` mit elf Zeilen und `linelength: {max-chars: 40}`. Gestartet
wurde mit `docker run --rm --network none -v <fixture>:/repo:ro d-check:latest --enable
linelength`. Die Zeichenzahlen in Codepunkten habe ich mit `awk` unter `C.UTF-8` gezählt.

| Probe | Zeile im Fixture | Ergebnis | zugesagt |
|---|---|---|---|
| P1 knapp unter | Z. 2, 39 Zeichen | kein Befund | ✓ |
| P2 genau auf | Z. 3, 40 Zeichen | kein Befund | ✓ (Grenze befundfrei) |
| P3 knapp über | Z. 4, 41 Zeichen, kein Wort über 40 | `probe.md:4  max-chars=40  line-chars-exceeded  Zeile hat 41 Zeichen, erlaubt sind 40` | ✓ |
| P4 Mehrbyte an der Grenze | Z. 5, `äö😀` ×10, 40 Codepunkte, 90 Bytes | kein Befund | ✓ (Codepunkte, nicht Bytes) |
| P5 Mehrbyte über der Grenze | Z. 6, 41 Codepunkte, 91 Bytes | Befund, `41 … 40` | ✓ |
| P6 Fenced-Code | Z. 8, 60 Zeichen zwischen Backtick-Zäunen (Z. 7/9) | kein Befund; mit `exempt-code-blocks: false` Befund `probe.md:8 … 60 …`, geprüft 10 Zeilen | ✓ an und aus |
| P7 Tabellenzeile | Z. 10, `\|`-Zeile mit 64 Zeichen | kein Befund; mit `exempt-tables: false` Befund `probe.md:10 … 64 …` | ✓ an und aus |
| P8 lange URL | Z. 11, 66 Zeichen, Wort mit 55 | kein Befund, zählt in U; mit `exempt-unbreakable: false` Befund `probe.md:11 … 66 …`, U = 0 | ✓ an und aus |
| P9 Bezugsmenge | ganze Datei | stderr `linelength: 1 Datei(en), 7 Zeile(n) geprüft, 1 davon mit einem Wort über der Grenze ausgenommen`; unter `--json` in `summary.notes`; Exit 1 | ✓ (11 − 3 Zaun/Code − 1 Tabelle = 7) |
| P10 Happy Path | `max-chars: 66`, alle Ausnahmen aus | `0 Befund(e)`, `11 Zeile(n) geprüft`, Exit 0 | ✓ |
| P11 Determinismus | zweimal P9 | stdout, stderr und Exit byte-gleich | ✓ |
| P12 inert | Block ohne `--enable` | `0 Befund(e)`, keine Bezugsmengen-Zeile, Exit 0 | ✓ |
| P13 ohne Grenze | `exempt-tables: true`, `--enable` | `linelength.max-chars fehlt …`, Exit 2 | ✓ |
| P14 Grenze 0 / `"40"` | ohne bzw. mit `--enable` | `max-chars 0 muss >= 1 sein` bzw. `"40" ist keine ganze Zahl`, je Exit 2 | ✓ |
| P15 Schalter ohne Wert | `exempt-tables:` ohne `--enable` | `steht ohne Wert (YAML-null)`, Exit 2 | ✓ |
| P16 unbekannter Schlüssel | `max-lines: 3` im Block | `line 3: field max-lines not found`, Exit 2 | ✓ (V-4) |
| P17 unlesbare Datei | `zu.md` mit Modus 000 | `zu.md nicht lesbar: … permission denied`, Exit 2 | ✓ |
| P18 Diagnose | `--doctor` | Klartext „Eine Zeile einer Markdown-Datei trägt mehr Zeichen …“ je Befund | ✓ |

Die Zusage hält: Am Image entstehen genau die Befunde P3 und P5. Fenced-Code, Tabellenzeile und
URL-Zeile sind im Default ausgenommen und mit dem jeweiligen Schalter einzeln meldbar.

## Byte-Identität ohne `linelength`-Konfiguration

`make blackbox-probe REF=v0.91.0`, Exit 0: `byte-identisch über 20 Vergleiche (Vorher 58c9486d,
stdout/stderr/Exit getrennt, Kanarienlauf vorher und nachher)`. Verglichen wurden die Fixtures
`ids`, `links`, `sauber`, `targets` und der Eigenlauf `repo`, je in den Formen Default, `--json`,
`--yaml` und `--doctor`. Das eigene Profil enthält keinen `linelength`-Block
(`grep -c linelength .d-check.yml` ⇒ 0). Zusätzlich habe ich das Probe-Fixture ohne
Konfigurationsdatei gegen das Vorher-Image `d-check:probe-vorher` und das Nachher-Image gefahren,
in denselben vier Formen: `4/4 gleich` (Exit 0). Einen gültigen Block **ohne** `--enable` kann das
Vorher-Image nicht lesen, es endet mit Exit 2 (Grenze 3 von `harness/sensors/blackbox-probe.md`).
Diese Hälfte belegen `TestLinelengthInaktivIstByteIdentisch`, `TestCLILinelengthByteIdentisch` und
P12.

## Findings

- **HIGH:** keine. **MEDIUM:** keine.
- **LOW — V-1, Rot-Belege nur behauptet.** DoD 2 verlangt „für jeden Test ist ein Rot ohne die Regel
  gezeigt“. Risiko 3 in §6 nennt dafür eine Zuordnung im Plan-Nachtrag als Bedingung für den
  Ausgang *entfallen*. Der Plan trägt diese Zuordnung nicht. Die Commit-Botschaften behaupten die
  Mutationen („Jeder neue Test mit einer rot faerbenden Mutation gelaufen“), ihre Ausgabe liegt
  nicht vor. Mit den Mutationen M01–M24 dieses Berichts ist das Rot für jedes Kriterium jetzt
  gezeigt. Der DoD-Punkt ist damit im Ergebnis erfüllt, die Lücke liegt in der Form des Belegs.
  **Träger:** Planner bei der Closure (Ausgang von Risiko 3). Die Zuordnung oben kann er zitieren,
  statt sie nachzutragen.
- **LOW — V-2, Reihenfolge-Test fängt den Ausfall der Sortierung allein nicht.**
  `TestLinelengthReihenfolgeUndDeterminismus` bleibt grün, wenn `model.SortFindings` entfällt
  (M20). Die Reihenfolge entsteht schon aus der sortierten Datei-Liste und den Befunden je Datei in
  Zeilen-Reihenfolge, die Sortierung ist eine zweite Sicherung. Rot wird der Test erst, wenn beide
  Sicherungen fallen (M24). Das Fehlen der Sortierung allein fangen acht Tests anderer Module
  (`TestCodepathsModul`, sechs `TestHostpaths*`, `TestVersionsAusgabeFolgtDerSortierung`). Die
  Zusage hält also repo-weit, nur der benannte Test kann sie allein nicht belegen. Die Hälfte
  „Determinismus“ lässt sich per Mutation nicht brechen: Das Modul hat keine nicht-deterministische
  Quelle, keine Map-Iteration über die Ausgabe und keine Zeit. Ihr Beleg sind die zehn Läufe im Test
  und P11. **Träger:** Planner, als Lerneintrag oder Annahme bei der Closure. Ein Fix ist nicht
  nötig.
- **INFO — V-3.** Das Rot von M11 entsteht als nil-Panic in `CheckLinelength` (`*cfg.MaxChars`), nicht
  als falscher Klartext. Genau das wäre der reale Vorzustand ohne die Prüfung. Die Behauptung des
  Tests (`err` enthält `linelength.max-chars`) finge auch eine Prüfung mit falscher Meldung. Kein
  Träger.
- **INFO — V-4.** Der unbekannte Schlüssel meldet `line 3: field max-lines not found`. Schritt 9
  verlangt „Klartext mit Schlüssel und Ursache“, beides steht in der Meldung, der Block
  `linelength` dagegen nicht. Das ist die strikte Dekodierung, die jedes Modul teilt. Keine
  Abweichung von `DC-FA-LLEN-001`, kein Träger.

## Offene DoD-Punkte (nicht bewertet)

Nicht bewertet sind: ADR-0114 `Accepted` samt Index (Architect, nach diesem Bericht; das
`Schärft:`-Feld nennt `SPEC-112` bereits), der Review-Report (R1 und R2 liegen vor), die
Closure-Notiz mit Lerneintrag, das Beobachtungs-Register (dazu R1 F-5 und F-6), die Ausgänge der
Risiken aus §6 und die drei Paarungen. Alles Planner. Für Risiko 1 belegt
`TestLinelengthOffenerZaunNimmtDenRestAus` mit M22 das Verhalten beim offenen Zaun. Für Risiko 2
belegen die beiden Kopplungs-Tests, dass die Helfer für das Modul nicht geändert wurden. Für Risiko
3 gilt V-1.

**Gesamt:** Die DoD-Punkte 1–4 sind bestätigt. Jedes der 17 Akzeptanzkriterien hat einen Test, und
für jedes ist ein Rot ohne die Regel gezeigt; für die Reihenfolge braucht es zwei Mutationen (V-2).
Die Entscheide tragen am Image, und ohne `linelength`-Konfiguration ist der Lauf byte-identisch zu
`v0.91.0`. Eine DoD-Verletzung habe ich nicht gefunden.
