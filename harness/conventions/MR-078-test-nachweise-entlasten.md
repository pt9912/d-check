# MR-078 — Test-Nachweise entlasten eine Anforderung, nicht nur der Slice

- **Datum:** 2026-10-09
- **Geltungsbereich:** `trace.coverage` in der [`.d-check.yml`](../../.d-check.yml),
  `make trace` und `make completeness-check`
- **Ersetzt-Baseline-Regel:** Baseline-Regelwerk
  [`grundlagen-traceability.md` §Die zweite Richtung](../../.harness/baseline/v6.17.0/regelwerk/grundlagen-traceability.md#die-zweite-richtung-anforderung--beleg)
  — der Slice als die Zusage, die eine Anforderung schließt
- **Adaption:** Neben dem Slice schließt ein Test-Nachweis eine Anforderung:
  ein Eintrag in `docs/user/abdeckung-tests.md` (Go-Suite) oder
  `docs/user/abdeckung-e2e.md` (`make image-test`). Getragen von
  [ADR-0104](../../docs/plan/adr/0104-test-nachweise-entlasten-in-der-rtm.md).
- **Grenze:** Ein Eintrag ist eine Deklaration — eine Kennung im
  Doc-Kommentar eines Tests —, kein Beleg, dass der Test die Anforderung
  prüft.
- **Begründung:** Auftraggeber-Entscheid: die RTM soll die Test-Nachweise
  zeigen, und das Produkt kennt keine Coverage-Quelle, die zeigt, ohne zu
  entlasten.
- **Auflösungs-Trigger:** Das Produkt bekommt eine Coverage-Quelle, die zeigt,
  ohne zu entlasten, und dieses Repo stellt darauf um.
