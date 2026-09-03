# Slice slice-177: Eine `structure`-Regel sagt, was zu tun ist

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Sein Closure-Grund geht über die eigene DoD nicht
hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht); der **Anlass** liegt in
[welle-86](welle-86-closure-uebergang-durchsetzen.md), die **Fähigkeit**
gehört jeder `structure`-Regel.

**Bezug:**
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(das erweiterte Modul);
[`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)
(die Befund-Zeile ist ein Vertrag);
[`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)
(die Diagnose und ihre `fixCandidate`-Abgrenzung);
[ADR-0069](../../adr/0069-zellenlaenge-als-strukturbedingung.md) (dieselbe
Bauform: eine neue `structure`-Bedingung samt Lastenheft-Bump).

**Berührte Spec-Stellen:**
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
[`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate),
[`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)
sowie ihre `.a`-Verfeinerungen in der Spezifikation.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-08-29.

---

## 1. Ziel

**Ein Befund sagt heute, was falsch ist, aber nicht, was der Regel-Autor
wollte.**

`section-forbidden` heißt „hier steht eine verbotene Wendung". Welche Zusage
diese Regel hütet und was der Leser jetzt tun soll, steht nirgends im Lauf — es
steht bestenfalls als Kommentar in der Konfiguration, also an einem Ort, den
niemand öffnet, wenn ein Gate rot wird.

**Der Grund-Code kann das nicht tragen, und zwar aus einem Struktur-Grund.**
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
führt zu jeder Bedingung eine **Reparatur**-Spalte — für `forbid-pattern` lautet
sie *„die Wendung ersetzen"*. Das ist richtig und bleibt richtig: der Grund-Code
sagt die **Art** des Defekts. Aber **eine** Art bedient viele Regeln. Jede
`forbid`-Regel dieses Repos bekäme denselben Satz, obwohl sie verschiedene
Zusagen hüten — die eine einen erfundenen Risiko-Ausgang, die andere einen
offenen DoD-Haken.

**Die fehlende Größe ist regel-lokal, nicht code-lokal.** Sie gehört an die
Regel, weil nur die Regel weiß, welche Zusage sie hütet.

**Beim Bau gemessen, und es verschiebt den Schnitt:** Das Feld gibt es
längst. [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--befund) führt
`message` als *„menschenlesbare Erläuterung (nicht stabilitätsgarantiert)"*,
**22** der Regel-Dateien setzen es — und **gerendert wird es für Menschen
nirgends**: weder in der Befund-Zeile noch in `--doctor`. Es erreicht
ausschließlich `--json`/`--yaml`, also den Maschinen-Konsumenten. Wer
`make doc-check` rot laufen lässt, sieht die Erläuterung nicht, die das
Produkt für ihn geschrieben hat.

**Damit ist die Aufgabe eine andere und eine kleinere:** nicht ein neues Feld
neben `message`, sondern (a) `message` für Menschen sichtbar machen und (b)
einer `structure`-Regel erlauben, es aus der Konfiguration zu **verfassen**.
Ein zweites Feld daneben hätte zwei Slots für dieselbe Frage geschaffen — die
Bauform, die dieses Repo an anderer Stelle als Redundanz zurückgebaut hat
([ADR-0070](../../adr/0070-tabellen-klammer-und-spaltenliste.md)).

**Der Anlass ist gemessen und liegt in diesem Repo:**
[slice-172](../done/slice-172-closure-uebergang-waechtern.md) hat den Sensor für
den offenen DoD-Haken fertig entworfen und beidseitig gemessen (37 Befunde ohne
Bestands-Ausnahme, null mit ihr, Positiv-Probe grün) — und ist an genau diesem
Punkt zurückgeführt worden: seine eigene §2 verlangt eine Meldung, die sagt,
*was zu tun ist*, und das Schema kennt sie nicht.

## 2. Vorgehen

