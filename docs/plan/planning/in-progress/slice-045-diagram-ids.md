# Slice slice-045: Modul `diagrams` — Kennungs-Konsistenz in Mermaid-Fences

**Status:** in-progress (welle-34-diagram-ids).

**Welle:** welle-34-diagram-ids (Trigger: belief-agent-Architektur — `ARC-NN`/
`LH-*`-Kennungen leben in `mermaid`-Diagrammen und entgehen heute allen Modulen,
weil d-check Fences opak behandelt).

**Bezug:** Neue Anforderung
[`DC-FA-DIAG-001`](../../../../spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in)
(opt-in Modul `diagrams`). Verwandt mit
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
(teilt die Kennungs-Muster, **nicht** die Link-Policy) und
[`DC-FA-LINK-001`](../../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
(dessen Vorverarbeitung Fences opak macht — genau die Invariante, die dieser
Slice gezielt durchbricht). **Eine ADR** zur Preprocessing-Fence-Ausnahme entsteht
mit der Implementierung (Nummer beim Schreiben vergeben).

**Autor:** pt9912. **Datum:** 2026-06-23.

---

## 1. Ziel

Ein opt-in Modul `diagrams`, das **gezielt benannte Diagramm-Fences** (Default
`mermaid`) öffnet, daraus **Kennungen per Regex** zieht und prüft, ob jede in
ihrer `defined-in`-Quelle **definiert** ist. Befund `diagram-id-undefined` fängt
Tipp-/Drift-Fehler (`ARC-99`, gelöschtes `LH-*`) in Diagrammen, die heute
unsichtbar sind. **Kein** Mermaid-Parser, **keine** Syntax-Validierung — reine
Token-Extraktion, read-only/netzlos.

## 2. Entscheidungen

- **Existenz statt Link-Policy (der Kern-Unterschied zu `ids`).** In Mermaid
  *kann* keine Markdown-Verlinkung stehen — eine Diagramm-Kennung wäre unter
  `link-policy: always` ein unbehebbares `id-unlinked`. Daher prüft `diagrams`
  **Existenz** der Kennung in `defined-in`, nicht deren Verlinkung. Eigenes Modul
  statt Überladung von
  [`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids).
- **Scoped Fence-Öffnung (ADR-würdig).** Heute ignorieren **alle** Module
  Fence-Inhalt (Vorverarbeitung des Moduls
  [`links`](../../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)).
  `diagrams` öffnet **nur** die in `diagrams.fences` gelisteten Sprachen — eine
  bewusste Ausnahme einer Kern-Invariante, daher **eine ADR** (warum nur
  Diagramm-Fences, warum Existenz statt Link-Policy). Andere Fences (`bash`/`go`)
  bleiben opak (False-Positive-Schutz: ID-ähnliche Tokens in Code-Fences sind
  keine Referenzen).
- **Kein Grammatik-Parsing.** Der `mermaid`-Block ist Rohtext; die `ids`-Regex
  läuft darüber. Deterministisch ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)),
  read-only/netzlos ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)),
  keine Toolchain — Syntax-Validierung bleibt bewusst out of scope (Diagramm-Linter).
- **Config (`model.DiagramsConfig`):** `fences []string` (Default `[mermaid]`) +
  `patterns` (je `regex` + `defined-in`, Gestalt wie `ids.patterns`). Decode/
  Validierung in `configyaml` analog `codepaths`. Default aus → byte-identisch.
- **Vorlagen-Parität:** `--print-config`-Gerüst und ggf. `--suggest-config`
  führen den `diagrams`-Block mit (sonst Schema-/Vorlagen-Drift wie bei früheren
  opt-in-Modulen).

## 3. Definition of Done

- [x] **Lastenheft-CR**
  [`DC-FA-DIAG-001`](../../../../spec/lastenheft.md#dc-fa-diag-001--kennungs-konsistenz-in-diagramm-fences-modul-diagrams-opt-in)
  (Bereich `DIAG`, Modul-Liste in [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl), Version 0.25.0, §7) — dieser CR.
- [ ] **ADR** (Preprocessing-Fence-Ausnahme): Begründung scoped Fence-Öffnung +
  Existenz-statt-Link-Policy; schärft die `spec/spezifikation.md`-Vorverarbeitung;
  ADR-Index ergänzt.
- [ ] `spec/spezifikation.md`: die `.a`-Algorithmus-Sektion (Fence-Öffnungsliste, Token-
  Extraktion, Existenz-Auflösung, Grund-Code).
- [ ] Glossar-Term `Diagramm-Fence`/`DIAG` im Lastenheft (Konvention wie `CODE`/`SPAN`).
- [ ] `model.DiagramsConfig` + `configyaml` Decode/Validierung; Vorverarbeitung
  öffnet gelistete Fences; `CheckDiagrams` (Token-Extraktion + Existenz-Auflösung,
  Grund-Code `diagram-id-undefined`).
- [ ] **Definition-Auflösung:** „Kennung kommt in `defined-in` vor" muss auch
  **tabellen-definierte** IDs (z. B. `ARC-NN`) abdecken, nicht nur Headings
  (heutige Auflösung ist heading-zentriert) — siehe Risiko.
- [ ] Tests: `configyaml`-Decode; Rule-Test (Happy 0 / Boundary 1 / **ohne**
  Config byte-identisch / `bash`-Fence unberührt).
- [ ] `--print-config`/`--suggest-config`: `diagrams`-Block sichtbar (Parität).
- [ ] [`CHANGELOG.md`](../../../../CHANGELOG.md).
- [ ] `make gates` grün (Coverage-Schwelle gehalten); unabhängiges Review R1;
  Closure (Slice → `done/` mit Roadmap-Flip,
  [`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)).

