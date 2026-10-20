# ADR-0109: Der Push-Lauf der CI fällt bei nicht auflösbarer Basis auf die Default-Branch-Range zurück

**Status:** Proposed

**Datum:** 2026-10-20

**Autor:** pt9912

**Bezug:** [ADR-0013](0013-pr-ci-und-traceability-gate.md) (die PR-/Push-CI und
ihr Traceability-Gate, dessen Range-Bildung hier geändert wird);
[ADR-0027](0027-commits-traceability-modul.md) (`trace-check` über eine
Commit-Range);
[`DC-FA-COMMITS-001`](../../../spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in)
(die Prüfung über eine Range);
[ADR-0067](0067-dependabot-als-hebender-kanal.md) (der Kanal, dessen Rebases den
Fall auslösen);
[`AGENTS.md`](../../../AGENTS.md) §3.6 (Prüfregeln werden nur per ADR geändert);
slice-275 <!-- d-check:status-provenance -->.

**Schärft:** — *(keine Spec-Stelle; eine Entscheidung über die CI-Konfiguration)*

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Schritt „Traceability + ADR-Immutable + Slice-Closure-Übergang" der CI bildet
seine Commit-Range bei einem Push als `github.event.before..github.sha`. Nach
einem Force-Push — jeder Rebase eines Dependabot-PR ist einer — ist `before` ein
Commit, den der Klon nicht mehr kennt. `history-range-guard` meldet die Range
dann zu Recht als nicht auflösbar (Exit 2, [`SPEC-102`](../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)), und der Push-Lauf wird
rot, obwohl keine geprüfte Eigenschaft verletzt ist. Der Pflicht-Check trägt den
Namen des Jobs, nicht des Ereignisses: der rote Push-Lauf sperrt den Merge eines
PR, dessen eigener Lauf grün ist. Ein Merge ist dann nur noch mit
Administrator-Rechten möglich, also durch Umgehung der Branch-Policy.

## Entscheidung

### 1. Eine nicht auflösbare Push-Basis fällt auf `origin/<default>..HEAD` zurück

Ist `github.event.before` kein auflösbarer Commit, prüft der Push-Lauf gegen den
Default-Branch — wie er es schon bei einem neuen Branch tut. Die Range ist damit
in der Regel größer als die des Pushs: sie umfasst alles, was der Branch gegenüber
dem Default-Branch trägt — nach einem Rebase auf einen neueren Default-Stand fallen
Commits heraus, die der Default-Branch schon enthält. Die übrigen Fälle behalten ihre Range (Pull-Request:
`base..head`, neuer Branch: `origin/<default>..head`, normaler Push:
`before..head`).

### 2. Auf dem Default-Branch bleibt die nicht auflösbare Basis ein Fehler

Dort ist `origin/<default>..HEAD` leer — der Rückfall ergäbe ein Grün über leerem
Prüfbereich. Ein Push auf den Default-Branch ohne auflösbare Basis endet deshalb
mit Exit 2; ebenso ein Rückfall ohne `origin/<default>`.

### 3. Die Auflösung liegt in einem Skript, nicht im Workflow

`tools/harness/ci-range.sh` trägt die Fallunterscheidung; die Probe
`tools/harness/ci-range-probe.sh` (`make ci-range-probe`, kein Gate) fährt sie
gegen ein Wegwerf-Repo, je Fall eine Probe. Im Workflow-YAML wäre dieselbe Logik
nur gegen eine echte CI-Ausführung prüfbar.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| Den Push-Lauf für Dependabot-Branches auslassen (`branches-ignore`) | kein roter Lauf | schaltet das Gate für genau die Commits ab, die es prüfen soll; ein Branch-Name wäre die Bedingung |
| `history-range-guard` lockern (nicht auflösbar ⇒ Exit 0) | kein Eingriff in die CI | macht aus jedem fehlenden Klon-Zustand ein stilles Grün ([`SPEC-102`](../../../spec/spezifikation.md#7-festlegungen-der-harness-werkzeuge)); der Wächter existiert, um das zu verhindern |
| Den Pflicht-Check auf den Pull-Request-Lauf einschränken | der Check zählt nur den passenden Lauf | liegt in den GitHub-Einstellungen, ist aus dem Klon nicht auditierbar (`harness/README.md` §Gate-Taxonomie) |
| Mit `--admin` mergen | sofort | umgeht die Policy, statt die Ursache zu beheben; jeder nächste Rebase wiederholt es |
| Gewählt: Rückfall auf die Default-Branch-Range, laut auf dem Default-Branch | prüft mehr statt weniger; keine Umgehung | ein Feature-Branch mit alter, schon geprüfter Historie prüft sie erneut und kann Altlast melden |

**Fitness-Funktion:**

- `make ci-range-probe` fährt fünf Fälle und eine Rückfall-Lücke: Pull-Request,
  neuer Branch, normaler Push, Force-Push auf einen Feature-Branch (Range
  `origin/<default>..head`), Force-Push auf den Default-Branch (Exit 2),
  Force-Push ohne `origin/<default>` (Exit 2).
- Die Probe wird rot, wenn der Rückfall entfällt (die Range bleibt
  `before..head`) und wenn die Default-Branch-Ausnahme entfällt (der Rückfall
  liefert dort eine leere Range) — beides am Skript gebrochen und gesehen.

## Konsequenzen

- Ein Dependabot-Rebase färbt den Push-Lauf nicht mehr rot; der Merge geht ohne
  Umgehung der Branch-Policy.
- Der Push-Lauf eines Force-Pushs prüft in der Regel mehr Commits (alle seit dem
  Default-Branch). Meldet `trace-check` dabei Altlast, ist es kein Fehler des
  Rückfalls, sondern der Historie des Branchs.
- **Grenze:** Geprüft wird die Auflösbarkeit des `before`-Commits, nicht seine
  Eignung als Basis; ob die Workflow-Datei das Skript ruft, liest keine Probe.
  Ein Repo ohne Remote `origin` bekommt im Rückfall keine Range (Exit 2).

## Re-Evaluierungs-Trigger

- Ein Force-Push auf den Default-Branch wird erlaubt (Ruleset `main-nur-admin`
  geändert) — dann braucht Entscheidung 2 eine neue Antwort.
- GitHub liefert für Force-Pushs eine auflösbare `before`-Basis.
- Der Rückfall meldet bei einem Adopter wiederholt Altlast.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-10-20 | Entwurf zu slice-275 nach dem Befund beim Merge des Dependabot-PR #9: der Rebase machte den Push-Lauf rot und sperrte den Merge. Status Proposed, bis die Umsetzung geprüft ist. |
