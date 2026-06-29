# ADR-0026 — Completeness-Gate dogfoodet den in-Produkt-Flag: `make completeness-check` ruft `--trace --require-complete` und löst die Skript-Mechanik von ADR-0017 ab

**Status:** Accepted
**Datum:** 2026-06-29
**Autor:** pt9912
**Bezug:** der in-Produkt-Flag
[`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
(`--trace --require-complete`) und die RTM-Quelle
[`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix);
**Supersedes die Skript-als-Gate-Quelle-Teilentscheidung** von
[ADR-0017](0017-requirements-completeness-gate.md) (**Policy und Bindepunkt bleiben** —
Teil-Supersede, wie [ADR-0002](0002-distribution-ghcr-image.md)/[ADR-0014](0014-latest-tag-fuer-stabile-releases.md));
Mechanik-Verwandtschaft [ADR-0013](0013-pr-ci-und-traceability-gate.md) (eine Quelle, zwei
Quadranten) und Präzedenz [ADR-0024](0024-vcs-immutable-gate.md) (Modul `vcs` löste das
Skript `adr-immutable-check.sh` ab); Verteilungs-Kern
[`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)
(„verteilen statt kopieren"); die immutable Skript-Referenz wird ein
[`codepaths.ignore-refs`](../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code)-Tombstone
([ADR-0025](0025-codepaths-ignore-refs.md)).
**Schärft:** keine neue Spec-Stelle — Prozess-/Mechanik-ADR; verbindlich für die
Verdrahtung des Closure-Gates `make completeness-check` auf den in-Produkt-Flag
[`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code).

## Kontext

[ADR-0017](0017-requirements-completeness-gate.md) führte `make completeness-check` als
**dünnen Wrapper** über `tools/completeness-check.sh` ein: das Skript ruft
`d-check --trace --json` und parst das top-level-Feld `orphans` in Bash (kein jq/python),
`orphans>0` → FAIL. Zu jenem Zeitpunkt gab es **keine** in-Produkt-Durchsetzung.

Später kam genau die: der Flag `--trace --require-complete`
([`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code))
bindet ≥1 Requirements-Waise an Exit 1 — **ausdrücklich**, „damit Konsumenten (a-check-Bootstrap)
die Vollständigkeits-Invariante als Makefile-Gate binden, **ohne die
`completeness-check.sh`-Parsing-Logik zu kopieren**" (Lastenheft-Historie). d-check **selbst**
aber bindet sein Closure-Gate weiter an das Skript und führt den Flag nur als
`make doc-complete` („Dogfood, **kein** Gate-Bindepunkt"). Damit **kopiert d-check für sich
die Logik, die es Konsumenten zu vermeiden empfiehlt** — der
[`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)-Copy-Drift
der Werkzeug-Familie, hier am eigenen Closure-Gate. Skript und Flag sind **zwei Quellen
derselben Wahrheit**.

## Entscheidung

`make completeness-check` ruft künftig **denselben in-Produkt-Flag** wie der Dogfood und die
Konsumenten: `d-check --trace --require-complete` im Image (`--network none`,
read-only-Mount). `tools/completeness-check.sh` wird **entfernt**. **Eine** Wahrheit (der
Flag), im Image verteilt — keine Skript-Kopie.

**Policy und Bindepunkt unverändert.** ≥1 Requirements-Waise → Exit 1, sonst grün; gebunden
an `make fullbuild`, **nicht** an `make gates`/`ci` (Greenfield duldet transiente Waisen, bis
der umsetzende Slice landet). Nur die **Mechanik** (Skript → in-Produkt-Flag) wechselt — wie
[ADR-0024](0024-vcs-immutable-gate.md) nur die Mechanik von
[ADR-0016](0016-adr-immutable-gate.md) wechselte, nicht dessen Policy.

**Waisen-Sichtbarkeit gewahrt.** Der Flag meldet die **Anzahl** (stderr) und Exit 1; die
**einzelnen** Waisen-IDs trägt die `--trace`-Tabelle als Status **`WAISE`**. Die konzise
Skript-Waisenliste entfällt zugunsten der kontext-reicheren Tabelle — **keine** Erweiterung
von [`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
nötig.

**Tombstone statt Edit.** Die immutable [ADR-0017](0017-requirements-completeness-gate.md)
referenziert `tools/completeness-check.sh` in Inline-Code; die Löschung bräche `make
doc-check` (`codepath-missing`) an unveränderlicher Doku. Ein Eintrag in `.d-check.yml`
`codepaths.ignore-refs` deklariert den Pfad als entfernt — der **zweite** reale
Anwendungsfall des Ventils aus [ADR-0025](0025-codepaths-ignore-refs.md).

**Kein neues `DC-*`, kein Release.** Die Mechanik
([`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)/[`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix))
existiert; es ändert sich **kein** Produkt-Code (`internal/`/`cmd/`) → das GHCR-Image bleibt
byte-identisch. Reiner Harness-Refactor: kein Versions-Bump, kein GHCR-Release.

## Verglichene Alternativen

| Alternative | Pro | Contra |
| --- | --- | --- |
| **`completeness-check` → in-Produkt-Flag (gewählt)** | **eine** Wahrheit, im Image verteilt (löst [`MR-007`](../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding)); kein Skript-Copy-Drift; Mechanik schon akzeptanz-getestet | Image-Lauf pro `completeness-check` (schwerer); Waisen-Ausgabe weniger konzis |
| Skript behalten (Status quo [ADR-0017](0017-requirements-completeness-gate.md)) | schon da, konzise Waisenliste | d-check kopiert die Logik, die es Konsumenten zu vermeiden empfiehlt — Dogfood-Bruch bleibt |
| Skript **und** Flag koexistieren (heute: `completeness-check` vs. `doc-complete`) | beide Bindepunkte vorhanden | zwei Quellen derselben Wahrheit, Drift-Risiko, unklar welches „die" Closure-Wahrheit ist |
| Flag erst um konzise Waisenliste erweitern, dann umstellen | volle Ausgabe-Parität | unnötig — die `WAISE`-Zeilen + Anzahl benennen die Waisen schon; eigener CR bei echtem Bedarf |

## Konsequenzen

- **`make completeness-check` braucht jetzt das Image** (wie `make adr-check` seit
  [ADR-0024](0024-vcs-immutable-gate.md)) — schwerer pro Lauf; dafür eine verteilte,
  kopie-freie Wahrheit. Die CI baut das Image ohnehin.
- **`tools/completeness-check.sh` entfernt;** sein Negativ-Selbsttest (orphans-Parsing-
  Vektoren, fail-closed) lebt als in-Produkt-Akzeptanztest der CLI
  ([`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code))
  weiter.
