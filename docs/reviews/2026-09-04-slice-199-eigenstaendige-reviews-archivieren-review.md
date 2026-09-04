# Review-Report: slice-199 — Eigenständige Reviews archivieren — 2026-09-04

**Review-Art:** Unabhängiger Code-/Bestands-Review (Modul 8, Reviewer-Rolle; kein geteilter Kontext mit der Implementierung).
**Gegenstand:** Commit `6c4146f` ("feat(planning): slice-199 -- alle eigenstaendigen Reviews archiviert (slice-199, welle-90)"), Diff gegen `b499d88`. 33 geänderte Dateien: 11 Original-Reviews unter `docs/reviews/` gelöscht, je ein Stub + ein `-archiv.zip` unter `docs/reviews/archiv/` neu.
**Skill/Modell/Datum:** `reviewer.md` v1.13.0, Claude Sonnet 5, Reviewer-Rolle, 2026-09-04.
**Eingangs-Kontext:** `AGENTS.md` (§3.3, §3.7, §5), `.harness/skills/reviewer.md`, Slice-Plan `slice-199` und `slice-198` (`done/`), `tools/archive-wave/stub.go`/`archive.go`/`rewrite.go`/`main.go`, `.d-check.yml` (`ignore-refs`), `docs/plan/adr/0013-pr-ci-und-traceability-gate.md`. Kein Zugriff auf den Implementierungs-Kontext dieser Session.

## Unabhängig nachgefahrene Prüfungen (echte Ausgabe, kein behaupteter Exit-Code)

- `find docs/reviews -maxdepth 1 -name "*.md" | grep -v 'slice-[0-9]'` → leer: kein eigenständiger Review mehr außerhalb `archiv/`.
- `ls docs/reviews/archiv/*.md` → 11, `ls docs/reviews/archiv/*.zip` → 11 — Bestand deckungsgleich mit der im Slice-Plan genannten Zahl.
- `unzip -l` gegen alle 11 `.zip`: je genau eine Datei, Pfad `docs/reviews/<original>.md`, Größe > 0.
- `unzip -p docs/reviews/archiv/2026-06-21-adr-0013-pr-ci-traceability-archiv.zip … | diff - <Original aus b499d88>` → **identisch** (Byte-für-Byte-Stichprobe).
- Alle 11 Stubs gegen `ReviewStub()` (`stub.go:169-178`) geprüft: Form, Zeilenreihenfolge, Archiv-Zeiger-Pfad, `d-check:ignore`-Marker stimmen für jeden der 11 exakt überein.
- Titel-Stichprobe (3 von 11: MR-016/017-R1, Backlog-Schnitt, Lastenheft-Doctor-Repair) gegen die erste Zeile der jeweiligen Originaldatei (`git show b499d88:…`) — Wortlaut deckungsgleich, keine Verstümmelung.
- `git show 6c4146f | grep "^+" | grep -v "^+++"` — volle Diff-Prüfung auf neue Kommentare: **kein** Neuzugang trägt Slice-Nummer-/Plan-Absatz-Provenienz (dieselbe Fehlerklasse wie slice-196/slice-198 F-1); alle 11 `d-check:ignore`-Kommentare sind wortgleich die vom Werkzeug selbst erzeugte Form.
- `Hervorgegangen:`-Feld aller 11 Stubs gegen eine eigene Regex-Extraktion (`DC-*`/`ADR-*`, fenced+inline-code-bereinigt wie `ExtractSurvivingIDs`) der jeweiligen Originaldatei gegengeprüft (siehe F-2).
- Repo-weite `git grep` nach allen 11 Original-Dateinamen: kein lebender Markdown-**Link** darauf; ein plain-text-Pfadverweis gefunden (siehe F-3).
- `.d-check.yml` nach allen 11 Dateinamen und den beiden Datums-Präfixen `2026-08-09`/`2026-08-10` durchsucht (siehe F-1/F-2 unten — hier als Fund gegen die neue Rewrite-Grenze, nicht als eigener Punkt doppelt gezählt).
- **Nicht selbst nachgefahren:** `make gates` und `make fullbuild`. Die Commit-Botschaft behauptet zehn grüne Gates und `--require-complete` mit 0 Trace-Waisen; das ist plausibel (reine Doku-Löschung/-Kürzung ohne Code-Änderung, `Hervorgegangen:` als Gegenprobe füllt genau die Lücke, die `--require-complete` prüft) und wird gemäß Aufgabenstellung als bereits gelaufener Beleg akzeptiert, aber **nicht** durch einen eigenen Lauf bestätigt — diese Zeile ist bewusst kein Positivbefund.

## Findings

