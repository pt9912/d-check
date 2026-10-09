# Verifikation — slice-260: Festlegungen der lokalen Wächter, Hooks und Prüfer in die Spezifikation

- **Rolle:** Verifier (Modul 11). Die Frage ist: „Bauen wir es richtig?" Maßstab sind DoD, Spec und Plan.
- **Gegenstand:** `slice-260` · Range `b6a1c21e~1..HEAD` (18 Commits, Spitze `4db45aa3`). Code-tragend
  sind `ffedcb95` (Übergangs-Erkennung), `5bef299c` (Nachweis im Rezept), `e4505c7f` (§7, Sensor-Dateien,
  `MR-076`), `596630db` (R1-Einarbeitung) und `4db45aa3` (R2-Einarbeitung).
- **Datum:** 2026-10-09 · **Modell-ID:** claude-opus-5-5
- **Eingangs-Kontext:** Slice-Plan (§1–§3 mit fünf Plan-Änderungen, §8), Reports R1
  (`2026-10-08-…-r1.md`) und R2 (`2026-10-09-…-r2.md`), die Commit-Botschaften als Sensor-Belege des
  Implementers, `spec/spezifikation.md` §7/§8, `slice-262` und `slice-263` in `open/`.
- **Vorbehalt:** §7 (Closure-Notiz) ist leer, DoD-Haken stehen offen. Das ist der Stand vor der Closure.
  DoD 4 ist deshalb Closure-Schuld und keine DoD-Verletzung.
- **Proben:** Alle liefen in Wegwerf-Kopien unter dem Scratchpad (`vf/`): Klon, Modell-Makefile,
  Wegwerf-Repo. Die zwei echten `make`-Läufe mit rotem bzw. ignoriertem Glied liefen im Repo selbst. Den
  Nachweis `.harness/state/gates-passed.diffsha` habe ich vorher gesichert und danach byte-gleich
  zurückgelegt; der Hash stimmt mit dem Arbeitsbaum überein. `git status` war vor und nach allen Proben
  leer.

---

## DoD-Prüfung

### DoD 1 — Je Werkzeug mit eigener Festlegung ein §7-Eintrag, am Code geprüft — **BESTÄTIGT**

Ich habe jede Aussage am Code gelesen und jede Liste nachgezählt.

**`SPEC-093` (Tool-Call-Wächter)** gegen `.claude/hooks/pretooluse-command-guard.sh`,
`tools/harness/extract-command.awk` und `.claude/settings.json`:

- Sperrliste: `BLOCKED` hat 26 Wörter (gezählt), dazu das Muster `^python[0-9]*(\.[0-9]+)*$`.
- Präfixe: 8 (`sudo env command exec nice time xargs eval`), dazu `{` `}` und Zuweisungen.
- Shells: 5. Trenner: `&&`, `&`, `||`, `|`, `;`, `$(`, Backtick, `(`, Wagenrücklauf und Zeilenende.
- Der Pfad zählt mit dem letzten Segment (`${head##*/}`).
- Rekursion: `depth -gt 3` blockt.
- Fail-closed: Extraktor-Exit 3, fehlendes `awk`, fehlender Extraktor. Die Zwei-Kanal-Antwort ist
  `deny` plus Exit 2, der Durchlass gibt nichts aus und endet mit 0.
- `settings.json` `deny`: dieselben 26 Wörter plus `python`, `python3` (gezählt), je als `Bash(<wort> *)`.
- Eigene Proben (`vf/g.sh`, 78 Fälle, 78 ok):
  - Jeder Trenner und jedes Präfix vor `pip` blockt, ebenso alle 26 Wörter und `python`/`python2`/`python3.12`.
  - Alle fünf Shells mit `-c`, dazu `-lc`, blocken.
  - Tiefe: drei verschachtelte `bash -c` lassen durch, vier blocken.
  - Die genannten Durchlass-Klassen lassen durch: `then`, `do`, `!`, `timeout`, `nohup`, `env -i`,
    `sudo -u`, `find -exec`, `p"i"p`, escapte Quotes, `php`.
  - JSON-Zweifel (kein Objekt, abgeschnitten, `\u`, zwei Strings ohne Trenner, `@` außerhalb) blockt.
  - `nul`/`tru` passieren: Es wird das Zeichen geprüft, nicht das Literal, wie die Festlegung seit
    R2-F-5 sagt.
