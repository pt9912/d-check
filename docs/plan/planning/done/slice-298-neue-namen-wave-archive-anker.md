# Slice slice-298: Neue Namen — `wave-<NN>`, `archive.zip`, `seit wave-<NN>` werden erkannt und erzeugt, der Bestand bleibt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`wave-98-baseline-v700`](../wave-98-baseline-v700.md).

**Bezug:** [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(die Wellen-Invariante liest ein literales Präfix; die Defaults bleiben),
[`ADR-0055`](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md) (`wave-drift` und die Register-Aussagen),
[`ADR-0117`](../../adr/0117-wellen-namen-ueber-die-konfiguration.md) (die Senkung von `planning`, `Proposed` bis zur Closure dieses Slice),
[`MR-069`](../../../../harness/conventions.md#mr-069) (eine Senkung wird benannt),
[`MR-050`](../../../../harness/conventions.md#mr-050) (der Herkunfts-Anker `seit <welle>`),
[`MR-053`](../../../../harness/conventions.md#mr-053) und
[`MR-054`](../../../../harness/conventions.md#mr-054) (Vorprüfungen in den Slice-Plänen),
[`slice-297`](../done/slice-297-spezifikation-heisst-specification.md) (die andere Hälfte der Namensregel).
Anlass: Entscheid 1 der [Welle](../wave-98-baseline-v700.md) — Baseline-Regelwerk `grundlagen-harness-dateien.md`
§Verzeichniskonvention, Absatz *Dateinamen sind englisch*, seit `v7.0.0`.

**Berührte Spec-Stellen:** — (die Defaults von `planning.waves` in der Spezifikation und im Lastenheft bleiben; der Slice
stellt die **Konfiguration dieses Repos** um, nicht das Produkt-Angebot).

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

**Ziel:** Eine neue Welle heißt `wave-<NN>-<titel>.md` mit `wave-<NN>-results.md`, ihr Archiv `archive.zip`, ihr
Herkunfts-Anker `seit wave-<NN>` — das Werkzeug `tools/archive-wave`, die Sensoren in `.d-check.yml` und
`.d-check.closure.yml` und die Command-Dateien erkennen und erzeugen sie, der Bestand (`welle-<NN>`, `archiv.zip`,
`seit welle-<NN>`) bleibt unberührt, und die Koexistenz beider Schreibweisen ist ein Eintrag in `harness/conventions/`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein Bestand wird umbenannt** (Entscheid 1, „Bestand bleibt“): 344 Pfade mit `welle-<NN>`, 188 `archiv.zip`, 24
  `seit welle-<NN>`, 38 Register-Zeilen — jede Umbenennung wäre ein Beleg-Eingriff und ein eigener Beschluss.
- **Kein neues Produkt-Merkmal.** Das Modul `planning` liest **ein** Präfix; es mehrere lesen zu lassen ist die größere
  Antwort von [wave-98 §5](../wave-98-baseline-v700.md) Frage 3 (Lastenheft, Spezifikation, ADR, Handbuch, Release) und
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

- [x] **Werkzeug:** `tools/archive-wave` liest `wave-<NN>` im `**Welle:**`-Feld der Slices und `-welle=wave-<NN>`
      **neben** `welle-<NN>` und schreibt für eine neue Welle `done/<wave-id>/archive.zip`; der Slice- und der Review-Modus schreiben
      `-archive.zip` (Verdikt Frage 8); ein Bestands-Lauf (`make archive-wave WELLE=welle-97`, Dry-Run) ist vorher und nachher
      **byte-gleich**; jeder neue Test läuft ohne die Änderung **aus dem richtigen Grund** rot (die Meldung ist gelesen, §3);
      `make archive-wave-test` grün.
      **Vorbehalt zum Haken:** Die Klausel „jeder neue Test läuft ohne die Änderung rot“ ist **wörtlich nicht
      erfüllt**: unter der Gesamt-Rücknahme sind 9 der 14 neuen Tests rot und 5 grün; jeder der 14 ist gegen eine
      benannte andere Änderung rot gefahren (Verifikation V-2, §7). Der Rest des Punkts ist bestätigt.
- [x] **Sensoren:** `planning.waves.glob` und `results-glob` stehen auf `wave-*`, die `matrix`-Klasse `wave` mit ihren drei
      Verboten steht **neben** `welle`, beide `exempt-paths` und das Closure-Profil tragen die neue Form; an der Umstellung
      steht die **Grenze** im Kommentar (der Bestand ist ungewächtert), und die Probe aus
      [wave-98 §3 (g)](../wave-98-baseline-v700.md) läuft in einer Wegwerf-Kopie grün — `wave-drift` ohne Zeiger,
      Exit 0 mit Zeiger. `make doc-check` und `make planning-check` melden 0 Befunde.
- [x] **Commands, Briefing, Eintrag:** `plan-welle.md`, `close-welle.md`, `implement-slice.md` und `tools/archive-wave/README.md`
      nennen für Neues `wave-<NN>`, `wave-<NN>-results.md`, `archive.zip`, `seit wave-<NN>` (und für den Bestand weiter die
      alten Formen); `AGENTS.md` §3.7 nennt beide Anker-Formen; die Sensor-Dateien `harness/sensors/planning-check.md` und
      `harness/sensors/archive-wave.md` tragen die Grenze; der **Koexistenz-Eintrag** steht unter `harness/conventions/` in der
      Form, die das Architect-Verdikt festlegt (Frage 1), im Index; die Senkung von `planning` trägt
      [`ADR-0117`](../../adr/0117-wellen-namen-ueber-die-konfiguration.md), der Eintrag zeigt auf sie; sie geht mit der Closure nach `Accepted`.
- [x] `make gates` grün; `make archive-wave-test` grün.
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

**Die Zählung nennt ihre Form** ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 14): gezählt wird die **Zeile**, die ein
**Muster** der alten Form trägt (`welle-<…>`, `welle-*`, `archiv.zip`, `seit welle-<…>`), nicht eine Instanz
(`welle-97`); der Bestand der Instanzen steht in [wave-98 §1](../wave-98-baseline-v700.md). Stand `590dde3d`; der
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

**Plan-Änderung (Implementer, vor dem ersten Code-Commit; Stand `db329ff2`, alle Zahlen am Arbeitsbaum dieses Standes gezählt):**

1. **Die Vorbedingung vor Commit 2 ist erfüllt.** [`ADR-0117`](../../adr/0117-wellen-namen-ueber-die-konfiguration.md) trägt seit dem
   Nachtrag V2 des [Verdikts](../../../reviews/2026-11-26-welle-98-architect-verdikt.md) in Z. 34 und 77 den Platzhalter `wave-<NN>` statt einer
   Beispiel-Kennung (Commit `db329ff2`); die Rückführung `in-progress` → `open` aus §4 entfällt. Risiko 5 trägt damit wieder: die Treffermenge von
   `git grep -l -E 'wave-[0-9]{2,}' -- . ':!.harness/baseline'` ist **vier** Dateien (dieser Plan, die Welle, das Verdikt, `planning_test.go`) statt
   fünf, keine liegt in Spec, Sicht oder ADR; die Aussage „nennt nur `planning_test.go`“ aus §6 bleibt überholt, an ihre Stelle tritt die Zählung mit
   ihrem Kommando.
2. **Die Zählungen vor der ersten Änderung** (die Tabelle bleibt Plan-Stand, diese Zahlen überschreiben sie dort, wo sie abweichen):

   | Gegenstand | Tabelle / Welle sagt | gemessen | Kommando |
   |---|---|---|---|
   | `archiv.zip` im Werkzeug (Quelltext) | 12 (Plan-Änderung oben) | **12** (`archive.go` 7, `stub.go` 4, `main.go` 1) | `git grep -c 'archiv\.zip' -- tools/archive-wave/archive.go tools/archive-wave/stub.go tools/archive-wave/main.go` |
   | Zeilen mit `welle` in den Tests | 194 in 7 Dateien | **194 in 7 Dateien** (nur klein geschrieben: 149 in 6) | `git grep -i -c 'welle' -- 'tools/archive-wave/*_test.go'` |
   | Zeilen mit `archiv.zip` in den Tests | 16 | **16** (davon **8** in den Tests des Slice- und des Review-Modus, sie erwarten `-archiv.zip` und ändern sich mit dem Verdikt Frage 8) | `git grep -c 'archiv\.zip' -- 'tools/archive-wave/*_test.go'` |
   | Pfade mit `welle-<NN>` im Namen | 344 | **346** (dasselbe Kommando zählt am Stand `590dde3d` 345; das Verdikt vom 2026-11-26 kam hinzu) | `git ls-files \| grep -c 'welle-[0-9]'` |
   | `archiv.zip` im Bestand | 188 | **188** | `git ls-files \| grep -c 'archiv\.zip$'` |
   | Register-Zeilen `welle-` | 38 | **38** | `grep -c '^\| welle-' docs/plan/planning/in-progress/roadmap.md` |
   | `seit welle-` im lebenden Bestand | 24 | **36** Zeilen, davon **17** die Platzhalter-Form `seit welle-<NN>` (Hard-Rule-Texte, Commands) und **19** mit Ziffern (1 × `seit welle-98`, der Anker an der Ventil-Zeile der `.d-check.yml`); die 24 reproduziert **kein** Kommando, die 36 sind die gezählte Form | `git grep -h -o 'seit welle-[0-9]*' -- . ':!.harness/baseline' ':!docs/plan/planning/done' ':!docs/reviews' \| wc -l` |
   | Nennungen von `welle-98` (alle Formen) | „mehr als die 14 Zeilen“ (Verdikt Frage 2) | **71 Zeilen in 13 Dateien**, **100 Nennungen** (33 × Welle-Datei, 5 × Link-Text ohne `.md`, 5 × Ergebnisnotiz, 13 × Dateiname des Verdikts, 44 × der bare Name) | `git grep -c 'welle-98' db329ff2 -- . ':!.harness/baseline'` und `git grep -o -h 'welle-98[A-Za-z0-9_.-]*' db329ff2 -- . ':!.harness/baseline' \| sed 's/^db329ff2://' \| sort \| uniq -c` |

3. **Der Nachzug von `welle-98` ist nicht mechanisch** ([`MR-070`](../../../../harness/conventions.md#mr-070)): ein Tausch über alle 71 Zeilen träfe drei Klassen,
   die den alten Namen als **Gegenstand** tragen. (a) Der **Dateiname** des Verdikts (`2026-11-26-welle-98-architect-verdikt.md`, 13 Nennungen) bleibt, er ist
   kein Name der Welle, sondern der eines Reports, den dieser Slice nicht umbenennt. (b) **Aussagen über den Umzug** selbst — dieser Plan (§1, §3), Frage 2 des
   Verdikts und der Welle, die Vorbedingungs-Sätze in den Closure-Notizen — bleiben, weil sie den alten Namen nennen, um ihn zu benennen. (c) **Suchmuster**
   in Zähl-Kommandos behalten das Muster, das sie zählen. Alles andere (Link-Ziele, Link-Texte, `**Welle:**`-Felder, Roadmap-Zeiger, der Anker an der Ventil-Zeile der
   `.d-check.yml`, die Vergangenheits- und Zustandssätze, die die Welle als Bezug nennen) wird auf `wave-98` gezogen; das Ergebnis ist eine **Liste je Datei** im
   Hand-back, nicht eine Summe.
4. **Commit-Zerlegung, wie sie hier gefahren wird:** (1) das Werkzeug samt Tests; (2a) der **reine** `git mv` der Welle (`R100`); (2b) der Nachzug der Verweise **und**
   die Umstellung der Konfiguration in **einem** Commit, weil Roadmap-Zeiger, Datei und `planning.waves` im selben Arbeitsbaum stehen müssen (Verdikt Frage 2); (3) Commands,
   Briefing, Sensor-Dateien, der Koexistenz-Eintrag samt Index. **Zwischenstand, der nicht die Spitze eines Push werden darf:** (2a) allein — die Roadmap zeigt dort noch auf die
   alte Datei. Der Übergang von [`ADR-0117`](../../adr/0117-wellen-namen-ueber-die-konfiguration.md) nach `Accepted` ist **nicht** Teil dieses Laufs: er gehört der Closure (§2, dritter Punkt), und nach `Accepted` ist die ADR immutabel — ein
   Review-Befund an ihr bliebe sonst ohne den Weg über die Änderung.

**Plan-Änderung (Implementer, vor dem Konfigurations-Commit; derselbe Stand, nach dem Werkzeug-Commit):**

1. **Die Zeile zu [`.d-check.closure.yml`](../../../../.d-check.closure.yml) in der Tabelle** („eine `structure`-Regel `files: …/wave-*-results.md` **neben** der Regel für
   `welle-*-results.md`“) **ist nicht ausführbar.** Gemessen am Arbeitsbaum mit der zweiten Regel: `make verify-closure-notes` endet mit Exit 2 und einem Befund,
   `docs/plan/planning/done/wave-*-results.md:1 … section-missing … Regel trifft keine Datei (auch nach Abzug von exempt-paths) — das Gate liefe leer`
   (1182 Dateien geprüft): Das Modul `structure` meldet eine Regel ohne Kandidaten fail-closed
   ([`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in), Spezifikation [`SPEC-049`](../../../../spec/specification.md#4-grund--und-fehler-codes), „Null Kandidaten“),
   und eine `wave-<NN>-results.md` gibt es erst bei der Closure der Welle. Die Regel erst dort anzulegen ließe das Profil bis dahin ohne die neue Form (DoD, Punkt *Sensoren*).
   **Entscheidung (Konfiguration, kein Produkt-Eingriff):** die **eine** bestehende Regel liest `docs/plan/planning/done/*-results.md` und damit beide Schreibweisen; die
   Menge ist heute dieselbe wie vorher (`git ls-files | grep -cE 'done/[^/]*-results.md$'` = 38, davon ohne `welle-`: `… | grep -vc 'done/welle-'` = 0). Die **Grenze** steht im Kommentar der
   Konfiguration: das Muster liest jede `*-results.md` unter `done/`. Eine Senkung ist es nicht — die Prüfung gilt weiter für alle 38 Bestandsnotizen und kommt für die neue
   Schreibweise hinzu; die Probe (g) fährt beide Seiten (eine `welle-97-results.md` und eine `wave-99-results.md` je **ohne** den Abschnitt `## Was wurde geliefert`).
2. **Commit-Zerlegung, Feinschnitt zu Punkt 4 der Plan-Änderung oben:** (2a) der reine `git mv`; (2b) **diese Plan-Datei** — die Notiz und ihr eigener Verweis-Nachzug auf `wave-98` —
   als eigener Commit **vor** dem Code der Konfiguration; (2c) der Verweis-Nachzug in allen übrigen Dateien **und** die Umstellung der Konfiguration.

**Reihenfolge und Commit-Zerlegung** ([`AGENTS.md`](../../../../AGENTS.md) §3.3):

1. **Vorbedingung, kein Commit dieses Slice:** das Architect-Verdikt zu den Fragen der Welle liegt vor,
   [`slice-295`](../done/slice-295-baseline-v700-bump.md) (die Vorlagen heißen jetzt `wave.template.md`, …) und
   [`slice-297`](../done/slice-297-spezifikation-heisst-specification.md) sind `done/`.
2. **Commit 1: das Werkzeug** samt Tests (`make archive-wave-test`), unabhängig von der Konfiguration — es liest beide Formen,
   solange nichts umgestellt ist.
3. **Commit 2: der Name der Welle, die Konfiguration und die Sensoren** — **erst jetzt**, weil die Umstellung von `planning.waves` die
   `welle-`-Wächter des Bestands abschaltet (Probelauf Punkt 5 in [wave-98 §1](../wave-98-baseline-v700.md)). Verdikt Frage 2: `welle-98`
   heißt ab hier `wave-98`; **keine Lücke**, weil Move, Roadmap-Zeiger und Konfiguration im selben Arbeitsbaum stehen und der Hook ihn prüft.
   Zuerst der **reine** `git mv` der Welle-Datei (R100, [`AGENTS.md`](../../../../AGENTS.md) §3.3), dann der Nachzug **aller** Verweise
   auf den alten Namen (`git grep -c 'welle-98' -- . ':!.harness/baseline'` neu zählen: Roadmap-Zeiger, `**Welle:**`-Felder, Closure-Notizen und
   Reports von [`slice-295`](../done/slice-295-baseline-v700-bump.md) und [`slice-297`](../done/slice-297-spezifikation-heisst-specification.md), der Plan selbst; darin
   `welle-98-results.md` → `wave-98-results.md`) und die Umstellung der Konfiguration. Pfad-begrenzt committen (`BEO-ALL/path-scoped-commit-carries-staged-rest`).
4. **Commit 3: Commands, Briefing, Sensor-Dateien und der Koexistenz-Eintrag.**
5. Review (höchstens zwei Runden), Verifikation, Closure. `git commit -F <datei> -- <Pfade>`, vorher `git diff --cached --stat`.

**Gates:** `make archive-wave-test`, `make doc-check`, `make planning-check`, `make verify-closure-notes` (das Closure-Profil),
`make gates`; die **Probe** aus [wave-98 §3 (g)](../wave-98-baseline-v700.md) in einer Wegwerf-Kopie (nicht im Repo).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`slice-295`](../done/slice-295-baseline-v700-bump.md) und
[`slice-297`](../done/slice-297-spezifikation-heisst-specification.md) liegen in `done/`; `in-progress/` trägt keinen Slice
(nur die Roadmap); das **Architect-Verdikt** zu den zwölf Fragen in [wave-98 §5](../wave-98-baseline-v700.md) (für diesen
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
der unveränderte Bestand — steht in [ihrem Closure-Trigger](../wave-98-baseline-v700.md) (g) und nicht in dieser DoD
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
  — **Ausgang:** *entfallen.* §7 nennt je Sensor beide Seiten, und Review und Verifikation haben sie in Wegwerf-Klonen
  gefahren; die einzige einseitige Stelle ist die gewollte Senkung von `planning` (Risiko 4).
- **Ein Test läuft grün, ohne die Regel zu fahren**
  ([`BEO-ALL/testfall-faehrt-seine-regel-nicht`](../observations/BEO-ALL/testfall-faehrt-seine-regel-nicht/observation.md),
  2×). — **Ausgang:** *entfallen*, wenn die Closure-Notiz je neuem Test die Mutation (die Änderung zurücknehmen) und die
  gelesene Meldung nennt und jeder rot war; sonst *eingetreten* (dann 3×).
  — **Ausgang:** *eingetreten.* Nicht jeder neue Test war unter der Rücknahme rot (9 von 14, §7), und
  `TestRun_PlanOhneEigeneSlices_Wave` nannte eine Rot-Bedingung, die er nicht fuhr (Review F-3; Kommentar behoben in
  `256a84b3`). Dazu ein bestehender Test, der mit abgeschalteter Regel grün bleibt (F-9, nicht behoben). Evidenz-Datei
  `evidence/slice-298.md`, damit 3×. Folge-Kennung: die Closure von [`wave-98`](../wave-98-baseline-v700.md), deren
  Lese-Schritt dem Eintrag seinen Ausgang gibt; ein Folge-Slice ist nicht geschnitten.
- **Das Werkzeug ändert den Bestands-Lauf** — `make archive-wave WELLE=welle-97` liefert danach etwas anderes. — **Ausgang:**
  *entfallen*, wenn der Dry-Run vorher und nachher mit `diff` **byte-gleich** ist (stdout, stderr, Exit); sonst *eingetreten*,
  und die Änderung geht zurück.
  — **Ausgang:** *entfallen.* Werkzeug von `db329ff2` gegen den Stand nach der Änderung, derselbe Baum, `welle-97`,
  `welle-96` und `welle-95` im Dry-Run: stdout mit `cmp` gleich, Exit gleich (0), keine Fehlerausgabe des Werkzeugs
  (Review und Verifikation, unabhängig). **Grenze:** stderr als Ganzes ist nicht byte-gleich, es trägt das
  Docker-Build-Log; verglichen ist die Ausgabe des Werkzeugs.
- **Die Umstellung von `planning.waves` schaltet die Wächter des Bestands ab** (38 Register-Zeilen, 38 Ergebnisnotizen; gemessen
  im Probelauf) — das ist die **gewollte** Folge der kleinsten Antwort, keine Panne. — **Ausgang:** *weiter offen* → Register
  (neue Beobachtung mit Beleg `evidence/slice-298.md`: „eine Umstellung des Präfixes schaltet die Wächter der anderen
  Schreibweise ab“); die Senkung selbst und ihre Wiedervorlage „das Modul liest mehrere Präfixe“ trägt
  [`ADR-0117`](../../adr/0117-wellen-namen-ueber-die-konfiguration.md) ([wave-98 §5](../wave-98-baseline-v700.md) Frage 3, größere Antwort).
  — **Ausgang:** *weiter offen.* Register:
  [`BEO-ALL/praefix-umstellung-schaltet-waechter-der-anderen-schreibweise-ab`](../observations/BEO-ALL/praefix-umstellung-schaltet-waechter-der-anderen-schreibweise-ab/observation.md),
  neu angelegt, 1×; die Senkung steht in der angenommenen
  [`ADR-0117`](../../adr/0117-wellen-namen-ueber-die-konfiguration.md).
- **Ein `wave-<NN>`-Token trifft einen falschen Treffer** in Spec, Sicht oder ADR (die Klasse `wave` braucht Ziffern:
  `wave-drift` ist kein Treffer). — **Ausgang:** *entfallen*, wenn `make doc-check` nach dem Anlegen der Klasse 0 Befunde
  meldet — am Stand `590dde3d` trägt kein Dokument außerhalb dieser Pläne ein `wave-<NN>`-Token (`git grep -l -E
  'wave-[0-9]{2,}' -- . ':!.harness/baseline'` nennt nur `planning_test.go`); sonst *eingetreten*.
  — **Ausgang:** *entfallen.* `make doc-check` meldet mit der Klasse 0 Befunde, und
  `git grep -n -E 'wave-[0-9]{2,}' -- spec docs/plan/adr` ist leer. Der falsche Treffer war angelegt (zwei
  Beispiel-Kennungen in der damals vorgeschlagenen ADR) und ist **vor** der Konfiguration beseitigt worden (Plan-Änderung
  und Verdikt-Nachtrag V2, `db329ff2`).

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

- **Was hat funktioniert:** Der Zuschnitt „Konfiguration statt Produkt, Bestand bleibt“ hielt: `internal/`, `spec/` und die
  Defaults sind ohne Diff, die einzige umbenannte Datei ist die Welle (`9d295355`, genau eine `R100`-Zeile). Je Sensor sind
  **beide** Seiten in Wegwerf-Klonen gefahren (Verifikation, Probe (g); Review): `planning` meldet eine flache `wave-99`-Datei
  ohne Zeiger als `wave-drift` (Exit 2, mit Zeiger Exit 0) und schweigt zu einer `welle-77`-Datei ohne Zeiger; `matrix` meldet
  `wave-98` in Spezifikation, Lastenheft, Sicht und ADR und weiter `welle-97` in der Sicht; das Closure-Profil meldet eine
  `wave-99-results.md` und eine `welle-97-results.md` ohne den Abschnitt je mit `section-missing`; die `ids`-Ausnahme trägt
  beide Verzeichnisformen (Kontrolle `done/xwave-99/`: zwei `id-unlinked`); `tools/archive-wave` liest beide und schreibt das
  Archiv nach dem Präfix. Review und Verifikation ohne HIGH und MEDIUM.
- **Die Mutation je neuem Test (Risiko 2), im Wegwerf-Klon, Meldung gelesen (Verifikation):** `wave_names_test.go` führt
  **14** neue Tests. Unter der Gesamt-Rücknahme (`collect.go` auf `db329ff2`, `ArchiveName` immer `archiv.zip`, Suffix
  `-archiv.zip`) sind **9** rot und **5** grün, dazu 2 bestehende rot. Je Einzel-Mutation: `collect.go` zurück („wave-98
  sammelt []“), Vergleich nur über die Ziffern („welle-98 sammelt [slice-701-a.md slice-702-b.md …]“), `ArchiveName` je
  Richtung („ArchiveName(\"wave-98\") = \"archiv.zip\", erwartet \"archive.zip\"“ und umgekehrt), Suffix („erwartet genau
  slice-601-archive.zip, got [slice-601-archiv.zip]“). Die fünf grünen sind drei Regressionswächter an präfix-generischem
  Code (`TestExtractTitle_Wave`, `TestFindWellePlan_Wave`, `TestRun_DryRun_NoWrites_Wave`: rot, wenn `FindWellePlan` oder
  `h1RE` `wave-` ablehnt oder der Dry-Run schreibt) und zwei Gegenproben, die unter der Rücknahme grün bleiben **müssen**
  (`TestFixture_Bestand_SchreibtWeiterArchiv`: rot bei `ArchiveName` immer `archive.zip`;
  `TestRunSlice_AcceptsContextualWaveMention`: rot, wenn die Zugehörigkeit aus dem ganzen Feld gelesen wird). Damit ist
  Review F-7 berichtigt: es nennt 13 neue und 7 neue rote, gemessen sind 14, 9 und 5 (V-5).
- **Was ging anders als geplant:** **(1) Das Closure-Profil bekam keine zweite Regel**, sondern die eine liest jetzt
  `done/*-results.md`: eine Regel ohne Kandidaten meldet `section-missing`, und eine `wave-<NN>-results.md` gibt es erst bei
  der Closure der Welle (Plan-Änderung `21d6f711`). Grenze, am Eintrag kommentiert: das Muster liest jede `*-results.md`
  unter `done/`, nicht nur Wellen-Notizen. **(2) Zahlen des Plans** stimmten nicht und sind in den Plan-Änderungen
  überschrieben; die dort „17 Platzhalter-Form `seit welle-<NN>`“ ist zu eng beschrieben (V-3: von den 17 Treffern ohne
  Ziffern sind 10 `<NN>`, 2 `<Kennung>`, 5 Prosa oder Zählkommando; 36 und 19 reproduzieren sich). **(3) `wave-98` in zwei
  ADR-Beispielen** hätte die neue Klasse getroffen; vor der Konfiguration vom Architect auf den Platzhalter geschrieben
  (`db329ff2`). **(4) Review F-1/F-2** (Kommentar am Closure-Profil: „heute“-Zahl und verworfene Alternative, `062e0c2a`),
  **F-3** (Test-Kommentare, `256a84b3`), **F-5** (der Nachzug machte einen Mess-Satz im Report zu `slice-297`
  selbstwidersprüchlich, `12fbc567`), **F-6 und Verifikation V-1** (Bestandszahlen der ADR und das Anker-Kommando in
  [`MR-082`](../../../../harness/conventions.md#mr-082), das Zitate mitzählte: 19 → 14, `f30cb082` und `208de9c6`). Bei der
  Closure erneut gemessen: 14 Anker und 435 Pfade nach den Kommandos des Eintrags.
- **Vorbehalt zum DoD-Punkt „Review durchgeführt“:** Eine Review-Runde. Die Fixrunde (`f30cb082`, `062e0c2a`, `256a84b3`,
  `12fbc567`) bekam keine zweite; die Verifikation hat sie nachgesehen. `062e0c2a` ändert in der Gate-Konfiguration nur
  einen Kommentar. `208de9c6` (V-1) entstand **nach** der Verifikation und ist von keiner der beiden Rollen gesehen; der
  Planner hat das Diff gelesen und die zwei Zahlen reproduziert.
- **Befunde ohne Umsetzung (INFO):** **F-8** die Wortgrenze der Kennungs-Regex trägt keinen Test (ohne `\b` bleiben alle
  Tests grün; die Regex-Probe verhält sich richtig). **F-9** `TestRunSlice_RejectsWelleSlice` bleibt mit abgeschalteter
  Zugehörigkeits-Prüfung grün: dem Fixture fehlt `docs/reviews`, der Lauf endet an anderer Stelle mit einem Fehler, und der
  Test verlangt nur irgendeinen. **F-10** die Kommentare zum Verzeichnis von `make slice-mv` nennen nur `welle-<NN>`; das
  Target reicht den Wert durch, das Verhalten stimmt. **V-4** zwei Zwischenstände sind einzeln rot, der Plan nannte einen
  (Review F-4): `9d295355` (28 `target-missing`, ein `wave-drift`) und `21d6f711` (18 `target-missing`, ein `wave-drift`).
- **Die Senkung ist benannt und angenommen:** [`ADR-0117`](../../adr/0117-wellen-namen-ueber-die-konfiguration.md)
  (`Accepted` mit dieser Closure). `planning` bewacht die Bestands-Wellen nicht mehr (`wave-drift`,
  `wave-results-missing`: 38 Register-Zeilen, 38 Ergebnisnotizen); `matrix`, Closure-Profil und Werkzeug lesen beide Formen.
  Zwei Angaben der ADR sind enger zu lesen, als sie dastehen, und tragen die Entscheidung nicht: „verlangt Gleichheit“ — die
  Konfigurations-Prüfung verlangt, dass das Präfix von `results-glob` mit dem von `glob` **beginnt**
  ([`DC-FA-PLAN-001.a`](../../../../spec/specification.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning), Schritt W1); „acht
  Stellen von `wavePrefix`“ ist mit `git grep -c wavePrefix -- internal` (5 Zeilen) nicht zu reproduzieren.
- **Steering-Loop-Eintrag:** mit diesem Slice ist keine Regel und kein Sensor verkörpert. **Benannte Spec-Lücke:**
  [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) —
  `planning.waves` kennt **ein** Kennungs-Präfix; ein Repo, in dem zwei Schreibweisen nebeneinander stehen, kann nur eine
  bewachen. Nicht in die Spezifikation geschrieben; Träger sind der Re-Evaluierungs-Trigger der
  [`ADR-0117`](../../adr/0117-wellen-namen-ueber-die-konfiguration.md), [wave-98 §5](../wave-98-baseline-v700.md) Frage 3
  (Entscheidung des Auftraggebers) und der neue Registereintrag. **Lernsignal 1, mit Ursache:** Die DoD verlangte „jeder neue
  Test läuft ohne die Änderung rot“ und war damit für Gegenproben und Regressionswächter unerfüllbar, **weil** die Klausel
  eine Mutation für alle Tests setzt statt je Test die Änderung zu verlangen, die ihn färbt. **Lernsignal 2, mit Ursache:**
  Ein Kommando, das eine Zeichenkette zählt, zählt den Plan mit, der sie zitiert, **weil** der Ausschluss nach
  Verzeichnissen (`done/`) formuliert war und der Plan beim Zählen noch nicht dort lag (V-1).
- **Beobachtungs-Register (`../observations/`):** neu
  [`BEO-ALL/praefix-umstellung-schaltet-waechter-der-anderen-schreibweise-ab`](../observations/BEO-ALL/praefix-umstellung-schaltet-waechter-der-anderen-schreibweise-ab/observation.md)
  (1×, `offen`; Risiko 4). Je eine `evidence/slice-298.md` in
  [`BEO-ALL/testfall-faehrt-seine-regel-nicht`](../observations/BEO-ALL/testfall-faehrt-seine-regel-nicht/observation.md)
  (**3×, `offen`: Schwelle erstmals erreicht**; Risiko 2, F-3 und F-9),
  [`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
  (13×, `verkörpert`; V-1, F-6, V-5) und
  [`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
  (7×, `verkörpert`; F-5).
- **Architektur-Beobachtung:** Eine Kennungs-Form hat in diesem Repo vier Leser, und jeder antwortet anders auf zwei
  Schreibweisen: der reguläre Ausdruck des Werkzeugs liest beide, `matrix` führt zwei Klassen, `structure` ein weiteres
  Glob, und `planning` **muss** wählen, weil es aus dem Glob-Präfix nicht nur die Zugehörigkeit, sondern die **Identität**
  (die Kennung) ableitet. Die Senkung folgt aus dieser Kopplung, nicht aus der Konfiguration.
- **Push-Pflicht:** Alle Commits dieses Slice ab `73f52639` bis zum Move nach `done/` gehören in **einen** Push, zusammen mit
  den ungepushten davor (`git log --oneline origin/main..HEAD`); `9d295355` und `21d6f711` sind einzeln rot und dürfen nicht
  Spitze eines Push werden. Nach dem Push die CI an der **Spitze** lesen.
- **Aufgaben für die Closure von [`wave-98`](../wave-98-baseline-v700.md):** (1) **Lese-Schritt**
  `testfall-faehrt-seine-regel-nicht` (3×): Ausgang zuweisen (`verkörpert`, `geplant` oder `gestrichen`); der Eintrag darf
  die Welle-Closure nicht als `offen` überstehen. (2) Die Messungen ihres Closure-Triggers, darunter (g) am Endstand,
  `make fullbuild` und `make blackbox-probe` — hier **nicht** gefahren. (3) Push und CI. (4) Ihre eigenen Artefakte tragen
  die neuen Namen (`wave-98-results.md`, `done/wave-98/archive.zip`, `seit wave-98`).
- **Release:** keiner; kein Bestandteil des Images ändert sich, keine `CHANGELOG`-Notiz.
- **Folge-Slices:** keine neuen. [`slice-296`](../open/slice-296-baseline-v700-neuerungen-sichtung.md) (Frage B: ob die
  Defaults auf `wave-*` wechseln) besteht schon in `open/`. F-8, F-9 und F-10 sind benannt, nicht geschnitten.
- **Risiken aus §6:** 1 *entfallen*, 2 *eingetreten*, 3 *entfallen*, 4 *weiter offen*, 5 *entfallen* — Begründung je Risiko
  in §6.
- **Die drei Paarungen**, hier vorab geprüft (getragen von der Closure von wave-98): (a) Anker — kein Eintrag mit
  `liegt in`-Feld, keine Paarung auszulösen; (b) Folge-Slice — `slice-296` besteht als Datei in `open/`; (c) Register —
  die vier genannten Verzeichnisse bestehen und tragen nicht leere `evidence/`.

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
