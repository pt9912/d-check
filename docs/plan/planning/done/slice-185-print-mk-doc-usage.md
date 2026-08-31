# Slice slice-185: `--print-mk` — `doc-usage` (die Hilfe des Werkzeugs selbst)

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Sein Closure-Grund geht über die eigene DoD nicht
hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht); der Anlass ist ein **Auftraggeber-Wunsch**, keine Welle.

**Bezug:** Change Request (Auftraggeber) an
[`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
(Lastenheft **0.79.0 → 0.80.0**). Exponiert den bestehenden Hilfe-Modus aus
[`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel)
als Fragment-Target. **Kein ADR** — additiv, Fragment-Erweiterung; dieselbe
Einordnung wie beim direkten Vorgänger
[slice-047](../done/slice-047-print-mk-doctor-repair-help-digest.md), der
`doc-doctor`/`doc-repair`/`doc-help` auf demselben Weg nachzog.

**Berührte Spec-Stellen:**
[`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
und seine `.a`-Verfeinerung in der
[Spezifikation](../../../../spec/spezifikation.md) (§[`DC-FA-CLI-010.a`](../../../../spec/spezifikation.md#dc-fa-cli-010a--makefile-fragment),
Punkt 5).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-08-31.

---

## 1. Ziel

**Wer das Fragment einbindet, hat keinen Weg zu den Optionen des Werkzeugs.**
`d-check.mk` verteilt zwölf Targets, aber die Frage *„was kann dieses Image
eigentlich?"* beantwortet nur, wer die `docker run`-Beschwörung selbst
zusammensetzt. `doc-usage` macht daraus ein Target — dieselbe Klasse wie
`doc-doctor` und `doc-repair`: ein **bestehender** Modus bekommt eine
Fragment-Oberfläche, kein neues Verhalten.

**Der Name ist eine Entscheidung, kein Zufall.** `doc-help` gibt es bereits und
listet die `doc-*`-Targets — das ist die **Make**-Ebene. `doc-usage` fragt nach
dem **Werkzeug**. Ein `doc-check-help` stand zur Wahl und sortiert in der Liste
hübsch neben `doc-check`, liest sich aber wie *„Hilfe zum doc-check-Target"*
statt wie *„Hilfe des Prüfers"*; die zwei Ebenen sollen im Namen
auseinanderliegen, nicht in vier Buchstaben.

**Der eigentliche Gegenstand ist die Aufzählung, nicht das Target.** Zwölf wird
dreizehn, und diese Zahl steht samt Target-Liste an mehreren Stellen in
Lastenheft, Spezifikation, Code-Kommentar und Handbuch. Dieselbe Aufzählung ist
an dieser einen Anforderung bereits **dreimal** saniert worden — Lastenheft
0.37.1 und 0.57.1, Spezifikation §7 am 2026-08-31. Der Zusatz ist billig; das
vollständige Nachziehen ist die Arbeit.

## 2. Vorgehen

1. **Spiegel zuerst auflisten, dann editieren**
   ([`BEO-002`](../observations.md)): `grep` nach dem **alten** Wortlaut
   (`zwölf`, die Target-Aufzählungen) über
   [Lastenheft](../../../../spec/lastenheft.md),
   [Spezifikation](../../../../spec/spezifikation.md),
   `internal/adapter/driving/cli/print_mk.go` und
   [Handbuch](../../../user/benutzerhandbuch.md). Die Liste ist das
   Arbeitsblatt und geht in die Closure-Notiz.
2. **Lastenheft-Bump 0.80.0:**
   [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
   Beschreibung (zwölf → dreizehn, `doc-usage` benannt), **beide**
   Akzeptanzkriterien (Happy Path und Boundary) und die Out-of-Scope-Klausel,
   die die Menge heute abschließend führt; Historie-Zeile mit CR-Bezug.
3. **Spezifikation:** §[`DC-FA-CLI-010.a`](../../../../spec/spezifikation.md#dc-fa-cli-010a--makefile-fragment) Punkt 5 — Zahl **und**
   Aufzählungspunkt, in dieser Reihenfolge gedacht; Historie-Zeile.
4. **Generator:** Template-Block in `internal/adapter/driving/cli/print_mk.go`,
   `@`-echo-unterdrückt wie `doc-help` und `doc-repair` (die Ausgabe **ist** die
   Nutzlast, nicht der Befund), plus der Kopfkommentar der Funktion, der die
   Target-Zahl und die Target-Liste führt. **Der Block darf kein `%` tragen** —
   das Template geht durch `fmt.Sprintf`, und der Kommentar sagt die Zahl der
   Verben zu.
5. **Test:** Akzeptanz-Fall analog zu den Vorgängern, plus **Umkehr-Probe**
   ([`BEO-023`](../observations.md)): der Block wird entfernt, genau dieser Test
   wird rot, mit Ausgabe.
6. **Den Ausgabe-Strom messen, nicht behaupten**
   ([`BEO-020`](../observations.md)): `--help` schreibt nach Lage des Codes auf
   **stderr** und endet mit Exit 0. Das ist gelesen, nicht gemessen — vor der
   Handbuch-Zeile fällt der Lauf.
7. **Handbuch** §4.16: Zahl, Liste, und der Satz zum Ausgabe-Strom (ein
   `make doc-usage > datei` fängt nichts). **§11-Zeile und Handbuch-Kopf nicht** —
   die tragen eine Software-Version, die erst beim Release-Schnitt feststeht;
   gemessen an der Historie der Datei entstehen sie im Release-Prep-Commit
   ([`AGENTS.md`](../../../../AGENTS.md) §5, dieselbe Grenze wie beim CHANGELOG).
8. `make gates`; **Review** und **Verifikation** als getrennte Läufe; Closure.
   `CHANGELOG.md` und der Handbuch-Kopf bleiben unberührt — die gehören in die
   Release-Prep ([`AGENTS.md`](../../../../AGENTS.md) §5).

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Umbau von `doc-help`.** Es bleibt die Liste der `doc-*`-Targets. Zwei
  Targets, zwei Ebenen — das ist die Absicht, nicht ein Rest.
- **Keine Änderung der `--help`-Ausgabe selbst.** Die gehört
  [`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel);
  dieses Target exponiert sie, es formt sie nicht.
