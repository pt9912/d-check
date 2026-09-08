# Review-Report: slice-216 — 2026-09-08

**Review-Art:** Plan-Review — geprüft wird ein **Entscheid** gegen die Bitte
eines Adopters: trägt die Messung, auf der das Nein ruht, und ist die Antwort
je Argument fair?

**Gegenstand:** slice-216 · Commit-Range `HEAD~2..HEAD` (`9c26494a` Plan +
Vorprüfungen · `215585fb` Beanspruchung · `a33b8547` Entscheid). Geändert: das
CR-Dokument, der Slice-Plan, der Ruhe-Marker der Roadmap.

**Skill:** `.harness/skills/reviewer.md` @ 1.16.0 ·
**Modell:** `claude-opus-5[1m]` · **Datum:** 2026-09-08

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als Tag plus Pfad in Inline-Code
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>) statt als Link.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan slice-216 (`in-progress/`), §1 Abgrenzung, §2 DoD, §3 Plan,
  §6 Risiken, §8 Vorprüfungen
- das CR-Dokument `docs/plan/cr/2026-09-07-cr-eingehend-adopter-links-anchors-exempt-paths.md`
  in voller Länge (Wortlaut des Absenders, Messung 2026-09-07, §Messung, §Entscheid)
- `DC-FA-REF-001`, `DC-FA-LINK-001`, `DC-FA-ANCH-001`
- `MR-069` (das Ventil als deklarierte Gate-Senkung)
- `AGENTS.md` §3.1, §3.7, §3.8, §5, §6
- Reviewer-Anker 8 (Reichweite einer Botschaft), 9 (Geltungsbereich eines
  Zitats), 17 (Proxy-Messung), 18 (Grenzen-Liste)
- Vorgänger-Entscheid: der Zeilenlängen-CR (entschieden 2026-09-07, nicht
  umgesetzt) — die von DoD (3) benannte Form-Vorlage

**Eigene Läufe** (echte Ausgaben, gekürzt auf die tragenden Zeilen):

- `make gates` — grün:
  `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`
- `make doc-check` — `d-check: 743 Datei(en) geprüft, 0 Befund(e)`
- **eigener Mini-Bestand** unter dem Scratchpad (nicht der des Slice), gefahren
  über das lokal gebaute Image mit `--config` je Konfiguration. Bestand: eine
  eingefrorene Datei mit sechs Referenz-Formen (toter Verweis *mit* und *ohne*
  Baseline-Bezug, Verweis über einen **Symlink**, Verweis, der die
  **Repo-Wurzel verlässt**, Anker in eine **vorhandene** Datei, Anker auf die
  **eigene** Datei), dazu eine zweite eingefrorene Datei mit einem Anker in
  einen **entfernten** Baum und ein lebendes Dokument.

  | Konfiguration | Befunde | was in der eingefrorenen Datei übrig bleibt |
  |---|---|---|
  | **C** — kein Ventil | 8 | alle sechs |
  | **A** — `in: frozen/**` · `refs: ["**"]` | 3 | **`symlink`** |
  | **B** — `in: frozen/**` · `refs: ["gone/**"]` | 7 | alles außer dem Baseline-Verweis |
  | **D** — `refs: ["../**", "symtarget.md"]` | 5 | **`symlink`** (`repo-escape` **weg**) |

---

## Findings

