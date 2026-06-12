# Slice slice-012: Pilot-Migrationen (drei Tool-Familien)

**Status:** in-progress.

**Welle:** welle-04-distribution-und-migration (Abschluss).

**Bezug:** [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)
(Messmethode: Pilot-Migration in mindestens drei Repos);
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(repo-spezifisches Verhalten per Config); Meilenstein M3 (Teil 2:
≥ 1 Repo migriert — hier: drei).

**Autor:** pt9912. **Datum:** 2026-06-11.

---

## 1. Ziel

Drei Schwester-Repos — je ein Vertreter der Shell-Familie
(`verify-doc-refs.sh`), der Python-Familie (`check_refs.py`, inkl.
u-boot-Vollausbau) und der JS-Familie (`docs-check.js`) — nutzen das
veröffentlichte `d-check`-Image statt ihrer Tool-Kopie; die
Vergleichsläufe sind dokumentiert und belegen `DC-QA-04`.

## 2. Definition of Done

- [x] Pro Pilot-Repo: passende `.d-check.yml` geschrieben
  (Scan-Wurzeln, Module, ggf. ids-/matrix-Regeln äquivalent zur
  Alt-Tool-Abdeckung).
- [x] Pro Pilot-Repo ein dokumentierter **Vergleichslauf** Alt-Tool
  vs. `d-check` auf demselben Repo-Stand:
  `d-check` meldet mindestens dieselben echten Befunde und erzeugt
  keine False-Positives, die eine bislang grüne CI brechen
  (Akzeptanz-Kern von
  [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools));
  Differenzen werden triagiert: echter Mehr-Befund (Fix im
  Ziel-Repo oder dortige Ignore-Regel) vs. False-Positive
  (Spec-Fortschreibung/Fix hier, vor Abschluss des Slices).
- [x] Pro Pilot-Repo: CI-/Make-Schritt auf
  `ghcr.io/pt9912/d-check@sha256:…` (Digest-Pin gemäß
  Release-Doku aus slice-011) umgestellt; die Alt-Tool-Kopie ist im
  Ziel-Repo gelöscht, auf einen Rest-Sensor für repo-spezifische
  Prüfungen geschrumpft (Fall Kurs-Repo: Modul-Nummern-Checks) oder
  als deprecated markiert — Entscheidung liegt beim Ziel-Repo, wird
  dokumentiert.
- [x] Vergleichstabellen (Repo, Familie, Alt-Befunde, d-check-Befunde,
  Differenzen + Triage) in der Closure-Notiz — der zweite und dritte
  Datenpunkt nach dem Eigenlauf aus
  [`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding).
- [x] `make gates` grün (dieses Repo);
  [`CHANGELOG.md`](../../../../CHANGELOG.md); Closure-Notiz mit
  Steering-Loop-Lerneintrag — damit ist zugleich der
  welle-04-Closure-Trigger erfüllt (M3 erreicht).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `<pilot-repo>/.d-check.yml` (3×, extern) | neu | Familien-äquivalente Konfiguration |
| `<pilot-repo>`-CI/Makefile (3×, extern) | update | Digest-gepinnter `d-check`-Step, Alt-Tool-Ablösung |
| Closure-Notiz dieses Slices | neu | Vergleichsläufe + Triage als `DC-QA-04`-Beleg |
| [`CHANGELOG.md`](../../../../CHANGELOG.md) | update | Migrations-Stand |

## 4. Trigger

slice-011 done (erfüllt 2026-06-11) **und** slice-013 done — die
Pilot-Repos konsumieren das **veröffentlichte** Image per Digest-Pin,
und die JS-Familie braucht das Modul `codepaths`
([`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in),
Change Request 0.3.0): Der `DC-QA-04`-Vergleichslauf gegen das
erweiterte `docs-check.js` (2026-06-11) zeigte die
Inline-Code-Pfad-Prüfung als Lücke; die drei `docs-check:ignore`-
Marker des Kurs-Repos werden bei der Migration zu `d-check:ignore`.
Der kurs-spezifische Modul-Nummern-Sensor bleibt als Rest-Sensor im
Kurs-Repo (semantische Prüfung, out-of-scope laut Lastenheft).

## 5. Closure-Trigger

DoD vollständig (drei dokumentierte Vergleichsläufe, drei umgestellte
CI-Steps) + Closure-Notiz — erfüllt zugleich den
welle-04-Closure-Trigger und M3.

## 6. Risiken und offene Punkte

- Die Hauptarbeit liegt **außerhalb dieses Repos**; hier landen nur
  die Belege. Die Ziel-Repos haben eigene Gates/Prozesse — Änderungen
  dort folgen deren Konventionen.
- u-boot-Vollausbau (`check_refs.py`) ist der anspruchsvollste
  Vergleich (größter Funktionsumfang der Alt-Familie); Differenzen
  dort können Spec-Fortschreibungen hier auslösen — bewusst als
  letzter der drei Vergleiche einplanen.
- `DC-QA-04` verlangt „keine False-Positives, die eine grüne CI
  brechen" — bei legitimen Mehr-Befunden ist die Triage-Doku
  entscheidend (Mehr-Befund ≠ False-Positive).

