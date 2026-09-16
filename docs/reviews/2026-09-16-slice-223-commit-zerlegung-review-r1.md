# Review-Report: slice-223 — 2026-09-16 (R1)

**Review-Art:** Code-Review — geprüft wird der Diff gegen Slice-Plan,
Hard Rules und die zitierten `MR`-Einträge.

**Gegenstand:** Commit-Kette `4c29298c^..815bea69` (`4c29298c`, `b2aef315`,
`a148466d`, `67082e6f`, `815bea69`). Der Slice liegt zum Prüfzeitpunkt noch
in `in-progress/`; `HEAD` des Repos ist inzwischen auf `99f1ca88`
weitergewandert (ein unabhängiger `slice-221`-Commit), der nicht Gegenstand
dieses Reports ist.

**Skill:** `.harness/skills/reviewer.md` @ v1.16.0 / 2026-09-07 ·
**Modell:** `claude-sonnet-5` · **Datum:** 2026-09-16

> **Zitier-Form.** Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb Kennung statt Adresse: `slice-NNN` statt Lifecycle-Pfad,
> `make <target>` als Token statt Link auf die Sensor-Datei, eine
> Baseline-Stelle als Tag plus Pfad in Inline-Code
> (`v6.9.0` · `regelwerk/<datei>.md` §<Abschnitt>).

**Eingangs-Kontext:**

- Slice-Plan `slice-223` (Fassung `in-progress/`, Stand des Prüfzeitpunkts,
  inkl. §7 Closure-Notiz und der Plan-Änderungs-Passage in §1)
- `MR-013` (getrimmte Fassung) · aufgelöst: `MR-059`/`MR-061`/`MR-062`/
  `MR-063`/`MR-064` · Präzedenz-Zitate `MR-014`/`MR-027`/`MR-038`
- `AGENTS.md` §3.3 (neue, gekürzte Form) · §3.7 · §5 · §6
- `harness/conventions.md` §Aktive/§Aufgelöste Adaptionen
- Beobachtungs-Register: `BEO-ALL/citation-stretched-beyond-scope` (19×,
  neuer Beleg `evidence/slice-223.md`)
- Vorgänger-Kontext: `docs/plan/planning/done/slice-222-baseline-v660-bump.md`,
  `docs/plan/planning/done/slice-224-baseline-v690-bump.md` (`ignore-refs`
  auf `slice-223` in `open/`)
- `.harness/baseline/v6.9.0/templates/AGENTS.template.md` §3.3 (Kanon-Vergleich)
- `.harness/baseline/v6.9.0/regelwerk/grundlagen-traceability.md`
  §Herkunfts-Anker (Ruheort-Regel-Zitat)
- Baseline `v6.9.0` · `.harness/baseline/v6.9.0/`

---

## Findings

### F-1 — Die zentrale Prämisse für den Fortbestand von MR-013 wird durch die eigene Commit-Historie dieses Slice widerlegt

