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
   schreibt `loose-<Hash>.pack`. Je nachdem, was der unsichtbare Pack
   verschluckt, meldet der Lauf `staged-Basis "HEAD" nicht auflösbar:
   reference not found`, `HEAD-Tree nicht lesbar: object not found` oder
   `nicht lesbares Objekt zu "<pfad>" in "<ref>": file not found` — **alle mit
   Exit 2**. Wer nur eine der Formen kennt, erkennt die anderen nicht wieder.
   **Abhilfe:** `git repack -A -d` (die `-A`-Form, damit unerreichbare Objekte
   lose werden statt verworfen).

   **Die dritte Form gab es bis slice-218 nicht — dort war der Lauf still
   grün.** go-git meldet ein *unlesbares* Objekt mit demselben Fehler wie eine
   *im Tree fehlende* Datei (`object.ErrFileNotFound`); der zweite Fall ist
   legitim (Datei später angelegt) und muss befundfrei bleiben, der erste ist
   ein Umgebungsfehler. Der Adapter behandelte beide gleich und übersprang die
   Datei — **eine echte `core-drift-vcs`-Verletzung verschwand** (gemessen:
   Exit 1 mit kanonischem Pack-Namen, Exit 0 nach dem Umbenennen derselben
   Datei, während `git diff` die Änderung unverändert zeigte). Seit slice-218
   trennt der Adapter die beiden über den **Tree-Eintrag**, der Name und Hash
   trägt und das Blob nicht braucht. Ein Regressionstest in `make test`
   (`TestFileAtUnlesbaresObjekt`) fällt ohne diese Unterscheidung — verifiziert
   durch Zurücknehmen des Fixes.

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
