# Slice slice-168: Die Titel-Spalte trägt einen Titel, nicht die Entscheidung

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`docs/plan/adr/README.md`](../../adr/README.md) (der Index selbst);
Baseline-Vorlage `templates/docs/plan/adr/README.template.md` (*„Derivativ:
Quelle der Wahrheit sind die ADR-Dateien"*);
[ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md) (die
Präzedenz: ein Zeichenketten-Wächter weicht einer `structure`-Regel).

**Berührte Spec-Stellen:** [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) (neunte Bedingung) und ihre Verfeinerung in der Spezifikation.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-28.

---

## 1. Ziel

**Die Spalte heißt „Titel" und trägt im Median das Fünffache eines Titels — und
das Gate dafür soll das Produkt stellen, nicht ein Regex.**

Gemessen über alle 68 Zeilen: die Titel-Spalte hat **Median 442**, Maximum
**2206** Zeichen. Die **H1-Titel der ADRs** liegen bei **Median 77**, Maximum
**158** — und **alle 68 sind genau ein Satz**. Ein Satz reicht für einen Titel;
das ist keine Schätzung, sondern der gemessene Bestand.

**Das ist mehr als Kosmetik.** Der Kopf derselben Datei sagt: *„Derivativ:
Quelle der Wahrheit sind die **ADR-Dateien**; dieser Index ist eine
Bequemlichkeits-Sicht."* Eine 2206-Zeichen-Zelle macht den Index zu einer
**zweiten Quelle** — sie wiederholt die Entscheidung, statt auf sie zu zeigen,
und kann von der ADR abdriften, ohne dass ein Gate es merkt.

**Warum kein Regex.** Ein `forbid-pattern` kann eine Zellenlänge ausdrücken —
ich habe es gemessen, es zählt Zeichen und trifft 43 der 44 zu langen Zeilen.
Aber es taugt nicht als Gate: der Befund nennt den **Abschnitt**, nicht die
Zeile und nicht die Spalte; das Muster hängt an der **Form der ID-Zelle** und
greift still nicht mehr, wenn die sich ändert; und eine Titel-Zelle mit einem
`|` zerteilt den Lauf und entkommt (im Bestand genau eine). Ein Wächter, dessen
Treffgenauigkeit an einer Zeichenkette hängt, ist die Bauform, die dieses Repo
zuletzt in [ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md)
zugunsten einer `structure`-Regel aufgegeben hat.

## 2. Vorgehen

1. **Eine neunte `structure`-Bedingung, kein Muster.** Sie misst die
   **Zeichenzahl** einer Tabellenzelle in einer **benannten Spalte** und meldet
   die **Zeile**, in der sie steht. Damit sagt der Befund, wo etwas zu tun ist —
   das kann ein `forbid-pattern` nicht.
2. **Zeichen, nicht Bytes.** Gemessen an einer Probe mit sechzig Umlauten (120
   Byte, 60 Zeichen): eine Zeichen-Schwelle ist die, die zu einem Titel passt.
3. **Die Schwelle kommt aus dem Bestand, und ihre Herkunft wird benannt.** Der
   längste echte H1 misst 158; 200 lässt Luft. Das ist eine Regel aus dem
   heutigen Bestand ([`BEO-011`](../observations.md)) — ein künftiger längerer
   Titel fiele auf, und das ist gewollt: „ein Satz" ist die Regel, 200 ihr
   grober Wächter.
4. **Erst die Bedingung, dann ihr erster Konsument.** Die 44 zu langen Zeilen
   werden auf den H1 der zugehörigen ADR zurückgeschnitten — die Sicht gibt
   wieder, was die Datei sagt, statt sie zu wiederholen.
5. **Der Rest der Zelle geht nicht verloren, er ist schon da.** Wo eine Zeile
   etwas trägt, das **nur** dort steht, ist das der eigentliche Befund und
   wandert in die ADR, nicht in den Papierkorb.
6. Lastenheft, ADR, Implementierung, Tests, Handbuch; `make gates`; unabhängiger
   Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein `forbid-pattern` als Ersatz.** Die Messung dazu ist gefahren und steht
  in §1; sie ist die Begründung für die Bedingung, nicht ihre Alternative.
- **Keine anderen Spalten.** `Bezug` ist mit Median 251 lang, aber zu Recht —
  Markdown-Links auf Lastenheft-Anker. Die Bedingung ist deshalb
  **spalten**-gebunden, nicht zeilen-pauschal.
- **Keine Änderung an den ADRs selbst.** Sie sind `Accepted` und immutabel; ihr
  H1 ist die Quelle, nicht der Gegenstand.
- **Kein Fremd-Repo.** Für `ai-harness-init` gilt die ausdrückliche Weisung
  *„nichts ändern"*.

## 4. Definition of Done

- [ ] Die neunte `structure`-Bedingung existiert: Spalte benannt, Schwelle in
      **Zeichen**, Befund auf der **Zeile** der zu langen Zelle, eigener
      Grund-Code.
- [ ] [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) trägt sie mit Akzeptanzkriterien-Trio; ADR mit
      Fitness Function; Spezifikation, Handbuch und `operations.md` nachgezogen.
- [ ] Negativtests: zu lange Zelle, Grenzfall genau auf der Schwelle, Umlaute
      (Zeichen ≠ Bytes), Zelle mit `|`, fehlende Spalte, leere Prüfmenge.
- [ ] Jede Titel-Zelle gibt den H1 der zugehörigen ADR wieder; die 24 bereits
      kurzen Zeilen sind unverändert.
- [ ] Kein Inhalt ist verloren: was nur im Index stand, ist benannt und
      versorgt — oder es gab keinen solchen Fall, und das ist gemessen.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Kürzen kann Inhalt vernichten, der nur hier stand.** Der Index ist
  *derivativ* — trägt er etwas Eigenes, ist das der Defekt, nicht die Kürzung.
  Aber „müsste in der ADR stehen" ist eine Annahme, bis sie geprüft ist. —
  **Ausgang:** *(bei Closure)*
- **Eine Schwelle aus dem heutigen Bestand ist aus dem Anlass gezogen**
  ([`BEO-011`](../observations.md)). 200 liegt über dem längsten **heutigen**
  H1; ein künftiger längerer Titel fiele auf. — **Ausgang:** *(bei Closure)*
- **Eine Bedingung, die eine Spalte über ihre Position anspricht, hängt an der
  Spalten-Reihenfolge.** Wird eine Spalte eingefügt, misst sie still die
  falsche. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `open`, falls sich beim Kürzen zeigt, dass
mehrere Zeilen Inhalt tragen, den die ADR **nicht** hat — dann ist der Befund
ein anderer (unvollständige ADRs), und die Kürzung wäre Datenverlust.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Register-/Doku-Form (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-28):
  [`BEO-011`](../observations.md) — die Regel aus dem **Bestand**, nicht aus dem
  Anlass: die Schwelle 200 stammt aus den heutigen H1-Längen und ist damit
  genau diese Bauform;
  [`BEO-013`](../observations.md) — ein Wächter, der nichts mehr fängt, bleibt
  stehen: das Muster hängt an der Form der ID-Zelle;
  [`BEO-020`](../observations.md) — die gezählte Menge benennen: „Median 442"
  gilt der **Titel-Spalte**, nicht der Zeile.
- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): `upstream-drift.yml`
  zuletzt **manuell** ausgelöst (13/13 Achsen gelaufen, ein Befund —
  `golangci-lint` 2.13.2, inzwischen gehoben). Der **geplante** Lauf des 28.
  ist ausgefallen, Ursache ungeklärt; der nächste entscheidet, ob es
  systematisch ist. `image-scan.yml` hatte noch keinen geplanten Lauf.

Slice-ID: slice-168. Betroffene IDs: [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in).
Module: `structure`. Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Formkorrektur an einem vorhandenen Register
plus eine Regel; kein Fremdsystem, keine Reconciliation.

## 9. Closure-Notiz (nach `done/`)
