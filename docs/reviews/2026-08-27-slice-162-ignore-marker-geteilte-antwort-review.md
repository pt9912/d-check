# Review slice-162 — `d-check:ignore` beantwortet dieselbe Frage anders

**Gegenstand:** [slice-162](../plan/planning/done/slice-162-ignore-marker-geteilte-antwort.md), Stand des Feat-Commits vor der Nacharbeit.
**Datum:** 2026-08-27. **Reviewer:** unabhängiger Subagent.

---

## Urteil

**Schließbar nach Nacharbeit.** Die Code-Änderung ist klein, korrekt und grün.
Nachzuarbeiten sind **die Zahlen**, **die tragende Begründung von Entscheidung
2**, **zwei unbenannte Verhaltens-Grenzen** und **die Straten-Verortung**.

Der Reviewer hat die Produkt-Lexik in `awk` nachgebaut und die Methode
**validiert**: sie reproduziert die Messung von
[ADR-0060](../plan/adr/0060-citations-marker-scan-geteilte-prosa-antwort.md) am
damaligen Stand exakt (173 / 63).

## Befunde

| # | Rang | Befund |
|---|---|---|
| F-1 | **HIGH** | *„`versions` liest keine Prosa"* trägt nicht — es liest **alle** Zeilen, eine **Obermenge**. Auf **derselben** Prosa-Zeile meldet `codepaths` und schweigt `versions`: die Zweifach-Antwort ist **verschoben, nicht behoben**. Der eigene Testname behauptet das Gegenteil dessen, was sein Fixture zeigt |
| F-2 | **HIGH** | *„160 davon ausschließlich in Inline-Code"* ist die falsche Grundgesamtheit — 160/161 ist die **Bare-Form**-Zahl. Richtig sind **183** (HEAD) bzw. 184 (Eltern); die Dateizahl 544 stammt aus der Ära von ADR-0060, richtig sind 552/553. Sieben Stellen betroffen |
| F-3 | **HIGH** | Zwei Grenzen des Strippens wurden mitgeerbt und nirgends benannt: **Verschluckung** (ein gesetzter Marker wirkt nicht, wenn eine Spanne desselben Absatzes ihn umschließt — Falsch-Rot) und **Erzeugung** (ein unpaariger Backtick kippt die Parität, die Erwähnung wirkt doch — stilles Grün). Die Zusage gilt nur bei gerader Parität |
| F-4 | **HIGH** | [ADR-0019](../plan/adr/0019-versions-pin-fence-ausnahme.md) und ADR-0054 §2 skopieren den **Pin**-Scan, nicht die Marker-Erkennung — über den Geltungsbereich hinaus zitiert (`BEO-012`) |
| F-5 | MEDIUM | Der Wächter fängt **2 von 6** konstruierten Umgehungen; die Erlaubnis-Liste ist **datei-granular**; die Compiler-Behauptung gilt an genau einer Stelle, nicht allgemein |
| F-6 | MEDIUM | Die normative Prosa landete in `DC-FA-REF-001(.a)` statt in den geänderten Anforderungen; die Algorithmus-Schritte blieben unangetastet; `DC-FA-CODE-001` bekam **kein** neues Akzeptanzkriterium; die Historie-Zeile ist misattribuiert |
| F-7 | MEDIUM | Der `GRENZE`-Kommentar an der Konstanten beschreibt den **behobenen** Defekt (§3.7) |
| F-8 | MEDIUM | *„alle fünf Befunde sind echt"* widerspricht der Definition des Slice-Plans **und** dem eigenen §Entscheidung 4. Ehrliche Bilanz: **null** gefundene Doku-Defekte, fünf implizite Ausnahmen wurden explizit |
| F-9 | LOW | Die ADR-Index-Zeile bricht mitten im Satz ab |
| F-10 | LOW | `diagrams` liest eine **dritte** Eingabe-Klasse (die Fence-Öffnungszeile), die die ADR-Tabelle nicht kennt (§3.8) |
| F-11 | LOW | Import-Block in `codepaths_test.go` nicht sortiert |

**Was geprüft wurde und trägt:** die Zahl **5** stimmt exakt; alle drei
gesetzten Marker sind **berechtigt**, und nichts anderes wurde mitgeschaltet;
der Wächter sitzt im richtigen Paket, verletzt weder §3.2 noch `arch-check`,
seine Erlaubnis-Liste fängt genau die zwei richtigen, und die **Gegenrichtung
läuft nachweislich**. §3 ist **nicht** verletzt.

## Erledigung

- **F-1**, **F-2**, **F-3**, **F-4**, **F-8**, **F-10** durch
  [ADR-0062](../plan/adr/0062-ventil-marker-versions-ist-eine-benannte-grenze.md),
  die die **Herleitung** von ADR-0061 ablöst — deren Kern ist gepusht und
  `Accepted`, also nach §3.5 unantastbar. Das Produktverhalten bleibt gleich.
  Die zwei HIGH-Proben (Divergenz auf derselben Prosa-Zeile; Paritäts-Loch) sind
  **eigens nachgefahren**, nicht übernommen.
- **F-2** zusätzlich an allen lebenden Stellen: Lastenheft-Historie,
  Spezifikations-Historie, `CHANGELOG`, ADR-Index.
- **F-5** teilweise geschlossen: das **Literal** neben der Konstanten und
  `strings.Index` neben `strings.Contains` werden jetzt gefangen — Gegenprobe
  gefahren. Die vier übrigen Formen und die Datei-Granularität stehen als
  Grenze in ADR-0062 Entscheidung 6; die Compiler-Behauptung ist zurückgenommen.
- **F-6** durch Verankerung in `DC-FA-CODE-001.a` Schritt 1, `DC-FA-ID-001.a`
  Bedingung 4 und der Ventil-Zusage in `DC-FA-CODE-001`, samt eigenem
  Boundary-Kriterium dort; Historie-Zeile korrigiert.
- **F-7**, **F-9**, **F-11** direkt behoben.
