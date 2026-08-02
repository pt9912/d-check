# Review — slice-071 Implementierung R4 (Vertrags- + Fix-Review zu R3)

**Datum:** 2026-07-17
**Review-Art:** unabhängiger Vertrags- und Fix-Review (kontext-getrennt;
adversarial, nicht bestätigend)
**Gegenstand:**
[`slice-071`](../plan/planning/done/slice-071-trace-cross-consistency-gate.md) —
`815a2a7` (Auflösung des R3-Vertragswiderspruchs zugunsten **post**-Ausschluss:
Spec Schritt 5 + Fehlerpräzedenz + Code + Tests + Diagnose);
Slice gesamt `6c4ccf5..HEAD`
**Vor-Reviews:** [R1](2026-07-17-slice-071-implementation-r1.md) (REJECT),
[R2](2026-07-17-slice-071-implementation-r2.md) (REJECT),
[R3](2026-07-17-slice-071-implementation-r3.md) (REJECT)
**Reviewer:** Claude (kontext-unabhängiger Lauf)
**Skill:** `.harness/skills/reviewer.md` v1.2.0
**Modell:** Opus 4.8 (1M context)

## Eingangs-Kontext

- [`DC-FA-XREF-001`](../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
  (Lastenheft 0.44.1, **unverändert** durch `815a2a7`),
  [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
  Schritte 1–7 + §Fehlerpräzedenz (geändert),
  [ADR-0038](../plan/adr/0038-trace-cross-consistency.md) Entscheidungen 5 + 8
  (**unverändert** durch `815a2a7`)
- [`AGENTS.md` §5 (Anforderungs-Anlege-Prozess)](../../AGENTS.md#5-dokumentations-regeln),
  [`MR-001`](../../harness/conventions.md#mr-001--source-precedence-mit-eigener-spezifikations-schicht),
  [`AGENTS.md`](../../AGENTS.md) §3.4/§3.5/§3.6
- Prüfgegenstand: `internal/hexagon/core/app/trace_cross.go`,
  `internal/hexagon/core/app/trace_cross_test.go`, `spec/spezifikation.md`

**Ausgeführte Sensoren:** `make lint` (0 issues), `make test` (grün),
`make arch-check` (0 Befunde), `make doc-check` (223 Dateien, 0 Befunde),
`make adr-check RANGE=9027514..HEAD` (0 Befunde), `gofmt -l` (kein neuer
Treffer), `make build` + vier Repro-Läufe gegen `d-check:latest`
(`--network none`, `:ro`).

## Vorbemerkung: die post-Ausschluss-Auflösung ist richtig — meine R3-Lesart war es nicht

R3 hat den Widerspruch korrekt gefunden, aber die Auflösung falsch geneigt. Ich
korrigiere das ausdrücklich, weil ein Review, das seine eigene Empfehlung nicht
revidiert, wertlos ist:

- Mein R3-Argument stützte sich darauf, dass die §Fehlerpräzedenz **keine**
  Ausschluss-Stufe listete, und las diese Stille als Absichtserklärung. Das war
  ein Fehlschluss: die Stille war eine **Auslassung** (genau die, die `815a2a7`
  jetzt nachträgt), kein Signal.
- Ich habe [ADR-0038](../plan/adr/0038-trace-cross-consistency.md) Entscheidung 8
  einseitig gelesen. Ihr **illustrativer** Satz („ein `design-pattern`, das …
  vorbeigreift") ist nicht ihr **normativer**: „hier ist der Bezugspunkt nicht
  eine Sicht, sondern der **Vergleich** — geguardet wird, **was
  konstruktionsbedingt nie einen Befund liefern kann**." Ein Ventil, das alle
  Anforderungen verschluckt, ist exakt das. Entscheidung 8 trug die
  post-Ausschluss-Lesart bereits; ich habe sie überlesen.
- Die Begründung „maßgeblich ist, was am Ende tatsächlich verglichen wird" ist
  zudem die einzige, die den Guard über seine **Wirkung** statt über eine
  **Ursachenliste** definiert — und damit die einzige, die bei der nächsten,
  noch unbekannten Ursache nicht wieder reißt.

Post-Ausschluss ist die schärfere **und** die prinzipiellere Lesart. Der Code
folgt ihr korrekt (verifiziert). Die verbleibenden Findings betreffen
ausschließlich, **wo** diese Entscheidung dokumentiert ist — und eine
Diagnose-Regression.

## Findings

### R4-F-1 — Der Ventil-Trigger fehlt in den Lastenheft-Akzeptanzkriterien: der Bump auf 0.44.2 ist nötig

- **kategorie:** MEDIUM
- **quelle:** [`DC-FA-XREF-001`](../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
  (§Vakuum-Absatz, §Ableitungssprung-Absatz, AK „Negative (Vakuum,
  fail-closed)");
  [`AGENTS.md` §5 (Anforderungs-Anlege-Prozess)](../../AGENTS.md#5-dokumentations-regeln)
  („Drei Akzeptanzkriterien … **Versions-Bump + Historie-Zeile**" als
  Pflicht-Bausteine bei **geänderten** `DC-*`-Anforderungen);
  [`MR-001`](../../harness/conventions.md#mr-001--source-precedence-mit-eigener-spezifikations-schicht)
- **pfad:** `spec/lastenheft.md:739-749` (Vakuum-Absatz), `:733-737`
  (Ableitungssprung-Absatz), `:778-782` (AK „Negative (Vakuum)")
- **befund:** `815a2a7` erzeugt einen **neuen, ursächlich eigenständigen
  Exit-2-Auslöser** — eine fehlerfreie Config mit korrekten Mustern, aber
  übergriffigem `exclude-req`. Kein Lastenheft-Kriterium beschreibt ihn: die
  AK „Negative (Vakuum, fail-closed)" ist in **beiden** Zweigen ursachen-gefasst
  („Given ein `design-pattern`, das … vorbeigreift … — oder eine kantenleere
  Rück-Sicht unter `mode: superset`"); der Ventil-Fall ist keiner von beiden.
  Der Vakuum-Absatz bindet den Zustand an dieselbe Ursache („**greifen die
  Muster am Inhalt vorbei**, fällt aus beiden Sichten keine Kante"), und der
  Ableitungssprung-Absatz beschreibt `exclude-req` als Ventil, ohne dass es je
  Exit 2 auslösen könnte. Ein Konsument, der nur Rang 1 liest, kann das
  Verhalten nicht herleiten; eine aus den AK abgeleitete Abnahme prüft es nie.
  Die Argumentation „Rang 1 definiert den Zustand, nicht den Messpunkt" trägt
  für die **Stufenfolge** — aber die Stufenfolge erzeugt hier einen Trigger, und
  Trigger sind Rang-1-Material: 0.44.1 selbst wurde genau dafür gebumpt und gab
  der Vakuum-Klasse **beide** Kriterien (Negative + Boundary). Der zweite
  Trigger derselben Klasse bleibt ohne.
- **verifizierbar:** ja — gegen `d-check:latest` reproduziert: Repo mit
  Trigger-088-Drift, fehlerfreien Mustern und `exclude-req: '.'` ⇒ Exit 2. Kein
  AK in
  [`DC-FA-XREF-001`](../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
  deckt diesen Given-Zustand ab.

### R4-F-2 — ADR-0038 trägt die post-Ausschluss-Entscheidung nicht; ihre Begründung lebt in Rang 2

- **kategorie:** MEDIUM
- **quelle:** [`AGENTS.md`](../../AGENTS.md) §3.4 („Die sprachkonkrete
  Übersetzung … und **die Begründungen leben in den ADRs**"), §3.5 (Accepted ⇒
  immutable); [ADR-0038](../plan/adr/0038-trace-cross-consistency.md)
  Entscheidungen 5 + 8
- **pfad:** `docs/plan/adr/0038-trace-cross-consistency.md:80-100`
  (Entscheidung 8, durch `815a2a7` unberührt), `:64-66` (Entscheidung 5);
  `spec/spezifikation.md:450-454` (Begründung in Rang 2)
- **befund:** Die Entscheidung „Vakuität wird **nach** dem Ausschluss gemessen"
  war strittig — R3 argumentierte für prä, der Auftraggeber entschied für post.
  Genau dafür existieren ADRs. Aufgezeichnet ist sie jedoch nur in der
  Spezifikation („denn maßgeblich ist, was am Ende tatsächlich verglichen wird:
  ein `exclude-req`, das alle Anforderungen verschluckt …") und in der
  Commit-Message; ADR-0038 ist unverändert. Damit steht eine **Begründung** in
  Rang 2, wo §3.4 sie in der ADR verortet. Zwei Folgen: (1) Entscheidung 8
  beschreibt das implementierte Verhalten nur noch teilweise — sie nennt allein
  die Muster-Ursache, während der Guard zwei ursächlich verschiedene Auslöser
  hat; (2) Entscheidung 5 („`exclude-req` … **kein** gelöstes Problem") ist
  überholt, seit der Totalfall des Ventils geguardet ist (die Teil-Drift bleibt
  es). Das Fenster schließt sich: ADR-0038 steht auf `Proposed` und wird bei der
  Closure auf `Accepted` gehen — danach ist sie nach §3.5 immutabel und eine
  Korrektur verlangt eine **neue ADR mit `Supersedes`**.
- **verifizierbar:** ja — `git show 815a2a7 -- docs/plan/adr/` ist leer; die
  §Geschichte von ADR-0038 endet mit dem 1665714-Eintrag, obwohl die dort
  festgehaltene Entscheidung 8 seither in ihrem Messpunkt geschärft wurde.

### R4-F-3 — Die Vakuitäts-Meldungen sind grammatisch gebrochen (Regression aus `815a2a7`)

- **kategorie:** LOW
- **quelle:** [`DC-FA-XREF-001`](../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
  (nutzersichtbare Diagnose); Maintainability
- **pfad:** `internal/hexagon/core/app/trace_cross.go:308-310`, `:313-315`
  (Hint-Argumente), `:321-327` (`vacuityHint`)
- **befund:** `vacuityHint` setzt `"prüfe, ob "` vor einen Hint, der in
  **Hauptsatz**-Stellung formuliert ist (`"… .design-pattern trifft den
  Artefakt-Namensraum beider Sichten"`). Im `ob`-Nebensatz muss das Verb ans
  Ende; heraus kommt „prüfe, ob …design-pattern **trifft den Artefakt-Namensraum
  beider Sichten**". Beide Zweige sind betroffen (auch „prüfe, ob
  …edge-column/req-pattern **treffen die Rück-Kanten**"). Es ist eine
  **Regression**: `9027514` formulierte „prüfe, ob %s.design-pattern den
  Artefakt-Namensraum beider Sichten **trifft**" (korrekt), die Extraktion in die
  Hint-Variable hat die Wortstellung zerlegt. Der angehängte Ventil-Teil („und ob
  `exclude-req` nicht alle Anforderungen verschluckt") ist korrekt — der Bruch
  betrifft nur den Muster-Hint. Die Tests prüfen Teilstrings, die die
  Wortstellung nicht berühren.
- **verifizierbar:** ja — gegen `d-check:latest`:
  `d-check: error: trace.cross-consistency: beide Sichten ergaben 0 Kanten — der
  Abgleich verglich nichts; prüfe, ob trace.cross-consistency.forward.design-pattern
  trifft den Artefakt-Namensraum beider Sichten`.

### R4-F-4 — Die `patternHint`-Argumente beider `crossVacuity`-Zweige sind ungepinnt (Antwort auf Frage 4)

- **kategorie:** LOW
- **quelle:** Maintainability; Reviewer-Skill (Mutations-Härte)
- **pfad:** `internal/hexagon/core/app/trace_cross.go:308-315`;
  `internal/hexagon/core/app/trace_cross_test.go:307-321` (Vakuum a), `:334-341`
  (Vakuum b), `:432-460` (Ventil)
- **befund:** `vacuityHint`s **Verzweigung** (Ventil gesetzt / nicht gesetzt) ist
  in beide Richtungen gepinnt — das ist erledigt. Ungepinnt ist der **Hint-Inhalt
  je Aufrufstelle**: alle drei Tests assertieren auf Präfixe aus dem
  `fmt`-Rahmen („beide Sichten ergaben 0 Kanten", „Rück-Sicht ergab 0 Kanten")
  bzw. auf `"exclude-req"` — keiner auf `"design-pattern"` oder
  `"edge-column"`. Vertauscht man die beiden `patternHint`-Argumente an den
  Aufrufstellen (a) und (b), bleibt die Suite grün, obwohl dann das
  Namensraum-Vakuum in die Rück-Kanten-Config und das superset-Vakuum in das
  `design-pattern` schickt — genau die Fehldiagnose-Klasse, die dieser Commit
  beseitigen wollte. Zusätzlich unbelegt: die Kombination (b) + gesetztes Ventil
  (superset-Vakuum mit `exclude-req`) wird von keinem Test durchlaufen.
- **verifizierbar:** ja — die beiden `vacuityHint`-Hint-Argumente in `:309` und
  `:314` tauschen, `make test`: grün.

## Negativbefunde (geprüft, ohne Befund)

### Widerspruchsfreiheit post-Ausschluss (Frage 1)

- **Spezifikation ↔ Code:** ohne Befund, deckungsgleich. Schritt 5 („**nach**
  dem Ausschluss (Schritt 4) gemessen") ↔ `crossConsistency:109-113`
  (`excludeReqs` ×2 → `crossVacuity` → `diffViews`). Die §Fehlerpräzedenz
  („… → Range-Expansion → **Ausschluss (Schritt 4)** → Vakuität (Schritt 5) →
  Diff (Schritt 6)") ist jetzt Stufe für Stufe die Ausführungsreihenfolge. Der
  in R3 gemeldete Widerspruch zwischen Schrittfolge und Präzedenz ist restlos
  aufgelöst — beide sagen dasselbe.
- **R3-F-2 geschlossen:** Schritt 3 verweist jetzt korrekt auf „Mengen-Diff
  (Schritt 6)" und ergänzt sinnvoll „die Vakuitäts-Prüfung (Schritt 5) fängt
  genau diesen Fall". Ich habe alle Schritt-Querverweise des Abschnitts erneut
  durchgezählt (`Schritt 2 von §DC-FA-CLI-009.a`, `Schritt 3`, `Schritt 4`,
  `Schritt 5`, `Schritt 6`): alle stimmig, keine stale Referenz mehr.
- **ADR-0038 ↔ Spezifikation — kein Widerspruch:** ohne Befund. Entscheidung 8
  legt den Messpunkt nicht fest und wird von der Präzisierung nicht verletzt;
  ihr normativer Kern („geguardet wird, was konstruktionsbedingt nie einen
  Befund liefern kann") **stützt** post-Ausschluss. Die Lücke ist eine
  Auslassung, kein Konflikt (R4-F-2).
- **Lastenheft ↔ Spezifikation — kein Widerspruch:** ohne Befund. Die
  Zustandsdefinition („Vakuum heißt: **beide** Sichten kantenleer — oder, unter
  `mode: superset`, die **Rück**-Sicht kantenleer") ist messpunkt-frei und wird
  von Schritt 5 präzisiert, nicht bestritten. Insoweit trägt die Begründung des
  Auftraggebers. Was sie **nicht** trägt, ist die AK-Seite (R4-F-1) — die
  Unterscheidung ist: die *Definition* deckt den Ventil-Fall, die
  *Akzeptanzkriterien* decken ihn nicht, und Letztere sind die Abnahme-Fläche.
- **[`AGENTS.md`](../../AGENTS.md) §3.6:** ohne Befund — `815a2a7` **verschärft**
  (neuer Exit-2-Auslöser), lockert nichts; §3.6 ist nicht berührt. Die
  Verschärfung ist zudem doc-first begründet (Spec vor Code im selben Commit,
  Reihenfolge im Body nachvollziehbar).
- **[`AGENTS.md`](../../AGENTS.md) §3.5 / `make adr-check`:** ohne Befund —
  ADR-0038 ist unangetastet, Status `Proposed`. `make adr-check
  RANGE=9027514..HEAD`: 0 Befunde.

### Legitime Fälle unter post-Ausschluss (Frage 2)

- **Kein legitimer Fall geht verloren** — geprüft, ohne Befund. Der tragende
  Grund: `excludeReqs` wendet **dasselbe** Prädikat auf **beide** Sichten an
  (`:279-289`), also bedeutet post-Ausschluss-Leere „jede Anforderung, die dieser
  Abgleich überhaupt vergleichen könnte, ist ausgeschlossen" — ein konfigurierter
  Check, der beweisbar nie einen Befund liefern kann. Das **ist** die
  Vakuum-Definition. Und weil der ganze Block **opt-in** ist, lautet der
  korrekte Ausdruck von „ich will hier nicht prüfen" *Block weglassen*, nicht
  *Block konfigurieren und alles ausschließen*. Drei Kandidaten durchgespielt:
  (i) beide Sichten enthalten ausschließlich Mittelschicht-IDs, die zu Recht alle
  ausgeschlossen werden ⇒ der Check ist inert, rot ist richtig; (ii)
  Migrations-Zwischenstand ⇒ Block temporär weglassen, nicht leerventilieren;
  (iii) Repo schrumpft, bis nur noch ausgeschlossene Familien übrig sind ⇒ das
  Gate vergleicht nichts mehr, rot ist richtig. Keiner ist ein Fehlalarm.
- **Keine asymmetrische Verzerrung:** ohne Befund — da das Prädikat auf beide
  Sichten wirkt, kann `r ∈ F'` nie bedeuten, dass die zugehörige Rück-Kante
  ausgeschlossen wurde; ein Ventil kann keine Phantom-`F\B`-Befunde erzeugen.

### Regressionen in `815a2a7` (Frage 3)

- **`exclude-req`-Semantik sonst unverändert (Ableitungssprung-Fall):** ohne
  Befund. `TestCrossConsistencyExcludeReq` (unverändert, grün) belegt: `F =
  {ARCH-006:{CORE}}`, `B = {ARCH-006:{CORE}, SPEC-042:{MID}}`, `exclude-req:
  ^GG-SPEC-` ⇒ `F' = B' = {ARCH-006:{CORE}}`, kein Vakuum, 0 Differenzen. Die
  Verschiebung aus `diffViews` in eine eigene Stufe ist für diesen Fall
  verhaltensneutral — das Prädikat und sein Wirkort (Schlüsselmengen beider
  Sichten) sind identisch. Zusätzlich durch die Gegenprobe in
  `…VakuumDurchUebergriffigesVentil:453-460` belegt (enges Ventil ⇒ Drift steht,
  2 Befunde).
- **`diffViews` ohne Ausschluss korrekt:** ohne Befund. Die Signatur nimmt nur
  noch `mode`; die Vereinigung über `keys(F) ∪ keys(B)`, die Richtungs-Zweige
  und die Sortierung nach `(R, Artefakt, Richtung)` sind unverändert. Da die
  Stufe vorher lief, ist `keys` bereits bereinigt — die entfernte
  `ExcludeReq.MatchString(r)`-Zeile wäre nur noch tote Wiederholung gewesen.
- **Determinismus ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)):**
  ohne Befund. `excludeReqs` iteriert zwar eine Map, baut daraus aber nur eine
  Map (kein ordnungsabhängiges Ergebnis); `crossVacuity` liest zwei `len()` und
  `mode`; `vacuityHint` verzweigt über einen nil-Check. Keine neue
  Nichtdeterminismus-Quelle. Sortier-/Fundstellen-Determinismus aus R1–R3
  unberührt.
- **Aliasing:** ohne Befund — `excludeReqs` gibt bei `exclude == nil` die
  Original-Map zurück und teilt sonst die inneren `map[string]crossEdge`-Werte;
  alle Konsumenten stromabwärts lesen nur. Kein Failure-Szenario.
- **Byte-Identität ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)):**
  ohne Befund — `815a2a7` fasst `report.go`, `cli.go`, `model/config.go`,
  `configyaml.go` nicht an; `crossConsistency` läuft nur bei `cc != nil`.
  `TestCLI071_Cross_DefaultByteIdentisch` grün.
- **Tote Pfade / erfundener Vertrag:** ohne Befund — `crossVacuity` nimmt jetzt
  `cc` statt `mode` (beide Felder genutzt); kein verwaister Parameter, kein Rest
  der alten Guards (`grep` über `internal/`: kein `crossNullGuard`,
  `loadCrossSource`). Keine Prüfung, die der Vertrag nicht verlangt.
- **Hexagon / §3.2 / gofmt:** ohne Befund — keine neuen Importe, keine
  `//nolint`-Direktive; `make arch-check` 0 Befunde; `gofmt -l` listet keine
  Datei dieses Slice (die vier `app`-Treffer sind Import-Gruppierungs-Artefakte
  und Bestand).

### Mutations-Härte (Frage 4)

- **Die drei benannten Mutationen kippen je genau einen Test — bestätigt.**
  (i) `excludeReqs` zurück hinter `crossVacuity` (bzw. in `diffViews`) ⇒ nur
  `TestCrossConsistencyVakuumDurchUebergriffigesVentil` kippt; (ii)
  `readCrossFile`/`spanCrossSource` wieder je Sicht bündeln ⇒ nur
  `TestCrossConsistencyPraezedenzQuelleVorAbschnitt` kippt (R3-F-3 geschlossen);
  (iii) `count != 1` → `count == 0` ⇒ nur
  `TestCrossConsistencyBackwardIDHeaderDoppelt` kippt (R3-F-4 geschlossen; die
  Fixture `| Kennung | Bezug | Kennung |` bindet `Bezug` sauber und lässt allein
  `backwardIDColumn` fehlschlagen, Assertion auf `"2-mal"`).
- **`vacuityHint`-Verzweigung beidseitig gepinnt:**
  `TestCrossConsistencyVakuumBeideSichtenLeer:320-323` assertiert **NICHT**-
  Vorkommen von `"exclude-req"` ohne gesetztes Ventil,
  `…VakuumDurchUebergriffigesVentil:447-450` das Vorkommen mit. Die
  Fehldiagnose-Klasse ist damit in beide Richtungen belegt — genau die
  Zusicherung, die der Commit beansprucht. (Restlücke: der Hint-**Inhalt**,
  R4-F-4.)
- **Ventil-Gegenprobe:** `…VakuumDurchUebergriffigesVentil:453-460` pinnt, dass
  ein eng gefasstes `exclude-req` den Drift stehen lässt (2 Befunde) — ohne sie
  wäre „Ventil ⇒ Vakuum" trivial durch ein kaputtes `excludeReqs` erfüllbar.
- **Vakuum-Zweige (a)/(b) und die Boundary bleiben gepinnt:** die R3-Analyse
  gilt unverändert; `TestCrossConsistencyLeereVorwaertsSichtIstKeinVakuum` und
  die `equal`-Gegenprobe in `…VakuumRueckSichtLeerUnterSuperset` halten beide
  `&&`-Konjunktionen fest.

### Vertragstreue des Gesamt-Slice (Frage 5)

- **Erneut über den vollen Diff `6c4ccf5..HEAD` geprüft, ohne neuen Befund:**
  Byte-Identität ohne Block in allen drei Ausgabeformaten (Markdown/JSON/YAML,
  `omitempty` auf nil-Slice, `CrossActive` nicht serialisiert); Gate-Bindung
  ausschließlich über das globale `--require-complete` mit getrennter
  Ursachen-Meldung; Modi `equal`/`superset`; `artifact-id-column`
  `first`-Sentinel **und** Header-Name, beidseitig fail-closed; geteiltes
  `design-pattern` als Namensraum-Mechanik; Sortierung `(R, Artefakt,
  Richtung)`; Config-Validierung inkl. `KnownFields(true)`, Pflichtfeldern,
  Repo-Wurzel-Bindung, `mode`-Enum, RE2-Kompilierung; `--print-config`-Vorlage
  auskommentiert und schema-deckend;
  [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  (reiner Lese-Pfad, alle Repros netzlos/read-only); Hexagon-Import-Regeln.
- **Alle R1–R3-Findings geschlossen:** R1-F-1…F-4, R2-F-1…F-4, R3-F-1…F-4 —
  jedes einzeln gegen `d-check:latest` bzw. per Code-Lesung nachgeprüft. Kein
  Silent-Green-Pfad ist im Kreuzverweis-Abgleich verblieben; die vier
  historischen Repros (Namensraum-Fehlgriff, `superset`/`B=∅`, partieller
  ID-Header-Verlust, übergriffiges Ventil) liefern alle Exit 2, der
  Bootstrap-Zustand weiterhin Exit 1 mit zwei lauten `B\F`-Befunden.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 2 | R4-F-1, R4-F-2 |
| LOW | 2 | R4-F-3, R4-F-4 |
| INFO | 0 | — |

## Verdikt

**REJECT** — und zum ersten Mal in dieser Kette **ohne jeden Vorbehalt gegen den
Code**.

Der Kern ist erledigt. Die Auflösung zugunsten post-Ausschluss ist die richtige,
sie ist besser begründet als meine R3-Neigung, und der Code setzt sie exakt um:
Spezifikations-Schrittfolge und Fehlerpräzedenz sagen jetzt dasselbe, die
Stufenfolge `excludeReqs → crossVacuity → diffViews` ist die des Vertrags,
`diffViews` ist von der Ventil-Logik befreit, alle zwölf Findings aus R1–R3 sind
geschlossen und die Mutations-Zusicherungen halten der Nachprüfung stand. Auch
die unaufgeforderte Diagnose-Schärfung ist die richtige Reaktion auf die richtige
Einsicht — dass eine Fehldiagnose dieselbe Klasse ist wie ein falscher Exit-Code.

Blockierend sind zwei **Dokumentations**-Befunde, und beide sind ein einziger
Editier-Durchgang:

**R4-F-1 — ja, der Bump auf 0.44.2 ist nötig; das ist meine klare Antwort auf
Frage 1.** Die Begründung „Rang 1 definiert den Zustand, nicht den Messpunkt"
trägt für die *Definition* — sie ist messpunkt-frei, und Rang 2 darf sie
präzisieren. Sie trägt **nicht** für die *Akzeptanzkriterien*: die sind in beiden
Zweigen ursachen-gefasst, und ein fehlerfreies Muster mit übergriffigem Ventil
ist keiner der beiden Given-Zustände. Die AK sind die abnahmebindende Fläche —
was dort nicht steht, wird bei der Abnahme nicht geprüft und ist für den
Konsumenten nicht herleitbar. Das entscheidende Argument liefert die eigene
Präzedenz: 0.44.1 wurde für **genau diese Klasse** gebumpt und gab ihr Negative
**und** Boundary. Ein zweiter, ursächlich eigenständiger Trigger derselben Klasse
ohne AK ist die Inkonsistenz. Der Aufwand ist minimal — ein AK-Zweig, ein Satz im
Ableitungssprung-Absatz, eine 0.44.2-Historie-Zeile; der Code bleibt unberührt.

**R4-F-2** ist die zweite Hälfte desselben Gedankens: die Entscheidung, die einen
Review-Dissens überstimmt hat, steht in der Spezifikation und in einer
Commit-Message, aber nicht in der ADR — obwohl
[`AGENTS.md`](../../AGENTS.md) §3.4 die Begründungen dort verortet und ADR-0038
für exakt diesen Zweck schon einmal (Entscheidung 8) nachgetragen wurde. Das
Fenster ist eng: mit `Accepted` wird die ADR immutabel, und danach kostet die
Nachbesserung eine eigene `Supersedes`-ADR. Entscheidung 5 („kein gelöstes
Problem") verdient im selben Zug die Präzisierung, dass der Totalfall des Ventils
seit diesem Slice geguardet ist und nur die Teil-Drift offen bleibt.

R4-F-3/F-4 sind Nits ohne Blockade-Anspruch; R4-F-3 ist der ärgerlichere, weil er
ausgerechnet die Meldung trifft, die dieser Commit ehrlich machen wollte, und
weil `9027514` sie noch korrekt formulierte.
