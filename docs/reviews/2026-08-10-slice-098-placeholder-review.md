# Review-Report: slice-098 — 2026-08-10

**Review-Art:** Code — geprüft gegen Slice-Plan, Lastenheft/Spezifikation
(Vertrag) und ADR-0052 (Entscheidungen + Fitness Function).

**Gegenstand:** Diff-Range `3ed7666..b48a08b` (drei Commits: Lifecycle-Move ·
CR + ADR · Implementierung), Modul `planning`, vierte Closure-Bedingung
`closure-note-placeholder`.

**Skill:** `.harness/skills/reviewer.md` @ 1.3.0 ·
**Modell:** `claude-opus-5[1m]` · **Datum:** 2026-08-10

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan [slice-098](../plan/planning/in-progress/slice-098-closure-note-placeholder.md)
- [ADR-0052](../plan/adr/0052-platzhalter-erkennung-inline-code.md) (Proposed),
  verfeinert [ADR-0048](../plan/adr/0048-closure-note-struktur-im-planning-modul.md)
- [`DC-FA-PLAN-001`](../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  (Lastenheft 0.54.0, vierte Bedingung) und
  [`DC-FA-PLAN-001.a`](../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
  Schritt C4b, §2-Schema, §4-Grund-Code-Tabelle
- [`AGENTS.md`](../../AGENTS.md) (Hard Rules), Beobachtungs-Register
  [observations.md](../plan/planning/observations.md)

**Messmittel** (alle Läufe reproduziert, kein Host-Go): `make build`,
`make test`, `make lint`, `make verify-closure-notes`; ein aus dem Vor-Commit
gebautes Vergleichs-Image (`make build IMAGE=d-check-old` auf einem
`git archive`-Export von `3ed7666`); Fixture-Repos in einem Temp-Verzeichnis
**außerhalb** des Repos (je Fall eine Closure-Notiz mit vier Satzende-Zeichen
plus die Probe-Zeile); Mutationen in einer zweiten Temp-Kopie des Arbeitsstands,
nicht im Repo.

---

## Findings

### F-1 — Spezifikation C5 (Befund-Form) blieb auf der Fassung vor C4b stehen

- `kategorie`: MEDIUM
- `quelle`: `DC-FA-PLAN-001.a` (Spezifikation, Schritte C4b/C5)
- `pfad`: `spec/spezifikation.md` §`DC-FA-PLAN-001.a`, Schritt C5 (Absatz
  „C5. Befund-Form"), zusammen mit Schritt C4b
- `befund`: C5 sagt `line` = Zeile der Abschnitts-Überschrift, `target` =
  `planning.closure.dir` und zählt als `message`-Inhalte nur „fehlender
  Abschnitt · Ist-/Soll-Satzzahl · getroffene Floskel" auf; der neue Schritt C4b
  sagt für denselben Befund „mit der Zeile des Treffers und dem Treffer als
  **Ziel** (auf 40 Runen gekappt)". Das Produkt meldet die **Treffer-Zeile**
  (`line: 7` statt Überschrift `line: 3`), setzt `target` auf das
  Closure-Verzeichnis und schreibt den Treffer in `message` — C5 ist damit für
  den vierten Code falsch, C4b für das Ziel-Feld ebenfalls, und beide
  widersprechen sich gegenseitig. Ebenso unverändert blieb C5s
  Kombinierbarkeits-Satz („Ein Kandidat kann `closure-note-thin` **und**
  `closure-note-boilerplate` tragen … `closure-note-missing` schließt beide
  aus"), obwohl der neue Code als dritter Kombinationspartner hinzukommt und
  von `closure-note-missing` ebenfalls ausgeschlossen wird.
- `verifizierbar`: ja — ein Lauf mit `--json` auf einer Notiz mit
  Vorlagen-Rumpf zeigt `"line": 7`, `"target": "docs/done"` und den Treffer im
  `message`-Feld; gemessen gegen das aus dem Diff gebaute Image.
- `klasse`: Rand einer Semantik-Änderung bleibt auf der widerrufenen Fassung

### F-2 — Lastenheft behauptet weiter, die Substanz-Schwelle decke Platzhalter ab

- `kategorie`: MEDIUM
- `quelle`: `DC-FA-PLAN-001` (Lastenheft 0.54.0)
- `pfad`: `spec/lastenheft.md`, Bullet „**Substanz**" der Closure-Note-Struktur
  (Satz „Ein zurückgelassener Platzhalter fällt damit auf.")
- `befund`: Derselbe Anforderungs-Text trägt rund zwanzig Zeilen später die
  gemessene Gegenaussage: „Ein Template-Rumpf ist syntaktisch vollständig und
  passiert die drei Bedingungen oben … und bleibt grün." Die ältere Zusage steht
  unverändert im normativen Bullet und ist gegen die Umsetzung falsifizierbar:
  eine Notiz aus vier Platzhalter-Sätzen erreicht die Schwelle und meldet ohne
  den neuen Schalter nichts. Denselben Rand hat die Anforderung schon einmal
  gebraucht (Historie 0.53.1 zog ein doppelt geführtes Akzeptanzkriterium
  zusammen).
- `verifizierbar`: ja — der Vergleichslauf des Vor-Commit-Images auf einer
  Notiz mit vier Platzhalter-Sätzen endet mit Exit 0 und null Befunden.
- `klasse`: Rand einer Semantik-Änderung bleibt auf der widerrufenen Fassung

### F-3 — Zusage „Vergleichszeichen sind durch die Form ausgeschlossen" ist mit einer Zeile widerlegbar

- `kategorie`: MEDIUM
- `quelle`: `DC-FA-PLAN-001` (Lastenheft 0.54.0) / ADR-0052 §Fitness Function
- `pfad`: `spec/lastenheft.md`, Absatz nach den drei Einschränkungen
  („Vergleichs- und Größer-Zeichen der Messwert-Prosa … sind bereits durch die
  Form ausgeschlossen"); Erkennung in
  `internal/hexagon/core/rules/planning.go:221`
- `befund`: Ausgeschlossen ist nur die **mit Leerzeichen** geschriebene Form.
  Der Satz `Die Latenz blieb <1 s und der Recall >0,9 im Median.` erzeugt einen
  Befund `closure-note-placeholder` mit dem „Treffer" `<1 s und der Recall >` —
  dieselbe Messwert-Prosa, nur ohne Leerzeichen nach dem Vergleichszeichen, wie
  sie im Deutschen üblich ist. Dieselbe Klasse trifft Tabellenzeilen: die Zeile
  `| p95 | <1 s | Recall | >0,9 |` meldet, weil das Muster über die Zellgrenzen
  hinweg bis zum nächsten Größer-Zeichen der **Zeile** greift. Der Testfall der
  Klasse prüft ausschließlich die Leerzeichen-Form („p95 < 1 s und Recall >
  0,9"), und die Fitness Function von ADR-0052 („Technische Prosa bleibt grün:
  Vergleichszeichen …") gilt damit nur für die getestete Schreibweise.
- `verifizierbar`: ja — zwei Fixture-Notizen (Prosa-Form und Tabellen-Form)
  gegen das gebaute Image mit `placeholder: true`; beide melden.
- `klasse`: Zusage-Reichweite größer als die Messung

### F-4 — HTML-Tag-Nachfilter: Liste unvollständig **und** nicht testgehalten

- `kategorie`: MEDIUM
- `quelle`: `DC-FA-PLAN-001.a` Schritt C4b („das erste Token … ist ein bekannter
  HTML-Tag-Name") / ADR-0052 Entscheidung 2
- `pfad`: `internal/hexagon/core/rules/planning.go:225-235` (`htmlTagNames`)
- `befund`: Die Liste führt 35 Namen, nennt aber weder der Vertrag noch ein Test
  ihren Umfang. Gemessen melden gängige, nicht gelistete Tags als Platzhalter:
  `<section>`, `<meta charset="utf-8">`, `<abbr title="x">` und
  `<input value="platzhalter">` erzeugen je einen Befund — Markup, das die
  Fitness Function ausdrücklich als grün zusagt („HTML-Tag … einzeln geprüft").
  Zugleich ist die Liste mutations-offen: werden **33 der 35** Einträge
  ersatzlos gestrichen (alles außer `a` und `br`), bleibt `make test` grün. Die
  beiden gehaltenen Einträge sind genau die, die in den beiden Testfällen
  vorkommen; jeder weitere Eintrag ist unbewachte Behauptung.
- `verifizierbar`: ja — (a) Fixture-Notizen mit den vier genannten Tags gegen
  das gebaute Image; (b) Rückbau der Liste in einer Temp-Kopie des
  Arbeitsstands, danach `make test` (grün).
- `klasse`: Nachfilter-Liste ohne Vertrag und ohne Test

### F-5 — Abschnitts-Grenze nachgebaut statt geteilt, Äquivalenz unbewacht

- `kategorie`: MEDIUM
- `quelle`: Maintainability / Beobachtungs-Register BEO-003 („geteilte Lexik
  driftet an den Rändern … es braucht je Konsument eine Assertion gegen
  Wieder-Divergenz")
- `pfad`: `internal/hexagon/core/rules/planning.go:199-214`
  (`closureSectionEnd`) gegenüber `internal/hexagon/core/rules/planning.go:284-305`
  (`closureSectionProse`)
- `befund`: Beide Funktionen beantworten dieselbe Frage („wo endet der
  Closure-Abschnitt?") mit derselben, **kopierten** Schleife — Fence-Automat,
  `inFence`-Übersprung, ATX-Ebenenvergleich. Heute ziehen sie dieselbe Grenze
  (gemessen an tieferer Überschrift, Rauten-Zeile im Fence, Setext-Zeile,
  fehlender Schluss-Zeile); bewacht ist die Gleichheit nicht: macht man **nur**
  `closureSectionEnd` fence-blind (Übersprung entfernt), bleibt `make test`
  grün, obwohl die Substanz-Zählung ab dann einen anderen Bereich misst als die
  Platzhalter-Suche — eine Rauten-Zeile in einem Beispielblock beendete den
  Abschnitt für die eine Sicht und nicht für die andere. Der Slice-Plan hatte
  BEO-003 für diese Bedingung ausdrücklich als einschlägig erklärt („diese
  Bedingung wird sie benutzen, nicht nachbauen") — eingelöst ist das für die
  Inline-Code-Paarung (`PreprocessMarkdown`), nicht für die Abschnitts-Grenze.
- `verifizierbar`: ja — Rückbau `if inFence` zu `if false` **nur** in
  `closureSectionEnd`, danach `make test` (grün); die beiden gröberen Rückbauten
  derselben Funktion (Ebenenvergleich verschärft, Grenze ganz entfernt) sind
  dagegen rot.
- `klasse`: geteilte Grenz-Logik nachgebaut statt geteilt, ohne Äquivalenz-Assertion

### F-6 — Drei weitere Rückbauten der neuen Erkennung überleben

- `kategorie`: LOW
- `quelle`: Maintainability (Testlast der Portierung, Slice §4 DoD „acht
  Rückbauten geprüft")
- `pfad`: `internal/hexagon/core/rules/planning.go:221` und
  `internal/hexagon/core/rules/planning.go:240-249`
- `befund`: Neben der Tag-Liste (F-4) bleiben drei weitere Mutanten grün:
  (a) der Schrägstrich in der Vorzeichen-Klasse — entfernt man ihn, meldet
  `docs/<datei>` neu, kein Test bemerkt es; (b) die Kappungslänge — von 40 auf
  400 gesetzt, bleibt alles grün, obwohl die Spezifikation „auf 40 Runen
  gekappt" zusagt; (c) die Kleinschreibung im Tag-Vergleich — entfernt man sie,
  meldet ein großgeschriebenes Tag neu. Damit überleben vier der elf von mir
  gefahrenen Rückbauten. Ob sie unter den acht des Plans waren, ist nicht
  entscheidbar: der Plan zählt sie, benennt sie aber nicht einzeln.
- `verifizierbar`: ja — je Mutation eine Temp-Kopie des Arbeitsstands,
  `make test` grün; die Verhaltensänderung je Mutation ist mit einer
  Fixture-Notiz gegen das Image zeigbar.
- `klasse`: neuer Code-Pfad ohne mutationsfeste Bewachung

### F-7 — Doc-Kommentar von `closureSectionProse` steht jetzt über `closureSectionEnd`

- `kategorie`: LOW
- `quelle`: Maintainability (`revive` §exported/package-comments-Konvention des
  Lint-Profils; godoc-Zuordnung)
- `pfad`: `internal/hexagon/core/rules/planning.go:192-199`
- `befund`: Die neue Funktion wurde **zwischen** den bestehenden Doc-Kommentar
  und seine Funktion gesetzt. Der Block beginnt mit „closureSectionProse liefert
  den Abschnitts-Text …" und dokumentiert nach der Einfügung
  `closureSectionEnd`; `closureSectionProse` (Zeile 284) steht seither ohne
  Kommentar. Wer die Zählung ändert, liest die Beschreibung der Grenz-Funktion
  und umgekehrt — dieselbe Naht, die F-5 unbewacht lässt.
- `verifizierbar`: ja — Sichtprüfung der Datei; `make lint` bleibt grün (das
  Profil führt weder `gofmt` noch eine Doc-Zuordnungs-Regel), das Gate deckt es
  also nicht ab.
- `klasse`: Kommentar-Anker zeigt auf die falsche Einheit

### F-8 — Zwei verschiedene Schritte tragen im Code dasselbe Label „C4b"

- `kategorie`: LOW
- `quelle`: Maintainability (Konvention: Kommentare zitieren die Schritt-Kennung
  der Spezifikation)
- `pfad`: `internal/hexagon/core/rules/planning.go:143` und
  `internal/hexagon/core/rules/planning.go:153`
- `befund`: Die Floskel-Prüfung trägt seit jeher den Kommentar „C4b: Floskel",
  obwohl die Spezifikation sie als Teil von **C4** führt; der neue Schritt
  heißt in der Spezifikation C4b und trägt im Code denselben Kommentar-Präfix.
  In derselben Funktion stehen damit zwei „C4b"-Blöcke, die auf verschiedene
  Vertragsstellen zeigen.
- `verifizierbar`: ja — Textsuche nach „C4b" in der Datei liefert zwei Treffer
  mit verschiedener Bedeutung; die Spezifikation kennt nur einen.
- `klasse`: Kommentar-Anker zeigt auf die falsche Einheit

### F-9 — Grund-Code-Enumerationen an vier Rändern, und die Release-Prep-Checkliste kennt den Auslöser nicht

- `kategorie`: LOW
- `quelle`: Beobachtungs-Register BEO-002 / Release-Prozess
  (`docs/user/releasing.md` §Release-Prep Punkt 4)
- `pfad`: `harness/README.md` (Sensor-Zeile zu `make verify-closure-notes`),
  `README.de.md:86`, `README.md:87`, `docs/user/benutzerhandbuch.md:900` und
  `docs/user/benutzerhandbuch.md:1591`
- `befund`: Alle vier Stellen zählen die Closure-Bedingungen als
  „missing/thin/boilerplate" auf. Für Handbuch und README ist der Nachzug
  Release-Prep — die Checkliste kennt aber nur die Auslöser „neues Modul" und
  „neue CLI-Option", nicht „neuer Grund-Code in einem bestehenden Modul"; für
  diesen Diff greift also kein Punkt der Liste. Die Sensor-Zeile in
  `harness/README.md` ist zusätzlich **jetzt** überholt und nicht erst zum
  Release: mit `placeholder: true` in `.d-check.closure.yml` erzwingt
  `make verify-closure-notes` seit diesem Commit eine vierte Bedingung, die
  seine Autoritäts-Zeile nicht nennt.
- `verifizierbar`: ja — die Sensor-Zeile gegen einen Lauf von
  `make verify-closure-notes` mit einer Notiz halten, die nur die vierte
  Bedingung verletzt: das Gate meldet rot, die Zeile erklärt es nicht.
- `klasse`: Rand einer Semantik-Änderung bleibt auf der widerrufenen Fassung

### F-10 — Weitere gemessene Falsch-Positiv-Klassen, die der Vertrag nicht nennt

- `kategorie`: LOW
- `quelle`: ADR-0052 §Kontext („Ein Falsch-Positiv ist hier teurer als ein
  übersehener Platzhalter") / `DC-FA-PLAN-001` §Drei Einschränkungen
- `pfad`: `internal/hexagon/core/rules/planning.go:255-282`
  (`checkClosurePlaceholder`)
- `befund`: Über die drei genannten Einschränkungen hinaus melden fünf weitere
  Formen, die kein Platzhalter sind: (a) ein **eingerückter** Code-Block (vier
  Leerzeichen) — d-check modelliert ihn nicht, sein Inhalt ist Prosa;
  (b) ein Markdown-Linkziel in Winkelklammern, etwa ein Ziel mit Leerzeichen im
  Pfad; (c) der Rumpf eines mehrzeiligen HTML-Kommentars (die Öffnungszeile
  fällt über das Ausrufezeichen raus, die Folgezeilen nicht); (d) ein
  Nicht-ASCII-Zeichen unmittelbar vor der Klammer — die Wortzeichen-Klasse ist
  ASCII, `Café<T>` meldet, `Cafe<T>` nicht; (e) ungerade Backtick-Parität im
  Absatz, die einen als Inline-Code gemeinten Ausdruck literal werden lässt.
  Keine der fünf trifft den eigenen Bestand (0 von 96 Notizen, nachgemessen:
  `make verify-closure-notes` meldet 0 Befunde über 329 Dateien), aber sie sind
  die Klassen, an denen ein Adopter die Erkennung erlebt; die Messung deckt nur
  die eigene Schreibkultur ab.
- `verifizierbar`: ja — fünf Fixture-Notizen gegen das gebaute Image mit
  `placeholder: true`; alle fünf melden, die Gegenproben (mehrzeiliger
  Code-Span, doppelte Backticks, `a<b>c`, `2<3>1`, umgekehrte Zeichen-Reihenfolge)
  bleiben grün.
- `klasse`: Zusage-Reichweite größer als die Messung

### F-11 — Falsch-Negative sind vertragskonform, der Re-Evaluierungs-Trigger nennt aber nur zwei der gemessenen Formen

- `kategorie`: INFO
- `quelle`: ADR-0052 §Re-Evaluierungs-Trigger
- `pfad`: `docs/plan/adr/0052-platzhalter-erkennung-inline-code.md`
  §Re-Evaluierungs-Trigger
- `befund`: Gemessen unerkannt bleiben: doppelte geschweifte Klammern,
  `TODO:`/`XXX`, ein über zwei Zeilen gebrochener Platzhalter und
  `_Ausstehend._`. Die ersten beiden nennt der Trigger; die Zeilen-Bindung sagt
  die Spezifikation ausdrücklich zu („bis zur nächsten schließenden Klammer
  derselben Zeile"), und die unterstrichene Vorlagen-Form fängt in der Praxis
  die Substanz-Schwelle, solange sie allein im Abschnitt steht. Erkannt werden
  dagegen — geprüft — Großbuchstaben-Felder, Felder mit Leerzeichen und
  Platzhalter in Tabellenzellen des Abschnitts. Die Zusage ist damit ehrlich;
  undokumentiert ist allein, dass die Zeilen-Bindung und die Unterstrich-Form
  bewusste Lücken sind.
- `verifizierbar`: ja — je eine Fixture-Notiz pro Form gegen das gebaute Image.
- `klasse`: bewusste Lücke ohne Notiz im Vertrag

### F-12 — Ältere Images lehnen den neuen Schlüssel fail-closed ab (Konsumenten-Hinweis fehlt)

- `kategorie`: INFO
- `quelle`: `DC-FA-CONF-001` (strikte Konfigurations-Dekodierung) /
  Release-Kommunikation
- `pfad`: `internal/adapter/driven/configyaml/configyaml.go:155-162`
  (`rawClosure`)
- `befund`: Eine Konfiguration mit `planning.closure.placeholder` bricht auf dem
  Vor-Commit-Image mit Exit 2 und `field placeholder not found` ab. Das ist die
  zugesagte Strenge (das Handbuch führt den unbekannten Schlüssel als
  Exit-2-Fall), trifft aber genau den Konsumenten, für den der Change Request
  gestellt wurde: setzt er den Schlüssel und pinnt weiter das alte Image, steht
  sein Gate rot, ohne dass eine Zeile im Änderungs-Protokoll darauf hinweist —
  bei der Vorgänger-Anforderung derselben Welle wurde ein solcher
  Konsumenten-Hinweis noch mitgeschrieben.
- `verifizierbar`: ja — dieselbe Fixture-Konfiguration gegen beide Images:
  Exit 1 mit Befund (neu) gegen Exit 2 mit Konfigurationsfehler (alt).
- `klasse`: Vertrags-Erweiterung ohne Konsumenten-Hinweis

## Negativbefunde

- geprüft, ohne Befund: **Byte-Identität ohne den Schalter.** Vor-Commit-Image
  gegen Diff-Image, drei Läufe je Ausgabeform — Repo mit der konventionellen
  Konfiguration, Fixture-Repo mit gesetztem `closure.dir` **ohne**
  `placeholder`, jeweils Klartext und `--json`: alle drei Paare byte-identisch,
  gleiche Exit-Codes. Auch `--doctor` ist ohne Befunde identisch (die Diagnose
  erklärt nur aufgetretene Grund-Codes, nicht den Katalog).
- geprüft, ohne Befund: **`--doctor`-Klartext des neuen Codes.** Erscheint mit
  Befund korrekt, benennt die Opt-in-Eigenschaft, keine Reichweiten-Behauptung
  über andere Module (die Klasse, die in Lastenheft 0.52.3 nachgezogen werden
  musste).
- geprüft, ohne Befund: **Abschnitts-Grenze im Ist-Zustand.** Tiefere
  Überschrift, Rauten-Zeile in einem Fence, Setext-Zeile, Datei ohne
  Schluss-Zeile, Platzhalter vor und hinter dem Abschnitt — die neue Grenze und
  die Prosa-Grenze der Zählung liefern in allen fünf Fällen dasselbe Ergebnis;
  Divergenz erst nach Mutation (F-5).
- geprüft, ohne Befund: **Inline-Code-Behandlung.** Mehrzeiliger Code-Span
  innerhalb eines Absatzes und doppelte Backticks werden korrekt geleert; die
  Paarungsregel ist die geteilte (`PreprocessMarkdown`), nicht nachgebaut.
- geprüft, ohne Befund: **Erster Treffer / Verwerfungs-Schleife.** Der
  Vorlagen-Rumpf meldet genau einen Befund an der Zeile des ersten Treffers;
  ein vom Nachfilter verworfener Treffer beendet die Suche nicht.
- geprüft, ohne Befund: **Fitness Function „eigener Bestand bleibt bei null".**
  `make verify-closure-notes` mit `placeholder: true`: 329 Dateien, 0 Befunde.
- geprüft, ohne Befund: **Gates.** `make test`, `make lint` und
  `make verify-closure-notes` grün gegen den Arbeitsstand; die Lockstep-Prüfung
  Grund-Code-Katalog gegen §4-Tabelle ist Teil von `make test` und hält.
- geprüft, ohne Befund: **SemVer.** Minor ist korrekt — neuer Grund-Code, kein
  bestehender Befund ändert sich, Default `false`, Byte-Identität gemessen.
- geprüft, ohne Befund: **Ränder mit Nachzug.** §2-Schema-Zeile,
  §4-Grund-Code-Zeile, Änderungs-Tabellen beider Spec-Straten, ADR-Index-Zeile,
  `print-config`-Vorlage und `.d-check.closure.yml` sind vorhanden, gegenseitig
  konsistent und nennen dieselbe Semantik (Default aus, Inline-Code
  ausgenommen, zwei Nachfilter).
- geprüft, ohne Befund: **Referenz-Richtung.** ADR-0052 nennt den Slice nur im
  Kopf-/Geschichte-Bereich; der Kandidaten-Weg ADR → Slice ist nicht als
  Entscheidungsgrundlage benutzt.
- geprüft, ohne Befund: **Modul-Grenze (BEO-004).** Die Bedingung liest keine
  zusätzliche Eingabe; sie sieht denselben Abschnitt enger.
- geprüft, ohne Befund: **Formatierung.** Die neue Zeile im Diagnose-Katalog
  ist nicht ausrichtungs-stabil (ein Leerzeichen zu viel), es gibt dafür aber
  weder einen Linter im Profil noch eine Konventions-Regel — kein Finding
  (Stil ohne Anker).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 5 |
| LOW | 5 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Rand einer Semantik-Änderung bleibt auf der
widerrufenen Fassung (3×: F-1, F-2, F-9) · Zusage-Reichweite größer als die
Messung (2×: F-3, F-10) · Nachfilter-Liste ohne Vertrag und ohne Test ·
geteilte Grenz-Logik nachgebaut statt geteilt, ohne Äquivalenz-Assertion ·
neuer Code-Pfad ohne mutationsfeste Bewachung · Kommentar-Anker zeigt auf die
falsche Einheit (2×: F-7, F-8) · bewusste Lücke ohne Notiz im Vertrag ·
Vertrags-Erweiterung ohne Konsumenten-Hinweis

Zwei dieser Klassen sind Wiederholungen aus dem Register: die
Rand-Klasse (dort BEO-002) tritt in **diesem** Lauf dreimal auf und erreicht
damit die Schwelle, ab der das Register „verkörpern statt weiterzählen"
vorsieht; die Grenz-Logik-Klasse ist die Register-Klasse BEO-003, deren eigene
Lehre („je Konsument eine Assertion gegen Wieder-Divergenz") hier nicht
angewandt wurde.

## Verdikt

**Merge-blockierend:** ja — fünf MEDIUM. Keine davon stellt die Mechanik in
Frage: die Portierung ist sauber, die Byte-Identität gemessen, die
Fitness-Function-Kernaussage (eigener Bestand bei null, Vorlagen-Rumpf meldet
genau einmal) hält. Blockierend sind die **Vertrags-Ränder**: zwei Stellen des
ausgelieferten Vertrags sagen etwas anderes als das Produkt (F-1, F-2), eine
dritte sagt mehr zu, als die Erkennung hält (F-3), und zwei Zusagen sind
unbewacht (F-4, F-5).

**Release-Empfehlung: noch nicht taggen.** Vor dem Release nötig:

1. F-1 bis F-3 im Vertrag heilen (Spezifikation C5 + C4b-Ziel-Feld, das
   Substanz-Bullet des Lastenhefts, die Vergleichszeichen-Zusage samt
   Testfall für die Form ohne Leerzeichen).
2. F-4 und F-5 bewachen — Tag-Liste und Abschnitts-Grenze so belegen, dass ihr
   Rückbau rot wird.
3. Release-Prep, die der Diff bewusst offenlässt: `version.md`,
   `CHANGELOG.md`-Schnitt, Handbuch (Header-Stempel, **§4.17-Text**,
   Modul-Tabelle, `closure`-Schlüsseltabelle in §5, neue Versionsverlauf-Zeile
   chronologisch **unter** die letzte), beide README-Fassungen (DE zuerst) und
   die Sensor-Zeile in `harness/README.md` — Letztere ist schon vor dem Release
   überholt (F-9). Der Digest-Pin folgt nach dem Tag.
4. ADR-0052 bei der Closure auf `Accepted` heben (heute korrekt `Proposed`).

**Übergabe:** Die Findings gehen an den Implementer; die Klassen-Zeile geht in
die Slice-Closure und von dort in den Zähler des Beobachtungs-Registers. Dieser
Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill, dieses Modell, dieses
Verdikt) und ersetzt keine Verifikation — die DoD-Abhakung prüft der Verifier
separat.
