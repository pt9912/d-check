# Review: slice-109 — v5.6.0-Konventions-Nachzüge C-2…C-6 (Erst-Review)

- **Review-Art:** Code-/Harness-Doku-Review (unabhängiger Erst-Review, frischer Kontext)
- **Gegenstand:** slice-109 (Etappe C-2…C-6 der Baseline-Migration v5.6.0), Commit `3229c10`
- **Skill:** `reviewer.md` @ 1.5.0 (Stand `3229c10`; der neue HIGH-Anker wurde mitgeprüft)
- **Modell-ID:** claude-fable-5
- **Datum:** 2026-08-21
- **Eingangs-Kontext:** `docs/plan/planning/in-progress/slice-109-v560-konventions-nachzuege.md`
  (Plan + §1a-Ergebnis), Audit-Findings C-2…C-6 aus slice-107 §9 samt F-4/F-9 der
  Audit-Reviews, `AGENTS.md` §3, Baseline `.harness/baseline/v5.6.0/regelwerk/`
  (`grundlagen-harness-dateien.md` §Was ein Kommentar trägt,
  `grundlagen-source-precedence.md` §ID-Schema/§Vergabe, `modul-03-spec.md`,
  `modul-14-docker-harness.md`), Kurs-Repo `/Development/KI/ai-harness-course`
  (read-only, Tags v5.0.0/v5.6.0). Kein `DC-*`-Produktvertrag berührt; die
  Kommentar-Bereinigung wurde auf Verhaltens-Neutralität geprüft.

## Findings

### F-1

