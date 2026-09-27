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

**Verantwortlich:** pt9912.

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

- [x] **Reviewer-Agent** unter `.claude/agents/`: Rolle, Eingang, Ausgang,
      Zeiger auf `.harness/skills/reviewer.md`, Kontext-Trennung; Größe in Bytes
      in §7.
- [x] **Verifier-Agent** unter `.claude/agents/`: dieselbe Form. **Vor dem Anlegen
      entschieden:** wo der Verifikations-Bericht liegt, ohne dass die
      Review-Deckung (`make review-coverage`) ihn als Report eines Slice
      missdeutet oder den Report-Zähler verschiebt — belegt an einer Probe, nicht
      vermutet.
- [x] **`implement-slice`:** Schritt 11 gleicht `AGENTS.md` §5 an; die Übergaben
      an Reviewer und Verifier stehen als Zeiger auf `AGENTS.md` §6 und die
      beiden Typen, ohne neue Regel.
- [x] Jede Datei besteht die Prüffrage aus §1; der unabhängige Review geht sie
      **Satz für Satz** durch und führt sie in seinem Report.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8, kein
      Self-Review. Löst der neue Agent-Typ in der laufenden Sitzung nicht auf,
      steht in §7, unter welchem Typ der Review lief.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — oder „keine
      Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
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
  **Ausgang: eingetreten, weiter offen** →
  [`BEO-ALL/agent-typ-loest-nicht-in-derselben-sitzung-auf`](../observations/BEO-ALL/agent-typ-loest-nicht-in-derselben-sitzung-auf/observation.md).
  Genau eingetreten (R1-Aufruf mit `subagent_type: "reviewer"` scheiterte),
  löste sich aber innerhalb derselben Sitzung ohne Repo-Änderung von selbst
  (R2-Aufruf erfolgreich) — kein Folge-Slice, da nichts im Repo zu ändern
  ist; als Beobachtung für künftige Implementer festgehalten.
- Die Dateien wachsen beim nächsten Anfassen wieder (Herleitungen, Zahlen,
  Kandidatenläufe) — die Klasse, die in den Schwester-Repos zu 17 und 28 KB
  geführt hat. **Ausgang: weiter offen** →
  [`BEO-ALL/duennes-werkzeug-artefakt-waechst-beim-naechsten-anfassen`](../observations/BEO-ALL/duennes-werkzeug-artefakt-waechst-beim-naechsten-anfassen/observation.md).
- Ein Verifikations-Bericht in `docs/reviews/` verschiebt die Review-Deckung.
  **Ausgang: entfallen.** Zweifach unabhängig am Code verifiziert
  (`hasMatchingReview` in `internal/hexagon/core/rules/reviews.go` prüft nur
  Dateinamens-Substring, nicht Inhalt/Anzahl/Suffix) und durch einen neuen,
  dauerhaften Test (`TestReviewsVerifierNamedReportSatisfiesCoverage`)
  belegt.

## 7. Closure-Notiz

**Geliefert:** zwei dünne Claude-Code-Agent-Träger unter `.claude/agents/`
(`reviewer.md` 2257 Byte, `verifier.md` 2627 Byte, nach den Nachzügen aus R2)
für die Kontext-Trennung von Reviewer und Verifier — beide verweisen auf
ihren Anweisungssatz (`.harness/skills/reviewer.md` bzw. das Baseline-
Regelwerk) statt Inhalt zu duplizieren. `.claude/commands/implement-slice.md`
Schritt 11 nennt CHANGELOG nicht mehr im Feature-Commit (`AGENTS.md` §5);
ein neuer Schritt 13 zeigt auf den Reviewer-/Verifier-Handoff.

**Agent-Typ-Auflösung (DoD-Pflicht):** Der erste Review-Aufruf mit
`subagent_type: "reviewer"` schlug fehl — der Typ löste in der laufenden
Sitzung noch nicht auf, der Review lief unter einem generischen Typ (R1). Ein
späterer Aufruf in **derselben** Sitzung (R2) löste denselben Typ bereits
erfolgreich auf, ohne dass sich am Repo etwas geändert hätte. Beide Reviews
liegen unter `docs/reviews/` vor.

