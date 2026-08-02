# Review-Report: slice-085 (Etappe B — Modul-Delta) — 2026-08-02

**Review-Art:** Plan/Analyse-Review (unabhängiger Frischkontext) — Prüf-Objekt
ist eine *Analyse*, kein Diff. Ein „Fehler" ist: Behauptung ohne Quell-Beleg,
falsche Quell-Zitierung, falsche C/D-Zuordnung, falsche Aussage über den
d-check-Ist, Über-Synthese oder Unter-Deckung.

**Gegenstand:** `docs/plan/planning/in-progress/slice-085-etappe-b-modul-delta.md`
§3 (C-1…C-8 / D-1…D-10), §3.3, §3.4, §4.

**Quellen geprüft gegen:** `.harness/baseline/v5.0.0/regelwerk/` (netzlos vendored),
`.harness/baseline/v5.0.0/templates/`, und der d-check-Ist (`harness/conventions.md`,
`spec/lastenheft.md`, `spec/spezifikation.md`, `.d-check.yml`, `docs/plan/adr/`,
`.harness/skills/reviewer.md`, `AGENTS.md`, `docs/plan/planning/in-progress/roadmap.md`).

**Kontext:** abgenommene Analyse `docs/plan/planning/done/slice-083-regelwerk-v500-migration-analyse.md`
§2.3/§2.4. **Modell:** Opus 4.8. **Datum:** 2026-08-02.

**Verfahren:** `.harness/skills/reviewer.md` (Kategorien, kein Finding ohne
Failure-Szenario, kein Lösungsvorschlag im Befund, REFUTED nur mit Beleg,
Negativbefunde Pflicht). Read-only; keine Repo-Datei außer diesem Report geändert.

---

## Findings

### F-1 — D-5: MR-`Status:`-Feld als „zweite Zustandsquelle" widerspricht §3.4/C-6 und dem Ist

- `kategorie`: MEDIUM
- `quelle`: interne Konsistenz (D-5 vs §3.4 vs C-6) + falsche Ist-Aussage
- `pfad`: `docs/plan/planning/in-progress/slice-085-etappe-b-modul-delta.md:71` (D-5) gegen `:59` (C-6), `:109-112` (§3.4)
- `befund`: D-5 stuft „Slice- **und** MR-`Status:`-Feld" als eine zu entfernende
  zweite Zustandsquelle ein (Handlung D). Für Slices ist das korrekt
  (`modul-05-planning-harness.md:28` verbietet das Slice-`Status:`-Feld; d-check-Slices
  tragen es, z. B. Zeile 3 des geprüften Slice). Für MRs trägt der Befund nicht:
  (a) d-checks aktuelle Einträge `MR-000`…`MR-023` in `harness/conventions.md`
  führen **kein** `Status:`-Feld — es gibt nichts zu entfernen; (b) die vendorte
  Vorlage `.harness/baseline/v5.0.0/templates/harness/conventions/MR-NNN-titel.template.md:21`
  schreibt genau `- **Status:** Accepted` vor, und `grundlagen-harness-dateien.md`
  verortet den *Lifecycle*-Zustand in der Verzeichnis-Position (`:136-138`) — die
  Template-`Status: Accepted` ist ein konstanter Annahme-Marker, **nicht** die
  Lifecycle-Quelle. §3.4 (`:109-112`) sagt selbst, das Status-Feld stamme aus dem
  Template und sei in Etappe C *daraus zu belegen*. D-5 (Feld entfernen) und
  §3.4/C-6 (Feld aus dem Template übernehmen/belegen) geben damit einander
  widersprechende Handlungen für dasselbe MR-Feld.
- `verifizierbar`: ja — Template Zeile 21 + `grundlagen-harness-dateien.md:129`/`:136-138`
  gegen den Feldbestand von `harness/conventions.md`; das Failure-Szenario: Etappe C
  legt MR-Dateien per Template mit `Status: Accepted` an, Etappe D streicht es laut
  D-5 wieder — Rework/Widerspruch. Der Slice-Teil von D-5 bleibt gültig.

### F-2 — D-8: „d-check-ADR-Kopf ohne Trigger-Feld" mischarakterisiert den Ist

