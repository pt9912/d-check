# Review R2 — slice-269: `releasing.md` zieht nach docs/maintainer/

**Review-Art:** Code (Diff gegen Plan, ADRs, Konventionen und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-269, Range `c81f3545..HEAD` — `1633d299` (Plan-Änderung nach R1),
`28552fcb` (Fix: MR-077, Rang 6, Tombstone-Kommentar); dazu die Form-Änderung am
R1-Report in `c81f3545`
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan slice-269 inkl. Plan-Änderung; R1-Report zu slice-269;
[ADR-0025](../plan/adr/0025-codepaths-ignore-refs.md) (Tombstone-Register);
`AGENTS.md` §2, §3.7, §5 (Regeln 13–16), §6 Schritt 4;
[`MR-021`](../../harness/conventions.md#mr-021), [`MR-051`](../../harness/conventions.md#mr-051);
Baseline `v6.17.0` · `templates/harness/conventions/MR-NNN-titel.template.md`,
`regelwerk/grundlagen-harness-dateien.md` §harness/conventions.md als Konventionsspeicher,
`regelwerk/grundlagen-source-precedence.md` §Source Precedence,
`regelwerk/modul-01-entwicklungszyklus.md` §Ziel-Form: Source-Precedence-Block.

## Messungen (Kommando und Ergebnis)

- **Abweichung oder Default?** Kanon §Source Precedence führt Rang 6 als
  `docs/user/*.md` mit „Quality-, Releasing- und Runbook-*Sichten*"; die Vorlagen
  `AGENTS.template.md` und `harness/README.template.md` ebenso. `modul-01` §Ziel-Form:
  „Die konkrete Rangordnung ist projektspezifisch […]; Wahl und Begründung gehören in den
  Adaptions-Block". Ein zweites Verzeichnis im Rang 6 ist damit eine Abweichung vom
  Default, kein Fork — ein `MR`-Eintrag ist die vorgesehene Form.
- **Form der Ersetzt-Baseline-Regel im Bestand:** `grep -n "Ersetzt-Baseline-Regel" -A2 harness/conventions/MR-0*.md`
  → jeder Eintrag, der eine **benannte** Regel ersetzt (MR-004, -005, -006, -007, -015,
  -021, -023, -031, -032, -033, -0098), trägt sie als Markdown-Link in
  `.harness/baseline/v6.17.0/regelwerk/…`; die übrigen schreiben „keine". MR-077 ersetzt
  eine benannte Regel und trägt sie als Inline-Code ohne Link.
- **Weitere lebende Spiegel von Rang 6:** `git grep` nach „Rang 6", nach Links mit Text `docs/user/` und nach `docs/user/*`,
  und `git grep -n -i "source.precedence"` außerhalb `done/`, `docs/reviews/`, `.harness/`
  → nur `AGENTS.md` §2, `harness/README.md` §Source precedence, der Index in
  `harness/conventions.md`, MR-077 selbst und der Plan. `README.md`/`README.de.md` nennen
  die Precedence nur als Zeiger auf `harness/README.md`; `CHANGELOG.md:2739` ist ein
  Lauf-Beleg. `harness/conventions/done/MR-009`/`MR-010` sind aufgelöst und eingefroren.
- **Neun Ränge:** beide Tabellen zählen 1–9 unverändert; nur Zeile 6 geändert, je eine Zeile
  im Diff (`git diff c81f3545..HEAD -- AGENTS.md harness/README.md`).
- **`d-check:cite` auf AGENTS.md:** `git grep -n "d-check:cite AGENTS.md"` außerhalb
  `done/`/Reviews → eine Spanne, `.harness/skills/reviewer.md` → `AGENTS.md:287-287`;
  `sed -n 287p AGENTS.md` → „Halluzinierte Gates sind die häufigste Form von Harness-Lüge".
  Die Änderung ersetzt eine Zeile durch eine; keine Verschiebung. Weder `AGENTS.md` noch
  `harness/README.md` tragen eigene `d-check:cite`-Direktiven.
- **Grenze des MR-Eintrags gegen die Konfiguration:** `grep -n "AGENTS.md\|harness/README.md" .d-check*.yml`
  → `structure`-Regeln lesen in `harness/README.md` nur §Sensors (Zellbreiten) und in
  beiden Dateien das Zeilenbudget; keine liest die Rangtabelle. Die Grenz-Aussage
  „Kein Gate liest die Rangtabelle gegen die Ablage" stimmt.
- **Plan-Messung „keine Scan-Menge hängt an `docs/user/` außer einer `structure`-Regel":**
  `grep -n "docs/user" .d-check*.yml Makefile` → außer den Tombstone-Zeilen nur
  `.d-check.yml:810` (`files: docs/user/benutzerhandbuch.md`). Bestätigt.
- **Tombstone-Grenze gegen den Gegenstand:** frischer Klon nach Scratchpad, in
  `docs/user/operations.md` eine Zeile mit den vier Inline-Code-Pfaden
  `docs/user/releasing.md`, `docs/user/maintainer`, `docs/user/maintainer/`,
  `docs/user/maintainer/releasing.md`, dann
  `docker run --rm --network none -v <klon>:/repo:ro -w /repo d-check:latest` →
  `1096 Datei(en) geprüft, 1 Befund(e)`:
  `operations.md:189  docs/user/maintainer/releasing.md  codepath-missing`.
  Genau das, was der neue Kommentar sagt.

## Findings

### LOW-1 — Die ersetzte Baseline-Regel steht ohne Link

- **kategorie:** LOW
- **quelle:** Baseline-Vorlage `MR-NNN-titel.template.md` (Regeln-Absatz: „als Link mit Abschnitts-Anker in die vendored Fassung"); [`MR-021`](../../harness/conventions.md#mr-021)
- **pfad:** `harness/conventions/MR-077-releasing-doku-unter-docs-maintainer.md` · „`grundlagen-source-precedence.md` §Source Precedence, Rang 6 —"
- **befund:** Das Feld nennt eine benannte Regel, aber als Inline-Code statt als Link in `.harness/baseline/v6.17.0/regelwerk/`; damit fallen die beiden Wächter weg, auf die die Vorlage die Bewegung dieses Feldes abwälzt (Existenzprüfung des Links bei Anker-Umbenennung, pin-gebundener Verweis beim nächsten Bump), und der Eintrag weicht vom Bestand ab, in dem jede ersetzte Regel verlinkt ist.
- **verifizierbar:** ja — mit Link würde `make doc-check` einen umbenannten Anker als `anchor-missing` melden; ohne Link bleibt er still
- **klasse:** ersetzt-baseline-regel-ohne-link

### INFO-1 — Auflösungs-Trigger und Grenze sprechen über verschiedene Mengen

- **kategorie:** INFO
- **quelle:** `AGENTS.md` §5 Regel 13
- **pfad:** `harness/conventions/MR-077-releasing-doku-unter-docs-maintainer.md` · „oder die Datei kehrt nach `docs/user/` zurück"
- **befund:** Die Grenze sagt, jede weitere Datei unter `docs/maintainer/` sei durch den Eintrag gerankt; der Trigger feuert aber schon, wenn **die eine** Datei zurückkehrt, auch wenn das Verzeichnis dann noch andere trägt. Und die erste Trigger-Hälfte („die Baseline nennt einen eigenen Ort") löst den Eintrag nur auf, wenn dieser Ort `docs/maintainer/` ist — sonst ist sie ein Re-Evaluierungs-Anlass, kein Auflösungs-Anlass. Heute ohne Folge (eine Datei im Verzeichnis).
- **verifizierbar:** nein (Urteil)
- **klasse:** trigger-menge-ungleich-geltungsbereich

### INFO-2 — Plan-Änderung als Prosa, Tabelle §3 und §1-Ziel unverändert

- **kategorie:** INFO
- **quelle:** `AGENTS.md` §6 Schritt 4
- **pfad:** slice-269 §3 · „*(Plan-Änderung vor dem nächsten Code-Commit, nach R1 MEDIUM-1:"
- **befund:** Die Änderung steht vor dem Code-Commit (`1633d299` vor `28552fcb`) und nennt die neuen Gegenstände, aber die Plan-Tabelle in §3 führt `AGENTS.md`, `harness/README.md`, `harness/conventions.md` und die neue MR-Datei nicht als Zeilen, und das Ziel in §1 beschreibt weiter zwei Commits. Der Nachsatz „`make mention-coverage` deckt alle Artefakte" stützt nichts über `docs/user/`: die Menge von `mentions` sind die ADR-Dateien gegen den ADR-Index (`.d-check.yml` `mentions:`).
- **verifizierbar:** nein (Urteil)
- **klasse:** plan-aenderung-nur-in-prosa

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| MR-077: Pflichtfelder | geprüft, ohne Befund außer LOW-1. Datum, Geltungsbereich, Ersetzt-Baseline-Regel, Adaption, Begründung, Auflösungs-Trigger vorhanden; `Löst auf`/`Ausgelöst durch` zu Recht weggelassen (keine Ablösung); `Grenze` als Zusatzfeld zulässig. |
| MR-077: Abweichung oder Baseline-konform | geprüft, ohne Befund. Abweichung vom Default-Rang 6, vom Kanon ausdrücklich an den Adaptions-Block verwiesen (`modul-01` §Ziel-Form); genau eine ersetzte Regel; Zitat des Kanons sinngemäß richtig („Betriebs-, Quality-, Releasing- und Runbook-Sichten"). |
| MR-077: Begründung (Frage 20) | geprüft, ohne Befund. Auftraggeber-Wunsch trifft am Gegenstand zu (Anwender- vs. Veröffentlicher-Doku). |
| Index-Zeile in `harness/conventions.md` | geprüft, ohne Befund. `<a id="mr-077">` vorhanden, Geltungsbereich und Ersetzt-Spalte deckungsgleich mit der Datei; beide Verweise in `AGENTS.md`/`harness/README.md` zeigen auf den Index-Anker, nicht auf die Datei. |
| Rangtabellen zueinander und zur Ablage | geprüft, ohne Befund. Beide nennen `docs/user/` und `docs/maintainer/` in Rang 6, neun Ränge unverändert, `ls docs/maintainer` → `releasing.md`; die Formulierungen unterscheiden sich (AGENTS ordnet je Verzeichnis zu, README fasst zusammen), widersprechen sich nicht. Der Link mit Text `docs/maintainer/` zielt auf `releasing.md` — in beiden Tabellen gleich. |
| Weitere lebende Spiegel der alten Rang-6-Aussage | geprüft, ohne Befund (siehe Messungen). |
| `d-check:cite`-Spannen und Zitate auf `AGENTS.md` | geprüft, ohne Befund. Die einzige Spanne (`AGENTS.md:287`) zeigt weiter auf die zitierte Zeile. |
| Tombstone-Kommentar §3.7 (R1 LOW-1) | geprüft, ohne Befund. Keine Slice-Nummer mehr; Kommentar trägt Kopplung und Grenze. |
| Tombstone-Kommentar §5 Regel 13 (R1 LOW-2) | geprüft, ohne Befund. Gegen das gemessene Verhalten im Klon bestätigt: Ur-Pfad und Verzeichnis still, Datei-Pfad unter dem Zwischen-Ort gemeldet. |
| Plan-Zahl CHANGELOG (R1 INFO-1) | geprüft, ohne Befund. „CHANGELOG einer" — Aufzählung summiert jetzt auf 16. |
| R1-Report, Form-Änderung in `c81f3545` | geprüft mit Grenze. Der Report kam erst mit diesem Commit ins Repo; der Vorzustand der umformulierten Zeile ist in git nicht vorhanden, „Inhalt unverändert" ist daher nicht gegen einen Diff prüfbar. Geprüft wurde der Inhalt gegen den Gegenstand: die Bruchprobe zum Tombstone ist reproduziert (siehe Messungen), die Quellen-Angaben von MEDIUM-1 decken sich mit Kanon und Vorlagen. |

## Kategorie-Summary

HIGH 0 · MEDIUM 0 · LOW 1 · INFO 2

## Verdikt

R1 MEDIUM-1, LOW-1 und LOW-2 sind eingelöst. Kein blockierendes Finding. LOW-1 annehmen
oder begründen; INFO-1 und INFO-2 ohne Handlungsbedarf vor der Closure.
