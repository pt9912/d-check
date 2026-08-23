# Review-Report: slice-135 — 2026-08-23

**Review-Art:** Code — geprüft gegen Slice-Plan, `AGENTS.md` §3.9/§3.1/§3.8/§4,
Zensus (slice-132), Sensors-Konvention (`harness/README.md`), `.d-check.yml`
(Block `targets`).

**Gegenstand:** Commit `b96de42` (Diff `HEAD~1..HEAD`) — Feature-Commit von
slice-135 (`AGENTS.md`, `Makefile`, `harness/README.md`,
`tools/harness/workflow-pins.sh`, neu).

**Skill:** `.harness/skills/reviewer.md` @ Version 1.10.0
**Modell:** claude-sonnet-5 · **Datum:** 2026-08-23

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan [`slice-135-uses-pin-sensor.md`](../plan/planning/done/slice-135-uses-pin-sensor.md)
- Welle [`welle-84-durchsetzung.md`](../plan/planning/done/welle-84-durchsetzung.md)
- Zensus [`slice-132-hard-rule-zensus.md`](../plan/planning/done/slice-132-hard-rule-zensus.md) (§3.9-Zeile)
- `tools/harness/workflow-pins.sh` (vollständig, Quelltext + isolierte Shell-Proben)
- `.github/workflows/{ci,release,upstream-drift}.yml` (vollständig)
- `Makefile` (Target, `.PHONY`, `gates`-Kette)
- `AGENTS.md` §3.1, §3.8, §3.9, §4
- `harness/README.md` §Sensors
- `.d-check.yml` Block `targets`
- Beobachtungs-Register [`observations.md`](../plan/planning/observations.md) BEO-010/BEO-011/BEO-012

---

## Findings

### F-1 — `.yaml`-Workflows sind für den Sensor unsichtbar, obwohl §3.9 sie einschließt

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.9
- `pfad`: `tools/harness/workflow-pins.sh:26,49` (`dir="${dir}"/*.yml`); `AGENTS.md:261,268`
- `befund`: §3.9 bindet „jeder `uses:`-Eintrag in `.github/workflows/`" ohne
  Endungs-Einschränkung, und die neue Zeile „Durchgesetzt: `make workflow-pins`
  in `make gates`" beansprucht dieselbe ungequalifizierte Reichweite. Der
  Sensor liest nur `"${dir}"/*.yml` — GitHub akzeptiert für Workflow-Dateien
  gleichrangig die Endung `.yaml`. Eine `.github/workflows/foo.yaml` mit einem
  beweglichen Tag bliebe vollständig ungeprüft; `checked` zählt weiter über 0
  (die drei bestehenden `.yml`-Dateien tragen bereits Treffer), also greift
  auch der Fail-closed-Zweig nicht. `make workflow-pins` und `make gates`
  blieben grün. Die Doku nennt „drei Grenzen, alle drei benannt" (Form,
  Prosa-Filter, leere Menge) — die Endungs-Lücke ist eine vierte, unbenannte.
- `verifizierbar`: ja — eine `.github/workflows/foo.yaml` mit
  `uses: actions/checkout@v6` anlegen; `make workflow-pins` bleibt Exit 0.
- `klasse`: gate-scope-glob-endungslücke

### F-2 — Die `file:line`-Rekonstruktion ist ein Zufallstreffer, kein korrekter Bau

- `kategorie`: MEDIUM
- `quelle`: Maintainability / Messmethode des Sensors
- `pfad`: `tools/harness/workflow-pins.sh:49`
- `befund`: `grep -n … "${dir}"/*.yml | sed "s|^|${dir}/|; s|^${dir}/${dir}/|${dir}/|"`
  verlässt sich darauf, dass GNU-grep bei **mehreren** File-Argumenten den
  Dateinamen voranstellt (dann entfernt der zweite `sed`-Ausdruck das doppelte
  Präfix wieder). Bei **genau einer** matchenden `.yml`-Datei lässt grep den
  Dateinamen weg — der erste `sed` hängt `${dir}/` an, der zweite greift nicht
  (kein doppeltes Präfix), und die spätere Feld-Zerlegung
  (`file="${hit%%:*}"`) reißt Verzeichnis und Zeilennummer zusammen. Isoliert
  nachgestellt (drei-Datei- vs. Ein-Datei-Verzeichnis, echtes GNU grep 3.11):
  bei einer Datei wird aus `uses: actions/checkout@…` die Fundstelle
  `…/one/6` statt `…/one/only.yml:6`, `line` wird zu `"        uses"` statt
  `6`. Erkennung und Exit-Code bleiben korrekt (das `text`-Feld ist
  unversehrt), nur die gemeldete Fundstelle ist falsch. Heute nicht sichtbar
  (drei Workflow-Dateien mit Treffern in allen dreien), aber latent: sinkt die
  Zahl der `.yml`-Dateien mit `uses:`-Treffern je auf eine, kippt die
  Ausgabe.