- **Kein Release:** kein Produkt-Code geändert → Image byte-identisch zu v0.34.0.
- **Doc-check bleibt grün** trotz entferntem Skript: die immutable
  [ADR-0017](0017-requirements-completeness-gate.md)-Referenz ist als
  `codepaths.ignore-refs`-Tombstone deklariert ([ADR-0025](0025-codepaths-ignore-refs.md)).
- **Restlücke wie [ADR-0017](0017-requirements-completeness-gate.md):** `completeness-check`
  bleibt ein Closure-/`fullbuild`-Bindepunkt, kein `gates`/`ci`-Gate — Greenfield-bewusst.

## Fitness Function

- `make completeness-check` läuft **rot** (Exit 1) bei ≥1 Requirements-Waise (adversariale
  Fixture: Anzahl-Meldung + sichtbare `WAISE`-Zeile), **grün** sonst — wie zuvor, jetzt über
  den Flag.
- `make gate-consistency` grün: das Target existiert und ist in
  [`AGENTS.md`](../../../AGENTS.md) §4 **und** [`harness/README.md`](../../../harness/README.md)
  §Sensors gelistet (beide Richtungen).
- `make doc-check` grün trotz entferntem Skript (Tombstone-Ventil greift).
- Kein Produkt-Code geändert (`git diff internal/ cmd/` leer) → Image byte-identisch.

## Re-Evaluierungs-Trigger

- Bedarf einer **knappen Nur-Waisen-Ausgabe** statt der vollen `--trace`-Tabelle →
  Erweiterung von
  [`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code),
  eigene Anforderung.
- Die [ADR-0017](0017-requirements-completeness-gate.md)-**Policy** ändert sich (anderer
  Bindepunkt, andere Waisen-Definition) → dort entscheiden, nicht hier.
- **Leerer/kaputter Scan-Scope** (findet `--trace` keine Anforderung, ist `Orphans==0` ⇒
  Exit 0 — ein silent-green ohne bewiesene Vollständigkeit) wird zum realen Risiko, etwa in
  Konsumenten-Repos ohne stabilen Lastenheft-Pfad → ein Vollständigkeits-Boden (`total>0`) als
  [`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)-Erweiterung.
  **Vorbestehend** — das abgelöste Skript bestand `{"total":0,"orphans":0}` ebenso grün; diese
  ADR führt das Verhalten nicht ein und ändert es nicht.

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-06-29 | Entwurf (slice-055, Auftraggeber-Audit „welche `tools/*.sh` noch in d-check mechanisieren?"): `make completeness-check` dogfoodet `--trace --require-complete` ([`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)), `tools/completeness-check.sh` entfernt und als `codepaths.ignore-refs`-Tombstone deklariert. **Supersedes** die Skript-als-Gate-Quelle-Teilentscheidung von [ADR-0017](0017-requirements-completeness-gate.md) (Policy/Bindepunkt unverändert). Kein Produkt-Code, kein Release. Status Proposed. |
| 2026-06-29 | Angenommen mit der slice-055-Closure: `make completeness-check` ruft `$(DCHECK_RUN) $(COMPLETE_FLAGS)` (geteilte Variable, kein Divergenz-Risiko zu `doc-complete`), `tools/completeness-check.sh` per `git rm` entfernt + als `codepaths.ignore-refs`-Tombstone deklariert; AGENTS §4 + harness/README §Sensors nachgezogen. Adversariale Waisen-Probe trieb das Gate rot (`WAISE`-Zeile + Anzahl, danach revertiert); zwei unabhängige Reviews (R1 doc 0H/1M/1L + R2 mechanik 0H/0M/1L/1I/1 REFUTED, alle Befunde behoben — fail-closed-Regression mit Skript-Zitat widerlegt). Kein Produkt-Code → Image byte-identisch, kein Release. Status Accepted. |
