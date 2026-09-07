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

## Bindung

**nicht** Teil von `gates`/`ci` — Diff-/Commit-Zeit-Bindepunkt.
[ADR-0016](../../docs/plan/adr/0016-adr-immutable-gate.md) ·
[ADR-0024](../../docs/plan/adr/0024-vcs-immutable-gate.md) ·
[ADR-0025](../../docs/plan/adr/0025-codepaths-ignore-refs.md)
