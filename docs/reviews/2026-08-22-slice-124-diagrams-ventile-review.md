# Review-Report: slice-124 — `diagrams`: Ventil-Parität und die fehlenden Schema-Zeilen

**Datum:** 2026-08-22 · **Review-Art:** Code- und Design-Review (Lastenheft/Spezifikation/ADR gegen Implementierung), unabhängiger Reviewer ohne Anteil an der Arbeit
**Gegenstand:** die vier Slice-Commits der Kette (Beanspruchung · Lastenheft-CR 0.65.0 · ADR-0058 Entscheidung 3 · Implementierung samt Spezifikation, §2-Schema-Nachtrag und Config-Kommentar)
**Skill:** `.harness/skills/reviewer.md` @ 1.9.0 · **Modell-ID:** `claude-opus-5[1m]`
**Vom Reviewer selbst gefahren:** `make gates` Exit 0, `make test` Exit 0, dazu eigene Image-Läufe gegen das gepinnte Vorgänger-Image und ein aus dem Vor-Slice-Stand gebautes Image

**Verdikt des Reviews: blockierend** — kein HIGH, sieben MEDIUM, acht LOW. Alle eingearbeitet.

## Befunde und Einarbeitung

### F-1 · MEDIUM — „das einzige Modul, das Befunde an Zeilen hängt"

Die Aussage stand in Lastenheft (Rang 1), §7-Historie, ADR und beiden
Commit-Botschaften — und ist falsch: `hostpaths`, `pins` und `spans` hängen
Befunde ebenfalls an Zeilen und tragen **auch nach diesem Slice** kein Ventil.
Die ADR erhebt die Klasse zum **Prinzip** und indiziert damit drei Module,
ohne es zu sagen. Dieselbe Klasse wie F-2 des slice-123-Reports.

**Eingearbeitet:** die Reichweite ist auf das Belegbare gezogen — `diagrams`
war das einzige Modul mit **konfigurierbaren Mustern** ohne Ventil. Die drei
übrigen sind als **offene Fläche** benannt (ihr Befund folgt aus einer festen
Lexik, nicht aus eigenen Mustern), mit einem eigenen Re-Evaluierungs-Trigger.

### F-2 · MEDIUM — die Ventile gelten auf der Ziel-Achse nicht, und der Vertrag schwieg

Der Slice erzeugt eine Asymmetrie: scan-seitig nimmt der Marker eine Zeile
aus, auf der **Definitions**-Seite nicht. Eine im `defined-in` als illustrativ
markierte Kennung definiert weiter. Der Reviewer hat es gemessen (mit und ohne
Marker in der Quelle: identisch) — genau die `BEO-004`-Frage, die der
Slice-Plan selbst als einschlägig benannt hatte.

**Eingearbeitet:** als **Festlegung** ausgesprochen, nicht als Auslassung
stehengelassen — in Anforderung und Algorithmus: die Ventile unterdrücken
Befunde, sie ändern nicht die Wahrheit über die Definitionsmenge.

### F-3 · MEDIUM — die Marker-Hälfte ist nicht opt-in

