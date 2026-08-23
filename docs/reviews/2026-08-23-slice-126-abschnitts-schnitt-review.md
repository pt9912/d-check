# Review-Report: slice-126 — Ventil-Gefälle `citations`, §5-Abschnitts-Schnitt, BEO-011

**Review-Art:** Design-/Doku-Review — geprüft wird die Nutzer-Doku, die
Release-Prozedur und das Beobachtungs-Register gegen Lastenheft, Spezifikation,
ADR, die Hard Rules **und gegen den Code**; unabhängiger Reviewer ohne Anteil an
der geprüften Arbeit.

**Gegenstand:** `d2123e3..HEAD` — die drei Commits des Slice:
`aefdd08` (Slice angelegt), `3100089` (Lifecycle-Move `open/` → `in-progress/`),
`1930197` (die Arbeit).

**Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) @ 1.9.0 ·
**Modell-ID:** `claude-opus-5[1m]` · **Datum:** 2026-08-23

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan [slice-126](../plan/planning/in-progress/slice-126-handbuch-abschnitts-schnitt.md),
  Welle [welle-82](../plan/planning/welle-82-config-flaechen.md),
  [Roadmap](../plan/planning/in-progress/roadmap.md)
- [`spec/lastenheft.md`](../../spec/lastenheft.md) — `DC-FA-CITE-001`,
  `DC-FA-CODE-001`, `DC-QA-02`
- [`spec/spezifikation.md`](../../spec/spezifikation.md) —
  §`DC-FA-CITE-001.a` (Schritte 1–5), §`DC-FA-CODE-001.a` (Schritt 1 Marker,
  Schritt 6 Zeilen-Check), §2 `<modul>.scope.roots`/`.ignore`
- [ADR-0058](../plan/adr/0058-konfigurations-flaechen-additiv-weiten.md)
  (Status `Proposed`, §Re-Evaluierungs-Trigger)
- [`docs/user/releasing.md`](../user/releasing.md) §Release-Prep §4,
  [`docs/user/benutzerhandbuch.md`](../user/benutzerhandbuch.md),
  [`docs/user/benutzerhandbuch-standard.md`](../user/benutzerhandbuch-standard.md),
  [`docs/user/operations.md`](../user/operations.md),
  [`README.de.md`](../../README.de.md), [`README.md`](../../README.md),
  [`CHANGELOG.md`](../../CHANGELOG.md)
- Beobachtungs-Register [`observations.md`](../plan/planning/observations.md)
  (BEO-002, BEO-006, BEO-007, BEO-008, BEO-009, BEO-011)
- Hard Rules [`AGENTS.md`](../../AGENTS.md) §3 (§3.3, §3.7), §5;
  [`harness/conventions.md`](../../harness/conventions.md)
  ([MR-013](../../harness/conventions/MR-013-lifecycle-move-buendelung.md),
  [MR-025](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md))
- Vorgänger-Review [2026-08-23 slice-125](2026-08-23-slice-125-release-v0630-review.md)
- Code: `internal/hexagon/core/rules/citations.go`,
  `internal/hexagon/core/rules/codepaths.go`,
  `internal/hexagon/core/model/config.go`,
  `internal/adapter/driven/configyaml/configyaml.go`

**Vom Reviewer selbst gefahren** (Kommando in eine Datei umgeleitet, `$?` direkt
gelesen — BEO-007; Fixture innerhalb der Repo-Wurzel unter `.rv-fixture/`,
danach gelöscht, `git status` sauber):

- `make build` → **Exit 0**
- Probe 1 — eigene Vier-Dateien-Fixture, Module `[codepaths, citations]`,
  `codepaths.check-lines: true`, ein Lauf → **Exit 1**, zwei Befunde:
  `b_nomarker.md:1 citation-out-of-range` (Kontrolle, ohne Marker) und
  `c_cite_marker.md:3 citation-out-of-range` (Direktive **mit**
  `d-check:ignore` auf derselben Zeile). Die Datei mit
  Inline-Code-Zeilen-Referenz **und** Marker (`a_marker.md`) schweigt.
  ⇒ Ventil-Gefälle des Slice **bestätigt** (A stumm, B meldet, C meldet).
