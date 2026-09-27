# Review-Report — slice-232 (`links`: Ziel hinter dem Zeilenumbruch nach `](`), R2

**Review-Art:** Code (gegen ADR-0092, Lastenheft 0.92.1/Spezifikation, Hard Rules,
und explizit: *trägt die R1-H1-Korrektur wirklich?*).

**Gegenstand:** `slice-232`, Commit-Range `45d55827..bb8d6fe8`, fünf Commits;
Fokus dieser Runde die neuesten zwei (`353b310f` Vertrag/ADR-0092,
`bb8d6fe8` Code-Fix). R1 (`docs/reviews/2026-09-27-slice-232-links-zeilenumbruch-review-r1.md`)
deckte die ersten drei Commits ab (Lifecycle-Move, ADR-0091, Feature-Code) und
fand zwei HIGH (R1-H1, R1-H2), ein MEDIUM (R1-M1), ein LOW (R1-L1).

**Skill:** `.harness/skills/reviewer.md` @ `8dcc0452` (Version 1.16.0, 2026-09-07)
— unverändert gegenüber R1.

**Modell-ID:** `claude-sonnet-5`.

**Datum:** 2026-09-27.

**Eingangs-Kontext:**

- R1-Report vollständig gelesen — insbesondere der exakte Rot-Beleg-Testfall
  für R1-H1 (unbalanciertes `[` + spätere unabhängige `](…)`-Sequenz).
- `docs/plan/adr/0091-…md` (Superseded-Übergang geprüft) und
  `docs/plan/adr/0092-…md` (neue ADR, vollständig gelesen).
- `docs/plan/planning/in-progress/slice-232-…md` (§1–§8).
- `spec/lastenheft.md` (`DC-FA-LINK-001`, 0.92.1) und `spec/spezifikation.md`
  (`DC-FA-LINK-001.a` Schritt 3), beide Diffs gegen 0.92.0 gelesen.
- `AGENTS.md` §3 (§3.5 ADR-Immutabilität, §3.7 Kommentar-Klassen) und §5.
- Vollständiger Diff `git diff 377311d8 bb8d6fe8` gelesen (nicht nur die
  Commit-Botschaft) sowie der komplette, geänderte Bereich von
  `internal/hexagon/core/rules/markdown.go`.

**Gates unabhängig nachgefahren** (Docker/`make`, kein Host-Go — AGENTS.md §3.1):

| Gate | Ergebnis |
|---|---|
| `make test` | grün — alle Pakete `ok` |
| `make gates` (voll, alle zehn Glieder) | grün — Coverage `94.30%` (Schwelle 93 %), `semgrep` 0 Befunde/55 Regeln/65 Dateien, `doc-check` `840 Datei(en) geprüft, 0 Befund(e)` |
| `make adr-check` (`RANGE=45d55827..bb8d6fe8`) | grün — `840 Datei(en) geprüft, 0 Befund(e)` über die ganze Slice-Range |
| Runtime-Image gegen dieses Repo (`docker run --network none`, Default-Module) | `840 Datei(en) geprüft, 0 Befund(e)` — reproduziert |

---

## Zentraler Prüfpunkt 1: Ist die R1-H1-Fehlerklasse strukturell ausgeschlossen?

**Ja — verifiziert am Code, nicht nur am Kommentar, und durch eigene
adversarielle Fälle bestätigt.**

`matchBracket(s, start, '[', ']')` in `parseLinkAt` (Zeile ~591) wird
**ausschließlich mit `s`** aufgerufen — niemals mit `full` (der
next-erweiterten Variante). `full`/`next` kommen strukturell erst **danach**
ins Spiel, ausschließlich für die Adress-Klammer
(`matchBracket(full, textEnd+1, '(', ')')`, Zeile ~599). Es gibt im
gesamten Kontrollfluss von `parseLinkAt` keinen Pfad, auf dem die
Linktext-Klammer mit `next`-erweitertem Text matchen könnte — `textEnd`
wird vor der `next`-Verzweigung berechnet und ist bereits fixiert. Das ist
kein Zufall des Testfalls, sondern eine strukturelle Eigenschaft der
Funktion: Eliminierung durch Konstruktion, nicht durch Vermeidung.

