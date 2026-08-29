# MR-055 — Der offene DoD-Haken wird gewächtert, der Altbestand ausgenommen (Nachtrag zu MR-049)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon macht die DoD-Häkchen zur
  **Bedingung** des Übergangs nach `done/`
  ([`modul-05-planning-harness.md`](../../.harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md)),
  sagt aber nichts darüber, **welches** Werkzeug sie hält — und nichts über den
  Umgang mit einem Bestand, der vor der Regel entstanden ist. Diese Adaption ist
  die Werkzeug-Wahl plus die Behandlung des Altbestands, keine Abweichung.
- **Datum:** 2026-08-29
- **Geltungsbereich:** [`.d-check.closure.yml`](../../.d-check.closure.yml), der
  Abschnitt `## N. Definition of Done` der Slice-Dateien unter
  `docs/plan/planning/done/`. **Nicht** der DoD-Abschnitt lebender Slices — dort
  ist ein offener Haken der Normalzustand.

## Adaption

**Eine `structure`-Regel im Closure-Profil verbietet das Task-Item-Muster
`- [ ]` im DoD-Abschnitt abgeschlossener Slices.** Sie läuft in
`make verify-closure-notes`, also am Closure-Bindepunkt und nicht im inneren
Loop — ein Slice in Arbeit trägt offene Haken, und die Regel darf ihn nicht rot
machen.

**Sie ist die zweite Instanz derselben Form wie
[`MR-049`](../conventions.md#mr-049)** — urteilsfreie Hälfte einer
Closure-Regel als `structure`-Regel, Altbestand über eine Ausnahme mit **fester
Ziffernzahl**. Damit ist die Form keine Einzelentscheidung mehr, sondern die
Haus-Antwort auf „neue Closure-Bedingung über gewachsenem Bestand".

**Der Hinweis trägt die Reparatur, nicht der Grund-Code.**
`section-forbidden` sagt „hier steht eine verbotene Wendung" — für jede
`forbid`-Regel gleich. Was zu tun ist, sagt der `hint`
([`DC-FA-STRUCT-001`](../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)):
*den Haken setzen, wenn die Zusage erfüllt ist — sonst den Slice zurückführen
und sie einlösen*. Genau dafür wurde dieser Slice einmal zurückgeführt; die
Fähigkeit entstand als eigener Schnitt.

## Die Grenze der Ausnahme ist ein Datum, kein Aufräum-Rest

**Sie endet bei `slice-171`, und das ist keine willkürliche Nummer:** es ist der
erste Slice, der nach der korrigierten Praxis geschlossen wurde. Die Ausnahme
zeigt damit nicht auf „was wir noch nicht aufgeräumt haben", sondern auf „was
vor der Regel entstanden ist" — und ihre Grenze ist überprüfbar, statt
verhandelbar.

**Drei Muster statt zwei**, weil die Treffer nicht bis zu einer runden Grenze
laufen: `slice-0??-*`, `slice-1[0-6]?-*`, `slice-170-*`. Die naheliegende
Kurzform mit `*` statt fester Ziffernzahl nähme zusätzlich vierstellige Nummern
heraus — `matchGlob` matcht segmentweise mit `path.Match`, und `*` frisst dort
den Rest. Dieselbe Klasse lautlosen Ausfalls, die
[`MR-049`](../conventions.md#mr-049) für seine Ausnahme benennt; die gewählte
Form hat sie nicht.

**Ein Unterschied zu MR-049, und er ist der Grund für den knapperen Zuschnitt:**
dort musste die Ausnahme **107** `section-missing` abfangen, weil die älteren
Slices ihren §5 unter anderem Titel führen. Hier gibt es **keinen einzigen** —
das Überschriften-Muster `^#{2,3} [0-9]+\. Definition of Done` deckt alle sieben
im Bestand vorkommenden Formen, einschließlich der Zusätze `(vorläufig)` und
`(R1 eingearbeitet)`. Die Ausnahme trägt hier also **nur** offene Haken, nichts
Formales.

## Der Anlassfall liegt in der Ausnahme, mit Absicht

**`slice-168`, `-169` und `-170` tragen genau den offenen Review-Haken, der
[welle-86](../../docs/plan/planning/welle-86-closure-uebergang-durchsetzen.md)
ausgelöst hat** — und sie bleiben ausgenommen. Ein nachträglich gesetzter Haken
behauptete einen Review, den es nicht gab; das wäre die schlechtere Lüge. Ihr
Befund steht im
[Beobachtungs-Register](../../docs/plan/planning/observations.md) und in den
Closure-Notizen der Folge-Slices, nicht in einem stillen Haken.

**Das ist die unbequeme Hälfte dieser Adaption:** der Wächter entsteht aus drei
Fällen, die er selbst nie melden wird. Wer die Ausnahme liest, ohne das zu
wissen, hält sie für Bequemlichkeit.

## Grenze

**Ein Haken ist eine Selbstauskunft.** Wer ihn setzt, ohne dass die Zusage
erfüllt ist, passiert das Gate. Der Sensor verschiebt die Lücke von
**unsichtbar** nach **behauptet** — das ist besser und es ist nicht dasselbe wie
eine Prüfung. Was die Zusage inhaltlich trägt, bleibt Urteil; die Deckung
zwischen Review-Zusage und Review-Report ist ein eigener Schnitt.

**Die Ausnahme altert.** Sie nimmt eine feste Nummernspanne heraus; jeder neue
Slice fällt unter die Regel. Wächst sie je wieder, ist das der Befund — die
Klasse, die dieses Repo als
[`BEO-013`](../../docs/plan/planning/observations.md) führt.

## Auflösungs-Trigger

**Wenn die Ausnahme wächst.** Ein zusätzlicher Eintrag hieße, dass ein Slice
nach `done/` gegangen ist, dessen Haken nicht gesetzt waren — dann ist nicht die
Ausnahme zu erweitern, sondern der Übergang zu prüfen.

**Wenn der Übergang selbst gewächtert wird.** Diese Regel prüft den **Zustand**
nach dem Move, nicht den Move. Sobald ein Wächter am `mv`-Commit hängt, ist zu
entscheiden, ob die Zustands-Regel daneben bestehen bleibt (sie fängt den
Bestand, er den Neuzugang) oder in ihm aufgeht.
