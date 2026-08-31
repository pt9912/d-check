# d-check

> 🇩🇪 **Deutsch** · 🇬🇧 [English](README.md)

Konsistenz-Checker: prüft die Dokumentation eines Repos gegen dessen
tatsächlichen Zustand — Markdown-Referenzen (Links, Bilder, Anker, Kennungen)
plus Deklarationen (Build-Targets, Versions-Pins, Commit-Trace, Planning,
Workflow-Referenzen).
Deterministisch, seiteneffektfrei, als Container-Image.

## Was ist d-check?

**d-check** prüft Markdown-Dokumentation als prüfbares **Invarianten-Netz**: jede
maschinell entscheidbare Doku-Invariante ist ein einzeln aktivierbares Regelmodul
mit eigener Anforderung im [Lastenheft](spec/lastenheft.md) — vom **Referenz-Netz**
(Links, Anker, ID-Linkpflicht, Referenzmatrix) über Markdown-Hygiene (Span-Artefakte,
Host-Pfad-Leaks), Content-Drift und Immutabilität (Content-/Core-Pins, git-Diff) bis
zu Versions-Pin-, Commit-Traceability-, Planning-Lifecycle- und
Getrackt-Status-Konsistenz, bis hin zu Struktur-Invarianten **innerhalb** eines
Dokuments:

- `links` — lokale Link- und Bildreferenzen: Ziel existiert, kein
  Repo-Escape; opt-in `resolve-from`: Dateien in **wandernden**
  Lifecycle-Verzeichnissen lösen jedes relative Ziel von jedem Ort ihrer
  Gruppe auf — gemeldet **vor** dem `git mv`
  ([`DC-FA-LINK-001`](spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links))
- `anchors` — Heading-Anker (GitHub-Slug-Verfahren) und Inline-HTML-Anker
  (`<a name>`, `id=`)
  ([`DC-FA-ANCH-001`](spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors))
- `ids` — Linkpflicht für Kennungen (z. B. `ADR-NNNN`) nach
  deklarierten Mustern ([`DC-FA-ID-001`](spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids))
