# Slice slice-234: `--suggest-config ai-harness` kennt die Randbedingungs-Reihe `-RB-`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** — **wellenlos**. Die Closure-Bedingung geht über die eigene DoD
nicht hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine
Welle braucht).

**Bezug:** [`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten)
(Konfigurations-Vorschlag, reservierte Quellen `ai-harness`/`ai-harness-init`),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(Determinismus). Change Request des Konsumenten `ai-harness-course`
(2026-09-27, Priorität niedrig, additiv). Der Beleg für die Reihe liegt dort in
einer noch **ungetaggten** Welle; die hier gepinnte Baseline kennt sie nicht.
d-check trägt die Form als Konsument dieses CR, nicht als adoptierte
Baseline-Regel.

**Berührte Spec-Stellen:**
[`DC-FA-CLI-006.a`](../../../../spec/spezifikation.md#dc-fa-cli-006a--konfigurations-vorschlag)
(reservierte Quellen: kanonisches Anforderungs-`ids`-Muster, Präfix-Ableitung).

**Verantwortlich:** pt9912.

**Autor:** claude-sonnet-5. **Datum:** 2026-09-27.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Ein Repo, dessen Lastenheft `<PREFIX>-RB-<NN>` führt, bekommt aus
`--suggest-config ai-harness[-init]` einen `ids`-Block, der die Reihe auf
Existenz und Link prüft; ohne `-RB-` im Repo bleibt der erzeugte Block
byte-gleich.

**Befund, gegen den Code gelesen — der CR stimmt in beiden Punkten.**
`internal/hexagon/core/app/suggest.go` trägt die Reihenfolge `FA-…|QA` an
zwei Stellen: `reqShape` (Präfix-Ableitung aus den Lastenheft-Headings) und
`harnessIDPatterns` (das kanonische Anforderungs-Muster). Eine
Randbedingungs-Überschrift wird in der Ableitung nicht als Anforderung
gelesen, und das erzeugte Muster matcht `<PREFIX>-RB-07` nicht — ein Verweis
auf eine nicht vorhandene Kennung bliebe still. Im **nicht**-reservierten
Modus (Quelle = Pfad) ist die Reihe nicht betroffen: die allgemeine
Kennungs-Gestalt matcht sie schon.

**Ein dritter Fundort, den der CR nicht nennt:** `reqIDFull` in
`internal/hexagon/core/app/trace.go` hat dieselbe Gestalt (`FA-…|QA`) und ist
der **Default** der Requirements Traceability Matrix. Er wird hier
**nicht** geändert (siehe Abgrenzung).

**Abnahme (aus dem CR, als Break-Test):**

- Lastenheft mit `### <PREFIX>-RB-01 — …` und einem Verweis auf `-RB-07`
  (existiert nicht): der erzeugte `ids`-Block enthält `RB`, der Lauf meldet
  die fehlende Kennung.
- Lastenheft nur mit `-FA-`/`-QA-`: der erzeugte Block ist **byte-gleich** zum
  Stand vor der Änderung.
- Mehrere verschiedene Präfixe: wie bisher Fehler, `--id-prefix` erforderlich
  — auch, wenn das zweite Präfix nur über eine `-RB-`-Überschrift ins Spiel
  kommt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`reqIDFull` (RTM-Default) erweitern** — anderer Vorgang mit eigener
  Folge: die RTM gehört zu
  [`DC-FA-CLI-009`](../../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix)
  und zählt jede Kennung ihrer Menge als Anforderung, die einen Slice braucht;
  `make completeness-check` meldete jede Randbedingung als **Waise**. Ob eine
  Randbedingung durch einen Slice „gebaut" wird, ist eine Fachfrage des
  Konsumenten, nicht eine Nebenwirkung eines Musters — und der Konsument kann
  die Menge heute schon per `id-pattern` erweitern. Diese Grenze steht in der
  Spezifikation als benannte Grenze, nicht als Stille.
- **Ein generisches `[A-Z]{2}` im Generator** — vom CR ausdrücklich
  nicht gewünscht: die Muster sind konventionsfest, ein offenes Muster
  vermischte die Reihen.
- **Weitere Kennungsklassen und die Verfeinerungs-Form `-FA-<NN>.<Buchstabe>`**
  — nicht Teil des CR; im Kurs offen.
- **Das statische Gerüst von `--print-config`** — es führt ein Beispiel für
  ein eigenes Repo (Präfix `DC`), keine Konvention; eine Randbedingungs-Reihe
  steht dort nicht zur Entscheidung. Bestand bleibt bewusst stehen.
- **Handbuch, CHANGELOG, Release** — Release-Prep-Vorgang, kein Feature-Commit
  (`AGENTS.md` §5). Der CR nennt beide; sie werden dort nachgezogen.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — höchstens drei Liefer-Punkte.

- [ ] **Vertrag:** `spec/lastenheft.md` (Konfigurations-Vorschlag: Beschreibung
      und Akzeptanzkriterien um die dritte Reihe, drei Break-Tests aus §1,
      Versions-Bump mit Historie-Zeile nach
      [MR-032](../../../../harness/conventions/MR-032-historie-vor-accepted.md)),
      `spec/spezifikation.md` (Präfix-Ableitung und kanonisches Muster;
      RTM-Grenze benannt). **Form nach dem Kanon** (`modul-03-spec.md`): die
      Historie-Zeile nennt **weder ADR noch Slice** (Decken-Regel), die
      Spezifikation nennt in keinem Abschnitt eine ADR oder einen Slice. Ob die
      RTM-Abgrenzung eine eigene ADR braucht (Entscheidung mit Alternativen:
      RTM erweitern, nicht erweitern, per `id-pattern` dem Konsumenten
      überlassen), wird **vor** dem ersten Edit entschieden und in §7 benannt.
- [ ] **Generator:** `reqShape` und das Anforderungs-Muster in
      `harnessIDPatterns` kennen `RB`; Tests: RB-Lastenheft (Muster enthält
      `RB`, fehlende `-RB-07` wird gemeldet — **Rot-Beleg gegen den alten
      Stand**), FA/QA-Lastenheft (Ausgabe byte-gleich, gegen den Stand vor der
      Änderung verglichen), mehrere Präfixe (Fehler).
- [ ] **Bestandsprobe:** `--suggest-config ai-harness-init` gegen dieses Repo
      und gegen die Fixtures ist vor und nach der Änderung byte-gleich
      (Vergleichs-Ausgabe in §7).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine
      Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | Erweiterung einer bestehenden Anforderung (kein neues Kürzel); Bump + Historie |
| `spec/spezifikation.md` (Konfigurations-Vorschlag) | update | Schritt „reservierte Quellen": dritte Reihe; benannte RTM-Grenze |
| `internal/hexagon/core/app/suggest.go` | update | zwei Muster: `reqShape` und Anforderungs-Zeile in `harnessIDPatterns` (Doc-Kommentar von `reqShape` nennt die Gestalt mit) |
| `internal/adapter/driving/cli/cli_acceptance_test.go` | update | drei Break-Tests nach den Akzeptanzkriterien; die Generator-Logik hat keine eigene Unit-Testdatei, ihre Belege sind die CLI-Akzeptanztests der `--id-prefix`-Kriterien |

**Reihenfolge der Vertrags-Änderung.** Das Lastenheft steht auf `Draft`: der
Kanon (`grundlagen-source-precedence.md`, *Wann die CR-Pflicht beginnt*) lässt
es vor `Accepted` frei änderbar und die Trennung von Entscheidung und Umsetzung
„greift noch nicht"; [MR-032](../../../../harness/conventions/MR-032-historie-vor-accepted.md)
verlangt Bump und Historie-Zeile trotzdem (`Verweis` bleibt `—`). Die
Änderung darf deshalb **im** Slice liegen; sie steht in einem eigenen Commit vor
dem Generator-Code, damit sie einzeln lesbar bleibt.

**Vor dem Editieren — Spiegel listen** ([MR-025](../../../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)):
die Gestalt `FA-…|QA` steht an vier Fundorten (`reqShape`,
`harnessIDPatterns`, `reqIDFull`, das statische Gerüst in
`config_template.go`); nur die ersten beiden ändern sich, die anderen zwei sind
in §1 als bewusst stehend begründet. Dazu Lastenheft-AKs und Handbuch-Beispiele,
die „FA/QA" aufzählen — Handbuch im Release-Prep.

## 4. Trigger

**Start** (`next` → `in-progress`): keine Abhängigkeit zu
[slice-232](../open/slice-232-links-ziel-hinter-zeilenumbruch.md) und
[slice-233](../open/slice-233-links-referenz-definitionen.md); alle drei
bumpen das Lastenheft, laufen aber nacheinander (WIP-Limit 1). Bei der
Beanspruchung entsteht der dritte Vorprüfungs-Block (Nachtlauf-Stand).

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): wenn sich zeigt, dass die Präfix-Ableitung
  mit `-RB-` mehr als die Gestalt ändert (etwa weil ein Verfeinerungs-Suffix
  oder ein Bereichs-Segment nach `RB` nötig wird).
