# Review — slice-071 Implementierung R3 (Vertrags- + Fix-Review zu R2)

**Datum:** 2026-07-17
**Review-Art:** unabhängiger Vertrags- und Fix-Review (kontext-getrennt;
adversarial, nicht bestätigend)
**Gegenstand:**
[`slice-071`](../plan/planning/open/slice-071-trace-cross-consistency-gate.md) —
`1665714` (doc-first: Lastenheft 0.44.1, `DC-FA-XREF-001.a` Schritt 5,
ADR-0038 Entscheidung 8) + `9027514` (Code folgt dem Vertrag);
Gesamt-Diff `6c4ccf5..HEAD`
**Vor-Reviews:** [R1](2026-07-17-slice-071-implementation-r1.md) (REJECT),
[R2](2026-07-17-slice-071-implementation-r2.md) (REJECT)
**Reviewer:** Claude (kontext-unabhängiger Lauf)
**Skill:** `.harness/skills/reviewer.md` v1.2.0
**Modell:** Opus 4.8 (1M context)

## Eingangs-Kontext

- **Neuer Vertrag** (Rang vor allem anderen):
  [`DC-FA-XREF-001`](../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
  Lastenheft 0.44.1 (§Vakuum-Absatz, Negative-Kriterium „Vakuum",
  Boundary-Kriterium „einseitig leere Sicht", erweiterte fail-closed-Enumeration)
- [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
  Schritte 1–7 + §Fehlerpräzedenz
- [ADR-0038](../plan/adr/0038-trace-cross-consistency.md) Entscheidung 8 +
  §Fitness-Funktion + §Geschichte;
  [ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md)
  (Abgrenzungs-Anker)
- [`MR-006`](../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs),
  [`AGENTS.md`](../../AGENTS.md) §3.4/§3.5/§3.6
- Prüfgegenstand: `internal/hexagon/core/app/trace_cross.go`,
  `internal/hexagon/core/app/trace_cross_test.go`, `spec/lastenheft.md`,
  `spec/spezifikation.md`, `docs/plan/adr/0038-trace-cross-consistency.md`

**Ausgeführte Sensoren:** `make lint` (0 issues), `make test` (grün),
`make arch-check` (0 Befunde), `make doc-check` (222 Dateien, 0 Befunde),
`make adr-check RANGE=d11398f..HEAD` (0 Befunde), `make build` + fünf
Repro-Läufe gegen `d-check:latest` (`--network none`, `:ro`).

## Findings

### R3-F-1 — Die Vakuitäts-Stufe läuft **vor** dem Ausschluss; die Spec-Schrittfolge ordnet sie **danach** — ein alles tilgendes `exclude-req` meldet einen echten Drift als `0 Differenz(en).`

- **kategorie:** MEDIUM
- **quelle:** [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
  Schritte 4/5/6 vs. §Fehlerpräzedenz;
  [`DC-FA-XREF-001`](../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
  (§Vakuum-Absatz, Negative-Kriterium „Vakuum");
  [ADR-0038](../plan/adr/0038-trace-cross-consistency.md) Entscheidung 5 + 8
- **pfad:** `internal/hexagon/core/app/trace_cross.go:109-111`
  (`crossVacuity` vor `diffViews`), `:284-297` (`crossVacuity`), `:311-325`
  (`diffViews` wendet `ExcludeReq` erst intern an); `spec/spezifikation.md:447-462`
- **befund:** Der Vertrag ist an dieser Stelle **in sich widersprüchlich**. Die
  nummerierte Schrittfolge mutiert `F`/`B` in Schritt 4 („Ausschluss …
  werden … aus den Schlüsselmengen von `F` und `B` entfernt") und prüft in
  Schritt 5 „wenn (a) `F` **und** `B` leer sind" — dieselben Symbole, nach der
  Mutation, also **post-Ausschluss**. Die §Fehlerpräzedenz dagegen listet
  „… → Range-Expansion → Vakuität (Schritt 5) → Diff" und kennt **keine**
  Ausschluss-Stufe, was **prä**-Ausschluss impliziert. Der Code implementiert die
  zweite Lesart. Unter der ersten ist ein `exclude-req`, das alle Anforderungen
  tilgt, ein Vakuum (Exit 2); tatsächlich liefert der Lauf `0 Differenz(en).` und
  Exit 0 auf einem realen Drift. Das Ventil ist von
  [ADR-0038](../plan/adr/0038-trace-cross-consistency.md) Entscheidung 5
  ausdrücklich als drift-anfällig markiert („kuratierte Kante … kann driften —
  kein gelöstes Problem"), womit ein zu weit gefasstes Muster (`'GG-'` statt
  `'^GG-SPEC-'`) kein konstruierter, sondern ein benannter Realfall ist.
- **verifizierbar:** ja — reproduziert gegen `d-check:latest`. Repo mit
  Trigger-088-Drift (`F = {GG-AR-COMP-CORE}`, `B = {GG-AR-COMP-SCHED}`,
  Schnittmenge null) und `exclude-req: '.'`:

  ```text
  ## Kreuzverweis-Konsistenz

  0 Differenz(en).
  EXIT=0                       # --trace --require-complete
  ```

  Nach der Schrittfolge-Lesart (4 → 5 → 6) wäre `F' = B' = ∅` und damit Exit 2.

### R3-F-2 — Stale Schritt-Referenz: Schritt 3 verweist auf den „Mengen-Diff (Schritt 5)", der seit `1665714` die Vakuitäts-Prüfung ist

- **kategorie:** LOW
- **quelle:** [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
  Schritt 3 (§Vorbedingung)
- **pfad:** `spec/spezifikation.md:445` („sonst wäre der Mengen-Diff (Schritt 5)
  inhärent leer/voll und bedeutungslos")
- **befund:** Der Einschub der Vakuitäts-Stufe als Schritt 5 hat die Schritte 6/7
  nachgezogen, aber nicht die **Rückwärts-Referenz** in Schritt 3. Sie zeigt
  jetzt auf die Vakuitäts-Prüfung statt auf den Diff (Schritt 6) — die Aussage
  „der Mengen-Diff (Schritt 5)" ist damit sachlich falsch. Schritt 5 selbst
  verweist korrekt auf „der Diff (Schritt 6)" und „(es ist geteilt, Schritt 3)";
  die §Fehlerpräzedenz nennt „Vakuität (Schritt 5)" korrekt. Genau diese
  Referenz ist der Anker, an dem R1 die Namensraum-Vorbedingung festgemacht hat —
  sie führt einen Nachleser künftig auf den falschen Schritt.
- **verifizierbar:** nein maschinell (kein Gate prüft Schritt-Querverweise);
  durch Lesen von `spec/spezifikation.md:445` gegen `:450` belegt.

### R3-F-3 — Der R2-F-2-Fix (Trennung `readCrossFile` / `spanCrossSource`) ist von keinem Test gepinnt

- **kategorie:** LOW
- **quelle:** Maintainability; [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
  §Fehlerpräzedenz („Quellen lesen … → Abschnitts-Spannung … Jede Stufe läuft
  über beide Sichten, bevor die nächste beginnt")
- **pfad:** `internal/hexagon/core/app/trace_cross_test.go:222-242`
  (`TestCrossConsistencyFailClosed`)
- **befund:** Die Tabelle führt „fehlende Rückwärts-Datei" und „Abschnittsname
  ohne Heading-Treffer" als **getrennte** Fälle; keiner kombiniert einen
  Vorwärts-Sektionsfehler mit einer fehlenden Rückwärts-Datei. Genau diese
  Kombination ist aber der einzige beobachtbare Unterschied zwischen der neuen
  Zwei-Stufen-Form und der alten, je Sicht bündelnden `loadCrossSource`. Führt
  man die beiden Aufrufe wieder zusammen, bleibt die Suite grün — der Fix ist
  regressionsfrei rückbaubar.
- **verifizierbar:** ja — `readCrossFile`/`spanCrossSource` wieder zu einem
  Aufruf je Sicht bündeln und `make test` laufen lassen: grün. Gegen
  `d-check:latest` ist das Verhalten dagegen belegt (fehlende
  `backward.file` schlägt jetzt korrekt vor dem `forward`-Sektionsfehler zu).

### R3-F-4 — `backwardIDColumn`: der `count > 1`-Zweig ist ungepinnt

- **kategorie:** LOW
- **quelle:** [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
  §Fehlerpräzedenz („`artifact-id-column`, wenn ≠ `first`, je genau einmal");
  [`DC-FA-XREF-001`](../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
  (Negative-Kriterium „ID-Header nicht genau einmal")
- **pfad:** `internal/hexagon/core/app/trace_cross.go:212-214`
  (`if count != 1`); `internal/hexagon/core/app/trace_cross_test.go:380-403`
- **befund:** `TestCrossConsistencyBackwardIDHeaderFehlt` deckt `count == 0`
  („0-mal"); kein Test führt eine Rück-Tabelle mit **zweimal** dem
  konfigurierten ID-Header. `count != 1` zu `count == 0` mutiert lässt die Suite
  grün, obwohl dann ein doppelter `Kennung`-Header still an sein **erstes**
  Vorkommen bindet — eine geratene Spalte, die falsche Artefakt-IDs und damit
  falsche Befunde erzeugen kann. Das Vorwärts-Pendant (`bindCrossColumns`,
  `counts[n] > 1`) ist über `TestCrossConsistencyDuplicateHeader` gepinnt; die
  Rück-Seite nicht.
- **verifizierbar:** ja — `count != 1` → `count == 0` mutieren, `make test`: grün.

## Negativbefunde (geprüft, ohne Befund)

### Vertrags-Konsistenz (Frage 2)

- **Lastenheft ↔ Spezifikation ↔ ADR-0038 — Vakuitäts-Definition:** ohne Befund,
  wortgleich in allen drei Strata. Lastenheft: „beide Sichten kantenleer — oder,
  unter `mode: superset`, die Rück-Sicht kantenleer"; Spec Schritt 5: „(a) `F`
  **und** `B` leer … (b) `B` leer **und** `mode: superset`"; ADR-0038
  Entscheidung 8 identisch. Die Abgrenzung „einseitig leere Sicht = Ergebnis" ist
  in allen dreien deckungsgleich formuliert. Kein Widerspruch (der einzige
  Vertrags-Riss ist die Stufen-**Reihenfolge**, R3-F-1).
- **Fehlerpräzedenz ↔ Code-Reihenfolge:** ohne Befund — deckungsgleich.
  `crossConsistency:77-112` fährt `readCrossFile` ×2 (`:77-84`) →
  `spanCrossSource` ×2 (`:85-92`) → `bindForwardTables`/`bindBackwardTables`
  (`:93-100`) → `forwardEdges`/`backwardEdges` (`:101-108`) → `crossVacuity`
  (`:109`) → `diffViews` (`:112`). Das ist exakt „Quellen lesen →
  Abschnitts-Spannung → Header-Bindung → Range-Expansion → Vakuität → Diff", und
  jede Stufe läuft über beide Sichten, bevor die nächste beginnt. Der neue
  Präzedenz-Satz („Jede Stufe läuft **über beide Sichten**…") beschreibt, was der
  Code tut.
- **Schritt-Nummerierung (Frage 2, Teil 3):** bis auf R3-F-2 stimmig. Die
  Einfügung von Schritt 5 und das Nachziehen von 6/7 ist korrekt; die
  Binnen-Referenzen in Schritt 5 („Schritt 3", „Schritt 6") und in der
  §Fehlerpräzedenz („Schritt 5") sind richtig; die Einleitung („Schritt 2 von
  §`DC-FA-CLI-009.a`") verweist auf eine andere Sektion und ist unberührt. Einzig
  Schritt 3 blieb stehen (R3-F-2).
- **[ADR-0038](../plan/adr/0038-trace-cross-consistency.md)-Immutabilität
  ([`AGENTS.md`](../../AGENTS.md) §3.5):** ohne Befund — Status ist weiterhin
  `Proposed`, damit ist die inhaltliche Ergänzung zulässig; die Änderung ist
  additiv (Entscheidung 8, ein Fitness-Funktions-Punkt, ein Geschichts-Eintrag),
  bestehende Entscheidungen 1–7 sind unangetastet. `make adr-check
  RANGE=d11398f..HEAD`: 0 Befunde.
- **[`AGENTS.md`](../../AGENTS.md) §3.6 (Gate-Lockerung nur per ADR):** ohne
  Befund — die Lockerung gegenüber `d11398f` (einseitig leere Sicht nicht mehr
  Exit 2) ist genau der Weg, den §3.6 verlangt: doc-first über Lastenheft-CR
  0.44.1 + ADR-0038 Entscheidung 8, **vor** dem Code-Commit (`1665714` vor
  `9027514`). Vorbildlich sequenziert.

### MR-006 / Referenz-Richtung (Frage 3)

- **Lastenheft-Historie ohne ADR-IDs:** ohne Befund — konsistent mit dem Bestand
  **und** mit
  [`MR-006`](../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs).
  Die 0.44.1-Zeile schreibt „Begründung + Abgrenzung … in der begleitenden ADR"
  — dieselbe ID-freie Formel wie 0.44.0 („in der Spezifikation"), 0.42.0/0.41.0
  („in begleitender ADR") und 0.34.0. Keine Abweichung.
- **Spezifikations-Historie mit `[ADR-0038](…)`-Link:** ohne Befund, aber als
  Bestands-Beobachtung notiert. `spec/spezifikation.md` ist Rang 2 und damit
  Spec-Stratum; die Historie-Zeile verlinkt abwärts auf die ADR. Das ist
  **gate-legal** — die maschinelle MR-006-Kodierung in
  [`.d-check.yml`](../../.d-check.yml) trägt
  `matrix.exclude-sections: [Historie, "7. Historie", Geschichte]`, die
  Historie ist also bewusst aus der `{from: spec-straten, to: adr, allow: false}`-
  Regel herausgenommen — und **Bestand**: die Zeilen 2026-07-14 (ADR-0037) und
  2026-07-16 (ADR-0038) tun dasselbe. Die neue Zeile setzt das Muster fort, führt
  es nicht ein. (Dass MR-006s Prosa diesen Historie-Carveout nicht nennt, ist eine
  vorbestehende Doku-/Mechanik-Lücke außerhalb dieses Slice.)
- **Token-Referenz-Richtung (`DC-FA-MTX-003`):** ohne Befund — der
  `slice-071`-Token in ADR-0038 steht ausschließlich in `## Geschichte`, die
  `matrix.exclude-sections` ohnehin ausnimmt; Entscheidung 8 und die
  Fitness-Funktion nennen keinen Slice-Token. Die Referenzen der
  ADR-Geschichte auf Lastenheft/Spezifikation sind aufwärts (adr →
  spec-straten, keine Verbotsregel). `make doc-check`: 0 Befunde über 222
  Dateien.
- **Marker-Ehrlichkeit (Reviewer-Skill MEDIUM-Anker):** ohne Befund — der Diff
  setzt keinen `<!-- d-check:status-provenance -->`-Marker; es gibt keine als
  Provenance getarnte Entscheidungsgrundlage.

### Code gegen den neuen Vertrag (Frage 1)

- **`crossVacuity` (a) ↔ Lastenheft/Spec:** ohne Befund — `len(fwd) == 0 &&
  len(bwd) == 0` (`:290`) ist die wörtliche Umsetzung von „`F` **und** `B` leer".
  Gegen `d-check:latest` verifiziert: `design-pattern: 'GG-ARCH-COMP-[A-Z]+'`
  (Tippfehler-Klasse) ⇒ Exit 2 mit „beide Sichten ergaben 0 Kanten".
- **`crossVacuity` (b) ↔ Lastenheft/Spec:** ohne Befund — `len(bwd) == 0 && mode
  == model.TraceCrossModeSuperset` (`:294`) ist wörtlich „`B` leer **und**
  `mode: superset`". Verifiziert: `F ≠ ∅`, Rück-Erste-Spalte ohne ID, `superset`
  ⇒ Exit 2 mit „Rück-Sicht ergab 0 Kanten und mode: superset gatet allein B \ F".
  Die Begründung im Meldungstext benennt die Ursache korrekt.
- **Boundary „einseitig leere Vorwärts-Sicht ist ein Ergebnis" (R2-F-1
  geschlossen):** ohne Befund — verifiziert. Bootstrap-Repo (Vorwärts-Zelle
  „alle Scheduler-Komponenten (siehe Architektur)", zwei gepflegte Rück-Kanten)
  liefert jetzt **zwei laute `B\F`-Befunde mit `Datei:Zeile` und Exit 1** statt
  Exit 2 — genau das Lastenheft-Boundary-Kriterium. Die Config-Fehldiagnose ist
  weg. Der von R2 als vertragswidrig nachgewiesene Pro-Sicht-Guard existiert nicht
  mehr.
- **`superset`/`F = ∅`/`B ≠ ∅`:** ohne Befund — kein Vakuum ((a) scheitert an
  `bwd`, (b) ebenso), der Diff meldet `|B|` laut. Vertragstreu: „Nicht vakuum ist
  eine einseitig leere Sicht mit nicht-leerer Gegenseite."
- **`equal`/`B = ∅`/`F ≠ ∅`:** ohne Befund — kein Vakuum, `F\B` meldet laut.
  Durch die Gegenprobe in
  `TestCrossConsistencyVakuumRueckSichtLeerUnterSuperset:334-341` gepinnt.
- **„Zu streng" (Frage 6, erste Hälfte):** **nein** — `crossVacuity` ist an
  keinem Rand strenger als der Vertrag. Beide Zweige sind wörtliche Umsetzungen
  der Konjunktionen (a)/(b); die prä-Ausschluss-Auswertung macht den Guard
  strikt **schwächer** als die post-Ausschluss-Lesart (prä-Vakuum ⊆
  post-Vakuum), nie stärker. Die einzige Abweichung liegt damit auf der
  Lasch-Seite (R3-F-1).
- **`len(out) == 0`-Guard bleibt tragend:** ohne Befund — dass eine Quelle
  **gar keine** Tabelle bindet, bleibt in der Bind-Stufe Exit 2 (`:158-161`,
  `:198-201`) und wird nicht mit dem Datenzustand „gebunden, aber kantenleer"
  verwechselt. Genau diese Trennung macht den Bootstrap-Fall vom Config-Fehler
  unterscheidbar; ohne sie wäre eine völlig unbindbare Vorwärts-Quelle als
  Bootstrap fehlgedeutet worden.

### Regressionen in `9027514` (Frage 4)

- **Determinismus ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)):**
  ohne Befund — `crossVacuity` liest nur zwei `len()` und `mode`, hat keine
  Map-Iteration, und beide Meldungstexte sind konstant. Extraktions- und
  Sortier-Determinismus aus R1/R2 unberührt.
- **Tote Pfade:** ohne Befund — `crossNullGuard` und `loadCrossSource` sind
  restlos entfernt (`grep` über `internal/`: keine Treffer); `readCrossFile` ist
  in der neuen Rolle reaktiviert, nicht dupliziert; `crossSource.file` wird in
  beiden Bind-Meldungen weiterhin genutzt. Kein verwaister Parameter.
- **Byte-Identität ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)):**
  ohne Befund — `9027514` fasst weder `diffViews`, `report.go`, `cli.go`,
  `model/config.go` noch `configyaml.go` an; `crossVacuity` läuft nur bei
  `cc != nil`. Die R1-Negativbefunde gelten unverändert;
  `TestCLI071_Cross_DefaultByteIdentisch` bleibt grün.
- **Erfundener Vertrag:** ohne Befund — `crossVacuity` fügt keine Prüfung
  hinzu, die Lastenheft/Spec nicht wörtlich verlangen. Der Vertrag ist diesmal
  *vor* dem Code entstanden, und der Code unterschreitet ihn eher (R3-F-1), als
  ihn zu überschreiten.
- **Hexagon-Import-Regeln / [`AGENTS.md`](../../AGENTS.md) §3.2:** ohne Befund —
  keine neuen Importe, keine `//nolint`-Direktive, keine
  `.golangci.yml`-Ausnahme. `make arch-check`: 0 Befunde.

### Mutations-Härte (Frage 5)

- **Die drei benannten Zweige sind einzeln gepinnt — Behauptung bestätigt.**
  Nachvollzogen: (i) Vakuum (a) entfernt ⇒ nur
  `TestCrossConsistencyVakuumBeideSichtenLeer` kippt (dessen `crossCfg()` fährt
  `equal`, also greift (b) nicht); (ii) Vakuum (b) entfernt ⇒ nur
  `TestCrossConsistencyVakuumRueckSichtLeerUnterSuperset` kippt (dort ist `F`
  nicht leer, (a) greift also nicht — die Isolation ist bewusst konstruiert und
  im Kommentar benannt); (iii) Relevanzbindung auf `ArtifactIDColumn`
  zurückgedreht ⇒ nur `TestCrossConsistencyBackwardIDHeaderFehlt` kippt.
- **`&&`-Konjunktionen gepinnt:** `TestCrossConsistencyLeereVorwaertsSichtIstKeinVakuum`
  kippt, sobald (a) zu `||` mutiert; die `equal`-Gegenprobe in
  `…VakuumRueckSichtLeerUnterSuperset:334-341` kippt, sobald `mode ==
  superset` aus (b) fällt. Beide Konjunktionen sind damit belegt — die
  R2-F-3-Lücke ist geschlossen.
- **R2-F-4 geschlossen:** `TestCrossConsistencyBackwardIDHeaderFehlt` führt jetzt
  **zwei** `Bezug`-Tabellen mit heterogenen ID-Headern und `mode: superset`. Unter
  dem Alt-Code bindet Tabelle 1 (`B` bleibt nicht-leer, kein Vakuitäts-Guard
  greift), Tabelle 2 verschwindet still, `B\F` ist leer ⇒ Exit 0. Der Test trennt
  alt und neu jetzt am **Verhalten**, nicht mehr nur am Meldungstext; die
  zusätzliche `"0-mal"`-Assertion hält die Diagnose fest. Die R2-F-4-Beobachtung
  ist vollständig adressiert.
- **Verbleibende ungepinnte Zweige:** R3-F-3 (Präzedenz-Split) und R3-F-4
  (`count > 1`) — siehe Findings. Darüber hinaus keine gefunden: der
  `first`-Sentinel-Zweig (`:208-210`) ist über den Default aller übrigen
  Cross-Tests gepinnt, `checkCrossSections` über
  `TestCrossConsistencyFailClosed:239-242`, `crossBadRow` über den
  `DC-FA-REQ-001`-Pfad.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 1 | R3-F-1 |
| LOW | 3 | R3-F-2, R3-F-3, R3-F-4 |
| INFO | 0 | — |

## Verdikt

**REJECT** — erneut knapp, und diesmal fast ausschließlich am **Text**, nicht am
Code.

Zuerst das Positive, weil es substanziell ist: der Weg über `1665714` vor
`9027514` ist genau das, was
[`AGENTS.md`](../../AGENTS.md) §3.6 und der Doc-first-Prozess dieses Repos
verlangen — die Lockerung („einseitig leere Sicht ist kein Fehler mehr") ist als
Lastenheft-CR und ADR-Entscheidung begründet, **bevor** der Code ihr folgt, nicht
als PR-Kommentar. Der Vertragstext ist in Lastenheft, Spezifikation und
ADR-0038 wortgleich, die Abgrenzung zu
[ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md) ist sauber
gezogen („dort die explizite Quelle, hier der Vergleich"). Der Code deckt beide
Vakuitäts-Zweige wörtlich, die Fehlerpräzedenz stimmt jetzt Stufe für Stufe mit
der Ausführungsreihenfolge überein, R2-F-1 bis R2-F-4 sind **alle vier**
geschlossen und einzeln gegen das Image verifiziert — insbesondere meldet der
Bootstrap-Zustand jetzt seine Rück-Kanten laut (Exit 1) statt ihn mit einer
Config-Fehldiagnose abzuwürgen. Und die Mutations-Härte-Behauptung hält der
Nachprüfung stand.

Blockierend bleibt **R3-F-1**, und es ist kein Nachtreten: die Vakuitäts-Stufe
ist als Schritt 5 **hinter** den Ausschluss (Schritt 4) einsortiert, während die
§Fehlerpräzedenz sie ohne Ausschluss-Stufe direkt hinter die Range-Expansion
setzt. Zwei Stellen desselben Abschnitts ordnen dieselbe fail-closed-Stufe
verschieden ein — in einem Repo, in dem das Dokument führt, ist ein Abschnitt,
der sich selbst widerspricht, kein Detail. Meine Lesart der *Absicht*: gemeint
ist prä-Ausschluss (so argumentiert
[ADR-0038](../plan/adr/0038-trace-cross-consistency.md) Entscheidung 8
durchgehend über Muster, die „am Inhalt vorbeigreifen", und so ordnet es die
Fehlerpräzedenz) — dann ist der **Code richtig** und die Stufe steht an der
falschen Nummer. Unter der wörtlichen Schrittfolge dagegen ist der Code zu lasch
und ein alles tilgendes `exclude-req` meldet einen echten Drift als
`0 Differenz(en).` — reproduziert. Welche Lesart gilt, entscheidet der Vertrag,
nicht das Review; nur widersprechen darf er sich nicht.

Zur ausdrücklichen Frage 6: `crossVacuity` ist an **keinem** Rand strenger als der
Vertrag (beide Zweige sind wörtliche Umsetzungen der Konjunktionen, und
prä-Ausschluss ist strukturell die schwächere Prüfung). Die einzige Abweichung
liegt auf der Lasch-Seite und ist R3-F-1.

R3-F-2 (stale „Schritt 5") gehört sinnvollerweise in denselben Handgriff wie
R3-F-1 — beide sitzen in der Schritt-Nummerierung. R3-F-3/F-4 sind
Test-Restlücken ohne Blockade-Anspruch; R3-F-3 ist die pikantere, weil sie genau
den R2-Fix ungepinnt lässt, den sie absichern soll.