**Fünf eigene adversarielle Fälle** wurden per temporärem Testfile
(`internal/hexagon/core/rules/zz_r2_temp_test.go`, `go test` im
gepinnten `golang:1.27.1`-Image, danach entfernt — Arbeitsbaum clean,
`git status` bestätigt leer) gegen den **aktuellen** Stand geprüft:

1. **Chained Spillover** (Zeile 1 spillt in Zeile 2, Zeile 2 spillt nach
   Abzug des verbrauchten Präfixes erneut in Zeile 3): beide Links korrekt
   und unabhängig gefunden (`{Line:1 Target:teil1 Text:a}`,
   `{Line:2 Target:teil2 Text:b}`).
2. **Spillover, dann eigener vollständiger Link auf dem Rest derselben
   Zeile** (Zeile 2 trägt sowohl den verbrauchten Adress-Rest als auch
   `[c](d.md)`): beide korrekt gefunden.
3. **Mehrere unbalancierte `[` vor einem echten, einzeiligen Link**
   (`a[i und b[j und [echt](ziel.md)`): nur der echte Link wird gefunden,
   keine Fusion.
4. **Unbalanciertes `[` unmittelbar vor einem Link, dessen Adresse selbst
   spillt**: nur der echte, spillende Link wird gefunden, korrekt.
5. **Zwei vollständige Links auf einer Zeile, wobei nur der zweite
   spillt**: beide korrekt und unabhängig gefunden — der frühe `break`
   nach einem Spillover-Fund verliert nichts, weil nach einer
   spillenden Adresse in der aktuellen Zeile strukturell kein weiterer
   Link-Start mehr folgen kann (die Tiefe der Adress-Klammer bleibt bis
   zum Zeilenende > 0).

**Rot-Beleg unabhängig reproduziert:** `TestExtractLinks_UnbalancierteKlammerVerschmilztNicht`
wurde gegen die ALTE (`377311d8`) Fassung von `markdown.go` gefahren (per
`git worktree`, Testdatei aus `bb8d6fe8` kopiert, `pins.go`/`sources.go`
mussten dafür **nicht** zurückgesetzt werden, weil der Test nur die
öffentliche `ExtractLinks`-Signatur nutzt, die auf beiden Ständen
identisch ist). Ergebnis: **fällt** mit exakt dem in ADR-0092 behaupteten
Fusions-Befund (`Target:"echt.md"`, `Text:"i steht in der Doku.\nEin Verweis"`)
— an der neuen Fassung grün. Damit ist der Rot-Beleg selbst (nicht nur
seine Behauptung) verifiziert.

## Zentraler Prüfpunkt 2: Consumed/Spillover-Buchführung

Geprüft wie oben (Fälle 1–2, 5) — Kettenbildung über zwei aufeinanderfolgende
Spillover-Fälle funktioniert korrekt, ein eigener Link auf dem nach Abzug
verbleibenden Rest einer Zeile wird gefunden, und ein früher `break` nach
einem Spillover-Treffer verliert keinen nachfolgenden Fund auf derselben
Zeile (es kann strukturell keinen geben). Kein Befund.