- Fehlender Extraktor (Hook-Kopie ohne `tools/`): `deny`, Exit 2.
- `make guard-probe`: `Fehlschläge: 0`.

**`SPEC-094` (Handoff-Gate)** gegen das `Makefile` (Rezept von `gates`), `record-gates.sh`,
`working-tree-hash.sh` und `stop-require-gates.sh`:

- Hash: `git ls-files -z --cached --others --exclude-standard | sort -zu`, je Datei `sha256sum`,
  `LINK … -> ziel`, `GONE …`. Das deckt sich mit der Festlegung. Eine unlesbare Datei bricht den Hash ab,
  weil `set -e` in der Pipeline-Subshell erbt; gemessen unten.
- Die Stop-Hook-Zweige decken sich mit dem Text: Schleifen-Schutz, kein Nachweis plus `git status`, Hash
  und Nachweis lesen, Vergleich.

**`SPEC-095` (git-Hooks)** gegen `Makefile:hooks`, `.githooks/{commit-msg,pre-commit}` und
`.github/workflows/ci.yml`: `core.hooksPath .githooks`; `commit-msg` → `make trace-check MSGFILE`;
`pre-commit` → `adr-check STAGED=1`, `doc-check`, unter `set -e`. Die Erkennung
`--diff-filter=AR … done/` mit `^docs/plan/planning/done/(.+/)?slice-[0-9]+[^/]*\.md$` steht wortgleich in
Hook und CI, ohne `-q`.

**`SPEC-096` (`blackbox-probe`)** gegen `tools/blackbox-probe.sh` und `Makefile:blackbox-probe`:

- Die Abbruch-Fälle der Festlegung entsprechen den 11 `fail`-Stellen des Skripts (gezählt).
- `PROBE_FORMS` hat den Default `${…:-- --json --yaml --doctor}`; reiner Leerraum bricht ab.
- Kanarienlauf `sauber:0 links:1` vorher und nachher auf beiden Images; ein Exit außer 0/1/2 bricht ab.
- `VERSION=0.0.0-dev` gilt auf beiden Seiten.
- Gefahren: `REF` leer, `REF` kein Commit, `PROBE_FORMS=' '` und ein fehlendes Nachher-Image enden je
  mit Exit 2 und der Meldung aus der Festlegung.

Referenz-Richtung (§3.4): Die Zeilen `SPEC-093` bis `SPEC-096` und die §8-Zeilen nennen weder ADR noch
Slice, Welle, MR oder Hash (grep, leer). Die §8-Zeile vom 2026-10-08 ist byte-gleich mit `e4505c7f`
(R2-F-3 geschlossen). Der Nachzug steht als eigene Zeile vom 2026-10-09.

### DoD 2 — Sensor-Dateien verlinken die Kennung; `make gates` grün — **BESTÄTIGT, mit Rest (V-2, V-4)**

- Die Kennung verlinken `guard-probe.md` (`SPEC-093`), `hooks.md` (`SPEC-095`), `blackbox-probe.md`
  (`SPEC-096`, zweimal), `verify-closure-notes.md` (`SPEC-095`) und `harness/README.md` in der Zeile
  `record-gates` (`SPEC-094`). Keine dieser Dateien widerspricht den Kennungen.
