# ADR-0024 — Git-Diff-Immutabilität: das Modul `vcs` prüft den Core über eine Commit-Range und löst die Skript-Mechanik von ADR-0016 ab

**Status:** Proposed
**Datum:** 2026-06-29
**Autor:** pt9912
**Bezug:** [`DC-FA-VCS-001`](../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
(opt-in Modul `vcs`); Schwester-Hälfte/Vorläufer [ADR-0023](0023-immutable-core-pin.md)
(hermetischer `immutable`-Pin — die git-Stufe, die ADR-0023 bewusst vertagte);
vorgezeichnet von [ADR-0008](0008-reparatur-ableitbarkeit.md) (künftiges
„VCS-Modul" als nicht-hermetischer Driven-Port, analog `external`); **Supersedes
die Skript-Mechanik** von [ADR-0016](0016-adr-immutable-gate.md)
(`adr-immutable-check.sh` → d-check-Modul `vcs`; die **Policy und das Gate
`adr-check` bleiben** unverändert — Teil-Supersede, wie [ADR-0002](0002-distribution-ghcr-image.md)/[ADR-0014](0014-latest-tag-fuer-stabile-releases.md));
Port-/Ordner-Konvention [ADR-0005](0005-modul-layout-hexagon-ordner.md);
Distroless-Vertrag [ADR-0002](0002-distribution-ghcr-image.md) (bleibt — **kein**
git-Binary); Pin-Disziplin der neuen Dependency [ADR-0011](0011-digest-pins-build-gate-images.md);
Mechanik-Verwandtschaft [ADR-0013](0013-pr-ci-und-traceability-gate.md) (CI-Range +
lokaler Hook, eine Quelle).
**Schärft:** die neue Algorithmus-Sektion
[`spec/spezifikation.md` §DC-FA-VCS-001.a](../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
(Range-Semantik, Core-Definition, Grund-Code `core-drift-vcs`) und die
VCS-Port-Rolle in [`spec/architecture.md` §1–§2](../../../spec/architecture.md#1-komponenten-übersicht).

## Kontext

Die Immutabilität von `Accepted`-ADRs hat heute **zwei** Wächter, die dieselbe
Frage aus zwei Richtungen stellen:

- [ADR-0016](0016-adr-immutable-gate.md) — das Shell-Skript `adr-immutable-check.sh`
  vergleicht über eine git-Range den **Core** jeder geänderten Accepted-ADR
  (`core(BASE)` vs. `core(HEAD)`). **Harte** Garantie, fängt jede Körper-Änderung
  am Integrationspunkt — aber als **kopiertes Skript** trägt es das
  Copy-Drift-Problem der Werkzeug-Familie ([`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Klasse).
- [ADR-0023](0023-immutable-core-pin.md) — das hermetische Modul `immutable`
  prüft denselben Core gegen einen **im Dokument hinterlegten Pin**, im
  read-only-Arbeitsbaum, **ohne git**. Verteilbar im Image (löst Copy-Drift),
  aber die Garantie ist **schwächer** (neu-pinn-bar → Reviewer als Boden).

[ADR-0023](0023-immutable-core-pin.md) entschied bewusst **nur** die hermetische
Pin-Hälfte und hielt fest: „Die git-Stufe ist nicht verworfen, nur vertagt … ein
neuer nicht-hermetischer Driven-Port (ein `git`-Adapter neben
`httpcheck`/`external`) als eigene Anforderung/ADR." Genau diese Stufe löst diese
ADR ein. [ADR-0008](0008-reparatur-ableitbarkeit.md) hatte den Port schon benannt
(„ein künftiges VCS-Modul … als opt-in, analog `external`").

Die offene Spannung war: d-check läuft ausschließlich über das **distroless-Image**
([ADR-0002](0002-distribution-ghcr-image.md)) — kein Host-Go, **kein git-Binary**,
`--network none`. Damit das Modul die git-Garantie **verteilen** kann (statt sie als
Skript zu kopieren), muss der Adapter `.git` **im Image** lesen, ohne ein git-Binary
hinzuzufügen.

## Entscheidung

Ein opt-in Modul `vcs` (Default aus) mechanisiert die git-historienbasierte
Immutabilität als **Core-Diff über eine Commit-Range**: `core(BASE)` ≟ `core(HEAD)`
für jede in der Range geänderte Datei der Klasse `vcs.paths`, deren BASE die
Bedingung `vcs.immutable-when` erfüllt. Abweichung, unzulässiger Status-Übergang
(`vcs.head-allow`) oder Löschung/Umbenennung → `core-drift-vcs`. Die Core-Semantik
ist in **Parität** zum abgelösten Skript: `vcs.exclude-sections` strippt
Abschnitte (`Geschichte`), `vcs.status-line` strippt **nur die Kopf-Status-Zeile**
(eine gleichlautende Körper-Zeile bleibt Core). **Diagnose-only**, kein
`--repair`-Hunk.

**Neuer Driven-Port + git-Adapter, reine-Go.** Der Kern bekommt einen
**VCS-Port** (`FileAtRef(ref, pfad)` / `ChangedFiles(base, head)`); die **einzige**
git-berührende Stelle ist ein neuer **Driven-Adapter** `git`, exakt parallel zu
`httpcheck` für `external`. Der Adapter liest `.git` über eine **reine-Go**
git-Objekt-Bibliothek — **kein** git-Binary, das **distroless-Image bleibt
unangetastet** ([ADR-0002](0002-distribution-ghcr-image.md)). `.git` ist im
bestehenden read-only-Mount (`-v <repo>:/repo:ro`) bereits vorhanden; der
Mount-Vertrag ändert sich **nicht**. `make arch-check` (R1–R5) hält die
git-Bibliothek im Adapter isoliert; der Kern bleibt rein und fakebar (ein
Fake-VCS-Port in den Tests).

**Lokal/lesend/deterministisch, nicht „unhermetisch wie `external`".** Anders als
[ADR-0008](0008-reparatur-ableitbarkeit.md)/[ADR-0023](0023-immutable-core-pin.md)
es vorzeichneten (Port „analog `external`"), relativiert `vcs` **weder**
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus) **noch**
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit):
`external`s Eingabe ist der **nicht-reproduzierbare Remote** (Netz); `vcs`'
Eingabe ist das **lokale `.git`** — reproduzierbar, read-only, netzlos. Der einzige
Unterschied zu den hermetischen Modulen ist der **erweiterte Eingabe-Scope** (git
+ Range statt nur des gescannten Markdown-Baums). Darum ist `vcs` **strikt opt-in**
und **fail-closed**: fehlt `.git` oder die Range, bricht der Lauf laut ab (Exit 2),
statt still grün zu sein.

**Dogfood: `adr-check` läuft auf das Modul um.** `make adr-check`, der
`pre-commit`-Hook und die PR-/Push-CI rufen künftig **das d-check-Image** mit
`--enable vcs` (`--range <base>..<head>` bzw. `--staged`) statt
`bash tools/adr-immutable-check.sh`. Die `vcs`-Klasse (ADR-Pfade,
`immutable-when: Status Accepted`, `exclude-sections: Geschichte`, `status-line`,
`head-allow`) liegt als `vcs:`-Block in [`.d-check.yml`](../../../.d-check.yml),
bleibt aber **aus der Default-`modules`-Liste**: `make doc-check` aktiviert `vcs`
nicht → der Default-Lauf ist byte-identisch
([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)). d-check isst
damit für seine **eigenen** Accepted-ADRs sein eigenes Futter — die Policy von
[ADR-0016](0016-adr-immutable-gate.md) bleibt, nur die Mechanik (Skript → Modul)
wechselt.

**Skript bleibt pfad-stabil.** `tools/adr-immutable-check.sh` wird **nicht
gelöscht**: [ADR-0016](0016-adr-immutable-gate.md) ist `Accepted`/immutable und
referenziert das Skript in Inline-Code — eine Löschung bräche `make doc-check`
(`codepath-missing`) an einer unveränderlichen Datei. Das Skript bleibt als
funktionierende Referenz/Fallback erhalten, ist aber aus dem Gate ausgehängt. Sein
Negativ-Selbsttest wandert als Akzeptanztest ins Modul (`make test`).

**Eigenes Modul** (nicht Strategie in `immutable`): anderer Bindepunkt
(Commit-Range statt Arbeitsbaum), anderer Eingabe-Scope (git statt Datei),
anderer Lebenspunkt (Integration/PR statt `make gates`). `immutable` bleibt sauber
hermetisch; beide koexistieren als **Defense-in-Depth** — der Pin als Boden
überall (`make gates`), der git-Diff als harter Deckel an der Integration
(`adr-check`).

## Verglichene Alternativen

| Alternative | Pro | Contra |
| --- | --- | --- |
| **Modul `vcs` + reine-Go-VCS-Port (gewählt)** | volle Range-Garantie; verteilt im Image (löst Copy-Drift); distroless bleibt (kein git-Binary); git im Adapter isoliert (`arch-check`); lokal/deterministisch/read-only | neue große Go-Dependency (Supply-Chain, Pin-Disziplin [ADR-0011](0011-digest-pins-build-gate-images.md), Image-Größe); erweiterter Eingabe-Scope (`.git`) |
| git-**Binary** ins Runtime-Image | vertrauter `git`-Exec | bräche distroless/static ([ADR-0002](0002-distribution-ghcr-image.md)); größere Angriffsfläche; eigener Supersede an ADR-0002 |
| Skript belassen, nicht mechanisieren | volle Garantie, schon da | Copy-Drift in Schwester-Repos bleibt — der Auftrag bleibt ungelöst |
| nur `immutable` (Pin), keine git-Stufe | hermetisch, schon da | Garantie schwächer (neu-pinn-bar); fängt kein getarntes Re-Pin |
| Strategie `mode: vcs` in `immutable` | ein Modul | mischt hermetisch + git-Eingabe in einem Modul; verwischt Bindepunkt und QA-Zusage |

## Konsequenzen

- **Neue Dependency.** Eine reine-Go git-Objekt-Bibliothek tritt in `go.mod`/`go.sum`
  ein; sie unterliegt der Pin-/Reproduzierbarkeits-Disziplin
  ([ADR-0011](0011-digest-pins-build-gate-images.md)) und dem `semgrep`-Scope. Der
  Import ist **ausschließlich** im `git`-Adapter erlaubt (`arch-check` R2-Klasse:
  wie `net/http` nur im HTTP-Adapter).
- **`adr-check` braucht jetzt das Image.** Der `pre-commit`-Hook und `make adr-check`
  bauen/rufen das d-check-Image statt eines reinen Shell-Skripts — schwergewichtiger
  pro Commit; dafür eine **verteilte** Wahrheit ohne Skript-Kopie. Die CI baut das
  Image ohnehin.
- **Restlücke unverändert.** Wie [ADR-0016](0016-adr-immutable-gate.md)/[ADR-0013](0013-pr-ci-und-traceability-gate.md):
  ohne Branch Protection (außerhalb des Repos) ist die CI nur *advisory*.
- **Determinismus/Read-only gehalten.** `vcs` liest `.git` nur lesend, netzlos,
  schreibt nicht ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit));
  gleiche Historie + Range ⇒ gleicher Befundsatz
  ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- **`adr-immutable-check.sh` bleibt im Baum** (Pfad-Stabilität für die immutable
  [ADR-0016](0016-adr-immutable-gate.md)-Inline-Referenz), aus dem Gate ausgehängt.

## Fitness Function

- `make adr-check` läuft **rot** bei einem Körper-Edit / Status-Rückfall an einer
  `Accepted`-ADR (bzw. Löschung/Umbenennung), **grün** sonst — wie zuvor, jetzt
  über das Modul. Die `Proposed → Accepted`-Reifung einer frischen ADR bleibt frei.
- Ohne aktives `vcs` ist der Befundsatz byte-identisch (opt-in-Selbsttest;
  `make doc-check` aktiviert `vcs` nicht).
- `make arch-check` bleibt grün: die git-Bibliothek wird **nur** im `git`-Adapter
  importiert, der Kern bleibt rein.
- fail-closed: fehlendes `.git` / fehlende oder unauflösbare Range ⇒ Exit 2 (kein
  stilles Grün) — Akzeptanztest.
- Core-Paritäts-Tests (alle **sieben** Klassen des Skript-Selbsttests): Geschichte-Anhang,
  Superseded-Übergang und ein Körper-Edit auf einer **`Proposed`-BASE**
  (`immutable-when` nicht erfüllt → Grandfathering) feuern **nicht**; Körper-Edit,
  Status-Rückfall, Edit an einer Körper-`**Status:**`-Zeile und Edit im Abschnitt
  nach `## Geschichte` feuern.
- `make semgrep`/`make lint` grün mit der neuen Dependency.

## Re-Evaluierungs-Trigger

- ADR-Vorlage ändert die Struktur (`## Geschichte`-Überschrift,
  `**Status:**`-Zeilenformat) → `vcs`-Klassen-Config (`exclude-sections`,
  `status-line`, `head-allow`) anpassen — jetzt **Config**, nicht Skript-Code.
- Bedarf, die git-Garantie auch in `make gates` (Arbeitsbaum-Bindepunkt) zu führen
  → bleibt die Domäne des hermetischen `immutable` ([ADR-0023](0023-immutable-core-pin.md));
  `vcs` ist der Integrations-Bindepunkt.
- Eine zweite git-bedürftige Befund-/Reparatur-Klasse (z. B. git-basierte
  Move-Erkennung, [ADR-0008](0008-reparatur-ableitbarkeit.md)) → nutzt denselben
  VCS-Port, eigene Anforderung.

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-06-29 | Entwurf nach Auftraggeber-Auftrag („`adr-immutable-check` regelkonform vollständig mechanisieren"). Vier Entscheidungen: eigenes Modul `vcs`; volle Parität `--range` + `--staged`; Dogfood-Ersatz des `adr-check`-Gates; git-Zugang reine-Go (distroless bleibt). Löst die von [ADR-0023](0023-immutable-core-pin.md) vertagte git-Stufe ein, Teil-Supersede der Skript-Mechanik von [ADR-0016](0016-adr-immutable-gate.md). Begleitet slice-053. Status Proposed. |
