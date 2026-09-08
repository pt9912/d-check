# `make semgrep` — hermetisches Security-/Static-Analysis-Gate über den Go-Code

## Vertrag

`docker run --network none` mit **gepinntem** Image und **gepinntem**,
außerhalb des Repos gecachtem Regelset (`semgrep/semgrep-rules` auf
Commit-Pin, Umfang `go/lang/security`). Ein Befund bricht das Gate
(`--error`).

**Reproduzierbar und netzlos zugleich:** Das Cache-Holen am Pin ist **Setup**
— Netz, wie ein Image-Pull —, nicht Teil der Analyse. Der Scan selbst läuft
ohne Netz.

## Grenze — was das Grün nicht abdeckt

1. **Der Umfang ist `go/lang/security`**, nicht das ganze Regelset. Ein
   grüner Lauf sagt etwas über diesen Ausschnitt.

2. **Das Regelset ist gepinnt und altert** — der Lauf misst gegen den Stand des
   Pins, nicht gegen den heutigen. Eine Regel, die upstream nach der Hebung
   entstand, existiert für diesen Gate nicht. Das ist der Preis der
   Netzlosigkeit und der Grund, warum es
   [`make freshness-semgrep`](freshness-go.md) und `make semgrep-digest`
   gibt — **beide fail-open und außerhalb von `gates`**. Ein grüner Lauf sagt
   also „nichts nach dem gepinnten Regelstand", nicht „nichts Bekanntes".
3. **Der Regel-Cache wird genau einmal geprüft — beim Holen.** Der Bezug läuft
   über einen **git-Commit-Pin**, und das ist die **stärkste** Bindung im
   ganzen Pin-Bestand dieses Repos: Ein Commit-SHA ist ein Hash über den
   **Baum**, nicht über eine Liste, die mit ihm geliefert wird. **Danach prüft
   ihn nichts mehr.** Die Bedingung für einen erneuten Bezug ist die
   **Existenz** des Cache-Verzeichnisses; wer seinen Inhalt anschließend
   ändert, wird von nichts bemerkt. Der Cache liegt **außerhalb des Repos**,
   also sieht ihn auch kein Gate, das den Arbeitsbaum prüft — er ist das
   einzige gepinnte Fremd-Artefakt dieses Repos, auf das das zutrifft.
   **Umgekehrt zu [`baseline-verify`](baseline-verify.md):** dort jeder Lauf,
   aber die Bindung beweist die Echtheit nicht; hier beweist sie sie, und der
   Lauf prüft nicht. Permanent, solange der einmalige Bezug die gewollte
   Eigenschaft ist ([ADR-0010](../../docs/plan/adr/0010-semgrep-hermetisches-gate.md)).

**Wie groß der Ausschnitt ist, sagt das Kommando:** Der Lauf nennt die Zahl
der gescannten Dateien und Regeln in seiner Zusammenfassung.

## Bindung

Bestandteil von `make gates`.
[ADR-0010](../../docs/plan/adr/0010-semgrep-hermetisches-gate.md) ·
[`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) ·
[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
