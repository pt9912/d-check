# Slice slice-099: Modul `structure` — Implementierung

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`) — kein `Status:`-Feld; Wechsel nur per `git mv`
(Baseline-Regelwerk `modul-05-planning-harness.md`).

**Welle:** noch keiner Welle zugeordnet — **nicht** welle-69: deren
Out-of-Scope schließt die Implementierung ausdrücklich aus, und sie ist mit
[slice-096](../done/slice-096-structure-modul-analyse.md) geschlossen. Dieser
Slice ist ihr **derivatives** Ergebnis und bekommt seine Welle bei der
Eröffnung der Umsetzungs-Welle.

**Bezug:** [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
[`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Preset-Kopplung + `closure-note-ambiguous`),
[ADR-0049](../../adr/0049-structure-modul-schnitt-und-preset.md).

**Autor:** pt9912. **Datum:** 2026-08-09.

---

## 1. Ziel

Das Modul `structure` vollständig — Regel-Liste, Kandidaten-Menge,
Abschnitts-Findung mit beiden Kardinalitäts-Modi, Bereinigung und **alle sechs**
Bedingungen — plus die Preset-Kopplung der Closure-Fähigkeit samt
`closure-note-ambiguous`.

## 2. Warum in **einem** Slice, nicht in zweien

Der erste Zuschnitt trennte Kern und Marken/Zählung. Das ist an der
**Release-Grenze** gescheitert (Befund des Umsetzbarkeits-Reviews): das
veröffentlichte Schema führt alle Schlüssel, und die Config-Dekodierung ist
strikt — ein Release nach dem Kern allein wiese eine spezifikationskonforme
Konfiguration mit Exit 2 ab. Ein Modul, das die Hälfte seines eigenen Schemas
ablehnt, ist kein lieferbarer Zwischenstand.

## 3. Vorgehen

1. **Gemeinsame Mechanik zuerst herausziehen**, nicht kopieren:
   Abschnitts-Findung, Kardinalitäts-Behandlung, Fence-Bereinigung. Die
   Spezifikation weist die Closure-Fähigkeit als Preset aus — zwei Kopien
   könnten driften, ohne dass ein Test es merkt.
2. **Grund-Codes im Lockstep — es sind neun, nicht sieben:** die sechs
   Bedingungs-Codes (`section-empty`/`-thin`/`-oversized`/`-forbidden`/
   `-pattern-missing`/`-marker-missing`) **plus** die zwei Struktur-Codes
   (`section-missing`, `section-ambiguous`) **plus** `closure-note-ambiguous`.
   Alle neun gehören mit `AllReasons()`, den Doctor-Klartexten und der
   Spezifikations-§4-Tabelle in **denselben** Commit.
3. **CLI-Mit-Modifikation**, die jedes bisherige Modul mitgebracht hat:
   `--print-config`-Gerüst, `--suggest-config`-Vorlage und die
   `--print-mk`-Target-Liste.

## 3a. Spiegel dieser Semantik-Änderung ([`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten))

**Vor** dem ersten Editor aufgeschrieben, wie die Regel es verlangt. Dieser Slice
ändert die Grund-Code-Menge um **neun** Einträge und ist damit der Slice mit den
meisten Spiegeln bisher.

