# Verifikation — slice-291: `vcs` folgt einer Umbenennung, Anforderung, Spezifikation und ADR

- **Rolle:** Verifier (Modul 11) — DoD, Plan und Spec gegen den gelieferten Text. Die Reviews R1–R4
  dienten nur als Kontext; Zuordnungen und Zählungen sind selbst gemacht.
- **Gegenstand:** Plan
  [`slice-291`](../plan/planning/done/slice-291-vcs-umbenennung-anforderung-und-adr.md) §1, §2,
  §3 (Entscheide 1–14, Plan-Nachtrag, Plan-Änderungen nach R1–R3) und §6;
  [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  (Lastenheft 0.114.3),
  §[`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
  Schritt 2 und 4 samt Schema-Zeile, [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md).
- **Range:** `c0959fff..1c571ac9` auf `main`, lokal, 32 Commits. Der Diff berührt nur `spec/lastenheft.md`,
  `spec/spezifikation.md`, die ADR, drei Planungs-Dokumente und die vier Review-Reports — kein Code, keine
  `.d-check.yml`, kein CHANGELOG, kein README, kein Handbuch (Abgrenzung §1 gehalten).
- **Datum:** 2026-11-24.

## Gates am Repo (selbst gefahren, Spitze `1c571ac9`, Arbeitsbaum sauber)

| Kommando | Ergebnis |
|---|---|
| `make gates` | Exit 0. Letzte Zeile: `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`. Davor `targets` und `planning` je `1277 Datei(en) geprüft, 0 Befund(e)`, Semgrep `Ran 55 rules on 75 files: 0 findings.`; `git status` danach leer |
| `make doc-check` nach dem Schreiben dieses Reports | Exit 0, `d-check: 1284 Datei(en) geprüft, 0 Befund(e)` |

## DoD-Punkte

### 1 — `DC-FA-VCS-001` im Lastenheft geändert: bestätigt

- **Version und Historie nach [`MR-032`](../../harness/conventions.md#mr-032):** Kopf `0.113.2` → `0.114.3`.
  Vier Historie-Zeilen (0.114.0 bis 0.114.3) stehen oben in der Tabelle, die Verweis-Spalte trägt `—`, wie
  MR-032 es verlangt. Die gezählten Kriterien gehen auf: Vorher 5, heute 30; 0.114.0 nennt +15,
  0.114.1 +4, 0.114.2 +4, 0.114.3 +2, zusammen 25. Die Zeile 0.114.0 sagt noch „mit ihm erst gleichgesetzt,
  dann normiert“; das ist der Stand des ersten Inhalts-Commits und wird von 0.114.2 abgelöst (Plan, R2 L-2).
  Das ist Historie, keine Drift.
- **Absatz und Grenzen nennen das Folgen einer Umbenennung:** Im Absatz „Pfad-Nachzug (opt-in)“ ist der
  Satz zum anderen Namen und zu „dieselbe Datei“ auf „ohne `vcs.follow-renames`“ eingeschränkt. Neu sind
  der Absatz „Umbenennung folgen (opt-in)“ (Erkennung, Kette, Mehrdeutigkeit, Gleichsetzung, Lesemenge)
  und der Block „Grenzen von `vcs.follow-renames` (benannt)“ mit zehn Punkten.
- **Kriterium „Boundary (Pfad-Nachzug)“ eingeschränkt:** Titel `Boundary (Pfad-Nachzug, ohne
  vcs.follow-renames)`, Given `vcs.ignore-link-targets: true ohne vcs.follow-renames`. Der Schlusssatz nennt
  jetzt den Schlüssel `vcs.ignore-link-targets` statt „den Schlüssel“. Das Kriterium wurde umgeschrieben;
  ein zweites, widersprechendes steht nicht daneben.
- **Out-of-Scope geschärft:** zwei Punkte statt einem. „Umbenennungs-Erkennung für die immutable Datei
  selbst … bleibt Drift, auch mit `vcs.follow-renames`“ und „Umbenennungs-Erkennung über
  Inhalts-**Ähnlichkeit** — für Link-Ziele zählt nur Inhalts-**Gleichheit**, nie eine
  Ähnlichkeits-Schwelle“. Beides verlangt die DoD.
- **Je Entscheid mit Produktwirkung ein Kriterium.** Die Zuordnung habe ich selbst gemacht, nach Lesen der
  30 Kriterien:

| Entscheid | Kriterium im Lastenheft |
|---|---|
| 1 Inhalts-Gleichheit, keine Schwelle | Negative (gelöscht und neu angelegt) · Negative (umbenennen und ändern), erster Teil |
| 2 jeder Commit, Kette, `--staged` | Happy Path (Umbenennung folgen) „im selben oder in einem anderen Commit“ · Negative (umbenennen und ändern), zweiter Teil · Boundary (Kette) · Boundary (`--staged`) |
| 3 Mehrdeutigkeit | Negative (Mehrdeutigkeit) |
| 4 eigener Schalter, aus ⇒ byte-identisch | Boundary (ohne `vcs.follow-renames` byte-identisch) · Boundary (Pfad-Nachzug, ohne `vcs.follow-renames`) |
| 5 nur das aufgelöste Datei-Ziel | Negative (Umbiegen) · Negative (Anker, Linktext) · Boundary (Art des Eintrags, Verzeichnis) |
| 6 Grenzen, fail-closed, kein neuer Code, immutable Datei | fail-closed (Zwischen-Commit) · Negative (immutable Datei umbenannt); „eine Umbenennung irgendwo im Repo zählt“ ohne eigenes Kriterium (V-3) |
| 9 seiten-gleich, Teilnahme | Boundary (seiten-gleich) · Negative (Name freigemacht und neu belegt) · Negative (zwei Quellen auf ein Ende) · Boundary (Rotation) · Negative (Ende gelöscht und neu angelegt) · Tausch in Negative (Mehrdeutigkeit) |
| 10 Zusammenspiel der Schalter | Boundary (beide Schlüssel) · Boundary (Negativ-Kriterien mit beiden Schlüsseln) · Negative (beide Schlüssel, Umbiegen auf den alten Namen) · Negative (beide Schlüssel, Umbenennung ohne Namenswechsel) |
| 11 Kette entlang der Abstammung, Merge | Boundary (Kette), alle drei Sätze · Negative (Merge-Seite löscht, legt später neu an) · Negative (Zwischenpfad neu belegt) · Boundary (Löschen und Neuanlegen im selben Commit); Reihenfolge-Unabhängigkeit ohne brechbares Kriterium (V-2) |
| 12 Basis kein Vorfahr | Boundary (Basis kein Vorfahr) |
| 13 Port-Vertrag (Art, Ausführbar-Bit, fail-closed) | Boundary (Art des Eintrags, Verzeichnis) · fail-closed (Zwischen-Commit); der Verzicht auf die Bibliotheks-Erkennung wirkt über Negative (Mehrdeutigkeit) |
| 7, 8, 14 | keine Produktwirkung (eigenes Profil, ADR-Status, Commit-Reihenfolge) — kein Kriterium nötig |

  Die Pflicht-Liste der DoD ist vollständig. Die „vier aus dem Issue“ lese ich aus §1 von
  [`welle-96`](../plan/planning/done/welle-96-vcs-folgt-umbenennung.md): Nachzug auf den neuen Namen
  (Happy Path), Löschen und neu Anlegen (gelöscht und neu angelegt), Umbiegen auf eine bestehende Datei
  (Umbiegen), Anker oder Linktext (Anker, Linktext). Dazu kommen Kette, Mehrdeutigkeit, „Basis kein
  Vorfahr“ und „ohne Schalter byte-identisch“. Das Issue selbst habe ich nicht gelesen; das Repo verlinkt
  es nur nach außen.

### 2 — Spezifikation führt Verfahren und Schema, Grenz-Aussagen am Gegenstand: bestätigt, mit V-1

Schritt 2 hat einen neuen Absatz „Umbenennungen der Range“ mit Commit-Menge, Basis kein Vorfahr, Erkennung
je Schritt, Kette sowie „Gelesen und fail-closed“. Schritt 4 führt die drei normierten Formen, die
Teilnahme, die Liste der Ziele, die nicht auflösen, und die geschärfte Grenze. In §2 stehen die neue Zeile
`vcs.follow-renames` und die eingeschränkte Zeile `vcs.ignore-link-targets`. Die drei Stellen aus der
Plan-Änderung vor dem Inhalt sind geliefert.

Stichproben gegen den Gegenstand, selbst gelesen (Code unverändert seit `c0959fff`):

| # | Aussage der Spezifikation | Gegenstand | Ergebnis |
|---|---|---|---|
| 1 | drei Arten zählen als Datei, Verzeichnis und Gitlink nicht | `internal/adapter/driven/git/git.go` Z. 165 `case filemode.Dir` (rekursiv), Z. 173 `case filemode.Regular, filemode.Executable, filemode.Symlink`; ein Gitlink fällt durch kein `case` | hält |
| 2 | im Index nimmt die Enumeration jeden Eintrag | `git.go` Z. 95–98: Schleife über `idx.Entries` ohne Filter nach Modus | hält für den Ist-Stand, steht aber gegen den Soll-Satz desselben Schritts (V-1) |
| 3 | Commit-Menge: von HEAD erreichbar, BASE samt Vorfahren ausgeschlossen, dieselbe wie bei `commits`; die Aufzählung nach Kennung ist nicht topologisch | `git.go` `CommitMessages`: `excluded := a.ancestors(base)`, Walk ab HEAD, `len(c.ParentHashes) <= 1`, `sort.Slice` nach `ShortSHA` (Z. 199–226) | hält |
| 4 | Basis kein Vorfahr ⇒ Urteil wie ohne Schalter, kein Exit 2 | `internal/hexagon/core/rules/vcs.go` ruft nur `AllPaths` und `FileAt`; `AllPaths` (Z. 115–124) prüft nur den Leerfall, keine Abstammung; `isAncestor` ruft nur `rangeTrees` aus `LineChanges` (Z. 402), das allein `metrics` nutzt | hält — ohne Schalter wirft `vcs` bei fremder Basis heute kein Exit 2, „wie ohne den Schlüssel“ ist damit wohldefiniert |
| 5 | sechs Ziel-Formen lösen nicht auf: reiner Anker, `/`, `://`, `?`, außerhalb des Repos, in keinem Baum | `internal/hexagon/core/rules/vcs_link_targets.go` Z. 84 (`t == ""`, `/`, `://`, `?`) und Z. 88 (`..`, `!tree[p]`) | hält, sechs Bedingungen |
| 6 | eine Funktion des aufgelösten Pfads, gegen die Vereinigung beider Bäume, auf beiden Seiten gleich | `vcs.go` Z. 53–56 (`pathTree` über `baseAll` und `headAll`), Z. 130–132 (ein `resolve` für beide `vcsCore`-Aufrufe) | hält |
| 7 | ohne Schlüssel wird kein Commit zwischen BASE und HEAD gelesen | `AllPaths` liest `ancestors(base)` (BASE und davor) und zwei Bäume, `FileAt` zwei Blobs | hält |
| 8 | ADR Alternative C und I: Vorbelegung der Bibliothek 60, exakte Erkennung wählt nach Namensähnlichkeit und verlangt gleichen Modus | Modul-Cache im Image `d-check:test`, go-git v5.19.3: `difftree.go` `RenameScore: 60`; `rename.go` `bestNameMatch(...)` und `sameMode(...)` | hält |

Die Plan-Zählungen „107 `Accepted`-ADRs, 59 mit Link-Ziel `spec/spezifikation.md`“ habe ich mit den
Kommandos aus Plan §3 nachgefahren: 107 und 59.

### 3 — ADR-0115: bestätigt

- `**Status:** Proposed` (Z. 3); Entscheidung 8 bindet `Accepted` an `slice-292`.
- Im [ADR-Index](../plan/adr/README.md) Z. 125 als `Proposed`, Datum 2026-11-23, mit den Bezügen.
- `## Re-Evaluierungs-Trigger` mit fünf beobachtbaren Auslösern.
- `## Verglichene Alternativen` mit A–M und der gewählten Option K. Darunter **C** (Ähnlichkeits-Schwelle),
  **D** (nur die Endpunkte BASE und HEAD) und **B** (Bedeutung von `ignore-link-targets` erweitern).
- Die ADR ist `Proposed`; ihre Änderungen über die Range (sieben Architect-Commits) berührt
  [`AGENTS.md`](../../AGENTS.md) §3.5 nicht.

### 4 — `make gates` grün: bestätigt

Selbst gefahren, Ausgabe oben. Exit 0 auf dem Stand `1c571ac9`.

### 5 — Review-Reports liegen vor, R4 ohne offenen MEDIUM oder HIGH: bestätigt

Vier Reports: R1 und R2
[`2026-11-23-slice-291-vcs-umbenennung-r1.md`](2026-11-23-slice-291-vcs-umbenennung-r1.md) und
[`…-r2.md`](2026-11-23-slice-291-vcs-umbenennung-r2.md), R3
[`…-r3.md`](2026-11-23-slice-291-vcs-umbenennung-r3.md), R4
[`2026-11-24-slice-291-vcs-umbenennung-r4.md`](2026-11-24-slice-291-vcs-umbenennung-r4.md).
R4 meldet HIGH 0, MEDIUM 1 (M-1), LOW 3 und INFO 1. Nachgesehen habe ich, ob die Nachbesserungen
wirklich im Text stehen:

| R4-Befund | Commit | im Text |
|---|---|---|
| M-1 (Plan erzählt die drei Fälle anders) | `2abd7848` | Plan Entscheid 9 erzählt nicht mehr nach. Er verweist auf die Grenze „Drei Fälle lassen etwas durch …“ in 0.114.3 und auf die ADR, mit der Bestätigung vom 2026-11-24. Geschlossen |
| L-1 (Fitness-Zeile, Testfall der nachsichtigen Hälfte) | `1c571ac9`, `2abd7848` | Die Fitness-Zeile nennt jetzt „die nicht teilnehmende Quelle … (C in Fall 1 und 3, A und B in Fall 2)“ und die Ausnahme A → X. `slice-292` §3 führt den neuen Testfall „nachsichtige Hälfte der drei Fälle mit beiden Schaltern“ |
| L-2 (Probe ohne Bedingung) | `1c571ac9` | „mit anderem Inhalt, in einem Commit gelöscht und in einem späteren neu angelegt ⇒ Befund“ |
| L-3 (Rot-Beleg am Tausch) | `2abd7848` | `slice-292`: den Rot-Beleg der BASE-Bedingung trägt der Fall „A gelöscht, dann C → A“ |
| I-1 (Geschichte-Zeile) | `1c571ac9` | berichtigt: „nannte ‚nur in Richtung Drift‘ für Entscheidung 10 und für Entscheidung 9 keine Richtung“ |

Einen eigenen Folge-Review über `2abd7848` und `1c571ac9` gibt es nicht. Die Nachbesserung habe ich hier
gegen den Befundtext geprüft, nicht als Review.

### 6 — Konsistenz in den drei nachsichtigeren Fällen: bestätigt

Ich habe jeden Fall gegen die Teilnahme-Regel der Spezifikation (Schritt 4) durchgerechnet, mit allen drei
Schlüssel-Kombinationen. Lastenheft-Absatz „Was gleichgesetzt wird“ und Spezifikation Schritt 4 haben
denselben Wortlaut für die Teilnahme, die Rotation eingeschlossen.

| Fall | Teilnahme nach Spezifikation | Lastenheft (Grenze, Kriterium) | ADR §Konsequenzen „Grenze“ |
|---|---|---|---|
| (1) A gelöscht, danach C → A | A: Linie endet mit „gelöscht“, keine Quelle; C → A: Ende lag in BASE ⇒ weder C noch A. `fr`: Nachzug C → A roh ⇒ Drift. Beide: `c.md`/`a.md` ⇒ Drift; Umbiegen C → andere `…/c.md` ⇒ gleich ⇒ kein Befund | Grenze: „weder C noch A … kein Befund“; Kriterium „Name freigemacht“: Drift | „C nimmt nicht teil“; „mit beiden … so grün wie ohne `follow-renames`“ — gleichsinnig |
| (2) A → E und B → E auf zwei Merge-Seiten | E ist Ende zweier Quellen ⇒ keine nimmt teil. `fr`: Umbiegen A → B roh ⇒ Drift. Beide: `a.md` ≠ `b.md` ⇒ Drift; Umbiegen A → andere `…/a.md` ⇒ kein Befund | Grenze: „keine der Quellen … kein Befund“; Kriterium „zwei Quellen“: Drift mit und ohne `ignore-link-targets` | „keine nimmt teil (vorher beide)“ — gleich |
| (3) Rotation A → X, danach C → A | A → X nimmt teil (X nicht in BASE, eine Quelle); C → A nicht (A lag in BASE); A nimmt als Quelle teil ⇒ A ≡ X. Nachzug A → X ⇒ `\0\0X` beidseitig ⇒ kein Befund, auch `fr` allein; Nachzug C → A ⇒ Drift | Grenze: „A nimmt als Quelle von A → X teil … auch mit `vcs.follow-renames` allein“; Kriterium „Rotation“: beide Hälften | „A nimmt als Quelle von A → X teil, ein Nachzug von A auf X folgt“ — gleich |

Der Rahmensatz des Lastenhefts („jedes dieser Ergebnisse liefert `vcs.ignore-link-targets` allein ebenso,
oder eine teilnehmende Umbenennung mit demselben Ende trägt es“) hält in allen drei Fällen. Die Zusage
„mit beiden Schlüsseln entfällt ein Befund … nur dort, wo beide Ziele teilnehmen und dasselbe Ende haben“
hält auch bei der Rotation: A und X nehmen teil und haben das Ende X. Der Plan nennt die ADR
„gleichlautend“; in der Sache ist sie gleich, im Wortlaut nicht (V-4). Die bestätigte Fassung 0.114.3 und
die ADR widersprechen sich nicht.

### 7 — Referenz-Richtung ([`AGENTS.md`](../../AGENTS.md) §3.4): bestätigt

In den `+`-Zeilen von `spec/` passt nur eine Zeile auf das Muster
`ADR-?[0-9]|adr/|slice-[0-9]|welle-[0-9]|[0-9a-f]{7,40}|Issue`: die Historie-Zeile 0.114.0 mit
„GitHub-Issue 18“. Eine Issue-Nummer gehört zu keiner der fünf Kategorien, und sie gibt den Anlass an.
`spec/architecture.md` ist unberührt. `make doc-check` (Modul `matrix`) ist im Gate-Lauf grün.

### 8 — Risiken aus §6: Ausgänge vorgeschlagen (Träger der Eintragung: Closure)

| Risiko | Vorschlag | Beleg |
|---|---|---|
| 1 — „Boundary (Pfad-Nachzug)“ umschreiben statt ein widersprechendes Kriterium danebenstellen | **eingetreten** | Das Kriterium selbst ist umgeschrieben, nicht verdoppelt (R1 Negativbefund, hier nachgelesen). Die Klasse trat aber ein: R1 M-1 meldete zwei Rang-1-Stellen, die für dieselbe Eingabe Befund und Nicht-Befund zusagen, und schrieb ausdrücklich „Risiko 1 des Plans: eingetreten“; Fix `4317dc0a` (0.114.1). Danach stand der Plan noch zweimal gegen Lastenheft und ADR (R3 M-1, R4 M-1); Fixes `90fd2340` (nach R3; `7a57fa1e` war die Berichtigung, die R3 M-1 bemängelte) und `2abd7848`. Laut Ausgangsregel fehlt noch eine Evidenz-Datei in [`BEO-ALL/plan-entscheid-steht-gegen-lastenheft`](../plan/planning/observations/BEO-ALL/plan-entscheid-steht-gegen-lastenheft/observation.md). Dort liegen heute `slice-289.md` und `slice-290.md`; `slice-291.md` wäre die dritte, und die Schwelle wäre erreicht |
| 2 — `--staged` sieht einen früheren Move-Commit nicht | **weiter offen** | Die Anforderung nennt die Grenze (Lastenheft, Grenze `--staged`; Kriterium „Boundary (`--staged`)“, zweiter Teil: Drift). Der Architect hat keinen Weg bestätigt, der den früheren Commit mitliest: ADR Alternative J ist verworfen. Damit gilt die erste Hälfte der Ausgangsregel |
| 3 — Commit-Menge bei Force-Push oder fremder Basis | **entfallen** | Die Spezifikation legt in Schritt 2 die Commit-Menge und „Basis kein Vorfahr ⇒ leere Abbildung, kein Exit 2“ fest; das Lastenheft hat Grenze und Kriterium „Boundary (Basis kein Vorfahr)“. [`slice-292`](../plan/planning/done/slice-292-vcs-umbenennung-modulkern.md) §3 Z. 144 führt den Testfall. Der Test selbst entsteht erst mit `slice-292`; dort gehört er in die DoD |

## Bewusstes Brechen

Kein DoD-Punkt dieses Slice stützt sich auf einen Test; der Slice liefert nur Text. Die Rot-Belege
der Kriterien liegen bei `slice-292` (Testfall-Liste §3, ADR §Fitness Function). Hier habe ich die
Grenz-Aussagen am Gegenstand gegengelesen (DoD 2) und die drei Fälle von Hand durchgerechnet (DoD 6).

## Befunde des Verifiers

### V-1 — LOW: Schritt 2 beschreibt für den Index den Ist-Code, gegen den Soll-Satz desselben Schritts

- **Stelle:** §DC-FA-VCS-001.a Schritt 2, Spiegelstrich „Eine Umbenennung in einem Schritt“: „… und zwar
  auch im Index, **wo die Enumeration jeden Eintrag nimmt**; ein Verzeichnis und ein Gitlink nehmen nicht
  teil.“
- **Befund:** Der Satz gibt den heutigen Code wieder: `git.go` Z. 95–98 nimmt jeden Index-Eintrag, auch
  einen Gitlink. Weiter oben im selben Schritt steht dagegen als Soll: „enumeriert **jeden** Datei-Pfad
  (regulär/ausführbar/symlink) an BASE und an HEAD (bzw. am staged Index)“. `slice-292` §3 (Zeile zu
  `git.go`) zieht den Code an dieses Soll nach, mit „ein Filter für Enumeration und
  Umbenennungs-Abbildung“, und setzt den Testfall „Gitlink im Index … kein Kandidat“ dazu. Danach ist der
  Nebensatz falsch. Schon heute sagt Rang 2 über die Index-Enumeration zweierlei. ADR-0115 Entscheidung 13
  nennt denselben Ist-Stand als Kontext („im Index nimmt `AllPaths` jeden Eintrag“); in der ADR ist das
  unschädlich.
- **Warum keine DoD-Verletzung:** DoD 2 verlangt, dass eine Grenz-Aussage am Gegenstand gemessen ist. Das
  ist sie, heute stimmt sie. Den Mechanismus berührt sie nicht, denn die Abbildung schließt den Gitlink
  in beiden Lesarten aus.
- **Vorschlag:** den Nebensatz streichen („Teil haben die drei Arten …, auch im Index; ein Verzeichnis und
  ein Gitlink nehmen nicht teil“). Alternativ zieht `slice-292` die Spezifikation im selben Zug nach, in dem
  es `AllPaths` filtert. Träger: Implementer von `slice-292`; oder ein Nachzug vor der Closure dieses
  Slice, wenn der Planner das will.

### V-2 — LOW: die Reihenfolge-Unabhängigkeit (Entscheid 11) hat kein Kriterium, das brechen kann

- **Stelle:** Lastenheft, Kette: „Das Ergebnis hängt nicht von der Reihenfolge ab, in der Commits
  aufgezählt werden“. Spezifikation Schritt 2: dieselbe Aussage samt „eine Aufzählung nach
  Commit-Kennung … ist nicht topologisch“.
- **Befund:** Das einzige Kriterium dazu ist „Derselbe Lauf zehnmal ergibt byte-gleiche Ausgabe“
  (Boundary (Kette)). Gegen eine reihenfolge-abhängige Implementierung kann es nicht rot werden: Die
  Aufzählung ist fest nach Kurz-SHA sortiert (`git.go` Z. 226), zehn Läufe sehen also zehnmal dieselbe
  Reihenfolge. Auch die Testfall-Liste von `slice-292` hat keinen Fall dazu, nur „zehn Läufe …
  byte-gleich“ (Z. 186).
- **Vorschlag:** In `slice-292` einen Testfall aufnehmen, der dieselbe Topologie mit vertauschter
  SHA-Ordnung fährt, etwa mit anderen Commit-Botschaften, sodass ein Nachfahre kleiner sortiert als sein
  Vorfahr; erwartet wird dasselbe Ergebnis. Träger: Planner bzw. Implementer von `slice-292`.

### V-3 — INFO: „eine Umbenennung irgendwo im Repo zählt“ (Entscheid 6) ohne ausdrückliches Kriterium

Kein Kriterium legt die umbenannte Datei ausdrücklich außerhalb von `vcs.paths` und der Scan-Menge. Der
Happy Path spricht von „eine Datei A“ ohne Einschränkung und deckt den Fall damit sinngemäß; der Anlass
selbst (eine Spec-Datei außerhalb von `vcs.paths`) ist die Probe im Closure-Trigger von `welle-96`. Kein
Handlungsbedarf in diesem Slice.

### V-4 — INFO: „gleichlautend“ heißt hier gleichsinnig

Plan Entscheid 9 nennt ADR §Konsequenzen „gleichlautend“ mit der Grenze im Lastenheft. Die Sätze sind
verschieden: Die ADR sagt in Fall 1 nur „C nimmt nicht teil“, das Lastenheft „weder C noch A“, und die
ADR rahmt mit „urteilt über die nicht teilnehmende Quelle wie ohne `follow-renames`“. In der Sache
stimmen sie überein (Tabelle zu DoD 6). Das ist kein Widerspruch, nur ein Wort im Plan, das mehr sagt,
als der Text trägt.

## Plan-vs-Code-Diff (Plan gegen gelieferten Text)

| Plan | Geliefert | Abweichung |
|---|---|---|
| §3 Tabelle: Lastenheft, Spezifikation, ADR (Angleichung durch den Architect), ADR-Index, Plan `slice-292` | alle fünf; die ADR sieben Mal angeglichen, der Index unverändert (stand schon vor der Range) | keine |
| Plan-Änderung vor dem Inhalt: drei zusätzliche Stellen der Spezifikation | Schritt 4 nennt sechs Ziel-Formen; §2-Zeile `vcs.ignore-link-targets` „(ohne `vcs.follow-renames`)“; Schritt 2 „wirkt nur auf Link-Ziele in Schritt 4, nicht auf diese Mengen“ | keine |
| Plan-Änderungen nach R1 (+3 Kriterien), R2 (+4), R3 (+2) | in der Historie gezählt +4, +4, +2. Nach R1 sind es vier, nicht drei: das vierte ist „Negative (Merge-Seite löscht, legt später neu an)“ aus `18157ed7`, eine ADR-Folge von Entscheidung 11, die R2 L-4 benannt hat | eine Stelle mehr als die Plan-Notiz, in der Historie korrekt gezählt |
| §3 Tabelle nennt den Welle-Plan nicht | `welle-96-vcs-folgt-umbenennung.md` von vier Planner-Commits geändert (`6ccf5ac5`, `7d4643c4`, `7a57fa1e`, `90fd2340`) | Planner-Träger außerhalb der Tabelle, inhaltlich nur Nachzug der Bestätigungen; keine Abgrenzungsverletzung |
| §1 Abgrenzung: kein Code, kein Profil, keine Release-Prep, Grenze „gleichnamige andere Datei“ von `ignore-link-targets` allein bleibt | Diff ohne Code, Profil und Release-Dateien; mit `ignore-link-targets` allein ist das Verhalten unverändert, mit beiden Schaltern geschärft (Entscheid 10, in `welle-96` §6 als Ausnahme geführt) | keine |
| Entscheid 9, Kern im Plan | verweist seit `2abd7848` auf Lastenheft und ADR, statt nachzuerzählen | „gleichlautend“ siehe V-4 |
| DoD 2: Spezifikation führt das Verfahren | geführt; Index-Satz siehe V-1 | V-1 |

## Gesamturteil

**Closure-reif für die Liefer-Punkte.** DoD 1–5 sind mit eigenem Beleg bestätigt, die Konsistenz in den
drei Fällen und die Referenz-Richtung ebenso. V-1 und V-2 sind LOW, gehören zu `slice-292` oder einem
kleinen Spec-Nachzug und blockieren nicht. Offen sind nur die Closure-Pflichten selbst: Closure-Notiz
mit Lerneintrag, Fortschreibung des Beobachtungs-Registers (bei Risiko 1 mit dem dritten Beleg zu
`BEO-ALL/plan-entscheid-steht-gegen-lastenheft`; R4 nennt dazu zwei Klassen über der Schwelle), die
Ausgänge in §6 nach dem Vorschlag oben, die Paarungen bei der Welle-Closure und der `git mv` nach
`done/`. `make verify-closure-notes` greift erst, wenn der Slice nach `done/` wandert, und lief hier
deshalb nicht.
