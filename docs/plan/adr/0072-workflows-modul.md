# ADR-0072: Der Workflow-Wächter wird das Modul `workflows`

**Status:** Accepted

**Datum:** 2026-08-29

**Autor:** pt9912

**Bezug:** [`DC-FA-WF-001`](../../../spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in)
(die neue Anforderung),
[ADR-0071](0071-lokale-workflow-referenz-rechte-pruefung.md) und
[ADR-0068](0068-lokale-workflow-referenzen-ohne-pin.md) (die Entscheide, deren
Mechanik hierher wandert),
[`AGENTS.md`](../../../AGENTS.md) §3.9 (die Hard Rule, die das Modul trägt) und
§3.8 (die Grenz-Frage, die es beantworten muss),
[ADR-0025](0025-codepaths-ignore-refs.md) (das Tombstone-Muster);
die fünf Präzedenzfälle derselben Bewegung —
[ADR-0024](0024-vcs-immutable-gate.md),
[ADR-0026](0026-completeness-in-product-gate.md),
[ADR-0027](0027-commits-traceability-modul.md),
[ADR-0028](0028-planning-lifecycle-modul.md),
[ADR-0031](0031-targets-deklarations-konsistenz-modul.md)

**Schärft:** [`DC-FA-WF-001`](../../../spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in)

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v5.12.0/regelwerk/modul-04-adrs.md).

---

## Kontext

Der Wächter über die `uses:`-Referenzen war ein `bash`-Skript. Er hat sich
bewährt — und zweimal gezeigt, wo eine Textsuche an ihre Grenze kommt:

1. **[ADR-0071](0071-lokale-workflow-referenz-rechte-pruefung.md) musste die
   `permissions:`-Blöcke mit `awk` über die Einrückung zerlegen.** Was die
   Block-Form nicht trug — `read-all`, ein Flow-Mapping mit Inhalt, ein
   Anker —, meldete das Skript als `uses-local-perms-unreadable`: fail-closed
   richtig, aber im Zweifel ein Falsch-Positiv.
2. **Die Prüfung galt nur diesem Repo.** Die Schwester-Repos konsumieren
   `d-check` über `d-check.mk`; ein Skript in `tools/harness/` erreicht sie
   nicht.

**Das Repo hat diese Bewegung fünfmal gemacht** — ein bewährtes Gate-Skript
wird ein Modul: [ADR-0024](0024-vcs-immutable-gate.md) (`vcs`),
[ADR-0027](0027-commits-traceability-modul.md) (`commits`),
[ADR-0028](0028-planning-lifecycle-modul.md) (`planning`),
[ADR-0031](0031-targets-deklarations-konsistenz-modul.md) (`targets`) und
[ADR-0026](0026-completeness-in-product-gate.md) (in-Produkt-Flag).

**Ein Einwand aus [ADR-0071](0071-lokale-workflow-referenz-rechte-pruefung.md)
Option D wird hier ausdrücklich zurückgenommen.** Dort stand als Contra:
*„ein Go-Modul im Produkt wäre die falsche Anforderung — `d-check` prüft
Dokumentation gegen Zustand, nicht Actions-Semantik"*. Das war schon beim
Schreiben nicht wahr: `targets` liest **Makefiles**, `commits` und `vcs` lesen
**git**, `tracked` den **git-Index**, und beide READMEs führen als Zusage
*„Markdown-Referenzen … plus Deklarationen (Build-Targets, Versions-Pins,
Commit-Trace, Planning)"*. Ein `uses:`-Eintrag ist eine Deklaration wie ein
Make-Target.

## Entscheidung

Die Mechanik wandert in das Modul `workflows`
([`DC-FA-WF-001`](../../../spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in)),
das Skript entfällt.

1. **Ein Modul für beide Fragen an dieselbe Zeile.** Pin-Form (§3.9) und
   Rechte-Deckung ([ADR-0071](0071-lokale-workflow-referenz-rechte-pruefung.md))
   lesen denselben Eintrag und teilen die Datei-Auswahl; zwei Module hießen,
   `.github/workflows/` zweimal zu durchlaufen und die Aussagen getrennt
   verfallen zu lassen.
