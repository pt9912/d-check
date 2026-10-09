# Verifikation: slice-265, `reviews` — Zusage-Muster, benannte Kennungen, Leerlauf, Unterverzeichnisse

- **Rolle:** Verifier (Modul 8/11) — geprüft wird der Stand gegen DoD (§2), Plan (§3 samt
  Plan-Änderungen vor dem Code, nach R1 bis R4) und §6, **nicht** der Diff gegen die Entscheidungen
  (das haben R1 bis R5 getan).
- **Gegenstand:** `slice-265`, Commits `86a0d887` bis `c6e3f5f9` (alle tragen `slice-265`;
  dazwischen die Pin-Commits `cc818666`, `798d256d` ohne Slice-Bezug). Code-Stand `HEAD` =
  `c6e3f5f9`, Arbeitsbaum sauber.
- **Eingang:** Slice-Plan in `in-progress/`; eingehender CR von `sf-connector`; ausgehender CR an
  den Kurs samt Antwort; Reports R1 bis R5; Lastenheft 0.101.3 `DC-FA-RVW-001`; Spezifikation
  `DC-FA-RVW-001.a`, §2-Zeilen `reviews.*`, `SPEC-081`.
- **Datum:** 2026-10-09 · **Modell-ID:** claude-opus-5-5

## Verdikt

**Liefer-Punkte der DoD (§2, Punkte 1 bis 3) bestätigt — mit Belegen am Code, am Image und durch
bewusstes Brechen.** Punkt 4 (Review, Verifikation) ist mit R1 bis R5 und diesem Bericht erfüllt.
**Offen ist nur Punkt 5** — die Closure selbst (siehe unten). Keine DoD-Verletzung, kein Befund
über LOW.

## Sensor-Belege (selbst gefahren)

| Lauf | Ergebnis |
|---|---|
| `make gates` (auf `c6e3f5f9`) | Exit 0; Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; `lint` „0 issues."; `coverage-gate: OK — Coverage 94.70% erfüllt Schwelle 93%`; `doc-check` 1041/1047 Dateien, 0 Befunde |
| `make review-coverage` | `d-check: 1041 Datei(en) geprüft, 0 Befund(e)`, Exit 0 |

## DoD Punkt 1 — Zusage- und Kennungs-Muster, Leerlauf-Schalter

Am Code: `model.DefaultPromisePattern` =
``[Uu]nabhängiger Review|(?:\[[ xX]\]|[;,])[ \t\r\n]*Review durchgeführt``; `ReviewsConfig` trägt
`PromisePattern`, `Match`, `RequirePromises`, `Recursive`, `SkipPattern`; `applyReviews` weist
explizit leeres und nicht kompilierendes `promise-pattern`, `match` außerhalb `id`/`name` und nicht
kompilierendes `skip-pattern` mit Exit 2 und Schlüsselnamen ab; YAML-`null` ist abwesend (Zeiger).
Unter `match: id` liefert `reviewFinding` für eine Zusage ohne `slice-<NNN>` einen Befund mit Hinweis
auf `match: name` (vorher `continue`).

Am Image (Exit 2): `promise-pattern: ""` ⇒ `reviews.promise-pattern ist leer — es träfe jeden
DoD-Punkt …`, Exit 2.

**Akzeptanzkriterien ↔ Tests** (Lastenheft 0.101.3): jedes Kriterium hat einen Test —
Happy/Negative/Bullet/kein Kandidat/leere Menge/exempt/Stubs/Modul-aus in `reviews_test.go`
(Bestand); Vorlagen-Form → `TestReviewsDefault_ErkenntVorlagenForm`,
`TestReviewsDefault_VorlagenFormNurAmTeilanfang` (inkl. `;` + Umbruch); eigenes Muster →
`TestReviewsPromisePattern_EigeneForm`; Negative (Default) → `TestReviewsDefault_Negativfaelle` +
Negativteil von `…NurAmTeilanfang` (inkl. umbrochener Verneinung); Slug-Kennung →
`TestReviewsMatch_BenannteKennung`, `…NameWortgrenze`, `…NameWortgrenzeUnicode`; keine Zusage →
`TestReviewsRequirePromises_LeerlaufRot`; Unterverzeichnisse/Stubs → `TestReviewsRecursiveUndSkipPattern`,
`…SkipDirsBleibenUnbetreten`, `…UnlesbaresUnterverzeichnis`, `…OhneLeerlauf`; Exit 2 →
`TestDecode_ReviewsMusterFehler`, `TestDecode_ReviewsPromisePatternNullIstAbwesend`.

## DoD Punkt 2 — Abstieg und Inhalts-Ausnahme

Dieselbe Semantik wie `closureCandidates`/`closureSkip` (slice-263): `SKIP_DIRS` ausgenommen,
Symlink auf ein Verzeichnis nicht verfolgt, `skip-pattern` gegen den rohen Inhalt, unlesbare Datei
bleibt Kandidatin. Die eine Abweichung (unlesbares Unterverzeichnis meldet je Verzeichnis und bricht
den Walk nicht ab) ist Plan-Änderung nach R1 und im Modell-Kommentar sowie in Spezifikation Schritt
2/5 benannt. Symlink-Probe am Image: `done/link -> ../extern` mit Zusage in `extern/` ⇒ 0 Befunde
unter `recursive: true` (nicht verfolgt).

