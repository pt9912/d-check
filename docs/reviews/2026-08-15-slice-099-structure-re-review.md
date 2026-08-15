# Re-Review-Report: slice-099 (Modul `structure`) — 2026-08-15

**Review-Art:** **Re-Review** (Code) — bestätigend geprüft wird, ob die
Befunde des ersten Laufs
([Erst-Report](2026-08-15-slice-099-structure-review.md): 1 HIGH, 7 MEDIUM,
3 LOW, 2 INFO) echt geheilt sind und ob die Heilung Neues aufgemacht hat.
Gemessen wird gegen Lastenheft, Spezifikation,
[ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md), den
Slice-Plan und
[`MR-025`](../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten).
Die DoD-Abhakung ist ausdrücklich **nicht** Gegenstand.

**Gegenstand:** `59a73a2..HEAD` (`8258a56`, der Heilungs-Commit) im Kontext des
Gesamt-Slice `64c62cb..HEAD`.

**Skill:** `.harness/skills/reviewer.md` @ 1.3.0 ·
**Modell:** `claude-opus-5[1m]` · **Datum:** 2026-08-15

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [slice-099](../plan/planning/in-progress/slice-099-structure-modul.md),
  besonders §3a (Spiegel-Liste) und **§3b** (die Bilanz),
  [welle-73](../plan/planning/welle-73-structure-umsetzung.md)
- [`MR-025`](../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten)
  in der **nach** dem Erst-Review geschärften Fassung
