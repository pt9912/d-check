# Slice slice-163: Baseline-Zitate außerhalb des Konventionsspeichers

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`MR-039`](../../../../harness/conventions.md#mr-039) (der
Geltungsbereich, der diese Dokumente ausdrücklich mitnennt);
[`MR-051`](../../../../harness/conventions.md#mr-051) (das Neu-Ankern beim
Bump); [slice-152](../done/slice-152-citations-scharfschalten.md) (der Zug, der
den Speicher abgearbeitet hat).

**Berührte Spec-Stellen:** — (Doku-Bestand; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

[`MR-039`](../../../../harness/conventions.md#mr-039) §Geltungsbereich nennt
**alle lebenden** Dokumente dieses Repos — *„Konventionsspeicher, `AGENTS.md`,
`harness/README.md`, die Skills und die lebenden Planungs-Dokumente"*.
[slice-152](../done/slice-152-citations-scharfschalten.md) hat den
**Konventionsspeicher** abgearbeitet und den Rest ausdrücklich ausgeklammert
(§3 dort: *„Keine Ausweitung auf Zitate außerhalb des Konventionsspeichers in
diesem Zug"*).

**Der Rest ist benannt, nicht geschätzt:** der Review von slice-152 hat
mindestens zwei **wörtliche, auszeichenbare** Baseline-Zitate außerhalb des
Speichers gefunden — je eines in [`AGENTS.md`](../../../../AGENTS.md) und im
Reviewer-Skill. Wie viele es insgesamt sind, ist die erste Messung dieses Slice
und hier **nicht** vorweggenommen.

## 2. Vorgehen

1. **Den Bestand zählen**, mit derselben Methode wie
   [slice-152](../done/slice-152-citations-scharfschalten.md): Zitate ab 16
   Zeichen, Produkt-Lexik (Fence-Automat, absatzweise Backtick-Spannen), je
   Zitat die drei Klassen — wörtlich / abweichend / gegen den gepinnten Stand
   nicht prüfbar.
2. **Bei Abweichung die Richtung entscheiden**, nach der Linie aus
   [slice-152](../done/slice-152-citations-scharfschalten.md): was die Quelle
   **hat** und das Zitat weglässt, ist ein Transkriptions-Fehler und wird
   korrigiert; was das Zitat **hinzufügt**, ist ein Autoren-Akt und wird
   deklariert statt korrigiert.
3. **Die Skills sind der heikle Teil.** Sie werden von Agenten gelesen, nicht
   von Menschen gerendert; eine Direktive am Zeilenende ist dort unauffällig,
   aber die Zitate tragen die Beweislast der Anker. Vor der Auszeichnung prüfen,
   ob das Zitat überhaupt den gepinnten Stand meint.
4. Bewusstes Brechen je gedeckter Form, **Ursache gelesen**; Rückbau grün.
5. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Auszeichnung in `done/`, `docs/reviews/` oder `conventions/done/`.**
  Eingefrorene Dokumente zitieren den Stand ihrer Zeit.
- **Keine Änderung an der Paarungsregel** des Moduls, auch nicht für die
  Delta-Tabelle aus [`MR-039`](../../../../harness/conventions.md#mr-039), die
  [slice-152](../done/slice-152-citations-scharfschalten.md) als nicht
  adressierbar ausgewiesen hat.
- **Keine neue Ventil-Achse.**

## 4. Definition of Done

- [x] Der Bestand ist **gezählt**, nach den drei Klassen getrennt — und die
      mittlere ist erst in der Nacharbeit berichtet worden: **49** Zitate ab 16
      Zeichen über neun lebende Dokumente, **10** wörtlich, **eines** abweichend
      (jetzt korrigiert), der Rest gegen den gepinnten Stand nicht prüfbar.
      Die erste Fassung maß **sechs** Dateien und nannte es *„die lebenden
      Dokumente"*.
- [x] Je Abweichung eine Entscheidung nach der Richtungs-Linie: der eine Fall
      lässt weg, was die Quelle **hat** ⇒ Transkriptions-Fehler, korrigiert.
- [x] Was ausgezeichnet wird, ist ausgezeichnet; was nicht, ist **benannt** —
      und der Grund war zunächst **falsch** (siehe §5).
- [x] Ein konstruierter Verstoß je gedeckter Form, **Ursache gelesen** — beide
      Platzierungsformen, die zweite erst in der Nacharbeit.
- [x] Doku-Currency geprüft: keine `CHANGELOG`-Pflicht (weder Produkt- noch
      nutzersichtbar), aber zwei Spiegel-Sätze waren durch das erste
      Nicht-Baseline-Ziel unvollständig geworden und sind nachgezogen.
- [x] `make gates` grün (Exit explizit), `make fullbuild` grün; unabhängiger
      Review.

## 5. Abnahme-Punkte / Risiken

- **Jede neue Direktive erhöht die Bump-Last.**
  [`MR-051`](../../../../harness/conventions.md#mr-051) legt das Neu-Ankern in
  die Prozedur, aber die Zahl der neu zu ankernden Spannen wächst mit jedem Zug.
  — **Ausgang: eingetreten, und anders als gedacht.** Die Last ist **nicht nur**
  Bump-Last: eine der sechs Direktiven zeigt auf
  [`AGENTS.md`](../../../../AGENTS.md), nicht auf die Baseline. Sie wandert bei
  **jeder** Änderung ihres Ziels, und `AGENTS.md` ist das meistgeänderte
  Dokument des Repos.
  [`MR-051`](../../../../harness/conventions.md#mr-051) deckte das nicht; jetzt
  tut er es, samt der zwei Spiegel, die *„die planmäßige Rot-Quelle ist der
  Bump"* sagten.
- **Ein Zitat im Skill trägt Beweislast.** Wird es beim Korrigieren gekürzt oder
  verschoben, ändert sich womöglich, was der Anker dem Reviewer sagt. —
  **Ausgang: entfallen.** Der Diff ist in beiden Skills rein **additiv** — kein
  Zitat wurde angefasst. Die eine Korrektur eines Wortlauts traf das
  Beobachtungs-Register, nicht einen Skill, und stellte dort die Fettung der
  Quelle wieder her.
- **Die Richtungs-Linie ist neu und ungeprüft.** Sie stammt aus einem einzigen
  Zug über einen einzigen Bestand; ob sie außerhalb des Speichers trägt, ist
  offen. — **Ausgang: eingetreten, und sie hat getragen.** Genau ein Fall stand
  zur Entscheidung, und er fiel eindeutig: das Zitat lässt eine Fettung weg, die
  die Quelle hat ⇒ Transkriptions-Fehler, korrigiert. **Ein Fall ist kein
  Beleg für eine Linie** — sie ist damit weniger ungeprüft als vorher, nicht
  bestätigt.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Messung zeigt, dass ein
Skill-Zitat den gepinnten Stand gar nicht meint.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Doku (GF), Konventionsspeicher (GF), Skills (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27):
  [`BEO-020`](../observations.md) — die gemessene Menge muss die sein, über die
  geredet wird; [`BEO-011`](../observations.md) — die Regel gehört aus dem
  Bestand, nicht aus dem Anlass; [`BEO-012`](../observations.md) — ein Zitat
  trägt nur, was in seinem Geltungsbereich steht.

Slice-ID: slice-163. Betroffene IDs: — (Doku-Bestand; keine Anforderung).
Module: `citations`. Gates: `make doc-check`, `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Fortsetzung eines abgeschlossenen Zugs auf
den benannten Rest.

## 9. Closure-Notiz (nach `done/`)

**Sechs Zitate ausgezeichnet, eines korrigiert, eine Regel geschärft — und die
Begründung, warum zwei angeblich nicht gehen, war dreifach falsch.**

**Die Messung, und ihre erste Fassung war zu eng.** Gemessen habe ich zuerst
**sechs** Dateien und es *„die lebenden Dokumente außerhalb des Speichers"*
genannt. [`MR-039`](../../../../harness/conventions.md#mr-039)
§Geltungsbereich nennt die lebenden **Planungs**-Dokumente mit. Mit ihnen sind
es **neun** Dateien und **49** Zitate ab 16 Zeichen — nicht 44. Der Schluss
kippt dadurch nicht, die Aussage schon; als Klasse ist das
[`BEO-020`](../observations.md).

| Größe | Wert |
|---|---|
| Zitate ab 16 Zeichen, neun lebende Dokumente | **49** |
| wörtliche Teilstrings des gepinnten Stands | **10** |
| davon Marker/Term statt Zitat (*„Nichts in Arbeit"*) | **3** |
| **echte Zitate** | **7** |
| davon ausgezeichnet | **6** |
| nicht adressierbar | **1** |

**Die dritte Klasse hatte ich gar nicht berichtet.** §2 verlangt *wörtlich /
abweichend / nicht prüfbar*; mein Trichter nannte nur die erste und las sich
trotzdem vollständig. Sie hat genau ein Mitglied: das `BEO-015`-Kanon-Zitat
lässt die **Fettung des zweiten Halbsatzes** weg, die die Quelle hat. Nach der
Richtungs-Linie ein Transkriptions-Fehler — korrigiert, und damit ist die Klasse
jetzt wirklich leer.

**„Nur vier sind adressierbar" war falsch, und zwar dreifach.** Ich hatte
argumentiert, eine Tabellenzeile trage mehrere Zitate und eine Direktive davor
bräche die Tabelle. Eine Tabelle **ist** ein Absatz — es gibt keinen
*„folgenden"*. `AGENTS.md` §4 trägt in der fraglichen Zeile **genau ein** Zitat.
Und die Form, die eine Tabellenzeile nicht bricht, hat dieses Repo längst
entschieden:
[ADR-0040](../../adr/0040-kommentar-suffix-in-tabellenzeilen.md), Kommentar-Zelle
hinter der schließenden Pipe. Beide sind jetzt ausgezeichnet, beide mit
Bruchprobe belegt.

**Das Argument stimmt — für einen anderen Fall.** Das
`BEO-015`-Kanon-Zitat ist das **vierte** seiner Zeile, und dort greift die
Erst-Zitat-Paarung tatsächlich nicht. Ich hatte ein richtiges Argument den
falschen Fundstellen zugeordnet.

**Ein Fall war lehrreicher als erwartet.** Der Reviewer-Skill zitiert
*„Halluzinierte Gates sind die häufigste Form von Harness-Lüge"* und schreibt es
[`AGENTS.md`](../../../../AGENTS.md) §4 zu. Der Satz steht wörtlich in
**beiden** — Baseline und `AGENTS.md`. Die Direktive folgt der
**Zuschreibung**: wer den Skill liest, soll die Stelle prüfen können, die er
zitiert zu haben behauptet. Das ist zugleich die erste `cite`-Direktive auf ein
**Nicht-Baseline-Ziel**, und dafür gab es keine Regel — siehe §5.

**Der eigentliche Ertrag ist
[`MR-052`](../../../../harness/conventions.md#mr-052).** Der offene Befund aus
[slice-152](../done/slice-152-citations-scharfschalten.md) war, dass zwei
Stellen dasselbe `v5.11.0`-Zitat verschieden wiedergeben — und das galt als
*nicht entscheidbar*, weil der alte Pin nicht mehr vendored ist. Der Bump
entfernt den alten Baum aus dem **Arbeitsbaum**, nicht aus dem **Repo**; ein
`git show <lösch-commit>^:<alter-pfad>` beantwortet die Frage in einer Zeile.
Gemessen: [`MR-039`](../../../../harness/conventions.md#mr-039)s Tabelle ist
wörtlich, [`MR-033`](../../../../harness/conventions.md#mr-033) hat die Fettung
verschoben. Beide bleiben stehen; festgehalten ist nur, **welche**.

**Auch dabei musste der Review korrigieren.** Meine erste Fassung schrieb die
Annahme *„der alte Stand ist nicht mehr zugänglich"*
[`MR-039`](../../../../harness/conventions.md#mr-039) zu. Der Eintrag sagt das
nirgends — er begründet das Einfrieren mit der historisch korrekten Aussage,
einem Grund, der von der Zugänglichkeit unabhängig ist. Träger der Annahme war
die **Praxis**, und zwar meine eigene in
[slice-152](../done/slice-152-citations-scharfschalten.md).
[`BEO-012`](../observations.md), und diesmal gegen einen Eintrag des eigenen
Speichers.

**Eine Methoden-Grenze, benannt:** die 49 sind der **rohe** Lauf. §2 nennt die
gestrippte Produkt-Lexik; für den **Zitattext** ist roh die richtige Wahl
([ADR-0060](../../adr/0060-citations-marker-scan-geteilte-prosa-antwort.md):
gestrippt wird die Marker-Suche, nicht der Vergleich), aber die Methodenangabe
sagte das Gegenteil.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 561 Dateien, 0 Befunde),
`make fullbuild` (Exit 0, 48 Anforderungen / 0 Waisen, bench Median 819 ms).
**Proben, Ursache je gelesen:** je Platzierungsform ein verfälschtes Zitat ⇒
`citation-mismatch` mit richtigem Ziel und Grund-Code, Rückbau grün; die zwei
Tabellenzeilen eigens gebrochen, nachdem ihre Nicht-Adressierbarkeit widerlegt
war. Ein unabhängiger Review ist gelaufen; sein Urteil war *„schließbar nach
Nacharbeit"*, seine neun Befunde sind eingearbeitet, und seine zwei HIGH sind
eigens nachgemessen statt übernommen.