### F-1 · HIGH · `AGENTS.md` §3.3 (git mv + Inhaltsänderung = zwei Commits; enumerierte Ausnahmen) · Commit `6c4146f` insgesamt

**Befund:** `6c4146f` ersetzt für 11 Dateien in einem einzigen Commit sowohl den Ort (`docs/reviews/` → `docs/reviews/archiv/`) als auch den Inhalt (Volltext → 5-Zeilen-Stub). Das ist exakt das Muster, das §3.3 grundsätzlich verbietet ("Wenn eine Datei verschoben **und** der Inhalt umgeschrieben wird" → zwei Commits) — zulässig nur über eine der fünf dort **enumerierten** Ausnahmen (Slice-Lifecycle-Move, Beanspruchung, MR-/Wellen-Lifecycle-Move, Wellen-Archiv-Stub-Move/`MR-059`, Register-Formatmigration/`MR-061`, Wellenloser-Einzel-Slice-Archiv-Move/`MR-062`). Keine dieser fünf trifft den vorliegenden Fall: Der `-review=<datei>`-Modus (`tools/archive-wave`, gebaut in slice-198) ist eine **eigene**, in `AGENTS.md` §3.3 namentlich **nicht genannte** Betriebsart. Weder der Slice-198-Plan/-Closure noch der Slice-199-Plan/-Commit legen eine solche Ausnahme neu an (kein `MR-06x`-Nachtrag, kein neuer `AGENTS.md`-Absatz).

Das ist derselbe Fehler-Typ, der in diesem Repo bereits einmal einen eigenen Korrektur-Commit brauchte: `9980aa9` ("fix(planning): slice-197 -- MR-062 statt ueberdehnter MR-059-Analogie …") korrigierte genau diese Lücke für den `-slice=<id>`-Modus, indem er die zuvor nur analog (nicht namentlich gedeckte) Berufung auf `MR-059` durch die eigens geschriebene `MR-062` ersetzte. Für den `-review=<datei>`-Modus fehlt dieser Schritt bislang vollständig — weder als eigene `MR-063` noch als sechste Ausnahme-Bullet in §3.3.

Technisch ist die Begründung dieselbe wie bei `MR-059`/`MR-062` (keine Phase mit unverändertem Inhalt, da der Stub den Volltext im selben Akt ersetzt, der ihn verschiebt — Git zeigt entsprechend elf reine `D`/`A`-Paare, keine Renames), aber die Regel selbst verlangt, dass das **benannt** wird, nicht dass es zutrifft: "Ein Wellen-Archivierungs-Commit bleibt deshalb bewusst ein Commit, in der Botschaft **ausdrücklich als solcher deklariert**" (`MR-059`). Die Commit-Botschaft von `6c4146f` deklariert nichts dergleichen und zitiert keine der fünf Ausnahmen.

**Verifizierbar:** ja — `AGENTS.md` §3.3 gegen die fünf dort aufgeführten Bullets lesen (keine nennt `-review`); `grep -n "MR-063\|Review-Archiv-Move" harness/conventions.md` liefert nichts; `git show 6c4146f` zeigt elf `D`+`A`-Paare ohne Rename-Erkennung, ohne begleitende Deklaration in der Botschaft.
**Klasse:** `hard-rule-ausnahme-nicht-benannt`

### F-2 · MEDIUM · `.d-check.yml` `ignore-refs` (geteiltes Referenz-Ventil) · Zeilen 81, 94, 121, 259

**Befund:** Vier `ignore-refs`-Einträge in `.d-check.yml` sind durch diesen Commit **verwaist**:

- `docs/reviews/2026-08-09-*.md` (Zeile 81–93) — kein Original mit diesem Datumspräfix existiert nach dem Commit mehr unter `docs/reviews/` (nur noch der Stub unter `archiv/`, den der Glob nicht trifft).
- `docs/reviews/2026-08-10-*.md` (Zeile 94–100) — dieselbe Lage.
- `docs/reviews/2026-08-09-backlog-schnitt-review.md` (Zeile 121) — namentlicher Eintrag auf exakt eine der 11 archivierten Dateien.
- `docs/reviews/2026-08-31-release-prep-v0.71.0-review.md` (Zeile 259) — dieselbe Lage.

