# Review Release-Prep v0.71.0 — `3254e8f`

**Review-Art:** Release-Prep (gegen [`docs/user/releasing.md`](../user/releasing.md)
§Release-Prep) — **kein** Slice-Review

**Gegenstand:** `3254e8f` (fünf Dateien:
[`version.md`](../../version.md), [`CHANGELOG.md`](../../CHANGELOG.md),
[`README.de.md`](../../README.de.md), [`README.md`](../../README.md),
[`docs/user/benutzerhandbuch.md`](../user/benutzerhandbuch.md)); der Inhalt von
slice-185 ist **nicht** Gegenstand (bereits reviewt und verifiziert:
[`2026-08-31-slice-185-code-r1.md`](2026-08-31-slice-185-code-r1.md),
[`2026-08-31-slice-185-verifikation.md`](2026-08-31-slice-185-verifikation.md)).
Zusätzlich geprüft, ob der Dependabot-Commit `1c32203` und die
Lifecycle-/Harness-Commits `66fcc03`, `7c7ea5a`, `3c29f76` außerhalb dieses
Diffs etwas berühren, das Release-Prep-Relevanz hätte.

**Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) @ 1.13.0
**Modell:** claude-sonnet-5 · **Datum:** 2026-08-31

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`docs/user/releasing.md`](../user/releasing.md) §Release-Prep, vollständig gelesen
- [`AGENTS.md`](../../AGENTS.md) §3.1 (Docker/make-only), §5 (Botschaft ≤ Messung, CHANGELOG-Regel)
- [`spec/lastenheft.md`](../../spec/lastenheft.md)
  [`DC-FA-CLI-010`](../../spec/lastenheft.md#dc-fa-cli-010--makefile-fragment-ausgeben) (0.80.0)
- [`CHANGELOG.md`](../../CHANGELOG.md) `[0.71.0]`, [`version.md`](../../version.md) §Aktuell/§Verlauf
- die beiden slice-185-Reports oben (Review + Verifikation, nicht erneut geprüft)

---

## Eigener Lauf

| Lauf | Ausgabe |
|---|---|
| `make build` | Image `sha256:43f6f246f6af…` (`d-check:latest`, `VERSION=0.0.0-dev`) |
| `make doc-check` | `d-check: 633 Datei(en) geprüft, 0 Befund(e)` — Exit 0 |
| `make ci` | `[gates] baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green` · `image-test: OK — DC-FA-DIST-001-Akzeptanzkriterien erfüllt` · `[ci] gates + image-test green` |
| `--print-mk` erzeugtes Fragment gegen das lokale Image, `doc-usage` gefahren (`make -f <frag> doc-usage DCHECK_IMAGE=d-check:latest`) | stdout 0 Byte, stderr 2488 Byte, Exit 0 — reproduziert exakt die CHANGELOG-/Handbuch-Zahl |
| `##`-annotierte Targets im erzeugten Fragment gezählt | 13, Namen und Reihenfolge identisch mit Handbuch §4.16-Klammerliste |
| `grep -oE 'ghcr\.io/pt9912/d-check:v0\.71\.0'` über die drei berührten Prosa-Dateien | 26 Treffer gesamt (24 Handbuch, 1 README.de, 1 README.md) — s. M1 |
| `grep` auf verbliebene `ghcr`-Pins ungleich `v0.71.0` außerhalb der Exempt-Pfade (`CHANGELOG.md`, `done/`, `spec/lastenheft.md`, `docs/reviews/**`) | keine Treffer |
| `git tag -l v0.71.0` | leer — Tag ist noch nicht gesetzt, Review läuft vor dem Tag |

`git status --short` zeigt am Ende nur diesen Report. `make fullbuild`, `make record-gates`
und `make image-scan` bewusst **nicht** gefahren (Nachweis-Datei bzw. Netz).

---

## Urteil

**TAG FREIGEBEN — mit einer Empfehlung.** 0 HIGH · 2 MEDIUM · 0 LOW · 1 INFO.

Die mechanische Hälfte der Prep ist vollständig und richtig: `version.md` §Aktuell/§Verlauf
inkl. wanderndem Anker, alle 26 `ghcr`-Pins und beide nackten Tags gezogen, `CHANGELOG.md`
korrekt unter der neuen Version geschnitten (kein `[Unreleased]`), Handbuch-Kopf und
§11-Zeile chronologisch unten mit dem Kalenderdatum des Commits, beide READMEs zu Recht
unverändert (kein neues Modul, Dogfooding-Zahl und Modul-Liste stimmen weiter), die
Operations-Referenz korrekt (Fix lag bereits im Feature-Commit, das CHANGELOG-`Fixed`
beschreibt ihn zutreffend), und `doc-usage` braucht keine eigene §4-Aufgabe — es reiht sich
so ein wie `doc-doctor`/`doc-repair` vor ihm. Der CHANGELOG-Eintrag ist gegen das aus HEAD
gebaute Produkt **exakt** nachgefahren: 0 Byte stdout, 2488 Byte stderr, Exit 0, 13 Targets,
kein `--disable` in der `doc-usage`-Recipe, `@`-Präfix vorhanden.

Die beiden MEDIUM betreffen **nicht** die ausgelieferte Dokumentation, sondern die
Commit-Botschaft selbst: eine falsch gezählte Pin-Zahl und ein fehlender Beleg für den in
§Release-Prep Punkt 5 verlangten `make ci`-Lauf (statt nur `make gates`). Beide sind
inhaltlich folgenlos — ich habe `make ci` unabhängig gefahren und es ist grün —, aber genau
die Klasse, die dieser Review-Typ fängt, bevor sie als falsche Selbstauskunft in der
Historie steht. Keine blockiert den Tag.

---

## MEDIUM

### M1 — Die Commit-Botschaft zählt die gezogenen `ghcr`-Pins falsch (28 statt 26, davon 26 statt 24 im Handbuch)

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §5 („Eine Commit-Botschaft … behauptet nicht mehr, als die Arbeit
  trägt") — Reviewer-Skill Frage 8
- `pfad`: Commit-Botschaft `3254e8f` (kein Datei-Pfad — die Zahl steht nur im Log)
- `befund`: Die Botschaft behauptet *„28 ghcr-Pins gezogen (26 Handbuch, je einer in beiden
  READMEs) — die meldet das versions-Gate."* Nachgezählt (`grep -oE
  'ghcr\.io/pt9912/d-check:v0\.71\.0'` über den Diff der drei Dateien) sind es **26**
  Treffer gesamt: **24** im Handbuch, je einer in `README.de.md` (Zeile 270) und `README.md`
  (Zeile 266). Die behauptete Zahl liegt um 2 zu hoch, ausschließlich beim Handbuch-Anteil.
  Zusätzlich benennt die Botschaft das `versions`-Gate als Quelle der Zahl — das Modul
  meldet aber nur **veraltete** Pins (`version-stale`), keine Gesamtzahl; die 28/26 ist eine
  reine Selbstzählung, keine Gate-Ausgabe.
- `verifizierbar`: ja — `git show 3254e8f -- docs/user/benutzerhandbuch.md README.de.md
  README.md | grep -oE 'ghcr\.io/pt9912/d-check:v0\.71\.0' | wc -l` liefert 26; je Datei
  aufgeschlüsselt 24/1/1.
- `klasse`: commit-botschaft-fehlzaehlung-ohne-gate-beleg

**Wirkung:** keine auf den ausgelieferten Stand — alle Pins sind tatsächlich korrekt
gezogen (bestätigt durch `make doc-check`: 0 Befunde, kein `version-stale`). Die falsche
Zahl steht ausschließlich in der Git-Historie und wird hier für die Nachwelt korrigiert,
nicht im Repo selbst (Commit-Botschaften werden nicht nachträglich verändert).

## M2 — Die Botschaft belegt nur `make gates`, obwohl §Release-Prep Punkt 5 `make ci` verlangt

- `kategorie`: MEDIUM
- `quelle`: [`docs/user/releasing.md`](../user/releasing.md) §Release-Prep Punkt 5 (*„`make
  ci` lokal grün fahren … erst dann taggen"*)
- `pfad`: Commit-Botschaft `3254e8f` (kein Datei-Pfad)
- `befund`: Die Botschaft schließt mit *„make gates gruen (633 Dateien, Exit 0)."* — das ist
  die Erfolgszeile des `doc-check`-Teilgates (`d-check: 633 Datei(en) geprüft, 0
  Befund(e)`), nicht die von `make gates` (`[gates] … green`) oder `make ci` (`[ci] gates +
  image-test green`). Ein Beleg für `make ci` — insbesondere `image-test`, das die
  `DC-FA-DIST-001`-Kriterien des Docker-Images prüft — fehlt. Das Vorgänger-Release
  (`11d8a60`, v0.70.0) hat denselben Checklisten-Punkt explizit erfüllt und benannt: *„make
  ci gruen (Exit 0): gates + image-test."*
- `verifizierbar`: ja — `make ci` lokal fahren; ich habe es getan (s. §Eigener Lauf): grün,
  `[gates] … green` gefolgt von `image-test: OK — DC-FA-DIST-001-Akzeptanzkriterien
  erfüllt` und `[ci] gates + image-test green`.
- `klasse`: prep-checkliste-punkt-nicht-belegt

**Wirkung:** keine auf die Freigabefähigkeit — `make ci` ist zum Zeitpunkt dieses Reviews
grün, das Image-Kriterium ist erfüllt. Der Fund ist die fehlende **Selbstauskunft**, nicht
ein rotes Gate.

---

## INFO

### I1 — `make ci` lief in diesem Review nicht gegen das Release-Prep-Diff isoliert, sondern gegen den vollen Arbeitsbaum inkl. slice-185-Code

Das ist Absicht: `docs/user/releasing.md` Punkt 5 verlangt den Lauf gegen den Stand vor dem
Tag insgesamt, nicht gegen den Prep-Diff isoliert. Genannt, damit klar ist, dass M2 sich auf
das **Fehlen des Belegs in der Botschaft** bezieht, nicht auf einen tatsächlich roten Lauf.

---

## Negativbefunde

- geprüft, ohne Befund: `version.md` §Aktuell + §Verlauf + `<a id>`-Wanderung (Zeilen 22, 35;
  alte `v0.70.0`-Zeile hat den Anker korrekt verloren)
- geprüft, ohne Befund: Vollständigkeit der `ghcr`-Pin-Ziehung — keine `version-stale`-Funde
  (`make doc-check`, 0 Befunde) und kein verbliebener `v0.70.0`-`ghcr`-Pin außerhalb der
  Exempt-Pfade (Zahl selbst siehe M1)
- geprüft, ohne Befund: die zwei baren Tags (Docker-Hub-Pull, `:vX.Y.Z`-Beispiel in
  §Versionen und Tags) — beide auf `v0.71.0` gezogen
- geprüft, ohne Befund: `CHANGELOG.md` — kein `[Unreleased]`-Rest, `[0.71.0]` korrekt unter
  `[0.70.0]` geschnitten, Datum `2026-08-31` stimmt mit `git log -1 --date=short` überein
- geprüft, ohne Befund: Handbuch-Kopf (Version 1.65, Software-Version v0.71.0, Stand
  2026-08-31) und §11-Zeile — letzte Zeile der Tabelle, chronologisch unter `1.64`
- geprüft, ohne Befund: Frage nach einer eigenen §4-Aufgabe für `doc-usage` — verneint
  zu Recht; die Dokumentation reiht sich in §4.16 wie `doc-doctor`/`doc-repair` ein,
  keine geänderte Bestandszusage, kein neuer Abschnitt nötig
- geprüft, ohne Befund: beide READMEs — Dogfooding-Zeile weiterhin „elf Module"/„eleven
  modules" (`.d-check.yml` `modules:`-Liste unverändert bei 11), Modul-Liste unverändert,
  nur die `docker pull`-Zeile auf v0.71.0 gezogen (Zeile 270 DE / 266 EN)
- geprüft, ohne Befund: `docs/user/operations.md` — die CHANGELOG-`Fixed`-Zeile beschreibt
  akkurat, was der vorangehende Feature-Commit (`8e75841`) dort geändert hat (drei-von-zwölf-
  Aufzählung durch Fragment-Zeiger ersetzt, `DCHECK_DIGEST` ergänzt); die Release-Prep selbst
  musste die Datei nicht mehr anfassen
- geprüft, ohne Befund: Exempt-Pfade — `CHANGELOG.md` (historische Einträge unverändert),
  `done/`, `spec/lastenheft.md`, `docs/reviews/**` — keiner trägt einen fälschlich gezogenen
  Pin, keiner wurde vom Prep-Commit außerhalb seines legitimen Zwecks angefasst
- geprüft, ohne Befund: Release-Prep als **ein** Commit vor dem Tag, getrennt von den
  Slice-Commits (`8e75841` … `835605e`) und vom Dependabot-Bump (`1c32203`, zu Recht nicht im
  CHANGELOG — reine indirekte Abhängigkeitshebung ohne Nutzersicht)
- geprüft, ohne Befund: CHANGELOG-Zusage gegen das Produkt — 0 Byte stdout/2488 Byte
  stderr/Exit 0 exakt reproduziert, 13 Targets exakt wie im Fragment, keine
  `--enable`/`--disable`-Flags in der `doc-usage`-Recipe, `@`-Echo-Unterdrückung vorhanden
- geprüft, ohne Befund: SemVer-Einordnung `0.71.0` (MINOR) — rein additives Target, keine
  entfernte/umbenannte Fähigkeit, kein geänderter Default

---

## Was ich nicht geprüft habe

- Den **Inhalt** von slice-185 (bereits reviewt/verifiziert, ausdrücklich nicht Gegenstand).
- Die **Gültigkeit** eines Digest-Pins gegen die Registry — es gibt noch keinen; er entsteht
  laut Prozess erst nach dem Tag als Folge-Commit.
- `make fullbuild`, `make image-scan`, die Nacht-Läufe (`freshness-*`, Digest-Achsen) — außerhalb
  des Release-Prep-Diffs und teils netzgebunden.
- Übersetzungstreue der EN-Fassung jenseits der von diesem Commit berührten Zeilen (unverändert,
  daher kein neuer Anlass).

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** commit-botschaft-fehlzaehlung-ohne-gate-beleg ·
prep-checkliste-punkt-nicht-belegt

## Verdikt

**Tag freigeben.** 0 HIGH, 2 MEDIUM — beide betreffen ausschließlich die Selbstauskunft der
Commit-Botschaft (eine falsch gezählte Pin-Zahl, ein fehlender `make ci`-Beleg statt nur
`make gates`), nicht die ausgelieferte Dokumentation oder das Produkt. Für beide gilt: der
zugrunde liegende Zustand ist geprüft und in Ordnung (`make doc-check` 0 Befunde, `make ci`
unabhängig grün gefahren). Die Commit-Botschaft selbst wird nicht nachträglich verändert
(keine Amend-Historie); der Befund steht hier für den Steering Loop.
