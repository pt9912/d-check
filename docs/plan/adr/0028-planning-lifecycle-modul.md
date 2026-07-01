# ADR-0028 — Planning-Lifecycle-Konsistenz: das hermetische Modul `planning` mechanisiert `planning-consistency.sh`

**Status:** Accepted
**Datum:** 2026-07-01
**Autor:** pt9912
**Bezug:** [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(opt-in Modul `planning`); **Familie** [ADR-0024](0024-vcs-immutable-gate.md)/[ADR-0026](0026-completeness-in-product-gate.md)/[ADR-0027](0027-commits-traceability-modul.md)
(dieselbe „Skript → verteilbares d-check-Feature"-Linie des `tools/*.sh`-Audits); der
Vorläufer `tools/planning-consistency.sh` entstand **ohne** begleitende ADR — diese
ADR ist die **erste** des Planning-Gates, sie supersedet daher keine ADR-Mechanik,
sondern verkörpert die bislang skript-/ADR-lose Durchsetzung; Port-/Ordner-Konvention
[ADR-0005](0005-modul-layout-hexagon-ordner.md); Verteilungs-Kern
[`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)
(„verteilen statt kopieren"); die durchgesetzte Policy lebt in
[`AGENTS.md` §3.3](../../../AGENTS.md#33-git-mv--inhaltsänderung--zwei-commits) (Slice-Lifecycle)
und [`MR-013`](../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise);
die immutable Skript-Referenz wird ein
[`codepaths.ignore-refs`](../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code)-Tombstone
([ADR-0025](0025-codepaths-ignore-refs.md)).
**Schärft:** die neue Algorithmus-Sektion
[`spec/spezifikation.md` §DC-FA-PLAN-001.a](../../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
(Invariante, Heading-Guard, Grund-Code `planning-drift`).

## Kontext

Der `tools/*.sh`-Audit (Auftraggeber) hat drei Gate-Skripte als d-check-Feature
mechanisiert: `adr-immutable-check.sh` → Modul `vcs` ([ADR-0024](0024-vcs-immutable-gate.md)),
`completeness-check.sh` → in-Produkt-Flag ([ADR-0026](0026-completeness-in-product-gate.md)),
`trace-check.sh` → Modul `commits` ([ADR-0027](0027-commits-traceability-modul.md)). Der
**letzte** Kandidat ist `tools/planning-consistency.sh` (das `make planning-check`-Meta-Gate):
ein kopiertes Shell-Skript, das die **Planning-Lifecycle-Invariante** erzwingt — die
Roadmap darf im `## Aktuelle Welle`-Abschnitt „Keine aktive Welle" genau dann behaupten,
wenn kein `slice-*` in `in-progress/` liegt (und umgekehrt; fail-closed, Heading-Guard).
Als kopiertes Skript trägt es denselben Copy-Drift
([`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Klasse).

Anders als die drei Vorgänger ist diese Prüfung **hermetisch**: sie liest **nur** den
Arbeitsbaum (eine Roadmap-Datei + das Verzeichnis-Listing) — **kein** git, **kein** Netz.
Damit ist sie **d-check-nativ** (Existenz-/Datei-Prüfung wie `codepaths`), ohne den
Eingabe-Scope zu erweitern (kein VCS-Port wie `vcs`/`commits`). Der Verteilungswert ist
aber **bewusst nachrangig**: die „Keine aktive Welle"-Konvention ist harness-layout-
spezifisch (nur Adopter-Repos mit demselben Roadmap-/Slice-Layout profitieren).

## Entscheidung

Ein opt-in Modul `planning` (Default aus, 15. Modul) mechanisiert die Invariante als
**hermetischen Post-Pass**: es liest `planning.roadmap`, bestimmt `hasSlices` (≥1
`planning.slice-glob`-Datei im Roadmap-Verzeichnis) und `hasActive` (der
`planning.heading`-Block trägt den `planning.marker` **nicht**) und meldet bei
`hasActive ≠ hasSlices` — oder bei fehlender kanonischer Überschrift/Roadmap-Datei
(**Heading-Guard, fail-closed**) — den Grund-Code `planning-drift`. **Diagnose-only.**

**Hermetisch, keine Port-Erweiterung.** `planning` nutzt **nur** den vorhandenen
Filesystem-Port (`ReadFile` + Verzeichnis-`List`) — **kein** git, **kein** VCS-Port,
**keine** neuen CLI-Flags. Es relativiert **weder**
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus) **noch**
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit):
gleicher Arbeitsbaum ⇒ gleicher Befund, read-only, netzlos. Das ist der **saubere**
Unterschied zu `vcs`/`commits` (die `.git` + Range brauchten) — `planning` ist ein
normales hermetisches Modul.

**Konfigurierbar, damit verteilbar.** `planning.roadmap` (Pflicht zum Aktivieren;
leer ⇒ inert), `planning.heading`/`marker`/`slice-glob` mit Konventions-Defaults
(`## Aktuelle Welle` / „Keine aktive Welle" / `slice-*.md`). So ist die harness-
spezifische Invariante ein **parametrierbares** Modul statt eines hartcodierten Skripts —
Adopter mit abweichendem Layout überschreiben die Keys.

**Anderer Bindepunkt als `vcs`/`commits`.** `make planning-check` bleibt ein
**Meta-Gate in `make gates`** (Arbeitsbaum-Bindepunkt, Inner-Loop) — anders als
`adr-check`/`trace-check` (Commit-/Range-Bindepunkt, **nicht** in `gates`). Nur die
**Mechanik** wechselt (Skript → Modul im Image); Policy und Bindepunkt bleiben.

**Dogfood + Skript entfernt.** `make planning-check` ruft künftig **das d-check-Image**
mit `--enable planning` (`FOCUS_DISABLE`, nur `planning`) statt
`bash tools/planning-consistency.sh`. Der `planning:`-Block liegt in
[`.d-check.yml`](../../../.d-check.yml), **außerhalb** der Default-`modules`-Liste →
`make doc-check` aktiviert `planning` nicht (default-aus byte-identisch). d-check prüft
so seine **eigene** Roadmap-Konsistenz. `tools/planning-consistency.sh` wird **atomar**
per `git rm` entfernt (wie bei [ADR-0027](0027-commits-traceability-modul.md)); die
immutable Inline-Referenz des abgelösten Skripts wird ein `codepaths.ignore-refs`-Tombstone
([ADR-0025](0025-codepaths-ignore-refs.md)). Der 3-Richtungs-Negativ-Selbsttest des
Skripts wandert als Akzeptanztest ins Modul (`make test`).

**`--print-mk doc-planning`.** `--print-mk` trägt ein `doc-planning`-Target (`--enable
planning`, hermetisch **ohne** `--range`; [`DC-FA-CLI-010`](../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
8→9). Verteilt die Invariante an Adopter mit demselben Layout — bewusst mit
begrenztem Wert (s. Kontext).

## Verglichene Alternativen

| Alternative | Pro | Contra |
| --- | --- | --- |
| **Hermetisches Modul `planning` (gewählt)** | d-check-nativ (kein git/Port); verteilt im Image (löst Copy-Drift); parametrierbar für Adopter; [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)/[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) unberührt; Familie komplett | begrenzter Verteilungswert (layout-spezifisch); Image-Lauf pro `planning-check` (schwerer als Skript) |
| Skript belassen, nicht mechanisieren | schon da, leichtgewichtig | Copy-Drift bleibt; letztes Familien-Skript unmechanisiert (Auftrag offen) |
| Als Strategie in `codepaths` falten | ein Modul weniger | mischt Inline-Code-Pfad-Existenz mit Roadmap-Invariante; verwischt Grund-Code/Zweck |
| In `make gates` als in-Produkt-Flag (wie completeness) | kein neues Modul | die Invariante ist kein CLI-Modus-Kandidat (kein RTM-artiges Feature); ein Modul ist der natürliche Ort |

## Konsequenzen

- **`make planning-check` braucht jetzt das Image** (wie `adr-check`/`trace-check`/`completeness-check`
  seit ihrer Mechanisierung) — schwerer pro Lauf; dafür verteilte, kopie-freie Wahrheit. `make gates`
  baut das Image ohnehin (`doc-check`).
- **`tools/planning-consistency.sh` entfernt;** sein 3-Richtungs-Negativ-Selbsttest (Slice+idle,
  kein-Slice+aktiv, kaputte Überschrift) lebt als Modul-Akzeptanztest (`make test`).
- **Kein neuer Eingabe-Scope.** `planning` bleibt hermetisch — die [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Modulliste des
  `gate-consistency`-Selbsttests ist unberührt (opt-in, nicht in `modules:`).
- **Doc-check bleibt grün** trotz entferntem Skript: die immutable Inline-Referenz des abgelösten Skripts ist als
  `codepaths.ignore-refs`-Tombstone deklariert ([ADR-0025](0025-codepaths-ignore-refs.md)).
- **Release nötig.** Produkt-Code (`internal/`) ändert sich (neues Modul) → Minor-Bump v0.36.0,
  GHCR-Release (wie [ADR-0024](0024-vcs-immutable-gate.md)/[ADR-0027](0027-commits-traceability-modul.md)).
- **`tools/*.sh`-Audit abgeschlossen.** `planning` ist das letzte mechanisierte Familien-Skript;
  die verbleibenden `tools/*.sh` (`arch-check`, `coverage-gate`, `semgrep`, `gate-consistency`,
  `image-test`, `bench-fixture`, `harness/*`) sind bewusst **nicht** d-check-fähig (Go-Build/SAST/
  Meta-Spannung Markdown↔Makefile↔YAML) — Roadmap §Nächste Wellen.

## Fitness Function

- `make planning-check` läuft **rot** (Exit 1) bei Drift (Slice+idle-Marker; kein-Slice+aktiv;
  kaputte/fehlende `## Aktuelle Welle`-Überschrift), **grün** bei Konsistenz — wie zuvor, jetzt
  über das Modul.
- Ohne aktives `planning` ist der Befundsatz byte-identisch (`make doc-check` aktiviert `planning`
  nicht) — opt-in-Selbsttest.
- Hermetik: der `planning`-Lauf berührt **kein** `.git`, **kein** Netz (`--network none`), schreibt
  nicht — [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  gehalten; gleicher Baum ⇒ gleicher Befund ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- Selbsttest-Paritäts-Tests (die drei Skript-Richtungen + die zwei konsistenten Gegenproben +
  Heading-Guard) feuern/schweigen korrekt (Modul-Akzeptanztest via Mem-FS).
- `make semgrep`/`make lint`/`make arch-check` grün (kein neuer Import, kein git).

## Re-Evaluierungs-Trigger

- Die Roadmap-Struktur ändert sich (andere Überschrift, anderer Marker, anderes Slice-Layout) →
  `planning.heading`/`marker`/`slice-glob` in [`.d-check.yml`](../../../.d-check.yml) anpassen —
  jetzt **Config**, nicht Skript-Code.
- Bedarf, mehrere Roadmaps/Slice-Verzeichnisse zu prüfen → Erweiterung von
  [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  (heute genau eine, wie das Skript).
- Eine git-basierte Lifecycle-Prüfung (z. B. Move-Commit-Bündelung, [`MR-013`](../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise))
  → das ist die Domäne des VCS-Ports (`vcs`/`commits`), nicht des hermetischen `planning`.

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-07-01 | Entwurf (slice-057, Auftraggeber-`tools/*.sh`-Audit — **letzter** Kandidat): opt-in Modul `planning` mechanisiert `tools/planning-consistency.sh` als **hermetischen** Post-Pass (Filesystem-Port, **kein** git/VCS-Port), `make planning-check` dogfood-umgestellt, Skript per `git rm` entfernt + `codepaths.ignore-refs`-Tombstone. Erste Planning-Gate-ADR (slice-040-Skript war ADR-los). `--print-mk doc-planning` ([`DC-FA-CLI-010`](../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben) 8→9). Bewusst nachrangig (layout-spezifischer Verteilungswert). Produkt-Code geändert → Release v0.36.0. Status Proposed. |
| 2026-07-01 | Angenommen mit der slice-057-Closure: Modul `planning` (`CheckPlanning` über den Filesystem-Port) + Config (`planning.roadmap`/`heading`/`marker`/`slice-glob`, config-zeitig validiert) implementiert/getestet, `make planning-check` dogfood-umgestellt (bleibt in `make gates`), `tools/planning-consistency.sh` per `git rm` entfernt + Tombstone; `--print-mk doc-planning` verteilt die hermetische Prüfung. **Drei** unabhängige Reviews (R1 doc 1M/1L + R2 code 1M/2L/2I + R3 verifikation NACHBESSERN→behoben, Mutations-Belege für beide Fail-closed-Guards: slice-glob-Validierung + Duplikat-Überschrift); keine Paritäts-Divergenz zum abgelösten Skript. `make ci`/`fullbuild` grün, Release v0.36.0. **Schließt den `tools/*.sh`-Audit ab** (viertes/letztes Familien-Skript). Status Accepted. |
