# Architect-Verdikt — welle-95 (Modul `linelength`), Plan gegen Entscheidungslage

**Rolle:** Architect (Modul 8), frischer Kontext, **vor** Code.
**Gegenstand:** [welle-95](../plan/planning/done/welle-95/welle-95-linelength-modul.md) und ihre zwei Slice-Pläne
[slice-289](../plan/planning/done/welle-95/slice-289-linelength-anforderung-und-adr.md) (Entscheide 1–8 in §3,
vom Auftraggeber am 2026-11-22 bestätigt) und
[slice-290](../plan/planning/done/welle-95/slice-290-linelength-modulkern.md).
**Anlass:** Auftrag des Auftraggebers in der Sitzung vom 2026-11-22 (kein Issue).
**Datum:** 2026-11-22.
**Geändert wurde kein Plan.** Mit diesem Verdikt entsteht die ADR des Moduls als
[ADR-0114](../plan/adr/0114-linelength-modul.md) (Status `Proposed`) samt Index-Zeile; sie trägt die
Folge-Entscheide aus diesem Verdikt als Empfehlung.

## Gesamturteil

Die Richtung trägt: eigenes opt-in-Modul über der Scan-Menge, Zählung in Codepunkten, Fence- und
Tabellen-Antwort aus der geteilten Lexik, hermetisch, kein Schreiben. Keine `Accepted`-ADR wird verletzt,
keine Schwelle gesenkt. **Ein Entscheid widerspricht dem Lastenheft** (F-1: Entscheid 8 gegen
[`DC-FA-CONF-002`](../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope)) und muss neu
bestätigt werden. **Zwei Entscheidungen fehlen** (F-2 Zeilen-Marker, F-7 Bezugsmenge). **Vier Stellen
behaupten über die Lexik oder die Nachbar-Module mehr oder anderes, als der Gegenstand hergibt** (F-3, F-4,
F-5, F-6). Dazu eine Lücke im ADR-Lebenslauf (F-8). Alle Folge-Entscheidungen stehen in der **einen** ADR
des Moduls; eine zweite ADR, ein Beobachtungs-Eintrag oder ein Folge-Slice ist für keinen Befund nötig.

Verdikt-Kürzel (Modul 8, Konflikt-Pfad): **V1** = Entscheidung gilt, Plan hat falsch behauptet/ausgelassen ·
**V2** = Folge-Entscheidung nötig · **V3** = Lockerung legitim, aber undokumentiert.

## Befunde

### F-1 — Entscheid 8 (leere Kandidaten-Menge ⇒ Exit 2) widerspricht `DC-FA-CONF-002` und beruft sich auf eine falsche Präzedenz — **V1**

- Entscheid 2 macht `linelength` zu einem **scannenden** Modul mit modul-lokalem `scope`
  ([slice-289](../plan/planning/done/welle-95/slice-289-linelength-anforderung-und-adr.md) Z. 117–123). Für jedes
  scannende Modul sagt das Lastenheft (Rang 1) zu: „Given ein Modul-`scope` mit explizit leerer `roots`-Liste
  … then prüft dieses Modul nichts — bewusste Setzung, kein Fehler“
  ([lastenheft.md](../../spec/lastenheft.md) Z. 4074, Beschreibung Z. 4066–4068). Entscheid 8
  (slice-289 Z. 154–155) macht genau diesen Fall zu Exit 2. Eine ADR kann das Lastenheft nicht schärfen, und
  eine zweite Lastenheft-Anforderung, die einer bestehenden widerspricht, wäre eine Drift im obersten Rang.
- Die Präzedenz „fail-closed wie `structure`/`file`“ trägt nicht: beide sind **nicht** scannende Module mit
  eigenen Globs (`DC-FA-CONF-002` Z. 4049–4054, `DC-FA-FILE-001` Z. 4143–4145), und ihre Nullmenge ist ein
  **Befund** mit Exit 1, kein Exit 2 (`file-no-match`, [lastenheft.md](../../spec/lastenheft.md) Z. 4163–4167;
  [spezifikation.md](../../spec/spezifikation.md) Z. 3348–3351, `SPEC-084`). Ihre Härte folgt daraus, dass eine
  Regel mit Glob **behauptet**, Dateien zu treffen; `linelength` behauptet nichts über seine Dateien, es erbt
  die Scan-Menge wie `links`, `ids` und `anchors`, und die laufen über einer leeren Menge ohne Fehler (die
  Default-Wurzeln sind optional, `DC-FA-SCAN-001`).
