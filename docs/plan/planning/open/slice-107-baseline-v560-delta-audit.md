# Slice slice-107: Baseline-Bump v5.0.0 → v5.6.0 — Etappe B (Stufen-Audit)

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** welle-78-baseline-v560-migration.

**Bezug:** Präzedenz slice-085 (das v5.0.0-Modul-Delta-Audit: 18 Findings,
je mit Zuordnung); Grundlage ist der in
[slice-106](slice-106-baseline-v560-vendoring.md) vendorte Baum. Kein
`DC-*`-Bezug — Lese-/Planungs-Arbeit.

**Autor:** pt9912. **Datum:** 2026-08-21.

---

## 1. Ziel

Das Regelwerks-Delta v5.0.0 → v5.6.0 **je Stufe gegen die Tag-Notizen** lesen
(nicht pauschal) und für jede neue oder geänderte Regel eine von drei
Antworten festhalten — **konform bereits** · **anzupassen** (mit Fundstelle
und Etappe-C-Kandidat) · **nicht anwendbar** (mit Begründung). Das Ergebnis
ist ein Findings-Register in diesem Slice, aus dem die Etappe-C-Slices
geschnitten werden.

Die sechs Stufen:

1. **v5.1.0** — §Vergabe unter §ID-Schema (wer vergibt die nächste Kennung).
2. **v5.2.0** — Straten-IDs, Bestands-Stichprobe, Reconciliation-Register.
3. **v5.3.0** — Kommentar-Regel (`grundlagen-harness-dateien.md`).
4. **v5.3.1/v5.4.0** — Korrekturen („Zwei Sensoren an derselben Aussage",
   „Zwei Kopien, zwei Antworten") + drei Regel-Ergänzungen.
5. **v5.5.0** — **Team-Fähigkeit** (größter Block: Rolleninhaber,
   Konflikt-Terminal, `team.md` dreistufiges SOLL; dazu `lab/team-sim/`).
6. **v5.6.0** — TA-7 nennt seine Wirkung (Hauptzweig-Regel, ein Absatz).

## 2. Leitfragen je Stufe

- Verlangt die Regel ein **Artefakt**, das d-check nicht führt (z. B.
  `team.md`, Reconciliation-Register)? Gilt sie auch für ein
  Ein-Operator-Repo, oder deklariert sie selbst ihre Grenze?
- Widerspricht eine **bestehende** d-check-Konvention der neuen Regel
  (Adaption nötig — `MR-*` — oder Anpassung)?
- Ist etwas, das d-check **bereits lebt**, jetzt Baseline-Default (dann ist
  ggf. eine bestehende Adaption auflösbar — Präzedenz: MR-018/019/020/022
  wurden bei v5.0.0 „Baseline-Stand"-aufgelöst)?
- **Wiedervorlage aus slice-090:** die upstream notierte 5-vs-6-Finding-
  Feld-Drift der Baseline — im v5.6.0-Stand behoben oder weiter offen?

## 3. Ausdrücklich NICHT in diesem Slice

Kein Editieren der eigenen Artefakte (außer diesem Slice selbst) — Etappe C
setzt um, dieser Slice liest und schneidet.

## 4. Definition of Done

- [ ] Je Stufe ein Abschnitt mit Findings-Tabelle (Regel · Antwort ·
      Fundstelle/Begründung); **kein** „pauschal konform" ohne je Regel eine
      Zeile (die welle-74-Lehre: eine Aufzählung, die vollständig heißt,
      braucht je Kandidat einen Negativbefund).
- [ ] Etappe-C-Schnitt: die „anzupassen"-Findings sind zu Slices gebündelt
      (Dateien in `open/`), die Wellendokument-§4-Tabelle ist nachgeführt
      (Roadmap-Drift-Log-Eintrag).
- [ ] Unabhängiger Review des Audits (Frischkontext gegen das vendored
      Delta); `make gates` grün (dieser Slice ändert nur Planungs-Doku).

## 5. Risiken

- **Die Team-Stufe könnte formal greifen, obwohl operativ ein
  Ein-Operator-Repo vorliegt** — dann ist die ehrliche Antwort eine
  deklarierte Adaption (`MR-*`) statt eines leeren Pflicht-Artefakts.
- **Additiv heißt nicht folgenlos:** auch eine neue Regel ohne
  Widerspruch kann ein stehendes Artefakt zur Pflicht machen (Präzedenz:
  das Beobachtungs-Register kam bei v5.0.0 genauso additiv).

## 6. Trigger

**Start** (`open` → `in-progress`):
[slice-106](slice-106-baseline-v560-vendoring.md) in `done/` (der Audit liest
den **vendorten** Baum, nicht das Kurs-Repo) **und** WIP-Slot frei.

**Rückführungen:** `in-progress` → `next`, falls das Delta eine Vorfrage
aufwirft, die der Auftraggeber entscheiden muss (z. B. Team-Stufe adoptieren
vs. adaptieren).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-/Planungs-Doku, Repo-Default GF.
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-21: keine
  unverkörperte offen): **BEO-002/MR-025** ist die Arbeits-Brille des Audits
  selbst — jede „anzupassen"-Zeile benennt ihre Spiegel gleich mit, statt
  sie Etappe C suchen zu lassen.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Lese-Audit gegen die adoptierte
Konvention, dokumentierte Präzedenz (slice-085).

## 9. Closure-Notiz (nach `done/`)

_Ausstehend._
