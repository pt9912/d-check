# slice-265: `reviews` — Zusage-Muster, benannte Kennungen, Leerlauf und Unterverzeichnisse

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD selbst.

**Bezug:** eingehender CR von `sf-connector`
([`docs/plan/cr/2026-10-09-cr-eingehend-sf-connector-reviews-zusage.md`](../../cr/2026-10-09-cr-eingehend-sf-connector-reviews-zusage.md),
Punkte 1 bis 3); Befund aus der Planung von slice-263 (`reviews.done-dir` liest
keine Unterverzeichnisse); Auftraggeber-Entscheide 2026-10-09 (ein Release mit
slice-263; der Default erkennt zusätzlich die Vorlagen-Form, ausgehender CR an
den Kurs).

**Berührte Spec-Stellen:** [`DC-FA-RVW-001`](../../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in),
`spec/spezifikation.md` §2 (Schlüssel `reviews.*`) und die Grund-Code-Zeile von
`review-missing`.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Das Modul `reviews` prüft einen Bestand, dessen DoD-Zeilen und
Kennungen anders lauten als der eigene Default, statt grün über einer leeren
Menge zu laufen — mit opt-in-Schlüsseln und einem Default, der zusätzlich die
Form der Baseline-Slice-Vorlage erkennt. Ohne die Schlüssel ändert sich der
Befundsatz nur an den zwei benannten Default-Änderungen (§8 und die
Plan-Änderungen nach R1), sonst ist er byte-identisch:

- **Zusage-Muster** — ein RE2, das statt der festen Phrase „unabhängiger
  Review" den Text eines DoD-Punkts als Review-Zusage erkennt (CR-Punkt 1).
- **Kennungs-Muster** — die Kennung eines Slice auch in benannter Form
  (`slice-<welle>-<titel>`); eine Datei, deren Kennung nicht lesbar ist, wird
  nicht mehr still übersprungen (CR-Punkt 2).
- **Leerlauf-Schalter** — Kandidaten ohne eine einzige Zusage werden rot
  (CR-Punkt 3).
- **Unterverzeichnisse und Stubs** — `done-dir` mit Abstieg und Inhalts-
  Ausnahme, in derselben Form wie die Closure-Prüfung aus slice-263.
- **Die mitgelieferte Beschreibung** — Vorlage von `--print-config` und
  Spezifikation sagen, was erkannt wird; die falsche Aussage „eine Zeile, die
  *Review* nennt" fällt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Doku der Module im Image** (CR-Punkt 4) — slice-266; ein eigener
  Liefer-Gegenstand mit eigener Oberfläche.
- ~~**Der Default der Phrase**~~ — entfallen durch die Plan-Änderung nach R1 (zweiter Teil); ursprünglich: bleibt „unabhängiger Review"
  (Auftraggeber-Entscheid): ein breiterer Default änderte den Befundsatz
  bestehender Nutzer.
- **Die Konfiguration dieses Repos** für `reviews` — sie setzt die neuen
  Schlüssel erst, wenn sie released sind; geschnitten bei der Closure.
- **Das Benutzerhandbuch** — Release-Prep (`AGENTS.md` §5 Regel 17).

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] Zusage- und Kennungs-Muster samt Leerlauf-Schalter: Lastenheft,
      Spezifikation, Konfig-Validierung (Exit 2), Tests, die ohne die
      Änderung aus dem richtigen Grund rot sind; eine nicht lesbare Kennung
      ist ein Befund statt eines stillen Übersprungs.
- [ ] Abstieg und Inhalts-Ausnahme für `done-dir`, dieselbe Semantik wie
      slice-263.
- [ ] Ohne die neuen Schlüssel ist die Ausgabe unverändert bis auf die zwei
      benannten Default-Änderungen, gemessen altes gegen neues Image mit
      eingeschaltetem `reviews`, und jede der beiden ist durch einen Test
      belegt; Vorlage von `--print-config` beschreibt die
      Erkennung richtig; die Abnahme-Fälle des CR (rot: Zusage ohne Report,
      grün: 0 Befunde) als Test; `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/`; Verifikation.
- [ ] Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben;
      jedes Risiko aus §6 mit Ausgang; drei Paarungen hier geprüft; der CR
      trägt seine Entscheidung je Punkt.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Kern-Regel `reviews` | update | Muster, Kennung, Leerlauf, Abstieg |
| Konfig-Modell, YAML-Adapter, `--print-config`-Vorlage | update | Schlüssel und Beschreibung |
| `spec/lastenheft.md` (Anforderung des Moduls), `spec/spezifikation.md` | update | Anforderung und Verfeinerung |
| CR-Datei | update | Entscheidung je Punkt |
| ausgehender CR an den Kurs unter `docs/plan/cr/` | neu | Vorlage der `.d-check.yml` setzt `require-promises` |
| `harness/sensors/review-coverage.md` | update | Vertrag und Grenzen des eigenen Gates gegen den gemessenen Stand |
| `Makefile` (Hilfe-Zeile `review-coverage`) | update | nennt die Erkennung |
| `.d-check.yml` (Kommentar über dem `reviews`-Block) | update | nennt die Erkennung |

