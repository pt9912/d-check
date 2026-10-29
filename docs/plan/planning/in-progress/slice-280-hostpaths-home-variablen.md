# Slice slice-280: `hostpaths` erkennt `$HOME` und `${HOME}` als Home-Verweis

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-92

**Bezug:** [`DC-FA-HOST-001`](../../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in) (Scope),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus) (Byte-Identität ohne den neuen Schlüssel). Keine *aktive* ADR trägt die
Erkennungsform; eine begleitende ADR hält die Entscheidung fest (Nummer bei der Anlage).

**Berührte Spec-Stellen:** [`DC-FA-HOST-001`](../../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in)
(Beschreibung, Out-of-Scope, Akzeptanzkriterien) ·
[`DC-FA-HOST-001.a`](../../../../spec/spezifikation.md#dc-fa-host-001a--host-pfad-erkennung) und die
Schema-Zeilen `hostpaths.*` (`spezifikation.md §2`).
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-29.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das Modul `hostpaths` meldet auf Wunsch `$HOME/…`, `${HOME}/…` und
`"$HOME"/…` als Home-Verweis — in Prosa und Inline-Code wie die Tilde —, über
einen neuen opt-in-Schlüssel; ohne ihn bleibt der Befundsatz byte-identisch.

**Gemessener Anlass** (GitHub-Issue 12, gegen `v0.87.0`): ein Tilde-Pfad
und ein absoluter Pfad unter einem Host-Präfix melden `hostpath-forbidden`, `$HOME/Dev/x/y.log`,
`${HOME}/Dev/z` und `"$HOME"/Dev/w` melden nichts. Ursache am Code: die
Wortgrenzen-Vorbedingung des Unix-Musters sieht vor `/Dev…` das `E` von `HOME`.
Dass `hostpaths` das einzige Modul ist, das Inline-Code mitliest, macht die
Lücke zum blinden Fleck für Maschinen-Layout-Leaks. Der Vorgänger-Slice
[`slice-249`](../done/wellenlos/slice-249-hostpaths-tilde-exempt-targets.md)
hatte die Variablen-Formen ausgeschlossen, *weil kein belegter Anlass bestand
(Bestand: null Vorkommen)* — dieser Anlass ist jetzt belegt, der Ausschluss
fällt damit.

**Die Erkennung — Entwurf, im Review zu prüfen, nicht gesetzt:**

- **Schlüssel `hostpaths.variables`**: Liste von Variablen-Namen. Abwesend ⇒
  Verhalten und Befundsatz byte-identisch zu heute. Der Wunsch des Issues
  nennt „Standard `[HOME]`" **und** „damit ein bisher grüner Lauf nicht
  ungefragt rot wird"; beides zugleich geht nicht. Entwurf: abwesend = aus,
  `[HOME]` ist der Wert, den das Konfigurations-Gerüst vorschlägt.
- **Treffer:** `$NAME` oder `${NAME}` (NAME aus der Liste), dann optional ein
  schließendes `"`, dann `/` und ein erstes Segment, das **nicht** mit `.`
  beginnt — dieselbe Regel wie bei der Tilde (Werkzeug-Konventionen wie
  `$HOME/.config` bleiben still). Die Wortgrenzen-Vorbedingung ist die der
  Tilde, ergänzt um `$` (kein `$$HOME`).
- **Gemeldet** wird der Pfad in voller Form samt Variable, genau einmal; ein
  Unix-Treffer, der innerhalb eines Variablen-Treffers beginnt, entfällt
  (ein Variablen-Pfad, dessen erstes Segment zugleich ein Host-Präfix-Name
  ist, träfe sonst zweimal).
- **`hostpaths.exempt-targets`** akzeptiert Globs, die mit `$` beginnen — aber
  nur, wenn `variables` gesetzt ist und das Glob mit einer dort genannten
  Variablen beginnt; sonst bliebe es wie heute still wirkungslos und bleibt
  deshalb Exit 2 (Rand-Disziplin des Lastenhefts: ein Glob, das nie trifft, ist
  ein Konfigurationsfehler).
- **Variablen-Namen** sind Shell-Bezeichner (`[A-Za-z_][A-Za-z0-9_]*`); ein
  leerer oder anders geformter Name ist Exit 2.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`$USERPROFILE`, `%USERPROFILE%`, `$XDG_*` und andere Variablen als
  Produkt-Default** — nur `HOME` ist belegt. Die Liste trägt weitere Namen,
  ohne dass der Slice sie als Anforderung festschreibt. Die Windows-Schreibweise
  `%NAME%` ist eine andere Syntax und bleibt im Lastenheft ausgeschlossen.
- **Variablen mit Default-Syntax** (`${HOME:-/tmp}/x`) — kein Home-Verweis,
  sondern ein Ausdruck; in der Negativliste als bewusst still geführt.
- **Eine Variablen-Form ohne Schrägstrich** (nacktes `$HOME`) — wie die nackte
  Tilde: ein Verweis auf das Home-Verzeichnis leakt kein Layout.
- **`hostpaths` in `.d-check.yml` dieses Repos mit `variables` zu aktivieren** —
  eine eigene Entscheidung: das Repo führt `$HOME` in eingefrorenen
  `done/`-Slices und in diesem Plan; die Aktivierung färbte `make doc-check`
  rot an Stellen, die nicht mehr geändert werden. Der Slice liefert das Modul,
  nicht die Anwendung auf sich selbst.
- **Ein Zeilen-Marker** — bleibt ausgeschlossen
  ([`DC-FA-HOST-001`](../../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in)); die Auswege sind Fences,
  Präfixliste und Ziel-Ventil.
- **Handbuch, README und CHANGELOG** — Release-Prep, nicht Feature-Commit
  (`AGENTS.md` §5 Zeile 17).

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`DC-FA-HOST-001`](../../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in) erfüllt, Test referenziert:
      Happy (`$HOME/Dev/x`, `${HOME}/Dev/z`, `"$HOME"/Dev/w` melden mit
      gesetztem `variables`) · Boundary (ohne den Schlüssel kein Befund, Befundsatz
      byte-identisch; `$HOME/.config/x` und nacktes `$HOME` still) · Negative
      (Fence ausgenommen, `exempt-targets` mit `$`-Glob entfernt genau diesen
      Fund, `$`-Glob ohne passende Variable ist Exit 2). Rot-Beleg: der
      `$HOME`-Test läuft ohne die Änderung aus dem richtigen Grund rot (kein
      Befund statt `hostpath-forbidden`).