- **Verdikt:** Das Lastenheft gilt, der Plan hat die Präzedenz falsch behauptet. Empfehlung (in ADR-0114
  Entscheid 8): eine leere Menge ist **kein** Fehler; die Bezugsmenge (Entscheid 7) nennt `0 Datei(en)` und
  macht das Grün damit sichtbar — die Antwort, die Baseline-Regelwerk `modul-13-quality-gates.md`
  §Hard Rule (*Ein Gate ohne seine Grenze*) verlangt, ohne neuen Grund-Code. Hält der Auftraggeber an
  „fail-closed“ fest, ist die mit `DC-FA-CONF-002` verträgliche Form ein **Befund** (Exit 1) nur für eine nicht
  explizit leere Menge — als Alternative in der ADR geführt.
- **Plan-Änderung:** slice-289 §3 Entscheid 8 neu fassen und vom Auftraggeber neu bestätigen lassen; §4
  Rückführung `in-progress → open` um Punkt 8 ergänzen; slice-290 §3 Z. 143 („Keine Markdown-Datei ⇒
  fail-closed“) auf das neue Verhalten und um den Fall „explizit leere `scope.roots`“ ergänzen.

### F-2 — Zeilen-Marker `d-check:ignore`: Entscheidung fehlt — **V2** (Entscheid 9 in ADR-0114)

- Der Marker ist eine **benannte Liste** honorierender Module (`codepaths`, `ids`, `versions`, `diagrams`;
  [codepaths.go](../../internal/hexagon/core/rules/codepaths.go) Z. 15–18, Historie
  [spezifikation.md](../../spec/spezifikation.md) Z. 4184), kein ableitbares Kriterium. Das Lastenheft setzt die
  Regel für alle übrigen: „deterministische Befunde werden behoben oder strukturell ausgenommen …, nicht
  zeilenweise stummgeschaltet“ ([lastenheft.md](../../spec/lastenheft.md) Z. 1392–1396, `DC-FA-MTX-001`).
  Die Form des Markers folgt der Kommentar-Lexik der Eingabe
  ([ADR-0063](../plan/adr/0063-marker-form-folgt-der-kommentar-lexik-der-eingabe.md)).
- Der Plan sagt dazu nichts. Ohne Entscheid liest der Implementer den Marker entweder als geltend (dann wächst
  die benannte Liste still um ein fünftes Modul, und ihre zwei Spiegel in Spezifikation und Nutzer-Doku
  driften) oder als nicht geltend (dann fehlt die Aussage in der Anforderung).
- **Verdikt:** Folge-Entscheid. Empfehlung: `linelength` honoriert den Marker **nicht** — eine überlange Zeile
  ist deterministisch und behebbar, die strukturellen Ventile sind (a)–(c), `scope` und `exempt-paths`. Folge,
  in der ADR benannt: ein Marker-Kommentar verlängert seine Zeile und zählt mit, und weil er Leerzeichen trägt,
  greift (c) auf ihm nie.
- **Plan-Änderung:** slice-289 DoD-Punkt 1 nimmt ein Akzeptanzkriterium „Zeile über der Grenze mit
  `<!-- d-check:ignore … -->` ⇒ Befund“ auf; slice-290 §3 den Testfall.

### F-3 — Ausnahme (a) trägt über die geteilte Lexik, aber mit drei Grenzen, die der Plan nicht nennt — **V1**

Am Gegenstand ([markdown.go](../../internal/hexagon/core/rules/markdown.go)):

- `proseLineSet` (Z. 180–186) baut auf `proseLines` (Z. 88–103): `strings.Split` an `\n`, `TrimFenceIndent`
  (Z. 52–54, entfernt **jede** Einrückung aus Space und Tab), dann `FenceToggle` (Z. 34–46). Fence-Zeilen selbst
  entfallen (`continue`) — „samt ihrer Zäune“ (Entscheid 5) trägt. **Aber** der Automat **schaltet bei jeder** Zeile
  um, die mit ```` ``` ```` oder `~~~` beginnt; die strenge Schluss-Lesart `FenceCloses` (Z. 77–83) benutzt er nicht.
  Daraus, und das ist die geteilte Antwort, nicht ein Fehler des neuen Moduls
  ([ADR-0054](../plan/adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 1):
  1. Ein ````` ```` `````-Block mit einer inneren ```` ``` ````-Zeile, und ein ```` ``` ````-Block mit einer
     `~~~`-Zeile, schließen früh; die Zeilen dahinter bis zum nächsten Umschalter werden gemessen.
  2. Ein eingerückter Code-Block (vier Leerzeichen, ohne Zaun) ist kein Fenced-Code und wird gemessen; ein
     Zaun in einem Zitat (```` > ``` ````) schaltet nicht um.
  3. Ein nie geschlossener Zaun nimmt den Rest der Datei aus. Diesen stillen Ausfall meldet bereits
     `spans` mit `fence-unclosed` (`SPEC-030`, [spezifikation.md](../../spec/spezifikation.md) Z. 3990: „mindestens
     eine der beiden Fence-Lesarten endet offen“). Ein zweiter Wächter in `linelength` wäre eine zweite Antwort.
