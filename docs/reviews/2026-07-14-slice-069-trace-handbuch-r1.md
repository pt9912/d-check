# Review-Report: slice-069 Trace-Handbuch-Parsergrenzen R1

**Review-Art:** Doku-/Code-Review des uncommitteten Working-Tree-Diffs.
**Review-Unabhängigkeit:** Self-Review durch den Implementierungs-Agenten; kein
personell unabhängiges Review.
**Gegenstand:** `slice-069` — Handbuch-Präzisierung, Handbuch-E2E-Verankerung,
Roadmap und Changelog.
**Reviewer-Skill:** `.harness/skills/reviewer.md` v1.2.0.
**Modell:** Codex auf GPT-5.
**Datum:** 2026-07-14.

**Eingangs-Kontext:** Working-Tree-Diff; Slice-Plan ohne Bewertung der
DoD-Abhakung; [`DC-FA-CLI-009`](../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix),
[`DC-FA-CLI-011`](../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code),
[`DC-FA-COV-001`](../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in),
[`DC-FA-MOD-001`](../../spec/lastenheft.md#dc-fa-mod-001--modalitäts-klassifikation-der-anforderungen-tracerequirementsmodality-opt-in);
Spezifikation §§`DC-FA-CLI-009.a`, `DC-FA-COV-001.a`,
`DC-FA-MOD-001.a`, `DC-FA-CLI-011.a`; ADR-0034 bis ADR-0036; `AGENTS.md`
§3; Trace-Parser und Handbuch-Beispiel-Harness.

---

## Findings

### R1-MEDIUM-1 — Einleitende Waisen-Definition widerspricht der präzisierten Semantik

- **kategorie:** MEDIUM
- **quelle:** `DC-FA-CLI-011`, `DC-FA-COV-001`
- **pfad:** `docs/user/benutzerhandbuch.md:575`
- **befund:** Das Kapitelziel bezeichnet eine Waise weiterhin als Anforderung,
  die „**niemand** referenziert". Eine nur durch eine ADR referenzierte
  Anforderung ist jedoch referenziert und bleibt trotzdem Waise
  (`¬slice ∧ ¬coverage`); damit widerspricht der Einstieg der korrekten
  Definition in Zeile 611. Der neue Marker-Test prüft nur die Anwesenheit der
  korrekten späteren Aussage und bleibt trotz dieses Widerspruchs grün.
- **verifizierbar:** nein — die aktuellen Doku-/Test-Gates interpretieren keine
  widersprüchlichen Prosaaussagen; der Widerspruch ist nur durch Abgleich mit
  Spezifikation und `traceRow` feststellbar.

### R1-MEDIUM-2 — Allgemeine Gate-Regel ignoriert `modality.require-levels`

- **kategorie:** MEDIUM
- **quelle:** `DC-FA-MOD-001`, `DC-FA-CLI-011`
- **pfad:** `docs/user/benutzerhandbuch.md:629`
- **befund:** Die allgemeine Erklärung behauptet „mindestens eine Waise ⇒
  Exit 1". Bei aktivem `modality` endet `--require-complete` dagegen mit Exit 0,
  wenn alle Waisen außerhalb von `require-levels` liegen; eine sichtbare
  SOLLTE-/KANN-/`unknown`-Waise genügt also nicht. Ein Konsument kann die Zeile
  als strengere Gate-Garantie lesen, als v0.42.0 tatsächlich bietet.
- **verifizierbar:** ja — `make test` belegt mit
  `TestCLI068_Modality_KannAdvisory` den Exit 0 für eine KANN-Waise unter dem
  Default `[must]`, während die zitierte Handbuchzeile Exit 1 behauptet.

### R1-LOW-1 — Maschinenlesbare RTM-Feldliste lässt konditionale Felder aus

- **kategorie:** LOW
- **quelle:** `DC-FA-COV-001`, `DC-FA-MOD-001`
- **pfad:** `docs/user/benutzerhandbuch.md:625`
- **befund:** Die Feldliste für `--trace --json`/`--yaml` nennt je Requirement
  nur `id`/`title`/`adrs`/`slices`/`orphan`. Bei aktiver Coverage oder
  Modalität erscheinen zusätzlich `coverage` beziehungsweise `modality`; erst
  die späteren Unterabschnitte erwähnen diese Erweiterungen. Ein Nutzer, der
  an der zentralen Ausgabeerklärung ein Schema ableitet, erhält damit eine
  unvollständige Feldbeschreibung.
- **verifizierbar:** ja — `make test` enthält JSON-Akzeptanzfälle für beide
  konditionalen Felder; die Handbuch-E2E-Strukturprüfung deckt diese
  Trace-Varianten derzeit nicht ab.

## Negativbefunde

- **Definitionssyntax geprüft, ohne weiteren Befund:** ATX-only,
  erstes Ganz-Token, Fence-/Tabellen-/Listen-/Setext-Grenze stimmen mit
  `traceRequirements` und dem gemeinsamen Heading-Scanner überein.
- **Nullmengenwarnung geprüft, ohne weiteren Befund:** fehlende Quelle,
  unpassendes Format und null ID-Treffer liefern laut Spezifikation und Code
  eine leere RTM; `--require-complete` endet bei null Waisen mit Exit 0.
- **Modalitäts-Body geprüft, ohne weiteren Befund:** `modality: {}` aktiviert
  Defaults; Body-Span, Normalisierung und Gleich-/Höherrang-Grenze stimmen mit
  ADR-0036 und `requirementModality` überein.
- **Referenzscan geprüft, ohne weiteren Befund:** rekursiver Markdown-Scan,
  Basisnamen-Regex, Capture-Gruppe 1 plus Präfix, Ganzdatei-ID-Suche,
  Deduplizierung und Sortierung stimmen mit Spezifikation und `traceRefs`
  überein.
- **Brownfield-Migration geprüft, ohne Befund:** beide Wege behaupten keine
  native Tabellenunterstützung und benennen die Drift-Grenze einer Projektion.
- **Testverankerung geprüft, ohne weiteren Befund:** die E2E-Fixture reproduziert
  die tabellarische Nullmenge samt Exit 0; Whitespace-normalisierte Marker sind
  deterministisch und der vorhandenen Handbuch-Harness-Struktur konform.
- **Planning/Changelog geprüft, ohne Befund:** aktive Welle, Slice-Pfad,
  No-CR-/No-Release-Abgrenzung und Changelog-Eintrag sind konsistent.
- **Hard Rules geprüft, ohne Befund:** keine Accepted-ADR geändert, keine
  Gate-Lockerung, keine Inline-Suppression, keine Architektur-/Spec-Abwärtsreferenz.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
| ---: | ---: | ---: | ---: |
| 0 | 2 | 1 | 0 |

## Verdikt

**NACHBESSERN.** Die zwei MEDIUM-Befunde betreffen die zentrale Waisen- und
Exit-Code-Erklärung eines Vollständigkeits-Gates und sollten vor Merge und
Slice-Closure geklärt werden. R1-LOW-1 blockiert allein nicht. Wegen des
Self-Review-Status bleibt ein personell unabhängiges Review zusätzlich
empfehlenswert.
