# Review-Report: slice-130 — Historie-Form auf vier Spalten, eigene Strenge deklariert

**Datum:** 2026-08-23 · **Review-Art:** Doku-/Konventions-Review (geprüft gegen Slice-Plan
slice-130 §1/§2/§4/§5/§7, Baseline-Regelwerk `grundlagen-source-precedence.md`
§Spec-Stratifizierung, `harness/conventions.md` §Purpose, `MR-031`, `AGENTS.md` §3.4,
`.d-check.yml` `structure`-Block, Beobachtungs-Register `BEO-002`/`BEO-009`/`BEO-011`),
unabhängiger Reviewer ohne Anteil an der Arbeit

**Gegenstand:** `git diff e41850d..HEAD` — zwei Commits: `951a9e1` (Lifecycle-Move
`open/` → `in-progress/`, MR-013), `9987cf7` (vierte Historie-Spalte im Lastenheft +
Strenge-Deklaration)

**Skill:** `.harness/skills/reviewer.md` @ 1.9.0 (`5331466`) · **Modell-ID:** `claude-sonnet-5`

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-130-lastenheft-historie-form.md` (§1–§5, §7)
- `spec/lastenheft.md` Kopf (Version/Status-Feld) und §7 Historie; `spec/spezifikation.md` §7;
  `spec/architecture.md` (§Charakter, „keine Historie")
- Vendorte Vorlagen `.harness/baseline/v5.11.0/templates/spec/lastenheft.template.md`,
  `spezifikation.template.md`
- Kanon `.harness/baseline/v5.11.0/regelwerk/grundlagen-source-precedence.md`
  §Spec-Stratifizierung (CR-Pflicht am Status, Tatsachenberichtigung, Zurückgezogene
  Anforderungen, „Welche Stelle der Version steigt …")
- `harness/conventions.md` §Purpose, §Adaptions-Block, `MR-000`;
  `harness/conventions/MR-031-schritt-3-benennen.md`
- `.d-check.yml` `structure:`-Block (Regeln auf `spec/lastenheft.md` §7 und
  `spec/spezifikation.md` §7)
- `docs/plan/planning/observations.md` (`BEO-011`, `BEO-002`, `BEO-009`)
- `AGENTS.md` §3.4
- Nicht erhalten: die DoD-Abhakung (Verifikations-Rolle, getrennter Kontext)

**Vom Reviewer selbst gefahren** (Exit je Lauf direkt in Datei umgeleitet und gelesen,
`BEO-007`): `make doc-check` Exit 0 (463 Datei(en), 0 Befund(e)) · `make completeness-check`
Exit 0 (48 Anforderung(en), 0 Waise(n)) · `make verify-closure-notes` Exit 0 (427 Datei(en),
0 Befund(e)) · Volltext-Diff `git diff e41850d..HEAD -- spec/lastenheft.md` (Hunk-Zählung) ·
Zeilen-/Zellen-Zensus über alle 95 Historie-Datenzeilen (`awk`) · Byte-Vergleich Header-/
Separator-Zeile gegen `lastenheft.template.md:152-153` (`diff`, Exit 0) · Grep über
`spec/architecture.md` und alle `spec/*.md` nach `^## .*Historie` · Grep nach
„Tatsachenberichtigung“ im ganzen Baum (außerhalb Baseline/Planung).

**Verdikt: blockierend** — kein HIGH, zwei MEDIUM, kein LOW, kein INFO.

---

## Findings

### F-1 — Die „bewusst strenger“-Erklärung ist eine Adaption ohne Adaptions-Eintrag; die Selbst-Abgrenzung von MR-031 hält der eigenen Begründung nicht stand

- **kategorie:** MEDIUM
- **quelle:** `harness/conventions/MR-031-schritt-3-benennen.md` (Präzedenz: „Wer ihn
  verschärft, weicht von einem Baseline-Default ab — auch wenn er nur mehr verlangt“) ·
  `grundlagen-source-precedence.md` §Spec-Stratifizierung, letzter Absatz vor
  §ID-Schema als Klammer („Welche Stelle der Version steigt, entscheidet das Repo und
  gehört in den Adaptions-Block von `harness/conventions.md` — wie die Rangwahl innerhalb
  des Vertrags-Stratums“) · `harness/conventions.md` §Purpose („Default-Ort für:
  Adaptionen ggü. der Baseline“) · `modul-02-harness-bootstrap.md` §Freshness-Audit
  („Der Review geht durch die Adaptions-Liste, nicht nur durch den Diff“)
- **pfad:** `spec/lastenheft.md:7-15` (die Blockquote); Commit-Botschaft `9987cf7`
  (die MR-031-Abgrenzung)
- **befund:** Die Commit-Botschaft begründet den Verzicht auf einen `MR-`Eintrag damit,
  MR-031 habe „eine Abweichung von einem BENANNTEN kanonischen Schritt“ getragen, hier
  werde „kein Default verändert, sondern eine Freiheit NICHT GENUTZT“. Der Kanon nennt
  denselben Sachverhalt selbst als benannten Default („Vor `Accepted` … frei änderbar,
  ohne Change Request, ohne Historie-Zeile“), und MR-031 begründet die Eintrags-Pflicht
  ausdrücklich auch für den Fall, dass die Abweichung nur *mehr* verlangt als der Kanon —
  exakt der hier vorliegende Fall (immer Bump+Historie statt nur ab `Accepted`). Die
  behauptete Unterscheidung trägt damit nicht gegen ihre eigene Präzedenz. Zusätzlich
  weist der Kanon diese *Art* von Entscheidung (welche Versionierungs-/Historie-Disziplin
  das Repo gegenüber dem Lastenheft fährt) im unmittelbar benachbarten Absatz explizit dem
  Adaptions-Block von `harness/conventions.md` zu. Die Behauptung „ein Freshness-Audit
  hätte nichts zu vergleichen“ ist am eigenen Kanon widerlegbar: Der Freshness-Audit läuft
  laut Modul 2 gerade durch die Adaptions-**Liste** in `harness/conventions.md`, nicht
  durch die Lastenheft-Prosa — ohne Eintrag ist die Abweichung für genau diesen Mechanismus
  unsichtbar, nicht redundant.
- **verifizierbar:** teilweise — kein Gate prüft MR-Vollständigkeit (Review-Griff laut
  Kanon selbst); die Widerlegung der Selbst-Abgrenzung ist Text-gegen-Text verifizierbar
  (MR-031-Body gegen Commit-Botschaft `9987cf7`, gegen den zitierten
  Source-Precedence-Absatz).
- **klasse:** adaption-nicht-im-konventionsspeicher-registriert

### F-2 — Die zwei Tatsachenberichtigungen der Historie sind nicht „als solche ausgewiesen“, wie Kanon und Vorlage es für den CR-Verzicht verlangen

- **kategorie:** MEDIUM
- **quelle:** `grundlagen-source-precedence.md` §Spec-Stratifizierung („Die eine Ausnahme
  ist die Tatsachenberichtigung … unter zwei Bedingungen: Sie wird in der Historie **als
  solche ausgewiesen**, und sie berührt keine Aussage einer Anforderung“) ·
  `lastenheft.template.md:132-135` (identischer Wortlaut) · Slice-Plan §Bezug (nennt
  „Tatsachenberichtigung“ ausdrücklich als eines von drei Elementen aus Kurs-Welle 90)
- **pfad:** `spec/lastenheft.md:2972-2973` (Zeilen 0.65.2 und 0.65.1)
- **befund:** Beide Zeilen korrigieren nachweislich falsche Bestandsaussagen der
  `DC-FA-REF-001`-Beschreibung (welche Module den `d-check:ignore`-Marker honorieren) und
  erfüllen damit inhaltlich die Definition der Tatsachenberichtigung — sie tragen aber nur
  das Etikett „redaktionell, keine Anforderungs-Änderung“, nirgends den vom Kanon
  verlangten Begriff „Tatsachenberichtigung“ selbst oder eine gleichwertig eindeutige
  Selbstauskunft, dass hier die CR-Ausnahmeregel in Anspruch genommen wird. Eine Suche im
  gesamten Baum zeigt: Der Begriff kommt nur im Kanon, in der Vorlage und in den
  Planungs-Dateien vor, nie in einer tatsächlichen Historie-Zeile. Der Slice selbst nennt
  „Tatsachenberichtigung“ im `Bezug:`-Feld als eines der drei Kurs-Welle-90-Elemente, klärt
  in §1–§5 aber nur die vierte Spalte und die eigene Strenge — die Frage, ob der
  Bestand die Kennzeichnungspflicht der zweiten Ausnahme erfüllt, bleibt unadressiert
  (weder umgesetzt noch unter §3 als Out-of-Scope benannt).
- **verifizierbar:** ja — `grep -n "Tatsachenberichtigung" spec/lastenheft.md` liefert 0
  Treffer trotz zweier Zeilen, die die Definition erfüllen.
- **klasse:** kanon-pflichtformat-nicht-erfuellt-trotz-benennung-im-bezug

## Negativbefunde

- geprüft, ohne Befund: Header- und Separator-Zeile der Lastenheft-Historie
  (`spec/lastenheft.md:2970-2971`) sind byte-genau identisch mit
  `lastenheft.template.md:152-153` (`diff`, Exit 0)
- geprüft, ohne Befund: `spec/spezifikation.md` §7 ist von diesem Diff nicht berührt,
  führt weiterhin zwei Spalten (`| Datum | Änderung |`) und trifft
  `spezifikation.template.md:122` exakt
- geprüft, ohne Befund: keine dritte Historie-Tabelle außerhalb der beiden `## 7. Historie`
  in `spec/lastenheft.md`/`spec/spezifikation.md` — `spec/architecture.md` erklärt explizit
  „Und keine Historie“, Grep über alle `spec/*.md` findet keine weitere `## .*Historie`
- geprüft, ohne Befund: alle 95 Historie-Datenzeilen des Lastenhefts tragen die vierte
  Zelle `—`, ausnahmslos (`awk`-Zensus über den gesamten Abschnitt)
- geprüft, ohne Befund: die einzige Zeile mit trailing `d-check:ignore`-Kommentar
  (0.2.1, `spec/lastenheft.md:3064`) trägt die neue Zelle korrekt **vor** dem Kommentar;
  kein Toleranzgrenzen-Bruch
- geprüft, ohne Befund: `structure`-Regel auf `spec/lastenheft.md` §7 (`.d-check.yml`)
  nutzt die Default-Spalte 1 (Version), kein `table-column`-Schlüssel zeigt auf eine durch
  die neue Spalte verschobene Position; `make doc-check` bestätigt empirisch Exit 0,
  463 Dateien, 0 Befunde (Botschafts-Zahl exakt)
- geprüft, ohne Befund: `make completeness-check` Exit 0, „48 Anforderung(en), 0
  Waise(n)“ — Botschaft exakt
- geprüft, ohne Befund: `make verify-closure-notes` Exit 0, „427 Datei(en) geprüft, 0
  Befund(e)“ — Botschaft exakt
- geprüft, ohne Befund: §3-Scope-Treue — Lastenheft-Status bleibt `Draft` (kein
  Status-Wechsel), keine Zeile trägt den Vermerk „zurückgezogen“ (Grep, 0 Treffer), keine
  bisherige Strenge zurückgebaut (Diff ist rein additiv: neue Spalte + neue Blockquote)
- geprüft, ohne Befund: Diff-Umfang von `spec/lastenheft.md` besteht aus genau zwei
  Hunks (Kopf-Blockquote, §7-Tabelle) — keine Anforderungs-Beschreibung außerhalb der
  Historie wurde durch diesen Slice angefasst
- geprüft, ohne Befund: Lifecycle-Move-Commit `951a9e1` — reine Pfad-Retargets
  `open/` → `in-progress/` in `slice-129`, `roadmap.md`, `welle-83`, Review-Report;
  keine inhaltliche Nebenwirkung, MR-013-Bündelung korrekt (git mv + Verweise im selben
  Commit)
- geprüft, ohne Befund: die neue Blockquote in `spec/lastenheft.md:7-15` referenziert
  weder ADR noch Welle noch Slice noch Commit-Hash — kein Verstoß gegen `AGENTS.md` §3.4
  in der wörtlichen Lesart (Referenzverbot); der Fund F-1 betrifft die **Stratum-Wahl**
  für eine Adaptions-Entscheidung, nicht eine verbotene Abwärts-Referenz

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** `adaption-nicht-im-konventionsspeicher-registriert` ·
`kanon-pflichtformat-nicht-erfuellt-trotz-benennung-im-bezug`

## Verdikt

**Merge-blockierend:** ja — zwei MEDIUM offen. F-1 ist der Kernpunkt: die
Selbst-Abgrenzung von MR-031 in der Commit-Botschaft hält der eigenen zitierten
Begründung nicht stand, und der Kanon weist diese Entscheidungsklasse ausdrücklich dem
Adaptions-Block zu — die Blockquote im Vertrags-Stratum ersetzt den `MR-`Eintrag nicht,
sie verschiebt ihn nur an einen Ort, den der Freshness-Audit nicht abläuft. F-2 ist enger,
aber real: zwei Zeilen erfüllen die Tatsachenberichtigungs-Ausnahme inhaltlich, ohne sie
als solche auszuweisen — und der Slice hat diese vom Kanon selbst im `Bezug:`-Feld
genannte Frage weder beantwortet noch als Out-of-Scope deklariert.

Die übrige Umsetzung ist sauber: Spaltenzahl, Bestandszeilen-Behandlung, Kopf-Byte-Treue,
`d-check:ignore`-Sonderfall, `structure`-Regel-Unberührtheit und alle drei genannten
Gate-Botschaften (`make doc-check`, `make completeness-check`, `make verify-closure-notes`)
sind exakt wie behauptet — kein BEO-009-Richtung-(a)-Fund in diesem Lauf.

**Übergabe:** Findings gehen an den Implementer (Rückkante Review → Plan bei
Plan-Defekt); die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort
in den Zähler. Dieser Report selbst ist ein Lauf-Beleg — DoD-/Spec-Konformität prüft der
Verifier separat.
