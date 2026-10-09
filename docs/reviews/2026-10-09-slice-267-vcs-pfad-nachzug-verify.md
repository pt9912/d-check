# Verifikation: slice-267, `vcs` lässt einen reinen Pfad-Nachzug in immutablen Dateien durch

- **Rolle:** Verifier (Modul 8/11). Geprüft wird der Stand gegen DoD (§2), Plan (§3 samt beiden
  Plan-Änderungen vor dem Code), Abgrenzung (§1) und §6, **nicht** der Diff gegen die
  Entscheidungen. Das haben R1 bis R8 getan.
- **Gegenstand:** `slice-267`, Commits `39f91a74` bis `00cc1f45` (`git log 039eaafa..HEAD`, alle
  tragen `slice-267`). Code-Stand `HEAD` = `00cc1f45`, Arbeitsbaum sauber.
- **Eingang:** Slice-Plan in `in-progress/`; Reports R1 bis R8; Lastenheft 0.102.6
  `DC-FA-VCS-001` (Absatz „Pfad-Nachzug (opt-in)", Kriterium „Boundary (Pfad-Nachzug)");
  Spezifikation `DC-FA-VCS-001.a` Schritt 4 und §2-Zeile `vcs.ignore-link-targets`; `ADR-0103`;
  `harness/sensors/adr-check.md` Grenze 4.
- **Datum:** 2026-10-09 · **Modell-ID:** claude-opus-5-5

## Verdikt

**Die Liefer-Punkte der DoD (§2, Punkte 1 und 2) sind bestätigt.** Belegt am Code, am gebauten
Image und durch bewusstes Brechen. Punkt 3 (Review, Verifikation) ist mit R1 bis R8 und diesem
Bericht erfüllt, mit einer Lücke: Der letzte Fix-Commit `00cc1f45` liegt hinter R8 und ist noch
nicht reviewt (V-2). **Offen ist Punkt 4**, die Closure selbst. Keine DoD-Verletzung. Ein Befund
LOW zum Plan-vs-Code-Diff (V-1), sonst INFO.

## Sensor-Belege (selbst gefahren)

| Lauf | Ergebnis |
|---|---|
| `make gates` (auf `00cc1f45`) | Exit 0. Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `lint` „0 issues."; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; `doc-check` 1067 bzw. 1061 Dateien, 0 Befunde |
| `make blackbox-probe REF=039eaafa` | 16/16 Fixture-Vergleiche gleich. Die 4 Abweichungen betreffen nur `repo/*`, alt Exit 2 (`field ignore-link-targets not found`), neu Exit 0. Das ist erwartet: Die `.d-check.yml` des Repos trägt den neuen Schlüssel. Der eigentliche Opt-in-Beleg steht unten |

## DoD Punkt 1: Nachzug geht durch, Linktext und Zeile bleiben Befund

**Am Code:** `VCSConfig.IgnoreLinkTargets` wird vom YAML-Adapter durchgereicht
(`rawVCS.IgnoreLinkTargets`). `CheckVCS` bildet nur mit Schlüssel einen Pfad-Baum
(`pathTree` über `baseAll` ∪ `headAll`, also **alle** Pfade beider Stände, nicht nur die
geschützte Klasse). `vcsCore` ersetzt nur Zeilen, in denen `normalizedLinkTargetLines` ein
auflösendes Ziel durch Marke + `path.Base` + Anker ersetzt hat. Ohne Schlüssel ist `resolve` nil
und `vcsCore` läuft wie zuvor. Lastenheft-Absatz, Spezifikation Schritt 4, §2-Zeile und
`--print-config`-Vorlage stehen da.

**Kriterium „Boundary (Pfad-Nachzug)" am gebauten Image.** Gebaut mit `make build` vom Stand
`00cc1f45`. Wegwerf-Klon des Repos, Lauf wie `make adr-check` (`--enable vcs` +
`FOCUS_DISABLE`, `--range <base>..<head>` mit SHAs). Basis-Commit: eine zusätzliche
`Accepted`-Probe-ADR mit Referenz-Definition `[rel]: ../../user/releasing.md` und Bild
`![Logo](../../user/releasing.md)` (im Bericht in Code-Spans):

