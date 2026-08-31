# MR-056 — Der DoD-Haken eines geschlossenen Slice wird gewächtert

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon macht die DoD-Häkchen zur
  **Bedingung des Übergangs** und sagt das wörtlich:   <!-- d-check:cite .harness/baseline/v5.15.0/regelwerk/modul-05-planning-harness.md:33-34 -->
  „DoD-Häkchen und Closure-Notiz
  sind die Bedingung dafür, dass die Datei überhaupt nach `done/` darf."
  ([`modul-05-planning-harness.md`](../../.harness/baseline/v5.15.0/regelwerk/modul-05-planning-harness.md)).
  Gehalten war davon nur die zweite Hälfte. Diese Adaption ist die
  Werkzeug-Wahl für die erste, keine Abweichung.
- **Datum:** 2026-08-30
- **Geltungsbereich:** [`.d-check.closure.yml`](../../.d-check.closure.yml),
  der Abschnitt `## N. Definition of Done` der Slice-Dateien unter
  `docs/plan/planning/done/`. **Nicht** der Übergang selbst — diese Regel prüft
  den **Zustand am Ruheort**; die Bindung an den `mv`-Commit ist ein eigener
  Schnitt.
- **Adaption:** Eine `structure`-Regel im Closure-Profil hält
  `max-open-tasks: 0` über den DoD-Abschnitt; sie läuft in
  `make verify-closure-notes`, nicht in `gates` — sonst meldete sie beim
  Arbeiten an einem laufenden Slice.

  **Die Bedingung ist `max-open-tasks`, nicht `forbid-pattern`, und das ist
  keine Geschmacksfrage.** Die bereinigt lesende Form fiel an einem einzelnen
  überzähligen Backtick im Absatz auf **null** Befunde — gemessen — und deckte
  nur die Bullet-Form, die ihr Autor aufschrieb; `* [ ]` und `+ [ ]` liefen
  still hindurch. Die rohe Form
  ([ADR-0074](../../docs/plan/adr/0074-offene-tasks-auf-rohen-zeilen.md)) zählt
  über die Modul-Lexik und meldet **je Haken auf seiner Zeile**.

  **Der `hint` trägt die Reparatur.** `section-tasks-open` sagt die **Art** des
  Defekts; was zu tun ist — Haken setzen oder Slice zurückführen — weiß nur die
  Regel ([ADR-0073](../../docs/plan/adr/0073-befund-erlaeuterung-fuer-menschen.md)).

  **Die Bestands-Ausnahme nennt den abgeschlossenen Altbestand mit fester
  Ziffernzahl**, wie [MR-049](MR-049-ausgangs-wortschatz.md) es für die
  Drei-Ausgänge-Regel führt, und aus demselben Grund: das Regelwerk wurde
  mehrfach gehoben, und die Dokumente sind nicht durchgängig nachgezogen. Ein
  Befund auf einem Slice, der nach damaliger Form korrekt war, wäre kein Befund,
  sondern Lärm. Drei Muster statt zwei, weil die Treffer nicht in einer Spanne
  liegen: `slice-0??-*`, `slice-1[0-6]?-*`, `slice-170-*`.

  **Die Grenze ist ein Datum, kein Aufräum-Rest.** Sie endet bei
  [slice-171](../../docs/plan/planning/done/slice-171-vorpruefungen-belegen.md),
  dem ersten Slice unter der korrigierten Praxis. Gemessen halten die Slices
  171–182 ihre Haken gesetzt: **acht** Dateien liegen dort (172 ist in Arbeit, 173–175 gibt es noch nicht), und keine trägt einen offenen Haken.

  **slice-168, -169 und -170 liegen bewusst *in* der Ausnahme.** Sie tragen
  genau den offenen Review-Haken, der
  [welle-86](../../docs/plan/planning/welle-86-closure-uebergang-durchsetzen.md)
  ausgelöst hat. Ein nachträglich gesetzter Haken behauptete einen Review, den
  es nicht gab; ihr Befund steht im Register und in der Closure-Notiz von
  slice-171. Wer die Ausnahme ohne diesen Satz liest, hält sie für
  Bequemlichkeit.