| Spiegel | berührt? | was genau |
|---|---|---|
| Anforderung [`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) | **nein** | liegt seit welle-69 vollständig vor — dieser Slice erfindet nichts |
| Anforderung [`DC-FA-PLAN-001`](../../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) | **ja** | die Mehrdeutigkeits-Härte (`closure-note-ambiguous`) ist zugesagt, aber **nicht** implementiert |
| Algorithmus [`DC-FA-STRUCT-001.a`](../../../../spec/spezifikation.md#dc-fa-struct-001a--struktur-invarianten-innerhalb-eines-dokuments-structure) | **nein** | acht Schritte liegen vor |
| Algorithmus [`DC-FA-PLAN-001.a`](../../../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning) | **ja** | Schritt C3 — die Mehrdeutigkeit |
| §2-Config-Schema | **nein** | `structure[]` steht vollständig |
| §4-Grund-Code-Tabelle | **ja, neun Zeilen** | sechs Bedingungen + `section-missing` + `section-ambiguous` + `closure-note-ambiguous` |
| `AllReasons()` / `reasonTexts()` | **ja, neun** | im **selben** Commit wie §4 (Lockstep) |
| `--print-config`-Vorlage | **ja** | `structure`-Gerüst |
| `--suggest-config`-Vorlage | **ja** | Aufnahme prüfen — situatives opt-in |
| `--print-mk`-Target-Liste | **ja** | **und** die **Zahl** in der Out-of-Scope-Zeile von [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben) |
| Benutzerhandbuch | **ja** | Modul-Tabelle (Zahl **und** Zeile), §5-Config, §11-Zeile, ggf. §4-Aufgabe |
| README (beide Sprachen) | **ja** | Status-Zeile (Zahl · Enumeration · „zuletzt das Modul X“) **und** Modul-Liste; DE zuerst |
| `operations.md` | **ja** | Modul-Enumeration der `--enable`/`--disable`-Zeile |
| `AGENTS.md` / `harness/README.md` | **nur falls** ein Gate-Target dazukommt | Gate-Beschreibungen |
| [ADR-0049](../../adr/0049-structure-modul-schnitt-und-preset.md) | **ja** | Status auf `Accepted` bei der Closure |

Nicht auf der Liste, weil unberührt: `.d-check.yml` (das eigene Repo aktiviert
`structure` erst nach eigener Messung), `.d-check.closure.yml`, die
Referenzmatrix.

### 3b. Bilanz der Spiegel-Liste — was sie **nicht** enthielt

Ein unabhängiger Review hat die Liste aus §3a gegen die Wirklichkeit geprüft.
Sie hat gewirkt (die **Zahl** der Targets in der Out-of-Scope-Zeile, der Intro-Satz beider READMEs
und `operations.md` standen darauf und wären sonst vergessen worden) — aber sie
war **selbst lückenhaft**, und das gehört in die erste Anwendung ihrer Regel:

| fehlender Spiegel | Folge |
|---|---|
| `Makefile` | das Closure-Gate musste ein zweites Modul aktivieren — nicht gelistet |
| `validModules()` | ein neues Modul muss in die Modul-Registry, sonst lehnt die CLI `--enable structure` ab |
| **Akzeptanzkriterien** derselben Anforderung [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben) | ich hatte Beschreibung **und** Out-of-Scope nachgezogen, die Kriterien nicht — dieselbe Stelle war in 0.37.1 schon einmal als Selbstwiderspruch saniert worden |
| [`DC-FA-CLI-006`](../../../../spec/lastenheft.md#dc-fa-cli-006--konfigurations-vorschlag-aus-autoritäts-dokumenten) + `--suggest-config` | die Enumeration der nicht aktivierten situativen Module |

Falsch **eingetragen** war außerdem eine Zeile: `.d-check.closure.yml` stand
ausdrücklich als „unberührt“ auf der Liste und wurde am Ende um 37 Zeilen
erweitert.

**Eine bestätigende Re-Review fand eine fünfte Lücke** — und die ist die
interessanteste: die Modul-Enumeration der emittierten `--print-config`-Vorlage
lag in **genau der Datei**, die ohnehin bearbeitet wurde. Sie stand deshalb nicht
auf der Liste und blieb doch stehen; sie war sogar schon vorher veraltet
(`citations` fehlte seit seiner Einführung). Der Datei-`grep` hätte sie nicht
gefunden. Die Antwort ist kein längerer Listeneintrag, sondern eine **Bindung**:
ein Test vergleicht die Zeile jetzt gegen die Modul-Registry. Wo eine Aufzählung
eine Menge spiegelt, gehört sie an ihre Quelle gebunden — die Liste ist der
Notbehelf für alles, was sich nicht binden lässt.

**Die Lehre ist nicht „die Regel taugt nicht“, sondern „die Liste ist
selbst ein Artefakt“:** sie war aus dem Gedächtnis geschrieben, nicht aus dem
Repo abgeleitet. Drei der vier Lücken hätte ein `grep` nach dem vorigen
Modul-Namen gefunden. Diese Erfahrung gehört in
[`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten).

