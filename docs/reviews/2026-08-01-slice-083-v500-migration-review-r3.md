# Re-Review-Report (R3, eng fokussiert): slice-083 — N-1-Fix (`conventions.md`-Anker-Kaskade) — 2026-08-01

**Review-Art:** Eng fokussierte Ziel-Re-Review (R3), unabhängiger Frischkontext.
Prüft **ausschließlich einen Punkt**: ob der Fix für den R2-Befund **N-1** den
`conventions.md`-Anker-Bruch der Migration vollständig und korrekt heilt, ohne
neuen Fehler. **Kein Voll-Review** (die übrigen R1/R2-Achsen sind nicht Gegenstand).

**Gegenstand:** `docs/plan/planning/open/slice-083-regelwerk-v500-migration-analyse.md`
(uncommittet) — konkret §2.3 (conventions-Bullet „Anker-Kaskade (Review-F-2/N-1)")
und Etappe C Schritt 2 („Anker-Erhalt") + Schritt 8.

**Vorlauf:** R2-Report `docs/reviews/2026-08-01-slice-083-v500-migration-review-r2.md`
(1 offener MEDIUM: N-1, Re-Eröffnung von R1-F-2c).

**Skill:** `.harness/skills/reviewer.md` v1.2.0 · **Modell:** claude-opus-4-8 · **Datum:** 2026-08-01

**Eingangs-Kontext (geprüfte Verträge/Quellen):**

- Der N-1-Abschnitt aus dem R2-Report (Anker-Bruch, Voll-Slug, zehn immutable ADRs).
- Der Fix in slice-083 §2.3 / Etappe C Schritt 2 + 8.
- Anker-Engine: `internal/hexagon/core/rules/anchors.go` (DC-FA-ANCH-001) —
  Slug-Bildung, Inline-HTML-`<a id>`-Erfassung, Fragment-Auflösung.
- Doku-Gate: `.d-check.yml` (scan-Skopus, `modules`, `ids`/`anchors`-Ventile).
- Der Ist-Link-Bestand: alle `conventions.md#mr-…`-Links repo-weit
  (ohne `.harness/{baseline,cache}`), die zwölf Accepted-ADRs unter `docs/plan/adr/`,
  die Review-Reports unter `docs/reviews/`.

**Methode:** eigener Beleg statt Zusage — Repo-Greps mit Datei-/MR-Klassifikation,
Anker-Engine am Code gelesen, Rohbyte-Prüfung eines Links. REFUTE nur mit Zitat.

---

## N-1-Verdikt: **GEHEILT** — mit einem verbleibenden LOW (Reviews-Kante)

Der Fix adressiert den Kern von N-1 korrekt: `conventions.md` behält je von
immutabler/eingefrorener Doku referenziertem MR einen **Voll-Slug**-`<a id>`-Anker
— **auch für die aufgelösten/entfallenen** — in einem **Anker-Kompatibilitäts-Block
unabhängig vom aktiv-/`done/`-Schnitt**, deklariert als migrationsspezifischer `MR`.
Beleg zu den vier Prüfpunkten:

### Prüfpunkt 1 — Deckt der Fix wirklich alle Links (inkl. der nicht-aktiven)?

**Ja für die immutablen ADRs; explizit auch die nicht-aktiven benannt.** Gemessen
(Grep, ohne `.harness/{baseline,cache}`): **174** `conventions.md#mr-…`-Vorkommen
(R2 zählte 173; die eine Differenz ist der seither hinzugekommene R2-Report selbst,
`docs/reviews/2026-08-01-…-r2.md:71` mit `#mr-007` — bestätigt den Mechanismus, kein
Widerspruch). **Alle** tragen den **Voll-Slug** (`#mr-NNN--voller-titel-slug`);
**null** die Kurzform `#mr-NNN` — der R2-Einwand gegen die frühere Kurz-Anker-Fassung
ist damit bestätigt und im Fix behoben (§2.3 „Die Links sind **Voll-Slug**"; Etappe C
Schritt 2 schreibt `<a id="mr-NNN--voller-titel-slug">`).

Die zwölf Accepted-immutablen ADRs und die von ihnen referenzierten MRs (eigener Grep):

| ADR (Status Accepted) | referenzierte MR | aktiv? |
|---|---|---|
| 0019 | mr-003, mr-007 | nein / nein |
| 0021 | mr-006 | ja |
| 0022 | mr-006 | ja |
| 0023, 0024, 0026, 0027, 0029, 0031 | mr-007 | nein |
| 0028 | mr-007, mr-013 | nein / ja |
| 0030 | mr-017, mr-019 | nein / nein |
| 0046 | mr-019, mr-022 | nein / nein |

**Zehn** der zwölf ADRs (alle außer 0021/0022) zeigen auf **nicht-aktive** MRs
(mr-003 aufgelöst, mr-007 historisch, mr-017 entfällt, mr-019/mr-022 verschmilzt) —
deckungsgleich mit §2.4. Der Fix nennt genau diese Kohorte explizit als abgedeckt:
§2.3 „**zehn** dieser ADR-Links zeigen auf **nicht-aktive** MRs … **auch für die
aufgelösten**" und Etappe C Schritt 2 „**auch die aufgelösten/entfallenen** (zehn der
zwölf immutablen ADRs zeigen auf nicht-aktive MRs)". **Alle** ADR-Links sind damit
gedeckt — **keine** Accepted-ADR muss editiert werden. Der harte N-1-Kern (Links in
nicht-editierbaren ADRs) ist vollständig aufgelöst.

### Prüfpunkt 2 — Technisch tragfähig (Anker-Engine)?

**Ja, am Code belegt.** `htmlAnchors` (anchors.go:120–135) liest `id`-Werte an
**beliebigen** Elementen **wörtlich** in die Anker-Menge (`set[v] = true`); `AnchorSet`
(155–161) vereinigt sie mit den Heading-Slugs; `CheckAnchors` löst per **Exakt-Match**
auf (`slugs[a.frag]`, Zeile 219 — kein Präfix/Teilstring). Rohbyte-Prüfung eines
ADR-Links (`0023-immutable-core-pin.md:25`): das Fragment ist **literales UTF-8**
`mr-007--auflösung-…` (kein `%`-Escape repo-weit — Grep: 0 perzent-kodierte Links);
`resolveAnchorRef` dekodiert per `url.PathUnescape` nur `%XX` und lässt das literale
`ö` unangetastet (anchors.go:185). Die aktuelle Überschrift
`### MR-007 — Auflösung von MR-003: doc-check als Dogfooding` slugged via `Slugify`
(65–79; Unicode-Kleinschreibung, `ö` als Buchstabe erhalten, `—`/`:` verworfen) zu
`mr-007--auflösung-von-mr-003-doc-check-als-dogfooding` — genau dem Link-Fragment.
Ein `<a id="mr-007--auflösung-von-mr-003-doc-check-als-dogfooding">` (literales `ö`)
erzeugt denselben Set-Eintrag und löst den Link exakt auf. **Der Migrations-Schritt
ist so `anchors`-gate-grün** — Voraussetzung: der Umsetzer schreibt den `<a id>`-Wert
**byte-genau** als den bestehenden Heading-Slug (inkl. literalem `ö`), was der Fix mit
„Voll-Slug" fordert.

### Prüfpunkt 3 — C2/C5-Widerspruch aufgelöst?

**Ja.** R2 rügte, C2 (Anker für alle 173) widerspreche C5 (Streichen/`done`-Verschieben
der nicht-aktiven Einträge). Der Fix entkoppelt beides ausdrücklich: der Anker-Block ist
„**unabhängig** vom aktiv-/`done/`-Schnitt der Index-Zeilen — **damit Schritt 5
(Streichen/Verschieben) die Anker nicht mitnimmt**" (Etappe C Schritt 2). Die Anker leben
in einem eigenen Block, den Schritt 5 (der auf Index-Zeilen und Einzeldateien wirkt)
nicht anfasst; nach Schritt 5 verbleiben sie in `conventions.md`. Die Etappe-C-Reihenfolge
ist damit in sich konsistent; Schritt 8 stützt sich korrekt auf beide Mechaniken
(Anker-Block + A6-Tombstone).

### Prüfpunkt 4 — Neuer Fehler durch den Fix?

**Kein blockierender neuer Fehler; drei Achsen geprüft:**

- **Default-Konflikt „Index = nur aktive Adaptionen":** aufgelöst. Der Block wird „als
  eigener `MR` deklariert (migrationsspezifisch)" — das ist unter v5.0.0 die
  **kanonisch korrekte** Auflösung (eine repo-spezifische Abweichung wird als MR
  geführt). Die `<a id>`-Tags tragen keinen lesbaren Adaptions-Text und wachsen den
  Pflicht-Lesekontext nicht spürbar — der Sinn des Splits bleibt gewahrt.
- **Doppel-Anker auf aktiven MRs:** unkritisch. In der Index-**Tabellen**-Form haben
  aktive MRs keine Überschrift, also keinen Heading-Slug — der `<a id>` ist der einzige
  Anker (kein Duplikat). Selbst bei versehentlicher Doppelung kollabiert `AnchorSet` als
  Map (`set[v]=true`) verlustfrei; es gibt keine „duplicate-anchor"-Regel.
- **`ids`-Interaktion:** kein neuer Befund. Der `<a id="mr-…">`-Wert ist kleingeschrieben;
  das `ids`-Muster ist `MR-\d{3}` (Großschreibung, .d-check.yml:30) — kein Treffer, keine
  Linkpflicht auf den Anker.

---

## Neue Findings

### F-R3-1 (LOW) — Der Anker-Block-Erhalt nennt nur ADRs/`done/`; die von einem Review-Report referenzierten `mr-015`/`mr-016` fallen durch die Enumeration

- **kategorie:** LOW
- **quelle:** slice-083 §2.3 (conventions-Bullet) / Etappe C Schritt 2 („Anker-Erhalt") / DC-FA-ANCH-001
- **pfad:** `docs/plan/planning/open/slice-083-regelwerk-v500-migration-analyse.md`
  (§2.3-Bullet + Etappe C Schritt 2); Beleg-Fundstellen
  `docs/reviews/2026-06-25-baseline-v140-mr016-mr017-r1.md:8` (`#mr-016`) und `:29`
  (`#mr-015`); `internal/hexagon/core/rules/anchors.go:219`; `.d-check.yml:16` (`anchors`
  aktiv, kein `docs/reviews`-Ausschluss für `anchors`)
- **befund:** Der Fix erhält Anker „je von immutabler/eingefrorener Doku referenziertem
  MR", enumeriert als Beleg aber nur „zehn der zwölf immutablen ADRs" und die
  `done/`-Einträge. `mr-015` (aktiv) und `mr-016` (historisch) werden von **keiner** ADR
  und **keinem** `done/`-Slice referenziert, wohl aber vom eingefrorenen Review-Report
  `2026-06-25-baseline-v140-mr016-mr017-r1.md`. Wird „eingefrorene Doku" nicht als
  Review-Reports mitgelesen, fehlen `mr-015`/`mr-016` im Anker-Block; da `anchors` auch
  `docs/reviews/` prüft (nur `ids`/`codepaths`/`versions` haben dort ein Ventil, `anchors`
  nicht), meldet `make gates` nach Etappe C `anchor-missing` für die zwei Links in jenem
  Report. `mr-015` ist zudem ein nicht-offensichtlicher Fall: obwohl **aktiv**, verliert es
  in der Index-**Tabellen**-Form seinen Heading-Anker und braucht den Kompatibilitäts-Anker
  ebenso — die Fix-Betonung „auch die **aufgelösten**" kann darüber hinwegtäuschen.
- **verifizierbar:** ja — nach Etappe C meldet das `anchors`-Modul `anchor-missing` für
  `2026-06-25-…-r1.md:8` und `:29`, sofern der Anker-Block-Bestand strikt aus ADR- und
  `done/`-Referenzen gebildet wird und die zwei Review-Links weder mitgedeckt noch
  retargetet werden.

**Warum LOW, nicht blockierend:** (a) Die **Regel** des Fixes ist bereits weit genug
(„je von immutabler/**eingefrorener** Doku referenziertem MR") — ein Umsetzer, der
Review-Reports als eingefroren behandelt (die Skill-Konvention verlangt „Nie
überschreiben"), deckt `mr-015`/`mr-016` **ohne** Spec-Änderung ab. (b) Reviews stehen
unter **keinem** Immutabilitäts-Gate (kein `vcs`/adr-immutable auf `docs/reviews/`), das
Retargeten der zwei Links bleibt als Ausweg. (c) Der Skopus ist zwei Links in einer Datei;
die für Etappe C Pflicht-vorgesehene Frischkontext-Review (§4) fängt ein solches
`anchor-missing` am Gate. Es ist eine **Vollständigkeits-Lücke der Enumeration**, kein
Fehler des Mechanismus — der Anker-Block-Inventar-Schritt sollte Review-referenzierte MRs
ausdrücklich einschließen.

---

## Negativbefunde (je Prüfachse — geprüft)

**Achse 1 (Vollständigkeit ADRs).** 174/12-ADR-Bestand am Grep nachgemessen; alle zehn
nicht-aktiven ADR-Ziele (mr-003/007/017/019/022) sind vom Fix explizit als gedeckt
benannt. Geprüft, ohne Befund.

**Achse 2 (Anker-Engine).** `htmlAnchors` wörtlich, Exakt-Match, literales `ö` unversehrt,
Slug-Rekonstruktion deckungsgleich — der `<a id>`-Voll-Slug trägt. Geprüft, ohne Befund.

**Achse 3 (C2/C5-Kohärenz).** Block ausdrücklich schnitt-unabhängig, Schritt 5 nimmt die
Anker nicht mit; Reihenfolge konsistent. Geprüft, ohne Befund.

**Achse 4 (neuer Fehler).** Default-Konflikt via migrationsspezifischer MR-Deklaration
aufgelöst; kein Doppel-Anker-Problem; keine `ids`-Interaktion. Geprüft, ein LOW (F-R3-1),
kein blockierender Befund.

**Nicht-Gegenstand (R3-Charter):** die `ids`-`target`-Umkonfiguration (MR-Verweise linken
künftig auf `harness/conventions/MR-NNN-…md` statt `harness/conventions.md`) ist ein
eigener Migrations-Aspekt außerhalb der N-1-Anker-Frage und wurde nicht bewertet.

---

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 0 | — |
| LOW | 1 | F-R3-1 (Review-referenzierte mr-015/mr-016 nicht enumeriert) |
| INFO | 0 | — |

N-1-Heilungsstand: **GEHEILT** — Voll-Slug-Anker (Prüfpunkt 2 code-belegt),
nicht-aktive ADR-MRs explizit gedeckt (Prüfpunkt 1), C2/C5 entkoppelt (Prüfpunkt 3),
kein blockierender Neu-Fehler (Prüfpunkt 4).

---

## Gesamt-Verdikt

**Abnahmereif** — kein offener HIGH/MEDIUM. Der N-1-Fix heilt den `conventions.md`-
Anker-Bruch am Kern korrekt und technisch tragfähig: die Links der zwölf immutablen ADRs
(inkl. der zehn auf nicht-aktive MRs) lösen nach Etappe C **ohne ADR-Edit** über einen
Voll-Slug-`<a id>`-Anker-Kompatibilitäts-Block auf, der vom aktiv-/`done/`-Schnitt
entkoppelt ist und als migrationsspezifischer `MR` deklariert wird; der zuvor
widersprüchliche C2/C5-Ablauf ist konsistent. Verbleibend ist ein **LOW** (F-R3-1): die
Beleg-Enumeration nennt nur ADRs/`done/` und übergeht die von einem eingefrorenen
Review-Report referenzierten `mr-015`/`mr-016` — heilbar durch die bereits weit genug
formulierte Regel („eingefrorene Doku" schließt Review-Reports ein) oder durch Aufnahme
dieser zwei MRs in den Anker-Block-Bestand; ein Blocker ist es nicht. Empfehlung an die
Abnahme: F-R3-1 als Anker-Block-Inventar-Hinweis (Review-referenzierte MRs mitzählen) in
Etappe C vermerken, dann Analyse annehmen.
