# Verifikation welle-96 — Closure-Trigger, Schritt 1

**Rolle:** Verifier (Modul 8, §Rollen-Sequenz für eine Welle: Verifier → Planner) ·
**Gegenstand:** Closure-Trigger aus
[welle-96](../plan/planning/welle-96-vcs-folgt-umbenennung.md) §3, geprüft in einem frischen Kontext ·
**Bezug:**
[`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in),
[`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus),
[ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) · Datum 2026-11-24.

Jeder Beleg stammt aus einem Lauf, den ich selbst gefahren habe. Das Trigger-Audit (Schritt 2) ist hier
nicht geprüft. Der Beleg geht über die Slice-DoDs hinaus; deren Verifikation steht in
[slice-291](2026-11-24-slice-291-vcs-umbenennung-verify.md) und
[slice-292](2026-11-24-slice-292-vcs-umbenennung-verify.md).

**Stand:** `make fullbuild` lief auf `274c6ba4` mit sauberem Arbeitsbaum. Während des Laufs kam `66eb94a3`
auf `main` (ein Commit nur mit Dokumentation: Regel 18 in `AGENTS.md`, eine Beobachtung, der
Trigger-Audit-Report — kein Go-Code, kein Makefile, kein Image-Bestandteil außer dem, was `docker build`
ohnehin liest). Die übrigen Gates, der Klon-Lauf und die Probe liefen auf `66eb94a3`, Arbeitsbaum sauber
bis auf diesen Report. Der Probe-Vergleich mit dem Eigenlauf lief auf `274c6ba4` und im Wegwerf-Klon auf
`66eb94a3`; beide ergaben dasselbe.

## Trigger-Punkte

| Punkt | Urteil | Beleg |
|---|---|---|
| 1 — alle Slices in `done/` | bestätigt | Die Welle trägt zwei Slices: [slice-291](../plan/planning/done/slice-291-vcs-umbenennung-anforderung-und-adr.md) und [slice-292](../plan/planning/done/slice-292-vcs-umbenennung-modulkern.md), beide in `done/` (Move von slice-292: `274c6ba4`). `open/` und `next/` tragen nur `.gitkeep`, `in-progress/` nur `roadmap.md`. Das `Welle:`-Feld beider Slices zeigt auf die Welle-Datei; ein dritter Slice mit diesem Zeiger existiert in keinem Lifecycle-Verzeichnis. |
| 2 — `make fullbuild` grün | bestätigt | Exit 0. Auszug: `verify ok (54 Dateien, vollständig)`; `0 issues.`; `gesamt: 0 Befund(e)` (arch-check); `coverage-gate: OK — Coverage 94.90% erfüllt Schwelle 93%`; `Ran 55 rules on 75 files: 0 findings.`; `[gates] … green`; `image-test: OK — Akzeptanzkriterien erfüllt (DC-FA-DIST-001, … DC-FA-LLEN-001)`; `bench: Median 1142 ms (Pass-Kriterium < 5000 ms — DC-QA-01)`; `58 Anforderung(en), 0 Waise(n).`; verify-closure-notes `1157 Datei(en) geprüft, 0 Befund(e)`. Abschluss: `[fullbuild] green — image-hash sha256:6d1195bcbb53a57ec3e72056b0d5f680332c9af080e2df7dfcdebaa494e6eadb` (gleich `docker image inspect d-check:latest` direkt danach). |
| 3 — ohne den Schalter nichts Neues | bestätigt, mit zwei Anmerkungen | Probe und Eigenlauf unten, dazu der zusätzliche Vorher/Nachher-Vergleich des Moduls `vcs` (60 Vergleiche). Die eine gewollte Abweichung ist die einzige. Anmerkung 1: der Eigenlauf der Probe weicht mit dem **eingeschalteten** Profil ab (siehe Abweichungen). Anmerkung 2: die Probe selbst fährt `vcs` nicht. |
| 4 — der Anlass trägt am gebauten Image | bestätigt | Tabelle Probe → Ergebnis unten: ohne den Schalter 59 `core-drift-vcs`, mit ihm 0, Gegenprobe 59. |

## Ohne den Schalter nichts Neues

### Die Probe

`make blackbox-probe REF=v0.92.0` auf `274c6ba4`, Skript-Exit 1, über `make` Exit 2:
`blackbox-probe: 4 von 20 Vergleichen weichen ab (Vorher 61a412eb).` Gleich sind die Fixtures `ids`, `links`,
`sauber`, `targets` (je vier Formen). Abweichend ist der Eigenlauf `repo` in allen vier Formen
(Standard, `--json`, `--yaml`, `--doctor`), jeweils Exit 2 vorher, Exit 0 nachher, mit der Ursache auf stderr
des Vorher-Images: `.d-check.yml: yaml: unmarshal errors: line 868: field follow-renames not found`.