- **Verdikt:** Die Lexik gilt; der Plan lässt ihre Grenzen aus. Kein neuer Mechanismus.
- **Plan-Änderung:** slice-289 DoD-Punkt 1/2: die Anforderung nennt die drei Grenzen und verweist für den
  offenen Zaun auf `spans`. Damit hat slice-289 Risiko 2 seinen Ausgang *entfallen* ohne Register-Eintrag, und
  slice-290 Risiko 1 (Z. 183–186) braucht nicht „Meldung oder Bezugsmenge macht den Teil sichtbar“, sondern den
  Test, der das Verhalten festhält, plus den Verweis. slice-290 §3 Z. 134: Testfälle „`~~~`-Zeile im Backtick-Block“
  und „eingerückter Code-Block zählt“ ergänzen. Der Kopplungs-Test (slice-290 Risiko 2,
  Z. 191) ist **keine** Option, sondern Pflicht nach ADR-0054 Entscheidung 4 („je Konsument eine Assertion“) und
  gehört in DoD-Punkt 2.

### F-4 — Ausnahme (b): die Lexik kennt „Zeile, die mit `|` beginnt“, nicht „Zeile einer Pipe-Tabelle“ — **V1**

- `tableRowLine` (Z. 192–194) ist `prose[i+1] && strings.HasPrefix(lines[i], "|")`: **keine** Trennzeile
  verlangt, **kein** Trimmen. Damit:
  1. Eine `|`-Zeile ohne Trennzeile ist für die Lexik eine Tabellenzeile und ausgenommen. slice-290 §3 Z. 138
     („Tabelle ohne Trennzeile (keine Tabelle)“) erwartet das Gegenteil.
  2. Eine eingerückte Tabellenzeile (in einem Listenpunkt), eine Zeile in einem Zitat (`> |`) und eine
     GFM-Zeile ohne führende Pipe (`a | b`) sind **keine** Tabellenzeilen und werden gemessen.
  3. Die Schalter sind unabhängig: Mit `exempt-code-blocks: false` und `exempt-tables: true` ist eine `|`-Zeile
     **im** Zaun keine Tabellenzeile (sie liegt nicht in `prose`) und wird gemessen.
- **Verdikt:** Die Lexik gilt; Entscheid 5 (slice-289 Z. 140) und die Testliste beschreiben eine andere Antwort.
- **Plan-Änderung:** Entscheid 5 (b) als „Zeile außerhalb von Fenced-Code, die mit `|` beginnt“ fassen (so in
  ADR-0114); die Anforderung nennt die Grenzen 2 und 3; slice-290 Z. 138 auf „`|`-Zeile ohne Trennzeile ⇒
  ausgenommen; eingerückte Zeile ⇒ gemessen“ korrigieren.

### F-5 — Ausnahme (c): die „Alternative“ aus Risiko 1 ist dieselbe Regel — **V1**, und die Ränder fehlen — **V2** (Entscheid 5 (c) in ADR-0114)

- Ausnahme (c): ausgenommen, wenn der Teil ab der Grenze kein Leerzeichen enthält. Das gilt genau dann, wenn
  das **letzte** Leerzeichen der Zeile auf einer Position ≤ `max-chars` steht. Die in slice-289 Z. 202 genannte
  Alternative „nur der Überschuss nach dem letzten Leerzeichen zählt“ misst die Länge bis zum letzten
  Leerzeichen und meldet, wenn sie die Grenze überschreitet — dieselbe Bedingung, bis auf eine einzige
  Randposition. Ein Vergleich mit ihr kann den Ausgang *entfallen* nicht tragen, denn er vergleicht die Regel
  mit sich selbst.
