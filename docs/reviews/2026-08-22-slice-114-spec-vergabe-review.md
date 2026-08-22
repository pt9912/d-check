# Review-Report: slice-114 — `SPEC-*`-Vergabe in der Spezifikation + `structure`-Wächter für §2

**Datum:** 2026-08-22 · **Review-Art:** Code-/Konfigurations-Review (geprüft gegen Slice-Plan, Wellendokument welle-80 D1/D2, Baseline v5.7.0 `grundlagen-source-precedence.md` §ID-Schema als Klammer + §Vergabe und `templates/spec/spezifikation.template.md`, `DC-FA-STRUCT-001` samt Spezifikation §DC-FA-STRUCT-001.a, `DC-FA-ID-001`/§DC-FA-ID-001.a, `DC-FA-ANCH-001.a`, Hard Rules AGENTS §3.4/§3.5/§3.7 + §5, MR-000, MR-025) mit eigener Gegenprobe am gebauten Image
**Gegenstand:** Commit `ed8daa7` (Range `9407b37..ed8daa7`) — Vollvergabe 66 Struktur-IDs in `spec/spezifikation.md`, neue `structure`-Regel in `.d-check.yml`, angepasster AllReasons-↔-§4-Lockstep in `internal/hexagon/core/app/diagnose_test.go`, zwölf retargetete Anker-Verweise; **vor** der Closure, kein Release
**Skill:** `.harness/skills/reviewer.md` @ 1.5.0 · **Modell-ID:** `claude-opus-5[1m]`
**Eingangs-Kontext:** Slice-Plan `docs/plan/planning/in-progress/slice-114-spec-vergabe-spezifikation.md` (§2 Vorgehen, §3 NICHT-Liste, §5 Risiken); Wellendokument `docs/plan/planning/welle-80-struktur-ids.md` (D1/D2, §3 Closure-Trigger, §6 Out-of-Scope); vendorte Baseline `.harness/baseline/v5.7.0/regelwerk/grundlagen-source-precedence.md` §ID-Schema als Klammer + §Vergabe und `.harness/baseline/v5.7.0/templates/spec/spezifikation.template.md`; `.harness/baseline/v5.7.0/regelwerk/grundlagen-harness-dateien.md` §Was ein Kommentar trägt; Spezifikation §DC-FA-STRUCT-001.a (Schritte 2–7), §DC-FA-ANCH-001.a, §DC-FA-ID-001.a; `harness/conventions.md` §MR-000, `harness/conventions/MR-025-spiegel-vor-dem-editieren.md`; `internal/hexagon/core/rules/markdown.go` (`parseATXHeading`) als Ist-Wahrheit der Heading-Lexik; Vorgänger-Report `docs/reviews/2026-08-22-slice-113-struktur-id-konvention-review.md`. Nicht erhalten: DoD-Abhakung (Verifikations-Rolle). Kein `make`-Target im echten Repo (paralleler Gate-Lauf) — alle Proben liefen als Image-Lauf gegen eine `.git`-freie Baum-Kopie außerhalb des Repos, Exit je Lauf explizit in eine Datei umgeleitet und geprüft (Arbeitsregel BEO-007).

## Findings

