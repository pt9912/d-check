## Modul 0 — Einführung

*Quelle: [00-einfuehrung/modul-00-einfuehrung.md](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/00-einfuehrung/modul-00-einfuehrung.md)*

### Kernidee (Modul 0)

Ein Chatbot antwortet. Ein Agent handelt. Engineering-Systeme handeln
**reproduzierbar** und **auditierbar** — das ist nicht dasselbe wie
"antwortet besser". Der Harness ist genau das System, das aus einem
handelnden Agenten einen reproduzierbar handelnden Agenten macht.

### Regeln gegen typische Fehlannahmen (Modul 0)

- In den dokumentierten Scheiterfällen war meist nicht das Modell die Ursache, sondern eine Spec-Lücke oder ein fehlender Sensor. Das Modell rät, *weil nichts in der Eingabe widerspricht*. Lopopolo (OpenAI 2026) und die Fallstudien in [`grundlagen/fallstudien.md`](https://github.com/pt9912/ai-harness-course/blob/v1.4.0/kurs/de/grundlagen/fallstudien.md) belegen das.
- Der Prompt wird zur Anti-Spec, die niemand pflegt. Was in *jedem* Lauf relevant ist, gehört in AGENTS.md oder eine Fitness Function, nicht in den Prompt.
- Genau das macht Auditierbarkeit unmöglich. Engineering-Systeme sind *reproduzierbar*, nicht kreativ.
- Falsche Attribution. Eine Halluzination ist ein *Output-Symptom*, dessen *Ursache fast immer im Kontext liegt*: fehlende Spec-Aussage, fehlende ADR, fehlende AGENTS.md-Regel, fehlende Tool-Allowlist. Die richtige Frage ist nicht "warum hat das Modell das erfunden", sondern "was *im Kontext* hätte das Erfinden verhindert" — und genau das ist eine Harness-Frage. Wer Halluzinationen als Modell-Bug klassifiziert, kann sie nur durch Modellwechsel adressieren; wer sie als Kontext-Bug klassifiziert, kann sie durch Spec/ADR/Sensor reduzieren. Empirie: dieselbe Klasse von Halluzinationen kommt nach Modellwechsel oft *wieder* — weil das Kontext-Loch nicht zugefüllt wurde.

