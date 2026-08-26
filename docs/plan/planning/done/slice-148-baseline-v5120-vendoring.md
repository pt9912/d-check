# Slice slice-148: Baseline-Pin auf v5.12.0 — Etappe A der Migration

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-85-baseline-v5120-migration](../welle-85-baseline-v5120-migration.md).

**Bezug:** [`MR-011`](../../../../harness/conventions.md#mr-011) (Pin auf
Release-Tag, diese Hebung ist ihre Fortschreibung),
[`MR-021`](../../../../harness/conventions.md#mr-021) (in-Repo-Verweise sind
pin-gebunden), [`MR-023`](../../../../harness/conventions.md#mr-023)
(self-contained Bundle-Layout), [`MR-030`](../../../../harness/conventions.md#mr-030)
(der abzulösende Vorgänger-Pin), [`BEO-008`](../observations.md) (drei
Spiegel-Klassen).

**Berührte Spec-Stellen:** — (Harness-Bestand; das Produkt bleibt unberührt).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-26.

---

## 1. Ziel

Der vendorte Baum wandert von `v5.11.0` auf `v5.12.0`, und **alle** Verweise
wandern mit. Das Delta ist bereits gemessen (Wellendokument §1): 28 von 52
Dateien unterscheiden sich, davon **fünf** mit echtem Regel-Inhalt.

**Der Inhalt ist nicht Gegenstand dieses Slice** — nur der Stand. Was die fünf
Änderungen für dieses Repo bedeuten, beantwortet
[slice-149](../open/slice-149-baseline-v5120-delta-audit.md).

## 2. Vorgehen

1. Bundle vendoren und **verifizieren** (`--verify`, Manifest über beide Bäume).
2. **Je Spiegel-Klasse aus [`BEO-008`](../observations.md) eine Zählung vor und
   nach der Hebung:** Pfad-Verweise (gate-gedeckt), Release-/Tree-URLs,
   Prosa-/Ellipsen-Pins. Die Zahl kommt in die Closure-Notiz, nicht nur der
   grüne Lauf.
3. **Die vierte Klasse prüfen, die `BEO-008` nicht führt:** ein Verweis, der
   nicht nur auf eine Datei zeigt, sondern deren **Wortlaut zitiert**. Bei
   Punkt 2 des CR hat sich genau dieser Wortlaut geändert
   ([`MR-033`](../../../../harness/conventions.md#mr-033) zitiert die alte
   Fassung). Ob das eine eigene Klasse ist, gehört gemessen und benannt.
4. Ein neuer Adaptions-Eintrag trägt die Pin-Hebung; [`MR-030`](../../../../harness/conventions.md#mr-030)
   nach `conventions/done/` samt Link-Tiefen-Fix im Move-Commit
   ([`MR-013`](../../../../harness/conventions.md#mr-013)).
5. Alt-Baum `v5.11.0` entfernen; eingefrorene Zitate in `done/` bleiben, ihr
   Pfad geht bei Bedarf ins Quell-skopierte Ventil.
6. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine inhaltliche Folge.** Was die fünf geänderten Dateien für dieses Repo
  heißen, ist Etappe B.
- **Kein Freshness-Audit der Adaptionen.** Ebenfalls Etappe B.
- **Kein Retrofit eingefrorener Dokumente.** `done/` und Review-Reporte
  zitieren den Stand ihrer Zeit.

## 4. Definition of Done

- [x] Bundle vendored und **offline verifiziert**: 51 Dateien, Manifest über
      beide Bäume vollständig (`fetch-baseline-cache.sh --verify` Exit 0).
- [x] Je Spiegel-Klasse eine Zahl: Pfad-Verweise **118** gesamt / **42 lebend
      gehoben** / 76 eingefroren · Release-/Tree-URLs **5 lebend**, alle
      gehoben · nackte Nennungen **24** / **4 gehoben** / 20 als
      Vergangenheits-Aussagen stehengelassen. **Die Messvorschrift steht neben
      der Zahl** — ohne sie war die dritte nicht nachrechenbar (§9).
- [x] Die Zitat-Klasse ist gemessen und beantwortet: **ja, eine eigene Klasse**,
      und `BEO-008` führt sie nicht. Drei Dokumente, sechs Paare, davon **genau
      eines lebend** ([`MR-033`](../../../../harness/conventions.md#mr-033)).
- [x] [`MR-037`](../../../../harness/conventions.md#mr-037) geschrieben,
      [`MR-030`](../../../../harness/conventions.md#mr-030) aufgelöst und nach
      `conventions/done/` verschoben; beide Index-Tabellen nachgeführt.
- [x] Alt-Baum entfernt; kein hängender Verweis — die eingefrorenen Zeiger
      tragen ein quell-skopiertes Ventil, dessen Preis benannt ist.
- [x] `make gates` Exit 0 (zehn Glieder), `make fullbuild` Exit 0; unabhängiger
      Review ([Report](../../../reviews/2026-08-26-slice-148-baseline-v5120-review.md)),
      blockierend mit **zwei MEDIUM** und fünf LOW, alle sieben eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **Ein Pfad-`grep` kennt keine Zeitform.** [`BEO-008`](../observations.md)
  nennt die Über-Hebung ausdrücklich: eine Aussage über die **Vergangenheit**
  darf nicht mitgehoben werden. Jede Fundstelle ist auf Gegenwart oder
  Historie zu prüfen. — **Ausgang:** *entfallen — nicht durch Glück, sondern
  weil kein nacktes Versions-Ersetzen lief.* Gehoben wurde ausschließlich über
  **Muster** (Pfad, Release-URL, Tree-URL); die 24 nackten Nennungen bekamen je
  ein Urteil, 4 gehoben, 20 stehengelassen. Der Review hat **beide** Richtungen
  unabhängig geprüft: nichts Lebendes übersehen, nichts über-gehoben.
- **Ein gehobener Link kann auf einen geänderten Wortlaut zeigen.** Der Pfad
  löst auf, das Zitat daneben stimmt nicht mehr — und kein Gate sieht das. —
  **Ausgang:** *eingetreten, genau einmal — aufgefangen von
  [slice-149](../open/slice-149-baseline-v5120-delta-audit.md).*
  [`MR-033`](../../../../harness/conventions.md#mr-033) zitiert zweimal die
  Fassung, die CR-Punkt 2 hat ändern lassen; der Pfad wurde mitgehoben und löst
  sauber auf. Das Schicksal des Eintrags — bleibt, abgelöst oder aufgelöst —
  steht ausdrücklich in der DoD von slice-149, nicht hier.

## 6. Trigger

**Start** (`open` → `in-progress`): Welle eröffnet.

**Rückführungen:** `in-progress` → `next`, falls die Verifikation des Bundles
scheitert.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Bestand (GF), Konventionsspeicher (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-26):
  [`BEO-008`](../observations.md) ist der Anlass;
  [`BEO-011`](../observations.md) für jede Aussage darüber, dass eine
  Spiegel-Klasse „vollständig" gehoben sei.

Slice-ID: slice-148. Betroffene IDs: — (kein `DC-`-Bezug). Module:
Harness-Bestand. Gates: `make doc-check`, `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Fortschreibung eines etablierten Vorgangs.

## 9. Closure-Notiz (nach `done/`)

Geliefert: Der Pin steht auf `v5.12.0`, 51 Dateien offline verifiziert, alle
lebenden Verweise gehoben, [`MR-030`](../../../../harness/conventions.md#mr-030)
aufgelöst, der Alt-Baum entfernt.

**Eine vierte Spiegel-Klasse, die [`BEO-008`](../observations.md) nicht führt.**
Ein Verweis kann nicht nur **zeigen**, sondern **zitieren**. Der Pfad wird
mitgehoben, löst sauber auf — und der Wortlaut daneben existiert am neuen Ziel
nicht mehr. Kein Gate sieht das, weil beide Hälften für sich in Ordnung sind.
Gemessen: drei Dokumente, sechs Paare; ein eingefrorener Report, ein durch seine
eigene Kopfzeile historischer CR, und **ein lebender Eintrag**.

**Zwei Zahlen dieses Slice waren falsch gerahmt, nicht falsch gemessen** — und
beide Male fehlte derselbe Satz. Der Zensus der vierten Klasse sagte „drei
Dokumente" nicht, sondern „zwei", weil meine Messung die eingefrorenen
Verzeichnisse ausschloss und der **Skopus** nirgends dabeistand. Die dritte
Spiegel-Klasse nannte „24", ohne die **Messvorschrift** — der Reviewer kam
mit anderer Abgrenzung auf 60, 22 und 14. Eine Zahl ohne ihre Vorschrift ist
keine Messung, sondern eine Behauptung mit Ziffern. Das ist
[`BEO-009`](../observations.md) Richtung (b), an einem Tag zum zweiten Mal.

**Die erste Zählung der ersten Klasse war zu eng** und hätte den Reviewer-Skill
auf einen gelöschten Baum zeigen lassen: das Muster verlangte
`.harness/baseline/…`, der Skill verweist relativ mit `../baseline/…`. Fünf
Vorkommen, ein Dokument. Gefunden hat es nicht das Muster, sondern die Liste
der nackten Nennungen daneben — ein zweiter Blickwinkel auf dieselbe Menge.

**Beim Nachrechnen kam ein eigener Fund dazu:** das Bundle führt zwei
Nicht-Markdown-Vorlagen, die meine `*.md`-Delta-Schleife nie gesehen hat. Beide
sind unverändert — geprüft, nicht angenommen. Dieselbe Klasse wie oben: ein
Filter, der still weniger sieht als die Menge, die er zu messen behauptet.

**Zwei Pflichtfelder fehlten in einem Eintrag, dessen Vorlage ich im selben
Commit vendored habe.** Die MR-Vorlage verlangt bei einer Ablösung `Löst auf`
und schreibt im selben Satz, warum der Verweis auf die **Index-Zeile** geht:
*„die wandert bei Auflösung nach `done/` und ein Pfad-Link bricht genau dann."*
Mir fehlte das Feld, mein Zeiger war ein Pfad-Link, und `Begründung` fehlte
ebenfalls. Der Commit hat dieselbe Reparatur zwölfmal an
[`MR-030`](../../../../harness/conventions.md#mr-030) vorgenommen.

**Das Ventil ist ein Glob statt vierzehn Datei-Einträgen**, und der Preis steht
im Kommentar: ein Report, der morgen entsteht und auf den entfernten Baum
zeigt, bleibt stumm. Bewusst in Kauf genommen — ein Report zitiert den Stand
seiner Zeit —, mit der Grenze daneben: zum **gepinnten** Baum melden tote Links
weiterhin, auch aus `done/`.

**Register:** [`BEO-008`](../observations.md) auf Zähler **4**.