- `kategorie`: MEDIUM
- `quelle`: falsche Ist-Aussage; mislozierte Baseline-Anforderung
- `pfad`: `docs/plan/planning/in-progress/slice-085-etappe-b-modul-delta.md:74` (D-8)
- `befund`: D-8 behauptet, d-check-ADRs trügen keinen Re-Evaluierungs-Trigger, und
  begründet die Einstufung „gering" mit „kein Template-Zwang". Beides ist am Ist
  unscharf: 27 von 46 ADRs unter `docs/plan/adr/` tragen eine eigene Body-Sektion
  `## Re-Evaluierungs-Trigger` mit Inhalt (z. B. `docs/plan/adr/0046-sources-upstream-content-drift.md:115`,
  die jüngste ADR). Der Baseline-Ort ist ebenfalls eine **Body**-Sektion, kein
  Kopf-Feld (`.harness/baseline/v5.0.0/templates/docs/plan/adr/NNNN-titel.template.md:94`
  führt ein **Pflicht**-Kapitel „## Re-Evaluierungs-Trigger", `modul-04-adrs.md:28-33`
  die Regel). Damit ist „kein Template-Zwang" für die v5.0.0-Vorlage schlicht falsch,
  und das reale Rest-Delta ist ein anderes: der Block `0032`–`0045` ohne die Sektion —
  alle `Accepted`/immutabel und daher nicht nachrüstbar (grandfathered). Der als D-Handlung
  implizierte Zustand („Trigger-Feld ergänzen") trifft die modernen ADRs nicht.
- `verifizierbar`: ja — `grep -lE '^#+ .*Re-Evaluierungs-Trigger' docs/plan/adr/[0-9]*.md`
  liefert 27 Treffer; Failure-Szenario: Etappe D handelt auf ein weitgehend
  nicht existentes Delta. Bounded durch ADR-Immutabilität und die Selbst-Einstufung „gering".

### F-3 — Unter-Deckung: §2.3-Zugang „AGENTS.md gegen Template angleichen" fehlt in der Finding-Liste

- `kategorie`: MEDIUM
- `quelle`: Unter-Deckung (DoD-Vollständigkeit der §2.3-Zugänge)
- `pfad`: `docs/plan/planning/in-progress/slice-085-etappe-b-modul-delta.md:47` (Spot-Check-Zusage) + §3-Findingliste
- `befund`: slice-085 §1 und DoD (`:143-144`) verpflichten die Analyse, **jeden**
  slice-083-§2.3-Zugang gegen die Quelle zu bestätigen/korrigieren und mit
  C/D-Handlung zu versehen. Der Zugang „AGENTS.md gegen `AGENTS.template.md`
  angleichen" (`done/slice-083-…md` §2.3, kanonisch
  `modul-09-implementierung.md:154-175`/`:186-191` §AGENTS.md-Regeln + §Ziel-Form,
  und expliziter Etappe-D-Schritt in slice-083 §2.7 D3) taucht in keinem der
  18 Findings auf — AGENTS.md erscheint in slice-085 nur in D-7 (`§4`-Gate) und als
  MR-Link. Zugleich behauptet der §3-Kopf (`:47`), die niedrig-priorisierten Module
  inkl. `modul-09` seien spot-gecheckt „keine neue Pflicht über die Liste hinaus".
  Das ist eine Aussage im Modul, das die AGENTS.md-Angleichung kanonisiert. Real
  offen: d-checks Hard Rules `AGENTS.md` §3.1–§3.6 tragen keinen Herkunfts-Anker
  (unkritisch, da meist ADR-/Anforderungs-abgeleitet), und der §1-Kanon-Zeiger auf
  Modul 9 fehlt (Etappe A/`MR-023` zog nur Pin/Layout nach, nicht die Struktur).
  (Sekundär gleicher Klasse: der §2.3-Zugang „Freshness-Audit als eigener
  Konventions-Abschnitt" / `MR-022`-Verschmelzung ist ebenfalls nicht als Finding
  geführt.)
- `verifizierbar`: ja — `grep -n AGENTS docs/plan/planning/in-progress/slice-085-…md`
  zeigt keine Angleichungs-Zeile; Failure-Szenario: Etappe D nimmt die Finding-Liste
  als delta-vollständig und unter-skopt AGENTS.md. **Bounded:** slice-083 §2.7 D3 trägt
  den Auftrag weiter, d. h. das Gesamt-Vorhaben verliert ihn nicht — Maintainer kann
  auf LOW herabstufen, wenn slice-085 nicht als alleinige C/D-Grundlage gilt.

### F-4 — C-7: Regel-Quelle auf `grundlagen-durchsetzungsschicht` statt `modul-13` §Guard-Härtung gelegt

