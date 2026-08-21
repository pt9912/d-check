# Review-Report: slice-110 — Baseline-Bump v5.6.0 → v5.7.0

- **Review-Art:** Code-/Harness-Review (adversarial, unabhängiger Kontext) — geprüft gegen Slice-Plan, MR-Kette, Baseline v5.7.0, Hard Rules `AGENTS.md` §3
- **Gegenstand:** slice-110 (`docs/plan/planning/in-progress/slice-110-baseline-v570-bump.md`), welle-79-zwei-haelften-ein-waechter; Commit-Range `ed2a4af..78360eb` (vier Commits: `6baa083` Wellen-Eröffnung · `77d2496` Bundle vendored · `f845e8b` Pin/MR-028/Verweis-Hebung/Baum-Entfernung · `78360eb` Zwei-Hälften-Prosa + §9)
- **Skill:** `reviewer.md` @ 1.5.0 (Output-Schema sechs Felder, HIGH-/MEDIUM-Anker, Negativbefund-Pflicht)
- **Modell-ID:** claude-fable-5
- **Datum:** 2026-08-21
- **Eingangs-Kontext:** `AGENTS.md` §3 (Hard Rules), `harness/conventions.md` (MR-Index, §Baseline), `MR-021`/`MR-023`/`MR-026`/`MR-028`, vendorter Baum `.harness/baseline/v5.7.0/`, Gegenprobe am Kurs-Repo (`git diff v5.6.0..v5.7.0`), Wellendokument welle-79, slice-111. Hinweis zur Rollengrenze: die DoD-Punkt-für-Punkt-Prüfung ist laut Skill Verifikations-Territorium; sie wurde hier auf ausdrückliche Beauftragung mitgeleistet und ist unten als eigener Abschnitt ausgewiesen.

## Findings

**F-1**
- **kategorie:** MEDIUM
- **quelle:** `MR-021` (i. V. m. MR-028-Geltungsbereich „pin-gebundene Verweise")
- **pfad:** `harness/README.md:60` · `AGENTS.md:32` · `harness/conventions/MR-021-vendored-verweise-pin-gebunden.md:17`
- **befund:** Drei lebende Pin-Spiegel tragen nach der Hebung weiterhin v5.6.0: die `lab-regelwerk.zip`-Download-URL (`releases/download/v5.6.0/…`) in der `harness/README.md`-Tabellenzeile, deren Link-Ziel in **derselben Zeile** bereits auf `baseline/v5.7.0` zeigt, dieselbe URL in `AGENTS.md` §1, und die Angabe „aktuell `…/v5.6.0/…`" im MR-021-Körper. Die Commit-Behauptung von `f845e8b` („Verweis-Hebung … in allen lebenden Trägern") trifft nur auf das Grep-Muster `baseline/v5.6.0` zu. Exakt diese Klasse war im slice-106-Review Befund F-1/F-3 und wurde dort per Auflagen-Commit (`77b9c15`, „URL-Spiegel") geheilt — dies ist die Wiederholung eine Pin-Hebung später.
- **verifizierbar:** nein — kein Gate deckt URL-Versionen oder Ellipsen-Inline-Pins (`make doc-check` prüft externe URLs nicht; genau das hält der slice-106-Fix-Commit selbst fest).
- **klasse:** pin-spiegel-nicht-gehoben

**F-2**
- **kategorie:** LOW
- **quelle:** `MR-028` / Maintainability
- **pfad:** `harness/conventions/MR-028-baseline-v570.md:24-25` · `docs/plan/planning/in-progress/slice-110-baseline-v570-bump.md:22-24` (gleichlautend Roadmap-Drift-Log-Zeile 2026-08-21 und welle-79 §1)
- **befund:** Die Delta-Benennung „drei Regelwerks-Dateien, +5/−3 … vollständig gelesen" unterschlägt die Template-Hälfte des vendorten Bundles: in derselben Stufe änderten sich `templates/.harness/skills/reviewer.template.md` (+3) und `templates/docs/plan/planning/roadmap.template.md` (+15/−3) — das Bundle-Delta sind fünf Dateien, nicht drei; §9 auditiert und MR-028 benennt nur den Regelwerks-Baum. Der neue BEDIENHINWEIS des Roadmap-Templates trägt operative Sensor-Bauer-Regeln (Substring-/Code-Fence-Falle), die ungeprüft blieben; mildernd: beide Template-Änderungen spiegeln die zwei auditierten Regeln, und die Fence-Ausnahme steht bereits in slice-111 §2 („Fence-Inhalte zählen nicht").
- **verifizierbar:** nein — kein Gate misst Bundle-Delta ↔ Audit-Deckung (`--check-latest` prüft nur Byte-Identität, `--verify` nur das Manifest).
- **klasse:** delta-benennung-unvollstaendig

