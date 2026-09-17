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

3. **Pack-Dateien ohne gültiges Hash-Suffix oder ohne passenden Index sind
   unsichtbar — hier bricht der Lauf ab.** Das Modul liest die
   Objektdatenbank über go-git (`v5.19.2`), und go-gits `DotGit`-Schicht
   findet einen Pack nur unter git's kanonischem Namen
   `pack-<Hash-des-Packs>.{idx,pack}` (gemessen; `git` selbst liest weiter).
   Ein Pack ohne gültiges Hash-Suffix meldet
   `Range-Basis "<base>" nicht auflösbar: reference not found` oder
   `commit <sha> nicht lesbar: object not found` — beides **Exit 2**.
   **Abhilfe:** `git repack -A -d`.

   **Seit [`slice-226`](../../docs/plan/planning/done/slice-226-vcs-pack-alias-fremdes-praefix.md)
   <!-- d-check:status-provenance --> löst der Adapter einen Pack zusätzlich
   unter seinem kanonischen Namen auf, wenn seine Datei einen gültigen
   SHA1/SHA256-Hash als Namens-Suffix trägt **und** eine passende
   `.idx`-Datei existiert.** Schreibt etwa
   `git maintenance run --task=loose-objects` einen `loose-<Hash>.pack`, liest
   dieses Target ihn jetzt wie jeden anderen — kein `Exit 2` mehr allein wegen
   des Präfixes; die Abhilfe oben bleibt nur für einen Pack ohne gültiges
   Hash-Suffix oder ohne passenden Index nötig.

   **Dieses Target war dabei immer fail-closed** — ein unlesbarer Commit
   bricht den Lauf ab (gemessen). [`make adr-check`](adr-check.md#grenze--was-das-grün-nicht-abdeckt)
   war es **nicht**: Bis slice-218 meldete es im `RANGE=`-Modus still grün,
   wenn der unsichtbare Pack nur einzelne Objekte verschluckte, und bis
   slice-220 blieben zwei weitere Ausprägungen offen bzw. fehldiagnostiziert
   — beide Fixe liegen dort. Für Konsumenten, die ein Image `v0.76.0` oder
   älter pinnen, führte
   [`CO-001`](../../docs/plan/carveouts/done/CO-001-vcs-range-stiller-skip.md)
   (inzwischen aufgelöst) die Reststrecke. Die vollständige Messung steht bei
   `adr-check`, damit sie an *einem* Ort gepflegt wird; dass sie für dieses
   Target **günstiger** ausfiel, ist gemessen und nicht angenommen
   *(seit slice-218)*.

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
