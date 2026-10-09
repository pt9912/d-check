# Review R1: slice-267, `vcs` lässt einen reinen Pfad-Nachzug durch

- **Review-Art:** Code. Geprüft wird der Diff `039eaafa..b64772a2` (Plan-Änderung `39f91a74`,
  `feat(vcs)` `64cbcfa7`, `docs(harness)` `b64772a2`) gegen den Slice-Plan `slice-267` samt
  Plan-Änderung, gegen die Entscheidungen (`ADR-0103` neu, `ADR-0024`, `ADR-0016`, `MR-025`) und
  gegen die Hard Rules `AGENTS.md` §3.4/§3.5/§3.6/§3.7 sowie §5 Regel 13. Die DoD-Abhakung prüft
  dieses Review nicht.
- **Gegenstand:** `slice-267` · `039eaafa..b64772a2`.
- **Skill:** `reviewer.md` @ 1.19.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-09
- **Eingangs-Kontext:** `DC-FA-VCS-001` im Lastenheft 0.102.0 (Absatz „Pfad-Nachzug (opt-in)",
  Kriterium „Boundary (Pfad-Nachzug)", Historie); Spezifikation `DC-FA-VCS-001.a` Schritt 4,
  §2-Zeile `vcs.ignore-link-targets`, §8-Historie; `ADR-0103` (Kontext, Entscheidung,
  Alternativen, Konsequenzen, Fitness Function); Baseline `v6.17.0` · `regelwerk/modul-04-adrs.md`
  §Nachzug ist keine Überschreibung. Im Code: `inlineLinkTargetRE`, `refDefTargetRE`,
  `blankLinkTargets`, `vcsCore` in `internal/hexagon/core/rules/vcs.go`; `model.VCSConfig`;
  `rawVCS`/`applyVCS` im YAML-Adapter; `config_template.go`; `vcs_pfad_nachzug_test.go`;
  `TestDecode_VCSIgnoreLinkTargets`; `.d-check.yml` (`vcs`-Block); `harness/sensors/adr-check.md`.
  Vorherige Findings am Modul: R1 zu `slice-247` (`vcs`, leere Range).
- **Proben** (Image `dcr267` per `make build IMAGE=dcr267` vom Stand `b64772a2`; Wegwerf-git-Repo im
  Scratchpad mit neun `Accepted`-ADRs, je eine Zeile zwischen zwei Commits geändert, Lauf
  `--enable vcs --disable links --range HEAD~1..HEAD`; danach Repo und Image entfernt, Arbeitsbaum
  unverändert):

  | ADR | BASE → HEAD | mit Schlüssel | ohne |
  |---|---|---|---|
  | 0001 | `[^1]: Nicht gilt X.` → `[^1]: Immer gilt X.` | **still** | Befund |
  | 0002 | `[Hinweis]: Verboten ist X.` → `[Hinweis]: Erlaubt ist X.` | **still** | Befund |
  | 0003 | `Feld f](alt bleibt.` → `Feld f](neu bleibt.` | **still** | Befund |
  | 0004 | `[A](<a b.md>)` → `[A](<x/a b.md>)` | still | Befund |
  | 0005 | `[A](<a b.md>)` → `[A](<a c.md>)` | Befund | Befund |
  | 0006 | `[A](` + Umbruch + `a.md)` → … `x/a.md)` | Befund | Befund |
  | 0007 | `[A](a.md)` → `[A](x/a.md)` | still | Befund |
  | 0008 | `[A](a.md) und [B](b.md)` → Ziele vertauscht | still | Befund |
  | 0009 | Codeblock `x = a](alt` → `x = a](neu` | still | Befund |

  Mit dem Schlüssel: 2 Befunde (0005, 0006), Exit 1. Ohne: 9 Befunde.

## Findings

### H-1 — `refDefTargetRE` leert das erste Wort jeder Zeile der Form `[…]:`, `inlineLinkTargetRE` das Wort hinter jedem `](`

- **Kategorie:** HIGH
- **Quelle:** Prüffrage 1 und 2 · `DC-FA-VCS-001` („Jede Änderung am Linktext, an der übrigen Zeile
  … bleibt Drift") · `ADR-0103` Entscheidung 3
- **Pfad:** `internal/hexagon/core/rules/vcs.go` · „`refDefTargetRE     = regexp.MustCompile(`^(\s{0,3}\[[^\]]+\]:)\s*\S+`)`"
  und „`inlineLinkTargetRE = regexp.MustCompile(`\]\([^)\s]*`)`"