- **Kein Umlenken des Stroms.** Ein `2>&1` im Recipe wäre eine Aussage über den
  CLI-Vertrag, getroffen im Fragment — der falsche Ort. Die Eigenschaft wird
  dokumentiert, nicht überschrieben.
- **Kein Weglassen des `-v`-Mounts**, obwohl `--help` ihn nicht braucht
  (Kurzschluss vor jedem Repo-Zugriff). Der Vertrag sagt die Form für **alle**
  Targets zu; eine Ausnahme kostete einen Satz und brächte nichts.
- **Kein ADR** (siehe Kopf) und **kein CHANGELOG-Eintrag** im Feature-Commit.

## 4. Definition of Done

- [x] `doc-usage` steht im Fragment — `@`-echo-unterdrückt, `##`-annotiert, mit
      `--help` gegen `$(DCHECK_REF)`; `make doc-help` listet es auf.
- [x] Der **Vertrag** führt dreizehn Targets: Lastenheft (Beschreibung, beide
      Akzeptanzkriterien, Out-of-Scope, Bump + Historie) und Spezifikation
      (Punkt 5: Zahl **und** Aufzählung, Historie).
- [x] Das **Handbuch** §4.16 nennt dreizehn Targets, `doc-usage` in der Liste,
      und den **Ausgabe-Strom** — gemessen, nicht gelesen. **§11-Zeile und
      Handbuch-Kopf bleiben der Release-Prep vorbehalten** (Software-Version).
- [x] **Umkehr-Probe** ([`BEO-023`](../observations.md)): ohne den
      Template-Block wird genau der neue Test rot, mit Ausgabe in der
      Closure-Notiz.
