# Slice slice-283: Modulkern `metrics` — Block lesen, Obergrenze gegen Formel, leerer Grund (hermetisch)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-93

**Bezug:** die Anforderung des Moduls `metrics` aus [`slice-282`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) (Kennung dort vergeben; dieser Slice trägt
die Hälften „Obergrenze gegen Formel“ und „leerer Grund“),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus) (Byte-Identität ohne `metrics`-Konfiguration),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) (hermetisch: kein git, kein Netz). Die begleitende ADR aus
[`slice-282`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) ist die aktive Entscheidungsgrundlage.

**Berührte Spec-Stellen:** die in [`slice-282`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) angelegten Stellen
(Anforderung, Schema-Zeilen `metrics.*`, Grund-Codes) — gelesen, nicht verändert,
außer dem Fassungsstand der Historie.
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-30.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Bei aktiviertem Modul `metrics` liest d-check aus jedem Plan unter dem
konfigurierten Verzeichnis den Block, hält die darin genannte Obergrenze gegen die
konfigurierte Formel über den Zahlenfeldern des Blocks (eine
Aufrundung, Schwellenregel) und meldet einen Ausnahme-Eintrag
mit leerem Grund als Befund — rein arithmetisch, ohne git, einzeln lieferbar.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Config-Schlüssel `estimate-inputs`, `table.*` und `excluded`** — das Schema der Spezifikation
  kommt stückweise: `excluded` mit der Messung ([`slice-284`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle)), `estimate-inputs` und `table.*` mit der
  Tabelle ([`slice-285`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle)). Bis dahin kennt der Config-Rand sie nicht (unbekannter Schlüssel, Exit 2);
  der Hinweis „Marke ausstehend“ nennt hier noch keine Zahlen und prüft nicht auf die Eingaben der
  Schätzung.
- **Die Tabellenzählung (Schätzung aus der Tabelle, Prüfung 2, Hinweis „Marke ausstehend“)**
  — Folge-Slice: [`slice-285`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) trägt sie, weil Prüfung 2 den Zustand „Messung 0 / ≥ 1“
  aus der Messung braucht. Hier ist `schaetzung` ein Zahlenfeld des Blocks, kein Ergebnis
  einer Tabelle. Die Prosa-Form-Erkennung ist nach der Antwort des Issue-Autors kein
  Zustand des Moduls: ein Plan ohne Block gibt nur den Hinweis „kein Block“, der den Plan
  nennt (Urteil gleich, nicht gemessen).
- **Die Messung aus git und der Vergleich mit der Obergrenze** — Schicht-Abgrenzung
  und Folge-Slice: [`slice-284`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) trägt VCS-Port, Adapter und den Vergleich;
  dieser Slice bleibt hermetisch (kein `.git` nötig), damit seine Hälfte ohne die andere
  prüfbar ist.
- **Release-Prep: `CHANGELOG.md`, Handbuch-Kopf, Digest-Pins, Prosa über die
  Modul-Aufzählung hinaus** — eigener Vorgang nach der Welle
  ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17). **Ausgenommen sind die Orte, die
  `registry_mirror_test.go` in beide Richtungen gegen `model.ValidModules()` hält** —
  Modulsatz in [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl), Glossar-Zeile „Regelmodul“, Optionen-Tabelle in
  `operations.md`, Modulliste in `README.md` und `README.de.md`, §6-Tabelle des
  Handbuchs. Sie gehören hierher, weil `make gates` sonst ab der Registrierung rot ist;
  mehr als die Aufzählung schreibt dieser Slice dort nicht. (Der Test gewinnt hier
  gegen die Prosa-Regel der Zeile 17: das gated Gebot ist das engere.)
