# Review-Report — slice-053 (Modul `vcs`, git-Diff-Immutabilität), doc-first-Stratum

| Feld | Wert |
|---|---|
| **Gegenstand** | slice-053 doc-first (kein Code): `spec/lastenheft.md` (DC-FA-VCS-001), `spec/spezifikation.md` (§DC-FA-VCS-001.a + 5 `vcs.*`-Keys + `core-drift-vcs`), `spec/architecture.md` (VCS-Port-Rolle), `docs/plan/adr/0024-vcs-immutable-gate.md` (NEU), `docs/plan/adr/README.md`, `docs/plan/adr/0016-adr-immutable-gate.md` (Geschichte-Append), `docs/plan/planning/in-progress/slice-053-vcs-modul.md` (NEU), `…/roadmap.md` |
| **Datum** | 2026-06-29 |
| **Reviewer** | unabhängig (Subagent) |
| **Diff-Umfang** | 6 M + 2 ?? (git status); Code existiert noch nicht — nur Dokumentation auf Regelkonformität/Konsistenz/Design-Stimmigkeit geprüft |
| **Rolle** | Reviewer, **nicht** Verifier (kein DoD-Abhaken, kein Gate-Lauf-Bestätigen) |

## Findings

### F-1 · MEDIUM · DC-QA-04 (Alt-Tool-Parität) / ADR-0024
- **pfad:** `docs/plan/adr/0024-vcs-immutable-gate.md:155-157` (Fitness Function) · `docs/plan/planning/in-progress/slice-053-vcs-modul.md:96-98` (DoD-Testliste)
- **befund:** Das abgelöste `tools/adr-immutable-check.sh` führt 7 Selbsttest-Klassen; die Paritäts-Testaufzählung im ADR und in der DoD nennt nur 6 — die Klasse „BASE `Proposed` + Körper-Edit feuert **nicht**" (Skript `self_test` (5), die den `immutable-when`-Filter negativ absichert) fehlt. Die `Proposed → Accepted`-Reifung (Status-Richtung) ist genannt, der gefährlichere Vektor (Körper-Edit an einer `Proposed`-BASE darf nicht als `core-drift-vcs` feuern) ist als Test nicht aufgezählt; das Verhalten steht zwar in `spec/spezifikation.md` §DC-FA-VCS-001.a Schritt 3.
- **verifizierbar:** ja — `make test` (ein fehlender Negativtest lässt einen `immutable-when`-Filter-Regress unentdeckt durch).

### F-2 · LOW · DC-QA-03
- **pfad:** `spec/lastenheft.md` (DC-QA-03, Messmethode — im Diff **nicht** angefasst)
- **befund:** DC-QA-03 Messmethode lautet „… **alle Module außer `external` aktiv**". `vcs` ist das erste Modul, das ohne `--range`/`--staged` **fail-closed** (Exit 2) statt nur inert zu sein (anders als `pins`/`immutable`/`diagrams`); ein wörtlich „alle Module außer `external`"-Lauf bräche damit ab. Der doc-first-Change gleicht die QA-03-Messmethoden-Prosa nicht nach (der reale Gate-Lauf nutzt die explizite `modules`-Liste in `.d-check.yml`, daher kein realer Bruch — reine Prosa-Drift).
- **verifizierbar:** ja — `make doc-check` (realer Lauf bleibt grün); ein literaler „alle-außer-external"-Lauf würde Exit 2 zeigen.

### F-3 · LOW · Anforderungs-Anlege-Prozess (Akzeptanzkriterien-Trio)
- **pfad:** `spec/lastenheft.md` DC-FA-VCS-001, Akzeptanzkriterien (Happy/Boundary/Negative)
- **befund:** Der durchgängig betonte Sicherheits-Vertrag „fail-closed → Exit 2 bei fehlendem `.git`/fehlender Range" hat kein eigenes Kriterium im verbindlichen Trio; der „Boundary"-Slot belegt stattdessen den Modul-aus-Pfad (byte-identisch). Der fail-closed-Pfad ist nur in der Anforderungs-Prosa, in spezifikation §DC-FA-VCS-001.a Schritt 1, der ADR-Fitness-Function und der DoD-Testliste belegt — ein Verifier, der nur das Lastenheft-Trio abprüft, testet ihn nicht.
- **verifizierbar:** ja — `make test` (DoD-Test existiert; Lücke ist die fehlende verbindliche Kriteriums-Verankerung, kein Verhaltens-Loch).

