# Slice slice-037: `--suggest-config` — Kennungs-Präfix als Option

**Status:** done (Closure 2026-06-21;
[ADR-0015](../../adr/0015-suggest-config-id-prefix.md) `Accepted`).

**Welle:** welle-26-suggest-prefix (Trigger: a-check-Bootstrap 2026-06-20 —
`--suggest-config ai-harness-init` emittierte d-checks **eigenes** `DC-`-Muster
in ein Fremd-Repo; a-check musste `DC-(FA-[A-Z]+|QA)` von Hand auf `AC-`
umschreiben).

**Bezug:**
[`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)
(`--suggest-config`, wird erweitert),
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
(das emittierte `ids`-Muster),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(deterministisch),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(read-only).

**Autor:** pt9912. **Datum:** 2026-06-20.

---

## 1. Ziel

`--suggest-config ai-harness-init` ist „Voll-Kanon fürs **leere** Repo" —
also explizit zum Bootstrappen *fremder* Harness-Repos gedacht, backt aber
d-checks eigenes Kennungs-Präfix `DC-` ein. Jedes Nicht-d-check-Repo
(a-check=`AC`, b-cad=`BC`, …) bekommt das falsche Muster und muss
nacheditieren. Das Präfix wird **parametrisierbar**; ohne Angabe wird **nicht
still `DC-`** emittiert.

## 2. Design (Vorschlag)

- **`ai-harness` (repo-bewusst):** Präfix aus dem vorhandenen Lastenheft
  **ableiten** (erste `…-FA-`-Kennung) — null Konfiguration im typischen Fall.
- **`ai-harness-init` (leeres Repo):** kann nichts ableiten → explizite
  Option `--id-prefix <PREFIX>` (oder Token-Form `ai-harness-init:<PREFIX>`).
- **Ohne Angabe:** markierter Platzhalter `<PREFIX>` plus `# TODO`-Kommentar
  in der Ausgabe — **kein** stiller `DC-`-Default (Anti-Footgun, „kein
  stiller falscher Default").

## 3. Zu entscheiden (im Slice)

- **Flag vs. Token:** `--id-prefix AC` (neues Flag) oder
  `ai-harness-init:AC` (Quelle trägt das Präfix). Empfehlung: Flag — die
  Quelle bleibt ein Bezeichner, das Präfix ist orthogonale Konfiguration.
- **Ableitungsquelle (`ai-harness`-Modus):** erste FA-Kennung im
  Lastenheft? Konfliktfall (mehrere Präfixe) → Fehler oder erstes gewinnt?
- **Default-Verhalten:** Platzhalter `<PREFIX>` ist eine **Verhaltensänderung**
  ggü. heute (`DC-` fix) → ADR-pflichtig (Gate/Default-Änderung,
  `AGENTS.md` §3.6 sinngemäß).
- **Lastenheft-CR:** Umfang der Erweiterung von
  [`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)
  (Präfix-Parameter + Ableitung) mit Happy/Boundary/Negative.

## 4. Definition of Done (vorläufig)

- [x] Lastenheft-CR (0.20.0):
  [`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)
  um Präfix-Parameter/Ableitung + Platzhalter-Default erweitert.
- [x] CLI: **Flag** `--id-prefix` (Entscheidung: Flag, nicht Token) wird von
  beiden `ai-harness`-Quellen konsumiert; `ai-harness` leitet aus dem
  Lastenheft ab (Konflikt ⇒ Fehler).
- [x] Ohne Präfix: Platzhalter `<PREFIX>` + `# TODO`, kein `DC-`.
- [x] Tests Happy/Boundary/Negative (+ Ableitung, Konflikt, ungültiger
  Wert); Determinismus
  ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)),
  read-only ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- [x] Benutzerhandbuch §`--suggest-config` ergänzt; `CHANGELOG` (Breaking).
- [x] [ADR-0015](../../adr/0015-suggest-config-id-prefix.md) (Default-Änderung)
  geschrieben; `make gates` grün.
- [x] Unabhängiges Review R1 (2 MEDIUM behoben) + R2 (bestätigt); Closure;
  [ADR-0015](../../adr/0015-suggest-config-id-prefix.md) → `Accepted`.

