# MR-029 — Baseline-Pin-Hebung auf v5.9.0 (fünfter Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** — *(keine; Pin-Fortschreibung innerhalb des von
  [`MR-023`](../conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
  festgelegten self-contained Bundle-Layouts)*
- **Datum:** 2026-08-22
- **Geltungsbereich:** [§Baseline](../conventions.md#baseline), [§Adoptierte
  Konventions-Quellen](../conventions.md#adoptierte-konventions-quellen), die
  pin-gebundenen Verweise
  ([`MR-021`](../conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden))
  in [`AGENTS.md`](../../AGENTS.md), [`harness/README.md`](../README.md), den
  aktiven `MR-*`-Dateien, [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md)
  und den Planning-Docs
- **Adaption:** Der Baseline-Pin ist von `v5.7.0` auf **`v5.9.0`** gehoben
  (Kurs-Tags vom 2026-08-22) — die von
  [`MR-011`](../conventions.md#mr-011--baseline-auf-release-tag-gepinnt)
  vorgesehene Fortschreibung, fünfter Nachtrag der Serie; ersetzt
  [`MR-028`](done/MR-028-baseline-v570.md) nach dessen eigenem
  Auflösungs-Trigger. Kein Layout-Wechsel: dasselbe self-contained Bundle
  (`lab-regelwerk.zip`, beide Bäume, `SHA256SUMS`, 51 Dateien — gemessen),
  dasselbe Materialisierungs-Skript, unverändertes Pfadschema
  `.harness/baseline/<tag>/{regelwerk,templates}/`.

  **Das Delta sind zwei Stufen, gemessen statt gezählt:** 33 Bundle-Dateien
  unterscheiden sich; **22 davon ändern ausschließlich den Quell-URL-Stempel**
  ihrer Kopfzeile, elf haben Änderungen im Rumpf. Von diesen elf trägt eine
  (`conventions.template.md`) im Rumpf ebenfalls nur einen Versions-Zeiger
  (die Download-URL) — **zehn Dateien tragen also echten Regel-Inhalt**: vier
  im Regelwerk (`grundlagen-harness-dateien` +34/−1, `modul-03-spec` +11,
  `modul-06-roadmap` +5/−5, der Index) und sechs Templates. Sie tragen
  **eine** Regel: ein Feld, das einen *Zustand*
  trägt, nennt Zustand und Beleg als auflösbaren Anker, nicht die Chronik; die
  Kopfzeile eines lebenden Registers entfällt; das Technik-Stratum trägt kein
  Kopf-Datum und keinen Kopf-Status, die Sicht behält ihren Frische-Marker; das
  Drift-Log führt nur Umplanungen. Der Reviewer-Skill der Vorlage bekommt dafür
  einen HIGH-Eintrag.

  - **Hebungs-Zensus (Checkliste für den Nachfolger):** die Spiegel einer
    Pin-Hebung sind **drei Klassen**, nicht eine — (1)
    `baseline/<tag>`-**Pfad**-Verweise (grep-bar, gate-gedeckt), (2)
    Release-/Tree-**URLs** mit dem Tag (`releases/tag/`,
    `releases/download/`, `tree/` — kein Gate deckt sie), (3)
    **Prosa-/Ellipsen-Pins** (`…/vX.Y.Z/…`, „Stand"-Angaben in MR-Körpern).
    Alle drei sind hier einzeln durchgegangen worden; Klasse 3 traf genau eine
    lebende Stelle (die Ellipse in
    [`MR-021`](../conventions.md#mr-021--in-repo-verweise-auf-das-vendored-regelwerk-sind-pin-gebunden))
    und eine Sektionsregel der Roadmap. **Historische** Nennungen des alten
    Tags in `done/`-Dokumenten, Review-Reports und Historie-Zeilen bleiben —
    sie sagen, was damals galt.
  - **Alter Baum entfernt, Historie via Tombstone.** `.harness/baseline/v5.7.0/`
    ist entfernt (ein Pin, eine netzlose Lese-Form). Der einzige **eingefrorene
    Link**-Bestand darauf — die aufgelöste
    [`MR-027`](done/MR-027-struktur-id-verzicht.md) — ist über das geteilte
    Referenz-Ventil
    ([`DC-FA-REF-001`](../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus))
    quell-skopiert ausgenommen; die Review-Reports nennen den Pfad nur in
    Inline-Code und sind dort ohnehin ausgenommen.
  - **Der Konformitäts-Abgleich ist nicht Teil dieser MR:** die eine Regel des
    Deltas wird in eigenen Slices verkörpert und angewandt. Diese MR hebt den
    Pin, sie behauptet keine Konformität.
- **Begründung:** Auftraggeber-Anstoß 2026-08-22 („danach können wir auf das
  neue Regelwerk v5.9.0 migrieren"); der Baseline-Default sticht die
  repo-lokale Adaption. Vendored wird das **Release-Asset am Tag**
  (`--check-latest` ist die Currency-/Authentizitäts-Gegenprobe), nicht der
  Kurs-Arbeitsbaum.
- **Löst auf:** [`MR-028`](../conventions.md#mr-028) *(der Verweis geht auf die
  Index-Zeile, nicht auf die Eintrags-Datei — die wandert bei Auflösung nach
  `done/`, und ein Pfad-Link bräche genau dann)*
- **Ausgelöst durch Baseline-Stand:** v5.9.0
- **Auflösungs-Trigger:** die nächste Pin-Hebung ersetzt diesen Eintrag durch
  ihren Nachfolger — wie [`MR-028`](done/MR-028-baseline-v570.md) durch diesen
  Eintrag ersetzt wurde.
  [`MR-023`](../conventions.md#mr-023--baseline-pin-hebung-auf-v500-samt-self-contained-bundle-layout)
  bleibt daneben **aktiv** stehen: es trägt das Bundle-Layout, nicht den
  Pin-Wert.
