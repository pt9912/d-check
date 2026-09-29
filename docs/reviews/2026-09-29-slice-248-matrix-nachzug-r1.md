# Review R1 — slice-248: matrix-Klassen `aussen` + `adaptionsblock` adoptieren, Bestands-Nachzug

- **Review-Art:** Code — geprüft gegen den Slice-Plan (`slice-248`), `ADR-0097`/`ADR-0047`,
  `MR-006`/`MR-025`/`MR-032`, Hard Rules `AGENTS.md` §3.4/§3.5/§3.6/§5 (Maintainability;
  DoD-Abhakung ist nicht Gegenstand dieses Reviews)
- **Gegenstand:** `slice-248` · Range `5a45b0b2..HEAD` — `42e264c3` (Claim, reiner
  Move), `95f40dc1` (Implementierung: `.d-check.yml`, Bestands-Nachzug in drei
  Spec-Straten, `ADR-0097`, ADR-Index-Zeile)
- **Skill:** `reviewer.md` @ 1.16.0
- **Modell-ID:** glm-5.3-flash (Z.ai)
- **Datum:** 2026-09-29
- **Eingangs-Kontext:** Slice-Plan (`slice-248`); die Messgrundlage
  (`slice-246`-Probe: 40 Befunde — 28 `aussen`-Links, 1 nacktes MR-Token,
  11 `matrix-inactive`); `ADR-0097` (neu) und `ADR-0047` (Vorgeschichte der
  Historie-Prüfung); `MR-006` (Referenz-Richtung), `MR-025` (Spiegel),
  `MR-032`; Hard Rules `AGENTS.md` §3.4/§3.5/§3.6, §5 Regel 15; Baseline
  `v6.13.0` · `regelwerk/grundlagen-referenz-richtung.md` §Referenz-Richtung
  (SDP), Regel 5; vorherige Findings am Gegenstand: `slice-247` R1 (M-1:
  commits-Spiegel), `slice-246` Closure (Probe-Ordnung)

---

## Findings

### M-1 — Botschaft und ADR behaupten 24 entfernte Verweise; entfernt wurden 11

- **Kategorie:** MEDIUM
- **Quelle:** `AGENTS.md` §5 Regel 15 · Prüffrage 8 (`BEO-ALL/commit-message-overclaims-work`)
- **Pfad:** `docs/plan/adr/0097-matrix-aussen-adaptionsblock-historie-status-ausnahmen.md:59-62`
  (§Entscheidung) und Commit-Botschaft `95f40dc1`
- **Befund:** Beide Stellen behaupten, „13 MR-Provenance-Verweise,
  6 Agenten-Briefing-Verweise, 2 Baseline-Kopf-Zitate und 3 Einzelfälle"
  (Summe 24) seien „aus den Straten entfernt" bzw. „werden … bereinigt".
  Der Diff weist genau **11** Link-Entfernungen aus (3 Konventionsspeicher:
  `MR-032`/`MR-013`/`MR-022`; 4 `AGENTS.md`; 2 Baseline-Köpfe; 2 Einzelfälle:
  Harness-README, Packaging) — die Gegenprobe bestätigt die Zahl (Negativbefund
  unten). Die übrigen **17** gemessenen Links stehen unverändert in den
  §7-Historie-Abschnitten und sind per `exclude-sections` **ausgenommen**, nicht
  entfernt — ein Zustand, den dieselben Stellen an anderer Stelle korrekt
  „Zeitdokument" nennen, hier aber als Entfernung verbuchen. Auch die
  Größenangaben passen so nicht: die 13 MR-Verweise umfassen rechnerisch 3
  entfernte plus 12 in der Historie stehengebliebene; ein „Einzelfall" der
  Drei findet im lebenden Text kein Pendant (CR ×2, Carveout, Register stehen
  alle in der Historie). Die erreichte Sachlage (alle 40 Befunde entschieden,
  Probe grün) ist tragfähig — die behauptete Verteilung der Arbeit ist es
  nicht, und sie ist in eine `Accepted`-ADR eingefroren.
- **Verifizierbar:** ja — Vorher-Baum (`42e264c3`) gegen die neue
  `.d-check.yml` ⇒ 11 Befunde, Exit 1; Ist-Baum ⇒ 0 Befunde, Exit 0.
- **Klasse:** `commit-message-overclaims-work`

### M-2 — ADR-0097 kehrt ADR-0047 um, ohne die Vorgeschichte zu nennen (§3.5)

