# ADR-0092: Der Zeilenumbruch-Lookahead der Link-Extraktion bleibt auf eine Zeile begrenzt — die Linktext-Klammer bleibt strikt zeilenlokal (supersedes ADR-0091)

**Status:** Accepted

**Supersedes:** ADR-0091

**Datum:** 2026-09-27

**Autor:** claude-sonnet-5

**Bezug:** Unabhängiger Review-Report
`docs/reviews/2026-09-27-slice-232-links-zeilenumbruch-review-r1.md` <!-- d-check:status-provenance -->,
Befund R1-H1; slice-232 <!-- d-check:status-provenance -->.

**Schärft:**
[`DC-FA-LINK-001`](../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
(präzisiert die Erweiterung aus ADR-0091, keine neue Anforderung).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

ADR-0091 entschied: `ExtractLinks` gruppiert `[]Line` zu **Absätzen**
(Leerzeile/Fenced-Lücke trennt), fügt deren Zeilen mit `"\n"` zusammen und
ruft `forEachLink` **einmal auf dem ganzen zusammengefügten Absatz** auf.

Der unabhängige Review (R1-H1, HIGH) hat diese Umsetzung empirisch widerlegt:
`matchBracket` (in `internal/hexagon/core/rules/markdown.go`) zählt
`[`/`]`- bzw. `(`/`)`-Tiefe ohne Rücksicht auf Zeilengrenzen. Sobald ein
ganzer Absatz zu einem String zusammengefügt wird, kann ein **unbalanciertes
`[` in gewöhnlicher Prosa** (kein `]` auf seiner Zeile — z. B. „Das Array
a[i steht in der Doku.") mit einer **späteren, unabhängigen**
`]…(…)`-Sequenz auf einer folgenden Zeile desselben Absatzes zu einem
**erfundenen Link** verschmelzen, dessen erfundener Linktext mehrere Zeilen
überspannt. Das ist eine andere, **breitere** Fehlerklasse als die in
ADR-0091 benannte Grenze „Zeilenumbruch im Linktext" (dort als *unpromised,
aber harmlos mitgezogen* eingeordnet) — hier entsteht ein **Befund, der ohne
die Änderung nicht existierte**, aus zwei inhaltlich unabhängigen
Textstellen. Der Review hat den Fall an einem konstruierten Testfall
empirisch bestätigt (siehe Fitness Function unten).

Die Ursache liegt darin, dass ADR-0091 **beide** Klammer-Suchen
(Linktext `[…]` **und** Adresse `(…)`) gemeinsam über den ganzen Absatz
laufen ließ, obwohl der Change Request nur die Adresse betrifft.

## Entscheidung

### Nur die Adress-Klammer darf die Zeile verlassen — um genau eine Zeile

`parseLinkAt` bekommt einen dritten Parameter `next` (die unmittelbare
Folgezeile im selben Absatz, sonst `""`):

- Die **Linktext-Klammer** (`matchBracket(s, start, '[', ']')`) bleibt
  **strikt auf `s`** (die aktuelle Zeile) beschränkt — unverändert
  gegenüber dem Stand vor ADR-0091/-0092. Ein `[` ohne `]` auf derselben
  Zeile ist kein Link-Öffner, Punkt.
- Erst **danach**, wenn `[…](` innerhalb der Zeile feststeht, darf die
  **Adress-Klammer** (`(…)`), falls sie innerhalb von `s` nicht schließt,
  um **genau eine** Zeile (`next`) verlängert werden — die konkrete, vom
  Change Request verlangte Form.

`forEachLink` verarbeitet dadurch weiterhin **eine Zeile nach der anderen**
(kein Absatz-Join mehr); liefert eine Adresse einen Übergriff in `next`
(„spillover"), meldet `forEachLink` die verbrauchte Byte-Zahl zurück, und
`ExtractLinks` schneidet diesen Anteil von der **nächsten** Zeile ab, bevor
sie als „aktuelle" Zeile verarbeitet wird — dieselbe Zeile wird nie zweimal
gescannt. `linkParagraphs` (ADR-0091) entfällt ersatzlos: die
Absatz-Zugehörigkeit wird nur noch für die EINE Nachbarzeile geprüft
(`fencedBlockBetween` + nicht-leer), nicht für einen ganzen Absatz.

**Konsequenz für die in ADR-0091 benannten Grenzen:** Linktext-Umbruch und
Titel-hinter-Umbruch „ziehen" jetzt **nicht mehr automatisch mit** — sie
waren in ADR-0091 als unpromised, aber technisch miterfasst beschrieben;
mit der Beschränkung auf einen Ein-Zeilen-Adress-Lookahead ist Linktext-
Umbruch strukturell ausgeschlossen (die Linktext-Klammer verlässt die Zeile
nie). Das ist **keine neue Grenze**, sondern derselbe bereits benannte
Fall, jetzt technisch garantiert statt zufällig.

**Bilder ziehen weiterhin mit** (unverändert aus ADR-0091): `parseLinkAt`
behandelt Links und Bilder über denselben Codepfad; ein Bild mit
zeilenumbrechender Adresse wird identisch erkannt.

Alle vier Aussagen zu den String-Konsumenten (`ids`, `pins`, `--repair`,
`planning`) aus ADR-0091 bleiben **unverändert gültig** — `ExtractLinkSpans`
ruft `forEachLink(text, "", …)` mit leerem `next` auf, der Spillover-Zweig
greift nie, ihr Verhalten ist byte-identisch zum Stand vor ADR-0091.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **Absatzweise Faltung beibehalten** (ADR-0091, verworfen) | einfachste Umsetzung, bereits implementiert | nachweislich fehlerhaft (R1-H1): erfundene Links aus unabhängigen Textstellen |
| **Ganzen Absatz falten, aber Linktext-Klammer nachträglich validieren** (Fund verwerfen, wenn `[`/`]` über eine Zeilengrenze reicht) | bliebe näher an ADR-0091s Struktur | Reparatur eines fehlerhaften Fundes ist fehleranfälliger als ihn gar nicht erst zuzulassen; zwei Stellen müssten synchron bleiben (Erzeugung + nachträglicher Filter) |
| **Ein-Zeilen-Lookahead nur für die Adress-Klammer, Linktext-Klammer strikt zeilenlokal** (gewählt) | eliminiert die Fehlerklasse **durch Konstruktion** (die Linktext-Klammer verlässt `s` nie, unabhängig von Absatz-Inhalt); kleinerer, präziserer Blast-Radius als „ganzer Absatz"; deckt exakt die vom Change Request verlangte Form | eine Zieladresse, die über **zwei** Zeilenumbrüche verteilt ist, bleibt unerkannt (war auch in ADR-0091 nicht zugesagt: „einen einzigen Zeilenumbruch") |

## Konsequenzen

- `parseLinkAt(s, next string, i int)` (Signaturänderung), `forEachLink(text,
  next string, fn) (spillover int)` (Signaturänderung plus Rückgabewert).
  Beide sind unexported; betroffene Aufrufer außerhalb von `ExtractLinks`:
  `ExtractLinkSpans` (`next=""`), `pins.go`, `sources.go` (beide `next=""`,
  Verhalten unverändert).
- `linkParagraphs` (ADR-0091) entfernt — nicht mehr gebraucht.
- Drei neue Tests in `markdown_test.go`
  (`TestExtractLinks_UnbalancierteKlammerVerschmilztNicht`,
  `TestExtractLinks_ZeilenumbruchHinterKlammer`,
  `TestExtractLinks_ZweiVollstaendigeLinksImSelbenAbsatz`) belegen die
  Korrektur direkt an der ursprünglich fehlerhaften Konstruktion.
- Bestandsmessung wiederholt gegen **alle sechs** `ExtractLinks`-Konsumenten
  (nicht nur die vier im eigenen `.d-check.yml` aktiven) — R1-M1 des Reviews
  bemängelte, dass die ursprüngliche Messung `external`/`tracked` ausließ.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `go test` | `TestExtractLinks_UnbalancierteKlammerVerschmilztNicht` — Rot-Beleg: schlägt an der ADR-0091-Fassung von `markdown.go` fehl (Review hat das unabhängig an einem eigenen Testprogramm nachgewiesen), grün an dieser Fassung | `make test` |
| `go test` | `TestExtractLinks_ZeilenumbruchHinterKlammer`, `TestExtractLinks_ZweiVollstaendigeLinksImSelbenAbsatz` (Grenzwert-Kontrollen) | `make test` |
| `make doc-check` (`--enable external --enable tracked`) | Bestandsmessung über alle sechs Konsumenten, nicht nur die vier aktiven | Closure-Notiz |

## Re-Evaluierungs-Trigger

Zwei Bedingungen, jede für sich hinreichend (identisch zu ADR-0091, dort
weiterhin gültig):

1. Ein Konsument bittet um dieselbe Erkennung für einen der vier
   String-Konsumenten.
2. Ein Konsument bittet um eine Zieladresse, die über **mehr als einen**
   Zeilenumbruch verteilt ist.

Ohne eines von beiden: permanent.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-09-27 | Proposed → Accepted (`slice-232`, Nachzug nach unabhängigem Review R1-H1) |
