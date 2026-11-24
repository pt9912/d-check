# Welle welle-97: `linelength` nimmt eingerückte Tabellenzeilen aus — eine `|`-Zeile mit bis zu drei Leerzeichen davor ist eine Tabellenzeile

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-11-24.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Das Modul `linelength` nimmt mit `linelength.exempt-tables` an (Default) auch eine `|`-Zeile aus,
vor der **höchstens drei Leerzeichen** stehen — die Tabelle in einem Listenpunkt, wie sie unter
`1. ` mit drei Leerzeichen Einrückung steht. Heute ist nur eine `|` in Spalte 0 eine
Tabellenzeile; die eingerückte Zeile wird gemessen, und die Anforderung sagt das ausdrücklich zu
([`DC-FA-LLEN-001`](../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in),
Abschnitt „Was Fenced-Code und was eine Tabellenzeile ist“; Entscheidung 5 (b) von
[`ADR-0114`](../adr/0114-linelength-modul.md)). Ab vier Leerzeichen, mit einem Tab davor und hinter
`>` bleibt die Zeile ein Befund, wie heute; ein Lauf ohne eine eingerückte `|`-Zeile ändert sich in
keinem Byte. Das Akzeptanzkriterium, an dem sich dieses Ziel spiegelt: *Given eine überlange `|`-Zeile mit
drei Leerzeichen davor, when der Lauf im Default erfolgt, then kein Befund; mit
`exempt-tables: false` ein Befund.*

