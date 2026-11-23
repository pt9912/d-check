# Verifikation slice-289 — Anforderung, Spezifikation und ADR für das Modul `linelength`

**Rolle:** Verifier (Modul 8/11) · **Gegenstand:** Plan und DoD aus
`docs/plan/planning/in-progress/slice-289-linelength-anforderung-und-adr.md`, geprüft auf dem
Stand `01849cfa` · **Bezug:**
[`DC-FA-LLEN-001`](../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in),
[`DC-FA-LLEN-001.a`](../../spec/spezifikation.md#dc-fa-llen-001a--zeilenlänge-in-markdown-dateien-linelength),
[ADR-0114](../plan/adr/0114-linelength-modul.md) · Datum 2026-11-23.
Architect-Verdikt und Reviews r1/r2 sind gelesen und nicht neu bewertet. Alle Belege unten habe ich
selbst gefahren; übernommen ist keiner.

## DoD-Punkte

| DoD | Urteil | Beleg |
|---|---|---|
| 1 — Anforderung im Lastenheft | bestätigt | `spec/lastenheft.md`: Kopf `**Version:** 0.112.4`. Das Kürzel `LLEN` steht in der Liste von §3 (Z. 82–83). Die Anforderung beginnt in Z. 4423, die fünf Historie-Zeilen 0.112.0–0.112.4 stehen in Z. 4636–4640, die `Verweis`-Spalte ist durchgehend `—` ([`MR-032`](../../harness/conventions.md#mr-032)). Akzeptanzkriterien gibt es zu jedem Entscheid. Zähleinheit: Umlaut/Emoji, `\r\n`, `\r` ohne `\n`, Tab, ungültiges UTF-8. Jede Ausnahme an und aus, samt ihren Rändern. Grenzwert 80/81, leere Menge, leere Datei. Marker-Zeile bei mindestens 14 ⇒ Befund, bei `max-chars: 10` ⇒ kein Befund. Dazu die Ränder mit Exit 2 (sieben Ursachen), die vier Grenzen der geteilten Lexik als Liste und das Out-of-Scope (1)–(8). |
| 2 — Spezifikation | bestätigt | `DC-FA-LLEN-001.a` (`spec/spezifikation.md` Z. 3575 ff.) führt die Schritte 1–9: Config, Eingabemenge, Zeilen, (a)/(b), Länge, (c), Befund mit `line-chars-exceeded`, Bezugsmenge, Eingabefehler. Dazu kommen *Nicht Gegenstand* und *Grenzen, am Gegenstand gelesen*. Das Schema `linelength.*` (fünf Zeilen) steht in `SPEC-005` (Z. 3983–3987); der Schlüssel `scope` läuft über die generische `<modul>.scope`-Zeile. In der Tabelle von §4 steht der Grund-Code nicht, wie §1 des Plans es vorsieht. Jede Grenz-Aussage habe ich am Code nachgelesen, siehe §Grenz-Aussagen am Gegenstand. |
| 3 — ADR-0114 `Proposed` | bestätigt | `**Status:** Proposed`. Die Index-Zeile in `docs/plan/adr/README.md` Z. 124 führt `Proposed` und nennt `DC-FA-LLEN-001` an erster Stelle. Das `Bezug:`-Feld verlinkt `DC-FA-LLEN-001`, das `Schärft:`-Feld `DC-FA-LLEN-001.a` und `SPEC-005`. `## Re-Evaluierungs-Trigger` ist vorhanden, mit (a)–(e). Dass die Links auflösen, belegt `make doc-check` (grün, siehe DoD 4). |
| 4 — `make gates` grün | bestätigt | Selbst gefahren auf `01849cfa`, sauberer Arbeitsbaum, Exit 0. Die Ausgabe: `fetch-baseline-cache: verify ok (54 Dateien, vollständig)`; `d-check: 1245 Datei(en) geprüft, 0 Befund(e)` (doc-check); `0 issues.` (lint); alle Test-Pakete `ok`; `coverage-gate: OK — Coverage 94.90% erfüllt Schwelle 93%`; Semgrep `0 findings`; gate-consistency/planning-check jeweils `1239 Datei(en) geprüft, 0 Befund(e)`. Die Schlusszeile lautet `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`. |

Ein **Bewusstes Brechen** war nicht nötig: Kein DoD-Punkt dieses Slice beruft sich auf einen Test,
denn der Slice liefert keinen Code. Die Testbehauptungen liegen in `slice-290`, und dessen DoD
verlangt das gezeigte Rot schon.

## Plan gegen Artefakt

**§3-Tabelle.** Laut `git log --grep=slice-289 --name-only` haben 29 Commits nur diese Dateien
berührt: Lastenheft, Spezifikation, ADR-0114, ADR-Index, die Pläne von slice-289 und slice-290,
den Welle-Plan, die Roadmap und drei Review-Dateien. Code, CHANGELOG, README und Handbuch sind
unberührt, die Abgrenzung aus §1 ist damit eingehalten. Die Modul-Aufzählung in `DC-FA-CLI-002`,
das Glossar und die Zeile in §4 nennen `linelength` nicht, wie §1 es vorsieht.

**Entscheide 1–9**, je gegen Lastenheft (LH), Spezifikation (SP) und ADR-0114:

| # | LH | SP | ADR | Bemerkung |
|---|---|---|---|---|
| 1 Eigenes Modul, `LLEN` | ✓ | — | ✓ E1, Alt. B/C | Die Abgrenzung der Längen-Achsen steht in LH und ADR |
| 2 Eingabemenge, §3.8 „keine“ | ✓ | ✓ Schritt 2 + letzter Satz der Grenzen | ✓ E2 | |
| 3 Codepunkte, `\r`, Tab, UTF-8 | ✓ | ✓ Schritt 3/5 | ✓ E3 | Den Verweis auf `DC-FA-PSYNC-001.a` Schritt 4 und ADR-0113 E2 habe ich nachgelesen, er trägt |
| 4 `max-chars` ≥ 1, `>` | ✓ | ✓ Schritt 1, Schema | ✓ E4 | |
| 5 (a)/(b)/(c) Variante H, Überschriften zählen | ✓ | ✓ Schritt 4/6, Grenzen | ✓ E5, Alt. M/H | Die Zahlen sind bestätigt, siehe §Gegenprobe |
| 6 Befund, Reihenfolge | ✓ | ✓ Schritt 7 | ✓ E6 | SP legt zusätzlich `target = max-chars=<N>` fest; das ist eine Präzisierung, keine Abweichung |
| 7 Bezugsmenge | ✓ | ✓ Schritt 8 | ✓ E7 | |
| 8 Leere Menge kein Fehler | ✓ | ✓ Schritt 2 | ✓ E8, Alt. J/K | |
| 9 Kein Marker, 14er-Rand | ✓ | ✓ Schritt 6/7 | ✓ E9, Konsequenzen | |

Keiner der neun Entscheide weicht ab. Die Spec-Straten enthalten keine Verweise abwärts: In
`DC-FA-LLEN-001`, `DC-FA-LLEN-001.a` und den fünf Historie-Zeilen steht kein `slice-`, keine
`ADR-<NNNN>` und kein `Review`.

## Grenz-Aussagen am Gegenstand

Die Greps aus dem ersten Plan-Nachtrag habe ich neu gefahren. Alle Zeilennummern stimmen:
`markdown.go` Z. 91 (`strings.Split(…, "\n")`), Z. 53 (`TrimLeft(raw, " \t")`), Z. 193
(`HasPrefix(lines[i], "|")`); `scan.go` Z. 94, 112, 119, 120; `run.go` Z. 232
(`nicht lesbar: `); `applyScopes` ⇒ 16; `codepaths.go` Z. 15 (die vier Marker-Module).

Darüber hinaus habe ich jede Aussage der Anforderung zu den Grenzen am Code gelesen:

- **Jede Zaun-Zeile schaltet um, ohne Abgleich:** `proseLines` ruft nur `FenceToggle`, nicht
  `FenceCloses`. Eine Backtick-Zeile mit Backtick im Rest schaltet nicht um (`FenceToggle`, Z. 33–46).
  ✓
- **Zaun hinter `>` schaltet nicht um, eingerückter Code-Block ist kein Fenced-Code:**
  `TrimFenceIndent` entfernt nur Space und Tab, eine Erkennung eingerückter Blöcke gibt es nicht. ✓
- **Offener Zaun, gemeldet von `spans`:** `checkUnclosedFence` (`spans.go` Z. 29–50) führt die
  naive Lesart (`stepNaive`, dieselbe wie `proseLines`) neben der strengen und meldet, sobald eine
  der beiden am Dateiende offen steht. Damit meldet `spans` auch den Fall, den nur der naive
  Automat offen sieht, etwa einen Block mit Vier-Backtick-Zaun, den eine innere Drei-Backtick-Zeile früh schließt. ✓
- **Tabellenzeile = in `proseLineSet` und rohe Zeile beginnt mit `|`:** `tableRowLine`
  (Z. 192–194) trimmt nicht und verlangt keine Trennzeile. Mit (a) aus liegt eine `|`-Zeile im Zaun
  nicht in `prose`. ✓
- **Eingabemenge:** `.md` mit Groß-/Kleinschreibung, kein Abstieg in Symlinks (`fs.go`: `Lstat`,
  Symlink ⇒ `KindSymlink`), Skip-Liste `isSkipDir`. ✓
- **Leere Scan-Wurzel ⇒ Exit 2:** `cli.go` Z. 415–418 (`Scan-Wurzel ist leer`). ✓
- **Codepunkte wie `cell-max-chars`:** `structure_tablecell.go` Z. 127 `utf8.RuneCountInString`.
  Ein ungültiges Byte zählt als eine Rune. ✓
- **Zeilenzahl wie bei `file`:** `DC-FA-FILE-001.a` Schritt 5 (Zeilenumbrüche plus unvollständige
  Schlusszeile, leere Datei 0). Das deckt sich mit LLEN Schritt 3, auch für `…\n\r` und `\r`. ✓

Keine Aussage widerspricht dem Gegenstand.

## Gegenprobe

Das Kommando aus dem Plan stammt aus der Repo-Wurzel eines `git archive f3dc94da`, entpackt im
Scratchpad. Gelaufen ist es mit GNU Awk 5.2.1 unter `LC_ALL=C.UTF-8`, das Programm aus einer Datei
(`-f`) mit identischem Text. Es lief als eine einzige `awk`-Instanz, denn jede Ausgabe ist eine
Zeile.

| Größe | Plan | gemessen |
|---|---|---|
| Dateien · Zeilen | 1237 · 86275 | 1237 · 86275 |
| über 120, mit / ohne Fenced-Code | 6108 / 6082 | 6108 / 6082 |
| Tabellenzeilen | 3130 | 3130 |
| (c) alte Fassung / Befunde | 828 / 2124 | 828 / 2124 |
| (c) Variante H / Befunde | 513 / 2439 | 513 / 2439 |
| Bytes (`LC_ALL=C`), mit Fenced-Code | 6233 | 6233 |

Die Aufschlüsselung in ADR-0114 Alternative M/H habe ich mit einem eigenen `awk`-Lauf über dieselbe
Menge nachgezählt: 151 Zeilen ohne Leerraum, 316 mit einem Wort über 120, 361 ohne ein solches
Wort, darunter 141 mit einem letzten Wort von höchstens 20 Zeichen. 46 Befunde der alten Fassung
fallen unter H weg. Daraus folgen 828 − 361 + 46 = 513 und netto 315; beide stimmen. Der Hinweis
aus Review R1 F-9 ist ebenfalls bestätigt: `mawk` liefert ohne Warnung die Bytes-Zahlen
(`… 6233 6206 3159 880 2167`). Auf `01849cfa` ergibt dieselbe `find`-Menge 1239 Dateien, und die
Gates melden 1239 bzw. 1245. Die Differenz von 6 sind die Fixtures, wie im Plan begründet.

## Akzeptanzkriterien gegen die Testfälle in slice-290 §3

17 Kriterien. Vollständig abgedeckt: Happy Path, Grenzwert, Zähleinheit, (a) an/aus, (a) Grenzen,
(b) Grenzen, (c) an/aus, Kein Zeilen-Marker, Leere Menge, Reihenfolge/Determinismus. Teilweise oder
gar nicht abgedeckt sind:

| Kriterium | in §3 von slice-290 |
|---|---|
| (b) an **und aus** | nur „an“ als Testfall; „aus“ allein über DoD 2 („jede Ausnahme … an und aus“) |
| (c) Ränder: ein Wort ≤ `max-chars`, das **nur mit seiner Einrückung** über die Grenze reicht ⇒ Befund | **fehlt** (die Präfix-Fälle `- `, `> `, `- foo `, `-   ` stehen da) |
| Überschrift: überlange ATX-Überschrift ⇒ Befund | **fehlt** |
| Eingabemenge: Datei unter `exempt-paths` und Datei außerhalb von `scope` ⇒ kein Befund | **fehlt**. Nur die Nicht-Markdown-Datei steht da; DoD 1 nennt `scope`/`exempt-paths` nur als Exit-2-Fall |
| Bezugsmenge: stderr (Default) bzw. `summary.notes` (`--json`/`--yaml`), Wortlaut der Zeile | teilweise: Zeilenzahl und eigene Zahl für (c) stehen da, der Ausgabekanal nicht |
| Negative: `max-chars: "80"`, `exempt-tables: "nein"`, ein leeres Glob, `scope` ohne `roots` | teilweise: DoD 1 deckt sie über „jede Fehlkonfiguration der Anforderung“; die Typfehler stehen nicht in der Liste |
| Default byte-identisch (kein Block, kein `--enable`) | **fehlt** |

## Findings

- **HIGH:** keine.
- **MEDIUM — V-1, Widerspruch im Lastenheft selbst.** §5 *Globale Out-of-Scope-Punkte* (Z. 4615)
  schließt „Rechtschreib-, Stil- oder Markdown-Lint-Prüfungen (dafür existieren etablierte Tools)“
  aus. Eine Grenze der Zeilenlänge in Markdown ist eine klassische Markdown-Lint-Regel. Weder
  `DC-FA-LLEN-001` noch ADR-0114 nennt diese Spannung: Alternative A verwirft den fremden Linter,
  ohne §5 zu erwähnen. Verdikt und Reviews behandeln sie ebenfalls nicht
  (`grep -i 'lint\|§5'` über die drei Reports: kein Treffer zu §5). Kein DoD-Punkt ist verletzt,
  denn DoD 1 verlangt das Out-of-Scope der Anforderung, nicht die Verträglichkeit mit §5. Nach
  `AGENTS.md` §1 ist der Widerspruch aber zu melden und nicht stillschweigend aufzulösen.
  Bestehende Module wie `spans` und `structure` mit `cell-max-chars` berühren dieselbe Linie. Das
  ist Präzedenz, aber keine Auflösung. **Träger:** Architect. Er entscheidet, ob §5 geschärft wird,
  etwa „Stilprüfungen, die keine deklarierte, deterministische Repo-Invariante sind“, oder ob
  `DC-FA-LLEN-001` sich in einem Satz gegen §5 abgrenzt; den Text im Lastenheft zieht der
  Implementer nach, als neue Historie-Zeile. Spätester sinnvoller Punkt: vor `Accepted` von
  ADR-0114 mit der Closure von slice-290.
- **LOW — V-2, Plan-Abweichung im Testfall-Nachzug.** Die Zeile „Plan von slice-290 | update —
  Planner“ in §3 verlangt die Testfälle, die aus den bestätigten Entscheiden folgen. Vier
  Kriterien fehlen dort ganz: Einrückungs-Rand von (c), ATX-Überschrift, Eingabemenge über
  `exempt-paths`/`scope` und Default byte-identisch. Drei weitere sind nur über die DoD von
  slice-290 gedeckt (siehe Tabelle oben). Gemildert ist das dadurch, dass §3 von slice-290 sagt
  „maßgeblich ist die Anforderung“ und Risiko 3 dort eine Zuordnung je Randbedingung im
  Plan-Nachtrag verlangt. Gerade dieses Risiko
  ([`BEO-ALL/spec-randbedingung-ohne-test`](../plan/planning/observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md))
  spricht aber dafür, die Liste jetzt zu schließen. **Träger:** Planner (Plan von slice-290).
- **INFO — V-3.** Die Aussage „eine Datei der Eingabemenge, die nicht gelesen werden kann, beendet
  den Lauf mit Exit 2“ stimmt am Code (`run.go` Z. 232). Eine FIFO mit der Endung `.md` meldet
  `List` aber als `KindFile`, und `ReadFile` blockiert. Diese Grenze benennt die Spezifikation nur
  für `pinsync`. Sie gilt für jedes Modul, das die Scan-Menge liest, und ist kein Gegenstand dieses
  Slice; ein Träger ist nicht nötig.

## Offene DoD-Punkte (Planner, nicht bewertet)

Nicht bewertet sind: Review durchgeführt und Report, Closure-Notiz mit Lerneintrag,
Beobachtungs-Register, Ausgänge der Risiken aus §6 und die drei Paarungen.

**Gesamt:** Alle vier geprüften DoD-Punkte sind bestätigt, die Entscheide 1–9 stimmen in allen drei
Artefakten überein, die Grenz-Aussagen stimmen am Code und die Zahlen beider Plan-Nachträge sind
reproduziert. Offen bleiben ein Widerspruch im Lastenheft (V-1, Architect) und eine Lücke im
Testfall-Nachzug von slice-290 (V-2, Planner).
