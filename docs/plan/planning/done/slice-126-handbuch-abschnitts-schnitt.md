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

- [x] Das Ventil-Gefälle steht im Handbuch **und** in der §6-Tabelle —
      **korrigiert nach Review**: nicht „kein Ventil", sondern **keine feine
      Achse**. Gemessen an einer Fixture: ohne `scope` Exit 1
      (`citation-out-of-range`), mit `citations.scope.roots`+`ignore` Exit 0,
      ohne `roots` fail-closed Exit 2; für `codepaths.check-lines` Marker
      stumm (Exit 1 nur auf der Kontrollzeile). Die Tabelle führt jetzt fünf
      Achsen.
- [x] Der §5-Abschnitt ist geschnitten: 183 Zeilen / sechs Module → fünf
      Abschnitte. **Eine der fünf Überschriften war zunächst selbst unehrlich**
      (zwei von drei `planning`-Fähigkeiten) — nach Review nennt sie auch das
      Wellen-Register. §Inhalt braucht nichts (nur `##`-Ebene); der Reviewer
      hat unabhängig belegt, dass kein Anker und kein Prosa-Verweis eine neue
      Kante überquert.
- [x] `citations` steht in **allen** Aufzählungen — CHANGELOG, Handbuch,
      beide READMEs, `operations.md`, Spezifikation, Lastenheft (0.65.2) und
      im Re-Evaluierungs-Trigger von [ADR-0058](../../adr/0058-konfigurations-flaechen-additiv-weiten.md).
      Sechs davon fand erst der Review.
- [x] `releasing.md`-Regel gezogen — **auf die Klasse**, nicht auf zwei
      Kapitel: jeder gegliederte Fließtext, den ein Release anfasst.
- [x] Beobachtungs-Register trägt die Klasse als **BEO-011**, Zähler **3** mit
      drei Belegen — nach der vom Review erzwungenen Neuzählung.
- [x] `make gates` Exit 0 (acht Gates, 450 Dateien, 0 Befunde, Coverage
      94,80 %); unabhängiger Review
      ([Report](../../../reviews/2026-08-23-slice-126-abschnitts-schnitt-review.md),
      Verdikt blockierend, 0 HIGH · 6 MEDIUM · 3 LOW · 1 INFO, alle zehn
      eingearbeitet); §11-Zeile 1.57 mit unveränderter Software-Version
      `v0.63.0`.

## 5. Abnahme-Punkte / Risiken

- **Ein reiner Abschnitts-Schnitt sieht im Diff aus wie eine Umschreibung.**
  Wird beim Auftrennen zugleich formuliert, ist nicht mehr prüfbar, ob Inhalt
  verloren ging. — **Ausgang:** *nicht eingetreten, und der Beleg war der
  Diff selbst.* Der Schnitt-Commit löschte genau **drei** Zeilen im Handbuch,
  alle drei erweitert wieder eingefügt; der Reviewer hat das unabhängig
  nachgezählt. Die Trennung von Bewegung und Formulierung hat getragen — die
  eine Botschaft, die daraus „vier Überschriften, sonst nichts" machte, war
  dagegen zu knapp (es waren +28/−3 Zeilen).
