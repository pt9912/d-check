# Review-Report — slice-179 (`structure`: erklärte Grundmenge)

**Review-Art:** Code (gegen Plan, ADR-0075, `DC-FA-STRUCT-001`, Hard Rules)
**Gegenstand:** `a7e1cb4~1..78a1783` (6 Commits), HEAD `78a1783`
**Skill:** `.harness/skills/reviewer.md` v1.13.0 · **Modell:** claude-opus-5[1m] · **Datum:** 2026-08-30
**Eingangs-Kontext:** Slice-Plan, eingehender CR + Antwort, ADR-0075/0069/0070/0073, `AGENTS.md` §3, `spec/lastenheft.md` 0.76.0, `spec/spezifikation.md`, Handbuch

## Eigener Lauf (Ausgabe, nicht behauptet)

| Target | Ergebnis |
|---|---|
| `make test` | Exit 0 — alle Pakete `ok` |
| `make lint` | Exit 0 |
| `make coverage-gate` | `Coverage 94.70% erfüllt Schwelle 93%` |
| `make doc-check` | `604 Datei(en) geprüft, 0 Befund(e)` |
| `make arch-check` | `gesamt: 0 Befund(e)` |
| `make planning-check` | `604 Datei(en) geprüft, 0 Befund(e)` |
| `make gate-consistency` | `604 Datei(en) geprüft, 0 Befund(e)` |

`make gates` habe ich **nicht** gefahren (es zieht `record-gates` und schriebe `.harness/state/gates-passed.diffsha`). Dazu: 7 Code-Mutationen in einer Repo-**Kopie** unter scratchpad, Vorher/Nachher-Vergleiche gegen `d-check:vor-0075`, ~30 Probe-Läufe des Images gegen Probe-Repos und gegen den echten Planungsbaum (Config via zweitem Read-only-Mount über `.d-check.yml`).

## Urteil: **blockiert** — 1 HIGH · 4 MEDIUM · 4 LOW · 1 INFO

---

## HIGH

### H-1 · Die tragende Messung „verankert 26 Treffer" ist nicht reproduzierbar und arithmetisch unmöglich
`docs/plan/adr/0075-…md:59,100` · `spec/lastenheft.md:2542` · `spec/spezifikation.md:3039` · `docs/user/benutzerhandbuch.md:2104` · `docs/plan/cr/2026-08-30-antwort-…md:74` · `docs/plan/planning/done/slice-179-…md:61` · Commit-Body `b9e6052`

**Zugesagt:** *„Gefahren über dieselben 444 Items: frei, wie im CR → 13 ignoriert, davon 2 falsch; verankert (das Item **beginnt** mit dem Ausdruck) → 26 ignoriert, 0 falsch."* ADR-0075 Entscheidung 3 stützt die Asymmetrie der beiden Muster ausdrücklich darauf: *„gemessen ist die verankerte Form die tragfähige (26 gegen 13 Treffer, 0 gegen 2 falsche)"*.

**Gemessen** (Modul `structure` gegen `docs/plan/planning/**/slice-*.md`, `section-pattern: '^## 4\\. Definition of Done'`, `sections: each`, `max-tasks: 0` — 89 Abschnitte, 444 Items; unabhängig nachgezählt mit einer awk-Nachbildung der Vorverarbeitung, die exakt dieselben 444/13/1 liefert):

| Musterform | ignorierte Items |
|---|---|
| frei `(make gates\\|Closure-Notiz\\|Beobachtungs-Register\\|Risiko aus)` | **13** ✓ (2 davon echte Liefer-Punkte ✓) |
| **verankert** `^(make gates\\|Closure-Notiz\\|Beobachtungs-Register\\|Risiko aus)` | **1** ✗ (behauptet: 26) |

Ein `^` vor demselben Ausdruck kann die Treffermenge **nur verkleinern** — 13 frei und 26 verankert ist für dieselbe Wendung nicht bloß falsch gemessen, sondern unmöglich. Ich habe zusätzlich ~15 alternative Lesarten von „verankert" durchprobiert (bold-Präfix `^\\*\\*` → 47, `^grün` → 42, `^(Closure|Register|Risiko|Gate)` → 1, gegen die **rohen** Zeilen → 78, die im Plan als überholt bezeichnete Vorgänger-Zahl 24/23 → 88/78); **keine** ergibt 26.

