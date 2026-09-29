# Verifikation slice-240 — Baseline-Pin-Hebung auf `v6.13.0` (DoD)

- **Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?"; geprüft gegen §2 DoD des Slice-Plans, nicht gegen Review-Entscheidungen.
- **Gegenstand:** `604bc25d..ae9f111f` — 4 Commits: `60aa9b18` (reiner Move MR-072 → `conventions/done/`, R100), `94b80f35` (MR-072 Link-Tiefen, ±9 Zeilen), `d8a7baf2` (Baum-Swap + MR-073), `ae9f111f` (Review-R1-Fixes + Review-Report).
- **Sensor-Evidence:** alle numerischen Aussagen unten selbst gemessen (`sha256sum -c`, `git grep`/`git diff`-Zählung, `make gates`, `make baseline-freshness`) — nicht aus dem Implementer-Bericht übernommen.
- **Modell-ID:** glm-5.3-flash (Claude Agent SDK)
- **Datum:** 2026-09-29

---

## DoD-Prüfung je Punkt (§2)

### 1. Vendierter Baum auf `v6.13.0`, `SHA256SUMS` verifiziert, `v6.9.0` entfernt, §Baseline-Pin nachgezogen — **ERFÜLLT**

- `.harness/baseline/` enthält nur noch `v6.13.0` (54 `SHA256SUMS`-Einträge; `sha256sum -c` im Baum: **RC=0**, 54× OK).
- `v6.9.0`-Baum im Range-Diff entfernt: 52 Dateien als Rename `R062`–`R100` erkannt, `modul-00`/`modul-04` als `D`/`A`-Paar (< 50 % Similarity) — Re-Vendor inkl. beider Bäume belegt.
- `harness/conventions.md:27-29`: §Baseline-Stand auf Release-URL `v6.13.0` retargetet, gepinnt mit `MR-073`; `Datum der Adoption: 2026-06-10` korrekt unverändert (Adoptionsdatum ≠ Pin-Stand).
- Neue MR-Datei `harness/conventions/MR-073-baseline-v6130.md` vorhanden; MR-Vorgänger als reiner Move nach `conventions/done/` (`60aa9b18`, R100, 0 Zeilen) mit Link-Tiefen-Korrektur als **eigenem** Folgacommit (`94b80f35`) — `AGENTS.md` §3.3 sauber zerlegt.

### 2. Alle lebenden Referenz-Klassen nennen `v6.13.0` — **ERFÜLLT** (Substanz gemessen; ein Zahlen-Befund am MR-073-Text, siehe V-1)

