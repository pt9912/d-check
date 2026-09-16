# Review-Report: slice-221 — 2026-09-16 (R1)

**Review-Art:** Plan-/Closure-Review — geprüft wird die "Gegenstand
entfallen"-Argumentation, der neue Carveout `CO-002`, der Folge-Slice
`slice-225` und die Closure-Vorbereitung von `slice-221` gegen
`modul-05-planning-harness.md` (§Ein Slice, dessen Gegenstand ein anderer
übernimmt), `modul-07-carveouts.md` (§Ziel-Form: Carveout) und `AGENTS.md`.

**Gegenstand:** `docs/plan/planning/in-progress/slice-221-agents-md-tabellenzellen.md`
(vor dem `git mv` nach `done/`), `docs/plan/carveouts/CO-002-slice-221-gegenstand-entfallen.md`,
`docs/plan/planning/open/slice-225-gegenstand-entfallen-uebernommen.md`,
`.d-check.closure.yml`, `docs/plan/carveouts/README.md`. Commit-Kette
`99f1ca88..afdba6bc` (`99f1ca88`, `58bda448`, `ebb912d4`, `afdba6bc`).

**Skill:** `.harness/skills/reviewer.md` @ v1.16.0 / 2026-09-07 ·
**Modell:** `claude-sonnet-5` · **Datum:** 2026-09-16

> **Zitier-Form.** Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb Kennung statt Adresse: `slice-NNN` statt Lifecycle-Pfad,
> `make <target>` als Token statt Link auf die Sensor-Datei, eine
> Baseline-Stelle als Tag plus Pfad in Inline-Code
> (`v6.9.0` · `regelwerk/<datei>.md` §<Abschnitt>).

**Eingangs-Kontext:**

- Slice-Plan `slice-221` (Fassung `in-progress/`, Stand des Prüfzeitpunkts)
- `CO-002` (neu), `slice-225` (neu, in `open/`)
- `MR-072` (Baseline-Pin `v6.9.0`, Delta-Punkt 1: vierter
  Slice-Lifecycle-Zweig), `MR-049` (Ausgangs-Wortschatz), `MR-013`
  (aufgelöst — Lifecycle-Move-Bündelung)
- `DC-FA-TGT-001`, `DC-FA-STRUCT-001`
- `AGENTS.md` §3.3 · §3.7 · §5 · §6
- Baseline `v6.9.0` · `modul-05-planning-harness.md` §Ein Slice, dessen
  Gegenstand ein anderer übernimmt · `modul-07-carveouts.md` §Ziel-Form:
  Carveout
- Beobachtungs-Register: `BEO-ALL/semantic-change-body-only-edges-stale`
  (zum Prüfzeitpunkt 17 Evidence-Dateien), `BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`
  (zum Prüfzeitpunkt 5 Evidence-Dateien, bestätigt durch
  `docs/reviews/2026-09-16-slice-224-baseline-v690-review-r1.md`)
- Precedent: `CO-001` (`docs/plan/carveouts/CO-001-vcs-range-stiller-skip.md`)
  als bestehender Carveout ohne Register-Route für seinen Fund

---

## Verifikationen dieses Laufs

- `make doc-check`: **788 Datei(en) geprüft, 0 Befund(e).**
- `make gates`: grün — `baseline-verify + workflow-pins + doc-check + lint +
  test + arch-check + coverage-gate + semgrep + gate-consistency +
  planning-check`, Coverage 94,60 % (Schwelle 93 %), Semgrep 0 Findings.
- `make verify-closure-notes` gegen den aktuellen Stand: **684 Datei(en)
  geprüft, 0 Befund(e)** (slice-221 liegt noch in `in-progress/`, die
  `done/`-Scan-Menge trifft es deshalb noch nicht).
