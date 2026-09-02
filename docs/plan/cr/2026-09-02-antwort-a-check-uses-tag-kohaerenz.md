# Antwort auf den `a-check`-CR vom 2026-08-30 — ein SHA, ein Tag-Kommentar

**Absender der Antwort:** d-check
**Datum:** 2026-09-02
**Bezug:** [eingehender CR](2026-08-30-cr-a-check-uses-tag-kohaerenz.md); Vorgänger
[CR 4](2026-08-30-cr-a-check-leermenge.md) samt
[Antwort](2026-08-30-antwort-a-check-leermenge.md)
**Ergebnis:** **in der Sache und in der Form angenommen** — eine Ergänzung, die
der Antrag nicht stellte, sonst unverändert
**Einordnung:** MINOR — Lastenheft `0.83.0`, Begründung in
[ADR-0080](../adr/0080-uses-pin-tag-conflict.md)

---

## Vorab

Der Antrag trägt, und er trägt in der Form, in der ihr ihn gestellt habt:
`uses-pin-tag-conflict`, dateiübergreifend, ein Befund je beteiligter Zeile,
Wiederholung ist kein Befund, welcher Wert stimmt bleibt Netz. Eure
Einordnung — *„für eine Zusage ist ‚vorhanden' die schwächste denkbare
Prüfung"* — ist die Begründung, die wir übernommen haben, nicht bloß zur
Kenntnis genommen.

Euer Beleg-Anhang ist nachgelesen: die Herkunfts-Spalte hält, `d-check
v0.69.0` gegen `a-check` meldete zur genannten Zeit genau
`uses-local-perms-undeclared`, und der eigene Bestand ist mit euren Zahlen
deckungsgleich — drei distinkte SHAs, `actions/checkout` fünffach mit
identischem Kommentar.

## Was wir ergänzt haben, das ihr nicht beantragt habt

**Die Grenze bei zwei legitimen Tags auf demselben Commit ist jetzt benannt,
nicht nur stillschweigend in Kauf genommen.** Ein Major-Tag (`v4`) neben
seinem Patch-Release (`v4.2.0`), kurz nach dessen Erscheinen auf denselben
Commit gezeigt, ist textuell ein Widerspruch und meldet als einer — obwohl
beide Werte wahr sind. Ihr habt diesen Fall nicht angesprochen; er stand für
uns bei der Umsetzung an, weil die Regel ihn nicht unterscheiden kann, ohne
eine SemVer-Kompatibilitäts-Frage zu beantworten, die niemand beantragt hat.
Wir haben ihn in Lastenheft, Spezifikation und Handbuch als **benannte
Grenze** aufgenommen, nicht als Ausnahme in der Regel — die Regel meldet ihn
weiterhin, jetzt mit Ansage.

## Angenommen, unverändert

- **Grund-Code-Name `uses-pin-tag-conflict`** — euer Vorschlag, übernommen.
- **Dateiübergreifend.** Die Gruppierung läuft über alle Kandidaten-Dateien,
  nicht je Datei — genau der Fall, den euer zweites Beispiel
  (`actions/checkout`, `v5.0.0` gegen `v6.0.2`) über zwei Dateien zeigt.
- **Ein Befund je beteiligter Zeile.** Ein SHA mit zwei Kommentar-Werten über
  drei Zeilen meldet dreimal, nicht einmal — eure Formulierung *„weil keine
  von ihnen für sich falsch ist"* ist die Begründung, die wir übernehmen.
- **Wiederholung ist kein Befund.** Identischer Kommentar über beliebig
  viele Zeilen bleibt befundfrei — geprüft an eurem und unserem Bestand.
- **Keine Aussage, welcher Wert stimmt.** Bleibt Netz, wie ihr es selbst
  ausgeschlossen habt.
- **Keine Ausweitung auf Referenzen ohne Tag-Kommentar.** Die deckt
  `uses-pin-untagged` bereits; wir haben nichts hinzugefügt.
- **Kein neues Konfigurationsfeld.** Die Bedingung teilt die
  Aktivierungsschranke von `uses-pin-untagged` — `workflows.dir` gesetzt,
  Modul aktiv. Ein eigener Schalter wäre eine Fähigkeit, die niemand
  beantragt hat.

## Am eigenen Bestand gemessen

`.github/workflows/` dieses Repos liefert vor und nach dem Bau **0** Befunde
dieses Grund-Codes — deckungsgleich mit eurer Messung, dass die Regel bei uns
heute nichts fände. Der Antrag ist, wie ihr selbst schreibt, eine
Regressions-Bremse, kein Bestandsräumer.

## Was wir nicht tun

**Wir bauen keine SemVer-Kompatibilitäts-Prüfung** für die benannte Grenze
(Major- neben Patch-Tag). Sie ist [in der ADR](../adr/0080-uses-pin-tag-conflict.md#re-evaluierungs-trigger)
als Re-Evaluierungs-Trigger geführt: tritt der Fall real auf — bei euch oder
bei uns —, ist neu zu entscheiden, ob sich der Aufwand lohnt.

**Wir prüfen nicht gegen die GitHub-API, welcher Wert stimmt** — dieselbe
Gültigkeitsfrage, die `uses-pin-missing`/`uses-pin-untagged` schon an das
Netz abgeben, und aus demselben Grund: das Modul bleibt hermetisch.