- **kategorie:** MEDIUM
- **quelle:** AGENTS §3.7 (Kommentar-Klassen) / slice-109 C-3 („per grep **alle**
  Fundstellen … ziehen oder streichen")
- **pfad:** `tools/semgrep.sh:64` · `tools/bench-fixture.sh:58` · `Makefile:125` ·
  `.github/workflows/ci.yml:63`
- **befund:** Die C-3-Räumung deckt nur `internal/*.go` (16 Fundstellen); vier
  Fundstellen derselben Klasse (nacktes Review-Token als Herkunfts-Feld:
  `(Review R1, HIGH-1)`, `(Review R1 zu slice-009/010)`, `(slice-055-R2-LOW)`,
  `Review R1 HIGH-1`) stehen in Skripten und Konfiguration — Flächen, die die im
  selben Commit installierte Hard Rule §3.7 ausdrücklich nennt („Code,
  Konfiguration oder Skript"). Die Commit-Aussage „**alle** 16
  Review-Token-Fundstellen in Code-Kommentaren bereinigt" ist damit in ihrem
  „alle" falsch; das Slice-Mandat kennt keine Verzeichnis-Grenze.
- **verifizierbar:** nein (kein Gate; grep `'\(Review |R[0-9]-F-|-R[0-9]-LOW'`
  über `tools/ Makefile .github/` reproduziert die vier Treffer)
- **klasse:** räumung-scope-schmaler-als-regel

### F-2

- **kategorie:** LOW
- **quelle:** AGENTS §3.7 / reviewer.md 1.5.0 HIGH-Anker
- **pfad:** `internal/**/*_test.go` (≈40 Fundstellen, u. a.
  `internal/adapter/driving/cli/cli_acceptance_test.go:2322,2346,2502`,
  `internal/hexagon/core/app/trace_cross_test.go:132–541`,
  `internal/adapter/driving/cli/cli_config_path_test.go:77–137`)
- **befund:** §3.7 gilt seinem Wortlaut nach für „Code" ohne Test-Ausnahme; in
  Testdateien verbleiben ≈40 nackte Review-Tokens als Herkunfts-Feld, weder
  geräumt noch als deklarierte Grenze der Rule notiert. Latente Wartungsfalle:
  jeder künftige Diff, der eine dieser Zeilen berührt, kollidiert mit dem neuen
  HIGH-Anker, ohne dass irgendwo steht, ob Tests drin oder draußen sind.
- **verifizierbar:** nein (grep; kein Sensor)
- **klasse:** scope-grenze-undeklariert

### F-3

- **kategorie:** LOW
- **quelle:** Baseline §Was ein Kommentar trägt (Verbots-Klasse
  „Ersetzungs-Trümmer") / Maintainability
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go:1720–1722`
- **befund:** Die Streichung von „(Review R1/B2)." ersetzte die Kommentarzeile
  durch eine Leerzeile statt sie zu entfernen; der Doc-Kommentar
  „decodeStrict dekodiert mit KnownFields; nil-raw bei leerem Dokument" ist
  jetzt durch eine Leerzeile vom Funktionskopf getrennt und damit kein
  Go-Doc-Kommentar mehr — exakt der in der Baseline gemessene Modus („einen
  ließ er als Trümmer stehen").
- **verifizierbar:** teilweise (gofmt/lint schlagen nicht an — `make gates`
  blieb grün; sichtbar nur im Datei-Blick)
- **klasse:** ersetzungs-trümmer

### F-4

- **kategorie:** LOW
- **quelle:** slice-109 C-4 („Neue Verweise nutzen die Kennungs-Form")
- **pfad:** `harness/conventions.md:88` (MR-000-Bullet, Link auf MR-027)
- **befund:** Der einzige im selben Commit neu geschriebene MR-Index-Verweis
  (MR-000 → MR-027) nutzt den Voll-Slug
  `#mr-027--struktur-ids-spec-arc--werden-nicht-vergeben` statt der soeben
  eingeführten Kennungs-Form `#mr-027` — der Commit befolgt die von ihm selbst
  deklarierte Verweis-Regel nicht.
- **verifizierbar:** nein (beide Anker existieren; Links-Gate grün in beiden Formen)
- **klasse:** neue-form-nicht-selbst-genutzt

### F-5

- **kategorie:** INFO
- **quelle:** Baseline `grundlagen-harness-dateien.md` §harness/README.md als
  Einstiegspunkt (Pflichtgliederung)
- **pfad:** `harness/README.md:17`
- **befund:** Das Baseline-Skelett führt `## Leseordnung` als **letzte** Sektion
  („Die sieben Sektionen darüber sind eine Referenzfläche"); d-check setzt sie an
  Position 2 und formuliert die Prosa gespiegelt um („die Sektionen darunter").
  Die Positions-Wahl ist inhaltlich vertretbar, aber nirgends als Auslegung
  notiert; kein Sensor prüft die Sektions-Reihenfolge (`structure` prüft in
  `harness/README.md` nur `doc-tables`).
- **verifizierbar:** nein
- **klasse:** form-abweichung-undeklariert

### F-6

- **kategorie:** INFO
- **quelle:** slice-109 §1a (C-6-Tabelle, Zeile „Image-Hash im Build-Output")
- **pfad:** `docs/plan/planning/in-progress/slice-109-v560-konventions-nachzuege.md:69`
- **befund:** Das ✓ deckt die Regel-Überschrift (Hash im Build-Output —
  `make fullbuild` echot ihn, `Makefile:169`), nicht den Regel-Mechanismus:
  modul-14 verlangt ein **persistiertes** einzeiliges Beleg-Artefakt
  (`--metadata-file` → `harness/image-hash.txt`); ein solches existiert im Repo
  nicht. Der Konsument des Artefakts (Replay-Manifest-Slot) ist per n.-a.-Zeile
  derselben Tabelle abwesend, die Verkürzung also folgerichtig — aber die
  Tabelle notiert sie nicht.
- **verifizierbar:** nein
- **klasse:** stichproben-antwort-glättet-mechanismus

### F-7

- **kategorie:** INFO
- **quelle:** slice-109 §1a (C-6-Tabelle, Zeile „Devcontainer/Compose-Kriterium")
- **pfad:** `docs/plan/planning/in-progress/slice-109-v560-konventions-nachzuege.md:73`
- **befund:** modul-14 formuliert „Faustregel: Compose ist *Pflicht*
  (CI-Vertrag)" ohne die Vorbedingung „Service-Verbund", auf die sich die
  n.-a.-Begründung stützt. Die Lesart „Faustregel ≠ Hard Rule; verbundlose
  Ein-Container-CLI liegt außerhalb des Kriteriums" ist plausibel und
  wahrscheinlich richtig — sie steht aber nur implizit in der Tabelle; ein
  strenger Zweit-Leser könnte die Zeile als undeklarierte Abweichung lesen.
- **verifizierbar:** nein
- **klasse:** n-a-lesart-implizit

## Negativbefunde (geprüft, ohne Befund)

- **MR-027-Zahlen, selbst nachgezählt:** `grep -c '^### DC-FA-'
  spec/lastenheft.md` = 44; `.a`-Sektionen in `spec/spezifikation.md` = 36;
  `.a`+`.b` = 37 mit genau einer `.b` (`DC-FA-ANCH-001.b`, deren Anforderung
  auch die `.a` trägt → 36 **distinkte** Anforderungen). Alle drei Zahlen der
  Begründung korrekt.
- **MR-027 Auflösungs-Trigger:** beobachtbar (Schreibmoment einer ADR: zwei
  ununterscheidbare Festlegungen im selben Abschnitt) und deckungsgleich mit
  der modul-03-eigenen Rationale für `SPEC-*` („zwei Defaults im selben
  Abschnitt werden ununterscheidbar") — der Trigger greift exakt dort, wo die
  verzichtete Kennung ihren einzigen Konsumenten bekäme.
- **MR-027 Ersetzt-Baseline-Links:** beide auflösbar
  (`grundlagen-source-precedence.md:175` „### ID-Schema als Klammer" →
  `#id-schema-als-klammer`; `modul-03-spec.md` existiert und schreibt `SPEC-*`
  real vor). Die zitierten Baseline-Stellen existieren wörtlich: „ersatzweise
  der Abschnitt" (`grundlagen-source-precedence.md:233`), „ein Anker je Zeile
  wäre Pflege ohne Gegenwert" (ebd.:252).
- **Abgrenzung MR-000/MR-027:** sauber. Die Vergabe-Aussage in MR-000 (dichte
  repo-weite Nummern, ein Schreiber, kein Bereichssegment, „Verzeichnis und
  offene Welle-Dateien") ist Baseline-konforme Konkretisierung — §Vergabe
  bindet die Segment-Frage selbst an „genau ein Mensch schreibt" — und gehört
  daher in MR-000; der Struktur-ID-Verzicht als echte Abweichung ist nach
  MR-027 ausgelagert und in MR-000 verlinkt. Die frühere stille Lüge (Blanket
  „keine Adaptionen fürs ID-Schema" ohne SPEC/ARC-Erwähnung) ist aufgelöst;
  die Rest-Spannung im MR-000-Bullet löst der Abweichungs-Link im selben Satz.
- **C-3 Substanz-Stichprobe:** `cli.go` reorderArgs behält die Zusage („sonst
  würde das Pfad-Argument als Wert verschluckt"); `configyaml.go` behält
  KnownFields-/nil-raw- und Typnamen-Filter-Aussage; `trace.go`/`markdown.go`
  behalten das Failure-Szenario (Klammern im Linkziel, Backtick-Wrapper);
  `planning.go`/`structure_tableorder.go`/`trace_cross.go` behalten
  Kopplung/Abgrenzung. Kein Substanzverlust. Ersatz-Formen regelkonform
  (`seit slice-057`, `seit slice-059`, `seit slice-073` bzw. ersatzlose
  Streichung reiner Herkunft).
- **C-3 Verhaltens-Neutralität:** `git show 3229c10 -- internal/` enthält
  ausschließlich Kommentarzeilen-Änderungen (einzige Nicht-Kommentar-Textänderung
  ist die Leerzeile aus F-3); kein Code-Diff.
- **C-4 Anker-Mechanik:** exakt 27 Kurz-Anker `mr-001`…`mr-027`, einer je
  Index-Zeile (11 aktive + 16 aufgelöste), keine Dublette, keine doppelte
  `<a id>` in `harness/conventions.md`; Voll-Slugs byte-identisch erhalten;
  Kurz- und Voll-Form kollidieren nie (Voll-Slugs tragen `--`-Suffixe). Der
  Griff ist Baseline-gedeckt (explizite `<a id>` genau bei den MR-Zeilen des
  Adaptions-Index, `grundlagen-source-precedence.md:247–250`).
- **C-5 Leseordnung:** vier geordnete Zeiger (Soll 3–5), selektiv statt
  vollständig („eine Leseordnung, die alles nennt, ist keine" — erfüllt); alle
  vier Ziele existieren (`AGENTS.md`, `harness/conventions.md`,
  `docs/plan/planning/in-progress/` + `roadmap.md`,
  `.harness/baseline/v5.6.0/regelwerk/README.md`).
- **C-6 Delta-Freiheit, im Kurs-Repo gegengeprüft:**
  `lab/regelwerk/modul-14-docker-harness.md` existiert an **beiden** Tags und
  `git diff v5.0.0..v5.6.0` ist leer (0 Zeilen; auch die
  `kurs/de/05-betrieb/`-Variante delta-frei) — kein Datei-fehlt-Artefakt. Die
  vendorte Kopie ist bis auf die pin-gebundene Quelle-Zeile identisch
  (MR-021-konform).
- **C-6 Verkörperungs-Behauptungen, am Artefakt geprüft:** `Dockerfile` pinnt
  alle drei FROM-Zeilen per `@sha256:` (golang, golangci-lint, distroless);
  deps-Stage kopiert `go.mod`/`go.su[m]` vor `COPY . .`; Multi-Stage
  deps→compile/lint/test/coverage/build→runtime; Runtime
  `distroless/static-debian12:nonroot`, `USER 65532:65532`, nur `/d-check`
  kopiert. `Makefile:168–169`: `fullbuild` schließt mit dem Image-Hash;
  `Makefile:99`: `versions` zeigt die Runtime-Image-ID. Der
  `AGENTS.md#31-dockermake-only`-Anker löst auf. Die Replay-n.-a.-Zeile ist
  konsistent mit der modul-12-Skopierung des Stufen-Audits (deterministischer
  Kern, `DC-QA-02`).
- **AGENTS §3.7:** nächste freie Nummer, keine Kollision mit §3.1–§3.6 (kein
  bestehender § regelt Kommentare); Trigger-Form „permanent" erfüllt die
  v5.5.0-Anforderung; die Feld-Formen (`DC-*`, `ADR-*`, `seit slice-<NNN>`/
  `seit welle-<NN>`) sind zulässige Repo-Konkretisierung der Baseline-Liste
  (`LH-*` → `DC-*`; Slice-IDs sind laut Baseline stabile Tokens). Die
  `Verantwortlich:`-Regel steht als §5-Halbsatz, going-forward deklariert
  (kein Retrofit); slice-109 trägt das Feld selbst.
- **reviewer.md 1.5.0 / Skill-Form des neuen HIGH-Eintrags:** Versions- und
  Datums-Bump korrekt (1.4.0→1.5.0, 2026-08-21); der Eintrag trägt
  Einführungs-Marke und Trigger („neuer HIGH-Eintrag seit 1.5.0,
  Auflösungs-Trigger: permanent") und einen pin-gebundenen, auflösbaren
  Baseline-Verweis (`../baseline/v5.6.0/regelwerk/grundlagen-harness-dateien.md`
  `#was-ein-kommentar-trägt--code-konfiguration-skripte` — Datei und Heading
  verifiziert, Slug korrekt). Genügt der Skill-Form.
- **Gate-Lauf:** eigener `make doc-check` auf `3229c10` (== HEAD, Tree clean):
  „404 Datei(en) geprüft, 0 Befund(e)" — identisch zur Commit-Behauptung;
  alle neuen Links/Anker (MR-027, §3.7, Leseordnung, Kurz-Anker) mechanisch
  aufgelöst (`links`/`anchors` aktiv; `.harness/baseline/**` nur als Quelle
  ausgenommen, als Ziel geprüft).

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 1 | F-1 |
| LOW | 3 | F-2, F-3, F-4 |
| INFO | 3 | F-5, F-6, F-7 |

**Finding-Klassen dieses Laufs:** räumung-scope-schmaler-als-regel ·
scope-grenze-undeklariert · ersetzungs-trümmer · neue-form-nicht-selbst-genutzt ·
form-abweichung-undeklariert · stichproben-antwort-glättet-mechanismus ·
n-a-lesart-implizit.

## Verdikt

**Blockierend (eine MEDIUM-Auflage), kein APPROVE.** Die Substanz aller fünf
Etappen ist sauber: MR-027 ist zahlengeprüft, trigger-beobachtbar und sauber
von MR-000 abgegrenzt; die C-4-Anker-Mechanik ist fehlerfrei; C-5 und C-6 sind
belegt (C-6 inklusive eigener Kurs-Repo- und Dockerfile-Gegenprobe); §3.7 und
der neue HIGH-Anker genügen der geforderten Form. Blockierend ist allein F-1:
die im selben Commit installierte Hard Rule nennt Skripte und Konfiguration
ausdrücklich, vier grep-bare Fundstellen genau der geräumten Klasse stehen dort
weiter, und die Commit-Aussage „alle … bereinigt" hält der Nachmessung nicht
stand — dieselbe Ehrlichkeits-Klasse, die C-3 räumen sollte. F-2 (Test-Grenze
deklarieren oder räumen) sollte im selben Zug entschieden werden; F-3/F-4 sind
Kleinst-Fixes, F-5–F-7 Notiz-Kandidaten ohne Handlungszwang.
