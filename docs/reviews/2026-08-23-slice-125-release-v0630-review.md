# Review-Report: slice-125 — Release-Prep `v0.63.0` (Doku-Currency, Versions-Register, CHANGELOG)

**Review-Art:** Design-/Doku-Review — geprüft wird die Nutzer-Doku und das
Release-Register gegen Lastenheft, Spezifikation, ADR, die Release-Prozedur
und **gegen den Code**; unabhängiger Reviewer ohne Anteil an der Arbeit.

**Gegenstand:** `d87c4d1..HEAD` — die zwei Commits des Slice:
`254e8c5` (Lifecycle-Move `open/` → `in-progress/` + Roadmap-Flip) und
`f4a7181` (Release-Prep `v0.63.0`).

**Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) @ 1.9.0 ·
**Modell-ID:** `claude-opus-5[1m]` · **Datum:** 2026-08-23

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan [slice-125](../plan/planning/in-progress/slice-125-release-v0630.md),
  Welle [welle-82](../plan/planning/welle-82-config-flaechen.md)
- [`spec/lastenheft.md`](../../spec/lastenheft.md) — `DC-FA-DIST-001`, `DC-QA-02`,
  `DC-FA-VER-001` (0.63.0), `DC-FA-STRUCT-001` (0.64.0), `DC-FA-DIAG-001` (0.65.0)
- [`spec/spezifikation.md`](../../spec/spezifikation.md) — §2-Schema zu
  `versions.patterns*`, `structure[].headings-match`/`headings-level`,
  `diagrams.exempt-paths`; §4 Grund-Code `SPEC-067`
- [ADR-0058](../plan/adr/0058-konfigurations-flaechen-additiv-weiten.md) (drei Entscheidungen)
- Release-Prozedur [`docs/user/releasing.md`](../user/releasing.md) §Release-Prep, §4-Checkliste
- Hard Rules [`AGENTS.md`](../../AGENTS.md) §3 (§3.3, §3.7),
  [`harness/conventions.md`](../../harness/conventions.md) (MR-013, MR-025)
- Beobachtungs-Register [`observations.md`](../plan/planning/observations.md)
  (BEO-002, BEO-006, BEO-007, BEO-008, BEO-009)
- Vorgänger-Review [2026-08-22 slice-124](2026-08-22-slice-124-diagrams-ventile-review.md)
  §„Für den Release-Prep benannt"
- Code: `internal/hexagon/core/rules/*.go`,
  `internal/adapter/driven/configyaml/configyaml.go`,
  `internal/adapter/driving/cli/config_template.go`

**Vom Reviewer selbst gefahren** (Exit-Codes in eine Datei umgeleitet und
direkt gelesen — BEO-007):

- `make doc-check` → **Exit 0**, „447 Datei(en) geprüft, 0 Befund(e)"
- `make gates` → **Exit 0** (acht Gates)
- `make coverage-gate` → **Exit 0**, „Coverage 94.80% erfüllt Schwelle 93%"
- vier eigene Image-Läufe gegen selbst gebaute Config- und Fixture-Varianten
  innerhalb der Repo-Wurzel (danach gelöscht; `git status` sauber)

**Verdikt: blockierend** — kein HIGH, **2 MEDIUM**, **4 LOW**, **2 INFO**.

---

## Findings

### F-1 — Das neue Begründungs-Kriterium für die Marker-Reichweite ist falsch

- `kategorie`: MEDIUM
- `quelle`: `DC-FA-CODE-001` / `DC-FA-ID-001`; BEO-009 (b) i. V. m. der
  Wellen-Lehre „Exklusivitäts-Aussage aus dem Anlass statt aus dem Bestand"
- `pfad`: [`README.de.md`](../../README.de.md):178–180,
  [`README.md`](../../README.md):176–178,
  [`docs/user/operations.md`](../user/operations.md):128–130
  (dazu [`CHANGELOG.md`](../../CHANGELOG.md):75–77)
