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

- [ ] Die `welle`-Klasse steht in der Config, mit begründeter Entscheidung über
      die zweite mögliche Regel (`adr → welle`).
- [ ] Der Bestand ist **gemessen**, nicht erwartet; jede Fundstelle ist
      ent-tokenisiert, keine ausgenommen.
- [ ] Der konstruierte Verstoß meldet `matrix-forbidden` — **Ursache gelesen**,
      nicht nur der Exit; Rückbau grün.
- [ ] §3.4 verspricht nur noch, was es halten kann; für die zwei übrigen
      Kategorien steht dort, **was** ihnen fehlt.
- [ ] Lastenheft-Bump + Historie-Zeile; `make gates` grün (Exit explizit);
      unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Eine neue Token-Klasse trifft mehr als die Spec-Straten.** `matrix` prüft
  Klassen gegeneinander; kommt `welle` dazu, ist zu prüfen, welche **bestehenden**
  Regeln sie berührt — und ob dabei Befunde auf ADRs oder Slices entstehen, die
  niemand erwartet hat. — **Ausgang:** *(bei Closure)*
- **Das Token `welle-\d{2}` ist zweistellig.** Bei `welle-100` griffe es nicht
  mehr — oder schlimmer: es griffe auf die ersten zwei Stellen. Die Grenze
  gehört geprüft, nicht angenommen. — **Ausgang:** *(bei Closure)*
- **Eine Historie-Zeile zu ändern berührt ein Protokoll.** Der Kanon ändert
  Historie-Zeilen nicht rückwirkend; hier geht es um die **Form** eines
  Verweises, nicht um die Aussage. Die Unterscheidung ist zu belegen, nicht zu
  behaupten. — **Ausgang:** *(bei Closure)*

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

*(wird mit dem Closure-Body gefüllt)*