### F-1 — „A ist `exempt-paths`" ist gemessen falsch: die Klasse `symlink` überlebt jedes `refs`-Glob

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §5 (*„ihr Schluss reicht nicht weiter als die gemessene
  Menge"*) · Reviewer-Anker 8 und 17 · `DC-FA-REF-001`
- `pfad`: `docs/plan/cr/2026-09-07-cr-eingehend-adopter-links-anchors-exempt-paths.md:239-241`
  (*„**A ist `exempt-paths`.** Ein `refs: ["**"]` mit `in:` nimmt die genannte
  Datei **vollständig** aus der Prüfung"*), `:264-265` (*„erreicht dessen
  Wirkung **exakt**"*), `:277` (*„**kann**, was die fehlende könnte — und
  mehr"*); Commit `a33b8547` (*„Das IST exempt-paths"*, *„und mehr"*)
- `befund`: Der Entscheid vergleicht ein **gemessenes** Verhalten von
  `ignore-refs` mit einem **beschriebenen** Verhalten von `exempt-paths` und
  setzt beide gleich; die zweite Hälfte wurde nie gefahren. Gemessen bleibt in
  Konfiguration **A** (`refs: ["**"]`, also die behauptete Voll-Ausnahme) der
  Befund `frozen/report.md:5 ../symtarget.md symlink Linkziel ist oder enthält
  einen Symlink` stehen, und auch ein gezielt auf den Symlink gerichtetes
  `refs` (Konfiguration **D**) unterdrückt ihn nicht — die Symlink-Prüfung in
  `links` läuft vor dem Ventil. Ein `exempt-paths` nach der Semantik der
  übrigen Module (*„Dateien ganz ohne `<modul>`-Prüfung"*) nähme ihn mit. Die
  beiden Ventile stehen damit nicht in der behaupteten Dominanz, sondern quer:
  `ignore-refs` ist feiner auf der Ziel-Achse, `exempt-paths` breiter auf der
  Befund-Klassen-Achse.
- `verifizierbar`: ja —
  `docker run --rm --network none -v <mini>:/repo:ro d-check:latest --config cfgA.yml`
  über einen Bestand mit einem Link auf einen Symlink; der Befund `symlink`
  steht in der Ausgabe.
- `klasse`: `ventil-gleichgesetzt-ohne-die-befundklassen-zu-messen`

### F-2 — Der Entscheid deckt zwei Module, die Messung deckte eines; auf der `anchors`-Achse ist das Rezept für den beschriebenen Fall leer

- `kategorie`: MEDIUM
- `quelle`: `DC-FA-ANCH-001` · `AGENTS.md` §5 · Reviewer-Anker 8
  (*„suche die N+1-te Form"*) — derselbe Ableiter, den §8 des Plans selbst als
  *„für DoD (1) der wichtigste"* benennt
- `pfad`: `docs/plan/cr/…-links-anchors-exempt-paths.md:226-245` (§Messung,
  Tabelle mit 3/1/2 Befunden) und `:310` (*„Das deckt `links` **und**
  `anchors`"*)
- `befund`: Die drei Konfigurationen der §Messung ergeben 3/1/2 Befunde über
  „zwei tote Verweise plus ein lebendes Dokument" — die Zahlen gehen ohne Rest
  in `links` auf; `anchors` kommt im Bestand nicht vor und wurde nicht
  gefahren, obwohl der Entscheid über beide Module fällt. Gemessen kommt
  hinzu, dass das ausgeschriebene Rezept auf der `anchors`-Achse für genau den
  Fall des CR **nichts** bewirkt: ein Anker in einen **entfernten** Baum
  erzeugt nie `anchor-missing` (die Zieldatei fehlt, `anchors` ist nicht
  zuständig — gemessen: `frozen/anker.md:3 …#irgendein-anker target-missing`,
  kein zweiter Befund), und wo `anchors` wirklich meldet, matcht ein `refs`
  auf den entfernten Baum nicht — in Konfiguration **B** bleiben beide
  `anchor-missing`-Befunde der eingefrorenen Datei stehen. Die für `anchors`
  wirksamen `refs` müssten die **lebenden** Zieldateien nennen; die im Nein
  hervorgehobene Trennschärfe ist damit eine `links`-Eigenschaft.
- `verifizierbar`: ja — derselbe Lauf mit `--config cfgB.yml`; die zwei
  `anchor-missing`-Zeilen der eingefrorenen Datei stehen in der Ausgabe.
- `klasse`: `entscheid-breiter-als-die-gemessene-modulmenge`

### F-3 — „Die sechs Module stimmen" bestätigt die Grep-Form des Absenders, nicht die Fähigkeit — es sind acht, und der Absatz widerspricht sich selbst

- `kategorie`: MEDIUM
- `quelle`: Reviewer-Anker 17 (Proxy statt Gegenstand) · `AGENTS.md` §5
- `pfad`: `docs/plan/cr/…-links-anchors-exempt-paths.md:120-125`
  (*„**Die sechs Module stimmen.** Nachgezählt über …"*) und `:274`
  (*„`codepaths` führt drei (Zeile · Datei · Ziel), `ids` zwei"*); Slice-Plan
  slice-216 §3 (*„die sechs Module mit `exempt-paths` sind nachgezählt"*) und
  §2 DoD (2) (*„die sechs Module"* als Name des Arguments)
- `befund`: Die Nachzählung reproduziert die Schreibform `<modul>.exempt-paths`
  aus dem Grep des Absenders statt der Frage *„welches Modul führt das
  Ventil"*. Die Konfigurations-Referenz führt acht: `matrix`, `codepaths`,
  `diagrams`, `versions` (Kurz- **und** Paar-Form), `workflows`, `reviews`
  **sowie** `ids.patterns[].exempt-paths` und `structure[].exempt-paths` — die
  beiden letzten fallen aus dem Muster, weil vor dem Punkt eine eckige Klammer
  steht. Der Fehltreffer wird im selben Abschnitt sichtbar: die Antwort auf (b)
  schreibt `ids` **zwei** Ventil-Achsen zu, und die zweite ist genau
  `exempt-paths` — also führt `ids` das Ventil und steht trotzdem nicht in den
  sechs.
- `verifizierbar`: ja — die Konfigurations-Referenz in `spec/spezifikation.md`
  nach Tabellenzeilen mit `exempt-paths` durchsuchen; sie listet die acht
  Schlüssel. Die Gegenprobe steht in `internal/adapter/driven/configyaml/configyaml.go`.
- `klasse`: `zaehlung-uebernimmt-die-form-der-fremden-messung`

### F-4 — „dieses Repo fährt 25 Einträge ohne Konflikt" misst die Abwesenheit der Regel, nicht die Breite der Einträge

- `kategorie`: MEDIUM
- `quelle`: Reviewer-Anker 17 · Beobachtungs-Register
  `BEO-ALL/eigene-menge-gemessen-fremde-behauptet` (16×), das §8 des Plans
  selbst für dieses Argument benennt
- `pfad`: `docs/plan/cr/…-links-anchors-exempt-paths.md:294-296`
- `befund`: Der Satz steht als Gegengewicht zum Breiten-Wächter des Absenders,
  unmittelbar nach der Zusage, dessen Kappung sei *„nicht beurteilbar von
  hier"*. Dieses Repo führt aber **keinen** Breiten-Wächter: der
  Zeilenlängen-CR ist am 2026-09-07 — einen Tag vor diesem Entscheid — mit
  *„nicht umgesetzt"* beschieden worden, und weder `.d-check.yml` noch
  `Makefile` noch `.golangci.yml` tragen eine Zeichen-je-Zeile-Bedingung.
  „Ohne Konflikt" ist damit von der Breite der 25 Einträge unabhängig und
  liest sich als Vergleich mit dem fremden Repo, den der Satz davor ausschließt.
  Hinzu kommt, dass die Zahl die falsche Menge ist: die Liste trägt **45**
  Einträge, 25 ist die Baseline-Teilmenge aus `MR-069`.
- `verifizierbar`: ja — `grep -cE '^  - in:' .d-check.yml` ⇒ 45;
  `grep -rn "MD013\|line-length\|max-line" .d-check.yml Makefile .golangci.yml`
  ⇒ leer; Kopf des Zeilenlängen-CR ⇒ *„entschieden am 2026-09-07 — nicht
  umgesetzt"*.
- `klasse`: `vergleich-ohne-gemeinsame-messgroesse`

### F-5 — Umkehr-Bedingung (1) war beim Schreiben bereits erfüllt, und die Liste nennt ihre größte Lücke nicht

- `kategorie`: MEDIUM
- `quelle`: Reviewer-Anker 18 (Grenzen-Liste ohne ihre größte Lücke) ·
  Beobachtungs-Register `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`
  (6×, verkörpert), von §8 des Plans selbst für diesen Abschnitt benannt
- `pfad`: `docs/plan/cr/…-links-anchors-exempt-paths.md:317-330`
  (§Was den Entscheid umkehren würde, Bedingung 1 und die Schluss-Grenze)
- `befund`: Bedingung (1) lautet *„Ein Fall, den `ignore-refs` nicht ausdrücken
  kann. Ein gemessenes Beispiel genügt"* — F-1 ist ein solches Beispiel und
  stand zum Zeitpunkt des Entscheids schon fest, ungemessen. Eine
  Umkehr-Bedingung, die bei ihrer Niederschrift bereits zutrifft, löst keine
  Wiedervorlage aus, weil sie als offen gelesen wird. Als benannte Grenze führt
  der Abschnitt nur *„Kein Sensor wacht über die drei"*; die schwerere Lücke —
  Bedingung (1) nennt, anders als die Vorlage des Zeilenlängen-CR, **keine
  Population**, gegen die sie nachmessbar wäre (*„Nachmessbar mit derselben
  Population wie oben — ohne sie ist die Bedingung nicht prüfbar"*) — steht
  nicht da. Gemeinsame Wurzel mit F-1: dieselbe unterlassene Klassen-Messung.
- `verifizierbar`: ja — derselbe Lauf wie F-1 belegt den Fall, den Bedingung (1)
  verlangt.
- `klasse`: `umkehr-bedingung-bei-niederschrift-bereits-erfuellt`

### F-6 — „die schwächere ist die bequemere" ist behauptet, und die eigene Konfiguration ist der ungenannte Gegenfall

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §5 · Reviewer-Anker 8
- `pfad`: `docs/plan/cr/…-links-anchors-exempt-paths.md:278-281`
- `befund`: Der Satz *„Wer `exempt-paths` schreiben kann, schreibt es, und
  verliert die Unterscheidung"* ist eine Verhaltensaussage über künftige
  Konfigurations-Autoren; gemessen ist sie nicht, und sie trägt die Antwort auf
  das Argument, das der Entscheid zuvor ausdrücklich zugesteht. Der nächste
  Gegenfall steht in der eigenen Datei: `codepaths` führt **alle drei** Achsen
  (Zeile · Datei · Ziel), und `.d-check.yml:522` nimmt `docs/reviews/**` über
  die **Datei**-Achse aus — für dieselbe Artefakt-Klasse, um die der CR bittet;
  der Kommentar bei `.d-check.yml:76` beruft sich sogar auf diese Wahl als
  Präzedenz. Der Entscheid nennt den Fall nicht und trennt ihn nicht vom
  Fall des Absenders (dort genügt ein Glob, hier müsste man aufzählen) — er
  verallgemeinert stattdessen.
- `verifizierbar`: ja — `grep -n 'exempt-paths' .d-check.yml` (Zeilen 76, 522)
  gegen `DC-FA-CODE-001` §Ventil-Achsen.
- `klasse`: `argument-gegen-den-absender-ohne-die-eigene-praxis-zu-nennen`

### F-7 — `DC-FA-REF-001` und der Entscheid widersprechen sich über die unterdrückten Klassen; die Messung gibt keinem von beiden recht

- `kategorie`: MEDIUM
- `quelle`: `DC-FA-REF-001` · Reviewer-Anker 10 (Messmethode klafft gegen die
  Spec-Stelle) und Anker 9 (Quelle für das gelesen, was sie gewährt, nicht für
  das, was sie begrenzt)
- `pfad`: `spec/lastenheft.md:981-986` (*„das Ventil unterdrückt **nur** die
  Auflösungs-Klasse … keine anderen Befunde (Symlink-Ablehnung, Repo-Escape
  u. Ä. bleiben)"*) gegen `internal/hexagon/core/rules/links.go:40-49` und
  `docs/plan/cr/…-links-anchors-exempt-paths.md:239-241`
- `befund`: Die Anforderung sagt, Symlink-Ablehnung **und** Repo-Escape blieben
  hinter dem Ventil stehen; der Entscheid sagt, die Datei falle
  *„vollständig"* aus der Prüfung. Gemessen trifft keines von beidem: `symlink`
  bleibt (F-1), `repo-escape` nicht — in Konfiguration **A** wie in **D**
  verschwindet `frozen/report.md:6 ../../ausserhalb.md repo-escape` aus der
  Ausgabe, weil `refIgnored` im Modul `links` **vor** dem Escape-Zweig läuft.
  Die Divergenz ist älter als dieser Slice und nicht von ihm verursacht; sie
  ist aber der Rahmen, aus dem der Entscheid seine Klassen-Aussage nimmt, und
  keine der beiden Hälften wurde gefahren. Versagensbild: wer der Anforderung
  glaubt und für ein eingefrorenes Verzeichnis ein breites `refs` schreibt,
  verliert die Escape-Meldung still mit.
- `verifizierbar`: ja — `--config cfgD.yml` mit `refs: ["../**", "symtarget.md"]`:
  `repo-escape` fehlt in der Ausgabe, `symlink` steht darin.
- `klasse`: `spec-satz-ueber-nicht-unterdrueckte-klassen-vom-produkt-widerlegt`

### F-8 — „Was die Messung dazu beiträgt" schreibt der Messung eine Aussage zu, die sie nicht erzeugt hat

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §5 · Reviewer-Anker 8
- `pfad`: `docs/plan/cr/…-links-anchors-exempt-paths.md:296-298`
- `befund`: Die Messung hat Befund-**Anzahlen** über drei Konfigurationen
  erzeugt; die Aussage *„Form **B** ist eine Zeile `refs:` mit einem Glob —
  deutlich kürzer als ein Paar je Ziel"* ist eine Betrachtung der
  Konfigurations-Form und stammt nicht aus ihr. Der Schluss selbst ist
  plausibel und ausdrücklich als nachzumessen markiert; falsch ist nur die
  Herkunfts-Zuschreibung.
- `verifizierbar`: nein — Textbefund; kein Gate-Lauf entscheidet ihn.
- `klasse`: `aussage-der-messung-zugeschrieben-die-sie-nicht-traegt`

### F-9 — „die sieben toten `in:`-Skopen" übernimmt eine teilmengen-skopierte Zahl aus `MR-069` als Gesamtaussage

- `kategorie`: LOW
- `quelle`: Reviewer-Anker 9 (Geltungsbereich lesen, nicht den Titel) ·
  `MR-069` §Grenze (2)
- `pfad`: Slice-Plan slice-216 §1, dritter Abgrenzungs-Punkt
- `befund`: `MR-069` zählt *„sieben der **25 Baseline-Einträge**"* mit totem
  `in:`-Skopus. Der Abgrenzungs-Punkt schreibt *„die 25 Baseline-Einträge und
  **die** sieben toten `in:`-Skopen"* und macht die Teilmengen-Zahl mit dem
  bestimmten Artikel zur Gesamtaussage. Gemessen über die ganze Liste: 45
  Einträge, alle mit `in:`, davon 18 Globs, 11 auf existierende und **16** auf
  nicht mehr existierende literale Pfade. Der Punkt grenzt den Bestand
  ausdrücklich aus, verkleinert ihn aber im Vorbeigehen.
- `verifizierbar`: ja — `grep -E '^  - in:' .d-check.yml` und je Wert `test -e`.
- `klasse`: `teilmengen-zahl-als-gesamtaussage-uebernommen`

### F-10 — Auf der `anchors`-Achse fällt die Ziel-Achse für Eigen-Datei-Anker mit der Datei-Achse zusammen

- `kategorie`: INFO
- `quelle`: `DC-FA-ANCH-001` · `DC-FA-REF-001`
- `pfad`: `internal/hexagon/core/rules/anchors.go:249` (`refIgnored(ignoreRefs, file, a.rel)`)
  in Verbindung mit `:223-224` (`idx == 0` ⇒ `out.rel = file`)
- `befund`: Für einen Anker auf die **eigene** Datei (`#abschnitt`) ist das
  Ventil-Ziel die Quelldatei selbst; ein `refs`-Glob, das ihn stumm schaltet,
  ist ein Glob auf die Quelldatei — also exakt die Granularität, die der
  Entscheid `exempt-paths` als gröber vorhält. Gemessen: in Konfiguration **B**
  meldet `frozen/report.md:8 #gibt-es-auch-nicht anchor-missing` weiter, in
  **A** (`refs: ["**"]`) nicht. Dokumentationswürdige Annahme, kein Defekt —
  sie begrenzt die Reichweite des Auflösungs-Arguments auf eine Klasse, die der
  Entscheid nicht ausnimmt.
- `verifizierbar`: ja — derselbe Lauf mit `cfgA.yml`/`cfgB.yml`.
- `klasse`: `ziel-achse-kollabiert-auf-die-quelldatei`

---

## Negativbefunde

- **geprüft, ohne Befund: die tragende Messung reproduziert.** Über einen
  unabhängig gebauten Mini-Bestand ergeben C/A/B qualitativ dasselbe Bild wie
  im CR: ohne Ventil melden beide toten Verweise der eingefrorenen Datei, mit
  `refs: ["**"]` schweigen beide, mit einem auf den entfernten Baum skopierten
  `refs` schweigt nur der Baseline-Verweis und der zweite meldet weiter. Der
  **Kern** des Entscheids — die Ziel-Achse kann den beschriebenen Fall enger
  ausdrücken als eine Datei-Achse — trägt.
- **geprüft, ohne Befund: keine Sammelantwort.** Die drei Argumente bekommen
  drei getrennte Antworten; (b) wird ausdrücklich zugestanden, bevor es
  beantwortet wird, und nicht mit der widerlegten Prämisse (a) erledigt — genau
  das Risiko, das §6 als erstes führt.
- **geprüft, ohne Befund: die Korrektur an (a) ist präzise.** *„keine
  modul-lokale Options-Sektion — der Knopf sitzt eine Ebene höher"* trifft den
  Sachverhalt und nimmt dem Absender nur die Prämisse, nicht die Beobachtung.
- **geprüft, ohne Befund: (c) konzediert, was zu konzedieren ist.**
  `scan.ignore` als falsches Instrument und die Koexistenz als bloßer Aufschub
  werden ohne Einschränkung zugestimmt; das Rezept beantwortet die Bitte mit
  einer ausgeschriebenen Konfiguration statt mit einem Doku-Verweis, und für
  die `links`-Hälfte des Anliegens trägt es.
- **geprüft, ohne Befund: der Breiten-Wächter des Absenders wird nicht
  bewertet.** *„Ob diese Kappung eine bewusste Regel ist, weiß nur der
  Absender"* und *„wäre neu zu messen"* halten die Zurückhaltung durch; die
  offene Rückfrage ist im Schluss-Absatz **benannt**, nicht nur erwähnt, und
  der Entscheid sagt ausdrücklich, dass er ohne sie gefallen ist. Der Mangel
  sitzt allein im Vergleichs-Satz (F-4).
- **geprüft, ohne Befund: §1-Abgrenzung eingehalten.** Der Range fasst genau
  drei Dateien an (CR-Dokument, Slice-Plan, Roadmap-Ruhe-Marker); keine
  Umsetzung, kein Eingriff in `.d-check.yml` oder den Ventil-Bestand, keine
  versandte Antwort. Das Rezept im Entscheid ist Inhalt der Entscheidung, kein
  Anschreiben.
- **geprüft, ohne Befund: §8 Vorprüfungen.** Beide kanonischen Blöcke tragen
  eine `d-check:cite`-Direktive mit wörtlichem Zitat, der dritte
  (Nachtlauf, `MR-053`) bewusst keine. Nachgezählt: **38**
  Beobachtungs-Verzeichnisse (37 unter `BEO-ALL`, 1 unter `BEO-HARN`); der
  `BEO-HARN`-Eintrag ist der `--check-latest`-Eintrag und berührt den Slice
  nicht. Die vier zitierten Einträge tragen die genannten Stände: 11× /
  18× / 16× (*gemischt*) / 6×.
- **geprüft, ohne Befund: §8 Sub-Area-Wahl.** `harness/conventions.md`
  §Modus-Deklaration führt genau zwei Zeilen (`*` und `tools/harness/`); die
  Begründung, für Produkt-Code sei keine eigene deklariert und der Slice fahre
  das Produkt nur, statt es zu ändern, deckt sich mit dem Diff und erfindet
  keine dritte Sub-Area.
- **geprüft, ohne Befund: `Stand:`-Zeile als Zustandsfeld (`AGENTS.md` §3.7).**
  Sie nennt Zustand (*entschieden … nicht umgesetzt*) und Beleg als auflösbaren
  Anker (§Entscheid, §Messung) und erzählt keine Chronik; die Historie steht in
  den datierten Abschnitten darunter, wo sie hingehört.
- **geprüft, ohne Befund: Umkehr-Bedingungen vorhanden und beobachtbar
  formuliert**, mit der Grenze *„Kein Sensor wacht über die drei"* — die von
  DoD (3) verlangte Form ist da; der Einwand aus F-5 betrifft ihren Inhalt,
  nicht ihr Vorhandensein.
- **geprüft, ohne Befund: Commit-Botschaft `215585fb`** (Beanspruchung) —
  *„Gefahren: `make planning-check` — 743 Dateien, 0 Befunde"* deckt sich mit
  dem hier nachgefahrenen Lauf; der Move ist rein, kein Pfad-Verweis war
  nachzuziehen.
- **geprüft, ohne Befund: Gate-Stand.** `make gates` grün (zehn Gates),
  `make doc-check` `743 Datei(en) geprüft, 0 Befund(e)` — beide Zahlen der
  Botschaft `a33b8547` treffen zu.
- **geprüft, ohne Befund: `AGENTS.md` §3.1.** Kein Host-Go, kein
  Host-Interpreter in diesem Lauf; ein versehentlicher `python3`-Aufruf wurde
  vom Tool-Call-Wächter blockiert und nicht umgangen.
- **geprüft, ohne Befund: Referenz-Richtung und Provenance-Marker.** Weder
  Slice-Plan noch CR-Dokument tragen einen `d-check:status-provenance`-Marker;
  `matrix` führt `docs/plan/cr/` und `docs/reviews/` in keiner Klasse, die
  Abwärts-Sperre ist nicht berührt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 7 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:**
`ventil-gleichgesetzt-ohne-die-befundklassen-zu-messen` ·
`entscheid-breiter-als-die-gemessene-modulmenge` ·
`zaehlung-uebernimmt-die-form-der-fremden-messung` ·
`vergleich-ohne-gemeinsame-messgroesse` ·
`umkehr-bedingung-bei-niederschrift-bereits-erfuellt` ·
`argument-gegen-den-absender-ohne-die-eigene-praxis-zu-nennen` ·
`spec-satz-ueber-nicht-unterdrueckte-klassen-vom-produkt-widerlegt` ·
`aussage-der-messung-zugeschrieben-die-sie-nicht-traegt` ·
`teilmengen-zahl-als-gesamtaussage-uebernommen` ·
`ziel-achse-kollabiert-auf-die-quelldatei`

## Verdikt

**Merge-blockierend:** ja — sieben MEDIUM.

**Der Entscheid selbst kippt dadurch nicht.** Die tragende Messung reproduziert
auf einem unabhängigen Bestand, und für den Fall, den der CR beschreibt, ist
die Ziel-Achse tatsächlich das schärfere Instrument. Was nicht trägt, ist die
**Reichweite**, in der der Entscheid das sagt: „exakt", „vollständig", „und
mehr" und „das IST `exempt-paths`" sind Aussagen über **alle** Befund-Klassen
und **beide** Module, gemessen wurde eine Klasse in einem Modul. Gegen die
Bitte des Absenders steht damit ein Nein, dessen Begründung an drei Stellen
weiter reicht als sein Beleg — und eine seiner eigenen Umkehr-Bedingungen ist
bereits erfüllt, ohne dass es jemand bemerkt hätte.

**Übergabe:** Findings an den Implementer (Rückkante Review → Plan, weil F-1,
F-2 und F-5 die Messung von DoD (1) und die Entscheid-Substanz betreffen, nicht
nur ihre Formulierung). F-7 betrifft `DC-FA-REF-001` und `links` und ist älter
als dieser Slice — er gehört als eigener Vorgang benannt, nicht hier
mitgenommen. Die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7
und von dort in den Zähler; `commit-message-overclaims-work`,
`zaehlmethode-misst-proxy-statt-gegenstand`,
`eigene-menge-gemessen-fremde-behauptet` und
`grenzen-liste-wird-als-vollstaendig-gelesen` sind in diesem Lauf jeweils
erneut belegt. Dieser Report ersetzt keine Verifikation — DoD-Konformität
prüft der Verifier separat.
