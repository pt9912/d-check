# Slice slice-183: Die Baseline steht auf v5.15.0

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Sein Closure-Grund geht über die eigene DoD nicht
hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht); der Anlass ist ein **Upstream-Release**, keine Welle.

**Bezug:**
[`MR-021`](../../../../harness/conventions.md#mr-021) (In-Repo-Verweise sind
pin-gebunden — die Pfad-Hälfte des Bumps);
[`MR-051`](../../../../harness/conventions.md#mr-051) (die `cite`-Spannen, die
zweite pin-gebundene Größe, samt den drei Fallunterscheidungen);
[`MR-055`](../../../../harness/conventions.md#mr-055) (die Symlink-Aliase, die
denselben Pin binden und von keinem Modul gescannt werden);
[`MR-039`](../../../../harness/conventions.md#mr-039) (ein zitierter Wortlaut
wird nicht rückwirkend umgeschrieben — der Ausgang für ein echtes Delta).

**Berührte Spec-Stellen:** — (Adoptions-Stand einer externen Konvention; keine
Produkt-Anforderung berührt).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-08-30.

---

## 1. Ziel

**Der vendorte Baseline-Baum steht auf `v5.12.0`, upstream stehen vier
Releases weiter.** Gemessen mit `make baseline-freshness`: `v5.13.0`,
`v5.13.1`, `v5.14.0`, `v5.15.0` — und der Content am gepinnten Tag ist **unverändert**
(`Bytes == vendored SHA256SUMS`). Es ist also ein reiner **Currency**-Rückstand,
kein Drift: nichts ist falsch, aber die verkörperte Form beruft sich auf eine
Fassung, die vier Releases alt ist.

**Der Herausgeber hat den Inhalt für `v5.13.0` angekündigt** — fünf
Regelwerk-Dateien mit echtem Inhalt (`grundlagen-durchsetzungsschicht`,
`grundlagen-begriffe`, `modul-13`, `modul-14`, `modul-02`), dazu zwei Punkte,
die uns direkt betreffen: der Abschnitt `§Referenz-Implementierung` heißt jetzt
`§Das vollständige Artefakt-Set`, und die `MR`-Vorlage hat ihr `Status:`-Feld
verloren.