- **Kategorie:** MEDIUM
- **Quelle:** `AGENTS.md` §3.5 (ADRs sind nach `Accepted` immutable)
- **Pfad:** `docs/plan/adr/0097-matrix-aussen-adaptionsblock-historie-status-ausnahmen.md:43-50`
  (Entscheidung 2)
- **Befund:** Entscheidung 1 von `ADR-0047` lautet: „`matrix.exclude-sections`
  auf `[Geschichte]` verengen — `Historie` und `7. Historie` entfallen. Das
  prüft ab sofort die Spec-§7-Historie". Entscheidung 2 von `ADR-0097` nimmt
  genau `7. Historie` wieder in `exclude-sections` auf (`v6.13.0` ·
  `.d-check.yml:614`) — die Richtung der Bestands-Entscheidung wird
  ins Gegenteil gekehrt. `ADR-0097` trägt weder ein `Supersedes: ADR-0047`
  noch einen Bezug, und `ADR-0047` wird in ihrem Text nicht erwähnt — im
  Unterschied zur eigenen Haus-Form (`ADR-0047` selbst verhandelt ihr
  Verhältnis zu `ADR-0022` ausdrücklich). `make adr-check` fängt das nicht
  (er bewacht den Körper der Bestands-ADR, nicht die Beziehung der Neuen);
  die Konflikt-Pflicht aus §3.5 bleibt dadurch unbehandelt.
- **Verifizierbar:** nein (Form-/Beziehungslücke; `make adr-check` bleibt grün).
- **Klasse:** `adr-nachfolger-ohne-vorgeschichte`

### M-3 — Die Historie-Ausnahme ist eine Baseline-Abweichung, die der Konventionsspeicher nicht führt

- **Kategorie:** MEDIUM
- **Quelle:** Source Precedence / Adaptions-Block-Disziplin (`harness/conventions.md` §Adaptions-Block)
- **Pfad:** `.d-check.yml:614` (`exclude-sections: [Geschichte, "7. Historie"]`)
- **Befund:** Die gepinnte Baseline trägt die Historie-Prüfung ausdrücklich:
  `v6.13.0` · `regelwerk/grundlagen-referenz-richtung.md` §Referenz-Richtung
  (SDP), Regel 5 — „Kein Spec-Dokument nennt eine ADR oder einen Slice, in
  keinem Abschnitt, auch nicht in seiner Historie" — und begründet sie mit
  genau dem Szenario, das hier eingerichtet wird: „kein Gate meldet es, wenn
  die Sektion von der Prüfung ausgenommen ist". Die neue Ausnahme gilt
  global über alle Klassen (auch `adr`/`slice`/`welle`/`commit-hash`), hebt
  damit die Baseline-Enforcement in der Historie auf und ist eine Abweichung
  von der kanonischen Regel — aber `harness/conventions.md` (Adaptions-Block,
  Ort der `MR-<NNN>`-Deklarationen) wird vom Commit nicht berührt, und
  `ADR-0097` führt die neue Begründung (Zeitdokument, CR-Zulässigkeit nach
  v5.11.0), ohne die Unreparierbarkeits-Argumentation der Baseline zu
  widerlegen oder die Abweichung zu deklarieren. Die Referenz-Richtung der
  Spec-Straten ist damit „vollständig mechanisiert" (ADR-Konsequenzen) nur
  unter einer Ausnahme, die als solche nirgends geführt wird.
- **Verifizierbar:** nein (Deklaration; keine Maschine hält
  Baseline-Regel ↔ Config gegeneinander).
- **Klasse:** `baseline-abweichung-undeclariert`

### L-1 — ADR-0097-Kopf ohne `**Autor:**`-Feld

- **Kategorie:** LOW
- **Quelle:** Nutzer-Vorgabe 2026-09-27 (Autor-Feld auch bei ADRs, Wert `pt9912`)
- **Pfad:** `docs/plan/adr/0097-matrix-aussen-adaptionsblock-historie-status-ausnahmen.md:3-8`
- **Befund:** Der ADR-Kopf trägt Status, Datum, Betroffen — aber kein
  `**Autor:**`-Feld; unmittelbare Vorgängerin `ADR-0096` (2026-09-27) führt es.
  Der Feldbestand ist Deklaration, kein Sensor hält ihn.
- **Verifizierbar:** nein (Deklaration).
- **Klasse:** `adr-kopf-feld-fehlt`

