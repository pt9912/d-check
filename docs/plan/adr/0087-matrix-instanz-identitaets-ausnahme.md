# ADR-0087: `matrix` bekommt eine Instanz-Identitäts-Ausnahme über Regel-Feld statt Klassen-Feld

**Status:** Accepted

**Datum:** 2026-09-18

**Autor:** pt9912

**Bezug:** [`DC-FA-MTX-003`](../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix),
[`DC-FA-MTX-001`](../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
(Präzedenzfall Supersede-Lineage-Ausnahme)

**Schärft:** [`DC-FA-MTX-001.a`](../../../spec/spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung)
Schritt 7 (die Instanz-ID-Extraktion und der Korrelations-Vergleich — der
eingehende CR und die begleitenden Messungen liegen in
`docs/plan/cr/2026-09-18-cr-eingehend-pg-change-feed-matrix-instanz-identitaet.md` <!-- d-check:status-provenance -->)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Ein eingehender CR (`pg-change-feed`, 2026-09-18,
`docs/plan/cr/2026-09-18-cr-eingehend-pg-change-feed-matrix-instanz-identitaet.md`)
meldet: `matrix` behandelt jede verbotene Token-Referenz gleich, ob eine
Quelldatei die **eigene** Instanz des Ziels zitiert (ein Slice, der seinen
eigenen Review-Report nennt — harmlos, beide wandern gemeinsam ins selbe
Wellen-Archiv) oder eine **fremde** (ein Zitat auf eine andere, ggf. bereits
archivierte Instanz — ein reales Risiko für Archivierungs-Werkzeuge wie
`tools/archive-wave`). Gelesen: `ruleFor`/`classOf` vergleichen ausschließlich
Klassennamen; `tokenFindings` korreliert einen Fund nie mit einer Kennung der
Quelldatei. Der Unterschied ist strukturell nicht ausdrückbar.

Eine begleitende Messung fand einen Präzedenzfall im selben Modul: die
Supersede-Lineage-Ausnahme (`matrix.status.allow-supersede-lineage`,
[`DC-FA-MTX-001`](../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix))
nimmt eine `matrix-inactive`-Kante aus, wenn ein **deklariertes Feld** der
Quelldatei (`**Supersedes:** ADR-NNNN`) das Ziel nennt — dieselbe Grundfigur
(instanzbasierte Freigabe einer sonst verbotenen/inaktiven Kante), aber (a)
nur für die Status-Prüfung, nicht für die Klassen-Regelprüfung, und (b) über
ein deklariertes Feld korreliert, nicht über eine aus dem Dateinamen
extrahierte ID. Ein Slice trägt keine „**MeinReview:**"-Zeile — die
Zugehörigkeit liegt ausschließlich in der Namenskonvention
(`slice-<NNN>-x.md` ↔ `review-slice-<NNN>-y.md`).

## Entscheidung

Wir fügen `matrix.rules[].allow-if-same-id: bool` hinzu (Default `false`).
Ist die Regel gesetzt, wird das **bereits vorhandene** `token`-Regex der
beteiligten Klassen **zweifach** genutzt: wie bisher gegen den Fließtext
(Fund-Erkennung), zusätzlich **einmalig gegen den repo-wurzel-relativen Pfad
der Quelldatei** (Instanz-Ermittlung der Quelle). Trägt das Regex genau eine
Capture-Gruppe, ist ihr Wert die Instanz-ID; stimmen Quell- und Ziel-ID überein
(Trim, case-sensitiv), wird der Fund ausgenommen. Kein neues Klassen-Feld
(`id-pattern` o. Ä.) — die Wiederverwendung von `token` ist die schmalere
Umsetzung gegenüber dem CR-Vorschlag, der ein eigenes Feld skizzierte. Die
Ausnahme wirkt **ausschließlich** auf die Token-Form von `matrix-forbidden`;
Link-Referenzen und `matrix-inactive` sind unberührt. Fehlkonfiguration ist
fail-closed am Config-Rand: `allow-if-same-id: true` auf einer Regel, deren
`from`- oder `to`-Klasse kein `token` mit genau einer Capture-Gruppe trägt, ist
Exit 2 (Regel und fehlende Klasse benannt).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, den CR ablehnen | keine neue Konfigurationsfläche | der gemeldete Unterschied bleibt strukturell unausdrückbar; ein real aufgetretener Workaround-Versuch (pauschale Klassenregel) hatte bereits 11 harmlose Selbst-Zitate fälschlich getroffen (CR-Anlass) |
| B — neues Klassen-Feld `id-pattern` (CR-Vorschlag wörtlich) | folgt dem CR-Wortlaut direkt, keine Doppel-Semantik für `token` | zusätzliche Schema-Fläche für dieselbe Information, die `token` bereits trägt — zwei Regex-Felder pro Klasse, die bei denselben Instanz-Präfixen einig gehalten werden müssten (Drift-Risiko), ohne Erkenntnisgewinn gegenüber der Wiederverwendung |
| C — Instanz-Korrelation über ein deklariertes Feld, analog zur Lineage-Ausnahme (`same-id-fields`) | konsistent mit dem einzigen bestehenden Präzedenz-Mechanismus | passt nicht zum Anlassfall: ein Slice deklariert keine „gehört zu"-Zeile, die Korrelation liegt einzig in der Dateibenennung — ein Feld-Mechanismus verlangte eine neue Deklarationspflicht, die kein Repo heute führt |
| **D — Rule-Feld `allow-if-same-id`, Wiederverwendung von `token` für die Quell-Pfad-Extraktion (gewählt)** | kleinste Schema-Erweiterung (ein `bool` je Regel, kein neues Klassen-Feld); dieselbe Wahrheit (das Instanz-Muster) an einer Stelle deklariert; passt zur Pfad-basierten Natur des Anlassfalls | `token` trägt jetzt zwei Rollen (Fund-Erkennung **und** Instanz-Extraktion) — wer die Klasse künftig ändert, muss beide Verwendungen im Kopf behalten; die Fail-closed-Validierung (Capture-Gruppen-Pflicht) ist eine neue, regel-übergreifende Config-Randbedingung |

## Konsequenzen

- **Positiv:** Der CR ist strukturell lösbar, ohne die Klassen-Schema-Fläche
  zu vergrößern — `token` bekommt eine zweite, eng begründete Verwendung
  statt eines Parallel-Feldes.
- **Positiv:** Fail-closed am Config-Rand hält die „keine Freigabe ohne
  eindeutige Korrelation"-Zusage des CR wörtlich ein — eine Regel, die nie
  greifen könnte, lädt nicht still.
- **Negativ, benannt:** `token` ist ab jetzt nicht mehr nur „erkennt Referenzen
  auf diese Klasse im Fremdtext", sondern trägt optional auch „erkennt die
  eigene Instanz-ID am eigenen Pfad" — zwei Lesarten desselben Feldes. Ein
  Reviewer, der `token` künftig ändert, muss beide Verwendungen prüfen.
- **Negativ, benannt:** Die Instanz-Identitäts-Ausnahme deckt nur die
  Token-Form. Ein Slice, der seinen eigenen Review als **Markdown-Link**
  zitiert, bleibt ohne Ausnahme — dieselbe Grenze, die die Spec (Out-of-Scope)
  bereits benennt. Ein Adopter, der überwiegend Links statt Token nutzt,
  bekommt aus diesem ADR keine Erleichterung.
- **Folgepflicht:** keine — `spec/lastenheft.md` und `spec/spezifikation.md`
  sind bereits geschrieben (Teil der CR-Messung, nicht dieser ADR).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Tests | Fünf neue Testfälle nach den [`DC-FA-MTX-003`](../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)-Akzeptanzkriterien (Instanz-Identität Happy/Boundary/Negative/Fehlkonfiguration/Default) | `make test` |
| Config-Adapter-Test | `allow-if-same-id: true` ohne Capture-Gruppe auf einer beteiligten Klasse ⇒ Exit 2 | `make test` |
| `matrix` (Selbstanwendung) | das erweiterte Config-Beispiel in `spec/spezifikation.md` bleibt syntaktisch gültig (`make doc-check`) | `make doc-check` |

## Re-Evaluierungs-Trigger

**Ein zweiter Anlassfall, der dieselbe Instanz-Korrelation für die Link-Form
von `matrix-forbidden` braucht.** Dann ist zu prüfen, ob die Pfad-Extraktion
aus `token` auch dort trägt, oder ob die Link-Form eine andere Korrelations-
Quelle braucht (der Linktext selbst, wie bei der Lineage-Ausnahme).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-18 | Accepted | `slice-229` |
