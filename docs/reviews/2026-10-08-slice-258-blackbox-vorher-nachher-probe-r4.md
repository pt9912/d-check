# Review R4 — slice-258: Black-Box-Vorher/Nachher-Probe als make-Target

- **Review-Art:** Code. Runde 4, nur über den Umbau nach R3 (Kanarienlauf statt
  Ausgabe-Klassifikation). Geprüft gegen den Slice-Plan (`slice-258`, §1 und die
  Plan-Änderung nach R3), die R3-Befunde R3-1 bis R3-3 und die Hard Rules
  [`AGENTS.md`](../../AGENTS.md) §3.1, §3.7 sowie §5 Regel 13 und 15. Die
  DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-258` · Range `cf59722e..3c18791c`, eingeschränkt auf
  `tools/blackbox-probe.sh`, `tools/blackbox-probe/README.md`,
  `harness/sensors/blackbox-probe.md` und den Slice-Plan. Plan-Änderung
  `e83697b3`, Fix-Commit `3c18791c`. Mitgelesen, nicht im Diff: das
  `blackbox-probe`-Target im `Makefile`, das Fixture `sauber`, die
  Commit-Botschaft von `3c18791c`.
- **Skill:** `reviewer.md` @ 1.18.0.
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-08
- **Eingangs-Kontext:** R1-, Verifikations-, R2- und R3-Report zu `slice-258`, alle
  vom 2026-10-08. Vorherige Findings am selben Werkzeug: F-1, V-1, R2-1 und R3-3 zur
  Klasse *stilles Grün bei gleichem Scheitern*, R3-2 zur Klasse *echter Lauf als
  Ausfall abgewiesen*.
- **Proben** (Docker-Shims und Klon im Scratchpad, kein Code geändert):
  - `make blackbox-probe REF=HEAD`: 20× `gleich`, `byte-identisch über 20
    Vergleiche (…, Kanarienlauf vorher und nachher)`, rc 0, 53 s.
  - `make blackbox-probe REF=v0.83.0`: 20× `gleich`, rc 0, 1 min 20 s.
  - `REF=v0.40.0` und `REF=v0.10.0`, `PROBE_FORMS=-`: beide Kanarienläufe grün,
    danach `4 von 5 Vergleichen weichen ab`, rc 1 (das alte Binary versteht die
    heutige `.d-check.yml` nicht, Grenze 3). Die Zusammenfassungs-Zeile
    `Datei(en) geprüft` steht laut `git log -S` seit dem CLI-Kern unverändert.
  - `PROBE_FORMS="-h --bogus --require-complete"`: 15× `gleich` (Exit 0 bzw. 2),
    rc 0.
  - Shim **Panic nur im Nachher** bei `--json` und **geänderte
    Nutzungsfehler-Meldung nur im Nachher** bei `--bogus`: 10× `ABWEICHUNG`, für
    `--bogus` mit Diff `-flag provided but not defined` / `+unbekannte Option`,
    für `--json` `exit(1→2)` bzw. `exit(0→2)`, rc 1.
  - **Klon mit `umask 077`** (`git clone` im Scratchpad, Verzeichnisse `drwx------`),
    Probe dort gestartet, `PROBE_FORMS="- --json"`: `gleich repo/- (exit 2)`,
    `gleich repo/--json (exit 2)`, `byte-identisch über 10 Vergleiche`, **rc 0**.
    Die Ausgabe des Werkzeugs für diesen Mount: `d-check: error: lstat
    /repo/docs: permission denied`, Exit 2.
  - Shim **leerer Mount nur für die Repo-Quelle** (`$PWD:/repo:ro`
    → leeres Verzeichnis, die Fixture-Mounts unter `$WORK` unverändert),
    `PROBE_FORMS="- --json"`: `gleich repo/- (exit 2)`, `gleich repo/--json
    (exit 2)`, **rc 0**.
  - Shim **geänderte Zusammenfassungs-Zeile nur im Nachher** (`Datei(en)
    geprüft` → `Dateien geprüft` auf stderr von `:latest`), `PROBE_FORMS=-`:
    `GESCHEITERT — Kanarienlauf vorher auf d-check:latest: Exit 0 — d-check: 1
    Dateien geprüft, 0 Befund(e); Umgebung traegt keinen Vergleich`, rc 2, **kein
    einziger Vergleich gelaufen**.

## Findings

### R4-1 — MEDIUM — Der Kanarienlauf prüft nur die Mount-Quelle der Fixtures, nicht die des Repos: unlesbares oder leeres Repo meldet wieder „byte-identisch"

- **quelle:** Skill Prüffragen 1 (Kontext: Werkzeug, kein Gate), 8 und 18;
  [`AGENTS.md`](../../AGENTS.md) §5 Regel 13 und 15; R2-1 (geschlossen in R3, durch
  diesen Umbau für das Repo wieder offen)
- **pfad:** `tools/blackbox-probe.sh` · „`# Daemon, Mounts und den Inhalt, den der Container sieht`"
  und „`fixtures+=("repo=$PWD")`"; `harness/sensors/blackbox-probe.md` · „Jeder andere
  Ausgang — auch Exit 2 oder ein Absturz — ist Verhalten des Werkzeugs und wird
  verglichen"; Commit `3c18791c` · „R3-3: ein Mount ohne Markdown faellt am
  Kanarienlauf auf"
