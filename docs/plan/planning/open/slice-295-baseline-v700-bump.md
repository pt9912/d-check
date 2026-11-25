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

**Verantwortlich:** —

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
(Delta vollständig nach Datei, **keine Neuerung übernommen**), dazu den Adaptions-Eintrag zum Dateinamen der
Spezifikation (Entscheid 1 der [Welle](../welle-98-baseline-v700.md)).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Keine Übernahme einer Neuerung.** Jede ist eine Konventions- oder Produkt-Entscheidung mit eigener
  Begründung (Entscheid 2); der Bump ist ein anderer Vorgang. Der Bump-Eintrag nennt das Delta vollständig
  nach Datei und markiert jede Bewegung als *nicht übernommen*; die Klassifikation mit Empfehlung liefert
  [`slice-296`](../open/slice-296-baseline-v700-neuerungen-sichtung.md), die Übernahme ein eigener Folge-Slice nach der
  Antwort des Auftraggebers. Das schließt die **Bezeichner** ein: `welle-`, `archiv.zip`, `seit welle-` bleiben
  in unseren Artefakten wie sie sind; nachgezogen werden nur die **Namen der Vorlagen-Dateien**, weil ein Verweis
  auf eine Datei, die es nicht mehr gibt, ein toter Verweis ist.
- **Keine Umbenennung von `spec/spezifikation.md`** (Entscheid 1). Der Eintrag hält die Abweichung fest; Pfade in
  `manual.go`, [`.d-check.yml`](../../../../.d-check.yml) und die Spec-Verweise bleiben.
- **Keine Änderung am Produkt.** Die eine Berührung ist die Pin-Nennung in `spec/spezifikation.md` (Zeile 8), die
  `manual.go` einbettet; sie ist in der Welle (d) benannt und gemessen. Quelltext, Handbuch, Standardwerte bleiben.
- **Frozen-Dateien bleiben byte-stabil** — `done/`-Slices und -Wellen, `docs/reviews/`, `CHANGELOG.md`,
  `conventions/done/`, eingehende CRs und der Kern der `Accepted`-ADRs nennen den alten Pin als
  Vergangenheits-Aussage (Bestand bewusst). Die Ausnahme ist allein die Antwort auf Frage 3 in
  [welle-98 §5](../welle-98-baseline-v700.md): sieben eingefrorene Dateien **verlinken** den alten Baum.
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

- [ ] **Vendoring:** `.harness/baseline/v7.0.0/` steht über `tools/harness/fetch-baseline-cache.sh v7.0.0`
      (56 Dateien, 55 im Manifest), `make baseline-verify` Exit 0, der Baum `v6.18.1` (55 Dateien) ist entfernt;
      der Pin steht in `harness/conventions.md` §Baseline und §Adoptierte Konventions-Quellen.
- [ ] **Spiegel:** Jeder lebende pin-gebundene Verweis nennt `v7.0.0` (die Klassen und ihre Zahlen in §3), die
      7 Symlinks unter `.claude/rules/` lösen auf den neuen Baum auf, jede `d-check:cite`-Spanne steht gegen den
      **Wortlaut** des neuen Baums (die 7 verschobenen gelesen, nicht nur verschoben), die Vorlagen-Namen in den
      Verweisen sind die neuen, und `make doc-check` meldet 0 Befunde.
