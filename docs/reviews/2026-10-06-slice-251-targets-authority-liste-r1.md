# Review R1 — slice-251: `targets.authority` nimmt eine Liste an

- **Review-Art:** Code. Geprüft wurde gegen den Slice-Plan (`slice-251`, §1 Abgrenzung,
  §3 Spiegel-Liste, §6 Risiken), gegen den eingehenden CR von `ai-harness-init` vom
  2026-10-06 (sieben Akzeptanzkriterien, Abgrenzung, zwei offene Fragen), gegen
  `ADR-0100` (Proposed) und `ADR-0099` (Vorgängerin am Modul), gegen `MR-025`/`MR-032`
  und gegen die Hard Rules `AGENTS.md` §3.1/§3.2/§3.7/§3.8 sowie §5 Regel 13/15/17.
  Gegenstand ist die Maintainability; die DoD-Abhakung prüft dieses Review nicht.
- **Gegenstand:** `slice-251` · Commit `92dde5e8`. Er umfasst `decodeTargetsAuthority`
  und `applyTargets` (Config-Rand), `undocumentedFindings` (Kern), `TargetsConfig`,
  das `--print-config`-Gerüst, Tests, Lastenheft 0.96.0, Spezifikation (Schritt 5,
  Schema-Zeile, Grund-Code-Zeile `SPEC-061`), `ADR-0100` und den ADR-Index.
- **Skill:** `reviewer.md` @ 1.16.0
- **Modell-ID:** claude-opus-5-5
- **Datum:** 2026-10-06
- **Eingangs-Kontext:** Slice-Plan (`slice-251`); aus dem Lastenheft `DC-FA-TGT-001`
  (0.96.0), aus der Spezifikation `DC-FA-TGT-001.a` (Schritte 1 und 5), `SPEC-005`
  (Schema-Zeile) und `SPEC-061`; `ADR-0100`, `ADR-0099`; der CR-Wortlaut; die
  Klartext-Tabelle in `internal/hexagon/core/app/diagnose.go`. Vorherige Findings am
  Modul `targets`: R1 zu `slice-250` (Klassen `grenze-gegen-beschreibung-statt-gegenstand-geprueft`,
  `grenze-ohne-negativtest`, `dublette-ohne-pfad-normalisierung`, `release-prep-nachzug`).
- **Proben:** Zwei Images aus `git archive` gebaut — Vorher (`92dde5e8^`) und Nachher
  (`92dde5e8`, Image-ID identisch mit `d-check:latest`) —, je Probe
  `docker run --rm --network none -v <fx>:/repo:ro <img> --enable targets` mit
  getrennt erfasstem stdout/stderr/Exit-Code und Byte-Vergleich. Fixture:
  `Makefile` mit `own`/`tool`/`ghost`, `harness/README.md` (`own`),
  `harness/targets.md` (`tool`, `own`), `harness/[b].md` (`own`). Zusatz-Images
  wieder entfernt; der Arbeitsbaum ist bis auf diesen Report unverändert. `make test`/
  `make gates` wurden nicht gefahren (Sache des Verifiers).

Probe-Ergebnisse (`authority: <Wert>`, `makefiles: [Makefile]`):

```text
authority-Wert                                   Vorher            Nachher
harness/README.md                                rc 1, 2 Befunde   byte-identisch (auch --json, --doctor, --doctor --json)
(Schlüssel fehlt) / "" / '' / (leer)             rc 0              byte-identisch
"null" / !!str null / 123 / true / "  "          rc 2 (Datei)      byte-identisch
./harness/README.md, !!str harness/README.md     rc 1              byte-identisch
null / ~ / Null                                  rc 0 (entfällt)   rc 2 „Doku-Datei "null" nicht lesen"
'harness/[b].md' (existierende Datei)            rc 1, 2 Befunde   rc 2 „ist ein Muster"
*a  (Alias auf einen Pfad-Scalar)                rc 1, 2 Befunde   rc 2 „muss ein Pfad oder eine Liste … sein"
{a: b}                                           rc 2, „line 3: …" rc 2, ohne Zeilenangabe
[harness/README.md, harness/targets.md]          rc 2 (Schema)     rc 1, genau ghost, Meldung nennt beide
[README, targets, README] / [README, ./README]   rc 2 (Schema)     Dublette einmal gelesen; Singular-Wortlaut bei [README, ./README]
[README, nope.md] / [nope.md, README]            rc 2 (Schema)     rc 2 mit Namen nope.md
[README, "*.md"] / [/abs.md] / [[README]]        rc 2 (Schema)     rc 2 (Muster / Pfad-Regel / keine Liste)
[]                                               rc 2 (Schema)     rc 0 (Richtung 2 entfällt, zugesagt)
[~]                                              rc 2 (Schema)     rc 0 — Richtung 2 entfällt still
authority:\n  -                                  rc 2 (Schema)     rc 0 — Richtung 2 entfällt still
authority:\n  - harness/README.md\n  -          rc 2 (Schema)     rc 1 — leerer Eintrag still verworfen
[harness/README.md, null]                        rc 2 (Schema)     rc 1 — leerer Eintrag still verworfen
```