- **Befund:** Die Muster sind nicht an die Link-Syntax gebunden: `\[[^\]]+\]:` trifft auch eine
  Fußnoten-Definition (`[^1]: …`) und jede Prosa-Zeile, die mit `[Wort]:` beginnt, und leert ihr
  erstes Wort; `\]\(` trifft jedes `](`, auch ohne Link, und leert das angehängte Wort. Proben 0001,
  0002, 0003: eine inhaltliche Umkehr („Nicht" → „Immer", „Verboten" → „Erlaubt") in einer
  `Accepted`-ADR passiert `make adr-check` mit dem in diesem Repo eingeschalteten Schlüssel grün,
  ohne den Schlüssel ist jede ein Befund. Im heutigen ADR-Bestand steht keine solche Zeile
  (`grep -nE '^\s{0,3}\[[^]]+\]:' docs/plan/adr/*.md` leer); für jede künftige ADR und jeden
  Konsumenten des opt-in-Schlüssels ist der Pfad offen.
- **Verifizierbar:** ja — Probe oben (Image-Lauf, 0001–0003 still mit, gemeldet ohne Schlüssel); ein
  Kern-Test mit einer Fußnoten-Zeile liefe heute grün, wo er rot sein müsste.
- **Klasse:** `muster-leert-mehr-als-gegenstand`

### M-1 — Die Zusage „die übrige Zeile bleibt Drift" steht in fünf Artefakten, die Grenze nennt nur den Code-Kontext

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 18 · `AGENTS.md` §5 Regel 13
- **Pfad:** `docs/plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md` · „Weiterhin ein Befund bleibt
  jede Änderung am Linktext, an der übrigen Zeile"; `spec/lastenheft.md` · „Jede Änderung am
  Linktext, an der übrigen Zeile, am Titel eines Links oder an"; `spec/spezifikation.md` ·
  „eine Änderung am Linktext, an der übrigen Zeile, am Titel oder an der Zahl der Links bleibt Drift"
  (Schritt 4) und die §2-Zeile `vcs.ignore-link-targets`; `internal/hexagon/core/model/config.go` ·
  „jede Aenderung"; `harness/sensors/adr-check.md` · „Ein Pfad-Nachzug wird nur an seiner Form erkannt."
- **Befund:** Alle fünf sagen zu, dass außerhalb des Link-Ziels jede Änderung Drift bleibt, und
  nennen als einzige Mechanik-Grenze „Link-Ziel in Inline-Code oder Codeblock". Der Code leert mehr
  (H-1): das erste Wort jeder `[…]:`-Zeile und das Wort hinter jedem `](` auch außerhalb eines Links
  — die größte Lücke des Schlüssels fehlt in jeder Grenzen-Liste. `ADR-0103` ist `Accepted` und
  friert die falsche Zusage ein.
- **Verifizierbar:** ja — Proben 0001–0003 gegen den Wortlaut; kein Gate fängt die Zusage.
- **Klasse:** `grenzen-liste-ohne-groesste-luecke`

### M-2 — `ADR-0103` liest aus der Baseline ein „darf ihm folgen", das der Abschnitt nicht sagt, und lässt dessen eigene Form aus

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 9 · `ADR-0103` §Kontext und §Verglichene Alternativen
- **Pfad:** `docs/plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md` · „Verschiebt sich das Ziel
  eines Verweises, darf die ADR ihm folgen"; Gegenstelle `v6.17.0` · `regelwerk/modul-04-adrs.md`
  §Nachzug ist keine Überschreibung · „trägt die ADR einen deklarierten Vermerk: Der Verweis *verfällt*"