- `in-progress` → `open` (blockiert): wenn der Kurs die Form der Reihe vor
  dem Tag ändert und der CR neu gefasst werden muss.

## 5. Closure-Trigger

DoD vollständig (§2) + Review-Report liegt vor + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Der Kurs hat die Reihe noch nicht getaggt; ändert er ihre Form, trägt d-check
  eine Kennungsreihe, die die Baseline nicht führt. **Ausgang:** bei Closure zu
  vergeben.
- Der Konsument liest „Randbedingungen im Generator" als „Randbedingungen in der
  RTM" und wundert sich, dass `completeness-check` sie nicht kennt.
  **Ausgang:** bei Closure zu vergeben (die benannte Grenze in der
  Spezifikation ist die Antwort; Rückfrage an den Konsumenten, falls er die RTM
  will).
- Ein Repo mit **zwei** Präfixen, von denen eines nur über `-RB-`-Überschriften
  auftaucht, wechselt vom stillen Ein-Präfix-Ergebnis zum Fehler „mehrdeutig".
  **Ausgang:** bei Closure zu vergeben (gewollt nach CR, aber eine sichtbare
  Verhaltensänderung für genau diesen Fall).

## 7. Closure-Notiz

*(Bei der Closure zu füllen.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:363-364 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine berührte Sub-Area: das Produkt
(`spec/`, `internal/`) unter dem Default `*` (Kürzel `ALL`); bereits
deklariert, keine Ausdifferenzierung nötig.

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:369-369 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(Stichworte Kennungsreihe, `--suggest-config`, `--id-prefix`,
`ai-harness-init`, `reqShape`): **keine Treffer**.

**Vorgelagert — Nachtlauf-Stand lesen** (bei der Beanspruchung, 2026-09-27):
`make nightly-state` meldet `upstream-drift.yml` **ROT** (Lauf
2026-09-27T06:03:41Z, derselbe Lauf wie bei `slice-237`s Beanspruchung,
unverändert), `image-scan.yml` **grün** (2026-09-26T08:38:12Z). Dieselben
vier planmäßigen Fremd-Release-Meldungen (`golangci-lint`, `semgrep`,
`a-check` VERALTET; `golang`-Basis-Digest ABWEICHEND), keine unerwarteten;
sie berühren dieses Modul (`suggest`) nicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

### Sub-Area: `*` (Produkt: Spezifikation und Code)

- **Modus:** GF
- **Konventionen-Dichte:** Hoch — Lastenheft und Spezifikation führen die
  reservierten Quellen samt Akzeptanzkriterien für den Präfix-Parameter.
- **Phase-Reife:** Phase 5.
- **Evidenz-/Diskrepanz-Risiko:** Niedrig — additiv, ohne `-RB-` byte-gleich;
  das Risiko liegt an der ungetaggten Quelle (§6).
- **Reconciliation-Aufwand:** Keiner — kein Brownfield-Bestand.
