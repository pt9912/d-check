# Slice slice-152: `citations` scharfschalten — und vorher die eigene Doku entschärfen

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-85-baseline-v5120-migration](../welle-85-baseline-v5120-migration.md)
— **Etappe C**, geschnitten von
[slice-150](../in-progress/slice-150-pin-gebundene-zitate.md).

**Bezug:** [`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in);
[ADR-0045](../../adr/0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md)
(das Modul und seine Fail-closed-Entscheidung);
[`MR-038`](../../../../harness/conventions.md#mr-038) (was ein Zitat beim Bump
tut); [`BEO-008`](../observations.md) (vierte Spiegel-Klasse).

**Berührte Spec-Stellen:** `spec/lastenheft.md` / `spec/spezifikation.md` — nur
falls die Lösung ein **Produkt**-Delta ist (Inline-Code-Bewusstsein des
Direktiven-Scans); Bump und Historie dann nach
[`MR-032`](../../../../harness/conventions.md#mr-032).

---

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-26.

## 1. Ziel

Die vierte Spiegel-Klasse hat eine mechanische Form, und sie liegt seit
`v0.50.0` im Produkt: das Modul `citations` prüft ein per
`d-check:cite`-Direktive ausgezeichnetes Zitat **gegen die von ihm zitierte
Quelle**, whitespace-normalisiert. Das ist genau die Prüfung, an der ein
Korpus-Test scheitert.

**Sie ist heute nicht aktivierbar, und das ist gemessen.** Ein Probelauf über
den Bestand bricht fail-closed ab:

```text
d-check: error: CHANGELOG.md:592: malformte d-check:cite-Direktive
```

Die Fundstelle ist die **Dokumentation der Direktive selbst** — der Kopfteil des
Musters steht dort in Inline-Code. Der Scan ist **fence**-bewusst, aber nicht
inline-code-bewusst; eine Fenced-Darstellung wäre immun, die gewählte ist es
nicht. Betroffen sind neben dem CHANGELOG auch `README.md`/`README.de.md`,
Lastenheft, Spezifikation und Benutzerhandbuch.

## 2. Vorgehen

1. **Den Bestand zählen**, bevor entschieden wird: wie viele Stellen schreiben
   die Direktiv-Syntax in Prosa, und wie viele davon in Inline-Code gegen
   Fenced-Block?
2. **Die Wegwahl treffen und begründen** — Doku-Konvention (Syntax nur noch in
   Fenced-Blöcken) oder Produkt-Delta (der Scan überspringt Inline-Code wie die
   übrigen Module). Der zweite Weg ist die geteilte Lexik, der erste eine Regel
   für Menschen; beide haben einen Preis, und er gehört benannt.
3. Erst danach `citations` in den aktiven Modul-Satz — **mit** Bestandsmessung
   vor dem Scharfschalten.
4. Die Zitate aktiver `MR-*`-Einträge auszeichnen, so weit sie den **gepinnten**
   Stand zitieren. Historisch gestempelte Fassungen nach
   [`MR-038`](../../../../harness/conventions.md#mr-038) bleiben **ohne**
   Direktive — ihre Quelle existiert nicht mehr, und ein Wächter darauf wäre
   dauerhaft rot.
5. Bewusstes Brechen: ein verfälschtes Zitat ⇒ `citation-mismatch`, **Ursache
   gelesen**; Rückbau grün.
6. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Auszeichnung eingefrorener Dokumente.** `done/` und Review-Reporte
  zitieren den Stand ihrer Zeit.
- **Keine Direktive auf eine verschwundene Quelle.** Das wäre ein Wächter, der
  konstruktionsbedingt rot bleibt.
- **Keine Ausweitung auf Zitate außerhalb des Konventionsspeichers** in diesem
  Zug — der Bestand dort ist gemessen und klein.

## 4. Definition of Done

- [ ] Der Bestand der Direktiv-Erwähnungen ist **gezählt**, getrennt nach
      Inline-Code und Fenced-Block.
- [ ] Die Wegwahl ist begründet, mit dem benannten Preis des verworfenen Wegs.
- [ ] `citations` läuft über den Bestand — Exit und Befundzahl genannt.
- [ ] Die Zitate aktiver Einträge sind ausgezeichnet, soweit sie den gepinnten
      Stand zitieren; die Ausnahmen sind **benannt**, nicht übergangen.
- [ ] Ein konstruierter Verstoß meldet `citation-mismatch` mit gelesener
      Ursache; Rückbau grün.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Die Direktive trägt Zeilennummern in den vendorten Baum.** Jeder Bump
  verschiebt sie, und dann bricht die Prüfung — laut, aber sie bricht. Ob das
  Wartungs-Last oder gewollter Alarm ist, gehört entschieden statt erlitten. —
  **Ausgang:** *(bei Closure)*
- **Fail-closed heißt: ein Fehler in EINER Direktive nimmt den ganzen Lauf
  mit.** Bei einem Gate im inneren Loop ist das eine andere Zumutung als bei
  einem Closure-Gate; die Bindepunkt-Frage gehört mitentschieden. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Wegwahl ein Produkt-Delta
verlangt, das eine eigene Anforderung braucht.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Module (GF), Konventionsspeicher (GF), Doku (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-26):
  [`BEO-008`](../observations.md) ist der Anlass;
  [`BEO-011`](../observations.md) für jede Aussage darüber, dass der Bestand
  „vollständig" ausgezeichnet sei.

Slice-ID: slice-152. Betroffene IDs:
[`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in).
Module: `citations`. Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Scharfschalten einer vorhandenen,
opt-in-Fähigkeit.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
