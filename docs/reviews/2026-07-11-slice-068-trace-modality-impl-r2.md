# Impl-Review slice-068 — `trace.requirements.modality` (DC-FA-MOD-001)

**VERDICT: ACCEPT WITH NITS**

Reviewer: unabhängiger Impl-Reviewer (R2). Datum: 2026-07-11.
Grundlage: uncommitteter Working-Tree-Diff vs. HEAD (`b811a12`), Lastenheft
`DC-FA-MOD-001`/`DC-FA-CLI-009`/`DC-FA-CLI-011`, Spezifikation
`DC-FA-MOD-001.a` (Schritte 0–5), ADR-0036, Slice-068-DoD.

Verifikation: `make test` (Docker, `go test ./...`) **grün, Exit 0** — alle
Pakete `ok`, inkl. `internal/adapter/driving/cli` und `internal/hexagon/core/app`
(TestCLI068_* und TestModality* laufen und bestehen). Die Byte-Identitäts-Tests
der Nachbar-Slices (066/067) bleiben grün.

Zusammenfassung: die Implementierung ist spec-konform, der Klassifikator ist in
den kniffligen Fällen (Längster-Treffer, Wortgrenze, Normalisierung, früheste
Position) korrekt, die Default-aus-Byte-Identität ist strukturell sauber, und die
Config-Validierung ist config-zeitig fail-closed. Es bleiben ein
Test-Abdeckungs-Loch gegen ein **explizit gelistetes** AK (leerer Stufen-Name)
und mehrere kleinere Punkte. Keiner blockiert; alle sind vor Closure billig
nachziehbar.

---

## Dimensions-Urteil (Kurz)

1. **Spec-Konformanz** — erfüllt. Alle AK von DC-FA-MOD-001 sind im Code
   abgebildet; DC-FA-CLI-009 (Modality-Spalte, konditional, vor Status) und
   DC-FA-CLI-011 (Gating nur `require-levels`, stderr „gatende von") sind
   umgesetzt. Eine AK-Zeile ist implementiert, aber **ungetestet** (M1).
2. **Klassifikator-Korrektheit** — korrekt. Longest-first-Alternation +
   Go-Default-Leftmost-first ⇒ frühester/längster Treffer; `\b`-Wortgrenze
   trennt `musste` von `MUSS`; Normalisierung fängt `**MUSS** NICHT` / `MUSS\nNICHT`.
   Siehe Analyse unten.
3. **Byte-Identität (DC-QA-02)** — sauber. `omitempty` + `json:"-"`-Gating +
   5-Spalten-Fallback im Reporter erzeugen ohne `modality` byte-gleiche Ausgabe;
   von den 066/067-Tests belegt.
4. **Fail-closed (Exit 2)** — vollständig, config-zeitig (in `Decode`→`applyTrace`
   →`applyModality`). Alle fünf Fehlerklassen abgedeckt; kein Ordering-Bug.
5. **Gating-Semantik** — korrekt. `GatingOrphans` zählt Waisen mit Stufe ∈
   `require-levels`; ohne `modality` == `Orphans` (alle). KANN gatet nicht,
   MUSS gatet, Default `[must]` stimmt.
6. **Hexagonale Grenzen** — eingehalten. Kein git/IO im Kern außer dem
   Filesystem-Port; Matcher wird **einmal** je Lauf gebaut, nicht je Zeile.
7. **Test-Adäquanz** — überwiegend gut (Längster-Treffer, Wortgrenze, Emphasis,
   Negative-Config, Byte-Identität gepinnt). Zwei Lücken (M1, L2).
8. **Go-Qualität** — solide. Nil-Pfad des `*TraceModality` sicher; `MustCompile`
   auf QuoteMeta-Eingabe unerreichbar-panisch; keine Shadowings gesichtet.

---

## Findings

### MEDIUM

**M1 — AK „leerer Stufen-Name ⇒ Exit 2" ist implementiert, aber ungetestet
(mutations-unverriegelt).**
`internal/adapter/driven/configyaml/configyaml.go:439-441` (`validateModalityLevels`,
`if strings.TrimSpace(name) == "" { … }`) und
`internal/adapter/driving/cli/cli_acceptance_test.go:1958-1978`
(`TestCLI068_Modality_NegativeConfig`).

Die AK „Negative (Config)" von DC-FA-MOD-001 listet **vier** Exit-2-Fälle
explizit auf: leerer Stufen-Name, leeres Keyword, reserviertes `unknown`,
dasselbe Keyword in zwei Stufen. `TestCLI068_Modality_NegativeConfig` deckt drei
davon (+ ungültiges `require-levels`), aber **nicht** den leeren Stufen-Namen
(`levels: { "": [X] }`). Kein Unit-Test trifft `validateModalityLevels` direkt.
Damit ist der `name==""`-Guard nicht mutations-verriegelt — würde er entfernt
(oder zu `!= ""` mutiert), bliebe jeder Test grün. Das verstößt gegen die
projekteigene Bar (slice-057-R3/slice-060: jeder Fail-closed-Guard braucht einen
Mutations-Beleg, und die AK nennt den Fall wörtlich).

**Fix:** einen Fall `{"empty level name", "    modality:\n      levels:\n        '': [X]\n", "leerer Stufen-Name"}`
in die `cases`-Tabelle aufnehmen (analog zu „empty keyword"). Eine Zeile.

### LOW

**L1 — `modality:` (YAML-Null) wird still inaktiv, nur `modality: {}` aktiviert.**
`internal/adapter/driven/configyaml/configyaml.go:411-413` (`applyModality`, `raw==nil ⇒ nil`).

Die Aktiv-Erkennung hängt am dekodierten `*rawModality`-Zeiger. yaml.v3 setzt
den Zeiger bei einem **Null-Wert** (`modality:` ohne Inhalt) auf `nil` — nicht
unterscheidbar von „Schlüssel fehlt". Ein Nutzer, der `modality:` (leer,
einrückungslos) schreibt in der Erwartung „Präsenz ⇒ aktiv" (so formuliert Spec
Schritt 0), bekommt **still** keine Spalte, keinen Fehler. Spec/ADR/Tests nutzen
durchgängig `modality: {}`, daher kein Vertragsbruch — aber ein Footgun.

