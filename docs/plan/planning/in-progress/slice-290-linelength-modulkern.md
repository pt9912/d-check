# Slice slice-290: Modulkern `linelength`: Konfiguration lesen, Zeilen messen, Ausnahmen anwenden, Grund-Code

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-95-linelength-modul`](../welle-95-linelength-modul.md).

**Bezug:** [`DC-FA-LLEN-001`](../../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in) und
[ADR-0114](../../adr/0114-linelength-modul.md) aus
[`slice-289`](../done/slice-289-linelength-anforderung-und-adr.md),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit),
[`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl).

**Berührte Spec-Stellen:** `spezifikation.md` §4 (die Zeile des Grund-Codes),
`lastenheft.md` (Modul-Aufzählung in
[`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl) und
Glossar). Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie.

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-11-22.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das opt-in-Modul `linelength` ist registriert, liest seinen Block am
Konfigurations-Rand (Fehlkonfiguration ⇒ Exit 2), misst jede Zeile der Markdown-Dateien
seiner Scan-Menge nach den in [`slice-289`](../done/slice-289-linelength-anforderung-und-adr.md)
bestätigten Entscheiden und meldet jede überlange Zeile mit ihrem Grund-Code — jede
Randbedingung der Anforderung hat einen Test, der ohne die Regel aus dem richtigen
Grund rot läuft.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Anforderung, Verfahren, ADR ändern** — gehören
  [`slice-289`](../done/slice-289-linelength-anforderung-und-adr.md). Zeigt der Code, dass ein
  Entscheid nicht trägt, geht das als Plan-Änderung zurück an die Anforderung, nicht als
  stille Abweichung im Code.
- **Das Modul im eigenen Profil (`.d-check.yml`) einschalten** — Welle-Out-of-Scope:
  der Bestand führt bewusst lange Zeilen; das Einschalten wäre eine Umformatierung,
  keine Lieferung dieses Slice.
- **Eine eigene Fence- oder Tabellen-Erkennung** — die Ausnahmen kommen aus der
  geteilten Lexik; ein zweites Muster daneben ist genau die Klasse
  [`BEO-ALL/gemeinsame-leser-sind-zwei-familien`](../observations/BEO-ALL/gemeinsame-leser-sind-zwei-familien/observation.md).
- **`--fix-patch`, `--doctor`-Reparaturvorschlag, `--suggest-config`-Erkennung** — eigene
  Vorgänge; der Diagnose-Klartext des Grund-Codes kommt mit (Lockstep-Test), mehr nicht.
- **CHANGELOG, Handbuch-Kopf, Release** — Release-Prep
  ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17). Was ein Spiegel-Test schon im
  Feature-Commit verlangt (Modul-Bullets der READMEs, Modul-Abschnitt und
  Modul-Aufzählung im Handbuch), kommt mit und steht in §3.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Konfigurations-Rand: der Block `linelength` wird gelesen, jede Fehlkonfiguration der
      Anforderung (fehlende oder < 1 Schwelle, unbekannter Schlüssel, ungültiger
      `scope`/`exempt-paths`) ist Exit 2 mit Test; `--print-config` und `--enable` kennen
      das Modul.
- [x] Regel: Zählung in Codepunkten, Grenze genau auf der Schwelle befundfrei, jede
      Ausnahme einzeln schaltbar und getestet (an und aus), ein Befund je Zeile mit Ort,
      Länge und Grenze, Reihenfolge nach Pfad und Zeile, Bezugsmenge (Zeilen wie bei `file`
      gezählt, die über (c) ausgenommenen als eigene Zahl), leere Menge kein Fehler; die
      Negativfälle aus dem Plan von
      [`slice-289`](../done/slice-289-linelength-anforderung-und-adr.md) stehen als Testfälle, und
      für jeden Test ist ein Rot ohne die Regel gezeigt (Baseline-Regelwerk
      `modul-11-verification.md` §Bewusstes Brechen). Der Kopplungs-Test fährt dieselbe Eingabe
      durch `linelength` und die bisherigen Konsumenten von `proseLineSet`/`tableRowLine` —
      Pflicht nach [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
      Entscheidung 4, keine Option (Architect-Verdikt F-3).
- [x] Spiegel: Grund-Code in `AllReasons()`/Diagnose-Klartext und §4 der Spezifikation,
      Modul in der Aufzählung von
      [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl),
      im Glossar, in den READMEs und im Handbuch — die Spiegel-Deckungstests sind grün und
      fangen das neue Modul (vor dem Editieren aufgelistet,
      [`MR-025`](../../../../harness/conventions.md#mr-025)).
- [x] `make gates` grün.
- [x] [ADR-0114](../../adr/0114-linelength-modul.md) `Accepted`, Index nachgezogen — der
      Status-Übergang durch den Architect, nachdem Review und Verifikation dieses Slice keine
      Änderung an der Entscheidung mehr verlangen (Architect-Verdikt F-8); vorher nennt das
      `Schärft:`-Feld die `SPEC-<NNN>`-Kennung der neuen §4-Zeile des Grund-Codes (Review R1 zu
      slice-289, F-11).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **mit** Wellen von der nächsten Welle-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/core/model/config.go` | update | Konfigurations-Typ des Moduls; Modulliste (`ValidModules()`) |
| `internal/adapter/driven/configyaml/configyaml.go` | update | Block lesen und am Rand prüfen (Exit 2) |
| Regel-Datei `linelength` im Verzeichnis `internal/hexagon/core/rules/` | neu | Zeilen messen, Ausnahmen über `proseLineSet`/`tableRowLine`, Befunde, Bezugsmenge |
| `internal/hexagon/core/rules/run.go` | update | Modul im Lauf verdrahten (Scan-Menge mit `scope`/`exempt-paths`) |
| `internal/hexagon/core/app/diagnose.go` | update | `AllReasons()` und Diagnose-Klartext des Grund-Codes (Lockstep-Test) |
| `internal/adapter/driving/cli/config_template.go` | update | `--print-config`: Zeile „Verfügbar“ und Block `linelength` |
| `spec/spezifikation.md` §4 | update | Zeile des Grund-Codes — kommt **mit** dem Code, der Lockstep-Test hält die Tabelle gegen `AllReasons()` |
| `spec/lastenheft.md` (Modul-Aufzählung in [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl), Glossar) | update | nennt `linelength` erst, wenn das Modul registriert ist; Version und Historie nach [`MR-032`](../../../../harness/conventions.md#mr-032) |
| `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`, `docs/user/operations.md` | update | Modul-Bullets und -Aufzählungen, die ein Spiegel-Test schon im Feature-Commit verlangt; der Handbuch-Kopf bleibt Release-Prep |
| `docs/user/abdeckung-tests.md` | update | von `make abdeckung` abgeleitet; `make test` hält sie gegen die Ableitung |
| Testdateien neben Regel und Konfigurations-Rand | neu | Happy/Boundary/Negative je Randbedingung der Anforderung; Spiegel-Deckung |

**Spiegel vor dem Editieren auflisten** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
die Tabelle oben ist der Stand der Planung, nicht der Suche. Vor dem ersten Edit läuft die
Aufzählungs-Suche nach einem bestehenden Modulnamen (etwa `pinsync`) über das ganze Repo;
jeder Treffer, der in der Tabelle fehlt, kommt als Plan-Nachtrag vor den Code.

**Testfälle aus der Anforderung** (Negativliste nach Schritt 19 des Workflow-Skeletts; die
Liste folgt den bestätigten Entscheiden 1–9 und dem
[Architect-Verdikt](../../../reviews/2026-11-22-welle-95-architect-verdikt.md); maßgeblich ist
die Anforderung, sobald sie steht):

- Zeile mit `max-chars` Zeichen ⇒ befundfrei; mit `max-chars + 1` ⇒ ein Befund — mit
  `exempt-unbreakable: false` oder ohne Wort über der Grenze; eine Zeile aus einem einzigen Wort
  mit `max-chars + 1` Zeichen ist unter Variante H im Default ausgenommen ([ADR-0114](../../adr/0114-linelength-modul.md) Entscheid 5 (c)).
- Zeile aus Mehrbyte-Zeichen (Umlaut, Emoji) an der Grenze ⇒ gezählt in Codepunkten,
  nicht in Bytes.
- `\r\n`-Zeilenende ⇒ das `\r` zählt nicht; Tab zählt als ein Zeichen; ungültiges UTF-8
  nach der Regel der Anforderung.
- Überlange Zeile in Fenced-Code (Backtick- und Tilde-Zaun, eingerückter Zaun) ⇒ mit
  Ausnahme an befundfrei, mit Ausnahme aus ein Befund; die Zaunzeile selbst ebenso.
- `~~~`-Zeile in einem Backtick-Block ⇒ der Block schließt früh, die Zeilen danach werden
  gemessen; eingerückter Code-Block (vier Leerzeichen, ohne Zaun) ⇒ gemessen (Verdikt F-3).
- Nicht geschlossener Fence bis Dateiende ⇒ der Rest der Datei ist ausgenommen, wie die
  Anforderung es als Grenze benennt; gemeldet wird der offene Zaun von `spans`
  (`fence-unclosed`), nicht von `linelength` (Verdikt F-3).
- Überlange `|`-Zeile ohne Trennzeile ⇒ ausgenommen (die Lexik verlangt keine Trennzeile);
  eingerückte `|`-Zeile, Zeile in einem Zitat (`> |`) und GFM-Zeile ohne führende Pipe ⇒
  gemessen; `exempt-code-blocks: false` mit `exempt-tables: true` ⇒ eine `|`-Zeile im Zaun wird
  gemessen (Verdikt F-4).
- Ausnahme (c) nach Variante H (Entscheid 5 (c), neu bestätigt 2026-11-23): Zeile mit einem Wort
  (längste Folge ohne Space und Tab), das allein länger als `max-chars` ist (lange URL) ⇒
  ausgenommen, auch mit umbrechbarem Text davor; Zeile über der Grenze, deren längstes Wort genau
  `max-chars` lang ist ⇒ Befund; Zeile nur
  aus kurzen Wörtern, deren letztes über die Grenze ragt ⇒ Befund; ein Tab trennt Wörter wie ein
  Space; ein geschütztes Leerzeichen (U+00A0) trennt nicht; letzte Zeile `…\r` ohne `\n` mit
  genau `max-chars` Zeichen ⇒ befundfrei ([ADR-0114](../../adr/0114-linelength-modul.md) Entscheid 3).
- Zeile über der Grenze mit `<!-- d-check:ignore -->` bei einer Grenze von mindestens 14 ⇒ Befund
  (Entscheid 9); bei `max-chars: 10` und sonst nur kurzen Wörtern ⇒ kein Befund, (c) nimmt die
  Zeile über das Marker-Wort aus (Review R1 zu slice-289, F-8).
- Form-Präfix (Review R1 zu slice-289, F-2): bei `max-chars: 120` eine Zeile aus `- ` und einem Wort
  mit 119 Zeichen ⇒ Befund; dasselbe mit `> ` ⇒ Befund; `- foo ` plus ein Wort mit 119 Zeichen ⇒
  Befund, die Zeile ist vor dem Wort umbrechbar (Review R2 zu slice-289, F-3); `-   ` plus ein Wort
  mit 117 Zeichen ⇒ Befund (Rand mehrerer Leerzeichen hinter dem Marker, [ADR-0114](../../adr/0114-linelength-modul.md)).
- Leere Markdown-Datei als einzige Datei ⇒ kein Befund, Bezugsmenge `1 Datei(en), 0 Zeile(n)`
  (Review R1 zu slice-289, F-1).
- Bezugsmenge: Datei mit und ohne abschließendes `\n` ⇒ dieselbe Zeilenzahl wie bei `file`;
  die über (c) ausgenommenen Zeilen stehen als eigene Zahl (Verdikt F-7).
- Aus der Verifikation zu slice-289 (V-2), je ein Akzeptanzkriterium der Anforderung: ein Wort, das
  nur mit seiner Einrückung über die Grenze reicht ⇒ Befund; eine überlange ATX-Überschrift ⇒
  Befund; eine Datei unter `linelength.exempt-paths` und eine außerhalb von `linelength.scope` ⇒
  nicht gemessen; ohne `linelength`-Block und ohne `--enable linelength` ⇒ byte-identische Ausgabe
  zu einem Lauf ohne das Modul. Dazu als eigene Testfälle, nicht nur über die DoD: (b) aus, die
  Typfehler `"80"` für `max-chars` und `"nein"` für einen Schalter ⇒ Exit 2, und der Ausgabekanal
  der Bezugsmenge (stderr im Default, `summary.notes` unter `--json`/`--yaml`).
- Zwei Dateien mit Befunden ⇒ Reihenfolge Pfad, dann Zeile; zweimal derselbe Lauf ⇒
  byte-gleiche Ausgabe.
- Keine Markdown-Datei in der Menge ⇒ kein Befund, Exit 0, Bezugsmenge `0`; explizit leere
  `linelength.scope.roots` ⇒ prüft nichts, kein Fehler (Entscheid 8,
  [`DC-FA-CONF-002`](../../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope));
  Nicht-Markdown-Datei in `scan.roots` mit langer Zeile ⇒ kein Befund.

**Plan-Änderungs-Notiz — Spiegel-Suche vor dem ersten Edit (2026-11-23, Implementer).**
Kommando aus der Repo-Wurzel:
`git grep -l pinsync -- . ':!docs/plan/planning/done' ':!docs/reviews' ':!docs/plan/planning/observations'`
⇒ 30 Dateien. Der Grund-Code steht in `diagnose.go` als Konstante (`ReasonPinsync…`) und fällt
deshalb aus diesem Muster; die Tabelle führt die Datei schon. Zuordnung:

- **In der Tabelle oben** (14): `model/config.go`, `configyaml.go`, `rules/run.go`,
  `config_template.go`, `spezifikation.md`, `lastenheft.md`, `README.md`, `README.de.md`,
  `benutzerhandbuch.md`, `operations.md`, `abdeckung-tests.md` sowie Regel-Datei und Testdateien
  des Moduls (`rules/pinsync.go`, `rules/pinsync_test.go`, `configyaml/pinsync_test.go` als Muster
  der neuen Dateien). Im Handbuch heißt das: ein Modul-Abschnitt §4.28 und die Zeile in §6;
  Kopf und §11 bleiben Release-Prep.
- **Neu in den Plan** (zwei Testdateien, beide `update`):
  `internal/adapter/driving/cli/cli_acceptance_test.go` — die Vorlage von `--print-config`
  (der Block `linelength` entkommentiert ist gültig), `--enable linelength` am echten Lauf
  (Bezugsmenge auf stderr und in `summary.notes`, Exit 2 ohne `max-chars`) und der byte-identische
  Lauf ohne Modul; `internal/hexagon/core/rules/lexikon_kopplung_test.go` — der Kopplungs-Test aus
  DoD-Punkt 2 hängt `linelength` als Konsument an die Tabellenzeilen-Frage und stellt die
  Fence-Frage („liegt die Zeile außerhalb von Fenced-Code?“) an `linelength`, `targets` und
  `structure` mit derselben Eingabe.
- **Bleibt unberührt** (14), je mit Grund: `CHANGELOG.md` und `packaging/dockerhub/overview.md`
  (Modul-Liste ohne Spiegel-Test; Release-Prep, wie bei `pinsync` in der Release-Prep von v0.91.0);
  `.d-check.yml`, `Makefile` (`doc-check`-Kommentar, `FOCUS_DISABLE`),
  `configyaml/gate_consistency_test.go`, `harness/README.md`, `harness/sensors/doc-check.md`
  (das eigene Profil — sie spiegeln dessen `modules:`, nicht `ValidModules()`, und das Modul wird
  dort nicht eingeschaltet, §1); `tools/image-test.sh`, `harness/sensors/image-test.md`,
  `docs/user/abdeckung-e2e.md` (der Container-Lauf von `pinsync`; dieser Plan verlangt keinen
  Container-Fall für `linelength`); [ADR-0113](../../adr/0113-pinsync-modul.md), [ADR-0114](../../adr/0114-linelength-modul.md) und der ADR-Index (Architect);
  Roadmap und diese Datei (Planung).

**Plan-Änderungs-Notiz — zwei Festlegungen am Konfigurations-Rand.** `max-chars` gilt nur als
YAML-Ganzzahl (Tag `!!int`), die drei Schalter nur als `true`/`false` (Tag `!!bool`): der
Dekoder kürzte `80.5` sonst still auf `80` und läse ein quotiertes `"yes"` als `true` — beides
widerspräche „keine ganze Zahl“ bzw. „kein Wahrheitswert“ in
[`DC-FA-LLEN-001`](../../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in).
~~Ein leerer Wert (`max-chars:`) gilt wie ein fehlender.~~ **Zurückgenommen** (Planner, nach
[Code-Review R1](../../../reviews/2026-11-23-slice-290-code-r1.md) F-1 und dem Verdikt des Architect
in [ADR-0114](../../adr/0114-linelength-modul.md) Entscheid 4): Ein Schlüssel des Blocks ohne Wert
(YAML-`null`) ist gesetzt und kein gültiger Wert ⇒ Exit 2, für `max-chars` und die drei Schalter,
auch ohne aktives Modul. Die Notiz stand gegen das Lastenheft. Die Wort-Zerlegung für (c) ist kein
regulärer Ausdruck, aber ein Erkennungs-Muster: ihre Negativfälle (Tab, U+00A0, Einrückung,
Marker-Kommentar, `\r` am Ende) stehen in der Testfall-Liste oben, und der Lauf über den eigenen
Korpus ist die Gegenprobe gegen die `awk`-Zahlen aus dem Plan von
[`slice-289`](../done/slice-289-linelength-anforderung-und-adr.md) (Stand `f3dc94da`, `max-chars: 120`:
1237 Dateien, 513 über (c) ausgenommen, 2439 Befunde). Kommando: `git archive f3dc94da` in ein
leeres Verzeichnis entpackt, dort eine Konfigurationsdatei mit dem `scan`-Block des damaligen
Profils, `modules: []` und `linelength: {max-chars: 120}`, dann das mit `make build` gebaute Image
mit `--network none`, `--config <datei>` und `--enable linelength`; verglichen werden die
Bezugsmenge auf stderr und die Zahl der Befund-Zeilen auf stdout.

**Plan-Änderungs-Notiz — eine Handbuch-Aufzählung ohne Spiegel-Test.** Im Handbuch zählt §5
(Ventile, Bullet `exempt-paths`) die Module mit `exempt-paths` auf und schließt mit „Die übrigen
Module haben keinen der beiden Schlüssel“ — mit dem registrierten Modul wäre der Satz falsch.
`linelength` kommt dort in die Liste; die Aufzählung der Module, die `d-check:ignore` kennen,
bleibt bei vier. Gefunden mit
``git grep -n 'codepaths`, `diagrams`, `workflows`' -- spec docs/user README.md README.de.md harness``
⇒ ein Treffer (`benutzerhandbuch.md`).

**Plan-Änderungs-Notiz — Fix nach [Code-Review R1](../../../reviews/2026-11-23-slice-290-code-r1.md)
(F-1 bis F-4, F-7; Implementer).** Was der Fix-Commit tut, je Befund, mit den Dateien:

- **F-1** (Verdikt in [ADR-0114](../../adr/0114-linelength-modul.md) Entscheid 4): Ein Schlüssel des
  Blocks ohne Wert ist gesetzt und ungültig ⇒ Exit 2, für `max-chars` und die drei Schalter, auch
  ohne aktives Modul; abwesend ist nur ein fehlender Schlüssel. `configyaml.go` bekommt dafür eine
  eigene Abwesenheits-Prüfung für `linelength`, die übrigen Module bleiben unverändert. Tests in
  `configyaml/linelength_test.go`, dazu die Gegenprobe, dass bei `versions` (`current-from:`) und
  `planning` (`marker:`) ein Schlüssel ohne Wert abwesend bleibt. Die Spezifikation
  ([`DC-FA-LLEN-001.a`](../../../../spec/spezifikation.md#dc-fa-llen-001a--zeilenlänge-in-markdown-dateien-linelength)
  Schritt 1, Historie) und das Handbuch §4.28 bekommen den Satz des Architect; das Lastenheft
  bleibt unverändert.
- **F-2:** gewählt ist die **gemeinsame Funktion**, nicht nur ein Test: `pinsyncLines` wird zu
  `textLines` und steht neben `countLines` in `rules/codepaths.go`; `pinsync` und `linelength`
  zerlegen beide damit (neu im Plan: `rules/pinsync.go`, `rules/codepaths.go`). Grund: [ADR-0114](../../adr/0114-linelength-modul.md)
  Entscheid 3 sagt „dieselbe Zerlegung“ — eine Funktion macht das zum Bau, ein Test nur zur
  Beobachtung. `countLines` bleibt, wie es ist, weil `codepaths`, `citations` und `file` es
  nutzen und dieser Slice deren Code nicht ändert; ein Test in `rules/linelength_test.go` hält
  `len(textLines(x)) == countLines(x)` über die Randformen.
- **F-3:** ein Test für `\r\r\n` und für eine letzte Zeile auf `\r\r`: nur das letzte `\r` fällt
  weg (`rules/linelength_test.go`).
- **F-4:** die Fence-Kopplung fragt zusätzlich `planning.waves` (Registerzeile, vorhandene
  Probe), `planning` (Ruhe-Marker im Block der Überschrift) und `vcs` (`baseImmutable` und
  `vcsHeadStatusLineNo`). Gezählt mit
  `git grep -n 'proseLineSet(content)' -- internal/hexagon/core/rules ':!*_test.go'` ⇒ 8 Aufrufe
  in fünf Modulen. Nicht gefragt wird `waveBijection` (`planning_waves.go`, Kennungen im Block):
  dasselbe Modul und dieselbe Prosa-Menge wie die Registerzeile, aber eine Probe bräuchte die
  flachen Wellendokumente samt Verzeichnis — der Gewinn ist eine zweite Lesestelle desselben
  Moduls, keine neue Antwort.
- **F-7:** der Klartext nennt bei einem nicht skalaren Wert dessen Art (Liste, Abbildung, Alias)
  statt eines leeren Zitats; die Grenzen von Schritt 1 sagen, dass der YAML-Typ zählt, nicht der
  Zahlwert (`80.0` ⇒ Exit 2).

Die Suche nach der Aussage jedes Befunds (R1): außerhalb von `linelength` trug keiner der Treffer
dieselbe Aussage. F-1 (`!!null`, „ohne Wert“): sonst nur die eigenen Null-Prüfungen von `trace`
und `targets`, `versions` (`pin-pattern:`), `reviews.promise-pattern` und ein Flag ohne Wert in
der CLI. F-2 (Zerlegung am `\n`): `immutable`, `matrix`, `pins`, `sources`, `versions`,
`markdown`, `repair`, `trace_table` und `vcs` zerlegen roh, ohne `\r`-Regel, und behaupten keine
Gleichheit mit `pinsync`. F-3 (`TrimSuffix … "\r"`): sonst `metrics_table.go`, eine eigene
Tabellenzeile. F-4: [ADR-0114](../../adr/0114-linelength-modul.md) und §2/§6 dieses Plans. F-7: die zwei Meldungen in `configyaml.go`.
Die Treffer für R2 stehen in der folgenden Notiz.

**Plan-Änderungs-Notiz — Fix nach [Code-Review R2](../../../reviews/2026-11-23-slice-290-code-r2.md)
(F-1 bis F-4; Implementer).** Was der Fix-Commit tut, je Befund:

- **F-1** (Verdikt in [ADR-0114](../../adr/0114-linelength-modul.md) Entscheid 4, eingeengt auf die
  vier eigenen Werte): Am Verhalten ändert sich nichts. Der Code prüft schon nur `max-chars` und die
  drei Schalter; `scope:` und `exempt-paths:` ohne Wert sind abwesend, wie bei jedem scannenden
  Modul. Eingeengt werden die Sätze, die mehr sagen:
  [`DC-FA-LLEN-001.a`](../../../../spec/spezifikation.md#dc-fa-llen-001a--zeilenlänge-in-markdown-dateien-linelength)
  Schritt 1 (Wortlaut des Architect, am Code geprüft) mit Historie-Zeile, Handbuch §4.28, der
  Kommentar an `applyLinelength`/`linelengthSet` und der Klartext („ein Schlüssel ohne Wert ist
  gesetzt“ wird zu „gesetzt, aber kein gültiger Wert“). Das Lastenheft trägt keine pauschale Form
  (`grep -c 'ohne Wert' spec/lastenheft.md` ⇒ 2, beide zu `reviews.promise-pattern`)
  und bleibt unverändert. Neue Tests in `configyaml/linelength_test.go` und
  `rules/linelength_test.go`: `scope:` ohne Wert ⇒ dieselbe Menge wie ohne `scope`;
  `exempt-paths:` ohne Wert ⇒ keine Datei ausgenommen; Gegenprobe `ids` mit `scope:` ⇒ abwesend.
- **F-2:** gewählt ist **Auflösen**: ein Alias wird auf seinen Zielknoten aufgelöst und dieser
  geprüft — `max-chars: *g` auf eine ganze Zahl ist gültig, ein Alias auf eine Liste meldet
  „eine Liste“, einer auf `null` ist ohne Wert. Grund: ein Alias steht in YAML für denselben Wert,
  und der Klartext soll die wirkliche Ursache nennen; dieselbe Auflösung machen `targets.authority`
  und `trace.requirements` am Config-Rand. Die Grenzen von Schritt 1 sagen es.
- **F-3:** die Gegenprobe bekommt `structure` (`max-lines:` ohne Wert ⇒ abwesend).
- **F-4:** der Satz oben verweist jetzt auf Repo-Artefakte statt auf den Bericht.

Treffer der Suche nach der Aussage (R2), vor dem Fix-Commit:

- F-1, „Schlüssel des Blocks ohne Wert“/„Schlüssel ohne Wert“ über `spec docs/user internal
  harness docs/plan`: pauschal sind `spezifikation.md` Schritt 1 und Historie-Zeile, Handbuch
  §4.28, der Kommentar an `applyLinelength`, der Klartext in `linelengthSet` und der Kommentar an
  `TestDecode_LinelengthSchluesselOhneWert` — alle eingeengt. Die ADR hat der Architect schon
  eingeengt; `configyaml.go` (versions, `pin-pattern:`) und `spezifikation.md` (`patterns:`)
  sagen anderes und bleiben.
- F-2, „ein Alias“: nur `yamlValueKind` in `configyaml.go`.
- F-3, „`versions`, `planning` und `structure`“: Spezifikation, Handbuch, [ADR-0114](../../adr/0114-linelength-modul.md) und der
  Kommentar an `linelengthSet` — die Aussage stimmt am Code (`rawStructure` mit `*int`), der Test
  kommt dazu.
- F-4, „Bericht des Implementers“: nur dieser Plan.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`slice-289`](../done/slice-289-linelength-anforderung-und-adr.md)
liegt in `done/` (Anforderung, Spezifikation und ADR stehen, die Entscheide sind bestätigt),
und das WIP-Limit ist frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Spiegel-Suche ergibt so viele
  Stellen, dass Regel und Spiegel nicht in einer Review-Sitzung prüfbar sind — dann wird
  der Spiegel-Nachzug ein eigener Slice derselben Welle.
- `in-progress` → `open` (blockiert): die geteilte Lexik liefert für die Ausnahmen (a) oder
  (b) eine Antwort, die der Anforderung widerspricht, und eine Änderung der Lexik verschöbe
  andere Module — dann geht die Frage an die Anforderung zurück.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (Modul registriert, Tests je Randbedingung mit gezeigtem Rot, Spiegel grün,
`make gates` grün), Review- und Verifikations-Report liegen unter `docs/reviews/`,
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Ein nicht geschlossener Fence nimmt über die geteilte Lexik den Rest der Datei aus der
  Prüfung — ein stilles Grün über einem Teil der Datei
  ([`BEO-ALL/prosa-verschwindet-in-inline-code-spanne`](../observations/BEO-ALL/prosa-verschwindet-in-inline-code-spanne/observation.md),
  Fence-Hälfte, 2×) — **Ausgang:** *entfallen* — `TestLinelengthOffenerZaunNimmtDenRestAus`
  belegt das in der Anforderung festgelegte Verhalten (der Rest der Datei ist ausgenommen,
  `linelength` meldet kein `fence-unclosed`), und die Anforderung verweist für die Meldung auf
  `spans`; die [Verifikation](../../../reviews/2026-11-23-slice-290-verifikation.md) zeigt das Rot
  der Mutation M04. Kein Beleg in `BEO-ALL/prosa-verschwindet-in-inline-code-spanne`.
- `proseLineSet`/`tableRowLine` sind paket-interne Helfer anderer Module; sie für ein neues
  Modul zu nutzen bindet dessen Verhalten an deren Änderungen
  ([`BEO-ALL/shared-lexicon-drifts-at-edges`](../observations/BEO-ALL/shared-lexicon-drifts-at-edges/observation.md),
  verkörpert über Kopplungs-Tests; Pflicht nach [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 4, DoD-Punkt 2) —
  **Ausgang:** *entfallen* — `TestTabellenzeilenFrageHatEineAntwort` (vier Konsumenten) und
  `TestFencedCodeFrageHatEineAntwort` (sieben der acht Aufrufe von `proseLineSet`, der achte
  begründet ausgenommen) fahren dieselbe Eingabe durch `linelength` und die bisherigen
  Konsumenten; `proseLineSet` und `tableRowLine` blieben unverändert (Code-Review
  [R1](../../../reviews/2026-11-23-slice-290-code-r1.md) F-4,
  [R2](../../../reviews/2026-11-23-slice-290-code-r2.md)).
- Eine Randbedingung der Anforderung bleibt ohne Test
  ([`BEO-ALL/spec-randbedingung-ohne-test`](../observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md),
  1×) — **Ausgang:** *entfallen* — jedes der 17 Akzeptanzkriterien hat einen Test mit gezeigtem
  Rot; die Zuordnung Kriterium → Test → Mutation steht im
  [Verifikationsbericht](../../../reviews/2026-11-23-slice-290-verifikation.md) (24 Mutationen,
  alle rot), nicht wie hier vorgesehen im Plan-Nachtrag — dort standen die Rot-Belege nur als
  Behauptung in den Commit-Botschaften (Verifikation V-1). Kein Beleg in
  `BEO-ALL/spec-randbedingung-ohne-test`.

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

- **Was hat funktioniert:** Die Anforderung aus slice-289 trug den Code: Die Testfall-Liste in
  §3 übersetzte sich eins zu eins in Tests, und die Korpus-Gegenprobe am gebauten Image traf die
  `awk`-Zahlen aus slice-289 exakt (1237 Dateien, 513 ausgenommen, 2439 Befunde am Stand
  `f3dc94da`). Die geteilte Lexik blieb unverändert; das Modul liest `proseLineSet` und
  `tableRowLine` und baut nichts nach, die Kopplungs-Tests halten es fest. Implementer,
  Reviewer und Verifier haben das Rot je Testfall unabhängig voneinander gezeigt (30, 20 und 24
  Mutationen).
- **Was ging anders als geplant:** **(1)** Eine Plan-Notiz des Implementers („ein leerer Wert
  gilt wie ein fehlender“) stand gegen das Lastenheft, der Code folgte ihr. Code-Review R1 F-1
  fand es, der Architect entschied V1 ([ADR-0114](../../adr/0114-linelength-modul.md)
  Entscheid 4), die Notiz ist zurückgenommen. Review R2 fand dann die Reichweite zu pauschal
  formuliert; der Architect begrenzte Entscheid 4 auf die vier eigenen Werte. Das ist das
  zweite Auftreten dieser Klasse in welle-95 (Register unten). **(2)** Die Zeilen-Zerlegung
  ist jetzt eine geteilte Funktion (`textLines`, R1 F-2) statt einer Kopie bei `pinsync`. **(3)**
  `make slice-mv` scheiterte am eigenen pre-commit-Hook; die Übergänge liefen von Hand.
  **(4)** Die Zuordnung Testfall → Rot stand nicht im Plan-Nachtrag, sondern nur behauptet in den
  Commit-Botschaften; der Verifier hat sie nachgeliefert (V-1).
- **Steering-Loop-Eintrag:** mit diesem Slice ist keine Regel, kein Sensor und keine benannte
  Spec-Lücke verkörpert. Die Klasse „Plan-Festlegung gegen das Lastenheft“ steht bei 2× im
  Register; ein drittes Auftreten ist die Schwelle für den Lese-Schritt der Welle-Closure.
- **Beobachtungs-Register (`../observations/`):** neu
  [`BEO-ALL/plan-entscheid-steht-gegen-lastenheft`](../observations/BEO-ALL/plan-entscheid-steht-gegen-lastenheft/observation.md)
  (2×: slice-289 und slice-290),
  [`BEO-ALL/yaml-dekoder-kuerzt-zahlen-still`](../observations/BEO-ALL/yaml-dekoder-kuerzt-zahlen-still/observation.md)
  (Code-Review R1 F-5),
  [`BEO-ALL/regeltext-gegen-spiegel-test`](../observations/BEO-ALL/regeltext-gegen-spiegel-test/observation.md)
  (R1 F-6: [`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17 gegen den Spiegel-Test der READMEs)
  und
  [`BEO-HARN/slice-mv-scheitert-am-pre-commit-doc-check`](../observations/BEO-HARN/slice-mv-scheitert-am-pre-commit-doc-check/observation.md),
  je 1× außer dem ersten. Die Zähler folgen aus den Dateien.
- **Akzeptierte Negative:** Verifikation V-2 — der Reihenfolge-Test von `linelength` bleibt
  grün, wenn `SortFindings` wegfällt, weil die Reihenfolge schon aus Datei- und Zeilen-Reihenfolge
  entsteht; den Ausfall der gemeinsamen Sortierung fangen acht Tests anderer Module. V-3 und V-4
  (INFO) ohne Träger: das Rot von M11 als nil-Panic ist der reale Vorzustand, die Meldung zum
  unbekannten Schlüssel ist die geteilte strikte Dekodierung.
- **Folge-Slices:** keine. Release-Prep (CHANGELOG, Handbuch-Kopf, Image-Tags) ist ein eigener
  Vorgang nach der Welle.
- **Risiken aus §6:** alle drei *entfallen* — siehe §6.
- **Trigger-Audit:** läuft mit der Closure von welle-95; dieser Slice setzt keinen Carveout und
  kein bootstrap-aware Gate. [ADR-0114](../../adr/0114-linelength-modul.md) ist mit ihm
  `Accepted` geworden; ihre Re-Evaluierungs-Trigger (a)–(f) sind nicht eingetreten.
- **Drei Paarungen:** getragen von der Closure von welle-95 — kein Anker (nichts verkörpert),
  kein Folge-Slice, Register die vier Einträge oben.
- **Reviews und Verifikation:** [Code-Review R1](../../../reviews/2026-11-23-slice-290-code-r1.md)
  (1 MEDIUM, 3 LOW, 3 INFO — MEDIUM und LOW behoben, F-5/F-6 ins Register),
  [Code-Review R2](../../../reviews/2026-11-23-slice-290-code-r2.md) über die Fix-Commits
  (1 MEDIUM, 1 LOW, 2 INFO — behoben),
  [Verifikation](../../../reviews/2026-11-23-slice-290-verifikation.md) (DoD 1–4 bestätigt,
  Probe am Image, Byte-Identität gegen v0.91.0 über 20 Vergleiche).

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
Regelmodul, Konfigurations-Rand, CLI-Vorlage, Spiegel-Tests; keine Sub-Area mit eigener
Modus-Deklaration (`tools/harness/` bleibt unberührt).

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand
2026-11-22), Sub-Area `*`. Treffer:
[`BEO-ALL/modulliste-spiegel-ungegated`](../observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md)
(verkörpert in den Spiegel-Deckungstests — die DoD verlangt, dass sie das neue Modul
fangen),
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/observation.md)
(verkörpert in Schritt 19 des Workflow-Skeletts — die Testfall-Liste in §3 ist seine
Anwendung),
[`BEO-ALL/gemeinsame-leser-sind-zwei-familien`](../observations/BEO-ALL/gemeinsame-leser-sind-zwei-familien/observation.md)
(offen, 2× — ein eigenes Fence-Muster wäre der dritte Beleg),
[`BEO-ALL/prosa-verschwindet-in-inline-code-spanne`](../observations/BEO-ALL/prosa-verschwindet-in-inline-code-spanne/observation.md)
(offen, 2× — Risiko 1 in §6),
[`BEO-ALL/shared-lexicon-drifts-at-edges`](../observations/BEO-ALL/shared-lexicon-drifts-at-edges/observation.md)
(verkörpert über Kopplungs-Tests — Risiko 2 in §6),
[`BEO-ALL/spec-randbedingung-ohne-test`](../observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md)
(offen, 1× — Risiko 3 in §6) und
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
(verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.8 — das Modul liest nur seine
Scan-Menge). Der Zähler folgt aus den Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-11-22: **beide rot** — `upstream-drift` mit fünf veralteten
Pins (Go, Semgrep, a-check, `docker/login-action`, Digest `golang:1.27.3`), `image-scan` mit
behebbaren CRITICAL/HIGH im publizierten Image; Einzelheiten im Plan von
[`slice-289`](../done/slice-289-linelength-anforderung-und-adr.md) §8. Berührt das Modul nicht;
eigene Pin-Hebung vor dem nächsten Release. Beim Start dieses Slice neu lesen.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
