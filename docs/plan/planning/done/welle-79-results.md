# Welle 79 — Baseline v5.7.0 + das eigene Prädikat der Listen-Hälfte — Closure-Notiz

**Welle:** welle-79-zwei-haelften-ein-waechter
**Abschluss:** 2026-08-21
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- **Bump** ([slice-110](welle-79/slice-110-baseline-v570-bump.md)): Baseline
  v5.6.0 → **v5.7.0** vendored (`--verify` 51 Dateien, `--check-latest`
  beidseitig OK), [`MR-028`](../../../../harness/conventions.md#mr-028)
  als vierter Pin-Nachtrag mit aufgelöster
  [`MR-026`](../../../../harness/conventions.md#mr-026), Verweis-Hebung
  über **alle drei Spiegel-Klassen** (die dritte — URL-/Ellipsen-Pins —
  erst auf Review-Auflage, zum zweiten Mal in Folge), zwei exakt
  skopierte Tombstones, die Zwei-Hälften-Prosa der Roadmap, Delta-Audit
  über fünf Bundle-Dateien.
- **Produkt** ([slice-111](welle-79/slice-111-wave-drift-zwei-haelften.md)):
  `planning.waves.mode: one | many` auf **formalen Konsumenten-CR**
  („Bijektion statt Singleton", ai-harness-course) — der CR landete als
  reiner Lastenheft-Commit (0.62.0), die Entscheide als
  [ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md)-Fortschreibung,
  dann `waveBijection` über die geteilte Block-Grenze (`target` =
  Kennung, Marker außen vor, Default byte-identisch mit
  Live-Beweis gegen das gepinnte Alt-Image). **Release v0.62.0**
  (Digest `sha256:3996a593…4cacf`), danach die eigenen Profile auf
  `many`.
- **Neben den Slices, auf Auftraggeber-Audits desselben Tages:** AGENTS
  §1 (Nachschlag-Doktrin + Pflicht-Blick-Trigger), §3.3 auf das
  [`MR-013`](../../../../harness/conventions.md#mr-013)-Muster (samt
  Erweiterung des Eintrags um MR-/Wellen-Moves), §3.7 auf die
  Baseline-Feld-Formen geschärft, §5-Zeiger auf
  [ADR-0027](../../adr/0027-commits-traceability-modul.md); zwei
  Heilungen an den eigenen Slice-Köpfen (Verantwortlich, Vorprüfungen);
  **BEO-007 verkörpert** (der `pre-commit`-Hook trägt jetzt den vollen
  `doc-check`).

## Was hat funktioniert?

- **Die Konsumenten-Schleife trug in beide Richtungen:** beide eigenen
  Upstream-Notizen aus welle-78 sind in der Kurs-Welle 81 gelandet
  (modul-06 „zusätzlich", modul-10 `klasse`), und der Rück-CR des
  Konsumenten kam mit Messung (team-sim, 11/11 PASS) und einem
  Landungs-Protokoll, das exakt dem eigenen Prozess entsprach — der
  Zuschnitt von slice-111 wurde am offenen Slice nachgezogen statt
  verteidigt.
- **Die Lieferung bewies sich am eigenen Prozess:** die Closure dieser
  Welle lief selbst durch das Fenster „Marker zusätzlich zum Zeiger",
  das vor der Welle gate-rot gewesen wäre — der erste gelebte
  `many`-Zustand ist der Schluss-Zustand des liefernden Slice.
- **Review-Live-Proben statt Behauptungen:** der stärkste Beleg der
  Welle war die byte-identische Gegenprobe des `one`-Defaults gegen das
  gepinnte v0.61.0-Image — grün **und** rot; dazu neun Fixture-Proben
  am frisch gebauten Image (Exit 2 inkl. explizit leerem Modus).

## Was ging anders als geplant?

- **Das Register ist ehrlicher geworden, nicht leerer:** BEO-007 trat am
  Tag der gelebten Arbeitsregel ein **drittes** Mal ein (der Gate-Exit
  stand in der Ausgabe, die `;`-Kette band ihn nicht) und ist
  **verkörpert** — nicht als Regel, sondern als Hook, der den Commit
  blockiert. **BEO-009 neu bei 2** (Commit-Botschaften behaupteten eine
  fehlgeschlagene Probe bzw. einen still verfehlten Edit als erledigt —
  einmal gepusht, im Folge-Commit korrigiert). **BEO-008 neu bei 2**
  (Pin-Spiegel-Klassen, jetzt Checkliste in
  [`MR-028`](../../../../harness/conventions.md#mr-028)).
- **Der CR traf mitten in der Umsetzung ein** und ersetzte den geplanten
  Default-Umbau durch das opt-in-Modell — weniger invasiv, Default-treu,
  und mit fertigem Vertrags-Text; der ursprüngliche slice-111-Schnitt
  hätte konsumentensichtbar liberalisiert.
- **Geplant zwei Slices, geliefert zwei Slices** — der Wellen-Schnitt
  trug unverändert durch; Eröffnung und Closure am selben Tag.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 4 — was aus dieser Welle in die
Steuerung zurückfließt.

- **Ein expliziter Exit ohne Bindung ist Dekoration.** Zweimal stand der
  rote Gate-Exit sichtbar in der Ausgabe und band nichts. Wächter
  gehören dorthin, wo sie den Vollzug blockieren (Hook, Gate), nicht in
  die Erzähl-Kette des Arbeitenden — die BEO-007-Verkörperung ist die
  Konsequenz.
- **Eine Botschaft über eine Probe entsteht nach deren Exit.** Zwei
  Botschaften beschrieben den beabsichtigten statt den eingetretenen
  Zustand (BEO-009). Vor dem Commit gehört jede Botschafts-Zeile gegen
  `git diff --stat` bzw. das Proben-Ergebnis gehalten.
- **Kommentare beschreiben, was da ist — mit Baseline-Feld-Formen.** Die
  Auftraggeber-Korrektur (keine Befund-Marker, Slice-Nummern,
  Mess-Labels in Code-/Test-Kommentaren) ist jetzt Regeltext in §3.7,
  nicht nur Gedächtnis.

## Beobachtungs-Register (Zeiger)

Der Lese-Schritt dieser Welle: **BEO-007 hat 3× erreicht und ist
verkörpert** (der `pre-commit`-Hook trägt den vollen `doc-check`; die
Rest-Klasse — Gates jenseits von `doc-check` — bleibt benannt).
**BEO-006, BEO-008 und BEO-009 stehen bei je 2** mit gelebtem
Gegenmittel und benannter mechanischer Form; BEO-002/003/004 bleiben
verkörpert, BEO-001/005 gestrichen. Die Verkörperungen haben getragen:
[`MR-025`](../../../../harness/conventions.md#mr-025) erzwang die
Spiegel-Liste des `mode`-Umbaus, der BEO-004-Anker blieb ohne neuen
Fund.

## Folge-Slices

- **Keiner.** `open/` ist leer, §Nächste Wellen trägt `— keine —`. Die
  benannten Grenzen warten mit beobachtbaren Triggern: der
  Lastenheft-Wortlaut „Menge der flachen Wellendokumente" meint die
  Kennungs-Menge (Präzisierungs-Kandidat der nächsten
  Lastenheft-Redaktion), der Mehr-Wellen-**Betrieb** dieses Repos ist
  ein eigener Roadmap-Entscheid (die Bijektion macht ihn prüfbar), und
  die drei offenen Beobachtungen tragen ihre 3×-Formen im Register.
- **Toolchain-Kandidat außerhalb der Welle:** Go 1.27.0 ist erschienen;
  der Bump des Build-Args ist ein eigener kleiner Schritt auf
  Auftraggeber-Freigabe.

## Verifikation

- **Closure-Trigger erfüllt:** beide Slices in `done/`; Pin v5.7.0
  vendored und beidseitig auditiert; kein lebender Verweis nennt
  `baseline/v5.6.0` (eingefrorene getombstoned, `make doc-check` belegt
  beides); Delta-Audit je geänderter Regel; unter `mode: many` beide
  Bijektions-Richtungen als Tests rot und die drei legitimierten
  Zustände grün belegt, der Default byte-identisch; **Release v0.62.0**
  auf GHCR samt Digest-Backfill; die eigenen Profile laufen auf `many`
  (Gates + GUARD lokal und mit dem released Digest-Image grün);
  `make fullbuild` grün, Image-Hash `sha256:56d27f06…2965`.
- **Trigger-Audit** über die drei Artefaktklassen: keine offenen
  Carveouts, keine stehengebliebene Gate-Reifestufe, kein eingetretener
  Re-Evaluierungs-Trigger — geprüft auch für
  [`MR-028`](../../../../harness/conventions.md#mr-028) (nächste
  Pin-Hebung), [`MR-027`](../../../../harness/conventions.md#mr-027)
  (Schärft-Eindeutigkeit) und die fortgeschriebene
  [ADR-0055](../../adr/0055-wellen-invariante-artefakt-und-grund-codes.md)
  (dritte Kardinalitäts-Semantik, Bijektions-Prädikat der Register) —
  keiner fällig.
- **Zwei unabhängige Frischkontext-Reviews** (je APPROVE mit Auflagen,
  0 HIGH; alle Auflagen vor Closure bzw. Tag eingearbeitet), dazu die
  Auftraggeber-Audits von AGENTS §1/§3.3/§3.4/§3.7/§5 mit vier
  Nachzügen und zwei Slice-Kopf-Heilungen.