**Fix (optional):** entweder im Handbuch/`--print-config` klarstellen, dass `{}`
zum Aktivieren nötig ist (der Template-Kommentar tut das implizit), oder — falls
gewünscht — Präsenz über den Roh-Map-Schlüssel statt den Zeiger erkennen. Für
diesen Slice genügt die Doku-Klarstellung.

**L2 — Strikt-Modus `require-levels: [must, unknown]` (unknown *gatet*) ist
ungetestet.**
`internal/hexagon/core/app/trace.go:164`.

`TestCLI068_Modality_KannAdvisory` pinnt die „unknown/may gatet **nicht**"-
Richtung, aber nicht die von der AK „Unknown" ebenfalls geforderte Gegenrichtung
(„sie gatet **nur**, wenn `require-levels` `unknown` enthält"). Der ADR nennt den
`[must, unknown]`-Strikt-Modus als zentrales Gegenmittel zur Fail-open-Grenze; er
sollte einen Regressionstest haben (Waise ohne Modal-Verb + `require-levels:
[unknown]` ⇒ Exit 1).

### INFO

**I1 — ADR-0036 ist noch `Status: Proposed` und seine Default-Beispielliste hinkt
Spec/Impl hinterher.** `docs/plan/adr/0036-…md:3` und `:42-44`. Die ADR-Liste
nennt `must` ohne `SHALL NOT` und `should` ohne `SOLLTEN NICHT`; Spezifikation
(§DC-FA-MOD-001.a Schritt 1) und Impl (`config.go:DefaultModalityLevels`) führen
beide. Spec ist normativ und Impl == Spec — also **kein** Impl-Defekt, aber die
ADR-Aufzählung und ihr Status sollten bei Closure angeglichen bzw. auf `Accepted`
geflippt werden (Prozess, nicht Code).

**I2 — Quelldatei wird bei aktivem `modality` zweimal gelesen/geparst.**
`internal/hexagon/core/app/trace.go:128` (`traceRequirements`) und `:153`
(`requirementModality`) lesen beide `rt.source` und laufen `extractHeadingLines`.
Kein Korrektheits- oder Determinismusproblem (reiner Lese-Pfad, DC-QA-03 gewahrt),
nur eine bescheidene Doppel-Arbeit auf einer Datei. Bewusst schlichte Trennung —
akzeptabel; ein späteres Zusammenlegen (Titel + Body in einem Pass) wäre
Aufräumen, kein Muss.

**I3 — Klassifikation ignoriert Fenced-Code im Body und `_…_`-Emphasis.**
`normalizeBody` (`trace.go:517`) entfernt nur `*`/`` ` `` (Spec-konform), nicht
`_`; und der Body-Span (`HeadingSections`) schließt Fenced-Code nicht aus, sodass
ein Keyword in einem Code-Fence im Anforderungs-Body mitzählt. Beides deckt sich
mit dem Vertrag (Spec strippt nur `*`/Backtick; keine Fence-Ausnahme genannt) und
ist für reale Anforderungs-Prosa unkritisch — nur zur Kenntnis.

**I4 — `TestCLI068_Modality_KlassifikationUndGating` kann bei nicht gefundener ID
panicken statt sauber zu failen.**
`cli_acceptance_test.go:~1908` (`if i < 0 || doc.Requirements[i].Modality != lvl`).
Bei `i == -1` wertet `Fatalf` das Argument `doc.Requirements[i]` aus ⇒
Index-Panic. Nur ein Test-Robustheits-Nit (der Fall tritt nur bei bereits
kaputtem Verhalten ein); `doc.Requirements[i]` im Fatalf durch `id` ersetzen.

---

## Detail-Belege (verifiziert)

### Klassifikator (Dimension 2) — korrekt

- **Longest-first + Go-Semantik.** `resolveModality` sortiert die Keywords
  längster-zuerst (`trace.go:475-480`, Tiebreak byte-asc) und baut
  `(?i)\b(?:kw1|kw2|…)\b` (`:495`). Go-`regexp` ist per Default **leftmost-first**
  (kein `Longest()`): die früheste Start-Position gewinnt, und an dieser Position
  wird die **erstgelistete** matchende Alternative gewählt — durch die
  Längen-Sortierung also die **längste**. Damit ist die Spec-Semantik „frühester
  Treffer, an gleicher Position längste Phrase" exakt getroffen. Konfliktpaare
  (`MUSS`⊂`MUSS NICHT`, `MÜSSEN`⊂`MÜSSEN NICHT`, `SOLLTE(N)`⊂`SOLLTE(N) NICHT`,
  `MUST`⊂`MUST NOT`, `SHALL`⊂`SHALL NOT`, `SHOULD`⊂`SHOULD NOT`) haben je den
  längeren Eintrag mit größerer **Byte**-Länge (nur ASCII-Suffix ` NICHT`/` NOT`
  angehängt), sodass die Umlaut-Byte-Zählung die Sortierung nicht verdreht.
- **Wortgrenze.** `\bMUSS\b` matcht `musste` nicht (nach `MUSS` folgt `t`,
  Wort-Zeichen ⇒ kein `\b`). Von `TestModalityClassify` gepinnt; das RE2-ASCII-
  `\b`-Caveat für konfigurierte Umlaut-Rand-Keywords ist bewusst und **1:1** so
  im Vertrag (Spec Schritt 3) vermerkt — die Default-Menge hat nur ASCII-Ränder.
- **Normalisierung.** `normalizeBody` (`:517`) entfernt Emphasis-/Code-Marker (`*` und Backtick) und zieht via
  `strings.Fields`+`Join` alle Whitespace-/Umbruch-Folgen auf ein Leerzeichen
  zusammen ⇒ `**MUSS** NICHT` und `MUSS\nNICHT` klassifizieren als `may`
  (`TestModalityClassify`, TestNormalizeBody). Die Multi-Wort-Keywords tragen ein
  literales Leerzeichen (QuoteMeta escaped es nicht), das exakt zum normalisierten
  Einzel-Leerzeichen passt.
- **Rückabbildung.** `classify` lowercased den Treffer und schlägt in
  `kwToLevel` nach (`:506-509`); Schlüssel und Regex-Alternativen entstammen
  derselben `TrimSpace`-Keyword-Menge ⇒ Lookup trifft immer.

### Byte-Identität (Dimension 3) — sauber

- `TraceRow.Modality` trägt `json/yaml:",omitempty"` und wird in `traceRow` nur
  bei `modalityActive` gesetzt (`trace.go:185-189`) ⇒ ohne Block leer ⇒ im
  JSON/YAML entfällt das Feld.
- `TraceMatrix.ModalityActive`/`GatingOrphans` sind `json:"-" yaml:"-"` ⇒ nicht
  serialisiert.
- Reporter (`report.go:249-272`): ohne Coverage/Modality ergibt sich die
  identische 5-Spalten-Kopf- und Trenner-Zeile (`|---|---|---|---|---|`) und das
  identische Zeilenformat wie vor der Erweiterung; mit nur Coverage die
  identische 6-Spalten-Form aus slice-067. `orDash`/`joinOrDash` nutzen denselben
  Em-Dash `—`.
- CLI-stderr (`cli.go:205-215`): bei `!ModalityActive` exakt die alte Zeile mit
  `matrix.Orphans`; `GatingOrphans == Orphans` in diesem Fall ⇒ Exit-Verhalten
  unverändert.

### Fail-closed (Dimension 4) — vollständig, config-zeitig

`Decode` → `applyModules` → `applyTrace` → `applyModality` (`configyaml.go:317,
373, 411`). `validateModalityLevels` prüft leeren Namen (M1: ungetestet),
reserviertes `unknown`, leeres Keyword, Keyword-Dublette (case-insensitiv über
`seenKW`); danach validiert `applyModality` jeden `require-levels`-Eintrag gegen
{`unknown`} ∪ effektive Stufen-Namen. Reihenfolge korrekt (Levels vor
require-levels; `effective` = konfigurierte Levels sonst Defaults). Alle Pfade
liefern `error` ⇒ Exit 2 (von `TestCLI068_Modality_NegativeConfig` für 4 Fälle
belegt).

### Gating (Dimension 5) — korrekt

`BuildTraceMatrix` (`trace.go:162-167`): je Waise `Orphans++`; `GatingOrphans++`
gdw. `!ModalityActive` **oder** `mm.requireLevels[row.Modality]`. Kurzschluss
schützt den Nil-`requireLevels`-Lesezugriff im Inaktiv-Fall (und ein Nil-Map-Read
gäbe ohnehin `false`). Default `require-levels=[must]` in `resolveModality:486-492`.
KANN(`may`)-Waise gatet nicht, MUSS-Waise gatet — von
`TestCLI068_Modality_KannAdvisory`/`_KlassifikationUndGating` belegt.
