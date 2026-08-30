# Slice slice-181: Das Handbuch ist aus dem Werkzeug heraus auffindbar

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Sein Closure-Grund geht über die eigene DoD nicht
hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht); der Anlass ist ein Auftraggeber-Wunsch, keine Welle.

**Bezug:**
[`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel)
(die Hilfe-Ausgabe),
[`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
(das Makefile-Fragment),
[`DC-FA-CLI-005`](../../../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben)
(die Präzedenz: eine Ausgabe verweist auf einen anderen Einstieg).

**Berührte Spec-Stellen:**
[`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel)
(Akzeptanzkriterium *Hilfe*) und
[`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben).

**Verantwortlich:** — (bis zur Priorisierung).

**Autor:** pt9912. **Datum:** 2026-08-30.

---

## 1. Ziel

**Wer d-check aus dem Werkzeug heraus kennenlernt, findet die Dokumentation
nicht.** Gemessen: `--help` enthält **null** URLs, der Kopf des erzeugten
`d-check.mk` ebenfalls. Beide Ausgaben verweisen heute nur auf **andere
Ausgaben** — `--help` auf `--print-config`/`--suggest-config`, `d-check.mk` auf
die Release-Notes. Der Weg zum aufgabenorientierten Einstieg existiert im
Werkzeug nicht.

**Das trifft genau die Lage, für die das Fragment gedacht ist.** Ein Adopter
schreibt `d-check --print-mk > d-check.mk` und hat danach eine Datei im
fremden Repo, die von d-check erzählt, ohne zu sagen, wo man nachliest. Der
Kopf ist der einzige Ort, an dem so ein Zeiger dauerhaft mitreist.

## 2. Vorgehen

1. **Beide Ausgaben bekommen den Handbuch-Zeiger** — die Hilfe als eigene
   Trailer-Zeile neben den beiden Config-Hinweisen, der `d-check.mk`-Kopf als
   Kommentarzeile neben dem Digest-Hinweis.
2. **Die URL zeigt auf `main`, nicht auf eine Version.** Eine versionierte URL
   wäre eine **neue Release-Prep-Fläche**, die kein Gate deckt: der
   `versions`-Gate hält ausschließlich `ghcr`-präfixierte Pins gegen
   `version.md`, und nackte Tags in Prosa driften still — dieses Repo führt
   dafür bereits zwei benannte Stellen im Handbuch. Ein `blob/main`-Ziel kann
   nicht veralten; es ist dieselbe Form, die die Docker-Hub-Overview schon
   nutzt.
3. **Die Sprache braucht hier keinen Marker** — anders als in
   [`README.md`](../../../../README.md), das englisch ist und das Handbuch
   ausdrücklich als *(German)* führt. `--help` und der `d-check.mk`-Kopf sind
   **selbst deutsch**; ein Marker beantwortete eine Frage, die sich nicht
   stellt.
4. **Beide Akzeptanzkriterien mitziehen** — [`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel) nennt heute den
   Verweis auf `--print-config` als Teil der Hilfe, [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben) die
   Bestandteile des Fragments. Ein Zeiger, den kein Kriterium fordert, ist eine
   Zeile, die beim nächsten Umbau still verschwindet.
5. **Die Handbuch-Illustration von `--print-mk` zieht mit** (§4.16). Sie trägt
   `not-replayable` (abgekürzte Illustration), die Replay-Harness bricht also
   nicht — aber eine Illustration, die den Kopf zeigt und die neue Zeile nicht,
   ist ab dem Commit falsch.
6. `make gates`; **Review** und **Verifikation** als getrennte Läufe; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine ADR.** Es wird keine Schwelle gesenkt (§3.6) und keine Entscheidung
  getroffen, die über den Text hinausreicht; die eine nicht-offensichtliche
  Wahl — `main` statt Version — steht mit ihrer Begründung im Lastenheft, wo
  ein Leser sie sucht. **Wird daraus mehr** (etwa ein zweiter Zeiger, eine
  Sprach-Auswahl oder eine versionierte Form), ist das der Anlass für eine.
- **Keine zweite URL.** Nicht README, nicht die Release-Notes, nicht das
  Repository — der Auftrag ist der Handbuch-Zeiger, und drei Links in einem
  Hilfe-Trailer sind keiner.
- **Keine englische Fassung der Ausgaben.** Dass CLI und Handbuch deutsch sind,
  während `README.md` und die Docker-Hub-Seite englisch sind, ist eine
  bestehende Spannung. Sie ist größer als dieser Slice und wird hier weder
  gelöst noch verschärft.
- **Kein Zeiger in den erzeugten Targets** selbst — nur im Kopf. Ein Kommentar
  je Target wäre Rauschen in einer Datei, die eingebunden wird.

## 4. Definition of Done

- [ ] `d-check --help` nennt das Handbuch mit voller URL; `d-check --print-mk`
      trägt sie im Kopfkommentar. **Beides gegen das echte Binary gemessen**,
      nicht gegen den Quelltext.
- [ ] Die URL zeigt auf `blob/main` und enthält **keine** Versionsangabe —
      geprüft, dass ein Release sie nicht anfassen muss.
- [ ] [`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel)
      und [`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben)
      fordern den Zeiger im Akzeptanzkriterium; Lastenheft-Bump samt
      Historie-Zeile.
- [ ] Je ein Test hält die Zusage, und **jeder wird von genau der Mutation rot**,
      gegen die er steht (Zeile entfernt ⇒ rot); der Vorzustand ist mitgeprüft
      ([`BEO-023`](../observations.md)).
- [ ] Das [Benutzerhandbuch](../../../user/benutzerhandbuch.md) §4.16 zeigt den
      Kopf mit der neuen Zeile; `--print-mk` bleibt `not-replayable`.
- [ ] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.

## 5. Abnahme-Punkte / Risiken

- **Eine URL im Werkzeug ist ein Versprechen über ein fremdes System.** Wird das
  Repository umbenannt oder das Handbuch verschoben, zeigt jedes ausgelieferte
  Binary ins Leere — und kein Gate dieses Repos merkt es, weil `external`
  strikt opt-in und nie im Default aktiv ist. — **Ausgang:** *(bei Closure)*
- **`blob/main` zeigt auf den Kopf, nicht auf die ausgelieferte Version.** Wer
  ein altes Image fährt, liest ein neueres Handbuch. Das ist der bewusste Preis
  gegen die Drift-Fläche einer versionierten URL. — **Ausgang:** *(bei Closure)*
- **Der `d-check.mk`-Kopf wächst, und er reist in fremde Repos.** Jede Zeile
  dort ist eine Zeile, die ein Adopter nicht geschrieben hat und pflegen muss,
  wenn er das Fragment je von Hand anpasst. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — `in-progress/` trägt
keinen Slice.

**Rückführungen:** `in-progress` → `open`, falls sich beim Bauen zeigt, dass
der Zeiger eine Sprach- oder Versions-Entscheidung erzwingt, die über den
Auftrag hinausgeht — dann ist zuerst diese Frage zu beantworten.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `internal/adapter/driving/cli/` (die beiden
  Ausgabe-Oberflächen) und `spec/` (Anforderung und Verfeinerung). Beide fallen
  unter den Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration).
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-30, höchste Kennung
  `BEO-024`): [`BEO-022`](../observations.md) — eine Regel tritt in Kraft, bevor
  ihre Zustellung existiert: dieser Slice ist die **Umkehrung** davon, ein
  Lesepfad ohne neue Regel, und seine Prozedur (den Adressaten und dessen
  Lesepfad benennen) ist hier der Gegenstand selbst.
  [`BEO-020`](../observations.md) — **Zähler 4**: die Aussage „null URLs" ist
  über `--help` und den `--print-mk`-Kopf gemessen, nicht über die
  CLI-Ausgaben insgesamt; die Menge steht in §1.
  [`BEO-023`](../observations.md) — **Zähler 4**: die DoD verlangt deshalb je
  Test die Mutation, gegen die er steht, nicht nur einen grünen Lauf.
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen:** entfällt in `open/` — der Block entsteht
  **spätestens bei der Beanspruchung** (`open→in-progress`)
  ([`MR-053`](../../../../harness/conventions.md#mr-053)).

Slice-ID: slice-181. Betroffene IDs:
[`DC-FA-CLI-001`](../../../../spec/lastenheft.md#dc-fa-cli-001--aufruf-und-scan-wurzel),
[`DC-FA-CLI-010`](../../../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben).
Module: — (CLI-Oberfläche, kein Regelmodul). Gates: `make gates`, `make test`,
`make doc-check`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter den
Default: Doc führt, Code folgt. Zwei Textzeilen in zwei Ausgaben, zwei
Akzeptanzkriterien, zwei Tests; kein Fremdsystem, keine Reconciliation, kein
Bestand, der umgestellt werden müsste.

## 9. Closure-Notiz (nach `done/`)
