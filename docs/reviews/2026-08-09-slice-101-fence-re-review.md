# Review-Report: slice-101 (Re-Review) — 2026-08-10

**Review-Art:** Code — **Re-Review** (bestätigend). Der Erst-Review war
merge-blockierend; geprüft wird, ob die Heilung **am Lauf** hält, nicht ob sie
plausibel klingt (Modul 10 §Drei Review-Arten). Der Dateiname trägt den
Datums-Stamm des Erst-Reports, der Lauf fand am 2026-08-10 statt.

**Gegenstand:** [slice-101](../plan/planning/in-progress/slice-101-fence-unbalanciert.md),
Heilungs-Commit `8f16f2b` (Diff-Range `aba7792..HEAD`); Gesamt-Slice
`38c36b2..HEAD`.

**Skill:** `.harness/skills/reviewer.md` @ 1.3.0 ·
**Modell:** claude-opus-5[1m] · **Datum:** 2026-08-10

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [Erst-Report vom 2026-08-09](2026-08-09-slice-101-fence-review.md) (F-1 bis F-7
  als Checkliste)
- [slice-101](../plan/planning/in-progress/slice-101-fence-unbalanciert.md) §5 DoD
  und §6 Risiken (die dort behaupteten Ausgänge)
- [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md) (Entscheidungen 1–5,
  Fitness Function, Konsequenzen) und
  [ADR-0042](../plan/adr/0042-markdown-lexik-folgt-commonmark.md)
