# Review slice-295 — Baseline-Pin-Hebung auf `v7.0.0`

- **Review-Art:** Code/Doku (gegen Plan, Architect-Verdikt, `MR-081`, Hard Rules)
- **Gegenstand:** `slice-295`, Diff `738d2bca..HEAD` (Plan-Änderung `214a980f`, neuer Baum `79a14c01`, Hauptcommit `04dca313`, Move `d67b6221`, Link-Tiefen `3d8f3555`, Treffermenge `50913578`)
- **Skill:** `reviewer.md` @ 1.20.0
- **Modell-ID:** claude-sonnet-5-5
- **Datum:** 2026-11-26
- **Eingangs-Kontext:** `slice-295` (§1 Abgrenzung, §3 Plan samt Plan-Änderung, §6 Risiken), `welle-98` (§3, §5), das Architect-Verdikt zu `welle-98` (Fragen 5, 9, 12), `MR-081`, Vorgänger `MR-080` und `MR-074` als Form-Vergleich, `MR-069`, `AGENTS.md` §3.1 / §3.3 / §3.5 / §3.6 / §3.7 / §3.8 / §5 / §6.
- **Zusätzlich gefahren (Wegwerf im Scratchpad, nichts im Repo):** alter Baum per `git archive 738d2bca` entpackt und Datei für Datei gegen den vendorten Baum verglichen (ohne `<!-- Quelle:`-Zeilen, Versionsstempel normalisiert); das Release-Bundle von `v7.0.0` per `diff -r` gegen den vendorten Baum; jede geänderte lebende Datei gegen ihren Vorzustand mit zurückgetauschtem Tag verglichen (nur was *nicht* Tag-Tausch ist, bleibt übrig).

## Findings

Keine HIGH. Ein MEDIUM, zwei LOW, drei INFO.

### F-1 — MEDIUM — Begründung des Ventils trägt für keine der sieben Dateien; der Vorgänger behandelte `slice-274` anders

- **quelle:** Prüffrage 20 (Begründung neben richtiger Entscheidung), `MR-069`, `MR-080` (Ausnahme von der Frozen-Regel)
- **pfad:** `harness/conventions/MR-081-baseline-v700.md` · „und ein Anker kann im neuen Baum fehlen (die Überschrift in `modul-14`)“
- **befund:** Der Grund für die sieben `ignore-refs`-Einträge nennt neben „Nachziehen verfälschte den Stand ihrer Zeit“ einen fehlenden Anker (`modul-14`); keine der sieben Dateien verlinkt `modul-14` — sie zeigen auf `modul-04` (fünf ADRs), `modul-11` §Bewusstes Brechen (`slice-274`) und `modul-05` (Verdikt zu `welle-93`), und alle drei Ziele samt Anker lösen im Baum `v7.0.0` auf (`modul-11` ist byte-gleich bis auf den Stempel). Dieselbe Datei `slice-274` hat `MR-080` noch retargetet, weil sonst `ignore-refs` wüchse; `MR-081` lässt sie diesmal im Ventil („diesmal bleibt die Datei unberührt“) ohne zu sagen, warum der damalige Grund nicht mehr gilt, und der Plan wendet für den Review-Link („Link, der weiter auflösen kann“) den gegenteiligen Maßstab an.
- **Failure-Szenario:** Das Ventil wächst um zwei Einträge (`slice-274`, Verdikt `welle-93`), die ein Retarget wie bei `MR-080` vermeiden könnte; jeder künftige Link auf `v6.18.1` in diesen zwei Dateien wird still, und der Zähler des Auflösungs-Triggers von `MR-069` („kein neuer Eintrag über drei Hebungen“) startet wegen vermeidbaren Wachstums neu. Bei den fünf `Accepted`-ADRs bliebe die Entscheidung dieselbe (Stand ihrer Zeit, `AGENTS.md` §3.5) — das ist der Test von Frage 20: die Entscheidung stimmt, die zweite Begründung nicht.
- **Hinweis zur Stufe:** Die Ventil-Form für alle sieben stammt aus dem Architect-Verdikt (Frage 5); der Implementer hat sie wie verlangt umgesetzt. Die Prämisse „Anker fehlt“ steht dort ebenfalls — das Finding geht deshalb an Architect/Planner (Konflikt-Pfad, Verdikt 1 oder 2 der Rollen-Sequenz), nicht als Mangel der Umsetzung.
- **verifizierbar:** ja — `git grep -n 'baseline/v6\.18\.1' -- docs/plan/adr docs/plan/planning/done/wellenlos/slice-274-trace-requirements-als-liste-von-quellen.md docs/reviews/2026-10-30-welle-93-architect-verdikt.md` nennt je ein Ziel (`modul-04`, `modul-11`, `modul-05`); `grep -c '^### Bewusstes Brechen für DoD-Testbehauptungen'` im vendorten `modul-11` ergibt 1
- **klasse:** begruendung-trifft-gegenstand-nicht

