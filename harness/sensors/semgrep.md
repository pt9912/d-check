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

**Wie groß der Ausschnitt ist, sagt das Kommando:** Der Lauf nennt die Zahl
der gescannten Dateien und Regeln in seiner Zusammenfassung.

## Bindung

Bestandteil von `make gates`.
[ADR-0010](../../docs/plan/adr/0010-semgrep-hermetisches-gate.md) ·
[`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) ·
[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
