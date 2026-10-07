# slice-252: `targets` prüft opt-in die Disjunktheit der Autoritäts-Dateien

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [`DC-FA-TGT-001`](../../../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
(Scope — erweitert, kein neues Kürzel: Einzelmodul-Frage nach
[ADR-0044](../../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md));
Anlass ist der eingetretene Re-Evaluierungs-Trigger von
[ADR-0100](../../../adr/0100-targets-authority-liste.md), gemeldet im
[Hinweis der Baseline](../../../cr/2026-10-07-hinweis-eingehend-ai-harness-course-disjunktheit.md);
Beobachtung
[`BEO-ALL/disjunktheit-geteilter-gate-index-ungeprueft`](../../observations/BEO-ALL/disjunktheit-geteilter-gate-index-ungeprueft/state.md);
[`DC-QA-02`](../../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(ohne den Schalter byte-identisch).

**Berührte Spec-Stellen:** [§DC-FA-TGT-001.a](../../../../../spec/spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets)
(neuer Schritt), [`SPEC-005`](../../../../../spec/spezifikation.md#spec-005--d-checkyml)
(Schema-Zeile), [§4 Grund-Codes](../../../../../spec/spezifikation.md#4-grund--und-fehler-codes)
(neue Zeile).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.
**Ziel:** Mit `targets.authority-disjoint: true` meldet `targets` ein Target,
das in mehr als einer `targets.authority`-Datei von einer Tabellenzeile in
der ersten Zelle geführt wird, mit dem neuen Grund-Code `gate-declared-twice`
— an der führenden Zeile jeder weiteren Datei, die Meldung nennt die Datei
der ersten Nennung. Ohne den Schalter ist der Befundsatz byte-identisch.
*(Plan-Änderung nach R1: zuerst zählte jede Zelle, und die Prüfung hing an
`targets.makefiles`; beides widersprach der Baseline-Lesart bzw. ließ einen
eingeschalteten Schalter still.)*
Schalter ist der Befundsatz byte-identisch.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Doppelungen innerhalb einer Datei** — andere Frage: die Regel spricht von
  zwei *Teilen* des Index; eine Tabelle, die ein Target zweimal führt, ist ein
  Struktur-Mangel einer Datei.
- **`doc-tables`** — sind keine Teile des Gate-Index, sondern Orte von
  Behauptungen (Richtung 1).
- **Default-an** — die Disjunktheit ist eine Regel ab Baseline `v6.16.0`;
  Adopter älterer Baselines sollen nicht rot werden.
- **Baseline-Hebung dieses Repos** auf `v6.16.0` — anderer Vorgang
  (Nachtlauf-Meldung, eigener Pin-Slice).
- **make-Target für eine Vorher-/Nachher-Black-Box-Probe** — Entscheidung
  des Auftraggebers offen; dieser Slice fährt die Probe von Hand vor dem
  Review.
- **Benutzerhandbuch, README, CHANGELOG** — Release-Prep
  (`AGENTS.md` §5 Regel 17).

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Liefer-Punkte (3):

- [x] Kern: Disjunktheits-Prüfung mit `gate-declared-twice` (Fundstelle =
      Tabellenzeile in der späteren Datei, Konfigurations-Reihenfolge;
      `exempt-targets` wirkt nicht; mit einer Datei wirkungslos); Tests rot
      ohne die Änderung aus dem richtigen Grund (Bewusstes Brechen).
- [x] Config-Rand und Lexik: Schlüssel `targets.authority-disjoint`
      (Bool, strikt), `--print-config`-Gerüst, Grund-Code in `AllReasons`
      und im `--doctor`-Klartext.
- [x] Lastenheft (Erweiterung, Bump + Historie nach
      [`MR-032`](../../../../../harness/conventions.md#mr-032)), Spezifikation
      (Schritt, Schema, §4); ADR; Antwort im Hinweis-Dokument; eigene
      Black-Box-Probe gegen ein Vorher-Image vor dem Review (ohne Schalter
      byte-identisch), Ergebnis im Plan notiert.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor —
      Rollenwechsel nach Schritt 8, kein Self-Review.
- [x] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben
      (Ausgang des Eintrags zur Disjunktheit); jedes Risiko aus §6 mit
      Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | Anforderung erweitert, Akzeptanzkriterien, Out-of-Scope, Bump + Historie |
| `spec/spezifikation.md` | update | neuer Schritt, Schema-Zeile, §4-Zeile (nächste freie Struktur-ID), Historie |
| `internal/hexagon/core/model/config.go` | update | `AuthorityDisjoint` |
| `internal/hexagon/core/rules/targets.go` | update | Prüfung, Grund-Code |
| `internal/hexagon/core/app/diagnose.go` | update | `AllReasons`, Klartext |
| `internal/adapter/driven/configyaml/configyaml.go` | update | Schlüssel |
| `internal/adapter/driving/cli/config_template.go` | update | Gerüst |
| `internal/hexagon/core/rules/targets_test.go`, `configyaml_test.go` | update | Happy/Negative/Boundary |
| `docs/plan/adr/0101-…`, `docs/plan/adr/README.md` | neu / update | Entscheidung, Bezug auf die Vorgänger-ADR |

**Spiegel vor dem Editieren** ([`MR-025`](../../../../../harness/conventions.md#mr-025)):
Lastenheft (Beschreibung, Out-of-Scope, der Satz „eine Doppelnennung ist kein
Befund"), Spezifikation (Schritt 5 sagt dasselbe), die Vorgänger-ADR, Entscheidung 2
(Geschichte-Anhang, Körper frozen), CR-Antwort zur authority-Liste (sie nennt
den Folgeschritt), `--print-config`, `--doctor`, `AllReasons`/§4-Lockstep,
Handbuch §6-Tabelle (Release-Prep).

## 4. Trigger

**Start** (`next` → `in-progress`): Implementer übernimmt; `in-progress/`
leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): verlangt die Fundstellen-Regel eine
  Änderung an der Tabellen-Extraktion über `targets` hinaus.
- `in-progress` → `open` (blockiert): keiner bekannt.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- Die Aussage „eine Doppelnennung ist kein Befund" steht in Lastenheft,
  Spezifikation, Vorgänger-ADR und CR-Antwort; bleibt eine davon stehen, sagt das
  Repo zwei Dinge. — **Ausgang:** entfallen — alle Stellen auf „für `gate-undocumented` bzw. ohne Schalter" eingeschränkt (Lastenheft, Spezifikation Schritt 5, Anhang der Vorgänger-ADR, Nachtrag der CR-Antwort); Handbuch und READMEs zieht die Release-Prep nach.
- Die Black-Box-Probe von Hand kann wieder nur die Fälle prüfen, an die der
  Autor denkt. — **Ausgang:** weiter offen — eingetreten in diesem Vorgang (die Probe bestätigte die Byte-Identität, fand aber weder die Baseline-Lesart noch den Panic bei einer leeren Tabellenzeile); als Beleg ins Register zu `BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Baseline-Quelle wurde vor dem Schnitt gelesen
  (`v6.16.0`, Bedingung *Disjunkt*), der Re-Evaluierungs-Trigger der
  Vorgänger-ADR war damit belegt. Die eigene Black-Box-Probe vor dem Review
  bestätigte die Byte-Identität ohne Schalter; nach allen Korrekturen erneut
  gefahren: **78 Vergleiche** gegen `59d21c07^` (drei Durchläufe, vier
  Fixtures mit einer `|`-Zeile, je sechs Ausgabeformen, dazu sechs
  Selbstläufe über das Repo), stdout und stderr getrennt — byte-identisch. Ein
  erster Durchlauf mit zusammengeführten Streams zeigte zwei Abweichungen,
  die sich als Verschränkung der Streams herausstellten.
- **Was ging anders als geplant:** Zweimal ein HIGH. R1 fand, dass die erste
  Fassung jede Zelle zählte, während die Baseline die **führende** Zeile
  meint, und dass die Prüfung ohne `targets.makefiles` nicht lief. Die
  Behebung setzte im geteilten Extraktor an und brachte eine Regression mit:
  eine Tabellenzeile aus nur `|` ließ jeden `targets`-Lauf mit Panic
  abbrechen, auch ohne Schalter (Verifier V1). Zwei Plan-Stellen (Ziel, DoD)
  tragen Änderungs-Vermerke.
- **Steering-Loop-Eintrag:** Die Klasse
  [`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/state.md)
  traf den vierten Slice in Folge an diesem Modul. Gemessen an diesem Vorgang:
  Die Hand-Probe vor dem Review bestätigte, was sie prüfte, und fand nichts,
  woran der Autor nicht gedacht hatte. Der in slice-251 benannte
  Sensor-Kandidat (Vorher-/Nachher-Black-Box als make-Target) hätte die
  Byte-Identität mechanisch gehalten, die **Lesart** einer Baseline-Regel
  aber ebenso wenig gefunden — dafür trug der fremde Leser. Die Entscheidung
  über den Sensor bleibt beim Auftraggeber.
- **Beobachtungs-Register (`../../observations/`):** Evidence `slice-252` unter
  `grenzen-liste-wird-als-vollstaendig-gelesen` und
  [`BEO-ALL/fix-aendert-unbenannten-zweiten-pfad-mit`](../../observations/BEO-ALL/fix-aendert-unbenannten-zweiten-pfad-mit/state.md)
  (jetzt 2×);
  [`BEO-ALL/disjunktheit-geteilter-gate-index-ungeprueft`](../../observations/BEO-ALL/disjunktheit-geteilter-gate-index-ungeprueft/state.md)
  bekommt den Ausgang *verkörpert*.
- **Folge-Slices:** keine geschnitten. Benannt: die Baseline-Hebung dieses
  Repos auf `v6.16.0` (Nachtlauf), der Sensor-Kandidat (Auftraggeber).
- **Risiken aus §6:** Risiko 1 (widersprüchliche Aussage zur Doppelnennung):
  entfallen — alle Stellen eingeschränkt, Handbuch/READMEs folgen in der
  Release-Prep. Risiko 2 (Hand-Probe sieht nur bekannte Fälle): weiter offen
  — eingetreten, Beleg im Register. Trigger-Audit: kein Carveout, kein
  bootstrap-aware Gate; die Vorgänger-ADR löste ihren Re-Evaluierungs-Trigger
  über die neue ADR ein (Anhang dort), die neue ist `Accepted`; keine Hard
  Rule mit eingetretenem Auflösungs-Trigger. Nachtlauf-Stand
  ([`MR-053`](../../../../../harness/conventions.md#mr-053)): wie in §8.
  Antwort an die Baseline im
  [Hinweis-Dokument](../../../cr/2026-10-07-hinweis-eingehend-ai-harness-course-disjunktheit.md).
- **Drei Paarungen:** (a) Anker — kein `liegt in`-Feld; (b) Folge-Slice —
  keine; (c) Register — alle zitierten Beobachtungen existieren und tragen
  Belege.

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: das Produkt
samt Spec-Stratum (`*`, Kürzel `ALL`); bereits deklariert.

<!-- d-check:cite .harness/baseline/v6.13.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** drei Einträge treffen den
Gegenstand —
[`BEO-ALL/disjunktheit-geteilter-gate-index-ungeprueft`](../../observations/BEO-ALL/disjunktheit-geteilter-gate-index-ungeprueft/state.md)
(der Gegenstand selbst; bekommt seinen Ausgang),
[`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/state.md)
(dreimal in Folge an diesem Modul — deshalb die eigene Black-Box-Probe vor
dem Review) und
[`BEO-ALL/commit-message-overclaims-work`](../../observations/BEO-ALL/commit-message-overclaims-work/state.md)
(Zahlen in Botschaften aus der Messung, nicht abgeschrieben).

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-07 — jüngster Lauf vom Vortag, unverändert: `image-scan` grün,
`upstream-drift` rot mit Fremd-Meldungen (darunter die Baseline-Hebung, die
diesen Slice ausgelöst hat — ihr Pin-Vorgang bleibt eigener Slice).

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
