# ADR-0050 — Unbalancierter Fence: Zustand melden statt Paarung reparieren, im Modul `spans`

**Status:** Proposed
**Datum:** 2026-08-09
**Autor:** pt9912
**Bezug:** [`DC-FA-SPAN-001`](../../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
(dritte Artefakt-Klasse); die betroffene Fence-Lexik stammt aus
[ADR-0042](0042-markdown-lexik-folgt-commonmark.md) (dort wurde der
längenabgeglichene Fence-Schluss **bewusst offen gelassen**);
Schnitt-Kriterium [ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md);
betroffener Konsument der Lexik
[`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Closure-Note-Struktur, [ADR-0048](0048-closure-note-struktur-im-planning-modul.md)).
**Schärft:** die Erweiterung von
[`spec/spezifikation.md` §DC-FA-SPAN-001.a](../../../spec/spezifikation.md#dc-fa-span-001a--span-artefakt-erkennung)
um Schritt 3.

## Kontext

Ein unabhängiger Review fand einen **ausgelieferten** stillen Grün-Pfad: hinter
einem nicht geschlossenen Fenced-Code-Block überspringt jede Vorverarbeitung den
Rest des Dokuments. Reproduziert gegen das veröffentlichte Image — dieselbe
Notiz meldet mit Fence **0 Befunde / Exit 0**, ohne Fence einen Befund und
Exit 1. Ein Autor, der einen Code-Block nicht schließt, schaltet damit unbemerkt
die Prüfung des Restdokuments ab.

Die Lexik ist nicht neu: [ADR-0042](0042-markdown-lexik-folgt-commonmark.md) hat
den naiven Fence-Toggle bewusst behalten und den längenabgeglichenen Schluss als
offene Grenze benannt. Das war vertretbar, solange Module über **ganze** Dateien
urteilten. Mit einer Bedingung, die **innerhalb** eines Abschnitts misst, wurde
aus der Grenze ein Silent-Grün.

**Bestandsmessung (2026-08-09)** mit d-checks eigener `FenceToggle`-Lexik über
drei Repos — das eigene und die zwei, die die offenen Change Requests gestellt
haben:

| Repo | Markdown-Dateien | unbalanciert | gemischte Fence-Längen | `~~~`-Fences |
|---|---|---|---|---|
| d-check | 347 | 0 | 0 | 0 |
| a-check | 184 | 0 | 0 | 0 |
| ai-harness-course | 245 | 0 | 0 | 0 |

**776 Dateien, null Vorkommen.** Der Defekt ist **latent**: kein Dokument im
Ökosystem löst ihn heute aus.

## Entscheidung

1. **Gemeldet wird der Zustand, nicht die Paarung — und zwar in beiden
   Lesarten.** Neuer Grund-Code `fence-unclosed`: eine Fence-Öffnung, die bis
   zum Dateiende keinen Schluss findet. Die Paarungsregel selbst bleibt
   **unverändert**.

   Das Produkt trägt zwei Schluss-Lesarten: den naiven Toggle (Vorverarbeitung
   und alles, was auf ihr aufsetzt) und den längenabgeglichenen
   CommonMark-Schluss (Tabellen-Leser). Der Wächter wertet **beide** aus und
   meldet, sobald **eine** von beiden am Dateiende offen steht — genau dann
   überspringt mindestens ein Modul den Rest. Nur eine auszuwerten hieße, die
   andere unbewacht zu lassen: ein Backtick-Block, den eine Tilden-Zeile
   „schließt“, ist für den naiven Toggle balanciert und für den
   Tabellen-Leser bis Dateiende offen.
   Damit diese Lesart nicht als Kopie neben ihrem Konsumenten liegt, ist der
   längenabgeglichene Schluss zu einem geteilten Prädikat geworden.
   Das ist keine Bequemlichkeit, sondern folgt aus der Messung und aus einer
   Beobachtung, die beim ersten Zuschnitt fehlte: **ein strengerer Schluss löst
   diesen Fall gar nicht.** Er korrigiert *Fehlpaarungen* — ein Fence, der
   **nie** geschlossen wird, bleibt unter jeder Paarungsregel offen. Der zuerst
   als „die Wurzel" notierte Kandidat war die Wurzel eines **anderen** Befundes.
   Zusätzlich zeigt die Messung, dass er auf dem realen Bestand wirkungslos wäre:
   ohne gemischte Fence-Längen und ohne `~~~`-Fences gibt es nichts zu
   korrigieren.

2. **Der Befund wohnt im Modul `spans`, nicht in `planning`.** Der Fundort war
   die Closure-Note-Struktur, aber ein offener Fence ist kein Planning-Thema.
   `spans` sagt genau diese Frageform bereits zu — „ungeschlossene Code-Spans …
   kippen die Backtick-Parität des restlichen Absatzes". Ein unbalancierter
   **Fence** ist dieselbe Aussage eine Ebene höher: eine Öffnung ohne Schluss,
   die alles Folgende umdeutet. Nach dem Kriterium aus
   [ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md) ist das kein neues
   Kürzel, sondern eine **Erweiterung der bestehenden Anforderung**: dasselbe
   Modul, dieselbe Frage, andere Ebene.

3. **Reichweite ist die Datei, nicht der Absatz.** `span-unclosed` misst
   absatzweise; für einen Fence ist das nicht übertragbar, weil er selbst eine
   Absatzgrenze **ist**. **Genau ein** Befund je Datei: hinter einem offenen
   Fence kann keine zweite Öffnung mehr gemeldet werden, ohne zu raten.

   Die gemeldete Zeile ist eine **Fundstelle**, keine Reparaturstelle: die
   strenge Lesart hat Vorrang, weil sie eine tatsächlich offene Öffnung kennt,
   aber welche Öffnung *fehlt*, ist grundsätzlich nicht entscheidbar. Gleich
   lange Öffner sind ununterscheidbar, und auch die strenge Lesart zeigt
   daneben, wenn eine längere Fence-Zeile eine kürzere Öffnung geschlossen hat.
   Diese Ehrlichkeit ist der Preis dafür, den Zustand überhaupt zu melden.

4. **Kein neuer Config-Schlüssel.** `spans` ist opt-in und trägt seine beiden
   bestehenden Klassen ohne Schalter; eine dritte bekommt keinen. Wer sie nicht
   will, deaktiviert das Modul. Ein Schalter je Klasse wäre eine Oberfläche, die
   niemand bewusst setzt — dieselbe Begründung, mit der die Zähl-Semantik der
   Closure-Struktur ohne Ventil auskommt.

5. **SemVer: Minor.** d-check findet danach mehr; ein grüner Konsumentenlauf kann
   rot werden — dieselbe Einordnung wie
   [ADR-0042](0042-markdown-lexik-folgt-commonmark.md) und
   [ADR-0043](0043-tabellengrenze-am-relevanten-header.md). Laut Messung ist die
   erwartete Fallzahl im Ökosystem **null**, aber die Zusage ändert sich, und
   das entscheidet die Einordnung — nicht die Trefferzahl.

## Verglichene Alternativen

| Alternative | Warum verworfen |
|---|---|
| Längenabgeglichener Fence-Schluss im geteilten Automaten | Löst den belegten Fall **nicht** (ein nie geschlossener Fence bleibt offen) und ist auf 776 gemessenen Dateien wirkungslos — Vertragsfläche ohne Wirkung |
| Behandlung nur in den Closure-/Struktur-Bedingungen | Deckt den Fundort, lässt die Klasse in jedem anderen Modul stehen; dieselbe Lexik, zwei Verhaltensweisen |
| Befund im Modul `planning` | Falsch einsortiert: ein Markdown-Artefakt in einem Planning-Lifecycle-Modul; der nächste Konsument der Lexik fände es dort nicht |
| Neues Modul für Fence-Artefakte | Nicht querschnittlich im Sinne von [ADR-0044](0044-geteiltes-referenz-ventil-quell-skopus.md) — es ist genau die Frage, die `spans` schon stellt |
| Absatzweise Erkennung wie `span-unclosed` | Ein Fence *ist* eine Absatzgrenze; absatzweise gemessen wäre er per Definition nie ungeschlossen |
| Befund an der letzten Zeile der Datei | Der Ort der Reparatur ist die Öffnung, nicht das Ende — ein Befund am Dateiende zeigt auf eine Zeile, die niemand ändern muss |
| Opt-in-Schalter je Artefakt-Klasse | Eine Oberfläche, die niemand bewusst setzt; `spans` ist bereits als Ganzes opt-in |

## Konsequenzen

- **Die von [ADR-0042](0042-markdown-lexik-folgt-commonmark.md) offen gelassene
  Grenze bleibt offen** — bewusst, und jetzt mit einem Wächter davor: nicht die
  Paarung wird streng, sondern der ungepaarte Zustand wird laut.
- **Jedes Modul profitiert**, weil beide Lesarten geteilt sind: der Befund
  entsteht einmal in `spans`, die stille Verkürzung betraf aber alle. Die
  strenge Lesart wurde dafür aus dem Tabellen-Leser herausgezogen — der Wächter
  bewacht damit **die** Lesart seines Konsumenten und nicht eine zweite Kopie,
  die davon abdriften kann.
- **Die Zusage der Closure-Note-Struktur wird eingelöst**, ohne sie zu ändern:
  ihr Silent-Grün-Pfad wird von einem anderen Modul aufgedeckt. Die Grenze ist
  dabei **schärfer als zunächst notiert**, gemessen beim End-to-End-Beleg: der
  Wächter greift nur, wenn `spans` aktiv ist **und** die Datei im **Scan-Scope**
  liegt. Zwei Klassen fallen dadurch heraus: Post-Pässe über selbst benannte
  Verzeichnisse und **Zieldateien** außerhalb der Scan-Wurzeln, aus denen Module
  lesen — `matrix` den Status, `anchors` und `codepaths` die Slugs, `diagrams`
  und `versions` ihre Quellen, `pins` die gehashte Heading-Section. Bei `pins`
  ist die Folge **still**: ein offener Fence verdeckt die Überschrift, der Anker
  wird unauflösbar, und die Drift-Prüfung entfällt kommentarlos. Ebenfalls
  außerhalb liegt der eingerückte Code-Block — er ist in keiner Lesart
  modelliert. Das gehört in die Release-Notiz, nicht in eine stille Annahme.
- **Der Bestand bleibt grün** (776 Dateien, null Vorkommen). Der Wert liegt
  vollständig in der Zukunft; das ist bei einem latenten Defekt der Normalfall
  und kein Argument gegen ihn.

## Fitness Function

- **Der belegte Reproduktionsfall meldet**, und der Test ist mutations-echt: der
  Rückbau des Fixes macht ihn wieder grün.
- **Balanciert bleibt grün:** eine Datei mit vielen, unter **beiden** Lesarten
  geschlossenen Fences erzeugt keinen Befund — einschließlich einer legalen
  Verschachtelung, in der ein längerer Fence einen kürzeren zeigt.
- **Eine kippende Lesart genügt:** ein Backtick-Block, den eine Tilden-Zeile
  „schließt“, meldet trotz balancierten Toggles; ein vier Zeichen langer
  Block, der eine dreizeichige Fence-Zeile zeigt, meldet trotz schließender
  strenger Lesart.
- **Genau ein Befund je Datei.**
- **Der Wächter trimmt wie das Bewachte:** eine mit Unicode-Whitespace
  eingerückte Fence-Zeile zählt für keinen von beiden.
- **Der eigene Bestand bleibt bei null** — jede Abweichung ist ein echter Fund,
  kein Rauschen.

## Re-Evaluierungs-Trigger

- Wenn `fence-unclosed` in der Praxis Falsch-Positive erzeugt (etwa durch
  bewusst offene Fences in generierten Dokumenten), ist die Erkennung zu
  verengen — nicht abzuschalten.
- Wenn die Fence-**Paarung** eines Tages reale Fehlpaarungen zeigt (gemischte
  Längen oder `~~~` im Bestand), ist Entscheidung 1 neu zu bewerten: dann wäre
  der längenabgeglichene Schluss keine folgenlose Änderung mehr.
- Wenn ein weiteres Modul **innerhalb** eines Abschnitts misst, ist zu prüfen, ob
  die geteilte Vorverarbeitung noch weitere still verkürzende Grenzen trägt.

## Geschichte

- 2026-08-09: Proposed (doc-first, `slice-101`).
- 2026-08-09: nach unabhängigem Code-Review überarbeitet — Entscheidung 1 wertet
  **beide** Schluss-Lesarten aus (vorher nur den naiven Toggle, wodurch der
  Tabellen-Leser unbewacht blieb), Entscheidung 3 nennt die gemeldete Zeile
  ehrlich eine Fundstelle, die Fitness Function ist um die Fälle ergänzt, an
  denen die vorherige Fassung nachweislich scheiterte, und die Grenze deckt
  jetzt auch die Ziel-Achse.
- 2026-08-09: nach bestätigender Re-Review nachgezogen — die Trimmung ist zu
  einem geteilten Prädikat geworden (das Modul `planning` trimmte weiter
  unicode-weit und trug den Anlassfall des Slice unverändert), das CR einer
  CRLF-Zeile fällt aus dem Befund-Ziel, der Vorrang der strengen Lesart hat eine
  Assertion, die Fundstellen-Klausel gilt für **beide** Lesarten statt nur für
  die Parität, und die Grenze nennt `pins`, `codepaths` und den eingerückten
  Code-Block.
