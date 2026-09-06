# Review-Report — slice-202 (Baseline-Pin-Hebung auf v6.3.1)

**Review-Art:** Code/Diff (gegen Slice-Plan, MR-021/MR-051/MR-054/MR-055, `AGENTS.md` §3/§5, Baseline-Kanon)
**Gegenstand:** `d92e19c` (Plan) · `89c981d` (Beanspruchung) · `9eeb39f` (MR-060-Move) · `bd585ee` (Hauptbump)
**Skill:** `.harness/skills/reviewer.md` v1.13.0 @ `bd585ee`
**Modell-ID:** claude-opus-5[1m] · **Datum:** 2026-09-06
**Eingangs-Kontext:** der Slice-Plan; `harness/conventions/MR-065-baseline-v631.md`; MR-013/021/035/036/051/054/055; `.d-check.yml`; vendorter Baum `v6.0.0` (aus `bd585ee^`) gegen `v6.3.1`
**Eigene Läufe:** `make baseline-verify` (ok, 54 Dateien) · `make doc-check` (617 Dateien, 0 Befunde) · `make gates` (zehn Gates grün, Coverage 94,70 %)

---

## Findings

### F-1 · MEDIUM

- **quelle:** `MR-021` / `AGENTS.md` §5 („behauptet nicht mehr, als die Arbeit trägt")
- **pfad:** `harness/conventions/MR-065-baseline-v631.md:56` (gespiegelt in der Botschaft von `bd585ee`)
- **befund:** Der Zensus nennt „**52** lebende Dateien gehoben"; die im selben Satz geschlossene Aufzählung summiert sich auf 51 (1+1+1+32+2+2+4+8), und gemessen tragen 51 Nicht-Baseline-Dateien in `bd585ee` einen Versions-Lift — 53 berührte Dateien minus `.d-check.yml` (bekommt nur den Tombstone, der den **alten** Tag nennt) und minus der neu angelegten `MR-065`.
- **verifizierbar:** nein (kein Gate); nachzählbar über `git show --name-only --format= bd585ee | grep -v '^\.harness/baseline/'`
- **klasse:** `selbstauskunfts-zahl-weicht-von-gemessener-menge-ab`

### F-2 · MEDIUM

- **quelle:** `MR-051`
- **pfad:** `harness/conventions/MR-065-baseline-v631.md:83-87`
- **befund:** „**15** lebende Direktiven in 10 Dateien … **12** unverändert bestätigt"; gemessen liegen **16** baseline-gebundene `d-check:cite`-Direktiven in **11** lebenden Dateien (also 13 unverändert, bei korrekt gezählten 3 neu verankerten). Die elfte Datei ist der Slice-Plan selbst mit zwei Direktiven — er liegt in `in-progress/` und ist von `citations.scope` nicht ausgenommen.
- **verifizierbar:** teilweise — `git grep -c` außerhalb der drei eingefrorenen Verzeichnisse; das Gate prüft jede einzelne Direktive, nie ihre Anzahl
- **klasse:** `selbstauskunfts-zahl-weicht-von-gemessener-menge-ab`

### F-3 · MEDIUM

- **quelle:** Baseline `modul-05` §Zwei Schritte vor der Modus-Begründung
- **pfad:** Slice-Plan §8
- **befund:** „Register durchgegangen (gemergter Stand, **34** Verzeichnisse)"; das Beobachtungs-Register führt 33 Beobachtungs-Verzeichnisse, sowohl auf dem Plan-Commit `d92e19c` als auch auf `HEAD`.
- **verifizierbar:** nein (kein Gate)
- **klasse:** `selbstauskunfts-zahl-weicht-von-gemessener-menge-ab`

> **Steering-Loop-Signal:** F-1/F-2/F-3 sind die **dritte** Wiederholung derselben Klasse in einer Sitzung. Sie ist als `BEO-ALL/eigene-menge-gemessen-fremde-behauptet` bereits registriert und wird in §8 dieses Slice als *abgefangen* geführt — die drei Instanzen zeigen, dass die Abfangstelle nur für die Bundle-Zahl griff, nicht für die selbst erzeugten Mengen.

### F-4 · MEDIUM

- **quelle:** `MR-021` / `BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`
- **pfad:** Slice-Plan §3
- **befund:** Die §3-Zeile lautete vor `bd585ee` `| .harness/baseline/v6.0.0/ | entfällt | …` und benannte damit den zu **entfernenden** Baum; der Zensus hat sie mitgehoben. Die Tabelle führt jetzt dieselbe Pfad-Angabe zweimal — einmal als „neu", einmal als „entfällt". MR-065 zählt drei gefangene Übersetzungsfehler-Klassen auf; diese vierte, die identifizierende Nennung des alten Baums in einem lebenden Dokument, ist nicht darunter.
- **verifizierbar:** nein — beide Pfade lösen auf, `make doc-check` bleibt grün
- **klasse:** `mechanischer-lift-trifft-identifizierende-nennung`

### F-5 · MEDIUM

- **quelle:** `MR-054` („Die Spanne trägt die *vorschreibende* Zeile, nicht irgendeine aus demselben Absatz")
- **pfad:** Slice-Plan §8, zweiter Vorprüfungs-Block
- **befund:** Der Block ankert `modul-05-planning-harness.md:234-236` und zitiert die Notier-/Lese-**Nebenregel** des Schritt-Punkts, exakt die von MR-054 benannte Form, nicht die vorschreibende Zeile „2. **Offene Beobachtungen sichten.**". Die beiden Vorgänger-Slices (`done/slice-200`, `done/slice-201`) ankern dafür `:229-229`.
- **verifizierbar:** nein — `citations` prüft Wortgleichheit, nicht welche Zeile des Absatzes gewählt wurde
- **klasse:** `cite-belegt-nebenregel-statt-vorschrift`

### F-6 · MEDIUM

- **quelle:** `MR-051` §Geltungsbereich
- **pfad:** `harness/conventions/MR-065-baseline-v631.md`, Adaptions-Review
- **befund:** „Der Adaptions-Review ist durch alle 38 aktiven Einträge gelaufen … **alle 38 bleiben gültig**" schließt MR-051 ein, dessen Geltungsbereich wörtlich sagt, in den eingefrorenen Verzeichnissen stehe **keine** `cite`-Direktive. Tatsächlich tragen `done/slice-200` und `done/slice-201` je zwei sowie `conventions/done/MR-057` und `MR-058` je eine; die vier in `done/` zeigen seit `bd585ee` auf den entfernten `v6.0.0`-Baum. Die Index-Zeile in `conventions.md` weist die Feld-Aussage bereits als „durch MR-054 überholt" aus — der Eintrag selbst wurde als gültig durchgereicht.
- **verifizierbar:** nein — kein Gate liest `Geltungsbereich`-Felder
- **klasse:** `adaptions-review-bestaetigt-veraltetes-geltungsfeld`

### F-7 · MEDIUM *(Kontext-Eskalation: Basis LOW, im Gate-Pfad eine Stufe höher)*

- **quelle:** `AGENTS.md` §3.7 (Kommentar-Klassen Abgrenzung/Grenze) / `MR-065`
- **pfad:** `.d-check.yml` (wortgleich gespiegelt in `MR-065`)
- **befund:** Der Tombstone-Kommentar begründet den einen Glob damit, die übrigen eingefrorenen Klassen — namentlich `done/`-Slices — nennten `v6.0.0` „entweder als reine Versions-Aussage **ohne Pfad** oder in Inline-Code". `done/slice-200` und `done/slice-201` tragen jedoch vier `d-check:cite`-Direktiven mit vollständig auflösendem Pfad in den entfernten Baum — weder Versions-Aussage noch Inline-Code. Der Lauf schweigt dort, weil `citations.scope` `docs/plan/planning/done/**` ausnimmt und `links` keine HTML-Kommentare liest, nicht aus dem angegebenen Grund; die **Entscheidung** (ein Glob) trägt, die **Begründung** nicht.
- **verifizierbar:** teilweise — der grüne Lauf bestätigt die Entscheidung, nicht die Begründung
- **klasse:** `gate-kommentar-nennt-falschen-grund`

### F-8 · LOW

- **quelle:** `MR-055`
- **pfad:** `harness/conventions/MR-055-symlink-als-pin-traeger.md`
- **befund:** „Und der Ausfall wäre **still**: **vier** tote Aliase, ein grüner Lauf" beschreibt den Bestand zur Zeit von slice-176; heute stehen acht Baseline-Aliase unter `.claude/rules/` (plus drei repo-interne). Der Befund ist vorbestehend — die Datei wurde in `bd585ee` angefasst (Pfad-Lifts) und im Adaptions-Review als gültig geführt.
- **verifizierbar:** nein
- **klasse:** `zahl-in-mr-prosa-veraltet`

### F-9 · INFO

- **quelle:** Maintainability
- **pfad:** `harness/conventions/MR-065-baseline-v631.md`, Delta-Tabelle
- **befund:** Die Tabelle ordnet jedes Delta einem der vier Tags zu. Vendored liegt nur der `v6.3.1`-Endstand; die Zwischen-Bundles sind im Repo nicht vorhanden, die Tag-**Zuordnung** ist damit aus dem Repo nicht nachvollziehbar (der **Inhalt** jeder Zeile ist gegen den Aggregat-Diff belegt). Undokumentierte Annahme, kein Defekt.
- **verifizierbar:** nein (Netz)
- **klasse:** `per-tag-zuordnung-ohne-repo-beleg`

---

## Negativbefunde (geprüft, ohne Befund)

1. **Bundle-Delta-Zahlen** — nachgemessen mit `diff -rq -I` gegen `bd585ee^`: 54 Inhalts-Dateien, 13 mit echtem Delta, 1 neu, naiver Diff meldet 53. Alle vier Zahlen korrekt.
2. **Zensus-Klassen im Detail** — 32 gehobene aktive `MR-*`-Einträge, 8 Baseline-Aliase, 4 Planning-Docs, 2 Skills, 2 Spec-Straten: je einzeln nachgezählt, alle korrekt. Die 5 nicht angefassten aktiven MR-Dateien tragen nachweislich keinen pin-gebundenen Verweis.
3. **Zensus-Vollständigkeit vorwärts** — kein lebender Verweis trägt mehr `v6.0.0`.
4. **Zensus-Vollständigkeit rückwärts** — keine Datei unter den fünf eingefrorenen Verzeichnissen und kein `observation.md`/`evidence/*.md` wurde in `bd585ee` angefasst.
5. **Die zwei von MR-060 benannten Übersetzungsfehler-Klassen** — `../baseline/`-relative Pfade im Reviewer-Skill (fünf Vorkommen) und die Release-Download-URL ohne `.harness/baseline/`-Segment (drei Vorkommen): alle auf `v6.3.1`.
6. **Die dritte, neu benannte Klasse** (Link-Text vs. Link-Ziel) — in sich konsistent.
7. **Live-Zeiger vs. historische Aussage** — beide als Live-Zeiger gehobenen Vorkommen sind welche; die bewusst stehen gelassenen sind ebenfalls korrekt klassifiziert.
8. **cite-Spannen** — alle 16 lebenden Direktiven lösen wortgleich auf, `citations` ist fail-closed. Die drei neu verankerten zusätzlich von Hand gegen den Zieltext gehalten. „0 entfallen" bestätigt.
9. **Tombstone-Umfang** — nur ADR-0083 trägt einen Markdown-Link in den entfernten Baum; „EIN Glob, nicht fünf" ist als Entscheidung gemessen richtig (Begründung: F-7).
10. **MR-013-Konformität von `9eeb39f`** — reiner Move plus Link-Tiefenfixes plus der eine externe Referrer; alle neun relativen Ziele lösen vom neuen Ort auf, kein Referrer zeigt mehr auf den alten Pfad.
11. **Index-Konsistenz** — 38 Dateien ↔ 38 Zeilen §Aktive Adaptionen; 27 ↔ 27 §Aufgelöste.
12. **MR-035/MR-036** — Auflösungs-Trigger nicht eingetreten; beide bleiben zu Recht aktiv.
13. **Delta-Inhalt gegen das Bundle** — jeder der vier beschriebenen Stränge ist im Aggregat-Diff belegt.
14. **Abgrenzung Hebung vs. Adoption** — als Reihenfolge-Bedingung begründet und in Plan §1, DoD und MR-065 gleichlautend geführt; keine „mitgemacht"-Behauptung.
15. **Gate-Läufe** — `make baseline-verify` ok, `make doc-check` 617/0, `make gates` grün über alle zehn Gates (Coverage 94,70 % ≥ 93 %).
16. **HIGH-Prüffragen** — kein Produktcode berührt, kein Gate-Pfad mit stillem Grün, keine Suppression, keine Schwellen-Senkung, kein Netzzugriff außerhalb `external`.
17. **Kommentar-/Zustandsfeld-Klassen** — der einzige neue Config-Kommentar trägt Abgrenzung und Grenze, keine Slice-Nummer.
18. **Beobachtungs-Zähler im Plan** — je gegen die `evidence/`-Dateien gezählt, alle korrekt.

**Nicht geprüft** (außerhalb der Rolle bzw. offline nicht möglich): die DoD-Abhakung und die Risiko-Ausgänge (Verifikation, getrennter Kontext); die Zuordnung der Deltas zu den Zwischen-Tags (Netz, siehe F-9); `make baseline-freshness` (Netz).

---

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 7 (F-1 … F-7) |
| LOW | 1 (F-8) |
| INFO | 1 (F-9) |

## Verdikt

**Blockiert** (sieben MEDIUM). Die **mechanische** Hälfte des Slice ist sauber und nachgemessen: der Baum steht vollständig auf `v6.3.1`, die acht Aliase hängen um, der Zensus ist in beide Richtungen dicht, die Tombstone-Entscheidung ist richtig, der Move-Commit MR-013-konform, und alle zehn Gates sind im eigenen Lauf grün.

Was blockiert, ist die **Selbstauskunft** darüber: drei Zahlen, die von der gemessenen Menge abweichen (F-1/F-2/F-3), eine überhobene identifizierende Nennung, die der Slice selbst als Risiko benannt hatte (F-4), ein Vorprüfungs-Beleg auf einer Nebenregel statt der Vorschrift (F-5), und zwei Begründungen, die ihren Gegenstand nicht tragen (F-6/F-7). Keiner der sieben verlangt, den Bump zu wiederholen — alle sind an Text korrigierbar. F-1/F-2/F-3 gehören zusammen als Steering-Loop-Signal behandelt, nicht als drei Einzelkorrekturen: die Klasse ist registriert und in diesem Slice erneut dreimal aufgetreten, und zwar genau dort, wo die Menge **selbst erzeugt** statt fremd übernommen wurde.
