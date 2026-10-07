# Review slice-253 — Baseline-Pin-Hebung `v6.13.0` → `v6.17.0` (R1)

- **Review-Art:** Code/Doku-Review — geprüft gegen den Slice-Plan `slice-253` (§1 Ziel und Abgrenzung, §3 Plan, §6 Risiken), die Hard Rules (`AGENTS.md` §3, §5 Regel 15) und die Konventionen (`MR-021`, `MR-039`, `MR-051`, `MR-055`, `MR-069`, `MR-070`); **nicht** gegen die DoD (Verifikation, getrennter Kontext).
- **Gegenstand:** `0c0d46dc..957aeedc` — 3 Commits: `f8382b60` (reiner Move MR-073 → `conventions/done/`, 100 % Similarity, mit `--no-verify`), `c7ace79b` (MR-073 Link-Tiefen ±5), `957aeedc` (Baum-Swap 55/55 Dateien, §Baseline-Pin, MR-074, MR-Index, 8 Symlinks, Retargets in den lebenden Dokumenten, 4 Cite-Neu-Ankerungen).
- **Skill:** `reviewer.md` @ 1.16.0
- **Modell-ID:** claude-opus-5-5 (Claude Agent SDK)
- **Datum:** 2026-10-07
- **Eingangs-Kontext:** Slice-Plan `slice-253`; `AGENTS.md` §3/§5; `harness/conventions.md` samt `MR-021`, `MR-039`, `MR-051`, `MR-055`, `MR-069`, `MR-070`; Vorgänger `MR-073` und `slice-240` samt den Reviews vom 2026-09-29 (dort F-1 „record-claim-vs-diff" und F-2 „frozen-list-misclassification" am selben Vorgangstyp); Vorher-Baum aus `git archive 0c0d46dc .harness/baseline/v6.13.0`, Nachher-Baum `.harness/baseline/v6.17.0/`; `v6.17.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.
- **Gelaufene Sensoren (Reviewer-seitig):** `make baseline-verify` → `verify ok (54 Dateien, vollständig)`; `make doc-check` → `951 Datei(en) geprüft, 0 Befund(e)`. `make gates` nicht gefahren (Gate-Lauf-Bestätigung ist Verifier-Rolle).

---

## Findings

### F-1 — MEDIUM

- **Kategorie:** MEDIUM
- **Quelle:** Anker „Quelle über ihren Geltungsbereich hinaus zitiert" / „Botschaft verallgemeinert über die Messung" (`AGENTS.md` §5 Regel 15/16); `v6.17.0` · `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt (Absatz *Ein Index, mehrere Eigentümer*)
- **Pfad:** `harness/conventions/MR-074-baseline-v6170.md:31` · „dieses Repo hat keine Werkzeug-Fragmente"
- **Befund:** Die Abgrenzung der ersten Schlagzeile stützt sich auf eine falsche Prämisse: das `Makefile` bindet per `include a-check.mk` ein Make-Fragment ein, das laut eigenem Kopf „aus `a-check --print-mk`" erzeugt ist (Target `a-check:`), und `.d-check.yml` `targets.makefiles: [Makefile]` liest es nicht („gate-consistency parst nur das Makefile, keine includes"). Genau diese Konstellation — Werkzeug-Fragment, dessen Targets der Deklarations-Sensor nicht sieht — ist der Gegenstand des neuen `v6.17.0`-Abschnitts. **Failure-Szenario:** Ein Folge-Slice (oder der nächste Bump) liest MR-074 als Inventar und schneidet die Adoptions-Frage nicht, weil sie hier als „betrifft dieses Repo nicht" abgeschlossen dasteht.
- **Verifizierbar:** ja — `grep -n 'include a-check.mk' Makefile`; Kopf von `a-check.mk` Zeile 1–2; `.d-check.yml` `targets`-Block. Kein Gate fängt die Falschbehauptung.
- **Klasse:** record-claim-vs-diff (Abwesenheits-Behauptung gegen den Repo-Bestand)

