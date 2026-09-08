# Slice slice-222: Baseline-Pin auf `v6.6.0` — mechanisch, ohne Urteil über den Delta

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`MR-011`](../../../../harness/conventions.md#mr-011)-Pin-Serie
(aktueller Eintrag: [`MR-067`](../../../../harness/conventions.md#mr-067)),
[`MR-021`](../../../../harness/conventions.md#mr-021) (pin-gebundene
Verweise), [`MR-051`](../../../../harness/conventions.md#mr-051)
(`d-check:cite`-Spannen neu ankern),
[`MR-055`](../../../../harness/conventions.md#mr-055) (Symlink als Träger),
[`MR-069`](../../../../harness/conventions.md#mr-069) (`ignore-refs` als
deklarierte Gate-Senkung),
[`MR-070`](../../../../harness/conventions.md#mr-070) (Frozen-Klassen vor
mechanischer Ersetzung). Der neue Eintrag der Serie wird **`MR-071`** <!-- d-check:ignore (entsteht erst mit diesem Slice) -->.

**Berührte Spec-Stellen:** — (der Slice berührt keine Spec-Stelle).

**Verantwortlich:** pt9912 (Implementer-Rolle).

**Autor:** pt9912.

## 1. Ziel und Abgrenzung

**Ziel.** Den vendorten Baseline-Bestand von `v6.5.0` auf `v6.6.0` heben und
alle pin-gebundenen Verweise nachziehen — **mechanisch, ohne den Regel-Delta
zu beurteilen**. Was der Delta inhaltlich verlangt, entscheidet der
Folge-Slice; dieselbe Zerlegung wie beim Vorgänger-Paar
[slice-207](../done/slice-207-baseline-v650-bump.md) /
[slice-208](../done/slice-208-v650-regel-adoption.md).

**Der Anlass ist am Sensor bestätigt, nicht übernommen:**
`make baseline-freshness` meldet *„NEUER RELEASE verfügbar (Pin v6.5.0):
v6.6.0"* und zugleich, dass der **gepinnte** Tag upstream inhaltlich
unverändert ist (Bytes == vendored `SHA256SUMS`).

**Abgrenzung — vier Punkte, jeder mit Grund:**

1. **Kein Urteil über den Regel-Delta.** Was `v6.6.0` an Regeln ändert, wird
   **gemessen und gelistet**, aber nicht beantwortet — *ein Folge-Slice
   übernimmt es*, und zwar mit einer Antwort je Regel (übernommen · nicht
   anwendbar mit Begründung · abweichend als Adaption).
2. **Keine Template-Adoption.** Die neue `AGENTS.template.md` führt
   **weniger Tabellen** (Auftraggeber-Information) — genau der Grund, aus dem
   [slice-221](../next/slice-221-agents-md-tabellenzellen.md) nach `next/`
   zurückging. **Dieser Slice übernimmt die Form nicht**; er macht sie nur
   lesbar. *Es wäre ein anderer Vorgang*, und slice-221 wartet auf **ihn**,
   nicht auf diesen hier.
3. **slice-221 wird nicht angefasst.** Weder beansprucht noch nachgezogen —
   *Schicht-Abgrenzung*: Dieser Slice hebt einen Pin, er räumt keine Tabelle
   auf.
4. **Kein Aufräumen der `ignore-refs`-Einträge.** Die Hebung **fügt** einen
   hinzu (der `v6.5.0`-Baum verschwindet, eingefrorene Artefakte zitieren ihn
   weiter); dass die Liste damit wächst, ist als deklarierte Gate-Senkung
   bereits geführt ([`MR-069`](../../../../harness/conventions.md#mr-069)) —
   *Bestand bleibt bewusst stehen*.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** `.harness/baseline/v6.6.0/` ist materialisiert (`regelwerk/`,
      `templates/`, `SHA256SUMS`), der `v6.5.0`-Baum entfernt, §Baseline in
      [`harness/conventions.md`](../../../../harness/conventions.md) zeigt auf
      den neuen Tag, und **`MR-071`** <!-- d-check:ignore (entsteht erst mit diesem Slice) --> trägt die Hebung als nächster Eintrag der
      [`MR-011`](../../../../harness/conventions.md#mr-011)-Serie.
      `make baseline-verify` grün.
- [ ] **(2)** **Alle vier Spiegel-Klassen sind nachgezogen, nicht nur die
      grep-bare.**
      [`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
      (5×) nennt sie: Pfad-Verweise (gate-gedeckt) · Release-/Tree-**URLs** mit
      dem Tag · **Prosa-/Ellipsen-Pins** · der **zitierende Verweis**, dessen
      Wortlaut am neuen Ziel nicht mehr stehen muss
      ([`MR-051`](../../../../harness/conventions.md#mr-051)). Je Klasse steht
      im Slice, **wie** sie gesucht wurde — drei davon deckt kein Gate.
- [ ] **(3)** **Die Frozen-Klassen sind VOR der mechanischen Ersetzung
      aufgelistet** ([`MR-070`](../../../../harness/conventions.md#mr-070),
      Geltungsbereich trifft hier zu: eine mechanische Ersetzung über mehr als
      eine Datei) — über die **Eigenschaft**, nicht über Verzeichnisse. Was
      eingefroren ist, wird **nicht** retargetet, sondern über `ignore-refs`
      abgefangen.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

**Die Reihenfolge ist nicht beliebig, und der Grund steht in
[`MR-070`](../../../../harness/conventions.md#mr-070):** Die Frozen-Liste
entsteht **vor** der ersten Ersetzung, sonst ist sie eine Rechtfertigung
hinterher.

1. **Frozen-Klassen auflisten** — über die Eigenschaft *„würde ein
   korrigierter Wert verfälschen, was dieses Artefakt zu seinem Datum
   festgehalten hat?"*, nicht über Verzeichnisnamen.
2. **Delta messen, nicht beurteilen:** `diff -I '<!-- Quelle:'` gegen den
   alten Baum — die Herkunftszeile in Zeile 3 jeder Regelwerk-Datei trägt den
   Tag und meldete sonst *jede* Datei als geändert. Die Liste wandert in den
   Slice und ist die Eingabe des Folge-Slice.
3. **Re-vendoren:** `bash tools/harness/fetch-baseline-cache.sh v6.6.0`, alten
   Baum entfernen, `make baseline-verify`.
4. **Vier Spiegel-Klassen nachziehen**, je mit benannter Suchform.
5. **`d-check:cite`-Spannen neu ankern** — der Bump verschiebt Zeilennummern,
   und `citations` ist fail-closed im inneren Loop
   ([`MR-051`](../../../../harness/conventions.md#mr-051)).
6. `MR-071` <!-- d-check:ignore (entsteht erst mit diesem Slice) --> schreiben, §Baseline umstellen, `make gates`, Handoff.

## 4. Trigger

**Beanspruchung:** WIP-Limit frei (`in-progress/` ist leer, gemessen nach der
Rückführung von slice-221), `make baseline-freshness` meldet den neuen
Release, und der Auftraggeber hat den Vorgang am 2026-09-08 beauftragt.

**Rückführung nach `next/`** (`in-progress→next`): wenn der gemessene Delta
so groß ist, dass das **Nachziehen der Spiegel** selbst mehrere Sitzungen
braucht — dann trennt sich der Bump in Vendoring und Retargeting.

**Rückführung nach `open/`** (`in-progress→open`): wenn `v6.6.0` eine
**Struktur**-Änderung am vendorten Baum mitbringt (andere Verzeichnisnamen,
anderes Bundle-Layout), die die Pfad-Verweise nicht mechanisch abbildbar
macht. Dann ist vorher eine Entscheidung über die Verweis-Form fällig, wie sie
[`MR-023`](../../../../harness/conventions.md#mr-023) beim
Bundle-Layout-Wechsel gebraucht hat.

## 5. Closure-Trigger

DoD (1) bis (3) abgehakt, `make gates` und `make baseline-verify` grün mit
echter Ausgabe, unabhängiger Review durchgeführt und eingearbeitet,
Closure-Notiz geschrieben, Register fortgeschrieben, jedes Risiko aus §6 mit
einem der drei Ausgänge.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Drei der vier Spiegel-Klassen deckt kein Gate — und der Eintrag dazu steht
  bei 5×, ohne formgültigen Ausgang.** Ein grüner `make gates`-Lauf nach dem
  Bump sagt über Release-URLs, Prosa-Pins und zitierende Verweise **nichts**.
  Die einzige Gegenmaßnahme ist, je Klasse die Suchform aufzuschreiben und die
  Trefferliste anzusehen — DoD (2) verlangt genau das, und es ist eine
  Disziplin, kein Wächter. — **Ausgang:** \<offen\>
- **Der `citations`-Bruch ist die planmäßige Rot-Quelle, nicht ein Unfall.**
  Der Bump verschiebt Zeilenspannen; `citations` läuft fail-closed im inneren
  Loop und nimmt den `pre-commit`-Hook mit. Wer das nicht erwartet, hält es für
  einen Defekt und sucht an der falschen Stelle
  ([`MR-051`](../../../../harness/conventions.md#mr-051)). — **Ausgang:**
  \<offen\>
- **Die Über-Hebung ist so teuer wie die vergessene.** Eine
  Vergangenheits-Aussage („der Stand v6.5.0 führte …") mitzuheben macht sie
  falsch, und kein Gate meldet es — der Registereintrag nennt die Blindheit
  ausdrücklich **in beide Richtungen**. Die Frozen-Liste aus DoD (3) ist die
  Antwort darauf, und ob sie vollständig ist, bleibt Urteil. — **Ausgang:**
  \<offen\>

## 7. Closure-Notiz

\<wird vor dem `git mv` nach `done/` gefüllt\>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende.

**Die drei Vorprüfungen sind bei der Beanspruchung am 2026-09-08 bestätigt**
— Nachtlauf und Register unverändert gegenüber dem Anlege-Stand desselben
Tages; die Sub-Area-Wahl trägt der Block unten. Ursprünglich notiert war:
([`AGENTS.md`](../../../../AGENTS.md) §5). **Zwei sind schon gelaufen**, weil
sie den Zuschnitt tragen, und werden bei der Beanspruchung gegen den dann
gültigen Stand wiederholt:

- **Nachtlauf-Stand** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
  am 2026-09-08 beide grün. **Für diesen Slice ausnahmsweise nicht bezugslos:**
  `make baseline-freshness` ist Teil desselben Nachtlaufs, und sein Exit 3 ist
  der Auslöser dieses Vorgangs.
- **Register** (40 Verzeichnisse): **drei** Einträge sind einschlägig.
  [`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
  (5×) ist **der tragende** und formt DoD (2) — er nennt die vier Klassen und
  benennt, dass drei davon gate-blind sind, **in beide Richtungen**. Sein Stand
  ist *kein formgültiger Ausgang*, geführt bei
  [`registerzeile-ohne-ausgang-nach-schwelle`](../observations/BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle/observation.md);
  dieser Slice ändert daran nichts und erfindet auch keinen.
  [`mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
  (3×, verkörpert als [`MR-070`](../../../../harness/conventions.md#mr-070))
  ist hier **anwendbar** — anders als in slice-221: Der Geltungsbereich nennt
  *„jede mechanische Ersetzung über mehr als eine Datei"*, und genau das ist
  das Retargeting. DoD (3) löst die Pflicht ein.
  [`form-vom-nachbarn-statt-von-der-vorlage`](../observations/BEO-ALL/form-vom-nachbarn-statt-von-der-vorlage/observation.md)
  (1×) ist der Grund, warum slice-221 wartet — und die Mahnung, in diesem Slice
  die **Form** der neuen Vorlage nicht nebenbei zu übernehmen (Abgrenzung 2).

**Modus-Begründungsblock.** Alle berührten Sub-Areas GF — ein Block genügt.

### Sub-Area: `*`

- **Modus:** GF (Repo-Default).
- **Konventions-Dichte:** **sehr hoch.** Die Pin-Serie ist zwölfmal gelaufen
  ([`MR-011`](../../../../harness/conventions.md#mr-011) bis
  [`MR-067`](../../../../harness/conventions.md#mr-067)), das Vorgehen ist in
  fünf `MR`-Einträgen verankert, und der Vorgänger-Slice liegt als Vorlage vor.
- **Phase-Reife:** Phase 5.
- **Evidenz-/Diskrepanz-Risiko:** **niedrig für den vendorten Baum**
  (`baseline-verify` prüft ihn hart), **mittel für die Spiegel** — drei der
  vier Klassen sind gate-blind, und der Registereintrag dazu steht bei 5×.
- **Reconciliation-Aufwand:** keiner (GF).