- Zwei Ränder sind offen: was ein Leerzeichen ist (nur U+0020, oder auch Tab) und ab welcher Position „ab der
  Grenze“ zählt.
- **Verdikt:** Die bestätigte Regel bleibt; sie wird präzisiert und gegen echte Alternativen verglichen.
  Empfehlung in ADR-0114: Leerraum sind Space und Tab (dieselbe Menge wie `TrimFenceIndent`); geprüft werden die
  Codepunkte an den Positionen `max-chars + 1` bis Zeilenende. Verglichen: (c-streng) nur ein Wort, das allein
  länger als die Grenze ist, wird ausgenommen; keine Ausnahme (c), sondern ein Marker (gegen F-2); eine
  lexikalische Link-/URL-Erkennung (brauchte eine eigene Link-Lexik, slice-289 Rückführungs-Trigger Z. 176–178).
- **Plan-Änderung:** slice-289 Risiko 1, Ausgang *entfallen*: „die Anforderung nennt die Grenze, die ADR hat die
  strenge Variante verglichen“ statt der gleichwertigen. slice-290 §3 Z. 139: Randfälle „Leerzeichen genau auf
  Position `max-chars + 1` ⇒ Befund“, „letztes Leerzeichen auf Position `max-chars` ⇒ ausgenommen“, „Tab nach
  der Grenze ⇒ Befund“.

### F-6 — Entscheid 1: der Grund gegen `structure` ist am Gegenstand falsch — **V1**

- slice-289 Z. 114–116: `structure` „meldet auf der Überschriftszeile“. Die zeilen-gebundenen Bedingungen von
  `structure` melden auf der Zeile, an der repariert wird (`cell-max-chars` ⇒ `section-cell-oversized` an
  **ihrer** Zeile, [spezifikation.md](../../spec/spezifikation.md) Z. 2724–2728 und Historie Z. 4178).
- Die Gründe, die tragen: `structure` wählt seine Dateien über eigene Globs und kennt kein `scope`
  (`DC-FA-CONF-002` Z. 4052–4054), jede Regel braucht einen Abschnitts-Selektor, und es misst den
  **bereinigten** Text ([ADR-0088](../plan/adr/0088-file-modul-groessengrenzen.md) (a),
  [ADR-0089](../plan/adr/0089-structure-max-lines-zwoelfte-bedingung.md) (b)). Eine Zeilen-Grenze der ganzen
  Datei über der Scan-Menge hätte keinen dieser drei Anker.
- **Verdikt:** Die Wahl (eigenes Modul) bleibt, die Begründung im Plan ist falsch. ADR-0114 führt die
  tragenden Gründe.
- **Plan-Änderung:** slice-289 §3 Entscheid 1 an die ADR angleichen (oder auf sie verweisen).

### F-7 — Bezugsmenge: Zählung und die Ausnahme (c) sind unbestimmt — **V2** (Entscheid 7 in ADR-0114)

- „Geprüfte Zeilen nach Abzug der Ausnahmen“ (slice-289 Z. 150–153) hat zwei offene Stellen. (1) `proseLines`
  zerlegt mit `strings.Split` an `\n`; eine Datei, die mit `\n` endet, liefert ein leeres Schlusselement. Zählt
  das Modul die Elemente, ist die Zahl um eins größer als `wc -l` und als die Zählung von `file`
  ([lastenheft.md](../../spec/lastenheft.md) Z. 4155–4158). (2) (a) und (b) sind Klassen von Zeilen, (c) hängt
  vom Messergebnis ab: eine Zeile wird erst gemessen, dann ausgenommen. Zieht die Zahl (c) ab, verschwindet
  genau die Nachsicht, deren Grobheit Risiko 1 benennt.
- **Verdikt:** Folge-Entscheid. Empfehlung: Zeilen zählen wie bei `file` (Zeilenumbrüche plus eine
  unvollständige Schlusszeile); „geprüft“ ist die Menge nach `scope`, `exempt-paths`, (a) und (b); die über (c)
  ausgenommenen Zeilen stehen als eigene Zahl daneben. Form wie bei `mentions`/`pinsync`: stderr und
  `summary.notes` ([spezifikation.md](../../spec/spezifikation.md) Z. 3312–3318, 3549–3552).
- **Plan-Änderung:** slice-290 DoD-Punkt 2 und §3: Test für die Zahl bei einer Datei mit und ohne
  abschließendes `\n` und für die (c)-Zahl.