## 5. Risiken / offene Punkte

- **Rückwärtskompatibilität:** wer heute auf `ai-harness-init` → `DC-`
  baut, sieht künftig den Platzhalter. Bewusst (der alte Default war für
  Fremd-Repos schlicht falsch), aber als Breaking-Change in der Historie
  auszuweisen.
- **Ableitungs-Heuristik** kann bei gemischten Präfixen danebenliegen →
  klar definierter Konfliktfall nötig.

## 6. Trigger

a-check-Bootstrap 2026-06-20: das Init-Template ist zum Bootstrappen
*fremder* Repos da, emittiert aber d-checks Eigen-Präfix — die stille
Fehlsetzung, die das Regelwerk gerade verbietet. Schwester-Repo a-check ist
der erste reale Konsument.

## 7. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (CLI-/Core-/Doku-Arbeit; Greenfield-Default).

## 8. Closure-Notiz (nach `done/`)

**Umsetzung.** Richtung der drei Auftraggeber-Entscheidungen: **Flag**
`--id-prefix` (kein Token), **Platzhalter `<PREFIX>` + TODO** ohne
Präfix/Ableitung (kein stiller `DC-`; **Breaking** ggü. 0.18.1), **Fehler**
bei mehrdeutigem Präfix. Nur das Anforderungs-Muster ist parametrisiert
(`harnessIDPatterns(reqPrefix)` in `internal/hexagon/core/app/suggest.go`);
`ADR-`/`MR-`/`slice`/Carveout bleiben konventions-fest. Quelle des Präfix:
`--id-prefix` (gewinnt immer) → im Modus `ai-harness` Ableitung aus
`spec/lastenheft.md` (`deriveReqPrefix`, Konflikt ⇒ Exit 2) → sonst
Platzhalter. Spec-CR
[`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)
(Lastenheft 0.20.0) + Technik in [`spec/spezifikation.md` §DC-FA-CLI-006.a](../../../../spec/spezifikation.md#dc-fa-cli-006a--konfigurations-vorschlag);
[ADR-0015](../../adr/0015-suggest-config-id-prefix.md) `Accepted`.

**Belege.** `make gates` grün (doc-check, lint, test, arch-check, coverage
94,20 %, semgrep 55/0, gate-consistency); 5+2 neue `TestCLI037_*`-Fälle;
kein Regress an den ai-harness-Alt-Tests (ai-harness leitet für d-check
selbst weiter `DC` ab). Read-only
([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)),
deterministisch
([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus), 10×
byte-identisch im R1/R2-Lauf belegt).

**Review R1** (`docs/reviews/2026-06-21-slice-037-suggest-config-id-prefix.md`):
0 HIGH/2 MEDIUM/1 LOW/1 INFO. MEDIUM-1 (stale
[`spec/spezifikation.md` §DC-FA-CLI-006.a](../../../../spec/spezifikation.md#dc-fa-cli-006a--konfigurations-vorschlag)
+ falsches ADR-`Schärft`) und MEDIUM-2 (ungetesteter
Override/Konflikt-Bypass) behoben. **Review R2**
(`docs/reviews/2026-06-21-slice-037-r2-verifikation.md`): beide bestätigt
behoben, keine neuen Findings, keine Regression. LOW-1 (aktiver
`<PREFIX>`-Platzhalter — mildes Silent-Green bei unverändert übernommener
Vorlage) und INFO-1 (`--id-prefix` außerhalb der ai-harness-Modi ignoriert)
won't-fix: entschiedener advisory-Default mit TODO bzw. im Flag-Hilfetext
verortet.

**Lerneintrag.** Eine CLI-Default-Änderung berührt **beide** Spec-Straten:
das Lastenheft (Vertrag, CR) **und** `spec/spezifikation.md` (Technik) — das
Technik-Stratum mitzuziehen wurde erst im unabhängigen Review sichtbar
(MEDIUM-1). Das `Schärft`-Feld einer ADR muss auf die real existierende
geschärfte Spec-Stelle zeigen, nicht reflexhaft „keine Spec-Stelle".
