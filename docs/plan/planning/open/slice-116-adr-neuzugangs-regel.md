# Slice slice-116: ADR-Neuzugangs-Regel + Erstanwendung an den `Proposed`-ADRs

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-80-struktur-ids](../welle-80-struktur-ids.md) (zugeordnet bei
der Eröffnung).

**Bezug:** AGENTS §3.5 (ADR-Immutabilität — `Accepted` bleibt unverändert,
`make adr-check`), die ADR-Index-Konvention aus slice-113, die beiden
`Proposed`-ADRs [ADR-0050](../../adr/0050-fence-unclosed-in-spans.md) und
[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md),
Baseline-ADR-Template (`Schärft:` nennt die Kennung, wo das Zielelement eine
trägt; Link auf Überschrift oder Sektion), Reviewer-Skill
`.harness/skills/reviewer.md` (Referenz-Richtungs-Anker), Entscheid D4.

**Berührte Spec-Stellen:** — (ADR-Felder und Harness-Skills; keine Spec-Zeile).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

Die neue Form wird **angewandt**, nicht nur deklariert: die beiden
`Proposed`-ADRs ziehen `Schärft:`/`Bezug:` auf die jetzt vergebenen
`SPEC-*`/`ARC-*` (kein Immutabilitäts-Bruch — `adr-check` greift nur bei
`Accepted`), das Slice-Kopf-Feld „Berührte Spec-Stellen" im Haus-Stil nennt
Kennungen statt `§N`, der Reviewer-Skill trägt den MEDIUM-Anker „Kennung
statt §-Anker, wo das Zielelement eine trägt" (Zwei-Formen-Welt: `Accepted`-
Bestand bleibt §, Neuzugänge Kennung), und der ADR-Index-Kopf zeigt ein
Beispiel beider Formen. Damit hat die Welle ihren ersten lebenden Konsumenten
auch auf der ADR-Seite — und der Review hat den Anker, um die Regel künftig zu
halten.

## 2. Vorgehen

1. **Messen:** `Schärft:`-Formen über alle ADRs (Stand der Inventur 2026-08-22:
   46 Felder; 28 auf `.a`-Sektionen, 10 „keine Spec-Stelle", 5 Architektur
   §-basiert, 2 ohne Anker, 1 `—`; 11 ohne Feld) — gemessen, die Zahlen in die
   Closure-Notiz.
2. **[ADR-0050](../../adr/0050-fence-unclosed-in-spans.md)/[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md) nachziehen:** Felder auf Kennungen, wo das Zielelement
   jetzt eine trägt (sonst bleibt der §-Anker); `## Geschichte`-Zeile je ADR.
   `make adr-check` muss **stumm** bleiben (Proposed) — Probe.
3. **ADR-Index-Kopf** (`docs/plan/adr/README.md`): Beispiel-Zeile für beide
   Formen; **Reviewer-Skill** 1.5.0 → 1.6.0: Anker im MEDIUM-Block, Version
   und Datum.
4. **Haus-Stil der Slices:** das Kopf-Feld „Berührte Spec-Stellen" in diesem
   und den künftigen Slices nennt Kennungen; die welle-80-Slices 114/115
   tragen sie nach der Vergabe bereits in ihrer Closure-Notiz.
5. **Spiegel** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
   ADR-Index-Kopf, Reviewer-Skill, AGENTS §5 (Satz aus slice-113 prüfen),
   `harness/README.md` Guides-Zeile zum Reviewer-Skill (Version?).
6. Unabhängiger Review; Closure — und die **Wellen-Closure** (Ergebnisnotiz,
   Register-Lese-Schritt, `make fullbuild`).

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Retarget eines `Accepted`-ADR** — auch nicht „nur das Schärft-Feld";
  `adr-check` würde es melden, und die Konvention erlaubt die alte Form.
- **Keine Template-Änderung** im vendored Baseline-Baum (derivativ).
- **Kein Produkt-Code.**

## 4. Definition of Done

- [ ] [ADR-0050](../../adr/0050-fence-unclosed-in-spans.md)/[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md) mit Kennungen in `Schärft:`/`Bezug:` + Geschichte-Zeile;
      `make adr-check` stumm (Probe dokumentiert).
- [ ] ADR-Index-Kopf mit Beispiel beider Formen; Reviewer-Skill 1.6.0 mit
      Anker.
- [ ] `make gates` grün; unabhängiger Review; Closure-Notiz; Register
      gesichtet.

## 5. Abnahme-Punkte / Risiken

- **[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md) ist gleichzeitig Gegenstand laufender Fortschreibungen** (zuletzt
  slice-112); ein Nachzug der Felder darf keine Entscheidung umschreiben — nur
  Adress-Form. — **Ausgang:** *(bei Closure)*
- **Reviewer-Anker als MEDIUM** könnte Alt-ADR-Zitate in Reviews als Befund
  lesen — der Anker nennt ausdrücklich die Zwei-Formen-Regel. — **Ausgang:**
  *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): slice-114 **und** slice-115 in `done/`
(Kennungen existieren in beiden Straten).

**Rückführungen:** `in-progress` → `next`, falls ein `Proposed`-ADR vorher
`Accepted` wird (dann entfällt sein Nachzug — Konvention, kein Bruch).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** ADR-Bestand (`docs/plan/adr/`, GF), Harness-Skills
  (`.harness/skills/`, GF), Planungs-Haus-Stil (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-22): BEO-002
  ([`MR-025`](../../../../harness/conventions.md#mr-025)), BEO-006/009 Arbeitsregeln; nichts Einschlägiges darüber hinaus.

Slice-ID: slice-116. Betroffene IDs:
[ADR-0050](../../adr/0050-fence-unclosed-in-spans.md),
[ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
(Proposed), [`MR-025`](../../../../harness/conventions.md#mr-025). Module:
ADR-Index, Reviewer-Skill, Slice-Haus-Stil. Gates: `make adr-check` (Probe),
`make doc-check`, `make gates`; Wellen-Closure `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Anwendung einer Konvention an eigenen
Artefakten; kein Legacy-Import.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