### F-8 — ADR-Status: `Accepted` schon mit slice-289 — **V1**

- slice-289 DoD Z. 81–83 verlangt die ADR `Accepted` bei der Closure von slice-289. slice-290 §4 Z. 160–162
  sieht aber eine Rückführung vor, „wenn die geteilte Lexik … eine Antwort liefert, die der Anforderung
  widerspricht“, und §1 Z. 50–53 eine „Plan-Änderung an die Anforderung“. Mit `Accepted` wäre jede solche
  Korrektur eine Folge-ADR mit `supersedes` ([AGENTS.md](../../AGENTS.md) §3.5), wo eine Änderung am Entwurf
  genügt hätte. Dazu tragen F-1, F-2 und F-7 Entscheide, die der Auftraggeber noch nicht bestätigt hat.
- Präzedenz: [ADR-0112](../plan/adr/0112-metrics-modul.md) wurde mit dem **umsetzenden** Slice `Accepted`
  (Geschichte dort, nach demselben Befund im Verdikt zu welle-93).
- **Verdikt:** Plan-Lücke. ADR-0114 bleibt `Proposed` bis zur Closure von slice-290; der Übergang ist dort
  DoD-Punkt (für `make adr-check` erlaubt). Die ADR nennt `DC-FA-LLEN-001` heute ohne Link und mit
  `d-check:ignore`, weil die Anforderung noch nicht im Lastenheft steht; der Link ersetzt den Marker, sobald sie
  steht (Entwurf, also änderbar).
- **Plan-Änderung:** slice-289 DoD-Punkt 3: „ADR-0114 liegt `Proposed` vor, Index geführt, die Kennung der
  Anforderung und das `Schärft:`-Feld sind verlinkt“; slice-290 DoD: „ADR-0114 `Accepted`, Index nachgezogen“.

## Bestätigt, mit akzeptiertem Negativ

- **Spiegel-Falle aus welle-93 F-1 — vermieden.** slice-289 Z. 49–53 lässt die Modul-Aufzählung in
  [`DC-FA-CLI-002`](../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl) (Lastenheft Z. 111), das Glossar
  (Z. 4460) und §4 bei slice-290. Die Spiegel-Tests
  ([registry_mirror_test.go](../../internal/hexagon/core/app/registry_mirror_test.go) Z. 108–160) und
  `TestPrintConfigVerfuegbarDecktRegistry`
  ([cli_acceptance_test.go](../../internal/adapter/driving/cli/cli_acceptance_test.go) Z. 3300–3326) lesen genau
  die Fundorte, die slice-290 §3 führt. Einen Test, der die §2-Schema-Zeilen gegen den Decoder hält, gibt es
  nicht (gesucht in allen `_test.go`, die `spezifikation.md` lesen); die Schema-Zeilen dürfen deshalb mit
  slice-289 kommen, wie bei `pinsync`.
- **Reihenfolge-Abweichung: slice-289 liegt in `in-progress/`, bevor dieser Schritt lief.** Es gibt keinen Code
  und kein Implementer-Artefakt; der Schaden ist null, solange der Implementer erst beginnt, nachdem der Planner
  dieses Verdikt in die Pläne gezogen hat. **Akzeptiertes Negativ, einmalig**: kein Beobachtungs-Eintrag, weil
  der Fall ohne Wiederholung keine Lücke belegt.
- **Schlüsselnamen `exempt-code-blocks`/`exempt-tables`/`exempt-unbreakable` (Bool) neben `exempt-paths`
  (Liste).** Gleiches Präfix, anderer Typ. **Akzeptiertes Negativ**: eine Verwechslung ist ein Typfehler am
  Konfigurations-Rand und damit laut; die Namen hat der Auftraggeber bestätigt.
- **Überschriften zählen mit (Entscheid 5).** Eine ATX-Überschrift lässt sich nicht umbrechen; die Reparatur
  ist Kürzen, und Kürzen ändert den Anker. **Akzeptiertes Negativ** mit Grund: die Folge ist laut (`links`/
  `anchors` melden jeden gebrochenen Verweis), und die Anforderung nennt es. Steht als Konsequenz in ADR-0114.
- **Entscheid 2 und die Frage aus [AGENTS.md](../../AGENTS.md) §3.8.** Die Scan-Menge enthält nur Markdown
  (`DC-FA-SCAN-001`, Lastenheft Z. 1007–1008); „andere Dateiarten liest das Modul nicht“ gilt deshalb per
  Konstruktion. Gelesen wird außer der Konfiguration keine Datei außerhalb der Scan-Menge. Die Antwort „keine“
  stimmt.