- `kategorie`: LOW
- `quelle`: falsche/ungenaue Quell-Zitierung
- `pfad`: `docs/plan/planning/in-progress/slice-085-etappe-b-modul-delta.md:60` (C-7)
- `befund`: C-7 zitiert als Quelle „`grundlagen-durchsetzungsschicht` §Guard-Härtung:
  ‚Grenz-Zeile mitziehen'" samt Enumeration `python -c`/`env`/`Wrapper`. Die
  ausformulierte Regel „Jeder Wächter-`MR` trägt, was der Wächter *nicht* kann
  (`python -c`, `env`-Umwege, Wrapper-Skripte)" steht kanonisch in
  `modul-13-quality-gates.md:189-192` (§Guard-Härtung, Anker `guard-haertung`).
  `grundlagen-durchsetzungsschicht.md` hat **keinen** Abschnitt „§Guard-Härtung";
  er nennt „Mitziehen der Grenz-Zeile" nur beiläufig (`:86`) und verweist auf
  Modul 13, plus §Grenzen (`:65-75`, `python -c`). Die Regel existiert also (Substanz
  korrekt), aber die Primär-Zitierung deutet auf die falsche Datei/Sektion.
- `verifizierbar`: ja — Datei-Vergleich; kein inhaltliches Fehlurteil, nur
  Fundstellen-Präzision (rottet, wenn Etappe C den Anker-Link setzt).

### F-5 — C-3: D-Komponente lebt nur in §3.3-Prosa, nicht als Zeile in der D-Tabelle (§3.2)

- `kategorie`: INFO
- `quelle`: C/D-Zuordnungs-Form
- `pfad`: `docs/plan/planning/in-progress/slice-085-etappe-b-modul-delta.md:56` (C-3 in der C-Tabelle) / `:88-91` (§3.3)
- `befund`: C-3 wird in §3.1 (Etappe-C-Tabelle) geführt; seine D-Handlung — die
  `slice-NNN`-/ADR-Verweise aus `spec/spezifikation.md` §7 und `spec/lastenheft.md`
  §7 entfernen (berührt die kanonische Spec, ausdrücklicher Abnahme-Punkt) — steht
  nur im §3.3-Fließtext, nicht als eigene D-Zeile in §3.2. Wer §3.2 als
  Etappe-D-Arbeitsliste liest, findet die einzige die Spec-Straten editierende
  Handlung nicht in der Tabelle. §3.3 hebt sie prominent hervor, daher INFO, nicht höher.
- `verifizierbar`: ja — Struktur der §3.2-Tabelle enthält keine C-3-D-Zeile.

### F-6 — §3.4: „= Adaptionen, keine Forks" ist minimal zu weit (Konformität vs. Ersetzung)

