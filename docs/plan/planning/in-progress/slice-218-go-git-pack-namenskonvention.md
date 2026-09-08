# Slice slice-218: Die go-git-Pack-Namenskonvention als Grenze der git-lesenden Module

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle. Es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden wäre (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
und [`DC-FA-COMMITS-001`](../../../../spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in)
— **als die zwei Zusagen, deren Reichweite hier benannt wird**, nicht als
geänderte Anforderungen. Keine ADR: Der Slice trifft keine Entscheidung, er
schreibt eine gemessene Eigenschaft auf.

**Berührte Spec-Stellen:** — (der Slice berührt keine Spec-Stelle).

**Verantwortlich:** pt9912 (Implementer-Rolle).

**Autor:** pt9912.

## 1. Ziel und Abgrenzung

**Ziel.** Eine gemessene Produkt-Grenze bekommt einen Ort: **Die git-lesenden
Module finden Objekte nur in Pack-Dateien, die git's kanonischen Namen tragen**
(`pack-<Hash-des-Packs>.{idx,pack}`). Weicht der Name ab, ist der Pack
unsichtbar, und Objekte, die nur dort liegen, sind *nicht gefunden*.

**Der Anlass war ein realer Ausfall, kein Gedankenspiel.** Während slice-217
brach der `pre-commit`-Hook mit `HEAD-Tree nicht lesbar: object not found` ab,
nachdem `git maintenance` Packs mit `loose-`-Präfix angelegt hatte —
`git repack -A -d` behob es. **Das trifft jeden Adopter**, der `vcs` oder
`commits` fährt und dessen git eine Wartungs-Aufgabe laufen lässt; es ist
keine Eigenschaft dieses Repos.

**Sechs Proben, isoliert gefahren** (Probe-Repo, alle Objekte im Pack, `.idx`,
`.pack` und `.rev` gemeinsam umbenannt, sonst nichts verändert):

| Pack-Name | Ergebnis |
|---|---|
| `pack-<eigener Hash>` (Kontrolle) | `2 Datei(en) geprüft, 0 Befund(e)`, Exit **0** |
| `loose-<eigener Hash>` | `staged-Basis "HEAD" nicht auflösbar`, Exit **2** |
| `pack-0123…4567` (gültige Form, **fremder** Hash) | Exit **2** |
| `pack-zzzzzzzz` | Exit **2** |
| `packXYZ` | Exit **2** |
| `xpack-abc` | Exit **2** |

**Damit zählen beide Namensteile**, und das ist mehr, als der Anlass zeigte:
Nicht nur das Präfix (`loose-` scheitert mit korrektem Hash), sondern auch der
Hash (die gültige Form scheitert mit fremdem Hash). **`git` selbst liest
weiter** — es enumeriert `*.idx` unabhängig vom Namen; nur go-git tut es
nicht.

**Abgrenzung — vier Punkte, jeder mit Grund:**

1. **Kein Repack durch das Produkt, und der Grund ist ein Kernvertrag.** Die
   naheliegendste Abhilfe wäre, dass d-check selbst repackt — das schriebe in
   das geprüfte Repository und bräche
   [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).
   *Es wäre ein anderer Vorgang*, und zwar einer mit ADR-Last.
2. **Kein Sensor auf Pack-Namen.** Ein Wächter darüber wäre ein neues Modul
   mit eigener Scan-Zusage. *Bestand bleibt bewusst stehen.*
3. **Keine Lastenheft-Änderung — die Spezifikation dagegen schon.**
   [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
   verspricht schon heute, dass eine Core-Änderung gemeldet wird; der Fix löst
   die Zusage ein, statt eine neue zu vereinbaren. **Die erste Fassung dieses
   Punktes schloss auch die Spezifikation aus** — mit der Begründung, ein
   Spec-Eintrag käme erst bei einem **neuen Grund-Code** in Frage. Das ist am
   eigenen Bestand widerlegt: Die Historie derselben Spec-Sektion führt zum
   2026-08-31 einen strukturgleichen Fall (Rename-Erkennung, *„still war
   ausgerechnet der Modus, den die CI fährt"*), der **ohne** neuen Grund-Code
   einen Spezifikations-Eintrag bekam und ausdrücklich keinen Lastenheft-Bump.
   §[`DC-FA-VCS-001.a`](../../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) Schritt 2 trägt den Mechanismus jetzt als Geschwister
   zur Rename-Zusage.
4. **Zwei Dinge bleiben ungelöst, beide bewusst und beide in
   [`CO-001`](../../carveouts/CO-001-vcs-range-stiller-skip.md) geführt.**
   Erstens wirkt der Fix erst mit dem nächsten Release — jeder Adopter auf
   `v0.75.0` oder früher behält die zwei geschlossenen Pfade blind
   (*Folge-Slice*: [`slice-219`](../open/slice-219-release-loest-co-001.md)).
   Zweitens ist die **dritte** Ausprägung gar nicht behoben und trifft auch
   dieses Repo (*Folge-Slice*:
   [`slice-220`](../open/slice-220-vcs-pfadmenge-statt-diff.md)). Keines von
   beiden ist stillschweigend hingenommen.

**Plan-Änderung (2026-09-08, nach Review-Runde 1) — der Fix kommt hinzu.**
Die erste Fassung schloss **Produkt-Code aus** und wollte den stillen Pfad nur
beschreiben. Der Review fand, dass die beschriebene Grenze in Wahrheit ein
**Defekt** ist: `vcs` überspringt im `RANGE=`-Modus eine echte
`core-drift-vcs`-Verletzung, wenn der unsichtbare Pack einzelne Objekte
verschluckt — gemessen, Exit 0 statt Exit 1. **Auftraggeber-Entscheid vom
2026-09-08: der Fix wird mitgenommen**, statt ihn zu vertagen; ein Gate, der
im CI-Modus schweigt, ist kein Gate (Baseline-Regelwerk
`modul-13-quality-gates.md` §Kernidee). Die Abgrenzung ist damit
**ausdrücklich geöffnet** — nicht still ausgeweitet, und das steht hier vor
dem Code.

**Der Slice überschreitet damit die Ein-Sitzungs-Review-Grenze**
([`MR-066`](../../../../harness/conventions.md#mr-066)), und beide Pflichten
sind hier einzulösen:

- **Grund für die Nicht-Rückführung:** Eine Teilung träfe genau die Naht, an
  der der Wert entsteht. Die korrigierte Doku und der Fix beschreiben
  **denselben** Sachverhalt in zwei Registern; getrennt geschnitten stünde die
  Doku eine Zeit lang als *„so verhält es sich"* da, während der Fix sie schon
  überholt — und die Beschreibung eines Defekts, den es nicht mehr gibt, ist
  schlechter als keine. Zudem ist der Fix klein (**eine** Fallunterscheidung
  im git-Adapter), aber sein Beleg ist der Bruch-Test, der auch die Doku trägt.
- **Ersatz-Form der Prüfung, vorab benannt:** **zwei Review-Runden gegen je
  einen abgeschlossenen Stand** — Runde 1 lag über der Doku (erledigt, 1 HIGH
  / 2 MEDIUM / 1 LOW), Runde 2 liegt über dem Code und der korrigierten Doku.
  Runde 2 trägt einen **deklarierten Fokus**: die Fallunterscheidung im
  Adapter gegen die **Gegenrichtung** — eine Datei, die im BASE-Tree
  wirklich nicht existiert (neu angelegte ADR), darf **weiterhin** befundfrei
  bleiben. *Ein Fix, der beide Fälle gleich behandelt, tauscht ein stilles
  Übersehen gegen einen Fehlalarm.*

**Zweite Plan-Änderung (2026-09-08, nach Review-Runde 2).** Der deklarierte
Fokus hat **genau getroffen, wovor er warnte** — und noch etwas dazu. Zwei
HIGH, beide selbst nachgemessen und beide eingearbeitet:

- **Der stille Pfad war nur zur Hälfte zu.** Liegt statt des BASE-**Blobs** das
  BASE-**Tree** des geschützten Verzeichnisses im unsichtbar benannten Pack,
  kommt die Datei als `A` an — und der `Added`-Zweig ruft `FileAt` gar nicht,
  der Adapter-Fix konnte dort nicht greifen. Gemessen: weiterhin
  `0 Befund(e)`/Exit 0. Der Zweig fasst den BASE-Stand jetzt an.
- **Der befürchtete Fehlalarm trat ein**, und zwar für Einträge **ohne**
  Datei-Inhalt: Ein wandernder Gitlink (Submodul) in der `vcs.paths`-Klasse
  brach mit Exit 2 ab, wo der Vor-Fix-Stand korrekt schwieg. `entryUnreadable`
  prüft jetzt den **Modus** des Tree-Eintrags.

**Damit wächst der Slice ein zweites Mal**, und die [`MR-066`](../../../../harness/conventions.md#mr-066)-Pflichten gelten
erneut — und die **dritte Plan-Änderung steht gleich darunter**, weil Runde 3
sie erzwungen hat: Der **Grund** ist derselbe (die Doku beschriebe sonst einen
Defekt, den es halb nicht mehr gibt); die **Ersatz-Form** wird um eine
**dritte Runde**
ergänzt, deren Fokus die *Vollständigkeit der Klasse* ist — nicht mehr die
Gegenrichtung einer Bedingung, sondern die Frage, welche **weiteren**
Diff-Zustände (`D`, Typänderung) denselben Weg nehmen könnten.


**Dritte Plan-Änderung (2026-09-08, nach Review-Runde 3) — und diesmal wird
NICHT erweitert.** Runde 3 fand eine **dritte** Ausprägung: Ein gelöschtes
oder umbenanntes geschütztes Verzeichnis mit unsichtbarem BASE-Tree erreicht
die Änderungsliste **gar nicht** — der Tree-Walker der Bibliothek macht aus
dem nicht ladbaren Unterbaum ein `io.EOF`. Selbst nachgemessen: `0 Befund(e)`,
Exit 0, gegen `1 Befund`/Exit 1 bei kanonischem Namen.

**Der Fix dafür kommt bewusst nicht in diesen Slice**, und der Grund ist der
Befund selbst: Dreimal wurde hier der gefundene **Pfad** geschlossen und nicht
die **Klasse**, jedes Mal vom Review und nie vom Autor. Ein vierter Patch auf
derselben Schicht wäre derselbe Fehler — und an der entscheidenden Stelle
ohnehin unmöglich: Der naheliegende Wachposten eine Schicht tiefer
(`tree.Files()`) benutzt denselben Walker und schweigt ebenso, **gemessen**.
Die dritte Ausprägung braucht einen **Schicht-Wechsel** und liegt deshalb bei
[`slice-220`](../open/slice-220-vcs-pfadmenge-statt-diff.md);
[`CO-001`](../../carveouts/CO-001-vcs-range-stiller-skip.md) führt sie bis
dahin. **Auftraggeber-Entscheid vom 2026-09-08.**

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [x] **(1)** [`harness/sensors/adr-check.md`](../../../../harness/sensors/adr-check.md)
      und [`harness/sensors/trace-check.md`](../../../../harness/sensors/trace-check.md)
      tragen die Grenze in ihrem Abschnitt *Grenze — was das Grün nicht
      abdeckt*: **Symptom, Exit-Code, Auslöser und Abhilfe**, jeweils mit der
      gemessenen Ausgabe statt einer Beschreibung.
- [x] **(2)** Das [Benutzerhandbuch](../../../../docs/user/benutzerhandbuch.md)
      nennt sie bei **beiden** git-basierten Modulen (§`vcs`, §`commits`) —
      denn sie ist eine Eigenschaft des **Produkts**, nicht dieses Repos, und
      ein Adopter liest das Handbuch, nicht unsere Sensor-Dateien.
- [x] **(3)** **Das Negativ-Ergebnis steht dabei**, und zwar als gemessenes:
      `tracked` ist **nicht** betroffen (es liest den git-**Index**, nicht die
      Objektdatenbank). Belegt durch **Positiv-Kontrolle** — derselbe
      `target-untracked`-Befund unter beiden Pack-Namen —, nicht durch einen
      grünen Lauf, der auch „nichts geprüft" heißen könnte.
- [x] **(4)** *(nach der Plan-Änderung, §1)* **Der stille Pfad ist behoben:**
      Ein **unlesbares** Objekt wird nicht mehr wie eine **im Tree fehlende**
      Datei behandelt. Beleg ist der Bruch-Test in **beide** Richtungen —
      derselbe partielle Pack, der heute Exit 0 liefert, meldet danach; **und**
      eine im BASE-Tree wirklich fehlende Datei bleibt weiterhin befundfrei.
      Dazu ein **Regressionstest** in `make test`, der ohne den Fix fällt.
- [x] **(5)** *(nach der Plan-Änderung, §1)* **`CO-001` ist angelegt** und im
      Carveout-Index eingetragen: Der Defekt steckt im **veröffentlichten**
      Bild und bleibt dort bis zum nächsten Release; Auflösungs-Trigger und
      Folge-Slice sind benannt.
- [x] `make gates` grün.
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

**Die Messung ist vor dem Plan gelaufen und steht in §1** — sie ist die
Grundlage, nicht ein Schritt. Was folgt, ist das Aufschreiben.

**Die Form der Aussage, bevor sie an drei Orte geht.** Eine Grenze trägt hier
vier Teile, und jeder muss aus einer Probe stammen:

1. **Symptom** — die wörtliche Meldung, und dass sie **variiert**:
   `staged-Basis "HEAD" nicht auflösbar: reference not found` (Probe-Repo, alle
   Objekte im Pack) gegen `HEAD-Tree nicht lesbar: object not found` (dieses
   Repo, wo der HEAD-Commit noch lose lag). **Ein Leser, der nur die eine
   kennt, erkennt die andere nicht wieder.**
2. **Exit-Code** — **2**, also fail-closed. Das ist die entlastende Hälfte des
   Befunds und gehört genauso hin wie die belastende: Es gibt **kein stilles
   Grün**.
3. **Auslöser** — ein Pack-Name, der von `pack-<eigener Hash>` abweicht;
   erzeugt etwa von `git maintenance run --task=loose-objects`.
4. **Abhilfe** — `git repack -A -d` (die `-A`-Form, damit unerreichbare
   Objekte lose werden statt verworfen).

**Schritte:**

1. Die zwei Sensor-Dateien ergänzen — als **nummerierter Punkt** in ihrem
   vorhandenen Grenzen-Abschnitt, in dessen Form.
2. Das Handbuch bei §`vcs` und §`commits` ergänzen. Beide Stellen brauchen die
   Aussage, weil ein Adopter nur eines der Module fahren kann.
3. `make doc-check`, dann `make gates`, dann Handoff an den unabhängigen
   Review.

**Was der Plan nicht mitnimmt**, steht in §1 — vier Punkte, die den Lauf
binden.

## 4. Trigger

**Beanspruchung:** WIP-Limit frei (`in-progress/` leer nach der Closure von
slice-217), und der Auftraggeber hat den Zuschnitt am 2026-09-08 beauftragt,
nachdem der Befund in slice-217 ausdrücklich als eigener Vorgang benannt
worden war.

**Rückführung nach `next/`** (`in-progress→next`): wenn sich beim Schreiben
zeigt, dass die Grenze eine **Verhaltens**-Antwort verlangt statt einer
Beschreibung — etwa weil die Meldung so irreführend ist, dass ein Adopter sie
für einen Defekt seines Repos hält. Dann ist der Gegenstand Produkt-Code und
der Slice falsch geschnitten (§1 Punkt 2 schließt Code aus).

**Rückführung nach `open/`** (`in-progress→open`): wenn die Messung sich beim
Aufschreiben als nicht reproduzierbar erweist — etwa weil ein weiterer
Faktor mitspielte, den die sechs Proben nicht isoliert haben.

## 5. Closure-Trigger

DoD (1) bis (3) abgehakt, `make gates` grün mit echter Ausgabe, unabhängiger
Review durchgeführt und seine Befunde eingearbeitet, Closure-Notiz geschrieben,
Register fortgeschrieben, jedes Risiko aus §6 mit einem der drei Ausgänge.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Der Mechanismus ist erschlossen, nicht gelesen.** Gemessen ist das
  **Verhalten** über sechs Proben; **warum** go-git so entscheidet, steht in
  seiner Quelle, und die ist ohne Host-Go nicht lesbar (§3.1). Der Slice darf
  deshalb nur sagen, *was* passiert, und muss die Erklärung als das kennzeichnen,
  was sie ist. **Die erste Erklärung war bereits falsch:** *„go-git listet
  Packs nach dem `pack-`-Präfix"* — widerlegt durch die Probe mit gültiger Form
  und fremdem Hash, die ebenfalls scheitert. Wer den zweiten Erklärungsversuch
  für sicherer hält als den ersten, hat nichts gelernt. — **Ausgang:**
  eingetreten, und zwar genau so. Der zweite Erklärungsversuch (*„beide
  Namensteile zählen"*) hielt — Review-Runde 2 hat ihn in der go-git-Quelle
  **gelesen** (Namens-Filter plus Checksummen-Vergleich) und bestätigt. Der
  **dritte** war dann wieder falsch: Für die Tree-Hälfte führten alle vier
  Doku-Orte dieselbe `ErrFileNotFound`-Zusammenfassung als Ursache an,
  obwohl dort ein anderer Mechanismus wirkt (der Tree-Walker macht `io.EOF`).
  **Die Warnung des Risikos war also richtig und hat trotzdem nicht
  geschützt** — sie stand im Plan, und der Fehler passierte erneut.
- **Drei Orte, eine Aussage — das ist eine Drift-Quelle.** Sensor-Datei,
  zweite Sensor-Datei und Handbuch tragen dieselbe Grenze; ändert sich das
  Verhalten (go-git-Bump), müssen alle drei mit. **Kein Sensor koppelt sie**,
  und `mentions` wäre das falsche Werkzeug: Es misst *erwähnt*, nicht
  *übereinstimmend*. — **Ausgang:** eingetreten, dreimal. Jede Review-Runde
  fand mindestens eine Stelle, an der einer der Orte der Wirklichkeit
  hinterherhinkte: die widerlegte fail-closed-Zusage (R1), der Abschluss über
  die ganze Klasse bei gemessener Hälfte (R2), die Ursachen-Angabe und der
  `RANGE=`-only-Geltungsbereich (R3). **Die Gegenmaßnahme ist nicht mehr
  Sorgfalt, sondern Struktur:** Die vollständige Messung steht jetzt an
  **einem** Ort (`harness/sensors/adr-check.md`), die anderen verweisen
  darauf, statt sie zu wiederholen.
- **Die Grenze könnte weiter reichen als die drei gemessenen Module.** Geprüft
  sind `vcs`, `commits` und `tracked`. Ob eine andere Stelle des Produkts die
  Objektdatenbank liest, ist **nicht** geprüft — und §3.8 verlangt genau diese
  Frage. Der Slice darf deshalb nicht *„die git-basierten Module"* schreiben,
  wenn er drei gemessen hat. — **Ausgang:** eingetreten, in beide Richtungen.
  Zu **eng** war die erste Fassung nicht — sie nannte drei Module und hatte
  drei gemessen. Zu **weit** war die Aussage über die *Diff-Zustände*: „geprüft
  sind `A`, `M`, `D`" las sich als Deckung der Klasse, während der dritte
  stille Pfad **vor** dem Diff entsteht und in keinem der drei Zustände
  ankommt. Der Text sagt das jetzt ausdrücklich und führt die offene Zeile
  als solche.

## 7. Closure-Notiz

**Geliefert.** Eine gemessene Produkt-Grenze hat einen Ort bekommen — und
unterwegs stellte sich heraus, dass sie ein **Defekt** ist. Zwei von drei
Ausprägungen sind geschlossen (Adapter und Regel), drei Regressionstests
halten sie, die Spezifikation trägt den Mechanismus, `CO-001` führt die dritte,
und `slice-220` geht sie strukturell an. **Drei** Review-Runden, blockierend:
1 HIGH / 2 MEDIUM / 1 LOW · 2 HIGH / 4 MEDIUM / 4 LOW · 1 HIGH / 6 MEDIUM /
4 LOW. `make gates` grün.

**Was funktioniert hat: die vorab deklarierte Ersatz-Form der Prüfung.** Der
Plan hat den Fokus jeder Runde **vor** der Arbeit benannt, und jede Runde traf
genau dort. Runde 2 sollte die Gegenrichtung prüfen — und fand den Fehlalarm
beim Gitlink, den §1 wörtlich vorhergesagt hatte (*„ein Fix, der beide Fälle
gleich behandelt, tauscht ein stilles Übersehen gegen einen Fehlalarm"*).
Runde 3 sollte die Vollständigkeit der Klasse prüfen — und fand die dritte
Form. **Ohne die Deklaration wäre beides Zufall gewesen.**

**Was Friktion war, und es ist ein Muster statt eines Fehlers: dreimal habe
ich den gefundenen Pfad geschlossen und nicht die Klasse.** Jeder Fix hatte
einen Bruch-Test in beide Richtungen und eine Kontroll-Matrix; jeder sah
vollständig aus. Die nächste Form lag jedes Mal **eine Schicht davor** —
hinter dem Adapter die Regel, hinter der Regel die Diff-Erzeugung. **Gefunden
hat es jedes Mal der Review, nie der Autor.** Der neue Registereintrag
[`fix-schliesst-pfad-nicht-klasse`](../observations/BEO-ALL/fix-schliesst-pfad-nicht-klasse/observation.md)
trägt den Ableiter: *erst die Schicht benennen, auf der ein Fix wirkt, dann
prüfen, was davor liegt* — eine Frage, die dreimal fällig gewesen wäre und
zwei Sekunden gekostet hätte.

**Steering-Loop-Lerneintrag, und er korrigiert den Review.** Der Eintrag steht
bei **1×**, nicht bei 3×: Die drei Instanzen liegen in **einem** Vorgang, und
der Zähler misst Wiederholung über Vorgänge hinweg. Runde 3 nannte die
Schwelle erreicht; das trifft nicht zu. **Die Zahl passend zu machen wäre die
bequemere Antwort gewesen** — der Eintrag benennt die Korrektur stattdessen in
seinem Stand.

**Zweiter Lerneintrag: eine Zusage über Ausfälle stimmt selten beim ersten
Mal.** Ich schrieb *„Exit 2 ist fail-closed — es gibt kein stilles Grün"* als
**entlastende** Hälfte einer Grenze, und sie war falsch. Gemessen sind heute
**drei** Ausgänge: Exit 2 (die meisten Formen), Exit 1 mit einer
`core-drift-vcs`-**Fehldiagnose** (unsichtbarer HEAD-Tree) und Exit 0 (die
offene Form). Beleg bei
[`wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md)
(10×).

**Dritter: drei Orte für eine Aussage sind eine Drift-Quelle, und §6 hat es
vorhergesagt.** Jede Runde fand einen Ort, der hinterherhing. Die Antwort ist
nicht mehr Sorgfalt, sondern Struktur — die vollständige Messung steht jetzt
an **einem** Ort, die anderen verweisen. Beleg bei
[`semantic-change-body-only-edges-stale`](../observations/BEO-ALL/semantic-change-body-only-edges-stale/observation.md)
(16×).

**Was offen bleibt — benannt, nicht verschwiegen.** **(1)** Die dritte
Ausprägung: ein gelöschtes/umbenanntes geschütztes Verzeichnis mit
unsichtbarem BASE-Tree erreicht die Änderungsliste nie. **Sie trifft auch
dieses Repo**, nicht nur gepinnte Konsumenten — `CO-001` sagt das jetzt, nachdem
zwei frühere Fassungen des Feldes in **beide** Richtungen falsch lagen.
**(2)** Eine `Accepted`-ADR, die durch einen **Gitlink** ersetzt wird, bleibt
befundfrei; das ist der Preis der Gegenrichtungs-Regel und steht bei ihr.
**(3)** Der unsichtbare HEAD-Tree wird als Inhalts-Befund fehldiagnostiziert —
laut, aber falsch; `slice-220` DoD (2) nimmt ihn mit.

**Die drei Paarungen, gemessen.** **(a) Anker** — vakant: Der Slice verkörpert
keine Steering-Loop-Regel; sein neuer Registereintrag steht bei 1× und trägt
folgerichtig **keinen** Ausgang. **(b) Folge-Slice** — zwei, beide als Datei
im Lifecycle: [`slice-219`](../open/slice-219-release-loest-co-001.md) (das
Release) und [`slice-220`](../open/slice-220-vcs-pfadmenge-statt-diff.md) (der
Klassen-Fix); `CO-001` nennt beide in der geforderten Reihenfolge. **(c)
Register** — alle zitierten Pfade lösen auf; ein neues Verzeichnis mit
`observation.md`, `state.md` und nicht leerem `evidence/`, dazu zwei Belege an
bestehenden Einträgen.

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

**Eine** Sub-Area: `*` (Repo-Default). Der Slice ändert drei Dokumente. Der
**Gegenstand** ist zwar Produkt-Verhalten, aber es wird **gemessen, nicht
geändert** — und für Produkt-Code führt dieses Repo keine eigene Sub-Area; die
Modus-Deklaration kennt `*` und `tools/harness/`. Eine dritte hier zu erfinden
wäre eine Sub-Area ohne Deklaration.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:274-274 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, **39** Verzeichnisse über beide
Kürzel; `BEO-HARN` einzeln geöffnet — der eine offene Eintrag dort betrifft
`--check-latest` und berührt diesen Slice nicht). **Vier** Einträge sind
einschlägig:

- [`module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
  (2×) — **der tragende, und er ist die Quelle von §6 Risiko 3.** Die Frage
  *„welche Eingaben liest ein Modul, die es nicht scannt?"* ist hier wörtlich
  der Gegenstand: `vcs` und `commits` lesen die **Objektdatenbank**, die in
  keiner Scan-Wurzel steht. Die Grenze **ist** die Antwort auf diese Frage.
- [`wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md)
  (9×) — die Gegenrichtung: Eine Grenze, die mehr Module nennt als gemessen
  wurden, behauptet eine Prüfung, die nicht stattfand. §6 führt es als drittes
  Risiko.
- [`grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md)
  (8×) — dieser Slice **schreibt** Grenzen-Listen, an drei Orten. Der Ableiter
  gilt doppelt: den Vertrags-Teil umdrehen, und wo der Gegenstand Verhalten
  ist, gegen das Verhalten prüfen statt gegen die Prosa.
- [`form-vom-nachbarn-statt-von-der-vorlage`](../observations/BEO-ALL/form-vom-nachbarn-statt-von-der-vorlage/observation.md)
  (1×, gestern angelegt) — die **Form** des Grenzen-Punkts wird aus dem
  vorhandenen Abschnitt derselben Datei übernommen. Das ist hier **richtig**:
  Der Eintrag warnt vor der Form-Übernahme, wo eine **Vorlage** existiert; für
  einen Punkt in einer bestehenden Liste ist die Liste selbst die Form. Notiert,
  damit die Übernahme eine Entscheidung ist und keine Gewohnheit.

**Keiner der vier erreicht mit diesem Slice die Schwelle erstmalig.**

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-08 gelesen: **beide Nachtläufe grün** —
`upstream-drift.yml` (jüngster Lauf 2026-09-08T05:31:31Z) und `image-scan.yml`
(2026-09-07T08:21:32Z). **Für diesen Slice ausnahmsweise nicht bezugslos:**
Der Gegenstand ist an `go-git v5.19.2` gebunden, und ein Dependabot-Bump dieser
Abhängigkeit könnte das gemessene Verhalten ändern. Die Grenze nennt deshalb
die Version.

**Modus-Begründungsblock.** Alle berührten Sub-Areas GF — ein Block genügt.

### Sub-Area: `*`

- **Modus:** GF (Repo-Default).
- **Konventions-Dichte:** hoch für die **Form** einer Sensor-Grenze — beide
  Zieldateien führen den Abschnitt *Grenze — was das Grün nicht abdeckt* mit
  nummerierten Punkten bereits.
- **Phase-Reife:** Phase 5. Die Sensor-Dateien sind gewachsen und wurden in
  slice-212 zuletzt systematisch auf fehlende Grenzen durchgesehen.
- **Evidenz-/Diskrepanz-Risiko:** **niedrig für den Bestand, mittel für die
  Aussage.** Zu inventarisieren ist nichts. Das Risiko sitzt in der Reichweite:
  drei gemessene Module, aber eine Formulierung, die leicht für alle spricht —
  und in der Erklärung, die schon einmal falsch war.
- **Reconciliation-Aufwand:** keiner (GF).
