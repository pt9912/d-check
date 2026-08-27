# Slice slice-163: Baseline-Zitate außerhalb des Konventionsspeichers

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`MR-039`](../../../../harness/conventions.md#mr-039) (der
Geltungsbereich, der diese Dokumente ausdrücklich mitnennt);
[`MR-051`](../../../../harness/conventions.md#mr-051) (das Neu-Ankern beim
Bump); [slice-152](../done/slice-152-citations-scharfschalten.md) (der Zug, der
den Speicher abgearbeitet hat).

**Berührte Spec-Stellen:** — (Doku-Bestand; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

[`MR-039`](../../../../harness/conventions.md#mr-039) §Geltungsbereich nennt
**alle lebenden** Dokumente dieses Repos — *„Konventionsspeicher, `AGENTS.md`,
`harness/README.md`, die Skills und die lebenden Planungs-Dokumente"*.
[slice-152](../done/slice-152-citations-scharfschalten.md) hat den
**Konventionsspeicher** abgearbeitet und den Rest ausdrücklich ausgeklammert
(§3 dort: *„Keine Ausweitung auf Zitate außerhalb des Konventionsspeichers in
diesem Zug"*).

**Der Rest ist benannt, nicht geschätzt:** der Review von slice-152 hat
mindestens zwei **wörtliche, auszeichenbare** Baseline-Zitate außerhalb des
Speichers gefunden — je eines in [`AGENTS.md`](../../../../AGENTS.md) und im
Reviewer-Skill. Wie viele es insgesamt sind, ist die erste Messung dieses Slice
und hier **nicht** vorweggenommen.

## 2. Vorgehen

1. **Den Bestand zählen**, mit derselben Methode wie
   [slice-152](../done/slice-152-citations-scharfschalten.md): Zitate ab 16
   Zeichen, Produkt-Lexik (Fence-Automat, absatzweise Backtick-Spannen), je
   Zitat die drei Klassen — wörtlich / abweichend / gegen den gepinnten Stand
   nicht prüfbar.
2. **Bei Abweichung die Richtung entscheiden**, nach der Linie aus
   [slice-152](../done/slice-152-citations-scharfschalten.md): was die Quelle
   **hat** und das Zitat weglässt, ist ein Transkriptions-Fehler und wird
   korrigiert; was das Zitat **hinzufügt**, ist ein Autoren-Akt und wird
   deklariert statt korrigiert.
3. **Die Skills sind der heikle Teil.** Sie werden von Agenten gelesen, nicht
   von Menschen gerendert; eine Direktive am Zeilenende ist dort unauffällig,
   aber die Zitate tragen die Beweislast der Anker. Vor der Auszeichnung prüfen,
   ob das Zitat überhaupt den gepinnten Stand meint.
4. Bewusstes Brechen je gedeckter Form, **Ursache gelesen**; Rückbau grün.
5. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Auszeichnung in `done/`, `docs/reviews/` oder `conventions/done/`.**
  Eingefrorene Dokumente zitieren den Stand ihrer Zeit.
- **Keine Änderung an der Paarungsregel** des Moduls, auch nicht für die
  Delta-Tabelle aus [`MR-039`](../../../../harness/conventions.md#mr-039), die
  [slice-152](../done/slice-152-citations-scharfschalten.md) als nicht
  adressierbar ausgewiesen hat.
- **Keine neue Ventil-Achse.**

## 4. Definition of Done

- [ ] Der Bestand ist **gezählt**, nach den drei Klassen getrennt.
- [ ] Je Abweichung eine Entscheidung nach der Richtungs-Linie, mit Begründung.
- [ ] Was ausgezeichnet wird, ist ausgezeichnet; was nicht, ist **benannt**.
- [ ] Ein konstruierter Verstoß je gedeckter Form, **Ursache gelesen**.
- [ ] Doku-Currency geprüft (Handbuch, `README`-Fassungen, `CHANGELOG`).
- [ ] `make gates` grün (Exit explizit), `make fullbuild` grün; unabhängiger
      Review.

## 5. Abnahme-Punkte / Risiken

- **Jede neue Direktive erhöht die Bump-Last.** [`MR-051`](../../../../harness/conventions.md#mr-051)
  legt das Neu-Ankern in die Prozedur, aber die Zahl der neu zu ankernden
  Spannen wächst mit jedem Zug. — **Ausgang:**
- **Ein Zitat im Skill trägt Beweislast.** Wird es beim Korrigieren gekürzt oder
  verschoben, ändert sich womöglich, was der Anker dem Reviewer sagt. —
  **Ausgang:**
- **Die Richtungs-Linie ist neu und ungeprüft.** Sie stammt aus einem einzigen
  Zug über einen einzigen Bestand; ob sie außerhalb des Speichers trägt, ist
  offen. — **Ausgang:**

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Messung zeigt, dass ein
Skill-Zitat den gepinnten Stand gar nicht meint.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Doku (GF), Konventionsspeicher (GF), Skills (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27):
  [`BEO-020`](../observations.md) — die gemessene Menge muss die sein, über die
  geredet wird; [`BEO-011`](../observations.md) — die Regel gehört aus dem
  Bestand, nicht aus dem Anlass; [`BEO-012`](../observations.md) — ein Zitat
  trägt nur, was in seinem Geltungsbereich steht.

Slice-ID: slice-163. Betroffene IDs: — (Doku-Bestand; keine Anforderung).
Module: `citations`. Gates: `make doc-check`, `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Fortsetzung eines abgeschlossenen Zugs auf
den benannten Rest.

## 9. Closure-Notiz (nach `done/`)

— (offen)
