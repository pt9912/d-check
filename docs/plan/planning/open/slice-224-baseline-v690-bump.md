# Slice slice-224: Baseline-Pin auf `v6.9.0` — mechanisch, ohne Urteil über den Delta

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`MR-011`](../../../../harness/conventions.md#mr-011)-Pin-Serie
(Vorgänger-Eintrag: [`MR-071`](../../../../harness/conventions.md#mr-071)),
[`MR-021`](../../../../harness/conventions.md#mr-021) (pin-gebundene
Verweise), [`MR-051`](../../../../harness/conventions.md#mr-051)
(`d-check:cite`-Spannen neu ankern),
[`MR-055`](../../../../harness/conventions.md#mr-055) (Symlink als Träger),
[`MR-069`](../../../../harness/conventions.md#mr-069) (`ignore-refs` als
deklarierte Gate-Senkung),
[`MR-070`](../../../../harness/conventions.md#mr-070) (Frozen-Klassen vor
mechanischer Ersetzung). Der neue Eintrag der Serie wird **`MR-072`** <!-- d-check:ignore (entsteht erst mit diesem Slice) -->.

**Berührte Spec-Stellen:** — (der Slice berührt keine Spec-Stelle).

**Verantwortlich:** — (wird bei der Beanspruchung gesetzt).

**Autor:** pt9912. **Datum:** 2026-09-16.

---

## 1. Ziel und Abgrenzung

**Ziel.** Den vendorten Baseline-Bestand von `v6.6.0` auf `v6.9.0` heben und
alle pin-gebundenen Verweise nachziehen — **mechanisch, ohne den Regel-Delta
zu beurteilen**. Was der Delta inhaltlich verlangt, entscheidet ein
Folge-Slice; dieselbe Zerlegung wie beim Vorgänger-Paar
[slice-207](../done/slice-207-baseline-v650-bump.md) /
[slice-208](../done/slice-208-v650-regel-adoption.md).

**Der Anlass ist am Sensor bestätigt, nicht übernommen:**
`make baseline-freshness` meldet *„NEUER RELEASE verfügbar (Pin v6.6.0):
v6.7.0, v6.7.1, v6.7.2, v6.8.0, v6.9.0"* und zugleich, dass der **gepinnte**
Tag upstream inhaltlich unverändert ist (Bytes == vendored `SHA256SUMS`).
`make nightly-state` bestätigt denselben Befund am jüngsten Lauf von
`upstream-drift.yml` (2026-09-16, planmäßig rot — Fremd-Release, keine
Störung).

**Fünf Releases liegen zwischen dem Pin und upstream, nicht eines.** Das ist
kein neuer Fall: Die [`MR-011`](../../../../harness/conventions.md#mr-011)-Kette hat wiederholt mehrere Zwischen-Tags in
einem Sprung übersprungen (`v5.7.0`→`v5.9.0`, `v5.9.0`→`v5.11.0`,
`v5.15.0`→`v5.18.0`, `v6.0.0`→`v6.3.1`) — ein Bump zielt auf den **jeweils
aktuellen** Tag, nicht auf jeden dazwischen. Dasselbe gilt hier: Ziel ist
`v6.9.0`, kein Zwischenstopp.

**Abgrenzung — vier Punkte, jeder mit Grund:**

1. **Kein Urteil über den Regel-Delta.** Was `v6.9.0` gegenüber `v6.6.0` an
   Regeln ändert, wird **gemessen und gelistet**, aber nicht beantwortet —
   *ein Folge-Slice übernimmt es*, mit einer Antwort je Regel (übernommen ·
   nicht anwendbar mit Begründung · abweichend als Adaption).
2. **Keine Template-Adoption.** `AGENTS.md` §4 trägt bereits die Ziel-Form aus
   der `v6.6.0`-Vorlage (§4 als Zeiger auf `harness/README.md`, slice-222).
   Ändert `v6.9.0`s `AGENTS.template.md` diese Form **erneut**, ist das ein
   Fund dieses Slice (§3, Delta-Messung), aber seine Übernahme bleibt dem
   Folge-Slice vorbehalten — *es wäre ein anderer Vorgang*, und die Lehre aus
   slice-222 (§6 dieses Plans) ist gerade, diesen Punkt **vorab** zu
   entscheiden statt ihn während der Implementierung neu aufzurollen. Trifft
   der Fall doch ein, ist die in slice-222 gewählte Antwort (Bump + Adoption
   im selben Slice, per ausdrücklicher Auftraggeber-Weisung) die Präzedenz,
   an der zu messen ist — keine stillschweigende Wiederholung ohne dieselbe
   Weisung.
3. **slice-221 wird nicht angefasst.** Weder beansprucht noch nachgezogen —
   *Schicht-Abgrenzung*: Dieser Slice hebt einen Pin, er räumt keine Tabelle
   auf. Es liegt weiter in `next/` und wird **nach** diesem Bump gegen die
   dann geltende Vorlage neu zugeschnitten (sein eigener Trigger, §4 dort).
4. **Kein Aufräumen der `ignore-refs`-Einträge.** Wächst die Liste durch
   diesen Bump (neuer Eintrag, weil der `v6.6.0`-Baum verschwindet und
   eingefrorene Artefakte ihn weiter zitieren), ist das als deklarierte
   Gate-Senkung bereits geführt ([`MR-069`](../../../../harness/conventions.md#mr-069))
   — *Bestand bleibt bewusst stehen*.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** `.harness/baseline/v6.9.0/` ist materialisiert (`regelwerk/`,
      `templates/`, `SHA256SUMS`), der `v6.6.0`-Baum entfernt, §Baseline in
      [`harness/conventions.md`](../../../../harness/conventions.md) zeigt auf
      den neuen Tag, und **`MR-072`** <!-- d-check:ignore (entsteht erst mit diesem Slice) --> trägt die Hebung als nächster Eintrag der
      [`MR-011`](../../../../harness/conventions.md#mr-011)-Serie.
      `make baseline-verify` grün.
- [ ] **(2)** **Alle vier Spiegel-Klassen sind nachgezogen, nicht nur die
      grep-bare** —
      [`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
      (6×, weiterhin ohne formgültigen Ausgang) nennt sie: Pfad-Verweise
      (gate-gedeckt) · Release-/Tree-**URLs** mit dem Tag ·
      **Prosa-/Ellipsen-Pins** · der **zitierende Verweis**, dessen Wortlaut
      am neuen Ziel nicht mehr stehen muss
      ([`MR-051`](../../../../harness/conventions.md#mr-051)). Je Klasse steht
      im Slice, **wie** sie gesucht wurde — Suchform **nach der Version**, mit
      Gruppierung nach Präfix (Lehre aus slice-222, nicht nach dem Pfad
      suchen).
- [ ] **(3)** **Die Frozen-Klassen sind VOR der mechanischen Ersetzung
      aufgelistet** ([`MR-070`](../../../../harness/conventions.md#mr-070),
      Geltungsbereich trifft hier zu: eine mechanische Ersetzung über mehr als
      eine Datei) — über die **Eigenschaft**, nicht über Verzeichnisse. Was
      eingefroren ist, wird **nicht** retargetet, sondern über `ignore-refs`
      abgefangen. **Lebende** Artefakte mit einer Vergangenheits-Aussage über
      einen früheren Pin (wie [`MR-067`](../../../../harness/conventions.md#mr-067)
      in slice-222) fallen nicht automatisch unter „Bestand" — jeder Eintrag
      der [`MR-011`](../../../../harness/conventions.md#mr-011)-Kette wird einzeln gegen seinen `Geltungsbereich` gelesen,
      nicht pauschal als „liegt in `conventions/`, also lebend, also
      retargeten".
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
   festgehalten hat?"*, nicht über Verzeichnisnamen. Jede [`MR-011`](../../../../harness/conventions.md#mr-011)-Kettenglied
   einzeln gegen seinen `Geltungsbereich` lesen (Lehre aus slice-222s
   Über-Hebung von [`MR-067`](../../../../harness/conventions.md#mr-067)).
2. **Neuen Baum materialisieren** (`fetch-baseline-cache.sh v6.9.0`), **ohne**
   den alten zu löschen — beide koexistieren kurz, damit Schritt 3 diffen
   kann.
3. **Delta messen, nicht beurteilen:** `diff -I '<!-- Quelle:'` gegen den
   alten Baum (die Herkunftszeile in Zeile 3 jeder Regelwerk-Datei trägt den
   Tag und meldete sonst *jede* Datei als geändert). Die Liste wandert in den
   Slice und ist die Eingabe des Folge-Slice. **Hier fällt die Entscheidung**,
   ob Abgrenzung 2 (keine Template-Adoption) hält oder — wie in slice-222 —
   eine erzwungene Ausnahme braucht; ohne Weisung des Auftraggebers bleibt es
   bei Abgrenzung 2.
4. **Alten Baum entfernen**, `make baseline-verify`.
5. **Vier Spiegel-Klassen nachziehen**, je mit benannter Suchform — Version
   statt Pfad, Gruppierung nach Präfix.
6. **`d-check:cite`-Spannen neu ankern** — der Bump verschiebt Zeilennummern,
   und `citations` ist fail-closed im inneren Loop
   ([`MR-051`](../../../../harness/conventions.md#mr-051)).
7. `MR-072` <!-- d-check:ignore (entsteht erst mit diesem Slice) --> schreiben, §Baseline umstellen, `make gates`, Handoff.

## 4. Trigger

**Start** (`next` → `in-progress`): WIP-Limit frei (`in-progress/` ist leer,
Ruhe-Marker steht), `make baseline-freshness` meldet den neuen Release, und
der Auftraggeber hat den Vorgang am 2026-09-16 beauftragt (Entscheidung
zwischen drei Optionen — voller Bump zuerst statt Weiterarbeit an slice-221
oder nur Sichtung des Deltas).

**Rückführung nach `next/`** (`in-progress→next`): wenn der gemessene Delta
so groß ist, dass das **Nachziehen der Spiegel** selbst mehrere Sitzungen
braucht — dann trennt sich der Bump in Vendoring und Retargeting.

**Rückführung nach `open/`** (`in-progress→open`): wenn `v6.9.0` eine
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

- **Drei der vier Spiegel-Klassen deckt kein Gate** — dieselbe Lücke wie bei
  jedem Vorgänger-Bump, zuletzt gemessen bei 6×
  ([`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md),
  weiterhin ohne formgültigen Ausgang: die mechanische Form `versions.patterns`
  existiert seit slice-122, ist aber bewusst nicht scharfgeschaltet). Ein
  grüner `make gates`-Lauf nach dem Bump sagt über Release-URLs, Prosa-Pins
  und zitierende Verweise **nichts**. — **Ausgang:** \<offen\>
- **Die Über-Hebungs-Falle bei lebenden `MR`-Einträgen ist bekannt, aber die
  Frozen-Liste fängt sie nicht.** slice-222 hat gezeigt: Eine mechanische
  Ersetzung über den `Geltungsbereich`-Text kann einen Eintrag treffen, der
  *„lebend"* ist (liegt in `harness/conventions/`), aber inhaltlich eine
  **Vergangenheits-Aussage** über einen früheren Pin trägt
  ([`mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md),
  bereits als [`MR-070`](../../../../harness/conventions.md#mr-070)
  verkörpert — die Grenze der Regel ist benannt, nicht behoben). Gegenmittel
  dieses Slice: jedes [`MR-011`](../../../../harness/conventions.md#mr-011)-Kettenglied einzeln lesen statt pauschal
  „`conventions/` ⇒ retargeten". — **Ausgang:** \<offen\>
- **Der `citations`-Bruch ist die planmäßige Rot-Quelle, nicht ein Unfall.**
  Der Bump verschiebt Zeilenspannen; `citations` läuft fail-closed im inneren
  Loop und nimmt den `pre-commit`-Hook mit
  ([`MR-051`](../../../../harness/conventions.md#mr-051)). — **Ausgang:**
  \<offen\>
- **Ändert `v6.9.0`s `AGENTS.template.md` §4 (oder eine andere Sektion mit
  Pointer-Charakter) erneut die Form**, entsteht dieselbe Zwickmühle wie bei
  slice-222: Bump ohne Adoption ließe `AGENTS.md` einer Vorlage widersprechen,
  die jeder Lauf lädt. Abgrenzung 2 hält dagegen fest, dass die Übernahme
  einer erneuten Weisung bedarf, keiner stillen Wiederholung. — **Ausgang:**
  \<offen\>

## 7. Closure-Notiz

\<wird vor dem `git mv` nach `done/` gefüllt\>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende. Dieses Repo führt **drei** Prüfungen — die
zwei kanonischen und, als Adaption, den Nachtlauf-Stand
([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.6.0/regelwerk/modul-05-planning-harness.md:268-269 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area: `*` (Repo-Default). Der Slice hebt einen vendorten
Bestand und zieht pin-gebundene Verweise nach; `tools/harness/` ist nicht
berührt (das Werkzeug selbst ändert sich nicht, nur sein Ziel-Tag).

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.6.0/regelwerk/modul-05-planning-harness.md:274-274 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, **40** Verzeichnisse). **Vier
Einträge sind einschlägig, alle vier aus dem Vorgänger-Bump bekannt:**

- [`pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
  (6×, **kein formgültiger Ausgang**) — der tragende Eintrag für DoD (2).
- [`mechanical-id-rewrite-misses-frozen-classes`](../observations/BEO-ALL/mechanical-id-rewrite-misses-frozen-classes/observation.md)
  (4×, verkörpert als [`MR-070`](../../../../harness/conventions.md#mr-070))
  — trägt DoD (3); die Grenze der Regel (Frozen-Liste über Eigenschaften kann
  eine im *Inhalt* versteckte Vergangenheits-Aussage nicht fangen) ist
  benannt, nicht behoben — §6 trägt sie deshalb erneut als Risiko.
- [`form-vom-nachbarn-statt-von-der-vorlage`](../observations/BEO-ALL/form-vom-nachbarn-statt-von-der-vorlage/observation.md)
  (2×, offen unter der Schwelle) — Mahnung, die Form aus der **Vorlage**
  abzuleiten, nicht aus dem Bestand; hier vor allem für die
  Abgrenzung-2-Frage in §3 Schritt 3 relevant.
- [`liefer-punkt-in-fremdem-commit`](../observations/BEO-ALL/liefer-punkt-in-fremdem-commit/observation.md)
  (2×, offen unter der Schwelle) — Mahnung, jeden Commit gegen `git show`
  statt aus der Erinnerung zu beschreiben, besonders bei einem
  Kontext-Wechsel mitten in der Arbeit.

**Keiner der vier erreicht mit diesem Slice die Schwelle erstmalig.**

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-16 gelesen: `upstream-drift.yml` **planmäßig
rot** (2026-09-16T05:33:02Z) — genau der Fremd-Release-Befund, der dieses
Slice auslöst, keine unerwartete Störung. `image-scan.yml` grün
(2026-09-16T08:37:54Z), ohne Bezug zu diesem Slice.

**Modus-Begründungsblock.** Alle berührten Sub-Areas GF — ein Block genügt.

### Sub-Area: `*`

- **Modus:** GF (Repo-Default).
- **Konventions-Dichte:** **sehr hoch.** Die Pin-Serie ist dreizehnmal
  gelaufen ([`MR-011`](../../../../harness/conventions.md#mr-011) bis
  [`MR-071`](../../../../harness/conventions.md#mr-071)), das Vorgehen ist in
  fünf `MR`-Einträgen verankert, und der Vorgänger-Slice
  ([slice-222](../done/slice-222-baseline-v660-bump.md)) liegt als Vorlage
  samt sieben dokumentierten Lerneinträgen vor.
- **Phase-Reife:** Phase 5.
- **Evidenz-/Diskrepanz-Risiko:** **niedrig für den vendorten Baum**
  (`baseline-verify` prüft ihn hart), **mittel für die Spiegel** — drei der
  vier Klassen sind gate-blind, und der Registereintrag dazu steht
  unverändert bei 6× ohne Ausgang.
- **Reconciliation-Aufwand:** keiner (GF).
