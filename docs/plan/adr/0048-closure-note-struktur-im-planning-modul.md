# ADR-0048 — Closure-Note-Struktur als zweite `planning`-Fähigkeit + Konfigurations-Pfad-Flag

**Status:** Accepted
**Datum:** 2026-08-09
**Autor:** pt9912
**Bezug:** [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
(Modul `planning`, um die Closure-Fähigkeit geschärft),
[`DC-FA-CLI-012`](../../../spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben)
(neu), [`DC-FA-CONF-001`](../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(Pfad als Konvention statt Zwang); Modul-Fundament
[ADR-0028](0028-planning-lifecycle-modul.md) (das `planning`-Modul selbst),
[ADR-0005](0005-modul-layout-hexagon-ordner.md),
[ADR-0012](0012-kern-paketschnitt-model-rules-app.md); Bindepunkt-Vorbild
[ADR-0026](0026-completeness-in-product-gate.md) (Closure-Gate als
in-Produkt-Flag außerhalb von `gates`).
**Schärft:** die Erweiterung von
[`spec/spezifikation.md` §DC-FA-PLAN-001.a](../../../spec/spezifikation.md#dc-fa-plan-001a--planning-lifecycle-konsistenz-planning)
(Schritte C1–C5) und die neue Sektion
[§DC-FA-CLI-012.a](../../../spec/spezifikation.md#dc-fa-cli-012a--konfigurations-pfad---config).

## Kontext

Die adoptierte Baseline vendored **zwei** Formen für den
Closure-Note-Qualitäts-Nachlauf: ein *computational* Struktur-Gate und den
*inferentiellen*
[`closure-note-reviewer.template.md`](../../../.harness/baseline/v5.0.0/templates/.harness/skills/closure-note-reviewer.template.md)
darüber. Das Regelwerk begründet die Zweiteilung ausdrücklich damit, dass
Floskeln wie „war ganz okay, läuft jetzt" **syntaktisch** zwei Sätze sind
([`modul-11-verification.md`](../../../.harness/baseline/v5.0.0/regelwerk/modul-11-verification.md)
§Fitness Function ohne Standard-Tool).

d-check hatte **keine** der beiden Formen — obwohl es die Closure-Notiz selbst zur
Pflicht erklärt und in 92 abgeschlossenen Slices führt. Die Pflicht stand damit
ohne jede maschinelle Entsprechung; ein zurückgelassener `_Ausstehend._`-Platzhalter
wäre durch alle Gates gelaufen.

Eine Messung gegen den eigenen Bestand (2026-08-09) rahmt den Entwurf:

- **92 von 92** abgeschlossenen Slices tragen einen Closure-Notiz-Abschnitt; die
  Überschrift variiert nur in Nummer und Suffix — ein Muster deckt alle.
- Das **Minimum** an Satzende-Zeichen außerhalb Code-Blöcken ist **5**. Ein
  Struktur-Gate braucht also **keinen** Bestands-Retrofit.
- Der reale Fehlerfall hat trotzdem Zähne: der Platzhalter eines noch offenen
  Slice zählt **1**.

## Entscheidung

1. **Kein neues Modul — `planning` bekommt eine zweite Fähigkeit.** Die
   Aktiv-Status-Invariante und die Closure-Notiz sind **dieselbe**
   Lifecycle-Invariante von zwei Seiten: der Eintritt in den Lifecycle (Slice
   liegt in Arbeit ⇔ Roadmap benennt die Welle) und der Austritt (Slice ist
   abgeschlossen ⇒ er trägt eine Closure-Notiz). Ein zweites Modul hätte
   dieselbe Config-Achse (Slice-Verzeichnis, `slice-glob`) ein zweites Mal
   deklariert. Verworfen wurde damit auch das im Baseline-Template gezeigte
   `check_closure_notes.py` — d-checks Identität ist gerade die Ablösung
   kopierter `tools/*.sh`-Skripte durch verteilbare Go-Regelmodule
   ([ADR-0028](0028-planning-lifecycle-modul.md) hat als letztes der Familie
   genau diesen Schritt gemacht; ein neues Skript wäre eine Rolle rückwärts).

2. **Opt-in innerhalb des opt-in Moduls.** Die Fähigkeit ist an
   `planning.closure.dir` gebunden. Ohne den Schlüssel wird **keine** Slice-Datei
   geöffnet und der Befundsatz ist byte-identisch
   ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)) — dasselbe
   Muster wie `codepaths.check-lines`
   ([ADR-0045](0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md)).

3. **Drei Grund-Codes statt einem.** `closure-note-missing` (Abschnitt bzw.
   gesetztes Verzeichnis fehlt), `closure-note-thin` (zu wenig Substanz),
   `closure-note-boilerplate` (Floskel-Treffer) — wie `targets` zwei Richtungen
   getrennt meldet ([ADR-0031](0031-targets-deklarations-konsistenz-modul.md)).
   Ein Sammel-Code hätte drei verschiedene Reparaturen unter eine Meldung
   gelegt.

4. **Struktur, nicht Bedeutung — als Vertragsgrenze.** Zugesagt ist ausschließlich
   die strukturelle Prüfung. Die Grenze steht ausdrücklich in
   [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
   §Out-of-Scope, damit ein grünes Gate nicht als „Closure-Notizen sind gut"
   missverstanden wird. Die semantische Schicht ist der Reviewer-Skill; er
   meldet nicht doppelt, was die Struktur schon abdeckt.

5. **Floskel-Liste per Default leer.** Der Vertrag bringt **keine**
   sprach-spezifischen Phrasen mit; das adoptierende Repo deklariert seine
   eigenen. Anders als bei den Modalitäts-Keywords
   ([ADR-0036](0036-trace-modality-klassifikation.md)), wo ein DE/EN-Built-in die
   *Klassifikation* trägt, entscheidet eine Floskel-Liste über **rot oder grün**
   — ein mitgeliefertes deutsches Vokabular wäre in einem englischen Repo
   entweder wirkungslos oder falsch-positiv. Ein leerer Listen-Eintrag ist
   fail-closed (Exit 2), weil er jeden Text träfe.

6. **Schwelle 4 statt der Baseline-2.** Das Regelwerk nennt „mindestens zwei
   Sätze" als Beispiel einer operationalisierbaren Aussage, nicht als Zahl-Zusage.
   Der gemessene Bestand liegt bei mindestens 5; eine Schwelle von 4 bleibt
   darunter (kein Retrofit) und ist zugleich mehr als der Platzhalter-Boden. Ein
   **späteres Anheben** ist frei, ein **Senken** ist nach
   [`AGENTS.md` §3.6](../../../AGENTS.md#36-gates-dürfen-nicht-ohne-adr-gelockert-werden)
   ADR-pflichtig — die Richtung der Reibung zeigt damit auf die strengere Seite.
   `min-sentences` < 1 ist Exit 2 (eine Null-Schwelle wäre ein stilles Grün).

7. **Eigener Bindepunkt statt `gates` — und deshalb `--config`.** Eine
   Closure-Frage gehört nicht in den inneren Loop
   ([`modul-11-verification.md`](../../../.harness/baseline/v5.0.0/regelwerk/modul-11-verification.md)
   §Fitness Function; d-check hat dieselbe Trennung schon bei
   [ADR-0026](0026-completeness-in-product-gate.md) gezogen: das
   Vollständigkeits-Gate hängt an der Closure, nicht an `gates`). Da die
   Konfiguration konventionell aus **einer** Datei in der Scan-Wurzel kommt, liefe
   alles, was im `planning`-Modul wohnt, automatisch dort mit, wo das Modul läuft
   — also in `gates`. Die Trennung braucht deshalb eine **Herkunfts**-Umschaltung:
   `--config <datei>` erlaubt zwei disjunkte Prüf-Profile im selben Repo.
   Bewusst **keine** Semantik-Änderung: der Schalter verschiebt nur, **woher** die
   Konfiguration kommt; er ersetzt (statt zu ergänzen), bleibt innerhalb der
   Scan-Wurzel (read-only-Mount-Grenze) und fällt bei fehlender Datei **nicht**
   still auf Defaults zurück — ein vertipptes Profil darf keinen anderen
   Prüfumfang fahren.

8. **Fail-closed auch bei null Kandidaten — und geteilte Heading-Lexik.** Zwei
   Nachschärfungen aus dem Review, beide gegen dieselbe Klasse (stilles Grün):
   (a) Ein existierendes, aber **kandidatenfreies** Closure-Verzeichnis meldet
   `closure-note-missing` statt zu schweigen. Den Schlüssel zu setzen **ist** die
   Behauptung, dass dort Notizen liegen; zieht der Bestand in Unterordner um,
   liefe das Gate sonst fortan leer und grün — dieselbe Nullmengen-Logik, die
   [ADR-0037](0037-trace-tabellenquellen-nullmengen-guard.md) für die
   RTM-Anforderungsquellen gezogen hat. Ein Repo ohne abgeschlossene Slices setzt
   den Schlüssel schlicht noch nicht.
   (b) Was als Überschrift zählt, entscheidet der **geteilte** ATX-Parser des
   Pakets (derselbe, den `anchors`/`matrix` nutzen) — nicht eine eigene
   `#`-Zählung. Eine Zeile wie `#1 war ein Thema` ist Fließtext; die eigene
   Heuristik hätte sie als H1 gelesen, den Abschnitt dort abgeschnitten und alles
   dahinter unsichtbar gemacht. Dieselbe Lexik an beiden Abschnitts-Grenzen ist
   damit konstitutiv, nicht kosmetisch.

## Verglichene Alternativen

| Alternative | Warum verworfen |
|---|---|
| Neues Regelmodul `closure` | Verdoppelt die Slice-Verzeichnis-/`slice-glob`-Achse; die Invariante ist dieselbe wie in `planning`, nur die andere Lifecycle-Seite |
| Ein `check_closure_notes.py`-Skript wie im Baseline-Template | Ein kopiertes Skript ist genau die Form, die d-check ablöst ([ADR-0028](0028-planning-lifecycle-modul.md)); nicht verteilbar, nicht dogfood-fähig, kein Konsumenten-Nutzen |
| Closure-Prüfung in `gates` mitlaufen lassen (ohne `--config`) | Billiger, aber vermischt Inner-Loop und Closure; die Baseline trennt beide ausdrücklich, und d-check hat die Trennung mit `completeness-check` schon etabliert |
| Ein Sammel-Grund-Code `closure-note-drift` | Drei verschiedene Reparaturen unter einer Meldung; `targets` zeigt den Gegenentwurf |
| Deutsche Floskel-Defaults ausliefern | Entscheidet über rot/grün und wäre in fremdsprachigen Adopter-Repos wirkungslos oder falsch-positiv |
| Semantische Floskel-Erkennung im Produkt | d-check ist ein deterministischer Struktur-Checker; „Inhalt vs. Floskel" ist inferentiell und gehört in den Reviewer-Skill |
| Config-**Merge** statt Ersetzen bei `--config` | Zwei Dateien, die sich überlagern, machen den effektiven Prüfumfang unlesbar; [`DC-FA-CONF-001`](../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei) schließt Vererbung schon aus |

## Konsequenzen

- Das `planning`-Modul liest ab jetzt **Slice-Inhalte** — bisher zählte nur die
  Datei-Existenz. Die Out-of-Scope-Zeile von
  [`DC-FA-PLAN-001`](../../../spec/lastenheft.md#dc-fa-plan-001--planning-lifecycle-konsistenz-modul-planning-opt-in)
  ist entsprechend präzisiert: Inhalt liest ausschließlich die opt-in
  Closure-Fähigkeit, und nur im Closure-Verzeichnis.
- `--config` ist eine **neue öffentliche CLI-Fläche**. Jede Stelle, die heute den
  konventionellen Dateinamen als Befund-Provenance nennt, muss die *tatsächlich
  geladene* Datei nennen — sonst zeigt eine Meldung auf eine Datei, die der Lauf
  nie gelesen hat.
- Zwei Profil-Dateien im Repo bedeuten zwei Pflegestellen. Die Netzlos-Modullisten-
  Integrität ([ADR-0032](0032-gate-consistency-tombstone.md)) muss beide abdecken,
  sonst entsteht eine ungeprüfte zweite Tür.
- Der Bestand bleibt ohne Retrofit grün (gemessenes Minimum 5 ≥ Schwelle 4). Das
  ist ein **Boden, keine Decke**: die Schwelle sagt nichts über Qualität, nur
  über Substanz-Untergrenze.
- Konsumenten erben die Fähigkeit über das verteilte Image, ohne sie zu
  aktivieren (Default inert).

## Fitness Function

- **Negativ-Selbsttest** als Akzeptanztest im Modul: ein Slice mit Platzhalter
  ⇒ `closure-note-thin`; ein Slice ohne passenden Abschnitt ⇒
  `closure-note-missing`; ein konfigurierter Floskel-Treffer ⇒
  `closure-note-boilerplate` — jeweils Exit 1, kein stilles Grün.
- **Code-Block-Boundary:** eine Notiz, deren Satzende-Zeichen überwiegend in
  einem Fenced-Block stehen, bleibt unter der Schwelle und meldet
  `closure-note-thin` (die Bereinigung ist wirksam, nicht dekorativ).
- **Byte-Identität:** ohne `planning.closure.dir` ist der Befundsatz identisch
  zum Stand vor dieser Entscheidung, und keine Slice-Datei wird geöffnet.
- **`--config`-fail-closed:** fehlende Datei ⇒ Exit 2 (kein Rückfall); Pfad
  außerhalb der Scan-Wurzel ⇒ Exit 2.
- **Dogfood:** d-check fährt die Prüfung über den eigenen `done/`-Bestand.

## Re-Evaluierungs-Trigger

- Wenn eine Floskel-Art dreimal als Reviewer-Befund auftritt, gehört sie in die
  `boilerplate`-Liste — und die Frage zurück auf den Tisch, ob die
  Struktur-Schicht mehr tragen kann, als hier zugesagt ist.
- Wenn `--config` über den Closure-Bindepunkt hinaus zum Profil-Mechanismus wird
  (mehrere benannte Profile, Vererbung), ist das eine neue Anforderung — nicht
  eine Ausweitung dieser Entscheidung.
- Wenn die Schwelle 4 sich als zu grob oder zu fein erweist (Realdaten aus
  Konsumenten-Repos), ist ein **Anheben** frei; ein Senken bleibt ADR-pflichtig.
- Sollte ein zweites Modul dieselbe Closure-Achse brauchen, ist die
  „kein neues Modul"-Entscheidung neu zu bewerten.

## Geschichte

- 2026-08-09: Proposed (doc-first, `slice-093`).
- 2026-08-09: Entscheidung 8 nach zwei unabhängigen Frischkontext-Reviews
  ergänzt (Nullmengen-Guard + geteilte Heading-Lexik); die Wurzel-Grenze von
  `--config` wurde von lexikalisch auf symlink-fest nachgezogen und ein leerer
  Flag-Wert als Nutzungsfehler ausgewiesen.
- 2026-08-09: Accepted (Closure `slice-093`, Release **v0.52.0**).
