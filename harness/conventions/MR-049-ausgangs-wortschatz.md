# MR-049 — Der Ausgang eines Risikos trägt eines der drei Wörter

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine. Der Kanon benennt die urteilsfreie Hälfte
  der Drei-Ausgänge-Regel ausdrücklich — *dass* zu jedem Risiko ein Ausgang   <!-- d-check:cite .harness/baseline/v5.15.0/regelwerk/modul-05-planning-harness.md:146-147 -->
  dasteht **und welcher der drei** es ist — und schließt mit: „*Welches* Werkzeug
  die urteilsfreie Hälfte prüft, ist Repo-Entscheidung; dass sie eine hat, ist
  es nicht."
  ([`modul-05-planning-harness.md`](../../.harness/baseline/v5.15.0/regelwerk/modul-05-planning-harness.md)). Die Auszeichnung folgt hier der **Quelle**, nicht der Haus-Form `*„…"*` — sie fasst das erste Wort kursiv, und eine Kursiv-Klammer darum verschluckte es. Das Zitat ist per Direktive gebunden; wer die Klammer zurücksetzt, macht das Gate rot.
  Diese Adaption ist die Werkzeug-Wahl, nicht eine Abweichung.
- **Datum:** 2026-08-27
- **Geltungsbereich:** [`.d-check.closure.yml`](../../.d-check.closure.yml),
  Abschnitt §5 der Slice-Dateien unter `docs/plan/planning/done/`.
- **Adaption:** Der Ausgangs-Wortschatz ist **geschlossen** und umfasst alle
  drei Kanon-Ausgänge: `eingetreten` · `entfallen` · `weiter offen`. Eine
  `structure`-Regel im Closure-Profil hält ihn; sie läuft in
  `make verify-closure-notes`.

  **Sie ergänzt die vorhandene Regel, sie ersetzt sie nicht.** Die
  Platzhalter-Regel fängt den **vergessenen** Ausgang, diese den **erfundenen** —
  die Gestalt aus [`BEO-015`](../../docs/plan/planning/observations.md).

  **Die Form ist ein `forbid-pattern` über das Komplement der drei Wörter**, kein
  `require-pattern` über ihr Vorkommen. Der Unterschied entscheidet die
  Reichweite: `forbid-pattern` ist über **jedes Vorkommen** des Markers
  quantifiziert und trifft damit **je Risiko**, während ein `require-pattern`
  nur sagt, dass *irgendwo* im Abschnitt ein erlaubtes Wort steht — ein §5 mit
  einem kanonischen und einem Freitext-Ausgang liefe darunter grün durch.
- **Grenzen — gemessen, nicht geschätzt.**

  **Der Lookahead-Weg ist versperrt, der Komplement-Weg nicht.** RE2 weist
  `(?!` mit *„invalid or unsupported Perl syntax"* ab. Das schließt aber nur
  eine Schreibweise aus: reguläre Sprachen sind unter Komplement abgeschlossen,
  und das Komplement einer endlichen Wortmenge ist als Präfix-Alternation
  darstellbar. Genau die steht in der Regel — bis **Tiefe fünf**. Eine
  Abweichung erst ab dem sechsten Zeichen (`eingetreXen`) läuft durch.

  **Drei Formen sieht die Regel nicht.** Ein Risiko **ganz ohne** Marker liefert
  kein Vorkommen und damit kein Match — im Bestand betrifft das `slice-106`,
  `slice-110` und `slice-111`, deren §5 Risiken führt und keinen einzigen   <!-- d-check:cite .harness/baseline/v5.15.0/regelwerk/modul-05-planning-harness.md:136-136 -->
  Ausgang; das ist der Kanon-Kernsatz *„Ein Slice geht nicht nach `done/`,
  während ein Risiko ohne Ausgang dasteht"*, dreimal verletzt und von **keiner**
  der beiden Regeln erreichbar. Ein Ausgangswort in **Inline-Code** wird von der
  Abschnitts-Bereinigung geleert, und die Regel meldet dann **falsch**. Und der
  Befund zeigt auf die **Abschnitts-Überschrift**, nicht auf das verletzende
  Risiko — bei einer Abschnitts-Regel geht es nicht anders.

  **Der Altbestand ist ausgenommen, und die Ausnahme zeigt nach hinten.**
  Ohne sie meldet die Regel **122** Befunde: **107** `section-missing`, weil die
  älteren Slices ihr §5 unter einem anderen Titel führen (`## 5. Trigger` 43 ·
  `## 5. Closure-Trigger` 31 · `## 5. Risiken / offene Punkte` 19 · weitere) —
  eine `## 5.`-Überschrift tragen **alle 150** —, plus **15** Freitext-Ausgänge.
  Ein Retrofit über 122 Dateien hat niemand beschlossen, und der Kurs hat die
  Ausweitung ausdrücklich als Repo-Entscheidung und **kein** Konformitätsthema
  bezeichnet.

  Die Ausnahme nennt deshalb den **abgeschlossenen** Altbestand mit **fester
  Ziffernzahl** (`slice-0??-*`, `slice-1[0-3]?-*`). Die naheliegende Kurzform
  `slice-1[0-3]*` nähme zusätzlich `slice-1000` bis `slice-1399` heraus:
  `matchGlob` matcht segmentweise mit `path.Match`, und `*` frisst dort den
  ganzen Rest. Ein Glob auf die **neuen** Nummern (`slice-1[4-9][0-9]`) hätte
  denselben Fehler mit näherem Horizont — er hörte bei `slice-200` auf. Beides
  ist dieselbe Klasse lautlosen Ausfalls; die gewählte Form hat sie nicht.
- **Begründung:** Der Wortschatz war offen, und das ist gemessen: über alle 144
  Ausgangs-Vorkommen im Bestand stehen **28** verschiedene Anfangswörter,
  darunter `behoben`, `gemessen`, `gehalten`, `erledigt`, `benannt`,
  `aufgelöst` — und blanker Freitext. Jedes einzelne liest sich vernünftig;
  zusammen ist die geschlossene Menge des Kanons keine mehr, und ein vierter
  Ausgang fällt nicht auf, weil er wie der fünfte klingt.
- **Auflösungs-Trigger:** das Produkt kann eine Bedingung **je Listen-Eintrag**
  ausdrücken. Dann ist auch das Risiko **ohne** Marker erreichbar — die Lücke,
  die diese Regel bauartbedingt offen lässt — und ein Nachfolge-Eintrag löst
  sie ab.
