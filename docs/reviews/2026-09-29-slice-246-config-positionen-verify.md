# Verifikation — slice-246 (`.d-check.yml`-Positionen aus ai-harness-init evaluieren)

**Rolle:** Verifier (Modul 11) — Frage: „Bauen wir es richtig?" (gegen DoD/Spec/Plan)
**Gegenstand:** slice-246, Range `fc0d53ba..HEAD` (`e5787517` Claim/Move · `d21f7a76`
Evaluierung + slice-248 · `cf6f6908` R1-Report · `4d2244e3` R1-Einarbeitung ·
`eb998b60` Ruhe-Marker); Welle: welle-91.
**Eingangs-Kontext:** Slice-Plan §2/§7, R1-Report
(`2026-09-29-slice-246-config-positionen-r1.md`), eigene Config, Snapshot
`/tmp/aih-v6.13.0/.d-check.yml`. Nicht geprüft: Plan-vs-ADR-Qualität (Reviewer).
**Datum:** 2026-09-29 · **Modell-ID:** glm-5.3-flash
**Methode:** jeder DoD-Punkt mit eigener Messung nachgefahren, nicht aus dem
Implementer-Bericht übernommen.

---

## DoD 1 — Je Position (5) eine belegte Entscheidung · **BESTÄTIGT**

| # | Entscheidung (§7) | Eigene Messung | Ergebnis |
|---|---|---|---|
| 1 | welle-Klasse bereits Bestand, schärfer (MR-034) | `.d-check.yml`: Klasse `welle` mit `paths: ["docs/plan/planning/**/welle-*.md"]` und `token: 'welle-\d{2,}'`; Regeln `spec-straten→welle`, `sicht→welle`, `adr→welle` je `allow: false`; Kommentar nennt MR-034/MR-006 | bestätigt — d-check führt das schärfere `{2,}`-Token mit **drei** Regeln gegen die Schwester (Präfix-Token, zwei Regeln) |
| 2 | aussen — Grundsatz bejaht, Umfang an slice-248 delegiert | `docs/plan/planning/open/slice-248-matrix-aussen-adaptionsblock.md` existiert und trägt die vollständige gemessene Menge (28 Links, 1 Token, 11 Status-Fälle) in §1 und DoD | bestätigt |
| 3 | adaptionsblock — dieselbe Lage, im selben Folge-Slice | wie oben; slice-248 §3 plant Klassen-Ordnung `adaptionsblock` vor `aussen` | bestätigt |
| 4 | segment-tolerante ids abgelehnt — null Treffer | `grep -rEn '\bADR-[A-Z]+-[0-9]{4}\b' spec/ docs/ harness/ --include='*.md'` → **0 Treffer** | bestätigt — die Ablehnungsbasis stimmt |
| 5 | exclude-sections-Scoping abgelehnt — kein eigener Fall | Headings gemessen: `spec/lastenheft.md:3948` und `spec/spezifikation.md:3525` heißen `## 7. Historie`; ADRs heißen `## Geschichte`; matrix `exclude-sections: [Geschichte]` ist global | bestätigt — die globale Ausnahme trifft die Spec-Historie nicht, die Ablehnung „kein eigener Fall" hält |

## DoD 2 — Befund-Wirkung gemessen · **BESTÄTIGT (Probe selbst wiederholt)**

Probe nach der Aufgaben-Vorgabe reproduziert: Klassen `aussen` (`paths: ["**"]`) als
**letzte** Klasse, `adaptionsblock` davor, Regeln `spec-straten→aussen`/`sicht→aussen`
und `spec-straten→adaptionsblock`/`sicht→adaptionsblock`, ids-Regex
`ADR-([A-Z]+-)?\d{4}` — Einfügung per awk, Lauf, Revert (`git checkout -- .d-check.yml`,
Baum danach wieder sauber).

**Lauf:** `docker run --rm --network none … d-check:latest --enable matrix --enable ids
--disable …` — Exit 1 (erwartet, Probe), Ausgabe: **`d-check: 916 Datei(en) geprüft,
40 Befund(e)`**.

| Klasse | gemessen (ich) | behauptet (§7/R1) |
|---|---|---|
| matrix-forbidden gesamt | **29** | 29 |
| — aussen-Links aus Straten/Sicht | **28** (21 lastenheft, davon 20 Links + 1 Token · 7 spezifikation · 1 architecture) | 28 (20/7/1) |
| — nacktes MR-Token | **1** — `spec/lastenheft.md:4069` (`MR-004`) | 1, dieselbe Zeile |
| matrix-inactive gesamt | **11** | 11 |
| — Verteilung | `docs/plan/adr/README.md` ×7 · `CHANGELOG.md` ×3 · `docs/user/releasing.md` ×1 | identisch |
| **Gesamt** | **40** | **40** |

