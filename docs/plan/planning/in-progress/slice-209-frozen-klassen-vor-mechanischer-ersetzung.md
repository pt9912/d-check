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
Beim ersten Mal (Register-Formatmigration) blieben `Accepted`-ADR-Kerne und
gesendete CRs unbedacht, obwohl die drei *benannten* Frozen-Verzeichnisse
korrekt ausgenommen waren. Beim zweiten (Pin-Hebung) dasselbe Muster. Beim
dritten (Pin-Hebung) gab es **gar keine** Ausnahme-Liste, und gehoben wurden
zehn `d-check:cite`-Direktiven in `done/`-Slices, deren Zitat gegen den alten
Tag geschrieben war. **Die Verzeichnis-Liste ist die falsche Abstraktion** —
maßgeblich ist die Eigenschaft *„zitiert den Stand seiner Zeit"*, und die
tragen auch Zeilen, die in keinem der drei Verzeichnisse stehen.

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

**Beleg 1 — eine kanon-konforme Datei wurde trotzdem beschädigt.** Die dritte
Form des Kanons lautet: *„Eine Stelle der vendored Baseline heißt Tag **und**
Pfad in Inline-Code, nicht als Link."* Genau diese Form haben die zehn
`d-check:cite`-Direktiven, die slice-207 mitgehoben hat — sie sind
HTML-Kommentare mit Tag und Pfad, kein Link. Sie waren **konform** und wurden
beschädigt. Damit ist gezeigt: Vollständige Kanon-Befolgung schließt den
Fehler nicht aus, sie ist gegen ihn wirkungslos. Schlimmer noch, die Richtung
kehrt sich um — die vom Kanon **vorgeschriebene** Form ist genau das, was eine
Tag-Ersetzung greift, während der abgeratene Link von einem Sensor gemeldet
würde.

**Beleg 2 — eine der drei Instanzen liegt außerhalb des Kanon-Geltungsbereichs.**
Der Fund in slice-202 betraf eine identifizierende Nennung in einem
**lebenden** Dokument (die Plan-Tabelle, die den zu entfernenden Baum
benennt). Der Kanon nimmt lebende Artefakte ausdrücklich aus. Eine Regel, die
nur die einfrierenden deckt, hätte diese Instanz nie gefangen.

**Beleg 3 — die Kanon-Aufzählung ist für dieses Repo unvollständig.** slice-195
traf `Accepted`-ADR-Kerne (vom Kanon genannt) **und gesendete CRs** unter
`docs/plan/cr/` (vom Kanon nicht genannt). Ein gesendeter CR ist einfrierend
nach derselben Eigenschaft — er hält eine Bitte zu ihrem Datum fest —, steht
aber in keiner der fünf Klassen. Das ist wörtlich der Punkt dieser Regel: **die
Liste über die Eigenschaft, nicht über eine Aufzählung.**

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
- [ ] **(2)** Fällt die Antwort auf **Nein**: Ein Konventions-Eintrag trägt die
      Regel samt der **Klassen-Liste über die Eigenschaft** und der Grenze, was
      sie nicht leistet.
- [ ] **(3)** Der Registereintrag trägt den Ausgang *verkörpert* mit
      auflösbarem Zielort und Herkunfts-Anker `seit slice-209`.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §5 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 5. Abnahme-Punkte / Risiken

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Die Regel könnte überflüssig sein, und das wäre das beste Ergebnis.**
  `v6.5.0` führt die Unterscheidung *einfrierend / lebend* mit einer
  Aufzählung, die alle drei Anlässe abdeckt. Ein eigener Eintrag daneben wäre
  dann eine zweite Quelle für dieselbe Regel — genau das, wovor die
  Source-Precedence warnt. Die Vorfrage steht deshalb **als DoD-Punkt**, nicht
  als Vorbemerkung. — **Ausgang:** \<offen\>
- **Eine Regel aus drei Anlässen ist keine Inventur**
  ([`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md),
  7×). Drei ist die Kanon-Schwelle, aber die drei Anlässe sind alle vom selben
  Typ (Pin-/Format-Migration). Ob die Regel für andere mechanische Ersetzungen
  trägt, ist unbelegt und gehört als Grenze in den Eintrag. — **Ausgang:** \<offen\>
- **Kein Gate, und das ist eine Aussage über die Wirkung.** Die Regel hängt an
  der Disziplin des Ausführenden; sie verschiebt einen Fehler von *unsichtbar*
  nach *vermeidbar*, nicht nach *unmöglich*. Das gehört ausgeschrieben, sonst
  liest sie sich stärker, als sie ist. — **Ausgang:** \<offen\>

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

\<wird vor dem `git mv` nach `done/` gefüllt\>
