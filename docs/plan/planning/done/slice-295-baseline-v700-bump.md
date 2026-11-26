# Slice slice-295: Baseline-Pin auf `v7.0.0` heben — mechanisch, ohne eine Neuerung zu übernehmen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-98-baseline-v700`](../welle-98-baseline-v700.md).

**Bezug:** [`MR-080`](../../../../harness/conventions.md#mr-080) (die vorige Pin-Hebung und ihre Form),
[`MR-011`](../../../../harness/conventions.md#mr-011) (die Pin-Kette),
[`MR-021`](../../../../harness/conventions.md#mr-021) (pin-gebundene Verweise),
[`MR-039`](../../../../harness/conventions.md#mr-039) (Zitat-Delta),
[`MR-051`](../../../../harness/conventions.md#mr-051) (`d-check:cite`-Spannen),
[`MR-052`](../../../../harness/conventions.md#mr-052) (historisches Zitat, alter Baum in der git-Historie),
[`MR-055`](../../../../harness/conventions.md#mr-055) (Symlinks als Träger),
[`MR-069`](../../../../harness/conventions.md#mr-069) (`ignore-refs`),
[`MR-070`](../../../../harness/conventions.md#mr-070) (Frozen-Klassen vor mechanischer Ersetzung).
Anlass: der Nachtlauf `upstream-drift` meldet den neuen Release `v7.0.0` (Pin `v6.18.1`).

**Berührte Spec-Stellen:** der Kopf von `spec/spezifikation.md` und von `spec/architecture.md` — je die Zeile
**Rolle** mit der Pin-Nennung (`Baseline-Regelwerk v6.18.1`); keine Kennung trägt das Zielelement, der Abschnitt
ist der Kopf. Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie.

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-11-25.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Den committet vendorten Baseline-Bestand von `v6.18.1` auf
[`v7.0.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v7.0.0) heben, den Pin und alle lebenden
pin-gebundenen Verweise samt der `d-check:cite`-Spannen nachziehen, die sieben umbenannten Vorlagen in den
Verweisen des Repos beim neuen Namen nennen — und die Hebung als Eintrag in `harness/conventions/` führen
(Delta vollständig nach Datei, **keine Neuerung übernommen** — die Namensregel führt nicht dieser Slice, sondern
[`slice-297`](slice-297-spezifikation-heisst-specification.md) und
[`slice-298`](../next/slice-298-neue-namen-wave-archive-anker.md); Entscheid 1 der [Welle](../welle-98-baseline-v700.md)).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Keine Übernahme einer Neuerung.** Jede ist eine Konventions- oder Produkt-Entscheidung mit eigener
  Begründung (Entscheid 2); der Bump ist ein anderer Vorgang. Der Bump-Eintrag nennt das Delta vollständig
  nach Datei und markiert jede Bewegung als *nicht übernommen*; die Klassifikation mit Empfehlung liefert
  [`slice-296`](../open/slice-296-baseline-v700-neuerungen-sichtung.md), die Übernahme ein eigener Folge-Slice nach der
  Antwort des Auftraggebers. Die **Bezeichner** `wave-`, `archive.zip`, `seit wave-` sind **keine** Neuerung dieser
  Art mehr — der Auftraggeber hat sie entschieden — und gehören in
  [`slice-298`](../next/slice-298-neue-namen-wave-archive-anker.md); dieser Slice zieht nur die **Namen der
  Vorlagen-Dateien** nach, weil ein Verweis auf eine Datei, die es nicht mehr gibt, ein toter Verweis ist.
- **Keine Umbenennung von `spec/spezifikation.md`.** Sie ist entschieden (Entscheid 1) und ein eigener Slice,
  [`slice-297`](slice-297-spezifikation-heisst-specification.md), mit eigenem Commit-Zwang. Der Bump nennt die
  Pin-Zeile der Spezifikation (Zeile 8) beim **heutigen** Namen; Pfade in `manual.go` und der
  [`.d-check.yml`](../../../../.d-check.yml) bleiben.
- **Keine Änderung am Produkt.** Die eine Berührung ist die Pin-Nennung in `spec/spezifikation.md` (Zeile 8), die
  `manual.go` einbettet; sie ist in der Welle (d3) Klasse (iii) benannt und gemessen. Quelltext, Handbuch,
  Standardwerte bleiben.
- **Frozen-Dateien bleiben byte-stabil** — `done/`-Slices und -Wellen, `docs/reviews/`, `CHANGELOG.md`,
  `conventions/done/`, eingehende CRs und der Kern der `Accepted`-ADRs nennen den alten Pin als
  Vergangenheits-Aussage (Bestand bewusst). Die Ausnahme ist allein die Antwort auf Frage 5 in
  [welle-98 §5](../welle-98-baseline-v700.md): sieben eingefrorene Dateien **verlinken** den alten Baum; nach dem Nachtrag V1 des Verdikts sind es fünf hinter dem Ventil, die zwei übrigen werden retargetet (Plan-Änderung in §3).
- **`tools/harness/selbstpruefung.sh`** — sein Stand-Vermerk nennt die Version des Schwester-Werkzeugs
  `ai-harness-init`, nicht den Baseline-Pin (anderer Gegenstand; dieselbe Abgrenzung wie in
  [`slice-277`](../done/wellenlos/slice-277-baseline-v6181-bump.md)).
- **Die Tool-Pins** (Go, golangci-lint, semgrep, a-check) und Dependabot-Hebungen — eigene Vorgänge mit eigener Messung.
- **Kein Release** — der Bump berührt keine Distributions-Fläche außer der eingebetteten Pin-Nennung.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Vendoring:** `.harness/baseline/v7.0.0/` steht über `tools/harness/fetch-baseline-cache.sh v7.0.0`
      (56 Dateien, 55 im Manifest), `make baseline-verify` Exit 0, der Baum `v6.18.1` (55 Dateien) ist entfernt;
      der Pin steht in `harness/conventions.md` §Baseline und §Adoptierte Konventions-Quellen.
- [x] **Spiegel:** Jeder lebende pin-gebundene Verweis nennt `v7.0.0` (die Klassen und ihre Zahlen in §3), die
      7 Symlinks unter `.claude/rules/` lösen auf den neuen Baum auf, jede `d-check:cite`-Spanne steht gegen den
      **Wortlaut** des neuen Baums (die 7 verschobenen gelesen, nicht nur verschoben), die Vorlagen-Namen in den
      Verweisen sind die neuen, und `make doc-check` meldet 0 Befunde.
- [x] **Einträge:** Der Bump-Eintrag steht unter `harness/conventions/` (Delta **vollständig nach Datei**,
      Frozen-Klassen am echten Vorzustand gemessen, Zitat-Delta nach
      [`MR-039`](../../../../harness/conventions.md#mr-039), jede Bewegung als *nicht übernommen* markiert);
      [`MR-080`](../../../../harness/conventions.md#mr-080) liegt per reinem `git mv` in `conventions/done/`; der
      Bump-Eintrag steht im Index.
- [x] `make gates` grün; `make archive-wave-test` grün (der Kommentar-Nachzug in
      `tools/archive-wave/stub.go` ist bestätigt, Verdikt Frage 12).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

**Die Zählung nennt ihre Form** ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 14): gezählt wird die **git-getrackte
Datei**, die die Zeichenkette `v6.18.1` trägt (`git grep -l`), nicht der Verweis auf den Baum; Symlinks liest
`git grep` nicht und stehen eigens. Stand `76cf7d51`, vor Anlage der Pläne dieser Welle; der Implementer
**wiederholt** die Zählung vor dem ersten Tag-Tausch ([`MR-070`](../../../../harness/conventions.md#mr-070),
[`BEO-ALL/pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)).

| Datei / Komponente | Änderungs-Art | Begründung und Zahl |
|---|---|---|
| `.harness/baseline/v7.0.0/` | neu | **56** Dateien (`find .harness/baseline/v7.0.0 -type f \| wc -l`), davon 55 im Manifest; über `tools/harness/fetch-baseline-cache.sh v7.0.0`, nicht aus dem Scratchpad |
| `.harness/baseline/v6.18.1/` | entfernen | **55** Dateien (`git ls-files .harness/baseline/v6.18.1 \| wc -l`); [`MR-021`](../../../../harness/conventions.md#mr-021) (1) |
| `.claude/rules/*.md` | update | **7** der 9 Symlinks zeigen in den Baum (`readlink .claude/rules/*`); `AGENTS.md` und `harness-README.md` zeigen nicht hinein. Alle 7 Ziel-Dateinamen gibt es im neuen Baum |
| `harness/conventions.md` | update | §Baseline (Stand, Pin-Zeiger), §Adoptierte Konventions-Quellen (Release-, Asset- und Baum-Verweis), der Pfad-Link im Konventionsspeicher-Absatz, der Vorlagen-Name `MR-<NNN>-titel.template.md` (Zeile 59), Index: zwei neue Zeilen, [`MR-080`](../../../../harness/conventions.md#mr-080) von *Aktive* nach *Aufgelöste* |
| neue Eintrags-Datei unter `harness/conventions/` (Bump) | neu | aus `MR-<NNN>-title.template.md` des **neuen** Baums, per `cp` ([`BEO-ALL/form-vom-nachbarn-statt-von-der-vorlage`](../observations/BEO-ALL/form-vom-nachbarn-statt-von-der-vorlage/observation.md)); Nummer und Form: im [Architect-Verdikt](../../../reviews/2026-11-26-welle-98-architect-verdikt.md) (Frage 1; höchste bisher vergebene: [`MR-080`](../../../../harness/conventions.md#mr-080)); der alte Baum fällt im Bump, die Abweichung vom Kanon trägt [`MR-021`](../../../../harness/conventions.md#mr-021) (1) schon — der Eintrag zeigt dorthin und nennt den Vor-Bump-Commit, ohne eigenen Satz (Frage 9); Delta nach Datei: **31** inhaltlich verschieden + **7** umbenannt + **1** neu, die **7** Stempel-Dateien in einer Zeile |
| die Eintrags-Datei von [`MR-080`](../../../../harness/conventions.md#mr-080) unter `harness/conventions/` | move | reiner `git mv` nach `conventions/done/`, dann die Link-Tiefen als eigener Commit ([`AGENTS.md`](../../../../AGENTS.md) §3.3) |
| 35 aktive `MR-*`-Dateien unter `harness/conventions/` | update | Tag-Tausch je Datei; die Mischzeilen **einzeln gelesen**: [`MR-021`](../../../../harness/conventions.md#mr-021) Z. 20 („aktuell …/v6.18.1/…“), [`MR-049`](../../../../harness/conventions.md#mr-049) Z. 14, [`MR-053`](../../../../harness/conventions.md#mr-053) Z. 9 |
| `AGENTS.md`, `harness/README.md`, 3 Dateien unter `harness/rules/`, 2 unter `.harness/skills/`, `.claude/agents/reviewer.md` | update | 8 Dateien; `AGENTS.md` Z. 345 („Baseline-Form `v6.18.1` §1“) ist eine Mischzeile und wird gegen die Gliederung der neuen `slice.template.md` gelesen (die Abschnitte §1 bis §8 bleiben, gemessen) |
| `docs/plan/planning/README.md`, `in-progress/roadmap.md`, `observations/README.md` | update | 3 Dateien; `roadmap.md` Z. 13 und `observations/README.md` Z. 18 sind Prosa-Pins, einzeln gelesen |
| `spec/spezifikation.md`, `spec/architecture.md` | update | 2 Dateien, je die **Rolle**-Zeile; die Spezifikation ist eingebettet (Welle (d)) |
| 15 lebende `d-check:cite`-Spannen | neu ankern | `.harness/skills/reviewer.md` 3, `closure-note-reviewer.md` 3, [`MR-031`](../../../../harness/conventions.md#mr-031) 3, [`MR-049`](../../../../harness/conventions.md#mr-049) 2, [`MR-005`](../../../../harness/conventions.md#mr-005) 1, [`MR-035`](../../../../harness/conventions.md#mr-035) 1, [`MR-039`](../../../../harness/conventions.md#mr-039) 1, [`MR-043`](../../../../harness/conventions.md#mr-043) 1. **7 verschieben sich** (gemessen am Vergleich derselben Zeilenbereiche beider Bäume): `grundlagen-begriffe` 50, `grundlagen-harness-dateien` 196, `modul-05` 170 und 180–181, `modul-10` 99, `reviewer.template` 87–88 und 99–100. Die 8 übrigen stehen auf unverschobenen Zeilen. **Plus 8** in den vier Plänen dieser Welle (je Slice zwei, `modul-05` 373–374 und 379; sie verschieben sich ebenfalls). Zwei weitere Spannen unter `conventions/done/` ([`MR-057`](../../../../harness/conventions.md#mr-057), [`MR-058`](../../../../harness/conventions.md#mr-058)) zielen auf `v5.18.0` und liegen außerhalb des Scopes von `citations` |
| `.claude/commands/plan-welle.md` (2 Zeilen), `close-welle.md` (4 Zeilen) | update | die vier umbenannten Vorlagen: `welle.template.md`→`wave.template.md`, `welle-results.template.md`→`wave-results.template.md`, `archiv-stub-slice`/`-welle`→`archive-stub-slice`/`-wave`; **9 Zeilen in 4 Dateien** mit `stub.go` und `conventions.md` Z. 59 (`git grep -n "welle\.template\|welle-results\.template\|archiv-stub-\(slice\|welle\)\.template\|titel\.template" -- . ':!.harness/baseline' ':!docs/plan/planning' ':!docs/reviews' ':!harness/conventions/done'`; die Treffer in `conventions/done/` und in den Plänen nennen die Namen als Vergangenheits-Aussage und bleiben) |
| `tools/archive-wave/stub.go` | update (Kommentar) | zwei Kommentarzeilen nennen `archiv-stub-slice.template.md`/`-welle.template.md`; Nachzug im Bump bestätigt ([Verdikt](../../../reviews/2026-11-26-welle-98-architect-verdikt.md), Frage 12); kein Verhalten |
| [`.d-check.yml`](../../../../.d-check.yml) | update | nur, was die Antwort auf Frage 5 verlangt: je eingefrorener Datei mit Link auf den alten Baum **ein** `ignore-refs`-Eintrag (gemessen 7; nach dem Nachtrag V1 des Verdikts **5**, Plan-Änderung in §3); heute 30 Einträge über 12 Tags |
| dieser Plan und die Pläne [`slice-296`](../open/slice-296-baseline-v700-neuerungen-sichtung.md), [`slice-297`](slice-297-spezifikation-heisst-specification.md), [`slice-298`](../next/slice-298-neue-namen-wave-archive-anker.md) | update | die je zwei `d-check:cite`-Spannen in §8 auf den neuen Baum (zusammen 8) |

**Plan-Änderung (Implementer, vor dem Code):** Zusätzlich zur Tabelle trägt dieser Slice **eine** Datei aus dem
Frozen-Bestand, die die Spiegel-Messung nicht sah, weil sie nicht den Pin, sondern den **Pfad** der Eintrags-Datei
von [`MR-080`](../../../../harness/conventions.md#mr-080) nennt: `docs/reviews/2026-10-20-slice-277-baseline-v6181-verify.md`
Z. 4 verlinkt die Eintrags-Datei `conventions/…-baseline-v6181.md` (`git grep -n 'conventions/MR-0[8]0-baseline-v6181' -- . ':!docs/plan/planning/done'
':!.harness/baseline'` nennt außer der Index-Zeile in `harness/conventions.md` nur Inline-Code in drei Reviews, die `codepaths`
ausnimmt, und diesen einen Link). Nach dem Move nach `conventions/done/` meldet `make doc-check` ihn als `target-missing`.
**Antwort:** der Link wird auf `conventions/done/` nachgezogen — ein reiner Pfad-Nachzug, das Ziel lebt, nur sein Ort wechselt
(dieselbe Begründung wie in Frage 5 für die Spezifikation; Vorgang und Zähler wie bei `slice-274` in
[`MR-080`](../../../../harness/conventions.md#mr-080): 1 Datei, 1 Zeile). Die Alternative, ein achter `ignore-refs`-Eintrag, wäre
eine Gate-Senkung ([`MR-069`](../../../../harness/conventions.md#mr-069)) für einen Link, der weiter auflösen kann. Der Nachzug
reist im Link-Tiefen-Commit (Schritt 4), der Move bleibt rein; der Bump-Eintrag führt ihn als Ausnahme von der Frozen-Regel.
**Folge für Risiko 2 (§6):** `git diff --stat <Vor-Bump>..HEAD -- docs/plan/planning/done docs/reviews CHANGELOG.md harness/conventions/done`
ist nicht leer, sondern nennt genau diese eine Review-Datei und die zugezogene Eintrags-Datei von [`MR-080`](../../../../harness/conventions.md#mr-080) unter `conventions/done/`.

**Plan-Änderung (Implementer, Fixrunde nach Review r1, gedeckt durch den Nachtrag V1 im
[Verdikt](../../../reviews/2026-11-26-welle-98-architect-verdikt.md) zu Frage 5):** (1) Das Ventil in der
[`.d-check.yml`](../../../../.d-check.yml) trägt **fünf** Einträge (die `Accepted`-ADRs 0112 bis 0116), nicht sieben: ihr
Zitat (`modul-04` §Ziel-Form: ADR) ist in `v7.0.0` nicht wortgleich. (2) Zwei Dateien kommen zur Tabelle, je **eine Zeile,
nur das Link-Ziel** (`v6.18.1` → `v7.0.0`, Anker unverändert; die zitierte Stelle ist gelesen wortgleich):
`docs/plan/planning/done/wellenlos/slice-274-trace-requirements-als-liste-von-quellen.md` und
`docs/reviews/2026-10-30-welle-93-architect-verdikt.md`; ihr Ventil-Eintrag entfällt. (3) Eine Datei kommt zur Tabelle auf
Befund F-3 des Reviews: `harness/sensors/baseline-verify.md` Z. 17 nennt die Ausgabe des Laufs mit der Dateizahl des alten
Baums (`54 Dateien`); sie wird auf `55 Dateien` gesetzt — die Zahl, die `make baseline-verify` jetzt meldet. **Folge für
Risiko 2 (§6):** `git diff --stat <Vor-Bump>..HEAD -- docs/plan/planning/done docs/reviews CHANGELOG.md harness/conventions/done`
nennt damit drei eingefrorene Dateien mit einer geänderten Zeile (`slice-274`, das Verdikt zu welle-93 und den Verifikations-Report
von `slice-277`) neben der zugezogenen Eintrags-Datei von [`MR-080`](../../../../harness/conventions.md#mr-080) und den Dateien,
die diese Welle selbst anlegt.

**Reihenfolge und Commit-Zerlegung** ([`AGENTS.md`](../../../../AGENTS.md) §3.3, [`BEO-ALL/path-scoped-commit-carries-staged-rest`](../observations/BEO-ALL/path-scoped-commit-carries-staged-rest/observation.md)):

1. **Vorbedingung, kein Commit dieses Slice:** das Architect-Verdikt zu den Fragen der Welle liegt vor.
2. **Spiegel-Messung vor dem ersten Tag-Tausch:** `git grep -l 'v6\.18\.1' -- . ':!.harness/baseline'` über das **ganze**
   Repo (nicht `.harness/`: dort liegen die Skills), geteilt in Frozen- und lebende Klassen **über die Eigenschaft**
   (liegt in `done/`, `docs/reviews/`, `conventions/done/`, ist eine `Accepted`-ADR, ist `CHANGELOG.md` oder ein
   eingehender CR), nicht über Verzeichnisnamen. Erwartung: 72 Dateien außerhalb des Baums, davon 17 Frozen (8 in
   `done/`, 8 Reviews, 1 `CHANGELOG.md`) und 55 weitere, von denen 5 `Accepted`-ADRs sind.
3. **Ein Hauptcommit:** Baum-Tausch, Pin, Tag-Tausch, Symlinks, Zitat-Spannen, Vorlagen-Namen, `ignore-refs`, **und** die
   den neuen Eintrag samt Index — **nicht** in mehreren Commits, weil der lokale `pre-commit`-Hook den Arbeitsbaum
   mit `doc-check` prüft und ein Zwischenstand mit totem Pfad-Link ihn nicht passiert
   ([`slice-277`](../done/wellenlos/slice-277-baseline-v6181-bump.md) §7: ein `--no-verify` war dort der Fehler).
   `git diff --cached --stat` vor dem Commit, Commit mit `-F <datei> -- <Pfade>`.
4. **Reiner Move-Commit** [`MR-080`](../../../../harness/conventions.md#mr-080) → `conventions/done/`, danach **Link-Tiefen-Commit** (eigener Commit); beide
   gehören in **denselben Push** — der Move allein darf nicht die Spitze eines Push werden
   (`grundlagen-traceability.md` §Herkunfts-Anker, Ruheort-Regel).
5. Review (höchstens zwei Runden), Verifikation, Closure.

**Gates:** `make baseline-verify`, `make doc-check` (darin `links`, `anchors`, `citations`), `make gates`;
`make baseline-freshness` (Netz) am Ende; `make archive-wave-test` bei der Kommentar-Änderung; **nicht** hier:
`make blackbox-probe` und `make fullbuild` — sie sind Closure-Trigger der Welle (§3 (d)).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen Slice (nur die Roadmap); das **Architect-Verdikt**
zu den zwölf Fragen in [welle-98 §5](../welle-98-baseline-v700.md) (für diesen Slice tragend: 5, 9, 12) liegt vor: [Verdikt](../../../reviews/2026-11-26-welle-98-architect-verdikt.md);
Priorisierung gesetzt, `Verantwortlich:` beim Übergang `open` → `next` gesetzt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Bump ist nicht mehr mechanisch — die Messung am
  vendorten Baum zeigt gegenüber §1 der Welle mehr als die 7 umbenannten und die 1 neue Datei (etwa entfallene,
  verlinkte Dateien oder ein umbenannter Anker, den ein lebender Verweis nennt), oder der Diff wächst über eine
  Review-Sitzung; dann gilt die Form von [`MR-066`](../../../../harness/conventions.md#mr-066).
- `in-progress` → `open` (blockiert): Das Release-Bundle ist nicht ladbar, oder sein `SHA256SUMS` stimmt nicht mit dem
  des Release überein (`make baseline-freshness`, Teil *Content*); oder das Verdikt ist nicht eingetroffen.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (§2) und `make gates` grün; Review- und Verifikations-Report liegen unter `docs/reviews/`;
Closure-Notiz mit Lerneintrag geschrieben. **Review-Vorgabe der Welle:** höchstens zwei Runden, eingearbeitet
werden nur HIGH und MEDIUM; jede Nachbesserung zitiert die Stelle, die sie trägt
([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 18). Was die **Welle** einlöst — Gegenmessung am Endstand,
Produkt unverändert, Trigger-Audit — steht in [ihrem Closure-Trigger](../welle-98-baseline-v700.md) und nicht in
dieser DoD ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 10).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Ein lebender Verweis zeigt auf einen Anker oder eine Datei, die der neue Baum nicht mehr trägt.** Gemessen vor
  dem Start: eine Überschrift entfällt (`modul-14`, „Zwei Formen des Reproduzierbarkeits-Ankers“; kein Verweis
  darauf außerhalb des Baums, `git grep -n 'zwei-formen-des-reproduzierbarkeits' -- . ':!.harness/baseline'` ist
  leer), sieben Vorlagen tragen einen neuen Namen. — **Ausgang:** *entfallen.* `make doc-check` meldet nach dem
  Entfernen des alten Baums außer den in §3 benannten Folgen nichts (Verifikation M4: 1319 Dateien, 0 Befunde);
  die Probe ohne die fünf Ventil-Einträge liefert genau fünf `target-missing` auf `modul-04-adrs.md`, sonst nichts.
- **Eine Mischzeile wird pauschal ersetzt und fälscht eine Vergangenheits-Aussage**
  ([`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md),
  5×). — **Ausgang:** *entfallen.* Jede geänderte Zeile ohne Pfad- oder URL-Form steht im Bump-Eintrag einzeln; der
  Verifier fuhr die Schleife des Eintrags nach (14 Dateien, 40 Vorkommen byte-stabil). Die Frozen-Verzeichnisse
  tragen nur die Ausnahme der Antwort auf Frage 5: je eine Zeile, nur das Link-Ziel, in `slice-274`, im Verdikt zu
  welle-93 und im Verifikations-Report von `slice-277`.
- **Eine verschobene Zitat-Spanne wird nur verschoben, nicht gelesen**
  ([`BEO-ALL/citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md),
  20×). — **Ausgang:** *entfallen.* Der Bump-Eintrag nennt je verschobener Spanne alte und neue Zeile und den Befund
  „Wortlaut gleich“; der Verifier las die sieben Spannen und die zwei Plan-Bereiche Zeile für Zeile alt gegen neu und
  fand sie wortgleich. Die Mutation einer Spanne (`115-115` → `114-114`) wird als `citation-mismatch` rot.
- **Eine Zahl im Bump-Eintrag stimmt nicht** — Zählung am Proxy statt am Gegenstand
  ([`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md),
  10×; in [`MR-080`](../../../../harness/conventions.md#mr-080) waren es drei). — **Ausgang:** *eingetreten.* Der
  Reviewer zählte das Zitat-Delta nach (31 / 30 statt „51 / 50“, ohne Kommando, Review R1 F-2); der Verifier fand
  bei der Nachzählung, dass das Kommando zur Treffermenge der Vorlagen-Namen („danach 0“) den Eintrag selbst
  durchsucht und wörtlich sieben Zeilen zeigt (Verifikation V-1). Beide sind behoben (Verfahren im Eintrag; Ausnahme
  im Kommando); Evidenz-Datei `evidence/slice-295.md` unter dem genannten Eintrag, jetzt 11×.
- **Die kopierte Welle-Vorlage trägt nach dem Bump `wave-<Kennung>-results.md`** in ihrer `Lifecycle:`-Note, und
  `plan-welle.md` verlangt, diese Note zu behalten — dieses Repo führt aber bis
  [`slice-298`](../next/slice-298-neue-namen-wave-archive-anker.md) `welle-NN-results.md`. Das ist eine **bekannte
  Folge der Reihenfolge**, kein Fehler des Bumps. — **Ausgang:** *weiter offen* →
  [`slice-298`](../next/slice-298-neue-namen-wave-archive-anker.md), das die Namen neuer Wellen nach der Regel setzt;
  bis dahin ersetzt der Planner den Namen beim Füllen.

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

- **Was hat funktioniert:** Der Zuschnitt „mechanisch, ohne eine Neuerung“ hielt: Der Diff von
  `738d2bca..HEAD` berührt außerhalb von Baum, Plänen, Reviews, `harness/`, `.claude/` und den Skills nur die
  Pin-Zeilen der beiden Spec-Köpfe, zwei Kommentarzeilen in `stub.go` und die [`.d-check.yml`](../../../../.d-check.yml)
  (Verifikation, Abschnitt Abgrenzung). Den Bytebestand des Baums lud der Verifier unabhängig vom Upstream-Asset und
  verglich ihn mit `diff -rq`: einziger Unterschied `SHA256SUMS`. Die Wirkung der Gates brach er am Klon, je aus dem
  richtigen Grund rot: eine verschobene Zitat-Spanne (`citation-mismatch`), ein Ziel im alten Baum
  (`citation-out-of-range`), die fünf Ventil-Einträge entfernt (genau fünf `target-missing`), ein Symlink auf den
  alten Baum (`baseline-verify` Exit 2, `doc-check` bleibt grün). Eine Review-Runde (1 MEDIUM, 2 LOW, 3 INFO) und
  eine Fixrunde.
- **Was ging anders als geplant:** **(1) Das Ventil trägt fünf Einträge statt sieben** (Review R1 F-1 MEDIUM,
  Nachtrag V1 im [Verdikt](../../../reviews/2026-11-26-welle-98-architect-verdikt.md)): Das Zitat von
  `modul-04` §Ziel-Form: ADR ist in `v7.0.0` für die fünf `Accepted`-ADRs nicht wortgleich; für `slice-274` und das
  Verdikt zu welle-93 ist es gelesen wortgleich, deshalb Retarget statt Ventil. **(2) Ein Pfad-Nachzug in einem
  eingefrorenen Report** (`slice-277`-Verifikation, Link auf den verschobenen [`MR-080`](../../../../harness/conventions.md#mr-080)),
  beide Änderungen als Plan-Änderung vor dem Code in §3. **(3) Zwei Zwischenstände sind einzeln rot** (V-4, unten).
  **(4) Die Fixrunde bekam keine zweite Review-Runde**, obwohl sie Gate-Konfiguration änderte.
- **Vorbehalt zum DoD-Punkt „Review durchgeführt“ (V-2, wörtlich aus der Verifikation):** „Die Fixrunde war **nicht
  rein dokumentarisch** (V-2) und hat keine zweite Review-Runde.“ Der Verifier hat das Delta statt dessen selbst gemessen
  (`make doc-check` 0 Befunde; ohne die fünf Einträge fünf `target-missing`; beide Retargets je eine Zeile, die
  zitierten Abschnitte wortgleich) und schließt: „Kein Befund am Inhalt — nur die Lücke im Verfahren ist zu benennen.“
  Der Haken steht mit diesem Vorbehalt.
- **Steering-Loop-Eintrag:** mit diesem Slice ist keine Regel, kein Sensor und keine benannte Spec-Lücke verkörpert;
  zwei Beobachtungen sind belegt (nächster Punkt). **Lernsignal 1, mit Ursache:** Die Treffermenge „danach 0“ stimmte
  am Gegenstand und nicht im Kommando, **weil** der Eintrag die alten Vorlagen-Namen selbst nennt (Z. 128–133 und 216
  von [`MR-081`](../../../../harness/conventions.md#mr-081)) und im Suchraum seines eigenen Kommandos liegt, dessen
  Ausschlussliste ihn nicht nannte; wörtlich gefahren zeigte es sieben Zeilen. Gefunden hat es der Verifier, der das
  Kommando fuhr, statt die Zahl zu lesen — der Eintrag trägt jetzt die Ausnahme im Kommando (Commit `e584680e`).
  **Lernsignal 2, mit Ursache:** `make doc-check` blieb über acht mehrzeilig geöffneten Inline-Code-Spans des Eintrags
  grün, **weil** `codepaths` zeilenlokal urteilt: derselbe erfundene Pfad meldet einzeilig `codepath-missing`,
  mehrzeilig nichts (Verifikation V-3, Klon-Probe) — eine Prüflücke, die weder die Anforderung
  [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
  noch die Spezifikation benennt. Ob daraus eine Spec-Lücke wird, ist eine Lastenheft-Frage des Auftraggebers
  (`AGENTS.md` §5 Regel 3) und hier nicht entschieden.
- **Beobachtungs-Register (`../observations/`):** neu angelegt
  [`BEO-ALL/codepaths-prueft-zeilenlokal-mehrzeilige-spanne-ungeprueft`](../observations/BEO-ALL/codepaths-prueft-zeilenlokal-mehrzeilige-spanne-ungeprueft/observation.md)
  mit `evidence/slice-295.md` (1×, `offen`); `evidence/slice-295.md` in
  [`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
  (11×, `verkörpert`; Risiko 4: „51 / 50“ ohne Kommando, R1 F-2, und das Kommando mit der Selbsttreffer-Menge, V-1).
  Benannt, nicht gezählt: V-2 (die Fixrunde änderte `.d-check.yml` ohne zweite Review-Runde; verwandt mit
  [`BEO-ALL/fix-commit-ausserhalb-review-range`](../observations/BEO-ALL/fix-commit-ausserhalb-review-range/observation.md),
  dessen Regel „ein Fix, der mehr schreibt als seine Befunde“ hier nicht eindeutig zutrifft — der Fix schrieb, was F-1
  verlangte —, weshalb kein Beleg angelegt ist); V-5 (das Zitat-Delta nennt den Bezug der 25 Zeichen nicht, mit dem
  Inneren in Zeichen reproduziert es sich exakt); V-6 beziehungsweise Review F-4 (Bewegung 1, 5, 6 des Eintrags ohne
  eigenes „nicht übernommen“, DoD in der Sammelform erfüllt), F-5 (Baum-Zuwachs als eigener Commit statt „ein
  Hauptcommit“, `79a14c01` einzeln grün) und F-6 (Inline-Code-Nennungen des alten Pfads von [`MR-080`](../../../../harness/conventions.md#mr-080) in drei
  eingefrorenen Reports, `codepaths` nimmt `docs/reviews/**` aus). Keine davon ist umgesetzt.
- **Architektur-Beobachtung:** Der Spiegel eines Pin-Bumps hängt an drei Wächtern mit verschiedenen Gegenständen, und
  keiner ersetzt einen anderen: **Zitat-Fäule** hält `citations` (die Spanne gegen den Wortlaut), **Link-Bruch**
  `links` mit dem Ventil in der [`.d-check.yml`](../../../../.d-check.yml), den **Symlink als Träger** allein
  `baseline-verify` ([`MR-055`](../../../../harness/conventions.md#mr-055)) — `doc-check` liest den Inhalt eines
  Symlinks nicht und blieb in der Mutation grün. Die Prosa-Pins in `roadmap.md` Z. 13 und `observations/README.md`
  Z. 18 hält gar kein Gate; sie sind von Hand gegen den neuen Baum gelesen
  ([`BEO-ALL/pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md), ohne neuen
  Beleg: nichts war falsch, nur ungeschützt).
- **Push-Pflicht:** Die Commits `04dca313` (Pin-Hebung, `conventions.md:186` zeigt auf `conventions/done/` vor dem Move:
  `target-missing`) und `d67b6221` (der reine Move von [`MR-080`](../../../../harness/conventions.md#mr-080), seine
  relativen Anker auf den Konventions-Index lösen nicht mehr auf) sind **einzeln rot** — der Zwischenstand, den `AGENTS.md` §3.3 für
  den Lifecycle-Move zulässt, solange er nicht die Spitze eines Push wird. Der Verifier maß je Commit im Klon:
  `79a14c01` grün, `04dca313` ein Befund, `d67b6221` rot, `3d8f3555` und alles danach grün. **Alle Commits bis
  mindestens `0985933f` gehören in einen Push**; `main` steht vor `origin/main`, nichts ist gepusht. Nach dem Push
  `make ci` beziehungsweise die CI lesen, nicht aus dem lokalen Stand schließen.
- **Folge-Slices:** keine neuen. [`slice-296`](../open/slice-296-baseline-v700-neuerungen-sichtung.md),
  [`slice-297`](slice-297-spezifikation-heisst-specification.md) und
  [`slice-298`](../next/slice-298-neue-namen-wave-archive-anker.md) bestehen schon in `open/`; Risiko 5 geht an
  `slice-298`. Für die Prüflücke von `codepaths` (Lernsignal 2) ist **kein** Slice geschnitten — erst die
  Antwort des Auftraggebers auf die Spec-Lücken-Frage entscheidet, ob einer entsteht.
- **Risiken aus §6:** Risiko 1 *entfallen*, 2 *entfallen*, 3 *entfallen*, 4 *eingetreten* (Beleg: die Evidenz-Datei
  oben), 5 *weiter offen* → `slice-298` — Begründung je Risiko in §6.
- **Die drei Paarungen**, hier vorab geprüft (getragen von der Closure von welle-98): (a) Anker — kein
  Steering-Loop-Eintrag mit `liegt in`-Feld, keine Paarung auszulösen; (b) Folge-Slice — `slice-298` besteht als Datei
  in `open/`; (c) Register — beide genannten Verzeichnisse existieren und tragen nicht leere `evidence/`.

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

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md:416-417 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Die Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md) führt **zwei** Zeilen: den Default `*` (`ALL`) und
`tools/harness/` (`HARN`), beide Greenfield. Berührt ist der Default `*`: `.harness/` (Baum und Skills), `harness/`,
`.claude/`, `docs/`, `spec/` und `tools/archive-wave/` tragen keine eigene Deklaration (gemessen an der Tabelle
§Modus-Deklaration pro Sub-Area). `tools/harness/` bleibt **unberührt** — `fetch-baseline-cache.sh` wird gelaufen, nicht
geändert (anders als in [`slice-277`](../done/wellenlos/slice-277-baseline-v6181-bump.md), wo die Under-Copy-Barriere
einen Eingriff verlangte); wird er nötig, ist das eine Plan-Änderung vor dem Code, und die Sub-Area `HARN` kommt
hinzu.

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md:422-422 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Stand 2026-11-25,
`ls docs/plan/planning/observations/BEO-ALL/`; Zähler = Dateien unter `evidence/`), Sub-Area `*`. Treffer:
[`BEO-ALL/pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
(kein formgültiger Ausgang, 8× — Risiko 1: die Spiegel-Messung nimmt nur `.harness/baseline/` aus),
[`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
(verkörpert in [`MR-070`](../../../../harness/conventions.md#mr-070), 5× — Risiko 2),
[`BEO-ALL/citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
(verkörpert, 20× — Risiko 3: sieben Spannen verschieben sich),
[`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
(verkörpert, 10× — Risiko 4 und die Zählform in §3),
[`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md)
(verkörpert, 19× — das Delta-Inventar im Eintrag wird am vendorten Baum gezählt, nicht aus dem Scratchpad
übernommen),
[`BEO-ALL/path-scoped-commit-carries-staged-rest`](../observations/BEO-ALL/path-scoped-commit-carries-staged-rest/observation.md)
(verkörpert, 5× — `git diff --cached --stat` vor jedem Commit, §3),
[`BEO-ALL/pipe-swallows-gate-exit-code`](../observations/BEO-ALL/pipe-swallows-gate-exit-code/observation.md)
(verkörpert, 5× — kein Gate-Lauf hinter einer Pipe),
[`BEO-ALL/form-vom-nachbarn-statt-von-der-vorlage`](../observations/BEO-ALL/form-vom-nachbarn-statt-von-der-vorlage/observation.md)
(offen, 2× — der Bump-Eintrag entsteht per `cp` aus der Vorlage des **neuen** Baums, nicht nach [`MR-080`](../../../../harness/conventions.md#mr-080)),
[`BEO-ALL/commit-message-overclaims-work`](../observations/BEO-ALL/commit-message-overclaims-work/observation.md)
(verkörpert, 26× — die Botschaft des Hauptcommits behauptet nur, was der Diff trägt). **Kein Treffer:** keine
Beobachtung zu Vorlagen-Umbenennungen. Der Zähler folgt aus den Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-11-25: `image-scan` **grün** (jüngster Lauf 2026-11-25T12:59:28Z). `upstream-drift`
**rot** (jüngster Lauf 2026-11-25T09:26:13Z); der Log war über `gh run view` nicht abrufbar (leere Ausgabe).
Gelesen ist stattdessen `make baseline-freshness` am selben Tag: Exit 3, „NEUER RELEASE verfügbar (Pin v6.18.1):
v7.0.0“, Content „upstream unverändert“ — das ist der Gegenstand dieses Slice, **planmäßig gelesen**. Ob der Lauf
zugleich veraltete Tool-Pins meldet, ist **nicht** gelesen und nicht Gegenstand (eigene Vorgänge).

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
