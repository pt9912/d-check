# Review-Report: slice-215 — 2026-09-08

**Review-Art:** Code — geprüft wird der Diff gegen Slice-Plan, die
`MR-011`-Kette (deklarierte fail-open-Linie) und die Hard Rules
(`AGENTS.md` §3.1, §3.6, §3.7, §5, §6). Erster Slice dieser Folge mit einer
**Verhaltens-Änderung an einem Skript**; die Prüfung ist entsprechend am
Verhalten gefahren, nicht an der Beschreibung.

**Gegenstand:** slice-215 · Commit-Range `HEAD~3..HEAD`
(`2f5e4fff`, `fea736e0`, `73f1098a`) · geändert: `tools/harness/fetch-baseline-cache.sh`,
`harness/sensors/baseline-freshness.md`, der Registereintrag
`BEO-HARN/check-latest-blind-before-pin/state.md`, der Slice-Plan,
`docs/plan/planning/in-progress/roadmap.md`

**Skill:** `.harness/skills/reviewer.md` @ 1.16.0 ·
**Modell:** claude-opus-5[1m] · **Datum:** 2026-09-08

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan slice-215 (Ziel/Abgrenzung §1, DoD §2, Plan §3, Trigger §4,
  Closure-Trigger §5, Risiken §6, Vorprüfungen/Modus §8)
- `harness/sensors/baseline-freshness.md` (Vertrag · Grenze · Ausgänge · Bindung)
- `MR-011`-Kette (Baseline-Pin, fail-open), `MR-004`, `MR-013`, `MR-025`,
  `MR-053`, `MR-054`
- `AGENTS.md` §3.1 (Docker/make-only), §3.6 (Gate-Lockerung braucht ADR),
  §3.7 (Kommentar-/Zustandsfeld-Klassen), §5 (gemessene Menge · Geltungsbereich
  eines Zitats), §6, §4 (Target-Tabelle)
