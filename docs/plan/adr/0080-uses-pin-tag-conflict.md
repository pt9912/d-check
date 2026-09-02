# ADR-0080: Ein SHA trägt innerhalb der Scan-Menge überall denselben Tag-Kommentar

**Status:** Accepted

**Datum:** 2026-09-02

**Autor:** pt9912

**Bezug:** [`DC-FA-WF-001`](../../../spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in)
(die erweiterte Anforderung);
[ADR-0072](0072-workflows-modul.md) (das Modul, seine Hexagon-Aufteilung und
seine bisherigen sechs Grund-Codes);
[ADR-0071](0071-lokale-workflow-referenz-rechte-pruefung.md) (dieselbe
Bauform: eine bestehende Bedingung bekommt eine Nachbar-Bedingung, kein neues
Modul)

**Schärft:** [`DC-FA-WF-001`](../../../spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in)

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v5.18.0/regelwerk/modul-04-adrs.md).

---

## Kontext

**`uses-pin-untagged` prüft nur, *dass* ein Tag-Kommentar dasteht — nicht, *ob*
er über die Scan-Menge hinweg konsistent ist.** Ein SHA, der an zwei Stellen
mit unterschiedlichen Kommentaren steht, sagt zwei sich widersprechende Dinge
über denselben Commit, und mindestens eine der beiden Stellen liegt falsch.
Bisher meldet keine der sechs bestehenden Bedingungen das.

**Anlass ist ein eingehender Change Request** der Schwester-Anwendung
`a-check` (`docs/plan/cr/2026-08-30-cr-a-check-uses-tag-kohaerenz.md`). Ihr
eigenes, abzulösendes Bash-Sensor-Äquivalent fand genau diesen Fall: derselbe
SHA von `docker/login-action` zweimal gepinnt, einmal mit `# v4.2.0`, einmal
mit `# v3.6.0` — aufgelöst über die GitHub-API (`v4.2.0` korrekt, `v3.6.0`
falsch). Ein zweites, unabhängiges Vorkommen desselben Musters lag
**dateiübergreifend** vor (`actions/checkout`, `v5.0.0` gegen `v6.0.2`). Der CR
benennt ausdrücklich, was er **nicht** beantragt: kein Urteil, welcher Wert
stimmt (das bliebe Netz), keine Ausweitung auf ungetaggte Referenzen.

**Am eigenen Bestand gemessen, vor dem Bau:** `.github/workflows/` dieses
Repos führt drei distinkte SHAs, jeder mit genau einem Tag-Kommentar;
`actions/checkout` erscheint fünfmal mit identischem Kommentar. Der Fall
liegt hier nicht vor — die Regel ist für einen fremden, nicht den eigenen
Bestand gebaut ([`BEO-011`](../planning/observations.md)).

**Die Bedingung ist die erste des Moduls, die nicht je Datei urteilt.** Alle
sechs bestehenden Bedingungen (Pin-Form, Tag-Präsenz, lokale Existenz, beide
Rechte-Bedingungen, Parse-Fehler) werten eine Referenz oder eine Datei
isoliert aus. Ein Tag-Konflikt existiert erst im Vergleich zweier Referenzen
— im Extremfall über zwei Dateien hinweg.

## Entscheidung

**Ein neuer Grund-Code `uses-pin-tag-conflict`.** Alle fremden, voll
gepinnten, tag-kommentierten Referenzen der gesamten Scan-Menge werden nach
ihrem SHA gruppiert; führt eine Gruppe mehr als einen distinkten
Tag-Kommentar-Text, meldet **jede** Zeile der Gruppe.

1. **Dateiübergreifend, nicht je Datei.** Die Sammlung läuft über **alle**
   Kandidaten-Dateien hinweg, bevor gruppiert wird — der CR verlangt
   ausdrücklich diese Reichweite, sein zweiter Realfall lag zwischen zwei
   Dateien.
2. **Ein Befund je beteiligter Zeile, nicht je SHA.** Ein SHA mit zwei
   Kommentar-Werten über drei Zeilen ist **drei** Befunde. Jede Zeile trägt
   ihre eigene Reparaturstelle, und eine Sammel-Meldung verstünde nicht, wie
   viele Stellen zu korrigieren sind.
