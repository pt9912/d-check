# Verifikations-Bericht — slice-288 (DoD-/Plan-Konformität)

**Rolle:** Verifier (gegen DoD §2, Plan §1/§3 samt Plan-Nachtrag/§5/§6, Anforderung, Verfahren und ADR; nicht Reviewer, nicht Validator)
**Gegenstand:** Stand `b3b197b3` (Commits `529d357e`, `a03977a3`, `b3b197b3`; Basis `a03977a3~2`; 6 Dateien)
**Anforderung:** [`DC-FA-PSYNC-001`](../../spec/lastenheft.md#dc-fa-psync-001--gleichheit-eines-wertes-an-mehreren-stellen-modul-pinsync-opt-in)
**Verfahren:** [`DC-FA-PSYNC-001.a`](../../spec/spezifikation.md#dc-fa-psync-001a--gleichheit-eines-wertes-an-mehreren-stellen-pinsync)
**ADR:** [ADR-0113](../plan/adr/0113-pinsync-modul.md)
**Vorgänger:** [`2026-10-31-slice-288-code-r1.md`](2026-10-31-slice-288-code-r1.md)
**Datum:** 2026-10-31

## Selbst gefahrene Sensoren (echte Ausgabe)

- `make doc-check` am Bestand → Exit 0:
  `d-check: pinsync: 2 Pin(s) geprüft, 4 Quelle(n) gelesen` und `d-check: 1241 Datei(en) geprüft, 0 Befund(e)`.
- `make gates` → Exit 0, Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; darin die pinsync-Zeile, `1241 Datei(en) geprüft, 0 Befund(e)` und semgrep `Ran 55 rules on 74 files: 0 findings`.
- `make image-test` → `image-test: OK — Akzeptanzkriterien erfüllt (…)`, Exit 0 (Teile 1 bis 7; **kein pinsync-Teil**, siehe V-2).
- Selbstlauf `docker run --rm --network none --user 4242:4242 -v "$PWD":/repo:ro d-check:latest` → Exit 0, `d-check: pinsync: 2 Pin(s) geprüft, 4 Quelle(n) gelesen`, `1241 Datei(en) geprüft, 0 Befund(e)`.
- `make blackbox-probe REF=v0.90.0` → **Exit 2, 4 von 20 Vergleichen weichen ab** — alle vier sind `repo/*`, siehe V-1. Die 16 Vergleiche der Fixtures `ids`, `links`, `sauber`, `targets` (je `-`, `--json`, `--yaml`, `--doctor`) sind gleich.
- Bestandszählung selbst wiederholt: die vier `grep -cE` aus dem Plan ⇒ 1 / 1 / 1 / 1; `grep -c $'\r'` ⇒ 0 in `Dockerfile` und `Makefile`.

## Brechproben (beide Richtungen)

**Durchführung:** Wegen der Dockerfile-/Makefile-Pins ist ein `make doc-check` auf mutiertem Stand nicht tauglich (es baute das Image mit einem erfundenen Go-Wert neu). Die Proben liefen deshalb mit dem **bereits gebauten** `d-check:latest` auf einer Kopie des Arbeitsbaums im Scratchpad (`docker run --rm --network none -v <Kopie>:/repo:ro d-check:latest`, dieselbe Aufrufform wie `DCHECK_RUN`), je Probe aus dem Original zurückgesetzt. Das Repo selbst blieb unberührt (`git diff` leer).

| # | Mutation (nur eine Datei) | Beobachtet | Urteil |
|---|---|---|---|
| 1 | `Makefile` `GO_VERSION` 1.27.3 → 1.27.4 | `Makefile:13  GO_VERSION  pinsync-mismatch  Pin "GO_VERSION": Makefile hat "1.27.4", der Bezug Dockerfile hat "1.27.3" (Werte: Dockerfile=1.27.3, Makefile=1.27.4)`; `1 Befund(e)`, Exit 1 | rot, richtiger Grund, Fundstelle |
| 2 | `Dockerfile` `ARG GO_VERSION` → 1.27.4 | `Makefile:13  GO_VERSION  pinsync-mismatch … (Werte: Dockerfile=1.27.4, Makefile=1.27.3)`, Exit 1 | rot |
| 3 | `Makefile` `GOLANGCI_LINT_VERSION` → `v2.16.0` | `Makefile:14  GOLANGCI_LINT_VERSION  pinsync-mismatch … Makefile hat "2.16.0", der Bezug Dockerfile hat "2.15.0"`, Exit 1 | rot |
| 4 | `Dockerfile` `ARG GOLANGCI_LINT_VERSION` → `v2.16.0` | `Makefile:14  GOLANGCI_LINT_VERSION  pinsync-mismatch … (Werte: Dockerfile=2.16.0, Makefile=2.15.0)`, Exit 1 | rot |
| 5 | `Dockerfile` ohne `v` (`2.15.0`) | `0 Befund(e)`, Exit 0 | grün, `v` gilt gleich |
| 5b | `Makefile` ohne `v` (`2.15.0`) | `0 Befund(e)`, Exit 0 | grün |
| 6 | `Dockerfile`: `ARG GO_VERSION=` auskommentiert (0 Treffer) | `Dockerfile:1  GO_VERSION  pinsync-no-match  Pin "GO_VERSION": der Ausdruck "^ARG GO_VERSION=(\S+)$" trifft in Dockerfile keine Zeile mit nicht leerem Wert`, Exit 1 | rot |
| 7 | `Dockerfile`: `ARG GO_VERSION=1.27.3` dupliziert (2 Treffer) | `Dockerfile:31  GO_VERSION  pinsync-ambiguous  … trifft in Dockerfile mehrere Zeilen (31, 32)`, Exit 1 | rot |
| 8 | Rücknahme aller Proben | `diff` gegen Original leer; Bestandslauf grün (siehe oben) | grün |

Die Fundstelle liegt bei Abweichung jeweils auf der abweichenden Datei relativ zum Bezug (erste Quelle, `Dockerfile`); beide Werte und beide Dateien stehen in der Meldung. Beide Dateien werden tatsächlich gelesen (Proben 1 und 3 ändern nur das `Makefile`, 2 und 4 nur das `Dockerfile`).

### Spiegel der Modulliste

| Spiegel | Mutation | Beobachtet |
|---|---|---|
| `netlessDocModules()` | `"pinsync"` entfernt, `make test` | rot: `TestQA03_NetlessModuleList_Live`: `DC-QA-03-Netzlos-Modulliste verletzt: modules aktiviert "pinsync" — weder in netlessDocModules() noch in forbiddenInNetless() eingeordnet`; `make test` Exit 2. Richtiger Grund. |
| `FOCUS_DISABLE` (Makefile) | `--disable pinsync` entfernt, `make adr-check` | **nichts rot**: Exit 0; der Fokuslauf gibt lediglich zusätzlich `d-check: pinsync: 2 Pin(s) geprüft, 4 Quelle(n) gelesen` aus und zählt 1235 Dateien. Mit dem Disable fehlt diese Zeile. |

Beide Mutationen zurückgenommen (`git checkout`), `git diff` leer.

**Befund zu `FOCUS_DISABLE` (V-3):** Kein Test und kein Gate hält die Variable gegen die Modulliste. Das ist die schon geführte Beobachtung `BEO-ALL/modulliste-spiegel-ungegated`, hier nicht neu; der Plan-Nachtrag zählt den Spiegel zwar, kann ihn aber nur durch Hinsehen halten. Folge eines Fehlens ist nur eine Zeile mehr im Fokuslauf, kein falsches Grün — solange das Modul in einem fokussierten Lauf nichts findet. Findet es etwas (z. B. Pin-Drift), schlüge ein fokussierter Lauf wie `adr-check` mit fremdem Befund an.

## Review-Befunde R1

| Finding | Urteil | Beleg |
|---|---|---|
| F-1 LOW (Zählangabe „Zwölf Module“) | behoben | `grep 'Zwölf\|Elf Module' harness/sensors/doc-check.md` ⇒ kein Treffer; Vertrag lautet „Die Module der Liste `modules` der `.d-check.yml` …, darunter …" (keine Zahl mehr, die driften kann). `file` wird weiter nicht genannt, das deckt „darunter". |
| F-2 LOW (Makefile-`##` von `doc-check`) | behoben | Beschreibung nennt jetzt „Zitate und die Gleichheit der Doppel-Pins (pinsync)"; `file` fehlt weiterhin (Rest, kosmetisch). |
| F-3 INFO (`?=` überschreibbar) | behoben | `doc-check.md` Punkt 5: „Es liest den Dateistand: das `Makefile` setzt beide Werte mit `?=`, und ein Aufruf mit `GO_VERSION=…` oder eine Umgebungsvariable überschreibt sie zur Laufzeit, ohne dass das Modul es sieht." |

## Urteil je DoD-Punkt

Alle Haken im Plan stehen noch `[ ]` (Closure noch nicht erfolgt).

1. **`.d-check.yml` führt `pinsync` mit zwei Pins zu je zwei Quellen; je Ausdruck genau ein Treffer mit Zählbefehl im Plan-Nachtrag; `v` außerhalb der Fanggruppe:** **bestätigt.** Vier Ausdrücke im Block (`^ARG GO_VERSION=(\S+)$`, `^GO_VERSION\s*\?=\s*(\S+)$`, `^ARG GOLANGCI_LINT_VERSION=v?(\S+)$`, `^GOLANGCI_LINT_VERSION\s*\?=\s*v?(\S+)$`); eigene Zählung 1/1/1/1; Plan-Nachtrag nennt Befehle, Ergebnis und Zeilen (`Dockerfile:31/32`, `Makefile:13/14`); Probe 5/5b belegt die `v`-Gleichheit.
2. **Brechprobe beidseitig mit echter Ausgabe, Grund-Code und Fundstelle:** **bestätigt** (Proben 1 bis 8; zusätzlich Ausdruck-Brüche `no-match`/`ambiguous`). Vorbehalt: gefahren auf Kopie mit dem fertigen Image, nicht als `make doc-check` auf mutiertem Baum (Begründung oben). Der Selbstlauf desselben Images am Original ist grün; damit ist der Beleg gleichwertig.
3. **Gate-Verzeichnis nennt die Zusage mit Bindung:** **bestätigt.** `harness/README.md` Zeile `make doc-check` führt `DC-FA-PSYNC-001` samt `ADR-0113`; `harness/sensors/doc-check.md` Vertrag und Grenze Punkt 5 (zwei Paare, Wertgleichheit statt Aktualität, Digests/Spiegel ungeprüft, Quellen außerhalb der Scan-Menge — §3.8 beantwortet). Beide Links klickbar, `make doc-check` grün.
4. **`make gates` grün:** **bestätigt** (eigener Lauf, Exit 0).
5. **Review durchgeführt, Report liegt vor:** **bestätigt** (R1: 0 HIGH, 0 MEDIUM, 2 LOW, 1 INFO; alle behoben). Der Report ist noch untracked.
6. **Closure-Notiz mit Lerneintrag:** **offen** — §7 „Entsteht bei der Closure".
7. **Beobachtungs-Register fortgeschrieben:** **offen.** §8 nennt `BEO-ALL/pin-bump-mirrors-ungated` (löst der Slice nicht auf; Closure-Notiz muss das sagen) und `BEO-ALL/registry-vs-authority-table-drift`; Kandidat für Evidence: `BEO-ALL/modulliste-spiegel-ungegated` (V-3).
8. **Risiken §6 mit Ausgang:** **offen** (entstehen bei Closure). Beide Risiken sind erkennbar *entfallen*: Risiko 1 ist die Zusage selbst, Risiko 2 ist durch die Proben 1 bis 4 belegt (beide Dateien gelesen).
9. **Drei Paarungen:** **offen** (nächste Welle-Closure).

Plan §1 (Abgrenzung): Der Diff berührt nur `.d-check.yml`, `Makefile` (`FOCUS_DISABLE`, `##`-Text), `harness/README.md`, `harness/sensors/doc-check.md`, `gate_consistency_test.go` und den Plan — im Plan-Nachtrag vorab gemeldet; keine weiteren Pins, kein Modul-Code, keine Schwelle gesenkt.

Closure-Trigger §5: „DoD vollständig, `make gates` grün, beide Richtungen der Brechprobe mit echter Ausgabe im Report" — Brechprobe und Gates erfüllt durch diesen Bericht; DoD-Rest (Haken, Notiz, Register) Closure-Sache.

## Abweichungen und offene Punkte

- **V-1 (kein Befund am Slice, aber zu beachten):** `blackbox-probe REF=v0.90.0` ist **nicht** byte-identisch (4 von 20, Exit 2). Ursache: die Fixture `repo` liest die Konfiguration des Repos; das Vorher-Image kennt `pinsync` nicht und endet mit `.d-check.yml: … line 1088: field pinsync not found`, Exit 2 → Nachher Exit 0 mit pinsync-Zeile (`--json`: `summary.notes` mit der Bezugsmengen-Zeile). Die Annahme des Auftrags, die Repo-Konfiguration sei in den Fixtures nicht enthalten, trifft für `repo` nicht zu. Alle Läufe ohne `pinsync`-Konfiguration (16 Vergleiche) sind gleich. Der Default-Bytegleichheit der Anforderung ist damit weiter belegt (und in slice-287 mit 20/20); der Slice ändert allerdings die Eingabe der Fixture `repo`. Bei der Closure-Notiz benennen; die Probe ist kein DoD-Punkt dieses Slice.
- **V-2 (INFO):** `make image-test` hat keinen `pinsync`-Teil. Der Selbstlauf unter `:ro` mit `--user 4242:4242` und `--network none` ist per Hand grün (Exit 0, pinsync-Zeile); damit ist das aus slice-287 verschobene Risiko 3 (Betrieb unter `:ro`/fremdem Eigentümer) **belegt**, aber nicht als wiederholbarer Test. Das Gate selbst (`doc-check` in `gates`) läuft jedes Mal unter `:ro` und `--network none` und trägt die Messung dauerhaft; das Fremd-Benutzer-Merkmal tragen nur die `metrics`-Teile des `image-test`.
- **V-3 (INFO):** `FOCUS_DISABLE` ist von keinem Wächter gedeckt (s. o.); die Brechprobe zeigt keine rote Stelle. Nicht neu, `BEO-ALL/modulliste-spiegel-ungegated`.
- **V-4 (Kosmetik):** Makefile-`##`-Text nennt `file` nicht; alle DoD-Haken offen; Review-Report R1 und dieser Bericht untracked.

## Verdikt

Keine DoD-Verletzung. DoD-Punkte 1 bis 5 sind durch eigene Läufe bestätigt: `make doc-check` und `make gates` Exit 0, Brechproben rot aus dem richtigen Grund (`pinsync-mismatch`, `-no-match`, `-ambiguous`) und grün bei der `v`-Variante, Rücknahme sauber, R1-Befunde F-1 bis F-3 behoben. Die Punkte 6 bis 9 sind erwartungsgemäß offen (Closure). Bei der Closure festzuhalten: `blackbox-probe` weicht wegen der Fixture `repo` ab (V-1), der `:ro`-Fremdbenutzer-Selbstlauf ist manuell belegt (V-2), `FOCUS_DISABLE` ist ungewächtert (V-3).
