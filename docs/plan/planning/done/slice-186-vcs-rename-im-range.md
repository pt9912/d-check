# Slice slice-186: `vcs` meldet den reinen Rename auch über `--range`

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Reaktiv: ein Befund liegt vor, und sein Closure-Grund
geht über die eigene DoD nicht hinaus (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** eingehender **Werkzeug-Befund** aus `ai-harness-course` (kein CR, keine
neue Fähigkeit beantragt) gegen
[`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in);
[ADR-0024](../../adr/0024-vcs-immutable-gate.md) (das Modul),
[ADR-0016](../../adr/0016-adr-immutable-gate.md) (das Gate, dessen CI-Hälfte
betroffen ist).

**Berührte Spec-Stellen:**
[`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
und seine `.a`-Verfeinerung in der
[Spezifikation](../../../../spec/spezifikation.md). **Kein Lastenheft-Bump** —
die Zusage steht bereits und wurde nur nicht eingelöst.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-08-31.

---

## 1. Ziel

**Ein Gate wird auf dem CI-Pfad still grün.** [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in) sagt zu:
*gelöschte oder umbenannte immutable Datei ⇒ `core-drift-vcs`*, ohne einen Modus
einzuschränken. Der `--staged`-Pfad löst das ein, der `--range`-Pfad nicht — und
`--range` ist der, den die PR-CI fährt.

**Reproduziert, alle vier Fälle:**

| Lauf | Ergebnis |
|---|---|
| `--staged`, reiner Rename | `core-drift-vcs` |
| `--range`, **derselbe** Rename committet | **0 Befunde** |
| `--range`, Rename **mit** Umformulierung | `core-drift-vcs` |
| `--range`, reines Löschen | `core-drift-vcs` |

**Die Ursache steht im Code und der Kommentar daneben behauptet das Gegenteil.**
`diffTrees` ruft `base.Diff(head)`; go-git fährt dort Rename-Erkennung, die
Änderung kommt mit `From` **und** `To` an und landet im `default:`-Zweig als
`VCSModified` auf dem **neuen** Pfad. Die Delete-Hälfte, auf die die Anforderung
zielt, entsteht nie. Fällt die Ähnlichkeit unter die Schwelle, zerlegt go-git
wieder in Delete und Add — daher meldet Fall 3.

**Wir sind selbst exponiert:** `.githooks/pre-commit` fährt `STAGED=1` (fängt),
`.github/workflows/ci.yml` fährt `RANGE=` (fängt nicht). Der laute Pfad ist der
überspringbare. [`AGENTS.md`](../../../../AGENTS.md) §3.5 nennt beide als
Durchsetzung der ADR-Immutabilität; für die CI-Hälfte stimmt das bei einem
reinen Rename nicht.

## 2. Vorgehen

1. **Den Range-Pfad ohne Rename-Erkennung diffen**, damit die Delete-Hälfte
   entsteht — der Code tut dann, was sein Kommentar seit jeher behauptet. Die
   Alternative (die Rename-Änderung in Delete + Add zerlegen) ist zu vergleichen
   und die Wahl zu begründen; der Absender beantragt ausdrücklich **keine**
   Ähnlichkeits-Erkennung.
2. **Umkehr-Probe zuerst schreiben, dann fixen:** ein Test im **Range**-Modus mit
   einem reinen Rename, der ohne den Fix rot ist. Ohne ihn wäre nicht belegt,
   dass der Fix die Lücke trifft und nicht nur etwas ändert.
3. **Die drei Kontrollfälle festhalten** (Löschen, Rename mit Umformulierung,
   `--staged`) — sie dürfen sich nicht bewegen.
4. **Die zwei falschen Zusagen korrigieren:** der Kommentar über `diffTrees` und
   die Zeile *„EINE Wahrheit: ruft dasselbe Gate wie CI"* in
   `.githooks/pre-commit`. Beide behaupten heute etwas, das nicht gilt.
5. **`AGENTS.md` §3.5 prüfen:** der Fix stellt die Aussage wieder her. Bleibt ein
   Rest, gehört er als benannte Grenze dorthin — nicht stillschweigend.
6. **Spezifikation** §[`DC-FA-VCS-001.a`](../../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs): der Mechanismus (Rename-Erkennung im
   Range-Pfad aus) gehört in den Ablauf, plus Historie-Zeile. **Kein**
   Lastenheft-Bump.
7. **Die Dokumente in Ordnung bringen:** den Befund als eingehendes Dokument
   ablegen, und die beiden CR-Dokumente vom 2026-08-31 korrigieren — sie geben
   den Abschnitt *„Bereits gelöst: Immutabilität"* wieder, den der Absender
   inzwischen zurückgezogen hat.
8. `make gates`; **Review** und **Verifikation** als getrennte Läufe; Closure.
   `CHANGELOG.md` erst in der Release-Prep.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Ähnlichkeits-Erkennung.** Sie ist Out-of-Scope der Anforderung, und der
  Absender beantragt sie ausdrücklich nicht.
- **Kein neuer Grund-Code.** `core-drift-vcs` trägt den Fall bereits; er entsteht
  nur nicht.
- **Keine Vertrags-Änderung.** Die Zusage steht; sie wird eingelöst.
- **Keine Änderung am `--staged`-Pfad.** Der ist korrekt und bleibt die
  Vergleichsgröße.
- **Kein CHANGELOG-Eintrag im Feature-Commit.**

## 4. Definition of Done

- [ ] Der reine Rename einer immutablen Datei meldet `core-drift-vcs` **auch über
      `--range`** — mit Test, und der Test ist **vor** dem Fix rot gewesen
      (Ausgabe in der Closure-Notiz).
- [ ] Die drei Kontrollfälle sind unverändert: `--staged`-Rename, `--range`-Löschen,
      `--range`-Rename-mit-Umformulierung.
- [ ] **Die zwei falschen Zusagen sind weg** (`git.go`-Kommentar,
      `.githooks/pre-commit`), und [`AGENTS.md`](../../../../AGENTS.md) §3.5 sagt
      wieder die Wahrheit — oder trägt die verbleibende Grenze benannt.
- [ ] Spezifikation §[`DC-FA-VCS-001.a`](../../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) nennt den Mechanismus; Historie-Zeile
      gesetzt; **kein** Lastenheft-Bump.
- [ ] Der Befund ist als eingehendes Dokument abgelegt, und die zwei
      CR-Dokumente vom 2026-08-31 sind um den zurückgezogenen Abschnitt
      korrigiert.
- [ ] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register
      fortgeschrieben; jedes Risiko aus §5 mit Ausgang; die drei Paarungen
      geprüft.

## 5. Abnahme-Punkte / Risiken

- **Der Fix könnte andere Befunde bewegen.** Rename-Erkennung auszuschalten
  ändert die Diff-Menge für **jede** `--range`-Prüfung des Moduls, nicht nur für
  immutable Pfade. Gegenprobe: der eigene Bestand über eine echte Range, vor und
  nach dem Fix. — **Ausgang:** *(bei Closure)*
- **Eine umbenannte Datei erzeugt künftig zwei Einträge** (Delete + Add). Ob das
  irgendwo als Doppel-Befund sichtbar wird, ist zu messen, nicht zu vermuten.
  — **Ausgang:** *(bei Closure)*
- **Die Kommentar-Korrektur ist die eigentliche Lehre und die am leichtesten
  vergessene.** Der Code war falsch, aber der Kommentar hat den Fehler *gedeckt*
  — wer ihn liest, prüft nicht nach. — **Ausgang:** *(bei Closure)*
- **`BEO-024` oder nicht:** der Absender fragt, ob dies eine zweite Instanz ist.
  Die Entscheidung fällt bei der Closure, mit dem vollen Eintrag vor Augen — eine
  vorhandene Kennung zu dehnen wäre selbst ein
  [`BEO-012`](../observations.md)-Fall. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — `in-progress/` trägt keinen
Slice.

**Rückführungen:** `in-progress` → `open`, falls die Gegenprobe aus §5 zeigt, dass
das Abschalten der Rename-Erkennung die Befund-Menge des Moduls breit verschiebt
— dann ist der Schnitt die Zerlegung in Delete + Add, und die Wahl braucht eine
eigene Begründung statt einer Zeile im Vorgehen.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `internal/adapter/driven/git` (der VCS-Adapter samt
  Tests), `spec/` (die `.a`-Verfeinerung) und `docs/plan/cr/` (die eingehende
  Korrespondenz). Alle drei fallen unter den Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration pro Sub-Area); jede trägt eigene Sensoren und einen eigenen
  Änderungs-Takt. Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-31, höchste Kennung
  `BEO-024`): [`BEO-023`](../observations.md) (Zähler 7) — *ein Wächter, der nie
  fangen konnte, liest sich wie einer, der fängt*: das **ist** dieser Slice, und
  zwar in seiner schärfsten Form, weil hier zusätzlich ein Kommentar die Deckung
  behauptet; [`BEO-024`](../observations.md) (Zähler 1) — der Absender fragt, ob
  dies eine zweite Instanz ist (Bedingung klingt nach dem Gegenstand, hängt am
  Werkzeugweg); die Entscheidung gehört an die Closure, nicht an den Plan;
  [`BEO-012`](../observations.md) (Zähler 11) — genau deshalb: eine vorhandene
  Kennung zu dehnen, weil sie ungefähr passt, ist die Klasse selbst. Die Regel,
  die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): `upstream-drift.yml`
  meldet **ROT** (Lauf 2026-08-31T06:31:40Z), `image-scan.yml` grün. Gelesen
  statt weggeklickt — drei Schritte fielen: `baseline-freshness` (bekannt,
  geschnitten als [slice-183](../open/slice-183-baseline-v5140.md)),
  `freshness-a-check` (**heute behoben**, Pin steht auf v0.19.0) und
  `go-base-digest` (**neu**: der `golang:1.27.0`-Tag ist upstream neu gebaut,
  ABWEICHEND). Keiner der drei berührt diesen Slice; der dritte ist ein eigener
  Punkt. **Dieser Block trägt bewusst keine `cite`-Direktive** — sein Ziel ist
  eine Repo-Adaption, kein Baseline-Abschnitt
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-186. Betroffene IDs:
[`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in),
[ADR-0024](../../adr/0024-vcs-immutable-gate.md),
[ADR-0016](../../adr/0016-adr-immutable-gate.md). Module: `vcs`. Gates:
`make test`, `make adr-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — alle drei berührten Sub-Areas fallen unter
den Default: Doc führt, Code folgt. Der Fix löst eine **vorhandene** Zusage ein;
es entsteht keine neue Fähigkeit, kein Fremdsystem, keine Reconciliation. Das
Evidenz-Risiko ist hier ausnahmsweise **nicht** niedrig, sondern gemessen: Code
und Kommentar standen auseinander, und der Kommentar hat gewonnen.

## 9. Closure-Notiz (nach `done/`)