- [`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
  und [`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  (acht Schritte), §2-Schema, §4-Grund-Code-Tabelle
- [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  und [`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritte C3–C5
- [`DC-FA-CLI-002`](../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl),
  [`DC-FA-CLI-005`](../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben),
  [`DC-FA-CLI-006`](../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten),
  [`DC-FA-CLI-010`](../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
- [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md) (Accepted,
  Fitness Function), [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md)
- [`AGENTS.md`](../../AGENTS.md) Hard Rules

**Eigene Läufe dieses Re-Reviews:** `make build`, `make test` (Baseline plus
**acht** Mutationsläufe), `make lint` und `make gates` (grün, Coverage 94,20 %),
`make verify-closure-notes` (338 Dateien, 0 Befunde), `make fullbuild` (grün,
Image-Hash `sha256:5dc3e484…`); dazu ein aus dem **Vor-Heilungs**-Commit
`59a73a2` gebautes Vergleichs-Image (`make build IMAGE=d-check-old` in einem
`git archive`-Auszug) und Image-Läufe gegen acht Fixture-Bäume in einem
Temp-Verzeichnis außerhalb des Repos, alle mit `--network none` und
`:ro`-Mount ([`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
Der Arbeitsbaum ist unverändert (`git status --short` leer, Beleg am Ende).

---

## Erst-Befunde: Status

| Erst-Befund | R1 | Status | Beleg (am Lauf, nicht am Code) |
|---|---|---|---|
| F-1 — `closure-note-ambiguous` von keinem Codepfad erzeugt | HIGH | **geheilt** | Fixture mit zwei Überschriften (Z. 3 leer, Z. 7 substanziell): **genau ein** Befund, `closure-note-ambiguous` in **Z. 7**, kein `-thin`. Drei Treffer (Z. 3/7/11) ⇒ ein Befund in **Z. 7** (zweiter Treffer). Zwei Treffer **verschiedener Ebene** (`###` Z. 3, `##` Z. 7) ⇒ ambiguous Z. 7. Ein Treffer **im Fence** plus ein echter ⇒ **kein** Befund (gemessen). Fixture mit Floskel **und** Platzhalter im ersten Abschnitt ⇒ **nur** ambiguous. Mutationen M5/M6 rot |
| F-2 — drei Flächen sagen eine Verhaltensänderung zu, die nicht eintritt | MEDIUM | **geheilt**, mit Rest | Gegen das aus `59a73a2` gebaute Image sind die Befundsätze jetzt **verschieden** statt byte-identisch. Rest: die Richtung „vorher grün ⇒ jetzt rot" fehlt in allen drei Flächen ⇒ **N-3** |
| F-3 — Marken-Grenze prüft ASCII-Bytes | MEDIUM | **geheilt** | 13-Zeilen-Matrix gegen `require-all` mit der Marke `Beleg`: Umlaut, Kyrillisch und arabisch-indische Ziffer gelten **nicht** mehr als Grenze (Befund `section-marker-missing`); Bindestrich, Doppelpunkt, Punkt, Klammer bleiben Grenze. Mutation M2 rot |
| F-4 — Schritte 5/6 beschreiben Bereinigung und Zählung anders als der Code | MEDIUM | **teilweise** | Algorithmus-Schritte 5 und 6 sind nachgezogen (Inline-Code-Leerung samt Folge, Satzende nur vor Whitespace). **Nicht** nachgezogen: die §2-Schema-Zeile zu `structure[].min-sentences` und der Handbuch-§5-Block ⇒ **N-4** |
| F-5 — `DC-FA-CLI-010`-Akzeptanzkriterien führen `doc-structure` nicht | MEDIUM | **geheilt** | `--print-mk` liefert 12 `doc-*`-Targets; die Mengen aus Happy Path und Boundary sind jetzt **deckungsgleich** mit der Ausgabe (beide 12, gleiche Namen) |
| F-6 — 20. Modul fehlt in der Nicht-aktiviert-Enumeration | MEDIUM | **teilweise** | Die erzeugte Vorlage nennt `structure` jetzt, ebenso das Akzeptanzkriterium. Die normative Klausel §Aufnahme ins Modulset derselben Anforderung klassifiziert `structure` weiterhin **nicht** ⇒ **N-2** |
| F-7 — vier Rückbauten bleiben grün | MEDIUM | **geheilt** | Alle vier selbst nachgestellt und **rot**: Schwelle `<` auf `<=` (M1), Ziffernbereich der Wortgrenze entfernt (M2), `TrimSpace` im Selektor entfernt (M3), fail-closed-Zweig auf `return nil` (M4). Zwei **neue** Rückbauten bleiben grün ⇒ **N-6** |
| F-8 — Spiegel-Liste nennt Unberührtes und lässt Spiegel aus | MEDIUM | **geheilt**, mit Rest | §3b bilanziert vier Auslassungen und die falsch eingetragene Zeile; `MR-025` ist um die Ableitungs-Regel geschärft. Die Bilanz ist **selbst** noch unvollständig ⇒ **N-1** |
| F-9 — Zähl-Wächter über dem `--print-mk`-Template | LOW | **nicht geheilt** | `internal/adapter/driving/cli/print_mk.go:22-30` sagt weiter „elf `##`-annotierte Targets" und „genau SECHS fmt-Verben"; das Template führt 12 Targets, `makefileFragment` übergibt **sieben** Argumente |
| F-10 — `planning`-Zeile der Handbuch-Modultabelle ohne den neuen Code | LOW | **nicht geheilt**, jetzt schwerer | `docs/user/benutzerhandbuch.md:1660` listet weiterhin fünf Codes ohne `closure-note-ambiguous`, während der Code seit `8258a56` **wirklich** feuert (belegt) — §4-Tabelle (Z. 508) und §11-Zeile (Z. 1841) führen ihn |
| F-11 — doppelte Aufzählungs-Konjunktion in der deutschen README | LOW | **nicht geheilt** | `README.de.md:15-19` trägt unverändert „von … bis zu … bis zu"; `README.md:16-20` nicht |
| F-12 — Symlink-Kandidaten still übergangen | INFO | **nicht geheilt, nicht dokumentiert** | Fixture mit Datei-Symlink (`docs/link.md`) und Verzeichnis-Symlink unter `files: "docs/**"`: nur `docs/a.md` erscheint im Befundsatz, die beiden Symlink-Kandidaten ohne Befund und ohne Hinweis. Keine Vertragsfläche sagt es |
| F-13 — `CheckStructure` verwirft einen Walk-Fehler | INFO | **nicht geheilt, nicht dokumentiert** | `internal/hexagon/core/rules/structure.go:23-26` kehrt bei Fehler aus `structureTree` weiterhin mit `nil` zurück; keine Notiz in Slice, ADR oder Spezifikation |

## Neue Findings

### N-1 — Das Gerüst von `--print-config` nennt sein eigenes 20. Modul nicht in der Verfügbar-Liste

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-CLI-005`](../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben)
  §Beschreibung („es dokumentiert die verfügbaren Module und Optionen als
  Kommentare, damit sie sichtbar sind");
  [`MR-025`](../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten)
  Spiegel „Emittierte Vorlage";
  [slice-099](../plan/planning/in-progress/slice-099-structure-modul.md) §3a
  Zeile `--print-config`-Vorlage und §3b
- `pfad`: `internal/adapter/driving/cli/config_template.go:21`
- `befund`: Die Verfügbar-Zeile der erzeugten Vorlage nennt 18 Module und ist
  vor und nach dem Slice **byte-identisch** (`git show 64c62cb` gegen HEAD);
  `structure` fehlt, obwohl derselbe Slice in **derselben Datei** einen
  `structure`-Gerüstblock ergänzt hat, und `citations` fehlt seit älterer
  Herkunft. Die Bilanz §3b führt diesen Spiegel nicht. Der im Heilungs-Commit
  reparierte `--suggest-config`-Kommentar verweist für das Voll-Schema
  ausdrücklich auf `--print-config` — der Leser landet damit auf einer Liste,
  die das eben beworbene Modul nicht führt. Die Verriegelung greift nicht: der
  Akzeptanztest prüft die Verfügbar-Zeile per Teilstring
  (`internal/adapter/driving/cli/cli_acceptance_test.go:838`), und ein
  Teilstring bleibt gültig, wenn ein Eintrag **fehlt**.
- `verifizierbar`: ja — `--print-config` gegen die Modul-Liste aus
  [`DC-FA-CLI-002`](../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
  bzw. `docs/user/operations.md:23`: dort stehen 20 Module, in der
  Verfügbar-Zeile 18.
- `klasse`: „neues Modul nicht in der Verfügbar-Enumeration der eigenen
  emittierten Vorlage"

### N-2 — `DC-FA-CLI-006`: die als geschlossen deklarierte Modul-Klassifikation kennt `structure` nicht, das Akzeptanzkriterium schon

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-CLI-006`](../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)
  §Aufnahme ins Modulset (**geschlossene Menge**);
  [`MR-025`](../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten)
  Spiegel „Anforderung: Beschreibung **und** Akzeptanzkriterien"
- `pfad`: `spec/lastenheft.md:237-242` gegen `spec/lastenheft.md:297`
- `befund`: Die Klausel teilt jedes Modul in „fixe Aktiv-Menge" oder „alle
  übrigen" mit Grund; sie nennt sieben inaktive Module und begründet jedes über
  ein K-Kriterium. `structure` steht in keiner der beiden Aufzählungen, obwohl
  das Akzeptanzkriterium derselben Anforderung es seit dem Heilungs-Commit
  aufführt. Die beiden Hälften einer Anforderung widersprechen sich damit
  wieder — dieselbe Klasse wie der geheilte F-5, nur mit vertauschten Rollen
  (dort war die Beschreibung nachgezogen und das Kriterium nicht, hier
  umgekehrt). `citations` und `sources` fehlen in derselben Klausel aus älterer
  Herkunft.
- `verifizierbar`: ja — die Klausel gegen die Modul-Liste aus
  [`DC-FA-CLI-002`](../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl):
  17 klassifizierte gegen 20 existierende Module.
- `klasse`: „Zahl im Beschreibungstext nachgezogen, Enumeration in den
  Akzeptanzkriterien nicht" (invers)

### N-3 — Die Release-Notiz nennt den Code-Tausch, nicht die von der ADR verlangte Richtung „vorher grün, jetzt rot"

- `kategorie`: MEDIUM
- `quelle`: [ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md)
  §Konsequenzen („Ein Repo mit zwei Closure-Abschnitten wird danach **rot**, wo
  es vorher still den ersten las. Das gehört in die Release-Notiz.")
- `pfad`: `CHANGELOG.md:25-33`, `spec/lastenheft.md:2610`,
  `docs/user/benutzerhandbuch.md:1841`
- `befund`: Die Notiz beschreibt ausschließlich den Tausch eines bestehenden
  Befunds („meldete bisher `closure-note-thin`/`-boilerplate`/`-placeholder`
  über den ersten Abschnitt; jetzt … `closure-note-ambiguous`"). Gemessen gegen
  das aus `59a73a2` gebaute Vergleichs-Image tritt daneben der Fall auf, den die
  ADR benennt: ein Dokument mit zwei passenden Überschriften, dessen **erster**
  Abschnitt substanziell ist, war unter dem Vor-Stand **befundfrei** und meldet
  jetzt `closure-note-ambiguous`. Ein Konsument, dessen Closure-Lauf heute grün
  ist, liest aus allen drei Flächen keinen Anlass, ihn vor dem Update erneut zu
  fahren. Für die vergleichbare Änderung in 0.56.0 führt dieselbe Handbuch-Spalte
  die Richtungen ausdrücklich („ein grüner Lauf kann rot werden").
- `verifizierbar`: ja — Fixture mit zwei Überschriften und substanziellem
  erstem Abschnitt gegen beide Images: Vor-Stand 0 Befunde, HEAD 1 Befund.
- `klasse`: „Release-Notiz nennt nur eine Richtung einer zweiseitigen
  Verhaltensänderung"

### N-4 — §2-Schema und Handbuch-§5 tragen die in Schritt 5/6 nachgezogene Bereinigung und Zählung nicht

- `kategorie`: LOW
- `quelle`: [`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  Schritte 5–6 (neu) gegen §2-Schema; Erst-Befund F-4 nannte drei Flächen
- `pfad`: `spec/spezifikation.md:2336` und
  `docs/user/benutzerhandbuch.md:1608-1640`
- `befund`: Die Schema-Zeile zu `structure[].min-sentences` sagt weiterhin nur
  „Satzende-Zeichen außerhalb Fenced-Code" — weder die Inline-Code-Leerung noch
  die Bedingung „nur vor Whitespace oder Zeilenende" steht dort, obwohl Schritt 6
  sie seit dem Heilungs-Commit führt. Der Handbuch-Block zum Modul nennt unter
  „Drei Dinge, die überraschen können" beides ebenfalls nicht, während der
  Parallel-Abschnitt zum Preset-Partner es seit 0.56.0 tut. Wer die Schwelle nach
  dem Schema wählt, zählt für „a.b.c.d." vier und bekommt eins.
- `verifizierbar`: ja — Abschnitt mit Punktkette und Inline-Code-Marke,
  `--enable structure`, gegen die Schema-Zeile gelesen.
- `klasse`: „Preset-Partner-Spezifikationen driften an derselben geteilten
  Mechanik"

### N-5 — Zwei Wort-Begriffe in derselben geteilten Mechanik, nur einer steht im Vertrag

- `kategorie`: LOW
- `quelle`: [`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  Schritt 6 („mit einem nicht-alphanumerischen Zeichen weitergeht") gegen
  [`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritt C4 (Wortzeichen ausdrücklich als ASCII definiert); Beobachtung
  BEO-003 im [Register](../plan/planning/observations.md)
- `pfad`: `internal/hexagon/core/rules/structure.go:206-213` gegen
  `internal/hexagon/core/rules/planning.go:249-261`
- `befund`: Die Marken-Grenze ist seit dem Heilungs-Commit unicode-weit
  (`unicode.IsLetter` bzw. `unicode.IsDigit`), die Floskel-Wortgrenze bleibt
  ASCII einschließlich Unterstrich. Beobachtbar an derselben Eingabe: ein
  Umlaut setzt eine **Marke** fort, beendet aber eine **Floskel**; ein
  Unterstrich beendet eine Marke, setzt aber eine Floskel fort
  (`- **Beleg_2:**` erfüllt die Marke `Beleg`). Der Spezifikationstext zu
  Schritt 6 ist vor und nach dem Heilungs-Commit **identisch** — er trennt die
  beiden Implementierungen nicht, und die 0.56.1-Historie schreibt für die
  Nachbarbedingung ausdrücklich die entgegengesetzte Menge fest. Wer von dort
  auf `require-all` schließt, erwartet das Gegenteil des Gemessenen.
- `verifizierbar`: ja — Fixture-Matrix mit Umlaut und Unterstrich, jeweils
  gegen `require-all` und gegen `boilerplate`.
- `klasse`: „geteilte Lexik mit zwei Wort-Begriffen, nur einer im Vertrag"

### N-6 — Zwei weitere Rückbauten am neuen Modul bleiben grün, beide in Richtung Falsch-Grün

- `kategorie`: MEDIUM
- `quelle`: Reviewer-Anker „fehlende Negativtests bei neuem öffentlichen
  Vertrag";
  [`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  Schritt 6 (Task-Item-Definition und Marken-Definition)
- `pfad`: `internal/hexagon/core/rules/structure.go:172` und
  `internal/hexagon/core/rules/structure.go:195-198`, gegen
  `internal/hexagon/core/rules/structure_test.go` (gesamt)
- `befund`: (1) Streicht man in `taskItemRE` die **nummerierte** Listen-Form,
  zählt `max-tasks` Zeilen wie `1. [ ] offen` nicht mehr — `make test` bleibt
  grün; kein Fixture des Moduls trägt ein nummeriertes Task-Item, obwohl die
  Spezifikation den Ziffern-Marker ausdrücklich nennt und das Produkt ihn
  zählt (im Fixture gemessen: 2 Task-Items). (2) Ersetzt man in `hasMarker` den
  Präfix-Vergleich durch einen Teilstring-Vergleich, erfüllt
  `- **Zwischenbeleg:**` fortan die Marke `Beleg` — `make test` bleibt
  ebenfalls grün; die Zusage „dessen Inhalt mit `M` **anfängt**" ist von keinem
  Test bewacht, nur ihre Fortsetzungs-Hälfte. Beide Rückbauten liegen außerhalb
  der vom Erst-Review gemeldeten vier und wirken in Richtung stilles Grün.
- `verifizierbar`: ja — jede Mutation einzeln, `make test` grün (Läufe
  durchgeführt, Datei danach aus einer Kopie wiederhergestellt).
- `klasse`: „Schwellen- und fail-closed-Ränder eines neuen Moduls ohne Test"

### N-7 — `closureHeadingLine` trägt nach der Signatur-Änderung zwei gestapelte Doc-Kommentare

- `kategorie`: LOW
- `quelle`: Maintainability (latente Wartungsfalle in einer Funktion, deren
  Vertrag der Heilungs-Commit gerade erweitert hat)
- `pfad`: `internal/hexagon/core/rules/planning.go:169-182`
- `befund`: Über der Funktion stehen zwei mit dem Funktionsnamen beginnende
  Kommentarblöcke unmittelbar hintereinander: der ältere beschreibt die
  Rückgabe als „die 1-basierte Zeilennummer der **ersten** Überschrift … 0 ⇒
  kein Treffer" und kennt den dritten Rückgabewert nicht, der jüngere
  beschreibt ihn. Wer nur den ersten liest — er steht zuerst und nennt die
  Funktion beim Namen —, hält die Mehrdeutigkeits-Behandlung für nicht
  vorhanden; genau diese Fehleinschätzung war der Erst-Befund F-1.
- `verifizierbar`: nein — kein Gate deckt Doc-Kommentar-Dopplung; durch Lesen
  belegbar.
- `klasse`: „Doc-Kommentar nach Signatur-Änderung gestapelt statt ersetzt"

### N-8 — Der `one`-Abbruch in der Treffer-Schleife ist nach dem frühen Return wirkungslos

- `kategorie`: INFO
- `quelle`: [`DC-FA-STRUCT-001.a`](../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure)
  Schritt 4 (Kardinalität)
- `pfad`: `internal/hexagon/core/rules/structure.go:107-109`
- `befund`: Im Modus `one` kehrt die Funktion bei mehr als einem Treffer schon
  vorher zurück; die Schleife wird in diesem Modus daher nur mit **genau einem**
  Treffer betreten, und das `break` an ihrem Ende kann keinen Durchlauf sparen.
  Es ist damit ein Zweig, den keine Mutation von seiner Abwesenheit
  unterscheiden kann — er sieht wie eine tragende Verriegelung der
  Kardinalitäts-Zusage aus und ist keine.
- `verifizierbar`: nein — verhaltensgleich mit und ohne; nur durch Lesen bzw.
  eine Mutation belegbar, die grün bleibt.
- `klasse`: „wirkungsloser Zweig in Verriegelungs-Gestalt"

## Negativbefunde

- geprüft, ohne Befund: **Erreichbarkeit aller Grund-Codes (die Klasse hinter
  F-1).** Methode statt Stichprobe: aus `AllReasons()` die 44
  Konstanten-Namen gezogen, für jede die Deklarationsstelle bestimmt und alle
  Nicht-Test-Go-Dateien außer `internal/hexagon/core/app/diagnose.go` (dem
  Katalog selbst) nach Referenzen abgesucht. **Jede** der 44 Konstanten hat
  mindestens eine Referenz außerhalb von Deklaration und Katalog; 40 stehen
  direkt an einer Befund-Konstruktion, die vier übrigen (`external-status`,
  `external-timeout`, `external-redirects`, `source-unreachable`) an einer
  Klassifikator-Funktion, deren Rückgabewert in einen Befund fließt. Kein
  zweiter Code ist deklariert, dokumentiert und tot.
- geprüft, ohne Befund: **Erreichbarkeit der neun neuen Codes am Lauf.** Ein
  Fixture mit neun Regeln erzeugt **alle acht** `section-*`-Codes in einem Lauf;
  `closure-note-ambiguous` erzeugt das Closure-Fixture. Damit sind die neun
  Codes dieses Slice nicht nur statisch referenziert, sondern beobachtet.
- geprüft, ohne Befund: **Nebenwirkungen des frühen Returns.** Genau eine
  Überschrift wird gemessen (`closure-note-thin` an ihrer Zeile); keine
  Überschrift meldet `closure-note-missing`; ein fehlendes Closure-Verzeichnis
  und ein Verzeichnis ohne Kandidaten melden `closure-note-missing` mit dem
  Verzeichnis als Ort; eine vorhandene, aber unlesbare Datei meldet weiterhin
  fail-closed **vor** der Mehrdeutigkeits-Prüfung (Reihenfolge im Lauf
  bestätigt). Alle vier Fälle sind gegenüber dem Vor-Stand unverändert.
- geprüft, ohne Befund: **Warum `verify-closure-notes` auf dem eigenen Bestand
  grün bleibt.** Von 99 abgeschlossenen Slices trägt **genau einer**
  (`docs/plan/planning/done/slice-094-closure-zaehl-paritaet.md`) zwei Zeilen,
  die auf das `heading-pattern` passen — die erste liegt **innerhalb** eines
  Fenced-Blocks und ist damit keine Überschrift. Gegenprobe an einer Kopie
  außerhalb des Repos: entfernt man die beiden Fence-Zeilen, meldet derselbe
  Lauf `closure-note-ambiguous`. Die Grünheit hängt also an der Fence-Lexik und
  nicht am Zufall — und sie ist damit belegt statt behauptet.
- geprüft, ohne Befund: **`--doctor`-Erreichbarkeit des neuen Klartexts.** Ein
  echter `--doctor`-Lauf über das Fixture zeigt „Mehrere
  Closure-Notiz-Überschriften — ohne eindeutigen Abschnitt wird nicht gemessen"
  mit Zeile und Stelle; Exit 1.
- geprüft, ohne Befund: **Wirkung der unicode-weiten Grenze auf die eigenen
  aktivierten Regeln.** Die drei in `.d-check.closure.yml` aufgenommenen
  `structure`-Regeln nutzen ausschließlich `non-empty`; `hasMarker` wird von
  keiner erreicht. `make verify-closure-notes` läuft unverändert über 338
  Dateien mit null Befunden. Die abgelöste Byte-Grenze hatte keinen zweiten
  Aufrufer (die Funktion ist vollständig entfernt), die Floskel-Wortgrenze ist
  unberührt.
- geprüft, ohne Befund: **Marken-Gegenrichtungen.** Über eine 16-Fälle-Matrix:
  Ziffer, Umlaut, kyrillischer Buchstabe und arabisch-indische Ziffer setzen die
  Marke fort (Befund); Bindestrich, Doppelpunkt, Punkt, Leerzeichen mit Klammer
  sowie das Marken-Ende beenden sie (kein Befund). Eine Marke, die Präfix einer
  längeren ist, trifft nicht auf die längere (`Beleg` gegen `**Belege:**` ⇒
  Befund). Bei `require-all` mit zwei Marken meldet ein fehlender Eintrag genau
  einen Befund, der die **erste** fehlende Marke benennt; fehlen beide, bleibt
  es bei einem — verhaltensgleich, weil der Dedup-Schlüssel die Meldung nicht
  enthält. Emoji und ein kombinierendes Zeichen gelten als Grenze; das ist
  spezifikationskonform („nicht-alphanumerisch"), aber in N-5 als
  Vertrags-Unschärfe mitgeführt.
- geprüft, ohne Befund: **`MR-025` nach der Schärfung.** Das von der Regel
  vorgeschlagene `grep` nach dem vorigen Vertreter (`targets`) über den Baum
  liefert 56 Dateien. Nach Abzug der historischen Klassen (`done/`-Slices,
  Review-Reports, ADRs) und der vom Slice ohnehin berührten Dateien bleiben
  zehn; von diesen sind acht die Implementierungs- und Testdateien des
  **vorigen** Moduls (deren Analogon dieser Slice als `structure.go` und
  `structure_test.go` neu angelegt hat), einer ist ein Hook-Skript ohne
  Modul-Enumeration, und die Einträge in `.d-check.yml` sind bewusste
  `ignore-refs`-Ventile auf umgezogene Pfade. Die vier in §3b bilanzierten
  Lücken decken sich mit dem Befund des Erst-Reviews; die **fünfte** Lücke
  liegt in einer Datei, die der Slice berührt hat, und ist deshalb weder vom
  Datei-`grep` noch von der Bilanz erfasst worden ⇒ N-1.
- geprüft, ohne Befund: **ADR-0049 Fitness Function.** Alle sieben Punkte
  gemessen; insbesondere „Mehrdeutigkeit schlägt Messung: ein Dokument mit zwei
  passenden Abschnitten und einem zu dünnen ersten meldet **nur** den
  Mehrdeutigkeits-Code" trifft jetzt wörtlich zu (genau ein Befund im Fixture).
  Der beim Erst-Review offene Punkt ist damit geschlossen.
- geprüft, ohne Befund: **Mutations-Echtheit.** Acht Rückbauten selbst
  nachgestellt, sechs davon rot: Schwelle am Rand, Ziffern-Wortgrenze,
  `TrimSpace` im Selektor, fail-closed bei unlesbarer Datei, Mehrdeutigkeit erst
  ab drei Treffern, Mehrdeutigkeit meldet den ersten statt den zweiten Treffer.
  Die beiden grün gebliebenen stehen in N-6. Alle Dateien nach jedem Lauf aus
  einer Kopie wiederhergestellt, nicht per Versionskontrolle.
- geprüft, ohne Befund: **Gate-Läufe.** `make lint`, `make test`, `make gates`
  (einschließlich `arch-check`, `semgrep`, `gate-consistency`, `planning-check`,
  Coverage 94,20 % ≥ 93 %), `make verify-closure-notes` und `make fullbuild`
  (mit `image-test`, `bench` Median 718 ms, `completeness-check`) sind auf dem
  unveränderten Arbeitsbaum grün.
- geprüft, ohne Befund: **Vertragsflächen der Release-Koordinaten.** Beide
  README-Pins, alle Handbuch-Beispiele, `version.md` (Aktiv-Zeile **und**
  wandernder Anker auf `v0.57.0`, neue Tabellenzeile), Handbuch-Kopfstempel
  1.48, `docs/user/operations.md`-Optionstabelle (führt alle 20 Module,
  einschließlich `structure`), Modul-Liste in
  [`DC-FA-CLI-002`](../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl),
  `validModules()`, §4-Grund-Code-Tabelle, `AllReasons()`, `reasonTexts()`,
  Handbuch-§4- und -§11-Zeile.
- geprüft, ohne Befund: **SemVer-Einstufung.** **Minor** ist korrekt und war
  es auch schon vor der Heilung: additives Modul, additive Grund-Codes, nichts
  entfernt, Default-Befundsatz ohne aktives `structure` unverändert. Die
  Mehrdeutigkeits-Härte macht einen **opt-in**-Lauf strenger, entfernt aber
  keine zugesagte Fähigkeit; der Grund-Code war seit 0.51.0 spezifiziert. Die
  Einstufung ist damit richtig — unvollständig ist nur die Beschreibung der
  Richtung (N-3).
- geprüft, ohne Befund: **Architektur und Import-Regeln.** `structure.go`,
  `sections.go` und die geänderte `planning.go` liegen im Kern und importieren
  nur `model` und `port/driven`; `make arch-check` innerhalb `make gates` grün.
- geprüft, ohne Befund: **Seiteneffektfreiheit.** Alle Läufe dieses
  Re-Reviews mit `--network none` und `:ro`-Mount; Fixtures ausschließlich
  außerhalb des Repos; `git status --short` am Ende leer.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** neues Modul nicht in der
Verfügbar-Enumeration der eigenen emittierten Vorlage · Zahl im
Beschreibungstext nachgezogen, Enumeration in den Akzeptanzkriterien nicht
(invers) · Release-Notiz nennt nur eine Richtung einer zweiseitigen
Verhaltensänderung · Preset-Partner-Spezifikationen driften an derselben
geteilten Mechanik · geteilte Lexik mit zwei Wort-Begriffen, nur einer im
Vertrag · Schwellen- und fail-closed-Ränder eines neuen Moduls ohne Test ·
Doc-Kommentar nach Signatur-Änderung gestapelt statt ersetzt · wirkungsloser
Zweig in Verriegelungs-Gestalt

Aus dem Erst-Lauf **fortbestehend** (unverändert, daher mit ihrer dortigen
Klasse zu zählen): Zähl-Kommentar als Wächter formuliert, aber nicht
mitgezogen (F-9) · Grund-Code in nur einer von mehreren Handbuch-Enumerationen
nachgezogen (F-10) · Doku-Drift in der READMEs-Statuszeile (F-11) ·
undokumentierte Kandidaten-Ausnahme im Baum-Walk (F-12) · Fehler-Rückgabe im
Post-Pass verworfen statt gemeldet (F-13).

## Verdikt

**Merge-blockierend:** nein — kein HIGH; die vier MEDIUM blockieren den Merge
des Codes nicht, weil keiner von ihnen die Mechanik betrifft. Das ist die
begründete Abweichung von der Regel „MEDIUM blockiert typischerweise": N-1 und
N-2 sind Enumerationen, N-3 ist eine Release-Notiz, N-6 sind fehlende
Negativtests an einer Stelle, deren **Produktverhalten** ich am Lauf als korrekt
gemessen habe.

**Die Heilung ist echt.** Der HIGH-Befund ist nicht kosmetisch, sondern
mechanisch geschlossen: `closure-note-ambiguous` entsteht, schließt die drei
Messungs-Codes aus, meldet den **zweiten** Treffer, verhält sich bei drei
Treffern und bei gemischten Ebenen wie zugesagt, zählt eine Überschrift im
Fence nicht, und der Klartext erscheint im `--doctor`-Lauf. Sechs von acht
Rückbauten sind rot. Die Marken-Grenze ist unicode-weit und in beide Richtungen
belegt. Der Grund, aus dem der eigene Bestand grün bleibt, ist nachgemessen und
nicht angenommen — er hängt an einer Fence-Zeile in genau einem von 99 Slices.

**Die Katalog-Lücke des Erst-Reviews ist eine Einzelfall-Lücke gewesen.** Über
alle 44 Grund-Codes geprüft: keiner ist deklariert, dokumentiert und tot. Die
Verriegelung prüft weiterhin nur Katalog-Deckung — dass daraus kein zweiter
Schaden entstanden ist, ist gemessen, nicht mechanisiert.

**Release-Empfehlung: `v0.57.0` taggen, nachdem N-1 und N-3 geschlossen sind.**
Die beiden sind billig und sie sind die einzigen, die einen **Konsumenten**
treffen: N-1 versteckt das neue Modul vor genau dem Aufruf, auf den die
reparierte Vorlage verweist, und N-3 lässt einen heute grünen Closure-Lauf
unangekündigt rot werden — das ist die Zusage, die
[ADR-0049](../plan/adr/0049-structure-modul-schnitt-und-preset.md)
ausdrücklich in die Release-Notiz verlangt hat. N-2, N-4, N-5, N-7, N-8 und die
fünf fortbestehenden Erst-Befunde sind Vertrags- und Doku-Präzision; sie
gehören in die Closure, nicht vor den Tag. N-6 ist vor dem Tag zu **entscheiden**
(zwei Tests oder eine benannte Grenze), blockiert das Modul aber nicht.

Bemerkenswert für den Steering-Loop: die geschärfte
[`MR-025`](../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten)
verlangt, die Liste per `grep` **aus dem Repo** abzuleiten. Genau dieser `grep`
hätte N-1 **nicht** gefunden, weil die Datei berührt wurde — nur die Stelle in
ihr nicht. Die Regel unterscheidet Datei und Stelle noch nicht; das ist der
nächste Schärfungs-Punkt, nicht ein Versagen der Regel.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in
die Slice-Closure §7 und von dort in den Steering-Loop-Zähler. Dieser Report ist
ein Lauf-Beleg (dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) und
ersetzt keine Verifikation.
