# Review slice-165 — Die zugesagte Gleichheit war die falsche Größe

**Gegenstand:** [slice-165](../plan/planning/done/slice-165-dockerhub-spiegel.md), Stand des Feat-Commits vor der Nacharbeit.
**Datum:** 2026-08-27. **Reviewer:** unabhängiger Subagent.

---

## Urteil

**Nicht schließbar.** Die tragende Zusage der Anforderung ist **empirisch
widerlegt**, und die daraus gebaute fail-closed-Prüfung hätte jedes Release rot
gemacht. Zwei HIGH, sieben MEDIUM, vier LOW. Eigene Gate-Läufe: `make gates`,
`make ci`, `make completeness-check` — alle Exit 0; alle drei Zahlen der
Commit-Botschaft stimmen.

## Befunde

| # | Rang | Befund |
|---|---|---|
| F-1 | **HIGH** | *„Derselbe Manifest-Digest"* ist für `docker tag` + `push` **falsch**. Am Schwester-Repo gemessen: gleiche Image-ID, gleicher Config-Digest, aber **drei von drei** Tag-Paaren mit verschiedenen Manifest-Digests (neu komprimierte Layer-Blobs). Beide Pushes gelingen, der Vergleich schlägt fehl, `Create GitHub Release` läuft nicht mehr — der Spiegel liegt öffentlich, das Release existiert nicht |
| F-2 | **HIGH** | Das **Negative-Akzeptanzkriterium wird in seinem eigenen Fall nicht erfüllt**: die Zugangsdaten werden in einem `uses:`-Schritt verbraucht, der mit der generischen Action-Meldung scheitert; der `trap`, der den GHCR-Digest nennt, steht im **nächsten** Schritt und läuft nie an |
| F-3 | MEDIUM | Das Lastenheft sagt, ein Fehlschlag der Beschreibungsseite lasse das Release grün — `Prepare Docker Hub metadata` hat **kein** `continue-on-error` und steht **vor** `Create GitHub Release`. Source-Precedence-Widerspruch |
| F-4 | MEDIUM | `releasing.md` behauptet, die Hub-Seite pflege *„niemand automatisch"*, und belegt das mit einer ADR-Sektion, die etwas anderes sagt — dreifach falsch im selben Commit |
| F-5 | MEDIUM | `d-migrate` wird über seinen Geltungsbereich hinaus zitiert: es überspringt nur bei **fehlendem Secret**; bei gesetztem Secret und gestörtem Docker Hub wird es ebenfalls rot |
| F-6 | MEDIUM | Drei Nutzer-Oberflächen behaupten im Indikativ einen Spiegel, den es nicht gibt — `docker pull pt9912/d-check:v0.64.0` scheitert, weil v0.64.0 **vor** diesem Commit veröffentlicht wurde |
| F-7 | MEDIUM | Der `-z`-Zweig ist **unerreichbar**: findet `grep` nichts, kippt `pipefail` die Kommandosubstitution und `set -e` beendet den Lauf vorher. Dazu bricht der `^`-Anker, sobald `DOCKERHUB_IMAGE` die Schreibweise `docker.io/…` trägt, die Lastenheft und Handbuch selbst verwenden |
| F-8 | MEDIUM | Dritter Action-Pin ohne Frische-Achse; `AGENTS.md` und `harness/README.md` zählen weiter *„die zwei Action-Pins"* |
| F-9 | MEDIUM | DoD-Haken unerfüllt: `operations.md` zeigt den zweiten Bezugsweg nicht |
| F-10 | MEDIUM | Der *„fail-fast statt still"*-Schritt hat selbst zwei stille Pfade: leere Datei passiert den Wächter, fehlender `__VERSION__`-Platzhalter ebenso |
| F-11 | LOW | Ein Bestands-Kommentar wurde mitten entzweigeschnitten; beide Hälften tragen einzeln keine der fünf Klassen mehr |
| F-12 | LOW | Das benannte `versions`-Risiko ist ungemindert — die Auffanglinie in `releasing.md` §4 wurde nicht ergänzt |
| F-13 | LOW | Der Wächter misst **Bytes**, die Meldung behauptet Zeichen; die eigene Datei zeigt 96 vs. 95 |
| F-14 | LOW | Fehlende Leerzeile; die Pipeline-Liste endet bei Schritt 6 und kennt den neuen Ausfallweg nicht |

**Negativbefunde.** Der `trap` feuert bei explizitem `exit` **nicht** (keine
Doppelmeldung) und bei fehlschlagenden Top-Level-Kommandos korrekt; das
`${VAR#*@}`-Stripping trägt beidseitig; `RepoDigests` ist nach dem Push gefüllt;
das Boundary-AK (Prerelease lässt `:latest` unberührt) ist gedeckt; der
Digest-Vergleich ist **keine** Tautologie; §3.4 Abwärts-Sperre unverletzt; alle
`uses:` SHA-gepinnt; ADR-Form, Roadmap-Kopplung und die drei §7-Vorprüfungen in
Ordnung; die Zahlen in `overview.md` (20 Module, zwei Netz-Module) stimmen.

## Erledigung

Alle vierzehn Befunde sind eingearbeitet; die zwei HIGH sind **eigens
nachgemessen** statt übernommen.

- **F-1** Selbst nachgemessen und bestätigt. [ADR-0065](../plan/adr/0065-spiegel-gleichheit-ist-der-config-digest.md)
  löst [ADR-0064](../plan/adr/0064-dockerhub-spiegel-fail-closed.md) ab: die
  Gleichheit ist der **Config**-Digest, aus **beiden Registries** gelesen.
  Lastenheft 0.71.0 zieht nach; der registry-lokale Manifest-Digest steht jetzt
  ausdrücklich in Anforderung, Handbuch, beiden READMEs und `operations.md`.
- **F-2** Der Login ist ein `run`-Schritt mit Vorab-Prüfung — nur so trägt die
  Meldung den GHCR-Stand, den das Kriterium verlangt.
- **F-3**, **F-10**, **F-13** Das Zeichen-Limit prüft `make gates` als Go-Test
  (Zeichen, nicht Bytes); beide Darstellungs-Schritte tragen
  `continue-on-error`; leere Datei und fehlender Platzhalter sind verriegelt.
- **F-4**, **F-5**, **F-6**, **F-12**, **F-14** direkt behoben.
- **F-7** `sed -n …p` schweigt statt zu scheitern — der `-z`-Zweig ist
  erreichbar; die Schreibweisen-Falle entfällt, weil jetzt aus der Registry
  gelesen wird statt `RepoDigests` zu grep­pen.
- **F-8** `make hubdesc-pin-freshness` ergänzt, im Nachtlauf verdrahtet
  (`ok` gegen 5.0.0); beide Tabellen sagen jetzt **drei**.
- **F-9**, **F-11** behoben.

**Was der Review zusätzlich ausgelöst hat:** die Korrektur des ADR-Status legte
offen, dass `matrix.status.forbidden` auch die vom Kanon ausdrücklich erlaubte
**ADR→ADR-Lineage** verbot — wer den Status korrekt setzte, brach die
Supersede-Kette. `allow-supersede-lineage` ist jetzt eingeschaltet.
