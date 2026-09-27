# Slice slice-231: `AGENTS.md` — Chronik und Planungs-Verweise aus den Hard Rules streichen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** — **wellenlos**. Die Closure-Bedingung geht über die eigene DoD
nicht hinaus — kein repo-weiter Beleg, den dieser Slice allein nicht liefert
(Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).

**Bezug:** — kein `DC-*`, keine aktive ADR. Harness-Meta-Dokumentation:
[`AGENTS.md`](../../../../AGENTS.md) §3.7 (Zustandsfelder tragen Zustand und
Anker, nicht die Chronik) und die Adaptionen
[MR-015](../../../../harness/conventions/MR-015-agents-md-routet.md) (AGENTS.md
routet, spiegelt nicht),
[MR-045](../../../../harness/conventions/MR-045-slice-verweise-nicht-im-briefing.md)
und [MR-050](../../../../harness/conventions/MR-050-herkunfts-anker-ist-kein-verweis.md)
(keine Slice-Verweise; der Herkunfts-Anker ist ausgenommen).

**Berührte Spec-Stellen:** — (`AGENTS.md` ist Rang 8 der Source Precedence,
kein Spec-Stratum).

**Verantwortlich:** claude-sonnet-5.

**Autor:** claude-sonnet-5. **Datum:** 2026-09-27.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Aus den Hard Rules in [`AGENTS.md`](../../../../AGENTS.md) die
Sätze entfernen, die erzählen, wie eine Regel entstand oder was früher galt,
und die beiden Links auf Planungs-Artefakte; Regel, Grenze und
Herkunfts-Anker bleiben unverändert stehen. Anlass ist die Auftraggeber-
Beobachtung, die Datei (rund 38 KB, jeder Lauf lädt sie) trage Chronik — etwa
den Satz „Seit slice-172 hält das ein Sensor…" in §5.

**Form der Chronik, ausgeschrieben vor der Messung.** Ein **Chronik-Satz**
beantwortet nicht, *was gilt*, sondern *wie es dazu kam*: Er erzählt die
Genese einer Regel, eine frühere Fassung, einen Messvorfall als Herleitung,
oder er verweist per Link auf ein Planungs-Artefakt (`slice-<NNN>`-Pfad). Die
zwei Tests aus `AGENTS.md` §3.7 tragen das Urteil: **Adressat** (wer liest den
Satz, um zu handeln?) und **Zeitform** (Indikativ über das, was ist). **Kein**
Chronik-Satz: ein Zustand, eine Grenze (auch eine *Bestandsgrenze*), ein
Rang-Zeiger, ein Herkunfts-Anker der Form `seit slice-<N>` / `seit welle-<N>`
in der `(Hard Rule aus dem Steering Loop, …)`-Klammer. Eine Grenze, die in
einem Chronik-Absatz steht, wird herausgelöst und bleibt.

**Kandidaten aus einer Stichwort-Suche** (nur ein Suchraster, kein Umfang —
sie findet Wörter, nicht Sätze; gelesen wird §3 und §5 vollständig):

- §3.3 „Historische Klärung: Diese Datei nannte hier früher eine eigene
  ‚Ausnahme'…" samt Link auf einen Slice (Verstoß gegen die Regel „keine Slice-Verweise“).
- §3.9 „gemessen am Tag-Push von `v0.66.0`, während dieses Gate grün
  meldete" — die Regel selbst trägt bereits [ADR-0071](../../adr/0071-lokale-workflow-referenz-rechte-pruefung.md).