Kein Gate scheitert daran — `ignore-refs` ist reine additive Erlaubnis ohne Lebendigkeits-Prüfung des `in:`-Feldes, ein verwaister Eintrag ist für keinen Sensor sichtbar (dieselbe stille Klasse wie `AGENTS.md` §3.8: das Modul verspricht nur über das, was es scannt — hier über Referenzen, nicht über die Lebendigkeit seiner eigenen Ausnahmeliste). Die Datei selbst dokumentiert für **exakt diese Situation** (Quelldatei einer `ignore-refs`-Zeile wird archiviert) bereits ein etabliertes Muster — siehe Zeilen 211–260, wo die Wellen-/Register-Archivierungen von slice-190/191/195/197 jeweils mit einem eigenen "Tombstone"-Kommentarblock und/oder einer angepassten Zeile nachgezogen wurden. Für die vier hier betroffenen Zeilen fehlt dieser Nachzug.

**Verifizierbar:** ja — `ls docs/reviews/2026-08-09-*.md docs/reviews/2026-08-10-*.md` (Exit 2, kein Treffer) gegen die vier genannten `.d-check.yml`-Zeilen.
**Klasse:** `ignore-refs-eintrag-verwaist-durch-archivierung`

### F-3 · MEDIUM · `docs/plan/adr/0013-pr-ci-und-traceability-gate.md:180` (Geschichte-Tabelle, `Accepted`/immutable)

**Befund:** Die Geschichte-Zeile `| 2026-06-21 | Review (docs/reviews/2026-06-21-adr-0013-pr-ci-traceability.md) eingearbeitet: … |` nennt den Review-Pfad als **Klartext in Klammern**, nicht als Markdown-Link. `RewriteRepo`/`RewriteFile` (`rewrite.go:66-77`) erkennt ausschließlich `](…)`-Linkziele; ein bloßer Pfad in Prosa wird nicht erfasst und blieb unangetastet. Der genannte Pfad existiert seit diesem Commit nicht mehr (Datei liegt jetzt unter `docs/reviews/archiv/…`) — die Zeile ist damit sachlich veraltet, und zwar in einer Weise, die **kein** Modul scannt (weder `links`/`ids` noch `tracked`, da kein Link-Syntax vorliegt).

Erschwerend: `docs/plan/adr/0013-…` ist `Accepted` und nach §3.5 immutable — der Kern-Vergleich von `adr-check` lässt nur `## Geschichte`-**Anhänge** und den `**Status:**`-Übergang zu, keine Änderung an einer bestehenden Geschichte-Zeile. Eine nachträgliche Korrektur des Pfads in Zeile 180 selbst wäre also vermutlich ein Verstoß gegen §3.5; zulässig wäre höchstens ein neuer Anhang, der die Verlegung vermerkt. Weder der Slice-Plan noch der Commit erwähnt diesen Fall.

**Verifizierbar:** ja — `git grep -n "2026-06-21-adr-0013-pr-ci-traceability.md" -- '*.md'` zeigt die Fundstelle außerhalb `docs/reviews/archiv/`; `ls docs/reviews/2026-06-21-adr-0013-pr-ci-traceability.md` (Exit 2).
**Klasse:** `plain-text-referenz-uebersehen-bei-archivierung`

### F-4 · LOW · Slice-Plan `slice-199` §5, zweiter Punkt ("Hervorgegangen: fängt das strukturell auf")

**Befund:** Eigene Regex-Gegenprobe der Originaldatei `2026-06-21-adr-0013-pr-ci-traceability.md` (vor Archivierung) fand sechs `DC-*`/`ADR-*`-Kennungen; der erzeugte Stub trägt in `Hervorgegangen:` nur `ADR-0004, ADR-0007, ADR-0011, ADR-0012, ADR-0013, DC-QA-02` — **vier fehlen**: `DC-FA-CODE-001`, `DC-FA-DIST-001`, `DC-FA-ID-001`, `DC-FA-MTX-001`. Ursache ist by design: alle vier stehen im Original nur als **bloßer Inline-Code-Span** ohne Link-Label (z. B. `` `DC-FA-DIST-001` ``), und `stripStandaloneInlineCode`/`ExtractSurvivingIDs` (`stub.go:58-107`) entfernt genau solche Spannen bewusst (Schutz gegen erfundene Illustrations-Kennungen, siehe Kommentar dort). Im vorliegenden Fall ist das **folgenlos**: alle vier Kennungen sind repo-weit reichlich anderswo zitiert (README.md, README.de.md, `spec/lastenheft.md` selbst u. a.), also keine Trace-Waise entstanden — die `--require-complete`-Behauptung der Commit-Botschaft bleibt korrekt.

