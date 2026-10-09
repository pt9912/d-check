# Review R2 — slice-262: Festlegungen der Netz- und Nachtlauf-Werkzeuge in die Spezifikation

**Review-Art:** Code (Nachprüfung der R1-Befunde gegen Plan, Konventionen und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-262, Commits `40f9c434` (Plan-Änderung) und `4f6803ea` (Korrekturen), `git show HEAD~1 HEAD`
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0 (`94f798ef`)
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Review R1 zu slice-262 (F-1 bis F-7); Slice-Plan slice-262 §3 in der
geänderten Fassung; `spec/spezifikation.md` §7 (`SPEC-102`, `SPEC-103`);
`AGENTS.md` §3.4, §3.7, §5 (Regeln 13, 15), §6 Schritt 4;
[`MR-025`](../../harness/conventions.md#mr-025). Gegenstand wie in R1:
`tools/harness/history-range-guard.sh`, `tools/harness/selbstpruefung.sh`,
`tools/harness/fetch-baseline-cache.sh`, das Rezept von `trace-check`, `.githooks/commit-msg`,
`.github/workflows/upstream-drift.yml`.

## Messungen (Kommando und Ergebnis)

| # | Kommando | Ergebnis |
|---|---|---|
| A | `make workflow-pins doc-check gate-consistency planning-check` | alle vier: `0 Befund(e)` |
| B | `grep` nach den alten Wortlauten („JEDER Pin", „zwei Action", „vier Versions", „Werkzeug-/Manifest", „fängt erst dieser", „wer die Achse") über den Baum ohne Baseline, `done/`, Reviews, CHANGELOG | ein Rest: Kopf von `fetch-baseline-cache.sh` (siehe F-2) |

## Abgleich der R1-Befunde

| R1 | Stand | Beleg |
|---|---|---|
| F-1 MEDIUM | **benannt, nicht behoben — angenommen.** SPEC-102 und `trace-check.md` nennen die Sperre im Message-Modus; die Code-Korrektur schließt die Abgrenzung aus („Werkzeuge"). Die Übergabe ins Beobachtungs-Register erfolgt bei der Closure und ist hier nicht prüfbar | `spezifikation.md` · „bricht der installierte Hook dort jeden Commit mit Exit 2 ab" |
| F-2 MEDIUM | erledigt: vier Versions-Achsen, fünf Digest-Achsen, drei Action-Pins, `GRENZE:` mit Trivy-Pins und `release.yml`-Buildern | `upstream-drift.yml` · „GRENZE: nicht gewacht werden der Scanner-Pin" |
| F-3 LOW | erledigt in `upstream-drift.yml` und `docker-make-only.md`; Rest im Skript-Kopf, siehe F-2 unten | `docker-make-only.md` · „mit einer Ausnahme: fehlt" |
| F-4 INFO | erledigt (Stand von `HEAD`, Netz der Image-Builds); Nachschärfung siehe F-3 unten | `spezifikation.md` · „des Stands von `HEAD` — nicht des Arbeitsbaums" |
| F-5 INFO | erledigt | `image-scan.md` · „ist nicht gemessen — beide Wege enden mit Exit 2" |
| F-6 INFO | erledigt | „wenn jemand die Achse" (beide Dateien) |
| F-7 INFO | erledigt; die Plan-Änderung steht in einem eigenen Commit **vor** dem Code (`AGENTS.md` §6 Schritt 4) | Plan §3 · „Plan-Änderung nach Review R1 (F-2, F-3, F-7)" |

## Findings

### F-1 — INFO — SPEC-102 nennt das leere Repo nicht

- **quelle:** Maintainability
- **pfad:** `spec/spezifikation.md` · „In einem Repo mit nur einem Commit oder einem Klon der Tiefe 1"
- **befund:** In einem Repo **ohne** Commit ist `HEAD` nicht auflösbar, und der Vorlauf endet ebenso mit Exit 2 — mit installiertem Hook scheitert dort schon der erste Commit. Die Aussage nennt den Fall mit einem Commit, nicht den mit null.
- **verifizierbar:** ja — `git init`, Makefile und Hook anlegen, `make hooks`, erster Commit
- **klasse:** randform-unvollstaendig

### F-2 — INFO — Zwei Rest-Stellen tragen „SKIP je Teil" ohne die `curl`-Ausnahme

- **quelle:** [`MR-025`](../../harness/conventions.md#mr-025)
- **pfad:** `tools/harness/fetch-baseline-cache.sh` · „KEIN fail-closed (Ausfall → SKIP je Teil)"; `Makefile` · „+ Content-Drift am gepinnten Tag (Netz, NICHT in gates, fail-open)"
- **befund:** Der Abschnittskommentar zu `check_latest` und der Hilfetext von `baseline-freshness` sagen weiter pauschal fail-open, obwohl dieselbe Funktion drei Zeilen tiefer bei fehlendem `curl` mit Exit 1 endet. Beide liegen in Werkzeug-Dateien, die die Abgrenzung ausschließt. Deshalb INFO und nicht LOW: Die Abgrenzung ist bewusst gezogen, ein Spiegel bleibt stehen. Der Kopf-Kommentar (Teil B, „Netz-/Werkzeug-/Manifest-Ausfall") ist im Kontext von `unzip`/`sha256sum` richtig.
- **verifizierbar:** nein
- **klasse:** spiegel-ausserhalb-abgrenzung

### F-3 — INFO — SPEC-103 verortet die Image-Builds nur in `make gates`

- **quelle:** Maintainability
- **pfad:** `spec/spezifikation.md` · „die Image-Builds von `make gates` im Klon können Basis-Images aus dem Netz ziehen"
- **befund:** Schon die zwei Commit-Versuche bauen ein Image: Der Hook ruft `make trace-check`, und das Target hängt an `build`. Die Netz-Aussage trifft deshalb früher zu, als der Satz sagt; die Folge, dass Netz nötig sein kann, bleibt richtig.
- **verifizierbar:** ja — `make -n trace-check MSGFILE=/dev/null` zeigt den Build
- **klasse:** grenze-zu-eng-verortet

## Negativbefunde

- **§3.7 am Workflow-Kopf:** Die neuen Zeilen sind Zusage (welche Achsen) und `GRENZE:` (was nicht gewacht wird). `SPEC-098`/`099`/`100` stehen als Rang-Zeiger auf das höher rangierende Stratum. Der Kopf enthält keine Review-Historie, keine Befund-Nummern und keine Slice-Kennung. Ohne Befund.
- **Workflow-Schritte und Gates unberührt:** Der Diff an `upstream-drift.yml` betrifft nur Kommentarzeilen, kein `uses:`, `run:` und keinen Job-Schlüssel. `make workflow-pins`, `make doc-check`, `make gate-consistency` und `make planning-check` liefern 0 Befunde (Messung A). Ohne Befund.
- **Referenz-Richtung (`AGENTS.md` §3.4):** Die geänderten Einträge SPEC-102/103 enthalten keine Slice-, ADR-, Wellen- oder Hash-Kennung; den Nachzug trägt die Historie-Zeile. `doc-check` (matrix) ist grün. Ohne Befund.
- **Plan-Treue:** Die Plan-Änderung (`40f9c434`) liegt vor dem Code-Commit; der Code-Commit bleibt in den dort genannten Dateien. `trace-check.md` fällt unter „`harness/sensors/*.md`". Kein Werkzeug und kein Makefile wurde geändert. Ohne Befund.
- **Commit-Botschaft `4f6803ea` (Regel 15):** Jede genannte Korrektur ist im Diff vorhanden. Die Botschaft sagt „benannt", nicht „behoben". Ohne Befund.
- **Alte Wortlaute (Messung B):** außer F-2 keine weitere Stelle in einem lebenden Dokument.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 0 | 0 | 3 |

## Verdikt

**Freigegeben.** Alle R1-Befunde sind erledigt oder, im Fall von F-1, mit Begründung
angenommen. Die drei INFO-Befunde blockieren nicht. Die Register-Eintragung zu R1 F-1 ist
eine Closure-Pflicht dieses Slice und hier nicht prüfbar.