### L-2 — Slice-Plan §1: Ausschluss-Bullet endet abgebrochen

- **Kategorie:** LOW
- **Quelle:** Maintainability (Slice-Plan-Form §1 — je Punkt mit Begründung)
- **Pfad:** `docs/plan/planning/in-progress/slice-248-matrix-aussen-adaptionsblock.md:58-60`
- **Befund:** Der dritte Ausschluss („Status-Prüfungs-Weitung") endet nach
  „Die Schwester nimmt den ADR-Index, `docs/reviews/**` und `done/welle-*.md`
  in status.exempt-paths —" mitten im Satz; der Gegengedanke (hiesige
  Entscheidung: je §3.6-Fall per ADR statt Schwester-Freifeld) ist nicht
  ausgeschrieben. Die Begründung des Bullets ist dadurch unvollständig,
  obwohl die Implementierung sie einholt (`ADR-0097`, Entscheidung 3).
- **Verifizierbar:** nein (Prosa).
- **Klasse:** `plan-satz-abgebrochen`

### L-3 — ADR-Index-Zeile: Selbst-Beleg statt Bezugs-Kennung

- **Kategorie:** LOW
- **Quelle:** Maintainability (ADR-Index-Form, Dokumentations-Regel 4)
- **Pfad:** `docs/plan/adr/README.md:107`
- **Befund:** Die Belege-Spalte der neuen Index-Zeile verlinkt die
  `ADR-0097` selbst — ein Beleg, der die Zeile, die ihn trägt, bereits in
  Spalte 1 führt. Die reale Bezüge (die ADR nennt `MR-006` und mechanisiert
  `DC-FA-MTX-003`s Modul) fehlen; der Bestand führt Belege entweder als
  Fremd-Kennung oder als `—` (vgl. `ADR-0067`), nie als Selbstverweis.
  Wer den Index nach der DC-Bindung der Matrix-Senkung befragt, findet eine
  gefüllte Spalte ohne Inhalt.
- **Verifizierbar:** nein (Prosa/Index-Form).
- **Klasse:** `index-selbstbeleg`

### I-1 — ADR-0097 trägt eine fremde Grenze (shallow-Klon) in ihren Konsequenzen

