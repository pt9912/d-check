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

**Verantwortlich:** —.

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

- [ ] `--help` und der `--print-mk`-Kopf nennen **beide** Formen, beschriftet,
      und der Test hält beide.
- [ ] Beide Anforderungen und beide `.a`-Verfeinerungen sagen zwei URLs zu;
      Historie-Zeilen gesetzt.
- [ ] Die vier Spiegel sind nachgezogen; über die ungekoppelte Zweitfassung ist
      **entschieden**, nicht hinweggegangen.
- [ ] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** — beide in eigenen Kontexten.
- [ ] Closure-Notiz mit Lerneintrag; Register fortgeschrieben; jedes Risiko aus
      §5 mit Ausgang; die drei Paarungen geprüft.

## 5. Abnahme-Punkte / Risiken

- **Zwei URLs sind mehr Zeilen in einer Hilfe-Ausgabe, die knapp bleiben soll.**
  Wer beide nennt, muss auch sagen, wofür jede da ist — sonst kostet die
  Ergänzung Aufmerksamkeit statt sie zu sparen. — **Ausgang:** *(bei Closure)*
- **Die ungekoppelte Zweitfassung.** Sie ist seit slice-181 benannt und bleibt
  eine stille Spiegel-Stelle ([`BEO-002`](../observations.md), Zähler 7); dieser
  Slice verdoppelt die zu pflegende Zeichenkette. — **Ausgang:** *(bei
  Closure)*
- **Die Messung stammt aus zwei Abrufen, nicht aus einer Reihe.** 174,6 KB
  gegen 1,2 MB ist ein Verhältnis aus je einem Lauf; die Größenordnung trägt,
  die exakte Zahl ist kein Mittelwert. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — `in-progress/` trägt keinen
Slice. Heute hält [slice-174](../done/slice-174-register-deckung.md) den
Slot.

**Rückführungen:** `in-progress` → `open`, falls die Beschriftungs-Frage aus §2.1
eine Entscheidung über die Hilfe-Ausgabe insgesamt verlangt (Umfang, Reihenfolge,
Adressat) — die wäre größer als dieser Zeiger.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `internal/adapter/driving/cli` (die Konstante und beide
  Ausgaben), `spec/` (beide Anforderungen) und `packaging/` (die Zweitfassung).
  Alle drei fallen unter den Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration pro Sub-Area). Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.15.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-31, höchste Kennung
  `BEO-024`): [`BEO-002`](../observations.md) (Zähler 7) — die ungekoppelte
  Zweitfassung in `packaging/` ist genau diese Klasse, und dieser Slice
  vergrößert sie; [`BEO-012`](../observations.md) (Zähler 11) — die Versuchung,
  slice-181 als „hat sich für `blob` entschieden" zu lesen, obwohl sein
  Kommentar nur den **Zweig** begründet. Die Regel, die diesen Schritt
  vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.15.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — alle drei berührten Sub-Areas fallen unter
den Default. Eine Zeiger-Form ändert sich; kein Fremdsystem, keine
Reconciliation, kein Bestand, der umgestellt werden müsste.

## 9. Closure-Notiz (nach `done/`)
