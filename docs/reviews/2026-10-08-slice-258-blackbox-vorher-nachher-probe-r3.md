# Review R3 — slice-258: Black-Box-Vorher/Nachher-Probe als make-Target

- **Review-Art:** Code. Runde 3, nur über die Einarbeitung der R2-Befunde.
  Geprüft gegen den Slice-Plan (`slice-258`, §1 und die Plan-Änderung nach R2),
  die R2-Befunde R2-1 bis R2-4 und die Hard Rules
  [`AGENTS.md`](../../AGENTS.md) §3.1, §3.7 sowie §5 Regel 13. Die DoD-Abhakung
  prüft dieses Review nicht.
- **Gegenstand:** `slice-258` · Range `f24fd1ef..1b5d9b60`, eingeschränkt auf
  `tools/blackbox-probe.sh`, `harness/sensors/blackbox-probe.md` und den
  Slice-Plan. Fix-Commit `1b5d9b60`, Plan-Änderung `3520840f`. Mitgelesen, nicht
  im Diff: `spec/spezifikation.md` §DC-FA-CLI-001.a (Fehlermodi des Werkzeugs),
  die CLI-Hilfe des Images, `cmd/` und `internal/` (Suche nach `recover()`).
- **Skill:** `reviewer.md` @ 1.18.0.
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-08
- **Eingangs-Kontext:** R1-, Verifikations- und R2-Report zu `slice-258`, alle vom
  2026-10-08. Vorherige Findings am selben Werkzeug: R1 F-1, V-1 und R2-1 zur
  Klasse *stilles Grün bei gleichem Scheitern*, R2-3 zur Klasse
  *Allaussage über gemessene Teilmenge*.
- **Proben** (Docker-Shims im Scratchpad, kein Code geändert):
  - `make blackbox-probe REF=HEAD`: 20× `gleich`, `byte-identisch über 20
    Vergleiche`, rc 0, 1 min 23 s.
  - Shim **leerer Mount** (Quelle jedes `…:/repo:ro` durch ein leeres
    Verzeichnis ersetzt), `PROBE_FORMS="- --json"`: 10× `NICHT PRUEFBAR … Scan-Wurzel
    ist leer`, `GESCHEITERT — 10 Fall/Faelle …`, **rc 2**. R2-1 ist damit geschlossen.
  - Shim **toter Socket** (`DOCKER_HOST=unix:///nonexistent.sock` nur für `run`):
    rc 2, `Exit 1 ohne Ausgabe des Werkzeugs — kein Lauf des Werkzeugs: failed to
    connect to the docker API …`.
  - Shim **Exit 125**: rc 2, `Exit 125 — kein Lauf des Werkzeugs`.
  - `PROBE_FORMS="--repair --repair-broad --trace"`: 15× `gleich`, rc 0. R2-3 ist
    damit geschlossen.
  - `PROBE_FORMS=--bogus` und `PROBE_FORMS=--require-complete`: je 5× `NICHT
    PRUEFBAR` (`flag provided but not defined` bzw. `--require-complete erfordert
    --trace`), rc 2.
  - Shim **geänderte Nutzungsfehler-Meldung nur im Nachher** (stderr des
    `:latest`-Laufs: `flag provided but not defined` → `unbekannte Option`),
    `PROBE_FORMS=--bogus`: 5× `NICHT PRUEFBAR … unbekannte Option: -bogus`,
    `0 verglichen, 0 davon abweichend`, rc 2. Kein `ABWEICHUNG`, kein Diff.
  - `PROBE_FORMS=-h`: rc 2, `Exit 0 ohne Ausgabe des Werkzeugs — kein Lauf des
    Werkzeugs: https://raw.githubusercontent.com/…`. Die Usage steht auf stderr,
    ihre erste Zeile lautet `d-check — prüft …`, nicht `d-check:`.
  - Shim **Laufzeitabbruch nur im Nachher** (Exit 2, stderr `panic: runtime
    error …` / `goroutine 1 [running]:`, stdout leer): rc 2, `d-check:latest …:
    Exit 2 ohne Ausgabe des Werkzeugs — kein Lauf des Werkzeugs: goroutine 1
    [running]:`. `grep -rn 'recover()' cmd/ internal/` liefert keinen Treffer: Ein
    Go-Panic des Werkzeugs endet mit genau dieser Form, Exit 2 und ohne
    `d-check:`-Zeile.
  - Shim **Mount auf ein nicht leeres Verzeichnis ohne Markdown**,
    `PROBE_FORMS="- --json"`: 10× `gleich (exit 0)`, `byte-identisch über 10
    Vergleiche`, **rc 0**.