**F-3**
- **kategorie:** LOW
- **quelle:** Hard Rule `AGENTS.md` §3.3 / `MR-013`
- **pfad:** Commit `f845e8b` (`harness/conventions/done/MR-026-baseline-v560.md`, Rename-Score R060)
- **befund:** Der MR-026-Move bündelt `git mv` mit 34 geänderten Zeilen (17 Link-Tiefen-Fixes) in einem Commit; §3.3 verlangt zwei Commits und deklariert die Ausnahme nur für Slice-Lifecycle-Moves. Für MR-Lifecycle-Moves ist die Bündelung zweimal gelebte, aber nirgends deklarierte Praxis (Präzedenz MR-024-Move in slice-108, R076); der Score von 60 % liegt nahe der 50-%-Schwelle, die §3.3 selbst als Begründung nennt — eine künftige Hebung mit mehr Fixes reißt `git log --follow` still ab, ohne dass eine Regel den Konflikt auflöst. (Die Bündelung selbst ist gate-erzwungen: ein reiner Move wäre wegen der relativen Links doc-check-rot — genau das trat in slice-108 ein.)
- **verifizierbar:** nein — kein Gate prüft Rename-Scores oder die Zwei-Commit-Form.
- **klasse:** lifecycle-move-ausnahme-undeklariert

## DoD §4 gegen den Ist-Stand (auftragsgemäß mitgeprüft)

1. **Pin vendored + `--verify` + `--check-latest`:** Baum `.harness/baseline/v5.7.0/` mit beiden Bäumen + `SHA256SUMS` vorhanden, einziger Baseline-Baum; `fetch-baseline-cache.sh --verify` im Review **selbst ausgeführt: Exit 0, „51 Dateien, vollständig"**. `--check-latest` ist netzlos nicht reproduzierbar — liegt nur als Commit-Behauptung vor (beidseitig OK). Erfüllt mit dieser Beleg-Grenze.
2. **v5.6.0 entfernt, Verweise:** Baum entfernt ✓. Repo-weiter Zensus (19 Fundstellen außerhalb des vendorten Baums): genau **zwei** Markdown-Links (`done/slice-109…:35`, `done/MR-024…:8`) — beide eingefroren und von den neuen `.d-check.yml`-Tombstones **exakt** gedeckt (`in:`-Pfade stimmen dateigenau, `refs`-Glob deckt beide Ziele); alle übrigen Fundstellen sind Inline-Code (Präzedenz slice-106), Config selbst oder Historien-Prosa. `make doc-check` im Review selbst gelaufen: **Exit 0, 410 Dateien, 0 Befunde**. Erfüllt — die Pin-Spiegel aus F-1 liegen außerhalb des DoD-Wortlauts, gehören aber zur Sache.
3. **MR-Kette:** MR-028 Status Accepted, Index-Zeile mit Doppel-Anker (Voll-Slug + `mr-028`); MR-026 in `conventions/done/`, Aufgelöste-Zeile mit **wanderndem Doppel-Anker** (Voll-Slug + `mr-026`) — die eingefrorenen Verweise auf beide Anker-Formen (welle-78-results `#mr-026`, done/slice-106 und Roadmap-Drift-Log Voll-Slug) lösen weiter auf; Link-Tiefen in der verschobenen MR-026 stichprobenhaft korrekt (`../../conventions.md`, `../../../spec/…`, `../../../docs/plan/adr/…`). Erfüllt.
4. **Zwei-Hälften-Sektionsregel:** erfüllt, siehe Negativ-Proben 6/7.
5. **Delta-Audit:** je Regelwerks-Datei eine belegte Antwort ✓ (Template-Hälfte → F-2).
6. **`make gates` grün / GUARD:** in allen vier Botschaften mit expliziten Zahlen behauptet (55 Regeln / 409 bzw. 410 GUARD-Dateien); im Review per `make doc-check` und `make planning-check` nachvollzogen (beide Exit 0, 410/0). Der volle `gates`-Lauf und die GUARD-Läufe wurden nicht wiederholt — Behauptung plausibilisiert, nicht vollständig reproduziert.
7. **Unabhängiger Review:** dieser Lauf.

