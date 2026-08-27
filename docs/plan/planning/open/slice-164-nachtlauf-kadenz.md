# Slice slice-164: Der Nachtlauf hat keinen Adressaten — wer liest das Rot, und wann?

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`upstream-drift.yml`](../../../../.github/workflows/upstream-drift.yml)
(der Nachtlauf und seine bereits benannte Grenze);
[slice-142](../done/slice-142-freshness-weitere-achsen.md) (die zwölf Achsen);
[slice-161](../done/slice-161-sechs-pins-heben.md) (die Hebung, die den
dritten §5-Punkt offen ließ); [`MR-051`](../../../../harness/conventions.md#mr-051)
(die Bump-Nachwirkung, die planmäßig meldet).

**Berührte Spec-Stellen:** — (Harness-Betrieb; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

Zwölf Achsen wachen über jeden gepinnten Fremd-Bestand. Sie melden **korrekt**
— und **an niemanden**. Der Job fällt rot aus und ist nur in der
Actions-Übersicht des Repos sichtbar; es gibt keinen Push-Kanal, keine
Benachrichtigung, keinen festen Lese-Schritt.

**Die Lücke ist nicht neu und nicht verdeckt:** der Workflow-Kopf trägt sie als
benannte Grenze — *ein dauerroter Nachtlauf ist wieder derselbe verwaiste
Sensor*. [slice-161](../done/slice-161-sechs-pins-heben.md) hat das Rot geräumt
und den Punkt als **weiter offen** an dieses Register zurückgegeben. Genau
deshalb ist er jetzt fällig: solange alles grün ist, kostet die Entscheidung
nichts; beim nächsten Rot ist sie wieder dringend und wieder unbequem.

**Der nächste rote Lauf ist absehbar**, nicht hypothetisch: nach
[`MR-051`](../../../../harness/conventions.md#mr-051) meldet die Zitat-Prüfung
bei jedem Baseline-Bump planmäßig, und Fremd-Releases erscheinen ohne
Ankündigung.

## 2. Vorgehen

1. **Zuerst die Frage, nicht das Werkzeug:** wer ist der Adressat, und in
   welchem Takt soll er lesen? Eine Kadenz ohne benannten Leser ist dieselbe
   Verwaisung eine Ebene höher.
2. **Den vorhandenen Lese-Schritt prüfen, bevor ein neuer gebaut wird.** Die
   Slice-Planung hat bereits einen Sichtungs-Schritt (Beobachtungs-Register).
   Ob der Nachtlauf-Stand dort hineingehört, ist billiger zu haben als jeder
   Benachrichtigungs-Kanal — und braucht kein Geheimnis, keinen Dienst, keine
   zweite Fehlerquelle.
3. **Erst dann die Kanal-Frage**, falls Schritt 2 nicht trägt: was GitHub ohne
   Zusatz-Dienst kann (Actions-Benachrichtigung je Watch-Einstellung, ein Issue
   je rotem Lauf) gegen das, was einen Dienst bräuchte.
4. **Die Rausch-Frage mitentscheiden.** Ein Kanal, der bei jedem Bump meldet,
   erzieht zum Wegklicken. Ob eine planmäßige Meldung
   ([`MR-051`](../../../../harness/conventions.md#mr-051)) anders behandelt
   wird als eine unerwartete, gehört zur Kadenz.
5. Nur bauen, was die Antwort trägt; eine Entscheidung **gegen** einen Kanal
   ist ebenso auszuweisen wie eine dafür.
6. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Auto-Bump.** Die Achsen melden; wer hebt, entscheidet.
- **Keine Aufnahme netzhaltiger Achsen in `gates`.** Der Nachtlauf bleibt der
  Bindepunkt.
- **Kein Fremd-Dienst mit Geheimnis**, solange Schritt 2 oder 3 trägt — ein
  Wächter, der ein Token braucht, hat eine neue Ausfall-Achse.

## 4. Definition of Done

- [ ] Der **Adressat** ist benannt, nicht nur der Kanal.
- [ ] Der vorhandene Lese-Schritt ist geprüft, bevor ein neuer gebaut wird.
- [ ] Die Rausch-Frage ist entschieden (planmäßige gegen unerwartete Meldung).
- [ ] Was gebaut wird, ist gefahren — mit gelesener Ausgabe, nicht behauptet.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein Kanal ist billig zu bauen und teuer zu ignorieren** — dieselbe Klasse
  wie bei den Achsen selbst, nur eine Ebene höher. — **Ausgang:**
- **Ein Lese-Schritt in der Slice-Planung greift nur, wenn ein Slice geplant
  wird.** In einer Pause liest niemand. — **Ausgang:**
- **Die Entscheidung könnte den Zustand mit dem Werkzeug verwechseln.** „Wir
  bauen eine Benachrichtigung" beantwortet nicht, wer sie liest. — **Ausgang:**

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Adressaten-Frage ein
Auftraggeber-Entscheid ist.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** CI (GF), Harness-Prozess (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27):
  [`BEO-020`](../observations.md) — über die gemessene Menge reden, nicht über
  die naheliegende; [`BEO-011`](../observations.md) — die Regel aus dem
  Bestand, nicht aus dem Anlass.

Slice-ID: slice-164. Betroffene IDs: — (Harness-Betrieb; keine Anforderung).
Module: CI, Harness-Prozess. Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Betriebs-Entscheid an vorhandener Mechanik.

## 9. Closure-Notiz (nach `done/`)

— (offen)
