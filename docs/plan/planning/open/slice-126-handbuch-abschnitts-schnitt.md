# Slice slice-126: Der §5-Abschnitt, dessen Überschrift eine Teilmenge nennt — und das ungesagte Ventil-Gefälle bei `citations`

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-82-config-flaechen](../welle-82-config-flaechen.md) (zugeordnet
bei der Eröffnung; die Welle war für ihren Closure noch offen, als der Befund
entstand).

**Bezug:**
[`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in),
[`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
(dieselben Grund-Codes, verschiedene Ventil-Lage);
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus) (die
Korrektur ändert kein Verhalten); der Benutzerhandbuch-Standard und die
§4-Checkliste der Release-Prozedur.

**Berührte Spec-Stellen:** — (Nutzer-Doku, Release-Prozedur und
Beobachtungs-Register; keine Anforderung und kein Spec-Stratum ändert sich).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

Zwei Befunde derselben Welle, beide **nach** der Closure von
[slice-125](../done/slice-125-release-v0630.md) entstanden, beide gemessen statt
vermutet:

**Erstens: dasselbe Grund-Code-Paar ist im einen Modul stummschaltbar, im
anderen nicht — und die Doku sagt es nicht.** `citation-out-of-range` und
`citation-inverted-range` entstehen in **zwei** Modulen. Aus
`codepaths.check-lines` heraus liegen sie **innerhalb** der Zeilen-Schleife, die
den Marker `d-check:ignore` auswertet, sind also stummschaltbar; aus dem Modul
`citations` heraus nicht — es ist **parameterlos** (`CitationsConfig struct{}`),
ohne `exempt-paths`, ohne Marker, ohne Konfiguration überhaupt. Gemessen an
einer Drei-Fälle-Fixture: dieselbe Zeilen-Konstruktion, derselbe Grund-Code, der
Marker hilft einmal und einmal nicht. Wer den Befund sieht, kann aus der Doku
nicht ableiten, welches Modul spricht und ob das Ventil greift.

**Zweitens: ein §5-Abschnitt nennt in seiner Überschrift eine Teilmenge seines
Inhalts.** *„Zitate und Zeilen-Referenzen gegen ihre Quelle prüfen
(`codepaths.check-lines` / Modul `citations`)"* trägt **183 Zeilen** und
**sechs** Module: `codepaths.check-lines`, `citations`, `vcs`, `commits`,
`planning` (zwei Fähigkeiten) und `tracked`. Vier davon stehen nicht in der
Überschrift. Wer in §5 die `vcs`- oder `tracked`-Konfiguration sucht, findet sie
unter einer Überschrift über Zitate.

**Der Befund ist aus dem Bestand gezogen, nicht aus dem Anlass** — das ist die
Lehre dieser Welle, hier angewandt. Zensus über **alle** dreizehn
§5-Abschnitte: drei sind über 90 Zeilen lang (`trace` 203, *Weitere Module* 155,
dieser 183), aber nur **einer** nennt in der Überschrift weniger, als er trägt.
Die Länge ist nicht der Defekt, die **Unauffindbarkeit** ist es.

## 2. Vorgehen

1. **Das Ventil-Gefälle aussprechen**, wo der Nutzer den Befund trifft: im
   Handbuch bei beiden Fähigkeiten und in der §6-Modul-Tabelle, deren Zeilen für
   `codepaths` und `citations` heute dieselben Grund-Codes listen, ohne den
   Unterschied zu nennen. Vorher an einer Fixture messen, nachher gegenprüfen.
2. **Den Abschnitt schneiden:** je Modul eine eigene `###`-Überschrift, die
   nennt, was darunter steht. Kein Umschreiben des Inhalts — reine Auftrennung,
   damit die Änderung lesbar bleibt und der Diff die Bewegung zeigt.
   Anschließend die §Inhalt-Liste und alle Quer-Verweise auf den alten
   Abschnitt prüfen (`anchors` fängt gebrochene Anker, **nicht** die Prosa, die
   auf „der Abschnitt oben" zeigt).
3. **Die zwei unvollständigen Aufzählungen aus slice-125 nachziehen:** die
   CHANGELOG-Zeile und die „benannte Liste"-Notiz zählen die ventillosen
   Zeilen-Melder auf und lassen `citations` aus.
4. **Die Regel auf ihre Klasse weiten.** `docs/user/releasing.md` §4 trägt die
   Anti-Anlagerungs-Regel bereits — aber nur für **§4-Aufgaben**, geschrieben
   nach der Lehre, dass §4.12 auf ~330 Zeilen / 8 Themen anwuchs. Die
   Anlagerung ist danach in **§5** passiert. Die Regel ist aus dem Anlass
   gezogen statt aus der Klasse und wird auf beide Kapitel gezogen.
5. **Beobachtungs-Register:** die Klasse *„eine Aussage/Regel aus dem Anlass
   statt aus dem Bestand"* ist in dieser Welle sechsmal aufgetreten (slice-122,
   -123, -124, zweimal in slice-125, hier). Sie hat keinen Eintrag. Eintragen
   mit Zähler und der Frage, ob eine mechanische Form existiert.
