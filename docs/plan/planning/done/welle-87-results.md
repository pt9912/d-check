# Welle 87 — Die Wellen-Archivierung nachgerüstet, an 85 echten Wellen bewiesen — Closure-Notiz

**Welle:** welle-87-wellen-archivierung
**Abschluss:** 2026-09-03
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

Dieses Repo hatte, gemessen bei [slice-188](slice-188-register-gegen-neuen-kanon.md),
noch nie eine Welle archiviert — der Kanon verlangt das seit Einführung der
Regel für jede Wellen-Closure. `welle-87` liefert das fehlende Werkzeug und
wendet es sofort auf den **gesamten** Alt-Bestand an, nicht nur auf den
ursprünglich vorgesehenen Ausschnitt:

- [slice-190](slice-190-wellen-archiv-werkzeug.md): `tools/archive-wave`,
  ein eigenständiges Go-Programm (eigenes `go.mod`, eigenes `Dockerfile`,
  eigenes `Makefile`) — Sammeln nach `**Welle:**`-Feld, ZIP-Bau, Stub-Erzeugung
  nach Template, repo-weiter Verweis-Nachzug. An einem konstruierten Fixture
  bewiesen, nicht am echten Bestand.
- [slice-191](slice-191-alt-bestand-archivieren.md): das Werkzeug auf
  **alle 85** nummerierten Wellen dieses Repos angewendet (`welle-01` bis
  `welle-85`) — nicht nur `welle-60…85`, wie ursprünglich geplant. Die
  Scope-Erweiterung war ein Nutzer-Entscheid, ausgelöst durch die Messung,
  dass die Wellen-Nummerierung bei `welle-01` beginnt, nicht `welle-60`.

**Alle drei Closure-Trigger-Bedingungen sind erfüllt:** das Werkzeug führt die
Operation aus (Sammeln, ZIP, Stubs, Verweis-Nachzug in beiden Pfad-Formen);
alle vor dieser Welle geschlossenen Wellen sind archiviert (85 von 85, mehr
als die geforderten `welle-60…85`); die vermeintlich wellenlosen Alt-Slices
sind einer Zuordnung zugeführt — gemessen stellte sich heraus, dass **keiner**
von ihnen echt wellenlos war (vier trugen ihre Welle nur in Prosa statt im
Feld, nachgetragen; 52 weitere ab `slice-137` sind **bewusst** wellenlos nach
der Baseline-Konvention „wellenlos heißt nicht wächterlos" und liegen damit
außerhalb der Archivierungspflicht, kein Sammel-Archiv nötig).

`make gates` und `make fullbuild` sind auf dem vollständig archivierten
Bestand grün (50 Requirements, 0 Waisen).

## Was hat funktioniert?

**Die Trennung „Werkzeug bauen, nicht anwenden" (slice-190) hat den ersten
Slice klein gehalten**, während die eigentliche Belastungsprobe erst in
slice-191 kam — mit exakt der Vorsicht, die diese Trennung selbst benannte
([`BEO-011`](../observations.md), sechste Instanz): eine Fixture-Beweisführung
generalisiert nicht automatisch auf einen gewachsenen, 85 Wellen tiefen realen
Bestand.

**Die Wellen-für-Welle-Disziplin (Commit je Welle, `doc-check` danach) hat
jeden der insgesamt fünf während der Ausführung gefundenen Fehler sofort am
Ort seines Auftretens sichtbar gemacht** — nie erst Dutzende Wellen später,
nie musste eine der 85 Wellen zurückgenommen werden.

**Drei unabhängige Prüf-Runden fingen, was das andere Auge nicht sah:** ein
Review fand eine Staging-Lücke (zwei korrekt nachgezogene Verweise, nie
committet) und zwei Dokumentations-Lücken; eine erste Verifikation fand einen
Formfehler im `Hervorgegangen:`-Rückbau (56 von 59 Stubs); eine zweite,
gezielte Nachprüfung fand, dass die ERSTE Korrektur selbst eine Shell-Näherung
war, die punkt-suffigierte Kennungen falsch behandelte — aufgelöst durch einen
Wechsel des Mittels (das echte, getestete Produkt direkt aufrufen statt einer
weiteren Handnäherung), nicht durch eine weitere Regex-Iteration.

## Was ging anders als geplant?

**Der Scope wuchs zweimal, beide Male durch Messung statt Vorsatz:**
`welle-60…85` wurde zu `welle-01…85`, als sich zeigte, dass die
Wellen-Nummerierung länger zurückreicht als angenommen; und die
„Zuordnung wellenloser Alt-Slices" löste sich in eine Feld-Nachtrags-Aufgabe
für vier Slices auf, statt das im Kanon vorgesehene Sammel-Archiv zu bauen.

