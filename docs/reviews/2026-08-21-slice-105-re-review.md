# Re-Review-Report: slice-105 — Heilung der Review-Auflagen (Chronologie-Monotonie) — 2026-08-21

**Review-Art:** bestätigende Re-Review (Code, gegen die Vertragsflächen) —
frischer Kontext; geprüft wird ausschließlich, ob die Heilung jeden Befund des
Erst-Reviews schließt, ob sie neue Defekte einführt und ob die gewählte
Lesart konsistent über alle Vertragsflächen gezogen ist. Nicht geprüft:
DoD-Abhakung, Gate-Lauf-Bestätigung, Plan-Konformität (Verifikation).

**Gegenstand:** Heilungs-Commit `8d56c3c` (auf Feat-Commit `9b33d8b`,
Erst-Report `docs/reviews/2026-08-21-slice-105-chronologie-review.md`,
Findings F-1…F-5); Arbeitsbaum-Stand `HEAD` = `8d56c3c` — der Arbeitsbaum
trägt zusätzlich **ungecommittete Release-Prep-Entwürfe** (CHANGELOG,
READMEs, Benutzerhandbuch, `version.md`), die nicht Gegenstand dieser
Re-Review sind, aber auf Lesart-Konsistenz mitgelesen wurden (unten).

**Skill:** `.harness/skills/reviewer.md` @ 1.4.0 ·
**Modell:** claude-fable-5 · **Datum:** 2026-08-21

**Eingangs-Kontext:** dieselben Verträge wie im Erst-Review —
`DC-FA-STRUCT-001` (Lastenheft 0.61.1), §`DC-FA-STRUCT-001.a`
(§6-Bedingungs-Tabelle, Schritt-6-Fließtext, §2-Schema, §4-Grund-Codes),
ADR-0057 (Proposed), ADR-0054 (E1/E2/E4), Hard Rules `AGENTS.md` §3.

---

## Heilungs-Prüfung je Befund

### F-1 (MEDIUM — Typ-Mischung, zwei Vertragslesarten) · **geschlossen**

Gepinnt wurde die **Paar-Lesart mit Anker-Reset**: die Misch-Zelle meldet
sich selbst, die gesunde Folge-Zeile dahinter meldet nicht.

1. **Alle Vertragsstellen sagen jetzt dieselbe Lesart — geprüft, ohne
   Befund.** Sechs Stellen, wortverschieden, bedeutungsgleich:
   - `spec/spezifikation.md:2032` (§6-Bedingungs-Tabelle): „den Typ ihres
     typisierbaren **Vorgängers** fortführt (Paar-Lesart; nach jedem Befund
     setzt der Vergleich neu auf)" — die alte „Spalten-Typ"-Zelle ist ersetzt;
     die Klammer bindet das Wort „Vorgänger" an den Vergleichs-Anker, sodass
     die Tabellenzeile allein nicht mehr die zweite Lesart hergibt.
   - `spec/spezifikation.md:2070–2074` (Schritt-6-Fließtext): Anker-Reset
     nach **jedem** Befund, „die gesunde Folge-Zeile hinter einer Misch-Zelle
     meldet **nicht**" — explizit.
   - `spec/spezifikation.md:2653` (§4-Zeile): „weicht vom Typ ihrer
     typisierbaren **Vorgänger-Zelle** ab … genau einer je gemeldeter Zelle".
   - `spec/lastenheft.md:2346–2350` (Chronologie-Absatz): unverändert, aber
     jetzt **wahr** — „an dieser Zeile … hinter der gemeldeten Zelle setzt
     der Vergleich beim nächsten typisierbaren Nachbar-Paar wieder auf" war
     unter dem alten Code für den Misch-Fall falsch (der Anker blieb stehen)
     und beschreibt das neue Verhalten exakt.
   - `spec/lastenheft.md:2390` (Akzeptanzkriterium): „an **genau dieser**
     Zeile — der Vergleichs-Anker wird nach jedem Befund zurückgesetzt: eine
     gesunde Folge-Zeile hinter einer Misch-Zelle meldet **nicht**".
   - `docs/plan/adr/0057-…md:120–125` (E7): „nach **jedem** Befund — auch
     nach einer Typ-Mischung … sie meldet sich **selbst**, und die gesunde
     Folge-Zeile dahinter meldet **nicht** (Paar-Lesart, im Review als einzig
     verbleibende Lesart gepinnt)".
