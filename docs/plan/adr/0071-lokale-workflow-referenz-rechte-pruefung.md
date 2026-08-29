# ADR-0071: Die lokale Workflow-Referenz wird auch auf ihre Rechte-Anforderung geprüft

**Status:** Accepted (teil-superseded: ADR-0072)

**Datum:** 2026-08-29

**Autor:** pt9912

**Bezug:** [ADR-0068](0068-lokale-workflow-referenzen-ohne-pin.md) (die Ausnahme,
die den Pin-Check durch die Existenz-Frage ersetzt — diese ADR **ergänzt** sie,
sie löst sie nicht ab),
[`AGENTS.md`](../../../AGENTS.md) §3.9 (die Hard Rule) und §3.8 (die Regel, an
der der Wächter scheiterte)

**Schärft:** [`AGENTS.md`](../../../AGENTS.md) §3.9 — Prozess-/Harness-Entscheid
ohne Spec-Stratum.

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v5.12.0/regelwerk/modul-04-adrs.md).

---

## Kontext

**Das Gate war grün, und das Release lief nicht an.** Am Tag-Push von `v0.66.0`
(Run 33235586073) meldete GitHub `startup_failure` — **null Jobs**, kein Log,
keine Check-Runs. `release.yml` führt `permissions: {}` im Kopf; der Job
`hub-description` erbte das und ruft
[`hub-description.yml`](../../../.github/workflows/hub-description.yml) auf,
dessen Job `contents: read` verlangt. Ein aufgerufener Workflow bekommt nur, was
der aufrufende Job selbst führt; verlangt er mehr, lehnt GitHub den **ganzen
Lauf vor dem ersten Job** ab.

`make workflow-pins` meldete währenddessen *„8 uses:-Einträge geprüft, davon 1
lokale Referenz(en) ohne Pin-Pflicht"*.

**Das ist die Klasse, die dieses Repo als Hard Rule führt.**
[`AGENTS.md`](../../../AGENTS.md) §3.8 verlangt von jedem Modul die Antwort auf
*„welche Eingaben liest es, die es nicht scannt — und gilt dort dieselbe
Zusage?"*. Der Wächter liest die `uses:`-**Zeile**; die referenzierte **Datei**
öffnete er (`[ -f "$target" ]`) und sah sie nie an. Im Beobachtungs-Register
steht die Klasse als [`BEO-004`](../planning/observations.md).

**Was [ADR-0068](0068-lokale-workflow-referenzen-ohne-pin.md) zusagte und was
nicht.** Ihr Kern-Entscheid ist unverändert richtig: eine lokale Referenz kann
keinen SHA tragen und braucht keinen, weil sie auf denselben Commit auflöst wie
ihr Aufrufer — sie ist **stärker** gebunden als ein Pin. An die Stelle des
Pin-Checks setzte sie die Existenz-Frage und schloss: *„Die Ausnahme lässt damit
**keinen Eintrag ungeprüft**."* Über den **Pin** ist dieser Satz wahr; über die
**Lauffähigkeit** ist er zu weit. Die Rechte-Frage kommt in ADR-0068 nirgends
vor — sie war nicht verworfen, sondern nicht im Blick.

## Entscheidung

Der `./`-Zweig von
[`workflow-pins.sh`](../../../tools/harness/workflow-pins.sh) stellt an dieselbe
Referenz eine **zweite** Frage: hat der aufrufende Job die Rechte, die das Ziel
verlangt?

1. **Zwei Bedingungen, beide urteilsfrei.**
   - `uses-local-perms-undeclared`: das Ziel verlangt Rechte, der aufrufende
     **Job** trägt kein eigenes `permissions:`. Das ist der Fall, der
     ausgefallen ist — die stille Vererbung vom Workflow-Kopf.
   - `uses-local-perms-narrow`: der Aufrufer führt einen geforderten Scope
     niedriger als das Ziel ihn verlangt (`none` < `read` < `write`). Ein Scope,
     den der Aufrufer **nicht nennt**, zählt als `none` — das ist GitHubs eigene
     Semantik, kein Zugeständnis.
2. **Fail-closed bei jeder Form, die der Wächter nicht sicher liest**
   (`uses-local-perms-unreadable`). `permissions: read-all`/`write-all`, ein
   Flow-Mapping mit Inhalt, ein Anker, eine Job-Einrückung außerhalb der
   2/4-Form: alles meldet sich, statt zu passieren. **Ein Wächter, der bei
   unbekannter Schreibweise schweigt, wäre derselbe stille Grün-Pfad, den diese
   Prüfung schließt.**
