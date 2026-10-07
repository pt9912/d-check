# slice-257: CVE-Scan je Plattform des Image-Index

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** [ADR-0066](../../adr/0066-cve-scan-gegen-das-publizierte-image.md)
(CVE-Scan gegen das publizierte Image); Folge von
slice-256 (Multi-Arch-Index).

**Berührte Spec-Stellen:** — *(Nachtlauf-Sensor, keine Spec-Aussage)*.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** `make image-scan` scannt **jede Plattform** des publizierten Index,
nicht nur die, die Trivy für den Host auswählt. Heute nennt
`tools/image-scan.sh` je Registry eine Referenz ohne Plattform; bei einem
Index wählt Trivy die Variante des Runners (`amd64`), und die
`arm64`-Variante — eigene distroless-Pakete — bliebe ungescannt, ohne dass
der Lauf es sagt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Index selbst** — slice-256; dieser Slice setzt ihn voraus.
- **Eine Schwelle oder ein rotes Urteil auf Befunde** — [ADR-0066](../../adr/0066-cve-scan-gegen-das-publizierte-image.md) lässt den
  Scan berichten, nicht gaten; das ändert sich hier nicht.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] `tools/image-scan.sh` scannt je Referenz `linux/amd64` und
      `linux/arm64` (Plattform im Bericht genannt); die Plattformliste steht
      an einer Stelle.
- [ ] `harness/sensors/image-scan.md` nennt die Plattformen und die Grenze;
      `make gates` grün; Nachtlauf manuell ausgelöst und grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/`.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/image-scan.sh` | update | Plattform-Schleife |
| `harness/sensors/image-scan.md` | update | Vertrag und Grenze |

## 4. Trigger

**Start** (`next` → `in-progress`): slice-256 in `done/` und ein
Multi-Arch-Release veröffentlicht; `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `open` (blockiert): Trivy wählt bei einem Index keine
  Plattform per Flag — dann je Plattform-Digest scannen.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft.

## 6. Risiken und offene Punkte

- Der Nachtlauf verdoppelt seine Scan-Zeit. — **Ausgang:** *(offen)*

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

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** eine berührte Sub-Area: die
Distribution unter dem Default `*` (`ALL`); deklariert.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** keine Treffer für CVE-Scan
oder Nachtlauf.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-07 — `upstream-drift` und `image-scan` grün.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
