# Review-Report — slice-107 Stufen-Audit (Baseline v5.0.0 → v5.6.0), unabhängiger Erst-Review

- **Review-Art:** Plan-/Audit-Review — geprüft werden die **Aussagen** des
  Stufen-Audits (`docs/plan/planning/in-progress/slice-107-baseline-v560-delta-audit.md` §9)
  gegen die Quellen, nicht die Form.
- **Gegenstand:** slice-107 §9 samt Etappe-C-Schnitt
  (`docs/plan/planning/open/slice-108-roadmap-offene-wellen.md`,
  `docs/plan/planning/open/slice-109-v560-konventions-nachzuege.md`),
  welle-78 §4 und Roadmap-Drift-Log; Commit `a0e6ecc`.
- **Skill:** `.harness/skills/reviewer.md` @ 1.4.0.
- **Modell-ID:** `claude-fable-5`.
- **Datum:** 2026-08-21.
- **Eingangs-Kontext:** vendorter Baum `.harness/baseline/v5.6.0/regelwerk/`
  (netzlos, Deckungsprobe gegen den Kurs-Tag `v5.6.0` siehe Negativbefunde);
  Stufen-Diffs und Tag-Notizen v5.0.0→v5.1.0→v5.2.0→v5.3.0→v5.3.1→v5.4.0→
  v5.5.0→v5.6.0 im Kurs-Repo (`git diff vX..vY -- lab/regelwerk`, read-only);
  `harness/conventions.md` (MR-000/MR-013/MR-024), `AGENTS.md`,
  `harness/README.md`, `spec/spezifikation.md`, `.d-check.yml`,
  `.d-check.closure.yml`, `internal/hexagon/core/model/config.go`,
  `internal/hexagon/core/rules/planning.go`, `planning_waves.go`,
  `internal/adapter/driving/cli/print_mk.go`. Kein Gate-Lauf durch den
  Reviewer (Verifikation ist getrennte Rolle); keine Host-Toolchain benutzt.

---

## Findings

### F-1