3. **Kein YAML-Parser als vierte Toolchain.** Die Zerlegung ist `awk` über die
   Block-Form — das trägt der Bestand, und alles darüber ist ein **Entscheid**,
   kein Werkzeug nebenbei
   ([`MR-046`](../../../harness/conventions.md#mr-046)). Trüge die Block-Form
   nicht weit genug, wäre der Umfang zu schneiden, nicht die Toolchain zu
   weiten.
4. **Die Prüfung bleibt im selben Wächter, statt ein eigenes Target zu
   bekommen.** Sie gilt derselben Zeile, wird an derselben Stelle aufgelöst und
   teilt die Datei-Auswahl; ein zweites Target hieße, `.github/workflows/`
   zweimal zu durchlaufen und die Pin- und die Rechte-Aussage getrennt
   verfallen zu lassen.
5. **Proben statt Erinnerung.** `bash tools/harness/workflow-pins.sh --selftest`
   fährt sieben konstruierte Fälle mit Erwartung und Ergebnis. Der Anlassfall
   ist **ein** Exemplar; eine Regel, die nur an ihrem Anlass geeicht ist, trägt
   nicht über ihn hinaus ([`BEO-011`](../planning/observations.md)).

**Was diese ADR NICHT tut:** sie löst
[ADR-0068](0068-lokale-workflow-referenzen-ohne-pin.md) nicht ab. Deren
Entscheid — lokale Referenz ohne Pin, dafür Existenz-Prüfung — steht unverändert;
korrigiert wird nur ihre zu weite Reichweiten-Aussage.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v5.12.0/regelwerk/modul-04-adrs.md)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, der Fall ist behoben | kein Code; der Bestand trägt genau **eine** lokale Referenz | die Reparatur ist am Symptom: der nächste lokale Aufruf mit engeren Aufrufer-Rechten fällt genauso aus, bei grünem Gate. Genau die Bauform, die dieses Repo als [`BEO-013`](../planning/observations.md) führt |
| B — die Rechte im aufgerufenen Workflow weglassen, damit nichts verlangt wird | kein Wächter nötig | dreht die Härtung um: `permissions` je Job ist die Least-Privilege-Zusage dieses Repos. Ein Defekt im Wächter wäre mit einer Schwächung der CI „behoben" |
| C — ein eigenes Target `workflow-permissions` | trennt die zwei Fragen sauber | zweiter Durchlauf über dasselbe Verzeichnis, zweite Datei-Auswahl, zwei Aussagen, die getrennt verfallen — und die Frage entsteht ohnehin nur dort, wo der Pin-Wächter die Referenz schon aufgelöst hat |
| D — echter YAML-Parser (Go-Modul oder Fremd-Tool) | trägt jede Schreibweise | eine vierte Toolchain für einen Fall ([`MR-046`](../../../harness/conventions.md#mr-046)); ein Go-Modul im Produkt wäre zudem die falsche Anforderung — `d-check` prüft Dokumentation gegen Zustand, nicht Actions-Semantik |
| **E — zweite Frage im selben Wächter, Block-Form in `awk`, fail-closed bei Unlesbarem (gewählt)** | fängt den Anlassfall (retro gemessen); kein neues Werkzeug; die unlesbare Form meldet statt zu schweigen; Proben sind wiederholbar | deckt **eine** Fehlerklasse, nicht die Lauffähigkeit; die Block-Form-Grenze ist real und muss benannt bleiben |

## Konsequenzen

- **Positiv, retro gemessen:** gegen den Stand vor dem Fix (`8beac2d`) meldet der
  Wächter `uses-local-perms-undeclared` auf Zeile 268 und endet mit Exit 1;
  gegen den heutigen Stand ist er grün. Die Regel fängt ihren Anlassfall.
- **Positiv:** die in §3.8 gestellte Frage ist für diesen Wächter beantwortet —
  die Zieldatei ist jetzt eine Eingabe, über die er etwas verspricht.
- **Negativ, benannt:** der Wächter deckt **eine** Fehlerklasse. Ein grüner Lauf
  sagt „diese Klasse liegt nicht vor", nicht „der Workflow läuft". Wer mehr
  hineinliest, wiederholt genau den Fehler von ADR-0068.
- **Negativ, benannt:** die Zerlegung ist eine Näherung über die Block-Form.
  YAML-Anker, Flow-Mappings mit Inhalt, Mehrfach-Dokumente und abweichende
  Einrückung sind **nicht** gedeckt — sie melden sich als unlesbar, was
  fail-closed richtig, aber im Zweifel ein Falsch-Positiv ist.
- **Negativ, benannt:** `make workflow-pins` läuft in `make gates`. Eine
  Schreibweise, die der Wächter nicht kennt, macht damit den **inneren Loop**
  rot, nicht erst die CI. Das ist gewollt und trotzdem ein Preis
  ([`BEO-018`](../planning/observations.md)).

## Fitness Function (falls maschinell prüfbar)

`bash tools/harness/workflow-pins.sh --selftest` — sieben Proben mit Erwartung
und Ergebnis je Fall: die stille Vererbung (der v0.66.0-Fall), der ausreichend
deklarierte Aufrufer, `read` gegen `write`, der beim Aufrufer fehlende Scope,
das Ziel ohne Forderung, die unlesbare Form (`read-all`) und die Kombination
„keine Forderung, keine Deklaration". Netzlos, kein git, kein Docker.

`make workflow-pins` selbst bleibt der Bindepunkt in `make gates`.

**Was keine Fitness Function prüft:** ob die Block-Form-Zerlegung für einen
künftigen Workflow ausreicht. Das meldet sich als `uses-local-perms-unreadable`
— und ist dann ein Entscheid, kein Bug.

## Re-Evaluierungs-Trigger

**Der erste `uses-local-perms-unreadable`-Befund an einer legitimen
Schreibweise.** Dann ist zu entscheiden, ob die Block-Form erweitert wird oder
ob der Fall die Grenze aus [`MR-046`](../../../harness/conventions.md#mr-046)
erreicht — ein echter Parser ist ein Entscheid, kein Nachbessern.

**Zweiter Trigger: eine zweite lokale Workflow-Referenz im Repo.** Solange es
**eine** gibt, ist jede Aussage über die Klasse an einem Exemplar geeicht
([`BEO-011`](../planning/observations.md)); mit der zweiten ist zu prüfen, ob
die Regel auch dort greift, ohne dass sie dafür gebaut wurde.