- **Entscheid 3.** Dieselbe Einheit wie `cell-max-chars` (`utf8.RuneCountInString`,
  [structure_tablecell.go](../../internal/hexagon/core/rules/structure_tablecell.go) Z. 127;
  [ADR-0069](../plan/adr/0069-zellenlaenge-als-strukturbedingung.md) Entscheidung 3), mit einem Unterschied, der in
  die Abgrenzung gehört: die Zelle wird getrimmt gemessen, die Zeile roh (Einrückung zählt).
- **Entscheid 4.** `max-chars` meidet die dritte Lesart von `max-lines`
  ([`BEO-ALL/max-lines-namenskollision-structure-file`](../plan/planning/observations/BEO-ALL/max-lines-namenskollision-structure-file/observation.md)).
  Der Name liegt nahe an `cell-max-chars`, meint aber dieselbe Einheit — kein Befund.
- **Entscheid 6.** Länge und Grenze in `message` folgen der Präzedenz von `file-lines-exceeded`; `message` ist
  nicht stabilitätsgarantiert (`SPEC-001`, [spezifikation.md](../../spec/spezifikation.md) Z. 3580–3589). Was
  `target` trägt, ist Spezifikations-Sache in slice-289 und braucht keine Entscheidung hier.
- **Rohe Messung gegen ADR-0054 Entscheidung 2.** Die Länge einer Zeile ist eine **andere Frage** als „ist das
  Prosa“ — wie bei `file`; die Fragen „ist das Fenced-Code“ und „ist das eine Tabellenzeile“ bekommen die
  geteilte Antwort. Inline-Code wird deshalb zu Recht nicht geleert.
- **Hermetik und Determinismus.** Nur Dateizugriff über die Scan-Menge, Reihenfolge nach Pfad und Zeile —
  deckt [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) und
  [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit). Kein neuer Port;
  [architecture.md](../../spec/architecture.md) nennt keine Modulnamen und bleibt zu Recht unberührt.
- **Abwärts-Sperre.** Die Pläne lassen die Spec-Straten nie auf Slices oder ADRs zeigen (slice-289 Z. 26–27);
  die Kennung entsteht im Lastenheft ([AGENTS.md](../../AGENTS.md) §5 Zeile 3).

## Nicht geprüft

- Die Gegenprobe am eigenen Bestand (Zahl der Zeilen über 120 mit und ohne Fenced-Code) — Auftrag von slice-289.
- Ob der Konfigurations-Decoder einen falschen Typ unter `linelength.*` wirklich mit Exit 2 abweist; gelesen,
  nicht gefahren.
- Ob die Aufrufer von `tableRowLine` ihre Zeilen ebenfalls an `\n` zerlegen; gelesen ist nur die Funktion selbst.
- Keine Gates gegen die Pläne gefahren außer `make doc-check` über dieses Verdikt und die ADR.

## Zusammenfassung

| Befund | Verdikt | Träger der Korrektur |
|---|---|---|
| F-1 leere Menge ⇒ Exit 2 gegen `DC-FA-CONF-002` | V1 Plan falsch, Neu-Bestätigung | Plan 289/290; ADR-0114 Entscheid 8 |
| F-2 Zeilen-Marker | V2 Folge-Entscheid | ADR-0114 Entscheid 9; Plan 289/290 |
| F-3 Fence-Ausnahme, drei Grenzen, Verweis auf `spans` | V1 Plan lässt aus | Plan 289/290 |
| F-4 Tabellen-Ausnahme = „beginnt mit `\|`“ | V1 Plan falsch | Plan 289/290; ADR-0114 Entscheid 5 (b) |
| F-5 Ausnahme (c): Schein-Alternative, Ränder | V1 + V2 | Plan 289/290; ADR-0114 Entscheid 5 (c) |
| F-6 Grund gegen `structure` | V1 Plan falsch | Plan 289; ADR-0114 Kontext/Alternativen |
| F-7 Bezugsmenge: Zählung, (c)-Zahl | V2 Folge-Entscheid | ADR-0114 Entscheid 7; Plan 290 |
| F-8 ADR `Proposed` bis slice-290 | V1 Plan-Lücke | Plan 289/290 |

Kein Befund fällt unter V3: keine Schwelle, kein Gate wird gelockert.