3. **Die Meldung nennt beide (bzw. alle) widersprechenden Werte, sortiert.**
   Ohne sie müsste, wer den Befund liest, selbst suchen, worin der
   Widerspruch besteht.
4. **Wiederholung ist kein Befund.** Derselbe SHA mit identischem Kommentar
   über beliebig viele Zeilen — der häufige Fall, `actions/checkout` fünffach
   in diesem Repo — bleibt befundfrei. Die Bedingung zählt **distinkte**
   Werte, nicht Vorkommen.
5. **Referenzen ohne Tag-Kommentar gehen nicht ein.** Sie tragen bereits
   `uses-pin-untagged`; ohne Kommentar gibt es nichts, das widersprechen
   könnte, und eine Doppelmeldung wäre eine zweite Zusage über denselben
   Mangel.
6. **Welcher der widersprechenden Werte stimmt, sagt die Regel nicht.** Das
   wäre dieselbe Gültigkeitsfrage, die schon `uses-pin-missing`/
   `uses-pin-untagged` an das Netz abgibt (`AGENTS.md` §3.9). Die Regel
   deckt die **Deklarations**-Klasse „widersprüchlich", nicht die
   Entscheidung „richtig".
7. **Benannte Grenze: zwei legitime Tags auf demselben Commit.** Ein
   unversionierter Major-Tag (`v4`) neben seinem spezifischen Patch-Release
   (`v4.2.0`), kurz nach dessen Erscheinen auf denselben Commit gezeigt, sind
   textuell verschieden und lösen **denselben** Befund aus wie ein echter
   Fehler. Die Unterscheidung wäre ein Versions-Kompatibilitäts-Urteil
   („umfasst dieser Tag jenen"), keine Deklarations-Prüfung — sie bleibt
   bewusst außerhalb, wie der CR es auch nicht beantragt.
8. **Kein neues Konfigurationsfeld.** Die Bedingung teilt die
   Aktivierungsschranke von `uses-pin-untagged`: `workflows.dir` gesetzt und
   das Modul aktiv. Ein eigener Schalter wäre eine Fähigkeit, die niemand
   beantragt hat, gegen eine Bedingung, die sich nicht abschalten lässt, ohne
   das ganze Modul abzuschalten.

## Verglichene Alternativen

| Option | Pro | Contra | Verworfen, weil |
|---|---|---|---|
| **A: eine Meldung je SHA-Gruppe statt je Zeile** | kompakter Befundsatz | verschweigt, wie viele Stellen zu korrigieren sind; widerspricht der bestehenden Zusage „ein Befund auf **ihrer** Zeile" für jede andere Bedingung des Moduls | eine Reparatur pro betroffene Zeile ist die Reparatur, die tatsächlich ansteht |
| **B: `v4`/`v4.2.0`-Kompatibilität erkennen (SemVer-Präfix-Vergleich)** | löst die benannte Grenze auf, keine Falsch-Positiven bei Major/Patch-Paaren | verlangt eine SemVer-Parse- und Kompatibilitäts-Logik, die der CR nicht beantragt hat, und deckt nur den SemVer-Fall — ein Alias-Tag ohne SemVer-Form (`stable`, `latest-v4`) bliebe ungedeckt | ein Versions-Urteil ist eine andere Frage als eine Deklarations-Prüfung; der CR benennt das selbst als Nicht-Antrag |
| **C: Gültigkeit gegen die GitHub-API prüfen (welcher Tag stimmt)** | löste die Ursache, nicht nur das Symptom | Netz — widerspricht der hermetischen Zusage des Moduls (`AGENTS.md` §3.1/§3.8, [ADR-0072](0072-workflows-modul.md)); gehörte, wenn überhaupt, zur Freshness-Familie | dieselbe Grenze wie bei `uses-pin-missing`: Gültigkeit ist Netz, Form ist hermetisch |
| **D: nichts tun, den CR ablehnen** | kein neuer Code, keine neue Fläche | der gemessene Anlass ist real (zwei belegte Fälle bei `a-check`) und die Reichweite ist gemessen klein (ein Grund-Code, keine neue Konfiguration); Ablehnen verwehrt eine Fähigkeit, die für den eigenen Bestand ohnehin blind bleibt (0 Konflikte) | die Kosten sind gemessen niedrig, der Nutzen für einen adoptierenden Fremd-Bestand belegt |

## Konsequenzen

**Positiv.** Ein SHA-Tag-Widerspruch — belegt an einem Fremd-Bestand, real
und zweimal unabhängig aufgetreten — wird jetzt gemeldet, dateiübergreifend.
Wiederholung bleibt befundfrei, die Meldung nennt die widersprechenden Werte,
und die Reparaturstelle ist die betroffene Zeile.

**Negativ, und benannt.** Zwei legitime Tags auf demselben Commit (Major
neben Patch) melden wie ein echter Fehler — ein Adopter, der diese Praxis
pflegt, sieht Befunde, die keine Reparatur brauchen. Das ist derselbe Preis,
den `uses-pin-missing`/`uses-pin-untagged` mit der Gültigkeitsfrage schon
zahlen: Form, nicht Wahrheit.

**Der Grund-Code-Raum des Moduls wächst auf sieben.** Ein Sammel-Code
schiede aus: die Befund-Deduplikation vergleicht (Datei, Zeile, Regel, Ziel,
Grund), und zwei verschiedene SHA-Konflikte auf derselben Zeile gäbe es
ohnehin nie — pro Zeile genau ein SHA.

**Am eigenen Bestand ohne Wirkung.** Gemessen 0 Konflikte vor und nach dem
Bau; die Fähigkeit ist für einen Fremd-Bestand gebaut und bleibt hier
unbelegt bis zum ersten echten Treffer.

## Fitness Function (falls maschinell prüfbar)

**Ja, und die Proben sind gefahren.** Kern-Tests decken: denselben SHA mit
zwei Kommentar-Werten über drei Zeilen einer Datei (drei Befunde);
denselben SHA über zwei Dateien mit unterschiedlichem Kommentar (ein Befund
je Datei); denselben SHA mit identischem Kommentar über fünf Zeilen/zwei
Dateien (null Befunde); die Gruppierung direkt (`checkTagConflicts`,
Determinismus über sortierte SHA-Schlüssel, [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).

**Eine Verdrahtungs-Probe belegt die Kopplung an den Einstiegspunkt**
(Prozedur aus [`BEO-023`](../planning/observations.md)): ein Konflikt, der
nur über `CheckWorkflows` selbst sichtbar wird (nicht nur über die isoliert
aufgerufene Gruppierungs-Funktion), beweist, dass die Sammel-Schleife
tatsächlich an `checkTagConflicts` angeschlossen ist — die Lücke, die
[`BEO-023`](../planning/observations.md) an slice-182 <!-- d-check:status-provenance --> benannt hat
(ein Konfigurationsschlüssel, der nie ins Kernmodell durchgereicht wurde,
weil kein Test den vollen Pfad nahm).

**Am eigenen Bestand gemessen:** `make workflow-pins` gegen dieses Repo
liefert vor und nach dem Bau **0** Befunde dieses Grund-Codes (3 distinkte
SHAs, je ein Tag-Kommentar; `actions/checkout` 5× identisch).

## Re-Evaluierungs-Trigger

**Wenn die benannte Grenze (Major-/Patch-Tag-Paar) am eigenen oder einem
adoptierenden Bestand real auftritt.** Dann ist zu entscheiden, ob eine
SemVer-Präfix-Erkennung (Option B) den Aufwand rechtfertigt — bis dahin ist
sie unbelegt.

**Wenn der Sender des CR eine dritte, verwandte Bedingung nachreicht** (etwa
eine Ausweitung auf ungetaggte Referenzen, die der aktuelle CR ausdrücklich
ausschließt) — dann ist neu zu entscheiden, ob das Modul wächst oder die
Grenze bestätigt wird.

## Geschichte