- **Empirische Probe der `CO-002`-Kernbehauptung**, gegen eine Datei-Kopie
  (kein Git-Repo, nicht committet) mit `slice-221` nach
  `docs/plan/planning/done/` verschoben:
  - **Mit** dem `exempt-paths`-Eintrag: kein `section-tasks-open` für
    `slice-221` (zwei unabhängige Störsignale aus der unvollständigen
    Simulation selbst — fehlender Review-Report, fehlender
    Ruhe-Marker-Nachzug — sind Artefakte des Testaufbaus, keine echten
    Befunde).
  - **Ohne** den Eintrag: **fünf** `section-tasks-open`-Treffer, exakt auf
    den fünf offenen DoD-Häkchen-Zeilen (156, 163, 168, 173, 177).
  - Damit ist `CO-002`s Auflösungs-Trigger-Beleg ("vorher
    `section-tasks-open`, danach `0 Befund(e)`") **bestätigt**, nicht nur
    behauptet — für genau die Menge der fünf betroffenen Zeilen.

## Findings

### F1 — MEDIUM

- **Quelle:** `BEO-ALL/semantic-change-body-only-edges-stale`
  (dieselbe Klasse, die `slice-221` §8 selbst als "tragenden" Fund für die
  Sub-Area führt).
- **Pfad:** `docs/plan/carveouts/CO-002-slice-221-gegenstand-entfallen.md:9-10`
- **Befund:** Der Feld-Wert `**Geltungsbereich:**` zeigt als **sichtbares
  Link-Label** weiterhin `docs/plan/planning/next/slice-221-agents-md-tabellenzellen.md`,
  während das **Linkziel** (der `href`-Teil) in Commit `ebb912d4` bereits
  auf `../planning/in-progress/slice-221-agents-md-tabellenzellen.md`
  aktualisiert wurde. Der Commit hat den Ziel-Pfad nachgezogen, aber den
  Anzeigetext des Links nicht — Label und Ziel widersprechen sich. Zum
  Zeitpunkt der Carveout-Anlage (`58bda448`) stimmten beide auf `next/`
  überein; die Divergenz entstand erst mit dem nächsten `git mv`.
- **Verifizierbar:** ja — Diff von `ebb912d4` auf `CO-002-*.md` zeigt die
  Ziel-Aktualisierung ohne Label-Aktualisierung; `make doc-check` deckt es
  nicht ab, weil das Modul `links` nur den `href`, nie das Label prüft.
- **Klasse:** `semantic-change-body-only-edges-stale` (Link-Label nicht
  mitgezogen bei Ortswechsel).

### F2 — LOW

- **Quelle:** Maintainability.
- **Pfad:** `.d-check.closure.yml` (Geltungs-Konfiguration-Zeile in
  `CO-002-slice-221-gegenstand-entfallen.md:63`)
- **Befund:** Der `d-check:ignore`-Kommentar zum `exempt-paths`-Eintrag
  begründet sich mit „slice-221 liegt zum Anlegen dieses Carveouts noch in
  `next/`, wandert erst mit ihm nach `done/`" — diese Ortsangabe ist seit
  Commit `ebb912d4` (Move `next/` → `in-progress/`) veraltet; die
  eigentliche Ignore-Rechtfertigung (das Linkziel `done/…` existiert noch
  nicht) bleibt davon unberührt gültig.
- **Verifizierbar:** ja — Vergleich der Kommentar-Formulierung mit dem
  aktuellen Dateipfad `docs/plan/planning/in-progress/slice-221-agents-md-tabellenzellen.md`.
- **Klasse:** stale-Erklärung in einem Ignore-Kommentar (verwandt mit F1,
  anderer Fundort, kein Gate deckt Kommentarinhalte gegen den realen
  Dateistand).

### F3 — LOW

- **Quelle:** `AGENTS.md` §5 (Messmethode-vor-Zahl-Regel), verwandt mit
  `BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`.
- **Pfad:** `docs/plan/planning/in-progress/slice-221-agents-md-tabellenzellen.md:335,341`
  (§8, Block "Vier Einträge sind wirklich einschlägig")
- **Befund:** Zwei der vier zitierten Beobachtungs-Register-Zählstände sind
  gegenüber dem gemergten Stand zum Prüfzeitpunkt veraltet:
  `semantic-change-body-only-edges-stale` wird als "16×" geführt
  (`ls .../evidence/ | wc -l` liefert **17** — die zusätzliche Datei
  `slice-222.md` kam am selben Tag über einen anderen Slice hinzu), und
  `mechanical-id-rewrite-misses-frozen-classes` wird als "3×" geführt
  (tatsächlich **5** Evidence-Dateien; der taggleiche Schwester-Review
  `docs/reviews/2026-09-16-slice-224-baseline-v690-review-r1.md` zitiert
  für denselben Eintrag bereits korrekt "5×"). Beide Zahlen ändern keine
  Schlussfolgerung — beide Einträge sind längst `verkörpert`, ein weiteres
  Auftreten löst keine neue 3×-Schwelle aus —, aber derselbe Abschnitt
  korrigiert im unmittelbar vorangehenden Absatz explizit eine andere
  veraltete Zitierung (`registry-vs-authority-table-drift`, dort korrekt
  als "gestrichen" erkannt) und lässt diese beiden stehen.
- **Verifizierbar:** ja — `ls docs/plan/planning/observations/BEO-ALL/<slug>/evidence | wc -l`
  gegen den aktuellen Stand.
- **Klasse:** veralteter Registerzähler in einer als "durchgegangen"
  deklarierten Sichtung.

### F4 — INFO

- **Quelle:** Maintainability / `.d-check.closure.yml` (`structure`-Regel,
  `forbid-pattern` für Vorlagen-Platzhalter).
- **Pfad:** `docs/plan/planning/in-progress/slice-221-agents-md-tabellenzellen.md:318`
- **Befund:** `**Review-Runde 1:** \<wird nach dem Review ergänzt\>` ist ein
  unaufgelöster Platzhalter, den keines der vier `forbid-pattern`-Literale
  (`\(bei Closure\)`, `wird mit dem Closure-Body gefüllt`, `<…>`,
  `<eingetreten:`) trifft — eine in `.d-check.closure.yml` selbst benannte
  Grenze ("die Alternation ist eine LISTE"). Kein Gate hält den Autor davon
  ab, mit diesem Platzhalter nach `done/` zu verschieben.
- **Verifizierbar:** ja — Regex-Probe der vier Literale gegen die Zeile
  ergibt keinen Treffer.
- **Klasse:** bekannte Gate-Lücke, hier als Erinnerung: vor dem `git mv`
  durch die tatsächliche Zusammenfassung dieses Reviews ersetzen.

## Negativbefunde

- **§1 Gegenstand-entfallen-Argumentation** — geprüft gegen den
  `slice-222`-Diff (`99f1ca88`) und `MR-072`s Delta-Messung: Die Behauptung
  ("§4 ist seither ein Sieben-Zeilen-Zeiger", "keine erneute Änderung an §4s
  Form" gegen `v6.9.0`) ist durch den zitierten Commit und `MR-072` gedeckt,
  kein Befund.
- **CO-002 — sechs Pflicht-Header-Felder** (Status · Datum angelegt ·
  Letzte Prüfung · betroffenes Gate · Geltungsbereich · Folge-Slice) —
  geprüft, alle sechs vorhanden, kein Befund.
- **CO-002 — Begründung ohne "noch nicht geschafft"-Aussage** — geprüft,
  rein technische Herleitung (Baseline-Delta vs. Gate-Config-Lücke), kein
  Befund.
- **CO-002 — Auflösungs-Trigger, beobachtbar und messbar** — geprüft und
  empirisch bestätigt (siehe Verifikationen oben), kein Befund.
- **CO-002 — Gate-Konfiguration nennt die `CO-<NNN>`** — geprüft:
  `.d-check.closure.yml:229-235` trägt einen Kommentarblock mit
  `# CO-002: …` direkt am `exempt-paths`-Eintrag, kein Befund.
- **`exempt-paths`-Eintrag — Zielgerichtetheit** — geprüft: exakt ein
  Pfad (`docs/plan/planning/done/slice-221-agents-md-tabellenzellen.md`),
  kein Glob, keine Überdeckung anderer Slices, kein Befund.
- **§6 — alle drei Risiken mit Ausgang** — geprüft: alle drei tragen
  `Ausgang: entfallen` mit individueller, am verschwundenen Gegenstand
  festgemachter Begründung (kein pauschales Copy-Paste), kein Befund.
- **§7 — drei Paarungen** — geprüft einzeln: (a) Anker korrekt als
  "vakant" geführt (kein `liegt in`-Feld, keine Verkörperung, also
  zurecht kein Paarungs-Gegenstand); (b) `slice-225` existiert als Datei in
  `open/` (verifiziert); (c) keine neue Register-Beobachtung nötig — der
  Fund wurde direkt über `CO-002`+`slice-225` aufgelöst (dieselbe
  Risiko-eingetreten-Route wie bei `CO-001`, die den Zähler-Registerweg
  bewusst umgeht, siehe `modul-05-planning-harness.md` §Offene Risiken
  werden bei Closure aufgelöst) — kein Befund.
- **AGENTS.md §3.3 — Move-Commit-Reinheit** — geprüft über die vier
  relevanten Commits: `99f1ca88` (reiner Inhalt, keine Bewegung),
  `ebb912d4` (reiner `git mv` von `slice-221`, 0 Insertions/Deletions an
  der bewegten Datei selbst — andere Dateien reisen zulässig mit),
  `afdba6bc` (reiner Inhalt, keine Bewegung) — kein Befund.
- **Citation-Direktiven in §8** (`d-check:cite` auf
  `modul-05-planning-harness.md:363-364` und `:369-369`) — Wortlaut gegen
  die Baseline-Datei verglichen, exakte Übereinstimmung, kein Befund
  (zusätzlich durch `make gates`/Modul `citations` gedeckt).
- **Spec-ID-Referenzen** (`DC-FA-TGT-001`, `DC-FA-STRUCT-001`) — beide
  lösen in `spec/lastenheft.md` auf, kein Befund.
- **`docs/plan/carveouts/README.md`-Tabellenzeile für `CO-002`** —
  geprüft: Zeile vorhanden, Zähler "zwei Carveouts" korrekt, Verweise lösen
  auf; abweichende Kurzformulierung gegenüber dem Volltitel ist Paraphrase,
  kein Widerspruch, kein Befund.
- **`slice-225` — Selbstständigkeit als Folge-Slice** — geprüft: eigene
  Abgrenzung (drei Punkte, je mit Begründung), eigene DoD (≤ 3
  Liefer-Punkte), referenziert `CO-002` korrekt als aufzulösenden Carveout,
  kein Befund.
- **Kommentar-Klassen der neuen `d-check:ignore`-Direktiven** — beide
  tragen die Klasse *Grenze* (erklären einen noch nicht auflösenden Link),
  keine Review-Historie, keine Deliberation, kein Befund (Formfrage
  getrennt von F2s Inhalts-Staleness).
- **Netzzugriff / Hexagon-Import / Suppression ohne ADR** — nicht
  einschlägig, reine Planungs-/Konfigurationsänderung ohne Code-Diff, kein
  Befund.
- **Review-Report-Namenskonvention** (`docs/reviews/<datum>-slice-221-…-review-r1.md`,
  Substring `slice-221` im Dateinamen) — geprüft gegen `DC-FA-RVW-001`
  (Modul `reviews`, Substring-Match auf die erste `slice-<NNN>`-Kennung),
  passt, kein Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
| --- | --- |
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 1 |

## Verdikt

**Kein Blocker.** Die "Gegenstand entfallen"-Argumentation ist sauber
belegt, `CO-002` erfüllt die Baseline-Form für einen Carveout (sechs
Pflichtfelder, technische Begründung, konkreter und empirisch bestätigter
Auflösungs-Trigger, Folge-Slice, im Gate sichtbar benannt), der
`exempt-paths`-Eintrag ist eng genug gefasst, und alle drei Risiken sowie
alle drei Paarungen sind korrekt aufgelöst. `make gates`,
`make doc-check` und `make verify-closure-notes` laufen grün gegen den
aktuellen Stand.

Der einzige MEDIUM-Fund (F1: stales Link-Label in `CO-002`s
Geltungsbereich-Feld) ist eine Ein-Zeilen-Korrektur, aber **vor** dem
`git mv` nach `done/` zu beheben — sonst schreibt der Move-Commit einen
dritten, dann endgültig falschen Pfadstand fest (das Label müsste ohnehin
auf `done/` aktualisiert werden). Die beiden LOW-Funde (F2, F3) sollten im
selben Aufwasch mitgezogen werden, blockieren die Closure aber nicht. F4
ist eine reine Erinnerung: den Platzhalter `<wird nach dem Review
ergänzt>` durch die Zusammenfassung dieses Reports ersetzen, bevor der
`git mv` erfolgt — sonst hält ihn kein Sensor auf.

**Empfehlung an den Implementer:** F1 und F4 vor dem `git mv` beheben
(beide betreffen exakt die Zeilen, die der Move ohnehin anfasst); F2/F3
können im selben Commit mitlaufen.
