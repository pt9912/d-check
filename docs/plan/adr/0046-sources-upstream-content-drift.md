# ADR-0046 — Netz-Modul `sources`: Upstream-Content-Drift externer Quellen (zweite Netz-Tür)

**Status:** Proposed
**Datum:** 2026-07-19
**Autor:** pt9912
**Bezug:** [`DC-FA-SRC-001`](../../../spec/lastenheft.md#dc-fa-src-001--upstream-content-drift-externer-quellen-modul-sources-opt-in-netz) (opt-in Modul `sources`); Modul-/Paket-Fundament [ADR-0005](0005-modul-layout-hexagon-ordner.md), [ADR-0012](0012-kern-paketschnitt-model-rules-app.md); Vorläufer als Harness-Tooling [`MR-022`](../../../harness/conventions.md#mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019) (`fetch-baseline-cache.sh --check-latest`).
**Schärft:** die neue Algorithmus-Sektion [`spec/spezifikation.md` §DC-FA-SRC-001.a](../../../spec/spezifikation.md#dc-fa-src-001a--upstream-content-drift-externer-quellen-sources) und die Erweiterung von [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) um eine zweite Netz-Tür.

## Kontext

d-check pinnt externe Quellen bislang **nicht auf Inhalt**: `external`
([`DC-FA-EXT-001`](../../../spec/lastenheft.md#dc-fa-ext-001--externe-links-modul-external-opt-in))
prüft nur Erreichbarkeit, `pins`
([`DC-FA-PIN-001`](../../../spec/lastenheft.md#dc-fa-pin-001--content-pin-gegen-inhaltlichen-drift-modul-pins-opt-in))
nur **repo-interne** Ziele. Ein Adopter will aber erkennen, ob eine adoptierte
**externe** Form-Quelle inhaltlich gedriftet ist (Kurs-Beispiel
`check_regelwerk_drift.py`, das genau dafür da ist, für die Verzeichnis-Quelle
aber DEFERRED bleibt). d-check deckt seine **eigene** Baseline schon per
Harness-Tooling ab (der Bash-Helfer `fetch-baseline-cache.sh --check-latest`,
[`MR-022`](../../../harness/conventions.md#mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019)),
doch das ist nicht reusable. Wunsch: die Fähigkeit als d-check-**Produkt-Modul**.

Die Spannung: [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
erklärt `external` zur **einzigen** Netz-Tür. Ein zweites Netz-Modul amendiert
diesen QA-Kernvertrag — bewusst, nicht beiläufig.

## Entscheidung

Neues **opt-in Netz-Modul `sources`** (das 19.,
[`DC-FA-SRC-001`](../../../spec/lastenheft.md#dc-fa-src-001--upstream-content-drift-externer-quellen-modul-sources-opt-in-netz)):
pinnt eine externe `http(s)`-Quelle auf einen `sha256`, holt sie über den
Netz-Port (wie `external`, als Post-Pass nach dem Datei-Scan), hasht, vergleicht
→ `source-drift` bzw. `source-unreachable`.

- **Zwei Deklarations-Flächen — Marker *und* Config.** Marker
  `<!-- source-pin: [archive] sha256:<hex> -->` (dpin-Stil, per-Referenz,
  scope-treu, d-check-nativ) **und** ein Config-Block `sources: [{url, sha256,
  unpack}]` (zentral, URL+Hash sauber). Beide, weil beide reale Nutzungsmuster
  bedienen: der in-Doc-Verweis auf eine externe Quelle vs. die zentrale
  Pin-Liste ohne Doku-Anker.
- **Zwei Quelltypen — Einzeldatei *und* Archiv.** Einzeldatei: `sha256` der
  Roh-Bytes. Archiv (`unpack: zip`): `sha256` eines **byte-genau definierten
  Content-Manifests**, nicht der Zip-Roh-Bytes. Begründung: Zip-Framing/
  Recompression ist instabil; der Adopter pinnt **Inhalt** (die enthaltenen
  Dateien), nicht Zip-Bytes. Das Manifest ist **nach dem normalisierten Pfad**
  sortiert (nicht nach der ganzen Hash-Zeile) → **reihenfolge-invariant** und
  vollständig determiniert
  ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)). Es folgt
  **konzeptionell** dem committet-vendored `SHA256SUMS`-Muster
  ([`MR-019`](../../../harness/conventions.md#mr-019--regelwerk-lese-form-committet-statt-gecacht-nachtrag-zu-mr-017)),
  ist aber **eigenständig** kanonisiert (Pfad-Sortierung + -Normalisierung) und
  **nicht** byte-identisch zur unsortierten `sha256sum`-Ausgabe — die byte-genaue
  Form legt die Spezifikation fest.
- **Das Archiv-Keyword ist explizit** (Marker `zip`, parallel zum Config-Wert
  `unpack: zip`) — kein Ableiten aus der `.zip`-Endung; beide
  Deklarations-Flächen nennen denselben Format-Namen (erweiterbar auf
  `tar`/`gz`).
- **`source-drift` emittiert den vollen Ist-`sha256`.** Damit ist Pinnen „einmal
  laufen, gemeldeten Hash kopieren" — die dpin-Ergonomie-Sackgasse (kein Weg an
  den vollen Pin-Hash) tritt für dieses Modul gar nicht erst auf.
- **`source-unreachable` ist von `source-drift` getrennt** — unerreichbar ≠
  gedriftet (ehrliche Diagnose, keine Falsch-Drift bei Netzfehler).
- **Robustheit als Vertrag.** Der Fetch ist **größenbegrenzt** (Body ≤ 64 MiB;
  beim Entpacken Gesamtgröße ≤ 256 MiB, ≤ 10 000 Einträge — Zip-Bomben-Schutz),
  folgt Redirects wie `external` (bis fünf, Inhalt der finalen Antwort) und
  behandelt eine 2xx-Antwort, die **kein gültiges Zip** ist, wie „unerreichbar"
  (`source-unreachable`). Der `sha256` wird case-insensitiv verglichen und stets
  in Kleinbuchstaben emittiert.
- **Amendment der Netz-Sparsamkeit.** Netz jetzt in `external` **und** `sources` — beide
  opt-in, nie im netzlosen Default-Lauf. Der getippte Netzlos-Modullisten-Test
  führt `sources` als zweite Netz-Ausnahme; ein netzloser Lauf mit aktivem
  `sources` öffnet keine Verbindung zum Markdown-Baum, nur zu den explizit
  gepinnten `url`.

## Verglichene Alternativen

- **Nur der Bash-Helfer `--check-latest`.** Deckt nur d-checks **eigene**
  Baseline, nicht reusable. Verworfen — der Nutzer will ein Produkt-Modul.
- **Nur Marker *oder* nur Config.** Je ein reales Nutzungsmuster bliebe
  ungedeckt (zentrale Pin-Liste bzw. in-Doc-Verweis). Verworfen zugunsten
  beider Flächen.
- **Zip-Roh-Byte-Hash statt Content-Manifest.** Bräche bei Recompression und
  Eintrags-Umsortierung; prüfte Framing statt Inhalt. Verworfen.
- **`archive` aus der `.zip`-Endung ableiten.** Fragile Magie. Verworfen
  zugunsten des expliziten Keywords.
- **Currency (neuerer Tag) ins Modul ziehen.** Braucht die Release-API, keinen
  Content-Hash; bleibt im Bash-Helfer
  ([`MR-022`](../../../harness/conventions.md#mr-022--baseline-currency-audit-modus-nachtrag-zu-mr-019)).
  Abgegrenzt.

## Konsequenzen

- **Positiv:** reusable Upstream-Content-Drift für jeden Adopter; **verifiziert**
  die Tag-Immutabilitäts-Annahme (prüfen statt trauen); schließt die
  Hash-Ergonomie für das neue Modul von Anfang an.
- **Kosten:** die zweite Netz-Tür weicht die Netz-Sparsamkeit ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)) auf (bewusst, opt-in, nie
  Default); der Go-Kern gewinnt Netz-Download + Zip-Entpacken (`archive/zip` aus
  der Standardbibliothek — **keine** neue Dependency); mehr Test-/Doku-Fläche
  (Marker + Config × Einzeldatei + Archiv).
- Der Hash-Ergonomie-Fix des **Bestands**-Moduls `pins`/`dpin` bleibt bewusst
  **separat** (eigener Kleinst-Change am Bestandsmodul `pins`) — dieser Slice
  schneidet nicht quer.

## Fitness Function

- Der Netzlos-Test (`forbiddenInNetless` + die
  [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)-Messmethode)
  bleibt grün: `sources` ist nie im netzlosen Default; ein netzloser Lauf öffnet
  keine Verbindung zum Markdown-Baum.
- Boundary: ein Archiv mit **umsortierten** Zip-Einträgen, inhaltlich identisch,
  ergibt **kein** `source-drift` (Manifest reihenfolge-invariant).
- Die `source-drift`-Meldung trägt den **vollständigen** 64-Hex-`sha256`
  (Re-Pin-Vorlage).

## Re-Evaluierungs-Trigger

- Bedarf nach weiteren Archiv-Formaten (`tar`/`gz`) oder nach
  Authentifizierung/Custom-Headern → neue Anforderung.
- Sollte die Netz-Sparsamkeit je „genau **eine** Netz-Tür" erzwingen wollen → diese
  Entscheidung neu bewerten.

## Geschichte

- 2026-07-19: Proposed (doc-first, `slice-080`).
