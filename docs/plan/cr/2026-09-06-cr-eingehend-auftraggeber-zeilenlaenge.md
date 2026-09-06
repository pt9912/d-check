# Eingehender Change Request — Zeilenlänge als prüfbare Größe

**Absender:** Auftraggeber · **Eingegangen:** 2026-09-06
**Richtung:** eingehend — dieses Repo ist der **Empfänger**.
**Ziel-Dokument:** [`spec/lastenheft.md`](../../../spec/lastenheft.md)
**Berührt:** [`DC-FA-STRUCT-001`](../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in) (Modul `structure`)
**Stand:** eingegangen, **noch nicht entschieden**.

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