- **befund:** Der Kanarienlauf mountet `$WORK/fx/sauber`, eine Kopie unter
  `mktemp -d`, die das Skript lesbar macht. Der Repo-Fall mountet dagegen `$PWD`,
  eine zweite Quelle, die weder kopiert noch lesbar gemacht wird. Sieht der
  Container dort nichts Lesbares, endet das Werkzeug auf beiden Seiten gleich mit
  Exit 2, und der Fall gilt seit dem Umbau als „gleich". Gemessen: Klon mit
  `umask 077` (`lstat /repo/docs: permission denied`) und Shim mit leerem Mount nur
  für die Repo-Quelle, beide mit `gleich repo/… (exit 2)` und rc 0. Ein Slice, der
  das Verhalten auf dem Repo ändert, bekommt dort ein Grün über einen Fall, in dem
  nichts geprüft wurde. Grenze 4, der Code-Kommentar und die Commit-Botschaft
  behaupten aber „Mounts" im Plural. Die Restgrenze nennt nur den zeitlichen Fall
  (Ausfall zwischen den Kanarien), nicht den quellen-abhängigen. Der Satz zur
  `umask 077`, der in R2 genau diesen Fall benannte, ist mit dem Umbau aus Grenze 4
  verschwunden. R3-3 ist damit nur für die Fixtures geschlossen. Die gemessenen
  Shims ersetzten jeden `/repo`-Mount gleichzeitig und konnten die Asymmetrie
  deshalb nicht zeigen. Nicht gemessen, aber dieselbe Klasse: Ein Ausfall, der vom
  Umfang der Eingabe abhängt (etwa ein Speicher-Kill mit 137 auf beiden Seiten beim
  großen Repo-Lauf), passiert den Kanarienlauf über das kleinste Fixture ebenfalls.
- **verifizierbar:** ja. `(umask 077; git clone -q <repo> <dir>)`, im Klon
  `REF=HEAD PROBE_FORMS=- bash tools/blackbox-probe.sh` liefert `gleich repo/-
  (exit 2)` und rc 0.
- **klasse:** stilles-gruen-bei-gleichem-scheitern

### R4-2 — MEDIUM — Ein Vorgang, der die Zusammenfassungs-Zeile oder das Verhalten auf `sauber` ändert, bekommt keinen Vergleich, sondern „Umgebung trägt keinen Vergleich"; die genannte Abhilfe kann nicht greifen

- **quelle:** Skill Prüffragen 18 und 19; [`AGENTS.md`](../../AGENTS.md) §5
  Regel 13; R3-2 (dieselbe Fehldiagnose, andere Ursache)
- **pfad:** `tools/blackbox-probe.sh` · „`|| ! grep -q '^d-check: 1 Datei(en) geprüft, 0 Befund(e)' "$WORK/kanarie.err"; then`"
  und „`kanarie "vorher"`"; `harness/sensors/blackbox-probe.md` · „Ändert ein
  Vorgang die Zusammenfassungs-Zeile des Werkzeugs, passt er den Kanarienlauf mit
  an."
- **befund:** Der Kanarienlauf hält beide Images gegen dieselbe feste Zeile, und
  er läuft auf beiden Images, bevor der erste Vergleich startet. Ändert ein Slice
  die Zusammenfassungs-Zeile oder das Ergebnis auf `sauber` (etwa ein neues
  Default-Modul, das dort einen Befund meldet), scheitert der Nachher-Kanarienlauf
  sofort. Die Probe endet dann mit rc 2 und „Umgebung traegt keinen Vergleich",
  und zwar ohne eine einzige `ABWEICHUNG`-Zeile (Shim gemessen). Die Verhaltensänderung,
  die die Probe finden soll, erscheint wieder als Infrastruktur-Ausfall. Die
  Abhilfe in Grenze 4 trägt nicht: Wer den Kanarienlauf an die neue Zeile anpasst,
  lässt den Vorher-Kanarienlauf scheitern, weil das Vorher-Binary aus `REF` die alte
  Zeile schreibt. Gegen den Stand vor der Änderung ist ein solcher Slice mit der
  Probe nicht vergleichbar. Gegen alte REFs ist der Kanarienlauf dagegen stabil
  (gemessen bis `v0.10.0`), falsch rot aus dem Vorher-Stand allein trat nicht auf.
- **verifizierbar:** ja. Shim, der nur auf stderr von `d-check:latest` `Datei(en)
  geprüft` umschreibt, `PROBE_FORMS=-`: rc 2, Abbruch im Kanarienlauf „vorher",
  null Vergleiche.
- **klasse:** echter-lauf-als-ausfall-abgewiesen

## Schließung R3

