# Review-Report: slice-228 — 2026-09-17

**Review-Art:** Code — geprüft gegen Slice-Plan, `AGENTS.md` §3 (insbesondere
§3.1, §3.3, §3.7, §3.8) und die Präzedenz-Findings von slice-200.

**Gegenstand:** `git log 3adcf206..HEAD` (vier Commits:
`1e756019` Beanspruchung, `20290bdd` 28er-Archivierung, `14b1527c`
`.d-check.yml`-Bereinigung, `47324ab9` Closure-Content).

**Skill:** `.harness/skills/reviewer.md` @ `1b689b1f` (Version 1.16.0)

**Modell:** claude-sonnet-5 · **Datum:** 2026-09-17

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-228-verbliebenen-bestand-200-227-archivieren.md`
  (vollständig gelesen)
- Präzedenz-Findings slice-200 (F-1/F-2/F-4, via Commit-Botschaft `49665e6e`
  und `docs/plan/planning/done/wellenlos/slice-200-*`-Stub/Archiv)
- `AGENTS.md` §3.1, §3.3, §3.7, §3.8, §5
- `BEO-ALL/batch-slice-archival-zips-post-rewrite-content` (Register-Eintrag,
  2. Auftreten in diesem Slice)
- keine neuen ADRs (reine Bestandspflege, wie im Plan deklariert)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Acht neu eingefügte Kommentare im `ignore-refs`-Block erklären in Fließtext, *warum* ein Eintrag gegenstandslos wurde (Bezug auf `slice-228`, benannte Slice-Nummern `slice-217/218/220/221/222/223/224/225`, und ein Review-Befund-Marker „dieselbe Klasse wie slice-200 F-4"). Das ist Herkunfts-Prosa und ein Review-Befund-Marker — keine der fünf zulässigen Klassen (Zusage · Kopplung · Abgrenzung · Rang-Zeiger · Grenze); Herkunft ist nur als ein auflösbares `DC-*`/`ADR-*`/`MR-*`/`seit welle-<NN>`-Feld zulässig, nicht als Slice-Nummer-Prosa. Neuzugänge fallen unter den Anker, unabhängig vom Stil umliegender Alt-Kommentare. | `AGENTS.md` §3.7 | `.d-check.yml:242-244,259-260,265-266,283-286,299-301,306-307,311-313` | ja — kein Gate prüft das (Urteil, kein `grep`, wie der Skill-Anker selbst festhält); Verifikation ist eine manuelle Zeilen-Lektüre gegen §3.7 | Kommentar mit Review-Befund-Marker/Slice-Nummer statt einer der fünf Klassen |
| F-2 | MEDIUM | Slice-Plan-§2/§3, Commit-Botschaft `20290bdd` und die Closure-Notiz behaupten übereinstimmend „37 zugehörige Review-Reports" archiviert. Die tatsächliche Zahl der archivierten, vormals flachen Review-Reports (Gegenprobe: `git diff --diff-filter=D` gegen `docs/reviews/` außerhalb `archiv/`, gegengeprüft durch Auszählen der `docs/reviews/*.md`-Einträge in allen 28 `archiv.zip`-Dateien `slice-200`–`slice-227`) ist **36**, nicht 37. Die Botschaft behauptet damit eine Menge, die die eigene Messung nicht trägt (AGENTS.md §5: „wer N Formen geprüft hat, berichtet N"). | AGENTS.md §5 (`BEO-ALL/commit-message-overclaims-work`) | `docs/plan/planning/done/wellenlos/slice-228-verbliebenen-bestand-200-227-archivieren.md` §2/§3 (DoD-Punkt 1, Plan-Tabelle); Commit `20290bdd` | ja — `git diff --diff-filter=D --name-only -- docs/reviews/ \| grep -v archiv/ \| wc -l` liefert 36 | Zählung im Bericht weicht von der tatsächlichen archivierten Menge ab (Off-by-one) |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Vollständigkeit der Archivierung (alle 28 Slices `slice-200`–`slice-227`) | geprüft, ohne Befund — alle 28 Stub+Zip-Paare vorhanden, `docs/plan/planning/done/*.md` (flach) enthält nur `welle-*-results.md` (unberührter Altbestand), `docs/reviews/*.md` (flach) ist leer |
| Zehn entfernte `ignore-refs`-Einträge in `.d-check.yml` (Commit `14b1527c`) | geprüft, ohne Befund — alle zehn entfernten `in:`-Ziele existieren nachweislich nicht mehr am alten Pfad; Stichprobe von sieben verbliebenen Einträgen (Zeilen 78, 92, 94, 144, 172/174/180, 246, 267) zeigt: die weiterhin gültigen (Evidence-Dateien zu slice-217/222, die laut Kommentar bewusst nicht mitwandern) sind korrekt stehengeblieben |
| Vier manuell nachgezogene tote Verweise (Commit `20290bdd`: CO-002, zwei Beobachtungs-Register-Evidenzen, MR-072) | geprüft, ohne Befund — alle vier Ziele lösen nach der Retarget-Korrektur auf existierende Dateien unter `done/wellenlos/` auf |
| Weitere, nicht in `.d-check.yml` erfasste tote Verweise außerhalb der Modul-Scan-Achsen (§3.8) — gezielter Grep über `CHANGELOG.md`, `README*.md`, `spec/*.md`, `harness/README.md` nach Pfaden auf die archivierten Slices/Reviews | geprüft, ohne Befund — Treffer sind ausschließlich bloße Slice-Nummer-Erwähnungen (Prosa/Status-Provenance-Marker), keine Pfad-Links; kein toter Link außerhalb der bereits geprüften Fundstellen |
| Kommentar-Klassen (§3.7) der neu hinzugefügten `.d-check.yml`-Kommentare | siehe F-1 |
| Scope-Treue (§1-Abgrenzung: kein Anfassen von `tools/archive-wave/`, `docs/reviews/archiv/`) | geprüft, ohne Befund — `git diff --stat` gegen beide Pfade ist leer; die vier Dateien außerhalb der erwarteten Pfad-Muster (`CO-001`, `CO-002`, ein eingehender CR, `MR-070`) sind reine `RewriteRepo()`-Pfad-Retargets, kein inhaltlicher Eingriff |
| Determinismus/Order-Abhängigkeit (DC-QA-02, §6-Risiko) | geprüft, ohne Befund — das Risiko ist als „weiter offen" korrekt in `BEO-ALL/batch-slice-archival-zips-post-rewrite-content` (2. Auftreten, `evidence/slice-228.md`) dokumentiert; der Effekt bleibt wie beim ersten Auftreten inhaltlich harmlos (reiner Pfad-Nachzug), keine über das dokumentierte Maß hinausgehende Instanz gefunden |
| Vorherige Findings am selben Vorgang (slice-200 F-1/F-2/F-4) korrekt behandelt | geprüft, ohne Befund für F-2/F-4 (siehe oben); F-1 (Commit-Bündelungs-Deklaration) ist laut Auftrag durch MR-064 bereits generell geschlossen und hier nicht erneut zu prüfen |
| Vorgelagerte Sub-Area-Prüfungen (§8, `d-check:cite`-Belege) | geprüft, ohne Befund — beide `d-check:cite`-Direktiven (Zeilen 193, 206 des Slice-Plans) zitieren `modul-05-planning-harness.md` wortgleich an den referenzierten Zeilen |
| `make doc-check`-Gegenprobe | geprüft, ohne Befund — 791 Dateien, 0 Befunde auf dem aktuellen Stand, bestätigt die im Plan behauptete Grün-Meldung |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Kommentar mit Review-Befund-Marker/Slice-Nummer
statt einer der fünf Klassen · Zählung im Bericht weicht von der tatsächlich
archivierten Menge ab (Off-by-one)

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) verstößt gegen eine Hard Rule
(`AGENTS.md` §3.7) und ist vor Closure zu bereinigen (die acht betroffenen
Kommentare auf eine der fünf zulässigen Klassen zurückführen, ohne
Slice-Nummern/Review-Befund-Marker). F-2 (MEDIUM) verlangt eine Korrektur der
Zahlenangabe („37" → „36") in Slice-Plan §2/§3; die Commit-Botschaft
`20290bdd` selbst bleibt als Lauf-Beleg unverändert (§3.5-analoge
Unveränderlichkeit historischer Commits), der Fehler gehört aber in die
Closure-Notiz als „ging anders als geplant"-Eintrag benannt, nicht
stillschweigend übernommen.

**Übergabe:** Beide Findings gehen an den Implementer zurück (Rückkante
Review → Plan). Die DoD-Häkchen für „Alle 28 … archiviert" und „Closure-Notiz"
sind mit der unkorrigierten Zahlenbehauptung bereits als `[x]` gesetzt — das
gehört vor der eigentlichen Closure (`git mv` nach `done/`) noch einmal
angefasst. Dieser Report ersetzt keine Verifikation — DoD-/Spec-Konformität
prüft der Verifier separat.
