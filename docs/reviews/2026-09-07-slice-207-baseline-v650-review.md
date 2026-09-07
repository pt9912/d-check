# Review-Report — slice-207 (Baseline-Pin-Hebung auf v6.5.0)

**Review-Art:** Code/Diff (gegen Slice-Plan, MR-011/021/023/039/051/055, `AGENTS.md` §3/§5/§6, Baseline-Kanon)
**Gegenstand:** `9270e4e` (Vorprüfungen) · `db7d80b` (Beanspruchung) · `4ff020e` (Pin-Hebung) · `f112185` (MR-065-Move)
**Skill:** `.harness/skills/reviewer.md` v1.13.0 @ `f112185`
**Modell-ID:** claude-opus-5[1m] · **Datum:** 2026-09-07
**Eingangs-Kontext:** der Slice-Plan; `harness/conventions/MR-067-baseline-v650.md`; MR-013/021/023/031/035/039/049/051/054/055/056/066; `harness/conventions/done/MR-065-baseline-v631.md`; `.d-check.yml`; die Vorlage `templates/harness/conventions/MR-NNN-titel.template.md` des **neuen** Baums; der Vorgänger-Report `2026-09-06-slice-202-baseline-v631-review.md`; vendorter Baum `v6.3.1` (aus `9270e4e`) gegen `v6.5.0`
**Eigene Läufe:** `make gates` (zehn Gates grün, 686 Dateien / 0 Befunde, Coverage 94,60 %) · `make baseline-verify` (ok, 54 Dateien, vollständig) · `make baseline-probe` (ok, 9 Proben) · `make adr-check RANGE=HEAD~3..HEAD` (686/0) · `make trace-check RANGE=HEAD~3..HEAD` (686/0) · `make verify-closure-notes` (587/0) · `make review-coverage` (686/0) · eigener Bruch-Test (Modul `citations` über `docs/plan/planning/done/` in einem isolierten Probe-Baum: **14 Befunde**)

**Zitier-Form.** Baseline-Stellen stehen als Tag + Pfad in Inline-Code statt als Link — die Form, die `v6.5.0` · `regelwerk/grundlagen-harness-dateien.md` §Ein einfrierendes Artefakt … für einfrierende Artefakte verlangt. Dieser Report ist eines.

---

## Findings

### F-1 · HIGH

