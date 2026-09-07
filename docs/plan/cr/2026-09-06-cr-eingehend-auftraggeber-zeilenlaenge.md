# Eingehender Change Request — Zeilenlänge als prüfbare Größe

**Absender:** Auftraggeber · **Eingegangen:** 2026-09-06
**Richtung:** eingehend — dieses Repo ist der **Empfänger**.
**Ziel-Dokument:** [`spec/lastenheft.md`](../../../spec/lastenheft.md)
**Berührt:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) (Modul `structure`)
**Stand:** eingegangen, **noch nicht entschieden** — der Entscheid ist am
2026-09-07 **bis zum Obermengen-Nachweis vertagt** (Auftraggeber-Entscheid);
was bis dahin gemessen ist, steht unter §Zwischenstand.

**Ablage-Hinweis.** Dies ist der **kanonische** Change Request — der Kanon
kennt genau einen, den eingehenden, als „externen Vorgang, in dem eine
Vertragsänderung mit dem Auftraggeber vereinbart wird", ausdrücklich als
„bewusst kein Harness-Konstrukt". Er liegt hier aus demselben Grund, den
[`MR-035`](../../../harness/conventions.md#mr-035) für ausgehende CRs trägt:
Die Bitte und ihre Begründung sollen den Vorgang überleben, unabhängig davon,
wie entschieden wird.

---

## Anlass, gemessen

Beim Auslagern überladener Tabellenzellen (slice-203) kam die Frage auf, ob
dieselbe Größen-Frage auch für **Fließtext-Zeilen** prüfbar sei. Der Anlass ist
eine konkrete Stelle im Adopter-Repo `ai-harness-init`,
`harness/README.md:99`: ein **Fließtext-Absatz als eine einzige Zeile mit
7613 Zeichen** und rund 29 Satzenden.

Gemessen über alle Markdown-Zeilen beider Repos, Tabellenzeilen
ausgenommen:

| Repo | Zeilen | > 400 | > 1000 | > 2000 | Maximum |
|---|---|---|---|---|---|
| d-check | 35 906 | 99 | 7 | 1 | 2045 |
| ai-harness-init | 145 120 | 696 | 68 | 14 | 7613 |

**Was die Zahlen sagen und was nicht.** Die *Rate* ist in beiden Repos
ähnlich (0,02 % gegen 0,05 % über 1000 Zeichen) — die **Spitzen** sind es
nicht: 2045 gegen 7613, und d-checks längste Zeile ist ein vendorter
Baseline-Alias, also Fremdtext. Ein Verhältnis von 1:3,7 an der Spitze ist ein
Unterschied in der Sache, die Rate allein wäre kein Befund.

## Warum das heute niemand fängt

Das Modul `structure` misst die Größe eines Abschnitts (`min-sentences`,
`max-tasks`) und seit slice-203 die einer **Tabellenzelle**
(`cell-max-chars`/`cell-min-chars`). Eine Zeichenzahl **pro Zeile** ist in
keiner Bedingung enthalten, und keines der zweiundzwanzig Module hat sie als
Gegenstand. Die Lücke ist also strukturell, nicht konfigurativ.

## Warum es zählt — und warum nicht Ästhetik

Der Schaden ist nicht optisch, sondern **diffbar**: `git` zeigt eine geänderte
7613-Zeichen-Zeile als vollständig ausgetauscht. Ein Review kann nicht sehen,
welches Wort sich bewegt hat; ein Merge-Konflikt betrifft den ganzen Absatz.
Das ist dieselbe Klasse, die slice-203 in den Tabellenzellen aufgelöst hat:
Prosa wächst, bis niemand mehr hinsieht — nur ohne Sensor, der es meldet.

## Vorschlag

Eine Bedingung `max-line-chars` im Modul `structure`, abschnitts-skopiert wie
die übrigen, opt-in wie alle Bedingungen dieses Moduls, mit eigenem Grund-Code
(etwa `section-line-oversized`).

### Akzeptanzkriterien

- **Happy Path:** Given ein Abschnitt mit `max-line-chars: N` und einer Zeile
  über N Zeichen, when das Modul läuft, then ein Befund **auf der Zeile**, mit
  ihrer gemessenen Länge.
- **Boundary (Struktur zählt nicht):** Given dieselbe Konfiguration und eine
  überlange **Tabellenzeile** oder eine Zeile in einem **Fenced Block**, when
  das Modul läuft, then **kein** Befund — die Bedingung gilt Fließtext.
- **Boundary (unteilbares Token):** Given eine Zeile, die nur aus einer langen
  URL oder einem langen Inline-Code besteht, when das Modul läuft, then **kein**
  Befund — sie ist nicht umbrechbar, und eine Regel, die Unmögliches verlangt,
  wird abgeschaltet.
- **Negativ:** Given keine `max-line-chars`-Angabe, when das Modul läuft, then
  byte-identisches Verhalten zu heute.
- Determinismus und Seiteneffektfreiheit wie alle Module
  ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus),
  [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

## Abgrenzung — und der eine Einwand, der ernst zu nehmen ist

**Die Zeile ist ein Proxy, nicht der Gegenstand.** Gemeint ist „dieser Absatz
ist zu lang"; die Zeile misst das nur in einem Repo, das seinen Markdown-Quelltext
**umbricht**. Ein Repo, das jeden Absatz als eine Zeile schreibt, hätte
durchweg lange Zeilen, ohne dass ein einziger Absatz zu lang wäre — die Regel
würde dort Formatierung mit Substanz verwechseln. Das ist derselbe Fehler, den
slice-203 bei der Tabellen-Spaltenpolsterung gemacht hat, und er gehört vor der
Entscheidung ausgeräumt, nicht danach.

**Damit ist es eine Stil-, keine Struktur-Frage** — und d-check hält bisher
die Linie, Verweise und Struktur zu prüfen, nicht Formatierung. Der Schritt
ist klein, aber er ist ein Schritt in Formatter-Territorium.

**Es gibt ein Standard-Werkzeug.** `markdownlint` deckt genau das (Regel
MD013), samt der Ausnahmen für Tabellen, Code-Blöcke und lange URLs. Der Kanon
verlangt vor einem selbstgebauten Gate die Frage, ob ein vorhandenes Werkzeug
eine **Obermenge** ist — und diese Frage ist hier offen, nicht beantwortet.

**Nicht Gegenstand:** automatisches Umbrechen (d-check ist ein Lese-Tool,
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit));
Zeilenlänge in Code-Dateien (dafür ist der Linter zuständig).