## 7. Closure-Notiz (nach `done/`)

**Abgeschlossen:** 2026-06-12. Release v0.2.0 veröffentlicht
(`ghcr.io/pt9912/d-check@sha256:f2e0ac7b…d311846`) — alle sechs
Regelmodule im Image; die Pilot-CIs pinnen diesen Digest.

### Vergleichsläufe (Datenpunkte 2–4 nach [`MR-007`](../../../../harness/conventions.md#mr-007--auflösung-von-mr-003-doc-check-als-dogfooding), Beleg für [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools))

| Repo | Familie | Alt-Befunde | d-check (Erstlauf) | Differenzen + Triage |
|---|---|---|---|---|
| d-migrate | Shell (`verify-doc-refs.sh`, 98 Z.) | 0 | 11 | 8 absolute Host-Pfade (nur auf der Build-Maschine auflösbar, auf GitHub kaputt) + 3 veraltete Anker — alle **echte Mehr-Befunde** (Alt-Tool prüfte weder die GitHub-Semantik absoluter Ziele noch Anker oder Bilder); Fix im Ziel-Repo (61cfae08). Keine False-Positives. |
| ai-harness-course | JS (`docs-check.js`, 550 Z.) | 0 (116 Dateien) | 5 (116 Dateien) | Exakt die fünf `docs-check:ignore`-Markerzeilen (Slice-Plan nannte drei — seither gewachsen) → zu `d-check:ignore` konvertiert; identischer Scan-Umfang, `codepaths`-Parität inkl. Anker-Prüfung. Keine False-Positives. |
| u-boot | Python-Vollausbau (`check_refs.py`, 609 Z.) | 0 | 35 → 106 → 0 | Dreiteilig: (a) **5 False-Positives in d-check** — legale mehrzeilige CommonMark-Code-Spans invertierten zeilenbasiert die Backtick-Parität; Spec-Fortschreibung §`DC-FA-LINK-001.a` Schritt 2 + Fix hier vor Abschluss (546674e, in v0.2.0). (b) 1 Config-Artefakt (fehlende `\b`-Wortgrenzen im übersetzten Muster) — Config-Fix. (c) Nach Parser-Fix **~100 echte Mehr-Befunde**: neun ungeschlossene Titel-Code-Spans im CHANGELOG kippten die Parität ganzer Einträge, dazu Link-im-Code-Span-Nesting in Planning-JSON-Pins — GitHub-Rendering nachweislich kaputt, vom zeilenbasierten Alt-Parser durch Text-Verschmelzung übersehen; Fix im Ziel-Repo (470be86). **Finaler Gegenlauf: d-check 0 Befunde, Alt-Tool 8 False-Positives** auf den korrekten Code-Spans — der Migrations-Nutzen, am Alt-Tool selbst gemessen. |

### Umstellungen (CI-/Make-Schritt + Alt-Tool-Entscheidung)

- **d-migrate:** `make docs-check` → Digest-Pin v0.1.0 (`links`+`anchors` genügen der Familie; Pin-Hebung auf v0.2.0 ist dortige Routine), Alt-Skript **gelöscht** (ae2d13e9).
- **ai-harness-course:** `make docs-check` zweistufig — d-check (Digest v0.2.0, `codepaths` aktiv) + **Rest-Sensor** für die semantischen Modul-Nummern-Checks A/B (`docs-check.js` von 550 auf ~330 Zeilen geschrumpft) (d16a92f).
- **u-boot:** `make docs-check` → Digest-Pin v0.2.0; `check_refs.py` **deprecated** (kein Gate mehr — 8 False-Positives auf dem korrekten Stand); drei u-boot-spezifische Lints (verschachtelte Link-Artefakte, LH-Shorthand-Suffixe, Reference-Definition-Targets) als Rest-Sensor-Extraktion dort notiert (841ff7d).

### Lerneintrag (Steering Loop)

Die zeilenbasierte Inline-Code-Erkennung war als normative Grenze
dokumentiert, aber gegen reale Korpora falsch kalibriert — und der
Fehler wirkte **beidseitig**: d-check erzeugte 5 False-Positives auf
legalen mehrzeiligen Spans, während derselbe Parserfehler im Alt-Tool
~100 echte Befunde *versteckte* (Text-Verschmelzung zerstörte die
`\b`-Wortgrenzen). Geschärfte Regel: Parser-Vereinfachungen gegen die
Rendering-Wahrheit (CommonMark/GitHub) sind keine konservativen
Näherungen, sie verschlucken und erfinden Befunde zugleich;
Vergleichsläufe gegen fremde Korpora sind der wirksamste Sensor dafür
— die [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Messmethode
hat damit ihren zweiten Spec-Fortschritt ausgelöst (nach
[`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
aus dem JS-Vergleich).

## 8. Sub-Area-Modus-Begründung

Dieses Repo: GF (Doku-/Beleg-Arbeit). Pilot-Repos: außerhalb des
Geltungsbereichs dieses Harness (deren Modus gilt dort).
