# ADR-0031 — Deklarations-Konsistenz Doku ↔ Build-Targets: das hermetische Modul `targets` mechanisiert den `gate-consistency.sh`-Kern

**Status:** Accepted
**Datum:** 2026-07-05
**Autor:** pt9912
**Bezug:** [`DC-FA-TGT-001`](../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
(opt-in Modul `targets`); Muster-Vorbild und **Prämissen-Wende**
[ADR-0028](0028-planning-lifecycle-modul.md) (hermetischer Post-Pass über den
Filesystem-Port); Verteilungs-Kern
[`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)
(„verteilen statt kopieren"); Port-/Ordner-Konvention
[ADR-0005](0005-modul-layout-hexagon-ordner.md).
**Schärft:** die neue Algorithmus-Sektion
[`spec/spezifikation.md` §DC-FA-TGT-001.a](../../../spec/spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets)
(zwei Richtungen `gate-phantom`/`gate-undocumented`, `^|`-Tabellen-Scoping,
Zeilen-Heuristik, fail-closed).

## Kontext

`tools/gate-consistency.sh` ist ein **kopiertes** Meta-Gate über die Repo-Familie
(d-check, a-check, belief-agent, …) und **driftet** bereits real auseinander:
a-checks Fassung trägt einen Pin-Block, eine Utility-Allowlist und `d-check.mk`
als zweite Target-Quelle, d-checks nicht (Cross-Repo-Drift-Beleg, geprüft
2026-07-05 in a-checks `tools/gate-consistency.sh`). Sein
prüfwürdiger, cross-repo **gemeinsamer Kern** ist die Doku-↔-Makefile-
Deklarations-Konsistenz — „jedes in der Doku als ` `make X` ` behauptete Gate
existiert real, und jedes reale Gate ist dokumentiert" — dieselbe
[`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Copy-Drift-Klasse
wie die fünf bereits mechanisierten Skripte (`adr-immutable→vcs`,
`completeness→Flag`, `trace→commits`, `planning→planning`, `arch→a-check`).

**Die Prämissen-Wende.** [ADR-0028](0028-planning-lifecycle-modul.md) hatte in
seinen Konsequenzen `gate-consistency` **explizit** als „bewusst **nicht**
d-check-fähig" eingestuft — Begründung: „Meta-Spannung Markdown ↔ Makefile ↔
YAML". Genau diese Einstufung ist **überholt**, und zwar durch ADR-0028 selbst:
das dort eingeführte Modul `planning` etablierte den **hermetischen Post-Pass**,
der nicht den Markdown-Baum scannt, sondern über den Filesystem-Port eine
**deklarierte Datei + Repo-Struktur** liest (Roadmap ↔ `in-progress/`-Listing),
um eine **Doku-Behauptung gegen die Repo-Realität** zu prüfen. Ein Makefile über
denselben `ReadFile`-Port zu lesen und seine Regelnamen gegen die Doku-Tabellen
zu prüfen ist **dieselbe Klasse** — kein Makefile-*Ausführen*, kein neuer Port,
kein erweiterter Netz-/git-Scope. Die „nicht d-check-fähig"-Einstufung war an
d-checks **damaligen** Fähigkeiten (Markdown-only) gemessen; `planning` hat die
Prämisse verschoben.

**Identität.** Damit vollzieht d-check bewusst die Ausweitung vom
**Markdown-Doku-Referenz-Checker** zum **Repo-Deklarations-Konsistenz-Checker**
(Doku ↔ Makefile ↔ Config ↔ git ↔ Verzeichnis). Diese Ausweitung ist nicht neu —
`vcs`/`commits`/`tracked` (git) und `planning` (Verzeichnis-Struktur) haben die
Markdown-only-Grenze längst überschritten; `targets` liest zusätzlich das
Makefile. Der Name **d-check** liest sich fortan als *declaration-checker*
(Auftraggeber-Framing 2026-07-05).

## Entscheidung

Ein opt-in Modul `targets` (Default aus, 17. Modul) mechanisiert den
Doku-↔-Makefile-Kern als **hermetischen Post-Pass** (Algorithmus s. **Schärft**
oben):

- **Zwei Richtungen.** Ein ` `make X` ` in einer Doku-Tabellenzeile ohne
  Makefile-Regel ⇒ `gate-phantom`; eine Makefile-Regel (minus
  `targets.exempt-targets`) ohne Eintrag in der Autoritäts-Doku ⇒
  `gate-undocumented`.
- **Tabellen-Scoping (`^|`).** ` `make X` ` zählt **nur** aus Tabellenzeilen —
  Prosa (z. B. „Richtig: ` `make gates` `") ist keine Existenz-Behauptung. Das
  ist die Parität zum abgelösten Skript (`grep -E '^\|'`) und verhindert
  spuriöse `gate-phantom`.
- **Hermetisch, keine Port-Erweiterung.** Nur der vorhandene Filesystem-Port
  (`ReadFile` auf die konfigurierten Dateien) — **kein** git, **kein Ausführen**
  des Makefile, **kein** Netz, **keine** neuen CLI-Flags. Relativiert **weder**
  [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus) **noch**
  [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).
  Die Makefile-Regel-Erkennung ist eine **statische Zeilen-Heuristik** (literale
  Namen; keine Pattern-Rules/variablen/`include`-Targets) — in Parität zum
  Skript, mit derselben dokumentierten Grenze.
- **Konfigurierbar, damit verteilbar.** `targets.makefiles` (Pflicht zum
  Aktivieren; leer ⇒ inert), `targets.doc-tables` (leer ⇒ Richtung 1 entfällt),
  `targets.authority` (leer ⇒ Richtung 2 entfällt — die Richtungen sind
  unabhängig), `targets.exempt-targets` (Utility-Regeln ohne Doku-Pflicht). So ist
  die harness-spezifische Invariante ein **parametrierbares** Modul statt eines
  hartcodierten Skripts — Konsumenten (a-check/belief-agent) überschreiben die
  Keys (z. B. `makefiles: [Makefile, d-check.mk]`).

**Scope = Kern.** `targets` mechanisiert **nur** die zwei cross-repo-driftenden
Richtungen. Der `.d-check.yml`-Modul-Listen-Selbstcheck (heutige
gate-consistency-Prüfung 3, Netzlos-Gate-Integrität) bleibt **draußen**:
repo-spezifische Selbst-Konsistenz der eigenen Config, **kein** Kopie-Drift-Kern.
Er verbleibt im **Rest** von `gate-consistency.sh` — anders als bei `planning`
wird das Skript also **nicht vollständig** entfernt, nur sein Kern.

**Dogfood + Selbstbezug.** `make gate-consistency` ruft künftig das d-check-Image
mit `--enable targets` für die zwei Richtungen (behält den kleinen Prüfung-3-Rest)
statt der zwei Kern-Prüfungen im Skript. Der `targets:`-Block liegt in
[`.d-check.yml`](../../../.d-check.yml) **außerhalb** der Default-`modules`-Liste
→ `make doc-check` aktiviert `targets` nicht (default-aus byte-identisch). d-check
prüft so seine **eigene** Doku-↔-Makefile-Konsistenz über sein eigenes Binary —
Bootstrap-Zirkularität wie beim doc-check-Dogfooding, abgesichert durch einen
modul-internen **Negativ-Selbsttest** (Guard entfernt ⇒ Test rot, wie
`planning`/`commits`/`vcs`).

**`--print-mk doc-targets`.** `--print-mk` trägt ein `doc-targets`-Target
(`--enable targets`, hermetisch **ohne** `--range`;
[`DC-FA-CLI-010`](../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
10→11). Verteilt die Prüfung an Konsumenten ohne Skript-Kopie.

## Verglichene Alternativen

| Alternative | Pro | Contra |
| --- | --- | --- |
| **Hermetisches Modul `targets` (gewählt)** | d-check-nativ (kein git/Port, wie `planning`); verteilt im Image (löst Copy-Drift); parametrierbar; [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)/[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) unberührt; jedes Repo fährt d-check ohnehin | Identitäts-Ausweitung (Makefile-Lesen); Bootstrap-Selbstbezug; Image-Lauf pro `gate-consistency` |
| Neues Schwester-Tool (analog a-check) | scharfe Domänen-Trennung | zahlt d-checks Bootstrap erneut (Repo/Image/Pipeline/Config) für **einen** Check; jedes Repo fährt d-check schon — unverhältnismäßig ([`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Ökonomie) |
| Skript belassen, nicht mechanisieren | schon da, leichtgewichtig | Cross-Repo-Drift bleibt (a-check ≠ d-check ist real) |
| Content-Pin des Skripts (Modul `immutable`) | vorhandene Mechanik | falsch für **parametriert-divergente** Skripte (a-check divergiert legitim) — Content-Pin fängt nur *identisch-geteilte* Dateien |

## Konsequenzen

- **Identitäts-Ausweitung wird kanonisch:** d-check liest ab jetzt auch das
  Makefile (Deklarations-Konsistenz-Checker). Der Re-Evaluierungs-Grund von
  ADR-0028 („`gate-consistency` nicht d-check-fähig") ist **eingetreten und
  aufgelöst** — nicht per Supersede (ADR-0028s `planning`-Entscheidung steht),
  sondern weil dessen eigenes Muster die Einstufung überholt.
- **`make gate-consistency` braucht jetzt das Image** für die zwei Kern-
  Richtungen (behält den kleinen Prüfung-3-Rest im Skript) — schwerer pro Lauf,
  dafür kopie-freie, verteilte Wahrheit. `make gates` baut das Image ohnehin.
- **`tools/gate-consistency.sh` bleibt (reduziert)** — anders als bei
  `planning`/`commits`/`vcs` (voll entfernt): nur der Doku-↔-Makefile-Kern wandert
  ins Modul, die repo-lokale Prüfung-3 bleibt. **Kein** `codepaths.ignore-refs`-
  Tombstone (das Skript existiert weiter).
- **Bootstrap-Selbstbezug:** d-check prüft seine eigene Doku-↔-Makefile-
  Konsistenz — abgesichert durch den modul-internen Negativ-Selbsttest.
- **Paritäts-Pflicht:** der Modul-Befundsatz muss dem Kern des abgelösten
  Skripts auf dem aktuellen Baum entsprechen (Mutations-Proben je Richtung, wie
  `arch→a-check`) — ein aktives Gate wird nur mit belegter Parität stillgelegt.
- **a-checks Pin-Block bleibt a-check-lokal** (nur a-check pinnt per Digest;
  kein cross-repo-Drift) — nicht Teil dieses Moduls.
- **Release nötig.** Produkt-Code (`internal/`) ändert sich (neues Modul) →
  Minor-Bump **v0.38.0**, GHCR-Release (wie [ADR-0028](0028-planning-lifecycle-modul.md)).

## Fitness Function

- `make gate-consistency` läuft **rot** (Exit 1) bei einem dokumentierten
  Phantom-Target (`gate-phantom`) oder einem undokumentierten Makefile-Target
  (`gate-undocumented`), **grün** bei Konsistenz — wie zuvor, jetzt über das
  Modul (Kern).
- Ohne aktives `targets` ist der Befundsatz byte-identisch (`make doc-check`
  aktiviert `targets` nicht) — opt-in-Selbsttest.
- Prosa-`make X` löst **kein** `gate-phantom` aus (Tabellen-Scoping-Test).
- Hermetik: der Lauf berührt **kein** `.git`, **kein** Netz (`--network none`),
  führt das Makefile **nicht** aus, schreibt nicht —
  [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  gehalten; gleicher Baum ⇒ gleicher Befund
  ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- Paritäts-Test (Kern ≡ Skript) + Negativ-Selbsttest des Bootstrap-Guards
  (Guard entfernt ⇒ rot) feuern/schweigen korrekt.
- `make semgrep`/`make lint`/`make arch-check` grün (kein neuer Import, kein git).

## Re-Evaluierungs-Trigger

- Bedarf, die Autoritäts-Vollständigkeit auf einen **Abschnitt** zu scopen
  (statt whole-file) → optionaler `targets.authority`-Section-Anker (heute
  whole-file, in Parität zum Skript).
- Bedarf, die Prüfung-3 (Modul-Listen-Selbstkonsistenz) ebenfalls zu
  mechanisieren → eigene Anforderung (heute bewusst im Skript-Rest).
- Pattern-Rules/variabel benannte Targets müssen erkannt werden → das
  überschritte die statische Zeilen-Heuristik (Makefile-Ausführen ⇒
  [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Bruch,
  eigene Anforderung/ADR).
- Nicht-Make-Build-Systeme → eigene Anforderung.

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-07-05 | Entwurf (slice-063, Auftraggeber „driftende Skripte einfangen" + „d-check = declaration-checker, Klasse A regelkonform"): opt-in Modul `targets` mechanisiert den Doku-↔-Makefile-Kern von `tools/gate-consistency.sh` als **hermetischen** Post-Pass (Filesystem-Port, **kein** git/Netz/Makefile-Ausführen). Zwei Richtungen `gate-phantom`/`gate-undocumented`, `^|`-Tabellen-Scoping, statische Zeilen-Heuristik (Parität zum Skript). Löst die von [ADR-0028](0028-planning-lifecycle-modul.md) notierte „nicht d-check-fähig"-Einstufung auf (dessen `planning`-Muster verschob die Prämisse) — Identitäts-Ausweitung zum Deklarations-Konsistenz-Checker. Skript **reduziert** (Prüfung-3-Rest bleibt), kein Tombstone. `--print-mk doc-targets` ([`DC-FA-CLI-010`](../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben) 10→11). Produkt-Code geändert → Release v0.38.0. Status Proposed. |
| 2026-07-05 | **Angenommen** mit der Closure des umsetzenden Slice (Modul `targets` implementiert, reviewt und als **v0.38.0** released; `make gates`/`ci`/`fullbuild` grün). Status **Accepted** (Closure-Backfill, R1-F-4 zu welle-53). Der Prüfung-3-Rest-Verbleib im Skript wird mit welle-53/[ADR-0032](0032-gate-consistency-tombstone.md) revidiert (Voll-Tombstone). |
