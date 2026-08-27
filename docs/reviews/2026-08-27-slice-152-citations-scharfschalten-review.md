# Review slice-152 — `citations` scharfschalten

**Gegenstand:** [slice-152](../plan/planning/done/slice-152-citations-scharfschalten.md), Stand des Feat-Commits vor der Nacharbeit.
**Datum:** 2026-08-27. **Reviewer:** unabhängiger Subagent.

---

## Urteil

**Schließbar nach Nacharbeit.** Der Kern trägt — die Direktiven sind sachlich
richtig, das Modul läuft grün, die Wegwahl ist begründet. Zwei Befunde sind vor
der Closure zu beheben: eine **gebrochene, im `Makefile` selbst dokumentierte
Kopplung** und eine **Zahl, die über eine andere Menge spricht als die
gemessene**.

Selbst gefahren: `make doc-check` (550 Dateien, 0 Befunde), `make test`,
`make gate-consistency`, `make planning-check`, `make trace-check`,
`make adr-check`, `make verify-closure-notes` (494 Dateien),
`make completeness-check` (48/0), `citations` solo — je Exit 0.

## Befunde

| # | Rang | Befund |
|---|---|---|
| F-1 | **HIGH** | `FOCUS_DISABLE` nicht nachgezogen — vier fokussierte Gates feuern `citations` mit, fail-closed, und `adr-check` hängt im `pre-commit`-Hook. Der Makefile-Kommentar darüber sagt die Kopplung wörtlich |
| F-2 | **HIGH** | *„8 wörtliche Teilstrings → ausgezeichnet"* spricht über zwei verschiedene Mengen. Am Eltern-Commit: **8** verbatim, davon **6** ausgezeichnet. Am Commit-Stand: **10** verbatim, **8** ausgezeichnet. Dass 8 = 8 dasteht, sind zwei gegenläufige Fehler |
| F-3 | MEDIUM | *„am Repo nicht entscheidbar"* ist durch `MR-037`s Delta-Tabelle widerlegt — die Quelldatei des `MR-035`-Zitats steht nicht darunter, kann also nicht gedriftet sein |
| F-4 | MEDIUM | Die Tragbarkeits-Begründung ist eine Exklusivitäts-Aussage mit drei belegten Gegenwegen: der Bump (planmäßig, `MR-051` sagt es selbst zu), ein Merge, und `d-check:ignore` greift hier **nicht** |
| F-5 | MEDIUM | Alle acht Direktiven stehen als **HTML-Block** am Zeilenanfang und unterbrechen nach CommonMark den Absatz — mitten im Satz |
| F-6 | MEDIUM | `MR-051` kennt den dritten Fall nicht: eine **umbenannte** Zieldatei meldet `citation-out-of-range`, und „dieselbe Datei" ist dann gar nicht auswertbar |
| F-7 | MEDIUM | Der Netzlos-Modul-Guard im Go-Test hält `citations` nicht — er prüft eine **Teilmenge** |
| F-8 | LOW | Beide Bucket-Etiketten sagen mehr als die Menge trägt (Zeichensetzung; mindestens sechs der „25" sind sehr wohl Baseline-Zitate) |
| F-9 | LOW | Der Beleg des bewussten Brechens nennt `MR-049:41`; am Commit-Stand steht die Direktive auf einer anderen Zeile |
| F-10 | LOW | *„Der Punkt gehört außerhalb der Anführung"* beschreibt nicht, was der Diff zeigt — der Punkt wurde gelöscht |
| F-11 | LOW | `MR-049` ist jetzt das einzige Zitat im Speicher ohne Kursiv-Klammer; ohne Erklärung stellt die nächste Session es zurück und macht das Gate rot |
| F-12 | LOW | `.d-check.yml` trägt für `citations` keinen Kommentar, obwohl es das einzige **fail-closed** Modul im Profil ist (§3.7) |
| F-13 | LOW | Slice-Stand: DoD ungehakt, §5 ohne Ausgänge, §9 leer; §1 trägt zwei überholte Mess-Generationen. `slice-152` fällt **nicht** unter die Altbestands-Ausnahme des Closure-Profils — der Ausgang muss mit einem der drei Wörter beginnen |
| F-14 | LOW | Vorbestehend: die `v5.11.0`-Spalte in `MR-039` gibt `MR-033`s Zitat nicht wörtlich wieder |

**Was geprüft wurde und trägt:** alle acht Direktiven paaren mit dem
**beabsichtigten** Zitat (auch der heikle Fall mit zwei Zitaten in einem
Absatz), alle Spannen sind exakt und minimal, keine liegt in einer Tabelle oder
zwischen Zitat-Anfang und -Ende. Die zwei Korrekturen sind zulässig —
[`MR-039`](../../harness/conventions.md#mr-039) regelt den **Drift**-Fall, nicht
die Fehltranskription einer unveränderten Quelle —, und die Aussage der
Einträge ändert sich nicht. `MR-051` ist formal in Ordnung, und die MR-Form ist
die richtige: keine Gate-Lockerung, und alle früheren Modul-Scharfschaltungen
liefen ohne eigene ADR. §3 ist eingehalten.

## Erledigung

Alle vierzehn Befunde sind eingearbeitet:

- **F-1** `FOCUS_DISABLE` nachgezogen; die gebrochene Spiegel-Liste ist als
  zweiter Beleg von [`BEO-010`](../plan/planning/observations.md) eingetragen.
- **F-2**, **F-8**, **F-9**, **F-10** auf die gemessene Form gebracht.
- **F-3** durch eine **stärkere** Messung ersetzt: der alte Baum liegt in der
  git-Historie, und der Diff zeigt in beiden Quelldateien nur die geänderte
  Provenienz-URL. Alle fünf Abweichungen sind damit Transkription. Die Linie
  läuft jetzt über die **Richtung** statt über das Datum — was die Quelle hat
  und das Zitat weglässt, ist ein Fehler; was das Zitat hinzufügt, ein
  deklarierter Autoren-Akt.
- **F-4** in `AGENTS.md` §4, §Sensors und am `modules`-Eintrag korrigiert.
- **F-5** alle Direktiven ans Zeilenende gehängt; Paarung nachgeprüft.
- **F-6** dritter Fall in `MR-051` ergänzt.
- **F-7** `netlessDocModules()` erweitert.
- **F-11**, **F-12** als Kommentar bzw. Erklärung gesetzt.
- **F-13** im Closure-Body.
- **F-14** als benannter Rest, zusammen mit den Zitaten außerhalb des Speichers:
  [slice-163](../plan/planning/open/slice-163-zitate-ausserhalb-des-speichers.md).
