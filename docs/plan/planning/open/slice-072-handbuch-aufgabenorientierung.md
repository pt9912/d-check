# Slice slice-072: Handbuch — Aufgabenorientierung der §4-Kapitel

**Status:** open (Backlog; noch keiner Welle zugeordnet).

**Welle:** keine; wartet auf Aufnahme in eine Welle (Roadmap §Nächste Wellen).

**Bezug:** Redaktionelle Korrektur des
[Benutzerhandbuchs](../../../user/benutzerhandbuch.md) gegen den
[Benutzerhandbuch-Standard](../../../user/benutzerhandbuch-standard.md) §2
(„Aufgaben statt Funktionen beschreiben") und §5 („Schritt-für-Schritt-
Anleitungen"). **Kein Change Request** (kein Verhaltens- oder Vertragsdelta —
das Produkt bleibt unberührt), **kein ADR** (keine neue Architekturentscheidung),
**kein Release** (nur abgeleitete Nutzer-Dokumentation). Betrifft ausschließlich
`docs/user/benutzerhandbuch.md`.

**Autor:** pt9912. **Datum:** 2026-07-17.

---

## 1. Ziel

§4 des Benutzerhandbuchs ist stellenweise **funktions-** statt aufgabenorientiert:
Blöcke tragen Config-Schlüssel als Titel, begründen eine Fähigkeit aus d-checks
**Produkt-Innensicht** (welches Modul was tut oder nicht tut) und zählen Felder
auf, wo der Leser eine Handlungsfolge braucht. Der Standard verlangt das
Gegenteil: *„Schreibe nicht, was die Software kann. Schreibe, wie der Nutzer damit
seine Arbeit erledigt."*

Der Slice zieht die betroffenen §4-Blöcke auf die Ziel-zuerst-Form nach —
Ausgangslage → Ziel → Vorgehen → Ergebnis — und räumt die dabei sichtbare
Duplikation zwischen §4 (Aufgabe) und §5 (Referenz) auf.