6. **Kein Tag.** Doku-Korrektur ohne Software-Änderung; §11-Zeile mit
   **unveränderter** Software-Version, wie bei den Handbuch-Ständen 1.54 und
   1.55. Reist mit dem nächsten Release.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Funktions-Änderung.** Insbesondere bekommt `citations` **kein**
  Ventil — ob ein parameterloses Modul eines braucht, ist eine eigene Frage mit
  eigener Messung und wäre ein Change Request am Lastenheft, kein Doku-Slice.
- **Kein Kürzen der beiden anderen langen Abschnitte** (`trace`, *Weitere
  Module*). Sie sind lang, aber ihre Überschriften sind ehrlich; Länge allein
  ist hier nicht der Defekt.
- **Keine mechanische Prüfung** der Abschnitts-Aufteilung. Ob `structure` das
  überhaupt messen kann (es zählt Tasks, nicht Themen), ist im
  Register-Eintrag als offene Frage benannt, nicht hier entschieden.

## 4. Definition of Done

- [ ] Das Ventil-Gefälle steht im Handbuch **und** in der §6-Tabelle; vorher/
      nachher an einer Fixture gemessen, beide Läufe mit Exit-Code belegt.
- [ ] Der §5-Abschnitt ist so geschnitten, dass **jede** Überschrift nennt, was
      unter ihr steht; §Inhalt und Quer-Verweise nachgezogen.
- [ ] CHANGELOG-Enumeration und „benannte Liste"-Notiz nennen `citations`.
- [ ] `releasing.md` §4-Regel auf §4 **und** §5 gezogen.
- [ ] Beobachtungs-Register trägt die Klasse mit Zähler und Beleg-Slices.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review; §11-Zeile mit
      unveränderter Software-Version.

## 5. Abnahme-Punkte / Risiken

- **Ein reiner Abschnitts-Schnitt sieht im Diff aus wie eine Umschreibung.**
  Wird beim Auftrennen zugleich formuliert, ist nicht mehr prüfbar, ob Inhalt
  verloren ging. — **Ausgang:** *(bei Closure)*
- **Quer-Verweise auf den alten Abschnitt sind teils Prosa** („siehe oben",
  „derselbe Abschnitt") und damit gate-unsichtbar; `anchors` fängt nur die
  Anker. — **Ausgang:** *(bei Closure)*
- **Die Register-Klasse könnte zu breit geraten.** „Aus dem Anlass statt aus
  dem Bestand" beschreibt fast jeden Übergeneralisierungs-Fehler; ein Eintrag,
  der alles trifft, steuert nichts. Er braucht einen Ableiter, der eine
  Handlung nennt. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): sofort — [slice-125](../done/slice-125-release-v0630.md)
ist geschlossen, `v0.63.0` ausgeliefert, `in-progress/` frei.

**Rückführungen:** `in-progress` → `next`, falls der Review zeigt, dass das
Ventil-Gefälle keine Doku-Frage ist, sondern ein Change Request am Lastenheft —
dann geht der Anforderungs-Text vor.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Nutzer-Doku (GF), Release-Prozedur (GF),
  Beobachtungs-Register (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23): **BEO-002**
  (Semantik-Änderung nur im Körper nachgezogen, Ränder bleiben stehen) ist
  unmittelbar einschlägig — die zwei unvollständigen Aufzählungen aus slice-125
  sind genau solche Ränder. **BEO-009** (Botschaft behauptet mehr, als die
  Arbeit trägt) gilt für jede Zahl, die dieser Slice nennt; die
  Vorgänger-Botschaft hat hier bereits eine Menge falsch benannt. **BEO-006**
  beim pfad-selektiven Commit.

Slice-ID: slice-126. Betroffene IDs:
[`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in),
[`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus). Module:
Nutzer-Doku, Release-Prozedur, Beobachtungs-Register. Gates: `make doc-check`,
`make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Doku-Korrektur nach etablierter Form; der
Anforderungs-Text bleibt unberührt.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
