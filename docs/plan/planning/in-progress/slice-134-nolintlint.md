# Slice slice-134: Die verbotene Direktive wirkt — `nolintlint` ins Profil

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-84-durchsetzung](../welle-84-durchsetzung.md).

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §3.2 (Suppression-Verbot;
Auflösungs-Trigger `nolintlint`), §3.6 (Linter-Strenge nur per ADR);
[ADR-0006](../../adr/0006-lint-profil-solid.md) (Lint-Profil);
geschnitten vom Zensus in [slice-132](../done/slice-132-hard-rule-zensus.md).

**Berührte Spec-Stellen:** — (Lint-Profil; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

Der Zensus hat §3.2 als **einzige** Regel gemessen, deren verbotene Form nicht
nur ungeprüft bleibt, sondern **wirkt**: ein echter Lint-Verstoß mit passendem
`//nolint` lässt `make lint` mit **Exit 0** durchlaufen. `nolintlint` — der
Linter, der genau das meldet — steht nicht im Profil.

Das ist die schwerste Klasse, die der Zensus gefunden hat: nicht eine Regel ohne
Wächter, sondern eine Regel, deren Umgehung **funktioniert und still bleibt**.

## 2. Vorgehen

1. `nolintlint` in [`.golangci.yml`](../../../../.golangci.yml) aufnehmen,
   Einstellungen bewusst wählen: verlangt wird mindestens eine **Begründung**
   und ein **spezifischer** Linter-Name; ob eine unbenutzte Direktive gemeldet
   wird, ist eigens zu entscheiden.
2. **Am Bestand messen, bevor scharfgeschaltet wird:** wie viele Befunde meldet
   das Profil heute? Bei 0 ist der Zug frei; bei >0 gehört jede Fundstelle
   einzeln beurteilt — geräumt oder als zentrale Ausnahme mit `Why:` geführt.
3. **Bewusstes Brechen — drei Formen, weil zwei Schalter getrennt greifen.**
   Der Zensus-Verstoß trug keine Begründung; erst die **vollständige** Form
   (benannter Linter **und** nachgestellter `// Grund`) passiert alle drei
   Schalter. Zu messen sind deshalb alle drei: nackte `//nolint` ⇒ rot ·
   benannt, aber unbegründet ⇒ rot · vollständig ⇒ **grün**. Der letzte, grüne
   Exit ist der Beleg für die verbleibende Lücke — und die Zusage des Gates ist
   genau die Differenz zwischen den ersten beiden und ihm.
4. `AGENTS.md` §3.2 nachziehen: der Auflösungs-Trigger ist eingelöst, die Regel
   wechselt von *einseitig* auf **teilgedeckt** — nicht auf *gedeckt*, denn
   `nolintlint` prüft die **Form** der Direktive, nicht ihre Berechtigung (§5).
   Der Zensus-Eintrag in
   [slice-132](../done/slice-132-hard-rule-zensus.md) bleibt als
   historische Messung stehen.
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Lockerung anderer Linter**, um Platz zu schaffen. Eine Senkung wäre
  nach §3.6 ADR-pflichtig.
- **Keine Sanierung von Alt-Suppressions über die Messung hinaus.** Findet
  Schritt 2 Bestand, wird er benannt und einzeln entschieden — nicht pauschal
  ausgenommen.
- **Kein zweiter Wächter** auf andere Suppression-Formen (`//lint:ignore`,
  Build-Tags). Erst messen, ob es sie gibt.

## 4. Definition of Done

- [ ] `nolintlint` ist im Profil, mit begründeter Einstellung je Schlüssel.
- [ ] Der Bestand ist **gemessen** und die Zahl steht in der Closure-Notiz;
      jede Fundstelle ist geräumt oder zentral mit `Why:` geführt.
- [ ] **Drei** Direktiv-Formen gemessen, jede mit gelesenem Exit **und**
      geprüfter Ursache (nicht „rot" allein): nackt ⇒ rot durch `nolintlint` ·
      benannt-unbegründet ⇒ rot durch `nolintlint` · vollständig ⇒ grün.
- [ ] `AGENTS.md` §3.2 sagt hin, was das Gate trägt **und was nicht**.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein Linter, der Bestand meldet, wird gern per Ausnahme stumm gestellt** —
  und dann steht die Ausnahme da, wo vorher die Direktive stand. Jede zentrale
  Ausnahme braucht ihr `Why:` und einen Namen, keinen Sammelbegriff. —
  **Ausgang:** *(bei Closure)*
- **`nolintlint` prüft die Form der Direktive, nicht ihre Berechtigung.** Eine
  begründete, spezifische `//nolint` besteht ihn — die Regel verbietet sie
  trotzdem. Die Deckung wird also **teilweise** sein, und das gehört benannt
  statt aufgerundet. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-132](../done/slice-132-hard-rule-zensus.md)
in `done/`, WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Bestandsmessung so viele
Fundstellen zeigt, dass ihre Räumung ein eigener Slice ist.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Lint-Profil (GF), Gate-Landschaft (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23):
  [`BEO-011`](../observations.md) für die Bestandsaussage in Schritt 2 — „null
  Fundstellen" ist eine Messung, keine Erwartung. [`BEO-007`](../observations.md)
  für jeden Beleg-Lauf.

Slice-ID: slice-134. Betroffene IDs: — (Lint-Profil; keine Anforderung).
Module: Lint-Profil, Gate-Landschaft. Gates: `make lint`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Profil-Erweiterung an bestehender Mechanik.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
