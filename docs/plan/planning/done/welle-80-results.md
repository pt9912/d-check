# Welle 80 — Struktur-IDs nach Baseline, die Umkehr von MR-027 — Closure-Notiz

**Welle:** welle-80-struktur-ids
**Abschluss:** 2026-08-22
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- **Konvention zuerst** ([slice-113](welle-80/slice-113-struktur-id-konvention.md)):
  die deklarierte Abweichung ist per `git mv` aufgelöst (Index-Zeile umgezogen,
  beide Voll-Slug-Anker mitgenommen, damit eingefrorene Verweise weiter
  auflösen), die Vergabe-Aussage steht in
  [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage)
  — fortlaufend je Datei, keine Lücken-Nachbelegung, kein Bereichssegment, der
  Link trägt den Abschnitt und der Text die Kennung, nicht in
  Commit-Botschaften. Der ADR-Index und `AGENTS.md` §5 sagen dasselbe in ihrer
  Rolle. Das Prüf-Profil kennt seither zwei Kennungs-Muster mit Linkpflicht.
- **Vergabe im Technik-Stratum**
  ([slice-114](welle-80/slice-114-spec-vergabe-spezifikation.md)): 66 Kennungen —
  fünf Schema-Überschriften, sieben Defaults, 51 Grund-Codes, drei externe
  Verträge, die drei Tabellen je mit einer neuen ersten Spalte. Zwei
  Konsumenten wurden im selben Commit scharf: eine Abschnitts-Invariante über
  die Schema-Überschriften und der bestehende Grund-Code-Lockstep-Test, der
  jetzt zusätzlich die Kennungs-Spalte auf Eindeutigkeit prüft.
- **Vergabe in der Sicht**
  ([slice-115](welle-80/slice-115-arc-vergabe-architektur.md)): elf Kennungen — sieben
  Komponenten in einer neuen Komponenten-Tabelle (und im Kasten-Label des
  Flowcharts), vier externe Berührungspunkte als Fortsetzung; die
  Schichten-Tabelle referenziert dieselben sieben, ohne neue zu vergeben. Der
  Konsument ist das Modul `diagrams`, auf die Spec-Straten gescopt und in der
  Netzlos-Modulliste verankert.
- **Anwendung statt Deklaration**
  ([slice-116](welle-80/slice-116-adr-neuzugangs-regel.md)): die beiden `Proposed`-ADRs
  tragen die neue Adressierungs-Form (eine bekam ihr `Schärft:`-Feld erstmals
  überhaupt), der ADR-Index zeigt beide Formen, und der Reviewer-Skill hat
  einen Anker dafür — mit dem `Accepted`-Bestand ausdrücklich ausgenommen.

## Was hat funktioniert?

- **Die Reihenfolge Konvention → Vergabe → Anwendung.** Die Regel stand,
  bevor die erste Kennung existierte; jede Vergabe konnte sich darauf berufen,
  und der letzte Slice hatte etwas zum Anwenden. Die Alternative — erst
  vergeben, dann begründen — hätte die Deklaration zur Nacherzählung gemacht.
- **Der Gate-Konsument im selben Commit wie die Vergabe.** Vorher gemessen rot,
  danach grün: so ist die Zusage belegt und nicht behauptet. Der `pre-commit`-
  Hook macht diese Reihenfolge zur einzig möglichen — ein roter Zwischenstand
  ist nicht committebar.
- **Vier unabhängige Reviews, vier Mal mit Substanz.** Kein HIGH, aber acht
  MEDIUM über die Welle, und **jedes einzelne** hat etwas Reales getroffen:
  einen still entkommenden Fall, einen falsch gewordenen Lauf-Beleg, ein
  über-feuerndes Gate, ein gebrochen ausgezeichnetes Beispiel, ein
  unvollständiges Feld.

## Was ging anders als geplant?

- **Die Form, die man einführt, gilt zuerst für einen selbst.** Dreimal
  blockierte der eigene Wächter den eigenen Text, bevor ein Commit durchging:
  Beispiel-Kennungen in Slice-Plänen, dann in einer Closure-Notiz, dann das
  Beispiel im ADR-Index — dieses sogar in gebrochener Auszeichnung, die kein
  Gate sah. Wer eine Kennungs-Pflicht einführt, macht auch das *Reden über*
  Kennungen pflichtig; Beispiele gehören in Platzhalter-Form oder hinter das
  Zeilen-Ventil.
