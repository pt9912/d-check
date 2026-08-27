# MR-051 — Ein Baseline-Bump ankert die `d-check:cite`-Spannen neu (Nachtrag zu MR-021)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon kennt den Bump-Drift-Audit und
  die Regel, dass ein zitierter Wortlaut nicht rückwirkend umgeschrieben wird
  ([`MR-039`](../conventions.md#mr-039)). Er kennt keine Zeilen-Anker, weil er
  kein Werkzeug voraussetzt, das sie prüft. Die Form-Frage tritt die Rangliste
  an diesen Speicher ab.
- **Datum:** 2026-08-27
- **Geltungsbereich:** jede `d-check:cite`-Direktive in den **lebenden**
  Dokumenten dieses Repos, deren Ziel unter `.harness/baseline/<tag>/` liegt.
  Nicht `done/`, nicht `docs/reviews/`, nicht `conventions/done/` — dort steht
  keine.
- **Adaption:** [`MR-021`](../conventions.md#mr-021) verpflichtet den
  Bump-Drift-Audit, alle vendored-Pfad-Links auf den neuen Tag zu ziehen. Eine
  `cite`-Direktive trägt **zwei** pin-gebundene Größen: den Pfad **und die
  Zeilenspanne**. Der Pfad wandert mit der Link-Regel; die Spanne muss **neu
  angekert** werden, weil sich Zeilennummern schon durch eine Änderung weiter
  oben verschieben.

  **Der Alarm ist gewollt, der Preis ist benannt.** Verschiebt sich die Spanne,
  meldet `citations` `citation-mismatch` — genau die vierte Spiegel-Klasse aus
  [`BEO-008`](../../docs/plan/planning/observations.md), die sonst still
  alterte. Der Preis: die Direktive meldet **auch dann**, wenn nur die
  Zeilennummern gewandert sind und der Wortlaut unverändert gilt. Das ist keine
  Falschmeldung — die Direktive zeigt dann tatsächlich auf die falschen Zeilen
  —, aber es ist Arbeit ohne inhaltlichen Anlass.

  **Wie unterschieden wird — drei Fälle, und der Grund-Code sagt, welcher.**
  Meldet das Modul `citation-mismatch`, existiert die Zieldatei und der
  Vergleich schlug fehl: steht der zitierte Wortlaut **anderswo in derselben
  Datei** unverändert, ist es eine reine Verschiebung ⇒ **Spanne nachziehen**,
  Zitat unangetastet. Steht er dort **gar nicht mehr**, ist es ein
  Wortlaut-Delta ⇒ [`MR-039`](../conventions.md#mr-039) greift: der Wortlaut
  bleibt stehen, die **Direktive wird entfernt** (ihre Quelle existiert nicht
  mehr), und der Bump-Eintrag hält das Delta fest.

  Meldet es dagegen `citation-out-of-range`, **fehlt die Zieldatei** (oder sie
  ist kürzer geworden). Dann ist zuerst zu prüfen, ob der Abschnitt nur in eine
  **andere Datei gewandert** ist — die Baseline hat ihre Module schon einmal
  aufgeteilt. Ist er das, wandert der **Pfad** mit, wie beim Link, und die
  Spanne wird dort neu angekert. Erst wenn der Wortlaut **nirgends** im neuen
  Baum steht, greift [`MR-039`](../conventions.md#mr-039). Ohne diesen dritten
  Fall entfernte man eine Direktive, deren Quelle sehr wohl existiert.
- **Begründung:** Ohne diesen Nachtrag hätte der Bump eine gate-rote Nachwirkung
  ohne benannten Adressaten — und die erste Reaktion auf ein dauerrotes Gate ist
  erfahrungsgemäß, es abzuschalten, nicht es zu lesen. Der Schritt gehört in die
  Prozedur, die den Bump ohnehin beschreibt, statt in die Erinnerung dessen, der
  ihn fährt.
- **Auflösungs-Trigger:** entfällt, sobald keine `cite`-Direktive mehr auf den
  vendorten Baum zeigt — oder sobald das Produkt eine Anker-Form trägt, die
  ohne Zeilennummern auskommt.