- §5 „Seit slice-172 hält das ein Sensor…" und der Absatz dahinter, darin
  ein Link auf `slice-170` (derselbe Regelverstoß) und dieselbe Aussage über den
  Fenced-Block **zweimal** („INNERHALB eines wohlgeformten Fenced-Blocks" /
  „IM Fenced-Block") — ein Defekt, keine Chronik.
- §5 „Gemessen an sieben Fundstellen: In sechs …" (Grenzen-Regel) und
  „Gemessen, nicht vereinbart" (CHANGELOG-Regel).
- Die Herleitung vor den fünf `(Hard Rule aus dem Steering Loop, …)`-
  Klammern; der Anker in der Klammer bleibt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Regeltext nach `harness/conventions/` oder `harness/sensors/` auslagern**
  (§3.1-Netz-Targets, §3.3-Sonderfall, §3.9-Rechte-Absatz) — das wäre ein
  anderer Vorgang: er bewegt Regeln statt Prosa zu streichen und braucht eine
  Inventur der Rückverweise (Überschriften-Anker von `AGENTS.md` stehen in
  [`.d-check.yml`](../../../../.d-check.yml), in Hooks, Workflows und
  Adaptionen; [MR-025](../../../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  verlangt die Spiegel-Liste **vor** dem Editieren). **Diese Grenze hat noch
  keine Adresse:** ein Folge-Slice wird erst geschnitten, wenn die Messung
  dieses Slice zeigt, wie viel Größe danach bleibt.
- **Überschriften umbenennen, verschieben oder zusammenlegen** — acht
  Überschriften-Anker werden von Dokumenten außerhalb der eingefrorenen
  Bestände verlinkt (etwa §3.5 und §3.6); jede Änderung risse dort Links.
  Bestand bleibt bewusst stehen.
- **Eine Regel inhaltlich lockern oder streichen** — Streichen einer Hard
  Rule ist an ihren Auflösungs-Trigger gebunden (Modul 13 §Hard Rule), und für
  Regeln mit Herkunfts-Anker gilt der Retirement-Check. Dieser Slice entfernt
  nur, was **keine** Regel ist.
- **`harness/README.md`, `harness/conventions.md` und das vendorte Regelwerk**
  — dieselbe Beobachtung (Größe) gilt dort, ist aber ein anderer Gegenstand
  mit eigener Messung; ein Sammel-Slice sprengte die Ein-Sitzungs-Review.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

- [x] Alle Chronik-Sätze nach der Form aus §1 sind aus `AGENTS.md` §3 und §5
      entfernt; §7 führt sie **einzeln** (Fundstelle, Satz-Anfang, Test, der
      ihn als Chronik ausweist) und nennt, welche Grenze aus einem
      Chronik-Absatz herausgelöst wurde und stehen blieb.
- [x] `AGENTS.md` trägt **keinen** Link mehr auf ein Planungs-Artefakt
      (`slice-`- oder `welle-`-Pfad); Herkunfts-Anker `(seit …)` bleiben, ihre
      Zahl vorher und nachher steht in §7 (Kommando `grep`, Form der Anker
      aus §1).
- [x] Der Retirement-Check ist je berührter verankerter Regel gelaufen: die
      `state.md` der Beobachtung wurde gelesen, ihr benannter Zielort
      (`AGENTS.md` §5 bzw. §3.8) und der Anker stehen unverändert; §7 nennt das
      Ergebnis je Regel.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag; sie nennt die Größe der
      Datei vorher und nachher (Bytes, mit `wc -c`) **und** was danach noch
      nicht gelöst ist.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — oder
      „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      wellenlos hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `AGENTS.md` §3 (3.3, 3.7, 3.9) | update | Chronik-Sätze und Planungs-Links nach der Form aus §1 entfernen; Bestandsgrenzen und Herkunfts-Anker bleiben wörtlich |
| `AGENTS.md` §5 (Haken-Absatz, Grenzen-Regel, Zähl-Regel, CHANGELOG-Regel, fünf Steering-Loop-Klammern) | update | dasselbe; die doppelt stehende Fenced-Block-Aussage wird zu **einer** |
| `.harness/skills/reviewer.md` (eine `d-check:cite`-Direktive auf `AGENTS.md`-Zeilen) | update | Spiegel, im ersten Wurf des Plans **nicht** gelistet und erst vom `doc-check` gefunden (`citation-mismatch`): jede Kürzung vor der zitierten Stelle verschiebt die Zeilen-Spanne; sie wird auf den neuen Stand nachgezogen, der zitierte Wortlaut bleibt unverändert |

Kein Test, keine Code-Datei. Belege sind `make gates` (doc-check prüft Links
und Anker, `citations`, `gate-consistency`) und der Vorher/Nachher-Vergleich
der Anker-Zahl.

**Ansatz:**

- Vor dem Editieren: `AGENTS.md` §3 und §5 **vollständig** lesen und je Absatz
  das Urteil (Chronik / Zustand / Grenze / Anker) notieren; die Stichwort-Liste
  aus §1 nur als Gegenprobe benutzen.
- Anker-Zahl und Zeichenzahl vor dem ersten Edit messen und in §7 festhalten.
- Edits ohne Shell-Werkzeug (Edit-Tool); kein `sed` über den Baum, denn es
  ist **eine** Datei, und [MR-070](../../../../harness/conventions/MR-070-frozen-klassen-vor-mechanischer-ersetzung.md)
  gilt der Massen-Ersetzung.

## 4. Trigger

**Start** (`next` → `in-progress`): keine Abhängigkeit; der Übergang landet
auf dem Hauptzweig vor der Arbeit. Bei der Beanspruchung entsteht der dritte
Vorprüfungs-Block (Nachtlauf-Stand, `make nightly-state`) in §8.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn die vollständige
  Lektüre zeigt, dass die Chronik-Frage in mehr als einem Dutzend Absätzen
  nicht ohne Regel-Umbau entscheidbar ist, oder die Kürzung Regel-Inhalt
  berührt — dann gehört sie in den Auslager-Vorgang.
- `in-progress` → `open` (blockiert): wenn `make gates` aus einem Grund rot
  bleibt, der nicht an `AGENTS.md` hängt.

## 5. Closure-Trigger

DoD vollständig (§2) + Review-Report liegt vor + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Chronik und Grenze stehen im selben Absatz (§5 Haken-Absatz: die „drei
  Grenzen" folgen dem Chronik-Einstieg); die Kürzung nimmt eine Grenze mit —
  eine Regel, deren Reichweite still schrumpft. **Ausgang:** entfallen — der
  unabhängige Review hat den Diff Hunk für Hunk gegen den Ausgangsstand
  gehalten und keine verlorene Regel oder Grenze gefunden; die drei
  Fenced-Block-Grenzen des Haken-Absatzes stehen unverändert. Die verwandte
  Abweichung in der Gegenrichtung (F-2, eine Reichweite wurde zur Allaussage
  **erweitert**) ist im Review gefunden und behoben.
- Die Stichwort-Suche findet Chronik nur, wo ein Stichwort steht; ein Satz wie
  „Ohne diesen Satz meldet jede Verifikation den Rückstand erneut" fiele
  durch. **Ausgang:** eingetreten — der Dependabot-Punkt in `AGENTS.md` §5 (eine
  Konjunktiv-Erwägung über eine verworfene Alternative) fiel durch das Raster
  und wurde vom unabhängigen Review gefunden (F-1). Behoben im selben Slice
  (Fix-Commit der Review-Runde 1); ein Folge-Slice ist nicht nötig, weil nichts
  offen bleibt — der Ausgang ist damit *geschlossen im Slice*, nicht mit einer
  Folge-Slice-Kennung belegt.
- Nach der Kürzung bleibt `AGENTS.md` der Größenordnung nach groß — die
  Chronik ist nur ein Teil der 38 KB; die Beobachtung „zu groß" bliebe damit
  zu weiten Teilen offen. **Ausgang:** weiter offen — Register-Eintrag
  [`BEO-ALL/briefing-datei-ueberschreitet-lade-budget`](../observations/BEO-ALL/briefing-datei-ueberschreitet-lade-budget/observation.md)
  (1×). Gemessen: 38258 → 36917 Bytes.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Form der Chronik stand **vor** der Messung im
  Plan (§1: Genese, frühere Fassung, Messvorfall als Herleitung, Planungs-Link;
  Gegenprobe über die zwei Tests aus §3.7), und der Diff ließ sich gegen sie
  Hunk für Hunk prüfen — der unabhängige Review konnte jede Entfernung als
  Chronik oder Regelinhalt einordnen und fand **keine** verlorene Regel. Der
  Sensor (`doc-check`, `citations`) meldete den gebrochenen Spiegel im inneren
  Loop, bevor der Feature-Commit entstand.
- **Was ging anders als geplant:** (1) Der Plan listete einen Spiegel nicht: eine
  `d-check:cite`-Direktive im Reviewer-Skill zitiert `AGENTS.md` per
  Zeilen-Spanne, und jede Kürzung davor verschiebt sie. Der Plan wurde **vor**
  der Skill-Änderung ergänzt (§3), nicht nachträglich. (2) Das Stichwort-Raster
  fand nicht alles: der Dependabot-Punkt (Konjunktiv-Erwägung über eine
  verworfene Alternative) und ein Satz samt Link zur abgelösten
  Skript-Mechanik fielen durch — beide vom Review (F-1) bzw. beim Lesen von §5
  gefunden. (3) Ich habe beim Kürzen eine gemessene Reichweite zur Allaussage
  **erweitert** statt sie zu streichen (F-2); der Review fand es.
- **Entfernte Chronik-Sätze** (Fundstelle, Satz-Anfang, Test — je Satz Adressat
  und Zeitform):
  - §3.3 „Historische Klärung: Diese Datei nannte hier früher …" — Vergangenheit
    über die Datei selbst, niemand handelt daraus.
  - §3.3 „— genau das hat die eigene Commit-Historie von …" samt Planungs-Link —
    Genese, Verstoß gegen die Regel „keine Slice-Verweise".
  - §3.3 „ist seither vollständig aufgelöst" → „ist vollständig aufgelöst" —
    Zustand bleibt, die Zeitangabe entfällt.
  - §3.9 „das stand hier zu weit" — Selbstkorrektur.
  - §3.9 „gemessen am Tag-Push von `v0.66.0`, während dieses Gate grün meldete"
    — Messvorfall als Herleitung; die **Grenze** („die Existenz-Prüfung allein
    sieht das nicht") ist herausgelöst und bleibt.
  - §3.9 „das frühere Skript ist darin aufgegangen" und „seit ADR … via Modul" —
    Genese; der Zeiger auf die Modul-ADR bleibt.
  - §4 „seit der Umstellung auf den einen Index" — Genese; „nur in
    `harness/README.md`" bleibt.
  - §5 (Commits) „seit dem Modul `commits` dogfooded" und „die abgelöste
    Skript-Mechanik trug …" samt ADR-Link — Genese (der zweite Satz steht in
    keiner der beiden Plan-Aufzählungen und ist hier einzeln geführt).
  - §5 (Haken-Absatz) „Seit slice-172 hält das ein Sensor" → „Ein Sensor hält
    das" — Genese; ein Fließtext-Satz, **kein** Herkunfts-Anker.
  - §5 (Haken-Absatz) „Der Altbestand bis `slice-170`" samt Planungs-Link →
    „Der Altbestand"; die Abgrenzung trägt
    [MR-056](../../../../harness/conventions/MR-056-dod-haken-waechter.md) (ein
    Hop mehr, benannt in Review-Befund F-3).
  - §5 (Haken-Absatz) der zweite, wortgleich wiederholte Fenced-Block-Satz —
    Duplikat, kein Chronik-Satz.
  - §5 (Grenzen-Regel) „Gemessen an sieben Fundstellen: In sechs …" und „gemessen
    fand ihn in allen sieben Fällen …" — Messvorfall als Herleitung; als
    „in den belegten Fällen" an die Messmenge gebunden (F-2), die Zahlen stehen
    in der Beobachtung.
  - §5 (CHANGELOG) „Gemessen, nicht vereinbart: die Feature-Commits der letzten
    Slices …" und „die Regel darüber sagte nur …" — Genese; der operative
    Kern („ein fehlender Eintrag im Feature-Commit ist kein Rückstand") bleibt.
  - §5 (Dependabot) „Die naheliegende Alternative … hätte … wäre also
    dokumentiert zulässig gewesen …" — Konjunktiv über Verworfenes (F-1); der
    Grund gegen die Alternative steht im Indikativ weiter da.
- **Retirement-Check je berührter verankerter Regel** (`state.md` gelesen,
  Zielort und Anker verglichen): die Regeln zu überzogenen Botschaften (§5,
  Anker `seit welle-82`), zitierter Quelle (§5, `seit slice-147`), Zählmethode
  (§5, `seit slice-210`) und Grenzen-Liste (§5, `seit slice-213`) tragen
  Zielort und Anker **unverändert**; der Text der Grenzen-Regel ist berührt
  (Herleitung gekürzt, Reichweite an die belegten Fälle gebunden), ihr Ableiter
  und ihre Grenzen stehen. Die Regel zur Scan-Achse (§3.8) ist **nicht
  berührt**. Anker vorher/nachher: 5/5 (`grep -oE 'seit (slice|welle)-[0-9]+'`).
- **Größe:** 38258 → 36917 Bytes (`wc -c`), 591 → 576 Zeilen. **Nicht gelöst:**
  die Datei bleibt im Wesentlichen so groß; entfernt sind rund 3,5 %, der
  Rest ist Regeltext (§3 rund 16 KB, §5 rund 12 KB).
- **Steering-Loop-Eintrag:** gezählt, nicht verkörpert. Lernsignale: (a) die
  Spiegel-Suche vor dem Editieren muss **Zitat-Spannen** einschließen — wer vor
  einer zitierten Stelle kürzt, verschiebt sie; der Sensor fängt es, aber erst
  nach dem Edit; (b) ein Stichwort-Raster ist ein Suchraster, kein Umfang, und
  eine Kürzung neigt zum **Umformulieren** statt zum Streichen — beides fand
  der Review, nicht der Lauf.
- **Beobachtungs-Register (`../observations/`):**
  `BEO-ALL/briefing-datei-ueberschreitet-lade-budget/` neu angelegt, Beleg
  `evidence/slice-231.md` (1×); `evidence/slice-231.md` in
  `BEO-ALL/commit-message-overclaims-work/` ergänzt (Klasse (b) an einem
  Regeltext statt an einer Botschaft, benannt).
- **Folge-Slices:** keine. Der Auslager-Vorgang hat keine Kennung; er steht als
  offene Beobachtung im Register.
- **Risiken aus §6:** alle drei mit Ausgang (§6).
- **Drei Paarungen:** Anker — kein `liegt in`-Feld verwendet, kein Gegenstand.
  Folge-Slice — keiner genannt, kein Gegenstand. Register — beide
  genannten Beobachtungs-Verzeichnisse existieren und tragen je einen Beleg
  für diesen Vorgang.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:363-364 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: `AGENTS.md`
fällt unter den Default `*` (Kürzel `ALL`, `harness/conventions.md`
§Modus-Deklaration); sie ist als Sub-Area bereits deklariert, keine
Ausdifferenzierung nötig.

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:369-369 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register unter
`docs/plan/planning/observations/` durchgegangen. Zur Sub-Area `ALL` gibt es
keinen Eintrag, der den Gegenstand — Chronik in einem Zustandsfeld der
Hard-Rule-Datei — trägt. **Nächster Treffer:**
[`kommentar-traegt-herkunfts-prosa-statt-fuenf-klassen`](../observations/BEO-ALL/kommentar-traegt-herkunfts-prosa-statt-fuenf-klassen/observation.md)
(Stand 1×) — dieselbe Klasse Herkunfts-Prosa, aber am **Code-Kommentar**
(§3.7), nicht am Zustandsfeld; kein Beleg für diesen Slice, wird bei der
Closure erneut gegen den Befund gehalten. **Berührt, weil verkörpert:** vier
Regeln in `AGENTS.md` §5 und §3.8 tragen ihren Zielort in einer
`state.md` (`commit-message-overclaims-work`, `citation-stretched-beyond-scope`,
`grenzen-liste-wird-als-vollstaendig-gelesen`,
`zaehlmethode-misst-proxy-statt-gegenstand`, dazu
`module-promise-only-on-scan-axis`); ihr Zielort und Anker bleiben stehen
(DoD, Retirement-Check).

**Modus-Begründungsblock:** Die berührte Sub-Area ist GF — alle berührten
Sub-Areas GF.

### Sub-Area: `*` (`AGENTS.md`)

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — [MR-015](../../../../harness/conventions/MR-015-agents-md-routet.md),
  [MR-045](../../../../harness/conventions/MR-045-slice-verweise-nicht-im-briefing.md)
  und [MR-050](../../../../harness/conventions/MR-050-herkunfts-anker-ist-kein-verweis.md)
  bestimmen, was in der Datei stehen darf.
- **Phase-Reife:** Phase 5 (jeder Lauf lädt sie, seit vielen Wellen stabil).
- **Evidenz-/Diskrepanz-Risiko:** Niedrig bis mittel — die Datei widerspricht
  an zwei Stellen ihrer eigenen Regel (§3.7 Zustandsfelder ohne Chronik;
  keine Planungs-Links, siehe Bezug); der Slice stellt die Konformität her, führt
  keine neue Diskrepanz ein. Das Risiko liegt in der verschränkten Grenze (§6).
- **Reconciliation-Aufwand:** Keiner — kein Brownfield-Bestand.

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung, 2026-09-27):
`make nightly-state` meldet `upstream-drift.yml` **ROT** (Lauf
2026-09-26T05:44Z) und `image-scan.yml` **grün**. Die Ausgabe der vier roten
Achsen, lokal nachgefahren: `golangci-lint` VERALTET (Pin 2.13.2, upstream
2.14.0), `semgrep` VERALTET (Pin 1.177.0, upstream 1.178.0), `a-check`
VERALTET (Pin 0.19.0, upstream 0.20.0), `golang:1.27.1` ABWEICHEND
(Digest unter demselben Tag neu gebaut). Das sind **planmäßige**
Fremd-Release-Meldungen, keine unerwarteten; sie berühren `AGENTS.md` nicht und
sind nicht Gegenstand dieses Slice.
