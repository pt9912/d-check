## Modul 2 — Harness-Bootstrap

<!-- Quelle: [01-spec-und-architektur/modul-02-harness-bootstrap.md](https://github.com/pt9912/ai-harness-course/blob/v6.0.0/kurs/de/01-spec-und-architektur/modul-02-harness-bootstrap.md) -->

### Harness-Einordnung (Modul 2)

Bootstrap ist die **initiale Anwendung des Steering-Loops**, laufend
ausgeführt, bis das Repo *steady state* erreicht. Was
[Modul 11 — Verification Harness](modul-11-verification.md)
als laufende Praxis lehrt (Beobachtung → Guide/Sensor → Diff →
Closure), lehrt dieses Modul als *initiale Aufsetzungs-Praxis*:
gleiche Sensoren und Guides, andere Anwendungsphase. Die abstrakte
Verbindung steht in
[`grundlagen-bootstrap.md` §Verbindung zum Steering-Loop](grundlagen-bootstrap.md#verbindung-zum-steering-loop).

[Modul 1 §Source-Precedence-Block, Schritt 0](modul-01-entwicklungszyklus.md#source-precedence-block)
hat den Bootstrap-Modus als Kurz-Vorgriff eingeführt (Baseline und Modus als
Voraussetzung für den Lebenszyklus); dieses Modul ist die Vollform —
die *Diagnose-Praxis*, die Schritt 0 zur Vorbedingung jeder
modusabhängigen Aktion macht.

Gegen die vier Harness-Linsen:

* **Drift** — Bootstrap ist der erste *Drift-Sensor*: ohne formellen
  Bootstrap bleibt jede spätere Drift unmessbar.
* **Reproduzierbarkeit** — der dokumentierte Modus pro Sub-Area ist
  Voraussetzung dafür, dass ein zweiter Lauf am selben Repo dieselben
  Klassifikationsentscheidungen trifft.
* **Auditierbarkeit** — der Bootstrap-Modus macht die *Begründung*
  jeder Folge-Entscheidung explizit.
* **Steering-Loop** — Bootstrap = initiale Steering-Loop-Anwendung
  (siehe oben).

### Kernidee (Modul 2)

**Bootstrap ist ein fortlaufender Modus, kein einmaliges Event — und
er gilt pro Sub-Area, nicht pro Repo.** Der Modus ist ein beobachtbares
Verhältnis zwischen Code und Doku, kein Etikett auf dem Repo. Wer den
Modus pro Sub-Area diagnostizieren kann, weiß, *welcher* Trigger die
nächste sinnvolle Aktion auslöst — und spart sich das Ausprobieren auf
der Phasen-Ebene. Genau dieses Diagnose-Vermögen ist der Kern der
Bootstrap-Diagnose; die Modus-Wahl als Planungs-Entscheidung
(also: *welcher Modus gilt für jede vom nächsten Slice berührte
Sub-Area — und warum?* — der Slice selbst trägt keinen Modus, er
berührt nur Sub-Areas, die je einen tragen) folgt später in
[Modul 5 — Planning Harness](modul-05-planning-harness.md).

#### Wann wechselt der Modus? Drei Anzeichen

Im laufenden Betrieb verändern sich Sub-Areas zwischen den Modi. Drei
beobachtbare Anzeichen, an denen sich ein Modus-Wechsel ankündigt:

1. **Diskrepanz-Häufung ändert sich** (Indikator in beide
   Richtungen). *BF → GF Graduation:* der Reconciliation-Backlog
   einer BF-Sub-Area schrumpft über mehrere Slices, neue Inventur-
   Schritte melden keine Diskrepanzen mehr. *GF → BF Drift:* in
   einer als GF gemeldeten Sub-Area werden plötzlich Diskrepanzen
   sichtbar (Tests ohne Spec-Anker, ADRs ohne Code-Bezug).
2. **Test-Bestand übertrifft Spec-Anker.** Wenn das Test-Bestand
   strukturell mehr prüft als die Spec behauptet (z. B. Edge-Case-
   Tests ohne `LH-*`-ID), ist die Sub-Area de facto **von GF nach BF
   gedriftet** — der Code "weiß" mehr als die Doku. Symptom: bei
   jeder Code-Änderung muss die Spec nachgezogen werden, statt
   umgekehrt.
3. **Carveout-Auflösung schließt eine BF-Sub-Area.** Wenn ein
   Inventur-Carveout durch einen Reconciliation-Slice geschlossen
   wird (orphan code bekommt nachträglich seinen Anforderungs-Anker
   in der Spec oder als retroaktiver ADR), nähert sich die zugehörige
   Sub-Area der **Graduation zu GF**. Symptom: der
   Reconciliation-Backlog sinkt um genau einen Eintrag, der
   schließende Slice trägt einen "CO-NNN aufgelöst durch
   LH-FA-MMM"-Hinweis im Closure-Block, und die nächste Inventur
   meldet die Sub-Area mit einer Diskrepanz weniger.

   **Umgekehrt: eine Diskrepanz-Häufung *eröffnet* eine BF-Sub-Area.**
   Wenn mehrere Carveouts denselben Geltungsbereich tragen oder sich
   ein systemisches "Code existiert vor Doku"-Muster zeigt, ist die
   richtige Antwort eine BF-Sub-Area-Markierung mit Graduation-Plan,
   nicht eine Carveout-Kaskade — siehe
   [Modul 7 §Werkzeug-Wahl bei Diskrepanz](modul-07-carveouts.md#werkzeug-wahl).
   Die Carveout↔BF-Klammer trägt damit in beide Richtungen: Auflösung
   schließt eine BF-Sub-Area, Häufung eröffnet eine.

Diese drei Anzeichen sind die Sensor-Seite der Bootstrap-Diagnose.

### Regeln gegen typische Fehlannahmen (Modul 2)

* Bootstrap ist ein fortlaufender
  Modus, der sich über Sub-Areas und Phase-Reife entwickelt. Jeder
  Trigger ist ein Bootstrap-Mikro-Event. Wer Bootstrap als Setup
  versteht, übersieht die Modus-Wechsel, die im laufenden Betrieb
  passieren, und produziert daher dauerhaft Findings ohne
  Modus-Bewusstsein.
* Modus gilt **pro Sub-Area**. Ein Repo kann in den
  *Konventionen* BF und in der *Spec-Schreibung* GF sein.
  Sub-Area-Heterogenität ist der Normalfall, nicht die Ausnahme — ein
  einzelnes Modus-Etikett auf dem Repo ist deshalb immer zu grob.
* auch im GF-Modus entstehen Trigger
  (Diskrepanz, Promotion-Auslöser etc.), nur nicht aus
  *Bestandsinventur*, sondern aus *Konsistenzprüfung des neu
  Geschaffenen*. Die vier Klassen gelten in jedem Modus; was sich
  ändert, ist die Frequenz und die typische Auslöse-Quelle.
* BF ist der *typische* Ausgangspunkt
  realer Repos.
  BF kann systematisch in Richtung GF graduieren — *Graduation* ist
  eine ausgewiesene Bedingung mit Konvergenz-Auftrag, kein
  Wunschdenken. Die Frage ist nicht *ob* graduiert wird, sondern
  *wie weit* das Repo schon ist und *welche Sub-Area* als nächstes
  graduations-reif wird.
* Eine Struktur qualifiziert erst über die drei Inklusions-Achsen
  (Schwelle ≥ 2, siehe
  [`grundlagen-bootstrap.md` §Was ist eine Sub-Area?](grundlagen-bootstrap.md#was-ist-eine-sub-area)).
  Der übersprungene Qualifikations-Schritt erzeugt **beide**
  Granularitäts-Fehler zugleich: *zu grob* — ein Aggregat wie *"Backend"*
  wird als *eine* Sub-Area gelabelt, statt in mehrere aufgeteilt; *zu
  fein* — ein substanzloses Verzeichnis (*"Struktur ohne Substanz"*, nur
  eine Achse erfüllt) wird zur Sub-Area erhoben, obwohl es eine
  Sub-Area-*Aspirantin* bleibt. Gemeinsame Wurzel beider Fehler:
  Reife/Substanz wird aus einem Oberflächenmerkmal (Existenz, Größe)
  *abgelesen* statt über die Achsen *geprüft*.

### Greenfield-Bootstrap: Schritt-Sequenz (Modul 2)

Neun Schritte (0–8) in drei Phasen: *Orient* (0–1, Modus + Baseline),
*Action* (2–3, vendoren + Konventionen), *Content* (4–8, normativer
Inhalt). Endpunkt: bereit für die ersten Code-Slices.

#### Detail-Tabelle (Schritte 0–4: Setup-Phase)

Trigger-Anker (T1, T2, T4, T5, T6, T7) sind Instanz-Beispiele der
vier Trigger-Klassen aus
[`grundlagen-bootstrap.md` §Vier Trigger-Klassen](grundlagen-bootstrap.md#vier-trigger-klassen) —
die abstrakten Definitionen stehen dort, hier nur die Instanzen.

*Hinweis zur T-Nummerierung:* Die Trigger sind durch das
Beispiel nummeriert (T1 = Pointer in README,
T2 = Pointer in AGENTS, T4 = Promotion, T5 = erste ADR-Vorschläge,
T6 = Cross-Reference ADR → Spec (aufwärts), T7 = ADR-Review-Auslöser).
**T3 ist BF-spezifisch** und tritt im Greenfield-Walkthrough nicht
auf — er erscheint im Brownfield-Walkthrough als Sync-Trigger in
BF-Diskrepanz-Auslöse-Variante.

| # | Aktion | Berührte Dateien (Phasen-Übergang) | Trigger |
|---|---|---|---|
| 0 | Modus pro Sub-Area entscheiden: GF für *Konventionen*, *Spec*, *Architektur*, *ADR* (alle vier Doku-führt). | keine | keine — Vorbedingung |
| 1 | Baseline-Auswahl (Kurs-Harness) + Repo-Klasse (Tooling) + ID-Schema festlegen: Vertrags-Präfix wählen (`LH-*`), `SPEC-*`/`ARC-*`/`MR-*` aus der Baseline übernehmen | keine | reift 2/3 |
| 2 | **Baseline vendoren** — Regelwerk *und* Templates nach `.harness/baseline/<tag>/{regelwerk,templates}/` (+ `SHA256SUMS`, netzlos) als präsente, gepinnte Referenz; **Tooling** (`Makefile` mit d-check-Doku-Gate, `.d-check.yml`) als Startgerüst übernehmen; **Dokument-Skelette** aus der vendored Baseline (`…/templates/`) kopieren *und ausfüllen* | Dokument-Skelette **0 → 1**; vendored Baseline + Tooling tragen keine Phase-Reife | keine |
| 3 | `harness/conventions.md` mit MR-000 (Baseline-Aussage inkl. ID-Schema) + leerer Index-Tabelle — dieses Repo weicht nirgends ab, also entsteht **keine** `MR-<NNN>`-Datei | `conventions.md` 0 → 1 | **T1** (Pointer auf `conventions.md` in `harness/README.md`), **T2** (Pointer in `AGENTS.md`) |
| 4 | `spec/lastenheft.md` Outline mit `LH-FA-*`/`LH-QA-*` | `lastenheft.md` 1 → 2 | keine direkt |

#### Detail-Tabelle (Schritte 5–8: Inhalts-Phase)

| # | Aktion | Berührte Dateien (Phasen-Übergang) | Trigger |
|---|---|---|---|
| 5 | `docs/plan/planning/roadmap.md` mit Welle + Release-Mapping; `releasing.md` mit Release-Strategie | `roadmap.md` 1 → 2; `releasing.md` 1 → 2 | keine |
| 6 | Sensors-Roster im "Nicht behauptet"-Block (Prosa-Pointer-Liste, kein Status) | `harness/README.md` §Sensors Sub-1 → Sub-2; `AGENTS.md` §4 Sub-1 → Sub-2 | **T4** (Promotion-Auslöser bei erstem Code-Slice) |
| 7 | `spec/architecture.md` + `spec/spezifikation.md` Outline mit `ARC-*`/`SPEC-*` | beide 1 → 2 | **T5** (erste ADR-Vorschläge aus Architektur-Outline) |
| 8 | `docs/plan/adr/0001-doc-source-of-truth.md` mit Status *Proposed* | `0001-…md` 0 → 2; ADR-Index 1 → 2 | **T6** (Cross-Reference: ADR → Spec, aufwärts), **T7** (ADR-Review-Auslöser) |
| Bootstrap-Ende | Bereit für ersten Code-Slice — Workflow-Übergang | — | — |

**Anmerkung zur Tabellen-Splittung.** Die Setup-Phase (0–4) etabliert
*was wo lebt*; die Inhalts-Phase (5–8) füllt mit *normativem Inhalt*.
Trigger entstehen ab Schritt 3 (Konventionen-Adoption); die Setup-Phase
legt die *Quellen* an, die Inhalts-Phase erzeugt die *Folge-Bezüge*
zwischen Dokumenten und ADRs. Die Tabellen-Trennung macht das
kognitiv lesbar — die Phasen verschwimmen sonst.

**Anmerkung zum Instanziierungs-Zeitpunkt (Schritt 2).** Beim Bootstrap
entstehen **nur die Gründungs-Dokumente** — je genau eines pro Repo, gefüllt und
behalten: `spec/lastenheft.md`, `spec/spezifikation.md`, `spec/architecture.md`,
`harness/conventions.md` (Index; einzelne `MR-<NNN>` entstehen pro
Adaption, nicht beim Bootstrap), `harness/README.md`, `AGENTS.md`,
`docs/plan/planning/roadmap.md` und der Gründungs-ADR `0001` (Skelette in
Schritt 2, gefüllt in 3–8). Die **wiederkehrenden Artefakte** — `slice`,
`welle`, weitere ADRs (`NNNN-*`), `carveout`, `review-report` — werden **nicht**
beim Bootstrap vorab kopiert, sondern **pro Instanz** aus der vendored Baseline
(`.harness/baseline/<tag>/templates/…`), wenn der Workflow sie erreicht
(Modul 4–10). Nur der Gründungs-ADR `0001` entsteht schon hier; jeder weitere
ADR ist ein wiederkehrendes Vorkommen. Die vendored Baseline ist deren
**einzige Referenz-Form** — **keine Blank-Kopie im Repo vorhalten.**

**Anmerkung zur vendored Baseline (Schritt 2).** Regelwerk *und* Templates
werden beim Bootstrap **committet vendored**
(`.harness/baseline/<tag>/{regelwerk,templates}/`
+ `SHA256SUMS`, netzlos materialisiert), nicht pro Lauf extern gefetcht — es
ist die *präsente, nachschlagbare Vertiefung* zur verkörperten Form: pro
Entscheidung, deren operative Detailtiefe Briefing und Konventionen nicht
tragen (Trigger-Klassen, Sub-Area-Qualifikation, Carveout-vs-Reconciliation,
Modus-Diagnose), wird der *relevante Abschnitt* nachgeschlagen (README =
Index), ohne das ganze Regelwerk im Kontext zu halten. Das ist das
Modul-0-Prinzip — *Per-Lauf-Relevantes gehört verkörpert, nicht extern
nachgeladen* — angewandt auf das Regelwerk selbst; ein realer Konsument-Repo
(d-check) führt es netzlos vendored und schlägt im Slice-Betrieb routinemäßig
nach. Die Templates tragen dabei zwei Rollen: **vendored** als Referenz-Form
(worauf das Regelwerk als „Ziel-Form" verweist — `../templates/…` löst netzlos
lokal auf) und **kopiert-und-ausgefüllt** als deine eigenen Artefakte.

#### Freshness-Audit der vendored Baseline (Schritt 2)

Eine vendored Kopie
driftet still von der Quelle weg, sobald ein neues Kurs-Release erscheint;
Pinnen ohne Überwachung ist die halbe Maßnahme (Doktrin „pinnen und
überwachen", Modul 12/14). Der Freshness-Audit hat sieben Eigenschaften:

* **Beobachtbarer Auslöser, keine Kalenderpflicht** — die Frage „ist mein
  `<tag>` noch das aktuelle Kurs-Release?" an ein Ereignis (z. B.
  Konventions-Änderung mit Baseline-Audit) oder periodisch gebunden.
* **Netz-Operation, außerhalb der Gates** — fragt den Upstream, bleibt daher
  aus den netzlosen Gates heraus (die Integritätsprüfung der Arbeitskopie
  prüft nie den Upstream); Wartung, kein Feedback-Gate.
* **Release-*Liste* prüfen, nicht das Asset** — der Hash des gepinnten Assets
  fängt nur ein nachträglich verändertes Release, **nicht einen neuen Tag**;
  ein Sensor auf den gepinnten Tag meldet „kein Drift", während upstream ein
  neuer Tag steht. Auf einen neueren Tag in der Release-Liste prüfen. d-check
  `sources` automatisiert die Asset-Prüfung (`source-pin` auf den `sha256` einer
  http(s)-Quelle, `source-drift` bei Inhaltsabweichung); deckt die
  *Integritäts*-Hälfte ab, ersetzt die Release-Listen-Prüfung nicht.
* **Der Review geht durch die Adaptions-Liste, nicht nur durch den Diff** —
  Frage pro Eintrag in
  [`harness/conventions.md`](grundlagen-harness-dateien.md#harnessconventionsmd-als-konventionsspeicher):
  *Regelt die neue Fassung das, wofür diese Adaption angelegt wurde?* Davor steht
  ein Formcheck — **den erledigt das Doku-Gate**: Die vendored Baseline liegt im
  Repo, ihre Dateien sind gültige Link-Ziele, und ein `Geltungsbereich`, der nach
  dem Update ins Leere zeigt, meldet sich als `anchor-missing`. Ein toter Anker
  ist kein Ausgang, sondern ein Formfehler und wird zuerst repariert. Einmal
  prüfen, dass `.harness/baseline/` im Prüfumfang liegt; danach automatisch. Dann **fünf Ausgänge** — sie beziehen sich auf
  das *Delta* der neuen Fassung, nicht auf den Zustand der Baseline: Die
  Adaption wird **gegenstandslos** (die neue Fassung regelt das **jetzt neu**
  selbst → Rückbau), sie **bleibt gültig** (die neue Fassung
  ändert daran nichts → stehen lassen; Normalfall, auch wenn die Baseline das
  Thema regelt — eine Adaption weicht ja gerade von ihr ab), sie ist **teilweise
  überholt** (die neue Fassung regelt **jetzt neu** einen Teil der
  abgewichenen Regel → ablösen,
  engere Nachfolgerin; Abgrenzung gegen die *Schärfung* in
  [§Konventionsspeicher](grundlagen-harness-dateien.md#harnessconventionsmd-als-konventionsspeicher)),
  ihr **Bezug ist entfallen** (die Baseline regelt das Thema gar nicht mehr —
  dann ist der Eintrag keine Adaption mehr; ein Nachfolge-Eintrag löst ihn auf
  und hält fest, dass die Baseline dazu seit `<tag>` schweigt — wo die Setzung
  selbst künftig lebt, entscheidet das Repo), oder sie **widerspricht** der
  neuen Fassung — die regelt das
  Thema **neu und anders**, als die Adaption es setzt. Das ist der einzige
  Ausgang, an dem das Delta die Antwort **nicht** vorgibt, denn hier hat das
  Repo eine Wahl: Entweder die Adaption gilt in ihrem Geltungsbereich weiter —
  dann gehört der Widerspruch benannt, sonst adoptiert das Repo eine Regel, die
  es nicht befolgt —, oder das Repo **übernimmt** die neue Regel. Übernehmen ist
  ein Rückbau wie oben, nur aus einem anderen Grund: Bei *gegenstandslos* hat
  die Baseline dem Repo recht gegeben, hier gibt das Repo der Baseline recht.
  Dieser Unterschied gehört in die `Begründung` des Nachfolge-Eintrags, weil er
  eine **Entscheidung** ist und kein Befund. Und will das Repo übernehmen, kann
  es aber noch nicht, ist das keine `MR` mehr: Eine Adaption sagt *„diese Regel
  gilt hier nicht"*, ein Carveout sagt *„sie gilt, wir erfüllen sie noch
  nicht"* — befristet, mit Auflösungs-Trigger und Folge-Slice
  ([Modul 7](modul-07-carveouts.md)). **Auch `permanent`-Einträge werden mitgeprüft:**
  *permanent* heißt „kein automatischer Auflösungs-Trigger", nicht
  „unauflösbar". War die Adaption eine *Lockerung* und die neue Baseline
  verschärft, ist die richtige Antwort ein **Carveout** mit Auflösungs-Trigger
  ([Modul 7](modul-07-carveouts.md)), keine stille Dauer-`MR`. Ohne diesen
  Durchgang bleibt eine präzise formulierte Auflösungs-Bedingung jahrelang
  stehen, obwohl sie längst erfüllt ist: Ein Trigger, den niemand abfragt, ist
  kein Wächter.
* **Der Review vergleicht auch die Form, nicht nur die Regeln** — ein neuer
  Stand kann die *Struktur* der Artefakte ändern, nicht nur die Regeln über sie;
  dafür gibt es kein Trigger-Feld, das sich melden könnte. Nach dem
  Re-Vendoring steht die neue Referenz-Form unter
  `.harness/baseline/<tag>/templates/` und ist die Vergleichsgrundlage. Weil der
  Vendoring-Pfad `<tag>`-gescopt ist, liegen alte und neue Form nebeneinander:
  `diff -r .harness/baseline/<alt>/templates .harness/baseline/<neu>/templates`
  zeigt umbenannte Sektionen und neue Felder direkt. Das alte Verzeichnis fällt
  erst, wenn der Review durch ist.
  Ob ein Feld **Pflicht** ist, entscheidet nicht die Feldzahl im Template, sondern
  die Pflichtgliederung im vendored Regelwerk — für `harness/conventions.md` in
  [§Konventionsspeicher](grundlagen-harness-dateien.md#harnessconventionsmd-als-konventionsspeicher),
  für `harness/README.md` in
  [§Einstiegspunkt](grundlagen-harness-dateien.md#harnessreadmemd-als-einstiegspunkt).
  Für `AGENTS.md` und das Lastenheft gibt es **keine** Pflichtgliederung — dort
  entscheidet die Referenz-Form selbst, was ein neues Feld ist. Für
  **Singletons** (`harness/conventions.md`, `AGENTS.md`, `harness/README.md`,
  Lastenheft, …): neue *optionale* Felder verlangen keine Nacharbeit am gefüllten
  Artefakt, neue **Pflicht**-Felder und umbenannte Sektionen schon — sonst
  behauptet das Repo
  eine Baseline-Konformität, die seine Artefakte nicht tragen. Für
  **wiederkehrende** Templates (ADR, Slice, Welle, Carveout, Review-Report)
  gilt die Append-only-Logik: Neue Instanzen folgen der neuen Form, bestehende
  werden nicht rückwirkend umgeschrieben.
* **Rückbau ist ein neuer Eintrag, kein Edit** — eine aufgelöste `MR-<NNN>`
  wird nicht überschrieben, sondern bekommt einen Nachfolger, der sie auflöst
  und den Baseline-Stand nennt, der den Trigger gefeuert hat. Die alte Zeile
  ist die historisch korrekte Aussage über den damaligen Zustand
  (Append-only-Disziplin wie bei ADRs, [Modul 4](modul-04-adrs.md)).
* **Eine Stichprobe gegen den Bestand, nicht gegen das Delta** — sie läuft
  unabhängig vom Ergebnis der Tag-Frage, auch wenn der Pin aktuell ist. Die
  sechs Eigenschaften darüber haben dann nichts zu tun, denn ihr Gegenstand ist
  die *Änderung*; der Gegenstand dieser hier ist der *Bestand*, und den gibt es
  immer. Eine Baseline-Regel,
  die **nie** ins ausgefüllte Artefakt übernommen wurde und sich seither
  **nie** geändert hat, erzeugt keinen Template-Diff und hat keinen
  `MR`-Eintrag, den der Adaptions-Durchgang abschreiten könnte — sie ist
  unsichtbar, *weil* sie alt und stabil ist. Genau über sie behauptet die
  `MR-000`-Aussage „keine inhaltlichen Adaptionen ggü. Baseline-Default" etwas,
  und nichts prüft es.
  **Auswahlkriterium:** aus den Baseline-Abschnitten ziehen, die seit dem
  adoptierten `<tag>` **kein Delta** hatten — die Komplementärmenge zu
  `git diff <alt> <neu> -- .harness/baseline/`; kein Delta-Review hat sie je
  berührt. **Umfang: ein Abschnitt pro Audit**, rotierend — keine
  Vollinventur, sonst verliert die Welle ihr Closure-Kriterium. Die Frage pro
  Regel: *Steht sie im ausgefüllten Artefakt — oder als deklarierte
  Abweichung?* Zweimal nein heißt: nie übernommen.
  **Ausgang:** ein einzelner Fund geht den Weg jeder Diskrepanz
  ([Modul 7 §Werkzeug-Wahl](modul-07-carveouts.md#werkzeug-wahl)) — Übernahme
  im nächsten Slice oder Carveout mit Auflösungs-Trigger. **Mehrere** Funde
  treffen nicht die einzelne Regel, sondern die `MR-000`-Aussage: Sie ist dann
  falsch und wird korrigiert — durch Übernahme oder durch `MR`-Einträge, die
  die Abweichungen deklarieren. Eine BF-Markierung ist hier **nicht** die
  Antwort; sie regelt Doc ↔ Code und trifft die Achse nicht
  ([`grundlagen-bootstrap.md` §Modus pro Sub-Area](grundlagen-bootstrap.md#modus-pro-sub-area-greenfield-vs-brownfield)).

Ein neuer Tag löst einen **Review** aus (Re-Vendoring mit eigenem Diff),
keinen stillen Auto-Bump.

#### Gate-Fragment `d-check.mk` (Schritt 2)

Das d-check-Doku-Gate wird nicht
handgeschrieben: `d-check.mk` beim Bootstrap aus der gepinnten d-check erzeugen
(`d-check --print-mk > d-check.mk`), per `-include` einbinden; den
`DCHECK_DIGEST`-Platzhalter der `Makefile` mit dem eigenen Digest füllen. Nicht
mitgeliefert (tool-/versionsspezifisch); bei jedem d-check-Bump neu erzeugen.
Das Tool pflegt die Recipe-Form (`--network none`, Target-Set); die
advisory-Targets sind verfügbar, nicht als Gate behauptet (Modul 13).

**Und das Fragment mountet.** Der Prüflauf ist hermetisch — Quellen per `COPY`
ins Image, Ergebnisse über `stdout`
([Modul 14 §Der Prüflauf ist hermetisch](modul-14-docker-harness.md#der-prüflauf-ist-hermetisch--kein-mount)).
Ein generiertes Fragment kennt diese Form nicht; es reicht den Baum read-only
herein. Das ist **zulässig, solange das Werkzeug nur liest** — der Schaden der
Mount-Form entsteht am Schreiben. Wer auch diesen Lauf hermetisch will, bindet
das Fragment **nicht** ein und schreibt die Recipe aus. Dann pflegt nicht mehr
das Tool die Recipe-Form, sondern du, und die Abweichung gehört als `MR-<NNN>`
deklariert.

#### T1/T2 als Sync-Trigger konkret (Schritt 3)

Sobald
`harness/conventions.md` existiert, müssen die auf sie verweisenden
Dokumente einen Eintrag bekommen: ein Pointer in `harness/README.md`
(T1) und in der Source-Precedence-Liste von `AGENTS.md` (T2). Wer den
Eintrag vergisst, hat einen **Sync-Drift**, den der Doku-Konsistenz-Agent
([Modul 15 §Observability](modul-15-observability.md)) später als
Inkonsistenz findet.

### Brownfield-Bootstrap: Schritt-Sequenz (Modul 2)

Neun Schritte (1–9); gegenüber GF neu sind die Phasen *Discover*
(Code-Inventur) und *Diskrepanz-Schock* (Reconciliation) sowie das
Konvergenz-Ziel *Graduation BF → GF* pro Sub-Area.

#### Detail-Tabelle (Schritte 1–4: Inventur-Phase, was bei BF anders ist)

*Hinweis zur Schritt-Nummerierung:* BF nimmt die Modus-Setzung in
Schritt 1 mit (wo GF Schritt 0 trägt), weil in BF die
Repo-Klassen-Wahl und die Modus-Antizipation zusammenfallen — beide
hängen von der ersten Inventur-Beobachtung ab und können nicht *vor*
dem ersten Hinschauen entschieden werden. Daher die asymmetrische
Nummerierung GF 0–8 vs. BF 1–9.

| # | Aktion | BF-Besonderheit gegenüber GF |
|---|---|---|
| 1 | GF-Schritte 0 und 1 in einem Schritt zusammengefasst: Modus-Antizipation "BF pro Sub-Area" + Baseline-Auswahl + Repo-Klasse + ID-Schema festlegen (Vertrags-Präfix wählen, `SPEC-*`/`ARC-*`/`MR-*` aus der Baseline) | + explizite Modus-Setzung mit Sub-Area-Aufzählung; Repo-Klassen-Wahl und Modus-Antizipation fallen zusammen |
| 2 | **Code-Inventur (Discovery):** Makefile, CI, Tests, README, Commit-Messages inventarisieren als Lerner-Schritt | **neu in BF** — kein Repo-Artefakt entsteht, nur Lerner-Wissen |
| 3 | **Baseline vendoren** — Regelwerk + Templates (`.harness/baseline/<tag>/{regelwerk,templates}/` + `SHA256SUMS`) + Skelette kopieren-und-ausfüllen | wie GF |
| 4 | `harness/conventions.md` mit Modus = BF pro Sub-Area, MR-000-Aussage | Modus-Block anders strukturiert (BF-Deklarationen + Konvergenz-Auftrag pro Sub-Area) |

#### Detail-Tabelle (Schritte 5–9: Reconciliation-Phase)

| # | Aktion | BF-Besonderheit gegenüber GF |
|---|---|---|
| 5 | Sensors-Haupt-Tabelle direkt aus Makefile-Kommentaren entstehen lassen | **gegenteilig zu GF** — Targets existieren, keine "Nicht behauptet"-Promotion nötig; **T3** (Sync-Trigger in BF-Diskrepanz-Auslöse-Variante: Sensor-Lücke = impliziter Pointer-Mismatch zwischen Makefile-Realität und Sensor-Tabelle) wird sichtbar |
| 6 | Lastenheft aus Code/Tests/CI rückbauen | Inventur-Umkehr; Diskrepanz-Material entsteht (Code ohne Anforderung, Test ohne LH-Bezug) |
| 7 | Architektur + Spezifikation aus `src/` rückbauen; **retroaktive ADRs** für implizite Entscheidungen | ADRs teils retroaktiv mit Status *Accepted* (oder *Superseded*, falls Entscheidung schon revidiert) |
| 8 | **Diskrepanz-Schock:** `docs/plan/planning/reconciliation.md` anlegen (Ziel-Form [`reconciliation.template.md`](../templates/docs/plan/planning/reconciliation.template.md)) und jeden Fund als Zeile klassifizieren: (a) Code ohne Anforderung → Carveout, (b) Anforderung ohne Code → Reconciliation-Slice, (c) Entscheidung ohne Beleg → retroaktive ADR | **BF-spezifischer Schritt** — die Reconciliation-Pflicht macht hier den Modus-Übergang sichtbar; **T3** (Sync-Trigger zwischen Code-Realität und Anforderungs-Anker, in BF-typischer Diskrepanz-Auslöse-Variante) als typische Trigger-Quelle |
| 9 | Roadmap als Reconciliation-Plan; letzte Welle = Graduation pro Sub-Area | Inhalt anders als GF-Roadmap — Plan ist Diskrepanz-Auflösungs-Sequenz, nicht Feature-Sequenz |
| Bootstrap-Ende | Reconciliation-Backlog steht — jede Zeile in `reconciliation.md` trägt ihre Auflösung, je Sub-Area ist die Graduation-Welle benannt | — |

#### Das Reconciliation-Register

Der Rückbau (Schritte 6–7) erzeugt Funde, die niemand sonst festhält: Weder die
Kennung eines Carveouts noch seine sechs Kopffelder sagen, dass er aus der
Inventur stammt, und einem Slice in `open/` sieht man nicht an, ob er eine
Zusage einlöst oder eine Lücke schließt. Ohne eigenen Ort ist die Menge der
offenen Funde deshalb **nicht bestimmbar** — und ein Closure-Kriterium über eine
unbestimmbare Menge ist keines. Der Ort ist
`docs/plan/planning/reconciliation.md`, flach neben dem Beobachtungs-Register,
das dieselbe Form hat: eine stehende Datei mit einer Zeile je Fund
(Ziel-Form: [`templates/docs/plan/planning/reconciliation.template.md`](../templates/docs/plan/planning/reconciliation.template.md)).

**Kennungen sind `RC-<NNN>`**, fortlaufend je Datei, ohne Bereichssegment — die
Sub-Area steht als eigene Spalte, und alle Zeilen liegen in *einer* Datei
([`grundlagen-source-precedence.md` §Vergabe](grundlagen-source-precedence.md#vergabe-woher-die-nächste-nummer-kommt)).

**Drei Klassen, aus Schritt 8 — sie bestimmen, was in der Spalte *Auflösung* steht:**

| Klasse | Was gefunden wurde | Auflösung |
|---|---|---|
| Code ohne Anforderung | Code, den keine Anforderung deckt | `CO-<NNN>` samt dessen Folge-Slice |
| Anforderung ohne Code | zugesagt, aber nicht geliefert | der Reconciliation-Slice |
| Entscheidung ohne Beleg | im Code getroffen, nirgends begründet | die retroaktive `ADR-<NNNN>` |

Die dritte Klasse ist mit dem Eintrag erledigt: Eine retroaktive ADR entsteht in
Schritt 7 mit Status *Accepted* (oder *Superseded*) und wartet auf nichts; ihre
Zeile geht sofort nach *Aufgelöste Einträge*.

**Fortgeschrieben wird beim Auflösen, nicht beim Bootstrap-Ende.** Schließt ein
Slice einen Fund, wandert dessen Zeile mit Datum und auflösendem Artefakt in die
zweite Tabelle — dieselbe Append-only-Disziplin wie beim Beobachtungs-Register:
Wer eine Zeile still löscht, macht sie ununterscheidbar von einer, die es nie
gab. Der Slice-Plan trägt das als Closure-Pflicht.

**Der Backlog *steht***, wenn jede Zeile ihre Auflösung trägt — nicht, wenn das
Register leer ist; beim Bootstrap-Ende ist es im Gegenteil voll. **Leer wird es
je Sub-Area bei der Graduation**, und genau das ist dort die
Graduation-Bedingung: keine offene Zeile mehr für diese Sub-Area. Wo der Modus
einer Sub-Area steht und wie ihre Graduation-Bedingung lautet, bleibt
`harness/conventions.md` §Modus-Deklaration — das Register liefert nur den
Zählstand, es führt keine zweite Modus-Tabelle.

### Phasen × Modus — die zweidimensionale Reife-Matrix

Beide Walkthroughs bewegen Artefakte durch **Phase-Reife** (0–5)
pro Sektion (sechs Stufen, siehe
[`grundlagen-bootstrap.md` §Sektionsweise Reife](grundlagen-bootstrap.md#sektionsweise-reife-phasen-pro-dokument)).
Die folgende Matrix macht sichtbar, *was Phase-N in GF bedeutet
versus was sie in BF bedeutet* — dieselbe Phase-Stufe,
unterschiedliche Bewegungsrichtung:

| | Greenfield (Doc → Code) | Brownfield (Code → Doc) |
|---|---|---|
| **Phase 0 leer** | Datei existiert nicht — Pflicht zur Anlage aus Konvention | Datei existiert nicht — Inventur stellt fest, dass Doku-Anker fehlt |
| **Phase 1 Skelett** | Template kopiert, *Versprechen* zu füllen | Template kopiert, *Inventur-Auftrag* an Code |
| **Phase 2 Outline** | Top-Level-Wunschbild | Top-Level-Bestandsaufnahme |
| **Phase 3 partiell** | Sektionen versprochen, Code folgt | Sektionen dokumentiert, andere unentdeckt |
| **Phase 4 kohärent** | Vertrag steht, Code wird daran gemessen | Inventur abgeglichen, **Diskrepanz-Schock sichtbar** |
| **Phase 5 stabil** | Change-Process aktiv | Reconciliation-Slice oder Carveout aktiv |

Die Matrix ist *die* Lesart für jede Bootstrap-Aktion: erst
identifizieren, *welche Sub-Area* sich bewegt; dann *welche Phase*
sie erreicht; dann *in welchem Modus* — daraus folgt der nächste
Schritt fast deterministisch.

### Ziel-Form: Projekt-README

Die Root-`README.md` des Adopter-Repos folgt der Vorlage
[`templates/project-readme.template.md`](../templates/project-readme.template.md):
Überblick, Ist-Stand, Vertrauens-Signale. Sie ist **Rang 6** der Source
Precedence ([`grundlagen-source-precedence.md` §Source Precedence](grundlagen-source-precedence.md#source-precedence))
— beschreibend, keine kanonische Quelle. Entsteht in Schritt 2 der
Bootstrap-Sequenz aus der vendored Baseline.
