# slice-251: `targets.authority` nimmt eine Liste an

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [`DC-FA-TGT-001`](../../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
(Scope — erweitert, kein neues Kürzel: Einzelmodul-Frage nach
[ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md));
Anlass ist der eingehende
[CR von `ai-harness-init`](../../cr/2026-10-06-cr-eingehend-ai-harness-init-targets-authority-liste.md);
Vorgänger am selben Modul [ADR-0099](../../adr/0099-targets-makefiles-glob.md)
(hält `authority` wörtlich — bleibt so);
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(String-Form byte-identisch).

**Berührte Spec-Stellen:** [§DC-FA-TGT-001.a](../../../../spec/spezifikation.md#dc-fa-tgt-001a--deklarations-konsistenz-doku-und-build-targets-targets)
(Schritte 1 und 5), [`SPEC-005`](../../../../spec/spezifikation.md#spec-005--d-checkyml)
(Schema-Zeile `targets.authority`).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-06.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** `targets.authority` nimmt neben einem String eine Liste wörtlicher
Pfade an; `gate-undocumented` misst gegen die Vereinigung der dort
dokumentierten Targets. Die String-Form bleibt byte-identisch, eine fehlende
Datei ist Exit 2, eine leere Liste lässt Richtung 2 entfallen wie
`doc-tables: []` Richtung 1.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Disjunktheits-Prüfung (ein Target in zwei Autoritäts-Dateien als
  Befund)** — die Regel „kein Target steht in beiden Teilen" existiert nur in
  einem offenen CR des Absenders an den Kurs; die gepinnte Baseline kennt sie
  nicht. Benannter Folgeschritt mit Auslöser: die Kurs-Baseline nimmt die
  Regel an (dann opt-in-Schlüssel mit eigenem Grund-Code).
- **Glob für `authority`** — die Autoritäts-Menge ist klein und fest, der CR
  sagt selbst, wörtliche Pfade genügen; die Vorgänger-ADR hält `authority` wörtlich.
- **`doc-tables`, `exempt-targets`, Regel-Extraktion** — vom CR
  ausgeschlossen, Bestand bleibt.
- **Werkzeug-eigener Index in diesem Repo** — anderer Vorgang: d-check hat
  keine werkzeug-eigenen Fragmente; die Kurs-Regel wäre Thema einer
  Baseline-Hebung.
- **Benutzerhandbuch, README, CHANGELOG** — Release-Prep
  (`AGENTS.md` §5 Regel 17).

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Liefer-Punkte (3):

- [ ] Kern: Vereinigung über alle Autoritäts-Dateien, Meldung nennt bei einer
      Datei den heutigen Wortlaut und bei mehreren alle; Dokument-Dubletten
      einmal gelesen; Tests nach den CR-Akzeptanzkriterien 1–6, darunter das
      rote Gegenbeispiel — rot ohne die Änderung aus dem richtigen Grund
      (Bewusstes Brechen, Modul 11).
- [ ] Config-Rand: String oder Liste; leerer, Null- oder Nicht-Pfad-Listeneintrag
      ⇒ Exit 2, Pfad-Regel je Eintrag; die String-Form dekodiert wie zuvor;
      `--print-config`-Gerüst mit Listen-Beispiel. *(Plan-Änderung nach R1:
      die ursprünglich geplante Ablehnung von Glob-Zeichen entfällt — sie
      brach die zugesagte Byte-Identität eines wörtlichen Pfads mit `[`.)*
- [ ] Lastenheft (Erweiterung, Bump + Historie nach
      [`MR-032`](../../../../harness/conventions.md#mr-032)), Spezifikation,
      Schema; ADR; Antwort im CR-Dokument samt Folgeschritt-Vermerk.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor —
      Rollenwechsel nach Schritt 8, kein Self-Review.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | Anforderung erweitert, Akzeptanzkriterien, Out-of-Scope, Bump + Historie |
| `spec/spezifikation.md` | update | Schritte 1/5, Schema-Zeile, Historie |
| `internal/hexagon/core/model/config.go` | update | `Authority` wird Liste |
| `internal/hexagon/core/rules/targets.go` | update | Vereinigung, Meldung |
| `internal/adapter/driven/configyaml/configyaml.go` | update | String-oder-Liste, Validierung |
| `internal/adapter/driving/cli/config_template.go` | update | Gerüst-Beispiel |
| `internal/hexagon/core/rules/targets_test.go`, `configyaml_test.go` | update | CR-Akzeptanzkriterien 1–6 |
| `docs/plan/adr/0100-…`, `docs/plan/adr/README.md` | neu / update | Entscheidungen |

**Spiegel vor dem Editieren** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
Lastenheft (Beschreibung, fail-closed-Satz, Out-of-Scope), Spezifikation
(Schritte 1 und 5, Schema-Zeile `targets.authority`), `--print-config`-Gerüst,
Meldungstext von `gate-undocumented`, `--doctor`-Klartext von
`gate-undocumented` (nennt er „die Autoritäts-Doku" im Singular?), Handbuch
§5-Beispiel und §6-Modultabelle (Release-Prep), `suggest.go` (erzeugt es einen
`authority`-Wert?), Kommentare an `rawTargets`/`TargetsConfig`.

## 4. Trigger

**Start** (`next` → `in-progress`): Implementer übernimmt; `in-progress/`
leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): verlangt die String-oder-Liste-Form eine
  Änderung am gemeinsamen Decoder über `targets` hinaus, wird neu geschnitten.
- `in-progress` → `open` (blockiert): keiner bekannt.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- Ein Konsument, der den Meldungstext von `gate-undocumented` auswertet,
  sähe bei einer Liste einen anderen Wortlaut. Die String-Form muss ihn
  byte-identisch lassen. — **Ausgang:** *(offen)*
- Nimmt der Kurs die Disjunktheits-Regel an, entsteht eine Lücke, die dieser
  Slice bewusst offen lässt. — **Ausgang:** *(offen)*

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

**Vorgelagert — offene Beobachtungen sichten:** drei Einträge der Sub-Area
treffen den Gegenstand —
[`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/state.md)
(zuletzt in slice-250: jede Grenze dieses Slice wird gegen den Code geprüft,
nicht gegen die Beschreibung),
[`BEO-ALL/commit-message-overclaims-work`](../observations/BEO-ALL/commit-message-overclaims-work/state.md)
(Rot-Gründe in Botschaften nur, wie der Test sie zeigt) und
[`BEO-ALL/semantic-change-body-only-edges-stale`](../observations/BEO-ALL/semantic-change-body-only-edges-stale/state.md)
(Spiegel in §3 vorab gelistet, einschließlich `--doctor`-Klartext).

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-06 — unverändert gegenüber slice-250 (`image-scan` grün;
`upstream-drift` rot mit drei Fremd-Meldungen). Keine berührt diesen
Gegenstand.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