- **Quer-Verweise auf den alten Abschnitt sind teils Prosa** („siehe oben",
  „derselbe Abschnitt") und damit gate-unsichtbar; `anchors` fängt nur die
  Anker. — **Ausgang:** *nicht eingetreten.* Kein Dokument verlinkt überhaupt
  auf einen Handbuch-Anker, und kein Prosa-Verweis überquert eine der vier
  neuen Kanten — unabhängig geprüft, nicht angenommen.
- **Die Register-Klasse könnte zu breit geraten.** „Aus dem Anlass statt aus
  dem Bestand" beschreibt fast jeden Übergeneralisierungs-Fehler; ein Eintrag,
  der alles trifft, steuert nichts. Er braucht einen Ableiter, der eine
  Handlung nennt. — **Ausgang:** **eingetreten.** Der erste Zählstand war 5;
  zwei davon fielen unter andere Klassen (`BEO-004`-nah bzw. `BEO-002`), und
  die dritte der drei benannten Ausprägungen war unter den fünf gar nicht
  vertreten. Der Zähler steht jetzt auf 3. Das benannte Risiko war also real,
  und die Vorsorge — drei *handlungsförmige* Ableiter statt einer
  Definition — hat es nicht verhindert, sondern nur nachweisbar gemacht.

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

Geliefert sind beide Befunde: die Ventil-Lage bei `citations` steht als
Fünf-Achsen-Tabelle im Handbuch, und der 183-Zeilen-Abschnitt ist in fünf
Abschnitte geschnitten, deren Überschriften nennen, was unter ihnen steht.
Dazu die Regel-Weitung in der Release-Prozedur, `citations` in sieben
Aufzählungen und `BEO-011` im Register.

**Die Lehre steht nicht in dem, was der Slice geliefert hat, sondern in dem,
was er dabei falsch gemacht hat.** Der unabhängige Review fand drei Befunde,
die je eine der drei Ausprägungen von `BEO-011` sind — **im selben Commit, der
`BEO-011` anlegt**. (a) „`citations` trägt gar kein Ventil, keinen
Konfigurations-Block": am Code falsch, `citations.scope` wirkt, gemessen. (b)
„Das ist Absicht — eine ausdrücklich gesetzte Prüfung soll nicht zeilenweise
zurückgenommen werden können": eine Begründung, die in keinem Vertrag steht —
genau das, was der zweite Ableiter desselben Eintrags verbietet. (c) Die
Regel-Weitung nannte wieder Orte (§4 und §5) statt der Klasse.

**Eine Klasse aufzuschreiben schützt nicht davor, sie zu begehen.** Das ist
kein Nebenbefund, sondern die eigentliche Aussage dieses Slice: der Eintrag
war frisch formuliert, präsent und selbst geschrieben — und hat dreimal nicht
gegriffen. Was gegriffen hat, war ein zweiter Leser mit einem Prüfauftrag.
Deshalb bleibt die mechanische Form im Register ausdrücklich **offen** statt
aus dem Zähler abgeleitet; ein Wortlisten-Lint hätte (a) vielleicht gefangen,
(b) und (c) sicher nicht.

**Der Zähler des Eintrags war selbst aus dem Anlass gebildet.** Fünf
behauptete Instanzen, zwei davon aus anderen Klassen ([slice-122](slice-122-versions-musterliste.md):
eine von der geteilten Nachrunde widerlegte Zusage; die Spiegel-Liste in
[slice-125](slice-125-release-v0630.md): `BEO-002`), und die Ausprägung (c) war
unter den fünf nicht vertreten. Wer eine Klasse benennt, zählt ihre Instanzen
danach — nicht die Vorfälle, die ihn auf sie gebracht haben.

**Ein Schnitt kann den Defekt erzeugen, den er behebt.** Eine der fünf neuen
Überschriften nannte zwei von drei `planning`-Fähigkeiten; das Wellen-Register
belegt rund ein Drittel des Abschnitts. Der Prüfsatz, den dieser Slice in die
Release-Prozedur geschrieben hat — *nennt die Überschrift alles, was unter ihr
steht?* —, war auf die alte Überschrift angewandt worden und auf die neuen
nicht.

**Und eine bereits notierte Lehre hat trotzdem nicht getragen:** der
[`MR-025`](../../../../harness/conventions.md#mr-025)-`grep` nach dem **alten
Wortlaut**, repo-weit, steht seit
[slice-125](slice-125-release-v0630.md) als Ableiter in dessen Closure-Notiz.
Hier wurde er wieder nicht gefahren — sechs Fundstellen blieben stehen, zwei
davon in den Spec-Straten. Eine Regel für Menschen kann weiter verfehlt
werden; das ist genau die Aussage, die `BEO-002` seit drei Instanzen trägt.

**Nachträglich korrigiert, nicht amendiert:** der bereits getaggte
`0.63.0`-CHANGELOG-Block war retro-editiert worden und steht wieder auf dem
Tag-Stand; die Änderungen dieses Slice liegen unter `[Unreleased]`. Und der
Spec-Bump (Lastenheft 0.65.2 samt Historie-Zeilen in beiden Straten) fehlte am
Review-Nachzug und wurde als eigener Commit nachgeholt.
