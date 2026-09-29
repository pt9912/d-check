# Review — slice-246 (`.d-check.yml`-Positionen aus ai-harness-init evaluieren)

**Review-Art:** Plan/Design — geprüft gegen Slice-Plan und dessen Messung (Reproduktion
der Probe), gegen die Kandidaten-Konfiguration des Snapshots `/tmp/aih-v6.13.0`, gegen
die eigene `.d-check.yml` und gegen die Hard Rules (`AGENTS.md` §3). Nicht geprüft: die
DoD-Abhakung als solche (Verifier).
**Gegenstand:** slice-246, Range `fc0d53ba..HEAD` (`e5787517` Claim/Move · `d21f7a76`
Evaluierung + slice-248); Welle: welle-91.
**Skill:** `reviewer.md` @ Version 1.16.0
**Modell-ID:** glm-5.3-flash
**Datum:** 2026-09-29
**Eingangs-Kontext:** `AGENTS.md` §3/§4/§6; `harness/conventions.md` (MR-000, MR-034,
MR-053, MR-054); `v6.13.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form,
`v6.13.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice,
`v6.13.0` · `regelwerk/grundlagen-traceability.md` §Ruheort-Regel; welle-91-Plan;
Snapshot `/tmp/aih-v6.13.0` (`.d-check.yml` der Schwester); eigene `.d-check.yml`
(green: `make doc-check` — 915 Dateien, 0 Befunde, vor und nach der Probe).

---

## Messungs-Reproduktion (Entscheidungsbasis für Positionen 2+3)

Die Probe wurde je zweimal reproduziert — Kandidaten-Klassen (`aussen` `paths: ["**"]`,
`adaptionsblock` mit `token: 'MR-\d{3}'`) samt Regeln (`spec-straten→aussen/sicht→aussen`,
`spec-straten→adaptionsblock/sicht→adaptionsblock`) in `.d-check.yml` eingesetzt, Lauf,
Revert:

| Probe | Ordnung | matrix-forbidden | matrix-inactive | Gesamt |
|---|---|---|---|---|
| A | `aussen` vor `adaptionsblock` (wie §7 beschrieben) | 29 | 11 | **40** |
| B | `adaptionsblock` vor `aussen` (Snapshot-Ordnung) | 29 | 11 | **40** |

Aufschlüsselung (A und B identisch in der Menge):

- **28 aussen-Links** aus den Straten bzw. der Sicht: 20 aus `spec/lastenheft.md`,
  7 aus `spec/spezifikation.md`, 1 aus `spec/architecture.md` — Ziele: 14 ×
  `harness/conventions.md`, 6 × `AGENTS.md`, 2 × Baseline-Zitat
  (`spec/spezifikation.md:9`, `spec/architecture.md:7`), je 1 `harness/README.md`,
  `packaging/dockerhub/README.md`, 2 × CR, 1 Carveout, 1 Register.
- **1 nacktes MR-Token:** `spec/lastenheft.md:4069` (`MR-004`, durch `d-check:ignore`
  für `ids` entschärft — `matrix` ehrt den Marker nicht, ganz wie die Schwester es
  dokumentiert).
- **11 matrix-inactive** (der Seiteneffekt: `paths: ["**"]` zieht jede Datei in die
  Status-Prüfung): `docs/plan/adr/README.md` ×7 (Superseded-Index-Zeilen),
  `CHANGELOG.md` ×3, `docs/user/releasing.md` ×1.

Probe A vs. B bestätigt die Aussage in §7, die Ordnung verändere die Zählung nicht
(nur die Labels: 14 Konventions-Links wechseln zu `adaptionsblock`) — der
First-Match-Unterschied ist label-seitig, nicht zählseitig.

---

## Findings