- **Kategorie:** INFO
- **Quelle:** Maintainability (Artefakt-Grenze)
- **Pfad:** `docs/plan/adr/0097-matrix-aussen-adaptionsblock-historie-status-ausnahmen.md:70-73`
- **Befund:** Die Konsequenzen wiederholen die Shallow-Klon-Grenze der
  leere-Range-Prüfung (`DC-FA-VCS-002.a`) — ein Vertragsgegenstand des
  Vorgänger-Slices, der dort bereits deklariert ist
  (`spec/lastenheft.md:2212` Out-of-Scope „auch bei nicht-leerer Range;
  Abhilfe ist `fetch-depth: 0`") und mit der Matrix-Adoption nichts teilt.
  Die Doppel-Deklaration lebt im immutablen Artefakt: ändert sich der
  Erreichbarkeits-Walk, verrottet die Stelle, ohne dass ein Gate sie liest.
- **Verifizierbar:** nein (Prosa).
- **Klasse:** `grenze-doppelt-im-fremden-adr`

---

## Negativbefunde (geprüft, ohne Befund)

- **Probe produktiv grün:** `d-check` mit `--enable matrix --enable ids`
  (übrige Module `--disable`) gegen den Ist-Baum: **915 Dateien, 0 Befunde,
  Exit 0**. Die Config ist adoptiert, nicht mehr Probe.
- **Regeln sind scharf, das Grün ist nicht Inertia:** derselbe Lauf gegen den
  Vorher-Baum (`42e264c3`) mit der neuen `.d-check.yml` (Worktree-Kopie):
  **11 Befunde `matrix-forbidden` (spec-straten → aussen), Exit 1** — genau
  die 11 lebend-entfernten Links; ohne sie wäre der Slice nicht grün.
- **Alle 40 gemessenen Befunde entschieden:** 11 entfernt (Diff), 17
  Historie-Links durch `exclude-sections: [Geschichte, "7. Historie"]`
  ausgenommen, 1 nacktes MR-Token in der Historie (lebender Text token-frei —
  die 11 Befunde der Gegenprobe sind alles Links), 11 Status-Fälle per
  `matrix.exempt-paths` (ADR-0097) ausgenommen. Restmenge grün (Probe oben).
- **exclude-sections heading-genau:** `## 7. Historie` trägt Lastenheft
  (Zeile 3967) und Spezifikation (Zeile 3545); keine ADR führt das Heading
  (geprüft: `grep -l "^## 7. Historie" docs/plan/adr/*.md` leer), die
  ADR-`Geschichte` bleibt über `[Geschichte]` abgedeckt.
- **exempt-paths file-weit ist kein Verlust (ADR-0097-Behauptung hält):**
  `CHANGELOG.md` (slice-Token als Prosa), `docs/plan/adr/README.md`
  (`welle-80`) und `docs/user/releasing.md` führen als FROM-Dokumente der
  Klasse `aussen` keine verbotene Regelkante — ihre einzigen
  Matrix-Treffer waren die 11 Status-Befunde.
- **§3.4 lebender Text:** kein einziger aussen-Link vor §7 in Lastenheft und
  Spezifikation, Architektur-Sicht komplett frei (grep + Probe). Die
  verbleibenden relativen Links beider Straten stehen ausnahmslos in der
  Historie.
- **Baseline-Kopf-Zitate in Zitier-Form:** `spec/spezifikation.md:7-8` und
  `spec/architecture.md:5-6` nennen `Baseline-Regelwerk v6.13.0,
  modul-03-spec.md §…` als Text mit Version statt Link; keine verwaiste
  `d-check:cite`-Direktive (die Straten führen keine).
- **Editorial-Entfernungen ohne Sinn-Verlust:** alle 11 Stellen paraphrasieren
  das Ziel im Kontext und erhalten die §-Nummern („Agenten-Briefing §3.3/§3.8");
  die MR-032-Herkunft bleibt in der Lastenheft-Historie auflösbar (Zeile 4012).
- **MR-025-Spiegel (slice-247 M-1) intakt:** die Korrektur in
  `spec/spezifikation.md:1989-1998` („die **commits-eigene** Range-Semantik
  … der commits-Vertrag, der von `vcs` bewusst **nicht** geteilt wird") ist
  vom slice-248-Commit unberührt; die letzte Berührung liegt vor dem
  Review-Range-Anfang.
- **Rückführungs-Trigger §4 (Urteilsfrage):** nicht ausgelöst, und die
  Teilmenge wäre die schlechtere Entscheidung gewesen — `aussen` zieht die
  Status-Weitung nach sich und erzwingt damit die ADR unabhängig von
  `adaptionsblock`; ohne `aussen` bliebe die Träger-Klasse der 28 Befunde
  unmechanisiert. Der Abschluss in einer Sitzung belegt den Zuschnitt.
- **Klassen-Ordnung und Regeln (DoD-Form):** `adaptionsblock` (Zeile 558) vor
  `aussen` (Zeile 568, `paths: ["**"]`, letzte Klasse); vier Regeln
  `spec-straten`/`sicht` → `aussen`/`adaptionsblock` je `allow: false`
  (Zeilen 572-575); Begründungskommentare am eigenen Bestand, ADR-0097 als
  auflösbares Herkunfts-Feld — §3.7-Form gewahrt.

---

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 (M-1, M-2, M-3) |
| LOW | 3 (L-1, L-2, L-3) |
| INFO | 1 (I-1) |

## Verdikt

**MEDIUM blockiert typischerweise** — hier: die drei MEDIUM blockieren die
Closure, nicht die Mechanik. Die Adoption selbst ist sauber gearbeitet und
belegt: die Klassen-Ordnung entspricht der slice-246-Lektion (First-Match,
`aussen` zuletzt), die vier Regeln feuern nachweisbar (Gegenprobe: 11 Befunde
am Vorher-Baum), die 40 gemessenen Befunde sind alle entschieden und die
Restmenge ist produktiv grün — die Ausschluss-Mechanik ist heading-genau und
ihr file-weiter Zuschnitt verliert keine zweite Prüfung. Was blockiert, ist
die **Buchführung**: die Befund-Verteilung in Commit-Botschaft und
`ADR-0097` (24 entfernt statt 11 entfernt + 17 ausgenommen), das fehlende
Verhältnis zu `ADR-0047` und die undeklarierte Baseline-Abweichung in
Regel 5. Das sind keine Code-Reworks — sie sind Deklarations-Arbeit vor der
Closure, und M-2/M-3 müssen vor dem `git mv` entschieden sein, weil
`ADR-0097` mit `Accepted` in denselben Commit gefroren ist, der sie noch
nicht trägt.