- [ ] **Einträge:** Der Bump-Eintrag steht unter `harness/conventions/` (Delta **vollständig nach Datei**,
      Frozen-Klassen am echten Vorzustand gemessen, Zitat-Delta nach
      [`MR-039`](../../../../harness/conventions.md#mr-039), jede Bewegung als *nicht übernommen* markiert);
      [`MR-080`](../../../../harness/conventions.md#mr-080) liegt per reinem `git mv` in `conventions/done/`; der
      Adaptions-Eintrag zum Dateinamen der Spezifikation steht in der Form, die das Architect-Verdikt festlegt;
      beide Einträge stehen im Index.
- [ ] `make gates` grün; `make archive-wave-test` grün (nur, wenn Frage 7 der Welle den Kommentar-Nachzug in
      `tools/archive-wave/stub.go` bestätigt).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

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
| neue Eintrags-Datei unter `harness/conventions/` (Bump) | neu | aus `MR-<NNN>-title.template.md` des **neuen** Baums, per `cp` ([`BEO-ALL/form-vom-nachbarn-statt-von-der-vorlage`](../observations/BEO-ALL/form-vom-nachbarn-statt-von-der-vorlage/observation.md)); Nummer vergibt der Architect (höchste vergebene: [`MR-080`](../../../../harness/conventions.md#mr-080)); Delta nach Datei: **31** inhaltlich verschieden + **7** umbenannt + **1** neu, die **7** Stempel-Dateien in einer Zeile |
| neue Eintrags-Datei unter `harness/conventions/` (Dateiname) | neu | Form und Umfang aus dem Verdikt (Frage 1 und 2 der Welle); der **Architect** schreibt sie, nach dem Re-Vendor-Commit |
| die Eintrags-Datei von [`MR-080`](../../../../harness/conventions.md#mr-080) unter `harness/conventions/` | move | reiner `git mv` nach `conventions/done/`, dann die Link-Tiefen als eigener Commit ([`AGENTS.md`](../../../../AGENTS.md) §3.3) |
| 35 aktive `MR-*`-Dateien unter `harness/conventions/` | update | Tag-Tausch je Datei; die Mischzeilen **einzeln gelesen**: [`MR-021`](../../../../harness/conventions.md#mr-021) Z. 20 („aktuell …/v6.18.1/…“), [`MR-049`](../../../../harness/conventions.md#mr-049) Z. 14, [`MR-053`](../../../../harness/conventions.md#mr-053) Z. 9 |
| `AGENTS.md`, `harness/README.md`, 3 Dateien unter `harness/rules/`, 2 unter `.harness/skills/`, `.claude/agents/reviewer.md` | update | 8 Dateien; `AGENTS.md` Z. 345 („Baseline-Form `v6.18.1` §1“) ist eine Mischzeile und wird gegen die Gliederung der neuen `slice.template.md` gelesen (die Abschnitte §1 bis §8 bleiben, gemessen) |
| `docs/plan/planning/README.md`, `in-progress/roadmap.md`, `observations/README.md` | update | 3 Dateien; `roadmap.md` Z. 13 und `observations/README.md` Z. 18 sind Prosa-Pins, einzeln gelesen |
| `spec/spezifikation.md`, `spec/architecture.md` | update | 2 Dateien, je die **Rolle**-Zeile; die Spezifikation ist eingebettet (Welle (d)) |
| 15 lebende `d-check:cite`-Spannen | neu ankern | `.harness/skills/reviewer.md` 3, `closure-note-reviewer.md` 3, [`MR-031`](../../../../harness/conventions.md#mr-031) 3, [`MR-049`](../../../../harness/conventions.md#mr-049) 2, [`MR-005`](../../../../harness/conventions.md#mr-005) 1, [`MR-035`](../../../../harness/conventions.md#mr-035) 1, [`MR-039`](../../../../harness/conventions.md#mr-039) 1, [`MR-043`](../../../../harness/conventions.md#mr-043) 1. **7 verschieben sich** (gemessen am Vergleich derselben Zeilenbereiche beider Bäume): `grundlagen-begriffe` 50, `grundlagen-harness-dateien` 196, `modul-05` 170 und 180–181, `modul-10` 99, `reviewer.template` 87–88 und 99–100. Die 8 übrigen stehen auf unverschobenen Zeilen. **Plus 4** in den zwei Plänen dieser Welle (je Slice zwei, `modul-05` 373–374 und 379; sie verschieben sich ebenfalls). Zwei weitere Spannen unter `conventions/done/` ([`MR-057`](../../../../harness/conventions.md#mr-057), [`MR-058`](../../../../harness/conventions.md#mr-058)) zielen auf `v5.18.0` und liegen außerhalb des Scopes von `citations` |
| `.claude/commands/plan-welle.md` (2 Zeilen), `close-welle.md` (4 Zeilen) | update | die vier umbenannten Vorlagen: `welle.template.md`→`wave.template.md`, `welle-results.template.md`→`wave-results.template.md`, `archiv-stub-slice`/`-welle`→`archive-stub-slice`/`-wave`; **9 Zeilen in 4 Dateien** mit `stub.go` und `conventions.md` Z. 59 (`git grep -n "welle\.template\|welle-results\.template\|archiv-stub-\(slice\|welle\)\.template\|titel\.template" -- . ':!.harness/baseline' ':!docs/plan/planning' ':!docs/reviews' ':!harness/conventions/done'`; die Treffer in `conventions/done/` und in den Plänen nennen die Namen als Vergangenheits-Aussage und bleiben) |
| `tools/archive-wave/stub.go` | update (Kommentar) | zwei Kommentarzeilen nennen `archiv-stub-slice.template.md`/`-welle.template.md`; nur mit Bestätigung der Frage 7 der Welle; kein Verhalten |
| [`.d-check.yml`](../../../../.d-check.yml) | update | nur, was die Antwort auf Frage 3 verlangt: je eingefrorener Datei mit Link auf den alten Baum **ein** `ignore-refs`-Eintrag (gemessen **7**); heute 30 Einträge über 12 Tags |
| dieser Plan und [`slice-296`](../open/slice-296-baseline-v700-neuerungen-sichtung.md) | update | die je zwei `d-check:cite`-Spannen in §8 auf den neuen Baum |

**Reihenfolge und Commit-Zerlegung** ([`AGENTS.md`](../../../../AGENTS.md) §3.3, [`BEO-ALL/path-scoped-commit-carries-staged-rest`](../observations/BEO-ALL/path-scoped-commit-carries-staged-rest/observation.md)):

1. **Vorbedingung, kein Commit dieses Slice:** das Architect-Verdikt zu den Fragen der Welle liegt vor.
2. **Spiegel-Messung vor dem ersten Tag-Tausch:** `git grep -l 'v6\.18\.1' -- . ':!.harness/baseline'` über das **ganze**
   Repo (nicht `.harness/`: dort liegen die Skills), geteilt in Frozen- und lebende Klassen **über die Eigenschaft**
   (liegt in `done/`, `docs/reviews/`, `conventions/done/`, ist eine `Accepted`-ADR, ist `CHANGELOG.md` oder ein
   eingehender CR), nicht über Verzeichnisnamen. Erwartung: 72 Dateien außerhalb des Baums, davon 17 Frozen (8 in
   `done/`, 8 Reviews, 1 `CHANGELOG.md`) und 55 weitere, von denen 5 `Accepted`-ADRs sind.
3. **Ein Hauptcommit:** Baum-Tausch, Pin, Tag-Tausch, Symlinks, Zitat-Spannen, Vorlagen-Namen, `ignore-refs`, **und** die
   zwei neuen Einträge samt Index — **nicht** in mehreren Commits, weil der lokale `pre-commit`-Hook den Arbeitsbaum
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
zu den sieben Fragen in [welle-98 §5](../welle-98-baseline-v700.md) liegt als Datei unter `docs/reviews/` vor;
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
  leer), sieben Vorlagen tragen einen neuen Namen. — **Ausgang:** *entfallen*, wenn `make doc-check` nach dem
  Entfernen des alten Baums außer den in §3 benannten Folgen (Pfad-Links, Zitat-Spannen, die sieben Frozen-Links)
  nichts meldet; sonst *eingetreten* mit einer Evidenz-Datei zu
  [`BEO-ALL/pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md).
- **Eine Mischzeile wird pauschal ersetzt und fälscht eine Vergangenheits-Aussage**
  ([`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md),
  5×). — **Ausgang:** *entfallen*, wenn jede geänderte Zeile ohne Pfad- oder URL-Form im Bump-Eintrag einzeln
  genannt ist und `git diff --stat <Vor-Bump>..HEAD -- docs/plan/planning/done docs/reviews CHANGELOG.md harness/conventions/done`
  leer ist (bis auf die Ausnahme, die die Antwort auf Frage 3 festlegt); sonst *eingetreten*.
- **Eine verschobene Zitat-Spanne wird nur verschoben, nicht gelesen**
  ([`BEO-ALL/citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md),
  20×). — **Ausgang:** *entfallen*, wenn der Bump-Eintrag je verschobener Spanne alte und neue Zeile und den
  Befund „Wortlaut gleich“ oder „Wortlaut geändert, Spanne angepasst“ nennt; sonst *eingetreten* mit einer
  Evidenz-Datei.
- **Eine Zahl im Bump-Eintrag stimmt nicht** — Zählung am Proxy statt am Gegenstand
  ([`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md),
  10×; in [`MR-080`](../../../../harness/conventions.md#mr-080) waren es drei). — **Ausgang:** *entfallen*, wenn
  Reviewer und Verifier jede Zahl mit dem im Eintrag genannten Kommando nachzählen und keine abweicht; sonst
  *eingetreten* mit einer Evidenz-Datei.
- **Die kopierte Welle-Vorlage trägt nach dem Bump `wave-<Kennung>-results.md`** in ihrer `Lifecycle:`-Note, und
  `plan-welle.md` verlangt, diese Note zu behalten — dieses Repo führt aber `welle-NN-results.md`. Das ist eine
  **bekannte Folge** der nicht übernommenen Neuerung (§1), kein Fehler des Bumps. — **Ausgang:** *weiter offen* →
  Punkt 1 der Entscheidungsvorlage in [`slice-296`](../open/slice-296-baseline-v700-neuerungen-sichtung.md); bis zur Antwort
  ersetzt der Planner den Namen beim Füllen.

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

*(Wird bei der Closure gefüllt.)*

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Die Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md) führt **zwei** Zeilen: den Default `*` (`ALL`) und
`tools/harness/` (`HARN`), beide Greenfield. Berührt ist der Default `*`: `.harness/` (Baum und Skills), `harness/`,
`.claude/`, `docs/`, `spec/` und `tools/archive-wave/` tragen keine eigene Deklaration (gemessen an der Tabelle
§Modus-Deklaration pro Sub-Area). `tools/harness/` bleibt **unberührt** — `fetch-baseline-cache.sh` wird gelaufen, nicht
geändert (anders als in [`slice-277`](../done/wellenlos/slice-277-baseline-v6181-bump.md), wo die Under-Copy-Barriere
einen Eingriff verlangte); wird er nötig, ist das eine Plan-Änderung vor dem Code, und die Sub-Area `HARN` kommt
hinzu.

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

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
