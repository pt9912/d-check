# Eingehender Change Request — `matrix` kennt nur Klassenpaare, keine Instanz-Identität

**Absender:** Adopter `pg-change-feed` · **Eingegangen:** 2026-09-18, gegen
`ghcr.io/pt9912/d-check:v0.75.0` <!-- d-check:ignore (Lauf-Beleg: der Stand, gegen den der CR geschrieben wurde) -->
(`sha256:18e9cd857f8db3569526d1f9a3cbeba8af51e9f2dd84c17a22444028b977c3da`).
**Richtung:** eingehend — dieses Repo ist der **Empfänger**, nicht der Bittsteller.
**Ziel-Dokument:** [`spec/lastenheft.md`](../../../spec/lastenheft.md)
**Berührt:** [`DC-FA-MTX-001`](../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix),
[`DC-FA-MTX-003`](../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)
**Stand:** **offen** — eingegangen, erste Messung unten, **keine Entscheidung
getroffen**.

**Ablage-Hinweis.** Ein **eingehender** CR ist die dritte Klasse neben
[`MR-035`](../../../harness/conventions.md#mr-035) (ausgehend) und
[`MR-036`](../../../harness/conventions.md#mr-036) (Antworten darauf); der
Kanon führt sie als *externen Vorgang* und ausdrücklich als „bewusst kein
Harness-Konstrukt". Die Datei liegt hier aus demselben Grund wie ihre
Vorgänger: Der erste Konsumenten-CR dieses Repos ging verloren, und mit ihm
die Frage, was genau gebeten und mit welcher Begründung entschieden wurde.

---

## Wortlaut (unverändert übernommen)

> **Betreff:** Das `matrix`-Modul kennt nur Klassenpaare, keine Instanz-Identität
> — verhindert eine mechanische Unterscheidung „eigenes Zitat" vs. „fremdes <!-- d-check:ignore (zitierter Fremd-Wortlaut, kein Verweis) -->
> Zitat" innerhalb derselben Ziel-Klasse.
>
> **Einreicher:** pg-change-feed (Adopter), 2026-09-18, gegen d-check
> `ghcr.io/pt9912/d-check:v0.75.0` <!-- d-check:ignore (zitierter Fremd-Wortlaut, kein Verweis) -->
> (`sha256:18e9cd857f8db3569526d1f9a3cbeba8af51e9f2dd84c17a22444028b977c3da`).
>
> ### Anlass
>
> Zwei bereits `Accepted`-ADRs dieses Repos (ADR-0094, ADR-0097) lehnten <!-- d-check:ignore (fremde ADR-Kennungen des Absenders, kein Verweis in dieses Repo) -->
> unabhängig voneinander eine pauschale Matrix-Regel `{from: slice, to:
> review, allow: false}` ab — mit derselben Begründung: Ein Slice zitiert
> **routinemäßig seinen eigenen** Review-Report; das ist harmlos, weil beide
> gemeinsam in dasselbe Wellen-Archiv wandern. Ein Zitat auf einen **fremden**
> Review (aus einer anderen, ggf. bereits archivierten Welle) trüge dagegen
> real ein Hänger-Risiko für `archive-welle`-artige Werkzeuge. Der Unterschied
> ist **instanzbasiert** (dieselbe Slice-/Welle-ID), nicht klassenbasiert — die
> bestehende `matrix`-Engine kann ihn strukturell nicht ausdrücken.
>
> Real aufgetreten: Ein Versuch, die Unterscheidung doch über eine pauschale
> Klassenregel zu erzwingen, musste rückgängig gemacht werden (eigene
> Supersedes-ADR nötig), nachdem er 11 echte, ausnahmslos harmlose
> Selbst-Zitate (7 Slice-, 5 Welle-Dateien) umschreiben ließ, weil die Regel
> sie nicht von — bislang nicht aufgetretenen, aber möglichen — Fremd-Zitaten
> unterscheiden konnte.
>
> ### Befund
>
> Gelesen: `internal/hexagon/core/rules/matrix.go`. `ruleFor(cfg.Rules,
> srcClass, dstClass)` (Zeile 278) vergleicht ausschließlich **Klassennamen**.
> `classOf()` (Zeile 219) liefert nur den Klassennamen der Quelldatei, nie eine
> aus ihrem Pfad/Namen extrahierte Instanz-ID. Der `token`-Treffer im
> Zieltext (`tokenFindings`, Zeile 97 ff.) wird gegen die Regel-Tabelle
> geprüft, aber nie gegen eine Kennung der Quelldatei korreliert. Es gibt
> keinen Pfad durch den bestehenden Code, der „gleiche ID" ausdrücken könnte.
>
> ### Bitte
>
> Eine optionale Erweiterung der Regel-Form um ein Feld, das eine verbotene
> Kante **freigibt, wenn Quelle und Ziel dieselbe Instanz-ID tragen** —
> skizziert:
>
> ```yaml
> rules:
>   - {from: slice, to: review, allow: false, allow-if-same-id: true}
> ```
>
> Ausführung: Beide beteiligten Klassen (`slice`, `review`) tragen bereits ein <!-- d-check:ignore (zitierter Fremd-Wortlaut, kein Verweis) -->
> `token`-Regex mit einer Capture-Gruppe für die Instanz-ID (z. B.
> `slice-(\d{3})`). Ist `allow-if-same-id: true` gesetzt, wird eine sonst
> verbotene Kante nur dann als `matrix-forbidden` gemeldet, wenn die
> Capture-Gruppe aus dem Dateinamen/Pfad der Quelldatei nicht mit der
> Capture-Gruppe eines Tokens im Zieltext übereinstimmt. Fehlt einer Klasse
> die Capture-Gruppe, bleibt die Regel unverändert streng (fail-closed — <!-- d-check:ignore (zitierter Fremd-Wortlaut, kein Verweis) -->
> keine Freigabe ohne eindeutige ID-Korrelation).
>
> **Gegenprobe für den Test**
>
> | Quelle | Ziel-Zitat | Erwartung |
> |---|---|---|
> | `slice-036-x.md` | referenziert `review-slice-036-x.md` (Token `slice-036`) | grün — gleiche ID, Freigabe greift |
> | `slice-036-x.md` | referenziert `review-slice-099-y.md` (Token `slice-099`) | rot — `matrix-forbidden`, wie heute |
> | `slice-036-x.md` | referenziert einen Review ohne erkennbares `slice-\d{3}`-Token | rot — keine Korrelation möglich, fail-closed |
>
> **Was ausdrücklich nicht gebeten ist**
>
> - Keine Lockerung der bestehenden Klassenregeln ohne das neue Feld —
>   `allow-if-same-id` ist strikt opt-in, Default aus, rückwärtskompatibel.
> - Keine Deaktivierung des Fail-closed-Verhaltens, wenn die
>   ID-Korrelation nicht eindeutig auflösbar ist (fehlende Capture-Gruppe,
>   mehrdeutiges Token) — dann bleibt es beim heutigen `matrix-forbidden`.
> - Keine Pflicht, das Feld für bestehende Regeln zu setzen — Repos ohne
>   diesen Bedarf ändern nichts an ihrer `.d-check.yml`.
>
> **Wie der Adopter bis dahin arbeitet**
>
> Fremd-Zitat bleibt vollständig Review-Prüfpflicht, kein Gate — die
> bestehende Nicht-Regel für slice → review/welle → review (ADR-0094, <!-- d-check:ignore (fremde ADR-Kennung des Absenders, kein Verweis in dieses Repo) -->
> ADR-0097, ADR-0099) gilt unverändert für beide Fälle (eigen wie fremd), <!-- d-check:ignore (fremde ADR-Kennungen des Absenders, kein Verweis in dieses Repo) -->
> bis diese Erweiterung existiert.

---

## Erste Messung am eigenen Werkzeug (2026-09-18, **kein Entscheid**)

Festgehalten, weil sie den Entscheid vorbereitet. **Sie entscheidet nichts** —
der Entscheid ist ein eigener Vorgang und liegt beim Auftraggeber.

**Die drei zitierten Codestellen stimmen exakt**, gegengelesen gegen
`internal/hexagon/core/rules/matrix.go` (Stand: HEAD dieses Repos, 509
Zeilen):

- `classOf()` steht auf **Zeile 219** und liefert ausschließlich `c.Name`
  (den Klassennamen) — keine Instanz-ID-Extraktion aus `rel` (dem Datei-Pfad).
- `ruleFor()` steht auf **Zeile 278** und vergleicht ausschließlich
  `r.From == from && r.To == to` — beides Klassennamen, kein Instanz-Bezug.
- `tokenFindings()` beginnt auf Zeile 91, die Token-Fundstellen-Schleife
  (`for _, loc := range c.Token.FindAllStringIndex(...)`) liegt auf Zeile
  **112** (CR nennt „Zeile 97 ff.", trifft den Funktionskörper, nicht exakt
  die Schleife selbst). Sie nutzt
  `FindAllStringIndex` — nur Start-/End-Offset des Gesamt-Treffers, **keine**
  Capture-Gruppen. Eine Korrelation zur Quelldatei findet nicht statt.

**Die Struktur-Prämisse trägt.** `model.MatrixClass.Token` ist ein einfaches
`*regexp.Regexp` (`internal/hexagon/core/model/config.go` Zeile 165–167);
`model.MatrixRule` trägt nur `From, To string` und `Allow bool` (Zeile
176–179) — kein Feld, das eine Instanz-Kennung tragen könnte. Die Aussage
„es gibt keinen Pfad durch den bestehenden Code, der 'gleiche ID' ausdrücken
könnte" ist damit **bestätigt**, nicht nur plausibel.

**Was der CR voraussetzt und dieses Repo nicht prüfen kann:** die drei
zitierten ADR-Kennungen `ADR-0094`/`ADR-0097`/`ADR-0099` <!-- d-check:ignore (fremde ADR-Kennungen des Absenders, kein Verweis in dieses Repo) -->
sind **nicht** Teil dieses Repos — dessen eigene ADR-Reihe endet bei
[ADR-0086](../adr/0086-pack-alias-fs-os-kapsel-erweiterung.md)
([`docs/plan/adr/README.md`](../adr/README.md)). „Dieses Repos" im
Wortlaut bezeichnet folglich den **Bestand des Absenders**
(`pg-change-feed`), nicht d-check selbst — dieselbe Lesart wie bei früheren
eingehenden CRs, die fremde ADR-Kennungen als Absender-Präzedenz zitieren.
Weder die drei ADRs noch der „real aufgetretene" Rückgängig-Vorgang (11
Selbst-Zitate, 7 Slice-/5 Welle-Dateien) sind von hier aus nachprüfbar und
werden nicht vermutet.

**Der Digest-Verweis auf `v0.75.0` ist von hier nicht gegengeprüft** — das
CHANGELOG dieses Repos führt für `[0.75.0]` kein Digest-Feld an dieser
Stelle; eine Prüfung gegen die veröffentlichte GHCR-Signatur wäre ein
eigener, Netz-gebundener Schritt und ist hier nicht gelaufen.

**Technische Machbarkeit der Bitte, ungeprüft gegen Aufwand/Tragweite:** Die
Erweiterung verlangt zwei neue Fähigkeiten, die heute fehlen: (1) eine
Instanz-ID-Extraktion aus dem **Pfad/Namen der Quelldatei** — heute matcht
`classOf` nur Globs, liest aber keine Capture-Gruppe daraus; (2) den Wechsel
von `FindAllStringIndex` auf eine Capture-Gruppen-fähige Suche
(`FindAllStringSubmatchIndex`) für den **Ziel-Token**-Treffer. Beides ist
lokal in `matrix.go` und dem zugehörigen Config-Adapter umsetzbar, ohne
andere Module zu berühren — das deckt sich mit der Abgrenzung des CR
(„kein neues Modul, keine geänderte Semantik für bestehende Regeln"). Ob das
Feld tatsächlich **so** (opt-in je Regel, `allow-if-same-id`) oder anders
geführt werden sollte, ist eine Design-Frage und Teil des ausstehenden
Entscheids.

**Was diese Messung nicht beantwortet:** ob die Bitte in den Explizit-Scope
von `matrix` (Klassenpaare, keine Instanzen — bisher eine bewusste
Vereinfachung, kein Versehen) passt, oder ob eine Instanz-Korrelation
grundsätzlich woanders hingehört (z. B. `vcs`/`commits`, die bereits
Instanz-Bezüge über Commit-Ranges tragen). Das ist die eigentliche Frage
für den Entscheid.

## Zweite Messung (2026-09-18) — ein vorhandener Präzedenzfall im selben Modul

**`matrix` kennt bereits eine Instanz-Korrelation — für einen anderen
Zweck.** Die Supersede-Lineage-Ausnahme
(`matrix.status.allow-supersede-lineage`/`supersede-fields`,
[`spec/spezifikation.md` §DC-FA-MTX-001.a Schritt 4](../../../spec/spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung),
Implementierung `lineageValues`/`supersedeFieldValues`/`supersedeFieldValue`
in `matrix.go` Zeile 126–166) nimmt eine sonst als `matrix-inactive` gemeldete Kante aus, **wenn**
ein deklariertes Feld der Quelldatei (`**Supersedes:** ADR-0007`) <!-- d-check:ignore (Beispielwert aus der Spec, keine Fundstellen-Referenz) --> den
Linktext oder den aufgelösten Zielpfad der Referenz als Teilzeichenkette
enthält — dieselbe Grundidee wie die Bitte: eine sonst verbotene/inaktive
Kante wird **instanzbasiert** freigegeben, nicht klassenbasiert.

**Zwei Unterschiede zur Bitte, beide relevant für den Entscheid:**

1. **Geltung.** Die Lineage-Ausnahme wirkt nur auf `matrix-inactive`
   (Status-Prüfung); die Klassen-Regelprüfung (`matrix-forbidden`) „bleibt
   unberührt" — genau die Prüfung, um die die Bitte geht, ist bei der
   vorhandenen Ausnahme **ausdrücklich ausgenommen**. Es gibt also *keinen*
   bestehenden Weg, der die Bitte bereits erfüllt — nur ein
   Konstruktionsprinzip, das sich wiederverwenden ließe.
2. **Quelle der Korrelation.** Die Lineage-Ausnahme liest ein **deklariertes
   Feld** im Fließtext der Quelldatei (Kopplung über Inhalt); die Bitte
   verlangt eine **aus dem Dateinamen/Pfad extrahierte** ID (Kopplung über
   Namenskonvention). Ein Slice deklariert keine „**MeinReview:**
   review-slice-036"-Zeile — die Korrelation existiert nur in der
   Namenskonvention (`slice-036-x.md` ↔ `review-slice-036-x.md`). Die
   Pfad-Extraktion der Bitte ist damit keine Vereinfachung des
   Lineage-Mechanismus, sondern eine andere Informationsquelle für dieselbe
   Grundfigur.

**Eine Wiederverwendungs-Möglichkeit, die die Bitte selbst nicht benennt:**
`MatrixClass.Token` ist bereits ein `*regexp.Regexp` mit (laut Bitte)
potenzieller Capture-Gruppe. Er wird heute ausschließlich auf **fremden
Prosa-Text** angewendet (`tokenFindings`, um Bare-Token-Zitate zu finden).
Dieselbe Regex ließe sich für die **Quell-ID-Extraktion** auf den
**Pfad der Quelldatei selbst** anwenden — kein neues Schema-Feld für die
Quellseite nötig, nur eine zweite Anwendung des bereits vorhandenen
`token`-Feldes auf einen anderen Input. Für die Zielseite genügt der
Wechsel von `FindAllStringIndex` auf `FindAllStringSubmatchIndex`. Das
wäre eine schmalere Erweiterung als das in der Bitte skizzierte Schema
(kein zusätzliches `id-pattern`-Feld je Klasse) — bei identischer
Außenwirkung für den beschriebenen Fall.

## Empfehlung (kein Entscheid — zur Bestätigung)

**Einordnung:** Die Bitte passt in den bestehenden Scope von `matrix` — sie
ist keine Erweiterung auf „Instanzen statt Klassen", sondern eine **zweite
Instanz-Ausnahme neben einer bereits existierenden** (Lineage), nur für
`matrix-forbidden` statt `matrix-inactive` und mit Pfad- statt
Feld-Korrelation. Der in der ersten Messung erwogene Alternativ-Standort
(`vcs`/`commits`) trägt nicht: jene Module korrelieren Commit-Ranges, nicht
Dokumentklassen-Paare — die Bitte bleibt eine `matrix`-Frage.

**Vorschlag, falls angenommen:** ein neues Rule-Feld
`allow-if-same-id: bool` (Default `false`, byte-identisch ohne Nutzung —
dieselbe Zusage wie bei `allow-supersede-lineage`), das bei `true` **beide**
beteiligten Klassen zwingt, ein `token`-Regex mit **genau einer**
Capture-Gruppe zu tragen (sonst Exit 2 bei Config-Laden — fail-closed vor
dem Lauf, nicht erst beim ersten Fund); die Quell-ID wird durch Anwendung
von `srcClass.Token` auf den Quell-Dateipfad gewonnen, die Ziel-ID durch
`FindAllStringSubmatchIndex` auf den Prosa-Treffer; eine Übereinstimmung
(nach Trim, case-sensitiv wie die übrigen `matrix`-Vergleiche) nimmt den
Fund aus `matrix-forbidden` aus, jeder Nicht-Treffer bleibt gemeldet.

**Was das für die Umsetzung heißt, sollte diese Empfehlung bestätigt
werden:** eine neue ADR (`Schärft:` →
[`DC-FA-MTX-001`](../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)/[`DC-FA-MTX-003`](../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix),
mit `Supersedes` keiner bestehenden ADR — echte Erweiterung, keine Korrektur),
danach ein Slice mit dem in §Bitte skizzierten Gegenprobe-Test als
Akzeptanzkriterien-Trio. **Nicht ohne Bestätigung ausgelöst** — die ADR ist
nach `Accepted` immutable (`AGENTS.md` §3.5), und diese Empfehlung ist ein
Vorschlag, kein Entscheid.