- **Vor-Swap-Messung** (HEAD `94b80f35`, außerhalb des vendierten Baums): `v6.9.0` in **75 Dateien mit 161 Vorkommen** — exakt die in MR-073 genannten Zahlen.
- **Ist-Bestand nach dem Swap:** die 27 „ganz frozen"-Träger der MR-073-Liste stimmen **Datei für Datei** mit den gemessenen verbliebenen Nennungen überein (13 `done/`-Slices, 2 Review-Reports, 3 aufgelöste MR-Dateien, CO-002, CR 2026-09-17, 1 Evidence-Datei, ADR-0085, `spec/lastenheft.md`, 3 Code-Dateien unter `internal/hexagon/core/`, `.d-check.closure.yml`); Rest-Mentionen in slice-240-Plan, MR-072/MR-073 und R1-Report sind Gegenstands-/Vergangenheits-Nennungen. **Keine lebende Pin-Referenz auf `v6.9.0` bleibt.**
- Mischfälle zeilenweise getrennt belegt (Stichproben): `spec/spezifikation.md:9` lebt (→ `v6.13.0`), Zeile 3540 Historie (bleibt `v6.9.0`); MR-056 Baseline-Link retargetet, Wortlaut-Vermerk Prosa bleibt.
- **Symlinks (`MR-055`):** 8 gemessen statt 7 geplant; alle 8 lösen gegen `v6.13.0` auf (kein toter Alias). Abweichung in MR-073 dokumentiert („gemessen; der Plan sagte 7").
- **`ignore-refs` (`MR-069`):** genau **4** neue Einträge (`.d-check.yml:382-388`: ADR-0085, MR-059, MR-061, MR-062) — deckungsgleich mit den 4 gemessenen Markdown-Links auf den entfernten Baum; Begründung als Kommentar im Tombstone.
- **Zitat-Delta (`MR-039`):** MR-056s Zitatzeile „…die Bedingung dafür, dass die Datei überhaupt nach `done/` darf…" ist in `v6.9.0`:36 und `v6.13.0`:36 byte-identisch — „kein Zitat-Delta" verifiziert.

### 3. `make gates` grün; `make baseline-freshness` Exit 0 — **ERFÜLLT** (beide selbst gefahren, HEAD `ae9f111f`, sauberer Baum)

- `make gates`: **Exit 0**, alle zehn Glieder grün; `doc-check`/`planning-check`: „886 Datei(en) geprüft, 0 Befund(e)"; `coverage-gate`: „Coverage 94.60% erfüllt Schwelle 93%"; `semgrep`: 55 Regeln / 0 Findings; Abschlusszeile „[gates] … green".
- `make baseline-freshness`: **Exit 0** — „Pin v6.13.0 ist der neueste Release-Tag" (Currency) und „gepinnter Tag v6.13.0 upstream unverändert (Bytes == vendored SHA256SUMS)" (Content).

## Abgrenzung (§1) — gehalten

- **Keine Adoption der inhaltlichen Deltas:** einzige `.d-check.yml`-Änderung ist das `ignore-refs`-Ventil; MR-073 führt die vier inhaltlichen Schlagzeilen des `v6.13.0`-Deltas ausdrücklich als nicht übernommen.
- **Frozen-Dateien byte-stabil:** kein `done/`-Slice, kein Alt-Review-Report, kein Carveout, kein CR, kein `spec/lastenheft.md`, kein ADR, keine Code-Datei im Range-Diff (nur der R100-Move + Link-Tiefen auf MR-072 als geplanter Vorgang).
- **Kein Release/Tag/GHCR-Lauf:** keine Änderungen an `Dockerfile` oder `.github/workflows/` (grep-verified: keine Treffer in der Range).

## Verifier-Befunde

### V-1 — LOW: MR-073s korrigierte `d-check:cite`-Zahlen weichen vom Diff ab (F-1-Nachklapp)

- **Pfad:** `harness/conventions/MR-073-baseline-v6130.md:72-80`; derselbe falsche Wert steht in der Commit-Botschaft `ae9f111f` („15 d-check:cite-Direktiven").
- **Messung** (okkurrenzgenau): vor dem Swap **38** `d-check:cite`-Direktiven auf `v6.9.0`, danach **24** (frozen bleibend) + **14** auf `v6.13.0`; Diff: 14 hinzugefügte / 14 entfernte Direktiven.
  - **Re-anchored (Spanne geändert): 7** — die 6 in MR-073 gelisteten (MR-005 `grundlagen-durchsetzungsschicht.md` 50→66; MR-031 `modul-09-implementierung.md` 196→203; MR-049 `modul-05` 170-171→180-181 und 160→170; slice-240-Plan 363-364→373-374 und 369→379) **plus MR-043** (`grundlagen-durchsetzungsschicht.md` 100-102→116-118) — genau die Spanne, die R1-F-1 benannt hatte.
  - **Tag-only (Spanne identisch): 7** (MR-031 `modul-09:20` und `grundlagen-harness-dateien:194`; MR-039 `modul-02:283-284`; MR-035 `grundlagen-begriffe:49`; Reviewer-Skill `modul-10:82`; Closure-Note-Reviewer-Skill `modul-11:83` und `template:83`).
- **Urteil:** MR-073 behauptet „**15** Direktiven … **6** neu geankert … **9** tag-only" — gemessen **14 / 7 / 7**. Die Reviewer-Zerlegung (9 re-anchored + 5 tag-only = **14**) traf die Gesamtzahl, zählte aber die beiden Plan-Spannen doppelt (zusätzlich unter `.harness/skills/reviewer.md` gelistet; der Skill trägt nur `modul-10:82-82`). **Beide Zerlegungen sind falsch; die Gesamtzahl des Reviewers (14) stimmt.** Die Substanz der F-1-Auflösung (Abwesenheits-Behauptung entfernt, Neu-Ankern als MR-051-konform benannt) ist unverändert gültig — betroffen ist nur die Messzahl im Konventionsspeicher, dieselbe record-claim-vs-diff-Klasse wie F-1 selbst (zweite Instanz am selben Vorgang, `BEO-ALL/pin-bump-mirrors-ungated`).
- **Auflösung vor Closure-Move:** eine Zeile in MR-073 (15→14, 6→7, 9→7, MR-043-Spanne ergänzen); die Closure-Notiz darf „15" nicht wiederholen.

### V-2 — Beobachtung (kein Befund): Plan-Kopf „`spec/` bleibt unverändert" vs. zwei Retarget-Zeilen

`spec/architecture.md` und `spec/spezifikation.md` wurden je um genau eine Zeile retargetet (§Rolle-Pointer, Anker identisch). Das war **DoD-2-erzwungen** (entfernter Baum ⇒ sonst `target-missing` im lebenden Doc) und ist in MR-073 als Mischfundstelle dokumentiert; gegen den wörtlichen Plan-Kopf fällt es dennoch auf. Rückführung nach `next` (§4) nicht nötig: die Klassen waren je Fundstelle trennbar.

### V-3 — Beobachtung (kein Befund): `.d-check.closure.yml` in Plan §3 als Update genannt, im Vollzug frozen

Die 2 Provenanz-Kommentare sind Vergangenheits-Aussagen; die Ausführung hat sie zu Recht nicht angehoben und MR-073 führt die Datei in der Frozen-Liste. Plan §3 war hier breiter als der Vollzug — dokumentiert, kein Schaden.

## Verdict

**DoD je Punkt erfüllt (3/3), Abgrenzung gehalten.** Der mechanische Vorgang (Swap, Pin, Retargets, Ventil, Move-Zerlegung) ist messrichtig und gate-grün — alle Kern-Zahlen des Implementers (75/161, 28 Deltas, 8 Symlinks, 4 ignore-refs, 36/10/27) haben sich bei eigener Messung bestätigt. Offen ist **V-1 (LOW)**: die in `ae9f111f` „gemessen" korrigierten `d-check:cite`-Zahlen in MR-073 sind selbst nicht diff-gemessen (15/6/9 statt 14/7/7). Empfehlung: Zeilenkorrektur in MR-073 **vor** dem Closure-Move, danach kann der Slice in den Closure-Vorgang (`git mv` + §7 Notiz + `make fullbuild` inkl. `verify-closure-notes`) übergehen. §7 ist leer und slice-240 liegt in `in-progress/` — Soll-Zustand vor Verifier-Freigabe.