1. **`message` wird für Menschen sichtbar.** Vierte tab-getrennte Spalte in der
   Befund-Zeile, **nur wenn gefüllt**; eigene `Hinweis:`-Zeile in `--doctor`
   unter `Stelle:`. `--json`/`--yaml` tragen das Feld bereits — dort ändert
   sich nichts. Das ist der Teil, von dem **22** Regel-Dateien sofort
   profitieren, nicht nur die neue Zusage.
2. **Warum eine vierte Spalte und keine Fortsetzungszeile.**
   [`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)
   sagt **ein Befund pro Zeile** zu, und ein Akzeptanzkriterium zählt Zeilen
   („genau zwei Befund-Zeilen"). Eine Fortsetzungszeile bräche beides. Eine
   vierte Spalte bricht es nicht: wer auf Tab trennt und die Felder 1–3 liest,
   liest weiter dasselbe. **Nur wenn gefüllt**, damit ein Befund ohne
   Erläuterung byte-identisch bleibt.
3. **Ein neues Feld `structure[].hint`**: freier Text an **einer** Regel, der
   `message` für die Befunde dieser Regel **verfasst**. Leer gesetzt ⇒ Exit 2
   — ein leerer Hinweis sagt nichts zu, dieselbe Härte, die
   `planning.closure.boilerplate` für den leeren Eintrag führt.
4. **Vorrang, ausdrücklich entschieden:** setzt die Bedingung selbst schon ein
   `message` (die Zellen-Bedingungen tun das), **gewinnt das der Regel** —
   `hint` ist die Zusage des Konfigurations-Autors, die modul-eigene Meldung
   ist die des Werkzeugs, und die Zusage steht näher am Leser. Der Fall ist zu
   messen, nicht zu behaupten.
5. **Abgrenzung gegen `fixCandidate`, ausdrücklich.** Der Fix-Kandidat ist
   **abgeleitet** und nur dort, wo er *eindeutig ableitbar* ist
   ([`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus));
   er speist `--repair` und wird zu einem anwendbaren Patch. Ein `hint` ist
   **verfasst**. Beides zu vermischen hieße, Autoren-Prosa in die
   Patch-Pipeline zu geben.
6. **ADR** für die drei Entscheide (vierte Spalte statt Fortsetzungszeile;
   `hint` schreibt in `message` statt in ein zweites Feld; Vorrang gegenüber
   der modul-eigenen Meldung), Lastenheft-Bump, Spezifikations-Verfeinerung,
   Handbuch.
7. `make gates`; **Review** und **Verifikation** als getrennte Läufe; Closure.
## 3. Ausdrücklich NICHT in diesem Slice

- **Kein `hint` für andere Module.** Sie **sehen** ihre `message` ab jetzt —
  das ist Punkt 1 —, aber sie können sie nicht aus der Konfiguration verfassen.
  `links`, `ids`, `matrix` und die übrigen konfigurieren keine Regel-Liste mit
  Autoren-Text; ihre Befunde entstehen aus dem Dokument, nicht aus einer
  benannten Zusage. Wer das später will, misst erst den Bedarf.
- **Keine Änderung an `fixCandidate` oder `--repair`.** Siehe §2 Punkt 5.
- **Keine Prüfung des Hinweis-Textes.** Ob er stimmt, ist Urteil — dieselbe
  Klasse wie ein Kommentar (`AGENTS.md` §3.7). Der Sensor prüft, dass er nicht
  leer ist.
- **Keine Anwendung auf den DoD-Haken-Sensor.** Das ist
  [slice-172](../done/slice-172-closure-uebergang-waechtern.md); dieser Slice
  liefert die Fähigkeit, nicht ihren ersten Konsumenten.

## 4. Definition of Done

- [x] `message` erscheint als **vierte Spalte** der Befund-Zeile (nur wenn
      gefüllt) und als eigene Zeile in `--doctor`; `--json`/`--yaml` bleiben
      unverändert. Je mit Test.
