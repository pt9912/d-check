# Review-Report: slice-291 — `vcs` folgt einer Umbenennung, Fix-Commits nach R3 (R4) — 2026-11-24

**Review-Art:** Design. Geprüft wird der Diff gegen Slice-Plan, ADR und Hard Rules, also Maintainability
(`v6.18.1` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten). Gegen die DoD prüft danach der
Verifier.

**Gegenstand:** `slice-291`, Commit-Range `9ccfcda4..e2819614`. Die sechs Commits liegen nur lokal auf
`main`:

- `90fd2340` und `e2819614`: Planner. Pläne `slice-291` und `slice-292`, Welle `welle-96`.
- `b297d9a4`: Architect, [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md), weiter `Proposed`.
- `48e81131`: Plan-Änderung des Implementers.
- `10788cf7`: Lastenheft
  [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in),
  Version 0.114.3.
- `81622ae5`: Spezifikation
  §[`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
  Schritt 4, Teilnahme.

Nebeneinander gelesen wurden im Stand `e2819614`:

- Absatz „Was gleichgesetzt wird“, Grenzen und alle 30 Kriterien (Lastenheft Z. 2292–2418);
- Spezifikation Schritt 2 und 4 und die Schema-Zeile;
- ADR Entscheidung 9 und 10, §Konsequenzen, Fitness Function und Geschichte;
- die Testfall-Liste von `slice-292`.

Gegengelesen wurde die zuerst bestätigte Fassung von Entscheidung 9 und 10 im Stand `c89e2dee`.

**Skill:** `.harness/skills/reviewer.md` @ 1.20.0 (`d1c37c14`).
**Modell:** claude-opus-5-5 · **Datum:** 2026-11-24.

**Eingangs-Kontext:**

- Report R3 ([`2026-11-23-slice-291-vcs-umbenennung-r3.md`](2026-11-23-slice-291-vcs-umbenennung-r3.md))
  samt seiner Widerspruchstabelle; R1 und R2 über dessen Statustabelle.
- Plan `slice-291` §3, Entscheid 9 und 10 und alle Plan-Änderungen; Plan `slice-292` mit der
  Testfall-Liste; Welle `welle-96` §5.
- [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md).
- [`MR-032`](../../harness/conventions.md#mr-032).
- [`AGENTS.md`](../../AGENTS.md) §3.3, §3.4, §3.5, §5 Regeln 13, 15 und 16, §6 Schritt 4.
- `v6.18.1` · `regelwerk/modul-11-verification.md` §Bewusstes Brechen für DoD-Testbehauptungen.

**Lesart „nachsichtiger“.** Die drei Fälle wurden je Schlüsselkombination gegen die zuerst bestätigte
Fassung ausgewertet. Die ist `c89e2dee`: Teil hat eine Quelle mit Ende in HEAD oder ein solches Ende,
zugleich Quelle und Ende nicht; mit beiden Schlüsseln wird ein teilnehmendes Ziel auf den Dateinamen
des Endes normiert. Das Ergebnis, auf dem M-1 und L-1 aufbauen, zeigt die Tabelle. „fr“ steht für
`vcs.follow-renames` allein, „beide“ für beide Schlüssel.

| Fall | fr: gegen die erste Fassung | beide: gegen die erste Fassung |
|---|---|---|
| (1) A gelöscht, danach C → A | **strenger**: Nachzug C → A war grün, ist Drift (Kriterium Z. 2409) | nachsichtiger: Umbiegen von C auf eine andere `…/c.md` war Drift, ist grün |
| (2) A → E und B → E auf zwei Seiten eines Merges | **strenger**: Umbiegen A → B war grün, ist Drift (Z. 2410) | nachsichtiger: Umbiegen von A auf eine andere `…/a.md` war Drift, ist grün |
| (3) Rotation A → X, danach C → A | nachsichtiger: Nachzug A → X war Drift (A nahm als Quelle und Ende nicht teil), folgt jetzt (Z. 2411). Für C ändert sich nichts: C → A war Drift, weil A roh blieb, und ist es weiter | nachsichtiger für C wie in (1); Nachzug A → X folgt; Nachzug C → A war grün und ist Drift |

---

## Stand der R3-Befunde

| R3 | Stand | Beleg |
|---|---|---|
| M-1 (Richtung am falschen Ziel, zwei Fälle fehlen) | **teilweise geschlossen**, Rest M-1 unten | Alle drei Fälle sind jetzt benannt, im Plan, im Welle-Plan, in der ADR und im Lastenheft. Richtig beschreiben sie die ADR (Z. 205–217) und das Lastenheft (Z. 2338–2349). Der Plan (Z. 203–212) nennt bei der Rotation die falsche Hälfte und bei (1) und (2) nicht den Schlüssel, unter dem der Fall nachsichtiger ist |
| M-2 (Bestätigung ohne Vorbehalt in der ADR) | **geschlossen** | Vorspann Z. 78–83: die Vorlage vom 2026-11-23 hält in drei Fällen nicht, die Bestätigung vom 2026-11-24 gilt mit diesen Fällen und verweist auf §Konsequenzen. Dort steht die Grenze Z. 205–217, dazu die Geschichte-Zeile Z. 278. Rest I-1 |
| M-3 (Tausch-Testfall fährt seine Regel nicht) | **in der Sache geschlossen**, Rest L-3 | `slice-292` Z. 148–151: Tausch über einen Zwischennamen in drei Commits; der Tausch der Inhalte in einem Commit ist als Inhaltsänderung benannt. Die BASE-Bedingung bricht verlässlich der Testfall Z. 152 (freigemachter Name). Der neue Rot-Beleg am Tausch ist nicht aus der Regel ableitbar (L-3) |
| M-4 (zweite Quelle ohne Kriterium, Testfall, Grenze) | **geschlossen** | Kriterium Z. 2410, mit oder ohne `vcs.ignore-link-targets`, verschiedene Dateinamen. Grenze Z. 2345–2347. `slice-292` Z. 155–157 mit derselben Bedingung (`e2819614`). Fitness Function Z. 251. Ohne die Klausel nähmen A und B mit Ende E teil, und das Umbiegen wäre grün; der Testfall fährt also die Klausel |
| L-1 (Botschaft `e84632b1` ohne Widerruf) | **geschlossen** | Botschaft `b297d9a4`: „Berichtigt die Botschaft von e84632b1: ‚nur in Richtung Drift‘ traegt fuer Entscheidung 9 nicht.“ Dazu die ADR-Geschichte Z. 278 |
| L-2 (Neubelegung ohne Commit-Bedingung in `slice-292`) | **geschlossen**, Klasse an anderer Stelle wiederholt | `slice-292` Z. 165–167: „in einem Commit gelöscht und in einem späteren“. Die Inhalts-Suche fand dieselbe Klasse in der ADR (L-2 unten) |
| I-1 (Rotation, „auch“-Satz als Ausschluss lesbar) | **geschlossen** | Lastenheft Z. 2296–2302 und Spezifikation Z. 2155–2161 gleichlautend, Kriterium Z. 2411, ADR Entscheidung 9 Z. 118–122. Zur Eindeutigkeit siehe die Negativbefunde |
| I-2 (Kopf 0.114.0 in zwei Commits) | **unverändert**, Won't-Fix | Nicht wiederholt: `10788cf7` hebt Kopf und Zeile 0.114.3 im selben Commit |

---

## Findings

### M-1 — Der Plan beschreibt die drei nachsichtigeren Fälle anders als ADR und Lastenheft, und auf diese Beschreibung stützt er die Bestätigung

- **Kategorie:** MEDIUM
- **Quelle:** [`AGENTS.md`](../../AGENTS.md) §5 Regel 15; Prüffrage 8;
  [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 9 und §Konsequenzen „Grenze“
- **Pfad:** `docs/plan/planning/in-progress/slice-291-vcs-umbenennung-anforderung-und-adr.md` ·
  „dort gilt das Ergebnis wie ohne den Schalter, nicht Drift“ und „C nimmt nicht teil, auch mit
  `follow-renames` allein“ (Z. 205–212)
- **Befund:** Bei der Rotation hängt der Plan „auch mit `follow-renames` allein“ an C. Mit dem Schlüssel
  allein ändert sich für C gegen die erste Fassung aber kein Ergebnis. Nachsichtiger wird dort der
  Nachzug A → X: er war Drift und folgt jetzt. Die ADR (Z. 211–213) und das Lastenheft (Z. 2348–2349)
  nennen genau das, der Plan nicht. Bei (1) und (2) gilt die Nachsicht nur mit beiden Schlüsseln. Mit
  dem Schlüssel allein sind beide Fälle strenger als die erste Fassung: der Nachzug C → A und das
  Umbiegen A → B waren grün und sind Drift (Tabelle oben). Der Rahmensatz „dort gilt das Ergebnis wie
  ohne den Schalter, nicht Drift“ ist für den Lauf mit dem Schlüssel allein in (1) falsch, denn ohne den
  Schalter ist der Nachzug C → A Drift (Kriterium Z. 2409). Für A → X in (3) ist er ebenfalls falsch: das
  Ergebnis trägt eine teilnehmende Umbenennung, nicht der Lauf ohne den Schalter.
- **Failure-Szenario:** Der Plan hält fest, die Fassung sei am 2026-11-24 „mit diesen drei Fällen“
  bestätigt. Der Welle-Plan (Z. 114–116) verweist für die drei Fälle auf diesen Abschnitt. Damit gibt es
  zwei Aufzeichnungen dessen, was bestätigt wurde: die ADR (§Konsequenzen) und den Plan. Sie
  unterscheiden sich genau in dem Fall, den R3 M-1 als nicht vorgelegt benannt hat, der Nachsicht mit
  dem Schlüssel allein. Liest der Verifier den Plan gegen `slice-292` oder der Retirement-Check die
  Herkunft, so ist nach dem Plan die Rotation mit dem Schlüssel allein für C nachsichtig. Das
  Kriterium Z. 2411 verlangt für C → A Drift, und die tatsächliche Nachsicht (A → X) steht im Plan
  nicht. Ob die Bestätigung auf die ADR-Fassung oder auf die Plan-Fassung hin gegeben wurde, geht aus
  dem Repo nicht hervor.
- **Warum nicht HIGH:** Es ist kein Code betroffen. ADR und Lastenheft sind richtig und untereinander
  gleich, und die ADR bindet die Bestätigung an ihre eigene §Konsequenzen-Stelle (Z. 82–83). Der Plan
  wandert später in den Archiv-Bestand.
- **Verifizierbar:** nein
- **Klasse:** richtungs-aussage-am-falschen-ziel (zweites Auftreten nach R3 M-1)
- **Träger:** Planner (Plan §3). Ob die Bestätigung trägt, entscheidet der Auftraggeber.

### L-1 — Die Fitness Function sagt für alle drei Fälle mit beiden Schaltern „dasselbe Ergebnis wie mit `ignore-link-targets` allein“, und `slice-292` führt dazu keinen Testfall

- **Kategorie:** LOW
- **Quelle:** [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) §Konsequenzen „Folgepflicht“ („Jede
  Randbedingung hat einen Test“);
  [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  Kriterium „Boundary (Rotation)“
- **Pfad:** `docs/plan/adr/0115-vcs-folgt-umbenennung.md` · „die drei Fälle der Grenze zu Entscheidung 9
  mit beiden Schaltern ⇒ dasselbe Ergebnis wie mit `ignore-link-targets` allein“ (Z. 252);
  `docs/plan/planning/open/slice-292-vcs-umbenennung-modulkern.md` · Testfall-Liste (Z. 130–181)
- **Befund:** Für die nicht teilnehmende Quelle gilt der Satz. Für die Rotation als ganzen Fall gilt er
  nicht: mit beiden Schaltern folgt der Nachzug A → X (A nimmt teil). Mit `ignore-link-targets` allein
  ist er Drift, weil der Dateiname wechselt. Das Lastenheft fasst das mit „oder eine teilnehmende
  Umbenennung mit demselben Ende trägt es“ (Z. 2339–2340), die Fitness-Zeile nicht. `slice-292` hat
  keinen Testfall für die nachsichtige Hälfte mit beiden Schaltern. Das betrifft das Umbiegen von C bzw.
  A auf eine gleichnamige andere Datei, das die Grenze Z. 2342–2347 als „kein Befund“ zusagt. Es gibt nur
  die allgemeine Zeile Z. 160–161.
- **Verifizierbar:** ja, ab `slice-292`. Mit dem Fixture Rotation, beiden Schaltern und dem Nachzug A → X
  widerspricht die Fitness-Zeile dem Kriterium Z. 2411.
- **Klasse:** grenz-aussage-ohne-ihre-bedingung

### L-2 — Die Probe der Fitness Function sagt „gelöscht und neu angelegt ⇒ Befund“ ohne die Bedingung, unter der das gilt

- **Kategorie:** LOW
- **Quelle:** [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 1 (Inhalts-Gleichheit)
  und 11 (Löschen und Neuanlegen im selben Schritt)
- **Pfad:** `docs/plan/adr/0115-vcs-folgt-umbenennung.md` · „gelöscht und neu angelegt ⇒ Befund“
  (Z. 253, Zeile „Probe am gebauten Image“)
- **Befund:** Nach Entscheidung 1 ist eine gelöschte Datei, die unter neuem Namen mit demselben Inhalt
  angelegt wird, eine Umbenennung, und der Nachzug folgt ihr. Nach Entscheidung 11 ist Löschen und
  Neuanlegen am selben Pfad im selben Commit eine Inhaltsänderung. Ein Befund entsteht nur mit
  **anderem** Inhalt (Kriterium Z. 2395) oder bei der Neubelegung über zwei Commits. Die Zeile nennt
  keine der beiden Bedingungen; `slice-292` Z. 135 und der Welle-Plan Z. 70 tragen „mit anderem
  Inhalt“. Die Zeile liegt außerhalb des Diffs und wurde bei der Inhalts-Suche zu R3 L-2 gefunden.
- **Verifizierbar:** ja, ab `slice-292`. Wird die Probe mit gleichem Inhalt gebaut, folgt der Nachzug,
  und die Probe ist rot, obwohl der Code stimmt.
- **Klasse:** grenz-aussage-ohne-ihre-bedingung (in diesem Slice nach R1 M-2, R2 M-2 und R3 L-2 jetzt
  zwei weitere Stellen)

### L-3 — Der Rot-Beleg am Tausch setzt einen Ausgang voraus, den die Regel ohne die BASE-Bedingung nicht festlegt

- **Kategorie:** LOW
- **Quelle:** `v6.18.1` · `regelwerk/modul-11-verification.md` §Bewusstes Brechen für
  DoD-Testbehauptungen
- **Pfad:** `docs/plan/planning/open/slice-292-vcs-umbenennung-modulkern.md` · „Rot-Beleg: ohne die
  BASE-Bedingung wäre der Nachzug grün“ (Z. 150)
- **Befund:** Ohne die BASE-Bedingung nehmen A → B (über T) und B → A beide teil. Dann ist A Quelle der
  einen und Ende der anderen Umbenennung, und für B gilt dasselbe. „Quelle und Ende gelten beide als das
  Ende“ ergibt so für A und für B je zwei Normalformen. Ob der Nachzug A → B dann grün oder Drift ist,
  hängt davon ab, wie der Code die Doppelrolle auflöst. Mit der Regel der ersten Fassung („zugleich
  Quelle und Ende nimmt nicht teil“) bliebe er Drift, und der Test bliebe grün. Den Bruch der
  BASE-Bedingung fängt verlässlich der Testfall Z. 152.
- **Verifizierbar:** ja, ab `slice-292`: die Mutation „BASE-Bedingung entfernt“ gegen das Tausch-Fixture
  fahren.
- **Klasse:** test-begruendung-faehrt-ihre-regel-nicht (drittes Auftreten nach R2 I-3 und R3 M-3)

### I-1 — Die neue Geschichte-Zeile schreibt der Zeile zu Review R2 die Aussage für beide Änderungen zu

- **Kategorie:** INFO
- **Quelle:** [`AGENTS.md`](../../AGENTS.md) §5 Regel 16
- **Pfad:** `docs/plan/adr/0115-vcs-folgt-umbenennung.md` · „Die Zeile vom 2026-11-23 zu Review R2 und
  die Botschaft des Commits, der sie schrieb, nannten beide Änderungen“ (Z. 278)
- **Befund:** Für die Botschaft `e84632b1` stimmt das: „Beide Aenderungen wirken nur in Richtung
  Drift“. Die Zeile Z. 276 setzt „nur in Richtung Drift“ dagegen in den Abschnitt zu Entscheidung 10,
  durch Semikolon vom Abschnitt zu Entscheidung 9 getrennt; R3 M-2 hat sie so gelesen. Die Zeile Z. 278
  legt die weitere Lesart fest, und mit `Accepted` wird diese Lesart unveränderlich.
- **Verifizierbar:** nein
- **Klasse:** quelle-ueber-geltungsbereich

---

## Negativbefunde

### Teilnahme-Satz gegen die R3-Widerspruchstabelle

Gelesen wurden Lastenheft Z. 2292–2312 und Spezifikation Z. 2152–2173. Beide Fassungen sind
gleichlautend.

| Fall aus R3 | Neuer Satz | Ergebnis |
|---|---|---|
| Tausch über einen Zwischennamen | Beide Umbenennungen enden in BASE; kein Pfad ist Quelle einer teilnehmenden | unverändert, kein Pfad nimmt teil |
| freigemachter und neu belegter Name | C → A endet in BASE; A ist Quelle keiner Umbenennung | unverändert, weder C noch A nimmt teil (Grenze Z. 2341–2344) |
| zwei Quellen auf ein Ende | „über sie“: weder A noch B noch E nimmt teil; E ist nicht in BASE und deshalb keine Quelle | unverändert; jetzt mit Kriterium |
| Zwischenpfad | Das neue B ist weder Quelle noch Ende | unverändert, Drift (Z. 2403) |
| Löschen und Neuanlegen im selben Commit, getrennte Commits, Merge-Seiten | nicht berührt (Schritt 2) | unverändert |
| Verzeichnis, Wurzel-Datei | nicht berührt | unverändert |
| Rotation | A nimmt als Quelle von A → X teil und gilt als X; C nimmt nicht teil | jetzt eindeutig; Kriterium Z. 2411 |

Geprüft ohne Befund. Kein Pflichtfall kippt. Die Abbildung bleibt eine Funktion des Pfads: eine Quelle
liegt in BASE, ein teilnehmendes Ende nicht, also kann kein Pfad zugleich teilnehmende Quelle und
teilnehmendes Ende sein. Zwei teilnehmende Enden einer Quelle schließt die Mehrdeutigkeit aus, zwei
Quellen eines Endes die neue Klausel. „Ist ein Pfad aus BASE Ende einer solchen …“ deckt beide Arten
nicht teilnehmender Umbenennungen ab. Bei der zweiten Art liegt das Ende nicht in BASE, sie kommt also
nur zusammen mit der ersten in Frage.

### Übrige Bereiche

| Bereich | Ergebnis |
|---|---|
| Zusage „kein Befund entfällt, den es ohne `follow-renames` gibt, außer über eine teilnehmende Umbenennung“ | Geprüft, ohne Befund, für Plan Z. 204–205, Lastenheft Z. 2350–2355 und Spezifikation Z. 2168–2187. Mit dem Schlüssel allein bleiben nicht teilnehmende Ziele roh, und ein teilnehmendes Ziel gleicht nur einem teilnehmenden mit demselben Ende. Mit beiden Schlüsseln werden nicht teilnehmende Ziele wie mit `vcs.ignore-link-targets` allein normiert. Ein Befund entfällt also nur über ein gemeinsames Ende. Einzeln gegengehalten: der Nachzug A → X in der Rotation, das Umbiegen X → A in der Rotation (A gilt in HEAD als X) und die nachsichtigen Fälle (1) und (2) mit beiden Schlüsseln. Alle stehen innerhalb der Zusage |
| Drei Fälle in vier Artefakten (Schwerpunkt 1) | Die Fälle heißen überall gleich: Plan Z. 207–209 und 386, Welle Z. 114–116, ADR Z. 205–217, Lastenheft Z. 2338–2349. Als Testfall in `slice-292` stehen sie in Z. 152, 153–154 und 155–157, im Lastenheft als Kriterium Z. 2409–2411 und als Grenze. Der Inhalt weicht ab: M-1 (Plan) und L-1 (Fitness Function, keine nachsichtigen Testfälle mit beiden Schlüsseln) |
| Lastenheft-Grenze zu den drei Fällen | Geprüft, ohne Befund. „liefert `vcs.ignore-link-targets` allein ebenso, oder eine teilnehmende Umbenennung … trägt es“ stimmt je Fall mit der Tabelle oben überein. Bei (1) und (2) steht „Mit beiden Schlüsseln“ ausdrücklich |
| Historie nach [`MR-032`](../../harness/conventions.md#mr-032) (Schwerpunkt 4) | Geprüft, ohne Befund. Der Kopf steht auf 0.114.3 (Z. 3), die Zeile 0.114.3 ist eigen (Z. 4786), die Verweis-Spalte trägt `—`. Die Zeile 0.114.2 ist byte-gleich, denn der Diff fügt nur oberhalb ein. Die Zählung stimmt: zwei neue Kriterien, 28 → 30 (Z. 2389–2418). Der Patch-Rang folgt 0.114.1 und 0.114.2 als Präzisierung nach Review. Eine Spiegelung der Versionsnummer gibt es außerhalb der Historie nicht (Suche nach `0.114.` ohne `done/` und `docs/reviews/`). Die Spezifikation hat eine eigene Datumszeile (Z. 4323), für sie gilt MR-032 nicht |
| Referenz-Richtung ([`AGENTS.md`](../../AGENTS.md) §3.4) | Geprüft, ohne Befund. In den hinzugefügten Zeilen beider Spec-Straten stehen weder `ADR-`, `slice-`, `welle-` noch ein Commit-Hash, ein Link nach `docs/plan/`, „R3“, „Auftraggeber“ oder „bestätigt“; die Suche über die `+`-Zeilen ist leer. Die Grenze sagt „das man als Drift erwarten könnte“ und verweist nicht auf die erste Fassung |
| ADR-Form ([`AGENTS.md`](../../AGENTS.md) §3.5) | Geprüft, ohne Befund. Die ADR bleibt `Proposed`; die neue Geschichte-Zeile trägt den Provenance-Marker als Herkunftsangabe und begründet keine Entscheidung. Rest I-1 |
| Commit-Disziplin: Pfad-Begrenzung (Schwerpunkt 5) | Geprüft, ohne Befund. Jeder Commit berührt nur die Dateien seiner Rolle: `90fd2340` die drei Pläne, `48e81131` den Plan, `b297d9a4` die ADR, `10788cf7` das Lastenheft, `81622ae5` die Spezifikation, `e2819614` `slice-292`. Kein `git mv` |
| Commit-Disziplin: Reihenfolge | Geprüft, ohne Befund. Bestätigung in den Plänen, dann Plan-Änderung des Implementers vor dem Spec-Text ([`AGENTS.md`](../../AGENTS.md) §6 Schritt 4), dann ADR, Lastenheft, Spezifikation und Nachtrag `slice-292` |
| Commit-Disziplin: Botschaften (Regel 15) | Geprüft, ohne Befund. `48e81131`, `10788cf7` (zwei Kriterien, 30), `81622ae5` und `e2819614` decken sich mit ihrem Diff. `b297d9a4` berichtigt `e84632b1` ausdrücklich. `90fd2340` sagt, was geändert wurde, nicht, dass M-1 geschlossen sei; der Rest liegt in M-1 |
| Commit-Disziplin: Kennungen | Geprüft, ohne Befund. Jede Botschaft nennt `slice-291` oder `slice-292` und [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in); `b297d9a4` und `90fd2340` nennen auch [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) |
| Abgrenzung des Plans (§1) | Geprüft, ohne Befund. Kein Code, keine `.d-check.yml`, kein README, CHANGELOG oder Handbuch |
| Aussagen-Suche nach Inhalt | Geprüft über Lastenheft, Spezifikation, ADR und die drei Pläne. Gesucht wurde nach „zweiten Quelle“, „genau einer Quelle“, „in BASE lag“, „Rotation“, „freimacht“, „zwei Quellen“, „nachsichtiger“, „auch mit `follow-renames` allein“, „wie ohne den Schalter“, „nur in Richtung Drift“, „neu angelegt“, „Neuanlegen“ und „neu belegt“. Treffer außerhalb der Historie-Zeilen sind die oben belegten Stellen oder Zitate der Suche selbst. Neu gefunden: L-2 (ADR Z. 253). Der Plan führt „nur in Richtung Drift“ bei Entscheid 10 (Z. 220–221) weiter, dort zutreffend |

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 1 | M-1 |
| LOW | 3 | L-1, L-2, L-3 |
| INFO | 1 | I-1 |

## Verdikt

**Kleine Nachbesserung vor der Verifikation, am Plan.** Es gibt kein HIGH. Von den acht R3-Befunden sind
sechs geschlossen. M-1 ist teilweise geschlossen, I-2 ist ein unveränderter Won't-Fix.

Die Spec-Straten und die ADR sind für alle Pflichtfälle widerspruchsfrei, und der neue Teilnahme-Satz
lässt sich eindeutig lesen: die Rotation ist geregelt, kein Fall aus der R3-Tabelle kippt. Die Zusage
„außer über eine teilnehmende Umbenennung“ hält. Die Historie folgt MR-032, die Referenz-Richtung hält,
und `e84632b1` ist berichtigt.

Offen bleibt M-1. Der Plan beschreibt die drei Fälle anders als ADR und Lastenheft: bei der Rotation
fehlt der Nachzug A → X mit dem Schlüssel allein, bei (1) und (2) fehlt „mit beiden Schaltern“. Gerade
mit dieser Beschreibung hält der Plan die Bestätigung vom 2026-11-24 fest. Das geht an den Planner. Ob
die Bestätigung auf die ADR-Fassung gegeben wurde, kann nur der Auftraggeber sagen; war es die
ADR-Fassung, genügt eine Berichtigung des Plantexts.

L-1 bis L-3 betreffen die Fitness Function und die Testfall-Liste von `slice-292`. Sie sind vor dessen
erstem Edit behebbar und blockieren die Verifikation von `slice-291` nicht.

Zwei Klassen haben in diesem Slice die Steering-Loop-Schwelle überschritten:
`grenz-aussage-ohne-ihre-bedingung` mit fünf Stellen über vier Runden und
`test-begruendung-faehrt-ihre-regel-nicht` mit drei. Beide gehören in die Closure-Notiz.
