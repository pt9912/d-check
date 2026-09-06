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

**Berührte Spec-Stellen:** neue `DC-FA-*`-Anforderung in
[`spec/lastenheft.md`](../../../../spec/lastenheft.md) §3 (Bereichs-Kürzel wird
mit ihr vergeben); Versions-Bump 0.84.0 → 0.85.0 samt Historie-Zeile.

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
- **Die Gegenprobe an d-check läuft grün:** Alle 20 Module sind im Handbuch
  erwähnt. Wir bauen für einen Bedarf, den wir selbst nicht haben — mit der
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
      Bereichs-Kürzel in §3, Versions-Bump auf 0.85.0 und Historie-Zeile
      (Baseline-Regelwerk `modul-03-spec.md`).
- [ ] **(2)** Eine neue ADR <!-- d-check:ignore (wird mit diesem Slice angelegt) --> trägt die **zwei** Entscheide mit ihren
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
`make gates` ist grün; (b) die begleitende ADR <!-- d-check:ignore (wird mit diesem Slice angelegt) --> ist `Accepted` und im ADR-Index verlinkt.

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