- **kategorie:** MEDIUM
- **quelle:** slice-107 §4 DoD („je Regel eine Zeile", welle-74-Lehre); Tag-Notiz v5.4.0 (NEU 2)
- **pfad:** `docs/plan/planning/in-progress/slice-107-baseline-v560-delta-audit.md:132-139`
- **befund:** Die v5.3.1/v5.4.0-Tabelle deckt nur zwei der drei
  Regel-Ergänzungen, die die Tag-Notiz v5.4.0 selbst als „NEU 1–3" nennt.
  NEU 2 fehlt als Zeile: der erweiterte Freshness-Audit-Ausgang „Widerspruch"
  (Wahl behalten-mit-benanntem-Widerspruch **oder** übernehmen; „übernehmen
  wollen, aber noch nicht können" ist Carveout, keine `MR` —
  `modul-02-harness-bootstrap.md`, v5.4.0-Diff). Das ist ausgerechnet die
  Regel, nach der dieser Audit selbst seine Abweichungs-Fälle (C-1/MR-024,
  Struktur-ID-Frage) einordnen müsste.
- **verifizierbar:** nein (Lese-Abgleich `git diff v5.3.1..v5.4.0 -- lab/regelwerk/modul-02-harness-bootstrap.md` gegen die Tabelle)
- **klasse:** aufzählung-unvollständig

### F-2

- **kategorie:** MEDIUM
- **quelle:** Baseline `modul-02-harness-bootstrap.md` §Freshness-Audit, siebte Eigenschaft (Auswahlkriterium der Bestands-Stichprobe)
- **pfad:** `docs/plan/planning/in-progress/slice-107-baseline-v560-delta-audit.md:120`; Spiegel `docs/plan/planning/open/slice-109-v560-konventions-nachzuege.md:35-38`
- **befund:** Der als „delta-freier Kandidat" benannte Stichproben-Gegenstand
  `modul-07-carveouts.md` ist **nicht** delta-frei: die Datei änderte sich in
  Stufe v5.2.0 (+13/−2: neuer Absatz „Eine Diskrepanz-Klasse liegt außerhalb
  dieses Trichters" + `SL-CO-AUDIT`→`slice-CO-AUDIT`) und wurde damit von
  genau diesem Delta-Review berührt. Das Auswahlkriterium verlangt die
  Komplementärmenge zu `git diff <alt> <neu>`; die tatsächlich delta-freie
  Menge v5.0.0..v5.6.0 ist `grundlagen-klassifikation.md`,
  `modul-00-einfuehrung.md`, `modul-01-entwicklungszyklus.md`,
  `modul-14-docker-harness.md`, `modul-15-observability.md`,
  `modul-16-produktiver-betrieb.md`. slice-109 §1.5 (C-6) erbt den falschen
  Kandidaten.
- **verifizierbar:** nein (Nachrechnung: `git diff --stat v5.0.0..v5.6.0 -- lab/regelwerk` enthält `modul-07-carveouts.md`)
- **klasse:** stichprobe-kandidat-falsch

### F-3

- **kategorie:** MEDIUM
- **quelle:** [`MR-025`](../../harness/conventions.md#mr-025--semantik-änderung-die-spiegel-vor-dem-editieren-auflisten) (vom Audit in §7 selbst als Arbeits-Brille zugesagt: „jede ‚anzupassen‘-Zeile benennt ihre Spiegel gleich mit")
- **pfad:** `docs/plan/planning/in-progress/slice-107-baseline-v560-delta-audit.md:145`; `docs/plan/planning/open/slice-108-roadmap-offene-wellen.md:30-45`
- **befund:** Die C-1-Zeile benennt als Spiegel nur Roadmap, `MR-024` und die
  `planning`-Selbstkonfiguration. Fünf weitere Repo-eigene Träger der alten
  Form („§Aktuelle Welle" / „Keine aktive Welle") bleiben unbenannt und
  stehen auch nicht in slice-108 §2: `AGENTS.md:99-100` (§3.3
  Lifecycle-Move-Ausnahme, Wortlaut des Roadmap-Flips), `AGENTS.md:144`
  (§4-Zeile `make planning-check`), `Makefile:81` (Target-Kommentar),
  `harness/README.md:60` (Sensors-Zeile) und
  `harness/conventions/MR-013-lifecycle-move-buendelung.md:15,20`
  (kanonische Quelle der §3.3-Ausnahme). Nach C-1 rotten diese fünf still —
  exakt die Klasse, für die MR-025 geschaffen wurde. (Produkt-Doku-Nennungen
  des *Defaults* — `spec/spezifikation.md`, Handbuch, READMEs, `CHANGELOG.md`
  — sind korrekt keine Spiegel: der Produkt-Default bleibt unverändert.)
- **verifizierbar:** teilweise (nach C-1-Umsetzung meldet `make doc-check`
  die MR-013-/AGENTS-Divergenz **nicht** — Prosa-Drift ist kein Gate-Gegenstand;
  genau darum muss die Liste im Plan stehen)
- **klasse:** spiegel-liste-unvollständig

### F-4

- **kategorie:** MEDIUM
- **quelle:** Baseline `modul-03-spec.md` §Ziel-Form Spezifikation („Alles Übrige … trägt eine `SPEC-<NNN>`"); `modul-02` §Freshness-Audit (Ausgang „Widerspruch": deklarierte Abweichung = `MR`)
- **pfad:** `docs/plan/planning/in-progress/slice-107-baseline-v560-delta-audit.md:118`; Spiegel `docs/plan/planning/open/slice-109-v560-konventions-nachzuege.md:19-22`
- **befund:** Die Struktur-ID-Zeile stützt sich auf „der Rückfallweg ist
  ausdrücklich zulässig" — die zitierte Klausel deckt aber Referenzen auf
  Elemente, die keine Kennung *tragen* („wo ein Element keine Kennung
  trägt"), nicht den repo-weiten Verzicht auf die Vergabe. modul-03 (v5.2.0)
  schreibt `SPEC-*` für genau die Sektionstypen vor, die
  `spec/spezifikation.md` §2–§6 führt (Datenschemas, Defaults, Fehler-Codes,
  Metrik-Felder, externe Verträge), und `modul-02` Schritt 1 übernimmt
  `SPEC-*`/`ARC-*` „aus der Baseline". Der Voll-Verzicht ist damit eine
  Abweichung vom Baseline-Default; der Audit ordnet ihn als bloße
  Deklarations-Zeile in der MR-000-Aussage ein — derselben Aussage, die
  „keine inhaltlichen Adaptionen ggü. Baseline-Default"
  behauptet (`harness/conventions.md:80-82`). Ob Deklarations-Zeile genügt
  oder der Freshness-Audit-Ausgang „Widerspruch" (behalten-als-`MR` /
  übernehmen / Carveout) greift, stellt der Audit nicht als Frage.
- **verifizierbar:** nein (Textabgleich `spec/spezifikation.md` §2–§6 — null `SPEC-*`/`ARC-*` — gegen `modul-03-spec.md` und `grundlagen-source-precedence.md` §ID-Schema-Tabelle)
- **klasse:** konform-basis-überdehnt

### F-5

- **kategorie:** MEDIUM
- **quelle:** slice-107 §4 DoD („**kein** ‚pauschal konform‘ ohne je Regel eine Zeile")
- **pfad:** `docs/plan/planning/in-progress/slice-107-baseline-v560-delta-audit.md:141-151`
- **befund:** Die v5.5.0-Tabelle (12-Dateien-Diff) führt sieben Zeilen; zwei
  normative Blöcke haben keine: (a) die §Vergabe-Ergänzungen in
  `grundlagen-source-precedence.md` — Welle zählt repo-weit dicht (fällt aus
  dem Sub-Area-Schema), die „lokal ableitbar"-Grenze („wer die nächste Nummer
  zieht, liest Verzeichnis **und** offene Welle-Dateien" — von d-check
  faktisch gelebt: welle-78 §4 vergab slice-108/109 vor Datei-Existenz),
  `MR`-Hybrid-Kollisionsklasse, Abgrenzung Nummern-Ableitung vs.
  Beanspruchen; (b) die modul-10-Dissens-Regel (Abweichung zwischen
  Rolleninhabern ⇒ Skill schärfen, nie die mildere Lesart). Beide stehen nur
  unter der Pauschale „Mehr-Schreiber-Teile sind laut Baseline selbst
  ‚entworfen, nicht belegt‘" — (a) ist davon gar nicht gedeckt, denn die
  Vergabe-Grenze gilt auch dem Ein-Schreiber-Repo, sobald eine Welle §4-Slices
  vor ihren Dateien nennt.
- **verifizierbar:** nein (Lese-Abgleich `git diff v5.4.0..v5.5.0 -- lab/regelwerk/grundlagen-source-precedence.md lab/regelwerk/modul-10-review-harness.md` gegen die Tabelle)
- **klasse:** aufzählung-unvollständig

### F-6

- **kategorie:** LOW
- **quelle:** slice-107 §4 DoD; Baseline `modul-09-implementierung.md` §Was der Agent in den Kommentar schreibt
- **pfad:** `docs/plan/planning/in-progress/slice-107-baseline-v560-delta-audit.md:125-130`
- **befund:** Die v5.3.0-Tabelle führt nur die beiden
  `grundlagen-harness-dateien`-Regeln; der zweite Teil des v5.3.0-Deltas —
  `modul-09` §Was der Agent in den Kommentar schreibt, insbesondere die
  Emissions-Regel („Werkzeuge, die Repos erzeugen: emittierte Skripte,
  Makefiles und Konfigurationen tragen keine Slice- oder Befund-Nummer des
  Erzeuger-Repos") — hat keine Zeile. Für d-check als Producer
  (`--print-mk` emittiert `d-check.mk` in fremde Repos) ist das die
  einschlägigste Teilregel der Stufe. Gegenprobe dieses Reviews:
  `internal/adapter/driving/cli/print_mk.go` emittiert keine Slice-/
  Befund-Nummern, nur `DC-FA-*`-Produkt-IDs (Werkzeug-Vertrag, kein
  Erzeuger-Prozess-Kontext) — der Befund wäre also „konform", aber die Frage
  wurde nicht gestellt.
- **verifizierbar:** ja (`d-check --print-mk`-Ausgabe enthält kein `slice-`/`R?-F-`-Token)
- **klasse:** aufzählung-unvollständig

### F-7

- **kategorie:** LOW
- **quelle:** slice-107 §4 DoD
- **pfad:** `docs/plan/planning/in-progress/slice-107-baseline-v560-delta-audit.md:113-123`
- **befund:** Zwei kleinere v5.2.0-Kandidaten ohne Zeile: (a)
  `grundlagen-bootstrap.md` „Beide Modi regeln dieselbe Achse: Doc ↔ Code"
  (Regelwerks-Migration ist kein BF-Fall, zuständig ist der Freshness-Audit)
  — beschreibt genau das Werkzeug dieses Slices und wäre eine billige
  konform-Zeile; (b) `grundlagen-traceability.md`: „Architektur-ID" zählt
  nicht mehr als Traceability-Bezug (die `ARC-*` der Sicht adressiert nur
  innerhalb der Spec) — eine *geänderte*, nicht nur ergänzte Regel.
- **verifizierbar:** nein (Lese-Abgleich gegen `git diff v5.1.0..v5.2.0`)
- **klasse:** aufzählung-unvollständig

### F-8

- **kategorie:** LOW
- **quelle:** Maintainability (ehrlich zählen — dieselbe Klasse wie der 15/19-Nachzug in slice-095)
- **pfad:** `docs/plan/planning/in-progress/roadmap.md:98`
- **befund:** Der Drift-Log-Eintrag fasst „Sechs Konform-Belege ohne
  Handlung" zusammen; die §9-Tabellen führen zwölf konform-Zeilen (v5.1.0: 1,
  v5.2.0: 3, v5.3.0: 1, v5.4.0: 3, v5.5.0: 3, v5.6.0: 1). Die Zahl im
  Drift-Log ist keine Zählung der Tabelle.
- **verifizierbar:** nein (Nachzählen der Antwort-Spalte)
- **klasse:** zahl-ohne-zählung

### F-9

- **kategorie:** LOW
- **quelle:** Baseline `grundlagen-harness-dateien.md` §Was ein Kommentar trägt (Hard Rule: Herkunft als **ein** auflösbares Feld — `LH-*`, `ADR-*`, `· seit welle-<NN>`)
- **pfad:** `docs/plan/planning/in-progress/slice-107-baseline-v560-delta-audit.md:129`; Belege `internal/adapter/driving/cli/cli.go:331`, `cli.go:438`, `internal/hexagon/core/rules/planning.go:116`
- **befund:** Die Zeile behauptet „konform gelebt — exakt die im Repo
  etablierte Kommentar-Disziplin" ohne Stichprobe. Ein Spot-Check findet eine
  Herkunfts-Feld-Form außerhalb der drei zugelassenen: nackte
  Review-Finding-Tokens („(R2-F-2)", „(R2-F-4)", „— R1-F-5") — repo-weit
  mehrdeutig (dieselbe Kennung existiert in Dutzenden Review-Reports) und
  laut modul-10 zeigen sie auf Lauf-Belege, die „über Läufe hinweg nicht
  wieder gelesen werden". Die Varianten mit Slice-Kontext
  (`trace.go:389`, `markdown.go:447`: „slice-073 R1-F-1") lösen dagegen auf.
  Die Kommentar-Substanz (Zusage/Abgrenzung/Kopplung) ist in allen fünf
  Fällen konform — betroffen ist nur die Herkunfts-Feld-Form; der künftige
  C-3-HIGH-Anker wird genau diese Klasse treffen müssen.
- **verifizierbar:** nein (Urteils-Hälfte der Kommentar-Regel; kein Sensor gebaut)
- **klasse:** konform-ohne-stichprobe

### F-10

- **kategorie:** INFO
- **quelle:** Maintainability (dokumentationswürdige Annahme)
- **pfad:** `internal/hexagon/core/model/config.go:502-516`; `spec/spezifikation.md:2508-2509`
- **befund:** Der Produkt-Default `planning.heading` bleibt
  `## Aktuelle Welle` / `Keine aktive Welle` (bewusste C-1-Grenze „kein
  Produkt-Code"). Nach C-1 lebt d-check dauerhaft auf Nicht-Default-Config,
  und jedes Konsumenten-Repo, das die v5.6.0-Baseline-Form adoptiert, braucht
  denselben Override. Ob der Produkt-Default in einem künftigen Release der
  Baseline-Form folgt, ist eine offene Produkt-Frage, die nirgends notiert
  ist.
- **verifizierbar:** nein
- **klasse:** default-form-drift

### F-11

- **kategorie:** INFO
- **quelle:** slice-107 §9 Stufe v5.6.0
- **pfad:** `docs/plan/planning/in-progress/slice-107-baseline-v560-delta-audit.md:157`
- **befund:** „`main` ist nicht push-geschützt" ist eine Aussage über eine
  GitHub-Einstellung und netzlos nicht verifizierbar. Die gelebte
  Direkt-Commit-Praxis auf `main` ist belegt (durchgängige Commit-Historie
  ohne Merge-PRs, `HEAD` == `origin/main`); die Schutz-Aussage selbst bleibt
  unbelegt im Repo.
- **verifizierbar:** nein (außerhalb des Repos)
- **klasse:** externe-aussage-unbelegt

---

## Negativbefunde (geprüft, ohne Befund)

- **Vendorter Baum ≡ Kurs-Tag v5.6.0:** alle 26 Regelwerk-Dateien
  byte-identisch bis auf die bekannten Bundle-Rewrites (Quelle-Kommentar und
  zwei README-Links auf pin-gebundene GitHub-URLs). Der Audit hat gegen die
  richtige Quelle gelesen.
- **Stufe v5.1.0:** beide Zeilen decken das Delta (7 Dateien = §Vergabe +
  reine Heading-Ebenen-Absenkungen). Die konform-Begründung trägt: GitHub-
  Slugs sind ebenen-invariant, d-check-Anker in den vendorten Baum lösen
  weiter auf. MR-000 führt tatsächlich keine Vergabe-Deklaration
  (`harness/conventions.md:80-84`) — die C-2-Einstufung ist korrekt.
- **Stufe v5.2.0, `.a`-Behauptung:** bestätigt — alle 39 Sektionen unter
  `spec/spezifikation.md` §1 tragen `.a`/`.b`-Verfeinerungs-Suffixe, keine
  Ausnahme.
- **Stufe v5.2.0, Reconciliation-Register n. a.:** bestätigt — die Ziel-Form
  steht ausschließlich im BF-Walkthrough (`modul-02` §Das
  Reconciliation-Register, Schritt 8) und die Verzeichniskonvention
  annotiert „nur im Brownfield-Bootstrap"
  (`grundlagen-harness-dateien.md:15`). Für ein durchgängiges GF-Repo ohne
  Inventur-Funde ist das leere Pflicht-Artefakt korrekt abgelehnt.
- **Stufe v5.2.0, Golden-Set n. a.:** die Skopierung trägt — `modul-12`
  definiert „Modell" ausdrücklich als den *nicht-deterministischen Kern*;
  d-checks Kern ist per `DC-QA-02` (Anker existiert im Lastenheft)
  deterministisch, die Regeln haben keinen Gegenstand.
- **Stufe v5.2.0, `d-check.mk`-Producer-Zeile:** bestätigt (Generator
  `internal/adapter/driving/cli/print_mk.go`, eigene Nutzung über
  Makefile-Targets).
- **Stufe v5.3.0, Träger-Fehlstellen:** bestätigt — `AGENTS.md` führt keine
  Kommentar-Hard-Rule (einziger Treffer „kein PR-Kommentar", anderer
  Kontext), `reviewer.md` 1.4.0 keinen Fünf-Klassen-HIGH-Anker. C-3 korrekt.
- **Stufe v5.3.1/v5.4.0, C-4-Befund:** bestätigt — die Index-Zeilen tragen
  Voll-Slug-Anker (z. B. `harness/conventions.md:94-102`), keine kurzen
  Kennungs-Anker `mr-<NNN>`; die Baseline-Form (Kennungs-Anker zusätzlich,
  Alt-Slugs bleiben) ist korrekt wiedergegeben
  (`grundlagen-harness-dateien.md`, v5.4.0-Diff).
- **Stufe v5.4.0, Gate-Obermenge „gelebte Praxis":** Fundstellen existieren —
  Paritäts-Nachweis als DoD-Punkt der Skript-Ablösung
  (`docs/plan/planning/done/slice-063-targets-modul.md:125-129,202`,
  „Modul-Befundsatz ≡ Skript-Befundsatz, Gate nur mit belegter Parität
  stillgelegt"), arch-check-Ablösung
  (`docs/plan/planning/done/slice-058-arch-check-via-a-check.md`).
- **Stufe v5.5.0, C-1-Produkt-Behauptung:** bestätigt —
  `PlanningConfig.EffectiveHeading()`/`EffectiveMarker()`
  (`internal/hexagon/core/model/config.go:502-516`) mit den Config-Schlüsseln
  `planning.heading`/`planning.marker` (`spec/spezifikation.md:2508-2509`);
  kein Produkt-Code nötig. Die Wächter-Aussage stimmt: die erste Fähigkeit
  (`planning.go:22-49`) hält den Marker gegen `in-progress/` — exakt der
  Baseline-Wächter der neuen Form.
- **Stufe v5.5.0, Drift-Kopplung:** die Behauptung „bleibt im
  Ein-Wellen-Betrieb wahr" hält — W3 prüft `hasActive == (len(flach)==1)`
  (`planning_waves.go:91-111`), und die gelebte Kombination aus
  MR-013-Move-Bündelung und Eröffnung-mit-erstem-Slice (Drift-Log-Präzedenz
  welle-68/71/78) hält beide Seiten synchron; die Mehr-Wellen-Grenze ist in
  slice-108 §2.4 als benannte, nicht gelöste Grenze korrekt platziert.
- **Stufe v5.5.0, `team.md`/`lab/team-sim`:** bestätigt — null Vorkommen von
  `team.md` im gesamten vendorten Baum (Regelwerk und Templates), `team-sim`
  reist nicht mit. Die n.-a.-Begründung trägt.
- **Stufe v5.5.0, Leseordnung:** bestätigt — `harness/README.md` führt die
  sieben Referenz-Sektionen, keine §Leseordnung. C-5 korrekt.
- **Stufe v5.5.0, `Verantwortlich:` going-forward:** in slice-108/109 bereits
  gelebt (beide Köpfe tragen das Feld).
- **Stufe v5.6.0:** der eine Absatz (Wirkung statt Mittel, `MR`-Träger bei
  Push-Schutz) ist vollständig erfasst; Delta = 7 Zeilen in `modul-05`,
  nichts weiter. (Zur Schutz-Aussage siehe F-11.)
- **Wiedervorlage slice-090:** bestätigt — `modul-10` §Output-Schema zählt
  fünf Felder ohne `klasse`
  (`.harness/baseline/v5.6.0/regelwerk/modul-10-review-harness.md:71-72`,
  ebenso `reviewer.template.md`), das Report-Template sechs
  (`.harness/baseline/v5.6.0/templates/docs/reviews/review-report.template.md:41-46`).
  Die Drift besteht in v5.6.0 fort; „Upstream-Notiz, keine d-check-Handlung"
  ist die richtige Antwort.
- **Etappe-C-Schnitt, Deckung:** alle acht „anzupassen"-Zeilen sind auf
  C-1…C-6 gemappt und von slice-108 (C-1) bzw. slice-109 (C-2…C-6) gedeckt;
  welle-78 §4 und der Drift-Log-Eintrag sind nachgeführt. Keiner der beiden
  Slices trägt eine Behauptung, die der Audit nicht deckt (der
  AGENTS-§5-Halbsatz in slice-109 C-3 ist die Verkörperung der
  going-forward-Zeile) — vorbehaltlich F-2 (geerbter Stichproben-Kandidat)
  und F-4 (geerbte Deklarations-Lesart).

---

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 5 | F-1, F-2, F-3, F-4, F-5 |
| LOW | 4 | F-6, F-7, F-8, F-9 |
| INFO | 2 | F-10, F-11 |

**Finding-Klassen dieses Laufs:** aufzählung-unvollständig (4×) ·
stichprobe-kandidat-falsch · spiegel-liste-unvollständig ·
konform-basis-überdehnt · zahl-ohne-zählung · konform-ohne-stichprobe ·
default-form-drift · externe-aussage-unbelegt.

## Verdikt

**Blockierend (MEDIUM-Auflagen), kein APPROVE.** Der Audit ist in Substanz
und Methode tragfähig — die geprüften n.-a.-Begründungen halten, die
C-1-Produkt-Behauptung ist code-belegt, der Etappe-C-Schnitt deckt die
Findings. Blockierend sind fünf MEDIUMs derselben Wurzel: Die vom Audit
selbst angelegten Maßstäbe (§4 „je Regel eine Zeile", §7 „Spiegel gleich
mitbenennen", das Stichproben-Auswahlkriterium der Baseline) sind an vier
Stellen nicht eingehalten (F-1/F-5/F-3/F-2), und eine Konform-Basis ist
überdehnt (F-4). Alle fünf sind im Audit-Dokument bzw. in slice-108/109
behebbar, bevor Etappe C startet; danach steht einem APPROVE im Folgelauf
nichts entgegen.
