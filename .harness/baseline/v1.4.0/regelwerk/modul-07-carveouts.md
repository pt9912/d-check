## Modul 7 — Carveout Management

*Quelle: [02-planung/modul-07-carveouts.md](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/02-planung/modul-07-carveouts.md)*

### Harness-Einordnung (Modul 7)

Carveout-Pflege ist ein Pfeiler von *Entropy Management* (siehe
[`klassifikation.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/klassifikation.md)):
ein Carveout-Audit pro Welle verhindert, dass temporäre Ausnahmen zu
permanenten Lügen werden.

### Kernidee (Modul 7)

Jeder temporäre Carveout benötigt einen Plan. Ein Carveout ohne
Auflösungs-Trigger ist ein permanenter Carveout, der lügt.

### Worked Example A: einen Carveout dokumentieren

**Ausgangssituation:** Das Coverage-Gate `coverage-gate-critical` ist
rot. Der Index-Layer (`internal/index/`) hat 76 % statt der geforderten
90 %. Grund: Binär-Format-Parser mit Fehlerpfaden (`E099` bei korrupter
Datei), die nur partiell durch Unit-Tests abgedeckt sind. Eine
Property-Test-Suite wird die verbleibenden Pfade abdecken — ist aber
erst in Welle 2 eingeplant.

Die Versuchung: das Gate in der CI-Konfiguration herunterdrehen. Das
ist eine *stille* Senkung — sie taucht weder in `harness/README.md`
noch in der Spec auf. Der bessere Weg: ein Carveout.

**Schritt 1 — Carveout-Datei anlegen.** Konvention:
`docs/plan/carveouts/CO-<NNN>-<kurztitel>.md`. ID läuft in `CO-*`-Reihe
(separat von `LH-`, `ADR-`, `SL-` — siehe
[`konventionen.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/konventionen.md#id-schema-als-klammer)).
Für unseren Fall: `CO-001-index-coverage.md`.

**Schritt 2 — Pflichtfelder im Frontmatter / Header festlegen.** Ein
temporärer Carveout, der nicht heimlich permanent werden soll, braucht
sechs Felder:

```markdown
# CO-001: Bootstrap-Coverage `internal/index/`

**Status:** Aktiv.
**Datum angelegt:** 2026-05-20. **Letzte Prüfung:** 2026-06-01.
**Betroffenes Gate:** `coverage-gate-critical`.
**Geltungsbereich:** `internal/index/` (Index-Layer, alle Sprachen).
**Folge-Slice:** [`slice-013-property-tests.md`](../planning/in-progress/slice-013-property-tests.md)
```

Wenn `Folge-Slice` fehlt oder leer ist, ist der Carveout *de facto*
permanent — und gehört dann offen so markiert (siehe Schritt 6).

**Schritt 3 — Auflösungs-Trigger als beobachtbare Bedingung
formulieren.** Anti-Form: *"sobald wir Zeit haben"*, *"nach dem nächsten
Refactoring"*, *"wenn das Team Kapazität hat"*. Gute Form: eine
Bedingung, die ein anderer Mensch ohne Rückfrage als eingetreten oder
nicht eingetreten beurteilen kann.

```markdown
## Auflösungs-Trigger

Welle 2 (welle-2-qualitaet) done — Property-Test-Suite läuft 100
Generationen und deckt die Fehlerpfade.

Konkret: `internal/index/`-Coverage erreicht ≥ 90 %, geprüft in
`make coverage-gate-critical` ohne Ausnahmen.
```

Zwei Sätze: einer für den Welle-Bezug (Roadmap-Anker), einer für die
*messbare* Schwelle. Die messbare Schwelle ist der wichtigere — sie
ist es, was die CI prüft.

**Schritt 4 — Geltungs-Konfiguration mit ID-Kommentar verdrahten.** Die
Gate-Konfiguration *zeigt* auf den Carveout, damit der Carveout im
`make gates`-Output nicht versteckt ist:

```diff
  # <sprache>/coverage.config
  critical_paths:
-   exceptions: []
+   exceptions:
+     - "internal/index/"  # CO-001 — bis Welle 2 done
```

Der `# CO-001`-Kommentar ist nicht Kosmetik: er ist die Brücke zwischen
Gate-Konfiguration und Carveout-Datei. Ohne ihn weiß niemand, *warum*
diese Pfad-Ausnahme existiert.

**Schritt 5 — Verifikations-Checkliste für den Auflösungs-Zeitpunkt
hinterlegen.** Damit nach Trigger-Eintritt klar ist, was zu tun ist:

```markdown
## Verifikation (nach Auflösung)

- [ ] `internal/index/`-Coverage in allen Sprach-Skeletten ≥ 90 %.
- [ ] Carveout-Konfiguration aus Coverage-Config entfernt.
- [ ] `make coverage-gate-critical` grün ohne Ausnahmen.
- [ ] Diese Datei nach `done/CO-001-index-coverage.md` bewegt (reiner `git mv`).
- [ ] slice-013 Closure-Notiz schließt diese Auflösung mit ein.
```

Vier Häkchen, eines davon ein `git mv`. Auflösung ohne Verschiebung in
`done/` ist eine zweite Lüge — der Carveout wirkt "aufgelöst", liegt
aber weiter im aktiven Verzeichnis.

**Schritt 6 — Carveout, BF-Sub-Area-Markierung oder ADR?** Bevor du
Schritt 1–5 als endgültige Form annimmst, prüfe, ob das Werkzeug
*Carveout* überhaupt passt. Drei legitime Alternativen, getrennt durch
**zwei sequenzielle Wenn-Dann-Fragen** statt einen flachen
Drei-Wege-Vergleich — Granularität *vor* Temporalität (so vermeidest
du den Reflex, jede entdeckte Diskrepanz als Carveout zu führen):

1. **Granularität — einzelne Diskrepanz oder Cluster?** Trifft *eine*
   konkrete Gate-/Regelausnahme isoliert auf, oder zeigt sich ein
   Diskrepanz-Cluster im selben Geltungsbereich (mehrere Carveouts auf
   denselben Pfad/dieselbe Sub-Area) bzw. ein systemisches *"Code
   existiert vor Doku"*-Muster?

   - *Cluster oder Muster* → **BF-Sub-Area-Markierung mit
     Graduation-Plan** als Modus-Deklaration im Adaptions-Block von
     `harness/conventions.md` (Mechanik in
     [`konventionen.md` §Modus pro Sub-Area](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/konventionen.md#modus-pro-sub-area-greenfield-vs-brownfield)).
     Frage 2 entfällt — eine Sub-Area-weite Markierung ist eine
     andere Werkzeug-Klasse als ein punktueller Carveout.
   - *Einzelne Diskrepanz* → weiter zu Frage 2.

   *Kein harter Schwellwert für "Cluster"* — Faustregel, nicht Zahl:
   wer die Diskussion am Zähler aufhängt ("ab 3 Carveouts BF?"), hat
   das Symptom-Muster (gemeinsamer Geltungsbereich) mit einer Quote
   verwechselt.

2. **Temporalität — Trigger ernst zu erreichen?** Bei einer einzelnen
   Diskrepanz: wäre der Auflösungs-Trigger aus Schritt 3 mit
   absehbarem Aufwand zu erreichen, und stünde das in realistischem
   Verhältnis zum Nutzen?

   - *Ja* → **Carveout** (Schritt 1–5 wie eben durchgegangen — der
     Carveout dokumentiert die temporäre Ausnahme bis zum Trigger).
   - *Nein, ehrliche Antwort "nichts davon werden wir in absehbarer
     Zeit tun"* → **Permanent**, übergeführt in eine ADR. Permanente
     Carveouts gehören nicht in `carveouts/`, sondern als
     Architekturentscheidung mit Begründung in eine ADR.

Die folgende Tabelle bildet die drei Wahl-Pfade für eine *konkrete
Diskrepanz* ab. Bootstrap-aware Gate (Modul 13) erscheint hier
absichtlich nicht — es regelt *Gate-Reifestufung*, nicht
*Diskrepanz-Auflösung*. BF-Markierung wirkt eine Ebene höher als
Carveout und ADR: sie kippt den Sub-Area-Kontext, in dem die
Diskrepanz erst entsteht.

| Wahl                       | Symptom-Indikator                                                                                                  | Träger                            | Folge-Artefakt                                                                            |
|----------------------------|--------------------------------------------------------------------------------------------------------------------|-----------------------------------|-------------------------------------------------------------------------------------------|
| **Carveout**               | Eine konkrete Gate-/Regelausnahme, klar abgrenzbar, mit Folge-Slice und ernst erreichbarem Auflösungs-Trigger.     | einzelne Diskrepanz               | `docs/plan/carveouts/CO-<NNN>-*.md` (Schritt 1–5)                                          |
| **BF-Sub-Area-Markierung** | Diskrepanz-Cluster im selben Geltungsbereich, oder generelles *"Code-vor-Doku"*-Muster.                            | ganze Sub-Area                    | Modus-Deklaration im Adaptions-Block von `harness/conventions.md`, mit Graduation-Trigger |
| **ADR (permanent)**        | Trigger ist ehrlich nie zu erreichen — die Senkung ist Architekturentscheidung, nicht Übergang.                    | dauerhafte Architekturregel       | `docs/architecture/ADR-<NNNN>-*.md`; `Status: Permanent — übergeführt in ADR-<NNNN>`      |

Drei verwandte Begriffe waren hier nebeneinander im Spiel und meinen
Verschiedenes: *Disambiguierung* ist die kognitive Operation, die du in
Frage 1–2 gerade ausgeführt hast (Symptom auf Werkzeug abbilden);
*Werkzeug-Wahl* ist deren Ergebnis (Carveout, BF-Markierung oder ADR —
die linke Tabellenspalte); *Werkzeug-Klasse* ist die Achse, auf der sich
die drei unterscheiden (punktuell vs. Sub-Area-weit vs. dauerhaft — die
Träger-Spalte).

> **Hinweis zum Lab-Beispiel:** Das Lab unter
> [`lab/example/docs/plan/carveouts/`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/docs/plan/carveouts/)
> trägt heute nur den einzelnen `CO-001-index-coverage.md` — Frage 1
> führt dort folglich auf den Einzeldiskrepanz-Pfad und weiter zu
> Frage 2 (Trigger erreichbar — ja: Welle-2-Property-Tests). Ein
> echter Cluster entstünde, wenn zusätzlich `CO-002`/`CO-003` für
> Boundary-Tests und Type-Coverage auf demselben `internal/index/`-
> Pfad lägen; dann sprängen Frage 1 und Werkzeug-Wahl auf
> BF-Sub-Area-Markierung um. Die Markierungs-Mechanik selbst ist im
> Lab strukturell bereits vorhanden: [`lab/example/harness/conventions.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/harness/conventions.md)
> trägt einen `## Adaptions-Block` mit `MR-000` Baseline-Aussage und
> `MR-001` Source Precedence — die BF-Sub-Area-Markierung wäre ein
> neuer `MR-NNN`-Eintrag im selben Block. Konkret-Format:
>
> ```markdown
> ### MR-002 — `internal/index/`-Sub-Area im Brownfield-Modus
>
> **Modus:** Brownfield bis Welle-3-Graduation.
> **Geltungsbereich:** `internal/index/` (Index-Layer, alle Sprach-Skelette).
> **Graduation-Trigger:** Property-Test-Suite läuft + Coverage ≥ 90 %
> über alle Pfade.
> **Sync-Trigger:** nach Graduation einen Pointer-Eintrag in
> `harness/README.md` §Sensors, der die Sub-Area als GF-bewertet
> ausweist.
> **Folge-Slice (für Reconciliation):** `slice-014-bf-index-reconciliation.md` (legt die Inventur und die ersten Reconciliation-Häppchen fest).
> ```
>
> Das ersetzt eine Carveout-Kaskade (`CO-001`/`CO-002`/`CO-003` auf
> denselben Pfad) durch eine einzelne Sub-Area-weite Aussage mit
> klarem Graduations-Pfad.

**Was passiert mit dem Schritt-1–5-Entwurf, wenn der
Trichter nicht auf Carveout führt?** Der Inhalt ist nicht verloren, nur verschoben.
Trigger-Formulierung (Schritt 3), Geltungsbereichs-Präzision
(Schritt 2 Geltungsbereich-Feld), Verifikations-Checkliste (Schritt 5)
wandern in das gewählte andere Werkzeug:

- *BF-Sub-Area-Markierung*: der Geltungsbereich wird zum
  Sub-Area-Geltungsbereich, der Trigger wird zum Graduation-Trigger,
  die Verifikations-Checkliste zur Liste der
  Reconciliation-Akzeptanzkriterien.
- *ADR-Überführung*: der Geltungsbereich wird zum architektonischen
  Wirkungsbereich, der Trigger fällt weg (permanent), die
  Verifikations-Checkliste reduziert sich auf die Architektur-Folgen.

Der angelegte Schritt-1-Stub `CO-<NNN>-*.md` wird nicht aktiviert,
sondern entweder gelöscht (Inhalt vollständig in BF-Markierung/ADR
aufgegangen) oder mit einem `Status: Überführt in <Ziel>`-Header
nach `done/` verschoben, damit die Werkzeug-Wahl-Spur im Repo lesbar
bleibt.

Wenn die Wahl im Trichter auf "permanent" fällt, ist der
`Status:`-Wechsel:

```markdown
**Status:** Permanent — übergeführt in ADR-0009.
```

Das ist nicht Aufgabe; das ist Ehrlichkeit. Vergleich:
[`CO-001-index-coverage.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/lab/example/docs/plan/carveouts/CO-001-index-coverage.md).

### Drei Werkzeuge für gelockerte Gate-Disziplin (Modul 7)

**Carveout** = punktuelle Ausnahme für *einen* Fall, mit Folge-Slice
und Auflösungs-Trigger. **BF-Sub-Area-Markierung** = Sub-Area-weiter
Übergangs-Modus mit Graduation-Plan im Adaptions-Block von
`harness/conventions.md` — *Sub-Area-Kontext, kein Closure-Werkzeug*
(siehe [Modul 13 §Bootstrap-aware Gates](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/04-qualitaet/modul-13-quality-gates.md#bootstrap-aware-gates)).
**Bootstrap-aware Gate** = Stufung *des Gates selbst* (z. B. 40 % heute
→ 70 % bei M2).

Disambiguierungs-Reflex: bei Diskrepanz-Häufung im selben
Geltungsbereich ist die Wahl BF-Markierung, nicht Carveout-Kaskade —
wer ein Dutzend Carveouts für dieselbe Sub-Area anlegt, hat den
Mechanismus auf das falsche Werkzeug skaliert. Bootstrap-aware Gate
skaliert mit dem Repo; Carveout ist punktueller Vertrag. Verwechslung
jeder Achse führt zu "Bootstrap-Schlupfloch" — Stufung ohne Trigger ist
Carveout-Wildwuchs, Carveout-Kaskade ohne BF-Markierung ist
verschleierte Sub-Area-BF.

### Worked Example B: ein Carveout-Audit als wiederkehrenden Slice entwerfen

**Ausgangssituation:** Das Repo hat sechs Carveouts. Drei sind seit
über sechs Monaten "aktiv". Niemand hat sie kürzlich geprüft. Faktisch
sind sie permanent — aber im Repo lügen sie weiter unter `aktiv`. Dies
ist die genaue Doku-Drift, die der Carveout-Mechanismus eigentlich
verhindern sollte.

**Schritt 1 — Audit-Slice als ID-Reihe einplanen.** Konvention: ein
Slice `SL-CO-AUDIT-<welle>` pro Welle-Closure, *bevor* die Welle nach
`done/` wandert. ID-Schema-Ergänzung in
[`konventionen.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/konventionen.md): Audit-Slices haben
ein Präfix, das sie vom regulären Implementierungs-Slice unterscheidet
— sie liefern *keinen Code*, nur Doku-Updates.

```markdown
# SL-CO-AUDIT-welle-2: Carveout-Audit vor Welle-2-Closure

**DoD:**
- Jeder aktive Carveout in `docs/plan/carveouts/` hat ein aktuelles
  `Letzte Prüfung:`-Datum (≤ heute).
- Jeder Carveout, dessen Trigger eingetreten ist, ist nach `done/`
  verschoben.
- Jeder Carveout, der seit > 2 Wellen "aktiv" ist, wurde *explizit*
  entweder als weiter-gültig bestätigt oder in eine ADR überführt.
- Audit-Bericht als Closure-Notiz in `done/welle-2-results.md`.
```

Vier DoD-Punkte: drei Status-Aktionen plus ein Belegartefakt. Mehr
braucht es nicht — Audit ist *Disziplin*, nicht *Forschung*.

**Schritt 2 — Audit-Bericht-Schablone festlegen.** Damit der Audit
nicht jedes Mal neu erfunden wird, eine Tabelle als
Closure-Notiz-Block:

```markdown
## Carveout-Audit — Welle 2 (2026-06-12)

| Carveout | Status vorher | Status nachher | Aktion |
|---|---|---|---|
| CO-001 (Index-Coverage) | aktiv, Trigger Welle 2 | aufgelöst | git mv nach `done/`; coverage.config-Ausnahme entfernt |
| CO-004 (Compose-Devmode) | aktiv, Trigger "Compose v2.20" | permanent | überführt in ADR-0014 (Devmode als bewusste Architektur) |
| CO-005 (Lock-File-Pin) | aktiv, Letzte Prüfung 2025-12 | aktiv, geprüft | Datum 2026-06-12 nachgetragen, Folge-Slice slice-018 angelegt |
```

Drei Status-Übergänge sind möglich: *aufgelöst* (Trigger eingetreten),
*permanent* (Trigger wird nie eintreten — in ADR überführen),
*weiterhin aktiv* (Trigger weiterhin sinnvoll — Datum nachtragen).

**Schritt 3 — Wer führt den Audit aus?** Rollen-Bezug (Modul 8):
*Planner* identifiziert die fälligen Carveouts vor Welle-Closure,
*Architect* entscheidet bei "permanent" über die ADR-Überführung,
*Implementer* führt die `git mv`-Operationen und Config-Updates aus.
Der Audit-Slice landet damit über drei Rollen verteilt — was *nicht*
ein Defekt ist, sondern Absicht: Carveout-Aufräumen ohne Architect-Blick
verlängert das Lügen.

**Schritt 4 — Audit-Lauf-Gate optional einbauen.** Maschinell prüfbar
ist mindestens die "Letzte Prüfung"-Frische:

```makefile
verify-carveout-freshness:  ## Modul 7 — Audit-Pflicht pro Welle
	@python tools/check_carveout_freshness.py --max-age-days 90

verify: verify-carveout-freshness
```

Ein Carveout, dessen letzte Prüfung > 90 Tage zurückliegt, ist ein
HIGH-Warnsignal — egal, ob der Trigger nominell noch gilt. *Beobachtung
schlägt Behauptung*: ein nicht geprüfter Carveout ist ein nicht
existierender Audit.

**Schritt 5 — Audit-Slice als Schablone festschreiben.** Damit der
Slice in jeder Welle wiederverwendet wird, kommt eine Vorlage unter
`docs/plan/planning/templates/carveout-audit.md`. Der Planner kopiert
sie für jede neue Welle und passt nur das Datum, die betroffenen
Carveouts und den Welle-Bezug an. Ohne Vorlage wird der Audit-Slice
beim dritten Mal vergessen — und die Drift kehrt zurück.

Der Carveout-Mechanismus
hält nur, wenn er von einem *zweiten* Mechanismus auditiert wird;
sonst ist er eine schöne Konvention, die niemand prüft.

### Regeln gegen typische Fehlannahmen (Modul 7)

- **Gegen "Carveout = Workaround":** Carveout = *dokumentierter* Workaround mit Trigger. Ohne Trigger ist es eine versteckte Annahme.
- **Gegen "Carveouts gehören ins Issue-Tracker":** Sie gehören ins Repo, neben Spec und ADRs. Tracker können vergessen werden, Repo-Files kommen mit beim Klonen.
- **Gegen "Wenn der Trigger eintritt, lösen wir den Carveout auf":** Realität: er bleibt liegen. Deshalb braucht jeder temporäre Carveout einen *Folge-Slice mit ID*, der das Auflösen plant. Slice schlägt Memo.
- **Gegen "Jede entdeckte Diskrepanz ist ein eigener Carveout":** Carveouts sind für **punktuelle** Ausnahmen mit Folge-Slice. Eine Diskrepanz-**Häufung** in einer Sub-Area (Symptom: mehrere Carveouts mit demselben Geltungsbereich, oder die Diskrepanz folgt aus generellem *"Code existiert vor Doku"*-Muster) gehört nicht in eine Carveout-Kaskade, sondern in eine **BF-Sub-Area-Markierung mit Graduation-Plan** (siehe [Modul 2 §Kernidee](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/01-spec-und-architektur/modul-02-harness-bootstrap.md#kernidee)). Maßgeblich ist das **Symptom-Muster** (gemeinsamer Geltungsbereich), nicht die Carveout-Zahl; die Wahl, welches Werkzeug bei welchem Symptom greift, leistet [§Worked Example A Schritt 6](#worked-example-a-einen-carveout-dokumentieren).
- **Gegen "Wenn Diskrepanz-Häufung BF-Markierung verlangt, ist auch jede einzelne Diskrepanz eine BF-Markierung wert":** BF-Markierung lohnt sich erst beim **Cluster im selben Geltungsbereich** oder beim systemischen *"Code existiert vor Doku"*-Muster — eine einzelne, gut abgrenzbare Diskrepanz mit klarem Folge-Slice ist und bleibt ein Carveout. Das Frage-Schema in [§Worked Example A Schritt 6](#worked-example-a-einen-carveout-dokumentieren) trennt diese Fälle: Frage 1 leitet einzelne Diskrepanzen explizit auf den Carveout-/ADR-Pfad, nicht auf BF-Markierung.