## Findings

### H1 — HIGH — Ein leerer Listeneintrag in Block-Form (`-`) oder als `~`/`null` wird still verworfen; `authority: [~]` legt Richtung 2 still

- **kategorie:** HIGH
- **quelle:** `DC-FA-TGT-001` (Lastenheft 0.96.0: „ein leerer Eintrag … ist ein Konfigurationsfehler"); Spezifikation `DC-FA-TGT-001.a` Schritt 5 („ein leerer Eintrag … ⇒ Exit 2 beim Laden"); `ADR-0100` Entscheidung 4
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go:1703` (`n.Decode(&out)`) und `:1706` (Leer-Prüfung)
- **befund:** `yaml.v3` überspringt beim Dekodieren in `[]string` jedes Null-Element, statt es als `""` abzulegen; die Leer-Prüfung sieht deshalb nur den Fall `""` aus dem Test. Gemessen: `authority:` mit einem leeren Block-Eintrag `-` (die natürlichste YAML-Form eines leeren Eintrags) oder `[~]` ergibt Exit 0 mit 0 Befunden — Richtung 2 entfällt, obwohl ein Eintrag konfiguriert ist; `[harness/README.md, null]` läuft still als einelementige Liste. Das ist ein Stilles-Grün-Pfad im Gate (Prüffrage 1), und die Grenze „leerer Eintrag ⇒ Exit 2" ist gegen die Beschreibung statt gegen den Gegenstand geprüft (`AGENTS.md` §5 Regel 13); der Negativtest deckt nur `""` (Prüffrage 13).
- **verifizierbar:** ja — Black-Box-Lauf mit `targets:\n  makefiles: [Makefile]\n  authority:\n    -\n` liefert Exit 0; ein Unit-Fall `"targets:\n  authority: [~]\n"` in `TestDecode_TargetsAuthority` liefe heute grün statt mit „leeren Eintrag" rot.
- **klasse:** `grenze-gegen-beschreibung-statt-gegenstand-geprueft`

### M1 — MEDIUM — Die String-Form ist nicht byte-identisch: `null`/`~`, ein Alias und ein wörtlicher Pfad mit `[`/`*`/`?` verhalten sich anders als vorher

- **kategorie:** MEDIUM
- **quelle:** CR-Akzeptanzkriterium 1 („verhält sich byte-identisch wie heute"); Lastenheft 0.96.0 Akzeptanzkriterium „Autorität (String unverändert)" und Satz „Mit einer Datei — als String … — byte-identisch"; `ADR-0100` Entscheidung 1; `DC-QA-02`
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go:1697-1701` (Scalar-Zweig übernimmt `n.Value` ohne Tag-Auswertung), `:1714-1718` (Muster-Prüfung auch auf den Scalar), `:1711-1712` (`default` fängt `AliasNode`)
- **befund:** Gemessen gegen das Vorher-Image: `authority: null`, `~` und `Null` liefen vorher als „Richtung 2 entfällt" (Exit 0), jetzt Exit 2 mit „Doku-Datei "null" nicht lesen" — der Scalar-Zweig liest den `!!null`-Tag als Dateinamen. `authority: *a` (Alias auf einen Pfad) lief vorher (Exit 1), jetzt Exit 2. `authority: 'harness/[b].md'` auf eine existierende Datei lief vorher (Exit 1, zwei Befunde), jetzt Exit 2 „ist ein Muster"; eine Datei mit Glob-Zeichen im Namen ist als Autorität überhaupt nicht mehr ausdrückbar. Alle Brüche sind laut, aber das Lastenheft sagt in zwei Sätzen Byte-Identität für die String-Form zu und nennt keinen dieser Fälle als Grenze; die Muster-Regel steht im selben Absatz als allgemeine Regel, sodass sich der Vertrag für einen wörtlichen Pfad mit `[` selbst widerspricht.
- **verifizierbar:** ja — Byte-Vergleich Vorher/Nachher-Image mit den drei genannten Werten (Probe-Tabelle oben).
- **klasse:** `grenze-gegen-beschreibung-statt-gegenstand-geprueft`

### L1 — LOW — Der Kopfkommentar von `CheckTargets` nennt weiter „die Autoritäts-Doku" im Singular

- **kategorie:** LOW
- **quelle:** `MR-025` (Spiegel); `AGENTS.md` §3.7 (ein Kommentar beschreibt, was da ist)
- **pfad:** `internal/hexagon/core/rules/targets.go:45`
- **befund:** Der Diff zieht die Kommentare an `rawTargets`, `TargetsConfig` und `undocumentedFindings` auf die Vereinigung nach, nicht aber die Zusage im Modul-Kopf („ohne Eintrag in der Autoritäts-Doku ⇒ gate-undocumented"). Wer das Modul vom Einstieg her liest, findet die Ein-Datei-Semantik.
- **verifizierbar:** nein — kein Gate liest Kommentar-Semantik.
- **klasse:** `spiegel-kommentar-nicht-nachgezogen`

### L2 — LOW — Eine Abbildung als `authority`-Wert verliert die Zeilenangabe in der Fehlermeldung

- **kategorie:** LOW
- **quelle:** Maintainability (Diagnose-Qualität des Config-Rands)
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go:1712`
- **befund:** Vorher meldete `authority: {a: b}` „line 3: cannot unmarshal !!map into string", jetzt „targets.authority muss ein Pfad oder eine Liste von Pfaden sein" ohne Zeile; dasselbe gilt für den Alias-Fall aus M1. Der Zweig für eine verschachtelte Liste behält die Zeile, weil er die yaml-Meldung durchreicht — innerhalb derselben Funktion zwei Formen.
- **verifizierbar:** ja — Probe `{a: b}` gegen beide Images.
- **klasse:** `fehlerort-verschluckt`

### I1 — INFO — Der `--doctor`-Klartext bleibt im Singular; das ist tragfähig, aber nirgends als Entscheidung festgehalten

- **kategorie:** INFO
- **quelle:** `MR-025`; Spezifikation `SPEC-061`
- **pfad:** `internal/hexagon/core/app/diagnose.go:173`
- **befund:** `SPEC-061` sagt jetzt „in einer der `targets.authority`-Dokus", der Klartext „in der Autoritäts-Doku". Gemessen: `--doctor` und `--doctor --json` sind für die String-Form byte-identisch, und die Befund-Meldung (`Hinweis:`) nennt bei mehreren Dateien alle — die Lesart „die Autoritäts-Doku als Gesamtheit" trägt. Eine Änderung würde `reasonText` für jeden Konsumenten verschieben. Die Abwägung steht weder in `ADR-0100` noch in der Commit-Botschaft; der Plan hat die Frage in §3 gestellt, die Antwort fehlt.
- **verifizierbar:** nein
- **klasse:** `bewusste-nichtaenderung-undokumentiert`

### I2 — INFO — Null-Elemente verschwinden in `makefiles`/`doc-tables` genauso (Bestand)

- **kategorie:** INFO
- **quelle:** Maintainability; Prüffrage 11 (gleiche Eingabe-Klasse)
- **pfad:** `internal/adapter/driven/configyaml/configyaml.go:598-599`
- **befund:** Mit dem Vorher-Image gemessen: `doc-tables:` mit einem leeren `-` und `makefiles: [Makefile, ~]` laufen still (Exit 0). Das ist Bestand außerhalb der Plan-Abgrenzung (§1: `doc-tables` bleibt unverändert) und hier nur benannt, damit die Behebung von H1 die Asymmetrie bewusst entscheidet.
- **verifizierbar:** ja — Black-Box-Lauf gegen das Vorher-Image.
- **klasse:** `null-element-still-verworfen`

### I3 — INFO — Release-Prep-Flächen tragen die Ein-Datei-Semantik

- **kategorie:** INFO
- **quelle:** `AGENTS.md` §5 Regel 17
- **pfad:** `docs/user/benutzerhandbuch.md:1514` (§5-Beispiel), `docs/user/benutzerhandbuch.md:2645` (§6-Modultabelle „jede Regel steht in der Autoritäts-Doku")
- **befund:** Bewusst nicht im Feature-Commit (Plan §1); in der Release-Prep nachzuziehen, samt Hinweis auf die Listenform.
- **verifizierbar:** nein
- **klasse:** `release-prep-nachzug`

## Negativbefunde

- **CR-Akzeptanzkriterien 1–7:** 1 (String byte-identisch) für jeden Pfad-Wert ohne Glob-Zeichen erfüllt — Text, `--json`, `--doctor`, `--doctor --json` byte-gleich —, Ausnahmen siehe M1. 2 (Vereinigung) erfüllt (Probe, `TestCheckTargetsAuthorityListe`). 3 (fehlender Listeneintrag Exit 2) erfüllt, unabhängig von der Position, mit Dateinamen. 4 (Doppelt dokumentiert kein Befund) erfüllt (`own` in beiden Dateien). 5 (Pfad-Regel je Eintrag) erfüllt — `applyTargets` prüft die dekodierten Einträge (`/abs.md`, `../x.md` ⇒ Exit 2). 6 (rotes Gegenbeispiel, Fundstelle im Fragment) erfüllt (`ghost` an seiner Regelzeile, auch im Fragment-Test). 7 (Glob) entschieden: nein, Muster ⇒ Konfigurationsfehler (Folge siehe M1). Offene Fragen 1 (`--print-config` mit Listen-Beispiel) und 2 (Doppelnennung kein Befund) sind beantwortet.
- **CR-Abgrenzung und Plan §1:** `doc-tables`, `exempt-targets`, Regel-Extraktion und Fence-/Tabellen-Erkennung unverändert; keine Disjunktheits-Prüfung, kein Glob für `authority`; Handbuch/README/CHANGELOG nicht angefasst. Der Diff berührt nur die Dateien aus Plan §3 (dazu `lexikon_kopplung_test.go`, reine Typ-Anpassung). Der gemeinsame Decoder (`decodeStrict`) ist unverändert — die Rückführungs-Bedingung aus §4 ist nicht eingetreten. Geprüft, ohne Befund.
- **`undocumentedFindings` — Dubletten, Reihenfolge, Fehlerfall:** Dubletten über `path.Clean` (`README`/`./README` ⇒ einmal gelesen, Singular-Wortlaut mit der Schreibweise der ersten Nennung — die Klasse `dublette-ohne-pfad-normalisierung` aus der Vorrunde wiederholt sich nicht). Die Meldung folgt der Konfigurations-Reihenfolge, die Befunde der Regel-Reihenfolge; `DC-QA-02` gewahrt. Eine fehlende Datei bricht vor jedem Befund ab; `len(authority)==0` schließt den Zugriff auf `docs[0]` aus, und `docs` ist nach der Schleife nie leer. Geprüft, ohne Befund.
- **Strikter Decoder:** `yaml.Node` als Feldtyp hebt `KnownFields` für den Wert nicht auf eine Lücke — Abbildungen werden im `default`-Zweig abgelehnt, in Listen durch `n.Decode`. Unbekannte Schlüssel unter `targets` lehnt `decodeStrict` weiter ab. Geprüft, ohne Befund außer L2.
- **`authority: ""`, leerer Wert und fehlender Schlüssel:** byte-identisch zum Vorher-Image (Exit 0, Richtung 2 entfällt). Geprüft, ohne Befund — anders `null`/`~` (M1).
- **Hexagon-Richtung (`ADR-0005`):** Der Kern nutzt nur `path`/`strings`; `yaml` bleibt im Adapter. Geprüft, ohne Befund.
- **Netz außerhalb `external`, Inline-Suppression, Schwellen-Senkung (§3.2/§3.6):** keine. Geprüft, ohne Befund.
- **Kommentare (§3.7):** Die neuen Kommentare an `rawTargets`, `decodeTargetsAuthority`, `TargetsConfig`, `undocumentedFindings`, im Gerüst und an den Tests tragen Zusage, Abgrenzung oder Grenze; keine Review-Historie, Slice-Nummern oder Mess-Labels. Ohne Befund außer L1.
- **§3.8 (Zusage nur über die Scan-Menge):** Die Autoritäts-Dateien sind gelesene, nicht gescannte Eingaben — das war vorher so und gilt jetzt je Datei gleich (dieselbe `extractDocTargets`-Strecke mit Fence-Erkennung). Keine neue Achse. Geprüft, ohne Befund.
- **`MR-032`:** Lastenheft (Draft) auf 0.96.0 gehoben, Historie-Zeile mit CR-Verweis vorhanden; Spezifikations-Historie nachgezogen. Geprüft, ohne Befund.
- **`MR-025` (Spiegel):** Lastenheft (Beschreibung „in keiner Autoritäts-Doku", neuer Absatz, drei Akzeptanzkriterien, Out-of-Scope), Spezifikation (Schritt 5, Schema-Zeile, `SPEC-061`), Gerüst, Meldungstext und die Kommentare an Typ/Kern sind nachgezogen; `suggest.go` erzeugt keinen `authority`-Wert. Offen: L1, I1, I3.
- **`ADR-0100` Form:** `Bezug:` nennt `DC-FA-TGT-001`, `Schärft:` die Kennungen `DC-FA-TGT-001.a`, `SPEC-005`, `SPEC-061`. Der Marker `d-check:status-provenance` an `slice-251` zeigt die Entstehung und begründet keine Entscheidung. Re-Evaluierungs-Trigger vorhanden, `Proposed` bewusst; Index-Zeile ergänzt. Geprüft, ohne Befund.
- **Commit-Botschaft (§5 Regel 15):** Die Rot-Belege der Tests kann dieses Review nicht nachmessen (Sache des Verifiers). Überdehnt ist „leerer Eintrag … ein Konfigurationsfehler" (H1) und „Mit einer Datei bleibt der Meldungstext wie bisher" insofern, als einige String-Werte gar keinen Meldungstext mehr erreichen (M1); der Meldungstext selbst ist für jeden erreichten Fall byte-gleich.

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 1 | H1 |
| MEDIUM | 1 | M1 |
| LOW | 2 | L1, L2 |
| INFO | 3 | I1, I2, I3 |

Wiederkehrende Finding-Klasse für die Closure: `grenze-gegen-beschreibung-statt-gegenstand-geprueft`
(H1, M1) — zum dritten Mal in Folge am Modul-Rand (slice-249, slice-250, slice-251;
`BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen`, `AGENTS.md` §5 Regel 13). Neu
hier ist der Träger: die Grenze sitzt nicht im eigenen Code, sondern im Dekodier-
Verhalten der YAML-Bibliothek (Null-Elemente, Tags, Aliase) — die Black-Box-Probe
gegen das Vorher-Image fand alle drei, keiner der Unit-Fälle.

## Verdikt

**Nicht merge-bereit in dieser Form.** H1 ist ein Stilles-Grün-Pfad im Gate: Ein
konfigurierter, aber leerer Listeneintrag legt Richtung 2 ohne Meldung still, entgegen
der ausdrücklichen Zusage in Lastenheft, Spezifikation und `ADR-0100`. M1 widerspricht
der zweifach zugesagten Byte-Identität der String-Form; vor der Closure muss entweder
der Code oder der Vertrag die Fälle `null`/`~`, Alias und wörtlicher Pfad mit
Glob-Zeichen so tragen, wie der jeweils andere sie sagt. L1/L2 sollten behoben werden,
blockieren nicht; I1–I3 dienen der Information.