- **Grenzen — gemessen, nicht geschätzt.**
  - **Ein Haken ist eine Selbstauskunft.** Wer ihn setzt, ohne dass ein Review
    stattfand, passiert die Regel. Sie verschiebt die Lücke von *unsichtbar*
    nach *behauptet* — das ist besser und keine Prüfung des Reviews.
  - **Ein Haken IM Fenced-Block ist unsichtbar, und zwar wohlgeformt.**
    Gemessen: `- [ ]` innerhalb eines geschlossenen Fences im DoD-Abschnitt ⇒
    **0 Befunde**, und auch kein `fence-unclosed` — der Fence ist ja in
    Ordnung. Das ist dieselbe Fence-Treue, die verhindert, dass ein Slice über
    Task-Items schreibt und seine eigene Illustration meldet; sie ist zugleich
    der Weg, einen offenen Haken zu verstecken. **Kein Sensor deckt das**, und
    keiner kann es: ob ein Fence Beispiel oder Versteck ist, ist ein Urteil.
  - **Der `hint` erscheint auch auf `section-missing`.** Gemessen: fehlt der
    DoD-Abschnitt oder ist er umbenannt, meldet die Regel `section-missing`
    **mit demselben Hinweis**. Deshalb nennt er die **Zusage** (*„ein
    geschlossener Slice trägt einen DoD-Abschnitt, dessen Haken alle gesetzt
    sind"*) und nicht den Defekt — so stimmt er für beide Grund-Codes.
    [ADR-0073](../../docs/plan/adr/0073-befund-erlaeuterung-fuer-menschen.md)
    nimmt nur zwei Nicht-Bedingungs-Befunde vom `hint` aus; `section-missing`
    und `section-ambiguous` gehören derselben Klasse an und sind es **nicht**.
    Ob sie es sollten, ist ein eigener Entscheid — hier nur benannt.
  - **Ein vergessener Schluss-Fence macht die BEDINGUNG blind**, und die rohe Lesung
    behebt das nicht. Isoliert gemessen: ein offener Haken hinter einem
    unbeendeten Fence ⇒ **0 Befunde, Exit 0**. Gefangen wird der Fall von
    `fence-unclosed`, also von `spans` — das Profil führt es seit
    [ADR-0077](../../docs/plan/adr/0077-spans-am-bindepunkt-die-begruendung-traegt-anders.md).
    **Diese Regel hängt damit an einer Eigenschaft des Profils, die sie nicht
    selbst herstellt.**

    **Der Bindepunkt als Ganzes wird davon nicht grün, und das gehört
    unterschieden.** Im heutigen Profil melden Nachbarregeln (etwa
    `section-empty`, weil der Fence den Abschnitt leert), und `spans` nennt die
    Ursache. Das ist aber eine Eigenschaft **dieses** Regelsatzes, nicht der
    Bedingung: fiele die Nachbarregel weg, bliebe nur `fence-unclosed`. Blind
    ist die Bedingung, nicht der Lauf.
  - **Sie prüft den Zustand, nicht den Übergang.** Ein Slice, der mit offenem
    Haken nach `done/` wandert, wird erst beim nächsten
    `make verify-closure-notes` gesehen — nicht beim `mv`-Commit.
  - **Die Ausnahme altert.** Sie nimmt eine feste Nummernspanne heraus; jeder
    neue Slice fällt unter die Regel. Wächst die Ausnahme je wieder, ist das der
    Befund ([`BEO-013`](../../docs/plan/planning/observations.md)).
- **Auflösungs-Trigger:** wenn die Bindung an den **Übergang** steht (der
  `mv`-Commit wird geprüft, nicht der Zustand danach) — dann ist die
  Zustands-Regel die schwächere Hälfte und ihre Rolle neu zu lesen.