**F-1 · MEDIUM · Maintainability (Anker 8/17/18) ·
`docs/plan/planning/in-progress/slice-246-dcheck-yml-positionen-evaluieren.md:100-102`**
Die Closure-Notiz behauptet „29 Befunde (25 `aussen`-Links aus den Straten, 1 nacktes
MR-Token, 1 `matrix-inactive`-Seiteneffekt)“ — die Aufzählung summiert 27 ≠ 29, und die
getreue Reproduktion misst 40 Befunde: 28 aussen-Links (nicht 25) und **11**
matrix-inactive-Seiteneffekte (nicht 1), verteilt auf drei Dateien. Der gemeldete
Einzelfall (`docs/user/releasing.md`) ist der kleinste; die strukturell größte Klasse
(`docs/plan/adr/README.md` ×7 — der Index *muss* superseded ADRs führen; die Schwester
nimmt exakt ihn in `status.exempt-paths`) und `CHANGELOG.md` ×3 fehlen in der Notiz.
*verifizierbar:* ja — Probe-Lauf wie oben (A/B), `docker run` mit Kandidaten-Konfiguration;
Rückstand ist der fehlende Blick auf die `matrix-inactive`-Hälfte des Laufs.
*klasse:* `messung-unterbestreibt-seiteneffekte`

**F-2 · MEDIUM · Maintainability (Anker 18), AGENTS.md §3.6 ·
`docs/plan/planning/open/slice-248-matrix-aussen-adaptionsblock.md:51-55,66-69`**
slice-248 erbt die unvollständige Messbasis: §1 nennt als Status-Fall nur
`docs/user/releasing.md`, DoD 2 verlangt „jeder der in slice-246 gemessenen Befunde ist
entschieden“ auf Basis von „29“ — die zehn übrigen matrix-inactive-Funde sind nicht
geplant, obwohl jede nötige Ausnahme nach slice-248s eigener Abgrenzung ein §3.6-Fall
(Senkung nur per ADR) ist. Der Reconciliation-Aufwand (§8: „Gering“) steht damit auf der
Falsch-Underlage; der Rückführungs-Trigger (§4) wird dafür nicht ausgeschlossen.
*verifizierbar:* ja — derselbe Probe-Lauf; der Umfang wird sichtbar, sobald die Adoption
läuft und `make doc-check` ohne die Status-Ausnahmen rot bleibt. *klasse:*
`grenzen-liste-ohne-groesste-luecke`

**F-3 · LOW · Maintainability ·
`docs/plan/planning/in-progress/slice-246-dcheck-yml-positionen-evaluieren.md:59-61` vs. `:105-107`**
Widersprüchliche Lesart desselben Ausgangs: die DoD-2-Notiz sagt „die **Ablehnungen**
2+3“, §7 entscheidet für dieselben Positionen „Grundsatz **bejaht**“ mit Delegation an
slice-248. Der Plan (§1) kannte dafür nur die zwei Zustände adoptiert/abgelehnt; der
dritte Zustand ist in §7 ehrlich dokumentiert, friert aber nicht in die DoD-Zeile.
*verifizierbar:* nein (Urteil). *klasse:* `closure-widerspruch`

**F-4 · LOW · Maintainability ·
`docs/plan/planning/in-progress/slice-246-dcheck-yml-positionen-evaluieren.md:38-39`**
Position 4 deckt nur die ids-Hälfte der „Bereichs-Form der ADR-Kennung“ ab; die
adr-Klassen-Glob-Hälfte (`docs/plan/adr/[A-Z]*-[0-9]*.md`) und die Schwester-Abweichung
beim slice-Token (`token: 'slice-'` gegen d-checks `slice-\d{3}`) sind gegen den
Anspruch „die fünf Positionen“ nicht namentlich evaluiert. Beide Ablehnungen folgen aus
derselben Null-Treffer-Messung (segmentierte ADR-Form: 0 Treffer im Bestand, bestätigt)
bzw. aus d-checks schärferem Bestand — der Evaluiierungsanspruch ist nur ungenau.
*verifizierbar:* ja — Snapshot-/Template-Vergleich (`v6.13.0` · `templates/.d-check.yml`).
*klasse:* `abgrenzung-unvollstaendig`

**INFO-1 · Maintainability (Modul 5/6) ·
`docs/plan/planning/open/slice-248-matrix-aussen-adaptionsblock.md:8`**
slice-248 läuft „ohne Welle“, obwohl welle-91 §1 die Übernahme („an d-check angepasst“)
zum Welle-Ziel zählt — die Welle kann vor der eigentlichen Adoption schließen.
Wellenlos ist ein legitimer Modus; eine Zeile in welle-91 (Out-of-Scope oder Drift-Log)
macht die Trennung explizit. *verifizierbar:* nein. *klasse:*
`wellenzugehoerigkeit-unausgesprochen`

