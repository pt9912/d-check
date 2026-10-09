# Review R1 — slice-268: `releasing.md` zieht nach docs/user/maintainer/

**Review-Art:** Code (Diff gegen Plan, ADRs und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-268, Range `a9d62bab..HEAD` — `bcb104fe` (Move + eingehende
Verweise), `420e401a` (eigene Links der bewegten Datei)
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan slice-268 (inkl. Plan-Änderung §3);
[ADR-0025](../plan/adr/0025-codepaths-ignore-refs.md) (Tombstone-Register),
[ADR-0103](../plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md) (Pfad-Nachzug),
[ADR-0097](../plan/adr/0097-matrix-aussen-adaptionsblock-historie-status-ausnahmen.md)
(`matrix.exempt-paths`); `AGENTS.md` §3.3, §3.5, §3.6, §3.7, §6;
[`MR-070`](../../harness/conventions.md#mr-070); Register
`BEO-ALL/pfad-nachzug-gleicher-name-andere-datei`,
`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`; vorherige Findings am
selben Gegenstand: Reviews zu slice-267 (r4, r7, r8, verify).

## Messungen (Kommando und Ergebnis)

- **Alter Pfad, volle Form:** `git grep -n 'docs/user/releasing\.md' HEAD -- .`
  → nach Klasse:
  - *Konfiguration (gewollt):* `.d-check.yml` Tombstone-Kommentar + `ignore-refs`-Eintrag.
  - *Inline-Code, eingefroren:* ADR-0097 (1), ADR-0103 (1), `done/`-Slices 230 (4),
    236 (1), 248 (2), 256 (3), 267 (2), Register-Beleg `citation-stretched-beyond-scope/evidence/slice-181.md` (1),
    `MR-010` in `conventions/done/` (1), `CHANGELOG.md` Abschnitt 0.84.0 (1).
  - *Inline-Code, lebend:* slice-268 selbst (2) — benennt den Umzug.
  - *Review-Reports* (`docs/reviews/**`, exempt): 14 Zeilen.
  - *Go-Test-Fixtures* (`internal/…/vcs_pfad_nachzug_test.go`): Strings, kein Verweis.
- **Relative alte Form / Link-Ziele:** `git grep -nE '\]\([^)]*releasing\.md[^)]*\)' HEAD -- . ':!docs/reviews' ':!internal' | grep -v maintainer/releasing.md`
  → nur `spec/spezifikation.md:1974`, ein Beispiel in Inline-Code (`[R](../user/releasing.md#prep)`), kein Verweis.
- **Links, Ist:** kein Link im Repo (außer Reviews/Tests) zielt mehr auf den alten Pfad.
- **Gegenprobe der 19 `codepath-missing`:** Wegwerf-Worktree auf HEAD, Tombstone-Eintrag
  entfernt, `docker run --rm --network none -v <wt>:/repo:ro d-check:latest`
  → **19 `codepath-missing`, sonst nichts** — Verteilung genau wie oben (inkl. CHANGELOG und slice-268 ×2). Zahl in Plan §3 und Commit-Botschaft stimmt.
- **16 eingehende Links (Commit-Botschaft):** CHANGELOG 1 + README 2 + README.de 2 +
  Handbuch 2 + Hub-README 1 + ADR-0014 4 + ADR-0067 1 + drei Wellen-Notizen 3 = 16 ✓.
  Plan §1 „12 lebende Stellen + 9 eingefrorene Links": 7 lebende Links + 5 Stellen ohne
  Doku-Gate = 12; ADR 5 + Wellen 3 + CHANGELOG 1 = 9 ✓.
- **Nur Link-Ziele geändert:** je Datei `diff <(git show a9d62bab:F | sed -E 's/\]\([^)]*\)/](X)/g') <(git show HEAD:F | sed …)`
  → nur README.md / README.de.md mit sichtbarem Text geändert (lebend, gewollt);
  CHANGELOG, ADR-0014, ADR-0067, welle-73/74/82-results, Benutzerhandbuch (inkl. §11-Zeile 1.29),
  Hub-README: **nur Ziele**. Ebenso `420e401a` (23 Ziele, sonst nichts).
- **Rename-Erkennung:** `git show -M --summary bcb104fe` → `rename docs/user/{ => maintainer}/releasing.md (100%)`.
- `make adr-check RANGE=a9d62bab..HEAD` → `1078 Datei(en) geprüft, 0 Befund(e)`.
- `make doc-check` → `1084 Datei(en) geprüft, 0 Befund(e)`.

## Findings

### LOW-1 — Commit-Schnitt weicht vom Plan ab, die Plan-Änderung nennt es nicht

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §6 Schritt 4 (Plan-Änderung vor den Code); `AGENTS.md` §3.3
- **pfad:** slice-268 §1 · „ein reiner `git mv`, danach ein eigener Commit, der die eigenen relativen Links der Datei und alle eingehenden Verweise nachzieht"; Commit `bcb104fe`
- **befund:** Die eingehenden Verweise (15 Dateien) reisen im Move-Commit mit, nicht im Folge-Commit wie in §1 geplant; die Plan-Änderung in §3 erwähnt den geänderten Schnitt nicht, nur die Commit-Botschaft begründet ihn nachträglich (gestagter `adr-check` gegen BASE ohne alten Pfad — die Begründung trifft am Gegenstand zu, ADR-0103 normiert nur Ziele, die in BASE oder HEAD auflösen). `AGENTS.md` §3.3 nennt das Mitreisen anderer Dateien ausdrücklich nur für Fall 2 (Übergang nach `done/`). Die Rename-Erkennung ist unberührt (100 %).
- **verifizierbar:** nein (Urteil; kein Gate sieht die Commit-Zerlegung)
- **klasse:** plan-schnitt-ohne-plan-aenderung

### LOW-2 — Der Tombstone deckt einen umgezogenen, weiter lebenden Pfad repo-weit ab

- **kategorie:** LOW
- **quelle:** [ADR-0025](../plan/adr/0025-codepaths-ignore-refs.md) §Konsequenzen „Globaler Scope"
- **pfad:** `.d-check.yml` · „`- docs/user/releasing.md: umgezogen nach docs/user/maintainer/releasing.md`"
- **befund:** Der Eintrag ist durch ADR-0025 gedeckt (der alte Pfad existiert nicht mehr; Refactoring ist der Anlassfall der ADR) und damit keine undeklarierte Senkung nach §3.6. Er wirkt aber referenz-weit auch in lebenden Dateien: Nennt künftig ein Slice-Plan, das Handbuch oder ein CHANGELOG-Abschnitt den alten Pfad in Inline-Code — naheliegend, weil 19 eingefrorene Stellen ihn als Kopiervorlage tragen —, bleibt `make doc-check` grün. ADR-0025 nimmt das für **entfernte** Artefakte mit „lebende Verweise räumt man ohnehin auf" in Kauf; für eine umgezogene, weiter benutzte Datei ist die Wahrscheinlichkeit einer Neunennung höher, und weder Plan noch Commit noch Kommentar benennen die Lücke. Die Link-Achse bleibt gedeckt (`links` meldet einen Link auf den alten Pfad).
- **verifizierbar:** ja — Inline-Code `` `docs/user/releasing.md` `` in einer lebenden Datei ergänzen, `make doc-check` bleibt bei 0 Befunden
- **klasse:** tombstone-fuer-umzug-deckt-lebende-neunennung

### INFO-1 — Commit-Botschaft nennt alle 19 Inline-Code-Stellen „historisch"

- **kategorie:** INFO
- **quelle:** `AGENTS.md` §5 Regel 15
- **pfad:** Commit `bcb104fe` · „Die 19 Inline-Code-Erwaehnungen des alten Pfads -- in zwei Accepted-ADRs, geschlossenen Slices, Register-Belegen, einem aufgeloesten MR, einem alten CHANGELOG-Eintrag -- bleiben als historischer Text"
- **befund:** Zwei der 19 stehen im lebenden Plan slice-268 selbst (Plan §3 nennt sie richtig); „Register-Belege" ist ein einziger Beleg. Die Zahl stimmt, nur die Aufzählung ist unvollständig — kein Versagen.
- **verifizierbar:** ja (Gegenprobe oben)
- **klasse:** botschaft-aufzaehlung-unvollstaendig

### INFO-2 — Risiko §6 nennt die Docker-Hub-Beschreibung als Außen-Verweis

- **kategorie:** INFO
- **quelle:** slice-268 §6
- **pfad:** slice-268 §6 · „Links von außen (Docker-Hub-Beschreibung, Adopter) zeigen auf den alten Pfad"
- **befund:** Die auf Docker Hub hochgeladene Seite ist `packaging/dockerhub/overview.md` (`hub-description.yml`, `OVW=`), und sie verlinkt `releasing.md` nicht; `packaging/dockerhub/README.md` wird nur im Repo gerendert und ist nachgezogen. Für den Ausgang des Risikos bei der Closure relevant: der Hub-Teil trägt nicht.
- **verifizierbar:** ja (`grep releasing packaging/dockerhub/overview.md` leer)
- **klasse:** risiko-beispiel-trifft-nicht-zu

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Vollständigkeit lebender Stellen | geprüft, ohne Befund. Kein lebender Link, keine lebende Konfiguration und keine Stelle ohne Doku-Gate (`.github/dependabot.yml`, `release.yml` ×2, `hub-description.yml`, `tools/image-test.sh`) nennt den alten Pfad; verbleibend nur slice-268 (gewollt), ein Spec-Beispiel und Go-Test-Fixtures (Illustration, kein Verweis). |
| Frozen-Klassen (MR-070) | geprüft, ohne Befund. In ADR-0014, ADR-0067, welle-73/74/82-results, CHANGELOG-Altabschnitt und Handbuch-§11-Zeile nur das Link-Ziel geändert, sichtbarer Text identisch (Kommando oben). Eingefrorener Inline-Code unverändert. |
| ADR-Nachzug auf dieselbe Datei (`BEO-ALL/pfad-nachzug-gleicher-name-andere-datei`) | geprüft, ohne Befund. Alle fünf ADR-Ziele (ADR-0014 ×4, ADR-0067 ×1) zeigen auf `docs/user/maintainer/releasing.md`, dieselbe umgezogene Datei (Rename 100 %). |
| Tombstone-Deckung §3.6 und Kommentar §3.7 | geprüft — gedeckt durch ADR-0025 (Grenze siehe LOW-2). Kommentar folgt der Form der Nachbar-Einträge (Was + Kopplung an die eingefrorenen Nennungen), keine Review-Historie, keine Slice-Nummer. |
| `matrix.exempt-paths` | geprüft, ohne Befund. Eintrag auf neuen Pfad gezogen; Begründungs-Kommentar („die Release-Doku referenziert den aktuellen Stand") trifft weiter zu; keine weiteren `docs/user`-Globs in `.d-check*.yml`/`.a-check.yml`, die die Datei durch den Umzug verlören. |
| Commit-Zerlegung §3.3 / Rename | geprüft — bewegte Datei im Move-Commit unverändert, `R100`; eigene Links im Folge-Commit (23 Ziele, nur Ziele). Plan-Abweichung siehe LOW-1. |
| Anker `#release-prep-vor-dem-tag` | geprüft, ohne Befund. Überschrift unverändert, welle-73/74-Links lösen auf (`make doc-check` 0 Befunde, Modul `anchors` aktiv). |
| Eigene Links der bewegten Datei | geprüft, ohne Befund. `../../` → `../../../`, `../plan/` → `../../plan/`, Geschwister → `../`; alle auflösend. |
| Docker-Hub-README | geprüft, ohne Befund. `packaging/dockerhub/README.md` wird im Repo (GitHub) gerendert, der relative Link löst dort auf; die Hub-Seite selbst ist `overview.md` mit absoluten GitHub-URLs, ohne `releasing.md`. |
| Gates | `make adr-check RANGE=a9d62bab..HEAD` 0 Befunde; `make doc-check` 0 Befunde. |
| Kommentare in Workflows/Skripten | geprüft, ohne Befund. Nur Pfad im bestehenden Rang-Zeiger bzw. in Meldungstexten geändert. |

## Kategorie-Summary

HIGH 0 · MEDIUM 0 · LOW 2 · INFO 2

## Verdikt

Kein blockierendes Finding. LOW-1 und LOW-2 annehmen oder begründen;
INFO-2 gehört in den Ausgang des Risikos bei der Closure.
