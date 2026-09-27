# Slice slice-235: Claude-Code-Agents für Reviewer und Verifier; `implement-slice` an `AGENTS.md` angleichen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** — **wellenlos**. Die Closure-Bedingung geht über die eigene DoD
nicht hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine
Welle braucht).

**Bezug:** — kein `DC-*`, keine aktive ADR. Harness-Werkzeug-Dateien unter
`.claude/`; Gate-Grenze dazu:
[MR-048](../../../../harness/conventions/MR-048-gate-ueber-werkzeug-datei.md).
Kanon: `modul-08-agentenrollen.md` §Welche Rolle braucht welche Artefaktklasse,
`grundlagen-source-precedence.md` §Source Precedence (Absatz *Vollständigkeit*).

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** claude-sonnet-5. **Datum:** 2026-09-27.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Der Handoff an Reviewer und Verifier (`AGENTS.md` §6) läuft über
zwei benannte Agent-Typen in frischem Kontext statt über einen von Hand
geschriebenen Auftrag, und der bestehende `implement-slice` widerspricht
`AGENTS.md` nicht mehr.

**Befund.** (1) Es gibt in d-check keinen `.claude/agents/`; die Reviews dieser
Sitzung liefen unter `general-purpose` mit langem Ad-hoc-Auftrag. (2) Bei
`slice-231` ist der Verifier-Schritt ausgefallen — es gab keinen Träger, der ihn
anfordert. (3) `.claude/commands/implement-slice.md` Schritt 11 nennt
„CHANGELOG if required"; `AGENTS.md` §5 verlangt ihn **nicht im Feature-Commit**.

**Keine Chronik, keine Waisen.** Die neuen und geänderten Dateien tragen **keinen
Chronik-Satz** im Sinn der Form aus
[slice-231 §1](../done/slice-231-agents-md-chronik-streichen.md) (Genese, frühere
Fassung, Messvorfall als Herleitung, Planungs-Verweis, „seit slice-N" im
Fließtext) — und legen **nichts fest**: nach dem Kanon dürfen Werkzeug-Dateien
verweisen, ausführen und einen gerankten Ablauf ausbuchstabieren, aber keine
Regel einführen. Prüffrage je Satz: *Steht diese Aussage auch in einer gerankten
Quelle?* Nein ⇒ sie gehört nicht in die Datei.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Agents für Planner, Architect, Implementer und Validator** — nach dem Kanon
  Attrappen: deren Urteil läuft über Template bzw. Briefing, oder die
  Prüfgrundlage liegt im Slice; eine Zusatzdatei trägt keinen nicht-ableitbaren
  Inhalt und kann nur driften. Reviewer (Skill) und Verifier (Träger für den
  ausgefallenen Schritt) tragen dagegen die **Kontext-Trennung**, nicht Inhalt.
- **`plan-welle` und `close-welle` als Commands** — das Verfahren steht in
  `modul-06` und ist als Symlink in jedem Lauf geladen; eine Kopie wäre doppelte
  Pflege, und es ist keine Welle offen. Kein Folge-Slice; entsteht der Bedarf,
  ist das ein eigener Vorgang.
- **Rollen-Telemetrie** (`span-report` u. ä.) — existiert in d-check nicht; ein
  Verweis darauf wäre ein Phantom-Target.
- **Die Größe von `AGENTS.md` und die Regelwerk-Symlinks unter `.claude/rules/`**
  — anderer Gegenstand, offene Entscheidung des Auftraggebers.
- **Ein `model`-Pin** — bindet an eine Kennung, die in der ausführenden Umgebung
  fehlen kann.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — höchstens drei Liefer-Punkte.

- [ ] **Reviewer-Agent** unter `.claude/agents/`: Rolle, Eingang, Ausgang,
      Zeiger auf `.harness/skills/reviewer.md`, Kontext-Trennung; Größe in Bytes
      in §7.
- [ ] **Verifier-Agent** unter `.claude/agents/`: dieselbe Form. **Vor dem Anlegen
      entschieden:** wo der Verifikations-Bericht liegt, ohne dass die
      Review-Deckung (`make review-coverage`) ihn als Report eines Slice
      missdeutet oder den Report-Zähler verschiebt — belegt an einer Probe, nicht
      vermutet.