### F-2 — MEDIUM

- **Kategorie:** MEDIUM
- **Quelle:** Anker „Botschaft verallgemeinert über die Messung hinaus" (`BEO-ALL/commit-message-overclaims-work`(b), `AGENTS.md` §5 Regel 15); Plan §1 (der Bump ist mechanisch, MR-074 „nennt die Schlagzeilen" — er ist nach dem Entfernen des alten Baums das einzige Delta-Inventar im Arbeitsbaum)
- **Pfad:** `harness/conventions/MR-074-baseline-v6170.md:26` · „Vier inhaltliche Schlagzeilen" und `:38` · „Viertens ziehen die Templates (… Welle-Ergebnis …) diese drei nach"
- **Befund:** Gemessen (`diff -I '<!-- Quelle:'` Vorher-/Nachher-Baum) tragen die 22 Delta-Dateien mehr als vier Linien, und die vierte ist falsch zugeordnet: (a) `regelwerk/modul-15-observability.md` ändert die Audit-Span-Regel (Pflicht-Feld ohne Wert bleibt Pflicht und wird als „nicht bekannt" gekennzeichnet) — in MR-074 nicht genannt; (b) `welle-results.template.md`, `slice.template.md`, `conventions.template.md` und teils `review-report.template.md` ziehen die Register-Kennung `BEO-<NNN>` → `BEO-<KUERZEL>/<slug>` bzw. `observations.md` → `observations/` nach — das ist keine der „drei" Schlagzeilen; `welle-results.template.md` trägt **ausschließlich** diesen Nachzug; (c) `spezifikation.template.md` schiebt die Historie von `## 7.` nach `## 8.` (siehe F-5). Zählung und Diff-Zeilenzahlen stimmen (22, 65/17/17/17/15 nachgemessen); überdehnt ist die Lesart als vollständige Inhaltsliste. **Failure-Szenario:** Eine Adoptions-Planung auf Basis von MR-074 übersieht (a) und (c), weil der vendorte Vorher-Baum nur noch in der Historie liegt. Gleiche Klasse wie F-1 des Vorgänger-Reviews (`slice-240` R1) — zweites Auftreten am selben Vorgangstyp.
- **Verifizierbar:** ja — `git diff 0c0d46dc 957aeedc -M -- .harness/baseline` bzw. `diff -rq -I '<!-- Quelle:'` gegen den archivierten Vorher-Baum.
- **Klasse:** record-claim-vs-diff (Delta-Inventar unvollständig/fehlzugeordnet)

### F-3 — LOW

- **Kategorie:** LOW
- **Quelle:** `MR-039` (Zitat-Delta im Bump-Eintrag) — hier am Rand seines Geltungsbereichs
- **Pfad:** `.d-check.yml:933` · „Der Gate-Index steht einmal ... Diese Datei fuehrt die Liste nicht."
- **Befund:** Der Kommentar begründet „dem EINEN Gate-Index" mit einem **elidierten** Zitat der „adoptierten Vorlage". Beide Zitat-Enden stehen in `v6.17.0` · `templates/AGENTS.template.md` §4 weiter wörtlich, aber genau die ausgelassene Mitte hat sich geändert: dort steht jetzt „Targets aus Werkzeug-Fragmenten stehen in dem Teil des Werkzeugs, den §Sensors verlinkt". MR-074 (`:76-79`) erklärt „Kein Zitat-Delta" anhand von zwei namentlich geprüften Zitaten; dieses steht nicht in der Prüfmenge. `MR-039` nennt Konfigurationsdateien in seinem Geltungsbereich nicht ausdrücklich — deshalb LOW statt MEDIUM. **Failure-Szenario:** Das Zitat altert still in der Richtung, die F-1 betrifft.
- **Verifizierbar:** ja — `sed -n 200,207p .harness/baseline/v6.17.0/templates/AGENTS.template.md` gegen `.d-check.yml:933-934`.
- **Klasse:** elided-quote-spans-changed-text

### F-4 — LOW

- **Kategorie:** LOW
- **Quelle:** `v6.17.0` · `templates/harness/conventions/MR-NNN-titel.template.md` (Regelzeile: „`Löst auf` und `Ausgelöst durch Baseline-Stand` nur, wenn dieser Eintrag einen früheren ablöst"; „kein Status-Feld"); Beobachtung `BEO-ALL/form-vom-nachbarn-statt-von-der-vorlage`
- **Pfad:** `harness/conventions/MR-074-baseline-v6170.md:3` · „**Status:** Accepted" (und fehlende Felder `Löst auf`/`Ausgelöst durch Baseline-Stand`)
- **Befund:** MR-074 löst MR-073 ab, trägt aber weder `Löst auf` noch `Ausgelöst durch Baseline-Stand`, dafür ein `Status`-Feld, das die Vorlage ausschließt. Die Form ist vom Nachbarn übernommen: MR-067 führte beide Ablöse-Felder noch, MR-071/072/073 nicht mehr — die Serie driftet seit drei Einträgen von der Vorlage weg, die in beiden Pins (`v6.13.0` und `v6.17.0`) dieselbe ist. Kein Gate prüft die Feldmenge; die Kette bleibt über §Aufgelöste Adaptionen auffindbar, deshalb LOW.
- **Verifizierbar:** ja — `grep -n 'Löst auf' harness/conventions/done/MR-0{67,71,72,73}*.md harness/conventions/MR-074*.md`.
- **Klasse:** form-vom-nachbarn (MR-Feldmenge)

### F-5 — INFO

- **Kategorie:** INFO
- **Quelle:** `MR-0098` (`matrix.exclude-sections` = `[Geschichte, "7. Historie"]`)
- **Pfad:** `.harness/baseline/v6.17.0/templates/spec/spezifikation.template.md:115` · „## 7. Festlegungen der Harness-Werkzeuge"
- **Befund:** Die neue Vorlage setzt §7 auf die Harness-Werkzeug-Festlegungen und die Historie auf `## 8.`. Unadoptiert bleibt das folgenlos (d-checks `spec/spezifikation.md` führt `## 7. Historie`); eine spätere Adoption der zweiten Schlagzeile verschiebt die Historie aus der Ausnahme von `MR-0098` heraus. Das wäre laut (matrix-Befunde), nicht still — deshalb nur als undokumentierte Kopplung notiert.
- **Verifizierbar:** ja — `.d-check.yml:614` gegen die Vorlage.
- **Klasse:** latent-section-number-coupling

### F-6 — INFO

- **Kategorie:** INFO
- **Quelle:** Hard Rule `AGENTS.md` §3.3 (MR-/Wellen-Lifecycle-Move, Regelfall 1); `make hooks`-Sensor-Grenze 1
- **Pfad:** Commit `f8382b60` · „reiner Move" (gesetzt mit `--no-verify`)
- **Befund:** Bewertung wie erbeten: §3.3 ist **eingehalten** — `f8382b60` ist ein reiner Rename (100 %, 0 Zeilen), `c7ace79b` trägt ausschließlich die fünf Tiefen-Korrekturen, und der rote Zwischenstand ist nicht die Spitze eines Push (`origin/main..HEAD` enthält alle drei, Spitze `957aeedc`, `doc-check` grün). Das `--no-verify` war allerdings nicht nötig: §3.3 beschreibt selbst den Weg, auf dem der Hook passiert — die Korrektur liegt im Arbeitsbaum, nur der Move ist gestagt. `--no-verify` überspringt zudem den `commit-msg`-Hook (`trace-check`); die Botschaft trägt `slice-253`, der Bypass hatte also keine Folge. Notiert, weil der Umweg über „Inhalt zurücksetzen + Hook umgehen" eine Lernlage ist, kein Verstoß.
- **Verifizierbar:** ja — `git diff f8382b60^ f8382b60 -M --summary`; `git log origin/main..HEAD`.
- **Klasse:** unnoetiger-hook-bypass

### F-7 — INFO

- **Kategorie:** INFO
- **Quelle:** `MR-070` / `BEO-ALL/pin-bump-mirrors-ungated` (Messung vor der Ersetzung)
- **Pfad:** `harness/conventions/MR-074-baseline-v6170.md:43-44` · „72 Dateien außerhalb des vendorten Baums nannten `v6.13.0` — 180 Vorkommen"
- **Befund:** Nachgemessen (`git grep -o 'v6\.13\.0' 0c0d46dc -- . ':!.harness'`) sind es 73 Dateien / 184 Vorkommen; die Differenz ist genau der Slice-Plan `slice-253` selbst (4 Vorkommen, davon 2 retargetete Cite-Zeilen). Der Ausschluss ist sachlich vertretbar, aber nicht deklariert. Symlinks zählt `git grep` nicht mit — die 8 Baseline-Symlinks (von 11 unter `.claude/rules/`) sind separat genannt und stimmen.
- **Verifizierbar:** ja — Kommando oben.
- **Klasse:** undeklarierter-mess-ausschluss

---

## Negativbefunde (geprüft, ohne Befund)

1. **Vendorter Baum, Delta-Zählung:** 55/55 Dateien, keine neu, keine entfallen; 22 inhaltlich geändert (8 Regelwerk, 14 Templates) + `SHA256SUMS` — exakt wie behauptet; Diff-Zeilen 65/17/17/17/15 nachgemessen und gleich. `make baseline-verify` grün (54 Dateien). Kein Heading im Regelwerk-Delta umbenannt (einziges Heading-Delta: `spezifikation.template.md`, F-5).
2. **Nichts versehentlich adoptiert:** Außerhalb `.harness/baseline/` enthält der Range-Diff nur Tag-Swaps, vier Cite-Zeilennummern, die MR-073-Tiefen und den neuen MR-074 (zeilenweise abgeglichen: Minus-Zeilen mit `v6.13.0→v6.17.0` ersetzt sind bis auf die Cite-Nummern identisch mit den Plus-Zeilen). Kein `harness/mk/`, keine Spec-§7-Festlegungen, keine Skill-Formänderung (`pfad` als Kurzzitat) übernommen.
3. **Lebende Retargets vollständig:** Kein lebendes Dokument trägt mehr einen `v6.13.0`-Pfad, -Link oder -URL; die verbleibenden 31 Dateien mit `v6.13.0` sind `done/`-Slices/-Welle, Reviews, `conventions/done/MR-073`, ADR-0100 (`Accepted`), der CR vom 2026-10-06, `CHANGELOG.md:103`, `.d-check.yml:369` (Tombstone-Kommentar der Vorgänger-Stufe), `tools/harness/selbstpruefung.sh:3` (Stand des Schwester-Werkzeugs, anderer Gegenstand), `observations/README.md:21` (Herkunft der Trigger-Audit-Adoption — korrekt stehen gelassen) und der Plan selbst (Vorher-Bezeichnung). Keine davon ist eine falsch entschiedene Mischfundstelle.
4. **Frozen byte-stabil:** `git diff 0c0d46dc 957aeedc` berührt kein `done/`, kein `docs/reviews/`, keine ADR, keinen CR, kein `CHANGELOG.md`, `.d-check.yml` nicht.
5. **Baseline-Form-Aussagen im neuen Baum:** `AGENTS.md` §6 „Baseline-Form `v6.17.0` §1 *Ziel und Abgrenzung*", `observations/README.md:18` und `MR-053` „§8 *Sub-Area-Prüfungen und Modus-Begründung*", `MR-049` „*Risiken und offene Punkte*" — gegen `v6.17.0` · `templates/docs/plan/planning/slice.template.md` (Headings §1, §6, §8) bestätigt; `roadmap.md` „§Roadmap-Struktur (v6.17.0)" gegen `modul-06-roadmap.md` (Abschnitt im Delta unberührt); Spec-Rolle-Zeilen gegen `modul-03-spec.md` §Ziel-Form: Spezifikation / §Ziel-Form: Architektur-Sicht (beide Headings vorhanden).
6. **Anker:** Alle 15 eindeutigen `#…`-Anker lebender Links in den `v6.17.0`-Baum lösen gegen die Headings auf (Handabgleich, zusätzlich `make doc-check` mit `anchors` grün).
7. **Symlinks (`MR-055`):** 8 Baseline-Symlinks zeigen auf `v6.17.0` und lösen auf; die 3 übrigen (`AGENTS.md`, `conventions.md`, `harness-README.md`) sind keine Pin-Träger und unverändert.
8. **`d-check:cite` (`MR-051`):** 14 lebende Direktiven im Vorzustand (bei `0c0d46dc` gezählt, ohne `citations.scope`-Ausnahmen), 4 Zeilennummern geändert (194→196, 49→50, 82→87, 83→84), 10 tag-only — wie behauptet; `citations` grün im `doc-check`.
9. **`ignore-refs` (`MR-069`):** wächst nicht — kein Markdown-Link auf `.harness/baseline/v6.13.0/` im Repo; die 20 Cite-Direktiven in `done/` liegen unter `citations.scope.ignore`. Behauptung bestätigt.
10. **Wörtliche Zitate (`MR-039`):** „Halluzinierte Gates sind die häufigste Form von Harness-Lüge", das MR-049-Kernsatz-Zitat und die beiden MR-069-Zitate stehen im `v6.17.0`-Wortlaut unverändert (Ausnahme: F-3).
11. **Index `harness/conventions.md`:** §Baseline zeigt auf MR-074; MR-074 aktiv, MR-073 unter Aufgelöste Adaptionen mit beiden Alt-Ankern und Nachfolger MR-074 — Ketten-Form wie die Vorgänger.
12. **Commit-Botschaften (`AGENTS.md` §5 Regel 15):** `f8382b60`/`c7ace79b` beschreiben exakt ihren Diff; `957aeedc` deckt sich mit dem Diff in Zahlen (54 verifizierte Dateien, 22/8/14, 8 Symlinks, 3 Rules-Dateien, 14/4/10) — die Überdehnung liegt im MR-Text (F-1, F-2), nicht in einer zusätzlichen Behauptung der Botschaft. „make gates gruen" und „baseline-freshness" nicht nachgefahren (Verifier).
13. **Kommentare (`AGENTS.md` §3.7):** Der Range fügt keinen Code-/Konfigurations-Kommentar hinzu.

---

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 2 | F-1, F-2 |
| LOW | 2 | F-3, F-4 |
| INFO | 3 | F-5, F-6, F-7 |

Wiederkehrende Klasse: **record-claim-vs-diff** (F-1, F-2) — dieselbe Klasse wie F-1 im `slice-240`-R1; zweites Auftreten am Vorgangstyp Pin-Bump, Zuordnung zu `BEO-ALL/pin-bump-mirrors-ungated` bzw. `BEO-ALL/commit-message-overclaims-work` bei der Closure.

## Verdikt

**Nicht freigegeben, bis F-1 und F-2 adressiert sind.** Der mechanische Teil des Bumps ist sauber: Baum, Retargets, Frozen-Klassen, Symlinks, Cite-Spannen und Anker halten nachgemessen. Beide MEDIUM betreffen den Text von MR-074 als Delta-Inventar — eine falsche Abwesenheits-Prämisse (F-1) und eine als vollständig lesbare, aber unvollständige Schlagzeilen-Liste (F-2). LOW/INFO: annehmen oder begründen.
