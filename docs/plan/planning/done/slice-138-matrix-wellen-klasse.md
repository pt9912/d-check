# Slice slice-138: §3.4 nennt fünf Kategorien, `matrix` deckt zwei — eine dritte ist baubar

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos** (Baseline-Regelwerk
[`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht](../../../../.harness/baseline/v5.11.0/regelwerk/modul-06-roadmap.md)):
seine Closure-Bedingung wäre seine eigene DoD.

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §3.4 (Abwärts-Sperre, fünf
Kategorien, zwei gedeckt);
[`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)
und [`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix);
[ADR-0047](../../adr/0047-matrix-spec-historie-nicht-provenance-exempt.md)
(die §7-Historie ist **nicht** provenance-exempt); der Zensus in
[slice-132](../done/slice-132-hard-rule-zensus.md) und seine Berichtigung in
[slice-136](../done/slice-136-agents-34-klaerung.md).

**Berührte Spec-Stellen:** [`spec/lastenheft.md`](../../../../spec/lastenheft.md)
§7 (eine Historie-Zeile wird ent-tokenisiert) — keine Anforderung ändert ihre
Aussage.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

§3.4 verbietet den Spec-Straten **fünf** Referenz-Kategorien; `matrix` trägt
Klassen für **zwei** (ADRs, Slices). Der Zensus hatte das als *gedeckt*
ausgewiesen, [slice-136](../done/slice-136-agents-34-klaerung.md) hat es auf
*zwei von fünf* berichtigt und einen Trigger gesetzt: *die drei fehlenden
Kategorien als Token-Klassen.*

**Dieser Trigger war unbelegt, und die Messung widerlegt ihn zu zwei Dritteln:**

- **Wellen** — baubar. Es gibt 43 Wellendokumente unter
  `docs/plan/planning/**/welle-*.md`; sie bilden eine Dokumentklasse wie
  `adr` und `slice`, und `welle-\d{2}` ist ihr Token.
- **Commit-Hashes** — **nicht** als `matrix`-Klasse ausdrückbar. Ein Hash ist
  kein Dokument; `matrix` verbietet Referenzen **auf eine Klasse von Dateien**,
  und dafür gibt es hier keine.
- **Closure-Daten** — **nicht vom legitimen Bestand trennbar.** Ein Datum ist
  ein Datum: das Lastenheft trägt 98 Vorkommen, 96 davon sind seine **eigenen**
  Historie-Zeilen. Ein Muster, das „Closure-Datum" von „Historie-Datum"
  unterscheidet, existiert nicht.

**Und die Lücke ist nicht theoretisch.** Der Bestand trägt **einen** echten
Verstoß: eine Historie-Zeile des Lastenhefts nennt eine Wellen-Kennung. Sie steht
dort seit 0.60.0 und ist nie aufgefallen, weil die Klasse fehlt, die sie melden
würde.

## 2. Vorgehen

1. **Die `welle`-Klasse in [`.d-check.yml`](../../../../.d-check.yml)** ergänzen —
   Pfade, Token, und die Regel `{from: spec-straten, to: welle, allow: false}`.
   Zu entscheiden und zu begründen: ob auch `{from: adr, to: welle}` fällt.
2. **Am Bestand messen, bevor scharfgeschaltet wird** — die eine bekannte
   Fundstelle ist zu bestätigen, und es ist zu prüfen, ob es weitere gibt.
3. **Die Fundstelle ent-tokenisieren**, nicht ausnehmen. Der Präzedenzfall ist
   [ADR-0047](../../adr/0047-matrix-spec-historie-nicht-provenance-exempt.md)
   Entscheidung 2: die Aussage bleibt, der Abwärts-Verweis geht. Lastenheft-Bump
   und Historie-Zeile nach
   [`MR-032`](../../../../harness/conventions.md#mr-032).
4. **Bewusstes Brechen:** eine injizierte Wellen-Kennung im Spec-Körper ⇒
   `matrix-forbidden` mit gelesener Ursache; Rückbau ⇒ grün.
5. **§3.4s Trigger berichtigen** — er verspricht drei Kategorien und kann eine
   halten. Was die anderen zwei bräuchten, gehört benannt statt als Trigger
   stehen gelassen.
6. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine neue Modul-Fähigkeit.** Ein Muster-Verbot (für Commit-Hashes) wäre ein
  Produkt-Delta mit ADR und Release — hier wird nur benannt, dass es fehlte.
- **Keine Ausnahme für die gefundene Zeile.** Ein `exempt` an dieser Stelle
  machte den neuen Wächter im selben Zug wieder blind.
- **Keine Rückwirkung auf die `## Geschichte` der ADRs.** Sie bleibt
  provenance-exempt ([ADR-0047](../../adr/0047-matrix-spec-historie-nicht-provenance-exempt.md) Entscheidung 3).

## 4. Definition of Done

- [x] Die `welle`-Klasse steht in der Config — **mit beiden Regeln.** Die
      Entscheidung über `adr → welle` ist im ersten Anlauf **falsch** gefallen
      und nach dem Review umgekehrt worden (§9); geführt als
      [`MR-034`](../../../../harness/conventions.md#mr-034).
- [x] Der Bestand ist gemessen: **zwei** Befunde, nicht der eine erwartete. Der
      im Lastenheft ist ent-tokenisiert; der in einer `Accepted`-ADR ist
      **grandfathered** — mit gemessenem Preis (ihr Körper trägt keinen
      weiteren Token, den die Ausnahme mit stummschaltete).
- [x] **Vier** konstruierte Verstöße, jeder mit gelesener Fundstelle:
      `welle-99` und `welle-123` im Spec-Körper, `welle-99` im ADR-Körper, und
      `welle-98` im ausgenommenen ADR — der bleibt **still**, wie er soll.
      Rückbau je 0 Befunde.
- [x] §3.4 verspricht nur noch, was es halten kann: **drei von fünf**, und für
      die zwei übrigen steht der jeweilige Grund da — beim Commit-Hash die
      fehlende **Präzision** (nicht die fehlende Fähigkeit), beim Closure-Datum
      die Ununterscheidbarkeit vom legitimen Bestand.
- [x] Lastenheft auf **0.65.4** mit Historie-Zeile; `make gates` Exit 0 (zehn
      Glieder, 485 Dateien), `make fullbuild` Exit 0; unabhängiger Review
      ([Report](../../../reviews/2026-08-23-slice-138-matrix-wellen-klasse-review.md)),
      blockierend mit einem HIGH, alle drei Befunde eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **Eine neue Token-Klasse trifft mehr als die Spec-Straten.** — **Ausgang:**
  *eingetreten, sofort und mit voller Wucht.* Die erste Messung meldete einen
  zweiten Befund in einer `Accepted`-ADR. Mein erster Umgang damit war der
  falsche — siehe §9. Die Fassung, die jetzt steht, bewacht die Kante **und**
  nimmt den einen immutablen Bestandsfall gemessen aus.
- **Das Token ist zweistellig.** — **Ausgang:** *vor dem Schreiben abgewendet,
  danach belegt.* Die Klasse trägt `welle-\d{2,}`; die Probe mit `welle-123`
  meldet `welle-123` und nicht `welle-12`.
- **Eine Historie-Zeile zu ändern berührt ein Protokoll.** — **Ausgang:**
  *eingetreten, und meine erste Formulierung hat es verdeckt.* Die 0.65.4-Zeile
  behauptete zunächst *„die Aussage bleibt unverändert"* — das stimmte nicht:
  der Zeiger auf die zeitliche Schicht **entfällt ersatzlos**, und die
  Ersatz-Formulierung ließ zunächst sogar ein bezugloses *„jener"* stehen. Beides
  berichtigt; die Zeile sagt jetzt, was wirklich passiert ist.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Bestandsmessung zeigt, dass
die neue Klasse den Bestand breit trifft — dann ist die Räumung ein eigener
Slice und die Scharfschaltung eine Auftraggeber-Frage.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Konfigurations-Profil (GF), Spec-Straten (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23):
  [`BEO-011`](../observations.md) für die Bestandsaussage — „**eine** Fundstelle"
  ist eine Messung und bleibt eine, bis der Wächter läuft.
  [`BEO-002`](../observations.md) für die Ränder der ent-tokenisierten Zeile.
  [`BEO-007`](../observations.md) für jeden Beleg-Lauf.

Slice-ID: slice-138. Betroffene IDs:
[`DC-FA-MTX-001`](../../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix),
[`DC-FA-MTX-003`](../../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix).
Module: `matrix`, Konfigurations-Profil. Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Konfigurations-Erweiterung an bestehender
Modul-Mechanik.

## 9. Closure-Notiz (nach `done/`)

Geliefert: `matrix` trägt die Klasse `welle` und **zwei** Regeln, §3.4 steht auf
drei von fünf gedeckten Kategorien, und der eine Verstoß, den der Bestand seit
Monaten trug, ist weg.

**Der Slice hat eine Regel gebrochen, die dieses Repo selbst geschrieben hat —
und sie stand seit Juni da.** [`MR-006`](../../../../harness/conventions.md#mr-006)
§Scope-Grenze (C-4) nennt die Kante `ADR→Welle` *„bewusst unbewacht"* und
begründet das ausdrücklich damit, dass *„d-check Carveout/Welle/Roadmap **nicht
als `matrix`-Klassen** modelliert; eine Erweiterung wäre ein eigener Change"*.
Dieser Slice macht `welle` zur Klasse — und falsifiziert damit die Bedingung, auf
der die Ausnahme beruhte. Ich habe sie nicht gelesen.

**Stattdessen habe ich die Regel mit einem Non-sequitur wieder herausgenommen.**
Mein Satz lautete: *„§3.4 verlangt sie ohnehin nicht."* Das stimmt — und trägt
nichts, denn die bestehende `adr→slice`-Regel verlangt §3.4 genauso wenig. Sie
kommt aus jenem Eintrag und der kanonischen Referenz-Matrix, die `ADR→Welle` als
**flaches** Verbot führt, ohne den Marker-Ausweg, den `ADR→Slice` offenlässt.
Und die Dichotomie, mit der ich es begründete — *Falsch-Verdikt oder pauschale
Ausnahme* — war eine falsche Alternative: der enge `exempt-paths`-Eintrag, den
die Reihe `0001`–`0021` längst benutzt, stand die ganze Zeit offen.

**Zweimal habe ich in diesem Slice am eigenen Werkzeug vorbeigemessen, und beide
Male aus demselben Grund: eine Abschnittsgrenze ist eine Überschrift, keine
Zeichenkette.**

1. Mein `split('## Geschichte')` traf ein **Zitat im Kopf** der ADR. Der
   „Körper" schrumpfte auf 846 von 4620 Zeichen, und die Messung meldete null,
   wo das Modul einen Treffer fand.
2. Meine erste Gegenprobe hängte den Verstoß **ans Dateiende** — und landete
   damit in `## Geschichte`, dem einzigen Abschnitt, den `matrix` ausnimmt. Der
   Wächter schwieg zu Recht, und ich hätte daraus beinahe geschlossen, er greife
   nicht.

Beide Male hätte ein grünes Ergebnis eine falsche Sicherheit erzeugt. Was half,
war nicht mehr Sorgfalt, sondern das Werkzeug selbst: `matrix` hat den Treffer
gefunden, den mein Skript nicht fand.

**Ein Nebenbefund, der nicht hierher gehört und trotzdem benannt ist.** Ein
`make gates`-Lauf wurde rot, der `make fullbuild` unmittelbar danach auf
demselben Arbeitsbaum grün. Die Ursache ist identifiziert und steht als
[`BEO-014`](../observations.md) im Register: ein Fixture schreibt dieselbe Datei
von `"v1\n"` auf `"v2\n"` — gleiche Größe, gleiche Sekunde, also der
*racily-clean*-Fall der stat-basierten Änderungserkennung. **Bewusst nicht hier
repariert:** ein Fixture-Fix ohne eigene Probe wäre die stille Reparatur, die
diese Arbeit mehrfach beanstandet hat.

**Offen und benannt:** Von den fünf Kategorien des §3.4 bleiben zwei ungedeckt,
und ihre Gründe sind verschieden. Der Commit-Hash **wäre** ausdrückbar — die
Token-Mechanik existiert —, ihm fehlt die **Präzision**; das ist ein Slice mit
einer Messung, kein Verzicht. Das Closure-Datum bleibt Urteil: es ist von den
Daten, die die Spec-Straten in ihren eigenen Historie-Zeilen führen, nicht
unterscheidbar.