- [ ] **`implement-slice`:** Schritt 11 gleicht `AGENTS.md` §5 an; die Übergaben
      an Reviewer und Verifier stehen als Zeiger auf `AGENTS.md` §6 und die
      beiden Typen, ohne neue Regel.
- [ ] Jede Datei besteht die Prüffrage aus §1; der unabhängige Review geht sie
      **Satz für Satz** durch und führt sie in seinem Report.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8, kein
      Self-Review. Löst der neue Agent-Typ in der laufenden Sitzung nicht auf,
      steht in §7, unter welchem Typ der Review lief.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — oder „keine
      Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      wellenlos hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.claude/agents/reviewer.md` | neu | dünner Träger für die Kontext-Trennung; der Anweisungssatz bleibt der Skill |
| `.claude/agents/verifier.md` | neu | Träger für den Schritt, der bei `slice-231` ausfiel; Ablage-Entscheid siehe §2 |
| `.claude/commands/implement-slice.md` | update | Schritt 11 an `AGENTS.md` §5; Übergaben als Zeiger auf §6 |

**Vor dem Editieren — Spiegel listen**
([MR-025](../../../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)):
`AGENTS.md` §6 und `harness/README.md` §Guides nennen Reviewer-Skill und
Rollen; ob sie die neuen Typen nennen müssen, wird gelesen, nicht vermutet.
Die Agent-Dateien liegen im Scan-Bereich von `doc-check`: Kennungen darin sind
Links oder entfallen.

## 4. Trigger

**Start** (`next` → `in-progress`): keine Abhängigkeit. Bei der Beanspruchung
entsteht der dritte Vorprüfungs-Block (Nachtlauf-Stand).

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): wenn der Verifikations-Bericht einen eigenen
  Ablage- und Gate-Entscheid verlangt, der über die Probe hinausgeht.
- `in-progress` → `open` (blockiert): wenn der Auftraggeber die
  Symlink-Entscheidung vorzieht und sie die Rollen-Dateien berührt.

## 5. Closure-Trigger

DoD vollständig (§2) + Review-Report liegt vor + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Ein neu angelegter Agent-Typ löst in der laufenden Sitzung nicht auf.
  **Ausgang:** bei Closure zu vergeben.
- Die Dateien wachsen beim nächsten Anfassen wieder (Herleitungen, Zahlen,
  Kandidatenläufe) — die Klasse, die in den Schwester-Repos zu 17 und 28 KB
  geführt hat. **Ausgang:** bei Closure zu vergeben.
- Ein Verifikations-Bericht in `docs/reviews/` verschiebt die Review-Deckung.
  **Ausgang:** bei Closure zu vergeben.

## 7. Closure-Notiz

*(Bei der Closure zu füllen.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:363-364 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: die
Harness-Werkzeug-Dateien unter `.claude/` fallen unter den Default `*`
(Kürzel `ALL`); bereits deklariert, keine Ausdifferenzierung nötig.

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:369-369 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
**Treffer:**
[`zustellkanal-haengt-an-werkzeugweg`](../observations/BEO-ALL/zustellkanal-haengt-an-werkzeugweg/observation.md)
(Stand offen) — eine Anweisung, die am Werkzeugweg hängt statt am Inhalt; ein
Agent-Typ ist genau so ein Zustellkanal. Und
[`kommentar-traegt-herkunfts-prosa-statt-fuenf-klassen`](../observations/BEO-ALL/kommentar-traegt-herkunfts-prosa-statt-fuenf-klassen/observation.md)
(1×) für die Chronik-Freiheit der Dateien.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

### Sub-Area: `*` (Harness-Werkzeug-Dateien)

- **Modus:** GF
- **Konventionen-Dichte:** Mittel — der Kanon führt Rollen und Artefaktklassen;
  eine Konvention für Claude-Code-Agent-Dateien führt d-check nicht.
- **Phase-Reife:** Phase 3 (Durchsetzungsschicht steht, Rollen-Träger fehlen).
- **Evidenz-/Diskrepanz-Risiko:** Mittel — Attrappen- und Wachstums-Risiko (§6),
  belegt an zwei Schwester-Repos.
- **Reconciliation-Aufwand:** Keiner.
