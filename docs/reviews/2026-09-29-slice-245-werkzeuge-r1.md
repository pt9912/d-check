# Review-Report — slice-245, Runde 1

- **Review-Art:** Code (geprüft wird der Diff gegen Plan, Entscheidungen und Hard Rules — nicht gegen die DoD; das ist die Verifikation)
- **Gegenstand:** slice-245 — `tools/harness`-Werkzeuge aus ai-harness-init evaluieren (welle-91); Range `e6c51c14..HEAD`, Commits `aad1c9e7`, `a0d53819`, `bff7b762`, `8637f08a`
- **Skill:** `reviewer.md` @ 1.16.0
- **Modell-ID:** glm-5.3-flash
- **Datum:** 2026-09-29
- **Eingangs-Kontext:** slice-245-Plan (`in-progress/`), slice-247-Plan (`open/`), `AGENTS.md` §3 (§3.1, §3.3, §3.7, §3.9), `harness/conventions.md` (`MR-045`, `MR-053`, `MR-054`, `MR-040/042/044`), `harness/README.md` §Sensors/§Werkzeuge, `v6.13.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form, `v6.13.0` · `regelwerk/modul-13-quality-gates.md` §Hard Rule und §Die dritte Lage. Empirische Gegenproben: scratch-Repos für slice-mv (drei Fälle) und history-range-guard (Exit-Codes); ci.yml-Checkout (`fetch-depth: 0`).

---

## Befunde

### F-1 — MEDIUM · slice-mv: Outgoing-Verweise werden bei gesetztem `SLICE_MV_DONE_UNTERORDNER` in falscher Tiefe umgeschrieben

- **quelle:** Maintainability — Zusage-Bruch des Werkzeugs (Skriptkopf ZUSAGE, `slice-mv.sh:5-15`) plus Grenzen-Liste ohne die durch die eigene Anpassung entstandene Lücke (`slice-mv.sh:41-64`)
- **pfad:** `tools/harness/slice-mv.sh:178` (Ersetzungsform `../$from/`), `:204-207` (d-check-Anpassung `ziel`), `:284` (Aufruf); Anpassung aus `a0d53819`
- **befund:** Die ausgehende Ersetzung hängt praefixlosen Zielen `../$from/` vor — eine Ebene. Mit der d-check-Anpassung liegt die bewegte Datei aber in `done/welle-<NN>/` (zwei Ebenen tief); der geschriebene Link `../in-progress/X.md` löst dort zu `done/in-progress/X.md` auf, das nicht existiert. Die eingehende Ersetzung nutzt korrekt `$ziel` mit Unterordner, die ausgehende nicht; die Grenzen-Liste im Skriptkopf deckt den Fall nicht. Empirisch belegt (scratch-Repo, `SLICE_MV_DONE_UNTERORDNER=welle-91`, Move `in-progress → done`): die bewegte Datei trägt `](../in-progress/slice-244-nachbar.md)` statt `](../../in-progress/slice-244-nachbar.md)`. Der danach entstehende Commit 2 behauptet in seiner Botschaft „Verweise … nachgezogen". Ohne Unterordner (wellenlos, flaches `done/`) ist das Verhalten korrekt — der Bug feuert genau im adaptierten Fall.
- **verifizierbar:** ja — `make doc-check` nach einem `make slice-mv` mit `SLICE_MV_DONE_UNTERORDNER` (toter Link); oder der scratch-Repro oben.
- **klasse:** `adaptions-luecke-relativer-pfad` (Anpassung ändert die Pfadtiefe, nur eine der beiden Ersetzungsrichtungen zieht mit)

### F-2 — MEDIUM · harness/README.md trägt drei Slice-Verweise (`MR-045`)

- **quelle:** `MR-045` (Geltungsbereich: `AGENTS.md` und `harness/README.md` — „Beide Dateien tragen **keine** `slice-<NNN>`-Verweise")
- **pfad:** `harness/README.md:127-129`
- **befund:** Die drei neuen Werkzeuge-Zeilen nennen je „adoptiert aus ai-harness-init (slice-245)". `MR-045` verlangt für diesen Ort Regel/Grenze/ADR-/MR-Zeiger statt Planungs-Artefakt-Bezug; die Begründung des Eintrags (Alterung des Verweises beim Lifecycle-Wandel) trifft auf den nackten Token genauso zu.
- **verifizierbar:** ja — `grep -n "slice-2" harness/README.md`; kein Gate, das die `MR-045`-Hälfte prüft.
- **klasse:** `slice-verweis-im-einstieg`

### F-3 — MEDIUM · Slice-Nummern in Code- und Makefile-Kommentaren (§3.7)

- **quelle:** `AGENTS.md` §3.7 („keine Slice-Nummern"; Herkunft nur als `DC-*`/`ADR-*`/`MR-*`/`seit welle-<NN>`-Feld) — dieselbe Klasse wie die zweifach früher gemeldete Instanz
- **pfad:** `tools/harness/slice-mv.sh:99` · `tools/harness/selbstpruefung.sh:3` · `Makefile:404,408,411,414`
- **befund:** Sechs Kommentare tragen `(slice-245)` — eine Slice-Nummer, die keine der zulässigen Herkunfts-Formen ist (der Baseline-Herkunfts-Anker für Make-Targets wäre die Form `seit slice-<Kennung>`, wie sie `Makefile:313` für slice-215 bereits benutzt; Adoptions-Provenance trägt die Commit-Botschaft von `aad1c9e7`, die sie schon führt). Die Kommentare tragen im Übrigen reguläre Klassen (Zusage/Grenze) — es ist das verbotene Zusatzelement, nicht ein fehlender Inhalt.
- **verifizierbar:** ja — `grep -n "slice-245" tools/harness/*.sh Makefile`; kein Gate.
- **klasse:** `slice-id-in-kommentar` — **Steering-Loop-Hinweis:** dritte Instanz dieser Klasse über Sessions hinweg (zwei frühere Korrekturen); ein nachziehender Guide/Sensor verdient ein Register-Ereignis.

### F-4 — LOW · slice-mv Identity-Fallback: leeres Array unter `set -u` bricht auf bash < 4.4

- **quelle:** Maintainability — latente Portabilitäts-Falle im Widerspruch zur eigenen Zielsetzung des Skripts
- **pfad:** `tools/harness/slice-mv.sh:239` (`local -a ident=()`) gegenüber `:247` und `:292` (`git "${ident[@]}" commit`)
- **befund:** Auf bash < 4.4 (macOS liefert bash 3.2 aus) bricht die Expansion `"${ident[@]}"` eines leeren Arrays unter `set -euo pipefail` mit „unbound variable" ab — genau im Normalfall (konfigurierte Identität, Array bleibt leer), und zwar erst nach dem `git mv`, so dass ein gestagter Rename liegen bleibt, dessen Gefahr der Skriptkopf selbst benennt. Das Skript investiert an anderer Stelle ausdrücklich in macOS-Portabilität (`psed_i`, BSD-sed-Kommentar `:70-79`); lokal (bash 5.2) läuft der Fallback korrekt — gemessen.
- **verifizierbar:** ja — bash ≤ 4.3 bzw. macOS; auf bash ≥ 4.4 reproduziert der Befund nicht.
- **klasse:** `bash-versionen-leeres-array`

### F-5 — LOW · history-range-guard: GRENZE-Zeile unterschlägt, dass der Wächter unauflösbare Basen selbst behandelt

- **quelle:** `AGENTS.md` §3.7 (ein Kommentar beschreibt, was da ist)
- **pfad:** `tools/harness/history-range-guard.sh:10-13` (GRENZE) gegenüber `:89-93` (Behandlung)
- **befund:** Die GRENZE-Zeile führt eine unauflösbare Basis unter „was dieser Wächter NICHT deckt" und begründet das mit „bricht schon ohne ihn ab". Der Wächter behandelt den Fall aber selbst — `git rev-list --count … || exit 2` mit eigener Meldung; gemessen: `bash tools/harness/history-range-guard.sh 0000000..HEAD` → Exit 2. Der Leser der GRENZE hält den Wächter für passiv in dem Fall, in dem er aktiv abbricht; die Substanz (das Modul bricht ebenfalls laut, Probe C) stimmt, die Abgrenzungs-Formulierung zum Code nicht.
- **verifizierbar:** ja — der Lauf oben; Gegenlesung `:89-93`.
- **klasse:** `grenze-gegen-code-drift`

### F-6 — LOW · Sensor-Doku zu `make trace-check`/`make adr-check` nennt den neuen Vorlauf-Wächter nicht

- **quelle:** Maintainability — Doku-Drift am Target-Vertrag
- **pfad:** `harness/sensors/trace-check.md`, `harness/sensors/adr-check.md` gegenüber `Makefile:368` und `Makefile:382`
- **befund:** Beide Targets führen seit `aad1c9e7`/`8637f08a` einen Host-Schritt vor dem Container-Lauf aus, der sie bei leerer, auflösbarer Range in shallow-Klonen laut abbrechen lässt — ein neues, benanntes Versagensbild, das in die sonst sehr detaillierten Grenze-Sektionen der beiden Sensor-Dokumente (Pack-Verhalten, Geschichte-Ausnahme) nicht eingezogen ist. Die README-Werkzeuge-Zeile beschreibt den Wächter als Einzellauf, nicht seine Verdrahtung in die beiden Gates.
- **verifizierbar:** ja — Gegenlesung der beiden Dokumente gegen die Makefile-Recipes.
- **klasse:** `sensor-doku-hinter-recipe`

### F-7 — LOW · `make slice-mv` fährt ohne Probe

- **quelle:** Negativtest-Pflicht zu neuem Vertrag (Prüffrage 13, hier auf Werkzeug-Ebene gewichtet)
- **pfad:** `tools/harness/slice-mv.sh:306-310` (BASH_SOURCE-Wächter „ein Test sourced dieses Skript …"), Makefile-Umfeld (kein slice-mv-Probe-Target)
- **befund:** Das Werkzeug, das künftig alle Lifecycle-Moves dieses Repos trägt, wurde ohne Probe adoptiert — anders als die Adoptions-Pendants `make guard-probe`/`make baseline-probe`, die ihre Prüfgegenstände gegen Proben fahren. Der BASH_SOURCE-Wächter bereitet Testbarkeit ausdrücklich vor, nichts nutzt ihn. F-1 ist genau die Fehlerklasse, die eine Probe mit dem d-check-Layout (Unterordner) gefangen hätte.
- **verifizierbar:** ja — `ls tools/harness/`, Makefile-Index; keine slice-mv-Probe.
- **klasse:** `neuer-vertrag-ohne-negativprobe`

### F-8 — INFO · Fehlermeldung des Bereits- dort-Checks nennt bei gesetztem Unterordner den falschen Ort

- **quelle:** Maintainability
- **pfad:** `tools/harness/slice-mv.sh:228-230`
- **befund:** `[ "$from" = "$TO" ]` bricht korrekt ab, nennt in der Meldung aber `$ziel` — bei gesetztem `SLICE_MV_DONE_UNTERORDNER` und flach liegendem Altbestand (slice-228…239 in `done/`) behauptet sie „liegt bereits in done/welle-91/", während die Datei in `done/` liegt. Exit und Abbruchverhalten sind in allen Fällen korrekt.
- **verifizierbar:** ja — Aufruf mit `TO=done` + Unterordner gegen einen flachen Bestand.
- **klasse:** `meldung-nennt-anderen-ort`

---

## Negativbefunde (geprüft, ohne Befund)

- **§3.1 Docker/make-only:** Die Host-bash-Frage ist im Plan (§6, Risiko 1) sauber verhandelt; die Skripte rufen nur die benannte Klasse (bash, git, sed, grep, coreutils) — dieselbe Klasse wie das bestehende `tools/harness/`. Kein Befund.
- **§3.3 Zwei-Commits-Regel:** Empirisch geprüft — Move-Commit und Verweis-Commit sind getrennt, Reihenfolge und `git mv`-Reinheit stimmen; die `SLICE_MV_DONE_UNTERORDNER`-Durchreichung ist im Makefile vorhanden. Kein Befund.
- **§3.9 SHA-Pins:** Kein Workflow im Diff. Kein Befund.
- **Hexagon/Import-Regeln (ADR-0005):** Kein Go-Code im Diff. Kein Befund.
- **§3.2/§3.6 Suppression/Schwellen:** Der Wächter ergänzt eine Prüfung; nichts wird unterdrückt oder gesenkt. Kein Befund.
- **Netzzugriff (DC-QA-03-Raum):** Keiner in den Skripten; `make gates` im Wegwerf-Klon (selbstpruefung) läuft netzlos über den committeten Bestand. Kein Befund.
- **gate-consistency-Form:** Drei neue README-Zeilen ↔ drei Makefile-Regeln, in beiden Richtungen deckungsgleich; die „kein Gate"-Kennzeichnung steht in der Zeile selbst, und die Einordnung des Vorlauf-Wächters als kein Gate entspricht `v6.13.0` · `regelwerk/modul-13-quality-gates.md` §Die dritte Lage wörtlich. Kein Befund.
- **slice-247-Form:** Lifecycle-Hinweis statt Status-Feld; `**Berührte Spec-Stellen:**` nennt `DC-FA-VCS-001` bzw. den Abschnitt; beide kanonischen Vorprüfungen tragen ehrliche `d-check:cite`-Anker (Zielzeilen 373-374/379 im vendorten Baum verifiziert — `MR-053`/`MR-054` erfüllt); Abgrenzung je Punkt mit Begründungsklasse; `**Verantwortlich:** —` korrekt für `open/`. Kein Befund.
- **Commit-Botschaften:** Alle vier tragen Kennungen; `8637f08a` benennt die gate-phantom-Behlung selbst (Zwischenstand roter Gates ist als Push-Regel-Form zulässig, Push-Tip ist grün). Kein Overclaim-Befund.
- **history-range-guard-Verhalten:** Leere auflösbare Range → Exit 1, unauflösbar → Exit 2, `--staged` → Exit 0, voller Klon `HEAD~1..HEAD` → Exit 0 — alles gemessen, deckungsgleich mit README-Zeile und ANLASS (nur die GRENZE-Formulierung ist F-5). Kein Verhaltens-Befund.
- **selbstpruefung-Design:** Rot- und Grün-Fall in einem Lauf schließen einen abgestürzten (statt haltenden) Traeger aus; `core.hooksPath` wird lokal, nicht mit Scope-Mischung gelesen; Träger-Abgleich `-ef` gegen die gerufene Datei. Kein Befund.
- **Identity-Fallback-Funktion (bash ≥ 4.4):** Gemessen — setzt `d-check slice-mv <d-check@local>`, ohne eine konfigurierte Identität zu überschreiben. Kein Befund (die Versions-Grenze ist F-4).

## Kategorie-Summary

**HIGH 0 · MEDIUM 3 (F-1, F-2, F-3) · LOW 4 (F-4, F-5, F-6, F-7) · INFO 1 (F-8)**

## Verdikt

MEDIUM blockiert typischerweise — hier mit einer Abstufung im Grad: **F-1** ist substantiell und sollte vor dem ersten produktiven Einsatz von `make slice-mv` in diesem Repo (Wellen-Closure von welle-91 steht unmittelbar bevor) gelöst sein, denn der erzeugte Commit behauptet nachgezogene Verweise und baut tote Links ein. **F-2** und **F-3** sind Form-Verstöße gegen benannte Regeln (`MR-045`, §3.7) und vor dem Commit des Implementers trivial behebbar — sie blockieren den Abschluss-Slice-Commit, nicht den Sachinhalt. Die LOW/INFO-Findungen sind Annahme- oder Begründungsfall, kein Konflikt-Pfad.
