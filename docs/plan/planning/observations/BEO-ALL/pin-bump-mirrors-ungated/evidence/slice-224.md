**Vorgang:** slice-224
**Fund:** **Dieselben vier Spiegel-Klassen, siebtes Auftreten — und zum
ersten Mal ohne neuen `ignore-refs`-Eintrag.** Pfad-Verweise blieben
gate-gedeckt; Release-/Tree-URLs, bare Versionsnennungen und
`d-check:cite`-Direktiven wieder nicht — gefunden durch gezielte,
klassen-eigene Suchen (`baseline/v6\.6\.0` für Pfade,
`ai-harness-course/[a-z]*/v6\.6\.0` für URLs, `d-check:cite` repo-weit
für die vierte Klasse), nicht durch eine einzelne pauschale Suchform.

**Neuer Datenpunkt gegen die bisherige Annahme, jeder Bump brauche einen
`ignore-refs`-Nachtrag:** `make doc-check` meldet nach dem Entfernen des
alten Baums **0** Befunde, ganz ohne neuen Eintrag. Der Unterschied zum
Vorgänger (slice-222, ein Eintrag für vier `target-missing`): dort trugen
zwei Lauf-Belege den entfernten Baum als Markdown-**Link**; hier tragen alle
neun eingefrorenen Fundstellen den Pfad nur in Inline-Code oder Prosa, die
kein Modul auflöst. Die Klasse ist also nicht bei jedem Bump gleich groß —
sie hängt davon ab, in welcher **Form** die eingefrorenen Belege zufällig
zitieren.

**Der unabhängige Review fand eine neunte Stelle, die die eigene
Frozen-Liste ausgelassen hatte** (F-1, Report jetzt archiviert, siehe
[Slice-Stub](../../../../done/wellenlos/slice-224-baseline-v690-bump.md)):
der Kommentar bei `.d-check.yml:235` trägt `v6.6.0` als Begründung für einen
**anderen**, bereits bestehenden `ignore-refs`-Eintrag (den elften, aus dem
`v6.5.0`→`v6.6.0`-Bump) — korrekt unverändert, aber weder in der ersten
Frozen-Liste dieses Slice benannt noch retargetet. Die
Vollständigkeits-Rechnung in [`MR-072`](../../../../../../../harness/conventions.md#mr-072) ging dadurch nicht auf (8 statt 9
Frozen-Stellen bei 59 Dateien außerhalb des vendorten Baums). Gefunden durch
den Review, nicht durch die eigene Suchform — dieselbe Klasse wie beim
Vorgänger, nur an einer noch feineren Stelle: eine Konfigurationsdatei mit
genau **einem** frozen-artigen Kommentar wird von einer datei-basierten
Frozen-Liste ebenso leicht übersehen wie eine Tabellenzeile innerhalb einer
lebenden Doku-Datei (vgl. `evidence/slice-224.md` bei
`mechanical-id-rewrite-misses-frozen-classes`).

**Weiterhin kein formgültiger Ausgang** — die mechanische Form
`versions.patterns` existiert seit slice-122, bleibt aber unscharf
geschaltet. Siebtes Auftreten, Zähler-Stand unverändert bei
[`BEO-ALL/registerzeile-ohne-ausgang-nach-schwelle`](../../registerzeile-ohne-ausgang-nach-schwelle/observation.md).
