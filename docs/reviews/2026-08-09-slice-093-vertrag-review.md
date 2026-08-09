# Review-Report: slice-093 (Vertrags-/Dokumentations-Seite) — 2026-08-09

**Review-Art:** Design — geprüft wird die **Vertragsseite** gegen Architektur-
und Prozess-Invarianten: Lastenheft-Zusage ↔ Spezifikations-Algorithmus ↔ ADR-
Begründung ↔ Harness-Deklarationen. Die Go-Implementierung prüft ein zweiter
Reviewer parallel; sie wird hier nur als **Beleg** herangezogen, wo eine
Vertragsaussage sonst nicht falsifizierbar wäre.

**Gegenstand:** [slice-093](../plan/planning/done/slice-093-closure-note-gate.md),
Commit-Range `18489ee..ae4cc09` (sechs Commits)

**Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) @ 1.3.0 ·
**Modell:** `claude-opus-5[1m]` · **Datum:** 2026-08-09

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne diese Liste
ist der Lauf nicht reproduzierbar):

- [slice-093](../plan/planning/done/slice-093-closure-note-gate.md) (Slice-Plan, Abnahme-Punkte 1 + 2)
- [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md) sowie der [ADR-Index](../plan/adr/README.md)
- [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in) (geschärft),
  [`DC-FA-CLI-012`](../../spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben) (neu),
  [`DC-FA-CONF-001`](../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei),
  [`DC-FA-SRC-001`](../../spec/lastenheft.md#dc-fa-src-001--upstream-content-drift-externer-quellen-modul-sources-opt-in-netz),
  [`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
- [`spec/spezifikation.md`](../../spec/spezifikation.md) §Planning (C1–C5), §Konfigurations-Pfad, §2-Schema, §4-Grund-Codes, §7-Historie
- [`AGENTS.md`](../../AGENTS.md) (Hard Rules, §4-Gate-Tabelle), [`harness/README.md`](../../harness/README.md) (§Guides, §Sensors, §Gate-Taxonomie)
- `Makefile`, [`.d-check.closure.yml`](../../.d-check.closure.yml), [`.d-check.yml`](../../.d-check.yml), `CHANGELOG.md`,
  [`.harness/skills/closure-note-reviewer.md`](../../.harness/skills/closure-note-reviewer.md)

---

## Findings

### F-1 — Zwei normative Stellen der Spezifikation widersprechen sich zur Befund-Provenance

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-SRC-001`](../../spec/lastenheft.md#dc-fa-src-001--upstream-content-drift-externer-quellen-modul-sources-opt-in-netz) / [`DC-FA-CLI-012`](../../spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben)
- `pfad`: `spec/spezifikation.md:1875`
- `befund`: §[`DC-FA-SRC-001.a`](../../spec/spezifikation.md#dc-fa-src-001a--upstream-content-drift-externer-quellen-sources)
  Schritt 5 legt für einen Config-Pin weiterhin fest, der Befund trage `file` = `.d-check.yml`;
  die neue §[`DC-FA-CLI-012.a`](../../spec/spezifikation.md#dc-fa-cli-012a--konfigurations-pfad---config)
  Schritt 6 legt für dieselbe Ausgabe die *tatsächlich geladene* Datei fest. Beide Sätze
  stehen unverändert nebeneinander im selben, technisch verbindlichen Stratum, und
  `CHANGELOG.md` kündigt die Verhaltensänderung als „Fixed" an.
- `verifizierbar`: nein — kein Gate vergleicht zwei Prosa-Stellen desselben Dokuments;
  sichtbar erst, wenn jemand die ältere Stelle als Abnahme-Orakel benutzt.
- `klasse`: Zwei normative Stellen desselben Stratums widersprechen sich

### F-2 — Leerer `--config`-Wert fällt still auf die konventionelle Datei zurück

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-CLI-012`](../../spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben)
- `pfad`: `spec/lastenheft.md:537`
- `befund`: Der Vertrag sagt zu, es gebe „**keinen** stillen Rückfall auf Defaults und
  **keinen** auf die konventionelle Datei"; für den **Leerwert** (`--config ""`, in der
  Praxis eine nicht expandierte Make-/CI-Variable) sagt weder die Beschreibung noch ein
  Akzeptanzkriterium etwas, und der Lauf verwendet dann wortlos die konventionelle Datei.
  Ein so aufgerufener Closure-Bindepunkt meldet grün, während er das Inner-Loop-Profil
  gefahren hat — genau der Fall, den der Absatz auszuschließen behauptet.
- `verifizierbar`: ja — reproduziert am eigenen Repo gegen das Runtime-Image: der Lauf mit
  leerem `--config` prüft 329 Dateien (Datei-Zahl der konventionellen Konfiguration samt
  ihrer Ignorier-Muster) statt der 300 des Closure-Profils, Exit 0, ohne Hinweis auf stderr.
- `klasse`: Leerwert einer fail-closed-Option fällt still zurück

### F-3 — Die Exit-2-Aufzählung der Closure-Konfiguration ist enger als die Validierung

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
- `pfad`: `spec/lastenheft.md:1879`
- `befund`: Lastenheft und §C1 (`spec/spezifikation.md:1656`) zählen die Exit-2-Ursachen
  **abschließend** auf (nicht kompilierendes `heading-pattern`, `min-sentences` < 1, leerer
  Floskel-Eintrag) und ordnen ein „gesetztes, aber fehlendes oder unlesbares"
  `planning.closure.dir` ausdrücklich dem Befund `closure-note-missing` (Exit 1) zu. Ein
  absolutes oder `..`-haltiges `planning.closure.dir` ist jedoch ein vierter Exit-2-Fall;
  die §2-Schema-Zeile nennt zwar „Wurzel-relativ, innerhalb der Repo-Wurzel", benennt aber
  — anders als die Schwester-Zeile zu `codepaths.roots` — keine Folge.
- `verifizierbar`: ja — eine Konfiguration mit absolutem `planning.closure.dir` liefert
  Exit 2, die Vertragslesart erwartet Exit 1 mit Befund.
- `klasse`: Geschlossene Fehler-Aufzählung enger als die Validierung

### F-4 — Fence-Behandlung nur für die halbe Messmethode spezifiziert

- `kategorie`: MEDIUM
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
- `pfad`: `spec/spezifikation.md:1667`
- `befund`: §C4 legt ausdrücklich fest, dass die Fenced-Code-Blöcke **vor** dem Zählen
  entfernt werden; §C3 — die Bestimmung der Abschnitts-Überschrift — sagt dazu nichts und
  beschreibt schlicht „die **erste** Zeile, deren getrimmte Fassung auf
  `planning.closure.heading-pattern` passt". Eine Beispiel-Überschrift in einem
  Fenced-Block eines abgeschlossenen Slice eröffnete nach dieser Lesart den geprüften
  Abschnitt; kein Akzeptanzkriterium schließt den Fall aus. (Die aktuelle Implementierung
  überspringt Fences auch in C3 — die Zusage bleibt trotzdem hinter dem Verhalten zurück,
  und ein Refactoring, das sich am Vertrag orientiert, dürfte sie entfernen.)
- `verifizierbar`: ja — ein Slice mit einer eingezäunten Beispiel-Überschrift vor der
  echten Notiz unterscheidet die beiden Lesarten im Befundsatz.
- `klasse`: Fence-Behandlung nur für halbe Messmethode spezifiziert

### F-5 — Das Gate hat keine Untergrenze der geprüften Menge (latenter Silent-Green-Pfad)

- `kategorie`: MEDIUM
- `quelle`: [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md) §Konsequenzen / [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
- `pfad`: `spec/spezifikation.md:1660`
- `befund`: §C2 erklärt ein **leeres** Closure-Verzeichnis ausdrücklich für befundfrei, und
  die Kandidatenmenge entsteht aus einem Basisnamen-Glob. Existiert das konfigurierte
  Verzeichnis, enthält aber keinen Treffer — etwa nach einer Umgliederung des Bestands in
  Wellen-Unterverzeichnisse —, meldet `make verify-closure-notes` grün, ohne eine einzige
  Notiz gelesen zu haben. Die Zusammenfassungszeile trägt die Zahl der **gescannten**
  Markdown-Dateien (im Slice-Plan als „299 Dateien / 0 Befunde" notiert, heute 300), nicht
  die Zahl der geprüften Closure-Notizen (92); der Unterschied ist aus der Gate-Ausgabe
  nicht erkennbar. Der Selbsttest des Profils verriegelt nur, dass `planning.closure.dir`
  **gesetzt** ist, nicht dass es Kandidaten liefert.
- `verifizierbar`: ja — `planning.closure.dir` im Closure-Profil auf ein existierendes
  Verzeichnis ohne `slice-*.md` zeigen lassen: das Gate bleibt grün.
- `klasse`: Gate ohne Untergrenze der geprüften Menge

### F-6 — Die Historie-Zeile der Spezifikation widerspricht dem Dokument-Körper

- `kategorie`: LOW
- `quelle`: [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
- `pfad`: `spec/spezifikation.md:2229`
- `befund`: Die §7-Zeile vom 2026-08-09 sagt, die drei Grund-Codes „folgen mit der
  Implementierung (AllReasons-↔-§4-Lockstep)". Im selben Dokument führt die §4-Tabelle sie
  seit dem Implementierungs-Commit bereits (`spec/spezifikation.md:2199`–`2201`); eine
  zweite, spätere Historie-Zeile existiert nicht. Wer die Chronik als Stand-Aussage liest,
  bekommt das Gegenteil des Dokument-Körpers.
- `verifizierbar`: nein — kein Gate liest die Historie gegen §4; belegbar allein durch
  Vergleich beider Abschnitte im selben Commit.
- `klasse`: Historie-Zeile widerspricht dem Dokument-Körper

### F-7 — Zwei implementierte und getestete Randfälle ohne Akzeptanzkriterium

- `kategorie`: LOW
- `quelle`: [`DC-FA-CLI-012`](../../spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben)
- `pfad`: `spec/lastenheft.md:551`
- `befund`: Der Akzeptanzkriterien-Satz deckt fünf Fälle (Happy, ohne Flag, fehlend,
  außerhalb der Wurzel, ungültiges Schema). Zwei weitere Verhaltensweisen sind zugesagt
  bzw. behauptet und getestet, aber ohne Kriterium: `CHANGELOG.md` verspricht Exit 2 auch,
  wenn der Pfad auf ein **Verzeichnis** zeigt (der Lastenheft-Text kennt nur „Fehlt sie"),
  und ein **absoluter Pfad innerhalb** der Wurzel wird akzeptiert, obwohl die Beschreibung
  nur „relativ zur Scan-Wurzel aufgelöst" sagt und §C1 dazu lediglich die vage Formel „wird
  auf die Wurzel bezogen geprüft" führt.
- `verifizierbar`: ja — beide Fälle sind heute als Go-Test vorhanden, aber an kein
  Lastenheft-Kriterium gebunden; ein Verhaltenswechsel bräche keinen Vertrag.
- `klasse`: Getestetes Verhalten ohne Akzeptanzkriterium

### F-8 — Formfehler in einer Tabellenzeile eines gleich einzufrierenden Artefakts

- `kategorie`: LOW
- `quelle`: [`AGENTS.md` §3.5](../../AGENTS.md#35-adrs-sind-nach-accepted-immutable)
- `pfad`: `docs/plan/adr/0048-closure-note-struktur-im-planning-modul.md:120`
- `befund`: Die zweite Zeile der Alternativen-Tabelle beginnt ohne führendes Pipe-Zeichen,
  alle übrigen Zeilen tragen es. Sobald der Status auf `Accepted` wechselt, ist der
  ADR-Körper per `make adr-check` eingefroren (erlaubt bleiben nur `## Geschichte`-Anhänge
  und der Status-Übergang) — die Korrektur wäre danach nur noch über eine Folge-ADR
  erreichbar.
- `verifizierbar`: ja — nach dem Status-Übergang meldet `make adr-check` jede
  Körper-Änderung dieser Datei als `core-drift-vcs`.
- `klasse`: Formfehler in einem gleich immutablen Artefakt

### F-9 — Die Aufzählung der Durchsetzungsgrenzen kennt weiterhin nur einen Closure-Bindepunkt

- `kategorie`: LOW
- `quelle`: Maintainability (Harness-Ehrlichkeit)
- `pfad`: `harness/README.md:111`
- `befund`: Die Tabelle §Gate-Taxonomie nennt seit diesem Slice zwei Closure-Bindepunkte;
  die darunter stehende Liste „ihre Kraft ist real begrenzt" führt weiterhin nur
  `completeness-check` als das Closure-Gate, das die CI bewusst nicht erzwingt. Wer diese
  Liste als Inventar der nicht CI-erzwungenen Prüfungen liest — genau ihr erklärter Zweck
  —, zählt eine Lücke zu wenig.
- `verifizierbar`: nein — `make gate-consistency` prüft nur die Target-Tabellenzeilen
  gegen das Makefile, nicht die Prosa darunter.
- `klasse`: Aufzählung der Durchsetzungslücken unvollständig

### F-10 — Entscheidungs-Status hinter der bereits verdrahteten Umsetzung

- `kategorie`: INFO
- `quelle`: [`AGENTS.md` §3.5](../../AGENTS.md#35-adrs-sind-nach-accepted-immutable) / [ADR-Index](../plan/adr/README.md)
- `pfad`: `docs/plan/adr/0048-closure-note-struktur-im-planning-modul.md:3`
- `befund`: Die ADR steht auf `Proposed` und ist damit der einzige nicht angenommene
  Eintrag des Index, während die Entscheidung bereits als Produkt-Code und als
  `fullbuild`-Bindepunkt fährt. Solange der Status nicht wechselt, greift das
  Immutabilitäts-Gate nicht: die Begründung, gegen die gebaut wurde, ist nachträglich
  unbemerkt änderbar.
- `verifizierbar`: ja — `make adr-check` lässt Körper-Änderungen an dieser Datei heute
  passieren und würde sie nach dem Status-Übergang blockieren.
- `klasse`: Entscheidungs-Status hinter der Umsetzung

### F-11 — Enumerations-Currency der Operations-Referenz noch offen

- `kategorie`: INFO
- `quelle`: [`docs/user/releasing.md`](../user/releasing.md) §Release-Prep Punkt 4
- `pfad`: `docs/user/operations.md:19`
- `befund`: Die Optionen-Tabelle der Operations-Referenz führt `--config` nicht. Die
  Release-Prep-Checkliste verlangt genau diese Ergänzung „bei einer neuen CLI-Option" und
  hält fest, dass kein Gate die Enumerationen prüft (Präzedenzfall: die Modul-Liste stand
  über zwölf Minor-Versionen still). Der Slice-Plan hat den Release-Punkt noch offen; der
  Befund ist damit terminiert, nicht versäumt — er wird hier nur festgehalten, weil ihn
  nichts Maschinelles auffangen würde.
- `verifizierbar`: nein — für diese Datei existiert keine Deklarations-Konsistenz-Prüfung
  (die Autoritäts-Tabelle des Moduls `targets` umfasst sie nicht).
- `klasse`: Enumerations-Currency ohne Gate

## Negativbefunde

- geprüft, ohne Befund: Zahlen- und Mengen-Behauptungen des Slice-Plans und der ADR —
  92 abgeschlossene Slices sind nachgezählt (92 Dateien im Closure-Verzeichnis), die
  Aussage „92/92 tragen einen Abschnitt" ist mit dem gemeldeten Null-Befund-Lauf
  konsistent, und die Schwellen-Begründung (Bestands-Minimum 5 ≥ Schwelle 4 ⇒ kein
  Retrofit) trägt.
- geprüft, ohne Befund: die Floskel-Liste des Closure-Profils — alle fünf Phrasen haben
  im abgeschlossenen Bestand null Treffer (case-insensitiv nachgezählt); die im Kommentar
  begründete Nicht-Aufnahme zweier treffender Phrasen ist damit belegt statt behauptet.
- geprüft, ohne Befund: Referenz-Richtung der Spec-Straten — weder die neuen
  Lastenheft-Abschnitte noch die neuen Spezifikations-Abschnitte noch die beiden neuen
  §7-/Historie-Zeilen nennen ADR-, Wellen-, Slice- oder Commit-Kennungen; die Begründung
  wird durchgängig als „begleitende ADR" ohne Kennung referenziert.
- geprüft, ohne Befund: das `Schärft:`-Feld der ADR zeigt aufwärts ausschließlich auf zwei
  Spezifikations-Abschnitte, nicht auf das Lastenheft; die Anforderungs-Änderungen sind
  im Lastenheft entstanden, die ADR referenziert sie nur unter `Bezug:`.
- geprüft, ohne Befund: der ADR-Index-Eintrag — Zeile vorhanden, Datei-Link auflösend,
  Kurzfassung deckungsgleich mit der Entscheidung, Anforderungs-Spalte vollständig
  (beide berührten Anforderungen plus die Pfad-Notiz).
- geprüft, ohne Befund: die Alternativen-Tabelle der ADR ist inhaltlich vollständig
  (eigenes Modul, Skript-Variante der Vorlage, Mitlaufen in `gates`, Sammel-Grund-Code,
  mitgelieferte Floskel-Defaults, semantische Erkennung, Config-Merge) und benennt zu
  jeder Alternative einen prüfbaren Verwerfungsgrund; die Konsequenzen benennen auch die
  unangenehmen (zwei Pflegestellen, Boden statt Decke, Slice-Inhalte werden jetzt gelesen).
- geprüft, ohne Befund: Deklarations-Konsistenz des neuen Targets zwischen `Makefile`,
  [`AGENTS.md`](../../AGENTS.md) §4 und [`harness/README.md`](../../harness/README.md)
  §Sensors — Target-Name, `fullbuild`-Zugehörigkeit, Ausschluss aus `gates`/`ci`,
  Profil-Datei und Grund-Codes stimmen an allen drei Stellen überein.
- geprüft, ohne Befund: das Closure-Profil als zweite Konfigurations-Tür — es aktiviert
  kein Netz-Modul, und ein getippter Test verriegelt beides (kein Netz-/Range-Modul,
  gesetzte Closure-Wurzel); die im Slice-Plan als Risiko benannte „ungeprüfte zweite Tür"
  ist damit adressiert.
- geprüft, ohne Befund: die Struktur-vs-Inferenz-Grenze des
  [Closure-Note-Reviewer-Skills](../../.harness/skills/closure-note-reviewer.md) — er
  verlangt den Gate-Lauf als Eingangs-Kontext, verbietet die Doppelmeldung explizit als
  Anti-Pattern, nimmt offene Slices ausdrücklich aus und beschreibt den Rückweg in die
  Floskel-Liste samt Vorher-Messung; die Kategorien überschneiden sich nicht mit den drei
  Grund-Codes.
- geprüft, ohne Befund: Out-of-Scope-Ehrlichkeit beider Anforderungen — die
  Struktur-nicht-Bedeutung-Grenze steht an drei Stellen (Beschreibung, Out-of-Scope,
  Historie), die frühere Zusage „nur Datei-Existenz, wie `codepaths`" ist korrekt
  zurückgezogen statt stehen gelassen, und Merge/Vererbung/Profil-Registry sind bei der
  neuen Option ausgeschlossen. Mehrfache `--config`-Angaben sind ebenfalls
  ausgeschlossen — der Vertrag verspricht dort also nichts, was er nicht hält.
- geprüft, ohne Befund: die §2-Schema-Tabelle gegen das beschriebene Verhalten — Typen,
  Defaults (leere Wurzel, Muster, Schwelle 4, leere Liste) und Folgen decken sich mit den
  Schritten C1–C5; die Inert-Zusage („keine Slice-Datei wird gelesen, Befundsatz
  byte-identisch") ist in beiden Straten gleichlautend.
- geprüft, ohne Befund: die §4-Grund-Code-Tabelle — drei Zeilen, Modul-Spalte `planning`,
  Beschreibungen deckungsgleich mit C2–C5 inklusive der Ausschluss-Regel
  („missing" schließt die beiden anderen aus).
- geprüft, ohne Befund: der `CHANGELOG.md`-Eintrag gegen die Zusagen — Added-Block und
  Fixed-Block beschreiben nichts, was über Lastenheft und Spezifikation hinausgeht
  (Ausnahme: der in F-7 genannte Verzeichnis-Fall).
- geprüft, ohne Befund: Beobachtungs-Register und Slice-Plan-Vorprüfungen — die während
  des Slice entstandene Beobachtung ist eingetragen, mit Beleg, Zähler und begründeter
  Scope-Abgrenzung; der Slice-Plan verweist darauf und begründet die Nicht-Aufnahme.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 5 |
| LOW | 4 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zwei normative Stellen desselben Stratums
widersprechen sich · Leerwert einer fail-closed-Option fällt still zurück ·
Geschlossene Fehler-Aufzählung enger als die Validierung · Fence-Behandlung nur für
halbe Messmethode spezifiziert · Gate ohne Untergrenze der geprüften Menge ·
Historie-Zeile widerspricht dem Dokument-Körper · Getestetes Verhalten ohne
Akzeptanzkriterium · Formfehler in einem gleich immutablen Artefakt · Aufzählung der
Durchsetzungslücken unvollständig · Entscheidungs-Status hinter der Umsetzung ·
Enumerations-Currency ohne Gate

## Verdikt

**Merge-blockierend:** ja — fünf MEDIUM. Der Vertrag ist in der Substanz tragfähig
(Zusagen belegt, Alternativen fair, Konsequenzen ehrlich, Zahlen nachgeprüft), aber
er ist an vier Stellen **enger oder weiter als das, was gebaut wurde** (F-1, F-3, F-4,
F-7) und beschreibt an zwei Stellen einen Weg, auf dem grün gemeldet wird, ohne dass
geprüft wurde (F-2, F-5). Beides ist genau die Klasse, die ein Struktur-Gate
unterlaufen würde, das gerade gegen stille Grün-Pfade gebaut wird.

F-8 und F-10 hängen zusammen und sind zeitkritisch: beide verlieren ihre Reparierbarkeit
in dem Moment, in dem der ADR-Status auf `Accepted` wechselt.

**Übergabe:** Findings gehen an den Implementer (Rückkante Review → Plan, hier
überwiegend Rückkante Review → Spec). Die **Finding-Klassen** gehen zusätzlich in die
Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein **Lauf-Beleg**
(dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) und ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat, ebenso die
Go-Implementierung der parallele Code-Review.