### F-4 · LOW · ADR-0016 / AGENTS §3.5 (Sequenzierung)
- **pfad:** `docs/plan/adr/0016-adr-immutable-gate.md:121` (Geschichte-Append) · `docs/plan/adr/README.md` (Status-Zelle ADR-0016)
- **befund:** Geschichte-Append und Index-Status-Zelle erklären den Teil-Supersede als vollzogen („… teil-superseded durch ADR-0024: das Gate `adr-check` **läuft** auf das Modul `vcs` **um**"), während ADR-0024 `Proposed` ist und der Gate-Umbau in der DoD offen (`[ ]`) ist. Stillt der Slice nach doc-first, trägt eine **append-only**, immutable ADR dauerhaft eine Geschichte-Zeile über einen nicht eingetretenen Effekt; zugleich liest die Index-Status-Zelle „teil-superseded", obwohl der Wächter noch das Skript fährt.
- **verifizierbar:** nein — kein Gate fängt Tempus/Finalität; `adr-check` akzeptiert den Geschichte-Append unabhängig (Kern unverändert). Semantische Konsistenz-Beobachtung (Querbezug DoD §3 „Gate-Umbau `[ ]`").

## Negativbefunde (geprüft, ohne Befund)

- **Referenz-Richtung / MR-006:** Die bindenden Spec-Körper (lastenheft DC-FA-VCS-001, spezifikation §.a, architecture) verweisen nur auf gleich-/höherrangige `DC-*`-IDs + interne §; **kein** Abwärts-Verweis auf ADR/Slice/Commit/Planning im bindenden Text. Die `slice-053`-Token in den §7-Historien-Zeilen liegen in `exclude-sections` (Historie/Geschichte) — etablierte Provenance-Form. Geprüft, ohne Befund.
- **ADR-0016-Immutabilität (was angefasst):** Nur Geschichte-Append (Z. 121) + Index-Annotation; Körper, `## Geschichte`-Kern und Kopf-`**Status:**`-Zeile unverändert → adr-check-core-compare hält. (Sequenzierung separat als F-4.) Geprüft, ohne Befund.
- **Teil-Supersede-Ehrlichkeit:** Nur die **Skript-Mechanik** abgelöst, Policy + Gate `adr-check` + zwei Quadranten + Grandfathering bleiben — deckungsgleich mit ADR-0016s 4-teiliger Entscheidung; Teil-Supersede hat Präzedenz (Index: „ADR-0014 ratifiziert ADR-0002 §4"). Geprüft, ohne Befund.
- **„Skript bleibt pfad-stabil"-Begründung:** stichhaltig, **nicht** vorgeschoben — `tools/adr-immutable-check.sh` steht in ADR-0016 Inline-Code (Z. 40); `codepaths.roots` enthält `tools`; eine Löschung bräche `make doc-check` (`codepath-missing`) an einer immutablen Datei. REFUTED als Pretext.
- **QA-02/QA-03-Behauptung („relativiert weder"):** konsistent mit DC-QA-03-Wortlaut (kein Schreiben, kein Netz — `vcs` tut beides nicht) und DC-QA-02 (Eingabe = Repo-Stand inkl. `.git` + Optionen, reproduzierbar). Die Abweichung von ADR-0008s Prognose („würde die Determinismus-Zusage relativieren") ist **offengelegt** (ADR-0024: „Anders als ADR-0008/ADR-0023 es vorzeichneten") und **stichhaltig**: ADR-0008s Sorge galt der Ähnlichkeits-Rename-Erkennung, die `vcs` ausdrücklich Out-of-Scope stellt. Geprüft, ohne Befund.
- **Strata-Konsistenz:** §DC-FA-VCS-001.a deckt das Akzeptanz-Trio (Happy/Boundary/Negative); die 5 `vcs.*`-Keys decken alle in Anforderung/Algorithmus genannten Felder; ADR-0024 `Schärft:`-Ziele (§.a, architecture §1–§2) existieren; DoD §3 beschreibt die real autorisierten Artefakte. Geprüft, ohne Befund (modulo F-3).
- **Skript-Parität (Mechanik):** Kopf-only-Status-Strip, `Geschichte`-Strip, A=frei, D/R=Befund, `head-allow`, `..`-Range-Pflicht-Guard, unauflösbare-Basis→Exit 2, `--staged`-erster-Commit=No-op — alle in spezifikation §.a Schritt 1–5 erfasst (modulo F-1-Testaufzählung). Geprüft, ohne Befund.
- **Index-README-Ausschluss-Parität:** `docs/plan/adr/README.md` trägt keine `**Status:** Accepted`-Kopfzeile → `vcs.immutable-when` nimmt ihn automatisch aus; ADR-0016s „Der Index wird nicht geprüft" bleibt unabhängig vom `vcs.paths`-Glob erhalten. REFUTED als Lücke.
- **Akzeptanz-Trio + Out-of-Scope:** Trio vorhanden; Out-of-Scope umfassend (Pin-Form/IMM-001, git-Binary, Forge-/Netz-API, Ähnlichkeits-Rename, Schreiben/Neu-Pinnen, mehrere Hash-Algos). Geprüft (Kriteriums-Lücke separat als F-3).
- **Architektur (Diagramm/Sequenz/Fehlermodell):** VCS-Port (`FileAtRef`/`ChangedFiles`), `opt vcs aktiv`-Block, Exit-2-Zeile konsistent; sprach-frei respektiert („rein in der Implementierungssprache" statt „Go"). Geprüft, ohne Befund.
- **Roadmap/planning-check:** `slice-053` in `in-progress/` + Roadmap „Aktive Welle: welle-42-vcs" konsistent (Invariant erfüllt). Geprüft, ohne Befund.
- **Maintainability/Won't-Fix:** reine-Go-git-Dependency-Gewicht (`go.mod`/Image/`semgrep`-Scope, Pin-Disziplin), „`adr-check` braucht jetzt das Image", Bootstrap-Reihenfolge — in ADR-0024 Konsequenzen + slice §4 Risiken dokumentiert. Geprüft, ohne Befund.
- **Anker/Links neuer Sektionen:** gate-verifizierbar via `make doc-check` (anchors/links); kein manueller Slug-Mismatch entdeckt. Geprüft, ohne Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 (F-1) |
| LOW | 3 (F-2, F-3, F-4) |
| INFO | 0 |

## Verdikt

**Mergebar nach Klärung von F-1.** Der doc-first-Stratum ist intern konsistent
und regelkonform: MR-006-Referenz-Richtung gewahrt, ADR-0016 nur per
Geschichte-Append + Index-Annotation angefasst (Kern immutabel), Teil-Supersede
ehrlich abgegrenzt und mit stichhaltiger Pfad-Stabilitäts-Begründung, die
QA-02/QA-03-„relativiert-weder"-Zusage sauber und ihre Abweichung von ADR-0008
offengelegt. **Kein HIGH.** Das eine MEDIUM (F-1) ist eine
Test-**Aufzählungs**-Lücke, kein Design-Defekt — das fehlende Negativ-Paritäts-
Kriterium (Skript-Selbsttest 5: Körper-Edit auf `Proposed`-BASE feuert nicht)
gehört vor dem Code in die Paritäts-Testliste, da der Slice „volle Parität"
verspricht (DC-QA-04-Klasse). Die drei LOW sind beratend (QA-03-Messmethoden-
Prosa, fail-closed-Kriteriums-Verankerung, ADR-0016-Sequenzierungs-Finalität).
Das Verdikt blockiert MEDIUM gemäß Reviewer-Konvention; die Abweichung
(„mergebar nach Klärung") ist begründet, weil das Verhalten in
§DC-FA-VCS-001.a Schritt 3 bereits spezifiziert ist und die Lücke ausschließlich
die Test-Enumeration betrifft.