- **Das Aktivieren des Moduls im eigenen `.d-check.yml` mit einem Block** — dieses Repo
  führt keinen Block in seinen Plänen; ein künstlicher Block wäre eine Selbstbehauptung
  ohne Gegenstand. Die `modules:`-Liste bleibt damit unverändert, und
  `FOCUS_DISABLE` im `Makefile` — das diese Liste spiegelt und von
  `TestFocusDisable_DecktDCheckYmlModules` dagegen gehalten wird — bleibt es auch; dass
  der Test ein nicht aufgenommenes opt-in-Modul durchlässt, ist vor dem Code am
  Test selbst nachzulesen, nicht hier behauptet.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Registrierung und Konfiguration:** das Modul steht in `model.ValidModules()`,
      `metrics.*` wird gelesen und validiert (Config-Rand: fehlende Pflichtschlüssel, Wertebereiche,
      Teiler 0, doppelte Feldnamen, unbekannter Schlüssel → Exit 2), das Gerüst (`--print-config`) führt den Block; der
      Spiegel-Test (`registry_mirror_test.go`) ist grün, weil alle von ihm gehaltenen Orte
      (§1) `metrics` nennen. `gate_consistency_test.go` bleibt unverändert: das Modul
      steht nicht in der `modules:`-Liste von `.d-check.yml`, die dritte Richtung des
      Tests greift nicht; ob `metrics` später in `forbiddenInNetless()` gehört (es
      braucht `.git`), entscheidet die Aufnahme in diese Liste.
- [x] **Prüfung „Obergrenze = Formel“:** pro Plan mit Block wird aus dem Schätzungs-Feld (`E = Konstante + Σ Faktor · Feld`)
      die Obergrenze mit genau einer ganzzahligen Aufrundung berechnet und gegen die
      genannte verglichen (`metrics-limit-mismatch`); die Schwellenregel meldet ein
      fehlendes Pflichtfeld (`metrics-cut-reason-missing`). Ein formal ungültiger Block
      (unvollständig, Zahl ungültig, Fence ungeschlossen) ist ein **Eingabefehler**:
      Exit 2 mit Klartext, kein Befund, und der Lauf bricht beim ersten Fehler ab.
      Test referenziert, Negativliste (§3) getestet; die Testfälle sind die belegten
      Zahlen der Anforderung (E = 6790 ⇒ 9200; Grenzfälle E = 2000 ⇒ 2700, 2001 ⇒ 2800,
      2222 ⇒ 3000, 2223 ⇒ 3100 samt Schwellenregel; 5036 ⇒ 6800; 5430 ⇒ 7400; 9215 ⇒ 12500; 4670 ⇒ 6400); unbekannte Block-Felder werden ignoriert.
- [x] **Prüfung „Eintrag mit leerem Grund“:** ein Ausnahme-Eintrag `<N> | <Grund>` mit
      leerem oder nur-Leerraum-Grund ist ein Befund (`metrics-reason-empty`),
      **unabhängig** von einer Messung; Plan-Suche nur unmittelbar im konfigurierten
      Verzeichnis (`slice-*.md`, ohne Datei grün); Test referenziert.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update für die Spiegel, die ein Test hält (§1, nur die Aufzählung);
- [x] Grund-Codes dieses Slice (`metrics-limit-mismatch`, `metrics-reason-empty`, `metrics-cut-reason-missing`): Zeile in der Tabelle von `spec/spezifikation.md` §4, Eintrag in `AllReasons()` und `--doctor`-Klartext (`TestAllReasonsDeckungGegenSpezifikationGrundCodes`); `internal/hexagon/core/model/finding.go` und `internal/hexagon/core/app/diagnose.go` (`AllReasons()`) kommen in die Plan-Tabelle (§3).
      `CHANGELOG.md`, Handbuch-Kopf und Digest-Pins **nicht** (Release-Prep).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `metrics.go` im Regel-Paket | neu | Block aus dem Plan-Text lesen (Fence mit konfiguriertem Info-String, Zeilen-Syntax ohne YAML, auf `metrics` begrenzte Fence-Ausnahme — Entscheide 6/7 der ADR), Formel ganzzahlig auswerten (Summe über Block-Felder, eine Aufrundung, Schwellenregel), Befunde `metrics-limit-mismatch`, `metrics-cut-reason-missing`, `metrics-reason-empty`; Eingabefehler als Exit 2 |
