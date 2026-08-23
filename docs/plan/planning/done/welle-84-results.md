# Welle 84 — Zehn Regeln befragt, zwei Wächter gebaut, und sechsmal ein Exit-Code falsch gelesen — Closure-Notiz

**Welle:** welle-84-durchsetzung
**Abschluss:** 2026-08-23
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief · Steering-Loop · Register-Lese-Schritt.

- [slice-133](slice-133-baseline-sensor-verdrahten.md) — **der verwaiste Sensor
  ist eingesteckt.** `make baseline-verify` ist erstes Glied von `make gates`,
  `make baseline-freshness` läuft im eigenen Nachtlauf `upstream-drift.yml`,
  getrennt von `ci.yml`. Gebaut wurde nichts — das Skript konnte all das schon
  und wurde von nichts gerufen.
- [slice-132](slice-132-hard-rule-zensus.md) — **der Zensus.** Elf Zeilen: neun
  Abschnitte in §3 tragen zehn Regeln, dazu die Botschafts-Regel aus §5. Jede
  Deckungs-Aussage steht auf einem konstruierten Verstoß mit rotem Exit.
- [slice-134](slice-134-nolintlint.md) — **`nolintlint` im Profil.** §3.2 von
  *einseitig* auf *teilgedeckt*.
- [slice-135](slice-135-uses-pin-sensor.md) — **`make workflow-pins`**, zehntes
  Glied von `gates`. §3.9s auflösender Trigger ist eingelöst.
