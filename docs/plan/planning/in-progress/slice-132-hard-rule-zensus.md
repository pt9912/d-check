# Slice slice-132: Je Hard Rule eine Antwort — welcher Gate-Lauf trägt sie?

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-84-durchsetzung](../welle-84-durchsetzung.md).

**Bezug:** Baseline-Regelwerk
[`modul-09-implementierung.md` §AGENTS.md-Regeln](../../../../.harness/baseline/v5.11.0/regelwerk/modul-09-implementierung.md#agentsmd-regeln-modul-9)
(*„Jede Hard Rule liegt in zwei Quadranten … nur in einem ist halb
durchgesetzt"*) und
[`modul-13-quality-gates.md` §Hard Rule (Doku-Disziplin)](../../../../.harness/baseline/v5.11.0/regelwerk/modul-13-quality-gates.md);
[`AGENTS.md`](../../../../AGENTS.md) §3 und §5.

**Berührte Spec-Stellen:** — (Harness-Dateien; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

`AGENTS.md` §3 trägt neun Hard Rules, dazu die Botschafts-Regel in §5. Dieser
Slice beantwortet für **jede einzelne**: welcher Gate-Lauf trägt ihre
Feedback-Hälfte — mit Target **und** Befund-Code — oder ist sie einseitig?

**Was dieser Slice ausdrücklich nicht behauptet:** dass die Mehrheit einseitig
ist. Belegt sind heute genau **drei** (§3.8, §3.9 und die §5-Zeile, alle aus
welle-83 und dort als ohne Gate ausgewiesen). Für die übrigen ist die Frage
**offen**, nicht beantwortet — und beide Ausgänge sind Ergebnisse.

**Der Beleg ist ein Gate-Lauf, keine Zuordnung.** Ein Target zu nennen, das
plausibel klingt, ist die Klasse, die in welle-82 achtmal gekippt ist
([`BEO-011`](../observations.md)). Wo eine Zeile behauptet, Gate X trage Regel Y,
muss ein **konstruierter Verstoß** X rot färben — genau das *„Bewusste Brechen"*,
das
[`modul-13`](../../../../.harness/baseline/v5.11.0/regelwerk/modul-13-quality-gates.md#adr-zur-fitness-function)
für die Fitness Function verlangt.

## 2. Vorgehen

1. **Je Regel eine Zeile** in einer Tabelle: Regel · tragender Gate-Lauf
   (Target + Befund-Code) · Beleg-Form · Verdikt (*gedeckt* / *einseitig* /
   *teilgedeckt*).
2. **Teilgedeckt ist ein eigenes Verdikt**, kein aufgerundetes *gedeckt*: §3.3
   etwa koppelt Lifecycle-Move und Verweise über `planning-check`, sagt aber
   auch etwas über **Commit-Zerlegung**, das kein Gate sieht.
3. **Belegen, nicht zuordnen:** je *gedeckt*-Zeile ein konstruierter Verstoß,
   der das genannte Gate rot färbt, und der Rückbau grün.
4. **Die einseitigen ausweisen** — im Regeltext selbst, in der Form, die §3.7
   und §3.9 schon führen (*Auflösungs-Trigger*), ohne Forensik.
5. **Schneiden:** was baubar ist, wird als Slice benannt; was nicht baubar ist,
   bleibt ausgewiesen. Beides kommt in die Welle-Datei §4.
6. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Gates bauen.** Der Zensus misst und schneidet; das Bauen ist Folge.
- **Kein Heuristik-Wächter**, um eine Zeile von *einseitig* auf *gedeckt* zu
  heben. Ein behauptetes Gate ist schlechter als eine ausgewiesene Lücke.
- **Nicht die Regeln ändern.** Wer beim Zensus merkt, dass eine Regel unscharf
  ist, notiert es — geändert wird sie in einem eigenen Slice.

## 4. Definition of Done

- [ ] Je Hard Rule in §3 **und** der Botschafts-Regel in §5 genau eine Zeile;
      keine Regel ohne Zeile, keine Zeile ohne Regel.
- [ ] Jede *gedeckt*-Zeile trägt einen **konstruierten Verstoß mit rotem
      Gate-Exit** als Beleg, nicht eine Zuordnung.
- [ ] Jede *einseitig*-Zeile ist im Regeltext als solche ausgewiesen.
- [ ] Der Schnitt steht in der Welle-Datei §4 — baubar als Slice benannt, nicht
      baubar als ausgewiesen.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Eine Vollständigkeits-Aussage über eine ganze Datei ist die Form, die in
  welle-82 achtmal gekippt ist** ([`BEO-011`](../observations.md)). Diesmal
  zusätzlich gefährdet durch die feinere Auflösung, die slice-127 gelehrt hat:
  eine *Regel* ist nicht dasselbe wie ein *Abschnitt* — §3.7 trägt mehrere
  Aussagen. — **Ausgang:** *(bei Closure)*
- **„Gate genannt" ist nicht „Regel getragen".** Ein Gate kann einen Teil der
  Regel prüfen und der Zeile trotzdem ein *gedeckt* verschaffen. Genau davor
  schützt der konstruierte Verstoß — und nur, wenn er die **Regel** bricht, nicht
  irgendetwas, das dasselbe Gate rot färbt. — **Ausgang:** *(bei Closure)*
- **Der Zensus kann den Zuschnitt der Welle umwerfen.** Findet er mehr
  Baubares als erwartet, wird die Welle zu groß; findet er weniger, war die
  Produkt-Welle eine Vermutung. Beides ist ein Ergebnis, kein Fehler. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): [welle-84](../welle-84-durchsetzung.md)
eröffnet, WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls der Zensus ergibt, dass eine
Hard Rule vor der Durchsetzung **umformuliert** werden muss — dann ist das eine
Auftraggeber-Frage und kein Nachzug.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Dateien (GF), Gate-Landschaft (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23):
  [`BEO-011`](../observations.md) ist zentral — dieser Slice produziert genau
  eine Vollständigkeits-Aussage. [`BEO-012`](../observations.md) ebenso, denn
  jede Zeile zitiert eine Regel und behauptet eine Reichweite; das
  Geltungs-Feld ist zu lesen, nicht der Titel. [`BEO-007`](../observations.md)
  für jeden Beleg-Lauf: der Exit gehört gelesen, nicht hinter eine Pipe.

Slice-ID: slice-132. Betroffene IDs: — (Harness-Dateien; keine Anforderung,
keine ADR, keine Adaption). Module: Harness-Dateien, Gate-Landschaft.
Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Bestandsaufnahme an eigenen Dateien.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
