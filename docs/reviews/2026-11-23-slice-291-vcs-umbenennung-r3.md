# Review-Report: slice-291 — `vcs` folgt einer Umbenennung, Fix-Commits nach R2 (R3) — 2026-11-23

**Review-Art:** Design (Diff gegen Slice-Plan, ADR und Hard Rules, Maintainability;
`v6.18.1` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten). Nicht gegen die DoD — die
prüft danach der Verifier.

**Gegenstand:** `slice-291`, Commit-Range `7dc0b423..7a57fa1e` (sieben Commits, nur lokal auf `main`):
`e84632b1` und `d03cda9c` ([ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md), Architect,
`Proposed`), `7d4643c4` und `7a57fa1e` (Planner: Pläne `slice-291`, `slice-292`, Welle `welle-96`),
`8bcbbc67` (Plan-Änderung des Implementers), `e3e5d4eb` (Lastenheft
[`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in),
Versionen 0.114.1 und 0.114.2), `5b90b304` (Spezifikation
§[`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
Schritt 2 und 4, Schema-Zeile). Gelesen wurden Absatz, Grenzen, alle 28 Kriterien
(Lastenheft Z. 2373–2400), Spezifikation Schritt 2 und 4 und die ADR im Stand `7a57fa1e`
nebeneinander, dazu die ADR im bestätigten Stand `c89e2dee` und das Lastenheft im Stand `7dc0b423`.

**Skill:** `.harness/skills/reviewer.md` @ 1.20.0 (`d1c37c14`).
**Modell:** claude-opus-5-5 · **Datum:** 2026-11-23.

**Eingangs-Kontext:**

- Reports R1 ([`2026-11-23-slice-291-vcs-umbenennung-r1.md`](2026-11-23-slice-291-vcs-umbenennung-r1.md))
  und R2 ([`2026-11-23-slice-291-vcs-umbenennung-r2.md`](2026-11-23-slice-291-vcs-umbenennung-r2.md))
- Plan `slice-291` §1, §3 samt allen Plan-Änderungen; Plan `slice-292` (Testfall-Liste);
  Welle `welle-96` §5 und §6
- [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) (Entscheidung 1–14, Alternativen K–M,
  Konsequenzen, Fitness Function, Geschichte)
- [`MR-032`](../../harness/conventions.md#mr-032) und der Kopf des Lastenhefts;
  [`AGENTS.md`](../../AGENTS.md) §3.3, §3.4, §3.5, §5 (Regeln 13, 15, 16), §6 Schritt 4;
  `v6.18.1` · `regelwerk/modul-08-agentenrollen.md` §Konflikt-Pfad als Rollen-Sequenz;
  `v6.18.1` · `regelwerk/modul-11-verification.md` §Bewusstes Brechen für DoD-Testbehauptungen
- Gelesen am Gegenstand: `internal/adapter/driven/git/git.go` (`AllPaths`, `ancestors`, Index-Zweig),
  `internal/hexagon/core/rules/vcs_link_targets.go` (`normalizedTargetMark`, `markFreeTree`,
  `linkTargetResolver`) und `vcs.go`
- `.d-check.yml` Block `vcs` (setzt `ignore-link-targets: true`; mit Entscheidung 7 laufen im eigenen
  Profil beide Schlüssel)

**Lesart von Schwerpunkt 2(a).** Wörtlich („öffnet kein Grün“) kann die Zusage nicht gemeint sein: ein
Nachzug auf den neuen Namen ist ohne den Schlüssel Drift und mit ihm keine — das ist der Zweck. Geprüft
wurde deshalb die Zusage, die die Straten selbst geben: ohne `vcs.ignore-link-targets` kein Befund, den
es ohne `vcs.follow-renames` nicht gibt (Ausnahme Exit 2), mit beiden Schlüsseln entfällt ein Befund von
`vcs.ignore-link-targets` allein nur, wo beide Ziele teilnehmen und dasselbe Ende haben; und jedes Grün
mit dem Schlüssel geht auf eine teilnehmende Umbenennung zurück.

---

## Stand der R2-Befunde

| R2 | Stand | Beleg |
|---|---|---|
| M-1 (Doppelrolle: gelöscht, danach Ende einer anderen Umbenennung) | **geschlossen** | Alle drei Straten schließen den Pfad aus: Lastenheft Z. 2294–2298, Spezifikation Z. 2152–2161, ADR Entscheidung 9 (Z. 112–116) und §Konsequenzen (Z. 196–199); Kriterium Z. 2393. Die Zwischenpfad-Doppelrolle ist mit „weder Quelle noch Ende“ und der BASE-Bedingung an das Ende aufgelöst (Z. 2282–2283, Spezifikation Z. 2088–2089). Eine dritte Doppelrolle (Rotation) steht als I-1 unten |
| M-2 (Löschen und Neuanlegen ohne Schritt-Bedingung) | **geschlossen** | Lastenheft Z. 2278–2280, 2323–2327, 2340–2344, 2351–2355, Kriterien Z. 2385, 2386, 2394; Spezifikation Z. 2085–2086; ADR Entscheidung 11 (Z. 136–137), §Konsequenzen Z. 212–220, Fitness Function Z. 233–234. Die einzige Stelle ohne Bedingung liegt in `slice-292` (L-2) |
| M-3 (beide Schlüssel, Umbiegen auf eine umbenannte Quelle) | **geschlossen** | Eigene Marke für teilnehmende Ziele: Lastenheft Z. 2305–2308, Grenze Z. 2334–2339, Kriterien Z. 2390–2392; Spezifikation Z. 2144–2150, 2176–2182; ADR Entscheidung 10, §Konsequenzen Z. 202–207, Alternative L. „nie grüner“ und „strenger als … und nur“ stehen in keinem Stratum mehr |
| L-1 („nur im Fall fail-closed (Zwischen-Commit)“) | **geschlossen** | Lastenheft Z. 2309–2316, Kriterien Z. 2388 und 2399 nennen Commit der Range, Baum eines Schritts oder seines Eltern-Commits (auch vor BASE) und geänderten Teilbaum; Spezifikation Z. 2101–2108; ADR Entscheidung 13 (Z. 146–150). Gegen den Code gehalten: `AllPaths` liest ohne Schlüssel von den Vorfahren der Basis nur Commit-Objekte, Bäume nur von BASE und HEAD |
| L-2 (Historie-Zeile 0.114.0 fortgeschrieben) | **geschlossen** | Zeile 0.114.0 byte-gleich zu `f0622065` (`diff` leer), eigene Zeilen 0.114.1 und 0.114.2 (Z. 4768–4769), Kopf 0.114.2 (Z. 3). Rest I-2 |
| L-3 (Botschaften „nur in Richtung Drift“) | **in der Sache überholt, Klasse wiederholt** | Die Botschaften `3a2d7e3f`/`6295c2d4` bleiben, wie sie sind; Entscheidung 9 ist in neuer Fassung bestätigt. Dieselbe Behauptung kehrt in `e84632b1`, `e3e5d4eb` und in den Plänen wieder: M-1, M-2, L-1 |
| L-4 (Plan-Änderung zählt den Inhalt nicht) | **geschlossen** | `8bcbbc67` nennt vier neue und fünf angeglichene Kriterien und die Versionsform; `e3e5d4eb` trägt genau diese (Wortdiff der Kriterienzeilen). Die Frage des Implementers zur Richtung liegt jetzt als Artefakt vor (Botschaft `e3e5d4eb`, §Richtung) |
| L-5 (Satz „nur noch in BASE“ ohne seinen Schlüssel) | **geschlossen** | Spezifikation Z. 2206–2210 bindet den Satz an zwei nicht teilnehmende Ziele und an `vcs.ignore-link-targets` |
| I-1 (Linie ohne den Pfad ohne Klasse) | **geschlossen** | „auf der der Pfad nicht liegt“ ist unberührt: Lastenheft Z. 2276–2278, Spezifikation Z. 2083–2085, ADR Z. 132–133 |
| I-2 („gelöscht“ als Ende in der ADR) | **geschlossen** | ADR Entscheidung 11 sagt „ohne Ende“ und nennt die Bedingung wie die Spezifikation Z. 2093–2097 („mindestens eine Linie ein Ende …, keine Linie ohne Ende“) |
| I-3 (Testfall „Tausch in einem Commit“) | **offen, verschärft** | `7d4643c4` schreibt den Testfall um und behält „in einem Commit“; er nennt jetzt eine Abbildung, die es nach Schritt 2 nicht gibt: M-3 |

---

## Findings

### M-1 — Die Berichtigung der Richtung nennt das falsche Ziel und lässt zwei Fälle aus

- **Kategorie:** MEDIUM
- **Quelle:** [`AGENTS.md`](../../AGENTS.md) §5 Regel 15; Prüffrage 8 und 18;
  [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 9
- **Pfad:** `docs/plan/planning/done/slice-291-vcs-umbenennung-anforderung-und-adr.md` ·
  „mit beiden Schaltern fällt ein Ziel, dessen Pfad in BASE lag“ (Z. 205–209);
  `docs/plan/planning/welle-96-vcs-folgt-umbenennung.md` · „ist Entscheidung 9 an einer“ (Z. 113–117);
  Botschaft `7a57fa1e` · „faellt mit beiden Schaltern ein Ziel, dessen Pfad in BASE lag und Ende einer
  Umbenennung wird“
- **Befund:** Gegen die bestätigte Fassung (`c89e2dee`: „Teil hat ein Pfad P, wenn er Quelle … oder wenn
  er ein solches Ende ist. Ein Pfad, der zugleich Quelle der einen und Ende einer anderen Umbenennung ist,
  nimmt nicht teil“) ändert sich für das Ende, das in BASE lag, nichts: es wurde vorher mit beiden
  Schlüsseln auf seinen eigenen Dateinamen normiert (Ende seiner selbst) und wird es jetzt wieder. Was
  zurückfällt, ist die **Quelle**: A gelöscht, C → A — C wurde vorher zu Marke und `a`, jetzt zu Marke
  und `c`; ein Umbiegen von C auf eine andere Datei `…/c.md` war Drift und ist kein Befund mehr. Zwei
  weitere Fälle fehlen: (a) zwei Quellen auf ein Ende (Merge, A → E und B → E) — vorher nahmen beide teil,
  jetzt fallen beide Quellen auf ihren eigenen Dateinamen zurück, mit beiden Schlüsseln wird ein
  Umbiegen von A auf eine andere `…/a.md` grün; (b) Rotation (A → X, danach C → A): vorher nahm A als
  „zugleich Quelle und Ende“ nicht teil, jetzt nimmt A als Quelle von A → X teil — der Nachzug von A auf
  X ist auch mit `vcs.follow-renames` **allein** nicht mehr Drift. Der Plan sagt dagegen „nicht nur
  strenger: mit beiden Schaltern …“ und lässt damit den Lauf mit dem Schlüssel allein als nur strenger
  stehen; die Welle zählt „an einer Stelle“.
- **Failure-Szenario:** Der Plan-Absatz ist die Grundlage, auf der entschieden wird, ob die Bestätigung
  vom 2026-11-23 trägt (er schließt mit „Die Vorlage … nannte … ‚nur in Richtung Drift‘“). Wer ihn liest,
  sieht eine Abweichung, die nur mit beiden Schlüsseln und nur als Rückfall auf
  `vcs.ignore-link-targets` allein auftritt, und hält eine erneute Vorlage für entbehrlich. Mit
  Entscheidung 7 laufen im eigenen Profil beide Schlüssel (`.d-check.yml`, `ignore-link-targets: true`);
  die Abweichung, die auch mit dem Schlüssel allein grüner ist (Rotation), hat niemand vorgelegt
  bekommen. In `slice-292` sucht ein Test, der die Richtung belegen soll, am Ende statt an der Quelle
  und findet nichts.
- **Warum nicht HIGH:** Jedes neue Grün geht auf eine teilnehmende Umbenennung mit demselben Ende oder
  auf das Ergebnis von `vcs.ignore-link-targets` allein zurück — gegen den Lauf ohne
  `vcs.follow-renames` öffnet keiner der drei Fälle ein Grün außerhalb der Zusage. Kein Code.
- **Verifizierbar:** ja, ab `slice-292` — Fixtures „A gelöscht, C → A, Umbiegen C → …/c.md“,
  „Merge A → E | B → E, Umbiegen A → …/a.md“ und „A → X, C → A, Nachzug A → X“, je gegen die Fassung
  `c89e2dee` gelesen.
- **Klasse:** richtungs-aussage-am-falschen-ziel (neu); trägt
  `BEO-ALL/commit-message-overclaims-work`, zweite Richtung
- **Träger:** Planner (Plan §3, Welle §5)

### M-2 — Die ADR verzeichnet die Bestätigung von Entscheidung 9 ohne den Vorbehalt, unter dem sie gegeben wurde

- **Kategorie:** MEDIUM
- **Quelle:** [`AGENTS.md`](../../AGENTS.md) §5 Regel 16; Prüffrage 9
- **Pfad:** `docs/plan/adr/0115-vcs-folgt-umbenennung.md` · „sind nach Review R2 in der Sache geändert
  und vom Auftraggeber am 2026-11-23 in dieser Fassung bestätigt“ (Z. 78–79) und Geschichte ·
  „der Auftraggeber hat Entscheidung 9 und 10 in der Fassung nach Review R2 bestätigt“ (Z. 259)
- **Befund:** Der Plan hält fest, dass die Vorlage an den Auftraggeber die Änderung „nur in Richtung
  Drift“ nannte und dass das gegen die vorher bestätigte Fassung nicht hält (Z. 208–209). Die ADR zitiert
  die Bestätigung für die Fassung als Ganzes, ohne diesen Vorbehalt und ohne einen Beleg, dass die
  nachsichtigere Seite (M-1) vorgelegt oder ihre Vorlage als entbehrlich entschieden wurde. Die
  Geschichte-Zeile Z. 258 trägt die Richtungs-Aussage nur für Entscheidung 10 (dort stimmt sie: mit
  gleicher Teilnahme wirkt die eigene Marke nur in Richtung Drift); für Entscheidung 9 steht keine.
- **Failure-Szenario:** Beim Übergang nach `Accepted` (Entscheidung 8, nach `slice-292`) wird die ADR
  immutabel und sagt „in dieser Fassung bestätigt“. Der Plan, der den Vorbehalt trägt, wandert nach
  `done/` und ins Archiv. Ein späterer Retirement- oder Re-Evaluierungs-Lauf liest eine informierte
  Bestätigung, wo eine auf eine unzutreffende Vorlage gegebene steht — die Abweichung in Richtung Grün
  hat dann keinen Ort mehr, an dem jemand über sie entschieden hätte.
- **Warum nicht HIGH:** Die ADR ist `Proposed` und der Fall ist nachholbar; die nachsichtigere Seite
  landet auf schon akzeptiertem Verhalten (`vcs.ignore-link-targets` allein) oder auf einer echten
  Umbenennung.
- **Verifizierbar:** nein
- **Klasse:** bestaetigung-ueber-geaenderten-wortlaut (zweites Auftreten in diesem Slice nach R1 I-3)
- **Träger:** Architect (ADR), Entscheidung beim Auftraggeber

### M-3 — Der Testfall „Tausch in einem Commit“ nennt eine Abbildung, die Schritt 2 nicht kennt, und fährt die Regel nicht, die er belegen soll

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 20; `v6.18.1` · `regelwerk/modul-11-verification.md` §Bewusstes Brechen für
  DoD-Testbehauptungen; [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 1 und 9
- **Pfad:** `docs/plan/planning/open/slice-292-vcs-umbenennung-modulkern.md` · „Tausch A ↔ B in einem
  Commit (A → B und B → A, gleiche Blobs gekreuzt)“ (Z. 148–149)
- **Befund:** In einem Commit bestehen A und B vor und nach dem Schritt; nach Lastenheft Z. 2263–2266
  und Spezifikation Z. 2073–2077 ist das keine Umbenennung — es gibt weder A → B noch B → A. Der Testfall
  begründet „kein Pfad nimmt teil“ mit „weil jedes Ende in BASE existiert (Entscheid 9)“; er bliebe grün,
  wenn die BASE-Bedingung fehlte. Der Fall, den die Regel entscheidet — Tausch über einen Zwischennamen
  (ADR §Konsequenzen Z. 196–197) —, steht in keiner Zeile der Liste.
- **Failure-Szenario:** `slice-292` baut das Fixture wie geplant, der Test ist grün und gilt als Beleg
  für Entscheidung 9. Fehlt im Code die Bedingung „Ende nicht in BASE“, nehmen beim Tausch über T
  A → B und B → A beide teil, A ist zugleich Quelle (Ende B) und Ende (gilt als A) — die Normalform ist
  nicht mehr eine Funktion des Pfads, und kein Test der Liste wird rot. Kriterium Z. 2384 („Ein Tausch
  zweier Namen samt Nachzug ist `core-drift-vcs`“) gilt dann als getestet, ohne dass ein Tausch lief.
- **Warum nicht HIGH:** Kein Code; die Lücke liegt im Plan des Folge-Slice und ist vor dessen erstem Edit
  behebbar.
- **Verifizierbar:** ja, ab `slice-292` — Mutation „BASE-Bedingung entfernt“ gegen das geplante Fixture:
  bleibt grün.
- **Klasse:** test-begruendung-faehrt-ihre-regel-nicht (zweites Auftreten nach R2 I-3)
- **Träger:** Planner (`slice-292`)

### M-4 — „Ende keiner zweiten Quelle“ ist neu und hat kein Kriterium, keinen Testfall und keine Grenze

- **Kategorie:** MEDIUM
- **Quelle:** Prüffrage 13 und 18;
  [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  Absatz „Was gleichgesetzt wird“
- **Pfad:** `spec/lastenheft.md` · „existiert und Ende keiner zweiten Quelle ist“ (Z. 2295), Kriterien
  Z. 2373–2400; `docs/plan/planning/open/slice-292-vcs-umbenennung-modulkern.md` · Testfall-Liste
  (Z. 130–173); `docs/plan/adr/0115-vcs-folgt-umbenennung.md` · Fitness Function (Z. 233–234)
- **Befund:** Die Bedingung steht gleichlautend in Lastenheft, Spezifikation (Z. 2153–2154), ADR
  Entscheidung 9 (Z. 113–114) und Schema-Zeile. Erreichbar ist sie nur über einen Merge (A → E auf der
  einen, B → E auf der anderen Seite). Keines der 28 Kriterien, keine Zeile der Testfall-Liste und keine
  Zeile der Fitness Function fährt diesen Fall; die Grenzen-Liste (Z. 2319–2357) nennt ihn nicht, nur
  Plan `slice-291` §3 Entscheid 9.
- **Failure-Szenario:** `slice-292` implementiert ohne die Bedingung — jedes der 28 Kriterien und jede
  Zeile der Testfall-Liste bleibt grün. In einem Repo, in dem zwei Zweige zwei verschiedene Dateien auf
  denselben Namen verschieben, wird ein Link, der in einer `Accepted`-ADR von A auf B umgebogen wird,
  zu E = E, `make adr-check` bleibt grün, obwohl A und B verschiedene Dokumente waren und der Lauf ohne
  den Schlüssel Drift meldet.
- **Warum nicht HIGH:** Kein Code, und der Vertrag selbst ist in allen Straten gleich.
- **Verifizierbar:** ja, ab `slice-292` — Fixture mit Merge zweier Seiten, die A bzw. B nach E
  umbenennen, Umbiegen von A auf B.
- **Klasse:** neue-vertragsklausel-ohne-negativtest
- **Träger:** Implementer (Lastenheft), Planner (`slice-292`)

### L-1 — Zwei Botschaften behaupten weiter „nur in Richtung Drift“; berichtigt ist eine, in der nächsten Botschaft

- **Kategorie:** LOW
- **Quelle:** [`AGENTS.md`](../../AGENTS.md) §5 Regel 15
- **Pfad:** Commit `e84632b1` · „Beide Aenderungen wirken nur in Richtung Drift“; Commit `e3e5d4eb` ·
  „mit follow-renames allein wirkt die Aenderung nur in Richtung Drift“
- **Befund:** `5b90b304` widerruft die Aussage von `e3e5d4eb` ausdrücklich („traegt der Diff nicht“) und
  nennt ein Gegenbeispiel (Löschen und Neuanlegen im selben Commit); ein zweites ist die Rotation aus M-1.
  Für `e84632b1` gibt es keinen Widerruf in einer Botschaft; `7a57fa1e` berichtigt nur den Plan- und den
  Welle-Text. Abgewogen: Ein Amend ist hier ausgeschlossen, der Widerruf in der unmittelbar folgenden
  Inhalts-Botschaft desselben Pushs ist die Form, die das Repo hat, und er steht in `git log` neben der
  falschen Aussage — für `e3e5d4eb` genügt er; `git show e3e5d4eb` allein zeigt die Aussage weiterhin
  ohne ihn. Die dritte Wiederholung derselben Klasse in diesem Slice (R2 L-3 mit zwei Botschaften, hier
  zwei weitere) ist ein Steering-Loop-Signal.
- **Verifizierbar:** nein
- **Klasse:** botschaft-ueber-messung (`BEO-ALL/commit-message-overclaims-work`, zweite Richtung)

### L-2 — `slice-292` nennt „Ende gelöscht und als andere Datei neu angelegt“ ohne die Commit-Bedingung, zwei Zeilen unter dem Fall, der sie umkehrt

- **Kategorie:** LOW
- **Quelle:** [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 11 („Löschen und
  Neuanlegen im **selben** Schritt ist … eine Inhaltsänderung“)
- **Pfad:** `docs/plan/planning/open/slice-292-vcs-umbenennung-modulkern.md` · „Ende gelöscht und als
  andere Datei neu angelegt; Zwischenpfad einer Kette neu belegt“ (Z. 158)
- **Befund:** Z. 155–156 (neu in `7d4643c4`) sagt für Löschen und Neuanlegen im selben Commit „der
  Nachzug folgt“, Z. 158 für „gelöscht und … neu angelegt“ ohne Commit-Angabe „Nachzug ist Befund“.
  Lastenheft Z. 2386 und Fitness Function Z. 234 tragen „in einem späteren“. Z. 158 selbst liegt
  außerhalb des Diffs; der Widerspruch entsteht durch die neue Nachbarzeile. Die Klasse steht in diesem
  Slice damit beim dritten Auftreten (R1 M-2, R2 M-2).
- **Verifizierbar:** nein
- **Klasse:** grenz-aussage-ohne-ihre-bedingung

### I-1 — Rotation: der „auch“-Satz zur Teilnahme lässt sich als eigener Ausschluss lesen

- **Kategorie:** INFO
- **Quelle:** Maintainability
- **Pfad:** `spec/lastenheft.md` · „Jeder andere Pfad nimmt nicht teil — auch ein“ (Z. 2296–2298);
  gleichlautend `spec/spezifikation.md` Z. 2155–2157
- **Befund:** A → X, danach C → A: A liegt in BASE und ist Ende von C → A, zugleich Quelle der
  teilnehmenden Umbenennung A → X. Nach dem Satzbau („Jeder andere Pfad …“) und nach ADR Entscheidung 9
  („hat dort eine eigene Geschichte … und wird nicht zum Ende einer anderen Datei“) nimmt A als Quelle
  teil. Der „auch“-Einschub beschreibt A aber wörtlich und liest sich ohne den Satzanfang als Ausschluss.
  Beide Lesarten öffnen kein Grün gegen den Lauf ohne den Schlüssel; ein Kriterium für den Fall gibt es
  nicht.
- **Verifizierbar:** nein
- **Klasse:** pfad-mit-zwei-rollen-in-der-abbildung (zweites Auftreten nach R2 M-1, diesmal ohne
  Widerspruch zwischen den Straten)

### I-2 — `4317dc0a` und `18157ed7` tragen im Kopf 0.114.0 mit dem Inhalt von 0.114.1

- **Kategorie:** INFO
- **Quelle:** [`MR-032`](../../harness/conventions.md#mr-032)
- **Pfad:** `spec/lastenheft.md` · Historie-Zeile „| 0.114.1 |“ (Z. 4769)
- **Befund:** Die Zeile 0.114.1 ordnet die Fixes nach R1 richtig zu; die beiden Commits, die sie
  geschrieben haben, nennen im Kopf weiter 0.114.0. Ohne Umschreiben der Historie ist das nicht zu
  ändern; bewusste Won't-Fix-Notiz.
- **Verifizierbar:** nein
- **Klasse:** version-im-kopf-vor-der-zeile

---

## Negativbefunde

### Widerspruchsfreiheit je Pflichtfall und Schlüsselkombination

Referenz ist `vcs.ignore-link-targets` allein (ilt). „fr“ = nur `vcs.follow-renames`, „beide“ = beide
Schlüssel. Gelesen gegen Absatz, Grenzen, alle 28 Kriterien, Spezifikation Schritt 2 und 4 und ADR.

| Fall | fr | beide | Stellen | Ergebnis |
|---|---|---|---|---|
| Tausch über einen Zwischennamen | kein Pfad nimmt teil, roh, Nachzug Drift | wie ilt (Drift bei anderem Dateinamen) | Lastenheft Z. 2296–2298, 2330, 2384; Spezifikation Z. 2155–2157; ADR Entscheidung 9, Z. 196–199 | ohne Befund im Vertrag; Tausch in einem Commit ist keine Umbenennung (Z. 2263–2266, Spezifikation Z. 2073–2077). Rest M-3 (Plan) |
| freigemachter und neu belegter Name | Drift (Kriterium Z. 2393) | wie ilt | Lastenheft Z. 2297–2298, 2330–2331; Spezifikation Z. 2157–2161; ADR Z. 196–199, Alternative M | ohne Befund; Richtung gegen die bestätigte Fassung: M-1 |
| zwei Quellen auf ein Ende | kein Pfad nimmt teil, Drift | wie ilt | Lastenheft Z. 2294–2295; Spezifikation Z. 2153–2154; ADR Z. 113–114 | gleichlautend; Rest M-4 (kein Kriterium), M-1 (Richtung) |
| Zwischenpfad | A → Ende, neues B roh, Drift (Z. 2387) | Drift auch bei gleichem Dateinamen | Lastenheft Z. 2282–2283; Spezifikation Z. 2088–2089; ADR Z. 112 | ohne Befund |
| Löschen und Neuanlegen im selben Commit | folgt (Z. 2394) | folgt | Lastenheft Z. 2278–2280, 2326, 2352–2353; Spezifikation Z. 2085–2086; ADR Z. 136–137, 215 | ohne Befund |
| … in getrennten Commits | „gelöscht“, keine Umbenennung, Drift (Z. 2386) | wie ilt, bei gleichem Dateinamen kein Befund (Z. 2390) | Lastenheft Z. 2283–2286, 2340–2344; Spezifikation Z. 2089–2093; ADR Z. 134–136, 218–220 | ohne Befund. Rest L-2 (Plan `slice-292`) |
| Merge-Seiten | unberührte Seite trägt nichts bei, folgt (Z. 2383); spätere Neuanlage ⇒ Mehrdeutigkeit (Z. 2385); Überschreiben an Ort und Stelle folgt | wie fr für teilnehmende, sonst wie ilt | Lastenheft Z. 2273–2278, 2351–2355; Spezifikation Z. 2080–2097; ADR Entscheidung 11, Z. 212–217 | ohne Befund; „ohne ihn zu löschen — auch durch Löschen und Neuanlegen“ (Z. 2352–2353) ist nach Z. 2278–2280 lesbar |
| Umbenennung nur des Verzeichnisses | Datei darin folgt; Verzeichnis-Link roh (Z. 2395) | Umbiegen auf andere gleichnamige Datei Drift (Z. 2392); Verzeichnis-Link wie ilt | Lastenheft Z. 2303–2305, 2336–2338; Spezifikation Z. 2161–2163, 2179–2181; ADR Z. 98–99, 204–205 | ohne Befund |
| Wurzel-Datei | Ende `x.md` wird zu NUL NUL `x.md` | nicht teilnehmend NUL `x.md`, teilnehmend NUL NUL `x.md` — verschieden im zweiten Zeichen | Spezifikation Z. 2141–2150, 2169–2171; Code: `normalizedTargetMark`, `markFreeTree` | ohne Befund. Ein Dateiname kann nicht mit NUL beginnen, und eine Fassung mit NUL wird nicht normiert — die Formen gleichen sich auch dort nicht, wo Pfad und Dateiname gleich sind |

### Übrige Bereiche

| Bereich | Ergebnis |
|---|---|
| Richtung (a) gegen den Lauf ohne `vcs.follow-renames` | geprüft, ohne Befund. fr allein: die Normalform ist eine Funktion des Ziels auf beiden Seiten, ein nicht teilnehmendes Ziel bleibt roh, eine rohe Fassung kann die doppelte Marke nicht tragen — kein Befund, den es ohne den Schlüssel nicht gibt; Ausnahme Exit 2 (Kriterien Z. 2388, 2399). Beide: nicht teilnehmende Ziele wie ilt, teilnehmende gleichen nur teilnehmenden mit gleichem Ende — ein Befund von ilt entfällt nur dort (Lastenheft Z. 2334–2336, Spezifikation Z. 2176–2179, ADR Z. 119–121) |
| Richtung (b) gegen die bestätigte Fassung | geprüft. Entscheidung 10 für sich wirkt nur in Richtung Drift (ADR Z. 258 und Plan Z. 217–218 stimmen); Entscheidung 9: M-1, M-2. Das Löschen und Neuanlegen im selben Commit ist gegen den Lastenheft-Text 0.114.1 nachsichtiger, gegen Spezifikation und ADR 11 nicht — `5b90b304` nennt es |
| Historie nach [`MR-032`](../../harness/conventions.md#mr-032) | geprüft, ohne Befund. 0.114.0 byte-gleich zu `f0622065`, 0.114.1 und 0.114.2 eigene Zeilen, Kopf 0.114.2; die Zählungen stimmen (0.114.1: vier neue Kriterien, 20 → 24; 0.114.2: vier neu, fünf angeglichen, 24 → 28). Die Spezifikation führt eine eigene Datumszeile (Z. 4318); MR-032 gilt dort nicht. Rest I-2 |
| Referenz-Richtung ([`AGENTS.md`](../../AGENTS.md) §3.4) | geprüft, ohne Befund. In den hinzugefügten Zeilen beider Spec-Straten kein `ADR-`, `slice-`, `welle-`, kein Commit-Hash, kein Link nach `docs/plan/`, kein „R1“/„R2“ (grep über die `+`-Zeilen leer); „Anlass: Review“ ist die gelebte Form |
| ADR-Form ([`AGENTS.md`](../../AGENTS.md) §3.5) | geprüft, ohne Befund. `Proposed`; zwei neue Geschichte-Zeilen mit Provenance-Marker als Herkunft; Alternativen L und M führen die verworfenen Fassungen. Rest M-2 |
| Grenz-Aussagen gegen den Code (Regel 13) | geprüft, ohne Befund. `AllPaths` liest ohne Schlüssel die Commit-Objekte der Basis-Vorfahren und die Bäume von BASE und HEAD bzw. jeden Index-Eintrag — „ohne den Schlüssel liest das Modul dieses Objekt nicht“ (Z. 2399) und „kein Commit zwischen BASE und HEAD“ (Z. 2400) stimmen; die Form für `vcs.ignore-link-targets` allein (Marke, `path.Base`, Anker) gleicht Spezifikation Z. 2141–2143 |
| Commit-Disziplin: Kennungen | geprüft, ohne Befund. Jede Botschaft nennt `slice-291`, die meisten [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in) und/oder [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md); `7d4643c4` auch `slice-292` |
| Commit-Disziplin: Reihenfolge | geprüft, ohne Befund. ADR (`e84632b1`) → Bestätigung in den Plänen (`7d4643c4`) → Plan-Änderung (`8bcbbc67`) → Lastenheft (`e3e5d4eb`) → Spezifikation (`5b90b304`) → Bestätigung in der ADR (`d03cda9c`) → Richtungs-Berichtigung (`7a57fa1e`) |
| Commit-Disziplin: Pfad-Begrenzung | geprüft, ohne Befund. Je Commit nur die Dateien seiner Rolle; die ADR-Datei liegt nur in `e84632b1` und `d03cda9c`; kein `git mv` in der Range |
| Commit-Disziplin: Botschaften (Regel 15) | geprüft. `7d4643c4`, `8bcbbc67`, `5b90b304`, `d03cda9c` decken sich mit ihrem Diff; Rest L-1 und, für `7a57fa1e`, M-1 |
| Abgrenzung des Plans (§1) | geprüft, ohne Befund. Kein Code, kein `.d-check.yml`, kein README/CHANGELOG/Handbuch |
| Aussagen-Suche nach Inhalt | geprüft über Lastenheft, Spezifikation, ADR und die drei Pläne: „nie grüner“, „zugleich Quelle“, „Dateiname(n) des Endes“, „erst gleich…“, „wechselt den Dateinamen“, „strenger als“, „nur noch in BASE“, „Zwischen-Commit“, „Neuanlegen“/„neu an(gelegt)“, „zweite(n) Quelle“ — Treffer außerhalb der Historie-Zeilen sind die oben belegten oder Plan-Zitate der Suche selbst; Rest L-2 |

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 4 | M-1, M-2, M-3, M-4 |
| LOW | 2 | L-1, L-2 |
| INFO | 2 | I-1, I-2 |

## Verdikt

**Nachbessern vor der Verifikation.** Kein HIGH. Zehn der elf R2-Befunde sind geschlossen oder in der
Sache überholt; I-3 ist offen und als M-3 verschärft. Die Spec-Straten und die ADR sind für alle
Pflichtfälle und alle drei Schlüsselkombinationen widerspruchsfrei, die Historie folgt MR-032, die
Referenz-Richtung hält. Offen sind zwei Dinge. Das Erste ist der Nachweis der Richtung gegen die
bestätigte Fassung: Der Plan nennt das falsche Ziel und lässt zwei Fälle aus, einen davon mit
`vcs.follow-renames` allein (M-1). Die ADR verzeichnet die Bestätigung ohne den Vorbehalt, unter dem sie
gegeben wurde (M-2). Das geht an Planner und Architect; ob die Bestätigung trägt, entscheidet der
Auftraggeber. Das Zweite ist die Belegbarkeit für `slice-292`: Der Tausch-Testfall fährt seine Regel
nicht (M-3), und die neue Klausel „Ende keiner zweiten Quelle“ hat weder ein Kriterium noch einen
Testfall (M-4). Plan-Risiko 1 bleibt **eingetreten**. Drei Klassen stehen in diesem Slice beim zweiten
oder dritten Auftreten: `grenz-aussage-ohne-ihre-bedingung` 3×, `botschaft-ueber-messung` 3×,
`bestaetigung-ueber-geaenderten-wortlaut` 2×. Die ersten beiden sind ein Steering-Loop-Signal für die
Closure.