- `verify-closure-notes.md` (Vertrag „direkt", Grenze 8, Bindung) und `hooks.md` (Vertrag, Grenze 4)
  sind seit R2-F-1 ehrlich. Die Aussage habe ich am Gegenstand gemessen: Im Klon trägt eine Kopie des
  offenen Plans als `done/slice-998-…` 2 Befunde; dieselbe Kopie als `done/wellenlos/slice-999-…` trägt 0.
- Die Lücke trägt `slice-263`. §3 dort nennt die vier Aussagen, die er zurücknimmt (R2-F-4).
- **`make gates`, eigener Lauf:** Exit 0. Abschlusszeile:
  `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`.
  Der geschriebene Nachweis `61af7400…28f7` ist gleich dem Hash des Arbeitsbaums.
- Ergänzend über die Slice-Range `b6a1c21e~1..HEAD`: `make trace-check` ergab 18 Commits und 0 Befunde,
  `make adr-check` 0 Befunde.

### DoD 3 — Review, Report, Verifikation — **BESTÄTIGT** (dieser Bericht ist die Verifikation)

R1 und R2 liegen vor. R2-F-1 bis R2-F-5 sind laut Plan-Änderung nach R2 und `4db45aa3` eingearbeitet.
Nachgefahren habe ich F-2 (Stop-Hook mit kaputtem Index, siehe unten), F-3 (Historie) und F-5
(Extraktor-Satz). R2-F-6 bis R2-F-8 (INFO) sind nicht eingearbeitet; F-6 führe ich als V-3, F-8 ist
Closure-Sache.

### DoD 4 — Closure-Notiz, Register, Risiko-Ausgänge, Paarungen, MR-074-Vermerk — **OFFEN (stand-bedingt)**

§7 ist leer. Das ist Schuld vor dem `git mv`, kein Befund. Für die Closure siehe V-5 und V-6.

---

## Bewusstes Brechen der DoD-relevanten Fixes

### (a) Übergangs-Erkennung — alt (`ffedcb95^`) gegen neu

| Probe | alt | neu |
|---|---|---|
| Diffs echter Commits: `482f1047` (wellenlos), `b175e4a5` (wellenlos), `58328990` (welle-91) | 0 | 1 |
| `5c4a0a07` (direkt unter `done/`) | 1 | 1 |
| `c4049dad` (Plan in `open/`) | 0 | 0 |
| alle Closure-Moves seit 2026-09-29 | 0/19 | 19/19 |
| echter `pre-commit` im Klon, `make` als Shim, `git mv` nach `done/wellenlos/` | kein `verify-closure-notes` | `verify-closure-notes` gerufen |
| dasselbe nach `done/welle-99/` | kein Ruf | gerufen |
| dasselbe direkt nach `done/` | gerufen | gerufen |
| Änderung in `open/`, `M` eines `done/`-Slice | — | nicht gerufen (richtig) |
| SIGPIPE-Modell unter `pipefail`: Treffer in Zeile 1, 20000 Folgezeilen | `grep -q` 0/10 | `grep > /dev/null` 10/10 |

Die Zählung im Plan habe ich nachgezählt: 20 Closure-Renames seit 2026-09-29 (17 nach `wellenlos/`, 3 nach
`welle-91/`), keiner direkt unter `done/`. Sie verteilen sich auf 19 Commits. Die Ursache des Rots ist
gelesen: Das alte Muster verlangt `done/slice-…` ohne Zwischensegment.

### (b) Nachweis unter `make -k`/`make -i`

Modell-Makefile (`vf/mk/`) mit der **wörtlichen** Rezeptzeile aus `Makefile:356-360`, `.NOTPARALLEL` und
einem schaltbar roten Glied. Daneben steht die alte Form mit `record-gates` als letztem Prerequisite.
Gefahren mit GNU Make 4.3.

| Form, Glied rot | alt: Nachweis | neu: Nachweis / Exit / Meldung |
|---|---|---|
| ohne Flag | nein | nein / 2 / — |
| `-k`, `-j4 -k` | **ja** | nein / 2 / — |
| `-i`, `-ik`, `--ignore-errors`, `-j2 -i` | **ja** (mit `green`) | nein / 0 / „kein Nachweis unter make -i" |
| `-s -w`, `FOO=iii` | nein | nein / 2 / — |
| `MAKEFLAGS=i` aus der Umgebung | — | nein / 0 |
| `make -i gates MAKEFLAGS=` | — | **ja** / 0 (siehe V-3) |

Bei grünem Glied schreibt die neue Form den Nachweis ohne Flag, mit `-k`, `-s -w` und einer Variablen mit
`i`. Unter `-i` schreibt sie ihn nicht, mit Meldung.

**Echte Läufe im Repo:**

- `make -k gates THRESHOLD=99`: `coverage-gate: FAIL — Coverage 94.70% unter Schwelle 99%`, Exit 2,
  `.harness/state/` danach leer. Kein Nachweis.
- `make -i gates`: Meldung „gates: kein Nachweis unter make -i …", Exit 0, kein Nachweis.

### (c) Stop-Hook — neu (`HEAD`) gegen alt (`596630db^`)

Wegwerf-Repo mit den echten Skripten (`vf/stop.sh`):

| Fall | neu | alt |
|---|---|---|
| kein Nachweis, sauber | approve | approve |
| kein Nachweis, Änderung | block | block |
| Nachweis passt; danach Commit ohne Inhaltsänderung | approve; approve | approve; approve |
| Inhalt nach dem Nachweis geändert | block | block |
| `stop_hook_active: true` | approve | approve |
| unlesbare Datei (`chmod 000`), mitten in der Sortierung bzw. als letzte | block | Exit 1 ohne Ausgabe (**fail-open**) |
| Datei beim Nachweis write-only, Inhalt danach geändert | block | Exit 1 ohne Ausgabe |
| unlesbarer Nachweis | block | Exit 1 ohne Ausgabe |
| leerer Nachweis | block | block |
| **kaputter `.git/index`, kein Nachweis, Änderung** | **block** | **approve** |
| kaputter Index, mit Nachweis | block | Exit 128 ohne Ausgabe |
| Hash-Skript fehlt | block | Exit 127 ohne Ausgabe |
| außerhalb eines git-Repositorys | block | approve |

Jeder Lauf der neuen Fassung endet mit Exit 0 und gültigem `decision`-JSON. Die alte Fassung wird in
sechs Fällen aus dem richtigen Grund rot (fail-open); die neue schließt alle sechs. Die Zusage
„jede Inhaltsänderung … ändert ihn" hält auch für eine Datei, die beim Nachweis unlesbar war. Das liegt
nicht am Hash: `record-gates` scheitert dort, die Umleitung hat den alten Nachweis aber schon geleert,
und der leere Nachweis blockt.

---

## Plan-Konformität und Abgrenzung (§1, §3)

- **Geänderte Dateien** (`git diff --stat b6a1c21e~1..HEAD`, 20 Dateien) stehen alle in §3 einschließlich
  der Plan-Änderungen. Ausnahmen sind Planungs-, Review- und Lifecycle-Dateien sowie die Zeile in
  `MR-074` (der Teilungs-Vermerk aus `b6a1c21e`).
- **Code-Änderungen:** Jede Plan-Änderung steht im Commit **vor** dem Code, den sie trägt (`04b9c8c7` →
  `ffedcb95`, `eea9b05f` → `5bef299c`, `98037999` → `e4505c7f`, `01d7c500` → `596630db`, `4f561612` →
  `4db45aa3`). Der Nachtrag in §8 nennt `HARN` als berührt (GF).
- **Abgrenzung eingehalten:**
  - §7 trägt keine Netz- oder Nachtlauf-Werkzeuge; `image-scan`, Versions- und Digest-Achsen,
    `baseline-freshness` und `nightly-state` übernimmt `slice-262`, das in `open/` liegt.
  - §7 trägt keine Festlegung zu `image-test`, `trace-check` oder `adr-check`. `SPEC-095` nennt nur,
    welcher Hook welche Prüfung ruft, so wie §1 es ankündigt.
  - Die Behebung der Blindheit von `verify-closure-notes` steht nicht in diesem Slice, sondern in
    `slice-263` (`open/`, Ziel, DoD und die vier Rücknahmen in §3).
  - Bewegende und sagende Werkzeuge sind nicht berührt.
- **`MR-076`** folgt der Vorlage mit allen Pflichtfeldern und der Grenze-Zeile. Die Index-Zeile in
  `harness/conventions.md` hat Anker und Geltungsbereich. `MR-004` ist unverändert und wird geschärft,
  nicht umgeschrieben. Der Kopfkommentar von `record-gates.sh` nennt den Ort des Aufrufs. Zum Umfang
  siehe V-1.

---

## Befunde

### V-1 — MEDIUM: Die Härtung des Stop-Hooks hat keinen Adaptions-Eintrag; `MR-076` deckt nur das Rezept von `gates`

- **Gegen:** Plan §3, Zeile „die Härtung am Handoff-Gate landet als eigener Eintrag, der `MR-004`
  schärft". Baseline `modul-13-quality-gates.md` §Guard-Härtung: „Jede Härtung landet als neuer
  `MR-<NNN>`" und „Die Grenz-Zeile wird mitgezogen".
- **Befund:** `SPEC-094` fasst Rezept **und** Stop-Hook als ein Handoff-Gate. Der Slice hat den Stop-Hook
  zweimal gehärtet: R1-F-5, fail-closed bei Hash- oder Nachweis-Fehler, und R2-F-2, fail-closed bei
  scheiterndem `git status` und außerhalb eines Repos. Der Bruchtest (c) zeigt sechs vorher fail-offene
  Fälle. `MR-076` nennt im Geltungsbereich nur `gates` und `record-gates.sh`, und seine Adaption nennt den
  Hook nicht. `MR-005` (Geltungsbereich `.claude/hooks/`) beschreibt die Restlücke „frischer Klon …
  freigegeben", nicht die neue Fehlerpolitik.
- **Szenario:** Wer den Hook später lockert, etwa „bei Hash-Fehler freigeben, das nervt", findet die
  Entscheidung und ihren Anlass nur in einer Spec-Zeile und in Review-Reports, nicht im Konventionsspeicher.
  Der Retirement-Check hat dann keine Herkunft.
- **Vorschlag:** Vor der Closure `MR-076` erweitern; der Eintrag entsteht in diesem Slice, das hat R2
  ausdrücklich als zulässig bewertet. Geltungsbereich um `.claude/hooks/stop-require-gates.sh` ergänzen,
  die Adaption um die Fehlerpolitik ergänzen. Alternativ ein eigener `MR`.

### V-2 — LOW: Die neue Durchlass-Klasse „Präfix mit Flag" fehlt im Grenze-Kommentar des Wächters und in der Grenze von `guard-probe.md`

- **Gegen:** DoD 2 („widersprechen ihnen nicht") und Baseline `modul-13` §Guard-Härtung („Die
  Grenz-Zeile wird mitgezogen"). Die Klasse gehört zu `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`.
- **Befund:** `SPEC-093` nennt als gemessene Durchlass-Klasse „hinter einem Flag eines Präfixes
  (`env -i`, `sudo -u`)". Ich habe sie bestätigt: `env -i pip x` und `sudo -u root pip x` lassen durch.
  - `guard-probe.md` Grenze 1 spricht von „die vier Umgehungs-Klassen der Segmentierung" und verweist auf
    die Tabelle in `MR-042`, die genau vier führt, ohne diese.
  - Der Kopfkommentar `GRENZE, Umfang — gemessen` im Hook zählt vier Klassen plus `find -exec`, `awk` und
    unbekannte Interpreter auf, ebenfalls ohne diese.
  - Die Spec (Rang 2) ist vollständiger als die beiden Spiegel.
- **Vorschlag:** In `guard-probe.md` Grenze 1 die Aufzählung durch den Verweis auf `SPEC-093` ersetzen,
  statt „vier" zu sagen. Im Hook-Kommentar die Klasse ergänzen.

### V-3 — LOW: `SPEC-094` sagt ohne Einschränkung „unter `make -i` … kein Nachweis"; `make -i gates MAKEFLAGS=` schreibt ihn

- **Gegen:** `SPEC-094` und `AGENTS.md` §5 Regel 13 (Grenze gegen den Gegenstand). Das ist R2-F-6 (INFO),
  nicht eingearbeitet.
- **Befund:** Am Modell gemessen: Glied rot, `make -i gates MAKEFLAGS=` ergibt Nachweis, Exit 0. Die
  Erkennung liest die überschreibbare Variable. `MR-076` deckt das unter „keine Sperre gegen Absicht", die
  Festlegung sagt es nicht.
- **Vorschlag:** Einen Halbsatz in `SPEC-094` aufnehmen: Eine Zuweisung an `MAKEFLAGS` auf der
  Kommandozeile umgeht die Erkennung. Alternativ bewusst stehen lassen und in der Closure begründen.

### V-4 — INFO: `blackbox-probe.md` Grenze 4 führt die Randformen des Kanarienlaufs weiter selbst

Die Grenze nennt `sauber` 0 und `links` 1 sowie den Abbruch bei jedem Exit außer 0, 1 und 2 (125, 126,
127, 137). Das ist dieselbe Randform wie in `SPEC-096`. Heute ist beides konsistent, DoD 2 („statt
Schwelle und Randform zu führen") ist hier aber nur teilweise umgesetzt. Die Erklärung der Restgrenze
braucht einen Teil davon; die Exit-Liste könnte auf `SPEC-096` zeigen. Das Drift-Risiko gehört zu
`BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`, die seit R1 F-7 bei drei
Treffern steht (R2-F-8).

### V-5 — INFO: Drei Werkzeuge liegen weder im Ziel noch in einer Abgrenzung

`history-range-guard`, `selbstpruefung` und `baseline-probe` stehen in der Werkzeug-Tabelle von
`harness/README.md`. Weder `slice-259`, `slice-260` noch `slice-262` nennen sie im Ziel oder in der
Abgrenzung. `history-range-guard` trifft eine Entscheidung: Eine leere oder unauflösbare Range ist rot.
Der Satz in `MR-074` Bewegung 2 („Die übrigen Werkzeuge übernehmen slice-260 … und slice-262 …") liest
sich als vollständig. Beim Vermerk von Bewegung 2 in der Closure sollte das zugeordnet werden: als
„setzt eine Anforderung durch" oder „Probe eines Gates", oder als Rest an `slice-262`.

### V-6 — INFO: Die eigene Closure dieses Slice prüft kein Sensor

`slice-260` schließt nach dem Muster der letzten 20 Closures nach `done/wellenlos/`. Dort sieht ihn
`make verify-closure-notes` nicht (`SPEC-095`, Messung oben). Der Übergang löst den Lauf aus, und der
Lauf wird grün, ohne die Notiz gelesen zu haben. Bis `slice-263` schließt, sollte das Closure-Profil von
Hand gegen die Notiz laufen, wie in dieser Probe: eine Kopie direkt unter `done/` in einem Wegwerf-Klon,
dann `docker run … --config .d-check.closure.yml`. Ebenfalls Closure-Sache sind der dritte Treffer von
`guide-doku-…` (R2-F-8) und der Plan-Hinweis auf `begruendung-traegt-entscheidung-nicht` (drei Belege,
ohne Ausgang).

---

## Negativbefunde (geprüft, ohne Befund)

- **Zählungen in Plan und Botschaften:**
  - Die 26 Wörter der Sperrliste und die 26 + 2 in `settings.json` stimmen.
  - Die 11 Abbrüche von `blackbox-probe.sh` stimmen.
  - Die 20 Closures (17/3) stimmen.
  - Die Messwerte der Botschaft `ffedcb95` stimmen: `482f1047`, `58328990`, `b175e4a5` alt 0 / neu 1,
    `5c4a0a07` 1/1, `c4049dad` 0.
- **Rezept-Logik:** Scheitert `record-gates.sh`, fehlt die `green`-Zeile und `gates` endet rot (`&&`
  im `case`-Zweig). `record-gates` als Target bleibt Werkzeug und steht in keiner Prerequisite-Liste
  mehr (grep).
- **Konsumenten:** `harness/README.md` („`record-gates` als letzter Schritt"),
  `.claude/commands/{plan,close}-welle.md` und `MR-004` (geschärft, nicht umgeschrieben) stimmen weiter.
- **Kommentare (§3.7)** an `.NOTPARALLEL`, `gates`, Hook, `pre-commit` und `ci.yml` tragen Zusage,
  Kopplung, Grenze und Rang-Zeiger. Neue Slice-Nummern habe ich keine gefunden.
- **Traceability:** Alle 18 Commits tragen eine Kennung (`make trace-check` über die Range).

## Kategorie-Summary

| Kategorie | Anzahl | Befunde |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 1 | V-1 |
| LOW | 2 | V-2, V-3 |
| INFO | 3 | V-4, V-5, V-6 |

## Verdikt

**DoD 1–3 bestätigt, DoD 4 offen (Closure, stand-bedingt). Eine Lücke in der Plan-Konformität (V-1) ist
vor der Closure zu schließen.**

Die vier Festlegungen stimmen mit dem Code überein, jede Liste ist nachgezählt. Die drei
Korrektheits-Fixes werden im Bruchtest aus dem richtigen Grund rot und schließen die gemessenen Lücken:
die Übergangs-Erkennung (0/19 → 19/19), den Nachweis unter `-k`/`-i` (Modell und echter Lauf) und den
fail-offenen Stop-Hook (sechs Fälle). `make gates` ist im eigenen Lauf grün. Die Abgrenzung ist
eingehalten, und das Ausgeschlossene hat in `slice-262` und `slice-263` eine Adresse, die es annimmt.