### F-1 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** `DC-FA-STRUCT-001` / Spezifikation §DC-FA-STRUCT-001.a Schritt 3 (Abschnitts-Findung über die ATX-Lexik aus §DC-FA-ANCH-001.a)
- **pfad:** `.d-check.yml:251` (die `forbid-pattern`-Zeile der §2-Regel)
- **befund:** Die Bedingung verankert mit `(?m)^###` am Zeilen-Anfang und verlangt danach genau ein Leerzeichen, während die Heading-Lexik desselben Moduls (`parseATXHeading`, `internal/hexagon/core/rules/markdown.go:391-403`) die Zeile erst per `TrimLeft(" \t")` beschneidet und Leerzeichen **oder** Tabulator als Trenner akzeptiert. Eine §2-Überschrift mit bis zu drei führenden Leerzeichen oder mit Tabulator nach `###` ist damit für `anchors` eine echte Überschrift mit eigenem Slug, entgeht der Kennungs-Pflicht aber lautlos: in der Gegenprobe blieb der Lauf bei eingefügtem `   ### Eingerueckt ohne Kennung` grün (Exit 0, 0 Befunde) und der aus `README.md` daraufgesetzte Verweis `#eingerueckt-ohne-kennung` löste auf — der Wächter meldet grün, obwohl eine kennungslose, verlinkbare §2-Sektion existiert. Der Kommentar über der Regel sagt die Zusage in ihrer vollen Form („jede ###-Sektion des Schema-Abschnitts traegt eine SPEC-Kennung"); die Grenze, die er benennt, ist eine andere (RE2 ohne Lookahead).
- **verifizierbar:** ja — Negativ-Probe 6, Läufe H/T/V: `make doc-check` bleibt Exit 0, obwohl die Invariante verletzt ist.
- **klasse:** waechter-regex-enger-als-die-lexik-des-waechters

### F-2 · MEDIUM

- **kategorie:** MEDIUM
- **quelle:** Repo-Praxis „Review-Reports sind Lauf-Belege" (`.d-check.yml:42-49`, Ventil `ignore-refs`) / `MR-025`
- **pfad:** `docs/reviews/2026-06-19-slice-029-doctor-json.md:131,139-140`
- **befund:** Die drei Vorkommen im Report sind **Inline-Code**, keine Markdown-Links, und `docs/reviews/**` ist für `ids` wie für `codepaths` per `exempt-paths` ausgenommen — kein Gate verlangte den Retarget (Negativ-Probe 5: mit dem Vor-Vergabe-Stand der Spezifikation und dem retargeteten Report meldet der Lauf fünf Befunde, davon **null** aus `docs/reviews/`). Der Edit macht zugleich zwei Aussagen des Lauf-Belegs faktisch falsch: der Report paart jetzt die wörtlich zitierte Überschrift `### JSON-Diagnose (--doctor --json)` mit dem Slug `#spec-003--json-diagnose---doctor---json` und `### JSON-Ausgabe (--json)` mit `#spec-002--json-ausgabe---json` — Slugs, die aus den zitierten Überschriften nach §DC-FA-ANCH-001.a nicht entstehen können. Wer den Beleg nachrechnet (genau sein Zweck), findet ein Paar, das nicht aufgeht; die Zeile behauptet außerdem für den 19. Juni 2026 eine Adresse, die es damals nicht gab.
- **verifizierbar:** ja — Negativ-Probe 5: der auf den Vor-Commit-Stand zurückgesetzte Report lässt den Lauf bei Exit 0 / 0 Befunden.
- **klasse:** lauf-beleg-nachgezogen-statt-beventilt

### F-3 · LOW

- **kategorie:** LOW
- **quelle:** Slice-Plan §2 Schritt 4 („Ist die Invariante mit den heutigen Schlüsseln nicht präzise ausdrückbar, wird das als Grenze benannt und der neue Schlüssel als CR-Kandidat notiert") / `MR-025`
- **pfad:** `.d-check.yml:242-248` (Kommentar der neuen Regel)
- **befund:** Von den vier bevergabten Abschnitten hält nur §4 einen Konsumenten (Go-Test: Kennung vorhanden **und** innerhalb §4 eindeutig); §3, §6 und die Eindeutigkeit in §2 sind ungebunden — gemessen bleiben eine §3-Zeile ohne Kennung, eine §6-Zeile ohne Kennung und eine zweite `### SPEC-001`-Überschrift in §2 jeweils bei Exit 0 / 0 Befunden. Dass `structure` das gar nicht leisten kann, ist eine Folge von §DC-FA-STRUCT-001.a Schritt 5 (Inline-Code-Spans werden geleert, die Kennung steht in Backticks — Lauf S belegt es: eine Überschrift `### ` mit backtick-gefasster Kennung wird als kennungslos gemeldet). Der Kommentar benennt die RE2-Grenze, nicht diese; ein CR-Kandidat für den fehlenden Schlüssel ist nirgends notiert, und die §4-Eindeutigkeit endet an der Abschnitts-Grenze — eine in §3 und §4 doppelt vergebene Nummer passiert beide Konsumenten.
- **verifizierbar:** ja — Negativ-Probe 6, Läufe N/O/P: Exit 0 trotz fehlender bzw. doppelter Kennung.
- **klasse:** vergabe-ohne-konsument-in-der-haelfte-der-abschnitte

### F-4 · LOW

- **kategorie:** LOW
- **quelle:** `MR-025` („Der Spiegel ist die Stelle, nicht die Datei") / Hard Rule AGENTS §3.7 (ein Kommentar beschreibt, was da ist)
- **pfad:** `.d-check.yml:212-220` (Kopf-Kommentar des `structure:`-Schlüssels)
- **befund:** Der Kopf führt den Schlüssel als „Modul structure — Chronologie-Monotonie … die sechs eigenen chronologischen Bestandstabellen" ein und trägt die zugehörige Vorab-Messung („heutiger Bestand 0; retro … 27 Befunde"). Der Schlüssel trägt seit dem Commit sieben Regeln, die siebte ist keine Chronologie-Regel und war von jener Messung nie erfasst. Die Stelle liegt in der Datei, die der Commit ohnehin editiert — genau die Lücken-Klasse, die `MR-025` als „Spiegel ist die Stelle" beschreibt. Wer den Kopf liest, um zu beantworten, was `structure` in diesem Repo prüft, bekommt eine Antwort, die eine Regel unterschlägt.
- **verifizierbar:** nein — Kommentar; kein Gate vergleicht Kommentar-Aussagen mit der Regel-Liste.
- **klasse:** block-kopf-enumeration-nach-erweiterung-nicht-nachgezogen

### F-5 · INFO

- **kategorie:** INFO
- **quelle:** `MR-025` (Spiegel „Autoritäts-Doku: AGENTS.md, harness/README.md — Gate-Beschreibungen") / Slice-Plan §2 Schritt 6
- **pfad:** `AGENTS.md:190` (`make doc-check`-Zeile), `harness/README.md:74` (Sensors-Zeile `make doc-check`)
- **befund:** Beide Autoritäts-Stellen zählen für `doc-check` „Links, Anker, Kennungs-Linkpflicht, Referenzmatrix, Inline-Code-Pfade" auf; `structure` — und damit die im Commit scharfgeschaltete §2-Invariante — kommt dort nicht vor, ebenso wenig `spans`, `hostpaths`, `versions`. Ein Agent, der vor dem Editieren der Spezifikation nach der Regel sucht, findet sie nur im Config-Kommentar. Die Enumerations-Lücke ist **Bestand außerhalb der Range** (sie besteht seit den Modulen `spans`/`hostpaths`/`versions`/`structure`); neu ist, dass jetzt eine Invariante an der Spezifikation daran hängt.
- **verifizierbar:** nein — Prosa; `targets` prüft die Existenz der Targets, nicht die Vollständigkeit ihrer Beschreibung.
- **klasse:** modul-enumeration-der-gate-beschreibung-unvollstaendig

### F-6 · INFO

- **kategorie:** INFO
- **quelle:** Baseline `grundlagen-source-precedence.md` §ID-Schema als Klammer („Deshalb bekommt auch eine Sektion mit eigener Schlüsselspalte Kennungen — sonst hätte eine ADR … kein Ziel") / welle-80 §1 (Scope-Entscheid)
- **pfad:** `spec/spezifikation.md:2475-2597` (Schema-Tabelle unter `### SPEC-005 — .d-check.yml`)
- **befund:** Die größte Schlüsseltabelle der Datei — 122 Datenzeilen, je Zeile eine einzelne Festlegung mit Default und Constraint — bleibt ohne Kennungen, während die sieben Zeilen der §3-Defaults sie bekommen. Die Wellen-Entscheidung („§2-Schemas als Überschrift-Kennung") deckt das, aber weder Spezifikation noch Commit noch Slice halten fest, warum die eine Schlüsselspalte Kennungen trägt und die andere nicht. Eine künftige ADR, die einen einzelnen Config-Schlüssel schärft, hat als Ziel weiterhin nur den 190-Zeilen-Abschnitt `SPEC-005` — die Unschärfe, gegen die die Welle angetreten ist.
- **verifizierbar:** nein — Scope-Frage; kein Gate misst Vergabe-Vollständigkeit außerhalb §2.
- **klasse:** vergabe-endet-an-der-groessten-schluesseltabelle

## Negativ-Proben (geprüft, ohne Befund)

1. **Vollständigkeit und Lückenlosigkeit der Vergabe (Leitfrage 1).** `grep -o 'SPEC-[0-9]\{3\}' spec/spezifikation.md` liefert 66 Treffer, `sort -u` ebenfalls 66, `sort | uniq -d` ist leer; die Treffer stehen in Dokumentreihenfolge streng aufsteigend 001…066. Zeilen-Zensus je Abschnitt: §2 fünf `###`-Überschriften (SPEC-001…005), §3 sieben Datenzeilen (006–012), §4 51 (013–063), §6 drei (064–066) — Summe 66, deckungsgleich mit der Commit-Botschaft. Jede Datenzeile der drei Tabellen beginnt mit einer backtick-gefassten `SPEC`-Zelle; die einzigen Zeilen ohne Kennung sind die drei Kopfzeilen. §1 trägt ausschließlich `DC-*-.a`-Verfeinerungen (Baseline: „Ein Abschnitt, der eine einzelne Anforderung verfeinert, trägt die Verfeinerung"), §5 ist eine Prosa-Sektion ohne Tabelle. Ohne Befund; die Scope-Grenze an der §2-Schema-Tabelle → F-6.

2. **Slug-Ableitung gegen §DC-FA-ANCH-001.a (Leitfrage 2).** Die fünf neuen Slugs Schritt für Schritt nachgerechnet (Auszeichnung entfernen, kleinschreiben, Nicht-Wort-Zeichen entfernen, Leerzeichen einzeln zu `-`): `spec-001--befund`, `spec-002--json-ausgabe---json`, `spec-003--json-diagnose---doctor---json`, `spec-004--yaml-ausgabe---yaml`, `spec-005--d-checkyml`. Die je zwei Bindestriche nach der Kennung entstehen aus Leerzeichen-Gedankenstrich-Leerzeichen, die drei vor `-doctor`/`-json` aus Leerzeichen plus dem literalen `--` der Option — dieselbe Rechnung, die schon die Alt-Slugs erklärte. Alle zwölf Verweise verwenden genau diese Form; der Lauf bestätigt sie (kein `anchor-missing`). Ohne Befund.

3. **Anker-Wanderung, Zensus im ganzen Baum (Leitfrage 2).** Fixed-String-`grep -rn` über den gesamten Baum inklusive `.harness/` und Fenced-Code nach den fünf Alt-Slugs `#befund`, `#json-ausgabe---json`, `#json-diagnose---doctor---json`, `#yaml-ausgabe---yaml`, `#d-checkyml`: **ein** Treffer, die Historien-Zeile `spec/spezifikation.md:2698`, dort als Inline-Code im Satz „die §2-Anker wandern damit" — kein Link, also kein `anchors`-Gegenstand. Gegenzensus der fünf neuen Slugs: zwölf Fundstellen, verteilt exakt wie die Botschaft sagt (Spezifikation intern 5, `README.md`/`README.de.md`/`docs/user/operations.md` je 1, `docs/plan/planning/done/slice-031-yaml-ausgabe.md` 1, Review-Report 3); `spec-004--yaml-ausgabe---yaml` wird nirgends referenziert. Ohne Befund im Zensus; die Bewertung des Review-Reports → F-2.

4. **Eingefrorene Verweise / `<a id>`-Migrationsschuld (Slice §2 Schritt 1, §5 Risiko 1).** Kein `Accepted`-ADR referenziert einen der fünf Alt-Slugs (`grep` über `docs/plan/adr/` ⇒ 0 Treffer); die Abwägung „Zeilenanker gegen Baseline-Verbot" ist damit gegenstandslos und der Verzicht auf `<a id>` korrekt. Der einzige `done/`-Slice-Verweis wurde retargetet, wie es der Lifecycle-Umgang mit `done/`-Verweisen in diesem Repo vorsieht. Die Rückführungs-Bedingung des Slice §6 („mehr eingefrorene Verweise, als `ignore-refs` sauber trägt") ist nicht eingetreten. Ohne Befund.

5. **Rot-vorher / Grün-nachher, am Produkt nachgestellt (Leitfrage 3, DoD).** Baum ohne `.git` nach `<scratch>/review-114` kopiert, Image `d-check:latest` (`e8454c88`), Ausgabe je Lauf in eine Datei, Exit explizit geprüft:

   ```bash
   docker run --rm --network none -v "<scratch>/review-114":/repo:ro d-check:latest > lauf.txt 2>&1; echo "exit=$?"
   # P0 (Baum unverändert):                          exit=0  — "421 Datei(en) geprüft, 0 Befund(e)"
   # P1 (Review-Report auf den Vor-Commit-Stand):    exit=0  — "421 Datei(en) geprüft, 0 Befund(e)"
   # W  (Spezifikation auf den Vor-Commit-Stand):    exit=1  — "421 Datei(en) geprüft, 5 Befund(e)"
   #    spec/spezifikation.md:2316  section-forbidden      (genau EINER, je Abschnitt)
   #    README.md:232 / README.de.md:235 / operations.md:52 / slice-031:45  anchor-missing
   ```

   Damit sind drei Botschafts-Aussagen belegt: die Regel ist am Vor-Vergabe-Bestand rot mit **genau einem** `section-forbidden` (ein Befund je Abschnitt, nicht je Zeile), nach der Vergabe grün bei 421 Dateien / 0 Befunden — und Lauf P1 zeigt, dass der Report-Retarget für kein Gate nötig war (→ F-2). Ohne Befund an der Messung selbst.

6. **Präfix-Negation, Positiv- und Negativfälle am Produkt (Leitfrage 3).** Je Fall eine Überschrift ans Ende von §2 eingefügt, Lauf, Rückbau per Dateikopie (md5 nach jedem Rückbau identisch mit der Ausgangskopie):

   ```text
   ### Sonderfall                     exit=1  1 section-forbidden
   ### SPEC-12 — x                    exit=1  1 section-forbidden
   ### SPEC-1234 — x                  exit=1  1 section-forbidden
   ### SPEC-00 — x                    exit=1  1 section-forbidden
   ### SPEC-001x — y                  exit=1  1 section-forbidden
   ### SPEC-Regel                     exit=1  1 section-forbidden
   ### SPECIAL — Fall                 exit=1  1 section-forbidden
   ### S            /  ### SP         exit=1  1 section-forbidden
   ###              (nackt)           exit=1  1 section-forbidden
   ### SPEC-067     (ohne Titel)      exit=1  1 section-forbidden
   Ueberschrift nur aus Inline-Code   exit=1  1 section-forbidden
   Kennung in Backticks (Lauf S)      exit=1  1 section-forbidden
   ### SPEC-067 — Neu  (Sollform)     exit=0  0 Befunde   — kein Falsch-Positiv
   ###-Zeile in einem Fence (Lauf U)  exit=0  0 Befunde   — kein Falsch-Positiv
   #### Unterabschnitt ohne Kennung   exit=0  0 Befunde   — H4 ist nicht Gegenstand
   ### mit drei Leerzeichen davor (H) exit=0  0 Befunde   — siehe F-1
   ###<TAB>Ohne Kennung        (T)    exit=0  0 Befunde   — siehe F-1
   ```

   Die ausgeschriebene Präfix-Negation ist damit für alle am Zeilen-Anfang stehenden `###`-Formen korrekt und erzeugt weder am Bestand noch an der Sollform ein Falsch-Positiv; die Bereinigung nach Schritt 5 hält Fenced-Code zuverlässig heraus. Der Fall `#### …` ist bewusst außerhalb (die Zusage spricht von `###`-Sektionen). Befunde: F-1 (H/T), F-3 (Kennung in Backticks als Beleg für die Schritt-5-Grenze).

7. **Bindepunkt der Regel (Leitfrage 3).** Die Regel steht in `.d-check.yml` (Inner Loop, `make doc-check`, seit welle-79 zusätzlich im `pre-commit`-Hook), nicht in `.d-check.closure.yml`. Das ist der richtige Ort und deckt sich mit der Begründung, die der bestehende `structure`-Block schon trägt („Inner-Loop-Bindepunkt (doc-check), weil die Zeilen in Feat-/Release-Prep-Commits entstehen"): `spec/spezifikation.md` ist ein lebendes Dokument, die kennungslose Sektion entsteht im Feat-Commit, nicht am Ruheort. Das Closure-Profil führt ausschließlich Regeln über `docs/plan/planning/done/`. Regel-Identität (`files`-Glob plus Abschnitts-Selektor) unterscheidet sich von der bestehenden §7-Regel derselben Datei — kein Exit 2 durch Identitäts-Kollision, im Lauf P0 bestätigt. Ohne Befund.

8. **Der Test bindet, was er behauptet (Leitfrage 4).** Parser-Logik von `grundCodesAusSpezifikation` (`internal/hexagon/core/app/diagnose_test.go:71-116`) gegen die reale Tabelle nachgestellt: Überschrift eindeutig gefunden (Zeile 2619), 51 Kennungen SPEC-013…SPEC-063 lückenlos, 51 Codes, keine Dublette; `AllReasons()` trägt 51 Einträge, alle verschieden — der Lockstep ist deckungsgleich. Die Kopfzeilen-Erkennung ist mit dem Spaltenwechsel korrekt von `| Code |` auf `| Kennung |` gezogen; die Trennzeilen-Regex trifft die um eine Spalte verlängerte Trennzeile weiterhin. Der fail-closed-Charakter bleibt vollständig: jede Body-Zeile, die die neue Zwei-Gruppen-Regex nicht trifft, ist `t.Fatalf` — auch die Kopfzeile selbst, falls die Spalte je wieder verschwindet. Die Eindeutigkeits-Prüfung ist wirksam (Map-Test vor der Aufnahme, Fatal beim zweiten Vorkommen), aber auf §4 skopiert → F-3. Ohne Befund an der Bindung.

9. **Kommentar-Klassen am geänderten Test-Kommentar (Hard Rule §3.7, Leitfrage 4).** `internal/hexagon/core/app/diagnose_test.go:59-70` trägt Zusage (fail-closed-Aufzählung), Abgrenzung („Annahme: alle Tabellenzeilen unter §4 sind Grund-Codes"), Grenze (die künftige zweite Tabelle macht laut rot), Kopplung („dieselbe Kopplung bindet damit auch die Vergabe an die Tabelle") und Herkunft als **ein** auflösbares Feld (`MR-000`). Die beiden Befund-Marker der Bestandsfassung (`R1-LOW-1`, `R1-INFO-1`) sind geräumt — die Bestandsgrenze aus §3.7 ist an der angefassten Zeile korrekt eingelöst; keine Slice-Nummer, kein Mess-Label, keine Deliberation über Verworfenes. Der neue Config-Kommentar `.d-check.yml:242-248` trägt Zusage, Grenze (RE2 ohne Lookahead) und Rang-Zeiger (`MR-000`) plus Herkunft als ein Feld („seit welle-80"), in derselben Form, die der Vorgänger-Report für `.d-check.yml:126-133` bereits geprüft hat. Ohne Befund; inhaltliche Lücken der beiden Kommentare → F-3/F-4.

10. **Hard Rules §3.4/§3.5 (Leitfrage 5).** Die Spezifikation nennt im gesamten Diff keine ADR, keine Welle, keinen Slice, keinen Commit-Hash und kein Closure-Datum; die neue §7-Zeile verweist auf `harness/conventions.md#mr-000--baseline-aussage`. Das ist **kein** Abwärtsverweis im Sinne §3.4: `harness/conventions.md` ist der Konventionsspeicher (Form-Autorität), keines der dort verbotenen Planungs-/Entscheidungs-Artefakte, und `matrix` kennt für ihn keine Klasse — der Lauf P0 meldet nichts. Präzedenz besteht mehrfach und in beiden Richtungen: `spec/lastenheft.md` verweist in seiner Historie auf `MR-007`/`MR-013`/`MR-017`/`MR-019`/`MR-022`, und `spec/spezifikation.md` selbst trug den Verweis auf `MR-006` schon seit dem 10. Juni 2026. Keine ADR wurde inhaltlich verändert (§3.5); ADR-0012 blieb unangetastet, die Messung zu seinem `Schärft:`-Ziel steht in der Botschaft, nicht in der ADR. Ohne Befund.

11. **MR-025-Spiegel, selbst abgeleitet (Leitfrage 6).** Ableiter statt Gedächtnis, drei `grep`-Achsen über den ganzen Baum: (a) wer die Spezifikation **maschinell parst** — genau eine Stelle, `diagnose_test.go` (angefasst); alle übrigen Nennungen von `spezifikation.md` in Go-Dateien sind Rang-Zeiger auf `§2`/`§3`/`§4`/`§DC-*-.a` in Kommentaren, die auf Abschnitts-Ebene zeigen und von der Spaltenlage unberührt bleiben. (b) Wer die **Abschnitts-Titel** nennt — `anchors_test.go:20` (Slug-Fixture für „Grund- und Fehler-Codes", Titel unverändert) und zwei Review-Reports mit §-Ankern (`#4-grund--und-fehler-codes`, `#2-datenstrukturen-und-schemas`), die alle unverändert auflösen. (c) Wer auf **§2-Anker** zeigt — die zwölf Fundstellen aus Probe 3. Nicht berührte Spiegel der `MR-025`-Tabelle: Lastenheft (Struktur-IDs sind Form, nicht Anforderung — Slice §3), Klartexte `reasonTexts()` (Code-Spalte unverändert), emittierte Vorlage `--print-config`/`--suggest-config` (welle-80 §6 schließt Produkt-Änderung aus, Diff berührt sie nicht), Benutzerhandbuch (dokumentiert das Modul generisch, spiegelt die Regel-Liste dieses Repos nicht), ADR (keine neue). Der Datei-Zensus ist damit vollständig; die Stellen-Lücken in ohnehin bearbeiteten Dateien → F-4, die Autoritäts-Doku → F-5. Ohne Befund im Datei-Zensus.

12. **Slice §5 Risiko 3 — Tabellen-Spalte vorn (Leitfrage 7).** Alle sechs bestehenden `structure`-Chronologie-Regeln geprüft: zwei zielen auf `## 7. Historie` (Lastenheft und Spezifikation — beide Tabellen haben **keine** neue Spalte bekommen, Schlüsselspalte bleibt 1), zwei auf die Roadmap (andere Datei; die zweite mit `table-column: 2` auf das Abschluss-Register), eine auf `version.md` §Verlauf und eine auf das Handbuch §11 — keine davon berührt §3/§4/§6. Umgekehrt trägt keiner der drei bevergabten Abschnitte eine Chronologie-Regel. Die neue §7-Zeile (2026-08-22) steht oben und bricht die `desc`-Monotonie nicht, weil die Vorgängerzeile dasselbe Datum trägt (nicht-strikte Monotonie). Das Risiko ist damit gemessen ausgegangen: keine Regel verschiebt sich. Ohne Befund.

13. **Botschaft gegen Baum (BEO-009-Probe).** Alle zählbaren Aussagen der Commit-Botschaft nachgerechnet: 66 Kennungen mit der genannten Verteilung ✓, „421 Dateien / 0 Befunde" ✓ (Lauf P0), „vorher 1 Befund `section-forbidden`" ✓ (Lauf W), „zwölf Verweise retargetet" mit der genannten Aufteilung ✓ (Probe 3), „ADR-0012 `Schärft:`-Ziel §Kern ist eine Tabellenzeile der Architektur-Sicht ohne Anker" ✓ (`docs/plan/adr/0012-kern-paketschnitt-model-rules-app.md:13-14` zeigt auf `spec/architecture.md` ohne Fragment; „Kern" steht dort in der Komponenten-Tabelle, `spec/architecture.md:55`, nicht als Überschrift), „Befund-Marker im Kommentar geräumt" ✓, „§1 behält seine .a-Verfeinerungen, §5 ist leer" ✓. Die eine nicht selbst nachgestellte Aussage ist Gegenprobe B (`make test` rot bei einer §4-Zeile ohne Kennung) — sie folgt zwingend aus der Parser-Logik (Probe 8: jede nicht passende Body-Zeile ist `t.Fatalf`), ein Test-Lauf im echten Repo war wegen des parallelen Gate-Laufs ausgeschlossen. Ohne Befund.

14. **Spalten-Vokabular gegen die Baseline-Vorlage.** Die Vorlage `spezifikation.template.md` betitelt die Kennungs-Spalte mit `ID`, der Commit mit `Kennung`. Das ist keine Abweichung, sondern Baseline-Vokabular: die Straten-Tabelle im Regelwerk selbst führt ihre erste Spalte als `Kennung`, der Slice-Plan benennt die Wahl vorab, das Dokument ist durchgehend deutschsprachig, und der Go-Test ist auf denselben Wortlaut gezogen. Ohne Befund.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 (F-1, F-2) |
| LOW | 2 (F-3, F-4) |
| INFO | 2 (F-5, F-6) |

## Verdikt: APPROVE MIT AUFLAGEN

Die Vergabe selbst ist sauber und belegbar: 66 Kennungen, lückenlos, doppelfrei, in Dokumentreihenfolge, mit der Verteilung, die die Botschaft behauptet; kein Abschnitt übersehen, den die Baseline-Vorlage verlangt; die Anker-Wanderung ist vollständig — kein Alt-Slug steht mehr im Baum, alle zwölf neuen Adressen sind nach §DC-FA-ANCH-001.a korrekt gerechnet und lösen im Produkt-Lauf auf; kein eingefrorener Verweis war betroffen, die `<a id>`-Abwägung entfällt zu Recht. Die Rot-vorher/Grün-nachher-Messung ist nachgestellt und stimmt auf den Befund genau. Der Test bindet, was er behauptet: 51 ↔ 51 im Lockstep, Kopf- und Trennzeilen-Erkennung korrekt nachgezogen, fail-closed in jedem Zweig erhalten, Eindeutigkeit wirksam. Hard Rules §3.4/§3.5/§3.7 sind gewahrt, das Risiko der vorangestellten Spalte ist gemessen ausgegangen, der Bindepunkt ist der richtige.

**Zwei Punkte gehören vor die Closure**, beide klein und beide außerhalb der eigentlichen Vergabe: **F-1** ist ein stiller Grün-Pfad im Gate — der Wächter-Regex ist enger als die Heading-Lexik des Moduls, das er bedient, und die Gegenprobe zeigt eine kennungslose, per Anker erreichbare §2-Sektion, die grün durchgeht; die Zusage im Kommentar steht in ihrer vollen Form da. **F-2** betrifft keinen Vertrag, sondern die Ehrlichkeit eines Lauf-Belegs: der Retarget war für kein Gate nötig (Lauf P1, Exit 0) und hinterlässt im Report zwei Überschrift-zu-Slug-Paare, die nicht aufgehen — die Klasse, für die dieses Repo das `ignore-refs`-Ventil eingerichtet und im Config-Kommentar begründet hat. Beide blockieren den Merge nicht im Sinne eines Defekts an der Lieferung, aber die Welle begründet ihren Wert ausdrücklich mit dem Gate-Konsumenten; ein Konsument mit belegtem Loch und ein verfälschter Beleg sollten die Closure nicht passieren.

**Nice-to-fix mit der Closure:** F-3 (Grenze und CR-Kandidat für §3/§6 benennen — der Slice-Plan verlangt genau das in §2 Schritt 4) und F-4 (Kopf-Kommentar des `structure:`-Blocks). F-5 und F-6 sind Notizen für den nächsten Release-Prep bzw. für die Closure-Notiz der Welle.