**INFO-2 · Maintainability · `.d-check.yml:591`**
Die Ablehnung von Position 5 („kein eigener Fall“) ist eine Zustandsaussage über
heutige Headings: trägt künftig ein Spec-Stratum ein `## Geschichte`, weitet die global
wirkende `exclude-sections: [Geschichte]` die Ausnahme still auf die Straten aus —
die von der Schwester dokumentierte Grenze, die dort nur der Reviewer fängt.
Dokumentationswürdige Annahme, kein Handlungsanlass. *verifizierbar:* nein.
*klasse:* `dokumentationswuerdige-annahme`

---

## Negativbefunde (Pflicht)

- **Position 1 (welle-Klasse):** geprüft — „bereits Bestand in schärferer Form“
  verifiziert: d-check führt `token: 'welle-\d{2,}'` mit **drei** Regeln
  (spec-straten/sicht/adr→welle, `.d-check.yml:515-527,558,567,572`), die Schwester ein
  Prefix-Token `welle-` mit zwei Regeln; MR-034-Bezug stimmt. Kein Befund.
- **Position 5 (exclude-sections-Scoping):** geprüft — Ablehnung sauber; kein Gate
  geändert oder gelockert, §3.5/§3.6 nicht berührt, kein Folge-ADR/CR geboten. Kein Befund.
- **First-Match-Ordnungs-Aussage:** durch Probe A/B gegeneninander geprüft — identische
  Summen, nur Label-Verschiebung; „die Zählung ist gültig“ bestätigt. Kein Befund.
- **Commit-Form:** `e5787517` reiner Move (R100, 0 Zeilen — §3.3); beide Botschaften
  tragen Kennungen (`make trace-check`-Form). Kein Befund.
- **MR-053/054-Form von slice-248:** drei Vorprüfungen vorhanden, `d-check:cite` auf den
  beiden kanonischen Blöcken; `Verantwortlich: —` korrekt für `open/`; je
  Abgrenzungspunkt eine Begründung. Kein Befund.
- **Probe-Hygiene:** `.d-check.yml` trägt die Probe nicht (Baum sauber;
  `make doc-check` 915 Dateien / 0 Befunde nach Revert). Kein Befund.
- **Risiko-Ausgänge §6 von slice-246:** beide mit Ausgang (eingetreten/entfallen),
  Form konform. Kein Befund.

## Hinweise (keine Findings)

- **Ruheort-Formen vor der Closure:** slice-246 §7 führt `../open/slice-248-…`,
  `../observations/` und den 4-stufigen `../../../../harness/…`-Prefix in
  Schreibort-Form. Vom Ruheort `done/welle-91/` (Geschwister-Form von slice-244/245)
  brechen sie; retargeten im Closure-Arbeitsbaum **vor** dem `git mv` — der
  pre-commit-`doc-check` fängt ein Vergessen (nicht still).
- Review-Report ist als Lauf-Beleg unverbindlich für Lebend-Verweise
  (`ignore-refs`-Klasse); slice-248-Zitate in §7 wandern wie gewohnt als Lauf-Beleg.

## Kategorie-Summary

**HIGH: 0 · MEDIUM: 2 · LOW: 2 · INFO: 2**

## Verdikt

Die fünf Entscheidungen sind in ihrer **Richtung** allesamt durch die Messung gedeckt
— die Delegation von Positionen 2+3 an slice-248 ist durch die vollständige gemessene
Menge (40, nicht 29) sogar *gestärkt*: eine Adoption im Slice wäre an den 11
Status-Seiteneffekten und 28 Link-Entscheidungen ohnehin gescheitert (DoD 1 erlaubt
Adoption nur bei grünen Gates). Die beiden MEDIUM blockieren die **Closure** von
slice-246 nicht der Sache, aber der Zahl nach: §7 friert eine Befund-Enumeration ein,
die intern nicht summiert und ihre größte Klasse verschweigt. Korrektur vor dem
`git mv`: Zahlen in §7 auf 40/28/1/11 heben (F-1), slice-248s §1/DoD 2 auf die
vollständige Status-Fall-Menge stellen (F-2). Danach ist der Slice closure-fähig;
DoD-Abhakung und `make fullbuild`-Belege bleiben Sache des Verifiers.