- `verifizierbar`: ja — zwei der drei Workflow-Dateien temporär entfernen/leeren
  und `make workflow-pins` gegen die verbleibende mit einem `uses:`-Treffer
  laufen lassen; die gemeldete `Datei:Zeile` ist dann falsch.
- `klasse`: grep-filename-omission-single-file

### F-3 — Die Form-Prüfung erkennt nur Git-Commit-SHA-Pins, keine gleichwertigen Alt-Formen

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `tools/harness/workflow-pins.sh:40,45`
- `befund`: `@[0-9a-f]{40}([[:space:]]|$)` setzt einen unquotierten, direkt auf
  `@` folgenden 40-Hex-String voraus. Ein zitierter Wert
  (`uses: "actions/checkout@<sha>" # tag`) oder eine Docker-Digest-Referenz
  (`uses: docker://image@sha256:<64-hex>`) — beides von GitHub Actions
  akzeptierte, bereits gepinnte Formen — würden fälschlich als
  `uses-pin-missing` gemeldet. Fail-safe-Richtung (Falsch-Positiv, kein
  Silent-Pass); im heutigen Bestand kommt keine der beiden Formen vor.
- `verifizierbar`: ja — testweise eine der beiden Formen in eine Workflow-Datei
  einfügen und `make workflow-pins` laufen lassen.
- `klasse`: sha-form-nur-git-pin

## Negativbefunde

- geprüft, ohne Befund: Prosa-Filterung — die drei Kopfzeilen-Kommentare, die
  `` `uses:` `` in den drei Workflow-Dateien erklärend nennen, beginnen mit
  `#` und werden vom verankerten `^[[:space:]]*(- )?uses:[[:space:]]` korrekt
  nicht getroffen (Anker auf Zeilenanfang, nicht Substring-Suche).
- geprüft, ohne Befund: `jobs.<id>.uses:` (reusable-workflow-Aufruf) — die
  Regex ist einrückungs-/kontext-unabhängig und würde diese Form fangen; kein
  blinder Fleck (anders als zunächst vermutet).
- geprüft, ohne Befund: Composite-Actions unter `.github/actions/**` — liegen
  außerhalb der Regel selbst (§3.9 bindet ausdrücklich nur
  `.github/workflows/`), keine Überdehnung des Sensors gegenüber der Regel.
- geprüft, ohne Befund: eingebetteter Doppelpunkt im Treffer-Text (z. B. im
  Tag-Kommentar) — bricht die Feld-Zerlegung nicht, `text` wird über den
  zweiten Doppelpunkt hinaus vollständig übernommen und nur per `grep -qE`
  geprüft, nicht weiter zerlegt.
- geprüft, ohne Befund: Dateiname mit Leerzeichen im Mehrdatei-Fall — isoliert
  nachgestellt (zwei Dateien, eine mit Leerzeichen im Namen); die
  Colon-basierte Zerlegung bleibt korrekt. Das eigentliche Risiko ist die
  Dateizahl (F-2), nicht das Leerzeichen.
- geprüft, ohne Befund: `set -euo pipefail` + `while read … done < <(…)` — ein
  `grep`-Exit ≠ 0 innerhalb der Prozess-Substitution (keine Treffer, leeres
  Verzeichnis, fehlende Dateien) bricht das Skript **nicht** ab (Exit-Status
  einer Prozess-Substitution wird von `set -e` nicht verfolgt); isoliert
  nachgestellt für alle drei Fail-closed-Fälle aus §5 der Aufgabe
  (Verzeichnis fehlt / leer / Dateien ohne `uses:`) — alle drei erreichen
  korrekt den `checked -eq 0`-Zweig und schließen mit Exit 1.
- geprüft, ohne Befund: §3.1 Docker/make-only — das Skript ruft nur `git`,
  `grep`, `sed`, `printf`/`echo` und Bash-Builtins; `sed` ist bereits Teil der
  Werkzeug-Klasse (Präzedenz: `tools/harness/fetch-baseline-cache.sh`, selbst
  Teil von `gates`). §3.1 nennt seine Liste ausdrücklich als Klasse, nicht als
  Aufzählung — kein Verstoß.
- geprüft, ohne Befund (BEO-011/BEO-012): „4 uses:-Einträge, alle formgerecht"
  und „zehn Glieder" (`gates`-Kette) sind nachgezählt und stimmen exakt (1
  `uses:` in `ci.yml`, 2 in `release.yml`, 1 in `upstream-drift.yml`; je 40
  Hex-Zeichen + Tag-Kommentar; die `gates:`-Prerequisite-Liste trägt genau
  zehn benannte Glieder vor `record-gates`) — keine Bestandsaussage
  reicht über das Gemessene hinaus.