- **quelle:** [`MR-021`](../../harness/conventions/MR-021-vendored-verweise-pin-gebunden.md) §Geltungsbereich · [`MR-039`](../../harness/conventions/MR-039-zitat-delta-im-neuen-eintrag.md) §Geltungsbereich · [`BEO-ALL/pin-bump-mirrors-ungated`](../plan/planning/observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md) (Klasse 4, Richtung *Über-Hebung*)
- **pfad:** `docs/plan/planning/done/slice-202-baseline-v631-bump.md:160,183` · `slice-203-…:154,168` · `slice-204-…:127,139` · `slice-205-…:250,260` · `slice-206-…:125,138` (zehn Direktiven, gesetzt in `4ff020e`)
- **befund:** Die zehn `d-check:cite`-Direktiven in fünf **eingefrorenen** `done/`-Slices wurden von `v6.3.1` auf `v6.5.0` umgehängt, ohne die Zeilen-Spannen mitzuziehen: sie zeigen weiter auf `modul-05-planning-harness.md:223-224` bzw. `:229-229`. Dort steht im neuen Baum *„Der Bootstrap-Modus ist Eigenschaft *pro Sub-Area*…"* bzw. *„1. **Konventionen-Dichte** …"*; die zitierten Sätze sind auf `:268-269` und `:274-274` gewandert. Vor dem Bump war jede der zehn Direktiven korrekt — der Lift hat sie falsch gemacht. Beide Regeln, die der Commit für den Lift anruft, decken diese Dateien nicht: `MR-021`s Geltungsbereich nennt „die **lebenden** Planungs-Dokumente", `MR-039` nimmt `done/` ausdrücklich aus. Der Repo-Lauf bleibt grün, weil `citations.scope.ignore` genau `docs/plan/planning/done/**` ausnimmt.
- **verifizierbar:** ja — isolierter Probe-Baum mit den sieben `done/slice-20*.md` und `.harness/baseline/v6.5.0/`, `citations.scope` ohne `ignore`: das Produkt meldet **10 × `citation-mismatch`** („Zitat-Fäule") auf genau diesen zehn Zeilen (dazu 4 × `citation-out-of-range` für die vom Vorgänger bewusst stehen gelassenen `v6.0.0`-Direktiven in `done/slice-200`/`slice-201`). Ohne Werkzeug: `sed -n '223,224p;229p;268,269p;274p' .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md`
- **klasse:** `mechanischer-lift-trifft-eingefrorene-zitat-spanne`
- **Kategorie-Begründung:** Basis wäre MEDIUM (`commit-message-overclaims-work`, erste Richtung — eine behauptete Messung fand nicht statt). Eskaliert, weil der Befund zweierlei zugleich ist: zehn **falsifizierte Lauf-Belege** im eingefrorenen Audit-Bestand, und eine **Zusage in der Gate-Konfiguration**, die das Gegenteil behauptet (F-2). Die dritte Instanz von [`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../plan/planning/observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md) (bisher 2×) — mit diesem Slice erreicht der Eintrag die Schwelle.

### F-2 · MEDIUM

- **quelle:** `AGENTS.md` §3.7 (Kommentar-Klassen) · [`BEO-ALL/begruendung-traegt-entscheidung-nicht`](../plan/planning/observations/BEO-ALL/begruendung-traegt-entscheidung-nicht/observation.md) · Wiederholung von `2026-09-06-slice-202-baseline-v631-review.md` F-7
- **pfad:** `.d-check.yml:212-221`
- **befund:** Der neue Tombstone-Kommentar trägt drei prüfbar falsche Aussagen. (a) *„`done/`-Slices … tragen ihre `d-check:cite`-Direktiven auf dem neuen Tag"* — `done/slice-200` und `done/slice-201` tragen vier Direktiven auf `v6.0.0`, also auf dem Tag von **zwei** Hebungen zuvor. (b) *„Reviews … tragen ihre `d-check:cite`-Direktiven auf dem neuen Tag"* — `docs/reviews/` enthält **null** `d-check:cite`-Direktiven; eine leere Menge wird als gehobene beschrieben. (c) *„und die zitierten Zeilen stimmen dort weiter (gemessen: …)"* — für zehn von ihnen stimmen sie nicht (F-1). Die **Entscheidung** (ein Glob für ADR-0084) trägt; die Begründung daneben trägt nicht. Das ist wortgleich die Klasse, die der Vorgänger-Report bereits an derselben Stelle meldete.
- **verifizierbar:** teilweise — (a)/(b) über `git grep -n 'd-check:cite \.' -- 'docs/plan/planning/done/*' 'docs/reviews/*'`; (c) über den Bruch-Test aus F-1. Kein Gate liest Kommentare.
- **klasse:** `gate-kommentar-nennt-falschen-grund`

### F-3 · MEDIUM

- **quelle:** `templates/harness/conventions/MR-NNN-titel.template.md` (vendored, `v6.5.0`) — *„`Löst auf` und `Ausgelöst durch Baseline-Stand` nur, wenn dieser Eintrag einen früheren ablöst"* / *„Pflicht zusammen mit „Löst auf""*
- **pfad:** `harness/conventions/MR-067-baseline-v650.md:1-10`
- **befund:** MR-067 löst MR-065 ab — der Index in `harness/conventions.md` führt MR-065 unter §Aufgelöste Adaptionen mit „aufgelöst durch [MR-067]", und MR-065s Datei wandert in `conventions/done/`. MR-067 selbst trägt weder `Löst auf:` noch `Ausgelöst durch Baseline-Stand:`. Alle vier Vorgänger der Serie (MR-057, MR-058, MR-060, MR-065) tragen beide Felder. Die Kette ist damit nur noch in einer Richtung navigierbar: von der Index-Zeile zum Nachfolger, nicht vom Nachfolger zum Vorgänger.
- **verifizierbar:** nein (kein Gate liest MR-Felder); nachzählbar über `grep -l 'Löst auf' harness/conventions/done/MR-0{57,58,60,65}-*.md harness/conventions/MR-067-*.md`
- **klasse:** `pflichtfeld-paar-des-templates-fehlt-im-neuzugang`

### F-4 · MEDIUM

- **quelle:** [`MR-039`](../../harness/conventions/MR-039-zitat-delta-im-neuen-eintrag.md) §Adaption (*„hält der **Bump-Eintrag** das fest"*) · [`MR-021`](../../harness/conventions/MR-021-vendored-verweise-pin-gebunden.md) §Geltungsbereich (*„die Menge bestimmt der Zensus der Bump-Prozedur"*)
- **pfad:** `harness/conventions/MR-067-baseline-v650.md`
- **befund:** MR-067 führt vier Rechenschaften nicht, die der unmittelbare Vorgänger MR-065 im selben Feld führte und die von aktiven Adaptionen verlangt sind: **kein Zensus** (welche und wie viele lebenden Dateien gehoben wurden), **kein `cite`-Direktiven-Konto** (die Zahlen stehen nur im Commit-Text und — falsch — im `.d-check.yml`-Kommentar), **kein MR-039-Vermerk** zum Zitat-Delta (auch die ehrliche Leermenge *„kein zitierter Wortlaut hat sich geändert"* fehlt), und **kein Adaptions-Review**. Der Review war diesmal materiell fällig: das Delta hat genau die Abschnitte umgeschrieben, die drei aktive Einträge als ihre Baseline-Regel nennen — `modul-09-implementierung.md` §Minimal Agent Workflow ([`MR-031`](../../harness/conventions/MR-031-schritt-3-benennen.md)), `grundlagen-begriffe.md` ([`MR-035`](../../harness/conventions/MR-035-cr-ablage.md)) und `modul-05-planning-harness.md` §Ziel-Form: Slice ([`MR-066`](../../harness/conventions/MR-066-slice-wachstum-ohne-rueckfuehrung.md)). Ob einer ihrer Auflösungs-Trigger gefeuert hat, beantwortet kein Artefakt des Slice.
- **verifizierbar:** nein (kein Gate); die drei Berührungen sind über `diff -w -B -I` gegen `9270e4e:.harness/baseline/v6.3.1/` nachweisbar
- **klasse:** `bump-eintrag-ohne-zensus-und-adaptions-review`

### F-5 · MEDIUM

- **quelle:** Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-06-roadmap.md` §Roadmap-Struktur · Roadmap-Entscheid des Auftraggebers 2026-08-22
- **pfad:** `docs/plan/planning/in-progress/roadmap.md:31` (gesetzt in `db7d80b`)
- **befund:** Der Beanspruchungs-Commit hat zwei Zeilen gelöscht, wo nur der Ruhe-Marker `Nichts in Arbeit.` zu löschen war. Mit ihm sind der Abschluss der Klammer `§Wann Arbeit eine Welle braucht)` und der Satz *„Parallelität ist Erlaubnis, kein Ziel."* verschwunden. §Offene Wellen endet jetzt mitten im Satz: *„Eröffnet wird weiterhin nur, was einen eigenen Closure-Grund hat (Baseline-Regelwerk `modul-06-roadmap.md`". Die Roadmap ist Rang 5 der Source Precedence; verloren ist ausgerechnet die Grenze, die der Mehr-Wellen-Entscheid sich selbst gesetzt hat. Beim nächsten Ruhe-Marker (Closure dieses Slice) wird `Nichts in Arbeit.` hinter den abgebrochenen Satz zurückgeschrieben.
- **verifizierbar:** nein — `make planning-check` prüft nur den Marker-String, `make doc-check` sieht keinen Link in der gekappten Zeile; beide sind grün. Nachweis: `git show db7d80b -- docs/plan/planning/in-progress/roadmap.md`
- **klasse:** `marker-entfernung-nimmt-nachbarsatz-mit`

### F-6 · MEDIUM

- **quelle:** `AGENTS.md` §5 (*„behauptet nicht mehr, als die Arbeit trägt"*) · [`BEO-ALL/eigene-menge-gemessen-fremde-behauptet`](../plan/planning/observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md) (9×) · Wiederholung von `2026-09-06-slice-202-baseline-v631-review.md` F-3
- **pfad:** `docs/plan/planning/in-progress/slice-207-baseline-v650-bump.md:229` (gesetzt in `9270e4e`)
- **befund:** *„Register durchgegangen (gemergter Stand, **35** Verzeichnisse — nachgezählt)."* Das Beobachtungs-Register führt **36** Verzeichnisse, sowohl auf dem Plan-Commit `9270e4e` als auch auf `HEAD`; die 36. entstand mit der Closure von slice-205 (`a0a3e61`). Die Zahl ist aus slice-206 übernommen, wo sie ebenfalls schon 35 lautete — das ausdrückliche *„nachgezählt"* ist damit die Behauptung, die nicht trägt.
- **verifizierbar:** nein (kein Gate); `find docs/plan/planning/observations -mindepth 2 -maxdepth 2 -type d | wc -l`
- **klasse:** `selbstauskunfts-zahl-weicht-von-gemessener-menge-ab`

### F-7 · MEDIUM

- **quelle:** [`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../plan/planning/observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md) · `AGENTS.md` §5
- **pfad:** `docs/plan/planning/in-progress/slice-207-baseline-v650-bump.md:53-73` (§2 *Der gemessene Delta*) und `harness/conventions/MR-067-baseline-v650.md`
- **befund:** Zwei Aussagen im selben Abschnitt widersprechen einander. §2 sagt, die größte Datei *„(`grundlagen-begriffe.md`, 81 Zeilen) war **vollständig** Tabellen-Padding"*; Punkt 7 desselben Abschnitts führt dieselbe Datei mit *„RTM im Glossar (die einzige nicht-Rauschen-Zeile dieser Datei)"*. Gemessen ist Punkt 7 richtig: die Datei überlebt `diff -w -B -I` und ist deshalb eine der zwölf, nicht eine der fünfzehn. MR-067 spiegelt die falsche Hälfte in abgeschwächter Form (*„enthielt **keine** Regel-Änderung"*) und stellt sie direkt hinter den Satz über die fünfzehn Padding-Dateien, wodurch sie sich wie eine von ihnen liest.
- **verifizierbar:** nein (kein Gate); `diff -w -B -I 'v6\.[0-9]' -I '20[0-9][0-9]-[0-9][0-9]-[0-9][0-9]'` über die beiden Bäume zeigt genau eine inhaltliche Zeile in dieser Datei
- **klasse:** `messung-widerspricht-sich-im-selben-abschnitt`

### F-8 · MEDIUM

- **quelle:** MR-067 §Adaption (*„Die Messmethode … ist der Ertrag dieser Hebung"*) · [`BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`](../plan/planning/observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
- **pfad:** `docs/plan/planning/in-progress/slice-207-baseline-v650-bump.md:56-70` · `harness/conventions/MR-067-baseline-v650.md`
- **befund:** Der erklärte Ertrag des Slice ist eine Messmethode, und die Methode ist nicht aufgeschrieben. Plan und MR nennen `diff -I` und `diff -w -B -I`, aber keine der `-I`-Regeln — die Tabellenzeile sagt nur *„(Versionen, Daten)"*. Die drei Zahlen 27/12/15 sind ohne diese Regeln nicht reproduzierbar: mit dem naheliegenden `-I '<!-- Quelle:'` (der Methode des Vorgängers, wörtlich in MR-065 protokolliert) ergeben sich 30 und 15 statt 27 und 12. Für diesen Review waren drei Anläufe nötig, bis das Regelpaar `-I 'v6\.[0-9]' -I '20[0-9][0-9]-[0-9][0-9]-[0-9][0-9]'` die berichteten Zahlen exakt traf. Damit ist die Lehre, die weitergegeben werden soll, im Artefakt nicht vollständig enthalten.
- **verifizierbar:** nein (kein Gate)
- **klasse:** `messmethode-als-ertrag-erklaert-parameter-nicht-protokolliert`

### F-9 · LOW

- **quelle:** `AGENTS.md` §3.6 (Gate-Lockerung) · Präzedenz in derselben Datei (`.d-check.yml:91-95`, `:105-113`)
- **pfad:** `.d-check.yml:221-222`
- **befund:** Der Tombstone ist quell-skopiert auf `in: docs/plan/adr/**`, obwohl der Kommentar selbst feststellt, dass **eine einzige** ADR den Verweis trägt. Für dieselbe Klasse existiert in derselben Datei die datei-skopierte Form (`in: docs/plan/adr/0022-…md`, `in: docs/plan/adr/0047-…md`). Da `.harness/baseline/v6.3.1/` jetzt entfernt ist, passiert damit jeder **künftige** ADR-Link in den toten Baum still — auch einer in einer noch nicht `Accepted`-ADR, den `links` sonst als `target-missing` meldete.
- **verifizierbar:** ja — eine Probe-ADR mit einem `v6.3.1`-Link erzeugt keinen Befund; mit `in: docs/plan/adr/0084-mentions-eigenes-modul.md` schon
- **klasse:** `tombstone-breiter-als-sein-gegenstand`

### F-10 · LOW

- **quelle:** `AGENTS.md` §3.7 (*„keine Deliberation über Verworfenes"*, *„keine Herkunfts-Prosa"*)
- **pfad:** `.d-check.yml:215-216`
- **befund:** *„§3.5 schuetzt ihren Kern, und `adr-check` hat den Versuch gemeldet"* erzählt, wie der Autor auf die Ausnahme kam — ein verworfener Zwischenschritt, kein Zustand. Der Satz kann außerdem still falsch werden: ändert sich der Prüfumfang von `adr-check`, behauptet der Kommentar weiter eine Meldung, die es nicht mehr gäbe. Der Vorgänger-Kommentar (`:200-210`) kommt ohne diesen Halbsatz aus.
- **verifizierbar:** nein (kein Gate; Urteil)
- **klasse:** `config-kommentar-traegt-deliberation`

### F-11 · LOW

- **quelle:** Maintainability
- **pfad:** `docs/plan/planning/in-progress/slice-207-baseline-v650-bump.md:42`
- **befund:** *„Der Sprung überspringt zwei Releases (`v6.4.0`, `v6.5.0`)"* — `v6.5.0` ist das Ziel des Sprungs, nicht ein übersprungenes Release; übersprungen wird genau eines. MR-067 formuliert es korrekt (*„Der Sprung übergeht `v6.4.0`"*), §5 und §6 des Plans ebenfalls (*„Zwei Releases in einem Sprung"*, *„upstream liegen zwei Releases"*). Nur diese eine Zeile behauptet, der vendorte Stand sei nicht materialisiert worden.
- **verifizierbar:** nein
- **klasse:** `spanne-und-luecke-verwechselt`

### F-12 · LOW

- **quelle:** [`MR-050`](../../harness/conventions.md#mr-050) (Herkunfts-Anker als auflösbares Feld) · Präzedenz MR-065/MR-061
- **pfad:** `harness/conventions/MR-067-baseline-v650.md` §Adaption (*„Die Liste steht im auslösenden Slice"*)
- **befund:** MR-067 ist ein **lebendes** Dokument und verweist auf „den auslösenden Slice", ohne dessen Kennung zu nennen. MR-065 nennt `slice-203` an der entsprechenden Stelle ausdrücklich, MR-061 seinen Slice ebenfalls. Wandert slice-207 nach `done/` und später ins Archiv, ist die Delta-Liste — der einzige Ort, an dem sie steht — aus MR-067 heraus nur noch über `git`-Archäologie auffindbar.
- **verifizierbar:** nein
- **klasse:** `lebender-eintrag-verweist-ohne-kennung`

### F-13 · INFO

- **quelle:** Maintainability · Präzedenz MR-065 §*Grenze der Tabelle, benannt*
- **pfad:** `docs/plan/planning/in-progress/slice-207-baseline-v650-bump.md:80-113`
- **befund:** Die Zweiteilung *(A) `v6.4.0`* / *(B) `v6.5.0`* ist eine Zuordnung je Tag, und die Zwischen-Bundles liegen nicht im Repo. Für vier der fünf (A)-Dateien trägt die in-repo abgelegte CR-Antwort den Beleg (*„Landung: Modul 5 · Modul 9 · beide Spiegel · `slice.template.md` · `templates/README.md`"*); die fünfte (`modul-06-roadmap.md`) ist eine plausible, aber unbelegte Ableitung. Der Vorgänger hat genau diese Grenze bei derselben Konstruktion ausdrücklich benannt; hier fehlt sie, und der Text verschärft stattdessen (*„die Trennlinie ist scharf"*). Der **Inhalt** beider Gruppen ist gegen den Aggregat-Diff belegt — geprüft und korrekt; nur die Tag-Zuordnung ist aus dem Repo heraus nicht vollständig nachvollziehbar.
- **verifizierbar:** nein (Netz)
- **klasse:** `per-tag-zuordnung-ohne-repo-beleg`

### F-14 · INFO

- **quelle:** Maintainability
- **pfad:** `harness/conventions/MR-067-baseline-v650.md` · `docs/plan/planning/in-progress/slice-207-baseline-v650-bump.md:71-72`
- **befund:** *„55 Dateien vorher wie nachher"* zählt den Bestand auf Platte inklusive `SHA256SUMS`; MR-065 zählte dieselbe Menge als **54** (Manifest-Einträge), und `make baseline-verify` meldet *„54 Dateien, vollständig"*. Beide Zahlen sind richtig, sie messen verschiedene Mengen — der Wechsel der Zählbasis gegenüber dem Vorgänger ist nicht markiert. Undokumentierte Annahme, kein Defekt.
- **verifizierbar:** ja — `find .harness/baseline/v6.5.0 -type f | wc -l` = 55, `make baseline-verify` = 54
- **klasse:** `zaehlbasis-still-gewechselt`

---

## Negativbefunde (geprüft, ohne Befund)

1. **Der Datei-Bestand des Bundles** — beide Bäume aus `git archive` extrahiert und die Namenslisten verglichen: **55 ↔ 55**, identische Mengen, keine neue, keine entfallene Datei. Bestätigt.
2. **`diff -rq` roh** — **35** Pfade. Bestätigt.
3. **`-I`-Stufe** — **27** Markdown-Dateien (plus `SHA256SUMS`). Bestätigt (mit dem rekonstruierten Regelpaar, siehe F-8).
4. **`-w -B -I`-Stufe** — **12** Markdown-Dateien (plus `SHA256SUMS`). Bestätigt.
5. **Die fünfzehn Weißraum-Dateien** — als Differenz der beiden Mengen einzeln aufgelistet; alle fünfzehn haben **unveränderte Zeilenzahl**, ein stiller Spannen-Versatz durch `-B` ist damit ausgeschlossen. Bestätigt.
6. **Die 81 Zeilen** — `diff -I 'v6\.[0-9]'` im klassischen Format liefert exakt **81** geänderte Zeilen (unified: 83). Zahl reproduziert.
7. **Klassifikation der zwölf** — jede Datei einzeln gegen den `-w -B -I`-Diff gehalten. Gruppe (A) fünf Dateien, alle CR-Umsetzung (§1 *Ziel und Abgrenzung* mit den vier Klassen samt der ungebetenen Schärfung *„Die Adresse muss die Sendung annehmen"*, §8-Umbenennung, `modul-09`-Pflicht, `modul-06`-Querverweis, `templates/README.md`). Gruppe (B) sieben Dateien, alle Zitier-Form bzw. RTM. **Keine Datei ist falsch einsortiert**, und keine trägt Anteile beider Gruppen — insbesondere `slice.template.md` (26 Zeilen) ist rein (A), die vier Vorlagen aus (B) tragen ausschließlich den *Zitier-Form*-Block.
8. **Vollständigkeit der Pin-Hebung, vorwärts** — `git grep '\.harness/baseline/v[0-9]'` über den ganzen Baum: **kein lebendes Dokument** nennt mehr einen Pfad in einen entfernten Baum. Der einzige `v6.3.1`-Pfad außerhalb der Tombstone-Zeile steht in `docs/plan/adr/0084-mentions-eigenes-modul.md:20` und ist gedeckt.
9. **Vollständigkeit rückwärts (Über-Hebung)** — alle übrigen `v6.3.1`-Nennungen sind reine Versions-Aussagen ohne Pfad: die beiden CR-Dokumente (*„Baseline-Stand des CR"*), `done/slice-202`/`slice-203` (Titel und Prosa), der Vorgänger-Report, MR-065 selbst und die zwei Tombstone-Kommentare. Alle bleiben zu Recht stehen. **Die einzige Über-Hebung ist F-1.**
10. **Die vier Spiegel-Klassen** — Release-Download-URL (`releases/download/v6.5.0/…`, drei Fundstellen: `AGENTS.md`, `conventions.md`, und die Vorlagen im Bundle), Tree-/Tag-URLs in `conventions.md` §Baseline und §Adoptierte Konventions-Quellen (Link-**Text** und Link-**Ziel** beide gehoben — die dritte Klasse aus MR-060 ist nicht wieder aufgetreten), Prosa-Pins (`roadmap.md` *„§Roadmap-Struktur (v6.5.0)"* — korrekt gehoben, denn §Roadmap-Struktur ist im Delta unverändert, die Aussage gilt weiter). Klasse 4 ist F-1.
11. **Symlinks** — elf Aliase unter `.claude/rules/`, davon **acht** baseline-gebunden. Alle acht zeigen auf `v6.5.0` und lösen auf; `make baseline-probe` fährt seine **neun** Proben grün. `git grep` sieht Symlink-Ziele nicht (Mode 120000) — die Klasse ist zu Recht eigens genannt.
12. **Lebende `cite`-Direktiven** — **16** baseline-gebundene in 11 Dateien, plus eine repo-interne (`reviewer.md` → `AGENTS.md:408-409`). **Vier** neu geankert (slice-207 `223-224`→`268-269` und `229-229`→`274-274`, MR-031 `modul-09:166`→`:175`, MR-035 `grundlagen-begriffe:47`→`:48`), **zwölf** unverändert, **null** entfallen. Zahl bestätigt.
13. **Geltungsbereich der Spannen** — jede der sechzehn Spannen von Hand gegen den Zieltext gehalten: alle zeigen auf die **vorschreibende** Zeile, keine auf eine Nebenregel. Die beiden neu geankerten Vorprüfungs-Belege in slice-207 treffen genau *„1. **Sub-Area-Wahl prüfen.**"* und *„2. **Offene Beobachtungen sichten.**"* — damit ist der Befund F-5 des Vorgänger-Reports (`cite-belegt-nebenregel-statt-vorschrift`) hier nicht wiederholt.
14. **Wortlaut der Zitate ohne `cite`-Direktive** — die verbatim zitierten Baseline-Sätze in MR-032 (`grundlagen-source-precedence.md`) und MR-066 (*„Dann zurück zum Schneiden (`in-progress→next`), nicht still weiterschieben."*) existieren im neuen Baum unverändert. **Kein** zitierter Wortlaut hat sich mit `v6.5.0` geändert — MR-039 hat inhaltlich nichts zu vermerken (dass der Vermerk als erklärte Leermenge fehlt, ist Teil von F-4).
15. **`Accepted`-ADR-Ausnahme** — `docs/plan/adr/0084-mentions-eigenes-modul.md:20` ist der einzige Markdown-Link einer `Accepted`-ADR in den entfernten Baum; keine weitere ADR und keine andere unantastbare Klasse trägt einen. Der Tombstone ist die richtige Auflösung: ein `d-check:ignore`-Marker oder die neue Kennung-statt-Adresse-Form hätten beide den ADR-Kern angefasst und `adr-check` rot gemacht — dieselbe Grenze, die der Commit protokolliert. Die **Breite** des Globs ist F-9.
16. **MR-065-Move (`f112185`)** — `R079`-Rename plus Link-Tiefen-Fixes plus Index-Umhängung, in der Botschaft ausdrücklich als MR-013-Ausnahme deklariert. **Neun** Verweis-Formen nachgezogen, alle neun lösen vom neuen Ort auf (geprüft mit `test -e` je Ziel). Die Index-Zeile trägt beide `<a id>`-Anker (Voll-Slug und Kurzform) und nennt MR-067 als Nachfolger; MR-065 ist aus §Aktive Adaptionen entfernt und in §Aufgelöste eingetragen. Bijektion Dateien ↔ Zeilen stimmt in beide Richtungen. *(Die Aufzählung in der Botschaft nennt acht der neun Formen — die weggelassene, `../conventions.md` → `../../conventions.md`, ist mit 14 Vorkommen die größte; die **Zahl** neun stimmt.)*
17. **MR-013 beim Beanspruchungs-Move (`db7d80b`)** — `R100`-Rename, Roadmap-Flip und die drei Pfad-Verweise in `slice-208` in **einem** Commit. Konform; der Inhalts-Defekt darin ist F-5, nicht die Bündelung.
18. **Auflösungs-Trigger der übrigen aktiven Einträge** — alle 39 durchgesehen. Keiner außer MR-065 hat mit dieser Hebung gefeuert: MR-031s Trigger verlangt, dass die Baseline die **Vorab-Nennung** selbst fordert (der neue `modul-09`-Absatz fordert etwas anderes — Out-of-Scope-Erweiterung als Plan-Änderung); MR-035/MR-036 verlangen einen Kanon-Ruheort für den ausgehenden CR (`grundlagen-begriffe.md` hat nur die RTM-Zeile bekommen); MR-066 verlangt, dass der Kanon selbst sagt, was an die Stelle der Stille tritt (der neue §1-Text regelt Scope-Creep, nicht die Ein-Sitzungs-Review-Grenze); MR-054s beide Vorprüfungs-Blöcke bestehen fort. **MR-065 ist der einzige.** Dass diese Prüfung im Slice nicht protokolliert ist, ist F-4.
19. **Slice-Plan in Haus-Form** — kein Verstoß. §3 des Plans schließt die Auflösung der Haus-Form ausdrücklich aus und adressiert sie an slice-208; die in-repo abgelegte CR-Antwort trägt denselben Entscheid (*„Die Haus-Form bleibt bis zum Bump … und wird dann aufgelöst"*), und die Reihenfolge-Bedingung ist dieselbe wie bei slice-202 → slice-203: die neuen Vorlagen liegen erst **nach** dem Bump im Repo. slice-208 existiert in `open/`, trägt den Bezug und ist damit eine Adresse, die die Sendung annimmt.
20. **Drei Vorprüfungen** — alle drei vorhanden, in der von `AGENTS.md` §5 verlangten Reihenfolge; die beiden kanonischen tragen ihre `d-check:cite`-Direktive (MR-054), die dritte bewusst keine. Die Nachtlauf-Lesung ist inhaltlich korrekt: rot ist allein die Currency-Achse.
21. **Beobachtungs-Zähler im Plan** — `pin-bump-mirrors-ungated` 5×, `semantic-change-body-only-edges-stale` 11×, `citation-stretched-beyond-scope` 15×, `zaehlmethode-misst-proxy-statt-gegenstand` 1×, `modulliste-spiegel-ungegated` 2×: alle fünf gegen die `evidence/`-Dateien nachgezählt, **alle korrekt**. Falsch ist nur die Gesamtzahl (F-6).
22. **Gate-Läufe** — `make gates` grün über zehn Gates (686 Dateien / 0 Befunde, Coverage 94,60 % ≥ 93 %), `make baseline-verify` ok (54 Dateien), `make baseline-probe` ok (9 Proben), `make adr-check RANGE=…` und `make trace-check RANGE=…` je 686/0, `make verify-closure-notes` 587/0, `make review-coverage` 686/0. Der Arbeitsbaum ist nach diesem Review unverändert.
23. **HIGH-Prüffragen der Standard-Liste** — kein Produkt-Code berührt (`internal/`, `cmd/`, `tools/` unangetastet), kein neuer Netzzugriff, keine Inline-Suppression, keine Schwellen-Senkung, keine Import-Richtung berührt. Die einzige Gate-Ausnahme ist der Tombstone (F-9).
24. **Zustandsfelder** — die berührten Kopf-/Index-Zellen (`**Stand:**` in §Baseline, die neue Index-Zeile, MR-067s Felder) nennen Zustand und Beleg, keine Chronik. Kein Drift-Log-Eintrag geschrieben, zu Recht: eine Pin-Hebung ist keine Umplanung.
25. **Pin-Träger außerhalb Markdown** — `tools/harness/fetch-baseline-cache.sh` liest den Tag aus `harness/conventions.md` §Baseline (`grep -m1 '\*\*Stand:\*\*'`); kein Workflow, kein Makefile und kein Go-Test trägt den Tag hart. `versions.pin-pattern` deckt nur `ghcr.io`-Tags — der Baseline-Pin bleibt bewusst ungewächtert.

**Nicht geprüft** (außerhalb der Rolle bzw. offline nicht möglich): die DoD-Abhakung und die Risiko-Ausgänge (Verifikation, getrennter Kontext — angemerkt sei nur, dass DoD (1) die Liste *„getrennt nach Regelwerk, Templates und reines Rauschen"* verlangt und §2 nach `v6.4.0`/`v6.5.0` gruppiert, ohne die fünfzehn Rausch-Dateien aufzuzählen); die Zuordnung der Deltas zu den beiden Zwischen-Tags (Netz, F-13); `make baseline-freshness` und `make nightly-state` (Netz); die Closure-Notiz (noch nicht geschrieben).

---

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 (F-1) |
| MEDIUM | 7 (F-2 … F-8) |
| LOW | 4 (F-9 … F-12) |
| INFO | 2 (F-13, F-14) |

## Verdikt

**Blockiert** (ein HIGH, sieben MEDIUM).

**Der Pin ist gehoben, und die mechanische Hälfte ist dicht.** Der vendorte Baum steht vollständig auf `v6.5.0`, `baseline-verify` und `baseline-probe` sind grün, die acht Aliase hängen um, kein lebendes Dokument nennt mehr einen Pfad in den entfernten Baum, kein historisches Vorkommen wurde fälschlich mitgehoben, alle sechzehn lebenden Zitat-Spannen lösen wortgleich und auf der **vorschreibenden** Zeile auf, und die vier neu geankerten sind genau die vier, die es sein mussten. **Und der gemessene Delta ist, was er zu sein behauptet:** 55/35/27/12/15 sind einzeln nachgemessen und stimmen; die Klassifikation der zwölf in fünf CR-Umsetzung und sieben Neuzugänge ist Datei für Datei geprüft und **korrekt** — keine ist falsch einsortiert.

**Was blockiert, sind drei Dinge.** Erstens F-1: der Blanket-Retarget hat nicht nur die lebenden Verweise erwischt, sondern zehn `cite`-Direktiven in fünf eingefrorenen `done/`-Slices, deren Zitate am neuen Ziel nicht mehr stehen — vor dem Bump waren alle zehn richtig, jetzt behaupten Lauf-Belege ein Kanon-Zitat, das dort nicht steht, und das einzige Gate, das es fände, nimmt genau diese Verzeichnisse aus. Zweitens die **Selbstauskunft**: der Tombstone-Kommentar sagt „gemessen", wo nicht gemessen wurde (F-2), MR-067 verzichtet auf Zensus, cite-Konto, Zitat-Delta-Vermerk und Adaptions-Review (F-4) und auf das Pflichtfeld-Paar der Vorlage (F-3), die Register-Zahl ist übernommen statt nachgezählt (F-6), und der Abschnitt über die Messung widerspricht sich selbst (F-7) und protokolliert seine eigenen Parameter nicht (F-8). Drittens F-5, der einzige Befund außerhalb der Bump-Mechanik: eine Zeile zu viel gelöscht, und §Offene Wellen endet mitten im Satz.

**Keiner der zwölf verlangt, den Bump zu wiederholen.** F-1 ist mit den korrekten Spannen (`268-269`, `274-274`) oder mit dem Rückbau auf den unangetasteten Zustand behoben — welcher der beiden Wege richtig ist, ist ein Entscheid und gehört in den Slice, denn `MR-021` und `MR-039` decken den Lift in `done/` nicht; alle übrigen sind an Text korrigierbar.

**Steering-Loop-Signale, zwei.** F-1 ist die **dritte** Instanz von `BEO-ALL/mechanical-id-rewrite-misses-frozen-classes` (bisher 2×) — mit diesem Slice erreicht der Eintrag die Schwelle und ist keine Notiz mehr, sondern eine Lücke mit eigenem Folge-Slice. Und F-2 wiederholt wörtlich den Befund F-7 des Vorgänger-Reports an derselben Datei: der Tombstone-Kommentar trifft die richtige Entscheidung und begründet sie mit einem Sachverhalt, der nicht zutrifft (`BEO-ALL/begruendung-traegt-entscheidung-nicht`, bisher 1×).
