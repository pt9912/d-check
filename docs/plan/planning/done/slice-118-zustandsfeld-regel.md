# Slice slice-118: Die Zustandsfeld-Regel verkörpern — Briefing und Reviewer-Anker

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-81-zustandsfelder](../welle-81-zustandsfelder.md) (zugeordnet
bei der Eröffnung).

**Bezug:** [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage)
(die Baseline-Aussage, unter der die Regel gilt),
[`MR-015`](../../../../harness/conventions.md#mr-015) (das Briefing routet und
spiegelt nicht), Reviewer-Skill; Baseline-Regelwerk
`grundlagen-harness-dateien.md` §Was ein Kommentar trägt („Dieselbe Regel für
Zustandsfelder", „Die Kopfzeile eines lebenden Registers ist derselbe Fall")
und `modul-06-roadmap.md` §Das Beobachtungs-Register.

**Berührte Spec-Stellen:** — (Briefing und Harness-Skills; keine Spec-Zeile).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

Die Regel steht, bevor sie angewandt wird — dieselbe Reihenfolge, die in der
Vor-Welle getragen hat. Die Hard Rule des Briefings zu Kommentar-Klassen gilt
ausdrücklich auch für **Zustandsfelder** (`Stand`/`Status` in Register,
Roadmap und Meilenstein-Tabelle): Zustand und Beleg als auflösbarer Anker
statt Chronik, und das Drift-Log trägt nur Umplanungen. Der Reviewer-Skill
bekommt den zugehörigen **HIGH**-Anker, wie die Baseline ihn führt — denn kein
Gate fängt das.

## 2. Vorgehen

1. **Briefing:** die Kommentar-Hard-Rule um den Zustandsfeld-Satz erweitern —
   Geltung (Register-, Roadmap-, Meilenstein-Zellen), Form (Zustand + Beleg
   als auflösbarer Anker), Abgrenzung (Chronik gehört nicht dorthin; ihre Orte
   sind Vorhaben, Closure-Log, Drift-Log, `git`). Kurz und routend, nicht
   spiegelnd.
2. **Reviewer-Skill:** neuer HIGH-Eintrag *Zustandsfeld trägt Chronik* mit
   beiden Ausprägungen (erzählende `Stand`-Zelle; Drift-Log, das Schließungen
   oder erreichte Meilensteine protokolliert) und dem ausdrücklichen Hinweis,
   dass kein Gate ihn stützt; Version und Datum mitziehen.
3. **Spiegel-Liste vor dem Editieren** (Semantik-Fläche): Briefing,
   Reviewer-Skill, `harness/README.md` (Guides-Zeile), Konventionsspeicher —
   per `grep` nach der bestehenden Kommentar-Regel abgeleitet, nicht erinnert.
4. Unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Anwendung:** keine Kopfzeile, keine `Stand`-Zelle, keine
  Drift-Log-Zeile wird hier angefasst (slice-119, slice-120).
- **Kein neuer Sensor** — die Baseline benennt Briefing und Reviewer als
  Träger; ein Gate ist ausdrücklich nicht vorgesehen.
- **Kein Produkt-Code.**

## 4. Definition of Done

- [ ] Briefing-Hard-Rule deckt Zustandsfelder; Geltung, Form und Abgrenzung
      stehen dort, ohne die Baseline zu duplizieren.
- [ ] Reviewer-Skill trägt den HIGH-Anker samt beider Ausprägungen; Version
      und Datum nachgezogen.
- [ ] Spiegel-Liste abgeleitet und abgehakt.
- [ ] `make gates` grün; unabhängiger Review; Closure-Notiz; Register
      gesichtet.

## 5. Abnahme-Punkte / Risiken

- **Ein HIGH-Anker ohne Gate ist eine Zusage an Menschen** — er muss so
  formuliert sein, dass ein Reviewer ihn ohne Nachschlagen anwenden kann, und
  darf keine Scheinbefunde auf legitime Formen erzeugen (ein Datum, das ein
  benannter Trigger pflegt, ist kein Chronik-Feld). — **Ausgang:** *(bei
  Closure)*
- **Das Briefing routet, es spiegelt nicht:** der Satz muss kurz bleiben und
  auf die Baseline zeigen, statt sie nachzuerzählen. — **Ausgang:** *(bei
  Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): slice-117 in `done/` (die Regel muss
vendored sein, bevor das Briefing sie zitiert).

**Rückführungen:** keine erwartet.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Briefing und Harness-Skills (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-22): BEO-002 wirkt
  als Spiegel-Pflicht (§2 Schritt 3); BEO-006/009/010 als Arbeitsregeln.

Slice-ID: slice-118. Betroffene IDs:
[`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage),
[`MR-015`](../../../../harness/conventions.md#mr-015). Module: Briefing,
Reviewer-Skill. Gates: `make doc-check` (eng), `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Regel-Adoption in eigenen Artefakten.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
