# Review-Report: slice-153 — Sagt das Lastenheft noch, was `matrix` tut?

**Datum:** 2026-08-26 · **Review-Art:** Vertrags-/Text-Review (Spec-Diff gegen
Slice-Plan, Kanon, Konventionsspeicher und die lebenden Spiegel im Baum) ·
unabhängiger Reviewer ohne Anteil an der Arbeit
**Gegenstand:** Commit-Kette `1343863..d9fbfe0` von slice-153 — `1343863`
(Anspruchs-/Lifecycle-Move), `d9fbfe0` (Feat: `spec/lastenheft.md` +8/−2,
`spec/spezifikation.md` +1/−1, `AGENTS.md` +3/−1, `.d-check.yml` +6/−4)
**Skill:** `.harness/skills/reviewer.md` @ 1.10.0 (`9ee805b`) · **Modell-ID:**
`claude-opus-5[1m]`

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/done/slice-153-lastenheft-token-ziel-klasse.md`
  vollständig — §1 Ziel, §2 (die vier Schritte), §3 „Ausdrücklich NICHT",
  §4 DoD, §5 Risiken, §7 Vorgelagert
- `spec/lastenheft.md` vorher (`git show d9fbfe0^:spec/lastenheft.md`) und
  nachher: `DC-FA-MTX-001` (Zeilen 1188–1231 inkl. Akzeptanzkriterien und
  Out-of-Scope), `DC-FA-MTX-002` (1233–1263), `DC-FA-MTX-003` (1267–1312),
  §1 Zweck (22–34), das Glossar (2960–2974) und §7 Historie (2975–…)
- `spec/spezifikation.md` — §2-Schema `SPEC-005`, die fünf
  `matrix.classes[]`-Zeilen (2576–2580), §`DC-FA-MTX-001.a` Schritte 1–6
  (1042–1125), §4 `SPEC-018` (2725) und die Historie ab 2790
- `.d-check.yml` — `matrix` vollständig (acht Klassen, vierzehn `rules`,
  `status`, `exclude-sections`, `exempt-paths`), dazu `ids`, `structure`,
  `scan`, `modules`
- `AGENTS.md` §3.4 vorher (`git show d9fbfe0^:AGENTS.md`) und nachher, dazu
  §3.7, §3.8, §4, §5
- `harness/conventions/MR-032-historie-vor-accepted.md`,
  `harness/conventions/MR-025-spiegel-vor-dem-editieren.md`,
  `harness/conventions.md` (Adaptions-Tabelle, `MR-000`)
- Baseline `v5.12.0`: `grundlagen-source-precedence.md` §Spec-Stratifizierung
  und §Wann die CR-Pflicht beginnt (vollständig, inkl. der Absätze über
  zusammenfallende Rollen und über die Tatsachenberichtigung);
  `modul-03-spec.md` §Ziel-Form: Akzeptanzkriterium und §Ziel-Form:
  Spezifikation
- `docs/plan/planning/observations.md` — `BEO-002`, `BEO-009`, `BEO-011`,
  `BEO-012`
- Quellcode zur Gegenprobe: `internal/hexagon/core/rules/matrix.go`
  (`CheckMatrix`, `classOf`, `tokenFindings`),
  `internal/adapter/driven/configyaml/configyaml.go` (`applyMatrix`),
  `internal/hexagon/core/model/config.go` (`MatrixClass`),
  `internal/hexagon/core/app/diagnose.go` (`reasonTexts`),
  `internal/adapter/driving/cli/config_template.go`
- **Vorherige Findings am gleichen Modul (Pflicht):**
  `docs/reviews/2026-08-26-slice-144-commit-hash-muster-review.md` — F-6 MEDIUM
  (`hard-rule-stuetzt-sich-auf-undokumentierte-config-freiheit`, das Finding,
  das diesen Slice ausgelöst hat; es nennt neben Lastenheft und Schema-Zeile
  ausdrücklich `internal/hexagon/core/model/config.go:144`) sowie die
  Finding-Klassen `zahl-ohne-messvorschrift` und
  `risiko-klasse-schmaler-benannt-als-das-muster`; und
  `docs/reviews/2026-08-26-slice-145-pfad-token-sicht-review.md` — F-1 HIGH
  (`sensor-luecke-als-erlaubnis-in-die-regel-geschrieben`), F-3 MEDIUM
  (`wegwahl-begruendet-mit-einer-widerlegten-modul-eigenschaft`)

**Nicht erhalten** (Skill §Eingangs-Kontext): die DoD-Abhakung. Die vier Haken
in slice-153 §4 werden hier **nicht** bewertet; der Slice liegt in
`in-progress/`. Ebenso nicht bewertet: die Bestätigung der Gate-Läufe aus der
Botschaft (`make gates`, `make fullbuild`) — das ist Verifikations-, nicht
Reviewer-Rolle.

**Vom Reviewer selbst gefahren** (nur lesend; Sonden ausschließlich als
Overlay über den read-only-Mount bzw. in einem eigenen Sonden-Repo im
Scratchpad — **keine Repo-Datei wurde verändert**, `git status --short` vor
und nach jedem Lauf leer, kein `git checkout` nötig):

- `make doc-check` ⇒ **Exit 0**, „518 Datei(en) geprüft, 0 Befund(e)"; das
  Image wurde dabei aus HEAD gebaut (`sha256:f41f99ad9a1e…`) und ist die Basis
  aller folgenden Läufe.
- **Sonde A** (Overlay-Config über `.d-check.yml`, Regel
  `{from: commit-hash, to: adr, allow: false}` ergänzt) ⇒ „518 Datei(en)
  geprüft, 0 Befund(e)", **Exit 0**. **Sonde B** (dieselbe Regel mit
  undeklariertem Klassennamen) ⇒ „`.d-check.yml`: matrix.rules[12]
  referenziert undeklarierte Klasse", **Exit 2**.
- **Sonde C** (eigenes Sonden-Repo: eine Klasse mit `paths`, eine
  Token-Ziel-Klasse `commit-hash`, verbotene Kante, eine Datei mit einem
  Hash-Token) ⇒ `--doctor` meldet „Z. 1 · **Referenz zwischen Dokumentklassen
  nicht erlaubt** (Referenzrichtung) [matrix]", Exit 1.
- `--print-config` über das Repo ⇒ Zeile 40 „`# --- matrix: erlaubte
  Referenzen zwischen Dokumentklassen ---`".
- **Token-Gegenprobe** der neuen Zeilen (Link-Spans entfernt, dann
  `slice-\d{3}`, `welle-\d{2,}`, `\b[0-9a-f]{7,40}\b`, `\b(internal|cmd)/[a-z]`
  über 1275–1279, 2970 und 2979): kein Treffer.
- `grep` nach dem **alten** Wortlaut („Dokumentklasse") über den ganzen Baum
  ohne `done/`, `docs/reviews/` und `.harness/baseline/` — der von `MR-025`
  benannte Ableiter für eine Wortlaut-Änderung.

---

## Findings

### F-1

- **kategorie:** MEDIUM
- **quelle:** `MR-025` §Adaption (Spiegel-Tabelle: Zeile *Anforderung* →
  `spec/lastenheft.md`, Zeile *Klartexte* → „`AllReasons()` und
  `reasonTexts()` (`--doctor`)", Zeile *Emittierte Vorlage* →
  „`--print-config`/`--suggest-config`", Zeile *Nutzer-Doku*) und §*„Die Liste
  wird aus dem Repo abgeleitet"* (*„Für eine Wortlaut-Präzisierung ist der
  Ableiter das `grep` nach dem **alten** Wortlaut über den ganzen Baum"*) ·
  `BEO-002` · Reviewer-Skill §MEDIUM *Botschaft verallgemeinert über die
  Messung hinaus* · slice-144-Report F-6, der
  `internal/hexagon/core/model/config.go:144` mit Datei:Zeile übergeben hat
- **pfad:** Commit-Botschaft `d9fbfe0`, Absatz *„DIE SPIEGEL SIND VORHER
  GEZÄHLT (MR-025), nicht danach gesucht: vier lebende Stellen … Alle vier
  tragen jetzt denselben Begriff und zeigen auf dieselbe Anforderung.
  Eingefroren bleiben die `done/`-Slices und die Review-Reporte"* gegen
  `spec/lastenheft.md:28`, `spec/lastenheft.md:2971`,
  `internal/hexagon/core/model/config.go:144`,
  `internal/hexagon/core/rules/matrix.go:28`,
  `internal/hexagon/core/app/diagnose.go:113`,
  `internal/adapter/driving/cli/config_template.go:48`, `README.de.md:31`,
  `docs/user/benutzerhandbuch.md:33`, `:319`, `:1874`
- **befund:** Der Ableiter, den `MR-025` für genau diesen Fall vorschreibt
  (`grep` nach dem alten Wortlaut), findet **zehn** weitere lebende Zeilen in
  sieben Dateien, die die Knoten der Referenzmatrix weiterhin ausnahmslos als
  Dokumentklassen führen — zwei davon in der bearbeiteten Datei selbst (§1
  Zweck *„Referenzrichtungs-Regeln zwischen Dokumentklassen
  (Referenzmatrix)"* und die Glossar-Zeile `Referenzmatrix` **eine Zeile
  unter** dem neuen Eintrag), eine als Typ-Definition im Kern (*„MatrixClass
  ist eine über Pfad-Globs deklarierte Dokumentklasse"*), die der auslösende
  Review bereits mit Datei:Zeile übergeben hatte. Die Aufzählung der
  eingefrorenen Menge (`done/`-Slices, Review-Reporte) nennt zudem die
  Spezifikations-Historie nicht, die als Protokollzeile denselben Sachverhalt
  in einer fünften Formulierung führt (`spec/spezifikation.md:2794`, *„reines
  **Token-Ziel**"*).
- **verifizierbar:** ja — `grep -rn "Dokumentklasse" --include=*.md
  --include=*.go .` ohne `done/`, `docs/reviews/` und `.harness/baseline/`;
  kein Gate meldet es (`make doc-check` Exit 0, 518 Dateien, 0 Befunde, vor
  wie nach der Änderung — die Spiegel sind Prosa, `BEO-002`).
- **klasse:** `spiegel-zaehlung-vor-dem-editieren-unvollstaendig`

### F-2

- **kategorie:** MEDIUM
- **quelle:** `spec/lastenheft.md:2970` (neuer Glossar-Eintrag: *„sie ist damit
  **keine** Dokumentklasse"*) · `MR-025` Spiegel-Zeilen *Klartexte* und
  *Emittierte Vorlage* · `internal/hexagon/core/app/diagnose.go:100-103`
  (Kommentar bindet die Klartexte an `spec/spezifikation.md` §4) ·
  `spec/spezifikation.md:2725` (`SPEC-018` sagt dort neutral *„Referenz
  zwischen **Klassen** nicht erlaubt"*)
- **pfad:** `internal/hexagon/core/app/diagnose.go:113`;
  `internal/adapter/driving/cli/config_template.go:48`
- **befund:** Der `--doctor`-Klartext für `matrix-forbidden` nennt die
  verletzte Beziehung *„Referenz zwischen Dokumentklassen nicht erlaubt"*,
  obwohl **drei der vierzehn** Regeln dieses Repos auf einer Klasse enden, die
  das neue Glossar ausdrücklich nicht als Dokumentklasse führt
  (`spec-straten → commit-hash`, `sicht → commit-hash`, `sicht → modul-pfad`);
  wer den Befund nachschlägt, findet im Glossar die Auskunft, dass es die
  genannte Beziehung nicht gibt. Dieselbe Formulierung emittiert
  `--print-config` als Überschrift des `matrix`-Blocks in jede Adopter-Config.
- **verifizierbar:** ja — Sonde C: `--doctor` im Sonden-Repo gibt für einen
  Commit-Hash-Token wörtlich *„Z. 1 · Referenz zwischen Dokumentklassen nicht
  erlaubt (Referenzrichtung) [matrix]"* aus (Exit 1); `--print-config` über
  dieses Repo liefert Zeile 40.
- **klasse:** `nutzer-klartext-widerspricht-dem-neuen-vertragsbegriff`

### F-3

- **kategorie:** MEDIUM
- **quelle:** Baseline `modul-03-spec.md` §Ziel-Form: Akzeptanzkriterium
  (*„drei Pfade im Given/When/Then-Stil — Happy · Boundary · Negative"*) ·
  slice-153 §2 Schritt 1 (*„Er steht unter **Beschreibung**, nicht unter
  Akzeptanzkriterien — das ist ein Unterschied, und er gehört gelesen"*) ·
  Commit-Botschaft (*„die Akzeptanzkriterien sagen über Pfade nichts"* und
  *„Minor, weil eine Anforderung wächst"*)
- **pfad:** `spec/lastenheft.md:1275-1279` (die neue Zusage) gegen
  `spec/lastenheft.md:1306-1308` (die drei unveränderten Akzeptanzkriterien
  von `DC-FA-MTX-003`)
- **befund:** Die neue Zusage steht ausschließlich in der Beschreibung; die
  drei Akzeptanzkriterien sprechen weiter nur über `token` plus Marker, über
  den fehlenden Marker und über `exempt-paths` und kennen die pfadlose Klasse
  nicht. Damit trägt genau die Textsorte die Zusage, die dieselbe Botschaft
  einen Absatz zuvor als **nicht**-zusagend gegen `DC-FA-MTX-001` in Anschlag
  bringt — und das Wachstum, mit dem der Minor-Bump begründet wird, ist an
  keinem Kriterium messbar.
- **verifizierbar:** ja — die in slice-144-Report F-6 erwogene Härtung (ein
  fail-closed-Rand auf leere `paths`, Exit 2) verletzt **keines** der drei
  Akzeptanzkriterien von `DC-FA-MTX-003` und hebelt die 0.66.0-Zusage
  dennoch aus; heute belegt durch Lesen von `spec/lastenheft.md:1304-1308`
  gegen `:1275-1279`.
- **klasse:** `neue-zusage-ohne-akzeptanzkriterium`

### F-4

- **kategorie:** MEDIUM
- **quelle:** `spec/lastenheft.md:1276-1277` (neu: *„hat keine Mitglieder und
  ist reines *Ziel*"*) · `spec/spezifikation.md:2579`
  (`matrix.classes[].direction`: *„keine still wirkungslose
  Richtungs-Deklaration"* — dasselbe Modul, derselbe Schema-Block) ·
  `internal/adapter/driven/configyaml/configyaml.go:1960-1962` (`applyMatrix`
  prüft für `rules` ausschließlich, ob der Klassenname deklariert ist) ·
  `internal/hexagon/core/rules/matrix.go:37` und `:88` (`srcClass` stammt
  immer aus `classOf`, also aus `paths`)
- **pfad:** `spec/lastenheft.md:1276`
- **befund:** Der Vertrag sagt seit 0.66.0 zu, dass eine Token-Ziel-Klasse
  *ausschließlich* Ziel ist, während der Konfigurations-Rand sie
  widerspruchslos als `from` einer Regel annimmt; die Regel kann dann nie
  feuern, weil `classOf` für eine Klasse ohne `paths` nie einen Treffer
  liefert. Ein Adopter, der die Kante in dieser Richtung deklariert, bekommt
  kein Config-Signal, sondern ein stilles Grün — dieselbe Klasse
  „wirkungslose Deklaration", die `order`/`direction` im selben Schema-Block
  ausdrücklich fail-closed ausschließen.
- **verifizierbar:** ja — Sonde A (`{from: commit-hash, to: adr, allow:
  false}` als Overlay) ⇒ „518 Datei(en) geprüft, 0 Befund(e)", Exit 0; Sonde B
  (derselbe Eintrag mit undeklariertem Namen) ⇒ „matrix.rules[12] referenziert
  undeklarierte Klasse", Exit 2.
- **klasse:** `vertragszusage-ohne-fail-closed-rand`

### F-5

- **kategorie:** LOW
- **quelle:** `BEO-012` (*„Eine Quelle wird über ihren Geltungsbereich hinaus
  zitiert … der Text der Quelle stimmt, die in Anspruch genommene Reichweite
  nicht"*) · Commit-Botschaft, Absatz *„DIE BEISPIELE SIND BEISPIELE, DIE
  DEFINITION IST EINE ZUSAGE"*
- **pfad:** `spec/lastenheft.md:1190-1191` (*„Die Konfiguration deklariert
  Dokumentklassen über Pfad-Muster (z. B. Contract-Spec …)"*) und
  `spec/lastenheft.md:1270-1273` (*„`matrix` erkennt verbotene Referenzen
  bisher nur als **Links** … Eine Klasse kann **daher** zusätzlich ein
  `token`-Muster tragen"*)
- **befund:** Zwei der drei Textbelege tragen die Behauptung nicht, für die
  sie zitiert werden: das „z. B." steht vor der Aufzählung der Klassen, nicht
  vor dem Mechanismus *„über Pfad-Muster"* — dieser Halbsatz ist ungehedgt —,
  und das „zusätzlich" steht in einem Satz, den ein „daher" an den
  vorangehenden Kontrast Token **gegen Link** bindet (*„bisher nur als
  Links"*), nicht an einen Kontrast Token gegen Pfade; die Schema-Zeile sagt
  dazu passend *„Ohne `token` nur Link-Erkennung"*. Getragen wird das Ergebnis
  allein vom dritten Beleg, der Glossar-Definition.
- **verifizierbar:** nein — Lesart am Text, kein Gate; nachlesbar an den vier
  genannten Zeilen und an `spec/spezifikation.md:2580`.
- **klasse:** `beleg-traegt-die-behauptung-nicht-fuer-die-er-zitiert-wird`

### F-6

- **kategorie:** LOW
- **quelle:** Baseline `modul-03-spec.md` §Ziel-Form: Spezifikation (*„Kein
  Kopf-Datum, kein Kopf-Status — die Spezifikation trägt ihre Änderungen in
  der Historie"*) · `spec/spezifikation.md:2794-2796` (die drei jüngsten
  Bestandszeilen; zwei davon sind ausdrücklich *„redaktionell, keine
  Semantik-Änderung"* und haben je eine Zeile bekommen)
- **pfad:** `spec/spezifikation.md:2577` (die geänderte Schema-Zelle) gegen
  `spec/spezifikation.md:2790-2794` (Historie-Tabelle)
- **befund:** Der Commit ändert die Schema-Zelle `matrix.classes[].paths`,
  ohne der Spezifikations-Historie eine Zeile zu geben; die oberste
  Bestandszeile beschreibt dieselbe Zelle vom selben Tag noch mit dem
  Wortlaut *„reines **Token-Ziel**"*, den der neue Begriff gerade abgelöst
  hat. Wer die Historie des Technik-Stratums liest, findet den Begriff, den
  die Zelle heute führt, dort nicht eingeführt.
- **verifizierbar:** ja — `git show d9fbfe0 -- spec/spezifikation.md` zeigt
  genau einen geänderten Hunk und keine Historie-Zeile; Gegenprobe
  `git log --oneline -- spec/spezifikation.md` gegen die Datums-Spalte der
  Tabelle.
- **klasse:** `technik-stratum-ohne-historie-zeile-geaendert`

### F-7

- **kategorie:** LOW
- **quelle:** `MR-032` §Adaption, die den Kanon zitiert (*„Welche Stelle der
  Version steigt, entscheidet das Repo und gehört in den Adaptions-Block von
  `harness/conventions.md`"*) · Baseline `grundlagen-source-precedence.md`
  §Wann die CR-Pflicht beginnt, letzter Absatz · `harness/conventions.md`
  (Adaptions-Tabelle: kein Eintrag zur Stellenwahl) ·
  `spec/lastenheft.md:2988` (0.62.1: *„Präzisierung ohne
  Verhaltensänderung"* ⇒ Patch), `:2980`–`:2983` (0.65.1–0.65.4: *„keine
  Anforderungs-Änderung"* ⇒ Patch)
- **pfad:** `spec/lastenheft.md:2979` (Zeile `0.66.0`, Marker *„**Kein
  Verhaltens-Delta:** die Fähigkeit bestand in der Umsetzung, war aber nicht
  zugesagt"*)
- **befund:** Die neue Zeile trägt zugleich die Minor-Stelle und die
  Marker-Formel, mit der dieses Repo in fünf der letzten zehn Zeilen
  **Patch**-Bumps ausgewiesen hat; welches der beiden Kriterien vorgeht
  („eine Anforderung wächst" gegen „kein Verhaltens-Delta"), steht nirgends
  geschrieben — der Adaptions-Block, den der Kanon dafür ausdrücklich
  benennt, führt zur Stellenwahl keinen Eintrag. Die nächste Zeile derselben
  Art ist damit aus der Historie in beide Richtungen ableitbar.
- **verifizierbar:** nein — Urteil an der Historie-Praxis; die Abwesenheit des
  Konventions-Eintrags ist per `grep -rn -i "minor\|patch\|Stelle der
  Version" harness/conventions.md harness/conventions/` belegbar (nur
  `MR-032` und `MR-039`, beide zu anderen Fragen).
