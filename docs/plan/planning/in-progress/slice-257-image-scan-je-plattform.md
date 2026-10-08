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

- [x] `tools/image-scan.sh` scannt je Referenz jede Plattform des Index — gelesen
      aus dem Index selbst, nach der Plan-Änderung nach R1 statt einer festen
      Liste `linux/amd64`/`linux/arm64` (Plattform im Bericht genannt); die
      Übersteuerung `IMAGE_SCAN_PLATFORMS` steht an einer Stelle.
- [x] `harness/sensors/image-scan.md` nennt die Plattformen und die Grenze;
      `make gates` grün; Nachtlauf manuell ausgelöst und grün.
- [x] Review durchgeführt, Report unter `docs/reviews/`.
- [x] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/image-scan.sh` | update | Plattform-Schleife |
| `harness/sensors/image-scan.md` | update | Vertrag und Grenze |

*(Plan-Änderung vor dem Code-Commit: ein **Plattform-Nachweis** je Lauf — gemessen scannt Trivy ein Einzel-Manifest-Image bei `--platform linux/arm64` still als amd64 mit Exit 0; ohne Nachweis wäre der arm64-Scan dort eine Behauptung. Dazu die Funktion `arch_aus_json` mit vier Selbsttest-Proben und der Spiegel „BEIDE Trivy-Läufe" im Skriptkopf und in der Sensor-Datei. Nach R1: die Plattformen kommen aus dem Index selbst statt aus einer Kopie (F-1), ein Ref ohne lesbaren Index gilt als gescheitert; Geschichte-Anhänge an [ADR-0066](../../adr/0066-cve-scan-gegen-das-publizierte-image.md) und [ADR-0102](../../adr/0102-multi-arch-index-und-spiegel-per-index-digest.md) (F-2, F-3). Nach R2: die Plattformliste gilt nur vollständig — ein Abbruch von `imagetools` oder ein Index-Eintrag ohne Plattform machen sie leer (R2-1, HIGH), die Ursache kommt aus demselben Aufruf (R2-3).)*

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

- Der Nachtlauf verdreifacht seine Trivy-Läufe je Image (Plattform-Nachweis, Vollbericht, Entscheidung, je Plattform). — **Ausgang:** entfallen — gemessen: der Nachtlauf mit vier Scans brauchte rund 45 s bei 30 min Timeout, lokal 19 s.

## 7. Closure-Notiz

- **Was hat funktioniert:** Jede Zusage des Scans ist an einem Bruch belegt,
  nicht nur an einem grünen Lauf: das alte Einzel-Manifest mit verlangtem
  arm64 (still als amd64 gescannt — der Grund für den Plattform-Nachweis), eine
  fehlende Plattform, ein fehlender Ref, ein präparierter Index mit einem
  Eintrag ohne Plattform, eine fehlerhafte `config.json`. Die Auswertung ist in
  drei netzlos prüfbare Funktionen geschnitten (Zählung, Architektur,
  Plattformliste) mit 17 Selbsttest-Proben.
- **Was ging anders als geplant:** Der Plan sah eine feste Plattformliste vor;
  gemessen scannt Trivy ein Einzel-Manifest bei verlangtem arm64 still als
  amd64 (Plan-Änderung vor dem Code: Plattform-Nachweis). R1 zeigte, dass die
  feste Liste eine Kopie des Release-Pfads ist (F-1); ihr Ersatz — die
  Plattformen aus dem Index lesen — führte einen neuen Leseweg ein, der Fehler
  von `imagetools` verschluckte: eine Teil-Antwort galt als vollständige Liste,
  der Scan meldete Exit 0 ohne arm64 (R2-1 HIGH, gemessen). Die Behebung ließ
  stderr in die Liste laufen (R3-1). Drei Runden für einen Fix, der jeweils
  den nächsten Rand anfasste.
- **Steering-Loop-Eintrag:** Reviewer-Skill ergänzt: Prüffrage 19 — ein neuer
  Leseweg einer Härtung wird gegen seine Fehlerformen gefahren, nicht nur
  gegen den Fall, den er beheben soll — liegt in `.harness/skills/reviewer.md`.
  Auslöser: `BEO-ALL/haertung-kippt-fehlerpolitik-ungeprueft` (slice-156, slice-206, slice-257 — 3×).
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-257.md` in
  [`BEO-ALL/haertung-kippt-fehlerpolitik-ungeprueft`](../observations/BEO-ALL/haertung-kippt-fehlerpolitik-ungeprueft/state.md)
  ergänzt — Zähler 3×, Ausgang verkörpert (Auftraggeber-Entscheid);
  [`BEO-ALL/fix-commit-ausserhalb-review-range`](../observations/BEO-ALL/fix-commit-ausserhalb-review-range/state.md)
  neu angelegt mit Belegen aus slice-256 und slice-257 (2×).
- **Folge-Slices:** keine.
- **Risiken aus §6:** entfallen (Begründung in §6). Trigger-Audit: kein
  Carveout, kein bootstrap-aware Gate, keine Hard Rule mit eingetretenem
  Trigger; [ADR-0066](../../adr/0066-cve-scan-gegen-das-publizierte-image.md)
  und [ADR-0102](../../adr/0102-multi-arch-index-und-spiegel-per-index-digest.md)
  tragen je einen Geschichte-Anhang (drei Läufe je Plattform; die
  Offen-Konsequenz von [ADR-0102](../../adr/0102-multi-arch-index-und-spiegel-per-index-digest.md) eingelöst). Nachtlauf-Stand
  ([`MR-053`](../../../../harness/conventions.md#mr-053)): wie in §8.
- **Drei Paarungen:** (a) Anker — `.harness/skills/reviewer.md` trägt
  `(seit slice-257)` in Prüffrage 19; (b) Folge-Slice — keine; (c) Register —
  beide zitierten Beobachtungen existieren und tragen Belege.

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
