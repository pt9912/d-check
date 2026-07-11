# Review — slice-067 (`trace.coverage`)

**Datum:** 2026-07-11 · **Slice:** slice-067 · **Bezug:** `DC-FA-COV-001`, ADR-0035,
Mods `DC-FA-CLI-009`/`DC-FA-CLI-011` · **Reviewer:** unabhängiger Subagent (2 Läufe)

Zwei unabhängige Läufe: **R1 Doc-first** (vor der Implementierung) und **R2 Impl**.

## R1 — Doc-first · Verdikt NACHBESSERN → behoben

Prüfte das Fundament (Lastenheft-CR, Spezifikation, ADR-0035, Slice) gegen die
9-Punkte-Nutzer-Spec, den Anforderungs-Anlege-Prozess und die Span-Semantik.

- **MEDIUM (behoben):** Die wiederverwendete Span-Mechanik (`plainHeadingText`)
  vergleicht den **vollen Heading-Text exakt**; die Beispiele in ADR/Spezifikation/
  Code-Kommentar nutzten die Kurzform `exclude-sections: [27.1.1]`, die die reale
  Überschrift `### 27.1.1 Anforderungen ohne Design-Artefakt` **nicht** matcht
  (§27.1.1 nicht ausgenommen ⇒ falsche Kredite; Whitelist-Tippfehler ⇒ still alle
  Zeilen leer ⇒ alle Anforderungen falsche Waisen). Behoben: alle Beispiele auf
  vollen Heading-Text, explizite Aussage „Sektionsname = voller Heading-Klartext",
  **plus fail-closed-Guard** (`checkSectionNames`: Sektionsname ohne Heading-Treffer
  ⇒ Exit 2 — macht die stille Falle laut).
- **INFO (eingearbeitet):** Spalten-Position (`Coverage` vor `Status`),
  `files`-Repo-Escape-Validierung, „`trace.coverage` führt kein eigenes Regex"
  klargestellt; die Whitelist-ohne-Treffer-fail-open-Falle über denselben Guard
  geschlossen.
- **Negativbefunde sauber:** treue Abbildung aller 9 Punkte; Range-Kanten
  (`GG-DNP3-001`, Breiten-Erhaltung, `AAA>BBB`); Byte-Identität/Konditional-Spalte;
  DC-FA-CLI-011-Vertrag; Harness-Regeln (Bereich `COV`, AKs, ADR-Geschichte-only).

## R2 — Impl · Verdikt ACCEPT-WITH-NITS → behoben

Prüfte den Code gegen die (nachgebesserte) Spec; kein Korrektheitsfehler, kein
falscher Exit-Code, keine Panik/Off-by-one.

- **MEDIUM (behoben):** die **positive `sections`-Whitelist** (Include-Zweig von
  `SelectSections`) war ungetestet — der einzige `sections`-Test brach vorher im
  Guard ab. Ergänzt: `TestCLI067_Coverage_IncludeSection` (§27.2 in Whitelist ⇒
  gedeckt, §27.1 außerhalb ⇒ Waise).
- **LOW (behoben):** `ranges: false` end-to-end ungetestet ⇒
  `TestCLI067_Coverage_RangesFalse` (keine `..`-Expansion, nur exakt).
- **LOW (behoben):** `--require-complete`-Exit-1-Meldung sagte „ohne Slice" ⇒ bei
  aktiver Coverage jetzt „ohne Slice und ohne Coverage".
- **INFO (bewusst nicht):** die Range-Schleife läuft `end-start`-mal unabhängig vom
  `id-pattern`; für die 3-stellige Konvention ≤1000, Atoi-Overflow entschärft die
  Extremform. Trusted-Repo-Eingabe, kein Sicherheitsproblem — kein Cap eingebaut.
- **Negativbefunde sauber (per Fallanalyse mutations-hart):** Byte-Identität
  (`omitempty` json+yaml, `CoverageActive json:"-"`, 5-Spalten-Fallback,
  `Slices==nil`-Normalisierung verhaltensgleich); Range-Parser (DNP3/Breite/
  AAA>BBB, Voll-Match nicht durch Teilstring unterlaufbar); Sektions-Filter
  kein Off-by-one (beide `proseLines`/`i+1`), Guard trifft über alle files;
  Fail-closed vollständig verdrahtet bis Exit 2; Waise-Neubestimmung; Spec↔Code-
  Parität; Arch/Lint/Reuse (`sortedSets` verhaltensgleich zum Inline-Dedup).
