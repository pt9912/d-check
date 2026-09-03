# Antwort des Kurses auf den Konsumenten-CR vom 2026-09-03

**Absender:** `ai-harness-course` (Herausgeber der Baseline)
**Datum:** 2026-09-03
**Bezug:** [Konsumenten-CR vom 2026-09-03](2026-09-03-cr-ai-harness-course-wellenlose-review-archivierung.md),
Zeitdokumente-Archivierung ohne Träger im wellenlosen Betrieb
**Ergebnis:** **angenommen** — Punkt 1 in der von euch beantragten Form
**Erwartete Einordnung:** MINOR — `v5.19.0` (Erstfassung), Korrektur in
`v5.20.0`. **Pin-Empfehlung ist inzwischen `v6.0.0`** (siehe *Umfang und
Auslieferung*).

**Fassung 2 (2026-09-03, nach Auslieferung von `v5.19.0`).** Die Erstfassung
lehnte eure sechste Tabellenzeile ab und nannte das Archivieren im wellenlosen
Betrieb eine Lücke ohne Auslöser. **Das war falsch.** Die Zeile kommt — mit
einem echten Träger. Was sich gegenüber Fassung 1 ändert, steht unter *Punkt 1*
und im Abschnitt zum Sammel-Archiv; der Befund und alles Übrige bleiben.

**Fassung 3 (2026-09-03, nach Auslieferung von `v6.0.0`).** An eurem Punkt
ändert sich nichts — er ist seit `v5.20.0` erledigt und bleibt es. Geändert hat
sich, worauf ihr pinnen solltet: Am selben Tag ist ein **MAJOR** nachgekommen,
der das Beobachtungs-Register umbaut. Der neue Abschnitt *Was `v6.0.0` für euch
bedeutet* steht am Ende; der Rest dieses Dokuments ist unverändert.

---

## Vorab: der Befund war schärfer als der Antrag

Ihr habt eine fehlende Zeile gemeldet. Gefunden habt ihr einen **Widerspruch
innerhalb eines Moduls**.

`modul-06` §*Wann Arbeit eine Welle braucht* sagte über das, was ohne
wellenlosen Ersatz-Träger offen bleibt:

> **Was offen bleibt.** Das ist genau eine Sache: die **Carveout-Frist**.

Dieselbe Datei, Wellen-Closure-Prozedur Schritt 4, unter *Drei Grenzen, benannt
statt überspielt*:

> Und in einem Repo **ohne** Wellen-Betrieb fehlt der Auslöser ganz: Dort gibt
> es keine Closure, an der das Archivieren hängen könnte, und die Frage bleibt
> offen.

Zwei Abschnitte eines Moduls, uneins über die Zahl der eigenen Lücken. Die
Lücke war also nicht unbenannt — sie war an der falschen Stelle benannt, und
die richtige Stelle behauptete zusätzlich, es gebe sie nicht. Für ein wellenlos
arbeitendes Repo ist das die schlechtere Hälfte des Fehlers: Wer nachschlägt,
schlägt in §*Wann Arbeit eine Welle braucht* nach und liest dort eine
abschließende Aufzählung, die keine ist.

Eure Messung trägt den Punkt unabhängig davon: 45 bewusst wellenlose Slices mit
57 Review-Reports, denen nach dem Wortlaut kein Archivierungs-Ereignis
zusteht, gegen 85 wellengebundene Wellen, deren Zeitdokumente vollständig
archivierbar waren. Das ist der Nachweis, dass der wellenlose Betrieb im
Dauerbetrieb ankommt und nicht als Randfall.

## Punkt 1 — angenommen, als sechste Tabellenzeile

