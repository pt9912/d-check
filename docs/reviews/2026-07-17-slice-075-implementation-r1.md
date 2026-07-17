# Review slice-075 (R1) — Komma-Kurzform fail-closed

**Datum:** 2026-07-17 · **Rolle:** unabhängiger Reviewer (kontext-getrennt,
nicht der Autor, ohne Zugriff auf die Implementierungs-Analyse) · **Lauf:** R1,
**vor** Release v0.46.0.

**Gegenstand:** Commit-Range `0edee0a..HEAD`, `HEAD` = `250f270`. Kern-Code
`internal/hexagon/core/app/trace.go` (`expandRange`, `commaShortform`,
`expandNumericRange`, `rangeAwareIDs`). Feat `319b0f2`, Schärfung `909c781`
(inline-Verengung), Refactor `85261fd`, Doc-first `a30318c`/Release-Prep
`9a6062b`, SDP-Cleanup `250f270`.

**Kontext-Stabilität:** `git status --porcelain` am Ende leer, `HEAD` unverändert
`250f270`. Der Scratch-Worktree für die Mutationen wurde mit
`git worktree remove --force` entfernt; der Produktivbaum blieb unberührt.

**Quellen:**
[`DC-FA-COV-001`](../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
(Lastenheft 0.46.0),
[`DC-FA-XREF-001`](../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in),
[`DC-FA-COV-001.a`](../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)
Schritt 3, [ADR-0041](../plan/adr/0041-komma-kurzform-fail-closed.md) (Proposed),
[`AGENTS.md`](../../AGENTS.md) §3.

**Verifikations-Basis:** `make build` (HEAD `250f270`), Image
`sha256:826c7ea7ac02b516e53f70c8fec2da5e091b8abdf225d01eb0d6259da52f1a62`.
Fixtures gefahren mit
`docker run --rm --network none -v <fixture>:/repo:ro -w /repo d-check:latest --trace`
(Config: `modules: []` + `trace.requirements.source`/`id-pattern` +
`trace.coverage[]` bzw. `trace.cross-consistency`). Mutationen via
`make -C <worktree> test` (Docker `--target test`, `--no-cache-filter test`) in
einem Scratch-Worktree; Baseline dort zuvor grün.

---

## Findings

### F-1 · LOW · Handbuch-Kopf lag hinter der eigenen Release-Prep zurück — Handbuch-Version und Software-Version widersprechen der jüngsten Historienzeile und dem eigenen Anker

**quelle:** Maintainability (Doku-Currency, Release-Prep-Hygiene);
[`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) berührt nicht,
kein Gate deckt die Kopf-Zeile ab.

**pfad:** `docs/user/benutzerhandbuch.md:3`.

**befund:** Die Kopf-Zeile lautet
`**Handbuch-Version:** 1.33 · **Software-Version:** [v0.45.1](../../version.md#v0.46.0)`.
Der Release-Prep-Commit `9a6062b` hat im selben Zug (a) die Historienzeile `1.34`
für `v0.46.0` ergänzt (`benutzerhandbuch.md:1507`), (b) `version.md` auf
`Aktuelle Version: v0.46.0` gehoben und (c) den **Anker** der Kopf-Zeile von
`#v0.45.1` auf `#v0.46.0` umgebogen — den **Linktext** `v0.45.1` und die
`Handbuch-Version: 1.33` aber stehen gelassen. Folge: der sichtbare Text `v0.45.1`
widerspricht seinem eigenen Anker `#v0.46.0`, der Historie (`1.34`) und
`version.md` (`v0.46.0`). Failure-Szenario: der Leser des ausgelieferten
Handbuchs entnimmt dem Kopf „Software-Version v0.45.1 / Handbuch 1.33", obwohl das
Dokument die Komma-Kurzform (v0.46.0, §5 Zeile 743 „Ab v0.46.0.") bereits
beschreibt — eine Currency-Lüge im Kopf. `git log -S "Handbuch-Version:** 1.34"`
liefert keinen Treffer; die Kopf-Zeile wurde nie auf `1.34/v0.46.0` gezogen.

**verifizierbar:** ja, aber **nicht** durch einen Gate-Lauf (der `versions`-Gate
prüft nur ghcr-/Toolchain-Pins, nicht die Handbuch-Kopf-Zeile; die
`version.md#aktuell`-Kopplung gilt dem CHANGELOG, nicht dem Handbuch). Manuell:
`sed -n '3p;1507p' docs/user/benutzerhandbuch.md` gegen `version.md` §Aktuell.

### F-2 · INFO · Eine dritte, undokumentierte Rest-Klasse „stiller Drop": die Komma-Kurzform feuert nur hinter ziffern-endenden Kennungen — eine buchstaben-suffigierte Kennung (vom Default-`id-pattern` erlaubt) plus `, <Ziffern>` fällt weiter still

**quelle:** [`DC-FA-COV-001.a`](../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)
Schritt 3 (spricht allgemein von „Folgt der Fundstelle … ein Komma und
unmittelbar darauf Ziffern", ohne die Kennung auf Ziffern-Ende zu qualifizieren);
[ADR-0041](../plan/adr/0041-komma-kurzform-fail-closed.md) §Konsequenzen
(dokumentiert **zwei** bewusste Scope-Grenzen derselben Klasse, diese dritte
nicht).

**pfad:** `internal/hexagon/core/app/trace.go:519-523` (`expandRange`:
`d := trailingDigits.FindStringIndex(id); if d == nil { return nil, nil }`) im
Zusammenspiel mit `trace.go:550` (der Komma-Check sitzt **hinter** diesem
Früh-Return).

**befund:** Der Komma-Kurzform-Check liegt in `expandRange` **nach** dem
Früh-Return für Kennungen, die nicht auf Ziffern enden (`trailingDigits = \d+$`).
Das Default-`requirements.id-pattern`
(`[A-Z][A-Z0-9]*-(?:FA-[A-Z]+|QA)-\d+[A-Za-z]?`) erlaubt aber ein
Buchstaben-Suffix. Eine solche Kennung, gefolgt von Komma und Ziffern, umgeht den
Check und lässt die Ziffern **still** fallen — genau der Defekt, den der Slice
tilgen soll. Belegt gegen das Image (Fixture `fx/letter`, Default-Pattern,
`ranges: true`):

```text
Quelle:  Abdeckung: DC-QA-05a, 007 deckt.
Ergebnis: Exit 0
  | DC-QA-05a | … | Trace | ok   |
  | DC-QA-007 | … | —     | WAISE |   ← 007 still gefallen, kein Exit 2
```

Failure-Szenario: ein Konsument mit buchstaben-suffigierten Anforderungs-IDs
schreibt `XX-YY-05a, 007` in eine `trace.coverage`-Quelle; `007` verschwindet
lautlos, kein Signal. Die Klasse ist eng (grid-gyms Realform nutzt `\d{3}`, das
Ziffern-Suffix ist kein Realdaten-Auslöser), aber sie ist von derselben Art wie
die beiden in ADR §Konsequenzen **bewusst** offen gelassenen Rest-Klassen
(Vor-Komma-Whitespace, Zeilenumbruch) — und dort ausdrücklich als „ein schmaler
Rest derselben Klasse" benannt. Diese dritte ist weder in der ADR noch im
Code-Kommentar erwähnt.

**verifizierbar:** ja, empirisch reproduziert (Fixture oben). Kein bestehender
Test deckt die buchstaben-suffigierte Kennung ab — die Suite gatet den Punkt
nicht.

---

## Negativbefunde (geprüft, ohne Befund)

- **Drei Positionen (nackte Kennung / hinter Range / hinter Enum):** geprüft an
  Code (`trace.go:527-552`: `rest` wird nach `..`- und `/`-Konsum vorgeschoben,
  der Komma-Check greift auf dem Rest) **und** empirisch —
  `GG-SCN-001, 007` / `GG-SCN-001..005, 007` / `GG-SCN-001/003, 007` liefern je
  Exit 2 mit dem Notations-Hinweis. Kein Befund.
- **Realer grid-gym-Fall `GG-SCN-001..005, 007, 008`:** Exit 2, Meldung
  `trace.coverage: Komma-Kurzform hinter GG-SCN-001 …`; die Range wird konsumiert,
  der Komma-Schwanz `007, 008` bricht ab, **kein** stiller 3-Waisen-Drop. Der
  partielle `out` aus der Range-Expansion wird mit dem Fehler verworfen
  (`rangeAwareIDs` propagiert `nil, err`) — keine Teil-Coverage leckt. Kein Befund.
- **Saubere Range / Range + volle Kennung ⇒ Exit 0:** `GG-SCN-001..005` deckt
  001–005 (007/008 korrekt Waisen); `GG-SCN-001..005, GG-SCN-007` deckt 001–005
  **und** 007 (Komma vor vollständiger Kennung hinter Range unberührt, beide
  Kennungen regulär gefunden), 008 Waise. Beide Exit 0. Kein Befund.
- **Scope-Grenze (a) Zeilenumbruch (`[ \t]` statt `\s`):** Komma am Zeilenende,
  ziffern-beginnende Folgezeile (`GG-QA-001,\n2026 …`) ⇒ Exit 0, `GG-QA-001`
  gedeckt, `2026` erzeugt keine Falsch-Coverage. Verhalten = ADR-Doku; die Ziffer
  fällt still, wie dort ehrlich beschrieben. Kein Befund.
- **Scope-Grenze (b) Whitespace vor dem Komma:** `GG-QA-001 , 007` ⇒ Exit 0,
  `GG-QA-001` gedeckt, `GG-QA-007` **Waise** (007 fällt still). Verhalten =
  ADR-Doku (§Konsequenzen „die `007` fällt dort weiter still"). Die Dokumentation
  ist ehrlich. Kein Befund.
- **Beide Konsumenten (coverage + cross-consistency, beide Richtungen):**
  `trace.coverage` ⇒ Exit 2 (Fixtures oben). `trace.cross-consistency` **backward**
  (Bezug-Zelle `GG-SIM-001, 007`) ⇒ Exit 2 mit Präfix
  `trace.cross-consistency.backward`; **forward** (Anforderungs-Zelle
  `GG-SIM-001, 007`) ⇒ Exit 2 mit Präfix `trace.cross-consistency.forward`;
  saubere Zellen ⇒ Exit 0, „0 Differenz(en)". Beide Richtungen laufen über den
  geteilten `rangeAwareIDs` (`trace_cross.go:252`/`:278`). Kein Befund.
- **`ranges: false` unberührt:** die Komma-Kurzform gehört zur Range-Fähigkeit
  (`rangeAwareIDs` kehrt vor der Expansion zurück, `trace.go:501`); mit
  `ranges: false` kein Fehler. Getestet (`TestCoverageRefsCommaShortform`, dritter
  Zweig) und code-belegt. Kein Befund.
- **Fehlermeldung:** verweist auf die zwei zugesagten Notationen (`…..BBB oder
  …/BBB`) und trägt das Config-Feld (`trace.coverage` /
  `trace.cross-consistency.forward|backward`) — actionable. Im gemischten Fall
  nennt sie die Grund-Kennung `GG-SCN-001` (nicht die Position des Komma nach
  `..005`); leichte Unschärfe, aber der Hinweis bleibt umsetzbar. Kein Befund.
- **Mutations-Härte (Scratch-Worktree, `make test`):**
  - Regel entfernt (Komma-Check neutralisiert) ⇒ 9 Zeilen kippen:
    `TestExpandRangeCommaShortform` (7 Subtests: Komma-Kurzform, ohne Space, Tab,
    Prosa-Zahl, verlinkt, Schwanz hinter Range, Schwanz hinter Enum) +
    `TestCoverageRefsCommaShortform` + `TestCrossConsistencyCommaShortform`.
  - `\d`-Grenze entfernt (`^,[ \t]*\d` → `^,`) ⇒ 4 Subtests kippen: „Komma vor
    voller Kennung", „Komma vor Prosa ohne Ziffer", „Range dann volle Kennung",
    „Newline dann Ziffer".
  - `[ \t]`-Grenze aufgeweitet (`[ \t]` → `\s`) ⇒ genau „Komma dann Newline dann
    Ziffer (zwei Zeilen)" kippt.
  Jede der drei zurückgedrehten Grenzen kippt mindestens einen Test. Kein Befund.
- **Referenz-Richtung (SDP), Marker-Ehrlichkeit:** `slice-075` erscheint in
  ADR-0041 ausschließlich in `## Geschichte` (Zeilen 111/112, „Umsetzender Slice
  slice-075") — ein Provenance-/Verifikations-Zeiger (wo umgesetzt), **keine**
  getarnte Entscheidungsgrundlage; der Entscheidungs-Körper (Zeilen 32–56) trägt
  keinen Slice-Token. Der jüngste Commit `250f270` ist genau diese SDP-Bereinigung.
  Kein Provenance-Marker `<!-- d-check:status-provenance -->` gesetzt (auch nicht
  nötig, da §Geschichte). Kein Befund.
- **Vertrags-/Doku-Konsistenz (außer F-1):** Lastenheft §DC-FA-COV-001 (Negative-
  + Boundary-Kriterium, Version 0.46.0 + Historienzeile), Spezifikation
  §DC-FA-COV-001.a Schritt 3 (drei Positionen, `[ \t]`-Begründung), ADR-0041
  Fitness-Funktion, CHANGELOG [0.46.0] und Handbuch §5 (Zeilen 734–743, inkl.
  „hinter einer Range") stimmen mit dem beobachteten Verhalten überein. Kein
  Befund.

---

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH   | 0 |
| MEDIUM | 0 |
| LOW    | 1 (F-1) |
| INFO   | 1 (F-2) |

---

## Verdikt

**ACCEPT-WITH-NITS.**

Die Kern-Regel ist spec- und ADR-treu: die Komma-Kurzform feuert an allen drei
Positionen (nackte Kennung, hinter Range, hinter Enum), ein Komma vor einer
vollständigen Kennung — auch hinter einer Range — bleibt unberührt, der reale
grid-gym-Fall `GG-SCN-001..005, 007, 008` liefert Exit 2 statt drei stiller
Waisen, saubere Range und Range + volle Kennung liefern Exit 0. Beide bewusst
offen gelassenen Scope-Grenzen (Zeilenumbruch, Vor-Komma-Whitespace) sind
empirisch bestätigt und in ADR §Konsequenzen ehrlich dokumentiert. Beide
Konsumenten (`trace.coverage` und `trace.cross-consistency`, letzterer in beiden
Richtungen) sind laufzeit-korrekt; die Sensor-Härte hält gegen das Zurückdrehen
der Regel und beider Grenzen.

Kein HIGH/MEDIUM ⇒ kein Blocker. Empfehlung: **F-1** (Handbuch-Kopf) vor dem
Tag ziehen — es ist eine reine Release-Prep-Kosmetik derselben Sitzung, und der
Kopf soll beim Auslieferungsstand nicht `v0.45.1` behaupten. **F-2** ist ein
INFO-Hinweis (undokumentierte dritte Rest-Klasse der stillen Drops); er blockiert
nicht, gehört aber — der Ehrlichkeit der beiden schon dokumentierten Grenzen
halber — als Won't-Fix-/Scope-Notiz sichtbar gemacht.