## DoD Punkt 3 — Ausgabe ohne die neuen Schlüssel, Vorlage, CR-Abnahme

**Altes gegen neues Image**, selbst gebaut: alt aus `git archive 97be2feb` (Commit vor
`1837d8d2`), neu aus `git archive HEAD`, beide `VERSION=0.0.0-dev`. Je Bestand zwei Läufe (voll mit
`--enable reviews`; Fokus wie `review-coverage`) × drei Formen (Standard, `--json`, `--doctor`),
stdout, stderr und Exit getrennt mit `cmp`:

| Bestand | Exit alt/neu | Ergebnis |
|---|---|---|
| Repo (`git archive HEAD`) | 0/0 | byte-identisch, 6/6 |
| klein — Phrase-Zusagen, nummeriert, ein Slug ohne Zusage, Stub unter `welle-01/`, ohne Vorlagen-Form | 1/1 | byte-identisch, 6/6 |
| derselbe ohne `reviews-dir` | 1/1 | byte-identisch, 6/6 |
| leere `done/` | 1/1 | byte-identisch, 6/6 |
| klein + ein Slice mit Vorlagen-Zeile | 1/1 | **Abweichung, erwartet:** genau ein zusätzlicher Befund `slice-005-v.md:5 … review-missing …` (Default-Änderung 1) |
| klein + Slug-Slice `slice-w3-login.md` mit Phrase | 1/1 | **Abweichung, erwartet:** genau ein zusätzlicher Befund `Review-Zusage, aber keine slice-<NNN>-Kennung … — match: name …` (Default-Änderung 2) |

Jede der beiden Default-Änderungen ist durch einen Test belegt und durch Brechen bestätigt (M1, M2
unten).

**Vorlage `--print-config`** (am neuen Image gelesen): beschreibt den DoD-Punkt mit Folgezeilen,
beide Default-Formen samt Position, „nicht ‚kein Review durchgeführt'", alle fünf neuen Schlüssel,
`match: name` für Slug-Kennungen. Die falsche Aussage „eine Zeile, die *Review* nennt" ist weg.

**CR-Abnahme am gebauten Image** — Bestand im Stil von `sf-connector`: `slice-abholzustand-und-lesepfad-publication.md`,
`slice-w2-login-flow.md`, `slice-039-abholzustand-und-lesepfad-publication.md`, DoD-Zeile
`- [x] Review durchgeführt, Report unter \`docs/reviews/\` liegt vor` mit Folgezeile; Konfiguration
`match: name`, `require-promises: true`:

- ohne Reports: 3 × `review-missing`, Exit 1 (**rot**);
- Reports `2026-10-08-<basis>-review.md`, `…-korrekturen-review.md`, `…-verifikation.md`: 0 Befunde,
  Exit 0 (**grün**), `--json` `"exitCode": 0`;
- einen Report entfernt: genau dieser Slice meldet, Exit 1;
- das CR-Wunschmuster `'^\s*- \[[ x]\] Review durchgeführt'` als `promise-pattern`: greift;
- `promise-pattern: 'Code-Review erledigt'` ohne Treffer: Leerlauf-Befund `keine Review-Zusage unter 3
  Kandidat(en) — require-promises …`;
- ohne `match: name` (Default `id`): beide Slug-Slices melden mit Hinweis auf `match: name` — die
  Beobachtung des Kurses, reproduziert;
- altes Image auf derselben Konfiguration: Exit 2 (`field match not found`) — erwartet.

**Eigenes Gate, nachgemessen** (Probe-Konfiguration mit leerem `reviews-dir` auf einer Kopie):
alt 0 erkannte Zusagen, neu 11 in `done/` flach, 26 mit `recursive: true`; mit dem echten
`docs/reviews/` und `recursive: true` 0 Befunde — alle 26 sind gedeckt. Deckt sich mit R5 (26).

## Bewusstes Brechen (12 Mutationen, je eine Wegwerf-Kopie aus `git archive HEAD`, `make test` dort)

