# Reviewer-Skill — d-check

**Version:** 1.2.0 · **Datum:** 2026-06-28 ·
**Baseline:** Agents-Digest Kurs-Welle 18 §8 (Output-Schema,
Kategorien-Semantik, Report-Pflicht); Referenz-Richtung (SDP) aus
`grundlagen-konventionen.md` §Referenz-Richtung — seit
[`DC-FA-MTX-003`](../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)
(token-mechanisiert) auf Marker-Ehrlichkeit abgespeckt.

## Eingangs-Kontext (Pflicht — sonst nicht reproduzierbar)

Der Reviewer erhält: Diff/Commit-Range, den Slice-Plan, die
betroffenen `DC-*`-Anforderungen, die referenzierten ADRs und die
Hard Rules ([`AGENTS.md`](../../AGENTS.md) §3). **Nicht** erhalten:
die DoD-Abhakung — Plan-/DoD-Konformität prüft die Verifikation
(getrennter Kontext, anderes Prüf-Artefakt).

## Repo-spezifische Anker pro Kategorie

- **HIGH** (blockiert Merge): Stilles-Grün-Pfad in einem Gate oder
  Gate-Skript (Harness-Lüge); Korrektheitsfehler in Kern-Modulen mit
  falschen Befunden/Exit-Codes; Verstoß gegen
  [ADR-0005](../../docs/plan/adr/0005-modul-layout-hexagon-ordner.md)-Import-Regeln;
  Gate-Suppression ohne ADR; Netzzugriff außerhalb `external`
  ([`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- **MEDIUM** (vor Merge zu klären): Spec-Treue-Lücke einer
  Messmethode; Konsistenz-Lücke **zwischen** Modulen derselben
  Eingabe-Klasse; Erkennungs-Differenz zur Alt-Tool-Familie
  ([`DC-QA-04`](../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)); fehlende Negativtests bei neuem öffentlichen Vertrag;
  **Referenz-Richtung (SDP) — Marker-Ehrlichkeit.** Seit
  [`DC-FA-MTX-003`](../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)
  mechanisiert `matrix` die Token-Referenz-Richtung: ein **undeklarierter**
  Abwärts-Token (Slice/Plan im ADR-/Spec-Körper) ist ein `matrix-forbidden`-Befund
  des Linters — nicht mehr deine Aufgabe. **Deine** ist die nicht grep-bare
  Resthälfte: trägt eine Referenz den Provenance-Marker
  `<!-- d-check:status-provenance -->`, prüfe, ob die Deklaration **ehrlich** ist —
  *zeigt* sie, wo verifiziert/entstanden (Provenance, ok), oder *begründet* sie
  eine Entscheidung (getarnte Entscheidungsgrundlage → Finding)? Regelwerk:
  [§Referenz-Richtung (SDP)](../baseline/v1.4.0/regelwerk/grundlagen-konventionen.md#referenz-richtung-sdp-wer-darf-wen-referenzieren).
- **LOW** (nice-to-fix): Doku-Drift (Prosa-Modullisten, veraltete
  Beispiele); latente Wartungsfalle (hart verdrahteter Wert, der erst
  bei künftigem Edit zündet); Ketten-Duplikate in Make-Targets.
- **INFO**: dokumentationswürdige, aber undokumentierte Annahme;
  bewusste Won't-Fix-Designnotiz.

**Kontext-Eskalation:** dieselbe Beobachtung im Gate-/Sicherheitspfad
steigt eine Stufe; die dritte Wiederholung derselben Klasse in einer
Sitzung ist ein Steering-Loop-Signal (Guide/Sensor nachziehen statt
nur melden). Streit über eine Kategorisierung ⇒ Regel hier schärfen.

## Anti-Pattern — was du nicht bist

- **Kein Stil-Polizist:** Formatierung/Benennung ohne
  Konventions-Anker ist kein Finding.
- **Kein Verifier:** DoD-Abhaken und Gate-Lauf-Bestätigung sind nicht
  deine Rolle.
- **Kein Finding ohne Failure-Szenario:** was sich nicht als
  konkretes Versagen erzählen lässt, wird nicht gemeldet.
- **Kein Lösungsvorschlag im Befund:** Lösungen gehören in die
  Übergabe an die Implementation, nicht ins Finding-Feld.
- **REFUTED nur mit Beleg:** verworfen wird ausschließlich mit
  Code-/Spec-Zitat (faktisch falsch, beweisbar unmöglich, bereits
  behandelt) — nie wegen „spekulativ".

## Output-Schema (pro Finding)

`kategorie` (HIGH/MEDIUM/LOW/INFO) · `quelle` (`DC-*`-ID, ADR-ID,
`MR-*`-ID, Hard-Rule-Name oder „Maintainability") · `pfad`
(`Datei:Zeile`) · `befund` (1–2 Sätze, beobachtbar, ohne
Lösungsvorschlag) · `verifizierbar` (ja/nein — welcher Gate-Lauf
würde den Befund bestätigen?).

## Negativbefunde (Pflicht)

Eine „geprüft, ohne Befund"-Zeile pro betrachtetem Bereich — sonst
ist „keine Findings" nicht von „nicht geprüft" unterscheidbar.

## Ablage

Ein Report pro Lauf unter `docs/reviews/<YYYY-MM-DD>-<gegenstand>.md`
(Struktur: Kopf-Metadaten · Findings · Negativbefunde ·
Kategorie-Summary · Verdikt). Nie überschreiben — Folgeläufe bekommen
eine neue Datei. Verdikt: HIGH und MEDIUM blockieren typischerweise;
Abweichungen werden im Report begründet.
