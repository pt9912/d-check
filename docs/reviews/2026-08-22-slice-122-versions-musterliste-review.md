# Review-Report: slice-122 — `versions`: mehrere Muster-Quellen-Paare statt eines

**Datum:** 2026-08-22 · **Review-Art:** Code- und Design-Review (Lastenheft/Spezifikation/ADR gegen Implementierung), unabhängiger Reviewer ohne Anteil an der Arbeit
**Gegenstand:** die fünf Slice-Commits der Kette (Beanspruchung · Lastenheft-CR · ADR-0058 · Spezifikation · Implementierung); der Prüfbaum wurde per `git archive` außerhalb des Repos ausgepackt, am Repo hat der Reviewer nichts geändert
**Skill:** `.harness/skills/reviewer.md` @ 1.8.0 · **Modell-ID:** `claude-opus-5[1m]`
**Eingangs-Kontext:** [`DC-FA-VER-001`](../../spec/lastenheft.md#dc-fa-ver-001--versions-pin-konsistenz-modul-versions-opt-in) (Lastenheft 0.63.0), [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus), [`DC-FA-VER-001.a`](../../spec/spezifikation.md#dc-fa-ver-001a--versions-pin-konsistenz-versions) samt [`SPEC-005`](../../spec/spezifikation.md#spec-005--d-checkyml), [ADR-0058](../plan/adr/0058-konfigurations-flaechen-additiv-weiten.md), `AGENTS.md` §3/§5, [Slice-Plan](../plan/planning/in-progress/slice-122-versions-musterliste.md), [Wellendokument](../plan/planning/welle-82-config-flaechen.md)

**Verdikt des Reviews: blockierend** — ein HIGH, vier MEDIUM, drei LOW, drei INFO. Alle acht Befunde sind eingearbeitet; die drei INFO sind zwei behoben, einer ist hier festgehalten.

## Befunde und Einarbeitung

### F-1 · HIGH — die zweite Erwartung verschwand still

`internal/hexagon/core/rules/versions.go` i. V. m. `internal/hexagon/core/model/finding.go`

Zwei Paare, die auf **derselben Zeile denselben Pin-Wert gegen verschiedene
erwartete Versionen** treffen, erzeugten **einen** Befund — die Erwartung des
zweiten Paares ging verloren. Die Befund-Adresse der geteilten Nachrunde ist
`(Datei, Zeile, Regel, target, Grund-Code)`; die `message`, in der allein die
erwartete Version steht, gehört nicht dazu. Damit war das Akzeptanzkriterium
„erwarten sie verschiedene Versionen, then **zwei** Befunde" nicht erfüllt, und
die tragende ADR-Begründung („was die Paare unterscheidbar macht, ist die
erwartete Version — sie steht bereits in der Nachricht") beruhte auf einer
falschen Prämisse: sie steht dort, aber sie zählt dort nicht.

**Messung des Reviews:** zwei Paare mit überlappenden Mustern auf der Zeile
`werkzeug:v1.0.0` ⇒ ein Befund, überlebende Nachricht „erwartet v0.27.0"; die
Erwartung v5.9.0 fehlte vollständig. Das Modul erzeugte beide, die Sortierung
verwarf den zweiten.

**Eingearbeitet — nicht die Adresse erweitert, sondern die Nachricht
vollständig gemacht:** je Zeile entsteht höchstens ein Befund pro gefundenem
Pin-Wert, und seine Nachricht nennt **jede** Erwartung mit ihrer Quelle, in
Deklarationsreihenfolge. Die Alternative — ein Adress-Feld für die Quelle —
läge in [`SPEC-001`](../../spec/spezifikation.md#spec-001--befund) und beträfe
jedes Modul; sie ist in der ADR als verworfene Alternative samt
Re-Evaluierungs-Trigger festgehalten, nicht verschwiegen. Nachgemessen an
derselben Fundstelle: ein Befund, beide Erwartungen benannt.

Die Klasse ist im Repo bereits zweimal belegt — `structure` und `planning`
umgehen dieselbe Enge über eigene Grund-Codes und sagen das in ihren
Kommentaren. `versions` folgte ihr nicht; jetzt ist die Enge benannt statt
umgangen.

### F-2 · MEDIUM — der Test bestätigte eine Zusage, die schon vorher galt

Der eingebaute paar-lokale Dedup hatte **keine beobachtbare Wirkung**: sein
Schlüssel war feiner als der der nachgelagerten Sortierung, die ohnehin alles
kollabierte. Der Reviewer hat die Map ersatzlos entfernt und die gesamte Suite
laufen lassen: **grün**. Eine Mutation, die niemand bemerkt, ist kein
gedeckter Vertrag.

**Eingearbeitet:** die Bündelung ist jetzt der Mechanismus, der die
Vollständigkeit der Nachricht herstellt — und damit beobachtbar. Gegenprobe
nachgeholt: Bündelung ausgeschaltet (kompilierbar mutiert) ⇒ **zwei** Tests
rot; wiederhergestellt ⇒ grün.

### F-3 · MEDIUM — eine Reihenfolge, die ein nachgelagerter Sort überschreibt

Die Zusage „zwei Befunde **in Deklarationsreihenfolge**" war als beobachtbare
Eigenschaft falsch: die Ausgabe ist nach `(Datei, Zeile, Regel, Ziel, Grund)`
sortiert, also nach dem Pin-Wert. Der pinnende Test hatte Ziele gewählt, deren
Sortierreihenfolge zufällig der Deklarationsreihenfolge entsprach — er konnte
die beiden Ordnungen nicht unterscheiden. Die ADR-Begründung „ohne diese Sätze
hinge der Befundsatz an einer Map-Iteration" traf zusätzlich nicht zu: die
Paare stehen in einem Slice.

**Eingearbeitet:** Lastenheft, Spezifikation und ADR sagen jetzt, was gilt —
die **Ausgabe**-Reihenfolge ist die der geteilten Sortierung; die
Deklarationsreihenfolge entscheidet, in welcher Folge die Erwartungen **in der
Nachricht** stehen. Der Test ist ersetzt durch einen, der die beiden Ordnungen
trennt: das zuerst deklarierte Paar findet den lexikografisch größeren Wert,
gemeldet wird der kleinere zuerst.

### F-4 · MEDIUM — Anwesenheit gegen Wert bei der Mischform

Die Mischform-Erkennung fragte **Werte** ab. `pin-pattern: ""` neben `patterns`
zählte nicht als gesetzt, und die Prüfung lief still mit der Liste weiter —
genau das, was die Anforderung ausschließt („keine der beiden Schreibweisen
gewinnt still"). Gemessen: vier Fälle, zwei davon fälschlich Exit 1.

**Eingearbeitet:** die drei Kurzform-Schlüssel sind Zeiger, die
Mischform-Erkennung hängt an der **Anwesenheit**. Vier neue Config-Rand-Fälle
gepinnt (leerer String, leere Liste, leere `patterns`-Liste, Vollform). Die
Rest-Grenze ist benannt statt geschlossen: ein Schlüssel **ohne Wert** ist im
YAML von einem fehlenden nicht unterscheidbar und zählt als fehlend — das
steht jetzt in Anforderung, Spezifikation und ADR.

### F-5 · MEDIUM — ADR ohne Re-Evaluierungs-Trigger

`AGENTS.md` §5 verlangt die Sektion; ADR-0048 bis ADR-0057 tragen sie
durchgehend, ADR-0058 als einzige neuere nicht. Da das Feld im ADR-Core liegt,
wäre es nach dem Übergang auf `Accepted` **nicht nachrüstbar**, ohne
`make adr-check` zu brechen.

**Eingearbeitet:** vier Trigger ergänzt, solange die ADR `Proposed` ist —
darunter der, der aus F-1 folgt (ein dritter Konsument, der zwei Befunde an
derselben Adresse braucht, macht die Befund-Adresse selbst zur Frage).

### F-6 · LOW — Zeiger auf einen Schlüssel, den es in dieser Schreibweise nicht gibt

Befund-Nachricht und fail-closed-Meldungen nannten `versions.current-from` —
einen Schlüssel, den die Listen-Form nicht kennt. Der Config-Rand benannte das
Paar bereits, die Lauf-Zeit-Hälfte zog nicht nach.

**Eingearbeitet:** der Wortlaut hängt an der **Paar-Zahl**, nicht an der
Schreibweise: bei genau einem Paar der Kurzform-Schlüssel (byte-identisch), ab
zwei Paaren `versions.patterns[i].current-from` — in der Befund-Nachricht wie
im Abbruch.

### F-7 · LOW — ein Vorlagen-Block, der einkommentiert bricht

Wer den `versions:`-Block der `--print-config`-Vorlage als Ganzes
einkommentierte — der übliche Umgang mit einer Vorlage —, bekam Exit 2, weil
Kurzform und `patterns` im selben Block standen. Kein anderer Block hat diese
Eigenschaft.

**Eingearbeitet:** zwei getrennte Blöcke, der zweite ausdrücklich als
Alternative überschrieben, seine Prosa doppelt kommentiert (damit sie das
Einkommentieren übersteht). Ein Test kommentiert **beide** Blöcke aus der
echten `--print-config`-Ausgabe ein und schickt sie durch den eigenen Parser.

### F-8 · LOW — Default-Spalte beschreibt einen verbotenen Wert

Die Schema-Zeile sagte Default „leer", der Zellentext „explizit leere Liste ⇒
Exit 2". **Eingearbeitet:** Default ist jetzt „— (abwesend ⇒ Kurzform-Pfad)".

### INFO-1 · Ausrichtung des `runState`-Literals

Nach dem Einfügen des neuen Feldes war der Block nicht neu ausgerichtet; das
Lint-Profil enthält keinen Formatter, es gab also keinen Anker. **Behoben** —
sonst zieht der nächste Editor-Save fünf unbeteiligte Zeilen in einen fremden
Diff.

### INFO-2 · Exportierte Funktion mit unexportiertem Parametertyp

`CheckVersions` nahm einen unexportierten Typ und war damit von außerhalb des
Pakets nicht mehr aufrufbar — anders als alle übrigen `Check*`-Funktionen der
Modul-Familie. **Behoben:** der Typ ist exportiert.

### INFO-3 · Zahl in der Commit-Botschaft

Die Botschaft des Implementierungs-Commits sagt „elf neue Tests", tatsächlich
waren es dreizehn. Untertreibung, keine Falschaussage; der Commit ist
veröffentlicht und wird nicht umgeschrieben. Hier festgehalten, damit die Zahl
nicht als gemessene Größe weiterwandert.

## Negativbefunde des Reviews (geprüft, ohne Befund)

Der Reviewer hat fünfzehn Punkte ohne Befund abgeschlossen; die tragenden:

- **Die Kern-Zusage hält.** Ganzer Baum mit unveränderter `.d-check.yml` gegen
  das gepinnte Vorgänger-Image: stdout byte-identisch (105 B), stderr
  identisch. Roter Fall: byte-identisch (786 B). Die Behauptung der
  Commit-Botschaft ist nachgemessen, nicht geglaubt.
- **Die Gegenprobe hält:** die Zwei-Paar-Konfiguration weist das
  Vorgänger-Image mit `field patterns not found` (Exit 2) zurück.
- **Die Meldungen der Kurzform sind stabil:** vier Config-Rand-Fälle vor und
  nach dem Umbau byteweise identisch.
- **Paar-Lokalität der Datei-Ventile, paar-übergreifender Zeilen-Marker,
  fail-closed je Paar, strikte Schlüssel auch in `patterns[]`** — je gemessen.
- **Kein neuer Grund-Code**, `--doctor`-Klartext unberührt, die beidseitige
  Kopplung zwischen Grund-Code-Test und §4 bleibt grün.
- **Das Handbuch widerspricht nicht** — es dokumentiert nur die Kurzform und
  behauptet nirgends „genau ein Muster"; die Ergänzung gehört wie geplant zum
  Release-Prep-Slice.

## Nachmessung nach der Einarbeitung

- `make gates` grün (acht Gates, Exit 0 explizit gelesen).
- Byte-Identität **erneut** gemessen, nachdem der Nachrichten-Code angefasst
  war: Fixture rot 563 B und das ganze Repo grün — beide byte-identisch gegen
  das Vorgänger-Image.
- Die F-1-Fundstelle des Reviews nachgestellt: ein Befund, **beide**
  Erwartungen mit ihrer je eigenen Quelle in der Nachricht.
- Mutations-Gegenprobe zu F-2: Bündelung ausgeschaltet ⇒ zwei Tests rot.