- geprüft, ohne Befund: Zensus-Eintrag `slice-132` §9 (Closure-Notiz) bleibt in
  diesem Commit unangetastet. Das ist korrekt: die Zensus-Zeile zu §3.9 war
  zum Zeitpunkt der Zensus-Messung wahr (einseitig) und bleibt es als
  historischer Beleg eines `done/`-Slice — sie widerspricht dem Verzeichnis
  nicht. Die Auflösung gehört in `welle-84-results.md` (Wellen-Closure) bzw.
  `slice-135` §9, beide noch offen; kein Nachzieh-Fehler in diesem Commit.
- geprüft, ohne Befund: Doku-Drei-Fach-Konsistenz (`BEO-010`-Klasse) —
  `Makefile` (Target + `.PHONY`-Eintrag + `gates`-Kette + Echo-Zeile),
  `AGENTS.md` §4-Tabelle und `harness/README.md` Sensors-Tabelle nennen
  dieselben zwei Befund-Codes (`uses-pin-missing`/`uses-pin-untagged`) und
  dieselben drei genannten Grenzen wortgleich in der Substanz. `.d-check.yml`
  Block `targets` (`makefiles: [Makefile]`, `doc-tables: [AGENTS.md,
  harness/README.md]`) leitet die Konsistenzprüfung dynamisch aus genau
  diesen drei Dateien ab — keine zusätzliche manuelle Registrierung nötig,
  keine Lücke wie in BEO-010 beobachtet.
- geprüft, ohne Befund: Ort-Entscheidung (make-Target vs. CI-Schritt vs.
  d-check-Modul) — die Ablehnung des Modul-Wegs verweist korrekt auf §3.8
  (Scan- vs. Lese-Menge) als ADR-Frage statt sie im Slice selbst zu
  entscheiden; deckt sich mit der bereits in `welle-84-durchsetzung.md`
  etablierten Messung („d-check scannt nur Markdown, 467 = 467").
- geprüft, ohne Befund: „Ausdrücklich NICHT" aus dem Slice-Plan (§3) — keine
  SHA-Existenz-/Gültigkeitsprüfung, kein Auto-Pinning, keine Ausweitung auf
  andere YAML-Felder im Skript vorhanden; deckt sich mit der Zusage.
- geprüft, ohne Befund: bewusstes Brechen (drei Formen aus der Commit-Botschaft)
  — Logik am Quelltext nachvollzogen und mit isolierten Shell-Proben
  bestätigt: beweglicher Tag ⇒ `uses-pin-missing`, SHA ohne Tag-Kommentar ⇒
  `uses-pin-untagged`, geleertes Verzeichnis ⇒ Fail-closed Exit 1,
  unveränderter Bestand ⇒ Exit 0 (`checked=4, findings=0` gegen eine
  Kopie der drei echten Workflow-Dateien reproduziert).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** gate-scope-glob-endungslücke ·
grep-filename-omission-single-file · sha-form-nur-git-pin

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) ist ein Stilles-Grün-Pfad in einem
Gate, das sich selbst gerade als „Durchgesetzt" (§3.9) einträgt: eine
`.yaml`-Workflow-Datei mit einem ungepinnten `uses:`-Eintrag wird von
`make workflow-pins`/`make gates` nicht erkannt, während die Doku exakt diese
Reichweite (jeder Workflow, keine Endungs-Einschränkung) beansprucht und ihre
Grenzen als vollständig benannt ausweist ("drei Grenzen, alle drei benannt").
F-2 (MEDIUM) ist kein Silent-Pass — Erkennung und Exit-Code bleiben korrekt —,
aber die `sed`-Rekonstruktion der Fundstelle ist nachweislich kein robuster
Bau, sondern funktioniert nur, solange mindestens zwei `.yml`-Dateien
Treffer beitragen; das sollte vor Merge geklärt werden, weil sonst künftige
Fundstellen-Ausgaben irreführen. F-3 (LOW) blockiert nicht.

**Übergabe:** Findings gehen an den Implementer. Empfehlenswerte Reichweite
für den Fix von F-1: entweder den Glob auf `*.yml`/`*.yaml` erweitern, oder
die Endungs-Einschränkung explizit als vierte benannte Grenze in Skript-Kopf,
`AGENTS.md` §3.9 und der Sensors-Tabelle ausweisen (dann wäre „Durchgesetzt"
ehrlich geschnitten statt überdehnt) — welche der beiden Optionen gewählt
wird, ist Implementer-Entscheidung. Dieser Report ersetzt keine Verifikation;
DoD-/Gate-Lauf-Bestätigung (`make gates` Exit, DoD-Checkboxen) prüft der
Verifier separat.
