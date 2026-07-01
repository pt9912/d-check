# ADR-0027 — Traceability-Kennung in Commit-Messages: das Modul `commits` mechanisiert `trace-check.sh` über den VCS-Port und löst dessen Skript-Mechanik von ADR-0013 ab

**Status:** Accepted
**Datum:** 2026-07-01
**Autor:** pt9912
**Bezug:** [`DC-FA-COMMITS-001`](../../../spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in)
(opt-in Modul `commits`); **Präzedenz** [ADR-0024](0024-vcs-immutable-gate.md)
(Modul `vcs` löste das Skript `adr-immutable-check.sh` über den reine-Go-VCS-Port ab —
dieselbe Mechanik, hier auf Commit-**Messages** statt Datei-Inhalt) und
[ADR-0026](0026-completeness-in-product-gate.md) (dritte Skript-Ablösung derselben
Familie); **Supersedes die Skript-Mechanik** von [ADR-0013](0013-pr-ci-und-traceability-gate.md)
(`tools/trace-check.sh` → d-check-Modul `commits`; die **Policy, der Bindepunkt und
die CI-Topologie** — Traceability-Regel, `commit-msg`-Hook + PR-/Push-CI — **bleiben
unverändert**, Teil-Supersede wie [ADR-0016](0016-adr-immutable-gate.md)/[ADR-0024](0024-vcs-immutable-gate.md));
Port-/Ordner-Konvention [ADR-0005](0005-modul-layout-hexagon-ordner.md); Distroless-Vertrag
[ADR-0002](0002-distribution-ghcr-image.md) (bleibt — **kein** git-Binary); Pin-Disziplin
der VCS-Dependency [ADR-0011](0011-digest-pins-build-gate-images.md); die immutable
Skript-Referenz wird ein
[`codepaths.ignore-refs`](../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code)-Tombstone
([ADR-0025](0025-codepaths-ignore-refs.md)); Verteilungs-Kern
[`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)
(„verteilen statt kopieren").
**Schärft:** die neue Algorithmus-Sektion
[`spec/spezifikation.md` §DC-FA-COMMITS-001.a](../../../spec/spezifikation.md#dc-fa-commits-001a--traceability-kennung-in-commit-messages-über-eine-commit-range-commits)
(Quellen/Modi, Bereinigung, Grund-Code `commit-untraceable`) und die um Commit-Message-Lesen
erweiterte VCS-Port-Rolle in [`spec/architecture.md` §2](../../../spec/architecture.md#2-schichten-und-constraints).

## Kontext

Die Traceability-Regel — „PRs/Commits **müssen** mindestens eine `DC-*`-, `ADR-*`-,
`MR-*`- oder `slice-*`-ID nennen"
([`harness/README.md` §Traceability rules](../../../harness/README.md#traceability-rules))
— wird heute von `tools/trace-check.sh` durchgesetzt ([ADR-0013](0013-pr-ci-und-traceability-gate.md)):
ein **kopiertes** Shell-Skript, das (a) der `commit-msg`-Hook mit `--message`, (b) die
PR-/Push-CI mit `--range` und (c) `make trace-check` lokal rufen. Es ist das **letzte**
Gate-Skript der Familie, das noch nicht mechanisiert ist — nach `adr-immutable-check.sh`
([ADR-0024](0024-vcs-immutable-gate.md)) und `completeness-check.sh`
([ADR-0026](0026-completeness-in-product-gate.md)). Als kopiertes Skript trägt es
denselben Copy-Drift über die Repo-Familie
([`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Klasse):
Schwester-Repos, die dieselbe Traceability-Invariante wollen, kopieren das Skript, statt
die Prüfung als gepinntes Image zu beziehen.

Der Antrieb ist derselbe wie bei [ADR-0024](0024-vcs-immutable-gate.md): der dort gebaute
**reine-Go-VCS-Port** (liest das read-only `.git` **ohne git-Binary**, **ohne Netz**)
macht die git-berührende Hälfte mechanisierbar und **im Image verteilbar**. `trace-check`
liest Commit-**Messages** über eine Range — genau eine Operation, die der Port heute noch
nicht anbietet, aber trivial ergänzen kann (der Diff-Pfad löst den Commit ohnehin schon
auf).