**Das Delta ist inzwischen gelesen, nicht angenommen** — netzlos aus dem
Kurs-Klon, Tag für Tag: `v5.13.0` = Wellen 99–101, `v5.13.1` = Welle 102
(kein Delta in `kurs/de/`, geändert ist nur die Konfigurations-Vorlage),
`v5.14.0` = Welle 103, `v5.15.0` = Wellen 105/106. Von 52 Bundle-Dateien sind
20 unverändert, 16 tragen nur den Versions-Stempel, 14 echten Regel-Inhalt —
darunter `templates/.d-check.yml` mit 31 geänderten Zeilen, die eine
Delta-Schleife über `*.md` allein übersehen hätte ([`MR-037`](../../../../harness/conventions.md#mr-037)
hatte genau davor gewarnt).

**Das Ziel ist während der Wartezeit gewandert, und der Grund ist neu.** Der
Plan zielte bei seiner Anlage auf `v5.14.0`. Inzwischen ist `v5.15.0`
erschienen, und **erst dort** liegen die zwei Quell-Wellen, die unser
Beobachtungs-Register betreffen: die Vergabe des Bereichssegments und die
Schwelle mit ihren drei Ausgängen. Nachgemessen: beide Wellen sind **sieben
Commits jünger als `v5.14.0`** und in diesem Tag nicht enthalten. Ein Bump
dorthin hätte den Pin gehoben und den eigentlichen Anlass verfehlt.

**Damit hängt ein zweiter Slice an diesem hier:**
[slice-188](../done/slice-188-register-gegen-neuen-kanon.md) korrigiert zwei Zähler
unseres Registers gegen die neue Beleg-Definition — und kann die Regel erst
zitieren, wenn sie im vendorten Baum steht. Diese Abhängigkeit ist neu und war
bei der Anlage nicht absehbar.

## 2. Vorgehen

1. **Re-vendorn und das Delta sichten.** `fetch-baseline-cache.sh v5.15.0`,
   dann `make baseline-verify` (Integrität **und** Manifest-Deckung **und**
   Alias-Auflösung). Erst danach steht fest, was die Schritte 3 und 4 zu tun
   haben — der Plan nimmt es nicht vorweg.
2. **Alle Pfad-Verweise ziehen** ([`MR-021`](../../../../harness/conventions.md#mr-021)):
   gemessen bei der Anlage **76** Dateien mit `baseline/v5.12.0`, dazu die
   Symlinks unter `.claude/rules/`
   ([`MR-055`](../../../../harness/conventions.md#mr-055)) — sie binden
   denselben Pin, werden von keinem Modul gescannt und brächen still.
   `make baseline-verify` fängt genau das. **Bei der Ausführung nachgemessen:
   es sind VIER Aliase, nicht zwei**, und der Wächter hat alle vier gemeldet,
   bevor irgendetwas committet war.
3. **Die `cite`-Spannen neu ankern**
   ([`MR-051`](../../../../harness/conventions.md#mr-051)). **Vorab-Schätzung
   bei der Anlage** (noch auf `v5.14.0` gezielt, §1): rund 33 echte Direktiven
   in den vendorten Baum — die rohe Zahl 133 zählt Erwähnungen mit, darunter
   Doku-Beispiele wie `<pfad>`. Geschätzte Verteilung: 25 in `modul-05`, je 2 in
   `modul-09` und `grundlagen-durchsetzungsschicht`, je 1 in `modul-02`,
   `modul-10`, `modul-11` und den Templates. **Diese Schätzung war überholt:**
   Bei der Ausführung gegen `v5.15.0` gemessen sind es **24** lebende
   Direktiven, davon **14** in `modul-05` — das deckt sich mit dem
   tatsächlichen Ergebnis (§4: 2 nachgezogen + 0 entfernt + 22 unverändert =
   24). Die Differenz zur Vorab-Schätzung erklärt sich aus dem gewechselten
   Ziel-Tag (§1) und aus der groben Ersteinschätzung; sie ändert nichts an
   Schritt 4. **Unterschieden wird
   nach Grund-Code, nicht nach Gefühl:** `citation-mismatch` mit dem Wortlaut
   anderswo in derselben Datei ⇒ Spanne nachziehen; Wortlaut nirgends mehr ⇒
   [`MR-039`](../../../../harness/conventions.md#mr-039), Direktive entfernen,
   Delta im Bump-Eintrag festhalten; `citation-out-of-range` ⇒ **zuerst** prüfen,
   ob der Abschnitt in eine andere Datei gewandert ist.
4. **Der Adaptions-Review durch die Liste, nicht durch den Diff.** Der Kanon
   verlangt ihn ausdrücklich, und zwar je Eintrag mit einem von **fünf**
   Ausgängen (gegenstandslos · bleibt gültig · teilweise überholt · Bezug
   entfallen · Schärfung). Gegenstand sind die **33 lebenden** `MR`-Einträge;
   die 23 in `conventions/done/` sind aufgelöst und nicht Gegenstand.
   **Ein Ausgang ist schon bekannt:** drei lebende Einträge zitieren
   `§Referenz-Implementierung`; der Abschnitt heißt ab `v5.13.0`
   `§Das vollständige Artefakt-Set`. **Bei der Ausführung präzisiert:**
   [`MR-048`](../../../../harness/conventions.md#mr-048) ist der eigenständige
   Fund dieses Adaptions-Reviews — sein Linktext war seit seiner Anlage
   korrekt und wurde erst durch die Umbenennung stale.
   [`MR-042`](../../../../harness/conventions.md#mr-042) und
   [`MR-043`](../../../../harness/conventions.md#mr-043) tragen denselben
   mechanischen Nachzug, aber ihre Situation war schon bekannt: der
   `F-7`-Befund aus dem slice-155-Review hatte genau ihren Linktext gemeldet
   (in der Gegenrichtung — sie nannten den neuen Namen, als noch der alte
   galt) und ist damit jetzt geschlossen, nicht neu gefunden (siehe
   [`MR-057`](../../../../harness/conventions.md#mr-057)).
5. **Die Bestands-Stichprobe fahren**, die
   [`AGENTS.md`](../../../../AGENTS.md) §1 auch bei aktuellem Pin verlangt —
   sie hängt nicht am Delta.
6. `make gates` und `make fullbuild`; **Review** und **Verifikation** als
   getrennte Läufe; Closure.

## 2a. Ergebnis der Bestands-Stichprobe (ausgeführt 2026-08-31)

**Gezogen:** [`modul-07-carveouts.md` §Ziel-Form: Carveout](../../../../.harness/baseline/v5.15.0/regelwerk/modul-07-carveouts.md#ziel-form-carveout).
**Auswahl nach Kanon-Kriterium** — aus den inhaltlich delta-freien Abschnitten
(15 Module tragen nur den Versions-Stempel), rotierend: `modul-14` war die
Stichprobe des `v5.6.0`-Bumps. Der Grund für gerade diesen Abschnitt ist
[`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage): dort
steht *„keine inhaltlichen Adaptionen ggü. Baseline-Default für …
Carveout-Disziplin"* — eine Behauptung, die dieses Repo nie geprüft hat, weil
es **nie einen Carveout eröffnet** hat.

**Befund: drei der vier operativen Regeln stehen nicht im ausgefüllten
Artefakt.** [`docs/plan/carveouts/README.md`](../../carveouts/README.md) ist
die einzige Stelle, die die Disziplin trägt:

| Kanon-Regel | Im ausgefüllten Artefakt? |
|---|---|
| Sechs Pflicht-Header-Felder (Status · Datum angelegt · Letzte Prüfung · betroffenes Gate · Geltungsbereich · Folge-Slice) | **nein** — genannt sind zwei (Auflösungs-Trigger, Folge-Slice) |
| Auflösungs-Trigger als beobachtbare, messbare Bedingung | teilweise — das Feld ist genannt, die Beobachtbarkeits-Anforderung nicht |
| Die Gate-Konfiguration nennt die `CO-<NNN>` im Gate-Output | **nein** — nirgends |
| Auflösung ist ein `git mv` nach `done/` | **nein** — es gibt kein `carveouts/done/` |

**Das ist genau die Klasse, die die Stichprobe finden soll:** eine alte,
stabile Baseline-Regel erzeugt keinen Template-Diff und hat keinen
`MR`-Eintrag, den der Adaptions-Durchgang abschreiten könnte — sie ist
unsichtbar, *weil* sie sich nie geändert hat. **Kein Gate sieht das**, und die
Lücke ist folgenlos, solange kein Carveout existiert; sie wird es in dem
Moment, in dem einer eröffnet wird. **Ausgang: Folge-Slice** — entweder das
`README` auf die sechs Felder ziehen, oder die Abweichung als `MR` erklären.
Nicht in diesem Slice: die DoD nennt die Stichprobe nicht, nur das Vorgehen,
und ihr Zweck ist das Aufdecken.

## 2b. Zwei Funde des Delta-Durchgangs, die dieser Slice nicht einlöst

- **Der Kanon verlangt einen Versions-Sensor auf den Baseline-Pin, wir haben
  ihn nicht.** Kurs-Welle 103
  ([`grundlagen-traceability.md`](../../../../.harness/baseline/v5.15.0/regelwerk/grundlagen-traceability.md)):   <!-- d-check:cite .harness/baseline/v5.15.0/regelwerk/grundlagen-traceability.md:89-89 -->
  der adoptierte Stand steht **einmal** im Adaptions-Block, ein Sensor prüft
  jeden Pin dagegen — *„Ein vergessener Nachzug ist dann ein Befund, kein toter
  Link."* Bei uns fängt das heute die **Existenz** (`links`/`codepaths` melden
  den gelöschten Pfad), nicht der **Vergleich**; solange beide Bäume
  nebeneinander liegen, ist ein stehengebliebener Pin unsichtbar. Die
  Kurs-Vorlage schlägt dafür die **Kurzform** unter `versions:` vor — die ist
  hier von [`DC-FA-VER-001`](../../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in)
  belegt. Über `versions.patterns` ginge es ohne Produktänderung
  ([ADR-0058](../../adr/0058-konfigurations-flaechen-additiv-weiten.md)). **Der
  Vorschlag der Vorlage ist damit für jeden Adopter mit belegter Kurzform
  falsch** — das ist zugleich ein CR-Kandidat an den Kurs.
- **[`MR-013`](../../../../harness/conventions.md#mr-013) ist teilweise
  überholt.** Der Widerspruch samt beider Begründungen steht in
  [`MR-057`](../../../../harness/conventions.md#mr-057); die Auflösung ist ein
  Nachfolge-Eintrag plus Nachzug in [`AGENTS.md`](../../../../AGENTS.md) §3.3,
  nicht ein Edit am bestehenden Eintrag.

  **Der Herausgeber hat darauf geantwortet, und die Antwort ist nachgemessen.**
  Seine Feststellung: keine Stufe unserer CI checkt je einen Zwischen-Commit
  aus. Unabhängig geprüft in
  [`ci.yml`](../../../../.github/workflows/ci.yml) — `checkout` holt ohne
  `ref:` den Tip, `make trace-check` und `make adr-check` bekommen eine
  **Range** und lesen deren Enden (Zwei-Baum-Diff), `make ci` läuft über den
  ausgecheckten Baum, also wieder den Tip. **Die Begründung von [`MR-013`](../../../../harness/conventions.md#mr-013) ist
  damit zu breit formuliert:** der gefürchtete Fall entsteht in genau einer
  Lage — wenn jemand den Move **allein** pusht. Genau das war 2026-06-21 bei
  slice-040 passiert.

  **Zitierbar ist die neue Bedingung aber noch nicht.** Sie stammt aus
  Kurs-Welle 110 und liegt in **keinem Release**: `v5.16.0` existiert bereits,
  trägt Welle 109 und enthält sie nicht. Wer sie jetzt zitiert, zitiert etwas,
  das im vendorten Baum nicht steht — dieselbe Lage, aus der dieser Slice
  [slice-188](../done/slice-188-register-gegen-neuen-kanon.md) gerade befreit
  hat.

  **Und „Trigger statt Nachfolger" greift zu kurz.** [`MR-013`](../../../../harness/conventions.md#mr-013) deckt **drei**
  Move-Klassen; nur die dritte (`MR`-/Wellen-Lifecycle) kollidiert. Die beiden
  Slice-Klassen hängen an der Roadmap-Kopplung, die `make planning-check`
  erzwingt — ein anderer Mechanismus, der von der Bedingung unberührt bleibt.
  Verliert ein Eintrag durch die Baseline **einen Teil** seines
  Geltungsbereichs, ist das nach dem Kanon eine **Ablösung mit engerem
  Nachfolger**, keine Auflösung. Der Folge-Slice hat also beides zu tun: den
  engeren Nachfolger schreiben **und** den Pin auf ein Release heben, das
  Welle 110 trägt.

  **`v5.16.0` ist während dieser Arbeit erschienen** und damit Sache der
  nächsten Pin-Hebung, nicht dieser — ein Wechsel des Ziels hätte den
  gemessenen Delta-Durchgang entwertet, und Welle 109 (Wellen-Closure
  archiviert ihre Zeitdokumente) berührt den Anlass dieses Slice nicht.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Entfernen der `Status:`-Felder.** Die `MR`-Vorlage des Herausgebers
  hat das Feld verloren, und unsere eigene Regel
  ([`AGENTS.md`](../../../../AGENTS.md) §3.7: der Zustand ist die
  Verzeichnis-Position) spricht fürs Entfernen. Es sind **56** Dateien; zusammen
  mit den 76 des Bumps wäre der Diff nicht mehr in einer Sitzung prüfbar. Der
  Herausgeber sagt ausdrücklich, dass daraus nichts folgt — also eigener Slice.
- **Keine Antwort auf die `--print-mk`-Frage.** Der Herausgeber fragt, ob eine
  **mount-freie** Recipe-Form ins Werkzeug gehört; ein Adopter führt sie als
  Adaption mit **unserem** Werkzeug als Auflösungs-Trigger. Das ist eine
  Produkt-Entscheidung mit ADR und Lastenheft-Bump, kein Bump-Nebenprodukt.
  **Dieser Slice macht sie aber erst beurteilbar** — die Kanon-Stelle, auf die
  sie sich beruft (`modul-14 §Der Prüflauf ist hermetisch`), liegt in unserem
  Baum heute **nicht** vor.
- **Kein Nachziehen der Wortlaute in `done/`.** Ein zitierter Wortlaut wird
  nicht rückwirkend umgeschrieben
  ([`MR-039`](../../../../harness/conventions.md#mr-039)); [`MR-051`](../../../../harness/conventions.md#mr-051)s
  Geltungsbereich nennt die **lebenden** Dokumente.

## 4. Definition of Done

- [x] Der Pin steht auf `v5.15.0`: vendorter Baum re-vendored,
      `make baseline-verify` grün (Integrität, Manifest-Deckung, Alias-Auflösung),
      alle Pfad-Verweise und die zwei Symlinks gezogen. Belegt in `08373c9`:
      129 lebende `v5.12.0`-Vorkommen, 95 gehoben, 34 als Aussagen über die
      Vergangenheit stehen gelassen; alle vier `.claude/rules`-Aliase von
      `baseline-verify` vor dem Commit gemeldet.
- [x] **Jede `cite`-Direktive ist entschieden**, nicht nur grün: je Direktive
      steht fest, ob sie nachgezogen, umgehängt oder nach
      [`MR-039`](../../../../harness/conventions.md#mr-039) entfernt wurde. Die
      Zahl der drei Fälle steht in der Commit-Botschaft. Belegt in `08373c9`:
      2 nachgezogen ([`MR-035`](../../../../harness/conventions.md#mr-035),
      [`MR-043`](../../../../harness/conventions.md#mr-043)), 0 entfernt, 22
      gegen den Datei-Diff geprüft und unverändert bestätigt.
- [x] **Der Adaptions-Review ist gefahren und dokumentiert:** je lebendem
      `MR`-Eintrag einer der fünf Ausgänge. **Keine Treffer sind ebenfalls eine
      Antwort** und werden notiert. Belegt in `08373c9`/[`MR-057`](../../../../harness/conventions.md#mr-057):
      33 lebende Einträge geprüft, 31 bleiben gültig — darunter
      [`MR-042`](../../../../harness/conventions.md#mr-042) und
      [`MR-043`](../../../../harness/conventions.md#mr-043), deren
      Linktext-Nachzug bereits der `F-7`-Befund aus slice-155 trägt, keine
      neue Antwort dieses Reviews. **Zwei** tragen eine neue Antwort:
      [`MR-048`](../../../../harness/conventions.md#mr-048) (eigenständiger
      Fund: sein Linktext war seit Anlage korrekt und wurde erst durch die
      Umbenennung stale) und
      [`MR-013`](../../../../harness/conventions.md#mr-013), das teilweise
      überholt ist (§2b).
- [x] [`MR-048`](../../../../harness/conventions.md#mr-048),
      [`MR-042`](../../../../harness/conventions.md#mr-042) und
      [`MR-043`](../../../../harness/conventions.md#mr-043) zeigen auf
      `§Das vollständige Artefakt-Set`; der `F-7`-Befund aus slice-155 (er
      betraf die beiden zuletzt genannten, nicht das MR davor) ist damit
      geschlossen vermerkt. Belegt in `08373c9`.
- [x] **Das Delta ist gelesen, nicht angenommen:** was `v5.13.1`, `v5.14.0` und `v5.15.0`
      gegenüber `v5.13.0` tragen, steht im Slice — der Herausgeber hat nur
      `v5.13.0` angekündigt. Siehe §1: Wellen 99–106 je Tag zugeordnet, 52
      Bundle-Dateien nach unverändert/Stempel/echtem-Inhalt klassifiziert.
- [x] `make gates` und `make fullbuild` grün (Exit explizit); **unabhängiger
      Review**; **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.
      `make gates` grün: zehn Gates, 647 Dateien, 0 Befunde, Coverage 94,70 %
      (2026-09-01). `make fullbuild` grün: Image-Hash
      `sha256:7b8cb29549c4553300d87ee9e382dfbb3630068b63cf945f3a095f2f2caa3b1e`
      (2026-09-01). **Review** (eigener Kontext): 0 HIGH, 1 MEDIUM, 1 LOW —
      beide behoben in `d32c13d`. **Verifikation** (eigener Kontext, nach
      `d32c13d`): alle acht geprüften Punkte konform, Zahlen/Hashes/Zitate
      unabhängig nachgerechnet und deckungsgleich; Freigabe für `done/`
      erteilt.

## 5. Abnahme-Punkte / Risiken

- **Der Bump-Alarm ist gewollt, aber er unterscheidet nicht.**
  [`MR-051`](../../../../harness/conventions.md#mr-051) nennt den Preis: eine
  Direktive meldet **auch dann**, wenn nur die Zeilennummern gewandert sind. Bei
  33 Direktiven ist das Arbeit ohne inhaltlichen Anlass, und die Versuchung ist,
  sie mechanisch nachzuziehen statt den Grund-Code zu lesen. — **Ausgang:
  entfallen.** Das differenzierte Ergebnis (2 nachgezogen, 0 entfernt, 22 gegen
  den Datei-Diff geprüft statt pauschal übernommen) belegt, dass je Direktive
  der Grund-Code gelesen wurde, nicht die Zeilennummer.
- **Der Adaptions-Review ist Urteil, kein `grep`.** Fünf Ausgänge je Eintrag,
  33 Einträge — die Gefahr ist nicht, einen falsch zu entscheiden, sondern alle
  reflexhaft auf „bleibt gültig" zu setzen, weil das der Normalfall ist. —
  **Ausgang: entfallen.** Von 33 Einträgen wurden 2 als nicht trivial
  identifiziert ([`MR-042`](../../../../harness/conventions.md#mr-042)/`043`/[`MR-048`](../../../../harness/conventions.md#mr-048)
  aktualisiert, [`MR-013`](../../../../harness/conventions.md#mr-013)
  teilweise überholt) — ein rein reflexhafter Durchgang hätte alle 33 auf
  „bleibt gültig" gesetzt.
- **Vier Releases auf einmal.** Angekündigt ist nur `v5.13.0`; `v5.13.1`,
  `v5.14.0` und `v5.15.0` kamen ohne Ankündigung. Ein Sprung über vier Fassungen macht das Delta
  größer als den Diff, den ein Reviewer in einer Sitzung prüft. — **Ausgang:
  entfallen.** Die Datei-Klassifikation (unverändert · nur Versions-Stempel ·
  echter Regel-Inhalt) hat die 52 Bundle-Dateien auf 14 inhaltlich relevante
  reduziert und den Diff dadurch review-fähig gemacht — als geschärfte Regel
  im Lerneintrag (§9) festgehalten.
- **Der Bump macht eine fremde Adaption beurteilbar, die auf uns zeigt.** Die
  `--print-mk`-Frage hängt an `modul-14 §Der Prüflauf ist hermetisch`; solange
  der Abschnitt nicht im Baum liegt, wäre jede Antwort darauf ein Urteil aus dem
  Zitat eines Dritten ([`BEO-012`](../observations.md)). — **Ausgang:
  entfallen.** Die Vorbedingung des Risikos (unvollständige Quelle) ist mit dem
  Vendoring von `v5.15.0` behoben — `modul-14` liegt jetzt vollständig im Baum.
  Die Frage selbst bleibt unbeantwortet, aber das ist eine bewusste
  Scope-Entscheidung (§3), kein offenes Risiko dieses Slice mehr.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — `in-progress/` trägt
keinen Slice; `v5.15.0` ist upstream verfügbar (gemessen mit
`make baseline-freshness`).

**Rückführungen:** `in-progress` → `open`, falls das Delta einen Ausgang
**gegenstandslos** oder **teilweise überholt** ergibt, dessen Rückbau selbst ein
Slice ist — dann trägt dieser Slice den Bump und der Rückbau folgt getrennt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `.harness/baseline/` (vendorte Fremd-Konvention) und
  `harness/` (der Konventionsspeicher). Beide fallen unter den Default `*` =
  **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration). Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.15.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** — **bei der Beanspruchung aufgefrischt**, weil
  der Plan seit dem 2026-08-30 in `open/` lag und das Register sich seither
  bewegt hat (Stand 2026-08-31, höchste Kennung `BEO-025`):
  [`BEO-012`](../observations.md) (**12**, war 11) — ein Zitat über seinen
  Geltungsbereich hinaus: ein Bump ankert `cite`-Spannen neu, und eine Spanne,
  die zufällig wieder wortgleich trifft, ist kein Beleg dafür, dass sie noch
  dieselbe Regel zitiert; [`BEO-008`](../observations.md) (**4**) — die drei
  Klassen einer Pin-Hebung, von denen nur die grep-bare gehoben wird, plus die
  vierte (der zitierende Verweis), die genau bei diesem Vorgang eintritt;
  [`BEO-002`](../observations.md) (**7**) — die Spiegel des Pins, und der
  geschärfte Ableiter dazu: die Liste entsteht aus einem `grep` nach dem
  **alten** Wortlaut, nicht aus dem Gedächtnis; [`BEO-013`](../observations.md)
  (**1**) — ein Wächter, der nichts mehr fängt, bleibt stehen: die
  Bestands-Stichprobe des Freshness-Audits gehört gefahren, auch wenn der Pin
  aktuell ist. **Neu seit der Anlage und für diesen Slice einschlägig:**
  [`BEO-025`](../observations.md) (**1**) — ein Liefer-Punkt landet im Commit
  eines fremden Slice; ein Bump zerfällt in mehrere Commits, und das ist genau
  die Lage, in der die Zuordnung verrutscht.
  <!-- d-check:cite .harness/baseline/v5.15.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)) — **bei der
  Beanspruchung neu gelesen:** `upstream-drift.yml` meldet **ROT** (jüngster
  Lauf 2026-08-31T06:31:40Z), `image-scan.yml` grün. Gelesen statt weggeklickt:
  drei Schritte fielen — `freshness-a-check` und `go-base-digest` sind seither
  **behoben** (beide Pins gehoben, alle acht Achsen melden `ok`), und der
  dritte ist `baseline-freshness`, **also genau dieser Slice**. Der Nachtlauf
  wird grün, wenn er geschlossen ist. Die benannte Grenze des Targets zeigt
  sich dabei: es liest den **jüngsten** Lauf, nicht sein Alter — zwei der drei
  Ursachen sind längst weg, der Lauf ist nur nicht wiederholt. **Dieser Block
  trägt bewusst keine `cite`-Direktive** — sein Ziel ist eine Repo-Adaption
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-183. Betroffene IDs:
[`MR-021`](../../../../harness/conventions.md#mr-021),
[`MR-039`](../../../../harness/conventions.md#mr-039),
[`MR-048`](../../../../harness/conventions.md#mr-048),
[`MR-051`](../../../../harness/conventions.md#mr-051),
[`MR-055`](../../../../harness/conventions.md#mr-055). Module: `links`,
`anchors`, `citations`. Gates: `make baseline-verify`, `make gates`,
`make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter den
Default: Doc führt, Code folgt. Ein Adoptions-Stand und sein Konventionsspeicher;
kein Produkt-Code, keine Reconciliation. Das **Evidenz-Risiko** ist die einzige
Achse mit Substanz: der vendorte Baum ist Fremd-Inhalt, und ob eine Adaption
noch trägt, entscheidet sein Delta — nicht unser Bestand.

## 9. Closure-Notiz (nach `done/`)

**Geliefert.** Der vendorte Baseline-Baum steht auf `v5.15.0` — vier Releases
in einem Bump statt der ursprünglich geplanten einen (§1), Wellen 99–106 je
Tag zugeordnet, nicht angenommen. Alle 129 lebenden `v5.12.0`-Pfadverweise
gezogen (95 gehoben, 34 als Aussagen über die Vergangenheit stehen gelassen),
die vier `.claude/rules`-Symlink-Aliase mitgezogen, 24 `d-check:cite`-Direktiven
entschieden (2 nachgezogen, 0 entfernt, 22 bestätigt), ein Adaptions-Review
über alle 33 lebenden `MR`-Einträge gefahren (31 bleiben gültig,
[`MR-048`](../../../../harness/conventions.md#mr-048) ist ein eigenständiger
Fund, [`MR-013`](../../../../harness/conventions.md#mr-013) ist teilweise
überholt), [`MR-037`](../../../../harness/conventions/done/MR-037-baseline-v5120.md)
nach `conventions/done/` migriert,
[`MR-057`](../../../../harness/conventions.md#mr-057) als Nachfolge-Eintrag
angelegt.

**Was funktioniert hat.** Die Datei-Klassifikation (unverändert · nur
Versions-Stempel · echter Regel-Inhalt) hat den Vier-Releases-Sprung auf 14
inhaltlich relevante Dateien reduziert und damit den Diff review-fähig
gehalten, obwohl der Herausgeber nur eines der vier Releases angekündigt
hatte. Die differenzierten Ausgänge in §5 (2 von 33 MR-Einträgen als
nicht-trivial markiert, 2 von 24 Direktiven nachgezogen) belegen, dass die
Versuchung zum reflexhaften Pauschal-Urteil nicht eingetreten ist.

**Was anders lief.** Der unabhängige Review (dde0916..HEAD) fand 0 HIGH, 1
MEDIUM, 1 LOW — beide behoben in `d32c13d`: eine seit der Plan-Anlage
stehengebliebene Vorab-Schätzung unter dem Wort „Gemessen" (§2 Schritt 3;
Ursache und Reparatur jetzt als siebte Instanz von
[`BEO-009`](../observations.md) registriert, erstmals an einem Slice-Plan statt
an einer Commit-Botschaft), und eine unbegründete Klassifikations-Asymmetrie
in [`MR-057`](../../../../harness/conventions.md#mr-057) zwischen
[`MR-048`](../../../../harness/conventions.md#mr-048) und
[`MR-042`](../../../../harness/conventions.md#mr-042)/[`MR-043`](../../../../harness/conventions.md#mr-043)
(jetzt mit dem `F-7`-Bezug präzisiert — der ursprüngliche Plan-Text hatte
`F-7` fälschlich dem einen statt den beiden anderen zugeordnet). Die
unabhängige Verifikation (nach `d32c13d`) fand
alle acht geprüften DoD-Punkte konform, mit unabhängig nachgerechneten
Zahlen, Hashes und Zitat-Proben.

**Ein benannter Verfahrens-Fehler, klein aber real.** Der letzte DoD-Haken
(Review/Verifikation) wurde in der Datei gesetzt, **bevor** der `git mv`
folgte — dieser Edit blieb ungecommittet, bis der `git mv` ihn mitnahm. Der
Move-Commit (`11c614d`) ist damit **kein** byte-identischer Rename (98 %
Similarity statt 100 %), anders als es
[`AGENTS.md`](../../../../AGENTS.md) §3.3 für die Slice-Lifecycle-Move-Ausnahme
vorsieht ("die Slice-Datei selbst ist im Move-Commit unverändert"). Praktisch
folgenlos — `git log --follow` hält (die Similarity liegt weit über der
50-%-Schwelle, gemessen), `make gates` blieb grün —, aber es ist eine
Abweichung von der etablierten Form (verglichen mit dem Vorgehen bei
[slice-182](wellenlos/slice-182-erklaerte-leermenge.md): dort blieben alle DoD-Haken bis
zum Closure-Body-Commit offen). **Prozedur, präzisiert:** der letzte DoD-Haken
gehört wie alle anderen erst in den Body-Commit nach dem Move, nicht in einen
Zwischen-Commit davor.

**Zwei offene Punkte, benannt statt eingelöst (§2b).** Der vom Kanon
verlangte Versions-Sensor auf den Baseline-Pin fehlt weiterhin — die
Kurzform der Kurs-Vorlage ist hier von
[`DC-FA-VER-001`](../../../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in)
belegt, ein Folge-Slice ist nicht geschnitten. Die
[`MR-013`](../../../../harness/conventions.md#mr-013)/[`MR-057`](../../../../harness/conventions.md#mr-057)-Kollision
(Kurs-Welle 103 verlangt für den `MR`-Lifecycle-Move den umgekehrten
Commit-Schnitt) bleibt gemeldet statt aufgelöst — die neue Kanon-Bedingung
liegt in keinem Release (`v5.16.0` trägt Welle 109, nicht 110) und ist damit
noch nicht zitierbar; ein Folge-Eintrag mit engerem Nachfolger für die dritte
Move-Klasse ist fällig, sobald sie es ist.

**Verifikation.** `make gates` Exit 0 (zehn Gates, 647 Dateien, 0 Befunde,
Coverage 94,70 %) · `make fullbuild` Exit 0, Image-Hash
`sha256:7b8cb29549c4553300d87ee9e382dfbb3630068b63cf945f3a095f2f2caa3b1e`, 50
Anforderungen / 0 Waisen · unabhängiger Review (0/1/1, behoben) · unabhängige
Verifikation (8/8 konform).