Anlass ist [Issue #20](https://github.com/pt9912/d-check/issues/20): Am Pin `v0.92.0` meldet das Modul
im Bestand des Melders 41 von 85 Befunden an solchen Zeilen (Angabe des Issues, von hier nicht
gemessen — der Bestand liegt außerhalb dieses Repos). **Der Zuschnitt folgt der Lehre aus**
[`BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass`](observations/BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass/observation.md):
je offener Frage des Issues die kleinste Antwort, die die Kriterien des Anlasses trägt (§5), und
höchstens etwa acht Akzeptanzkriterien — geplant sind sechs.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Der Auftraggeber hat [Issue #20](https://github.com/pt9912/d-check/issues/20) als Welle angefordert
  (2026-11-24, in der Sitzung).
- Die Welle für [Issue #18](https://github.com/pt9912/d-check/issues/18) ist geschlossen: ihre
  Welle-Datei liegt unter `done/welle-96/` (`ls docs/plan/planning/done/welle-96/`), ihre
  Closure-Notiz als `done/welle-96-results.md`.
- `git fetch`: `origin/main` liegt nicht vor dem lokalen Hauptzweig (0 Commits hinter, 65 voraus,
  ungepusht — `git rev-list --left-right --count origin/main...HEAD`). `in-progress/` trägt nur die
  Roadmap, `next/` ist leer, `open/` trägt nur die zwei Slices dieser Welle.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle sind `done/`.
- `make fullbuild` grün.
- **Das Mehr gegenüber den Slice-DoDs** (repo-weit, in keiner einzelnen DoD):
  - **(a) Ohne den Wunsch nichts Neues.** Zwei Läufe, die zusammen tragen, was jeder allein nicht trägt.
    **(a1)** `make blackbox-probe REF=69c53e0f` — der Self-Close von welle-96, der letzte Stand **vor**
    dieser Welle — endet mit „byte-identisch“ über seine Vergleiche (stdout, stderr und Exit
    getrennt). **Nicht `REF=v0.92.0`:** das Profil dieses Repos schaltet seit welle-96
    `vcs.follow-renames` ein, das Vorher-Image `v0.92.0` kennt den Schlüssel nicht und lehnt die
    `.d-check.yml` ab; der Eigenlauf `repo` endet dort vorher mit Exit 2, nachher mit Exit 0, und
    „4 von 20 Vergleichen weichen ab (Vorher 61a412eb)“
    ([Verifikation zu welle-96](../../reviews/2026-11-24-welle-96-verifikation.md)). Diese Abweichung
    kommt aus welle-96 und nicht aus dieser Welle
    ([Architect-Verdikt](../../reviews/2026-11-24-welle-97-architect-verdikt.md) F-2). Das Target fährt `linelength` **nicht**: unter
    `tools/blackbox-probe/fixtures/` liegen `ids`, `links`, `sauber` und `targets`
    (`ls tools/blackbox-probe/fixtures/`), und `grep -rl linelength tools/blackbox-probe/` findet
    nichts. Es belegt deshalb nur, dass die Änderung kein Modul verschiebt, das **das Profil dieser
    Fixtures und das des Repos selbst** (der Lauf `repo` mit eigener `.d-check.yml`) einschaltet —
    darunter `targets`, das die Tabellenzeilen-Antwort der Lexik liest; ob dieses Profil `structure`
    und `planning.waves` einschaltet, liest die Closure an der `.d-check.yml` und nennt es in der
    Ergebnis-Notiz. Der Beleg für beide bleibt der Kopplungs-Test aus
    [`slice-294`](open/slice-294-linelength-eingerueckte-tabellen-kern.md). **(a2)** Ein Vorher-/Nachher-Lauf mit `--enable linelength`: dasselbe Vorher-Image
    wie in (a1) (`git archive 69c53e0f`), gebaut wie in `tools/blackbox-probe.sh` (`docker build --target runtime`),
    das Nachher-Image aus `make build`; beide über dasselbe Wegwerf-Fixture, eine Kopie von
    `docs/user/` mit einem `linelength`-Block (`max-chars` kleiner als die längste Zeile dort). Das
    Fixture trägt keine Zeile mit ein bis drei Leerzeichen vor `|` (`grep -rEc '^ {1,3}\|'` über
    die Kopie ergibt 0; am 2026-11-24 gemessen für `docs/user/`). stdout, stderr und Exit sind
    **gleich**. **Die eine gewollte Abweichung** liegt in einem zweiten Fixture mit einer langen
    Tabelle in einem Listenpunkt: Vorher ein Befund, Nachher keiner, und Nachher mit
    `exempt-tables: false` wieder derselbe Befund wie Vorher. Die Bezugsmengen-Zeile auf stderr
    nennt dort Nachher eine kleinere Zahl geprüfter Zeilen — gewollt, Kriterium „Bezugsmenge“ in
    [`slice-293`](done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md) §3.
  - **(b) Der Anlass trägt am gebauten Image.** Ein Wegwerf-Fixture aus der Form des Issues: ein
    nummerierter Listenpunkt, darunter eine Tabelle mit drei Leerzeichen Einrückung und einer
    überlangen Zeile. Default: kein Befund. Mit `exempt-tables: false`: ein Befund. Dieselbe Tabelle
    mit vier Leerzeichen, mit einem Tab und hinter `>` bleibt je ein Befund — die Grenze, die der
    Wunsch enger zieht als sein Wortlaut („2 bis 5 Leerzeichen“), steht damit in der Probe, nicht
    nur im Text. Die Tabelle **Probe → Ergebnis** (Fixture, Schalter, erwartet, gemessen) steht in
    `welle-97-results.md`.
  - **(c) Das Trigger-Audit** (seit slice-241) ist über die Welle-Closure vollzogen und in
    `welle-97-results.md` belegt. Fällig ist in dieser Welle der Re-Evaluierungs-Trigger **(b)** von
    [`ADR-0114`](../adr/0114-linelength-modul.md) („ein Konsument eine Zeile ausnehmen muss, die weder
    Code noch Tabelle noch unumbrechbar ist“): Die eingerückte Tabellenzeile ist nach Entscheidung 5 (b)
    dort keine Tabelle. **Nicht (a)** („die geteilte Lexik ihre Fence- oder Tabellen-Antwort ändert“):
    `tableRowLine` bleibt zeichengleich
    ([Verdikt](../../reviews/2026-11-24-welle-97-architect-verdikt.md) F-5). Die Closure-Notiz nennt
    [`ADR-0116`](../adr/0116-linelength-eingerueckte-tabellenzeilen.md) als Antwort (sie wählt den Helfer
    in der Lexik und verwirft den Marker, Alternative G).
- Closure-Notiz in `welle-97-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-293`](done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md) | Anforderung, Spezifikation und Folge-ADR zu [`ADR-0114`](../adr/0114-linelength-modul.md): `exempt-tables` nimmt eine `\|`-Zeile mit bis zu drei Leerzeichen davor aus; Entscheide festschreiben | [`DC-FA-LLEN-001`](../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in) |
| [`slice-294`](open/slice-294-linelength-eingerueckte-tabellen-kern.md) | Kern und Tests: die Erkennung in `linelength`, Kopplungs-Test, Handbuch und Abdeckung | [`DC-FA-LLEN-001`](../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in) in der Fassung aus [`slice-293`](done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine Welle. Das Issue bleibt bis zum Release offen; Release, Tag und Antwort auf das
  Issue sind ein eigener Vorgang nach der Welle (§6).
- Wird blockiert von: keiner.
- **Reihenfolge der Slices ist fachlich bindend:** erst
  [`slice-293`](done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md) (Anforderung
  und Entscheide stehen, bevor Code entsteht), dann
  [`slice-294`](open/slice-294-linelength-eingerueckte-tabellen-kern.md). Das WIP-Limit 1 hält sie
  ohnehin nacheinander.
- **Die kleinste Antwort auf jede offene Frage des Issues** — im Wortlaut, mit der Stelle, die sie
  trägt, und ihrem Preis. Die Lehre aus
  [`BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass`](observations/BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass/observation.md)
  verlangt, sie vor der Empfehlung auszuschreiben; die größere Antwort steht daneben, damit der
  Auftraggeber ihren Mehraufwand sieht, bevor er bestätigt. [`slice-293`](done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md)
  schreibt sie fest; ohne Bestätigung beginnt kein Code.
  1. *Eigener Schalter (etwa `exempt-indented-tables`) oder `exempt-tables` erweitern?* —
     **„`linelength.exempt-tables` wird selbst erweitert; es gibt keinen neuen Schlüssel.“** Trägt:
     die Ausnahme (b) in der Anforderung
     ([`DC-FA-LLEN-001`](../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in),
     Absatz „Drei Ausnahmen, je einzeln schaltbar“) und das Schema `linelength.*`
     ([`SPEC-005`](../../../spec/spezifikation.md#spec-005--d-checkyml)). Ein zweiter Schalter
     brächte eine Schema-Zeile, einen Fall am Konfigurations-Rand (Wahrheitswert, Schlüssel ohne
     Wert ⇒ Exit 2) und Kriterien je Wert. Preis der kleinen Antwort: wer `exempt-tables` heute
     im Default fährt, bekommt für eingerückte Tabellenzeilen keinen Befund mehr — keiner bekommt
     einen **zusätzlichen**; die Bezugsmengen-Zeile nennt weniger geprüfte Zeilen.
  2. *Gilt `> |` (Tabelle im Zitat) ebenso?* — **„Nein. Eine Zeile hinter `>` bleibt gemessen, wie
     heute.“** Trägt: derselbe Absatz der Anforderung — „… eine Zeile in einem Zitat (`> |`) und eine
     Tabellenzeile ohne führende Pipe (`a | b`) sind keine Tabellenzeilen und werden gemessen.“ Das
     Issue benennt `> |` nicht als Anlass; wer es will, stellt eine eigene Anforderung (§6).
  3. *Ein eingerückter Code-Block (vier und mehr Leerzeichen ohne Listenkontext) bleibt unberührt* —
     die geteilte Lexik kann den Listenkontext **nicht** hergeben: sie beantwortet „liegt die Zeile
     außerhalb von Fenced-Code“ (`proseLineSet`) und „beginnt sie mit `|`“ (`tableRowLine`), und
     die Listen-Muster im Code (`taskItemRE` in `structure.go`, `listCodeRE` in
     `vcs_link_targets.go`) sind private Einzelmuster mit je eigenem Zweck. Deshalb die Fassung,
     die im Zweifel **mehr misst**: **„Ausgenommen wird eine `|`-Zeile, vor der höchstens drei
     Leerzeichen stehen. Ab vier Leerzeichen und mit einem Tab davor bleibt sie ein Befund, auch in
     einem Listenpunkt.“** Der Schnitt bei vier Leerzeichen ist der, den die Anforderung für den
     eingerückten Code-Block schon nennt („Ein eingerückter Code-Block (vier Leerzeichen, ohne Zaun)
     …“); bis drei Leerzeichen ist die Zeile nach dieser Lesart keine
     Code-Zeile, ohne dass die Lexik einen Listenkontext kennen muss. **Sie ist enger als der
     Wortlaut des Issues („2 bis 5 Leerzeichen“) — der eine Punkt, den der Auftraggeber
     bestätigt hat (2026-11-24, „bis 3 Leerzeichen“):** Eine Tabelle unter einer Marke von vier Spalten (`10. `) oder in einem
     verschachtelten Listenpunkt steht mit vier und mehr Leerzeichen und bleibt ein Befund. Ob der
     Bestand des Melders solche Zeilen trägt, ist von hier nicht gemessen. Die **größere Antwort**
     (4 und 5 Leerzeichen nur im Listenkontext, mit einer neuen Listen-Erkennung in der geteilten
     Lexik und Kriterien für ihre Ränder) ist nicht geplant; sie wird ein eigener Vorgang, wenn der
     Bestand des Melders nach der Hebung Befunde dieser Form behält.
  4. *Wo wohnt die Antwort?* — **entschieden** im
     [Architect-Verdikt](../../reviews/2026-11-24-welle-97-architect-verdikt.md) (F-1) und in
     [`ADR-0116`](../adr/0116-linelength-eingerueckte-tabellenzeilen.md) Entscheid 4, die ihn trägt:
     „Ein neuer Helfer in der geteilten Lexik neben `tableRowLine` beantwortet die Frage aus Entscheid 3.
     Er **schließt `tableRowLine` ein**: Er ruft sie, statt die Spalte-0-Antwort ein zweites Mal zu
     schreiben, und erweitert sie um die Einrückung. `linelength` ruft für (b) nur diesen Helfer.“ So
     bewegt sich weder `targets` noch `planning.waves` noch `structure`.
     Preis: Die Zusage „dieselbe Antwort, die die übrigen Module auf diese Frage geben, keine zweite“
     (Anforderung, Absatz „Was Fenced-Code und was eine Tabellenzeile ist …“; Out-of-Scope (4);
     [`ADR-0114`](../adr/0114-linelength-modul.md) Kontext) gilt für diese Frage dann nicht mehr
     wörtlich; sie wird in [`slice-293`](done/slice-293-linelength-eingerueckte-tabellen-anforderung-und-adr.md)
     **umgeschrieben**, nicht neben einem Widerspruch stehen gelassen
     ([`BEO-ALL/plan-entscheid-steht-gegen-lastenheft`](observations/BEO-ALL/plan-entscheid-steht-gegen-lastenheft/observation.md)).
     Die Alternative — `tableRowLine` selbst weiten — hält die eine Antwort, verschiebt aber drei
     Module und braucht Kriterien für jedes.
- **Review-Vorgabe der Welle** (Lehre aus welle-96): je Slice **höchstens zwei Runden**, eingearbeitet
  werden nur HIGH und MEDIUM; LOW und INFO bleiben im Report. Jede Nachbesserung zitiert die Stelle,
  die sie trägt, statt sie nachzuerzählen.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Zitate (`> |`)** — bleiben gemessen; kein Teil des Anlasses, und eine Zitat-Erkennung für
  Tabellenzeilen wäre eine zweite Antwort neben der geteilten Lexik. Ein eigener Vorgang, wenn der
  Auftraggeber ihn verlangt.
- **Ein eigener Schalter** (`exempt-indented-tables` o. ä.) — Frage 1 in §5; ein zusätzlicher
  Schlüssel trüge Schema, Konfigurations-Rand und Kriterien je Wert, ohne dass der Anlass ihn braucht.
- **Eine Listen-Erkennung** und damit Einrückungen von vier und fünf Leerzeichen — Frage 3 in §5;
  eine neue Antwort in der geteilten Lexik mit eigenen Kriterien für ihre Ränder. Ohne sie bleibt ein eingerückter
  Code-Block unberührt, auch mit `|` am Zeilenanfang.
- **Andere Module, die Tabellenzeilen lesen** (`targets`, `planning.waves`, `structure`) — sie lesen
  eine eingerückte `|`-Zeile weiter nicht als Tabellenzeile; `tableRowLine` bleibt, wie es ist
  ([`ADR-0116`](../adr/0116-linelength-eingerueckte-tabellenzeilen.md) Entscheid 4).
- **Ein Umbau der geteilten Lexik** über den einen neuen Helfer hinaus.
- **`--fix-patch`, Umbrechen oder Kürzen** — das Modul schreibt nie
  ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit),
  Out-of-Scope (5) der Anforderung); eine Tabellenzeile lässt sich ohnehin nicht umbrechen.
- **Das eigene Profil des Repos** — schaltet `linelength` nicht ein; die Welle ändert daran nichts
  ([`ADR-0114`](../adr/0114-linelength-modul.md), Konsequenzen).
- **Release, Tag, Release-Prep** (CHANGELOG, READMEs, Handbuch-Kopf, Digest-Pins) und die Antwort
  auf das Issue — eigener Vorgang nach der Welle
  ([`AGENTS.md`](../../../AGENTS.md) §5 Zeile 17).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-97-results.md` (Geschwister im Ruheort `done/`, entsteht bei der Closure)
Zähler: [Beobachtungs-Register](observations/README.md)
