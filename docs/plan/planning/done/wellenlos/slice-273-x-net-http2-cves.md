# slice-273: `golang.org/x/net` auf v0.60.0 — vier HTTP/2-CVEs im publizierten Image

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — reaktive Arbeit, die Closure-Bedingung ist die DoD
selbst.

**Bezug:** Fund eines Fremd-Scanners am publizierten Image `v0.86.0`:
CVE-2026-78663 (9.1, CRITICAL), CVE-2026-78659, CVE-2026-78660,
CVE-2026-78669 (je 7.5, HIGH), alle in `golang.org/x/net` v0.59.0, behoben in
v0.60.0; [ADR-0066](../../adr/0066-cve-scan-gegen-das-publizierte-image.md)
(der eigene Scan), [ADR-0067](../../adr/0067-dependabot-als-hebender-kanal.md)
(der hebende Kanal).

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-10.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Das Image trägt `golang.org/x/net` v0.60.0; die vier CVEs sind im
nächsten Release behoben. Gehoben wird über `go.mod` und `make tidy`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Hebungen aus dem Dependabot-PR #7** (go-git, go-winio) — ein eigener
  Vorgang über den hebenden Kanal; sie beheben keinen der vier Funde.
- **Die Schweregrad-Lücke des eigenen Scans** — Trivy führt drei der vier als
  `UNKNOWN`, `make image-scan` bleibt deshalb grün. Das ist die offene
  Beobachtung `BEO-ALL/scanner-vendor-severity-lag`; sie bekommt einen Beleg,
  keine Änderung am Gate.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [x] `go.mod` führt `golang.org/x/net` v0.60.0, `go.sum` passt; `make gates`
      grün; der Vollbericht von Trivy gegen das neu gebaute Image führt die vier
      CVEs nicht mehr.
- [x] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [x] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `go.mod`, `go.sum` | update | Hebung, `make tidy` |

## 4. Trigger

**Start** (`next` → `in-progress`): `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `open` (blockiert): v0.60.0 verlangt eine neuere
  Go-Toolchain als die gepinnte — dann zuerst die Toolchain heben.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft. Geht mit dem Release v0.86.1 hinaus.

## 6. Risiken und offene Punkte

- **Die Hebung zieht weitere Module nach** — `x/net` v0.60.0 kann neuere
  Fassungen von `x/sys`, `x/text` oder `x/crypto` verlangen. —
  **Ausgang:** entfallen — nur die `x/net`-Zeilen änderten sich, `make tidy`
  lässt den Baum unverändert (Review, Verifikation).

## 7. Closure-Notiz

- **Was hat funktioniert:** `golang.org/x/net` steht auf v0.60.0, gehoben über
  `go.mod` und `make tidy` — `go.sum` ist ein Fixpunkt von `tidy`, kein
  weiteres Modul zog nach. Trivy findet im neu gebauten Image keine der vier
  CVEs mehr, dazu ist CVE-2026-97032 (MEDIUM) behoben; die Gegenprobe gegen
  `v0.86.0` zeigt alle vier, der Scan sieht sie also. Ein verfälschter
  `go.sum`-Hash lässt den Build mit `checksum mismatch` fallen.
- **Was ging anders als geplant:** Plan und Commit-Botschaft sagen „drei davon
  `UNKNOWN`“ — gemessen waren es mit der Datenbank vom 2026-10-09 vier, mit der
  vom 2026-10-10 zwei; die zusammengefassten Zellen der Trivy-Tabelle waren
  falsch gelesen (Verifikation L1). Die Commit-Botschaft ist eingefroren, das
  Register trägt die richtigen Zahlen.
- **Steering-Loop-Eintrag:** keiner mit neuer Schwelle.
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-273.md` in
  [`BEO-ALL/scanner-vendor-severity-lag`](../observations/BEO-ALL/scanner-vendor-severity-lag/state.md)
  (jetzt 2×): der eigene Scan blieb rund einen Tag hinter einem Fremd-Scanner
  zurück.
- **Folge-Slices:** keiner. Geht mit dem Release v0.86.1 hinaus.
- **Risiken aus §6:** entfallen (siehe §6). Trigger-Audit: kein Carveout,
  kein bootstrap-aware Gate, keine neue ADR; keine Hard Rule mit eingetretenem
  Trigger. Nachtlauf-Stand
  ([`MR-053`](../../../../harness/conventions.md#mr-053)): wie in §8.
- **Drei Paarungen:** (a) Anker — kein Eintrag mit `liegt in`; (b)
  Folge-Slices — keiner genannt; (c) Register — die zitierte Beobachtung
  existiert und trägt zwei Belege.

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist nur der
Abhängigkeits-Stand des Produkts, unter dem Default `*` (`ALL`).

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:**
[`BEO-ALL/scanner-vendor-severity-lag`](../observations/BEO-ALL/scanner-vendor-severity-lag/observation.md)
(1×, Beleg slice-201) trifft genau diesen Fall: Trivy führt drei der vier CVEs
als `UNKNOWN`, `make image-scan` meldete `v0.86.0` grün, ein anderer Scanner
stuft sie als CRITICAL und HIGH ein. Mit diesem Slice 2×.

**Messung beim Anlegen:** `make image-scan` gegen `v0.86.0`: Exit 0; der
Vollbericht führt CVE-2026-78659, -78660, -78663, -78669 unter
`golang.org/x/net` v0.59.0, Fix 0.60.0, drei davon mit Schweregrad `UNKNOWN`.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`image-scan` des Nachtlaufs vom 2026-10-09 lag vor dem Release `v0.86.0`; der
Fund kam über einen Fremd-Scanner.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