- [x] `structure[].hint` ist im Schema, in
      [`spec/lastenheft.md`](../../../../spec/lastenheft.md) (Bump + Historie)
      und in [`spec/spezifikation.md`](../../../../spec/spezifikation.md)
      geführt; leerer Wert ⇒ Exit 2, mit Test.
- [x] **Der Vorrang ist gemessen:** eine `structure`-Regel mit `hint` auf einer
      Bedingung, die selbst ein `message` setzt, zeigt den `hint` — als Test,
      nicht als Behauptung.
- [x] **Byte-Identität gemessen, und ihre Grenze benannt:** ein grüner Lauf
      (null Befunde) ist unverändert; ein Befund **ohne** `message` ist
      unverändert; ein Befund **mit** `message` gewinnt die vierte Spalte —
      das ist die gewollte Änderung, und sie betrifft 22 Regel-Dateien. Beide
      Ausgaben stehen in der Commit-Botschaft.
- [x] Eine ADR begründet die drei Entscheide aus §2 Punkt 6 und ist im
      [ADR-Index](../../adr/README.md) eingetragen.
- [x] Das [Benutzerhandbuch](../../../user/benutzerhandbuch.md) führt das Feld
      dort, wo es die übrigen `structure`-Schlüssel führt, **und** die vierte
      Spalte dort, wo es das Ausgabeformat beschreibt.
- [x] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.

## 5. Abnahme-Punkte / Risiken

- **Ein Hinweis ist Autoren-Text und altert wie ein Kommentar.** Er kann
  unwahr werden, ohne dass ein Gate es merkt — die Klasse, die dieses Repo als
  [`BEO-013`](../observations.md) führt, nur ohne Wächter. —
  **Ausgang:** weiter offen. Der Text ist neu, seine Alterung liegt in der
  Zukunft; [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--befund) sagt für `message` ohnehin keine Stabilität zu, und der
  Reviewer-Skill trägt den Anker für Kommentar-Klassen.
- **Die Ausgabe-Zusage wird geweitet.** Ein Konsument, der auf genau drei
  Tab-Felder besteht, sieht ab jetzt bei manchen Regeln vier. Die
  Nicht-Änderung für Regeln ohne `hint` ist messbar, die Weitung selbst bleibt
  eine Vertrags-Änderung. — **Ausgang:** eingetreten, und schärfer als
  geplant. Der Review zeigte, dass die Weitung über **fremden Text** brechbar
  war: ein `hint` mit `\n` erzeugte zwei Zeilen, einer mit `\t` fünf Felder,
  und `commits` trägt den Commit-Betreff ohne jeden Config-Rand. **Antwort in
  zwei Hälften:** die Konfiguration weist ab, der Reporter saniert. Drei
  Mutationsproben halten es jetzt.
- **Das Feld lädt zur Ausrede ein.** Ein schlecht benannter Grund-Code lässt
  sich mit einem Hinweis zudecken, statt den Code zu schärfen. —
  **Ausgang:** weiter offen. Dagegen steht die Arbeitsteilung aus
  [ADR-0073](../../adr/0073-befund-erlaeuterung-fuer-menschen.md) — Code sagt
  die Art, Hinweis die Zusage —, nicht ein Sensor.
- **Der Hinweis ERSETZT, statt zu ergänzen — und das kostet Messwerte.** Mit
  gesetztem `hint` verlieren **alle** Befunde der Regel die quantitative
  Angabe der modul-eigenen Meldung, in Befund-Zeile, `--doctor` und `--json`.
  Gemessen: `section-cell-undersized` meldet dann nicht mehr *„Zelle der Spalte
  „Titel" hat 1 Zeichen, verlangt sind 10"*, sondern nur den Hinweis. Der
  Vorrang ist entschieden ([ADR-0073](../../adr/0073-befund-erlaeuterung-fuer-menschen.md)),
  sein Preis gehört daneben. — **Ausgang:** weiter offen, und der Preis ist
  jetzt benannt statt still: er steht in diesem §5 und im Handbuch. Ob er zu
  hoch ist, entscheidet der erste Konsument, der eine quantitative Meldung
  vermisst.
