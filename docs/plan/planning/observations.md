# Beobachtungs-Register

**Status:** Aktiv. **Letzte Änderung:** 2026-08-09.

Regeln dieses Registers: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — wer schreibt (Slice-Closure), wer liest (Welle-Closure =
Lese-Schritt bei 3×; Slice-Planung = Sichtungs-Schritt darunter), wann gestrichen wird,
welche Form ein Beleg hat (Slice-Kennung, so viele wie der Zähler, in `done/`), und dass
eine leere Tabelle `— keine —` trägt statt zu verschwinden.

| Kennung | Beobachtung | Sub-Area | Zähler | Belege | Stand |
|---|---|---|---|---|---|
| BEO-002 | **Eine Semantik-Änderung wird nur im Dokumentkörper nachgezogen, ihre Ränder bleiben stehen** — Historie-Zeilen, Index-Einträge, Konsequenz-Abschnitte, Entscheidungs-Protokolle. Sie sind Prosa und darum gate-unsichtbar: `links`/`anchors`/`ids` prüfen die Auflösung, nicht die Aussage. Der Rand referiert danach eine Fassung, die es nicht mehr gibt, und wirkt dabei so belastbar wie der Körper | `*` | 2 | [slice-093](done/slice-093-closure-note-gate.md), [slice-096](done/slice-096-structure-modul-analyse.md) | offen. Erstmals belegt bei slice-093 (die Befund-Provenance wurde beim neuen `--config`-Flag im Kern nachgezogen, in `sources.go` nicht — ein Review fand es). In slice-096 **zweimal** eingetreten (zählt als **ein** Beleg, weil die Beleg-Form Slice-Kennungen verlangt — die Klasse ist dichter als der Zähler): „erster Treffer" blieb nach der Mehrdeutigkeits-Härtung an drei Stellen stehen, und die Überarbeitung nach dem Review korrigierte Anforderung und Algorithmus, ließ aber Spezifikations-Historie, ADR-Index, ADR-Konsequenzen und das Entscheidungs-Protokoll auf der widerrufenen Fassung. **Bei 3× verkörpern statt weiterzählen:** die naheliegende Form ist eine Regel „vor dem Ändern einer Semantik ihre Spiegel auflisten" in der Autoritäts-Doku — dieselbe Bauform wie die Release-Prep-Aufgabenregel, die eine andere still wachsende Klasse geschlossen hat |
| BEO-001 | Ein Datei-Register und seine Autoritäts-Tabelle driften **unbemerkt** auseinander: [ADR-0047](../adr/0047-matrix-spec-historie-nicht-provenance-exempt.md) war sieben Tage `Accepted`, ohne im [ADR-Index](../adr/README.md) zu stehen (`AGENTS.md` §5 verlangt die Zeile). Kein Gate deckt die Richtung „Artefakt ⇒ registriert" ab — `links` prüft nur die Gegenrichtung (Index-Zeile ⇒ Datei existiert). **Dieselbe Bauform** tragen der Konventionsspeicher-Index und das Wellen-Register der Roadmap; beide sind heute konsistent, aber ebenso ungeprüft | `*` | 1 | [slice-087](done/slice-087-spec-historie-referenzrichtung.md) | offen. Gefunden 2026-08-09 während slice-093, **bewusst nicht** dorthin gezogen (WIP-Limit, und slice-093 ist bereits Produkt-Code + Release). Vorschlag für einen späteren Slice: opt-in Modul (Arbeitsname `registry`) mit `[{dir, file-pattern, authority}]`, **eine** Richtung, ein Grund-Code — Größenordnung `tracked`, deckt alle drei Register mit einem Config-Block. Abgrenzung: **kein** Ausbau von `targets` (dessen Eingabe-Achse sind Makefile-Regeln; ein Datei-Register wäre ein zweiter, disjunkter Eingabe-Modus) |

## Gestrichene Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — wer eine Zeile still löscht, macht sie ununterscheidbar von
einer, die es nie gab; darum hierher verschieben (mit Begründung), nicht löschen.

| Kennung | Beobachtung | Gestrichen am | Warum sie nicht mehr auftreten kann |
|---|---|---|---|
| — keine — | | | |