| R3 | Stand | Beleg |
|---|---|---|
| R3-1 | geschlossen | `--bogus`/`--require-complete`: 10× `gleich (exit 2)`, rc 0. Geänderte Meldung nur im Nachher: `ABWEICHUNG … stderr` mit Diff, rc 1. |
| R3-2 | geschlossen | `-h`: `gleich (exit 0)`. Panic nur im Nachher: `ABWEICHUNG … exit(1→2)`, rc 1, keine Fehldiagnose. |
| R3-3 | für die Fixtures geschlossen, für das Repo nicht | Fremder oder leerer Mount unter `$WORK` scheitert am Kanarienlauf. Für die Quelle `$PWD` gilt das nicht → R4-1. |

## Negativbefunde

- **Falsch rot durch alte REFs:** geprüft, ohne Befund. `v0.83.0`, `v0.40.0` und
  `v0.10.0` bestehen beide Kanarienläufe. Die Zeile `Datei(en) geprüft, 0
  Befund(e)` ist laut `git log -S` seit dem CLI-Kern unverändert. Das Fixture kommt
  aus dem Arbeitsbaum, nicht aus `REF`, und fehlt deshalb in keinem Vorher-Stand.
  `--network none` ist ein Docker-Schalter und hängt nicht am `REF`.
- **Zeitlicher Ausfall zwischen den Kanarien:** geprüft. Er ist als Restgrenze
  benannt. Fällt der Daemon für beide Seiten eines Falls aus und kehrt vor dem
  Nachher-Kanarienlauf zurück, bleibt der Fall gleich. Trifft der Ausfall nur
  eine Seite, ist das eine `ABWEICHUNG`, also laut.
- **Abbruch bei 125/126/127:** geprüft, ohne Befund. `fail` läuft in der
  Haupt-Shell (`lauf` ist kein Subshell-Aufruf), Exit 2 kommt an, und `trap`
  räumt `$WORK` auf.
- **Exit-Präzedenz:** geprüft, ohne Befund. Ein gescheiterter Nachher-Kanarienlauf
  nach Abweichungen endet mit 2, die `ABWEICHUNG`-Zeilen stehen davor in der
  Ausgabe.
- **Veraltetes Nachher-Image:** geprüft, ohne Befund. Das Target baut `:latest`
  vor dem Skript neu (`$(MAKE) build VERSION=0.0.0-dev`).
- **Exit-Tabelle gegen Code:** geprüft, ohne Befund außer R4-1. Alle Zeilen haben
  ihren `fail`-Pfad. Zeile 0 stimmt dem Wortlaut nach, liest sich wegen R4-1 aber
  stärker, als der Repo-Fall trägt.
- **Kommentare** (§3.7): geprüft, ohne Befund in der Form. Beide neuen Blöcke im
  Skript tragen Zusage bzw. Abgrenzung, keine Befund-Nummern und keine
  Review-Historie. Auch der README-Satz ist eine Kopplung. Inhaltlich überdehnt
  der Kanarien-Kommentar seine Zusage („Mounts") → R4-1.
- **Slice-Plan:** geprüft, ohne Befund. Die Plan-Änderung `e83697b3` liegt vor dem
  Code-Commit `3c18791c` und weitet die Abgrenzung nicht. Ihr Satz „belegt Daemon,
  Mounts und gesehenen Inhalt" teilt die Überdehnung aus R4-1.
- **Docker/make-only** (§3.1): geprüft, ohne Befund. Neu ist nur `grep`.

## Kategorie-Summary

HIGH 0 · MEDIUM 2 (R4-1, R4-2) · LOW 0 · INFO 0.
Wiederkehrende Finding-Klassen:

- `stilles-gruen-bei-gleichem-scheitern`, fünftes Auftreten in diesem Slice (F-1,
  V-1, R2-1, R3-3, R4-1), jetzt wieder als MEDIUM.
- `echter-lauf-als-ausfall-abgewiesen`, zweites Auftreten (R3-2, R4-2).

## Verdikt

R3-1 und R3-2 sind sachlich geschlossen. Der Kanarienlauf ist gegen alte REFs
robust, und der Vergleich von Exit 2 und Abstürzen funktioniert. Offen sind zwei
MEDIUM:

- **R4-1:** Das Repo, die zweite Mount-Quelle, ist von der Umgebungsprüfung nicht
  gedeckt, und dort kehrt das stille Grün aus R2-1 zurück. Das ist eine Regression
  gegenüber R3.
- **R4-2:** Gerade eine Änderung an der Zusammenfassungs-Zeile oder am Ergebnis auf
  `sauber` macht die Probe unbenutzbar und wird als Umgebungsfehler gemeldet.

Beide blockieren nach Skill typischerweise. Das Target ist deklariert kein Gate,
deshalb genügt die Einarbeitung vor der Closure.

**Steering-Loop-Signal** (Skill §Kontext-Eskalation): Die Klasse
`stilles-gruen-bei-gleichem-scheitern` erreicht in diesem Slice das fünfte
Auftreten. Auch dieser Umbau wurde mit Shims belegt, die alle Mounts gleichzeitig
ersetzten. Das Skript hat aber zwei Mount-Quellen mit unterschiedlicher
Behandlung (Kopie mit `chmod` gegen `$PWD`). Die Fallmenge, die eine
Umgebungsprüfung decken muss, ist die Menge der Quellen, nicht die Menge der
Ausfallarten.
