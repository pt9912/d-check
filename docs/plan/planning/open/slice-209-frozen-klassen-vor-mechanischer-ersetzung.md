# Slice slice-209: Frozen-Klassen vor einer mechanischen Ersetzung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre.

**Bezug:** [`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
(3× erreicht, Ausgang *geplant*),
[`MR-025`](../../../../harness/conventions.md#mr-025) (der Geschwister-Eintrag:
Spiegel **vor** dem Editieren auflisten),
[`MR-052`](../../../../harness/conventions.md#mr-052) (ein historisches Zitat
ist über die git-Historie prüfbar).

**Berührte Spec-Stellen:** — *(keine; der Slice verkörpert eine
Steering-Loop-Regel, er ändert keine Anforderung)*

**Verantwortlich:** — · **Autor:** pt9912. **Datum:** 2026-09-07.

---

## 1. Ziel

Die Regel schreiben, die drei Anlässe verlangt haben: **Wer mechanisch über
den Baum ersetzt, listet die Frozen-Klassen vorher auf — und die Liste ist
über eine Eigenschaft definiert, nicht über Verzeichnisse.**

## 2. Vorgehen

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| neuer Eintrag unter [`harness/conventions/`](../../../../harness/conventions/) | neu | die Regel als Adaption; die Kennung wird beim Schreiben vergeben |
| [`harness/conventions.md`](../../../../harness/conventions.md) | update | Index-Zeile |
| [`BEO-ALL/mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/state.md) | update | Ausgang von *geplant* auf *verkörpert*, mit Herkunfts-Anker |

**Was die drei Anlässe gemeinsam haben, und es ist nicht das Verzeichnis:**
Beim ersten Mal (Register-Formatmigration) blieben `Accepted`-ADR-Kerne und
gesendete CRs unbedacht, obwohl die drei *benannten* Frozen-Verzeichnisse
korrekt ausgenommen waren. Beim zweiten (Pin-Hebung) dasselbe Muster. Beim
dritten (Pin-Hebung) gab es **gar keine** Ausnahme-Liste, und gehoben wurden
zehn `d-check:cite`-Direktiven in `done/`-Slices, deren Zitat gegen den alten
Tag geschrieben war. **Die Verzeichnis-Liste ist die falsche Abstraktion** —
maßgeblich ist die Eigenschaft *„zitiert den Stand seiner Zeit"*, und die
tragen auch Zeilen, die in keinem der drei Verzeichnisse stehen.

**Der Kanon hat dazu seit `v6.5.0` etwas zu sagen**
(`grundlagen-harness-dateien.md` §Ein einfrierendes Artefakt …): Er
unterscheidet **einfrierend** von **lebend** und zählt die einfrierenden auf —
Review-Report, Closure-Notiz, Archiv-Stub, `Accepted`-ADR, geschlossener
Slice. Ob die Regel dieses Repos damit **entbehrlich** wird oder ob sie den
*Vorgang* regelt, den der Kanon nicht kennt, ist die erste Frage des Slice.

## 3. Ausdrücklich NICHT in diesem Slice

- **Ein Sensor auf die Regel.** Ob eine Ersetzung ihre Frozen-Klassen
  aufgelistet hat, ist eine Aussage über einen **Akt**, nicht über einen
  ruhenden Zustand — dieselbe Lage wie bei
  [`AGENTS.md`](../../../../AGENTS.md) §3.6. Ein Gate darauf gäbe es nicht zu
  bauen.
- **Ein Retrofit der drei Anlässe.** Was geschehen ist, steht in den
  Evidence-Dateien; die geschlossenen Slices bleiben, wie sie sind.
- **Die Adoption der `v6.5.0`-Artefakt-Unterscheidung.** Sie liegt in
  [slice-208](../done/slice-208-v650-regel-adoption.md); dieser Slice
  **liest** sie nur, um zu entscheiden, ob seine Regel noch gebraucht wird.

## 4. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** Die **Vorfrage ist beantwortet**: Macht die
      `v6.5.0`-Unterscheidung *einfrierend / lebend* die Regel entbehrlich? Ein
      **Nein** braucht die Differenz — was regelt der Kanon nicht —, ein **Ja**
      schließt den Slice ohne neuen Eintrag und setzt den Register-Ausgang auf
      *verkörpert* mit dem Kanon als Zielort.
- [ ] **(2)** Fällt die Antwort auf **Nein**: Ein Konventions-Eintrag trägt die
      Regel samt der **Klassen-Liste über die Eigenschaft** und der Grenze, was
      sie nicht leistet.
- [ ] **(3)** Der Registereintrag trägt den Ausgang *verkörpert* mit
      auflösbarem Zielort und Herkunfts-Anker `seit slice-209`.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §5 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 5. Abnahme-Punkte / Risiken

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Die Regel könnte überflüssig sein, und das wäre das beste Ergebnis.**
  `v6.5.0` führt die Unterscheidung *einfrierend / lebend* mit einer
  Aufzählung, die alle drei Anlässe abdeckt. Ein eigener Eintrag daneben wäre
  dann eine zweite Quelle für dieselbe Regel — genau das, wovor die
  Source-Precedence warnt. Die Vorfrage steht deshalb **als DoD-Punkt**, nicht
  als Vorbemerkung. — **Ausgang:** \<offen\>
- **Eine Regel aus drei Anlässen ist keine Inventur**
  ([`rule-drawn-from-occasion-not-inventory`](../observations/BEO-ALL/rule-drawn-from-occasion-not-inventory/observation.md),
  7×). Drei ist die Kanon-Schwelle, aber die drei Anlässe sind alle vom selben
  Typ (Pin-/Format-Migration). Ob die Regel für andere mechanische Ersetzungen
  trägt, ist unbelegt und gehört als Grenze in den Eintrag. — **Ausgang:** \<offen\>
- **Kein Gate, und das ist eine Aussage über die Wirkung.** Die Regel hängt an
  der Disziplin des Ausführenden; sie verschiebt einen Fehler von *unsichtbar*
  nach *vermeidbar*, nicht nach *unmöglich*. Das gehört ausgeschrieben, sonst
  liest sie sich stärker, als sie ist. — **Ausgang:** \<offen\>

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-208](../done/slice-208-v650-regel-adoption.md)
liegt in `done/` — die Vorfrage aus DoD (1) braucht den adoptierten
`v6.5.0`-Stand, sonst urteilt sie über einen Kanon, den das Repo noch nicht
führt. WIP-Limit frei.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Zeigt die Vorfrage, dass die
  Klassen-Liste selbst eine Inventur über den Bestand verlangt, wird die
  Inventur ein eigener Slice.
- `in-progress` → `open` (blockiert): Ergibt die Vorfrage ein **Ja**, ist der
  Slice nicht blockiert, sondern **fertig** — er schließt mit DoD (1) und ohne
  DoD (2). Das ist kein Rückführungs-Fall, sondern der kürzere Weg.

**Closure-Trigger.** Zwei beobachtbare Kriterien und ein Lerneintrag: (a) die
Vorfrage ist beantwortet und die Antwort steht im Slice; (b) der
Registereintrag trägt einen Ausgang mit auflösbarem Zielort.

## 7. Vorgelagert (vor der Modus-Begründung)

\<entsteht spätestens bei der Beanspruchung — ein Plan in `open/` trägt die drei
Vorprüfungen noch nicht\>

## 8. Sub-Area-Modus-Begründung

\<entsteht mit den Vorprüfungen bei der Beanspruchung\>

## 9. Closure-Notiz (nach `done/`)

\<wird vor dem `git mv` nach `done/` gefüllt\>
