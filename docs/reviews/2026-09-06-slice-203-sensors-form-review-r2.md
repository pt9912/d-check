# Review-Report — slice-203 (Runde 2, Endstand)

**Review-Art:** Code/Design — geprüft gegen Slice-Plan, Baseline-Kanon (`grundlagen-harness-dateien.md` §Sensors, `modul-13` §Hard Rule / *Die dritte Lage*, `gate.template.md`) und [`AGENTS.md`](../../AGENTS.md) §3/§4/§5
**Gegenstand:** `b521262..7ab22f2` — die vier Slice-Commits `13c46ab`, `b160184`, `aadb07e`, `b949e07` (`7ab22f2` ist ein CR-Dokument, das während des Reviews entstand und nicht zum Slice gehört)
**Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) @ 1.13.0 · **Modell-ID:** claude-opus-5[1m] · **Datum:** 2026-09-06
**Eingangs-Kontext:** Slice-Plan (neunteilig), `harness/sensors/*.md` (23), `harness/README.md` §Sensors, `harness/conventions.md` §Adaptions-Block, `.d-check.yml` (`structure`), `AGENTS.md` §6, Findings der ersten Runde (7 MEDIUM / 6 LOW / 2 INFO)

**Gefahrene Läufe.** `make doc-check` (650 Dateien, 0 Befunde), `make gates` **grün, Exit 0** (zehn Gates). Dazu **sieben eigene Break-Tests** der neuen Regel, alle zurückgerollt (`md5sum` vor/nach identisch, `git status` leer).

> Ein Zwischenbefund war fremd: Beim ersten `doc-check`-Lauf war der Zeilenlängen-CR untracked und rot (`anchor-missing`). Die Datei entstand **während** des Reviews und wurde als `7ab22f2` korrigiert committet. Kein Slice-Befund.

---

## Findings

### MEDIUM

**M-1 · Selbstauskunfts-Zahl reproduziert nicht**
`quelle` AGENTS.md §5 · `pfad` Slice-Plan §2 · `verifizierbar` ja (Nachmessung, kein Gate) · `klasse` selbstauskunft-zahl-reproduziert-nicht

§2 sagt *„**21 von 31** Zeilen lagen über 200"*, mit der eigens deklarierten Messgröße Zellinhalt. Nachgemessen mit genau dieser Messgröße ergibt `13c46ab^` **24 von 31** — und `13c46ab` ebenfalls 24. Die 21 reproduziert an keiner Commit-Grenze. Ebenso `b160184`s Botschaft *„Der größte Vertrag ging von 4997 auf **106** Zeichen"*: 4997 stimmt, die Zielzahl ist **60**.

**M-2 · Der neuen Regel werden zwei Funde zugeschrieben, die sie nicht finden kann**
`quelle` AGENTS.md §4, §5 · `pfad` Commit `aadb07e` · `verifizierbar` ja (Break-Test) · `klasse` sensor-fund-falsch-zugeschrieben

Die Botschaft sagt, die neue Regel habe drei vorbestehende Defekte gefunden, darunter die fehlende Bindung-Spalte bei `image-scan` und `workflow-pins`. Gemessener Break-Test: einer Zeile die Bindung-Zelle **ganz** entfernt — dieselbe Form, die beide vorher hatten — meldet **null Befunde**. Die gebundene Spalte ist Position 2 (`Vertrag`); eine Zeile mit zwei Zellen erfüllt sie. Nur der MR-054-Fund ist der Regel zurechenbar.

**M-3 · Eine Zahl widerspricht ihrer eigenen Aufzählung — fünfmal**
`quelle` Maintainability / AGENTS.md §5 · `verifizierbar` ja · `klasse` zahl-widerspricht-ihrer-aufzaehlung

- `doc-check.md` — *„**Sieben** Modul-Klassen"*, danach **acht** Glieder; der Lauf fährt **elf**. Die alte Zelle nannte keine Zahl; die Zahl ist beim Umzug entstanden.
- `test.md` zählt *„zwei Zusagen"*, führt die Grenzen aber mit *„**Die dritte Zusage** ist die mit den meisten Löchern"* ein — Folge des F-10-Nachzugs (3→2), nicht durchgezogen.
- `lint.md` — *„**24** kalibrierte Linter: 5 Default- plus 23 …, dazu `nolintlint`"* summiert sich zu 29. Die Quelle bindet die 24 an die Nicht-Default-Gruppe.
- `hooks.md` — *„`pre-commit`, **zwei Teile**"*; `.githooks/pre-commit` führt drei Schritte.
- `workflow-pins`: alte Zelle *„Drei benannte Grenzen"*, `AGENTS.md` *„Zwei"*, Träger-Datei vier Punkte ohne Zahl.

**M-4 · Eine Grenze wird „permanent" genannt, deren Quelle einen Auflösungs-Trigger führt**
`quelle` `MR-048` · `pfad` `harness/sensors/test.md` · `klasse` quelle-ueber-geltungsbereich

Die Datei nennt die Skip-Hälfte „permanent"; `MR-048` führt dafür einen Auflösungs-Trigger. Als permanent deklariert, wird die Lücke beim Eintritt des Triggers von keinem Retirement-Check mehr angefasst.