- `kategorie`: HIGH
- `quelle`: `MR-013` §Adaption/§Begründung · `AGENTS.md` §3.3 (Ausnahme
  MR-/Wellen-Lifecycle-Move) · `AGENTS.md` §5 (*„Eine Commit-Botschaft oder
  Closure-Notiz behauptet nicht mehr, als die Arbeit trägt"*)
- `pfad`: Commit `a148466d` (gesamt) · `harness/conventions.md:136-140` ·
  `harness/sensors/archive-wave.md:61-64`
- `befund`: Der gesamte verbleibende Geltungsbereich von `MR-013` (und damit
  ein Kernstück der Closure-Notiz und des Plans) stützt sich auf die
  Behauptung, der kanonische Zwei-Commit-Weg (reiner `git mv`, dann
  Korrektur-Commit, beide im selben Push) sei *„zwischen den beiden Commits
  lokal gar nicht committierbar"*, weil der lokale `pre-commit`-Hook
  `doc-check` auf **jedem** Commit fährt und einen Zwischenstand mit
  gebrochenen Links ablehnen würde. Diese Behauptung wird durch die eigene
  Commit-Sequenz des Slice widerlegt: Commit `a148466d` verschiebt fünf
  `MR`-Dateien nach `harness/conventions/done/` und korrigiert dabei nur die
  **internen** Querverweise der bewegten Dateien — die **externen**
  Rückverweise in `harness/conventions.md` (fünf Links, Zeilen 136-140) und
  `harness/sensors/archive-wave.md` (vier Links, Zeilen 61-64) bleiben auf
  den alten, jetzt nicht mehr existierenden Pfaden stehen. Nachgebaut mit
  `git worktree add … a148466d && make doc-check`: **9 Befunde,
  `target-missing`**, Exit 1 — der Commit ist für sich genommen doc-check-rot.
  Erst der nächste Commit `67082e6f` (eine Minute später) räumt das auf
  (verifiziert: 0 Befunde). `core.hooksPath` steht in diesem Klon auf
  `.githooks`, und `.githooks/pre-commit` ruft `make doc-check`
  bedingungslos vor jedem Commit auf — exakt der Mechanismus, den `MR-013`
  als Grund nennt, warum der kanonische Weg hier *nicht* gangbar sei. Dass
  `a148466d` trotzdem existiert, zeigt entweder, dass der Hook für diesen
  Commit nicht gegriffen hat (z. B. `--no-verify`, was das Git-Safety-Protokoll
  ohne ausdrücklichen Nutzer-Auftrag verbietet), oder dass die Prämisse
  schlicht falsch ist: ein roter Zwischenstand **war** lokal committierbar.
  In beiden Fällen trägt die Begründung, mit der `MR-013` statt vollständiger
  Auflösung nur getrimmt wurde, nicht das Gewicht, das ihr Plan und
  Closure-Notiz geben — und Risiko-Ausgang 3 in §6 des Plans (*„eine lokale
  Werkzeug-Grenze"*) verallgemeinert genau diese unbelegte Prämisse als
  gesicherten Fakt.
- `verifizierbar`: ja — `git worktree add /tmp/wt a148466d && (cd /tmp/wt &&
  make doc-check)` liefert `9 Befund(e)`, `target-missing`, gefolgt von
  `git worktree add /tmp/wt2 67082e6f && (cd /tmp/wt2 && make doc-check)` mit
  `0 Befund(e)`.
- `klasse`: `praemisse-durch-eigene-commit-historie-falsifiziert`

### F-2 — Die Scope-Reduktion (sechs → fünf Einträge) ist im Code vollzogen, bevor sie im Plan dokumentiert ist

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §6 (*„Wer im Lauf etwas mitnimmt, das der Plan
  ausschließt, hat den Plan geändert — und das gehört vor den Code, nicht in
  den Bericht danach"*)
- `pfad`: Commit `a148466d` (20:44:38) vs. Commit `67082e6f` (20:45:29),
  `docs/plan/planning/in-progress/slice-223-commit-zerlegung-ausnahmen-aufloesen.md`
- `befund`: Der Plan sah vor, **sechs** `MR`-Einträge aufzulösen (§1, §2
  DoD (2) vor der Änderung). Commit `a148466d` vollzieht bereits die
  reduzierte Menge — es verschiebt nur die **fünf** Einträge 059/061/062/063/064
  und lässt `MR-013` unangetastet in der Tabelle *Aktive Adaptionen* stehen.
  Die Plan-Datei selbst wird in `a148466d` **nicht** angefasst (`git show
  --stat a148466d` zeigt keinen Treffer für den Slice-Plan). Der Absatz
  *„Plan-Änderung bei der Implementierung"*, der diese Reduktion begründet
  und DoD (2) entsprechend umformuliert, wird erst im **nächsten** Commit
  `67082e6f` (eine Minute später) hinzugefügt — im selben Commit, der auch
  `harness/conventions.md`s Index-Zeile für `MR-013` aktualisiert. Für das
  Zeitfenster zwischen beiden Commits behauptete die Plan-Datei also noch
  „sechs", während der Code bereits „fünf" umgesetzt hatte, ohne dass die
  Abweichung dort vermerkt war. Die Regel, wonach eine Plan-Änderung *vor*
  dem sie umsetzenden Code steht, ist damit nur für die Hälfte der Änderung
  (den `AGENTS.md`-Umbau in `67082e6f`) eingehalten, nicht für den ersten
  Schritt der Umsetzung (`a148466d`).
- `verifizierbar`: ja — `git show --stat a148466d` (kein Treffer für den
  Slice-Plan) gegen `git show 67082e6f -- '*slice-223*'` (fügt den
  Plan-Änderungs-Absatz hinzu).
- `klasse`: `plan-aenderung-nach-dem-ersten-umsetzungs-commit`

### F-3 — Commit-Botschaft von `67082e6f` schreibt eine Fremdursache für den `archive-wave.md`-Link-Fix

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §5 (Commit-Botschaft behauptet nicht mehr, als die
  Arbeit trägt — hier: nicht die falsche Ursache)
- `pfad`: Commit `67082e6f` Botschaft, letzter Absatz; `harness/sensors/archive-wave.md:61-64`
- `befund`: Die Botschaft von `67082e6f` erklärt die Änderungen an
  `reviewer.md` und `harness/sensors/archive-wave.md` gemeinsam als
  *„Zitat-/Link-Nachzüge, die der Paragraph-3.3-Umbau ausgelöst hat"*. Für
  `reviewer.md` stimmt das (die Zeilennummern-Verschiebung im `d-check:cite`
  folgt direkt aus der Kürzung von `AGENTS.md` §3.3). Für
  `harness/sensors/archive-wave.md` stimmt es nicht: Diese vier Links waren
  bereits durch den MR-Datei-Move in `a148466d` gebrochen (siehe F-1) — sie
  haben nichts mit dem `AGENTS.md`-§3.3-Umbau zu tun, der erst im selben
  Commit `67082e6f` stattfindet. Die Botschaft bündelt zwei ursächlich
  getrennte Reparaturen unter einer gemeinsamen, für eine von beiden
  unzutreffenden Erklärung.
- `verifizierbar`: ja — `git show a148466d -- harness/conventions/` zeigt die
  Dateien bereits verschoben, `git show 67082e6f -- harness/sensors/archive-wave.md`
  zeigt den Link-Fix erst hier.
- `klasse`: `commit-botschaft-falsche-kausale-zuordnung`

---

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Prüffrage 1 (Gate wird still grün) | nicht anwendbar — kein Gate-Skript/Workflow in diesem Diff berührt |
| Prüffrage 2 (Kern-Modul meldet falsch) | nicht anwendbar — kein Go-Code im Diff |
| Prüffrage 3 (Hexagon-Import-Verstoß) | nicht anwendbar — kein Code geändert |
| Prüffrage 4 (Suppression ohne ADR) | nicht anwendbar — keine `//nolint`/Schwellen im Diff |
| Prüffrage 5 (Netzzugriff außerhalb `external`) | nicht anwendbar — keine Netz-Skripte berührt |
| Prüffrage 6 (Kommentar ohne der fünf Klassen) | geprüft, ohne Befund — die einzige neue Konfigurations-Kommentarzeile (`.d-check.yml`, neue `ignore-refs`-Begründung für `slice-222`/`slice-224`→`slice-223`-Zitate) trägt eine **Grenze**-Aussage (warum die Folge-Slice-Paarung hier toleriert wird), keine Chronik/Deliberation |
| Prüffrage 7 (Zustandsfeld erzählt Chronik) | geprüft, ohne Befund — `harness/conventions.md`s „Aufgelöst durch"-Zellen für `MR-059/061/062/063/064` nennen Zustand (`Baseline-Konformität`/`erschöpft`) und Beleg (`seit slice-223`), keine Chronik; Stil deckt sich mit dem Präzedenzformat von `MR-014`/`MR-027` |
| Prüffrage 8 (Botschaft verallgemeinert über Messung) | siehe F-1 (Risiko-Ausgang 3 überträgt die unbelegte Prämisse als Fakt) und F-3 (falsche Kausalzuordnung) |
| Prüffrage 9 (Quelle über Geltungsbereich hinaus zitiert) | geprüft, ohne Befund im gelieferten Stand: `MR-013`s neues `Ersetzt-Baseline-Regel`-Feld zitiert `grundlagen-traceability.md#herkunfts-anker` und dort wörtlich den Satz zur `MR-<NNN>`-Ruheort-Ausnahme (*„Der `git mv` zieht die Pfad-Berichtigung nach sich, als eigener Commit nach dem Umzug"*) — dieser Satz steht tatsächlich im Absatz über Adaptions-Einträge, die dauerhaft in `harness/conventions/` leben, exakt der hier vorliegende Fall. Die ursprünglich falsch zitierte Stelle (`modul-05`) ist korrigiert; die im Plan selbst schon vorab gemeldete Fehlzitierung der „Beide Commits gehören in denselben Push"-Stelle (Anker-Paarung, nicht Slice-Lifecycle) ist im gelieferten Stand nicht mehr als tragendes Argument für die zwei aufgelösten Fälle in Verwendung |
| Prüffrage 10 (Messmethode klafft gegen Spec) | nicht anwendbar — Slice-Kopf „Berührte Spec-Stellen: —" |
| Prüffrage 11 (zwei Module, dieselbe Eingabeklasse) | nicht anwendbar — kein Modul-Code geändert |
| Prüffrage 12 (Erkennung < Alt-Tool-Familie) | nicht anwendbar |
| Prüffrage 13 (fehlender Negativtest) | nicht anwendbar — kein neuer öffentlicher Vertrag |
| Prüffrage 14 (Provenance-Marker begründet statt zeigt) | geprüft, ohne Befund — keine `d-check:status-provenance`-Marker im Diff |
| Prüffrage 15 (Modul liest ungescannte Eingabe) | nicht anwendbar — kein Modul-Code geändert |
| Prüffrage 16 (neues `Schärft:`/`Bezug:` nur „§N") | geprüft, ohne Befund — `MR-013`s `Ersetzt-Baseline-Regel`-Feld trägt einen vollen Anker (`#herkunfts-anker`) plus Zitat, kein bloßes „§N"; keine neuen `SPEC-*`/`ARC-*`-Zielelemente im Diff |
| Prüffrage 17 (Messung zählt Proxy) | geprüft, ohne Befund — `citation-stretched-beyond-scope` „19×" gegen `ls .../evidence/ \| wc -l` nachgezählt: exakt 19 Dateien |
| Prüffrage 18 (Grenzen-Liste ohne größte Lücke) | Befund — weder Plan §1 Abgrenzung (4 Punkte) noch §6 Risiken benennen, dass die zentrale Prämisse für den `MR-013`-Fortbestand ungetestet/unbelegt blieb, obwohl die eigene Commit-Sequenz sie widerlegt (siehe F-1); dieselbe Lücke, nicht als separates Finding gezählt |
| Link-Tiefen in allen fünf verschobenen `MR-0XX`-Dateien | geprüft — jede Referenz auf `MR-013` trägt korrekt `../`, jede Referenz innerhalb des Batches (059↔061↔062↔063↔064) bleibt korrekt bar (gleiches Verzeichnis), alle `../../../`-Pfade auf `docs/plan/…`/`.harness/baseline/…` sind um eine Ebene erhöht — bestätigt durch `make doc-check` (0 Befunde am gelieferten Stand) |
| `harness/conventions.md`-Tabellenkonsistenz | geprüft — alle fünf Anchor-Slugs stimmen mit den jeweiligen Datei-Überschriften überein, keine Kennung kommt doppelt oder fehlend vor, `MR-013`-Zeile korrekt aktualisiert (Geltungsbereich, Ersetzt-Baseline-Regel: „keine — …") |
| `ignore-refs`-Einträge für `slice-223` in `open/` | geprüft — beide Einträge sind datei-skopiert (`in:` exakt eine Datei) mit exakter Ziel-Pfad-Zeichenkette, kein Glob; Begründung (Lauf-Beleg eingefrorener `done/`-Slices, Modul 6 „bis zur Prüfung kann er weitergewandert sein") ist sachlich korrekt und deckt sich mit der bestehenden Praxis für andere Tombstone-Einträge in derselben Datei |
| `AGENTS.md` §3.3 gegen `AGENTS.template.md` (`v6.9.0`) | geprüft — die zwei Kanon-Fälle sind byte-identisch übernommen (`diff` bestätigt Identität zwischen `v6.6.0`- und `v6.9.0`-Vorlage für §3.3); die Zusatz-Klarstellung „Andere Dateien … dürfen im selben Move-Commit mitreisen" widerspricht dem Kanon nicht und ist durch `MR-013`s Ersetzt-Baseline-Regel-Feld gedeckt |
| `AGENTS.md` §3.3 in sich schlüssig, kein Widerspruch zu `MR-013`-Datei | geprüft, ohne Befund — Wortlaut der „Ausnahme MR-/Wellen-Lifecycle-Move" in `AGENTS.md` und `MR-013`s Adaptions-Feld stimmen sachlich überein (Grund, Mechanismus, Geltungsbereich) |
| `reviewer.md`-Zitat-Zeilennummer nach §3.3-Kürzung | geprüft — `AGENTS.md:366-367` trägt exakt den zitierten Satz; kein weiterer `d-check:cite AGENTS.md`-Verweis im Repo, der nach der Kürzung hätte verschieben müssen |
| `make gates` am gelieferten Stand (`815bea69`) und am aktuellen `HEAD` (`99f1ca88`) | grün — `git worktree add … 815bea69 && make doc-check` → `785 Datei(en) geprüft, 0 Befund(e)`; `make gates` am aktuellen Arbeitsbaum → `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` |
| `make doc-check` an den Zwischen-Commits | `4c29298c`: nicht separat geprüft (Botschaft behauptet „zehn Gates, 784 Dateien, 0 Befunde", plausibel, keine Datei-Änderung mit Bruchpotential) · `b2aef315`: verifiziert grün (0 Befunde) · `a148466d`: **verifiziert rot**, 9 `target-missing` (siehe F-1) · `67082e6f`: verifiziert grün (0 Befunde) · `815bea69`: verifiziert grün (0 Befunde) |
| Kommentar-Fünf-Klassen-Regel (`AGENTS.md` §3.7) an den neu geschriebenen Stellen | geprüft — außer der einen `.d-check.yml`-Kommentarzeile (s. Prüffrage 6) sind alle Änderungen Markdown-Prosa (`AGENTS.md`, `MR`-Dateien, Slice-Plan) und fallen nicht unter §3.7 (gilt „für Code, Konfiguration und Skripte"); kein Go-Code, kein Makefile, keine Shell-Skripte im Diff |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:**
`praemisse-durch-eigene-commit-historie-falsifiziert` ·
`plan-aenderung-nach-dem-ersten-umsetzungs-commit` ·
`commit-botschaft-falsche-kausale-zuordnung`

## Verdikt

**Merge-blockierend:** F-1 (HIGH) blockiert die Closure in der jetzigen
Form — nicht weil der gelieferte Endstand fehlerhaft wäre (`make gates` ist
am `HEAD` von `815bea69` und am aktuellen `HEAD` grün, alle Link-Tiefen und
Tabellen sind korrekt), sondern weil die **Begründung**, mit der `MR-013`
statt vollständiger Auflösung nur getrimmt und im Repo belassen wird, durch
die eigene Commit-Sequenz des Slice widerlegt ist. Der Slice trifft damit
eine Architektur-Entscheidung (Abweichung bleibt bestehen) auf Basis einer
Prämisse, die sein eigener Beleg (`a148466d`) nicht trägt. Vor der Closure
zu klären: entweder wird belegt, dass der lokale Hook zum fraglichen
Zeitpunkt tatsächlich aktiv war und der Zwischenstand aus einem anderen
Grund zustande kam (z. B. dokumentierter `--no-verify`-Einsatz mit
Begründung), oder die Schlussfolgerung des Slices ist zu revidieren: Wenn
ein roter Zwischenstand lokal tatsächlich committierbar ist (wie `a148466d`
zeigt), fehlt der tragende Grund, `MR-013` nicht ebenfalls vollständig
aufzulösen.

F-2 (MEDIUM) und F-3 (LOW) sind eigenständig behebbar (Plan-Text ergänzen
bzw. Commit-Botschaft ist bereits historisch — hier nur zur Kenntnis für die
Closure-Notiz) und blockieren für sich genommen nicht.

**Übergabe:** Alle drei Findings gehen an den Implementer. F-1 gehört vor
den nächsten Schritt (Closure-Move) geklärt, da sonst eine unbelegte
Architektur-Prämisse in `done/` einfriert. F-3 sollte in der Closure-Notiz
§7 als Präzisierung der „Was Friktion war"-Passage nachgetragen werden.
Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation.