**Nicht** Gegenstand: §5 und §6. Sie sind laut Standard §3 (Punkt 5
„Konfiguration") legitime **Referenz**-Kapitel; eine Feld-/Optionsreferenz ist
dort korrekt. Sie werden nur dort angefasst, wo §4-Inhalt zu ihnen wandert.

## 2. Entscheidungen / Regel

- **Schablone ist der bestehende Bestand, nicht eine neue Erfindung.** §4.10
  (`--repair`) und der §4.12-Kreuzverweis-Block („Prüfen, ob Ihre RTM-Tabelle noch
  zum Design passt", slice-071) erfüllen den Standard bereits: Ausgangslage → Ziel
  → nummerierte Schritte → Ergebnis-Deutung → Fallstricke als Hinweis. Sie sind
  die Vorlage; es wird kein neuer Stil eingeführt.
- **Titel benennen ein Nutzerziel, das Flag/der Schlüssel steht in Klammern**
  (gelebtes Muster von §4.9/§4.10).
- **Referenz gehört nach §5, Aufgabe nach §4.** Wo §4 heute Grammatik-, Feld- oder
  Fail-closed-Enumerationen trägt, die §5 bereits führt, wird §4 gekürzt und
  verweist; **kein Text wird verdoppelt**.
- **Kompatibilitäts-/Versionsprosa gehört nicht in eine Aufgabe** — Fehlerbilder
  nach §7, Versionsaussagen nach §11.
- **Die verifizierten Beispiele bleiben verifiziert.** Jede Umstellung hält die
  beiden Handbuch-Harnesse grün: Config-Beispiele gegen den Validator
  (slice-061) und Ausgabeblöcke E2E gegen das echte Binary (slice-062). Ein
  Ausgabeblock verliert weder seinen Anker noch seine Form-Token.

## 3. Definition of Done

- [ ] **B-3 (HIGH):** §4.12 „Unterstützte Definitionssyntax" (~50 Zeilen
  Grammatik-/Feldreferenz **vor** dem „Vorgehen") gekürzt; die Tabellen-Binding-
  Regeln leben nur noch in §5, §4 verweist. Duplikat zu §5 aufgelöst.
- [ ] **B-1/B-2 (MEDIUM–HIGH):** §4.7 `order`/`direction` und `token` auf
  Ziel-Titel + Lesersituation umgestellt (heute: „Trägt eine Klasse zusätzlich
  `order` …, prüft d-check auch …" bzw. „`matrix` sieht standardmäßig nur
  Markdown-Links").
- [ ] **B-6 (MEDIUM, strukturell):** §4.12 (heute ~300 Zeilen = 20 % des
  Handbuchs, mindestens fünf Aufgaben in einer) aufgetrennt — RTM erzeugen/gaten
  bleibt §4.12, Coverage/Modalität/Kreuzverweis werden eigene §4.x mit eigener
  Ziel/Vorgehen/Ergebnis-Klammer. Doppelte WAISE-Definition entfernt.
- [ ] **B-4 (MEDIUM):** §4.12 Modalität als Aufgabe („Nur MUSS-Anforderungen
  gaten, SOLLTE/KANN advisory lassen"); Klassifikations-Mechanik und
  Fail-closed-Liste nach §5.
- [ ] **B-7 (MEDIUM):** §4.4 öffnet mit der Entscheidungsfrage (frisches vs.
  bestehendes Repo) statt mit elf Zeilen Fähigkeits-Inventar.
- [ ] **B-5 (MEDIUM):** §4.12-Versionsgrenze („Bis v0.42.0 …, ab v0.43.1 …") als
  Fehlerbild nach §7, Versionsaussage nach §11.
- [ ] **B-8 (LOW):** Flag-/Schlüssel-Titel auf Ziel-Titel (§4.5, §4.11, §4.4
  `--id-prefix`, §4.12 `trace`/`trace.coverage`) — nur im Zug der ohnehin
  berührten Abschnitte.
- [ ] **Kein Regress:** beide Handbuch-Harnesse grün (Config-Beispiele,
  E2E-Ausgabeblöcke); Handbuch-Version + §11-Zeile nachgezogen; `make gates` grün.
- [ ] **Gegenprobe:** ein Abschnitt-für-Abschnitt-Audit gegen den Standard meldet
  für §4 keine offenen HIGH/MEDIUM-Befunde mehr — und benennt ausdrücklich, was
  bewusst Referenz bleibt.

## 4. Risiken / offene Punkte

- **Kein Gate erzwingt Aufgabenorientierung.** Sie ist eine Erkenntnis-, keine
  Laufzeit-Eigenschaft; die Disziplin liegt beim Autor und beim Review. Dieser
  Slice räumt den Bestand auf, verhindert aber den Rückfall nicht (siehe unten).
- **Die eigentliche Ursache ist strukturell, nicht redaktionell:** §4.12 wuchs
  über die Slices 066/067/068/070/071, weil **jeder** Slice seine Fähigkeit an die
  bestehende Aufgabe anhängte, statt eine Aufgabe zu schreiben. Ohne eine Regel am
  Release-Prep-Punkt (analog zur Prosa-Currency-Liste in
  [`releasing.md`](../../../user/releasing.md) §Release-Prep) erodiert §4 erneut.
  **Offener Designpunkt:** ob dieser Slice die Regel gleich mit in die
  Release-Prep-Checkliste schreibt („neuer Handbuch-Abschnitt = eigene Aufgabe mit
  Ziel/Vorgehen/Ergebnis, keine Anhängung") — das wäre die billigste dauerhafte
  Sicherung und wäre konsistent damit, wie die Currency-Blindflecken dort schon
  behandelt werden.
- **Auftrennung von §4.12 (B-6) berührt die Inhalts-Navigation** (§Inhalt,
  Querverweise „siehe §4.12" aus §5/§6/operations.md). Die Verweise sind
  gate-geschützt (`anchors`), ein vergessener Verweis fällt also auf — aber die
  Umnummerierung ist Handarbeit.
- **Grenze zwischen Aufgabe und Referenz ist Ermessen.** Der Slice zieht sie
  bewusst so: Was der Leser braucht, um **einmal** ans Ziel zu kommen, gehört nach
  §4; was er **nachschlägt**, nach §5. Grenzfälle werden zugunsten von §5
  entschieden (kein Text zweimal).

## 5. Trigger

Nutzer-Frage 2026-07-17 während des slice-071-Release-Prep: ob der neu
geschriebene §4.12-Block den Benutzerhandbuch-Standard §2 einhält. Er tat es
nicht — er begründete die Fähigkeit aus d-checks Modul-Landschaft heraus
(„`matrix` prüft die Richtung, `trace.coverage` die Abdeckung, keines vergleicht
die Mengen") statt aus der Lesersituation. Der Block wurde in slice-071
nachgezogen (Commit `d3a8757`); die anschließende Frage „gibt es weitere solche
Stellen?" ergab ein Audit des gesamten Handbuchs mit sieben Befunden (B-1…B-8) —
dieser Slice arbeitet sie ab.

## 6. Sub-Area-Modus-Begründung

GF (Repo-Default): Der Standard führt, die Doku folgt. `docs/user/` ist kein
Brownfield-Bestand mit eigener Spec — die Zielform steht im
Benutzerhandbuch-Standard und ist im Bestand (§4.10, §4.12-Kreuzverweis) bereits
zweimal gelebt; der Slice zieht den Rest nach.

## 7. Closure-Notiz (nach `done/`)

_Ausstehend — wird bei Abschluss mit Commit-Hash, Review-Verdikt und Lerneintrag
gefüllt._
