# Review R2 — slice-264: Das Closure-Profil prüft die Slices in den Unterverzeichnissen von `done/`

**Review-Art:** Code (Diff gegen Plan, ADRs und Hard Rules — nicht gegen die DoD)
**Gegenstand:** slice-264, Plan-Änderung `1c6aa706` und Korrektur-Commit `f851c64f`
gegen die Findings aus Review R1 (`9da943ac`)
**Skill:** `.harness/skills/reviewer.md` @ 1.20.0
**Modell-ID:** claude-opus-5-5
**Datum:** 2026-10-09
**Eingangs-Kontext:** Slice-Plan slice-264 samt Plan-Änderung; Review R1 zu slice-264;
[ADR-0081](../plan/adr/0081-reviews-modul.md),
[ADR-0082](../plan/adr/0082-uebergangswaechter-reviews-observations.md);
[`MR-049`](../../harness/conventions.md#mr-049),
[`MR-025`](../../harness/conventions.md#mr-025);
[`DC-FA-RVW-001`](../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in),
[`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in);
`AGENTS.md` §3.5, §3.6, §3.7, §5, §6; Baseline `v6.17.0` ·
`regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register.

## Messungen (Kommando und Ergebnis)

Im Wegwerf-Klon des Stands `f851c64f`, Image aus `make build` desselben Stands.

- `make verify-closure-notes` → `978 Datei(en) geprüft, 0 Befund(e)`;
  `make review-coverage` → `1114 Datei(en) geprüft, 0 Befund(e)`.
- **Zwei Läufe, eine Frage (Wiederholung der R1-Probe):** die Reports zu slice-263
  entfernt — beide Läufe melden `review-missing` auf
  `done/wellenlos/slice-263-…md` und enden rot. R1 F-1 ist behoben.
- **Verengtes Stub-Muster:** `grep -rl '^> \*\*ARCHIVIERT\*\* — Volltext:' docs/plan/planning/done`
  → 254 Dateien, je ein Treffer; alle Stubs (Slice- und Wellen-Plan-Stubs). Der
  Erzeuger `tools/archive-wave/stub.go` schreibt genau diese Form.
- **`**` in `reviews.exempt-paths`:** Probe-Konfiguration mit `recursive`,
  `skip-pattern`, den Reports zu slice-263 entfernt und
  `exempt-paths: ["docs/plan/planning/done/**/slice-263-*"]` → 0 Befunde; ebenso
  mit dem ausgeschriebenen Pfad `…/done/wellenlos/slice-263-*`. Unter blankem
  `path.Match` steht `**` für **genau ein** Segment — die Form trifft also ein
  Unterverzeichnis eine Ebene tief, aber keinen Slice direkt unter `done/`.
- **Volltexte unter `done/wellenlos/` (zu R1 F-7):**
  `grep -L '^> \*\*ARCHIVIERT\*\* — Volltext:' docs/plan/planning/done/wellenlos/slice-*.md | wc -l`
  → **24**; `grep -L ARCHIVIERT docs/plan/planning/done/wellenlos/*.md | wc -l` → 23.
  Die Differenz ist slice-263: sein Volltext nennt das Wort `ARCHIVIERT` im
  Fließtext und fällt bei der Wort-Zählung heraus, ist aber kein Stub.

## Stand der R1-Findings

| R1 | Stand | Beleg |
|---|---|---|
| F-1 MEDIUM | behoben | Probe oben; `.d-check.yml` trägt `recursive`/`skip-pattern`, beide Kommentare nennen die Kopplung, `review-coverage.md` Grenze 1 nachgezogen |
| F-2 MEDIUM | behoben | slice-242 R2 trägt *weiter offen* mit Register-Eintrag `BEO-ALL/aufnahme-kriterium-ohne-bestands-beleg`; die Notiz „aus §7 nachgetragen" stimmt jetzt |
| F-3 MEDIUM | behoben | slice-240: drei Risiken *eingetreten — im selben Slice behoben/revertiert*, wie im Bestand (slice-231, slice-236); das erste *entfallen* mit Beleg aus dem R1-Report zu slice-240 (Negativbefund 3: „kein Über-Swap auf frozen Aussagen"), am Report nachgelesen |
| F-4 LOW | behoben, mit neuem Befund (F-11 unten) | Kommentar im Closure-Profil nennt KOPPLUNG und GRENZE |
| F-5 LOW | behoben | `SPEC-095` nennt die beiden stillen Ausfälle mit Verweis auf C2 |
| F-6 LOW | behoben | Muster auf `> **ARCHIVIERT** — Volltext:` verengt, alle sieben Stellen |
| F-7 INFO | bleibt | siehe Messung: 24 mit dem Muster des Profils, 23 mit der Wort-Zählung, die slice-263 verliert |
| F-8 INFO | zurückgestellt, akzeptiert | Handbuch-Pflege gehört in die Release-Prep (`AGENTS.md` §5 Regel 17) |

**Korrektur an R1 F-4:** Dort stand, die `**/`-Form „trifft dort gar nichts". Das
ist falsch — gemessen trifft sie einen Slice eine Ebene unter `done/`; sie verfehlt
nur Slices direkt in `done/` und solche zwei Ebenen tief. Der Rest von R1 F-4 bleibt.

## Findings

### F-9 — MEDIUM — Beide Profile scannen jetzt rekursiv, ADR-0081 entscheidet das Gegenteil

- **kategorie:** MEDIUM
- **quelle:** `AGENTS.md` §3.6 („Eine neue Fehlerklasse, ein neuer Scope oder ein
  Widerspruch zu einer bestehenden ADR braucht weiterhin eine eigene"); §3.5;
  [ADR-0081](../plan/adr/0081-reviews-modul.md) Entscheidung 4
- **pfad:** `docs/plan/adr/0081-reviews-modul.md` · „**Beide Verzeichnisse werden NICHT
  rekursiv gescannt.**"; `.d-check.yml` · `recursive: true` im `reviews`-Block
- **befund:** Die `Accepted`-ADR, die `make review-coverage` in der Sensor-Datei als
  Bindung führt, entscheidet den nicht-rekursiven Scan und begründet damit, dass ein
  Stub „natürlich" herausfällt. Seit `f851c64f` liest das Hauptprofil, seit `c55bae02`
  schon das Closure-Profil rekursiv, und der Stub fällt nur noch über `skip-pattern`
  heraus — ohne Folge-ADR und ohne `## Geschichte`-Eintrag an ADR-0081. Wer die
  Entscheidung an ihrer Quelle liest, liest die alte Kandidatenmenge und die alte
  Begründung, und hält etwa das `skip-pattern` für entbehrlich. R1 hat das am
  Closure-Profil übersehen.
- **verifizierbar:** nein — Urteil; `grep -n "NICHT rekursiv" docs/plan/adr/0081-reviews-modul.md`
  gegen `grep -n "recursive" .d-check.yml .d-check.closure.yml`.
- **klasse:** `config-widerspricht-accepted-adr`

### F-10 — LOW — Sensor-Datei nennt eine Bestands-Ausnahme, die der Config-Kommentar daneben „keine" nennt

- **kategorie:** LOW
- **quelle:** [`MR-025`](../../harness/conventions.md#mr-025) (Spiegel einer
  geänderten Semantik vorher auflisten)
- **pfad:** `harness/sensors/review-coverage.md` · „**Bestands-Ausnahme mit fester
  Dateiliste** (fünf Funde beim Scharfschalten"
- **befund:** Der im selben Commit umgeschriebene `reviews`-Kommentar in `.d-check.yml`
  sagt „BESTANDS-AUSNAHME: keine", und der Block trägt kein `exempt-paths`; Grenze 4
  derselben Sensor-Datei, deren Grenze 1 nachgezogen wurde, führt die Liste weiter.
  Die Drift ist älter als dieser Slice, steht jetzt aber zwei Absätze neben einer
  nachgezogenen Stelle.
- **verifizierbar:** ja — `grep -n "Bestands-Ausnahme" harness/sensors/review-coverage.md`.
- **klasse:** `spiegel-teilweise-nachgezogen`

### F-11 — LOW — Die neue GRENZE zu `reviews.exempt-paths` nennt die falsche Grenze

- **kategorie:** LOW
- **quelle:** Skill-Prüffrage 20; `AGENTS.md` §5 Regel 13
- **pfad:** `.d-check.closure.yml` · „eine Ausnahme fuer einen Slice unter einem
  Unterverzeichnis nennt dessen Pfad ausgeschrieben, nicht in der `**/`-Form"
- **befund:** Gemessen trifft `done/**/slice-…` einen Slice unter `done/wellenlos/`
  (Messung oben) — der Kommentar warnt also vor einer Form, die für genau diesen Fall
  funktioniert, und verschweigt den Fall, in dem sie still versagt: ein Slice direkt in
  `done/`, den die gleichlautende `structure`-Ausnahme derselben Datei trifft. Der Rat
  (Pfad ausschreiben) bleibt richtig, die Begründung nicht. R1 F-4 hatte denselben
  Fehler vorgegeben.
- **verifizierbar:** ja — die Probe-Konfiguration oben.
- **klasse:** `grenze-gegen-beschreibung-statt-gegenstand`

### F-12 — LOW — Der neue Register-Eintrag speichert den Zähler

- **kategorie:** LOW
- **quelle:** `docs/plan/planning/observations/README.md` · „der Zähler ist die Zahl
  dieser Dateien, kein gepflegtes Feld"; Baseline `v6.17.0` ·
  `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register · „Der Zähler wird
  abgeleitet, nicht geführt"
- **pfad:** `docs/plan/planning/observations/BEO-ALL/aufnahme-kriterium-ohne-bestands-beleg/state.md` ·
  „**Stand:** offen — 1×."
- **befund:** `state.md` schreibt die Zahl neben die Evidence-Liste; mit dem zweiten
  Beleg steht sie falsch, ohne dass ein Gate es meldet. Der Bestand führt dieselbe Form
  16-mal — sie ist nicht hier erfunden, wird aber hier fortgesetzt.
- **verifizierbar:** ja — `grep -rh "Stand:" docs/plan/planning/observations/BEO-*/*/state.md | grep -c "×"`.
- **klasse:** `register-zaehler-gespeichert`

### F-13 — INFO — Die Nachschärfung von `SPEC-095` hat keine Historie-Zeile

- **kategorie:** INFO
- **quelle:** Maintainability
- **pfad:** `spec/spezifikation.md` · „und ein Verzeichnis-Symlink fallen dabei still heraus"
- **befund:** Die Zeile vom selben Tag beschreibt die Fassung aus `c55bae02`; die in
  `f851c64f` ergänzten Ausfälle stehen in keiner. Die Spezifikation führt solche
  Nachzüge sonst als eigene Zeile „Nachzug nach Review an …".
- **verifizierbar:** nein.
- **klasse:** `spec-nachzug-ohne-historie`

## Negativbefunde

- **Plan-Änderung vor dem Code (`AGENTS.md` §6 Schritt 4):** `1c6aa706` liegt vor
  `f851c64f` und nennt beide neuen Gegenstände (Hauptprofil samt Sensor-Datei,
  Register-Eintrag) mit Befund-Bezug — ohne Befund.
- **Register-Eintrag (Form):** drei Dateien, Pfad aus dem deklarierten Kürzel `ALL`,
  Evidence-Dateiname ist ein geschlossener Vorgang (slice-242 in `done/wellenlos/`),
  Verweis aus slice-242 löst auf — ohne Befund (Zähler: F-12).
- **Ausgangs-Wortschatz (`MR-049`):** alle neuen Ausgänge in slice-240 und slice-242
  beginnen mit einem der drei Wörter; die Regel des Closure-Profils läuft grün.
- **Verlinkter R1-Report zu slice-240:** existiert, Pfad vom Ruheort aufgelöst, der
  zitierte Negativbefund trägt die Aussage.
- **Kommentare (`AGENTS.md` §3.7):** die neuen Kommentare in `.d-check.yml` und
  `.d-check.closure.yml` tragen Grenze, Kopplung und Zusage; die gestrichenen Absätze
  nahmen Herkunfts-Prosa und eine Slice-Nummer mit — ohne Befund (Inhalt der GRENZE: F-11).
- **Stub-Muster gegen den Erzeuger:** Marker-Form identisch mit `tools/archive-wave/stub.go`;
  weicht der Erzeuger ab, werden die Stubs Kandidaten und der Lauf rot — laut, nicht still.
- **`make gates`:** `reviews` ist in keiner `modules:`-Liste von `.d-check.yml`; die
  Änderung am Hauptprofil wirkt nur auf `make review-coverage`.
- **Weitere Stellen mit „nicht rekursiv" für `reviews`:** außer ADR-0081 (F-9) und dem
  Handbuch (R1 F-8, zurückgestellt) keine in lebenden Dateien.

## Kategorie-Summary

| HIGH | MEDIUM | LOW | INFO |
|---|---|---|---|
| 0 | 1 | 3 | 1 |

Wiederkehrende Klasse in diesem Slice: Spiegel einer Scope-Änderung unvollständig
nachgezogen (R1 F-1/F-5/F-8, hier F-9/F-10) — dritter Lauf in Folge nach slice-260 und
slice-263, ein Steering-Loop-Signal.

## Verdikt

**R1 ist abgearbeitet; F-9 vor der Closure klären.** Der Scope von `review-coverage`
widerspricht einer `Accepted`-ADR, und `AGENTS.md` §3.6 verlangt dafür eine eigene ADR,
keinen Config-Kommentar. F-10 bis F-12 sind klein und lassen sich
im selben Zug nachziehen.
