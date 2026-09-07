# Review-Report: slice-209 — 2026-09-07

**Review-Art:** Plan/Design — geprüft werden die **Vorfrage** (DoD 1: macht die
`v6.5.0`-Unterscheidung *einfrierend / lebend* die Repo-Regel entbehrlich?) und
der daraus abgeleitete **Konventions-Eintrag** `MR-070` gegen den Kanon, gegen
die Vorlage und gegen den eigenen Bestand. **Nicht** geprüft: die DoD-Abhakung
und die Closure-Notiz (§9 ist leer, Verifier-Sache).

**Gegenstand:** slice-209 · Diff-Range `b2b4acf..3ce7b7c` (vier Commits:
`84e30e4`, `b672348`, `f37af1e`, `3ce7b7c`)

**Skill:** `.harness/skills/reviewer.md` @ 1.14.0 · <!-- d-check:ignore (Adopter-spezifischer Skill-Pfad, existiert im Ziel-Repo ggf. nicht) -->
**Modell:** claude-opus-5[1m] · **Datum:** 2026-09-07

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v6.5.0` · `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- der Slice-Plan slice-209 (§1–§8), die vier Commits der Range und ihre Botschaften
- `harness/conventions/MR-070-frozen-klassen-vor-mechanischer-ersetzung.md`,
  die Index-Zeile in `harness/conventions.md`, die `state.md` von
  `BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`
- `AGENTS.md` §1, §3.1–§3.9, §4, §5, §6; `harness/README.md` §Sensors
- `MR-013`, `MR-025`, `MR-049`, `MR-051`, `MR-052`, `MR-053`, `MR-054`,
  `MR-066`, `MR-067`, `MR-069`, `MR-000`
- die Vorlage `v6.5.0` · `templates/harness/conventions/MR-NNN-titel.template.md`
  und `v6.5.0` · `templates/docs/plan/planning/archiv-stub-slice.template.md`