- [x] **Die Spiegel-Liste aus §2.1 steht in der Closure-Notiz** — vollständig
      abgearbeitet, mit der Zahl der angefassten Stellen.
- [x] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register
      fortgeschrieben; jedes Risiko aus §5 mit Ausgang; die drei Paarungen
      geprüft.

## 5. Abnahme-Punkte / Risiken

- **Die Aufzählung ist der Gegenstand, nicht das Target**
  ([`BEO-002`](../observations.md), Zähler 5). Dieselbe Menge ist an dieser
  Anforderung dreimal nachgezogen worden, zuletzt heute — die vierte Sanierung
  wäre die Bestätigung, dass eine Liste, die eine Menge spiegelt, ohne Bindung
  an ihre Quelle nicht stabil bleibt. Ein Sensor dafür existiert nicht; die
  Gegenmaßnahme ist §2.1 und sonst nichts. — **Ausgang:** *weiter offen* → [`BEO-002`](../observations.md), Zähler **6**. Es ist eingetreten, wenn auch als Beinahe-Treffer: die Spiegel-Liste war handverlesen und übersah `docs/user/operations.md`, das seit Langem bei drei von zwölf Targets stand. Gefunden per Seitenblick, nicht per Methode; die Datei spiegelt jetzt nicht mehr, sondern zeigt auf das erzeugte Fragment.
- **Der Ausgabe-Strom.** `--help` schreibt auf stderr; ein Konsument, der
  `make doc-usage > datei` schreibt, bekommt eine leere Datei und einen
  Erfolgs-Exit. Das ist die Eigenschaft des CLI, nicht des Targets — aber das
  Target macht sie erstmals bequem erreichbar. — **Ausgang:** *entfallen* — nicht weil die Eigenschaft verschwände, sondern weil sie gemessen und benannt ist: 0 Byte stdout, 2488 Byte stderr, Exit 0, und der Satz steht in Lastenheft, Spezifikation und Handbuch §4.16. Die Verifikation hat den Fall am eingebundenen Fragment nachgestellt (`make doc-usage > out.txt` ⇒ 0-Byte-Datei). Was bleibt, ist eine dokumentierte Eigenschaft des CLI, kein offener Punkt dieses Slice.
- **Das `%`-Verbot im Template.** Der Kopfkommentar sagt die Zahl der
  `fmt`-Verben zu; ein `%` im neuen Block bräche `fmt.Sprintf` zur Laufzeit,
  nicht beim Übersetzen. Der geplante Block trägt keines — aber die Zusage ist
  eine Falle für den nächsten Zusatz. — **Ausgang:** *entfallen* — nicht weil die Eigenschaft verschwände, sondern weil sie gemessen und benannt ist: 0 Byte stdout, 2488 Byte stderr, Exit 0, und der Satz steht in Lastenheft, Spezifikation und Handbuch §4.16. Die Verifikation hat den Fall am eingebundenen Fragment nachgestellt (`make doc-usage > out.txt` ⇒ 0-Byte-Datei). Was bleibt, ist eine dokumentierte Eigenschaft des CLI, kein offener Punkt dieses Slice.
- **Zwei Targets, die beide „Hilfe" bedeuten.** Wer `doc-help` tippt, bekommt
  die Target-Liste und nicht die Optionen. Die `##`-Annotationen tragen die
  Unterscheidung — das ist Text, kein Wächter. — **Ausgang:** *entfallen* — nicht weil die Eigenschaft verschwände, sondern weil sie gemessen und benannt ist: 0 Byte stdout, 2488 Byte stderr, Exit 0, und der Satz steht in Lastenheft, Spezifikation und Handbuch §4.16. Die Verifikation hat den Fall am eingebundenen Fragment nachgestellt (`make doc-usage > out.txt` ⇒ 0-Byte-Datei). Was bleibt, ist eine dokumentierte Eigenschaft des CLI, kein offener Punkt dieses Slice.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — `in-progress/` trägt
keinen Slice.