**M-5 · Exit-Vertrag ohne den `make`-Vorbehalt, den drei Schwesterdateien tragen**
`quelle` `DC-FA-CLI-003` · `pfad` `harness/sensors/test.md` · `klasse` exit-vertrag-ohne-make-normalisierung

`make test` ist ein Docker-Build-Recipe; GNU Make normalisiert jeden fehlgeschlagenen Recipe auf Exit 2 — über das Target ist die 1 nicht beobachtbar. Drei Schwesterdateien tragen diesen Vorbehalt, `test.md` nicht. F-3 aus Runde 1, an anderer Datei wieder aufgetaucht.

**M-6 · Beim Umzug verlorene Still-Grün-Zusage**
`quelle` ADR-0028 · `pfad` `harness/sensors/planning-check.md` · `klasse` grenz-aussage-beim-umzug-verloren

Alte Zelle: *„beide Richtungen, **fail-closed via Heading-Guard**"*. Die Träger-Datei nennt weder „fail-closed" noch „Heading-Guard", die Index-Zeile auch nicht. Damit steht die Aussage, dass eine fehlende Überschrift nicht still grün passiert, nur noch in der ADR — ausgerechnet in der Datei, deren Abschnitt *Grenze* heißt.

**M-7 · Das Nicht-Gate-Kriterium ist nicht auf alle Zeilen angewandt**
`quelle` das angeschriebene Kriterium in `harness/README.md` · `klasse` nicht-gate-kriterium-inkonsistent

- `make doc-complete` (Werkzeug) und `make completeness-check` (Gate) haben **byte-identische** Recipes; getrennt werden sie ausschließlich durch den Bindepunkt — das ausdrücklich disqualifizierte Kriterium.
- `make archive-wave-test` (Werkzeug) fährt eine Go-Testsuite, urteilt also fail-closed über Repo-Code — dieselbe Klasse wie `make test` (Gate).

Die F-4-Korrektur selbst ist tragfähig; sie wurde nur nicht auf `doc-complete`, `archive-wave-test` und `baseline-probe` angewandt.

**M-8 · Der Vorprüfungs-Block §7 beschreibt einen Zuschnitt, den es nicht mehr gibt**
`quelle` Baseline `modul-05`; `MR-054` · `klasse` vorpruefungsblock-nach-wachstum-veraltet

§7 sagt zur Beobachtung `large-migration-exceeds-session-review-limit`: *„Der Schnitt auf **drei Träger** ist die Antwort darauf. Ein dritter Beleg entsteht damit **nicht**."* `b949e07` zog §3 und §4 nach, §7 nicht. Es gibt keinen Schnitt auf drei Träger mehr. **Daran hängt die Register-Zählung der Closure:** Bleibt die Aussage, zählt das Register auf 2 und die Beobachtung erreicht die Schwelle nie.