| `internal/hexagon/core/model/finding.go` · `internal/hexagon/core/app/diagnose.go` | update | Grund-Code-Konstanten und `AllReasons()` samt Klartext (Spiegel zur Tabelle von `spec/spezifikation.md` §4) |
| `internal/hexagon/core/rules/run.go` | update | Modul in die Regel-Kette einhängen |
| `internal/hexagon/core/model/config.go` | update | `MetricsConfig`, Eintrag in `ValidModules()` |
| `internal/adapter/driven/configyaml/configyaml.go` | update | `metrics.*` lesen und validieren |
| `internal/adapter/driving/cli/config_template.go` | update | Gerüst führt den Block (Spiegel) |
| `metrics_test.go` im Regel-Paket | neu | Happy/Boundary/Negative nach der Anforderung; Negativliste unten |
| `internal/adapter/driven/configyaml/metrics_test.go` | neu | Config-Rand → Exit 2 (eigene Datei neben `mentions_test.go`, statt `configyaml_test.go` zu verlängern); `docs/user/abdeckung-tests.md` wird mit `make abdeckung` neu erzeugt |
| `internal/hexagon/core/app/registry_mirror_test.go` | update (falls nötig) | Spiegel-Deckung ([`BEO-ALL/modulliste-spiegel-ungegated`](../observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md)); der Test liest die Menge aus `ValidModules()`, eine Änderung ist nur bei einem neuen Ort nötig |
| `spec/lastenheft.md` | update | Modulsatz in [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl), §6 Glossar-Zeile „Regelmodul“, Fassungsstand und Historie; die Anforderung selbst kommt aus [`slice-282`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) |
| `spec/spezifikation.md` | update | nur Fassungsstand und Historie |
| `docs/user/operations.md` · `README.md` · `README.de.md` · `docs/user/benutzerhandbuch.md` | update | nur die Modul-Aufzählung, die der Spiegel-Test hält (§1); Handbuch ohne Slice-Verweise |

**Zählung der Spiegel (Schritt 17/18 des Workflow-Skeletts), vor dem Code
nachzuzählen:** `grep -rlE 'mentions' --include=*.go --include=*.md --include=*.yml
internal spec docs/user harness .claude .d-check.yml Makefile README.md README.de.md | grep -v '_test.go\|done/\|docs/reviews'`
als Vorbild-Menge für ein jüngeres Modul (16 Treffer bei der Prüfung dieses Plans, davon
14 ohne die beiden READMEs); jede
Datei dieser Menge ist für `metrics` zu beantworten — berührt, Release-Prep oder
begründet unberührt. `docs/user/abdeckung-tests.md` ist von `make abdeckung` abgeleitet
(`make test` hält sie gegen ihre Ableitung): sie wird neu erzeugt, nicht von Hand
geschrieben, und neue Testdateien verschieben sie. `harness/README.md` und `harness/sensors/` bekommen **kein**
Target (kein `make`-Ziel in diesem Slice; Gate-Index unverändert).

