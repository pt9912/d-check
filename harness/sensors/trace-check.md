# `make trace-check` — hält, dass jede Commit-Botschaft eine Traceability-Kennung nennt

## Vertrag

Via Modul `commits` (Image, dogfood). Jede Commit-Message nennt eine
`DC-*`/`ADR-*`/`MR-*`/`slice-*`-Kennung, sonst `commit-untraceable`.
Ausgenommen sind Merge- und Revert-Commits.

**Zwei Modi über dasselbe Modul:** `--range` für CI und Push, `--commit-msg -`
für den Hook über stdin.

## Grenze — was das Grün nicht abdeckt

1. **Geprüft ist die Nennung, nicht der Bezug** — dass eine Botschaft eine
   Kennung trägt, sagt nicht, dass der Commit sie betrifft. Permanent; das ist
   Review-Territorium.
2. **Der lokale Hook ist opt-in pro Klon** (`make hooks`); `--no-verify`
   umgeht ihn, nicht die PR-/Push-CI.

3. **Pack-Dateien mit unkanonischem Namen sind unsichtbar — hier bricht der
   Lauf ab.** Das Modul liest die Objektdatenbank über go-git (`v5.19.2`), und
   go-git findet einen Pack nur unter git's kanonischem Namen
   `pack-<Hash-des-Packs>.{idx,pack}` — **beide** Namensteile zählen (gemessen;
   `git` selbst liest weiter). Schreibt etwa
   `git maintenance run --task=loose-objects` einen `loose-<Hash>.pack`, meldet
   dieses Target `Range-Basis "<base>" nicht auflösbar: reference not found`
   oder `commit <sha> nicht lesbar: object not found` — beides **Exit 2**.
   **Abhilfe:** `git repack -A -d`.

   **`commits` ist dabei fail-closed, `vcs` nicht — die beiden sind hier
   ausdrücklich *nicht* gleich.** Ein unlesbarer Commit bricht diesen Lauf ab
   (gemessen); [`make adr-check`](adr-check.md#grenze--was-das-grün-nicht-abdeckt)
   dagegen meldet im `RANGE=`-Modus **still grün**, wenn der unsichtbare Pack
   nur einzelne Objekte verschluckt. Die vollständige Messung steht dort, damit
   sie an *einem* Ort gepflegt wird; dass sie für dieses Target **günstiger**
   ausfällt, ist gemessen und nicht angenommen *(seit slice-218)*.

**Dependabot braucht dafür keine Ausnahme:** Seine Botschaften tragen die
Kennung im Präfix und erfüllen die Regel wie jeder andere Commit
([ADR-0067](../../docs/plan/adr/0067-dependabot-als-hebender-kanal.md)) — eine
erweiterte `exempt-pattern` hätte den Gate für eine **ganze Commit-Klasse**
blind gemacht.

## Bindung

**nicht** Teil von `gates`/`ci` — Commit-Zeit-Bindepunkt; gerufen vom
`commit-msg`-Hook und der PR-/Push-CI.
[ADR-0013](../../docs/plan/adr/0013-pr-ci-und-traceability-gate.md) ·
[ADR-0027](../../docs/plan/adr/0027-commits-traceability-modul.md) ·
[`DC-FA-COMMITS-001`](../../spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in)