- `befund`: Die Vierer-Aufzählung stimmt, ihre neu hinzugefügte **Begründung**
  nicht: „er stellt genau die vier Module still, **die eigene Muster
  konfigurieren und ihre Befunde an Zeilen hängen**". `matrix` konfiguriert mit
  `classes[].token` ein eigenes RE2 (`config.go` Feld `Token *regexp.Regexp`,
  kompiliert in `configyaml.go` `compileMatrixToken`) und meldet auf Zeilen
  (`matrix.go`, `Line: ref.Line` bzw. `Line: pl.no`); `structure` konfiguriert
  `section-pattern`, `forbid-pattern`, `require-pattern` und neuerdings
  `headings-match` und meldet `section-heading-mismatch` auf der Zeile der
  verletzenden Überschrift (`structure.go`, `structureFinding(..., sh.Line, ...)`).
  Beide erfüllen das Kriterium und honorieren den Marker **nicht** — in
  `operations.md` steht der Satz sogar unter der Überschrift „Kein Zeilen-Marker
  für `matrix`". Die CHANGELOG-Zeile zieht dasselbe Kriterium in der
  Gegenrichtung („`hostpaths`, `pins` und `spans` … sie konfigurieren keine
  eigenen Muster"), während `hostpaths.prefixes` die Präfixliste ersetzt, aus
  der `unixHostpathRE` das Muster baut. Versagen: Ein Leser wendet das genannte
  Kriterium auf `matrix`/`structure`/`hostpaths` an, erwartet dort den Marker,
  und der Befund bleibt stehen — die Doku hat ihm eine Regel gegeben, die sie
  selbst zwei Absätze weiter widerlegt.
- `verifizierbar`: nein — kein Gate deckt Prosa-Aussagen; nachprüfbar per
  `grep -L ignoreMarker` über `internal/hexagon/core/rules/*.go` und per Lauf
  `--enable matrix` gegen eine Zeile mit dem Marker.
- `klasse`: Exklusivitäts-Kriterium statt Aufzählung — aus dem Anlass, nicht aus dem Bestand

### F-2 — Der Spiegel im höherrangigen Stratum steht weiter auf der zu engen Aussage

- `kategorie`: MEDIUM
- `quelle`: MR-025 / BEO-002 (Spiegel-Pflicht); Source Precedence aus `AGENTS.md`
- `pfad`: [`spec/lastenheft.md`](../../spec/lastenheft.md):966,
  [`spec/spezifikation.md`](../../spec/spezifikation.md):844 und 1084–1086
- `befund`: Nach diesem Commit sagen Handbuch, `operations.md` und beide READMEs
  „vier Module"; das Lastenheft führt den Marker in der Achsen-Abgrenzung des
  geteilten Referenz-Ventils weiter als „(Zeile, **nur** `codepaths`)", die
  Spezifikation zweimal als „`d-check:ignore` (**nur** `codepaths`, Schritt 1)"
  und „`matrix` trägt — anders als `codepaths`/`ids` — keinen
  `d-check:ignore`-Marker (… Schritt 1 **gilt nur für jene Module**)". Die
  Lastenheft-Stelle ist schon vor dieser Welle falsch (`ids` honoriert den Marker
  seit Lastenheft 0.8.0). Versagen: Wer der Source Precedence folgt und im
  ranghöheren Dokument nachschlägt, bekommt genau die Antwort, die dieser Slice
  als Fehler ausweist — die Korrektur endet an der Dokument-Rolle statt an der
  Aussage.
- `verifizierbar`: nein — kein Gate; nachprüfbar per `grep` nach dem alten
  Wortlaut über `spec/`.
- `klasse`: Spiegel-Liste endet an der Dokument-Rolle (Nutzer-Doku korrigiert, Spec-Stratum nicht)

### F-3 — Die schließende Fence-Zeile ist im Handbuch nicht als Nicht-Ventil benannt

- `kategorie`: LOW
- `quelle`: `DC-FA-DIAG-001`, Spezifikation §`DC-FA-DIAG-001.a` Schritt 5
- `pfad`: [`docs/user/benutzerhandbuch.md`](../user/benutzerhandbuch.md):1109–1127
- `befund`: Der neue Absatz nennt zwei Ventil-Orte (Diagramm-Zeile,
  Öffnungszeile). Dass der Marker auf der **schließenden** Fence-Zeile
  wirkungslos ist, sagt die Spezifikation ausdrücklich, die Nutzer-Doku nicht —
  dabei ist die Schlusszeile die symmetrische Fehlplatzierung zur beworbenen
  Öffnungszeile. Eigener Lauf (`--enable diagrams`, vier Fälle in einer
  Fixture): Marker auf der Öffnungszeile mit und ohne `%%` unterdrückt, Marker
  **nur** auf der Schlusszeile meldet die Kennung weiter (`diagram-id-undefined`),
  Kontrollfall meldet ebenfalls. Versagen: Ein Nutzer platziert den Marker
  symmetrisch auf der Schlusszeile, hält den Block für ausgenommen und findet in
  der Doku keinen Hinweis, warum der Befund bleibt.
- `verifizierbar`: ja — `--enable diagrams` gegen eine Datei mit Marker
  ausschließlich auf der schließenden Fence-Zeile (Befund bleibt, Exit 1).
- `klasse`: benannte Nicht-Wirkung nur im Spec-Stratum, nicht in der Nutzer-Doku

### F-4 — Das neue Beispiel steht nicht in dem Fence-Schutz, den seine Form suggeriert

- `kategorie`: LOW
- `quelle`: Maintainability; `DC-FA-ID-001`
- `pfad`: [`docs/user/benutzerhandbuch.md`](../user/benutzerhandbuch.md):1118–1123
- `befund`: Das Beispiel setzt einen Tilden-Fence **innerhalb** eines
  `markdown`-Fence. Der geteilte Fence-Automat (`markdown.go`, `FenceToggle`)
  schaltet bei jeder `~~~`-Zeile um — die Tilden-Zeile **schließt** damit den
  äußeren Block, und die beiden Diagramm-Zeilen des Beispiels sind für alle
  Prosa-Module Fließtext. Eigener Lauf mit demselben Konstrukt und
  dreistelligen Kennungen: zwei `id-unlinked`-Befunde auf der Body-Zeile,
  während dieselben Kennungen in einem gewöhnlichen Fence stumm bleiben. Heute
  grün ist das Beispiel allein deshalb, weil seine Kennungen **zweistellig**
  sind — die Beispiel-Config zwei Abschnitte höher führt `regex: 'ARC-\d{2}'`.
  Die Datei selbst trägt keinen Hinweis auf diese Bedingung. Versagen: Der
  nächste Edit an dem Beispiel (dreistellige Kennung, ein Pfad in Inline-Code,
  eine `DC-`/`MR-`-Kennung) macht `make doc-check` rot an einer Stelle, die wie
  geschützter Code-Block aussieht.
- `verifizierbar`: ja — dasselbe Konstrukt mit dreistelligen Kennungen, Modul `ids`.
- `klasse`: Beispiel-Block ohne den Fence-Schutz, den seine Form suggeriert

### F-5 — „alle 24 bare-Tag-Pins" benennt die gate-gedeckte Menge mit dem Namen der gate-blinden

- `kategorie`: LOW
- `quelle`: BEO-009; `docs/user/releasing.md` §Release-Prep Punkt 4
- `pfad`: Commit `f4a7181` (Botschaft, Handbuch-Bullet);
  [`docs/user/benutzerhandbuch.md`](../user/benutzerhandbuch.md):75
- `befund`: Die Botschaft sagt „alle 24 bare-Tag-Pins auf v0.63.0". Gezählt
  trägt das Handbuch **24** `ghcr`-**präfixierte** Pins — genau die Klasse, die
  das aktive `versions`-Gate erzwingt — und **27** Vorkommen von `v0.63.0`
  insgesamt. Die Prozedur definiert „bare Tag" ausdrücklich als `:vX.Y.Z`
  **ohne** `ghcr`-Präfix („vom `versions`-Gate nicht erfasst, driftet still");
  davon trägt das Handbuch genau **einen** (Zeile 75). Die einzige Stelle, für
  die die Handnachführung überhaupt nötig war, verschwindet damit in einer Zahl,
  die sie nicht enthält. Versagen: Der nächste Release-Prep nimmt diese Botschaft
  als Vorlage, hebt die 24 `ghcr`-Pins per `grep` und lässt den einen bare Tag
  stehen — der Drift, den §4 benennt.
- `verifizierbar`: nein — Botschaft gegen Datei-Bestand; nachzählbar per
  `grep -c` auf `ghcr.io/pt9912/d-check:v0.63.0` gegen `grep -c 'v0\.63\.0'`.
- `klasse`: Botschafts-Zahl benennt eine andere Menge als die gemessene

### F-6 — Die Beispiel-Vorlage kombiniert die neuen Schlüssel zu einer still inerten Bedingung

- `kategorie`: LOW
- `quelle`: `DC-FA-STRUCT-001`, Spezifikation §2 (`structure[].headings-level`)
- `pfad`: [`docs/user/benutzerhandbuch.md`](../user/benutzerhandbuch.md):1749 und 1760–1761
- `befund`: Der Beispiel-Block führt `section-pattern: '^#{2,3} .*Closure-Notiz'`
  und darunter `headings-level: 3`. Trifft die Regel einen Abschnitt auf Ebene 3,
  kann in ihm keine Überschrift der Ebene 3 vorkommen — die Bedingung ist
  wirkungslos wahr. Das ist genau die Grenze, die der Fließtext desselben
  Abschnitts zwei Absätze später benennt („eine Ebene **flacher** als der
  Abschnitt kann in ihm gar nicht vorkommen"); an der Vorlage steht sie nicht.
  Eigener Lauf mit einkommentierten `headings-match`/`headings-level`: Exit 1,
  kein Exit 2, kein `structure`-Befund — der Leerlauf meldet sich nicht.
  Versagen: Ein Nutzer kopiert die Vorlage, kommentiert beide Zeilen ein, sieht
  Grün und hält eine Bedingung für aktiv, die für einen Teil ihrer Treffer nie
  zünden kann.
- `verifizierbar`: ja — Lauf mit der einkommentierten Vorlage gegen eine Datei
  mit `### …`-Abschnitt; kein `section-heading-mismatch`.
- `klasse`: Beispiel-Vorlage demonstriert eine still inerte Bedingungs-Kombination

### F-7 — Das `%%` steht im Beispiel an einer Stelle, an der es kein Kommentar ist

- `kategorie`: INFO
- `quelle`: `DC-FA-DIAG-001`; Maintainability
- `pfad`: [`docs/user/benutzerhandbuch.md`](../user/benutzerhandbuch.md):1112–1119
- `befund`: Der Absatz begründet die Token-Form damit, dass die **Diagramm-Sprache**
  den Marker vor dem Renderer versteckt („in Mermaid `%%`"), und zeigt ihn dann
  auf der **Öffnungszeile**. Dort ist er Teil des CommonMark-Infostrings, nicht
  Diagramm-Text; `%%` ist an dieser Stelle kein Mermaid-Kommentar, sondern
  Infostring-Text. Das Modul liest ohnehin nur das Token: `fenceInfo` schneidet
  den Infostring am ersten Leerzeichen ab (Sprache bleibt `mermaid`), und
  `CheckDiagrams` prüft `strings.Contains(dl.fenceOpen, ignoreMarker)`. Eigener
  Lauf: Öffnungszeile **ohne** `%%` unterdrückt identisch. Unausgesprochen bleibt
  damit die Annahme, dass ein Renderer einen Infostring mit Zusatztext weiterhin
  als Diagramm liest.
- `verifizierbar`: ja — Fixture mit Marker ohne `%%` auf der Öffnungszeile,
  gleiche Unterdrückung.
- `klasse`: Begründung gilt nicht am demonstrierten Ort

### F-8 — Die Proben-Aussage zu den Fenced-Beispielen nennt weniger, als sie behauptet

- `kategorie`: INFO
- `quelle`: BEO-009 (b)
- `pfad`: Commit `f4a7181` (Botschaft, Absatz „Die fenced YAML-Beispiele …")
- `befund`: Die Botschaft sagt, „die fenced YAML-Beispiele … wurden deshalb
  **einzeln** gegen den Validator gefahren", und nennt anschließend zwei Läufe —
  beide auf dem §5-Block „Weitere Module". Der ebenfalls in diesem Commit
  geänderte `structure`-Beispiel-Block (zwei neue Schlüssel) ist nicht genannt.
  Der Reviewer hat diese N+1-te Form gefahren: sie hält (Exit 1, kein Exit 2),
  ein Defekt folgt aus der Lücke also nicht — die Reichweite der Aussage ist
  dennoch größer als die belegte Menge.
- `verifizierbar`: nein — Botschaft gegen Proben-Menge.
- `klasse`: Botschaft verallgemeinert über die genannte Proben-Menge hinaus

---

## Negativbefunde (geprüft, ohne Befund)

- **Ventil-Zahlen am Code, nicht an der Doku.** Zeilen-Marker: genau vier
  Module referenzieren `ignoreMarker` — `codepaths.go` (Definition),
  `ids.go`, `versions.go`, `diagrams.go`; kein weiterer Treffer in
  `internal/**` außerhalb der Regeln (die Fundstelle in `trace_table.go` ist
  eine Kommentar-Erwähnung, kein Ventil). `exempt-paths`: genau sechs
  YAML-Träger — `rawIDPattern`, `rawCodepaths`, `rawDiagrams`, `rawVersions`
  (+ `rawVersionPattern`), `rawStructure`, `rawMatrix`; `tracked.exempt-targets`,
  `targets.exempt-targets` und `commits.exempt-pattern` sind andere Schlüssel.
  **Beide Zahlen (4 und 6) stimmen**, ebenso die genannten Modulnamen.
- **Die Ränder der Aussage.** `structure`, `matrix`, `hostpaths`, `pins`,
  `spans`, `citations`, `tracked`, `targets`, `planning` kennen den Zeilen-Marker
  nicht (kein Code-Treffer) — die Aufzählung „alle übrigen kennen ihn nicht"
  trifft zu; nur ihre Begründung nicht (F-1).
- **`--print-config`.** `config_template.go` führt `diagrams.exempt-paths`,
  das Zeilen-Ventil als Kommentar und den `versions`-Alternativblock bereits
  (aus slice-122/124); keine Aussage über die Marker-Reichweite, also kein
  Spiegel-Rückstand. Modulliste im Template: 20 Module, deckungsgleich mit
  `operations.md`.
- **`--suggest-config` / Glossar / FAQ.** Kein weiterer Ort mit einer
  Reichweiten-Behauptung; die FAQ-Antwort ist mitgezogen.
- **§4-Checkliste der Release-Prozedur, Punkt 1 (`version.md`).** §Aktuell auf
  `v0.63.0`/2026-08-23, neue §Verlauf-Zeile, Anker **gewandert**: `grep -c 'a id='`
  über `version.md` = **1**, auf `v0.63.0`. Die `v0.62.0`-Zeile hat ihn verloren.
- **Punkt 2 (gepinnte `ghcr`-Verweise).** Kein `ghcr…:v0.62.0` mehr in lebender
  Doku; `versions` läuft im Default-Profil scharf und ist Teil des grünen
  `doc-check`.
- **Punkt 3 (CHANGELOG-Schnitt).** `[Unreleased]` ist leer, der bisherige
  `Changed`-Block (slice-112) liegt jetzt unter `[0.63.0]`, drei `Added`-Einträge
  für slice-122/123/124, ein `Changed`-Eintrag für slice-125. Keine
  Referenz-Link-Definitionen im Dokument — Form wie in allen Vorgänger-Releases.
- **Punkt 4 (Handbuch).** Kopfstempel 1.56 / `v0.63.0` / 2026-08-23; §11-Zeile
  **chronologisch unten** (letzte Tabellenzeile, zusätzlich durch die eigene
  `structure`-Regel `table-order: asc` gate-gedeckt); §5-Schema, §5-Ventil-Absatz,
  sechster Überraschungs-Punkt, §6-Zeilen für `diagrams`/`versions`/`structure`
  und §8-FAQ vorhanden und untereinander widerspruchsfrei.
- **Punkt 4 (beide READMEs).** DE und EN tragen dieselben drei Änderungen
  (Marker-Absatz, `diagrams`-/`versions`-Bullet, `structure`-Bedingungszahl
  sechs → acht) in gleicher Reihenfolge und gleicher Aussage; die
  Bedingungszahl acht deckt sich mit §6 und mit dem §2-Schema.
- **Punkt 4 (`operations.md`).** Modul-Enumeration der `--enable`/`--disable`-Zeile
  führt alle 20 Module (Abgleich gegen `config_template.go`); die Optionen-Tabelle
  braucht keinen Zuwachs — das Release bringt **keine** neue CLI-Option.
- **Punkt 4 (fünf Datumsstempel).** `version.md` §Aktuell, `version.md` §Verlauf,
  `CHANGELOG`-Überschrift, Handbuch-Kopf und Handbuch-§11 tragen einheitlich
  **2026-08-23**. Orakel: `git log -1 --date=short --format=%ad` = 2026-08-23
  (Commit-Zeitstempel `Sun Aug 23 06:31:24 2026 +0200`). Die Welle begann am
  2026-08-22 — das Datum ist also **nicht** von der Vorgängerzeile abgeschrieben,
  wie die Prozedur es fordert. Die Historien von Lastenheft und Spezifikation
  datieren korrekt auf den Tag ihrer Änderung (2026-08-22, aus slice-122/123/124).
- **Digest-Zeile.** Handbuch §2 trägt weiter den Vorgänger-Digest; das ist
  prozedurkonform (der Digest entsteht erst nach dem Tag) und in der Botschaft
  ausdrücklich als Backfill angekündigt. Der umgebende Text behauptet nicht,
  dass es der Digest der aktuellen Version sei.
- **Gate-Behauptungen der Botschaften.** `make planning-check` (im Move-Commit
  behauptet) ist Teil des selbst gefahrenen `make gates` — Exit 0, 447 Dateien,
  0 Befunde. `make doc-check` Exit 0 und `make gates` Exit 0 reproduziert;
  „Coverage 94,80 % gegen Schwelle 93 %" reproduziert (`coverage-gate`:
  „Coverage 94.80% erfüllt Schwelle 93%"). **Alle drei Gate-Behauptungen sind
  wahr.**
- **Fenced YAML §5 „Weitere Module" gegen den Validator.** Wie gedruckt:
  **Exit 1** (79 Befunde), kein Exit 2 — deckt die Botschaft. Mit
  einkommentierter `patterns`-Liste **zusätzlich** zur Kurzform: **Exit 2** mit
  der benannten Ursache („Kurzform … und `versions.patterns` zugleich gesetzt"),
  genau wie der Kommentar im Block warnt. Mit `patterns` **statt** Kurzform:
  **Exit 1** — die auskommentierte Alternative ist beim Einkommentieren gültig,
  die Einrückung trägt, und die zweite Quelle `harness/conventions.md#baseline`
  löst auf (sonst fail-closed Exit 2).
- **Kommentar-Texte gegen die echten Defaults.** `# geprüfte ATX-Ebene (Default:
  Abschnitts-Ebene + 1)` und `# JEDE Überschrift im Abschnitt matcht dieses
  Muster` decken sich wörtlich mit den §2-Schema-Zeilen; `exempt-paths gilt je
  Paar` deckt sich mit `versions.patterns`; „beides gesetzt ⇒ Exit 2" ist am
  Lauf belegt.
- **§6-Modul-Tabelle gegen die Spezifikation.** `diagrams`-, `versions`- und
  `structure`-Zeile decken sich mit §2/§4; `section-heading-mismatch` ist in der
  Grund-Code-Spalte ergänzt und existiert als `SPEC-067`.
- **Lifecycle-Move `254e8c5` (MR-013 / AGENTS.md §3.3).** Reiner `git mv`
  (Rename 100 %, Slice-Body byte-unverändert) plus genau die zwei gekoppelten
  Verweise: Roadmap-Ruhe-Marker entfernt (in-progress nicht mehr leer) und der
  Zeiger im Wellendokument von `open/` nach `in-progress/`. `grep` über
  `slice-125`: kein Verweis zeigt mehr auf den alten Pfad. Die MR-013-Ausnahme
  ist für die Gegenrichtung geschrieben, ihre Begründung (atomare Kopplung an
  `planning-check`) trägt hier gleichermaßen.
- **Out-of-Scope-Treue.** Der Diff berührt ausschließlich Doku, Register und
  CHANGELOG; kein `internal/**`, kein `Makefile`, **kein** `.d-check.yml` — die
  im Slice zugesagte „keine Funktions-Änderung, keine Profil-Umstellung" hält.
- **ADR-0058 auf `Proposed`.** Korrekt an dieser Stelle: Präzedenz ADR-0057 —
  dort wurde der Status erst im Closure-Commit (`29df048`, „Closure-Notiz
  slice-105 + welle-77-results, ADR-0057 Accepted") gesetzt, also **nach**
  Release-Prep (`2b5d7f4`) und Digest-Backfill (`6e97867`). welle-82 ist offen;
  es fehlt hier nichts.
- **Übergabe aus dem Vorgänger-Review.** Der einzige an slice-125 übergebene
  Punkt („der Zeilen-Marker wirkt nicht nur für `codepaths` und `ids`") ist in
  allen vier dort genannten Dokumenten eingearbeitet; F-1/F-2 betreffen die
  **Begründung** und die **nicht** genannten Spiegel, nicht die Übergabe selbst.
- **BEO-006 (pfad-selektiver Commit).** Beide Commits decken sich mit ihrem
  `--stat`: der Move-Commit trägt drei Pfade und nennt drei, der Release-Commit
  trägt sechs und nennt sechs. Kein undeklarierter Mitreisender.
- **§11-Zeile 1.55.** Sie enthält weiterhin die falsche Aussage — als
  historische Chronik dessen, was 1.55 änderte, und mit der ausdrücklichen
  Korrektur in der unmittelbar folgenden Zeile 1.56. Kein Befund: eine
  Historien-Zeile rückwirkend umzuschreiben wäre der schlechtere Zustand.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 4 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Exklusivitäts-Kriterium statt Aufzählung —
aus dem Anlass, nicht aus dem Bestand · Spiegel-Liste endet an der
Dokument-Rolle · benannte Nicht-Wirkung nur im Spec-Stratum, nicht in der
Nutzer-Doku · Beispiel-Block ohne den Fence-Schutz, den seine Form suggeriert ·
Botschafts-Zahl benennt eine andere Menge als die gemessene ·
Beispiel-Vorlage demonstriert eine still inerte Bedingungs-Kombination ·
Begründung gilt nicht am demonstrierten Ort · Botschaft verallgemeinert über
die genannte Proben-Menge hinaus

## Verdikt

**Merge-blockierend: ja** — zwei MEDIUM. Beide sind vor dem Tag zu klären,
nicht danach: F-1 und F-2 sind konsumentensichtbare Aussagen in genau den
Dokumenten, die dieses Release ausliefert, und beide tragen dieselbe Klasse
wie die Lehre, mit der die Welle sich selbst erklärt. Die vier LOW und zwei
INFO blockieren nicht.

Bemerkenswert für den Steering-Loop: die Klasse aus F-1 ist innerhalb dieser
Welle zum **vierten** Mal aufgetreten — dreimal als zu weite
Exklusivitäts-Aussage (slice-122/123/124), jetzt als zu weites
**Kriterium**, das dieselbe Aussage nur eine Abstraktionsstufe höher
wiederholt. Die Gegenmittel-Frage lautet damit nicht mehr „ist die
Aufzählung vollständig", sondern „gilt der genannte Grund exklusiv".

**Übergabe:** Findings an den Implementer; die Finding-Klassen zusätzlich in
die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein
Lauf-Beleg und ersetzt keine Verifikation.
