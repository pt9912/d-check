# Neuerungen von `v7.0.0` — Entscheidungsvorlage

**Adressat:** der Auftraggeber. **Gegenstand:** was sich im Baseline-Bundle zwischen `v6.18.1` und `v7.0.0`
bewegt hat. **Stand:** 2026-11-27, Repo bei `5a14fe90`. **Auftrag:**
[slice-296](in-progress/slice-296-baseline-v700-neuerungen-sichtung.md) in
[wave-98](wave-98-baseline-v700.md). **Übernommen ist hiermit nichts** — jede Zeile wartet auf Ihre Entscheidung.

## So lesen

- **Klasse:** *bereits befolgt* (Beleg im Repo) · *neu und anwendbar* (kleinste Antwort mit Preis) ·
  *nicht anwendbar* (Grund). In Frage A heißen die Antworten *befolgt* · *Lücke* · *nicht anwendbar*.
- **Empfehlung:** *übernehmen* · *nicht übernehmen* · *Eintrag anlegen* (die Abweichung als Adaptions-Eintrag
  benennen) · *folgt X* (dort entscheiden, hier nichts) · *nichts tun*.
- **Begriffe:** *Sicht* ist `spec/architecture.md`; *Check* der Tool-/Versions-Check aus A6; *Pin* eine
  festgeschriebene Version oder ein Image-Digest; `pinsync` das Modul, das einen Wert an mehreren Stellen vergleicht;
  `modul-…` und `grundlagen-…` sind Dateien des vendorten Regelwerks. *Spiegel* heißt: dieselbe Regel an zweiter Stelle.
  *Wortlaut X* ist das Zitat zur Zeile X am Ende.
- **Preise** sind Dateien, Zeilen und Gate; die Zeilenzahlen sind geschätzt, nicht gebaut.
- **Liegenlassen ist nicht immer kostenlos.** Der Pin steht auf `v7.0.0`
  ([Freshness-Audit](../../../.harness/baseline/v7.0.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2)).
  Für wiederkehrende Vorlagen (ADR, Slice, Welle, MR; hier per `cp` angelegt,
  [`implement-slice.md`](../../../.claude/commands/implement-slice.md)) <!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-02-harness-bootstrap.md:282-283 -->
  „gilt die Append-only-Logik: Neue Instanzen folgen der neuen Form" — ein Hinweis in einer solchen Vorlage wirkt
  also ohne Ihr Zutun. In `AGENTS.md`, Lastenheft und anderen Einzelstücken verlangt der Kanon Nacharbeit für <!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-02-harness-bootstrap.md:279-281 -->
  „neue **Pflicht**-Felder und umbenannte Sektionen", nicht für neue optionale Felder. Eine bewusst nicht befolgte
  Regel gehört in einen Adaptions-Eintrag ([`harness/conventions.md`](../../../harness/conventions.md); eine Datei
  plus Index-Zeile, rund 25 Zeilen).

## Wo wirklich zu entscheiden ist

| Gegenstand | Zeilen | Empfehlung in einem Satz |
|---|---|---|
| Umgebungs-Zusage und Tool-/Versions-Check | A3, A6, A7 | kein Check; die Host-Lücke und den Ort der Zusage in einem Adaptions-Eintrag benennen |
| Fund: `tools/archive-wave` baut mit einem anderen Go, als der Aufruf sagt | A8 | übernehmen |
| Bedrohungen und Vertrauensgrenzen | 4a, 4b, 7a, 7b, 7c | übernehmen, mit dem Abschnitt in der Sicht zuerst |
| Arbeitsbäume im Repo | 5b | übernehmen |
| Reste der Namensregel im Produkt | B1, B2 | B1 Handbuch-Satz, B2 übernehmen |
| Begründungssatz der Docker-only-Regel in `AGENTS.md` | 10b | nicht übernehmen |

Alle übrigen Zeilen folgen einer dieser Zeilen, sind entschieden, befolgt oder verlangen keinen Akt.

## Zählung

Gezählt ist die **Datei im Bundle** (`regelwerk/` und `templates/`), alt aus
`git archive 738d2bca .harness/baseline/v6.18.1`, neu aus dem vendorten Baum. **Ohne `SHA256SUMS`**; mit dieser
Datei wären es 55/56, 48 gleichnamige, 39 byte-verschiedene und 32 inhaltlich verschiedene. „Inhaltlich“ heißt:
ohne die `<!-- Quelle:`-Zeilen, beide Versionsstempel auf einen Platzhalter gesetzt. Diff-Zeilen sind die Zeilen
von `diff`, die mit `<` oder `>` beginnen.

| Gemessen | Zahl |
|---|---|
| Dateien alt / neu | 54 / 55 |
| gleichnamig · davon byte-verschieden (`cmp`) · Diff-Zeilen roh | 47 · 38 · 862 |
| inhaltlich verschieden (19 Regelwerk, 12 Vorlagen) · Diff-Zeilen | 31 · 804 |
| nur Stempel · umbenannt · neu | 7 · 7 · 1 |

Die Zahlen decken sich mit [wave-98 §1](wave-98-baseline-v700.md). Dateinummern unten: 1 bis 39, dazu die
Sammelzeile 40; `grep -o '^| [0-9][0-9]*' <diese Datei> | sort -u | wc -l` ergibt 40, `grep -c '^| [0-9]'` zählt die
Unterzeilen mit. Reihenfolge: größte Bewegung zuerst; die Zahl hinter dem Dateinamen sind die Diff-Zeilen (bei
umbenannten Dateien gegen den alten Namen, bei der neuen Datei ihre Länge).

## Inventar

