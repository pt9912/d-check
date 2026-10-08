# Review R1 — slice-258: Black-Box-Vorher/Nachher-Probe als make-Target

- **Review-Art:** Code. Geprüft gegen den Slice-Plan (`slice-258`, §1 Ziel und
  Abgrenzung, §3 Plan samt Plan-Änderung zu `scan.ignore`, §6 Risiken) und die
  Hard Rules [`AGENTS.md`](../../AGENTS.md) §3.1, §3.6, §3.7, §3.8 sowie §5
  Regel 13. Die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-258` · Range `231ce771..87b2ada9` (`af3fc934`
  Plan-Änderung, `87b2ada9` Skript, Fixtures, Target, Sensor-Datei, Index-Zeile,
  `scan.ignore`). Mitgelesen, nicht im Diff: `Dockerfile` (Build-Args,
  `-buildvcs=false`), `Makefile` `build`/`clean`, `harness/sensors/doc-check.md`,
  `MR-069`.
- **Skill:** `reviewer.md` @ 1.18.0 (`3814e269`).
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-08
- **Eingangs-Kontext:** Slice-Plan `slice-258`; Hard Rules; keine ADR berührt
  (Plan: „Harness-Werkzeug, keine Spec-Aussage"). Vorherige Findings am selben
  Modul: keine — das Werkzeug ist neu; die Klasse *stilles Grün* ist im Register
  geführt.
- **Proben:**
  - `make blackbox-probe REF=HEAD` — 20 Vergleiche (5 Fixtures × 4 Formen), alle
    `gleich`, Schlusszeile `byte-identisch über 20 Vergleiche`, ca. 50 s. Die
    Exit-Codes sind plausibel: `ids`/`links`/`targets` exit 1, `sauber`/`repo` exit 0.
  - Fixtures einzeln gegen `d-check:latest` gefahren: Jedes Fixture löst den
    vorgesehenen Befund aus (`id-unlinked`, `target-missing` + `anchor-missing`,
    `gate-undocumented` + `gate-phantom`, sauber 0). Keines endet in einem
    Konfigurationsfehler, der beide Seiten gleich leer laufen ließe.
  - **Bruch-Test Docker-Fehler:** `docker`-Shim (bash) im `PATH`, der für `run` die
    Zeile `docker: Error response from daemon: simulierter Mount-Fehler.` und
    exit 125 liefert und alles andere durchreicht. Ergebnis: 20× `gleich … (exit 125)`,
    `byte-identisch über 20 Vergleiche`, **Skript-rc 0**.
  - Fehlerformen von `docker run` direkt gemessen: Daemon nicht erreichbar
    (`DOCKER_HOST=unix:///nonexistent.sock`) → rc **1**. Fehlender Entrypoint → rc 127.
  - `docker run … --doctor`/`--yaml` auf `sauber`: Keine Ausgabe trägt die
    eingebettete Version.
  - Arbeitsbaum am Ende sauber, kein Code geändert (der Shim lag im Scratchpad).

## Findings

### F-1 — MEDIUM — Gleich gescheiterte Läufe zählen als „byte-identisch"