## Zu entscheiden

1. Trägt die Zeilenlänge als **Proxy** für Absatzlänge — oder ist die
   eigentliche Anforderung eine andere (Sätze je Absatz, Zeichen je Absatz)?
2. Ist `markdownlint` MD013 eine Obermenge? Falls ja, ist der Einbau als
   fremdes Gate billiger als eine eigene Bedingung.
3. Falls eigenes Modul: rechtfertigt der Nutzen den Schritt in
   Formatter-Territorium — und wo verläuft die Linie danach?

---

## Zwischenstand (2026-09-07) — gemessen, nicht entschieden

**Der Entscheid ist vertagt**, und zwar auf genau die Frage, die der CR selbst
als offen benennt: Ist `markdownlint` MD013 eine **Obermenge**? Der Kanon
verlangt dafür nicht den Datenblatt-Vergleich, sondern **je Verstoßklasse
einen Break-Test mit beiden Sensoren nebeneinander**, plus den unveränderten
Bestand, auf dem beide schweigen müssen
(Baseline-Regelwerk `modul-11-verification.md` §Fitness Function ohne
Standard-Tool). Der Nachweis steht aus; bis dahin ist über die Fragen 1 und 3
nicht entschieden.

**Was inzwischen gemessen ist: die Bedingung existiert bereits.** Ein
`forbid-pattern: '.{N,}'` im Modul `structure` ist ein Zeilenlängen-Wächter —
und zwar ein echter, nicht bloß ein zufällig passender:

