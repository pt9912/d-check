# Welle 81 — Zustandsfelder tragen Zustand, keine Chronik (Baseline v5.9.0) — Closure-Notiz

**Welle:** welle-81-zustandsfelder
**Abschluss:** 2026-08-22
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- **Pin-Hebung über zwei Stufen**
  ([slice-117](slice-117-baseline-v590-bump.md)): das Bundle am Tag `v5.9.0`
  vendored (51 Dateien beide Bäume plus Manifest), der Vorgänger-Baum entfernt,
  `--verify` offline grün und `--check-latest` meldet den Pin als neuesten
  Release. Nachfolge-Adaption angelegt, Vorgängerin aufgelöst, alle
  pin-gebundenen Verweise nach dem Drei-Klassen-Zensus gehoben.
- **Die Regel verkörpert, bevor sie angewandt wird**
  ([slice-118](slice-118-zustandsfeld-regel.md)): das Briefing sagt in seiner
  Rolle, dass ein Zustandsfeld ein Zustands-Artefakt mit **eigener** Form ist —
  es erbt vom Kommentar die zwei Tests, nicht dessen fünf Klassen; der
  Reviewer-Skill trägt den HIGH-Anker samt der Ausnahme für Daten, die ein
  **benannter Trigger** pflegt.
- **Die Kopf-Zustandsfelder**
  ([slice-119](slice-119-kopf-zustandsfelder.md)): drei von vier Zeilen
  ersatzlos entfernt (beide lebenden Register und das Technik-Stratum); die
  Sicht behält ihre und sagt an Ort und Stelle, warum sie dort ein
  Frische-Marker ist. Nebenbei tragen beide Spec-Straten die Rollen-Zeile der
  Vorlage nach, die im Bestand fehlte.
- **Die Register selbst**
  ([slice-120](slice-120-register-und-drift-log.md)): acht `Stand`-Zellen auf
  Zustand, Gegenmittel und mechanische Form gezogen (577–3 011 → 326–636
  Zeichen), das Drift-Log von **69 auf 10** Zeilen zurückgeschnitten, die
  Meilenstein-Status-Form ergänzt und der fünfte Treffer geheilt — die
  Ketten-Nacherzählung in der `Stand`-Zeile des Konventionsspeichers.

## Was hat funktioniert?

- **Vendoren → verkörpern → anwenden.** Dieselbe Reihenfolge wie in der
  Vor-Welle: die Regel stand im Briefing, bevor der erste Bestand angefasst
  wurde, und jede Anwendung konnte sich darauf berufen.
- **Der neue HIGH-Anker hat sofort gegriffen — gegen den eigenen Text.** In
  drei der vier Slices fand der Review ein Zustandsfeld, das eine Chronik
  erzählte oder etwas behauptete, das der Bestand nicht trägt. Ein Anker, der
  beim ersten Lauf den Autor trifft, ist ein guter Anker.
- **Messen vor dem Planen.** Das Bundle-Delta sah nach 33 geänderten Dateien
  aus; gemessen waren es 22 reine Stempel und zehn Dateien mit Regel-Inhalt.
  Ohne die Messung wäre die Welle nach Dateizahl geschnitten worden.
- **Die Gegenprobe vor dem Löschen.** Bevor 59 Drift-Log-Zeilen fielen, wurde
  mechanisch nachgezählt, dass jede genannte Welle im Closure-Log und jeder
  genannte Slice in `done/` steht. Keine entfernte Zeile war die einzige Spur.

## Was ging anders als geplant?

- **Ein `grep` kennt keine Zeitform.** Der Pfad-Zensus der Pin-Hebung hob eine
  **Vergangenheits**-Aussage mit; der Tombstone behauptete danach, der *neue*
  Baum sei entfernt worden — und machte damit seine eigene Begründung
  unlesbar. Die Beobachtung dazu hat eine zweite Richtung bekommen: nicht nur
  die vergessene, auch die **über**-gehobene Stelle.
- **Die Ausnahme-Begründung war selbst ein Chronik-Feld.** An der Sicht stand
  „ein benannter Trigger pflegt sie" — der Review widerlegte es am Bestand: 94
  Commits am Vertrag zwischen zwei Marker-Hebungen, der Marker stand still.
  Wer eine Ausnahme schreibt, muss ihren Träger benennen können; kann er es
  nicht, ist es keine Ausnahme, sondern eine Behauptung.
