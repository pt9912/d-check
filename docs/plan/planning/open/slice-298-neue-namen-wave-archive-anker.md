# Slice slice-298: Neue Namen — `wave-<NN>`, `archive.zip`, `seit wave-<NN>` werden erkannt und erzeugt, der Bestand bleibt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-98-baseline-v700`](../welle-98-baseline-v700.md).

**Bezug:** [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(die Wellen-Invariante liest ein literales Präfix; die Defaults bleiben),
[`ADR-0055`](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md) (`wave-drift` und die Register-Aussagen),
[`ADR-0117`](../../adr/0117-wellen-namen-ueber-die-konfiguration.md) (die Senkung von `planning`, `Proposed` bis zur Closure dieses Slice),
[`MR-069`](../../../../harness/conventions.md#mr-069) (eine Senkung wird benannt),
[`MR-050`](../../../../harness/conventions.md#mr-050) (der Herkunfts-Anker `seit <welle>`),
[`MR-053`](../../../../harness/conventions.md#mr-053) und
[`MR-054`](../../../../harness/conventions.md#mr-054) (Vorprüfungen in den Slice-Plänen),
[`slice-297`](../done/slice-297-spezifikation-heisst-specification.md) (die andere Hälfte der Namensregel).
Anlass: Entscheid 1 der [Welle](../welle-98-baseline-v700.md) — Baseline-Regelwerk `grundlagen-harness-dateien.md`
§Verzeichniskonvention, Absatz *Dateinamen sind englisch*, seit `v7.0.0`.

**Berührte Spec-Stellen:** — (die Defaults von `planning.waves` in der Spezifikation und im Lastenheft bleiben; der Slice
stellt die **Konfiguration dieses Repos** um, nicht das Produkt-Angebot).

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

**Ziel:** Eine neue Welle heißt `wave-<NN>-<titel>.md` mit `wave-<NN>-results.md`, ihr Archiv `archive.zip`, ihr
Herkunfts-Anker `seit wave-<NN>` — das Werkzeug `tools/archive-wave`, die Sensoren in `.d-check.yml` und
`.d-check.closure.yml` und die Command-Dateien erkennen und erzeugen sie, der Bestand (`welle-<NN>`, `archiv.zip`,
`seit welle-<NN>`) bleibt unberührt, und die Koexistenz beider Schreibweisen ist ein Eintrag in `harness/conventions/`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein Bestand wird umbenannt** (Entscheid 1, „Bestand bleibt“): 344 Pfade mit `welle-<NN>`, 188 `archiv.zip`, 24
  `seit welle-<NN>`, 38 Register-Zeilen — jede Umbenennung wäre ein Beleg-Eingriff und ein eigener Beschluss.
- **Kein neues Produkt-Merkmal.** Das Modul `planning` liest **ein** Präfix; es mehrere lesen zu lassen ist die größere
  Antwort von [welle-98 §5](../welle-98-baseline-v700.md) Frage 3 (Lastenheft, Spezifikation, ADR, Handbuch, Release) und
  gehört dem Auftraggeber. Dieser Slice stellt **die Konfiguration** um und benennt den Preis.
- **Keine Änderung der Defaults** des Werkzeugs (`welle-*.md`, `welle-*-results.md`): Konsumenten verlassen sich auf sie;
  ob sie auf `wave-*` wechseln, ist eine Zeile der Sichtung
  ([`slice-296`](../open/slice-296-baseline-v700-neuerungen-sichtung.md), Frage B).
- **Die Prosa bleibt deutsch:** das Wort *Welle*, das Slice-Feld `**Welle:**`, die Roadmap-Überschriften (`## Offene Wellen`
  …), die Grund-Codes (`wave-drift` heißt schon so).
- **`welle-98` zieht um** (neu und ungepusht; Verdikt Frage 2): die einzige Datei, deren Namen dieser Slice ändert.
- **Keine weitere Neuerung** (Entscheid 2), insbesondere nicht `.claude/worktrees/**` im Scan-Ignore und kein Vorlagen-Paar.
- **Kein Release.** Das Werkzeug `tools/archive-wave` ist kein Produkt-Bestandteil des Images
  (`make archive-wave-test` ist nicht Teil von `make gates`); sein Verhalten für Bestand bleibt byte-gleich.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Werkzeug:** `tools/archive-wave` liest `wave-<NN>` im `**Welle:**`-Feld der Slices und `-welle=wave-<NN>`
      **neben** `welle-<NN>` und schreibt für eine neue Welle `done/<wave-id>/archive.zip`; der Slice- und der Review-Modus schreiben
      `-archive.zip` (Verdikt Frage 8); ein Bestands-Lauf (`make archive-wave WELLE=welle-97`, Dry-Run) ist vorher und nachher
      **byte-gleich**; jeder neue Test läuft ohne die Änderung **aus dem richtigen Grund** rot (die Meldung ist gelesen, §3);
      `make archive-wave-test` grün.
- [ ] **Sensoren:** `planning.waves.glob` und `results-glob` stehen auf `wave-*`, die `matrix`-Klasse `wave` mit ihren drei
      Verboten steht **neben** `welle`, beide `exempt-paths` und das Closure-Profil tragen die neue Form; an der Umstellung
      steht die **Grenze** im Kommentar (der Bestand ist ungewächtert), und die Probe aus
      [welle-98 §3 (g)](../welle-98-baseline-v700.md) läuft in einer Wegwerf-Kopie grün — `wave-drift` ohne Zeiger,
      Exit 0 mit Zeiger. `make doc-check` und `make planning-check` melden 0 Befunde.
- [ ] **Commands, Briefing, Eintrag:** `plan-welle.md`, `close-welle.md`, `implement-slice.md` und `tools/archive-wave/README.md`
      nennen für Neues `wave-<NN>`, `wave-<NN>-results.md`, `archive.zip`, `seit wave-<NN>` (und für den Bestand weiter die
      alten Formen); `AGENTS.md` §3.7 nennt beide Anker-Formen; die Sensor-Dateien `harness/sensors/planning-check.md` und
      `harness/sensors/archive-wave.md` tragen die Grenze; der **Koexistenz-Eintrag** steht unter `harness/conventions/` in der
      Form, die das Architect-Verdikt festlegt (Frage 1), im Index; die Senkung von `planning` trägt
      [`ADR-0117`](../../adr/0117-wellen-namen-ueber-die-konfiguration.md), der Eintrag zeigt auf sie; sie geht mit der Closure nach `Accepted`.
- [ ] `make gates` grün; `make archive-wave-test` grün.
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

**Die Zählung nennt ihre Form** ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 14): gezählt wird die **Zeile**, die ein
**Muster** der alten Form trägt (`welle-<…>`, `welle-*`, `archiv.zip`, `seit welle-<…>`), nicht eine Instanz
(`welle-97`); der Bestand der Instanzen steht in [welle-98 §1](../welle-98-baseline-v700.md). Stand `590dde3d`; der
Implementer **wiederholt** die Zählung vor der ersten Änderung.

| Datei / Komponente | Änderungs-Art | Begründung und Zahl |
|---|---|---|
| `tools/archive-wave/collect.go` | update | **2** Stellen: `welleIDInFieldRE` (`\bwelle-(\d+)\b`, Z. 16) und der Vergleich `"welle-"+m[1] == welleID` (Z. 124). Ohne die Änderung endet `make archive-wave WELLE=wave-99` mit Exit 2 („weder Welle-Plan noch Slices … gefunden“, gemessen im Probelauf, fail-closed) |
| `tools/archive-wave/archive.go`, `stub.go`, `main.go` | update | `archiv.zip` an **14** Stellen (`archive.go` 7, `stub.go` 6, `main.go` 1: der Dateiname und der `unzip -p`-Hinweis im Stub). Der Dateiname hängt an der Welle: Bestand schreibt weiter `archiv.zip`, Neues `archive.zip` — **die Schreibweise folgt dem Präfix der Kennung**, nicht einem Schalter |
| `tools/archive-wave/*_test.go` | update | **194** Zeilen mit `welle` in 7 Dateien, **16** mit `archiv`: je Fall mit `welle-<NN>` ein Gegenstück mit `wave-<NN>` (Happy, Bestand unverändert, gemischt: eine `welle-`- und eine `wave-`-Welle im selben Repo — der Fall, der **nur** unter Koexistenz entsteht). Jeder neue Test wird ohne die Änderung rot gefahren (Mutation, §6) |
| `tools/archive-wave/README.md`, `Makefile` (Werkzeug), `Makefile` (Wurzel, Target `archive-wave`) | update | **11** Zeilen (Wurzel-`Makefile` 1, `tools/archive-wave/Makefile` 3, `README.md` 7; gezählt mit `git grep -c 'welle-NN\|<welle-id>\|archiv\.zip\|welle-42\|-welle='`): Aufruf-Beispiele (`WELLE=welle-NN`), `<welle-id>/archiv.zip`; die Hilfetexte nennen beide Formen |
| [`.d-check.yml`](../../../../.d-check.yml) | update | **10** Zeilen und **3** Kommentarzeilen: `planning.waves` (Z. 915ff.: `glob`, `results-glob`, **2**), `matrix`-Klasse `wave` neben `welle` (Z. 518–530: **3**: `name`, `paths`, `token`), drei Verbote `from: spec-straten / sicht / adr → to: wave` (Z. 580, 589, 594: **3**), `exempt-paths` `done/wave-*/*.md` (Z. 410, 418: **2**), dazu **3** Kommentarzeilen mit der **Grenze** an der Umstellung. Im Probelauf: mit `wave-*` meldet `planning-check` `wave-drift` für die unregistrierte Datei; mit `welle-*` (heute) ist sie unsichtbar |
| [`.d-check.closure.yml`](../../../../.d-check.closure.yml) | update | **3** Zeilen: eine `structure`-Regel `files: "docs/plan/planning/done/wave-*-results.md"` mit `section-pattern: '^## Was wurde geliefert'` und `non-empty: true` neben der Regel für `welle-*-results.md` (Z. 268) |
| `.claude/commands/plan-welle.md`, `close-welle.md`, `implement-slice.md` | update | **13** Muster-Zeilen (`welle-<…>`, `<welle-id>`, `welle-NN`, `archiv.zip`, `seit welle-<…>`; gezählt mit `git grep -n 'welle-<\|welle-NN\|<welle-id>\|welle-\*\|archiv\.zip' -- .claude/commands`); die Vorlagen-Namen zog [`slice-295`](../done/slice-295-baseline-v700-bump.md) schon nach |
| `AGENTS.md` | update | **1** Zeile (§3.7, Z. 230: Herkunft „nur als **ein** auflösbares Feld … `seit welle-<NN>`“) nennt beide Formen; die Instanzen (`seit welle-73`, `seit welle-80`) bleiben |
| `harness/sensors/planning-check.md`, `harness/sensors/archive-wave.md` | update | die **Grenze** der Umstellung (der Bestand ist ungewächtert) — [`AGENTS.md`](../../../../AGENTS.md) §3.8: ein Modul verspricht nur über das, was es scannt |
| neue Eintrags-Datei unter `harness/conventions/` (Koexistenz) | neu | aus `MR-<NNN>-title.template.md` des **vendorten** Baums (`cp`, `diff -q`); **Felder, Nummer und Trigger legt das [Verdikt](../../../reviews/2026-11-26-welle-98-architect-verdikt.md) fest (Frage 1)**, der **Implementer** schreibt sie in diesem Slice mit den dann gemessenen Zahlen; Index in `harness/conventions.md` |
| [`ADR-0117`](../../adr/0117-wellen-namen-ueber-die-konfiguration.md) und `docs/plan/adr/README.md` | update (Status) | `**Status:**` `Proposed` → `Accepted` bei der Closure, Geschichte-Zeile, Index-Zeile; der Status-Übergang ist dem ADR-Gate erlaubt |
| `docs/plan/planning/README.md` | prüfen | nennt „flache Welle-Dateien“ — prüfen, ob es einen Namen trägt (Prosa bleibt) |

**Plan-Änderung (Planner, vor dem Übergang `open` → `next`; Stand `439a49ff`, die Zählungen der Tabelle sind wiederholt):**
Die Wiederholung nennt dieselben Zahlen wie die Tabelle bis auf drei Stellen, die nachweislich nicht stimmen; die Tabelle bleibt
als Plan-Stand stehen, diese Zeilen überschreiben sie ausdrücklich.

| Stelle | Tabelle sagt | gemessen | Kommando |
|---|---|---|---|
| `archiv.zip` im Werkzeug | **14** Stellen (`archive.go` 7, `stub.go` 6, `main.go` 1) | **12** (`archive.go` 7, `stub.go` **4**, `main.go` 1); dieselbe Zahl schon am Stand `590dde3d` | `git grep -c 'archiv\.zip' -- tools/archive-wave/archive.go tools/archive-wave/stub.go tools/archive-wave/main.go` |
| Zeilen in der [`.d-check.yml`](../../../../.d-check.yml) | Z. 915ff. (`glob`, `results-glob` als „update“), 518–530, 580/589/594, 410/418 | Klasse `welle` Z. 536–548, Verbote 598/607/612, `exempt-paths` 428/436 (je **+18**); der Block `waves:` (Z. 936) führt **weder `glob` noch `results-glob`**, nur `dir` und `mode` — die beiden Zeilen sind **neu**, die Defaults (`welle-*.md`) gelten heute; es bleiben 10 Zeilen, davon **8 geändert und 2 neu** | `git grep -n 'welle' -- .d-check.yml` und `awk '/^  waves:/,/mode: many/' .d-check.yml` |
| Risiko 5, die Treffermenge „nennt nur `planning_test.go`“ | ein Treffer | **fünf**: `planning_test.go`, [`ADR-0117`](../../adr/0117-wellen-namen-ueber-die-konfiguration.md) (Z. 34 und 77 tragen `wave-99`), die Welle, das Architect-Verdikt und dieser Plan; die Aussage galt am Stand `590dde3d`, vor dem Verdikt | `git grep -l -E 'wave-[0-9]{2,}' -- . ':!.harness/baseline'` |

**Folge für Risiko 5 und Commit 2:** In einer Wegwerf-Kopie (Klasse `wave`, Token `wave-\d{2,}`, Verbot `adr → wave`) meldet `doc-check`
**2** `matrix-forbidden`, beide in [`ADR-0117`](../../adr/0117-wellen-namen-ueber-die-konfiguration.md) (Z. 34, 77), sonst 0. Der Ausgang
„*entfallen*, wenn `make doc-check` 0 Befunde meldet“ trägt damit nur, wenn die ADR vorher kein `wave-<Ziffern>`-Beispiel mehr führt oder die
Nennung als Provenance deklariert ist. Beides ist eine **Entscheidung des Architect**, nicht des Implementers: Die ADR ist `Proposed` und
damit noch änderbar ([`AGENTS.md`](../../../../AGENTS.md) §3.5, nur `Accepted` ist immutabel), und sie trägt die Senkung, die dieser
Slice konfiguriert. Sie ist **Vorbedingung vor Commit 2** (die Konfiguration), nicht vor Commit 1 (das Werkzeug); liegt sie dort nicht vor, gilt
die Rückführung `in-progress` → `open` aus §4.

**Reihenfolge und Commit-Zerlegung** ([`AGENTS.md`](../../../../AGENTS.md) §3.3):

1. **Vorbedingung, kein Commit dieses Slice:** das Architect-Verdikt zu den Fragen der Welle liegt vor,
   [`slice-295`](../done/slice-295-baseline-v700-bump.md) (die Vorlagen heißen jetzt `wave.template.md`, …) und
   [`slice-297`](../done/slice-297-spezifikation-heisst-specification.md) sind `done/`.
2. **Commit 1: das Werkzeug** samt Tests (`make archive-wave-test`), unabhängig von der Konfiguration — es liest beide Formen,
   solange nichts umgestellt ist.
3. **Commit 2: der Name der Welle, die Konfiguration und die Sensoren** — **erst jetzt**, weil die Umstellung von `planning.waves` die
   `welle-`-Wächter des Bestands abschaltet (Probelauf Punkt 5 in [welle-98 §1](../welle-98-baseline-v700.md)). Verdikt Frage 2: `welle-98`
   heißt ab hier `wave-98`; **keine Lücke**, weil Move, Roadmap-Zeiger und Konfiguration im selben Arbeitsbaum stehen und der Hook ihn prüft.
   Zuerst der **reine** `git mv` der Welle-Datei (R100, [`AGENTS.md`](../../../../AGENTS.md) §3.3), dann der Nachzug **aller** Verweise
   auf den alten Namen (`git grep -c 'welle-98' -- . ':!.harness/baseline'` neu zählen: Roadmap-Zeiger, `**Welle:**`-Felder, Closure-Notizen und
   Reports von [`slice-295`](../done/slice-295-baseline-v700-bump.md) und [`slice-297`](../done/slice-297-spezifikation-heisst-specification.md), der Plan selbst; darin
   `welle-98-results.md` → `wave-98-results.md`) und die Umstellung der Konfiguration. Pfad-begrenzt committen (`BEO-ALL/path-scoped-commit-carries-staged-rest`).
4. **Commit 3: Commands, Briefing, Sensor-Dateien und der Koexistenz-Eintrag.**
5. Review (höchstens zwei Runden), Verifikation, Closure. `git commit -F <datei> -- <Pfade>`, vorher `git diff --cached --stat`.

**Gates:** `make archive-wave-test`, `make doc-check`, `make planning-check`, `make verify-closure-notes` (das Closure-Profil),
`make gates`; die **Probe** aus [welle-98 §3 (g)](../welle-98-baseline-v700.md) in einer Wegwerf-Kopie (nicht im Repo).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`slice-295`](../done/slice-295-baseline-v700-bump.md) und
[`slice-297`](../done/slice-297-spezifikation-heisst-specification.md) liegen in `done/`; `in-progress/` trägt keinen Slice
(nur die Roadmap); das **Architect-Verdikt** zu den zwölf Fragen in [welle-98 §5](../welle-98-baseline-v700.md) (für diesen
Slice tragend: 1, 2, 3, 8) liegt vor: [Verdikt](../../../reviews/2026-11-26-welle-98-architect-verdikt.md); Priorisierung gesetzt, `Verantwortlich:` beim Übergang
`open` → `next` gesetzt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Wiederholung der Zählung zeigt mehr Stellen, die ein Muster der
  alten Form **fest** tragen, als §3 nennt (etwa ein Skript oder ein Hook, der `welle-` in einem regulären Ausdruck führt), oder
  der Auftraggeber wählt die größere Antwort von Frage 3 (das Modul liest mehrere Präfixe) — dann ist das ein
  **Produkt-Slice** mit eigener Anforderung und eigener ADR, und dieser Slice schrumpft auf das Werkzeug und die Commands.
- `in-progress` → `open` (blockiert): Das Architect-Verdikt zu Frage 1 oder 3 steht aus, oder `make planning-check` zeigt
  nach der Umstellung einen Befund, der nicht aus dem Bestand-Verzicht folgt (etwa ein Roadmap-Zeiger, den die neue Kennung
  nicht trifft).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (§2) und `make gates` grün; Review- und Verifikations-Report liegen unter `docs/reviews/`;
Closure-Notiz mit Lerneintrag geschrieben. **Review-Vorgabe der Welle:** höchstens zwei Runden, eingearbeitet
werden nur HIGH und MEDIUM; jede Nachbesserung zitiert die Stelle, die sie trägt
([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 18). Was die **Welle** einlöst — die Probe der neuen Namen am Endstand und
der unveränderte Bestand — steht in [ihrem Closure-Trigger](../welle-98-baseline-v700.md) (g) und nicht in dieser DoD
([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 10).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Ein Sensor erkennt `wave-` nur auf einer Seite** — Negativtest nur gegen die neue Form, der Bestand läuft stumm weiter
  ([`BEO-ALL/negativtest-deckt-nur-eine-regelseite`](../observations/BEO-ALL/negativtest-deckt-nur-eine-regelseite/observation.md),
  1×). — **Ausgang:** *entfallen*, wenn die Closure-Notiz je Sensor **beide** Seiten nennt (neue Form gesehen, Bestand nach der
  Umstellung ausdrücklich ungewächtert) und die Probe sie in einer Wegwerf-Kopie umdreht; sonst *eingetreten* mit einer
  Evidenz-Datei.
- **Ein Test läuft grün, ohne die Regel zu fahren**
  ([`BEO-ALL/testfall-faehrt-seine-regel-nicht`](../observations/BEO-ALL/testfall-faehrt-seine-regel-nicht/observation.md),
  2×). — **Ausgang:** *entfallen*, wenn die Closure-Notiz je neuem Test die Mutation (die Änderung zurücknehmen) und die
  gelesene Meldung nennt und jeder rot war; sonst *eingetreten* (dann 3×).
- **Das Werkzeug ändert den Bestands-Lauf** — `make archive-wave WELLE=welle-97` liefert danach etwas anderes. — **Ausgang:**
  *entfallen*, wenn der Dry-Run vorher und nachher mit `diff` **byte-gleich** ist (stdout, stderr, Exit); sonst *eingetreten*,
  und die Änderung geht zurück.
- **Die Umstellung von `planning.waves` schaltet die Wächter des Bestands ab** (38 Register-Zeilen, 38 Ergebnisnotizen; gemessen
  im Probelauf) — das ist die **gewollte** Folge der kleinsten Antwort, keine Panne. — **Ausgang:** *weiter offen* → Register
  (neue Beobachtung mit Beleg `evidence/slice-298.md`: „eine Umstellung des Präfixes schaltet die Wächter der anderen
  Schreibweise ab“); die Senkung selbst und ihre Wiedervorlage „das Modul liest mehrere Präfixe“ trägt
  [`ADR-0117`](../../adr/0117-wellen-namen-ueber-die-konfiguration.md) ([welle-98 §5](../welle-98-baseline-v700.md) Frage 3, größere Antwort).
- **Ein `wave-<NN>`-Token trifft einen falschen Treffer** in Spec, Sicht oder ADR (die Klasse `wave` braucht Ziffern:
  `wave-drift` ist kein Treffer). — **Ausgang:** *entfallen*, wenn `make doc-check` nach dem Anlegen der Klasse 0 Befunde
  meldet — am Stand `590dde3d` trägt kein Dokument außerhalb dieser Pläne ein `wave-<NN>`-Token (`git grep -l -E
  'wave-[0-9]{2,}' -- . ':!.harness/baseline'` nennt nur `planning_test.go`); sonst *eingetreten*.

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

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md:416-417 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Die Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md) führt zwei Zeilen: den Default `*` (`ALL`) und
`tools/harness/` (`HARN`), beide Greenfield. Berührt ist der Default `*`: `tools/archive-wave/` (eigenes `go.mod`),
die Konfiguration im Repo-Wurzelverzeichnis, `.claude/commands/`, `harness/` und `AGENTS.md` tragen keine eigene
Deklaration (gemessen an der Tabelle §Modus-Deklaration pro Sub-Area). `tools/harness/` bleibt **unberührt** — kein
Skript dort trägt ein `welle-`-Muster im Code (`git grep -n -i 'welle' -- tools/harness` nennt nur Kommentarzeilen in
`slice-mv.sh`); wird das anders gemessen, ist es eine Plan-Änderung vor dem Code und `HARN` kommt hinzu.

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md:422-422 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Stand 2026-11-25,
`ls docs/plan/planning/observations/BEO-ALL/`; Zähler = Dateien unter `evidence/`), Sub-Area `*`. Treffer:
[`BEO-ALL/negativtest-deckt-nur-eine-regelseite`](../observations/BEO-ALL/negativtest-deckt-nur-eine-regelseite/observation.md)
(offen, 1× — Risiko 1: neue **und** alte Form),
[`BEO-ALL/testfall-faehrt-seine-regel-nicht`](../observations/BEO-ALL/testfall-faehrt-seine-regel-nicht/observation.md)
(offen, 2× — Risiko 2),
[`BEO-ALL/grenze-der-scan-menge-nur-bei-einem-modul-benannt`](../observations/BEO-ALL/grenze-der-scan-menge-nur-bei-einem-modul-benannt/observation.md)
(Grenze der Scan-Menge — die Sensor-Dateien tragen die Grenze der Umstellung, §3),
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
(verkörpert als Hard Rule in `AGENTS.md` §3.8, 4× — das Modul `planning` verspricht nur über das Präfix, das es liest),
[`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
(verkörpert in [`MR-070`](../../../../harness/conventions.md#mr-070), 5× — hier wird **nichts** mechanisch ersetzt, der Bestand bleibt),
[`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
(verkörpert, 10× — die Zählung in §3 zählt Muster, nicht Instanzen),
[`BEO-ALL/commit-message-overclaims-work`](../observations/BEO-ALL/commit-message-overclaims-work/observation.md)
(verkörpert, 26× — die Botschaft sagt, dass der Bestand ungewächtert ist). **Kein Treffer:** keine Beobachtung zu zwei
Schreibweisen derselben Kennung. Der Zähler folgt aus den Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-11-25: `image-scan` **grün** (jüngster Lauf 2026-11-25T12:59:28Z). `upstream-drift`
**rot** (jüngster Lauf 2026-11-25T09:26:13Z), der Log war über `gh run view` nicht abrufbar; gelesen ist
`make baseline-freshness` (Exit 3, neuer Release `v7.0.0`) — der Gegenstand der Welle, **planmäßig gelesen**. Zum
Zeitpunkt dieses Slice sind [`slice-295`](../done/slice-295-baseline-v700-bump.md) und
[`slice-297`](../done/slice-297-spezifikation-heisst-specification.md) `done/`; der Stand wird beim Start neu gelesen.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
