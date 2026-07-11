# ADR-0035 — Kuratierte Coverage-Quellen der RTM (`trace.coverage`, range-aware)

**Status:** Accepted
**Datum:** 2026-07-11
**Autor:** pt9912
**Bezug:** [`DC-FA-COV-001`](../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
(die neue Anforderung; Spezifikation
[§`DC-FA-COV-001.a`](../../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage));
[`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
(RTM-Struktur — Coverage-Spalte) und
[`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
(Waisen-Definition), beide mit-modifiziert;
[ADR-0034](0034-trace-konfigurierbare-quellen.md) (die konfigurierbaren
RTM-Quellen, auf denen dies aufbaut);
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).

## Kontext

Nach [ADR-0034](0034-trace-konfigurierbare-quellen.md) liest die RTM Abdeckung
aus **Referenz-Scans** von ADR-/Slice-Dateien. Für ein Repo, dessen Abdeckung in
einer **kuratierten Matrix** liegt, misfeuert das: der Konsument **grid-gym**
lieferte den Großteil seiner Anforderungen in einer Wellen-Ära (`M*.md`) und
weist die Deckung in einer ausgelagerten Traceability-Matrix + ADRs nach — nichts davon sind
`NNN-…md`-Slices. Von 171 „Waisen" waren **≥122 anderswo belegt** (83 mit
ADR-Link, 80 exakt in `traceability.md`, 73 im Wellen-Archiv). Zwei Ursachen:

1. **Bereichs-Notation.** `traceability.md` schreibt `GG-QA-001..006`; der
   ID-Scanner erkennt nur `GG-QA-001`, nicht `002…006`.
2. **Keine Referenzklasse für kuratierte Quellen.** Man könnte `slices` auf
   `docs/plan` überladen — das kontaminiert aber mit ADR-Dateien (Owner `0013`)
   und vergibt irreführende Labels (`Slice traceability`).

Der `slices.file-pattern`-Trick von ADR-0034 löst nur die Wellen-Docs (`M\d+`
mit-matchen, 185→113 Waisen gemessen); die kuratierte Matrix bleibt außen vor.

## Entscheidung

### 1. Eine dritte, opt-in Referenzklasse `trace.coverage` (Liste)

`trace` erhält neben `adrs`/`slices` eine **Liste** `coverage` benannter Quellen.
Je Quelle: `files` (explizite Pfad-Liste), `label` (feste Owner-Kennung),
`ranges` (bool, Default true), `sections`/`exclude-sections`.