- [`DC-FA-SPAN-001`](../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
  (Klasse 3, Lastenheft-Version 0.52.1) und
  [§`DC-FA-SPAN-001.a`](../../spec/spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung)
  Schritt 3
- [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  (der Konsument, dessen Silent-Grün den Anlass gab)
- [`AGENTS.md`](../../AGENTS.md) (Hard Rules)

**Läufe dieses Reviews.** Alle Fixtures in einem Temp-Verzeichnis außerhalb des
Repos, alle Läufe netzlos und read-only. Gebaut wurden **zwei** Vergleichs-Images
aus dem Repo selbst: der Stand **vor** der Heilung (`aba7792`, über einen
git-Worktree im Temp-Verzeichnis) und HEAD; dazu zwei bewusst mutierte Images für
gezielte Sonden. Gefahren: `make build`, `make test`, `make lint`,
`make doc-check`, `make gates`, `make completeness-check`,
`make verify-closure-notes`, `make bench` sowie rund 60 Fixture-Läufe gegen die
Images. `make gates` ist grün (349 Dateien, 0 Befunde, alle acht Gates),
`make completeness-check` liefert 48 Anforderungen und 0 Waisen,
`make verify-closure-notes` 320 Dateien und 0 Befunde, `make bench` 651 ms Median
(Kriterium 5000 ms). Der Arbeitsbaum ist am Ende unverändert.

---

## Erst-Befunde: Status

| Erst-Befund | Kategorie | Status | Beleg (Lauf) |
|---|---|---|---|
| F-1 — Wächter wertet einen anderen Fence-Automaten aus als der Tabellen-Leser | HIGH | **geheilt** | Original-Szenario nachgebaut (Anforderungs-Tabelle, dahinter ein Backtick-Öffner mit Infozeile, den eine Tilden-Zeile „schließt", dahinter eine zweite Tabelle): Vollständigkeits-Lauf weiterhin 2 Anforderungen, 0 Waisen, Exit 0 — der `spans`-Lauf derselben Datei meldet jetzt `fence-unclosed` an der Öffnungszeile, Exit 1 (vorher 0 Befunde, Exit 0). Tabellen-Leser **unverändert**: 17 Fence-Randfälle gegen Alt- und Neu-Image, 17 von 17 byte-gleiche RTM-Ergebnisse |
| F-2 — Wächter und Vorverarbeitung trimmen verschieden | HIGH | **halb geheilt** | Original-Fixture (U+00A0 vor einer Fence-Zeile, dahinter ein echter, nie geschlossener Fence): Alt-Image 0 Befunde/Exit 0, Neu-Image meldet den echten Öffner. Aber die Trimmung ist nur mit `proseLines` und dem Tabellen-Leser angeglichen — das Modul `planning` trimmt weiter unicode-weit, und dort ist derselbe Silent-Grün-Pfad reproduzierbar (siehe N-1) |
| F-3 — Grenze deckt nur die Quell-Achse | MEDIUM | **als Grenze dokumentiert, Aufzählung unvollständig** | Original-Fixture reproduziert (Ziel-Datei außerhalb der Scan-Wurzeln, Status-Zeile hinter offenem Fence): 0 Befunde, Exit 0; Gegenprobe ohne Fence liefert `matrix-inactive`. Die vier genannten Module stimmen (im Code je über die fence-sensitive Vorverarbeitung belegt) — es fehlen `codepaths` und `pins`, und ausgerechnet `pins` ist die **stille** Richtung (siehe N-2) |
| F-4 — Parität statt Paarung | MEDIUM | **teilweise geheilt** | Die legale Verschachtelung ist grün: ein Vier-Zeichen-Block, der einen Drei-Zeichen-Block zeigt, erzeugt keinen Befund mehr (Alt-Image ebenfalls grün, Neu-Image grün, Test vorhanden). Die zweite Hälfte steht unverändert: zwei Fence-Zeilen in einem vier Leerzeichen eingerückten Codeblock verschlucken die Prosa dazwischen — Alt und Neu je 0 Befunde, Gegenprobe ohne die beiden Zeilen meldet `target-missing` (siehe N-6) |
| F-5 — zugesagtes Befund-Ziel ohne Assertion | MEDIUM | **überwiegend geheilt** | Selbst nachvollzogen: sechs Rückbauten einzeln angewendet (über eine Dateikopie, nicht über `git checkout`), `make test` je rot — Kappung 30 auf 300, Rechts-Trimmung entfernt, strenge Lesart entfernt, naive Lesart entfernt, Längen-Abgleich entfernt, Aufruf aus der Modul-Oberfläche entfernt. Ein siebter Rückbau bleibt **grün** (siehe N-4) |
| F-6 — gemeldete Zeile ist nicht die Reparaturstelle | LOW | **dokumentiert, Dokumentation im Regelfall falsch** | Original-Fixture reproduziert: Öffner ohne Schluss in Zeile 3, zwei intakte Blöcke dahinter, Befund `a.md:14` — die Schlusszeile eines intakten Blocks. Anforderung, Spezifikation und ADR nennen die Zeile jetzt **Fundstelle**, binden die Einschränkung aber an „nur die Parität des naiven Toggles kippt". Eine Sonde mit ausschließlich strenger Lesart meldet in derselben Datei ebenfalls Zeile 14 — die Einschränkung gilt auch dort (siehe N-5) |
| F-7 — Sammelsatz nicht mit der dritten Klasse mitgewachsen | LOW | **geheilt** | Der Sammelsatz nennt das Ziel jetzt je Klasse; Lauf gegen eine Datei mit langer Infozeile liefert genau die getrimmte Fence-Zeile auf 30 Runen, ohne Zeichen mitten zu zerschneiden |

### Mutations-Gegenprobe: eigener Nachvollzug

Verfahren: Datei-Kopie im Temp-Verzeichnis sichern, mutieren, `make test`, aus der
Kopie zurückschreiben; jede Ersetzung bricht ab, wenn sie nicht **genau einmal**
greift (der Fehler des ersten Anlaufs). Der Arbeitsbaum war nach jedem Durchlauf
leer (`git status --short`).

| Rückbau | Ergebnis | roter Testfall |
|---|---|---|
| strenge Lesart ersatzlos entfernt | rot | drei Fälle: Tilde schließt keinen Backtick-Block, Backtick keinen Tilde-Block, kürzerer Schluss schließt nicht |
| naive Lesart ersatzlos entfernt | rot | „nur die naive Lesart kippt" |
| Trimmung wieder unicode-weit | rot | „U+00A0 vor der Fence ist kein Toggle" |
| Kappung von 30 auf 300 | rot | „auf 30 Runen gekappt" |
| Längen-Abgleich im Schluss-Prädikat entfernt | rot | „kürzerer Schluss schließt nicht", „legale Verschachtelung bleibt grün" |
| Aufruf aus der Modul-Oberfläche entfernt | rot | Verdrahtungs-Test |
| Rechts-Trimmung des Ziels entfernt | rot | „Trailing-Whitespace weg" |
| **Vorrang der strengen Lesart getauscht** | **grün** | keiner — siehe N-4 |

Die im Slice behaupteten „sechs Rückbauten, alle rot" sind damit bestätigt und um
einen siebten ergänzt; der achte widerlegt die Vollständigkeit der Zusage.

## Findings

### N-1 — Der bewachte Konsument trimmt unicode-weit; der Anlass-Fall dieses Slice ist mit einem geschützten Leerzeichen weiterhin still grün

- `kategorie`: HIGH
- `quelle`: [§`DC-FA-SPAN-001.a`](../../spec/spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung)
  Schritt 3 („mit derselben Trimmung wie die Vorverarbeitung"),
  [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md) §Konsequenzen
  („Jedes Modul profitiert, weil beide Lesarten geteilt sind"),
  Reviewer-Skill §HIGH (Stilles-Grün in einem Gate)
- `pfad`: `internal/hexagon/core/rules/planning.go:168` und
  `internal/hexagon/core/rules/planning.go:195` gegen
  `internal/hexagon/core/rules/spans.go:34`
- `befund`: Die Closure-Fähigkeit des Moduls `planning` führt ihren
  Fence-Automaten über `strings.TrimSpace(raw)` — unicode-weit —, während der
  geheilte Wächter (und `proseLines`) `strings.TrimLeft(raw, " \t")` benutzt. Eine
  mit U+00A0 eingerückte Fence-Zeile ist deshalb **nur** für den Konsumenten ein
  Toggle: sie verschluckt den Rest des Closure-Abschnitts, und der Wächter sieht
  dort keinen Fence. Genau der Fall aus §2 des Slice — eine Floskel hinter einem
  ungeschlossenen Fence — ist damit unverändert reproduzierbar.
- `verifizierbar`: ja — Fixture mit einer Closure-Notiz aus fünf Satzenden, danach
  eine Zeile aus U+00A0 plus drei Backticks, danach die konfigurierte Floskel:
  `--enable planning` meldet **0 Befunde, Exit 0**, und der `spans`-Lauf derselben
  Datei ebenfalls **0 Befunde, Exit 0**. Gegenprobe A (Zeile entfernt):
  `closure-note-boilerplate`, Exit 1. Gegenprobe B (dieselbe Zeile ohne das
  geschützte Leerzeichen): `planning` bleibt still, aber `spans` meldet
  `fence-unclosed` — der Wächter greift also, wo die Trimmungen übereinstimmen,
  und nur dort. Der Effekt trägt bis in `make verify-closure-notes`.
- `klasse`: „Wächter und Bewachtes trimmen verschieden"

### N-2 — Die dokumentierte Ziel-Achsen-Grenze zählt vier Module auf; das Modul mit der stillen Richtung fehlt

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-SPAN-001`](../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
  Absatz „Grenze"; [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md)
  §Konsequenzen;
  [`DC-FA-PIN-001`](../../spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in)
- `pfad`: `internal/hexagon/core/rules/pins.go:114` (`spanHash`-Rückgabe wird
  verworfen) und `internal/hexagon/core/rules/codepaths.go:294`; Vertragstext in
  `spec/lastenheft.md`
- `befund`: Die Grenze nennt als Ziel-Achsen-Leser `matrix`, `anchors`,
  `diagrams` und `versions` — alle vier stimmen. Zwei weitere Module lesen
  Zieldateien über dieselbe fence-sensitive Vorverarbeitung und fehlen:
  `codepaths` (Anker-Menge der Zieldatei) und `pins` (Heading-Section des
  gepinnten Ziels). Bei `pins` ist die Richtung **still**: ein nicht auflösbarer
  Anker liefert kein Finding, also verschwindet die Drift-Prüfung, statt laut zu
  werden — dieselbe Eigenschaft, deretwegen die Grenze überhaupt aufgeschrieben
  wurde.
- `verifizierbar`: ja — Fixture mit `scan.roots` auf ein Unterverzeichnis, darin
  ein Link mit Content-Pin auf eine Ziel-Datei außerhalb der Wurzeln. Ohne Fence
  meldet der Lauf `link-stale` (Drift erkannt); mit einem offenen Fence vor der
  Ziel-Überschrift **0 Befunde, Exit 0** — und `spans` sieht die Ziel-Datei nicht.
  Für `codepaths` dieselbe Bauform: mit Fence `anchor-missing`, ohne Fence
  befundfrei.
- `klasse`: „Modul-Grenze nur auf der Quell-Achse gedacht"

### N-3 — Die neue Trimmung lässt ein Wagenrücklauf-Zeichen ins Befund-Ziel durch und zerstört die Befundzeile auf CRLF-Dokumenten

- `kategorie`: MEDIUM
- `quelle`: [§`DC-FA-SPAN-001.a`](../../spec/spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung)
  Schritt 3 („Das gemeldete Ziel ist diese Zeile, links- und rechtsbündig
  getrimmt")
- `pfad`: `internal/hexagon/core/rules/spans.go:47`
- `befund`: Vor der Heilung entstand das Ziel aus `strings.TrimSpace` und war
  damit auch von einem CR befreit; jetzt trimmen Links- und Rechtsseite nur Space
  und Tab. Auf einer Datei mit CRLF-Zeilenenden trägt das gemeldete Ziel deshalb
  ein rohes CR. In der Zeilen-Ausgabe steht es zwischen zwei Tabulatoren: das
  Terminal setzt den Cursor an den Zeilenanfang zurück, und die gerenderte Zeile
  zeigt am Ende Datei, Zeile und Grund — das Ziel selbst ist überschrieben und
  unsichtbar. In der JSON-Ausgabe erscheint das Steuerzeichen im Feld `target`.
- `verifizierbar`: ja — dieselbe CRLF-Datei gegen beide Images: Alt-Image liefert
  als Ziel die Fence-Zeile, Neu-Image dieselbe Zeile plus CR (in `--doctor` als
  `Stelle:`-Wert und in `--json` im Feld `target` sichtbar). Der Erst-Report hatte
  genau diese Eigenschaft als Negativbefund geprüft („CRLF, Ziel ohne CR") — sie
  ist mit der Heilung verloren gegangen.
- `klasse`: „Trim-Wechsel schleppt Steuerzeichen ins Befund-Ziel"

### N-4 — Der neu zugesagte Vorrang der strengen Lesart hat keine Assertion

- `kategorie`: MEDIUM
- `quelle`: [§`DC-FA-SPAN-001.a`](../../spec/spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung)
  Schritt 3 („Ist die strenge Lesart offen, steht der Befund an **ihrer**
  Öffnungszeile"); [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md)
  Entscheidung 3
- `pfad`: `internal/hexagon/core/rules/spans.go:38` gegen
  `internal/hexagon/core/rules/spans_test.go:100`
- `befund`: Kein Testfall unterscheidet die beiden Lesarten in der **Zeilenwahl** —
  in allen 15 Tabellenfällen liegen entweder beide offenen Zustände auf derselben
  Zeile, oder es ist ohnehin nur einer offen. Der Rückbau, der die Reihenfolge der
  beiden Lesarten vertauscht, lässt `make test` grün, obwohl er die gemeldete
  Zeile real verschiebt.
- `verifizierbar`: ja — Rückbau „Vorrang getauscht" angewendet: `make test`
  Exit 0. Aus demselben mutierten Stand ein Image gebaut und gegen eine Datei
  gefahren, in der beide Lesarten an **verschiedenen** Zeilen offen enden (Zeile 1
  Backtick-Öffner, danach zwei Tilden-Zeilen): HEAD-Image meldet Zeile 1, das
  mutierte Image Zeile 5. Arbeitsbaum danach wiederhergestellt.
- `klasse`: „zugesagtes Befund-Feld ohne Assertion"

### N-5 — Die neue Ehrlichkeits-Klausel hängt am seltenen Fall; im gemessenen Regelfall ist sie falsch

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-SPAN-001`](../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
  („Der Befund steht an der Öffnungszeile, wenn der CommonMark-Schluss sie kennt.
  Kippt **dagegen** nur die Parität …");
  [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md) Entscheidung 3 („die
  strenge Lesart kennt die tatsächlich offene Öffnung und hat deshalb Vorrang")
- `pfad`: `spec/lastenheft.md` (Klasse 3, Absatz zur Zeilenwahl) und
  `docs/plan/adr/0050-fence-unclosed-in-spans.md`
- `befund`: Beide Texte schreiben der strengen Lesart zu, die schuldige Öffnung zu
  **kennen**, und schränken die Aussage nur für den Fall ein, dass „nur die
  Parität des naiven Toggles kippt". Sind alle Fence-Zeilen eines Dokuments gleich
  lang und gleichzeichig — laut Bestandsmessung §3 des Slice **jedes** der 776
  gemessenen Dokumente —, fallen beide Lesarten zusammen; die strenge weiß dann
  nichts mehr als die naive und zeigt genauso auf eine Zeile, die niemand
  reparieren muss.
- `verifizierbar`: ja — F-6-Fixture (Öffner ohne Schluss in Zeile 3, danach zwei
  intakte Blöcke aus gleich langen Fences): der Befund steht auf Zeile 14. Ein
  eigens gebautes Image, in dem **nur** die strenge Lesart ausgewertet wird,
  meldet dieselbe Datei ebenfalls auf Zeile 14 — die strenge Lesart ist dort also
  offen und liegt trotzdem 11 Zeilen neben der Reparaturstelle. Arbeitsbaum
  danach wiederhergestellt.
- `klasse`: „Ehrlichkeits-Klausel am seltenen Fall aufgehängt"

### N-6 — Der zweite Teil von F-4 steht unverändert und ist nirgends als Grenze benannt

- `kategorie`: INFO
- `quelle`: [ADR-0050](../plan/adr/0050-fence-unclosed-in-spans.md)
  §Re-Evaluierungs-Trigger; Erst-Report F-4
- `pfad`: `internal/hexagon/core/rules/markdown.go:82` (Toggle ohne
  Einrückungsgrenze)
- `befund`: Eine **gerade** Zahl von Toggles aus Nicht-Fence-Quellen bleibt
  unsichtbar: zwei Fence-Zeilen innerhalb eines vier Leerzeichen eingerückten
  Codeblocks — nach CommonMark literaler Inhalt, kein Fence — verschlucken die
  Prosa dazwischen, ohne dass eine der beiden Lesarten unbalanciert endet. Weder
  Anforderung noch ADR nennen diese Restklasse; die Grenze ist dort ausschließlich
  über den Scan-Scope und die Ziel-Achse beschrieben.
- `verifizierbar`: ja — Fixture mit zwei eingerückten Fence-Zeilen und einem
  kaputten Link dazwischen: Alt- und Neu-Image je **0 Befunde, Exit 0**;
  Gegenprobe ohne die beiden Zeilen meldet `target-missing`.
- `klasse`: „Paritätszählung statt Paarung — legale Verschachtelung kippt beide
  Richtungen"

## Negativbefunde

- geprüft, ohne Befund: **Semantik-Erhalt des Tabellen-Lesers (F-1-Kernrisiko).**
  Der Refactoring-Schnitt ist eine wörtliche Substitution: die entfernte lokale
  Marker-Funktion ist bis auf den Parameternamen zeichengleich zur neuen
  geteilten, und das entfernte inline-Schluss-Prädikat prüft dieselben drei
  Bedingungen wie das neue geteilte (gleiches Zeichen, mindestens gleiche Länge,
  dahinter nur Whitespace) auf derselben Markierung. Empirisch abgesichert:
  17 Fence-Randfälle — Backtick gegen Tilde, längere gegen kürzere Folge,
  Infozeile mit Backtick, Schluss mit Text dahinter, Schluss mit
  Trailing-Whitespace, Schluss eingerückt mit Leerzeichen und mit Tab, Schluss
  mit geschütztem Leerzeichen eingerückt, Fence im Blockquote, legale
  Verschachtelung, Öffner ohne Schluss — je gegen Alt- und Neu-Image: **17 von 17
  identisch**. `make completeness-check` und `make trace` über das eigene Repo
  unverändert grün.
- geprüft, ohne Befund: **False-Positive-Suche über 28 Fixtures.** Eingerückte
  Fences in Listen, Fences in Blockquotes, vier Leerzeichen eingerückte
  Beispielblöcke, Tilden-Fences gleicher Länge, Infozeilen mit Sonderzeichen und
  mit Backtick, sehr lange Fence-Läufe (20 Zeichen, symmetrisch), leere Datei,
  Datei nur aus zwei Fence-Zeilen, Setext-Überschriften, YAML-Frontmatter, Fence
  hinter einer Tabelle, Tab-eingerückte Fences, CRLF mit sauber geschlossenem
  Fence — alle befundfrei, in Alt und Neu gleich. **Neu meldende Fälle sind
  ausschließlich echte Asymmetrien**: kürzerer Tilden-Schluss auf längeren
  Tilden-Öffner, kürzerer Backtick-Schluss auf 20-Zeichen-Öffner, Schlusszeile mit
  Text dahinter — in allen dreien überspringt der Tabellen-Leser tatsächlich den
  Rest der Datei.
- geprüft, ohne Befund: **Bestand des Ökosystems unter der geweiteten Lesart.**
  Das eigene Repo über `make doc-check`: 349 Dateien, 0 Befunde. Die beiden
  Schwester-Repos, die die offenen Change Requests gestellt haben, mit einem
  reinen `spans`-Profil: 224 bzw. 222 Dateien, je 0 Befunde. Die Zusage „der
  Bestand bleibt bei null" hält also auch für die zusätzlich bewachte Lesart —
  einschließlich der Klassen, die die Bestandsmessung §3 gar nicht als Spalte
  führt (Schlusszeile mit Text dahinter).
- geprüft, ohne Befund: **Zustands-Führung der neuen Struktur.** Beide
  Lesart-Zustände werden je Datei innerhalb der Prüffunktion angelegt; es gibt
  keinen Träger über Dateigrenzen. Zeichen und Länge werden beim Schließen nicht
  zurückgesetzt, aber ausschließlich im offenen Zustand gelesen und beim nächsten
  Öffnen vollständig überschrieben; Zeile und Text ebenso. Ein Zustand, der nie
  zurückgesetzt wird und dadurch wirkt, existiert nicht.
- geprüft, ohne Befund: **Aufwand bei fenceloser Datei.** Beide Schritt-Funktionen
  brechen am ersten Zeichen der getrimmten Zeile ab; die einzige nennenswerte
  Arbeit ist das Zerlegen in Zeilen, und das war vor der Heilung identisch.
  `make bench` über 1000 generierte Dateien: Median 651 ms gegen ein
  Pass-Kriterium von 5000 ms.
- geprüft, ohne Befund: **Vollständigkeit der „zwei Lesarten"-Behauptung im
  Schluss-Verhalten.** Im Produkt existieren genau zwei Schluss-Regeln: der naive
  Toggle (Vorverarbeitung, Diagramm-Fences, Closure-Abschnitt) und der
  längenabgeglichene Schluss (Tabellen-Leser). Eine dritte Schluss-Regel gibt es
  nicht. Abweichend ist nur die **Trimmung** vor dem Toggle — das ist N-1 und kein
  zusätzlicher Automat.
- geprüft, ohne Befund: **Grund-Code-Lockstep und Klartext.** Der Grund-Code steht
  unverändert in der Aufzählung, in den Klartexten und in der §4-Tabelle der
  Spezifikation; `--doctor` gibt ihn gruppiert mit Fundstelle aus,
  `--repair-broad` bricht nicht.
- geprüft, ohne Befund: **Import-Regeln und statische Analyse.** Die neuen
  geteilten Prädikate liegen im Regel-Paket und werden vom App-Paket konsumiert —
  Richtung unverändert. `make lint`, `make arch-check` und `make semgrep` grün
  (55 Regeln, 0 Findings).
- geprüft, ohne Befund: **Referenz-Richtung (SDP).** Der Heilungs-Commit trägt
  keinen Provenance-Marker und keinen neuen Abwärtsverweis aus ADR oder Spec auf
  Planning-Artefakte; der im Commit erwähnte, zwischenzeitlich eingebaute
  Abwärtsverweis vom Lastenheft auf eine ADR existiert im HEAD-Stand nicht mehr
  (`make doc-check` grün, das Referenzmatrix-Modul ist dort aktiv).
- geprüft, ohne Befund: **Versions- und Historien-Pflege.** Lastenheft-Version auf
  0.52.1 gehoben, Historien-Zeile in beiden Spec-Dokumenten ergänzt, ADR-Geschichte
  um den Überarbeitungs-Eintrag erweitert, Slice §6 mit den Ausgängen versehen.
- geprüft, ohne Befund: **Methodik-Behauptung des Slice.** Die Aussage, der erste
  Mutationsanlauf sei durch Zurücksetzen auf den Commit-Stand methodisch kaputt
  gewesen, ist plausibel und wurde durch das hier gefahrene Kopie-Verfahren
  ersetzt; die daraus abgeleitete Zahl „sechs Rückbauten, alle rot" ist
  eigenständig bestätigt (siehe Tabelle oben).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 4 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Wächter und Bewachtes trimmen verschieden ·
Modul-Grenze nur auf der Quell-Achse gedacht · Trim-Wechsel schleppt
Steuerzeichen ins Befund-Ziel · zugesagtes Befund-Feld ohne Assertion ·
Ehrlichkeits-Klausel am seltenen Fall aufgehängt · Paritätszählung statt Paarung
— legale Verschachtelung kippt beide Richtungen

**Wiederholungs-Signal.** Drei dieser Klassen sind **zum zweiten Mal** in
derselben Sitzung dran (Trim-Divergenz, Ziel-Achsen-Grenze, Zusage ohne
Assertion). Nach dem Reviewer-Skill ist die dritte Wiederholung das
Steering-Loop-Signal; die Klassen gehören deshalb bei der Closure §7 in das
Beobachtungs-Register, nicht nur in diesen Report.

## Verdikt

**Merge-blockierend:** ja — ein HIGH und vier MEDIUM.

Die Heilung ist im Kern echt und nicht bloß plausibel. Der Wächter wertet
tatsächlich beide Schluss-Lesarten aus, die strenge ist real aus ihrem
Konsumenten herausgezogen worden, und der Tabellen-Leser verhält sich auf 17
Randfällen byte-gleich wie vorher — das war das größte Regressionsrisiko des
Refactorings und es hat sich nicht materialisiert. Die zusätzliche Befundfläche
erzeugt auf 795 realen Dokumenten in drei Repos **kein** einziges Falsch-Positiv,
und die neu meldenden Fixtures sind samt und sonders Dokumente, bei denen ein
Modul den Rest der Datei wirklich überspringt. F-1, F-7 und der Kern von F-5 sind
sauber geschlossen.

Blockierend ist, dass die Heilung an derselben Stelle zu früh aufhört wie die
Erstfassung. F-2 wurde gegen die **Vorverarbeitung** angeglichen, nicht gegen den
Konsumenten, dessen Silent-Grün den Slice ausgelöst hat: die Closure-Prüfung
trimmt unicode-weit, und der Anlass-Fall dieses Slice — Floskel hinter offenem
Fence, alles grün — ist mit einem geschützten Leerzeichen unverändert
reproduzierbar, an einem Closure-Gate. Dieselbe Bewegung wiederholt sich in
kleinerem Maßstab: die Ziel-Achsen-Grenze zählt vier Module auf und lässt
ausgerechnet das mit der stillen Richtung weg, die neue Vorrang-Zusage hat keine
Assertion, und die neue Ehrlichkeits-Klausel beschreibt den seltenen Fall statt
des gemessenen Regelfalls. Dazu kommt eine handfeste kleine Regression: das
Befund-Ziel trägt auf CRLF-Dokumenten jetzt ein Steuerzeichen, das die
Befundzeile im Terminal unlesbar macht — eine Eigenschaft, die der Erst-Report
noch ausdrücklich als in Ordnung geprüft hatte.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist
ein Lauf-Beleg (dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) und
ersetzt keine Verifikation — DoD- und Spec-Konformität prüft der Verifier
separat.