| Nr | Datei · Neuerung | Fundstelle | Klasse | Beleg · kleinste Antwort mit Preis · Grund | Empfehlung | Entscheidung des Auftraggebers |
|---|---|---|---|---|---|---|
| 1a | `modul-14` (258): Umgebungs-Zusage als Anforderung | [§Umgebungs-Zusage](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#die-umgebungs-zusage--anforderung-vor-träger) | neu und anwendbar | Die Umgebung eines Laufs ist festgelegt, und das lässt sich belegen (Wortlaut 1a). Je Bestandteil in A1 bis A5, A7. | folgt A3, A7 |  |
| 1b | `modul-14`: zweiter Träger (Tool-Pin plus Tool-/Versions-Check) mit Inventar, Skriptpflichten, Probe | [§Zweiter Träger](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#der-zweite-träger-tool-pin-und-tool-versions-check) | nicht anwendbar | Bedingung: ein Check wird geführt (Wortlaut A6). Einzelheiten in A6, A9 bis A14. | folgt A6 |  |
| 1c | `modul-14`: zwei Prüfungen derselben Zusage (Stellen-Abgleich, Laufzeit-Gate) | [§Zwei Prüfungen](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#zwei-prüfungen-derselben-zusage) | neu und anwendbar | Stellen-Abgleich: A8; Laufzeit-Gate: A9. | folgt A8 |  |
| 1d | `modul-14`: Fehlannahmen und eine Überschrift umformuliert (*Zwei Formen des Belegs beim Docker-Träger*) | [§Fehlannahmen](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#regeln-gegen-typische-fehlannahmen-modul-14) | nicht anwendbar | Kein eigener Akt. Auf den alten Anker zeigt kein lebender Link ([`MR-081`](../../../harness/conventions.md#mr-081)). | nichts tun |  |
| 2a | `templates/.d-check.yml` (74): `.claude/worktrees/**` im Scan-Ausschluss | [Vorlage](../../../.harness/baseline/v7.0.0/templates/.d-check.yml) | neu und anwendbar | Vorlagen-Seite zu 5b. | folgt 5b |  |
| 2b | `.d-check.yml`-Vorlage: `doc-manual` in der Ausnahmeliste für Fragment-Targets | [Vorlage](../../../.harness/baseline/v7.0.0/templates/.d-check.yml) | nicht anwendbar | Die Liste nimmt Targets des Fragments `d-check.mk` aus; die Vorlage setzt das Fragment erst voraus, wenn es erzeugt ist (Wortlaut 2b). Dieses Repo bindet keines ein (`DCHECK_RUN` im [`Makefile`](../../../Makefile)); `targets.exempt-targets` ist leer. | nichts tun |  |
| 2c | `.d-check.yml`-Vorlage: Klasse `wave`, `specification.md` | [Vorlage](../../../.harness/baseline/v7.0.0/templates/.d-check.yml) | nicht anwendbar | Entschieden (12a); Klasse `wave` steht in [`.d-check.yml`](../../../.d-check.yml). | nichts tun |  |
| 2d | `.d-check.yml`-Vorlage: Beispielblock `planning` mit `wave-*`-Globs und `mode: many` | [Vorlage](../../../.harness/baseline/v7.0.0/templates/.d-check.yml) | bereits befolgt | [`.d-check.yml`](../../../.d-check.yml) `planning.waves` (beide Globs `wave-*`, `mode: many`); `make planning-check`. | nichts tun |  |
| 2e | `.d-check.yml`-Vorlage: Beispielblock `pinsync` | [Vorlage](../../../.harness/baseline/v7.0.0/templates/.d-check.yml) | bereits befolgt | [`.d-check.yml`](../../../.d-check.yml) `pinsync` (zwei Pins, `make doc-check`). Eine Lücke daneben: A8. | nichts tun |  |
| 3 | `environment-check.template.md` (74, neu): Sensor-Datei des Tool-/Versions-Checks | [Vorlage](../../../.harness/baseline/v7.0.0/templates/harness/sensors/environment-check.template.md) | nicht anwendbar | Bedingung: es gibt ein eigenes Prüf-Skript (Wortlaut A15); je Abschnitt A15 bis A20. | folgt A6 |  |
| 4a | `modul-04` (57): die Vertrauensgrenze steht in der Sicht | [§Bedrohung](../../../.harness/baseline/v7.0.0/regelwerk/modul-04-adrs.md#bedrohung-grenze-ausgang-modul-4) | neu und anwendbar | Die Grenze steht an ihrer Schnittstelle in der Sicht (Wortlaut 4a). [`spec/architecture.md`](../../../spec/architecture.md) hat keinen solchen Abschnitt (`git grep -i vertrauensgrenze -- spec`: 0 Zeilen), nimmt aber fremde Eingaben an (Dateisystem-, HTTP-, VCS-, Config-Adapter). **Kleinste:** ein Abschnitt mit Tabelle in der Sicht, 1 Datei, rund 15 Zeilen, kein Gate. Größer: dazu je Grenze die Bedrohungen mit Ausgang durchgehen. | übernehmen (kleinste) |  |
| 4b | `modul-04`: jede Bedrohung bekommt einen von drei Ausgängen, kein eigenes Register | [§Bedrohung](../../../.harness/baseline/v7.0.0/regelwerk/modul-04-adrs.md#bedrohung-grenze-ausgang-modul-4) | neu und anwendbar | Gemindert: Anforderung mit Prüfung; akzeptiert: ADR; offen: Slice in `open/` (Wortlaut 4b). **Kleinste:** 0 Dateien, die Regel wirkt über den Pin; durchgesetzt würde sie erst durch 7a. | übernehmen mit 7a |  |
| 4c | `modul-04`: eine Bedrohung ist kein Risiko des Slice | [§Bedrohung](../../../.harness/baseline/v7.0.0/regelwerk/modul-04-adrs.md#bedrohung-grenze-ausgang-modul-4) | neu und anwendbar | Ihr Ausgang steht nicht unter *Risiken* eines Slice (Wortlaut 4c). **Kleinste:** 0 Dateien; die Slice-Vorlage trägt den Hinweis in jeden neuen Plan (15b). | nichts tun |  |
| 4d | `modul-04`: für einen Wächter liefert das Bedrohungsmodell nur Kandidaten | [§Bedrohung](../../../.harness/baseline/v7.0.0/regelwerk/modul-04-adrs.md#bedrohung-grenze-ausgang-modul-4) | nicht anwendbar | Bedingung: es gibt ein Bedrohungsmodell (Wortlaut 4d). Das Repo führt keines. | folgt 4a |  |
| 5a | `modul-05` (49): je parallelem Lauf ein Zweig und ein Arbeitsbaum | [§Lifecycle](../../../.harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md#lifecycle-als-state-machine) | neu und anwendbar | Zwei Läufe in derselben Arbeitskopie überschreiben sich still (Wortlaut 5a). Wirkt erst bei einem zweiten Lauf; [`implement-slice.md`](../../../.claude/commands/implement-slice.md) schließt heute ein paralleles `in-progress/` aus. **Kleinste:** 0 Dateien. Größer: Absatz *Parallele Läufe* in `AGENTS.md` §6, rund 10 Zeilen. | nichts tun |  |
| 5b | `modul-05`: Arbeitsbaum neben dem Repo; legt ein Werkzeug ihn hinein, gehört der Pfad aus dem Scan und ins `.gitignore` | [§Lifecycle](../../../.harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md#lifecycle-als-state-machine) | neu und anwendbar | Ein Sensor, der das Dateisystem liest, sieht den Arbeitsbaum sonst als zweite Kopie (Wortlaut 5b). Gemessen am Wegwerf-Baum mit derselben Scan-Wurzel wie hier: d-check meldet einen Befund aus `.claude/worktrees/`. Weder `scan.ignore` noch `.gitignore` nennen den Pfad. **Kleinste:** [`.d-check.yml`](../../../.d-check.yml) und `.gitignore`, je 1 Zeile; das Closure-Profil hat keinen `scan`-Block und wäre eine dritte Stelle. Ob der Ausschluss eine ADR braucht ([`AGENTS.md`](../../../AGENTS.md) §3.6), entscheidet der Architect. | übernehmen |  |
| 5c | `modul-05`: ein Slice gehört einem Implementer-Lauf | [§Lifecycle](../../../.harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md#lifecycle-als-state-machine) | neu und anwendbar | Wer parallel arbeiten will, zerlegt den Slice (Wortlaut 5c). **Kleinste:** 0 Dateien; wirkt über den Pin. | nichts tun |  |
| 5d | `modul-05`: Bedrohung ist kein Slice-Risiko (Spiegel) | [§Offene Risiken](../../../.harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md#offene-risiken-werden-bei-closure-aufgelöst) | neu und anwendbar | Wie 4c; der Rest der Datei ist Namens-Nachzug (12a). | folgt 4c |  |
| 6 | `templates/README.md` (35): Bundle-Index | [Vorlagen-Index](../../../.harness/baseline/v7.0.0/templates/README.md) | nicht anwendbar | Verzeichnis des Bundles (Zählung, neue Vorlage, Namen); kein Akt im Repo. | nichts tun |  |
| 7a | `modul-10` (28): Prüffrage *Berührt der Schnitt eine Vertrauensgrenze?*, fehlender Ausgang ist MEDIUM | [§Review-Arten](../../../.harness/baseline/v7.0.0/regelwerk/modul-10-review-harness.md#drei-review-arten--wogegen-wird-geprüft) | neu und anwendbar | Eine offene Frage, kein Verstoß; ein Urteil im Review, kein Sensor (Wortlaut 7a). **Kleinste:** eine Frage in [`reviewer.md`](../../../.harness/skills/reviewer.md), rund 6 Zeilen, kein Gate. Setzt 4a voraus. | übernehmen mit 4a |  |
| 7b | `modul-10`: ein akzeptierter Ausgang ist kein Verstoß | [§Finding-Kategorien](../../../.harness/baseline/v7.0.0/regelwerk/modul-10-review-harness.md#finding-kategorien) | neu und anwendbar | Was eine aktive ADR an der Grenze akzeptiert und nennt, meldet der Reviewer nicht (Wortlaut 7b). **Kleinste:** Ausnahme in [`reviewer.md`](../../../.harness/skills/reviewer.md), rund 4 Zeilen. | übernehmen mit 7a |  |
| 7c | `modul-10`: Kontext-Eingang um Sicht, ADRs der berührten Schnittstelle und offene Slices erweitert | [§Reviewer-Skill](../../../.harness/baseline/v7.0.0/regelwerk/modul-10-review-harness.md#ziel-form-reviewer-skill) | neu und anwendbar | Pflicht-Eingang des Reviewers (Wortlaut 7c). Der Eingang in [`reviewer.md`](../../../.harness/skills/reviewer.md) nennt die drei nicht. **Kleinste:** rund 3 Zeilen dort. Der Rest der Datei ist Namens-Nachzug (12a). | übernehmen mit 7a |  |
| 8 | `modul-06` (28): `wave-`, `archive.zip` | [Modul 6](../../../.harness/baseline/v7.0.0/regelwerk/modul-06-roadmap.md) | nicht anwendbar | Entschieden (12a), nur Namen. | nichts tun |  |
| 9a | `wave.template.md` (28, umbenannt): Name und Platzhalter | [Vorlage](../../../.harness/baseline/v7.0.0/templates/docs/plan/planning/wave.template.md) | nicht anwendbar | Entschieden (12a). | nichts tun |  |
| 9b | `wave.template.md`: Slice-Tabelle verlinkt die Slice-Datei samt Lifecycle-Verzeichnis | [Vorlage](../../../.harness/baseline/v7.0.0/templates/docs/plan/planning/wave.template.md) | bereits befolgt | [wave-98 §4](wave-98-baseline-v700.md) verlinkt so; [`slice-mv.sh`](../../../tools/harness/slice-mv.sh) zieht die Pfade nach, `make doc-check` prüft sie. | nichts tun |  |
| 10a | `AGENTS.template.md` (24): Host braucht mindestens Docker und `make`, den Rest nennt das Inventar eines Checks | [Vorlage](../../../.harness/baseline/v7.0.0/templates/AGENTS.template.md) | bereits befolgt | [`AGENTS.md`](../../../AGENTS.md) §3.1 nennt `git`, `make`, `bash`, Docker und die POSIX-Werkzeuge als Klasse. Der Inventar-Halbsatz gilt nur, wenn das Repo einen Tool-/Versions-Check führt (Wortlaut 10a). | nichts tun |  |
| 10b | `AGENTS.template.md`: Begründung der Docker-only-Regel nennt die Umgebungs-Zusage statt *Supply-Chain-Defense* | [Vorlage](../../../.harness/baseline/v7.0.0/templates/AGENTS.template.md) | neu und anwendbar | [`AGENTS.md`](../../../AGENTS.md) §3.1 trägt den alten Wortlaut. **Kleinste:** 1 Zeile, kein Gate. Die alte Begründung bleibt wahr. | nicht übernehmen |  |
| 10c | `AGENTS.template.md`: Absatz *Parallele Läufe* | [Vorlage](../../../.harness/baseline/v7.0.0/templates/AGENTS.template.md) | neu und anwendbar | Text zu 5a. | folgt 5a |  |
| 10d | `AGENTS.template.md`: `specification.md`, `seit wave-` | [Vorlage](../../../.harness/baseline/v7.0.0/templates/AGENTS.template.md) | nicht anwendbar | Entschieden (12a). | nichts tun |  |
| 11 | `grundlagen-traceability` (24): `seit wave-`, Ergebnisnotiz-Namen | [Traceability](../../../.harness/baseline/v7.0.0/regelwerk/grundlagen-traceability.md#herkunfts-anker) | nicht anwendbar | Entschieden (12a), nur Namen. | nichts tun |  |
| 12a | `grundlagen-harness-dateien` (23): **Dateinamen sind englisch** — die Namensregel | [§Verzeichniskonvention](../../../.harness/baseline/v7.0.0/regelwerk/grundlagen-harness-dateien.md#verzeichniskonvention) | nicht anwendbar | **Entschieden** am 2026-11-25 („ab jetzt für Neues, der Bestand bleibt“, [wave-98 §1](wave-98-baseline-v700.md)), umgesetzt in [slice-297](done/slice-297-spezifikation-heisst-specification.md) und [slice-298](done/slice-298-neue-namen-wave-archive-anker.md), [`MR-082`](../../../harness/conventions.md#mr-082). Regel: englisch, wo es ein gängiges englisches Wort gibt (Wortlaut 12a). Offene Reste nur im Produkt: Frage B. | nichts tun |  |
| 12b | `grundlagen-harness-dateien`: die Reproduzierbarkeits-Bindung setzt eine Umgebungs-Zusage voraus | [§Einstiegspunkt](../../../.harness/baseline/v7.0.0/regelwerk/grundlagen-harness-dateien.md#harnessreadmemd-als-einstiegspunkt) | neu und anwendbar | Ein Gate mit Image-Hash oder Toolchain-Pin setzt die Zusage voraus (Wortlaut 12b). [`harness/README.md`](../../../harness/README.md) führt die Bindung bei `make fullbuild`. **Kleinste:** 0 Dateien, erledigt mit dem Eintrag aus A3. | folgt A3 |  |
| 13a | `modul-03` (21): die Bedrohung gehört nicht ins Lastenheft, die Eigenschaft schon | [§Fehlannahmen](../../../.harness/baseline/v7.0.0/regelwerk/modul-03-spec.md#regeln-gegen-typische-fehlannahmen-modul-3) | neu und anwendbar | Ins Lastenheft gehört die Eigenschaft, die die Bedrohung abwehrt (Wortlaut 13a). **Kleinste:** 0 Dateien; wirkt beim nächsten Sicherheits-Kriterium. | nichts tun |  |
| 13b | `modul-03`: Vertrauensgrenzen gehören in die Sicht (Spiegel) | [§Architektur-Sicht](../../../.harness/baseline/v7.0.0/regelwerk/modul-03-spec.md#ziel-form-architektur-sicht) | neu und anwendbar | Wie 4a; der Rest ist Namens-Nachzug (12a). | folgt 4a |  |
| 14 | `architecture.template.md` (21): neuer §6 *Vertrauensgrenzen und Datenflüsse* | [Vorlage](../../../.harness/baseline/v7.0.0/templates/spec/architecture.template.md) | neu und anwendbar | Form für 4a; ihre STRIDE-Fragen sind dort ausdrücklich Denkhilfe, keine Pflicht. | folgt 4a |  |
| 15a | `slice.template.md` (17): `Verantwortlich:` nennt bei parallelen Läufen Person und Zweig | [Vorlage](../../../.harness/baseline/v7.0.0/templates/docs/plan/planning/slice.template.md) | neu und anwendbar | Bedienhinweis, reist mit jedem neuen Plan (`cp`); 0 Dateien. Der Rest ist Namens-Nachzug (12a). | folgt 5a |  |
| 15b | `slice.template.md`: eine Bedrohung gehört nicht unter *Risiken* | [Vorlage](../../../.harness/baseline/v7.0.0/templates/docs/plan/planning/slice.template.md) | neu und anwendbar | Bedienhinweis, reist mit jedem neuen Plan; 0 Dateien. | folgt 4c |  |
| 16 | `wave-results.template.md` (16, umbenannt) | [Vorlage](../../../.harness/baseline/v7.0.0/templates/docs/plan/planning/wave-results.template.md) | nicht anwendbar | Entschieden (12a), nur Namen. | nichts tun |  |
| 17a | `reviewer.template.md` (15): Kontext-Eingang erweitert | [Vorlage](../../../.harness/baseline/v7.0.0/templates/.harness/skills/reviewer.template.md) | neu und anwendbar | Vorlagen-Form zu 7c; der eigene Skill [`reviewer.md`](../../../.harness/skills/reviewer.md) ändert sich nur über jene Zeile. | folgt 7c |  |
| 17b | `reviewer.template.md`: Ausnahme beim Sicherheits-Anti-Pattern | [Vorlage](../../../.harness/baseline/v7.0.0/templates/.harness/skills/reviewer.template.md) | neu und anwendbar | Vorlagen-Form zu 7b. | folgt 7b |  |
| 17c | `reviewer.template.md`: MEDIUM-Fall ohne Ausgang | [Vorlage](../../../.harness/baseline/v7.0.0/templates/.harness/skills/reviewer.template.md) | neu und anwendbar | Vorlagen-Form zu 7a. | folgt 7a |  |
| 18a | `modul-13` (14): Tool-/Versions-Check in der Tabelle der Gate-Typen | [§Gate-Typ](../../../.harness/baseline/v7.0.0/regelwerk/modul-13-quality-gates.md#gate-typ--fehlerbild) | nicht anwendbar | Beschreibt einen Gate-Typ, den das Repo nicht führt; Bedingung wie A6. | folgt A6 |  |
| 18b | `modul-13`: Bedrohungsmodell speist Wächter nur mit Kandidaten (Spiegel) | [§Guard-Härtung](../../../.harness/baseline/v7.0.0/regelwerk/modul-13-quality-gates.md#guard-haertung) | nicht anwendbar | Bedingung wie 4d. | folgt 4a |  |
| 18c | `modul-13`: der Check ergänzt den Wächter, ersetzt ihn nicht | [§Guard-Härtung](../../../.harness/baseline/v7.0.0/regelwerk/modul-13-quality-gates.md#guard-haertung) | nicht anwendbar | Bedingung: es gibt einen Check (Wortlaut 18c). Ohne Check ist nichts abzugrenzen. Der Rest ist Namens-Nachzug (12a). | folgt A6 |  |
| 19 | `modul-08` (14): Rolleninhaber hängt an Zweig und Arbeitsbaum | [§Rollen-Regeln](../../../.harness/baseline/v7.0.0/regelwerk/modul-08-agentenrollen.md#rollen-regeln-modul-8) | neu und anwendbar | Spiegel zu 5a; der Rest ist Namens-Nachzug (12a). | folgt 5a |  |
| 20 | `modul-02` (14): Datei- und Vorlagen-Namen | [Modul 2](../../../.harness/baseline/v7.0.0/regelwerk/modul-02-harness-bootstrap.md) | nicht anwendbar | Entschieden (12a), nur Namen. | nichts tun |  |
| 21 | `grundlagen-source-precedence` (14): `specification.md`, `<wave-name>` | [Source Precedence](../../../.harness/baseline/v7.0.0/regelwerk/grundlagen-source-precedence.md) | nicht anwendbar | Entschieden (12a), nur Namen. | nichts tun |  |
| 22 | `lastenheft.template.md` (13): Sicherheits-Anforderung als prüfbare Eigenschaft; `specification.md` | [Vorlage](../../../.harness/baseline/v7.0.0/templates/spec/lastenheft.template.md) | neu und anwendbar | Form zu 13a; das Lastenheft ist ein Einzelstück ohne neues Pflichtfeld. 0 Dateien. | folgt 13a |  |
| 23 | `grundlagen-referenz-richtung` (12): `specification.md` | [Referenz-Richtung](../../../.harness/baseline/v7.0.0/regelwerk/grundlagen-referenz-richtung.md) | nicht anwendbar | Entschieden (12a), nur Namen. | nichts tun |  |
| 24 | `roadmap.template.md` (10) | [Vorlage](../../../.harness/baseline/v7.0.0/templates/docs/plan/planning/roadmap.template.md) | nicht anwendbar | Entschieden (12a), nur Platzhalter. | nichts tun |  |
| 25 | `archive-stub-slice.template.md` (10, umbenannt) | [Vorlage](../../../.harness/baseline/v7.0.0/templates/docs/plan/planning/archive-stub-slice.template.md) | nicht anwendbar | Entschieden (12a), nur Namen. | nichts tun |  |
| 26 | `archive-stub-wave.template.md` (10, umbenannt) | [Vorlage](../../../.harness/baseline/v7.0.0/templates/docs/plan/planning/archive-stub-wave.template.md) | nicht anwendbar | Entschieden (12a), nur Namen. | nichts tun |  |
| 27a | `NNNN-title.template.md` (7, umbenannt): Name | [Vorlage](../../../.harness/baseline/v7.0.0/templates/docs/plan/adr/NNNN-title.template.md) | nicht anwendbar | Entschieden (12a). | nichts tun |  |
| 27b | ADR-Vorlage: eine ADR, die eine Bedrohung akzeptiert, nennt Wortlaut, Grenze und Grund | [Vorlage](../../../.harness/baseline/v7.0.0/templates/docs/plan/adr/NNNN-title.template.md) | neu und anwendbar | Bedienhinweis, reist mit jeder neuen ADR aus der Vorlage; 0 Dateien. | folgt 4b |  |
| 28 | `modul-12` (6): der Tool-/Versions-Check als Vorstufe der Drift-Diagnose | [§Regeln](../../../.harness/baseline/v7.0.0/regelwerk/modul-12-replay-evaluierung.md#regeln) | nicht anwendbar | Bedingung: es gibt einen Check (Wortlaut 28). Es gibt keinen. | folgt A6 |  |
| 29 | `planning/README.template.md` (6) | [Vorlage](../../../.harness/baseline/v7.0.0/templates/docs/plan/planning/README.template.md) | nicht anwendbar | Entschieden (12a), nur Namen. | nichts tun |  |
| 30 | `conventions.template.md` (4) | [Vorlage](../../../.harness/baseline/v7.0.0/templates/harness/conventions.template.md) | nicht anwendbar | Entschieden (12a), nur Namen. | nichts tun |  |
| 31 | `modul-16` (4): Container als Träger der Umgebungs-Zusage; `done/wave-…-closure.md` | [§Fehlannahmen](../../../.harness/baseline/v7.0.0/regelwerk/modul-16-produktiver-betrieb.md#regeln-gegen-typische-fehlannahmen-modul-16) | nicht anwendbar | Umformulierung ohne neue Pflicht (Wortlaut 31) und ein Name (12a); dieses Repo ist nach [`harness/README.md`](../../../harness/README.md) kein produktiver Service. | nichts tun |  |
| 32 | `modul-09` (4): Begründung der Docker-only-Regel; `seit wave-` | [§Hard Rules](../../../.harness/baseline/v7.0.0/regelwerk/modul-09-implementierung.md#hard-rules-repo-spezifisch) | neu und anwendbar | Spiegel zu 10b. | folgt 10b |  |
| 33 | `modul-07` (4): `<wave-id>`, `done/wave-` | [Modul 7](../../../.harness/baseline/v7.0.0/regelwerk/modul-07-carveouts.md) | nicht anwendbar | Entschieden (12a), nur Namen. | nichts tun |  |
| 34 | Harness-`README.template.md` (4) | [Vorlage](../../../.harness/baseline/v7.0.0/templates/harness/README.template.md) | nicht anwendbar | Entschieden (12a), nur Namen. | nichts tun |  |
| 35 | `observation.template.md` (4) | [Vorlage](../../../.harness/baseline/v7.0.0/templates/docs/plan/planning/observation.template.md) | nicht anwendbar | Entschieden (12a), nur Namen. | nichts tun |  |
| 36 | `regelwerk/README.md` (2): Kurs-Welle und Datum | [Index](../../../.harness/baseline/v7.0.0/regelwerk/README.md) | nicht anwendbar | Bundle-Stand; mit dem Pin erledigt ([`MR-081`](../../../harness/conventions.md#mr-081)). | nichts tun |  |
| 37 | `MR-NNN-title.template.md` (2, umbenannt) | [Vorlage](../../../.harness/baseline/v7.0.0/templates/harness/conventions/MR-NNN-title.template.md) | nicht anwendbar | Entschieden (12a), nur Namen. | nichts tun |  |
| 38 | `specification.template.md` (2, umbenannt) | [Vorlage](../../../.harness/baseline/v7.0.0/templates/spec/specification.template.md) | nicht anwendbar | Entschieden (12a), nur Namen. | nichts tun |  |
| 39 | `grundlagen-begriffe` (1): Glossarzeile *Umgebungs-Zusage* | [§Kernbegriffe](../../../.harness/baseline/v7.0.0/regelwerk/grundlagen-begriffe.md#kernbegriffe) | nicht anwendbar | Begriffsbestimmung ohne eigenen Akt; die Regel steht in 1a. | nichts tun |  |
| 40 | **Sammelzeile**, 7 Dateien nur mit neuem Stempel: `grundlagen-bootstrap`, `grundlagen-durchsetzungsschicht`, `grundlagen-klassifikation`, `modul-00`, `modul-01`, `modul-11`, `modul-15` | [Index](../../../.harness/baseline/v7.0.0/regelwerk/README.md) | nicht anwendbar | Je 2 rohe Diff-Zeilen, die `<!-- Quelle:`-Zeile; ohne sie 0. | nichts tun |  |

## Frage A — Umgebungs-Zusage und `environment-check` gegen das Vorhandene

Kurz: Für alles **im Container** ist die Zusage erfüllt, für den **Host** nicht. Einen Tool-/Versions-Check führt das
Repo nicht. Das Kapitel verlangt ihn von einem Docker-Repo nicht wörtlich, nennt ihn aber dessen Prüfung (A6).

| Nr | Bestandteil | Fundstelle | Antwort | Beleg · kleinste Antwort mit Preis · Grund | Empfehlung | Entscheidung des Auftraggebers |
|---|---|---|---|---|---|---|
| A1 | Zusage, Teil *festgelegt* — Container-Seite | [§Umgebungs-Zusage](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#die-umgebungs-zusage--anforderung-vor-träger) | befolgt | Jede `FROM`-Zeile im [`Dockerfile`](../../../Dockerfile) und die Scanner-Images tragen einen Digest ([`ADR-0011`](../adr/0011-digest-pins-build-gate-images.md)); `make workflow-pins` hält die Action-SHAs. Eine Abweichung zwischen Tag und Digest: A8; Host und CI-Runner: A7. | nichts tun |  |
| A2 | Zusage, Teil *adressierbar* | [§Umgebungs-Zusage](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#die-umgebungs-zusage--anforderung-vor-träger) | befolgt | Die Digests stehen im Commit; `make versions` gibt die Pins aus, `make fullbuild` endet mit dem Image-Hash. Grenze: `make gates` und die CI halten keine Host-Version fest (A7). | nichts tun |  |
| A3 | Die Zusage steht als Qualitäts-Anforderung im Lastenheft | [§Umgebungs-Zusage](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#die-umgebungs-zusage--anforderung-vor-träger) | Lücke | Das Kapitel verlangt sie als Qualitäts-Anforderung, mit dem Sensor als Messmethode (Wortlaut A3). Keine der vier Qualitäts-Anforderungen in [`spec/lastenheft.md`](../../../spec/lastenheft.md) sagt sie zu. **Kleinste:** ein Adaptions-Eintrag *die Zusage trägt hier die Digest-Pin-Entscheidung, nicht das Lastenheft* (1 Datei, 1 Index-Zeile). Größer: neue Anforderung im Lastenheft (Version, Historie, Nachweis, sonst meldet `make completeness-check` eine Waise); das Lastenheft beschreibt das Produkt, nicht dieses Harness. | nicht übernehmen, Eintrag anlegen |  |
| A4 | Ordnung der Pins nach Frage, letzte Zeile *Hält der Pin noch?* | [§Umgebungs-Zusage](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#die-umgebungs-zusage--anforderung-vor-träger) | nicht anwendbar | Ordnungstabelle ohne eigenen Akt: Ihre letzte Zeile ist kein weiterer Träger (Wortlaut A4). Die Prüfungen: A8, A9. | nichts tun |  |
| A5 | Docker als Träger (Digest, Lock-File, Stages, Beleg) | [§Zwei Formen des Belegs](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#zwei-formen-des-belegs-beim-docker-träger) | befolgt | [`Dockerfile`](../../../Dockerfile) mit Stages und Digests, `go.sum`, `make versions`, Image-Hash am Ende von `make fullbuild`. Eine Datei `image-hash.txt`, wie das Kapitel sie als Beispiel nennt, führt das Repo nicht. | nichts tun |  |
| A6 | Gilt der zweite Träger **neben** dem Docker-Träger? | [§Umgebungs-Zusage](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#die-umgebungs-zusage--anforderung-vor-träger) · [§Zweiter Träger](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#der-zweite-träger-tool-pin-und-tool-versions-check) | nicht anwendbar (als Pflicht) | Das Kapitel nennt zwei Wege, Docker **oder** Tool-Pin mit Check. Zugleich: bei Docker ist der Check die Prüfung der Zusage (Wortlaut A6). Lesart dieser Vorlage: wählbar; wer die Zusage geprüft haben will, braucht ihn. **Kleinste Antwort bei Übernahme:** Skript, Inventar-Datei, Target als erstes in `make gates`, Sensor-Datei, Index-Zeile, Adaptions-Eintrag, Festlegung in der Spezifikation (rund 8 Dateien, 200 bis 250 Zeilen), Break-Test, ADR für das neue Gate ([`AGENTS.md`](../../../AGENTS.md) §3.6). | nicht übernehmen |  |
| A7 | Host-Werkzeuge und CI-Runner | [§Zweiter Träger](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#der-zweite-träger-tool-pin-und-tool-versions-check) | Lücke | Zum Host sagt der Digest nichts aus (Wortlaut A7). [`AGENTS.md`](../../../AGENTS.md) §3.1 nennt die Host-Werkzeuge als Klasse ohne Version; die Workflows laufen auf `ubuntu-latest`; Docker am Host ist ungepinnt. **Kleinste:** die Lücke im Eintrag aus A3 mitbenennen (rund 10 Zeilen, kein Gate) — schließt sie nicht. Größer: A6. | Eintrag anlegen |  |
| A8 | Stellen-Abgleich (statisch) gegen `pinsync` | [§Zwei Prüfungen](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#zwei-prüfungen-derselben-zusage) | Lücke | Derselbe Pin steht oft an mehreren Stellen (Wortlaut A8). `pinsync` hält Go- und Lint-Version zwischen `Dockerfile` und `Makefile`, **nicht** `tools/archive-wave/`. **Fund:** `make archive-wave-test` übergibt Go 1.27.4, die `FROM`-Zeile dort trägt aber den Digest von 1.27.0, und der Digest gewinnt (`printenv GOLANG_VERSION` auf `golang:1.27.4@sha256:4013…`: `1.27.0`). `make archive-wave` übergibt dieselbe Version; `make go-base-digest` liest nur das `Dockerfile` der Wurzel. **Kleinste:** `tools/archive-wave/` auf den Stand der Wurzel heben (`Dockerfile`, `Makefile`; 3 Zeilen) und in [`.d-check.yml`](../../../.d-check.yml) einen Pin über beide Basis-Digests (rund 6 Zeilen; an Kopien der beiden Dockerfiles probiert: heute `pinsync-mismatch`). Grenze: Der Pin hält die zwei Digests gleich, nicht den Tag zum Digest; ob das Werkzeug mit 1.27.4 grün testet, ist nicht gemessen. Gate: `make doc-check`; ob der neue Pin eine ADR braucht, entscheidet der Architect. Größer: A9. | übernehmen |  |
| A9 | Laufzeit-Gate (dynamisch) gegen `make freshness-*` und `make *-digest` | [§Zwei Prüfungen](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#zwei-prüfungen-derselben-zusage) | nicht anwendbar | Teil des Checks aus A6: er prüft die installierten Werkzeuge gegen den Pin (Wortlaut A9). Die vorhandenen Targets ersetzen es nicht: Sie fragen, ob upstream neuer ist, nicht, ob das Installierte dem Pin gleicht. | folgt A6 |  |
| A10 | Woher das Soll kommt (ableiten, Manifest, Mischform) | [§Soll](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#woher-das-soll-kommt) | nicht anwendbar | Bedingung: es gibt ein Gate, das Ist gegen Soll vergleicht (Wortlaut A10). Es gibt keines. Bei Übernahme: Wahl als Adaptions-Eintrag. | folgt A6 |  |
| A11 | Inventar je Werkzeug, vier Prüfformen, Vergleichsregeln | [§Prüfform](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#prüfform-je-werkzeug--kein-werkzeug-bleibt-ungeprüft) | nicht anwendbar | Regel des Checks aus A6. Heute nennt das Repo Klassen statt eines Inventars (A7); dazu sagt das Kapitel: eine Pauschale für ungeprüfte Werkzeuge gibt es nicht (Wortlaut A11). | folgt A6 |  |
| A12 | Rot und Warnung, Regel in der Spezifikation, Ausführungsort | [§Inventar](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#inventar-stufen-kontext) | nicht anwendbar | Bedingung: es gibt einen Check (Wortlaut A12). Ohne Check kein Ort zu deklarieren. | folgt A6 |  |
| A13 | Check als erster in `make gates`, fünf Pflichten des Skripts samt Break-Test | [§Wo der Check läuft](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#wo-der-check-läuft--und-was-das-skript-leisten-muss) | nicht anwendbar | Bedingung: es gibt ein Prüf-Skript des Repos (Wortlaut A13). Es gibt keines. | folgt A6 |  |
| A14 | Probe im leeren Kontext (`PATH`-Käfig) gegen `make guard-probe` | [§Probe](../../../.harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md#ist-das-inventar-vollständig-die-probe-im-leeren-kontext) | nicht anwendbar | Sie prüft die Vollständigkeit eines Inventars, das es nicht gibt (Wortlaut A14). `make guard-probe` prüft etwas anderes: den Tool-Call-Wächter gegen seine Fälle. | folgt A6 |  |
| A15 | Vorlage, §Vertrag | [Vorlage](../../../.harness/baseline/v7.0.0/templates/harness/sensors/environment-check.template.md) | nicht anwendbar | Gilt für alle Abschnitte A15 bis A20: die Vorlage liefert nur das Muster, das Skript schreibt das Repo (Wortlaut A15). Bei Übernahme zu füllen: Festlegung in der Spezifikation, Anforderung (A3). | folgt A6 |  |
| A16 | Vorlage, §Inventar und Soll | [Vorlage](../../../.harness/baseline/v7.0.0/templates/harness/sensors/environment-check.template.md) | nicht anwendbar | Bei Übernahme: Pfad des Inventars, Soll-Quelle, Ausführungsort. | folgt A6 |  |
| A17 | Vorlage, §Grenze | [Vorlage](../../../.harness/baseline/v7.0.0/templates/harness/sensors/environment-check.template.md) | nicht anwendbar | Bei Übernahme: Drift zwischen Läufen, Werkzeuge außerhalb des Inventars, ungeprüfte Kontexte. | folgt A6 |  |
| A18 | Vorlage, §Ausgabe und Ausgänge | [Vorlage](../../../.harness/baseline/v7.0.0/templates/harness/sensors/environment-check.template.md) | nicht anwendbar | Bei Übernahme: Exit-Codes, je Werkzeug Soll und Ist. | folgt A6 |  |
| A19 | Vorlage, §Sperren | [Vorlage](../../../.harness/baseline/v7.0.0/templates/harness/sensors/environment-check.template.md) | nicht anwendbar | Bei Übernahme: Sperre für eine zu alte Bash. | folgt A6 |  |
| A20 | Vorlage, §Bindung | [Vorlage](../../../.harness/baseline/v7.0.0/templates/harness/sensors/environment-check.template.md) | nicht anwendbar | Bei Übernahme: die Kennungen und der Inventar-Pfad. | folgt A6 |  |

## Frage B — Reste der Namensregel im Angebot des Produkts

Die Namensregel selbst ist entschieden (12a). Offen ist, was **Konsumenten** des Produkts sehen. „Zeile“ heißt: eine
Zeile in einer git-getrackten Datei, die die genannte Zeichenkette trägt. Gemessen mit dem Image aus `make build`
am Stand oben, an Wegwerf-Bäumen.

| Nr | Rest · gezählte Reichweite | Fundstelle | Preis der Übernahme | Preis des Bleibens | Empfehlung | Entscheidung des Auftraggebers |
|---|---|---|---|---|---|---|
| B1 | **Defaults** `welle-*.md` / `welle-*-results.md`: 10 Zeilen in 4 Dateien (`config.go` 481, 490; `configyaml.go` 1567, 1572; `config_template.go` 168, 169; Spezifikation 2384, 2385, 4050, 4051), dazu `welle-<n>` in Lastenheft 2767 und Spezifikation 2389 (`git grep -n 'welle-\*' -- '*.go' spec docs/user ':!*_test.go'`; `git grep -n 'welle-<n>' -- spec docs/user`). Das Handbuch nennt die Default-Globs nirgends. | [Vorlage](../../../.harness/baseline/v7.0.0/templates/.d-check.yml): die Vorlage nennt beide Globs Pflicht (Wortlaut B1). | **Kleinste:** im Handbuch sagen, dass bei `wave-<NN>` beide Globs zu setzen sind (1 Datei, rund 6 Zeilen; kein Verhalten ändert sich). Größer 1: Defaults tauschen — die 12 Zeilen, in einer Wegwerf-Kopie, in der nur die 4 Code-Zeilen der Defaults getauscht waren, 21 Tests in 2 Paketen rot; Lastenheft mit Version und Historie; Release. Wer `welle-` ohne eigenen Glob führt, würde dann **still** nicht mehr geprüft (Spiegelbild der Messung rechts, nicht eigens gemessen). Größer 2: mehrere Präfixe lesen ([`ADR-0117`](../adr/0117-wellen-namen-ueber-die-konfiguration.md), Alternative A), das einzige, das beide Schreibweisen bedient. | Ein Konsument, der dem Regelwerk folgt und die Globs nicht setzt, bekommt **0 Befunde** — gemessen mit gesetztem `waves.dir` und `mode: many`: Welle `wave-99` ohne Zeiger und Register-Zeile ohne Notiz, Exit 0; mit beiden Globs `wave-drift` und `wave-results-missing`. | übernehmen (nur die kleinste) |  |
| B2 | **Beispiel-Pfade** mit `spec/spezifikation.md`: 8 Zeilen in 4 Dateien (Spezifikation 233, 234; Handbuch 372, 373; `config_template.go` 56, 57; `suggest.go` 379, 380; `git grep -n 'spec/spezifikation\.md'` über diese vier). `--print-config`, `--suggest-config` und `--manual` zeigen sie. | [Vorlage](../../../.harness/baseline/v7.0.0/templates/.d-check.yml) nennt `specification.md` | **Kleinste:** alle 8 Zeilen an allen Stellen oder an keiner ([`MR-025`](../../../harness/conventions.md#mr-025)); `make test` in einer Wegwerf-Kopie grün; wirkt mit dem nächsten Release. Danach trifft dasselbe Schweigen einen neu erzeugten Vorschlag bei einem Repo mit dem **alten** Namen (Spiegelbild, nicht eigens gemessen). Größer: `--suggest-config` nimmt den vorhandenen Namen (neue Logik, Test, Festlegung in der Spezifikation). | `--suggest-config` schlägt einem Repo mit `spec/specification.md` den deutschen Namen vor (gemessen). Mit diesem Vorschlag wird ein Abwärtsverweis in der Spezifikation **nicht gemeldet** (0 Befunde; mit dem richtigen Namen `matrix-downward`). | übernehmen (kleinste) |  |
| B3 | **Konsument mit `wave-<NN>` und `archive.zip`:** was er sieht. Das Produkt liest keinen Archiv-Namen (`git grep -n 'archive\.zip' -- '*.go' spec docs/user ':!tools'`: 0 Zeilen); betroffen ist allein das Modul `planning`. `tools/archive-wave` ist nicht Teil des Images. | [§Verzeichniskonvention](../../../.harness/baseline/v7.0.0/regelwerk/grundlagen-harness-dateien.md#verzeichniskonvention) | Nichts über B1 hinaus. | Wie B1. Ein Konsument mit **beiden** Schreibweisen kann nur eine bewachen ([Beobachtung](observations/BEO-ALL/praefix-umstellung-schaltet-waechter-der-anderen-schreibweise-ab/observation.md)). | folgt B1 |  |

## Wortlaute

*Wortlaut X* in den Tabellen meint das Zitat unter der Zeilennummer X. Die Zitate stehen hier, weil das Modul `citations`
ein Zitat nur auf einer Zeile **nach** seiner Direktive liest, also nicht in einer Tabellenzelle; `make doc-check`
vergleicht jedes mit seiner Quelle.

**1a**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md:14-15 -->

> *Die Umgebung, in der ein Lauf stattfindet, ist die behauptete — und das lässt sich belegen.*

**2b**

<!-- d-check:cite .harness/baseline/v7.0.0/templates/.d-check.yml:40-40 -->

> d-check.mk erst ergaenzen, wenn `d-check --print-mk` gelaufen ist

**4a**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-04-adrs.md:93-93 -->

> **Die Grenze steht in der Sicht**

**4b**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-04-adrs.md:96-97 -->

> **Jede Bedrohung bekommt einen von drei Ausgängen, und jeder hat schon ein Zuhause:**

**4c**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-04-adrs.md:113-113 -->

> **Eine Bedrohung ist kein Risiko des Slices.**

**4d**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-04-adrs.md:132-132 -->

> **Für einen Wächter liefert das Modell nur Kandidaten.**

**5a**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md:62-63 -->

> jeder parallele Lauf bekommt seinen eigenen, etwa einen `git worktree`.

**5a, weiteres**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md:64-67 -->

> Zwei Läufe in derselben Arbeitskopie schreiben in dieselben Dateien, der zweite überschreibt den ersten, und kein Merge meldet etwas — die Änderung des ersten ist still weg.

**5a, Stand im Repo**

<!-- d-check:cite .claude/commands/implement-slice.md:94-94 -->

> kein paralleles `in-progress/`

**5b**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md:81-84 -->

> Legt ein Werkzeug seine Arbeitsbäume ins Repo (etwa Claude Code unter `.claude/worktrees/`), nimmt die Gate-Konfiguration den Pfad aus dem Scan; für git gehört er trotzdem ins `.gitignore`.

**5c**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md:86-86 -->

> **Ein Slice gehört einem Implementer-Lauf.**

**5c, weiteres**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-05-planning-harness.md:91-92 -->

> Wer seine Arbeit auf mehrere Arbeitsbäume verteilen will, schneidet ihn

**7a**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-10-review-harness.md:24-25 -->

> Fehlt er, ist das ein MEDIUM: eine offene Frage, kein Verstoß. Die Frage ist ein Urteil im Review, kein Sensor.

**7b**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-10-review-harness.md:40-40 -->

> **Ein akzeptierter Ausgang ist kein Verstoß.**

**7c**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-10-review-harness.md:91-94 -->

> **Kontext-Eingang (Pflicht):** Diff · `spec/lastenheft.md` · `spec/architecture.md` (Vertrauensgrenzen, `ARC-*`) · ADRs, deren ID im PR/Commit vorkommt, und aktive ADRs, deren `Schärft:` eine berührte `ARC-*` nennt · Slices in `open/`, die eine berührte `ARC-*` nennen

**10a**

<!-- d-check:cite .harness/baseline/v7.0.0/templates/AGENTS.template.md:102-103 -->

> führt das Repo einen Tool-/Versions-Check

**12a**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/grundlagen-harness-dateien.md:34-35 -->

> **Dateinamen sind englisch.** Deutsch bleibt ein Name nur, wenn das Artefakt kein gängiges englisches Wort hat

**12b**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/grundlagen-harness-dateien.md:298-299 -->

> es setzt eine Umgebungs-Zusage voraus

**13a**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-03-spec.md:37-37 -->

> Ins Lastenheft gehört die *Eigenschaft*, die sie abwehrt

**18c**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-13-quality-gates.md:324-324 -->

> Der Check ergänzt den Wächter, er ersetzt ihn nicht.

**28**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-12-replay-evaluierung.md:221-221 -->

> Rang 1 hat eine maschinelle Vorstufe: Der Tool-/Versions-Check

**31**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-16-produktiver-betrieb.md:15-15 -->

> Deployment ist *eine* Anwendung des Containers als Träger der Umgebungs-Zusage (Modul 14), nicht das Ziel.

**A3**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md:18-19 -->

> Dass ein Projekt eine Umgebungs-Zusage braucht, steht als `LH-QA-*` in seinem Lastenheft, mit dem Sensor als Messmethode

**A4**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md:37-38 -->

> Die letzte Zeile ist kein weiterer Träger, sondern die Prüfung, dass die anderen noch tragen.

**A6**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md:42-44 -->

> Die Umgebungs-Zusage lässt sich auf zwei Wegen erfüllen: mit Docker und Dockerfile, die die Umgebung aus gepinnten Eingaben bauen, oder mit Tool-Pins und einem Tool-/Versions-Check

**A6, weiteres**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md:193-194 -->

> Wo Docker festlegt, ist der Check die **Prüfung** der Zusage

**A7**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md:201-201 -->

> nichts — der Host liegt vor dem Image und ruft es auf

**A8**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md:207-207 -->

> Derselbe Pin steht oft an mehreren Stellen

**A9**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md:213-214 -->

> Die installierten Werkzeuge tragen den Pin. Das läuft im Kontext, den es prüft — ein Skript des Repos.

**A10**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md:220-220 -->

> Das Gate vergleicht Ist gegen Soll.

**A11**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md:253-254 -->

> Eine Pauschale wie „Coreutils bleiben ungeprüft“ gibt es nicht

**A12**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md:285-285 -->

> Der Check prüft, was dort läuft, wo er läuft.

**A13**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md:306-306 -->

> Das Laufzeit-Gate ist ein **Skript des Repos**

**A14**

<!-- d-check:cite .harness/baseline/v7.0.0/regelwerk/modul-14-docker-harness.md:374-374 -->

> **Ohne Probe** bleibt das Inventar handgepflegt und unbelegt.

**A15**

<!-- d-check:cite .harness/baseline/v7.0.0/templates/harness/sensors/environment-check.template.md:6-6 -->

> Vorlage liefert das **Muster, kein Skript**

**B1**

<!-- d-check:cite .harness/baseline/v7.0.0/templates/.d-check.yml:106-109 -->

> Die beiden Globs sind PFLICHT: Die Wellen-Kennung beginnt mit `wave-`

## Abweichungen von den Erwartungen des Plans

- Die Bestands-Zahlen stimmen; der Plan nennt nicht, dass sie ohne `SHA256SUMS` zählen (§Zählung).
- B1: Der Plan erwartete die Defaults auch im Handbuch — dort stehen sie nicht. Zusätzlich zu `config.go` führt
  `configyaml.go` zwei eigene Rückfall-Literale. Die Anforderung des Moduls `planning` legt den Default-Glob nicht fest;
  das tut die Schema-Tabelle der Spezifikation.
- A8 ist ein Fund an diesem Repo, keine Neuerung der Baseline. Er ist hier **nicht behoben**.
