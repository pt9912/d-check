# Review-Report: welle-95 — image-test Phase (9) `linelength` — 2026-11-23

**Review-Art:** Code (Diff gegen Welle-Plan, Anforderung, Spezifikation und Hard Rules,
Maintainability; `v6.18.1` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten). Nicht gegen
die DoD.

**Gegenstand:** Commit `b792b305`. Er bringt Phase (9) in `tools/image-test.sh`, dazu
`harness/sensors/image-test.md` und `docs/user/abdeckung-e2e.md`. Er ist der Beleg für den
zweiten Punkt des „Mehr“ im Closure-Trigger von `welle-95` (§3, „Die Entscheide tragen am gebauten
Image“). Vorbild ist die `pinsync`-Phase (8) aus `5b669af0`.

**Skill:** `.harness/skills/reviewer.md` @ 1.20.0 (`d1c37c14`).
**Modell:** claude-opus-5-5 · **Datum:** 2026-11-23.

**Eingangs-Kontext:**

- Welle-Plan `welle-95` §3 (Closure-Trigger, die verlangten Proben)
- `DC-FA-LLEN-001` (Zählregel, Ausnahmen, Befund, Bezugsmenge, Akzeptanzkriterien), `DC-FA-LLEN-001.a`
  Schritte 3 bis 8, [ADR-0114](../plan/adr/0114-linelength-modul.md)
- Vorbefunde am Modul: Code-Review [R1](2026-11-23-slice-290-code-r1.md) und
  [R2](2026-11-23-slice-290-code-r2.md) zu `slice-290`
- `AGENTS.md` §3.1, §3.7, §4 und §5 (Regeln 13, 15)
- Die drei Restpunkte, die der Implementer selbst nennt: kein Lauf mit Exit 0, kein Vergleich `:ro`
  gegen beschreibbar, Meldungstext nur lose geprüft
- Gelesen am Gegenstand: `tools/image-test.sh` (Kopf und Phasen (8) und (9)),
  `harness/sensors/image-test.md` vollständig

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse | Träger |
|---|---|---|---|---|---|---|---|
| F-1 | INFO | Die Meldung prüft die Phase nur darauf, dass beide Zahlen vorkommen, nicht ihre Rollen. Am Image gefahren: Vertauscht der Code Länge und Grenze in der Meldung („Zeile hat 40 Zeichen, erlaubt sind 41“), bleibt `make image-test` grün. Der Kommentar an `ll_lauf` („die Länge und die Grenze in der Meldung“) und die Botschaft („mit Laenge und Grenze in der Meldung“) stimmen wörtlich, lesen sich aber stärker. Die genaue Meldung hält `TestLinelengthGrenzeGenauAufDerSchwelle` gegen denselben Code. An der Distributionsgrenze ändert sich die Formatierung nicht, deshalb bleibt es INFO. | `AGENTS.md` §5 Regel 15 | `tools/image-test.sh` · „`$4 ~ ("(^\|[^0-9])" n "([^0-9]\|$)")`“ | ja — die Mutation oben, `make image-test` grün | meldung-nur-auf-zahlen-geprueft | Implementer (optional) |
| F-2 | INFO | Der Abschnitt *Bindung* von `harness/sensors/image-test.md` führt `DC-FA-LLEN-001` nicht, obwohl der Vertrag jetzt einen Punkt dazu trägt. Dasselbe gilt für `DC-FA-MET-001` und `DC-FA-PSYNC-001`, schon vor diesem Commit. Das ist keine Überzusage, der Vertrag sagt eher mehr zu, als *Bindung* nennt. Wer aber über *Bindung* sucht, welche Anforderungen das Gate trägt, findet die drei Module nicht. | Maintainability | `harness/sensors/image-test.md` · „Schritt in `release.yml`, nicht in `ci`.“ (Abschnitt *Bindung*, danach die Kennungs-Liste) | nein | bindung-spiegelt-vertrag-nicht | Implementer |

## Beurteilung der drei Restpunkte des Implementers

