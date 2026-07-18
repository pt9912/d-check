# Review — slice-071 Implementierung R2 (Fix-Review zu R1)

**Datum:** 2026-07-17
**Review-Art:** unabhängiger Fix-Review (kontext-getrennt; adversarial, nicht
bestätigend)
**Gegenstand:**
[`slice-071`](../plan/planning/done/slice-071-trace-cross-consistency-gate.md) —
Fix-Commit `d11398f` (Einarbeitung der R1-Befunde F-1…F-4); Gesamt-Diff des Slice
`6c4ccf5..HEAD`
**Vor-Review:** [R1](2026-07-17-slice-071-implementation-r1.md) (Verdikt REJECT)
**Reviewer:** Claude (kontext-unabhängiger Lauf)
**Skill:** `.harness/skills/reviewer.md` v1.2.0
**Modell:** Opus 4.8 (1M context)

## Eingangs-Kontext

- Behauptete Fixes: F-1 (`crossNullGuard`), F-2 (`bindBackwardTables` +
  `backwardIDColumn`), F-3 (Phasen-Umbau), F-4 (gofmt); F-5 bewusst offen
- Lastenheft
  [`DC-FA-XREF-001`](../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
  (§Akzeptanzkriterien, §fail-closed-Enumeration),
  [`DC-FA-CLI-011`](../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code),
  [`DC-FA-REQ-001`](../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen),
  [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)
- Spezifikation
  [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
  (Schritte 2/3/4/5 + §Fehlerpräzedenz),
  [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
  (§Fehlerpräzedenz — Vergleichsanker)
- [ADR-0038](../plan/adr/0038-trace-cross-consistency.md),
  [ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md)
- Hard Rules [`AGENTS.md`](../../AGENTS.md) §3
- Prüfgegenstand: `internal/hexagon/core/app/trace_cross.go`,
  `internal/hexagon/core/app/trace_cross_test.go`

**Ausgeführte Sensoren:** `make lint` (0 issues), `make test` (grün),
`make arch-check` (0 Befunde), `make build` + Repro-Läufe gegen `d-check:latest`
(`--network none`, `:ro`), `gofmt -l` in der Toolchain.

## Antwort auf die Vertragsfrage: greift `crossNullGuard` zu Recht bei **einer** leeren Sicht?

Die Unsicherheit ist berechtigt — die Antwort lautet **nein, nicht symmetrisch**.
Der Vertrag trägt die Per-Sicht-Anwendung nicht:

1. **Spec Schritt 5 definiert den Diff über die Vereinigung:** „Für jede
   Anforderung `R ∈ keys(F) ∪ keys(B)`". Mit `F = ∅` ist das wohldefiniert und
   liefert `keys(B)` — jede B-Kante wird zu einem `B(R) \ F(R)`-Befund. Das ist
   **kein** degenerierter Zustand, sondern ein Ergebnis: „die RTM nennt nichts,
   die Rück-Kanten nennen alles".
2. **Die Fehlerpräzedenz kennt keine Nullmengen-Stufe:** „Config-Schema →
   Quellen lesen → Header-Bindung → Range-Expansion → Diff". Zum Vergleich führt
   [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
   ausdrücklich „… → Duplicate-ID → **Nullmenge** → Referenz-/Coverage-Scans".
   Die Auslassung bei `XREF` ist damit sprechend, nicht vergesslich.
3. **Die fail-closed-Enumeration des Lastenhefts** nennt „ungültiges Regex,
   fehlende Spalte, ID-Header nicht genau einmal, unbekannter `mode`" — keine
   Nullmenge. Slice §2 wiederholt genau diese Liste.
4. **Schritt 3 verlangt keinen Laufzeit-Check.** Der Satz lautet: „Vorwärts- und
   Rückwärts-Artefakt-IDs werden bewusst mit **demselben** `forward.design-pattern`
   extrahiert und liegen **damit** im selben Namensraum — sonst wäre der
   Mengen-Diff inhärent leer/voll und bedeutungslos." Das ist die **Begründung
   dafür, das Muster zu teilen** (das „sonst" ist kontrafaktisch: *hätte* man zwei
   Muster). Die Kongruenz ist durch die Konstruktion zugesichert, nicht als
   Prüfpflicht formuliert.
5. **[ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md) trägt
   die Analogie nur halb.** Dort ist der Guard **bedingt**
   (`strictSource = expliziter source ∨ format == table`) und betrifft die
   Anforderungs**quelle**, wo eine Nullmenge die RTM tatsächlich sinnlos macht.
   Hier macht `F = ∅` den Diff nicht sinnlos — es macht ihn maximal aussagekräftig.

**Welches Verhalten ist richtig?** Die zu schließende Lücke ist nicht „eine Sicht
ist leer", sondern „**der Diff ist vakuum und behauptet trotzdem Konsistenz**".
Vakuum ist er genau dann, wenn er null Befunde aus null Vergleichsmaterial
erzeugt:

| Zustand | `equal` | `superset` | vakuum? |
|---|---|---|---|
| `F = ∅`, `B = ∅` | 0 Befunde, Exit 0 | 0 Befunde, Exit 0 | **ja** (R1-F-1) |
| `F ≠ ∅`, `B = ∅` | `|F|` Befunde (laut) | 0 Befunde, Exit 0 | **ja** unter `superset` |
| `F = ∅`, `B ≠ ∅` | `|B|` Befunde (laut) | `|B|` Befunde (laut) | **nein** |
| `F ≠ ∅`, `B ≠ ∅` | laut/grün nach Daten | laut/grün nach Daten | nein |

Die dritte Zeile ist der Fall, den der Fix zu Unrecht mitnimmt. Der Guard müsste
an der **Vakuität** (und damit am Modus) hängen, nicht an der einzelnen Sicht.
Details als Finding R2-F-1.

## Findings

### R2-F-1 — `crossNullGuard` weist eine vertraglich definierte, meldenswerte Lage ab (`F = ∅`, `B ≠ ∅`) — mit falscher Ursachen-Zuschreibung

- **kategorie:** MEDIUM
- **quelle:** [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
  Schritt 5 + §Fehlerpräzedenz;
  [`DC-FA-XREF-001`](../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
  (§fail-closed-Enumeration, §Akzeptanzkriterium „Superset-Modus");
  [ADR-0038](../plan/adr/0038-trace-cross-consistency.md) Entscheidung 3 + 7;
  [`slice-071`](../plan/planning/done/slice-071-trace-cross-consistency-gate.md) §4
- **pfad:** `internal/hexagon/core/app/trace_cross.go:230-233` (`forwardEdges`
  ruft `crossNullGuard` unbedingt), `:254-257` (`backwardEdges` dito), `:259-270`
  (`crossNullGuard`)
- **befund:** Der Guard feuert je Sicht, sobald **eine** von beiden kantenleer
  ist, ohne Kenntnis von `mode` oder Befundmenge. Damit ist `F = ∅` bei
  gepflegten Rück-Kanten Exit 2, obwohl Schritt 5 diesen Zustand über
  `keys(F) ∪ keys(B)` definiert und `|B|` `B\F`-Befunde vorschreibt. Betroffen
  ist genau der Bootstrap-Zustand, den
  [`slice-071`](../plan/planning/done/slice-071-trace-cross-consistency-gate.md)
  §4 als Konsumenten-Vorarbeit benennt („die Vorwärts-Sicht (§27.1) muss auf
  konkrete IDs restrukturiert werden (Wildcards/Prosa raus), bevor das Gate grün
  wird") — dort ist die Befundliste die Migrations-Checkliste und laut
  [ADR-0038](../plan/adr/0038-trace-cross-consistency.md) Entscheidung 7 zugleich
  der Input des späteren Generators. Erschwerend: die Meldung schreibt die
  Ursache der **Config** zu („prüfe, dass req-column/design-column die Rollen
  treffen und design-pattern den Artefakt-Namensraum"), obwohl die Config korrekt
  und allein die Datenlage leer ist — eine Fehldiagnose, die den Nutzer an der
  falschen Stelle suchen lässt.
- **verifizierbar:** ja — reproduziert gegen `d-check:latest`. Repo mit
  unrestrukturierter Vorwärts-Zeile (`| GG-ARCH-006 | alle Scheduler-Komponenten
  (siehe Architektur) |`) und zwei gepflegten Rück-Kanten (`GG-AR-COMP-SCHED`,
  `GG-AR-P-005` → `GG-ARCH-006`), Config fehlerfrei:

  ```text
  d-check: error: trace.cross-consistency.forward: die gebundenen Tabellen in
  "docs/traceability.md" ergaben 0 Kanten — prüfe, dass req-column/design-column
  die Rollen treffen und design-pattern den Artefakt-Namensraum
  EXIT=2
  ```

  Vertragskonform wären zwei `B\F`-Befunde und Exit 1 unter `--require-complete`.

### R2-F-2 — Rest-Fehlerpräzedenz: der Sektionsnamen-Guard der Vorwärts-Sicht schlägt vor der Existenzprüfung der Rückwärts-Datei zu

- **kategorie:** LOW
- **quelle:** [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
  §Fehlerpräzedenz („→ Quellen lesen (fehlende `forward.file`/`backward.file` ⇒
  Exit 2) → …")
- **pfad:** `internal/hexagon/core/app/trace_cross.go:77-84`, `:123-135`
  (`loadCrossSource`)
- **befund:** Der F-3-Fix staffelt Bindung und Range-Expansion korrekt über beide
  Sichten, bündelt aber Existenzprüfung **und** `checkCrossSections` in einen
  Aufruf pro Sicht. Dadurch entsteht innerhalb der Stufe „Quellen lesen" dieselbe
  sicht-orientierte Reihenfolge, die R1-F-3 beanstandete: `forward`-Sektionsname
  ohne Heading-Treffer wird vor `backward.file: existiert nicht` gemeldet. Die
  Präzedenz benennt als Inhalt dieser Stufe ausdrücklich nur die
  Datei-Existenz beider Sichten; vor dem Fix liefen beide Existenzprüfungen
  (`readCrossFile` ×2) geschlossen voran, sodass die Reihenfolge hier gegenüber
  `1bfa7f8` neu ist. Beobachtbar allein an der Meldung (beide Fälle Exit 2).
- **verifizierbar:** ja — Config mit `forward.sections: ["27.1"]` (Kurzform-
  Tippfehler) und `backward.file: spec/fehlt-komplett.md`; gemeldet wird
  `Abschnitt "27.1" trifft keine Überschrift`, nach Vertrag die fehlende Datei.

### R2-F-3 — Kein Test pinnt eine der beiden `crossNullGuard`-Aufrufstellen einzeln

- **kategorie:** LOW
- **quelle:** Maintainability; Reviewer-Skill (Negativtest-Anker bei neuem
  öffentlichen Vertrag)
- **pfad:** `internal/hexagon/core/app/trace_cross_test.go:307-320`
  (`TestCrossConsistencyNamensraumVorbedingung`)
- **befund:** Die Fixture des einzigen Guard-Tests macht **beide** Sichten
  kantenleer (das verstellte `DesignPattern` greift auch an der Rück-ID-Spalte
  vorbei, da beide dasselbe Muster teilen). Der Test prüft `err != nil` +
  `"0 Kanten"` — beides bleibt erfüllt, wenn man den `crossNullGuard`-Aufruf in
  `forwardEdges` **oder** den in `backwardEdges` einzeln entfernt: der jeweils
  andere feuert und liefert denselben Textbaustein. Nur das Entfernen **beider**
  Aufrufe bricht die Suite. Der Rück-Sicht-Guard feuert damit in keinem Test; die
  Zeilen-Coverage verdeckt das, weil `crossNullGuard` über den Vorwärts-Pfad als
  abgedeckt zählt.
- **verifizierbar:** ja — `crossNullGuard`-Aufruf in `backwardEdges` entfernen und
  `make test` laufen lassen: grün. (Der Pfad ist nicht tot — gegen
  `d-check:latest` reproduziert: Vorwärts-Kanten vorhanden, Rück-Erste-Spalte
  ohne Artefakt-ID, `mode: superset` ⇒ Exit 2 aus `backwardEdges`.)

### R2-F-4 — Der F-2-Regressionstest reproduziert die F-2-Fehlgestalt nicht und hängt am Meldungstext

- **kategorie:** LOW
- **quelle:** Maintainability; [R1](2026-07-17-slice-071-implementation-r1.md) F-2
- **pfad:** `internal/hexagon/core/app/trace_cross_test.go:328-345`
  (`TestCrossConsistencyBackwardIDHeaderFehlt`)
- **befund:** Die Fixture führt **eine** Rück-Tabelle ohne den ID-Header. R1-F-2
  war aber der **partielle** Fall: eine Tabelle bindet, eine zweite (mit `Bezug`,
  ohne den ID-Header) wird still verworfen — nur dort entstand Exit 0 statt
  Exit 2. Mit der Ein-Tabellen-Fixture liefert auch der **alte** Code Exit 2
  (über den `found`-Guard, „keine Tabelle mit dem konfigurierten Header
  \"Bezug\""); der Test unterscheidet alt und neu allein daran, dass die Assertion
  auf den Teilstring `"Kennung"` prüft, den die alte Meldung zufällig nicht
  enthält. Eine spätere Umformulierung der `found`-Meldung (etwa: alle
  konfigurierten Header nennen) ließe den Test unter wieder-eingeführtem
  Silent-Skip grün laufen.
- **verifizierbar:** ja — Fixture um eine zweite, bindende Tabelle
  (`| Kennung | Bezug |`) ergänzen: der Test bleibt grün, deckt aber dann den
  Fall ab, den R1-F-2 beschrieb.

## Negativbefunde (geprüft, ohne Befund)

- **R1-F-1 Original-Repro geschlossen:** gegen `d-check:latest` verifiziert.
  Reeller Drift + `design-pattern: 'GG-ARCH-COMP-[A-Z]+'` (Tippfehler-Klasse) →
  Exit 2 mit `… ergaben 0 Kanten …` statt `0 Differenz(en).`/Exit 0. Der von R1
  reproduzierte Harness-Lüge-Pfad existiert nicht mehr.
- **R1-F-2 Original-Repro geschlossen:** verifiziert. Zwei-Tabellen-Fixture
  (`| Kennung | Bezug |` + `| Port-ID | Bezug |`), `artifact-id-column: Kennung`,
  `mode: superset` → `trace.cross-consistency.backward: Tabelle ab Zeile 11:
  konfigurierter Header "Kennung" kommt 0-mal vor (genau einmal erwartet)`,
  Exit 2. Die Relevanz entsteht in `bindBackwardTables:168` tatsächlich allein
  über `bc.EdgeColumn`; `backwardIDColumn:186-207` löst danach auf und ist für
  `count == 0` **und** `count > 1` fail-closed — die Spec-Formulierung „genau
  einmal" ist damit beidseitig getroffen.
- **Superset-Silent-Green (R1, Rand des F-1-Komplexes) geschlossen:** der
  Rück-Sicht-Guard ist erreichbar und feuert. `F ≠ ∅`, `B = ∅` (erste Spalte trägt
  Klartext statt ID) unter `mode: superset` ergab vor dem Fix `0 Differenz(en).`/
  Exit 0, jetzt Exit 2. (Testlücke dazu: R2-F-3.)
- **R1-F-4 (gofmt) geschlossen:** `gofmt -l internal/hexagon/core/app/` listet
  `trace_cross_test.go` nicht mehr; die vier verbliebenen Treffer
  (`diagnose.go`, `diagnose_test.go`, `repair.go`, `suggest.go`) sind
  Import-Gruppierungs-Artefakte des bloßen `gofmt` und Bestand vor diesem Slice.
- **Phasen-Umbau (F-3) — Bind-/Edge-Staffelung:** ohne Befund.
  `crossConsistency:77-101` durchläuft `loadCrossSource` ×2 → `bindForwardTables`
  → `bindBackwardTables` → `forwardEdges` → `backwardEdges` → `diffViews`; ein
  Rück-Sicht-Header-Defekt rangiert damit korrekt vor einem Vorwärts-Range-Defekt
  (die R1-F-3-Beobachtung ist auf Bind-/Range-Ebene geschlossen; der Rest liegt
  eine Stufe tiefer, R2-F-2).
- **`crossBadRow` in der Bind-Phase:** semantisch unverändert. Der Aufruf steht in
  beiden Bind-Funktionen **nach** der Relevanz-Prüfung (`:148`, `:179`), also
  bleibt ein Zellenzahl-Bruch in einer **nicht**-relevanten Tabelle folgenlos
  (bloßes Tabellenende) und in einer relevanten fail-closed — identisch zu
  `1bfa7f8` und zum `DC-FA-REQ-001`-Pendant. Die Einordnung in die Bind-Phase
  widerspricht der Präzedenz nicht: der Zellenzahl-Bruch ist dort gar nicht
  gelistet.
- **Relevanz-Semantik vorwärts:** unverändert (`bindCrossColumns` mit beiden
  Namen; alle Namen präsent ⇒ relevant). Der `len(out) == 0`-Guard ersetzt das
  alte `found`-Flag 1:1, inklusive Meldungstext.
- **Determinismus ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)):**
  ohne Befund. `bindForwardTables`/`bindBackwardTables` füllen `out` in
  `markdownTables`-Dokumentreihenfolge; `forwardEdges`/`backwardEdges` iterieren
  Tabellen- dann Zeilen-geordnet, `crossView.add` behält die erste Fundstelle.
  `backwardIDColumn` nimmt bei Mehrfach-Treffern den ersten Index und meldet
  deterministisch mit Zähler. `crossNullGuard`-Meldung ist konstant. Keine
  unsortierte Map-Iteration im neuen Code.
- **Index-Sicherheit `boundTable`:** ohne Befund. `primary`/`secondary` stammen aus
  Header-Indizes bzw. dem `first`-Sentinel (0); `t.header` hat garantiert ≥1 Zelle
  (`isTableDelimiter` verlangt `len(cells) > 0`), und `consumeTableRows` nimmt nur
  Zeilen mit `len(cells) == len(t.header)` auf — `row.cells[bt.primary]` kann nicht
  out-of-range laufen. Die rollen-agnostischen Feldnamen sind je Sicht konsistent
  belegt (vorwärts req/design, rückwärts ID/edge) und dokumentiert.
- **Tote Pfade:** ohne Befund. `bindBackwardColumns` ist restlos entfernt (kein
  Aufrufer verblieben); `bindCrossColumns` wird weiter variadisch genutzt (2 Namen
  vorwärts, 1 Name rückwärts), beide Fehlerzweige erreichbar. Kein toter Zweig
  neu entstanden.
- **`exclude-req`-Reihenfolge:** vertragstreu — die Ausnahme greift in `diffViews`
  **nach** der Extraktion, wie Spec Schritt 4 („vor dem Diff aus den Schlüsselmengen
  von `F` und `B` entfernt") es ordnet. Dass ein weit gefasstes `exclude-req` alle
  Befunde tilgen kann, ist das benannte Ventil
  ([ADR-0038](../plan/adr/0038-trace-cross-consistency.md) Entscheidung 5, „kein
  gelöstes Problem") und kein Defekt der Umsetzung.
- **Erfundener Vertrag (übrige Prüfungen):** außer R2-F-1 ohne Befund. Insbesondere
  `backwardIDColumn` prüft nur, was Spec und Lastenheft verlangen („genau einmal"),
  und der `first`-Sentinel bleibt prüfungsfrei (korrekt — er ist positionell).
- **Mutations-Härte der Gegenprobe:** `TestCrossConsistencyBackwardIDHeaderBenannt`
  ist **hart**. Die ID-Spalte steht bewusst an Position 2 (`| Schicht | Kennung |
  Bezug |`); gäbe `backwardIDColumn` fälschlich 0 zurück, träfe `design-pattern`
  die Zelle `Kern` nicht, die Rück-Sicht bliebe kantenleer und der Null-Guard
  brächte den Test zu Fall. Der Test kann sich nicht selbst grün testen.
- **Byte-Identität ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus))
  und Gate-Bindung:** unberührt — `d11398f` fasst weder `diffViews`, `report.go`,
  `cli.go`, `model/config.go` noch `configyaml.go` an. Die R1-Negativbefunde dazu
  gelten unverändert; `TestCLI071_Cross_DefaultByteIdentisch` bleibt grün.
- **Hexagon-Import-Regeln ([ADR-0005](../plan/adr/0005-modul-layout-hexagon-ordner.md),
  [ADR-0012](../plan/adr/0012-kern-paketschnitt-model-rules-app.md)):** ohne
  Befund — keine neuen Importe; `make arch-check` (a-check, digest-gepinnt):
  0 Befunde.
- **[`AGENTS.md`](../../AGENTS.md) §3.2 / §3.6:** ohne Befund — keine
  `//nolint`-Direktive, keine `.golangci.yml`-Ausnahme, keine Schwellen-Senkung.
- **[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit):**
  ohne Befund — `loadCrossSource` bleibt reiner Lese-Pfad über den
  Filesystem-Port; alle Repro-Läufe mit `--network none` und `:ro`.
- **F-5 (bewusst offen):** akzeptiert als Won't-Fix-Designnotiz. Die Einordnung
  „Fläche erweitern ginge über den Vertrag hinaus" trifft zu; die Zusage, es in
  der Nutzerdoku zu benennen, ist der richtige Ort. Kein Finding.

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 1 | R2-F-1 |
| LOW | 3 | R2-F-2, R2-F-3, R2-F-4 |
| INFO | 0 | — |

## Verdikt

**REJECT** — knapp, mit deutlich verengtem Umfang gegenüber R1.

Was gelungen ist, sei ausdrücklich festgehalten: **beide blockierenden
R1-Befunde sind geschlossen**, und zwar an der Wurzel, nicht am Symptom. Kein
Stilles-Grün-Pfad ist im Kreuzverweis-Abgleich verblieben — auch der von R1 nur
am Rand benannte `superset`/`B = ∅`-Fall ist jetzt fail-closed. Der F-2-Fix trifft
die Spec-Formulierung („zählt jede Tabelle mit einem `edge-column`-Header")
wörtlich und ist beidseitig fail-closed. Der Phasen-Umbau ist auf Bind-/Range-
Ebene sauber und hat weder Relevanz-, `badLine`- noch Determinismus-Semantik
gebrochen.

Blockierend bleibt allein **R2-F-1**: der F-1-Fix hat über das Ziel hinaus
geschossen. R1 hat einen Defekt benannt („der Diff behauptet Konsistenz, ohne je
verglichen zu haben"), nicht ein Mittel; das gewählte Mittel — je Sicht ≥1 Kante
verlangen — schließt den Defekt, nimmt aber den vertraglich definierten Fall
`F = ∅ ∧ B ≠ ∅` mit, in dem der Abgleich sein **wertvollstes** Ergebnis liefern
würde: die vollständige Liste der Kanten, die die RTM verschweigt. Genau dieser
Zustand ist der dokumentierte Einstiegspunkt des Konsumenten
([`slice-071`](../plan/planning/done/slice-071-trace-cross-consistency-gate.md)
§4) und der Input des in
[ADR-0038](../plan/adr/0038-trace-cross-consistency.md) Entscheidung 7
vorgesehenen Generators. Verschärfend ist nicht der Exit-Code allein, sondern die
**Fehldiagnose**: die Meldung schickt den Nutzer in die Config, während die Config
korrekt ist. Ein Gate, das bei intakter Konfiguration „prüfe deine Konfiguration"
sagt, verbrennt genau das Vertrauen, das die R1-Fixes wiederhergestellt haben.
Die vollständige Vertragsherleitung samt Vakuitäts-Tabelle steht oben; sie ist die
Grundlage, auf der die Guard-Bedingung neu zu fassen ist — die Entscheidung
gehört, wenn sie vom Vertrag abweichen soll, nach
[`AGENTS.md`](../../AGENTS.md) §3.6 in einen ADR-Nachtrag, nicht in den Code.

R2-F-2/F-3/F-4 sind Nits ohne Blockade-Anspruch, aber R2-F-3 ist der
erwähnenswerteste davon: die Zeilen-Coverage (93,9 %) meldet `crossNullGuard` als
abgedeckt, während keine der beiden Aufrufstellen einzeln gepinnt ist — Coverage
als Boden, nicht als Decke, in Reinform.
