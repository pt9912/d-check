# Slice slice-219: Das Release, das `CO-001` auflöst

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in),
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
— und [`CO-001`](../../carveouts/CO-001-vcs-range-stiller-skip.md), dessen
Auflösung dieser Slice trägt.

**Berührte Spec-Stellen:** — (der Slice berührt keine Spec-Stelle).

**Verantwortlich:** — (wird bei der Beanspruchung gesetzt).

**Autor:** pt9912.

## 1. Ziel und Abgrenzung

**Ziel.** Der Fix aus [slice-218](../done/slice-218-go-git-pack-namenskonvention.md)
erreicht das **publizierte** Image, und
[`CO-001`](../../carveouts/CO-001-vcs-range-stiller-skip.md) wird aufgelöst.
Solange das nicht geschehen ist, fährt jeder Konsument mit gepinntem
`v0.75.0` weiterhin den stillen Pfad — **der Fix im Quellstand hilft ihm
nicht.**

**Warum das ein eigener Vorgang ist und nicht der Rest von slice-218.**
Ein Carveout braucht einen Folge-Slice, der ihn **überlebt**
(Baseline-Regelwerk `modul-07-carveouts.md`). slice-218 schließt, sobald der
Fix committet ist — er kann die Auflösung also nicht tragen, und ein
Carveout, dessen Folge-Slice vor ihm schließt, hat faktisch keinen. Genau das
hat der unabhängige Review an der ersten Fassung von `CO-001` beanstandet.

**Abgrenzung — drei Punkte, jeder mit Grund:**

1. **Kein Produkt-Code.** Der Fix ist in slice-218 geliefert und dort
   gemessen; dieser Slice liefert ihn **aus** — *es wäre ein anderer Vorgang*,
   den Fix hier nachzubessern.
2. **Kein Vorziehen des Release-Zeitpunkts.** Wann geschnitten wird, ist eine
   Entscheidung des Auftraggebers, nicht dieses Plans; der Slice beschreibt,
   *was* dann zu tun ist — *Bestand bleibt bewusst stehen*.
3. **Keine weiteren Carveouts.** Es gibt heute genau einen; sollte bis dahin
   ein zweiter entstehen, löst dieser Slice ihn **nicht** mit —
   *Schicht-Abgrenzung*, sonst wächst er unbemerkt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** Ein Release nach `docs/user/releasing.md` ist veröffentlicht und
      am **gezogenen** Image verifiziert (Digest, OCI-Label, Smoke).
- [ ] **(2)** Die Auflösung ist **gemessen, nicht angenommen**: dasselbe
      Probe-Repo mit partiell unsichtbarem Pack, gegen das **publizierte**
      Image gefahren — vorher `0 Befund(e)`/Exit 0, danach
      `nicht lesbares Objekt …`/Exit 2. Beide Ausprägungen (unsichtbares
      BASE-Blob **und** unsichtbares BASE-Tree).
- [ ] **(3)** `CO-001` ist aufgelöst: Verifikations-Haken abgehakt, Datei per
      reinem `git mv` nach `docs/plan/carveouts/done/` <!-- d-check:ignore (entsteht erst mit dieser Auflösung) -->, Index in
      `docs/plan/carveouts/README.md` und die Bindung-Spalte in
      `harness/README.md` §Sensors nachgezogen.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

1. Release-Prep nach `docs/user/releasing.md` §Release-Prep in **einem**
   Commit; die Version trägt einen **Patch**- oder **Minor**-Bump je nachdem,
   was bis dahin sonst noch aufgelaufen ist.
2. `make ci` grün, Tag, Pipeline beobachten, `docker pull` als Beweis.
3. Die Auflösungs-Probe aus DoD (2) gegen das publizierte Image fahren —
   **vor** dem Abhaken der Verifikations-Liste in `CO-001`.
4. Carveout-Move und Index-Nachzug.

## 4. Trigger

**Beanspruchung:** slice-218 ist geschlossen **und** der Auftraggeber gibt das
Release frei.

**Rückführung nach `open/`** (`in-progress→open`): wenn die Auflösungs-Probe
gegen das publizierte Image den stillen Pfad **weiterhin** zeigt — dann trägt
der Fix nicht, und der Carveout bleibt aktiv.

## 5. Closure-Trigger

DoD (1) bis (3) abgehakt, `CO-001` liegt in `done/`, `make gates` grün.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Der Carveout altert, ohne dass etwas meldet.** Kein Sensor prüft, wie
  lange `CO-001` schon aktiv ist; die kanonische Frist misst in **Wellen**, und
  dieses Repo arbeitet wellenlos — eine benannte Lücke des Kanons
  (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).
  Der einzige Träger ist dieser Slice in `open/`. — **Ausgang:** \<offen\>
- **Ein Release kann aus anderen Gründen fällig werden und den Carveout
  „nebenbei" auflösen**, ohne dass jemand die Probe aus DoD (2) fährt. Dann
  wäre der Carveout formal offen, obwohl der Defekt weg ist — oder schlimmer,
  er würde geschlossen, ohne dass die Auflösung gemessen wurde. — **Ausgang:**
  \<offen\>

## 7. Closure-Notiz

\<wird vor dem `git mv` nach `done/` gefüllt\>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende.

**Die drei Vorprüfungen entstehen spätestens bei der Beanspruchung**
([`AGENTS.md`](../../../../AGENTS.md) §5) — dieser Plan liegt in `open/` und
trägt sie noch nicht. Das ist die Regel, nicht eine Auslassung: Ein
Register-Stand und ein Nachtlauf-Stand, die beim Anlegen gelesen werden,
wären beim Beanspruchen alt.
