# MR-027 — Struktur-IDs (`SPEC-*`/`ARC-*`) werden nicht vergeben

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** [`grundlagen-source-precedence.md` §ID-Schema als Klammer](../../../.harness/baseline/v5.7.0/regelwerk/grundlagen-source-precedence.md#id-schema-als-klammer)
  (Straten-Tabelle: `SPEC-<NNN>`/`ARC-<NNN>` als Struktur-IDs) und
  [`modul-03-spec.md`](../../../.harness/baseline/v5.7.0/regelwerk/modul-03-spec.md)
  (Struktur-ID-Vergabe für Sektionstypen der Spezifikation)
- **Datum:** 2026-08-21
- **Geltungsbereich:** `spec/spezifikation.md`, `spec/architecture.md`; die
  Adressierungs-Praxis von ADRs (`Schärft:`-Feld) und Slices (`Bezug:`-Feld)
- **Adaption:** d-check vergibt **keine** Struktur-IDs. Das Technik-Stratum
  wird über die gelebten **Verfeinerungs-Sektionen** adressiert (von 44
  FA-Anforderungen tragen 36 eine eigene `.a`-Sektion, 37 Sektionen inkl.
  einer `.b`); alles Übrige — §2-Schema, §4-Grund-Code-Tabelle, Defaults —
  wird nach innen per **`§`-Anker** adressiert, die Sicht
  (`spec/architecture.md`) ausschließlich per `§`-Anker. Der Baseline-eigene
  Rückfallweg („ersatzweise der Abschnitt") wird damit zur Regelform erhoben —
  das ist eine **Abweichung**, keine Auslegung, und steht deshalb hier statt
  still in der [`MR-000`](../../conventions.md#mr-000--baseline-aussage)-Aussage
  (Widerspruchs-Ausgang des Freshness-Audits, v5.4.0).
- **Begründung:** Die Adressierbarkeit, die `SPEC-*` schaffen soll, ist über
  die dichten `.a`-Verfeinerungen für die technischen Festlegungen bereits
  gegeben; ein `SPEC`-Retrofit über die gewachsene Spezifikation hätte keinen
  Konsumenten (kein Gate, keine ADR wartet auf eine Kennung, die es heute
  nicht gibt) und erzeugte genau die Pflege-ohne-Gegenwert, vor der die
  Baseline bei Zeilen-Ankern selbst warnt. `ARC-*` scheitert am selben
  Kriterium: die Architektur-Sicht ist meilensteinfrei und wird selten
  adressiert.
- **Auflösungs-Trigger:** eine ADR kann ihr `Schärft:`-Ziel nicht mehr
  **eindeutig** per `§`-Anker adressieren (zwei ununterscheidbare
  Festlegungen im selben Abschnitt) — dann ist die Struktur-ID-Vergabe für
  die betroffene Sektionsklasse einzuführen, beginnend bei den Neuzugängen.
