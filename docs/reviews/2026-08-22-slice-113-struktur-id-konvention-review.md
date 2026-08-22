# Review-Report: slice-113 — Struktur-ID-Konvention (MR-027 aufgelöst, `ids`-Muster SPEC/ARC, ADR-Index-Regel, AGENTS §5)

**Datum:** 2026-08-22 · **Review-Art:** Konventions-/Konfigurations-Review (Code-Review-Art nach Modul 10: geprüft gegen Slice-Plan, Baseline v5.7.0 `grundlagen-source-precedence.md` §ID-Schema als Klammer + §Vergabe und `modul-03-spec.md` §Zwei Kennungs-Arten, `DC-FA-ID-001` (Lastenheft + Spezifikation §DC-FA-ID-001.a), Hard Rules `AGENTS.md` §3.3/§3.7, `MR-000`/`MR-013`/`MR-015`/`MR-025`, Baseline §Konventionsspeicher) mit eigener Gegenprobe am gebauten Image
**Gegenstand:** Commit-Range `c7d9a52..5bdc6c3` — `b7b8a75` (Lifecycle-Move MR-027 nach `harness/conventions/done/`, Index-Zeile „Aktive" → „Aufgelöste Adaptionen", Slice-Pfad-Retarget) → `5bdc6c3` (MR-000-Vergabe-Aussage, `.d-check.yml` `ids.patterns` SPEC/ARC, ADR-Index-Kopf-Regel, AGENTS §5-Satz, Handbuch-Marker, Platzhalter-Kennungen in drei welle-80-Slices); **vor** der Closure, kein Release
**Skill:** `reviewer.md` @ 1.5.0 (Stand `f845e8b`) · **Modell-ID:** `claude-fable-5`
**Eingangs-Kontext:** Slice-Plan `docs/plan/planning/in-progress/slice-113-struktur-id-konvention.md` (§2 Schritt 5 Spiegel-Liste, §3 NICHT-Liste, §5 Risiken); Wellendokument `docs/plan/planning/welle-80-struktur-ids.md` (D1–D4, §6 Out-of-Scope); vendortes Regelwerk `.harness/baseline/v5.7.0/regelwerk/grundlagen-source-precedence.md` §ID-Schema als Klammer + §Vergabe, `modul-03-spec.md` §Zwei Kennungs-Arten, `grundlagen-harness-dateien.md` §Konventionsspeicher + §Was ein Kommentar trägt; `DC-FA-ID-001` (Linkpflicht-Mechanik, Ventile `exempt-paths`/`d-check:ignore`, Zieldatei frei); `MR-013` (MR-Move-Klausel), `MR-025` (Spiegel-Ableiter), `MR-015` (AGENTS routet); `internal/hexagon/core/rules/ids.go` als Ist-Wahrheit der Marker-Mechanik; Präzedenz-Moves MR-024 (`a46f09c`) und MR-026 (`f845e8b`). Nicht erhalten: DoD-Abhakung (Verifikations-Rolle); kein `make`-Target im Review (paralleler Gate-Lauf) — die Gegenprobe lief als Image-Lauf auf einer Baum-Kopie außerhalb des Repos.

## Findings

### F-1 · LOW

- **kategorie:** LOW
- **quelle:** Hard Rule `AGENTS.md` §3.7 (Zusage beschreibt, was da ist) / Maintainability
- **pfad:** `.d-check.yml:127-128`
- **befund:** Der Bruch-Satz der Zusage „Wortgrenzen sind Pflicht — ohne sie maskiert ARC-01 ein ARC-012 (Definitions-Maskierung + Fehl-Target)" beschreibt den Präzedenzfall des `diagrams`-Moduls (zweistelliges Muster, Definitions-Tabelle, `target`-Feld), nicht den Bruch des Musters, unter dem er steht: `ARC-\d{3}` kann `ARC-01` gar nicht treffen, und ohne `\b` bricht die Zusage an einer **vier**stelligen Kennung (`ARC-0123` liefert das Präfix-Token `ARC-012` als `id-unlinked`) — `ids` kennt weder Definitions-Maskierung noch Fehl-Target. Wer die Zeile ändert, liest einen Bruch, der hier nicht eintreten kann, und nicht den, der eintritt.
- **verifizierbar:** nein — Kommentar; die Mechanik selbst ist per Gegenprobe belegt (Negativ-Probe 6: `SPEC-0123`/`ARC-0123` ohne Befund mit `\b`).
- **klasse:** kommentar-bruchbeispiel-aus-fremdem-modul