## Negativ-Proben (geprüft, ohne Befund)

1. **Kurs-Gegenprobe Regelwerk-Delta:** `git diff v5.6.0..v5.7.0 -- lab/regelwerk/` = exakt drei Dateien, +5/−3; Inhalte (README-Stand Welle 81 · modul-06 Zwei-Hälften-Fassung · modul-10 `klasse` als sechstes Feld) decken sich wortgenau mit den Behauptungen in Slice/MR-028; vendorter Baum trägt beide Neufassungen (Stichprobe positiv).
2. **Gegenrichtung Verweis-Hebung:** alle **24** Links auf `baseline/v5.7.0` (inkl. versteckter Verzeichnisse, skriptgestützt) lösen auf — Datei existiert, Anker existiert (auch die Umlaut-/Doppelstrich-Slugs wie `#was-ein-kommentar-trägt--code-konfiguration-skripte`).
3. **MR-028 inhaltlich gegen die MR-026-Vorlage:** Feldbestand und -reihenfolge identisch, Serien-Zählung konsistent („vierter Nachtrag"), Layout-/Pin-Trennung zu MR-023 sauber wiederholt („diese MR hebt den Pin, sie behauptet keine Konformität"), Auflösungs-Trigger spiegelbildlich; Index-Spalten decken sich mit dem Datei-Inhalt.
4. **Tombstone-Kommentar `.d-check.yml`:** trägt Kommentar-Klassen (Abgrenzung eingefroren/lebend, Kopplung an MR-021), Herkunfts-Feld in der Präzedenz-Form der beiden Vorgänger-Tombstones; die Behauptung „die übrigen … nur in Inline-Code" wurde fundstellen-genau verifiziert (alle Review-/done-Nennungen sind Backtick-Spans; `codepaths` prüft sie wegen Root-Präfix-Regel bzw. `docs/reviews/**`-Exempt nicht — kein stilles Loch).
5. **Commit-Schnitt-Ehrlichkeit (BEO-006-Probe):** alle vier Stats gegen ihre Botschaften gehalten — `6baa083` exakt die vier deklarierten Planungs-Dateien; `77d2496` ausschließlich der neue Baum (52 Pfade); `f845e8b` deckt jede Botschafts-Zeile (3+3+1+1+4-Link-Zählung aufgegangen, „zehn MR-Dateien" = 9×1 + MR-027×2, keine Beifänge); `78360eb` nur Roadmap + Slice-§9. Kein undeklarierter Inhalt.
6. **Roadmap-Sektionsregel vs. v5.7.0-modul-06:** getreue Paraphrase (Marker zusätzlich, genau-dann-Bedingung, Normalfall nach Eröffnung, Wächter beidseitig nur für die Marker-Hälfte, Liste als Ableitung) plus ehrliche Repo-Lage („bis dahin hält `wave-drift` den Aktiv-Status gegen die Datei-Zahl") — kein Inhalts-Drift.
7. **Substring- und Kennungs-Wächter:** der Marker-Wortlaut steht **nicht** in der §Offene-Wellen-Sektion (nur zweimal im Drift-Log, außerhalb des gewächterten `heading`-Blocks); die Sektions-Prosa nennt keine Wellen-Kennung — `welle-79` nur in der Zeiger-Zeile (zukunftsfest für die `many`-Bijektion; die Slice-Nennung „slice-110 beansprucht" ist keine Wellen-Kennung).
8. **Planungs-Konsistenz welle-79/slice-111:** widerspruchsfrei gegen die v5.7.0-Fassung (Bijektion beider Richtungen, Marker-Orthogonalität, Default byte-identisch, fail-closed); `ADR-0055` existiert mit Status Proposed, `DC-FA-PLAN-001` §Wellen-Invariante existiert; der Konsumenten-CR ist als bewusst nicht-vorliegende Datei sauber deklariert („Landung ist der Lastenheft-Commit"); Reihenfolge-Bindung slice-110→111 in Welle und Slice identisch begründet.
9. **§9-Audit-Belege:** die modul-10-Konform-Behauptung stimmt — `reviewer.md` 1.5.0 führt `kategorie · quelle · pfad · befund · verifizierbar · klasse` (§Output-Schema), `klasse` seit slice-090 (Commit `419e82e`); Herkunfts-These „Landung eigener Upstream-Notizen" deckt sich mit dem Drift-Log (Wiedervorlage slice-090 in der welle-78-Zeile).
10. **Referenz-Richtung/Kommentar-Klassen (HIGH-Anker):** keine neuen `status-provenance`-Marker im Range außerhalb des vendorten Baums; keine neuen Kommentare ohne Klasse; die Grenz-Kommentare der Prüf-Profile blieben wie in §2.4 deklariert unangetastet und bleiben inhaltlich wahr.
11. **Gate-Nachvollzug:** `make doc-check` Exit 0 (410/0) und `make planning-check` Exit 0 (410/0) — beide mit explizit geprüftem Exit (BEO-007-Arbeitsregel), keine stillen Grün-Pfade beobachtet.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 (F-1) |
| LOW | 2 (F-2, F-3) |
| INFO | 0 |

## Verdikt: APPROVE mit Auflagen

Der Kern des Slice ist sauber und belegt: das Regelwerks-Delta stimmt gegen das Kurs-Repo aufs Zeichen, die Verweis-Hebung ist für das `baseline/v5.6.0`-Muster vollständig (zwei eingefrorene Links, beide exakt getombstoned; 24 v5.7.0-Links samt Ankern auflösbar), die MR-Kette inklusive wandernder Doppel-Anker hält alle eingefrorenen Verweise am Leben, die Zwei-Hälften-Prosa ist eine treue, wächter-sichere Paraphrase, und der Commit-Schnitt deckt jede Botschafts-Zeile. Gates wurden im Review teilreproduziert (doc-check, planning-check, `--verify` — alle Exit 0).

**Auflage vor der Closure (F-1, MEDIUM):** die drei stehengebliebenen Pin-Spiegel (`AGENTS.md:32`, `harness/README.md:60`, `MR-021:17`) auf v5.7.0 heben — es ist die Wiederholung des bereits in slice-106 als Review-Auflage geheilten Musters, und `harness/README.md:60` ist in sich widersprüchlich (Baum v5.7.0, Quelle v5.6.0). **Steering-Loop-Signal** nach `reviewer.md` §Kontext-Eskalation: dieselbe Klasse zum wiederholten Mal über zwei Pin-Hebungen, drei Fundstellen in diesem Lauf — der Spiegel-Zensus der Pin-Hebung gehört mechanisiert oder als Checkliste in die MR-023/MR-028-Prozedur, nicht nur gemeldet. F-2 und F-3 sind nice-to-fix (Delta-Benennung um die Template-Hälfte ehrlich machen bzw. quell-genau scopen; die MR-Move-Ausnahme deklarieren) und blockieren die Closure nicht.
