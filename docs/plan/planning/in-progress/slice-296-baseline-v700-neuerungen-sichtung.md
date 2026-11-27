# Slice slice-296: Neuerungen von `v7.0.0` sichten — eine Entscheidungsvorlage für den Auftraggeber, keine Übernahme

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`wave-98-baseline-v700`](../wave-98-baseline-v700.md).

**Bezug:** [`MR-011`](../../../../harness/conventions.md#mr-011) (die Pin-Kette),
[`MR-080`](../../../../harness/conventions.md#mr-080) (die vorige Hebung; ihr Delta stand „vollständig nach Datei,
keine davon übernommen“ — dieselbe Abgrenzung), [`slice-295`](../done/slice-295-baseline-v700-bump.md) (liefert den
vendorten Baum und den Vor-Bump-Commit, aus dem der alte Baum gelesen wird),
[`ADR-0113`](../../adr/0113-pinsync-modul.md) (der Gegenstand, an dem die Frage zur Umgebungs-Zusage sich misst).

**Berührte Spec-Stellen:** — (der Slice liest und schreibt eine Vorlage; er berührt keine Spec-Stelle).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-11-25.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Zu jeder Bewegung zwischen `v6.18.1` und `v7.0.0` — 39 Dateien und die Sammelzeile der 7 reinen
Stempel-Änderungen — steht in **einer** Entscheidungsvorlage eine Zeile mit Klasse (*bereits befolgt* · *neu und
anwendbar* · *nicht anwendbar*), Fundstelle im vendorten Baum, Empfehlung und Preis; dazu die zwei benannten
Fragen (Umgebungs-Zusage und `environment-check` gegen die vorhandenen Gates; die **Reste** der Namensregel im
Angebot des Produkts an Konsumenten). Der
Auftraggeber entscheidet je Zeile; **der Slice übernimmt nichts**.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Keine Übernahme** (Entscheid 2). Kein Code, keine Spec-Änderung, kein Vorlagen-Wechsel, keine
  Konfigurations-Änderung, kein neuer MR-Eintrag, kein Folge-Slice. Folge-Slices entstehen **nach** der Antwort des
  Auftraggebers, je Übernahme einer, und keiner wird hier vorab angelegt: Ein Slice ohne Entscheid wäre ein Plan
  für etwas, das niemand verlangt hat
  ([`BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass`](../observations/BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass/observation.md)).
- **Kein Durchgang durch die Adaptions-Liste** (Freshness-Audit) — Frage 10 in [wave-98 §5](../wave-98-baseline-v700.md).
  Die Vorlage nennt zu einer Bewegung die aktive Adaption, die sie berührt, **wenn** die Lektüre sie findet; sie
  prüft die 41 Einträge nicht je einzeln.
- **Kein Wort über die Namensregel selbst hinaus.** Entscheid 1 ist gefallen und in [`slice-297`](../done/slice-297-spezifikation-heisst-specification.md) und [`slice-298`](../done/slice-298-neue-namen-wave-archive-anker.md) umgesetzt; die Vorlage führt die
  Zeilen dazu mit „entschieden“ und fragt nicht erneut — gefragt wird nur nach den **Resten** (Frage B).
- **Kein CR an den Kurs.** Findet die Sichtung einen Widerspruch im Regelwerk, steht er in der Vorlage; ein
  ausgehender Change Request ([`MR-035`](../../../../harness/conventions.md#mr-035)) ist ein eigener Vorgang.
- **Keine Änderung am Bestand an Plänen und Slices**, auch keine Berichtigung von `welle-`-Kennungen: sie bleiben (Bestand;
  Entscheid 1, „Bestand bleibt“).

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

<!-- BEDIENHINWEIS: je Zeile ein pruefbares Kriterium. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Inventar:** Jede der 39 Dateien mit Bewegung (31 inhaltlich verschieden, 7 umbenannt, 1 neu) hat eine Zeile
      in der Vorlage, **größte Bewegung zuerst**; die 7 Dateien, die nur ihren Stempel ändern, stehen in **einer**
      Sammelzeile mit dem Beleg „je 2 rohe Diff-Zeilen, die `<!-- Quelle:`-Zeile“. Eine Datei mit mehreren
      Neuerungen (erwartet: `modul-14`, `modul-04`, `modul-05`) trägt je Neuerung eine Unterzeile. Die Zählung ist am
      **vendorten** Baum von `v7.0.0` und am alten Baum aus `git archive <Vor-Bump>` wiederholt (§3); weicht sie von §1
      der Welle ab, steht die Abweichung in der Vorlage, nicht stillschweigend in der Zahl.
- [ ] **Klasse und Empfehlung je Zeile:** *bereits befolgt* nennt den Beleg im Repo (Pfad oder Target — der
      Reviewer öffnet jeden), *neu und anwendbar* nennt die **kleinste** Antwort mit Preis (Dateien, Zeilen, Gate) und,
      wo es eine gibt, die größere daneben, *nicht anwendbar* nennt den Grund **an der Bedingung der Regel**
      ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 13). Die Fundstelle ist ein Anker-Link in den vendorten Baum;
      ein Wortlaut, der die Klasse trägt, ist zitiert und mit `d-check:cite` gehalten, nicht nacherzählt
      ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 18). Jede Zeile hat eine leere Spalte *Entscheidung des
      Auftraggebers*.
- [ ] **Die zwei benannten Fragen:** (A) Umgebungs-Zusage und `environment-check` — je Bestandteil des Kapitels
      gegen das, was das Repo hat (§3), mit einer Antwort *befolgt · Lücke · nicht anwendbar*; (B) die **Reste der
      Namensregel**, die in [`slice-297`](../done/slice-297-spezifikation-heisst-specification.md) und
      [`slice-298`](../done/slice-298-neue-namen-wave-archive-anker.md) nicht stehen — die Defaults des Werkzeugs
      (`planning.waves.glob` und `results-glob` heißen `welle-*`, die Vorlage `.d-check.yml` des Regelwerks nennt `wave-*`),
      die Beispiel-Pfade, die das Werkzeug und das Handbuch **Konsumenten** zeigen (`spec/spezifikation.md` in
      `--print-config`, `--suggest-config` und im Handbuch) und das Verhalten für einen Konsumenten mit den
      neuen Namen — je Rest die **gezählte** Reichweite und Preis der Übernahme gegen Preis des Bleibens.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `baseline-v700-neuerungen.md` unter `docs/plan/planning/` | neu | **die** Entscheidungsvorlage; der Ort ist die Antwort auf Frage 11 in [wave-98 §5](../wave-98-baseline-v700.md) (Vorschlag: flach in `planning/`, der Name trifft weder den Glob `welle-*.md` noch (nach Frage 3) `wave-*.md` des Moduls `planning`). Aus nichts kopiert — die Klasse hat keine Vorlage im Baum; die Spalten stehen in §3 unten |
| `docs/plan/planning/observations/` | update | nur, falls die Sichtung ein Auftreten belegt (Evidenz-Datei), sonst „keine Beobachtung angefallen“ in §7 |

Keine Testdatei: Der Slice schreibt Markdown, keinen Code; die Sensoren sind `make doc-check` (Links, Anker,
`citations`) und `make gates`.

**Spalten der Vorlage** (je Zeile; die Reihenfolge der Zeilen ist die der Diff-Größe, größte zuerst):
laufende Nummer · Datei und Neuerung · Fundstelle (Anker-Link in `.harness/baseline/v7.0.0/…`, Wortlaut zitiert, wo er die
Klasse trägt) · Klasse · Beleg im Repo **oder** kleinste Antwort mit Preis **oder** Grund der Nicht-Anwendbarkeit ·
Empfehlung · *Entscheidung des Auftraggebers* (leer). Die laufende Nummer in Spalte 1 macht die Zeilen zählbar:
`grep -c '^| [0-9]' docs/plan/planning/baseline-v700-neuerungen.md`.

**Vorgehen:**

1. **Zählung wiederholen**, am vendorten Baum von `v7.0.0` und am alten Baum aus `git archive <Vor-Bump>
   .harness/baseline/v6.18.1 | tar -x -C <Wegwerf-Verzeichnis>`: Dateien alt/neu (54/55), gleichnamig (47), byte-verschieden
   (38), inhaltlich verschieden nach Abzug der `<!-- Quelle:`-Zeilen und mit normalisiertem Versionsstempel (31), Diff-Zeilen
   roh (862) und normalisiert (804). Die Zahlen aus §1 der Welle sind Erwartung, kein Beleg.
2. **Den Diff lesen, nicht die Zusammenfassung** — nicht die des Kurses, nicht die dieses Plans. Die größten Bewegungen
   zuerst: `modul-14` (258 Diff-Zeilen), `templates/.d-check.yml` (74), `modul-04` (57), `modul-05` (49),
   `templates/README.md` (35), danach der Rest nach Größe.
3. **Klassifizieren**, mit dem Beleg, den die Klasse verlangt (§2). Eine Bewegung, die ein Entscheid des Auftraggebers
   oder eine aktive Adaption schon beantwortet, steht in *nicht anwendbar* mit dem Grund „entschieden“ und der
   Fundstelle der Entscheidung — das sind die Zeilen zur Namensregel (Entscheid 1: der Absatz *Dateinamen sind
   englisch* in `grundlagen-harness-dateien.md` und die sieben Vorlagen-Namen), deren Umsetzung
   [`slice-297`](../done/slice-297-spezifikation-heisst-specification.md) und
   [`slice-298`](../done/slice-298-neue-namen-wave-archive-anker.md) tragen.
4. **Frage (A), Umgebungs-Zusage.** Je Bestandteil des neuen Kapitels in `modul-14` und je Abschnitt der Vorlage
   `environment-check.template.md` eine Zeile: *die Zusage selbst* (festgelegt **und** adressierbar; als `LH-QA-*`
   geführt — trägt das Repo eine solche Anforderung? `git grep -n "Umgebung" -- spec/lastenheft.md` ist der Anfang der
   Suche, nicht ihr Ergebnis), *der Docker-Träger* (Digest-Pins im `Dockerfile`, `make versions`, der Image-Hash am
   Ende von `make fullbuild`), *der zweite Träger* (Tool-Pin und Tool-/Versions-Check), *die zwei Prüfungen derselben
   Zusage* gegen [`pinsync`](../../adr/0113-pinsync-modul.md) (Wertgleichheit an mehreren Stellen) und gegen
   `make freshness-*`/`make *-digest` (Aktualität, nicht Gleichheit), *Inventar, Stufen, Kontext*, *die Probe im
   leeren Kontext* gegen `make guard-probe`. Je Zeile: *befolgt* (Beleg), *Lücke* (kleinste Antwort mit Preis) oder
   *nicht anwendbar* (Grund an der Bedingung der Regel, am **Wortlaut** des Kapitels belegt und nicht aus diesem
   Plan übernommen: ob der zweite Träger neben dem Docker-Träger gilt, sagt das Kapitel, nicht der Planner).
5. **Frage (B), Reste der Namensregel.** Je Rest eine **gezählte** Reichweite, jede mit der Form ihrer Zählung
   ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 14): (i) die **Defaults** des Werkzeugs für `planning.waves`
   (`welle-*.md`, `welle-*-results.md`: `internal/hexagon/core/model/config.go` Z. 481 und 490; die Schema-Zeilen der
   Spezifikation; das Handbuch) gegen die Vorlage `.d-check.yml` des Regelwerks, die `wave-*` verlangt — ein Konsument,
   der dem Regelwerk folgt, muss heute beide Globs setzen; (ii) die **Beispiel-Pfade** für Konsumenten
   (`config_template.go` Z. 56–57, `suggest.go` Z. 379–380, `benutzerhandbuch.md` Z. 372–373: `spec/spezifikation.md`);
   (iii) was ein Konsument sieht, dessen Wellen `wave-<NN>` heißen und dessen Archiv `archive.zip` heißt (aus der Probe
   in [wave-98 §1](../wave-98-baseline-v700.md), Punkt 5 und 6, übernommen als Erwartung, am Image des Hauptzweigs
   neu gemessen). Je Rest *Preis der Übernahme* gegen *Preis des Bleibens* (ein Default ist ein Vertrag mit
   Konsumenten: [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in))
   — **beide**, nicht nur der eine, auf den die Empfehlung fällt.