Das Profil des Repos schaltet `vcs.follow-renames` seit dem Architect-Akt ein
([ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md), Entscheid 14); das Vorher-Image v0.92.0 kennt den
Schlüssel nicht und lehnt die Datei ab. Das ist kein Verhaltensunterschied des Moduls, sondern die Schärfe der
Konfigurations-Prüfung. Zur Gegenprobe habe ich in einem Wegwerf-Klon von `274c6ba4` genau diese eine Zeile
aus `.d-check.yml` entfernt (Arbeitsbaum-Änderung, `git diff --stat`: 1 Zeile gelöscht) und die Probe dort
gefahren: **`byte-identisch über 20 Vergleiche (Vorher 61a412eb, stdout/stderr/Exit getrennt, Kanarienlauf
vorher und nachher)`**, Exit 0. Für den Eigenlauf zeigt die Probe also: ohne den Schalter im Profil dieselbe
Ausgabe, mit ihm ein Profil, das ältere Images nicht lesen.

### Das Modul `vcs` selbst

Die Probe fährt `--range` und `--staged` nicht (Default-Belegung von `PROBE_FORMS`), das Modul `vcs` bleibt
dort unberührt. Deshalb habe ich es zusätzlich gegen das Vorher-Image verglichen: `d-check:w96-vorher`, gebaut
aus `git archive v0.92.0` (Commit `61a412eb`, `docker build --target runtime`, `VERSION=0.0.0-dev`) gegen das
Image aus `make fullbuild` (`sha256:6d1195bc…`). Beide laufen mit `--network none` und `-v <repo>:/repo:ro`;
stdout, stderr und Exit sind getrennt gespeichert und je mit `cmp` verglichen.

Fixtures: 28 Wegwerf-Repos, alle **ohne** `follow-renames`, je einmal ohne und einmal mit
`ignore-link-targets: true`; dazu eines ganz ohne Konfigurationsdatei und eines mit einem unlesbaren
Zwischen-Commit.

| Gruppe | Inhalt | Vergleiche | Ergebnis |
|---|---|---|---|
| r1 | Umbenennung auf neuen Namen samt Nachzug (drei ADRs: zwei `Accepted`, eine `Proposed`); Standard, `--json`, `--yaml` | 6 | gleich (Exit 1 ohne, Exit 1 mit `ignore-link-targets`) |
| r2 | Umzug mit gleichem Namen samt Nachzug; Umbiegen auf eine andere gleichnamige Datei; drei Ranges | 6 | gleich (Exit 1 ohne, Exit 0 mit `ignore-link-targets`) |
| r3 | Statuswechsel, Inhaltsänderung, Löschen, Neuanlage, Merge, Kette, fremde Basis (`HEAD..side`, Exit 2), `--json`; neun Ranges | 18 | gleich |
| r4 | Gitlink ersetzt eine immutable Datei (`--range`), neuer Gitlink, Gitlink ersetzt eine `Proposed`-ADR | 6 | gleich (Exit 1, 0, 0) |
| s1 bis s9 | `--staged`: Umbenennung mit Nachzug, Nachzug allein, Statuswechsel, Inhaltsänderung, leerer Index, Löschen, neuer Gitlink, Gitlink ersetzt `Proposed`-ADR | 16 | gleich |
| s6 | `--staged`: **Gitlink ersetzt eine immutable Datei**, Standard und `--json` | 4 | **vorher Exit 2, nachher Exit 1** |
| r6, r7 | ohne Konfigurationsdatei; Basis..Spitze über einen Commit mit gelöschtem Baum-Objekt (Exit 0 beidseitig — das Objekt liest ohne den Schlüssel niemand), Basis..Zwischen-Commit (Exit 2 beidseitig) | 4 | gleich |
| | | **60** | **56 gleich, 4 abweichend** |

Die vier Abweichungen sind **alle** dieselbe: gestagter Modus, ein Gitlink ersetzt an seinem Pfad eine immutable
Datei. Vorher stdout leer, stderr `d-check: error: object not found`, Exit 2; nachher
`adr/0001.md:1 … core-drift-vcs … immutable Datei gelöscht oder umbenannt — der Pfad einer immutablen Datei ist stabil`,
Exit 1. Das ist die im Welle-Plan genannte Abweichung
([Architect-Verdikt](2026-11-24-slice-292-architect-verdikt.md) F-4,
[slice-292](../plan/planning/done/slice-292-vcs-umbenennung-modulkern.md) §3). Dieselbe Konstellation im
`--range`-Modus (r4) ist beidseitig gleich, ebenso ein Gitlink, der keine immutable Datei ersetzt (r4, s7, s9).
Dass der Vergleich eine Abweichung findet, wenn eine existiert, zeigt er selbst: er hat genau diese vier
gefunden und keine andere.