Ziel-Verteilung der 28 Links in meinem Lauf (Ordnung korrigiert, `aussen` zuletzt):
14 × `harness/conventions*` (= die 14 adaptionsblock-gelabelten — die First-Match-
Labelverschiebung, die R1 Probe A/B gegenüberstellte), 6 × `AGENTS.md`, 2 ×
Baseline-Zitat, je 1 `harness/README.md`, `packaging/dockerhub/README.md`, Carveout
CO-001, Register, 2 × CR — Summe 28, deckt die §7-Enumeration exakt. **Die Zählung
entscheidet: 40 — bestätigt.** Die Messung trägt die Delegationsentscheidung der
Positionen 2+3: eine Adoption im Slice wäre an den 28 Link- und 11 Status-
Entscheidungen gescheitert (DoD 1 erlaubt Adoption nur bei grünen Gates).

## DoD 3 — `make gates` grün · **BESTÄTIGT**

Eigener Lauf, nicht der behauptete Exit-Code: `make gates` → **Exit 0**, Abschlusszeile
„[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check +
coverage-gate + semgrep + gate-consistency + planning-check green" (alle zehn Glieder).
Lauf auf **cleanem Arbeitsbaum** (`git status --porcelain` leer) exakt an `eb998b60`;
`record-gates`-Hash `.harness/state/gates-passed.diffsha` = `93658efb…` stimmt mit
`tools/harness/working-tree-hash.sh` überein. Commit-Datum `eb998b60` 14:17:18 +0200,
Gates-Lauf 14:23:52 — die Aussage ist nicht älter als der letzte Commit und zudem
unabhängig neu belegt.

## Review-Vorgeschichte (Stichproben)

- **F-1 (Zahlenbasis):** §7 trägt jetzt **40/29/28/1/11** konsistent (28 Links + 1
  Token = 29; 29 + 11 = 40) und nennt die größte unterbliebene Klasse (`ADR-README`
  ×7) wie CHANGELOG ×3 — aufgelöst in `4d2244e3`.
- **F-2 (slice-248 Status-Fälle):** slice-248 §1 nennt alle **11** matrix-inactive-Fälle
  (ADR-Index ×7, CHANGELOG ×3, `docs/user/releasing.md` ×1) und hält sie in §3 als
  eigene Plan-Zeile; DoD 3 verlangt je Ausnahme die §3.6-ADR — aufgelöst.
- **F-3 (LOW, Widerspruch DoD-2-Zeile vs. §7):** die DoD-2-Annotation sagt jetzt
  „Grundsatz bejaht, Umfang an slice-248 delegiert" — deckungsgleich mit §7. Aufgelöst.
- **INFO-1:** welle-91 Out-of-Scope trägt die Zeile zur wellenlosen Laufbahn von
  slice-248 — die Trennung ist explizit. Aufgelöst.

## Plan-vs-Code-Diff

- **Geplant, nicht umgesetzt — korrekt:** `.d-check.yml`-Update (§3). Die Zeile war
  conditional („nach Evaluierung … adoptierte Positionen"); keine Position wurde
  adoptiert, die Datei ist unverändert — der Plan ist nicht still erweitert, keine
  Gate-Senkung (`exclude-sections`, `ids`, `matrix` unverändert).
- **Umgesetzt:** Plan-Update (§2/§7, Anker/Evaluierung), Folge-Slice slice-248 in
  `open/` (167 Zeilen), R1-Report. Beides durch §4/§7 des Plans gedeckt.
- **Nicht geplant, aber prozessmechanisch:** Roadmap-Ruhe-Marker-Entfernung
  (`eb998b60` — planning-check-Kopplung an die Beanspruchung) und welle-91 +4 Zeilen
  (R1-INFO-1-Antwort). Keine Abgrenzungs-Ausweitung im Sachprogramm.

## Verdikt

**Alle drei DoD-Punkte bestätigt** — jede Entscheidung durch eigene Messung belegt,
die Probe reproduziert 40/29/28/11 exakt, `make gates` von mir an `eb998b60` grün
gefahren. Der Slice ist closure-fähig. Vor dem `git mv` gilt unverändert der
R1-Hinweis: die Ruheort-Formen in §7 (`../open/…`, `../observations/`,
`../../../../harness/…`) sind im Closure-Arbeitsbaum auf `done/welle-91/` umzuzielen —
der pre-commit-`doc-check` fängt ein Vergessen.

**Grenzen dieser Verifikation:** meine Probe lief in der korrigierten Ordnung
(`aussen` zuletzt) — identische Summen wie §7s Probe A, Label-Verschiebung wie von
R1 beschrieben; die Exit-0-Aussage des Implementers ist damit ersetzt durch einen
eigenen, neueren Lauf auf denselben Inhalt.
