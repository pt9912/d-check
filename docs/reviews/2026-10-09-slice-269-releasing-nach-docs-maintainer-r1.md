# Review R1 — slice-269: `releasing.md` zieht nach docs/maintainer/

**Review-Art:** Code (Diff gegen Plan, ADRs und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-269, Range `2b088bb1..HEAD` — `877af2d4` (Move + eingehende
Verweise + Tombstone), `d0e64cdb` (eigene Links der bewegten Datei)
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan slice-269;
[ADR-0025](../plan/adr/0025-codepaths-ignore-refs.md) (Tombstone-Register),
[ADR-0103](../plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md) (Pfad-Nachzug),
[ADR-0097](../plan/adr/0097-matrix-aussen-adaptionsblock-historie-status-ausnahmen.md)
(`matrix.exempt-paths`); `AGENTS.md` §2, §3.3, §3.5, §3.6, §3.7, §5, §6;
[`MR-070`](../../harness/conventions.md#mr-070); Baseline `v6.17.0` ·
`regelwerk/grundlagen-source-precedence.md` §Source Precedence und
`regelwerk/modul-01-entwicklungszyklus.md` §Ziel-Form: Source-Precedence-Block;
vorherige Findings am selben Gegenstand: Review R1 und Verifikation zu slice-268.

## Messungen (Kommando und Ergebnis)

- **Zwischen-Pfad, alle Nennungen:** `grep -rn 'user/maintainer' . -I` (ohne `.git`) →
  - *Konfiguration (gewollt):* `.d-check.yml` Tombstone-Kommentar + `ignore-refs`-Eintrag.
  - *Text, lebend:* slice-269-Plan (3) — benennt den Umzug.
  - *Eingefroren:* slice-268 (4: Titel, §1, §3 als Text; §7 Z. 106 als Inline-Code `` `docs/user/maintainer/` ``).
  - *Review-Reports* (exempt): slice-267 r4–r8/verify, slice-268 r1/verify.
  - *Go-Test-Fixtures* (`vcs_pfad_nachzug_test.go`): Strings, kein Verweis.
  - **Kein Link, keine Workflow-/Skript-Stelle und keine lebende Doku nennt den Zwischen-Pfad.**
- **Alle `releasing.md`-Nennungen je Datei:** `git grep -n 'releasing\.md' HEAD -- . ':!docs/reviews' | cut -d: -f2 | sort | uniq -c`
  → 32 Dateien; Reviews: 17 Dateien.
- **Eingehende Link-Ziele aufgelöst:** `git grep -noE '\]\([^)]*releasing\.md[^)]*\)' HEAD -- . ':!docs/reviews' ':!internal'`,
  je Treffer `realpath -m` relativ zur Quelldatei → **16 × `docs/maintainer/releasing.md`**,
  1 × `user/releasing.md` (Inline-Code-Beispiel in `spec/spezifikation.md`, kein Verweis).
  Verteilung: CHANGELOG 1, README 2, README.de 2, Handbuch 2, Hub-README 1, ADR-0014 4, ADR-0067 1, welle-73/74/82 je 1 = 16 ✓.
- **Eigene Links vorher/nachher:** Ziele aus `git show 2b088bb1:docs/user/maintainer/releasing.md` und
  `git show HEAD:docs/maintainer/releasing.md` je mit `realpath -m` aufgelöst, `diff` → **23 = 23, identisch**,
  alle Ziele existieren. Vorher-Klassen: 11 × `../../../` (wird `repo-escape`), 12 andere (werden `target-missing`) ✓ Plan §3.
- **Nur Link-Ziele geändert:** je Datei `diff <(git show A:F | sed -E 's/\]\([^)]*\)/](X)/g') <(git show B:F | sed …)`
  → CHANGELOG, ADR-0014, ADR-0067, welle-73/74/82-results, Benutzerhandbuch (inkl. §11-Zeile 1.29), Hub-README:
  **nur Ziele**; README.md/README.de.md mit sichtbarem Link-Text (lebend, gewollt); `d0e64cdb`: nur Ziele.
- **Rename:** `git show -M --summary 877af2d4` → `rename docs/{user => }/maintainer/releasing.md (100%)`.
- **Frischer Klon** (`git clone` nach Scratchpad, `docs/user/` dort ohne `maintainer/`):
  `make doc-check` → `1094 Datei(en) geprüft, 0 Befund(e)`;
  `make adr-check RANGE=2b088bb1..HEAD` → `1088 Datei(en) geprüft, 0 Befund(e)`;
  `make trace-check RANGE=2b088bb1..HEAD` → `0 Befund(e)`.
- **Tombstone, Bruchprobe im Klon:** Eintrag `"docs/user/maintainer"` entfernt,
  `docker run --rm --network none -v <klon>:/repo:ro d-check:latest` →
  **1 Befund:** `slice-268-…md:106  docs/user/maintainer/  codepath-missing`. Eintrag zurück, dazu in
  `docs/user/operations.md` Inline-Code `` `docs/user/maintainer/releasing.md` `` und `` `docs/user/maintainer/` `` →
  **1 Befund:** `operations.md:189  docs/user/maintainer/releasing.md  codepath-missing` — der Eintrag deckt
  nur das Verzeichnis selbst (`path.Match`, abschließender Slash normalisiert), nicht Pfade darunter.
- Lokal ist `docs/user/maintainer/` ebenfalls entfernt (`ls docs/user`).

## Findings

### MEDIUM-1 — Die Source Precedence führt `releasing.md` nicht mehr

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §2 (Rang 6); `harness/README.md` §Source precedence (Rang 6); `v6.17.0` · `regelwerk/modul-01-entwicklungszyklus.md` §Ziel-Form: Source-Precedence-Block („Wahl und Begründung gehören in den Adaptions-Block")
- **pfad:** `AGENTS.md` §2 · Rang 6, Link auf `docs/user/`, „Operations, Releasing."; `harness/README.md` §Source precedence · Rang 6, Link auf `docs/user/`, „Operations, Releasing"
- **befund:** Beide Rangtabellen ordnen „Releasing" dem Rang `docs/user/` zu (wie der Baseline-Default `docs/user/*.md` mit Releasing-Sichten); mit dem Umzug liegt die einzige Releasing-Doku außerhalb jedes Rangs, und keine `MR`-Adaption deklariert den neuen Ort. Failure-Szenario: Widersprechen sich `docs/maintainer/releasing.md` und `README.md` (beide beschreiben den Digest-Pin-Weg), findet ein Agent die Datei in keiner Rangzeile, sucht sie laut Rang 6 vergeblich unter `docs/user/` und kann nicht entscheiden, welche Seite angepasst wird — oder stuft die ungerankte Datei unter README (7) ein und kehrt die gemeinte Ordnung um. Bei slice-268 trat das nicht auf, weil `docs/user/maintainer/` noch unter Rang 6 lag; Plan §1 grenzt es nicht ab.
- **verifizierbar:** nein (Urteil; kein Gate liest die Rangtabelle gegen die Ablage)
- **klasse:** umzug-verlaesst-precedence-rang

### LOW-1 — Slice-Nummer im Tombstone-Kommentar

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §3.7 („keine Slice-Nummern"; Herkunft nur als ein Feld `DC-*`/`ADR-*`/`MR-*`/`seit welle-<NN>`)
- **pfad:** `.d-check.yml` · „`- docs/user/maintainer: Zwischen-Ort derselben Datei; die Closure-Notiz von slice-268 nennt ihn in Inline-Code`"
- **befund:** Der Kommentar trägt eine Kopplung (welche eingefrorene Stelle den Eintrag braucht) und ist damit nicht klassenlos (kein HIGH), nennt den Träger aber per Slice-Nummer; der Nachbar-Eintrag beschreibt dieselbe Kopplung ohne Nummer („geschlossene Slices nennen den alten Pfad"), und R1 zu slice-268 hielt genau dieses Fehlen der Nummer als Negativbefund fest.
- **verifizierbar:** nein (kein Gate prüft Kommentarinhalt)
- **klasse:** slice-nummer-im-kommentar

### LOW-2 — Grenz-Kommentar behauptet für den Verzeichnis-Eintrag mehr Stille, als er erzeugt

- **kategorie:** LOW
- **quelle:** `AGENTS.md` §5 Regel 13 (Grenze gegen den Gegenstand prüfen)
- **pfad:** `.d-check.yml` · „`GRENZE: beide Einträge wirken repo-weit -- nennt lebende Doku einen der`" / „`alten Pfade künftig in Inline-Code, bleibt das still`"
- **befund:** Für `docs/user/maintainer` bleibt nur das Verzeichnis selbst still; der alte **Datei**-Pfad `docs/user/maintainer/releasing.md` in Inline-Code wird weiter als `codepath-missing` gemeldet (Bruchprobe oben). „einen der alten Pfade" liest sich, als sei auch dieser gedeckt; die Grenze ist zu weit, nicht zu eng — kein Versagen, aber gegen ihre eigene Beschreibung statt gegen `path.Match` geschrieben.
- **verifizierbar:** ja — Bruchprobe im Klon (s. Messungen)
- **klasse:** grenze-nicht-am-gegenstand-gemessen

### INFO-1 — Plan §3 schreibt dem CHANGELOG zwei eingehende Links zu

- **kategorie:** INFO
- **quelle:** `AGENTS.md` §5 Regel 14
- **pfad:** slice-269 §3 · „(CHANGELOG, README und README.de je zwei, Benutzerhandbuch zwei,"
- **befund:** Der CHANGELOG trägt einen Link (Z. 893); seine zweite `releasing.md`-Nennung ist Inline-Code des Ur-Pfads (Z. 37, Tombstone 1). Die Aufzählung summiert so auf 17, die genannte und gemessene Zahl 16 stimmt.
- **verifizierbar:** ja (Link-Verteilung oben)
- **klasse:** botschaft-aufzaehlung-unvollstaendig

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Vollständigkeit lebender Stellen | geprüft, ohne Befund. Kein Link, keine Workflow-/Skript-/Dependabot-Stelle und keine lebende Doku nennt den Zwischen-Pfad; die fünf Stellen ohne Doku-Gate (`dependabot.yml`, `release.yml` ×2, `hub-description.yml`, `tools/image-test.sh`) zeigen auf `docs/maintainer/releasing.md`. |
| Plan- und Botschafts-Zahlen | geprüft — 40 = 23 (11 `repo-escape` + 12 `target-missing`) + 16 + 1 stimmt; Aufzählung siehe INFO-1. |
| Relative Tiefen | geprüft, ohne Befund. 16 eingehende und 23 eigene Ziele lösen auf dieselbe Datei auf wie vorher (Handbuch `../maintainer/`, ADRs `../../../docs/maintainer/`, Wellen `../../../maintainer/`). |
| Frozen-Klassen (MR-070) / Accepted-ADRs (ADR-0103) | geprüft, ohne Befund. ADR-0014 (4), ADR-0067 (1), welle-73/74/82, CHANGELOG-Altabschnitt, Handbuch-§11-Zeile: nur Link-Ziele; `make adr-check` über die Range 0 Befunde. |
| Tombstone-Deckung §3.6 | geprüft, ohne Befund. ADR-0025 trägt „Tombstone-Register bewusst entfernter/historischer Artefakte, deren Pfad eingefrorene Doku noch zitiert"; die Nennung steht in `done/` (eingefroren), der Eintrag wirkt nachweislich (Bruchprobe: ohne ihn genau der eine `codepath-missing`). |
| `matrix.exempt-paths` | geprüft, ohne Befund. Eintrag nachgezogen; keine weiteren `docs/user`-Globs in `.d-check*.yml`, `Makefile`, `.github/`, die die Datei durch den Umzug verlöre. |
| Commit-Zerlegung §3.3 / Rename | geprüft, ohne Befund. `R100`, bewegte Datei im Move-Commit unverändert; das Mitreisen der eingehenden Verweise ist diesmal in Plan §1 vor dem Code deklariert (vgl. R1 slice-268 LOW-1). |
| Frischer Klon (Risiko §6) | geprüft, ohne Befund. `doc-check`, `adr-check`, `trace-check` im Klon je 0 Befunde. |
| Kommentare in Workflows/Skripten | geprüft, ohne Befund. Nur der Pfad in bestehenden Rang-Zeigern bzw. Meldungstexten geändert. |
| Abgrenzung §1 | geprüft, ohne Befund. Nur `releasing.md` in `docs/maintainer/`; slice-268 und seine Reports unberührt; Inhalt der Datei außer Link-Zielen unverändert. |

## Kategorie-Summary

HIGH 0 · MEDIUM 1 · LOW 2 · INFO 1

## Verdikt

MEDIUM-1 blockiert typischerweise: der neue Ort hat in der Source Precedence keinen
Rang, und Plan §1 schließt die Frage nicht ab — vor der Closure annehmen (Plan-Änderung)
oder begründen. LOW-1 und LOW-2 annehmen oder begründen; INFO-1 ohne Handlungsbedarf.