- Probe 2 — dieselbe Fixture, zusätzlich
  `citations: {scope: {roots: [docs], ignore: ["docs/c_*.md"]}}` → **Exit 1**,
  **nur noch ein** Befund (die `codepaths`-Kontrolle); der
  `citations`-Befund ist **weg**. ⇒ `citations` **hat** einen
  Konfigurations-Block und ein datei-weites Ventil (siehe F-1).
- Probe 3 — `citation-mismatch` mit `d-check:ignore` auf der Direktiven-Zeile →
  **Exit 1**, Befund bleibt (Marker wirkt auch für den dritten Grund-Code nicht).
- `docker run … --print-config` gegen die Repo-Wurzel → `citations` steht in der
  emittierten Modul-Liste.
- `make gates` → **Exit 0** (acht Gates, `449 Datei(en) geprüft, 0 Befund(e)`,
  `coverage-gate: OK — Coverage 94.80% erfüllt Schwelle 93%`).

**Verdikt: blockierend** — kein HIGH, **6 MEDIUM**, **3 LOW**, **1 INFO**.

---

## Findings

### F-1 — „`citations` trägt gar kein Ventil" ist am Code falsch und widerspricht der Liste 140 Zeilen darüber

- `kategorie`: MEDIUM
- `quelle`: `DC-FA-CITE-001`, §2-Schema `<modul>.scope.roots`/`.ignore`,
  `DC-FA-CONF-002`