**Negativliste der Block-Erkennung (Schritt 19, vor dem Code-Commit zu testen):**
Block mit anderem Info-String · zwei Blöcke im selben Plan · Block in einer
Einrückung/Zitat-Ebene · Block im Fenced-Block eines Fensters (Fence in Fence) ·
Block ohne Schluss-Fence (⇒ Eingabefehler Exit 2, nicht „ohne Block“) · Backtick- und Tilde-Fence · zwei Blöcke im Plan (Eingabefehler) · mehrfacher Schlüssel, der nur einmal stehen darf · Zahl mit führender Null oder Tausendertrennung · Schwellenregel genau auf der Schwelle (E = 2222 ohne, E = 2223 mit Pflichtfeld) · Ausnahme-Eintrag mit N kleiner als die Messung · Pflichtfeld fehlt · Pflichtfeld doppelt · Eingabe negativ
oder nicht ganzzahlig · Überlauf der Ganzzahl-Arithmetik · Grund nur aus Leerraum ·
Plan ohne Block (kein Befund, kein Abbruch) · Plan außerhalb des konfigurierten
Verzeichnisses. Dazu der Lauf über den eigenen Bestand: mit aktiviertem Modul und
leerem Block-Satz **kein** Befund und gegen das Vorher-Image byte-gleich.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`slice-282`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) liegt in `done/` (Anforderung,
Schema und ADR stehen, die Entscheide 1–9 sind gesetzt), und welle-93 ist offen. Der
Übergang landet auf dem Hauptzweig, vor der Arbeit.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Config-Leser, Block-Lexik und
  Spiegel zusammen tragen mehr als eine Review-Sitzung (Plan-Tabelle über zehn Dateien
  im Diff, mehr als zwei Schichten) — dann wandert die Block-Lexik in einen eigenen
  Slice davor, ausdrücklich nicht still weiter. **Wird die Rückführung nicht gezogen**,
  steht im Review-Report nach [`MR-066`](../../../../harness/conventions.md#mr-066) der
  Grund **und** die Ersatz-Form der Prüfung: zwei Review-Durchgänge, Teil A
  Registrierung, Config-Rand und Spiegel, Teil B Block-Lexik, Formel und Grund.
- `in-progress` → `open` (blockiert — Carveout?): die Anforderung aus
  [`slice-282`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) wird beim Implementieren als nicht umsetzbar erkannt (z. B.
  Ganzzahl-Arithmetik oder die einzelne Aufrundung reicht für die Werte des Anlass-Repos nicht).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, Review-Report liegt vor, und am gebauten Image
meldet ein Fixture mit falscher Obergrenze `metrics-limit-mismatch`, eines mit leerem
Grund `metrics-reason-empty`, während der Lauf ohne `metrics`-Konfiguration gegen das
Vorher-Image byte-gleich ist. Dazu die Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die Block-Lexik liest einen Fence, den der Markdown-Leser des Repos anders begrenzt
  (Einrückung, Fence in Fence, Tilde-Fence); ein Block verschwindet still und der Plan
  gilt als „ohne Block“ — **Ausgang:** weiter offen — die Spezifikation legt die Leseweise der RTM-Leser fest (Zaun nach Space und Tab vorn); ob Pläne des Anlass-Repos davon abhängen, zeigt die Paritäts-Probe ([`BEO-ALL/gemeinsame-leser-sind-zwei-familien`](../observations/BEO-ALL/gemeinsame-leser-sind-zwei-familien/observation.md))
- Das Modul liest die Pläne unterhalb des konfigurierten Verzeichnisses, aber das
  Verzeichnis kann außerhalb der Scan-Wurzeln liegen; die Zusage der Scan-Menge
  ([`AGENTS.md`](../../../../AGENTS.md) §3.8) gilt dort nicht — **Ausgang:** weiter offen — die Grenze steht in Anforderung, Spezifikation, Handbuch, Gerüst und Code; die Messung des Verzeichnisses selbst bleibt eine Zusage der Anforderung ([`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md))
- Ein Plan ohne Block gilt als „nicht gemessen“ und meldet nur einen Hinweis; wer den Block
  vergisst, umgeht die Marke vollständig — **Ausgang:** entfallen — der Hinweis „kein Block“ nennt den Plan und steht in der Zusammenfassung; der Issue-Autor bestätigt für Pläne ohne Block dasselbe Urteil (grün, nicht gemessen, Exit 0) wie sein bisheriges Programm

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

- **Geliefert:** das Modul `metrics` in `ValidModules()`, `metrics.*` am Config-Rand, das Gerüst, die Regel
  (`rules/metrics.go`: Block lesen, Zustand bestimmen, Marke gegen Formel, Schwellenregel, leerer Grund,
  Plan-Suche, Eingabefehler als Exit 2), die drei Grund-Code-Zeilen samt Klartext und die Spiegel-Aufzählung
  in Lastenheft, `operations.md`, beiden READMEs und der Handbuch-Tabelle.
- **Gemessen:** `make gates` grün (Coverage 94,80 %). Am gebauten Image mit `--network none` und `:ro`: falsche
  Marke, leerer Grund und fehlendes Pflichtfeld ergeben je einen Befund mit Exit 1, eine Block-Zeile ohne `:` Exit 2,
  korrekte Pläne Exit 0. `make blackbox-probe REF=v0.89.0`: 20 Vergleiche byte-identisch zum Vorher-Image.
  Bewusstes Brechen durch den Verifier: Duplikat-Prüfung und Schwellen-Vergleich ausgebaut, je ein benannter Test rot.
- **Was anders lief als geplant:** (1) Das Config-Schema kommt stückweise (Plan-Änderung vorab): `excluded` mit der
  Messung, `estimate-inputs` und `table.*` mit der Tabelle; ein Zwischenstand wird nicht released. (2) Die Grund-Code-Zeilen
  der Spezifikation kamen erst hier, weil `TestAllReasonsDeckungGegenSpezifikationGrundCodes` sie an den Code bindet.
  (3) Zwei meiner Testfälle lagen über der Schwellenregel und lieferten einen zweiten Befund; die belegten Zahlen
  2222 und 2223 fangen die Mutation `<=` zu `<` nicht, erst ein Test mit Schwelle genau auf der Grenze (Review R1 F-3).
- **[`MR-066`](../../../../harness/conventions.md#mr-066) — Rückführung nicht gezogen:** der Slice lag mit 15 Dateien in vier Schichten über der Ein-Sitzungs-Grenze.
  Grund: Registrierung, Gerüst und Spiegel sind mechanisch und bilden das Modul nicht ohne den Kern ab; die Zerlegung hätte
  einen Slice ergeben, dessen Gates ohne den Kern rot sind (Spiegel-Test, Grund-Code-Test). Ersatz-Form der Prüfung: der
  Review R1 las die Regel und die Spiegel getrennt (Teil Regel: `rules/metrics.go`, `run.go`, Tests; Teil Rand: Config-Leser,
  Gerüst, Klartext, Spiegel), der Verifier fuhr das Image und zwei Mutationen. Ein Re-Review der Fix-Commits entfiel,
  weil der Fix genau die Befunde umsetzt (keine neue Zusage, kein neuer Mechanismus).
- **Steering-Loop-Eintrag:** keiner mit neuer Schwelle.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angefallen, die ein Beleg wäre. Benannt, nicht gezählt: zwei Gleichheits-Mutationen
  blieben ohne Test, bis Review und Verifikation sie fanden (Review R1 F-3); die Negativliste (Schritt 19) war angewandt.
- **Folge-Slices:** [`slice-284`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) und [`slice-285`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) in `open/`.
- **Risiken aus §6:** zwei weiter offen (mit Verweis), eines entfallen (mit Grund). Trigger-Audit: kein Carveout, kein
  bootstrap-aware Gate, keine ADR mit eingetretenem Trigger ([`ADR-0112`](../../adr/0112-metrics-modul.md) bleibt `Proposed`), keine Hard Rule mit
  eingetretenem Trigger, keine Schwelle gesenkt.
- **Drei Paarungen:** (a) Anker — kein `liegt in`-Feld; (b) Folge-Slice — die genannten Slices liegen in `open/`;
  (c) Register — die zitierten Beobachtungen existieren und tragen Belege.

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
Regelmodul, Konfigurations-Rand, Gerüst, Spiegel-Tests; keine Sub-Area mit eigener
Modus-Deklaration (`tools/harness/` bleibt unberührt).

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen, Sub-Area
`*`. Treffer:
[`BEO-ALL/modulliste-spiegel-ungegated`](../observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md)
(3×, verkörpert in den Spiegel-Deckungstests — die DoD verlangt, dass sie das neue
Modul fangen),
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/observation.md)
(7×, verkörpert in Schritt 19 des Workflow-Skeletts — die Negativliste in §3 ist seine
Anwendung) und
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
(3×, verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.8 — das Pläne-Verzeichnis ist eine
Eingabe außerhalb der Scan-Menge, §6). Der Zähler folgt aus den Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-10-30: `upstream-drift` grün (jüngster Lauf 2026-10-30T06:41:15Z), `image-scan`
grün (2026-10-30T10:04:41Z). Nichts zu lesen, nichts ausstehend.

### Sub-Area: `*` (Default)

Alle berührten Sub-Areas GF.
