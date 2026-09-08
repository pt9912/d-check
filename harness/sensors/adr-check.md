# `make adr-check` — hält, dass eine `Accepted`-ADR nicht inhaltlich überschrieben wird

## Vertrag

Via Modul `vcs` (Image, dogfood). Eine `Accepted`-ADR unter
`docs/plan/adr/NNNN-*.md` wird nicht inhaltlich geändert; verglichen wird
`core(BASE)` gegen `core(HEAD)` über die Commit-Range, mit reiner-Go-git im
read-only `.git`.

**Erlaubt bleiben zwei Dinge:** `## Geschichte`-Anhänge und der
`**Status:**`-Übergang — das Status-Feld ist ein Zustandsfeld wie jedes andere
und ausdrücklich **nicht** Teil des Kern-Vergleichs
([`AGENTS.md`](../../AGENTS.md) §3.5, §3.7).

Eine gelöschte oder umbenannte `Accepted`-ADR ist ein **FAIL**.

## Grenze — was das Grün nicht abdeckt

1. **Zwei Modi, ein Bindepunkt-Paar** — `STAGED=1` im `pre-commit`-Hook,
   `RANGE=` in der PR-/Push-CI. Der lokale Hook ist **opt-in pro Klon**
   (`make hooks`); `--no-verify` umgeht ihn, nicht die CI.

2. **`## Geschichte` ist vom Vergleich ausgenommen — bis zum Dateiende.**
   Das Profil führt `exclude-sections: [Geschichte]`, und die Ausnahme reicht
   von der Überschrift bis zum Schluss der Datei. **Gemessen: 79 von 84 ADRs
   führen `## Geschichte` als *letzte* Sektion** — bei ihnen liegt also der
   gesamte Rest hinter dem Wächter. Wer eine Kern-Aussage **dorthin** schreibt,
   passiert. Das ist der Preis dafür, dass Anhänge an die Geschichte erlaubt
   bleiben (`AGENTS.md` §3.5), und es ist die größte Grenze dieser Datei.
   Permanent — die Unterscheidung *Anhang gegen verkleidete Kern-Änderung* ist
   ein Urteil.

