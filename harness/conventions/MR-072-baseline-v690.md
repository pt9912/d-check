# MR-072 — Baseline-Pin-Hebung auf `v6.9.0` (vierzehnter Nachtrag zu MR-011, Nachtrag zu MR-023)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** — *(Pin-Fortschreibung im Bundle-Layout des
  Vorgängers; keine inhaltliche Abweichung)*
- **Datum:** 2026-09-16
- **Geltungsbereich:** §Baseline, pin-gebundene Verweise,
  `.harness/baseline/v6.9.0/`
- **Adaption:** Der vendorte Bestand steht auf
  [`v6.9.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v6.9.0).
  Fünf Zwischen-Tags (`v6.7.0`, `v6.7.1`, `v6.7.2`, `v6.8.0`) liegen zwischen
  dem alten Pin und diesem — ein Bump zielt auf den jeweils aktuellen Tag,
  nicht auf jeden dazwischen (Präzedenz: `v5.7.0`→`v5.9.0`,
  `v6.0.0`→`v6.3.1`, u. a.).

  **Der Delta, gemessen und nicht geschätzt** (`diff -I '<!-- Quelle:'`): 54
  Dateien vorher wie nachher, keine neu, keine entfallen. **16** von 26
  Regelwerk-Dateien tragen ein inhaltliches Delta (Schlagzeile: `grundlagen-
  source-precedence.md` mit 140 Diff-Zeilen, `modul-05-planning-harness.md`
  mit 111), dazu **15** von 28 Templates.

  **Zwei inhaltliche Schlagzeilen, keine davon übernommen** (Abgrenzung 1):
  Erstens generalisiert `grundlagen-source-precedence.md` §Vergabe das
  ID-Schema — Kennungen als **Namen statt Nummern** für Slice/Welle,
  Mehrfach-Schreiber-Kollisionsbehandlung, Bereichssegmente für ADR/Carveout.
  d-checks eigene Aussage (dichte, repo-weite Nummern, [MR-000](../conventions.md#mr-000--baseline-aussage))
  bleibt davon unberührt — die Vorlage generalisiert, sie zwingt nicht.
  Zweitens trägt `modul-05-planning-harness.md` einen **vierten
  Lifecycle-Zweig** — „Ein Slice, dessen Gegenstand ein anderer übernimmt"
  (ein Slice geht ohne eigene Lieferung nach `done/`, wenn sein Gegenstand
  entfällt oder aufgeht) —, gespiegelt in `slice.template.md` (neue Felder
  `Übernimmt:`/`Gegenstand:`). Ob und wie d-check das adoptiert, ist Sache
  des Folge-Slice.

  `AGENTS.template.md` §4 ändert die mit slice-222 übernommene Zeiger-Form
  **nicht erneut** — Abgrenzung 2 dieser Hebung hält damit ohne erzwungene
  Ausnahme, anders als beim Vorgänger.

  **Vier Spiegel-Klassen, drei davon gate-blind**
  ([`BEO-ALL/pin-bump-mirrors-ungated`](../../docs/plan/planning/observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md),
  weiterhin ohne formgültigen Ausgang). Gemessen am **echten** Vorzustand —
  direkt nach dem Materialisieren von `v6.9.0`, vor jeder Ersetzung, **beide**
  Bäume noch vorhanden: 87 Dateien nannten `v6.6.0` mit **181** Vorkommen.
  **28** davon sind der vendorte Baum selbst (26 Regelwerk- + 2
  Template-Dateien mit Selbstverweis) — Klasse *vendorter Baum*, durch
  Entfernen erledigt, keine Retargeting-Frage. **59** Dateien liegen
  außerhalb.

  Von diesen 59 blieben **acht eingefroren** — Lauf-Belege, die den alten
  Pin korrekt als ihren damaligen Stand nennen: die MR-071-Datei selbst
  (jetzt in `conventions/done/`; Vergangenheits-Aussage über die
  v6.6.0-Hebung, bleibt bei ihrer Auflösung unverändert stehen), ein
  `done/`-Slice
  ([slice-222](../../docs/plan/planning/done/slice-222-baseline-v660-bump.md)),
  zwei Review-Reports, drei Evidence-Dateien im Beobachtungs-Register und ein
  wörtliches Fremdzitat im CR.

  Die übrigen **51 lebenden** Dateien sind retargetet — mehrheitlich als
  einfacher Pfad-Verweis (`.harness/baseline/v6.6.0/…`, gate-gedeckt — ein im
  Index behaupteter Verweis ohne Ziel meldet `target-missing`, sobald der
  alte Baum fehlt); mehrere davon tragen **zusätzlich** eine der drei
  gate-blinden Klassen (kein Zählraster, sondern dieselbe Datei aus
  mehrfachem Anlass): vier Vorkommen sind Release-/Tree-URLs
  (`harness/conventions.md` zweifach, `AGENTS.md`, `harness/README.md`),
  vier sind `d-check:cite`-Direktiven (zwei je in `slice-221` und in diesem
  Slice, alle vier auf `modul-05-planning-harness.md` zeigend — neu
  geankert auf Zeilen 363–364 bzw. 369, weil diese Datei zum inhaltlichen
  Delta gehört), sieben sind bare Versionsnennungen (`AGENTS.md`,
  `.d-check.closure.yml`, `MR-049`, `MR-053`, `MR-021`,
  `observations/README.md`, `roadmap.md`) — diese letzte Klasse ist reines
  Urteil, kein Muster fängt sie.

  **Ein Fehler der mechanischen Ersetzung ist aufgetreten und sofort
  behoben, noch im selben Arbeitsschritt:** Ein pauschaler `sed` über
  `harness/conventions.md` traf auch die **Tabellenzeile** von MR-071 in
  §Aktive Adaptionen — dieselbe Über-Hebungs-Klasse wie bei `MR-067` in
  slice-222, diesmal am Index-Eintrag statt an der Datei selbst. Erkannt vor
  dem nächsten Schritt (kein Review nötig, um es zu finden), zurückgesetzt
  auf `v6.6.0`. **Grenze bestätigt:** [`MR-070`](../conventions.md#mr-070)s
  Frozen-Liste deckt Datei-Eigenschaften, keine Tabellenzeilen — wer über
  eine Datei sed't, die selbst einen Adaptions-Index trägt, muss dessen
  Zeilen einzeln gegen ihren Gegenstand lesen.

  **`ignore-refs` wächst diesmal NICHT** — anders als bei den beiden
  Vorgängern, und das ist gemessen, nicht angenommen: `make doc-check` nach
  dem Entfernen des `v6.6.0`-Baums meldet **0** Befunde ohne einen neuen
  Eintrag. Die acht eingefrorenen Fundstellen (oben) tragen den Pfad
  ausschließlich in Inline-Code, Prosa oder als bloße Zeichenkette, die kein
  Modul auflöst — anders als beim Vorgänger, wo zwei Lauf-Belege den
  entfernten Baum als Markdown-**Link** trugen.

  **Die `d-check:cite`-Spannen brauchten Neu-Ankern** — anders als beim
  Vorgänger: sieben Direktiven zeigten `citation-mismatch`, alle sieben
  wegen reiner Zeilenverschiebung in unverändertem Wortlaut, neu geankert
  (`modul-11-verification.md`, `modul-09-implementierung.md`,
  `grundlagen-harness-dateien.md`, `grundlagen-begriffe.md`, dreimal
  `modul-05-planning-harness.md`). Die beiden `d-check:cite`-Spannen in
  `slice-221`/diesem Slice (*Sub-Area-Wahl prüfen* · *Offene Beobachtungen
  sichten*) sind wortgleich geblieben, nur verschoben (268–269→363–364,
  274→369) — die Nummerierung wechselte von Bullet auf `1./2.`, was den
  zitierten Text selbst nicht berührt.

  **Ein Zitat-Delta, nach [`MR-039`](../conventions.md#mr-039) hier
  vermerkt, nicht am zitierenden Dokument:**

  | Quelle | Zitiert (Stand `v6.6.0`, in [`MR-056`](../conventions.md#mr-056)) | Seit `v6.9.0` |
  |---|---|---|
  | `regelwerk/modul-05-planning-harness.md` §Lifecycle als State Machine | *„DoD-Häkchen und Closure-Notiz sind die Bedingung dafür, dass die Datei überhaupt nach `done/` darf."* | *„…darf — mit einer Ausnahme für die Liefer-Häkchen ([§Ein Slice, dessen Gegenstand ein anderer übernimmt])."* — die Ausnahme spiegelt den neuen vierten Lifecycle-Zweig (oben) |

  Die `d-check:cite`-Direktive an dieser Stelle ist entfernt — der
  mechanisch geprüfte Wortlaut kann die neue Ausnahme nicht zugleich zitieren
  und unangetastet bleiben. Das Zitat in `MR-056` bleibt stehen, wie es zum
  Zeitpunkt seiner Einführung lautete.
- **Begründung:** Der Nachtlauf meldete den neuen Release
  (`make baseline-freshness`, Exit 3), und der gepinnte Tag war upstream
  inhaltlich unverändert — die Hebung ist damit eine reine Fortschreibung,
  kein Reparatur-Fall.
- **Auflösungs-Trigger:** die nächste Pin-Hebung.
