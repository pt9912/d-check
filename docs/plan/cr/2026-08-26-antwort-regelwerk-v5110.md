# Antwort des Kurses auf den Konsumenten-CR vom 2026-08-25

**Absender:** `ai-harness-course` (Herausgeber der Baseline)
**Datum:** 2026-08-26
**Bezug:** [Konsumenten-CR vom 2026-08-25](2026-08-25-cr-regelwerk-v5110.md),
vier Punkte gegen Regelwerk `v5.11.0`
**Ergebnis:** alle vier Punkte **angenommen**; Punkt 1 mit anderer Encodierung
**Erwartete Einordnung:** MINOR — `v5.12.0`

**Fassung:** enthält die drei vom Absender nachgereichten Korrekturen (Punkt 2
Zählung, Punkt 2 Gegenbeleg-Absatz, Abschnitt *Umfang und Auslieferung*). Ein
angrenzender Halbsatz in Punkt 2 (*„plus Spiegel"*) ist mitgezogen, weil die
korrigierte Lieferliste den Beispiel-Baum ausnimmt.

---

## Vorab zur Form

Dass keiner der vier Punkte ein Gate verlangt, war die richtige Konsequenz aus
der vorigen Runde. Die Belege sind an eigenem Bestand geführt und im Zweifel
gegen euch selbst gelesen — die drei Fallen unter Punkt 1 und die drei
Fehlläufe unter Punkt 3 sind der Grund, warum die Punkte tragen. Sie wandern
trotzdem nicht in den Kanon: ein Beleg begründet eine Regel, er wird nicht Teil
von ihr.

## Punkt 1 — angenommen, mit anderer Formulierung

Die Lücke ist bestätigt. `modul-06` trennt an dieser Lage Urteil von Deckung,
`modul-05` §*Offene Risiken werden bei Closure aufgelöst* tut es nicht — im
ganzen Abschnitt steht weder „Urteil" noch „maschinell". Der Parallelfall ist
sauber gezogen.

**Was wir anders machen.** Der Satz wird nicht *„ein Slice in `done/` trägt
keinen unaufgelösten Vorlagen-Platzhalter"* lauten. Zwei Gründe:

1. **Referenz-Richtung.** Der Satz stünde im Kursmodul und setzte damit die
   Platzhalter-Form von `slice.template.md` voraus. In diesem Korpus ist die
   Quelle der Anker und die Vorlage abgeleitet; eine Regel, die die Form ihrer
   eigenen Ableitung zur Bedingung macht, dreht die Richtung um. Die Vorlage
   ist zudem illustrativ und nicht normativ — ein Adopter, der sie nicht
   benutzt, fiele aus der Regel heraus.
2. **Reichweite.** Der Platzhalter ist der Träger, nicht die Wirkung.
   Urteilsfrei prüfbar ist: dass zu jedem notierten Risiko ein Ausgang dasteht
   und **welcher der drei** es ist — Form, kein Freitext. Ob der eingetragene
   Ausgang inhaltlich trägt, bleibt Urteil, da sind wir einig. Der
   stehengebliebene Platzhalter ist ein Fall von „kein Ausgang", kein eigener
   Regelinhalt.

**Was das für euer Target heißt: nichts Einschränkendes.** Es prüft einen
Spezialfall der Regel und bleibt gültig. Wenn ihr es breiter zieht — *Ausgang
fehlt ganz* · *Ausgang ist Freitext statt einer der drei Formen* —, deckt es
die urteilsfreie Hälfte vollständig ab, statt nur ihren häufigsten Auslöser.

Eure dritte Falle (der Fehler-Fallback hinter der Pipe, der ein Leseversagen zu
Exit 0 verschluckt) ist im Kanon bereits benannt, nur an anderer Stelle:
`modul-13` §*Guard-Härtung* führt das Fail-Closed-Verhalten als
Design-Eigenschaft des Wächters. Dass ihr es am eigenen Wächter geprobt habt,
ist genau der dort gemeinte Griff.

## Punkt 2 — angenommen; die Lesart ist Code-Modul-Pfad

Wir haben nachgezählt: genau die zwei Fundstellen, die ihr nennt. Keine dritte,
keine andere Schreibweise — eure Suche stimmt.

Gemeint sind Pfade zu **Code**-Modulen (`internal/service/…`), nicht
Dokument-Pfade. Der Aufzählungspunkt trägt zwei Eigenschaften, und der Satz
bedient beide: *„referenziert Modul-Pfade"* trägt die sprachfreie Hälfte — ein
Pfad statt eines Sprachkonstrukts —, *„aber keine Wellen, Slices, Commit-Hashes
oder Closure-Daten"* die meilensteinfreie. Gegen „Wellen" hätte ein
Dokument-Pfad keinen Kontrast; die Achse des Satzes ist räumlich gegen
zeitlich.

Dass `architecture.template.md` keine Code-Pfade führt, ist kein Gegenbeleg.
Der Satz ist asymmetrisch: der erste Teil erlaubt, der zweite verbietet. Die
Vorlage übt die Erlaubnis schlicht nicht aus — sie benennt Rollen, weil das für
ihren Zweck reicht. Aus einer nicht ausgeübten Erlaubnis folgt kein Verbot. Ihr
habt das selbst richtig eingeordnet: eine Vorlage ist kein Regeltext.

**Für euch heißt das: eure Annahme war richtig.** Ein Repo, das seiner Sicht
Code-Pfade verbietet, **verschärft** — der deklarierte Adaptions-Eintrag steht
und braucht keine Korrektur. Wir setzen das klärende Wort an beiden Stellen.

## Punkt 3 — angenommen wie beantragt

Der stärkste der vier. Die Präzisierung des sechsten Schritts kommt, und zwar
in beiden Hälften: **das Rot muss von der gebrochenen Regel kommen, und seine
Ursache gehört gelesen.**

Eine Ergänzung, die euch nützt: der Kanon führt die Gegenrichtung bereits.
`modul-11` §*Schritt 8* verlangt *„je Verstoßklasse ein Break-Test … plus der
unveränderte Bestand, auf dem beide schweigen müssen"* — das ist die
**Negativkontrolle** (der Wächter darf nicht grundlos rot werden). Was fehlte,
ist die **Positivkontrolle** (er muss aus dem richtigen Grund rot werden). Eure
drei Fehlläufe verteilen sich genau auf diese zwei Richtungen: Fall 1 und 2
sind Positivkontrolle, Fall 3 ist Negativkontrolle mit umgekehrtem Vorzeichen.
Wir schreiben den fehlenden Satz in `modul-13` Schritt 6 und verbinden ihn mit
der vorhandenen Stelle, damit die zwei Richtungen als Paar lesbar sind.

## Punkt 4 — angenommen; ein Satz, verortet in `grundlagen-source-precedence.md`

Die Lücke ist bestätigt: wir haben den Korpus auf eine allgemeine
Zitat-Disziplin durchsucht und keine gefunden. Die zwei Stellen, die ihr
zitiert, beantworten die Reichweitenfrage tatsächlich nur je für sich selbst.

**Verortung**, weil ihr sie offengelassen habt: der Satz landet in der
Nachbarschaft der zwei zitierten Stellen, nicht in
`grundlagen-referenz-richtung.md`. Die regelt die **Richtung** (welche
Artefakt-Klasse welche referenzieren darf) und ist strikt als Matrix gebaut;
ein Reichweiten-Satz wäre dort ein Fremdkörper und würde als Zeile der Matrix
missverstanden. Reichweite ist eine Eigenschaft der zitierten **Aussage**,
nicht der Referenz-Beziehung.

**Was der Satz encodiert, ist die Frage, nicht der Katalog.** Von euren drei
Formen ist die mittlere bereits geregelt und wurde nur überlesen — dafür
schreiben wir keine zweite Regel. Die anderen beiden sind Instanzen derselben
Frage, und die gehört gestellt, nicht aufgezählt. Ein Katalog von Verstoßformen
würde die vierte Form nicht abdecken, die Frage schon.

## Umfang und Auslieferung

Vier Wellen, eine je Punkt — ihr habt sie ausdrücklich einzeln annehmbar
gehalten, und so werden sie auch registriert. Im gelieferten Umfang ändern sich
vier Regelwerk-Dateien und bei Punkt 2 zusätzlich `AGENTS.template.md`.

Kein Gate, kein Werkzeug, keine Vorlagen-Änderung — auch nicht bei Punkt 1.

Erwartete Einordnung: **MINOR (`v5.12.0`)**. Der Baseline-Pin ist danach
nachzuziehen.

**Was ihr jetzt nicht tun müsst:** die Verschärfungs-Deklaration aus Punkt 2
bleibt korrekt. Euer Closure-Target aus Punkt 1 bleibt gültig; ob ihr es auf
die volle urteilsfreie Hälfte zieht, ist eure Entscheidung und kein
Konformitätsthema.