| # | Änderung | Ergebnis |
|---|---|---|
| 1 | `git mv docs/user/releasing.md docs/user/maintainer/releasing.md`; die vier Body-Links in `ADR-0014` und Ref-Def + Bild der Probe-ADR nachgezogen | 0 Befunde, Exit 0 |
| 1s | dasselbe gestagt, `--staged` (Pfad des `pre-commit`-Hooks) | 0 Befunde, Exit 0 |
| 1o | Fall 1 **ohne** den Schlüssel | 2 × `core-drift-vcs` (`ADR-0014`, Probe-ADR), Exit 1 |
| 1a | Fall 1 ohne Schlüssel, **altes** Image (`039eaafa`) | dieselben 2 Befunde, Exit 1 |
| 2 | ohne Umzug: Ziel `../../../docs/user/releasing.md` → `../../../docs/user/operations.md` | `core-drift-vcs`, Exit 1 |
| 3 | Ziel → `../../../docs/user/fehlt/releasing.md` (existiert in keinem Stand) | `core-drift-vcs`, Exit 1 |
| 4 | Umzug + Nachzug wie 1, zusätzlich Linktext `` `releasing.md` `` → `` `release-prozess.md` `` | `core-drift-vcs`, Exit 1 |

Damit ist jede Hälfte des Kriteriums am Image bestätigt: Inline-Link, Bild und Referenz-Definition,
anderer Name, nicht existierendes Ziel, Linktext und „ohne Schlüssel meldet schon der reine
Nachzug".

**Bewusstes Brechen** (10 Mutationen, je eine Wegwerf-Kopie aus `git archive HEAD`, `make test`
dort, jeweils Exit 2):