2. **Die Referenzen kommen aus dem YAML-Baum.** `gopkg.in/yaml.v3` ist seit
   [ADR-0003](0003-config-format.md) Dependency. Damit entfällt die
   Block-Form-Näherung: `read-all`/`write-all`, Flow-Mappings und Anker sind
   **gelesen** statt gemeldet, und die Zeilennummer stammt aus dem Knoten.
   `uses-local-perms-unreadable` entfällt ersatzlos; an seine Stelle tritt
   `workflow-unparsable` für den Fall, der wirklich unlesbar ist.
3. **Die Scan-Menge ist konfigurierbar** (`workflows.dir`), nicht verdrahtet.
   `.github/workflows` ist eine GitHub-Konvention; ein fest kodierter Pfad
   machte das Modul für jeden Adopter mit anderer Ablage wertlos. `dir` ist
   zugleich der **Aktivierungs-Schalter** — ohne ihn wird keine Datei geöffnet.
4. **Die §3.8-Frage ist in der Anforderung beantwortet, nicht im Kommentar
   versteckt.** Das Modul **scannt** die Dateien in `workflows.dir` und
   **liest** darüber hinaus die Ziele lokaler Referenzen, auch außerhalb dieses
   Verzeichnisses. Dort gilt **dieselbe** Zusage: sie werden geparst, und ein
   Parse-Fehler ist ein Befund. Ohne diese Aussage erbte das Modul genau den
   Defekt, gegen den es gebaut ist.
5. **Nur der Job-Level-`uses` trägt die Rechte-Frage.** Ein `uses:` in einem
   **Step** ist eine Action; sie erbt die Job-Rechte und deklariert nichts.
   Diese Unterscheidung konnte das Skript nicht treffen — sie fällt mit dem
   Baum ab.
6. **Verhaltens-Parität vor dem Tombstone, gemessen.** Modul und Skript liefen
   gegen denselben Bestand und gegen den Retro-Stand vor `4681835`; die
   Ergebnisse stehen unten.
7. **Das Skript geht pfad-stabil ins Tombstone-Register**
   ([ADR-0025](0025-codepaths-ignore-refs.md)-Muster). Sein Pfad steht in zwei
   **immutablen** ADRs, einem `done/`-Slice und drei Review-Reports — Belegen,
   die den Wächter zu ihrer Zeit geprüft haben. Sie nachzuziehen verfälschte
   sie; `ignore-refs` nimmt sie quell-skopiert aus.
8. **Das Parsen liegt hinter einem Port, nicht im Kern** — und diesen Schnitt
   hat `make arch-check` erzwungen, nicht die Planung. Der erste Bau legte
   `gopkg.in/yaml.v3` in `core/rules`; a-check meldete `app-impurity`. Die
   Auflösung ist der Port `driven.WorkflowParser` mit dem Adapter
   `internal/adapter/driven/workflowyaml/` — dieselbe Ansiedlung, die
   [ADR-0009](0009-yaml-im-report-adapter.md) für den report-Adapter getroffen
   hat. **Die yaml-Allowlist in [`.a-check.yml`](../../../.a-check.yml) wächst
   dafür von zwei auf drei Adapter.** Das ist eine bewusste Erweiterung der
   Regel R3, keine Lockerung: der Kern bleibt rein, die Ausnahme bleibt
   adapter-lokal, und `composition_root: forbid` gilt unverändert.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v5.12.0/regelwerk/modul-04-adrs.md)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, das Skript behalten | kein Eingriff; es funktioniert und ist retro belegt | die Block-Form-Näherung bleibt (ein `read-all` meldet als unlesbar), und die Prüfung erreicht kein Schwester-Repo |