- **klasse:** `bump-stelle-ohne-geschriebenes-kriterium`

### F-8

- **kategorie:** LOW
- **quelle:** `MR-025` §*„Der Spiegel ist die Stelle, nicht die Datei"* ·
  `.d-check.yml:223-225` (`modul-pfad`: *„TOKEN-ZIEL-KLASSE (DC-FA-MTX-003):
  kein Pfad-Muster, keine Mitglieder, reines Ziel … wie bei commit-hash"*) ·
  `.d-check.yml:308` (Regel `{from: sicht, to: modul-pfad, allow: false}`)
- **pfad:** `AGENTS.md:173-176` (Begriff gezogen) gegen `AGENTS.md:195-205`
  (derselbe Abschnitt, zweiter Absatz)
- **befund:** §3.4 trägt den neuen Begriff nur im Commit-Hash-Absatz; der
  zweite Absatz beschreibt mit dem Modul-Pfad-Sensor dieselbe Bauform — eine
  Klasse ohne `paths`, nur mit `token` — und nennt als Träger stattdessen
  *„die Sicht als **eigene Klasse**"*, während die Klasse `modul-pfad`, ohne
  die der Befund nicht entsteht, in der Hard Rule überhaupt nicht vorkommt.
  Wer §3.4 liest, sieht eine Token-Ziel-Klasse, wo die Konfiguration zwei
  fährt.
- **verifizierbar:** ja — `.d-check.yml:235-236` (Klasse `modul-pfad` mit
  `token`, ohne `paths`) und `:308` gegen `AGENTS.md:195-205`; ohne diese
  Klasse liefert der Sicht-Sensor keinen Befund (belegt in
  slice-145-Report F-3, Sonden A/B/C).
- **klasse:** `begriff-nur-in-einem-der-zwei-absaetze-nachgezogen`

### F-9

- **kategorie:** LOW
- **quelle:** slice-153 §Berührte Spec-Stellen (*„`spec/lastenheft.md`,
  §`DC-FA-MTX-001` — **falls** die Antwort so ausfällt"*) · slice-153 §1 (der
  Satz, der den Slice auslöste, stammt aus `DC-FA-MTX-001`) · slice-153
  §Betroffene IDs (nennt allein `DC-FA-MTX-001`)
- **pfad:** `docs/plan/planning/done/slice-153-lastenheft-token-ziel-klasse.md:15-17`
  und `:91-92` gegen `git show --stat d9fbfe0`
- **befund:** Geändert wurden `DC-FA-MTX-003`, das Glossar und die Historie;
  `DC-FA-MTX-001` — die einzige als berührt ausgewiesene Spec-Stelle und die
  einzige unter „Betroffene IDs" — ist unverändert geblieben. Wer über diese
  Felder sucht, welche Anforderung in 0.66.0 gewachsen ist, landet bei der
  falschen; die Commit-Botschaft trägt beide Kennungen und verdeckt die
  Differenz.
- **verifizierbar:** ja — `git show d9fbfe0 -- spec/lastenheft.md` zeigt
  keinen Hunk zwischen den Zeilen 1188 und 1231.
- **klasse:** `beruehrte-spec-stelle-benennt-die-nicht-geaenderte-anforderung`

### F-10

- **kategorie:** LOW
- **quelle:** `MR-025` §Adaption, Zeile *Config-Schema* ·
  `git show d9fbfe0^:spec/spezifikation.md` Zeile 2577 (*„eine Klasse **ohne
  `paths`** hat keine Mitglieder"*)
- **pfad:** `spec/spezifikation.md:2577`
- **befund:** Die Zelle beginnt nach *„Leer zulässig:"* mit *„eine **solche**
  Token-Ziel-Klasse"*; das Demonstrativpronomen hatte in der Vorfassung
  *„eine Klasse ohne `paths`"* als Bezugswort, das mit dem Begriffs-Tausch
  entfallen ist. Die Zelle nennt damit die Bedingung nicht mehr, unter der
  eine Klasse zur Token-Ziel-Klasse wird; sie steht nur noch im
  Lastenheft-Glossar, also im Stratum darüber.
- **verifizierbar:** nein — Lesbarkeit am Text; nachstellbar über
  `git show d9fbfe0^:spec/spezifikation.md` gegen den heutigen Stand.
- **klasse:** `demonstrativ-ohne-antezedens-nach-begriffs-tausch`

---

## Negativbefunde

1. **Geprüft, ohne Befund — Kernfrage D, der Verzicht auf einen CR.** Die
   Behauptung der Botschaft hält in beiden Hälften: `spec/lastenheft.md:5`
   trägt `**Status:** Draft`, und der Kanon sagt wörtlich *„Vor `Accepted` ist
   das Lastenheft ein Entwurf — frei änderbar, ohne Change Request, ohne
   Historie-Zeile; die Trennung von Entscheidung und Umsetzung greift noch
   nicht, weil noch nichts versprochen wurde."* Auch der Absatz, der bei
   zusammenfallenden Rollen einen **eigenen** Commit vor dem umsetzenden
   Slice verlangt, greift nicht: er beschreibt den *angenommenen Change
   Request*, also den Zustand ab `Accepted`. `MR-032` verlangt trotzdem Bump
   und Historie — beides ist da, und die `Verweis`-Spalte trägt `—`, wie
   `MR-032` es für die Zeit vor `Accepted` festhält. Der Rückzug der
   zwischenzeitlichen Vorlage als Auftraggeber-Entscheidung ist damit
   belegbar richtig, nicht bloß bequem.
2. **Geprüft, ohne Befund — Form und Position der Historie-Zeile.** Vier
   Spalten (`| 0.66.0 | 2026-08-26 | … | — |`, maschinell nachgezählt), oben
   eingefügt und damit chronologisch: 2026-08-26 über 2026-08-25. Die Position
   ist gate-gedeckt — `.d-check.yml` führt für `spec/lastenheft.md` §7 die
   `structure`-Regel `table-order: desc`; ein Einschub an falscher Stelle
   erzeugt `section-unordered`.
3. **Geprüft, ohne Befund — keine Abwärts-Referenz in den neuen Zeilen.** Nach
   Entfernen der Markdown-Link-Spans tragen die Zeilen 1275–1279, 2970 und
   2979 keinen Treffer für `slice-\d{3}`, `welle-\d{2,}`,
   `\b[0-9a-f]{7,40}\b` oder `\b(internal|cmd)/[a-z]`; die einzige Referenz
   der Historie-Zeile ist ein dokument-interner Anker auf `DC-FA-MTX-003`
   (gleiche Klasse, gleicher Rang, also weder `matrix-forbidden` noch
   `matrix-downward`). `make doc-check` bestätigt es mit Exit 0 über 518
   Dateien.
4. **Geprüft, ohne Befund — Kernfrage F, die beiden §3-Verbote.** Keine
   Rücknahme der Klasse: `classes` und `rules` in `.d-check.yml` sind
   unverändert, der Diff berührt dort ausschließlich Kommentarzeilen (`git
   show d9fbfe0 -- .d-check.yml`, sechs `#`-Zeilen). Keine Ausweitung auf
   andere `DC-*`-Beschreibungen: die Hunks im Lastenheft liegen in
   `DC-FA-MTX-003`, im Glossar, in der Historie und in der Kopf-Version — in
   keiner anderen Anforderung.
5. **Geprüft, ohne Befund — „kein Verhaltens-Delta" trägt.** Der Commit fasst
   keine Datei unter `internal/` oder `cmd/` an und keinen semantischen
   Schlüssel der Konfiguration. Die Fähigkeit bestand tatsächlich vorher:
   `applyMatrix` validiert je Klasse nur `Name != ""` und die Eindeutigkeit
   (`configyaml.go:1944-1958`), und `classOf` liefert für eine Klasse ohne
   `paths` nie einen Treffer (`matrix.go:219-228`) — sie hat damit keine
   Mitglieder, bleibt aber als Token-Ziel wirksam. Das ist die eine
   Zuschreibung der Botschaft, die eine `BEO-009`-Prüfung ohne Rest übersteht.
6. **Geprüft, ohne Befund — Kernfrage E, der Glossar-Eintrag gegen die
   gefahrene Konfiguration.** Die Definition *„Klasse ohne Pfad-Muster, die
   nur ein `token`-Muster trägt"* deckt genau die zwei Klassen, die dieses
   Repo pfadlos fährt (`modul-pfad`, `.d-check.yml:235`; `commit-hash`,
   `:266`), und schließt `sicht` (`paths: [spec/architecture.md]`,
   `:221-222`) korrekt aus. **Keine** lebende Stelle behauptet, es gebe nur
   eine solche Klasse — die Prüfung auf Exklusivitäts-Wörter („einzig",
   „nur", „erste", „als einzige") über `AGENTS.md`, `.d-check.yml` und die
   beiden Spec-Straten liefert keinen Treffer. F-8 ist eine Auslassung, keine
   Exklusivitäts-Behauptung; der Unterschied ist `BEO-011` Ausprägung (a) und
   ist hier **nicht** eingetreten.
7. **REFUTED — die Inertness-Aussage im Lastenheft ist kein Stratum-Verstoß.**
   Der naheliegende Verdacht (*„Ohne Pfad-Muster **und** ohne `token`-Muster
   ist eine Klasse inert"* sei eine technische Festlegung, die nach
   `modul-03-spec.md` §Ziel-Form: Spezifikation ins Technik-Stratum gehört)
   ist am Bestand widerlegt: das Lastenheft führt Inertness-Zusagen an sieben
   weiteren Stellen (`:259`, `:949`, `:1800`, `:2125`, `:2348`, `:2730`,
   `:2777`), und die wörtliche Doppelung mit der Spezifikation ist Hausform
   (`:949` *„fehlendes `refs` macht den Eintrag inert"* gegen
   `spec/spezifikation.md:832`). Kein Finding.
8. **REFUTED — `DC-FA-MTX-002` ist kein weiterer Spiegel.** Sein
   Beschreibungs-Satz *„Eine **Dokumentklasse** … kann zusätzlich eine
   geordnete Rangfolge und eine Richtungspolitik tragen"* setzt Mitglieder
   voraus (`order` ist eine Liste von Pfad-Globs, der Rang ist der Index des
   ersten Treffers **einer Datei**); für eine Klasse ohne `paths` ist die
   Fähigkeit definitionsgemäß leer. Die Einschränkung auf „Dokumentklasse"
   ist dort richtig und bleibt es. Ebenso unberührt: die Grund-Code-Tabelle
   `spec/spezifikation.md:2725`, die *„zwischen **Klassen**"* sagt und unter
   beiden Vokabularen wahr bleibt — sie ist der Gegen-Beleg dafür, dass die
   Formulierung in F-2 kein Zwang war.
9. **Geprüft, ohne Befund — Kernfrage A, die dritte Aussage.** Die
   Glossar-Definition *„Dokumentklasse — über Pfad-Muster definierte Gruppe
   von Dokumenten"* stand vor dem Commit wörtlich so da
   (`git show d9fbfe0^:spec/lastenheft.md`, Glossar) und ist unverändert
   geblieben; die Behauptung, eine Klasse mit `token` und ohne `paths` sei im
   eigenen Vokabular keine Dokumentklasse gewesen, trägt an diesem Text
   allein. Der Schnitt „neben der Definition statt an ihr" ist also nicht
   erfunden, sondern belegt — nur seine beiden Nebenbelege tragen nicht
   (F-5).
10. **Geprüft, ohne Befund — `DC-FA-MTX-003` im Übrigen.** Provenance-Marker,
    Grandfathering, `exempt-paths` und die Out-of-Scope-Liste sind
    unverändert; der eingefügte Text steht zwischen zwei Bestandssätzen und
    unterbricht die Kette *Token erkennen → Kante A→B → `matrix-forbidden`*
    nicht (nachgelesen an `:1270-1284`).

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 6 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:**
`spiegel-zaehlung-vor-dem-editieren-unvollstaendig` ·
`nutzer-klartext-widerspricht-dem-neuen-vertragsbegriff` ·
`neue-zusage-ohne-akzeptanzkriterium` ·
`vertragszusage-ohne-fail-closed-rand` ·
`beleg-traegt-die-behauptung-nicht-fuer-die-er-zitiert-wird` ·
`technik-stratum-ohne-historie-zeile-geaendert` ·
`bump-stelle-ohne-geschriebenes-kriterium` ·
`begriff-nur-in-einem-der-zwei-absaetze-nachgezogen` ·
`beruehrte-spec-stelle-benennt-die-nicht-geaenderte-anforderung` ·
`demonstrativ-ohne-antezedens-nach-begriffs-tausch`

## Verdikt

**Merge-blockierend:** ja — vier MEDIUM, kein HIGH.

**Die Antwort auf die Kernfrage hält, ihre Begründung nur zur Hälfte.** Der
Slice hat die richtige Frage gestellt und sie nicht bequem beantwortet: die
Klasse wurde nicht zurückgenommen, die Definition der Dokumentklasse nicht
geweitet, und der Verzicht auf einen CR ist am Kanon-Wortlaut belegbar (N-1).
Getragen wird die Entscheidung aber allein von der Glossar-Definition; die
beiden Nebenbelege — das „z. B." und das „zusätzlich" — sagen nicht, was ihnen
zugeschrieben wird (F-5). Das ist genau die Bewegung, für die dieses Repo
`BEO-012` führt: der Text stimmt, die in Anspruch genommene Reichweite nicht.

**Zur Gegenthese** (hätte man statt eines neuen Begriffs die Definition der
Dokumentklasse weiten können): Sie ist nicht offensichtlich schlechter, und
die Botschaft benennt ihren Preis nicht. Eine geweitete Definition — „über
Pfad-Muster **oder** `token`-Muster definierte Gruppe" — hätte die zehn
Bestandszeilen aus F-1 sämtlich wahr gelassen und keinen einzigen Spiegel
erzeugt; die gewählte Abgrenzung macht sie zu Teilaussagen und bindet den
Begriff dauerhaft an eine Pflege, die kein Gate stützt. Was für die Abgrenzung
spricht, steht in der Botschaft und ist stichhaltig (der Gegenstand *ist* eine
Zeichenkette, kein Dokument; eine Gruppe von Dokumenten ohne Mitglieder wäre
ein Widerspruch im Begriff). Was dagegen spricht, ist nirgends abgewogen — und
die Rechnung dafür steht in F-1 und F-2.

**Die Buchführung über den Bestand ist die eigentliche Lücke.** Die Botschaft
sagt „vier lebende Spiegel, alle vier gezogen" und führt `MR-025` als Beleg;
der von `MR-025` selbst vorgeschriebene Ableiter — `grep` nach dem **alten**
Wortlaut — findet zehn weitere Zeilen, zwei davon in der bearbeiteten Datei
und eine im Kern, die der auslösende Review bereits mit Datei:Zeile übergeben
hatte (F-1). Zwei dieser Stellen sind nicht Prosa, sondern Ausgabe: der
`--doctor`-Klartext und die emittierte Config-Vorlage sagen einem Nutzer
„Referenz zwischen Dokumentklassen", während das Glossar ihm sagt, dass es
diese Beziehung für seinen Befund nicht gibt (F-2). Das ist `BEO-002` in der
Form, für die `MR-025` geschrieben wurde, und es ist bemerkenswert, dass die
Botschaft die Regel korrekt zitiert und ihren Ableiter trotzdem nicht gefahren
hat.

**Die beiden verbleibenden MEDIUM betreffen die Zusage selbst.** Sie steht
ohne Akzeptanzkriterium in derselben Textsorte, die die Botschaft zwei Absätze
zuvor als nicht-zusagend abwertet (F-3) — wer der Argumentation des Commits
folgt, hat nach 0.66.0 genauso wenig eine Zusage wie vorher. Und sie sagt
„ausschließlich Ziel", ohne dass der Config-Rand die Gegenrichtung kennt: eine
Regel mit einer Token-Ziel-Klasse als Quelle läuft mit Exit 0 durch und feuert
nie (F-4, am Produkt nachgestellt). Beides ist Vertragstext, den eine spätere
Härtung oder ein Adopter beim Wort nimmt.

Die sechs LOW sind Ränder derselben Bewegung: eine Historie, die im
Technik-Stratum fehlt (F-6), eine Bump-Stelle ohne geschriebenes Kriterium
(F-7), ein Begriff, der in einem von zwei Absätzen derselben Hard Rule
ankommt (F-8), ein Plan-Feld, das auf die nicht geänderte Anforderung zeigt
(F-9), und ein Demonstrativpronomen ohne Bezugswort (F-10).

**Was ausdrücklich hält:** die Entscheidung, der Verzicht auf den CR, die Form
und Position der Historie-Zeile, die Token-Freiheit der neuen Zeilen, beide
§3-Verbote und die Aussage „kein Verhaltens-Delta" — letztere als einzige
Zuschreibung dieser Botschaft, die eine `BEO-009`-Prüfung ohne Rest übersteht.
