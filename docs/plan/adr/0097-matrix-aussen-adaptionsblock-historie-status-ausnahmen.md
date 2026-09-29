# ADR-0097 — Matrix: aussen- und adaptionsblock-Adoption mit Historie- und Status-Ausnahmen

**Status:** Accepted

**Datum:** 2026-09-29

**Betroffen:** `.d-check.yml` (Modul `matrix`), `spec/lastenheft.md`,
`spec/spezifikation.md`, `spec/architecture.md`

## Kontext

Die Schwester (ai-harness-init) geht mit drei `.d-check.yml`-Positionen über
die Baseline-Matrix hinaus: die Klassen `aussen` (First-Match, letzte —
`paths: ["**"]`), `adaptionsblock` (Konventionsablage samt Token
`MR-\d{3}`) und die zugehörigen Straten-Regeln
(`spec-straten→aussen`, `spec-straten→adaptionsblock`, jeweils auch für die
Sicht). slice-246 hat die Positionen evaluiert und die Adoption an <!-- d-check:status-provenance -->
slice-248 delegiert — die Probe maß 40 Befunde gegen den Bestand. <!-- d-check:status-provenance -->

Die Referenz-Richtung ([`MR-006`](../../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)) verlangt: die Spec-Straten verweisen nie
abwärts. Die Matrix bewacht diese Richtung heute für ADRs, Slices, Wellen
und Commit-Hashes — aber nicht für Harness-, Planning- und
Betriebs-Dokumente. Der gemessene Bestand trägt 28 solche Links: in den
Konventionsspeicher (14, von der ids-Link-Pflicht erzwungen), in das
Agenten-Briefing (6), in Baseline-Zitate (2), in Harness-README, Packaging,
Carveout, Register und Change Requests (je 1–2). Dazu 11 `matrix-inactive`-
Funde, weil die `aussen`-Klasse jede Datei in die Status-Prüfung zieht
(ADR-Index ×7, CHANGELOG ×3, Release-Doku ×1).

Zwei Ausnahmen sind funktional begründet und brauchen eine ADR (§3.6):
der ADR-Index führt superseded ADRs — das ist seine Aufgabe; die
Lastenheft-Historie und die Spezifikations-Chronik sind **Zeitdokumente**
(jede Zeile friert den Stand ihres Datums), deren Verweise zum
Zeitzpunkt-Text gehören.

## Entscheidung

1. Die Klassen `adaptionsblock` und `aussen` werden adoptiert — in dieser
   Reihenfolge (First-Match: `aussen` fasst alle Pfade und muss **letzte**
   Klasse sein). Die Regeln `spec-straten→aussen`, `sicht→aussen`,
   `spec-straten→adaptionsblock` und `sicht→adaptionsblock` (je
   `allow: false`) schließen die Referenz-Richtung nach außen.
2. Die Abschnitte `## 7. Historie` beider Spec-Straten werden in
   `matrix.exclude-sections` aufgenommen. Begründung: sie sind Zeitdokumente;
   ihre Zeilen tragen Provenance zum jeweiligen Datum (MR-Verweise,
   CR-Verweise — der Kanon lässt in der Lastenheft-Historie den externen
   Change Request ausdrücklich zu, v5.11.0). Die Ausnahme ist
   heading-genau (`7. Historie`) — die ADR-Geschichte bleibt separat
   ausgenommen, die Referenz-Richtung gilt im übrigen Dokumentkörper
   unverändert.
3. `matrix.exempt-paths` nimmt die drei gemessenen funktionalen Verweiser
   file-weit auf: `docs/plan/adr/README.md` (führt superseded ADRs), das
   `CHANGELOG.md` (Historie) und `docs/user/releasing.md` (verweist auf den
   zum Schreibzeitpunkt aktiven Stand). Das Produkt führt keine
   Status-Spitzen-Ausnahme — der file-weite Ausschluss ist die bestehende
   Mechanik; file-weit ist hier kein Verlust, denn keine der drei Dateien
   trägt eine andere matrix-überwachte Referenz.

Die lebenden Textstellen (13 MR-Provenance-Links, 6 Agenten-Briefing-Links,
2 Baseline-Kopf-Zitate, 3 Einzelfälle) werden auf die Regel hin bereinigt:
Baseline-Zitate in der Zitier-Form (Text statt Link), Provenance-Verweise
entfallen (die Herkunft lebt im Konventionsspeicher und in der Historie).

## Konsequenzen

- Die Referenz-Richtung der Spec-Straten ist vollständig mechanisiert:
  innen, nach oben, nach außen — jede Richtung hat ihre Regel.
- Die Status-Prüfung gilt repo-weit; die drei Ausnahmen sind per ADR
  gesichert. Eine neue Datei mit `Status:`-Feld fällt unter die Prüfung.
- Der shallow-Klon bleibt eine benannte Grenze: der
  Erreichbarkeits-Walk der leere-Range-Prüfung
  (§[`DC-FA-VCS-002`](../../../spec/lastenheft.md#dc-fa-vcs-002--leere-commit-range-im-modul-vcs-ist-laut-zu-melden-opt-in).a) braucht die volle Historie — Abhilfe
  `fetch-depth: 0`.

## Re-Evaluierungs-Trigger

- Trägt ein Spec-Stratum einen neuen Verweis nach außen, den die Regeln
  verbieten, und der Bestand hat sich um ≥ 5 Befunde derselben Klasse
  gehäuft, wird geprüft, ob eine weitere Klasse (statt
  Einzelausnahmen) die Lage mechanisiert.
- Löst die Baseline die `aussen`-First-Match-Form ab (z. B. durch ein
  je-Klasse-Scoping der `exclude-sections`), wird diese Adoption gegen den
  neuen Stand neu bewertet.