## 4. Risiken / offene Punkte

- **Definition-Auflösung (Haupt-Aufwand; Spec-Gap geschlossen, R1-F-1).** Der
  Lastenheft-CR definiert „definiert" jetzt explizit (Nicht-Fence-Token, Heading
  **oder** Tabelle). Offen bleibt der **Implementierungs**-Aufwand: d-checks heutige
  Definition-Auflösung (für `ids`-Link-Ziele) ist heading-/anker-zentriert und muss
  um die Token-Vorkommen-Prüfung erweitert werden, ohne einen Heading-Zwang
  einzuführen.
- **`spans`-Interaktion (R1-F-4).** Das aktive Modul `spans` meldet ungeschlossene
  Code-Spans. Geöffnete Diagramm-Fences dürfen dort keine False-Positives erzeugen
  — die ADR/`.a` legt fest, dass `diagrams` den Fence-**Rohtext** liest, ohne den
  Fence-Zustand für andere Module zu ändern.
- **False Positives.** Diagramm-Syntax-Tokens, die zufällig ein Kennungs-Muster
  treffen. Mitigation: spezifische Regex (`defined-in` je Muster); nur gelistete
  Fence-Sprachen; Gegentest mit `bash`-Fence.
- **Kern-Invariante.** „Fences sind opak" ist stark gelebt (selbst `hostpaths`
  exemptiert Fences, [`DC-FA-HOST-001`](../../../../spec/lastenheft.md#dc-fa-host-001--host-lokale-absolute-pfade-modul-hostpaths-opt-in)).
  Die Ausnahme ist scoped + opt-in + ADR-begründet — kein Default-Verhaltens-Delta.
- **Dogfooding-Grenze.** d-checks eigene Diagramme tragen **keine** Kennungen — der
  Payoff liegt bei Konsumenten (belief-agent). Tests laufen daher über Fixtures
  (wie `codepaths`/`spans`/`hostpaths`), nicht über die eigene Doku.
- **Determinismus.** Ohne `diagrams`-Block byte-identisch
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)) — Gegentest Pflicht.

## 5. Trigger

belief-agent-Architektur (2026-06-23, Nutzer-Hinweis): zwei nicht-triviale
`mermaid`-Diagramme — `ARC-01…ARC-09` im Flowchart, `LH-FA-POL-004` im
Sequenzdiagramm; Drift/Typo darin ist heute unsichtbar. Abgegrenzt von Mermaid-*Syntax*-Validierung (gehört
in einen Diagramm-Linter, nicht in d-checks distroless-Domäne).

## 6. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Produkt-Code `rules`/`model`/`configyaml` + Spec;
Greenfield-Default „Doc führt, Code folgt"). Keine BF-Sub-Area berührt.

## 7. Closure-Notiz (nach `done/`)

_(Wird beim Move nach `done/` gefüllt — Belege: `make gates`-Ausgabe, ADR-Nummer,
unabhängiges Review R1, ggf. Release-Digest.)_
