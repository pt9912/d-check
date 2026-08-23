# Slice slice-130: Historie-Form auf vier Spalten, und unsere Strenge deklarieren

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-83-baseline-v5110-migration](../welle-83-baseline-v5110-migration.md)
(Etappe C, geschnitten vom Delta-Audit).

**Bezug:** Baseline-Regelwerk
[`grundlagen-source-precedence.md`](../../../../.harness/baseline/v5.11.0/regelwerk/grundlagen-source-precedence.md)
(Kurs-Welle 90: CR-Pflicht am Lastenheft-Status, Tatsachenberichtigung,
zurückgezogene Anforderungen) und die Vorlage
`.harness/baseline/v5.11.0/templates/spec/lastenheft.template.md`;
[`spec/lastenheft.md`](../../../../spec/lastenheft.md) §7 Historie.

**Berührte Spec-Stellen:** `spec/lastenheft.md` §7 (Form der Historie-Tabelle),
Kopf-Feld `**Status:**` — keine Anforderung ändert ihre Aussage.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

Kurs-Welle 90 bindet die **CR-Pflicht an den Lastenheft-Status** und gibt der
Historie eine vierte Spalte. Zwei Deltas gegen unseren Bestand, beide vom
Delta-Audit belegt:

- **Die Lastenheft-Vorlage trägt `| Version | Datum | Änderung | Verweis |`** —
  die vierte Spalte nennt den externen CR-Vorgang, bei einer
  **Tatsachenberichtigung** ein `—`. Unser Lastenheft führt drei.
  **Nur das Lastenheft:** `spec/spezifikation.md` §7 führt **zwei** Spalten und
  trifft damit die vendorte `spezifikation.template.md` exakt — der Kanon gibt
  dem Technik-Stratum bewusst keine `Verweis`-Spalte, weil es keinen
  Change-Request-Vorgang kennt.
- **Vor `Accepted` verlangt der Kanon *nichts*:** *„frei änderbar, ohne Change
  Request, ohne Historie-Zeile"*. Unser Lastenheft steht auf **`Draft`** — wir
  fahren Versions-Bumps und Historie-Zeilen also **strenger als verlangt**.
  Das ist eine legitime Wahl, aber sie ist heute **undeklariert**: ein Leser
  kann nicht unterscheiden, ob unsere Disziplin Pflicht oder Vorsatz ist.

## 2. Vorgehen

1. Die vierte Spalte **einführen — nur im Lastenheft**; Bestandszeilen
   bekommen `—`, weil es für sie keinen externen Vorgang gibt. Die
   Spezifikations-Historie bleibt bei zwei Spalten.
2. **Die eigene Strenge deklarieren** — als Adaption im Konventionsspeicher
   (`MR-`Eintrag) oder als Satz im Lastenheft-Kopf; **im Slice entscheiden**,
   nicht vorab. Kriterium: eine Form-Frage gehört in den Konventionsspeicher,
   eine Vertrags-Aussage ins Lastenheft.
3. Prüfen, ob die `structure`-Regel über `## 7. Historie` (Chronologie-Monotonie
   auf Spalte 1) von der neuen Spaltenzahl berührt ist.
4. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Status-Wechsel des Lastenhefts.** Ob es von `Draft` auf `Accepted`
  geht, ist eine Auftraggeber-Entscheidung mit Vertragswirkung, kein Nachzug.
- **Kein Rückbau der bisherigen Strenge.** Der Kanon *erlaubt* weniger; er
  verbietet nicht mehr.
- **Keine Anwendung der Zurückgezogen-Regel** — dieses Repo hat keine
  entfallene Anforderung.

## 4. Definition of Done

- [ ] Die **Lastenheft**-Historie trägt die vierte Spalte, Bestandszeilen `—`;
      die Spezifikations-Historie ist **unverändert** bei zwei Spalten.
- [ ] Die eigene Strenge ist deklariert, mit begründeter Ortswahl.
- [ ] Die `structure`-Regel auf `## 7. Historie` läuft unverändert grün.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Eine Spalte in einer chronologischen Bestandstabelle zu ergänzen berührt
  die `table-order`-Regel.** Sie prüft Spalte 1 — die Zahl der Spalten ändert
  sich, die Schlüsselspalte nicht. Geprüft, nicht angenommen. — **Ausgang:**
  *(bei Closure)*
- **„Strenger als der Kanon" kann als Verstoß gelesen werden**, wenn die
  Deklaration fehlt oder am falschen Ort steht. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-129](../done/slice-129-baseline-v5110-delta-audit.md)
in `done/`.

**Rückführungen:** `in-progress` → `next`, falls die Ortswahl für die
Deklaration eine Auftraggeber-Entscheidung verlangt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Spec-Straten (GF), Konventionsspeicher (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23): **`BEO-002`**
  für die Ränder der Historie-Form — und ausdrücklich für die Frage, welches
  Stratum sie überhaupt trägt (die erste Fassung dieses Slice hat sie falsch
  beantwortet);
  **`BEO-011`** für jede Aussage darüber, welche Tabellen die Form tragen.

Slice-ID: slice-130. Betroffene IDs:
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus). Module:
Spec-Straten, Konventionsspeicher. Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Form-Nachzug an der adoptierten Vorlage.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
