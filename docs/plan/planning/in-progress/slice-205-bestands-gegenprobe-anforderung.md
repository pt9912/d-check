# Slice slice-205: Bestands-Gegenprobe — Anforderung und Entscheid

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** eingehender Change Request des Adopters `ai-harness-init`
([`docs/plan/cr/`](../../cr/2026-09-06-cr-eingehend-ai-harness-init-rtm-soll-ist.md)),
angenommen in Vorschlag **A**; berührt
[`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
(RTM) und
[`DC-FA-COV-001`](../../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
nur als **Abgrenzung**, nicht als Änderung.

**Berührte Spec-Stellen:** [`DC-FA-MENT-001`](../../../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in) (mit diesem Slice vergeben) in
[`spec/lastenheft.md`](../../../../spec/lastenheft.md) §3 (Bereichs-Kürzel wird
mit ihr vergeben); Versions-Bump samt Historie-Zeile — **auf 0.86.0, nicht wie hier zunächst geplant auf 0.85.0**: Der Plan las das Kopf-Feld, das auf 0.84.0 stand; die Historie darunter führte bereits eine 0.85.0-Zeile (§5).

**Verantwortlich:** pt9912 · **Autor:** pt9912. **Datum:** 2026-09-06.

---

## 1. Ziel

Die vom CR beschriebene Achse als **Anforderung** festschreiben und den einen
Entscheid treffen, den sie verlangt: Eine konfigurierte Artefakt-Menge wird
gegen ein Ist-Dokument gehalten, und gemeldet wird, welche Mitglieder dort
**nicht vorkommen**. Die RTM misst *verfolgt*, diese Achse misst *erwähnt* —
zwei Fragen, zwei Antworten.

## 2. Vorgehen

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`spec/lastenheft.md`](../../../../spec/lastenheft.md) | update | neue `DC-FA-*` mit Akzeptanzkriterien-Trio, Bereichs-Kürzel in §3, Versions-Bump + Historie |
| neue ADR unter [`docs/plan/adr/`](../../../../docs/plan/adr/) | neu | der Entscheid **eigene Achse gegen dritte Quell-Form von `trace.cross-consistency`**, plus Abgrenzung zu `targets` und zur RTM |
| CR-Dokument | update | Stand von *eingegangen* auf *angenommen (Vorschlag A)*, mit Verweis auf Anforderung und ADR |
> **Zweimal berichtigt (Review-Runden 1 und 2, M-8 / R2-M-1).** Der folgende
> Block steht im Zustand der Planung und ist in **beiden** Fassungen falsch;
> er bleibt als Lauf-Beleg stehen. Gemessen mit einer Methode, die verlangt,
> was eine Kennung ausmacht — Präfix **plus Ziffernblock** —, gilt: **zwölf**
> Repos führen beide Seiten, **neun** davon mit kennungstragenden
> Überschriften in **vier** Schemata (`LH` sechsmal, dazu `HSM`, `DC`, `AC`),
> **drei** ohne (`b-trace`, `cmake-xray`, `m-trace`). Die hier genannte Zahl
> zehn, vier der zehn Anforderungs-Zählungen und das Schema `S-*` sind falsch;
> das in der ersten Berichtigung ergänzte fünfte Schema `CLI`/`CORS`/`MVP`
> ebenfalls — es entstand daraus, dass die Zählmethode deutsche Komposita
> (`MVP-Abgrenzung`, `CORS- und CSP-Grundregeln`) für Kennungen hielt. Der
> Stand steht in [ADR-0084](../../adr/0084-mentions-eigenes-modul.md)
> §Geschichte. **Die Richtung des Arguments überlebt beide Berichtigungen** und
> wird stärker: vier Schemata über neun Repos plus drei ganz ohne
> Heading-Kennung sind für eine kennungsbasierte Achse genau das Problem, das
> die pfadbasierte nicht hat.


**Der Bestand, gemessen — und er trägt die Anforderung, nicht der Anlass.**
**Zehn** Repos unter `/Development` führen das Paar Soll (`spec/lastenheft.md`)
und Ist (`docs/user/`): `k-deskflight` (185 Anforderungen), `u-boot` (139),
`bedrock-eu-guard` (114), `c-hsm-doc` (89), `d-check` (51), `b-cad` (23),
`a-check` (21), `belief-agent` (15), `ai-harness-init` (14), `m-trace` (7);
zwei weitere (`cmake-xray`, `b-trace`) führen die Ist-Seite ohne erkennbares
Soll-Schema. Die Konvention ist damit keine Eigenheit des Einreichers.

**Und die Inventur liefert das stärkste Argument für A gleich mit:** Die
Soll-Seiten führen **fünf verschiedene ID-Schemata** (`LH-*`, `DC-*`, `AC-*`,
`HSM-*`, `S-*`). Eine Achse, die über **Kennungen** arbeitet, müsste jedes
davon kennen — `trace.coverage` tut das über Konfiguration. Vorschlag A
arbeitet über **Artefakte** und ist damit schema-unabhängig: derselbe
Mechanismus trägt alle zehn Repos ohne eine einzige Schema-Angabe. Diese
Eigenschaft nennt der CR nicht; sie ist der eigentliche Grund, warum die Achse
sich zu bauen lohnt.

**Drei Messungen, und die dritte korrigiert die erste:**

- **Die CR-Prämisse stimmt:** d-checks eigenes Handbuch trägt **5**
  `DC-*`-Nennungen bei **51** Anforderungen. Die Konvention „Ist-Dokument
  zitiert IDs" existiert auch hier nicht — `trace.coverage` kann diese Achse
  also prinzipiell nicht bedienen.
- **Die Gegenprobe an d-check läuft grün:** Alle Module sind im Handbuch
  erwähnt (zur Planungszeit 20, inzwischen 22). Wir bauen für einen Bedarf, den wir selbst nicht haben — mit der
  Folge, dass **eigenes Rauschen zum Kalibrieren fehlt**.
- **Eine erste Stichprobe an `u-boot` meldete „nichts gefunden" — sie prüfte
  nichts.** Der Glob traf dort keine Datei; null geprüfte Mitglieder lasen
  sich wie null Funde. Mit einer echten Menge (19 Make-Targets gegen **alle
  acht** Ist-Dokumente unter `docs/user/`) stehen **6** davon in keinem —
  `build-binaries`, `clean`, `compile`, `docs-check`, `fullbuild`,
  `verify-depguard`. Ob alle sechs ein *Fund* sind, entscheidet die
  Mengen-Wahl; dass die Achse etwas findet, ist damit belegt.

**Zwei Anforderungs-Merkmale folgen direkt aus diesen Messungen** und stehen
nicht zur Disposition:

1. **Die Ist-Seite ist eine Menge, kein Dokument.** `u-boot` führt acht
   Dateien unter `docs/user/`; wer gegen eine prüft, misst an sieben vorbei.
   Der CR spricht von „einem Ist-Dokument" — das ist zu eng.
2. **Fail-closed bei leerer Prüfmenge.** Der Fehler oben ist die gefährlichste
   Form eines grünen Laufs: eine Zusage über nichts. d-check hält das an
   anderer Stelle bereits so (`reviews`, `workflow-pins`, `requirements`).

## 3. Ausdrücklich NICHT in diesem Slice

- **Die Implementierung.** Modul, Grund-Codes und Konfiguration sind ein
  eigener Slice. Der Grund ist die Größenregel und
  [`MR-066`](../../../../harness/conventions.md#mr-066): Anforderung, ADR und
  Implementierung in einem Slice überschritten die Ein-Sitzungs-Grenze, und
  eine Ersatz-Form wäre hier ohne Not gewählt.
- **Vorschlag B des CR** (`trace.coverage` mit externer Mapping-Datei). Der
  CR nennt ihn selbst nachrangig: Er löst den **deklarierten** Fall und findet
  keine unbekannten Lücken — also genau das nicht, was den Fund produziert hat.
  Er wird im CR-Dokument als *zurückgestellt* vermerkt, nicht als abgelehnt.
- **Jede Änderung an `trace.coverage` oder der RTM.** Beide bleiben
  byte-identisch; die neue Achse steht daneben.

## 4. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** Die Anforderung steht in
      [`spec/lastenheft.md`](../../../../spec/lastenheft.md), nur dort, mit dem
      Akzeptanzkriterien-Trio (Happy Path · Boundary · Negativ), dem
      Bereichs-Kürzel in §3, Versions-Bump und Historie-Zeile
      (Baseline-Regelwerk `modul-03-spec.md`).
- [ ] **(2)** [ADR-0084](../../adr/0084-mentions-eigenes-modul.md) trägt die **zwei** Entscheide mit ihren
      verworfenen Alternativen: **(a) eigene Achse oder dritte Quell-Form von
      [`trace.cross-consistency`](../../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)**
      — dort existiert der Mengenabgleich samt Richtungslabel und Gate-Modi,
      die Mengen unterscheiden sich aber (ID-zentriert gegen
      artefakt-zentriert); **(b) welcher Gate-Modus der Default wird** — das
      Muster *Bericht per Default, Verdikt per Konfiguration* ist etabliert,
      also ist zu entscheiden, nicht zu erfinden. Die Abgrenzung zu `targets`
      gehört ebenfalls hinein.
- [ ] **(3)** Das CR-Dokument führt den Stand *angenommen (A), B
      zurückgestellt* mit Verweis auf Anforderung und ADR — die Antwort steht
      bei ihrer Bitte.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §5 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 5. Abnahme-Punkte / Risiken

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Die Mengen-Wahl ist das Urteil, nicht die Differenz.** Zwei Stichproben
  zeigen es: In `u-boot` fand die Gegenprobe über zwölf Werkzeuge **nichts**;
  in `ai-harness-init` fand sie `conventions.md` — ein Harness-Artefakt, das im
  **Benutzer**handbuch vermutlich nichts verloren hat. Eine Anforderung, die
  die Menge nicht als konfigurierte, begründete Wahl führt, produziert
  Falsch-Positive statt Funde. Dieselbe Stichprobe lehrte zugleich die
  **Leermengen-Falle**: Der erste Lauf traf keine Datei und las sich wie ein
  grünes Ergebnis. — **Ausgang:** <offen>
- **„Bericht statt Gate" ist in d-check KEINE neue Klasse** — diese Annahme des
  ersten Plan-Entwurfs war falsch, gefunden vom Auftraggeber.
  [`DC-FA-XREF-001`](../../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
  führt den Mengenabgleich zweier Sichten bereits **mit konfigurierbarem
  Verdikt**: Modus `equal` gatet beide Differenzen, `superset` nur eine. Dazu
  ist `--trace` selbst advisory (Exit 0), und `--require-complete` schaltet das
  Verdikt zu. Das Muster **Bericht per Default, Gate per Konfiguration** ist
  also etabliert — die Frage ist nicht *ob*, sondern welcher Modus der Default
  wird. — **Ausgang:** entfallen (die Annahme traf nicht zu)
- **Der nähere Verwandte ist XREF, nicht `targets`.** XREF sagt über sich
  selbst: *„Die einzige neue Logik ist: eine Sicht invertieren und die Mengen
  diffen — kein neuer Parser."* Genau das ist auch Vorschlag A. Der
  **Unterschied** liegt in den Mengen: XREF ist **ID-zentriert** (`F(R)` und
  `B(R)` je Anforderung, damit schema-abhängig), A ist **artefakt-zentriert**
  (Pfad-Globs gegen Dokument-Text, schema-frei). Ob A deshalb eine dritte
  Quell-Form von `trace.cross-consistency` wird statt eines eigenen Moduls, ist
  **der zentrale Entscheid dieses Slice** — und er ist mit dem Fund offen, nicht
  entschieden. — **Ausgang:** <offen>
- **Kein eigenes Rauschen zum Kalibrieren.** d-checks Gegenprobe läuft grün;
  jede Schwelle und Ausnahme-Klasse müsste am Fremd-Repo justiert werden. Das
  widerspricht der gelebten Praxis, jede Regel vor der Aufnahme am **eigenen**
  Bestand zu messen. — **Ausgang:** <offen>
- **Zwei Deklarations-Drifts im eigenen Lastenheft, während der Arbeit
  gefunden — und beide sind das Rauschen, dessen Fehlen der Punkt darüber
  beklagt.** (i) Die Modul-Aufzählung in
  [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
  nannte 20 Module und ließ `workflows` und `reviews` aus, obwohl beide seit
  0.83.0 bzw. 0.84.0 eigene Anforderungen tragen. (ii) Das `**Version:**`-Feld
  stand auf 0.84.0, während die Historie darunter bereits eine
  0.85.0-Zeile führte — der Kopf-Bump unterblieb bei jenem Eintrag. Beide
  wurden **nicht still behoben**: (i) mit der Anforderung berichtigt, weil die
  Liste ohnehin anzufassen war und eine bekannt falsche Aufzählung darin
  schlechter gewesen wäre; (ii) durch den Bump auf **0.86.0** aufgelöst, statt
  wie geplant auf 0.85.0. **Der Plan selbst ist der Beleg für die Klasse:**
  seine DoD nannte 0.85.0, weil sie das Kopf-Feld las und nicht die Historie.
  — **Ausgang:** <offen>
- **Der Sichtungs-Schritt in §7 hat den Bestand nicht nach dem durchsucht, was
  der Slice anfassen würde** — und der erste Nachzug hat daraus die falsche
  Zuordnung gezogen (unabhängiger Review, M-6). Gesichtet wurden vier
  Beobachtungen, alle nach dem **Gegenstand** des Slice; die Drift aus (i)
  betrifft dagegen, was er **anfasst**. Sie ist **nicht** eine dritte Instanz
  von
  [`modulliste-spiegel-ungegated`](../observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md):
  jener Eintrag führt als Sub-Area `.d-check.yml`, `Makefile`, Gate-Doku und
  meint die Spiegel der **Profil**-Modulliste (`FOCUS_DISABLE`,
  Netzlos-Modulliste, Gate-Prosa) — `workflows` und `reviews` stehen in
  `.d-check.yml` gar nicht. Die Aufzählung in
  [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
  spiegelt eine **andere** Menge: die im Produkt gültigen Modulnamen. **Die
  3×-Behauptung ist damit zurückgezogen** — sie wäre ein Ausgang für eine
  Schwelle gewesen, die eine andere Menge betrifft, und genau die Klasse
  [`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md).
  **Die Zuordnung war dann selbst wieder zu schnell** (Review-Runde 2,
  R2-M-6): „weder der eine noch der andere Eintrag deckt sie" stand gegen 35
  Registerverzeichnisse, von denen ich zwei geprüft hatte. Es gibt einen, der
  wörtlich passt —
  [`semantic-change-body-only-edges-stale`](../observations/BEO-ALL/semantic-change-body-only-edges-stale/observation.md)
  (9 Belege, verkörpert als
  [`MR-025`](../../../../harness/conventions.md#mr-025), bewusst weiter offen):
  *„wo eine **Aufzählung eine Menge spiegelt**, gehört sie an ihre Quelle
  gebunden statt in eine Liste geschrieben."* Genau das ist die
  Modul-Aufzählung des Lastenhefts gegenüber der Modulmenge des Produkts. Die
  Closure trägt den Beleg dort ein. **Zweimal in Folge habe ich in derselben
  Sache zu wenig gelesen** — erst den Geltungsbereich eines Eintrags, dann den
  Bestand des Registers; beides sind Lese-Fehler, keine Urteils-Fehler. —
  **Ausgang:** <offen>

## 6. Trigger

**Start** (`open` → `in-progress`): Der Auftraggeber hat Vorschlag A
angenommen (2026-09-06) und die Bestands-Grundlage bestätigt: mehrere Repos
mit Soll/Ist-Paar. WIP-Limit frei.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Zeigt sich, dass die Erkennungsform
  (voller Pfad · Dateiname · beides) eine eigene Vorfrage mit Messung ist,
  wird sie ein eigener Slice vor der Anforderung.
- `in-progress` → `open` (blockiert): Erweist sich „Bericht statt Gate" als
  Entscheid, der die Modul-Architektur berührt (ein Modul ohne Exit-Code),
  ruht der Slice bis zu einem Architektur-Entscheid.

**Closure-Trigger.** Zwei beobachtbare Kriterien und ein Lerneintrag: (a) die
Anforderung steht mit vollständigem Akzeptanzkriterien-Trio im Lastenheft und
`make gates` ist grün; (b) [ADR-0084](../../adr/0084-mentions-eigenes-modul.md) ist `Accepted` und im ADR-Index verlinkt.

## 7. Vorgelagert (vor der Modus-Begründung)

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.3.1/regelwerk/modul-05-planning-harness.md:223-224 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area: `*` (Repo-Default). Der Slice berührt Spec-Stratum und ADR,
nicht Code und nicht `tools/harness/`. Eine Ausdifferenzierung ist nicht nötig.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.3.1/regelwerk/modul-05-planning-harness.md:229-229 -->

> **Offene Beobachtungen sichten.**

Register durchgegangen (gemergter Stand, 35 Verzeichnisse). Vier Einträge
betreffen diesen Gegenstand:

- [`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md)
  (7×) — **das Gegenmittel steht in §2**: Die Anforderung wird nicht aus dem
  einen CR gezogen, sondern aus einer Inventur von acht Repos mit Soll/Ist-Paar.
  Der Eintrag ist mit slice-204 zuletzt gewachsen; hier wird ihm bewusst
  entgegengearbeitet.
- [`module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
  (1×) — unmittelbar einschlägig: Die neue Achse liest **zwei** Eingaben, die
  sie unterschiedlich behandelt (die Artefakt-Menge über Globs, das
  Ist-Dokument als Text). Die Frage des Eintrags — *welche Eingaben liest das
  Modul, die es nicht scannt* — gehört in die Anforderung, nicht in die
  Implementierung.
- [`wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md)
  (7×) — die Anforderung beschreibt eine Prüfung, die es noch nicht gibt; ihr
  Wortlaut darf nicht mehr zusagen, als die spätere Implementierung tragen
  kann. Konkret betroffen: die Frage, was „kommt vor" heißt.
- [`spec-randbedingung-ohne-test`](../observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md)
  (1×) — das Akzeptanzkriterien-Trio ist die Antwort darauf; jede Randbedingung
  der neuen Anforderung braucht ihr Kriterium, nicht nur der Happy Path.

Keiner der vier erreicht mit diesem Slice die Schwelle von 3× erstmalig.

**Diese Sichtung war unvollständig, und der erste Nachzug hat daraus die
falsche Zuordnung gezogen** (§5, unabhängiger Review M-6): Gesichtet wurde nach
dem **Gegenstand** des Slice, nicht nach dem, was er **anfassen** würde — und
angefasst wird die Modul-Aufzählung des Lastenhefts. Der Nachzug erklärte sie
zum vierten Spiegel von
[`modulliste-spiegel-ungegated`](../observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md)
und dessen Schwelle für erreicht. **Das ist zurückgezogen:** Jener Eintrag
meint die Spiegel der **Profil**-Modulliste, und die beiden fehlenden Namen
stehen im Profil gar nicht. Was bleibt, ist eine echte, ungegatete Drift ohne
passenden Registereintrag; die Zuordnung entscheidet die Closure. **Beide
Fassungen bleiben stehen** — die ursprüngliche Sichtung, weil sie Lauf-Beleg
ist, und der falsche Nachzug, weil eine nachträglich geglättete Korrektur die
Klasse verstecken würde, die sie selbst instanziiert.

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-06 gelesen. `image-scan.yml` **grün**.
`upstream-drift.yml` **ROT** — unverändert der Lauf **vor** der Pin-Hebung aus
slice-202; die Ursache ist behoben, `make baseline-freshness` meldet lokal
beide Teile grün. Dieselbe benannte Grenze wie in den beiden Slices davor: Das
Target liest den **jüngsten** Lauf, nicht sein **Alter**.

## 8. Sub-Area-Modus-Begründung

**Modus:** `*` ist **GF** (Greenfield, Repo-Default) — Doc führt, Code folgt,
und hier besonders ausgeprägt: Der Slice schreibt die Anforderung **vor** der
Implementierung, was der Reihenfolge des Lebenszyklus entspricht.
**Konventionen-Dichte** hoch (Anlege-Prozess vollständig in `modul-03`
verankert). **Phase-Reife** hoch. **Evidenz-/Diskrepanz-Risiko** liegt nicht
im Code-vs-Doku-Abstand, sondern in der **Reichweite der Zusage**: Eine
Anforderung, die mehr verspricht, als ein Modul über eine Mengendifferenz
halten kann, ist der Fehler, den die gesichteten Beobachtungen benennen.

## 9. Closure-Notiz (nach `done/`)

<wird vor dem `git mv` nach `done/` gefüllt>