**Zwei Review-Runden, je ein echter Fund.** R1 fand F-1 (HIGH,
merge-blockierend): zwei Sätze — wortgleich aus dem Referenz-Repo
`pg-change-feed` übernommen — bestanden die vom DoD verlangte Satz-für-
Satz-Prüfung nicht (keine Verankerung in einer gerankten Quelle dieses
Repos), sowie F-2 (MEDIUM): ein falscher Abschnitts-Anker in `verifier.md`.
Beide behoben. R2 verifizierte die Korrektur mit einer vollständigen
Neu-Prüfung **beider** Dateien (nicht nur der beiden benannten Stellen) und
fand zwei kleine Nachzüge (F-3 LOW: ein präziserer Anker war verfügbar; F-4
INFO: ein redundanter Satz) — beide sofort behoben, nicht mehr blockierend.

**Bestätigt, unabhängig verifiziert:** die Verifier-Ablage-Entscheidung
(Dateiname mit `-verify`-Suffix im selben `docs/reviews/`-Verzeichnis)
verfälscht `make review-coverage` nicht — `hasMatchingReview` prüft nur, ob
irgendein Dateiname die `slice-<NNN>`-Kennung trägt, nicht Inhalt, Anzahl
oder Suffix; ein neuer, dauerhafter Test belegt es. Die
[MR-025](../../../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)-
Spiegel-Prüfung (`AGENTS.md` §6, `harness/README.md` §Guides) bestätigt: kein
Nachzug nötig — beide Dokumente nennen auch `.claude/commands/
implement-slice.md` selbst nicht, Claude-Code-spezifische Werkzeug-Dateien
werden dort grundsätzlich nicht gespiegelt.

**Steering-Loop-Einträge:** zwei neue Beobachtungen —
[`BEO-ALL/agent-typ-loest-nicht-in-derselben-sitzung-auf`](../observations/BEO-ALL/agent-typ-loest-nicht-in-derselben-sitzung-auf/observation.md)
(1×, `offen`) und
[`BEO-ALL/duennes-werkzeug-artefakt-waechst-beim-naechsten-anfassen`](../observations/BEO-ALL/duennes-werkzeug-artefakt-waechst-beim-naechsten-anfassen/observation.md)
(1×, `offen`).

**Risiko-Ausgänge:** ein Risiko *entfallen* (Review-Deckung), eines
*eingetreten* mit Register-Verweis, eines *weiter offen* mit
Register-Verweis — siehe §6.

**Register-Sichtung bei Planung:** §8 nannte zwei Treffer
(`zustellkanal-haengt-an-werkzeugweg`,
`kommentar-traegt-herkunfts-prosa-statt-fuenf-klassen`); keiner davon trat in
diesem Slice als neue Instanz auf (thematisch verwandt, aber kein neuer
Beleg fällig).

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

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung, 2026-09-27):
`make nightly-state` meldet `upstream-drift.yml` **ROT** (Lauf
2026-09-27T06:03:41Z, unverändert seit slice-232/233), `image-scan.yml`
**grün** (Lauf 2026-09-27T09:18:12Z). Derselbe veraltete Stand: vier der
fünf gemeldeten Fremd-Release-Stände sind bereits gehoben, die Kurs-Baseline
(`v6.9.0` → `v6.10.0`) bleibt bewusst zurückgestellt bis `v6.11.0`. Keiner
der fünf Punkte berührt dieses Repo-Werkzeug (`.claude/`).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

### Sub-Area: `*` (Harness-Werkzeug-Dateien)

- **Modus:** GF
- **Konventionen-Dichte:** Mittel — der Kanon führt Rollen und Artefaktklassen;
  eine Konvention für Claude-Code-Agent-Dateien führt d-check nicht.
- **Phase-Reife:** Phase 3 (Durchsetzungsschicht steht, Rollen-Träger fehlen).
- **Evidenz-/Diskrepanz-Risiko:** Mittel — Attrappen- und Wachstums-Risiko (§6),
  belegt an zwei Schwester-Repos.
- **Reconciliation-Aufwand:** Keiner.
