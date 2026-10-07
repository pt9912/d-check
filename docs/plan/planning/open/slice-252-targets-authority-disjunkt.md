# slice-252: `targets` prüft opt-in die Disjunktheit der Autoritäts-Dateien

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [`DC-FA-TGT-001`](../../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
(Scope — erweitert, kein neues Kürzel: Einzelmodul-Frage nach
[ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md));
Anlass ist der eingetretene Re-Evaluierungs-Trigger von
[ADR-0100](../../adr/0100-targets-authority-liste.md), gemeldet im
[Hinweis der Baseline](../../cr/2026-10-07-hinweis-eingehend-ai-harness-course-disjunktheit.md);
Beobachtung
[`BEO-ALL/disjunktheit-geteilter-gate-index-ungeprueft`](../observations/BEO-ALL/disjunktheit-geteilter-gate-index-ungeprueft/state.md);
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(ohne den Schalter byte-identisch).

**Berührte Spec-Stellen:** [§DC-FA-TGT-001.a](../../../../spec/spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets)
(neuer Schritt), [`SPEC-005`](../../../../spec/spezifikation.md#spec-005--d-checkyml)
(Schema-Zeile), [§4 Grund-Codes](../../../../spec/spezifikation.md#4-grund--und-fehler-codes)
(neue Zeile).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Mit `targets.authority-disjoint: true` meldet `targets` ein Target,
das als Tabellenzeile in mehr als einer `targets.authority`-Datei steht, mit
dem neuen Grund-Code `gate-declared-twice` — an der Tabellenzeile jeder
weiteren Datei, die Meldung nennt die Datei der ersten Nennung. Ohne den
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

- [ ] Kern: Disjunktheits-Prüfung mit `gate-declared-twice` (Fundstelle =
      Tabellenzeile in der späteren Datei, Konfigurations-Reihenfolge;
      `exempt-targets` wirkt nicht; mit einer Datei wirkungslos); Tests rot
      ohne die Änderung aus dem richtigen Grund (Bewusstes Brechen).
- [ ] Config-Rand und Lexik: Schlüssel `targets.authority-disjoint`
      (Bool, strikt), `--print-config`-Gerüst, Grund-Code in `AllReasons`
      und im `--doctor`-Klartext.
- [ ] Lastenheft (Erweiterung, Bump + Historie nach
      [`MR-032`](../../../../harness/conventions.md#mr-032)), Spezifikation
      (Schritt, Schema, §4); ADR; Antwort im Hinweis-Dokument; eigene
      Black-Box-Probe gegen ein Vorher-Image vor dem Review (ohne Schalter
      byte-identisch), Ergebnis im Plan notiert.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor —
      Rollenwechsel nach Schritt 8, kein Self-Review.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben
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

**Spiegel vor dem Editieren** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
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
  Repo zwei Dinge. — **Ausgang:** *(offen)*
- Die Black-Box-Probe von Hand kann wieder nur die Fälle prüfen, an die der
  Autor denkt. — **Ausgang:** *(offen)*

## 7. Closure-Notiz

*(gefüllt vor dem `git mv` nach `done/`)*

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

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
[`BEO-ALL/disjunktheit-geteilter-gate-index-ungeprueft`](../observations/BEO-ALL/disjunktheit-geteilter-gate-index-ungeprueft/state.md)
(der Gegenstand selbst; bekommt seinen Ausgang),
[`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/state.md)
(dreimal in Folge an diesem Modul — deshalb die eigene Black-Box-Probe vor
dem Review) und
[`BEO-ALL/commit-message-overclaims-work`](../observations/BEO-ALL/commit-message-overclaims-work/state.md)
(Zahlen in Botschaften aus der Messung, nicht abgeschrieben).

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-07 — jüngster Lauf vom Vortag, unverändert: `image-scan` grün,
`upstream-drift` rot mit Fremd-Meldungen (darunter die Baseline-Hebung, die
diesen Slice ausgelöst hat — ihr Pin-Vorgang bleibt eigener Slice).

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
