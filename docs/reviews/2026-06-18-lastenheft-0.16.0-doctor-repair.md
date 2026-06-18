# Review — Lastenheft-CR 0.15.0/0.16.0 (`--doctor` + `--repair`)

## Kopf-Metadaten

- **Datum:** 2026-06-18
- **Gegenstand:** Working-Tree-Diff `spec/lastenheft.md` (Change Request,
  Dok-Version 0.14.0 → 0.16.0) — zwei neue Anforderungen DC-FA-CLI-007
  (`--doctor`) und DC-FA-CLI-008 (`--repair`) plus Mit-Schärfung
  DC-FA-CLI-003/DC-FA-CLI-004.
- **Reviewer:** Claude (Agent), Skill `.harness/skills/reviewer.md` v1.0.0.
- **Eingangs-Kontext:** Diff (uncommittet, vor Commit B); betroffene
  Anforderungen DC-FA-CLI-003/004/007/008; Hard Rules DC-QA-02
  (Determinismus), DC-QA-03 (Seiteneffektfreiheit/Netz); conventions
  §Anforderungs-Anlege-Prozess; ADRs: keine (CR berührt keine
  Architektur-Entscheidung — DC-QA-03 bleibt unverändert, kein ADR
  erforderlich).
- **Lauf:** Findings im selben Sitzungs-Revisionsstand erhoben und
  aufgelöst; Verdikt bezieht sich auf den finalen Stand (alle MEDIUM/LOW
  eingearbeitet). `make gates` grün.

## Findings

### MEDIUM-1 — Ausgabeformat-Vertrag nicht mit-geschärft

- **Quelle:** DC-FA-CLI-004 · **Pfad:** §3 DC-FA-CLI-004 / DC-FA-CLI-007/008
- **Befund:** DC-FA-CLI-004 deklariert stdout als zeilenweise Befunde
  bzw. (mit `--json`) JSON. `--doctor`/`--repair` schreiben ein weiteres
  stdout-Format, das CLI-004 nicht kannte; CLI-003 war geschärft, CLI-004
  nicht (asymmetrisch). Zusätzlich war `--doctor --json` / `--repair
  --json` undefiniert.
- **Verifizierbar:** ja (AK-Lücke; ein Kombinationsaufruf hatte kein Soll).
- **Resolution:** CLI-004 um die Ausgabe-Modi ergänzt; Modi untereinander
  und mit `--json` nicht kombinierbar → Nutzungsfehler exit 2 (nach
  CLI-003); JSON-Varianten out of scope.

### MEDIUM-2 — „Negative"-AKs waren Property-Checks statt Fehlerpfad

- **Quelle:** DC-FA-CLI-007, DC-FA-CLI-008 · **Pfad:** je Negative-AK
- **Befund:** Beide Negative-AKs prüften Seiteneffektfreiheit
  (read-only-Mount) — dieselbe Eigenschaft wie Beschreibung + DC-QA-03 —
  statt eines echten Fehler-/Missbrauchspfads (exit 2), wie ihn
  CLI-001/003/005/006 führen. Neuer öffentlicher Vertrag ohne
  Negativtest des Fehlerpfads.
- **Verifizierbar:** ja.
- **Resolution:** Negative-AKs auf exit-2-Pfade umgestellt
  (`--doctor --json` bzw. ungültige Stufen-Wahl/`--repair --json`), mit
  angehängter read-only-Note (Muster wie CLI-006).

### MEDIUM-3 — Boundary-AG überzog die breite Reparatur-Stufe

- **Quelle:** DC-FA-CLI-008 · **Pfad:** §3 DC-FA-CLI-008 Boundary
- **Befund:** Die Boundary behauptete einen Best-Guess-Hunk für
  `span-unclosed`; die Beschreibung etabliert Best-Guess nur für
  `target-missing`. Für `span-unclosed` (Schließposition mehrdeutig) war
  keine Reparatur zugesichert — AK setzte nicht zugesicherte Wirkung
  voraus.
- **Verifizierbar:** ja (AK potenziell nicht falsifizierbar).
- **Resolution:** Boundary nutzt jetzt `target-missing` (echt
  best-guess-fähig): konservativ leerer Patch, breite Stufe ein
  markierter Hunk.

### INFO-4 — „review-pflichtig markiert" vs. `git apply`-Kompatibilität

- **Quelle:** DC-FA-CLI-008 · **Pfad:** §3 DC-FA-CLI-008 Beschreibung
- **Befund:** Marker im Patch stehen in Spannung zu einem strikt
  `git apply`-kompatiblen unified diff.
- **Verifizierbar:** teils (erst in spezifikation/Slice).
- **Resolution:** Kennzeichnung läuft über stderr (wie Diagnose,
  CLI-004); der Patch auf stdout bleibt `git apply`-rein. Mechanik-Detail
  in spezifikation.md des Slice.

### LOW-1 — Zwei Anforderungen in einem Minor

- **Quelle:** conventions §Anforderungs-Anlege-Prozess · **Pfad:** §7 / Version
- **Befund:** Bündelung beider Anforderungen in 0.15.0 wich vom
  Ein-Minor-je-Anforderung-Muster ab; keine getrennte Provenance für
  `--repair`.
- **Resolution:** Versions-Split 0.15.0 (`--doctor` + CLI-003/004-Schärfung)
  / 0.16.0 (`--repair`); 1:1 zu slice-025 / slice-026.

### INFO-5 / INFO-6 (akzeptiert)

- INFO-5: §7 verweist auf noch nicht angelegte slice-025/026
  (Forward-Provenance, bare Token, kein Broken Link; Planung folgt).
- INFO-6: CLI-007 Happy-AG nennt jetzt explizit „aktivem Modul `ids`"
  (Präsupposition sichtbar gemacht).

## Negativbefunde (geprüft, ohne Befund)

- **Anker/Link-Politik** (`link-policy: always`): make gates grün; alle
  DC-Querverweise verlinkt, neue Anker lösen auf.
- **Matrix/Referenzrichtung:** keine Abwärts-Links Spec → ADR/Slice; nur
  intra-Spec-Anker + §7 (ausgenommen).
- **DC-QA-02/03-Zitate:** Definitionen decken die Aussagen wörtlich
  (byte-identisch / nie schreiben, kein Netz außer `external`).
- **ID-Vergabe/Struktur:** Bereich CLI vorhanden, CLI-007/008 lückenlos
  nach 006; je drei AK + Out-of-Scope + Versions-Bump + Historie.
- **Kernvertrag:** beide Modi stdout-only/read-only, kein In-place —
  keine Harness-Lüge, kein Vertrags-Widerspruch zur Decke.

## Kategorie-Summary

| Stand | HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|---|
| Erhebung | 0 | 3 | 1 | 3 |
| Nach Resolution | 0 | 0 | 0 | akzeptiert |

## Verdikt

**Freigegeben.** Keine HIGH; alle MEDIUM und das LOW wurden im selben
CR-Stand eingearbeitet, INFO-Punkte akzeptiert bzw. an die
spezifikation.md des Umsetzungs-Slice delegiert (INFO-4). `make gates`
grün. Changeset B ist commit-reif; die Umsetzung folgt über welle-15
(slice-025 `--doctor`, slice-026 `--repair`).