- **Ein Wächter muss die Lexik seines Moduls sprechen.** Die
  Abschnitts-Invariante verlangte Zeilenanfang und genau ein Leerzeichen; das
  Modul trimmt beliebigen Weißraum und nimmt auch Tab. Eine eingerückte
  Sektion entkam still — der Wächter behauptete eine Deckung, die er nicht
  hatte. Dieselbe Klasse wie die geteilte Lexik im Code, hier erstmals
  **konfigurationsseitig**.
- **Die Vorlage sticht den eigenen Plan.** Der Plan wollte auch die
  Fehlermodell-Tabelle bekennen; die Baseline vergibt an Komponenten und
  Berührungspunkte, nicht an Fehlerquellen.
- **Eine grüne Gegenprobe war der Fund, nicht der Fehler.** Entfernt man nur
  die Vergabe-Zeile, bleibt die Kennung definiert, weil die zweite Tabelle sie
  nennt: das Diagramm-Modul liest Token, nicht Struktur — Nennung und Vergabe
  sind für es dasselbe. Als Grenze benannt, statt als Zusage überdehnt.
- **Ein neues Modul hat Spiegel in der Ausführungsschicht.** Die
  `--disable`-Kette der fokussierten Gates spiegelt die Modulliste und sagt das
  im eigenen Kommentar — nachgezogen hat sie erst der Review. Das ist der neue
  Register-Eintrag dieser Welle.

## Steering-Loop-Einträge

- **Reviewer-Skill** geschärft (1.6.0): der MEDIUM-Anker
  *Adressierungs-Form eines Neuzugangs* — ein neues `Schärft:`/`Bezug:`-Feld
  nennt die Kennung, wo das Zielelement eine trägt; der `Accepted`-Bestand vor
  dieser Welle ist ausdrücklich ausgenommen. Liegt in
  `.harness/skills/reviewer.md §Repo-spezifische Anker`.
- **Prüf-Profil** um drei Sensoren gewachsen: zwei Kennungs-Muster mit
  Linkpflicht, eine Abschnitts-Invariante über die Schema-Überschriften und das
  Diagramm-Modul. Liegt in `.d-check.yml §ids`, `§structure`, `§diagrams`.
- **Autoritäts-Doku** trägt die Vergabe-Regel an drei Stellen in ihrer je
  eigenen Rolle: Konventionsspeicher (Aussage), ADR-Index (Feld-Form),
  `AGENTS.md` §5 (Slice-Kopf-Feld).

## Beobachtungs-Register (Zeiger)

Gelesen zur Closure ([`observations.md`](../observations.md)): **kein Eintrag
erreicht die Schwelle 3.** BEO-002 und BEO-003 sind je erneut eingetreten und
bleiben verkörpert — zitiert, nicht weitergezählt. Neu bei 1: **BEO-010** —
die Modulliste des eigenen Prüf-Profils hat Spiegel außerhalb der Config, und
keiner davon ist gate-gedeckt.

## Folge-Slices

- **Keiner.** `open/` ist leer, §Nächste Wellen trägt `— keine —`. Die
  Grenzen dieser Welle sind benannt und tragen ihre Trigger: die
  ausgeschriebene Präfix-Negation wäre mit einem eigenen Schlüssel klarer
  (Change-Request-Kandidat, kein Defekt), Tabellen-Kennungen bleiben für die
  Abschnitts-Invariante unerreichbar (nur die Grund-Code-Spalte ist über den
  Test gebunden), Diagramme außerhalb der Spec-Straten sind ungewacht, und die
  feste Dreistelligkeit lässt eine vertippte Kennung still durch.
- **Aus der Vor-Welle unverändert offen:** ob zwei flache Wellendokumente
  derselben Kennung selbst meldepflichtig sein sollten — eine Produktfrage für
  einen Change Request, kein Slice.

## Verifikation

- `make gates` nach jedem Slice grün; `make fullbuild` zur Closure grün
  (acht Gates, Image-Integrationstest, Benchmark, Requirements-Vollständigkeit,
  Closure-Note-Struktur).
- Vier CI-Läufe auf frischem Runner grün — jeder gepushte Stand der Welle.
- Jede Zusage der drei neuen Sensoren ist **vorher rot und nachher grün**
  gemessen, dazu je eine konstruierte Gegenprobe; alle Rückbauten per
  Dateikopie mit md5-Vergleich.
- Die Referenz-Richtung ist gemessen, nicht gesetzt: das Technik-Stratum nennt
  keine Kennung der Sicht.
