# Review slice-164 — Der Nachtlauf hat keinen Adressaten

**Gegenstand:** [slice-164](../plan/planning/done/slice-164-nachtlauf-kadenz.md), Stand des Feat-Commits vor der Nacharbeit.
**Datum:** 2026-08-27. **Reviewer:** unabhängiger Subagent.

---

## Urteil

**Schließbar nach Nacharbeit.** Der Kern trägt: die Reihenfolge *Adressat →
Takt → Kanal* ist richtig gefahren, das Werkzeug existiert und läuft, `make
gates` ist grün (eigener Lauf). `nightly-state` ist korrekt **nicht** in
`gates`, und §3 *„Kein Fremd-Dienst mit Geheimnis"* ist eingehalten.

## Befunde

| # | Rang | Befund |
|---|---|---|
| F-1 | **HIGH** | Der Zeit-Beleg in §7 vergleicht **UTC gegen Lokalzeit** — in einer Basis liegt die Pin-Hebung **2 min 49 s vor** dem Lauf, nicht zwei Stunden danach. Der Beleg widerlegt die Schlussfolgerung, statt sie zu stützen. Richtig belegbar über den `head_sha` |
| F-2 | **HIGH** | *„Jeder ohne Fremd-Dienst verfügbare Kanal …"* ist eine **All-Aussage**, am Bestand widerlegt (Status-Badge; ein einziges fortgeschriebenes Issue statt eines je Lauf). Und *„gemessen"* benennt **keine Messung** |
| F-3 | **HIGH** | Ein **laufender** Nachtlauf wird als `ROT (null)` gemeldet — die API liefert JSON-`null`, nicht ein fehlendes Feld; der SKIP-Zweig ist toter Code. Falsch-Rot in genau dem Werkzeug, das Wegklicken verhindern soll |
| F-4 | **HIGH** | `AGENTS.md` §3.1 zählt weiter **zwei** Netz-Skripte und sagt *„Beide"* — das neue ist das dritte, und sein Kopf beruft sich auf genau diesen Absatz. Derselbe Befund war am selben Tag schon geschlossen |
| F-5 | MEDIUM | Kein netzloser Prüfeinstieg und keine Probe — die Parse-Semantik ist am Repo nicht prüfbar, anders als bei den Schwester-Skripten |
| F-6 | MEDIUM | Kein Alters-Test: ein grüner Lauf von vor sieben Monaten meldet `gruen`, während der Wächter still ausgefallen ist |
| F-7 | MEDIUM | Das ganze Nicht-`success`-Vokabular bekommt den *„planmäßige Meldung"*-Rat — auch `cancelled`, `timed_out`, wo er falsch ist |
| F-8 | MEDIUM | `MR-053`s Pflichtfeld trägt einen **Datei-Link ohne Anker**; die Vorlage sagt ausdrücklich *„ein Datei-Link benennt keine Regel"* |
| F-9 | MEDIUM | *„dieselbe Sparsamkeit wie bei den Frische-Achsen (`DC-QA-03`)"* ist an beiden Enden ungenau — jene sind gerade **nicht** API-basiert, und `DC-QA-03` gilt dem **Produkt**, nicht dem inneren Lauf |
| F-10 | MEDIUM | Hartkodierter Repo-Slug und zwei undokumentierte Umgebungs-Schalter (§3.8); drei ungeprüfte Eingaben unbenannt |
| F-11 | MEDIUM | `AGENTS.md` §5 und `MR-053` haben **verschiedene** Geltungsbereiche — `open/slice-160` fällt in die Lücke |
| F-12 | LOW | *„Immer Exit 0"* ist als Verhalten in Ordnung, als **Begründung** ein Fehlschluss; und alles landet auf **stdout**, wo `pin-freshness.sh` stderr nutzt |
| F-13 | LOW | Der jüngste „Nachtlauf" lief mittags — die Abtastrate ist nachweislich unzuverlässig |
| F-14 | LOW | Doppelte Leerzeile im `Makefile` |

**Was der Reviewer zu brechen versuchte und nicht brechen konnte:** das
`tr ','`-Parsing hält — ein Komma **im Wert** kann die Feld-Extraktion nicht
verschieben, weil das Muster an beide Anführungszeichen verankert ist und ein
literales `"` im JSON-String maskiert wäre. *„Immer Exit 0"* hält ebenfalls,
gegen leeres JSON, HTML und `total_count: 0`. **Der Fehler sitzt nicht in der
Trennung, sondern in der fehlenden Wert-Prüfung.**

## Erledigung

Alle vierzehn Befunde sind eingearbeitet; die vier HIGH sind **eigens
nachgemessen**.

- **F-1** Beleg durch den `head_sha` ersetzt; die Vorhersage über den nächsten
  Lauf steht jetzt als Vorhersage da.
- **F-2** All-Aussage und *„gemessen"* zurückgenommen; drei betrachtete
  Kandidaten mit Grund, Liste ausdrücklich unvollständig, Nicht-Kandidaten
  benannt.
- **F-3**, **F-7** im `case` behoben, plus Form-Prüfung des Zeitstempels.
- **F-4** §3.1 auf *„Alle drei"* gezogen.
- **F-5** `--parse` und `--selftest` ergänzt — **und der erste Selbsttest fand
  sofort einen Fehler**, den der Live-Lauf verdeckte (`"}` nicht gestrippt).
- **F-6**, **F-10** als vier benannte Grenzen in beiden Deklarationen.
- **F-8**, **F-9**, **F-11**, **F-12**, **F-14** direkt behoben.
- **F-13** als Teil der Alters-Grenze.
