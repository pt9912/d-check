# Slice slice-184: Ein SHA trägt überall denselben Tag-Kommentar

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Sein Closure-Grund geht über die eigene DoD nicht
hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht); der Anlass ist ein **eingehender CR**, keine Welle.

**Bezug:** [der eingehende CR 5](../../cr/2026-08-30-cr-a-check-uses-tag-kohaerenz.md)
(Antrag und Beleg des Absenders);
[`DC-FA-WF-001`](../../../../spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in)
(das erweiterte Modul);
[ADR-0072](../../adr/0072-workflows-modul.md) (das Modul selbst);
[ADR-0074](../../adr/0074-offene-tasks-auf-rohen-zeilen.md) (dieselbe Bauform:
neue Bedingung, eigener Grund-Code, Lastenheft-Bump);
[ADR-0073](../../adr/0073-befund-erlaeuterung-fuer-menschen.md) (die
Erläuterung, falls die Meldung die widersprüchlichen Werte tragen soll).

**Berührte Spec-Stellen:**
[`DC-FA-WF-001`](../../../../spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in)
und seine `.a`-Verfeinerung in der Spezifikation.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-08-30.

---

## 1. Ziel

**Der Tag-Kommentar ist im Vertrag eine Zusage, und geprüft wird nur, dass er
dasteht.** `uses-pin-untagged` erzwingt seine **Existenz** — für eine Zusage
ist das die schwächste denkbare Prüfung. Zwei Zeilen, die denselben SHA
verschieden beschriften, widersprechen einander; mindestens eine ist falsch.

**Der Anlass liegt beim Absender, nicht bei uns, und das ist der Punkt.**
`a-check` pinnte denselben `docker/login-action`-Digest zweimal — `# v4.2.0`
und `# v3.6.0`, 83 Zeilen auseinander, derselbe SHA. Aufgelöst über die
GitHub-API war der zweite Kommentar falsch. Ein zweiter Fall lag zwischen zwei
Dateien.