| Fall | Verhalten | vom CR gefordert |
|---|---|---|
| Fließtext-Zeile über N | meldet | meldet |
| Abschnitt aus lauter kurzen Zeilen, 758 Zeichen gesamt | schweigt | schweigt — `.` matcht keinen Zeilenumbruch, die Bedingung ist **pro Zeile** |
| Zeile in einem Fenced Block | schweigt | schweigt |
| Tabellenzeile | **meldet** | soll schweigen |
| unteilbares Token (lange URL) | **meldet** | soll schweigen |
| Inline-Code-Spanne | **meldet** | soll schweigen |

Dazu zwei Formabweichungen von den Akzeptanzkriterien: Der Befund sitzt auf
der **Überschriften**-Zeile statt auf der langen Zeile, und er nennt die
gemessene Länge nicht — der Grund-Code ist das generische
`section-forbidden`.

**Der Bestand sagt, dass die Tabellen-Blindheit kein Randfall ist.** Über
`docs/`, `spec/` und `harness/`, Archiv ausgenommen:

| Schwelle | Tabellenzeilen | Fließtext-Zeilen |
|---|---|---|
| > 400 Zeichen | 524 | 281 |
| > 1000 Zeichen | 121 | 21 |

Knapp zwei Drittel bzw. rund sechs Siebtel dessen, was der Workaround meldete,
wären Falsch-Positive. **Nicht** verwendbar als Kosten-Schätzung sind
Sweep-Läufe über das ganze Repo: Zwei verschiedene Abschnitts-Selektoren
lieferten über derselben Schwelle 342 gegen 645 Befunde — die Zahl hängt am
Abschnitts-Schnitt, und das ist selbst ein Befund über den Workaround, keine
Messung des Bestands.

**Eine Vorbedingung des Nachweises, die der CR nicht kennt.**
`markdownlint` ist Node; [`AGENTS.md`](../../../AGENTS.md) §3.1 sperrt
Host-Skript-Interpreter, und der Tool-Call-Wächter blockt sie. Der Vergleich
braucht deshalb ein **digest-gepinntes Image**, denselben Weg, den `semgrep`
und `trivy` gehen. Das ist machbar und es ist ein eigener Vorgang — kein
Nebenbei-Lauf.

---

## Obermengen-Nachweis gegen `markdownlint` MD013 (2026-09-07, slice-211)

**Die Form des Nachweises gibt der Kanon vor** und schließt den
Datenblatt-Vergleich ausdrücklich aus: je Verstoßklasse ein Break-Test mit
**beiden** Sensoren nebeneinander, plus der unveränderte Bestand
(Baseline-Regelwerk `modul-11-verification.md` §Fitness Function ohne
Standard-Tool). Geprüft wird **Obermenge in drei Teilen** — dieselbe
Kandidaten-Menge, dieselben Bedingungen, dieselbe Schwelle, **wie die
Anforderung sie setzt**, nicht wie das Werkzeug sie vorbelegt.

**Reproduzierbar:** `markdownlint-cli2` als digest-gepinntes Image, netzlos:

```
docker run --rm --network none -v "<probenverzeichnis>":/workdir:ro \
  davidanson/markdownlint-cli2@sha256:173cb697a255a8a985f2c6a83b4f7a8b3c98f4fb382c71c45f1c52e4d4fed63a \
  "**/*.md"
```

mit `.markdownlint-cli2.jsonc`:
`{ "config": { "default": false, "MD013": { "line_length": N, "code_blocks": false, "tables": false, "headings": false } } }`

Das repo-eigene Mittel ist `d-check --enable structure` mit
`forbid-pattern: '.{N+1,}'` über den Abschnitt.

### Teil 1 — Bedingungen: MD013 ist Obermenge

Sieben Proben, Schwelle 300. **Zwei Proben mussten neu gebaut werden**, weil sie
ihre Klasse nicht isolierten — die erste Prosa-Probe endete zufällig ohne
Leerzeichen jenseits der Schwelle und traf damit MD013s Ausnahme für
unumbrechbare Zeilen; die Inline-Code-Probe enthielt Leerzeichen und war damit
kein *unteilbares* Token. Beide Male ähnelte die Probe der Klasse, ohne sie zu
sein.