**Ein zusätzlicher, in der Aufgabenstellung nicht ausdrücklich benannter
Fall wurde konstruiert und ist dokumentationswürdig (siehe R2-INFO-1
unten):** Eine Adresse, die hinter `](` unvollständig bleibt (kein `)` in
der eigenen Zeile), zieht **jeden** Inhalt der Folgezeile bis zum ersten
`)` heran — unabhängig davon, ob diese Folgezeile inhaltlich als
Fortsetzung der Adresse gedacht war. Fixture: Zeile 1
`Siehe die Dokumentation [hier](`, Zeile 2
`Wie in Kapitel 3 beschrieben) gilt das für alle Fälle.` (kein `(` in
Zeile 2 vor der `)`). Ergebnis: `{Line:1 Target:"Wie" Text:"hier"}` — ein
Link mit sinnlosem, abgeschnittenem Ziel entsteht aus einem eigentlich
unvollständigen/fehlerhaften Link-Versuch plus einer inhaltlich
unabhängigen Folgezeile, die zufällig eine schließende Klammer trägt. Das
ist **keine Wiederkehr der R1-H1-Klasse** (die Linktext-Klammer bleibt
zeilenlokal, der Linktext selbst ist korrekt), sondern eine engere,
symmetrische Variante auf der Adress-Seite — und sie ist die **exakte,
vom Change Request verlangte Form** (ADR-0092: „Erst danach, wenn `[…](`
innerhalb der Zeile feststeht, darf die Adress-Klammer … um genau eine
Zeile … verlängert werden"), nicht ein Konstruktionsfehler. Weder
ADR-0091 noch ADR-0092 noch die Spezifikation benennen diese Kehrseite
ausdrücklich (dass ein **unvollständiger** Link — Autor hat die Klammer
schlicht nicht geschlossen — durch einen zufälligen `)` in der
Folgezeile zu einem falschen, aber gefundenen Ziel wird, statt wie vor
ADR-0091 gar nicht erkannt zu werden). Da dies eine **bewusst
gewählte, vom Auftraggeber autorisierte Zuschnitts-Entscheidung** ist
(Reichweite: genau ein Zeilenumbruch, Adress-Klammer) und kein
Implementierungsfehler, werte ich es als **INFO**, nicht als Finding
gegen die Korrektheit dieser Runde — siehe Negativbefunde/INFO unten.

## Zentraler Prüfpunkt 3: Bestandsmessung nachgefahren

- Default-Module (`make doc-check`): **840 Datei(en) geprüft, 0 Befund(e)** —
  reproduziert, für Alt- (`377311d8`) und Neu-Image (`bb8d6fe8`) identisch.
- `--enable external --enable tracked` (Netz aus, `--network none`):
  Alt-Image **108**, Neu-Image **107** Befunde auf demselben (aktuellen,
  840-Datei-)Bestand. **Alle** 107/108 Befunde sind `external-status`
  (Netzwerk-nicht-erreichbar, ein reines Artefakt der netzlosen
  Sandbox — `tracked` allein liefert 0 Befunde). Der Unterschied von
  genau einem Fund ist erklärt und geprüft: `tools/archive-wave/README.md:3`
  trägt einen Link, dessen **Linktext** über eine Zeile umbricht
  (`[Baseline-Regelwerk … §Wellen-Closure-Prozedur,\nSchritt 4](https://…)`).
  Die alte (R1-)Fassung erkannte diesen Link als unbeabsichtigten
  Nebeneffekt der Absatz-Faltung (in ADR-0091 als „unpromised, aber
  harmlos mitgezogen" benannt); die neue Fassung erkennt ihn **nicht
  mehr**, weil die Linktext-Klammer jetzt strikt zeilenlokal ist — exakt
  der in ADR-0092 dokumentierte Rückbau („Linktext-Umbruch … strukturell
  ausgeschlossen … keine neue Grenze, sondern derselbe bereits benannte
  Fall"). Das entspricht dem Verhalten **vor** slice-232 überhaupt
  (Zeilen-only-Matching fand diesen Link auch vor ADR-0091 nicht). Kein
  neuer Fehlfund, keine verlorene Warnung (der Link selbst ist gültig und
  erreichbar).

**Aber:** Die Commit-Botschaft von `bb8d6fe8` formuliert das Ergebnis als
„840 Dateien, 0 Befunde vor und nach der Aenderung" für **„alle sechs
ExtractLinks-Konsumenten (inkl. external/tracked)"** — das ist so, wie
geschrieben, nicht die tatsächlich gelaufene Messung: Mit `external`/`tracked`
aktiv sind es 108 → 107, nicht 0 → 0. Siehe Finding R2-M1.

## Zentraler Prüfpunkt 4: `pins.go`/`sources.go`

Beide rufen jetzt `forEachLink(text, "", func(...) {...})` mit explizit
leerem drittem Argument auf (verifiziert per Diff, Zeilen wie oben
zitiert). Da `next == ""` den gesamten Spillover-Zweig in `parseLinkAt`
(`if !ok && next != ""`) nie betritt, ist ihr Verhalten **byte-identisch**
zum Aufruf vor ADR-0091/-0092 — keine Verhaltensänderung, nur eine
Signatur-Anpassung. Verifiziert am Code, nicht nur behauptet.

## Zentraler Prüfpunkt 5: ADR-Form/Prozess

- **ADR-0091 Supersession sauber:** `git diff 377311d8 bb8d6fe8 -- docs/plan/adr/0091-….md`
  zeigt **ausschließlich** die `**Status:**`-Zeile (`Accepted` →
  `Superseded by ADR-0092`) und einen neuen `## Geschichte`-Anhang. Kein
  Byte des Kerns (Kontext, Entscheidung, Alternativen, Konsequenzen,
  Fitness Function, Re-Evaluierungs-Trigger) verändert — konform zu
  AGENTS.md §3.5 und dem erlaubten Status-Übergang.
- **`make adr-check` unabhängig nachgefahren** über die ganze Slice-Range
  (`RANGE=45d55827..bb8d6fe8`): grün, `840 Datei(en) geprüft, 0 Befund(e)`.
- **ADR-0092 Form:** `Supersedes: ADR-0091` gesetzt; drei verglichene
  Alternativen mit Trade-offs; Fitness-Function-Tabelle inkl. Rot-Beleg-
  Zusage; `## Re-Evaluierungs-Trigger` (zwei Bedingungen); `Schärft:`
  zeigt aufwärts auf `DC-FA-LINK-001`. Vollständig nach
  `modul-04-adrs.md` §Ziel-Form.
- **ADR-Index (`docs/plan/adr/README.md`)** aktualisiert: ADR-0091-Zeile
  auf „Superseded by ADR-0092" korrigiert, neue ADR-0092-Zeile ergänzt.
- **Referenzrichtung:** `grep -n "ADR-009[12]\|slice-232" spec/lastenheft.md spec/spezifikation.md`
  — kein Treffer; beide Spec-Straten nennen ADR/Slice nicht im Körper.
- **MR-032-Decken-Regel:** neue Historie-Zeilen (Lastenheft 0.92.1,
  Spezifikation) nennen weder ADR-Nummer noch Slice-ID, „Verweis"-Spalte
  bleibt `—` (Lastenheft) — konform.

## Zentraler Prüfpunkt 6: Rot-Beleg gegen den alten Stand

Siehe Prüfpunkt 1 — unabhängig nachgefahren, bestätigt exakt den in
ADR-0092/Commit-Botschaft behaupteten Fehlschlag am alten Stand.

---

## Weitere Prüfungen

- **`make gates` und `make test`** unabhängig nachgefahren (s. o.), beide grün
  mit echter Ausgabe.
- **§3.7 in NEUEN Dateien dieser Runde:** Vier `// slice-232 …`-Präfixe in
  `cli_acceptance_test.go` (R1-H2) sind vollständig entfernt — verifiziert per
  Diff, alle vier Testabsichten bleiben inhaltlich unverändert erhalten.
  **Aber:** zwei **neue** Kommentare dieser Runde tragen einen
  **Review-Befund-Marker** — siehe R2-H1.
- **Zweite String-Konsumenten-Zählung (Doku-Präzision):** ADR-0092 spricht an
  zwei Stellen (Zeilen 85, 128) weiterhin von „vier String-Konsumenten
  (`ids`, `pins`, `--repair`, `planning`)" — dieselbe Aufzählung wie
  ADR-0091. Tatsächlich rufen aber **fünf** Stellen `forEachLink`/
  `ExtractLinkSpans` mit leerem `next` auf: `ids.go` (2 Aufrufe),
  `repair.go`, `planning_observations.go` — alle vier via
  `ExtractLinkSpans` — **und zusätzlich** `pins.go` **und** `sources.go`
  (Modul `sources`, `DC-FA-SRC-001`), beide **direkt** über `forEachLink`,
  nicht über `ExtractLinkSpans`. Die Konsequenzen-Sektion derselben ADR
  (Zeile 103) benennt `sources.go` korrekt als eigenen Aufrufer — die
  Zählung „vier" in der Entscheidungs- und Trigger-Sektion ist damit
  intern inkonsistent zur eigenen Konsequenzen-Sektion. Funktional folgenlos
  (Code für `sources.go` korrekt auf `next=""` umgestellt, verifiziert
  in Prüfpunkt 4) — reine Dokumentations-Präzision, dieselbe Klasse wie
  R1-L1, diesmal aber in einem in dieser Runde neu geschriebenen
  Dokument. Siehe R2-L1.
- **Commit-Zerlegung:** `353b310f` (Vertrag) ändert ausschließlich
  `docs/plan/adr/`, `spec/*.md` — kein Code. `bb8d6fe8` (Fix) ändert
  ausschließlich `internal/` — kein Vertragstext. Sauber.
- **Bilder/Grenzen (§1 „Ausdrücklich NICHT"):** unverändert gegenüber R1,
  keine neuen Berührungen in dieser Runde.

---

## Findings

| # | Kategorie | Quelle | Pfad | Befund | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| R2-H1 | HIGH | `AGENTS.md` §3.7 (Kommentar-Klassen, „keine Review-Befund-Marker") / Prüffrage 6 | `internal/hexagon/core/rules/markdown.go:577` (`parseLinkAt`-Doku-Kommentar), `internal/hexagon/core/rules/markdown_test.go:154` (Kommentar vor `TestExtractLinks_UnbalancierteKlammerVerschmilztNicht`) | Zwei **neue** Kommentare dieser Runde tragen einen Review-Befund-Marker im Fließtext: `„… zu einem erfundenen // Link (Review-Befund R1-H1, ADR-0091)."` in **Produktionscode** und `„// ADR-0091 (Review R1-H1): ein unbalanciertes …"` im Test. §3.7 listet „Review-Historie und keine Review-Befund-Marker" explizit als verbotene Klasse; erlaubt ist Herkunft nur als **ein** Feld nach dem Schema `DC-*`/`ADR-*`/`MR-*`/`seit welle-<NN>` — „Review-Befund R1-H1" fällt unter keine davon. Die Bestandsgrenze deckt nur *vor* der Einführung geschriebene Kommentare; beide sind Neuzugänge dieses Diffs (verifiziert: in `377311d8` noch nicht vorhanden). Der ADR-0091-Teil des Verweises wäre für sich zulässig — das Problem ist ausschließlich der eingeschobene Review-Befund-Bezug. | ja — Diff zeigt beide Zeilen als reine Neuzugänge; einfache Entfernung von „(Review-Befund R1-H1, " bzw. „(Review R1-H1): " ändert an der technischen Aussage nichts | review-befund-marker-in-kommentar |
| R2-M1 | MEDIUM | `AGENTS.md` §5 (Commit-Botschaft-Overclaim) / Prüffrage 8 | Commit `bb8d6fe8` (Botschaft) | Die Botschaft „Bestandsmessung wiederholt gegen alle sechs ExtractLinks-Konsumenten (inkl. external/tracked …): 840 Dateien, 0 Befunde vor und nach der Aenderung — schliesst R1-M1" ist als geschlossene Aussage über **alle sechs** Konsumenten formuliert. Die tatsächlich unabhängig reproduzierte Messung mit `--enable external --enable tracked` liefert **108 → 107** Befunde (alle `external-status`, ein reines Sandbox-Artefakt), nicht 0 → 0. Der „0 Befunde"-Wert stimmt nur für die **vier** in `.d-check.yml` aktiven Standard-Module — exakt die Teilmenge, deren Unvollständigkeit R1-M1 ursprünglich bemängelte. Die Aussage überdehnt die gelaufene Messung auf eine Menge, für die sie nicht zutrifft, obwohl die zugrunde liegende Arbeit (Vergleich Alt-/Neu-Binary über alle sechs Konsumenten) tatsächlich stattfand und der Delta-Fund selbst korrekt verstanden/benannt ist (Prüfpunkt 3). | ja — `docker run --rm --network none -v $(pwd):/repo:ro <image> --enable external --enable tracked` gegen Alt- und Neu-Image liefert 108 bzw. 107, nicht 0 | messmethode-botschaft-overclaim |
| R2-L1 | LOW | Doku-Präzision, ADR-0092 §Entscheidung/§Re-Evaluierungs-Trigger | `docs/plan/adr/0092-links-zeilenumbruch-begrenzter-lookahead.md` (Zeilen 85, 128) | Die ADR zählt weiterhin „vier String-Konsumenten (`ids`, `pins`, `--repair`, `planning`)", obwohl ihre eigene Konsequenzen-Sektion (Zeile 103) `sources.go` korrekt als zusätzlichen, direkten `forEachLink`-Aufrufer neben `pins.go` benennt — tatsächlich sind es fünf Aufrufstellen (drei über `ExtractLinkSpans`, zwei direkt: `pins.go`, `sources.go`). Funktional folgenlos (Code für `sources.go` korrekt behandelt, verifiziert), aber eine interne Inkonsistenz innerhalb desselben, in dieser Runde neu geschriebenen Dokuments — dieselbe Klasse wie R1-L1, hier aber kein Bestand, sondern ein Neuzugang. | ja — `grep -n "vier\|sources" docs/plan/adr/0092-….md` zeigt den Widerspruch zwischen Zeile 85/128 und Zeile 103 | adr-attribution-ungenau |

## Negativbefunde (geprüft, ohne Befund)

- **R1-H1-Fehlerklasse strukturell ausgeschlossen:** verifiziert am Code
  (`matchBracket` für `[`/`]` operiert nie auf `full`) und an fünf eigenen
  adversariellen Fällen plus dem unabhängig reproduzierten Rot-Beleg gegen
  den alten Stand (Prüfpunkt 1).
- **Consumed/Spillover-Buchführung:** Kettenbildung über mehrere Zeilen,
  eigener Link nach verbrauchtem Präfix, kein Verlust durch frühen
  `break` — alle drei Fälle korrekt (Prüfpunkt 2).
- **`pins.go`/`sources.go` byte-identisches Verhalten:** verifiziert am
  Code (`next=""` schließt den Spillover-Zweig strukturell aus,
  Prüfpunkt 4).
- **Bestandsmessung (Default-Module):** 840 Dateien, 0 Befunde vor und
  nach der Änderung — unabhängig reproduziert, korrekt für die
  tatsächlich aktive Modul-Menge (siehe R2-M1 für die Reichweiten-Kritik
  an der Formulierung).
- **ADR-0091-Supersession (§3.5):** ausschließlich Status-Feld + Geschichte-
  Anhang geändert, Kern byte-identisch; `make adr-check` über die ganze
  Slice-Range unabhängig grün.
- **ADR-0092-Form:** drei Alternativen, Fitness Function mit Rot-Beleg-
  Zusage, Re-Evaluierungs-Trigger, `Supersedes`-Feld, `Schärft:` aufwärts
  — vollständig.
- **Referenzrichtung/MR-032:** kein ADR-/Slice-Token im Körper der
  Spec-Straten; neue Historie-Zeilen ohne ADR-/Slice-Nennung, `Verweis`
  bleibt `—`.
- **R1-H2 vollständig behoben:** alle vier `// slice-232 …`-Präfixe aus
  `cli_acceptance_test.go` entfernt, Testabsicht erhalten.
- **Commit-Zerlegung:** Vertrags-Commit ohne Code, Fix-Commit ohne
  Vertragstext — sauber.
- **Gates:** `make test`, voller `make gates` (zehn Glieder) und
  `make adr-check` über die ganze Slice-Range unabhängig nachgefahren —
  alle grün mit echter Ausgabe (Coverage 94,30 %, semgrep 0/55/65,
  doc-check 840/0).

## INFO

- **R2-INFO-1 — Adress-Spillover kann mit inhaltlich unabhängiger
  Folgezeile fusionieren, wenn die Adresse unvollständig bleibt.** Ein
  Link mit `](` am Zeilenende und **keiner** schließenden Klammer in
  derselben Zeile zieht den gesamten Inhalt der Folgezeile bis zum
  ersten `)` heran — auch wenn diese Zeile inhaltlich nichts mit einer
  Adress-Fortsetzung zu tun hat (konstruierter Fall in Prüfpunkt 2:
  `[hier](` gefolgt von unabhängiger Prosa mit zufälligem `)` ergibt
  `Target:"Wie"`). Das ist die vom Auftraggeber explizit gewählte
  Reichweite (ein Zeilenumbruch, Adress-Klammer) und kein
  Implementierungsfehler — vor ADR-0091 wäre ein solcher unvollständiger
  Link gar nicht erkannt worden (stille Nicht-Erkennung), jetzt kann er
  zu einem falschen, aber gemeldeten Ziel führen. Keiner der beiden
  Verträge (Lastenheft, Spezifikation) benennt diese Kehrseite
  ausdrücklich; sie ist auch nicht das Gegenstück zu den bereits
  benannten zwei Grenzen (Linktext-Umbruch, Titel-Umbruch). Nicht
  merge-blockierend — kein Failure-Szenario mit realem Schaden auf dem
  aktuellen 838/840-Datei-Bestand (Bestandsmessung 0 einschlägige
  Befunde), und die Reichweite ist eine bewusste, autorisierte
  Zuschnitts-Entscheidung, kein Versehen. Dokumentationswürdig für einen
  künftigen Re-Evaluierungs-Anlass.

## Kategorie-Summary

- HIGH: 1 (R2-H1)
- MEDIUM: 1 (R2-M1)
- LOW: 1 (R2-L1)
- INFO: 1 (R2-INFO-1)

## Verdikt

**Merge-blockierend: ja, wegen R2-H1.**

Begründung: R2-H1 ist trivial behebbar (zwei Kommentar-Umformulierungen,
ADR-0091-Bezug bleibt erhalten, nur der Review-Befund-Marker entfällt),
blockiert aber formal nach §3.7/Prüffrage 6, ohne Ausnahme für Neuzugänge.
Die eigentliche fachliche Frage der Runde — trägt die R1-H1-Korrektur? —
ist mit **ja** beantwortet: Die absatzweise Faltung ist vollständig durch
einen strukturell auf die Adress-Klammer begrenzten Ein-Zeilen-Lookahead
ersetzt, die Linktext-Klammer verlässt ihre Zeile nie, der Rot-Beleg gegen
den alten Stand wurde unabhängig reproduziert, und fünf zusätzliche eigene
adversarielle Fälle zeigen keine Wiederkehr der R1-H1-Klasse. `pins.go`/
`sources.go` sind nachweislich byte-identisch zum Stand vor ADR-0091.
R2-M1 (Botschafts-Überdehnung bei der Bestandsmessung) und R2-L1
(Doku-Ungenauigkeit in ADR-0092) sind für sich beide nicht
merge-blockierend, sollten aber vor Closure durch eine präzisere
Formulierung bzw. eine Korrektur der Konsumenten-Zählung aufgelöst werden.
R2-INFO-1 ist eine benannte, aber nicht blockierende Randbeobachtung zur
Adress-Spillover-Reichweite.