Ihr habt die Zeile beantragt, und ihr bekommt sie. Nur trägt sie nicht, was ihr
vorgeschlagen habt („kein automatischer Träger, Repo-Entscheidung"), sondern
einen **echten Träger**:

| Vorgang | Träger im Repo **ohne** Wellen | Wann |
|---|---|---|
| **Zeitdokumente archivieren** (Closure-Schritt 4) | Slice-Closure | **nach** den Paarungen — sie lesen den Volltext in `done/`, den das Archiv dort schließt. Schlüssel ist der Slice: `done/slice-<NNN>-archiv.zip`, **flach** neben dem Stub |

**Warum das trägt — und warum Fassung 1 das Gegenteil behauptete.** Wir hatten
*„ein Trigger, der nichts beobachtet, ist Zeremonie"* gegen einen
Ersatz-Auslöser ins Feld geführt. Der Satz steht aber unmittelbar unter dem
Wellen-Kriterium und beantwortet genau eine Frage: **ob eine Welle vorliegt.**
Er ist keine allgemeine Regel, dass jeder Vorgang einen eigenen beobachtenden
Auslöser braucht — wir haben ein Wellen-Kriterium auf eine
Archivierungs-Operation angewandt, wo es nichts zu suchen hat.

Und die Erlaubnis zu archivieren hängt gar nicht an der Welle.
`grundlagen-traceability.md` begründet sie so: Der Volltext eines geschlossenen
Slice *„kommt in keinem lesenden Knoten vor [...] genau deshalb"* darf er ins
Archiv. Das *genau deshalb* zeigt auf **geschlossen**, nicht auf *Welle
geschlossen*. Die Eigenschaft entsteht mit der Slice-Closure; die Welle trägt
zur Rechtfertigung nichts bei — sie bündelt, mehr nicht.

Damit stand der Auslöser die ganze Zeit in derselben Tabelle, in der wir seine
Abwesenheit behauptet haben: Vier ihrer fünf Zeilen hängen ihren Wellen-Vorgang
an die Slice-Closure um. Die sechste tut jetzt dasselbe.

**Der Schlüssel ist flach, und das ist kein Geschmack.** Schritt 4 warnt selbst
davor, dass ein Archiv im Unterverzeichnis Sensoren blendet, die auf
`done/*.md` keilen. Flach bleibt der Stub auf `done/slice-<NNN>-*.md` liegen —
dort, wo die **Lage-Prüfung** der Register-Paarung ihn sucht, und dort trägt er
die überlebenden Kennungen, die das Register-Gate liest. Die Falle entsteht gar
nicht erst. Zu bündeln gibt es zudem nichts: Ein Wellen-Archiv fasst mehrere
Slices und einen Welle-Plan unter einen Schlüssel, ein Slice ist schon einer.

## Punkt 2 — angenommen wie beantragt

`modul-10` §*Reviewer berichtet auch, was er nicht gefunden hat* band die
Archivierung unbedingt an die Welle-Closure. Sie sagt jetzt ausdrücklich, dass
dieses Ereignis in einem Repo ohne Wellen nie eintritt, dass der Report dann
liegen bleibt und dass der Kurs keinen Ersatz-Auslöser vorschreibt — mit
Verweis auf die Stelle, an der die übrigen wellenlosen Träger stehen.

Ein Halbsatz, wie beantragt. Er nennt allerdings **nicht** das Sammel-Archiv
als zulässige Option, und das ist die zweite Abweichung.

## Warum das Sammel-Archiv trotzdem nicht die Antwort ist

Ihr habt die Lösungsform an der Nachbarstelle richtig erkannt — Schritt 4 kennt
das Sammel-Archiv für den Bestand vor der Einführung, und die chronologisch
nächste Welle scheidet für wellenlose Slices zu Recht aus, weil sie die
Wellenlosigkeit rückwirkend verfälschte. Nur überträgt sich die zweite Option
ebenfalls nicht, und zwar aus einem Grund, der im Antrag nicht vorkommt.

Der Altbestand ist eine **einmalige, abgezählte Menge mit bekanntem Ende**. Er
braucht ein Ziel, weil sein Auslöser schon feststeht: der bewusste, einmalige
Vorgang, den Schritt 4 beschreibt. Wellenloser Dauerbetrieb hat kein Ende. Ein
Sammel-Archiv dafür bräuchte einen **wiederkehrenden** Auslöser — und genau der
fehlt. Was der wellenlose Betrieb vermissen lässt, ist nicht das Ziel des
Umzugs, sondern sein Moment.

Erfinden musste ihn am Ende niemand — das ist der Punkt, den Fassung 1
übersehen hat. Die **Slice-Closure** ist ein vorhandenes, beobachtbares
Ereignis und trägt in derselben Tabelle bereits vier andere Vorgänge; der
wiederkehrende Auslöser war da, er stand nur nicht in dieser Zeile.

Beim Sammel-Archiv bleibt es also bei der Ablehnung, aber aus dem schwächeren
Grund: nicht weil es keinen Auslöser gäbe, sondern weil ein besserer Träger
danebenstand.

## Umfang und Auslieferung

Eine Welle, beide Punkte zusammen: Sie korrigieren dieselbe Aussage, und eine
davon allein zu registrieren hieße, den Widerspruch halb stehen zu lassen.

Im gelieferten Umfang ändern sich **drei** Regelwerk-Dateien —
`modul-06-roadmap.md`, `modul-10-review-harness.md` und
`modul-05-planning-harness.md` — sowie die `Stand:`-Zeile der Regelwerk-README
auf Kurs-Welle 113.

Die dritte stammt nicht aus eurem CR, sondern aus dem Review dieser Welle:
`modul-05` behauptete unbedingt, `done/` sei *„auch nicht die letzte Station
der Datei"*. Ohne Wellen ist es genau das. Derselbe Fehlerbau wie bei
`modul-10`, nur schwächer — und ihn stehen zu lassen, während der von euch
gemeldete korrigiert wird, wäre die halbe Korrektur.

Kein Gate, kein Werkzeug, keine Vorlagen-Änderung. Templates, Beispiel und
Lösungen sind unberührt.

**Zwei Releases, weil die Erstfassung schon draußen war.** `v5.19.0` (Wellen
112–113) trägt den Befund und die zu schwache Antwort. Die Korrektur — die
sechste Tabellenzeile mit ihrem Träger — kam als **Welle 114** in `v5.20.0`.
`v5.19.0` sagt an drei Stellen, das Archivieren habe im wellenlosen Betrieb
keinen Träger, und das stimmt nicht — **überspringt es.**

**Pinnt auf `v6.0.0`, nicht auf `v5.20.0`.** Am selben Tag ist ein MAJOR
nachgekommen. Auf `v5.20.0` zu pinnen hieße, zwei Migrationen kurz
hintereinander zu fahren.

## Was `v6.0.0` für euch bedeutet

Nicht Gegenstand eures CR, aber ihr seid davon betroffen, und zwar mehr als
von ihm.

**Das Beobachtungs-Register ist keine flache Tabelle mehr.** Je Beobachtung ein
Verzeichnis `observations/BEO-<KUERZEL>/<slug>/` mit `observation.md`
(unveränderlich), `state.md` (veränderlich) und `evidence/<vorgangs-id>.md`
(unveränderlich, **eine je Auftreten**). Der Zähler wird aus den
Evidence-Dateien **abgeleitet** statt geführt.

**Der Grund ist nicht Teamfähigkeit** — die fällt an. Der Grund ist, dass
Zähler *und* Belegliste zwei Quellen für denselben Zustand waren. Gemessen an
fremden Registern mit je **einem** Schreiber: eines stand auf 7 statt 6, ein
Eintrag auf 3 mit einem Beleg. Die Anzahl-Prüfung aus Welle 106 gab es nur,
weil der Zähler gespeichert wurde; sie ist jetzt ersatzlos entfallen.

**Drei Dinge, die das für euch heißt:**

1. Euer `observations.md` migriert in die Verzeichnisform. Für ein Repo mit
   einem Schreiber ist das mechanisch: eine Zeile → ein Verzeichnis, ein Beleg
   → eine Evidence-Datei. Die Gegenprobe ist hart — der **abgeleitete** Zähler
   muss die alte Spalte reproduzieren. Tut er es nicht, war die alte Spalte
   falsch, und das ist dann der eigentliche Fund.
2. Die **Kürzel-Spalte** eurer Modus-Deklaration ist nicht mehr optional. Ihr
   habt in eurer `conventions.md` ausdrücklich erklärt, warum ihr keine führt —
   das war korrekt, solange keine Kennung ein Segment trug. Jetzt trägt eine
   eins: der Pfad der Beobachtung. Die Bedingung von Welle 105 ist damit
   **erfüllt, nicht aufgehoben**.
3. **Was ihr aufgeschoben habt, ist jetzt der Kanon.** Auf unseren CR zur
   relationalen BEO-Ablage habt ihr angenommen und aufgeschoben — §1–§5 nach
   unseren zwei Quell-Wellen. Die sind seit `v5.15.0` draußen, und mit `v6.0.0`
   ist die Form nicht mehr Vorschlag, sondern Regelwerk. Die zwei
   Sensor-Fähigkeiten, die dabei offen blieben — die 3×-Schwelle an eine Aktion
   koppeln und Alias-Ketten beim Zählen auflösen —, haben damit einen
   Konsumenten in der Baseline und nicht mehr nur in einem Entwurf.

Was **nicht** in den Kanon gewandert ist, ist unser Entwurf als Ganzes: Der
Kurs verlangt die drei Eigenschaften (ein Ort je Beobachtung, ein Beleg je
Auftreten, abgeleiteter Zähler) und die Pfad-Grammatik. `alias-of`,
`invalidations/` und die erzeugte Übersichts-Sicht sind **nicht** normativ
geworden. Wer sie braucht, baut sie als Repo-Entscheidung.

**Was ihr jetzt tun könnt — nicht müsst.** Eine Pflicht schafft die Antwort
nicht; eure 57 Review-Reports dürfen liegen bleiben, und euer Repo ist so
konform. Wer sie loswerden will, hat jetzt aber Träger, Moment und Schlüssel
benannt statt „eure Entscheidung": Slice-Closure, nach den Paarungen,
`done/slice-<NNN>-archiv.zip`. Euer `tools/archive-wave` braucht dafür einen
wellenlosen Modus — dieselbe Operation, Schlüssel Slice statt Welle.

**Und eine Bitte in eigener Sache.** Fassung 1 hat euren Punkt kleingeredet und
ist als `v5.19.0` ausgeliefert worden. Ihr hattet in der Sache recht, und mit
der Form — der sechsten Zeile — ebenfalls. Was nicht trug, war allein der
Inhalt, den ihr für sie vorgeschlagen habt.
