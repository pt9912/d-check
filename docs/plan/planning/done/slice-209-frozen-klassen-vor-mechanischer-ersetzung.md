# Slice slice-209: Frozen-Klassen vor einer mechanischen Ersetzung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
(3× erreicht, Ausgang *geplant*),
[`MR-025`](../../../../harness/conventions.md#mr-025) (der Geschwister-Eintrag:
Spiegel **vor** dem Editieren auflisten),
[`MR-052`](../../../../harness/conventions.md#mr-052) (ein historisches Zitat
ist über die git-Historie prüfbar).

**Berührte Spec-Stellen:** — *(keine; der Slice verkörpert eine
Steering-Loop-Regel, er ändert keine Anforderung)*

**Verantwortlich:** pt9912 · **Autor:** pt9912. **Datum:** 2026-09-07.

---

## 1. Ziel

Die Regel schreiben, die drei Anlässe verlangt haben: **Wer mechanisch über
den Baum ersetzt, listet die Frozen-Klassen vorher auf — und die Liste ist
über eine Eigenschaft definiert, nicht über Verzeichnisse.**

## 2. Vorgehen

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| neuer Eintrag unter [`harness/conventions/`](../../../../harness/conventions/) | neu | die Regel als Adaption; die Kennung wird beim Schreiben vergeben |
| [`harness/conventions.md`](../../../../harness/conventions.md) | update | Index-Zeile |
| [`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/state.md) | update | Ausgang von *geplant* auf *verkörpert*, mit Herkunfts-Anker |

**Was die drei Anlässe gemeinsam haben, und es ist nicht das Verzeichnis:**
Beim ersten (Register-Formatmigration) blieben `Accepted`-ADR-Kerne und
gesendete CRs unbedacht, obwohl die drei *benannten* Frozen-Verzeichnisse
korrekt ausgenommen waren. Der zweite (Pin-Hebung) ist ein **anderer** Fall —
*(die erste Fassung schrieb „dasselbe Muster" und widersprach damit der eigenen
Evidence-Datei; der Review hat es gemessen)*: Getroffen war kein eingefrorenes
Verzeichnis, sondern eine **identifizierende Nennung in einem lebenden
Dokument** — die Plan-Tabelle nannte danach denselben Pfad zweimal, als neu und
als entfallend, und **beide lösten auf**, weshalb kein Gate es meldete. Beim
dritten (Pin-Hebung) gab es **gar keine** Ausnahme-Liste, und gehoben wurden
zehn `d-check:cite`-Direktiven in `done/`-Slices, deren Zitat gegen den alten
Tag geschrieben war. **Die Verzeichnis-Liste ist die falsche Abstraktion** —
maßgeblich ist die Eigenschaft *„zitiert den Stand seiner Zeit"*, und die
tragen auch Zeilen, die in keinem der drei Verzeichnisse stehen. Dass die drei
Anlässe **nicht** dasselbe Muster zeigen, ist dabei kein Einwand gegen die
Regel, sondern ihr Argument: Ein Muster ließe sich als Verzeichnis-Liste
fassen, drei verschiedene nicht.

**Der Kanon hat dazu seit `v6.5.0` etwas zu sagen**
(`grundlagen-harness-dateien.md` §Ein einfrierendes Artefakt …): Er
unterscheidet **einfrierend** von **lebend** und zählt die einfrierenden auf —
Review-Report, Closure-Notiz, Archiv-Stub, `Accepted`-ADR, geschlossener
Slice. Ob die Regel dieses Repos damit **entbehrlich** wird oder ob sie den
*Vorgang* regelt, den der Kanon nicht kennt, ist die erste Frage des Slice.

### Die Vorfrage, beantwortet (DoD 1)

**Nein — der Kanon macht die Regel nicht entbehrlich.** Er regelt eine andere
Größe, und die Differenz ist an drei Stellen belegbar statt behauptet.

**Was der Kanon regelt.** `v6.5.0` sagt: *„Ein einfrierendes Artefakt nennt ein
prozess-bewegtes bei seiner Kennung, nicht unter seiner Adresse"*, zählt die
einfrierenden auf (Review-Report, Closure-Notiz, Archiv-Stub, `Accepted`-ADR,
geschlossener Slice) und grenzt ausdrücklich ab: *„Die Grenze: Sie gilt für
einfrierende Artefakte."* Adressat ist der **Autor**, im **Moment des
Schreibens**; die Anweisung betrifft die **Form eines Verweises**.

**Was diese Regel regeln will.** Adressat ist, wer **mechanisch über den Baum
ersetzt**, in einem Moment **lange nach** dem Schreiben; die Anweisung betrifft
die **Ausschluss-Menge einer Massen-Operation**. Anderer Adressat, anderer
Zeitpunkt, andere Handlung.

**Beleg 1 — die vom Kanon vorgeschriebene Form schützt nicht.** *(Die erste
Fassung dieses Belegs war falsch und ist vom unabhängigen Review widerlegt
worden; sie behauptete, die zehn gehobenen `d-check:cite`-Direktiven **seien**
die Kanon-Form. Sind sie nicht: Der Kanon sagt *„Tag **und** Pfad in
Inline-Code"*, und eine `<!-- d-check:cite … -->`-Direktive ist ein
HTML-Kommentar, kein Inline-Code — sie ist eine dritte Form, weder Link noch
Backtick.)*

**Der Beleg trägt trotzdem, über eine andere Messung.** Die Kanon-Form ist
gegen eine Massen-Ersetzung **genauso wehrlos**, denn sie enthält den Pfad,
auf den eine Tag-Hebung skopiert. Gemessen im heutigen Bestand: **neun**
Vorkommen der Form `` `.harness/baseline/<alter-tag>/…` `` — mit Backticks,
also kanon-konform — stehen in eingefrorenen Artefakten, darunter
`done/slice-202`, `done/slice-203`, `done/welle-67-results` und vier
aufgelöste `MR`-Einträge. Sie tragen genau das Präfix, das eine Pin-Hebung
sucht. **Überlebt haben sie, weil ihr Verzeichnis ausgenommen war** — also
durch die Handlung, die diese Regel vorschreibt, und nicht durch ihre Form.

*(Auch die Gegenthese des Reviewers hält damit nicht: Er schrieb, die
Kanon-Form trage den Präfix `.harness/baseline/<tag>/` „gerade nicht" und
slice-207s Ersetzung hätte sie nicht berührt. Die neun Vorkommen tragen ihn.
Die Richtung seines Befundes stimmt, seine Begründung nicht — und das ändert
den Schluss nicht, sondern schärft ihn.)*

**Beleg 2 — eine der drei Instanzen liegt außerhalb des Kanon-Geltungsbereichs.**
Der Fund in slice-202 betraf eine identifizierende Nennung in einem
**lebenden** Dokument (die Plan-Tabelle, die den zu entfernenden Baum
benennt). Der Kanon nimmt lebende Artefakte ausdrücklich aus. Eine Regel, die
nur die einfrierenden deckt, hätte diese Instanz nie gefangen.

**Beleg 3 — der schwächste der drei, und er wird als solcher geführt.**
slice-195 traf `Accepted`-ADR-Kerne (vom Kanon genannt) **und gesendete CRs**
unter `docs/plan/cr/` (nicht genannt). *(Die erste Fassung schrieb, der Kanon
führe eine „reine Aufzählung" — das ist zu stark, und der Review hat es
gemessen: Der Kanon-Satz nennt die **Eigenschaft** mit, „Einfrierend sind die
**Zeitdokumente**", und leitet die fünf Klassen daraus ab. Ein Beleg, der dem
Kanon eine Aufzählung unterstellt und den Gegenbeleg dann aus der Eigenschaft
zieht, die im selben Satz steht, argumentiert gegen sich selbst.)*

**Was bleibt, ist schmaler und trägt trotzdem:** Der Kanon nennt die
Eigenschaft dem **Autor**, der sein eigenes Dokument einordnet — eine Frage
mit einer Datei als Gegenstand. Diese Regel verlangt dieselbe Ableitung vom
**Ersetzenden**, über den **ganzen Baum**, für Klassen, die er nicht
geschrieben hat. Dass die fünf Klassen den gesendeten CR nicht enthalten, ist
dafür ein Symptom, kein eigenständiger Beweis: Es zeigt, dass die Aufzählung
allein nicht reicht — nicht, dass der Kanon keine Eigenschaft kennt.

**Die Gegenprobe, und sie gehört dazu.** Die stärkste Lesart eines **Ja**
wäre: Stünde nirgends eine Adresse, gäbe es nichts zu über-heben. Beleg 1
widerlegt sie — der Kanon **schreibt** die Adressform hier vor, weil sie für
den *Leser* die sichere ist. Sie ist es nur nicht gegen `sed`. Damit steht
neben dem lokalen Eintrag eine Beobachtung, die den Kanon selbst betrifft; sie
gehört gemeldet, aber nicht in diesen Slice (§3, *anderer Vorgang*).

**Die direkteste Quelle wurde gesucht, nicht die bequemste**
([`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md),
15×). Das Regelwerk `v6.5.0` kennt **keine** Stelle über die mechanische
Massen-Ersetzung — geprüft über den ganzen `regelwerk/`-Baum; der einzige
Treffer zu *„repo-weit"* steht in `grundlagen-source-precedence.md` und gilt
dem Geltungsbereich eines `MR`-Eintrags, nicht einer Ersetzung. Es gibt also
keine nähere Quelle, die zu zitieren wäre.

