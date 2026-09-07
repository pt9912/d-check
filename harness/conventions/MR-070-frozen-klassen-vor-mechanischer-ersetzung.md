# MR-070 — Wer mechanisch über den Baum ersetzt, listet die Frozen-Klassen vorher auf

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine — der Eintrag **ergänzt** die Baseline um
  einen Vorgang, den sie nicht führt. Die nächstgelegene Regel ist
  [`grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt](../../.harness/baseline/v6.5.0/regelwerk/grundlagen-harness-dateien.md#harnessreadmemd-als-einstiegspunkt)
  (Absatz *„Ein einfrierendes Artefakt nennt ein prozess-bewegtes bei seiner
  Kennung"*). Sie regelt die **Form eines Verweises**, geschrieben vom Autor im
  Moment des Schreibens. Dieser Eintrag regelt die **Ausschluss-Menge einer
  Massen-Operation**, ausgeführt von jemand anderem lange danach. Anderer
  Adressat, anderer Zeitpunkt, andere Handlung — der Nachweis steht in
  [slice-209](../../docs/plan/planning/done/slice-209-frozen-klassen-vor-mechanischer-ersetzung.md)
  §2 und ruht auf drei Belegen, nicht auf einer Behauptung.
- **Datum:** 2026-09-07 · **Herkunft:** seit slice-209 (Steering Loop,
  `BEO-ALL/mechanical-id-rewrite-misses-frozen-classes` 3×)
- **Geltungsbereich:** jede **mechanische Ersetzung über mehr als eine Datei** —
  `sed`, ein Werkzeug-Lauf, der Token austauscht, jede Massen-Operation, die
  einen Wert durch einen anderen ersetzt. **Nicht** erfasst: die Änderung einer
  einzelnen Datei von Hand, die Wahl der Verweis-**Form** beim Schreiben (die
  regelt der Kanon) — und, ausdrücklich, der **Pfad-Nachzug eines
  Lifecycle-Moves** nach [`MR-013`](../conventions.md#mr-013); siehe die
  gemeldete Kollision unten.
- **Adaption:** Zwei Pflichten, und die erste ist die, die dreimal fehlte:

  1. **Die Ausschluss-Menge steht vor dem Lauf** — geschrieben, nicht erinnert.
  2. **Der `git diff` nach dem Lauf wird gelesen**, bevor committet wird.

  **Die Menge wird über eine Eigenschaft gebildet, nicht über Verzeichnisse.**
  Der Test lautet: *Würde ein korrigierter Wert verfälschen, was dieses Artefakt
  zu seinem Datum festgehalten hat?* Ist die Antwort ja, gehört es
  ausgenommen — unabhängig davon, wo es liegt.

  **Was heute unter den Test fällt, gemessen und ausdrücklich nicht
  abschließend:**

  | Klasse | warum eingefroren | vom Kanon genannt? |
  |---|---|---|
  | `docs/plan/planning/done/` | Lauf-Beleg zu seinem Datum | ja |
  | `docs/reviews/` | Befundstand zu seinem Datum | ja |
  | `harness/conventions/done/` | die Form, die das Repo einmal hatte | **nein** — der Kanon führt den aufgelösten `MR`-Eintrag als *aufbewahrt*, nicht als Zeitdokument |
  | `Accepted`-ADR-**Kerne** | [`AGENTS.md`](../../AGENTS.md#35-adrs-sind-nach-accepted-immutable) §3.5 | ja |
  | gesendete CRs in `docs/plan/cr/` | die Bitte zu ihrem Datum | **nein** |
  | identifizierende Nennungen in **lebenden** Dokumenten | sie benennen den Gegenstand, den die Ersetzung entfernt | **nein — der Kanon nimmt lebende Artefakte ausdrücklich aus** |

  **Drei der sechs Zeilen nennt der Kanon nicht** — der gesendete CR, der
  aufgelöste `MR`-Eintrag und die identifizierende Nennung in einem lebenden
  Dokument. **Die letzte ist der Grund für die Eigenschaft.** Eine Plan-Tabelle, die
  den zu entfernenden Baum benennt, ist ein lebendes Dokument und wird trotzdem
  falsch, wenn die Ersetzung sie mitnimmt: Danach steht derselbe Pfad zweimal
  da, als neu und als entfallend, und **beide lösen auf** — kein Gate meldet
  es. Wer die Menge über Verzeichnisse bildet, findet diesen Fall nie.


  **Gemeldete Kollision mit [`MR-013`](../conventions.md#mr-013) —
  [`AGENTS.md`](../../AGENTS.md#1-was-diese-datei-ist) §1 verlangt die Meldung,
  nicht die stille Auflösung.** Ein Lifecycle-Move zieht die Pfad-Verweise auf
  den bewegten Slice **auch in `done/`-Dateien** nach; genau das ist die
  `MR-013`-Ausnahme, und der Beanspruchungs-Commit dieses Slice hat es an sechs
  Stellen getan. Die erste Fassung dieses Eintrags führte *„`git mv` mit
  Nachzug"* im Geltungsbereich und nahm zugleich `done/` in Zeile 1 aus — wer
  das literal befolgte, ließe die Verweise stehen, und `links` meldete
  `target-missing`. **Der unabhängige Review hat es gefunden, kein Gate.**

  **Der Eigenschafts-Test löst den Widerspruch, und deshalb bleibt es bei zwei
  Regeln statt einer Ausnahme.** Ein Pfad-Nachzug verfälscht **nicht**, was das
  Artefakt festgehalten hat: Es nannte einen Slice, und es nennt ihn weiter —
  nur an seinem neuen Ort. Eine Tag- oder Kennungs-Ersetzung verfälscht sehr
  wohl: Sie behauptet, das Artefakt habe damals einen Wert gemeint, den es nicht
  gab. **Die Unterscheidung ist nicht *ob* geschrieben wird, sondern *ob die
  Aussage danach noch dieselbe ist*.** Der Geltungsbereich oben nimmt den
  Nachzug deshalb ausdrücklich aus, statt Zeile 1 aufzuweichen.

  **Und die Grenze der Auflösung:** Sie hält, solange der Nachzug wirklich nur
  Pfade berührt. Ein Lifecycle-Move, der zugleich eine **Kennung** ersetzt,
  fällt unter beide Regeln — dann gilt dieser Eintrag, und die
  `MR-013`-Bündelung entbindet nicht davon. Ein Fall dieser Art ist im Bestand
  **nicht** gemessen; die Aussage ist Ableitung, kein Beleg.

  **Was diese Regel nicht leistet, ausgeschrieben:**

  - **Kein Gate.** Ob jemand seine Ausschluss-Menge vor dem Lauf geschrieben
    hat, ist eine Aussage über einen **Akt**, nicht über einen ruhenden
    Zustand — dieselbe Lage wie bei
    [`AGENTS.md`](../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden)
    §3.6. Es gäbe nichts zu bauen.
  - **Die Tabelle oben altert.** Sie ist eine Momentaufnahme des Bestands, kein
    Ersatz für den Test darüber; eine neue Artefakt-Klasse steht am Tag ihrer
    Einführung nicht darin. Wer die Tabelle abarbeitet statt den Test
    anzuwenden, hat die Regel nicht befolgt, sondern ihre Beispiele.
  - **Drei Anlässe sind keine Inventur**
    (`BEO-ALL/rule-drawn-from-occasion-not-inventory`, 7×). Alle drei sind
    Doku-Migrationen — eine Register-Formatmigration und zwei Pin-Hebungen. Ob
    die Regel für eine **Code**-weite Ersetzung trägt, ist unbelegt.
  - **Sie verschiebt den Fehler von *unsichtbar* nach *vermeidbar*, nicht nach
    *unmöglich*.** Gefangen wurde er dreimal von Pflicht 2 oder vom
    unabhängigen Review — **nie** von einer Vorab-Liste, denn es gab keine. Die
    erste Pflicht ist damit die unbelegte von beiden, und das gehört hier
    hingeschrieben statt in den Bericht danach.
  - **Sie deckt nicht, was der Kanon selbst mitbringt.** Die vom Kanon
    **vorgeschriebene** Form für eine Baseline-Stelle — Tag und Pfad in
    Inline-Code — ist genau das Muster, das eine Tag-Ersetzung greift; der
    abgeratene Link würde von einem Sensor gemeldet. Das ist eine Beobachtung
    über die Baseline, kein lokaler Mangel, und sie ist hier **benannt, nicht
    aufgelöst**.
- **Begründung:** Dreimal wurde mechanisch über den Baum ersetzt, dreimal traf
  es eingefrorene Artefakte, und dreimal fiel es erst **nach** dem Lauf auf.
  Beim ersten Mal waren die drei benannten Frozen-Verzeichnisse korrekt
  ausgenommen und zwei weitere Klassen nicht; beim dritten gab es überhaupt
  keine Liste. Der gemeinsame Nenner ist nicht ein vergessenes Verzeichnis,
  sondern eine **falsche Abstraktion**: Die Verzeichnis-Liste ist bequem und
  unvollständig, die Eigenschaft ist unbequem und vollständig.
- **Auflösungs-Trigger:** die Baseline regelt den **Vorgang** — nicht die
  Verweis-Form, sondern die Massen-Ersetzung selbst; dann gilt ihre Fassung.
  Bis dahin permanent.