**Eine Mess-Korrektur (`Hervorgegangen:`-Backfill) brauchte zwei Anläufe, weil
der erste per Shell-Näherung statt per Produkt-Aufruf gerechnet wurde** — die
Lehre aus dem vorherigen Fehler (blindes Wegwerfen von Inline-Code verschluckt
echte Zitate) wurde beim zweiten Anlauf in einer anderen Form wiederholt
(punkt-suffigierte Kennungen), weil die Korrektur wieder als Handnäherung statt
als Produkt-Aufruf lief. Aufgelöst durch einen temporären Go-Test, der die
echte Funktion gegen alle 85 `archiv.zip`-Archive fuhr.

**Ein Wellen-Archivierungs-Commit passt nicht in `AGENTS.md` §3.3s bestehende
Zwei-Commit-Form** — es gibt keine Phase, in der die bewegte Datei ihren
Inhalt unverändert behält, weil der Stub den Volltext im selben Akt ersetzt,
der ihn verschiebt. Das war in slice-190 als Folgepunkt für slice-191
angekündigt, aber erst nach einem Review-Fund tatsächlich als vierte
`AGENTS.md`-§3.3-Ausnahme ([`MR-059`](../../../../harness/conventions.md#mr-059--wellen-archiv-stub-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)) geschrieben.

## Steering-Loop-Einträge

- **[`BEO-011`](../observations.md) sechste Instanz** (slice-190/191) — **weiter
  offen**, nicht verkörpert: eine einmalig sorgfältig gefahrene Prozedur
  (Welle-für-Welle mit Gate-Prüfung) ist keine Verkörperung ohne Zielort samt
  Herkunfts-Anker für künftige Läufe. Bleibt im Register offen.
- **[`MR-059`](../../../../harness/conventions.md#mr-059--wellen-archiv-stub-move-ist-ein-einziger-deklarierter-commit-nachtrag-zu-mr-013)**
  neu geschrieben, `AGENTS.md` §3.3 trägt jetzt eine vierte Ausnahme dafür.
- **Drei Prozedur-Lehren, in slice-191 §9 festgehalten** (keine gesonderten
  Regel-Dateien): (1) ein Archivierungswerkzeug muss überlebende Kennungen
  aktiv ins Stub-Feld übernehmen, sonst bricht `--require-complete` lautlos;
  (2) eine Kennungs-Extraktion muss echte Link-Label-Zitate von
  illustrativem Inline-Code unterscheiden; (3) eine Mess-**Korrektur** nutzt
  das getestete Produkt direkt, keine weitere Handnäherung seiner Logik.

## Beobachtungs-Register (Zeiger)

Lese-Schritt über die Bewegungen dieser Welle:

| Eintrag | Stand | Was daraus folgt |
|---|---|---|
| [`BEO-011`](../observations.md) | 5 → **6** | Sechste Instanz (slice-190/191) — Ausgang **weiter offen**, keine Verkörperung ohne Zielort samt Herkunfts-Anker |

Der Zähler steht in [`observations.md`](../observations.md). Kein weiterer
Eintrag erreichte während dieser Welle die 3×-Schwelle neu.

## Folge-Slices

Keine — `welle-87` schließt vollständig innerhalb ihrer zwei Slices, ohne
offenen Rest. `welle-86` bleibt eigenständig (siehe §5/§6 des Welle-Plans) und
profitiert vom hier gebauten und bewiesenen Werkzeug bei ihrer eigenen
Closure.

## Verifikation

- `make gates` Exit 0 (zehn Glieder) auf jedem der über 90 Commits dieser
  Welle.
- `make fullbuild` Exit 0 — 50 Requirements, 0 Waisen — auf dem vollständig
  archivierten Bestand (image-hash `sha256:fffe5ea422a3b55021c5b243faf0c72d7dce7773b2e65fdcde59352e46cca3af`).
- `make archive-wave-test` grün, inklusive fünf dokumentierter Umkehr-Proben
  über die Lebensdauer der Welle.
- Ein temporärer, nicht committeter Audit-Lauf (`ExtractSurvivingIDs` gegen
  jedes `archiv.zip` aller 85 Wellen) meldete zuletzt 0 Abweichungen.
- Ein unabhängiger Review und zwei unabhängige Verifikationsrunden je Slice,
  alle Befunde eingearbeitet.
