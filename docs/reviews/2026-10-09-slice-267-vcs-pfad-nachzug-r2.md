# Review R2: slice-267, `vcs` lässt einen reinen Pfad-Nachzug durch

- **Review-Art:** Code. Geprüft wird der Fix-Diff `36a99a4c..273bebea` (`fix(vcs)` `38ad8473`,
  `docs(adr)` `273bebea`) gegen den Slice-Plan `slice-267`, gegen die Findings aus R1
  (`2026-10-09-slice-267-vcs-pfad-nachzug-r1.md`), gegen `ADR-0103`, `DC-FA-VCS-001`,
  `MR-032` und die Hard Rules `AGENTS.md` §3.5/§3.7 sowie §5 Regel 13. Die DoD-Abhakung prüft
  dieses Review nicht.
- **Gegenstand:** `slice-267` · `36a99a4c..273bebea`.
- **Skill:** `reviewer.md` @ 1.19.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** `DC-FA-VCS-001` im Lastenheft 0.102.0 (Absatz „Pfad-Nachzug (opt-in)",
  Historie-Zeile 0.102.0); Spezifikation `DC-FA-VCS-001.a` Schritt 4, §2-Zeile
  `vcs.ignore-link-targets`, §8-Historie; `ADR-0103` samt neuer Geschichte-Zeile; Baseline
  `v6.17.0` · `regelwerk/modul-04-adrs.md` §Nachzug ist keine Überschreibung; `MR-032`. Im Code:
  `inlineLinkTargetRE`, `refDefTargetRE`, `blankLinkTargets`, `vcsCore` in
  `internal/hexagon/core/rules/vcs.go`; Kommentar von `VCSConfig.IgnoreLinkTargets` in
  `internal/hexagon/core/model/config.go`; `vcs_pfad_nachzug_test.go`;
  `harness/sensors/adr-check.md` Grenze 4. Vorherige Findings am Modul: R1 zu `slice-267`, R1 zu
  `slice-247`.
- **Proben** (Image `dcr267r2` per `make build IMAGE=dcr267r2` vom Stand `273bebea`; zwei
  Wegwerf-git-Repos im Scratchpad, je Probe eine `Accepted`-ADR mit einer zwischen zwei Commits
  geänderten Zeile, Lauf `--enable vcs --disable links --range HEAD~1..HEAD`, Schlüssel an und aus;
  danach Repos und Image entfernt, Arbeitsbaum unverändert):

  | # | BASE → HEAD | mit Schlüssel | ohne |
  |---|---|---|---|
  | 1 | `\[A\](Nicht) gilt.` → `\[A\](Immer) gilt.` | **still** | Befund |
  | 2 | `` Code: `fns[k](true)` gilt. `` → `` …`fns[k](false)`… `` | **still** | Befund |
  | 3 | `[Status]: Abgelehnt.` → `[Status]: Angenommen.` | still | Befund |
  | 4 | `\| [A](a.md) \| Nicht \|` → `\| [A](x/a.md) \| Immer \|` | Befund | Befund |
  | 5 | `[A](a.md (Nicht erlaubt))` → `[A](x/a.md (Immer erlaubt))` | Befund | Befund |
  | 6 | 3 Leerzeichen + `[x]: a.md` → `…x/a.md` | still | Befund |
  | 7 | 4 Leerzeichen + `[x]: a.md` → `…x/a.md` | Befund | Befund |
  | 8 | `[A](Nicht) gilt.` → `[A](Immer) gilt.` | still | Befund |
  | 9 | `[A](<Nicht erlaubt>)` → `[A](<Immer erlaubt>)` | still | Befund |
  | 10 | `[A](a.md) Nicht [B](b.md)` → `[A](x/a.md) Immer [B](x/b.md)` | Befund | Befund |
  | 11 | `[A](a.md"Nicht")` → `[A](a.md"Immer")` | still | Befund |
  | 12 | `- [x]: a.md` → `- [x]: x/a.md` | Befund | Befund |
  | 13 | `[x]: a.md (Nicht)` → `[x]: x/a.md (Immer)` | Befund | Befund |
  | 14 | `[x]: a.md [y]: b.md` → beide Ziele nachgezogen | Befund | Befund |
  | 15 | `[Regel]: nicht.erlaubt` → `[Regel]: immer.erlaubt` | still | Befund |
  | 16 | `> [x]: a.md` → `> [x]: x/a.md` | Befund | Befund |
  | 17 | `[A](a.md)[B](b.md)` → beide nachgezogen | still | Befund |
  | 18 | `[A](a.md)` → `[A](a.md "Nicht")` | Befund | Befund |
  | R1-0001 | `[^1]: Nicht gilt X.` → `[^1]: Immer gilt X.` | Befund | — |
  | R1-0002 | `[Hinweis]: Verboten ist X.` → `[Hinweis]: Erlaubt ist X.` | Befund | — |
  | R1-0003 | `Feld f](alt bleibt.` → `Feld f](neu bleibt.` | Befund | — |
  | R1-esc | `Regel \[Ausnahme\](keine) gilt.` → `…(alle) gilt.` | **still** | — |

  Lesart: 3, 8, 9, 11, 15 sind nach CommonMark Link-Ziele bzw. Referenz-Definitionen (die Zeile
  wird nicht gerendert) — gedeckt durch „das Gate sieht nur die Form" und die benannte
  `[Wort]: wort.md`-Grenze. 6 und 17 sind reine Nachzüge. 1, 2 und R1-esc sind **kein** Link nach
  Markdown und werden trotzdem geleert.