*(Plan-Änderung vor dem Code-Commit: die Sensor-Datei des eigenen Gates kommt
dazu — die Messung in diesem Slice zeigt, dass es heute keinen Slice prüft;
ihre Grenze sagte „Gewollt". Die Einordnung im CR wird mit derselben Messung
berichtigt.)*

*(Plan-Änderung nach R1, vor dem Code: Ein unlesbares Unterverzeichnis bricht
den Walk nicht mehr ab — es wird gemeldet, die übrigen Einträge werden weiter
gelesen, und der Leerlauf-Befund entfällt neben ihm (F-1). Kommentare und
Testtexte ohne Herkunft und Chronik, über die ganze Klasse gesucht (F-2, F-8,
F-11, F-7). Ein Negativtest hält die `SKIP_DIRS`-Ausnahme (F-3). Die
Byte-Identität wird mit eingeschaltetem `reviews` gemessen, altes gegen neues
Image (F-4). Die Grenzen nennen, dass das Muster jeden Checkbox-Punkt der Datei
sieht, auch in einem Codeblock, und den Punkt ab dem Bullet (F-5); YAML-`null`
gilt als abwesend (F-6). `match: name` verlangt nach dem Basisnamen ein
Zeichen, das kein Buchstabe und keine Ziffer ist — `slice-26` wird nicht mehr
vom Report zu `slice-265-…` gedeckt; die Präfix-Grenze mit Bindestrich bleibt
benannt (F-10). Die Sensor-Datei nennt statt eingefrorener Zahlen, wie die
Messung zu wiederholen ist (F-9).)*

