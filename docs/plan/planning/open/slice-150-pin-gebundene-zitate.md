# Slice slice-150: Ein Zitat der Baseline ist pin-gebunden wie ein Link

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-85-baseline-v5120-migration](../welle-85-baseline-v5120-migration.md)
— **Etappe C**, geschnitten vom Delta-Audit in
[slice-149](../done/slice-149-baseline-v5120-delta-audit.md).

**Bezug:** [`MR-021`](../../../../harness/conventions.md#mr-021) (in-Repo-Verweise
sind pin-gebunden — der Eintrag, den dieser Slice schärfen würde);
[`MR-033`](../../../../harness/conventions.md#mr-033) (der eine Bestandsfall);
[`BEO-008`](../observations.md) (vierte Spiegel-Klasse, seit
[slice-148](../done/slice-148-baseline-v5120-vendoring.md) im Register).

**Berührte Spec-Stellen:** — (Harness-Regeltext; ein Produkt-Delta nur, falls
die Messung eine mechanische Form trägt).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-26.

---

## 1. Ziel

[`MR-021`](../../../../harness/conventions.md#mr-021) bindet **Links** auf den
vendorten Baum an den Pin: beim Bump wandern sie mit. Ein **Zitat** des
vendorten Textes wandert nicht — der Pfad daneben wird gehoben, löst sauber auf,
und der zitierte Wortlaut existiert am neuen Ziel nicht mehr. Beide Hälften sind
für sich in Ordnung, und kein Gate sieht die Kombination.

Bestandsfall: [`MR-033`](../../../../harness/conventions.md#mr-033) zitiert
zweimal die Fassung, die der eigene Konsumenten-CR hat ändern lassen.

**Die Vorfrage ist keine Bauentscheidung, sondern eine Regelfrage:** Ein
`Accepted`-Eintrag wird nach dem Kanon **nie überschrieben**. Ist die Korrektur
eines veralteten Zitats ein Überschreiben — oder die Pflege einer
pin-gebundenen Referenz, die [`MR-021`](../../../../harness/conventions.md#mr-021)
ohnehin verlangt? Davon hängt ab, ob der Bestandsfall repariert werden **darf**.

## 2. Vorgehen

1. **Die Regelfrage zuerst und begründet beantworten** — mit dem Kanon, nicht
   mit der Bequemlichkeit. Der Geltungsbereich von
   [`MR-021`](../../../../harness/conventions.md#mr-021) nennt heute
   ausdrücklich **Links**; ihn auf Zitate zu lesen, wäre genau die
   Reichweiten-Dehnung, die [`BEO-012`](../observations.md) führt.
2. Fällt sie für die Pflege: ein **neuer** Eintrag, der
   [`MR-021`](../../../../harness/conventions.md#mr-021) **schärft** (Titel
   trägt das nach Kanon-Form) — nicht dessen Änderung.
3. Erst dann den Bestandsfall reparieren.
4. **Messen, ob eine mechanische Form trägt:** die Prüfung ist *Zitat gegen
   sein eigenes Link-Ziel*, nicht gegen einen Korpus — der Korpus-Test ist in
   [slice-148](../done/slice-148-baseline-v5120-vendoring.md) nachweislich
   gescheitert. Ob das Produkt das heute hergibt (`citations` ist ein Modul
   dieses Repos), ist zu **messen**.
5. Die Bump-Prozedur um den Schritt ergänzen, falls er trägt.
6. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Änderung an [`MR-021`](../../../../harness/conventions.md#mr-021)
  selbst.** Einträge werden nicht überschrieben.
- **Keine Reparatur eingefrorener Zitate.** `done/` und Review-Reporte zitieren
  den Stand ihrer Zeit; der Bestandsfall ist der **lebende** Eintrag.
- **Kein Gate ohne Messung.** Ein Wächter auf Zitate hat eine
  Falsch-Positiv-Last, und die gehört gezählt, bevor er gebaut wird.

## 4. Definition of Done

- [ ] Die Regelfrage ist **begründet** beantwortet, mit Kanon-Fundstelle.
- [ ] Bei Pflege-Antwort: neuer Eintrag in Kanon-Form; der Bestandsfall
      repariert und die Reparatur belegt (Zitat gegen das neue Ziel).
- [ ] Bei Überschreib-Antwort: der Bestandsfall ist **ausgewiesen** statt
      repariert, und die Ablage-Form der Korrektur ist benannt.
- [ ] Die mechanische Form ist **gemessen** beantwortet — trägt sie, oder ist
      die Falsch-Positiv-Last zu hoch?
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Die bequeme Antwort ist die Pflege-Antwort.** Sie erlaubt die Reparatur und
  spart einen Eintrag. Genau darum ist sie die verdächtige — die Begründung muss
  aus dem Kanon kommen, nicht aus dem Ergebnis. — **Ausgang:** *(bei Closure)*
- **Ein Zitat-Wächter misst Prosa gegen Prosa.** Umbrüche, Auszeichnung und
  Auslassungszeichen machen jede naive Gleichheit falsch; in
  [slice-148](../done/slice-148-baseline-v5120-vendoring.md) hat dieselbe
  Prüfung mehrere verworfene Fassungen gebraucht, bis sie trug — die dort
  genannte Zahl **sechs** zählt Dokument/Quell-Paare, nicht Anläufe. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Regelfrage eine Klärung des
Kanons braucht — dann ist sie ein CR-Kandidat, kein Slice.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Konventionsspeicher (GF), Produkt-Module (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-26):
  [`BEO-008`](../observations.md) ist der Anlass;
  [`BEO-012`](../observations.md) für die Frage, wie weit
  [`MR-021`](../../../../harness/conventions.md#mr-021) trägt.

Slice-ID: slice-150. Betroffene IDs: — (kein `DC-`-Bezug, solange kein
Produkt-Delta entsteht). Module: Konventionsspeicher. Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Regelfrage am eigenen Konventionsspeicher.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
