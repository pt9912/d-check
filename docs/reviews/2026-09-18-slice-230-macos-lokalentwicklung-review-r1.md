# Review-Report: slice-230 — 2026-09-18

**Review-Art:** Code — geprüft gegen den Slice-Plan (`slice-230-macos-lokalentwicklung-image-test-wc-fix.md`), `AGENTS.md` §3 (Hard Rules) und die Reviewer-Skill-Prüffragen (Modul 10).

**Gegenstand:** slice-230 — Arbeitsbaum-Diff, noch nicht committet, auf HEAD `9614973e` (`tools/harness/fetch-baseline-cache.sh`, `docs/user/releasing.md`, `docs/plan/planning/in-progress/roadmap.md`, neue Datei `docs/plan/planning/in-progress/slice-230-macos-lokalentwicklung-image-test-wc-fix.md`).

**Skill:** `.harness/skills/reviewer.md` @ Commit `1b689b1f` (Version 1.16.0)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-18

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-230-macos-lokalentwicklung-image-test-wc-fix.md` (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §6 Risiken)
- kein `DC-*`, keine aktive ADR (Slice-Kopf `**Bezug:**` — Harness-Meta-Tooling ohne Produkt-Spec-Berührung)
- [MR-004](../../harness/conventions/MR-004-gate-nachweis-mechanik.md) (Gate-Nachweis-Mechanik `tools/harness/`)
- `AGENTS.md` §3.1 (Docker/make-only), §3.7 (Kommentar-Klassen), §5 (Botschaft verallgemeinert nicht über die Messung hinaus)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der neue Kommentar in `check_aliases()` schreibt die Eigenschaft „`-f` meldet ein fehlendes Ziel nicht zuverlässig per Exit-Code" spezifisch BSD zu („BSD gibt selbst bei unauflösbarem Ziel einen Pfad aus"). Laut GNU-coreutils-Doku gilt dieselbe Eigenschaft auch für GNU `readlink -f` (nur die vorletzte Komponente muss existieren, nicht die letzte) — die Formulierung grenzt die Ursache enger ein, als sie ist. Kein Funktionsfehler: das nachfolgende `[ -e "$tgt" ]` läuft plattformunabhängig. | Maintainability | `tools/harness/fetch-baseline-cache.sh:128-130` | nein — Dokumentations-Präzision, kein Gate prüft Kommentar-Fakten | Kommentar überzieht Platform-Spezifität einer Tool-Eigenschaft, die auf beiden Plattformen gilt |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `verify()` — Datei-Anzahl-Vergleich `-eq` statt `=` (`fetch-baseline-cache.sh:219-225`) | geprüft, ohne Befund. Bug real reproduziert auf diesem BSD-Host: `find … \| wc -l` liefert `on_disk="      54"` (rechtsbündiges Padding), `grep -c . "$sums"` liefert `manifest="54"` — String-Vergleich (`=`) schlägt fehl, numerischer Vergleich (`-eq`) korrekt. `bash tools/harness/fetch-baseline-cache.sh --verify` und `make baseline-verify` laufen mit dem Fix grün (0 Befunde, 54 Dateien). Kein Risiko durch nicht-numerische Operanden: beide Quellen sind strukturell Zeilenzahlen, unter `set -euo pipefail` würde ein Nicht-Integer-Vergleich ohnehin fail-closed abbrechen (Exit ≠ 0 aus `[`), nicht still grün werden (Prüffrage 1). |
| `check_aliases()` — `readlink -f` + explizites `[ -e "$tgt" ]` statt `readlink -e` (`fetch-baseline-cache.sh:122-132`) | geprüft, ohne Befund. Vorher/Nachher unabhängig reproduziert: `--selftest` mit dem alten Skriptstand (`git show HEAD:…`) liefert exakt 2 von 9 rot — „gesunder Alias" und „Alias auf Verzeichnis", die beiden einzigen `ok`-Fälle mit real auflösbarem Ziel; mit dem Fix 9 von 9 grün. Kein Silent-Green-Pfad identifiziert: Symlink-Zyklus, toter Alias (flach/Unterbaum/Punkt-Name) bleiben in beiden Fassungen `rot` (Prüffrage 1). |
| Prüffrage 6 (Kommentar-Klassen) auf beide neue Kommentare | geprüft — bis auf F-1 kein Befund. Beide Kommentare beantworten eine Kopplung, die der Code selbst nicht zeigt: warum eine scheinbar redundante Prüfung (`[ -e "$tgt" ]` neben `readlink -f`; `-eq` statt des naheliegenderen `=`) tatsächlich nötig ist — Klasse *Grenze* (dokumentiert eine Werkzeug-Eigenschaft, die das Verhalten erzwingt) bzw. *Kopplung* (zwei Prüfungen, die zusammen leisten, was vorher eine allein tat). Keine Review-Historie, keine Deliberation über Verworfenes, keine Herkunfts-Prosa. |
| Prüffrage 15 (Modul liest Eingaben, die es nicht scannt) auf `docs/user/releasing.md` | geprüft, ohne Befund. `docs/user/` liegt innerhalb der Scan-Wurzeln von `codepaths` (`roots: [docs, spec, tools, harness, internal, cmd]`) und `hostpaths`; der neue Inline-Code-Pfad `tools/image-test.sh` existiert real. Der neue Fenced-Codeblock (Docker-Wrapper-Rezept) wird von `codepaths`/`hostpaths` bewusst als opak behandelt (Fence-Ausnahme, spec-konform) — kein verdeckter Anker, keine stillschweigend übergangene Zusage. Das genannte Image `docker:27-cli` matcht nicht `versions.pin-pattern` (`ghcr\.io/…`) und löst dort keine falsche Erwartung aus. `make doc-check` läuft mit dem vollständigen Diff grün (803 Dateien, 0 Befunde), `make planning-check` ebenso (0 Befunde) — bestätigt u. a. die Zitier-Direktiven (`d-check:cite`) im neuen Slice-Plan als wortgleich. |
| Prüffrage 8 (Botschaft verallgemeinert über die Messung hinaus) auf den Slice-Plan | geprüft, ohne Befund. Die drei zahlenscharfen Behauptungen des Plans wurden unabhängig nachvollzogen und stimmen exakt: (a) „on_disk=„      54"` vs. manifest=„54"`" — reproduziert, identische Werte; (b) „2 von 9 rot — beide gesunden Alias-Fälle — vor dem Fix, 9 von 9 grün danach" — reproduziert mit dem alten Skriptstand aus `git show HEAD:…`, exakte Übereinstimmung inkl. der benannten zwei Fälle; (c) „`--selftest` (neun Proben) bleibt grün" — bestätigt. Die Formulierung „jeder gesunde Symlink … wurde dadurch fälschlich als unauflösbar gemeldet" ist keine Verallgemeinerung über die Stichprobe hinaus, sondern eine notwendige Folge des beschriebenen Mechanismus (jede Verwendung von `readlink -e` schlägt auf BSD mit „illegal option" fehl, unabhängig vom Ziel) — kein Proxy-Schluss. |
| Technische Korrektheit der `docs/user/releasing.md`-Ergänzung (ELF/Mach-O/Rosetta-Aussage, distroless-Begründung, Docker-out-of-Docker-`TMPDIR`-Mechanik) | geprüft, ohne Befund. `tools/image-test.sh` bestätigt die beschriebene `docker create`/`docker cp`-Extraktion und den nativen Host-Lauf zum byte-identischen Vergleich; `Dockerfile` bestätigt die distroless-Runtime-Stage (`gcr.io/distroless/static-debian12:nonroot`, kein Shell/`tar`). Die `TMPDIR`-Umlenkung auf einen host-identisch gemounteten Pfad ist die korrekte Antwort auf das bekannte Docker-out-of-Docker-Problem (verschachtelte `-v`-Mounts lösen gegen den Host-Daemon auf, nicht gegen den Wrapper-Container). |
| Scope-/Abgrenzungs-Treue (§1 des Slice-Plans) | geprüft, ohne Befund. Nur `tools/harness/fetch-baseline-cache.sh` geändert (kein Sammel-Fix über weitere `tools/harness/*.sh`); kein neues Makefile-Target; keine Berührung der Upstream-Freshness-Befunde. Erweiterung (b) (readlink-Fix) ist im Plan selbst vorab als Erweiterung derselben Datei/Funktion deklariert, nicht nachträglich stillschweigend mitgenommen. |
| Roadmap-Ruhe-Marker-Entfernung (`docs/plan/planning/in-progress/roadmap.md`) | geprüft, ohne Befund. Konsistent mit dem tatsächlichen Verzeichnis-Zustand (`slice-230…md` liegt in `in-progress/`); `make planning-check` bestätigt (0 Befunde). Kein Zustandsfeld mit Chronik (Prüffrage 7) — reine Streichung einer nicht mehr zutreffenden Aussage, keine neue Formulierung. |
| Hexagon-/Netz-/Suppressions-Anker (Prüffragen 3–5) | geprüft, ohne Befund. Reiner Shell-/Doku-Diff, kein Go-Code, kein Import, kein Netzzugriff außerhalb `external`, keine Inline-Suppression, keine Gate-Schwelle gesenkt. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Kommentar überzieht Platform-Spezifität einer Tool-Eigenschaft, die auf beiden Plattformen gilt

## Verdikt

**Merge-blockierend:** nein — ein LOW-Finding, kein HIGH/MEDIUM.

**Übergabe:** F-1 ist eine optionale Präzisierung (Kommentartext), keine Korrektur-Pflicht vor Closure; sie kann beim nächsten Anfassen der Zeile mitgenommen werden oder offen bleiben. Die beiden Kern-Fixes (numerischer Vergleich, `readlink -f`) sind unabhängig reproduziert und korrekt; die Doku-Ergänzung ist technisch akkurat und modul-scan-seitig unauffällig (`make doc-check`/`make planning-check`/`make baseline-verify` grün). Aus Reviewer-Sicht spricht nichts gegen die Closure des Slice. DoD-Abhaken, Closure-Notiz und die drei Paarungen bleiben Aufgabe des Verifiers/Planners (andere Rolle, anderer Prüf-Kontext).