- Kanon: `v6.5.0` · `regelwerk/grundlagen-harness-dateien.md`
  §harness/README.md als Einstiegspunkt (Absatz *„Ein einfrierendes Artefakt
  nennt ein prozess-bewegtes bei seiner Kennung"*, Zeilen 293–326) ·
  `v6.5.0` · `regelwerk/grundlagen-source-precedence.md` §Fork ·
  `v6.5.0` · `regelwerk/modul-05-planning-harness.md` §Zwei Schritte vor der
  Modus-Begründung · `v6.5.0` · `regelwerk/modul-06-roadmap.md`
  §Das Beobachtungs-Register
- die drei Evidence-Dateien des Registereintrags (slice-195, slice-202,
  slice-207), die Review-Reports zu slice-202, slice-207 und slice-208,
  `.d-check.yml` (`citations.scope`, `ignore-refs`, `mentions`)
- **Vorherige Findings am gleichen Gegenstand:** die Runde zu slice-207
  (F-1 = die Über-Hebung, die der dritte Beleg dieses Registereintrags ist)
  und die zwei Runden zu slice-208 (F-4 zu `MR-068`, F-7/F-8 zu `MR-069`)

**Eigene Läufe:** `make gates` (Exit 0 — *„[gates] baseline-verify +
workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep
+ gate-consistency + planning-check green"*, `d-check: 702 Datei(en) geprüft, 0
Befund(e)`, `coverage-gate: OK — Coverage 94.60% erfüllt Schwelle 93%`, semgrep
*„Ran 55 rules on 63 files: 0 findings"*) · `make mention-coverage` (Exit 0,
*„84 von 84 Artefakt(en) erwähnt, über 1 Dokument(e)"*).

---

## Findings

### F-1 · MEDIUM · Beleg 1 nennt die `d-check:cite`-Direktive die vom Kanon vorgeschriebene Form — sie ist es nicht, und die Umkehr-These ist am eigenen Bestand widerlegt

- **kategorie:** MEDIUM
- **quelle:** `v6.5.0` · `regelwerk/grundlagen-harness-dateien.md`
  §harness/README.md als Einstiegspunkt (dritte Form) ·
  [`AGENTS.md` §5](../../AGENTS.md#5-dokumentations-regeln) (*„Eine zitierte
  Quelle trägt nur, was in ihrem Geltungsbereich steht"*) ·
  `BEO-ALL/citation-stretched-beyond-scope`
- **pfad:** slice-209 §2, Absatz *Beleg 1* und Absatz *Die Gegenprobe* ·
  `harness/conventions/MR-070-frozen-klassen-vor-mechanischer-ersetzung.md`
  §Adaption, fünfte Grenze (*„Sie deckt nicht, was der Kanon selbst
  mitbringt"*) · Botschaften `f37af1e` und `3ce7b7c`
- **befund:** Der Plan schreibt: *„Die dritte Form des Kanons lautet: ‚Eine
  Stelle der vendored Baseline heißt Tag **und** Pfad in Inline-Code, nicht als
  Link.' Genau diese Form haben die zehn `d-check:cite`-Direktiven … Sie waren
  **konform** und wurden beschädigt."* Die zehn gehobenen Vorkommen lauten
  gemessen `<!-- d-check:cite .harness/baseline/<tag>/regelwerk/<datei>.md:NNN-NNN -->`
  — HTML-Direktiven mit **vollem Vendoring-Pfad und Zeilenspanne**, weder
  Inline-Code noch Abschnitts-Anker. Sie sind nicht die Kanon-Form, sondern die
  Adress-Form, gegen die der Absatz gebaut ist; in diesem Repo regiert sie
  `MR-051` (Geltungsbereich: *lebende* Dokumente) und nicht der zitierte
  Kanon-Satz. Die daraus gezogene Umkehr-These — *„die vom Kanon
  **vorgeschriebene** Form ist genau das, was eine Tag-Ersetzung greift"* — ist
  am Bestand widerlegt: Die vorgeschriebene Form liegt vor (**15** Vorkommen in
  drei eingefrorenen Reports unter `docs/reviews/`, Muster
  `` `v6.5.0` · `regelwerk/<datei>.md` §<Abschnitt> ``) und trägt den Präfix
  `.harness/baseline/<tag>/` gerade **nicht**, während slice-207s Ersetzung
  laut eigener Evidence-Datei auf genau diesen Präfix skopiert war. Die
  vorgeschriebene Form wäre also unberührt geblieben; gegriffen wurde die
  nicht-vorgeschriebene. Das Nein der Vorfrage trägt weiter (Beleg 2 und das
  Adressat-/Zeitpunkt-Argument sind davon unabhängig) — was nicht trägt, ist
  der als *„der harte"* ausgewiesene Beleg und die Grenze, die daraus im
  Eintrag geworden ist.
- **verifizierbar:** ja — `git show 92390db6` zeigt die zehn Direktiven in
  ihrer Form; eine `grep -rE`-Suche nach dem Muster
  ``  `v<X.Y.Z>` · `regelwerk ``  über `docs/reviews/` zeigt die 15 Vorkommen
  der Kanon-Form ohne Vendoring-Präfix. **Kein Gate hält das** —
  `citations.scope` nimmt `docs/plan/planning/done/**` und `docs/reviews/**`
  aus.
- **klasse:** `beleg-traegt-andere-form-als-behauptet`

### F-2 · MEDIUM · `MR-070` kollidiert mit `AGENTS.md` §3.3/`MR-013`, und die Kollision ist nicht benannt — der eigene Beanspruchungs-Commit verletzt die Regel, die derselbe Slice schreibt

- **kategorie:** MEDIUM (Kontext-Eskalation erwogen: die literale Befolgung
  färbt `make doc-check` rot, also den Gate-Pfad; als Regel-Kollision ohne
  eigenen Sensor bleibt sie hier MEDIUM)
- **quelle:** [`AGENTS.md` §3.3](../../AGENTS.md#33-git-mv--inhaltsänderung--zwei-commits)
  (Ausnahme Beanspruchung) · [`MR-013`](../../harness/conventions.md#mr-013) ·
  [`AGENTS.md` §1](../../AGENTS.md#1-was-diese-datei-ist) (*„Melde den
  Widerspruch, statt ihn stillschweigend nach einer Seite aufzulösen"*)
- **pfad:** `harness/conventions/MR-070-frozen-klassen-vor-mechanischer-ersetzung.md`
  §Geltungsbereich und §Adaption, Tabellenzeile 1 · Commit `b672348`
- **befund:** Der Geltungsbereich erfasst *„jede mechanische Ersetzung über
  mehr als eine Datei — `sed`, **`git mv` mit Nachzug**, ein Werkzeug-Lauf, der
  Token austauscht"*, und Tabellenzeile 1 nimmt `docs/plan/planning/done/` als
  *„Lauf-Beleg zu seinem Datum"* aus. `AGENTS.md` §3.3 verlangt beim
  Beanspruchungs-Move das genaue Gegenteil: *„die Pfad-Verweise auf den Slice
  wandern von `open/` nach `in-progress/`"* — auch in `done/`-Dateien. Der
  Commit `b672348` dieses Slice hat das an **sechs** Stellen in
  `done/slice-207` und `done/slice-208` ausgeführt, also die Regel gebrochen,
  die `3ce7b7c` drei Commits später schreibt. Wer künftig Zeile 1 befolgt,
  lässt die Verweise stehen; `links` meldet dann `target-missing` und
  `make doc-check` wird rot. Weder der Eintrag noch der Plan nennt diesen
  Konflikt, obwohl §1 des Briefings genau dafür die Meldepflicht setzt und
  `MR-070` fünf andere Grenzen ausschreibt.
- **verifizierbar:** ja — `git show b672348` zeigt die sechs Nachzüge in
  `done/`; das Zurücknehmen einer der sechs Zeilen und `make doc-check` zeigt
  den `target-missing`-Ausgang. Kein Gate hält die Kollision selbst.
- **klasse:** `neue-regel-kollidiert-mit-hard-rule-unbenannt`

### F-3 · MEDIUM · Die Tabellenzeile `harness/conventions/done/` ordnet eine Kanon-Klasse falsch zu; die Zählung „zwei Zeilen nennt der Kanon nicht" ist damit drei

- **kategorie:** MEDIUM
- **quelle:** `v6.5.0` · `regelwerk/grundlagen-harness-dateien.md`
  §harness/README.md als Einstiegspunkt (die fünf Klassen) ·
  `v6.5.0` · `templates/docs/plan/planning/archiv-stub-slice.template.md` ·
  `v6.5.0` · `templates/harness/conventions/MR-NNN-titel.template.md`
- **pfad:** `harness/conventions/MR-070-frozen-klassen-vor-mechanischer-ersetzung.md`
  §Adaption, Tabellenzeile 3 (*„`harness/conventions/done/` | … | ja (als
  Archiv-Stub-Klasse)"*) · Botschaft `3ce7b7c` (*„Zwei Zeilen der Tabelle nennt
  der Kanon nicht"*)
- **befund:** Der Kanon nennt fünf einfrierende Klassen — Review-Report,
  Closure-Notiz, **Archiv-Stub**, `Accepted`-ADR, geschlossener Slice. Ein
  Archiv-Stub ist der gekürzte Platzhalter aus `v6.5.0` ·
  `regelwerk/modul-06-roadmap.md` §Wellen-Closure-Prozedur Schritt 4, mit
  eigener Vorlage; ein aufgelöster `MR`-Eintrag ist kein Stub, sondern der
  vollständige Text an einem anderen Ort. Der Kanon führt ihn im Absatz davor
  ausdrücklich als *aufbewahrt* (*„Ein aufgelöster `MR`-Eintrag erklärt weiter,
  welche Form das Repo einmal hatte, und wird deshalb aufbewahrt"*) — nicht als
  einfrierend —, und die `MR`-Vorlage verlangt nach dem `git mv` nach `done/`
  ausdrücklich einen **Pfad-Berichtigungs-Commit**, also genau das Nachziehen,
  das die Kanon-Definition ausschließt. Die Spalte müsste *nein* tragen; damit
  ist auch die Zahl in der Commit-Botschaft um eins zu niedrig. Die Zeile trägt
  im Bestand: `conventions/done/MR-057` und `MR-058` stehen weiter auf
  `v5.18.0` — die *Praxis* dieses Repos ist richtig, ihre *Herleitung* aus dem
  Kanon ist es nicht.
- **verifizierbar:** ja — Textvergleich der fünf Klassen gegen die Zeile;
  `grep -rn "d-check:cite" harness/conventions/done/` zeigt die Bestandspraxis.
  Kein Gate hält das.
- **klasse:** `kanon-klasse-falsch-zugeordnet`

### F-4 · MEDIUM · Beleg 3 unterstellt dem Kanon eine reine Aufzählung, leitet den Gegenbeleg aber aus der Eigenschaft ab, die im selben Kanon-Satz steht

- **kategorie:** MEDIUM
- **quelle:** `v6.5.0` · `regelwerk/grundlagen-harness-dateien.md`
  §harness/README.md als Einstiegspunkt ·
  `BEO-ALL/citation-stretched-beyond-scope` (15×, vom Slice selbst als der
  einschlägigste Eintrag gesichtet)
- **pfad:** slice-209 §2, Absatz *Beleg 3*
- **befund:** Der Plan überschreibt den Absatz mit *„die Kanon-Aufzählung ist
  für dieses Repo unvollständig"* und schließt mit *„Das ist wörtlich der Punkt
  dieser Regel: die Liste über die Eigenschaft, nicht über eine Aufzählung."*
  Der Kanon-Satz lautet jedoch: *„Einfrierend sind die **Zeitdokumente** —
  Review-Report, Closure-Notiz, Archiv-Stub, `Accepted`-ADR, geschlossener
  Slice: Sie halten eine Messung oder Entscheidung zu ihrem Datum fest und
  werden nicht nachgezogen."* Die Eigenschaft steht dort, und der Absatz nutzt
  sie selbst — *„Ein gesendeter CR ist einfrierend **nach derselben
  Eigenschaft**"*. Damit belegt Beleg 3 die Aussage nicht, die er tragen soll:
  Er zeigt, dass die **Aufzählung** den CR nicht nennt, nicht, dass der Kanon
  ohne Eigenschaft arbeitet — und der Test in `MR-070` (*„Würde ein
  korrigierter Wert verfälschen, was dieses Artefakt zu seinem Datum
  festgehalten hat?"*) ist eine Umformulierung derselben Eigenschaft. Das Nein
  der Vorfrage bleibt richtig; sein drittes Standbein trägt nicht, was der Plan
  ihm zuschreibt.
- **verifizierbar:** ja — Textvergleich des Kanon-Satzes gegen den Absatz. Kein
  Gate hält das.
- **klasse:** `kanon-enger-gelesen-als-er-ist`

### F-5 · MEDIUM · „Beim zweiten (Pin-Hebung) dasselbe Muster" widerspricht der eigenen Evidence-Datei und vier Absätze später dem eigenen Beleg 2

- **kategorie:** MEDIUM
- **quelle:** [`AGENTS.md` §5](../../AGENTS.md#5-dokumentations-regeln)
  (*„behauptet nicht mehr, als die Arbeit trägt"*) · die Evidence-Datei
  `evidence/slice-202.md` des Registereintrags
- **pfad:** slice-209 §2, Absatz *Was die drei Anlässe gemeinsam haben*, Satz 3
- **befund:** Der Satz steht direkt hinter der Beschreibung des ersten Anlasses
  (*„blieben `Accepted`-ADR-Kerne und gesendete CRs unbedacht, obwohl die drei
  benannten Frozen-Verzeichnisse korrekt ausgenommen waren"*) und überträgt sie
  auf den zweiten. Beides ist für slice-202 falsch: Die Evidence-Datei
  beschreibt eine *„identifizierende Nennung in einem lebenden Dokument"* und
  keine ADR-Kerne oder CRs, und slice-202s Plan führt überhaupt keine
  Ausschluss-Liste (`grep -n "Ausschluss\|ausgenommen\|Accepted\|docs/plan/cr"`
  über die Datei: null Treffer). Vier Absätze weiter sagt Beleg 2 das Gegenteil
  — die Instanz sei *„außerhalb des Kanon-Geltungsbereichs"*, also kategorial
  anders. Der Plan wandert bei der Closure nach `done/` und friert als
  Lauf-Beleg ein; ein Leser, der nur §2 Satz 3 liest, hält die drei Anlässe für
  gleichartig und damit die Regel für breiter belegt, als sie ist. `MR-070`s
  §Begründung ist an derselben Stelle korrekt — sie überspringt den zweiten
  Anlass; die Differenz sitzt allein im Plan.
- **verifizierbar:** ja — Vergleich des Satzes gegen `evidence/slice-202.md`
  und gegen Beleg 2 im selben Abschnitt. Kein Gate hält das.
- **klasse:** `selbstauskunft-widerspricht-eigenem-beleg`

### F-6 · MEDIUM · Die Vorfrage ist nur gegen den Kanon gestellt, obwohl §8 dieselbe Frage gegen `MR-069` aufwirft — und `MR-069` trägt bereits eine gemessene Klassenliste mit derselben Eigenschaft

- **kategorie:** MEDIUM
- **quelle:** [`MR-069`](../../harness/conventions.md#mr-069) §Adaption ·
  `v6.5.0` · `regelwerk/grundlagen-source-precedence.md` §Source Precedence
  (Kopien driften) · [`AGENTS.md` §5](../../AGENTS.md#5-dokumentations-regeln)
  (*„die **direkteste** Quelle wählen"*)
- **pfad:** slice-209 §2 (*Die Vorfrage, beantwortet* — vergleicht
  ausschließlich mit dem Kanon), §8 Bullet *Konventions-Dichte*, §5
  (Risiko-Liste) · `harness/conventions/MR-070-…md` §Adaption (Tabelle)
- **befund:** §8 benennt das Risiko wörtlich: *„Genau diese Dichte ist das
  Risiko — ein vierter Eintrag daneben muss sagen, was die drei nicht sagen,
  sonst ist er die zweite Quelle, vor der die Source-Precedence warnt."* Die
  Vorfrage in §2 beantwortet das nur gegen den **Kanon**; §5 führt das Risiko
  nicht, und `MR-070` nennt `MR-069` an keiner Stelle. `MR-069` — am selben Tag
  `Accepted` — misst über dasselbe Gelände: *„**25** Einträge über **zehn**
  entfernte Tags, dahinter **28** Dateien — 18 `Accepted`-ADRs, sechs
  aufgelöste Konventions-Einträge, drei `done/`-Slices und ein gesendeter CR"*,
  und formuliert die Eigenschaft: *„Alle vier Klassen zitieren den Stand ihrer
  Zeit."* `MR-070` legt daneben eine zweite Partition über sechs Zeilen, die
  `harness/conventions/done/` anders einordnet (F-3) und eine Klasse ergänzt,
  die `MR-069` nicht führt. Failure: Wird eine Klasse künftig in einem der
  beiden Einträge nachgetragen, driften die Listen; wer nur eine liest,
  übersieht eine Klasse — genau der Zustand, den §8 vorab benannt und die
  Vorfrage nicht geprüft hat.
- **verifizierbar:** ja — Textvergleich der beiden §Adaption-Abschnitte;
  `grep -n "MR-069" harness/conventions/MR-070-*.md` liefert null Treffer. Kein
  Gate hält das.
- **klasse:** `zweite-quelle-neben-naeherem-eintrag`

### F-7 · LOW · Die `state.md` kopiert `MR-070`s vierte Grenze, statt auf sie zu zeigen

- **kategorie:** LOW
- **quelle:** `v6.5.0` · `regelwerk/modul-06-roadmap.md` §Das
  Beobachtungs-Register (*„`state.md` … Zustand und Beleg als auflösbarer
  Anker … keine Chronik"*) · `BEO-ALL/semantic-change-body-only-edges-stale`
  (12×) · [`MR-025`](../../harness/conventions.md#mr-025)
- **pfad:** `docs/plan/planning/observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/state.md`
- **befund:** Zustand und Beleg stehen korrekt und auflösbar (*„verkörpert —
  `MR-070` (`seit slice-209`)"*). Der Zusatz *„Gefangen wurde der Fehler
  dreimal vom `git diff` danach oder vom unabhängigen Review, nie von einer
  Liste davor"* wiederholt jedoch sinngleich die vierte Grenze aus `MR-070`
  §Adaption. Wird sie dort revidiert — und F-1 gibt dafür einen Anlass —,
  bleibt die Registerzeile stehen und behauptet weiter, was der Eintrag nicht
  mehr sagt; das Register wird bei **jeder** Slice-Planung gelesen, der Eintrag
  nur bei Bedarf. Der Bestand kennt die Form (`large-migration-exceeds-session-review-limit`,
  `zaehlmethode-misst-proxy-statt-gegenstand`) — gemeldet als Klasse, die mit
  diesem Slice wächst, nicht als Neuerung dieses Slice.
- **verifizierbar:** nein (kein Sensor auf Register-Zeilen-Inhalt; Textvergleich
  gegen `MR-070` §Adaption, vierte Grenze).
- **klasse:** `registerzeile-kopiert-statt-zeigt`

### F-8 · INFO · `MR-070` erörtert die Zitier-Form des Kanons, ohne zu nennen, dass sein eigener Baseline-Verweis ein Link ist und nach seiner Auflösung `MR-069`s Ventil speist

- **kategorie:** INFO
- **quelle:** [`MR-069`](../../harness/conventions.md#mr-069) §Grenze (1)
  (*„Jede Hebung, die eine Adresse in einem einfrierenden Artefakt zurücklässt,
  braucht einen weiteren Eintrag"*)
- **pfad:** `harness/conventions/MR-070-frozen-klassen-vor-mechanischer-ersetzung.md`
  §Ersetzt-Baseline-Regel (Markdown-Link in
  `.harness/baseline/v6.5.0/regelwerk/…`) gegen §Adaption, fünfte Grenze
- **befund:** Der Link ist an dieser Stelle **korrekt** — die Vorlage verlangt
  ihn, und ein Eintrag in `harness/conventions/` ist ein lebendes Artefakt.
  Nach der Auflösung wandert er nach `conventions/done/` und wird eingefroren;
  `MR-069` misst bereits sechs solche Einträge hinter dem `ignore-refs`-Ventil.
  Der Eintrag diskutiert in seiner fünften Grenze genau die Zitier-Form und
  lässt diese Anschluss-Zusage unerwähnt. Undokumentierte Annahme, kein Mangel.
- **verifizierbar:** nein (Textvergleich; die Wirkung tritt erst beim übernächsten
  Bump ein).
- **klasse:** `eigener-verweis-speist-das-benannte-ventil`

---

## Negativbefunde (geprüft, ohne Befund)

- **`make gates`** — gefahren, Exit 0. Zehn Gates, `d-check: 702 Datei(en)
  geprüft, 0 Befund(e)`, `coverage-gate: OK — Coverage 94.60% erfüllt Schwelle
  93%`, semgrep *„Ran 55 rules on 63 files: 0 findings"*. Die Zahl in der
  Botschaft `3ce7b7c` (*„zehn Gates, 702 Dateien, 0 Befunde"*) reproduziert
  exakt; die 701 der drei früheren Botschaften sind mit dem einen neuen
  `MR-070`-File konsistent.
- **`make mention-coverage`** — gefahren, Exit 0, *„84 von 84 Artefakt(en)
  erwähnt"*. Der Block skopiert auf ADRs; ein neuer `MR`-Eintrag fällt nicht
  darunter, also kein verdeckter roter Fokus-Lauf.
- **Selbstauskunfts-Zahlen des Plans** — alle reproduziert: **37**
  Beobachtungs-Verzeichnisse (36 unter `BEO-ALL`, 1 unter `BEO-HARN`);
  **fünf** einschlägige Einträge aufgeführt; die Zähler **3×** (Gegenstand),
  **15×** (`citation-stretched-beyond-scope`), **2×**
  (`registerzeile-ohne-ausgang-nach-schwelle`), **7×**
  (`rule-drawn-from-occasion-not-inventory`), **12×**
  (`semantic-change-body-only-edges-stale`) und **3×**
  (`zaehlmethode-misst-proxy-statt-gegenstand`) entsprechen der Zahl der
  Evidence-Dateien. *„Keiner der fünf erreicht mit diesem Slice die Schwelle
  erstmalig"* trifft zu — der Slice legt keine Evidence-Datei an.
- **„zehn Direktiven in fünf Dateien"** — reproduziert: `git show 92390db6`
  nimmt je zwei `d-check:cite`-Direktiven in `done/slice-202` … `done/slice-206`
  zurück, dazu die zwei Prosa-Pfade, die die Evidence-Datei nennt. Die *Zahl*
  trägt; ihre *Einordnung* ist F-1.
- **Beleg 2** — geprüft und tragend: `bd585ee4` änderte die betroffene Zeile in
  `docs/plan/planning/in-progress/slice-202-…md`, also in einem lebenden
  Dokument im Sinne des Kanons (*„Ein *lebendes* — `AGENTS.md`, eine Spec, ein
  offener Slice — verlinkt weiter"*); der Review-Report zu slice-202 vom
  2026-09-06 benennt dieselbe Klasse. Der Kanon nimmt lebende Artefakte
  ausdrücklich aus, diese Instanz liegt also wirklich außerhalb seines
  Geltungsbereichs. Das ist der Beleg, auf dem das Nein der Vorfrage
  eigenständig steht.
- **Die Vorfrage als Urteil** — geprüft in beide Richtungen. Das Nein trägt:
  Kanon-Adressat ist der Autor im Moment des Schreibens und der Gegenstand die
  **Form eines Verweises**; `MR-070`-Adressat ist der Ausführende einer
  Massen-Operation und der Gegenstand ihre **Ausschluss-Menge**. Die
  Gegenrichtung (hat der Implementer den Kanon zu eng gelesen, um seine Regel
  zu retten?) ist an zwei Stellen belegt eingetreten und dort gemeldet (F-1,
  F-4); sie kippt das Ergebnis nicht, weil Beleg 2 und das
  Adressat-Argument unabhängig tragen.
- **`d-check:cite`-Direktiven in §7** — geprüft: `modul-05-planning-harness.md:268-269`
  und `:274-274` sind die **vorschreibenden** Zeilen der beiden kanonischen
  Vorprüfungen (Zeile 268 = *„Sub-Area-Wahl prüfen"*, Zeile 274 = *„Offene
  Beobachtungen sichten"*), die Zitate sind wortgleich — `citations` läuft
  fail-closed im inneren Loop und ist grün. Der dritte Block (Nachtlauf) trägt
  bewusst keine Direktive, wie `MR-053`/`MR-054` es vorsehen.
- **Anker-Paarung** — geprüft: `harness/conventions.md#mr-070` löst auf (die
  Index-Zeile trägt `<a id="mr-070">` **und** den Voll-Slug-Anker), und der
  Zielort trägt den Herkunfts-Anker (`**Herkunft:** seit slice-209` im Eintrag,
  `(seit slice-209)` in der Index-Zeile).
- **`MR-070` gegen die Vorlage** — geprüft: alle sechs Pflichtfelder (Datum,
  Geltungsbereich, Ersetzt-Baseline-Regel, Adaption, Begründung,
  Auflösungs-Trigger) sind vorhanden; `Löst auf` und `Ausgelöst durch
  Baseline-Stand` fehlen zu Recht, weil kein Vorgänger abgelöst wird. Der
  Abschnitts-Anker des Baseline-Links (`#harnessreadmemd-als-einstiegspunkt`)
  ist korrekt: Der zitierte Absatz steht bei Zeile 293 zwischen den
  Überschriften der Zeilen 164 und 328.
- **Fork-Frage** — geprüft, **kein** Fork. Die Vorlage sagt verkürzt *„Ein
  Eintrag, der keine benannte Regel ersetzt, ist ein **Fork**"*; die direktere
  Quelle `v6.5.0` · `regelwerk/grundlagen-source-precedence.md` schneidet enger:
  Fork ist, wer *„die Baseline **pauschal für nicht anwendbar** erklärt"*.
  `MR-070` erklärt nichts für unanwendbar, benennt die nächstgelegene Regel und
  grenzt gegen sie ab; **23 der 41** aktiven Einträge tragen dieselbe
  „keine — …"-Form. Kein Befund.
- **`AGENTS.md` §3.3** — geprüft: `b672348` ist ein reiner `git mv` plus
  Roadmap-Flip plus gekoppelte Pfad-Verweise, wie `MR-013` es für die
  Beanspruchung verlangt; der Slice-Body bleibt unangetastet, die
  Rename-Detection hält. Der Marker-Entfernung fehlt nichts — genau zwei
  Zeilen, `planning-check` grün; der Zwei-Zeilen-Fehler aus slice-207 (F-5)
  wiederholt sich **nicht**. Der Befund F-2 betrifft die Kollision der neuen
  Regel mit dieser Ausnahme, nicht die Form des Commits.
- **`AGENTS.md` §3.1** — geprüft: kein Host-Go, kein Host-Interpreter; die
  Range enthält ausschließlich Markdown, alle genannten Messungen laufen über
  `make`.
- **`AGENTS.md` §3.2 / §3.4 / §3.5 / §3.6 / §3.9** — geprüft, nicht berührt:
  keine Suppression, kein Spec-Stratum, keine ADR, keine Schwelle, keine
  Workflow-Datei im Diff.
- **`AGENTS.md` §3.7, Kommentar-Klassen** — geprüft: die Range enthält keinen
  Code-, Konfigurations- oder Skript-Kommentar. Das Zustandsfeld `**Stand:**`
  nennt Zustand und Beleg voran; der Chronik-Verdacht am Nachsatz ist als F-7
  unter der Kopie-Achse geführt, nicht als §3.7-Verstoß — der Nachsatz stützt
  eine gegenwärtige Rest-Risiko-Aussage und hat Bestands-Präzedenz.
- **`AGENTS.md` §3.8** — geprüft, nicht anwendbar: der Slice baut kein Modul
  und ändert keine Scan-Menge; §3 schließt einen Sensor ausdrücklich aus.
- **Commit-Botschaften, Traceability und Reichweite** — geprüft: alle vier
  tragen eine Kennung (`slice-209` plus `MR-013`/`MR-053`/`MR-054`/`MR-070`);
  jede genannte Probe ist gelaufen und reproduziert. Die einzige Überdehnung
  ist die Beleg-1-Aussage in `f37af1e`/`3ce7b7c` und dort als F-1 geführt statt
  hier doppelt.
- **`MR-068`** — geprüft: die Nummer ist nirgends belegt und wird nicht
  nachbesetzt; die Botschaft `3ce7b7c` sagt das korrekt.
- **Folge-Slice** — geprüft: `slice-210`, den §7 als Ausgang von
  `zaehlmethode-misst-proxy-statt-gegenstand` nennt, liegt in `open/`.
- **Nachtlauf-Block (dritte Vorprüfung, `MR-053`)** — nicht nachgefahren
  (Netz-Target, fail-open, außerhalb `gates`); die Aussage ist eine
  Zeitpunkt-Messung, die ein Review nicht reproduzieren kann. Als Nicht-Prüfung
  benannt statt als Negativbefund behauptet.

---

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
| --- | --- | --- |
| HIGH | 0 | — |
| MEDIUM | 6 | F-1 · F-2 · F-3 · F-4 · F-5 · F-6 |
| LOW | 1 | F-7 |
| INFO | 1 | F-8 |

Wiederkehrende Klassen dieser Sitzung: **drei** Findings (F-1, F-3, F-4) sind
Ausprägungen von `BEO-ALL/citation-stretched-beyond-scope` — eine Quelle stützt
eine Aussage, die sie nicht trägt, hier zweimal durch **zu enges** und einmal
durch **zu weites** Lesen derselben Kanon-Stelle. Nach der
Kontext-Eskalations-Regel des Skills (*„die dritte Wiederholung derselben
Klasse in einer Sitzung ist ein Steering-Loop-Signal"*) ist das ein Signal für
die Closure: Der Eintrag steht bei 15× und ist von diesem Slice vorab gesichtet
worden — der vierte Beleg entsteht in genau dem Slice, der ihn als *„den
einschlägigsten für die Vorfrage"* geführt hat.

---

## Verdikt

**Blockierend.** Sechs MEDIUM, kein HIGH. Das **Ergebnis** der Vorfrage — das
Nein — trägt und ist nach eigener Prüfung richtig: Beleg 2 (lebendes Dokument,
vom Kanon ausdrücklich ausgenommen) und die Adressat-/Zeitpunkt-/Handlungs-
Differenz stehen unabhängig. Was nicht trägt, ist die Hälfte der
**Begründung**, auf die sich `MR-070` ausdrücklich beruft (*„ruht auf drei
Belegen, nicht auf einer Behauptung"*): Beleg 1 misst eine andere Form als die
behauptete und ist in seiner Umkehr-These am eigenen Bestand widerlegt (F-1),
Beleg 3 leitet aus der Kanon-Eigenschaft ab, deren Fehlen er behauptet (F-4).
Dazu kommen eine unbenannte Kollision mit `AGENTS.md` §3.3/`MR-013`, die der
eigene Beanspruchungs-Commit bereits vorführt (F-2), eine falsch zugeordnete
Kanon-Klasse samt zu niedriger Zahl in der Botschaft (F-3), ein
Selbstwiderspruch in §2 (F-5) und die nur halb gestellte Vorfrage gegenüber
`MR-069` (F-6). Alle sechs sind Doku-/Regel-Befunde ohne Gate-Deckung; `make
gates` ist grün und bleibt es, was den Punkt eher schärft als entschärft.