| Restpunkt | Urteil |
|---|---|
| Kein Lauf mit Exit 0 in der Phase | **tragbar.** §3 verlangt „genau die Befunde, die die Anforderung zusagt“, kein Exit 0. Dass die Zeilen mit 39 und 40 Zeichen und die Mehrbyte-Zeile mit 40 Codepunkten befundfrei sind, belegt der exakte Vergleich der Befundmenge (`cut -f1-3` gegen die erwartete Liste). Am Image gefahren: eine Mutation `<=`→`<` färbt die Phase rot, mit `Befunde probe.md:3 … want Zeilen 4:41 6:41`. Den Exit 0 eines befundfreien Laufs am Image zeigt Phase (2); die Abbildung Befund → Exit ist modulunabhängig. Der Sensor-Vertrag sagt für `linelength` keinen Exit 0 zu |
| Kein Vergleich `:ro` gegen beschreibbar | **tragbar.** Die Phase hält die Ausgabe unter `:ro` gegen eine **feste Erwartung**. Das ist strenger als die Gleichheit mit einem zweiten Lauf, der denselben Fehler hätte. Der Vertrag sagt den Vergleich für `linelength`, anders als für `pinsync`, nicht zu. Dass das Modul nie schreibt, ist Eigenschaft des Kerns (`DC-FA-LLEN-001.a`, hermetisch); dass es unter `:ro` und mit fremdem Benutzer liest, belegt die Phase |
| Meldungstext nur lose geprüft | **tragbar, mit Rest F-1.** Die Phase prüft, dass beide Zahlen in der Meldung stehen; die genaue Meldung hält der Kern-Test |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Deckt das Fixture jede Probe aus `welle-95` §3? | geprüft, ohne Befund. Zeichen gezählt: Zeile 2 mit 39, Zeile 3 mit 40, Zeile 4 mit 41 Zeichen; Zeile 5 mit 40 Codepunkten in 90 Bytes und Zeile 6 mit 41 in 91 (`äö😀` = 3 Codepunkte, 8 Bytes); Zeile 8 mit 62 im Backtick-Zaun; Zeile 10 Tabellenzeile mit 63; Zeile 11 mit 66 und einer URL mit 55 Zeichen. Damit sind knapp unter, genau auf und knapp über der Grenze, Mehrbyte, Fenced-Code, Tabelle und lange URL alle da |
| Stimmen die Erwartungen mit `DC-FA-LLEN-001.a`? | geprüft, ohne Befund. Nachgerechnet: Default 11 − 3 (Zaun samt Zäunen) − 1 (Tabelle) = 7 geprüft, davon 1 über (c), Befunde 4 und 6. Bei (a) aus 10 geprüft, dazu Zeile 8, die Zaunzeilen mit 3 Zeichen bleiben befundfrei. Bei (b) aus 8 geprüft, dazu Zeile 10 (längstes Wort 4). Bei (c) aus 7 geprüft, 0 ausgenommen, dazu Zeile 11 |
| Prüft die Phase aus dem richtigen Grund? Rot-Stichprobe am Image per `sed` am Produkt, je einzeln, `make image-test`, danach zurückgesetzt | 4 von 4 rot in Phase (9), je mit der passenden Meldung. Bytes statt Codepunkten ergibt `Befunde probe.md:4 probe.md:5 probe.md:6`. Die Länge aus der Meldung entfernt ergibt `Meldung zu Zeile 4 nennt nicht 41 und 40`. `U` nicht gezählt ergibt `Bezugsmengen-Zeile … want 7 geprüft, 1 ausgenommen`. Die Grenze `<` statt `<=` ergibt `Befunde probe.md:3 …`. Die fünfte Mutation (Zahlen vertauscht) überlebt (F-1). Danach `make image-test` unverändert grün, `git diff --exit-code` sauber |
| Exakte Befundmenge statt Teilmenge | geprüft, ohne Befund. `cut -f1-3` der gesamten Ausgabe gegen die Liste in Pfad- und Zeilen-Reihenfolge; ein zusätzlicher Befund, auch einer aus `links`, färbt rot. Die Bezugsmengen-Zeile wird zeilengenau geprüft (`grep -qx`) |
| Isolierung: `--network none`, `--user 12345:12345`, `:ro`, `--enable linelength`, Plattform-Argument wie in den übrigen Phasen | geprüft, ohne Befund |
| `AGENTS.md` §3.7 (Kommentare) | geprüft, ohne Befund. Der Kopf-Eintrag (9) und der Abschnittskommentar sind Zusagen. Die Zeilenkommentare (`# Z. 4: 41`) beschreiben, was das Fixture trägt, sie sind keine Mess-Labels. `# abdeckung: DC-FA-LLEN-001` ist der Marker, den `make abdeckung` liest. Keine Review-, Slice- oder Befund-Marker |
| `AGENTS.md` §3.1 | geprüft, ohne Befund. Nur Container und POSIX-Werkzeuge (`printf`, `cut`, `awk`, `grep`) |
| Sensor-Vertrag ehrlich? | geprüft. Jeder Teilsatz des neuen Punkts ist durch eine Prüfung der Phase gedeckt: 39/40 befundfrei, 41 Befund, Codepunkte statt Bytes, die drei Ausnahmen im Default, je Schalter genau seine Zeile, die Bezugsmengen-Zeile auf stderr. Exit 0, `--json`/`--yaml` und Eingabefehler sagt er nicht zu. Rest: F-2 |
| `docs/user/abdeckung-e2e.md` | geprüft, ohne Befund. Die Zeile entspricht dem Abschnittskommentar der Phase; `make test` hält die Datei gegen ihre Ableitung |
| Welle-Plan §3: Tabelle Probe → Ergebnis in `welle-95-results.md` | nicht Gegenstand dieses Commits; sie entsteht mit der Welle-Closure (Planner) |
| `make doc-check` vor dem Report | grün: `1275 Datei(en) geprüft, 0 Befund(e)` |

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 0 | — |
| LOW | 0 | — |
| INFO | 2 | F-1, F-2 |

## Verdikt

Kein blockierender Befund. Das Fixture trägt jede Probe, die `welle-95` §3 verlangt. Die Phase
prüft die Befundmenge exakt und die Bezugsmenge zeilengenau, und sie wird aus dem richtigen Grund
rot: 4 von 4 Produkt-Mutationen schlagen mit der passenden Meldung an. Der Commit trägt damit den
zweiten Punkt des „Mehr“; die Tabelle Probe → Ergebnis steht bei der Welle-Closure aus. Die drei
Restpunkte des Implementers sind tragbar. F-1 und F-2 sind Hinweise ohne Wirkung auf den
Closure-Trigger.
