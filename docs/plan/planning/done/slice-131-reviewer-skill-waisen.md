# Slice slice-131: Der Reviewer-Skill trägt Regeln, die nirgends gerankt stehen

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-83-baseline-v5110-migration](../welle-83-baseline-v5110-migration.md)
(Etappe C, geschnitten vom Delta-Audit).

**Bezug:** Baseline-Regelwerk
[`grundlagen-source-precedence.md`](../../../../.harness/baseline/v5.11.0/regelwerk/grundlagen-source-precedence.md)
§Vollständigkeit (Kurs-Welle 94) und
[`modul-09-implementierung.md` §AGENTS.md-Regeln](../../../../.harness/baseline/v5.11.0/regelwerk/modul-09-implementierung.md#agentsmd-regeln-modul-9)
(*„Jede Hard Rule liegt in zwei Quadranten"* — **Durchsetzungs**-Verdopplung,
nicht Verortung; die Verortung regelt §Vollständigkeit);
[`.harness/skills/reviewer.md`](../../../../.harness/skills/reviewer.md);
[`AGENTS.md`](../../../../AGENTS.md) §3.

**Berührte Spec-Stellen:** — (Harness-Dateien; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

Der Vollständigkeits-Zensus aus [slice-129](../done/slice-129-baseline-v5110-delta-audit.md)
hat **fünf** Fundorte ergeben statt des einen bekannten. Der schwerste ist der
**Reviewer-Skill**: mehrere seiner Kategorien-Anker
sind **dort entstanden** („neuer HIGH-Eintrag seit 1.5.0", „neuer MEDIUM-Anker
seit 1.9.0"), statt einen in einer gerankten Quelle stehenden Ablauf
auszubuchstabieren.

Der Kanon erlaubt Artefakten außerhalb der Rangliste ausdrücklich, zu
**verweisen**, **auszuführen** und einen dort gerankten Ablauf
**auszubuchstabieren** — aber nichts **festzulegen**, was nicht dort steht. Ein
HIGH-Anker, der eine Prüf-Pflicht *einführt*, legt fest.

**Zwei kleinere Fundorte reiten mit**, weil sie dieselbe Frage stellen und je
zwei Zeilen sind: der `FOCUS_DISABLE`-Kommentar im
[`Makefile`](../../../../Makefile) (*„Spiegelt die `.d-check.yml`-modules-Liste;
wächst die dort, hier nachziehen"* — zugleich die Hälfte von
[`BEO-010`](../observations.md)) und die SHA-Pin-Konvention in den beiden
Workflow-Köpfen. Beide legen fest, beide stehen nirgends gerankt.

**Was dieser Slice nicht behauptet:** dass alle vierzehn normativ wirkenden
Stellen Waisen sind. Mehrere buchstabieren `AGENTS.md` §3.7 aus und sind damit
zulässig. Die Trennung ist Urteilsarbeit **je Anker** — genau die Prüffrage des
Kanons, und sie ist der Inhalt dieses Slice.

## 2. Vorgehen

1. **Je Kategorien-Anker eine Antwort:** buchstabiert er eine gerankte Regel aus
   (zulässig, Quelle nennen) oder legt er fest (Waise)? Ergebnis ist eine
   Tabelle, kein Urteil im Fließtext.
2. **Waisen umziehen, nicht löschen** — nach `AGENTS.md`, mit Herkunfts-Anker
   (`modul-09`: Hard Rules aus dem Steering Loop tragen ihn). Der Skill behält
   den Anker als **Verweis**.
3. Prüfen, ob dabei `AGENTS.md` zur Sammelstelle wächst. Der Wächter ist die
   **Prüffrage des Kanons** selbst — *steht die Aussage auch in einer gerankten
   Quelle?* —, denn genau dann ist ein Zuzug legitim: eine Waise steht nirgendwo
   sonst. Je Zuzug ist das zu belegen, nicht zu behaupten.
4. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Rückbau von Ankern.** Was festlegt, wandert; es verschwindet nicht.
- **Kein Gate für die Anker.** Ob eine Review-Regel mechanisierbar ist, ist
  eine eigene Frage — und nach `modul-09` bleibt sie ohne Gate **halb
  durchgesetzt**, was hier zu **benennen** und nicht zu heilen ist.
- **Nicht `CLAUDE.md`** ([slice-127](../done/slice-127-claude-md-pointer.md))
  und nicht das Workflow-Skelett (im Audit als vermutlich konform eingestuft).

## 4. Definition of Done

- [ ] Je Anker eine Antwort mit Quelle; die Tabelle ist vollständig **belegt**,
      nicht behauptet (`BEO-011`).
- [ ] Jede Waise steht in `AGENTS.md` mit Herkunfts-Anker; der Skill bzw. die
      Kommentar-Stelle verweist. Das gilt auch für die zwei kleineren Fundorte.
- [ ] `AGENTS.md` ist nicht zur Sammelstelle geworden — je Zuzug benannt, warum
      er nirgendwo sonst stehen kann.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Die Trennung „ausbuchstabieren vs. festlegen" ist Urteil, kein `grep`.**
  Ein zu weiter Schnitt macht `AGENTS.md` zur Sammelstelle, ein zu enger lässt
  Waisen stehen. — **Ausgang:** *(bei Closure)*
- **Der Skill ist selbst das Werkzeug des Reviews.** Wer ihn ändert, ändert das
  Instrument, mit dem die Änderung geprüft wird. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-129](../done/slice-129-baseline-v5110-delta-audit.md)
in `done/`.

**Rückführungen:** `in-progress` → `next`, falls die Prüfung ergibt, dass die
Anker in ihrer Mehrheit ausbuchstabieren — dann ist es kein Umzugs-, sondern ein
Deklarations-Slice.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Dateien (GF), Review-Infrastruktur (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23): **`BEO-011`**
  ist zentral — dieser Slice produziert eine Vollständigkeits-Aussage über die
  Anker, und genau diese Form ist in welle-82 achtmal gekippt. **`BEO-002`**
  für die Ränder jedes umgezogenen Ankers (Skill-Version, Index, Verweise).

Slice-ID: slice-131. Betroffene IDs: — (Harness-Dateien; keine Anforderung,
keine ADR, keine Adaption). Module: Harness-Dateien, Review-Infrastruktur. Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Umzug nach kanonischer Ortsregel.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