- **Eine Überschrift ist eine Aussage.** „Kommentare **und Zustandsfelder**
  tragen eine der fünf Klassen" behauptete genau das, was der Kanon nicht sagt,
  und widersprach ihrem eigenen Absatz drei Zeilen tiefer.
- **Eine Aufzählung belegt keinen Klassen-Abschluss.** Die Welle zählte vier
  Treffer; der Kanon eröffnet allgemein („ein Feld, das einen Zustand trägt").
  Der fünfte fiel durchs Raster und kam erst im dritten Review dazu.
- **Beim Kürzen ist die Frage nicht „ist das alt?", sondern „ist das
  Zustand?"** Zwei Register-Zellen verloren eine Deklaration, die Zustand war
  — der Review holte sie zurück. Und die neue Meilenstein-Prosa widersprach
  ihrer eigenen leeren Tabelle.

## Steering-Loop-Einträge

- **Briefing** geschärft: Hard Rule §3.7 deckt Zustandsfelder mit eigener Form,
  mit der Trigger-Ausnahme und ohne Bestandsgrenze. Liegt in `AGENTS.md` §3.7.
- **Reviewer-Skill** (1.7.0): der HIGH-Eintrag *Zustandsfeld trägt Chronik* mit
  beiden Ausprägungen und der Nicht-Melde-Ausnahme. Liegt in
  `.harness/skills/reviewer.md §Repo-spezifische Anker`.
- **Bump-Prozedur** geschärft: der Drei-Klassen-Zensus prüft jede gehobene
  Stelle zusätzlich auf ihre **Zeitform**. Liegt im Hebungs-Zensus der
  aktuellen Pin-Adaption.

## Beobachtungs-Register (Zeiger)

Gelesen zur Closure ([`observations.md`](../observations.md)): **BEO-008 steht
seit dieser Welle bei 3** — Schwelle erreicht. Die dort benannte mechanische
Form ist heute **nicht konfigurierbar**: das Versions-Modul trägt genau ein
Muster gegen eine Quelle; ein zweiter Abgleich braucht einen Change Request,
keinen Slice. Bis dahin trägt die Prozedur. Alle übrigen Einträge unverändert;
keine neue Beobachtung — die drei Selbst-Treffer dieser Welle sind die bereits
verkörperte Klasse BEO-002 in neuen Formen (Titel, Verweis, Kürzung).

## Folge-Slices

- **Keiner.** `open/` ist leer, §Nächste Wellen trägt `— keine —`.
- **Benannte Folgepunkte ohne Slice** (aus dem letzten Review): 90 von 118
  `done/`-Slices tragen ein historisches `Status`-Kopffeld, **elf** davon
  widersprechen ihrem Verzeichnis — das Briefing erklärt das historische Feld
  für Alt-Slices, die Widersprüche bleiben offen; ein doppelter Anker in der
  Versions-Datei; ein ADR-Statusfeld mit Chronik, das mit der
  Immutabilitäts-Regel kollidiert und darum eine Vorrangregel bräuchte.
- **Change-Request-Kandidat:** ein zweiter Versions-Abgleich (Baseline-Tag in
  URLs und Prosa gegen den Pin) — die 3×-Form von BEO-008, heute nicht
  konfigurierbar.

## Verifikation

- `make gates` nach jedem Slice grün; `make fullbuild` zur Closure grün.
- Vier unabhängige Reviews, **kein HIGH**; zehn MEDIUM, jedes mit realem
  Gehalt, alle eingearbeitet.
- Das Vendoring stimmt gegen das Manifest **und** gegen den Kurs-Tag — kein
  Handanlegen an den Bundle-Bäumen.
- Der Drift-Log-Schnitt ist gedeckt: jede genannte Welle im Closure-Log, jeder
  genannte Slice in `done/`, die vier Umplanungs-Verdachtsfälle einzeln
  gegengeprüft.
- Die Chronologie-Regeln beißen auf der zurückgeschnittenen Tabelle
  unverändert (konstruierte Gegenprobe).