## Der Anlass am gebauten Image

Wegwerf-Klon von `66eb94a3`, ein Commit: `git mv` der Spezifikationsdatei von `spezifikation.md` auf
`specification.md` im Verzeichnis `spec/`; in allen ADRs, die sie verlinken (63, davon 60 `Accepted`), wurde nur
das Link-Ziel nachgezogen. Die Rückübersetzung des neuen Stands auf den alten ergibt für jede der 63 Dateien
byte-gleich die Fassung vor dem Commit (`sed` zurück, `diff`): es ändert sich nichts außer dem Dateinamen im
Link-Ziel; Git meldet `R100`. Dazu die Verweise der Konfiguration auf die Datei in `.d-check.yml` (nötig,
siehe Abweichungen). Das Image ist das aus `make fullbuild` (`sha256:6d1195bc…`), aufgerufen über
`make -o build adr-check RANGE=HEAD~1..HEAD` (unveränderter Rezept-Teil: `history-range-guard`, dann
`docker run … --enable vcs <FOCUS_DISABLE> --range HEAD~1..HEAD`, `--network none`, `:ro`).

| Probe | Profil | Ergebnis |
|---|---|---|
| Nachzug nach Umbenennung, **ohne** `vcs.follow-renames` (Zeile aus dem Profil entfernt, `ignore-link-targets: true` bleibt) | ohne Schalter | Exit 1, `1297 Datei(en) geprüft, 59 Befund(e)`: 59 verschiedene ADRs, je ein `core-drift-vcs` (`Core einer immutablen Datei hat sich über die Commit-Range geändert`) |
| derselbe Commit, `follow-renames: true` wie im Repo | mit Schalter | Exit 0, `1297 Datei(en) geprüft, 0 Befund(e)` |
| Gegenprobe: dasselbe, die Datei aber gelöscht und `specification.md` mit einem zusätzlichen Absatz neu angelegt (Git meldet `R099`; das Modul vergleicht Inhalts-Gleichheit) | mit Schalter | Exit 1, 59 Befunde, dieselbe Menge ADRs wie ohne den Schalter (`comm` leer) |
| dasselbe Image aus dem Vorher-Stand (v0.92.0) gegen den Klon, Profil ohne Schalter | ohne Schalter | Exit 1, 59 Befunde, stdout und stderr **byte-gleich** zum Nachher-Image ohne Schalter |

Warum 59 und nicht 60: Von den 60 `Accepted`-ADRs mit einem Link auf die Datei trägt
[ADR-0103](../plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md) ihn nur in `## Geschichte`
(`vcs.exclude-sections`); die Änderung dort ist auch ohne Schalter keine Drift. Die Menge der 59 ist die der
60 abzüglich dieser einen (`comm` zwischen der Liste der geänderten `Accepted`-ADRs und der der gemeldeten).
Das Rot kommt also vom Schalter: dieselbe Eingabe ergibt ohne ihn 59, mit ihm 0, und mit ihm und einer
anderen Datei wieder 59. Ein stilles Grün über einem leeren Prüfbereich ist ausgeschlossen: `1297 Datei(en)
geprüft` steht in beiden Läufen, `history-range-guard` hat die Range akzeptiert.

## Weitere Belege

| Beleg | Kommando (auf `66eb94a3`) | letzte Zeile |
|---|---|---|
| Gates | `make gates`, Exit 0 | `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` (davor `d-check: 1297 Datei(en) geprüft, 0 Befund(e)`) |
| Image-Test | `make image-test`, Exit 0 | `image-test: OK — Akzeptanzkriterien erfüllt (DC-FA-DIST-001, DC-FA-CLI-007, DC-FA-CLI-008, DC-FA-CLI-013, DC-FA-MET-001, DC-FA-PSYNC-001, DC-FA-LLEN-001)` |
| Coverage | `make coverage-gate`, Exit 0 | `coverage-gate: OK — Coverage 94.90% erfüllt Schwelle 93%` |
| Vollständigkeit | `make completeness-check`, Exit 0 | `58 Anforderung(en), 0 Waise(n).` |
| Closure-Notizen | `make verify-closure-notes`, Exit 0 | `d-check: 1158 Datei(en) geprüft, 0 Befund(e)` |
| Doku nach dem Schreiben dieses Reports | `make doc-check`, Exit 0 | `d-check: 1304 Datei(en) geprüft, 0 Befund(e)` |

