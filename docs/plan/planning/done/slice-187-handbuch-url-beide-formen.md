# Slice slice-187: Der Handbuch-Zeiger nennt beide URL-Formen

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Sein Closure-Grund geht über die eigene DoD nicht
hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:**
[`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel)
(die Hilfe-Ausgabe) und
[`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
(der Fragment-Kopf) — beide sagen heute **eine** URL zu.

**Berührte Spec-Stellen:** beide Anforderungen oben und ihre
`.a`-Verfeinerungen in der
[Spezifikation](../../../../spec/spezifikation.md).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-08-31.

---

## 1. Ziel

**Der Zeiger ist für den falschen Leser optimiert.** `--help` und der
`--print-mk`-Kopf nennen die gerenderte GitHub-Seite. Hauptleser ist aber ein
**Code-Agent**, kein Mensch im Terminal — und für den ist die gerenderte Seite
die schlechtere Quelle.

**Gemessen, beide Formen abgerufen:**

| Form | Nutzlast | was ankommt |
|---|---|---|
| `raw.githubusercontent.com/…/refs/heads/main/…` | **174,6 KB** | der Quelltext: `[v0.71.1](../../version.md#v0.71.1)` |
| `github.com/…/blob/main/…` | **1,2 MB** | gerendert: `v0.71.1` — das Link-Ziel fehlt |

Beide sind maschinell lesbar; das übliche Argument *„Agenten können `blob` nicht
lesen"* trifft **nicht** zu und ist damit erledigt. Entscheidend sind die zwei
anderen Achsen: **siebenfache Nutzlast** für identischen Inhalt, und die
**verlorenen Link-Ziele** — ein Agent, der weitergehen will, sieht bei `blob`
nicht einmal, dass es ein Ziel gab.

**Was hier NICHT umgedreht wird:** der Kommentar über der Konstante begründet
den **Hauptzweig** ausführlich und nennt zwei Grenzen. Zur Wahl `blob` gegen
`raw` sagt er **nichts** — die Form war nie begründet. Dieser Slice füllt eine
Lücke, er widerspricht keiner Entscheidung.

## 2. Vorgehen

1. **Beide URLs nennen, nicht ersetzen.** Die gerenderte Seite bleibt für den
   Menschen, der klickt; die `raw`-Form kommt daneben für maschinelles Ziehen.
   Die Reihenfolge und die Beschriftung sind zu entscheiden — sie müssen einem
   Agenten sagen, welche er nehmen soll, ohne den Menschen zu verwirren.
2. **Die Zusage in beiden Anforderungen nachziehen.** Sie sagen heute *die* URL
   zu (Singular) — beide Akzeptanzkriterien und beide `.a`-Verfeinerungen.
3. **Alle vier Spiegel, per `grep` gefunden:** die Konstante in
   `internal/adapter/driving/cli/cli.go`, die **ungekoppelte** Zweitfassung in
   `packaging/dockerhub/overview.md` (der Kommentar warnt ausdrücklich davor),
   die Nennung im Handbuch und der Test.
4. **Die Kopplung prüfen, statt sie erneut zu dulden:** die Zweitfassung ist
   heute eine stille Spiegel-Stelle. Ob sie gebunden werden kann (Generierung,
   `versions`-Muster) oder benannt bleibt, ist zu entscheiden.
5. `make gates`; **Review** und **Verifikation** getrennt; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Wechsel des Zweigs.** Beide Formen zeigen weiter auf den Hauptzweig
  ohne Versionsangabe — die bestehende, begründete Zusage.
- **Kein Sensor für die URL.** Die Grenze bleibt, wie der Kommentar sie nennt:
  kein Lauf dieses Repos löst sie auf.
- **Keine Änderung am Handbuch-Inhalt**, nur am Zeiger darauf.

## 4. Definition of Done

- [x] `--help` und der `--print-mk`-Kopf nennen **beide** Formen, beschriftet,
      und der Test hält beide. Umgesetzt: `handbuchURLRaw`-Konstante neben
      `handbuchURL` in `cli.go`, beide Ausgaben in `writeUsage`/`mkTemplate`
      erweitert, `TestHandbuchURL_TraegtKeineVersion` prüft jetzt beide
      Zeilen einzeln (blob gegen `/blob/main/`, raw gegen
      `/refs/heads/main/`, keine der beiden mit Version).
- [x] Beide Anforderungen und beide `.a`-Verfeinerungen sagen zwei URLs zu;
      Historie-Zeilen gesetzt.
      [`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel)/[`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
      (Lastenheft 0.82.0) und ihre `.a`-Abschnitte in `spec/spezifikation.md` (Historie
      2026-09-01) umgeschrieben.
- [x] Die vier Spiegel sind nachgezogen; über die ungekoppelte Zweitfassung ist
      **entschieden**, nicht hinweggegangen. Test-Literal und
      Handbuch-Illustration tragen jetzt beide Formen (sie reproduzieren
      Werkzeug-Output wörtlich); `packaging/dockerhub/overview.md` bleibt
      bewusst bei der einen (blob-)Form — sie ist ein handkuratierter
      Link für einen Menschen, keine Output-Illustration, und braucht die
      raw-Form nicht ([`BEO-002`](../observations.md), achte Instanz).
- [x] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** — beide in eigenen Kontexten. `make gates` grün: zehn
      Gates, 651 Dateien, 0 Befunde (2026-09-02). **Review**
      (`docs/reviews/2026-09-02-slice-187-handbuch-url-code-r1.md`): 1 HIGH,
      1 MEDIUM, 1 INFO — beide ersten behoben in `96504a4`. **Verifikation**
      (eigener Kontext, nach `96504a4`): sechs von sieben Punkten konform bei
      Erstlauf; ein rekursiver Fund — der Fix-Commit `96504a4` selbst
      begründete den Nicht-Rewrite mit „nichts gepusht", was zu dem
      Zeitpunkt bereits falsch war (`origin/main` stand schon auf `1ae05ca`)
      — als Ergänzung zur neunten `BEO-009`-Instanz nachgetragen, kein
      Blocker für die Sachlage.
- [x] Closure-Notiz mit Lerneintrag; Register fortgeschrieben; jedes Risiko aus
      §5 mit Ausgang; die drei Paarungen geprüft. Siehe §9.

## 5. Abnahme-Punkte / Risiken

- **Zwei URLs sind mehr Zeilen in einer Hilfe-Ausgabe, die knapp bleiben soll.**
  Wer beide nennt, muss auch sagen, wofür jede da ist — sonst kostet die
  Ergänzung Aufmerksamkeit statt sie zu sparen. — **Ausgang: entfallen.** Die
  zweite Zeile trägt eine eigene Beschriftung („roh, für Werkzeuge/Agenten"),
  die erste bleibt unverändert — kein Mensch, der nur die erste Zeile liest,
  verliert etwas.
- **Die ungekoppelte Zweitfassung.** Sie ist seit slice-181 benannt und bleibt
  eine stille Spiegel-Stelle ([`BEO-002`](../observations.md), Zähler 7 bei
  Anlage); dieser Slice verdoppelt die zu pflegende Zeichenkette. —
  **Ausgang: eingetreten, aber enger als befürchtet.** Geprüft statt
  angenommen: die Docker-Hub-Seite (`packaging/dockerhub/overview.md`)
  reproduziert keinen Werkzeug-Output und braucht die neue rohe Form nicht —
  dort bleibt es bei einer Zeichenkette, nicht zwei. Die tatsächliche
  Verdopplung betrifft nur Test und Handbuch-Illustration, die den Bau
  ohnehin wörtlich zeigen — kein Folge-Slice, [`BEO-002`](../observations.md)
  auf 8 mit dieser Unterscheidung als Lerneintrag.
- **Die Messung stammt aus zwei Abrufen, nicht aus einer Reihe.** 174,6 KB
  gegen 1,2 MB ist ein Verhältnis aus je einem Lauf; die Größenordnung trägt,
  die exakte Zahl ist kein Mittelwert. — **Ausgang: entfallen.** Die
  tragende Entscheidung (`raw` zusätzlich nennen) hängt an der
  **Größenordnung** (siebenfache Nutzlast, verlorene Linkziele), nicht an der
  Nachkommastelle — eine Wiederholung der Messung könnte die exakte Zahl
  verschieben, nicht die Schlussfolgerung. Kein Risiko für die getroffene
  Entscheidung.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — `in-progress/` trägt keinen
Slice.

**Rückführungen:** `in-progress` → `open`, falls die Beschriftungs-Frage aus §2.1
eine Entscheidung über die Hilfe-Ausgabe insgesamt verlangt (Umfang, Reihenfolge,
Adressat) — die wäre größer als dieser Zeiger.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `internal/adapter/driving/cli` (die Konstante und beide
  Ausgaben), `spec/` (beide Anforderungen) und `packaging/` (die Zweitfassung).
  Alle drei fallen unter den Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration pro Sub-Area). Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:219-220 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** — **bei der Beanspruchung aufgefrischt**
  (Register-Stand 2026-09-02, höchste Kennung `BEO-025`):
  [`BEO-002`](../observations.md) (**8**, war 7 bei der Anlage) — die
  ungekoppelte Zweitfassung in `packaging/` ist genau diese Klasse; geprüft
  statt angenommen, ob sie sich durch diesen Slice tatsächlich vergrößert
  (Ergebnis: nein — die Docker-Hub-Seite reproduziert keinen Werkzeug-Output
  und braucht die neue rohe Form nicht, nur Test und Handbuch-Illustration
  tun das); [`BEO-012`](../observations.md) (Zähler 12) — die Versuchung,
  slice-181 als „hat sich für `blob` entschieden" zu lesen, obwohl sein
  Kommentar nur den **Zweig** begründet. Die Regel, die diesen Schritt
  vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:225-225 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)) — **bei der
  Beanspruchung gelesen:** `upstream-drift.yml` meldet **ROT** (jüngster Lauf
  2026-09-02T05:19:44Z) — aber aus zwei Gründen, die diesen Slice nicht
  betreffen: die Toolchain-Pins `go` (1.27.0→1.27.1) und `semgrep`
  (1.175.0→1.176.0) sind veraltet, `baseline-freshness` selbst ist seit
  slice-189 grün. `image-scan.yml` grün. **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption
  ([`MR-054`](../../../../harness/conventions.md#mr-054)).

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — alle drei berührten Sub-Areas fallen unter
den Default. Eine Zeiger-Form ändert sich; kein Fremdsystem, keine
Reconciliation, kein Bestand, der umgestellt werden müsste.

## 9. Closure-Notiz (nach `done/`)

**Geliefert.** `--help` und der `--print-mk`-Kopf nennen jetzt beide
Handbuch-URLs — die gerenderte GitHub-Seite (unverändert zuerst) und die
rohe `raw.githubusercontent.com`-Form danach, mit eigener Beschriftung für
Werkzeuge/Agenten. Gemessen, nicht angenommen: 174,6 KB gegen 1,2 MB
Nutzlast für identischen Inhalt, und die rohe Form allein erhält
Markdown-Linkziele. Beide Anforderungen
([`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel),
[`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben),
Lastenheft 0.82.0) und ihre `.a`-Verfeinerungen sagen jetzt zwei URLs zu.

**Was funktioniert hat.** Die vier gegrepten Spiegelstellen zwangen zur
Einzelfall-Prüfung statt zur pauschalen Verdopplung: Test und
Handbuch-Illustration reproduzieren Werkzeug-Output wörtlich und bekamen
beide Formen; `packaging/dockerhub/overview.md` ist ein handkuratierter
Link für einen Menschen und blieb bewusst bei der einen. Diese
Unterscheidung — **Illustration** gegen **Verweis** — ist die eigentliche
Antwort auf die seit slice-181 offene Kopplungsfrage
([`BEO-002`](../observations.md), achte Instanz): nicht jeder
Spiegel-Kandidat ist ein Spiegel.

**Was anders lief.** Diese Closure trägt eine ungewöhnlich dichte
Fehlerkette, alle um dieselbe Klasse
([`BEO-009`](../observations.md), neunte Instanz plus rekursive Ergänzung):

1. Der **Beanspruchungs-Commit** (`e25f8b0`) hakte DoD-Punkte ab und erhöhte
   das Register mit vollständiger Vergangenheitsform, bevor die
   Implementierung (`1ae05ca`) überhaupt committet war — Ursache war die
   tatsächliche Arbeitsreihenfolge dieser Sitzung (geplant → implementiert →
   Beanspruchung nachgetragen), die der Plan-Text nicht offenlegte. Gefunden
   vom **Review**.
2. Der **Fix-Commit** dafür (`96504a4`) begründete den Verzicht auf eine
   Commit-Umschreibung mit „nichts gepusht" — unbelegt, und zu dem Zeitpunkt
   bereits falsch: `origin/main` stand da schon auf `1ae05ca`. Gefunden von
   der **Verifikation**.
3. Eine Commit-Botschaft (`1ae05ca`) zitierte
   [`MR-021`](../../../../harness/conventions.md#mr-021) statt
   [`MR-025`](../../../../harness/conventions.md#mr-025) als
   Quelle der gelösten Kopplungsfrage. Ebenfalls vom **Review** gefunden.

**Drei unabhängige Rollen — Implementer, Review, Verifikation — haben in
diesem einen Slice denselben Fehlertyp je einmal übersehen, bevor die
letzte ihn fing.** Der Code-Diff selbst war bei jeder Prüfung sauber; alle
drei Funde lagen im **Beleg**, nicht in der Sache. Keine Commit-Umschreibung
in allen drei Fällen — die Inhalte sind an ihrem jeweiligen Ort richtig
gestellt, benannt statt verschwiegen.

**Steering-Loop-Eintrag.** [`BEO-009`](../observations.md) auf 9 (neunte
Instanz: Register/DoD vor dem Code) plus eine rekursive Ergänzung im selben
Vorgang (unbelegte „nichts gepusht"-Behauptung im eigenen Fix-Commit) —
neue Prozedur-Zeile: eine Aussage über den Repo-Zustand ist so
belegpflichtig wie jede andere Tatsachenbehauptung, auch in einem Commit,
der einen Beleg-Fehler korrigiert. [`BEO-002`](../observations.md) auf 8
(achte Instanz: Illustration-vs-Verweis als Kriterium, ob ein zweiter
String mitspiegeln muss).

**Verifikation.** `make gates` Exit 0 (zehn Gates, 651 Dateien, 0 Befunde) ·
`make test` Exit 0, alle Pakete grün · unabhängiger Review
(1 HIGH/1 MEDIUM/1 INFO, behoben in `96504a4`) · unabhängige Verifikation
(fand den rekursiven Beleg-Fehler in `96504a4`, behoben in `43d9a26`).

**Drei Paarungen** (wellenlos, hier geprüft): **Anker** — jede zitierte
`BEO-<NNN>` hat eine Registerzeile (BEO-002, BEO-009, BEO-012, BEO-025 —
alle existieren). **Folge-Slice** — keiner benannt, keiner fällig. **Register** —
BEO-002 (Zähler 8) und BEO-009 (Zähler 9) tragen je genau so viele Belege wie
ihr Zähler.
