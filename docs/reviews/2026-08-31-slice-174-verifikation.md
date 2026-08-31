# Verifikation slice-174 — die Register-Deckung bekommt ihren Wächter

**Art:** Verifikation.
**Gegenstand:** [slice-174](../plan/planning/in-progress/slice-174-register-deckung.md), Diff `f6bdb7e~1..HEAD` (Pfade `internal spec docs/plan docs/user .d-check.yml`) — `f6bdb7e` (Plan angelegt), `fa5c4dd`/`fa6d63d` (Beanspruchung, `open/` → `in-progress/`), `827d1da` (die drei offenen Entscheide getroffen), `c5a7569` (Implementierung + Lastenheft 0.81.0), `d3b0dac` (ADR-0079 + Konfiguration vorbereitet, **auskommentiert**), `f33a4ec` (fünf Review-Befunde F-1…F-5 eingearbeitet — dieser Commit-Hash weicht von dem im Auftrag genannten `c39f8ea` ab; `c39f8ea` ist im lokalen Repo nicht auflösbar, `f33a4ec` trägt exakt die im Auftrag beschriebenen fünf Befunde und den Report-Dateinamen), `fe4fe55` (§3-Korrektur der Anzahl-Achsen-Begründung).
**Modell-ID:** `claude-sonnet-5`.
**Datum:** 2026-08-31.
**Rolle:** unabhängiger Verifier — geprüft wird gegen **DoD und Spec** (Baseline-Regelwerk [`modul-08-agentenrollen.md` §Welche Rolle braucht welche Artefaktklasse](../../.harness/baseline/v5.12.0/regelwerk/modul-08-agentenrollen.md#welche-rolle-braucht-welche-artefaktklasse-modul-8)), nicht gegen Plan/ADR-Maintainability. Der Review-Report `docs/reviews/2026-08-31-slice-174-register-deckung-code-r1.md` wurde bewusst **nicht** gelesen.
**Prüfgrundlage:** Slice-Plan `docs/plan/planning/in-progress/slice-174-register-deckung.md` §4 (DoD), §1, §2, §2a, §3; [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) (Version 0.81.0, §Register-Deckung + vier Akzeptanzkriterien) im Lastenheft; `SPEC-079` in [`spec/spezifikation.md`](../../spec/spezifikation.md) §4; [ADR-0079](../plan/adr/0079-register-deckung-zaehlt-linktext.md); Handbuch §4.20 und die `planning`-Zeile der Modul-Tabelle in [`docs/user/benutzerhandbuch.md`](../user/benutzerhandbuch.md); das gebaute Produkt (`make build`, `d-check:latest`, Image-ID `sha256:7b8cb29549c4…`), nicht der Quelltext allein — zusätzlich das **publizierte** `ghcr.io/pt9912/d-check:v0.71.1` für einen Vergleichslauf (§3).
**Nicht verändert:** kein Repo-Artefakt außer dieser Datei — `git status --short` ist vor und nach jedem Lauf leer. Alle Mess-Läufe gegen Wegwerf-Repos und eine Kopie des Planning-Baums unter `/tmp/.../scratchpad/slice174-verify`, außerhalb dieses Repos.

---

## 1. DoD-Tabelle (§4 des Slice-Plans)

| # | Behauptet | Gemessen | Verdikt |
|---|---|---|---|
| 1 | Eine in `done/` zitierte `BEO-<NNN>` ohne Registerzeile erzeugt einen Befund mit eigenem Grund-Code; Richtung und Scan-Menge stehen in der Anforderung, nicht nur im Code | `observation-unregistered` existiert (`model/finding.go`), eigener Reason-Text in `AllReasons()`/`reasonTexts()`. Lastenheft nennt die Richtung ausdrücklich ausschließend („Nur diese Richtung … die Umkehrung … ist ausgeschlossen") und die Scan-Menge samt [`AGENTS.md` §3.8](../../AGENTS.md#38-ein-modul-verspricht-nur-über-das-was-es-scannt)-Verweis. Am Produkt bestätigt (§2 unten) | **erfüllt** |
| 2 | Umkehr-Probe gemessen: konstruiertes `BEO-999` ⇒ Befund; eigener Bestand ⇒ 0 Befunde. Beide Ausgaben **in der Closure-Notiz** | Beide Messungen selbst reproduziert (§2, Läufe „eigener Bestand" und „konstruierte Kennung") — Substanz stimmt exakt mit den Commit-Botschaften überein. Die **Ablage in der Closure-Notiz** ist ein Closure-Schritt: §9 des Plans ist eine leere Überschrift, der Slice liegt in `in-progress/` | Substanz **erfüllt**; Beleg-Ablage **noch nicht fällig** |
| 3 | Träger-Entscheid (`planning`-Bedingung ↔ eigenes Modul) begründet in einer ADR, samt gemessener Verwerfung des `ids`-Wegs | [ADR-0079](../plan/adr/0079-register-deckung-zaehlt-linktext.md) §Entscheidung + §Verglichene Alternativen trägt Träger-Wahl, Erkennungs-Form (mit Messung 366/293/112/5) und den verworfenen `ids`-Anker-Weg (mit Messung 1450/1042/708/626) | **erfüllt** |
| 4 | Anforderung, `.a`-Verfeinerung mit Grund-Code, Profil-Aktivierung und Handbuch sind nachgezogen | Anforderung (Lastenheft 0.81.0, vier AKs) **erfüllt**; `.a`-Verfeinerung (`SPEC-079` in `spezifikation.md` §4) **erfüllt**; **Profil-Aktivierung nicht eingelöst** — `planning.observations` steht in `.d-check.yml` vollständig **auskommentiert**, die Fähigkeit ist im eigenen Profil **inert**; **Handbuch enthält eine nachweislich fehlschlagende Anleitung** — §3 unten. Beides ist im Plan (§2a, §5, §4) an keiner Stelle als Abweichung angekündigt oder begründet | **nicht erfüllt** (zwei von vier Teilzusagen) |
| 5 | `make gates` und `make fullbuild` grün (Exit explizit); unabhängiger Review; Verifikation — beide in eigenen Kontexten | Selbst gefahren: `make gates` → **Exit 0**, Schlusszeile `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`, Coverage 94,70 % ≥ 93 %. `make fullbuild` → **grün**, `[fullbuild] green — image-hash sha256:7b8cb29549c4…`. Review-Commit `f33a4ec` trägt den Report in eigenem Kontext (bewusst ungelesen). Diese Datei ist die Verifikation, in eigenem, vom Review getrenntem Kontext | **erfüllt** (durch diesen Report selbst mit eingelöst) |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register fortgeschrieben; jedes Risiko aus §5 mit Ausgang; die drei Paarungen geprüft | §9 leer; alle vier §5-Risiken tragen `*(bei Closure)*`; `git diff f6bdb7e~1..HEAD -- docs/plan/planning/observations.md` ist leer (Register unangetastet) | **noch nicht fällig** |

**Zusammenfassung:** 3 von 6 Zeilen **erfüllt** (1, 3, 5), 1 Zeile **Substanz erfüllt / Beleg-Ablage noch nicht fällig** (2), 1 Zeile **noch nicht fällig** (6 — Closure-Schritt, korrekt: Slice liegt in `in-progress/`), **1 Zeile nicht erfüllt** (4 — zwei der vier Teilzusagen sind nicht eingelöst und nirgends im Plan als Abweichung dokumentiert).

## 2. Die vier Akzeptanzkriterien der Anforderung, selbst ausgeführt (Wegwerf-Repo, Produkt-Binary)

Aufbau: eigenes Repo mit `docs/plan/planning/observations.md` (Register), `docs/plan/planning/done/*.md` (Zitier-Verzeichnis), `.d-check.yml` mit `planning.observations.register`/`dirs`. Alle Läufe `docker run --rm -v <repo>:/repo:ro d-check:latest --enable planning [--json]` gegen das lokal aus HEAD gebaute `d-check:latest`.

| Kriterium | Aufbau | Ausgabe | Erwartung |
|---|---|---|---|
| **Happy Path** | Register mit `BEO-001`, eine Datei zitiert `[`BEO-001`](../observations.md)` | `d-check: 3 Datei(en) geprüft, 0 Befund(e)`, Exit 0 | ✓ kein Befund |
| **Happy Path, ohne Schlüssel** | dieselbe Datei, `.d-check.yml` ohne `observations`-Block | `d-check: 3 Datei(en) geprüft, 0 Befund(e)`, Exit 0 — byte-identisch | ✓ DC-QA-02 |
| **Boundary (Zitat gegen Beispiel)** | eine Datei mit drei Zeilen: Linktext `[`BEO-999`](../observations.md)` (Zeile 3), reines Code-Span `` `BEO-999` `` (Zeile 5), Fließtext `BEO-999` ohne Backticks (Zeile 7) | `--json`: **genau zwei** Befunde, Zeile 3 und Zeile 7, `target: "BEO-999"` — Zeile 5 (Code-Span) meldet **nicht** | ✓ exakt wie zugesagt |
| **Boundary (nur erste Zelle deklariert)** | Register-Zweitzeile `\| BEO-777 \| … erwähnt BEO-888 querverweisend … \|`, `dirs` leer (Default = Registerverzeichnis) | zwei Befunde: die Fließtext-Erwähnung `BEO-888` **im Register selbst** (Zeile 6, zweite Zelle) und `BEO-555` aus einem Fließtext-Satz im Register — beide melden, obwohl sie im Register stehen, weil nur die **erste** Zelle deklariert | ✓ exakt wie zugesagt |
| **Negative — Register fehlt** | `register: docs/plan/planning/fehlt.md` | 1 Befund, Message „Register … fehlt oder ist unlesbar (fail-closed)", Exit 1 | ✓ |
| **Negative — `dirs`-Eintrag unlesbar** | `dirs: [docs/plan/planning, docs/plan/gibtsnicht]` | 1 Befund, Message „Zitier-Verzeichnis … fehlt oder ist unlesbar (fail-closed)", Exit 1 | ✓ |
| **Negative — ungültiges `pattern`** | `pattern: '[a-'` | `d-check: error: .d-check.yml: planning.observations.pattern "[a-" ist kein gueltiger Ausdruck: …`, **Exit 2** | ✓ Konfigurationsfehler, kein Befund |

Alle sieben Läufe stimmen exakt mit der Lastenheft-AK-Formulierung überein, einschließlich der Zeilen- und Ziel-Angaben.

## 2a. Die zwei besonders geprüften Zusagen (Aufträge C und die Deduplikation)

- **Zwei verschiedene ungedeckte Kennungen auf einer Zeile ⇒ zwei Befunde.** Zeile `zwei: [BEO-777](../observations.md) und [BEO-888](../observations.md)` → `--json` liefert **zwei** Einträge mit `target: "BEO-777"` bzw. `"BEO-888"`, Exit 1. Deckt sich mit `TestObservationsTwoIDsOnOneLineDoNotCollide` (F-1-Fix) — selbst am Produkt bestätigt, nicht nur im Unit-Test gelesen.
- **Eigener Bestand, gegen die echte Baum-Kopie gemessen (nicht nur behauptet).** `docs/plan/planning` dieses Repos in ein Wegwerf-Verzeichnis kopiert, `.d-check.yml` mit `observations.register`/`dirs: [docs/plan/planning]` scharf gestellt (kein Repo-Artefakt verändert): `d-check: 573 Datei(en) geprüft, 0 Befund(e)`, Exit 0 — deckt sich mit der im Plan (§1) und in ADR-0079 behaupteten Messung „eigener Bestand 0 Befunde", jetzt aber gegen den **heutigen** Stand (statt gegen den Stand zur Bau-Zeit) reproduziert.
- **Konstruierte Kennung an derselben Kopie.** Eine Zeile `[`BEO-9999`](../observations.md)` an `slice-174-register-deckung.md` angehängt → 1 Befund, `target: "BEO-999"` (das `\d{3}`-Muster matcht nicht-verankert die ersten drei Ziffern — Eigenschaft des überall gleichen Default-Musters, kein Fehler dieser Fähigkeit), Exit 1.

## 3. Spec-vs-Produkt

- **Mechanismus.** Lastenheft-Text („Gezählt werden Prosa und Linktext; ein reines Inline-Code-Span nicht") und Code (`inCodeOnly`, ruft `IDOccurrenceExempt` — dieselbe Funktion, die `ids` für seine Linkpflicht nutzt) stimmen überein; §2 oben bestätigt es am Produkt in allen drei Formen.
- **Scan-Mengen-Grenze steht in der Anforderung, nicht nur im Code** ([`AGENTS.md` §3.8](../../AGENTS.md#38-ein-modul-verspricht-nur-über-das-was-es-scannt)): Lastenheft-Zeile „Das Modul verspricht nur über seine Scan-Menge … eine Kennung in einem nicht gelisteten Verzeichnis wird nicht geprüft" — vorhanden, wörtlich mit Verweis auf §3.8.
- **Fail-closed an drei Rändern** — alle drei live reproduziert (§2, Negative-Zeilen); Meldungsform und Exit-Codes stimmen mit der Lastenheft-Beschreibung überein.
- **Ein reproduzierter, konkreter Defekt: das Handbuch-Beispiel in §4.20 schlägt gegen das aktuell publizierte Image fehl.** §4.20 nennt wörtlich `ghcr.io/pt9912/d-check:v0.71.1` als Aufruf-Beispiel für `planning.observations`. `v0.71.1` (Commit `75e905b`) ist ein direkter Vorfahre von `f6bdb7e` (Start dieses Slice) — die Fähigkeit existiert in diesem Image **nicht**, und die Config-Dekodierung ist strikt (`KnownFields(true)`, bereits zum Stand `75e905b` so). Selbst gepullt und exakt das Handbuch-Beispiel nachgefahren:
  ```
  $ docker run --rm -v <repo>:/repo:ro ghcr.io/pt9912/d-check:v0.71.1 --enable planning
  d-check: error: .d-check.yml: yaml: unmarshal errors:
    line 6: field observations not found
  Exit 2
  ```
  Wer die Handbuch-Anleitung heute exakt befolgt, bekommt **nicht** das dokumentierte Verhalten, sondern einen Konfigurationsfehler. Das ist keine Marginalie — die Anleitung ist der zentrale Inhalt von §4.20. Der Fund ist die Kehrseite der (an sich vertretbaren) Entscheidung, die Profil-Aktivierung „nach Release und Digest-Backfill" zu vertagen (`.d-check.yml`-Kommentar, Commit `d3b0dac`): die Handbuch-Seite wurde **so geschrieben, als wäre die Fähigkeit mit dem aktuell gepinnten Tag bereits nutzbar** — sie ist es nicht. Dieselbe Klasse Blind-Spot ist für Handbuch-Tag-Beispiele bereits aus früheren Release-Preps bekannt (sie werden von keinem Gate geprüft, `make gates`/`make fullbuild` bauen und prüfen ausschließlich `d-check:latest` aus dem aktuellen Quellstand, nie das gepinnte publizierte Image).

## 4. Plan-vs-Code, beide Richtungen

**Richtung 1 — im Diff, aber vom Plan nicht angekündigt:**

- Der Diff-Bereich `f6bdb7e~1..HEAD` (wie im Auftrag angegeben) enthält **drei fremde Vorgänge**, die chronologisch zwischen slice-174-Commits liegen, aber zu anderen Slices/CRs gehören: `a4bc209` (slice-187, Handbuch-URL-Zeiger — **dieser Commit liefert auch den kompletten §4.20-Abschnitt und die Modul-Tabellenzeile für die Register-Deckung**, ohne slice-174/DC-FA-PLAN-001/ADR-0079 in Commit-Botschaft oder Trace-IDs zu nennen), `bcc8557` (CR-Antwort zu einem unabhängigen Kanon-Thema, im eigenen Commit ausdrücklich als „slice-174 ist davon NICHT berührt" markiert) und `4a45376` (slice-188 geschnitten). Die ersten beiden sind für die DoD-Bewertung relevant: **das Handbuch-Deliverable von slice-174 landet nicht unter slice-174s eigener Traceability**, sondern versteckt in einem fremden Commit — wer `git log --grep slice-174` oder die Trace-Matrix als Änderungs-Historie liest, findet die Handbuch-Änderung nicht. Inhaltlich ist sie korrekt (§2/§3 oben), die Zuordnung ist es nicht.
- `docs/plan/planning/in-progress/roadmap.md` (−2 Zeilen): die deklarierte MR-013-Begleitung des Claim-Moves, kein unangekündigter Fund.

**Richtung 2 — vom Plan zugesagt, aber im Diff nicht wiederzufinden:**

- **Profil-Aktivierung** (§2 Punkt 6, §4 DoD Punkt 4): der Plan sagt sie ohne Einschränkung zu; im Diff steht sie nur **vorbereitet, auskommentiert** — nirgends im (weiterhin offenen) Plantext als bewusste Abweichung vermerkt. Zum Vergleich: der einzige Präzedenzfall im selben Modul (`waves.mode: many`) wurde ebenfalls zunächst auskommentiert, aber über eine **eigene, benannte Folge-Slice** (`slice-111`) nachgezogen. Für die Deaktivierung der Register-Deckung existiert **keine** vergleichbare Spur — kein `BEO`-Eintrag, kein offener Folge-Slice, keine Erwähnung in §5 (Risiken) des Plans. Das Risiko „Wächter bleibt für immer inert" ist damit ungenannt und ungetragen.
- **Handbuch „nachgezogen"** trifft inhaltlich zu, aber mit dem in §3 belegten Fehler: das Beispiel ist zum aktuell gepinnten Tag nicht lauffähig. Weder der Plan noch die Commit-Botschaft von `a4bc209` (die das Handbuch tatsächlich lieferte) benennen dieses Spannungsverhältnis zwischen „Handbuch fertig" und „Fähigkeit im eigenen Profil noch nicht scharf".

**§3 „Ausdrücklich NICHT" — Grenzen gehalten:**

1. **Beleg-Form-Prüfungen (Form/Anzahl/Lage)** — kein entsprechender Code-Pfad im Diff; unverändert ausgeklammert, mit `slice-188` als benanntem Nachzügler für die Anzahl-Achse.
2. **Umkehrung** — kein Code, der „jede Zeile ist zitiert" prüft; bestätigt durch Lesen von `planning_observations.go` (nur `declaredObservationIDs` → `citedWithoutRow`, eine Richtung).
3. **Urteil, ob zwei Beobachtungen dieselbe sind** — kein Vergleichs-Code zwischen Kennungen; unberührt.
4. **Nachrüsten von Bestands-Zitaten** — `git diff f6bdb7e~1..HEAD -- docs/plan/planning/observations.md` ist leer, Register unangetastet.

## 5. Gate-/Mess-Läufe (echte Ausgabe)

| Lauf | Ergebnis |
|---|---|
| `make build` | Neubau aus HEAD, `d-check:latest`, Image-ID `sha256:7b8cb29549c4553300d87ee9e382dfbb3630068b63cf945f3a095f2f2caa3b1e` |
| `make test` | alle Pakete `ok`, u. a. `internal/hexagon/core/rules` |
| `make gates` | **Exit 0** — `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; Coverage 94,70 % ≥ 93 %; semgrep 55 Regeln, 0 Befunde; `--enable planning` gegen den eigenen Baum: `645 Datei(en) geprüft, 0 Befund(e)` |
| `make fullbuild` | **grün** — `[fullbuild] green — image-hash sha256:7b8cb29549c4…`; RTM: 50 Anforderungen, **0 Waisen** |
| Happy Path (§2) | `3 Datei(en) geprüft, 0 Befund(e)`, Exit 0 |
| Happy Path ohne Schlüssel (§2) | `3 Datei(en) geprüft, 0 Befund(e)`, Exit 0 — byte-identisch |
| Boundary Zitat/Beispiel (§2) | 2 Befunde (Zeile 3, Zeile 7), Zeile 5 (Code-Span) fehlt korrekt, Exit 1 |
| Boundary nur erste Zelle (§2) | 2 Befunde (`BEO-888`, `BEO-555`), Exit 1 |
| Negative — Register fehlt (§2) | 1 Befund, Exit 1 |
| Negative — `dirs` unlesbar (§2) | 1 Befund, Exit 1 |
| Negative — ungültiges `pattern` (§2) | Exit 2, Meldung nennt den Schlüssel |
| Zwei IDs auf einer Zeile (§2a) | 2 Befunde, verschiedene `target` | 
| Eigener Bestand, Baum-Kopie (§2a) | `573 Datei(en) geprüft, 0 Befund(e)`, Exit 0 |
| Konstruierte Kennung, Baum-Kopie (§2a) | 1 Befund, Exit 1 |
| Handbuch-Beispiel gegen `ghcr.io/pt9912/d-check:v0.71.1` (§3) | `d-check: error: .d-check.yml: yaml: unmarshal errors:\n  line 6: field observations not found`, **Exit 2** |
| `git status --short` (vor/nach allen Läufen) | leer |

## 6. Verdikt

**Nicht bestanden.**

Der Mechanismus selbst ist sauber gebaut und deckt sich präzise mit Spezifikation und ADR: alle vier Akzeptanzkriterien der Anforderung, die Boundary-Fälle, alle drei fail-closed-Ränder und die im Review gefundene Deduplikations-Falle (F-1) sind am **laufenden Produkt** — nicht nur im Unit-Test — bestätigt, mit eigens erzeugten Belegen einschließlich einer Baum-Kopie des heutigen Bestands (573 Dateien, 0 Befunde). `make gates` und `make fullbuild` sind mit echtem Exit 0 gegengeprüft.

Der Grund für „nicht bestanden" ist DoD-Punkt 4: **zwei der vier zugesagten Teil-Deliverables sind nicht eingelöst, und keine Abweichung ist im Plan dokumentiert.**

1. Die **Profil-Aktivierung** ist im eigenen `.d-check.yml` vollständig auskommentiert — die Fähigkeit bleibt im eigenen Repo inert. Das folgt einem plausiblen Präzedenzfall (`waves.mode: many`), aber anders als dort gibt es **keine** benannte Folge-Slice, keinen Beobachtungs-Register-Eintrag und keine Erwähnung im Plan (§2a/§5), die sicherstellt, dass die Aktivierung nicht schlicht vergessen wird.
2. Das **Handbuch-Beispiel in §4.20** ist gegen das aktuell publizierte, im Beispiel selbst genannte Image `ghcr.io/pt9912/d-check:v0.71.1` **nachweislich lauffähig-falsch** — selbst gepullt und ausgeführt, Ergebnis Exit 2 mit „field observations not found" statt des dokumentierten Befund-Verhaltens. Das ist kein hypothetisches Risiko, sondern ein reproduzierter Fehlschlag der genauen Anleitung, die Nutzer heute lesen würden.

Beide Punkte hängen zusammen: die Handbuch-Seite wurde geschrieben, als sei die Fähigkeit mit dem aktuellen Tag nutzbar, während dieselbe Arbeit an anderer Stelle (`.d-check.yml`-Kommentar) explizit festhält, dass sie es noch nicht ist. Vor der Closure sollte entweder (a) das Handbuch-Beispiel bis zur tatsächlichen Freischaltung zurückgehalten oder mit einer Verfügbarkeits-Anmerkung versehen werden, und (b) ein benannter Trigger (Folge-Slice oder Beobachtungs-Register-Eintrag) die spätere Scharfschaltung der Profil-Konfiguration tragen — analog zum `mode: many`-Präzedenzfall. Zusätzlich sollte die Handbuch-Herkunft (Commit `a4bc209`, unter slice-187 statt slice-174 gebucht) für die Traceability nachgezogen oder in der Closure-Notiz benannt werden.

Kein Befund betrifft die Kern-Logik der Register-Deckung selbst — dort ist alles Geprüfte korrekt.