**M-9 · Deckungsnachweis in der Träger-Datei — das §5-Risiko ist eingetreten**
`quelle` `gate.template.md` (*„Was hier NICHT steht: womit das Werkzeug selbst gedeckt ist"*) · `klasse` deckungsnachweis-in-der-sensor-datei

Fünf Dateien tragen ihn, drei unter der Überschrift *Grenze — was das Grün nicht abdeckt*, wo ein bestandener Selbsttest keine Lücke ist: `adr-check.md`, `trace-check.md`, `planning-check.md` (*„Der Negativ-Selbsttest lebt als Akzeptanztest im Modul"*), `gate-consistency.md`, `arch-check.md` (*„per Proben-Matrix belegt, **nicht behauptet**"* — der Verstärker kam beim Umzug hinzu).

### LOW

**L-1 · Sechs Targets fehlen in beiden Tabellen, während die DoD Vollständigkeit behauptet** — `build`, `run`, `deps`, `compile`, `help`, `clean`. `make build` ist Prerequisite von `image-test` und damit von `ci`/`fullbuild`. Kein Sensor deckt diese Richtung: `targets.authority` ist `AGENTS.md`.

**L-2 · Die `## Bindung`-Sektionen sind verkürzte Kopien der Index-Zelle** — die Vorlage verlangt dieselbe Angabe. `doc-check.md` führt drei statt acht Kennungen; `adr-check`, `trace-check`, `completeness-check` verlieren die Ablöse-Relationen.

**L-3 · Die Taxonomie behauptet eine Partition, die `image-scan` nicht enthält** — es steht in der Gate-Tabelle und in keiner der beiden Klassen; ausgerechnet der Grenzfall, den `13c46ab` benannt hat.

**L-4 · `harness/README.md` §Minimal agent workflow nicht mit `AGENTS.md` §6 nachgezogen** — Schritt 8 steht dort weiter als Abschluss.

**L-5 · Heilung falsch zugeordnet** — `hooks.md`: *„heilbar nur durch Gegenlesen"*. Die Quelle trennt: sichtbar durch `git show`, geheilt durch `--amend`.

**L-6 · Klassen-Marker „Netz" verloren** — `baseline-freshness.md` nennt ihn nur als Fehlerursache, die Schwesterdateien als Klasse.

**L-7 · Grenzen ohne Quelle** — `arch-check.md` und `semgrep.md` machen aus der generischen Gate-Klassen-Zeile der Baseline eine Zusage dieses Targets; die „Rollen-Treue"-Grenze ist repo-weit unikal.

**L-8 · Probenklasse breiter zugesagt als geprobt** — `guard-probe.md` sagt „Host-Toolchain", geprobt wird als Toolchain nur eine.

### INFO

**I-1 · Die Waisen-Richtung ist ungewächtert.** Bleibt nach einem Target-Retirement die Datei ohne Index-Zeile liegen, dokumentiert sie ein Target, das es nicht gibt. Heute keine Waise (23/23 verlinkt).

**I-2 · `harness/sensors/` ist eine neue Artefaktklasse, die `MR-045` nicht nennt.** `hooks.md` trägt *„Gemessen in slice-202"* — weder verboten noch von der Anker-Form gedeckt.

---

## Negativbefunde (geprüft, ohne Befund)

- **Die neue Regel greift, in beiden Tabellen und in beide Richtungen.** Vier Mutationen in einem Lauf: > 200 und „x" in Tabelle 1, dieselben zwei in Tabelle 2 (Spalte *Tut was*) — alle vier gemeldet, jeder Befund auf **seiner** Zeile.
- **Kein Still-Grün-Pfad bei umbenanntem Spaltenkopf.** `Vertrag` → `Contract` ⇒ `section-column-missing`. Die Regel schaltet sich nicht selbst ab.
- **Die `hint`-Zeile trägt beide Richtungen** — erscheint wörtlich bei `oversized` und `undersized`.
- **Die Schwelle 200 ist begründet, nicht kalibriert.** Der Config-Kommentar nennt sie als Größenordnung eines Satzes und schreibt hin, dass eine bestandskalibrierte Schwelle den Bestand eingefroren hätte. Die benannte Grenze steht dabei.
- **Der Mechanik-Anspruch stimmt:** `table.column` ist im Repo nicht neu (ADR-Index, vier Spalten).
- **Die conventions.md-Zahlen stimmen** — elf Titel > 200, längster 1186; `MR-054` hatte zwei fehlende Zellen. Heute: 38 Zeilen, max. Titel 192.
- **`max Zellinhalt 179` stimmt.**
- **23 Träger, alle verlinkt**; `links`/`anchors` grün.
- **Eingefrorene Zahlen sind in der größten Datei durch Kommandos ersetzt**, samt der Aussage, worüber die Zahl spricht (Ausschnitt, nicht Repo).
- **Aussagen-Treue der größten Zelle** (4997 Zeichen) Satz für Satz geprüft: kein Verlust; F-3 und F-5 aus Runde 1 sauber eingearbeitet.
- **Sechs weitere Träger** gegen ihre alten Zellen geprüft — bei `image-scan` und `workflow-pins` ist der Umzug ein Zuwachs.
- **Kein Produkt-Code berührt.**
- **Der Ausschluss der `Bindung`-Spalte ist ehrlich und gemessen** — eine verlängerte Bindung-Zelle meldet nichts, genau wie §3 ausweist.
- **Plan-Form:** neun Abschnitte, drei Vorprüfungen mit zwei `cite`-Belegen, `citations` grün; §4 benennt die Überschreitung der Größenregel ausdrücklich.
- **Die F-6-Korrektur trägt** — der Satz *„Meta-Gates sind fail-closed"* stimmt wieder.

---

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 9 |
| LOW | 8 |
| INFO | 2 |

**Wiederkehrende Klassen:** `zahl-widerspricht-ihrer-aufzaehlung` (M-3 fünffach, plus F-2 aus Runde 1 und M-1), `quelle-ueber-geltungsbereich` (M-4, L-5), `deckungsnachweis-in-der-sensor-datei` (M-9, fünf Dateien).

## Verdikt

**Blockiert.** Neun MEDIUM, davon vier mit eigenständigem Gewicht: die falsche Fund-Zuschreibung an den neuen Sensor (M-2, per Break-Test widerlegt), die verlorene Still-Grün-Zusage (M-6), das nicht durchgehaltene Nicht-Gate-Kriterium bei zwei byte-identischen Recipes (M-7) und der veraltete Vorprüfungs-Block, an dem die Register-Zählung hängt (M-8).

Zur Einordnung: Der **Kern** trägt. Die neue Regel ist begründet, fail-closed gegen ihre eigene Abschaltung, greift in beiden Tabellen und beiden Richtungen, und der Umzug der größten Zelle ist aussagen-treu. Was blockiert, ist die Schicht darüber: Zahlen und Zuschreibungen, die die Arbeit beschreiben, halten der Nachmessung an sechs Stellen nicht stand. Beide §5-Risiken sind damit entschieden: „Umschichtung ohne Gewinn" ist in fünf Dateien eingetreten (M-9), „Nicht-Gate-Zuordnung ist ein Urteil" in zwei Zeilen falsch gefallen (M-7) — in der gefährlicheren Richtung, die der Plan selbst benannt hat.