3. **Pack-Dateien mit unkanonischem Namen sind unsichtbar — der Lauf bricht
   dann ab.** Das Modul liest die Objektdatenbank über go-git (`v5.19.2`), und
   go-git findet einen Pack nur unter git's kanonischem Namen
   `pack-<Hash-des-Packs>.{idx,pack}`. **Beide Namensteile zählen** — gemessen
   an sechs Proben in einem isolierten Repo (`.idx`/`.pack`/`.rev` gemeinsam
   umbenannt, sonst nichts verändert): `pack-<eigener Hash>` ⇒ Exit 0;
   `loose-<eigener Hash>`, `pack-<gültige Form, fremder Hash>`,
   `pack-zzzzzzzz`, `packXYZ`, `xpack-abc` ⇒ **je Exit 2**. **`git` selbst
   liest weiter** — es enumeriert `*.idx` unabhängig vom Namen.

   **Auslöser in freier Wildbahn:** `git maintenance run --task=loose-objects`
   schreibt `loose-<Hash>.pack`. **Was der Lauf dann tut, hängt davon ab, was
   der unsichtbare Pack verschluckt — und es sind drei verschiedene Ausgänge,
   nicht einer:**

   | Verschluckt | Ausgang |
   | --- | --- |
   | eine ganze Referenz, ein Commit, ein Blob (BASE oder HEAD), der BASE-Tree | **Exit 2** mit `… nicht auflösbar` / `… nicht lesbar` |
   | der **HEAD**-Tree des geschützten Verzeichnisses | **Exit 1** mit `core-drift-vcs` *„gelöscht oder umbenannt"* — **eine Fehldiagnose**: laut, aber inhaltlich falsch |
   | der **BASE**-Tree, **ohne Pendant** auf der Gegenseite (Verzeichnis gelöscht/umbenannt) | **Exit 0, 0 Befunde** — die Löschung erreicht die Änderungsliste nie |

   **Die dritte Zeile ist weiterhin offen** und als
   [`CO-001`](../../docs/plan/carveouts/CO-001-vcs-range-stiller-skip.md)
   geführt: Der Tree-Walker der Bibliothek macht aus einem nicht ladbaren
   Unterbaum ein `io.EOF`, die Änderung entsteht also **vor** jeder Stelle, an
   der dieses Modul prüfen könnte. Gemessen ist auch, dass der naheliegende
   Wachposten nicht trägt: `tree.Files()` benutzt denselben Walker und
   schweigt ebenso. **Abhilfe für alle drei:** `git repack -A -d` (die
   `-A`-Form, damit unerreichbare Objekte lose werden statt verworfen).

   **Die dritte Form gab es bis slice-218 nicht — dort war der Lauf still
   grün, in zwei Ausprägungen.** go-git meldet ein *unlesbares* Objekt mit
   demselben Fehler wie eine *im Tree fehlende* Datei
   (`object.ErrFileNotFound`); der zweite Fall ist legitim (Datei später
   angelegt) und muss befundfrei bleiben, der erste ist ein Umgebungsfehler.
   Gemessen an derselben echten `core-drift-vcs`-Verletzung, kanonisch je
   Exit 1:

   | Der unsichtbar benannte Pack verschluckt … | Ankunft | vor slice-218 |
   | --- | --- | --- |
   | das BASE-**Blob** | `M` | `0 Befund(e)`, Exit 0 |
   | das BASE-**Tree** des Verzeichnisses | `A` | `0 Befund(e)`, Exit 0 |

   Die **zweite** Zeile fand erst Review-Runde 2; sie läuft über den
   `A`-Zweig, der `FileAt` gar nicht rief. Seit slice-218 trennt der Adapter
   die beiden Zustände über den **Tree-Eintrag** (Name und Modus, ohne den
   Datei-Inhalt), und der `A`-Zweig fasst den BASE-Stand an, statt ihn
   ungesehen für frei zu erklären. **Die Gegenrichtung gehört zur Zusage:** ein
   Eintrag **ohne** Datei-Inhalt — Verzeichnis oder Gitlink — ist kein
   unlesbares Objekt und bleibt befundfrei; ein wandernder Submodul-Zeiger in
   der Klasse brach in der Zwischenfassung fälschlich ab. Zwei
   Regressionstests halten beide Richtungen und fallen ohne ihren Fix
   (`TestFileAtUnlesbaresObjekt`, `TestFileAtEintragOhneBlob`,
   `TestVCSAddedMeldetUnlesbareBasis`) — verifiziert durch Zurücknehmen.

   **Geprüft sind die drei Diff-Zustände, die das Modul kennt** (`A`, `M`,
   `D`); alle drei fassen den BASE-Stand jetzt über denselben Adapter-Pfad an.
   **Das deckt aber nur, was im Diff ankommt** — die dritte Tabellenzeile oben
   entsteht davor und bleibt offen. Das ist die Menge, über die hier geurteilt
   wird; „jeder denkbare Repo-Zustand" wäre eine andere Aussage.

   **Und die Gegenrichtung hat einen Preis, der hierher gehört:** Weil ein
   Eintrag ohne Datei-Inhalt befundfrei bleibt, geht auch der Fall durch, in
   dem eine `Accepted`-ADR am geschützten Pfad durch einen **Gitlink** ersetzt
   wird — sie ist dann verschwunden, und niemand meldet es. Das ist kein
   Neuzugang dieses Slice (der Vor-Fix-Stand schwieg ebenso), aber die Regel
   sanktioniert es jetzt ausdrücklich, und deshalb steht es hier.

   **Für gepinnte Konsumenten gilt das noch nicht:** Wer ein veröffentlichtes
   Image ab `v0.75.0` abwärts fährt, hat den stillen Pfad weiterhin — geführt
   als [`CO-001`](../../docs/plan/carveouts/CO-001-vcs-range-stiller-skip.md).
   **Dieses Repo ist nicht betroffen**, weil `make adr-check` die Prerequisite
   `build` trägt und aus dem Quellstand baut.

   **Die Reichweite ist gemessen, nicht geschätzt** — und die Module verhalten
   sich **unterschiedlich**: `commits`
   ([`make trace-check`](trace-check.md)) war immer fail-closed (ein unlesbarer
   Commit ⇒ `commit … nicht lesbar`, Exit 2, gemessen). **`tracked` ist gar
   nicht betroffen**, weil es den git-**Index** liest: Es meldete unter beiden
   Pack-Namen denselben `target-untracked`-Befund — eine Positiv-Kontrolle,
   kein grüner Lauf, der auch *„nichts geprüft"* heißen könnte. **Über andere
   Stellen des Produkts sagt das nichts** — drei Module sind geprüft, nicht
   alle. Permanent, solange die Module über go-git lesen; ein Bump der
   Abhängigkeit ist der Anlass, hier nachzumessen — **und zwar mit einem
   *partiellen* unsichtbaren Pack**, denn nur der deckte den stillen Pfad auf
   *(seit slice-218)*.

## Bindung

**nicht** Teil von `gates`/`ci` — Diff-/Commit-Zeit-Bindepunkt.
[ADR-0016](../../docs/plan/adr/0016-adr-immutable-gate.md) ·
[ADR-0024](../../docs/plan/adr/0024-vcs-immutable-gate.md) ·
[ADR-0025](../../docs/plan/adr/0025-codepaths-ignore-refs.md)