- [ ] Fassung nachgezogen: Lastenheft ([`DC-FA-HOST-001`](../../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in): Beschreibung,
      Out-of-Scope-Satz zu Variablen-Formen, Kriterien, Version, Historie-Zeile),
      Spezifikation ([`DC-FA-HOST-001.a`](../../../../spec/spezifikation.md#dc-fa-host-001a--host-pfad-erkennung), Schema-Zeile
      `hostpaths.variables`, Historie), begleitende ADR samt ADR-Index.
- [ ] Spiegel nachgezogen: das Konfigurations-Gerüst (`--print-config`) führt
      den Schlüssel; die Zählung der übrigen Fundstellen steht in §3.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für die Schnittstelle `hostpaths` falls öffentlicher Vertrag berührt (Spec und Gerüst oben; Handbuch ist Release-Prep).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/core/rules/hostpaths.go` | update | Variablen-Muster aus der Namensliste; Treffer innerhalb eines Variablen-Treffers verdrängt den Unix-Treffer — [`DC-FA-HOST-001`](../../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in) |
| `internal/hexagon/core/model/config.go` | update | `HostpathsConfig` bekommt die Namensliste |
| `internal/adapter/driven/configyaml/configyaml.go` | update | `hostpaths.variables` lesen und validieren; `exempt-targets` akzeptiert `$`-Globs nur mit passender Variable |
| `internal/hexagon/core/rules/hostpaths_test.go` | update | Happy/Boundary/Negative nach [`DC-FA-HOST-001`](../../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in); Negativliste unten |
| `internal/adapter/driven/configyaml/configyaml_test.go` | update | Config-Rand: leerer/ungültiger Name, `$`-Glob ohne Variable → Exit 2 |
| `internal/adapter/driving/cli/config_template.go` | update | Gerüst führt `hostpaths.variables` (Spiegel) |
| `spec/lastenheft.md` | update | [`DC-FA-HOST-001`](../../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in): Beschreibung, Out-of-Scope, Kriterien, Version, Historie-Zeile ([`MR-032`](../../../../harness/conventions.md#mr-032)) |
| `spec/spezifikation.md` | update | [`DC-FA-HOST-001.a`](../../../../spec/spezifikation.md#dc-fa-host-001a--host-pfad-erkennung), Schema-Zeile, Historie |
| `docs/plan/adr/` · `docs/plan/adr/README.md` | neu · update | begleitende ADR (Alternativen: Default `[HOME]` an · opt-in abwesend = aus · Variablen-Ausnahme im Unix-Muster); Index nachgezogen |

**Zählung der Spiegel (Schritt 17/18 des Workflow-Skeletts), vor dem Code
nachzuzählen:**

- Dateien, die `hostpaths` beschreiben: `grep -rlE 'hostpaths' --include=*.go --include=*.md --include=*.yml internal spec docs/user harness .claude .d-check.yml | grep -v '_test.go\|done/\|docs/reviews'` → 13 Treffer: `.d-check.yml`, `docs/user/abdeckung-tests.md`, `docs/user/benutzerhandbuch.md`, `docs/user/operations.md`, `harness/sensors/doc-check.md`, `configyaml.go`, `config_template.go`, `suggest.go`, `model/config.go`, `rules/hostpaths.go`, `rules/run.go`, `spec/lastenheft.md`, `spec/spezifikation.md`. Davon trägt dieser Slice: Code, Spec, Gerüst; `docs/user/*` ist Release-Prep (§1); `suggest.go` und `run.go` sind auf Aussagen über Präfixliste und Ventil zu lesen, bevor sie als unberührt gelten; `harness/sensors/doc-check.md` ist als Sensor-Seite auf Grenz-Aussagen zu prüfen.
- Vorkommen von `$HOME` im lebenden Bestand: `grep -rlE '\$\{?HOME\}?' --include=*.md . | grep -v '^./.harness/baseline'` → 2 Dateien (ein eingefrorener Slice in `done/` und der Welle-Plan dieser Welle).

**Negativliste der Erkennung (Schritt 19, vor dem Code-Commit zu testen):**
`$HOMEDIR/x` (längerer Name) · `$HOME` ohne Schrägstrich · `$HOME/.config/x` ·
`$HOME//x` · `${HOME:-/tmp}/x` · `${HOME}x/y` · `x$HOME/y` (Wortgrenze) ·
`$$HOME/x` · `https://h/$HOME/y` (URL) · `**$HOME/x**` und `_$HOME/x_`
(Markdown-Hervorhebung am Rand) · Pfad im Fenced-Block · Zeilenumbruch
zwischen Variable und Schrägstrich · Backslash-Maskierung `\$HOME/x` · Variable
nicht in der Liste (`$OTHER/x`) · ein Variablen-Pfad mit Host-Präfix-Namen als erstem Segment (genau ein Befund, in
voller Form). Dazu der Lauf über den eigenen Bestand: Befundsatz von `hostpaths`
ohne den Schlüssel gegen `v0.87.1` byte-gleich, mit `variables: [HOME]` die
Menge der Treffer vorab gezählt (erwartet: die zwei Dateien oben).
**Plan-Änderungen (Schritt 19 und Review):**

- **Verdrängung nur bei `$NAME/` redundant.** Die Negativliste nannte „ein Variablen-Pfad
  mit Host-Präfix-Namen als erstem Segment“ und begründete die Verdrängung eines
  Unix-Treffers damit. Bei `$HOME/<Präfix>/x` trifft das Unix-Muster ohnehin nicht
  (vor dem Präfix steht das `E` des Namens); bei `${HOME}/<Präfix>/x` und
  `"$HOME"/<Präfix>/x` steht `}` bzw. `"` davor, das Muster träfe, und die Verdrängung
  verhindert den zweiten Fund. Die Tests fahren alle drei Formen und die `=`-Form.
- **Verdrängung ist allgemein** (Review F-3): der früher beginnende Treffer gewinnt in
  allen Richtungen zwischen Tilde, Variable und Unix. Ohne `variables` bleibt das
  Verhalten gleich, denn dann gibt es nur Tilde und Unix.
- **`**$HOME/x**` ist kein „kein Fund“** (Review F-4, Umkehrung der Negativliste): die
  Hervorhebung am Rand verhält sich wie bei der Tilde und meldet mit den Stern-Zeichen
  im Ziel. Die Negativliste hatte das offen gelassen; die Spezifikation führt die
  Hervorhebung als benannte Grenze.
- **Das Ziel eines Funds hinter einer Variablen wechselt** (Review F-2): mit gesetztem
  `variables` meldet `${HOME}/<Präfix>/x` unter dem Variablen-Pfad, ein Bestands-Glob
  auf den Unix-Teil trifft nicht mehr. Das gehört in Lastenheft, Spezifikation und ADR und
  steht dort.
- **Anlass-Beispiel:** `$HOME/Dev/…` aus dem Issue setzt ein konfiguriertes Präfix `Dev`
  voraus; im Default heißt das Präfix `Development` (Review F-7).


## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): der erste Slice der Welle (siehe deren §4)
liegt in `done/` (beide schreiben eine Lastenheft-Fassung; die Reihenfolge
verhindert den Versions-Konflikt), und welle-92 ist offen. Der Übergang landet
auf dem Hauptzweig, vor der Arbeit.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Negativliste
  verlangt eine eigene Lexik für Shell-Ausdrücke (Parameter-Expansion, Quoting)
  statt eines Variablen-Musters neben der Tilde — dann gehört die
  Ausdrucks-Erkennung in einen eigenen Slice.
- `in-progress` → `open` (blockiert — Carveout?): der Auftraggeber entscheidet
  die Kopplung von `exempt-targets` und `variables` (§1) anders, oder die
  Lastenheft-Fassung wird nicht freigegeben.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, Review-Report liegt vor, und am gebauten
Image meldet das Fixture aus Issue 12 mit `variables: [HOME]` alle drei Formen,
während derselbe Lauf ohne den Schlüssel gegen `v0.87.1` byte-gleich ist. Dazu
die Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Das Variablen-Muster liest Shell-Ausdrücke, die kein Home-Verweis sind
  (Parameter-Expansion, Befehlsersetzung), als Pfad und meldet falsch positiv.
  Die Negativliste in §3 deckt die benannten Formen; die Klasse ist offen —
  **Ausgang:** wird bei der Closure gesetzt (eingetreten / entfallen / weiter offen)
- Die Kopplung von `exempt-targets` an `variables` ist eine Entwurfs-Wahl, die
  das Issue anders formuliert (Globs mit `$` *akzeptieren*, ohne Bedingung).
  Wird sie beim Review umgestoßen, ändert sich die Config-Rand-Zusage im
  Lastenheft — **Ausgang:** wird bei der Closure gesetzt (eingetreten / entfallen / weiter offen)
- Wer `variables` aktiviert, bekommt Funde in Dokumenten, die bisher grün
  waren; das ist beabsichtigt (opt-in), aber die Meldung nennt den Schlüssel
  nicht — **Ausgang:** wird bei der Closure gesetzt (eingetreten / entfallen / weiter offen)

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist der Default `*` (`ALL`) —
Regelmodul, Konfigurations-Rand, Spec, Gerüst; keine Sub-Area mit eigener
Modus-Deklaration (`tools/harness/` bleibt unberührt).

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen, Sub-Area
`*`. Drei Treffer:
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/observation.md)
(5×, verkörpert in Schritt 19 des Workflow-Skeletts — die Negativliste in §3
ist seine Anwendung; das Issue nennt drei Positivfälle),
[`BEO-ALL/prosa-verschwindet-in-inline-code-spanne`](../observations/BEO-ALL/prosa-verschwindet-in-inline-code-spanne/observation.md)
(2×, offen — `hostpaths` liest Inline-Code, und die Beobachtung betrifft
absatzweit gepaarte Backticks; ob eine fehlgepaarte Spanne auch den
Variablen-Treffer verschwinden lässt, ist ein Fall für die Negativliste, nicht
vorab entschieden) und
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
(3×, verkörpert in `AGENTS.md` §3.8 — die Zusage gilt für die Scan-Menge des
Moduls; `exempt-targets` ist eine Eingabe, die nicht gescannt wird). Der Zähler
folgt aus den Dateien der Closure.

### Sub-Area: `*` (Default)

Alle berührten Sub-Areas GF.