6. **Die Vorlage schreiben** und dem Reviewer übergeben. Die Spalte *Entscheidung des Auftraggebers* bleibt leer;
   der Slice schließt, wenn sie dem Auftraggeber **vorliegt**, nicht wenn er geantwortet hat.

**Plan-Änderung (Planner, vor dem Übergang `open` → `next`; Stand `2831e071`).** Drei Stellen des Plans stimmen nach
[`slice-297`](../done/slice-297-spezifikation-heisst-specification.md) und
[`slice-298`](../done/slice-298-neue-namen-wave-archive-anker.md) nicht mehr; diese Zeilen überschreiben sie.

1. **Frage (B) (ii) nennt drei der vier Spiegel-Stellen.** Das Konsumenten-Beispiel mit dem alten deutschen Dateinamen der
   Spezifikation steht auch in `spec/specification.md` Z. 233–234, und `--manual` zeigt es Konsumenten. Es sind **acht
   Zeilen an vier Stellen** (Spezifikation, Handbuch, `config_template.go`, `suggest.go`):
   [Review zu slice-297](../../../reviews/2026-11-26-slice-297-spezifikation-heisst-specification-r1.md), Befund F-6.
   F-6 gehört in die Frage (B), samt seiner zweiten Hälfte: `--suggest-config` schlägt einem Repo mit
   Spezifikations-Verzeichnis den deutschen Namen vor, den dieses Repo nicht mehr führt. Für den Preis gilt
   [`MR-025`](../../../../harness/conventions.md#mr-025) (die Spiegel ändern sich an allen Stellen oder an keiner).
2. **Frage (B) (iii), die Erwartung aus wave-98 §1 Punkt 6 ist überholt.** Seit
   [`slice-298`](../done/slice-298-neue-namen-wave-archive-anker.md) liest `tools/archive-wave` `wave-<NN>` und schreibt
   `archive.zip`; das Werkzeug ist zudem kein Bestandteil des Images und am Image nicht zu messen. Am Image neu zu messen ist
   allein Punkt 5 (das Modul `planning` mit den Defaults).
3. **§1 nennt 41 Adaptionen**; es sind 42 aktive (`ls harness/conventions/*.md | wc -l`; am Stand `590dde3d` 41). Die
   Abgrenzung bleibt: kein Durchgang je Eintrag.

**Plan-Änderung (Implementer, vor dem Commit der Vorlage; Stand `5a14fe90`).**

1. **Wortlaute stehen nicht in der Fundstellen-Spalte.** Das Modul `citations` sucht das Zitat ab der Zeile **nach**
   der Direktive (`citationQuote` in `internal/hexagon/core/rules/citations.go`); in einer Tabellenzelle ist das die
   nächste Zeile der Tabelle. Die Zitate samt `d-check:cite` stehen deshalb in einem eigenen Abschnitt *Wortlaute*
   der Vorlage, die Zeilen zeigen mit ihrer Nummer dorthin.
2. **Zählkommando.** Die Unterzeilen tragen dieselbe Dateinummer mit Buchstaben; `grep -c '^| [0-9]'` zählt die
   Zeilen, die Dateien samt Sammelzeile zählt `grep -o '^| [0-9][0-9]*' <Vorlage> | sort -u | wc -l`.
3. **Zwei Erwartungen aus Schritt 5 treffen nicht:** Das Handbuch nennt die Default-Globs nicht, und
   [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
   legt den Default nicht fest (die Schema-Tabelle der Spezifikation tut es). Gemessen steht es in der Vorlage.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`slice-295`](../done/slice-295-baseline-v700-bump.md), [`slice-297`](../done/slice-297-spezifikation-heisst-specification.md) und [`slice-298`](../done/slice-298-neue-namen-wave-archive-anker.md) liegen in `done/` (der Baum von
`v7.0.0` ist vendort, `ls .harness/baseline` nennt nur ihn, und der Bump-Eintrag nennt den Vor-Bump-Commit);
`in-progress/` trägt keinen Slice; die Antwort auf Frage 11 (Ort der Vorlage) liegt im Architect-Verdikt vor;
`Verantwortlich:` beim Übergang `open` → `next` gesetzt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Vorlage trägt mehr Zeilen, als **eine** Review-Sitzung liest
  (Maß: der Reviewer öffnet jede Fundstelle und jeden Beleg, §2). Der Schnitt ist vorab benannt, nicht erst dann
  gesucht: **Regelwerk** (19 inhaltlich verschiedene Dateien) und **Vorlagen** (12 verschieden, 7 umbenannt, 1 neu) —
  zwei Slices, die Vorlage dann in zwei Dateien; dann gilt die Form von
  [`MR-066`](../../../../harness/conventions.md#mr-066).
- `in-progress` → `open` (blockiert): Der vendorte Baum von `v7.0.0` fehlt oder der Vor-Bump-Commit ist nicht
  auflösbar (`git cat-file -e <Vor-Bump>`), oder das Architect-Verdikt zu Frage 11 steht aus.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (§2) und `make gates` grün; Review- und Verifikations-Report liegen unter `docs/reviews/`;
Closure-Notiz mit Lerneintrag geschrieben. **Review-Vorgabe der Welle:** höchstens zwei Runden, eingearbeitet
werden nur HIGH und MEDIUM; jede Nachbesserung zitiert die Stelle, die sie trägt
([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 18). **Die Antworten des Auftraggebers sind keine Closure-Bedingung**;
dass die Vorlage ihm vorliegt, hält der Closure-Trigger der [Welle](../wave-98-baseline-v700.md), nicht diese DoD
([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 10).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Eine Zeile „bereits befolgt“ ist eine Behauptung ohne Beleg**
  ([`BEO-ALL/eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md),
  18×;
  [`BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md),
  10×). — **Ausgang:** *entfallen*, wenn der Reviewer **jeden** genannten Beleg öffnet und jeder trägt; sonst
  *eingetreten* mit einer Evidenz-Datei.
- **Die Empfehlung weitet die Anforderung über den Anlass hinaus** — eine Zeile bündelt mehrere Neuerungen oder nennt nur
  die große Antwort
  ([`BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass`](../observations/BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass/observation.md),
  1×, offen). — **Ausgang:** *entfallen*, wenn jede Empfehlung *übernehmen* die kleinste Antwort mit Preis nennt und
  keine Zeile mehr als eine Neuerung trägt; sonst *eingetreten* (dann 2×) mit einer Evidenz-Datei.
- **Eine Grenze wird ohne ihre Bedingung wiedergegeben** („gilt nicht für …“ statt „gilt nicht, wenn …“)
  ([`BEO-ALL/grenz-aussage-ohne-ihre-bedingung`](../observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/observation.md),
  verkörpert in Regel 13). — **Ausgang:** *entfallen*, wenn jede *nicht anwendbar*-Zeile die Bedingung des Kapitels
  wörtlich zitiert; sonst *eingetreten*.
- **Die Sichtung schleust eine Übernahme ein** — eine Änderung außerhalb der Vorlage. — **Ausgang:** *entfallen*, wenn
  `git diff --stat <Start>..HEAD` außer der Vorlage, den Plan-Dateien, `docs/reviews/` und dem Register nichts
  zeigt; sonst *eingetreten* und die Änderung geht zurück (Plan-Änderung vor dem Code,
  [`AGENTS.md`](../../../../AGENTS.md) §6 Schritt 4).
- **Ein Wortlaut wird gestreckt** — die Fundstelle trägt weniger, als die Zeile ihr zuschreibt
  ([`BEO-ALL/citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md),
  20×). — **Ausgang:** *entfallen*, wenn Reviewer und Verifier an jeder zitierten Stelle das Feld lesen, nicht den
  Titel ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 16), und keine abweicht; sonst *eingetreten*.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md:416-417 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Die Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md) führt zwei Zeilen: den Default `*` (`ALL`) und
`tools/harness/` (`HARN`), beide Greenfield. Berührt ist der Default `*` — `docs/plan/planning/` (die Vorlage) und
lesend `.harness/baseline/` und das Repo, aus dem die Beläge der Zeilen stammen. `tools/harness/` bleibt
unberührt; wird eine Zeile der Vorlage zu einer Übernahme in `tools/harness/`, ist das ein Folge-Slice mit eigener
Sub-Area-Prüfung.

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md:422-422 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Stand 2026-11-25,
`ls docs/plan/planning/observations/BEO-ALL/`; Zähler = Dateien unter `evidence/`), Sub-Area `*`. Treffer:
[`BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass`](../observations/BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass/observation.md)
(offen, 1× — die Empfehlung je Zeile ist die kleinste Antwort; Risiko 2),
[`BEO-ALL/eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
(gemischt, 18× — Risiko 1: „bereits befolgt“ braucht einen geöffneten Beleg),
[`BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md)
(geplant, 10× — derselbe Befund von der Wortlaut-Seite),
[`BEO-ALL/citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
(verkörpert, 20× — Risiko 5),
[`BEO-ALL/grenz-aussage-ohne-ihre-bedingung`](../observations/BEO-ALL/grenz-aussage-ohne-ihre-bedingung/observation.md)
(verkörpert in Regel 13, 3× — Risiko 3),
[`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md)
(verkörpert, 19× — das Inventar der 39 Dateien wird gezählt, nicht aus dieser Liste übernommen),
[`BEO-ALL/rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md)
(gemischt, 9× — die Sichtung geht den **Bestand** der 39 Dateien durch, nicht nur die zwei benannten Anlässe
`environment-check` und die Reste der Namensregel). **Kein Treffer:** keine Beobachtung zu Entscheidungsvorlagen als Klasse.
Der Zähler folgt aus den Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-11-25: `image-scan` **grün** (jüngster Lauf 2026-11-25T12:59:28Z). `upstream-drift`
**rot** (jüngster Lauf 2026-11-25T09:26:13Z), der Log war über `gh run view` nicht abrufbar; gelesen ist
`make baseline-freshness` (Exit 3, neuer Release `v7.0.0`) — der Gegenstand der Welle, **planmäßig gelesen**. Zum
Zeitpunkt dieses Slice sind [`slice-295`](../done/slice-295-baseline-v700-bump.md), [`slice-297`](../done/slice-297-spezifikation-heisst-specification.md) und [`slice-298`](../done/slice-298-neue-namen-wave-archive-anker.md) `done/`; der Stand wird beim Start neu gelesen.

**Beim Start neu gelesen** (Planner, 2026-11-26). **Register:** seit dem Stand oben ist
[`BEO-ALL/praefix-umstellung-schaltet-waechter-der-anderen-schreibweise-ab`](../observations/BEO-ALL/praefix-umstellung-schaltet-waechter-der-anderen-schreibweise-ab/observation.md)
dazugekommen (offen, 1×); er trifft Frage (B): ein Konsument mit beiden Schreibweisen kann nur eine bewachen.
**Nachtlauf-Stand:** `make nightly-state`: `image-scan` grün (jüngster Lauf 2026-11-25T12:59:28Z), `upstream-drift` rot
(jüngster Lauf 2026-11-26T09:05:03Z) — gelesen: der Lauf fuhr `origin/main` (`76cf7d51`, Pin vor der Hebung) und meldet
in `make baseline-freshness` den neuen Release (`gh run view --log-failed`); lokal endet dasselbe Target mit Exit 0
(Pin `v7.0.0` ist der neueste, Inhalt unverändert). **Planmäßig**, bis die Hebung gepusht ist.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