`exempt-paths` hängt an einem Schlüssel, der Marker am **Inhalt**. Eine
Diagramm-Zeile, die die Zeichenfolge ohnehin trägt, wird **ohne** jede
Konfigurations-Änderung stumm. Die Wellen-Zusage („keine bestehende
Konfiguration ändert ihr Verhalten") ist damit für die Hälfte der Erweiterung
falsch. Selbst nachgemessen, identische Konfiguration ohne neuen Schlüssel:

```
v0.62.0: 2 Befunde   |   neuer Stand: 1 Befund
```

**Eingearbeitet:** als benannte Grenze in Wellendokument, ADR-Konsequenz und
Akzeptanzkriterium — die Byte-Identität gilt für die Marker-Hälfte nur für
Bäume ohne die Zeichenfolge in einer gelisteten Fence. Das ist der Preis
dafür, denselben Marker zu benutzen statt einen eigenen zu erfinden.

**Register:** dieselbe Klasse wie `BEO-009`, zweite Richtung (Zusage reicht
weiter als die Messung). Der Eintrag steht seit slice-123 bei **3** und ist
verkörpert — zitiert, nicht weitergezählt.

### F-4 · MEDIUM — die `diagrams.scope`-Schema-Zeile war falsch und überflüssig

Default „leer" ist in beiden Lesarten falsch: *abwesend* heißt **globaler
Scope**, ein wörtlich leeres Objekt ist **Exit 2**. Zudem war es die einzige
modul-spezifische Scope-Zeile der Tabelle — alle übrigen Module leben von den
zwei generischen Zeilen. **Eingearbeitet:** Zeile gestrichen.

### F-5 · MEDIUM — `fences: []` fällt still auf den Default zurück

Die §2-Präambel sagt generisch, eine explizit leere Liste ersetze den Default
durch die leere Menge. Für `diagrams.fences` gilt das nicht — `len == 0`
liefert `["mermaid"]`. Wer die Liste leert, um das Modul stillzulegen, bekommt
weiter mermaid geprüft. **Eingearbeitet:** die Abweichung steht in der
Default-Spalte. Das Verhalten selbst zu ändern wäre eine Verhaltensänderung
und braucht eine eigene Entscheidung.

### F-6 · MEDIUM — `defined-in` sagte „Datei", der Rand akzeptierte ein Verzeichnis

Die Prüfung beim Lauf-Start ließ ein Verzeichnis durch; `definedTokens`
bekommt dann einen Lesefehler und liefert fail-closed die **leere** Menge —
jede Kennung im Diagramm wird „undefiniert". Wer die Gewohnheit von
`ids.patterns[].target` überträgt (dort ist ein Verzeichnis legal), bekommt
einen Befund-Sturm ohne Hinweis auf die Ursache.

**Eingearbeitet — Code, nicht Text:** ein Verzeichnis bricht jetzt mit Exit 2.
Nachgemessen: vorher zwei Befunde, jetzt
`diagrams.patterns[0].defined-in ist keine Datei: spec`. Die Verengung ist als
**Verhaltensänderung** in der §7-Historie benannt, samt der Abgrenzung zu
`ids`.

### F-7 · MEDIUM — das Spec-Stratum nannte eine Welle

`AGENTS.md` §3.4 ist unbedingt: kein Spec-Stratum referenziert ADRs, Wellen,
Slices oder Commit-Hashes. Dieselbe §7-Zeile vermeidet die ADR-Nennung
ausdrücklich — die Regel war also bewusst angewandt und an einer Stelle
gebrochen. Kein Gate fängt es (die `matrix`-Klassen kennen nur `slice-\d{3}`).
**Eingearbeitet:** „am eigenen Profil gemessen", ohne Wellen-Kennung.

### F-8 bis F-15 · LOW — eingearbeitet

- **F-8:** die Grenzen-Aufzählung des Config-Kommentars sprang von (2) auf
  (4), weil der Slice „Grenze 3" zu „Grenze" gemacht hatte. Nummerierung
  wiederhergestellt.
- **F-9:** `exempt-paths` wurde als Ganzes validiert, zur Laufzeit aber
  segmentweise gematcht — ein nur als Ganzes gültiges Muster war still
  wirkungslos; ein leerer Eintrag ebenso. **Eingearbeitet:** die `tracked`-Form
  übernommen (leerer Eintrag = Fehler, Validierung je `/`-Segment mit
  `**`-Ausnahme), als geteilte Helferin, damit es die dritte Stelle nicht neu
  erfindet.
- **F-10:** die Glob-Prüfung stand unter „Validierung beim Lauf-Start",
  liegt aber am Config-Rand. Verortet, wo sie ist.
- **F-11:** der fünfte Test war vollständig subsumiert — keine Mutation konnte
  ihn allein rot machen. **Ersetzt** durch einen, der eine eigene Aussage
  trägt: nur das **exakte** Token schaltet still (`d-check:ignor` und
  `dcheck:ignore` nicht).
- **F-12:** Herkunfts-Prosa im Test-Kommentar („der Fall, der welle-80 zum
  Scoping zwang"). Auf die Zusage gekürzt.
- **F-13:** „die beiden Module, die den Marker honorieren" — `versions`
  honoriert ihn ebenfalls, und zwar dokumentiert. Er ist jetzt als dritte
  Präzedenz genannt (Lastenheft und Algorithmus).
- **F-14:** der Marker auf der **schließenden** Fence-Zeile war still
  wirkungslos — genau das Versagen, das die ADR bei der Zeilen-only-Variante
  zu Recht als „schlimmer als keines" bezeichnet. **Eingearbeitet:** ein
  Negativtest und ein Halbsatz im Algorithmus.
- **F-15:** die Messung „446 Dateien, null Befunde" beantwortet nicht die
  Frage, der sie zugeordnet war: sie zeigt, dass heute kein Diagramm außerhalb
  `spec/` ein `ARC-\d{3}`-Token trägt — und das galt vorher genauso. In der
  Closure-Notiz als **Bestandsaufnahme** benannt; den Wellen-Closure-Trigger
  („eine konstruierte Gegenprobe, die ohne die Erweiterung stumm bliebe")
  erfüllt die `exempt-paths`-Exit-2-Probe, nicht diese Zahl.

## Negativbefunde des Reviews (geprüft, ohne Befund)

- **Die Öffnungszeilen-Semantik ist korrekt** — acht konstruierte Fence-Formen
  gemessen: nicht gelisteter Fence mit Marker auf der Öffnungszeile (kein
  Leck), zwei Fences hintereinander (nur der markierte frei), eingerückter
  Fence, `~~~mermaid`, unbalanciert, verschachtelt. Die zwei Abweichungen bei
  unbalanciert/verschachtelt sind **vorbestehend** und Folge des bewusst
  naiven Toggle-Automaten, nicht des Umbaus.
- **Kein unerwarteter Nebeneffekt für andere Module:** der Marker auf einer
  Fence-Öffnungszeile ist für `codepaths`/`ids` unsichtbar (Fence-Zeilen
  werden verworfen) und für `versions` genau eine Zeile, nie ein Block. Über
  das ganze Repo gibt es **keine** Fence-Öffnungszeile mit der Zeichenfolge —
  es wird nichts Bestehendes neu stillgeschaltet.
- **`DC-QA-02` gemessen, nicht geglaubt:** ein aus dem Vor-Slice-Stand
  gebautes Image gegen das neue auf demselben Arbeitsbaum — stdout und stderr
  byte-identisch, Exit 0/0; roter Fixture gegen `v0.62.0` byte-identisch
  (123 B), Exit 1/1.
- **Beide Ventile mutationsgeprüft** (die Botschaft stimmt), **jedes
  Akzeptanzkriterium hat einen Test**, das Datei-Ventil greift **nach** dem
  Scope und nutzt dieselbe `path.Match`-Semantik wie `scan.ignore`.
- **Die `defined-in`-Quelle überlebt ihr eigenes `exempt-paths`** — das Ventil
  ist scan-seitig, wie beabsichtigt.
- **Der `--print-config`-Block übersteht das Einkommentieren** (die
  slice-122-Lehre ist angewandt), **kein neuer Grund-Code**, `Schärft:`
  vollständig, Lifecycle-Commit korrekt nach MR-013.

## Für den Release-Prep benannt (slice-125)

Der Reviewer hat die Handbuch- und README-Widersprüche einzeln aufgeführt;
sie stehen auf der Arbeitsliste des Release-Slice. Der schwerste ist **älter
als diese Welle**: Handbuch und beide READMEs behaupten, der Zeilen-Marker
wirke ausschließlich für `codepaths` und `ids` — `versions` honoriert ihn seit
seiner Einführung, jetzt kommt `diagrams` dazu.

## Nachmessung nach der Einarbeitung

- `make gates` grün (acht Gates, Exit 0 explizit gelesen).
- F-3 nachgestellt: identische Konfiguration, `v0.62.0` zwei Befunde, neuer
  Stand einer.
- F-6 nachgestellt: Verzeichnis als `defined-in` — vorher zwei Befunde
  (Exit 1), jetzt Exit 2 mit benannter Ursache.
