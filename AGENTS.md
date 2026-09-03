# AGENTS.md — Briefing für AI-Coding-Agenten

## 1. Was diese Datei ist

Onboarding-Briefing für jede AI-Session, die in diesem Repo Code oder
Dokumentation ändert. Sie verweist auf die kanonischen Quellen und
formuliert die Hard Rules, die der Implementation-Agent immer
einhalten muss.

Diese Datei trägt **Hard Rules und Pointer** auf die kanonischen Quellen und
**dupliziert deren Inhalt nicht** — sonst entsteht Drift (Kanon:
[`modul-09-implementierung.md` §AGENTS.md-Regeln](.harness/baseline/v5.18.0/regelwerk/modul-09-implementierung.md#agentsmd-regeln-modul-9)).

**Bei Konflikt gilt die höherrangige Quelle, und die niedriger rangierte wird
angepasst** (Source Precedence — siehe
[`harness/README.md`](harness/README.md)). Das gilt zwischen **dieser Datei**
und einer kanonischen Quelle ebenso wie zwischen **zwei kanonischen Quellen**.
**Melde den Widerspruch**, statt ihn stillschweigend nach einer Seite
aufzulösen — wer ihn nur befolgt, lässt die falsche Stelle stehen.

Strukturregeln (ID-Schemata, Verzeichniskonvention, Adaptionen ggü.
Baseline, Modus-Deklarationen pro Sub-Area, Zusatzklassen für
Sensors-Bindung) leben in
[`harness/conventions.md`](harness/conventions.md) — **vor jeder Änderung an
Code oder Dokumentation zu lesen**, nicht nur vor Doku-Änderungen.

Das Betriebsregelwerk der adoptierten Baseline ist **committet vendored**:
das nach Modulen und Grundlagen-Abschnitten aufgeteilte Regelwerk liegt
entpackt unter `.harness/baseline/<tag>/regelwerk/` (die dortige `README.md`
ist der Index), samt `.harness/baseline/<tag>/SHA256SUMS`-Integritätsmanifest —
**netzlos auf jedem Checkout präsent**, offline materialisier-/verifizierbar
per `tools/harness/fetch-baseline-cache.sh` (`--verify` offline-Integrität;
`--check-latest` = Currency- + Content-Drift-Audit ggü. Upstream, informativ/kein Gate,
[`MR-022`](harness/conventions.md#mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019); Tag aus §Baseline;
Quelle ist das derivative Release-Bundle
[`lab-regelwerk.zip`](https://github.com/pt9912/ai-harness-course/releases/download/v5.18.0/lab-regelwerk.zip);
Pfadschema/Provenance siehe
[`harness/conventions.md`](harness/conventions.md) §Adoptierte Konventions-Quellen,
[`MR-019`](harness/conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)).
Die **verkörperte Form** (dieses Briefing, die Konventionen, die
ausgefüllten Artefakte) **führt**; das Regelwerk ist die präsente,
nachschlagbare **Vertiefung** und wird **pro Entscheidung** nachgeschlagen,
deren operative Detailtiefe das Briefing nicht trägt — Trigger-Klassen,
Sub-Area-Qualifikation, Carveout-vs-Reconciliation, Modus-Diagnose. Dabei
pro Session **nur den benötigten Abschnitt** lesen, bevor der Workflow (§6)
startet — nicht das gesamte Regelwerk im Kontext halten. **Breiterer
Pflicht-Blick** bleibt bei: Bootstrap, Änderung an
[`harness/conventions.md`](harness/conventions.md) (Adaptionen `MR-<NNN>`,
Source-Precedence, ID-Schema) und dem Drift-Audit gegen die Baseline
([`modul-02-harness-bootstrap.md` §Freshness-Audit](.harness/baseline/v5.18.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2)
— darunter die **Bestands-Stichprobe, die auch bei aktuellem Pin läuft**).
Die **Skelett-Vorlagen** der Baseline liegen aus demselben self-contained Bundle
**committet vendored** unter `.harness/baseline/<tag>/templates/` (parallel zum
`.harness/baseline/<tag>/regelwerk/`-Baum, netzlos) und tragen zwei Rollen: als
**Referenz-Form**, auf die das Regelwerk als „Ziel-Form" verweist, und als **Vorlage**
beim Anlegen neuer Artefakte (ADR, Slice, Welle, …). d-checks gelebte Slice-/ADR-Struktur
folgt dabei einer **Haus-Stil-Form** — in Etappe C als baseline-konforme Form-Wahl
aufgelöst, nicht als Fork. Das Bundle ist derivativ; bei Konflikt sticht die Quelldatei
das Bundle, über ihr die kanonischen Quellen (Source Precedence). Stand/Provenance führt
[`harness/conventions.md`](harness/conventions.md) (§Adoptierte Konventions-Quellen bzw.
§Baseline).

## 2. Kanonische Quellen (Source Precedence)

In dieser Reihenfolge:

1. [`spec/lastenheft.md`](spec/lastenheft.md) — vertraglich abnahmebindend.
2. [`spec/spezifikation.md`](spec/spezifikation.md) — technisch verbindlich, fortschreibbar.
3. [`spec/architecture.md`](spec/architecture.md) — Komponenten- und Sequenzsicht.
4. [`docs/plan/adr/README.md`](docs/plan/adr/README.md) — ADR-Index.
5. [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md) — Wellen-Sequenzierung (offene Wellen derivativ).
6. [`docs/user/`](docs/user/) — Operations, Releasing.
7. [`README.md`](README.md) — Projekt-Überblick.
8. **AGENTS.md (diese Datei).**
9. [`harness/README.md`](harness/README.md) — Harness-Einstieg.

## 3. Harte Regeln

### 3.1 Docker/make-only

Implementierungssprache ist **Go**
([ADR-0001](docs/plan/adr/0001-implementierungssprache.md)). Es gilt:
**kein Host-Go und keine Host-Paketmanager** (`go`, `pip`, `npm`,
`cargo`, `apt`, `brew`, …). Alle Checks laufen über `make`; die
Go-Toolchain läuft in Docker (Multi-Stage gemäß
[ADR-0002](docs/plan/adr/0002-distribution-ghcr-image.md)). Der Host braucht `git`, GNU `make`, `bash`,
Docker und die POSIX-Standardwerkzeuge, die die Gate-Skripte rufen
(coreutils, findutils, `grep`, `awk`) — als **Klasse**, nicht als Liste.

**Auch keine Host-Skript-Interpreter** (`python`, `perl`, `ruby`, `node`, `uv`, …)
([`MR-040`](harness/conventions.md#mr-040)). Datei-Änderungen macht das Werkzeug
ohne Shell; für Messungen gilt die Rangfolge **Produkt vor `grep`/`awk` vor
allem anderen**. Eine Stage für Skripte gibt es nicht, und diese Regel verweist
auch nicht auf eine: ein Fall, der `bash` und die genannte Host-Klasse
übersteigt, ist ein **Entscheid** — keine vierte Toolchain nebenbei
([`MR-046`](harness/conventions.md#mr-046)).

**Die Host-Klasse oben ist die Zusage an den Agenten, nicht die vollständige
Liste dessen, was Gate-Skripte rufen.** Die **Netz**-Targets holen sich mehr:
[`fetch-baseline-cache.sh`](tools/harness/fetch-baseline-cache.sh) braucht
`curl` und `unzip`, [`pin-freshness.sh`](tools/harness/pin-freshness.sh)
braucht `curl`, [`nightly-state.sh`](tools/harness/nightly-state.sh) ebenso,
und [`image-scan.sh`](tools/image-scan.sh) braucht **Docker mit Netz** (Trivy
zieht seine Vuln-DB). **Alle vier** stehen bewusst außerhalb von `gates`; wer
sie fährt, fährt sie mit dieser zusätzlichen Erwartung. **Die ersten drei sind
fail-open, das vierte nicht** — ein gescheiterter CVE-Scan meldet Exit 2 und
ausdrücklich keinen grünen Befundstand ([ADR-0066](docs/plan/adr/0066-cve-scan-gegen-das-publizierte-image.md)).

**Falsch:** `go build ./…`, `go test ./…`, `pip install …`, `python3 - <<EOF`
**Richtig:** `make gates`

**Begründung:** Toolchain-Reproduzierbarkeit + Supply-Chain-Defense.

**Die Regel gilt unabhängig von ihrer Durchsetzung.** Wer sich auf den Wächter
verlässt, verlässt sich auf nichts.

**Durchsetzung, zwei unabhängige Schichten.** Die zweite ist eine
**Permission-Sperrliste** in [`.claude/settings.json`](.claude/settings.json):
sie hängt an keinem Hook, matcht aber den **ganzen** Befehl ab dem Anfang und
sieht deshalb weder Präfixe noch Sub-Shells noch zusammengesetzte Kommandos.
Ihre git-/docker-Hälfte hat keine zweite Schicht unter sich. Grenzen und
Nicht-Zusagen: [`MR-047`](harness/conventions.md#mr-047).

Die erste ist ein Tool-Call-Wächter
([`.claude/hooks/pretooluse-command-guard.sh`](.claude/hooks/pretooluse-command-guard.sh)):
er prüft die Befehlsposition jedes Segments und Sub-Shell-Strings rekursiv. Er ist
**werkzeug-lokal**, kein Repo-Gate: keine CI ruft ihn, ein Lauf ohne dieses
Werkzeug ist ungebunden. Er ist in `bash` geschrieben und liest die Hook-Eingabe
mit `awk`, läuft also in derselben Klasse, die er durchsetzt
([`MR-042`](harness/conventions.md#mr-042)); `make guard-probe` (§4) fährt ihn
gegen seine Proben.

**Er ist ein Stolperdraht, keine Sandbox.** Ungeprüft bleiben: ein
Shell-Schlüsselwort als Segment-Kopf (`if … then pip …`), ein
Wrapper außerhalb seiner Präfix-Liste (`nohup`, `timeout`), ein wort-interner
Quote-/Backslash-Splice (`p"i"p`) und escapte Quotes in der Verschachtelung
erreichen ein gelistetes Werkzeug, ohne dass er es als Kopf sieht. Die Regel
gilt trotzdem — sie hängt nicht an ihm. Tabelle in
[`MR-042`](harness/conventions.md#mr-042).

### 3.2 Suppression-Verbot

Inline-Suppressions sind verboten: `//nolint`-Direktiven im Code
brechen das künftige Suppression-Gate. Ausnahmen leben zentral in
`.golangci.yml` (exclude-rules) mit Begründung.

**Teilweise durchgesetzt, und die Grenze:** `nolintlint` im Profil
meldet eine Direktive ohne benannten Linter, ohne Begründung oder ohne Wirkung —
sie wird damit sichtbar und zurechenbar. Eine **wohlgeformte** `//nolint`
unterdrückt einen echten Verstoß weiterhin, und `make lint` bleibt grün: der
Linter prüft die **Form** der Direktive, nicht ihre Berechtigung. Verboten
bleibt sie durch diese Regel, nicht durch das Gate. *(Auflösungs-Trigger:
permanent — die Berechtigungsfrage ist ein Urteil.)*

### 3.3 git mv + Inhaltsänderung = zwei Commits

Wenn eine Datei verschoben **und** der Inhalt umgeschrieben wird:

1. `git mv source target` → eigener Commit (reiner Move, Git erkennt R-Rename).
2. Inhalt umschreiben → zweiter Commit.

**Begründung:** Sonst fällt die Rename-Detection unter die
50%-Similarity-Schwelle und `git log --follow` wird unzuverlässig.

**Teilweise durchgesetzt:** `make planning-check` hält die **Kopplung**
Lifecycle-Verzeichnis ↔ Roadmap-Ruhe-Marker in beide Richtungen. Die
**Commit-Zerlegung** selbst — Move und Inhaltsänderung in einem Commit — sieht
kein Gate. *(Auflösungs-Trigger: permanent.)*

**Ausnahme Slice-Lifecycle-Move (`in-progress/` → `done/`):** Der
`git mv`-Commit trägt hier **zusätzlich** den Roadmap-Flip §Offene Wellen
(zurück auf den Ruhe-Marker „Nichts in Arbeit", sofern kein Slice mehr
beansprucht ist — der Marker steht dann **zusätzlich** zur Zeiger-Liste;
das flache Wellendokument darf offen bleiben, nur die Bijektion muss
stimmen: `wave-drift` misst unter `mode: many` Zeiger ⟺ Dateien in beide
Richtungen) und alle Pfad-Verweise auf den Slice
(Roadmap, §4, `harness/README.md` §Sensors) von `in-progress/` nach
`done/`. Sonst ist der Commit gate-rot: `make planning-check` koppelt
in-progress-Stand und Roadmap atomar, und die alten Verweise laufen ins
Leere (`target-missing`). Nur der **Slice-Body** (DoD-Haken + Closure-Notiz;
historische Slices auch die Status-Zeile) bleibt Commit 2 — die Slice-Datei selbst ist im Move-Commit
unverändert, also hält die Rename-Detection. Kanonisch:
[`MR-013`](harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise).

**Ausnahme Beanspruchung (`open/` → `in-progress/`):** dieselbe Kopplung mit
umgekehrtem Vorzeichen — der Ruhe-Marker **verlässt** §Offene Wellen, und die
Pfad-Verweise auf den Slice wandern von `open/` nach `in-progress/`. Auch hier
ist ein byte-reiner Move-Commit gate-rot: `planning-drift` meldet den Slice in
`in-progress/` gegen die Roadmap, die den Marker noch trägt. Kanonisch:
[`MR-013`](harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise).

**Ausnahme MR-/Wellen-Lifecycle-Move** (`conventions/` → `conventions/done/`,
flaches Wellendokument → `done/`): hier trägt der Move-Commit die
**Link-Tiefen-Fixes der bewegten Datei selbst** mit — ein reiner Move wäre
`doc-check`-rot, weil die relativen Verweise vom neuen Ort nicht mehr
auflösen. Alles Übrige bleibt Commit 2; sinkt der Rename-Score dadurch
Richtung 50 %, deklariert die Commit-Botschaft den Move ausdrücklich als
`git mv`. Kanonisch:
[`MR-013`](harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise).

**Ausnahme Wellen-Archiv-Stub-Move** (`tools/archive-wave`, Modul 6
Schritt 4): hier gibt es **keine** Zwei-Commit-Zerlegung, weil es keine
Phase gibt, in der die bewegte Datei ihren Inhalt unverändert behält — der
Stub *ersetzt* den Volltext im selben Akt, der ihn verschiebt. Ein
Wellen-Archivierungs-Commit bleibt deshalb bewusst **ein** Commit, in der
Botschaft ausdrücklich als solcher deklariert; git zeigt reine `D`/`A`-Paare,
keine Renames. Kanonisch:
[`MR-059`](harness/conventions.md#mr-059--wellen-archiv-stub-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013).

### 3.4 Architektur sprach-/meilensteinfrei; Spec-Straten nie abwärts

[`spec/architecture.md`](spec/architecture.md) benennt Schichten und
Rollen statt Technologie — keine Sprach-/Modul-Pfade. Kein
Spec-Stratum (auch [`spec/spezifikation.md`](spec/spezifikation.md))
referenziert ADRs, Wellen, Slices, Commit-Hashes oder Closure-Daten.
Die sprachkonkrete Übersetzung (Modul-Pfade, Import-Regeln) und die
Begründungen leben in den ADRs, deren `Schärft:`-Feld aufwärts zeigt;
die zeitliche Schicht lebt in `docs/plan/planning/`.

**Die Abwärts-Sperre nennt fünf Kategorien; gedeckt sind vier.** `make doc-check`
(Modul `matrix`) hält **ADRs** über die Link-Prüfung, **Slices**, **Wellen** und
**Commit-Hashes** zusätzlich als Token im Körper
([`DC-FA-MTX-003`](spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)).
Erkennungs-Formen und ihre benannten Grenzen stehen bei der Regel, die sie
trägt — im `matrix`-Block der [`.d-check.yml`](.d-check.yml).

**Ungedeckt bleibt das Closure-Datum:** es ist von einem legitimen Datum nicht
unterscheidbar, und die Spec-Straten führen eigene Historie-Zeilen voller Daten.
*(Auflösungs-Trigger: keiner — die Kategorie bleibt Urteil.)*

**Die Sprachfreiheit der Sicht zerfällt in zwei ungleiche Hälften.** Ob eine
Zeile **Rollen statt Technologie** benennt, ist ein Urteil — *permanent*. Ob sie
einen **Modul-Pfad** trägt, ist ein detektierbarer Zustand: `matrix` führt die
Sicht als eigene Klasse und die Code-Modul-Pfade als Token-Ziel daneben.
Gegenstand sind **Modul**-Pfade, nicht Pfade überhaupt — Dokument- und
Skript-Pfade bleiben erlaubt, denn weder ein Dokument noch ein Gate-Skript ist
ein Modul.

**Das Pfad-Verbot ist eine Verschärfung gegenüber der Baseline** — sie erlaubt
der Sicht Modul-Pfade ausdrücklich. Geführt als
[`MR-033`](harness/conventions.md#mr-033).

### 3.5 ADRs sind nach `Accepted` immutable

Eine ADR mit Status `Accepted` wird nicht inhaltlich überschrieben.
Korrekturen entstehen als neue ADR mit `Supersedes ADR-NNNN` (vierstellig).
Maschinell erzwungen über `make adr-check` (`pre-commit`-Hook + PR-/Push-CI;
erlaubt bleiben `## Geschichte`-Anhänge + der `**Status:**`-Übergang;
[ADR-0016](docs/plan/adr/0016-adr-immutable-gate.md)).

### 3.6 Gates dürfen nicht ohne ADR gelockert werden

Jede Schwellen-Senkung (Coverage, Linter-Strenge, Prüfregel) ist ein
ADR, kein PR-Kommentar.

**Kein Gate prüft das** — die Regel gilt einem **Akt**, nicht einem ruhenden
Zustand: ob eine gesenkte Schwelle eine ADR *hat*, steht in keiner Datei, die
ein Sensor gegen die Senkung halten könnte. *(Auflösungs-Trigger: permanent.)*

### 3.7 Kommentare tragen eine der fünf Klassen

Gilt für Code, Konfiguration und Skripte — und, mit **eigener** Form, für
Zustandsfelder (unten). **Ein Kommentar beschreibt, was da ist**
(Baseline-Merksatz). Er
beantwortet in Code, Konfiguration oder Skript, was der Code nicht
beantworten kann — **Zusage · Kopplung · Abgrenzung · Rang-Zeiger ·
Grenze** ([Baseline §Was ein Kommentar trägt](.harness/baseline/v5.18.0/regelwerk/grundlagen-harness-dateien.md#was-ein-kommentar-trägt--code-konfiguration-skripte)).
Keine Review-Historie und keine Review-Befund-Marker, keine Deliberation
über Verworfenes, keine Herkunfts-Prosa, keine Slice-Nummern und keine
Mess-Labels; Herkunft nur als **ein** auflösbares Feld nach dem
Baseline-Schema (`DC-*` — die Baseline-Form `LH-*` —, `ADR-*`, `MR-*`,
`seit welle-<NN>`). Der Reviewer-Skill trägt den HIGH-Anker dazu.

**Zustandsfelder** (seit dem v5.9.0-Bump) sind Zustands-Artefakte wie der
Kommentar, nur im Rumpf — sie tragen **nicht** dessen fünf Klassen, sondern
eine **eigene Form**; übertragen sind die **zwei Tests**: Adressat ist, wer
den Zustand liest, um zu handeln, und die Zeitform ist der Indikativ über das,
was ist. Ein Feld, das einen Zustand trägt — etwa eine `Stand`- oder
`Status`-Zelle in Roadmap, Beobachtungs-Register oder Meilenstein-Tabelle —,
nennt **den Zustand und den Beleg als auflösbaren Anker**, nicht die Chronik,
wie es dazu kam; das
Drift-Log der Roadmap trägt **nur Umplanungen**, keine Schließungen und keine
erreichten Meilensteine (die stehen im Closure-Log bzw. in der Status-Spalte).
Ein lebendes Register trägt **keine** Kopfzeile `Status: Aktiv. Letzte
Änderung: <Datum>` — sein Zustand ist sein Inhalt, sein Änderungsdatum hält
`git`; ein Datum, das ein **benannter Trigger** pflegt, ist davon ausgenommen
(der Frische-Marker der Architektur-Sicht). **Verhältnis zu §3.5:** das
`**Status:**`-Feld einer ADR ist ein Zustandsfeld wie jedes andere — `adr-check`
nimmt die Kopf-Status-Zeile ausdrücklich **aus** dem Kern-Vergleich und lässt
den Übergang zu. Es darf also korrigiert werden, solange der Wert die erlaubte
Form behält; §3.5 schützt den **Kern**, nicht dieses Feld. **Benannte
Bestands-Ausnahme:** die historischen `**Status:**`-Felder der `done/`-Slices
bleiben, wie sie sind — sie sind eingefrorene Lauf-Belege, ihr
Lifecycle-Zustand ist ohnehin das Verzeichnis, und das Feld hat dort keine
Funktion (§5). Gemeldet wird von ihnen nur, was dem Verzeichnis
**widerspricht**. Kanon:
[Baseline §Was ein Kommentar trägt](.harness/baseline/v5.18.0/regelwerk/grundlagen-harness-dateien.md#was-ein-kommentar-trägt--code-konfiguration-skripte).

**Kein Gate prüft das** — weder die fünf Klassen noch die Zustandsfeld-Form;
die Prüfung ist ein Urteil, kein `grep`. Der Reviewer-Skill trägt dazu **zwei**
HIGH-Anker.

**Bestandsgrenze:** vor der Einführung bzw. vor dieser Schärfung
geschriebene Kommentare (Test-Kommentare; ältere Config-Kommentare mit
Slice-Nummer) sind grandfathered — geräumt wird beim nächsten Anfassen
der Zeile; Neuzugänge fallen überall unter den Anker. **Für Zustandsfelder
gibt es keine Bestandsgrenze:** der vorhandene Bestand wird mit dem
v5.9.0-Bump umgestellt, nicht grandfathered. *(Hard Rule seit
dem v5.6.0-Bump, geschärft mit dem v5.7.0-Bump auf die
Baseline-Feld-Formen und mit dem v5.9.0-Bump auf Zustandsfelder;
Auflösungs-Trigger: permanent.)*

### 3.8 Ein Modul verspricht nur über das, was es scannt

Jedes Modul gibt seine Zusagen über seine **Scan-Menge**. Daneben liest es
Eingaben, die es nie scannt: Zieldateien außerhalb der Scan-Wurzeln, selbst
benannte Verzeichnisse eines Post-Passes, git-Revisionen. Dort gilt **keine**
dieser Zusagen — und die Folge kann **still** sein: ein verdecktes Heading
macht einen Anker unauflösbar, die Prüfung entfällt kommentarlos. Wer ein Modul
anlegt oder ändert, beantwortet deshalb: **welche Eingaben liest es, die es
nicht scannt — und gilt dort dieselbe Zusage?** Wo sie nicht gilt, gehört die
Grenze in die Anforderung.

**Begründung:** Eine Liste der Achsen wäre selbst unvollständig — deshalb steht
hier die Frage und keine Liste. Kein Gate fängt das; der Reviewer-Skill trägt
den MEDIUM-Anker dazu. *(Hard Rule aus dem
Steering Loop, [Beobachtungs-Register](docs/plan/planning/observations.md)
`BEO-004`, seit welle-73; Auflösungs-Trigger: permanent.)*

### 3.9 GitHub-Action-Referenzen sind SHA-gepinnt

Jeder `uses:`-Eintrag in [`.github/workflows/`](.github/workflows) nennt einen
**vollen Commit-SHA** mit Tag-Kommentar dahinter, nie einen beweglichen Tag.
Das gilt für jeden Workflow gleich und für jeden Neuzugang.

**Eine Ausnahme, und sie ist keine Lockerung** ([ADR-0068](docs/plan/adr/0068-lokale-workflow-referenzen-ohne-pin.md)):
eine **lokale** Workflow-Referenz (`uses: ./.github/workflows/x.yml`) kann keinen
SHA tragen und **braucht keinen**. Sie löst auf denselben Commit auf wie der
aufrufende Workflow und ist damit **stärker** gebunden als ein SHA-Pin — sie kann
per Konstruktion nicht driften. **An die Stelle des Pin-Checks tritt die Frage,
die hier trägt: existiert das Ziel?** Ein vertippter Verweis fiele sonst erst zur
Laufzeit auf; `make workflow-pins` meldet ihn als `uses-local-missing`. Sie gilt
**nur** dem `./`-Präfix; eine Referenz in ein fremdes Repository
(`owner/repo/.github/workflows/x.yml@ref`) fällt unter die Regel wie jede
Action. Die Zahl der lokalen Referenzen steht in der Erfolgsmeldung, statt
stillschweigend übergangen zu werden.

**Die Existenz ist nicht die einzige Frage — das stand hier zu weit**
([ADR-0071](docs/plan/adr/0071-lokale-workflow-referenz-rechte-pruefung.md)).
Ein aufgerufener Workflow bekommt nur die Rechte, die der aufrufende **Job**
selbst führt; verlangt er mehr, lehnt GitHub den **ganzen Lauf vor dem ersten
Job** ab (`startup_failure`, kein Log) — gemessen am Tag-Push von `v0.66.0`,
während dieses Gate grün meldete. Geprüft wird deshalb auch die
**Rechte-Anforderung des Ziels**: ein Job ohne eigenes `permissions:`, dessen
Ziel Rechte verlangt (`uses-local-perms-undeclared`), und ein Aufrufer, der
einen geforderten Scope zu niedrig führt (`uses-local-perms-narrow`). Was der
Wächter nicht sicher liest, meldet er (`uses-local-perms-unreadable`), statt es
zu übergehen. **Er deckt damit eine Fehlerklasse, nicht die Lauffähigkeit** —
und die Zerlegung ist eine Näherung über die YAML-Block-Form, keine
Parser-Zusage.

**Begründung:** Supply-Chain-Härtung — ein Tag lässt sich umhängen, ein SHA
nicht; dieselbe Härtung wie der Docker/make-only-Pfad in §3.1.

**Durchgesetzt:** `make workflow-pins` in `make gates` — seit
[ADR-0072](docs/plan/adr/0072-workflows-modul.md) **via Modul `workflows`**
(Dogfooding über das eigene Image; das frühere Skript ist darin aufgegangen).
Er trägt die **Form** —
voller SHA plus Tag-Kommentar —, nicht die **Gültigkeit**: ob der SHA existiert
und den Commit bezeichnet, den der Tag-Kommentar behauptet, prüft er nicht.
*(Auflösungs-Trigger: permanent — die Gültigkeitsfrage ist Netz und gehört zur
Freshness-Familie.)*

## 4. Quality Gates

Nur hier gelistete Targets existieren im Makefile. Halluzinierte
Gates sind die häufigste Form von Harness-Lüge.

| Target                       | Zweck                                                                                                                                                                                                                                              |
| ---------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `make lint`                  | golangci-lint mit dem Projekt-Profil (§3.2)                                                                                                                                                                                                        |
| `make test`                  | `go test ./...` — Akzeptanzkriterien der `DC-FA-*`; [`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Netzlos-Modullisten-Integrität der `.d-check.yml` (Go-Test, [ADR-0032](docs/plan/adr/0032-gate-consistency-tombstone.md)); **Wohlgeformtheit der Durchsetzungs-Konfiguration** (Go-Test, [`MR-048`](harness/conventions.md#mr-048)): `.claude/settings.json` ist gültiges JSON, jeder verdrahtete Hook-Pfad existiert, der Befehls-Wächter hängt am `Bash`-Werkzeug, und die Permission-Sperrliste deckt jeden Namen der Wächter-Sperrliste als ganze Befehlsklasse. **Er sagt nicht, dass die Durchsetzung läuft**, und er fordert die Dateien nicht ein — fehlen sie, überspringt er                                                                                                                                                                                                 |
| `make arch-check`            | Import-Regeln des Hexagon-Schnitts + Kern-Paket-Richtung **via digest-gepinntes a-check-Image** (Schwester-Tool, `a-check.mk` + `.a-check.yml`, netzlos/read-only) ([ADR-0005](docs/plan/adr/0005-modul-layout-hexagon-ordner.md), [ADR-0012](docs/plan/adr/0012-kern-paketschnitt-model-rules-app.md), [ADR-0029](docs/plan/adr/0029-arch-check-via-a-check.md) löst die Skript-/Stage-Mechanik ab)                                                      |
| `make coverage-gate`         | Coverage-Schwelle über `./internal/...` (Kalibrierungs-Bindung, siehe [`harness/README.md`](harness/README.md) §Sensors)                                                                                                                           |
| `make gate-consistency`      | Meta-Gate: Deklarations-Konsistenz Doku↔Makefile via Modul `targets` (Image, dogfood; [ADR-0031](docs/plan/adr/0031-targets-deklarations-konsistenz-modul.md); [ADR-0032](docs/plan/adr/0032-gate-consistency-tombstone.md) löst das Rest-Skript voll ab). Die [`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Modullisten-Integrität prüft jetzt ein getippter Go-Test in `make test`                                                                           |
| `make planning-check`        | Meta-Gate **via Modul `planning`** (Image, dogfood): Roadmap §Offene Wellen (Ruhe-Marker) ↔ `in-progress/slice-*` (`planning-drift`, hermetisch — kein git, in `gates`) ([ADR-0028](docs/plan/adr/0028-planning-lifecycle-modul.md) löst die frühere Skript-Mechanik ab; [`DC-FA-PLAN-001`](spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)) |
| `make doc-check`             | Doku-Links, Anker, Kennungs-Linkpflicht, Referenzmatrix, Inline-Code-Pfade, Abschnitts-Invarianten (Modul `structure`) + Kennungen in Diagramm-Fences (Modul `diagrams`) + **wortgleiche Zitate gegen ihre Quelle** (Modul `citations`) via `d-check` selbst (Dogfooding; netzlos — zugleich [`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Messmethode). **`citations` ist fail-closed und läuft im inneren Loop:** eine strukturell unbrauchbare Direktive nimmt den ganzen Lauf mit, auch den `pre-commit`-Hook. Das ist tragbar, seit die Direktive nur außerhalb von Inline-Code zählt ([ADR-0060](docs/plan/adr/0060-citations-marker-scan-geteilte-prosa-antwort.md)) — der fail-closed-Rest ist ein Autoren-Fehler an einer **freien** Direktive. **Die planmäßige Rot-Quelle ist ein anderer:** der Baseline-Bump verschiebt die Zeilenspannen, und das Neu-Ankern liegt in der Bump-Prozedur ([`MR-051`](harness/conventions.md#mr-051)). Ein Merge trägt eine fremde Direktive am Hook vorbei, und `d-check:ignore` greift hier **nicht** — grob wirken nur `scan.ignore` und `citations.scope` |
| `make workflow-pins`         | Deklarations-Konsistenz der `uses:`-Referenzen **via Modul `workflows`** (Image, dogfood; [`DC-FA-WF-001`](spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in), [ADR-0072](docs/plan/adr/0072-workflows-modul.md) löst die Skript-Mechanik von [ADR-0071](docs/plan/adr/0071-lokale-workflow-referenz-rechte-pruefung.md) ab). **Fremde** Referenz: voller 40-stelliger Commit-SHA (`uses-pin-missing`) plus Tag-Kommentar (`uses-pin-untagged`) — geprüft wird die **Form**, nicht die Gültigkeit des SHA (das wäre Netz). **Lokale** Referenz (`./…`): das Ziel existiert (`uses-local-missing`) **und** der aufrufende **Job** führt die Rechte, die es verlangt (`uses-local-perms-undeclared`, `uses-local-perms-narrow`; `none` < `read` < `write`, ein nicht genannter Scope ist `none`). Unlesbares YAML meldet (`workflow-unparsable`), statt übersprungen zu werden. Die Referenzen kommen aus dem **YAML-Baum**, nicht aus einer Textsuche; das Verzeichnis ist **konfigurierbar** (`workflows.dir`), nicht verdrahtet. Liest `.yml` **und** `.yaml`; **fail-closed auch bei leerer Prüfmenge**. **Zwei benannte Grenzen:** es liest die **Ziele** lokaler Referenzen, die es nicht scannt (dieselbe Parse-Zusage gilt dort), und es deckt **eine** Deklarations-Klasse, nicht die Lauffähigkeit. Netzlos, hermetisch, **Bestandteil von `gates`** (§3.9) |
| `make review-coverage`       | Review-Report-Deckung **via Modul `reviews`** (Image, dogfood; [`DC-FA-RVW-001`](spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in), [ADR-0081](docs/plan/adr/0081-reviews-modul.md)). Ein `done/`-Slice mit Review-Zusage (ein DoD-Haken, dessen Zeile „unabhängiger Review" nennt, jede Bullet-Form, Haken-Zustand egal) braucht mindestens einen Report unter `reviews.reviews-dir` mit derselben `slice-<NNN>`-Kennung im Dateinamen (`review-missing`) — Substring-Match, 1:N zulässig. Beide Verzeichnisse **nicht rekursiv** gescannt: ein archivierter Slice-Stub trägt keine DoD mehr und fällt natürlich aus der Kandidatenmenge. **Fail-closed bei leerer Kandidatenmenge oder unlesbarem `reviews-dir`**, **nicht** bei null gefundenen Zusagen unter vorhandenen Kandidaten (ein junger Bestand ohne jede Zusage ist legitim). **Bestands-Ausnahme, feste Dateiliste** (fünf Funde beim Scharfschalten, davon zwei mit geschlossenem Haken — für `structure`s Haken-Wächter unsichtbar). Netzlos, hermetisch, **bewusst noch nicht Teil von `gates`/`ci`** — eine neue Modul-Klasse startet als eigenständiger Fokus-Lauf, dieselbe Vorsicht wie bei `trace-check` |
| `make baseline-verify`       | Vendorte Baseline gegen `SHA256SUMS`: `sha256sum -c` (geänderte/gelöschte Datei) **plus Manifest-Deckung** (zusätzlich eingelegte Datei — eine untermengige, in sich konsistente `SHA256SUMS` passierte sonst grün) **plus Auflösung der Aliase** unter `.claude/rules/`: ein Symlink bindet denselben Pin, wird aber von keinem Modul gescannt und steht in keiner Manifest-Zeile — beim Bump bräche er still ([`MR-055`](harness/conventions.md#mr-055)). Rekursiv und dotfile-bewusst; **benannte Grenzen:** geprüft wird die **Auflösung**, nicht das Ziel, und ein **fehlendes** `.claude/rules/` ist von "hier gibt es keine Aliase" nicht unterscheidbar. Netzlos, fail-closed, **Bestandteil von `gates`** ([`MR-011`](harness/conventions.md#mr-011)-Kette)                                                    |
| `make baseline-probe`        | Fährt die **Alias-Auflösung** von `baseline-verify` gegen ihre Proben: neun Fälle mit Erwartung und Ergebnis, Fehlschlag-Zähler am Ende — gesunder Alias, toter Alias flach/im Unterbaum/als Punkt-Name, Alias auf ein Verzeichnis, Symlink-Zyklus, echte Datei, leeres und fehlendes Verzeichnis. Netzlos, ohne Wirkung auf das Repo, **bewusst nicht in `gates`** — dieselbe Begründung wie bei `make guard-probe`: eine Zusage ohne wiederholbare Probe ist eine Erinnerung ([`MR-055`](harness/conventions.md#mr-055)) |
| `make gates`                 | alle inneren Gates (mandatory vor Handoff)                                                                                                                                                                                                         |
| `make ci`                    | CI-äquivalenter Lauf: gates + image-test (fährt die Release-Pipeline)                                                                                                                                                                              |
| `make trace-check`           | Traceability-Gate **via Modul `commits`** (Image, dogfood): DC-/ADR-/MR-/slice-ID in Commit-Messages (`commit-untraceable`; `RANGE=`-Range für CI, `MSGFILE=` für den `commit-msg`-Hook via stdin; bewusst **nicht** Teil von `gates`/`ci`) ([ADR-0027](docs/plan/adr/0027-commits-traceability-modul.md) löst die Skript-Mechanik von [ADR-0013](docs/plan/adr/0013-pr-ci-und-traceability-gate.md) ab; [`DC-FA-COMMITS-001`](spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in)) |
| `make adr-check`             | ADR-Immutable-Gate **via Modul `vcs`** (Image, dogfood): `Accepted`-ADRs nicht inhaltlich ändern (`RANGE=`/`STAGED=`-Modi; `pre-commit`-Hook + PR-CI; bewusst **nicht** Teil von `gates`/`ci`) ([ADR-0024](docs/plan/adr/0024-vcs-immutable-gate.md) löst die Skript-Mechanik von [ADR-0016](docs/plan/adr/0016-adr-immutable-gate.md) ab, [ADR-0025](docs/plan/adr/0025-codepaths-ignore-refs.md) entfernt das Alt-Skript)                                                            |
| `make freshness-go` · `make freshness-golangci` · `make freshness-semgrep` · `make freshness-a-check` | **Versions-Achsen:** melden, wenn upstream einen neueren Release führt als der Pin — Go-Toolchain, `golangci-lint`, `semgrep`, das Schwester-Tool `a-check`. Zwei Quellen-Formen: `go.dev/VERSION?m=text` (Plaintext — golang/go führt **keine** Release-Objekte) und der `releases/latest`-Redirect. Ein führendes `go`/`v` wird **beidseitig** gestrippt; unsere Pins führen das `v` uneinheitlich. **Netz**, fail-open (Ausfall ⇒ `SKIP`, Exit 0), bewusst **nicht** in `gates`/`ci`; gerufen vom Nachtlauf ([`upstream-drift.yml`](.github/workflows/upstream-drift.yml)). Der Vergleicher ist über `--compare` **netzlos** prüfbar. Sie melden nur — die Hebung bleibt ein bewusster Akt, und beim Go-Bump zieht das `golangci`-Pendant mit |
| `make checkout-pin-freshness` · `make login-pin-freshness` · `make hubdesc-pin-freshness` | Dieselbe Versions-Frage für die **drei** **Action-Pins** der Workflows. Der Pin ist ein SHA, aber der **Tag-Kommentar daneben** trägt den Release-Tag — genau die Größe, die der `releases/latest`-Redirect liefert; es braucht dafür keine eigene Quellen-Form. Ergänzt `make workflow-pins`, das die **Form** hält (§3.9): ein SHA-Pin schließt das Umhängen eines Tags aus und macht zugleich blind für dessen Behebung. **Netz**, fail-open, **nicht** in `gates`; Nachtlauf. **Grenze:** verglichen wird der **Kommentar**, nicht der SHA — ob der SHA den Commit bezeichnet, den der Kommentar behauptet, bleibt die benannte Grenze aus §3.9 |
| `make runtime-base-digest` · `make go-base-digest` · `make lint-base-digest` · `make semgrep-digest` · `make a-check-digest` | **Digest-Achsen** — andere Frage als die Versions-Achsen: nicht „neuerer Tag", sondern „**neuer Bau desselben Tags**". Sie gelten allen fünf digest-gepinnten Fremd-Images. Eine Versions-Achse deckt das **nicht** ab: `freshness-go` meldet `ok`, während derselbe `golang`-Tag längst neu gebaut ist. Für das Runtime-Image sind sie zusätzlich die **einzige** Handhabe, weil sein Tag `nonroot` keine Version führt. Quelle ist `docker buildx imagetools`; **Netz**, fail-open (Ausfall ⇒ `SKIP`), bewusst **nicht** in `gates`, gerufen vom Nachtlauf. Melden `ABWEICHEND`, nicht `VERALTET` — Digests haben keine Ordnung. Melden nur — die Hebung bleibt ein bewusster Akt ([ADR-0011](docs/plan/adr/0011-digest-pins-build-gate-images.md)) |
| `make baseline-freshness`    | Upstream-Audit des Baseline-Pins: neuerer Release-Tag (Currency, gelesen aus der Release-**Liste** — nicht `releases/latest`, das Prereleases überspringt) + Content-Drift am gepinnten Tag. **Netz**, fail-open (Ausfall ⇒ `SKIP` je Teil), bewusst **nicht** in `gates`/`ci` — der netzlose innere Lauf ist eine Eigenschaft dieses Repos, keine Zusage des Produkts; gerufen vom Nachtlauf ([`upstream-drift.yml`](.github/workflows/upstream-drift.yml)). Meldet nur — die Hebung bleibt ein bewusster Akt |
| `make nightly-state`        | **Lese-Schritt, kein Gate:** liest den Ausgang des jüngsten Laufs **beider** Nachtläufe ([`upstream-drift.yml`](.github/workflows/upstream-drift.yml) und [`image-scan.yml`](.github/workflows/image-scan.yml)) über die GitHub-API und sagt, ob er gelesen werden muss. Er ersetzt den Nachtlauf nicht — er macht seinen Stand an dem Moment verfügbar, an dem ohnehin jemand hinsieht: als **dritte Vorprüfung** der Slice-Planung ([`MR-053`](harness/conventions.md#mr-053), §5). **Netz** (nur `curl`, kein `gh` — die Netz-Targets tragen diese Erwartung ohnehin), fail-open, **immer Exit 0**: der Ausgang steht in der **Ausgabe**, damit kein Exit-Code ihn verdecken kann. Bewusst **nicht** in `gates`; netzlos prüfbar über `--parse <datei>` und `--selftest` (sechs Proben). **Vier benannte Grenzen:** greift der Schritt nur bei einer Slice-Planung, liest in einer **Pause** niemand; er liest den **jüngsten** Lauf, nicht sein **Alter** — ein abgeschalteter Nachtlauf meldete weiter `gruen`; bei privatem Repo oder umbenanntem Workflow ist der `SKIP` von einer Netzstörung ununterscheidbar; und der Repo-Slug ist ein **Default** (`NIGHTLY_REPO`, `NIGHTLY_WORKFLOWS` überschreiben ihn (die Liste, seit dem zweiten Nachtlauf — der Singular ist tot)) — in einem Fork meldete es sonst den Nachtlauf des Originals |
| `make hooks`                 | git-Hooks installieren (`core.hooksPath` → `.githooks`; aktiviert `commit-msg`-Traceability + `pre-commit` mit ADR-Immutable via Modul `vcs` **und** dem vollen `doc-check` als Doku-GUARD — seit welle-79: ein roter Gate-Exit kann keine Shell-Kette mehr passieren) **plus** den Slice-Closure-Übergangs-Wächter: ein gestagter Rename/Add nach `docs/plan/planning/done/slice-*.md` (nicht rekursiv — ein archivierter Stub eine Ebene tiefer zählt nicht) löst zusätzlich `make verify-closure-notes` aus, die Vorbedingungen hängen damit am **Übergang**, nicht nur an gelegentlichen `fullbuild`-Läufen; dieselbe Bindung läuft in der PR-/Push-CI über die Commit-Range, weil `--no-verify` nur den lokalen Hook umgeht ([ADR-0013](docs/plan/adr/0013-pr-ci-und-traceability-gate.md), [ADR-0016](docs/plan/adr/0016-adr-immutable-gate.md), [ADR-0024](docs/plan/adr/0024-vcs-immutable-gate.md); Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine)    |
| `make completeness-check`    | Requirements-Completeness-Gate **via in-Produkt-Flag** `--trace --require-complete` (≥1 Waise ⇒ Exit 1, mit `WAISE`-Zeilen + Anzahl); **Closure-Bindepunkt** (in `make fullbuild`, **nicht** `gates`/`ci`) ([ADR-0026](docs/plan/adr/0026-completeness-in-product-gate.md) löst die Skript-Mechanik von [ADR-0017](docs/plan/adr/0017-requirements-completeness-gate.md) ab; [`DC-FA-CLI-011`](spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)) |   <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:142-142 -->
| `make verify-closure-notes`  | Struktur des `done/`-Bestands: die Closure-Notizen **via Modul `planning`** (Abschnitt vorhanden, Substanz außerhalb Code, keine deklarierte Floskel, opt-in kein Vorlagen-Platzhalter) **und** Abschnitts-Invarianten **via Modul `structure`** (`section-*`-Codes), darunter die **urteilsfreie Hälfte der Drei-Ausgänge-Regel** des Baseline-Regelwerks (`modul-05`: *„Ein Slice geht nicht nach `done/`, während ein Risiko ohne Ausgang dasteht"*) in **zwei** Regeln: der **vergessene** Ausgang als `forbid-pattern` über jeden H1-Abschnitt, der **erfundene** als Komplement-`forbid-pattern` über §5 — der Wortschatz ist geschlossen und umfasst alle drei Kanon-Ausgänge ([`MR-049`](harness/conventions.md#mr-049)). **Sie sieht ein Risiko ganz ohne Ausgangs-Marker nicht**, und der Altbestand vor slice-140 ist ausgenommen. **Ob** ein eingetragener Ausgang inhaltlich trägt, bleibt Urteil. **Vier Bereiche sieht sie nicht** — Überschriften-Zeile, Fenced Blocks, Inline-Code (die drei gewollt, sonst meldete ein Slice über Platzhalter seine eigene Doku) und, als ausgewiesener Preis, Prosa innerhalb einer **absatzweiten** Inline-Code-Spanne ([ADR-0059](docs/plan/adr/0059-closure-waechter-weicht-structure-regel.md)) — beide über dasselbe `--config`-Profil. **Seit slice-172 hält sie zusätzlich den DoD-Haken:** `max-open-tasks: 0` über den `## N. Definition of Done`-Abschnitt meldet jeden offenen Haken eines `done/`-Slice (`section-tasks-open`, je Haken auf **seiner** Zeile, mit verfasstem Reparatur-Hinweis), mit Bestands-Ausnahme bis `slice-170` ([`MR-056`](harness/conventions.md#mr-056)). Der Kanon macht daraus eine Bedingung des **Übergangs**; diese Regel prüft den **Zustand am Ruheort**. **Drei Grenzen:** ein Haken ist eine Selbstauskunft; ein **vergessener** Schluss-Fence macht die Bedingung blind (deshalb fährt dasselbe Profil `spans`); und ein Haken **innerhalb** eines wohlgeformten Fenced-Blocks ist unsichtbar, ohne dass irgendetwas meldet. **Ein drittes Modul macht den Bindepunkt selbstgenügsam:** `spans` meldet die Fence-Artefakte, die die Bereinigung der beiden anderen **still** machen — ein vergessener Schluss-Fence verschluckt alles dahinter, und die Zusage darüber wird wahr, ohne etwas gesehen zu haben. **Es findet dabei nichts, was `gates` nicht schon fände** (gemessen: `make doc-check` meldet denselben `fence-unclosed` beim Commit, und die Scan-Menge des Bindepunkts ist eine Teilmenge). Der Gewinn ist die **Unabhängigkeit** von einem fremden Profil, nicht neue Deckung ([ADR-0077](docs/plan/adr/0077-spans-am-bindepunkt-die-begruendung-traegt-anders.md), supersedes [ADR-0076](docs/plan/adr/0076-spans-am-closure-bindepunkt.md)). **Drei** Grund-Codes können den Bindepunkt rot machen: `fence-unclosed`, `span-unclosed` und `span-nested-link`. **Es deckt nicht** die vierte Grenze oben — ein **wohlgeformter** Span, der Prosa umschließt, ist kein Defekt, sondern Code. Fährt ein **eigenes** Prüf-Profil über `--config` ([`DC-FA-CLI-012`](spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben)); **Closure-Bindepunkt** (in `make fullbuild`, bewusst **nicht** `gates`/`ci`) ([ADR-0048](docs/plan/adr/0048-closure-note-struktur-im-planning-modul.md), [`DC-FA-PLAN-001`](spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in), [`DC-FA-SPAN-001`](spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)) |
| `make fullbuild`             | volle Closure: gates + image-test + bench + completeness-check + verify-closure-notes, schließt mit dem Image-Hash                                                                                                                                                                             |
| `make image-test`            | [`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image)-Akzeptanzkriterien gegen das lokale Image (nativ vs. Container)                                                                                                                |
| `make bench`                 | [`DC-QA-01`](spec/lastenheft.md#dc-qa-01--performance)-Benchmark gegen generiertes Fixture (Median aus 3 Läufen, kein Gate in `gates`)                                                                                                             |
| `make trace`                 | Requirements Traceability Matrix via `d-check` selbst auf stdout (`--trace`, Dogfooding; netzlos, **kein Gate** — informativ) ([`DC-FA-CLI-009`](spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)) |
| `make doc-complete`          | Vollständigkeits-Dogfood via `d-check` selbst (`--trace --require-complete`, Requirements-Waise ⇒ Exit 1; Dogfooding, netzlos) — **kein** Gate-Bindepunkt; die Closure-Wahrheit bleibt `make completeness-check` ([`DC-FA-CLI-011`](spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code), [ADR-0017](docs/plan/adr/0017-requirements-completeness-gate.md))                              |
| `make semgrep`               | Security-/Static-Analysis-**Gate**: gepinntes semgrep-Image + gepinntes, lokal gecachtes `go/lang/security`-Regelset, netzloser Scan (`--network none`); **Bestandteil von `gates`** ([ADR-0010](docs/plan/adr/0010-semgrep-hermetisches-gate.md)) |
| `make image-scan`            | **CVE-Scan gegen die PUBLIZIERTEN Images** via digest-gepinntes Trivy ([ADR-0066](docs/plan/adr/0066-cve-scan-gegen-das-publizierte-image.md)). Geprüft wird, was Anwender **ziehen** — gegen CVEs ist ein push-getriebenes Gate prinzipiell blind, weil sie ohne Commit auftauchen und das Bild zwischen zwei Releases weiter altert. **Netz ist hier der Zweck, nicht ein Zugeständnis:** eine gepinnte Vuln-DB fände nur die CVEs von gestern; der **Scanner** ist digest-gepinnt, die **DB** bewusst nicht. Zwei Läufe je Image: Vollbericht über alle Schweregrade (fällt nie) und der Entscheidungslauf `CRITICAL`/`HIGH` **mit verfügbarem Fix** — nur der macht rot. **Drei Exit-Codes — die des SKRIPTS:** 0 sauber, 1 behebbare Befunde, 2 **Scan gescheitert** (kein grüner Befundstand). **`make` kollabiert das:** GNU Make normalisiert jeden fehlgeschlagenen Recipe auf Exit 2, über das Target sind 1 und 2 also nicht unterscheidbar — der Nachtlauf liest deshalb die **Ausgabe**, nicht den Exit-Code. Netzlos prüfbar über `bash tools/image-scan.sh --selftest` (sieben Proben der Auswertung; die Trivy-**Feldnamen** deckt er nicht). Beide Trivy-Läufe fahren `--exit-code 0`, weil Trivy einen echten Fehler ebenfalls mit 1 quittiert — gemessen. Kein Docker-Socket. **Netz, bewusst NICHT in `gates`**; gerufen vom Nachtlauf ([`image-scan.yml`](.github/workflows/image-scan.yml)). **Benannte Grenze:** der Fund-Raum ist klein und gemessen — fünf OS-Pakete plus die Go-Modul-Liste des Binaries; ein grüner Lauf sagt „nichts Bekanntes in diesem Raum", nicht „das Image ist sicher" |
| `make freshness-trivy` · `make trivy-digest` | Die zwei Achsen des Scanner-Pins: neuerer Trivy-Release gegen `TRIVY_VERSION`, neuer Bau desselben Tags gegen `TRIVY_DIGEST` ([ADR-0011](docs/plan/adr/0011-digest-pins-build-gate-images.md)). **Netz**, fail-open, **nicht** in `gates`; Nachtlauf |
| `make versions`              | Reproduzierbarkeits-Pins ausgeben (Go, Lint, Basis-Images, Runtime-Image-ID)                                                                                                                                                                       |
| `make build` / `make run`    | Runtime-Image bauen / Selbst-Smoke-Test                                                                                                                                                                                                            |
| `make tidy`                  | `go.mod`/`go.sum` pflegen (`go mod tidy` in Docker; Dependency-Aufnahme/-Hebung — **kein** Gate, bewusster Akt am Dependency-Stand) |
| `make archive-wave-test`     | Testsuite von `tools/archive-wave/` (eigenes `go.mod` — **nicht** Teil von `make test`, das nur das Hauptmodul deckt). Fixture-Verifikation + die drei Umkehr-Proben aus slice-190 |
| `make archive-wave`          | Setzt Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur Schritt 4 um (`WELLE=<id>` Pflicht, `APPLY=1` optional): sammelt die Slices einer geschlossenen Welle (`**Welle:**`-Feld) und ihre Review-Reports, baut `done/<welle-id>/archiv.zip`, ersetzt die Volltexte durch Stubs (Review-Reports ohne Stub), zieht repo-weite Verweise nach. Eigenständiges Werkzeug unter `tools/archive-wave/` (eigenes `go.mod`, eigenes `Dockerfile`) — portabel für jedes Repo mit demselben Planning-Layout, kein Import aus d-checks internen Paketen. **Sicherer Default:** ohne `APPLY=1` wird nichts geschrieben, nur der geplante Umfang angezeigt. **Kein Gate**, bewusster, von Hand ausgelöster Vorgang (slice-190) |
| `make deps` / `make compile` | Cache-Layer / schnelles Compile-Feedback                                                                                                                                                                                                           |
| `make guard-probe`           | Fährt den Tool-Call-Wächter aus §3.1 gegen seine Proben: je Fall Erwartung und Ergebnis, Fehlschlag-Zähler am Ende. **Werkzeug-lokal und bewusst nicht in `gates`** — der Wächter ist keine Repo-Invariante, sondern eine Werkzeug-Einstellung; ohne wiederholbare Proben wäre seine Zusage aber eine Erinnerung |
| `make record-gates`          | Nachweis schreiben: Working-Tree-Hash für den Stop-Hook                                                                                                                                                                                            |
| `make help` / `make clean`   | Targets anzeigen / Images entfernen                                                                                                                                                                                                                |

Alle dokumentiert-geplanten Targets existieren; Details und Bindungen:
Sensors-Tabelle in [`harness/README.md`](harness/README.md).

## 5. Dokumentations-Regeln

- Commits/PRs müssen mindestens eine `DC-*`-, `ADR-*`-, `MR-*`- oder
  `slice-*`-ID nennen (maschinell erzwungen: `make trace-check` /
  `commit-msg`-Hook / PR-CI — seit dem Modul `commits` dogfooded über das
  eigene Image, [ADR-0027](docs/plan/adr/0027-commits-traceability-modul.md);
  die abgelöste Skript-Mechanik trug
  [ADR-0013](docs/plan/adr/0013-pr-ci-und-traceability-gate.md).
  Ausnahme: Merge-/Revert-Commits). Vergeben werden IDs nur beim
  Spec-/ADR-Schreiben nach dem deklarierten Schema
  ([`MR-008`](harness/conventions.md#mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage))
  — nie ad hoc im Commit/PR; Agenten referenzieren IDs, sie erfinden
  keine. Struktur-IDs (`SPEC-<NNN>`/`ARC-<NNN>`,
  [`MR-000`](harness/conventions.md#mr-000--baseline-aussage)) entstehen nur
  beim Schreiben der Spec-Straten — fortlaufend je Datei — und gehören
  **nicht** in Commit-Botschaften.
- **Dependabot-Commits tragen die Kennung im Präfix, nicht in einer Ausnahme.**
  `commit-message.prefix` in [`.github/dependabot.yml`](.github/dependabot.yml)
  lautet `build(deps) [ADR-0067]` bzw. `build(ci) [ADR-0067]`; damit erfüllt <!-- d-check:ignore (literale Konfigurationswerte, keine Verweise) -->
  jeder Bump-Commit dieselbe Regel wie jeder andere
  ([ADR-0067](docs/plan/adr/0067-dependabot-als-hebender-kanal.md)). Die
  naheliegende Alternative — `commits.exempt-pattern` erweitern — hätte den Gate
  für eine **ganze Commit-Klasse** blind gemacht; das ist der Grund, nicht §3.6.
  §3.6 **verbietet** eine Lockerung nicht, es verlangt eine ADR dafür — die
  Alternative wäre also dokumentiert zulässig gewesen und ist aus dem Sachgrund
  verworfen, nicht aus einem Verfahrensgrund. **Die Kennung gilt dem Kanal,
  nicht dem Inhalt des einzelnen Bumps**; wer mehr Bezug hineinliest, liest zu
  viel.
- Neue oder geänderte `DC-*`-Anforderungen entstehen nur in
  [`spec/lastenheft.md`](spec/lastenheft.md) — nie per ADR (ADRs
  schärfen die Spezifikation, nicht das Lastenheft). Der
  Anlege-Prozess (Akzeptanzkriterien-Trio, Versions-Bump + Historie,
  Beleg-Pflicht) folgt dem Baseline-Regelwerk
  ([`modul-03-spec`](.harness/baseline/v5.18.0/regelwerk/modul-03-spec.md)); das
  repo-spezifische ID-Schema steht in `spec/lastenheft.md` §3.
- Neue ADRs müssen den ADR-Index aktualisieren.
- Neue ADRs tragen die Sektion `## Re-Evaluierungs-Trigger` (oder „permanent");
  die vor Einführung `Accepted`-ADRs sind immutable und **grandfathered** (das
  Trigger-Feld liegt im ADR-Core, nachträgliches Ergänzen bräche `make adr-check`).
  Der Welle-Closure-Trigger-Audit (Baseline-Regelwerk Modul 6) bestätigt oder
  revidiert sie (Folge-ADR mit `supersedes`).
- Roadmap/Status-Geschichte lebt in `docs/plan/planning/`, nicht in der Architektur-Spec.
- Slice-Lifecycle (`open → next → in-progress → done`) ist reine Datei-Bewegung (`git mv`, siehe §3.3).
- Neue Slice-Köpfe tragen das Feld `**Verantwortlich:**` (Rolleninhaber der
  Implementer-Rolle, gesetzt **spätestens bei der Beanspruchung** — beim
  Move `open→next` bzw. direkt `open→in-progress`, wie dieses Repo ihn
  fährt; Deklaration, kein Sensor — Baseline v5.5.0, template-forward,
  kein Retrofit).
- Das Slice-Kopf-Feld `**Berührte Spec-Stellen:**` nennt die **Kennung**, wo
  das Zielelement eine trägt (`SPEC-<NNN>`, `ARC-<NNN>`,
  `<DC-ID>.<Buchstabe>`), sonst den Abschnitt; `—`, wenn der Slice keine
  Spec-Stelle berührt. Der Verweis zeigt **aufwärts** — die Spec nennt den
  Slice nie (§3.4). Feld-Form aus der Baseline-`slice.template.md`, die
  Kennungs-Regel aus
  [`MR-000`](harness/conventions.md#mr-000--baseline-aussage).
- **Was eine Welle einlöst, gehört in ihren Closure-Trigger — nicht in die DoD
  eines Slice.** Ein DoD-Punkt, den der Slice **selbst nicht abhaken kann** (das
  Release, das erst mit der Welle fällt), zwingt ihn, mit offenem Haken zu
  schließen. Damit ist der Haken als Zustandsfeld unbrauchbar: er sagt nicht
  mehr *„hier fehlt etwas"*, sondern *„hier fehlt vielleicht etwas"*. Der
  Wellen-Closure-Trigger trägt ihn; die Slice-DoD nennt ihn gar nicht.
  **Bestands-Grenze:** vor dieser Regel geschriebene `done/`-Slices behalten
  ihre Form — ein nachträglich umgeschriebener DoD-Punkt fälschte einen
  Lauf-Beleg.

  **Seit slice-172 hält das ein Sensor**, und zwar am **Ruheort**: eine
  `structure`-Regel im Closure-Profil meldet jeden offenen DoD-Haken eines
  `done/`-Slice (`max-open-tasks: 0` ⇒ `section-tasks-open`, je Haken auf
  seiner Zeile, mit verfasstem Reparatur-Hinweis). Sie läuft in
  `make verify-closure-notes`, **nicht** in `gates` — sonst meldete sie beim
  Arbeiten an einem laufenden Slice. Der Altbestand bis
  [`slice-170`](docs/plan/planning/done/slice-170-workflows-modul.md) ist mit
  fester Ziffernzahl ausgenommen
  ([`MR-056`](harness/conventions.md#mr-056)). **Drei Grenzen gehören dazu:**
  ein Haken ist eine **Selbstauskunft** — die Regel verschiebt die Lücke von
  *unsichtbar* nach *behauptet* und prüft keinen Review; und ein **vergessener
  Schluss-Fence** macht die **Bedingung** blind (isoliert gemessen: 0 Befunde, Exit 0), weshalb dasselbe Profil `spans` fährt — `fence-unclosed` meldet den Fall. **Der Bindepunkt als Ganzes wird davon nicht grün:** im heutigen Profil melden Nachbarregeln, und `spans` nennt die Ursache. Und **ein Haken INNERHALB eines wohlgeformten Fenced-Blocks ist unsichtbar** — dort meldet auch `fence-unclosed` nichts; dieselbe Fence-Treue, die eine Illustration schützt, ist der Weg, einen Haken zu verstecken. Und **ein Haken IM Fenced-Block ist unsichtbar** — wohlgeformt, also auch ohne `fence-unclosed`; das ist dieselbe Fence-Treue, die eine Illustration schützt, und zugleich der Weg, einen Haken zu verstecken.
- Slice-Pläne tragen **kein** `**Status:**`-Feld — der Lifecycle-Zustand **ist** die
  Verzeichnis-Position; neue Slices führen stattdessen den `**Lifecycle:**`-Hinweis
  (Baseline-`slice.template.md`). Alt-Slices in `done/` behalten ihr historisches Feld.
- Jeder Slice-Plan trägt **vor** der Sub-Area-Modus-Begründung die **drei**
  Vorprüfungen: Sub-Area prüfen · offene Beobachtungen im Register
  `observations.md` sichten (beide Baseline-Regelwerk Modul 5/6, unabhängig vom
  Sub-Area-Modus) · den **Nachtlauf-Stand** lesen (`make nightly-state`). Die
  dritte ist eine Adaption ([`MR-053`](harness/conventions.md#mr-053)): der
  Kanon kennt keinen Nachtlauf. Sie hängt an diesem Moment, weil dort ohnehin
  gelesen wird — **benannte Grenze:** in einer Pause liest niemand. Der dritte
  Block entsteht **spätestens bei der Beanspruchung**; ein Plan in `open/`
  trägt ihn noch nicht.
  **Die beiden kanonischen Blöcke belegen ihre Regel** mit einer
  `d-check:cite`-Direktive auf die **vorschreibende** Regelwerk-Zeile, samt
  wörtlichem Zitat darunter ([`MR-054`](harness/conventions.md#mr-054));
  `citations` prüft es wortgleich im inneren Loop, ein falsch angekerter Beleg
  wird rot. Der dritte Block trägt bewusst keine — sein Ziel ist repo-eigen und
  meldete bei jeder Änderung. **Kein Sensor hält das:** ein Plan ganz ohne
  Direktiven ist grün.
- Eine Commit-Botschaft oder Closure-Notiz behauptet **nicht mehr, als die
  Arbeit trägt**: eine genannte Probe muss gelaufen sein (§6 Schritt 8), und ihr
  Schluss reicht **nicht weiter als die gemessene Menge** — wer N Formen
  geprüft hat, berichtet N; „damit ist X allgemein" ist eine andere Aussage als
  die gemessene. Beides ist Urteil, kein `grep`; der Reviewer-Skill trägt die
  Anker dazu. *(Hard Rule aus dem Steering Loop,
  [Beobachtungs-Register](docs/plan/planning/observations.md) `BEO-009`, seit
  welle-82; Auflösungs-Trigger: permanent.)*
- Eine zitierte Quelle trägt **nur, was in ihrem Geltungsbereich steht**. Vor
  jedem Verweis das **Feld** lesen, nicht den Titel: bei `MR-<NNN>` den
  `Geltungsbereich` **und** `Ersetzt-Baseline-Regel`, bei einer Kanon-Stelle den
  Absatz, bei einer ADR die Unterscheidung **Akt gegen stehendes Verbot**. Und
  die **direkteste** Quelle wählen — für eine Regel dieser Datei ist das ihre
  Vorlage, nicht eine andere. Ein Zitat sieht aus wie ein Beleg, auch wenn es
  keiner ist; das macht die Klasse beim Schreiben unsichtbar und im Review
  auffindbar. Urteil, kein `grep`; der Reviewer-Skill trägt den Anker dazu.
  Kanon:
  [`grundlagen-source-precedence.md` §Wie weit trägt ein zitierter Satz](.harness/baseline/v5.18.0/regelwerk/grundlagen-source-precedence.md)
  — dort als Frage an **jede** zitierte Aussage, hier als operative Form für den
  Implementer. *(Hard Rule aus dem Steering Loop,
  [Beobachtungs-Register](docs/plan/planning/observations.md) `BEO-012`, seit
  slice-147; Auflösungs-Trigger: permanent.)*
- `CHANGELOG.md` wird bei nutzersichtbaren Änderungen gepflegt — **in der
  Release-Prep, nicht im Feature-Commit.** Die Datei führt **keinen**
  `[Unreleased]`-Abschnitt: jeder Eintrag steht unter seiner Versions-Nummer,
  und die steht erst fest, wenn das Release geschnitten wird. Ein Slice, der
  seine Zeile vorzieht, muss sie beim Bump wieder anfassen. Dieselbe Grenze
  gilt den beiden `README*.md` und dem Handbuch-Kopf. **Gemessen, nicht
  vereinbart:** die Feature-Commits der letzten Slices fassen `CHANGELOG.md`
  nicht an. Ohne diesen Satz meldet jede Verifikation den Rückstand erneut —
  zu Recht, denn die Regel darüber sagte nur *„wird gepflegt"* und nicht
  *wann*.

## 6. Minimal Agent Workflow

Pro Slice:

1. [`harness/README.md`](harness/README.md) lesen.
2. Relevante kanonische Quelle lesen (Source Precedence beachten).
3. Betroffene Requirement-/ADR-IDs identifizieren — und **vor der
   Implementierung benennen**: Slice-ID, betroffene `DC-*`-IDs, ADR-IDs,
   betroffene Module, auszuführende Gates
   ([`MR-031`](harness/conventions.md#mr-031)).
4. Kleinste sinnvolle Änderung planen.
5. Engsten nützlichen Sensor laufen lassen.
6. Repo-weiten Gate-Lauf vor Handoff (`make gates`).
7. Doku/Indizes aktualisieren, falls ein öffentlicher Vertrag berührt.
8. Ausgeführte Sensors und verbleibende Risiken berichten — keine Erfolgsmeldung ohne
   Gate-Ausführung **und ohne ihre echte Ausgabe**: ein behaupteter Exit-Code ist keiner.
