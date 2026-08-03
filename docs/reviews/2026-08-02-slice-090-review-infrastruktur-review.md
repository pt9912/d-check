# Review-Report: slice-090 (welle-67, Etappe D — Review-Infrastruktur) — 2026-08-02

**Review-Art:** Code/Doc-Form — *wogegen*: die Änderung an `.harness/skills/reviewer.md`
und den Planning-Dokumenten gegen die vendored v5.0.0-Baseline-Formen (Ziel-Formen
`review-report.template.md`, `modul-10-review-harness.md`, `grundlagen-referenz-richtung.md`)
und den Slice-Plan.

**Gegenstand:** slice-090, Commit-Range `d7b4d6f..HEAD` (`580dc48`); drei Commits
`91bf130` (open) · `419e82e` (reviewer.md) · `580dc48` (D-7-Entscheid + Roadmap).

**Skill:** `reviewer.md` @ 1.3.0 (Commit `580dc48`) ·
**Modell:** claude-opus-4-8 · **Datum:** 2026-08-02

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-090-review-infrastruktur.md` (Ziel/Vorgehen/Abnahme/DoD)
- Ziel-Formen unter `.harness/baseline/v5.0.0/`: `templates/docs/reviews/review-report.template.md`,
  `regelwerk/modul-10-review-harness.md`, `regelwerk/grundlagen-referenz-richtung.md`
- Bezugs-Slice `slice-085` §3.2 (Etappe-D-Finding-Liste D-6/D-7/D-10)
- Hard Rules (`AGENTS.md` §3); Repo-Konvention „Baseline-Default sticht repo-lokale Adaption"

---

## Findings

### F-1 — Baseline-Modul-Enumeration und Template divergieren beim Feld `klasse`

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `.harness/baseline/v5.0.0/regelwerk/modul-10-review-harness.md:67` vs.
  `.harness/baseline/v5.0.0/templates/docs/reviews/review-report.template.md:46`
- `befund`: Die vom `reviewer.md`-Kopf als Output-Schema-Quelle zitierte Modul-Sektion
  §Ziel-Form: Reviewer-Skill enumeriert fünf Finding-Felder
  (`kategorie · quelle · pfad · befund · verifizierbar`) ohne `klasse`; die Template-Form
  und das nun angeglichene `reviewer.md` §Output-Schema führen sechs Felder inkl. `klasse`.
  Der `klasse`/Finding-Klasse-Begriff ist im Modul-Fließtext (Register-Übergabe) sehr wohl
  verankert, nur die Ein-Zeilen-Enumeration hinkt hinterher. slice-090 hat sich korrekt an
  die Template-Form gehalten (die das Modul selbst als Form des „ganzen Reports" delegiert);
  der Widerspruch liegt in der **vendored Baseline**, die dieser Doc-Form-Slice
  regelkonform nicht anfasst.
- `verifizierbar`: nein — kein Gate vergleicht Modul-Enumeration gegen Template.
- `klasse`: Baseline-Enumeration-lagt-hinter-kanonischer-Form
- **Failure-Szenario:** Ein späterer Pflege-Lauf liest den `reviewer.md`-Kopf, folgt dem
  Zeiger auf §Ziel-Form (Modul-Zeile 67), sieht dort nur fünf Felder und „korrigiert"
  `reviewer.md` durch Streichen von `klasse` — reaktiviert genau die Drift, die D-6 gerade
  beseitigt hat. Auffangbar erst upstream in der Baseline (nicht in slice-090).

## Negativbefunde

- geprüft, ohne Befund: **D-10-Currency vollständigkeit** — `grep grundlagen-konventionen`
  und `grep "Kurs-Welle 18"`/`"Welle 18"` in `.harness/skills/reviewer.md` → je 0 Treffer;
  die neuen Ziele existieren im vendored Baum und tragen die zitierten Abschnitte
  (`modul-10-review-harness.md:44` §Ziel-Form: Reviewer-Skill;
  `grundlagen-referenz-richtung.md:4` §Referenz-Richtung (SDP), Anker
  `referenz-richtung-sdp-wer-darf-wen-referenzieren` löst auf); Version 1.2.0→1.3.0,
  Datum 2026-06-28→2026-08-02.
- geprüft, ohne Befund: **D-6 Abwärtskompatibilität** — der Diff zeigt `klasse` als
  **angehängtes** sechstes Feld hinter `verifizierbar`; kein bestehendes Feld umbenannt
  oder entfernt (Risiko §5 „Felder ergänzen, nicht umbenennen" eingehalten).
- geprüft, ohne Befund: **D-6 §Ablage-Kopf-Metadaten vs. Template** — Feld-für-Feld deckungs-
  gleich mit `review-report.template.md` Kopf: Review-Art · Gegenstand · Skill (@ Version/
  Commit) · Modell(-ID) · Datum · Eingangs-Kontext.
- geprüft, ohne Befund: **Kein kanonisches Duplikat** — §Ablage verweist explizit auf die
  Ziel-Form `review-report.template.md`; die Feldnennung ist Spiegelung (wie das Template
  selbst deklariert), keine konkurrierende Zweitdefinition.
- geprüft, ohne Befund: **D-7 Faktentreue** — kein `check_closure_notes.py` im Repo, kein
  `verify-closure-notes`/`closure-note`-Target im `Makefile`, kein aktives
  `.harness/skills/closure-note-reviewer.md` (nur die vendored `*.template.md`);
  ADR-0011 in d-check ist `0011-digest-pins-build-gate-images.md` — bestätigt die
  Slice-Plan-Aussage, dass die 0011er-Nummer die Digest-Pins trägt. „Fehlt komplett" ist
  faktentreu.
- geprüft, ohne Befund: **Roadmap-Kandidat** — D-7 als „Kandidat (auf Freigabe wartend)"
  unter §Nächste Wellen, widerspruchsfrei zum Slice-Plan-Entscheid (a); die genannten
  Ziel-Formen `closure-note-reviewer.template.md` und `review-report.template.md` existieren
  beide vendored; keine toten Verweise (Selbst-Scan 0 Befunde).
- geprüft, ohne Befund: **Scope-Reinheit** — Range `d7b4d6f..HEAD` berührt nur
  `.harness/skills/reviewer.md`, `roadmap.md` und `slice-090-…md`; keine `Accepted`-ADR,
  keine `spec/`-Datei, kein Go-/Python-Code — konform zur Doc-Form-Disziplin (Risiko §5.2).
- geprüft, ohne Befund: **Slice-Plan-Selbstkonsistenz** — Ziel/Vorgehen/DoD/Abnahme decken
  die Umsetzung; DoD-Punkte D-10/D-6/D-7-Entscheid entsprechen den drei Commits.
- geprüft, ohne Befund: **Gate-Lauf** — offline `docker run --network none … d-check:latest`
  → 318 Dateien, 0 Befunde; Working-Tree clean an `580dc48`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Baseline-Enumeration-lagt-hinter-kanonischer-Form

## Verdikt

**Merge-blockierend:** nein. Kein HIGH/MEDIUM/LOW; das einzige INFO betrifft eine
harmlose interne Inkonsistenz der **vendored** Baseline (Modul-Enumeration vs. Template),
die dieser Doc-Form-Slice regelkonform nicht anfassen darf — für den Upstream-Baseline-Pfad
vermerkt, nicht für slice-090.

**Verdikt: abnahmereif.**

**Übergabe:** slice-090 setzt D-6/D-10 vollständig und abwärtskompatibel um; D-7 ist sauber
als Folge-Produkt-Slice herausgeschnitten und im Roadmap-Backlog verankert. Restpunkt der
DoD (unabhängiger Frischkontext-Review) ist mit diesem Report erbracht; die DoD-Abhakung und
`make gates`-Bestätigung sind Sache der Verifikation (getrennter Kontext).