## Abweichungen vom Welle-Plan

1. **Eigenlauf der Probe.** Der Plan sagt, die Probe zeige für den Eigenlauf dieselbe Ausgabe. Mit dem
   eingeschalteten Profil trifft das nicht zu: das Vorher-Image lehnt den Schlüssel `follow-renames` ab
   (Exit 2 gegen Exit 0, vier von 20 Vergleichen). Der Plan nennt die Probe-Gleichheit für ein Profil ohne den
   Schalter; dass Entscheid 14 das eigene Profil in derselben Welle einschaltet, nimmt dem Eigenlauf diese
   Aussage. Belegt ist die Gleichheit ohne die eine Zeile (20 von 20 byte-identisch). Folgen für Adopter:
   ein Profil mit dem neuen Schlüssel liest nur ein Image ab dieser Version — das gilt für jeden neuen
   Schlüssel (strikte Dekodierung) und ist keine Eigenschaft von `follow-renames`.
2. **`make adr-check` im Klon baut nicht.** Die Umbenennung der Spezifikationsdatei bricht `make build` im
   Klon (`manual.go:22:12: pattern spec/spezifikation.md: no matching files found`: das Werkzeug bettet die
   Datei ein) und die Konfiguration (`konfiguriertes ids-Target existiert nicht`). Ich habe die Verweise in
   `.d-check.yml` im Klon nachgezogen und das Target mit `make -o build adr-check` gegen das fertige Image
   gefahren. Das Ergebnis des Plans (ein Befund je ADR ohne, keiner mit Schalter) ist davon unberührt. Für
   den Planner: die Umbenennung der Spec-Dateien braucht außer den ADR-Links auch den Einbettungs-Pfad im
   Go-Quelltext und die Verweise in `.d-check.yml` — das gehört nach dem Welle-Plan §6 zur Baseline-Hebung, die
   sie verlangt, und ist dort noch nicht vorgemerkt.
3. **59 statt je ADR.** Der Plan sagt „je ADR einen `core-drift-vcs`“; es sind 59 von 60 `Accepted`-ADRs mit
   Link, weil eine ihn nur in `## Geschichte` trägt. Der Bestand ist seit der Zählung im Plan gewachsen
   (60 von 108 gegen 59 von 107 — ADR-0115 ist hinzugekommen); die Aussage bleibt.
4. **Zwei Stände.** `make fullbuild` auf `274c6ba4`, der Rest auf `66eb94a3`; der Unterschied ist ein reiner
   Dokumentations-Commit (siehe Stand oben). Ich habe `make fullbuild` nicht wiederholt; die Gates liefen auf
   dem neuen Stand erneut grün.

Keine weitere Abweichung: die Slice-Zahl, die Reihenfolge und die gewollte Abweichung im gestagten Modus
entsprechen dem Plan.

## Grenzen dieser Verifikation

- Die 28 Fixtures decken Umbenennung, Kette, Merge, Gitlink, Status-, Inhalts- und Anhangs-Änderung, fremde
  Basis und einen unlesbaren Zwischen-Commit; sie sind eine Auswahl, kein Beweis. Die Gleichheit im ganzen
  Raum tragen die Tests zum Kriterium „Boundary (ohne `vcs.follow-renames` byte-identisch)“ von
  [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  und die Mutationen aus der Slice-Verifikation.
- Der Fall r6 (ohne Konfigurationsdatei) endet beidseitig ohne Befund; er belegt Gleichheit, nicht die
  Schärfe der Default-Belegung.
- Das Trigger-Audit, die Paarungen und `welle-96-results.md` sind nicht Gegenstand dieses Schritts.

## Ergebnis

Die vier Trigger-Punkte sind bestätigt. Die eine gewollte Abweichung im gestagten Modus ist die einzige; der
Anlass trägt am gebauten Image (59 Befunde ohne den Schalter, 0 mit ihm, 59 in der Gegenprobe). Image-Hash
des Laufs `make fullbuild`:
`sha256:6d1195bcbb53a57ec3e72056b0d5f680332c9af080e2df7dfcdebaa494e6eadb`. Offen für die folgenden
Closure-Schritte (Planner): `welle-96-results.md` mit der Tabelle oben, und der Hinweis zur Umbenennung der
Spec-Dateien (Abweichung 2).