- **quelle:** Skill Prüffragen 1 und 19; `harness/sensors/blackbox-probe.md`
  §Ausgabe und Ausgänge (Exit 2 = „Lauf gescheitert")
- **pfad:** `tools/blackbox-probe.sh` · „`echo "$rc" > "$4.rc"`" und
  „`if [ -z "$diffs" ]; then`"
- **befund:** `lauf()` hält den Exit von `docker run` fest, ohne zu unterscheiden, ob
  das Werkzeug lief oder der Container gar nicht startete. Scheitert der Start auf
  beiden Seiten gleich, sind stderr und Exit identisch, und der Vergleich zählt ihn
  als `gleich`. Das Skript endet mit „byte-identisch über N Vergleiche" und rc 0.
  Gemessen: 20× `exit 125`, Skript-rc 0. Fällt der Daemon während des Laufs aus, ist
  der rc 1 — derselbe Code wie ein Befund des Werkzeugs. Selbst in der Einzelzeile
  ist der Fall dann nicht mehr erkennbar; nur bei `sauber`/`repo` steht exit 1 statt 0
  da, und auch das als `gleich`. Der Vertrag sagt für einen gescheiterten Lauf Exit 2
  zu. Wer die Probe als Beleg für „ohne Schalter unverändert" ins Review gibt, legt
  dann einen Beleg ohne einen einzigen gelaufenen Vergleich vor. MEDIUM statt HIGH:
  Das Target ist deklariert kein Gate. Mit einem Bindepunkt in `gates` oder einer
  DoD-Pflicht stiege der Befund auf HIGH.
- **verifizierbar:** ja — der `docker`-Shim oben, oder `DOCKER_HOST` auf einen toten
  Socket nach dem Image-Inspect.
- **klasse:** stilles-gruen-bei-gleichem-scheitern

### F-2 — MEDIUM — Die Grenzen-Liste der Sensor-Datei nennt die größte Lücke nicht

- **quelle:** Skill Prüffrage 18; [`AGENTS.md`](../../AGENTS.md) §5 Regel 13
- **pfad:** `harness/sensors/blackbox-probe.md` · „| 0 | byte-identisch über alle Vergleiche |"
- **befund:** Dreht man die Zusage „stdout, stderr und Exit getrennt" um, folgt eine
  Grenze: Verglichen wird das Ergebnis des Container-Aufrufs, nicht das des
  Werkzeugs. Ein Start- oder Daemon-Fehler, der beide Seiten gleich trifft, ist
  „gleich". Diese Grenze steht nur im Code (F-1), nicht unter §Grenze, und die
  Exit-Tabelle sagt das Gegenteil (0 = byte-identisch, 2 = gescheitert). Die drei
  genannten Grenzen (Fixture-Ausschnitt, falsche Erwartung, kein Gate) sind kleiner
  als diese. Ein Leser der Datei hält ein rc 0 deshalb für einen gelaufenen
  Vergleich.
- **verifizierbar:** ja — wie F-1.
- **klasse:** grenzen-liste-ohne-groesste-luecke

### F-3 — LOW — Ein drittes Verkleinerungs-Ventil von `doc-check` bleibt in dessen Grenzen-Liste ungenannt

- **quelle:** `harness/sensors/doc-check.md` §Grenze Punkt 4; [`AGENTS.md`](../../AGENTS.md) §5 Regel 13
- **pfad:** `.d-check.yml` · „`"tools/blackbox-probe/fixtures/**"]`"
- **befund:** Die bisherigen `scan.ignore`-Einträge nehmen **Fremdinhalt** aus (der
  Kommentar sagt „Selbe Klasse wie die eingebauten SKIP_DIRS (Fremd-/Generiertes)";
  belegt über MR-017/MR-019). Der neue Eintrag nimmt erstmals **repo-eigenen** Inhalt
  aus. `harness/sensors/doc-check.md` führt unter Punkt 4 den verkleinerten
  Prüfbereich als „zwei Ventile" (`ignore-refs`, `d-check:ignore`), und die
  Vertragszeile sagt „gesamte Repo-Doku". Dass der Scan selbst repo-eigene Pfade
  auslässt, nennt die Datei nicht. Zu §3.6 siehe die Negativbefunde: keine ADR-Pflicht.
- **verifizierbar:** nein — Urteil über die Dokumentation.
- **klasse:** grenzen-liste-ohne-groesste-luecke

### F-4 — INFO — Plan sagt „dieselben Build-Args wie `make build`", das Skript nimmt die Dockerfile-Defaults des REF-Stands

- **quelle:** Slice-Plan `slice-258` §1 („dieselben Build-Args wie `make build`")
- **pfad:** `tools/blackbox-probe.sh` · „`docker build -q --build-arg VERSION=0.0.0-dev --target runtime`"
- **befund:** `make build` reicht `GO_VERSION`/`GOLANGCI_LINT_VERSION` aus dem
  Makefile durch. Der Vorher-Build verwendet die `ARG`-Defaults des Dockerfiles von
  REF. Heute wirkt sich das nicht aus: Beide `FROM`-Zeilen sind digest-gepinnt, und
  `GOFLAGS` setzt `-buildvcs=false`, also kann auch der `.git`-lose `git archive`-Kontext
  nichts einbetten. Die Sensor-Datei beschreibt den tatsächlichen Stand richtig
  („Build-Defaults dieses Stands"). Abweichend formuliert ist nur der Plan.
  `VERSION` ist auf der Vorher-Seite fest auf `0.0.0-dev` gesetzt, auf der
  Nachher-Seite aber `VERSION ?=` aus der Umgebung. Der Kommentar „gleich dem
  Nachher-Default" hält deshalb nur ohne exportiertes `VERSION`. Die Standard-Formen
  geben die Version nicht aus (gemessen), wohl aber eine `PROBE_FORMS` mit `--print-mk`.
- **verifizierbar:** ja — `VERSION=9.9.9 make blackbox-probe REF=HEAD PROBE_FORMS=--print-mk`.
- **klasse:** plan-formulierung-weicht-vom-mechanismus

### F-5 — INFO — Eingaben kommen aus dem Arbeitsbaum, nicht aus REF

- **quelle:** [`AGENTS.md`](../../AGENTS.md) §3.8 (als Frage auf das Werkzeug angewandt)
- **pfad:** `tools/blackbox-probe.sh` · „`fixtures+=("repo=$PWD")`"
- **befund:** Beide Images lesen dieselben Fixtures und denselben Arbeitsbaum samt
  dessen `.d-check.yml`. Verglichen werden also zwei Binaries auf den **heutigen**
  Eingaben, nicht zwei Stände. Ändert ein Slice Code **und** Konfiguration, liest das
  Vorher-Binary die neue Konfiguration. Nutzt sie einen neuen Schlüssel, endet die
  Vorher-Seite mit Exit 2, und das erscheint als Abweichung, ohne
  Verhaltensänderung bei unveränderter Konfiguration. Für die Zusage „ohne Schalter
  unverändert" ist das gewollt. In der Sensor-Datei steht es nur implizit
  („dieselben Fixtures").
- **verifizierbar:** nein.
- **klasse:** eingaben-nicht-aus-ref

### F-6 — INFO — `d-check:probe-vorher` überlebt `make clean`

- **quelle:** Maintainability
- **pfad:** `Makefile` · „`clean: ## Lokale Images entfernen.`"
- **befund:** Das Skript lässt den Tag `$(IMAGE):probe-vorher` stehen, und `clean`
  zählt ihn nicht mit. Das ist dieselbe Lage wie beim bereits vorhandenen `:arm64`.
- **verifizierbar:** ja — `docker images d-check` nach `make clean`.
- **klasse:** clean-zaehlt-nicht-alle-tags

## Negativbefunde

- **Leerer Fixture-Glob** (`for d in "$FIXTURES_DIR"/*/`): Ohne Treffer bleibt das
  Literal stehen, `[ -d ]` scheitert, `fixtures` ist leer, und `fail` beendet mit
  rc 2, bevor `repo` angehängt wird. Ein relativer Aufruf außerhalb der Wurzel endet
  ebenso fail-closed. Geprüft, ohne Befund.
- **Leere Formen:** `PROBE_FORMS` aus Leerraum führt zu `fail`. Geprüft, ohne Befund.
- **Leere Ausgaben** sind als solche kein stilles Grün: Jedes Fixture erzeugt
  nicht-leeres stderr oder stdout (gemessen). Das Restrisiko fällt unter F-1.
- **Neuer Leseweg `git archive | tar`:** `set -o pipefail` ohne `errexit`, das `||`
  hängt an der Pipeline, und ein Fehler auf einer der beiden Seiten führt zu `fail`.
  **`docker build -q > /dev/null`:** Verworfen wird nur die Image-ID auf stdout, stderr
  bleibt sichtbar, ein rc ≠ 0 führt zu `fail`. Ein veraltetes `:probe-vorher` wird
  nicht verwendet, weil das Skript beim Build-Fehler endet. Geprüft, ohne Befund.
  Den neuen Leseweg `docker run` behandelt F-1.
- **Vorher-Image = Stand von REF:** `git archive` liefert den committeten Baum von REF.
  Es gibt kein `.gitattributes` (kein `export-ignore`/`export-subst`) und kein
  `.dockerignore`. `-buildvcs=false` steht im Dockerfile. Der Nachher-Build aus dem
  Arbeitsbaum mit `.git` bettet deshalb ebenfalls nichts ein. Geprüft, ohne Befund
  (siehe F-4 zu den Args).
- **Rauschen:** `REF=HEAD` gegen den Nachher-Build desselben Stands ist
  byte-identisch über 20 Vergleiche. Zeit- und Pfad-Rauschen trat nicht auf.
  `--network none` ist auf beiden Seiten gleich gesetzt.
- **Kommentar-Klassen §3.7:** Kopf von `tools/blackbox-probe.sh`
  (ZUSAGE/ABGRENZUNG/GRENZE), die Inline-Kommentare (`-`-Form, nonroot-Kopie,
  VERSION-Kopplung), der Makefile-Kommentar und der neue `.d-check.yml`-Kommentar
  tragen Zusage, Abgrenzung oder Kopplung. Sie enthalten keine Slice-Nummern, keine
  Befund-Marker und keine Review-Historie. Die Klammer „(zusammengefuehrte Streams …
  zeigen Abweichungen, die keine sind)" begründet die Zusage „getrennt" am Gegenstand
  und ist keine Deliberation über eine verworfene Variante. Geprüft, ohne Befund.
- **§3.6 und `scan.ignore`:** keine ADR-Pflicht. §3.6 betrifft das Senken einer
  Schwelle oder Prüfregel. Hier entfällt keine Prüfregel, und nichts zuvor Geprüfte
  fällt heraus: Der Pfad existiert in `231ce771` nicht. `MR-069` gilt laut
  Geltungsbereich nur für `ignore-refs` auf entfernte Baseline-Bäume. Die Begründung
  steht im Plan, und zwar vor dem Code-Commit (`af3fc934` vor `87b2ada9`). Den
  Dokumentationsrest behandelt F-3.
- **Plan-Abgrenzung §1:** kein Gate-Bindepunkt (nicht in `gates`/`ci`/Nachtlauf),
  keine Workflow-Pflicht, keine Lesart-Prüfung. Eingehalten.
- **Index-Zeile `harness/README.md`:** steht in der Werkzeug-Tabelle, `kein Gate` in
  der Bindungsspalte, mit Halbsatz zur Funktion. Sie entspricht dem dortigen
  Kriterium „urteilt nicht über den Repo-Zustand". Geprüft, ohne Befund.
- **Sensor-Datei gegen den Code (Regel 13), außer F-1/F-2:** Formen, Fixture-Menge,
  Exit-2-Gründe und der Hinweis „make normalisiert auf 2" stimmen mit dem Skript
  überein. Geprüft, ohne Befund.
- **§3.1:** nur `bash`, `git`, `docker`, `tar`, `cmp`, `diff`, `sed`, `head`, kein
  Host-Interpreter. Geprüft, ohne Befund.
- **Netz:** Die Läufe verwenden `--network none`. Der Vorher-Build braucht Netz nur
  für einen ungecachten Basis-Layer, so wie `make build`. Geprüft, ohne Befund.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 2 | 1 | 3 |

Wiederkehrende Klassen: `grenzen-liste-ohne-groesste-luecke` (F-2, F-3; Register
`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`, der Anlass dieses Slice) und
stilles Grün (F-1).

## Verdikt

**Nicht freigegeben.** F-1 und F-2 blockieren. Es ist ein Befund mit zwei Orten: Das
Werkzeug, das das stille Grün der Hand-Messung ersetzen soll, meldet bei gleichem
Startfehler selbst still grün, und seine Grenzen-Liste nennt das nicht. F-3 bis F-6
blockieren nicht.