- Baseline `v6.5.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice,
  §Zwei Schritte vor der Modus-Begründung; `regelwerk/modul-06-roadmap.md`
  §Das Beobachtungs-Register; `regelwerk/modul-13-quality-gates.md` §Hard Rule
- `.github/workflows/upstream-drift.yml` (der einzige Konsument)
- Vorherige Findings am gleichen Gegenstand: slice-214-Review R1 (dort die
  Klasse `inventur-behauptet-vollstaendigkeit-die-die-stichprobe-bricht`);
  Registereinträge `BEO-ALL/eigene-menge-gemessen-fremde-behauptet` (15×),
  `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`,
  `BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`

**Eigene Läufe** (echte Ausgaben, keine behaupteten):

- `make gates` → **Exit 0**; `[gates] baseline-verify + workflow-pins +
  doc-check + lint + test + arch-check + coverage-gate + semgrep +
  gate-consistency + planning-check green`; d-check-Läufe je
  `737 Datei(en) geprüft, 0 Befund(e)`; `Ran 55 rules on 63 files: 0 findings`;
  `fetch-baseline-cache: verify ok (54 Dateien, vollständig)`.
- `make doc-check` → **Exit 0**, `737 Datei(en) geprüft, 0 Befund(e)`.
- `bash tools/harness/fetch-baseline-cache.sh --selftest` → **Exit 0**,
  `selftest ok (9 Proben)`, alle neun `OK`.
- Regressionslauf `--check-latest` (Netz) → **Exit 0**,
  `check-latest OK (Currency) — Pin v6.5.0 ist der neueste Release-Tag.` +
  `check-latest OK (Content) — gepinnter Tag v6.5.0 upstream unverändert`.
- Bruch-Test **nachher** `--check-latest v0.0.1` → **Exit 3**,
  `Pin v0.0.1 NICHT in der Release-Liste; Currency unbestimmt …`.
- Bruch-Test **vorher** (Skript aus `HEAD~1`, identischer Aufruf) → **Exit 0**
  mit derselben Meldung in der alten Fassung. Beide Richtungen reproduziert.
- Gegenprobe **Netz-Ausfall**, isolierte Kopie: siehe F-1 — teilweise
  übertragene Antwort landet **nicht** im `skip`-Zweig.
- Messung Release-/Tag-Liste (GitHub-API, `curl`+`grep`): 58 Release-`tag_name`,
  58 git-Tags, `comm` in **beide** Richtungen leer; 51 davon in strikter
  `vX.Y.Z`-Form, 7 `templates-v*`; alle 58 `"draft": false`,
  `"prerelease": false`.

---

## Findings

### F-1 · MEDIUM · Die fail-open-Zusage ist aus **einer** Ausfall-Form verallgemeinert; eine gemessene zweite Form macht den Nachtlauf jetzt falsch rot

- **kategorie:** MEDIUM — Anker 8 (*Botschaft/Notiz verallgemeinert über die
  Messung hinaus*), verstärkt durch Anker 18 (die Kehrseite steht im Code, nicht
  im Vertrags-Text). Keine Eskalation auf HIGH: `make baseline-freshness` ist in
  `harness/README.md` ausdrücklich **kein Gate**, und der Fehler zeigt nach
  *rot*, nicht nach *still grün*.
- **quelle:** `AGENTS.md` §5 (*„ihr Schluss reicht **nicht weiter als die
  gemessene Menge**"*) · `MR-011`-Kette (deklarierte fail-open-Linie) ·
  slice-215 §6, Risiko 2
- **pfad:** `tools/harness/fetch-baseline-cache.sh:30-34` (Kopf-Kommentar) ·
  `harness/sensors/baseline-freshness.md` §Grenze 5 · Commit `73f1098a`,
  Absatz *„Die entscheidende Gegenprobe ist gemessen und nicht erschlossen"*
- **befund:** Drei Artefakte tragen denselben Satz. Der Kopf-Kommentar sagt
  *„Ein Netz-/API-Ausfall bleibt SKIP mit exit 0"*, die Sensor-Datei
  *„Ein Netz- oder API-Ausfall landet im `skip`-Zweig (Exit 0), nicht hier"*,
  die Commit-Botschaft *„Die fail-open-Linie ist damit unberuehrt"*. Gemessen
  wurde **eine** Form: ein unerreichbarer API-Host. Die N+1-te Form ist ein
  **abgebrochener bzw. abgeschnittener Transfer** — `curl` schreibt fortlaufend
  in die Pipe, die Antwort ist heute 322 579 Bytes groß, und GitHub liefert
  **neueste zuerst**. Bricht der Transfer nach dem Kopf der Liste ab, ist
  `$tags` nicht leer, enthält aber den (älteren) Pin nicht → der `else`-Zweig
  `currency="ahead"` greift → seit diesem Slice `rc=3`.
  **Gemessen**, isolierte Kopie des Skripts, Byte-identischer Verzweigungs-Code,
  nur `api=` auf eine `file://`-Quelle mit den ersten 12 000 Bytes der echten
  Antwort umgebogen (Tags darin: `v6.5.0 v6.4.0`), Pin `v6.3.1`:

  ```
  fetch-baseline-cache: check-latest (Currency) — Pin v6.3.1 NICHT in der
  Release-Liste; Currency unbestimmt (…). Manuell pruefen.
  EXIT=3
  ```

  **Versagensszenario:** Der Pin ist zwischen zwei Bumps regelmäßig einige
  Releases alt (der Normalzustand). Eine gestörte Verbindung im Nachtlauf
  erzeugt dann ein Rot, das nach *„Pin unauffindbar"* aussieht und ein
  Netz-Ausfall ist — genau der Fall, den `skip` deklariert abfängt. Vor diesem
  Slice war er folgenlos (Exit 0), jetzt färbt er den Job. Der Slice-Plan hat
  das als Risiko 2 vorab benannt (*„Wer nur `ahead` hebt, muss sicher sein, dass
  **kein Netz-Ausfall in diesem Zweig landet**"*) und mit einer Ausfall-Form für
  erledigt erklärt. Der Workflow-Kopf benennt die Kosten selbst:
  *„ein DAUERROTER Nachtlauf ist wieder derselbe verwaiste Sensor"*.
- **verifizierbar:** ja — isolierte Kopie mit `file://`-Quelle und gekürzter
  Antwort, Ausgabe oben; kein Gate hält den Satz, weil er Prosa ist.
- **klasse:** `fail-open-zusage-aus-einer-ausfallform-verallgemeinert`

### F-2 · MEDIUM · Die Grenzen-Liste wächst um zwei Einträge und lässt genau die Lücke aus, die den Slice ausgelöst hat: Currency misst **Release-Objekte**, nicht Tags

- **kategorie:** MEDIUM — Anker 18 (*Grenzen-Liste ohne ihre größte Lücke*),
  beide Griffe der Arbeitsanweisung angewandt: Vertrag umgedreht **und** Code
  statt Prosa gelesen.
- **quelle:** `AGENTS.md` §5 · Registereintrag
  `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen` · Baseline `v6.5.0` ·
  `regelwerk/modul-13-quality-gates.md` §Hard Rule (*„Ein Gate ohne seine Grenze
  behauptet ebenfalls zu viel"*)
- **pfad:** `harness/sensors/baseline-freshness.md` §Vertrag (A) und
  §Grenze 1–5 · `tools/harness/fetch-baseline-cache.sh:238-248`
- **befund:** Der Vertrag sagt *„die Release-**Liste** des Kurs-Repos gegen den
  Pin"*. Umgedreht folgt daraus eine Grenze, die unten fehlt — und zwar in zwei
  gemessenen Ausprägungen:
  **(a)** Die Currency-Hälfte sieht **Release-Objekte**, nicht **Tags**. Ein
  upstream angelegter Tag, für den (noch) kein Release-Objekt existiert oder das
  als Entwurf geführt wird, ist unsichtbar — und dieser Fall meldet
  `current`/**Exit 0**, also *still grün*, nicht `ahead`. Das ist exakt die
  Beobachtung, deren Untersuchung dieser Slice ist: `BEO-HARN`-Evidenz aus
  slice-193 (*„meldete weiterhin nur `v5.19.0`/`v5.20.0` als neuere Tags, obwohl
  `v6.0.0` … bereits publiziert war"*).
  **(b)** Vor dem Vergleich filtert das Skript die Liste mit
  `grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$'`. Gemessen fallen dabei `v7.0.0-rc1`,
  `v7.0.0-beta.1`, `templates-v8`, `V7.0.0` und `v7.0` heraus. Der Vertrag
  begründet die Wahl des Listen-Endpunkts aber ausdrücklich damit, dass
  `releases/latest` *„Prereleases überspringt"* — ein Prerelease mit
  suffigiertem Tag-Namen überspringt der Filter genauso. `AGENTS.md` §4 trägt
  dieselbe Begründung wörtlich weiter.
  **Versagensszenario:** Upstream veröffentlicht `v7.0.0-rc1` als neuesten
  Release. `newest` bleibt `v6.5.0`, der Zweig ist `current`, der Nachtlauf
  meldet *„Pin v6.5.0 ist der neueste Release-Tag"* und **Exit 0** — während der
  Vertrag zwei Absätze weiter oben Prerelease-Deckung verspricht. Das neu
  hinzugefügte Fenster-Grenze (Grenze 4) ist gegenüber (a) und (b) der
  **kleinere** Fall: es greift erst bei > 100 Releases Rückstand, (a) ist bereits
  einmal eingetreten und belegt.
- **verifizierbar:** ja — `grep -E` gegen die fünf Tag-Formen (Ausgabe oben);
  Vergleich `releases`- gegen `tags`-Endpunkt (heute deckungsgleich, aber das
  ist eine Aussage über heute, nicht über die Zusage). Kein Gate hält die
  Vollständigkeit einer Grenzen-Liste.
- **klasse:** `grenzen-liste-nennt-die-eigene-anlass-luecke-nicht`

### F-3 · MEDIUM · Der Ausschluss der Ursache steht auf einer Messung von **heute**; die verfügbare entscheidende Evidenz liegt im Eintrag selbst und wurde nicht benutzt

- **kategorie:** MEDIUM — Anker 17 (*Messung zählt einen Proxy statt des
  Gegenstands*: gemessen ist der heutige Listen-Zustand, Gegenstand ist der
  Zustand am Tag der Beobachtung).
- **quelle:** `AGENTS.md` §5 · Registereintrag
  `BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand` · Baseline `v6.5.0` ·
  `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register (*Stand trägt
  Zustand und Beleg*)
- **pfad:** `BEO-HARN/check-latest-blind-before-pin/state.md` (ganze Zeile) ·
  slice-215 §1, erster Abgrenzungs-Punkt · Commit `73f1098a`, vorletzter Absatz
- **befund:** Der Beleg lautet: *„Release- und Tag-Liste des Kurs-Repos sind
  heute deckungsgleich (58/58), die naheliegende Erklärung ‚Tag ohne
  Release-Objekt‘ trägt also nicht."* Die Zahl stimmt — ich habe sie
  nachgemessen (58 zu 58, `comm` in beide Richtungen leer). Der **Schluss**
  trägt nicht: Ein Release-Objekt, das nach der Beobachtung für einen schon
  vorhandenen Tag angelegt wird — oder ein Entwurf, der später veröffentlicht
  wird —, erzeugt genau die heutige Deckungsgleichheit. Heutige Kongruenz kann
  eine damalige Inkongruenz nicht ausschließen; alle 58 Objekte tragen heute
  `"draft": false`, was denselben Nachweis nicht liefert.
  Zugleich liegt die **entscheidende** Evidenz im Eintrag selbst und ist
  ungenutzt geblieben: Die Evidence-Datei sagt, es seien `v5.19.0`/`v5.20.0`
  **als neuere Tags gemeldet** worden. Damit ist bewiesen, dass der Pin in der
  Liste **gefunden** wurde (`newer` ist nur nicht leer, wenn der `sed`-Bereich
  ab der Pin-Zeile greift) — die Fenster-Hypothese, die der Slice als möglichen
  gemeinsamen Grund erwägt, ist dadurch ausgeschlossen, und zwar ohne Netz.
  Nachgemessen: Der Pin war damals `v5.18.0` (`MR-058` → `MR-060`), heute
  Position 42 von 51 in der strikten Liste, und zum Zeitpunkt der Beobachtung
  lagen **genau drei** Releases über ihm (`v5.19.0`, `v5.20.0`, `v6.0.0`) — vom
  100er-Fenster weit entfernt. Übrig bleibt genau die Hypothese aus F-2 (a),
  die der Eintrag als „trägt nicht" abgelegt hat.
  **Versagensszenario:** Der Eintrag steht auf 1× mit *„Ursache nicht
  feststellbar"*. Tritt der Fall wieder auf, sucht der nächste Lauf die Ursache
  erneut — und findet in `state.md` die Auskunft, die naheliegende Erklärung sei
  bereits ausgeschlossen.
- **verifizierbar:** ja — Release- gegen Tags-Endpunkt (58/58), Positions- und
  Abstandsmessung `v5.18.0`↔`v6.0.0` (drei Releases), `MR-058`/`MR-060` für den
  damaligen Pin. Kein Gate prüft die Tragfähigkeit eines Belegs.
- **klasse:** `heutiger-zustand-als-beleg-fuer-einen-damaligen`

### F-4 · MEDIUM · Die Exit-3-Semantik ist in zwei von fünf Spiegeln nachgezogen; einer der drei verbliebenen sagt jetzt etwas Falsches

- **kategorie:** MEDIUM — Basis LOW (*Doku-Drift*), **eine Stufe eskaliert**
  nach §Kontext-Eskalation: Der falsche Spiegel steht am **Bindepunkt** des
  Sensors, also in dem Artefakt, das ein Operator beim roten Lauf zuerst liest.
- **quelle:** Maintainability. `MR-025` beschreibt dieselbe Praxis
  (*„Vor dem Editieren wird die Liste der Spiegel dieser Semantik
  aufgeschrieben"*) und führt in seiner Spiegel-Tabelle ausdrücklich
  *„Autoritäts-Doku | `AGENTS.md`, `harness/README.md` (Gate-Beschreibungen)"* —
  **als Analogie zitiert, nicht als Autorität**: sein Geltungsbereich zählt
  *Grund-Code, Algorithmus-Schritt, Config-Schlüssel, Schwellenwert,
  Erkennungs-Form* auf, ein Skript-Exit-Code steht dort nicht.
- **pfad:** `.github/workflows/upstream-drift.yml`, Kopf-Abschnitt
  *WAS ER MELDET* · `AGENTS.md` §4, Zeile `make baseline-freshness` ·
  `Makefile:313`
- **befund:** Nachgezogen sind der Skript-Kopf und `harness/sensors/baseline-freshness.md`.
  Nicht nachgezogen sind drei Stellen, die dieselbe Semantik tragen:
  **(a)** Der Workflow-Kopf sagt *„Exit 3 heisst ‚Pin und Upstream sind
  verschieden‘ … VERALTET bei den Versions-Achsen (ihre Reihen sind monoton,
  ‚anders‘ heisst dort ‚neuer‘)"*. Für die Baseline-Achse ist das seit diesem
  Commit **falsch**: Exit 3 kann jetzt *unbestimmt* heißen, und „neuer" ist dann
  gerade nicht gemeint. **(b)** `AGENTS.md` §4 beschreibt die Currency-Hälfte
  weiterhin als *„neuerer Release-Tag"*. **(c)** Der Target-Kommentar in
  `Makefile:313` ebenso.
  **Versagensszenario:** Der Nachtlauf ist rot. Wer den Workflow-Kopf als
  Legende liest — er ist als solche geschrieben —, schließt auf einen
  verfügbaren Bump und beginnt eine Pin-Hebung, während der Lauf sagt, dass der
  eigene Pin nicht auffindbar ist. `make gate-consistency` fängt das nicht: Es
  hält die **Existenz** der Deklaration, nicht ihre Aussage.
- **verifizierbar:** ja — `grep -n "Exit 3" .github/workflows/upstream-drift.yml`
  gegen `harness/sensors/baseline-freshness.md` §Ausgabe und Ausgänge.
  Kein Gate; die Spiegel sind Prosa.
- **klasse:** `exit-semantik-ohne-spiegel-nachzug`

### F-5 · LOW · Grenze 4 nennt eine **notwendige** Bedingung als zweite **hinreichende**

- **kategorie:** LOW — Doku-Drift in einer neu geschriebenen Grenze.
- **quelle:** `AGENTS.md` §5 (Form vor der Aussage)
- **pfad:** `harness/sensors/baseline-freshness.md` §Grenze 4, letzter Satz
- **befund:** *„… tritt ein, sobald der Pin um mehr als 100 Releases zurückfällt
  **oder** das Repo diese Zahl überschreitet."* Der zweite Disjunkt allein löst
  den Fall nicht aus: Ein Repo mit 400 Releases und einem Pin auf dem
  zweitneuesten liefert weiterhin `current`. Auslösend ist ausschließlich der
  erste Teil; die Release-Zahl ist dafür notwendige Vorbedingung, nicht
  alternativer Trigger.
  **Versagensszenario:** Beim Überschreiten der 100 wird eine Prüfung
  eingeplant, die nicht nötig ist — oder, umgekehrt gelesen, für erledigt
  gehalten, weil das Repo unter 100 liegt, während der Pin sehr weit
  zurückfällt.
- **verifizierbar:** ja — gegen `fetch-baseline-cache.sh:246-248`, wo die
  Zugehörigkeit des Pins zur Liste allein entscheidet.
- **klasse:** `grenze-nennt-notwendige-bedingung-als-hinreichende`

### I-1 · INFO · `gemessen:` im Kopf-Kommentar ist kein auflösbares Herkunfts-Feld

- **pfad:** `tools/harness/fetch-baseline-cache.sh:33-34`
- **befund:** Die Zusagen-Hälfte des neuen Kommentars trägt Klassen (Zusage:
  Exit-Bedeutung; Abgrenzung: welcher Ausfall in welchen Zweig fällt) — insoweit
  ist §3.7 erfüllt. Der angehängte Beleg *„-- gemessen: leere Liste faellt in
  den skip-Zweig …"* ist dagegen Herkunfts-Prosa; §3.7 lässt Herkunft nur als
  **ein** auflösbares Feld nach Baseline-Schema zu (`DC-*`, `ADR-*`, `MR-*`,
  `seit welle-<NN>`/`seit slice-<NNN>`). Ein Beleg, den man im Kommentar nicht
  nachschlagen kann, altert unbemerkt — F-1 ist der Beleg dafür, dass er es
  bereits getan hat.

### I-2 · INFO · Ein Pin mit Prerelease-Suffix in der §Baseline-Zeile würde stillschweigend auf `vX.Y.Z` gekürzt

- **pfad:** `tools/harness/fetch-baseline-cache.sh:63-64` (Extraktion), `:66` (Validierung)
- **befund:** Als **Argument** wird `v6.5.0-rc1` sauber abgewiesen
  (`ungültiger/leerer Tag`, Exit 1 — geprüft). Aus der `**Stand:**`-Zeile wird
  der Tag dagegen mit `grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+'` **herausgeschnitten**;
  aus `v6.6.0-rc1` wird `v6.6.0` (geprüft). Ein solcher Pin auditierte
  wortlos einen anderen Tag als den gepinnten. Heute folgenlos — das Repo pinnt
  nur volle Releases —, aber es ist dieselbe Prerelease-Achse wie in F-2 (b).

### I-3 · INFO · Der Pin wird als `sed`-Regex interpoliert

- **pfad:** `tools/harness/fetch-baseline-cache.sh:245`
- **befund:** `sed -n "/^${tag}\$/,\$p"` behandelt die Punkte des Tags als
  Metazeichen. Bei den heutigen Tag-Formen folgenlos; die Bemerkung steht hier,
  weil dieser Ausdruck jetzt darüber entscheidet, ob der Nachtlauf 0 oder 3
  liefert.

### I-4 · INFO · `seit slice-215` steht im `state.md` eines Eintrags, dessen Stand `offen` ist

- **pfad:** `BEO-HARN/check-latest-blind-before-pin/state.md`
- **befund:** Der Herkunfts-Anker gehört zum Ausgang *verkörpert*; hier steht er
  in einem Eintrag, der ausdrücklich **nicht** verkörpert ist und bei 1× bleibt.
  Er löst korrekt auf (§Grenze 5 trägt ihn), verbindet aber eine Regel mit einer
  Beobachtung, die sie nicht ausgelöst hat.

---

## Negativbefunde (geprüft, ohne Befund)

- **Bruch-Test, beide Richtungen** — selbst gefahren, nicht nachgelesen: alte
  Fassung Exit 0, neue Fassung Exit 3, identischer Aufruf, echte Ausgaben oben.
  DoD (2) trägt in dem Umfang, den es beansprucht.
- **Regression** — normaler `--check-latest`-Lauf weiterhin Exit 0 mit
  `current`+`ok`; `--selftest` weiterhin `ok (9 Proben)`, Exit 0. Die
  Verhaltens-Änderung ist auf den `ahead`-Zweig begrenzt.
- **Der Konsument liest wirklich nur den Exit-Code** — am Workflow geprüft, nicht
  angenommen: `upstream-drift.yml` fährt `run: make baseline-freshness`, ohne
  `continue-on-error`, ohne Log-Auswertung, ohne `grep` auf stderr; keine
  weitere Stelle im Repo ruft `--check-latest`
  (`grep -n baseline-freshness Makefile *.mk`). Die Behauptung im Plan §3 und in
  der Commit-Botschaft trägt.
- **`AGENTS.md` §3.6 — keine Gate-Lockerung, keine ADR-Pflicht.** Der Pfad wird
  **verschärft** (0 → 3), keine Schwelle gesenkt, keine Prüfregel entfernt,
  `gates` unverändert (`baseline-freshness` war und bleibt außerhalb). §3.6 gilt
  der *Senkung*; hier greift sie nicht. Die faktische Verschiebung der
  fail-open-Linie für **eine** Ausfall-Klasse steht in F-1 — sie ist kein
  §3.6-Fall, sondern eine unbelegte Zusage.
- **`AGENTS.md` §3.1** — der Diff bringt keine Host-Toolchain ins Spiel; das
  Skript bleibt `bash` + `curl`, wie es die Netz-Targets ohnehin erwarten. Kein
  Go, kein Interpreter, keine neue Stage.
- **§1-Abgrenzung gehalten**, alle drei Punkte: `baseline-freshness` ist nicht
  nach `gates` gezogen (Makefile im Diff unberührt), die übrigen
  Freshness-Achsen sind unberührt (Diff umfasst fünf Dateien), und die Ursache
  der Anlass-Beobachtung wird nicht *bewiesen* — dass ihr **Ausschluss** zu weit
  geht, ist F-3 und keine Abgrenzungs-Verletzung.
- **§8 vollständig und formgerecht** — drei Vorprüfungen vorhanden, zwei
  Sub-Areas mit je einem Modus-Block, beide GF. Beide `d-check:cite`-Spannen
  **wortgleich** gegen `v6.5.0` · `regelwerk/modul-05-planning-harness.md`
  Zeilen 268–269 bzw. 274 nachgeprüft; der dritte Block (Nachtlauf) trägt
  regelkonform keine.
- **Register-Zahl 38 nachgezählt** —
  `find …/observations -mindepth 2 -maxdepth 2 -type d | wc -l` → `38`, zwei
  Kürzel (`BEO-ALL`, `BEO-HARN`). Die Aussage, `BEO-HARN` sei einzeln geöffnet
  worden, ist plausibel gemacht: Das Verzeichnis enthält genau einen Eintrag,
  und der Plan nennt ihn als solchen.
- **„58 Releases von 100" und „58/58" sind korrekt gemessen** — nachgerechnet:
  58 `tag_name` in der Release-Liste, 58 Tags über den Tags-Endpunkt (Seite 2
  leer), `comm` in beide Richtungen ohne Differenz. Beanstandet ist in F-3 der
  **Schluss**, nicht die Zahl.
- **Der rohe Byte-Parse der API-Antwort ist heute sauber** — die 58 aus
  `grep -o '"tag_name"…'` gewonnenen Werte decken sich exakt mit dem
  Tags-Endpunkt; kein Phantom-Treffer aus einem Release-Text. Ein solcher
  könnte konstruktionsbedingt nur `newer` erzeugen, nie `ahead`.
- **Exit 3 für zwei Zustände trägt.** Für den erklärten Konsumenten
  (Exit-Code-only) ist die Kollision folgenlos, stderr unterscheidet, und die
  Ausgangs-Tabelle deklariert beide Bedeutungen. Auch die Ordnung
  „schlimmster Fall (4>3>0)" bleibt konsistent: `ahead`+`drift` ergibt 4, wie
  `newer`+`drift` es vorher schon tat — kein **neuer** Verdeckungsfall. Was
  daraus doch folgt, steht in F-4 und betrifft die Spiegel, nicht die Wahl.
- **Erfolgreiche API-Antwort → `ahead` ohne Befund:** durchgespielt und, bis auf
  F-1, nicht gefunden. Leere Release-Liste → `tags` leer → `skip` (nicht
  `ahead`) — der Kommentar sagt das richtig. Rate-Limit liefert bei GitHub
  403/429, `curl -f` bricht ab → `skip`. 200 mit fremdem Inhalt (Proxy, HTML)
  → kein `tag_name` → `skip`. Ein Pin, der die Tag-Regex nicht erfüllt,
  erreicht `check_latest` gar nicht (Exit 1, geprüft) — die Restlücke dieser
  Achse ist I-2. Der einzige gefundene Falsch-`ahead`-Pfad ist die **teilweise**
  erfolgreiche Antwort in F-1.
- **`state.md`-Form** — kein §3.7-Chronik-Befund. Die Haus-Form dieses Registers
  trägt erklärende Prosa im Stand-Feld (vergleichbar:
  `BEO-ALL/review-collection-misses-non-slice-filenames`,
  `BEO-ALL/eigene-menge-gemessen-fremde-behauptet`,
  `BEO-ALL/large-migration-exceeds-session-review-limit`); Zustand und Beleg
  stehen vorn. Beanstandet ist der **Inhalt** des Belegs (F-3), nicht die Form.
- **Kommentar-Klassen §3.7** — der neue Kopf-Kommentar trägt Zusage und
  Abgrenzung; keine Review-Historie, keine Slice-Nummer, kein Mess-Label, keine
  Deliberation über Verworfenes. Einzige Rest-Beanstandung ist der
  `gemessen:`-Beleg (I-1).
- **`AGENTS.md` §3.3 / `MR-013`** — der Beanspruchungs-Commit `fea736e0` bündelt
  `git mv` und Roadmap-Ruhemarker in einem Commit, wie `MR-013` es für diese
  Richtung verlangt; `make planning-check` grün.
- **Commit-Botschaften §5, im Übrigen** — die Gate-Angabe *„zehn Gates, 737
  Dateien, 0 Befunde"* stimmt (eigener Lauf: zehn Targets in der `[gates]`-Zeile,
  `737 Datei(en) geprüft, 0 Befund(e)`). Die Bruch-Test-Angaben stimmen
  wörtlich. Überdehnt ist ausschließlich der Gegenproben-Absatz (F-1).
- **Nicht berührt** und deshalb ohne Befund: Hexagon-Import-Richtung
  (`ADR-0005`, kein Go-Code im Diff), Netzzugriff außerhalb `external`
  (`DC-QA-03` gilt dem Produkt; `harness/sensors/baseline-freshness.md`
  §Grenze 3 sagt das korrekt), Inline-Suppressions (keine),
  Referenz-Richtung/Provenance-Marker (keine neuen Abwärts-Token; `matrix` grün),
  Negativtest zu einem neuen öffentlichen Vertrag (kein Produkt-Vertrag berührt,
  keine `DC-*`-Bindung — der Slice-Kopf sagt das).

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 4 | F-1, F-2, F-3, F-4 |
| LOW | 1 | F-5 |
| INFO | 4 | I-1, I-2, I-3, I-4 |

---

## Verdikt

**Blockiert** — die Verhaltens-Änderung selbst ist richtig, belegt und in beide
Richtungen reproduziert, aber die drei Zusagen, die sie tragen (fail-open
unberührt · Ursache ausgeschlossen · Grenzen benannt), reichen jede weiter als
ihre Messung: F-1 ist als Falsch-Rot-Pfad gemessen, F-2 lässt die belegte
Anlass-Lücke aus der frisch erweiterten Grenzen-Liste heraus, F-3 stützt einen
Ausschluss auf einen Zustand von heute, F-4 lässt drei Spiegel stehen, von denen
einer jetzt falsch ist.