**Bei uns fände die Regel heute nichts** — nachgemessen, nicht übernommen:
`.github/workflows/` führt **drei** distinkte SHAs mit je genau einem
Tag-Kommentar (`3d3c42e5…` fünfmal `v7.0.1`, `1b9a80c0…` `v5.0.0`,
`dbcb8138…` `v4.6.0`); **null** Konflikte. Damit ist die Regel eine
**Regressions-Bremse ohne Bestands-Ausnahme** — anders als
[`MR-049`](../../../../harness/conventions.md#mr-049) und
[`MR-056`](../../../../harness/conventions.md#mr-056), die beide eine brauchten.

**Die Abgrenzung gegen unser Out-of-Scope trägt.** Dort steht die
**Gültigkeit** eines SHA (Netz). Dieser Antrag fragt nicht, ob ein Kommentar
**wahr** ist, sondern ob zwei einander **widersprechen** — eine Aussage der
Scan-Menge gegen sich selbst, aus derselben Eingabe, die das Modul ohnehin
parst.

## 2. Vorgehen

1. **Eine dritte Bedingung der Pin-Familie:** derselbe 40-stellige SHA trägt
   innerhalb der Scan-Menge überall denselben Tag-Kommentar. Eigener Grund-Code
   `uses-pin-tag-conflict` mit eigener `SPEC-`Kennung — die Befund-Deduplikation
   läuft über (Datei, Zeile, Regel, Ziel, Grund), und die Reparatur ist eine
   andere als bei `uses-pin-untagged` (*Kommentar vereinheitlichen* statt
   *Kommentar nachtragen*).
2. **Ein Befund je beteiligter Zeile**, nicht einer je SHA. Die Regel kann
   nicht wissen, welcher Kommentar der richtige ist — das wäre die
   Gültigkeitsfrage und damit Netz. Also melden **alle** beteiligten Zeilen;
   dieselbe Wahl wie bei
   [ADR-0074](../../adr/0074-offene-tasks-auf-rohen-zeilen.md) (je Item) und
   [ADR-0069](../../adr/0069-zellenlaenge-als-strukturbedingung.md) (je Zelle).
3. **Die Meldung nennt die widersprüchlichen Werte**, nicht nur den Zustand —
   sonst sucht der Leser selbst. Ob das die modul-eigene Meldung leistet oder
   ein `hint`, ist beim Bauen zu entscheiden.
4. **Die Reichweite ist die Scan-Menge**, dateiübergreifend. Der zweite
   gemessene Fall des Absenders lag **zwischen zwei Dateien** — eine
   datei-lokale Prüfung fände ihn nicht. Das ist zugleich eine neue Eigenschaft
   des Moduls: bisher urteilt jede Bedingung je Datei.
5. **Die Falsch-Positiv-Klasse benennen, die der Antrag nicht nennt:** ein
   Commit kann **legitim zwei Tags tragen** — direkt nach einem Release zeigen
   `v4` und `v4.2.0` auf denselben SHA, und beide Kommentare wären **wahr**.
   [`AGENTS.md`](../../../../AGENTS.md) §3.9 verlangt „einen vollen Commit-SHA
   mit Tag-Kommentar" und schreibt **nicht** vor, welcher Tag. Der Entscheid
   ist zu treffen und zu begründen, nicht zu übergehen — die Zusage lautet dann
   nicht *„der Kommentar ist wahr"*, sondern *„der Kommentar ist innerhalb der
   Scan-Menge eindeutig"*.
6. **Vor dem Scharfschalten messen**, nicht behaupten: gegen den eigenen
   Bestand (erwartet null) **und** gegen einen konstruierten Konflikt, der
   meldet — beides mit Ausgabe.
7. **ADR** für die Entscheide aus 1–5, Lastenheft-Bump, `.a`-Verfeinerung
   samt §4-Grund-Code, `--print-config` prüfen, Handbuch.
8. **Antwort an den Absender** — angenommen, mit der Falsch-Positiv-Klasse, die
   er nicht genannt hat, und der Korrektur seiner Angabe über unseren Bestand
   (`v5.0.0` gehört bei uns zu `dockerhub-description`, nicht zu
   `actions/checkout`).
9. `make gates` und `make fullbuild`; **Review** und **Verifikation** als
   getrennte Läufe; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Aussage, welcher Kommentar der richtige ist.** Das ist die
  Gültigkeitsfrage, sie braucht Netz und gehört zur Freshness-Familie — der
  Absender beantragt sie ausdrücklich nicht.
- **Keine Ausweitung auf Referenzen ohne Tag-Kommentar.** Die deckt
  `uses-pin-untagged` bereits.
- **Keine Vereinheitlichung fremder Konventionen.** Ob ein Repo `v4` oder
  `v4.2.0` schreibt, entscheidet es selbst; die Regel verlangt nur
  Eindeutigkeit **innerhalb** der Scan-Menge.
- **Kein Bestandsräumen.** Es gibt keinen Bestand: null Konflikte, gemessen.

## 4. Definition of Done

- [ ] `uses-pin-tag-conflict` ist im Modul, im
      [Lastenheft](../../../../spec/lastenheft.md) (Bump + Historie mit
      CR-Bezug) und in der
      [Spezifikation](../../../../spec/spezifikation.md) (Ablauf, §2-Schema
      falls konfigurierbar, §4-Grund-Code) geführt.
- [ ] **Ein Befund je beteiligter Zeile**, gemessen: ein SHA mit zwei
      Kommentaren über **drei** Zeilen ⇒ drei Befunde, nicht einer.
- [ ] **Dateiübergreifend gemessen:** der Konflikt zwischen **zwei** Dateien
      wird gefunden — das ist der zweite Fall des Absenders und die neue
      Eigenschaft gegenüber allen bisherigen Bedingungen des Moduls.
- [ ] **Wiederholung ist kein Befund:** derselbe SHA mit **identischem**
      Kommentar über fünf Zeilen ⇒ null Befunde. Mit Test.
- [ ] **Der eigene Bestand ist gemessen**, vor und nach dem Scharfschalten:
      erwartet null, und die Zahl steht in der Commit-Botschaft.
- [ ] **Die Falsch-Positiv-Klasse ist entschieden und benannt** (zwei legitime
      Tags am selben Commit) — in ADR, Lastenheft und Handbuch, nicht nur im
      Kopf.
- [ ] **Umkehr-Proben** ([`BEO-023`](../observations.md)): je Zusage eine
      Mutation, die genau die Tests rot macht, die dagegen stehen — inklusive
      der **Verdrahtung** in den Modul-Einstieg, die bei slice-182 gefehlt hat.
- [ ] Eine ADR begründet die Entscheide; im [ADR-Index](../../adr/README.md)
      eingetragen.
- [ ] Der Absender bekommt eine **Antwort**.
- [ ] `make gates` und `make fullbuild` grün (Exit explizit); **unabhängiger
      Review**; **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.

## 5. Abnahme-Punkte / Risiken

- **Zwei legitime Tags am selben Commit.** `v4` und `v4.2.0` zeigen direkt nach
  einem Release auf denselben SHA; beide Kommentare wären wahr, die Regel
  meldete trotzdem. Die Reparatur ist billig (Kommentar vereinheitlichen), aber
  es ist eine Meldung ohne Defekt. — **Ausgang:** *(bei Closure)*
- **Die erste dateiübergreifende Bedingung des Moduls.** Alle bisherigen
  urteilen je Datei. Das ändert die Form des Befund-Ziels und die Frage, was
  bei einer **teilweise** gescannten Menge gilt. — **Ausgang:** *(bei Closure)*
- **Die Regel entsteht für einen fremden Bestand**
  ([`BEO-011`](../observations.md)): bei uns null Treffer, beim Absender zwei
  gemessene Fälle. Das ist der Normalfall für einen CR und trotzdem eine
  Aussage über einen Bestand, den wir nicht sehen. — **Ausgang:** *(bei
  Closure)*
- **Der Grund-Code-Raum des Moduls wächst auf sechs.** Ein Leser muss
  `uses-pin-missing`, `uses-pin-untagged` und `uses-pin-tag-conflict`
  auseinanderhalten — drei Codes für eine Familie. — **Ausgang:** *(bei
  Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — `in-progress/` trägt
keinen Slice.

**Rückführungen:** `in-progress` → `open`, falls sich beim Bauen zeigt, dass
die dateiübergreifende Aussage die Modul-Architektur bricht (bisher urteilt
jede Bedingung je Datei) — dann ist der Schnitt ein anderer, und die Aussage
gehört an eine Stelle, die die ganze Scan-Menge ohnehin hält.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `internal/hexagon/core/rules` (Regelmodul) und
  `spec/` (die Anforderung). Beide fallen unter den Default `*` =
  **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration). Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-30, höchste Kennung
  `BEO-024`): [`BEO-011`](../observations.md) — eine Regel aus dem Anlass statt
  aus dem Bestand: hier ist der Anlass **fremd** und unser Bestand leer, also
  ist die Frage schärfer als sonst; [`BEO-023`](../observations.md) — ein
  Wächter, der nie fangen konnte: dieser Slice **ist** ein Wächter, und seine
  DoD verlangt die Umkehr-Probe samt der Verdrahtungs-Mutation, die bei
  slice-182 gefehlt hat; [`BEO-020`](../observations.md) — die eigene Menge
  gemessen, über die fremde ausgesagt: der Antrag stützt sich auf einen
  Bestand, den wir **nicht sehen**, und unsere Antwort darf das nicht
  überspielen. Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): beide Achsen melden
  **gruen** — `upstream-drift.yml` zuletzt 2026-08-30T06:08:17Z,
  `image-scan.yml` 2026-08-30T09:16:25Z. **Der Baseline-Rückstand
  (`v5.12.0` gegen `v5.14.0`) ist bekannt und aufgeschoben** — er hängt an
  [slice-183](../in-progress/slice-183-baseline-v5150.md), nicht an diesem Slice.
  **Dieser Block trägt bewusst keine `cite`-Direktive** — sein Ziel ist eine
  Repo-Adaption, kein Baseline-Abschnitt
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-184. Betroffene IDs:
[`DC-FA-WF-001`](../../../../spec/lastenheft.md#dc-fa-wf-001--deklarations-konsistenz-von-workflow-referenzen-modul-workflows-opt-in).
Module: `workflows`. Gates: `make workflow-pins`, `make gates`,
`make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter den
Default: Doc führt, Code folgt. Eine neue Bedingung an einer vorhandenen
Anforderung; kein Fremdsystem, keine Reconciliation, kein Bestand, der
umgestellt werden müsste — gemessen null Konflikte im eigenen
`.github/workflows/`.

## 9. Closure-Notiz (nach `done/`)