### F-2 — LOW — Zitat-Delta „51 Fundstellen, 50 gleich“ ohne Kommando, am Vorzustand nicht nachzählbar

- **quelle:** `AGENTS.md` §5 Regel 14 (Form vor Zahl), `slice-295` §6 (Risiko „Eine Zahl im Bump-Eintrag stimmt nicht“: nachzählen mit dem im Eintrag genannten Kommando), Prüffrage 17
- **pfad:** `harness/conventions/MR-081-baseline-v700.md` · „51 Fundstellen, **50** im Wortlaut beider“
- **befund:** Die Form ist genannt („…“, ab 25 Zeichen, lebende Markdown/YAML ohne `done/`, Reviews, CHANGELOG, `conventions/done/`, ADRs, CRs), ein Kommando fehlt. Nachgezählt mit dieser Form (Zeilenumbrüche normalisiert, Fundstelle je Datei und Vorkommen, beide Bäume per Teilstring-Suche) komme ich auf 31 Fundstellen im Baum (30 in beiden, 1 nur im alten), nicht auf 51 (50 / 1); die eine Abweichung, die Überschrift *Zwei Formen des Reproduzierbarkeits-Ankers* im Plan dieses Slice, deckt sich. Die Schlussfolgerung (kein lebendes Zitat verliert seinen Wortlaut) trägt bei beiden Zählungen; die Zahl 51 nicht.
- **verifizierbar:** nein — ohne das Kommando des Implementers ist die Differenz nicht erklärbar; ein genanntes Kommando würde es klären
- **klasse:** zaehlung-ohne-kommando

### F-3 — LOW — Beispiel im Sensor-Text nennt die Dateizahl des alten Baums

- **quelle:** Reviewer-Skill LOW „veraltete Beispiele“ (Doku-Drift), `AGENTS.md` §5 Regel 13
- **pfad:** `harness/sensors/baseline-verify.md` · „`verify ok (54 Dateien, vollständig)` (Exit 0)“
- **befund:** Der Absatz „Gemessen: …“ nennt die Ausgabe von `make baseline-verify` mit 54 Dateien; der Lauf auf dem neuen Baum meldet `verify ok (55 Dateien, vollständig)` (das Manifest führt 55 statt 54 Zeilen). Der Text ist als einmalige Messung formuliert, liest sich aber als aktuelle Ausgabe; die Datei liegt außerhalb der Pin-Spiegel-Klassen, die der Bump zählt, und wurde nicht angefasst.
- **verifizierbar:** ja — `make baseline-verify` (letzte Zeile) gegen die Zeile in der Sensor-Datei
- **klasse:** beispiel-zahl-stammt-vom-vorgaenger-baum

### F-4 — INFO — Bewegungen 5 und 6 tragen kein eigenes „nicht übernommen“

- **quelle:** `slice-295` §1 („markiert jede Bewegung als *nicht übernommen*“)
- **pfad:** `harness/conventions/MR-081-baseline-v700.md` · „**Welle-Vorlage: Slice-Tabelle mit Links**“
- **befund:** Die Bewegungen 2, 3 und 4 schließen mit „Nicht übernommen“ und Verweis auf `slice-296`; 1 nennt die Zuständigkeit von `slice-297`/`slice-298`, 5 (Welle-Vorlage) und 6 (Bundle-Stand) tragen nur die Sammelüberschrift „keine davon übernommen“. Keine der sechs ist als betroffen klassifiziert — die Zeile „Dieses Repo ist betroffen und hat entschieden“ gilt allein der vom Auftraggeber entschiedenen Namensregel.
- **verifizierbar:** ja — Lesen des Eintrags
- **klasse:** marker-je-bewegung-nur-sammelweise

### F-5 — INFO — Der Baum-Zuwachs steht in einem eigenen Commit, der Plan nennt „ein Hauptcommit“