| # | Mutation | rot, Grund |
|---|---|---|
| B1 | Schlüssel wirkungslos (`if false && cfg.IgnoreLinkTargets`) | **Kernzusage, Richtung „geht durch":** jeder Nachzug-Fall mit Soll 0 fällt mit `ignore-link-targets=true: Befunde = 1, want 0` (`inline-ziel nachgezogen`, `ziel mit anker …`, `bild …`, `referenz-definition …`, …) |
| B2 | Existenzprüfung `!tree[p]` entfernt (Vorzustand „leeren", nur Form) | **Kernzusage, Richtung „bleibt Drift":** `…Umzug/neues ziel fehlt in beiden staenden`, `…/neues ziel nur dateiname`, `verzeichnis eines nicht aufloesenden ziels geaendert`, `escapte zielklammer geaendert`, je `Befunde = 0, want 1` |
| B3 | Normierung auf vollen Pfad statt `path.Base` | alle Nachzug-Fälle `Befunde = 1, want 0` |
| B4 | Anker nicht mehr Teil der normierten Form | `anker geaendert`: `Befunde = 0, want 1` |
| B5 | `markFreeTree` gibt immer den Baum zurück (Fix `00cc1f45` zurückgenommen) | `…Umzug/neues ziel traegt die marke`: `Befunde = 0, want 1` |
| B6 | Baum nur aus BASE statt Vereinigung | `…Umzug/nachzug nach umzug`, `…/nachzug auf datei nur in base`: `Befunde = 1, want 0` |
| B7 | `opaqueLines` leer | 13 Fälle, u. a. `eingerueckter code`, `html-block`, `code im zitat`, `fence-schliesser zu kurz`, `zitat-link nachgezogen (fail-safe)`, je `want 1` |
| B8 | Prüfung des Zielausdrucks (`linkDestRE`) entfernt | `klammertext mit leerraum`, `offener titel`, `spitzklammer mit rest`, je `Befunde = 0, want 1` |
| B9 | YAML-Durchreichung entfernt | `TestDecode_VCSIgnoreLinkTargets`: „ignore-link-targets nicht durchgereicht" |
| B10 | `containsLink` entfernt | `klammer um inneren link geaendert`: `Befunde = 0, want 1` |

Jede Mutation fällt mit der behaupteten Meldung, keine nur „irgendwie". Ein erster Versuch an B4
scheiterte am Compiler (ungenutzte Variable) und zählt nicht. Er wurde so korrigiert, dass der Test
aus dem richtigen Grund fällt.

## DoD Punkt 2: Opt-in byte-identisch, Repo-Nutzung, `AGENTS.md`, Gates

**Ohne den Schlüssel byte-identisch, mit aktivem `vcs`:** Im Wegwerf-Klon wurde
`ignore-link-targets: true` aus der `.d-check.yml` entfernt. Altes Image (`git archive 039eaafa`)
und neues Image (`HEAD`), beide `VERSION=0.0.0-dev`, liefen über fünf Ranges (Fälle 1, 1s, 2, 3, 4,
als SHAs) × vier Formen (Standard, `--json`, `--yaml`, `--doctor`) × zwei Modulsätze (Fokus wie
`adr-check`, voll mit `--enable vcs`). Das sind **40 Läufe, stdout, stderr und Exit getrennt mit
`cmp`: 0 Abweichungen**, alle mit echtem Exit 1, also keine Leerläufe.

`.d-check.yml` setzt `ignore-link-targets: true` im `vcs`-Block, Kommentar mit Verweis auf
`ADR-0103`. `AGENTS.md` §3.5, `harness/README.md` §Traceability rules und
`harness/sensors/adr-check.md` (Vertrag und Grenze 4) nennen den Nachzug. `make gates` ist grün
(oben).

## Plan-vs-Code-Diff

- **§3, alle Zeilen geliefert:** Kern-Regel (`vcs.go`, neu `vcs_link_targets.go`), Modell, Adapter,
  Vorlage, Lastenheft, Spezifikation, `.d-check.yml`, `AGENTS.md` §3.5, `harness/README.md`,
  `harness/sensors/adr-check.md`, `ADR-0103` + Index-Zeile.
- **Über §3 hinaus:** nur `.harness/skills/reviewer.md`, ein Zeilen-Nachzug einer
  `d-check:cite`-Spanne (`AGENTS.md:286` → `287`), erzwungen durch die eine neue Zeile in §3.5. Kein
  Inhalt.
- **§1-Abgrenzung gehalten:** `docs/user/releasing.md` liegt unverändert an seinem Ort (slice-268).
  Kein Template-Feld-Nachzug. Linktext- und Zeilen-Änderungen bleiben Drift (Image Fall 4, B1/B2).
- **Rückführung aus §4** (Zeilen-Paarung) nicht eingetreten. Die Normierung braucht keine Paarung,
  wie die zweite Plan-Änderung sagt.

## Benannte Grenzen gegen den Code

Spezifikation Schritt 4, `harness/sensors/adr-check.md` Grenze 4 und der Kopfkommentar von
`normalizedLinkTargetLines` stimmen untereinander und mit dem Code überein. Am Image geprüft
(Umzug von `releasing.md` + Nachzug in einer Probe-ADR, jeweils in Code-Spans beschrieben):

| Form | benannt als | gemessen |
|---|---|---|
| Link in einer Fence, die hinter einer Listenmarke öffnet (`- ` + drei Backticks) | geht durch | 0 Befunde ✓ |
| Referenz-Definition mit CRLF-Ende | Drift (fail-safe) | `core-drift-vcs` ✓ |
| Inline-Link mit CRLF-Ende (Kontrolle) | — | 0 Befunde |
| Link in Zitat `> [a](…)` | Drift | `core-drift-vcs` ✓ |
| eingerückter Listen-Folgeabsatz | Drift | `core-drift-vcs` ✓ |
| Absatz beginnt mit Inline-HTML | Drift | `core-drift-vcs` ✓ |
| Listenpunkt `- [a](…)` (Kontrolle) | geht durch | 0 Befunde ✓ |
| Ziel mit Query `…releasing.md?x=1` | nur im Code-Kommentar (`linkTargetResolver`) | `core-drift-vcs` (V-3) |
| prozent-kodiertes Ziel `…releasing%2Emd` | nirgends | `core-drift-vcs` (V-3) |

Die Grenze „Nachzug auf eine Datei, die nur noch in BASE existiert, geht durch" ist durch
`…Umzug/nachzug auf datei nur in base` gehalten (B6 zeigt, dass der Fall an der Vereinigung hängt).
„Gleicher Name, andere Datei" ist benannt und von R7 I-2 vermessen.

## Befunde

- **V-1 (LOW): Die zweite Plan-Änderung beschreibt eine Auflösung, die der Code nicht mehr hat.**
  §3 sagt: „nur, wenn es im jeweiligen Stand als Datei oder Verzeichnis auflöst (BASE gegen den
  BASE-Baum, HEAD gegen den HEAD-Baum)". Der Code löst seit `056dc4cf` (nach R7) beide Seiten gegen
  die **Vereinigung** auf. So steht es auch in Spezifikation, Lastenheft 0.102.6 und
  `ADR-0103`-Geschichte. Der Wechsel kam als Review-Fix, nicht als dritte Plan-Änderung. Er ist
  fail-safe begründet (ein unveränderter Link auf eine gelöschte Datei bliebe sonst Drift), aber er
  öffnet die benannte BASE-only-Grenze. Vorschlag: In §3 eine Zeile nachtragen oder die
  Formulierung in der Closure-Notiz unter „Was ging anders als geplant" benennen. Der Entwurf in §8
  („`[text]()`", Leerung) ist erkennbar Vor-Code-Entwurf und durch die Plan-Änderung überholt. Kein
  eigener Befund.
- **V-2 (LOW): `00cc1f45` ist von keinem Review gesehen.** R8 prüfte `b60bf321..dc1fa440`. Der Fix
  zu R8 M-1 (`markFreeTree`), die Kopfkommentar-Ergänzung und der neue Test kamen danach. Das
  Verhalten habe ich bestätigt (B5, B6, `make gates`), den Diff auf Kommentar-Klassen und Spiegel
  nicht, weil das nicht die Frage dieser Rolle ist. Vor der Closure zu entscheiden: eine kurze R9
  über den einen Commit, oder die Lesart „Verifikation deckt ihn" in der Closure-Notiz benennen.
- **V-3 (INFO): Zwei fail-safe Formen fehlen in den Grenzen-Listen.** Ein Ziel mit Query
  (`?`) löst nie auf, das steht nur im Kommentar von `linkTargetResolver`. Ein prozent-kodiertes
  Ziel löst nie auf, weil der Resolver nicht dekodiert, `ResolveTarget` des Moduls `links` aber
  schon (vgl. R7 I-1). Das steht nirgends. Beide bleiben Drift, also in sicherer Richtung. Die Liste
  „Fail-safe bleibt Drift" in Spezifikation und Sensor-Datei ist nur nicht vollständig.
- **V-4 (INFO): Der Körper von `ADR-0103` beschreibt noch die Leerung.** Entscheidung 1
  („wird … geleert") und die zweite Konsequenz-Grenze („ein Link-Ziel in Inline-Code oder einem
  Codeblock wird ebenso geleert") treffen auf den Code nicht mehr zu. Berichtigt ist das nur in der
  Geschichte, und das ist der zulässige Weg (§3.5, `make adr-check` hält es). Ein Leser des Körpers
  erfährt das aber erst am Tabellenende. Das ist eine Folge davon, dass die ADR vor dem
  Code-Stand `Accepted` wurde. Möglicher Lerneintrag, kein Handlungsbedarf am Artefakt.
- **Außerhalb des Slice beobachtet (kein Befund hier):** Im Wegwerf-Klon meldeten alte wie neue
  Images für `--range <sha>..<lokaler-Branchname>` „Range-Leerfall … 0 Commits", Exit 2, obwohl ein
  Commit dazwischen lag. Mit SHAs lief alles korrekt. Das Verhalten ist fail-closed und kein stilles
  Grün, aber womöglich eine Ref-Auflösungs-Lücke. Gesondert ansehen.

## Offen bis zur Closure (DoD Punkt 4, kein Verifier-Befund)

- §7 Closure-Notiz, Lerneintrag, Register und drei Paarungen sind leer.
- **§6-Risiko „Nachzug als Tarnung" ohne Ausgang.** Es ist durch die Normierung auf
  auflösende Ziele verengt, aber nicht beseitigt: Gleicher Dateiname in einem anderen Verzeichnis
  geht durch (R7 I-2: 111 von 1289 Ziel-Tokens in `Accepted`-ADRs haben Namensvettern), ebenso ein
  Ziel, das nur noch in BASE existiert. Benannt ist das in Lastenheft, Spezifikation, Sensor-Datei
  und `ADR-0103`. *Entfallen* trägt nicht, weil das Risiko weiter eintreten kann. Passend wirkt
  *weiter offen* (Register) oder *eingetreten* mit der Re-Evaluierungs-Bedingung von `ADR-0103` als
  Wächter.
- Die DoD-Haken in §2 sind noch nicht gesetzt.
- `ADR-0103` §Fitness Function nennt `make adr-check` über die Range des Umzugs (slice-268) als
  Beleg am Bestand. Der ist erst mit slice-268 einlösbar. Der Image-Fall 1 oben ist der
  vorweggenommene Beleg auf einem Klon.

## Negativbefund

Gesucht und nicht gefunden:

- eine Abweichung alt/neu ohne Schlüssel (40 Vergleichsläufe mit aktivem `vcs`, 16 Fixture-Läufe);
- ein Pfad, auf dem der Schlüssel ohne Konfiguration wirkt (`tree` nur unter
  `cfg.IgnoreLinkTargets`, sonst `resolve == nil`);
- eine Auflösung nur gegen die geschützte Klasse statt gegen alle Pfade (`AllPaths`, belegt durch
  Image-Fall 1 mit einem Ziel außerhalb von `docs/plan/adr/`);
- ein Unterschied zwischen Range- und `--staged`-Lauf (Fall 1/1s);
- ein Lastenheft-Kriterium „Boundary (Pfad-Nachzug)" ohne Test oder Image-Beleg;
- eine benannte Grenze, die am Image anders ausfällt als beschrieben;
- ein Test, der beim Brechen nur aus fremdem Grund fällt.

Klon, Mutationskopien und die vier Probe-Images (`dcheck-v267`, `-old`, `-dev`,
`d-check:probe-vorher`) sind entfernt.