- [slice-136](slice-136-agents-34-klaerung.md) — **eine Zensus-Zeile
  berichtigt**, und mit ihr eine falsche Konformitäts-Annahme:
  [`MR-033`](../../../../harness/conventions.md#mr-033).

## Was hat funktioniert?

**Der konstruierte Verstoß als Beweisform.** Der Zensus hat für jede
Deckungs-Aussage einen Verstoß gebaut und den Exit gelesen — und genau das hat
den schärfsten Fund der Welle sichtbar gemacht: §3.2 verbietet
Inline-Suppressions, und die verbotene Direktive **funktionierte**. Ein Zensus,
der nur gezählt hätte, hätte das nicht gefunden.

**Vier Verdikte statt drei.** *gedeckt · teilgedeckt · werkzeug-lokal ·
einseitig* — die Trennung zwischen den mittleren beiden ist auf Review-Hinweis
entstanden und trägt: §3.1 hat einen Wächter, aber außerhalb der Gates, und das
ist etwas anderes als ein Gate, das die halbe Regel prüft.

**Die Ausweisung als vollwertiges Ergebnis.** Sechs Regeln bleiben einseitig,
und für vier davon ist das kein Mangel: ihre Durchsetzung wäre ein
Heuristik-Wächter. Das steht jetzt im Regeltext statt in niemandes Kopf.

## Was ging anders als geplant?

**Die Lehre der Welle ist eine einzige Bewegung, sechsmal wiederholt: aus einem
richtigen Messwert den falschen Schluss ziehen.**

1. Die erste `//nolint`-Probe nannte den falschen Linter — rot, aber aus anderem
   Grund; beinahe hätte das als *Regel greift* gezählt.
2. §3.4 stand als *gedeckt* da, obwohl die Probe nur eine seiner zwei Aussagen
   brach.
3. Der Slice-Text nannte eine Direktiv-Form *wohlgeformt*, die nie gefahren
   worden war — sie wird gemeldet.
4. Eine DoD verlangte einen Beleg, den §5 desselben Dokuments vier Zeilen weiter
   ausschließt.
5. §3.4s Abwärts-Sperre nennt **fünf** Kategorien; belegt war **eine**, gedeckt
   sind **zwei** — dieselbe Form, die ein Review einen Slice zuvor im selben
   Abschnitt als HIGH gemeldet hatte.
6. Und der schwerste: *„keine Verschärfung"*, geschlossen aus der **Praxis**
   einer Vorlage, während die direkteste Quelle ungelesen blieb.
   `AGENTS.template.md` §3.4 sagt wörtlich, die Sicht *referenziere*
   Modul-Pfade. Wir verbieten sie. Das ist eine Verschärfung, und sie war
   undeklariert.

**Ein Exit-Code ist kein Beleg, solange seine Ursache ungelesen bleibt** — die
DoDs verlangen das seit slice-134 ausdrücklich, und seither ist jeder rote Lauf
mit seiner Meldung belegt, nicht mit seiner Zahl.

**Zwei Wächter hatten den Pfad, den sie verhindern sollten.**
`baseline-verify`s Manifest-Deckung zählte nur innerhalb der zwei Bäume — eine
Datei als Geschwister daneben blieb unsichtbar. `workflow-pins` las nur `*.yml`,
GitHub liest auch `*.yaml`. Beide Male stand im selben Kommentarblock die
Begründung, warum genau das nicht passieren dürfe.

**Ein Risiko zu benennen ersetzt nicht, es zu prüfen.** slice-136 hat sein
eigenes §5-Risiko wörtlich erfüllt, Wort für Wort, nachdem es geschrieben war.

**Zwei Commits waren defekt, und beide fand eine Maschine.** Der Claim-Commit
trug seine gekoppelten Verweise nicht (CI: drei `target-missing`), der
Closure-Body legte einen bereits gelöschten Slice wieder an (Stop-Hook →
`planning-drift`). Beide Male lag zwischen Ursache und Fund ein Commit, in dem
kein Gate lief.

**Und ein `--amend` auf einem bereits gepushten Commit.** Die Regel *„nie
amenden, was auf origin liegt"* habe ich angewandt, ohne zu prüfen, **ob** er
dort lag.

## Steering-Loop-Einträge

- **`AGENTS.md` §3 trägt jetzt an jeder Regel ihren Durchsetzungs-Stand** —
  tragender Gate-Lauf oder ausgewiesene Einseitigkeit, dazu je einen
  Auflösungs-Trigger.
- **Zwei neue Gates**, beide netzlos und beide in `gates`: `baseline-verify` und
  `workflow-pins`. Dazu ein Nachtlauf, der bewusst **nicht** in der CI hängt.
- [`MR-033`](../../../../harness/conventions.md#mr-033) deklariert die Verschärfung, die §3.4 seit jeher trug.
- **Nicht** mechanisiert: die vier Urteils-Regeln. Das ist Ergebnis, nicht Rest.

## Beobachtungs-Register (Zeiger)

Gelesen zur Closure ([`observations.md`](../observations.md)) — dreizehn
Einträge:

- **[`BEO-013`](../observations.md) neu angelegt**, Zähler **1**: *ein
  Unterdrückungs-Marker, der nichts mehr unterdrückt, bleibt stehen und deckt
  lautlos den nächsten Befund.* Für Go ist die Frage seit `nolintlint`
  (`allow-unused: false`) gemessen; für `<!-- d-check:ignore -->` gibt es kein
  Gegenstück. Bestand: **zwölf** aktive Marker, elf eingefroren in `done/`,
  **einer** in einem lebenden Dokument.
- **[`BEO-012`](../observations.md) hat dreimal erneut getroffen** — [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
  über seinen Geltungsbereich zitiert, `nolintlint` in die [ADR-0006](../../adr/0006-lint-profil-solid.md)-Zahl
  hineingezählt, und die Verschärfung aus der falschen Quelle widerlegt. Zähler
  bleibt bei **3**, die Schwelle ist seit welle-83 erreicht. Die Prozedur des
  Eintrags — *lies das Geltungs-Feld, nicht den Titel* — hat in allen drei
  Fällen gegriffen, sobald jemand sie angewandt hat; keiner davon war ich.
- **[`BEO-011`](../observations.md)** traf bei den Zählungen (neun Abschnitte /
  zehn Regeln / elf Zeilen in einem Satz vermischt) und bei einem **erfundenen**
  Restrisiko — der Nachtlauf benachrichtigt sehr wohl. Zähler unverändert.
- **[`BEO-007`](../observations.md)** blieb ohne Vorfall: jeder Beleg-Lauf
  dieser Welle hat seinen Exit direkt gelesen.
- **[`BEO-010`](../observations.md)** bleibt bei 1. Zwei neue Targets sind in
  drei Doku-Flächen eingetragen worden, und `gate-consistency` hat die
  Namens-Hälfte gehalten; die Mengen-Hälfte bleibt gate-blind.
- Alle übrigen Einträge unverändert; keine Streichung.

## Folge-Slices

- **Drei fehlende Token-Klassen** für §3.4s Abwärts-Sperre: Wellen,
  Commit-Hashes, Closure-Daten. Der Trigger steht im Regeltext; der Slice ist
  **noch nicht geschnitten** und gehört in die nächste Welle.
- **Ein Sensor auf Pfad-Token in der Architektur-Sicht** — §3.4s zweiter
  auflösender Trigger. Kein heutiges Modul trägt ihn.
- **`BEO-013` mechanisieren** — das `allow-unused`-Gegenstück für
  `<!-- d-check:ignore -->`. Ob das Produkt die nötige Information hält, ist zu
  **messen**.
- **CR-Kandidat an den Kurs:** *„referenziert Modul-Pfade"* sagt nicht, ob Code-
  oder Dokument-Pfade gemeint sind. Für uns entschieden ([`MR-033`](../../../../harness/conventions.md#mr-033)), für jedes
  Adopter-Repo offen. Auftraggeber-Entscheidung.
- Die aus welle-83 offenen Ränder bleiben offen: `diagrams.scope`-Rückbau, die
  3×-Form von [`BEO-008`](../observations.md), die `citations`-Ventil-Achse.

## Verifikation

- `make gates` nach jedem Slice Exit 0, `make fullbuild` zu jeder Closure Exit 0
  — Exit-Codes direkt gelesen ([`BEO-007`](../observations.md)). Die Kette wuchs
  von acht auf **zehn** Glieder.
- **Fünf unabhängige Reviews**, aus den Kategorie-Summaries der fünf Reporte
  gezählt: **6 HIGH, 13 MEDIUM, 2 LOW, 1 INFO**. Vier der fünf waren
  blockierend; alle Befunde eingearbeitet, jede Einarbeitung nachgemessen.
- **Jede Deckungs-Aussage steht auf einem konstruierten Verstoß**, nicht auf
  einer Zuordnung — insgesamt zwölf Proben über vier Gates, jede mit gelesener
  Ursache.
- **Kein Release.** Die Welle hat kein Produkt-Verhalten geändert; beide neuen
  Gates sind Harness-seitig.
