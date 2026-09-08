# Slice slice-216: Entscheid über den eingehenden `links`/`anchors`-CR

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [eingehender CR *`links` und `anchors` brauchen `exempt-paths`*](../../cr/2026-09-07-cr-eingehend-adopter-links-anchors-exempt-paths.md)
(abgelegt, unentschieden — Auftraggeber-Entscheid 2026-09-08: ohne Rückfrage
entscheiden), [`DC-FA-REF-001`](../../../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)
(das vorhandene Ventil), [`DC-FA-LINK-001`](../../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links),
[`DC-FA-ANCH-001`](../../../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors),
[`MR-069`](../../../../harness/conventions.md#mr-069) (das Ventil als
deklarierte Gate-Senkung).

**Berührte Spec-Stellen:** — *(voraussichtlich keine. Fällt der Entscheid auf
„umsetzen", ist die Anforderung ein Folge-Slice; dieser trägt den Entscheid.)*

**Verantwortlich:** pt9912 · **Autor:** pt9912. **Datum:** 2026-09-08.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Den eingehenden CR **entscheiden** — er liegt seit 2026-09-07
abgelegt und unentschieden. Der Entscheid ruht auf **einer** Messung, die ihn
trägt oder kippt: **Ist die erbetene Fähigkeit heute ausdrückbar?**
[`DC-FA-REF-001`](../../../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)
nennt `links` und `anchors` ausdrücklich und trägt mit `in:` einen
Quell-Skopus; ein Eintrag `in: <datei>` mit `refs: ["**"]` wäre **genau**
das erbetene datei-weite Ventil. **Ob das so funktioniert, ist gemessen zu
beantworten, nicht aus dem Schema zu erschließen.**

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Umsetzung, falls der Entscheid dafür ausfällt.** Eine neue
  Konfigurationsfläche an zwei Kern-Modulen ist eine Anforderung im Lastenheft
  mit Akzeptanzkriterien und Modul-Arbeit — ein Folge-Slice. Dieser trägt den
  **Entscheid**, wie es die beiden Vorgänger-CRs vorgemacht haben.
- **Eine Antwort an den Absender formulieren.** Der Zwischenbescheid ist
  geschrieben und liegt beim Auftraggeber; ob und wann er geht, ist dessen
  Sache und kein Liefer-Punkt hier.
- **Der Bestand hinter dem Ventil.** Die 25 Baseline-Einträge und die sieben
  toten `in:`-Skopen bleiben, wie sie sind
  ([`MR-069`](../../../../harness/conventions.md#mr-069) führt sie); das
  Aufräumen ist ein eigener Vorgang und war es schon vor diesem CR.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [x] **(1)** **Gemessen**, ob `ignore-refs` mit `in: <datei>` und
      `refs: ["**"]` die erbetene datei-weite Wirkung erzielt — mit echter
      Ausgabe, in beide Richtungen: der gedeckte Verweis schweigt, ein
      **nicht** gedeckter in derselben Datei meldet weiter. Fällt die Messung
      negativ aus, kippt der Entscheid.
- [x] **(2)** Der CR ist **entschieden**: Bitte beantwortet, Begründung je
      Argument des Absenders (drei: der Knopf, die sechs Module, die drei
      verworfenen Wege), `Stand:`-Zeile gesetzt.
- [x] **(3)** Der Entscheid trägt **Umkehr-Bedingungen**, wie der
      Zeilenlängen-CR sie führt — beobachtbar formuliert, mit der Grenze, dass
      kein Sensor über sie wacht.
- [x] `make gates` grün.
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

**Zuerst die Form des Gegenstands** ([`AGENTS.md`](../../../../AGENTS.md) §5).

**Was heißt *ausdrückbar*?** Nicht *„das Schema erlaubt es"*, sondern: Ein
Konfigurations-Eintrag erzeugt **dieselbe Wirkung** wie die erbetene Fähigkeit,
gemessen an einem Bruch-Test in **beide** Richtungen. Die Gegenrichtung ist
hier die wichtigere — ein Ventil, das **alles** verschluckt, hat die erbetene
Wirkung und wäre trotzdem die schlechtere Antwort.

**Was zählt als *Argument des Absenders*?** Der CR führt drei, und jedes
bekommt eine eigene Antwort: **(a)** *„kein Knopf"* — die Prämisse; **(b)**
*„sechs Module führen ihn"* — das Konsistenz-Argument; **(c)** die drei
verworfenen Wege (`ignore-refs`-Breite, `scan.ignore`-Kosten, Koexistenz des
alten Baums). **Eine Sammelantwort trägt nicht** — der CR ist sorgfältig
gebaut, und (b) steht unabhängig von (a).

**Was bereits gemessen ist** und im CR-Dokument steht, damit dieser Slice es
nicht zweimal misst: [`DC-FA-REF-001`](../../../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)
nennt `links` und `anchors` ausdrücklich
und unterdrückt `target-missing` **und** `anchor-missing`; dieses Repo fährt
25 Einträge; ein toter Verweis **ohne** Baseline-Bezug in einer gedeckten Datei
meldet weiter (slice-208); und die sechs Module mit `exempt-paths` sind
nachgezählt.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [der CR](../../cr/2026-09-07-cr-eingehend-adopter-links-anchors-exempt-paths.md) | update | Messung, Entscheid, Umkehr-Bedingungen |

## 4. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei; der Auftraggeber-Entscheid
vom 2026-09-08 liegt vor (*„jetzt entscheiden, ohne die Rückfrage"*).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Ergibt DoD (1), dass die Fähigkeit **nicht**
  ausdrückbar ist, trägt der Entscheid eine Produkt-Frage statt einer
  Konfigurations-Antwort — dann ist die Anforderung zu schneiden und dieser
  Slice zu klein.
- `in-progress` → `open` (blockiert): Zeigt sich, dass der Entscheid ohne die
  Antwort des Absenders auf **(c)** nicht zu treffen ist — etwa weil sein
  Breiten-Wächter eine Eigenschaft des Werkzeugs und nicht seines Repos ist —,
  ruht der Slice bis zur Rückmeldung.

## 5. Closure-Trigger

Zwei beobachtbare Kriterien und ein Lerneintrag: (a) die Messung aus DoD (1)
steht mit echter Ausgabe im CR, in beide Richtungen; (b) der CR trägt eine
`Stand:`-Zeile mit Entscheid, drei Antworten und Umkehr-Bedingungen, und
`make gates` ist grün.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Der bequeme Entscheid ist „nein", und er wäre halb begründet.** Die
  Prämisse des CR ist widerlegt — daraus folgt aber **nicht**, dass sein
  Konsistenz-Argument falsch ist. `links` und `anchors` sind tatsächlich die
  einzigen referenz-prüfenden Module mit **einer** Ventil-Achse. Wer (a)
  widerlegt und (b) damit erledigt glaubt, hat zwei Argumente in eines
  gefaltet. — **Ausgang:** eingetreten, in anderer Form als beschrieben.
  Gefaltet wurde **nicht** — der Entscheid antwortete (b) getrennt und
  schrieb ausdrücklich, es werde *nicht* mit (a) erledigt. Verfehlt hat
  ihn trotzdem, wovor das Risiko schützen sollte: (b) bekam keine
  tragende Antwort. Die Zahl war die **fremde** (sechs statt acht), und
  der Grund — *„die schwächere ist die bequemere"* — ist am eigenen
  Bestand widerlegt (`codepaths` nimmt `docs/reviews/**` über genau die
  Datei-Achse aus). Beides in Review-Runde 1 gefunden, der Grund
  **gestrichen** statt abgeschwächt; was (b) jetzt trägt, ist schmaler.
- **Die Messung könnte die Bitte bestätigen statt sie zu erübrigen.** Ergibt
  DoD (1), dass `refs: ["**"]` **nicht** wie `exempt-paths` wirkt, ist die
  Fähigkeit heute nicht ausdrückbar und der CR im Kern berechtigt. **Der Slice
  ist so geschnitten, dass er beide Ausgänge trägt** — die Messung steht vor
  dem Entscheid, nicht hinter ihm. — **Ausgang:** eingetreten, zur Hälfte.
  Für den **beschriebenen** Fall erübrigt die Messung die Bitte: Form **B**
  löst ihn, auf beiden Achsen gemessen. Für den **allgemeinen** Fall
  bestätigt sie ihn — die Symlink-Ablehnung überlebt jedes `refs`-Glob, die
  Ventile stehen quer zueinander statt in Dominanz. Der Entscheid kippt
  nicht, seine Reichweite ist die des beschriebenen Falls; die Bitte bleibt
  für die Symlink-Klasse unbeantwortet und ist als Umkehr-Bedingung (1)
  adressiert.
- **Ein Entscheid ohne die Antwort des Absenders bleibt einseitig.** Argument
  (c) stützt sich auf einen Breiten-Wächter in **seinem** Repo; ob der eine
  bewusste Regel oder ein Nebeneffekt ist, weiß nur er. Der Entscheid muss das
  benennen, statt es zu übergehen — sonst beantwortet er ein Argument, das er
  nicht gehört hat. — **Ausgang:** eingetreten, aber **nicht** in der
  Form, die das Risiko beschrieb. Der Entscheid ist ohne die Antwort
  gefallen — auf Auftraggeber-Entscheid, und das steht im Dokument statt
  verschwiegen zu werden. **Einseitig blieb er trotzdem**, und die
  Adresse dafür sind die drei Umkehr-Bedingungen: Bedingung (2) nennt
  genau diesen Punkt — erwiese sich der Breiten-Wächter als Eigenschaft
  des Werkzeugs statt des fremden Repos, fiele die Rechnung anders aus.

## 7. Closure-Notiz

**Geliefert.** Ein **entschiedener eingehender CR** — der zweite, den dieses
Repo ablehnt, und der erste von einem Adopter statt vom Auftraggeber. `links`
und `anchors` bekommen kein `exempt-paths`. Vor dem Entscheid steht die
Messung: drei Konfigurationen über einen Mini-Bestand, dazu eine zweite
Messung auf der `anchors`-Achse. Ein unabhängiger Review, blockierend, sieben
MEDIUM und zwei LOW. `make gates` grün (zehn Gates, 744 Dateien).

**Was funktioniert hat: die Messung stand vor dem Entscheid, nicht dahinter.**
§3 hatte den Slice ausdrücklich so geschnitten, dass er **beide** Ausgänge
trägt — fiele die Messung negativ aus, wäre der CR im Kern berechtigt gewesen.
Sie fiel positiv aus, und der Bruch-Test trug die Gegenrichtung mit: ein
**nicht** gedeckter toter Verweis in derselben Datei meldet weiter. Ohne diese
zweite Hälfte wäre Form **B** — das Rezept, das an die Stelle des Nein tritt —
gar nicht sichtbar geworden.

**Was Friktion war: die Begründung reichte weiter als ihr Beleg, an drei
Stellen.** Der Entscheid kippt dadurch nicht; seine **Reichweite** schon.

**Erstens die Äquivalenz.** Ich schrieb, `refs: ["**"]` erreiche die Wirkung
von `exempt-paths` *„exakt"* und *„und mehr"*. Gemessen war: zwei Verweise
einer Datei werden stumm. Die **Symlink**-Ablehnung überlebt jedes `refs`-Glob
— und das Produkt sagt es selbst, im Kommentar über dem Ventil in
`internal/hexagon/core/rules/links.go`: die Symlink-Prüfung habe **Vorrang**
und bleibe unberührt. Die beiden Ventile stehen **quer** zueinander, nicht in
Dominanz. **Und dieselbe Überdehnung eine Ebene tiefer:** gemessen mit *einem*
Modul, entschieden über *zwei*.

**Zweitens die fremde Zahl.** Der CR schrieb *„sechs Module"*, ich übernahm
die Sechs und argumentierte auf ihr — obwohl die gezählte Menge **unsere**
ist. Es sind **acht**. Die richtige Zahl stand im eigenen Text: Zwei Sätze
weiter nannte ich `ids` als Beispiel für *„eine zweite Ventil-Achse"*, ohne zu
bemerken, dass diese Achse genau `exempt-paths` **ist**.

**Drittens ein Grund, der am eigenen Bestand widerlegt ist.** Ich
argumentierte, die Datei-Achse sei *„die bequemere"* und verdränge die
schärfere. Dieses Repo nimmt `docs/reviews/**` für `codepaths` über die
**Datei**-Achse aus — obwohl dort alle drei Achsen zur Verfügung stehen, und
für **dieselbe** Artefakt-Klasse, um die der CR bittet. Wer das tut und dem
Adopter die Datei-Achse mit dem Bequemlichkeits-Argument verweigert, misst mit
zwei Maßen. **Der Grund ist gestrichen**, nicht abgeschwächt.

**Steering-Loop-Lerneintrag: eine Deckungs-Messung ist keine Äquivalenz.**
Der Eintrag
[`commit-message-overclaims-work`](../observations/BEO-ALL/commit-message-overclaims-work/observation.md)
(12×) trägt den Ableiter — *„suche die N+1-te Form"* —, seit welle-82 in
[`AGENTS.md`](../../../../AGENTS.md) §5 verkörpert. **§8 hatte ihn für diesen
Slice sogar ausdrücklich als den wichtigsten benannt** und geschrieben, der
Bruch-Test müsse die Gegenrichtung mitnehmen. Er nahm sie mit — **innerhalb
der einen Auflösungs-Klasse**, die ich im Blick hatte. Die N+1-te Form lag
nicht in der Gegenrichtung, sondern in einer **zweiten Klasse**, an die die
Vorprüfung nicht dachte. *Die Gegenrichtung zu messen ist nicht dasselbe wie
die Menge zu kennen, über die man aussagt.*

**Ein zweiter Lerneintrag, und er ist der billigere.** Bei einer fremden Zahl
über die **eigene** Menge ist die mitgelieferte Zählung ein Angebot, keine
Inventur — ein Grep-Muster, das jemand anders geschrieben hat, ist seine Sicht
auf unseren Bestand.
[`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
(17×) trägt jetzt diese Variante: nicht der fremde Bestand, sondern die eigene
Menge durch ein fremdes Muster gesehen.

**Und die Umkehr-Bedingungen haben es zum siebten Mal nicht selbst gesehen.**
Bedingung (1) verlangte *„einen Fall, den `ignore-refs` nicht ausdrücken
kann"* — der Fall stand beim Schreiben **im selben Dokument**, im Abschnitt
darüber. Genau der gemessene Mechanismus von
[`grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md)
(7×): die fehlende Zeile in der anderen Rahmung, in derselben Datei.
Korrigiert nennt sie jetzt, was ihr wirklich fehlt — nicht der Fall, sondern
eine **Population** im Bestand des Absenders.

**Was offen bleibt.** Zwei Dinge, beide benannt statt behoben. **(1)** Eine
**Spec/Produkt-Divergenz**, älter als dieser CR:
[`DC-FA-REF-001`](../../../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)
sagt, hinter dem Ventil blieben *„Symlink-Ablehnung, Repo-Escape u. Ä."*
stehen — gemessen bleibt die Symlink-Ablehnung, der Repo-Escape nicht. Das ist
kein Befund dieses CR und gehört als eigener Vorgang aufgenommen. **(2)** Die
**Rückfrage an den Absender** ist unbeantwortet: ob nach Form **B** ein Fall
übrig bleibt. Der Entscheid ist ohne sie gefallen, auf Auftraggeber-Entscheid,
und das steht im Dokument, statt verschwiegen zu werden.

**Die drei Paarungen, gemessen.** **(a) Anker** — vakant: Der Slice verkörpert
keine Steering-Loop-Regel; seine drei Lerneinträge liegen bei bestehenden.
**(b) Folge-Slice** — keiner genannt; die Spec/Produkt-Divergenz ist bewusst
**ohne** Kennung gelassen, wie die Prerelease-Frage in slice-215. **(c)
Register** — alle zitierten Pfade lösen auf, die drei neuen Belege liegen als
`evidence/slice-216.md` in ihren Verzeichnissen. Der Wachposten
[`kanal-kennung-als-inhalt-gelesen`](../observations/BEO-ALL/kanal-kennung-als-inhalt-gelesen/observation.md)
trägt weiterhin kein `evidence/` — unverändert die benannte Spannung aus
slice-208.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende. Dieses Repo führt **drei** Prüfungen — die
zwei kanonischen und, als Adaption, den Nachtlauf-Stand
([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:268-269 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area: `*` (Repo-Default). Der Slice ändert **ein** getracktes
Artefakt, das CR-Dokument. **Der Gegenstand der Messung ist zwar das Produkt**
(die Module `links` und `anchors`), aber es wird **gefahren, nicht geändert** —
und für Produkt-Code führt dieses Repo keine eigene Sub-Area; die
Modus-Deklaration kennt `*` und `tools/harness/`. Eine dritte hier zu erfinden
wäre eine Sub-Area ohne Deklaration.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:274-274 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, **38** Verzeichnisse über beide
Kürzel; `BEO-HARN` einzeln geöffnet — dort steht der eine offene Eintrag zu
`--check-latest`, der diesen Slice **nicht** berührt). **Vier** Einträge sind
einschlägig:

- [`commit-message-overclaims-work`](../observations/BEO-ALL/commit-message-overclaims-work/observation.md)
  (11×, verkörpert) — **der frischeste und für DoD (1) der wichtigste.** Sein
  Ableiter, gerade in slice-215 teuer gelernt: *„suche die N+1-te Form"*. Eine
  Messung, die zeigt, dass `refs: ["**"]` **einen** Verweis verschluckt, sagt
  nichts über **alle**. Der Bruch-Test muss die Gegenrichtung mitnehmen — ein
  nicht gedeckter Verweis in derselben Datei —, sonst ist der Entscheid auf
  eine Form gebaut.
- [`citation-stretched-beyond-scope`](../observations/BEO-ALL/citation-stretched-beyond-scope/observation.md)
  (18×) — für die Antwort auf Argument (b): Das Konsistenz-Argument des
  Absenders steht **unabhängig** von seiner widerlegten Prämisse. Zwei
  Argumente in eines zu falten, um beide mit einem Satz zu erledigen, ist
  genau diese Klasse; §6 führt es als erstes Risiko.
- [`eigene-menge-gemessen-fremde-behauptet`](../observations/BEO-ALL/eigene-menge-gemessen-fremde-behauptet/observation.md)
  (16×, Stand *gemischt*) — Argument (c) stützt sich auf einen Breiten-Wächter
  im **fremden** Repo. Was dort gilt, ist von hier nicht gemessen; der
  Entscheid darf darüber nicht urteilen, sondern nur benennen.
- [`grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md)
  (6×, verkörpert) — der Entscheid trägt Umkehr-Bedingungen, also eine Liste.
  Ihr Ableiter gilt: den Vertrags-Teil umdrehen und, wo der Gegenstand Code
  ist, gegen den Code prüfen.

**Keiner der vier erreicht mit diesem Slice die Schwelle erstmalig.**

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-08 gelesen: **beide Nachtläufe grün** —
`upstream-drift.yml` (jüngster Lauf 2026-09-07T05:33:45Z) und `image-scan.yml`
(2026-09-07T08:21:32Z). Für diesen Slice ohne Bezug: Er ändert kein
gepinntes Artefakt und keinen Sensor. **Notiert, weil die Prüfung unbedingt
ist.**

**Modus-Begründungsblock.** Alle berührten Sub-Areas GF — ein Block genügt.

### Sub-Area: `*`

- **Modus:** GF (Repo-Default).
- **Konventions-Dichte:** hoch für die **Form** eines CR-Entscheids — zwei
  Vorgänger haben sie vorgelebt ([`MR-035`](../../../../harness/conventions.md#mr-035),
  [`MR-036`](../../../../harness/conventions.md#mr-036) regeln die ausgehende
  Richtung; die eingehende ist die dritte Klasse ohne deklarierten Ort, und
  das steht in jedem der drei Dokumente).
- **Phase-Reife:** Phase 5 für den Vorgang (zwei entschiedene CRs), Phase 4
  für den Gegenstand: Ein Entscheid **gegen** eine Bitte ist erst einmal
  gefallen (Zeilenlänge), und dort war die Bitte des Auftraggebers — hier ist
  es die eines Adopters.
- **Evidenz-/Diskrepanz-Risiko:** **mittel.** Am Bestand ist wenig zu
  inventarisieren; das Risiko sitzt in der Messung (eine Form statt aller) und
  in der Versuchung, drei Argumente mit einem Satz zu erledigen.
- **Reconciliation-Aufwand:** keiner (GF).