## Ausgangs-Matrix eines Laufpaars

L heißt Lebenszeichen: stdout nicht leer oder eine Zeile `^d-check:` auf stderr.

| Vorher | Nachher | Klassifikation | Richtig? |
|---|---|---|---|
| Exit ∉ {0,1,2} | beliebig (oder umgekehrt) | `fail` „kein Lauf" → 2 | ja für Docker-Exits 125/126/127; ein OOM-Kill (137) des Werkzeugs selbst wird ebenfalls als „kein Lauf" gemeldet |
| {0,1,2} ohne L | beliebig (oder umgekehrt) | `fail` „kein Lauf" → 2 | ja für toten Daemon (1). **Nein** für `-h` (Exit 0, Usage auf stderr) und für einen Go-Panic (Exit 2): Das Werkzeug lief → R3-2 |
| 2 mit L | 2 mit L | `NICHT PRUEFBAR` → 2 | ja für leeren Mount, fehlende Leserechte, kaputte Konfiguration. **Nein** für einen Nutzungsfehler, den die Form selbst auslöst; abweichende Meldungen gehen ohne Diff unter → R3-1 |
| 2 mit L | 0/1 mit L (oder umgekehrt) | `ABWEICHUNG` → 1 | ja |
| 0/1 mit L | 0/1 mit L, gleicher Exit | byte-Vergleich | ja, solange die Eingabe die beabsichtigte ist. Der spezifizierte Pfad „0 Datei(en) geprüft, Exit 0" wird gleich → R3-3 |
| 0/1 mit L | 0/1 mit L, anderer Exit | `ABWEICHUNG` → 1 | ja |

Bei Exit 0 der Probe haben alle Fälle auf beiden Seiten denselben Exit 0 oder 1
mit Lebenszeichen. Ein realistischer stiller Grün-Pfad durch gleiches Scheitern
bleibt nur über R3-3, und der setzt einen Daemon voraus, der unter dem
Client-Pfad ein anderes, nicht leeres Verzeichnis sieht.

## Findings

### R3-1 — MEDIUM — „Exit 2 auf beiden Seiten = Eingabe nicht prüfbar" stimmt für Nutzungsfehler nicht; abweichende Fehlermeldungen verschwinden ohne Diff

- **quelle:** Skill Prüffragen 18 und 19; [`AGENTS.md`](../../AGENTS.md) §5
  Regel 13; `harness/sensors/blackbox-probe.md` Grenze 4
- **pfad:** `tools/blackbox-probe.sh` · „`# Exit 2 auf BEIDEN Seiten heisst: das Werkzeug konnte die Eingabe auf`"
  und `harness/sensors/blackbox-probe.md` · „hat das Werkzeug die Eingabe nirgends geprüft"
- **befund:** Exit 2 ist beim Werkzeug auch der Code für Nutzungsfehler
  (Spezifikation §DC-FA-CLI-001.a), und die Form selbst ist eine Eingabe.
  `PROBE_FORMS=--bogus` und `--require-complete` landen deshalb als „nicht
  prüfbar" (rc 2). Dabei hat das Werkzeug geprüft, und zwar die Argumente. Ändert
  ein Slice eine Fehlermeldung, liefert die Probe nur `NICHT PRUEFBAR` mit der
  Meldung des Nachher-Laufs und `0 davon abweichend`. Gemessen mit einem Shim, der
  nur die Nachher-Meldung umschreibt: kein `ABWEICHUNG`, kein Diff. Die Kehrseite
  der Regel fehlt in Grenze 4: Fehlerpfade (Exit 2 auf beiden Seiten) kann die
  Probe grundsätzlich nicht vergleichen. Laut, kein stilles Grün. Es ist dieselbe
  Klasse wie R2-3: Die Allaussage über „Exit 2" wurde nur gegen die Fälle geprüft,
  die die Regel beheben soll.
