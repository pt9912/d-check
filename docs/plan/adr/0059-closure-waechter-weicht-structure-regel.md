# ADR-0059 — Der Zeichenketten-Wächter der Drei-Ausgänge-Regel weicht einer `structure`-Regel, mit gemessener Reichweiten-Differenz

**Status:** Accepted
**Datum:** 2026-08-26
**Autor:** pt9912
**Bezug:**
[`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(`forbid-pattern`, `sections`),
[`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(die benachbarte Platzhalter-Fähigkeit),
[ADR-0048](0048-closure-note-struktur-im-planning-modul.md) (Closure-Bindepunkt
und die Vertragsgrenze *Struktur, nicht Bedeutung*),
[`AGENTS.md` §3.6](../../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden)
(der Grund, warum diese Entscheidung eine ADR ist)

## Kontext

Die urteilsfreie Hälfte der Drei-Ausgänge-Regel des Baseline-Regelwerks
(`modul-05` §Offene Risiken werden bei Closure aufgelöst) wurde als
Bash-Skript `tools/harness/closure-outcomes.sh` durchgesetzt: vier
Zeichenketten-Formen über jede Datei in `docs/plan/planning/done/`, gebunden an
`make fullbuild`. Das Skript nannte im eigenen Kopf einen Auflösungs-Trigger —
*„sobald der Abschnitts-Skopus des Moduls den ganzen Slice umfasst, fällt
dieses Skript ersatzlos"*.

Zwei Annahmen dieses Triggers tragen nicht, beide nachgesehen statt erinnert:

1. **Den Skopus zu weiten genügt nicht.** Der Platzhalter-Erkenner des Moduls
   `planning` verlangt **whitespace-freie** Winkelklammern und deckt damit
   genau **eine** der vier Formen. Die Vorlagen-Form
   `<eingetreten: … | entfallen: … | weiter offen: …>` und die zwei
   Prosa-Formen bleiben unsichtbar, gleich wie weit der Skopus reicht.
2. **Die vermisste Fähigkeit existiert bereits.** Das Modul `structure` trägt
   `forbid-pattern` (Grund-Code `section-forbidden`) samt `sections` und
   `exempt-paths`. Es brauchte kein Produkt-Delta.

## Entscheidung

1. **Das Skript entfällt; eine `structure`-Regel im Closure-Profil trägt die
   Zusage.** Selektor `^# ` mit `sections: each` — geprüft wird **jeder**
   H1-Abschnitt, nicht nur der Titel-Abschnitt. Ein einzelner Selektor auf die
   Titelzeile hätte bei einer zweiten H1 die Spanne still gekappt; `sections:
   one` fängt das nicht, weil es nur **Selektor-Treffer** zählt.

2. **Der Ort ist das Closure-Profil, nicht `gates`.** Die Regel gilt dem
   **Übergang** nach `done/`, nicht dem Arbeitsbaum — dieselbe Einordnung, die
   [ADR-0048](0048-closure-note-struktur-im-planning-modul.md) Entscheidung 7
   für die Closure-Frage getroffen hat.

3. **Die Muster-Liste wandert in die Konfiguration.** Sie bleibt eine Liste und
   damit so gut wie ihre Pflege; neu ist, dass sie **deklariert** neben den
   übrigen Profil-Regeln steht statt in einem Skript, und dass sie über die
   geteilte Lexik des Produkts liest statt über einen `sed`-Ausdruck.

4. **Die Reichweiten-Differenz wird ausgewiesen, nicht behauptet.** Das Skript
   paarte Backticks **je Zeile**; das Produkt paart sie **absatzweise**, was der
   korrekteren Markdown-Lesart entspricht. Folge: ein Platzhalter, der zwischen
   zwei Backticks **eines Absatzes über Zeilengrenzen hinweg** liegt, ist für
   die Regel unsichtbar — das Skript meldete ihn. Gemessen, beide Richtungen:
   dieselbe eingefügte Risikozeile ergibt beim Skript Exit 1, bei der Regel
   Exit 0.

   **Genau diese Differenz macht die Entscheidung ADR-pflichtig**
   ([`AGENTS.md` §3.6](../../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden)):
   auf einer gemessenen Achse deckt der Nachfolger weniger als der Vorgänger.
   Die Entscheidung fällt dennoch für die Regel, weil die Alternative eine
   dauerhafte Doppelführung wäre und weil der Vorteil des Skripts auf dieser
   Achse aus einer **falscheren** Lexik stammte, nicht aus einem Entwurf.

## Verglichene Alternativen

| Alternative | Warum verworfen |
|---|---|
| Skopus des `planning`-Platzhalter-Erkenners weiten | Deckt eine von vier Formen; der Auflösungs-Trigger des Skripts war insoweit zu optimistisch |
| Skript **und** Regel nebeneinander lassen | Dauerhafte Doppelführung zweier Mechaniken für dieselbe Frage; die schwächere Lexik bliebe der Maßstab |
| Muster-Verbot als neue Produkt-Fähigkeit bauen | `forbid-pattern` existiert; ein Neubau wäre ein zweiter Auswertungspfad für dieselbe Frage |
| Backtick-Paarung im Produkt auf Zeilen-Lokalität ändern | Zöge eine **falschere** Markdown-Lesart in alle prosa-lesenden Module; die Klasse ist breiter als dieser Wächter und gehört ins Beobachtungs-Register |
| Selektor auf die Titelzeile (`^# Slice slice-`) | Eine zweite H1 kappt die Spanne still; die Nullmengen-Härte deckt nur den Fall *kein* Treffer |

## Konsequenzen

- `make fullbuild` trägt **fünf** Closure-Glieder statt sechs; die Zusage
  wandert in `make verify-closure-notes`.
- Die Zusage ist an **vier** Stellen enger als die Datei: die H1-Zeile selbst
  ist nicht Teil des Abschnitts-Texts, Fenced Blocks sind entfernt,
  Inline-Code-Spannen sind geleert, und eine absatzweite Spanne über
  Zeilengrenzen leert mit. Die ersten drei sind gewollt — sonst meldete ein
  Slice, der über die Platzhalter **schreibt**, seine eigene Dokumentation.
  Die vierte ist der Preis aus Entscheidung 4.
- Der Skript-Pfad geht ins codepaths-Tombstone-Register, wie bei den fünf
  abgelösten Skripten vor ihm.

## Fitness Function

- **Bewusstes Brechen je Form:** jede der vier Formen einzeln eingefügt ⇒
  `section-forbidden`, Exit 1; Rückbau ⇒ Exit 0. Gefahren am `make`-Target,
  nicht nur am Container-Aufruf.
- **Negativkontrolle:** dieselbe Form in Inline-Code ⇒ grün.
- **Zweite H1:** Platzhalter hinter `# Anhang` ⇒ Befund auf der Zeile jenes
  Abschnitts.
- **Nullmengen-Härte:** greift der Selektor in einer Datei nicht, meldet die
  Regel `section-missing` — sie kann nicht still leerlaufen.
- **Bestand:** 460 Dateien, 0 Befunde.

## Re-Evaluierungs-Trigger

- Ändert die Kanon-Vorlage ihre Platzhalter-Form, schweigt die Alternation —
  sie gehört beim nächsten Vorlagen-Bump mitgeprüft.
- Bekommt das Produkt einen Sensor für **mehrzeilige** Inline-Code-Spannen,
  fällt die Reichweiten-Differenz aus Entscheidung 4 weg; dann ist diese
  Abwägung neu zu bewerten.
- Erreicht die Klasse *Prosa verschwindet in einer absatzweiten Spanne* die
  Register-Schwelle, ist sie nicht mehr eine Grenze dieses Wächters, sondern
  eine Produkt-Frage.

## Geschichte

- 2026-08-26: Accepted (Closure `slice-143`).