Der Slice-Plan (§5) formuliert die Absicherung dieses Risikos jedoch als: *"`Hervorgegangen:` fängt das strukturell auf, `--require-complete` bestätigt es empirisch"* — das überzeichnet, was `Hervorgegangen:` tatsächlich leistet: das Feld fängt **nicht** jede im Original zitierte Kennung strukturell auf, sondern nur die außerhalb von Fenced-/Inline-Code-Spannen bzw. als Link-Label stehenden. Der tatsächliche Schutz in diesem Slice ist ausschließlich der empirische `--require-complete`-Lauf, nicht die Feld-Konstruktion selbst. Kein Korrektheitsfehler am Werkzeug (das Verhalten ist so beabsichtigt, siehe `stub.go`-Kommentar), aber eine Übergeneralisierung im Plan, die bei einem künftigen Review mit ausschließlich inline-code-zitierten, sonst nirgends verlinkten Kennungen tatsächlich zu einer stillen Trace-Waise führen könnte.

**Verifizierbar:** ja — eigene Extraktion vs. `Hervorgegangen:`-Feld des Stubs; `grep -rl "DC-FA-DIST-001\|DC-FA-ID-001\|DC-FA-MTX-001\|DC-FA-CODE-001"` zeigt reichlich andere Fundstellen (kein Waisen-Nachweis in diesem konkreten Fall).
**Klasse:** `botschaft-ueberdehnt-mechanismus-ueber-tatsaechliche-deckung`

## Positiv geprüft (kein Finding)

- **Vollständigkeit des Bestands:** genau 11 archiviert, 0 eigenständige Reviews außerhalb `archiv/` verblieben — deckt sich mit der im Plan genannten Zahl.
- **Zip-Integrität:** alle 11 `.zip` entpackbar, je genau eine Datei, Stichprobe byte-identisch mit dem Original vor der Löschung.
- **Stub-Form:** alle 11 exakt nach `ReviewStub()` (Titel, Archiv-Zeiger, `Archiviert:`-Platzhalter, `Hervorgegangen:`-Zeile samt `d-check:ignore`-Marker) — keine Abweichung.
- **Titel-Extraktion:** an drei Stichproben unterschiedlicher Überschriftenformen ("Review-Report: …", "Review — …", "Review Release-Prep …") verbatim erhalten, keine Verstümmelung des führenden Worts (das war slice-198s eigentliches Risiko und trat hier nicht auf).
- **§3.7-Kommentar-Provenienz:** kein Neuzugang dieses Diffs trägt Slice-/Plan-Absatz-Herkunft — die Fehlerklasse, die slice-196 und slice-198 je einmal HIGH kostete, trat in slice-199 **nicht** ein.
- **Repo-weite Markdown-Links** auf die 11 alten Pfade: keine gefunden (der einzige Fund war die Klartext-Stelle in F-3, kein Link).
- **Commit-Botschaft:** bleibt innerhalb der gemessenen Menge (nennt genau die gelaufenen Proben: Voll-Dry-Run, `make gates`, `make fullbuild --require-complete`, 0 Waisen) — keine Verallgemeinerung über das Gemessene hinaus (§5 Hard Rule zu Commit-Botschaften).
- **DoD/§5 des Slice-Plans (Stand `in-progress`):** Checkbox-Stand ehrlich unchecked (Bestandszahl, Review, Verifikation, Closure alle noch offen) — kein verfrühtes Abhaken vor Abschluss.

## Verdikt

**Closure-blockierend: JA — wegen F-1.**

Die eigentliche Archivierungs-Operation ist mechanisch sauber: Bestand vollständig erfasst, alle 11 Stubs und Zips korrekt geformt, keine neue §3.7-Kommentar-Provenienz, kein Korrektheitsfehler im Werkzeug selbst (das war bereits Gegenstand des separaten slice-198-Reviews). F-1 ist jedoch ein Verstoß gegen eine benannte Hard Rule (§3.3) ohne die dafür in diesem Repo etablierte Deckung — genau die Lücke, die hier bereits einmal (slice-197 → `MR-062`) einen eigenen Korrektur-Commit brauchte. Vor Closure sollte entweder eine sechste Ausnahme-Bullet in `AGENTS.md` §3.3 plus zugehöriger `MR-06x`-Nachtrag für den `-review=<datei>`-Modus ergänzt werden (analog `MR-059`/`MR-062`), oder im Report/Closure-Text begründet werden, warum keiner nötig ist.

F-2 und F-3 sind real, aber von keinem Gate erzwungen; sie sollten mindestens in der Closure-Notiz (§9) als bekannte, bewusst nicht behobene Lücken vermerkt werden (F-3 zusätzlich mit dem Hinweis auf die ADR-Immutabilität als Grund, warum eine direkte Korrektur nicht ohne Weiteres möglich ist) — sonst verschwinden sie lautlos, wie es §3.8 grundsätzlich beschreibt. F-4 ist eine Formulierungs-Präzisierung im Plan, nicht blockierend.