**Ein sechster Spiegel derselben Bauart, gefunden erst vor dem Tag: das Datum.**
Die Welle läuft über zwei Kalendertage; jeder Datumsstempel dieses Slice — das
Wellendokument, die Roadmap-Zeile,
[`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten),
die Review-Reports samt Dateinamen,
`CHANGELOG`, `version.md`, der Handbuch-Kopf und die §11-Zeile, die Historien
von Lastenheft und Spezifikation — trug den Tag des **vorigen** Releases, weil
er von der Zeile darüber abgeschrieben war. Der Fall ist die Spiegel-Klasse in
Reinform, nur ohne Fachbezug: eine Angabe, die an zwölf Stellen wiederholt und
an keiner geprüft wird. Nachgezogen als eigener Punkt der
[Release-Prep-Checkliste](../../../user/releasing.md#release-prep-vor-dem-tag);
die Spezifikations-Historie war zusätzlich falsch **einsortiert** (dritte statt
erste Zeile) — die Fehldatierung hatte sie plausibel aussehen lassen.

## 4. Definition of Done

- [x] Modul vollständig (beide `sections`-Modi, alle sechs Bedingungen,
      Kandidaten-Menge, fail-closed Ränder inkl. Leerlauf trotz `exempt-paths`);
      Abschnitts-Mechanik **geteilt** mit der Closure-Fähigkeit — belegt durch
      den Preset-Kopplungs-Test (dieselbe Eingabe, beide Oberflächen, gleiche
      Befund-Zeilen).
- [x] **Neun** neue Grund-Codes im Lockstep mit `AllReasons()` und Spezifikation
      §4; jedes Akzeptanzkriterium als Test, insbesondere die drei Marken-Formen,
      `sections: each` und „Mehrdeutigkeit schlägt Messung".
- [x] CLI-Enumerationen nachgezogen — einschließlich der Target-Liste in
      [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben),
      deren Out-of-Scope-Zeile die **Zahl** der Targets festschreibt und daher
      mitzuändern ist. `make gates` grün; **Paritäts-Beleg** gegen die
      beigezogenen Adopter-Fixtures; Release als **Minor** (neues Modul +
      additiver Code — ein Repo mit zwei Closure-Abschnitten wird danach rot,
      das gehört in die Release-Notiz).

## 4a. Messung: welche Regeln aktivieren wir selbst?

Vor der Aufnahme in die eigene Konfiguration gemessen, je Kandidat ein Lauf —
dieselbe Disziplin wie bei der Floskel-Liste und der Platzhalter-Erkennung.

| Kandidaten-Regel | Befunde | aufgenommen |
|---|---|---|
| slice-Closure-Notiz nicht leer | **0** | **ja** |
| slice-Definition-of-Done vorhanden | **0** | **ja** |
| Wellen-Ergebnisnotiz: „Was wurde geliefert?“ | **0** | **ja** |
| slice-DoD **ohne offene Tasks** | 32 | **nein** |
| ADR: „Fitness Function“ | 15 | nein |
| ADR: „Re-Evaluierungs-Trigger“ | 20 | nein |
| ADR: „Verglichene Alternativen“ | 14 | nein |
| Wellen-Ergebnisnotiz: „Steering-Loop-Einträge“ | 7 | nein |

**Die 32 sind das lehrreiche Nein.** Sie sehen aus wie ein Bestandsschaden, sind
aber keiner: ein abgeschlossener Slice **darf** eine offene Box tragen, wenn die
**Welle** sie einlöst — [slice-094](../done/slice-094-closure-zaehl-paritaet.md)
tut genau das, mit ausdrücklicher Begründung in seiner DoD. Die Regel wäre am
ersten Tag falsch gewesen; ohne Messung hätte ich sie aktiviert.

Die vier ADR-/Wellen-Regeln melden **echte** Lücken: ältere Dokumente kennen die
Abschnitte noch nicht. Sie zu sanieren ist eigene Arbeit, kein Nebeneffekt dieses
Slice — die Regeln bleiben deshalb draußen, **benannt** statt vergessen.

**Bindepunkt:** die drei aufgenommenen Regeln laufen über das
**Closure-Profil**, nicht im inneren Loop. Sie fragen nach dem Ruheort, und das
ist dieselbe Begründung, mit der [ADR-0048](../../adr/0048-closure-note-struktur-im-planning-modul.md)
die Closure-Note-Struktur dorthin gelegt hat.

## 5. Risiken / offene Punkte

- **Das Herausziehen berührt ausgelieferten Code.** Die Closure-Fähigkeit ist
  seit v0.52.0 draußen; ein Refactor ihrer Abschnitts-Logik darf ihren Befundsatz
  nicht verändern. — **Ausgang:** offen; der Preset-Kopplungs-Test ist die
  Absicherung.
- **Die Marken-Syntax ist eine Konvention, keine Norm.** `**M:**` ist die Form
  der beiden vermessenen Repos; ein drittes schreibt vielleicht anders.
  — **Ausgang:** offen; zu entscheiden, ob die Marke konfigurierbar wird oder
  die Konvention Teil der Zusage bleibt.
- **Die Paritäts-Fixtures liegen im Adopter-Repo**, nicht hier. — **Ausgang:**
  offen; beizuziehen, nicht nachzubauen.
- **Die Fence-Lexik trägt einen bekannten Defekt**
  ([slice-101](../done/slice-101-fence-unbalanciert.md)): ein ungeschlossener Fence
  verschluckt still den Rest. `structure` erbt ihn über die geteilte Mechanik.
  — **Ausgang: entschieden.** Der Fix läuft zuerst; §6 macht ihn zur **bindenden**
  Start-Bedingung, nicht zur Empfehlung.

## 6. Trigger

**Start** (`next` → `in-progress`): [slice-096](../done/slice-096-structure-modul-analyse.md)
in `done/` **und** [slice-101](../done/slice-101-fence-unbalanciert.md) in `done/`;
WIP-Slot frei. Die zweite Bedingung ist **bindend**, nicht bevorzugt: liefe
dieser Slice zuerst, erbte das neue Modul über die geteilte Mechanik einen
**bekannten** stillen Grün-Pfad — und läge damit ab Release als zugesagte
Fähigkeit vor, die ihre Zusage an einer bekannten Stelle nicht hält.

**Rückführungen:** `in-progress` → `next`, falls das Herausziehen der
gemeinsamen Mechanik einen eigenen Refactor-Slice verlangt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Code (`internal/`) und Spec (`spec/`), beide unter
  dem Repo-Default GF (`harness/conventions.md` §Modus: `*`).
- **Offene Beobachtungen sichten** (bei Slice-Beginn erneut gelesen — das
  Register führt inzwischen **vier** Einträge): **BEO-002** ist seit 2026-08-15
  **verkörpert** als [`MR-025`](../../../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten);
  die Spiegel-Liste steht in §3a und ist damit erledigt, statt eine Warnung zu
  bleiben. **BEO-003** (geteilte Lexik driftet an den Rändern) ist **einschlägig
  als Auftrag**: Schritt 1 des Vorgehens zieht die Abschnitts-Mechanik heraus,
  statt sie zu kopieren — und braucht je Konsument eine Assertion, nicht nur ein
  geteiltes Prädikat. **BEO-004** (Modul-Grenze nur auf der Quell-Achse) ist
  einschlägig als **Frage**: `structure` benennt seine Eingabe selbst und ist
  damit genau die Klasse. Ferner **BEO-001** (andere
  Klasse — Referenz zwischen Dokumenten statt Form innerhalb eines; in
  [slice-096](../done/slice-096-structure-modul-analyse.md) ausdrücklich
  als Nicht-Ziel festgehalten) und **BEO-002** (Semantik-Änderungen werden nur im
  Dokumentkörper nachgezogen). BEO-002 **betrifft diesen Slice unmittelbar**: er
  ändert die Grund-Code-Menge und damit gleich mehrere Spiegel — `AllReasons()`,
  Doctor-Klartexte, Spezifikation §4, drei CLI-Enumerationen und die
  Handbuch-Modulliste. Sie gehören vor dem Editieren aufgelistet.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Anforderung, Algorithmus und ADR liegen
bereits; dieser Slice liefert, was sie versprechen.

## 9. Closure-Notiz (nach `done/`)

_Ausstehend._
