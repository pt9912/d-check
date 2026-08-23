# Slice slice-135: Ein Sensor auf die `uses:`-Pins der Workflows

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-84-durchsetzung](../welle-84-durchsetzung.md).

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §3.9 (Action-Referenzen
SHA-gepinnt; *Auflösungs-Trigger: ein Sensor, der `uses:`-Pins prüft*);
geschnitten vom Zensus in [slice-132](../in-progress/slice-132-hard-rule-zensus.md).

**Berührte Spec-Stellen:** — (Harness-Gate; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

§3.9 ist die einzige der einseitigen Regeln, die einen **auflösenden** Trigger
trägt statt *permanent* — weil sie als einzige mechanisierbar ist: ein
`uses:`-Eintrag ohne 40-stelligen Commit-SHA ist ein `grep`-barer Zustand, kein
Urteil.

Der Slice löst diesen Trigger ein. Er ist damit auch die Probe auf die
Wellen-Regel *„keine Heuristik-Wächter"*: hier ist ein Wächter **kein**
Heuristik-Wächter, und der Unterschied ist zu zeigen, nicht zu behaupten.

## 2. Vorgehen

1. **Ort entscheiden, bevor gebaut wird.** Drei Kandidaten, jeder mit einer
   Folge: ein `make`-Target mit `grep` (billig, aber ein weiteres Host-Skript);
   ein CI-Schritt (bindet nur die CI, nicht den lokalen Lauf); ein
   d-check-Modul (das Produkt scannt heute **nur Markdown** — YAML wäre eine
   neue Eingabe-Klasse und fiele unter §3.8). **Im Slice entscheiden**, nicht
   vorab.
2. **Zusage genau schneiden:** geprüft wird die **Form** (`@` gefolgt von 40
   Hex-Zeichen, dahinter ein Tag-Kommentar), nicht die **Gültigkeit** des SHA.
   Ob der Pin auf den Commit zeigt, den der Kommentar behauptet, ist eine
   Netz-Frage und ausdrücklich **nicht** Teil der Zusage.
3. **Bewusstes Brechen:** ein `uses:` auf einen beweglichen Tag ⇒ rot; ein
   SHA ohne Tag-Kommentar ⇒ rot; Rückbau ⇒ grün.
4. `AGENTS.md` §3.9 nachziehen (Trigger eingelöst) und, falls ein Target
   entsteht, §4 und die Sensors-Tabelle — sonst ist der Slice
   `gate-consistency`-rot.
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Prüfung, ob der SHA existiert oder zum Tag passt.** Das ist Netz und
  gehört zur Freshness-Familie, nicht hierher.
- **Kein Auto-Pinning.** Der Sensor meldet; das Pinnen bleibt ein bewusster Akt.
- **Keine Ausweitung auf andere YAML-Felder.** Ein Sensor, der „Workflows
  prüft", wäre eine Zusage, die dieser Slice nicht einlöst.

## 4. Definition of Done

- [ ] Der Ort ist **begründet** gewählt, mit der benannten Folge der beiden
      verworfenen Kandidaten.
- [ ] Beide Verstoß-Formen sind rot gesehen — beweglicher Tag **und** SHA ohne
      Tag-Kommentar —, der Rückbau grün.
- [ ] Die Zusage ist auf die **Form** beschränkt und sagt das hin.
- [ ] `AGENTS.md` §3.9 trägt den eingelösten Trigger; falls ein Target entsteht,
      sind §4 und die Sensors-Tabelle nachgezogen und `gate-consistency` grün.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein Gate auf drei Dateien ist billig zu bauen und leicht zu überdehnen.**
  „Die Workflows sind gehärtet" wäre die Aussage, die der Sensor **nicht**
  trägt. — **Ausgang:** *(bei Closure)*
- **Der d-check-Modul-Weg würde den Gegenstand des Produkts weiten** (YAML statt
  Markdown) und fiele damit unter §3.8 — eine Eingabe, die das Modul liest, aber
  nicht scannt. Wird er gewählt, ist das eine ADR-Frage, kein Slice-Entscheid. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-132](../in-progress/slice-132-hard-rule-zensus.md)
in `done/`, WIP-Limit frei. Hängt **nicht** an
[slice-134](../open/slice-134-nolintlint.md).

**Rückführungen:** `in-progress` → `next`, falls die Ortswahl auf das
d-check-Modul fällt — dann ist es ein Produkt-Slice mit ADR und gehört in eine
andere Welle.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** CI/Workflows (GF), Gate-Landschaft (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23):
  [`BEO-011`](../observations.md) für die Zusage des Sensors — sie ist auf die
  Form zu schneiden, nicht auf „gehärtet". [`BEO-010`](../observations.md),
  falls ein Target entsteht: es erscheint dann in mehreren Doku-Tabellen.

Slice-ID: slice-135. Betroffene IDs: — (Harness-Gate; keine Anforderung).
Module: CI/Workflows, Gate-Landschaft. Gates: `make gate-consistency`,
`make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — neuer Wächter auf eigenem Bestand.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