| # | Mutation | rot (Exit 2), Grund |
|---|---|---|
| M1 | Default zurück auf `[Uu]nabhängiger Review` | `…ErkenntVorlagenForm` („rot: review-missing auf Zeile 3 … erwartet, got []"), `…CheckboxImCodeblockZaehlt`, `…VorlagenFormNurAmTeilanfang` („drei Zusagen … got []") |
| M2 | Zusage ohne Kennung wieder still (`return nil, true`) | `…Match_BenannteKennung` („Befund mit Hinweis auf match: name erwartet, got []") |
| M3 | Default `…|Review durchgeführt` ohne Verankerung | `…NurAmTeilanfang` Negativteil: alle vier (`Adaptions-…`, `kein …`, `das …`, umbrochen) melden |
| M4 | Alternative „Zeilenanfang" (`|\n`) zurück | `…NurAmTeilanfang`: genau `slice-007-g.md` (umbrochene Verneinung) |
| M5 | Wortgrenze nach dem Basisnamen entfernt | `…NameWortgrenze`, `…NameWortgrenzeUnicode` |
| M6 | Wortgrenze nur ASCII | `…NameWortgrenzeUnicode` (`slice-a-grö` / `slice-a-größe`) |
| M7 | `require-promises`-Zweig tot | `…RequirePromises_LeerlaufRot` |
| M8 | `SKIP_DIRS` unter `recursive` ignoriert | `…SkipDirsBleibenUnbetreten` (`node_modules/slice-002-b.md` meldet) |
| M9 | Leerlauf-Unterdrückung neben unlesbarem Verzeichnis entfernt | `…UnlesbaresUnterverzeichnisOhneLeerlauf` (zweiter Befund „leere Pruefmenge") |
| M10 | Exit 2 für explizit leeres `promise-pattern` entfernt | `TestDecode_ReviewsMusterFehler` („reviews.promise-pattern ist leer … got <nil>") |
| M11 | `match`-Validierung entfernt | `TestDecode_ReviewsMusterFehler` („reviews.match \"pfad\" … got <nil>") |
| M12 | `recursive` wirkungslos | `…RecursiveUndSkipPattern`, `…UnlesbaresUnterverzeichnis`, `…OhneLeerlauf` |

Jede Mutation fällt mit der behaupteten Meldung, keine nur „irgendwie". Kopien und Test-Images
danach entfernt.

## Befunde

- **V-1 (LOW) — Die Abnahme des CR in der Form, die der Kurs braucht, steht in keinem einzelnen
  Test.** DoD Punkt 3 verlangt „die Abnahme-Fälle des CR als Test".
  `TestReviewsDefault_ErkenntVorlagenForm` deckt das nummerierte Beispiel aus dem CR (rot/grün, Default).
  Die Kombination Slug-Kennung + Vorlagen-Form + `match: name` + `require-promises` deckt kein einzelner
  Test, sondern nur die Teile (`…BenannteKennung` mit der Phrase, `…LeerlaufRot`). Am Image habe ich
  sie bestätigt (oben). Vorschlag: ein Test mit genau dieser Konfiguration, oder die Lesart „Teile
  genügen" in der Closure-Notiz benennen.
- **V-2 (INFO) — Kein CLI-Abnahmetest fährt die neuen Schlüssel von YAML bis Befund.** Der Decode-Test
  hält die Durchreichung, die Regeltests setzen das Modell direkt. Die Image-Läufe oben schließen die
  Lücke für diesen Stand; einen Test, der sie dauerhaft hält, gibt es nicht.
- **V-3 (INFO) — „Symlink auf ein Verzeichnis wird nicht verfolgt" ist für `reviews` nicht getestet**
  (Spezifikation Schritt 2). Die Image-Probe bestätigt es. Die Zusage hängt am Filesystem-Adapter, den
  `planning` ebenso nutzt.

## Offen bis zur Closure (DoD Punkt 5, kein Verifier-Befund)

- §7 Closure-Notiz, Lerneintrag, Register und drei Paarungen sind leer.
- **§6 Risiko ohne Ausgang:** Die Präfix-Deckung unter `match: name` (`slice-a-foo` wird durch den
  Report zu `slice-a-foo-bar` gedeckt) ist eingetreten, aber nur **als benannte Grenze**: Lastenheft,
  Spezifikation und Code-Kommentar nennen sie, und `TestReviewsMatch_NameWortgrenze` hält sie fest.
  Gewählt werden muss einer der drei Ausgänge. *Entfallen* trägt nicht, weil das Risiko weiter
  eintreten kann. Bleibt die Grenze stehen, passt eher *weiter offen* ins Register als *eingetreten*
  ohne Folge-Slice.
- **Eingehender CR:** Das Kopffeld `**Stand:**` sagt noch „angenommen, in Planung"; „Die
  Entscheidung je Punkt folgt mit der Umsetzung" ist noch nicht eingelöst (Punkte 1 bis 3: umgesetzt;
  Punkt 4: slice-266).
- Ausgeschlossen bleibt laut §1, die Konfiguration dieses Repos (`recursive`, `skip-pattern`,
  `require-promises`) zu ändern. Die Messung oben zeigt: Mit `recursive: true` sind alle 26 Zusagen
  gedeckt, die Umstellung nach dem Release wäre also grün.

## Negativbefund

Gesucht und nicht gefunden: weitere Abweichungen alt/neu über die zwei benannten hinaus (36
Vergleichsläufe); Plan-Änderungen ohne Spiegel im Code; Lastenheft-Kriterien ohne Test; ein Test,
der beim Brechen nur aus fremdem Grund fällt. Auch ohne Befund blieben der Spiegel in der
`.d-check.yml` und die Hilfe-Zeile von `review-coverage` im `Makefile` (beide nennen beide
Default-Formen) sowie die Sensor-Datei, deren Grenze 2 die Messung bestätigt (neu 11 bzw. 26
erkannte Zusagen statt 0).
