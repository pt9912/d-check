# Verifikation slice-179 — DoD, Spec-Konformität und Mutationsproben

**Gegenstand:** `docs/plan/planning/in-progress/slice-179-strukture-teilmenge.md`, Range `a7e1cb4~1..78a1783` (6 Commits, HEAD == origin/main == `78a1783`).
**Datum:** 2026-08-30. **Verifikation:** unabhängiger Subagent, eigener Kontext.
**Aufbau:** zwei Images. `d-check:latest` habe ich per `make build` neu gebaut (identische Image-ID `f52e98e402b2` ⇒ das Tag steht auf HEAD). Das Vorher-Image habe ich **selbst** aus `git archive e390c2b` gebaut — Ergebnis-ID `bc43d82b63c1`, **identisch** mit dem vorgegebenen `d-check:vor-0075`. Der Vergleichspunkt ist damit belegt, nicht übernommen. Mutationsproben in einer Kopie des Baums im Scratchpad, gefahren gegen ein aus dem `deps`-Stage gebautes Test-Image.
**Gefahrene make-Targets (alle lesend):** `build`, `doc-check`, `test`, `lint`, `coverage-gate`, `arch-check`, `semgrep`, `baseline-verify`, `workflow-pins`, `gate-consistency`, `planning-check`, `verify-closure-notes`, `nightly-state`. **`make gates` und `make record-gates` bewusst nicht** — `gates` zieht `record-gates` und schriebe `.harness/state/gates-passed.diffsha`.

---

## 1. DoD-Tabelle