**Rückführungen:** `in-progress` → `open`, falls der Review die Einordnung
*kein ADR* nicht trägt — dann ist die Entscheidung über die abschließende
Out-of-Scope-Klausel zuerst zu treffen, und dieser Slice wartet auf sie.
`in-progress` → `next`, falls die Spiegel-Liste aus §2.1 mehr als die drei
Liefer-Punkte aufmacht (Fragment · Vertrag · Handbuch).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `internal/adapter/driving/cli` (der Generator samt
  Akzeptanz-Tests), `spec/` (die Anforderung und ihre Verfeinerung) und
  `docs/user/` (das Handbuch). Alle drei fallen unter den Default `*` =
  **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration pro Sub-Area); keine ist zu grob geschnitten — jede trägt
  eigene Konventionen, eigene Sensoren und einen eigenen Änderungs-Takt. Die
  Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-31, höchste Kennung
  `BEO-024`): [`BEO-002`](../observations.md) (Zähler 5) — eine Semantik-Änderung
  wird nur im Körper nachgezogen, ihre Ränder bleiben stehen: das **ist** dieser
  Slice, und die Zahl 12 → 13 ist genau so eine Menge, die sich in Prosa
  spiegelt; [`BEO-023`](../observations.md) (Zähler 7) — ein Wächter, der nie
  fangen konnte: der neue Akzeptanz-Test muss ohne den Template-Block rot
  werden, sonst prüft er nur, dass etwas dasteht;
  [`BEO-020`](../observations.md) (Zähler 6) — die eigene Menge gemessen, über
  eine andere ausgesagt: die stderr-Eigenschaft ist **gelesen**, und gelesen ist
  nicht gemessen. Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): beide Achsen melden
  **gruen** — `upstream-drift.yml` zuletzt 2026-08-30T06:08:17Z,
  `image-scan.yml` 2026-08-30T09:16:25Z. **Beide Zeitstempel sind vom Vortag**,
  und der Lese-Schritt sieht den jüngsten Lauf, nicht sein Alter — die
  benannte Grenze des Targets, hier sichtbar statt bloß dokumentiert. Der
  Baseline-Rückstand (`v5.12.0` gegen `v5.14.0`) hängt an
  [slice-183](../open/slice-183-baseline-v5150.md), nicht an diesem Slice. **Dieser
  Block trägt bewusst keine `cite`-Direktive** — sein Ziel ist eine
  Repo-Adaption, kein Baseline-Abschnitt
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-185. Betroffene IDs:
[`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
(erweitert),
[`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel)
(exponiert, unverändert). Module: keines — der Generator ist CLI-Oberfläche, kein
Regelmodul. Gates: `make test`, `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — alle drei berührten Sub-Areas fallen unter
den Default: Doc führt, Code folgt. Ein zusätzliches Target an einer
vorhandenen Anforderung; kein Fremdsystem, keine Reconciliation, kein Bestand,
der umgestellt werden müsste. Der Aufwand liegt vollständig im Nachziehen der
Aufzählung, und das ist Doku-Arbeit im führenden Stratum.

## 9. Closure-Notiz (nach `done/`)

- **Was hat funktioniert:** Die **Reihenfolge** aus §2 — erst die Spiegel
  auflisten, dann editieren — hat den Slice getragen: die Zahl 12 → 13 stand am
  Ende an neun Stellen konsistent, und keine blieb zurück. Die **Umkehr-Probe**
  war ihr Geld wert und zweimal unabhängig gefahren (von mir und vom Review):
  ohne den Template-Block kippt genau ein Test. Und die Entscheidung, die
  `stderr`-Eigenschaft zu **messen** statt sie aus dem Code zu lesen, hat den
  Satz im Handbuch belastbar gemacht — die Verifikation konnte ihn am
  eingebundenen Fragment nachstellen, statt ihn zu glauben.
- **Was ging anders als geplant:** Drei Dinge. (1) Die Spiegel-Liste war
  **handverlesen**, nicht gegrept — `docs/user/operations.md` stand nicht darin
  und war seit Langem bei drei von zwölf Targets stehengeblieben; gefunden hat
  es ein Seitenblick. Sie spiegelt jetzt nicht mehr, sondern zeigt auf das
  erzeugte Fragment. (2) Der Plan sagte eine **§11-Zeile im Handbuch** zu; die
  Historie der Datei zeigt, dass §11 und der Handbuch-Kopf im
  Release-Prep-Commit entstehen, weil sie eine Software-Version tragen, die
  erst beim Schnitt feststeht — der Plan ist an dieser Grenze nachgezogen
  worden, nicht der Bestand. (3) Der Review fand eine **Exklusivitäts-Aussage**,
  die ich im selben Zug geschrieben hatte, in dem ich eine Aufzählung sanierte.
- **Steering-Loop-Eintrag:** Regel **geschärft** —
  [`MR-013`](../../../../harness/conventions/MR-013-lifecycle-move-buendelung.md)
  gilt jetzt für **beide** Slice-Lifecycle-Richtungen statt nur für die
  Closure; die Bündelungs-Pflicht bei `open/` → `in-progress/` ist damit
  belegt statt geduldet — liegt in [`MR-013`](../../../../harness/conventions/MR-013-lifecycle-move-buendelung.md) `§Geltungsbereich`.
  Auslöser: [`BEO-012`](../observations.md) — nicht der Schwellen-Übertritt
  (der liegt lange zurück und ist in [`AGENTS.md`](../../../../AGENTS.md) §5
  verkörpert), sondern eine **neue Ausprägung**: die Quelle widersprach sich
  selbst, Übersichtstabelle breit gegen Feld eng, und 71 von 83
  Beanspruchungs-Commits folgten der Tabelle.
- **Beobachtungs-Register (`../observations.md`):**
  [`BEO-002`](../observations.md) auf **6** erhöht (handverlesene statt
  gegrepte Spiegel-Liste), [`BEO-011`](../observations.md) auf **5**
  (Exklusivitäts-Aussage im verbindlichen Stratum),
  [`BEO-012`](../observations.md) auf **11** (Kennung für einen Übergang, den
  das Feld nicht nannte) — Beleg jeweils slice-185. **Alle drei fand der
  zweite Leser, keinen der Schreibende.**
- **Folge-Slices:** keine. Zwei Entscheidungen liegen daneben und gehören
  **nicht** diesem Slice: die Hebung des a-check-Pins (0.17.0 gegen upstream
  0.19.0, gemessen) und die Neubewertung von
  [slice-184](../open/slice-184-uses-tag-kohaerenz.md), seit der CR-Absender
  die beantragte Regel selbst gebaut hat.
- **Risiken aus §5:** vier, jedes mit genau einem Ausgang — einmal *weiter
  offen* (→ [`BEO-002`](../observations.md)), dreimal *entfallen* mit
  Begründung. Siehe §5.
- **Drei Paarungen:** (a) **Anker** — das Feld `liegt in` zeigt auf
  [`MR-013`](../../../../harness/conventions/MR-013-lifecycle-move-buendelung.md) `§Geltungsbereich`; der Zielort existiert und trägt den Anker
  `seit slice-185`. (b) **Folge-Slice** — keiner genannt, nichts zu paaren.
  (c) **Register** — die vier zitierten Kennungen (`BEO-002`, `BEO-011`,
  `BEO-012`, `BEO-023`) haben je genau eine Zeile. Zur zweiten Hälfte, gemessen
  statt pauschal behauptet: in der Haupttabelle trägt jede Zeile einen Beleg
  **außer** [`BEO-021`](../observations.md), deren Beleg-Zelle `—` bei Zähler
  **0** steht — das ist konsistent, keine Lücke. [`BEO-001`](../observations.md)
  und [`BEO-005`](../observations.md) stehen in der zweiten Tabelle
  *Gestrichene Einträge*, die keine Beleg-Spalte führt; auf sie zielt die
  Paarung nicht.
