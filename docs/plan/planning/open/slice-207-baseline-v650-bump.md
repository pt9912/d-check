# Slice slice-207: Baseline-Pin auf `v6.5.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden wäre. Die Belege sind je Slice (`baseline-verify`
grün, `gates` grün); ein Wellen-Trigger schriebe sie ab. Präzedenz: die
Pin-Hebung über **vier** Tags lief ebenso wellenlos.

**Bezug:** [`MR-011`](../../../../harness/conventions.md#mr-011)-Kette (Pin auf
Release-Tag), [`MR-023`](../../../../harness/conventions.md#mr-023)
(Bundle-Layout), [`MR-021`](../../../../harness/conventions.md#mr-021)
(pin-gebundene Verweise), [`MR-051`](../../../../harness/conventions.md#mr-051)
(`d-check:cite`-Spannen beim Bump neu ankern),
[`MR-055`](../../../../harness/conventions.md#mr-055) (Symlink als Träger).

**Berührte Spec-Stellen:** — *(keine; der Slice bewegt den Baseline-Pin und
die pin-gebundenen Verweise, keine Anforderung und keine Sicht)*

**Verantwortlich:** — · **Autor:** pt9912. **Datum:** 2026-09-06.

---

## 1. Ziel

Den vendorten Baseline-Bestand von `v6.3.1` auf `v6.5.0` heben und alle
pin-gebundenen Verweise nachziehen — **mechanisch, ohne den Regel-Delta zu
beurteilen**. Was der Delta inhaltlich verlangt, entscheidet der Folge-Slice.

## 2. Vorgehen

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`.harness/baseline/`](../../../../.harness/baseline/) | neu + entfernt | `v6.5.0` materialisieren, `v6.3.1` entfernen — `fetch-baseline-cache.sh` trägt beides |
| pin-gebundene Verweise | update | jeder `.harness/baseline/<tag>/`-Pfad in lebenden Dokumenten, plus die Release-/Tree-URLs mit dem Tag |
| `d-check:cite`-Direktiven | update | die Zeilen-Spannen verschieben sich mit dem Bump; `citations` ist fail-closed und läuft im inneren Loop |
| neuer Eintrag unter [`harness/conventions/`](../../../../harness/conventions/) | neu | Nachtrag zur Pin-Serie, mit dem gemessenen Delta-Umfang; die Kennung wird beim Schreiben vergeben, nicht hier |

**Der Sprung überspringt zwei Releases** (`v6.4.0`, `v6.5.0`) — und das ist der
Grund für die erste Zeile der DoD. Ein direkter Byte-Vergleich über zwei
Versionen produziert Rauschen (Versionsnummern, Datumszeilen, verschobene
Zeilen), das den echten Regel-Delta verdeckt; gemessen wird deshalb mit
`diff -I`, die Lehre aus der vorigen Pin-Hebung.

**Was `v6.4.0` bringt, ist bekannt und angekündigt:** die Umsetzung der beiden
angenommenen Bitten unseres ausgehenden CR
([Antwort](../../cr/2026-09-06-antwort-ai-harness-course-slice-formluecken.md))
— der Ausschluss-Abschnitt als zweite Hälfte von §1, §8 mit unbedingtem Kopf
und bedingtem Rumpf, beide Prosa-Pflaster entfernt. **Was `v6.5.0` bringt, ist
unbekannt**; das misst dieser Slice.

## 3. Ausdrücklich NICHT in diesem Slice

- **Jede Regel-Adoption.** Ob und wie ein Delta-Punkt übernommen wird, ist ein
  Urteil je Regel und gehört in
  [slice-208](../open/slice-208-v650-regel-adoption.md). Dieser Slice **misst** den
  Delta und **hebt den Pin**; er entscheidet nichts.
- **Die Auflösung der Slice-Haus-Form.** Sie folgt aus dem `v6.4.0`-Delta und
  ist der Kern des Folge-Slice — inklusive der Gate-Regeln, die heute auf
  Haus-Form-Titel keilen.
- **Jede Änderung an Anforderung, Sicht oder ADR.** Ein Pin bewegt keine
  Spec-Stelle.

## 4. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** Der vendorte Bestand steht auf `v6.5.0`, `make baseline-verify`
      und `make baseline-probe` sind grün, und der **Delta gegenüber `v6.3.1`
      ist mit `diff -I` gemessen** und als Liste im Slice festgehalten —
      getrennt nach *Regelwerk*, *Templates* und *reines Rauschen*.
- [ ] **(2)** Alle pin-gebundenen Verweise zeigen auf `v6.5.0`
      ([`MR-021`](../../../../harness/conventions.md#mr-021)), die
      `d-check:cite`-Spannen sind neu geankert
      ([`MR-051`](../../../../harness/conventions.md#mr-051)), und die Aliase
      unter `.claude/rules/` lösen auf
      ([`MR-055`](../../../../harness/conventions.md#mr-055)).
- [ ] **(3)** Ein neuer Konventions-Eintrag trägt die Hebung als Nachtrag zur
      [`MR-011`](../../../../harness/conventions.md#mr-011)-Kette, mit dem
      **gemessenen** Delta-Umfang statt einer Schätzung.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §5 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 5. Abnahme-Punkte / Risiken

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Zwei Releases in einem Sprung — der Delta ist größer als bei jeder
  bisherigen Hebung mit bekanntem Inhalt.** Zeigt die Messung, dass `v6.5.0`
  eine eigene Regel-Änderung trägt, ist der Folge-Slice **vor** seiner
  Beanspruchung neu zu schneiden. Das ist kein Fehler, sondern der Grund,
  warum die Messung in diesem Slice liegt und nicht im nächsten. —
  **Ausgang:** \<offen\>
- **Die `d-check:cite`-Spannen sind die planmäßige Rot-Quelle**
  ([`MR-051`](../../../../harness/conventions.md#mr-051)). `citations` ist
  fail-closed und läuft im inneren Loop: eine nicht neu geankerte Direktive
  nimmt den `pre-commit`-Hook mit. Über **zwei** Versionen verschieben sich
  mehr Zeilen als über eine. — **Ausgang:** \<offen\>
- **Ein Verweis, den kein Gate hält**
  ([`BEO-ALL/pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)):
  Release-/Tree-URLs, Prosa-Pins und der **zitierende** Verweis, dessen Zitat
  am neuen Ziel nicht mehr existiert. Gate-blind in beide Richtungen —
  vergessene Hebung wie Über-Hebung. — **Ausgang:** \<offen\>

## 6. Trigger

**Start** (`open` → `in-progress`): `make baseline-freshness` meldet
`v6.4.0` und `v6.5.0` als neuere Releases (gelesen 2026-09-06); der Content am
gepinnten Tag ist unverändert. WIP-Limit frei.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Zeigt der gemessene Delta, dass allein das
  Neu-Ankern der Zitat-Spannen eine eigene Sitzung füllt, wird es ein eigener
  Slice vor der Pin-Hebung.
- `in-progress` → `open` (blockiert): Ist der Upstream-Bestand am neuen Tag
  nicht integer (`SHA256SUMS` passt nicht), ruht der Slice bis zur Klärung mit
  der Baseline — still weiter zu vendoren wäre der Verlust der Integritäts-Zusage.

**Closure-Trigger.** Zwei beobachtbare Kriterien und ein Lerneintrag: (a)
`make baseline-verify`, `make baseline-probe` und `make gates` sind grün; (b)
`make baseline-freshness` meldet den Pin als aktuell.

## 7. Vorgelagert (vor der Modus-Begründung)

\<entsteht spätestens bei der Beanspruchung — ein Plan in `open/` trägt die drei
Vorprüfungen noch nicht\>

## 8. Sub-Area-Modus-Begründung

\<entsteht mit den Vorprüfungen bei der Beanspruchung\>

## 9. Closure-Notiz (nach `done/`)

\<wird vor dem `git mv` nach `done/` gefüllt\>