**Die zweite Hälfte der Vorfrage, vom Review nachgereicht: gegen**
[`MR-069`](../../../../harness/conventions.md#mr-069)**.**
§8 stellt sie wörtlich — *„ein vierter Eintrag daneben muss sagen, was die drei
nicht sagen"* —, und die erste Fassung beantwortete sie nur gegen den Kanon.
[`MR-069`](../../../../harness/conventions.md#mr-069) trägt bereits eine
**gemessene Vier-Klassen-Liste** derselben Eigenschaft: 18 `Accepted`-ADRs,
sechs aufgelöste Konventions-Einträge, drei `done/`-Slices, ein gesendeter CR.
Die Frage ist damit real und nicht rhetorisch.

**Die Antwort ist die Achse *Zustand gegen Akt*.** [`MR-069`](../../../../harness/conventions.md#mr-069) beschreibt einen
**Zustand**: Wie viele Verweise stehen heute hinter dem Ventil, über welche
Tags, in welchen Klassen — eine Bestandsaufnahme, die mit jedem Bump wächst
und die er selbst als Gate-Senkung deklariert. [`MR-070`](../../../../harness/conventions.md#mr-070) schreibt einen **Akt**
vor: wie die Ausschluss-Menge **vor** einer Ersetzung gebildet wird. Zwei
Belege dafür, dass das nicht dieselbe Aussage ist:

- **Die Mengen decken sich nicht.** [`MR-069`](../../../../harness/conventions.md#mr-069)s Liste ist auf Verweise **in die
  vendorte Baseline** beschränkt — das ist sein Geltungsbereich. Die Instanz
  aus slice-202 (identifizierende Nennung in einem **lebenden** Dokument) und
  die Register-Kennungen aus slice-195 fallen nicht darunter; sie stehen nur
  unter [`MR-070`](../../../../harness/conventions.md#mr-070)s Test.
- **Die Richtung ist entgegengesetzt.** [`MR-069`](../../../../harness/conventions.md#mr-069) rechtfertigt, dass ein
  bereits entstandener toter Verweis **stehen bleibt**. [`MR-070`](../../../../harness/conventions.md#mr-070) soll
  verhindern, dass ein lebender Verweis **falsch wird**. Der eine verwaltet
  Schaden, der andere vermeidet ihn.

**Und die Grenze dieser Antwort:** Die beiden berühren sich trotzdem — wer
[`MR-070`](../../../../harness/conventions.md#mr-070) befolgt, erzeugt weniger Einträge für [`MR-069`](../../../../harness/conventions.md#mr-069)s Ventil. Das ist eine
Wirkung, keine Dopplung; aber es heißt, dass die beiden Einträge zusammen
gelesen werden müssen und ein späterer Vorgang sie zusammenführen könnte.

**Folge:** DoD (2) greift — der Eintrag wird geschrieben, und er trägt die
Differenz oben als seine Existenzberechtigung.

## 3. Ausdrücklich NICHT in diesem Slice

- **Ein Sensor auf die Regel.** Ob eine Ersetzung ihre Frozen-Klassen
  aufgelistet hat, ist eine Aussage über einen **Akt**, nicht über einen
  ruhenden Zustand — dieselbe Lage wie bei
  [`AGENTS.md`](../../../../AGENTS.md) §3.6. Ein Gate darauf gäbe es nicht zu
  bauen.
- **Ein Retrofit der drei Anlässe.** Was geschehen ist, steht in den
  Evidence-Dateien; die geschlossenen Slices bleiben, wie sie sind.
- **Die Adoption der `v6.5.0`-Artefakt-Unterscheidung.** Sie liegt in
  [slice-208](../done/slice-208-v650-regel-adoption.md); dieser Slice
  **liest** sie nur, um zu entscheiden, ob seine Regel noch gebraucht wird.

## 4. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [x] **(1)** Die **Vorfrage ist beantwortet**: Macht die
      `v6.5.0`-Unterscheidung *einfrierend / lebend* die Regel entbehrlich? Ein
      **Nein** braucht die Differenz — was regelt der Kanon nicht —, ein **Ja**
      schließt den Slice ohne neuen Eintrag und setzt den Register-Ausgang auf
      *verkörpert* mit dem Kanon als Zielort.
- [x] **(2)** Fällt die Antwort auf **Nein**: Ein Konventions-Eintrag trägt die
      Regel samt der **Klassen-Liste über die Eigenschaft** und der Grenze, was
      sie nicht leistet.
- [x] **(3)** Der Registereintrag trägt den Ausgang *verkörpert* mit
      auflösbarem Zielort und Herkunfts-Anker `seit slice-209`.
- [x] `make gates` grün.
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §5 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 5. Abnahme-Punkte / Risiken

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Die Regel könnte überflüssig sein, und das wäre das beste Ergebnis.**
  `v6.5.0` führt die Unterscheidung *einfrierend / lebend* mit einer
  Aufzählung, die alle drei Anlässe abdeckt. Ein eigener Eintrag daneben wäre
  dann eine zweite Quelle für dieselbe Regel — genau das, wovor die
  Source-Precedence warnt. Die Vorfrage steht deshalb **als DoD-Punkt**, nicht
  als Vorbemerkung. — **Ausgang:** entfallen — die Vorfrage ist beantwortet,
  und zwar mit **Nein**. Die Regel ist nicht überflüssig: Der Kanon regelt die
  **Form eines Verweises** beim Schreiben, dieser Eintrag die
  **Ausschluss-Menge einer Massen-Operation** danach. Das Risiko kann nicht
  mehr eintreten, weil die Frage entschieden ist. **Was der Review daran
  korrigiert hat, ändert die Antwort nicht, sondern ihre Belege:** Von den drei
  Belegen trug einer die falsche Begründung und einer war zu stark formuliert;
  beide sind neu gefasst, der zweite als der schwächste ausgewiesen. Der Kern —
  anderer Adressat, anderer Zeitpunkt, andere Handlung — ist unangetastet
  geblieben.
- **Eine Regel aus drei Anlässen ist keine Inventur**
  ([`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md),
  7×). Drei ist die Kanon-Schwelle, aber die drei Anlässe sind alle vom selben
  Typ (Pin-/Format-Migration). Ob die Regel für andere mechanische Ersetzungen
  trägt, ist unbelegt und gehört als Grenze in den Eintrag. — **Ausgang:**
  weiter offen — eingetragen als achter Beleg bei
  [`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md).
  Die Grenze steht im Eintrag, wie der Plan es verlangt hat; belegt ist sie
  damit nicht. **Verschärfend, und erst im Schreiben aufgefallen:** Von den
  zwei Pflichten der Regel ist die **erste** — die Vorab-Liste — die unbelegte.
  Gefangen wurde der Fehler dreimal von der zweiten Pflicht oder vom Review,
  nie von einer Liste davor, denn es gab keine. Die Regel schreibt also eine
  Handlung vor, deren Wirksamkeit sie nicht belegen kann; das steht jetzt in
  ihr.
- **Kein Gate, und das ist eine Aussage über die Wirkung.** Die Regel hängt an
  der Disziplin des Ausführenden; sie verschiebt einen Fehler von *unsichtbar*
  nach *vermeidbar*, nicht nach *unmöglich*. Das gehört ausgeschrieben, sonst
  liest sie sich stärker, als sie ist. — **Ausgang:** entfallen — die Grenze
  ist ausgeschrieben, und zwar fünffach; der Eintrag sagt in eigenen Worten,
  dass er den Fehler von *unsichtbar* nach *vermeidbar* verschiebt, nicht nach
  *unmöglich*. **Die Liste war trotzdem unvollständig, und das ist die
  Pointe:** Der unabhängige Review fand eine sechste Grenze, die keine der fünf
  vorwegnahm — die **Kollision mit
  [`MR-013`](../../../../harness/conventions.md#mr-013)**, bei der die literale
  Befolgung der neuen Regel `make doc-check` rot gemacht hätte. Sie ist jetzt
  als eigener Block gemeldet statt still aufgelöst
  ([`AGENTS.md`](../../../../AGENTS.md) §1). Wer seine Grenzen aufzählt, hat
  sie damit nicht vollständig — das gilt auch für diesen Satz.

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-208](../done/slice-208-v650-regel-adoption.md)
liegt in `done/` — die Vorfrage aus DoD (1) braucht den adoptierten
`v6.5.0`-Stand, sonst urteilt sie über einen Kanon, den das Repo noch nicht
führt. WIP-Limit frei.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Zeigt die Vorfrage, dass die
  Klassen-Liste selbst eine Inventur über den Bestand verlangt, wird die
  Inventur ein eigener Slice.
- `in-progress` → `open` (blockiert): Ergibt die Vorfrage ein **Ja**, ist der
  Slice nicht blockiert, sondern **fertig** — er schließt mit DoD (1) und ohne
  DoD (2). Das ist kein Rückführungs-Fall, sondern der kürzere Weg.

**Closure-Trigger.** Zwei beobachtbare Kriterien und ein Lerneintrag: (a) die
Vorfrage ist beantwortet und die Antwort steht im Slice; (b) der
Registereintrag trägt einen Ausgang mit auflösbarem Zielort.

## 7. Vorgelagert (vor der Modus-Begründung)

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:268-269 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area: `*` (Repo-Default). Der Slice schreibt eine Regel über den
Umgang mit dem Doku-Bestand und ändert höchstens einen Konventions-Eintrag und
eine Register-Datei. `tools/harness/` ist **nicht** berührt: Der Slice baut
keinen Sensor, und §3 schließt das ausdrücklich aus.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:274-274 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, **37** Verzeichnisse über beide
Kürzel, `BEO-ALL` und `BEO-HARN`). **Fünf** Einträge sind einschlägig, und
zwei davon erreichen mit diesem Slice ihren Zielpunkt:

- [`mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
  (3×, Ausgang *geplant*) — **der Gegenstand selbst.** Der Eintrag nennt diesen
  Slice namentlich; sein Ausgang wird hier zu *verkörpert* oder, wenn die
  Vorfrage mit **Ja** endet, zu *verkörpert* mit dem **Kanon** als Zielort.
  Beides ist ein Ausgang, und keines der beiden ist ein Freitext.
- [`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
  (15× — der größte Eintrag des Registers) — **der einschlägigste für die
  Vorfrage.** DoD (1) fragt, ob eine Kanon-Stelle eine eigene Regel entbehrlich
  macht; das ist wörtlich die Frage *wie weit trägt ein zitierter Satz*. Der
  Ableiter gilt hier ungekürzt: den **Geltungsbereich** lesen, nicht den Titel,
  und die **direkteste** Quelle wählen. Ein **Ja** auf die Vorfrage, das den
  Kanon weiter zieht, als sein Absatz reicht, wäre der teuerste Fehler dieses
  Slice — er löschte eine Regel, statt sie zu ersetzen.
- [`registerzeile-ohne-ausgang-nach-schwelle`](../observations/BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle/observation.md)
  (2×, frisch erhöht) — der zweite Beleg entstand **einen Commit vor diesem
  Plan**: slice-208 hob einen Eintrag auf 3× und ließ ihn ohne Ausgang. Dieser
  Slice hat dieselbe Pflicht an zwei Stellen — DoD (3) und der Ausgang oben —
  und der Fehler ist frisch genug, um ihn nicht zu wiederholen.
- [`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md)
  (7×) — die Gegenrichtung, und §5 führt sie bereits als Risiko: Drei Anlässe
  sind die Kanon-Schwelle, aber alle drei sind vom selben Typ. Die Grenze
  gehört in den Eintrag, nicht in den Bericht danach.
- [`semantic-change-body-only-edges-stale`](../observations/BEO-ALL/semantic-change-body-only-edges-stale/observation.md)
  (12×) — der Geschwister-Eintrag zu
  [`MR-025`](../../../../harness/conventions.md#mr-025). Fällt die Vorfrage auf
  **Ja**, verschwindet nichts still: Wer eine Regel für entbehrlich erklärt,
  listet ihre Spiegel genauso auf wie beim Ändern.

**Geprüft und ausgeschlossen:**
[`zaehlmethode-misst-proxy-statt-gegenstand`](../observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/observation.md)
(3×, Ausgang *geplant* → [slice-210](../open/slice-210-zaehlmethode-vor-der-messung.md))
— dieser Slice zählt nichts; seine Vorfrage ist ein Urteil über einen Absatz,
keine Messung über eine Menge. **Keiner der fünf erreicht mit diesem Slice die
Schwelle erstmalig.**

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-07 gelesen: **beide Nachtläufe grün** —
`upstream-drift.yml` (jüngster Lauf 2026-09-07T05:33:45Z) und `image-scan.yml`
(2026-09-06T07:56:19Z). Damit ist die Meldung aufgelöst, die
[slice-208](../done/slice-208-v650-regel-adoption.md) noch als **ROT**
vorfand: Jener Lauf datierte von **vor** der Pin-Hebung, und die benannte
Grenze des Targets — es liest den jüngsten Lauf, nicht sein Alter — hat sich
als genau das erwiesen, was sie zu sein behauptete. Nichts zu tun.

## 8. Sub-Area-Modus-Begründung

**Modus:** `*` ist **GF** (Greenfield, Repo-Default).

- **Konventions-Dichte:** hoch. Der Umgang mit eingefrorenem Bestand ist über
  [`MR-025`](../../../../harness/conventions.md#mr-025) (Spiegel vor dem
  Editieren), [`MR-052`](../../../../harness/conventions.md#mr-052)
  (historisches Zitat über die git-Historie prüfbar) und
  [`MR-069`](../../../../harness/conventions.md#mr-069) (das Ventil als
  deklarierte Gate-Senkung) bereits dreifach berührt. **Genau diese Dichte ist
  das Risiko** — ein vierter Eintrag daneben muss sagen, was die drei nicht
  sagen, sonst ist er die zweite Quelle, vor der die Source-Precedence warnt.
- **Phase-Reife:** Phase 5. Ein Konventions-Eintrag aus einem 3×-Registerstand
  ist der eingespielteste Vorgang dieses Repos; die Form liegt in der vendorten
  Vorlage, der Ablauf in vier Vorgänger-Slices.
- **Evidenz-/Diskrepanz-Risiko:** **niedrig für den Bestand, hoch für die
  Vorfrage.** Am Bestand ist nichts zu inventarisieren — die drei Anlässe sind
  in den Evidence-Dateien belegt und werden nicht angefasst. Das Risiko sitzt
  allein im Urteil über die Kanon-Stelle, und es ist der 15×-Eintrag
  `citation-stretched-beyond-scope`, der es benennt.
- **Reconciliation-Aufwand:** keiner (GF). Graduation entfällt.
## 9. Closure-Notiz (nach `done/`)

**Geliefert.** Eine beantwortete Vorfrage, ein Konventions-Eintrag
([`MR-070`](../../../../harness/conventions.md#mr-070)) mit fünf
ausgeschriebenen Grenzen und einer gemeldeten Regel-Kollision, ein
Registereintrag auf *verkörpert* mit Zielort und Anker. Ein unabhängiger
Review, blockierend, sechs MEDIUM plus ein LOW und ein INFO — alle
eingearbeitet. `make gates` grün (zehn Gates, 706 Dateien).

**Was funktioniert hat: die Vorfrage als DoD-Punkt statt als Vorbemerkung.**
Der Plan hat *„macht der Kanon die Regel entbehrlich?"* nicht vorweg
beantwortet, sondern als abhakbaren Punkt geführt, mit beiden Ausgängen
vorab beschrieben. Das hat den Slice gezwungen, den Kanon-Absatz **vollständig
zu lesen**, statt seinen Titel zu zitieren — und der Review konnte gegen einen
benannten Anspruch prüfen statt gegen eine Stimmung. Von den sechs MEDIUM
betrafen vier die **Belege** der Antwort und keiner ihr **Ergebnis**; das ist
die Trennung, die ein DoD-Punkt möglich macht.

**Was Friktion war: der Slice ist genau der Klasse aufgesessen, die er vorab
gesichtet hatte.** §7 führt
[`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
(15×) als *„den einschlägigsten für die Vorfrage"* — und drei der sechs
Befunde sind Ausprägungen genau davon. Der schwerste: Der tragende Beleg las
den Kanon-Satz *„Tag und Pfad in Inline-Code, **nicht als Link**"* als erfüllt,
weil die geprüften Direktiven keine Links sind. Das ist die **Negation gelesen
statt der Vorschrift** — ein HTML-Kommentar ist weder das eine noch das andere.
**Zum zweiten Mal in Folge gilt damit:** Das Sichten benennt die Klasse und
verhindert ihr Eintreten nicht; gefangen hat sie der unabhängige Review.

**Steering-Loop-Lerneintrag:** *Die Entscheidung trägt, die Begründung daneben
nicht*
([`begruendung-traegt-entscheidung-nicht`](../observations/BEO-ALL/begruendung-traegt-entscheidung-nicht/observation.md),
jetzt 2×) — und dieser Beleg ist aus zwei Richtungen zugleich entstanden. **Der
Slice** kam zum richtigen Schluss mit einer falschen Begründung. **Der Review**
widerlegte die Begründung zu Recht und stellte daneben eine eigene auf, die
ebenfalls nicht zutrifft: die Kanon-Form trage den Präfix
`.harness/baseline/<tag>/` „gerade nicht". Gemessen tragen ihn **neun**
Vorkommen in eingefrorenen Artefakten, die nur überlebt haben, weil ihr
Verzeichnis ausgenommen war. **Ein Befund kann richtig sein und seine
Begründung falsch** — wer den Befund übernimmt, ohne nachzumessen, übernimmt
beides. Der Eintrag steht bei 2× und wartet.

**Der teuerste Einzelbefund war nicht die Vorfrage, sondern die Regel selbst.**
[`MR-070`](../../../../harness/conventions.md#mr-070)s Geltungsbereich führte
*„`git mv` mit Nachzug"*, während seine
Tabelle `done/` ausnimmt — und
[`MR-013`](../../../../harness/conventions.md#mr-013) verlangt genau diesen
Nachzug in `done/`-Dateien. **Der Beanspruchungs-Commit dieses Slice hat ihn an
sechs Stellen ausgeführt**, drei Commits bevor die Regel geschrieben wurde, die
ihn verbietet. Wer sie literal befolgt hätte, machte `make doc-check` rot. Die
Kollision ist jetzt gemeldet, nicht still aufgelöst
([`AGENTS.md`](../../../../AGENTS.md) §1): Der Eigenschafts-Test entscheidet
sie, weil ein Pfad-Nachzug die Aussage nicht verfälscht, eine
Kennungs-Ersetzung schon.

**Die Lehre daraus ist unbequemer als der Fix.** Der Eintrag zählt fünf
Grenzen auf, und keine davon nahm die sechste vorweg. Eine ausgeschriebene
Grenzen-Liste macht eine Regel ehrlicher, nicht vollständig — und die Liste
selbst braucht einen fremden Leser.

**Was offen bleibt.** Die erste der beiden Pflichten von
[`MR-070`](../../../../harness/conventions.md#mr-070) — die
Vorab-Liste — ist die **unbelegte**: Gefangen wurde der Fehler dreimal von der
zweiten Pflicht oder vom Review, nie von einer Liste davor, denn es gab keine.
Das steht im Eintrag und im Register
([`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md),
achter Beleg). Ob die Regel trägt, entscheidet die vierte Instanz — dieselbe
Wette wie bei [`MR-066`](../../../../harness/conventions.md#mr-066).

**Die drei Paarungen, gemessen.** **(a) Anker** — [`MR-070`](../../../../harness/conventions.md#mr-070) trägt
`seit slice-209`, der Registereintrag zeigt auf ihn, der Zielort löst auf.
**(b) Folge-Slice** — dieser Slice nennt keinen; die Kanon-Beobachtung aus dem
Beleg (die vom Kanon vorgeschriebene Form ist gegen `sed` unsicher) ist als
**ausgehender CR** vorgemerkt und bewusst ohne Kennung gelassen, weil §3 den
anderen Vorgang ausschließt — eine Kennung ohne Datei wäre eine Adresse, die
die Sendung nicht annimmt. **(c) Register** — alle zitierten
Beobachtungs-Pfade lösen auf; die drei neuen Belege liegen als
`evidence/slice-209.md` in ihren Verzeichnissen. Der Wachposten
[`kanal-kennung-als-inhalt-gelesen`](../observations/BEO-ALL/kanal-kennung-als-inhalt-gelesen/observation.md)
trägt weiterhin kein `evidence/` — unverändert die benannte Spannung aus
slice-208, nicht ein Defekt dieses Slice.
