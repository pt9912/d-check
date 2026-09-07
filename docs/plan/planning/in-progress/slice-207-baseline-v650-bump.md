# Slice slice-207: Baseline-Pin auf `v6.5.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre. Die Belege sind je Slice (`baseline-verify`
grün, `gates` grün); ein Wellen-Trigger schriebe sie ab. Präzedenz: die
Pin-Hebung über **vier** Tags lief ebenso wellenlos.

**Bezug:** [`MR-011`](../../../../harness/conventions.md#mr-011)-Kette (Pin auf
Release-Tag), [`MR-023`](../../../../harness/conventions.md#mr-023)
(Bundle-Layout), [`MR-021`](../../../../harness/conventions.md#mr-021)
(pin-gebundene Verweise), [`MR-051`](../../../../harness/conventions.md#mr-051)
(`d-check:cite`-Spannen beim Bump neu ankern),
[`MR-055`](../../../../harness/conventions.md#mr-055) (Symlink als Träger).

**Berührte Spec-Stellen:** — *(keine; der Slice bewegt den Baseline-Pin und
die pin-gebundenen Verweise, keine Anforderung und keine Sicht)*

**Verantwortlich:** pt9912 · **Autor:** pt9912. **Datum:** 2026-09-06.

---

## 1. Ziel

Den vendorten Baseline-Bestand von `v6.3.1` auf `v6.5.0` heben und alle
pin-gebundenen Verweise nachziehen — **mechanisch, ohne den Regel-Delta zu
beurteilen**. Was der Delta inhaltlich verlangt, entscheidet der Folge-Slice.

## 2. Vorgehen

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`.harness/baseline/`](../../../../.harness/baseline/) | neu + entfernt | `v6.5.0` materialisieren, `v6.3.1` entfernen — `fetch-baseline-cache.sh` trägt beides |
| pin-gebundene Verweise | update | jeder `.harness/baseline/<tag>/`-Pfad in lebenden Dokumenten, plus die Release-/Tree-URLs mit dem Tag |
| `d-check:cite`-Direktiven | update | die Zeilen-Spannen verschieben sich mit dem Bump; `citations` ist fail-closed und läuft im inneren Loop |
| neuer Eintrag unter [`harness/conventions/`](../../../../harness/conventions/) | neu | Nachtrag zur Pin-Serie, mit dem gemessenen Delta-Umfang; die Kennung wird beim Schreiben vergeben, nicht hier |

**Der Sprung überspringt zwei Releases** (`v6.4.0`, `v6.5.0`) — und das ist der
Grund für die erste Zeile der DoD. Ein direkter Byte-Vergleich über zwei
Versionen produziert Rauschen (Versionsnummern, Datumszeilen, verschobene
Zeilen), das den echten Regel-Delta verdeckt; gemessen wird deshalb mit
`diff -I`, die Lehre aus der vorigen Pin-Hebung.

**Was `v6.4.0` bringt, ist bekannt und angekündigt:** die Umsetzung der beiden
angenommenen Bitten unseres ausgehenden CR
([Antwort](../../cr/2026-09-06-antwort-ai-harness-course-slice-formluecken.md))
— der Ausschluss-Abschnitt als zweite Hälfte von §1, §8 mit unbedingtem Kopf
und bedingtem Rumpf, beide Prosa-Pflaster entfernt. **Was `v6.5.0` bringt, ist
unbekannt**; das misst dieser Slice.

### Der gemessene Delta (DoD 1)

**Die Messmethode zuerst, wie der Register-Eintrag es verlangt.** Ein
**Regel**-Delta ist eine Zeile, die ändert, was ein Implementer tun oder lassen
muss — eine Pflicht, ein Verbot, eine Form-Vorgabe, eine Ziel-Form. **Rauschen**
ist alles, was die Verpflichtung unberührt lässt: Versions- und Datumsangaben,
Tabellen-Padding, nachgezogene Querverweise.

**Und die Methode des Vorgängers reichte nicht.** Die Lehre aus der letzten
Hebung lautete *„Bundle-Delta nur mit `diff -I` messen"* — das filtert Versionen
und Daten. Damit blieben **27** Markdown-Dateien mit Delta stehen, und die
größte (`grundlagen-begriffe.md`, 81 geänderte Zeilen) war **bis auf eine Zeile**
Tabellen-Padding: Die Glossar-Tabelle wurde von Ausricht-Leerzeichen befreit; inhaltlich kam **eine** Zeile hinzu (RTM), und genau deshalb steht die Datei unten unter den zwölf.
Erst `diff -w -B -I` trennt das:

| Messung | Dateien mit Delta |
|---|---|
| roh (`diff -rq`) | 35 Pfade |
| mit `-I` (Versionen, Daten) | 27 Markdown-Dateien |
| **mit `-w -B -I`** | **12 Markdown-Dateien** |

**Die Filter, ausgeschrieben — ohne sie ist die Tabelle nicht nachrechenbar:**
`-I 'v6\.[0-9]\+\.[0-9]\+'` (Versionsnummern, vor allem die Quell-URL im Kopf
jeder Datei), `-I '20[0-9][0-9]-[0-9][0-9]-[0-9][0-9]'` (Datumsangaben), `-w`
(Weißraum innerhalb der Zeile — das Tabellen-Padding), `-B` (eingefügte
Leerzeilen). Verglichen wird Datei für Datei über die Namensmenge des alten
Baums, nicht als Verzeichnis-Diff.

Fünfzehn Dateien waren reines Weißraum-Rauschen. Ohne den zweiten Filter hätte
der Folge-Slice über sie geurteilt. **Der Datei-Bestand ist unverändert:** 55
Dateien vorher wie nachher, keine neu, keine entfallen.

**Die zwölf, klassifiziert.** Zwei Gruppen, und die Trennlinie ist scharf:

**(A) `v6.4.0` — die Umsetzung unseres ausgehenden CR** (fünf Dateien):

1. `modul-05` §Ziel-Form: Slice — §1 heißt **Ziel und Abgrenzung**, mit den
   **vier Klassen** und der Begründungs-Pflicht je Punkt. **Plus eine Schärfung,
   die wir nicht erbeten haben:** *„Die Adresse muss die Sendung annehmen"* — ein
   Folge-Slice, der den verwiesenen Punkt selbst ausschließt oder **vor** dem
   verweisenden schließt, ist keine Adresse.
2. `slice.template.md` — §1 und §8 umbenannt (`Sub-Area-Prüfungen und
   Modus-Begründung`), Bedienhinweise, die zwei *Vorgelagert*-Blöcke als
   unbedingter Kopf.
3. `modul-09` — die Plan-Ausgabe nennt Out-of-Scope, **und eine neue Pflicht**:
   Nimmt der Lauf etwas mit, das §1 ausschließt, ist das eine **Plan-Änderung**
   und gehört vor den Code, nicht in den Bericht danach.
4. `modul-06` — Querverweis auf die neue §1-Form.
5. `templates/README.md` — beide Abschnitte beschrieben.

**(B) `v6.5.0` — neu, und die erste Hälfte trifft d-check ins Zentrum** (sieben
Dateien):

6. `grundlagen-traceability.md` — **neue Sektion** *„Die zweite Richtung:
   Anforderung → Beleg"*: die **RTM** als Kanon-Begriff, *„sie wird **erzeugt**,
   nicht gepflegt"*, und die Setzung, dass es eine **Konfigurationsentscheidung**
   ist, welche Verweis-Quelle eine Anforderung *entlastet*. Das beschreibt
   `--trace` — unser eigenes Werkzeug — und ist damit der inhaltlich schwerste
   Punkt des Deltas.
7. `grundlagen-begriffe.md` — RTM im Glossar (die einzige nicht-Rauschen-Zeile
   dieser Datei).
8. `grundlagen-harness-dateien.md` — **einfrierendes gegen lebendes Artefakt**:
   Ein lebendes verlinkt `harness/sensors/<target>.md`, ein einfrierendes nennt
   `make <target>` als **Token**. Einfrierend sind Review-Report, Closure-Notiz,
   Archiv-Stub, `Accepted`-ADR und geschlossener Slice.
9. `review-report.template.md` — die **Zitier-Form** dazu: Kennung statt
   Adresse, und eine Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als
   Link. Begründung des Kanons: Der vendorte Baum trägt genau einen Tag, und ein
   Link darauf färbt beim nächsten Bump ein Artefakt rot, das niemand mehr
   anfassen darf.
10.–12. `archiv-stub-slice`, `archiv-stub-welle`, `welle-results` — dieselbe
   Zitier-Form in den übrigen einfrierenden Vorlagen.

**Der Punkt 8/9 hat diesen Slice bereits eingeholt**, und das ist kein
Nebenbefund: Die `Accepted`-ADR-Ausnahme unten (§5) ist genau der Fall, den die
neue Kanon-Regel künftig gar nicht erst entstehen lässt. Die **Adoption** liegt
im Folge-Slice; hier ist sie nur gemessen.

## 3. Ausdrücklich NICHT in diesem Slice

- **Jede Regel-Adoption.** Ob und wie ein Delta-Punkt übernommen wird, ist ein
  Urteil je Regel und gehört in
  [slice-208](../open/slice-208-v650-regel-adoption.md). Dieser Slice **misst** den
  Delta und **hebt den Pin**; er entscheidet nichts.
- **Die Auflösung der Slice-Haus-Form.** Sie folgt aus dem `v6.4.0`-Delta und
  ist der Kern des Folge-Slice — inklusive der Gate-Regeln, die heute auf
  Haus-Form-Titel keilen.
- **Jede Änderung an Anforderung, Sicht oder ADR.** Ein Pin bewegt keine
  Spec-Stelle.

## 4. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** Der vendorte Bestand steht auf `v6.5.0`, `make baseline-verify`
      und `make baseline-probe` sind grün, und der **Delta gegenüber `v6.3.1`
      ist mit `diff -I` gemessen** und als Liste im Slice festgehalten —
      getrennt nach *Regelwerk*, *Templates* und *reines Rauschen*.
- [ ] **(2)** Alle pin-gebundenen Verweise zeigen auf `v6.5.0`
      ([`MR-021`](../../../../harness/conventions.md#mr-021)), die
      `d-check:cite`-Spannen sind neu geankert
      ([`MR-051`](../../../../harness/conventions.md#mr-051)), und die Aliase
      unter `.claude/rules/` lösen auf
      ([`MR-055`](../../../../harness/conventions.md#mr-055)).
- [ ] **(3)** Ein neuer Konventions-Eintrag trägt die Hebung als Nachtrag zur
      [`MR-011`](../../../../harness/conventions.md#mr-011)-Kette, mit dem
      **gemessenen** Delta-Umfang statt einer Schätzung.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §5 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 5. Abnahme-Punkte / Risiken

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Zwei Releases in einem Sprung — der Delta ist größer als bei jeder
  bisherigen Hebung mit bekanntem Inhalt.** Zeigt die Messung, dass `v6.5.0`
  eine eigene Regel-Änderung trägt, ist der Folge-Slice **vor** seiner
  Beanspruchung neu zu schneiden. Das ist kein Fehler, sondern der Grund,
  warum die Messung in diesem Slice liegt und nicht im nächsten. —
  **Ausgang:** \<offen\>
- **Die `d-check:cite`-Spannen sind die planmäßige Rot-Quelle**
  ([`MR-051`](../../../../harness/conventions.md#mr-051)). `citations` ist
  fail-closed und läuft im inneren Loop: eine nicht neu geankerte Direktive
  nimmt den `pre-commit`-Hook mit. Über **zwei** Versionen verschieben sich
  mehr Zeilen als über eine. — **Ausgang:** \<offen\>
- **Ein Verweis, den kein Gate hält**
  ([`BEO-ALL/pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)):
  Release-/Tree-URLs, Prosa-Pins und der **zitierende** Verweis, dessen Zitat
  am neuen Ziel nicht mehr existiert. Gate-blind in beide Richtungen —
  vergessene Hebung wie Über-Hebung. — **Ausgang:** \<offen\>

## 6. Trigger

**Start** (`open` → `in-progress`): `make baseline-freshness` meldet
`v6.4.0` und `v6.5.0` als neuere Releases (gelesen 2026-09-06); der Content am
gepinnten Tag ist unverändert. WIP-Limit frei.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Zeigt der gemessene Delta, dass allein das
  Neu-Ankern der Zitat-Spannen eine eigene Sitzung füllt, wird es ein eigener
  Slice vor der Pin-Hebung.
- `in-progress` → `open` (blockiert): Ist der Upstream-Bestand am neuen Tag
  nicht integer (`SHA256SUMS` passt nicht), ruht der Slice bis zur Klärung mit
  der Baseline — still weiter zu vendoren wäre der Verlust der Integritäts-Zusage.

**Closure-Trigger.** Zwei beobachtbare Kriterien und ein Lerneintrag: (a)
`make baseline-verify`, `make baseline-probe` und `make gates` sind grün; (b)
`make baseline-freshness` meldet den Pin als aktuell.

## 7. Vorgelagert (vor der Modus-Begründung)

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:268-269 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Zwei** Sub-Areas, und diesmal ist die Ausdifferenzierung nicht künstlich:
`*` (Repo-Default) für die pin-gebundenen Verweise in den lebenden Dokumenten,
**und** `tools/harness/` für das Werkzeug, das die Hebung ausführt. Die zweite
ist in [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration eigens geführt und über
[`MR-004`](../../../../harness/conventions.md#mr-004) konventionsgetragen. Ob
`fetch-baseline-cache.sh` selbst angefasst werden muss, entscheidet die
Delta-Messung — der Slice führt die Sub-Area deshalb als berührt, auch wenn
die Antwort „keine Änderung" lauten kann.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:274-274 -->

> **Offene Beobachtungen sichten.**

Register durchgegangen (gemergter Stand, **36** Verzeichnisse). Die erste Fassung
nannte **35** und zählte damit nur `BEO-ALL/` — das Register führt ein zweites
Kürzel (`BEO-HARN`, ein Eintrag). Dieselbe Klasse wie die Delta-Messung unten:
die Zählung traf einen Teilbaum, nicht den Gegenstand.
Gesucht nach **beidem**: dem Gegenstand des Slice und dem, was er **anfasst**.
Vier Einträge sind einschlägig:

- [`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
  (5×) — der unmittelbarste. Er benennt **vier** Spiegel-Klassen einer
  Pin-Hebung und sagt, dass nur die grep-bare gehoben wird: Release-/Tree-URLs
  mit dem Tag, Prosa-/Ellipsen-Pins und der **zitierende** Verweis, dessen
  Zitat am neuen Ziel nicht mehr existiert. Gate-blind in **beide** Richtungen
  — vergessene Hebung wie Über-Hebung. Er ist der Grund, warum DoD (2) die
  Verweise **und** die Zitat-Spannen nennt und nicht nur die Pfade.
- [`semantic-change-body-only-edges-stale`](../observations/BEO-ALL/semantic-change-body-only-edges-stale/observation.md)
  (11×, mit dem Vorgänger-Slice zuletzt gewachsen) — dieselbe Klasse eine Ebene
  höher, und seine **Präzisierung aus slice-206** trifft hier direkt: Ein
  Spiegel, der eine **abgeleitete Aussage** ist statt eines Verweises, lässt
  sich nicht per `grep` finden. Bei einer Pin-Hebung sind das die Sätze, die
  über den *Inhalt* des gepinnten Stands reden, ohne ihn zu zitieren.
- [`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
  (15×) — die `d-check:cite`-Spannen werden neu geankert, und ein neu
  geankertes Zitat kann auf eine Zeile zeigen, die *ähnlich* aussieht, aber
  einen anderen Geltungsbereich hat. `citations` prüft die **Wortgleichheit**,
  nicht den Geltungsbereich.
- [`zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
  (1×, aus slice-205) — DoD (1) **ist** eine Zählmethode: `diff -I` misst den
  Delta, und aus seinen Treffern wird die Regel-Liste abgeleitet, auf der
  [slice-208](../open/slice-208-v650-regel-adoption.md) urteilt. Der Ableiter
  des Eintrags gilt wörtlich — vor der Messung die **Form** des Gegenstands
  ausschreiben (was ist eine *Regel*-Änderung, was Rauschen?) und die
  Trefferliste stichprobenweise dagegen halten, nicht nur die Zahl.

**Geprüft und ausgeschlossen:**
[`modulliste-spiegel-ungegated`](../observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md)
(2×) — der Slice fasst keine Modulliste an. Keiner der vier erreicht mit
diesem Slice die Schwelle von 3× erstmalig.

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-07 gelesen. `image-scan.yml` **grün**.
`upstream-drift.yml` **ROT**, und die Ursache ist der Anlass dieses Slice: Der
Pin steht auf `v6.3.1`, upstream liegen **zwei** Releases (`v6.4.0`, `v6.5.0`).
Alle fünf Versions-Achsen und der **Content-Drift am gepinnten Tag** sind
grün — rot ist allein die Currency-Achse. Der Slice schließt sie.

## 8. Sub-Area-Modus-Begründung

**Sub-Area `*` — Modus: GF** (Greenfield, Repo-Default).

- **Konventions-Dichte:** hoch. Die Pin-Serie ist über die
  [`MR-011`](../../../../harness/conventions.md#mr-011)-Kette elfmal gelebt,
  das Bundle-Layout steht in
  [`MR-023`](../../../../harness/conventions.md#mr-023), die Verweis-Bindung in
  [`MR-021`](../../../../harness/conventions.md#mr-021), das Neu-Ankern in
  [`MR-051`](../../../../harness/conventions.md#mr-051). Es gibt keine offene
  Form-Frage.
- **Phase-Reife:** Phase 5. Der Vorgang ist prozedural vollständig beschrieben
  und werkzeug-getragen.
- **Evidenz-/Diskrepanz-Risiko:** **erhöht, und zwar messbar höher als bei
  jeder bisherigen Hebung.** Der Sprung überspringt zwei Releases; die
  gesichteten Einträge `pin-bump-mirrors-ungated` (5×) und
  `semantic-change-body-only-edges-stale` (11×) beschreiben beide genau die
  Klasse, die dabei still bleibt. Das Risiko liegt nicht im Vendoring — das
  hält `baseline-verify` —, sondern in den Verweisen, die kein Gate deckt.
- **Reconciliation-Aufwand:** keiner (GF). Graduation entfällt.

**Sub-Area `tools/harness/` — Modus: GF**, Kürzel `HARN`.

- **Konventions-Dichte:** hoch —
  [`MR-004`](../../../../harness/conventions.md#mr-004) trägt die
  Harness-Mechanik.
- **Phase-Reife:** Phase 5; das Werkzeug hat elf Hebungen getragen.
- **Evidenz-/Diskrepanz-Risiko:** niedrig, aber **nicht null**: Ändert das
  Bundle sein Layout, greift die tolerante Entpackung ins Leere. Das zeigt die
  Delta-Messung, und `baseline-verify` fängt es fail-closed.
- **Reconciliation-Aufwand:** keiner (GF).

## 9. Closure-Notiz (nach `done/`)

\<wird vor dem `git mv` nach `done/` gefüllt\>
