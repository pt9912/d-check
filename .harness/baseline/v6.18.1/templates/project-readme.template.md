# <Projektname>

> **Template-Hinweis.** Vorlage für das Projekt-Root-`README.md`. Kopiere
> nach `README.md` deines Repos, ersetze `<Platzhalter>` und lösche diesen
> Block. Das README ist **Rang 6** der Source Precedence (Projekt-Überblick)
> — es *verweist* auf die kanonischen Quellen, es *dupliziert* sie nicht.
> Tipp: oft zuletzt in Phase 1 füllen, wenn die verlinkten Artefakte stehen.
> Hintergrund: [Baseline-Regelwerk Modul 2 — Harness-Bootstrap](../regelwerk/modul-02-harness-bootstrap.md).

**Rolle:** Rang 6 der Source Precedence — verweist auf die kanonischen
Quellen, dupliziert sie nicht. Regeln: Baseline-Regelwerk
`modul-02-harness-bootstrap.md` §Ziel-Form: Projekt-README.

## Was ist <Projektname>?

<!-- 2–3 Sätze: was leistet es, für wen, gegen welche Annahme. Überblick,
nicht Implementierung (die lebt in spec/). -->

<…>

## Was kann ich heute tun?

Regeln dieser Sektion: ehrlicher Ist-Stand — was **jetzt** läuft, nicht was
geplant ist. Keine Erfolgsmeldung ohne lauffähigen Beleg.

<!-- Konkrete Befehle/Fähigkeiten. -->

- <z. B. `make gates` läuft grün>
- <z. B. Befehl X liefert Y>

## Warum <Projektname>?

<!-- Welche Lücke / welcher Schmerz? Warum existiert es, was wäre die
Alternative? Ein Absatz. -->

<…>

## Kerngedanke

<!-- Die eine Leitidee / das Designprinzip in 1–2 Sätzen. Woran sich jede
Entscheidung messen lässt. -->

<…>

## Was macht es vertrauenswürdig?

<!-- Die Harness-Signale, auf die sich Mensch und Agent verlassen. Pointer
auf die kanonischen Quellen — Inhalt nicht wiederholen. -->

- **Prozess:** [`AGENTS.md`](AGENTS.md) (Hard Rules), [`harness/README.md`](harness/README.md) (Source Precedence, Gates).
- **Verträge:** [`spec/lastenheft.md`](spec/lastenheft.md) (`LH-*`-IDs mit Akzeptanzkriterien).
- **Gates:** <welche Sensors laufen — nur existierende nennen; halluzinierte
  Gates sind die häufigste Form von Harness-Lüge, Baseline-Regelwerk
  `modul-13-quality-gates.md`>.
- **Auditierbarkeit:** Entscheidungen in `docs/plan/adr/`, Planung in `docs/plan/planning/`.