- **`files`, nicht `dir`+`file-pattern`.** Kuratierte Quellen werden **einzeln
  benannt**, nicht über ein Verzeichnis-Muster abgeleitet — das verhindert die
  ADR-/Slice-Kontamination strukturell (man kann `traceability.md` nicht
  versehentlich als „Slice" zählen).
- **Eigene `label`-Owner + eigene Coverage-Spalte.** Coverage-Labels stehen in
  einer **separaten** RTM-Spalte, nicht in `Slices` — Coverage ist eine **andere
  Dimension** als Slice-Implementierung, keine Vermischung.

### 2. Range-Parser (die einzige neue Kernlogik)

Bei `ranges: true` werden nach jeder `requirements.id-pattern`-Fundstelle
unmittelbar folgende Range-/Enum-Suffixe expandiert:

- **`<FAM>-AAA..BBB`** (Kurzform): Familie = die Kennung ohne die
  Trailing-`-<Ziffern>`; expandiert **breiten-erhaltend** zu `<FAM>-{i}` für
  `i ∈ [AAA, BBB]` (`001..007` ⇒ `001`,…,`007`). Nur die Kurzform — ein
  fam-qualifiziertes Ende (`..<FAM>-BBB`) ist Out-of-Scope (kommt in den Daten
  nicht vor).
- **`<FAM>-AAA/BBB/CCC`** (Aufzählung): `<FAM>-AAA`, `<FAM>-BBB`, `<FAM>-CCC`.
- Jede expandierte ID wird gegen `requirements.id-pattern` geprüft; Nicht-Treffer
  **still verworfen** (Sicherheitsnetz gegen Über-Expansion).
- **Fail-closed:** `AAA>BBB` oder abweichende Ziffern-Breite ⇒ Exit 2.

Der Parser ist eine **isolierte, unit-getestete Funktion**, parametrisiert über
das `id-pattern` — wiederverwendbar für eine spätere `commits`/`ids`-Range.

### 3. Abschnitts-Scoping über die **bestehende** Span-Semantik

`exclude-sections` (Blacklist) und `sections` (Whitelist) nutzen **dieselbe**
Heading-Span-Mechanik wie `matrix.exclude-sections` (Abschnitt = Überschrift bis
zur nächsten gleich-/höherrangigen). **Ein Sektionsname ist der volle
Heading-Klartext** (exakter Vergleich wie bei `matrix.exclude-sections`, z. B.
`"27.1 Anforderung zu Design"`, **nicht** die Kurzform `"27.1"`). Das ist
korrektheits-kritisch, **nicht** Kosmetik: `traceability.md` hat die Überschrift
`### 27.1.1 Anforderungen ohne Design-Artefakt` — ein Ganz-Datei-Scan würde
**gerade die nicht gedeckten** IDs kreditieren.
`sections: ["27.1 Anforderung zu Design", "27.2 Anforderung zu Implementierung", "27.3 Anforderung zu Test"]`
(Whitelist, H2-Span) **plus** `exclude-sections: ["27.1.1 Anforderungen ohne Design-Artefakt"]`
(Blacklist, H3-Span) ergibt „§27.1 ohne §27.1.1" — mit der Standard-Span-Semantik,
**ohne** einen neuen „bis-beliebige-Ebene"-Span. Reine Wiederverwendung
(`excludedRanges`), kein Neucode am Span. **Fail-closed:** ein konfigurierter
Sektionsname, der in den Quell-Dateien **keine** Überschrift trifft (Tippfehler /
Kurzform), ⇒ Exit 2 — sonst blankt die Whitelist still die ganze Datei
(alle Anforderungen erneut falsche Waisen) bzw. greift die Blacklist nicht.

### 4. Coverage zählt zur Vollständigkeit (`DC-FA-CLI-011`-Änderung)

Eine Anforderung ist **Waise** nur, wenn sie **weder** Slice- **noch**
Coverage-Referenz trägt. `--require-complete` bricht (Exit 1) nur bei so
definierten Waisen. Das ist eine **bewusste Vertrags-Änderung** an
[`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
(dessen Out-of-Scope bisher „andere RTM-Eigenschaften als
Waisenfreiheit" ausschloss) — Coverage ist keine *neue* Gate-Eigenschaft, sondern
eine **breitere** Waisen-Definition. Eine bloße **ADR**-Referenz deckt weiterhin
**nicht** ab (ADR = Design, nicht Umsetzung/Nachweis).

### 5. Konditionale Spalte (Byte-Identität)

Die Coverage-Spalte und das `coverage`-json/yaml-Feld erscheinen **nur**, wenn
≥1 `trace.coverage`-Quelle konfiguriert ist. Ohne Coverage ist die RTM
**byte-identisch** ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)) —
die bestehenden Byte-Identitäts-Tests und die Handbuch-E2E-Beispiele
(5-Spalten-Tabelle) bleiben grün. `coverage` trägt `omitempty` (leer ⇒ entfällt).

### Verglichene Alternativen

| Option | Pro | Contra |
| --- | --- | --- |
| **(A) `slices` auf `docs/plan` überladen** | keine neue Klasse | kontaminiert mit ADR-Dateien (Owner `0013`); irreführendes Label `Slice traceability`; keine Range-Expansion; vermischt Design/Umsetzung/Nachweis |
| **(B) Waise auch bei ADR-Referenz „gedeckt"** | trivial | konflatiert **Design** (ADR) mit **Umsetzung/Nachweis**; gegen die [`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)-Intention; false-grün |
| **(C, gewählt) separate opt-in `coverage`-Klasse, eigene Spalte, range-aware, sektions-gescopt** | ehrlich getrennte Dimension; range-aware; keine Kontamination; byte-identisch im Default | eine neue Config-Fläche + Range-Parser (isoliert, testbar) |

**Fitness-Funktion:**

- grid-gym mit `coverage: [{files: [traceability.md], ranges: true}]` (+ `sections`/
  `exclude-sections` gegen §27.1.1): Waisen sinken **deutlich** gegen ~0; `make
  doc-check` bleibt grün (Beleg im umsetzenden Slice, an Realdaten gemessen).
- `GG-QA-001..006` in der Quelle deckt **alle sechs** ab (Range-Beweis).
- ADR-Dateien werden **nicht** als Coverage/Slice gezählt.
- Ungültige Range (`GG-RT-009..003`) / fehlende Datei ⇒ Exit 2.
- **Ohne** `trace.coverage`: RTM byte-identisch (kein Spalten-/Feld-Diff).

## Konsequenzen

- **Neue Anforderung + zwei Modifikationen (Lastenheft-CR, Versions-Bump +
  Historie).** [`DC-FA-COV-001`](../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
  neu (Bereich `COV`); [`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
  (Coverage-Spalte + json/yaml-Feld) und [`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)
  (Waisen-Definition) mit-modifiziert.
  Nutzersichtbare Config-Fläche + geänderter Output bei aktiver Coverage →
  **Release**.
- **Modell** (`config.go`): `TraceCoverage` (Files/Label/Ranges/Sections/
  ExcludeSections) + `TraceConfig.Coverage []TraceCoverage`; `TraceRow.Coverage
  []string` (`omitempty`); `TraceMatrix` trägt, ob Coverage aktiv ist (Spalten-
  Entscheidung des Reporters).
- **Config-Decode** (`configyaml.go`): `rawTrace.Coverage`; `label`/`files`
  nicht-leer, `ranges`-Default-true (Pointer), Regex/Range-Validierung fail-closed.
- **RTM** (`trace.go`): Coverage-Scan je Quelle (abschnitts-gefilterter Text →
  exakte + range-expandierte IDs → Label), Waise = ¬slice ∧ ¬coverage; neuer
  Range-Parser + Wiederverwendung von `rules`-Section-Span.
- **Reporter** (`report.go`): konditionale Coverage-Spalte (Markdown), `coverage`
  in json/yaml (omitempty).
- **`--print-config`**: kommentierter `coverage`-Block. `--print-mk` unverändert.
- **Determinismus/Read-only** unberührt (nur Markdown-Lesen, kein neuer Eingabe-
  Scope über den gemounteten Baum hinaus). **Reversibel** im Verhalten (Default
  byte-identisch), aber Vertrags-Änderung → Lastenheft-CR.

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-07-11 | Entwurf (slice-067, welle-56; Konsumenten-Analyse grid-gym: 171 „Waisen" zu ≥122 anderswo belegt — ADR/`traceability.md`/Wellen —, weil die slice-zentrische RTM die kuratierte Deckungs-Matrix mit Bereichs-Notation nicht erkennt). Dritte opt-in Referenzklasse `trace.coverage` (Liste benannter `files` + `label` + `ranges` + `sections`/`exclude-sections`), range-aware (`<FAM>-AAA..BBB`/`/`-Enum, breiten-erhaltend, gegen `id-pattern` validiert, fail-closed), Abschnitts-Scoping über die bestehende `matrix`-Span-Semantik (gegen die §27.1.1-„ohne Design-Artefakt"-Falle); Coverage zählt zur Waisen-Definition ([`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)), eigene RTM-Spalte konditional (byte-identisch ohne Quelle). `files`-statt-`dir` gegen ADR-Kontamination; ADR-Referenz deckt weiter nicht. Lastenheft-CR (v0.41.0), Release geplant. Status Proposed. |
| 2026-07-11 | Angenommen mit der slice-067-Closure: `model.TraceCoverage` + `configyaml.applyTraceCoverage` (fail-closed) + `app.coverageRefs`/`expandRange`/`checkSectionNames` + `rules.SelectSections`/`HeadingTexts` (Section-Span-Wiederverwendung); Reporter-konditionale Coverage-Spalte. End-to-End gegen grid-gyms echte `traceability.md` verifiziert: Waisen **113 → 10** (`GG-QA-001..006` via Range gedeckt). Doc-first-Review R1 NACHBESSERN (Sektionsname = voller Heading-Text + fail-closed-Guard) + Impl-Review R2 ACCEPT-WITH-NITS (Whitelist-/`ranges:false`-Test, Meldungs-Text) eingearbeitet. `make gates` grün, Release v0.41.0. Status Accepted. |