*(Plan-Änderung nach R1, zweiter Teil, vor dem Code — Auftraggeber-Entscheid
2026-10-09, der den ersten zum Default revidiert: Die Baseline-Vorlage des
Slice-Plans formuliert den DoD-Punkt „Review durchgeführt, Report unter
`docs/reviews/` liegt vor", und ihre mitgelieferte `.d-check.yml` schaltet
`reviews` ohne Muster ein — jedes Repo, das der Baseline folgt, läuft damit
grün über einer leeren Menge. Der Default erkennt deshalb **zusätzlich** die
Vorlagen-Form „Review durchgeführt"; das ändert den Befundsatz dort, wo heute
nichts geprüft wird, und steht als Default-Änderung im Lastenheft. Dazu ein
**ausgehender CR** an den Kurs (`docs/plan/cr/`, [`MR-035`](../../../../harness/conventions.md#mr-035)): die
`.d-check.yml`-Vorlage setzt `require-promises`, sobald v0.85.0 released ist,
damit ein künftiges Auseinanderlaufen von Vorlage und Default laut wird. Die
Abgrenzung „Der Default der Phrase bleibt" in §1 entfällt damit.)*

*(Plan-Änderung nach R2, vor dem Code: Die Vorlagen-Form im Default war zu
breit — „Adaptions-Review durchgeführt" und „kein Review durchgeführt" trafen.
Sie gilt nur noch, wo sie den Punkt **eröffnet**, direkt hinter der Task-Box,
wie in der Vorlage (R2 F-1). Vier Spiegel nennen noch die alte Phrase:
Exit-2-Meldung, Kopfkommentar der Regel, Hilfe-Zeile von `review-coverage`
im `Makefile`, Vertrag der Sensor-Datei (F-3; das `Makefile` kommt in §3).
Die Wortgrenze nach dem Basisnamen gilt für jeden Buchstaben und jede Ziffer,
nicht nur ASCII (F-4); die Vorlage nennt `null` (F-5); Kommentare ohne
Slice-Nummern als Beispiel (F-6); die Lastenheft-Historie bekommt eine
berichtigende Zeile (F-7). Die Stellen im Plan, die „ohne die Schlüssel
unverändert" sagten, sind nachgezogen (F-2).)*

*(Plan-Änderung nach R3, vor dem Code: Der enge Default schloss eine Form aus,
die im eigenen Bestand eine echte Zusage ist — „`make gates` grün; Review
durchgeführt, Report unter …" (R3 F-1). Die Vorlagen-Form gilt deshalb hinter
der Task-Box, hinter `;` oder `,` und am Anfang einer Folgezeile; Verneinung
und Kompositum bleiben ausgeschlossen. Die `.d-check.yml` ist ein Spiegel
(F-2, Zeile in §3); die beiden Plan-Stellen aus F-3 sind nachgezogen.)*

*(Plan-Änderung nach R4, vor dem Code: Die Alternative „am Anfang einer
Folgezeile" machte eine umbrochene Verneinung zur Zusage (R4 F-1) und hat im
Bestand keinen Beleg. Sie entfällt; hinter `;` oder `,` gilt die Form auch über
einen Zeilenumbruch. Ein Test hält die umbrochene Verneinung.)*

## 4. Trigger

**Start** (`next` → `in-progress`): slice-263 in `done/`; `in-progress/` leer.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): Abstieg und Inhalts-Ausnahme für
  `done-dir` reißen mehr auf als erwartet — dann in einen eigenen Slice.

## 5. Closure-Trigger

DoD vollständig (§2) + Closure-Notiz mit Lerneintrag; die drei Paarungen
wellenlos hier geprüft. Produkt-Verhalten — geht mit demselben Release wie
slice-263 hinaus.

## 6. Risiken und offene Punkte

- **Kennungs-Muster und Report-Zuordnung** — eine benannte Kennung
  `slice-<welle>-<titel>` kann Präfix einer anderen sein; die Zuordnung
  Report → Slice muss eindeutig bleiben. — **Ausgang:** *(offen)*

## 7. Closure-Notiz

*(gefüllt vor dem `git mv` nach `done/`)*

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** geändert werden Produkt-Kern
(`reviews`), Konfig-Modell, YAML-Adapter, Konfig-Vorlage, Lastenheft und
Spezifikation — alle unter dem Default `*` (`ALL`); deklariert.

**Spiegel vor dem Editieren** (Schritt 17,
[`MR-025`](../../../../harness/conventions.md#mr-025); gemessen mit
`grep -rln "reviews\.\(done-dir\|reviews-dir\)\|ReviewsConfig\|unabhängiger Review\|review-missing"`
über Code und Doku, ohne eingefrorene Verzeichnisse): Kern-Regel und ihr Test,
Konfig-Modell, YAML-Adapter, Konfig-Vorlage, CLI-Abnahmetest, Lastenheft,
Spezifikation, `harness/sensors/review-coverage.md`; README (beide Sprachen)
und Benutzerhandbuch ziehen die Release-Prep nach (`AGENTS.md` §5 Regel 17).
`verify-closure-notes.md` nennt das Modul, ohne seine
Erkennung zu beschreiben — kein Spiegel. Der Kommentar über dem `reviews`-Block
der `.d-check.yml` nennt die Erkennung — ein Spiegel (Nachtrag nach R3).

**Entwurf vor dem Code** (die Liefer-Punkte aus §2 in Schlüssel übersetzt):
`reviews.promise-pattern` (RE2 gegen den Text eines DoD-Punkts; abwesend ⇒ die
Phrase „unabhängiger Review"; explizit leer oder nicht kompilierend ⇒ Exit 2),
`reviews.match` (`id` — Default, die `slice-<NNN>`-Kennung wie bisher — oder
`name`: ein Report deckt einen Slice, dessen Dateiname den Basisnamen des
Slice ohne `.md` enthält; anderer Wert ⇒ Exit 2), `reviews.require-promises`
(Kandidaten ohne eine einzige Zusage ⇒ Befund), `reviews.recursive` und
`reviews.skip-pattern` (dieselbe Semantik wie in der Closure-Prüfung). Unter
`match: id` meldet ein Slice mit Zusage, aus dessen Namen keine Kennung zu
lesen ist, einen Befund statt still zu fallen — eine Änderung am
Default-Verhalten; sie betrifft nur Dateien, die heute ungeprüft bleiben. Die
zweite — der Default erkennt die Vorlagen-Form — trägt die Plan-Änderung nach
R1, zweiter Teil.

<!-- d-check:cite .harness/baseline/v6.17.0/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** gelesen am 2026-10-09.
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/state.md)
(offen) — das Zusage-Muster wird gegen Negativfälle getestet, nicht nur gegen
die Formen aus dem CR; [`BEO-ALL/review-fix-applied-only-at-cited-site`](../observations/BEO-ALL/review-fix-applied-only-at-cited-site/state.md)
(2×) — Review-Befunde werden über ihre Klasse gesucht, ein dritter Treffer
wäre eine Lücke; [`BEO-ALL/kommentar-traegt-herkunfts-prosa-statt-fuenf-klassen`](../observations/BEO-ALL/kommentar-traegt-herkunfts-prosa-statt-fuenf-klassen/state.md)
— nach Verkörperung schon wieder aufgetreten; vor jedem Code-Commit greppt der
Lauf die neuen Kommentare und Testtexte nach Chronik-Wörtern;
[`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/state.md)
— die Grenzen der neuen Schlüssel werden am Code gezählt.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
gelesen am 2026-10-09 aus dem jüngsten Lauf (`make nightly-state`) —
`upstream-drift` grün (2026-10-08 11:37 UTC), `image-scan` grün
(2026-10-08 10:38 UTC).

**Modus-Begründungsblock:** GF — alle berührten Sub-Areas GF.
