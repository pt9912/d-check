# Review slice-077 (R1) — Stiller Tabellen-Übersprung: Grenze am relevanten Header

**Datum:** 2026-07-18 · **Rolle:** unabhängiger Reviewer (kontext-getrennt,
nicht der Autor) · **Lauf:** R1, **vor** Release (SemVer-Minor).

**Gegenstand:** `3ea0090` (feat: Tabellengrenze am relevanten Header), Diff
`18d00de..HEAD`. Doc-first-Commits `af25012`/`a05b2f9`/`69fcddc`
(Plan/Spec/ADR), Code in `3ea0090`.

**Quellen:** [`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen)
Schritt 5, [`DC-FA-XREF-001.a`](../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency),
[ADR-0043](../plan/adr/0043-tabellengrenze-am-relevanten-header.md) (Proposed),
[ADR-0037](../plan/adr/0037-trace-tabellenquellen-nullmengen-guard.md),
[slice-077](../plan/planning/done/slice-077-stiller-tabellen-uebersprung.md),
[Review R3](2026-07-17-slice-074-implementation-r3.md) (F-1 Wiederaufsetz-Punkt,
F-5 Masken-Blindheit, F-2 Mutations-Härte), [`AGENTS.md`](../../AGENTS.md) §3.

**Verifikations-Basis:** `make test` (HEAD, grün) und dieselbe Suite mit **neutralisierter
Grenze** (`if false && isRelevantHeader …`), beide in einem Scratch-Worktree (entfernt;
Produktivbaum unverändert, `git status` leer). Runtime-Images HEAD
(`d-check:latest`, `sha256:80c432ccd360` — identisch mit dem HEAD-Build) und Elter
`18d00de` (`d-check-parent`, frisch gebaut). Läufe
`docker run --rm --network none -v <fixture>:/repo:ro -w /repo <image> --trace`.
Realdaten: grid-gym **read-only gemountet** (`:ro`; grid-gym selbst nur gelesen,
sein `git status` blieb leer).

**Vorbemerkung.** Der Kern-Fix trägt und ist gegen die Images belegt: `fx-s`
(irrelevante Tabelle verschluckt die folgende relevante) liefert unter dem Elter
`18d00de` still `1 Anforderung(en), 1 Waise(n)` (R-2/R-3 verschluckt) und unter HEAD
`3 Anforderung(en), 3 Waise(n)`. Die Grenze steht **vor** dem Zellenzahl-Check
(`consumeTableRows` Zeile 167 vor Zeile 176), greift also auch bei passender
Zellenzahl — konform zu Spec Schritt 5 und ADR-0043 Entscheidung 1. Sie ist an
**beide** Konsumenten durchgereicht (`bindTableColumns` für `format: table`;
`bindCrossColumns` für forward/backward). Die Befunde unten sind Randschärfe der
Mutations-Härte und ein **vorbestehender** Rest-Übersprung — kein Bruch des
Kern-Fixes.

---

## Findings

### R-F-1 · LOW · `fx-adj` (`TestCLI077_BenachbarteRelevanteTabellen`) behauptet einen Mutations-Pin, der empirisch **nicht** hält — die Grenze lässt sich dort geräuschlos zurückdrehen

**quelle:** [slice-077](../plan/planning/done/slice-077-stiller-tabellen-uebersprung.md)
§3 DoD („**Mutations-Härte:** jede neue Grenze kippt einen Test — **gemessen, nicht
zugesagt** (die R3-F-2-Lehre)"); Reviewer-Anker LOW (latente Wartungsfalle);
[R3-F-2](2026-07-17-slice-074-implementation-r3.md) (dieselbe Klasse: eine Zusage
„per Mutation gepinnt", die nicht hält).

**pfad:** `internal/adapter/driving/cli/cli_acceptance_test.go:2317-2336`
(Kommentar `:2319` „Mutations-Pin: ohne die Grenze kippt dieser Test auf Total 1").

**befund:** Mit neutralisierter Grenze (`make test` im Scratch-Worktree) kippen
**nur** `TestCLI077_StillerUebersprungGrenze` und
`TestCLI077_Cross_RuecksichtNichtVerschluckt`. `TestCLI077_BenachbarteRelevanteTabellen`
bleibt **grün** — entgegen seinem Kommentar. Grund: die verschluckte zweite Tabelle
(`| ID | Anforderung |` / Trennzeile / `| R-2 | b |`) hat exakt die Zellenzahl der
ersten; ohne Grenze wird `| R-2 | b |` als Datenzeile der ersten Tabelle gelesen und
liefert R-2 trotzdem. Der Test prüft nur `Total == 2` und sieht `2` **mit und ohne**
Grenze. Er pinnt die Grenze also nicht; er testet eine Eigenschaft, die
grenzenunabhängig gilt. Die Grenze selbst ist durch `fx-s` (`format: table`) und den
Cross-Test (backward) gepinnt — die Klasse ist nicht ungeschützt, aber die
Pin-Zusage an genau dieser Stelle ist falsch (dieselbe Bewegung wie R3-F-2).

**verifizierbar:** ja — `make test` mit `if false && isRelevantHeader …` in
`consumeTableRows`: 2 FAIL (fx-s, cross), `BenachbarteRelevanteTabellen` PASS.

---

### R-F-2 · LOW · Die Prädikat-Durchreichung an den **Vorwärts**-Cross-Pfad (`bindForwardTables`) ist ungetestet — nur der Rückwärts-Pfad wird adversarial gepinnt

**quelle:** [slice-077](../plan/planning/done/slice-077-stiller-tabellen-uebersprung.md)
§2 („das Relevanz-Prädikat an **beide** Konsumenten durchgereicht … `cross-consistency`
via `bindCrossColumns`") und §3 DoD (Akzeptanztests „auf **Konsumenten-Ebene**");
Reviewer-Anker MEDIUM/LOW (Negativtest-Lücke am neuen Vertrag) — hier LOW, weil der
Code-Pfad korrekt ist und nur der Sensor fehlt.

**pfad:** `internal/hexagon/core/app/trace_cross.go:162-166` (`bindForwardTables`,
Prädikat durchgereicht) gegen `internal/adapter/driving/cli/cli_acceptance_test.go:2347-2351`
(Fixture: `docs/traceability.md` ist eine **einzelne, saubere** Tabelle; nur
`spec/architecture.md` = backward trägt die irrelevante Vor-Tabelle).

**befund:** `TestCLI077_Cross_RuecksichtNichtVerschluckt` belegt die
Prädikat-Durchreichung ausschließlich für den **Rückwärts**-Pfad (irrelevante
`Werkzeug/Zweck`-Tabelle vor der relevanten `Komponente/Bezug`-Tabelle; unter der
Mutation Exit 2). Der **Vorwärts**-Pfad hat keine Fixture, in der eine irrelevante
Tabelle einer relevanten forward-Tabelle unmittelbar vorausgeht. Eine Mutation, die
`bindForwardTables` das Prädikat entzöge (nil oder falsche Spalten), bliebe damit
unbemerkt, obwohl sie das Verhalten ändert (verschluckte forward-Tabelle ⇒
`len(out)==0` ⇒ Exit 2 bzw. verlorene Kanten ⇒ Phantom-Differenzen).

**verifizierbar:** ja — `make test` mit `isRelevant`-Ersetzung durch `nil` allein in
`bindForwardTables` bleibt grün; eine forward-Swallow-Fixture (irrelevante Tabelle
ohne Leerzeile vor der relevanten forward-Tabelle) bände den Pfad.

---

### R-F-3 · LOW · Rest-Übersprung: ein **mehrdeutig-relevanter** Header (doppelte Rollen-Spalte) direkt hinter einer irrelevanten Tabelle bleibt still verschluckt — dieselbe Tabelle allein wäre Exit 2

**quelle:** [ADR-0043](../plan/adr/0043-tabellengrenze-am-relevanten-header.md)
§Entscheidung 3 („jede relevante Tabelle wird erkannt") und §Konsequenzen/„Bewusst
offen gelassen" (nennt nur **zwei benachbarte irrelevante** Tabellen);
[`DC-FA-REQ-001.a`](../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen);
Reviewer-Anker LOW (latente Wartungsfalle / undokumentierte Annahme).

**pfad:** `internal/hexagon/core/app/trace_table.go:34-37` und `:168`
(`isRelevantHeader := … return err == nil && ok`) im Zusammenspiel mit
`bindTableColumns` `:242-244` (doppelter Rollen-Header ⇒ `(_, false, err)`);
analog `trace_cross.go:162-165`/`193-196` über `bindCrossColumns` `:398-400`.

**befund:** Das Grenz-Prädikat wertet einen `bindTableColumns`/`bindCrossColumns`-**Fehler**
(Rollen-Header kommt mehrfach vor) als „nicht relevant" (`err == nil && ok` ⇒ false)
und **feuert nicht** — fail-**open** an dieser Kante. Folge (gegen HEAD gemessen,
`fxdup`: relevante Tabelle R-1, Leerzeile, irrelevante `Werkzeug`-Tabelle, direkt
gefolgt von `| ID | ID | Anforderung |` + Trennzeile + `| R-2 | … |`): still
`1 Anforderung(en), 1 Waise(n)` — R-2 verschwindet lautlos. Dieselbe Tabelle
**allein** (`fxdup2`, Leerzeile davor) liefert Exit 2 („konfigurierter Header \"ID\"
kommt mehrfach vor"). Also: Exit 2 standalone, still bei Nachbarschaft — genau die
Silent-Loss-Klasse, deretwegen dieser Slice existiert. Der Fall ist **vorbestehend**
(Elter `18d00de` liefert `fxdup` byte-identisch — **kein slice-077-Regress**) und
verlangt eine fehlgeformte (doppelt-spaltige) Eingabe; ADR-0043 ist in sich
konsistent (ein Duplikat „bindet" per Definition nicht), doch seine
„Bewusst-offen"-Aufzählung führt diesen Rest **nicht**.

**verifizierbar:** ja — `fxdup` (still `1/1`, HEAD == Elter) gegen `fxdup2`
(Exit 2), beide gegen `d-check:latest`.

---

## Negativbefunde (geprüft, ohne Befund)

- **Reihenfolge im Loop (Grenze vor Zellenzahl-Check): korrekt.** Der
  `isRelevantHeader`-Block (`trace_table.go:167-171`) steht **vor** dem
  `len(cells) != len(t.header)`-Check (`:176`). `fx-s` (HEAD `total 3` vs. Elter
  `total 1`) belegt, dass ein relevanter Header auch bei passender Zellenzahl trennt
  (Spec Schritt 5 / ADR-0043 Entscheidung 1).
- **Diskriminator „relevant" — `fx-t` (all-dashes) trennt nicht: korrekt.**
  `TestCLI077_AllDashesKeineFalschtrennung` bleibt unter der Mutation grün (Total 2
  mit **und** ohne Grenze): `| - | - |` ist kein neuer Header, weil die Folgezeile
  (`| R-2 | b |`) keine Trennzeile ist ⇒ `tableHeaderAt` false; die Grenze feuert
  nie. Rein strukturelle Grenzen scheiterten hier — diese nicht.
- **Datenzeile, die wie ein relevanter Header aussieht (in einer relevanten
  Tabelle): kein Anforderungsverlust.** Der `fx-adj`/`fx-t2`-Fall spaltet die Tabelle,
  liest aber alle Zeilen (Tabelle 1 **oder** 2); die als Header umgedeutete Zeile
  trägt die Spaltennamen (`ID`, …), deren ID-Zelle nicht auf `id-pattern` passt und
  ohnehin keine Anforderung definierte. Ein Verlust bräuchte einen Header-Namen, der
  selbst auf `id-pattern` passt (pathologische Config) — nicht realistisch.
- **Zwei benachbarte irrelevante Tabellen: benigne, wie in ADR-0043 dokumentiert**
  — keine Rolle gebunden, keine Anforderung betroffen.
- **Beide Konsumenten / richtige Spalten: korrekt.** `format: table` reicht
  `bindTableColumns(header, cfg)` durch (deckt sich mit `extractTable`s
  Relevanz); forward `bindCrossColumns(header, ReqColumn, DesignColumn)`; backward
  `bindCrossColumns(header, EdgeColumn)` — Letzteres nutzt **nur** die Kanten-Spalte,
  exakt die Relevanz-Definition aus `bindBackwardTables` (die Artefakt-ID-Spalte wird
  separat über `backwardIDColumn` aufgelöst und ist nicht Teil der Relevanz).
- **Masken-Blindheit (R3-F-5) adressiert; fail-closed.** `consumeTableRows` prüft
  `maskAllows(mask, lines[j].no)` (Schleifenkopf) **und** `maskAllows(mask,
  lines[j+1].no)` (Grenz-Guard `:167`) **vor** `tableHeaderAt(lines, j)` — dieselbe
  Zwei-Zeilen-Maskenprüfung wie `markdownTables` an der Spitze. Liegt `j+1`
  außerhalb der Maske, feuert die Grenze nicht; die Zeile wird höchstens als
  Datenzeile gelesen, und der Schleifenkopf schneidet an `j+1` — kein stiller Pfad
  (v0.47.0-Richtung).
- **Panic/Bounds: sicher.** Der Guard `j+1 < len(lines)` steht vor `tableHeaderAt(lines,
  j)`, das `lines[j+1]` liest; leere Datei (`markdownTables` iteriert nie), Tabelle am
  Dateiende (Guard greift) und Re-Scan (`i = next-1`, `next >= start = i+2` ⇒ `i`
  schreitet je Außeniteration streng fort) sind ohne Panic/Endlosschleife.
- **Regression grid-gym: keine.** HEAD und Elter `18d00de` liefern **byte-identisch**
  (beide brechen an der vorbestehenden slice-075-Komma-Kurzform in
  `traceability.md` ab — unabhängig von slice-077). Anmerkung: grid-gyms Config nutzt
  **weder** `format: table` **noch** `cross-consistency` (Headings + Coverage-Regex
  über `rules.SelectSections`), der geänderte Reader `markdownTables` wird dort **gar
  nicht** betreten — die Realdaten belegen Regressfreiheit, nicht den geänderten Pfad
  (deckt sich mit dem Slice-Risiko „die Vorbedingung kommt in grid-gym nicht vor").
- **Determinismus (DC-QA-02): gegeben** — die Grenze hängt nur an den konfigurierten
  Header-Namen und am Inhalt; `bindTableColumns`/`bindCrossColumns` sind seiteneffektfrei.
  Gleiche Config + gleiche Eingabe ⇒ gleiche Grenzen.
- **Mutations-Härte der Grenze **als solcher**: gegeben** — die neutralisierte Grenze
  kippt `fx-s` (table) **und** den Cross-Test (backward). Die Grenze ist also nicht
  wholesale still zurückdrehbar (die Einschränkungen sind R-F-1/R-F-2).
- **ADR-0005-Import-Regeln / Modul-Layout: ohne Befund** — keine neuen Imports;
  Änderungen ausschließlich in `core/app` (Closures + Signaturerweiterung).
- **[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit):
  ohne Befund** — nur `regexp`/`strings`; alle Läufe `--network none`.
- **Gate-Suppression / Schwellen-Senkung ([`AGENTS.md`](../../AGENTS.md) §3.6):
  ohne Befund** — keine `//nolint`, keine Gate-/Schwellen-Änderung.
- **Referenz-Richtung (SDP): ohne Befund** — ADR-0043 und die README-Zeile verweisen
  aufwärts (`DC-FA-REQ-001`/`DC-FA-XREF-001`/ADR-0037/0040); slice-077 erscheint nur
  als Provenance (`## Geschichte`, „Umsetzender Slice"). Keine offensichtliche
  `matrix-forbidden`-Kante (der Linter prüft es gate-seitig).
- **Doku-Konsistenz: ohne Befund** — Spec Schritt 5 (Grenze + Historie-Zeile),
  ADR-README-Zeile und Roadmap §Aktuelle Welle sind konsistent nachgezogen; nil-Prädikat
  wird von keinem Aufrufer übergeben (Grenze nie versehentlich inaktiv).
- **DoD-Abhakung / `make gates`: nicht geprüft** — Rolle der Verifikation
  (Reviewer-Skill §Anti-Pattern). R-F-1/R-F-2 betreffen die **Belastbarkeit** der
  Mutations-Zusage, nicht deren DoD-Abhakung.

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 0 | — |
| LOW | 3 | R-F-1, R-F-2, R-F-3 |
| INFO | 0 | — |

---

## Verdikt

**ACCEPT-WITH-NITS.**

Der Kern-Fix ist korrekt, zweckgenau und in der sicheren Richtung: eine irrelevante
Tabelle verschluckt die folgende **relevante** nicht mehr still (`fx-s` HEAD
`total 3` vs. Elter `total 1`, gegen die Images belegt). Die Grenze steht vor dem
Zellenzahl-Check (auch bei passender Breite), ist an **beide** Konsumenten mit den
richtigen Spalten durchgereicht, adressiert R3-F-5 (Masken-Prüfung auf `j+1` vor
`tableHeaderAt`, fail-closed), ist panic-/endlosschleifenfrei, deterministisch und
ohne grid-gym-Regress. Der Diskriminator „relevant" hält die `fx-s`≡`fx-t`-Mehrdeutigkeit
sauber (all-dashes trennt nicht). Die Grenze **als solche** ist mutations-gepinnt
(fx-s + Cross-backward kippen).

Kein HIGH, kein MEDIUM. Die drei LOW-Befunde blockieren das Release nicht, benennen
aber die Randschärfe, die die R3-F-2-Lehre gerade an diesem Slice einfordert:
**R-F-1** — der `fx-adj`-Kommentar behauptet einen Mutations-Pin, den die Messung
widerlegt (der Test ist grenzenblind); **R-F-2** — die Vorwärts-Cross-Durchreichung
ist ungetestet (nur backward adversarial gepinnt); **R-F-3** — ein vorbestehender,
undokumentierter Rest-Übersprung (mehrdeutig-relevanter Header hinter irrelevanter
Tabelle bleibt still, während standalone Exit 2). Empfehlung an die Übergabe: den
`fx-adj`-Kommentar korrigieren oder die Fixture grenzensensitiv machen (z. B.
unterschiedliche Breiten), eine Vorwärts-Swallow-Fixture ergänzen und den Rest aus
R-F-3 in ADR-0043 §„Bewusst offen gelassen" aufnehmen (oder das Grenz-Prädikat auf
den Bind-Fehler fail-closed feuern lassen). Keiner dieser Punkte liegt im
ausgelieferten stillen Gate-Pfad des Kern-Fixes — der ist zu.