| B — nur die YAML-Zerlegung im Skript verbessern | kleiner Eingriff | genau die vierte Toolchain, die [`MR-046`](../../../harness/conventions.md#mr-046) ausschließt — ein YAML-Parser in `bash` ist keiner |
| C — zwei Module (Pin und Rechte getrennt) | jede Zusage einzeln abschaltbar | zweiter Durchlauf über dasselbe Verzeichnis; die Rechte-Frage entsteht ohnehin nur dort, wo der Pin-Wächter die Referenz schon aufgelöst hat |
| D — Modul mit verdrahtetem `.github/workflows` | eine Konfigurationszeile weniger | für jeden Adopter mit anderer Ablage wertlos — und d-check ist ein Werkzeug für fremde Repos, nicht nur für dieses |
| **E — ein Modul, YAML-Baum, konfigurierbare Scan-Menge, Skript-Tombstone (gewählt)** | beide Zusagen an einer Stelle; keine Näherung mehr; über `d-check.mk` verteilbar; die §3.8-Grenze steht in der Anforderung | neue Anforderung samt Spec, Codes, Tests und Doku — deutlich mehr als die Skript-Erweiterung; und ein Modul ist eine Zusage an Fremde, nicht nur an dieses Repo |

## Konsequenzen

- **Positiv, gemessen (Parität):** auf dem heutigen Bestand melden Modul und
  Skript beide **null** Befunde. Gegen den Retro-Stand vor `4681835` melden
  **beide** `uses-local-perms-undeclared` auf `.github/workflows/release.yml:268`
  — gleicher Code, gleiche Datei, gleiche Zeile. Die einzige Differenz: das
  Modul führt den Ziel-Pfad im `target`, das Skript nur in der Meldung.
- **Positiv:** `uses-local-perms-unreadable` entfällt. Eine Schreibweise, die
  das Skript rot machte, ohne dass etwas falsch war, kann das nicht mehr.
- **Positiv:** die Prüfung ist über `d-check.mk` verteilbar — die
  Schwester-Repos bekommen sie, ohne ein Skript zu kopieren.
- **Negativ, benannt:** eine Anforderung ist ein Vertrag mit Adoptern. Was hier
  als CI-Hygiene begann, muss ab jetzt seine Grenzen schärfer führen als ein
  Skript — insbesondere *„eine Deklarations-Klasse, nicht die Lauffähigkeit"*.
- **Negativ, benannt:** der Bestand trägt **eine** lokale Workflow-Referenz.
  Alles, was das Modul über die Klasse behauptet, ist an einem Exemplar plus
  konstruierten Tests geeicht ([`BEO-011`](../planning/observations.md)).
- **Negativ, benannt:** ein Parse-Fehler in einer Workflow-Datei macht jetzt
  den **inneren Loop** rot, nicht erst die CI. Das ist gewollt und trotzdem ein
  Preis ([`BEO-018`](../planning/observations.md)).

## Fitness Function (falls maschinell prüfbar)

`make workflow-pins` (in `make gates`) fährt das Modul über die eigenen
Workflows — Dogfooding. Dazu die getippten Tests in
`internal/hexagon/core/rules/workflows_test.go`:

- `TestWorkflowsHappyPath` — gepinnte fremde plus gedeckte lokale Referenz.
- `TestWorkflowsPinForm` — beweglicher Tag und SHA ohne Tag-Kommentar, jeweils
  auf **ihrer** Zeile.
- `TestWorkflowsPermsUndeclared` — der belegte Ausfall: Job erbt vom Kopf.
- `TestWorkflowsPermsNarrow` — `read` gegen `write` und der beim Aufrufer
  **fehlende** Scope (= `none`).
- `TestWorkflowsReadAll` — `read-all` deckt `contents: read` und deckt
  `contents: write` **nicht**.
- `TestWorkflowsFremdesRepoNurPin` — eine Referenz in ein fremdes Repository
  wird auf den Pin geprüft, nicht auf Rechte.
- `TestWorkflowsLeereMenge` / `TestWorkflowsUnparsable` — beide fail-closed.
- `TestWorkflowsInert` — ohne `dir` wird keine Datei geöffnet.

**Was keine Fitness Function prüft:** ob die Grenze *„eine Deklarations-Klasse,
nicht die Lauffähigkeit"* von Adoptern auch so gelesen wird. Das ist der Grund,
warum sie in der Anforderung, im Handbuch und in beiden READMEs steht.

## Re-Evaluierungs-Trigger

**Der erste Adopter, dessen Workflow-Ablage nicht `.github/workflows` heißt.**
Dann zeigt sich, ob `dir` als einzelnes Verzeichnis reicht oder ob die
Scan-Menge ein Glob werden muss — die heutige Form ist an **einer** Ablage
geeicht.

**Zweiter Trigger: ein `workflow-unparsable` an einer Datei, die das CI-System
akzeptiert.** Dann liest `yaml.v3` strenger als der Workflow-Läufer, und die
Zusage „unlesbar" ist die falsche — zu entscheiden ist dann, ob das Modul die
Toleranz nachbildet oder die Grenze benennt.

**Dritter Trigger: eine zweite lokale Workflow-Referenz im Repo.** Solange es
**eine** gibt, ist jede Aussage über die Klasse an einem Exemplar geeicht.