| # | DoD-Punkt | behauptet | gemessen | Verdikt |
|---|---|---|---|---|
| 1 | Beide Schlüssel in Schema, `--print-config`, Lastenheft 0.76.0 + Historie, Spezifikation (Schritt 1/3/6, §2, Historie); RE2-Fehler ⇒ Exit 2; `tasks-ignore-pattern` ohne `max-tasks` ⇒ Exit 2, je mit Test | erfüllt | Alle fünf Flächen tragen beide Schlüssel. `'^(['` je Schlüssel ⇒ Exit 2 mit sprechender Meldung; `tasks-ignore-pattern` ohne `max-tasks` ⇒ `Exit 2: … ist ohne max-tasks wirkungslos (halbe Aktivierung)`. Drei neue Config-Fälle in `configyaml_test.go`, Mutationen M10/M11 machen sie rot | **erfüllt** |
| 2 | Default byte-identisch, 166 Befunde, `diff` leer | erfüllt | Eigenes Profil (`max-tasks: 3`, `section-pattern: '^## 4\\. Definition of Done'`, Glob `docs/plan/planning/**/slice-*.md`) gegen beide Images: **166 Befunde** beidseits (86 `section-missing` + 80 `section-oversized`), stdout **sha-gleich** (`05b3674c…`), stderr gleich. Das Profil löst **beide** betroffenen Bedingungen aus. Zusätzlich `--json` (52 044 B) und `--doctor` (45 316 B) byte-identisch — mehr als zugesagt | **erfüllt, übererfüllt** |
| 3 | Überdeckung sichtbar; Zeichenkette *„Abschnitt trägt 4 Task-Items (3 ignoriert), erlaubt sind 3"*, dieselbe Zeile in `--doctor`; „7 ignoriert" und „0 ignoriert" je als Test | erfüllt | Zeichenkette **zeichengleich** reproduziert (Probe-Repo, 4 Liefer + 3 Konstanten, `max-tasks: 3`); `--doctor` liefert sie als `Hinweis:`-Zeile; `--json` als `message`. Beide Ränder als Subtests vorhanden und von Mutation M2 rot | **erfüllt** (Rand: der „7 ignoriert"-Fall fährt `MaxTasks = -1`, ein Wert, den der Config-Rand mit Exit 2 abweist — er ist Modell-, nicht Konfigurations-Ebene. Das ist die benannte Grenze, kein Defekt) |
| 4 | Leere Menge beantwortet: `section-missing` mit Schlüssel und Zahl | erfüllt | `docs/ac.md:1 … section-missing  alle 3 passenden Abschnitte sind von exempt-section-pattern ausgenommen — die Regel liefe leer`, Exit 1. Mutation M7 macht genau einen Test rot | **erfüllt** |
| 5 | Fence-Treue **beider** Muster gemessen; Inline-Code-Hälfte | teilweise | **Verhalten korrekt, beide Hälften, end-to-end nachgefahren.** Item im Fence zählt nicht (Probe grün); Fence-Überschrift wird weder gewählt noch ausgenommen (diskriminierende Probe: `alle 1 passenden Abschnitte …`). Inline-Code: `tasks-ignore-pattern: 'make gates'` ⇒ `(0 ignoriert)`. **Aber:** der neue Test `TestExemptSectionPattern_FenceTreu` **kann nicht rot werden** — Mutation M9 (Fence-Blindheit in `FindSectionHeads`) lässt ihn grün (s. A-4) | **teilweise** |
| 6 | Umkehr-Proben je Zusage, jede von **genau einem** Test gefangen; die drei tragenden mit Vorzustand in derselben Funktion | teilweise | Die **drei Vorzustands-Assertions** existieren und beißen (bestätigt). „Genau ein Test" trifft für **3 von 12** Mutationen zu (M3, M7, M8; dazu M10/M11 am Config-Rand); vier Mutationen reddenen je 4 Tests. Eine Zusage — Fence-Treue der Abschnitts-Ausnahme — wird von **keiner** neuen Probe gefangen | **teilweise** |
| 7 | ADR-0075 begründet Verortung, Namen, Muster-Ziele, Sichtbarkeit; im Index | erfüllt | ADR vorhanden (211 Z., Status `Accepted`, `## Re-Evaluierungs-Trigger` vorhanden), Index-Zeile 85 mit Titel/Status/Datum/Bezügen | **erfüllt** — aber die ADR trägt zwei Zahlen, die die Messung nicht trägt (A-1, A-2), und sie ist immutabel |
| 8 | Handbuch führt beide Schlüssel, Beispiel verankert, zwei gemessene Fallen | erfüllt | Beide Schlüssel im `structure`-Block; Aufzählung jetzt „Acht Dinge" und trägt **acht** (Erstens…Achtens, Z. 2014–2092); beide Beispiel-Muster als YAML **gültig** und verhalten sich wie beschrieben | **erfüllt in der Sache**, trägt die unbelegte 26 (A-2) |
| 9 | `make gates` grün; unabhängiger Review; Verifikation | **offen** (im Plan ungehakt) | Neun der zehn `gates`-Glieder einzeln gefahren und grün: `baseline-verify` (51 Dateien), `workflow-pins` (604/0), `doc-check` (604/0), `lint`, `test`, `arch-check` (0), `coverage-gate` (**94,70 %** ≥ 93), `semgrep` (55 Regeln/55 Dateien, 0), `gate-consistency` (604/0), `planning-check` (604/0). `record-gates` bewusst nicht. Review-Artefakt zu slice-179 liegt **nicht** in `docs/reviews/` | **offen, wie deklariert** |
| 10 | Antwort an den Absender | erfüllt | `docs/plan/cr/2026-08-30-antwort-a-check-structure-teilmenge.md`, 164 Z., beide Optionen angenommen, drei Umstellungs-Schritte, Paritäts-Tabelle des CR (5 Zeilen) ist 1:1 auf Tests abgebildet | **erfüllt**, trägt A-1/A-2/A-5 |

---

## 2. Nachgeprüfte Zahlen

