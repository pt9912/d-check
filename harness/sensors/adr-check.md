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

3. **Pack-Dateien ohne gültiges Hash-Suffix oder ohne passenden Index sind
   unsichtbar — der Lauf bricht dann ab.** Das Modul liest die
   Objektdatenbank über go-git (`v5.19.2`), und go-gits `DotGit`-Schicht
   findet einen Pack nur unter git's kanonischem Namen
   `pack-<Hash-des-Packs>.{idx,pack}`. Gemessen an sechs Proben in einem
   isolierten Repo (`.idx`/`.pack`/`.rev` gemeinsam umbenannt, sonst nichts
   verändert): `pack-<eigener Hash>` ⇒ Exit 0; `pack-<gültige Form, fremder
   Hash>`, `pack-zzzzzzzz`, `packXYZ`, `xpack-abc` ⇒ **je Exit 2**. **`git`
   selbst liest weiter** — es enumeriert `*.idx` unabhängig vom Namen.

   **Seit [`slice-226`](../../docs/plan/planning/in-progress/slice-226-vcs-pack-alias-fremdes-praefix.md)
   <!-- d-check:status-provenance --> löst der Adapter einen Pack zusätzlich
   unter seinem kanonischen Namen auf, wenn seine Datei einen gültigen
   SHA1/SHA256-Hash als Namens-Suffix trägt **und** eine passende
   `.idx`-Datei existiert.** `loose-<eigener Hash>` liefert damit nicht mehr
   `Exit 2`, sondern den echten Befund — genau der Fall, den `git maintenance
   run --task=loose-objects` erzeugt (siehe unten). Ein Präfix allein macht
   einen Pack also nicht mehr unsichtbar; unsichtbar bleibt nur ein Pack ohne
   gültiges Hash-Suffix oder ohne passenden Index — die verbleibenden vier
   Proben oben.

   **Auslöser in freier Wildbahn:** `git maintenance run --task=loose-objects`
   schreibt `loose-<Hash>.pack`. **Vier Ausgänge waren gemessen, bevor
   slice-220 die Kandidaten-Bestimmung umbaute — heute brechen alle vier
   einheitlich ab:**

   | Verschluckt | vor slice-220 | **heute** |
   | --- | --- | --- |
   | eine ganze Referenz, ein Commit, das BASE-Blob | Exit 2 mit `… nicht auflösbar` / `… nicht lesbar` | **Exit 2** (unverändert) |
   | der **HEAD**-Tree des geschützten Verzeichnisses | **Exit 1** mit `core-drift-vcs` *„gelöscht oder umbenannt"* — eine Fehldiagnose: laut, aber inhaltlich falsch | **Exit 2** mit `Range-Spitze "…" nicht auflösbar: … nicht lesbarer Unterbaum …` |
   | der **BASE**-Tree, **mit Pendant** auf der Gegenseite | Exit 2 (seit slice-218) | **Exit 2** (unverändert) |
   | der **BASE**-Tree, **ohne Pendant** (Verzeichnis gelöscht/umbenannt) | **Exit 0, 0 Befunde** — die Löschung erreichte die Änderungsliste nie | **Exit 2** mit derselben Meldungsform |

   **Abhilfe, falls Sie trotzdem auf einen unlesbaren Pack treffen** — seit
   slice-226 nur noch für einen Pack ohne gültiges Hash-Suffix oder ohne
   passenden Index nötig, nicht mehr für einen bloß umbenannten:
   `git repack -A -d` (die `-A`-Form, damit unerreichbare Objekte lose werden
   statt verworfen).

   **Warum die letzten beiden Zeilen bis `v0.76.0` offen blieben.** Der
   Adapter bestimmte seine Kandidaten aus einem **Tree-Diff** (`BASE..HEAD`);
   dessen Walker (aus derselben Bibliothek, die auch `tree.Files()` trägt)
   macht aus einem nicht ladbaren Unterbaum ein stilles `io.EOF` statt eines
   Fehlers — eine Löschung ohne Pendant erreichte die Änderungsliste also
   **nie**, und ein unlesbarer HEAD-Tree ließ den Diff eine Löschung
   **erfinden**, die dann als `core-drift-vcs` gemeldet wurde. **slice-220
   hat die Kandidaten-Bestimmung umgebaut, nicht nur gepatcht:** Die
   geschützte Klasse (`vcs.paths`) wird jetzt über einen eigenen,
   fail-closed Tree-Walker (`repo.TreeObject`, statt des geteilten
   `Files()`-Walkers) **vollständig** gegen BASE **und** HEAD aufgelöst,
   unabhängig von jedem Diff — ein nicht ladbarer Unterbaum bricht diese
   Auflösung selbst ab, **bevor** irgendein Pfad als Added/Deleted/Modified
   eingeordnet wird. Damit fallen alle vier Zeilen auf denselben Codepfad;
   `CO-001` (jetzt [aufgelöst](../../docs/plan/carveouts/done/CO-001-vcs-range-stiller-skip.md))
   trug die dritte und vierte, bis der Fix im Release `v0.76.1` das
   publizierte Image erreichte.

   **Zwei frühere Ausprägungen — Blob und BASE-Tree mit Pendant — waren
   bereits mit slice-218 geschlossen**, über eine andere Reparatur:
   go-git meldet ein *unlesbares* Objekt mit demselben Fehler wie eine *im
   Tree fehlende* Datei (`object.ErrFileNotFound`); der zweite Fall ist
   legitim (Datei später angelegt) und muss befundfrei bleiben, der erste ist
   ein Umgebungsfehler. `FileAt` trennt beide seither über den
   **Tree-Eintrag** (Name und Modus, ohne den Datei-Inhalt) — dieser
   Mechanismus ist von slice-220 **unberührt** und trägt weiterhin. **Die
   Gegenrichtung gehört zur Zusage:** ein Eintrag **ohne** Datei-Inhalt —
   Verzeichnis oder Gitlink — ist kein unlesbares Objekt und bleibt
   befundfrei; Regressionstests halten beide Richtungen
   (`TestFileAtUnlesbaresObjekt`, `TestFileAtEintragOhneBlob` in
   `internal/adapter/driven/git/`) und fallen ohne ihren Fix — verifiziert
   durch Zurücknehmen.

   **Und die Gegenrichtung hat einen Preis, der hierher gehört:** Weil ein
   Eintrag ohne Datei-Inhalt befundfrei bleibt, geht auch der Fall durch, in
   dem eine `Accepted`-ADR am geschützten Pfad durch einen **Gitlink** ersetzt
   wird — sie ist dann verschwunden, und niemand meldet es. Das ist kein
   Neuzugang, sondern eine benannte, unverändert bestehende Grenze.

   **Für gepinnte Konsumenten galt das länger:** Wer ein veröffentlichtes
   Image `v0.76.0` oder älter fährt, hat die dritte/vierte Ausprägung
   weiterhin — Konsumenten lösen das durch einen Pin-Wechsel auf `v0.76.1`
   oder neuer. **Dieses Repo ist nicht betroffen**, weil `make adr-check` die
   Prerequisite `build` trägt und aus dem Quellstand baut, der den Fix seit
   slice-220 trägt.

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
