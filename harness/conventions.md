# Harness-Konventionen

## Purpose

Diese Datei deklariert die *repo-lokalen* Strukturregeln dieses Repos
gegenüber der adoptierten Harnesskonvention (Baseline). Sie ist der
Default-Ort für:

- **Adaptionen** ggü. der Baseline (mit Begründung und Auflösungs-Trigger).
- **ID-Schema-Deklaration** — welches Präfix-Schema dieses Repo nutzt.
  Der Baseline-Default wird als Teil der [`MR-000`](#mr-000--baseline-aussage)-Aussage
  festgehalten; ein abweichendes Präfix oder Schema ist ein eigener `MR`-Eintrag.
- **Zusatzklassen-Deklarationen** für repo-spezifische
  Bindung-Klassen in der Sensors-Tabelle, die über die vier kanonischen
  hinausgehen (ADR, Carveout, Schwelle, Reproduzierbarkeit).
- **Modus-Deklarationen** pro Sub-Area (Greenfield / Brownfield /
  Hybrid) inklusive Konvergenz-Auftrag bei BF.

Bei Konflikt zwischen dieser Datei und einer kanonischen Quelle gilt die
kanonische Quelle (Source Precedence, [`harness/README.md`](README.md#source-precedence)).
Diese Datei ist konformitätsbringend für *Form*-Fragen, nicht autoritativ
über Inhalt.

## Baseline

- **Konvention:** AI-Harness-Kurs
- **Stand:** [`v6.0.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v6.0.0),
  gepinnt mit
  [`MR-060`](#mr-060--baseline-pin-hebung-auf-v600-zehnter-nachtrag-zu-mr-011-nachtrag-zu-mr-023)
  — der jeweils aktuelle Eintrag der Pin-Serie. Die **Kette** der bisherigen
  Hebungen steht nicht hier, sondern in
  [§Aufgelöste Adaptionen](#aufgelöste-adaptionen): dort trägt jede Zeile ihren
  Nachfolger, und die Einträge selbst liegen in
  [`conventions/done/`](conventions/done/). Der Anfang der Serie ist
  [`MR-011`](#mr-011--baseline-auf-release-tag-gepinnt)
- **Datum der Adoption:** 2026-06-10

## Adoptierte Konventions-Quellen

- **Extern (Lehrmaterial):**
  [`ai-harness-course@v6.0.0`](https://github.com/pt9912/ai-harness-course/tree/v6.0.0)
  (`kurs/de/` — Konventionen in `grundlagen/`, Templates in `lab/templates/`).
  Kanonische Quelle; bei Konflikt maßgeblich.
- **Vendored Baseline (Regelwerk + Templates):** aus dem self-contained
  Release-Asset
  [`lab-regelwerk.zip`](https://github.com/pt9912/ai-harness-course/releases/download/v6.0.0/lab-regelwerk.zip)
  entpackt nach [`.harness/baseline/v6.0.0/`](../.harness/baseline/v6.0.0/regelwerk/)
  (`{regelwerk,templates}/` + `SHA256SUMS`) — der **netzlose** Lesepfad,
  materialisiert/verifiziert per
  [`fetch-baseline-cache.sh`](../tools/harness/fetch-baseline-cache.sh).
- **In-Repo (verkörperte Form):** die eigenen Artefakte dieses Repos (`AGENTS.md`,
  `harness/README.md`, ADRs, Slices, dieser Konventionsspeicher) — kopiert und
  ausgefüllt aus den vendored `templates/` (Referenz-/Ziel-Form).

## Adaptions-Block

Regeln dieser Sektion: Diese Datei trägt den **Index**, nicht die Einträge. Jede
Adaption ist eine eigene Datei unter `harness/conventions/`, kopiert aus der
vendored Vorlage `harness/conventions/MR-<NNN>-titel.template.md`; ist ihr
Auflösungs-Trigger eingetreten, wandert sie per `git mv` nach `conventions/done/`.
Der **Zustand ist die Verzeichnis-Position**, kein Status-Feld — was hier steht,
liest **jeder** Agentenlauf, aufgelöste Adaptionen gehören nicht in diesen Pfad
([Baseline-Regelwerk §Konventionsspeicher](../.harness/baseline/v6.0.0/regelwerk/grundlagen-harness-dateien.md#harnessconventionsmd-als-konventionsspeicher)).
Je Index-Zeile steht ein Voll-Slug-`<a id>` — **Migrations-Schuld**, damit die
eingefrorenen in-repo-Verweise auf `conventions.md#mr-…` (immutable ADRs, `done/`-
Slices, Reviews) ohne Retarget und ohne ADR-Edit auflösen; ein frisch unter der
Baseline aufgesetztes Repo braucht sie nicht.

### MR-000 — Baseline-Aussage

Bleibt hier: keine Adaption, sondern die Adoptions-Erklärung — sie gilt für jeden Lauf.

- **Status:** Accepted
- **Datum:** 2026-06-10
- **Geltungsbereich:** gesamtes Repo
- **Ersetzt-Baseline-Regel:** — *(keine; dieser Eintrag ist die Adoptions-Erklärung,
  keine Adaption)*
- **Adaption:** *keine inhaltlichen Adaptionen ggü. Baseline-Default für
  Verzeichniskonvention, Lifecycle-Regeln, Carveout-Disziplin, ID-Schema
  (`DC-FA-*`, `DC-QA-*`, `ADR-NNNN`, `CO-NNN`, `slice-NNN`, `MR-NNN`; Präfix `DC`).*
  **Vergabe** (deklariert mit dem v5.6.0-Bump, Baseline
  §Vergabe): **dichte, repo-weite Nummern** je Präfix — ein schreibender
  Mensch, kein Bereichssegment; die nächste Nummer liest Verzeichnis **und**
  offene Welle-Dateien. **Struktur-IDs** `SPEC-<NNN>` (`spec/spezifikation.md`)
  und `ARC-<NNN>` (`spec/architecture.md`) werden nach Baseline-Default
  vergeben (Auftraggeber-Entscheid 2026-08-22, welle-80): **fortlaufend je
  Datei**, Lücken werden nicht nachbelegt, kein Bereichssegment; der Link
  trägt den Abschnitt, der Text die Kennung; Struktur-IDs sind keine
  Anforderungen und gehören nicht in Commit-Botschaften. ADRs mit Status
  `Accepted` vor welle-80 adressieren weiter per `§`-Anker (immutabel) — zwei
  Formen, eine Regel. Den vorherigen Verzicht trug die aufgelöste
  Abweichung [`MR-027`](#mr-027).
- **Begründung:** Initial-Setzung. Spätere Adaptionen werden als `MR-<NNN>` nachgetragen.
- **Auflösungs-Trigger:** permanent.

### Aktive Adaptionen

Eine Zeile je Datei in `harness/conventions/`; Geltungsbereich und
Ersetzt-Baseline-Regel stehen hier, damit ein Agent ohne Öffnen entscheiden kann,
ob der Eintrag ihn betrifft.

| MR                                                                                                                                                | Titel                                         | Geltungsbereich                                   | Ersetzt-Baseline-Regel                                       |
| ------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------- | ------------------------------------------------- | ------------------------------------------------------------ |
| [MR-004](conventions/MR-004-gate-nachweis-mechanik.md) <a id="mr-004--gate-nachweis-mechanik-und-claude-hooks-nach-b-cad-vorbild"></a><a id="mr-004"></a>            | Gate-Nachweis-Mechanik + `.claude`-Hooks      | `tools/harness/`, `.claude/`, `make record-gates` | `grundlagen-durchsetzungsschicht` §Drei Bindepunkte          |
| [MR-005](conventions/MR-005-haertung-gate-nachweis.md) <a id="mr-005--härtung-ggü-b-cad-inhaltsbasierter-gate-nachweis-sub-shell-prüfung"></a><a id="mr-005"></a>    | Härtung: Content-Hash + Sub-Shell-Guard       | `working-tree-hash.sh`, `.claude/hooks/`          | `grundlagen-durchsetzungsschicht` §Vier Design-Eigenschaften |
| [MR-006](conventions/MR-006-referenzrichtung-matrix.md) <a id="mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs"></a><a id="mr-006"></a>         | Referenzrichtung/Matrix (+ C-4-Scope-Grenze)  | `spec/`-Straten, `matrix`-Config                  | `grundlagen-referenz-richtung` §SDP                          |
| [MR-007](conventions/MR-007-aufloesung-mr-003.md) <a id="mr-007--auflösung-von-mr-003-doc-check-als-dogfooding"></a><a id="mr-007"></a>                              | doc-check als Dogfooding                      | `make doc-check`, `.d-check.yml`                  | `modul-13` §Hard Rule (Doku-Disziplin)                       |
| [MR-013](conventions/MR-013-lifecycle-move-buendelung.md) <a id="mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise"></a><a id="mr-013"></a>                  | Lifecycle-Move-Commit bündelt Verweise        | Slice-Lifecycle-Moves (**beide** Richtungen: `open/`→`in-progress/`, `in-progress/`→`done/`), MR-/Wellen-Moves, `make planning-check` | `modul-05` §Lifecycle als State Machine                      |
| [MR-015](conventions/MR-015-agents-md-routet.md) <a id="mr-015--auflösung-der-mr-012-pointer-drift-agentsmd-routet-spiegelt-nicht-mehr"></a><a id="mr-015"></a>      | AGENTS.md routet (spiegelt nicht)             | `AGENTS.md` §1                                    | `grundlagen-harness-dateien` §Template-Schichtung            |
| [MR-021](conventions/MR-021-vendored-verweise-pin-gebunden.md) <a id="mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden"></a><a id="mr-021"></a> | vendored-Verweise pin-gebunden                | Live-Links auf die vendored Baseline              | `grundlagen-harness-dateien` §Verzeichniskonvention          |
| [MR-023](conventions/MR-023-baseline-v500.md) <a id="mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout"></a><a id="mr-023"></a>                 | self-contained Bundle-Layout (vendored; historische Pin-Hebung auf v5.0.0 — den aktuellen Pin trägt der Nachtrag darunter) | §Baseline, `fetch-baseline-cache.sh`              | `grundlagen-harness-dateien` §Template-Schichtung            |
| [MR-025](conventions/MR-025-spiegel-vor-dem-editieren.md) <a id="mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten"></a><a id="mr-025"></a> | Semantik-Änderung: Spiegel **vor** dem Editieren auflisten | jede Änderung an einer zugesagten Semantik (Grund-Code, Algorithmus-Schritt, Config-Schlüssel, Schwelle, Erkennungs-Form) | `modul-10` §Review-Arten |
| [MR-031](conventions/MR-031-schritt-3-benennen.md) <a id="mr-031--schritt-3-des-agenten-workflows-verlangt-benennen-nicht-nur-identifizieren"></a><a id="mr-031"></a> | Schritt 3 verlangt Benennen statt nur Identifizieren | `AGENTS.md` §6 Schritt 3; jeder Implementer-Lauf | `modul-09-implementierung.md` §Minimal Agent Workflow, Schritt 3 |
| [MR-032](conventions/MR-032-historie-vor-accepted.md) <a id="mr-032--versions-bump-und-historie-zeile-schon-vor-accepted"></a><a id="mr-032"></a> | Bump und Historie schon vor `Accepted` | `spec/lastenheft.md`, solange sein Status unter `Accepted` liegt | `grundlagen-source-precedence.md` §Wann die CR-Pflicht beginnt |
| [MR-033](conventions/MR-033-sicht-ohne-modul-pfade.md) <a id="mr-033--die-architektur-sicht-führt-auch-keine-modul-pfade"></a><a id="mr-033"></a> | Sicht ohne Modul-Pfade | `spec/architecture.md` und `AGENTS.md` §3.4 | `AGENTS.template.md` §3.4 / `modul-03-spec.md` §Ziel-Form: Architektur-Sicht |
| [MR-034](conventions/MR-034-matrix-scope-welle.md) <a id="mr-034--die-referenzmatrix-bewacht-auch-die-kante-adr--welle"></a><a id="mr-034"></a> | Matrix bewacht ADR → Welle | `matrix`-Block der `.d-check.yml` | schärft [`MR-006`](conventions.md#mr-006) §Scope-Grenze (C-4) |
| [MR-035](conventions/MR-035-cr-ablage.md) <a id="mr-035--ausgehende-change-requests-an-die-baseline-liegen-im-repo"></a><a id="mr-035"></a> | Ausgehende CRs liegen in `docs/plan/cr/` | `docs/plan/cr/` | keine — Form-Frage, vom Kanon an den Konventionsspeicher abgetreten |
| [MR-036](conventions/MR-036-cr-antwort-ablage.md) <a id="mr-036--die-antwort-auf-einen-ausgehenden-cr-liegt-bei-ihrem-cr"></a><a id="mr-036"></a> | Die Antwort liegt bei ihrem CR | `docs/plan/cr/` — eingehende Antworten | keine — der Kanon kennt den ausgehenden CR nicht, also auch nicht seine Antwort |
| [MR-039](conventions/MR-039-zitat-delta-im-neuen-eintrag.md) <a id="mr-039--ändert-die-baseline-einen-zitierten-wortlaut-hält-das-ein-neuer-eintrag-fest--nicht-der-zitierende"></a><a id="mr-039"></a> | Geändertes Baseline-Zitat wird im Bump-Eintrag vermerkt, nicht im zitierenden Dokument | wörtliche Baseline-Zitate in allen lebenden Dokumenten | keine — der Kanon gibt das Prinzip, nicht den Ort |
| [MR-040](conventions/MR-040-guard-skript-interpreter.md) <a id="mr-040--der-tool-call-wächter-blockiert-auch-host-skript-interpreter-schärft-mr-005"></a><a id="mr-040"></a> | Wächter blockiert Host-Skript-Interpreter (schärft [`MR-005`](conventions.md#mr-005)) | `.claude/hooks/pretooluse-command-guard.sh`, `AGENTS.md` §3.1 | keine — der Kanon verlangt die Härtung, nicht ihre Liste |
| [MR-042](conventions/MR-042-guard-in-eigener-klasse.md) <a id="mr-042--der-wächter-läuft-in-der-klasse-die-er-durchsetzt-löst-mr-041-auf"></a><a id="mr-042"></a> | Wächter in `bash` + `awk`; `node`-Sperre fortgeführt (löst [`MR-041`](conventions.md#mr-041) auf) | `.claude/hooks/pretooluse-command-guard.sh`, `tools/harness/extract-command.awk`, `tools/harness/guard-probe.sh`, `AGENTS.md` §3.1/§4 | keine — Form-Frage, vom Kanon an den Konventionsspeicher abgetreten |
| [MR-043](conventions/MR-043-einstieg-importiert-agents.md) <a id="mr-043--der-werkzeug-einstieg-importiert-agentsmd-statt-auf-ihn-zu-verweisen"></a><a id="mr-043"></a> | `CLAUDE.md` importiert `AGENTS.md` (`@`-Syntax), statt darauf zu verweisen | `CLAUDE.md` | keine — Form-Frage, vom Kanon an den Konventionsspeicher abgetreten |
| [MR-044](conventions/MR-044-guard-zwei-kanaele.md) <a id="mr-044--der-wächter-blockt-über-zwei-kanäle-und-die-antwortform-ist-die-aktuelle-schärft-mr-042"></a><a id="mr-044"></a> | Wächter blockt über zwei Kanäle; Antwortform aktuell (schärft [`MR-042`](conventions.md#mr-042)) | `.claude/hooks/pretooluse-command-guard.sh`, `tools/harness/guard-probe.sh` | keine — Form-Frage, vom Kanon an den Konventionsspeicher abgetreten |
| [MR-045](conventions/MR-045-slice-verweise-nicht-im-briefing.md) <a id="mr-045--agentsmd-und-harnessreadmemd-tragen-keine-slice-verweise-schärft-mr-013"></a><a id="mr-045"></a> | Briefing und Harness-Einstieg ohne Slice-Verweise (schärft [`MR-013`](conventions.md#mr-013)) | `AGENTS.md`, `harness/README.md` | keine — folgt aus der Rollenzuweisung, steht dort aber nicht ausdrücklich |
| [MR-046](conventions/MR-046-keine-skript-stage.md) <a id="mr-046--es-gibt-keine-skript-stage-weil-kein-fall-sie-verlangt-schärft-mr-040"></a><a id="mr-046"></a> | Keine Skript-Stage; §3.1 verweist nicht mehr auf eine (schärft [`MR-040`](conventions.md#mr-040)) | `AGENTS.md` §3.1, `Dockerfile` | keine — eine Regel zeigt auf Vorhandenes, nicht auf Vorgesehenes |
| [MR-047](conventions/MR-047-permission-schicht.md) <a id="mr-047--eine-zweite-durchsetzungs-schicht-neben-dem-wächter-schärft-mr-040"></a><a id="mr-047"></a> | Permission-Sperrliste als zweite Schicht (schärft [`MR-040`](conventions.md#mr-040)) | `.claude/settings.json` | keine — der Kanon verlangt die Härtung, nicht ihre Schichtenzahl |
| [MR-048](conventions/MR-048-gate-ueber-werkzeug-datei.md) <a id="mr-048--ein-repo-gate-über-eine-werkzeug-datei-prüft-wohlgeformtheit-nicht-anwesenheit-schärft-mr-047"></a><a id="mr-048"></a> | Gate über eine Werkzeug-Datei: Wohlgeformtheit ja, Anwesenheit nein (schärft [`MR-047`](conventions.md#mr-047)) | jede Repo-Prüfung mit Gegenstand unter `.claude/` | keine — der Kanon führt die Werkzeug-Artefakte, aber nicht als Bindepunkt |
| [MR-049](conventions/MR-049-ausgangs-wortschatz.md) <a id="mr-049--der-ausgang-eines-risikos-trägt-eines-der-drei-wörter"></a><a id="mr-049"></a> | Geschlossener Ausgangs-Wortschatz, gewächtert im Closure-Profil | `.d-check.closure.yml`, §5 der `done/`-Slices | keine — Werkzeug-Wahl für die urteilsfreie Hälfte, die der Kanon dem Repo überlässt |
| [MR-050](conventions/MR-050-herkunfts-anker-ist-kein-verweis.md) <a id="mr-050--der-herkunfts-anker-einer-hard-rule-ist-kein-slice-verweis-schärft-mr-045"></a><a id="mr-050"></a> | Herkunfts-Anker `(seit slice-<NNN>)` ist von [`MR-045`](conventions.md#mr-045) ausgenommen (schärft ihn) | `AGENTS.md`, `harness/README.md` | keine — stellt die Baseline-Vorgabe wieder her, die der Vorgänger zu breit gefasst hatte |
| [MR-051](conventions/MR-051-cite-spannen-beim-bump.md) <a id="mr-051--ein-baseline-bump-ankert-die-d-checkcite-spannen-neu-nachtrag-zu-mr-021"></a><a id="mr-051"></a> | Der Bump-Drift-Audit zieht nicht nur den **Pfad** einer `cite`-Direktive auf den neuen Tag, sondern ankert auch ihre **Zeilenspanne** neu; steht der Wortlaut gar nicht mehr da, greift [`MR-039`](conventions.md#mr-039) und die Direktive entfällt | `d-check:cite`-Direktiven auf `.harness/baseline/<tag>/` in lebenden Dokumenten — die Feld-Aussage, in den eingefrorenen Verzeichnissen stehe keine, ist durch [MR-054](#mr-054) überholt | keine — Nachtrag zu [`MR-021`](conventions.md#mr-021); der Kanon kennt keine Zeilen-Anker, weil er kein Werkzeug voraussetzt, das sie prüft |
| [MR-052](conventions/MR-052-historisches-zitat-ist-pruefbar.md) <a id="mr-052--ein-historisches-baseline-zitat-ist-prüfbar-der-alte-baum-liegt-in-der-git-historie-schärft-mr-039"></a><a id="mr-052"></a> | Der Bump entfernt den alten Baseline-Baum aus dem **Arbeitsbaum**, nicht aus dem **Repo** — ein historisches Zitat ist damit **prüfbar**, und „am Repo nicht entscheidbar" ist keine zulässige Antwort mehr. Der Wortlaut bleibt trotzdem eingefroren; gemessen ist nur, **welche** Wiedergabe wörtlich war | wörtliche Baseline-Zitate auf einen **früheren** Pin in lebenden Dokumenten | keine — Schärfung von [`MR-039`](conventions.md#mr-039); der Kanon stellt die Wörtlichkeits-Frage nicht |
| [MR-053](conventions/MR-053-dritte-vorpruefung-nachtlauf.md) <a id="mr-053--die-slice-planung-trägt-eine-dritte-vorprüfung-den-nachtlauf-stand"></a><a id="mr-053"></a> | Dritter `**Vorgelagert**`-Block im Slice-Plan: den **Nachtlauf-Stand** lesen (`make nightly-state`). Adressat ist der Implementer beim Planen des nächsten Slice, Takt ist jede Slice-Planung — ein vorhandener Lese-Moment statt eines neuen. **Kein Benachrichtigungs-Kanal**, mit Grund; **benannte Grenze:** in einer Pause liest niemand | jeder neue Slice-Plan, §Vorgelagert (kein Retrofit) | ergänzt die **zwei** Blöcke der Baseline-`slice.template.md` — der Kanon kennt keinen Nachtlauf |
| [MR-054](conventions/MR-054-vorpruefungen-belegen-ihre-regel.md) <a id="mr-054--slice-vorprüfungen-belegen-ihre-regel-mit-d-checkcite"></a><a id="mr-054"></a> | Die beiden **kanonischen** Vorprüfungs-Blöcke jedes neuen Slice-Plans (Sub-Area-Wahl, Beobachtungs-Sichtung) tragen eine `d-check:cite`-Direktive auf die Regelwerk-Zeilen, die sie vorschreiben, mit dem wörtlichen Zitat darunter — und zwar auf die **vorschreibende** Zeile, nicht auf eine Nebenregel aus demselben Absatz — `citations` prüft es wortgleich, fail-closed, im inneren Loop. **Anlass gemessen:** drei Slices liefen mit vollständig ausgefüllten Blöcken durch, während der zuständige Zyklus-Abschnitt ungelesen war; kein Block war falsch, sie waren alle nur Deklaration. Der **dritte** Block (Nachtlauf, [MR-053](#mr-053)) trägt bewusst keine — sein Ziel ist repo-eigen und meldete bei jeder Änderung. **Nachtrag zu [MR-051](#mr-051):** dessen Geltungsbereich-Aussage über die eingefrorenen Verzeichnisse wird absehbar unwahr, sobald ein Slice mit Direktiven schließt; `citations.scope` nimmt die drei eingefrorenen Verzeichnisse aus — der Beleg zählt zum Zeitpunkt seiner Prüfung, danach ist er Lauf-Beleg. **Grenze:** ein Zitat belegt Zugriff, nicht Verständnis |
| [MR-055](conventions/MR-055-symlink-als-pin-traeger.md) <a id="mr-055--ein-symlink-auf-den-vendored-baum-ist-ein-pin-gebundener-träger-nachtrag-zu-mr-021"></a><a id="mr-055"></a> | Ein **Symlink** auf `.harness/baseline/<tag>/` bindet denselben Pin wie ein Markdown-Link — und ist von [MR-021](#mr-021)s Zensus **nicht** erfasst, der nach Markdown-Links sucht. **Die Lücke ist gemessen:** der Scanner folgt Symlinks nicht in die Prüfmenge (eine *echte* Datei unter `.claude/rules/` wird gescannt, ein *Symlink* nicht), und `sha256sum -c` samt Manifest-Deckung bleibt grün, während der Alias ins Leere zeigt. **Die Antwort ist ein Sensor, keine Prozedur-Zeile:** `--verify` prüft als dritte Frage, dass jeder Symlink unter `.claude/rules/` auflöst, und läuft in `make gates`. Rekursiv und dotfile-bewusst; die Proben fahren als `make baseline-probe`. **Sechs Grenzen, ausgeschrieben:** geprüft wird die **Auflösung**, nicht das Ziel; ein **fehlendes** oder leeres `.claude/rules/` meldet nicht (wer die Aliase löscht statt sie umzuhängen, hat einen grünen Lauf); ein Alias auf ein **Verzeichnis** passiert; die dritte Frage läuft **nach** den beiden ersten und akkumuliert nicht; ein Symlink überlebt nicht jedes Dateisystem (`core.symlinks=false` macht Textdateien daraus); und ob die verlinkten Dateien wirklich in den Kontext laden, ist beobachtet, nicht gewächtert | Symlinks unter `.claude/rules/` mit Ziel in `.harness/baseline/<tag>/` | keine — Nachtrag zu [MR-021](#mr-021), **nicht** zum Kanon: dass In-Repo-Verweise pin-gebunden sind, ist die Adaption von [MR-021](#mr-021). Der Kanon kennt den Baum als gepinnte Referenz, aber keinen Symlink als Träger |
| [MR-056](conventions/MR-056-dod-haken-waechter.md) <a id="mr-056--der-dod-haken-eines-geschlossenen-slice-wird-gewächtert"></a><a id="mr-056"></a> | Der **DoD-Haken** eines geschlossenen Slice wird gewächtert: `max-open-tasks: 0` über den DoD-Abschnitt, mit `hint` für die Reparatur und einer Bestands-Ausnahme mit fester Ziffernzahl. **Nicht** `forbid-pattern` — die bereinigt lesende Form fiel an einem einzelnen Backtick auf null Befunde und deckte nur eine Bullet-Form. **Sie prüft den Zustand am Ruheort, nicht den Übergang**, und ein vergessener Schluss-Fence schaltet sie ab (gefangen von `spans` im selben Profil) | `.d-check.closure.yml`, Abschnitt `## N. Definition of Done` der Slice-Dateien unter `docs/plan/planning/done/` | keine — der Kanon macht die DoD-Häkchen zur Bedingung des Übergangs; gehalten war davon nur die Closure-Notiz-Hälfte. Werkzeug-Wahl, keine Abweichung |
| [MR-059](conventions/MR-059-wellen-archiv-stub-move.md) <a id="mr-059--wellen-archiv-stub-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013"></a><a id="mr-059"></a> | Ein Wellen-Archivierungs-Commit (`tools/archive-wave`) ersetzt Volltext durch Stub im selben Akt — anders als die drei [MR-013](#mr-013)-Fälle gibt es keine Phase, in der die bewegte Datei unverändert bleibt, also keine Zwei-Commit-Zerlegung. Bleibt ein einziger, in der Botschaft ausdrücklich deklarierter Commit; git zeigt reine D/A-Paare, keine Renames | `tools/archive-wave` (`Apply()`), jeder `make archive-wave WELLE=<id> APPLY=1`-Commit | keine — Nachtrag zu [MR-013](#mr-013); der Kanon kennt Modul 6 Schritt 4 als Operation, aber nicht ihre Commit-Granularität |
| [MR-060](conventions/MR-060-baseline-v600.md) <a id="mr-060--baseline-pin-hebung-auf-v600-zehnter-nachtrag-zu-mr-011-nachtrag-zu-mr-023"></a><a id="mr-060"></a> | Baseline-Pin-Hebung auf v6.0.0 (Pin-Fortschreibung, Nachtrag); wellenlose Zeitdokumente-Archivierung adoptiert (Regelwerk-Ebene), Beobachtungs-Register-Neugestaltung adoptiert (slice-194/195, selbe Welle) | §Baseline, pin-gebundene Verweise, `.harness/baseline/v6.0.0/` | — *(Pin-Fortschreibung im Bundle-Layout des Vorgängers)* |
| [MR-061](conventions/MR-061-register-migrations-move.md) <a id="mr-061--register-formatmigration-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013"></a><a id="mr-061"></a> | Register-Formatmigration (Beobachtungs-Register Tabelle → Verzeichnis) ist ein einziger, deklarierter Commit — korrigiert eine fehlgeschlagene [MR-059](#mr-059)-Analogie in der Commit-Botschaft von slice-195 | genau der eine Migrations-Commit von [slice-195](../docs/plan/planning/done/welle-88/slice-195-beobachtungsregister-migration.md); erschöpft, keine Blankovollmacht | keine — Nachtrag zu [MR-013](#mr-013) |
| [MR-062](conventions/MR-062-wellenloser-slice-archiv-move.md) <a id="mr-062--wellenloser-einzel-slice-archiv-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013"></a><a id="mr-062"></a> | Wellenloser Einzel-Slice-Archiv-Move (`tools/archive-wave -slice=<id>`) ist ein einziger, deklarierter Commit — korrigiert eine fehlgeschlagene [MR-059](#mr-059)-Analogie im Plan von slice-197 | `tools/archive-wave`s Einzel-Slice-Modus (`ApplySlice()`), jeder `-slice=<id> -apply`-Commit | keine — Nachtrag zu [MR-013](#mr-013) |

### Aufgelöste Adaptionen

Eine Zeile je Datei in `conventions/done/` — nur ID und Nachfolger, damit die Kette
auffindbar bleibt, ohne gelesen zu werden.

| MR                                                                                                                                                   | aufgelöst durch                                            |
| ---------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------- |
| [MR-001](conventions/done/MR-001-eigene-spec-schicht.md) <a id="mr-001--source-precedence-mit-eigener-spezifikations-schicht"></a><a id="mr-001"></a>                   | Baseline-Stand `v4.0.0` (drei Spec-Straten Default)        |
| [MR-002](conventions/done/MR-002-id-schema-bereichskuerzel.md) <a id="mr-002--id-schema-mit-bereichskürzeln-ab-initialer-fassung"></a><a id="mr-002"></a>               | [MR-000](#mr-000--baseline-aussage) (ID-Schema ist Teil der Baseline-Aussage)              |
| [MR-003](conventions/done/MR-003-vendorter-bootstrap-sensor.md) <a id="mr-003--vendorter-bootstrap-sensor-toolsverify-doc-refssh"></a><a id="mr-003"></a>               | [MR-007](conventions/MR-007-aufloesung-mr-003.md)          |
| [MR-008](conventions/done/MR-008-id-schema-deklaration.md) <a id="mr-008--id-schema-deklaration-nachtrag-zur-baseline-aussage"></a><a id="mr-008"></a>                  | [MR-000](#mr-000--baseline-aussage) (ID-Schema ist Teil der Baseline-Aussage)              |
| [MR-009](conventions/done/MR-009-source-precedence-ohne-docs-user.md) <a id="mr-009--source-precedence-ohne-docsuser-rang"></a><a id="mr-009"></a>                      | [MR-010](conventions/done/MR-010-aufloesung-mr-009.md)     |
| [MR-010](conventions/done/MR-010-aufloesung-mr-009.md) <a id="mr-010--auflösung-von-mr-009-docsuser-rang-eingefügt"></a><a id="mr-010"></a>                             | Baseline-Konformität (docs/user-Rang = Baseline-Default; kein Divergenz)  |
| [MR-011](conventions/done/MR-011-baseline-pin-release-tag.md) <a id="mr-011--baseline-auf-release-tag-gepinnt"></a><a id="mr-011"></a>                                  | [MR-012](conventions/done/MR-012-baseline-pin-hebung.md)   |
| [MR-012](conventions/done/MR-012-baseline-pin-hebung.md) <a id="mr-012--baseline-pin-hebung-nachtrag-zu-mr-011"></a><a id="mr-012"></a>                                 | [MR-016](conventions/done/MR-016-baseline-pin-hebung-2.md) |
| [MR-014](conventions/done/MR-014-slice-adr-haus-stil.md) <a id="mr-014--slice-adr-doc-struktur-repo-haus-stil-ggü-baseline-template"></a><a id="mr-014"></a>            | Baseline-Konformität (Form ist Baseline-Wahl; ADR-Alternativen-Tabelle seit `v1.4.0` Default)       |
| [MR-016](conventions/done/MR-016-baseline-pin-hebung-2.md) <a id="mr-016--baseline-pin-hebung-zweiter-nachtrag-zu-mr-011"></a><a id="mr-016"></a>                       | [MR-023](conventions/MR-023-baseline-v500.md)              |
| [MR-017](conventions/done/MR-017-cache-selbst-scan.md) <a id="mr-017--lokale-baseline-lese-form-cache-aus-dem-selbst-scan-ausgenommen"></a><a id="mr-017"></a>          | [MR-019](conventions/done/MR-019-regelwerk-vendored.md)    |
| [MR-018](conventions/done/MR-018-keine-templates.md) <a id="mr-018--d-check-verkörpert-als-producer-self-hoster-keine-templates"></a><a id="mr-018"></a>                | Baseline-Stand `v5.0.0` (Bundle vendored beide Bäume)      |
| [MR-019](conventions/done/MR-019-regelwerk-vendored.md) <a id="mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017"></a><a id="mr-019"></a>          | Baseline-Stand `v5.0.0` (Vendoring = Default)              |
| [MR-020](conventions/done/MR-020-template-propagation.md) <a id="mr-020--baseline-template-propagation-per-drift-audit-template-frei-bestätigt"></a><a id="mr-020"></a> | Baseline-Stand `v5.0.0` (Template-Schichtung = Default)    |
| [MR-022](conventions/done/MR-022-currency-audit.md) <a id="mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019"></a><a id="mr-022"></a>                            | Baseline-Stand `v5.0.0` (Freshness-Audit = Default)        |
| [MR-024](conventions/done/MR-024-aktuelle-welle-ruhe-marker-form.md) <a id="mr-024--aktuelle-welle-ruhe-marker-im-wellenlosen-zustand-aktive-welle-template-konform"></a><a id="mr-024"></a> | Baseline-Stand `v5.6.0` (Offene-Wellen-Form + Ruhe-Marker mit Wächter = Default; adoptiert mit slice-108) |
| [MR-026](conventions/done/MR-026-baseline-v560.md) <a id="mr-026--baseline-pin-hebung-auf-v560-dritter-nachtrag-zu-mr-011-nachtrag-zu-mr-023"></a><a id="mr-026"></a> | [MR-028](conventions/done/MR-028-baseline-v570.md) (nächste Pin-Hebung — der eigene Auflösungs-Trigger des Eintrags) |
| [MR-027](conventions/done/MR-027-struktur-id-verzicht.md) <a id="mr-027--struktur-ids-spec-arc--werden-nicht-vergeben"></a><a id="mr-027"></a> | Baseline-Konformität (Struktur-ID-Vergabe `SPEC-*`/`ARC-*` = Baseline-Default; Auftraggeber-Entscheid 2026-08-22, welle-80 — die Vergabe-Aussage trägt [MR-000](#mr-000--baseline-aussage)) |
| [MR-028](conventions/done/MR-028-baseline-v570.md) <a id="mr-028--baseline-pin-hebung-auf-v570-vierter-nachtrag-zu-mr-011-nachtrag-zu-mr-023"></a><a id="mr-028"></a> | [MR-029](#mr-029) (nächste Pin-Hebung — der eigene Auflösungs-Trigger des Eintrags) |
| [MR-029](conventions/done/MR-029-baseline-v590.md) <a id="mr-029--baseline-pin-hebung-auf-v590-fünfter-nachtrag-zu-mr-011-nachtrag-zu-mr-023"></a><a id="mr-029"></a> | [MR-030](#mr-030) (nächste Pin-Hebung — der eigene Auflösungs-Trigger des Eintrags) |
| [MR-030](conventions/done/MR-030-baseline-v5110.md) <a id="mr-030--baseline-pin-hebung-auf-v5110-sechster-nachtrag-zu-mr-011-nachtrag-zu-mr-023"></a><a id="mr-030"></a> | [MR-037](#mr-037) (nächste Pin-Hebung — der eigene Auflösungs-Trigger des Eintrags) |
| [MR-037](conventions/done/MR-037-baseline-v5120.md) <a id="mr-037--baseline-pin-hebung-auf-v5120-siebter-nachtrag-zu-mr-011-nachtrag-zu-mr-023"></a><a id="mr-037"></a> | [MR-057](#mr-057) (nächste Pin-Hebung — der eigene Auflösungs-Trigger des Eintrags) |
| [MR-038](conventions/done/MR-038-zitate-pin-gebunden.md) <a id="mr-038--ein-zitat-der-baseline-ist-pin-gebunden-wie-ein-link-aber-es-wird-ergänzt-statt-ersetzt-schärft-mr-021"></a><a id="mr-038"></a> | [MR-039](#mr-039) (der Kanon regelt den Fall doch — bestehende Einträge werden nicht rückwirkend umgeschrieben) |
| [MR-041](conventions/done/MR-041-guard-node-und-eigene-toolchain.md) <a id="mr-041--auch-node-ist-ein-host-interpreter-der-wächter-selbst-bleibt-die-benannte-ausnahme-schärft-mr-040"></a><a id="mr-041"></a> | [MR-042](#mr-042) (die benannte Inkonsistenz ist eingelöst; die `node`-Sperre trägt der Nachfolger fort) |
| [MR-057](conventions/done/MR-057-baseline-v5150.md) <a id="mr-057--baseline-pin-hebung-auf-v5150-achter-nachtrag-zu-mr-011-nachtrag-zu-mr-023"></a><a id="mr-057"></a> | [MR-058](#mr-058) (nächste Pin-Hebung, löst zugleich die dort gemeldete [MR-013](#mr-013)-Kollision auf) |
| [MR-058](conventions/done/MR-058-baseline-v5180.md) <a id="mr-058--baseline-pin-hebung-auf-v5180-neunter-nachtrag-zu-mr-011-nachtrag-zu-mr-023"></a><a id="mr-058"></a> | [MR-060](#mr-060) (nächste Pin-Hebung — der eigene Auflösungs-Trigger des Eintrags) |



## Zusatzklassen-Deklaration für Sensors-Bindung

Zusätzlich zu den vier kanonischen Klassen (ADR, Carveout, Schwelle,
Reproduzierbarkeit):

| Klasse     | Form   | Bedeutung                                       | Beispiel                                                                                                                                                                                                                                                          |
| ---------- | ------ | ----------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| DC-Bindung | `DC-…` | Gate prüft eine konkrete Lastenheft-Anforderung | [`DC-QA-02`](../spec/lastenheft.md#dc-qa-02--determinismus) für den Determinismus-Test in `make test`; [`DC-QA-03`](../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) für `make arch-check` und das Netzlos-Gate in `make doc-check` |

## Modus-Deklaration pro Sub-Area

| Sub-Area (Pfad / Modul)         | Modus      | BEO-Kürzel | Begründung                                                                                                   | Graduation-Bedingung / Folge-Slice |
| ------------------------------- | ---------- | ---------- | ------------------------------------------------------------------------------------------------------------ | ---------------------------------- |
| `*` (Default für gesamtes Repo) | Greenfield | `ALL`      | Projekt startet spec-first; Doc führt, Code folgt                                                            | n/a (GF)                           |
| `tools/harness/`                | Greenfield | `HARN`     | adoptierte Harness-Mechanik, konventionsgetragen über [MR-004](conventions/MR-004-gate-nachweis-mechanik.md) | n/a (GF)                           |
