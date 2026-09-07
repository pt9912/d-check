# Review-Report — slice-208, Runde 1 (Delta-Audit, DoD 1)

**Review-Art:** Plan/Design — geprüft wird der **Delta-Audit** gegen seine Quelle (den in `slice-207` gemessenen Delta), gegen den Baseline-Kanon `v6.5.0` und gegen den eigenen Bestand. **Nicht** geprüft: die Form-Migration (DoD 2/3) — sie ist Runde 2 und existiert noch nicht.
**Gegenstand:** `e35bb53` (Vorprüfungen, MR-066 vorab) · `83cd116` (Beanspruchung) · `3d15de7` (Delta-Audit, DoD 1)
**Skill:** `.harness/skills/reviewer.md` v1.13.0 @ `3d15de7`
**Modell-ID:** claude-opus-5[1m] · **Datum:** 2026-09-07
**Eingangs-Kontext:** der Slice-Plan `slice-208`; `slice-207` §2 (die Delta-Liste) und §9; `MR-066`, `MR-053`, `MR-054`, `MR-013`, `MR-000`; `AGENTS.md` §3.7/§5/§6; `harness/README.md` §Minimal agent workflow; `.d-check.yml` (`trace`, `ignore-refs`, `citations.scope`); `DC-FA-CLI-009`, `DC-FA-CLI-011`, `DC-FA-COV-001`; das Beobachtungs-Register (36 Verzeichnisse); die Vorgänger-Adoptionen `slice-107` (Archiv-Volltext) und `slice-203`; die vendorten Bäume `v6.3.1` (aus `4ff020e~1`) und `v6.5.0`.
**Eigene Läufe:** `make gates` (zehn Gates grün, 689 Dateien / 0 Befunde, Coverage 94,60 %) · vollständige Re-Messung des Deltas mit den in `slice-207` §2 ausgeschriebenen Filtern (`diff -rq` ⇒ 35 · `-I` ⇒ 27 · `-w -B -I` ⇒ **12**, Namensmenge identisch, 55/55 Dateien — reproduziert exakt) · Zeilen-Zählung je der zwölf Dateien · Volltext-Diff jeder der zwölf · repo-weite Messung der Sensor-Verweise (26 lebend / 1 eingefroren) · repo-weite Messung der Baseline-**Links** in eingefrorenen Klassen (19 `Accepted`-ADRs, 4 `done/`-Slices, 6 aufgelöste `MR`) · Register-Zählung (36 Verzeichnisse, Evidence-Dateien je Eintrag) · Zählung der gelebten Ausschluss-Abschnitte (8 `done/`-Slices) und ihrer Folge-Slice-Adressen (2).

**Zitier-Form.** Baseline-Stellen stehen als Tag + Pfad in Inline-Code statt als Link, Slices als Kennung statt als Lifecycle-Pfad — die Form, die `v6.5.0` · `regelwerk/grundlagen-harness-dateien.md` §Ein einfrierendes Artefakt … für einfrierende Artefakte verlangt. Dieser Report ist eines.

---

## Regel → Antwort des Audits → Urteil

| # | Antwort des Audits | Urteil | Kern |
|---|---|---|---|
| R1 | übernommen · fünf Träger | **trägt** | Regel korrekt gefasst; alle fünf Träger bestätigt (`modul-05`, `slice.template`, `templates/README`, `modul-06`, `modul-09`) |
| R2 | übernommen · Träger `modul-05` | **trägt mit Vorbehalt** | Regel korrekt; Träger **unvollständig** (`slice.template` trägt sie ebenfalls), und die empirische Begründung ist falsch gezählt (F-6) |
| R3 | übernommen · Träger `slice.template`, `templates/README` | **trägt mit Vorbehalt** | Regel korrekt; **`modul-05` fehlt als Träger** — dort steht sie normativ (F-4) |
| R4 | übernommen (ohne Handlung) · Träger `modul-09` | **trägt nicht** | Kein lebender Träger dieses Repos führt die Regel; „übernommen" ist eine Behauptung ohne Deckung (F-2). Träger zusätzlich unvollständig (F-4) |
| R5 | übernommen, mit Handlung · offen sei „die Setzung" | **trägt nicht** | Die Prämisse über das eigene Produkt ist falsch: ADRs **entlasten nicht** (F-1). Zwei weitere Hälften der Kanon-Regel übersprungen (F-1b) |
| R6 | übernommen, Bestand konform | **trägt in der Zahl, nicht im Umfang** | 26/1 reproduziert exakt; die Regel ist aber nur in **einer** ihrer drei Formen gemessen (F-3), und die Fundstelle ist ein Zitat, keine Referenz (F-9) |
| R7 | übernommen, template-forward · Träger die vier Vorlagen | **trägt nicht** | Der benannte Träger (Reviewer-Skill) führt die Regel **nicht**; die Begründung „eigene Review-Form" ist ohne Deckung (F-5). `grundlagen-harness-dateien` fehlt als Träger (F-4) |
| — | „sieben Regeln in zwölf Trägern" | **unvollständig** | Eine **achte** Regel des Deltas hat keine Antwort (F-3); R6 und R7 sind im Kanon **eine** Regel in drei Formen, nicht zwei |