- **In `--doctor` verdoppelt die Hinweis-Zeile bei manchen Modulen den
  Grund-Klartext.** Gemessen an `ids`: „Kennung im Fließtext ohne Markdown-Link
  auf ihre Definition" (Grund-Klartext) neben „Kennung ohne Link auf ihre
  Definition" (Erläuterung). Beide Texte tragen historisch dieselbe Aussage; in
  der **knappen** Befund-Zeile, die keinen Klartext hat, trägt die Erläuterung
  trotzdem. — **Ausgang:** weiter offen. Die Doppelung ist eine Eigenschaft
  der **modul-eigenen** Texte, nicht des Hinweises; sie zu räumen hieße, 22
  Meldungen gegen ihre Grund-Klartexte zu prüfen — ein eigener Schnitt.
- **Die Fähigkeit entsteht für einen einzigen Konsumenten**
  ([`BEO-011`](../observations.md)): slice-172. Ob weitere Regeln sie nutzen,
  ist nicht gemessen. — **Ausgang:** weiter offen. Der **Nutzen** ist dagegen
  breiter als der Anlass und gemessen: **21 von 23** Regel-Einstiegs-Dateien
  setzen eine Erläuterung, die ab jetzt sichtbar ist — ungedeckt bleiben nur
  `hostpaths` und `spans`.

## 6. Trigger

**Start** (`open` → `in-progress`): unmittelbar — der Slice entsteht als
Vorbedingung von [slice-172](../done/slice-172-closure-uebergang-waechtern.md),
der dafür nach `open/` zurückgeführt wurde.

