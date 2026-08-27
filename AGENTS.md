# AGENTS.md — Briefing für AI-Coding-Agenten

## 1. Was diese Datei ist

Onboarding-Briefing für jede AI-Session, die in diesem Repo Code oder
Dokumentation ändert. Sie verweist auf die kanonischen Quellen und
formuliert die Hard Rules, die der Implementation-Agent immer
einhalten muss.

Diese Datei trägt **Hard Rules und Pointer** auf die kanonischen Quellen und
**dupliziert deren Inhalt nicht** — sonst entsteht Drift (Kanon:
[`modul-09-implementierung.md` §AGENTS.md-Regeln](.harness/baseline/v5.12.0/regelwerk/modul-09-implementierung.md#agentsmd-regeln-modul-9)).

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
[`lab-regelwerk.zip`](https://github.com/pt9912/ai-harness-course/releases/download/v5.12.0/lab-regelwerk.zip);
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
([`modul-02-harness-bootstrap.md` §Freshness-Audit](.harness/baseline/v5.12.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2)
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
braucht `curl`. Beide sind fail-open und stehen bewusst außerhalb von `gates`;
wer sie fährt, fährt sie mit dieser zusätzlichen Erwartung.

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

**Ausnahme MR-/Wellen-Lifecycle-Move** (`conventions/` → `conventions/done/`,
flaches Wellendokument → `done/`): hier trägt der Move-Commit die
**Link-Tiefen-Fixes der bewegten Datei selbst** mit — ein reiner Move wäre
`doc-check`-rot, weil die relativen Verweise vom neuen Ort nicht mehr
auflösen. Alles Übrige bleibt Commit 2; sinkt der Rename-Score dadurch
Richtung 50 %, deklariert die Commit-Botschaft den Move ausdrücklich als
`git mv`. Kanonisch:
[`MR-013`](harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise).

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
Grenze** ([Baseline §Was ein Kommentar trägt](.harness/baseline/v5.12.0/regelwerk/grundlagen-harness-dateien.md#was-ein-kommentar-trägt--code-konfiguration-skripte)).
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
[Baseline §Was ein Kommentar trägt](.harness/baseline/v5.12.0/regelwerk/grundlagen-harness-dateien.md#was-ein-kommentar-trägt--code-konfiguration-skripte).

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

**Begründung:** Supply-Chain-Härtung — ein Tag lässt sich umhängen, ein SHA
nicht; dieselbe Härtung wie der Docker/make-only-Pfad in §3.1.

**Durchgesetzt:** `make workflow-pins` in `make gates`. Er trägt die **Form** —
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
| `make doc-check`             | Doku-Links, Anker, Kennungs-Linkpflicht, Referenzmatrix, Inline-Code-Pfade, Abschnitts-Invarianten (Modul `structure`) + Kennungen in Diagramm-Fences (Modul `diagrams`) via `d-check` selbst (Dogfooding; netzlos — zugleich [`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Messmethode)             |
| `make workflow-pins`         | Jeder `uses:`-Schlüssel in `.github/workflows/` nennt einen **vollen 40-stelligen Commit-SHA** mit Tag-Kommentar dahinter (`uses-pin-missing` / `uses-pin-untagged`). Geprüft wird die **Form**, nicht die Gültigkeit des SHA — das wäre Netz. Liest `.yml` **und** `.yaml` und nur echte YAML-Schlüssel, nicht die Prosa der Workflow-Köpfe; **fail-closed auch bei leerer Prüfmenge**. Netzlos, **Bestandteil von `gates`** (§3.9)                                                                                              |
| `make baseline-verify`       | Vendorte Baseline gegen `SHA256SUMS`: `sha256sum -c` (geänderte/gelöschte Datei) **plus Manifest-Deckung** (zusätzlich eingelegte Datei — eine untermengige, in sich konsistente `SHA256SUMS` passierte sonst grün). Netzlos, fail-closed, **Bestandteil von `gates`** ([`MR-011`](harness/conventions.md#mr-011)-Kette)                                                    |
| `make gates`                 | alle inneren Gates (mandatory vor Handoff)                                                                                                                                                                                                         |
| `make ci`                    | CI-äquivalenter Lauf: gates + image-test (fährt die Release-Pipeline)                                                                                                                                                                              |
| `make trace-check`           | Traceability-Gate **via Modul `commits`** (Image, dogfood): DC-/ADR-/MR-/slice-ID in Commit-Messages (`commit-untraceable`; `RANGE=`-Range für CI, `MSGFILE=` für den `commit-msg`-Hook via stdin; bewusst **nicht** Teil von `gates`/`ci`) ([ADR-0027](docs/plan/adr/0027-commits-traceability-modul.md) löst die Skript-Mechanik von [ADR-0013](docs/plan/adr/0013-pr-ci-und-traceability-gate.md) ab; [`DC-FA-COMMITS-001`](spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in)) |
| `make adr-check`             | ADR-Immutable-Gate **via Modul `vcs`** (Image, dogfood): `Accepted`-ADRs nicht inhaltlich ändern (`RANGE=`/`STAGED=`-Modi; `pre-commit`-Hook + PR-CI; bewusst **nicht** Teil von `gates`/`ci`) ([ADR-0024](docs/plan/adr/0024-vcs-immutable-gate.md) löst die Skript-Mechanik von [ADR-0016](docs/plan/adr/0016-adr-immutable-gate.md) ab, [ADR-0025](docs/plan/adr/0025-codepaths-ignore-refs.md) entfernt das Alt-Skript)                                                            |
| `make freshness-go` · `make freshness-golangci` | Melden, wenn upstream eine neuere Go-Version bzw. einen neueren `golangci-lint`-Release führt als der Pin. Zwei Quellen-Formen: `go.dev/VERSION?m=text` (Plaintext — golang/go führt **keine** Release-Objekte) und der `releases/latest`-Redirect. **Netz**, fail-open (Ausfall ⇒ `SKIP`, Exit 0), bewusst **nicht** in `gates`/`ci`; gerufen vom Nachtlauf ([`upstream-drift.yml`](.github/workflows/upstream-drift.yml)). Der Vergleicher ist über `--compare` **netzlos** prüfbar. Sie melden nur — die Hebung bleibt ein bewusster Akt, und beim Go-Bump zieht das `golangci`-Pendant mit |
| `make runtime-base-digest`   | Meldet, wenn der **Runtime-Basis-Tag** `distroless/static-debian12:nonroot` inzwischen einen **anderen Digest** trägt. Andere Frage als die zwei Achsen darüber: nicht „neuerer Tag", sondern „neuer Bau desselben Tags" — für diesen Pin die **einzige** Handhabe, weil sein Tag keine Version führt. Quelle ist `docker buildx imagetools`; **Netz**, fail-open (Ausfall ⇒ `SKIP`), bewusst **nicht** in `gates`, gerufen vom Nachtlauf. Meldet nur — die Hebung bleibt ein bewusster Akt ([ADR-0011](docs/plan/adr/0011-digest-pins-build-gate-images.md)) |
| `make baseline-freshness`    | Upstream-Audit des Baseline-Pins: neuerer Release-Tag (Currency, gelesen aus der Release-**Liste** — nicht `releases/latest`, das Prereleases überspringt) + Content-Drift am gepinnten Tag. **Netz**, fail-open (Ausfall ⇒ `SKIP` je Teil), bewusst **nicht** in `gates`/`ci` — der netzlose innere Lauf ist eine Eigenschaft dieses Repos, keine Zusage des Produkts; gerufen vom Nachtlauf ([`upstream-drift.yml`](.github/workflows/upstream-drift.yml)). Meldet nur — die Hebung bleibt ein bewusster Akt |
| `make hooks`                 | git-Hooks installieren (`core.hooksPath` → `.githooks`; aktiviert `commit-msg`-Traceability + `pre-commit` mit ADR-Immutable via Modul `vcs` **und** dem vollen `doc-check` als Doku-GUARD — seit welle-79: ein roter Gate-Exit kann keine Shell-Kette mehr passieren) ([ADR-0013](docs/plan/adr/0013-pr-ci-und-traceability-gate.md), [ADR-0016](docs/plan/adr/0016-adr-immutable-gate.md), [ADR-0024](docs/plan/adr/0024-vcs-immutable-gate.md))    |
| `make completeness-check`    | Requirements-Completeness-Gate **via in-Produkt-Flag** `--trace --require-complete` (≥1 Waise ⇒ Exit 1, mit `WAISE`-Zeilen + Anzahl); **Closure-Bindepunkt** (in `make fullbuild`, **nicht** `gates`/`ci`) ([ADR-0026](docs/plan/adr/0026-completeness-in-product-gate.md) löst die Skript-Mechanik von [ADR-0017](docs/plan/adr/0017-requirements-completeness-gate.md) ab; [`DC-FA-CLI-011`](spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)) |
| `make verify-closure-notes`  | Struktur des `done/`-Bestands: die Closure-Notizen **via Modul `planning`** (Abschnitt vorhanden, Substanz außerhalb Code, keine deklarierte Floskel, opt-in kein Vorlagen-Platzhalter) **und** Abschnitts-Invarianten **via Modul `structure`** (`section-*`-Codes), darunter die **urteilsfreie Hälfte der Drei-Ausgänge-Regel** des Baseline-Regelwerks (`modul-05`: *„Ein Slice geht nicht nach `done/`, während ein Risiko ohne Ausgang dasteht"*) in **zwei** Regeln: der **vergessene** Ausgang als `forbid-pattern` über jeden H1-Abschnitt, der **erfundene** als Komplement-`forbid-pattern` über §5 — der Wortschatz ist geschlossen und umfasst alle drei Kanon-Ausgänge ([`MR-049`](harness/conventions.md#mr-049)). **Sie sieht ein Risiko ganz ohne Ausgangs-Marker nicht**, und der Altbestand vor slice-140 ist ausgenommen. **Ob** ein eingetragener Ausgang inhaltlich trägt, bleibt Urteil. **Vier Bereiche sieht sie nicht** — Überschriften-Zeile, Fenced Blocks, Inline-Code (die drei gewollt, sonst meldete ein Slice über Platzhalter seine eigene Doku) und, als ausgewiesener Preis, Prosa innerhalb einer **absatzweiten** Inline-Code-Spanne ([ADR-0059](docs/plan/adr/0059-closure-waechter-weicht-structure-regel.md)) — beide über dasselbe `--config`-Profil. Fährt ein **eigenes** Prüf-Profil über `--config` ([`DC-FA-CLI-012`](spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben)); **Closure-Bindepunkt** (in `make fullbuild`, bewusst **nicht** `gates`/`ci`) ([ADR-0048](docs/plan/adr/0048-closure-note-struktur-im-planning-modul.md), [`DC-FA-PLAN-001`](spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)) |
| `make fullbuild`             | volle Closure: gates + image-test + bench + completeness-check + verify-closure-notes, schließt mit dem Image-Hash                                                                                                                                                                             |
| `make image-test`            | [`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image)-Akzeptanzkriterien gegen das lokale Image (nativ vs. Container)                                                                                                                |
| `make bench`                 | [`DC-QA-01`](spec/lastenheft.md#dc-qa-01--performance)-Benchmark gegen generiertes Fixture (Median aus 3 Läufen, kein Gate in `gates`)                                                                                                             |
| `make trace`                 | Requirements Traceability Matrix via `d-check` selbst auf stdout (`--trace`, Dogfooding; netzlos, **kein Gate** — informativ) ([`DC-FA-CLI-009`](spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)) |
| `make doc-complete`          | Vollständigkeits-Dogfood via `d-check` selbst (`--trace --require-complete`, Requirements-Waise ⇒ Exit 1; Dogfooding, netzlos) — **kein** Gate-Bindepunkt; die Closure-Wahrheit bleibt `make completeness-check` ([`DC-FA-CLI-011`](spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code), [ADR-0017](docs/plan/adr/0017-requirements-completeness-gate.md))                              |
| `make semgrep`               | Security-/Static-Analysis-**Gate**: gepinntes semgrep-Image + gepinntes, lokal gecachtes `go/lang/security`-Regelset, netzloser Scan (`--network none`); **Bestandteil von `gates`** ([ADR-0010](docs/plan/adr/0010-semgrep-hermetisches-gate.md)) |
| `make versions`              | Reproduzierbarkeits-Pins ausgeben (Go, Lint, Basis-Images, Runtime-Image-ID)                                                                                                                                                                       |
| `make build` / `make run`    | Runtime-Image bauen / Selbst-Smoke-Test                                                                                                                                                                                                            |
| `make tidy`                  | `go.mod`/`go.sum` pflegen (`go mod tidy` in Docker; Dependency-Aufnahme/-Hebung — **kein** Gate, bewusster Akt am Dependency-Stand) |
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
- Neue oder geänderte `DC-*`-Anforderungen entstehen nur in
  [`spec/lastenheft.md`](spec/lastenheft.md) — nie per ADR (ADRs
  schärfen die Spezifikation, nicht das Lastenheft). Der
  Anlege-Prozess (Akzeptanzkriterien-Trio, Versions-Bump + Historie,
  Beleg-Pflicht) folgt dem Baseline-Regelwerk
  ([`modul-03-spec`](.harness/baseline/v5.12.0/regelwerk/modul-03-spec.md)); das
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
- Slice-Pläne tragen **kein** `**Status:**`-Feld — der Lifecycle-Zustand **ist** die
  Verzeichnis-Position; neue Slices führen stattdessen den `**Lifecycle:**`-Hinweis
  (Baseline-`slice.template.md`). Alt-Slices in `done/` behalten ihr historisches Feld.
- Jeder Slice-Plan trägt **vor** der Sub-Area-Modus-Begründung die zwei
  **Vorprüfungen** (Sub-Area prüfen · offene Beobachtungen im Register
  `observations.md` sichten) — Baseline-Regelwerk Modul 5/6, unabhängig vom
  Sub-Area-Modus.
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
  [`grundlagen-source-precedence.md` §Wie weit trägt ein zitierter Satz](.harness/baseline/v5.12.0/regelwerk/grundlagen-source-precedence.md)
  — dort als Frage an **jede** zitierte Aussage, hier als operative Form für den
  Implementer. *(Hard Rule aus dem Steering Loop,
  [Beobachtungs-Register](docs/plan/planning/observations.md) `BEO-012`, seit
  slice-147; Auflösungs-Trigger: permanent.)*
- `CHANGELOG.md` wird bei nutzersichtbaren Änderungen gepflegt.

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