- **quelle:** `slice-295` §3 Reihenfolge, Schritt 3
- **pfad:** `docs/plan/planning/in-progress/slice-295-baseline-v700-bump.md` · „**Ein Hauptcommit:** Baum-Tausch, Pin, Tag-Tausch“
- **befund:** `79a14c01` legt `v7.0.0` allein an (56 Dateien, nichts außerhalb des neuen Baums), `04dca313` tauscht Pin, Verweise und entfernt den alten Baum. Die Begründung des Plans (Zwischenstand mit totem Link passiert den Hook nicht) trifft den Zuwachs-Commit nicht, da er keinen Verweis ändert; ob der Zwischenstand `79a14c01` für sich die Gates passiert, habe ich nicht gefahren. Die Teilung hält den 9613-Zeilen-Diff des Baums vom Retarget-Diff getrennt und ist in der Botschaft benannt.
- **verifizierbar:** nein — nur durch einen Gate-Lauf auf `79a14c01`
- **klasse:** commit-zerlegung-weicht-vom-plan-ab-ohne-schaden

### F-6 — INFO — Inline-Code-Nennungen des alten MR-080-Pfads bleiben

- **quelle:** `slice-295` §3 (Plan-Änderung), `MR-070`
- **pfad:** `docs/reviews/2026-10-20-slice-277-baseline-v6181.md` · „`harness/conventions/MR-080-baseline-v6181.md` · „38 aktive `MR-*`““
- **befund:** Außer dem einen nachgezogenen Link nennen drei eingefrorene Reports den Pfad `harness/conventions/MR-080-baseline-v6181.md` als `pfad`-Feld in Inline-Code; er löst nach dem Move nicht mehr auf. `codepaths` nimmt die Reports aus (`make doc-check` meldet 0), und das `pfad`-Feld zitiert den Stand des Laufs. Bewusst, im Plan benannt.
- **verifizierbar:** ja — `git grep -n 'conventions/MR-080-baseline-v6181' -- docs/reviews`
- **klasse:** historischer-pfad-in-eingefrorenem-report

## Negativbefunde

Gemessene Werte stehen bei „Ergebnis“; **jede Zahl des Eintrags `MR-081` außer dem Zitat-Delta (F-2) stimmt mit dem im Eintrag genannten Kommando überein.**

