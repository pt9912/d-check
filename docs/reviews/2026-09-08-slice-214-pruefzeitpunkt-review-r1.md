# Review-Report: slice-214 — 2026-09-08

**Review-Art:** Code — geprüft wird der Diff gegen Slice-Plan, ADR-0010/ADR-0011
und die Hard Rules (`AGENTS.md` §3.1, §3.9, §5, §6). Der Slice liefert eine
**Messung** und eine **Grenzen-Zeile**; beide Klassen haben eigene Anker
(17 bzw. 18).

**Gegenstand:** slice-214 · Commit-Range `HEAD~3..HEAD`
(`a4217514`, `5db30cfa`, `9bc58977`, `98d9e5ab`) · geändert: der Slice-Plan,
`docs/plan/planning/in-progress/roadmap.md`, `harness/sensors/semgrep.md`

**Skill:** `.harness/skills/reviewer.md` @ 1.16.0 ·
**Modell:** claude-opus-5[1m] · **Datum:** 2026-09-08

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan slice-214 (Ziel/Abgrenzung §1, DoD §2, Inventur §3, Risiken §6,
  Vorprüfungen §8)
- ADR-0010 (hermetisches semgrep-Gate, *„einmalig in einen Cache"*),
  ADR-0011 (Digest-Pins aller Build- und Gate-Images)
- `AGENTS.md` §3.1 (Docker/make-only, Supply-Chain-Defense), §3.9
  (SHA-gepinnte Action-Referenzen), §3.7/§5 (gemessene Menge · Form vor der
  Messung · Grenze gegen den Gegenstand), §6
- `MR-011`-Kette, `MR-053`, `MR-054`, `MR-013`
- Vorherige Findings am gleichen Gegenstand: slice-212-Review R1 (dort I-2,
  `pin-ohne-integritaets-gegenprobe`, ist der Anlass dieses Slice; dort F-3,
  `inventur-behauptet-vollstaendigkeit-die-die-stichprobe-bricht`)

**Eigene Läufe:** `make gates` → **Exit 0**, zehn Gates grün, vier
d-check-Läufe je `731 Datei(en) geprüft, 0 Befund(e)`, `coverage-gate: OK —
Coverage 94.60% erfüllt Schwelle 93%`, `Ran 55 rules on 63 files: 0 findings`.
`make doc-check` → **Exit 0**, `731 Datei(en) geprüft, 0 Befund(e)`.
`make nightly-state` → beide Nachtläufe `gruen`, Zeitstempel
`2026-09-07T05:33:45Z` / `2026-09-07T08:21:32Z`.

---

## Findings

### F-1 · HIGH · Die Inventur zählt sieben und lässt mindestens fünf Artefakte aus, die ihre **eigene** Form erfüllen — die daraus gezogene „einzige"-Aussage steht jetzt in einem Gate-Vertrag

- **kategorie:** HIGH — Basis MEDIUM (Anker 17 *Messung zählt einen Proxy*: gezählt
  sind drei Fundorte, nicht die in §3 definierte Form), **eine Stufe eskaliert**
  nach §Kontext-Eskalation: Der Schluss aus der Zahl ist eine
  Supply-Chain-Aussage und steht seit `98d9e5ab` als **Grenze** in einem
  Gate-Vertrag.
- **quelle:** `AGENTS.md` §5 (*Vor einer Messung steht die Form ihres
  Gegenstands*), `AGENTS.md` §3.9, ADR-0011 §Entscheidung 1
- **pfad:** slice-214 §3, Inventur-Tabelle Z. 115–124 und Fließtext Z. 111–113
  („**Sieben** Artefakte erfüllen die Form aus §3") ·
  `harness/sensors/semgrep.md:33-34`
- **befund:** §3 definiert die Form vorab und korrekt: nicht im Repo entstanden ·
  an fester Kennung bezogen · geht in einen Lauf ein. Gezählt wurde danach aber
  nicht gegen die Form, sondern gegen drei Fundorte (`Dockerfile`, `a-check.mk`,
  `tools/semgrep.sh`). Die Form erfüllen zusätzlich, im Bestand nachweisbar:
  **(a)** die drei SHA-gepinnten Actions `actions/checkout@3d3c42e5…`,
  `docker/login-action@dbcb8138…`, `peter-evans/dockerhub-description@1b9a80c0…`
  (`.github/workflows/`, fünf Fundstellen) — sie sind Fremd-Code an einer festen
  Kennung, der in **jeden** CI-, Release- und Nachtlauf eingeht, `AGENTS.md`
  §3.9 macht sie zur Hard Rule mit derselben Supply-Chain-Begründung wie §3.1,
  und `AGENTS.md` §4 führt für genau diese drei Pins drei eigene
  Freshness-Achsen; **(b)** das digest-gepinnte Trivy-Image
  (`tools/image-scan.sh:126`, `TRIVY_DIGEST` Z. 50), für das `AGENTS.md` §4
  ebenfalls beide Achsen (`make freshness-trivy`, `make trivy-digest`) führt;
  **(c)** `tools/archive-wave/Dockerfile:19` mit einem **anderen**
  `golang`-Digest (`sha256:4013ae0f…`) als der gezählten Zeile
  (`sha256:512690a5…`) und Z. 52 mit dem distroless-Digest. Aus der zu kleinen
  Menge folgt der Satz, der jetzt in `harness/sensors/semgrep.md` steht: der
  Regel-Cache sei *„das einzige gepinnte Fremd-Artefakt dieses Repos"*, das
  außerhalb des Repos liegt und von keinem Arbeitsbaum-Gate gesehen wird. Der
  Docker-Store und der Actions-Cache eines Runners liegen ebenso außerhalb des
  Repos und werden von keinem Arbeitsbaum-Gate gesehen.
  **Versagensszenario:** Wer die Grenze liest, schließt, für alle übrigen
  gepinnten Fremd-Artefakte gebe es einen Träger. Für die drei Action-Pins gibt
  es genau einen — `make workflow-pins`, und der prüft nach `AGENTS.md` §3.9
  ausdrücklich die **Form**, nicht die Gültigkeit des SHA.
- **verifizierbar:** ja — `grep -rn "uses:" .github/workflows/*.yml` (fünf
  fremde Referenzen, drei Actions), `grep -n "sha256:" tools/image-scan.sh`,
  `grep -n "^FROM" tools/archive-wave/Dockerfile`. Kein Gate hält die
  Vollständigkeit einer Inventur; das sagt der Slice selbst nicht.
- **klasse:** `inventur-behauptet-vollstaendigkeit-die-die-stichprobe-bricht`

### F-2 · HIGH · Unter der in §3 selbst gesetzten Definition von *Prüfzeitpunkt* trägt die Tabelle drei verschiedene Antworten für **einen** Mechanismus — und damit fällt die Achse des Slice

- **kategorie:** HIGH — Basis MEDIUM (Anker 17), eskaliert wie F-1: Die Aussage
  ist die tragende des Slice und steht als Grenze in einem Gate-Vertrag.
- **quelle:** `AGENTS.md` §5 · ADR-0011 §Entscheidung 2 · Baseline `v6.5.0` ·
  `regelwerk/modul-13-quality-gates.md` §Hard Rule (*„Ein Gate ohne seine Grenze
  behauptet ebenfalls zu viel"*)
- **pfad:** slice-214 §3 Z. 91–94 (Definition), Z. 117–121 (Tabelle),
  Z. 129–132 · `harness/sensors/semgrep.md:26`
- **befund:** §3 definiert: *„Der Moment, in dem etwas die Bindung
  **tatsächlich nachrechnet**"* und sagt zwei Zeilen später selbst, wo das ist:
  *„nachgerechnet wird sie von Docker beim **Pull**"*. Die Tabelle schreibt für
  dieselbe Mechanik dann **jeder Image-Bau** (drei `Dockerfile`-Zeilen) und
  **jeder Lauf** (`a-check`-, `semgrep`-Image); die Commit-Botschaft von
  `98d9e5ab` schreibt eine dritte Fassung: *„werden bei jedem Bezug
  nachgerechnet"*. Drei Antworten, ein Mechanismus. Gemessen im `make
  gates`-Lauf dieses Reviews: vier Image-Bauten, **kein** Pull — die Zeilen
  lauten `[internal] load metadata for …@sha256:…` und `CACHED` (22×). Ein Bau
  ist kein Bezug, und ein `docker run` auf ein lokal vorhandenes Image erst
  recht nicht. Nach der Definition des Slice ist der Prüfzeitpunkt aller fünf
  Images damit **derselbe wie beim Regelset**: einmal beim ersten Bezug auf
  dieser Maschine, danach Vertrauen in einen Store außerhalb des Repos. Genau
  darauf ruht aber die Achse des Slice (*„Zwei Artefakte fallen aus der Reihe"*,
  *„das einzige mit einmaligem Prüfzeitpunkt"*) und der Kopfsatz der neuen
  Grenze.
  **Versagensszenario:** Eine spätere Hebung eines der fünf Digests wird als
  risikoarm behandelt („wird ohnehin bei jedem Lauf nachgerechnet"), obwohl
  zwischen Pull und Hebung nichts nachrechnet — dieselbe Lage, die der Slice für
  den Regel-Cache ausdrücklich als Lücke aufschreibt.
- **verifizierbar:** ja — `make gates` und den Bau-Log lesen: `#5 [internal]
  load metadata …`, `#10 [deps 1/5] FROM …`, `CACHED`; keine `pulling`-Zeile.
- **klasse:** `pruefzeitpunkt-spalte-nennt-drei-werte-fuer-einen-mechanismus`

### F-3 · MEDIUM · Die Selbstbeschränkung *„nicht in diesem Repo gemessen"* deckt die schwachen Fremd-Aussagen und lässt die starke aus

- **kategorie:** MEDIUM (Anker 8 — die Messung stimmt, ihr Schluss reicht
  weiter; zusätzlich Anker 17 — der Superlativ ist eine Ordnung ohne
  ausgeschriebenes Kriterium)
- **quelle:** `AGENTS.md` §5 (beide Absätze, `seit slice-210` und
  `seit slice-213`) · slice-214 §6, zweites Risiko
- **pfad:** slice-214 §3 Z. 141–148 · `harness/sensors/semgrep.md:27-29`,
  `:36-37`
- **befund:** §3 quarantäniert korrekt zwei Fremd-Aussagen (*„Docker rechnet den
  Digest nach"* und die Modul-Prüfsummen) als *„in diesem Repo nicht gemessen"*.
  Die dritte Fremd-Aussage trägt keinen solchen Marker und ist die stärkste des
  Slice: *„ein Commit-SHA ist ein Hash über den Baum"* und, in der Sensor-Datei,
  *„hier beweist sie sie [die Echtheit]"*. Dass `git fetch --depth 1 origin
  <sha>` + `checkout FETCH_HEAD` den Inhalt gegen den SHA bindet, ist die
  dokumentierte Eigenschaft von git — dieselbe Klasse wie die
  Docker-Eigenschaft, und in diesem Repo ebensowenig gemessen (kein
  Bruch-Test, keine Probe). Gemessen ist nur, dass der Pin **im Skript steht**
  (`tools/semgrep.sh:24`). Zweitens trägt der Satz einen Superlativ über die
  ganze Menge — *„die **stärkste** Bindung im ganzen Pin-Bestand dieses Repos"*
  —, dessen Kriterium nirgends ausgeschrieben ist, obwohl §3 die Form jeder
  anderen Größe vorab ausschreibt. Unter dem im Satz selbst genannten Kriterium
  („Hash über den Inhalt statt über eine mitgelieferte Liste") liegen die sechs
  `@sha256:`-Digests gleichauf; unter dem naheliegenden zweiten (Stärke der
  Hashfunktion) ist der git-Objektname SHA-1 und damit der **schwächste**
  Primitiv im Bestand.
  **Versagensszenario:** Der Satz wird beim nächsten Pin-Entscheid als gemessene
  Rangordnung zitiert und begründet, den Regel-Cache anders zu behandeln als die
  Digest-Pins — auf einer Ordnung, die niemand aufgestellt hat.
- **verifizierbar:** nein — kein Gate; am Text und an `tools/semgrep.sh`
  ablesbar.
- **klasse:** `fremde-aussage-ohne-den-marker-den-die-nachbarzeile-traegt`

### F-4 · MEDIUM · *„Genau einmal geprüft"* gilt einem warmen Cache; auf dem CI-Pfad — dem einzigen klon-unabhängigen Träger des Gates — wird bei **jedem** Lauf geholt und damit geprüft

- **kategorie:** MEDIUM (Anker 18, zweiter Griff: *wo der Gegenstand Code oder
  Konfiguration ist, gegen diese prüfen* — hier gegen `ci.yml`)
- **quelle:** `AGENTS.md` §5 (*Wer eine Grenze aufschreibt, prüft sie gegen den
  Gegenstand*) · ADR-0010
- **pfad:** `harness/sensors/semgrep.md:26`, `:30-31` · slice-214 §3 Z. 122
  (Tabellenzelle *„**einmalig** — danach nur Verzeichnis-Existenz"*)
- **befund:** Der Cache-Pfad ist `$RULES_DIR = $CACHE_ROOT/$RULES_COMMIT`
  (`tools/semgrep.sh:31-32`) und `$CACHE_ROOT` ist über `SEMGREP_RULES_CACHE`
  bzw. `XDG_CACHE_HOME` verschiebbar. *„Genau einmal"* heißt damit **einmal je
  Pin und je Cache-Ort** — eine Pin-Hebung erzeugt einen neuen Pfad und löst
  denselben Bezug samt Prüfung erneut aus. Entscheidender ist der zweite Ort:
  `.github/workflows/ci.yml` läuft `runs-on: ubuntu-latest` und fährt
  `make ci` ohne jeden Cache-Schritt für `~/.cache` — auf jedem CI-Lauf ist der
  Cache kalt, wird geholt und dabei gegen den Pin geprüft. Für den Pfad, auf dem
  das Gate klon-unabhängig durchgesetzt wird (`harness/README.md`
  §Durchsetzungsgrenzen: *„Der klon-unabhängige Boden ist die PR-/Push-CI"*),
  ist der Prüfzeitpunkt also **jeder Lauf** — das Gegenteil der Zeile. Die
  Grenze nennt diesen Ort nicht.
  **Versagensszenario:** Die Zeile wird als Argument für eine
  Re-Verifikation gelesen („nie wieder geprüft") und der Aufwand landet an der
  Stelle, an der er schon läuft, während der wirklich ungeprüfte Ort — der warme
  Entwickler-Cache — davon nicht profitiert.
- **verifizierbar:** ja — `.github/workflows/ci.yml` enthält keinen
  Cache-Schritt; in einem kalten CI-Log steht die Zeile `semgrep: hole
  Regel-Cache am Pin … (einmalig, Netz; …)`, im lokalen Lauf dieses Reviews
  fehlt sie (warmer Cache unter
  `~/.cache/d-check/semgrep-rules/d41fb34c…/go/lang/security`).
- **klasse:** `grenze-gilt-nur-fuer-den-warmen-pfad`

### F-5 · MEDIUM · Die Grenze nennt als Bedingung die Existenz des **Cache-Verzeichnisses**; das Skript testet das **Subset**-Verzeichnis, und der Reparaturpfad dafür ist gebrochen

- **kategorie:** MEDIUM (Anker 18, zweiter Griff — gemessen stand in einem von
  sieben Fällen die echte Grenze nur im Code; hier ist es wieder so)
- **quelle:** `AGENTS.md` §5 · ADR-0010
- **pfad:** `harness/sensors/semgrep.md:30-32` · `tools/semgrep.sh:36`, `:44`
- **befund:** Die neue Grenze sagt: *„Die Bedingung für einen erneuten Bezug ist
  die **Existenz** des Cache-Verzeichnisses"*. Das Skript testet
  `[ ! -d "$RULES_DIR/$RULES_SUBSET" ]`, also das **Subset**-Verzeichnis
  `go/lang/security` **innerhalb** des Cache-Verzeichnisses — der
  slice-212-Review hatte diese Unterscheidung in I-2 noch genau so notiert, die
  neue Fassung verliert sie. Der Unterschied ist nicht kosmetisch: Existiert
  `$RULES_DIR`, fehlt aber das Subset, greift der Bezugs-Zweig und endet mit
  `mv "${RULES_DIR}.tmp" "$RULES_DIR"` — POSIX-`mv` verschiebt dann **in** das
  bestehende Verzeichnis hinein. Gemessen mit einer coreutils-Probe: aus
  `RD.tmp/go/lang/security` wird `RD/RD.tmp/go/lang/security`; das Subset
  entsteht nie an der erwarteten Stelle. Der Zustand repariert sich also nicht,
  sondern wiederholt sich bei jedem Lauf und schachtelt tiefer. Die Grenze
  beschreibt damit ein Verhalten („Verzeichnis weg ⇒ Bezug und Prüfung erneut"),
  das das Skript nur für den Fall hat, in dem das **ganze** Commit-Verzeichnis
  fehlt.
  **Versagensszenario:** Wer die Prüfung nach der Grenzen-Zeile erzwingen will
  und das Subset löscht, bekommt einen endlosen Netz-Bezug und ein
  `semgrep: FEHLER — 0 Regeln geladen` (laut, nicht still — deshalb kein
  HIGH), ohne dass die Ursache irgendwo steht.
- **verifizierbar:** ja — `mkdir -p RD RD.tmp/go/lang/security && mv RD.tmp RD
  && find RD` zeigt die Schachtelung; kein Gate deckt den Zweig.
- **klasse:** `grenze-beschreibt-mechanismus-der-nicht-existiert`

### F-6 · MEDIUM · §8 hält fest, für `HARN` stehe **kein** Eintrag im Register — es steht einer, und er betrifft genau die Achse dieses Slice

- **kategorie:** MEDIUM (Pflicht-Vorprüfung mit falschem Ergebnis; Baseline
  `v6.5.0` · `regelwerk/modul-05-planning-harness.md` §Zwei Schritte vor der
  Modus-Begründung: *„Steht eine der berührten Sub-Areas dort?"*)
- **quelle:** `AGENTS.md` §5 (drei Vorprüfungen) · Baseline `v6.5.0` ·
  `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register
- **pfad:** slice-214 §8 Z. 234–236 · Commit `5db30cfa` (Botschaft:
  *„Fuer das Kuerzel HARN steht kein Eintrag im Register"*)
- **befund:** Der Slice führt `tools/harness/` (Kürzel `HARN`) als zweite
  berührte Sub-Area und notiert die Nicht-Existenz eines Eintrags ausdrücklich
  als Ergebnis. `docs/plan/planning/observations/BEO-HARN/` enthält
  `check-latest-blind-before-pin` — **Stand `offen`**, ein Beleg (`slice-193`),
  Sub-Area `tools/harness/`. Die eigene Zählung im selben Absatz („**38**
  Verzeichnisse über beide Kürzel", nachgezählt: 37 unter `ALL` + 1 unter
  `HARN` = 38) setzt die Existenz dieses Kürzels bereits voraus. Der Eintrag ist
  zudem nicht beliebig: Er beschreibt, dass
  `fetch-baseline-cache.sh --check-latest` einen neueren Release-Tag **nicht**
  sieht, solange der eigene Pin dahinter liegt — also eine benannte Blindheit
  genau des Trägers, auf den sich die Inventur-Zeile *„vendorte Baseline …
  Echtheit nur im Nachtlauf"* stützt.
  **Versagensszenario:** Das Kriterium *Evidenz-/Diskrepanz-Risiko* für
  `tools/harness/` steht auf **niedrig** mit der Begründung, die Sub-Area werde
  nur gelesen — die eine offene Beobachtung, die das qualifiziert hätte, ist
  nicht eingegangen; ihr Zähler bleibt bei 1, obwohl der Slice denselben Träger
  erneut berührt.
- **verifizierbar:** ja — `ls docs/plan/planning/observations/BEO-HARN/` und
  `cat …/check-latest-blind-before-pin/state.md`.
- **klasse:** `nein-aussage-woertlich-unzutreffend`

### F-7 · MEDIUM · Die §1-Abgrenzung begründet den CR-Ausschluss mit einer Warterichtung, die die CR-Datei umkehrt

- **kategorie:** MEDIUM (Anker 9 — *lies das Feld, nicht den Titel*; die Quelle
  ist echt, ihre Aussage die entgegengesetzte)
- **quelle:** `AGENTS.md` §5 (zitierte Quelle trägt nur ihren
  Geltungsbereich) · Baseline `v6.5.0` ·
  `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice (*„Je Punkt eine
  Begründung — ein Ausschluss ohne Grund ist eine Behauptung, keine Grenze"*)
- **pfad:** slice-214 §1 Z. 52–54
- **befund:** Der zweite Abgrenzungspunkt lautet: *„Der
  `ignore-refs`-/`exempt-paths`-CR. Er liegt unentschieden in `docs/plan/cr/`
  und wartet auf eine Antwort **des Absenders**"*. Die gemeinte Datei ist
  `2026-09-07-cr-eingehend-adopter-links-anchors-exempt-paths.md`; ihre
  Kopf-Felder sagen: *„**Richtung:** eingehend — dieses Repo ist der
  **Empfänger**, nicht der Bittsteller"* und *„**Stand:** eingegangen, **noch
  nicht entschieden**"*. Es wartet also der Absender auf **dieses** Repo. Die
  Begründung des Ausschlusses dreht die Bringschuld um.
  **Versagensszenario:** Der Punkt wird in Folge-Slices wortgleich übernommen —
  die Praxis der letzten vier Slices —, und der CR bleibt liegen, weil jeder
  Plan ihn als „liegt bei den anderen" ausweist.
- **verifizierbar:** ja — Kopf-Felder der CR-Datei; kein Gate prüft die
  Richtungs-Aussage.
- **klasse:** `abgrenzung-nennt-die-falsche-warterichtung`

### F-8 · MEDIUM · Das einzige Fremd-Image **ohne** Bindung fällt durch die Form-Definition heraus und wird nirgends genannt — obwohl es in jeden Bau eingeht

- **kategorie:** MEDIUM (Anker 18 — die Liste liest sich als Menge und lässt den
  Eintrag weg, der die anderen überwiegt)
- **quelle:** ADR-0011 §Entscheidung 1 (*„**Jedes** extern bezogene Image wird
  per `@sha256:`-Digest gepinnt"*) · `AGENTS.md` §3.1, §3.9 (*„ein Tag lässt
  sich umhängen, ein SHA nicht"*)
- **pfad:** slice-214 §3 Z. 84–89 (Form) und Z. 96–101 (*„Der vermutete Bestand
  ist klein"*) · `Dockerfile:1` · `tools/archive-wave/Dockerfile:1`
- **befund:** Beide Dockerfiles beginnen mit `# syntax=docker/dockerfile:1.7`.
  Das ist ein **fremdes Image**, das BuildKit bezieht und als Frontend
  **ausführt** — im `make gates`-Lauf dieses Reviews viermal sichtbar als
  `#2 resolve image config for docker-image://docker.io/docker/dockerfile:1.7`.
  Gebunden ist es an einen **beweglichen Tag**, nicht an einen Digest; nach der
  Form aus §3 („an einer festen Kennung bezogen — Digest, Commit-SHA,
  Tag+Manifest") fällt es damit aus der gemessenen Menge heraus, und kein Satz
  im Slice sagt das. Die Inventur liest sich dadurch als vollständige Antwort
  auf *„was prüft die Unversehrtheit und wann"*, während der eine Bezug ohne
  jede Antwort in ihr nicht vorkommt. ADR-0011 §Entscheidung 1 formuliert die
  Pflicht ohne Einschränkung auf `FROM`-Zeilen.
  **Versagensszenario:** Ein umgehängter `docker/dockerfile:1.7` ersetzt den
  Interpreter aller Dockerfiles dieses Repos — die höchste Stufe im Bau —, und
  weder `make workflow-pins` (nur `.github/workflows/`) noch `make versions`
  (nur `FROM`-Digests und das semgrep-Image) noch die Digest-Achsen (fünf
  benannte Images) sehen ihn.
- **verifizierbar:** ja — `head -1 Dockerfile`, `head -1
  tools/archive-wave/Dockerfile`, und die `resolve image config`-Zeilen jedes
  `make build`/`make gates`-Laufs.
- **klasse:** `form-definition-schliesst-die-schwaechste-bindung-aus`

### F-9 · LOW · Die Referenz-Zeile schreibt die Modul-Prüfsummen dem `-mod=readonly`-Schalter zu

- **kategorie:** LOW (latente Wartungsfalle; die Zeile ist ausdrücklich die
  nicht mitgezählte Referenz-Antwort)
- **quelle:** Maintainability · slice-214 §3 („Referenz-Antwort")
- **pfad:** slice-214 §3 Z. 124 · `Dockerfile:38`
- **befund:** Die achte Tabellenzeile nennt als Prüfzeitpunkt *„jeder Bau, über
  den read-only-Schalter im `Dockerfile`"*. Der Schalter ist
  `GOFLAGS="-mod=readonly -buildvcs=false"`; er verbietet der Toolchain,
  `go.mod`/`go.sum` zu **verändern**, und ist nicht der Mechanismus, der die
  Prüfsummen **verifiziert** — die Verifikation gegen `go.sum` hängt nicht an
  diesem Schalter. Die Zeile benennt damit den falschen Träger für die einzige
  Aussage, die sie macht.
  **Versagensszenario:** Ein späterer Bau-Refactor entfernt `-mod=readonly`
  (etwa für einen Tidy-Lauf) und stützt sich auf diese Zeile für die Annahme,
  damit falle die Prüfsummen-Prüfung weg — oder umgekehrt, sie bleibe erhalten,
  weil der Schalter bleibt.
- **verifizierbar:** ja — `grep -n GOFLAGS Dockerfile`.
- **klasse:** `mechanismus-der-falschen-konfigurationsgroesse-zugeschrieben`

### F-10 · LOW · Dieselbe Klasse steht in der Schwester-Datei an Position 0 mit Vorrang-Satz und hier an Position 3 ohne

- **kategorie:** LOW (Doku-Drift zwischen zwei Dateien derselben Klasse)
- **quelle:** Maintainability · Baseline `v6.5.0` ·
  `regelwerk/modul-13-quality-gates.md` §Hard Rule
- **pfad:** `harness/sensors/semgrep.md:26-38` gegen
  `harness/sensors/baseline-verify.md` §Grenze, Eintrag 0
- **befund:** slice-212 hat in der Schwester-Datei die tragende Grenze bewusst
  an **Position 0** gesetzt und den Vorrang ausgeschrieben (*„und das ist die
  Grenze, die alle folgenden überwiegt"*). Die neue Grenze beansprucht
  ausdrücklich die spiegelbildliche Bedeutung (*„Umgekehrt zu
  `baseline-verify`"*), wird aber ans **Ende** der Liste gehängt, ohne Aussage
  über ihren Rang. Zwei Dateien derselben Klasse, zwei Formen, kein benannter
  Grund.
  **Versagensszenario:** Wer die Liste von oben liest, nimmt Umfang (1) und
  Alterung (2) als die wesentlichen Grenzen und die einzige ohne jede
  Nachprüfung als Nachtrag.
- **verifizierbar:** nein — Form-Frage, an beiden Dateien ablesbar.
- **klasse:** `gleiche-klasse-zwei-formen-ohne-benannten-grund`

### I-1 · INFO · `make gates` ist auf kaltem Cache nicht netzlos — der Slice berührt die Stelle, an der das steht, ohne sie zu berühren

- **kategorie:** INFO (Bestand, nicht durch diesen Diff entstanden; als
  dokumentationswürdige Annahme notiert)
- **quelle:** ADR-0010 (Setup-Netz ausdrücklich entschieden) · `AGENTS.md` §3.1 ·
  `harness/README.md` §Sensors
- **pfad:** `AGENTS.md` §3.1 (*„Alle vier stehen bewusst außerhalb von
  `gates`"*) · `harness/README.md`, Zeile zu `make baseline-freshness` (*„dass
  `gates` netzlos bleibt, ist eine Eigenschaft dieses Repos"*)
- **befund:** `make semgrep` ist Bestandteil von `gates` und holt auf kaltem
  Cache über das Netz (`tools/semgrep.sh:37-44`). ADR-0010 und der Vertrag in
  `harness/sensors/semgrep.md` erklären das als Setup und decken es damit; die
  beiden Stellen, an denen die Netzlos-Zusage für `gates` **pauschal** steht,
  nennen den Fall nicht. Der Slice hat den Prüfzeitpunkt genau dieses Bezugs
  neu beschrieben und die Gelegenheit nicht genutzt. Kein Versagen im Bestand
  gemessen; notiert, weil die Ergänzung aus F-4 dieselbe Stelle betrifft.
- **verifizierbar:** nein — Textbefund.
- **klasse:** `pauschale-netzlos-zusage-neben-einem-setup-netz-zweig`

---

## Negativbefunde

- geprüft, ohne Befund: **`make gates`** — Exit 0, zehn Gates, vier
  d-check-Läufe je `731 Datei(en) geprüft, 0 Befund(e)`, Coverage 94.60 % über
  Schwelle 93 %, `Ran 55 rules on 63 files: 0 findings`. Die
  Commit-Botschaften nennen genau diese Zahlen; nachgezählt und deckungsgleich.
- geprüft, ohne Befund: **`make doc-check`** — Exit 0, `731 Datei(en) geprüft,
  0 Befund(e)`; die drei `d-check:cite`-Direktiven des Plans laufen im inneren
  Loop und sind grün.
- geprüft, ohne Befund: **`d-check:cite`-Anker in §8** — beide zeigen auf die
  **vorschreibende** Zeile: `modul-05-planning-harness.md:268-269` trägt
  *„Sub-Area-Wahl prüfen … drei Achsen, Schwelle ≥ 2"*, `:274` trägt *„Offene
  Beobachtungen sichten. Das"*; die Blockzitate darunter sind wortgleich.
  `MR-054` verlangt sie nur für die beiden kanonischen Blöcke — der dritte
  (Nachtlauf) trägt korrekt keine.
- geprüft, ohne Befund: **Nachtlauf-Vorprüfung (`MR-053`)** — eigener Lauf von
  `make nightly-state` bestätigt beide Zeitstempel wortgleich
  (`2026-09-07T05:33:45Z`, `2026-09-07T08:21:32Z`) und beide Ausgänge `gruen`.
- geprüft, ohne Befund: **Register-Zahlen unter `ALL`** — die vier zitierten
  Einträge tragen exakt die genannten Zähler (4 / 14 / 8 / 4 Evidence-Dateien)
  und die genannten Stände (`verkörpert`, `gemischt`, `geplant`, `verkörpert`);
  die Aussage *„Keiner der vier erreicht mit diesem Slice die Schwelle
  erstmalig"* trägt, alle vier liegen bereits darüber. Gesamtzahl 38
  nachgezählt und richtig. (Die Teil-Aussage zu `HARN` ist F-6.)
- geprüft, ohne Befund: **ADR-0010-Zitat in §1** — *„nennt den einmaligen Bezug
  ausdrücklich als Eigenschaft"* trägt: ADR-0010 §Entscheidung 1 (*„wird
  einmalig in einen Cache"*) und §3 (*„Das einmalige Holen des Caches am Pin ist
  Setup"*), zusätzlich in der Alternativen-Tabelle als Contra. Der Verweis
  bleibt im Geltungsbereich.
- geprüft, ohne Befund: **Ausschluss der Modul-Abhängigkeiten** — begründet,
  nicht bequem: Sie bleiben als achte Zeile **in** der Tabelle stehen, sind
  ausdrücklich nicht mitgezählt, und ihre Fremd-Aussage trägt den
  Quarantäne-Marker aus §3. (Der Träger in der Zelle ist F-9.)
- geprüft, ohne Befund: **§1-Abgrenzung, Punkte 1 und 3** — „keine
  Re-Verifikation einbauen" ist eingehalten (`tools/semgrep.sh` unverändert, der
  Diff berührt keinen Produkt-Code und kein Skript); „keine erneute Inventur der
  `## Grenze`-Abschnitte" ebenfalls (der Slice misst **eine** Achse über wenige
  Artefakte und fasst keine Grenzen-Abschnitte an außer dem einen, den er
  ergänzt).
- geprüft, ohne Befund: **Sub-Area-Wahl `tools/harness/` gegen das
  Inklusionskriterium** — die Sub-Area qualifiziert (Achse 1: eigene
  `MR`-Adaptionen bestehen — `MR-004`, `MR-005`, `MR-042`; Achse 3: eigenes
  Verzeichnis), ist in `harness/conventions.md` §Modus-Deklaration mit Kürzel
  `HARN` geführt und wird korrekt als *gelesen, nicht geändert* eingeführt; der
  GF-Modus ist richtig. Anmerkung ohne Finding-Rang: Von den vier Dateien, die
  §8 als ihren Gegenstand nennt, liegt genau **eine** unter `tools/harness/`
  (`fetch-baseline-cache.sh`); `tools/semgrep.sh`, `Dockerfile` und `a-check.mk`
  fallen unter den Repo-Default `*`. Der Modus-Schluss ändert sich dadurch
  nicht (beide GF).
- geprüft, ohne Befund: **Lifecycle und `MR-013`** — `9bc58977` ist ein reiner
  `git mv` plus dem gekoppelten Ruhe-Marker-Flip in der Roadmap, keine
  Inhaltsänderung an der Slice-Datei; `make planning-check` grün im
  `gates`-Lauf. Die Botschaft benennt korrekt, dass keine Pfad-Verweise
  nachzuziehen waren.
- geprüft, ohne Befund: **Traceability der drei Commits** — jede Botschaft nennt
  `slice-214` und weitere auflösbare Kennungen; `AGENTS.md` §5 erfüllt.
- geprüft, ohne Befund: **`git`-Bindung des Regelsets, sachlich** — der Bezug
  über `fetch --depth 1 origin <sha>` + `checkout FETCH_HEAD` bindet den Inhalt
  an den Commit-SHA; die inhaltliche Aussage ist richtig. Bemängelt ist in F-3
  ausschließlich, dass sie ohne den Marker steht, den die Nachbarzeile trägt,
  und der Superlativ ohne Kriterium.
- geprüft, ohne Befund: **Zustandsfelder im Diff** — der Slice trägt kein
  `**Status:**`-Feld (`**Lifecycle:**`-Hinweis korrekt), die Roadmap-Änderung
  entfernt nur den Ruhe-Marker; keine Chronik in einer Zustandszelle.
- geprüft, ohne Befund: **Kommentar-Klassen** — der Diff ändert keinen Code,
  keine Konfiguration und kein Skript; §3.7 hat in diesem Diff keinen
  Gegenstand.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 6 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:**
`inventur-behauptet-vollstaendigkeit-die-die-stichprobe-bricht` ·
`pruefzeitpunkt-spalte-nennt-drei-werte-fuer-einen-mechanismus` ·
`fremde-aussage-ohne-den-marker-den-die-nachbarzeile-traegt` ·
`grenze-gilt-nur-fuer-den-warmen-pfad` ·
`grenze-beschreibt-mechanismus-der-nicht-existiert` ·
`nein-aussage-woertlich-unzutreffend` ·
`abgrenzung-nennt-die-falsche-warterichtung` ·
`form-definition-schliesst-die-schwaechste-bindung-aus` ·
`mechanismus-der-falschen-konfigurationsgroesse-zugeschrieben` ·
`gleiche-klasse-zwei-formen-ohne-benannten-grund` ·
`pauschale-netzlos-zusage-neben-einem-setup-netz-zweig`

**Wiederholte Klassen:** `inventur-behauptet-vollstaendigkeit-die-die-stichprobe-bricht`
und `grenze-beschreibt-mechanismus-der-nicht-existiert` traten bereits im
slice-212-Review auf — zweite Instanz je Klasse, beide unter derselben
Beobachtung `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen` bzw.
`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand` einzuordnen. Beide
Beobachtungen sind bereits verkörpert (`seit slice-213` / `seit slice-210`) —
die Wiederholung **nach** der Verkörperung ist das eigentliche Signal: Die
Regeln waren im Plan zitiert und wurden auf den eigenen Gegenstand nicht
angewandt.

## Verdikt

**Merge-blockierend:** ja — zwei HIGH und sechs MEDIUM. F-1 und F-2 treffen
DoD (1) unmittelbar: Die Inventur ist weder in ihrer **Menge** vollständig noch
in ihrer **Aussage je Mitglied** widerspruchsfrei, und der Schluss daraus steht
seit `98d9e5ab` als Grenze in einem Gate-Vertrag. F-6 und F-7 treffen die
Pflicht-Vorprüfungen bzw. die Abgrenzungs-Begründung und sind an der Quelle
entscheidbar, nicht Urteil.

**Übergabe:** Findings gehen an den Implementer (Rückkante Review → Plan, denn
F-1/F-2 sind Plan-Defekte, keine Umsetzungsfehler); die Finding-Klassen gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist
ein Lauf-Beleg und ersetzt keine Verifikation.
