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
   lose werden statt verworfen).

   **Und jetzt der schlimmere Teil: im `RANGE=`-Modus meldet dieses Target
   dann still grün.** Wie der unsichtbare Pack wirkt, hängt davon ab, *was* in
   ihm liegt — und die beiden Fälle sind nicht gleich schlimm:

   | Was der unsichtbare Pack verschluckt | Verhalten |
   | --- | --- |
   | eine ganze Referenz (BASE-Commit oder -Tree) | **Exit 2**, fail-closed |
   | **einzelne Objekte** — etwa nur den BASE-Blob einer ADR | **Exit 0, 0 Befunde** |

   **Gemessen, und die Verletzung war echt:** dieselbe Range, dieselbe Datei,
   eine geänderte ADR-Kern-Zeile. Mit kanonischem Pack-Namen
   `1 Datei(en) geprüft, 1 Befund(e)` (`core-drift-vcs`, Exit 1); nach dem
   Umbenennen **desselben** Packs `1 Datei(en) geprüft, 0 Befund(e)`, Exit 0 —
   `git diff` zeigt die Änderung unverändert an. **Der Grund liegt in einer
   Zusammenfassung zweier Zustände:** go-git meldet ein *unlesbares* Blob als
   `ErrFileNotFound`, also mit demselben Fehler wie *„diese Datei gibt es in
   diesem Tree nicht"* — ein völlig legitimer Zustand (eine später angelegte
   Datei). Der Adapter bildet das auf `ok=false, err=nil` ab, und
   `vcsModified` überspringt die Datei dann ohne Befund.

   **Das ist keine bloße Grenze, sondern ein Defekt**, und er trifft den
   Bindepunkt, der *nicht* opt-in ist: Die PR-/Push-CI fährt `RANGE=`
   ([`ci.yml`](../../.github/workflows/ci.yml)), der fail-closede
   `--staged`-Pfad ist der lokale Hook. **Bis er behoben ist, ersetzt kein
   grüner `adr-check`-Lauf die Gewissheit, dass die Objektdatenbank kanonisch
   gepackt ist.**

   **Die Reichweite ist gemessen, nicht geschätzt** — und die Module verhalten
   sich **unterschiedlich**, was drei Formulierungen hier zuvor gleichsetzten:
   `commits` ([`make trace-check`](trace-check.md)) **ist** fail-closed (ein
   unlesbarer Commit ⇒ `commit … nicht lesbar`, Exit 2, gemessen); `vcs` ist es
   nur in den zwei Zeilen der Tabelle oben. **`tracked` ist nicht betroffen**,
   weil es den git-**Index** liest: Es meldete unter beiden Pack-Namen denselben
   `target-untracked`-Befund — eine Positiv-Kontrolle, kein grüner Lauf, der
   auch *„nichts geprüft"* heißen könnte. **Über andere Stellen des Produkts
   sagt das nichts** — drei Module sind geprüft, nicht alle. Permanent, solange
   die Module über go-git lesen; ein Bump der Abhängigkeit ist der Anlass,
   hier nachzumessen — **und zwar mit einem *partiellen* unsichtbaren Pack**,
   denn nur der deckt den stillen Pfad auf *(seit slice-218)*.

## Bindung

**nicht** Teil von `gates`/`ci` — Diff-/Commit-Zeit-Bindepunkt.
[ADR-0016](../../docs/plan/adr/0016-adr-immutable-gate.md) ·
[ADR-0024](../../docs/plan/adr/0024-vcs-immutable-gate.md) ·
[ADR-0025](../../docs/plan/adr/0025-codepaths-ignore-refs.md)
