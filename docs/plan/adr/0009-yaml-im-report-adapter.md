# ADR-0009 — `gopkg.in/yaml.v3` auch im report-Adapter (Reporter-Serialisierung)

**Status:** Accepted
**Datum:** 2026-06-19
**Autor:** pt9912
**Bezug:** [ADR-0005](0005-modul-layout-hexagon-ordner.md) (erweitert dessen
Import-Regel 3), [ADR-0001](0001-implementierungssprache.md),
[`DC-FA-CLI-004`](../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)
(YAML-Ausgabeformat, Lastenheft 0.19.0)
**Schärft:** `spec/architecture.md` §2 (Reporter-Rolle: Serialisierung der
Befunde) — die sprachkonkrete Import-Erlaubnis; nicht das Lastenheft.

## Kontext

Mit [`DC-FA-CLI-004`](../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)
0.19.0 gibt `d-check --yaml` die Befunde als YAML aus — über denselben
Reporter-Adapter `internal/adapter/driven/report`, der JSON bereits über
`encoding/json` serialisiert. [ADR-0005](0005-modul-layout-hexagon-ordner.md)
Regel 3 beschränkt `gopkg.in/yaml.v3` jedoch **ausschließlich** auf den
Config-Adapter `internal/adapter/driven/configyaml`, und `make arch-check`
erzwingt das. Zugleich verbietet Regel 5 driven-Adaptern, einander zu
importieren — der report-Adapter kann die yaml.v3-Nutzung also nicht über
`configyaml` ausleihen. Es bleibt: yaml.v3 im report-Adapter erlauben, oder
YAML von Hand emittieren.

## Entscheidung

`gopkg.in/yaml.v3` ist zusätzlich im Reporter-Adapter
`internal/adapter/driven/report` erlaubt — ausschließlich zur
**Serialisierung** der Befund-Ausgabe (Marshal), nicht zum Dekodieren.
[ADR-0005](0005-modul-layout-hexagon-ordner.md) Regel 3 wird entsprechend
erweitert: yaml.v3 in `configyaml` (Decode) **und** `report` (Encode);
sonst unverändert. `tools/arch-check.sh` wird angepasst, sodass die Fitness
Function die neue Erlaubnis abbildet (yaml.v3 weiterhin in keinem anderen
Paket, insbesondere nicht im Kern `internal/hexagon/*`).

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **yaml.v3 im report-Adapter (gewählt)** | korrekte Serialisierung (Escaping, Sonderzeichen, mehrzeilig) durch die erprobte Lib; minimaler Code; Wiederverwendung der schon vorhandenen Dependency | yaml.v3 in einem zweiten Adapter — Import-Regel muss erweitert werden |
| Hand-gerollter YAML-Emitter im report | keine Dependency-Ausweitung | YAML-Escaping (Sonderzeichen in `message`/`target`) fehleranfällig; eigener Korrektheits- und Testaufwand |
| Neuer „serializer"-Port + Adapter | strikte Trennung | Über-Zeremonie für ein einzelnes Format; zusätzliche Verdrahtung |
| Serialisierung über `configyaml` ausleihen | keine zweite yaml.v3-Stelle | verboten — Regel 5 (driven-Adapter importieren einander nicht) |

## Konsequenzen

- `make arch-check` erlaubt yaml.v3 in genau zwei Adaptern (`configyaml`,
  `report`); der Kern `internal/hexagon/*` bleibt yaml- und I/O-frei.
- YAML-Marshal eines Structs ist deterministisch (feste Feld-Reihenfolge —
  [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)), analog
  zur JSON-Ausgabe.
- Keine weitere Adapter-Tür für yaml.v3; künftige Formate (SARIF/JUnit)
  bewerten ggf. einen eigenen serializer-Port (siehe Re-Evaluierung).

## Fitness Function

`tools/arch-check.sh` (`make arch-check`): die yaml.v3-Regel listet
`configyaml` und `report` als erlaubte Pakete; ein yaml.v3-Import anderswo
bricht den Build. Bindung:
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).

## Re-Evaluierungs-Trigger

- Weitere maschinenlesbare Formate (SARIF, JUnit-XML) → einen eigenen
  serializer-Port statt weiterer Lib-Streuung über Adapter neu bewerten.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-19 | Proposed → Accepted (slice-031) |