**Rückführungen:** `in-progress` → `open`, falls sich beim Bau zeigt, dass die
vierte Spalte doch Konsumenten bricht, die das Repo führt — dann ist der
Entscheid ein anderer (`--doctor`-only), und die ADR trägt ihn.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `internal/hexagon/core/` (Kern: Modell und Regel) und
  `internal/adapter/driven/` (Report-Rendering). Beide fallen unter den Default
  `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration); eine eigene Deklaration führt nur `tools/harness/`.
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-29, höchste Kennung
  `BEO-022`): [`BEO-011`](../observations.md) — Regel aus dem Anlass: die
  Fähigkeit entsteht für **einen** Konsumenten, und das steht als Risiko in §5;
  [`BEO-013`](../observations.md) — ein Wächter, der nichts mehr fängt: ein
  alternder Hinweis ist die textliche Spielart derselben Klasse;
  [`BEO-022`](../observations.md) — eine Regel tritt in Kraft, bevor ihre
  Zustellung existiert: hier gering, weil ein neues **optionales** Feld
  niemanden blockiert; die Handbuch-Zeile ist trotzdem DoD-Punkt.
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): beide Achsen melden
  **gruen** — `upstream-drift.yml` zuletzt 2026-08-29T07:34:35Z,
  `image-scan.yml` 2026-08-29T10:07:43Z. **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption, kein
  Baseline-Abschnitt ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-177. Betroffene IDs:
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
[`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate),
[`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus).
Module: `structure`. Gates: `make gates`, `make test`, `make doc-check`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter den
Default: Doc führt, Code folgt. Ein neues optionales Konfigurationsfeld plus
sein Rendering; kein Fremdsystem, keine Reconciliation, kein Bestand, der
umgestellt werden müsste.

## 9. Closure-Notiz (nach `done/`)

**Geliefert, und der Schnitt hat sich beim Bauen zweimal geändert.** Der Slice
war als „neues Feld" geplant. Beim Lesen des Codes zeigte sich, dass es das Feld
**gibt**: `message` steht seit jeher in [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--befund), **21 von 23**
Regel-Einstiegs-Dateien setzen es — und für Menschen wurde es **nirgends**
gerendert, nur in `--json`/`--yaml`. Die Aufgabe wurde damit kleiner und die
Wirkung größer: nicht ein zweiter Slot, sondern der vorhandene sichtbar gemacht,
plus `structure[].hint`, damit eine Regel ihn **verfassen** kann.

**Der Beleg ist gefahren, nicht behauptet.** `make gates` grün, Exit 0, 595
Dateien, 0 Befunde. Die vierte Spalte am selben roten Lauf gegen das alte Image
gemessen: `… target-missing` wird zu `… target-missing⇥Linkziel existiert
nicht`. Ein grüner Lauf ist unverändert, ein Befund ohne Erläuterung bleibt
dreispaltig.

**Das wichtigste Ergebnis kommt aus dem Review, nicht aus dem Bau.** Zwei
Befunde trafen die tragende Zusage selbst:

- **„Ein Befund bleibt eine Zeile" war brechbar.** Die Zusage galt nur für
  Text, den das Werkzeug schreibt. Ein `hint` mit `\n` erzeugte zwei Zeilen,
  einer mit `\t` fünf Felder — und der Fall braucht gar keinen `hint`, weil
  `commits` den Commit-Betreff einsetzt. Der Report-Adapter versprach über eine
  Eingabe, die er nicht kontrolliert (`AGENTS.md` §3.8). Antwort in zwei
  Hälften: die Konfiguration **weist ab**, der Reporter **saniert**.
- **Die benannte Grenze war untermengig.** Dokumentiert waren zwei ausgenommene
  Befunde, es sind drei: die unlesbare **Einzeldatei** bekam den Hinweis und
  verlor damit ihre fail-closed-Ursache — genau die Irreführung, gegen die die
  Ausnahme gebaut ist.

**Die Verifikation hat gemessen, was ich nur behauptet hatte.** Ihre
Mutationsprobe zeigte, dass die alte Suite ein bedingungsloses Anhängen der
vierten Spalte **überlebt** — der Test dafür prüfte mit `t.Logf` statt zu
assertieren, und seine Fixture erzeugte selbst vier Spalten. Nach der Nacharbeit
fängt je ein Test jede der drei Mutationen: bedingungslose Spalte, entfernte
Sanierung, Hinweis für Nicht-Mess-Befunde. Das `report`-Paket hat dafür seine
erste Testdatei bekommen.

**Zwei Befunde stehen bewusst als Risiko statt als Fix.** Der Hinweis
**ersetzt** und kostet damit die quantitative Angabe der modul-eigenen Meldung;
und in `--doctor` verdoppelt er bei manchen Modulen den Grund-Klartext fast
wörtlich. Beides ist gemessen, benannt und in §5 mit Ausgang geführt — das
Räumen der modul-eigenen Texte ist ein eigener Schnitt.

**Fortgeschrieben:** `BEO-009` (die Zahl „22 von 31" maß Dateien statt Regeln),
`BEO-012` (der Plan zitierte [ADR-0069](../../adr/0069-zellenlaenge-als-strukturbedingung.md) für einen Entscheid von [ADR-0070](../../adr/0070-tabellen-klammer-und-spaltenliste.md)),
`BEO-022` (`--print-config` führte den neuen Schlüssel nicht). Neu registriert
ist [`BEO-023`](../observations.md) — *„Ein Wächter, der nie fangen konnte,
liest sich wie einer, der fängt"*; die Mutationsprobe ist die Antwort darauf und
war der Weg, auf dem er gefunden wurde.

**Offen und ausdrücklich benannt:** `CHANGELOG.md` ist nicht gepflegt. Die
vierte Spalte ist nutzersichtbar, aber dieses Repo führt CHANGELOG-Einträge im
Release-Prep. Der nächste Konsument ist
[slice-172](../done/slice-172-closure-uebergang-waechtern.md): sein Sensor liegt beidseitig
gemessen im Plan und wartet nur auf `hint`.
