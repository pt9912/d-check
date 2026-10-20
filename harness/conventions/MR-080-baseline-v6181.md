# MR-080 — Baseline-Pin-Hebung auf `v6.18.1` (siebzehnter Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Datum:** 2026-10-20
- **Geltungsbereich:** §Baseline, pin-gebundene Verweise,
  `.harness/baseline/v6.18.1/`
- **Ersetzt-Baseline-Regel:** — *(Pin-Fortschreibung im Bundle-Layout des
  Vorgängers; keine inhaltliche Abweichung)*
- **Adaption:** Der vendorte Bestand steht auf
  [`v6.18.1`](https://github.com/pt9912/ai-harness-course/releases/tag/v6.18.1).
  Ein Zwischen-Tag (`v6.18.0`) liegt zwischen dem alten Pin und diesem — ein Bump
  zielt auf den jeweils aktuellen Tag, nicht auf jeden dazwischen.

  **Der Delta, gemessen und nicht geschätzt** (Vergleich beider Bäume ohne die
  `<!-- Quelle:`-Zeilen und mit normalisierten Versionsstempeln, beide Bäume noch
  vorhanden): 55 Dateien vorher wie nachher, keine neu, keine entfallen. **8**
  tragen ein inhaltliches Delta (63 Diff-Zeilen), davon 1 Regelwerk-Datei und 7 mit
  Vorlagen bzw. dem Bundle-Index; alle übrigen unterscheiden sich nur im
  Versionsstempel und in den Quell-Verweisen, dazu `SHA256SUMS`. Größte
  Bewegungen: `templates/.d-check.yml` mit 17 Diff-Zeilen,
  `regelwerk/modul-10-review-harness.md` mit 15,
  `templates/harness/README.template.md` mit 11 und
  `templates/docs/reviews/review-report.template.md` mit 10.

  **Inhaltliche Bewegungen — vollständig nach Datei, keine davon übernommen**
  (Abgrenzung, slice-277; Adoptieren ist Sache eigener Folge-Slices):

  1. **Review-Deckung: der Sensor prüft nur, was er als Zusage erkennt** —
     `modul-10-review-harness.md` (neuer Absatz: ein Deckungs-Sensor, dessen Muster
     die Review-Zeile der Slice-Vorlage nicht trifft, läuft grün über einer leeren
     Menge; er meldet deshalb den Leerlauf, ordnet einen Report dem Slice über
     dieselbe Kennung zu, die der Slice trägt, nimmt archivierte Slices aus und
     nennt die leere Menge nach dem Archivieren einen Ruhezustand, keinen
     Leerlauf; dazu: die volle Slice-Kennung im Dateinamen des Reports).
     Nachgezogen in `templates/.d-check.yml` (Block `reviews`: `match: name`,
     `require-promises`, `recursive`, `skip-pattern`, `skip-allows-empty`, mit der
     Erklärung der Grenze des `skip-pattern`),
     `templates/harness/README.template.md` (Einsatz ab d-check v0.86.0, „gehakt oder
     nicht") und `templates/docs/reviews/review-report.template.md` (Dateiname mit der
     **vollen** Slice-Kennung, Folgeläufe mit Suffix `-r2`). **Dieses Repo ist
     betroffen, und die Bewegung geht auf eigene Befunde zurück:** sie ist die
     Fassung der Change Requests dieses Repos zum Modul `reviews`. Der
     `reviews`-Block der [`.d-check.yml`](../../.d-check.yml) trägt von den Schlüsseln der
     Vorlage nur `recursive` und `skip-pattern`; `require-promises`, `match` und
     `skip-allows-empty` fehlen. `make review-coverage` meldet 0 Befunde; **nicht
     gemessen** ist, ob der Lauf eine nichtleere Menge von Zusagen sieht (der
     Leerlauf-Fall, den die Vorlage benennt). Kandidat für einen Folge-Slice.
  2. **Vorlagen-Pfade nennen die vendorte Baseline** — `templates/.harness/skills/
     reviewer.template.md`, `closure-note-reviewer.template.md`,
     `templates/harness/conventions.template.md` und
     `templates/harness/README.template.md`: aus dem blanken Pfad der Vorlage unter `docs/`
     wird der Pfad `.harness/baseline/<tag>/templates/docs/…` mit dem Baum-Präfix
     (entsprechend für die Slice-, MR- und Sensor-Vorlage). Reine Pfad-Form der
     Verweise im Skelett; dieses Repo führt seine Skills mit eigenen Pfaden
     (`.harness/skills/`) und ist nicht betroffen.
  3. **Bundle-Stand** — `regelwerk/README.md` (Kurs-Welle 160 → 161, Datum
     2026-10-07 → 2026-10-10).

  **Spiegel-Klassen, gemessen am echten Vorzustand** (HEAD vor dem Swap, ohne
  `.harness/baseline/`, **nicht** ohne `.harness/` — die Skills liegen dort):
  104 Dateien außerhalb des vendorten Baums nannten `v6.17.0` — 262 Vorkommen.
  **Frozen, byte-stabil (51 Dateien, 156 Vorkommen):** die `done/`-Slices und
  -Wellen (24 / 68), die Review-Reports (25 / 85), der eingehende CR (1 / 1) und die
  CHANGELOG-Historie (1 / 2); `Accepted`-ADRs und `conventions/done/` nannten den Tag
  nicht. **Lebend retargetet (49 Dateien):** `AGENTS.md`, `harness/README.md`,
  `harness/conventions.md`, die drei `harness/rules/`-Dateien, 38 aktive `MR-*`,
  `.claude/agents/reviewer.md`, die beiden Skills unter `.harness/skills/`,
  `roadmap.md`, `planning/README.md`, `observations/README.md`,
  `spec/architecture.md` und `spec/spezifikation.md` (Rolle-Zeile). **Zeilenweise
  getrennt:** `observations/README.md` und `harness/conventions.md` (Zeilen mit
  Pin-Aussage retargetet, Vergangenheitsaussagen nicht);
  [`MR-075`](../conventions.md#mr-075) („Seit `v6.17.0` trägt die Spezifikation einen
  Abschnitt …") bleibt stehen — die Aussage ist eine Tatsache über die Fassung, in der
  die Umnummerierung kam, nicht ein Pin; ebenso der Beleg
  `evidence/slice-261.md` im Register (unveränderlich ab Merge) und dieser
  Eintrag selbst. **Die 8 Baseline-Symlinks unter `.claude/rules/`** lösen jetzt
  gegen `v6.18.1`. **Eine Ausnahme von der Frozen-Regel:** `slice-274` (am selben
  Tag geschlossen) trägt einen Markdown-Link auf den `v6.17.0`-Baum
  (`modul-11` §Bewusstes Brechen); er ist auf `v6.18.1` retargetet, weil ein toter
  Link in einer eingefrorenen Datei sonst `ignore-refs` wachsen ließe
  ([`MR-069`](../conventions.md#mr-069)). Der Text des Slice ist unberührt.

  **`ignore-refs` wächst nicht** — gemessen: nach dem Entfernen des alten Baums
  meldet `doc-check` außer dem eben genannten Link und den Folgen des Moves von
  `MR-074` nach `conventions/done/` (Link-Tiefen, eigener Commit) und der
  verrutschten Zitat-Spanne (unten) 0 Befunde.

  **`d-check:cite`-Neu-Ankern** — gemessen **19** lebende Direktiven auf den
  Baum (0 tag-only), **1** davon neu geankert, reine Zeilenverschiebung bei
  wortgleichem Text: Reviewer-Skill `modul-10-review-harness.md` 87→99. Die übrigen
  18 zeigen auf Dateien ohne inhaltliches Delta oder auf unverschobene Zeilen.

  **Zitat-Delta** ([`MR-039`](../conventions.md#mr-039)): wörtliche Zitate ohne
  Direktive in Markdown stammen aus `modul-13` (`AGENTS.md`, „Halluzinierte Gates
  …") und `modul-05` (`MR-056`, §Lifecycle) — beide Dateien tragen kein inhaltliches
  Delta, die Zitate stehen im `v6.18.1`-Wortlaut unverändert. Aus den acht
  Delta-Dateien zitiert dieses Repo außerhalb der `d-check:cite`-Spannen nicht
  wörtlich (gemessen: die Bezüge sind Pfad-Verweise).
- **Begründung:** Der Nachtlauf `upstream-drift` meldete den neuen Release
  (`make baseline-freshness`); das Release `v6.18.1` lag vor. Die Hebung ist eine
  reine Fortschreibung; die inhaltlichen Bewegungen (oben) warten auf eigene
  Folge-Slices. Beim Vendor scheiterte die Under-Copy-Barriere von
  `fetch-baseline-cache.sh` an einem String-Vergleich der Dateizahl unter macOS (BSD-`wc`
  füllt auf); die Zahl stimmte (54 = 54), der Vergleich ist numerisch.
- **Auflösungs-Trigger:** die nächste Pin-Hebung.
- **Löst auf:** [`MR-074`](../conventions.md#mr-074)
- **Ausgelöst durch Baseline-Stand:** `v6.18.1`