**Warum es zählt:** Die Zahl steht in einer nach §3.5 **immutablen** `Accepted`-ADR, in **zwei Spec-Straten** (Lastenheft ist vertraglich abnahmebindend), im Benutzerhandbuch als Anwender-Anleitung und in einem **ausgehenden** Dokument an einen Adopter, dem damit begründet wird, warum er sein Muster umschreiben soll. `AGENTS.md` §5 verlangt, dass eine genannte Probe gelaufen ist; hier reproduziert die eine Hälfte der Tabelle exakt und die andere gar nicht. Die *Richtung* der Empfehlung („verankern") ist plausibel und durch die 2 falschen Treffer belegt — der *Beleg* für ihre Tragfähigkeit ist es nicht.

**Kategorie-Begründung:** Der Skill führt „behauptete Probe fand nicht statt" als HIGH-nahe erste Richtung von `BEO-009`; die Unkorrigierbarkeit der ADR und die Außenwirkung heben sie hier auf HIGH.
**Verifizierbar:** ja — der Lauf oben, 5 Zeilen Konfiguration.
**Klasse:** `mess-zahl-nicht-reproduzierbar`
**Billigster Fix:** Die Tabelle mit der tatsächlichen verankerten Zahl neu fahren und die Aussage auf das stützen, was trägt (0 falsche statt 2), in allen sechs Artefakten; für die ADR ist das eine Folge-ADR mit `supersedes`.

---

## MEDIUM

### M-1 · `hint` löscht beide neuen Sichtbarkeits-Zusagen
`internal/hexagon/core/rules/structure.go:140,223` (beide via `structureFinding` → `MessageFor`)

**Zugesagt:** ADR-0075 Entscheidung 4 — *„die Meldung nennt die Zahl der ignorierten Items … auch bei null"*; Entscheidung 5 — *„die Meldung nennt den Schlüssel, der es tat"*; Handbuch: *„**`(0 ignoriert)` heißt: Ihr Muster wirkt nicht.**"* Das ist laut ADR und Antwort die **einzige** Gegenmaßnahme gegen ein zu breites Muster.

**Gemessen** (Probe-Repo, dieselbe Regel einmal mit und einmal ohne `hint`):

```
ohne hint: … section-oversized  Abschnitt trägt 4 Task-Items (3 ignoriert), erlaubt sind 3
mit  hint: … section-oversized  Hoechstens drei Liefer-Punkte
--doctor:                       Hinweis: Hoechstens drei Liefer-Punkte
```
```
ohne hint: … section-missing  alle 2 passenden Abschnitte sind von exempt-section-pattern ausgenommen — die Regel liefe leer
mit  hint: … section-missing  Jede Anforderung braucht eine Boundary-Marke
```

**Warum es zählt:** ADR-0073 nimmt genau diese Klasse — *„der unlesbare Dateibaum und **die leer laufende Regel** … Sie verletzen keine Bedingung … Sie behalten ihre eigene Meldung"* — vom `hint` **aus**, und der Code führt diese Ausnahme über `structureRawFinding` bzw. Literal-Findings auch aus (`structure.go:73`, `:121`). Der neue Nullmengen-Befund sagt in seinem eigenen Text *„die Regel liefe leer"*, geht aber durch `structureFinding`. Ergebnis: wer `hint` (ausgeliefert eine Slice zuvor) und `tasks-ignore-pattern` kombiniert — die naheliegende Kombination —, verliert die Diagnose, mit der er das Muster prüfen soll; kein Test deckt `hint` × einen der beiden neuen Schlüssel.
**Verifizierbar:** ja — zwei Probe-Läufe. **Klasse:** `zusage-gilt-nicht-auf-allen-pfaden`

### M-2 · Lastenheft sagt „inline-code-treu" für beide Muster zu; `exempt-section-pattern` sieht Inline-Code
`spec/lastenheft.md:2564-2567` · `internal/hexagon/core/rules/structure.go:96`

**Zugesagt:** *„Und **beide Muster lesen fence- und inline-code-treu** wie die Bedingungen, die sie verkleinern"*. Die Antwort an den Absender führt unter **„Angenommen, unverändert"** (`…antwort…md:128`) *„Fence-Treue für beide Muster — geerbt, nicht nachgebaut, und mit je einer Probe belegt"* — gegen einen CR, der wörtlich *„Beide Muster dürfen Code-Blöcke und **Inline-Code** nicht sehen"* verlangt.

**Gemessen:** `structureExemptSections` matcht `strings.TrimSpace(lines[h.Line-1])`, also die **rohe** Zeile. Probe mit Überschrift ``### `AC-001` alt``:
- `exempt-section-pattern: '^### `AC-001`'` (mit Backticks) → **greift** (Abschnitt ausgenommen).
- `exempt-section-pattern: '^###  *AC-001'` (ohne Backticks) → **greift nicht**.

Fence-Treue hält (`TestExemptSectionPattern_FenceTreu`, nachgefahren ✓), Inline-Code-Treue existiert für dieses Muster **nicht**.

**Warum es zählt:** Ein Spec-Stratum sagt eine Eigenschaft zu, die der Code nicht hat, und die Antwort meldet dem Absender eine seiner Anforderungen als „unverändert angenommen", obwohl nur die halbe zutrifft. Nebenbefund, der die Klasse erhärtet: der **ursprüngliche** DoD-Punkt lautete *„ein Task-Item bzw. eine **Überschrift** im Fenced-Block bzw. **in Inline-Code** wird nicht getroffen"* und wurde im Schluss-Commit `78a1783` genau um die Überschrift-/Inline-Code-Hälfte gekürzt — also dort, wo er gescheitert wäre.
**Verifizierbar:** ja — zwei Probe-Läufe. **Klasse:** `zusage-weiter-als-code`

### M-3 · `exempt-section-pattern` fehlt in `Identity()` — die naheliegende Grandfathering-Konfiguration ist Exit 2
`internal/hexagon/core/model/config.go:549-566` (unverändert) · `internal/adapter/driven/configyaml/configyaml.go:288`

**Gemessen:** Zwei Regeln über dieselbe Datei-Glob und denselben Selektor, die sich **nur** durch `exempt-section-pattern` unterscheiden:

```
d-check: error: .d-check.yml: structure[1]: Regel-Identität "docs/*.md :: ^### " kommt doppelt vor
exit=2
```

**Warum es zählt:** Das ist die Form, die Grandfathering eigentlich braucht — *Bedingung A für alle Abschnitte, Bedingung B nur für die nicht-ausgenommenen*. Es gibt keinen Ausweg: `files` und Selektor sind identisch, und beide Bedingungen in **eine** Regel zu legen unterwirft auch A der Ausnahme. ADR-0069 hat für genau diese Frage die Gegenentscheidung getroffen und `table.order-column` in `Identity()` aufgenommen (*„zwei Chronologie-Zusagen ueber denselben Abschnitt sind verschiedene Zusagen und brauchen verschiedene Identitaeten"*, Kommentar `config.go:554-561`); `exempt-section-pattern` macht denselben Unterschied und wurde nicht aufgenommen. Weder ADR-0075 noch Lastenheft, Spezifikation oder Handbuch nennen die Einschränkung, und die Fehlermeldung zeigt nicht auf ihre Ursache.
**Verifizierbar:** ja — die Config oben. **Klasse:** `identitaet-traegt-die-unterscheidung-nicht`

### M-4 · Der „gemessene Defekt des Antrags" beim Abschnitts-Muster ist eine Fehllesung des CR
`docs/plan/adr/0075-…md:93` · `docs/plan/cr/2026-08-30-antwort-…md:105-107` · `internal/hexagon/core/model/config.go:487-489` · Commit-Body `b9e6052`

**Zugesagt/behauptet:** ADR: *„die Falle, in die der Antrag selbst gelaufen ist — sein Beispiel (`'^AC-…'`) trifft am realen Lastenheft **nichts**, weil dort `### AC-…` steht"*. Antwort: *„Das war der **dritte gemessene Punkt**"*. Commit: *„ZWEI DEFEKTE DES ANTRAGS SIND GEMESSEN UND BEHOBEN … genau das tut das Antrags-Beispiel."*

**Gemessen am Antragstext** (`…cr-a-check-structure-teilmenge.md:61`): der Absender deklariert seinen Vergleichsgegenstand ausdrücklich — *„Abschnitte, deren **UEBERSCHRIFTSTEXT** dieses RE2 trifft"*. Sein Muster `'^AC-[A-Z]+-0…'` ist gegen **diese** Semantik korrekt geschrieben. Dass es „nichts trifft", entsteht erst durch die hier getroffene Entscheidung, gegen die **rohe** Zeile zu vergleichen. Das ist eine legitime und gut begründete Entscheidung — aber kein Defekt des Antrags und nichts, was am Antrag *gemessen* wurde. Zusätzlich: „am realen Lastenheft" bezeichnet a-checks Lastenheft, das in diesem Repo nicht liegt und hier nicht messbar ist; die Aussage ist eine Vermutung im Gewand einer Messung.

**Warum es zählt:** Skill-Frage 9 (Quelle über ihren Geltungsbereich hinaus zitiert — Feld gelesen statt Titel): der CR trägt die zugeschriebene Aussage nicht. Die Zuschreibung steht in einer immutablen ADR, in einem Code-Kommentar und in einem Dokument, das an den Absender geht.
**Verifizierbar:** ja — Zeile 61 des CR gegen ADR-Zeile 93. **Klasse:** `fremd-quelle-fehlgelesen`

---

## LOW

### L-1 · Code-Kommentar nennt „129 reale Items" — eine Zahl, die in keinem Artefakt vorkommt
`internal/hexagon/core/rules/structure.go:289`: *„ein freies Substring-Muster nahm an **129** realen Items eine echte Zusage still mit (ADR-0075)"*. Überall sonst — Slice, ADR, Lastenheft, Spezifikation, Handbuch, Antwort, Commit — lautet die Grundmenge **444**. `grep` findet 129 in keinem der Artefakte; die Item-Zahlen je Lifecycle-Verzeichnis sind 419 (`done/`) + 15 (`open/`) + 10 (`in-progress/`) = 444. Der Kommentar altert damit ab dem ersten Tag gegen jede Deklarations-Fläche.

### L-2 · Das Zahlen-Tripel 86 / 80 / 444 kommt in keinem einzelnen Lauf zusammen vor
`docs/plan/adr/0075-…md` §Kontext · `…antwort…md:20-24` · Commit-Body. Gemessen:

| Selektor | Abschnitte | Items | `section-oversized` bei `max-tasks: 3` |
|---|---|---|---|
| `section: '## 4. Definition of Done'` (exakt) | **86** | **425** | **77** |
| `section-pattern: '^## 4\\. Definition of Done'`, `sections: each` | **89** | **444** | **80** |

Die 444 und die 80 gehören zu **89** Abschnitten (86 exakt + 3 × `## 4. Definition of Done (vorläufig)` in `slice-036/037/038`); über die 86 sind es 425/77. Ebenso: *„in **allen** 86 Abschnitten steht die Wendung in Inline-Code"* — gemessen führen **86 von 89** DoD-Abschnitten `make gates` (84 Items), „alle" trifft auf 3 nicht zu. Die eigentliche Aussage (*Zähler misst das Falsche*) bleibt davon unberührt.

### L-3 · Slice §1/§3 tragen noch „160 Befunde bei 223 Dateien"
`docs/plan/planning/done/slice-179-…md:38,117` gegen ADR-0075 („80 Befunde") und DoD („166 Befunde"). Ich habe sechs Selektor-/Glob-Kombinationen durchprobiert und keine gefunden, die 160 liefert (nächstliegend: 161/211 bzw. 166/175). 223 ist plausibel die Zahl der `.md` unter `docs/plan/planning/` **vor** diesem Slice (heute 224) — die Befundzahl daneben bleibt unbelegt. Der Plan sagt selbst, die Messung sei nachgezogen worden; die Kopfzahl ist dabei stehengeblieben.

### L-4 · Der „zu breit"-Test fährt eine Konfiguration, die das Produkt mit Exit 2 ablehnt
`internal/hexagon/core/rules/structure_teilmenge_test.go:77,86` setzt `MaxTasks = ptr(-1)`. Der Config-Rand weist `max-tasks < 0` ab (gemessen: `max-tasks -1 muss >= 0 sein`, Exit 2). Damit belegt `„zu breit": 7 ignoriert` eine Meldung, die über eine gültige Konfiguration **nie** entstehen kann — bei `max-tasks ≥ 0` und `got = 0` greift `got > *r.MaxTasks` nicht. Der DoD-Haken (`slice-179-…md` §4) führt *„Ein zu breites Muster (**7 ignoriert**) … je als Test gefangen"*, was mehr klingt als das, was das Produkt leisten kann. Der Testkommentar benennt die Grenze korrekt; der DoD-Text nicht. (Der `0 ignoriert`-Fall ist dagegen bei `max-tasks: 0` erreichbar und in `TestTasksIgnorePattern_LiestDenBereinigtenText` echt gedeckt.)

---

## INFO

### I-1 · Herkunfts-Prosa im Modell-Kommentar
`internal/hexagon/core/model/config.go:487-489`: *„Zwei RE2 einer Regel mit zwei verschiedenen Zielen waeren die Falle, in die **der antragstellende Adopter selbst gelaufen ist** (ADR-0075)."* Der Kommentar trägt eine Abgrenzung und ein auflösbares Feld — der Nebensatz über den fremden Antragsteller ist Herkunfts-Prosa im Sinne von §3.7, und inhaltlich ist er der Gegenstand von M-4. Die Abgrenzung steht ohne ihn vollständig da (so wie in `structure.go:90-91`, wo derselbe Gedanke ohne die Zuschreibung formuliert ist).

---

## Negativbefunde (geprüft, mit Messung)

- **Byte-Identität ohne die Schlüssel — gemessen, nicht geglaubt.** `d-check:vor-0075` gegen `d-check:latest`, `max-tasks: 3` über `docs/plan/planning/**/slice-*.md` (**166 Befunde**: 77 `section-oversized` + 89 `section-missing`), in **fünf** Ausgabeformen: `plain` (166 Zeilen), `--json` (1337), `--yaml` (1001), `--doctor` (831), `--repair` (leer) — `cmp` auf stdout **und** stderr in allen fünf **identisch**, Exit 1/1. Die DoD-Zahl 166 reproduziert exakt.
- **Umkehr-Proben beißen (`BEO-023`).** Sieben Mutationen in einer Repo-Kopie, jede macht die *richtigen* Tests rot: Ignorier-Muster gegen die rohe Zeile → 7 rot · `exempt` gegen den Überschriften-**Text** → 6 rot · Nullmengen-Härte entfernt → 2 rot (`LeereMengeMeldet`, `StructureSelektorFormen`) · `exempt` **nach** der Kardinalität → **genau 1** rot (`LaeuftVorDerKardinalitaet`) · Zusatz immer anhängen → 1 rot (`AbwesendIstByteIdentisch`) · Zusatz nur bei `>0` → 3 rot · beide Config-Ränder entfernt (RE2-Kompilat, halbe Aktivierung) → je `TestDecode_StructureFehler` rot. Keine Mutation blieb grün.
- **Zwei Exit-2-Ränder, gemessen:** `tasks-ignore-pattern: '('` und `exempt-section-pattern: '('` → je Exit 2 mit Schlüsselname in der Meldung; `tasks-ignore-pattern` ohne `max-tasks` → Exit 2 *„halbe Aktivierung"*. `''` gilt bei beiden als Abwesenheit und fällt in die strenge Richtung (Exit 1 statt stillem Grün) — wie ADR-0075 Entscheidung 8 zusagt.
- **Nullmengen-Härte:** `exempt-section-pattern`, das alle Treffer nimmt → `section-missing` auf `line = 1`, Meldung *„alle 2 passenden Abschnitte sind von exempt-section-pattern ausgenommen — die Regel liefe leer"*. Kein stilles Grün (Exit 1). Deckt sich mit dem Akzeptanzkriterium.
- **Fence-Treue:** Task-Item im Fence zählt nicht; Überschrift im Fence wird weder gewählt noch ausgenommen (`FindSectionHeads` filtert vor der Ausnahme, `sections.go:19-37`).
- **Die Inline-Code-Messung des Antrags reproduziert exakt.** Je Alternative über die 444 Items: `make gates` **0**, `Closure-Notiz` **12**, `Beobachtungs-Register` **1**, `Risiko aus` **0**; 0 Fundstellen von `make gates` im bereinigten Text bei 84 im rohen. Auch die beiden „falschen" Treffer sind genau die im Slice zitierten Liefer-Punkte (`§2/§3/§4/§6 tragen SPEC-NNN …` und `ADR-0012-§Kern-Messung in der Closure-Notiz`).
- **Die Beispiel-Befundzeile in Handbuch/ADR/Antwort ist byte-echt.** Nachgefahren gegen ein Probe-Repo: `docs/slice-x.md:3⇥docs/slice-*.md :: ## 4. Definition of Done⇥section-oversized⇥Abschnitt trägt 4 Task-Items (3 ignoriert), erlaubt sind 3` — vier Tab-Spalten, identisch; dieselbe Zeichenkette in `--doctor` als `Hinweis:`.
- **Reihenfolge vor der Kardinalität:** zwei Treffer, einer ausgenommen, `sections: one` → **kein** `section-ambiguous`, befundfrei. Gemessen und durch M4-Mutation abgesichert.
- **Kein neuer Grund-Code, keine Message-Auswertung stromabwärts.** `ReasonSectionOversized` erscheint außerhalb der Regel nur in `diagnose.go:88,148` als Code, nie als geparster Text — die geänderte Meldung kann `--doctor`/`--repair`/`--json` nicht brechen (empirisch durch die Byte-Identität in allen fünf Modi bestätigt).
- **Hexagon-Richtung / Suppressions:** keine neuen Imports (`regexp`, `strconv`, `strings` waren da), `arch-check` 0 Befunde, `make lint` grün, keine `//nolint` im Diff (§3.2 ✓). `regexp.MustCompile` zur Laufzeit ist sicher: `StructureRule` wird ausschließlich in `configyaml.applyStructureRule` gebaut, und `structureBedingungsFehler` prüft beide neuen Muster vor dem Rückgabewert (`grep` über alle Nicht-Test-Konstruktoren: genau eine Stelle).
- **`--print-config`:** beide Schlüssel im Gerüst (Zeilen 178, 196 der Ausgabe), auskommentiert, mit korrekt beschriebenem Vergleichsgegenstand je Schlüssel. `--suggest-config` erzeugt keine `structure`-Regeln (`suggest.go:405` verweist nur aufs Voll-Schema) — keine vergessene Fläche dort. README/README.de zählen „neun **Bedingungen**"; die neuen Schlüssel sind keine Bedingung und tragen keinen Grund-Code, also keine Listen-Drift (`BEO-022` insoweit sauber bedient).
- **Namens-Begründung im Bestand belegt:** `exclude-sections` existiert als `[]string` literaler Überschriften-Titel in `vcs`/`immutable` (`model/config.go:322,337`) und in der Cross-/Sources-Konfiguration (`:787,865`) — die ADR-Begründung für den dritten Namen trägt.
- **`§3.4` Spec-Straten:** die neuen Lastenheft-/Spezifikations-Abschnitte nennen die ADR nur als *„Begründung in begleitender ADR"* und keine Slice-/Wellen-/Hash-Token; `doc-check` (Modul `matrix`) grün über 604 Dateien.
- **`§3.8` Reichweite:** das Modul liest durch diese Änderung keine Eingabe hinzu, die es nicht scannt — beide Muster arbeiten auf Zeilen bzw. Abschnitts-Text derselben, bereits gescannten Datei. Die einzige neue Reichweiten-Frage ist die *Asymmetrie* der beiden Vergleichsgegenstände, und die ist in Spezifikation, Handbuch, ADR und Kommentar viermal ausdrücklich benannt.
- **Der bekannte stille Pfad ist deklariert, nicht verschwiegen:** ein `tasks-ignore-pattern`, das so breit ist, dass die Schwelle nie fällt, macht die Regel stumm. Das steht als benannte Grenze in Lastenheft, Spezifikation, ADR (Konsequenzen), Handbuch, Antwort und Slice-§5 — kein Finding. Erst die Kombination mit `hint` (M-1) hebelt auch die Diagnose aus, die dagegenstehen soll.
- **Die fünf Zeilen der Paritäts-Tabelle des Absenders** sind je einem Test zugeordnet und laufen (`make test` grün).

---

## Repo-unverändert-Beleg

```
$ git status --porcelain      → (leer)
$ git rev-parse HEAD          → 78a1783820bb6a709ddc74dbe0c6c578b81d6512
$ ls -l .harness/state/gates-passed.diffsha
  → mtime 2026-08-30 08:28:16 +0200 (Autoren-Lauf; mein Review lief ab 08:45 — kein `make gates`, kein `record-gates`)
```

Alle Proben, Mutations-Kopien und Fixtures lagen unter `…/scratchpad/rev179/` und sind gelöscht.
