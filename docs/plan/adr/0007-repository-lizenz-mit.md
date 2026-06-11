# ADR-0007 — Repository-Lizenz: MIT

**Status:** Accepted
**Datum:** 2026-06-11
**Autor:** pt9912
**Bezug:** Repository-Veröffentlichung und Weiterverwendung
**Schärft:** keine Spec-Stelle — Projektmetadaten-ADR; verbindlich für
[`LICENSE`](../../../LICENSE) und den Lizenzhinweis im
[`README.md`](../../../README.md).

## Kontext

`d-check` wird als öffentliches Repository und später als GHCR-Image
bereitgestellt. Ohne explizite Lizenz ist die rechtliche
Weiterverwendung durch andere Repositories und CI-Umgebungen unklar,
auch wenn der Quellcode sichtbar ist.

## Entscheidung

Das Repository steht unter der MIT-Lizenz. Der kanonische Lizenztext
liegt in [`LICENSE`](../../../LICENSE); das
[`README.md`](../../../README.md) verweist darauf.

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **MIT (gewählt)** | permissiv, kurz, im Open-Source-Ökosystem breit verstanden | keine Copyleft-Pflichten für Ableitungen |
| Apache-2.0 | explizite Patentregelung | umfangreicher, für dieses kleine CLI-Tool schwerergewichtig |
| Keine Lizenz | keine Entscheidung nötig | praktisch keine klare Nutzungs- und Weitergaberechtslage |

## Konsequenzen

- GitHub kann das Repository als MIT-lizenziert erkennen, sobald die
  Lizenzdatei auf dem Default-Branch liegt.
- Konsumenten dürfen `d-check` und die Dokumentation unter den
  Bedingungen der MIT-Lizenz nutzen, kopieren, ändern und weitergeben.
- Eine spätere Lizenzänderung ist eine neue ADR und muss die
  Kompatibilität bisheriger Veröffentlichungen berücksichtigen.

## Fitness Function

`make doc-check` prüft die Links auf [`LICENSE`](../../../LICENSE) und
[`README.md`](../../../README.md).

## Re-Evaluierungs-Trigger

- Ein konkreter Bedarf nach zusätzlicher Patentregelung oder
  Copyleft-Pflichten entsteht.
- Ein Konsumenten-Repo verlangt eine andere Lizenzkompatibilität.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-11 | Proposed → Accepted |
