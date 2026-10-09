# ADR-0104: Test-Nachweise entlasten eine Anforderung in der RTM

**Status:** Proposed

**Datum:** 2026-10-09

**Autor:** pt9912

**Bezug:** [`DC-FA-COV-001`](../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in),
[`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code);
[ADR-0026](0026-completeness-in-product-gate.md); `AGENTS.md` §3.6;
slice-270 <!-- d-check:status-provenance -->.

**Schärft:** [`DC-FA-COV-001`](../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
— welche Quelle in diesem Repo eine Anforderung von der Waise entlastet.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`make trace` zeigte je Anforderung ADRs und Slices, aber keine Tests — obwohl
die Go-Suite und `make image-test` die meisten Anforderungen prüfen. Die
Schwester-Repos binden ihre Test-Nachweise über `trace.coverage` ein; die RTM
zeigt sie dann als Spalte.

`trace.coverage` macht eine Anforderung zugleich waisenfrei
([`DC-FA-COV-001.a`](../../../spec/spezifikation.md#dc-fa-cov-001a--kuratierte-coverage-quellen-tracecoverage)
Schritt 5): Waise ist nur noch, was weder einen Slice noch einen
Coverage-Eintrag hat. Für `make completeness-check`
([ADR-0026](0026-completeness-in-product-gate.md)) heißt das: eine
Anforderung mit Test-Nachweis und ohne Slice besteht. Das lockert eine
Prüfregel; `AGENTS.md` §3.6 verlangt dafür eine ADR. Der Baseline-Vorschlag
lässt nur den Slice entlasten (Baseline-Regelwerk `grundlagen-traceability.md`
§Die zweite Richtung) und verlangt, eine andere Wahl zu deklarieren.

## Entscheidung

1. Die RTM dieses Repos bindet zwei Abdeckungs-Dateien über `trace.coverage`
   ein: `docs/user/abdeckung-tests.md` (Label `Tests`, die Go-Suite) und
   `docs/user/abdeckung-e2e.md` (Label `E2E`, `make image-test`).
2. Beide Dateien sind aus den Testquellen abgeleitet — eine Kennung im
   Doc-Kommentar einer Testfunktion, ein Anker `# abdeckung:` unter einer
   Phase von `tools/image-test.sh`. `make abdeckung` schreibt sie, ein Test in
   `make test` hält sie gegen die Ableitung.
3. Eine Anforderung mit einem Eintrag in einer der beiden Dateien ist damit
   auch ohne Slice keine Waise.

## Verglichene Alternativen

| Alternative | Warum nicht |
|---|---|
| Keine Test-Spalte in der RTM | Die RTM zeigt dann nicht, was geprüft wird, und der Auftraggeber fragt danach |
| Eine Spalte, die nicht entlastet | Das Produkt kennt eine solche Quelle nicht; sie wäre eine Produkt-Erweiterung vor diesem Slice (Auftraggeber-Entscheid: annehmen) |
| Von Hand gepflegte Abdeckungs-Dateien | Sie driften, sobald Tests sich ändern, und kein Gate merkt es |

## Konsequenzen

- Die Zusage, dass jede Anforderung einen Slice hat, gibt es nicht mehr; ein
  Slice oder ein Test-Nachweis schließt sie. Heute ohne Wirkung: die RTM
  meldet 0 Waisen.
- **Grenze:** Ein Eintrag ist eine Deklaration, kein Beleg — eine Kennung im
  Doc-Kommentar sagt, dass der Test die Anforderung prüfen soll, nicht, dass
  er es tut. Eine Anforderung kann damit über eine einzige Kommentarzeile
  waisenfrei werden.

## Fitness Function (falls maschinell prüfbar)

`TestAbdeckungsDateienFolgenIhrerAbleitung` hält beide Dateien gegen ihre
Ableitung; die Negativliste der Ableitung steht in
`TestGoTestZeilenZaehltNurDenDocKommentarEinerTestfunktion` und
`TestImageTestZeilenVerlangtDenAnkerJederPhase`.

## Re-Evaluierungs-Trigger

- Das Produkt bekommt eine Coverage-Quelle, die zeigt, ohne zu entlasten.
- `make completeness-check` lässt eine Anforderung durch, deren einziger
  Nachweis ein Test ist, der sie nicht prüft.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-10-09 | Proposed (Auftraggeber-Entscheid zum Schnitt von slice-270) |