- **Befund:** Der Baseline-Abschnitt nimmt den Nachzug von der Hard Rule aus, beschreibt als Form aber
  den Verweis über eine **Kennung** und, wo das nicht geht, einen **Verfall-Vermerk**; dass die ADR
  dem verschobenen Ziel durch Editieren folgt, steht dort nicht ausdrücklich (gedeckt ist nur „nur
  ihre Adresse"). Die Alternativen-Tabelle vergleicht `ignore-refs` und `--no-verify`, nicht die
  Form, die die zitierte Stelle selbst nennt. Die ADR ist immutabel, und ihr erster
  Re-Evaluierungs-Trigger („Das Baseline-Regelwerk ändert die Ausnahme") misst künftig gegen eine
  Lesart, die der Wortlaut nicht trägt. Die Entscheidung selbst trägt `AGENTS.md` §3.6 auch ohne
  diese Begründung (Prüffrage-20-Test).
- **Verifizierbar:** nein — Urteil über den Geltungsbereich eines Zitats.
- **Klasse:** `quelle-ueber-geltungsbereich`

### M-3 — Kein Negativtest auf der Referenz-Definitions-Seite außer dem Label

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 13 · `DC-FA-VCS-001` Kriterium „Boundary (Pfad-Nachzug)"
- **Pfad:** `internal/hexagon/core/rules/vcs_pfad_nachzug_test.go` · „`{"referenz-label geaendert"`"
- **Befund:** Die Negativfälle (Linktext, Text neben dem Link, Titel, Zahl der Links) sind alle
  Inline-Links; für `refDefTargetRE` gibt es nur den Label-Fall. Weder eine Fußnoten- oder
  `[Wort]:`-Prosazeile (H-1) noch ein Titel hinter einer Referenz-Definition noch ein `](` ohne Link
  ist getestet — genau die Klasse, an der das Muster mehr leert als zugesagt, läuft durch die Suite
  ungeprüft.
- **Verifizierbar:** ja — `make test` bleibt mit dem fehlerhaften Muster grün.
- **Klasse:** `negativtest-fehlt-oeffentlicher-vertrag`

### L-1 — Die Gegenrichtung ist ungenannt: zwei Nachzug-Formen bleiben Drift

- **Kategorie:** LOW
- **Quelle:** `AGENTS.md` §5 Regel 13
- **Pfad:** `spec/spezifikation.md` · „Ein **reiner\n   Pfad-Nachzug** ergibt so denselben Core"
  (Schritt 4); `vcs.go` · „GRENZE: zeilenweise und ohne Code-Kontext"
- **Befund:** Probe 0005 (Spitzklammer-Ziel mit Leerraum, Änderung hinter dem Leerraum) und 0006
  (Ziel auf der Folgezeile hinter `[A](`, eine Form, die das Modul `links` laut `ADR-0091` liest)
  sind reine Nachzüge und bleiben Befund. Die Grenze „zeilenweise" deckt 0006 nur implizit, 0005 gar
  nicht; im Bestand der ADRs kommt keine der beiden Formen vor (`grep -nE '\]\($'` leer). Fail-safe,
  deshalb LOW.
- **Verifizierbar:** ja — Proben 0005/0006.
- **Klasse:** `grenzen-liste-ohne-groesste-luecke`

### L-2 — Testkommentar trägt einen zeitbezogenen Rest

- **Kategorie:** LOW
- **Quelle:** `AGENTS.md` §3.7
- **Pfad:** `internal/hexagon/core/rules/vcs_pfad_nachzug_test.go` · „Ohne den Schlüssel meldet derselbe
  Nachzug — der Vergleich davor."
- **Befund:** Der Kommentar trägt eine Zusage; „der Vergleich davor" beschreibt das Verhalten vor
  der Änderung — Herkunft, keine der fünf Klassen, und altert mit dem nächsten Umbau des Vergleichs.
  Ein Fragment in einem sonst klassentragenden Kommentar, deshalb nicht HIGH.
- **Verifizierbar:** nein — kein Gate prüft Kommentar-Klassen.
- **Klasse:** `kommentar-herkunfts-prosa`

### I-1 — `rawVCS` und die Test-Struktur sind nicht gofmt-ausgerichtet

- **Kategorie:** INFO
- **Quelle:** Maintainability (kein Konventions-Anker: das Lint-Profil führt keinen Formatierer)
- **Pfad:** `internal/adapter/driven/configyaml/configyaml.go` · „``IgnoreLinkTargets bool    `yaml:"ignore-link-targets"` ``";
  `vcs_pfad_nachzug_test.go` · „`mitSchalter, ohne int`"
- **Befund:** Die übrigen Felder sind auf die alte Spaltenbreite ausgerichtet; `gofmt` würde den Block
  neu ausrichten. Kein Gate sieht das.
- **Verifizierbar:** nein.
- **Klasse:** `formatierung-ohne-gate`

### I-2 — `ADR-0103` ist im Commit ihres Entstehens `Accepted`; ihre Fitness Function nennt einen künftigen Beleg

- **Kategorie:** INFO
- **Quelle:** `AGENTS.md` §3.5 · `ADR-0103` §Fitness Function
- **Pfad:** `docs/plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md` · „ist der Beleg
  am Bestand"
- **Befund:** Eine Korrektur von Entscheidung 3 oder §Konsequenzen (M-1, M-2) trifft den lokalen
  `pre-commit`-Hook (`STAGED=1`: BASE = `HEAD`, dort `Accepted`); die CI-Range ab dem Push-Stand
  sieht die Datei ohne BASE. Der Satz zur Fitness Function benennt einen Lauf (slice-268), der noch
  nicht stattfand, als Beleg in einem immutablen Dokument.
- **Verifizierbar:** nein.
- **Klasse:** `adr-accepted-vor-review`

## Negativbefunde (geprüft, ohne Befund)

- **Opt-in, byte-identisch ohne Schlüssel:** `blankLinkTargets` läuft nur unter
  `if ignoreLinkTargets`; Default `false` im Adapter; Probe ohne Schlüssel meldet alle neun Fälle.
  Kein Befund.
- **Inline-Titel bleibt Teil des Core:** `[^)\s]*` stoppt am Leerraum vor dem Titel; getestet
  („titel geaendert"). Kein Befund.
- **Mehrere Links je Zeile, Bild, Anker im Ziel:** `ReplaceAllString` über die Zeile; getestet. Kein Befund.
- **Ziel mit eigener Klammer:** geleert bis zur ersten `)`; der Rest bleibt verglichen — fail-safe, so
  benannt in Spezifikation und Code. Kein Befund.
- **Vertauschte Ziele / anderes Dokument:** still (Probe 0008), in `ADR-0103` und
  `harness/sensors/adr-check.md` als Grenze benannt. Kein Befund.
- **`ADR-0005`/Hexagon:** `rules` nutzt nur `regexp`/`strings`; Modell-Feld im Kern, Dekodierung im
  Adapter. Kein Befund.
- **`DC-QA-03`:** kein neuer Lesepfad, kein Netz. Kein Befund.
- **`MR-025` Spiegel:** Modell, Adapter, `--print-config`-Vorlage, Spezifikation (Schritt 4 + §2),
  Lastenheft, `.d-check.yml`, `AGENTS.md` §3.5, `harness/README.md` §Traceability rules,
  `harness/sensors/adr-check.md` nachgezogen; `grep head-allow` findet sonst nur Handbuch und
  CHANGELOG (Release-Prep) sowie eingefrorene Artefakte. Kein Befund.
- **`AGENTS.md` §3.6:** die Lockerung trägt eine ADR; die Plan-Änderung `39f91a74` liegt vor dem
  Code-Commit. Kein Befund.
- **Plan-Abgrenzung:** kein Umzug von `releasing.md`, kein Template-Feld-Nachzug im Diff. Kein Befund.
- **`AGENTS.md` §3.4:** die neuen Lastenheft-/Spezifikations-Zeilen nennen weder ADR noch Slice. Kein Befund.
- **`MR-032`:** Bump 0.101.3 → 0.102.0 samt Historie-Zeile. Kein Befund.
- **Prüffrage 16:** `Schärft:` nennt die Kennung `DC-FA-VCS-001`. Kein Befund.
- **Prüffrage 14:** die Provenance-Marker in `ADR-0103` zeigen Entstehung bzw. Prüfort, begründen nichts
  (der künftige Charakter des zweiten steht in I-2). Kein Befund.
- **ADR-Index, Re-Evaluierungs-Trigger:** vorhanden. Kein Befund.
- **Kommentare `vcs.go`, `config_template.go`, `.d-check.yml`:** Zusage/Grenze mit auflösbarem Feld;
  die inhaltliche Unwahrheit der Zusage ist M-1. Kein weiterer Befund.
- **Suppressions:** keine `//nolint`. Kein Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 (H-1) |
| MEDIUM | 3 (M-1, M-2, M-3) |
| LOW | 2 (L-1, L-2) |
| INFO | 2 (I-1, I-2) |

## Verdikt

**Blockiert.** H-1 ist ein Stilles-Grün-Pfad im ADR-Gate dieses Repos: mit dem eingeschalteten
Schlüssel passiert eine inhaltliche Umkehr in einer Fußnote oder einer `[Wort]:`-Zeile `make
adr-check`. M-1 hängt an H-1 (die Zusage in fünf Artefakten, eines davon eine `Accepted`-ADR), M-3 ist
die fehlende Testseite dazu; M-2 betrifft die Begründung der ADR, nicht die Entscheidung. Die
Inline-Link-Hälfte des Schlüssels trägt, was Plan und ADR für sie zusagen.