### F-2 · LOW

- **kategorie:** LOW
- **quelle:** `MR-025` („Der Spiegel ist die Stelle, nicht die Datei") / Maintainability
- **pfad:** `.d-check.yml:256-257` (Kommentar des `commits`-Blocks)
- **befund:** „id-patterns = die drei ids-Muster (ADR/MR/DC) plus slice-NNN (eine ID-Definition neben ids.patterns …)" war vor dem Diff die Aussage „commits spiegelt alle ids-Muster"; seit `ids.patterns` fünf Muster trägt, fehlt an dieser Stelle die Abgrenzung, dass `SPEC-*`/`ARC-*` **bewusst** nicht in die Traceability gehören (MR-000, welle-80 §6) — in derselben Datei, die der Commit editiert. Wer das nächste `ids`-Muster anlegt, liest hier „ids-Muster ⇒ auch commits" und zieht Struktur-IDs in `commits.id-patterns`, womit eine Commit-Botschaft mit nur `SPEC-014` `trace-check` bestünde.
- **verifizierbar:** nein — Kommentar; `make trace-check` prüft Botschaften, nicht die Kopplung der Muster-Listen.
- **klasse:** spiegel-stelle-in-bearbeiteter-datei-nicht-nachgezogen

### F-3 · LOW

- **kategorie:** LOW
- **quelle:** `DC-FA-ID-001` (Ventil `d-check:ignore` wirkt auf `codepaths` **und** `ids`) / Doku-Drift
- **pfad:** `docs/user/benutzerhandbuch.md:1085` (§5, Überblick „vier Ventil-Achsen"; **Bestand außerhalb der Range**, durch den Diff erstmals im eigenen Repo für `ids` belegt)
- **befund:** Der Ventil-Überblick sagt „`d-check:ignore` (eine Zeile, nur `codepaths`)", während Lastenheft, Spezifikation §DC-FA-ID-001.a (Bedingung 4), die Handbuch-FAQ (Zeile 1878) und `versions.go` den Marker für `ids` (und `versions`) kennen — und der Slice setzt genau diesen Marker für `ids` in dasselbe Handbuch (Zeile 788). Ein Nutzer, der den Überblick liest, hält das `ids`-Zeilen-Ventil für nicht existent und greift zu `exempt-paths` (Ganzdatei). Slice §5 Risiko 1 verlangt, dass das Ventil „in der Konvention steht, damit niemand es als Fehlalarm liest" — die Nutzer-Doku sagt an dieser Stelle das Gegenteil.
- **verifizierbar:** nein — Prosa; kein Gate vergleicht Handbuch-Aussagen mit dem Ventil-Verhalten.
- **klasse:** ventil-geltungsbereich-in-nutzer-doku-veraltet

### F-4 · INFO

- **kategorie:** INFO
- **quelle:** `DC-FA-ID-001` („Begründung in Klammern empfohlen"), ADR-0040 („Direktiven-Konvention `<!-- d-check:ignore (Grund) -->`")
- **pfad:** `docs/user/benutzerhandbuch.md:788`
- **befund:** Der neue Marker trägt keine Begründungs-Klammer; Lastenheft, ADRs und ältere Slices tragen sie („Beispiel-ID, fiktiv"), die jüngeren Marker in Roadmap und slice-107/109 nicht — die Empfehlung ist Bestand-gemischt gelebt. Ohne Klammer weiß der nächste Handbuch-Editor nicht, dass der Marker die fremde Kennung `GG-SPEC-042` vor dem `SPEC-\d{3}`-Muster schützt; entfernt er ihn, wird es allerdings sofort rot (`pre-commit`-doc-check), also selbstheilend.
- **verifizierbar:** ja — Marker entfernen ⇒ `id-unlinked` `SPEC-042` (Negativ-Probe 6, Lauf 3).
- **klasse:** marker-ohne-begruendung

### F-5 · INFO

- **kategorie:** INFO
- **quelle:** `MR-013` / Hard Rule `AGENTS.md` §3.3 (MR-Lifecycle-Klausel)
- **pfad:** `harness/conventions/MR-013-lifecycle-move-buendelung.md:28-36`, `AGENTS.md:123-130`; Commit `b7b8a75`
- **befund:** Die MR-Move-Klausel nennt nur „Link-Tiefen-Fixes der bewegten Datei selbst … alles Übrige bleibt Commit 2"; der Move-Commit trägt notwendig **auch** die Index-Zeile in `harness/conventions.md` und den Slice-Pfad-Verweis (sonst `target-missing`, `pre-commit` rot) — wie die Präzedenz-Moves MR-024 (`a46f09c`) und MR-026 (`f845e8b`). Drei Moves in Folge bündeln gekoppelte Pfad-Verweise, die die Klausel für MR-Moves nicht deklariert (die Slice-Klausel tut es unter (b)); der Regeltext hinkt der gelebten, gate-erzwungenen Praxis um genau diesen Satz hinterher.
- **verifizierbar:** nein — kein Gate prüft die Commit-Form.
- **klasse:** regeltext-hinter-gelebter-praxis

### F-6 · INFO

- **kategorie:** INFO
- **quelle:** `MR-000` / Maintainability
- **pfad:** `harness/conventions.md:83-84`
- **befund:** Die Schema-Aufzählung im selben Bullet („ID-Schema (`DC-FA-*`, `DC-QA-*`, `ADR-NNNN`, `CO-NNN`, `slice-NNN`, `MR-NNN`; Präfix `DC`)") führt `SPEC-<NNN>`/`ARC-<NNN>` nicht, obwohl der Bullet sie zwei Sätze später als vergeben erklärt (sie war schon vor dem Diff nicht erschöpfend: `welle-NN` fehlt ebenfalls). Die Aussage ist nicht falsch, nur die Aufzählung und der Satz im selben Absatz decken sich nicht.
- **verifizierbar:** nein — Prosa.
- **klasse:** aufzaehlung-im-selben-absatz-nicht-nachgezogen

### F-7 · INFO

- **kategorie:** INFO
- **quelle:** Maintainability (undokumentierte Haus-Praxis)
- **pfad:** `harness/conventions/done/MR-027-struktur-id-verzicht.md:3`
- **befund:** Der Eintrag trägt in `done/` weiter `Status: Accepted` und kein „Aufgelöst durch"-Feld — identisch zur Präzedenz MR-026 (`f845e8b`) und template-konform (Template: „der Zustand ist die Verzeichnis-Position, kein Status-Feld"; das Feld kennt es nicht). 16 der 18 `done/`-Einträge tragen das Feld dennoch (MR-024 bekam es im Auflagen-Commit); die Auflösung steht für MR-027 nur in der Index-Zeile. Wer die Datei über den Pfad-Link der Spiegel-Liste direkt öffnet, liest „d-check vergibt keine Struktur-IDs" ohne Hinweis im Körper.
- **verifizierbar:** nein — Formfrage.
- **klasse:** aufgeloest-durch-feld-praezedenz-gemischt

### F-8 · INFO

- **kategorie:** INFO
- **quelle:** Baseline `grundlagen-harness-dateien.md` §Konventionsspeicher („Von außen wird der Index adressiert, nicht die Eintrags-Datei")
- **pfad:** `docs/plan/planning/in-progress/slice-113-struktur-id-konvention.md:70`
- **befund:** Die Spiegel-Liste verlinkt die MR-027-Eintragsdatei per Pfad; genau dieser Link brach im Move-Moment und musste im Move-Commit retargetet werden — der von der Baseline benannte Bruch. Als Spiegel-Liste (Datei = zu editierende Stelle) ist der Datei-Link vertretbar, und nach dem Move ist er stabil (`done/`-Einträge wandern nicht mehr); die Beobachtung ist eine Designnotiz, kein Defekt.
- **verifizierbar:** ja — `make doc-check` hätte den reinen Move ohne Retarget mit `target-missing` gemeldet.
- **klasse:** pfad-link-auf-wandernde-datei

## Negativ-Proben (geprüft, ohne Befund)

1. **Move-Commit-Schnitt (MR-013, AGENTS §3.3, Leitfrage 1).** `git diff --name-status -M c7d9a52 b7b8a75` ⇒ genau drei Pfade: `R087 harness/conventions/MR-027-struktur-id-verzicht.md → done/`, `M harness/conventions.md`, `M docs/plan/planning/in-progress/slice-113-struktur-id-konvention.md`. In der bewegten Datei ausschließlich drei Link-Tiefen-Fixes (`../../` → `../../../` für zwei Baseline-Links, `../conventions.md` → `../../conventions.md`); Index-Zeile von „Aktive" nach „Aufgelöste Adaptionen" mit beiden `<a id>`-Ankern; Slice-Pfad retargetet. Kein Übriges (MR-000-Satz, Muster, README, AGENTS folgen in `5bdc6c3`). Rename-Score 87 % > 50 %, Botschaft deklariert den Move als `git mv`. Ohne Befund.
2. **Anker-Erreichbarkeit (Slice §5 Risiko 2).** `grep -o 'id="mr-027[^"]*"' harness/conventions.md | sort | uniq -c` ⇒ `mr-027` 1×, `mr-027--struktur-ids-spec-arc--werden-nicht-vergeben` 1× (beide in der Aufgelöste-Tabelle). Eingefrorene Verweise per `git grep -n 'mr-027\|MR-027'`: done/slice-109, done/slice-110, welle-78-results (3×), welle-79-results, Roadmap (§Offene Wellen + Drift-Log) — alle auf `conventions.md#mr-027` (Kurz-Anker); kein Verweis auf den alten Dateipfad außerhalb der retargeteten Spiegel-Liste. Lauf 1 (unten) mit `links`+`anchors` aktiv: 0 Befunde. Ohne Befund.
3. **Status-Feld-Präzedenz.** Alle 18 Einträge in `harness/conventions/done/` tragen `Status: Accepted` (grep `^- \*\*Status`); MR-027 konsistent; Abweichung in der Feld-Praxis → F-7 (INFO). Ohne Befund.
4. **MR-000-Wortlaut gegen die Baseline (Leitfrage 2).** Satz für Satz gegen `grundlagen-source-precedence.md` §Vergabe/§ID-Schema gehalten: „fortlaufend je Datei" ✓ („Gezählt wird fortlaufend je Datei — die nächste freie Nummer ist die höchste vergebene plus eins"); „Lücken werden nicht nachbelegt" ✓ (wortgleich, Begründung Umlenkung älterer Verweise); „kein Bereichssegment" ✓ („Sie brauchen deshalb kein Bereichssegment"); „der Link trägt den Abschnitt, der Text die Kennung" ✓ (Überschrift der Baseline-Regel); „keine Anforderungen" ✓ („`SPEC-*` und `ARC-*` sind keine Anforderungs-IDs"); „nicht in Commit-Botschaften" ✓ („Die Klammer trägt die Anforderungs-ID, nicht jede Kennung … `SPEC-014` in einer Commit-Message sagt einem Reviewer nichts"); „Accepted vor welle-80 per §-Anker (immutabel)" ✓ („ersatzweise der Abschnitt" + Hard Rule §3.5). Zieldateien je Präfix stimmen mit der Straten-Tabelle (Technik = Spezifikation, Sicht = Architektur). Kein Bereichssegment, kein `DC-SPEC-*` (welle-80 §6). Ohne Befund.
5. **Alt-Wortlaut-Ableiter (MR-025, Leitfrage 2).** `git grep -n -i -E 'nicht vergeben|Struktur-ID|SPEC-\*|ARC-\*|SPEC-<NNN>|ARC-<NNN>' -- . ':!.harness/baseline' ':!docs/reviews' ':!*/done/*'` ⇒ Treffer nur in `.d-check.yml` (neuer Kommentar), `AGENTS.md` §5, ADR-README, Roadmap (welle-80-Zeiger „Umkehr von MR-027"), slice-113/114/115/116, welle-80, `harness/conventions.md` — **keine** Stelle trägt noch den Verzicht. `harness/README.md` Sensors-Zeile `doc-check` nennt die Linkpflicht ohne Muster-Aufzählung (kein Spiegel; Spiegel-Liste-Frage beantwortet: nein). Spezifikation §2 zeigt das generische `--suggest-config`-Gerüst, nicht die eigene Config (kein Spiegel; Produkt-Änderung per §3 ausgeschlossen). Datei-Liste deckungsgleich mit §2 Schritt 5; Stellen-Lücke in der bearbeiteten Datei → F-2, im selben Bullet → F-6. Ohne Befund im Datei-Zensus.
6. **`ids`-Muster — Gegenprobe am Image (Leitfrage 3, Slice §5 Risiko 3).** Baum ohne `.git` in ein Scratch-Verzeichnis außerhalb des Repos kopiert (469 Markdown-Dateien), Image `d-check:latest` (gebaut 2026-08-22 06:11, Produkt-Code in der Range unverändert), Ausgabe je Lauf in Datei, Exit explizit:

   ```bash
   docker run --rm --network none -v "<scratch>/tree":/repo:ro d-check:latest > lauf1.txt 2>&1; echo "exit=$?"
   # Lauf 1 (Baum unverändert):          exit=0  — "420 Datei(en) geprüft, 0 Befund(e)"
   # Lauf 2 (+ docs/probe.md, s. u.):     exit=1  — "421 Datei(en) geprüft, 4 Befund(e)"
   #   docs/probe.md:3  ARC-001  id-unlinked
   #   docs/probe.md:3  SPEC-001 id-unlinked
   #   docs/probe.md:3  SPEC-042 id-unlinked      (aus GG-SPEC-042 — Grenze bestätigt)
   #   docs/probe.md:5  SPEC-002 id-unlinked      (Inline-Code, link-policy always)
   # Lauf 3 (Probe entfernt, Handbuch Z. 788 ohne Marker): exit=1 — "420 Datei(en), 1 Befund(e)"
   #   docs/user/benutzerhandbuch.md:788  SPEC-042  id-unlinked
   ```

   Probe-Datei `docs/probe.md` trug: Zeile 3 „Siehe SPEC-001, ARC-001, GG-SPEC-042 und SPEC-0123."; Zeile 5 Inline-Code `SPEC-002` und `ARC-0123`; Zeile 7 verlinkte `SPEC-003`/`ARC-003` auf die Zieldateien; Zeile 9 `SPEC-004`/`ARC-004` mit Zeilen-Marker. **Erwartung exakt getroffen:** Befund für SPEC-001, ARC-001, GG-SPEC-042 (als Token `SPEC-042`) und den Inline-Code-Fall; **kein** Befund für die vierstelligen `SPEC-0123`/`ARC-0123` (beide `\b` greifen, RE2-ASCII-Wortgrenze, YAML einfach-quotiert), nicht für verlinkte, nicht für die Marker-Zeile. Lauf 3 belegt, dass der Handbuch-Marker **wirksam und notwendig** ist (Marker-Erkennung ist Substring auf der rohen Prosa-Zeile, `ids.go` `markerLines`; HTML-Kommentar am Zeilenende einer Prosa-Zeile, kein Fence, keine Tabellenzeile — ADR-0040-Pipe-Regel nicht berührt). Ohne Befund.
7. **Bestands-Zensus dreistelliger Token.** `git grep -nP '\b(SPEC|ARC)-[0-9]{3}\b' -- . ':!.harness/baseline'` ⇒ nur `docs/user/benutzerhandbuch.md:788` (markiert), `docs/reviews/*` (exempt), `.d-check.yml` (kein Markdown, nicht gescannt), Go-Tests (nicht gescannt). Zweistellige Vorkommen (`ARC-01…ARC-99`) in done/slice-045 und `diagrams_test.go` liegen außerhalb von `\d{3}`. Die Kommentar-Aussage „einziger Bestandsfall" stimmt; Slice §6 Rückführungs-Bedingung („am Bestand rot") nicht eingetreten. Ohne Befund.
8. **Ventil-Wahl (Leitfrage 3).** `exempt-paths` aufs Handbuch wäre ein Ganzdatei-Carve-out am Nutzer-Doku-Spiegel (künftige echte `SPEC-*`-Nennungen im Handbuch würden von der Pflicht genommen); Umformulierung der Beispiel-Kennung risse das `exclude-req`-Beispiel `^GG-SPEC-` (Fence Z. 751, `handbook_examples_test.go`, `cli_acceptance_test.go`) aus der Kohärenz; der Zeilen-Marker ist das spezifizierte Ventil genau für „bewusst illustrative Beispiel-IDs" (DC-FA-ID-001) und der kleinste Schnitt. Grenze im Kommentar korrekt beschrieben: `\b` trennt `-`|`S`, daher matcht `SPEC-042` in `GG-SPEC-042` (Lauf 2 zeigt den Token). Ohne Befund; Form → F-4.
9. **Hard Rule §3.7 am neuen Config-Kommentar (Leitfrage 4).** `.d-check.yml:126-131` trägt Zusage (Wortgrenzen Pflicht), Grenze (Bindestrich-Präfix, Ventil), Rang-Zeiger (MR-000; Baseline-Satz „Link trägt den Abschnitt"), Abgrenzung (Zieldatei frei) und Herkunft als **ein** Feld („seit welle-80"); keine Slice-Nummer, kein Befund-Marker, keine Deliberation über Verworfenes. Kein HIGH-Anker; Präzision des Bruch-Beispiels → F-1. Ohne Befund.
10. **ADR-Index-Regel ↔ AGENTS §5 ↔ MR-000 (Leitfrage 5, MR-015).** Drei Rollen, eine Aussage: ADR-README trägt die `Schärft:`-Form (Kennung wo vorhanden, sonst §-Anker; `Accepted` vor welle-80 bleibt) mit Zeiger auf MR-000; AGENTS §5 trägt die operative Kurzregel (entstehen nur beim Spec-Schreiben, nicht in Commit-Botschaften) mit Zeiger auf MR-000 und parallel zum bestehenden `MR-008`-Satz; MR-000 trägt die Vergabe-Regel. Ein gespiegeltes Detail („fortlaufend je Datei" in AGENTS) — Routing-konform, kein Widerspruch, keine Aussage, die in einer der drei Stellen anders lautet (Accepted-§-Anker in README und MR-000 identisch). Ohne Befund.
11. **Slice §5 Risiken (Leitfrage 6).** Risiko 1 (Linkpflicht `always` sofort repo-weit): Config-Kommentar benennt Wirkung, Grenze und Ventil; der §5-Text selbst wurde im Commit um „auch Beispiel-Kennungen in Inline-Code; das Ventil ist der Zeilen-Marker" geschärft; die drei welle-80-Slices tragen Platzhalter-Kennungen; Ausgang für die Closure-Notiz offen (Verifikation) — Nutzer-Doku-Seite → F-3. Risiko 2 (Anker) → Probe 2. Risiko 3 (Regex-Härte) → Probe 6 (`SPEC-0123` kein Befund). Ohne Befund.
12. **Commit-Botschaften gegen Stats (BEO-009-Probe).** `b7b8a75`: „drei Verweise eine Ebene tiefer" = exakt die drei Link-Fixes, Index-Zeile + Anker ✓, Retarget ✓; Drei-Klassen-Zensus plausibel (MR-027 trägt keine Tag-URL, Prosa-Nennungen des Verzichts außerhalb `done/` keine — Probe 5). `5bdc6c3`: jede genannte Datei im Stat (8 Pfade, keine Beifänge); die behauptete Negativ-Gegenprobe (SPEC-001/ARC-001 ⇒ Befund, SPEC-0123 ⇒ keiner) durch Lauf 2 unabhängig reproduziert. Ohne Befund.
13. **NICHT-Liste §3.** Kein Produkt-Code, keine `structure`-Regel, kein `commits.id-patterns`-Eintrag, keine Spec-Zeile im Diff (Stat: `.d-check.yml`, `AGENTS.md`, ADR-README, drei Slices, Handbuch, `conventions.md`, MR-027) — eingehalten. Ohne Befund.
14. **Referenz-Richtung / Provenance (MEDIUM-Anker).** Keine neuen `status-provenance`-Marker, keine Abwärts-Token in Spec-Straten (Spec unverändert); `matrix` in Lauf 1 ohne Befund. Ohne Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 (F-1, F-2, F-3) |
| INFO | 5 (F-4, F-5, F-6, F-7, F-8) |

## Verdikt: APPROVE

Der Slice liefert, was er verspricht, und belegbar: der Move ist ein MR-013-konformer `git mv` (R087) mit exakt den gekoppelten Verweisen und nichts Übrigem; beide Anker der MR-027-Zeile sind genau einmal vorhanden und jeder eingefrorene Verweis löst auf; die MR-000-Vergabe-Aussage deckt sich Satz für Satz mit Baseline §Vergabe/§ID-Schema, der Verzicht steht nirgends mehr außerhalb `done/`; die `ids`-Muster wirken im Image exakt wie deklariert (Lauf 1: 420/0 Exit 0; Lauf 2: die vier erwarteten `id-unlinked` inklusive der `GG-SPEC-042`-Grenze, keine Vierstelligen; Lauf 3: der Handbuch-Marker ist notwendig und wirksam); ADR-Index-Regel, AGENTS §5 und MR-000 sagen in drei Rollen dasselbe ohne Drift.

Nichts blockiert. **Nice-to-fix vor oder mit der Closure:** F-1 und F-2 sind zwei Kommentar-Zeilen in der ohnehin angefassten `.d-check.yml` (Bruch-Beispiel an das eigene Muster anpassen; `commits`-Abgrenzung benennen — das ist die „Spiegel ist die Stelle"-Lehre aus MR-025 in ihrer kleinsten Form). F-3 ist Handbuch-Bestand außerhalb der Range und gehört in den nächsten Release-Prep (Ventil-Überblick §5 um `ids`/`versions` korrigieren), nicht in diesen Slice. F-4 bis F-8 sind Notizen für Closure-Notiz bzw. Register (F-5: MR-013-Klausel um die gekoppelten Pfad-Verweise des MR-Moves ergänzen, wenn sie das nächste Mal angefasst wird).