- `kategorie`: INFO
- `quelle`: Über-Präzisierung, aber sachgerecht deferiert
- `pfad`: `docs/plan/planning/in-progress/slice-085-etappe-b-modul-delta.md:95-108` (§3.4)
- `befund`: §3.4 korrigiert §2.4 zu Recht: `grundlagen-durchsetzungsschicht.md`
  kanonisiert Content-Hash-Nachweis (Design-Eigenschaft 2, `:50-55`) und
  Guard-Härtung (`:84-87` → `modul-13:183-192`), also ist die §2.4-Aussage
  „von der Baseline nicht geregelt" widerlegt und „Fork" nicht haltbar. Die
  Schlussfolgerung „= **Adaptionen**" greift jedoch eine Stufe zu weit: das Feld
  `Ersetzt-Baseline-Regel` verlangt eine Regel, „an deren Stelle dieser Eintrag
  tritt" (Template `:15-17`); `MR-004`/`MR-005` *implementieren/schärfen* die
  kanonisierte Mechanik, treten aber nicht an ihre Stelle — sie könnten daher
  ebenso reine Konformität (weder Fork noch Adaption) sein. §3.4 hält das offen
  („Etappe C muss das Fork-Kriterium … schärfen **bevor** es reklassifiziert"),
  daher kein Blocker.
- `verifizierbar`: ja — Template-Feldsemantik gegen den MR-Charakter; reine
  Klassifikations-Nuance für Etappe C.

---

## Negativbefunde (geprüft, ohne Befund)

- **C-1 (drei Straten obligatorisch):** belegt — `grundlagen-source-precedence.md:113-114`,
  `grundlagen-referenz-richtung.md:266`, `modul-03-spec.md:81`. Prämisse „Kurs-Default
  = zwei Ränge" von `MR-001` entwertet. Korrekt.
- **C-2 (Spec-Decke = Default):** belegt — `grundlagen-referenz-richtung.md:42-44`,
  `modul-03-spec.md:46-51`; `MR-006` → Provenienz, `matrix`-Mechanismus bleibt. Korrekt.
- **C-3 / §3.3 (Historie-Provenance für Spec-Straten revoziert) — schärfster Fund, hält:**
  wörtlich belegt in `grundlagen-referenz-richtung.md:164-166` (Regel 5) und `:207-214`
  („ohne ausgenommene Sektion" / „gäbe es dort eine ausgenommene Sektion, wäre sie
  genau die Stelle …"), zusätzlich `modul-03-spec.md:53-59` und `:86-90`
  („Kein ADR-Verweis, auch nicht in der Historie" · „wer den auslösenden Slice in der
  Historie nennt, tut dasselbe"). Ist bestätigt: `.d-check.yml:84`
  `matrix.exclude-sections: [Historie, "7. Historie", Geschichte]`; `spec/spezifikation.md`
  §7 (Zeile 2135 ff.) trägt eine `slice-NNN`-Verweis-Spalte **und** ADR-Links
  (10 ADR-/63 slice-Vorkommen), `spec/lastenheft.md` §7 (Zeile 2191 ff.) die
  `slice-NNN`-Spalte. Die Planungs-Klassen behalten die Ausnahme (Regel 5) — korrekt
  differenziert. Doppelbeleg (Reader A/C) plausibel.
- **C-4 (8×8-Matrix, d-check zu eng):** belegt — Matrix `grundlagen-referenz-richtung.md:31-40`;
  `.d-check.yml` `matrix` kodiert nur 3 Klassen/3 Regeln, damit sind
  ADR→{Carveout,Welle,Roadmap} (`:36`) und Slice→Roadmap (`:37`) unbewacht. Korrekt.
- **C-5 (24 Inline-Einträge, Index+Datei je MR):** belegt — 24 `### MR-*`-Einträge
  (`MR-000`…`MR-023`) in `harness/conventions.md` gezählt (Zahl korrekt gegenüber
  slice-083 „23" nach Etappe-A-`MR-023` erhöht); `grundlagen-harness-dateien.md:114-157`.
- **C-6 (neue Pflichtfelder):** belegt — `grundlagen-harness-dateien.md:129`
  (`Ersetzt-Baseline-Regel` Pflicht; `Löst auf`/`Ausgelöst durch Baseline-Stand`
  bedingt; `(schärft …)` im Titel) und Template `:12-13`. `Status` bewusst zu §3.4
  verwiesen — konsistent.
- **C-8 (Producer-Ausnahme entfällt):** belegt — die v5.0.0
  `.harness/baseline/v5.0.0/templates/README.md` §Ein- vs. wiederkehrende (`:72-90`)
  hält wiederkehrende Templates co-located; kein „Self-Hosting/Producer/`harness.mk`"-Abschnitt
  mehr (grep: 0 Treffer). `MR-018`-Brücke damit erloschen. Korrekt.
- **§3.4 Fork-Definition:** belegt — „pauschale Nichtanwendbarkeit … ist ein Fork"
  wörtlich in `grundlagen-source-precedence.md:80-82`. Alle 24 MRs geprüft: keiner
  erklärt die Baseline pauschal für unanwendbar — „keine d-check-Adaption erfüllt das"
  plausibel.
- **§3.4 Status-Korrektur:** belegt — grundlagen-Pflichtfeldliste
  (`grundlagen-harness-dateien.md:129`) ohne `Status`; Template `:21` trägt
  `Status: Accepted`. Die Herkunfts-Zuordnung ist korrekt (siehe aber F-1 zur Folge-Handlung).
- **D-1 (Roadmap 5 vs 3):** belegt — `modul-06-roadmap.md:53-64`;
  `docs/plan/planning/in-progress/roadmap.md` führt nur *Aktuelle Welle*,
  *Nächste Wellen*, *Historische Trigger-Verschiebungen* (Closure-Bestand steckt in
  der `Vorgänger:`-Prosa). Korrekt.
- **D-2 (Wellen-Closure 5 Schritte):** belegt — `modul-06-roadmap.md:111-217`; keine
  flachen `welle-*.md`-Plandateien, keine `done/welle-*-results.md`; beide Vorlagen
  vendored, ungenutzt. Korrekt.
- **D-3 (observations.md/`BEO-`):** belegt — `modul-06-roadmap.md:77-109`,
  `grundlagen-traceability.md:50-58` (Anker-Paarung / Feld `liegt in <Zielort>`);
  `docs/plan/planning/observations.md` fehlt, kein `BEO-` in d-check-Artefakten. Korrekt.
- **D-4 (zwei Vorprüfungen):** belegt — `modul-05-planning-harness.md:161-184`; der
  §6-Bezug passt zum Haus-Stil (`MR-014`: Sub-Area-Modus = §6). Korrekt.
- **D-6 (Report-Kopf + `klasse`):** belegt — Template
  `.harness/baseline/v5.0.0/templates/docs/reviews/review-report.template.md:10/16/17`
  (Review-Art/Skill/Modell), `:46` (`klasse`-Feld), `:69` (Finding-Klassen-Zeile);
  `.harness/skills/reviewer.md` §Output-Schema (Z. 65-71) ohne `klasse`. Korrekt.
- **D-7 (closure-note-reviewer):** belegt — `modul-10-review-harness.md:50-52`,
  `modul-11-verification.md:66`; Template existiert, `.harness/skills/` enthält nur
  `reviewer.md`. Korrekt.
- **D-9 (Carveout-Audit-Slice):** belegt — `modul-07-carveouts.md:96-104`;
  `docs/plan/carveouts/` trägt nur `README.md` (0 aktive Carveouts). Korrekt.
- **D-10 (Reviewer-Skill-Kopf-Drift):** belegt — `.harness/skills/reviewer.md` Z. 4-6
  zitiert `grundlagen-konventionen.md` (entfallen) + „Kurs-Welle 18 §8"; Z. 39 zeigt
  bereits korrekt auf `grundlagen-referenz-richtung.md`. Korrekt.
- **§3.5 Negativbefunde (CR-Fußabdruck, ADR-Immutabilität, Traceability-Kern, Modul 13/15):**
  stichprobenartig geprüft (`modul-03-spec.md:36-66`, `modul-04-adrs.md:35-44`,
  `grundlagen-traceability.md:1-17`) — als konform plausibel, keine Über-Synthese.
- **§4 Flotten-Stand:** **nicht** unabhängig verifizierbar in diesem Checkout —
  u-boot/a-check/ai-harness-init sind Fremd-Repos außerhalb des Baums; die
  `v3.5.2`-/`d-check-first`-Aussage konnte nur intern-plausibel, nicht am Fremdstand
  geprüft werden (kein Widerspruch gefunden, aber offen).

---

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** C/D-Handlungs-Widerspruch bei Template-gedeckten
Feldern (F-1) · Ist-Aussage vs. Artefakt-Bestand (F-2) · abgenommener §2.3-Zugang
nicht in die Finding-Liste überführt (F-3).

---

## Verdikt

**Nicht abnahmereif** (3 MEDIUM blockieren typischerweise; §Verfahren).

Der analytische Kern ist stark und hält der Gegenprobe stand: die höchststakes-Belege
— C-1, C-3/§3.3 (Historie-Provenance-Revocation, der schärfste Fund), C-4 (8×8), C-5/C-6,
C-8 und beide §3.4-Korrekturen — sind am vendorten v5.0.0-Quelltext **wörtlich** belegt
und am d-check-Ist bestätigt. Die drei MEDIUM sind Rand-Findings, keine Kern-Einstürze:

1. **F-1 (D-5, MR-Status):** das `Status:`-Feld für MRs ist template-**gedeckt**
   (`MR-NNN-titel.template.md:21`) und im Ist gar nicht vorhanden — D-5 (entfernen)
   widerspricht §3.4/C-6 (aus dem Template belegen). Der Slice-Teil von D-5 bleibt korrekt.
2. **F-2 (D-8, ADR-Trigger):** die Ist-Aussage „ADRs ohne Trigger" ist für 27/46 ADRs
   (inkl. der modernen Haus-Stil-ADRs) und für den v5.0.0-Template-Zwang falsch.
3. **F-3 (Unter-Deckung AGENTS.md):** der abgenommene §2.3-Zugang „AGENTS.md angleichen"
   ist aus der als C/D-verbindlich deklarierten Finding-Liste gefallen, während der
   `modul-09`-Spot-Check „keine neue Pflicht" behauptet.

Alle drei sind punktuell klärbar (D-5 MR-Teil streichen bzw. auf „per Template führen"
umformulieren; D-8 auf das reale immutable-Alt-ADR-Rest-Delta + fehlendes Trigger-Audit
umschreiben; AGENTS.md-Angleichung als D-Finding nachtragen oder als bewusst nach
slice-083 §2.7 D3 delegiert ausweisen). F-3 lässt sich vom Maintainer auf LOW abstufen,
falls slice-085 nicht als alleinige C/D-Grundlage gilt — dann verbleiben zwei MEDIUM.
LOW (F-4, Quell-Zitierung C-7) und INFO (F-5 C-3-D-Zeile, F-6 Fork/Konformität) sind
Etappe-C-Feinschliff.