| Bereich | Ergebnis |
|---|---|
| Vendoring-Integrität | geprüft, ohne Befund — `make baseline-verify` Exit 0, letzte Zeile `fetch-baseline-cache: verify ok (55 Dateien, vollständig)`; `diff -rq` Bundle `v7.0.0` gegen `.harness/baseline/v7.0.0`: einziger Unterschied „Nur in .harness/baseline/v7.0.0: SHA256SUMS“ (das Manifest, erwartet); 56 Dateien (`find … -type f`); `ls .harness/baseline` nennt nur `v7.0.0`; Commit `79a14c01` berührt nichts außerhalb des neuen Baums |
| Delta nach Datei (`MR-081`, Zahlen) | geprüft, ohne Befund — nachgezählt am Archiv `738d2bca` und am vendorten Baum: 54 / 55 Bundle-Dateien, 47 gleichnamig, 7 umbenannt, 1 neu, 38 byte-verschieden, 31 inhaltlich verschieden (19 Regelwerk, 12 Vorlagen), 7 nur Stempel (die im Eintrag genannten sieben), 862 Diff-Zeilen roh, **804** normalisiert; alle Datei-Zahlen der Spitzen und der Aufzählung stimmen (258, 74, 57, 49, 35, 28, 28, 24, 24, 23, 21, 21, 17, 15, 14, 14, 14, 14, 13, 12, 10, 6, 6, 4, 4, 4, 4, 4, 4, 2, 1); die 31 Dateien sind in den Bewegungen 1–6 vollständig und ohne Doppelnennung außer den benannten Anteilen verteilt; die sieben Umbenennungen tragen 28 / 16 / 10 / 10 / 7 / 2 / 2 Diff-Zeilen, die neue Vorlage 74 Zeilen |
| Spiegel-Klassen (`MR-081`) | geprüft, ohne Befund — `git grep -l` / `-o` am Vorzustand `738d2bca`: 77 Dateien / 190 Vorkommen / 175 Zeilen; Frozen 17 / 51 (8 / 29 `done/`, 8 / 21 Reviews, 1 / 1 CHANGELOG); `conventions/done/` und `docs/plan/cr` nannten den Tag nicht (Treffer leer); ADRs 5 / 5; Pläne dieser Welle 5 / 25; `MR-080` 1 / 9; lebend 49 / 100 |
| Treffermenge nach dem Bump (`MR-081`, Commit `50913578`) | geprüft, ohne Befund — das im Eintrag genannte Kommando liefert 10 Dateien (`.d-check.yml`, `MR-081`, ADR 0112 bis 0116, `slice-295`, `slice-296`, `welle-98`); die Summe 17 + 5 + 5 + 1 + 49 = 77 geht auf; die 7 Symlinks nennen den alten Pin nicht (`readlink`-Schleife: 0), alle 9 Einträge unter `.claude/rules/` lösen auf (`[ -e ]`), 9 Einträge im Index haben Modus 120000 |
| Keine übrigen lebenden `v6.18.1`-Treffer | geprüft, ohne Befund — außerhalb Frozen-Klassen, Baum und der zehn benannten Dateien: 0 Treffer; auch nicht in der Form `6.18.1` ohne `v` oder `v6.18` |
| Retarget ist nur Tag-Tausch | geprüft, ohne Befund — die 58 geänderten Dateien außerhalb des Baums und der Pläne (51 Textdateien mit zurückgetauschtem Tag gegen `738d2bca` verglichen, 7 Symlinks per `readlink`): Abweichungen nur an den neun Vorlagen-Namen-Zeilen, den 7 + 1 Zitat-Spannen, der Ventil-Konfiguration, der Index-Zeile von `MR-081` und dem Pfad-Nachzug; die Mischzeilen (`MR-021` Z. 20, `MR-049` Z. 14, `MR-053` Z. 9, `roadmap.md` Z. 13, `observations/README.md`, `AGENTS.md` „Baseline-Form `v7.0.0` §1“) sind Pin-Aussagen über die Form und stehen gegen den neuen Baum (`slice.template.md` §1 bis §8 unverändert benannt; `modul-06` §Roadmap-Struktur trägt die Aussage zu *Offene Wellen* im Wortlaut) |
| `d-check:cite`-Spannen | geprüft, ohne Befund — 15 lebende Direktiven (`.harness/skills`, `MR-005/031/035/039/043/049`); die 7 verschobenen und die 8 unverschobenen Bereiche am Archiv und am neuen Baum Zeile für Zeile verglichen: alle 7 verschobenen wortgleich mit der genannten alten Zeile, alle 8 übrigen unverändert auf denselben Zeilen; die zwei Bereiche der vier Pläne (`373–374` → `416–417`, `379` → `422`) wortgleich, die 8 Direktiven zeigen alle auf `v7.0.0` |
| Vorlagen-Namen (Antwort Frage 12 des Verdikts, Bewegung 1) | geprüft, ohne Befund — 9 Zeilen in 4 Dateien (`close-welle.md` 4, `plan-welle.md` 2, `conventions.md` 1, `stub.go` 2); nur Dateinamen getauscht, `<welle-id>`, `archiv.zip`, `seit welle-`, `WELLE=` unberührt; das Kommando des Eintrags findet außerhalb der Pläne, Reviews, `conventions/done/` und `MR-081` selbst keine alten Namen mehr; `stub.go` nur Kommentarzeilen, keine Kommentar-Klasse verletzt (Kopplung an die Vorlagen-Namen) |
| `.d-check.yml` (sieben Einträge) | geprüft, mit F-1 — Form der von `MR-069` (datei-skopiert, je ein Eintrag, Muster `v6.18.1/**`, die `in:`-Dateien existieren, jede trägt genau den einen Link); Zählung 37 statt 30 Einträge über 13 statt 12 Tags stimmt (`grep -c` am Vorzustand 30 / 12); Kommentar trägt Zusage, Grenze und Rang-Zeiger (`MR-081`, `MR-069`, `MR-052`), keine Review-Historie; `make doc-check` Exit 0, 0 Befunde |
| Abgrenzung (nichts aus `slice-296/297/298`) | geprüft, ohne Befund — `spec/spezifikation.md` weder umbenannt noch inhaltlich berührt (nur die Zeile **Rolle** mit dem Pin, gleiches in `spec/architecture.md`); keine `wave-`-Konfiguration; `internal/`, `tools/` außer `stub.go`-Kommentar, `Makefile`, `Dockerfile`, `docs/user`, `README`, `CHANGELOG`, `spec/lastenheft.md` unberührt; `MR-081` nimmt keine Neuerung an (keine Zeile „betroffen/anwendbar“ außer der entschiedenen Namensregel) |
| Frozen-Stabilität | geprüft, ohne Befund — `git diff --stat 738d2bca..HEAD -- docs/plan/planning/done docs/reviews CHANGELOG.md harness/conventions/done docs/plan/adr` nennt genau eine Review-Datei (eine Zeile) und die zugezogene `MR-080`; `make adr-check RANGE=738d2bca..HEAD` Exit 0, „0 Befund(e)“; ADRs 0112 bis 0116 byte-stabil |
| Pfad-Nachzug im eingefrorenen Review (Frage 3 des Auftrags) | geprüft, ohne Befund — vertretbar: Plan-Änderung steht **vor** dem Code (`214a980f` zeitlich vor `79a14c01`), 1 Datei / 1 Zeile, reiner Link-Ziel-Wechsel auf eine Datei, die lebt; der Plan nennt die Alternative (achter Ventil-Eintrag, Gate-Senkung nach `MR-069`) und die Folge für Risiko 2; das Kommando des Plans (Treffer: Index-Zeile, Inline-Code in drei Reviews, dieser eine Link) stimmt; `MR-081` führt den Nachzug als Ausnahme, die Messung `git diff --stat` ist wie angekündigt nicht leer |
| §3.3-Zerlegung | geprüft, ohne Befund — `d67b6221` ist `R100` (`harness/conventions/{ => done}/MR-080-baseline-v6181.md`, 0 Zeilen); der Inhalt (6 Link-Tiefen in `MR-080`, 1 Link im Review) steht in `3d8f3555`; keine Datei wandert und ändert sich im selben Commit |
| Kommentare (§3.7) und Botschaften (Regel 15) | geprüft, ohne Befund — Kommentar in `.d-check.yml` und `stub.go` ohne Review-Historie oder Deliberation, Herkunft als ein Feld; die Botschaften nennen Maße, die sich nachmessen lassen (49 Dateien / 100 Vorkommen, 7 von 15 Spannen, 9 Zeilen in 4 Dateien, 55 Dateien des alten Baums) und grenzen ab, was nicht Teil des Commits ist; jede Botschaft trägt eine Kennung (`make trace-check RANGE=738d2bca..HEAD` Exit 0, „0 Befund(e)“) |
| `harness/conventions.md` | geprüft, ohne Befund — Pin-Zeiger auf `MR-081`, Index-Zeile in *Aktive*, `MR-080` mit Nachfolger in *Aufgelöste*; Form des Eintrags (Datum, Geltungsbereich, Ersetzt-Baseline-Regel `—`, Adaption, Begründung, Auflösungs-Trigger, Löst auf, Ausgelöst durch) wie `MR-074`/`MR-080`; `Löst auf` zeigt auf die Index-Zeile wie die Vorlage es verlangt; Release-Liste im Eintrag stimmt (`v6.18.1`, `v7.0.0` ohne Zwischen-Tag, gelesen mit `gh release list` am Kurs-Repo) |
| Gates, die ich gefahren habe | geprüft, ohne Befund — `make doc-check`, `make planning-check`, `make gate-consistency`, `make adr-check RANGE=738d2bca..HEAD`, `make trace-check RANGE=738d2bca..HEAD`, `make baseline-verify`: je Exit 0, 0 Befunde. Nicht gefahren (Sache der Verifikation): `make gates`, `make archive-wave-test` |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 (F-1) |
| LOW | 2 (F-2, F-3) |
| INFO | 3 (F-4, F-5, F-6) |

## Verdikt

Der Diff trägt den Plan: mechanischer Bump, Delta und Spiegel-Klassen stimmen mit den genannten Kommandos, keine Neuerung übernommen, Frozen-Bestand stabil bis auf die benannte Ausnahme. **F-1 blockiert typischerweise**: die Entscheidung (Ventil für die fünf `Accepted`-ADRs) ist richtig, die zweite Begründung trifft keine der sieben Dateien, und `slice-274` ist gegenüber `MR-080` ohne Grund anders behandelt — das ist im Konflikt-Pfad zu klären (Architect-Verdikt prüfen: entweder die Begründung im Eintrag auf „Stand ihrer Zeit“ kürzen und die Abweichung zu `MR-080` benennen, oder `slice-274` und das Verdikt zu `welle-93` wie damals retargeten). Ein Finding wird nicht herabgestuft, weil das Verdikt die Prämisse trägt. F-2 und F-3 sind LOW und kein Merge-Hindernis.