| Zahl | Fundstelle | behauptet | gemessen | |
|---|---|---|---|---|
| Korpus-Größe | slice §1 Z. 55; ADR-0075 Z. 37; `spec/lastenheft.md` Z. 3323; CR-Antwort Z. 20 | **86** DoD-Abschnitte / Slice-Pläne | **89**. 175 Slice-Dateien; 89 tragen `## 4. Definition of Done` (dreifach belegt: `grep -l` = 89; 175 − 86 `section-missing` = 89; `max-tasks: 0` meldet 89 Abschnitte) | ✗ **A-1** |
| Task-Items | überall | **444** | **444** (Werkzeug und unabhängige `awk`-Extraktion) | ✓ |
| `section-oversized` bei `max-tasks: 3` | ADR Z. 38, CR-Antwort Z. 20, Commit-Body | **80** | **80** | ✓ |
| Byte-Identitäts-Lauf | DoD, CR-Antwort Z. 24 | **166**, `diff` leer | **166**, stdout-sha identisch, auch `--json`/`--doctor` | ✓ |
| `make gates` | CR-Antwort Tabelle | **0** von 444 | **0** | ✓ |
| `Closure-Notiz` | CR-Antwort Tabelle | **12** | **12** | ✓ |
| `Beobachtungs-Register` | CR-Antwort Tabelle | **1** | **1** | ✓ |
| `Risiko aus` | CR-Antwort Tabelle | **0** | **0** | ✓ |
| freies Muster | slice §1 Z. 58, ADR Z. 58, LH Z. 2542, Handbuch Z. 2103 | **13**, davon **2** falsch | **13**; die beiden falschen sind **wörtlich** die genannten (slice-Zeile *„§2/§3/§4/§6 tragen `SPEC-NNN` fortlaufend; Zählung in der Closure-Notiz"* und *„ADR-0012-§Kern-Messung in der Closure-Notiz"*) | ✓ |
| verankertes Muster | slice §1 Z. 61, ADR Z. 59 + 100, LH Z. 2542, CR-Antwort Z. 74, Handbuch Z. 2105 | **26** Treffer, **0** falsch | **nicht reproduzierbar.** Kein Dokument nennt das verankerte Muster. Die verankerte Form **derselben** Alternativen (`^(make gates\\|Closure-Notiz\\|Beobachtungs-Register\\|Risiko aus)`) trifft **1** Item — und dieses eine ist ein **Liefer**-Punkt (`slice-126`, *„Beobachtungs-Register trägt die Klasse als **BEO-011**, Zähler **3** …"*), also nach der eigenen Definition ein **falscher** Treffer. Sieben weitere plausible Anker-Formen gemessen (0/1/2/42/47/64/82) — keine ergibt 26. Anker kann Treffer nur **verringern**; 26 > 13 ist mit „dieselben Alternativen, verankert" unvereinbar | ✗ **A-2** |
| Anlass-Lauf | slice §1 Z. 38, §3 Z. 117, §7 Z. 178–179 | **160** Befunde bei **223** Dateien | **166** Befunde; 224 Markdown-Dateien unter `docs/plan/planning` (223 galt bei `df78907`, vor Anlage dieses Plans). Widerspricht §4 **derselben Datei** (166) | ✗ **A-3** |
| „in **allen** 86 Abschnitten steht `make gates`" | CR-Antwort Z. 47, slice §1 Z. 75, ADR Z. 179 | alle | **86 von 89**. Drei DoD-Abschnitte tragen die Wendung gar nicht: `slice-096`, `slice-111`, `slice-125` | ✗ **A-5** |
| „null Fundstellen ohne Backticks" | CR-Antwort Z. 47 | 0 | **0** — 424 Vorkommen von `make gates` im Slice-Bestand, **alle** 424 in Backticks | ✓ |
| „jüngste Slices 5 bis 9 DoD-Items" | slice §1 Z. 39 | 5–9 | **5–9** für slice-150…178; slice-179 selbst trägt 10 | ✓ |
| Handbuch-Aufzählung | Commit `c4ae0b6` | „hieß Sechs, trug sieben, jetzt acht" | **acht** (Erstens…Achtens) | ✓ |

---

## 3. Mutationsproben

Alle Mutationen einzeln in einer Baumkopie, `go test ./internal/...` je Lauf. Basislauf: **grün**.

| # | Mutation (Datei) | rot geworden | welche |
|---|---|---|---|
| M1 | `taskIgnoreRE` liefert immer `nil` (structure.go:314) | **4** | `KonstantenZaehlenNichtMit`, `LiestDenBereinigtenText`, `MeldungNenntIgnorierte/zu breit`, `SiehtDenItemText/am Item-Text` |
| M2 | Sichtbarkeits-Zusatz entfernt (`zusatz = ""`, :221) | **2** | `LiestDenBereinigtenText`, `MeldungNenntIgnorierte` (**beide** Subtests) |
| M3 | Zusatz auch **ohne** Schlüssel (`if true`, :220) | **1** | `AbwesendIstByteIdentisch` — die Byte-Identitäts-Zusage beißt |
| M4 | Muster sieht die **rohe** Zeile statt des Item-Textes (:299) | **4** | wie M1 + `SiehtDenItemText/am Listen-Marker` |
| M5 | Abschnitts-Ausnahme wirkungslos (:135) | **4** | `NimmtDenBestandHeraus`, `LeereMengeMeldet`, `LaeuftVorDerKardinalitaet`, `SiehtDieRoheUeberschriftenZeile/mit #-Folge` |
| M6 | Abschnitts-Muster sieht die Überschrift **ohne** `#`-Folge (:96) | **4** | wie M5, dazu `…/ohne #-Folge` |
| M7 | Nullmengen-Härte entfernt (früh-Return tot, :139) | **1** | `LeereMengeMeldet` |
| M8 | Kardinalitäts-Prüfung **vor** die Ausnahme gezogen | **1** | `LaeuftVorDerKardinalitaet` |
| **M9** | **`FindSectionHeads` fence-blind** (sections.go:27) | **1 — und keiner der neuen** | nur `TestClosureHeadingImFenceZaehltNicht` (Altbestand). **`TestExemptSectionPattern_FenceTreu` bleibt grün.** Gegenprobe: eine um **eine Zeile** geänderte Variante (Ausnahme-Muster auf die **echte** statt auf die **gefencte** Überschrift) ist auf `base` grün und auf M9 rot — die diskriminierende Form existiert und ist billig |
| M10 | RE2-Prüfung beider Schlüssel entfernt (configyaml.go:307–308) | **1** | `TestDecode_StructureFehler` |
| M11 | Halbe-Aktivierungs-Prüfung entfernt (configyaml.go:365) | **1** | `TestDecode_StructureFehler` |
| M12 | `ignoriert` wird nie erhöht (:300) | **2** | `LiestDenBereinigtenText`, `MeldungNenntIgnorierte/zu breit` |
| M13 | `proseLines` fence-blind (markdown.go:94) — Kontrollprobe für die **Item**-Hälfte | **36** | darunter `LiestDenBereinigtenText`. Die Item-Fence-Zusage ist also gedeckt |

**Ergebnis:** 12 von 13 Mutationen werden gefangen. Die **eine** ungefangene ist genau die Zusage aus DoD-Punkt 5 („die Überschrift im Fence wird weder gewählt noch ausgenommen"). Das Verhalten ist korrekt — die **Probe** ist es nicht.

---

## 4. Spec-Konformität

Alle vier neuen Akzeptanzkriterien aus `spec/lastenheft.md` 0.76.0 einzeln als Probe gefahren (Probe-Repos, `--enable structure`):

- **Teilmenge (`tasks-ignore-pattern`)** — 7 Items, 4 getroffen, `max-tasks: 3` ⇒ **0 Befunde**; ohne Schlüssel ⇒ `Abschnitt trägt 7 Task-Items, erlaubt sind 3`, **byte-identisch** gegen `d-check:vor-0075`; fünfter, nicht getroffener Punkt ⇒ `Abschnitt trägt 4 Task-Items (4 ignoriert), erlaubt sind 3`; Muster ohne Treffer ⇒ `… 8 Task-Items (0 ignoriert) …`. **erfüllt.**
- **Teilmenge (Muster-Ziel)** — `^\\*\\*Konstante:` trifft; `^- \\[ \\] \\*\\*Konstante:` trifft **nicht** (`0 ignoriert`); Item mit `` `make gates` `` gegen Muster `make gates` ⇒ `0 ignoriert`. **erfüllt.**
- **Teilmenge (`exempt-section-pattern`)** — drei gleichartige `### AC-`-Abschnitte, zwei getroffen, `sections: each`, Marke fehlt überall ⇒ **genau 1** Befund (Z. 11); ohne Schlüssel **3**; Muster ohne `#`-Folge (`^AC-00`) ⇒ **3** (trifft keinen). **erfüllt.**
- **Teilmenge (Nullmenge und Kardinalität)** — Muster über alle ⇒ `docs/ac.md:1 … section-missing  alle 3 passenden Abschnitte sind von exempt-section-pattern ausgenommen …`; `sections: one` mit zwei Treffern ⇒ Vorzustand `section-ambiguous`, nach Abzug eines ⇒ **0 Befunde**, kein `section-ambiguous`. **erfüllt.**

**`spec/spezifikation.md` §`DC-FA-STRUCT-001.a`** — Schritt 1 (zwei neue Exit-2-Ränder), Schritt 3 (Abzug gegen die **rohe getrimmte** Überschriften-Zeile, **vor** Schritt 4, Nullmengen-Härte mit Schlüssel + Zahl) und Schritt 6 (Vergleichsgegenstand = getrimmter Item-Text hinter dem Marker; Zusatz nur bei gesetztem Schlüssel, auch bei null) decken sich **wörtlich** mit `structure.go:127–142`, `:209–226`, `:290–315`. §2-Schema führt beide Zeilen. §4 Grund-Codes unberührt — korrekt, es gibt keinen neuen Code.

**Deklarations-Vollständigkeit (C):** beide Schlüssel erscheinen in `--print-config`-Gerüst (Z. 178, 196 der Ausgabe), Handbuch, Lastenheft (§Schema-Tabelle + Historie), Spezifikation (Schritt 1/3/6 + §2 + Historie); ADR-0075 ist im Index. **Keine Fläche fehlt.** `README.md`/`README.de.md` beschreiben nur die *Bedingungen* (neun) — die zwei Schlüssel sind keine, `hint` aus slice-177 steht dort ebenfalls nicht; `docs/user/operations.md` und `benutzerhandbuch-standard.md` führen keine `structure`-Schlüssel. **`CHANGELOG.md` trägt nichts** — das entspricht der Repo-Praxis (jüngster Eintrag 0.67.0/slice-170; auch slice-177 fehlt dort noch), gehört also in die Release-Prep, nicht in diesen Slice.

---

## 5. Abweichungen Zusage ↔ Zustand

**A-1 — Die Korpus-Zahl ist falsch: 86 statt 89.**
`docs/plan/adr/0075-erklaerte-teilmenge-in-structure.md:37`, `spec/lastenheft.md:3323`, `docs/plan/cr/2026-08-30-antwort-a-check-structure-teilmenge.md:20`, `docs/plan/planning/in-progress/slice-179-strukture-teilmenge.md:55` (auch ADR:158, :179 und slice:75). Behauptet: „die **86** Slice-Pläne dieses Repos, die einen `## 4. Definition of Done` führen". Gemessen: **89**. Die 86 ist die Zahl der `section-missing`-Befunde (Dateien **ohne** die Sektion) — und zufällig auch die Zahl der Abschnitte, die `make gates` tragen. Die davon abhängigen Zahlen (444, 80, 166) sind korrekt; nur die Bezugsmenge nicht. **ADR-0075 ist `Accepted` und damit immutabel** (§3.5) — die Korrektur geht über einen `## Geschichte`-Anhang oder eine Folge-ADR, nicht über eine Kern-Änderung.

**A-2 — Die Zahl „26 Treffer, 0 falsche" ist nicht belegt und nicht reproduzierbar.**
`…/slice-179-…:61`, `…/0075-…:59` und `:100`, `spec/lastenheft.md:2542`, `…/antwort-…:74`, `docs/user/benutzerhandbuch.md:2105`. **Kein Dokument nennt das verankerte Muster**, obwohl es die tragende Empfehlung an den Absender ist („verankert das Muster"). Die verankerte Form derselben Alternativen misst **1** Treffer, nicht 26 — und dieser eine Treffer ist ein Liefer-Punkt aus `slice-126`, wäre also selbst ein „falscher". Da Anker die Treffermenge nur verkleinern kann, ist 26 > 13 unter der Tabellen-Überschrift *„An denselben 444 Items"* nicht möglich. Der Vorgänger-Stand dieser Messung (`24/23` über 129 Items) steht noch im Scratchpad der Hauptsitzung — die aktuelle Zeile sieht aus wie ihre Fortschreibung, ist aber eine andere Grundmenge. **Das ist die Gestalt aus `BEO-020`** („gemessen wird die eigene Menge, ausgesagt wird über die fremde") an der Zahl, die der Absender in seine eigene Konfiguration übernehmen soll.

**A-3 — Der Slice-Plan widerspricht sich selbst: 160/223 gegen 166.**
`…/slice-179-…:38` („**160 Befunde bei 223 Dateien**"), :117 („Die 160 Befunde"), :178–179 („160 von 223") gegen :139 („liefert **166** Befunde"). Gemessen an HEAD: **166**. Die 223 war der Markdown-Zählstand unter `docs/plan/planning` bei `df78907`, vor Anlage dieses Plans; sie ist außerdem der **Scan**-Dateizähler, nicht die Kandidatenmenge der Regel (175 Slice-Dateien). Die CR-Antwort zieht an der Parallelstelle (Z. 156 „Die **80** Befunde sind der Anlass") eine andere Zahl als der Slice (§3: „Die 160 Befunde").

**A-4 — `TestExemptSectionPattern_FenceTreu` kann für die Zusage, die er trägt, nicht rot werden.**
`internal/hexagon/core/rules/structure_teilmenge_test.go` (letzte Funktion). Sein Ausnahme-Muster `^### AC-001\\b` trifft **genau die gefencte** Überschrift — ob sie nie gewählt oder gewählt-und-dann-ausgenommen wurde, ist am Ergebnis nicht unterscheidbar. Mutation M9 (Fence-Blindheit) lässt ihn **grün**. Der Kommentar über ihm sagt ausdrücklich *„eine geerbte Zusage ohne Probe ist eine Erinnerung"* — das ist **BEO-023** wortwörtlich, an dem Test, der gegen BEO-023 gebaut wurde. Folgeaussagen, die dadurch zu weit reichen: CR-Antwort Z. 148 *„Fence-Treue für beide Muster — geerbt, nicht nachgebaut, und **mit je einer Probe belegt**"* und ADR-0075 §Fitness Function *„Sie halten **jede** Zusage dieser ADR einzeln"*. **Fix ist eine Zeile:** Ausnahme-Muster auf die **echte** Überschrift richten (`^### AC-042\\b`) — dann erwartet der Test `section-missing  alle 1 passenden Abschnitte …`, ist auf HEAD grün und auf M9 rot (beides gemessen).

**A-5 — „In **allen** 86 Abschnitten steht die Wendung" ist als Aussage über den Korpus falsch.**
`…/antwort-…:47`, `…/slice-179-…:75`, `…/0075-…:179`. Drei der 89 DoD-Abschnitte tragen `make gates` gar nicht (`slice-096`, `slice-111`, `slice-125`). Der load-bearing Teil („null von 444 Items, weil durchgängig in Backticks") ist unberührt und **korrekt**.

**A-6 — „jede Umkehr-Probe von genau einem Test gefangen" trifft die Messung nicht.**
`…/slice-179-…:153` und §2 Punkt 7 (:107). Gemessen: 3 von 12 Mutationen färben genau einen Test rot; vier färben je vier. Das ist kein Mangel an Deckung — aber die Zusage sagt etwas anderes, als die Suite tut, und §5 Anstrich „Eine Commit-Botschaft oder Closure-Notiz behauptet nicht mehr, als die Arbeit trägt" (AGENTS.md §5, `BEO-009`) gilt auch für den DoD-Haken.

**A-7 (INFO) — Der Sichtbarkeits-Rand „7 ignoriert" ist über Konfiguration nicht erreichbar.**
`structure_teilmenge_test.go`, `TestTasksIgnorePattern_MeldungNenntIgnorierte` fährt `MaxTasks = ptr(-1)`; der Config-Rand weist `max-tasks < 0` mit Exit 2 ab. Der Test misst die Meldungs-Form korrekt, aber der Fall „Muster nimmt alles" ist end-to-end **stumm** — genau die im Plan, in der ADR, im Lastenheft und im Handbuch benannte Grenze. Kein Widerspruch, nur der Hinweis, dass „als Test gefangen" hier nicht „als Verhalten erreichbar" heißt.

---

## 6. Was gegenüber der Zusage **besser** ist

- **Die Byte-Identität ist stärker belegt als behauptet.** Zugesagt war der Plain-Lauf; identisch sind auch `--json` (52 044 B) und `--doctor` (45 316 B).
- **Der Vergleichspunkt ist belastbar.** Mein unabhängig aus `e390c2b` gebautes Image hat dieselbe Image-ID wie `d-check:vor-0075` — das Vorher-Image ist kein Vertrauensvorschuss.
- **Die Inline-Code-Falle ist allgemeiner als gemessen.** Die Baseline-Vorlage `slice.template.md` schreibt selbst `` - [ ] `make gates` grün. `` — das Beispiel-Muster des Absenders ist damit in **jedem** Repo inert, das die Vorlage benutzt, nicht nur in diesem. Das stützt die Empfehlung stärker, als die Dokumente behaupten.
- **Beide dokumentierten Beispiel-Muster sind gültige Konfiguration und verhalten sich wie beschrieben** (`'^\\*\\*Gate'`, `'^## 9\\. Alt'`, `'^\\*\\*Konstante:'`, `'^### AC-0(0[1-9]|1[0-9])\\b'`) — die Falle aus dem Release-Prep-Gedächtnis (Fenced-YAML-Beispiele sind gate-frei) ist hier nicht zugeschlagen.
- **Die zwei „falschen" Treffer des freien Musters sind wörtlich die genannten** — die Diagnose in Plan, ADR und CR-Antwort ist nicht paraphrasiert, sondern belegt.
- **Die Lifecycle-Spur ist sauber.** `e390c2b` ist ein **R100**-Rename (reiner Move) und trägt den gekoppelten Roadmap-Flip (Ruhe-Marker entfernt) im selben Commit — MR-013 korrekt angewandt; `make planning-check` grün.
- **Beide `d-check:cite`-Direktiven sind korrekt angekert** (MR-054): `modul-05-planning-harness.md:213-214` trifft *„Sub-Area-Wahl prüfen. Jede Sub-Area … Schwelle ≥ 2"*, `:219-219` trifft *„Offene Beobachtungen sichten."* — beide auf der **vorschreibenden** Zeile; `citations` läuft im `doc-check` und ist grün. Der Nachtlauf-Block trägt bewusst keine Direktive, wie MR-054 verlangt.
- **Der Nachtlauf-Stand stimmt.** `make nightly-state`: beide Achsen **gruen** (`upstream-drift.yml` inzwischen 2026-08-30T06:08:17Z, `image-scan.yml` 2026-08-29T10:07:43Z — letzterer zeichengleich mit dem Plan).
- **Alle im Register zitierten Beobachtungen existieren** (`BEO-011`, `BEO-013`, `BEO-023`).

---

## 7. Repo unverändert — Beleg

```
git status --porcelain   → leer
git rev-parse HEAD       → 78a1783820bb6a709ddc74dbe0c6c578b81d6512
```

Hashes nach dem Lauf (identisch mit dem Stand vor meinen Messungen — ich habe nie in den Baum geschrieben):

```
6af0b14b…302dc  docs/plan/planning/in-progress/slice-179-strukture-teilmenge.md
3a762b39…3cad9  internal/hexagon/core/rules/structure.go
42c0c960…25c6c  spec/lastenheft.md
7eaec85e…faff8  spec/spezifikation.md
46b562dd…955d4  docs/user/benutzerhandbuch.md
f4285360…064ae  docs/plan/adr/0075-erklaerte-teilmenge-in-structure.md
ba43d5e3…7b862  docs/plan/cr/2026-08-30-antwort-a-check-structure-teilmenge.md
953bbf94…76dbc  docs/plan/planning/in-progress/roadmap.md
2276e488…22164  .harness/state/gates-passed.diffsha   (mtime 08:28, von mir NICHT berührt)
```

Scratchpad-Baum `…/scratchpad/verif179/` restlos entfernt (root-eigene Go-Cache-Dateien per Container gelöscht); die von mir angelegten Image-Tags `d-check:verif-pre` und `d-check:verifdeps` entfernt. `d-check:latest` steht auf derselben ID wie vorher.

---

## 8. Empfehlung

**Vor der Closure zu schließen:**

1. **A-2 zuerst** — die „26" steht in fünf Dokumenten, darunter der **ausgehenden** CR-Antwort, und ist die Grundlage der Empfehlung „verankert das Muster". Entweder das gemessene Muster nachliefern und die Zahl nachziehen, oder die Zeile auf das reduzieren, was trägt: das freie Muster nimmt 2 echte Zusagen mit, die verankerte Form nimmt sie nicht.
2. **A-1** — 86 → 89 in Slice, Lastenheft-Historie und CR-Antwort. **ADR-0075 nur per `## Geschichte`-Anhang oder Folge-ADR** (§3.5); ein Kern-Edit bräche `make adr-check`.
3. **A-3** — die drei 160/223-Stellen im Slice auf 166 ziehen, sonst widerspricht das Dokument sich selbst.
4. **A-4** — die eine Zeile im Fence-Test drehen. Gemessen: auf HEAD grün, auf der Fence-Mutation rot. Danach trägt die CR-Antwort-Zeile „mit je einer Probe belegt" wieder; die ADR-Zeile „jede Zusage einzeln" bleibt ohne diesen Fix zu weit.
5. **A-5** — „in allen 86 Abschnitten" → „in 86 der 89".
6. **A-6** — den DoD-Haken auf das umformulieren, was die Suite tut (jede Zusage ist von mindestens einer Probe gedeckt; die drei tragenden fahren ihre Umkehr in derselben Funktion).
7. **Der offene DoD-Punkt 9** — Review-Artefakt zu slice-179 fehlt in `docs/reviews/`; `make gates` selbst ist bis auf `record-gates` von mir vollständig grün gemessen.

**Nach der Closure / in der Release-Prep:** CHANGELOG-Eintrag zu ADR-0075 (Repo-Praxis, kein Slice-Defekt); optional die `section-missing`-Ursachenliste in der README-Modulbeschreibung.

**Zwei Beobachtungen für das Register, beide belegt:** A-4 ist ein Wiederauftreten von **`BEO-023`** (ein Wächter, der nie fangen konnte — diesmal in dem Test, der gegen BEO-023 geschrieben wurde); A-1/A-2 sind ein Wiederauftreten von **`BEO-020`** (gemessen wird eine Menge, ausgesagt wird über eine andere).