- **verifizierbar:** ja. `REF=HEAD PROBE_FORMS=--bogus bash tools/blackbox-probe.sh`
  liefert rc 2 mit `NICHT PRUEFBAR`. Mit dem Shim, der die Nachher-stderr
  umschreibt, bleibt die Ausgabe dieselbe und zeigt keinen Diff.
- **klasse:** allaussage-ueber-gemessene-teilmenge

### R3-2 — MEDIUM — Das Lebenszeichen weist echte Läufe als „kein Lauf des Werkzeugs" ab: `-h` und Laufzeitabbruch

- **quelle:** Skill Prüffrage 19; `harness/sensors/blackbox-probe.md` Grenze 4
  und Exit-Tabelle (1 = „mindestens eine Abweichung")
- **pfad:** `tools/blackbox-probe.sh` · „`0|1|2) [ -s "$4.out" ] || grep -q '^d-check:' "$4.err" || why="Exit ${rc} ohne Ausgabe des Werkzeugs" ;;`"
- **befund:** Zwei echte Ausgänge des Werkzeugs tragen das Lebenszeichen nicht.
  Erstens `-h`: Exit 0, die Usage steht auf stderr ohne `d-check:`-Präfix
  (gemessen, rc 2, „kein Lauf des Werkzeugs"). Zweitens ein Go-Panic oder
  `fatal error`: Exit 2, Stacktrace ohne Präfix, stdout leer. Im Code gibt es kein
  `recover()`. Ein Slice, der einen Absturz behebt oder einführt, bekommt deshalb
  statt `ABWEICHUNG` die Meldung „kein Lauf des Werkzeugs" mit einer
  Stacktrace-Zeile als Grund (simuliert, rc 2). Gerade die Verhaltensänderung, die
  die Probe finden soll, erscheint dann als Infrastruktur-Ausfall. Das ist laut,
  aber falsch diagnostiziert. Gemessen war das Kriterium über sieben Formen und
  über den Erfolgs- und Befund-Pfad, nicht über die Abbruch-Pfade des Werkzeugs.
- **verifizierbar:** ja. `REF=HEAD PROBE_FORMS=-h bash tools/blackbox-probe.sh`
  liefert rc 2. Für den Panic-Fall: Shim, der nur für `:latest` Exit 2 mit
  `panic:`-stderr liefert.
- **klasse:** echter-lauf-als-ausfall-abgewiesen

### R3-3 — INFO — Der spezifizierte Exit-0-Pfad „0 Datei(en) geprüft" ist die Restform des gleichen Scheiterns

- **quelle:** `harness/sensors/blackbox-probe.md` Grenze 4 (Restgrenze);
  Spezifikation §DC-FA-CLI-001.a („eine Wurzel ohne Markdown-Dateien, aber mit
  Inhalt, liefert „0 Datei(en) geprüft" und Exit 0")
- **pfad:** `harness/sensors/blackbox-probe.md` · „keine gemessene Form (Exit 125, Daemon-Ausfall
  mit 1, leerer Mount) tut das"
- **befund:** Sieht der Daemon unter dem Mount-Pfad ein nicht leeres Verzeichnis
  ohne Markdown, enden beide Seiten gleich mit Exit 0, und die Probe meldet
  „byte-identisch" (gemessen: 10× `gleich`, rc 0). Die Restgrenze deckt das dem
  Wortlaut nach ab („Ausgabe trüge"). Ihre Beispielliste nennt aber nur den
  leeren Mount, und gerade diese Form ist ein spezifizierter Erfolgspfad des
  Werkzeugs, kein Container-Ausfall. Ein realistischer Auslöser braucht einen
  Daemon, der am Client-Pfad anderen Inhalt sieht. Deshalb INFO.
- **verifizierbar:** ja. Shim, der die Mount-Quelle durch ein Verzeichnis mit einer
  einzelnen `.txt`-Datei ersetzt.
- **klasse:** stilles-gruen-bei-gleichem-scheitern

## Schließung R2

| R2 | Stand | Beleg |
|---|---|---|
| R2-1 | geschlossen | Shim „leerer Mount": rc 2, `NICHT PRUEFBAR`. Der Fall „beide Exit 2" zählt nicht mehr als gleich. |
| R2-2 | in der Sache geschlossen | Grenze 4 nennt den Fall „beide Exit 2" und die Restgrenze, die Exit-Tabelle ist nachgezogen. Die neue Begründung („nirgends geprüft") überdehnt aber selbst → R3-1. |
| R2-3 | geschlossen | `PROBE_FORMS="--repair --repair-broad --trace"`: 15× `gleich`, rc 0. |
| R2-4 | geschlossen | §Nebenwirkung nennt die Kommandozeile ausdrücklich. |

## Negativbefunde

- **Stilles Grün durch Docker-Ausfälle** (125, toter Daemon mit 1, leerer Mount):
  geprüft, ohne Befund. Alle drei enden mit rc 2.
- **Exit-Präzedenz** (`unvergleichbar` vor `abweichungen`): geprüft, ohne Befund.
  Ein Lauf mit beidem endet mit 2 und nennt beide Zahlen.
- **Exit-Tabelle gegen Code:** geprüft, ohne Befund. Die Zeile „0 … auf mindestens
  einer Seite prüfbar" ist schwächer als das Ist (bei Exit 0 waren beide Seiten
  prüfbar), stimmt also.
- **Sensor-Datei gegen Code, Grenze 4 (Kriterium):** geprüft, ohne Befund außer
  R3-1/R3-2. Bedingung „0, 1 oder 2 und stdout oder `d-check:`" ist wortgleich
  umgesetzt.
- **Kommentare** (§3.7): geprüft, ohne Befund. Beide neuen Kommentar-Blöcke tragen
  Zusage bzw. Grenze, keine Review-Historie, keine Befund-Nummern.
- **Slice-Plan, Plan-Änderung nach R2:** geprüft, ohne Befund. Sie liegt vor dem
  Code-Commit (`3520840f` vor `1b5d9b60`), nennt die gemessene Formen-Menge statt
  einer Allaussage und weitet die Abgrenzung nicht.
- **Docker/make-only** (§3.1): geprüft, ohne Befund. Das Skript ruft nur `git`,
  `docker`, `tar` und POSIX-Werkzeuge.

## Kategorie-Summary

HIGH 0 · MEDIUM 2 (R3-1, R3-2) · LOW 0 · INFO 1 (R3-3).
Wiederkehrende Finding-Klassen:

- `allaussage-ueber-gemessene-teilmenge`, zweites Auftreten in diesem Slice
  (R2-3, R3-1).
- `stilles-gruen-bei-gleichem-scheitern`, viertes Auftreten (F-1, V-1, R2-1,
  R3-3), jetzt nur noch als INFO-Restform.

## Verdikt

Kein stiller Grün-Pfad mehr bei realistischen Docker- oder Mount-Ausfällen. Die
Klasse aus R1, Verifikation und R2 ist geschlossen. Offen sind zwei MEDIUM,
beide in der Gegenrichtung (falsch rot mit falscher Diagnose), beide laut:

- Fehlerpfade sind als Vergleich ausgeschlossen, ohne dass das als Grenze
  dasteht.
- Abstürze des Werkzeugs gelten als „kein Lauf".

Nach Skill blockieren MEDIUM typischerweise. Weil das Target deklariert kein Gate
ist und keiner der beiden Befunde still grün meldet, ist die Einarbeitung vor
Closure ausreichend. Eine weitere Review-Runde ist nur nötig, wenn sich das
Klassifikations-Kriterium erneut ändert.

**Steering-Loop-Signal** (Skill §Kontext-Eskalation): Das Kriterium wurde in vier
Runden jedes Mal gegen die eben gemessene Fehlerform geschärft. Die Ausgangs-Matrix
oben (Exit × Lebenszeichen × Seite) ist die vollständige Fallmenge, und auch sie
wurde erst im Review aufgestellt. Ob eine solche Matrix in die Sensor-Datei
gehört, entscheidet der Implementer.
