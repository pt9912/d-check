# Slice slice-224: Baseline-Pin auf `v6.9.0` — mechanisch, ohne Urteil über den Delta

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`MR-011`](../../../../harness/conventions.md#mr-011)-Pin-Serie
(Vorgänger-Eintrag: [`MR-071`](../../../../harness/conventions.md#mr-071)),
[`MR-021`](../../../../harness/conventions.md#mr-021) (pin-gebundene
Verweise), [`MR-051`](../../../../harness/conventions.md#mr-051)
(`d-check:cite`-Spannen neu ankern),
[`MR-055`](../../../../harness/conventions.md#mr-055) (Symlink als Träger),
[`MR-069`](../../../../harness/conventions.md#mr-069) (`ignore-refs` als
deklarierte Gate-Senkung),
[`MR-070`](../../../../harness/conventions.md#mr-070) (Frozen-Klassen vor
mechanischer Ersetzung). Der neue Eintrag der Serie ist [`MR-072`](../../../../harness/conventions.md#mr-072).

**Berührte Spec-Stellen:** — (der Slice berührt keine Spec-Stelle).

**Verantwortlich:** pt9912 (Implementer-Rolle).

**Autor:** pt9912. **Datum:** 2026-09-16.

---

## 1. Ziel und Abgrenzung

**Ziel.** Den vendorten Baseline-Bestand von `v6.6.0` auf `v6.9.0` heben und
alle pin-gebundenen Verweise nachziehen — **mechanisch, ohne den Regel-Delta
zu beurteilen**. Was der Delta inhaltlich verlangt, entscheidet ein
Folge-Slice; dieselbe Zerlegung wie beim Vorgänger-Paar
[slice-207](../done/slice-207-baseline-v650-bump.md) /
[slice-208](../done/slice-208-v650-regel-adoption.md).

**Der Anlass ist am Sensor bestätigt, nicht übernommen:**
`make baseline-freshness` meldet *„NEUER RELEASE verfügbar (Pin v6.6.0):
v6.7.0, v6.7.1, v6.7.2, v6.8.0, v6.9.0"* und zugleich, dass der **gepinnte**
Tag upstream inhaltlich unverändert ist (Bytes == vendored `SHA256SUMS`).
`make nightly-state` bestätigt denselben Befund am jüngsten Lauf von
`upstream-drift.yml` (2026-09-16, planmäßig rot — Fremd-Release, keine
Störung).

**Fünf Releases liegen zwischen dem Pin und upstream, nicht eines.** Das ist
kein neuer Fall: Die [`MR-011`](../../../../harness/conventions.md#mr-011)-Kette hat wiederholt mehrere Zwischen-Tags in
einem Sprung übersprungen (`v5.7.0`→`v5.9.0`, `v5.9.0`→`v5.11.0`,
`v5.15.0`→`v5.18.0`, `v6.0.0`→`v6.3.1`) — ein Bump zielt auf den **jeweils
aktuellen** Tag, nicht auf jeden dazwischen. Dasselbe gilt hier: Ziel ist
`v6.9.0`, kein Zwischenstopp.

**Abgrenzung — vier Punkte, jeder mit Grund:**

1. **Kein Urteil über den Regel-Delta.** Was `v6.9.0` gegenüber `v6.6.0` an
   Regeln ändert, wird **gemessen und gelistet**, aber nicht beantwortet —
   *ein Folge-Slice übernimmt es*, mit einer Antwort je Regel (übernommen ·
   nicht anwendbar mit Begründung · abweichend als Adaption).
2. **Keine Template-Adoption.** `AGENTS.md` §4 trägt bereits die Ziel-Form aus
   der `v6.6.0`-Vorlage (§4 als Zeiger auf `harness/README.md`, slice-222).
   Ändert `v6.9.0`s `AGENTS.template.md` diese Form **erneut**, ist das ein
   Fund dieses Slice (§3, Delta-Messung), aber seine Übernahme bleibt dem
   Folge-Slice vorbehalten — *es wäre ein anderer Vorgang*, und die Lehre aus
   slice-222 (§6 dieses Plans) ist gerade, diesen Punkt **vorab** zu
   entscheiden statt ihn während der Implementierung neu aufzurollen. Trifft
   der Fall doch ein, ist die in slice-222 gewählte Antwort (Bump + Adoption
   im selben Slice, per ausdrücklicher Auftraggeber-Weisung) die Präzedenz,
   an der zu messen ist — keine stillschweigende Wiederholung ohne dieselbe
   Weisung.
3. **slice-221 wird nicht angefasst.** Weder beansprucht noch nachgezogen —
   *Schicht-Abgrenzung*: Dieser Slice hebt einen Pin, er räumt keine Tabelle
   auf. Es liegt weiter in `next/` und wird **nach** diesem Bump gegen die
   dann geltende Vorlage neu zugeschnitten (sein eigener Trigger, §4 dort).
4. **Kein Aufräumen der `ignore-refs`-Einträge.** Wächst die Liste durch
   diesen Bump (neuer Eintrag, weil der `v6.6.0`-Baum verschwindet und
   eingefrorene Artefakte ihn weiter zitieren), ist das als deklarierte
   Gate-Senkung bereits geführt ([`MR-069`](../../../../harness/conventions.md#mr-069))
   — *Bestand bleibt bewusst stehen*.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [x] **(1)** `.harness/baseline/v6.9.0/` ist materialisiert (`regelwerk/`,
      `templates/`, `SHA256SUMS`), der `v6.6.0`-Baum entfernt, §Baseline in
      [`harness/conventions.md`](../../../../harness/conventions.md) zeigt auf
      den neuen Tag, und [`MR-072`](../../../../harness/conventions.md#mr-072) trägt die Hebung als nächster Eintrag der
      [`MR-011`](../../../../harness/conventions.md#mr-011)-Serie.
      `make baseline-verify` grün.
- [x] **(2)** **Alle vier Spiegel-Klassen sind nachgezogen, nicht nur die
      grep-bare** —
      [`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
      (7×, weiterhin ohne formgültigen Ausgang) nennt sie: Pfad-Verweise
      (gate-gedeckt) · Release-/Tree-**URLs** mit dem Tag ·
      **Prosa-/Ellipsen-Pins** · der **zitierende Verweis**, dessen Wortlaut
      am neuen Ziel nicht mehr stehen muss
      ([`MR-051`](../../../../harness/conventions.md#mr-051)). Je Klasse steht
      im Slice, **wie** sie gesucht wurde — Suchform **nach der Version**, mit
      Gruppierung nach Präfix (Lehre aus slice-222, nicht nach dem Pfad
      suchen).
- [x] **(3)** **Die Frozen-Klassen sind VOR der mechanischen Ersetzung
      aufgelistet** ([`MR-070`](../../../../harness/conventions.md#mr-070),
      Geltungsbereich trifft hier zu: eine mechanische Ersetzung über mehr als
      eine Datei) — über die **Eigenschaft**, nicht über Verzeichnisse. Was
      eingefroren ist, wird **nicht** retargetet, sondern über `ignore-refs`
      abgefangen. **Lebende** Artefakte mit einer Vergangenheits-Aussage über
      einen früheren Pin (wie [`MR-067`](../../../../harness/conventions.md#mr-067)
      in slice-222) fallen nicht automatisch unter „Bestand" — jeder Eintrag
      der [`MR-011`](../../../../harness/conventions.md#mr-011)-Kette wird einzeln gegen seinen `Geltungsbereich` gelesen,
      nicht pauschal als „liegt in `conventions/`, also lebend, also
      retargeten".
- [x] `make gates` grün.
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

**Die Reihenfolge ist nicht beliebig, und der Grund steht in
[`MR-070`](../../../../harness/conventions.md#mr-070):** Die Frozen-Liste
entsteht **vor** der ersten Ersetzung, sonst ist sie eine Rechtfertigung
hinterher.

1. **Frozen-Klassen auflisten** — über die Eigenschaft *„würde ein
   korrigierter Wert verfälschen, was dieses Artefakt zu seinem Datum
   festgehalten hat?"*, nicht über Verzeichnisnamen. Jede [`MR-011`](../../../../harness/conventions.md#mr-011)-Kettenglied
   einzeln gegen seinen `Geltungsbereich` lesen (Lehre aus slice-222s
   Über-Hebung von [`MR-067`](../../../../harness/conventions.md#mr-067)).
2. **Neuen Baum materialisieren** (`fetch-baseline-cache.sh v6.9.0`), **ohne**
   den alten zu löschen — beide koexistieren kurz, damit Schritt 3 diffen
   kann.
3. **Delta messen, nicht beurteilen:** `diff -I '<!-- Quelle:'` gegen den
   alten Baum (die Herkunftszeile in Zeile 3 jeder Regelwerk-Datei trägt den
   Tag und meldete sonst *jede* Datei als geändert). Die Liste wandert in den
   Slice und ist die Eingabe des Folge-Slice. **Hier fällt die Entscheidung**,
   ob Abgrenzung 2 (keine Template-Adoption) hält oder — wie in slice-222 —
   eine erzwungene Ausnahme braucht; ohne Weisung des Auftraggebers bleibt es
   bei Abgrenzung 2.
4. **Alten Baum entfernen**, `make baseline-verify`.
5. **Vier Spiegel-Klassen nachziehen**, je mit benannter Suchform — Version
   statt Pfad, Gruppierung nach Präfix.
6. **`d-check:cite`-Spannen neu ankern** — der Bump verschiebt Zeilennummern,
   und `citations` ist fail-closed im inneren Loop
   ([`MR-051`](../../../../harness/conventions.md#mr-051)).
7. [`MR-072`](../../../../harness/conventions.md#mr-072) schreiben, §Baseline umstellen, `make gates`, Handoff.

## 4. Trigger

**Start** (`next` → `in-progress`): WIP-Limit frei (`in-progress/` ist leer,
Ruhe-Marker steht), `make baseline-freshness` meldet den neuen Release, und
der Auftraggeber hat den Vorgang am 2026-09-16 beauftragt (Entscheidung
zwischen drei Optionen — voller Bump zuerst statt Weiterarbeit an slice-221
oder nur Sichtung des Deltas).

**Rückführung nach `next/`** (`in-progress→next`): wenn der gemessene Delta
so groß ist, dass das **Nachziehen der Spiegel** selbst mehrere Sitzungen
braucht — dann trennt sich der Bump in Vendoring und Retargeting.

**Rückführung nach `open/`** (`in-progress→open`): wenn `v6.9.0` eine
**Struktur**-Änderung am vendorten Baum mitbringt (andere Verzeichnisnamen,
anderes Bundle-Layout), die die Pfad-Verweise nicht mechanisch abbildbar
macht. Dann ist vorher eine Entscheidung über die Verweis-Form fällig, wie sie
[`MR-023`](../../../../harness/conventions.md#mr-023) beim
Bundle-Layout-Wechsel gebraucht hat.

## 5. Closure-Trigger

DoD (1) bis (3) abgehakt, `make gates` und `make baseline-verify` grün mit
echter Ausgabe, unabhängiger Review durchgeführt und eingearbeitet,
Closure-Notiz geschrieben, Register fortgeschrieben, jedes Risiko aus §6 mit
einem der drei Ausgänge.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Drei der vier Spiegel-Klassen deckt kein Gate** — dieselbe Lücke wie bei
  jedem Vorgänger-Bump, zuletzt gemessen bei 7×
  ([`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md),
  weiterhin ohne formgültigen Ausgang: die mechanische Form `versions.patterns`
  existiert seit slice-122, ist aber bewusst nicht scharfgeschaltet). Ein
  grüner `make gates`-Lauf nach dem Bump sagt über Release-URLs, Prosa-Pins
  und zitierende Verweise **nichts**. — **Ausgang:** eingetreten, wie
  erwartet — und diesmal mit einem neuen Datenpunkt: kein neuer
  `ignore-refs`-Eintrag nötig, weil alle neun eingefrorenen Fundstellen (eine davon erst im Review gefunden, siehe unten) den
  entfernten Baum nur in Inline-Code/Prosa tragen, nicht als Markdown-Link
  (Beleg: `evidence/slice-224.md` bei `pin-bump-mirrors-ungated`).
- **Die Über-Hebungs-Falle bei lebenden `MR`-Einträgen ist bekannt, aber die
  Frozen-Liste fängt sie nicht.** slice-222 hat gezeigt: Eine mechanische
  Ersetzung über den `Geltungsbereich`-Text kann einen Eintrag treffen, der
  *„lebend"* ist (liegt in `harness/conventions/`), aber inhaltlich eine
  **Vergangenheits-Aussage** über einen früheren Pin trägt
  ([`mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md),
  bereits als [`MR-070`](../../../../harness/conventions.md#mr-070)
  verkörpert — die Grenze der Regel ist benannt, nicht behoben). Gegenmittel
  dieses Slice: jedes [`MR-011`](../../../../harness/conventions.md#mr-011)-Kettenglied einzeln lesen statt pauschal
  „`conventions/` ⇒ retargeten". — **Ausgang:** eingetreten, trotz
  korrekt vorab erstellter Frozen-Liste — nicht an der gelisteten Datei
  selbst (die stand korrekt darauf), sondern an der
  **Tabellenzeile**, mit der `harness/conventions.md` denselben Eintrag im
  Adaptions-Index führt. Ein pauschaler `sed` über diese lebende Datei traf
  die Zeile mit; erkannt und behoben im selben Arbeitsschritt, vor dem
  nächsten Schritt. Beleg: `evidence/slice-224.md` bei
  `mechanical-id-rewrite-misses-frozen-classes`.
- **Der `citations`-Bruch ist die planmäßige Rot-Quelle, nicht ein Unfall.**
  Der Bump verschiebt Zeilenspannen; `citations` läuft fail-closed im inneren
  Loop und nimmt den `pre-commit`-Hook mit
  ([`MR-051`](../../../../harness/conventions.md#mr-051)). — **Ausgang:**
  eingetreten: sieben `d-check:cite`-Direktiven zeigten `citation-mismatch`
  (sechs reine Zeilenverschiebung, neu geankert; eine — in
  [`MR-056`](../../../../harness/conventions.md#mr-056) —
  ein echtes Zitat-Delta, nach [`MR-039`](../../../../harness/conventions.md#mr-039)
  in [`MR-072`](../../../../harness/conventions.md#mr-072) vermerkt statt am
  zitierenden Dokument nachgezogen). Alle sieben behoben, `citations` grün.
- **Ändert `v6.9.0`s `AGENTS.template.md` §4 (oder eine andere Sektion mit
  Pointer-Charakter) erneut die Form**, entsteht dieselbe Zwickmühle wie bei
  slice-222: Bump ohne Adoption ließe `AGENTS.md` einer Vorlage widersprechen,
  die jeder Lauf lädt. Abgrenzung 2 hält dagegen fest, dass die Übernahme
  einer erneuten Weisung bedarf, keiner stillen Wiederholung. — **Ausgang:**
  entfallen: `diff` zwischen `v6.6.0`s und `v6.9.0`s `AGENTS.template.md`
  zeigt keine Änderung an §4 (nur Release-URL im Kopf und eine
  Kennungs-Terminologie fernab von §4) — Abgrenzung 2 hält ohne erzwungene
  Ausnahme.

## 7. Closure-Notiz

**Geliefert.** Der Baseline-Pin steht auf `v6.9.0` (54 Dateien, `verify ok`),
[`MR-072`](../../../../harness/conventions.md#mr-072) trägt die Hebung als
vierzehnter Nachtrag der Serie, [`MR-071`](../../../../harness/conventions.md#mr-071)
liegt in `harness/conventions/done/`. Alle vier Spiegel-Klassen sind
nachgezogen (Pfad-Verweise, Release-/Tree-URLs, bare Versionsnennungen,
`d-check:cite`-Direktiven), die neun eingefrorenen Lauf-Belege/Vergangenheits-Aussagen unangetastet (eine erst im Review gefunden).
`make gates` grün — zehn Gates, 783 Dateien, 0 Befunde.

**Was funktioniert hat: die Lehren aus slice-222 wurden vorab in den Plan
geschrieben, nicht erst im Review gefunden.** Die Suchform „nach der Version,
gruppiert nach Präfix" fand alle vier Mirror-Klassen auf Anhieb; das
Vier-Klassen-Raster aus [`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
strukturierte die Suche, statt sie dem Zufall zu überlassen.

**Was Friktion war, und wie beim Vorgänger: die mechanische Ersetzung ging
einmal daneben — diesmal aber am eigenen Fund erkannt, nicht erst im
Review.** Ein pauschaler `sed` über `harness/conventions.md` traf die
Tabellenzeile, mit der die Datei selbst [`MR-071`](../../../../harness/conventions.md#mr-071) beschreibt — dieselbe
Über-Hebungs-Klasse wie [`MR-067`](../../../../harness/conventions.md#mr-067) in slice-222, nur an einer Stelle, die
[`MR-070`](../../../../harness/conventions.md#mr-070)s Frozen-Liste
strukturell nicht abdeckt: eine Zeile *innerhalb* einer unbestreitbar
lebenden Datei, die über eine eingefrorene Datei berichtet. Fünftes
Auftreten von [`mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md).

**Steering-Loop-relevanter Fund: ein echtes Zitat-Delta, nicht nur eine
Zeilenverschiebung.** Sechs der sieben `citation-mismatch`-Befunde waren
reine Zeilenverschiebung (Wortlaut unverändert, neu geankert). Der siebte
([`MR-056`](../../../../harness/conventions.md#mr-056)) traf ein Zitat, dessen Quellsatz seit `v6.9.0` eine Ausnahme
trägt, die es vorher nicht gab (§Ein Slice, dessen Gegenstand ein anderer
übernimmt). Nach [`MR-039`](../../../../harness/conventions.md#mr-039)
bleibt das Zitat in [`MR-056`](../../../../harness/conventions.md#mr-056) unangetastet stehen, die Direktive ist entfernt,
und der Delta ist in [`MR-072`](../../../../harness/conventions.md#mr-072) vermerkt — die Regel hat zum ersten Mal seit
ihrer Einführung tatsächlich gegriffen.

**Zweiter Fund: die `versions.patterns`-Lücke bleibt real, aber ihre Größe
schwankt.** Anders als beim Vorgänger (ein neuer `ignore-refs`-Eintrag für
vier `target-missing`-Befunde) brauchte dieser Bump **keinen** neuen
Eintrag — `make doc-check` meldete 0 Befunde direkt nach dem Entfernen des
alten Baums. Der Unterschied liegt nicht an sorgfältigerer Arbeit, sondern
daran, welche **Form** die eingefrorenen Lauf-Belege zufällig zitieren
(Markdown-Link vs. Inline-Code/Prosa) — ein Datenpunkt gegen die Annahme,
jeder Bump brauche zwingend einen neuen Eintrag.

**Was diesmal NICHT eintrat: die Zwangslage aus slice-222.** Der
`AGENTS.template.md`-Diff zeigt keine erneute Änderung an §4; Abgrenzung 2
(keine Template-Adoption) hielt ohne Auftraggeber-Eingriff. Der große
inhaltliche Delta dieses Bumps (ID-Schema-Generalisierung in
`grundlagen-source-precedence.md`, neuer vierter Lifecycle-Zweig in
`modul-05-planning-harness.md`) ist gemessen und in [`MR-072`](../../../../harness/conventions.md#mr-072) gelistet,
bewusst nicht beurteilt — Sache des Folge-Slice.

**Review-Runde 1: 0 HIGH, 1 MEDIUM, 1 LOW, 0 INFO** — der Report liegt unter
[`docs/reviews/2026-09-16-slice-224-baseline-v690-review-r1.md`](../../../reviews/2026-09-16-slice-224-baseline-v690-review-r1.md).
Beide Findings sind eigene Klassen, die dieser Slice teils schon kannte:

**F-1 (MEDIUM) — ein neunter lebender Fundort war weder retargetet noch als
eingefroren deklariert.** Der Kommentar bei `.d-check.yml:235` trägt `v6.6.0`
als Begründung für einen anderen, bereits bestehenden `ignore-refs`-Eintrag
(den elften, aus dem `v6.5.0`→`v6.6.0`-Bump) — inhaltlich korrekt
unverändert, aber in der ersten Frozen-Liste dieses Slice schlicht
ausgelassen, nicht bewusst ausgenommen. Behoben: [`MR-072`](../../../../harness/conventions.md#mr-072)
nennt jetzt neun statt acht Frozen-Stellen (50 statt 51 lebende), und
`evidence/slice-224.md` bei
[`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
trägt den Fund nach — **derselbe Vorgang**, kein achtes Auftreten, da
Modul 6 einen Vorgang nur einmal zählt. **Lerneintrag:** Eine
Konfigurationsdatei mit genau einem frozen-artigen Kommentar entgeht einer
datei-basierten Frozen-Liste so leicht wie eine Tabellenzeile innerhalb
einer lebenden Doku-Datei (die Über-Hebungs-Falle oben) — dieselbe Grenze
von [`MR-070`](../../../../harness/conventions.md#mr-070), an einer noch
feineren Stelle.

**F-2 (LOW) — eine unbeteiligte Slice-Plan-Datei trägt jetzt eine falsche
Gegenwarts-Aussage.** [slice-223](../open/slice-223-commit-zerlegung-ausnahmen-aufloesen.md) §4 sagt
weiterhin *„der Pin steht auf `v6.6.0`"* als Trigger-Bedingung — mit diesem
Bump ist das nicht mehr aktuell. Der Fund liegt außerhalb dessen, was dieser
Slice an sich selbst bindet (kein Abgrenzungspunkt nennt slice-223), und
bleibt bewusst unangefasst — die nächste Beanspruchung von slice-223 liest
den Satz ohnehin als Datumsanker vom 2026-09-08, nicht als aktuelle
Bedingung, und trägt die Korrektur dann selbst nach.

**Die drei Paarungen, gemessen.** **(a) Anker** — vakant: Der Slice
verkörpert keine neue Steering-Loop-Regel; seine Lerneinträge liegen bei
bestehenden Registereinträgen. **(b) Folge-Slice** — kein neuer; der
bekannte Folge-Slice für den Regel-Delta ist noch nicht angelegt (Sache
eines künftigen Adoptions-Slice, analog slice-208). **(c) Register** — alle
zitierten Pfade lösen auf; zwei neue Belege liegen als
`evidence/slice-224.md` in ihren Verzeichnissen
(`pin-bump-mirrors-ungated`, `mechanical-id-rewrite-misses-frozen-classes`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende. Dieses Repo führt **drei** Prüfungen — die
zwei kanonischen und, als Adaption, den Nachtlauf-Stand
([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Die drei Vorprüfungen sind bei der Beanspruchung am 2026-09-16 bestätigt** —
Anlage und Beanspruchung fallen in dieser Sitzung zusammen; Nachtlauf und
Register sind gegen denselben Stand gelesen, der unten steht.

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:363-364 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area: `*` (Repo-Default). Der Slice hebt einen vendorten
Bestand und zieht pin-gebundene Verweise nach; `tools/harness/` ist nicht
berührt (das Werkzeug selbst ändert sich nicht, nur sein Ziel-Tag).

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:369-369 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, **40** Verzeichnisse). **Vier
Einträge sind einschlägig, alle vier aus dem Vorgänger-Bump bekannt:**

- [`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
  (7×, **kein formgültiger Ausgang**) — der tragende Eintrag für DoD (2).
- [`mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
  (5×, verkörpert als [`MR-070`](../../../../harness/conventions.md#mr-070))
  — trägt DoD (3); die Grenze der Regel (Frozen-Liste über Eigenschaften kann
  eine im *Inhalt* versteckte Vergangenheits-Aussage nicht fangen) ist
  benannt, nicht behoben — §6 trägt sie deshalb erneut als Risiko.
- [`form-vom-nachbarn-statt-von-der-vorlage`](../observations/BEO-ALL/form-vom-nachbarn-statt-von-der-vorlage/observation.md)
  (2×, offen unter der Schwelle) — Mahnung, die Form aus der **Vorlage**
  abzuleiten, nicht aus dem Bestand; hier vor allem für die
  Abgrenzung-2-Frage in §3 Schritt 3 relevant.
- [`liefer-punkt-in-fremdem-commit`](../observations/BEO-ALL/liefer-punkt-in-fremdem-commit/observation.md)
  (2×, offen unter der Schwelle) — Mahnung, jeden Commit gegen `git show`
  statt aus der Erinnerung zu beschreiben, besonders bei einem
  Kontext-Wechsel mitten in der Arbeit.

**Keiner der vier erreicht mit diesem Slice die Schwelle erstmalig.**

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-16 gelesen: `upstream-drift.yml` **planmäßig
rot** (2026-09-16T05:33:02Z) — genau der Fremd-Release-Befund, der dieses
Slice auslöst, keine unerwartete Störung. `image-scan.yml` grün
(2026-09-16T08:37:54Z), ohne Bezug zu diesem Slice.

**Modus-Begründungsblock.** Alle berührten Sub-Areas GF — ein Block genügt.

### Sub-Area: `*`

- **Modus:** GF (Repo-Default).
- **Konventions-Dichte:** **sehr hoch.** Die Pin-Serie ist dreizehnmal
  gelaufen ([`MR-011`](../../../../harness/conventions.md#mr-011) bis
  [`MR-071`](../../../../harness/conventions.md#mr-071)), das Vorgehen ist in
  fünf `MR`-Einträgen verankert, und der Vorgänger-Slice
  ([slice-222](../done/slice-222-baseline-v660-bump.md)) liegt als Vorlage
  samt sieben dokumentierten Lerneinträgen vor.
- **Phase-Reife:** Phase 5.
- **Evidenz-/Diskrepanz-Risiko:** **niedrig für den vendorten Baum**
  (`baseline-verify` prüft ihn hart), **mittel für die Spiegel** — drei der
  vier Klassen sind gate-blind, und der Registereintrag dazu steht
  unverändert bei 7× ohne Ausgang.
- **Reconciliation-Aufwand:** keiner (GF).