- `matrix` — Referenzrichtungs-Regeln zwischen Dokumentklassen
  ([`DC-FA-MTX-001`](spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)),
  **innerhalb** einer geordneten Klasse (`order`/`direction` ⇒ `matrix-downward`,
  [`DC-FA-MTX-002`](spec/lastenheft.md#dc-fa-mtx-002--verweisrichtung-innerhalb-einer-geordneten-dokumentklasse-modul-matrix))
  und auch als **bare ID-Token** im Körper (`token` ⇒ `matrix-forbidden`,
  [`DC-FA-MTX-003`](spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix))
  plus Status-Bedingungen; Ausnahmen sind strukturell (`exclude-sections`,
  `allow-supersede-lineage`) oder per Provenance-Marker
  `<!-- d-check:status-provenance -->` deklariert
- `external` — Erreichbarkeit externer URLs, strikt opt-in
  ([`DC-FA-EXT-001`](spec/lastenheft.md#dc-fa-ext-001--externe-links-modul-external-opt-in))
- `sources` — Content-Pin externer Quellen gegen Upstream-Drift, opt-in und —
  neben `external` — die **zweite** Netz-Tür: eine per Marker
  `<!-- source-pin: … -->` am `http(s)`-Link oder per Config-Block `sources:` auf
  einen `sha256` gepinnte Quelle wird geholt, gehasht und verglichen
  (`source-drift` mit vollem Ist-Hash, `source-unreachable` bei Netzfehler);
  Einzeldatei oder Archiv (`unpack: zip`)
  ([`DC-FA-SRC-001`](spec/lastenheft.md#dc-fa-src-001--upstream-content-drift-externer-quellen-modul-sources-opt-in-netz))
- `codepaths` — explizite Pfade in Inline-Code, opt-in; das opt-in `check-lines`
  verifiziert `datei:<von>-<bis>`-Zeilen-Referenzen (`citation-out-of-range`,
  `citation-inverted-range`)
  ([`DC-FA-CODE-001`](spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in))
- `spans` — Markdown-Span-Artefakte (ungeschlossene Code-Spans,
  verschachtelte Links), opt-in
  ([`DC-FA-SPAN-001`](spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in))
- `hostpaths` — host-lokale absolute Pfade (Maschinen-Layout-Leaks),
  opt-in
  ([`DC-FA-HOST-001`](spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in))
- `diagrams` — Kennungs-Existenz in Diagramm-Fences (z. B. `mermaid`): jede
  im Diagramm gefundene Kennung muss in ihrer `defined-in`-Quelle definiert
  sein, opt-in; **beide Ventile** wie `ids`/`codepaths` — `exempt-paths`
  (datei-weit) und der Zeilen-Marker, hier als **Token** statt HTML-Kommentar
  und auf der **Öffnungszeile** des Fence für den ganzen Block
  ([`DC-FA-DIAG-001`](spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in))
- `versions` — Versions-Pin-Konsistenz: gepinnte `ghcr`-Image-Verweise müssen
  die aktuelle Version (aus `version.md#aktuell`) tragen, liest auch
  Fenced-Code, opt-in; der **Anker** von `current-from` folgt dabei derselben
  Antwort wie in `anchors`, nicht der Fence-Ausnahme. `versions.patterns` trägt
  **mehrere** Muster-Quellen-Paare (eigenes Release **und** ein fremder
  gepinnter Stand) — die Kurzform ist die einelementige Liste, Ausnahmen gelten
  paar-lokal
  ([`DC-FA-VER-001`](spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in))
- `pins` — Content-Pin gegen inhaltlichen Drift: ein Link mit
  `<!-- dpin: … -->` wird gegen den Hash seines Ziel-Spans geprüft (Befund
  `link-stale` bei Drift), opt-in pro Link; **welcher Anker den Span adressiert,
  entscheidet dieselbe Antwort wie in `anchors`** (Duplikat-Slug,
  Prozent-Dekodierung, case-sensitiv, nur außerhalb von Code)
  ([`DC-FA-PIN-001`](spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in))
- `immutable` — Immutabilitäts-Pin gegen Core-Drift: eine Datei mit
  `<!-- immutable: … -->` wird gegen den Hash ihres normalisierten **Core**
  (ohne Marker-Zeile + `exclude-sections`) geprüft (Befund `core-drift` bei
  Drift), opt-in pro Datei; hermetisch (kein git, read-only-Arbeitsbaum)
  ([`DC-FA-IMM-001`](spec/lastenheft.md#dc-fa-imm-001--immutabilitäts-pin-gegen-core-drift-modul-immutable-opt-in))
- `vcs` — git-Diff-Immutabilität des Core über eine Commit-Range: mechanisiert die
  ADR-Immutabilität als verteilbares Modul (`core-drift-vcs`), reine-Go-git im
  read-only `.git` (**kein** git-Binary, **kein** Netz), opt-in
  ([`DC-FA-VCS-001`](spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in))
- `commits` — Traceability-Kennung in Commit-Messages über eine Range (`--range`)
  bzw. der Pending-Message (`--commit-msg`): jede Commit-Message trägt eine
  `DC-`/`ADR-`/`MR-`/`slice-`-ID (`commit-untraceable`), teilt den VCS-Port mit `vcs`,
  opt-in
  ([`DC-FA-COMMITS-001`](spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in))
- `planning` — Lifecycle-Konsistenz, alle drei Seiten. **Eintritt:** der Ruhe-Marker
  steht im `## Aktuelle Welle`-Block genau dann, wenn kein `slice-*` im
  Verzeichnis liegt (`planning-drift`). **Austritt** (zusätzlich opt-in über
  `closure.dir`): die **Struktur** der Closure-Notizen abgeschlossener Pakete —
  Abschnitt vorhanden, genug Satzende-Zeichen außerhalb von Code-Blöcken, keine
  deklarierte Floskel (`closure-note-missing`/`-thin`/`-boilerplate`). Prüft
  Struktur, nicht Bedeutung. Hermetisch (kein git), fail-closed bei
  fehlender/mehrdeutiger Überschrift, fehlendem Closure-Verzeichnis und bei null
  Kandidaten. **Wellen-Register** (zusätzlich opt-in über `waves.dir`): die
  Wellen-Abschnitte der Roadmap gegen die Wellen-Dateien — aktive Welle ⟺
  flaches Wellendokument (`waves.mode: one`, Default) **oder**
  Kennungs-Bijektion Aktiv-Block ⟺ Dateien, Ruhe-Marker außen vor
  (`waves.mode: many`); Vorschau ohne Datei, Abschluss-Register ⟺
  Ergebnisnotizen beidseitig (`wave-drift`, `wave-preview-exists`,
  `wave-results-missing`, `wave-unregistered`), opt-in
  ([`DC-FA-PLAN-001`](spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in))
- `tracked` — Getrackt-Status auflösbarer, **existierender** Link-/Bild-Ziele
  gegen den git-**Index** (`target-untracked`: ein untracktes/gitignoriertes
  Ziel fehlt auf jedem frischen Klon); Index-Wahrheit (gestagt = getrackt,
  keine `.gitignore`-Interpretation), liest `.git` read-only ohne Range, opt-in
  ([`DC-FA-TRK-001`](spec/lastenheft.md#dc-fa-trk-001--getrackt-status-auflösbarer-referenz-ziele-modul-tracked-opt-in))
- `targets` — Deklarations-Konsistenz zwischen Doku und Build-Targets: ein in
  einer Doku-**Tabellenzeile** als `make X` behauptetes Target ohne
  Makefile-Regel (`gate-phantom`), oder eine Makefile-Regel ohne Eintrag in der
  Autoritäts-Doku (`gate-undocumented`); **hermetisch** (kein git, kein
  Makefile-Ausführen), fail-closed, opt-in
  ([`DC-FA-TGT-001`](spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in))
- `citations` — Verbatim-Zitat-Verifikation: die Direktive
  `<!-- d-check:cite <pfad>:<von>-<bis> -->` markiert das folgende Zitat (ein
  `>`-Blockquote oder inline `„…"`/`"…"`; Leerzeilen dazwischen sind unschädlich,
  ein Code-Block dazwischen trennt); der whitespace-normalisierte Zitattext muss
  ein zusammenhängender Teilstring der Quell-Spanne sein (`citation-mismatch`); die
  Direktive zählt nur **außerhalb** von Code-Blöcken **und außerhalb von
  Inline-Code** — in Backticks ist sie eine Erwähnung, keine Direktive —, während
  die zitierte Quell-Spanne **roh** gelesen wird, opt-in
  ([`DC-FA-CITE-001`](spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in))
- `structure` — Struktur-Invarianten **innerhalb** eines Dokuments: je Regel eine
  Dokumentklasse über **eigene** Globs, ein Abschnitt (Klartext oder RE2) und bis
  zu **zehn** Bedingungen mit je eigenem Grund-Code — nicht leer (`section-empty`),
  Mindest-Sätze (`section-thin`), Task-Obergrenze (`section-oversized`),
  **offene** Task-Items auf den **rohen** Zeilen (`max-open-tasks` ⇒
  `section-tasks-open`, ein Befund je Haken auf **seiner** Zeile — anders als
  `max-tasks` immun gegen die absatzweise Inline-Code-Paarung, dabei ebenso
  fence-treu),
  verbotenes bzw. gefordertes Muster (`section-forbidden`,
  `section-pattern-missing`), geforderte Marken (`section-marker-missing`),
  Chronologie-Monotonie der Schlüsselspalte (`section-unordered`,
  `section-cell-untyped`), die Form **jeder** Überschrift des Abschnitts
  (`headings-match`/`headings-level` ⇒ `section-heading-mismatch`, gemeldet je
  Überschrift auf ihrer Zeile) und die **Zellenlänge** einer über ihren
  Kopfzeilen-**Namen** benannten Tabellenspalte (`table.column` ⇒
  `section-cell-oversized`/`section-cell-undersized` auf **ihrer** Zeile,
  `section-column-missing` für die nicht adressierbare Spalte);
  fehlt der Abschnitt oder trifft die Regel keine Datei ⇒ `section-missing`,
  mehrfach vorhanden bei `sections: one` ⇒ `section-ambiguous`. Jede Regel darf
  **ihre Grundmenge erklären** — `tasks-ignore-pattern` nimmt Task-Items aus der
  `max-tasks`-Zählung (geprüft wird der **Item-Text hinter der Checkbox**, damit
  `^` „das Item beginnt so" heißt), `exempt-section-pattern` nimmt **Abschnitte**
  aus der Regel (geprüft wird **dieselbe** rohe Überschriften-Zeile wie bei
  `section-pattern`, samt `#`-Folge). Beide **verkleinern** nur, tragen keinen
  eigenen Grund-Code, und ohne sie ist der Befundsatz byte-identisch; leert die
  Abschnitts-Ausnahme die Menge ⇒ `section-missing` statt stillem Grün —
  **es sei denn, die Regel sagt, wie viele sie ausnimmt:** `exempt-expect-count`
  (int ≥ 0, nur mit dem Muster) macht die erklärte Leermenge stumm, solange die
  Zahl stimmt, und meldet sonst `section-exempt-mismatch` — in **beide**
  Richtungen und **unabhängig** von einer Restmenge. Der Konfigurationsdefekt
  (der Selektor trifft nichts) bleibt `section-missing`, auch mit gesetzter
  Zahl. Die Zahl gilt **je Datei**, ihr Befund **bricht die Datei ab**, und sie
  altert wie jeder Autoren-Text. Und
  `hint` lässt eine Regel ihre **Erläuterung selbst verfassen** — was der Leser
  tun soll, weiß nur sie; der Grund-Code sagt die **Art** des Defekts.
  **Hermetisch**,
  opt-in; die Closure-Note-Struktur des Moduls `planning` ist ein **Preset**
  derselben Semantik
  ([`DC-FA-STRUCT-001`](spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in))
- `workflows` — Deklarations-Konsistenz der `uses:`-Referenzen von CI-Workflows
  unterhalb eines **konfigurierten** Verzeichnisses (`workflows.dir` — nicht
  verdrahtet, weil CI-System-spezifisch), opt-in. Eine **fremde** Referenz nennt
  einen vollen 40-stelligen Commit-SHA (`uses-pin-missing`) mit Tag-Kommentar
  dahinter (`uses-pin-untagged`) — ein Tag lässt sich umhängen, ein SHA nicht;
  geprüft wird die **Form**, nicht die Gültigkeit (das wäre Netz). Eine
  **lokale** Referenz (`./…`) braucht keinen Pin, dafür zwei andere Zusagen: das
  Ziel existiert (`uses-local-missing`), und der aufrufende **Job** führt die
  Rechte, die es verlangt (`uses-local-perms-undeclared`,
  `uses-local-perms-narrow`) — ein Job ohne eigenes `permissions:` erbt zwar den
  Workflow-Kopf, kann aber nichts weitergeben, was er nicht deklariert.
  Unlesbares YAML meldet (`workflow-unparsable`), statt übersprungen zu werden.
  **Hermetisch** (kein git, kein Netz, kein Ausführen); die Referenzen kommen aus
  dem **YAML-Baum**, nicht aus einer Textsuche
  ([`DC-FA-WF-001`](spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in))

Jeder Befund nennt Datei, Zeile, Ziel und Grund; Exit-Codes:
`0` sauber, `1` Befunde, `2` Umgebungs- oder Konfigurationsfehler.

## Warum d-check?

Zwölf funktional überlappende Tool-Kopien aus drei Familien
(`verify-doc-refs.sh` — Shell, `check_refs.py` — Python,
`docs-check.js` — JavaScript) sind in den Schwester-Repositories des
Entwicklungs-Workspace gewachsen — jede mit eigenem Funktionsstand,
eigener Drift, eigener Pflege. d-check ersetzt sie durch **ein** Tool:

- **Konfiguration statt Fork:** repo-spezifisches Verhalten lebt
  deklarativ in `.d-check.yml`
  ([`DC-FA-CONF-001`](spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)),
  nicht in kopierten Skripten.
- **Ersetzbarkeit ist messbar:** Jedes Alt-Tool muss durch d-check
  mit passender Konfiguration ablösbar sein — mindestens dieselben
  echten Befunde, keine False-Positives, die eine grüne CI brechen
  ([`DC-QA-04`](spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)).
- **Ein Distributionsweg:** Container-Image mit Digest-Pin statt
  n gepflegter Kopien
  ([`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image)).

## Kerngedanke

**Dokumentation ist ein Referenz-Graph mit prüfbaren Invarianten.**
Ob ein Link sein Ziel erreicht, ein Anker existiert, eine Kennung auf
ihre Definition verlinkt oder eine Dokumentklasse abwärts zeigt, ist
maschinell entscheidbar — d-check macht diese Invarianten zum Gate
statt zur Review-Meinung.

Dabei gilt: **berichten, nie reparieren.** d-check ist ein reines
Lese-Tool ([`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit));
deterministische Befunde werden behoben, nicht unterdrückt. Einen
Opt-out-Marker gibt es nur dort, wo ein nicht existierendes Ziel oder
eine illustrative Kennung dokumentierte Absicht sein kann
(`d-check:ignore`, zeilenweise) — er stellt vier Module still: `codepaths`,
`ids`, `versions` und `diagrams`. Das ist eine **benannte Liste**, kein
ableitbares Kriterium: `matrix`, `structure` und `citations` melden ebenfalls
auf Zeilen und kennen den Marker nicht. Bei `codepaths` und `ids` zählt er nur **außerhalb von Inline-Code** — in Backticks ist er eine Erwähnung, und er muss in einem HTML-Kommentar stehen; bei `versions` und `diagrams` bleibt er ein roh gelesenes **Token** — bei `diagrams` strukturell, bei `versions` als benannte Grenze (es liest auch Prosa-Zeilen)
([`DC-FA-CODE-001`](spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in),
[`DC-FA-ID-001`](spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids),
[`DC-FA-VER-001`](spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in),
[`DC-FA-DIAG-001`](spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in));
alle übrigen kennen ihn nicht.

## Was macht es vertrauenswürdig?

**Ein Prüf-Tool, dessen eigene Aussagen wackeln, ist wertlos** —
Determinismus und Seiteneffektfreiheit sind deshalb Kernverträge des
Lastenhefts, und beide werden gemessen, nicht behauptet:

- **Determinismus:** identische Eingabe ⇒ byte-identische Ausgabe,
  stabil sortiert; getestet über zehn wiederholte Läufe mit
  Hash-Vergleich
  ([`DC-QA-02`](spec/lastenheft.md#dc-qa-02--determinismus)).
- **Seiteneffektfrei und netzlos:** schreibt nie in das geprüfte
  Repository, öffnet außer in den opt-in-Modulen `external` und `sources`
  keine Netzwerkverbindungen — gemessen im Gate-Lauf mit read-only-Mount
  und `--network none`
  ([`DC-QA-03`](spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **Keine stillen Defaults:** Jede ungültige `.d-check.yml` bricht
  mit Exit 2 ab; geprüft wird nie mit geratener Konfiguration
  ([`DC-FA-CONF-001`](spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)).
- **Dogfooding:** d-check validiert die eigene Doku bei jedem
  Gate-Lauf — mit der [Selbstkonfiguration](.d-check.yml) im
  Vollausbau (**elf** Module inkl. Referenzmatrix, Span-Artefakten,
  Host-Pfad-Hygiene, Versions-Pin-Konsistenz, Struktur-Invarianten,
  Diagramm-Kennungen und wortgleichen Zitaten).
- **Container nativ-identisch:** Befund-Ausgabe und Exit-Code des
  Images sind byte-identisch zur nativen Ausführung, automatisiert
  getestet ([`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image));
  CI-Konsum per Digest-Pin.

## Nutzung

Verteilung als Container-Image über GHCR
([`DC-FA-DIST-001`](spec/lastenheft.md#dc-fa-dist-001--docker-image)); **seit `v0.67.0`**
zusätzlich nach Docker Hub gespiegelt als
`pt9912/d-check` — dasselbe Bild, kein zweiter Bau, gleicher **Config**-Digest
(der **Manifest**-Digest ist registry-lokal: per Digest pinnt man den der
Registry, aus der man zieht)
([`DC-FA-DIST-002`](spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel)):

```bash
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.71.1
```

CI-Pipelines pinnen auf den Digest aus den Release-Notes statt auf
bewegliche Tags. Der aufgabenorientierte Einstieg — von der ersten Prüfung bis
zu jedem einzelnen Modul — ist das
[Benutzerhandbuch](docs/user/benutzerhandbuch.md); die knappe Aufruf-Referenz
mit Optionen und Exit-Codes steht in
[`docs/user/operations.md`](docs/user/operations.md), der Release- und
Digest-Pin-Weg in [`docs/user/releasing.md`](docs/user/releasing.md).

## Konfiguration (`.d-check.yml`)

Optional in der Repo-Wurzel; ohne Datei laufen die Default-Module
`links` + `anchors` über `docs/`, `spec/` und die Root-`*.md`:

```yaml
scan:
  roots: ["."]                  # gesamte Repo-Wurzel
modules: [links, anchors, ids]  # external + sources bleiben strikt opt-in (Netz)
ids:
  patterns:
    - regex: 'ADR-\d{4}'
      target: docs/adr/         # Kennungen müssen hierhin verlinken
```

Das vollständige Schema mit allen Schlüsseln, Defaults und
Validierungs-Constraints steht in der
[Spezifikation §`.d-check.yml`](spec/spezifikation.md#spec-005--d-checkyml);
ein lebendes Beispiel im Vollausbau (inkl. Referenzmatrix) ist die
[Selbstkonfiguration dieses Repos](.d-check.yml). Jede ungültige
Konfiguration bricht mit Exit 2 ab — geprüft wird nie mit
stillschweigenden Defaults
([`DC-FA-CONF-001`](spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)).

## Einstieg

| Dokument | Inhalt |
|---|---|
| [`docs/user/benutzerhandbuch.md`](docs/user/benutzerhandbuch.md) | **Benutzerhandbuch**: aufgabenorientiert, jedes Regelmodul mit Beispiel-Konfiguration |
| [`docs/user/operations.md`](docs/user/operations.md) | Aufruf-Referenz: Optionen, Exit-Codes, Konfiguration |
| [`docs/user/releasing.md`](docs/user/releasing.md) | Release-Prozess, Digest-Pin-Konsum |
| [`spec/lastenheft.md`](spec/lastenheft.md) | Anforderungen (`DC-FA-*`, `DC-QA-*`), Akzeptanzkriterien |
| [`harness/README.md`](harness/README.md) | Harness-Einstieg: Source Precedence, Guides, Sensors |
| [`AGENTS.md`](AGENTS.md) | Briefing für AI-Coding-Agenten, Hard Rules |
| [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md) | Wellen und Slices |
| [`CHANGELOG.md`](CHANGELOG.md) | Änderungshistorie |

## Entwicklung

Der Host braucht nur `git`, GNU `make`, `bash` und Docker
(siehe [`AGENTS.md`](AGENTS.md) §3.1).

```bash
make help     # verfügbare Targets
make gates    # alle inneren Gates (mandatory vor Handoff)
```

## Lizenz

Dieses Projekt steht unter der [MIT-Lizenz](LICENSE).