- `pfad`: `docs/user/benutzerhandbuch.md:1107`, `:1258`, `:1261`, `:1866`
- `befund`: Das Handbuch sagt an vier Stellen, `citations` habe
  „**gar kein** Ventil" (1107), „keinen Konfigurations-Block, keine
  Ausnahmeliste" (1258), es bleibe nur „behoben … oder das Modul wird
  abgeschaltet" (1261) und in der §6-Zeile „**kein Ventil**" (1866). Gemessen
  (Probe 2) schaltet `citations.scope.ignore` den Befund still, das
  YAML-Schema kennt den Block (`configyaml.go`, `Citations *rawScopeOnly`
  mit `Scope *rawScope`), der Kern-Kommentar sagt „die Struktur trägt nur den
  Scope-Platzhalter, damit Enable/Scope wie bei jedem Modul greifen"
  (`config.go`), und die §5-Ventil-Achsen-Liste desselben Kapitels führt
  15 Zeilen über der neuen Aussage `scan.ignore` als erste von **vier**
  Ventil-Achsen, die jede Datei jedem Modul entzieht. Die 1.55-Zeile
  desselben Handbuchs benennt für `diagrams` denselben Sachverhalt ehrlich
  („dort schneidet nur `diagrams.scope` bzw. `scan.ignore`"); für `citations`
  fehlt diese Hälfte. Versagen: ein Leser mit **einer** unreparierbaren
  `d-check:cite`-Stelle folgt der als erschöpfend formulierten Alternative und
  schaltet das Modul repo-weit ab, statt die eine Datei zu skopieren.
- `verifizierbar`: ja — Probe 2 (eigener Image-Lauf mit
  `citations.scope.ignore`, Exit 1, Befund verschwindet); kein Gate fängt es,
  weil die Aussage Prosa ist.
- `klasse`: absolute Ventil-Verneinung ohne Zählung der Ventil-Achsen

### F-2 — Die neue `planning`-Überschrift nennt zwei von drei Fähigkeiten — genau der Defekt, den der Slice entfernt

- `kategorie`: MEDIUM
- `quelle`: Slice-Plan §4 DoD („**jede** Überschrift nennt, was unter ihr
  steht"), `DC-FA-PLAN-001`
- `pfad`: `docs/user/benutzerhandbuch.md:1297` (Überschrift), `:1311`
  („**Dritte Fähigkeit**"), `:2068` (§11-Zeile 1.57)
- `befund`: Die eingefügte Überschrift lautet „Planning-Lifecycle und
  Closure-Notizen prüfen (Modul `planning`)"; der Abschnitt darunter trägt in
  seinem eigenen Text drei Fähigkeiten — die Lifecycle-Invariante, die
  Closure-Notizen und, ausdrücklich so benannt, die „**Dritte Fähigkeit**,
  opt-in über `waves.dir`: die Wellen-Register gegen die Wellen-Dateien" mit
  `waves.mode` und den vier Grund-Codes `wave-drift`, `wave-preview-exists`,
  `wave-results-missing`, `wave-unregistered` (rund 34 der 101 Zeilen). Die
  §6-Modul-Tabelle (`:1875`) führt dieselbe dritte Fähigkeit. Slice-Plan §1,
  Drift-Log-Zeile, Commit-Botschaft und §11-Zeile 1.57 sagen übereinstimmend
  „`planning` (zwei Fähigkeiten)"; die §11-Zeile schließt mit „jede Überschrift
  nennt ihren Inhalt". Versagen: wer in §5 die `waves.dir`-Konfiguration sucht,
  findet sie unter einer Überschrift über Lifecycle und Closure-Notizen — der
  Zensus „genau **einer**" gilt nach dem Schnitt für einen neuen Fall, den der
  Schnitt selbst erzeugt hat.
- `verifizierbar`: ja — `grep -n "Fähigkeit" docs/user/benutzerhandbuch.md`
  im Abschnitt 1297–1397; kein Gate fängt es (`structure` zählt Tasks, nicht
  Themen — vom Slice §3 selbst als Out-of-Scope benannt).
- `klasse`: Überschrift nennt eine Teilmenge ihres Inhalts

### F-3 — Die Begründung des Ventil-Gefälles steht in keinem Vertrag und ist eine neue, ungezählte Behauptung

- `kategorie`: MEDIUM
- `quelle`: `DC-FA-CITE-001` / §`DC-FA-CITE-001.a`; BEO-011 Ableiter 2
  (im selben Commit geschrieben)
- `pfad`: `docs/user/benutzerhandbuch.md:1264`
- `befund`: Der neue Absatz schließt mit „Das ist Absicht — eine Direktive zu
  setzen ist ein ausdrücklicher Akt des Autors, und eine ausdrücklich gesetzte
  Prüfung soll nicht zeilenweise zurückgenommen werden können." Weder
  `DC-FA-CITE-001` (Beschreibung, Akzeptanzkriterien, Out-of-Scope) noch
  §`DC-FA-CITE-001.a` (Schritte 1–5) noch ADR-0058 sagen etwas über eine
  gewollte Marker-Freiheit; das Lastenheft sagt zu Direktiven nur, „die
  Platzierungsregeln folgen der bestehenden `d-check:ignore`-Konvention", und
  `DC-FA-CODE-001` §Out-of-Scope nennt „Opt-out-Marker für andere Module"
  schlicht als nicht zugesagt. Der im selben Commit angelegte Eintrag BEO-011
  formuliert den Ableiter „Wer eine falsche Exklusivität **repariert**,
  begründet die reparierte Menge **nicht** — sie ist eine *benannte Liste*,
  kein ableitbares Kriterium; eine Begründung ist eine neue, ungezählte
  Behauptung." Versagen: der Slice benennt in §3 selbst den Change Request
  „bekommt `citations` ein Ventil?" als offene Frage — die Nutzer-Doku
  beantwortet sie vorab mit „Absicht", ohne dass ein Stratum das trägt.
- `verifizierbar`: nein — kein Gate; nachprüfbar durch `grep` nach der Aussage
  in `spec/lastenheft.md`, `spec/spezifikation.md` und
  `docs/plan/adr/0058-konfigurations-flaechen-additiv-weiten.md` (kein Treffer).
- `klasse`: nachträgliches Kriterium für eine benannte Liste

### F-4 — Die Spiegel der korrigierten Aufzählung enden bei den zwei Stellen, die aufgefallen sind

- `kategorie`: MEDIUM
- `quelle`: [MR-025](../../harness/conventions/MR-025-spiegel-vor-dem-editieren.md)
  („für eine Wortlaut-Präzisierung ist der Ableiter das `grep` nach dem
  **alten** Wortlaut über den ganzen Baum"), BEO-002
- `pfad`: `README.de.md:179`, `README.md:177`, `docs/user/operations.md:130`,
  `spec/spezifikation.md:1088`, `spec/lastenheft.md:970`,
  `docs/plan/adr/0058-konfigurations-flaechen-additiv-weiten.md:318`,
  `CHANGELOG.md:92`
- `befund`: Der Slice erklärt die „benannte Liste"-Notiz und die
  CHANGELOG-Enumeration für unvollständig, weil sie `citations` auslassen, und
  zieht **zwei** Stellen nach (`docs/user/benutzerhandbuch.md:1107`,
  `CHANGELOG.md:75`). Derselbe Satz steht wortgleich in `README.de.md`,
  `README.md` (englisch), `docs/user/operations.md`, `spec/spezifikation.md`
  und — ranghöher als die Nutzer-Doku — `spec/lastenheft.md`; die
  CHANGELOG-Zeile 92 (slice-125-Eintrag) führt dieselbe Aufzählung. Die
  parallele Enumeration der musterlosen Module steht im
  §Re-Evaluierungs-Trigger von ADR-0058 („Ein Modul **ohne** konfigurierbare
  Muster (`hostpaths`, `pins`, `spans`) braucht wiederholt eine Ausnahme"),
  die ADR ist `Proposed` und damit editierbar. Der Slice-Plan enthält keine
  MR-025-Spiegel-Liste. Versagen: der Re-Evaluierungs-Trigger, der zünden
  müsste, wenn `citations` wiederholt eine Ausnahme braucht, nennt `citations`
  nicht; und wer der Source Precedence folgt und im Lastenheft nachschlägt,
  liest die Fassung ohne `citations` — dieselbe Mechanik, die der
  slice-125-Review als F-2 gemeldet hat.
- `verifizierbar`: ja —
  `grep -rn "benannte Liste" --include=*.md .` liefert die unbearbeiteten
  Spiegel; kein Gate fängt es (`links`/`anchors`/`ids` prüfen Auflösung, nicht
  Geltung).
- `klasse`: Spiegel-Liste endet an den aufgefallenen Stellen

### F-5 — Der Closure-Trigger der Welle zählt weiter vier Slices, während §4 fünf führt

- `kategorie`: MEDIUM
- `quelle`: BEO-002 („Semantik-Änderung nur im Körper nachgezogen, Ränder
  bleiben stehen" — vom Slice §7 selbst als „unmittelbar einschlägig" benannt)
- `pfad`: `docs/plan/planning/welle-82-config-flaechen.md:49` gegen `:67`
- `befund`: `aefdd08` fügt die fünfte Zeile in §4 „Slices in dieser Welle" ein
  (slice-126); §3 „Closure-Trigger (Welle schließt)" trägt unverändert „Alle
  **vier** Slices in `done/`; `make fullbuild` grün". Versagen: der
  Wellen-Closure prüft sein eigenes Kriterium wörtlich — mit vier Slices in
  `done/` ist es erfüllt, während slice-126 noch in `in-progress/` liegt; kein
  Gate koppelt die zwei Abschnitte (`planning-check`/`wave-drift` vergleichen
  Roadmap-Zeiger gegen Wellen-Dateien, nicht Zähler innerhalb eines
  Wellendokuments).
- `verifizierbar`: ja — `sed -n '49p;63,67p'` auf dem Wellendokument;
  gate-blind.
- `klasse`: Zähler im Rand nicht mit der Aufzählung im Körper nachgezogen

### F-6 — Die Regel-Weitung in `releasing.md` nennt wieder einen Ort statt der Klasse

- `kategorie`: MEDIUM
- `quelle`: BEO-011 Ableiter 3 („Wer eine **Regel** aus einem Vorfall schreibt,
  benennt zuerst die Klasse und erst dann den Ort" — im selben Commit
  geschrieben)
- `pfad`: `docs/user/releasing.md:51`
- `befund`: Die neue Passage lautet „**Dieselbe Regel gilt für §5**" und
  begründet die Weitung damit, dass die Vorfassung „nur für §4 da" stand,
  „geschrieben für das Kapitel, in dem sie wehtat, statt für die Klasse". Der
  Prüfsatz selbst („nennt die Überschrift alles, was unter ihr steht?") ist
  dokument- und kapitelunabhängig; gebunden wird er trotzdem an genau die zwei
  Kapitel, an denen die Anlagerung bisher gemessen wurde. Nicht erfasst sind
  Handbuch §7/§8, `docs/user/operations.md` (dieselbe Release-Prep-Checkliste
  pflegt dort Modul-Liste und Optionen-Tabelle) und
  `docs/user/benutzerhandbuch-standard.md` §3 „Klare Struktur verwenden", wo
  eine kapitelfreie Form der Regel ihren Ort hätte. Versagen: der nächste
  Release-Prep hängt einen Feature-Absatz an eine bestehende Überschrift in
  `operations.md` oder in Handbuch §7; die Checkliste, die genau das verhindern
  soll, deckt den Ort nicht.
- `verifizierbar`: nein — kein Gate; nachprüfbar durch Lesen der
  §4-Checkliste gegen den Handbuch-Standard.
- `klasse`: Regel aus dem Anlass statt aus der Klasse (eine Kerbe weiter)

### F-7 — Die Gegenüberstellungs-Tabelle lässt das dritte Ventil von `codepaths` aus

- `kategorie`: LOW
- `quelle`: `DC-FA-REF-001` (geteiltes Referenz-Ventil `ignore-refs`),
  `DC-FA-CODE-001`
- `pfad`: `docs/user/benutzerhandbuch.md:1252`
- `befund`: Die Tabelle wird als der für die Reparatur entscheidende
  Unterschied eingeführt und führt zwei Ventil-Zeilen (`exempt-paths`,
  Zeilen-Marker). `codepaths` hat für dieselben Grund-Codes ein drittes:
  `refIgnored(refs, file, rel)` in `checkCodepathTarget` kehrt **vor** dem
  Zeilen-Bereichs-Check zurück (`codepaths.go`), sodass ein
  `ignore-refs`-/`codepaths.ignore-refs`-Treffer auf dem aufgelösten Ziel
  `citation-out-of-range` und `citation-inverted-range` unterdrückt; `citations`
  hat keinen Gegenpart. Versagen: ein Leser, der die Tabelle als vollständig
  liest, sucht das Ziel-Achsen-Ventil nicht und hält den Unterschied für
  kleiner, als er ist.
- `verifizierbar`: ja — Fixture mit `ignore-refs` auf das Ziel und
  `codepaths.check-lines`; kein Gate fängt die Doku-Aussage.
- `klasse`: Vergleichstabelle unvollständig auf der geprüften Achse

### F-8 — „der Diff zeigt vier eingefügte Überschriften, sonst nichts" gegen `git show --stat`

- `kategorie`: LOW
- `quelle`: BEO-009 Richtung (a) / MR-025-Prozedur („`git diff --stat` gegen
  jede Botschafts-Zeile halten")
- `pfad`: Commit `1930197`, Botschafts-Absatz „BEFUND 2"
- `befund`: Die Botschaft sagt „Der Text ist UNVERÄNDERT BEWEGT, nicht
  umgeschrieben -- der Diff zeigt vier eingefügte Überschriften, sonst nichts."
  `git show 1930197 -- docs/user/benutzerhandbuch.md` zeigt **28** eingefügte
  und **3** gelöschte Zeilen; vier davon sind Überschriften, 21 weitere sind
  der neue Ventil-Absatz, der **innerhalb** des ersten der fünf Abschnitte
  liegt und aus dessen genannter Länge (33) herausgerechnet ist. Versagen: die
  in BEO-009 verankerte Gegenprobe (Botschaftszeile gegen `--stat`) schlägt
  fehl und lässt offen, welche der 28 Zeilen zum Schnitt und welche zu Befund 1
  gehören.
- `verifizierbar`: ja — `git show --stat 1930197`.
- `klasse`: Botschafts-Absolutum („sonst nichts") gegen den Datei-Diff

### F-9 — Nutzersichtbare Doku-Änderung ohne `[Unreleased]`-Eintrag, dafür Retro-Edit im veröffentlichten `0.63.0`-Block

- `kategorie`: LOW
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §5 („`CHANGELOG.md` wird bei
  nutzersichtbaren Änderungen gepflegt"), `docs/user/releasing.md` Schritt 3
- `pfad`: `CHANGELOG.md:7` (leerer `[Unreleased]`-Block) gegen `:75`
- `befund`: Der Slice liefert Handbuch-Stand 1.57 (neue Ventil-Tabelle, vier
  neue §5-Überschriften, korrigierte §6-Zeile) und eine neue Regel in der
  Release-Prozedur; `[Unreleased]` bleibt leer, während der bereits getaggte
  Block `## [0.63.0] — 2026-08-23` nachträglich editiert wird. Die Vorgänger
  derselben Art tragen je einen `### Changed`-Eintrag (1.54/slice-112,
  1.55/slice-125, beide unter `0.63.0` bzw. `0.62.0`). Versagen:
  `releasing.md` Schritt 3 schneidet beim nächsten Release den
  `[Unreleased]`-Stand unter die neue Version — dieser Slice erscheint dort
  nicht, obwohl §6 des Plans sagt, er „reist mit dem nächsten Release".
- `verifizierbar`: ja — `sed -n '1,12p' CHANGELOG.md`; kein Gate erzwingt
  CHANGELOG-Einträge.
- `klasse`: nutzersichtbare Änderung ohne Changelog-Eintrag

### F-10 — BEO-011 zählt zwei Instanzen mit, die unter keine seiner drei Ausprägungen fallen

- `kategorie`: INFO
- `quelle`: Slice-Plan §5 („Die Register-Klasse könnte zu breit geraten … ein
  Eintrag, der alles trifft, steuert nichts"), BEO-011 selbst
- `pfad`: `docs/plan/planning/observations.md:11`
- `befund`: Der Eintrag nennt drei Ausprägungen (a Exklusivitäts-Aussage, b
  nachträgliches Kriterium, c Regel nur für den Ort) und beziffert den Zähler
  auf 5. Nachgelesen in den Closure-Notizen: slice-123 und slice-124 tragen (a),
  slice-125 trägt (b). Die beiden übrigen gezählten Instanzen fallen unter
  keine der drei — slice-122 ist eine **unerfüllbare Zusage über die Anzahl von
  Befunden** („aus der Absicht abgeleitet statt aus dem Mechanismus"), und die
  zweite slice-125-Instanz ist eine **nach Dokument-Rolle statt nach Aussage
  gebildete Spiegel-Liste** (MR-025-Klasse, im Register bereits als BEO-002
  geführt). Ausprägung (c) ist unter den fünf gar nicht vertreten (sie liegt
  „in Arbeit"), obwohl der Eintrag „Drei Ausprägungen, alle eingetreten"
  schreibt. Die drei Ableiter — Nachbarn zählen, reparierte Menge nicht
  begründen, Klasse vor Ort — greifen für keine der beiden mitgezählten
  Instanzen. Kein Befund gegen die Zahl 5 als solche: sie ist mit der
  Belege-Spalte und der Aufschlüsselung in der Stand-Zelle konsistent
  auflösbar.
- `verifizierbar`: nein — Urteil, kein `grep`; nachprüfbar durch Lesen der vier
  Closure-Notizen in `docs/plan/planning/done/`.
- `klasse`: Register-Klasse absorbiert Instanzen außerhalb ihrer benannten
  Ausprägungen

---

## Negativbefunde

- **geprüft, ohne Befund — die Ventil-Aussage zu `codepaths` gegen den Code:**
  `CheckCodepaths` verlässt sich bei `exempt-paths` auf `ignored(file, …)` und
  überspringt jede Zeile mit `ignoreMarker`, **bevor** die Spans gelesen werden;
  der Zeilen-Bereichs-Check (Schritt 6) liegt in `checkCodepathTarget` hinter
  dieser Schleife. Beide Zusagen der Tabelle („ja"/„ja") sind wahr und in
  Probe 1 gemessen.
- **geprüft, ohne Befund — `CitationsConfig struct{}` und die Marker-Freiheit:**
  `citations.go` enthält keine Referenz auf `ignoreMarker` und keinen
  `ExemptPaths`-Zweig; `model.CitationsConfig` ist leer. Die Aussagen „kein
  `exempt-paths`", „kein Zeilen-Marker" sind wahr (Probe 1 und 3).
- **geprüft, ohne Befund — Rand `citation-mismatch`:** derselbe Grund-Code-Pfad,
  derselbe Marker-Befund (Probe 3, Exit 1); die Doku behauptet für
  `citation-mismatch` nichts Abweichendes.
- **geprüft, ohne Befund — Verlustfreiheit des §5-Schnitts:**
  `git show 1930197 -- docs/user/benutzerhandbuch.md` zeigt genau **3**
  gelöschte Zeilen (Kopf-Versionsstempel, der erweiterte Ventil-Satz, die
  erweiterte §6-`citations`-Zeile), alle drei in erweiterter Form
  wieder eingefügt. Kein Absatz umformuliert, kein Beispiel verändert.
- **geprüft, ohne Befund — Zensus vor dem Schnitt:** dreizehn `###`-Abschnitte
  in §5 (`awk`-Extraktion `## 5.` … `## 6.` auf `d2123e3`), Längen
  15/6/11/16/**155**/39/**183**/42/43/54/11/**203**/78 — drei über 90 Zeilen,
  die drei genannten (`trace` 203, *Weitere Module* 155, dieser 183). Der
  183er-Abschnitt trägt sechs Module. Alle Zahlen des Slice reproduzieren.
- **geprüft, ohne Befund — Zensus nach dem Schnitt:** siebzehn `###`-Abschnitte,
  die fünf neuen mit 54/13/16/101/28 Zeilen; abzüglich der 21 Zeilen des neuen
  Ventil-Absatzes ergibt der erste 33, in Summe 191 = 183 + 4×2 Zeilen für die
  neuen Überschriften. Die Botschafts-Zahlen „33 / 13 / 16 / 101 / 28" stimmen.
- **geprüft, ohne Befund — Anker-Quer-Verweise:** kein Dokument des Repos
  verlinkt auf einen Anker in `docs/user/benutzerhandbuch.md`
  (`grep -rn "benutzerhandbuch.md#"` außerhalb `docs/reviews/` leer); die
  dokumentinternen `](#…)`-Links zeigen ausschließlich auf `##`-Kapitel.
  `make doc-check` in `make gates` → Exit 0.
- **geprüft, ohne Befund — §Inhalt-Liste:** sie führt nur die elf
  `##`-Kapitel; die vier neuen `###`-Überschriften erzeugen dort keine
  Nachziehpflicht.
- **geprüft, ohne Befund — Prosa-Quer-Verweise über die neuen Schnittkanten:**
  die drei „oben"/„siehe unten"-Stellen im geschnittenen Bereich
  (`Singleton-Prädikat oben`, `dieselbe Abschnitts-Mechanik wie bei den
  Closure-Notizen`, `fail-closed, Exit 2, siehe unten`) zeigen sämtlich auf
  Text **innerhalb** ihres neuen Abschnitts; kein Verweis überquert eine neue
  Überschrift.
- **geprüft, ohne Befund — dieselbe Klasse in §4:** siebzehn `###`-Aufgaben,
  jede mit Modul- bzw. Flag-Nennung in der Überschrift; keine trägt ein
  zweites, ungenanntes Modul (die scheinbaren H2/H3 bei `F-1 —
  Repository-Struktur` und `## Kreuzverweis-Konsistenz` liegen in
  Beispiel-Fences). §6 ist eine Tabelle ohne Unter-Überschriften.
- **geprüft, ohne Befund — MR-013 für beide Lifecycle-Commits:** `aefdd08`
  legt die Slice-Datei in `open/` an und fügt Roadmap-Drift-Zeile und
  Wellen-Zeile hinzu, ohne den Ruhe-Marker zu ziehen (korrekt, `open/` ist
  nicht beansprucht). `3100089` ist ein reiner `git mv` (Slice-Datei
  `0` geänderte Zeilen, Rename erkannt) und bündelt genau die gekoppelten
  Verweise: Roadmap-Flip („Nichts in Arbeit" raus) und die
  `open/`→`in-progress/`-Pfadkorrektur im Wellendokument. `make planning-check`
  läuft in `make gates` grün.
- **geprüft, ohne Befund — Drift-Log-Zugehörigkeit:** die Sektion erlaubt
  ausdrücklich „Slice oder Welle umgehängt oder **neu geschnitten**"; die neue
  Zeile ist genau das und folgt der Form der drei Vorgängerzeilen
  („slice-110 neu geschnitten", „slice-106 neu geschnitten"). Sie protokolliert
  keine Schließung und keinen erreichten Meilenstein.
- **geprüft, ohne Befund — Zustandsfeld-Form (§3.7) der neuen
  BEO-011-Stand-Zelle:** sie nennt Zustand („Schwelle erreicht"), Zähler,
  Aufschlüsselung, die offene mechanische Form und drei Ableiter — dieselbe
  Form wie BEO-008/BEO-009/BEO-010; keine Entstehungs-Chronik, kein
  „Letzte Änderung"-Datum.
- **geprüft, ohne Befund — Out-of-Scope-Treue:** kein `.go`-, `Makefile`- oder
  `.d-check.yml`-Hunk im gesamten Diff (`git diff --stat d2123e3..HEAD`:
  7 Dateien, alle Markdown); `citations` bekommt kein Ventil; `trace` und
  *Weitere Module* sind unberührt; keine mechanische Prüfung der
  Abschnitts-Aufteilung eingeführt.
- **geprüft, ohne Befund — Botschafts-Zahlen gegen eigene Läufe:** „183 Zeilen /
  sechs Module", „33 / 13 / 16 / 101 / 28", „nur 3 gelöschte Zeilen", „Zensus
  über alle dreizehn", die Fixture-Ergebnisse (A stumm, B meldet, C meldet,
  beide Läufe Exit 1) und „`make gates` Exit 0 (acht Gates, 449 Dateien,
  0 Befunde, Coverage 94,80 % gegen Schwelle 93 %)" reproduzieren sämtlich.
  Die einzige Botschafts-Abweichung ist F-8.
- **geprüft, ohne Befund — Referenz-Richtung (SDP):** der Diff enthält keinen
  neuen Provenance-Marker und keinen Abwärts-Token in einem ADR- oder
  Spec-Körper; die Slice→ADR-/Spec-Verweise zeigen aufwärts.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 6 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** absolute Ventil-Verneinung ohne Zählung der
Ventil-Achsen · Überschrift nennt eine Teilmenge ihres Inhalts ·
nachträgliches Kriterium für eine benannte Liste · Spiegel-Liste endet an den
aufgefallenen Stellen · Zähler im Rand nicht mit der Aufzählung im Körper
nachgezogen · Regel aus dem Anlass statt aus der Klasse (eine Kerbe weiter) ·
Vergleichstabelle unvollständig auf der geprüften Achse · Botschafts-Absolutum
gegen den Datei-Diff · nutzersichtbare Änderung ohne Changelog-Eintrag ·
Register-Klasse absorbiert Instanzen außerhalb ihrer benannten Ausprägungen

---

## Verdikt

**Merge-blockierend: ja** — kein HIGH, aber **6 MEDIUM**.

Die tragende Beobachtung dieses Laufs ist, dass der Slice seine eigene Lehre
dreimal wiederholt, während er sie aufschreibt: F-1 setzt eine
Exklusivitäts-Aussage („gar kein Ventil"), ohne die vier Ventil-Achsen zu
zählen, die dasselbe Kapitel fünfzehn Zeilen darüber aufführt (Ausprägung a);
F-3 begründet die reparierte Menge, was BEO-011 Ableiter 2 im selben Commit
verbietet (Ausprägung b); F-6 zieht die Regel auf ein zweites Kapitel statt auf
ihre Klasse (Ausprägung c). F-2 erzeugt beim Beheben eines
Teilmengen-Titels einen neuen. F-4 und F-5 sind Ränder, die stehenblieben,
obwohl §7 des Plans BEO-002 als „unmittelbar einschlägig" benennt.

Der geprüfte Kern hält: das Ventil-Gefälle **existiert** und ist in beide
Richtungen gemessen, der §5-Schnitt ist verlustfrei (drei gelöschte Zeilen,
alle drei erweitert wieder eingefügt), der Zensus reproduziert Zahl für Zahl,
beide Lifecycle-Commits sind MR-013-konform, und `make gates` ist grün
(Exit 0, selbst gefahren).

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in
die Slice-Closure §7 und von dort in den Zähler — F-1/F-3/F-6 sind
BEO-011-Instanzen und betreffen dessen eigenen Zähler, F-4 ist eine
BEO-002-Instanz. Dieser Report ist ein Lauf-Beleg und ersetzt keine
Verifikation (DoD-/Spec-Konformität prüft der Verifier separat).