| Klasse | vom CR gefordert | MD013 (auf die Anforderung konfiguriert) | `forbid-pattern` |
|---|---|---|---|
| K1 Fließtext-Zeile über der Schwelle | Befund | **meldet** ✔ | **meldet** ✔ |
| K2 Abschnitt aus kurzen Zeilen, Summe darüber | Schweigen | schweigt ✔ | schweigt ✔ |
| K3 Zeile in einem Fenced Block | Schweigen | schweigt ✔ | schweigt ✔ |
| K4 Tabellenzeile | Schweigen | schweigt ✔ | **meldet** ✘ |
| K5 unteilbares Token (lange URL) | Schweigen | schweigt ✔ | **meldet** ✘ |
| K6a unteilbares Inline-Code-Token | Schweigen | schweigt ✔ | **meldet** ✘ |

**Sechs von sechs gegen drei von sechs.** Auf dieser Achse ist MD013 Obermenge,
und zwar deutlich.

**Ein Nebenfund, der eine Klasse teilt:** Eine Inline-Code-Spanne **mit**
Leerzeichen (K6b) meldet bei beiden. Sie ist umbrechbar und damit keine
Ausnahme — *„Inline-Code"* ist keine Klasse, *„unteilbares Token"* ist eine.
Der CR nennt beides in einem Atemzug.

### Teil 2 — Schwelle: das eigene Mittel kann die Anforderung nicht ausdrücken

`forbid-pattern` ist RE2, und RE2 begrenzt den Wiederholungszähler auf **1000**.
Gemessen: `.{1000,}` läuft, `.{1001,}` ist ein **Konfigurationsfehler** und
nimmt den ganzen Lauf mit (Exit 2) — es schweigt nicht, es fällt.

**Die längste Zeile des Bestands misst 2045 Zeichen.** Eine Schwelle oberhalb
von 1000 ist mit dem eigenen Mittel also **nicht formulierbar**, und der
Kanon-Test *„der unveränderte Bestand, auf dem beide schweigen müssen"* ist für
`forbid-pattern` **konstruktiv unerreichbar**. MD013 nimmt jedes `N`.

### Teil 3 — Kandidaten-Menge: hier ist MD013 **keine** Obermenge

Gemessen an Probe K7 (lange Zeile in einem **anderen** Abschnitt derselben
Datei, Regel auf `## Text` skopiert): MD013 meldet, `structure` schweigt.
MD013 kennt **keinen** Abschnitts-Begriff — es urteilt datei-weit.

**Der CR verlangt genau das Gegenteil:** *„Eine Bedingung `max-line-chars` im
Modul `structure`, **abschnitts-skopiert wie die übrigen**."* Die eigene
Vormessung hatte denselben Unterschied schon gezeigt: zwei Abschnitts-Selektoren
lieferten über derselben Schwelle 342 gegen 645 Befunde.

### Der unveränderte Bestand

Bei Schwelle 1000 über `docs/`, `spec/`, `harness/` (Archiv ausgenommen):
MD013 meldet **9** Dateien, `forbid-pattern` **22**. Die Differenz ist
einseitig — **null** Dateien meldet nur MD013, **13** nur `forbid-pattern`,
und das sind genau die drei Klassen aus Teil 1, die es zu Unrecht trifft.

### Antwort auf Frage 2

**Nein — aber knapp, und die Lücke sitzt an einer anderen Stelle als vermutet.**
MD013 ist Obermenge bei den **Bedingungen** und bei der **Schwelle**; bei der
**Kandidaten-Menge** ist es das nicht, weil es die Abschnitts-Skopierung nicht
ausdrücken kann, die der CR selbst fordert. Zwei von drei Teilen genügen dem
Kanon-Test nicht: *„Ist das Werkzeug Obermenge, wird das Skript retired — sonst
benennt man die fehlende Klasse und behält es."*