## Findings

### M-1 — Link-förmiger Text, der kein Link ist, wird geleert: Escape-Klammern ungenannt, Code nur als „Link in Code" genannt

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 1 und 18 · `DC-FA-VCS-001` · `AGENTS.md` §5 Regel 13
- **Pfad:** `internal/hexagon/core/rules/vcs.go` · „``inlineLinkTargetRE = regexp.MustCompile(`(!?\[[^\[\]]*\])\(``";
  `spec/lastenheft.md` · „Geleert wird nur, was die Link-Syntax trägt"; `spec/spezifikation.md` ·
  „eine Zeile ohne vollständigen Link wird nicht verändert"; `harness/sensors/adr-check.md` ·
  „Ein Link in Inline-Code oder einem Codeblock wird ebenso geleert"
- **Befund:** Das Muster prüft nicht, ob die öffnende Klammer escaped ist, und kennt keinen
  Code-Kontext: `\[Ausnahme\](keine)` (nach Markdown sichtbarer Fließtext, kein Link) und
  `` `fns[k](true)` `` (Code mit Index und Aufruf, kein Link) verlieren das Wort in der Klammer —
  Proben 1, 2, R1-esc passieren mit dem Schlüssel still eine inhaltliche Umkehr. Lastenheft und
  Spezifikation sagen das Gegenteil zu („nur, was die Link-Syntax trägt", „eine Zeile ohne
  vollständigen Link wird nicht verändert"); die Grenzen-Listen nennen den Code-Fall nur als
  „**Link** in Inline-Code", nicht als link-förmigen Code, und den Escape-Fall gar nicht. Kein Test
  hält eine der beiden Formen.
- **Warum nicht HIGH:** Die Code-Hälfte liegt in der Substanz der benannten Grenze („ohne
  Code-Kontext"), nur ihr Wortlaut ist zu eng; die Escape-Hälfte verlangt `\]` unmittelbar vor
  `(` und kommt im ADR-Bestand nicht vor (`grep -nE '\\\]\(' docs/plan/adr/*.md` leer). Der Pfad
  ist real, aber schmal; die Zusage im Vertrag ist für beide Formen falsch.
- **Verifizierbar:** ja — Proben 1, 2, R1-esc am Image; ein Kern-Test mit `\[A\](x)` liefe heute
  ohne Befund.
- **Klasse:** `muster-leert-mehr-als-gegenstand`

### M-2 — Die neue Geschichte-Zeile verwirft den Verfall-Vermerk mit einem Grund, der für ihn nicht gilt; die Kontext-Lesart aus R1 M-2 bleibt unkorrigiert

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 20 und 9 · `ADR-0103` §Geschichte · Baseline `modul-04-adrs.md` §Nachzug ist
  keine Überschreibung
- **Pfad:** `docs/plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md` · „trägt für künftige ADRs,
  nicht für den Bestand, dessen Verweise schon Adressen sind"
- **Befund:** Die Baseline nennt zwei Formen: Verweis über eine Kennung, und — **wo das nicht
  möglich ist**, also gerade bei einem Adress-Verweis — einen Verfall-Vermerk. „Die Verweise sind
  schon Adressen" trägt deshalb nur die Kennungs-Hälfte; für den Vermerk ist die Adresse die
  Voraussetzung, nicht der Ausschluss (er passt hier nicht, weil das Ziel umzieht statt zu
  verschwinden und ein toter Link das Link-Gate trifft — das steht nirgends). Die Entscheidung
  bliebe bei richtiger Begründung dieselbe. Zugleich bleibt der Kontext-Satz „darf die ADR ihm
  folgen", den R1 M-2 als über den Abschnitt hinaus gelesen meldete, ohne Einordnung; die
  Geschichte-Zeile ergänzt die Alternative, korrigiert die Lesart aber nicht. Failure-Szenario: ein
  späterer Leser von `ADR-0103` schließt für eine neue ADR mit Adress-Verweis, der Verfall-Vermerk
  komme nicht in Frage — das Gegenteil der Baseline; die Zeile ist immutabel.
- **Verifizierbar:** nein — Urteil über die Reichweite einer Begründung.
- **Klasse:** `begruendung-trifft-gegenstand-nicht`

### L-1 — Lastenheft-Änderung ohne Bump und Historie-Zeile

- **Kategorie:** LOW
- **Quelle:** `MR-032`
- **Pfad:** `spec/lastenheft.md` · „Geleert wird nur, was die Link-Syntax trägt — eine"; Historie ·
  „vor dem Vergleich des Core wird jedes Link-Ziel geleert"
- **Befund:** Der Fix-Commit ändert den Anforderungstext von `DC-FA-VCS-001` (neue Grenze), die
  Version bleibt 0.102.0 und die Historie-Zeile 0.102.0 sagt weiter „jedes Link-Ziel". `MR-032`
  verlangt Bump und Historie für **jede** Änderung, solange der Status unter `Accepted` liegt
  (Status `Draft`). Dieselbe Lage in der Spezifikation: Schritt 4 geändert, die §8-Historie-Zeile
  nennt weiter „jedes Link-Ziel … Grenze: zeilenweise, ohne Code-Kontext"
  (dort ohne `MR`-Anker, deshalb nur mitgenannt).
- **Verifizierbar:** nein — kein Gate hält Version und Historie gegen den Text.
- **Klasse:** `historie-ohne-aenderungszeile`

### L-2 — Kommentar am Konfig-Modell trägt die alte Zusage

- **Kategorie:** LOW
- **Quelle:** `AGENTS.md` §5 Regel 13 · §3.7 (Zusage muss stimmen)
- **Pfad:** `internal/hexagon/core/model/config.go` · „IgnoreLinkTargets bringt beim Vergleich des
  Core jedes Link-Ziel auf eine"
- **Befund:** R1 M-1 nannte diese Stelle unter den fünf Trägern der Zusage; der Fix hat sie nicht
  berührt. „Jedes Link-Ziel" ist nach dem Fix zu weit (Referenz-Ziele ohne Pfadzeichen, Ziel mit
  eigener Klammer bleiben ungeleert), „jede Änderung … an der übrigen Zeile bleibt eine" zu eng
  (M-1). Fail-safe in der ersten Hälfte, deshalb LOW.
- **Verifizierbar:** nein.
- **Klasse:** `grenzen-liste-ohne-groesste-luecke`

### L-3 — Das Beispiel `[Wort]: wort.md` ist enger als die Regel, die es illustriert

- **Kategorie:** LOW
- **Quelle:** `AGENTS.md` §5 Regel 13
- **Pfad:** `harness/sensors/adr-check.md` · „eine Zeile `[Wort]: wort.md` gilt als
  Referenz-Definition"; gleiche Form in der Geschichte-Zeile von `ADR-0103`
- **Befund:** Geleert wird jedes Ziel mit `.`, `/`, `#` oder `:` — auch ein Satzwort mit Punkt:
  Probe 3 (`[Status]: Abgelehnt.` → `Angenommen.`) und 15 passieren still. Die Spezifikation nennt
  die Zeichenklasse richtig; die Sensor-Doku und die ADR zeigen nur die Pfad-Form, ein Leser liest
  „`.md`-artig". Markdown rendert die Zeile nicht, deshalb LOW.
- **Verifizierbar:** ja — Probe 3.
- **Klasse:** `grenzen-beispiel-enger-als-regel`

### I-1 — Eine Verglichene Alternative steht jetzt in `## Geschichte`

- **Kategorie:** INFO
- **Quelle:** `AGENTS.md` §3.5 · Auftraggeber-Entscheid (Anhang statt Folge-ADR)
- **Pfad:** `docs/plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md` · „Alternative nachgetragen"
- **Befund:** Formal zulässig — §3.5 und `make adr-check` erlauben Geschichte-Anhänge ohne
  Inhaltsbeschränkung, die Zeile steht im ausgenommenen Abschnitt. Wer die Tabelle
  §Verglichene Alternativen liest, findet die Kennungs-/Vermerk-Alternative dort nicht; die
  Geschichte trägt damit Begründung, nicht nur Ereignis. Bewusste Entscheidung, deshalb INFO.
- **Verifizierbar:** nein.
- **Klasse:** `begruendung-im-geschichte-anhang`

### I-2 — gofmt-Ausrichtung aus R1 I-1 unverändert

- **Kategorie:** INFO
- **Quelle:** Maintainability (kein Konventions-Anker)
- **Pfad:** `internal/adapter/driven/configyaml/configyaml.go` · „``IgnoreLinkTargets bool    `yaml:"ignore-link-targets"` ``"
- **Befund:** Wie in R1; nicht adressiert, kein Gate sieht es.
- **Verifizierbar:** nein.
- **Klasse:** `formatierung-ohne-gate`

## Einlösung der R1-Findings

| R1 | Stand | Beleg |
|---|---|---|
| H-1 | eingelöst | Proben R1-0001/0002/0003 jetzt Befund; Tests „fussnote geaendert", „prosa in referenz-form geaendert", „klammer ohne link geaendert", „offener link ohne schliessende klammer". Restpfad für link-förmigen Nicht-Link: R2 M-1 |
| M-1 | weitgehend eingelöst | Spezifikation Schritt 4 + §2, Lastenheft, `adr-check.md` Grenze 4, Code-Kommentar präzisiert; `config.go` nicht (R2 L-2); Escape-/Code-Form fehlt (R2 M-1) |
| M-2 | teilweise | Alternative als Geschichte-Zeile nachgetragen; Begründung trägt nur die Kennungs-Hälfte, Kontext-Lesart unkommentiert (R2 M-2) |
| M-3 | eingelöst | Referenz-Seite getestet: Titel nachgezogen/geändert, Fußnote, Prosa, Einzelwort, Fußnote mit Pfadwort |
| L-1 | eingelöst | Spitzklammer mit Leerraum jetzt geleert (Test); Folgezeile und eigene Klammer als fail-safe benannt und getestet |
| L-2 | eingelöst | „Ohne den Schlüssel bleibt jeder Nachzug Drift." |
| I-1 | offen | R2 I-2 |
| I-2 | entfallen | betraf den Zeitpunkt des `Accepted`; der Fix ging über `## Geschichte` |

## Negativbefunde (geprüft, ohne Befund)

- **Linktext, übrige Zeile, Titel, Zahl der Links:** `$1`/`$2` erhalten Linktext und Titel samt
  Leerraum; Proben 4, 5, 10, 13, 18 Befund. Kein Befund.
- **Tabellenzellen, mehrere Links, angrenzende Links:** Proben 4, 10, 17 wie zugesagt. Kein Befund.
- **Einrückung/Container:** 0–3 Leerzeichen geleert (6), 4 Leerzeichen, Listen- und Zitat-Präfix
  bleiben Drift (7, 12, 16) — fail-safe. Kein Befund.
- **Titel in runden Klammern:** bei Inline-Link und Referenz-Definition Teil des Core (5, 13).
  Kein Befund.
- **Ziel mit Wörtern in Spitzklammern / Ziel mit Anführungszeichen (9, 11), Wort-Ziel (8):** nach
  CommonMark Link-Ziele, nicht gerendert; gedeckt durch „das Gate sieht nur die Form"
  (`ADR-0103` §Konsequenzen, `adr-check.md` Grenze 4). Kein Befund.
- **Reihenfolge der Ersetzungen:** die Inline-Leerung kann keine Zeile in eine Referenz-Definition
  mit Pfadzeichen verwandeln (geleertes `[a]()` trägt keines). Kein Befund.
- **Opt-in/byte-identisch:** ohne Schlüssel meldet jede der 22 Proben. Kein Befund.
- **Fail-safe-Grenzen gegen Code:** „Linktext mit eckiger Klammer", „Ziel mit eigener Klammer",
  „Referenz-Ziel ohne Pfadzeichen" stimmen mit den Mustern überein (Tests vorhanden). Kein Befund.
- **Kommentare §3.7:** neue Kommentare in `vcs.go` und im Test tragen Zusage bzw. Grenze, keine
  Herkunfts- oder Review-Prosa. Kein Befund (Wahrheitsfrage: M-1, L-2).
- **`ADR-0103` Geschichte-Zeile, Wahrheitsgehalt der Präzisierung:** „geleert wird nur ein
  vollständiger Link; was die Erkennung nicht trifft, bleibt Drift; `[Wort]: wort.md` gilt als
  Referenz-Definition" stimmt mit dem Code, bis auf die Escape-/Code-Form (M-1). Link auf
  `DC-FA-VCS-001.a` löst auf. Kein weiterer Befund.
- **`ADR-0005`/Hexagon, `DC-QA-03`, Suppressions:** nur `regexp`/`strings`, kein Netz, kein
  `//nolint`. Kein Befund.
- **Plan-Abgrenzung:** kein Umzug von `releasing.md`, kein Template-Feld-Nachzug. Kein Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 (M-1, M-2) |
| LOW | 3 (L-1, L-2, L-3) |
| INFO | 2 (I-1, I-2) |

## Verdikt

**Nicht blockiert durch HIGH; MEDIUM vor Merge zu klären.** H-1 aus R1 ist am Image eingelöst:
Fußnote, Prosa in Referenz-Form und `](` ohne Link bleiben Drift. Ein schmaler Restpfad bleibt —
link-förmiger Text, der nach Markdown kein Link ist (Escape-Klammern, Code wie `f[k](x)`), wird
geleert, während Lastenheft und Spezifikation das Gegenteil zusagen (M-1). M-2 betrifft die
Begründung in der immutablen ADR, nicht die Entscheidung.