Der Name `trace` schied aus: das `--trace`-Flag und die Requirements Traceability Matrix
([`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix))
belegen ihn bereits (Doc-/Requirements-Achse). Das Modul heißt darum `commits` — es
beschreibt den **geprüften Gegenstand** (die Commits), analog `links`/`anchors`/`pins`/`versions`,
und hält `trace` eindeutig für die RTM frei. Das **Gate-Target** bleibt `make trace-check`
(etablierter CI-/Harness-Begriff) und delegiert intern an das Modul.

## Entscheidung

Ein opt-in Modul `commits` (Default aus) mechanisiert das Traceability-Gate als
d-check-Modul: es prüft, dass jede **bereinigte** Commit-Message (git-`strip`: scissors-
und `#`-Zeilen entfernt) mindestens eine Kennung nach `commits.id-patterns` auf einer
Inhalts-Zeile trägt; sonst `commit-untraceable`. Betreff-Ausnahme über
`commits.exempt-pattern` (Selbstkonfig `^(Merge |Revert )`). **Diagnose-only**, kein
`--repair`-Hunk (die Korrektur ist ein neuer Commit / ein menschliches `--amend`).

**VCS-Port um Commit-Message-Lesen erweitert.** Der [ADR-0024](0024-vcs-immutable-gate.md)-
VCS-Port bekommt eine Operation `CommitMessages(base, head)` (Nicht-Merge-Commits der
Range, `git rev-list --no-merges`-Parität); die **einzige** git-berührende Stelle bleibt
der reine-Go-`git`-Adapter — **kein** git-Binary, das **distroless-Image bleibt
unangetastet** ([ADR-0002](0002-distribution-ghcr-image.md)). `make arch-check` (R2) hält
die git-Bibliothek weiter im Adapter isoliert; der Kern bleibt rein und fakebar.

**Zwei Quellen, ein Kern.** Ein reiner Modul-Kern prüft eine Liste von Messages; zwei
Quellen speisen ihn: der **Range**-Modus (`--enable commits --range <base>..<head>`, via
Port — CI/Push und `make trace-check`) und der **Message**-Modus (`--commit-msg <datei|->`,
ein **Kurzschluss-Modus** wie `--print-config`/`--trace`, ohne Repo-Scan/Port — der
`commit-msg`-Hook pipet die Pending-Message über stdin). Die Pending-Message existiert
noch nicht als Commit, ist also über keine Range erreichbar — darum der eigene Modus;
**dieselbe** Bereinigung/Prüfung wie im Range-Modus (keine Divergenz).

**Lokal/lesend/deterministisch, strikt opt-in, fail-closed.** Wie `vcs` relativiert
`commits` **weder** [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus) **noch**
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit):
Eingabe ist das **lokale `.git`** (bzw. eine Message-Datei) — reproduzierbar, read-only,
netzlos; nur der Eingabe-Scope ist erweitert. Nie Default-Modul; ohne aktives `commits`
byte-identisch. **fail-closed:** fehlende/unauflösbare Range, unlesbares `.git`, nicht
lesbare Message-Datei → Exit 2 (kein stilles Grün).

**Dogfood: `make trace-check` läuft auf das Modul um.** `make trace-check`, der
`commit-msg`-Hook und die PR-/Push-CI rufen künftig **das d-check-Image** mit
`--enable commits` (`--range` bzw. `--commit-msg -`) statt `bash tools/trace-check.sh`.
Die `commits`-Klasse (ID-Muster + Exempt-Muster) liegt als `commits:`-Block in
[`.d-check.yml`](../../../.d-check.yml), bleibt aber **aus der Default-`modules`-Liste**:
`make doc-check` aktiviert `commits` nicht → der Default-Lauf ist byte-identisch. d-check
isst damit für seine **eigenen** Commits sein eigenes Futter — die Policy von
[ADR-0013](0013-pr-ci-und-traceability-gate.md) bleibt, nur die Mechanik (Skript → Modul)
wechselt. Die **eine Wahrheit** von [ADR-0013](0013-pr-ci-und-traceability-gate.md) Punkt 3
wird gestärkt: Hook und CI rufen dasselbe Modul im **verteilten** Image statt eines
kopierten Skripts; die ID-Definition ([ADR-0013](0013-pr-ci-und-traceability-gate.md)
Punkt 4) wandert vom hartcodierten `ID_RE` in die `commits.id-patterns` **neben** den
`ids.patterns` — beide ID-Quellen co-lokalisiert und gemeinsam reviewbar in **einer**
Datei.

**Skript entfernt, Tombstone statt Edit.** `tools/trace-check.sh` wird **entfernt** (`git
rm`) — anders als bei [ADR-0024](0024-vcs-immutable-gate.md), wo das Skript zunächst
pfad-stabil blieb und erst [ADR-0025](0025-codepaths-ignore-refs.md) es entfernte; die
Lehre ist gezogen, der Umweg entfällt. Die immutable
[ADR-0013](0013-pr-ci-und-traceability-gate.md) referenziert `tools/trace-check.sh` in
Inline-Code; die Löschung bräche `make doc-check` (`codepath-missing`) an unveränderlicher
Doku. Ein Eintrag in `.d-check.yml` `codepaths.ignore-refs` deklariert den Pfad als
entfernt — der **dritte** reale Anwendungsfall des Ventils aus
[ADR-0025](0025-codepaths-ignore-refs.md).

**Eigenes Modul** (nicht Strategie in `vcs`): anderer Gegenstand (Commit-**Message** statt
Datei-**Inhalt**), anderer Grund-Code, anderer Betreff-/ID-Semantik. `vcs` bleibt sauber
auf Datei-Immutabilität; beide teilen den VCS-Port, aber nicht die Regel.

## Verglichene Alternativen

| Alternative | Pro | Contra |
| --- | --- | --- |
| **Modul `commits` + Port-Erweiterung (gewählt)** | verteilt im Image (löst Copy-Drift); eine Wahrheit für Hook + CI; ID-Muster neben `ids` co-lokalisiert; distroless bleibt (kein git-Binary); nutzt den vorhandenen VCS-Port; lokal/deterministisch/read-only | erweiterter Eingabe-Scope (`.git`/Range); Image-Lauf pro Commit (Hook schwerer); neue Port-Operation + CLI-Modus |
| Skript belassen, nicht mechanisieren | volle Garantie, schon da | Copy-Drift in Schwester-Repos bleibt — der Auftrag bleibt ungelöst; letztes Familien-Skript |
| Modulname `trace` (wie Roadmap) | an `trace-check`-Herkunft gebunden | kollidiert mit `--trace`/RTM ([`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)) — zementiert die Mehrdeutigkeit; Harness-Ehrlichkeits-Reibung |
| Strategie `mode: commits` in `vcs` | ein Modul | mischt Datei-Inhalt + Commit-Message in einem Modul; verwischt Gegenstand und Grund-Code |
| nur Range-Modus (commit-msg-Hook fallen lassen) | ein CLI-Modus weniger | verliert die Commit-Zeit-Feedforward-Hälfte von [ADR-0013](0013-pr-ci-und-traceability-gate.md) — Gate-Lockerung ohne Not (`AGENTS.md` §3.6) |

## Konsequenzen

- **VCS-Port-Erweiterung.** `CommitMessages(base, head)` tritt zum Port hinzu; die
  go-git-Nutzung bleibt **ausschließlich** im `git`-Adapter (`arch-check` R2). Keine neue
  Dependency (go-git ist seit [ADR-0024](0024-vcs-immutable-gate.md) da).
- **`make trace-check` braucht jetzt das Image** (wie `make adr-check` seit
  [ADR-0024](0024-vcs-immutable-gate.md), `make completeness-check` seit
  [ADR-0026](0026-completeness-in-product-gate.md)) — schwerer pro Lauf/Commit; dafür eine
  verteilte, kopie-freie Wahrheit. Die CI baut das Image ohnehin.
- **`tools/trace-check.sh` entfernt;** sein Negativ-Selbsttest (ID erkannt, fehlende ID
  gefangen, Merge ausgenommen, uniforme `#`-Bereinigung) lebt als Modul-Akzeptanztest
  (`make test`) weiter — dieselbe Verlagerung wie bei
  [ADR-0024](0024-vcs-immutable-gate.md)/[ADR-0026](0026-completeness-in-product-gate.md).
  Der Verlust der **Per-Lauf**-Selbsttest-Eigenschaft (das Skript testete sich bei jedem
  Aufruf) ist derselbe bewusste Trade wie dort.
- **Release nötig.** Anders als [ADR-0026](0026-completeness-in-product-gate.md) (reiner
  Harness-Refactor) ändert diese ADR **Produkt-Code** (`internal/`, `cmd/`) — neues Modul +
  Port-Operation + CLI-Modus. Minor-Bump und GHCR-Release (v0.35.0), wie
  [ADR-0024](0024-vcs-immutable-gate.md).
- **Doc-check bleibt grün** trotz entferntem Skript: die immutable
  [ADR-0013](0013-pr-ci-und-traceability-gate.md)-Referenz ist als
  `codepaths.ignore-refs`-Tombstone deklariert ([ADR-0025](0025-codepaths-ignore-refs.md)).
- **Verteilung an Konsumenten** schließt den
  [`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Kern
  („verteilen statt kopieren"): `--print-mk` ([`DC-FA-CLI-010`](../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben))
  trägt ein `doc-commits`-Target (`--enable commits` + aus `ValidModules` abgeleitete
  Fokus-`--disable`-Liste, `--range`) — parallel zu `doc-immutable` (Modul `vcs`); Schwester-Repos
  beziehen die Commit-Traceability über das gepinnte Image, ohne Skript-Kopie. Der Name
  `doc-commits` (nicht `doc-trace`, das die RTM [`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
  belegt) hält die Ausgabe eindeutig. Eine Erweiterung des `--print-mk`-Fragments
  (7→8 Targets).
- **Restlücke unverändert.** Wie [ADR-0013](0013-pr-ci-und-traceability-gate.md)/[ADR-0016](0016-adr-immutable-gate.md):
  ohne Branch Protection (außerhalb des Repos) ist die CI nur *advisory*; der lokale Hook
  ist opt-in pro Klon (`make hooks`).

## Fitness Function

- `make trace-check RANGE=a..b` läuft **rot** (Exit 1) bei einem Commit ohne Kennung in der
  Range, **grün** sonst — wie zuvor, jetzt über das Modul. Ein Merge-/Revert-Commit in
  derselben Range bleibt grün.
- Der `commit-msg`-Hook (`--commit-msg -`) läuft **rot** bei einer Pending-Message ohne
  Kennung, **grün** sonst; `#`-Kommentar-only-Kennung zählt **nicht** (uniforme Bereinigung).
- Ohne aktives `commits` ist der Befundsatz byte-identisch (`make doc-check` aktiviert
  `commits` nicht) — opt-in-Selbsttest.
- `make arch-check` bleibt grün: `CommitMessages` wird **nur** im `git`-Adapter über go-git
  bedient; der Kern bleibt rein.
- fail-closed: fehlende/unauflösbare Range, unlesbares `.git`, nicht lesbare Message-Datei ⇒
  Exit 2 (kein stilles Grün) — Akzeptanztest.
- Selbsttest-Paritäts-Tests (die Klassen des abgelösten Skripts): ID erkannt; fehlende ID
  gefangen; Merge/Revert ausgenommen; Kennung nur auf `#`-Zeile ⇒ Befund (Bereinigung);
  leere gültige Range ⇒ Exit 0.
- `make semgrep`/`make lint` grün.

## Re-Evaluierungs-Trigger

- Die ID-Muster ändern sich (neues Bereichskürzel, anderes Slice-Schema) → `commits.id-patterns`
  in [`.d-check.yml`](../../../.d-check.yml) anpassen — jetzt **Config**, nicht Skript-Code.
- Bedarf, auch den **commit-msg-Hook** (Message-Modus) an Konsumenten zu verteilen (heute
  verteilt `doc-commits` nur den Range-Modus, wie `doc-immutable`) → ein hook-installierendes
  `--print-mk`-Target bzw. eine Hook-Vorlage (eigener CR).
- Eine dritte git-Message-bedürftige Befund-Klasse (z. B. Commit-Body-Konventionen) → nutzt
  denselben VCS-Port, eigene Anforderung.
- Wechsel des `core.commentChar` (nicht `#`) im Konsumenten → die Bereinigung müsste
  parametrisiert werden (heute `#` fest, wie der Skript-Vorläufer); Folge-CR bei echtem Bedarf.

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-07-01 | Entwurf (slice-056, Auftraggeber-Audit „welche `tools/*.sh` noch in d-check mechanisieren?" — letztes Familien-Skript): opt-in Modul `commits` mechanisiert `tools/trace-check.sh` über den [ADR-0024](0024-vcs-immutable-gate.md)-VCS-Port (`CommitMessages`), `make trace-check`/`commit-msg`-Hook/PR-CI dogfood-umgestellt, Skript per `git rm` entfernt + als `codepaths.ignore-refs`-Tombstone deklariert. **Supersedes** die Skript-Mechanik von [ADR-0013](0013-pr-ci-und-traceability-gate.md) (Policy/Bindepunkt/CI-Topologie unverändert). Modulname `commits` statt `trace` (Kollision mit `--trace`/RTM). Produkt-Code geändert → Release v0.35.0. Status Proposed. |
| 2026-07-01 | Angenommen mit der slice-056-Closure: Modul `commits` + Port-Erweiterung `CommitMessages` + go-git-Adapter implementiert/getestet, `make trace-check`/`commit-msg`-Hook dogfood-umgestellt (der Dogfood bewies sich an den slice-056-Commits selbst), `tools/trace-check.sh` per `git rm` entfernt + Tombstone; `--print-mk doc-commits` verteilt die Range-Prüfung ([`DC-FA-CLI-010`](../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben) 7→8). Drei unabhängige Reviews (R1 doc 1M/2L + R2 code 1M/1L/3I + R3 verifikation VERIFIED, Mutations-Beleg — alle Befunde behoben; keine Paritäts-Divergenz zum abgelösten Skript), `make ci`/`fullbuild` grün, Release v0.35.0. Status Accepted. |