---

## Findings

### F-1 · HIGH · R5 stützt sich auf eine falsche Aussage über das eigene Produkt

- **quelle:** [`DC-FA-CLI-011`](../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code) §Out-of-Scope · [`DC-FA-CLI-009`](../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix) · `v6.5.0` · `regelwerk/grundlagen-traceability.md` §Die zweite Richtung: Anforderung → Beleg
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md:83-89` (Zeile 68 der Tabelle)
- **befund:** Der Audit schreibt, der `trace`-Block deklariere *„`adrs` und `slices` als Quellen, die eine Anforderung entlasten"*, und leitet daraus die einzige Handlung am Bestand ab. Das Lastenheft sagt das Gegenteil: Waise ist eine Anforderung *„ohne referenzierenden Slice **und** … ohne Coverage-Referenz"*, und `DC-FA-CLI-011` §Out-of-Scope nennt es ausdrücklich — *„eine bloße ADR-Referenz ohne Slice/Coverage deckt weiterhin **nicht** ab"*. Die `adrs:`-Untersektion konfiguriert eine **Anzeigespalte**, keine entlastende Quelle. Die gesetzte Entlastung ist damit exakt der Kurs-Vorschlag (Slice entlastet, ADR steht als Spalte), und sie ist bereits aufgeschrieben — in `DC-FA-CLI-011` und `DC-FA-COV-001`. Die als offen bezeichnete Handlung zielt auf eine Konfiguration, die es nicht gibt.
- **verifizierbar:** ja — `spec/lastenheft.md:568` gegen `.d-check.yml:824-836`; ein Gate-Lauf ist nicht nötig, aber `make completeness-check` belegt das Verhalten.
- **klasse:** `aussage-ueber-eigenes-produkt-ohne-spec-abgleich`
- **Warum HIGH:** Die Aussage liegt auf dem Gate-Pfad. Würde die abgeleitete „Handlung" ausgeführt, indem `adrs` als entlastend behandelt oder so dokumentiert wird, verschöbe sie die Schwelle von `make completeness-check` — eine Gate-Lockerung ohne ADR (`AGENTS.md` §3.6). Kontext-Eskalation nach Reviewer-Skill §Kontext-Eskalation.

### F-1b · HIGH · Die Kanon-Regel hinter R5 hat vier Aussagen; der Audit beantwortet zwei

- **quelle:** `v6.5.0` · `regelwerk/grundlagen-traceability.md` §Die zweite Richtung: Anforderung → Beleg (Zeilen 30, 36, 43, 52)
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md:68` · `:83-89`
- **befund:** Die neue Sektion trägt vier Setzungen. Der Audit nennt zwei (*erzeugt statt gepflegt*; *die Entlastung ist eine Setzung*). Übersprungen sind: (a) *„Der Vorschlag dieses Kurses: der **Slice** … deshalb steht die ADR als eigene Spalte da und nicht als Quittung"* — samt der Pflicht, einen anderen Schnitt **zu deklarieren**, „wie jede Abweichung von der Baseline"; (b) *„Bericht und Gate sind derselbe Lauf, nicht zwei Werkzeuge"*. Beide sind für dieses Repo einschlägig: (a) ist die Regel, an der sich `trace.coverage` als *dritte, opt-in Referenzklasse* (`DC-FA-COV-001`) messen lassen muss — der Kanon nennt genau diesen Fall („eine kuratierte Nachweis-Datei als entlastende Quelle") als deklarationspflichtige Abweichung; (b) beschreibt `--trace` gegen `--trace --require-complete` und ist ohne Zutun erfüllt. Ein Delta-Punkt ohne Antwort ist nach DoD (1) *„ein offener Punkt, kein stilles Übergehen"*.
- **verifizierbar:** ja — Volltext der Sektion gegen die vier Sätze der Audit-Zeile.
- **klasse:** `regel-teilweise-beantwortet`

