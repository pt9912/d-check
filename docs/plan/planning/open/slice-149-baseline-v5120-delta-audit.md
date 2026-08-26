# Slice slice-149: Delta-Audit v5.12.0 und Freshness-Audit aller Adaptionen — Etappe B

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-85-baseline-v5120-migration](../welle-85-baseline-v5120-migration.md).

**Bezug:** [`modul-02-harness-bootstrap.md` §Freshness-Audit](../../../../.harness/baseline/v5.12.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2);
[Antwort des Kurses](../../cr/2026-08-26-antwort-regelwerk-v5110.md);
[`MR-033`](../../../../harness/conventions.md#mr-033) (die eine Adaption, deren
Grundlage sich nachweislich geändert hat);
[slice-143](../done/slice-143-structure-abschnitts-skopus.md) (unsere Deckung
der urteilsfreien Hälfte).

**Berührte Spec-Stellen:** — vermutlich keine; ergibt das Audit ein
Produkt-Delta, wächst es in einem **geschnittenen** Slice, nicht hier.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-26.

---

## 1. Ziel

Zwei Audits, die sich eine Frage teilen — *was heißt der neue Stand für uns?*

**Delta-Audit:** je Kurs-Welle (95–98, eine je CR-Punkt) eine Antwort:
Handlung mit Slice-Kennung oder **belegt** folgenlos. Zwei Erwartungen stehen
schon (Wellendokument §1) — sie sind Hypothesen, keine Belege.

**Freshness-Audit:** je **aktiver** Adaption eine Antwort auf die Frage des
Kanons — *Regelt die neue Fassung das, wofür dieser Eintrag angelegt wurde?*
Keiner bleibt ungefragt, auch nicht die offensichtlich unberührten.

## 2. Vorgehen

1. **Delta zuerst, Adaptionen danach** — eine Adaption gegen einen
   ungelesenen Stand zu prüfen, ist Raten.
2. Je Kurs-Welle: den **vendorten** Text lesen, nicht die Antwort des Kurses.
   Sie ist Erwartung; der Beleg ist die Datei.
3. **[`MR-033`](../../../../harness/conventions.md#mr-033) ist der benannte
   Prüffall.** Der Kanon sagt jetzt *„Die Erlaubnis ist keine Pflicht"* — die
   Frage ist, ob unser **Verbot** danach noch eine Verschärfung ist oder eine
   nicht ausgeübte Erlaubnis. Beide Antworten sind zulässig; die Begründung
   entscheidet, nicht die Bequemlichkeit.
4. **Die urteilsfreie Hälfte gegenrechnen:** Der neue `modul-05`-Text nennt
   *dass ein Ausgang dasteht* **und** *welcher der drei*. Unsere Regel prüft den
   häufigsten Auslöser. Die Differenz gehört benannt und die Entscheidung
   getroffen — als **eigener** Slice, falls sie fällt.
5. Alles Geschnittene bekommt eine Kennung, eine Zeile in §4 des
   Wellendokuments und eine Drift-Log-Zeile.
6. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Umsetzung der geschnittenen Arbeit.** Das Audit schneidet, es baut
  nicht.
- **Keine Antwort ohne gelesene Stelle.** „Vermutlich folgenlos" ist kein
  Ergebnis; folgenlos **mit Fundstelle** ist eines.
- **Keine Änderung an eingefrorenen Dokumenten.**

## 4. Definition of Done

- [ ] Vier Kurs-Wellen, vier Antworten — je Handlung eine Kennung, je
      Folgenlosigkeit eine Fundstelle.
- [ ] **Jede** aktive Adaption ist gefragt; das Ergebnis steht je Eintrag, auch
      „unberührt".
- [ ] [`MR-033`](../../../../harness/conventions.md#mr-033) ist begründet
      entschieden — bleibt, wird abgelöst oder aufgelöst.
- [ ] Die Differenz zur vollen urteilsfreien Hälfte ist benannt und entschieden.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein Audit, das nur bestätigt, was die Antwort sagt, ist kein Audit.** Die
  Kurs-Antwort ist eine gute Erwartung und ein schlechter Beleg; sie kann in
  Details anders gelandet sein, als sie ankündigt. — **Ausgang:** *(bei Closure)*
- **Der Freshness-Audit ist lang und langweilig, und genau darum wird er
  abgekürzt.** Wer ihn auf die „offensichtlich betroffenen" Einträge verengt,
  hat die Klasse [`BEO-011`](../observations.md) begangen, bevor er anfängt. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-148](../in-progress/slice-148-baseline-v5120-vendoring.md)
in `done/` — ohne gehobenen Pin gibt es keinen Stand, gegen den geprüft wird.

**Rückführungen:** `in-progress` → `next`, falls das Audit ein Produkt-Delta
ergibt, das eine eigene Anforderung braucht.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Regeltext (GF), Konventionsspeicher (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-26):
  [`BEO-011`](../observations.md) für jede Vollständigkeits-Aussage,
  [`BEO-012`](../observations.md) für jedes Zitat aus dem neuen Stand.

Slice-ID: slice-149. Betroffene IDs: — (kein `DC-`-Bezug). Module:
Harness-Regeltext. Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Audit-Arbeit an vendortem Fremdtext.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