2. **Code implementiert die Lesart — geprüft, ohne Befund.**
   `internal/hexagon/core/rules/structure_tableorder.go:177–184`: der
   Mismatch-Zweig gibt `return &f, nil, 0` zurück — identisch zu den beiden
   anderen Fehlerfällen (Zeilen 166, 172). Die Meldung nutzt `prevLine` nur
   im Zweig mit gültigem Anker; nach dem Reset läuft die nächste Zeile über
   den `prev == nil`-Frühausstieg (174–176), eine „Zeile 0" kann in keiner
   Meldung erscheinen.
3. **Kaskaden-Test pinnt das Richtige — geprüft, ohne Befund.**
   `structure_tableorder_test.go:121–131` (`TestTableOrderTypMischungKaskade`):
   Datum–Version–Datum–Datum, `desc` ⇒ `len(f) == 1`,
   `Reason == section-cell-untyped`, `Line == 8` (die Misch-Zelle). Der Test
   **diskriminiert**: das alte Verhalten lieferte für genau diese Eingabe
   zwei Befunde (Misch-Zelle + gesunde Folge-Zeile) und machte den Test rot.
4. **Randfälle der Reset-Lesart — geprüft, ohne Befund** (kein neues stilles
   Loch):
   - *Misch-Zelle als letzte Zeile:* ein Befund, danach Tabellen-/
     Schleifenende — kein Sonderpfad.
   - *Zwei Misch-Zellen hintereinander / alternierend
     (Datum–Version–Datum–Version):* nach jedem Reset wird die nächste
     typisierbare Zelle neuer Anker; jede weitere Mischung relativ zu ihrem
     Anker meldet wieder — alternierende Spalten melden jede zweite Zeile,
     nichts verstummt.
   - *Version–Datum–Version (bzw. maskierter Ordnungs-Bruch über die
     Misch-Zelle hinweg):* die beiden Nicht-Nachbarn werden nicht verglichen —
     das ist **dieselbe**, bereits vor der Heilung vertragliche und getestete
     Nachbar-Paar-Grenze der anderen beiden Fehlerfälle
     (`TestTableOrderSetztNachUntypisierbarNeuAuf`; AK 2390: „benachbarte
     typisierbare Paare werden weiterhin verglichen"). Die Heilung erweitert
     diese benannte Grenze auf den Misch-Fall, statt eine neue zu reißen.

### F-2 (LOW — Overflow-Segment degradiert still) · **geschlossen**

- **Code:** `structure_tableorder.go:59–63` — `versionSegments`-Ergebnis
  `nil` ⇒ `typeTableKey` gibt `ok=false` zurück; der Aufrufer meldet
  `section-cell-untyped` über den bestehenden „kein Token"-Pfad hinaus
  (Zeilen 168–173). Geprüft, ohne Befund.
- **Kommentar korrigiert:** Zeilen 67–72 — der falsche
  Unerreichbarkeits-Satz ist ersetzt durch „ERREICHBAR (`\d+` ist
  unbegrenzt)" samt Wirkungs-Begründung. Geprüft, ohne Befund.
- **Test trifft den Zweig:** `TestTableOrderVersionsSegmentUeberlauf`
  (Test-Zeilen 135–145) — 20-stelliges Segment (> int64), erwartet genau ein
  `section-cell-untyped` an Zeile 8; das alte Verhalten (stiller
  Kleinst-Vergleich ⇒ **kein** Befund in dieser absteigenden Spalte) machte
  den Test rot. Geprüft, ohne Befund.
- **Keine nil-Segmentfolge im Vergleich mehr:** `versionSegments` hat genau
  einen Aufrufer (`typeTableKey`, Zeile 59, per grep über `internal/`);
  Versions-Keys entstehen nur noch mit non-nil `segs`, Datums-Keys
  (`segs == nil`, `isDate == true`) erreichen den Segment-Zweig von
  `compareTableKeys` nie (Vergleich nur bei Typ-Gleichheit, Datums-Zweig
  Zeile 90–91). Geprüft, ohne Befund.
- **Komposition mit der „erster Treffer"-Zusage:** eine Zelle, deren
  **erster** Muster-Treffer überläuft, ist untypisierbar, auch wenn später
  in der Zelle ein gültiger Token stünde — das ist die wörtliche
  Erster-Treffer-Semantik aus §6, kein Widerspruch. Vertraglich benannt in
  `spec/spezifikation.md:2068–2070` und §4 (2653: „Segment außerhalb des
  Zahlbereichs"). Geprüft, ohne Befund.

### F-3 (INFO — Kopplungs-Reichweite undokumentiert) · **geschlossen**

`lexikon_kopplung_test.go:113–119`: die Reichweiten-Notiz steht am
Kopplungs-Test und ist **faktisch korrekt** — per grep haben
`tableHeaderOrSeparator` und `tableCells` heute je genau zwei Konsumenten
(`planning_waves.go:208/211`, `structure_tableorder.go:136/161`),
`tableRowLine` drei (`targets.go:157` dazu); die Notiz benennt die
ADR-0054-Schwelle und die Handlungs-Anweisung beim dritten Konsumenten
(„wer einen hinzufügt, schreibt ihn hier daneben"). Geprüft, ohne Befund.

### F-4 (INFO — Ausdrucks-Grenze unbenannt) · **geschlossen**

Benannt an beiden verlangten Orten: Lastenheft-Out-of-Scope
(`spec/lastenheft.md:2414–2419`) und ADR-0057 §Konsequenzen (Zeilen
167–173), beide mit derselben Begründung (Identität trägt keine Spalte;
Öffnung änderte die Deduplikations-Semantik aller sieben Bedingungen; Ausweg
= Change Request). **Faktisch korrekt:** `model/config.go:451–457` —
`Identity() = Files + " :: " + (Section|SectionPattern)`, ohne Spalte;
`configyaml.go:225–228` bricht bei doppelter Identität mit Fehler ab
(Exit 2, laut). Geprüft, ohne Befund.

### F-5 (INFO — zweite Zell-Lesart unabgegrenzt) · **geschlossen**

- `markdown.go:214–220`: der `tableCells`-Kommentar beansprucht nicht mehr
  „die geteilte Zell-Antwort" des Produkts, sondern „DIESER Vertragsfläche",
  benennt den trace-Leser als andere Fläche und die konkrete Konsequenz
  (Pipe im Backtick-Span verschiebt die Spaltenadresse).
- ADR-0057 §Konsequenzen (Zeilen 174–181): die zweite Lesart ist als
  existent benannt, als „andere Frage im Sinn von ADR-0054 Entscheidung 2"
  eingeordnet, **und** der Text behält die Ehrlichkeit, die der Befund
  verlangte: träte die Divergenz real ein, „ist das der Latenz-Fall der
  geteilten-Lexik-Klasse und nach deren Regeln zu behandeln" — die
  Einordnung als „andere Frage" ist damit ausdrücklich revidierbar deklariert
  statt als endgültige Ein-Antwort-Behauptung. Keine falsche
  „eine Antwort"-Behauptung mehr. Geprüft, ohne Befund.

---

## Neue Defekte durch die Heilung / Konsistenz-Restprüfung

1. **Verhaltens-Delta der Heilung — geprüft, ohne Befund:** genau zwei
   Code-Änderungen (Misch-Zweig-Reset, Overflow ⇒ untypisierbar), beide nur
   in Rändern wirksam, die im heutigen Bestand nicht vorkommen (Kaskade
   hinter Mischung, ≥19-stelliges Segment); beide jetzt vertraglich und je
   mit diskriminierendem Test. Der Retro-Beleg (27 `section-unordered`)
   bleibt davon unberührt — beide Deltas betreffen ausschließlich
   `section-cell-untyped`-Pfade.
2. **Alte Lesart restlos entfernt — geprüft, ohne Befund:** grep
   „Spalten-Typ" über das Repo (ohne `docs/reviews/`) trifft nur noch die
   Spezifikations-Historienzeile (die den behobenen Zustand **beschreibt**)
   und den Kaskaden-Test-Kommentar (der die gesunde Folge-Zeile gerade als
   „trägt den Spalten-Typ" charakterisiert, um das Nicht-Melden zu betonen)
   — beides korrekte Verwendungen, keine Vertragsaussage der alten Lesart.
3. **Weitere Semantik-Flächen — geprüft, ohne Befund:** §2-Schema-Zeile
   (`spec/spezifikation.md:2533`) und `--doctor`-Klartexte
   (`app/diagnose.go:146–147`) nennen die Mischung ohne Lesart-Aussage
   (kein Widerspruch möglich); `--print-config`-Vorlage trägt nur die beiden
   Schlüssel-Kommentare; `.d-check.closure.yml` unberührt.
4. **Ungecommittete Release-Prep-Entwürfe (außerhalb des Prüf-Commits,
   nur Lesart-Konsistenz) — geprüft, ohne Befund:** der
   CHANGELOG-0.61.0-Entwurf sagt „der Vergleichs-Anker setzt nach jedem
   Befund neu auf", der Handbuch-Entwurf (§5-Absatz „Fünftens",
   `docs/user/benutzerhandbuch.md:1724–1727`, und die §11-Zeile 1.52)
   sagt „an **genau dieser** Zeile, dahinter setzt der Vergleich neu auf"
   bzw. „Anker-Reset dahinter" — beide tragen die gepinnte Lesart.
5. **Versions-/Historien-Hygiene — geprüft, ohne Befund:** Lastenheft-Kopf
   `**Version:** 0.61.1`; §7-Zeile 0.61.1 **oben**, über 0.61.0 (absteigend,
   `spec/lastenheft.md:2791–2792`); Spezifikations-§7-Nachzug-Zeile **oben**
   (`spec/spezifikation.md:2681`, über der 0.61.0-Zeile desselben Datums);
   beide Historien-Einträge benennen ehrlich Review-Anlass und Inhalt.
6. **Import-/Hermetik-Ränder der Heilung — geprüft, ohne Befund:** die
   Heilung fügt keine Imports, keine Eingaben, keine Gates hinzu; der Diff
   ist auf die fünf benannten Dateien plus zwei Spec-Straten begrenzt.
7. **Randnotiz (kein Finding, Verifikations-Zeiger):** die
   Commit-Botschaft nennt „390/0" nach zuvor „389/0" bei **zwei** neu
   hinzugekommenen Testfunktionen; Zähl-Konvention ist Verifikations-, nicht
   Review-Gegenstand. Ebenso trägt der Slice-Plan noch den Arbeitsnamen
   `section-key-untyped` (Plan-Historie, prüft die Verifikation).

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 0 | — |
| LOW | 0 | — |
| INFO | 0 | — |

## Verdikt

**APPROVE (bestätigend).** Alle fünf Befunde des Erst-Reviews sind
geschlossen: F-1 ist auf genau eine Lesart gepinnt — an sechs
Vertragsstellen bedeutungsgleich formuliert, im Code als Anker-Reset
implementiert und durch einen diskriminierenden Kaskaden-Test festgenagelt;
die Reset-Lesart reißt kein neues stilles Loch (alle geprüften Randfälle
laufen in die bereits benannte Nachbar-Paar-Grenze). F-2 ist fail-closed
geheilt, getestet und der Zweig hat nur einen Aufrufer. F-3/F-4/F-5 sind an
den verlangten Orten benannt und faktisch korrekt. Die Heilung führt keine
neuen Defekte ein; die ungecommitteten Release-Prep-Entwürfe tragen die
gepinnte Lesart bereits mit.