### F-2 · MEDIUM · R4 ist als „übernommen" beantwortet, obwohl kein lebender Träger dieses Repos sie führt

- **quelle:** `v6.5.0` · `regelwerk/modul-09-implementierung.md` §Minimal Agent Workflow · [`AGENTS.md` §6](../../AGENTS.md#6-minimal-agent-workflow)
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md:67` · `AGENTS.md:597` · `harness/README.md` §Minimal agent workflow, Schritt 4
- **befund:** `modul-09` bindet die Regel an **Schritt 4** des Acht-Schritt-Workflows: *„Die Plan-Ausgabe in Schritt 4 nennt Out-of-Scope"*, und *„Nimmt der Lauf etwas mit, das §1 ausschließt, ist das eine Plan-Änderung und gehört vor den Code, nicht in den Bericht danach"*. Dieses Repo führt denselben Workflow als adoptierte Kopie; sein Schritt 4 lautet vollständig *„Kleinste sinnvolle Änderung planen."* — ohne Out-of-Scope und ohne die Plan-Änderungs-Pflicht. Dasselbe in `harness/README.md`. Der Audit antwortet dennoch „übernommen" ohne Handlungs-Vermerk und schreibt daneben, R5 verlange *„als einzige eine Handlung am Bestand"*. Der Träger, an dem R4 fehlt, ist weder die Slice-Vorlage noch das Closure-Profil noch ein Konventions-Eintrag — er steht damit in keiner Zeile der §2-Tabelle und fällt aus DoD (2) und (3) heraus.
- **verifizierbar:** ja — `sed -n '589,602p' AGENTS.md`; kein Gate fängt es (`make gate-consistency` prüft Targets, nicht Workflow-Schritte).
- **klasse:** `adoption-behauptet-ohne-traeger`

### F-3 · HIGH · Eine achte Regel des Deltas hat keine Antwort — und sie trifft ein Ventil, das dieses Repo bei jedem Bump vergrößert

- **quelle:** `v6.5.0` · `regelwerk/grundlagen-harness-dateien.md` §Ein einfrierendes Artefakt … (Zeilen 293–326) · `AGENTS.md` §3.6 · DoD (1) des Slice
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md:64-70` (Tabelle) · `.d-check.yml:31-231` (`ignore-refs`)
- **befund:** Das Delta von `grundlagen-harness-dateien.md` (46 Zeilen für angeblich **eine** Regel) trägt drei Bestandteile, die in R6 nicht aufgehen. Erstens **erweitert es die Klasse**: einfrierend sind jetzt Review-Report, Closure-Notiz, **Archiv-Stub, `Accepted`-ADR und geschlossener Slice** — vorher nur die ersten beiden. Zweitens nennt es **drei Formen derselben Regel** (Gate-Token · `slice-NNN` statt Lifecycle-Pfad · Baseline-Stelle als Tag + Pfad in Inline-Code); der Audit spaltet sie in R6 und R7 und ordnet die dritte Form einer Datei zu, in der sie nicht steht. Drittens — und das ist die fehlende Antwort — schließt der Absatz *„Und die Reparatur ist teurer als die Vermeidung"* mit der Setzung, ein Ausnahme-Ventil im Prüfbereich sei *„eine Gate-Senkung mit eigener Begründungslast"*. Genau dieses Ventil betreibt dieses Repo: `ignore-refs` trägt **acht** Tombstone-Familien für entfernte Baseline-Bäume, zuletzt für `v6.3.1` mit `MR-067`, und der Bestand hinter dem Ventil ist gemessen **kein Einzelfall** — 19 `Accepted`-ADRs, 4 `done/`-Slices und 6 aufgelöste `MR`-Einträge tragen Markdown-Links in entfernte Bäume. Zu dieser Setzung steht im Audit nichts, weder „übernommen" noch „nicht anwendbar".
- **verifizierbar:** ja — `git show 4ff020e~1:.harness/baseline/v6.3.1/regelwerk/grundlagen-harness-dateien.md` gegen die heutige Fassung; Zählung der Links per `grep -rl "](\(\.\./\)*\.harness/baseline/"` über `docs/plan/adr`, `docs/plan/planning/done`, `harness/conventions/done`.
- **klasse:** `regel-ohne-antwort-im-delta-audit`

### F-4 · MEDIUM · Die Träger-Zuordnung ist datei-weise statt regel-weise gebildet; vier Zuordnungen fehlen

- **quelle:** DoD (1) („übernommen (**mit Träger**)") · `v6.5.0` · `regelwerk/modul-05-planning-harness.md` · `templates/docs/plan/planning/slice.template.md` · `regelwerk/grundlagen-harness-dateien.md`
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md:64-70`
- **befund:** Der Audit ordnet jeder Datei die Regel zu, die er als ihre Hauptaussage liest, statt jede Regel gegen jede Datei zu halten. Gemessen fehlen vier Träger: **R2** steht auch in `slice.template` (*„mit Kennung — und die Kennung muss den Punkt auch annehmen"*, Zeile 63); **R3** steht normativ in `modul-05` (der Blockquote Zeile 219–222 und der Absatz *„Der Begründungsblock dagegen ist bedingt"*, Zeile 258) — die beiden zugeordneten Träger ziehen sie nur nach; **R4** steht auch in `modul-05` (Zeile 214) und in `slice.template` (Zeile 68); **R7** steht normativ in `grundlagen-harness-dateien` (Zeile 309) — die vier zugeordneten Vorlagen wenden sie an. Die Folge ist nicht kosmetisch: R3 und R4 erscheinen dadurch als Vorlagen-Fragen, obwohl sie im Regelwerk stehen, und R4s fehlender Träger (F-2) bleibt unentdeckt.
- **verifizierbar:** ja — Volltext-Diff je Datei gegen die sieben Regel-Formulierungen.
- **klasse:** `zuordnung-nach-datei-statt-nach-regel`

### F-5 · MEDIUM · R7 benennt einen Träger, der die Regel nicht führt, und begründet das mit einer nicht belegten Aussage über den eigenen Bestand

- **quelle:** `.harness/skills/reviewer.md` §Ablage · [`MR-000`](../../harness/conventions.md#mr-000--baseline-aussage) (Abweichungen sind zu deklarieren)
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md:96-100`
- **befund:** Der Audit schreibt: *„Träger bei uns ist nicht die Baseline-Vorlage — dieses Repo führt eine eigene Review-Form —, sondern der Reviewer-Skill."* Beide Hälften halten nicht. (a) Der Reviewer-Skill nennt die Baseline-Vorlage selbst als Ziel-Form (*„**Kopf-Metadaten** (Ziel-Form `review-report.template.md`)"*), es existiert keine lokale Kopie der Vorlage, und kein `MR` deklariert eine abweichende Review-Form — eine „eigene Review-Form" wäre nach `MR-000` eine undeklarierte Abweichung, nicht ein Träger. (b) Der benannte Träger **führt die Regel nicht**: der Skill enthält keine Zitier-Form-Aussage. Die Antwort „übernommen, template-forward" und der benannte lokale Träger schließen einander aus — ein lokaler Träger, der die Regel nicht trägt, macht die Adoption zu einer Handlung, nicht zu einem Zustand.
- **verifizierbar:** ja — `grep -n "Zitier-Form" .harness/skills/reviewer.md` (0 Treffer) · `find . -name "review-report*"` (nur der vendorte Baum).
- **klasse:** `traeger-benannt-der-die-regel-nicht-fuehrt`

### F-6 · MEDIUM · Die empirische Begründung von R2 zählt einen veralteten Bestand und zählt in ihm falsch

- **quelle:** [`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../plan/planning/observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md) · [`BEO-ALL/eigene-menge-gemessen-fremde-behauptet`](../plan/planning/observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md:77-82` · `:37-42`
- **befund:** *„Von den sechs gelebten Ausschluss-Abschnitten nennen vier einen Folge-Slice."* Gemessen: **acht** `done/`-Slices tragen beide Haus-Form-Überschriften (`slice-200` bis `slice-207`) — die Sechs stammt aus `af29684`, als `slice-206` und `slice-207` noch nicht geschlossen waren, und wurde im Audit-Commit unverändert weitergetragen. Und von den acht Ausschluss-Abschnitten nennen **zwei** einen Folge-Slice als Adresse (`slice-202` → `slice-203`, `slice-207` → `slice-208`); die beiden übrigen Treffer sind Rückverweise auf **Vorgänger** (`slice-200` nennt `slice-197` als Gegenstand eines Ausschlusses, `slice-204` nennt `slice-195`/`197`/`203` als Fälle, auf die es *nicht* rückwirkt) — keine Adresse im Sinne von R2. Die Zahl, mit der R2s Relevanz belegt wird, misst damit weder den aktuellen Bestand noch den Gegenstand der Regel.
- **verifizierbar:** ja — `grep -rl "^## 3\. Ausdrücklich NICHT" docs/plan/planning/done/` ⇒ 8; Abschnitts-Extraktion je Datei ⇒ 2 Vorwärts-Adressen.
- **klasse:** `zaehlmethode-misst-proxy-statt-gegenstand`

### F-7 · MEDIUM · Die Vorprüfung nennt einen Zähler, den das Register nicht trägt

- **quelle:** `v6.5.0` · `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register (*„Der Zähler wird abgeleitet, nicht geführt — er ist die Zahl der gültigen Evidence-Dateien"*) · [`MR-053`](../../harness/conventions.md#mr-053)
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md:278-281` · `:322` · `docs/plan/planning/observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/`
- **befund:** Die Sichtung führt `zaehlmethode-misst-proxy-statt-gegenstand` mit **2×** und stützt darauf die Risiko-Dichte in §8. Das Verzeichnis trägt **eine** Evidence-Datei (`slice-205.md`), und sein `state.md` sagt weiterhin *„erstes Auftreten"*. Der zweite Beleg wurde in der Closure von `slice-207` **behauptet** (Notiz und Commit-Botschaft: *„bekommt seinen zweiten Beleg"*), aber nie geschrieben — `9d7140e` legt nur `mechanical-id-rewrite-misses-frozen-classes/evidence/slice-207.md` an. Die Vorprüfung hat die Prosa des Vorgängers gezählt statt das Register. Die übrigen vier genannten Zähler (11×, 3×, 7×, 15×) und die Gesamtzahl 36 reproduzieren exakt.
- **verifizierbar:** ja — `ls docs/plan/planning/observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/evidence/` ⇒ eine Datei; `git show --stat 9d7140e -- docs/plan/planning/observations/`.
- **klasse:** `zaehler-aus-prosa-statt-aus-register`

### F-8 · MEDIUM · Ein einschlägiger Register-Eintrag mit 9× ist nicht gesichtet

- **quelle:** `v6.5.0` · `regelwerk/modul-05-planning-harness.md` §Zwei Schritte vor der Modus-Begründung (Schritt 2)
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md:259-291`
- **befund:** Die Sichtung führt vier einschlägige Einträge und einen geprüft-ausgeschlossenen. Nicht gesichtet ist [`BEO-ALL/eigene-menge-gemessen-fremde-behauptet`](../plan/planning/observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md) — **9×**, Sub-Area `*`, Stand *gemischt* (also teils offen), und inhaltlich die genaue Gestalt dieses Slice: *„Test vor jeder Messung: wer ändert die Menge, die ich zähle — und wer ändert die, über die ich rede?"*. Der Audit misst die **eigene** Menge (26/1 Sensor-Verweise, sauber gemessen) und sagt daneben über **fremde** Mengen aus, ohne zu messen: über den Delta („keine Regel ist nicht anwendbar"), über die beiden Vorgänger-Adoptionen (F-10), über die eigene Review-Form (F-5) und über den Ausschluss-Bestand (F-6). Drei der Befunde dieses Reports fallen in genau diese Klasse. Der Eintrag hat mit `slice-202`, `slice-203` und `slice-205` bereits drei Belege aus der unmittelbaren Nachbarschaft dieses Slice.
- **verifizierbar:** ja — Register-Verzeichnis gegen die Sichtungs-Liste in §7.
- **klasse:** `register-sichtung-unvollstaendig`

### F-9 · LOW · Die „genau eine" eingefrorene Fundstelle von R6 ist ein Zitat, keine Referenz

- **quelle:** `v6.5.0` · `regelwerk/grundlagen-harness-dateien.md` §Ein einfrierendes Artefakt …
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md:91-95` · `docs/reviews/2026-09-06-slice-206-mentions-modul-review.md:65`
- **befund:** Der Audit meldet *„Aus dem eingefrorenen Bestand tut es genau einer, und der steht in einem Review-Report."* Die Fundstelle ist die wörtliche Wiedergabe der Zeile, die ein Bruch-Test entfernt hat — ``Bruch-Test C entfernt den Link `[`make test`](sensors/test.md)` aus der Sensors-Tabelle`` —, also ein Zitat des geprüften Gegenstands, das die neue Vorlage ausdrücklich vom Verbot ausnimmt. Der Bestand ist damit nicht „bis auf einen" konform, sondern in dieser Form vollständig konform; die Schlussfolgerung „kein Handlungsbedarf" bleibt richtig, ihr Beleg ist es nicht.
- **verifizierbar:** ja — `sed -n '65p' docs/reviews/2026-09-06-slice-206-mentions-modul-review.md`.
- **klasse:** `zitat-als-fundstelle-gezaehlt`

### F-10 · MEDIUM · Der Vergleich mit den beiden Vorgänger-Adoptionen reproduziert für eine von beiden nicht

- **quelle:** `AGENTS.md` §5 (*„Eine Commit-Botschaft oder Closure-Notiz behauptet nicht mehr, als die Arbeit trägt"*) · [`BEO-ALL/commit-message-overclaims-work`](../plan/planning/observations/BEO-ALL/commit-message-overclaims-work/observation.md)
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md:102-106`
- **befund:** *„Die beiden Vorgänger-Adoptionen trugen **je** mehrere Nicht-anwendbar-Antworten (Wellen-Betrieb, Mehr-Schreiber-Teile)."* Für `slice-107` trifft es zu (§9 Stufen-Audit, im Archiv `done/welle-78/archiv.zip`, führt fünf `n. a.`-Antworten — BF-Reconciliation-Register, `ARC-*`, Golden Set, `lab/team-sim/`, Mehr-Schreiber-Teile). Für `slice-203` trifft es nicht zu: der Slice führt überhaupt keinen Regel-für-Regel-Audit — kein Vorkommen von „nicht anwendbar", „n. a." oder „konform", keine Regel-Tabelle, neun Abschnitte in Haus-Form ohne Audit-Sektion. Auch die §2-Zeile *„die Form, die der vorige Adoptions-Slice gelebt hat"* trifft damit nur auf `slice-107` zu, nicht auf den unmittelbaren Vorgänger. „Wellen-Betrieb" erscheint in `slice-107` außerdem in keiner der `n. a.`-Begründungen.
- **verifizierbar:** ja — `grep -n "nicht anwendbar\|n\. a\.\|konform" docs/plan/planning/done/slice-203-v631-template-adoption.md` ⇒ 0 Treffer; `unzip -p docs/plan/planning/done/welle-78/archiv.zip …slice-107….md` für die Gegenprobe.
- **klasse:** `vergleich-mit-eigenem-bestand-ungemessen`

### F-11 · LOW · Die Botschaft des Beanspruchungs-Commits widerspricht ihrem eigenen Diff

- **quelle:** `AGENTS.md` §5 · [`BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt`](../plan/planning/observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md)
- **pfad:** Commit `83cd116`, Botschaft Absatz 2
- **befund:** Die Botschaft schreibt: *„geloescht sind genau die Marker-Zeile und die Leerzeile dahinter, und der Diff zeigt genau eine entfernte Zeile."* Der Diff zeigt **zwei** entfernte Zeilen (`-Nichts in Arbeit.` und die Leerzeile). Der Satz war die ausdrückliche Selbstkontrolle gegen den Fehler des Vorgänger-Moves; er widerspricht sich innerhalb desselben Satzes und stimmt mit dem Diff nicht überein. Das Ergebnis im Dokument ist korrekt (kein Doppel-Leerraum vor `## Nächste Wellen`) — falsch ist nur der Beleg.
- **verifizierbar:** ja — `git show 83cd116 -- docs/plan/planning/in-progress/roadmap.md`.
- **klasse:** `selbstkontrolle-widerspricht-ihrem-beleg`

### F-12 · LOW · „26 Verweise aus lebenden Artefakten" — es ist ein Artefakt

- **quelle:** Maintainability
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md:91`
- **befund:** Alle 26 Sensor-Links stehen in `harness/README.md`; kein zweites lebendes Artefakt verlinkt eine Sensor-Datei. Der Plural legt eine Streuung nahe, die es nicht gibt — und verdeckt, dass die R6-Konformität des lebenden Bestands an **einer** Datei hängt. (24 der 26 Links sind zudem distinkt; zwei Targets erscheinen doppelt.)
- **verifizierbar:** ja — `grep -o "](sensors/[a-z0-9-]*\.md)" harness/README.md | wc -l` ⇒ 26, `| sort -u | wc -l` ⇒ 24.
- **klasse:** `plural-ohne-streuung`

### F-13 · LOW · „R1–R4 sind die CR-Umsetzung" widerspricht „R2 ist neu und wurde nicht erbeten"

- **quelle:** Maintainability
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md:72` gegen `:77`
- **befund:** Zwei Absätze auseinander stehen beide Sätze unverbunden nebeneinander. Gemeint ist vermutlich *„R1–R4 kamen mit dem CR-umsetzenden Release"*; gelesen wird *„R1–R4 sind das, worum wir gebeten haben"*, und dann ist R2 ein Widerspruch. Für einen Leser, der später entscheiden muss, welche Regeln aus dem eigenen CR stammen (und damit bereits abgewogen sind) und welche nicht, ist das die falsche Auskunft.
- **verifizierbar:** nein — Lesbarkeit, kein Gate.
- **klasse:** `zwei-saetze-eine-menge-zwei-aussagen`

### F-14 · INFO · Die MR-066-Ersatz-Form ist formal erfüllt, macht den Audit aber unrevidierbar

- **quelle:** [`MR-066`](../../harness/conventions.md#mr-066) §Adaption
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md` §6
- **befund:** `MR-066` verlangt Grund **und** benannte Ersatz-Form; beides steht im Plan, und „mehrere Runden gegen je einen abgeschlossenen Stand" ist die erste der vier Formen, die `MR-066` selbst aufzählt. Der Vollzug ist mit diesem Report begonnen. Zwei Beobachtungen dazu, ohne dass `MR-066` verletzt wäre: (a) Der Grund lautet, eine Teilung zerrisse den Audit, *„und die Form-Migration hängt an derselben Liste"* — die Prüfung wird dann an genau dieser Naht geteilt. `MR-066` erlaubt das ausdrücklich (*„Nicht der Slice wird geteilt, sondern die Prüfung"*), aber der Grund argumentiert gegen eine Naht, die er anschließend benutzt. (b) Der Plan legt fest, Runde 2 nehme *„den Audit als gegeben"*. Damit ist Runde 1 die **einzige** Prüfung der sieben Antworten; jede hier übersehene Fehl-Antwort ist danach ungeprüft und wirkt über DoD (2)/(3) weiter. Die acht Befunde oben zeigen, dass diese Menge nicht leer war.
- **verifizierbar:** nein — Urteil über die Prüf-Architektur.
- **klasse:** `ersatz-form-ohne-zweite-lesung-des-fundaments`

---

## Negativbefunde (geprüft, ohne Befund)

- **Die Delta-Messung selbst reproduziert exakt.** Namensmenge identisch (55/55, keine neue, keine entfallene Datei); `diff -rq` ⇒ 35 Pfade, `-I` ⇒ 27 Markdown-Dateien, `-w -B -I` ⇒ **12** — und es sind dieselben zwölf, die `slice-207` §2 nennt. Die in §2 ausgeschriebenen Filter sind vollständig und nachrechenbar.
- **Alle zwölf Dateien haben mindestens eine Zuordnung.** Keine Datei des Deltas fällt aus dem Audit heraus; die Reduktion „zwölf Träger" deckt die Dateimenge vollständig ab. Die Fehler liegen nicht bei den Dateien, sondern bei den Regeln *innerhalb* von zweien (F-3, F-1b) und bei der Richtung der Zuordnung (F-4).
- **R1 trägt vollständig.** Regel-Formulierung (vier Klassen, Begründung je Punkt, keine Mindestzahl, kein Sensor) deckt sich wörtlich mit `v6.5.0` · `regelwerk/modul-05-planning-harness.md`; alle fünf genannten Träger führen sie tatsächlich, `modul-06` und `modul-09` als Nachzug.
- **Die Abschnitts-Arithmetik stimmt.** Haus-Form neun Abschnitte, Baseline-Vorlage acht (§1 *Ziel und Abgrenzung* … §8 *Sub-Area-Prüfungen und Modus-Begründung*); zwei Verschmelzungen (§1+§3, §7+§8) und eine Spaltung (§6 → §4+§5) ergeben 9 − 2 + 1 = 8. Die Aussage „keine Bijektion" trägt.
- **Die Zählung des Registers trägt.** 36 Verzeichnisse über beide Kürzel (35 `BEO-ALL`, 1 `BEO-HARN`) — reproduziert; die Korrektur gegenüber dem Vorgänger-Slice ist echt. Vier der fünf genannten Zähler stimmen (11×, 3×, 7×, 15×), der fünfte nicht (F-7).
- **Die `d-check:cite`-Spannen der beiden Vorprüfungs-Blöcke sind korrekt geankert** (`modul-05` Zeilen 268–269 und 274) und erfüllen `MR-054`; `citations` bestätigt es im inneren Loop.
- **`make gates` ist grün** — zehn Gates, 689 Dateien, 0 Befunde, Coverage 94,60 %, semgrep 0 Findings. Keiner der Befunde dieses Reports ist gate-sichtbar; alle sind Urteile über Aussagen.
- **Der Beanspruchungs-Move ist strukturell korrekt** (`MR-013`): Ruhe-Marker entfernt, Pfad-Verweise aus `slice-207` und `slice-209` nachgezogen, `git` erkennt den Rename; `make planning-check` grün. Nur die Botschaft stimmt nicht (F-11).
- **Kein Kommentar- oder Zustandsfeld-Verstoß** (`AGENTS.md` §3.7) in den drei Commits: die Plan-Datei trägt kein `**Status:**`-Feld, die Risiken tragen `Ausgang: <offen>` als Zustand, keine Chronik.
- **Keine Hexagon-, Netz-, Suppressions- oder Schwellen-Frage berührt** — die drei Commits ändern ausschließlich Planungs-Dokumente; `internal/`, `cmd/`, `Makefile` und `.d-check.yml` sind unangetastet.
- **Die Sub-Area-Wahl (`*`, nicht `tools/harness/`) trägt** — der Delta hat kein Werkzeug angefasst, und die Hebung liegt in `slice-207`.

---

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 3 | F-1, F-1b, F-3 |
| MEDIUM | 7 | F-2, F-4, F-5, F-6, F-7, F-8, F-10 |
| LOW | 4 | F-9, F-11, F-12, F-13 |
| INFO | 1 | F-14 |

Gesamt **15** Findings über 15 Einträge (`F-1b` ist ein eigener Befund, kein Zusatz zu `F-1`: andere Fundstelle, anderes Versagen).

## Verdikt

**Blockierend.** DoD (1) verlangt *„zu **jeder** Regel des gemessenen Deltas eine Antwort"* und nennt eine Regel ohne Antwort ausdrücklich *„ein offener Punkt, kein stilles Übergehen"*. Gemessen fehlt eine achte Regel ganz (F-3), eine neunte und zehnte Aussage derselben Kanon-Sektion sind übersprungen (F-1b), und von den sieben gegebenen Antworten halten drei nicht: R4 ist als „übernommen" verbucht, obwohl kein lebender Träger dieses Repos sie führt (F-2); R5s einzige „Handlung am Bestand" beruht auf einer falschen Aussage über das eigene Produkt (F-1); R7 benennt einen lokalen Träger, der die Regel nicht enthält, mit einer nicht belegten Begründung (F-5).

**Die Sieben-Regel-Zerlegung trägt als Idee, nicht in ihrer Ausführung.** Die Bewegung von zwölf Dateien auf Regeln ist die richtige — sie ist der Ableiter aus `BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`, korrekt angewandt. Ausgeführt ist sie aber datei-weise: jede Datei bekam die Regel, die als ihre Hauptaussage gelesen wurde. Deshalb fehlen vier Träger (F-4), deshalb zerfällt eine Regel, die der Kanon selbst als *„Drei Formen derselben Regel"* schreibt, in R6 und R7, und deshalb bleibt der Rest der beiden 40-Zeilen-Dateien (`grundlagen-harness-dateien`, `grundlagen-traceability`) unbeantwortet.

**Der Schlusssatz hält nicht.** *„Keine Regel des Deltas ist nicht anwendbar, und keine wird abweichend adoptiert"* ist als Aussage über eine Menge formuliert, die der Audit nicht vollständig gebildet hat; sein Beleg — der Vergleich mit den beiden Vorgänger-Adoptionen — reproduziert für `slice-203` nicht (F-10). Ob die Aussage am Ende stimmt, ist offen; sie ist heute nicht getragen.

**Was ausdrücklich gut ist:** die Delta-Messung des Vorgängers ist vollständig nachrechenbar und reproduziert bis auf die Datei genau; die Register-Zählung über beide Kürzel korrigiert einen echten Vorgänger-Fehler; R6 ist am Bestand **gemessen** statt vermutet, und die Zahl stimmt; `MR-066` ist vorab und formgerecht bedient. Die Befunde betreffen ausnahmslos Aussagen, die neben gemessenen Zahlen stehen — die Zahlen selbst halten.
