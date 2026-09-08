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
   ab, statt still grün zu melden.** Das Modul liest die Objektdatenbank über
   go-git (`v5.19.2`), und go-git findet einen Pack nur unter git's
   kanonischem Namen `pack-<Hash-des-Packs>.{idx,pack}`. **Gemessen, sechs
   Proben in einem isolierten Repo** (alle Objekte im Pack, `.idx`/`.pack`/
   `.rev` gemeinsam umbenannt, sonst nichts verändert): `pack-<eigener Hash>`
   ⇒ Exit 0; `loose-<eigener Hash>`, `pack-<gültige Form, fremder Hash>`,
   `pack-zzzzzzzz`, `packXYZ`, `xpack-abc` ⇒ **je Exit 2**. Es zählen also
   **beide** Namensteile. **`git` selbst liest weiter** — es enumeriert
   `*.idx` unabhängig vom Namen.

   **Auslöser in freier Wildbahn:** `git maintenance run --task=loose-objects`
   schreibt `loose-<Hash>.pack`. Danach meldet dieses Target — gemessen an
   diesem Repo am 2026-09-08 — `HEAD-Tree nicht lesbar: object not found`.
   **Das Symptom variiert mit dem Zufall**, welche Objekte noch lose auf der
   Platte liegen: Im Probe-Repo lautete es `staged-Basis "HEAD" nicht
   auflösbar: reference not found`. Wer nur eine der beiden kennt, erkennt die
   andere nicht wieder.

   **Abhilfe:** `git repack -A -d` (die `-A`-Form, damit unerreichbare Objekte
   lose werden statt verworfen). **Die entlastende Hälfte gehört dazu:** Exit
   **2** ist fail-closed — es gibt hier **kein stilles Grün**, der Lauf
   verweigert die Aussage, statt eine falsche zu treffen.

   **Die Reichweite ist gemessen, nicht geschätzt:** betroffen sind die Module,
   die die **Objektdatenbank** lesen — `vcs` (dieses Target) und `commits`
   ([`make trace-check`](trace-check.md)), beide mit Exit 2. **`tracked` ist
   nicht betroffen**, weil es den git-**Index** liest: Es meldete unter beiden
   Pack-Namen denselben `target-untracked`-Befund — eine Positiv-Kontrolle,
   kein grüner Lauf, der auch *„nichts geprüft"* heißen könnte. **Über andere
   Stellen des Produkts sagt das nichts** — drei Module sind geprüft, nicht
   alle. Permanent, solange die Module über go-git lesen; ein Bump der
   Abhängigkeit kann es verschieben und ist der Anlass, hier nachzumessen
   *(seit slice-218)*.

## Bindung

**nicht** Teil von `gates`/`ci` — Diff-/Commit-Zeit-Bindepunkt.
[ADR-0016](../../docs/plan/adr/0016-adr-immutable-gate.md) ·
[ADR-0024](../../docs/plan/adr/0024-vcs-immutable-gate.md) ·
[ADR-0025](../../docs/plan/adr/0025-codepaths-ignore-refs.md)
