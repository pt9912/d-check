# slice-267: `vcs` lässt einen reinen Pfad-Nachzug in immutablen Dateien durch

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** Baseline `v6.17.0` · `regelwerk/modul-04-adrs.md` §Nachzug ist
keine Überschreibung (ein Referenz-/Pfad-Nachzug in einer `Accepted`-ADR ist
keine inhaltliche Überschreibung); Anlass: der Umzug von
`docs/user/releasing.md` (slice-268) braucht ihn für zwei `Accepted`-ADRs;
Auftraggeber-Entscheid 2026-10-09.

**Berührte Spec-Stellen:** [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in),
`AGENTS.md` §3.5.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Das Modul `vcs` erkennt eine Änderung an einer immutablen Datei, die
nur das **Ziel** eines Markdown-Links ändert und Linktext wie übrige Zeile
gleich lässt, als Pfad-Nachzug und meldet sie nicht — opt-in, ohne den
Schlüssel byte-identisch. Dieses Repo schaltet es für die ADRs ein;
`AGENTS.md` §3.5 nennt den Nachzug neben `## Geschichte` und dem
Status-Übergang.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Umzug von `releasing.md`** — slice-268.
- **Ein Nachzug, der den Linktext oder andere Worte der Zeile ändert** — das
  wäre Inhalt; die Grenze liegt beim Ziel in den Klammern.
- **Der Template-Feld-Nachzug** (zweiter Fall im Regelwerk) — anderer
  Gegenstand, kein Anlass.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] `vcs` lässt einen reinen Link-Ziel-Nachzug durch (Schlüssel, Lastenheft,
      Spezifikation, Tests, die ohne die Änderung aus dem richtigen Grund rot
      sind); eine Änderung am Linktext oder an der übrigen Zeile bleibt ein
      Befund.
- [ ] Ohne den Schlüssel ist die Ausgabe unverändert; `.d-check.yml` schaltet
      ihn für die ADRs ein; `AGENTS.md` §3.5 nachgezogen; `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Kern-Regel `vcs` | update | Pfad-Nachzug erkennen |
| Konfig-Modell, YAML-Adapter, `--print-config`-Vorlage | update | Schlüssel |
| `spec/lastenheft.md`, `spec/spezifikation.md` | update | Anforderung und Verfeinerung |
| `.d-check.yml`, `AGENTS.md` §3.5 | update | Nutzung in diesem Repo |

## 4. Trigger

**Start** (`next` → `in-progress`): slice-265 in `done/`; `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): der Diff-Vergleich im Modul kennt keine
  Zeilen-Paarung, die einen Nachzug von einer Ersetzung unterscheidet — dann
  zuerst die Paarung.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft. Produkt-Verhalten — geht mit dem nächsten Release
hinaus.

## 6. Risiken und offene Punkte

- **Nachzug als Tarnung** — eine Link-Ziel-Änderung kann die Aussage einer
  Entscheidung verschieben, wenn das neue Ziel etwas anderes ist. Das Modul
  sieht nur die Form. — **Ausgang:** *(offen)*

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Produkt unter dem Default `*`
(`ALL`); beim Beanspruchen neu prüfen.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** beim Beanspruchen neu lesen.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
beim Beanspruchen aus dem jüngsten Lauf lesen.

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
